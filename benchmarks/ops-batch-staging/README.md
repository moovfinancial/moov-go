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

BigQuery is queried once for the fixed rollout range. Only selected `pk_id` and `transfer_id` mappings and their analysis dimensions are stored locally.

## Prepare

First verify the query and billing bound without downloading rows:

```sh
go run ./benchmarks/ops-batch-staging prepare -dry-run-only
```

After confirming the active backfill start checkpoint and maximum in Infra, prepare the manifest once:

```sh
go run ./benchmarks/ops-batch-staging prepare \
  -start-pk 3001 \
  -max-pk 89821930 \
  -range-size 10000 \
  -maximum-bytes-billed 50000000000
```

Change `start-pk` if another manual backfill advances the checkpoint before the automated schedule deploys. Preparation records the immutable start, maximum, range size, source table, and sample seed in SQLite. It refuses to reuse a database that already contains a manifest; use a separate `-db` path for another campaign.

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

The summary reports planned samples, completed ranges, total attempts, matches, mismatches, per-transfer errors, and request-level failures. Historical errors remain recorded after a successful retry.

The command does not monitor Honeycomb or invoke itself on a timer. Start operator-driven validation only after the automatic backfill deployment is confirmed.
