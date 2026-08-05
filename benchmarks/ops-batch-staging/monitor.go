package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"syscall"
	"time"
)

type monitorOptions struct {
	DBPath          string
	Version         string
	ServiceVersion  string
	ProducerVersion string
	SafeLagRanges   int64
	CacheChunkPKs   int64
	MaximumBQBytes  int64
	DryRun          bool
	MCPToken        string
	MCPURL          string
	TelemetryDelay  time.Duration
}

type rolloutHealth struct {
	ProducerCheckpoint     int64
	ProducerIterations     int64
	ProducerErrors         int64
	ProducerWorkflowErrors int64
	ConsumerUpserts        int64
	ConsumerErrors         int64
	ConsumerWorkflowErrors int64
	ProducerQueryPK        string
	ConsumerQueryPK        string
	ProducerErrorsQueryPK  string
	ConsumerErrorsQueryPK  string
}

type comparisonEvidence struct {
	Count      int64
	Matches    int64
	Mismatches int64
	Errors     int64
	TraceID    string
	QueryRunPK string
}

func monitorOnce(options monitorOptions) error {
	if options.MCPToken == "" {
		return errors.New("HONEYCOMB_MCP_API_KEY is required")
	}
	if err := requireCampaignDB(options.DBPath); err != nil {
		return err
	}
	release, err := acquireMonitorLock(options.DBPath + ".monitor.lock")
	if err != nil {
		return err
	}
	defer release()

	rangeSizeText, err := scalar(options.DBPath, "SELECT range_size FROM campaign WHERE id = 1;")
	if err != nil {
		return fmt.Errorf("reading campaign range size: %w", err)
	}
	rangeSize, err := strconv.ParseInt(rangeSizeText, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing campaign range size %q: %w", rangeSizeText, err)
	}

	honeycomb := newHoneycombMCPClient(options.MCPURL, options.MCPToken)
	healthCtx, cancelHealth := context.WithTimeout(context.Background(), 60*time.Second)
	health, err := readRolloutHealth(healthCtx, honeycomb, options)
	cancelHealth()
	if err != nil {
		return err
	}
	if health.ProducerErrors > 0 {
		return fmt.Errorf("producer reported %d recent errors", health.ProducerErrors)
	}
	if health.ProducerWorkflowErrors > 0 {
		return fmt.Errorf("Transfers backfill traces reported %d recent errors: %s", health.ProducerWorkflowErrors,
			honeycombQueryURL("transfers", health.ProducerErrorsQueryPK))
	}
	if health.ConsumerErrors > 0 {
		return fmt.Errorf("consumer reported %d recent errors", health.ConsumerErrors)
	}
	if health.ConsumerWorkflowErrors > 0 {
		return fmt.Errorf("transfersbff2 TransferUpdated consumer traces reported %d recent errors: %s", health.ConsumerWorkflowErrors,
			honeycombQueryURL("transfersbff2", health.ConsumerErrorsQueryPK))
	}
	if health.ProducerIterations > 0 {
		if health.ProducerCheckpoint == 0 {
			return fmt.Errorf("producer reported recent iterations without a checkpoint: query=%s",
				honeycombQueryURL("transfers", health.ProducerQueryPK))
		}
		if !options.DryRun {
			if err := recordProducerCheckpoint(options.DBPath, options.ProducerVersion, health.ProducerCheckpoint); err != nil {
				return err
			}
		}
	} else {
		checkpoint, ok, err := persistedProducerCheckpoint(options.DBPath, options.ProducerVersion)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("producer has no recent iterations or persisted checkpoint for %s: query=%s",
				options.ProducerVersion, honeycombQueryURL("transfers", health.ProducerQueryPK))
		}
		health.ProducerCheckpoint = checkpoint
	}

	safeCheckpoint := health.ProducerCheckpoint - rangeSize*options.SafeLagRanges
	if safeCheckpoint < 1 {
		return fmt.Errorf("producer checkpoint %d is too early for %d safe lag ranges", health.ProducerCheckpoint, options.SafeLagRanges)
	}
	extension, err := extendSampleCache(sampleCacheOptions{
		DBPath:             options.DBPath,
		SafeCheckpoint:     safeCheckpoint,
		ChunkPKs:           options.CacheChunkPKs,
		MaximumBytesBilled: options.MaximumBQBytes,
		DryRun:             options.DryRun,
	})
	if err != nil {
		return fmt.Errorf("extending sample cache: %w", err)
	}
	if options.DryRun && extension.End > 0 {
		fmt.Printf("monitor dry-run: producer=%d safeCheckpoint=%d consumerUpserts=%d sample-cache-extension=%d-%d\n",
			health.ProducerCheckpoint, safeCheckpoint, health.ConsumerUpserts, extension.Start, extension.End)
		return nil
	}
	if extension.End > 0 {
		fmt.Printf("sample cache extended: %d-%d\n", extension.Start, extension.End)
	}
	selected, ok, err := latestPendingRange(options.DBPath, safeCheckpoint)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Printf("monitor clean: producer=%d safeCheckpoint=%d consumerUpserts=%d no-new-range\n",
			health.ProducerCheckpoint, safeCheckpoint, health.ConsumerUpserts)
		return nil
	}
	readinessCtx, cancelReadiness := context.WithTimeout(context.Background(), 60*time.Second)
	consumerReady, err := projectionSampleConsumed(readinessCtx, honeycomb, options.DBPath, options.ServiceVersion, selected, !options.DryRun)
	cancelReadiness()
	if err != nil {
		return err
	}
	if !consumerReady {
		fmt.Printf("monitor waiting: producer=%d safeCheckpoint=%d selected=%d-%d consumer-sample-not-seen\n",
			health.ProducerCheckpoint, safeCheckpoint, selected.Start, selected.End)
		return nil
	}
	if options.DryRun {
		fmt.Printf("monitor dry-run: producer=%d safeCheckpoint=%d consumerUpserts=%d selected=%d-%d samples=%d consumerReady=true\n",
			health.ProducerCheckpoint, safeCheckpoint, health.ConsumerUpserts,
			selected.Start, selected.End, len(selected.Samples))
		return nil
	}

	auth, err := stagingCredentials()
	if err != nil {
		return err
	}
	validateOptions := validateOptions{
		DBPath:          options.DBPath,
		Version:         options.Version,
		ServiceVersion:  options.ServiceVersion,
		MaxRanges:       1,
		UserAgentSuffix: fmt.Sprintf("monitor-%d", time.Now().UTC().UnixNano()),
		EvidencePending: true,
	}
	validationCtx, cancelValidation := context.WithTimeout(context.Background(), 130*time.Second)
	outcome, validationErr := validateRange(validationCtx, &http.Client{Timeout: 120 * time.Second}, auth, validateOptions, selected)
	cancelValidation()
	if options.TelemetryDelay == 0 {
		options.TelemetryDelay = 5 * time.Second
	}
	var evidence comparisonEvidence
	var evidenceErr error
	if outcome.Result.StatusCode == http.StatusOK {
		evidenceCtx, cancelEvidence := context.WithTimeout(context.Background(), 90*time.Second)
		evidence, evidenceErr = waitForComparisonEvidence(
			evidenceCtx, honeycomb, options.ServiceVersion, outcome.UserAgent, outcome.Attempted, options.TelemetryDelay,
		)
		cancelEvidence()
		if evidenceErr == nil && validationErr != nil {
			evidenceErr = recordValidationTrace(options.DBPath, outcome.RunID, evidence.TraceID)
		}
	}
	if validationErr != nil {
		return errors.Join(validationErr, evidenceErr)
	}
	if evidenceErr == nil {
		evidenceErr = evidence.requireClean(outcome.Attempted, options.ServiceVersion)
	}
	if evidenceErr != nil {
		return errors.Join(evidenceErr, markValidationEvidenceError(options.DBPath, outcome.RunID))
	}
	if err := finalizeMonitorSuccess(options.DBPath, outcome.RunID, evidence.TraceID, selected, safeCheckpoint); err != nil {
		return err
	}

	fmt.Printf("monitor clean: producer=%d safeCheckpoint=%d consumerUpserts=%d range=%d-%d attempted=%d requestID=%s traceID=%s query=%s\n",
		health.ProducerCheckpoint, safeCheckpoint, health.ConsumerUpserts,
		selected.Start, selected.End, outcome.Attempted, outcome.Result.RequestID,
		evidence.TraceID, honeycombQueryURL("transfersbff2", evidence.QueryRunPK))
	return nil
}

