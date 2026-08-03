# Moov Financial public LIST baseline

Run: `baseline-public`, 2026-08-03 22:21:40–22:25:13 UTC. The complete machine-readable result is [`results/baseline-public-db04bf9d-20260803T222140Z.json`](results/baseline-public-db04bf9d-20260803T222140Z.json).

The run sent one warmup and three measured requests for each of 29 cases, sequentially, against Moov Financial staging account `db04bf9d-91f6-4206-ba38-6844636532ad`. Honeycomb observed all 116 attempts: 112 account-path spans and four global-path spans. The [projection-route verification](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/fm2rnWDudGS) confirms every request used `projection_enabled=false`.

## Baseline highlights

Medians below exclude requests that hit the client-side 30-second timeout.

| Case | Median | Range | Outcome |
| --- | ---: | ---: | --- |
| default account feed | 250 ms | 247–253 ms | 2 successes, 1 timeout |
| global `/transfers` feed | 497 ms | 367–540 ms | 3 successes |
| count 10 | 247 ms | 219–274 ms | 2 successes, 1 timeout |
| count 1000 | 794 ms | 769–940 ms | 3 successes |
| shallow skip 200 | 321 ms | 281–481 ms | 3 successes |
| deep skip 50,000 | 721 ms | 698–814 ms | 3 successes |
| skip past end | 2,713 ms | 2,372–28,109 ms | 3 successes; high variance |
| account no-match | 101 ms | 98–105 ms | 3 successes |
| payment-link match | 502 ms | 467–536 ms | 3 successes |
| status + account + date | 113 ms | 110–121 ms | 3 successes |
| calling-account party guard | 240 ms | 186–294 ms | 2 successes, 1 timeout |

Three of 87 measured requests timed out: default feed, count 10, and the party-guard case. Because unrelated query shapes each stalled once while their other iterations completed in 186–294 ms, treat these as staging outliers until trace inspection proves they are query-specific. The skip-past-end case is independently concerning: all iterations completed, but latency varied from 2.4 to 28.1 seconds, consistent with the optimization research's warning that deep offset needs keyset pagination rather than another index.

Open the [latency comparison query](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/o3EMPfXEJJ5) for P50/P95/P99 by exact URI, run label, API version, and response status. Rerun the same 29 cases with label `projection-public` after enabling the projection feature flag.

## Historical slow-shape supplement

The 60-day Honeycomb review added five cases after the original baseline. Run `baseline-historical-slow` on 2026-08-03 22:45:31–22:46:06 UTC recorded one warmup and three measured requests for each case. The complete artifact is [`results/baseline-historical-slow-db04bf9d-20260803T224531Z.json`](results/baseline-historical-slow-db04bf9d-20260803T224531Z.json).

| Case | Median | Range | Outcome |
| --- | ---: | ---: | --- |
| explicit count 5 | 154 ms | 152–172 ms | 3 successes |
| explicit count 20 | 192 ms | 184–246 ms | 3 successes |
| skip 140, count 20 | 167 ms | 166–267 ms | 3 successes |
| skip 200, count 20 | 191 ms | 171–212 ms | 2 successes, 1 timeout |
| failed status, page 2 | 304 ms | 284–317 ms | 3 successes |

Honeycomb observed all 20 attempts, including warmups. The [labeled latency query](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/giATAVgvjSR) shows that successful server spans were normally 78–230 ms at P50 by URI; the timed-out skip-200 request became a 29.98-second `499`. The [projection-route verification](https://ui.honeycomb.io/moov/environments/staging/datasets/transfersbff2/result/kcc6xGNLSVH) confirms all 20 attempts used `projection_enabled=false` and `account-transfer.list`.

After enabling projection reads, rerun these five cases with label `projection-historical-slow` as well as rerunning the complete 34-case suite. Keep the timeout in the baseline artifact: it is the behavior the historical scan was designed to preserve for comparison.
