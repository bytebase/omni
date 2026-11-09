# refaudit

The TiDB syntax-mapping ledger (design doc §3.6). It checks that every
`// Ref: ... parser.y:N` citation in omni's TiDB productions names a real line in
the vendored v8.5.5 `parser.y`, and reports which grammar rules no production
cites.

Like `grammarcorpus`, it is a plain root-module command that reads text and
writes text, importing nothing from omni or pingcap. `go test ./...` compiles it;
the companion test runs the same audit.

## Run

```bash
go run ./tidb/parser/cmd/refaudit            # rewrite testdata/parser_y_coverage.md
go run ./tidb/parser/cmd/refaudit -check     # validate citations, write nothing
```

A broken citation — a line number past the grammar section, or one that lands in
the gap between rules — exits non-zero. `TestRefCitationsResolve` gates the same
check with no server; `TestCoverageReport` keeps the report current (`-update`).

## What it is and isn't

A citation says "meant to implement," not "accepts exactly": omni is hand-written
to its own contract, not generated from `parser.y` (§2), so a matching rule can
still differ in precedence, FIRST-sets, backtracking, and lexing. The ledger
guides and audits coverage; it does not gate on it. The gate on over-acceptance
is the differential rejection suite (the negatives corpus and the fuzz sweep).

Coverage is coarse: one production often cites the line range of a whole rule
group, and many `parser.y` rules are option/suffix sub-productions omni handles
inline without their own citation. A rule in the uncited list is a lead to check,
not a proven gap.
