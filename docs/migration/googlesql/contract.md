# googlesql Parser-Consumption Contract (bytebase ANALYZE)

Scope: the omni `googlesql` parser is a **full cutover** for both engines that
consume `github.com/bytebase/parser/googlesql` in bytebase:

- `backend/plugin/parser/bigquery/` → `storepb.Engine_BIGQUERY`
- `backend/plugin/parser/spanner/`  → `storepb.Engine_SPANNER`

`grep -rln "bytebase/parser/googlesql"` confirms the legacy googlesql import is
confined to **exactly these two plugin dirs** (plus their `_test.go`). The only
external consumer of the packages is `backend/server/ultimate.go`, which does a
blank import (`_ "...parser/bigquery"`, `_ "...parser/spanner"`) purely to fire
the `init()` registrations. No exported parser symbol is called directly
cross-package — everything is reached through the `base.*` dispatcher maps keyed
by `storepb.Engine`.

Everything in the two plugin dirs that touches the parser is **P0** (full-cutover).

---

## 1. Contract table (must-work Go surface)

`p0_consumed` = the symbol is on the live product path for BIGQUERY and/or SPANNER
and the omni parser must reproduce its behavior.

| symbol | signature | callers[] | engines | p0_consumed |
|---|---|---|---|---|
| `bigquery.GetQuerySpan` | `func(ctx context.Context, gCtx base.GetQuerySpanContext, stmt base.Statement, database, _ string, ignoreCaseSensitive bool) (*base.QuerySpan, error)` | registered via `base.RegisterGetQuerySpan(Engine_BIGQUERY)` → reached by `base.GetQuerySpan` (masking / query-span / field-lineage path) | bigquery | yes |
| `spanner.GetQuerySpan` | `func(ctx context.Context, gCtx base.GetQuerySpanContext, stmt base.Statement, database, _ string, ignoreCaseSensitive bool) (*base.QuerySpan, error)` | registered via `base.RegisterGetQuerySpan(Engine_SPANNER)` → `base.GetQuerySpan` | spanner | yes |
| `bigquery.SplitSQL` | `func(statement string) ([]base.Statement, error)` | `base.RegisterSplitterFunc(Engine_BIGQUERY)` → `base.SplitMultiSQL`; also called internally by `ParseBigQuerySQL` | bigquery | yes |
| `spanner.SplitSQL` | `func(statement string) ([]base.Statement, error)` | `base.RegisterSplitterFunc(Engine_SPANNER)` → `base.SplitMultiSQL`; also called internally by `ParseSpannerGoogleSQL` | spanner | yes |
| `bigquery.Diagnose` | `func(_ context.Context, _ base.DiagnoseContext, statement string) ([]base.Diagnostic, error)` | `base.RegisterDiagnoseFunc(Engine_BIGQUERY)` → `base.Diagnose` (editor LSP diagnostics) | bigquery | yes |
| `spanner.Diagnose` | `func(_ context.Context, _ base.DiagnoseContext, statement string) ([]base.Diagnostic, error)` | `base.RegisterDiagnoseFunc(Engine_SPANNER)` → `base.Diagnose` | spanner | yes |
| `bigquery.ParseBigQuerySQL` | `func(statement string) ([]*base.ANTLRAST, error)` | shared parse wrapper; called by `querySpanExtractor.getQuerySpan` (bigquery). Not called cross-package. | bigquery | yes |
| `spanner.ParseSpannerGoogleSQL` | `func(sql string) ([]*base.ANTLRAST, error)` | shared parse wrapper; called by `querySpanExtractor.getQuerySpan` (spanner). Not called cross-package. | spanner | yes |
| `bigquery.parseSingleBigQuerySQL` | `func(statement string, baseLine int) (*base.ANTLRAST, error)` (unexported) | internal to `ParseBigQuerySQL` | bigquery | yes |
| `spanner.parseSingleSpannerGoogleSQL` | `func(statement string, baseLine int) (*base.ANTLRAST, error)` (unexported) | internal to `ParseSpannerGoogleSQL` | spanner | yes |
| `bigquery.parseGoogleSQLStatement` | `func(statement string) *base.SyntaxError` (unexported) | internal to `bigquery.Diagnose` | bigquery | yes |
| `spanner.parseGoogleSQLStatement` | `func(statement string) *base.SyntaxError` (unexported) | internal to `spanner.Diagnose` | spanner | yes |
| `(*querySpanExtractor)` (both pkgs) | extractor struct + `newQuerySpanExtractor(defaultDatabase string, gCtx base.GetQuerySpanContext, _ bool) *querySpanExtractor`, `getQuerySpan(ctx, stmt) (*base.QuerySpan, error)` | the engine of `GetQuerySpan`; walks the googlesql parse tree | both | yes |
| `queryTypeListener` (both pkgs) | `struct{ *parser.BaseGoogleSQLParserListener; allSystems bool; result base.QueryType; err error }` + `EnterStmts`, `getQueryTypeForUnterminatedSQLStatement` | walked inside `getQuerySpan` to classify the statement (Select / DDL / DML / Explain / SelectInfoSchema). **Internal-only — there is NO `RegisterQueryType*` in `base/`.** | both | yes |
| `accessTableListener` (both pkgs) | `struct{ *parser.BaseGoogleSQLParserListener; currentDatabase string; sourceColumnSet base.SourceColumnSet; err error }` + `EnterTable_path_expression`; entry `getAccessTables(defaultDatabase string, tree antlr.Tree) (base.SourceColumnSet, error)` | called inside `getQuerySpan`; produces `QuerySpan.SourceColumns` + system/user mix detection | both | yes |

