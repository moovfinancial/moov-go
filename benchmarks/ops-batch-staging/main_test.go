package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSampleQueryUsesConfiguredBackfillRanges(t *testing.T) {
	query := sampleQuery(3001, 3001, 89_821_930, 10_000)
	for _, expected := range []string{
		"DIV(pk_id - 3001, 10000)",
		"WHERE pk_id BETWEEN 3001 AND 89821930",
		"3001 + range_index * 10000 AS range_start",
		"MOD(range_index + 1, 10) = 0, 500, 100",
		"projection-backfill-sample-v1",
		"PARTITION BY range_index, COALESCE(account_id, '')",
		"WHEN account_rank = 1 THEN 'account-stratum'",
	} {
		if !strings.Contains(query, expected) {
			t.Fatalf("sample query missing %q", expected)
		}
	}
}

func TestSampleQueryKeepsCampaignRangeAlignmentAcrossCacheChunks(t *testing.T) {
	query := sampleQuery(3001, 253001, 263000, 10_000)
	for _, expected := range []string{
		"DIV(pk_id - 3001, 10000)",
		"WHERE pk_id BETWEEN 253001 AND 263000",
		"3001 + range_index * 10000 AS range_start",
	} {
		if !strings.Contains(query, expected) {
			t.Fatalf("sample cache query missing %q", expected)
		}
	}
}

