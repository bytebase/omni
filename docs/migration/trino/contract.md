# Trino Parser Migration Contract

## 1. Engine Mapping

| Field | Value |
|---|---|
| storepb.Engine value | `Engine_TRINO = 24` |
| v1pb.Engine value | `Engine_TRINO = 24` (same ordinal) |
| Presto | **No** separate engine constant exists. There is no `Engine_PRESTO`. Trino is the only engine value, and all bytebase Trino support is gated on `Engine_TRINO`. |
| Legacy parser import path | `github.com/bytebase/parser/trino` (ANTLR4-generated lexer + parser) |
| Bytebase parser package | `github.com/bytebase/bytebase/backend/plugin/parser/trino` |
| Bytebase schema plugin | `github.com/bytebase/bytebase/backend/plugin/schema/trino` |
| Bytebase DB driver | `github.com/bytebase/bytebase/backend/plugin/db/trino` |
| Side-loaded by | `backend/server/ultimate.go` via blank import of both `plugin/parser/trino` and `plugin/schema/trino` |

---

## 2. Registered Functions Table

All registrations happen in `init()` blocks inside the bytebase `plugin/parser/trino` package. The schema plugin registers separately.

| Registration call | Engine key | Handler | Source file |
|---|---|---|---|
| `base.RegisterSplitterFunc` | `Engine_TRINO` | `SplitSQL` | `split.go:13` |
| `base.RegisterQueryValidator` | `Engine_TRINO` | `validateQuery` | `query.go:9` |
| `base.RegisterGetQuerySpan` | `Engine_TRINO` | `GetQuerySpan` | `query_span.go:13` |
| `base.RegisterCompleteFunc` | `Engine_TRINO` | `Completion` | `completion.go:18` |
| `base.RegisterDiagnoseFunc` | `Engine_TRINO` | `Diagnose` | `diagnose.go:16` |
| `schema.RegisterGetDatabaseDefinition` | `Engine_TRINO` | `GetDatabaseDefinition` | `plugin/schema/trino/get_database_definition.go:12` |

**Not registered for Trino:**
- `RegisterParseStatementsFunc` — not present; Trino uses `ParseTrino` internally.
- `RegisterStatementRangesFunc` — not registered; `base.GetStatementRanges` returns empty for Trino.
- `RegisterQueryTypeFunc` / `RegisterStatementTypeGettersFunc` — not a base-registered hook; `getQueryType` / `GetStatementType` are package-internal helpers.
- `RegisterChangedResourcesGetter` — not registered.
- `RegisterTransformDMLToSelect` — not registered.
- `RegisterGenerateRestoreSQL` — not registered.
- SchemaDiff / masking advisor / plan-check — not registered.

---

## 3. Entry-Point Signatures

### 3.1 SplitSQL (split.go)

```go
func SplitSQL(statement string) ([]base.Statement, error)
```

Splits a multi-statement string into individual `base.Statement` values with byte range and position metadata. Tries parser-based splitting first (`splitByParser`), falls back to `tokenizer.NewTokenizer(...).SplitStandardMultiSQL()` on error.

### 3.2 validateQuery (query.go)

```go
func validateQuery(statement string) (bool, bool, error)
// registered as base.ValidateSQLForEditorFunc
// returns (canRunInReadOnly bool, returnsData bool, error)
```

Parses the full statement, classifies each sub-statement as SELECT/EXPLAIN/SelectInfoSchema (readonly+returns-data), or DML/DDL (not readonly). Returns `(false, false, nil)` on first non-read-only statement. EXPLAIN ANALYZE of a SELECT is allowed.

### 3.3 GetQuerySpan (query_span.go)

```go
func GetQuerySpan(
    ctx context.Context,
    gCtx base.GetQuerySpanContext,
    stmt base.Statement,
    database, schema string,
    ignoreCaseSensitive bool,
) (*base.QuerySpan, error)
```

Entry point registered with `base.RegisterGetQuerySpan`. Creates a `querySpanExtractor`, calls `extractor.getQuerySpan(ctx, stmt.Text)`, wraps errors.

