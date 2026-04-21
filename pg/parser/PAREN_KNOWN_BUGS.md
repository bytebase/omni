# PAREN Known Bugs

Tracked bugs discovered while building the `pg-paren-dispatch` oracle
fence (Phase 2). Each entry documents a real omni-vs-PG-17 divergence
that the oracle harness surfaces but that is scoped OUT of Phase 2 —
the dispatch decision itself is correctly routed; the bug sits
elsewhere. Entries here are the forward references that
`paren_oracle_from_degenerate_test.go` (t.Skip messages) and
`testdata/paren-fuzz-corpus/{seed-cases,known-mismatches}.txt` point at
so the signal isn't lost when future work reads those files.

---

## PAREN-KB-1 — `(T JOIN U)` accepted without ON/USING

- **Summary:** omni's `parseJoinedTable` accepts an inner JOIN with no
  `ON` or `USING` qualifier; PG 17 rejects with 42601
  ("syntax error at or near `)`").
- **omni current behavior:** returns a `*nodes.JoinExpr` with a nil
  qual — the FROM-clause AST is well-formed downstream but the input
  is not valid PG SQL.
- **PG correct behavior:** rejects with syntax_error 42601 at the
  closing `)` because `JOIN` (and `INNER JOIN`, `LEFT JOIN`,
  `RIGHT JOIN`, `FULL JOIN`) require `join_qual`. Only `CROSS JOIN`
  and the `NATURAL` family elide `join_qual`.
- **Discovery context:** `SCENARIOS-pg-paren-dispatch.md` §2.7
  ("FROM-clause degenerate / malformed"), §2.8 (fuzz). Reproduced at
  three depths by the fuzz corpus:
  - `SELECT * FROM (T JOIN U)`
  - `SELECT * FROM (foo LEFT JOIN W)`
  - `SELECT * FROM (foo RIGHT OUTER JOIN W)`
- **Recommended fix location:** `pg/parser/select.go` →
  `parseJoinedTable` / `parseJoinQual`. After the right-hand operand
  is consumed, if the join type is not in {`CROSS`, `NATURAL *`},
  require a `join_qual` (`ON <expr>` or `USING (<collist>)`) and emit
  a syntax error at the next token otherwise.
- **Priority:** medium. Accepting invalid SQL is a real leniency bug;
  low business impact because downstream analyze/plan would fail on
  the nil-qual JoinExpr, but the error message then points at a
  different site than PG's, making user-facing parity worse.
- **Test pins:**
  - `pg/parser/paren_oracle_from_degenerate_test.go` "inner JOIN
    missing qual" (skipped with a pointer to PAREN-KB-1).
  - `pg/parser/testdata/paren-fuzz-corpus/seed-cases.txt` under the
    PAREN-KB-1 section.
  - `pg/parser/testdata/paren-fuzz-corpus/known-mismatches.txt`
    entries for all three depths.

## PAREN-KB-2 — `Parse()` silently drops trailing statements

- **Summary:** omni's top-level `Parse(sql string)` accepts inputs
  like `FROM (SELECT 1) SELECT 1` by returning only the first
  `RawStmt` and ignoring the rest. PG rejects with 42601 at the
  second `SELECT` because `select_with_parens` can only be followed
  by the outer-query's terminator — a second bare `SELECT` is not a
  valid continuation.
- **omni current behavior:** `Parse` returns a `*nodes.List` with a
  single `RawStmt` (the first valid parse); the trailing tokens are
  discarded without surfacing an error.
- **PG correct behavior:** 42601 "syntax error at or near `SELECT`"
  at the start of the trailing statement.
- **Discovery context:** `SCENARIOS-pg-paren-dispatch.md` §2.7, §2.8.
  Originally surfaced in the fuzz mismatches list; on review the
  divergence was recognized as NOT a `(` dispatch bug — the first
  statement is correctly routed to `select_with_parens`. The bug sits
  at the statement-list boundary.
- **Recommended fix location:** `pg/parser/parser.go` `Parse` (the
  outer `for p.cur.Type != 0 { ... }` loop). After each statement is
  parsed, the next token must be `;` or EOF; any other leading token
  (especially `SELECT`/`VALUES`/`TABLE`/`WITH`) for a single-statement
  caller is invalid. Callers that want multi-statement behavior should
  use an explicit `;`. Alternatively, confirm PG's behavior: PG's
  `raw_parser` treats the whole input as a statement list with `;`
  separators; a bare trailing SELECT without `;` should be rejected.
- **Priority:** medium. The silent truncation is a correctness hazard:
  any user who submits a malformed multi-statement payload gets a
  partial result with no error.
- **Test pins:**
  - Previously pinned in
    `pg/parser/testdata/paren-fuzz-corpus/seed-cases.txt` — entries
    removed once the bug was re-attributed (scope: statement-list,
    not `(` dispatch). A removal rationale remains in the seed file
    so the context is preserved.
  - `pg/parser/testdata/paren-fuzz-corpus/known-mismatches.txt` —
    intentionally NOT listed; the fuzz generator no longer emits the
    trailing-SELECT pattern.
