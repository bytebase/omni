# Recursive Descent Parser Implementation Skill - T-SQL

You are implementing a recursive descent T-SQL (SQL Server) parser.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `mssql/parser/`
**Reference:** Microsoft T-SQL documentation (https://learn.microsoft.com/en-us/sql/t-sql/)
**AST definitions:** `mssql/ast/`
**Tests:** `mssql/parser/compare_test.go`

## Your Task

1. Read `mssql/parser/PROGRESS_SUMMARY.json` (lightweight version — done batches are compressed, only pending/failed/in_progress have full detail)
2. Pick the next batch to work on:
   - If any batch has `"status": "in_progress"`, **resume that batch** (it was interrupted mid-work -- read the existing code in its target file and continue from where it left off)
   - Otherwise, find the first batch with `"status": "pending"` whose dependencies (by id) are all `"done"`
   - If any batch has `"status": "failed"`, **retry it** (reset to `"in_progress"` and try again)
3. Implement that batch following the steps below
4. Update `mssql/parser/PROGRESS.json` (the full file, NOT the summary):
   - Set `"in_progress"` before starting work
   - Set `"done"` only after `go build` and `go test` pass
   - Set `"failed"` with `"error"` if you cannot make tests pass

If all batches are `"done"`, output `ALL_BATCHES_COMPLETE` and stop.

## Progress Logging (MANDATORY)

You MUST print progress markers to stdout at each step. This is how the pipeline operator monitors your work. Use this exact format:

```
[BATCH N] STARTED - batch_name
[BATCH N] STEP reading_refs - Reading BNF and AST definitions
[BATCH N] STEP writing_tests - Writing test cases
[BATCH N] STEP writing_code - Implementing parse functions
[BATCH N] STEP build - Running go build
[BATCH N] STEP test - Running go test (X passed, Y failed)
[BATCH N] STEP commit - Committing changes
[BATCH N] DONE
```

If a step fails, print:
```
[BATCH N] FAIL test - description of failure
[BATCH N] RETRY - what you're fixing
```

**Do NOT skip these markers.** They appear in the build log and are essential for debugging pipeline issues.

## Implementation Steps for Each Batch

### Step 1: Read Official Documentation (MANDATORY)

**This is the most critical step. Do NOT skip it. Do NOT write BNF from memory.**

Documentation has been **pre-fetched** to local files. For every grammar rule in the batch:

1. **First check local docs**: Read from `mssql/parser/docs/{statement}.txt` (e.g., `alter-table-transact-sql.txt`, `create-trigger-transact-sql.txt`)
2. **Only if the local file is missing**, use WebFetch as fallback:
   - `https://learn.microsoft.com/en-us/sql/t-sql/statements/{statement}-transact-sql`
   - For queries: `https://learn.microsoft.com/en-us/sql/t-sql/queries/{query}-transact-sql`
3. **Extract the COMPLETE BNF/syntax diagram** — every branch, every option, every sub-clause
4. **Do NOT abbreviate** — never write `...` or truncate the BNF
5. **For sub-clauses that have their own doc page**, check `mssql/parser/docs/` first, then WebFetch if missing

### Step 2: Read AST and Existing Code

- Read the AST node definitions from `mssql/ast/parsenodes.go` (and `node.go`)
- Read the existing parser code in `mssql/parser/` to understand available helpers and already-implemented parse functions
- If the existing AST types don't cover all BNF branches, add new node types / fields to `parsenodes.go`

### Step 3: Write Tests FIRST (Test-Driven Development)

**TEST-DRIVEN**: Write tests FIRST, then implement.

**Every branch of the BNF must have at least one test case.** For example, if ALTER TABLE has ADD COLUMN, DROP COLUMN, ALTER COLUMN, ADD CONSTRAINT, DROP CONSTRAINT, ENABLE TRIGGER, DISABLE TRIGGER, SWITCH PARTITION, REBUILD, SET — then you need at least 10 test cases, one per branch.

Add test cases to `compare_test.go` using SQL strings relevant to the batch's grammar rules.
Add a new `TestParse{BatchName}` function.

### Step 4: Write Parse Functions

Create or update the target file (e.g., `mssql/parser/select.go`).

**Every parse function MUST have the COMPLETE BNF in its comment. This is a hard requirement.**

```go
// parseAlterTableStmt parses an ALTER TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-table-transact-sql
//
//  ALTER TABLE [ database_name . [ schema_name ] . | schema_name . ] table_name
//  {
//      ALTER COLUMN column_name
//      {
//          [ new_data_type [ ( precision [ , scale ] ) ] ]
//          [ COLLATE collation_name ]
//          [ NULL | NOT NULL ]
//          | { ADD | DROP } { ROWGUIDCOL | PERSISTED | NOT FOR REPLICATION | SPARSE | HIDDEN }
//          | { ADD | DROP } MASKED [ WITH ( FUNCTION = 'mask_function' ) ]
//      }
//      | [ WITH { CHECK | NOCHECK } ] ADD
//          { column_definition | computed_column_definition | table_constraint } [ ,...n ]
//      | DROP { [ CONSTRAINT ] [ IF EXISTS ] constraint_name [ ,...n ] | COLUMN [ IF EXISTS ] column_name [ ,...n ] }
//      | [ WITH { CHECK | NOCHECK } ] { CHECK | NOCHECK } CONSTRAINT { ALL | constraint_name [ ,...n ] }
//      | { ENABLE | DISABLE } TRIGGER { ALL | trigger_name [ ,...n ] }
//      | SWITCH [ PARTITION source_partition_number_expression ] TO target_table
//          [ PARTITION target_partition_number_expression ] [ WITH ( ... ) ]
//      | SET ( FILESTREAM_ON = { partition_scheme_name | filegroup | "default" | "NULL" } )
//      | REBUILD [ [PARTITION = ALL] [ WITH ( rebuild_index_option [ ,...n ] ) ]
//               | [ PARTITION = partition_number [ WITH ( single_partition_rebuild_index_option [ ,...n ] ) ] ] ]
//      | <table_option>
//  }
func (p *Parser) parseAlterTableStmt() *nodes.AlterTableStmt {
```

**The comment BNF must match the official docs exactly. No abbreviation, no `...` for omitted branches. Every branch in the BNF must have a corresponding code path in the function.**

**Return type conventions:**
- Expression parsers return `ExprNode` (not `Node`)
- Table reference parsers return `TableExpr`
- Statement parsers return `StmtNode`

**Rules for parse function implementation:**

1. **Use the existing AST types** from `mssql/ast/` -- do NOT create new node types without updating both parsenodes.go and outfuncs.go
1a. **When implementing a new batch, you MUST also add serialization to `outfuncs.go` for any new node types used.** Every node type in parsenodes.go must have a corresponding case in `writeNode` and a `writeXxx` function in outfuncs.go.
1b. **Every branch in the BNF comment MUST have a corresponding implementation.** If the BNF says `{ ADD | DROP | ALTER COLUMN | ENABLE TRIGGER | DISABLE TRIGGER | SWITCH | REBUILD | SET }`, you must handle ALL of them, not just a subset. If a branch requires a new AST node type or field, add it.
1c. **Sub-clauses must be recursively complete.** If a statement's BNF references `column_definition`, and `column_definition` itself has a full BNF (with DEFAULT, IDENTITY, CONSTRAINT, COLLATE, GENERATED, MASKED, etc.), you must fetch that sub-clause's BNF and implement it completely too.
2. **Record positions** on EVERY AST node that has a `Loc` field.
   Set `Loc: nodes.Loc{Start: p.pos()}` at the beginning of parsing a node.
   Set `node.Loc.End = p.pos()` at the end of parsing a node.
   **This is a hard requirement.**
3. **Token constants** are in `mssql/parser/lexer.go`
4. **Keywords are case-insensitive** -- the lexer handles this
5. **Error recovery**: When encountering unexpected tokens, try to recover by:
   - Skipping to the next semicolon for statement-level errors
   - Recording errors but continuing to parse subsequent statements

### Step 4: Test

Run these commands in order:

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./mssql/...

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./mssql/parser/ -run "TestParse{BatchName}"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./mssql/...
```

### Step 5: Update Progress

Edit `PROGRESS.json`:
- **Before starting work**: set `"status"` to `"in_progress"` so a restart knows this batch was interrupted
- **After all tests pass**: set `"status"` to `"done"`
- **If tests fail and you cannot fix**: set `"status"` to `"failed"` and add `"error": "description"`
- **NEVER mark `"done"` if `go build ./mssql/...` or `go test ./mssql/...` fails**

### Step 6: Git Commit

```bash
../../scripts/git-commit.sh "mssql/parser: implement batch N - name" mssql/
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `mssql/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `mssql/parser/lexer.go` | Lexer, Token, keywords, token constants |
| `mssql/ast/parsenodes.go` | AST node struct definitions |
| `mssql/ast/node.go` | Node interface, List, String, Integer, etc. |
| `mssql/ast/outfuncs.go` | NodeToString for AST comparison |
| `mssql/parser/compare_test.go` | Test infrastructure |
| `mssql/parser/PROGRESS.json` | Implementation progress tracking |

## T-SQL Specific Notes

- **Bracketed identifiers**: `[column name]` uses square brackets for quoting
- **Variables**: `@local_var`, `@@global_var`
- **N-strings**: `N'unicode string'` for nvarchar literals
- **Comments**: `--` line and `/* */` block (nested), no `#` comments
- **Not-equal**: Both `!=` and `<>` are valid
- **Special operators**: `!<` (not less than), `!>` (not greater than)
- **String concat**: `+` operator (same as addition)
- **Static methods**: `type::Method()` syntax
- **GO**: Batch separator (not actually a SQL statement)
- **TOP**: Uses parentheses: `TOP (10)` (parentheses optional for literals)
- **CROSS APPLY / OUTER APPLY**: T-SQL specific join types
- **OUTPUT clause**: Available on INSERT, UPDATE, DELETE, MERGE
- **MERGE**: Full MERGE statement with WHEN MATCHED/NOT MATCHED/BY SOURCE
- **TRY...CATCH**: Error handling blocks
- **IIF**: T-SQL specific inline IF function
- **CONVERT**: T-SQL specific type conversion with optional style

## New AST Node Types Required for Phase 2 Batches (23+)

When implementing new batches, you will need to add new AST node types to `mssql/ast/parsenodes.go` and corresponding serialization in `mssql/ast/outfuncs.go`. Below is a guide for each batch area:

### Cursor Operations (batch 23)
- `OpenCursorStmt` -- OPEN cursor_name
- `FetchCursorStmt` -- FETCH [NEXT|PRIOR|FIRST|LAST|ABSOLUTE n|RELATIVE n] FROM cursor INTO @vars
- `CloseCursorStmt` -- CLOSE cursor_name
- `DeallocateCursorStmt` -- DEALLOCATE cursor_name
- Extend `VariableDecl.IsCursor` to support full cursor options (SCROLL, STATIC, KEYSET, DYNAMIC, FAST_FORWARD, etc.)

### CREATE TRIGGER (batch 24)
- `CreateTriggerStmt` -- with OrAlter, Name, Table, TriggerType (AFTER/INSTEAD OF/FOR), Events (INSERT/UPDATE/DELETE), Body
- For DDL triggers: Events are DDL event types (CREATE_TABLE, ALTER_TABLE, etc.), scope is DATABASE or ALL SERVER

### CREATE SCHEMA (batch 25)
- `CreateSchemaStmt` -- Name, Authorization, optional contained CREATE TABLE/VIEW/GRANT statements
- `AlterSchemaStmt` -- TRANSFER entity

### CREATE TYPE (batch 26)
- `CreateTypeStmt` -- for alias types (FROM base_type), table types (AS TABLE (...)), CLR types (EXTERNAL NAME)

### CREATE SEQUENCE (batch 27)
- `CreateSequenceStmt` / `AlterSequenceStmt` -- with DataType, StartWith, IncrementBy, MinValue, MaxValue, Cycle, Cache options

### CREATE SYNONYM (batch 28)
- `CreateSynonymStmt` -- Name, ForName (target object)

### ALTER Objects (batch 29)
- Extend `parseAlterStmt` dispatcher for DATABASE, INDEX, VIEW, PROCEDURE, FUNCTION
- `AlterDatabaseStmt` -- SET options, MODIFY NAME/FILE, ADD FILE, etc.
- `AlterIndexStmt` -- REBUILD, REORGANIZE, DISABLE, SET options

### Security Principals (batch 30)
- `CreateUserStmt`, `AlterUserStmt` -- WITH PASSWORD, DEFAULT_SCHEMA, LOGIN, etc.
- `CreateLoginStmt`, `AlterLoginStmt` -- WITH PASSWORD, DEFAULT_DATABASE, CHECK_POLICY, etc.
- `CreateRoleStmt`, `AlterRoleStmt` -- ADD/DROP MEMBER
- `CreateAppRoleStmt`

### Security Keys/Certs (batch 31)
- `CreateMasterKeyStmt`, `CreateSymmetricKeyStmt`, `CreateAsymmetricKeyStmt`, `CreateCertificateStmt`, `CreateCredentialStmt`
- `OpenSymmetricKeyStmt`, `CloseSymmetricKeyStmt`
- `BackupCertificateStmt`

### DBCC (batch 32)
- `DbccStmt` -- Command (string), Args (list of expr), WithOptions

### BULK INSERT (batch 33)
- `BulkInsertStmt` -- Table, FromFile, WithOptions

### BACKUP/RESTORE (batch 34)
- `BackupStmt` -- Type (DATABASE/LOG), Database, ToDevices, WithOptions
- `RestoreStmt` -- Type (DATABASE/LOG/HEADERONLY/FILELISTONLY/etc.), Database, FromDevices, WithOptions

### PIVOT/UNPIVOT (batch 35)
- `PivotClause` -- AggFunc, ForColumn, InValues, Alias
- `UnpivotClause` -- ValueColumn, ForColumn, InColumns, Alias
- These are table sources, so they should implement `tableExpr()`

### Rowset Functions (batch 37)
- `OpenRowsetExpr`, `OpenQueryExpr`, `OpenJsonExpr`, `OpenDataSourceExpr`, `OpenXmlExpr`
- These should implement both `tableExpr()` (for FROM) and possibly `exprNode()`
- OPENJSON and OPENXML support WITH clause for column definitions

### Grouping Sets (batch 38)
- `GroupingSetsClause`, `CubeClause`, `RollupClause` -- as expression nodes in GROUP BY

### Partition (batch 42)
- `CreatePartitionFunctionStmt`, `CreatePartitionSchemeStmt`
- `AlterPartitionFunctionStmt`, `AlterPartitionSchemeStmt`

### Fulltext (batch 43)
- `CreateFulltextIndexStmt`, `CreateFulltextCatalogStmt`
- `ContainsPredicate`, `FreetextPredicate` -- as expression nodes in WHERE

### Service Broker (batch 46)
- `CreateMessageTypeStmt`, `CreateContractStmt`, `CreateQueueStmt`, `CreateServiceStmt`
- `SendStmt`, `ReceiveStmt`, `BeginConversationStmt`, `EndConversationStmt`

### SET Options (batch 41)
- Extend existing `SetStmt` or add `SetOptionStmt` for ON/OFF options and special SET forms like SET TRANSACTION ISOLATION LEVEL

### Dispatcher (batch 49)
- Wire all new keywords into `parseStmt`, `parseCreateStmt`, `parseAlterStmt`
- Add new keyword constants to `lexer.go` if needed (check existing keywords first -- many are already defined)