### 3.4 Completion (completion.go)

```go
func Completion(
    ctx context.Context,
    cCtx base.CompletionContext,
    statement string,
    caretLine int,
    caretOffset int,
) ([]base.Candidate, error)
```

Tries `NewStandardCompleter` first; if it returns empty results, retries with `NewTrickyCompleter` (which additionally calls `skipHeadingSQLWithoutSemicolon`). Both use `base.CodeCompletionCore` (c3 algorithm).

### 3.5 Diagnose (diagnose.go)

```go
func Diagnose(
    _ context.Context,
    _ base.DiagnoseContext,
    statement string,
) ([]base.Diagnostic, error)
```

Parses `statement` with ANTLR. If any lexer or parser error listener fires, converts the `base.SyntaxError` to a `base.Diagnostic` via `base.ConvertSyntaxErrorToDiagnostic`. Returns empty slice on success.

### 3.6 ParseTrino (trino.go — internal, not base-registered)

```go
func ParseTrino(sql string) ([]*base.ANTLRAST, error)
```

Splits the SQL with `SplitSQL`, then calls `parseSingleTrino` for each non-empty statement. Each result is a `*base.ANTLRAST{StartPosition, Tree, Tokens}` where `Tree` is the ANTLR `SingleStatementContext`.

### 3.7 getQueryType / GetStatementType (query_type.go — internal)

```go
func getQueryType(node any) (base.QueryType, bool)
// returns (queryType, isExplainAnalyze)

func GetStatementType(tree any) StatementType
// returns a package-local StatementType constant
```

`getQueryType` walks the ANTLR tree using `queryTypeListener`, mapping statement kinds to `base.Select`, `base.Explain`, `base.SelectInfoSchema`, `base.DML`, `base.DDL`. Used by `validateQuery` and `getQuerySpan`.

### 3.8 GetDatabaseDefinition (plugin/schema/trino)

```go
func GetDatabaseDefinition(
    _ schema.GetDefinitionContext,
    metadata *storepb.DatabaseSchemaMetadata,
) (string, error)
```

Emits `CREATE TABLE IF NOT EXISTS "schema"."table" (...)` DDL from `DatabaseSchemaMetadata`. Called by `db/trino/dump.go` `(*Driver).Dump(...)`.

### 3.9 Exported helpers (trino.go — consumed by query_span_* internally)

```go
func NormalizeTrinoIdentifier(ident string) string
func ExtractQualifiedNameParts(ctx parser.IQualifiedNameContext) []string
func ExtractDatabaseSchemaName(ctx parser.IQualifiedNameContext, defaultDatabase, defaultSchema string) (string, string, string)
```

Not registered with `base`; used only within the parser package.

---

## 4. Capability Deep-Dive

### 4.1 query_span: extractor + listener + predicate

**Files:** `query_span_extractor.go`, `query_span_listener.go`, `query_span_predicate.go`

**What it extracts:**
- `SourceColumns base.SourceColumnSet` — fully-qualified `{Database, Schema, Table, Column}` tuples for every column accessed by the query. For SELECT queries, columns are resolved to their originating base-table columns via `getDatabaseMetadata`.
- `PredicateColumns base.SourceColumnSet` — columns appearing in WHERE / JOIN ON / JOIN USING predicates. Used by bytebase masking to detect when a sensitive column is used as a filter key.
- `Results []base.QuerySpanResult` — one entry per output column with `Name`, `SourceColumns`, `IsPlainField`, `SelectAsterisk`. SELECT * is deferred and expanded after all table sources are resolved.

**Walking strategy:**
1. `getQuerySpan` calls `ParseTrino`, checks query type. Non-SELECT/EXPLAIN queries return a basic span with just the accessed tables (via `extractAccessedTables`).
2. For SELECT: creates a `trinoQuerySpanListener`, walks the parse tree with `antlr.ParseTreeWalkerDefault.Walk`.
3. Listener methods: `EnterQuery` (WITH/CTE), `EnterTableName` (resolves table to metadata, builds `PhysicalTable` with column list), `EnterSelectSingle` / `EnterSelectAll` (builds result column entries), `EnterQuerySpecification` (WHERE predicates), `EnterJoinCriteria` (JOIN predicates), `EnterUnnest`, `EnterLateral`.
4. Post-walk: `expandTableReferencesToColumns` → `expandPredicateColumns` → `expandSelectAsteriskResults`.

