# Trino Migration Analysis

Migration target: full coverage of the legacy ANTLR4 `bytebase/parser/trino` grammar
by a hand-written recursive-descent parser under `omni/trino`, adjudicated for
correctness against a **live Trino 481 oracle** (Docker). Bytebase consumption is
used to **prioritize** the order of work; it is never used to **scope down**
legacy parser features.

Sources analyzed:

- **Legacy parser (truth2):** `/Users/h3n4l/OpenSource/parser/trino` — `TrinoParser.g4`
  (1,132 lines, 128 rules, 81 statement-rule labeled alternatives), `TrinoLexer.g4`
  (397 lines, ~220 keyword tokens, no lexer modes). Apache-2.0, derived from Trino's
  official `SqlBase.g4`. Catalog: [`antlr_rules.md`](antlr_rules.md).
- **Documented syntax (truth1):** Trino 481 official docs (`trino.io/docs/current`),
  ~172 syntax forms across [`truth1/ddl.md`](truth1/ddl.md) (38),
  [`truth1/query-dml.md`](truth1/query-dml.md) (40),
  [`truth1/types-expr-functions.md`](truth1/types-expr-functions.md) (51),
  [`truth1/admin-utility.md`](truth1/admin-utility.md) (43).
- **Oracle:** live Trino 481 via Docker, Tier 2. [`oracle.md`](oracle.md) +
  harness at `trino/internal/trinooracle`.
- **Bytebase consumption (contract):** `backend/plugin/parser/trino`,
  `backend/plugin/schema/trino`, `backend/plugin/db/trino`. Full surface in
  [`contract.md`](contract.md).

---

## 1. Scope and Significance

Trino (formerly PrestoSQL) is a distributed SQL query engine. The grammar is a
**full, standards-aligned SQL surface** — substantially larger than partiql and,
in query/expression depth, richer than doris:

- **81 statement forms**: full DDL (TABLE/SCHEMA/VIEW/MATERIALIZED VIEW/CATALOG/
  ROLE/FUNCTION), DML (INSERT/DELETE/UPDATE/MERGE/TRUNCATE), DCL (GRANT/REVOKE/
  DENY/roles), admin/session (USE/SET SESSION/SET ROLE/SET PATH/SET TIME ZONE),
  transactions, prepared statements (PREPARE/EXECUTE/DEALLOCATE), SHOW family,
  EXPLAIN, CALL, ANALYZE.
- **Deep query layer**: CTEs (incl. `WITH RECURSIVE`), set operations, all JOIN
  types, `GROUP BY` with `GROUPING SETS`/`CUBE`/`ROLLUP`, named `WINDOW` + full
  frame syntax (`ROWS`/`RANGE`/`GROUPS`), `ORDER BY`/`OFFSET`/`LIMIT`/`FETCH … WITH TIES`.
- **Large Trino-specific features** (no analogue in other omni engines):
  - **`MATCH_RECOGNIZE`** row-pattern recognition (PATTERN/DEFINE/MEASURES/SUBSET/
    AFTER MATCH SKIP/row-pattern quantifiers & anchors) — a grammar subsystem on its own.
  - **Polymorphic table functions** `TABLE(func(...) [PARTITION BY] [COPARTITION] ...)`.
  - **SQL/JSON**: `JSON_EXISTS/JSON_QUERY/JSON_VALUE/JSON_OBJECT/JSON_ARRAY` with
    JSON-path invocation, `PASSING`, and behavior clauses.
  - **Inline SQL routines**: `CREATE FUNCTION` / `WITH FUNCTION` with a full
    control-flow body language (`BEGIN/END`, `IF`, `CASE`, `LOOP`, `WHILE`, `REPEAT`,
    `ITERATE`, `LEAVE`, `DECLARE`).
  - **Time travel**: `FOR (TIMESTAMP|VERSION) AS OF` (`queryPeriod`).
  - Lambda expressions (`x -> expr`), `ARRAY[...]`/`ROW(...)`/subscript/dereference,
    `LISTAGG ... WITHIN GROUP`, `AT TIME ZONE`, `UNNEST [WITH ORDINALITY]`, `LATERAL`,
    `TABLESAMPLE`.

