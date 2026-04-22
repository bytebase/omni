# MSSQL Parser Syntax Gaps — Systematic Fix Plan

Date: 2026-04-21
Branch: `junyi/mssql-parser-optimize`
Goal: Close known T-SQL syntax gaps in omni's mssql parser and put an AST-level oracle fence in place so they don't reopen.

## Context

Migration-side audit surfaced 4 real gaps:
1. `SELECT a = c FROM t` — T-SQL `column_alias = expression` form not recognized; parser eats `=` as binary-eq.
2. `(SELECT ...) AS x(c1, c2)` — derived table column alias list not parsed.
3. `fn(...) AS x(v)` — TVF column alias list not parsed (affects `CROSS APPLY`/`OUTER APPLY`).
4. `FROM (VALUES ...)` — VALUES not accepted as a derived table source.

## Systematic Root-Cause Analysis

### Pattern A — `AS alias [ (col_list) ]` tail implemented inconsistently

Every derived-table-like source should end with the same optional tail. Today it is implemented locally twice (CTE, OUTPUT INTO) and missing/partial elsewhere:

| Source | alias | (cols) | Location |
|---|---|---|---|
| CTE | ✅ | ✅ | `select.go:494-515` (local impl) |
| OUTPUT INTO | ✅ | ✅ | `insert.go:244-262` (local impl) |
| Subquery `(SELECT…)` | ✅ | ❌ | `select.go:786` (**Gap 2**) |
| TVF `fn(…)` | ✅ | ❌ | `select.go:1085` (**Gap 3**) |
| `(VALUES …)` | ❌ | ❌ | **Gap 4**, not wired into FROM |
| `OPENJSON WITH(…)` | ✅ | ⚠️ parsed but discarded | `rowset_functions.go:56` |

`AliasedTableRef.Columns` (`parsenodes.go:2011`) already exists and is wired through walker + deparser. The parser just never populates it.

### Pattern B — `alias = expr` SELECT list item (Gap 1)

T-SQL–specific syntax for SELECT list; equivalent to `expr AS alias`. SqlScriptDOM models via `SelectScalarExpression.ColumnName` with 2-token lookahead (`IDENT '='`). omni calls `parseExpr` first, so `a = c` is consumed as `BinaryExpr{Eq}` before the alias detector runs.

### Not systematic risk (excluded)

- PIVOT/UNPIVOT `AS alias` does **not** take `(col_list)` in T-SQL (SqlScriptDOM's `PivotedTableReference.Alias` is a single Identifier). Don't add it.

## Reference: SqlScriptDOM Models

Source at `/Users/rebeliceyang/Github/SqlScriptDOM` (sibling checkout).

| Gap | SqlScriptDOM node | Disambiguation |
|---|---|---|
| 1 | `SelectScalarExpression { ColumnName, Expression }` | 2-token lookahead `IDENT '='` at list-item head |
| 2 | `QueryDerivedTable { QueryExpression, Alias, Columns }` | after `AS alias`, peek `(` |
| 3 | `SchemaObjectFunctionTableReference { …, Alias, Columns }` | same |
| 4 | `InlineDerivedTable { RowValues, Alias, Columns }` | dispatch on `VALUES` inside FROM `(` branch |

## Plan

### Phase 0 — Audit

Deliverable: `SCENARIOS-mssql-syntax-gaps.md` with numbered scenarios covering:
- all T-SQL SELECT list item forms (vs. SqlScriptDOM `SelectElement` subclasses)
- all T-SQL `TableReference` subclasses vs. omni parse dispatch

Each row marked ✅ / ⚠️ / ❌ against omni. Known 4 gaps + anticipated audit finds (CHANGETABLE, SEMANTICSIMILARITYTABLE, FOR SYSTEM_TIME, OPENJSON WITH discard, …).

Estimate scope: if > 20 scenarios, escalate to starmap (`mssql-syntax-gaps-driver/worker`). Else drive from this plan + SCENARIOS doc.

### Phase 1 — Foundation refactor (single commit, behavior unchanged)

1. New helper `parser/select.go::parseAliasAndOptionalColumnList() (alias string, cols *nodes.List)`
2. Replace the bare `parseOptionalAlias()` at subquery (L786) and TVF (L1085) with the helper — for now the helper still discards cols (matches current behavior exactly). Columns wiring lit up in Phase 2.
3. `mssql/ast/parsenodes.go`: `ValuesClause` gains `tableExpr() {}`.
4. Walker / deparser / outfuncs: add `*ValuesClause` as TableExpr branch.
5. Full existing test suite stays green.

### Phase 2 — Gap fixes, one commit per batch

| Batch | Scope | Corpus test |
|---|---|---|
| B1 | Gap 2 subquery col list (turn on cols in helper for subquery branch) | S-2* |
| B2 | Gap 3 TVF col list | S-3* |
| B3 | Gap 4 `(VALUES …) AS x(cols)` in FROM | S-4* |
| B4 | Gap 1 `alias = expr` in SELECT list | S-1* |
| B5 | OPENJSON WITH column list no longer discarded | S-5* |
| B6+ | Any audit findings | TBD |

Each batch: parser change + corpus entries + any needed deparser/walker adjustments. `go test ./mssql/... -short` must be green before moving on.

### Phase 3 — SqlScriptDOM AST diff harness

Layout:
```
harness/mssql-scriptdom/          # .NET console app
  Program.cs                       # stdin SQL → normalized JSON AST
  mssql-scriptdom.csproj           # references SqlScriptDom via ProjectReference or dll
mssql/parser/scriptdom_harness_test.go   # Go side: drives dotnet, compares AST shape
```

- SqlScriptDOM side emits compact JSON: `{kind, alias, columns[], children[], expr_kind, op, …}` — only fields relevant for the gaps we care about (not a full AST serialization, which would be churn).
- omni side emits the same shape from its AST via a small visitor in test helper.
- Diff = go-cmp on the JSON trees.
- Fixture file: one SQL per scenario, gated behind `-tags scriptdom` (optional build tag so CI without dotnet still passes base tests).

### Phase 4 — Verify + decide PR strategy

- `go build ./mssql/... && go test ./mssql/... -count=1 -short`
- Run oracle testcontainer suite
- Run ScriptDOM harness over full scenario fixture
- If scope grew: create starmap skills
- With user: one PR vs. split commits into several PRs (user said "commit first, decide after")

## Key decisions locked

1. **Oracle**: Use SqlScriptDOM (sibling checkout) for AST-level diff. SQL Server testcontainer (PARSEONLY) remains for accept/reject.
2. **Starmap**: Default no. Escalate if Phase 0 finds > 20 scenarios.
3. **PRs**: Single branch with per-batch commits. PR strategy decided after Phase 4.
