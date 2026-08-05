package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const stagingAPIHost = "api.moov-staging.io"

type validateOptions struct {
	DBPath          string
	Checkpoint      int64
	Version         string
	ServiceVersion  string
	MaxRanges       int
	UserAgentSuffix string
	EvidencePending bool
}

type sample struct {
	PKID            int64  `json:"pk_id"`
	TransferID      string `json:"transfer_id"`
	RangeStart      int64  `json:"range_start"`
	RangeEnd        int64  `json:"range_end"`
	SelectionReason string `json:"selection_reason"`
}

type sampleRange struct {
	Start   int64
	End     int64
	Samples []sample
}

type validationResponse struct {
	TransferIDsWithDiff []string          `json:"transferIDsWithDiff"`
	ErrorsByTransferID  map[string]string `json:"errorsByTransferID"`
}

type requestResult struct {
	Response   validationResponse
	StartedAt  time.Time
	StatusCode int
	RequestID  string
	Duration   time.Duration
}

type validationOutcome struct {
	Status    string
	UserAgent string
	Result    requestResult
	Attempted int
	RunID     int64
}

func validatePending(options validateOptions) error {
	if err := requireCampaignDB(options.DBPath); err != nil {
		return err
	}
	ranges, err := pendingRanges(options.DBPath, options.Checkpoint, options.MaxRanges)
	if err != nil {
		return err
	}
	if len(ranges) == 0 {
		fmt.Printf("no pending samples at checkpoint %d\n", options.Checkpoint)
		return nil
	}

	credentials, err := stagingCredentials()
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 120 * time.Second}
	for _, current := range ranges {
		if _, err := validateRange(context.Background(), client, credentials, options, current); err != nil {
			return err
		}
	}
	return nil
}

type credentials struct {
	Host       string
	Username   string
	Credential string
}

func stagingCredentials() (credentials, error) {
	values := credentials{
		Host:       os.Getenv("MOOV_HOST"),
		Username:   os.Getenv("MOOV_PUBLIC_KEY"),
		Credential: os.Getenv("MOOV_SECRET_KEY"),
	}
	if values.Host == "" || values.Username == "" || values.Credential == "" {
		return credentials{}, errors.New("MOOV_HOST, MOOV_PUBLIC_KEY, and MOOV_SECRET_KEY must be set; run moov_env staging --quiet")
	}
	if values.Host != stagingAPIHost {
		return credentials{}, fmt.Errorf("refusing to call host %q; expected %s", values.Host, stagingAPIHost)
	}
	return values, nil
}

func validateRange(ctx context.Context, client *http.Client, auth credentials, options validateOptions, current sampleRange) (validationOutcome, error) {
	userAgent := fmt.Sprintf("moov-go-projection-backfill-validator/1 range-%d-%d", current.Start, current.End)
	if options.UserAgentSuffix != "" {
		userAgent += " " + options.UserAgentSuffix
	}
	result, requestErr := requestValidation(ctx, client, auth, options.Version, userAgent, current.Samples)
	status := validationStatus(result, requestErr)
	recordedStatus := status
	if options.EvidencePending && status == "clean" {
		recordedStatus = "evidence_pending"
	}
	runID, err := recordValidation(options.DBPath, current, options, userAgent, result, recordedStatus, requestErr == nil)
	if err != nil {
		return validationOutcome{}, err
	}
	outcome := validationOutcome{Status: status, UserAgent: userAgent, Result: result, Attempted: len(current.Samples), RunID: runID}

	fmt.Printf("range=%d-%d attempted=%d status=%d duration=%s requestID=%s result=%s diffs=%d errors=%d\n",
		current.Start, current.End, len(current.Samples), result.StatusCode,
		result.Duration.Round(time.Millisecond), result.RequestID, status,
		len(result.Response.TransferIDsWithDiff), len(result.Response.ErrorsByTransferID))
	if requestErr != nil {
		return outcome, requestErr
	}
	if status != "clean" {
		return outcome, fmt.Errorf("range %d-%d requires investigation: %s", current.Start, current.End, status)
	}
	return outcome, nil
}

func requestValidation(ctx context.Context, client *http.Client, auth credentials, version, userAgent string, samples []sample) (requestResult, error) {
	transferIDs := make([]string, 0, len(samples))
	for _, current := range samples {
		transferIDs = append(transferIDs, current.TransferID)
	}
	body, err := json.Marshal(map[string]any{"transferIDs": transferIDs})
	if err != nil {
		return requestResult{}, fmt.Errorf("encoding request: %w", err)
	}

	endpoint := "https://" + auth.Host + "/ops/validate-api-projection-transfers"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return requestResult{}, fmt.Errorf("creating request: %w", err)
	}
	req.SetBasicAuth(auth.Username, auth.Credential)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Moov-Version", version)
	req.Header.Set("User-Agent", userAgent)

	started := time.Now().UTC()
	resp, err := client.Do(req)
	result := requestResult{StartedAt: started, Duration: time.Since(started)}
	if err != nil {
		return result, fmt.Errorf("calling projection validator: %w", err)
	}
	defer resp.Body.Close()
	result.StatusCode = resp.StatusCode
	result.RequestID = resp.Header.Get("X-Request-ID")

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	result.Duration = time.Since(started)
	if err != nil {
		return result, fmt.Errorf("reading projection validator response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("projection validator returned HTTP %d (request %s)", resp.StatusCode, result.RequestID)
	}
	if err := json.Unmarshal(responseBody, &result.Response); err != nil {
		return result, fmt.Errorf("decoding projection validator response: %w", err)
	}
	result.Duration = time.Since(started)
	return result, nil
}