### Legacy googlesql API surface the omni parser must reproduce (entry points)

The contract against `github.com/bytebase/parser/googlesql` (these are the exact
legacy entry points the handlers call):

- Lexer/parser constructors: `parser.NewGoogleSQLLexer(antlr.CharStream)`, `parser.NewGoogleSQLParser(antlr.TokenStream)`.
- Top rule: `p.Root()` returning `*parser.RootContext`; `RootContext.Stmts() *parser.StmtsContext`.
- Token constants: `parser.GoogleSQLLexerSEMI_SYMBOL`, `parser.GoogleSQLParserIDENTIFIER`, `parser.GoogleSQLParserDOT_SYMBOL`.
- Listener base: `parser.BaseGoogleSQLParserListener` (embedded by `queryTypeListener`, `accessTableListener`), walked with `antlr.ParseTreeWalkerDefault.Walk`.
- Grammar rule contexts traversed by the query-span extractor (the deep contract):
  `StmtsContext.AllUnterminated_sql_statement()`, `IUnterminated_sql_statementContext.Sql_statement_body()`,
  `ISql_statement_bodyContext` (`Query_statement`, `Alter_statement`, all `Create_*_statement`, `Rename_statement`, `Drop_*`, `Undrop_statement`, `Dml_statement`, `Merge_statement`, `Call_statement`, `Explain_statement`, `Set_statement`),
  `IQueryContext.Query_without_pipe_operators()`, `IQuery_without_pipe_operatorsContext` (`With_clause`, `Query_primary_or_set_operation`),
  `IWith_clauseContext` (`AllAliased_query`, `RECURSIVE_SYMBOL`), `IAliased_queryContext` (`Identifier`, `Parenthesized_query`),
  `IQuery_primary_or_set_operationContext` (`Query_primary`, `Query_set_operation`), `IQuery_set_operationContext.Query_set_operation_prefix()` (`Query_primary`, `AllQuery_set_operation_item`),
  `IQuery_primaryContext` (`Select_`, `Parenthesized_query`, `Opt_as_alias_with_required_as`), `IParenthesized_queryContext.Query()`,
  `ISelectContext` (`From_clause`, `Select_clause().Select_list().AllSelect_list_item()`), select-list item rules (`Select_column_star`, `Select_column_dot_star`, `Select_column_expr`, `Select_column_expr_with_as_alias`, `Star_modifiers`, `Star_except_list`, `Star_replace_list`, `Star_replace_item`),
  `IExpressionContext`, `IExpression_higher_prec_than_andContext` (the `a.b.c` DOT-chain walk in `getPossibleColumnResources`),
  `IFrom_clauseContext.From_clause_contents()` (`Table_primary`, `AllFrom_clause_contents_suffix`), `IFrom_clause_contents_suffixContext` (`COMMA_SYMBOL`, `JOIN_SYMBOL`, `Join_type`, `Table_primary`, `On_or_using_clause_list`),
  `ITable_primaryContext` (`Tvf_with_suffixes`, `Table_path_expression`, `Table_subquery`, `Join`, nested `Table_primary`), `IJoin_typeContext` (`CROSS/INNER/LEFT/RIGHT_SYMBOL`), `IJoin_itemContext.Join_type()`,
  `ITable_subqueryContext` (`Parenthesized_query`, `Opt_pivot_or_unpivot_clause_and_alias`), `ITable_path_expressionContext` (`Table_path_expression_base`, `Opt_pivot_or_unpivot_clause_and_alias`),
  `ITable_path_expression_baseContext` (`Unnest_expression`, `Maybe_slashed_or_dashed_path_expression`), `Maybe_dashed_path_expression`, `Dashed_path_expression`, `Slashed_path_expression`, `Path_expression.AllIdentifier()`,
  `IIdentifierContext`, `IUsing_clauseContext.AllIdentifier()`, `On_or_using_clause_list().AllOn_or_using_clause()`.

