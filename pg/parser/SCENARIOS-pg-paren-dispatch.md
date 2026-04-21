# pg-paren-dispatch Scenarios

> Goal: align omni's pg/parser `(` / `)` dispatch sites with PostgreSQL 17 grammar, and build a regression fence so future parser changes can't silently re-introduce `(`-ambiguity bugs.
> Verification: per-section handwritten unit tests + `pgregress` delta + (Phase 2+) PG 17 testcontainer oracle + (Phase 4) fuzz corpus.
> Reference sources:
> - PG grammar: `../postgres/src/backend/parser/gram.y` (tagged PG 17)
> - Plan: `docs/plans/2026-04-21-pg-paren-dispatch.md` (Codex-reviewed, 2 rounds)
> - Audit: `pg/parser/PAREN_AUDIT.md` + `PAREN_AUDIT.json`
> - Phase 0 tests: `pg/parser/paren_multi_join_test.go`
> - Oracle pattern precedent: `docs/plans/2026-04-14-pg-first-sets.md` §oracle design

Status: [ ] pending, [x] passing, [~] partial (reason)
Reserved suffixes: `(codex-deferred: ...)`, `(codex-override: ...)` — driver-only

---

## Phase 0 (complete) — BYT-9315 core + initial `(` routing

Already landed on branch `junyi/byt-9315-fix-postgresql-grammar-to-accept-multi-join-syntax-in-sql`; recorded here for audit completeness.

- [x] `((a JOIN b ON TRUE) JOIN c ON TRUE)` parses to left-nested JoinExpr matching PG AST
- [x] `(a JOIN b ON TRUE)` single-level parenthesized joined_table parses
- [x] `(((a JOIN b ON TRUE) JOIN c ON TRUE) JOIN d ON TRUE)` triple-nested joined_table parses
- [x] `((SELECT 1))` nested subquery parses as RangeSubselect
- [x] `(((SELECT 1)))` three-level nested subquery parses
- [x] `((SELECT 1) UNION (SELECT 2))` parenthesized set-op with parenthesized operands parses
- [x] `((SELECT 1) UNION ALL (SELECT 2))` set-op with ALL
- [x] `((SELECT 1) INTERSECT (SELECT 2))` INTERSECT variant
- [x] `((SELECT 1) EXCEPT (SELECT 2))` EXCEPT variant
- [x] `(((SELECT 1) UNION (SELECT 2)) UNION (SELECT 3))` nested set-ops
- [x] `(a)` single relation in parens rejected (PG rejects; content not joined_table)
- [x] `((a))` double-wrapped single relation rejected
- [x] `((SELECT 1) x)` aliased subquery in extra parens rejected (content not joined_table)
- [x] `((SELECT 1) x JOIN (SELECT 2) y ON TRUE)` subquery-operand joined_table accepted
- [x] `((SELECT 1) x CROSS JOIN b)` subquery CROSS JOIN relation accepted
- [x] `(a JOIN b ON TRUE) AS jt` aliased parenthesized joined_table accepted
- [x] `(a JOIN b ON TRUE) jt(x, y)` joined_table with column-list alias accepted
- [x] `pg_get_viewdef()` realistic shape with double parens around ON-expr parses
- [x] pgregress: 17 previously-failing PG official regress statements now pass (0 new failures)

---

## Phase 1: Close the 3 remaining misaligned `(` sites + 1 unclear site

Direct fixes for the only ambiguity-present sites the audit found as `aligned != yes`. Each section handles one site, applies the technique from PAREN_AUDIT, lands tests.

### 1.1 parseInExpr — `IN (expr_list)` vs `IN (subquery)`

PG grammar (gram.y:14973-14998): `in_expr` accepts either a `select_with_parens` or a parenthesized `expr_list`. omni currently uses a 1-token peek which can't distinguish `IN (1, 2)` from `IN (SELECT 1)` reliably — it depends on FIRST-set detection of subquery starts.

