package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func recordValidation(dbPath string, current sampleRange, options validateOptions, userAgent string, result requestResult, status string, includeResults bool) (int64, error) {
	startedAt := result.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}
	var statements strings.Builder
	statements.WriteString("PRAGMA foreign_keys = ON;\nBEGIN;\n")
	fmt.Fprintf(&statements, `INSERT INTO validation_runs (
  range_start, range_end, started_at, api_version, service_version, user_agent,
  request_id, trace_id, duration_ms, http_status, attempted_count, status
) VALUES (%d, %d, %s, %s, %s, %s, %s, '', %d, %d, %d, %s);
CREATE TEMP TABLE current_run (id INTEGER);
INSERT INTO current_run VALUES (last_insert_rowid());
`, current.Start, current.End, sqlQuote(startedAt.Format(time.RFC3339Nano)),
		sqlQuote(options.Version), sqlQuote(options.ServiceVersion), sqlQuote(userAgent),
		sqlQuote(result.RequestID), result.Duration.Milliseconds(), result.StatusCode,
		len(current.Samples), sqlQuote(status))

	if includeResults {
		appendResultStatements(&statements, current.Samples, result.Response)
	}
	statements.WriteString("SELECT id FROM current_run;\nDROP TABLE current_run;\nCOMMIT;\n")
	output, err := runSQLite(dbPath, statements.String(), "-noheader")
	if err != nil {
		return 0, fmt.Errorf("recording validation: %w", err)
	}
	runID, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("reading recorded validation ID: %w", err)
	}
	return runID, nil
}

func appendResultStatements(statements *strings.Builder, samples []sample, response validationResponse) {
	diffs := make(map[string]struct{}, len(response.TransferIDsWithDiff))
	for _, transferID := range response.TransferIDsWithDiff {
		diffs[transferID] = struct{}{}
	}
	for _, current := range samples {
		comparison, resultErr := sampleComparison(current.TransferID, diffs, response.ErrorsByTransferID)
		fmt.Fprintf(statements, `INSERT INTO validation_results (run_id, pk_id, comparison, diff_paths, error)
SELECT id, %d, %s, '', %s FROM current_run;
`, current.PKID, sqlQuote(comparison), sqlQuote(resultErr))
	}
}

func sampleComparison(transferID string, diffs map[string]struct{}, errorsByTransferID map[string]string) (string, string) {
	if resultErr, ok := errorsByTransferID[transferID]; ok {
		return "error", resultErr
	}
	if _, ok := diffs[transferID]; ok {
		return "mismatch", ""
	}
	return "match", ""
}

type summary struct {
	CachedThroughPK int64 `json:"cached_through_pk"`
	PlannedSamples  int   `json:"planned_samples"`
	ValidatedRanges int   `json:"validated_ranges"`
	Attempted       int   `json:"attempted"`
	Matched         int   `json:"matched"`
	Mismatched      int   `json:"mismatched"`
	Errored         int   `json:"errored"`
	FailedRequests  int   `json:"failed_requests"`
	SkippedRanges   int   `json:"skipped_ranges"`
}

func printSummary(dbPath string) error {
	if err := requireCampaignDB(dbPath); err != nil {
		return err
	}
	query := `
SELECT
  (SELECT cached_through_pk FROM sample_cache_state WHERE id = 1) AS cached_through_pk,
  (SELECT COUNT(*) FROM samples) AS planned_samples,
  (SELECT COUNT(*) FROM validation_runs WHERE status IN ('clean', 'mismatch', 'mismatch_error')) AS validated_ranges,
  (SELECT COUNT(*) FROM validation_results) AS attempted,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'match') AS matched,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'mismatch') AS mismatched,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'error') AS errored,
  (SELECT COUNT(*) FROM validation_runs WHERE status IN ('request_error', 'http_error', 'evidence_error')) AS failed_requests,
  (SELECT COUNT(*) FROM monitor_skipped_ranges) AS skipped_ranges;
`
	output, err := runSQLite(dbPath, query, "-json")
	if err != nil {
		return fmt.Errorf("reading summary: %w", err)
	}
	var summaries []summary
	if err := json.Unmarshal(output, &summaries); err != nil {
		return fmt.Errorf("decoding summary: %w", err)
	}
	if len(summaries) != 1 {
		return fmt.Errorf("expected one summary row, got %d", len(summaries))
	}
	current := summaries[0]
	fmt.Printf("cachedThroughPK=%d planned=%d validatedRanges=%d skippedRanges=%d attempted=%d matched=%d mismatched=%d errored=%d failedRequests=%d\n",
		current.CachedThroughPK, current.PlannedSamples, current.ValidatedRanges, current.SkippedRanges, current.Attempted,
		current.Matched, current.Mismatched, current.Errored, current.FailedRequests)
	return nil
}

func recordValidationTrace(dbPath string, runID int64, traceID string) error {
	if runID == 0 {
		return errors.New("validation run ID is required to record Honeycomb evidence")
	}
	output, err := runSQLite(dbPath, fmt.Sprintf(`
UPDATE validation_runs SET trace_id = %s WHERE id = %d;
SELECT changes();
`, sqlQuote(traceID), runID), "-noheader")
	if err != nil {
		return fmt.Errorf("recording validation evidence: %w", err)
	}
	if strings.TrimSpace(string(output)) != "1" {
		return fmt.Errorf("recording validation evidence updated %s rows", strings.TrimSpace(string(output)))
	}
	return nil
}

func markValidationEvidenceError(dbPath string, runID int64) error {
	output, err := runSQLite(dbPath, fmt.Sprintf(`
UPDATE validation_runs SET status = 'evidence_error' WHERE id = %d AND status = 'evidence_pending';
SELECT changes();
`, runID), "-noheader")
	if err != nil {
		return fmt.Errorf("recording evidence failure: %w", err)
	}
	if strings.TrimSpace(string(output)) != "1" {
		return fmt.Errorf("recording evidence failure updated %s rows", strings.TrimSpace(string(output)))
	}
	return nil
}

func finalizeMonitorSuccess(dbPath string, runID int64, traceID string, selected sampleRange, safeCheckpoint int64) error {
	output, err := runSQLite(dbPath, fmt.Sprintf(`
BEGIN;
UPDATE validation_runs SET status = 'clean', trace_id = %s
WHERE id = %d AND status = 'evidence_pending';
SELECT changes();
INSERT OR IGNORE INTO monitor_skipped_ranges (range_start, range_end, skipped_at, safe_checkpoint, reason)
SELECT s.range_start, s.range_end, %s, %d, 'sparse-monitor-sampling'
FROM samples s
WHERE s.range_end <= %d AND s.range_start < %d
  AND (SELECT status FROM validation_runs WHERE id = %d) = 'clean'
  AND NOT EXISTS (
    SELECT 1 FROM validation_runs vr
    WHERE vr.range_start = s.range_start AND vr.range_end = s.range_end
  )
GROUP BY s.range_start, s.range_end;
COMMIT;
`, sqlQuote(traceID), runID, sqlQuote(time.Now().UTC().Format(time.RFC3339Nano)), safeCheckpoint, safeCheckpoint, selected.Start, runID), "-noheader")
	if err != nil {
		return fmt.Errorf("finalizing monitor success: %w", err)
	}
	if strings.TrimSpace(string(output)) != "1" {
		return fmt.Errorf("finalizing monitor success updated %s validation rows", strings.TrimSpace(string(output)))
	}
	return nil
}