The omni parser must expose the equivalent shape for these (or the migration must
re-author the extractor against omni's AST — both extractors are 100% custom Go
tree-walks, not generated visitors).

---

## 2. Registered-functions table (per engine)

Confirmed by `grep -rn "base.Register" backend/plugin/parser/{bigquery,spanner}/`.
**Exactly 3 registrations per engine. No others.**

### BIGQUERY (`Engine_BIGQUERY`)

| registration func | engine key | handler | legacy googlesql API used |
|---|---|---|---|
| `base.RegisterGetQuerySpan` | `storepb.Engine_BIGQUERY` | `bigquery.GetQuerySpan` | `ParseBigQuerySQL` → `NewGoogleSQLLexer`/`NewGoogleSQLParser` → `p.Root()`; walks `RootContext`/`StmtsContext` + `accessTableListener` + `queryTypeListener` + custom select/from/expr tree-walk |
| `base.RegisterSplitterFunc` | `storepb.Engine_BIGQUERY` | `bigquery.SplitSQL` | `NewGoogleSQLLexer` + `CommonTokenStream.Fill()` → `base.SplitSQLByLexer(stream, GoogleSQLLexerSEMI_SYMBOL)` (**pure lexer/token-based split**) |
| `base.RegisterDiagnoseFunc` | `store.Engine_BIGQUERY` | `bigquery.Diagnose` | `parseGoogleSQLStatement`: `NewGoogleSQLLexer`/`NewGoogleSQLParser`, `BuildParseTrees=false`, `p.Root()`, collect `ParseErrorListener.Err` only |

### SPANNER (`Engine_SPANNER`)