- [ ] `WHERE x IN (1, 2, 3)` — expr_list path, literals
- [ ] `WHERE x IN (a, b, c)` — expr_list path, identifiers
- [ ] `WHERE x IN ('a', 'b')` — expr_list path, strings
- [ ] `WHERE x IN (1)` — single-element expr_list (not a subquery)
- [ ] `WHERE x IN (SELECT y FROM t)` — simple subquery
- [ ] `WHERE x IN (SELECT 1)` — scalar subquery
- [ ] `WHERE x IN (SELECT 1 UNION SELECT 2)` — set-op subquery
- [ ] `WHERE x IN ((SELECT 1) UNION (SELECT 2))` — parenthesized set-op operands
- [ ] `WHERE x IN (VALUES (1), (2))` — VALUES subquery
- [ ] `WHERE x IN (WITH cte AS (SELECT 1) SELECT * FROM cte)` — WITH subquery
- [ ] `WHERE x IN (TABLE foo)` — TABLE subquery
- [ ] `WHERE x IN ()` — empty list rejected (syntax error, matches PG)
- [ ] `WHERE (x, y) IN (SELECT a, b FROM t)` — row constructor on LHS with subquery
- [ ] `WHERE (x, y) IN ((1,2), (3,4))` — row constructor on LHS with row list
- [ ] `WHERE (x, y) IN (SELECT a FROM t)` — row-constructor LHS with arity-mismatch subquery (parses, rejected later)
- [ ] `WHERE x NOT IN (1, 2, 3)` — NOT IN expr_list
- [ ] `WHERE x NOT IN (SELECT y FROM t)` — NOT IN subquery
- [ ] `WHERE x NOT IN (SELECT 1 UNION SELECT 2)` — NOT IN set-op subquery
- [ ] `WHERE (x, y) NOT IN (SELECT a, b FROM t)` — NOT IN row-constructor
- [ ] AST parity: `IN (SELECT 1)` produces SubLink with subLinkType=ANY_SUBLINK (matches PG)
- [ ] AST parity: `IN (1, 2)` produces ScalarArrayOpExpr with operator `=` (matches PG)
- [ ] AST parity: `NOT IN (SELECT 1)` produces SubLink with subLinkType=ANY_SUBLINK wrapped in NOT
- [ ] AST parity: `NOT IN (1, 2)` produces ScalarArrayOpExpr with operator `<>`, useOr=false
- [ ] Out-of-scope: `x = ANY (SELECT ...)` / `x = SOME (...)` / `x = ALL (...)` (gram.y:14976-14998 `sub_type` family — tracked separately, not a `(` dispatch issue)
- [ ] pgregress: file `subselect.sql` IN-subquery cases no longer fail (delta report)

### 1.2 parseLateralTableRef — LATERAL dispatch

PG grammar (gram.y:13611-13620): LATERAL can prefix `func_table`, `select_with_parens`, `xmltable`, `json_table`. omni's current path goes straight to `(` → select_with_parens without distinguishing XMLTABLE / JSON_TABLE.

- [ ] `SELECT * FROM t, LATERAL (SELECT 1) x` — LATERAL select_with_parens accepted
- [ ] `SELECT * FROM t, LATERAL (SELECT * FROM u WHERE u.x = t.x) y` — LATERAL with outer ref
- [ ] `SELECT * FROM t, LATERAL XMLTABLE('/root' PASSING x COLUMNS a int PATH 'a') as xt` — LATERAL xmltable accepted
- [ ] `SELECT * FROM t, LATERAL JSON_TABLE(x, '$' COLUMNS(a int PATH '$.a')) as jt` — LATERAL json_table accepted
- [ ] `SELECT * FROM t, LATERAL f(t.x)` — LATERAL func_table accepted
- [ ] `SELECT * FROM t, LATERAL f(t.x) WITH ORDINALITY` — LATERAL func_table + ordinality
- [ ] `SELECT * FROM t, LATERAL f(t.x) AS fa(a int, b text)` — LATERAL func_table + column-definition alias
- [ ] `SELECT * FROM t, LATERAL ROWS FROM (f(t.x), g(t.y))` — LATERAL rows-from accepted
- [ ] `SELECT * FROM t, LATERAL ROWS FROM (f(t.x), g(t.y)) WITH ORDINALITY` — LATERAL rows-from + ordinality
- [ ] `SELECT * FROM t, LATERAL (a JOIN b ON TRUE)` — LATERAL joined_table REJECTED (PG rejects; grammar has no such production)
- [ ] `SELECT * FROM t, LATERAL u` — bare LATERAL relation REJECTED (LATERAL prefix requires subquery/xmltable/json_table/func_table; plain relation_expr is not a variant)
- [ ] AST: LATERAL subquery sets `RangeSubselect.Lateral=true`
- [ ] AST: LATERAL func_table sets `RangeFunction.Lateral=true`
- [ ] AST: LATERAL xmltable sets `RangeTableFunc.Lateral=true`
- [ ] AST: LATERAL json_table sets `JsonTable.Lateral=true`