func projectionSampleConsumed(ctx context.Context, client *honeycombMCPClient, dbPath, serviceVersion string, selected sampleRange, persist bool) (bool, error) {
	if len(selected.Samples) == 0 {
		return false, errors.New("selected range has no samples")
	}
	ready, err := consumedRangeRecorded(dbPath, serviceVersion, selected)
	if err != nil || ready {
		return ready, err
	}
	transferIDs := make([]string, 0, len(selected.Samples))
	for _, current := range selected.Samples {
		transferIDs = append(transferIDs, current.TransferID)
	}
	query, err := client.runQuery(ctx, map[string]any{
		"environment_slug": "staging",
		"dataset_slug":     "transfersbff2",
		"query_spec": map[string]any{
			"calculations": []map[string]any{{"op": "COUNT_DISTINCT", "column": "transfer_id", "name": "consumed"}},
			"filters": []map[string]any{
				{"column": "name", "op": "=", "value": "projection-upsert-transfer"},
				{"column": "service.version", "op": "=", "value": serviceVersion},
				{"column": "transfer_id", "op": "in", "value": transferIDs},
			},
			"time_range": "30d",
		},
	})
	if err != nil {
		return false, fmt.Errorf("querying selected range consumption: %w", err)
	}
	if len(query.Rows) != 1 {
		return false, fmt.Errorf("unexpected selected range consumption rows: %d", len(query.Rows))
	}
	count, err := number(query.Rows[0], "consumed")
	if err != nil {
		return false, err
	}
	if count != int64(len(transferIDs)) {
		return false, nil
	}
	if persist {
		if err := recordConsumedRange(dbPath, serviceVersion, selected); err != nil {
			return false, err
		}
	}
	return true, nil
}

