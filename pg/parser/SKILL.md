# Recursive Descent Parser Implementation Skill

You are implementing a recursive descent PostgreSQL parser.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `pg/parser/`
**Reference grammar:** `pg/yacc/gram.y`
**AST definitions:** `pg/ast/`
**Tests:** `pg/parsertest/` and `pg/parser/compare_test.go`

## Your Task

1. Read `pg/parser/PROGRESS.json`
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

- Read the yacc grammar rules from `pg/yacc/gram.y` for the rules listed in the batch
- Read the AST node definitions from `pg/ast/parsenodes.go` (and `enums.go`, `node.go`)
- Read the existing parser code in `pg/parser/` to understand available helpers and already-implemented parse functions
- Read the test file(s) in `pg/parsertest/` corresponding to the batch's `"tests"` list

### Step 2: Fetch PG Documentation BNF (MANDATORY)

- For each major statement type, **you MUST fetch the official PostgreSQL documentation** to get the complete BNF syntax
- URL pattern: `https://www.postgresql.org/docs/17/sql-{command}.html`
  - e.g., `sql-select.html`, `sql-createtable.html`, `sql-altertable.html`
- **Every BNF branch must have a corresponding implementation** in the recursive descent parser
- **Sub-clauses must be recursively complete** — if a BNF branch references another production, that production must also be fully implemented
- **Every BNF branch must have a test case** in `compare_test.go` that exercises it
- If the documentation shows syntax alternatives (e.g., `SET | RESET`), each alternative must be handled
- Cross-reference with `pg/yacc/gram.y` to ensure no grammar alternatives are missed

### Step 3: Write Parse Functions

Create or update the target file (e.g., `pg/parser/select.go`).

**Every parse function MUST have this comment format:**

```go
// parseSelectStmt parses a SELECT statement.
//
// Ref: https://www.postgresql.org/docs/17/sql-select.html
//
//	SELECT [ ALL | DISTINCT [ ON ( expression [, ...] ) ] ]
//	    [ * | expression [ [ AS ] output_name ] [, ...] ]
//	    [ FROM from_item [, ...] ]
//	    [ WHERE condition ]
//	    [ GROUP BY [ ALL | DISTINCT ] grouping_element [, ...] ]
//	    [ HAVING condition ]
//	    ...
func (p *Parser) parseSelectStmt() *ast.SelectStmt {
```

**Rules for parse function implementation:**

1. **Use the existing AST types** from `pg/ast/` — do NOT create new node types
2. **Record positions** on EVERY AST node that has a `Loc` field.
   Set Start at the beginning: `Loc: nodes.Loc{Start: p.pos(), End: -1}`.
   Set End at the end of parsing that node: `stmt.Loc.End = p.pos()`.
   For unknown locations use `nodes.NoLoc()` which returns `Loc{Start: -1, End: -1}`.
   Check `pg/ast/parsenodes.go` — most statement and expression nodes have `Loc Loc`.
   **This is a hard requirement.** Missing locations will cause bytebase features (SQL rewriting,
   backup/restore, error reporting) to break. The comparison tests verify locations are non-negative.
3. **Token constants** are in `pg/parser/tokens.go` (exported from yacc) and `pg/yacc/parser.go`
4. **Lexer types** (`Token`, `Lexer`, `NewLexer`) are in `pg/parser/lexer.go`
5. **Keyword token types** are in `pg/parser/keywords.go`
6. **Match the yacc grammar semantics exactly** — the recursive descent parser must produce identical AST to the yacc parser for all valid inputs
7. **Error recovery**: When encountering unexpected tokens, try to recover by:
   - Skipping to the next semicolon for statement-level errors
   - Returning a partial AST node with what was parsed so far
   - Recording errors but continuing to parse subsequent statements
8. **Operator precedence** in expressions: use Pratt parsing (precedence climbing) for `a_expr` and `b_expr`

### Step 4: Wire Up

- If implementing `stmt` dispatch (batch 34), wire all statement parsers into the `stmt` rule
- For `stmtmulti`, handle the `;`-separated list and wrap in `*ast.RawStmt` with position info
- Ensure `Parse()` in `parser.go` calls the top-level `parseStmtBlock()` or equivalent

### Step 5: Test