Because bytebase's `Diagnose` (P0) runs on **every** statement and must not emit
false syntax errors, the omni parser must accept the **entire** legacy grammar
before the import switch — grammar breadth is itself P0. Priority below governs
*ordering*, not scope.

---

## 2. Bytebase Contract (the must-work surface)

Engine: `storepb.Engine_TRINO = 24`. **No** separate Presto engine. The DB driver
(`plugin/db/trino`) does **not** call the parser directly; all dispatch is via the
`base.*` registry. Full detail in [`contract.md`](contract.md).

| Registered API | Handler | Consumed by (runtime) | P0 |
|---|---|---|---|
| `RegisterSplitterFunc` | `SplitSQL` | every query execute, export, LSP | **YES** |
| `RegisterQueryValidator` | `validateQuery` | SQL-editor read-only guard | **YES** |
| `RegisterGetQuerySpan` | `GetQuerySpan` | **masking** + lineage + export ACL | **YES** |
| `RegisterCompleteFunc` | `Completion` | LSP autocomplete (`EngineSupportAutoComplete=true`) | **YES** |
| `RegisterDiagnoseFunc` | `Diagnose` | LSP inline diagnostics (every keystroke) | **YES** |
| `schema.RegisterGetDatabaseDefinition` | `GetDatabaseDefinition` | schema dump / SDL | **YES** |

**Not registered** (not gaps): `ParseStatementsFunc` (Trino uses internal
`ParseTrino`), `StatementRangesFunc`, `QueryTypeFunc` (internal helper),
`ChangedResourcesGetter`, `TransformDMLToSelect`, `GenerateRestoreSQL`, masking
advisors, SchemaDiff.

**Consumption depth notes:**
- `GetQuerySpan` is the heaviest consumer: extracts `SourceColumns` (base-table
  column lineage for masking), `PredicateColumns` (WHERE/JOIN-ON/USING columns;
  computed but enforcement gated to MSSQL only today), and per-result-column spans
  with deferred `SELECT *` expansion. Walks SELECT/CTE/JOIN/UNNEST/LATERAL and
  recurses into subqueries. Requires a real SELECT + expression AST with table/
  column resolution against `DatabaseMetadata`.
- `Completion` is a ~1730-line c3 (`CodeCompletionCore`) implementation; preferred
  rules `identifier` + `qualifiedName`; produces catalog/schema/table/column/keyword
  candidates; re-parses FROM clauses and uses `GetQuerySpan` to resolve CTE/subquery
  column names. Depends on the parser's rule indices and a catalog.
- `validateQuery` and `query_type` classify SELECT / EXPLAIN / `SelectInfoSchema`
  (system/information_schema) / DML / DDL.
- `GetDatabaseDefinition` emits `CREATE TABLE` SDL from `DatabaseSchemaMetadata` —
  metadata→DDL, largely independent of the parser.

---

## 3. Legacy Parse API (truth2)

```go
parser.NewTrinoLexer(antlr.CharStream) *TrinoLexer
parser.NewTrinoParser(antlr.TokenStream) *TrinoParser
(*TrinoParser).SingleStatement() ISingleStatementContext   // primary entry
(*TrinoParser).Parse() IParseContext                        // statements* EOF
```

Bytebase wrappers (`plugin/parser/trino`): `ParseTrino(sql) ([]*base.ANTLRAST, error)`
(internal), `SplitSQL`, `validateQuery`, `GetQuerySpan`, `Completion`, `Diagnose`,
plus exported helpers `NormalizeTrinoIdentifier`, `ExtractQualifiedNameParts`,
`ExtractDatabaseSchemaName`. AST = ANTLR-generated context graph (~128 rule
contexts + labeled-alternative contexts). The omni AST will be hand-written and
need only preserve the semantic distinctions consumers care about.

