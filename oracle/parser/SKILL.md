# Recursive Descent Parser Implementation Skill (Oracle PL/SQL)

You are implementing a recursive descent Oracle PL/SQL parser using a BNF-first methodology.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `oracle/parser/`
**AST definitions:** `oracle/ast/`
**Tests:** `oracle/parser/compare_test.go`

## Reference Documentation

- **BNF definitions**: `oracle/parser/bnf/{STATEMENT-NAME}.bnf` — 171 statements, extracted from Oracle 26 official docs
- These BNF files are the **authoritative ground truth** — do NOT use WebFetch or other sources
- Also reference the ANTLR2-based open-source Oracle grammar for understanding edge cases

## Your Task

1. Read `oracle/parser/PROGRESS_SUMMARY.json` (lightweight version — done batches are compressed, only pending/failed/in_progress have full detail)
2. Pick the next batch to work on:
   - If any batch has `"status": "in_progress"`, **resume that batch** (it was interrupted mid-work — read the existing code in its target file and continue from where it left off)
   - Otherwise, find the first batch with `"status": "pending"` whose dependencies (by id) are all `"done"`
   - If any batch has `"status": "failed"`, **retry it** (reset to `"in_progress"` and try again)
3. Implement that batch following the **BNF-First Workflow** below
4. Update `oracle/parser/PROGRESS.json` (the full file, NOT the summary):
   - Set `"in_progress"` before starting work
   - Set `"done"` only after `go build` and `go test` pass
   - Set `"failed"` with `"error"` if you cannot make tests pass

If all batches are `"done"`, output `ALL_BATCHES_COMPLETE` and stop.

## Progress Logging (MANDATORY)

You MUST print progress markers to stdout at each step:

```
[BATCH N] STARTED - batch_name
[BATCH N] STEP reading_bnf - Reading BNF files: file1.bnf, file2.bnf
[BATCH N] STEP reviewing_code - Reviewing existing code in target files
[BATCH N] STEP writing_tests - Writing test cases (N tests for M BNF lines)
[BATCH N] STEP writing_code - Implementing/fixing parse functions
[BATCH N] STEP build - Running go build
[BATCH N] STEP test - Running go test (X passed, Y failed)
[BATCH N] STEP commit - Committing changes
[BATCH N] DONE
```

If a step fails:
```
[BATCH N] FAIL test - description of failure
[BATCH N] RETRY - what you're fixing
```

**Do NOT skip these markers.**

## BNF-First Workflow (Per-Batch)

### Step 1: Read ALL BNF Files (MANDATORY — DO NOT SKIP)

**This is the most critical step. Every batch starts here.**

For every BNF file listed in the batch's `bnf_files` array:

1. **Read the BNF file**: `oracle/parser/bnf/{STATEMENT-NAME}.bnf`
2. **Read ALL referenced sub-clause BNFs**: If the BNF references `storage_clause`, `physical_attributes_clause`, etc., check if they have their own `.bnf` file and read those too
3. **Count branches**: Every `|` alternative, every `[ optional ]` clause, every `{ required choice }` — these ALL need code paths
4. **Note the BNF line count**: This determines minimum test count (see Quality Requirements)

### Step 2: Review Existing Code (MANDATORY)

Read the target `.go` files and identify:

1. **Stubs**: `skipToSemicolon` calls — these MUST be replaced with proper parsing
2. **Stubs**: `parseAdminDDLStmt` dispatches — these MUST be replaced with dedicated parsers for any statement that has a `.bnf` file
3. **Missing branches**: Compare BNF branches against code `switch`/`if` branches. Every BNF alternative must have a code path
4. **Incorrect BNF comments**: Compare existing code comments against the `.bnf` file content. Fix any divergence
5. **Missing `Loc` recording**: Every AST node must have `Loc.Start` and `Loc.End` set with the single sentinel convention below

### Step 3: Write Tests FIRST (Test-Driven Development)

**TEST-DRIVEN**: Write tests BEFORE implementing.

**Minimum test count**: `max(3, bnf_lines / 15)` test cases per batch.

Rules:
- **Every branch in the BNF must have at least one test case**
- If ALTER TABLE has ADD, MODIFY, DROP COLUMN, DROP CONSTRAINT, RENAME, ENABLE, DISABLE, MOVE, SPLIT PARTITION, EXCHANGE PARTITION — you need at least 10 test cases
- Add test cases to `compare_test.go` using a new `TestParse{BatchName}` function
- For "review" effort batches: add any **missing** tests, even if the code is already correct
- Test SQL should be realistic Oracle SQL, not minimal stubs

### Step 4: Implement / Fix Parse Functions

Create or update the target file(s).

**BNF Comment Requirement (HARD REQUIREMENT):**

Every parse function MUST have the COMPLETE BNF from the `.bnf` file in its doc comment:

