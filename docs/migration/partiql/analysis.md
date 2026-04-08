# PartiQL Migration Analysis

PartiQL is the SQL++-flavored query language used by AWS DynamoDB and Azure Cosmos DB.
In bytebase it is wired to `storepb.Engine_DYNAMODB`. The legacy parser is ANTLR4-based
and lives at `bytebase/parser/partiql`; the bytebase wrapper lives at
`bytebase/backend/plugin/parser/partiql`. The omni migration target is **full coverage
of the legacy ANTLR grammar's scope**, with bytebase consumption used only for
prioritization.

## Source Locations

| Path | Role |
|------|------|
| `bytebase/parser/partiql/PartiQLParser.g4` | Parser grammar (~685 lines) |
| `bytebase/parser/partiql/PartiQLLexer.g4`  | Lexer grammar (~543 lines) |
| `bytebase/parser/partiql/partiqlparser_parser.go` | ANTLR-generated parser (~212 context types) |
| `bytebase/parser/partiql/partiqlparser_listener.go` | Generated listener interface (~1107 lines) |
| `bytebase/parser/partiql/simple.sql` | Tiny smoke fixture (DELETE / INSERT) |
| `bytebase/backend/plugin/parser/partiql/partiql.go` | `ParsePartiQL()` entry point |
| `bytebase/backend/plugin/parser/partiql/split.go`   | Statement splitter |
| `bytebase/backend/plugin/parser/partiql/query.go`   | DQL-only query validator |
| `bytebase/backend/plugin/parser/partiql/completion.go` | Auto-completion |
| `bytebase/backend/plugin/db/dynamodb/dynamodb.go` | Driver — only consumer of the wrapper outside `init` registration |

## Grammar Coverage

### Top-level structure

`script` → `statement (; statement)*`. `statement` is one of: `dql`, `dml`, `ddl`, `execCommand`.
`root` allows an optional `EXPLAIN` prefix.

### DDL

- `CREATE TABLE <name>` — name only, no column definitions, no constraints
- `CREATE INDEX ON <table> ( <path>, <path>, ... )`
- `DROP TABLE <name>`
- `DROP INDEX <index_name> ON <table>`

No `ALTER`, no `VIEW`, no constraints (PK/FK/CHECK/NOT NULL/DEFAULT). This is
intentional — PartiQL is schemaless and document-collection oriented.

### DML

PartiQL has two coexisting DML syntaxes:

- **Legacy path-based** (g4 lines 94–100): `UPDATE … SET …`, `FROM … DELETE …`,
  `INSERT INTO <path> VALUE … AT … ON CONFLICT …`, `REMOVE <path>`
- **RFC 0011 modern**: `INSERT INTO <collection> AS <alias> VALUE … ON CONFLICT …`,
  `REPLACE INTO <collection> AS <alias> VALUE …`,
  `UPSERT INTO <collection> AS <alias> VALUE …`
- `DELETE FROM <path> WHERE …`

`ON CONFLICT` (g4 lines 139–180):

- `ON CONFLICT WHERE <expr> DO NOTHING` (legacy form)
- `ON CONFLICT [( cols ) | ON CONSTRAINT <name>] [DO NOTHING | DO REPLACE … | DO UPDATE …]`
- **TODO in legacy grammar (g4 lines 170, 180):** `doReplace` and `doUpdate` are stubs —
  only the `EXCLUDED` keyword is parsed; full `SET …` / `VALUE …` clauses are missing.

`RETURNING` (g4 lines 194–200): `RETURNING [MODIFIED|ALL] [OLD|NEW] (* | expr) , …`.

### DQL — SELECT

SELECT variants (g4 lines 216–221):

- `SELECT [DISTINCT|ALL] *`
- `SELECT [DISTINCT|ALL] expr [AS alias], …`
- `SELECT [DISTINCT|ALL] VALUE <expr>` — single-value projection (PartiQL-unique)
- `PIVOT <expr> AT <expr>` — used as a SELECT replacement (PartiQL-unique)

Clauses:

- `LET <expr> AS <alias>, …` (g4 238–242) — intermediate variable bindings, PartiQL-unique
- `FROM <path> [AS <itemAlias>] [AT <posAlias>] [BY <keyAlias>]` — `AT` for array index,
  `BY` for object key (PartiQL-unique)
