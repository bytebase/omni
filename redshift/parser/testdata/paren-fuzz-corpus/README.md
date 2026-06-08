# paren-fuzz-corpus

Property-based fuzz corpus for SCENARIOS-pg-paren-dispatch.md §2.8. Driven by
`TestParenOracleFuzz` in `pg/parser/paren_oracle_fuzz_test.go` (build tag
`oracle`).

Files:

- `README.md` — this file; also guarantees the directory is committed.
- `seed-cases.txt` (optional) — one SQL statement per line. If present, the
  fuzz harness replays every line as an additional deterministic probe on top
  of the PRNG-generated batch. Use for known-mismatch regressions discovered
  during triage.
- `mismatches.txt` (generated) — written by `TestParenOracleFuzz` whenever a
  probe's PG/omni classifications disagree. Format per entry:

  ```
  --- mismatch N ---
  SQL:      <generated statement>
  pg:       accept|reject(SQLSTATE)
  omni:     subquery|joined_table|other|reject
  pg_err:   <PG error on reject, empty otherwise>
  omni_err: <omni parse error on reject, empty otherwise>
  ```

  The file is overwritten on every run, so it always reflects the latest
  mismatches. Check it into git when promoting a new finding into
  `seed-cases.txt`.

The harness is deterministic: same seed + same template set + same corpus
size => identical SQL batch. Bumping `fuzzSeed` or adding/removing templates
in `paren_oracle_fuzz_test.go` is how coverage is expanded. Do not randomize
the seed at runtime — CI needs reproducibility.

## CI integration (SCENARIOS §5.2)

Two CI entry points drive this fuzz corpus:

1. **PR/push gate** (`.github/workflows/ci.yml` — `paren-fuzz` job): runs
   on every PR with `PAREN_FUZZ_SIZE=1000` and `PAREN_FUZZ_DEFER=1`.
   New mismatches write to `pg/parser/testdata/paren-fuzz-defer/
   <timestamp>.txt` and the job passes; this keeps the fast PR pipeline
   green while still surfacing regressions for triage.
2. **Nightly** (`.github/workflows/container-tests.yml` — `paren-fuzz-nightly`
   job): runs once a day with `PAREN_FUZZ_SIZE=10000` and the same
   `PAREN_FUZZ_DEFER=1` setting. Wider N catches low-frequency drift that
   N=1000 misses.

Triage loop: periodically inspect `pg/parser/testdata/paren-fuzz-defer/`,
move each entry into either `known-mismatches.txt` (legit divergence,
tracked by allowlist) or fix the underlying parser bug, then delete the
timestamped file. Running the fuzz test locally WITHOUT
`PAREN_FUZZ_DEFER` reverts to strict set-equality against
`known-mismatches.txt` and fails on any drift — that's the gate
developers use when promoting a fix.

Environment variables consumed by `TestParenOracleFuzz`:

| Var | Effect |
|-----|--------|
| `PAREN_FUZZ_SIZE` | Overrides the default corpus size (100). Integer > 0. |
| `PAREN_FUZZ_DEFER` | Any non-empty value switches from strict gate to log+persist mode. |
