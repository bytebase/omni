# MySQL Recursive Descent Parser Implementation Skill

You are implementing a recursive descent MySQL 8.0 parser.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `mysql/parser/`
**AST definitions:** `mysql/ast/`
**Reference:** [MySQL 8.0 SQL Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html)
**Tests:** `mysql/parser/compare_test.go`

## Your Task

1. Read `mysql/parser/PROGRESS.json`
2. Pick the next batch to work on:
   - If any batch has `"status": "in_progress"`, **resume that batch** (it was interrupted mid-work — read the existing code in its target file and continue from where it left off)
   - Otherwise, find the first batch with `"status": "pending"` whose dependencies (by id) are all `"done"`
   - If any batch has `"status": "failed"`, **retry it** (reset to `"in_progress"` and try again)
3. Implement that batch following the steps below
4. Update `PROGRESS.json`:
   - Set `"in_progress"` before starting work
   - Set `"done"` only after `go build` and `go test` pass
   - Set `"failed"` with `"error"` if you cannot make tests pass

If all batches are `"done"`, output `ALL_BATCHES_COMPLETE` and stop.

## Implementation Steps for Each Batch

### Step 1: Fetch Complete BNF from Official Documentation (MANDATORY)

**This is the most critical step. Do NOT skip it. Do NOT write BNF from memory.**

For every grammar rule in the batch, you MUST:

1. **Fetch the official MySQL 8.0 documentation page** using the WebFetch tool
   - URL pattern: `https://dev.mysql.com/doc/refman/8.0/en/{statement}.html`
   - e.g., `select.html`, `insert.html`, `create-table.html`, `alter-table.html`
   - For compound statements: `https://dev.mysql.com/doc/refman/8.0/en/{statement}.html`
     - e.g., `if.html`, `case-statement.html`, `while.html`, `repeat.html`, `declare-local-variable.html`
   - For XA: `https://dev.mysql.com/doc/refman/8.0/en/xa-statements.html`