### 1.3 parseArraySubscript — `ARRAY[...]` vs `ARRAY(SELECT ...)`

PG grammar (gram.y:14845-14910): two distinct productions — `ARRAY '[' ... ']'` is array constructor, `ARRAY '(' select_no_parens ')'` is ARRAY_SUBLINK. omni's 1-token peek at `(` currently routes all `ARRAY(...)` to sublink; `ARRAY[...]` handling may be broken or underexercised.

- [ ] `SELECT ARRAY[1, 2, 3]` — array constructor, int literals
- [ ] `SELECT ARRAY['a', 'b']` — array constructor, strings
- [ ] `SELECT ARRAY[]::int[]` — empty array with cast accepted
- [ ] `SELECT ARRAY[]` — empty array without cast REJECTED (PG: "cannot determine type of empty array")
- [ ] `SELECT ARRAY[[1,2],[3,4]]` — nested array constructor
- [ ] `SELECT ARRAY[1, 2] || ARRAY[3]` — array concatenation in expr context
- [ ] `SELECT ARRAY(SELECT 1)` — ARRAY sublink scalar
- [ ] `SELECT ARRAY(SELECT x FROM t)` — ARRAY sublink column
- [ ] `SELECT ARRAY(SELECT DISTINCT x FROM t ORDER BY x)` — ARRAY sublink with DISTINCT+ORDER BY
- [ ] `SELECT ARRAY(VALUES (1), (2))` — ARRAY VALUES subquery
- [ ] AST parity: `ARRAY[1,2]` produces A_ArrayExpr with elements, multidims=false (matches PG)
- [ ] AST parity: `ARRAY[[1,2],[3,4]]` produces A_ArrayExpr with multidims=true
- [ ] AST parity: `ARRAY(SELECT 1)` produces SubLink with subLinkType=ARRAY_SUBLINK

### 1.4 parseArraySubscript second site — contract lock