---

## 4. omni Architecture Fit

omni engines are hand-written recursive-descent parsers (stateless lexer producing
`Token{Type,Str,Loc,End}`, parser with 2-token lookahead, AST nodes carrying
`Loc{Start,End}` byte offsets). Reference engines: `pg/` (most mature ast/parser
split), `mysql/` (DDL-heavy statement dispatch + c3 completion + live-DB oracle in
`mysql/parser/oracle_test.go`). Target package layout:

```
trino/
  ast/                 # AST node types, Loc tracking, walker
  parser/              # lexer + recursive-descent parser + statement splitter + Parse() entry
    testdata/
      legacy/          # lifted from bytebase/parser/trino/examples (94 .sql files)
      official/        # from Trino 481 docs (truth1)
  analysis/            # query span: lineage (SourceColumns) + predicate columns
  completion/          # c3-style completion (preferred rules: identifier, qualifiedName)
  internal/trinooracle/  # ✅ BUILT — differential oracle (Trino REST client + self-test)
```

`GetDatabaseDefinition` (metadata→SDL) currently lives in bytebase
`plugin/schema/trino`; planning to decide whether it moves to `trino/deparse` or
stays in bytebase (it does not depend on the omni parser).

---

## 5. Divergences (docs/oracle vs legacy grammar) — for the ledger

The legacy ANTLR grammar (truth2) **lags Trino 481** (truth1 + oracle). Oracle-
adjudicated findings to seed the divergence ledger:

| # | Case | Legacy grammar | Trino 481 (oracle) | Disposition | Confidence |
|---|---|---|---|---|---|
| 1 | `NEAREST JOIN` | absent | **REJECT-SYNTAX** | docs-crawl **error**; discard `join-nearest` from truth1 | high |
| 2 | `COLON_` token literal `'_:'` | `'_:'` (bug) | `:` | implement `:`; legacy literal is a copy/paste bug | high |
| 3 | `JSON_TABLE` | **absent** from `primaryExpression`/`relationPrimary` | exists (481) | docs-ahead-of-legacy; P1/P2 extension — planning decides | high |
| 4 | `WITH SESSION prop = v SELECT …` | absent (`with` = CTE only) | ACCEPT (sem) | docs-ahead-of-legacy; P1 extension | high |
| 5 | DML branch `INTO tbl@branch` | absent | ACCEPT (sem `BRANCH_NOT_FOUND`) | docs-ahead-of-legacy; P1 extension | high |
| 6 | set-op `UNION CORRESPONDING` | `setQuantifier` = DISTINCT/ALL only | ACCEPT | docs-ahead-of-legacy; P1 extension | high |
| 7 | `GROUP BY ALL` / `GROUP BY AUTO` | `groupBy` has `setQuantifier` (ALL) | both ACCEPT (AUTO clean, ALL semantic) | verify exact semantics vs identifier-named-`all`/`auto` | medium |
| 8 | `SKIP_`, `TRY_CAST_`, explicit `DOT_`/`COMMA_` | local additions vs official `SqlBase.g4` | n/a | lexer-shape divergences; reconcile token model during lexer node | medium |