2. **Extract the COMPLETE BNF/syntax diagram** from the page — every branch, every option, every sub-clause
3. **Do NOT abbreviate** — never write `...` or truncate the BNF. If a statement has 50 lines of BNF, write all 50 lines
4. **For sub-clauses that have their own doc page**, fetch those too (e.g., CREATE TABLE's `column_definition`, `partition_options` link to separate pages)

If WebFetch fails, use WebSearch to find the correct URL, then fetch it.

### Step 2: Read AST and Existing Code

- Read the AST node definitions from `mysql/ast/parsenodes.go` (and `node.go`)
- Read the existing parser code in `mysql/parser/` to understand available helpers and already-implemented parse functions
- If the existing AST types don't cover all BNF branches, add new node types / fields to `parsenodes.go`
- Read the test cases listed in the batch's `"tests"` list in `mysql/parser/compare_test.go`

### Step 3: Write Test Cases FIRST (Test-Driven Development)

**TEST-DRIVEN**: Write tests FIRST, then implement.

**Every branch of the BNF must have at least one test case.** For example, if ALTER TABLE has ADD COLUMN, DROP COLUMN, MODIFY COLUMN, CHANGE COLUMN, ADD INDEX, DROP INDEX, RENAME, CONVERT CHARSET, ALGORITHM, LOCK — then you need at least 10 test cases, one per branch.

1. **FIRST** write test cases in `mysql/parser/compare_test.go` covering every BNF rule and branch listed in the batch
2. Every BNF alternative/branch must have at least one test case
3. Include both positive tests (valid SQL that should parse) and negative tests (invalid SQL that should fail)
4. Include edge cases: empty clauses, maximum nesting, unusual but valid syntax

Add a new `TestParse{BatchName}` function.

### Step 4: Write Parse Functions

Create or update the target file (e.g., `mysql/parser/select.go`).

**Every parse function MUST have the COMPLETE BNF in its comment. This is a hard requirement.**

```go
// parseSelectStmt parses a SELECT statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/select.html
//
//	SELECT
//	    [ALL | DISTINCT | DISTINCTROW]
//	    [HIGH_PRIORITY]
//	    [STRAIGHT_JOIN]
//	    [SQL_CALC_FOUND_ROWS]
//	    select_expr [, select_expr] ...
//	    [FROM table_references
//	        [PARTITION partition_list]]
//	    [WHERE where_condition]
//	    [GROUP BY {col_name | expr | position}
//	        [ASC | DESC], ... [WITH ROLLUP]]
//	    [HAVING where_condition]
//	    [WINDOW window_name AS (window_spec)
//	        [, window_name AS (window_spec)] ...]
//	    [ORDER BY {col_name | expr | position}
//	        [ASC | DESC], ... [WITH ROLLUP]]
//	    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
//	    [into_option]
//	    [FOR {UPDATE | SHARE}
//	        [OF tbl_name [, tbl_name] ...]
//	        [NOWAIT | SKIP LOCKED]
//	      | LOCK IN SHARE MODE]
//	    [into_option]
//
//	into_option: {
//	    INTO OUTFILE 'file_name'
//	        [CHARACTER SET charset_name]
//	        export_options
//	  | INTO DUMPFILE 'file_name'
//	  | INTO var_name [, var_name] ...
//	}
func (p *Parser) parseSelectStmt() *ast.SelectStmt {
```

**The comment BNF must match the official docs exactly. No abbreviation, no `...` for omitted branches. Every branch in the BNF must have a corresponding code path in the function.**

**Rules for parse function implementation:**

1. **Use the existing AST types** from `mysql/ast/` — do NOT create new node types without updating both `parsenodes.go` and `outfuncs.go`
1a. **When implementing a new batch, you MUST also add serialization to `outfuncs.go` for any new node types used.** Every node type in `parsenodes.go` must have a corresponding case in `writeNode` and a `writeXxx` function in `outfuncs.go`.
1b. **Every branch in the BNF comment MUST have a corresponding implementation.** If the BNF says `{ ADD | DROP | MODIFY | CHANGE | RENAME | CONVERT | ALGORITHM | LOCK }`, you must handle ALL of them, not just a subset. If a branch requires a new AST node type or field, add it.
1c. **Sub-clauses must be recursively complete.** If a statement's BNF references `column_definition`, and `column_definition` itself has a full BNF (with DEFAULT, AUTO_INCREMENT, UNIQUE, PRIMARY KEY, COMMENT, COLLATE, GENERATED, ON UPDATE, etc.), you must fetch that sub-clause's BNF and implement it completely too.
2. **Record positions** on EVERY AST node that has a `Loc` field.
   Set `Loc: nodes.Loc{Start: p.pos()}` at the beginning of parsing that node, and
   set `node.Loc.End = p.pos()` at the end of parsing that node (after the last token is consumed).
   **This is a hard requirement.** Missing locations will cause features (SQL rewriting, error reporting) to break.
3. **Token constants** are in `mysql/parser/lexer.go` (keyword constants `kw*`, literal tokens `tok*`)
4. **Lexer types** (`Token`, `Lexer`, `NewLexer`) are in `mysql/parser/lexer.go`
5. **Error recovery**: When encountering unexpected tokens, try to recover by:
   - Skipping to the next semicolon for statement-level errors
   - Returning a partial AST node with what was parsed so far
6. **Operator precedence** in expressions: use Pratt parsing (precedence climbing)
7. **Incremental dispatch**: The `parseStmt` dispatch in `stmt.go` should be extended incrementally as each statement batch is implemented. Do not wait until batch 42 — wire in each statement parser as it is completed.

### Step 5: Wire Up

- For new statement types, wire them into `parseStmt()` / `parseCreateDispatch()` / `parseAlterDispatch()` / `parseDropDispatch()` in `stmt.go`
- Handle the `;`-separated list and wrap in `*ast.RawStmt` with position info
- Add any missing keyword constants to `lexer.go` as needed

### Step 6: Test

Run these commands in order:

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./mysql/...

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./mysql/parser/ -run "TestXxx|TestYyy"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./mysql/...
```

**Test strategy:**
- `mysql/parser/compare_test.go` contains `ParseAndCheck(t, sql)` which parses SQL and verifies:
  - Parse succeeds for valid SQL
  - AST structure matches expected output (via `ast.NodeToString()`)
  - Location fields are non-negative on all nodes
- For each batch, **add test cases to `compare_test.go`** using SQL strings relevant to the batch's grammar rules
- If a test fails, fix the parser implementation, never the test

### Step 7: Update Progress & Commit

Edit `PROGRESS.json`:
- **Before starting work**: set `"status"` to `"in_progress"`
- **After all tests pass**: set `"status"` to `"done"`
- **If tests fail and you cannot fix**: set `"status"` to `"failed"` and add `"error": "description"`
- **NEVER mark `"done"` if `go build ./mysql/...` or `go test ./mysql/...` fails**

After marking a batch as done, commit:
```bash
../../scripts/git-commit.sh "mysql/parser: implement batch N - name" mysql/
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `mysql/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `mysql/parser/lexer.go` | Lexer (Token, Lexer, NewLexer), keyword map, token constants |
| `mysql/ast/parsenodes.go` | AST node struct definitions |
| `mysql/ast/node.go` | Node interface, List, String, Integer |
| `mysql/ast/outfuncs.go` | AST serialization (NodeToString) |
| `mysql/parser/compare_test.go` | Test cases for parser validation |
| `mysql/parser/PROGRESS.json` | Batch tracking |

## MySQL-Specific Lexer Notes

The MySQL lexer handles these MySQL-specific constructs:
- **Backtick-quoted identifiers**: `` `table` ``, `` `column name` ``
- **Double-quoted strings**: When ANSI_QUOTES is off (default), `"hello"` is a string literal
- **`#` line comments**: In addition to `--` and `/* */`
- **`!=` operator**: Same as `<>`
- **No `::` typecast**: MySQL uses `CAST()` or `CONVERT()` instead
- **`@variable`**: User variables
- **`@@variable`**: System variables (@@global.var, @@session.var)
- **String charset introducer**: `_utf8 'hello'`, `_latin1 'hello'`
- **Hex literals**: `0xFF`, `X'FF'`
- **Bit literals**: `0b101`, `b'101'`
- **`<=>` operator**: NULL-safe equality
- **`DIV` operator**: Integer division
- **`:=` operator**: Assignment in SET and SELECT
- **`->` and `->>` operators**: JSON column-path extraction (MySQL 5.7+)

## Expression Parsing Strategy

For expressions, use **Pratt parsing / precedence climbing**:

```go
func (p *Parser) parseExpr(minPrec int) nodes.ExprNode {
    left := p.parsePrimary() // atoms: literals, column refs, subqueries, func calls, etc.
    for {
        prec, ok := p.infixPrecedence(p.cur.Type)
        if !ok || prec < minPrec {
            break
        }
        op := p.advance()
        right := p.parseExpr(prec + 1) // +1 for left-assoc
        left = &nodes.BinaryExpr{...}
    }
    return left
}
```

**Return type rules for parse functions:**
1. Expression parse functions (e.g., `parseExpr`, `parsePrimary`) should return `nodes.ExprNode`
2. Table reference parse functions (e.g., `parseTableRef`, `parseJoinClause`) should return `nodes.TableExpr`
3. Statement parse functions (e.g., `parseSelectStmt`, `parseInsertStmt`) should return `nodes.StmtNode`
4. When implementing a new batch, you MUST also add serialization for any new node types to `outfuncs.go` (add a `writeNode` switch case and a corresponding `writeXxx` function for each new struct)

Precedence levels (low to high, matching MySQL):
1. `OR`, `||`
2. `XOR`
3. `AND`, `&&`
4. `NOT` (prefix)
5. `BETWEEN`, `CASE`, `WHEN`, `THEN`, `ELSE`
6. `=`, `<=>`, `>=`, `>`, `<=`, `<`, `<>`, `!=`, `IS`, `LIKE`, `REGEXP`, `IN`
7. `|`
8. `&`
9. `<<`, `>>`
10. `-`, `+`
11. `*`, `/`, `DIV`, `%`, `MOD`
12. `^`
13. Unary `-`, `~`
14. `!`
15. `BINARY`, `COLLATE`
16. `INTERVAL`

## New AST Node Types Required for Phase 2 Batches (22+)

When implementing new batches, you will need to add new AST node types to `mysql/ast/parsenodes.go` and corresponding serialization in `mysql/ast/outfuncs.go`. Below is a guide for each batch area:

### CTE / WITH (batch 22)
- Extend `SelectStmt` with a `CTEs []*CommonTableExpr` field
- Add `CommonTableExpr` struct: Name, Columns, Select, Recursive bool

### Window OVER (batch 23)
- The `WindowDef` and `WindowFrame` types already exist
- Complete the TODO in `parseFuncCall` to actually parse the OVER clause and populate `FuncCallExpr.Over`

### JSON Operators (batch 24)
- Add `JsonExtractExpr` and `JsonUnquoteExtractExpr` (or use BinaryExpr with new op types `BinOpJsonExtract`, `BinOpJsonUnquoteExtract`)
- Add `MemberOfExpr` struct
- Add `JsonTableExpr` for JSON_TABLE() table function

### XA Transactions (batch 25)
- Add `XAStmt` struct with Type (START/END/PREPARE/COMMIT/ROLLBACK/RECOVER), Xid, and options

### CALL / HANDLER (batch 27)
- Add `CallStmt` struct: Name, Args
- Add `HandlerOpenStmt`, `HandlerReadStmt`, `HandlerCloseStmt`

### SIGNAL / RESIGNAL / GET DIAGNOSTICS (batch 28)
- Add `SignalStmt`, `ResignalStmt`, `GetDiagnosticsStmt`

### Compound Statements (batch 29)
- Add `CompoundStmt` (BEGIN...END), `DeclareVarStmt`, `DeclareConditionStmt`, `DeclareHandlerStmt`, `DeclareCursorStmt`

### Flow Control (batch 30)
- Add `IfStmt`, `CaseStmt` (distinct from CaseExpr), `WhileStmt`, `RepeatStmt`, `LoopStmt`, `LeaveStmt`, `IterateStmt`, `ReturnStmt`
- Add `OpenCursorStmt`, `FetchCursorStmt`, `CloseCursorStmt`

### Role Management (batch 31)
- Add `CreateRoleStmt`, `DropRoleStmt`, `SetDefaultRoleStmt`, `SetRoleStmt`

### TABLE / VALUES Statements (batch 33)
- Add `TableStmt` struct (simple form of SELECT * FROM)
- Add `ValuesStmt` struct with Rows

### Tablespace / Server (batch 37)
- Add `CreateTablespaceStmt`, `AlterTablespaceStmt`, `DropTablespaceStmt`
- Add `CreateServerStmt`, `AlterServerStmt`, `DropServerStmt`
- Add `CreateLogfileGroupStmt`, `AlterLogfileGroupStmt`, `DropLogfileGroupStmt`

### ALTER misc (batch 39)
- Add proper `AlterViewStmt`, `AlterEventStmt`, `AlterFunctionStmt`, `AlterProcedureStmt`
- Add proper `DropFunctionStmt`, `DropProcedureStmt`, `DropTriggerStmt`, `DropEventStmt` (replace current RawStmt stubs)

### Grouping / Lateral (batch 41)
- Add `GroupingSetsClause`, `CubeClause`, `RollupClause` as expression nodes in GROUP BY
- Add `LateralDerivedTable` as a table expression

## Known Gaps in Existing Batches (tracked as new batches)

- Batch 3 (expressions): Window OVER clause has TODO placeholder — tracked in batch 23
- Batch 4 (select): WITH/CTE not implemented — tracked in batch 22
- Batch 4 (select): LATERAL derived tables not implemented — tracked in batch 41
- Batch 4 (select): GROUP BY WITH ROLLUP listed but GROUPING SETS/CUBE/ROLLUP missing — tracked in batch 41
- Batch 12 (drop): DROP FUNCTION/PROCEDURE/TRIGGER/EVENT return RawStmt stubs — tracked in batch 39
- Batch 13 (set_show): Many SHOW variants missing — tracked in batch 32
- Batch 13 (set_show): EXPLAIN only handles SELECT, not INSERT/UPDATE/DELETE/REPLACE — tracked in batch 40
- Batch 14 (transaction): XA transactions missing — tracked in batch 25
- Batch 14 (transaction): SET TRANSACTION missing — tracked in batch 26
- Batch 21 (stmt_dispatch): ALTER dispatch missing VIEW/EVENT/FUNCTION/PROCEDURE — tracked in batch 39
- Batch 21 (stmt_dispatch): No dispatch for CALL, HANDLER, SIGNAL, RESIGNAL, GET, XA, TABLE, VALUES — tracked in batch 42

## Important Constraints

- Do NOT modify any files outside `mysql/`
- The `go.mod` at `/Users/rebeliceyang/Github/omni/go.mod` already exists
- Run `gofmt -w` on all created/modified files
- Use random sleep (1-3s) before git operations to avoid lock contention with other engine pipelines