| registration func | engine key | handler | legacy googlesql API used |
|---|---|---|---|
| `base.RegisterGetQuerySpan` | `storepb.Engine_SPANNER` | `spanner.GetQuerySpan` | `ParseSpannerGoogleSQL` → `NewGoogleSQLLexer`/`NewGoogleSQLParser` → `p.Root()`; same extractor/listeners as bigquery |
| `base.RegisterSplitterFunc` | `storepb.Engine_SPANNER` | `spanner.SplitSQL` | `NewGoogleSQLLexer`/`NewGoogleSQLParser`, `BuildParseTrees=true`, **full `p.Root()` parse**; walks `tree.Stmts().AllUnterminated_sql_statement()` and slices on `GoogleSQLLexerSEMI_SYMBOL` token boundaries (**parse-tree-based split — diverges from bigquery**) |
| `base.RegisterDiagnoseFunc` | `storepb.Engine_SPANNER` | `spanner.Diagnose` | `parseGoogleSQLStatement` (same as bigquery, but `ParseErrorListener.StartPosition` is set to `{Line:1}`) |

### Registrations NOT present for either engine (verified absent)

`RegisterQueryValidator` (read-only)*, `RegisterCompleteFunc`,
`RegisterStatementRangesFunc`, `RegisterParseStatementsFunc`,
`RegisterExtractChangedResourcesFunc`, `RegisterGenerateRestoreSQL`,
`RegisterTransformDMLToSelect`, `RegisterGetStatementTypes`, and there is no
`RegisterQueryType*` symbol in `base/` at all (query-type is an internal listener).

\* `RegisterQueryValidator` IS registered for both engines but **in
`backend/plugin/parser/standard/query.go`, not in the googlesql plugins** — see
Cross-engine reuse below.

---

## 3. Base interface/result types the omni parser must produce

From `backend/plugin/parser/base/` (signatures only).

### Dispatcher function types (`base/interface.go`)

```go
type GetQuerySpanFunc func(ctx context.Context, gCtx GetQuerySpanContext, stmt Statement, database, schema string, ignoreCaseSensitive bool) (*QuerySpan, error)
type SplitMultiSQLFunc func(string) ([]Statement, error)
type DiagnoseFunc     func(ctx context.Context, dCtx DiagnoseContext, statement string) ([]Diagnostic, error)
// (registered for other engines but NOT googlesql — listed for completeness)
type CompletionFunc            func(ctx context.Context, cCtx CompletionContext, statement string, caretLine, caretOffset int) ([]Candidate, error)
type StatementRangeFunc        func(ctx context.Context, sCtx StatementRangeContext, statement string) ([]Range, error)
type ParseStatementsFunc       func(statement string) ([]ParsedStatement, error)
type GetStatementTypesFunc     func([]AST) ([]storepb.StatementType, error)
type ValidateSQLForEditorFunc  func(string) (bool, bool, error)
type ExtractChangedResourcesFunc func(string, string, *model.DatabaseMetadata, []AST, string) (*ChangeSummary, error)
```

### Query span (`base/span.go`) — the heavy contract

```go
type QueryType int // QueryTypeUnknown, Select, Explain, SelectInfoSchema, DDL, DML

type QuerySpan struct {
    Type             QueryType
    Results          []QuerySpanResult
    SourceColumns    SourceColumnSet
    PredicateColumns SourceColumnSet
    PredicatePaths   map[string]*PathAST
    MongoDBAnalysis        *MongoDBAnalysis        // not used by googlesql
    ElasticsearchAnalysis  *ElasticsearchAnalysis  // not used by googlesql
    NotFoundError             error
    FunctionNotSupportedError error
}

type QuerySpanResult struct {
    Name             string
    SourceColumns    SourceColumnSet
    IsPlainField     bool
    SourceFieldPaths map[string][]*PathAST  // Cosmos only
    SelectAsterisk   bool                   // Cosmos only
}

type SourceColumnSet map[ColumnResource]bool
func MergeSourceColumnSet(m, n SourceColumnSet) (SourceColumnSet, bool)

type ColumnResource struct { Server, Database, Schema, Table, Column string }

type TableSource interface {            // produced by the extractor
    isTableSource()
    GetTableName() string
    GetSchemaName() string
    GetDatabaseName() string
    GetServerName() string
    GetQuerySpanResult() []QuerySpanResult
}
// concrete impls used by googlesql: PseudoTable{Name, Columns []QuerySpanResult}, PhysicalTable{Server,Database,Schema,Name, Columns []string}
// (PhysicalView, Sequence also implement TableSource but are unused by googlesql)

type GetQuerySpanContext struct {
    InstanceID                    string
    GetDatabaseMetadataFunc       GetDatabaseMetadataFunc
    ListDatabaseNamesFunc         ListDatabaseNamesFunc
    GetLinkedDatabaseMetadataFunc GetLinkedDatabaseMetadataFunc
    TempTables map[string]*PhysicalTable
    Engine     storepb.Engine
}
type GetDatabaseMetadataFunc func(context.Context, string, string) (string, *model.DatabaseMetadata, error)
type ListDatabaseNamesFunc   func(context.Context, string) ([]string, error)

var MixUserSystemTablesError = errors.Errorf("cannot access user and system tables at the same time")
```

