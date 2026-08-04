package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultBackfillStartPK = int64(3001)
	defaultBackfillMaxPK   = int64(89_821_930)
	defaultRangeSize       = int64(10_000)
	defaultAPIVersion      = "v2026.07.00"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usageError()
	}

	switch args[0] {
	case "prepare":
		return runPrepare(args[1:])
	case "validate":
		return runValidate(args[1:])
	case "summary":
		return runSummary(args[1:])
	default:
		return usageError()
	}
}

func runPrepare(args []string) error {
	flags := flag.NewFlagSet("prepare", flag.ContinueOnError)
	dbPath := flags.String("db", defaultDBPath(), "SQLite campaign database")
	startPK := flags.Int64("start-pk", defaultBackfillStartPK, "first PK included in the automated backfill")
	maxPK := flags.Int64("max-pk", defaultBackfillMaxPK, "last PK included in the backfill")
	rangeSize := flags.Int64("range-size", defaultRangeSize, "PKs advanced by each scheduled invocation")
	maxBytes := flags.Int64("maximum-bytes-billed", 50_000_000_000, "BigQuery billing limit")
	dryRunOnly := flags.Bool("dry-run-only", false, "validate the BigQuery query without downloading samples")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("prepare accepts flags only")
	}

	return prepareSamples(prepareOptions{
		DBPath:             *dbPath,
		StartPK:            *startPK,
		MaxPK:              *maxPK,
		RangeSize:          *rangeSize,
		MaximumBytesBilled: *maxBytes,
		DryRunOnly:         *dryRunOnly,
	})
}

func runValidate(args []string) error {
	flags := flag.NewFlagSet("validate", flag.ContinueOnError)
	dbPath := flags.String("db", defaultDBPath(), "SQLite campaign database")
	checkpoint := flags.Int64("checkpoint", 0, "fully consumed backfill PK checkpoint")
	version := flags.String("version", defaultAPIVersion, "X-Moov-Version value")
	serviceVersion := flags.String("service-version", "", "deployed transfersbff2 image or version")
	maxRanges := flags.Int("max-ranges", 1, "maximum pending PK ranges to validate")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("validate accepts flags only")
	}
	if *checkpoint <= 0 {
		return errors.New("checkpoint must be positive")
	}
	if *maxRanges < 1 {
		return errors.New("max-ranges must be positive")
	}

	return validatePending(validateOptions{
		DBPath:         *dbPath,
		Checkpoint:     *checkpoint,
		Version:        *version,
		ServiceVersion: *serviceVersion,
		MaxRanges:      *maxRanges,
	})
}

func runSummary(args []string) error {
	flags := flag.NewFlagSet("summary", flag.ContinueOnError)
	dbPath := flags.String("db", defaultDBPath(), "SQLite campaign database")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("summary accepts flags only")
	}
	return printSummary(*dbPath)
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "projection-backfill-staging.sqlite"
	}
	return filepath.Join(home, ".local", "state", "moov", "projection-backfill-staging.sqlite")
}

func usageError() error {
	return errors.New("usage: go run ./benchmarks/ops-batch-staging <prepare|validate|summary> [flags]")
}