**Metadata access:** `getDatabaseMetadata` calls `gCtx.GetDatabaseMetadataFunc` and caches results in `metaCache`. On `ResourceNotFoundError` the span still returns with `NotFoundError` set (not a hard error).

**Subquery handling:** `extractPredicateColumnFromSubquery` recursively creates a new `querySpanExtractor` and calls `getQuerySpan` on the subquery text. Predicate columns and result columns from subqueries are merged into the outer extractor.

**LATERAL / UNNEST:** dedicated ANTLR listener hooks create `PseudoTable` entries in `tableSourcesFrom` and propagate outer table sources for correlated column resolution.

**AST interfaces used:** `parser.IQualifiedNameContext`, `parser.IQueryContext`, `parser.IBooleanExpressionContext`, `parser.IJoinCriteriaContext`, `parser.TableNameContext`, `parser.QuerySpecificationContext`, `parser.SelectSingleContext`, `parser.SelectAllContext`, `parser.UnnestContext`, `parser.LateralContext` — all from `github.com/bytebase/parser/trino`.

### 4.2 Completion (completion.go — ~1730 lines)

**Strategy:** ANTLR3-style `CodeCompletionCore` (c3 / follow-set algorithm), identical in structure to other engines (TSQL, Redshift).

**preferredRules:**
```go
preferredRules = map[int]bool{
    trinoparser.TrinoParserRULE_identifier:    true,
    trinoparser.TrinoParserRULE_qualifiedName: true,
}
```

**ignoredTokens:** All identifier token types (IDENTIFIER, QUOTED_IDENTIFIER, DIGIT_IDENTIFIER, BACKTICK_IDENTIFIER), all literals (STRING, UNICODE_STRING, DECIMAL_VALUE, DOUBLE_VALUE, INTEGER_VALUE, BINARY_LITERAL), QUESTION_MARK, all operators and punctuation, EOF.

**Candidate types produced:** keywords, catalogs (`CandidateTypeDatabase`), schemas, tables, columns, views. Functions are categorized but the `functionEntries` map is never populated (placeholder).

**Two-pass strategy:**
1. `NewStandardCompleter` → `skipHeadingSQLs` (skip statements before caret)
2. `NewTrickyCompleter` → additionally `skipHeadingSQLWithoutSemicolon` (finds last `SELECT` at column 0)

**Table reference collection:** `collectLeadingTableReferences` / `collectRemainingTableReferences` scan the token stream for `FROM` tokens and parse each FROM clause with `parseTableReferences` (re-parses from the FROM keyword using `parser.Relation()`). Alias resolution handled by `tableRefListener`.

**CTE extraction:** `fetchCommonTableExpression` walks the full statement tree with `cteExtractor` listener. When no column aliases are specified for a CTE, it calls `GetQuerySpan` to derive column names.

**Subquery columns:** `tableRefListener.ExitSubquery` calls `GetQuerySpan(SELECT * FROM (<subquery>))` to derive column names for subquery aliases.

**Context determination:** `determineQualifiedNameContext` / `determineColumnReference` walk backward from the caret through the token stream to find dot-separated identifier chains, then return multiple possible `objectRefContext` interpretations (catalog.schema, schema.object, object.column) which are all tried.

**Identity quoting:** `quotedIdentifierIfNeeded` adds `"..."` if the identifier starts with a non-letter/underscore or contains non-word characters. Caret-inside-quoted-token is detected and suppressed.

### 4.3 query_type.go — Statement Classification

**Two classifiers:**

`getQueryType(node any) (base.QueryType, bool)` — used internally by `validateQuery` and `getQuerySpan`. Returns one of: `base.Select`, `base.Explain`, `base.SelectInfoSchema`, `base.DML`, `base.DDL`, `base.QueryTypeUnknown`. Second bool is `isExplainAnalyze`.