func recordProducerCheckpoint(dbPath, serviceVersion string, checkpoint int64) error {
	_, err := runSQLite(dbPath, fmt.Sprintf(`
INSERT INTO monitor_producer_state (id, service_version, checkpoint, observed_at)
VALUES (1, %s, %d, %s)
ON CONFLICT(id) DO UPDATE SET
  service_version = excluded.service_version,
  checkpoint = CASE
    WHEN monitor_producer_state.service_version = excluded.service_version
      AND monitor_producer_state.checkpoint > excluded.checkpoint
    THEN monitor_producer_state.checkpoint
    ELSE excluded.checkpoint
  END,
  observed_at = excluded.observed_at;
`, sqlQuote(serviceVersion), checkpoint, sqlQuote(time.Now().UTC().Format(time.RFC3339Nano))))
	if err != nil {
		return fmt.Errorf("recording producer checkpoint: %w", err)
	}
	return nil
}

func persistedProducerCheckpoint(dbPath, serviceVersion string) (int64, bool, error) {
	value, err := scalar(dbPath, fmt.Sprintf(
		"SELECT checkpoint FROM monitor_producer_state WHERE id = 1 AND service_version = %s;",
		sqlQuote(serviceVersion),
	))
	if err != nil {
		return 0, false, fmt.Errorf("reading persisted producer checkpoint: %w", err)
	}
	if value == "" {
		return 0, false, nil
	}
	checkpoint, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("parsing persisted producer checkpoint %q: %w", value, err)
	}
	return checkpoint, true, nil
}

func consumedRangeRecorded(dbPath, serviceVersion string, selected sampleRange) (bool, error) {
	count, err := scalar(dbPath, fmt.Sprintf(`
SELECT COUNT(*) FROM monitor_consumed_ranges
WHERE range_start = %d AND range_end = %d AND service_version = %s;
`, selected.Start, selected.End, sqlQuote(serviceVersion)))
	if err != nil {
		return false, fmt.Errorf("reading consumed range state: %w", err)
	}
	return count == "1", nil
}

func recordConsumedRange(dbPath, serviceVersion string, selected sampleRange) error {
	_, err := runSQLite(dbPath, fmt.Sprintf(`
INSERT INTO monitor_consumed_ranges (range_start, range_end, service_version, observed_at)
VALUES (%d, %d, %s, %s)
ON CONFLICT(range_start, range_end, service_version)
DO UPDATE SET observed_at = excluded.observed_at;
`, selected.Start, selected.End, sqlQuote(serviceVersion), sqlQuote(time.Now().UTC().Format(time.RFC3339Nano))))
	if err != nil {
		return fmt.Errorf("recording consumed range state: %w", err)
	}
	return nil
}

