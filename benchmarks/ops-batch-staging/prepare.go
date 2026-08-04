package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const (
	sourceTransfersTable = "moov-data-staging.transfers_public.transfers"
	sampleSeed           = "projection-backfill-sample-v1"
)

type prepareOptions struct {
	DBPath             string
	StartPK            int64
	MaxPK              int64
	RangeSize          int64
	MaximumBytesBilled int64
	DryRunOnly         bool
}

func prepareSamples(options prepareOptions) error {
	if options.StartPK < 1 || options.MaxPK < options.StartPK {
		return errors.New("invalid PK range")
	}
	if options.RangeSize < 1 || options.MaximumBytesBilled < 1 {
		return errors.New("range-size and maximum-bytes-billed must be positive")
	}
	if err := requireCommand("bq"); err != nil {
		return err
	}

	query := sampleQuery(options.StartPK, options.MaxPK, options.RangeSize)
	if err := runBigQueryDryRun(query, options.MaximumBytesBilled); err != nil {
		return err
	}
	if options.DryRunOnly {
		fmt.Println("BigQuery sample query dry run passed")
		return nil
	}
	if err := initializeDB(options.DBPath); err != nil {
		return err
	}
	if err := ensureCampaign(options); err != nil {
		return err
	}
	existing, err := scalar(options.DBPath, "SELECT COUNT(*) FROM samples;")
	if err != nil {
		return err
	}
	if existing != "0" {
		return fmt.Errorf("campaign database already contains %s samples; use a new -db path to prepare another manifest", existing)
	}

	temp, err := os.CreateTemp("", "projection-backfill-samples-*.csv")
	if err != nil {
		return fmt.Errorf("creating sample export: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	defer temp.Close()

	if err := exportSamples(query, options.MaximumBytesBilled, temp); err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("closing sample export: %w", err)
	}
	if err := importSamples(options.DBPath, tempPath); err != nil {
		return err
	}

	count, err := scalar(options.DBPath, "SELECT COUNT(*) FROM samples;")
	if err != nil {
		return err
	}
	fmt.Printf("prepared %s samples in %s\n", count, options.DBPath)
	return nil
}

func ensureCampaign(options prepareOptions) error {
	commands := fmt.Sprintf(`
INSERT INTO campaign (id, start_pk, max_pk, range_size, sample_seed, source_table)
SELECT 1, %d, %d, %d, %s, %s
WHERE NOT EXISTS (SELECT 1 FROM campaign);
`, options.StartPK, options.MaxPK, options.RangeSize, sqlQuote(sampleSeed), sqlQuote(sourceTransfersTable))
	if _, err := runSQLite(options.DBPath, commands); err != nil {
		return fmt.Errorf("recording campaign configuration: %w", err)
	}
	matching, err := scalar(options.DBPath, fmt.Sprintf(`
SELECT COUNT(*) FROM campaign
WHERE id = 1 AND start_pk = %d AND max_pk = %d AND range_size = %d
  AND sample_seed = %s AND source_table = %s;
`, options.StartPK, options.MaxPK, options.RangeSize, sqlQuote(sampleSeed), sqlQuote(sourceTransfersTable)))
	if err != nil {
		return err
	}
	if matching != "1" {
		return errors.New("campaign database was prepared with different sampling options; use a new -db path")
	}
	return nil
}

func runBigQueryDryRun(query string, maximumBytesBilled int64) error {
	args := append([]string{"query", "--dry_run"}, bigQueryArgs(query, maximumBytesBilled)...)
	cmd := exec.Command("bq", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("dry-running BigQuery sample query: %w: %s", err, bytes.TrimSpace(output))
	}
	return nil
}

func exportSamples(query string, maximumBytesBilled int64, output *os.File) error {
	args := append([]string{"query", "--format=csv", "--max_rows=5000000"}, bigQueryArgs(query, maximumBytesBilled)...)
	cmd := exec.Command("bq", args...)
	cmd.Stdout = output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_, _ = output.Seek(0, io.SeekStart)
		stdout, _ := io.ReadAll(io.LimitReader(output, 4<<10))
		details := bytes.TrimSpace(append(stderr.Bytes(), stdout...))
		return fmt.Errorf("exporting BigQuery samples: %w: %s", err, details)
	}
	return nil
}

func bigQueryArgs(query string, maximumBytesBilled int64) []string {
	return []string{
		"--project_id=moov-data-staging",
		"--use_legacy_sql=false",
		fmt.Sprintf("--maximum_bytes_billed=%d", maximumBytesBilled),
		query,
	}
}

