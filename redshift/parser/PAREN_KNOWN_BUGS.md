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

## PAREN-KB-3 — `LATERAL ()` accepted with empty parens

- **Summary:** omni's `parseLateralTableRef` routes `LATERAL ()` to
  `select_with_parens` and accepts it, producing a `RangeSubselect`
  with a nil (or empty-SELECT) body. PG 17 rejects with 42601
  ("syntax error at or near `)`") because `select_with_parens`
  requires a `select_no_parens` body and empty parens are not one.
- **omni current behavior:** parses successfully; classifier reports
  `OmniSubquery`. Downstream analyze would fail on the degenerate
  subquery body, but raw-parse silently accepts the token sequence.
- **PG correct behavior:** 42601 syntax_error at the closing `)`.
- **Discovery context:** `SCENARIOS-pg-paren-dispatch.md` §3.2
  ("parseLateralTableRef oracle corpus"), scenario 4 (invalid LATERAL
  shapes). Surfaced by the §3.2 test
  `TestParenOracleLateral/invalid_shapes_rejected/LATERAL_empty_parens`.
- **Recommended fix location:** `pg/parser/select.go` →
  `parseLateralTableRef` / the inner `parseSelectWithParens` path.
  After the opening `(` is consumed in the `LATERAL (...)` arm, if the
  next token is `)` (empty body), emit 42601 at the close paren. The
  non-LATERAL path (`SELECT * FROM ()`) is already correctly rejected
  — see `TestParenOracleHarness/empty_parens_rejected` — so the fix
  likely sits on the LATERAL dispatch side that skips the shared
  empty-body check.
- **Priority:** low. Degenerate input unlikely to appear in user SQL,
  but the accept-vs-reject drift is real and the oracle fence should
  eventually be able to drop the skip.
- **Test pins:**
  - `pg/parser/paren_oracle_lateral_test.go`
    `TestParenOracleLateral/invalid_shapes_rejected/LATERAL_empty_parens`
    (skipped with a pointer to PAREN-KB-3).

## PAREN-KB-2 — `Parse()` accepts statements without `;` separator — **CLOSED 2026-04-22**

**Status:** All 13 upstream blockers fixed (KB-2a/b/c/d), then `parser.go:Parse` `needSeparator` check re-applied. pgregress + oracle suites fully green. The KB-2 entry below documents the history of the attempt, the blockers surfaced, and the fix commits that unblocked it.

Closed by commit chain:
- 0f5b7d2 KB-2a: parseRuleActionStmt NotifyStmt case
- 6075da0 KB-2b: ALTER SEQUENCE SET LOGGED/UNLOGGED
- 1808371 KB-2c: CREATE SCHEMA inline schema_element list
- d7e9ba3 KB-2d: CREATE TEMP VIEW inside SCHEMA + DROP FUNCTION empty name
- <this commit> KB-2 reland: needSeparator check in parser.go Parse loop

## PAREN-KB-2 — `Parse()` accepts statements without `;` separator — **history**

**Status:** fix attempted 2026-04-22 and reverted. Enforcing "cur must be `;` or EOF after each parseStmt" at parser.go:Parse surfaced 13 new pgregress failures — all pre-existing omni parser gaps that had been **masked** by the silent-accept behavior (CREATE RULE `DO INSTEAD NOTIFY x` body, SET inside transaction / savepoint blocks, CREATE-chained DDL). Before the statement-list strictness can land, those 13 upstream grammar gaps need to be fixed so the corresponding single-statement parses emit the right AST instead of truncating and letting a second parseStmt pick up the tail.

### A1 KB-2 full-surface survey results (2026-04-22 re-run)

Captured via clean baseline + in-place patch apply + clean revert, across all omni test surfaces:

- `go test ./pg/... -count=1`: 13 NEW FAILURE entries (all in pg/pgregress).
- `go test -tags=oracle ./pg/parser/... -count=1`: 0 additional NEW FAILURE entries (full 188-probe oracle corpus still green under the patch).
- Other packages (pg, pg/ast, pg/catalog, pg/completion, pg/parsertest, pg/plpgsql/parser): all still green.

**Total unique blockers: 13.** No new classes beyond what the delta reveals — oracle fence ran clean even under the stricter Parse loop.

Concrete blocker list, grouped by underlying parser gap (earlier draft mis-attributed some files; this is the accurate breakdown from the 2026-04-22 re-run):

