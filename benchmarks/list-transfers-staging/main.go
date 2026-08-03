package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	moovFinancialAccountID = "db04bf9d-91f6-4206-ba38-6844636532ad"
	highVolumeAccountID    = "6f36ae5e-d715-4b5f-9089-11dd8f21c120"
	defaultVersion         = "v2026.07.00"
	noMatchID              = "11111111-1111-4111-8111-111111111111"
)

type fixture struct {
	AccountID            string `json:"accountID"`
	TransferCount        int64  `json:"transferCount"`
	FirstTransferOn      string `json:"firstTransferOn"`
	LatestTransferOn     string `json:"latestTransferOn"`
	SourceAccountID      string `json:"sourceAccountID"`
	DestinationAccountID string `json:"destinationAccountID"`
	GroupID              string `json:"groupID,omitempty"`
	ScheduleID           string `json:"scheduleID,omitempty"`
	PaymentLinkCode      string `json:"paymentLinkCode,omitempty"`
	ForeignID            string `json:"foreignID,omitempty"`
	AuthorizationID      string `json:"authorizationID,omitempty"`
	CaptureID            string `json:"captureID,omitempty"`
	HasRefunds           bool   `json:"hasRefunds"`
	HasDisputes          bool   `json:"hasDisputes"`
}

var fixtures = map[string]fixture{
	moovFinancialAccountID: {
		AccountID:            moovFinancialAccountID,
		TransferCount:        261129,
		FirstTransferOn:      "2021-08-13T21:40:54Z",
		LatestTransferOn:     "2026-08-03T21:51:58Z",
		SourceAccountID:      moovFinancialAccountID,
		DestinationAccountID: "bcbfa114-2475-487e-bdd6-5f3d7efeff20",
		GroupID:              "bbd5020b-bf29-48f6-b52f-85fc8d78418c",
		ScheduleID:           "8574051f-06c8-4045-9cd3-b24cd5984e69",
		PaymentLinkCode:      "1rRyWuW8hS",
		ForeignID:            "123",
		HasRefunds:           true,
	},
	highVolumeAccountID: {
		AccountID:            highVolumeAccountID,
		TransferCount:        86892676,
		FirstTransferOn:      "2022-03-04T22:02:56Z",
		LatestTransferOn:     "2026-08-03T21:23:15Z",
		SourceAccountID:      "0e69d339-c00f-4ce2-b3c2-790d5ac379dc",
		DestinationAccountID: "00d98253-6b4a-44de-aafc-47f88b30d34d",
	},
}

type benchmarkCase struct {
	Name          string     `json:"name"`
	Category      string     `json:"category"`
	PathAccountID string     `json:"pathAccountID"`
	Version       string     `json:"version"`
	Params        url.Values `json:"params"`
	FixtureMatch  bool       `json:"fixtureMatch"`
}

type measurement struct {
	CaseName      string    `json:"caseName"`
	Iteration     int       `json:"iteration"`
	StartedAt     time.Time `json:"startedAt"`
	DurationMS    float64   `json:"durationMs"`
	StatusCode    int       `json:"statusCode"`
	RequestID     string    `json:"requestID,omitempty"`
	ResponseCount *int      `json:"responseCount,omitempty"`
	ResponseBytes int       `json:"responseBytes"`
	Error         string    `json:"error,omitempty"`
}

type output struct {
	Label        string          `json:"label"`
	Account      fixture         `json:"account"`
	Host         string          `json:"host"`
	StartedAt    time.Time       `json:"startedAt"`
	FinishedAt   time.Time       `json:"finishedAt"`
	Warmups      int             `json:"warmups"`
	Iterations   int             `json:"iterations"`
	Cases        []benchmarkCase `json:"cases"`
	Measurements []measurement   `json:"measurements"`
}