- `UNPIVOT <expr> [AS …] [AT …] [BY …]` (g4 408)
- `WHERE <expr>`
- `GROUP [PARTIAL] BY <expr> [AS alias], … [GROUP AS alias]` — `PARTIAL` is PartiQL-unique
- `HAVING <expr>`
- `ORDER BY <expr> [ASC|DESC] [NULLS FIRST|LAST], …`
- `LIMIT <expr>`, `OFFSET <expr>`

Set operations (g4 449–454): `UNION`, `INTERSECT`, `EXCEPT`, each with optional
`OUTER` and `[DISTINCT|ALL]`.

### Joins

`CROSS`, `INNER`, `LEFT [OUTER]`, `RIGHT [OUTER]`, `FULL [OUTER]`, bare `OUTER` (natural
outer). `ON <expr>` condition. No `LATERAL`, no `USING`.

### Window Functions

Only `LAG(expr [, offset [, default]]) OVER (…)` and `LEAD(…)` are explicitly in the
grammar (g4 589–591). No `ROW_NUMBER`, `RANK`, `DENSE_RANK`, `NTILE`. Aggregate
functions inside `OVER(…)` are not explicit in the grammar.

### Aggregates

`COUNT(*)`, `COUNT([DISTINCT|ALL] expr)`, `MAX`, `MIN`, `SUM`, `AVG`. Nothing else
(no `STDDEV`, no `ARRAY_AGG`, no custom aggregates).

### Expressions

Operator hierarchy: `OR` → `AND` → `NOT` → predicate (`=`/`<>`/`<`/`>`/`<=`/`>=`/`IN`/
`BETWEEN`/`LIKE`/`IS [NOT] (NULL|MISSING|TRUE|FALSE)`) → `||` (concat) → `+/-` → `*///%`
→ unary → primary. Multiple `mathOp` rules give the precedence layers.

`expr` includes: literals, paths, function calls, `CASE [expr] WHEN … THEN … ELSE … END`,
`CAST(expr AS type)`, `CAN_CAST`, `CAN_LOSSLESS_CAST`, parameter `?`, named variables
`@ident` or `ident`, subqueries, graph match expressions.

### Path Navigation (PartiQL-critical)

Path steps (g4 618–623):

- `[<expr>]` — index/key by expression
- `[*]` — all-elements wildcard
- `.<symbol>` — dot field access
- `.*` — all-fields wildcard
- Chains: `expr.a[0].b.c[key]`

### Collection / Tuple Literals

- **List** `[ expr, expr, … ]` — ordered, indexed (SQL array)
- **Bag**  `<< expr, expr, … >>` — unordered, duplicate-preserving (PartiQL-unique)
- **Tuple/Struct** `{ key: value, key: value, … }` (PartiQL-unique)
- `LIST(…)`, `SEXP(…)` constructor function syntax

### Type System (g4 674–686)

Standard SQL types: `NULL`, `BOOL/BOOLEAN`, `SMALLINT`, `INTEGER`/`INT`, `BIGINT`,
`REAL`, `DOUBLE PRECISION`, `DECIMAL(p,s)`, `NUMERIC(p,s)`, `FLOAT(n)`, `VARCHAR(n)`,
`CHAR(n)`/`CHARACTER(n)`, `DATE`, `TIME`, `TIMESTAMP`.

PartiQL-specific types: `MISSING` (distinct from `NULL`), `STRING`, `SYMBOL` (Ion),
`STRUCT`, `TUPLE`, `LIST`, `BAG`, `SEXP`, `BLOB`, `CLOB`, `ANY`.

### Built-in Functions

String: `CHAR_LENGTH`, `CHARACTER_LENGTH`, `OCTET_LENGTH`, `BIT_LENGTH`, `UPPER`,
`LOWER`, `SIZE` (collection size), `EXISTS`, `TRIM([LEADING|TRAILING|BOTH] [sub] FROM target)`,
`SUBSTRING(expr FROM x [FOR y])`.

Type: `CAST(expr AS type)`, `CAN_CAST`, `CAN_LOSSLESS_CAST`, `COALESCE`, `NULLIF`.

Date/time: `EXTRACT(field FROM expr)`, `DATE_ADD(unit, expr, expr)`, `DATE_DIFF(unit, expr, expr)`,
`DATE 'yyyy-mm-dd'`, `TIME [(p)] [WITH TIME ZONE] 'HH:MM:SS'`, `TIMESTAMP '…'` literals.