func waitForComparisonEvidence(ctx context.Context, client *honeycombMCPClient, serviceVersion, userAgent string, attempted int, interval time.Duration) (comparisonEvidence, error) {
	var lastEvidence comparisonEvidence
	var lastErr error
	for {
		lastEvidence, lastErr = readComparisonEvidence(ctx, client, serviceVersion, userAgent)
		if lastErr == nil && lastEvidence.Count >= int64(attempted) {
			return lastEvidence, nil
		}

		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			if lastErr != nil {
				return lastEvidence, fmt.Errorf("waiting for Honeycomb comparison evidence: %w", lastErr)
			}
			return lastEvidence, fmt.Errorf("waiting for Honeycomb comparison evidence: got %d of %d spans: %w", lastEvidence.Count, attempted, ctx.Err())
		case <-timer.C:
		}
	}
}

func readRolloutHealth(ctx context.Context, client *honeycombMCPClient, options monitorOptions) (rolloutHealth, error) {
	producer, err := client.runQuery(ctx, map[string]any{
		"environment_slug": "staging",
		"dataset_slug":     "transfers",
		"query_spec": map[string]any{
			"calculations": []map[string]any{
				{"op": "MAX", "column": "current_int_row_offset", "name": "checkpoint"},
				{"op": "COUNT", "name": "iterations"},
				{"op": "COUNT", "name": "errors", "filters": []map[string]any{{"column": "error", "op": "=", "value": true}}},
			},
			"filters": []map[string]any{
				{"column": "name", "op": "=", "value": "backfill-iteration"},
				{"column": "service.version", "op": "=", "value": options.ProducerVersion},
			},
			"time_range": "30m",
		},
	})
	if err != nil {
		return rolloutHealth{}, fmt.Errorf("querying producer health: %w", err)
	}
	consumer, err := client.runQuery(ctx, map[string]any{
		"environment_slug": "staging",
		"dataset_slug":     "transfersbff2",
		"query_spec": map[string]any{
			"calculations": []map[string]any{
				{"op": "COUNT", "name": "upserts"},
				{"op": "COUNT", "name": "errors", "filters": []map[string]any{{"column": "error", "op": "=", "value": true}}},
			},
			"filters": []map[string]any{
				{"column": "name", "op": "=", "value": "projection-upsert-transfer"},
				{"column": "service.version", "op": "=", "value": options.ServiceVersion},
			},
			"time_range": "30m",
		},
	})
	if err != nil {
		return rolloutHealth{}, fmt.Errorf("querying consumer health: %w", err)
	}
	producerWorkflowErrors, err := readWorkflowErrorCount(ctx, client, "transfers", options.ProducerVersion,
		map[string]any{"column": "name", "op": "in", "value": []string{
			"backfill-produce-transfer-updated", "process-backfill", "backfill-iteration",
			"get-issuing-data-for-backfill", "produce-transfer-updated-for-batch",
		}})
	if err != nil {
		return rolloutHealth{}, fmt.Errorf("querying Transfers backfill errors: %w", err)
	}
	consumerWorkflowErrors, err := readWorkflowErrorCount(ctx, client, "transfersbff2", options.ServiceVersion,
		map[string]any{"column": "name", "op": "in", "value": []string{
			"consuming", "consume-event", "projection-upsert-transfer",
			"projection-update-transfer-failure-reason", "projection-emit-projection-transfer-event",
		}})
	if err != nil {
		return rolloutHealth{}, fmt.Errorf("querying transfersbff2 backfill-consumer errors: %w", err)
	}
	if len(producer.Rows) != 1 || len(consumer.Rows) != 1 {
		return rolloutHealth{}, fmt.Errorf("unexpected health result rows: producer=%d consumer=%d", len(producer.Rows), len(consumer.Rows))
	}
	health := rolloutHealth{
		ProducerQueryPK:       producer.RunPK,
		ConsumerQueryPK:       consumer.RunPK,
		ProducerErrorsQueryPK: producerWorkflowErrors.RunPK,
		ConsumerErrorsQueryPK: consumerWorkflowErrors.RunPK,
	}
	if health.ProducerIterations, err = number(producer.Rows[0], "iterations"); err != nil {
		return rolloutHealth{}, err
	}
	if health.ProducerIterations > 0 {
		if health.ProducerCheckpoint, err = number(producer.Rows[0], "checkpoint"); err != nil {
			return rolloutHealth{}, err
		}
	}
	if health.ProducerErrors, err = number(producer.Rows[0], "errors"); err != nil {
		return rolloutHealth{}, err
	}
	if health.ConsumerUpserts, err = number(consumer.Rows[0], "upserts"); err != nil {
		return rolloutHealth{}, err
	}
	if health.ConsumerErrors, err = number(consumer.Rows[0], "errors"); err != nil {
		return rolloutHealth{}, err
	}
	if health.ProducerWorkflowErrors, err = number(producerWorkflowErrors.Rows[0], "errors"); err != nil {
		return rolloutHealth{}, err
	}
	if health.ConsumerWorkflowErrors, err = number(consumerWorkflowErrors.Rows[0], "errors"); err != nil {
		return rolloutHealth{}, err
	}
	return health, nil
}