### Statement / AST (`base/statement.go`, `base/ast.go`)

```go
type Statement struct {
    Text  string
    Empty bool
    Start *storepb.Position // 1-based inclusive
    End   *storepb.Position // 1-based exclusive
    Range *storepb.Range    // byte offsets
}
func (s *Statement) BaseLine() int   // Start.Line-1, or 0

type AST interface { ASTStartPosition() *storepb.Position }
type ANTLRAST struct {               // what the shared parse wrapper returns
    StartPosition *storepb.Position
    Tree          antlr.Tree
    Tokens        *antlr.CommonTokenStream
}
type ParsedStatement struct { Statement; AST AST } // not used by googlesql (no RegisterParseStatementsFunc)
```

### Diagnostics & errors (`base/diagnose.go`, `base/base.go`, `base/errors.go`)

```go
type Diagnostic = lsp.Diagnostic
type DiagnoseContext struct{}
func ConvertSyntaxErrorToDiagnostic(err *SyntaxError, statement string) Diagnostic

type SyntaxError struct { Position *storepb.Position; Message, RawMessage string }
type ParseErrorListener struct { StartPosition *storepb.Position; Err *SyntaxError; Statement string }

type ResourceNotFoundError struct { Err error; Server, DatabaseLink, Database, Schema, Table, Column, Function *string }
type TypeNotSupportedError struct { Err error; Type, Name, Extra string }
```

### Completion / ranges (`base/complete.go`, `base/interface.go`) — defined but unused by googlesql

```go
type Candidate struct { Text string; Type CandidateType; Definition, Comment string; Priority int }
type Range = lsp.Range
type StatementRangeContext struct{}
```

### Shared lexer/split helpers consumed by the handlers (`base/lexer.go`)

```go
func SplitSQLByLexer(stream *antlr.CommonTokenStream, semiTokenType int) ([]Statement, error) // bigquery.SplitSQL
func IsEmpty(tokens []antlr.Token, semi int) bool                                             // spanner.SplitSQL
func CalculateLineAndColumn(statement string, byteOffset int) (line, column int)              // spanner.SplitSQL
```

---

## 4. Shared vs per-engine divergence

The two plugin dirs are near-clones. The omni implementation can be a single
shared googlesql core parameterized by engine, **except** for the divergences
listed.

### Identical (omni can share verbatim)

| area | files | notes |
|---|---|---|
| Shared parse wrapper | `bigquery.go` ≈ `spanner.go` | `ParseBigQuerySQL`/`ParseSpannerGoogleSQL` + `parseSingle*` are byte-identical except the func name and one comment. Both: split → per-stmt `TrimRightFunc(IsSpaceOrSemicolon)+"\n;"` → lex/parse → `ParseErrorListener` → return `[]*base.ANTLRAST`. |
| `query_type.go` | both | `queryTypeListener`, `EnterStmts`, `getQueryTypeForUnterminatedSQLStatement` are **100% identical** (same DDL/DML/Explain/Select rule lists). |
| `dignose.go` | both | `Diagnose` + `parseGoogleSQLStatement` identical **except** spanner sets `ParseErrorListener.StartPosition = {Line:1}` while bigquery leaves it nil. (Cosmetic — affects only error-position base.) |
| `query_span.go` | both | `GetQuerySpan` wrapper identical. |
| query-span extractor core | `query_span_extractor.go` | CTE handling, set-ops, select-list, star-modifiers, join logic, `getPossibleColumnResources` DOT-chain walk, `accessTableListener` — all structurally identical. |

