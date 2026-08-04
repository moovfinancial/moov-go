package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const schema = `
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS campaign (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    start_pk INTEGER NOT NULL,
    max_pk INTEGER NOT NULL,
    range_size INTEGER NOT NULL,
    sample_seed TEXT NOT NULL,
    source_table TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS samples (
    pk_id INTEGER PRIMARY KEY,
    transfer_id TEXT NOT NULL UNIQUE,
    range_start INTEGER NOT NULL,
    range_end INTEGER NOT NULL,
    sample_rank INTEGER NOT NULL,
    selection_reason TEXT NOT NULL,
    transfer_type TEXT,
    status TEXT,
    account_id TEXT
);
CREATE INDEX IF NOT EXISTS samples_range_idx ON samples(range_start, range_end);
CREATE TABLE IF NOT EXISTS validation_runs (
    id INTEGER PRIMARY KEY,
    range_start INTEGER NOT NULL,
    range_end INTEGER NOT NULL,
    started_at TEXT NOT NULL,
    api_version TEXT NOT NULL,
    service_version TEXT,
    user_agent TEXT NOT NULL,
    request_id TEXT,
    trace_id TEXT,
    duration_ms INTEGER,
    http_status INTEGER,
    attempted_count INTEGER NOT NULL,
    status TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS validation_runs_range_idx ON validation_runs(range_start, range_end, status);
CREATE TABLE IF NOT EXISTS validation_results (
    run_id INTEGER NOT NULL,
    pk_id INTEGER NOT NULL,
    comparison TEXT NOT NULL,
    diff_paths TEXT,
    error TEXT,
    PRIMARY KEY (run_id, pk_id),
    FOREIGN KEY (run_id) REFERENCES validation_runs(id),
    FOREIGN KEY (pk_id) REFERENCES samples(pk_id)
);
`

func initializeDB(dbPath string) error {
	if err := requireCommand("sqlite3"); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return fmt.Errorf("creating database directory: %w", err)
	}
	if _, err := runSQLite(dbPath, schema); err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	return nil
}

func requireCampaignDB(dbPath string) error {
	if err := requireCommand("sqlite3"); err != nil {
		return err
	}
	info, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("campaign database %s does not exist; run prepare first", dbPath)
		}
		return fmt.Errorf("checking campaign database: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("campaign database %s is not a regular file", dbPath)
	}
	campaigns, err := scalar(dbPath, "SELECT COUNT(*) FROM campaign;")
	if err != nil || campaigns != "1" {
		return fmt.Errorf("campaign database %s is not initialized; run prepare first", dbPath)
	}
	samples, err := scalar(dbPath, "SELECT COUNT(*) FROM samples;")
	if err != nil || samples == "0" {
		return fmt.Errorf("campaign database %s has no samples; run prepare first", dbPath)
	}
	return nil
}

func runSQLite(dbPath, input string, args ...string) ([]byte, error) {
	cmdArgs := append([]string{"-batch", "-bail"}, args...)
	cmdArgs = append(cmdArgs, dbPath)
	cmd := exec.Command("sqlite3", cmdArgs...)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("sqlite3: %w: %s", err, bytes.TrimSpace(stderr.Bytes()))
	}
	return stdout.Bytes(), nil
}

func scalar(dbPath, query string) (string, error) {
	output, err := runSQLite(dbPath, query, "-noheader")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func requireCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s is required: %w", name, err)
	}
	return nil
}

func sqlQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func sqliteDotQuote(path string) string {
	return "'" + strings.ReplaceAll(path, "'", "''") + "'"
}