`GetStatementType(tree any) StatementType` — returns a package-local enum (Select, Explain, Insert, Update, Delete, Merge, CreateTable, CreateView, AlterTable, DropTable, DropView, CreateSchema, DropSchema, RenameTable, CreateTableAsSelect, Set, Show, Unsupported). Not used by any external caller currently, and not registered with base.

System schema detection: `containsSystemSchema` checks for `system.`, `information_schema.`, `$system.`, `catalog.`, `metadata.` prefixes (case-insensitive substring match) to classify queries as `SelectInfoSchema`.

### 4.4 diagnose.go — Syntax Diagnostics

`Diagnose` calls `parseTrinoStatement` which constructs a fresh ANTLR lexer + parser with `BuildParseTrees = false` (for speed), attaches `base.ParseErrorListener` on both, parses via `p.SingleStatement()`, and returns the first error as a `base.Diagnostic`. The function always returns a non-nil `[]base.Diagnostic` slice (possibly empty).

---

## 5. External Call Sites

The parser package is loaded purely through blank imports. All runtime dispatch goes through `base.*` functions keyed on engine.

### SQL Editor / Query API (`backend/api/v1/sql_service.go`)

| Call | Line(s) | Which registered API | Runtime-critical |
|---|---|---|---|
| `parserbase.SplitMultiSQL(instance.Metadata.GetEngine(), statement)` | 781, 789, 1089 | `RegisterSplitterFunc` → `SplitSQL` | **YES** — executed for every query submitted by a Trino instance |
| `parserbase.GetQuerySpan(ctx, gCtx, engine, statements, ...)` | 571–580, 674–683 | `RegisterGetQuerySpan` → `GetQuerySpan` | **YES** — called for masking and predicate-column checks on SELECT queries |
| `parserbase.ValidateSQLForEditor(engine, statement)` (via `validateQueryRequest`) | 1686 | `RegisterQueryValidator` → `validateQuery` | **YES** — called before every query execution to reject non-read-only statements |

### Query Statement Classification (`backend/api/v1/query_statement_classification.go`)

| Call | Line | Which API | Critical |
|---|---|---|---|
| `parserbase.SplitMultiSQL(engine, statement)` | 45 | `SplitSQL` | YES |
| `parserbase.GetQuerySpan(ctx, ..., engine, statements, ...)` | 49 | `GetQuerySpan` | YES |
| `parserbase.ValidateSQLForEditor(engine, statement)` | 11 | `validateQuery` | YES |

### LSP Completion (`backend/api/lsp/completion.go`)

| Call | Line | Which API | Critical |
|---|---|---|---|
| `parserbase.Completion(ctx, engine, cCtx, statement, caretLine, caretOffset)` | 58 | `RegisterCompleteFunc` → `Completion` | YES — gated by `common.EngineSupportAutoComplete` which returns `true` for `Engine_TRINO` |

### LSP Diagnostics (`backend/api/lsp/performance_optimizer.go`)

| Call | Line | Which API | Critical |
|---|---|---|---|
| `base.Diagnose(ctx, base.DiagnoseContext{}, engineType, content)` | 119 | `RegisterDiagnoseFunc` → `Diagnose` | YES — called for every text document change in the LSP |

### Export / Resource Extraction (`backend/component/export/resources.go`)

| Call | Line(s) | Which API | Critical |
|---|---|---|---|
| `base.SplitMultiSQL(engine, statement)` | 62, 169, 241 | `SplitSQL` | YES — runs on every data export |
| `base.GetQuerySpan(ctx, gCtx, engine, statements, ...)` | 66, 173, 245 | `GetQuerySpan` | YES |

### DB Driver — Dump (`backend/plugin/db/trino/dump.go`)

| Call | Line | Which API | Critical |
|---|---|---|---|
| `schema.GetDatabaseDefinition(storepb.Engine_TRINO, ...)` | 15 | `schema.RegisterGetDatabaseDefinition` → `GetDatabaseDefinition` | YES — called when dumping the Trino schema |