```go
// parseCreateTableStmt parses a CREATE TABLE statement.
//
// BNF: oracle/parser/bnf/CREATE-TABLE.bnf
//
//  CREATE [ GLOBAL TEMPORARY | PRIVATE TEMPORARY | SHARDED | DUPLICATED ]
//      TABLE [ schema. ] table_name
//      [ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
//      { relational_table | object_table | XMLType_table }
//      [ MEMOPTIMIZE FOR READ ] [ MEMOPTIMIZE FOR WRITE ]
//      [ PARENT [ schema. ] table ] ;
//
//  relational_table:
//      [ ( relational_properties ) ]
//      ...
func (p *Parser) parseCreateTableStmt() *ast.CreateTableStmt {
```

**The BNF comment must match the `.bnf` file exactly. No abbreviation, no `...` for omitted branches.**

**Implementation Rules:**

1. **Every BNF branch → code path**: If the BNF says `{ ADD | MODIFY | DROP | RENAME | ENABLE | DISABLE | MOVE | SPLIT | MERGE | EXCHANGE }`, handle ALL of them. If a branch requires a new AST node type or field, add it to `parsenodes.go` and `outfuncs.go`.
2. **Sub-clauses must be recursively complete**: If `column_definition` has a full BNF with DEFAULT, IDENTITY, ENCRYPT, inline_constraint — implement it completely.
3. **NO `skipToSemicolon` in new code**: If a clause is truly optional/rare, parse it as a generic option list, not skip.
4. **NO `parseAdminDDLStmt` for any statement that has a `.bnf` file**: Replace with a dedicated parser.
5. **Record `Loc` on EVERY AST node**: Set `Loc: nodes.Loc{Start: p.pos(), End: -1}` at start, then set `node.Loc.End = p.prev.End` after the parser has consumed the node's last token. Use `-1` as the only unknown/end-not-yet-set sentinel. Do not use `0` as an unknown `End`.
6. **Serialization**: Add serialization in `oracle/ast/outfuncs.go` for every new node type. Every node with `nodeTag()` must have a `writeNode` case and `writeXxx` function.
7. **Use existing AST types** from `oracle/ast/` — do NOT create new node types unless absolutely necessary.
8. **Token constants** are keyword constants (kw*) and lexer tokens (tok*) in `oracle/parser/lexer.go`.
9. **Match Oracle grammar semantics exactly**.
10. **Error propagation**: Required parse functions must return syntax and lexer errors directly. Do not skip to semicolon or return partial success for required grammar positions. Optional probes must be explicitly named `tryParse...` or `parseOptional...` and may return `(nil, nil)` when the optional clause is absent.
11. **Operator precedence**: Use Pratt parsing (precedence climbing) for expressions.

**Return type conventions:**
- Required expression parsers return `(ExprNode, error)` (for example, `parseExpr() (ast.ExprNode, error)`)
- Required table reference parsers return `(TableExpr, error)` (for example, `parseTableRef() (ast.TableExpr, error)`)
- Required statement parsers return `(*ast.XStmt, error)` or `(ast.StmtNode, error)`
- Optional probes are named `tryParse...` or `parseOptional...` and may return `(nil, nil)` when the grammar branch is absent
- Callers must propagate nested errors; assigning a parse-call error result to `_` is forbidden