### Genuinely different (needs engine param or separate handling)

| divergence | bigquery | spanner | impact for omni |
|---|---|---|---|
| **Statement split** | `SplitSQL` = pure **lexer/token** split via `base.SplitSQLByLexer` (no parse) | `SplitSQL` = full **parse-tree** split (`BuildParseTrees=true`, `p.Root()`, walk `AllUnterminated_sql_statement`, slice on SEMI tokens), produces richer `Start`/`End`/`Range` and handles BEGIN/END/CASE blocks | **Two different split implementations.** Spanner needs the parse-tree splitter (BEGIN/END/CASE correctness — see `split_begin_end_test.go`); bigquery uses the cheap lexer split. Cannot fully share. |
| **getQuerySpan empty/zero-stmt handling** | tolerates 0 parse results → returns empty `QuerySpan`; sets `Type: base.Select` explicitly on success | errors on `!= 1` statement; does NOT set `Type` field on the success span (leaves zero value) | Minor behavioral diff in the GetQuerySpan wrapper; preserve per-engine. |
| **`getFieldColumnSource` signature** | `(_, tableName, fieldName string)` — 3 args, schema slot present | `(tableName, fieldName string)` — 2 args | Cosmetic; both ignore schema in practice. Spanner also iterates `tableSourceFrom` forward (not reversed) in the inner loop. |
| **`findTableSchema` resolution** | resolves **per-dataset**: calls `GetDatabaseMetadataFunc(instance, datasetName)`, uses `GetSchemaMetadata("")` (BigQuery dataset == database, no schema layer); `datasetName` defaults to `defaultDatabase` | resolves **per-schema within one DB**: always `GetDatabaseMetadataFunc(instance, defaultDatabase)` then `GetSchemaMetadata(schemaName)` (Spanner has named schemas under a single DB) | **Real metadata-model divergence** (BigQuery project.dataset.table vs Spanner db.schema.table). Keep engine-specific. |
| **`accessTableListener` schema classification** | 2nd-to-last identifier: only treated as `Schema` if it equals `INFORMATION_SCHEMA`, else it's `Database` (project/dataset) | 2nd-to-last identifier always set as `Schema` | Affects `SourceColumnSet` shape; keep per-engine. |
| **System-resource / info-schema detection** | `isSystemResource` = schema EqualFold `INFORMATION_SCHEMA` | `isSystemResource` = `INFORMATION_SCHEMA` **OR `SPANNER_SYS`** | Spanner adds `SPANNER_SYS`. Drives `SelectInfoSchema` classification + user/system mix error. |
| **default join type** (`getJoinTypeFromJoinType` fallthrough) | `default` → `crossJoin` | `default` → `innerJoin` | Subtle lineage diff for unqualified joins; keep per-engine. |
| **`PhysicalTable.Schema`** | `Schema: ""` (no schema layer) | `Schema: schemaName` | Consequence of the metadata-model divergence above. |

`query_type` and the query-span extractor *core* are shareable; **split**,
**metadata resolution**, and **system-schema detection** are the three real
engine-specific seams.

---

## 5. Cross-engine reuse

- **Read-only query validation (SQL-editor admission) is delegated to the
  `standard` regex parser, NOT googlesql.** `backend/plugin/parser/standard/query.go`
  `init()` registers `base.RegisterQueryValidator(Engine_BIGQUERY, ...)` and
  `(Engine_SPANNER, ...)` → `standard.ValidateSQLForEditor`, a quote/comment
  stripper + regex (`^SELECT`, `^EXPLAIN`, `^WITH` without DML). **The omni
  googlesql parser does NOT need to implement query validation** — it stays on
  the shared `standard` validator (same pattern as several non-ANTLR engines).