Generic: any `<ident>(args)` function call (extensibility hook).

### Graph Pattern Matching (GPML)

Significant feature (g4 314–382, 625–630). Examples of constructs:

- Nodes `(label? pattern? where?)`, edges `-[label?]->`, `<-[]-`, `~[…]~`
- Quantifiers `+`, `*`, `{m,n}`
- Selectors `ANY`, `ALL SHORTEST`, `SHORTEST k`
- Restrictors `TRAIL`, `ACYCLIC`, `SIMPLE` (parsed but not semantically enforced)
- Path variables: `p = node-[edge]-node`
- `exprGraphMatchOne` and `exprGraphMatchMany`

### Ion Embedded Literals

Backtick-delimited inline Ion values (g4 669, lexer modes lines 327–403). The lexer
has a dedicated ION mode that recognises Ion's symbols, structs, and binary forms.

### EXEC

`EXEC <name> [arg, arg, …]` — call stored procedures. No `CREATE PROCEDURE`.

## Parse API Surface

### `bytebase/backend/plugin/parser/partiql/partiql.go`

```go
func ParsePartiQL(statement string) ([]*base.ANTLRAST, error)
```

- Splits the input via `SplitSQL`, then parses each statement with the ANTLR4 lexer/parser
- Returns one `*base.ANTLRAST` per statement; each contains:
  - `StartPosition *storepb.Position`
  - `Tree antlr.ParseTree`
  - `Tokens *antlr.CommonTokenStream`
- Custom error listeners: `base.ParseErrorListener` on both lexer and parser
- Registered: `base.RegisterParseStatementsFunc(storepb.Engine_DYNAMODB, parsePartiQLStatements)` (partiql.go:16)

### `split.go`

```go
func SplitSQL(statement string) ([]base.Statement, error)
```

- Splits by `PartiQLLexerCOLON_SEMI` (`;`) tokens via the ANTLR lexer
- Tracks line/column/byte positions per statement
- Registered: `base.RegisterSplitterFunc(storepb.Engine_DYNAMODB, SplitSQL)` (split.go:12)

### `query.go`

```go
func validateQuery(statement string) (bool, bool, error)
```

- Walks the parse tree with a `BasePartiQLParserListener` to assert that every top-level
  statement is DQL (or `EXPLAIN` of DQL)
- Rejects DDL, DML, and `EXEC`
- Registered: `base.RegisterQueryValidator(storepb.Engine_DYNAMODB, validateQuery)` (query.go:12)

### `completion.go`

```go
func Completion(ctx context.Context, cCtx base.CompletionContext, statement string, …) ([]base.Candidate, error)
```

- Uses `base.CodeCompletionCore` (a Go port of the ANTLR4 c3 completion library)
- Preferred parser rules: `varRefExpr`, `selectClause`, `fromClause`, `exprSelect`
- Two completer strategies: standard and "tricky" (for trailing-state corner cases)
- Surfaces table names from FROM scope and column names from current SELECT scope
- Registered: `base.RegisterCompleteFunc(storepb.Engine_DYNAMODB, Completion)` (completion.go:177)

### Error Handling Pattern

All entry points share the bytebase parser convention: install
`base.ParseErrorListener` on the lexer and parser, run the parse, then return
the listener's accumulated error.

## AST Types

The legacy AST is the ANTLR-generated parse-tree context graph (~212 context types).
Significant categories:

- Top-level: `ScriptContext`, `RootContext`, `StatementContext`,
  `QueryDqlContext`, `QueryDmlContext`, `QueryDdlContext`, `QueryExecContext`
- Clauses: `SelectClauseContext`, `FromClauseContext`, `WhereClauseContext`,
  `GroupClauseContext`, `OrderByClauseContext`, `LimitClauseContext`,
  `OffsetClauseContext`, `LetClauseContext`, `HavingClauseContext`
- DDL/DML: `CreateTableContext`, `CreateIndexContext`, `DropTableContext`,
  `DropIndexContext`, `InsertCommandContext`, `UpdateClauseContext`,
  `DeleteCommandContext`, `ReplaceCommandContext`, `UpsertCommandContext`,
  `RemoveCommandContext`