### Step 5: Build + Test

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./oracle/...

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./oracle/parser/ -run "TestXxx"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./oracle/...
```

**NEVER mark `"done"` if `go build` or `go test` fails.**

### Step 6: Update Progress

Edit `PROGRESS.json`:
- **Before starting work**: set `"status"` to `"in_progress"`
- **After all tests pass**: set `"status"` to `"done"`
- **If tests fail and you cannot fix**: set `"status"` to `"failed"` and add `"error": "description"`

### Step 7: Git Commit

```bash
../../scripts/git-commit.sh "oracle/parser: implement batch N - name" oracle/
```

## Quality Requirements

### Code Quality Rules
- Every parse function MUST have the complete BNF from the `.bnf` file in its doc comment
- Every branch in the BNF comment MUST have a corresponding code path (no silent drops)
- NO `skipToSemicolon` in new code — parse properly or use a generic option list
- NO `parseAdminDDLStmt` for any statement that has a `.bnf` file
- Record `Loc` (Start/End) on every AST node using `p.prev.End` and the `-1` sentinel for unknown ends
- Add serialization in `outfuncs.go` for every new node type
- Test count requirement: minimum `max(3, bnf_lines / 15)` test cases per batch

### Review Criteria (for "review" effort batches)
Even review-only batches MUST:
- Verify every BNF branch has a code path
- Verify every code path has a test
- Add missing tests
- Fix any divergence between code comments and actual BNF
- Update BNF comments in code to match the `.bnf` file exactly

### Batch Effort Types
- **review**: Read BNF, compare against existing code, add missing tests, fix BNF comment divergence
- **review+fix**: Review + fix missing branches, replace stubs with proper parsing
- **fix+implement**: Fix existing partial implementations + implement new statements from scratch
- **implement**: Full implementation from scratch (statement was previously stubbed via `parseAdminDDLStmt` or `skipToSemicolon`)
- **implement+review**: Implement new + review existing related statements

## Key Files Reference

| File | Purpose |
|------|---------|
| `oracle/parser/bnf/*.bnf` | **BNF syntax definitions** (171 files, primary reference) |
| `oracle/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `oracle/parser/lexer.go` | Lexer (Token, Lexer, NewLexer), keywords |
| `oracle/ast/node.go` | Node interface, List, String, Integer |
| `oracle/ast/parsenodes.go` | AST node struct definitions, enums |
| `oracle/ast/outfuncs.go` | AST serialization (NodeToString) |
| `oracle/parser/compare_test.go` | Test infrastructure |
| `oracle/parser/PROGRESS.json` | Batch progress tracking |
| `oracle/parser/ORACLE_BNF_CATALOG.json` | BNF extraction catalog (statement list + status) |

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
- FLASHBACK queries (AS OF SCN/TIMESTAMP, VERSIONS BETWEEN)
- MODEL clause (spreadsheet-like calculations)
- SAMPLE clause (random row sampling)
- Analytic functions with OVER, KEEP (DENSE_RANK FIRST/LAST), WITHIN GROUP
- GROUPING SETS / CUBE / ROLLUP extensions to GROUP BY
- LATERAL inline views
- XMLTABLE / JSON_TABLE as table sources
- XML expressions (XMLELEMENT, XMLFOREST, XMLAGG)
- JSON expressions (JSON_OBJECT, JSON_ARRAY, JSON_VALUE, JSON_EXISTS, IS JSON)
- Cursor expressions: CURSOR(subquery)
- MULTISET operations (MULTISET UNION/INTERSECT/EXCEPT)
- TREAT expression for object type casting
- FORALL / BULK COLLECT / PIPE ROW in PL/SQL
- PRAGMA directives (AUTONOMOUS_TRANSACTION, EXCEPTION_INIT, etc.)
- PL/SQL CASE statement (distinct from CASE expression)
- PL/SQL EXIT [WHEN] / CONTINUE [WHEN]
- PL/SQL collection types in DECLARE (TABLE OF, VARRAY, RECORD, REF CURSOR)
- Named parameter notation (param => value)
- LOCK TABLE, CALL, RENAME statements
- SET ROLE, SET CONSTRAINT(S) session control
- AUDIT / NOAUDIT
- ASSOCIATE / DISASSOCIATE STATISTICS
- Compound triggers (multiple timing sections)
- DDL and database event triggers

### Expression Parsing Strategy

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

## Batch Summary

| ID | Name | BNF Lines | Effort | Deps |
|----|------|-----------|--------|------|
| 87 | Infrastructure & Shared Sub-rules | — | review | — |
| 88 | SELECT | 397 | review | 87 |
| 89 | INSERT/UPDATE/DELETE/MERGE | 232 | review | 87 |
| 90 | CREATE TABLE | 774 | review+fix | 87 |
| 91 | ALTER TABLE | 1,010 | review+fix | 90 |
| 92 | Index + Indextype + Operator | 459 | fix+implement | 87 |
| 93 | Views (all types) | 645 | fix+implement | 87 |
| 94 | Database + Controlfile | 445 | review+fix | 87 |
| 95 | Pluggable Database | 372 | implement | 94 |
| 96 | Diskgroup | 235 | implement | 87 |
| 97 | Tablespace (all) | 293 | review+fix | 87 |
| 98 | ALTER SESSION + SYSTEM | 135 | implement | 87 |
| 99 | User / Role / Profile | 207 | fix+review | 87 |
| 100 | Audit (all) | 180 | implement+review | 87 |
| 101 | PL/SQL Objects (all) | 202 | review+fix | 87 |
| 102 | Attr Dimension / Hierarchy / Domain | 217 | implement | 87 |
| 103 | Cluster / Dimension / Zonemap / InMem | 228 | review+implement | 87 |
| 104 | Graph / Vector / Lockdown / Outline | 188 | implement | 87 |
| 105 | Small Objects Bundle | 239 | implement+review | 87 |
| 106 | Review Bundle (Txn+Utility+DCL+Drop) | 880 | review | 87 |

## Verification

After each batch:
```bash
cd /Users/rebeliceyang/Github/omni
go build ./oracle/...
go test -v ./oracle/...
```

After all batches complete:
- Run `/audit-parser` skill to verify completeness
- Cross-check: every `.bnf` file has a code path and a test
