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