- Expressions: `ExprContext`, `ExprSelectContext`, `ExprOrContext`, `ExprAndContext`,
  `ExprNotContext`, `ExprPredicateContext`, `MathOp00/01/02Context`,
  `ValueExprContext`, `ExprPrimaryContext`
- Path: `PathStepContext`, `PathStepIndexExprContext`, `PathStepDotExprContext`,
  `PathSimpleContext`, `PathSimpleStepsContext`
- Collections: `ArrayContext`, `BagContext`, `TupleContext`, `PairContext`,
  `ValueListContext`, `ValuesContext`
- Graph: `GpmlPatternContext`, `MatchPatternContext`, `NodeContext`, `EdgeContext`,
  `PatternQuantifierContext`, `EdgeSpecContext`, `PatternPartLabelContext`
- Functions: `AggregatContext`, `WindowFunctionContext`, `FunctionCallContext`,
  `FunctionCallIdentContext`, `FunctionCallReservedContext`, `CastContext`,
  `ExtractContext`, `TrimFunctionContext`, `SubstringContext`

The omni AST will be hand-written and need not preserve every ANTLR context — only
the semantic distinctions consumers care about.

## Bytebase Import Sites

Bytebase has only **one direct user** of the wrapper package outside of `init()`-time
registration:

| File | Function | Line | API used |
|------|----------|------|----------|
| `backend/plugin/db/dynamodb/dynamodb.go` | `Execute` | 99 | `base.SplitMultiSQL(storepb.Engine_DYNAMODB, statement)` |
| `backend/plugin/db/dynamodb/dynamodb.go` | `QueryConn` | 131 | `base.SplitMultiSQL(storepb.Engine_DYNAMODB, statement)` |

The DynamoDB driver splits multi-statement input because the AWS PartiQL transaction
API cannot mix read and write in one call — splitting is operationally critical.

The wrapper package is auto-imported (for its `init()` registrations) from:

- `backend/server/ultimate.go:37`
- `backend/component/sheet/sheet.go:17`

`storepb.Engine_DYNAMODB` participates in many capability checks in
`backend/common/engine.go` (lines 38, 77, 115, 142, 174, 228, 266, 302, 328, 366, 420),
but those are flag checks, not parser calls.

## Feature Dependency Map

| Bytebase Feature | Parser API | Data Extracted | Status |
|------------------|------------|----------------|--------|
| DynamoDB driver multi-statement execution | `SplitMultiSQL` (Splitter) | Statement boundaries with positions | **Required** |
| SQL Editor read-only validation | `ValidateSQLForEditor` (QueryValidator) | DQL vs not-DQL boolean | **Required** |
| SQL Editor auto-complete | `Completion` (CompleteFunc) | Table/column/keyword candidates by cursor scope | **Required** |
| Generic statement parsing | `ParseStatements` (ParseStatementsFunc) | Parse tree + tokens | **Required (pipeline plumbing)** |
| Schema sync | — | — | Not implemented |
| Query span / lineage | — | — | Not implemented |
| Masking | — | — | Not implemented |
| Diagnostics / lint advisors | — | — | Not implemented |
| DML → SELECT transformation | — | — | Not implemented |
| Restore SQL generation | — | — | Not implemented |
| Statement ranges (LSP) | — | — | Not implemented |

Compared with mongo/cosmosdb wrappers, partiql is **strictly the smaller subset** —
no analysis, masking, span, or diagnostics layer exists.

## Gaps and Limitations in the Legacy Parser

These are constraints inherited from the ANTLR grammar that omni does not need to
preserve as bugs:

1. **`ON CONFLICT DO REPLACE/UPDATE` are stubs** — only `EXCLUDED` is parsed
   (g4 lines 170, 180). Full SET/VALUE clauses are missing.
2. **No `ALTER TABLE`** — by design (NoSQL).
3. **No views, no triggers, no `CREATE PROCEDURE`** — only `EXEC` invocation.
4. **No constraints in `CREATE TABLE`** — name only, no columns.
5. **No CTEs / `WITH` clause.**
6. **No recursive queries.**
7. **Limited window functions** — only `LAG`/`LEAD`, no `ROW_NUMBER`/`RANK`/`DENSE_RANK`/`NTILE`.
8. **Aggregate window functions ambiguous** — grammar doesn't explicitly allow
   `SUM(...) OVER (...)`, though expression nesting may permit it.
9. **Graph pattern restrictors not enforced** — `TRAIL`/`ACYCLIC`/`SIMPLE` parsed
   but no semantics.