func importSamples(dbPath, csvPath string) error {
	commands := fmt.Sprintf(`
CREATE TEMP TABLE samples_import (
    pk_id INTEGER,
    transfer_id TEXT,
    range_start INTEGER,
    range_end INTEGER,
    sample_rank INTEGER,
    selection_reason TEXT,
    transfer_type TEXT,
    status TEXT,
    account_id TEXT
);
.mode csv
.import --skip 1 %s samples_import
BEGIN;
INSERT INTO samples
SELECT pk_id, transfer_id, range_start, range_end, sample_rank, selection_reason, transfer_type, status, account_id
FROM samples_import;
COMMIT;
`, sqliteDotQuote(csvPath))
	if _, err := runSQLite(dbPath, commands); err != nil {
		return fmt.Errorf("importing samples: %w", err)
	}
	return nil
}

func sampleQuery(startPK, maxPK, rangeSize int64) string {
	return fmt.Sprintf(`
WITH source AS (
  SELECT
    pk_id,
    transfer_id,
    transfer_type,
    COALESCE(api_transfer_status, status) AS api_status,
    account_id,
    DIV(pk_id - %d, %d) AS range_index,
    FARM_FINGERPRINT(CONCAT(transfer_id, ':%s')) AS stable_hash,
    CASE
      WHEN authorization_id IS NOT NULL OR capture_id IS NOT NULL THEN 'auth-capture'
      WHEN refunded_amount > 0 THEN 'refunded'
      WHEN disputed_amount > 0 THEN 'disputed'
      WHEN line_items IS NOT NULL AND JSON_TYPE(line_items) != 'null' THEN 'line-items'
      WHEN parent_id IS NOT NULL OR group_id IS NOT NULL THEN 'linked'
      WHEN NOT is_public THEN 'non-public'
      ELSE 'uniform'
    END AS risk_reason
  FROM %s
  WHERE pk_id BETWEEN %d AND %d
), ranked AS (
  SELECT
    *,
    %d + range_index * %d AS range_start,
    LEAST(%d, %d + (range_index + 1) * %d - 1) AS range_end,
    ROW_NUMBER() OVER (PARTITION BY range_index ORDER BY pk_id) AS first_rank,
    ROW_NUMBER() OVER (PARTITION BY range_index ORDER BY pk_id DESC) AS last_rank,
    ROW_NUMBER() OVER (PARTITION BY range_index, risk_reason ORDER BY stable_hash, pk_id) AS risk_rank,
    ROW_NUMBER() OVER (
      PARTITION BY range_index, COALESCE(transfer_type, ''), COALESCE(api_status, '')
      ORDER BY stable_hash, pk_id
    ) AS stratum_rank,
    ROW_NUMBER() OVER (
      PARTITION BY range_index, COALESCE(account_id, '')
      ORDER BY stable_hash, pk_id
    ) AS account_rank
  FROM source
), prioritized AS (
  SELECT
    *,
    CASE
      WHEN first_rank = 1 OR last_rank = 1 THEN 0
      WHEN risk_reason != 'uniform' AND risk_rank <= 5 THEN 1
      WHEN stratum_rank = 1 THEN 2
      WHEN account_rank = 1 THEN 3
      ELSE 4
    END AS sample_priority,
    CASE
      WHEN first_rank = 1 THEN 'range-first'
      WHEN last_rank = 1 THEN 'range-last'
      WHEN risk_reason != 'uniform' AND risk_rank <= 5 THEN risk_reason
      WHEN stratum_rank = 1 THEN 'type-status-stratum'
      WHEN account_rank = 1 THEN 'account-stratum'
      ELSE 'uniform'
    END AS selection_reason
  FROM ranked
), sampled AS (
  SELECT
    *,
    ROW_NUMBER() OVER (
      PARTITION BY range_index
      ORDER BY sample_priority, stable_hash, pk_id
    ) AS sample_rank
  FROM prioritized
)
SELECT
  pk_id,
  transfer_id,
  range_start,
  range_end,
  sample_rank,
  selection_reason,
  transfer_type,
  api_status AS status,
  account_id
FROM sampled
WHERE sample_rank <= IF(MOD(range_index + 1, 10) = 0, 500, 100)
ORDER BY range_start, sample_rank
`, startPK, rangeSize, sampleSeed, "`"+sourceTransfersTable+"`", startPK, maxPK,
		startPK, rangeSize, maxPK, startPK, rangeSize)
}