func TestPlanSampleCacheExtensionUsesOnlyCompletedRanges(t *testing.T) {
	extension, ok, err := planSampleCacheExtension(3001, 10_000, 253_000, 270_000, 1_000_000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || extension.Start != 253_001 || extension.End != 263_000 {
		t.Fatalf("extension = %+v, %t", extension, ok)
	}

	_, ok, err = planSampleCacheExtension(3001, 10_000, 263_000, 270_000, 1_000_000)
	if err != nil || ok {
		t.Fatalf("unexpected extension for complete cache: ok=%t err=%v", ok, err)
	}
}

func TestPlanSampleCacheExtensionBoundsLargeProducerLead(t *testing.T) {
	extension, ok, err := planSampleCacheExtension(3001, 10_000, 3000, 2_500_000, 1_000_000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || extension.Start != 3001 || extension.End != 1_003_000 {
		t.Fatalf("extension = %+v, %t", extension, ok)
	}
}

func TestPlanSampleCacheExtensionRejectsUnalignedCacheCheckpoint(t *testing.T) {
	_, _, err := planSampleCacheExtension(3001, 10_000, 12_500, 30_000, 1_000_000)
	if err == nil || !strings.Contains(err.Error(), "not aligned") {
		t.Fatalf("unaligned checkpoint error = %v", err)
	}
}

func TestExtendSampleCacheSkipsBigQueryWhenCheckpointIsCovered(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := ensureCampaign(prepareOptions{DBPath: dbPath, StartPK: 3001, RangeSize: 10_000}); err != nil {
		t.Fatal(err)
	}
	if _, err := runSQLite(dbPath, "UPDATE sample_cache_state SET cached_through_pk = 263000 WHERE id = 1;"); err != nil {
		t.Fatal(err)
	}

	options := sampleCacheOptions{
		DBPath:         dbPath,
		SafeCheckpoint: 270_000,
		ChunkPKs:       1_000_000,
	}
	if _, err := extendSampleCache(options); err != nil {
		t.Fatalf("covered checkpoint unexpectedly queried BigQuery: %v", err)
	}
	assertScalar(t, dbPath, "SELECT cached_through_pk FROM sample_cache_state WHERE id = 1;", "263000")
}

func TestExtendSampleCacheDoesNotExtendFixedCampaign(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := ensureCampaign(prepareOptions{DBPath: dbPath, StartPK: 3001, MaxPK: 13_000, RangeSize: 10_000}); err != nil {
		t.Fatal(err)
	}

	extension, err := extendSampleCache(sampleCacheOptions{
		DBPath:         dbPath,
		SafeCheckpoint: 30_000,
		ChunkPKs:       1_000_000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if extension != (sampleCacheExtension{}) {
		t.Fatalf("fixed campaign extension = %+v", extension)
	}
}

func TestImportSampleCacheFailureDoesNotAdvanceCheckpoint(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := ensureCampaign(prepareOptions{DBPath: dbPath, StartPK: 3001, RangeSize: 10_000}); err != nil {
		t.Fatal(err)
	}

	err := importSamples(dbPath, filepath.Join(t.TempDir(), "missing.csv"), 13_000)
	if err == nil {
		t.Fatal("expected missing sample export to fail")
	}
	assertScalar(t, dbPath, "SELECT cached_through_pk FROM sample_cache_state WHERE id = 1;", "3000")
}

func TestImportSampleCacheConstraintFailureDoesNotAdvanceCheckpoint(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := ensureCampaign(prepareOptions{DBPath: dbPath, StartPK: 3001, RangeSize: 10_000}); err != nil {
		t.Fatal(err)
	}
	csvPath := filepath.Join(t.TempDir(), "samples.csv")
	csv := "pk_id,transfer_id,range_start,range_end,sample_rank,selection_reason,transfer_type,status,account_id\n" +
		"3001,duplicate,3001,13000,1,range-first,wallet-to-wallet,completed,account-one\n" +
		"3002,duplicate,3001,13000,2,uniform,wallet-to-bank,completed,account-two\n"
	if err := os.WriteFile(csvPath, []byte(csv), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := importSamples(dbPath, csvPath, 13_000); err == nil {
		t.Fatal("expected duplicate transfer ID to fail")
	}
	assertScalar(t, dbPath, "SELECT COUNT(*) FROM samples;", "0")
	assertScalar(t, dbPath, "SELECT cached_through_pk FROM sample_cache_state WHERE id = 1;", "3000")
}

func TestMonitorProgressStateIsDurableAndVersionBound(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := recordProducerCheckpoint(dbPath, "dev-one", 270_000); err != nil {
		t.Fatal(err)
	}
	if err := recordProducerCheckpoint(dbPath, "dev-one", 260_000); err != nil {
		t.Fatal(err)
	}
	checkpoint, ok, err := persistedProducerCheckpoint(dbPath, "dev-one")
	if err != nil || !ok || checkpoint != 270_000 {
		t.Fatalf("persisted checkpoint = %d, %t, %v", checkpoint, ok, err)
	}
	if _, ok, err := persistedProducerCheckpoint(dbPath, "dev-two"); err != nil || ok {
		t.Fatalf("unexpected checkpoint for another version: %t, %v", ok, err)
	}

	selected := sampleRange{Start: 253_001, End: 263_000}
	if err := recordConsumedRange(dbPath, "dev-bff", selected); err != nil {
		t.Fatal(err)
	}
	ready, err := consumedRangeRecorded(dbPath, "dev-bff", selected)
	if err != nil || !ready {
		t.Fatalf("persisted consumed range = %t, %v", ready, err)
	}
	ready, err = consumedRangeRecorded(dbPath, "dev-other", selected)
	if err != nil || ready {
		t.Fatalf("consumed range leaked across versions = %t, %v", ready, err)
	}
}

func TestPendingRangesResumeAfterCompletedRun(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)

	ranges, err := pendingRanges(dbPath, 12_999, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 0 {
		t.Fatalf("got %d ranges before the first checkpoint", len(ranges))
	}

	ranges, err = pendingRanges(dbPath, 23_000, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 1 || ranges[0].Start != 3001 || len(ranges[0].Samples) != 2 {
		t.Fatalf("unexpected first pending range: %+v", ranges)
	}
	recordCleanRun(t, dbPath, ranges[0])

	ranges, err = pendingRanges(dbPath, 23_000, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 1 || ranges[0].Start != 13_001 {
		t.Fatalf("unexpected resumed ranges: %+v", ranges)
	}
	assertScalar(t, dbPath, "SELECT COUNT(*) FROM validation_results WHERE comparison = 'match';", "2")
}

func TestPendingRangesDoNotRetryMixedMismatchAndError(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)
	if _, err := runSQLite(dbPath, `
INSERT INTO validation_runs (
  range_start, range_end, started_at, api_version, user_agent, attempted_count, status
) VALUES (3001, 13000, '2026-08-04T00:00:00Z', 'v2026.07.00', 'test', 2, 'mismatch_error');
`); err != nil {
		t.Fatal(err)
	}

	ranges, err := pendingRanges(dbPath, 13_000, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(ranges) != 0 {
		t.Fatalf("known mismatch remained pending: %+v", ranges)
	}
}

func TestLatestPendingRangeUsesNewestSafeRange(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)

	current, ok, err := latestPendingRange(dbPath, 23_000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || current.Start != 13_001 {
		t.Fatalf("latest pending range = %+v, %t", current, ok)
	}
	if _, err := runSQLite(dbPath, `
INSERT INTO monitor_skipped_ranges VALUES (13001, 23000, '2026-08-05T00:00:00Z', 23000, 'test');
`); err != nil {
		t.Fatal(err)
	}
	current, ok, err = latestPendingRange(dbPath, 23_000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || current.Start != 3001 {
		t.Fatalf("latest pending range after skip = %+v, %t", current, ok)
	}
}

func TestLatestPendingRangeRetriesErrorOnlyResponse(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)
	if _, err := runSQLite(dbPath, `
INSERT INTO validation_runs (
  range_start, range_end, started_at, api_version, user_agent, attempted_count, status
) VALUES (13001, 23000, '2026-08-05T00:00:00Z', 'v2026.07.00', 'test', 1, 'error');
`); err != nil {
		t.Fatal(err)
	}

	current, ok, err := latestPendingRange(dbPath, 23_000)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || current.Start != 13_001 {
		t.Fatalf("error-only range was not retryable: %+v, %t", current, ok)
	}
}

func TestMarkSkippedRangesOnlySkipsOlderPendingRanges(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)
	selected, ok, err := latestPendingRange(dbPath, 23_000)
	if err != nil || !ok {
		t.Fatalf("selecting latest range: %+v, %t, %v", selected, ok, err)
	}
	result := requestResult{StatusCode: http.StatusOK, RequestID: "monitor-request", Duration: time.Second}
	options := validateOptions{Version: defaultAPIVersion, ServiceVersion: "dev-test"}
	runID, err := recordValidation(dbPath, selected, options, "test-agent", result, "evidence_pending", true)
	if err != nil {
		t.Fatal(err)
	}
	if err := finalizeMonitorSuccess(dbPath, runID, "trace-id", selected, 23_000); err != nil {
		t.Fatal(err)
	}
	assertScalar(t, dbPath, "SELECT COUNT(*) FROM monitor_skipped_ranges;", "1")
	assertScalar(t, dbPath, "SELECT range_start FROM monitor_skipped_ranges;", "3001")
	assertScalar(t, dbPath, fmt.Sprintf("SELECT status FROM validation_runs WHERE id = %d;", runID), "clean")
}

func TestFinalizeMonitorSuccessPreservesPriorFailedRange(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)
	options := validateOptions{Version: defaultAPIVersion, ServiceVersion: "dev-test"}
	older, err := samplesForRange(dbPath, 3001, 13_000)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := recordValidation(dbPath, older, options, "failed-agent", requestResult{}, "request_error", false); err != nil {
		t.Fatal(err)
	}
	selected, err := samplesForRange(dbPath, 13_001, 23_000)
	if err != nil {
		t.Fatal(err)
	}
	runID, err := recordValidation(dbPath, selected, options, "monitor-agent", requestResult{StatusCode: http.StatusOK}, "evidence_pending", true)
	if err != nil {
		t.Fatal(err)
	}
	if err := finalizeMonitorSuccess(dbPath, runID, "trace-id", selected, 23_000); err != nil {
		t.Fatal(err)
	}
	assertScalar(t, dbPath, "SELECT COUNT(*) FROM monitor_skipped_ranges;", "0")

	retry, ok, err := latestPendingRange(dbPath, 23_000)
	if err != nil || !ok || retry.Start != 3001 {
		t.Fatalf("prior failed range was not preserved: %+v, %t, %v", retry, ok, err)
	}
}

func TestEvidenceTransitionsTargetValidationRunID(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	seedSamples(t, dbPath)
	options := validateOptions{Version: defaultAPIVersion, ServiceVersion: "dev-test"}
	first, err := samplesForRange(dbPath, 3001, 13_000)
	if err != nil {
		t.Fatal(err)
	}
	second, err := samplesForRange(dbPath, 13_001, 23_000)
	if err != nil {
		t.Fatal(err)
	}
	firstID, err := recordValidation(dbPath, first, options, "first", requestResult{}, "evidence_pending", true)
	if err != nil {
		t.Fatal(err)
	}
	secondID, err := recordValidation(dbPath, second, options, "second", requestResult{}, "evidence_pending", true)
	if err != nil {
		t.Fatal(err)
	}
	if err := recordValidationTrace(dbPath, firstID, "first-trace"); err != nil {
		t.Fatal(err)
	}
	if err := markValidationEvidenceError(dbPath, secondID); err != nil {
		t.Fatal(err)
	}
	assertScalar(t, dbPath, fmt.Sprintf("SELECT trace_id FROM validation_runs WHERE id = %d;", firstID), "first-trace")
	assertScalar(t, dbPath, fmt.Sprintf("SELECT status FROM validation_runs WHERE id = %d;", secondID), "evidence_error")
}

func TestValidationStatusPreservesMixedMismatchAndError(t *testing.T) {
	result := requestResult{Response: validationResponse{
		TransferIDsWithDiff: []string{"mismatch"},
		ErrorsByTransferID:  map[string]string{"error": "internal error"},
	}}
	if got := validationStatus(result, nil); got != "mismatch_error" {
		t.Fatalf("validation status = %q", got)
	}
}

func TestStagingCredentialsRequireExactHost(t *testing.T) {
	t.Setenv("MOOV_PUBLIC_KEY", "username")
	t.Setenv("MOOV_SECRET_KEY", "credential")
	t.Setenv("MOOV_HOST", stagingAPIHost)
	if _, err := stagingCredentials(); err != nil {
		t.Fatalf("expected staging host to pass: %v", err)
	}

	t.Setenv("MOOV_HOST", "staging@attacker.example")
	if _, err := stagingCredentials(); err == nil {
		t.Fatal("expected attacker host to be rejected")
	}
}

func TestCampaignDatabaseMustExistAndContainSamples(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := requireCampaignDB(dbPath); err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("missing database error = %v", err)
	}
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	options := prepareOptions{DBPath: dbPath, StartPK: 3001, MaxPK: 13_000, RangeSize: 10_000}
	if err := ensureCampaign(options); err != nil {
		t.Fatal(err)
	}
	if err := requireCampaignDB(dbPath); err == nil || !strings.Contains(err.Error(), "has no samples") {
		t.Fatalf("empty database error = %v", err)
	}
	seedSamples(t, dbPath)
	if err := requireCampaignDB(dbPath); err != nil {
		t.Fatalf("prepared database rejected: %v", err)
	}
}

func TestDynamicCampaignCanStartWithoutSamples(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	if err := ensureCampaign(prepareOptions{DBPath: dbPath, StartPK: 3001, RangeSize: 10_000}); err != nil {
		t.Fatal(err)
	}
	if err := requireCampaignDB(dbPath); err != nil {
		t.Fatalf("dynamic campaign rejected: %v", err)
	}
	assertScalar(t, dbPath, "SELECT cached_through_pk FROM sample_cache_state WHERE id = 1;", "3000")
}

func TestCampaignOptionsAreImmutable(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "campaign.sqlite")
	if err := initializeDB(dbPath); err != nil {
		t.Fatal(err)
	}
	options := prepareOptions{DBPath: dbPath, StartPK: 3001, MaxPK: 13_000, RangeSize: 10_000}
	if err := ensureCampaign(options); err != nil {
		t.Fatal(err)
	}
	options.StartPK = 4001
	if err := ensureCampaign(options); err == nil {
		t.Fatal("expected incompatible campaign options to be rejected")
	}
}

func TestRequestValidation(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, credential, ok := r.BasicAuth()
		if !ok || username != "username" || credential != "credential" {
			t.Error("missing staging credentials")
		}
		if r.Header.Get("X-Moov-Version") != defaultAPIVersion {
			t.Errorf("version = %q", r.Header.Get("X-Moov-Version"))
		}
		if r.Header.Get("User-Agent") != "validator-test" {
			t.Errorf("user agent = %q", r.Header.Get("User-Agent"))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var request struct {
			TransferIDs []string `json:"transferIDs"`
		}
		if err := json.Unmarshal(body, &request); err != nil {
			t.Fatal(err)
		}
		if len(request.TransferIDs) != 2 {
			t.Errorf("transfer IDs = %v", request.TransferIDs)
		}
		w.Header().Set("X-Request-ID", "request-id")
		_, _ = io.WriteString(w, `{"transferIDsWithDiff":["two"],"errorsByTransferID":{}}`)
	}))
	defer server.Close()

	auth := credentials{Host: strings.TrimPrefix(server.URL, "https://"), Username: "username", Credential: "credential"}
	samples := []sample{{TransferID: "one"}, {TransferID: "two"}}
	result, err := requestValidation(context.Background(), server.Client(), auth, defaultAPIVersion, "validator-test", samples)
	if err != nil {
		t.Fatal(err)
	}
	if result.RequestID != "request-id" || validationStatus(result, nil) != "mismatch" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestHoneycombMCPRunQueryReadsResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Error("missing MCP authorization")
		}
		var request struct {
			Method string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		switch request.Method {
		case "tools/call":
			w.Header().Set("Mcp-Session-Id", "session-id")
			writeSSE(t, w, map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{
				"content": []map[string]any{{"type": "resource_link", "uri": "honeycomb://query-run/test-run/json"}},
			}})
		case "resources/read":
			if r.Header.Get("Mcp-Session-Id") != "session-id" {
				t.Error("missing MCP session ID")
			}
			writeSSE(t, w, map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{
				"contents": []map[string]any{{"text": `{"results":[{"data":{"checkpoint":123000}}]}`}},
			}})
		default:
			t.Fatalf("unexpected method %q", request.Method)
		}
	}))
	defer server.Close()

	client := newHoneycombMCPClient(server.URL, "token")
	client.HTTPClient = server.Client()
	result, err := client.runQuery(context.Background(), map[string]any{"test": true})
	if err != nil {
		t.Fatal(err)
	}
	checkpoint, err := number(result.Rows[0], "checkpoint")
	if err != nil {
		t.Fatal(err)
	}
	if result.RunPK != "test-run" || checkpoint != 123_000 {
		t.Fatalf("unexpected query result: %+v", result)
	}
}