### DB Driver — trino.go, sync.go, query.go

The DB driver (`backend/plugin/db/trino/`) does **not** call any function from `backend/plugin/parser/trino` directly. It uses `util.SanitizeSQL` (a simple sanitizer, not the parser-based splitter) in `Execute`, and `util.RowsToQueryResult` for query results. The parser is reached only through the `base.*` dispatch layer.

### Masking

`common.EngineSupportMasking(Engine_TRINO)` returns `true` (`backend/common/engine.go:102`). This means `sql_service.go` will call `GetQuerySpan` and then `masker.MaskQueryResults` using `span.SourceColumns`. `PredicateColumns` are computed but predicate-column enforcement is only enabled for `Engine_MSSQL` (`predicate_columns_check.go:132`), so Trino predicate columns are extracted but not yet enforced.

---

## 6. Contract Entries

| Symbol | Signature | Callers | P0 consumed |
|---|---|---|---|
| `SplitSQL` | `func(statement string) ([]base.Statement, error)` | `sql_service.go:781,789,1089`; `query_statement_classification.go:45`; `export/resources.go:62,169,241`; LSP completion internally | **YES** |
| `validateQuery` | `func(statement string) (bool, bool, error)` | `sql_service.go:1686`; `query_statement_classification.go:11`; `sql_service.go:1840` | **YES** |
| `GetQuerySpan` | `func(ctx context.Context, gCtx base.GetQuerySpanContext, stmt base.Statement, database, schema string, ignoreCaseSensitive bool) (*base.QuerySpan, error)` | `sql_service.go:571,674`; `query_statement_classification.go:49`; `export/resources.go:66,173,245`; `completion.go:1397,1499` (internal) | **YES** |
| `Completion` | `func(ctx context.Context, cCtx base.CompletionContext, statement string, caretLine int, caretOffset int) ([]base.Candidate, error)` | `lsp/completion.go:58` | **YES** |
| `Diagnose` | `func(_ context.Context, _ base.DiagnoseContext, statement string) ([]base.Diagnostic, error)` | `lsp/performance_optimizer.go:119` | **YES** |
| `GetDatabaseDefinition` | `func(_ schema.GetDefinitionContext, metadata *storepb.DatabaseSchemaMetadata) (string, error)` | `db/trino/dump.go:15` | **YES** |
| `ParseTrino` | `func(sql string) ([]*base.ANTLRAST, error)` | Internal only (`query_span_extractor.go`, `query.go`, `diagnose.go`) | No (internal helper) |
| `getQueryType` | `func(node any) (base.QueryType, bool)` | Internal only | No (internal) |
| `GetStatementType` | `func(tree any) StatementType` | Internal only — not called by any external package | No |
| `NormalizeTrinoIdentifier` | `func(ident string) string` | Internal only | No |
| `ExtractQualifiedNameParts` | `func(ctx parser.IQualifiedNameContext) []string` | Internal only | No |
| `ExtractDatabaseSchemaName` | `func(ctx parser.IQualifiedNameContext, defaultDatabase, defaultSchema string) (string, string, string)` | Internal only | No |

### P0 summary

All six registered handlers are P0 — they are invoked at runtime by real bytebase features whenever a Trino instance is connected:

1. **`SplitSQL`** — mandatory for every query execute, export, and LSP session.
2. **`validateQuery`** — mandatory read-only guard on the SQL editor.
3. **`GetQuerySpan`** — mandatory for masking, resource extraction, and export ACL.
4. **`Completion`** — mandatory for LSP autocomplete (gated by `EngineSupportAutoComplete=true`).
5. **`Diagnose`** — mandatory for LSP inline diagnostics.
6. **`GetDatabaseDefinition`** (schema plugin) — mandatory for schema dump / SDL generation.

`GetStatementRangesFunc` is **not registered** for Trino (the base dispatcher returns an empty range list); this is **not** a gap that blocks the import switch.

There are no masking-specific registrations (no `ExtractChangedResourcesFunc`, no schema diff, no plan-check advisors) for Trino, so those features are absent and not part of the P0 surface.
