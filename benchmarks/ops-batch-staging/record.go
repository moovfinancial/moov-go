package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func recordValidation(dbPath string, current sampleRange, options validateOptions, userAgent string, result requestResult, status string, includeResults bool) error {
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
	statements.WriteString("DROP TABLE current_run;\nCOMMIT;\n")
	if _, err := runSQLite(dbPath, statements.String()); err != nil {
		return fmt.Errorf("recording validation: %w", err)
	}
	return nil
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
	PlannedSamples  int `json:"planned_samples"`
	ValidatedRanges int `json:"validated_ranges"`
	Attempted       int `json:"attempted"`
	Matched         int `json:"matched"`
	Mismatched      int `json:"mismatched"`
	Errored         int `json:"errored"`
	FailedRequests  int `json:"failed_requests"`
}

func printSummary(dbPath string) error {
	if err := requireCampaignDB(dbPath); err != nil {
		return err
	}
	query := `
SELECT
  (SELECT COUNT(*) FROM samples) AS planned_samples,
  (SELECT COUNT(*) FROM validation_runs WHERE status IN ('clean', 'mismatch', 'mismatch_error')) AS validated_ranges,
  (SELECT COUNT(*) FROM validation_results) AS attempted,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'match') AS matched,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'mismatch') AS mismatched,
  (SELECT COUNT(*) FROM validation_results WHERE comparison = 'error') AS errored,
  (SELECT COUNT(*) FROM validation_runs WHERE status IN ('request_error', 'http_error')) AS failed_requests;
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
	fmt.Printf("planned=%d validatedRanges=%d attempted=%d matched=%d mismatched=%d errored=%d failedRequests=%d\n",
		current.PlannedSamples, current.ValidatedRanges, current.Attempted,
		current.Matched, current.Mismatched, current.Errored, current.FailedRequests)
	return nil
}