Target policy (per doris/partiql convention): implement the **legacy grammar
scope** correctly (oracle-adjudicated), and treat docs-ahead-of-legacy items
(#3–#6) as explicit P1 extensions the planner can include since the oracle proves
them real in 481. Discard #1.

---

## 6. Priority Ranking (ordering within full-parity scope)

P0 = required before bytebase can switch its imports (all 81 statement forms must
parse without error for `Diagnose`/`SplitSQL`/`validateQuery`). P1 = legacy-parity
/ docs-ahead items not blocking the switch.

### Tier 0 — Foundation (blocks everything) · P0
1. **Lexer** — ~220 case-insensitive keywords, reserved vs `nonReserved` split,
   identifiers (unquoted/`"quoted"`/`` `backtick` ``/digit-leading), `U&'...'` unicode
   strings, `X'...'` binary, numeric (int/decimal/double), operators, `||`, `=>`/`->`,
   `{-`/`-}`, comments (hidden channel), `;`. Position metadata for splitter.
2. **AST core** — node interfaces, `Loc` tracking, walker.
3. **Statement splitter** — token-driven (`;`); drives `base.SplitMultiSQL`.
4. **Parser entry + Diagnose** — `Parse`/`ParseTrino`; error collection with positions.

### Tier 1 — Identifiers, types, expressions · P0
5. **`qualifiedName` / `identifier`** + normalization; `nonReserved` keyword-as-identifier.
6. **Type grammar** — `genericType`, `arrayType`/`ARRAY<>`, `MAP<>`, `ROW(...)`,
   `INTERVAL ... TO ...`, `TIMESTAMP/TIME [(p)] [WITH[OUT] TIME ZONE]`, `DOUBLE PRECISION`,
   `DECIMAL(p,s)`, plus Trino types (JSON, IPADDRESS, UUID, HLL, *digest, VARIANT).
7. **Expression layer** — precedence-climbing `booleanExpression`→`valueExpression`→
   `primaryExpression`: predicates (`comparison`/`quantifiedComparison`/`BETWEEN`/`IN`
   list+subquery/`LIKE`/`IS [NOT] NULL`/`IS [NOT] DISTINCT FROM`), arithmetic/concat/
   `AT TIME ZONE`, and the big `primaryExpression` set (literals, `CASE`, `CAST`/`TRY_CAST`,
   function calls with `FILTER`/`ORDER BY`/`nullTreatment`/`OVER`, `LISTAGG`, lambda,
   `ARRAY[]`/`ROW()`/subscript/dereference, `EXISTS`/subquery, `EXTRACT`/`SUBSTRING`/
   `TRIM`/`NORMALIZE`/`POSITION`, special datetime fns, `GROUPING`, **JSON functions**).

### Tier 2 — SELECT core · P0
8. **`querySpecification`** — SELECT list (`selectSingle`/`selectAll`), `FROM`,
   `relation`/`joinType`/`joinCriteria` (CROSS/INNER/LEFT/RIGHT/FULL/NATURAL), `WHERE`,
   `GROUP BY` (+ `GROUPING SETS`/`CUBE`/`ROLLUP`), `HAVING`, `WINDOW` + frames.
9. **`queryNoWith` wrapper** — `ORDER BY`/`OFFSET`/`LIMIT`/`FETCH FIRST … ONLY|WITH TIES`.
10. **Set operations** (`UNION`/`INTERSECT`/`EXCEPT` [ALL|DISTINCT]).
11. **CTEs** (`with`/`namedQuery`, `WITH RECURSIVE`).
12. **Relation primaries** — subquery, `UNNEST [WITH ORDINALITY]`, `LATERAL`,
    `TABLE(tableFunctionCall …)`, `TABLESAMPLE`, `queryPeriod` time travel.
13. **`MATCH_RECOGNIZE`** — full row-pattern subsystem (own node).

### Tier 3 — Query span analysis · P0 (heaviest consumer)
14. **`analysis` package** — `GetQuerySpan` parity: table/column lineage,
    `SELECT *` expansion, CTE tracking, set-op column merge, subquery recursion,
    JOIN/UNNEST/LATERAL source resolution, predicate-column extraction.

### Tier 4 — DML + classification · P0
15. **DML** — `INSERT`, `DELETE`, `UPDATE`, `MERGE` (`mergeCase`), `TRUNCATE`,
    `tableExecute`, `refreshMaterializedView`.
16. **Statement classification** — `query_type` parity (Select/Explain/
    SelectInfoSchema/DML/DDL) + `validateQuery` read-only guard.

### Tier 5 — DDL · P0
17. CREATE/DROP/ALTER **TABLE** (incl. `CREATE TABLE AS`, `columnDefinition`,
    `likeClause`, properties), **SCHEMA**, **VIEW**, **MATERIALIZED VIEW** (`GRACE PERIOD`),
    **CATALOG**, **COMMENT ON**, **ANALYZE**.
18. **Functions** — `CREATE [OR REPLACE] FUNCTION` / `DROP FUNCTION` with the full
    `controlStatement` routine body; `WITH FUNCTION` inline.

### Tier 6 — Admin / session / txn / prepared / DCL · P0
19. `USE`, `SET`/`RESET SESSION`, `SET SESSION AUTHORIZATION`, `SET ROLE`, `SET PATH`,
    `SET TIME ZONE`; `SHOW` family; `DESCRIBE`/`DESC`; `EXPLAIN`/`EXPLAIN ANALYZE`;
    `PREPARE`/`EXECUTE`/`EXECUTE IMMEDIATE`/`DEALLOCATE`/`DESCRIBE INPUT|OUTPUT`;
    `START TRANSACTION`/`COMMIT`/`ROLLBACK`; `CREATE/DROP ROLE`, `GRANT`/`REVOKE`/`DENY`;
    `CALL`.

### Tier 7 — Completion · P0
20. **`completion` package** — c3 over the omni parser (preferred rules identifier,
    qualifiedName), catalog/schema/table/column/keyword candidates, FROM-scope and
    CTE/subquery resolution.

### Tier 8 — Schema definition · P0
21. **`GetDatabaseDefinition`** parity (metadata→`CREATE TABLE` SDL). Parser-independent.

### P1 — docs-ahead-of-legacy extensions (oracle-proven real in 481)
22. `JSON_TABLE`; `WITH SESSION`; DML `@branch`; `UNION/INTERSECT/EXCEPT CORRESPONDING`.
    (Discard `NEAREST JOIN` — oracle rejects it.)

---

## 7. Test Corpus Strategy

- **Corpus A — legacy examples (regression baseline):** lift the 95 `.sql` files
  from `/Users/h3n4l/OpenSource/parser/trino/examples/` into `trino/parser/testdata/legacy/`.
- **Corpus B — official 481 docs (truth1):** every `example_sql` captured in the
  four `truth1/*.md` files, cross-referenced to the grammar rule it exercises.
- **Differential gate:** each grammar node runs its corpus through
  `trino/internal/trinooracle` (`assertOracleMatch`: omni `Parse` accept/reject ==
  Trino 481 accept/reject). The oracle self-test (`TestOracleClassification`,
  15 cases) already passes against Trino 481; it is gated/skipped without a server.

---

## 8. Notes for Planning

- **Grammar breadth is P0.** `Diagnose` runs on every statement, so all 81 forms
  must parse cleanly before the import switch. Sequence by dependency (foundation →
  expr → select → analysis → statements), but don't drop forms.
- **Three oversized subsystems** deserve their own DAG nodes: `MATCH_RECOGNIZE`,
  SQL/JSON functions+path, and the `CREATE FUNCTION` routine body language.
- **`GetQuerySpan` (Tier 3) is the critical-path consumer** for masking — it needs
  the full SELECT+expression AST and column resolution. It gates the masking feature.
- **Completion depends on the parser's rule-index model** (c3). The omni parser must
  expose rule indices / a completion-core integration analogous to mysql/tsql.
- **Live oracle available** (Trino 481, Docker) — every node should adjudicate
  accept/reject mechanically rather than by judgment. This is the highest-leverage
  correctness asset; prefer it over antlr-faithfulness where the two conflict
  (see §5 divergences).
- **`GetDatabaseDefinition`** is parser-independent; can be done in parallel and may
  remain in bytebase.