- **Within googlesql, bigquery and spanner do NOT delegate to each other**, and
  neither delegates to a third engine's parser (unlike Doris→MySQL). Each has its
  own copy of every handler. There is no `completion.go` or `resource_change.go`
  in either dir.
- Schema sync / `GetDatabaseMetadata` for both engines lives in the **DB drivers**
  (`backend/plugin/db/bigquery/`, `backend/plugin/db/spanner/`), built from live
  API/`information_schema` queries — **not** parser-derived, so out of the
  googlesql parser contract.

---

## 6. Feature support matrix (BIGQUERY & SPANNER)

Derived from the 3 registrations found + `backend/common/engine.go` `EngineSupport*`
gates. "Parser-wired" = backed by a googlesql registration; "declared" = the
`EngineSupport*` gate value; "—" = no parser involvement.

| product feature | gate (`common/engine.go`) | BIGQUERY | SPANNER | parser-wired? |
|---|---|---|---|---|
| Statement splitting | (always, via splitter map) | ✅ `SplitSQL` (lexer) | ✅ `SplitSQL` (parse-tree) | **yes — googlesql** |
| Syntax diagnostics (editor/LSP) | (via `base.Diagnose`) | ✅ `Diagnose` | ✅ `Diagnose` | **yes — googlesql** |
| Query span / field lineage | (via `base.GetQuerySpan`) | ✅ `GetQuerySpan` | ✅ `GetQuerySpan` | **yes — googlesql** |
| Query-type classification | (internal to GetQuerySpan) | ✅ `queryTypeListener` | ✅ `queryTypeListener` | **yes — googlesql (internal listener)** |
| Data masking | `EngineSupportMasking` = true | ✅ | ✅ | **yes — via GetQuerySpan** |
| Query ACL (new ACL) | `EngineSupportQueryNewACL` = true | ✅ | ✅ | indirect — relies on query span types |
| Read-only query validation | `RegisterQueryValidator` | ✅ | ✅ | **no — `standard` regex parser** |
| Auto-complete | `EngineSupportAutoComplete` = **false** | ❌ | ❌ | no (no `RegisterCompleteFunc`) |
| SQL review / lint | `EngineSupportSQLReview` = **false**; `EngineSupportStatementAdvise` = **false** | ❌ | ❌ | no (no advisor dir for googlesql) |
| Statement report (changed resources) | `EngineSupportStatementReport` = **false** | ❌ | ❌ | no (no `RegisterExtractChangedResourcesFunc`) |
| Statement ranges | (statement-range map) | ❌ | ❌ | no (no `RegisterStatementRangesFunc`) |
| Prior backup / restore | `EngineSupportPriorBackup` = **false** | ❌ | ❌ | no (no Transform/Restore registration) |
| Task-level syntax check | `EngineSupportSyntaxCheck` = **false** | ❌ | ❌ | no (only editor `Diagnose` is wired) |
| Schema sync / GetDatabaseMetadata | n/a | ✅ | ✅ | no — DB driver, not parser |
| Create database | `EngineSupportCreateDatabase` = false | ❌ | ❌ | — |
| Query-span plain-field | `EngineSupportQuerySpanPlainField` = false | ❌ | ❌ | (span returns lineage but plain-field flag off) |

**Net googlesql parser surface to satisfy for full cutover: 3 registered handlers
per engine (`GetQuerySpan`, `SplitSQL`, `Diagnose`) + their shared parse wrapper +
the query-type and access-table listeners + the full query-span extractor
tree-walk.** Everything else (completion, SQL review, statement report/ranges,
prior backup, task syntax check) is **not wired** for these engines today, and
query validation is owned by the `standard` parser.