func TestReadMCPResponseDataSupportsJSONAndMultilineSSE(t *testing.T) {
	jsonData, err := readMCPResponseData("application/json; charset=utf-8", strings.NewReader(`{"jsonrpc":"2.0"}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(jsonData) != `{"jsonrpc":"2.0"}` {
		t.Fatalf("JSON response = %q", jsonData)
	}

	sseData, err := readMCPResponseData("text/event-stream", strings.NewReader("event: message\ndata: {\"jsonrpc\":\ndata:\"2.0\"}\n\n"))
	if err != nil {
		t.Fatal(err)
	}
	if string(sseData) != "{\"jsonrpc\":\n\"2.0\"}" {
		t.Fatalf("SSE response = %q", sseData)
	}
}

func TestComparisonEvidenceRequiresEveryAttemptToMatch(t *testing.T) {
	evidence := comparisonEvidence{Count: 100, Matches: 99, Mismatches: 1, TraceID: "trace", QueryRunPK: "query"}
	if err := evidence.requireClean(100, "dev-test"); err == nil {
		t.Fatal("expected mismatch evidence to fail")
	}
	evidence = comparisonEvidence{Count: 100, Matches: 100, TraceID: "trace", QueryRunPK: "query"}
	if err := evidence.requireClean(100, "dev-test"); err != nil {
		t.Fatalf("clean evidence failed: %v", err)
	}
}

func TestSQLQuote(t *testing.T) {
	if got := sqlQuote("projection's value"); got != "'projection''s value'" {
		t.Fatalf("sqlQuote = %q", got)
	}
}

func seedSamples(t *testing.T, dbPath string) {
	t.Helper()
	_, err := runSQLite(dbPath, `
INSERT INTO samples VALUES
  (3001, 'one', 3001, 13000, 1, 'range-first', 'wallet-to-wallet', 'completed', 'account-one'),
  (3002, 'two', 3001, 13000, 2, 'uniform', 'wallet-to-bank', 'completed', 'account-two'),
  (13001, 'three', 13001, 23000, 1, 'range-first', 'card-to-wallet', 'completed', 'account-one');
`)
	if err != nil {
		t.Fatal(err)
	}
}

func recordCleanRun(t *testing.T, dbPath string, current sampleRange) {
	t.Helper()
	result := requestResult{StatusCode: http.StatusOK, RequestID: "request-id", Duration: time.Second}
	options := validateOptions{Version: defaultAPIVersion, ServiceVersion: "dev-test"}
	if _, err := recordValidation(dbPath, current, options, "test-agent", result, "clean", true); err != nil {
		t.Fatal(err)
	}
}

func assertScalar(t *testing.T, dbPath, query, expected string) {
	t.Helper()
	actual, err := scalar(dbPath, query)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("scalar = %q, want %q", actual, expected)
	}
}

func writeSSE(t *testing.T, writer io.Writer, value any) {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = fmt.Fprintf(writer, "event: message\ndata: %s\n\n", encoded)
}