**Category KB-2a — CREATE RULE + NotifyStmt action (4 failures, fix site: `publication.go:689 parseRuleActionStmt`)**
- `copydml.sql:69` stmt[48] — `create rule qqq as on insert to copydml_test do instead notify copydml_test`
- `rules.sql:1017` stmt[474] — `create rule r4 as on delete to rules_src do notify rules_src_deletion`
- `with.sql:1721` stmt[291] — `CREATE OR REPLACE RULE y_rule AS ON INSERT TO y DO INSTEAD NOTIFY foo`
- `with.sql:1726` stmt[293] — `CREATE OR REPLACE RULE y_rule AS ON INSERT TO y DO ALSO NOTIFY foo`

Root cause: `parseRuleActionStmt` enumerates 4 of 5 PG `RuleActionStmt` alternatives (SELECT/INSERT/UPDATE/DELETE) — NotifyStmt falls through `default: return nil, nil`. The 5th production is literally named in the function's own doc comment but not handled.

**Category KB-2b — ALTER SEQUENCE SET LOGGED/UNLOGGED (4 failures, fix site: ALTER SEQUENCE action parser)**
- `identity.sql:535` stmt[266] — `ALTER SEQUENCE identity_dump_logged_a_seq SET UNLOGGED`
- `identity.sql:537` stmt[268] — `ALTER SEQUENCE identity_dump_unlogged_a_seq SET LOGGED`
- `sequence.sql:277` stmt[154] — `ALTER SEQUENCE sequence_test_unlogged SET LOGGED`
- `sequence.sql:279` stmt[155] — `ALTER SEQUENCE sequence_test_unlogged SET UNLOGGED`

Root cause: PG 15+ added `SET LOGGED` / `SET UNLOGGED` to the ALTER SEQUENCE action list. omni's ALTER SEQUENCE action dispatcher stops at an earlier match, leaves `SET` / `UNLOGGED` / `LOGGED` as residual tail.

**Category KB-2c — CREATE SCHEMA schema_element inline body (3 failures, fix site: `parseCreateSchemaStmt`)**
- `create_schema.sql:20` stmt[5] — `CREATE SCHEMA AUTHORIZATION regress_create_schema_role` + following inline `CREATE TABLE ...` etc.
- `create_schema.sql:33` stmt[11] — `CREATE SCHEMA AUTHORIZATION CURRENT_ROLE` + schema_element list
- `create_schema.sql:45` stmt[16] — `CREATE SCHEMA regress_schema_1 AUTHORIZATION CURRENT_ROLE` + schema_element list

Root cause: PG CREATE SCHEMA grammar supports `CREATE SCHEMA [name] [AUTHORIZATION role] [schema_element [...]]` where schema_element is CREATE TABLE / CREATE VIEW / CREATE INDEX / etc. inline. omni's parser stops after AUTHORIZATION, leaves subsequent CREATE tokens as "next statement".

**Category KB-2d — fringe / TBD (2 failures, need spot-check):**
- `errors.sql:139` stmt[38] — trailing "CREATE" after some preceding stmt in the DROP AGGREGATE test context. SQL shown by regress harness is `-- should fail` (comment) but the real content is the preceding DROP AGGREGATE stmt whose body isn't fully consumed.
- `create_view.sql:162` stmt[36] — trailing `(` inside the "-- subqueries" CREATE VIEW test block. Likely a CREATE VIEW parser gap when the view body has subqueries.

Root cause: unknown without reading the extracted SQL precisely; spot-check at fix time.

**Recommended path:**
1. Open follow-up issues for each of the 3 pattern classes (CREATE RULE action body, transaction-block SET, CREATE TRIGGER body).
2. Fix each so the single statement consumes all its tokens.
3. Then re-apply the KB-2 fix — `needSeparator` flag in parser.go's Parse loop requiring `;` or EOF between parseStmt calls. The fix itself is ~5 lines and already designed in commit history (see reverted diff at pg/parser/parser.go Parse loop on 2026-04-22).

Keep the oracle signal:
- `pg/parser/testdata/paren-fuzz-corpus/seed-cases.txt` — NOT pinned (this is a parser.go-Parse bug, not a `(` dispatch bug).
- `pg/parser/testdata/paren-fuzz-corpus/known-mismatches.txt` — NOT pinned; fuzz generator's `genReservedMisuse` no longer emits trailing-SELECT patterns.

## PAREN-KB-2 — `Parse()` silently drops trailing statements (original entry, kept for history)

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
