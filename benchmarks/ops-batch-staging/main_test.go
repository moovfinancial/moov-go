package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSampleQueryUsesConfiguredBackfillRanges(t *testing.T) {
	query := sampleQuery(3001, 89_821_930, 10_000)
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
	if err := recordValidation(dbPath, current, options, "test-agent", result, "clean", true); err != nil {
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