func readWorkflowErrorCount(ctx context.Context, client *honeycombMCPClient, dataset, serviceVersion string, workflowFilters ...map[string]any) (honeycombQueryResult, error) {
	filters := []map[string]any{
		{"column": "error", "op": "=", "value": true},
		{"column": "service.version", "op": "=", "value": serviceVersion},
	}
	filters = append(filters, workflowFilters...)
	result, err := client.runQuery(ctx, map[string]any{
		"environment_slug": "staging",
		"dataset_slug":     dataset,
		"query_spec": map[string]any{
			"calculations": []map[string]any{{"op": "COUNT", "name": "errors"}},
			"filters":      filters,
			"time_range":   "30m",
		},
	})
	if err != nil {
		return honeycombQueryResult{}, err
	}
	if len(result.Rows) != 1 {
		return honeycombQueryResult{}, fmt.Errorf("unexpected workflow error result rows: %d", len(result.Rows))
	}
	return result, nil
}

func readComparisonEvidence(ctx context.Context, client *honeycombMCPClient, serviceVersion, userAgent string) (comparisonEvidence, error) {
	query, err := client.runQuery(ctx, map[string]any{
		"environment_slug": "staging",
		"dataset_slug":     "transfersbff2",
		"query_spec": map[string]any{
			"calculations": []map[string]any{{"op": "COUNT", "name": "comparisons"}},
			"filters": []map[string]any{
				{"column": "root.user_agent.original", "op": "=", "value": userAgent},
				{"column": "name", "op": "=", "value": "compare-api-transfer-for-ops"},
			},
			"breakdowns": []string{"projection_comparison", "trace.trace_id", "service.version"},
			"time_range": "15m",
			"limit":      20,
		},
	})
	if err != nil {
		return comparisonEvidence{}, fmt.Errorf("querying comparison evidence: %w", err)
	}
	evidence := comparisonEvidence{QueryRunPK: query.RunPK}
	for _, row := range query.Rows {
		count, numberErr := number(row, "comparisons")
		if numberErr != nil {
			return comparisonEvidence{}, numberErr
		}
		evidence.Count += count
		comparison, _ := row["projection_comparison"].(string)
		switch comparison {
		case "match":
			evidence.Matches += count
		case "mismatch":
			evidence.Mismatches += count
		case "error":
			evidence.Errors += count
		}
		if traceID, ok := row["trace.trace_id"].(string); ok && traceID != "" {
			evidence.TraceID = traceID
		}
		if version, ok := row["service.version"].(string); ok && version != "" && version != serviceVersion {
			return comparisonEvidence{}, fmt.Errorf("comparison evidence came from service version %s, expected %s", version, serviceVersion)
		}
	}
	return evidence, nil
}

func (e comparisonEvidence) requireClean(attempted int, serviceVersion string) error {
	if e.Count != int64(attempted) || e.Matches != int64(attempted) || e.Mismatches != 0 || e.Errors != 0 {
		return fmt.Errorf("Honeycomb evidence is not clean for %s: attempted=%d comparisons=%d matches=%d mismatches=%d errors=%d query=%s",
			serviceVersion, attempted, e.Count, e.Matches, e.Mismatches, e.Errors,
			honeycombQueryURL("transfersbff2", e.QueryRunPK))
	}
	if e.TraceID == "" {
		return errors.New("Honeycomb comparison evidence did not include a trace ID")
	}
	return nil
}

func acquireMonitorLock(path string) (func(), error) {
	file, err := syscall.Open(path, syscall.O_CREAT|syscall.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("opening monitor lock: %w", err)
	}
	if err := syscall.Flock(file, syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = syscall.Close(file)
		if errors.Is(err, syscall.EWOULDBLOCK) {
			return nil, errors.New("projection backfill monitor is already running")
		}
		return nil, fmt.Errorf("locking monitor: %w", err)
	}
	return func() {
		_ = syscall.Flock(file, syscall.LOCK_UN)
		_ = syscall.Close(file)
	}, nil
}

func honeycombQueryURL(dataset, runPK string) string {
	return fmt.Sprintf("https://ui.honeycomb.io/moov/environments/staging/datasets/%s/result/%s", dataset, runPK)
}