`expr.go:1610` was marked `unclear` in audit. Per PG 17 gram.y:14845-14850, `ARRAY '(' select_no_parens ')'` restricts the paren content to `select_no_parens` — therefore `TABLE foo` and `WITH cte AS ... SELECT ...` ARE valid (both are select_no_parens shapes); `ROWS FROM (...)` is NOT (it's only reachable via func_table in FROM, not in select_no_parens). Scenarios commit to that expected outcome.

- [ ] `SELECT ARRAY(TABLE foo)` — ACCEPT (TABLE is a simple_select form → select_no_parens)
- [ ] `SELECT ARRAY(WITH cte AS (SELECT 1) SELECT * FROM cte)` — ACCEPT (WITH + select_clause is select_no_parens)
- [ ] `SELECT ARRAY(VALUES (1), (2))` — ACCEPT (VALUES is simple_select)
- [ ] `SELECT ARRAY(ROWS FROM (f(1)))` — REJECT (ROWS FROM is func_table-only, not valid in ARRAY sublink)
- [ ] `SELECT ARRAY(1)` — REJECT (bare expr not a select_no_parens)
- [ ] `SELECT ARRAY()` — REJECT (empty parens not a select_no_parens)
- [ ] After fix/verify, PAREN_AUDIT.json row for expr.go:1610 flips `aligned: unclear` → `aligned: yes` with the 6 scenarios above as empirical proof

---

## Phase 2: PG-oracle harness for parenBeginsSubquery

Build a PG 17 testcontainer-based oracle that probes every conceivable `FROM (...)` shape and compares omni's routing decision (subquery vs joined_table vs reject) against PG's accept/reject. This is the durable regression fence the plan §5.2 calls for.

### 2.1 Oracle infrastructure

- [ ] PG 17 testcontainer boots and accepts connections from the test harness
- [ ] Oracle driver sends a probe SQL to PG, receives accept (row parsed) or reject (SQLSTATE)
- [ ] Oracle driver tokenizes omni's result: accepted-as-subquery, accepted-as-joined-table, accepted-as-other, rejected-syntax
- [ ] Probe result stored with `{sql, pg_status, omni_status, pg_error?, omni_error?}`
- [ ] Mismatch detector: reports when `pg_status != omni_status` with side-by-side diff
- [ ] Test skip when testcontainer unavailable (CI-local); never fails the suite for infra reasons
- [ ] Oracle harness runs under `go test -tags=oracle` or equivalent, not in default test pass
- [ ] Harness output includes per-probe timing to measure against plan §6.5 CI-cost budget (~2 min for 200 probes)

### 2.2 FROM-clause simple-shape corpus

One probe per canonical shape. `T`, `U`, `V` are pre-created tables on the oracle container; `t1(a)` is a func_table.

- [ ] `SELECT * FROM (T)` — PG: reject (not joined_table)
- [ ] `SELECT * FROM ((T))` — PG: reject
- [ ] `SELECT * FROM (T JOIN U ON TRUE)` — PG: accept, JoinExpr
- [ ] `SELECT * FROM (T CROSS JOIN U)` — PG: accept
- [ ] `SELECT * FROM (T LEFT JOIN U ON TRUE)` — accept
- [ ] `SELECT * FROM (T RIGHT JOIN U ON TRUE)` — accept
- [ ] `SELECT * FROM (T FULL JOIN U ON TRUE)` — accept
- [ ] `SELECT * FROM (T INNER JOIN U ON TRUE)` — accept
- [ ] `SELECT * FROM (T NATURAL JOIN U)` — accept
- [ ] `SELECT * FROM (T NATURAL LEFT JOIN U)` — accept
- [ ] `SELECT * FROM (T NATURAL FULL OUTER JOIN U)` — accept
- [ ] `SELECT * FROM (T JOIN U USING (a))` — accept
- [ ] `SELECT * FROM (T JOIN U USING (a) AS alias_clause)` — accept, join-using alias
- [ ] `SELECT * FROM (T LEFT OUTER JOIN U ON T.a = U.a)` — accept
- [ ] `SELECT * FROM (T)` as LATERAL prefix `LATERAL (T)` — PG: reject
- [ ] `SELECT * FROM (T AS alias)` — PG: reject (content is aliased relation, not joined_table)
- [ ] `SELECT * FROM (T, U)` — PG: reject (comma in paren is not joined_table)
- [ ] `SELECT * FROM (ONLY T)` — PG: reject
- [ ] `SELECT * FROM (T TABLESAMPLE BERNOULLI(10))` — PG: reject (tablesample in paren)
- [ ] `SELECT * FROM (T WITH ORDINALITY)` — PG: reject (WITH ORDINALITY only valid on func_table)
- [ ] `SELECT * FROM (T FOR UPDATE)` — PG: reject (locking clause not valid inside FROM-paren)

### 2.3 FROM-clause subquery-shape corpus

- [ ] `SELECT * FROM (SELECT 1)` — accept subquery
- [ ] `SELECT * FROM (SELECT 1) AS s` — accept with alias
- [ ] `SELECT * FROM (SELECT 1) s(x)` — accept with column alias
- [ ] `SELECT * FROM ((SELECT 1))` — accept double-wrapped
- [ ] `SELECT * FROM (((SELECT 1)))` — accept triple-wrapped
- [ ] `SELECT * FROM ((((SELECT 1))))` — accept four-wrapped
- [ ] `SELECT * FROM (VALUES (1))` — PG: reject (needs alias for VALUES-in-paren? verify)
- [ ] `SELECT * FROM (VALUES (1)) AS v(a)` — accept
- [ ] `SELECT * FROM ((VALUES (1)) AS v(a))` — double-wrap with alias
- [ ] `SELECT * FROM (WITH cte AS (SELECT 1) SELECT * FROM cte)` — accept
- [ ] `SELECT * FROM ((WITH cte AS (SELECT 1) SELECT * FROM cte))` — double-wrap
- [ ] `SELECT * FROM (TABLE T)` — accept, TABLE subquery
- [ ] `SELECT * FROM (SELECT 1 UNION SELECT 2)` — accept, set-op
- [ ] `SELECT * FROM (SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3)` — multi-operand set-op
- [ ] `SELECT * FROM (SELECT 1 INTERSECT SELECT 2)` — INTERSECT
- [ ] `SELECT * FROM (SELECT 1 EXCEPT SELECT 2)` — EXCEPT
- [ ] `SELECT * FROM ((SELECT 1) UNION (SELECT 2))` — parenthesized operands
- [ ] `SELECT * FROM (((SELECT 1) UNION (SELECT 2)) UNION (SELECT 3))` — nested set-op with parens

### 2.4 FROM-clause joined_table-shape corpus

- [ ] `SELECT * FROM ((T JOIN U ON TRUE))` — double-wrapped joined_table
- [ ] `SELECT * FROM (((T JOIN U ON TRUE)))` — triple-wrapped
- [ ] `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE)` — BYT-9315 shape
- [ ] `SELECT * FROM (T JOIN (U JOIN V ON TRUE) ON TRUE)` — right-side paren-join
- [ ] `SELECT * FROM ((T JOIN U ON TRUE) JOIN (V JOIN W ON TRUE) ON TRUE)` — both sides paren-join
- [ ] `SELECT * FROM (((T JOIN U ON TRUE) JOIN V ON TRUE) JOIN W ON TRUE)` — deeply nested left
- [ ] `SELECT * FROM (T JOIN (U JOIN (V JOIN W ON TRUE) ON TRUE) ON TRUE)` — deeply nested right
- [ ] `SELECT * FROM (T CROSS JOIN U JOIN V ON TRUE)` — mixed join types
- [ ] `SELECT * FROM (T LEFT JOIN U ON T.a = U.a RIGHT JOIN V ON U.b = V.b)` — chained outer joins
- [ ] `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt` — outer alias on paren-joined
- [ ] `SELECT * FROM ((T JOIN U ON TRUE) JOIN V ON TRUE) AS jt(col1)` — outer column-list alias
- [ ] `SELECT * FROM (T JOIN U ON TRUE) AS jt JOIN V ON TRUE` — alias then outer join
- [ ] `SELECT * FROM (T NATURAL JOIN U) JOIN V USING (a)` — natural then using
- [ ] `SELECT * FROM (T FULL OUTER JOIN U USING (a))` — FULL OUTER + USING (gram.y `join_type` branch)
- [ ] `SELECT * FROM (T LEFT OUTER JOIN U USING (a, b))` — multi-column USING
- [ ] `SELECT * FROM ((T JOIN U ON ((T.a = U.a))) JOIN V ON ((U.b = V.b)))` — double-paren ON-expr (actual `pg_get_viewdef()` output shape)
- [ ] `SELECT * FROM (T JOIN LATERAL (SELECT U.a FROM U WHERE U.x = T.x) v ON TRUE)` — LATERAL subquery as JOIN right operand inside paren-joined_table
- [ ] `SELECT * FROM ((T JOIN U ON T.a = U.a) JOIN LATERAL (SELECT 1) l ON TRUE)` — LATERAL as second join's right operand

### 2.5 FROM-clause mixed shapes

`((subquery) JOIN relation)` — the exact shape from Phase 0 fix.

- [ ] `SELECT * FROM ((SELECT 1) x JOIN U ON TRUE)` — subquery as join operand
- [ ] `SELECT * FROM ((SELECT 1) x CROSS JOIN U)` — subquery CROSS JOIN
- [ ] `SELECT * FROM ((SELECT 1) x NATURAL JOIN U)` — subquery NATURAL JOIN
- [ ] `SELECT * FROM ((SELECT 1) x JOIN (SELECT 2) y ON x.a = y.a)` — both operands subquery
- [ ] `SELECT * FROM ((VALUES (1)) v(a) JOIN U ON TRUE)` — VALUES as operand
- [ ] `SELECT * FROM ((TABLE T) JOIN U ON TRUE)` — TABLE as operand (PG: check — may reject due to TABLE's form)
- [ ] `SELECT * FROM ((T) JOIN U ON TRUE)` — PG: reject (paren-single-relation not joined_table operand... verify)
- [ ] `SELECT * FROM (T JOIN (SELECT 1) s(x) ON TRUE)` — subquery as right operand

### 2.6 FROM-clause LATERAL interactions

- [ ] `SELECT * FROM T, LATERAL (SELECT 1)` — accept
- [ ] `SELECT * FROM T, LATERAL (SELECT 1) x` — accept with alias
- [ ] `SELECT * FROM T, LATERAL ((SELECT 1))` — accept double-wrap
- [ ] `SELECT * FROM T, LATERAL ((SELECT 1) UNION (SELECT 2))` — LATERAL with set-op parens
- [ ] `SELECT * FROM T, LATERAL (a JOIN b ON t.x = a.x)` — PG: reject (LATERAL joined_table not a production)
- [ ] `SELECT * FROM T, LATERAL XMLTABLE('/root' PASSING T.doc COLUMNS a int PATH 'a')` — LATERAL xmltable
- [ ] `SELECT * FROM T, LATERAL JSON_TABLE(T.doc, '$' COLUMNS(a int PATH '$.a'))` — LATERAL json_table

### 2.7 FROM-clause degenerate / malformed

- [ ] `SELECT * FROM ()` — PG: reject (empty parens)
- [ ] `SELECT * FROM (` — PG: reject (unclosed paren)
- [ ] `SELECT * FROM )` — reject
- [ ] `SELECT * FROM (SELECT 1` — PG: reject (unclosed after subquery)
- [ ] `SELECT * FROM (T JOIN` — reject (JOIN without right operand)
- [ ] `SELECT * FROM (T JOIN U)` — reject (missing join qual for inner join — PG requires ON or USING)
- [ ] `SELECT * FROM (T CROSS JOIN U ON TRUE)` — reject (CROSS JOIN has no qual)
- [ ] `SELECT * FROM (T NATURAL JOIN U ON TRUE)` — reject (NATURAL has no qual)
- [ ] `SELECT * FROM (SELECT)` — reject (incomplete SELECT)
- [ ] `SELECT * FROM (SELECT 1))` — reject (extra close paren)
- [ ] `SELECT * FROM ((SELECT 1)` — reject (missing close paren)
- [ ] `SELECT * FROM (SELECT 1 FROM)` — reject (FROM without relation)
- [ ] `SELECT * FROM ( ) SELECT 1` — reject (leading empty parens + stray SELECT)

### 2.8 Fuzz-generated paren combinations

Property-based corpus comparing omni vs PG 17 accept/reject on randomly generated balanced-paren FROM-clauses. Depth bound, seed count, and corpus storage are Stage 2 design decisions; this section defines coverage intent.

- [ ] Fuzz corpus exercises nested balanced parens
- [ ] Fuzz corpus exercises SELECT/VALUES/WITH/TABLE at varying depths
- [ ] Fuzz corpus exercises UNION/INTERSECT/EXCEPT between operands
- [ ] Fuzz corpus exercises JOIN/CROSS JOIN/NATURAL JOIN at varying positions
- [ ] Fuzz corpus exercises LATERAL prefixes, aliases, column-lists
- [ ] Fuzz corpus exercises obvious-reject cases (unbalanced, empty, reserved-word misuse)
- [ ] Fuzz mismatch rate between omni and PG stays below an agreed threshold
- [ ] Fuzz mismatches persist to a golden file for human triage, not silent test failure

---

## Phase 3: Oracle extension to the 3 newly-fixed sites (after Phase 1 lands)

Phase 2 hardens the `FROM (...)` primitive. Phase 3 extends the oracle discipline to the 3 Phase 1 fix sites so the fixes themselves can't silently regress.

### 3.1 parseInExpr oracle corpus

- [ ] Oracle compares `WHERE x IN (1,2,3)` omni AST shape vs PG ScalarArrayOpExpr
- [ ] Oracle compares `WHERE x IN (SELECT 1)` omni SubLink vs PG SubLink
- [ ] Oracle probes list size variants: single element, 2-element, 10-element, 100-element (4 sizes)
- [ ] Oracle probes literal kind variants: int, float, string, bool, null (5 kinds)
- [ ] Oracle probes subquery kind variants: SELECT, VALUES, WITH ... SELECT, TABLE, SELECT UNION SELECT (5 kinds)
- [ ] Oracle probes row-constructor LHS `(x,y) IN (...)` against PG (both list and subquery RHS)
- [ ] Oracle probes `NOT IN` for expr_list path
- [ ] Oracle probes `NOT IN` for subquery path
- [ ] Oracle probes `IN` in JOIN ON: `FROM T JOIN U ON T.a IN (SELECT ...)`
- [ ] Oracle probes `IN` in HAVING: `HAVING count(*) IN (1, 2, 3)`
- [ ] Oracle probes `IN` in CASE WHEN: `CASE WHEN x IN (1,2) THEN 'y' ELSE 'n' END`
- [ ] Oracle probes `IN` in CHECK constraint: `CREATE TABLE t (x int CHECK (x IN (1,2,3)))`
- [ ] Oracle probes `IN` inside a subquery's WHERE: `SELECT (SELECT 1 WHERE x IN (SELECT y FROM t))`

### 3.2 parseLateralTableRef oracle corpus

- [ ] Oracle compares LATERAL (SELECT) vs LATERAL xmltable vs LATERAL json_table AST shapes
- [ ] Oracle probes LATERAL + column-list alias combinations
- [ ] Oracle probes LATERAL with outer-table reference (typical correlated use)
- [ ] Oracle probes invalid LATERAL shapes (LATERAL joined_table, LATERAL ROWS FROM without parens, etc.) and confirms omni rejects matching PG

### 3.3 parseArraySubscript oracle corpus

- [ ] Oracle compares ARRAY[...] A_ArrayExpr vs ARRAY(...) SubLink shapes
- [ ] Oracle probes nested ARRAY[ARRAY[...]] constructions
- [ ] Oracle probes ARRAY with type cast combinations
- [ ] Oracle probes ARRAY sublink with VALUES/TABLE/WITH variants
- [ ] Oracle probes negative cases — ARRAY() empty, ARRAY[SELECT], etc.

---

## Phase 4: §5.3 lock-in proofs for the remaining aligned ambiguity-present sites

The audit found 2 ambiguity-present sites already aligned in Phase 0 (`parseParenTableRef`, `parseSelectClausePrimary`). These have handwritten tests but no formal `§5.3 two-proof bar` entry.

### 4.1 Phase 0 sites audit-row proofs

- [ ] `parseParenTableRef` caller-context proof written into PAREN_AUDIT.md
- [ ] `parseParenTableRef` empirical pinned tests ≥ 5 (already have 18 in paren_multi_join_test.go — cite file:line)
- [ ] `parseSelectClausePrimary` caller-context proof written
- [ ] `parseSelectClausePrimary` empirical pinned tests ≥ 5
- [ ] `parseJoinedTable` caller-context proof written
- [ ] `parseJoinedTable` empirical pinned tests ≥ 5
- [ ] `parseSelectWithParens` (the post-Phase 0 left-factored form) caller-context proof written
- [ ] All 4 proof rows mirrored in PAREN_AUDIT.json

### 4.2 §5.3 aligned-without-code-change audit rows for non-ambiguous sites

The 85+ non-ambiguous sites ("expect `(` after keyword" / optional paren-list) don't need the formal two-proof treatment — grammar structure is itself the proof. But the audit should cite **which existing test exercises each one** for anti-regression traceability.

- [ ] Each of the 9 C3 (type parens) sites cites an existing test covering it
- [ ] Each of the 14 C4 (DDL element list) sites cites an existing test covering it
- [ ] Each of the 33 C5 sites cites an existing test covering it
- [ ] Uncovered sites (if any) get minimal handwritten tests added — target < 10 new tests
- [ ] PAREN_AUDIT.json `proof_notes` field populated with test citations for every row

---

## Phase 5: CI integration and long-term fence

### 5.1 Oracle harness in CI

- [ ] Oracle harness wired to GitHub Actions via build tag `oracle`
- [ ] Oracle runs on PR against pg/parser/ changes (file-scoped trigger)
- [ ] Oracle timing ≤ 5 min in CI (green signal for plan §6.5 budget)
- [ ] Oracle mismatches fail the CI check with side-by-side SQL / omni-AST / PG-accept diff
- [ ] Oracle baseline file committed — known-diff entries tracked like pgregress known_failures.json

### 5.2 Fuzz corpus in CI

- [ ] Fuzz corpus stored under `pg/parser/testdata/paren-fuzz-corpus/`
- [ ] Fuzz run sampled at N=1000 in CI, full N=10000 nightly
- [ ] New mismatches from fuzz auto-file to `testdata/paren-fuzz-defer/` for triage (not auto-fail)
- [ ] `go test -tags=fuzz ./pg/parser/` runs the fuzz seeds + corpus as regular tests

### 5.3 PAREN_AUDIT governance

- [ ] `PAREN_AUDIT.json` schema documented in plan §6.1 (machine-readable)
- [ ] CI lint: every `aligned=yes` row must have non-empty `proof_notes` pointing to either a test file:line or a sister-starmap dependency
- [ ] CI lint: every `p.cur.Type == '('` / `p.cur.Type == ')'` in pg/parser/*.go must have a matching row in PAREN_AUDIT.json (drift detection)
- [ ] New dispatch site added without updating PAREN_AUDIT.json fails CI
