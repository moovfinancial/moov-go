# Staging projection backfill validation

This command prepares and runs deterministic sampled validation for the Transfers projection backfill. It compares the fully hydrated legacy and projection API mappings through `POST /ops/validate-api-projection-transfers` without validating every backfilled transfer.

It is staging-only. The generated SQLite database defaults to `~/.local/state/moov/projection-backfill-staging.sqlite` and must not be committed.

## Sampling policy

The manifest starts at the first PK after the checkpoint that precedes the automated rollout. Ranges align to each scheduled 10,000-PK invocation from that starting point, rather than to decimal 10,000 boundaries.

- Ordinary ranges select 100 transfers.
- Every tenth range selects 500 transfers.
- Selection always prefers the first and last existing transfer in a range.
- Bounded high-risk candidates cover authorization/capture, refunds, disputes, line items, linked transfers, and non-public transfers.
- One candidate per transfer-type/status stratum is preferred before the remaining sample is filled by a stable hash.
- One candidate per partner-account stratum is preferred after transfer-type/status coverage.
- The seed `projection-backfill-sample-v1` makes preparation repeatable.

The monitor reads the safe producer checkpoint from Honeycomb and extends the local sample cache only when that checkpoint moves beyond the cached PK range. Cache extensions are bounded, aligned to the campaign's 10,000-PK ranges, and persisted in SQLite. Only selected `pk_id` and `transfer_id` mappings and their analysis dimensions are stored locally.

## Prepare

Initialize the campaign without choosing a final PK:

```sh
go run ./benchmarks/ops-batch-staging prepare \
  -start-pk 3001 \
  -range-size 10000
```

The first monitor run reads the safe producer checkpoint from Honeycomb, queries BigQuery for at most one cache chunk, and records the cache checkpoint in SQLite. Later runs append samples only after the producer advances. The default chunk is 1,000,000 PKs and each query has a 50 GB billing cap; override these deliberately with `-sample-cache-chunk-pks` and `-maximum-bytes-billed` on `monitor`.

For a fixed one-time manifest, pass an explicit maximum:

```sh
go run ./benchmarks/ops-batch-staging prepare \
  -start-pk 3001 \
  -max-pk 89821930 \
  -range-size 10000 \
  -maximum-bytes-billed 50000000000
```

Change `start-pk` if another manual backfill advances the checkpoint before the automated schedule deploys. Preparation records the immutable start, range size, source table, and sample seed in SQLite. A zero maximum marks a dynamic campaign; `max-pk` remains available only for fixed exports. Use a separate `-db` path for another campaign.

The command requires authenticated `bq` and `sqlite3` CLIs. It explicitly queries `moov-data-staging.transfers_public.transfers` and refuses to exceed the configured BigQuery billing cap.

## Validate

Do not infer a safe validation checkpoint from elapsed time. Confirm in Honeycomb that Transfers advanced the checkpoint and that `transfersbff2-projection-backfill` consumed it without errors. Stay one completed range behind when catch-up is uncertain.

Load the existing staging API credentials without printing them:

```fish
moov_env staging --quiet
```

Validate one pending range at or below a confirmed consumed checkpoint:

```sh
go run ./benchmarks/ops-batch-staging validate \
  -checkpoint 13000 \
  -version v2026.07.00 \
  -service-version dev-add229b
```

The default `-max-ranges 1` prevents a stale operator from unexpectedly replaying a large validation backlog. Increase it deliberately after checking service health.

Each request records its range, sample IDs, API version, service version, User-Agent, duration, status, request ID, matches, mismatches, and per-transfer errors in one SQLite transaction. Clean ranges and any range containing a mismatch are not automatically repeated, including responses that contain both mismatches and errors. Request failures and error-only responses remain eligible for retry.

The runner stops on the first mismatch or error. Inspect that request in the staging `transfersbff2` Honeycomb dataset using its User-Agent:

```text
moov-go-projection-backfill-validator/1 range-{start}-{end}
```

Confirm the number of `projection_comparison` spans equals the attempted count. For mismatches, inspect `compare.diff_paths` and the corresponding legacy/projection field values before expanding the sample or continuing the rollout.

## Summary

```sh
go run ./benchmarks/ops-batch-staging summary
```

The summary reports the durable sample-cache PK checkpoint, planned samples, completed and intentionally skipped ranges, total attempts, matches, mismatches, per-transfer errors, and request-level failures. Historical errors remain recorded after a successful retry.

## Monitor every 30 minutes

The one-shot monitor is designed for an external scheduler such as `launchd`; it does not create a conversational or internal timer. It uses the Honeycomb MCP API key to check the current producer checkpoint and consumer health, stays one completed range behind the producer, extends its local deterministic sample cache toward that checkpoint, and selects the newest safe unvalidated range rather than replaying every range that elapsed since its previous run.

```fish
moov_env staging --quiet
go run ./benchmarks/ops-batch-staging monitor \
  -once \
  -producer-version dev-6ede7f1 \
  -service-version dev-add229b \
  -version v2026.07.00
```

The monitor:

- takes a non-blocking file lock so scheduled runs cannot overlap;
- refuses to send traffic when the expected producer or consumer image is not healthy;
- extends the SQLite sample cache in bounded BigQuery queries without requiring a final campaign PK;
- fails closed on errors in Transfers backfill stages, TransferUpdated consumption, or projection upserts;
- requires every selected transfer to have a matching projection-upsert span before validation;
- retries prior failed ranges before sampling a newer range;
- validates the newest safe range and records older elapsed ranges as intentionally skipped;
- records clean API responses as pending until Honeycomb confirms one matching comparison span per attempted transfer; and
- exits nonzero on a request, telemetry, mismatch, or consumer-health anomaly.

Use `-dry-run` to inspect health and range selection without loading Moov credentials or sending validation traffic. The local LaunchAgent wrapper reads the cached Honeycomb MCP and Moov staging tokens, so normal scheduled runs do not contact 1Password.

The command never invokes itself on a timer. Update the expected image flags whenever the staging rollout changes.