10. **No `LATERAL` join.**
11. **Comments stripped** — lexer routes them to the hidden channel; not preserved
    on the AST.
12. **Tiny smoke fixture** — `simple.sql` is 5 lines covering only DELETE/INSERT.
    Coverage is grammar-driven, not example-driven.

## Priority Ranking

Migration order, derived from bytebase consumption first then legacy parser parity:

### Tier 0 — Foundations (must come first)

1. Lexer: keywords, identifiers, strings, numbers, `;` delimiter, Ion backtick mode
2. Token-stream / position tracking infrastructure
3. Statement splitter (drives `base.SplitMultiSQL`, the one feature the DynamoDB driver
   directly needs)

### Tier P0 — Required by bytebase

4. Recursive-descent parser scaffolding (precedence climbing for expressions)
5. SELECT path: `SELECT … FROM … WHERE … GROUP BY … HAVING … ORDER BY … LIMIT/OFFSET`,
   including `SELECT VALUE`, set operations
6. FROM path with `AS`/`AT`/`BY` aliases, joins
7. Path navigation expressions (`a.b[0].c`, `[*]`, `.*`)
8. Collection/tuple literals (`[…]`, `<<…>>`, `{k:v,…}`)
9. Variable refs (`@ident`, `?` parameter)
10. DML: `INSERT … VALUE …`, `UPDATE … SET …`, `DELETE FROM …`, `UPSERT`, `REPLACE`,
    `REMOVE`, `ON CONFLICT`, `RETURNING`
11. DDL: `CREATE TABLE`, `CREATE INDEX`, `DROP TABLE`, `DROP INDEX`
12. `EXEC` command
13. `EXPLAIN` prefix
14. Query validator (DQL-only enforcement matching `validateQuery`)
15. Completion engine (table-name and column-name candidates with FROM-scope tracking)

### Tier P1 — Legacy parity (not consumed by bytebase, still in scope)

16. `LET` clause
17. `PIVOT` / `UNPIVOT`
18. Window functions (`LAG`/`LEAD`)
19. Aggregates (`COUNT(*)`, `COUNT(DISTINCT …)`, `SUM`, `AVG`, `MIN`, `MAX`)
20. Built-in functions (`CAST`, `EXTRACT`, `TRIM`, `SUBSTRING`, `DATE_ADD`, `DATE_DIFF`,
    `COALESCE`, `NULLIF`, `CASE`)
21. Graph pattern matching (`MATCH` clauses, nodes/edges/quantifiers/selectors)
22. Ion-embedded literals (backtick mode)
23. `CAN_CAST` / `CAN_LOSSLESS_CAST`

### Tier P2 — Nice-to-have / future

24. Schema sync, query span, masking, diagnostics — register stubs only when bytebase
    decides to add them; not required for the import switch

## Full Coverage Target

