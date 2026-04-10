# MySQL Recursive Descent Parser - BNF Review Skill

You are reviewing a recursive descent MySQL 8.0 parser against its BNF specification files.

**Working directory:** `/Users/rebeliceyang/Github/omni`
**Parser source:** `mysql/parser/`
**AST definitions:** `mysql/ast/`
**BNF catalog:** `mysql/parser/MYSQL_BNF_CATALOG.json`
**BNF files:** `mysql/parser/bnf/{statement}.bnf` (140 files, 1,823 lines)
**Reference:** [MySQL 8.0 SQL Statements](https://dev.mysql.com/doc/refman/8.0/en/sql-statements.html)
**Tests:** `mysql/parser/compare_test.go`

## Current State

- 97 implementation batches (0-96) are **done** (audit_round 17)
- 20 BNF-review batches (97-116) are **pending** (audit_round 18)
- 27 parser .go files, ~28K LOC parser, ~10K LOC tests
- No stubs or skipToSemicolon patterns remain
- All code compiles and tests pass

## Your Task

1. Read `mysql/parser/PROGRESS_SUMMARY.json` (lightweight version)
2. Pick the next batch to work on:
   - If any batch has `"status": "in_progress"`, **resume that batch**
   - Otherwise, find the first batch with `"status": "pending"` whose dependencies (by id) are all `"done"`
   - If any batch has `"status": "failed"`, **retry it**
3. Execute the BNF review following the steps below
4. Update `mysql/parser/PROGRESS.json` (the full file, NOT the summary):
   - Set `"in_progress"` before starting work
   - Set `"done"` only after `go build` and `go test` pass
   - Set `"failed"` with `"error"` if you cannot make tests pass

If all batches are `"done"`, output `ALL_BATCHES_COMPLETE` and stop.

## Progress Logging (MANDATORY)

You MUST print progress markers to stdout at each step:

```
[BATCH N] STARTED - batch_name
[BATCH N] STEP reading_bnf - Reading BNF files
[BATCH N] STEP comparing - Comparing parser code against BNF
[BATCH N] STEP gaps_found - Found X gaps
[BATCH N] STEP writing_tests - Writing test cases for gaps
[BATCH N] STEP fixing_code - Implementing missing branches
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

## BNF Review Steps

### Step 1: Read the BNF Files (MANDATORY)

**This is the most critical step. The BNF is the source of truth.**

1. Read the BNF files listed in the batch's `bnf_files` array from `mysql/parser/bnf/`
2. These BNF files are the **authoritative ground truth** — do NOT use WebFetch or other sources
3. Extract the **complete grammar** from each BNF file
4. Note every production rule, every alternative, every optional clause

### Step 2: Read the Parser Code

Read the parser files listed in the batch's `review_files` array:
- Understand the parse function structure
- Identify the BNF comment at the top of each function
- Map each BNF production to its corresponding code path

### Step 3: Compare Parser vs BNF (Line by Line)

For each BNF production rule:
1. Find the corresponding parse function
2. Check that **every alternative** in the BNF has a code path
3. Check that **every optional clause** is handled
4. Check that **every keyword/token** is consumed (not skipped with bare `p.advance()`)
5. Check that parsed values are **stored in AST fields** (not discarded)
6. Check that **AST Loc positions** are set (Start and End)

**Common gap patterns to look for:**
- `p.advance()` without storing the consumed value
- Missing BNF alternatives (e.g., handling ADD COLUMN but not ADD INDEX)
- Optional clauses not implemented (e.g., IF EXISTS, IF NOT EXISTS)
- Sub-clauses that reference other BNF rules but are not parsed
- Keywords consumed but not stored in AST (e.g., LOW_PRIORITY, DELAYED)

### Step 4: Write Tests for Gaps

For each gap found:
1. Write a test case in `mysql/parser/compare_test.go`
2. Use `ParseAndCheck(t, sql)` for positive tests
3. Cover each missing BNF branch with at least one test
4. Include edge cases

### Step 5: Fix Gaps

For each gap:
1. Update the parser code to handle the missing BNF branch
2. Update AST nodes in `mysql/ast/parsenodes.go` if new fields are needed
3. Update `mysql/ast/outfuncs.go` for any new AST fields/nodes
4. Set `Loc` positions on all new AST nodes

### Step 6: Build and Test

```bash
# Must compile
cd /Users/rebeliceyang/Github/omni && go build ./mysql/...

# Run batch-specific tests
cd /Users/rebeliceyang/Github/omni && go test -v -count=1 ./mysql/parser/ -run "TestXxx"

# Run full test suite (no regressions)
cd /Users/rebeliceyang/Github/omni && go test ./mysql/...
```

### Step 7: Update Progress and Commit

Edit `PROGRESS.json`:
- Set `"status"` to `"done"` after all tests pass
- **NEVER mark `"done"` if `go build ./mysql/...` or `go test ./mysql/...` fails**

Commit:
```bash
../../scripts/git-commit.sh "mysql/parser: review batch N - name" mysql/
```

## Key Files Reference

| File | Purpose |
|------|---------|
| `mysql/parser/parser.go` | Parser struct, Parse() entry point, helpers |
| `mysql/parser/lexer.go` | Lexer, Token, keyword map, token constants |
| `mysql/parser/expr.go` | Expression parser (Pratt/precedence climbing) |
| `mysql/parser/select.go` | SELECT, CTE/WITH, TABLE, VALUES, set operations |
| `mysql/parser/insert.go` | INSERT, REPLACE |
| `mysql/parser/update_delete.go` | UPDATE, DELETE |
| `mysql/parser/create_table.go` | CREATE TABLE, column defs, constraints, partitions |
| `mysql/parser/alter_table.go` | ALTER TABLE with all operations |
| `mysql/parser/create_index.go` | CREATE INDEX |
| `mysql/parser/create_view.go` | CREATE/ALTER VIEW |
| `mysql/parser/create_database.go` | CREATE/ALTER/DROP DATABASE |
| `mysql/parser/create_function.go` | CREATE FUNCTION/PROCEDURE |
| `mysql/parser/trigger.go` | CREATE TRIGGER, CREATE EVENT |
| `mysql/parser/drop.go` | DROP TABLE/INDEX/VIEW/DATABASE, TRUNCATE, RENAME |
| `mysql/parser/set_show.go` | SET, SHOW (37+ variants), USE, EXPLAIN |
| `mysql/parser/transaction.go` | BEGIN, COMMIT, ROLLBACK, SAVEPOINT, XA, LOCK |
| `mysql/parser/grant.go` | GRANT, REVOKE, CREATE/ALTER/DROP USER, roles |
| `mysql/parser/load_data.go` | LOAD DATA, LOAD XML |
| `mysql/parser/prepare.go` | PREPARE, EXECUTE, DEALLOCATE |
| `mysql/parser/utility.go` | ANALYZE, CHECK, OPTIMIZE, REPAIR, FLUSH, KILL, etc. |
| `mysql/parser/compound.go` | BEGIN...END, DECLARE, IF/CASE/WHILE/LOOP |
| `mysql/parser/signal.go` | SIGNAL, RESIGNAL, GET DIAGNOSTICS |
| `mysql/parser/replication.go` | CHANGE REPLICATION SOURCE, START/STOP REPLICA |
| `mysql/parser/stmt.go` | Top-level statement dispatch |
| `mysql/parser/name.go` | Identifiers, qualified names, variables |
| `mysql/parser/type.go` | Data type parsing |
| `mysql/ast/parsenodes.go` | AST node struct definitions |
| `mysql/ast/node.go` | Node interface, List, String, Integer |
| `mysql/ast/outfuncs.go` | AST serialization (NodeToString) |
| `mysql/parser/compare_test.go` | Test cases for parser validation |
| `mysql/parser/PROGRESS.json` | Full batch tracking |
| `mysql/parser/PROGRESS_SUMMARY.json` | Lightweight batch tracking |
| `mysql/parser/bnf/*.bnf` | BNF grammar files (140 files) |
| `mysql/parser/MYSQL_BNF_CATALOG.json` | BNF file index and metadata |

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

## Expression Precedence (Low to High)

1. `OR`, `||`
2. `XOR`
3. `AND`, `&&`
4. `NOT` (prefix)
5. `BETWEEN`, `CASE`, `WHEN`, `THEN`, `ELSE`
6. `=`, `<=>`, `>=`, `>`, `<=`, `<`, `<>`, `!=`, `IS`, `LIKE`, `REGEXP`, `IN`, `SOUNDS LIKE`, `MEMBER OF`
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

## Batch Summary

| Batch | Name | Status | Effort | BNF Lines | Primary File |
|-------|------|--------|--------|-----------|--------------|
| 0-96 | (implementation) | done | implement | - | (various) |
| 97 | review_infrastructure_shared | pending | review | - | parser.go, lexer.go, expr.go, type.go, name.go |
| 98 | review_select | pending | review | 46 | select.go |
| 99 | review_dml | pending | review | 131 | insert.go, update_delete.go |
| 100 | review_create_table | pending | review | 165 | create_table.go |
| 101 | review_alter_table | pending | review | 121 | alter_table.go |
| 102 | review_index_view_trigger | pending | review | 87 | create_index.go, create_view.go, trigger.go |
| 103 | review_routines_events | pending | review | 145 | create_function.go, utility.go |
| 104 | review_compound_flow | pending | review | 80 | compound.go, signal.go |
| 105 | review_security | pending | review | 283 | grant.go |
| 106 | review_database_tablespace_admin | pending | review | 97 | create_database.go, utility.go |
| 107 | review_transaction_xa_locking | pending | review | 84 | transaction.go |
| 108 | review_show | pending | review | 207 | set_show.go |
| 109 | review_set_explain | pending | review | 67 | set_show.go |
| 110 | review_load_import | pending | review | 86 | load_data.go, select.go |
| 111 | review_maintenance_utility | pending | review | 109 | utility.go |
| 112 | review_plugin_clone_instance | pending | review | 52 | utility.go |
| 113 | review_replication | pending | review | 14 | replication.go |
| 114 | review_spatial_resource_group | pending | review | 54 | utility.go |
| 115 | review_prepare_drop | pending | review | 50 | prepare.go, drop.go |
| 116 | review_expression_deep | pending | review | 0 | expr.go |

**Total BNF lines to review:** ~1,823

## Important Constraints

- Do NOT modify any files outside `mysql/`
- The `go.mod` at `/Users/rebeliceyang/Github/omni/go.mod` already exists
- Run `gofmt -w` on all created/modified files
- Use random sleep (1-3s) before git operations to avoid lock contention with other engine pipelines
- Every parse function MUST have the COMPLETE BNF in its comment
- Every AST node with a `Loc` field MUST have Start and End positions set
- Every new AST node type MUST have a corresponding case in `outfuncs.go`