Run these commands in order:

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./pg/parser/

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./pg/parsertest/ -run "TestXxx|TestYyy"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./...
```

**Test strategy:**
- `pg/parser/compare_test.go` contains `CompareWithYacc(t, sql)` which parses with both
  parsers and compares AST via `ast.NodeToString()`. This is the primary validation.
- For each batch, **add test cases to `compare_test.go`** using SQL strings relevant to
  the batch's grammar rules. Add a new `TestCompare{BatchName}` function.
- Remove the `t.Skip()` from `TestCompareBasic` once the first statements can be parsed.
- The tests in `pg/parsertest/` use `yacc.Parse()` — do NOT modify them.
- If a comparison test fails, fix the recursive descent implementation, never the test.

### Step 6: Update Progress

Edit `PROGRESS.json`:
- **Before starting work**: set `"status"` to `"in_progress"` so a restart knows this batch was interrupted
- **After all tests pass**: set `"status"` to `"done"`
- **If tests fail and you cannot fix**: set `"status"` to `"failed"` and add `"error": "description"`
- **NEVER mark `"done"` if `go build ./pg/parser/` or `go test ./pg/parser/` fails**

### Step 7: Git Commit

```bash
../../scripts/git-commit.sh "pg/parser: implement batch N - name" pg/
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `pg/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `pg/parser/lexer.go` | Lexer (Token, Lexer, NewLexer) |
| `pg/parser/keywords.go` | Keyword constants and lookup |
| `pg/parser/tokens.go` | Token type constants (from yacc) |
| `pg/yacc/gram.y` | Reference grammar (17865 lines) |
| `pg/yacc/parser.go` | Generated yacc parser (reference + parse tables) |
| `pg/ast/parsenodes.go` | AST node struct definitions |
| `pg/ast/enums.go` | Enum constants (SortByDir, SetOperation, etc.) |
| `pg/ast/node.go` | Node interface, List, String, Integer, etc. |
| `pg/parsertest/*.go` | 746 test cases organized by feature |
| `pg/pgregress/` | PostgreSQL regression test compatibility |

## Token Type Mapping

The `Parser.advance()` method already maps lexer-internal token types (`lex_*` constants)
to parser token constants (from `tokens.go`). You do NOT need to handle this mapping.

When writing parse functions, use the constants from `tokens.go` and `keywords.go`:
- Single-char tokens: `'('`, `')'`, `','`, `';'`, etc. — use the byte value directly (e.g., `p.expect('(')`)
- Keywords: use constants like `SELECT`, `FROM`, `WHERE`, `CREATE`, etc.
- Literals: use `ICONST`, `FCONST`, `SCONST`, `IDENT`, etc.
- Operators: use `TYPECAST`, `NOT_EQUALS`, `LESS_EQUALS`, `Op`, etc.

The `Token.Str` field contains the string value for identifiers, string literals, and operators.
The `Token.Ival` field contains the integer value for `ICONST`.

## Completeness Requirements

The parser must achieve **100% coverage** of the PostgreSQL yacc grammar (`pg/yacc/gram.y`, 694 rules). The following requirements apply:

1. **Every statement type in the `stmt` rule must be dispatched** in `parseStmt()`
2. **Every BNF branch of each grammar rule must have a corresponding code path** in the recursive descent parser
3. **Sub-clauses must be recursively complete** — a partially-implemented sub-rule (e.g., returning nil for some branches) is a gap
4. **Every BNF branch must have a test case** in `compare_test.go` that exercises it and validates AST parity with the yacc parser
5. **Official PG documentation must be fetched** (URL: `https://www.postgresql.org/docs/17/sql-{command}.html`) for every statement being implemented, to ensure all documented syntax variants are covered

## Audit Gap Summary (2026-03-12)

### Missing Statement Types (not dispatched in parseStmt)
- `CallStmt` (gram.y line 13449) — CALL procedure_name(args)
- `AlterSystemStmt` (gram.y line 6709) — ALTER SYSTEM SET/RESET
- `CreateTableSpaceStmt` (gram.y line 14698) — CREATE TABLESPACE
- `DropTableSpaceStmt` (gram.y line 14719) — DROP TABLESPACE
- `CreateConversionStmt` (gram.y line 16578) — CREATE [DEFAULT] CONVERSION

### Partially Implemented (dispatched but incomplete)
- `AlterTblSpcStmt` (gram.y line 14736) — ALTER TABLESPACE SET/RESET reloptions returns nil
- `RemoveOperStmt` (gram.y line 13653) — DROP OPERATOR has TODO placeholder, returns nil
- `RemoveFuncStmt` (gram.y line 13563) — DROP FUNCTION/PROCEDURE/ROUTINE uses parseAnyNameList instead of function_with_argtypes_list (missing argument type parsing)
- `RemoveAggrStmt` (gram.y line 13626) — DROP AGGREGATE uses parseAnyNameList instead of aggregate_with_argtypes_list

### Coverage Statistics
- Total gram.y rules: 694
- Rules tracked in PROGRESS.json: 455 (batches 0-34, all done)
- Pending batches: 35-41 (7 new batches covering gaps)

## Expression Parsing Strategy

For `a_expr` (the most complex rule), use **Pratt parsing / precedence climbing**:

```go
func (p *Parser) parseAExpr(minPrec int) nodes.Node {
    left := p.parseCExpr() // atoms: literals, column refs, subqueries, etc.
    for {
        prec, ok := p.infixPrecedence(p.cur.Type)
        if !ok || prec < minPrec {
            break
        }
        op := p.advance()
        right := p.parseAExpr(prec + 1) // +1 for left-assoc, +0 for right-assoc
        left = &nodes.BoolExpr{...} // or A_Expr, depending on operator
    }
    return left
}
```

Precedence levels (from low to high, matching PostgreSQL):
1. OR
2. AND
3. NOT (prefix)
4. IS, ISNULL, NOTNULL, IS NULL, IS NOT NULL, IS TRUE, etc.
5. comparison: =, <, >, <=, >=, <>, !=
6. LIKE, ILIKE, SIMILAR TO, BETWEEN, IN
7. addition: +, -
8. multiplication: *, /, %
9. exponentiation: ^
10. unary: +, -, ~
11. subscript: [], typecast: ::
12. function call, column ref
