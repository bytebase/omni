# Recursive Descent Parser Implementation Skill (Oracle PL/SQL)

You are implementing a recursive descent Oracle PL/SQL parser.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `oracle/parser/`
**AST definitions:** `oracle/ast/`
**Tests:** `oracle/parser/compare_test.go`

## Reference Documentation

- **Oracle Database 23c SQL Language Reference**: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/
- **Oracle Database 23c PL/SQL Language Reference**: https://docs.oracle.com/en/database/oracle/oracle-database/23/lnpls/
- Also reference the ANTLR2-based open-source Oracle grammar for understanding edge cases

## Your Task

1. Read `oracle/parser/PROGRESS.json`
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

### Step 1: Read References

- Read the AST node definitions from `oracle/ast/parsenodes.go` (and `node.go`)
- Read the existing parser code in `oracle/parser/` to understand available helpers
- Read Oracle documentation BNF for each statement type

### Step 2: Fetch Oracle Documentation BNF

- For each major statement type, search the Oracle 23c documentation for its syntax
- URL pattern: `https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/{command}.html`

### Step 3: Write Parse Functions

Create or update the target file (e.g., `oracle/parser/select.go`).

**Every parse function MUST have this comment format:**

```go
// parseSelectStmt parses a SELECT statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	SELECT [ hint ] [ ALL | DISTINCT | UNIQUE ]
//	    select_list
//	    [ FROM table_reference [, ...] ]
//	    [ WHERE condition ]
//	    [ hierarchical_query_clause ]
//	    ...
func (p *Parser) parseSelectStmt() *ast.SelectStmt {
```

**Return type conventions:**
- Expression parsers return `ExprNode` (e.g., `parseExpr() ast.ExprNode`)
- Table reference parsers return `TableExpr` (e.g., `parseTableRef() ast.TableExpr`)
- Statement parsers return `StmtNode` (e.g., `parseSelectStmt() *ast.SelectStmt` which satisfies `StmtNode`)

**Rules for parse function implementation:**

1. **Use the existing AST types** from `oracle/ast/` — do NOT create new node types unless absolutely necessary
2. **Record positions** on EVERY AST node that has a `Loc` field.
   Set `Loc: nodes.Loc{Start: p.pos()}` at the beginning of parsing a node.
   Set `node.Loc.End = p.pos()` at the end of parsing a node.
   **This is a hard requirement.**
3. **Token constants** are keyword constants (kw*) and lexer tokens (tok*) in `oracle/parser/lexer.go`
4. **Match the Oracle grammar semantics exactly**
5. **Error recovery**: When encountering unexpected tokens, try to recover by:
   - Skipping to the next semicolon for statement-level errors
   - Returning a partial AST node with what was parsed so far
6. **Operator precedence** in expressions: use Pratt parsing (precedence climbing)
7. **Serialization**: When implementing a new batch, you MUST also add serialization to `oracle/ast/outfuncs.go` for any new node types used. Every node type that has a `nodeTag()` method must have a corresponding case in `writeNode` and a `writeXxx` function.
8. **Incremental dispatch**: The `parseStmt` dispatch in `parser.go` should be extended incrementally as each statement batch is implemented. Do not wait until batch 23 — wire in each statement parser as it is completed.

### Step 4: Test

**TEST-DRIVEN**: Write tests FIRST, then implement. Every BNF branch must have a test.

Run these commands in order:

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./oracle/...

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestXxx"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./oracle/...
```

### Step 5: Update Progress

Edit `PROGRESS.json`:
- **Before starting work**: set `"status"` to `"in_progress"`
- **After all tests pass**: set `"status"` to `"done"`
- **If tests fail and you cannot fix**: set `"status"` to `"failed"` and add `"error": "description"`
- **NEVER mark `"done"` if `go build ./oracle/...` or `go test ./oracle/...` fails**

### Step 6: Git Commit

```bash
../../scripts/git-commit.sh "oracle/parser: implement batch N - name" oracle/
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `oracle/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `oracle/parser/lexer.go` | Lexer (Token, Lexer, NewLexer), keywords |
| `oracle/ast/node.go` | Node interface, List, String, Integer |
| `oracle/ast/parsenodes.go` | AST node struct definitions, enums |
| `oracle/ast/outfuncs.go` | AST serialization (NodeToString) |
| `oracle/parser/compare_test.go` | Test infrastructure |
| `oracle/parser/PROGRESS.json` | Batch progress tracking |

## Oracle-Specific Parsing Notes

### Lexer Specifics
- `"double-quoted identifiers"` (case-sensitive when quoted)
- Single-quoted strings with `''` escape (no backslash escape)
- `q'[delimiter]...[delimiter]'` alternative quoting mechanism
- `--` line comments, `/* */` block comments, `/*+ hints */`
- `||` string concatenation, `:=` assignment, `=>` named param, `..` range, `@` dblink, `**` exponent
- Bind variables: `:name`, `:1`
- National character literals: `N'text'`

### Oracle Dialect Differences from PostgreSQL
- MINUS instead of EXCEPT
- CONNECT BY / START WITH hierarchical queries
- PRIOR operator
- DECODE() function (like CASE but different syntax)
- (+) outer join syntax
- PIVOT / UNPIVOT
- INSERT ALL / INSERT FIRST (multi-table insert)
- MERGE statement with different syntax
- PL/SQL blocks (DECLARE/BEGIN/END/EXCEPTION)
- No dollar-quoted strings
- No backtick identifiers
- ROWID, ROWNUM, LEVEL pseudo-columns
- FLASHBACK queries (AS OF)
- MODEL clause
- SAMPLE clause

## Expression Parsing Strategy

Use **Pratt parsing / precedence climbing**:

Precedence levels (from low to high, matching Oracle):
1. OR
2. AND
3. NOT (prefix)
4. IS, IS NOT
5. comparison: =, <, >, <=, >=, !=, <>
6. LIKE, LIKEC, LIKE2, LIKE4, BETWEEN, IN
7. string concatenation: ||
8. addition: +, -
9. multiplication: *, /
10. unary: +, -, PRIOR, CONNECT_BY_ROOT
11. exponentiation: **
12. function call, column ref, bind variable
