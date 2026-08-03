# Staging LIST transfers benchmarks

This suite compares identical `LIST /accounts/{accountID}/transfers` requests before and after enabling transfersbff2 projection reads in staging. It is based on the query-plan research in [Transfers LIST query index strategy](https://ampcode.com/threads/T-019ef1bb-e305-724a-93ba-8f7e75a96ff8) and the implementation merged in [transfersbff2#924](https://github.com/moovfinancial/transfersbff2/pull/924).

The runner is read-only. It does not create, mutate, cancel, or refund transfers. It refuses to run unless `MOOV_HOST` contains `staging`.

## Accounts and fixtures

The fixture values were queried from `moov-data-staging.transfers_public.transfers` on 2026-08-03. Keep them fixed between the baseline and projection runs.

| Account | Transfer rows | First transfer | Latest transfer | Useful sparse fixtures |
| --- | ---: | --- | --- | --- |
| Moov Financial `db04bf9d-91f6-4206-ba38-6844636532ad` | 261,129 | 2021-08-13 | 2026-08-03 | group, schedule, payment link, foreign ID, refunds |
| High-volume staging partner `6f36ae5e-d715-4b5f-9089-11dd8f21c120` | 86,892,676 | 2022-03-04 | 2026-08-03 | none found |

The high-volume account had no group, schedule, payment-link, foreign-ID, refund, dispute, authorization, or capture fixture. Its sparse cases deliberately use valid no-match values. This still measures the relevant query shape and zero-result fast-out, but it is not equivalent to a matching sparse lookup.

Moov Financial also had no `authorization_id`, `capture_id`, or disputed transfer in this table. Those cases are intentional no-match requests for both accounts.

## Query matrix

The runner defines 34 stable public API cases:

- **Common paths:** default account feed and the calling-account-scoped `/transfers` endpoint.
- **Pagination:** explicit counts 5, 10, 20, and 1000; count-20 offsets 140 and 200 observed in slow staging traces; shallow skip 200 with count 200; deep skip 50,000; and skip beyond the account's known row count.
- **Indexed filters:** completed status; failed status with page-2 pagination; recent, old narrow, and broad date windows; status plus date.
- **Account filters:** source match, destination match, source/destination union, and no match.
- **Authorization guard:** list a destination account while authenticating as its transfer's partner. This exercises the account-path plus calling-account party guard that review added to preserve legacy authorization behavior.
- **Sparse selectors:** group, schedule, payment link, foreign ID, authorization ID, and capture ID.
- **Sparse booleans:** refunded, disputed, and refunded plus disputed.
- **Combined selectors:** status plus group; status plus account plus date.
- **Empty results:** absent foreign ID and skip past the end.

The common indexed paths correspond to the ordinary account/status/type/date B-tree indexes introduced by transfersbff2#924. Sparse cases correspond to selective indexes or fallback predicates whose real traffic was too rare to justify broad combination indexes. Deep offset remains intentionally included because keyset pagination, not another index, is the structural fix for very deep pages.

All cases use `X-Moov-Version: v2026.07.00`. The optimization research also measured `transferTypes`, but that parameter only exists in the dev API shape. Staging returned `404` for the `v2026.10.00` route on 2026-08-03, while older versions ignore the parameter, so it is intentionally excluded until that public API version is routable. Adding an ignored filter would produce a misleading benchmark of the unfiltered query.

Inspect the exact generated query strings without sending requests:

```sh
go run ./benchmarks/list-transfers-staging -list
go run ./benchmarks/list-transfers-staging -account 6f36ae5e-d715-4b5f-9089-11dd8f21c120 -list
```

## Run protocol

Load the staging username and credential outside the repository. No secret is written to result files.

```fish
moov_env staging --quiet

go run ./benchmarks/list-transfers-staging \
  -account db04bf9d-91f6-4206-ba38-6844636532ad \
  -label baseline-public \
  -warmups 1 \
  -iterations 3
```

After the projection feature flag is enabled, rerun the same command with only the label changed:

```fish
go run ./benchmarks/list-transfers-staging \
  -account db04bf9d-91f6-4206-ba38-6844636532ad \
  -label projection-public \
  -warmups 1 \
  -iterations 3
```

Use `-case '^(base|page|status)'` to run a subset. Requests run sequentially to avoid client-side contention, use a 30-second timeout, and continue after individual failures. The command exits nonzero if any measured request failed, but still writes the complete result artifact. Every measured result records the exact case definition, timestamp, duration, HTTP status, response count and size, and `X-Request-ID`. The run label is also appended to the benchmark user agent, so Honeycomb can break down `baseline` and `projection` directly. Result JSON is written under `results/` and can be correlated exactly by request ID.

Load the dedicated high-volume account credential before running that fixture. Do not run its suite with Moov Financial credentials: transfersbff2 requires the calling account to be a party to every returned transfer.

```fish
moov_env staging 6f36ae5e-d715-4b5f-9089-11dd8f21c120 "gatling/load testing" --quiet

go run ./benchmarks/list-transfers-staging \
  -account 6f36ae5e-d715-4b5f-9089-11dd8f21c120 \
  -label baseline-public-gatling \
  -warmups 1 \
  -iterations 3
```

## Honeycomb comparison

Open the [staging LIST benchmark query](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/o3EMPfXEJJ5). It selects root `http` spans for the benchmark user-agent prefix and Moov Financial calling account, then breaks down by:

- run label in `user_agent.original`;
- exact query shape in `http.uri`;
- API version in `http.versioning.used`;
- response status in `http.response.status_code`.

The query reports request count and P50/P95/P99 `duration_ms`, ordered by the slowest P95. It intentionally does not use average latency. Refresh or widen its seven-day window after the projection run. Use the request IDs in a result JSON artifact when an exact trace-level correlation is needed.

### Historical slow-query coverage

A [60-day staging query grouped by normalized LIST shape](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/oFLUT3MvKkT) found these notable successful-request tails on `v2026.07.00`:

| Shape | Requests | P50 | P95 | P99 | Max | Benchmark coverage |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| status + pagination | 11 | 59 ms | 21.8 s | 21.8 s | 21.8 s | `status-failed-page-2` |
| count 1000 | 8 | 699 ms | 16.8 s | 16.8 s | 16.8 s | `page-large-1000` |
| payment link | 111 | 115 ms | 625 ms | 3.4 s | 3.9 s | `sparse-payment-link` |
| pagination offset | 2,734 | 45 ms | 212 ms | 636 ms | 64.4 s | count-20 offsets plus shallow, deep, and past-end cases |

The [exact requests over two seconds](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/HSPR22D1MwD) supplied the additional parameter combinations. The two worst missing offset shapes were `skip=200&count=20` at 64.4 seconds and `skip=140&count=20` at 58.5 seconds. Trace inspection attributed 64.35 seconds and 58.47 seconds respectively to the transfers service's list SQL. A `status=failed&skip=20&count=20` request took 21.8 seconds, with 21.74 seconds in the same list SQL path. Explicit count-5 and count-20 pages also appeared among the slow successful requests and are retained as separate cases even though count 20 is the API default.

The historical scan also found a 3.75-second `captureID` request on `v2024.00.00`. The suite already covers the current `v2026.07.00` `captureIDs` shape. It does not add the old singular parameter because mixing API versions would make the projection comparison ambiguous.

## BigQuery fixture query

The fixture scan processed about 14.2 GB and used an explicit 16 GB billing cap:

```sql
SELECT
  account_id,
  COUNT(*) AS transfer_count,
  MIN(created_on) AS first_transfer_on,
  MAX(created_on) AS latest_transfer_on,
  ARRAY_AGG(DISTINCT api_transfer_status IGNORE NULLS ORDER BY api_transfer_status) AS api_statuses,
  ARRAY_AGG(DISTINCT transfer_type IGNORE NULLS ORDER BY transfer_type) AS transfer_types,
  ARRAY_AGG(IF(source_account_id IS NOT NULL, STRUCT(created_on, source_account_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].source_account_id AS source_account_id,
  ARRAY_AGG(IF(destination_account_id IS NOT NULL, STRUCT(created_on, destination_account_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].destination_account_id AS destination_account_id,
  ARRAY_AGG(IF(group_id IS NOT NULL, STRUCT(created_on, group_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].group_id AS group_id,
  ARRAY_AGG(IF(schedule_id IS NOT NULL, STRUCT(created_on, schedule_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].schedule_id AS schedule_id,
  ARRAY_AGG(IF(checkout_code IS NOT NULL, STRUCT(created_on, checkout_code), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].checkout_code AS payment_link_code,
  ARRAY_AGG(IF(foreign_id IS NOT NULL, STRUCT(created_on, foreign_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].foreign_id AS foreign_id,
  ARRAY_AGG(IF(authorization_id IS NOT NULL, STRUCT(created_on, authorization_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].authorization_id AS authorization_id,
  ARRAY_AGG(IF(capture_id IS NOT NULL, STRUCT(created_on, capture_id), NULL) IGNORE NULLS ORDER BY created_on DESC LIMIT 1)[SAFE_OFFSET(0)].capture_id AS capture_id,
  COUNTIF(refunded_amount > 0) AS refunded_count,
  COUNTIF(disputed_amount > 0) AS disputed_count
FROM `moov-data-staging.transfers_public.transfers`
WHERE account_id IN (
  'db04bf9d-91f6-4206-ba38-6844636532ad',
  '6f36ae5e-d715-4b5f-9089-11dd8f21c120'
)
GROUP BY account_id;
```