| Area | Feature | Tier |
|------|---------|------|
| Lexer | All keywords, identifiers, numbers, strings, `;`, Ion backtick mode | **P0** |
| Splitter | Statement splitting by `;` | **P0** |
| Parser | `EXPLAIN` prefix | **P0** |
| Parser DDL | `CREATE TABLE`, `CREATE INDEX (path,…)`, `DROP TABLE`, `DROP INDEX … ON …` | **P0** |
| Parser DML | `INSERT INTO … VALUE …` (legacy + RFC 0011) | **P0** |
| Parser DML | `UPDATE … SET …`, `FROM … SET …` | **P0** |
| Parser DML | `DELETE FROM … WHERE …`, `FROM … DELETE …` | **P0** |
| Parser DML | `UPSERT INTO …`, `REPLACE INTO …` | **P0** |
| Parser DML | `REMOVE <path>` | **P0** |
| Parser DML | `ON CONFLICT … DO NOTHING / DO REPLACE EXCLUDED / DO UPDATE EXCLUDED` | **P0** |
| Parser DML | `RETURNING [MODIFIED\|ALL] [OLD\|NEW] *\|expr` | **P0** |
| Parser DQL | `SELECT [DISTINCT\|ALL] *\|items\|VALUE expr` | **P0** |
| Parser DQL | `PIVOT expr AT expr` projection | **P1** |
| Parser DQL | `FROM path [AS] [AT] [BY]`, `UNPIVOT` | **P0** (FROM) / **P1** (UNPIVOT) |
| Parser DQL | Joins: `CROSS`, `INNER`, `LEFT/RIGHT/FULL [OUTER]`, bare `OUTER` | **P0** |
| Parser DQL | `WHERE`, `GROUP [PARTIAL] BY … [GROUP AS]`, `HAVING` | **P0** |
| Parser DQL | `ORDER BY … [ASC\|DESC] [NULLS FIRST\|LAST]`, `LIMIT`, `OFFSET` | **P0** |
| Parser DQL | `LET expr AS alias, …` | **P1** |
| Parser DQL | Set ops `UNION/INTERSECT/EXCEPT [OUTER] [DISTINCT\|ALL]` | **P0** |
| Parser Expr | Operator precedence: `OR/AND/NOT/predicate/||/+−/*//%/unary/primary` | **P0** |
| Parser Expr | Predicates: `=,<>,<,>,<=,>=,IN,BETWEEN,LIKE [ESCAPE],IS [NOT] (NULL\|MISSING\|TRUE\|FALSE)` | **P0** |
| Parser Expr | `CASE [expr] WHEN … THEN … ELSE … END` | **P1** |
| Parser Expr | `CAST/CAN_CAST/CAN_LOSSLESS_CAST` | **P1** |
| Parser Expr | Path expressions: `.field`, `.*`, `[expr]`, `[*]`, chained | **P0** |
| Parser Expr | Variables: `@ident`, bare `ident`, `?` parameter | **P0** |
| Parser Expr | Literals: tuple `{}`, list `[]`, bag `<<>>`, ion backtick | **P0** (tuple/list) / **P1** (bag/ion) |
| Parser Expr | Subqueries | **P0** |
| Parser Expr | Aggregates: `COUNT(*)`, `COUNT/SUM/AVG/MIN/MAX([DISTINCT\|ALL] expr)` | **P1** |
| Parser Expr | Window: `LAG/LEAD(... ) OVER (PARTITION BY … ORDER BY …)` | **P1** |
| Parser Expr | Built-ins: `CHAR_LENGTH`, `OCTET_LENGTH`, `BIT_LENGTH`, `UPPER`, `LOWER`, `SIZE`, `EXISTS`, `TRIM`, `SUBSTRING`, `EXTRACT`, `DATE_ADD`, `DATE_DIFF`, `COALESCE`, `NULLIF` | **P1** |
| Parser Expr | Generic `ident(args,…)` function call | **P0** |
| Parser Expr | Graph patterns: `MATCH (n)-[e]->(m)`, quantifiers, selectors, restrictors | **P1** |
| Parser Cmd | `EXEC name args…` | **P1** |
| Type system | All standard SQL + PartiQL types (`MISSING`, `STRUCT`, `TUPLE`, `LIST`, `BAG`, `SEXP`, `BLOB`, `CLOB`, `SYMBOL`, `ANY`) | **P0** |
| Date/time literals | `DATE 'Y-M-D'`, `TIME [(p)] [WITH TIME ZONE] '…'`, `TIMESTAMP '…'` | **P1** |
| Analysis | Query validator (DQL-only) | **P0** |
| Completion | Table/column/keyword candidates with FROM-scope and SELECT-scope rules | **P0** |
| Catalog | DynamoDB-style metadata (tables = collections, columns = inferred fields) | **P2** (only when bytebase wires schema sync) |
| Semantic / quality / deparse | — | **P2** (no consumer today) |

## Notes for Planning

- **Bytebase coupling is shallow.** Only four `init()` registrations and one direct
  caller in the DynamoDB driver. Switching the import is mechanically small once
  splitter, validator, completer, and parser are in place.
- **Graph pattern matching is the largest single feature** in the legacy grammar that
  has no analogue in other omni engines — budget it as its own DAG node and don't
  underestimate it.
- **Two DML syntaxes coexist**, so the DML parser must accept both legacy path-based
  and RFC 0011 forms; tests should cover both spellings of `INSERT`/`UPDATE`/`DELETE`.
- **Ion mode in the lexer** is mode-switched on backticks — needs an explicit lexer
  state, not just a token rule.
- **The legacy `simple.sql` fixture is too small to be a useful regression suite.**
  Build the omni test corpus from the official PartiQL spec / AWS DynamoDB PartiQL
  reference instead, and cross-reference each clause back to the `.g4` rules above.