func main() {
	accountID := flag.String("account", moovFinancialAccountID, "partner account ID whose LIST queries should be benchmarked")
	label := flag.String("label", "baseline", "run label, such as baseline or projection")
	casePattern := flag.String("case", ".*", "regular expression selecting case names")
	warmups := flag.Int("warmups", 1, "unmeasured requests per case")
	iterations := flag.Int("iterations", 3, "measured requests per case")
	outputPath := flag.String("output", "", "result JSON path; defaults under benchmarks/list-transfers-staging/results")
	listCases := flag.Bool("list", false, "print the selected cases without making requests")
	flag.Parse()

	if err := run(*accountID, *label, *casePattern, *warmups, *iterations, *outputPath, *listCases); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(accountID, label, casePattern string, warmups, iterations int, outputPath string, listCases bool) error {
	accountFixture, ok := fixtures[accountID]
	if !ok {
		return fmt.Errorf("no stable fixture for account %s", accountID)
	}
	if warmups < 0 || iterations < 1 {
		return errors.New("warmups must be non-negative and iterations must be positive")
	}
	pattern, err := regexp.Compile(casePattern)
	if err != nil {
		return fmt.Errorf("compile case pattern: %w", err)
	}

	cases := selectCases(benchmarkCases(accountFixture), pattern)
	if len(cases) == 0 {
		return errors.New("case pattern selected no benchmarks")
	}
	if listCases {
		for _, benchmark := range cases {
			fmt.Printf("%-36s %-12s %s\n", benchmark.Name, benchmark.Category, benchmark.Params.Encode())
		}
		return nil
	}

	host := os.Getenv("MOOV_HOST")
	username := os.Getenv("MOOV_PUBLIC_KEY")
	credential := os.Getenv("MOOV_SECRET_KEY")
	if host == "" || username == "" || credential == "" {
		return errors.New("MOOV_HOST, MOOV_PUBLIC_KEY, and MOOV_SECRET_KEY must be set; run moov_env staging first")
	}
	if !strings.Contains(strings.ToLower(host), "staging") {
		return fmt.Errorf("refusing to benchmark non-staging host %q", host)
	}

	result := output{
		Label:      label,
		Account:    accountFixture,
		Host:       host,
		StartedAt:  time.Now().UTC(),
		Warmups:    warmups,
		Iterations: iterations,
		Cases:      cases,
	}
	client := &http.Client{Timeout: 30 * time.Second}
	ctx := context.Background()
	var runErr error

	for _, benchmark := range cases {
		fmt.Printf("%-36s", benchmark.Name)
		for i := 0; i < warmups; i++ {
			if _, err := request(ctx, client, host, username, credential, label, benchmark, 0); err != nil {
				fmt.Printf(" warmup-error:%v", err)
			}
		}
		for i := 1; i <= iterations; i++ {
			m, requestErr := request(ctx, client, host, username, credential, label, benchmark, i)
			result.Measurements = append(result.Measurements, m)
			fmt.Printf(" %d:%.0fms/%d", i, m.DurationMS, m.StatusCode)
			if requestErr != nil {
				runErr = errors.Join(runErr, requestErr)
			}
		}
		fmt.Println()
	}

	return writeResult(result, outputPath, runErr)
}

func request(ctx context.Context, client *http.Client, host, username, credential, label string, benchmark benchmarkCase, iteration int) (measurement, error) {
	path := "/transfers"
	if benchmark.PathAccountID != "" {
		path = "/accounts/" + url.PathEscape(benchmark.PathAccountID) + "/transfers"
	}
	u := url.URL{Scheme: "https", Host: host, Path: path, RawQuery: benchmark.Params.Encode()}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return measurement{}, err
	}
	req.SetBasicAuth(username, credential)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "moov-go-list-transfers-staging-benchmark/1 "+safeName(label))
	req.Header.Set("X-Moov-Version", benchmark.Version)

	started := time.Now().UTC()
	resp, err := client.Do(req)
	m := measurement{CaseName: benchmark.Name, Iteration: iteration, StartedAt: started, DurationMS: float64(time.Since(started).Microseconds()) / 1000}
	if err != nil {
		m.Error = err.Error()
		return m, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	m.StatusCode = resp.StatusCode
	m.RequestID = resp.Header.Get("X-Request-ID")
	m.ResponseBytes = len(body)
	if err != nil {
		m.Error = err.Error()
		return m, err
	}
	var transfers []json.RawMessage
	if json.Unmarshal(body, &transfers) == nil {
		count := len(transfers)
		m.ResponseCount = &count
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		m.Error = strings.TrimSpace(string(body))
		if len(m.Error) > 500 {
			m.Error = m.Error[:500]
		}
		return m, fmt.Errorf("%s returned HTTP %d (request %s): %s", benchmark.Name, resp.StatusCode, m.RequestID, m.Error)
	}
	return m, nil
}