func validationStatus(result requestResult, requestErr error) string {
	if requestErr != nil {
		if result.StatusCode == 0 {
			return "request_error"
		}
		return "http_error"
	}
	if len(result.Response.TransferIDsWithDiff) > 0 {
		if len(result.Response.ErrorsByTransferID) > 0 {
			return "mismatch_error"
		}
		return "mismatch"
	}
	if len(result.Response.ErrorsByTransferID) > 0 {
		return "error"
	}
	return "clean"
}

func pendingRanges(dbPath string, checkpoint int64, limit int) ([]sampleRange, error) {
	query := fmt.Sprintf(`
WITH pending_ranges AS (
  SELECT s.range_start, s.range_end
  FROM samples s
  WHERE s.range_end <= %d
    AND NOT EXISTS (
      SELECT 1 FROM validation_runs vr
      WHERE vr.range_start = s.range_start
        AND vr.range_end = s.range_end
        AND vr.status IN ('clean', 'mismatch', 'mismatch_error')
    )
  GROUP BY s.range_start, s.range_end
  ORDER BY s.range_start
  LIMIT %d
)
SELECT s.pk_id, s.transfer_id, s.range_start, s.range_end, s.selection_reason
FROM samples s
JOIN pending_ranges p USING (range_start, range_end)
ORDER BY s.range_start, s.sample_rank;
`, checkpoint, limit)
	output, err := runSQLite(dbPath, query, "-json")
	if err != nil {
		return nil, fmt.Errorf("loading pending samples: %w", err)
	}
	var samples []sample
	if len(bytes.TrimSpace(output)) > 0 {
		if err := json.Unmarshal(output, &samples); err != nil {
			return nil, fmt.Errorf("decoding pending samples: %w", err)
		}
	}
	return groupSamples(samples), nil
}

func latestPendingRange(dbPath string, checkpoint int64) (sampleRange, bool, error) {
	query := fmt.Sprintf(`
WITH pending_range AS (
  SELECT s.range_start, s.range_end
  FROM samples s
  WHERE s.range_end <= %d
    AND NOT EXISTS (
      SELECT 1 FROM validation_runs vr
      WHERE vr.range_start = s.range_start
        AND vr.range_end = s.range_end
        AND vr.status IN ('clean', 'mismatch', 'mismatch_error')
    )
    AND NOT EXISTS (
      SELECT 1 FROM monitor_skipped_ranges msr
      WHERE msr.range_start = s.range_start AND msr.range_end = s.range_end
    )
  GROUP BY s.range_start, s.range_end
  ORDER BY EXISTS (
    SELECT 1 FROM validation_runs attempted
    WHERE attempted.range_start = s.range_start AND attempted.range_end = s.range_end
  ) DESC, s.range_start DESC
  LIMIT 1
)
SELECT s.pk_id, s.transfer_id, s.range_start, s.range_end, s.selection_reason
FROM samples s
JOIN pending_range p USING (range_start, range_end)
ORDER BY s.sample_rank;
`, checkpoint)
	output, err := runSQLite(dbPath, query, "-json")
	if err != nil {
		return sampleRange{}, false, fmt.Errorf("loading latest pending samples: %w", err)
	}
	var samples []sample
	if len(bytes.TrimSpace(output)) > 0 {
		if err := json.Unmarshal(output, &samples); err != nil {
			return sampleRange{}, false, fmt.Errorf("decoding latest pending samples: %w", err)
		}
	}
	ranges := groupSamples(samples)
	if len(ranges) == 0 {
		return sampleRange{}, false, nil
	}
	return ranges[0], true, nil
}

func samplesForRange(dbPath string, rangeStart, rangeEnd int64) (sampleRange, error) {
	query := fmt.Sprintf(`
SELECT pk_id, transfer_id, range_start, range_end, selection_reason
FROM samples
WHERE range_start = %d AND range_end = %d
ORDER BY sample_rank;
`, rangeStart, rangeEnd)
	output, err := runSQLite(dbPath, query, "-json")
	if err != nil {
		return sampleRange{}, fmt.Errorf("loading range samples: %w", err)
	}
	var samples []sample
	if err := json.Unmarshal(output, &samples); err != nil {
		return sampleRange{}, fmt.Errorf("decoding range samples: %w", err)
	}
	return sampleRange{Start: rangeStart, End: rangeEnd, Samples: samples}, nil
}

func groupSamples(samples []sample) []sampleRange {
	var ranges []sampleRange
	for _, current := range samples {
		if len(ranges) == 0 || ranges[len(ranges)-1].Start != current.RangeStart {
			ranges = append(ranges, sampleRange{Start: current.RangeStart, End: current.RangeEnd})
		}
		ranges[len(ranges)-1].Samples = append(ranges[len(ranges)-1].Samples, current)
	}
	return ranges
}
