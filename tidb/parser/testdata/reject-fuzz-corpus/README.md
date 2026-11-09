# TiDB reject-fuzz corpus

The differential mutation fuzz sweep of the §3.6 rejection suite
(`tidb/parser/reject_fuzz_test.go`, `TestTiDBRejectFuzz`). It mutates the valid
(accept/accept) statements of the §3.3 completeness corpus and the §3.4 real-SQL
samples at the token and byte level, runs omni and live TiDB v8.5.5 on each
mutant, and flags every **over-acceptance** (omni accepts, TiDB rejects).

Where `parser_y_negatives_corpus.tsv` enumerates one-step-off near-misses of each
grammar alternative, this sweep reaches malformed strings enumeration never will.

## How it runs

Deterministic: a fixed seed (`rejectFuzzSeed`) so a rerun produces the same
corpus. Size is `TIDB_REJECT_FUZZ_N` (local default small; CI dials up). The gate
is set-equality against `known-mismatches.txt`; `TIDB_REJECT_FUZZ_DEFER=1`
downgrades it to log-and-persist for the larger CI batch, so a fresh finding
does not auto-fail before a human triages it. Each over-acceptance is minimised
(token-level ddmin) to the shortest disagreeing input.

The gating direction is over-acceptance. A mutant omni rejects and TiDB accepts
is a §3.3 completeness gap, logged not gated.

## Files

- `known-mismatches.txt` — the baseline allowlist, one `family<TAB>sql` line per
  minimised over-acceptance the sweep currently finds. A NEW over-acceptance (not
  listed) fails the per-PR run; a listed one that stops reproducing also fails
  (re-baseline after a fix).
- `mismatches.txt` — written every run with the current findings, for triage.

## Workflow when it finds something

1. Minimise is automatic; read `mismatches.txt`.
2. Fix the over-acceptance across `mysql`/`mariadb`/`oceanbase` like a gap,
   measured on live TiDB.
3. Pin the minimised input into `parser_y_negatives_corpus.tsv` as a permanent
   regression row (`fuzz#N` label), and drop it from `known-mismatches.txt`.