func writeResult(result output, outputPath string, runErr error) error {
	result.FinishedAt = time.Now().UTC()
	if outputPath == "" {
		accountPrefix := strings.Split(result.Account.AccountID, "-")[0]
		name := fmt.Sprintf("%s-%s-%s.json", safeName(result.Label), accountPrefix, result.StartedAt.Format("20060102T150405Z"))
		outputPath = filepath.Join("benchmarks", "list-transfers-staging", "results", name)
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return errors.Join(runErr, err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		err = os.WriteFile(outputPath, append(data, '\n'), 0o644)
	}
	if err != nil {
		return errors.Join(runErr, err)
	}
	fmt.Printf("wrote %s\n", outputPath)
	return runErr
}

func benchmarkCases(f fixture) []benchmarkCase {
	sparse := func(value string) (string, bool) {
		if value == "" {
			return noMatchID, false
		}
		return value, true
	}
	groupID, groupMatch := sparse(f.GroupID)
	scheduleID, scheduleMatch := sparse(f.ScheduleID)
	authorizationID, authorizationMatch := sparse(f.AuthorizationID)
	captureID, captureMatch := sparse(f.CaptureID)
	paymentLinkCode, paymentLinkMatch := f.PaymentLinkCode, true
	if paymentLinkCode == "" {
		paymentLinkCode, paymentLinkMatch = "amp-list-benchmark-no-match", false
	}
	foreignID, foreignMatch := f.ForeignID, true
	if foreignID == "" {
		foreignID, foreignMatch = "amp-list-benchmark-no-match", false
	}

	accountPath := f.AccountID
	cases := []benchmarkCase{
		newCase("base-default", "common", accountPath, defaultVersion, true),
		newCase("base-global-endpoint", "common", "", defaultVersion, true),
		newCase("page-tiny-5", "pagination", accountPath, defaultVersion, true, "count", "5"),
		newCase("page-small-10", "pagination", accountPath, defaultVersion, true, "count", "10"),
		newCase("page-default-20", "pagination", accountPath, defaultVersion, true, "count", "20"),
		newCase("page-large-1000", "pagination", accountPath, defaultVersion, true, "count", "1000"),
		newCase("page-offset-140-count-20", "pagination", accountPath, defaultVersion, true, "skip", "140", "count", "20"),
		newCase("page-offset-200-count-20", "pagination", accountPath, defaultVersion, true, "skip", "200", "count", "20"),
		newCase("page-shallow-skip-200", "pagination", accountPath, defaultVersion, true, "skip", "200", "count", "200"),
		newCase("page-deep-skip-50000", "pagination", accountPath, defaultVersion, true, "skip", "50000", "count", "200"),
		newCase("status-completed", "indexed", accountPath, defaultVersion, true, "status", "completed"),
		newCase("status-failed-page-2", "indexed-combined", accountPath, defaultVersion, true, "status", "failed", "skip", "20", "count", "20"),
		newCase("date-recent-week", "indexed", accountPath, defaultVersion, true, "startDateTime", "2026-07-27T00:00:00Z", "endDateTime", "2026-08-04T00:00:00Z"),
		newCase("date-old-month", "indexed", accountPath, defaultVersion, true, "startDateTime", "2023-01-01T00:00:00Z", "endDateTime", "2023-02-01T00:00:00Z"),
		newCase("date-broad-history", "indexed", accountPath, defaultVersion, true, "startDateTime", "2022-01-01T00:00:00Z", "endDateTime", "2026-08-04T00:00:00Z"),
		newCase("status-completed-recent", "indexed-combined", accountPath, defaultVersion, true, "status", "completed", "startDateTime", "2026-07-27T00:00:00Z", "endDateTime", "2026-08-04T00:00:00Z"),
		newCase("account-source", "account-filter", accountPath, defaultVersion, true, "accountIDs", f.SourceAccountID),
		newCase("account-destination", "account-filter", accountPath, defaultVersion, true, "accountIDs", f.DestinationAccountID),
		newCase("account-source-destination-union", "account-filter", accountPath, defaultVersion, true, "accountIDs", f.SourceAccountID+","+f.DestinationAccountID),
		newCase("account-no-match", "account-filter", accountPath, defaultVersion, false, "accountIDs", noMatchID),
		newCase("guard-party-path-account", "authorization-guard", f.DestinationAccountID, defaultVersion, true, "count", "200"),
		newCase("sparse-group", "sparse", accountPath, defaultVersion, groupMatch, "groupID", groupID),
		newCase("sparse-schedule", "sparse", accountPath, defaultVersion, scheduleMatch, "scheduleID", scheduleID),
		newCase("sparse-payment-link", "sparse", accountPath, defaultVersion, paymentLinkMatch, "paymentLinkCode", paymentLinkCode),
		newCase("sparse-foreign-id", "sparse", accountPath, defaultVersion, foreignMatch, "foreignID", foreignID),
		newCase("sparse-authorization", "sparse-auth-capture", accountPath, defaultVersion, authorizationMatch, "authorizationIDs", authorizationID),
		newCase("sparse-capture", "sparse-auth-capture", accountPath, defaultVersion, captureMatch, "captureIDs", captureID),
		newCase("boolean-refunded", "sparse-boolean", accountPath, defaultVersion, f.HasRefunds, "refunded", "true"),
		newCase("boolean-disputed", "sparse-boolean", accountPath, defaultVersion, f.HasDisputes, "disputed", "true"),
		newCase("boolean-refunded-disputed", "sparse-combined", accountPath, defaultVersion, f.HasRefunds && f.HasDisputes, "refunded", "true", "disputed", "true"),
		newCase("status-group", "sparse-combined", accountPath, defaultVersion, groupMatch, "status", "completed", "groupID", groupID),
		newCase("status-account-date", "indexed-combined", accountPath, defaultVersion, true, "status", "completed", "accountIDs", f.DestinationAccountID, "startDateTime", "2026-07-27T00:00:00Z", "endDateTime", "2026-08-04T00:00:00Z"),
		newCase("empty-foreign-id", "empty-result", accountPath, defaultVersion, false, "foreignID", "amp-list-benchmark-no-match"),
		newCase("page-past-end", "empty-result", accountPath, defaultVersion, false, "skip", fmt.Sprint(f.TransferCount+1000), "count", "200"),
	}
	return cases
}

func newCase(name, category, pathAccountID, version string, fixtureMatch bool, pairs ...string) benchmarkCase {
	params := make(url.Values)
	for i := 0; i < len(pairs); i += 2 {
		params.Set(pairs[i], pairs[i+1])
	}
	return benchmarkCase{Name: name, Category: category, PathAccountID: pathAccountID, Version: version, Params: params, FixtureMatch: fixtureMatch}
}

func selectCases(cases []benchmarkCase, pattern *regexp.Regexp) []benchmarkCase {
	selected := make([]benchmarkCase, 0, len(cases))
	for _, benchmark := range cases {
		if pattern.MatchString(benchmark.Name) {
			selected = append(selected, benchmark)
		}
	}
	sort.SliceStable(selected, func(i, j int) bool { return selected[i].Name < selected[j].Name })
	return selected
}

func safeName(value string) string {
	value = strings.ToLower(value)
	return regexp.MustCompile(`[^a-z0-9_-]+`).ReplaceAllString(value, "-")
}
