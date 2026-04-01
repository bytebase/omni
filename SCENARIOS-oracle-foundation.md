# Oracle Parser Foundation Alignment

> Goal: Align the Oracle parser's Loc tracking and error handling infrastructure with the PG parser. This is the prerequisite for all future Oracle parser ecosystem work (catalog, semantic, completion, migration).
> Verification: `go test ./oracle/...` passes. All existing tests continue to pass. New tests validate each alignment change.
> Reference: PG parser at `pg/parser/` and `pg/ast/` — the target behavior for each scenario.

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Loc Semantics & Utilities

### 1.1 Unify Loc Unknown Convention

Current Oracle: `End=0` means unknown/unset, `-1` also means unknown (dual semantics).
Target: both Start and End use `-1` exclusively for "unknown", matching PG.

- [ ] `NoLoc()` function added to `oracle/ast/node.go` returning `Loc{Start: -1, End: -1}`
- [ ] Loc struct doc comment updated: both Start and End use -1 for unknown (remove 0=unknown for End)
- [ ] All `Loc{Start: someVal}` initializations (where End defaults to 0) changed to explicitly set End=-1 or set End properly
- [ ] All existing code checking `End == 0` to mean "unknown" updated to check `End == -1`
- [ ] All existing code checking `Start == -1` remains correct (no change needed — verify only)
- [ ] `go test ./oracle/...` passes after Loc semantic unification

### 1.2 Token End Field

Current Oracle Token: `{Type, Str, Ival, Loc}` — no End field.
Target: add End field matching PG Token: `{Type, Str, Ival, Loc, End}`.

- [ ] Token struct gains `End int` field in `oracle/parser/lexer.go`
- [ ] Lexer.NextToken() sets `Token.End` to byte offset past end of token for every token type
- [ ] EOF token has `End == Loc` (zero-width)
- [ ] Identifier tokens: End = Loc + len(identifier)
- [ ] Keyword tokens: End = Loc + len(keyword)
- [ ] String literal tokens: End = Loc + full literal length (including quotes)
- [ ] Number literal tokens: End = Loc + len(number text)
- [ ] Operator tokens (single char): End = Loc + 1
- [ ] Multi-char operator tokens (||, <=, >=, !=, <>): End = Loc + len(operator)
- [ ] Q-quote string tokens: End includes closing quote delimiter
- [ ] National character literal tokens (N'...'): End includes full N'...' span
- [ ] Hint tokens (/*+ ... */): End includes closing */
- [ ] Label tokens (<<...>>): End includes closing >>
- [ ] Bind variable tokens (:name): End = Loc + 1 + len(name)
- [ ] `go test ./oracle/parser/...` — all existing lexer tests pass with Token.End added

### 1.3 Parser Source Field & prev.End

Current Oracle Parser struct: `{lexer, cur, prev, nextBuf, hasNext}` — no source field.
Target: add `source string` field, use `prev.End` for Loc.End tracking.

- [ ] Parser struct gains `source string` field
- [ ] Parse() function stores sql string in `p.source`
- [ ] `p.prev.End` is available after `advance()` for Loc.End tracking (enabled by 1.2)
- [ ] `go test ./oracle/...` passes — no behavior change yet

### 1.4 RawStmt Conversion

Current Oracle RawStmt: `{Stmt StmtNode, StmtLocation int, StmtLen int}`.
Target: `{Stmt StmtNode, Loc Loc}` matching PG.

- [ ] RawStmt struct changed to `{Stmt StmtNode, Loc Loc}` in `oracle/ast/parsenodes.go`
- [ ] RawStmt gains `stmtNode()` method (already a StmtNode — verify)
- [ ] Parse() in `oracle/parser/parser.go` creates RawStmt with `Loc: Loc{Start: stmtLoc, End: p.prev.End}`
- [ ] `writeRawStmt` in `oracle/ast/outfuncs.go` updated to use `n.Loc.Start` and `n.Loc.End`
- [ ] All test code referencing `StmtLocation` or `StmtLen` updated to use `Loc`
- [ ] `go test ./oracle/...` passes after RawStmt conversion

### 1.5 NodeLoc Utility

Current: no way to extract Loc from an arbitrary Node.
Target: `NodeLoc(n Node) Loc` function covering all 248 Loc-bearing struct types.

- [ ] `oracle/ast/loc.go` file created
- [ ] `NodeLoc(n Node) Loc` function implemented with type switch
- [ ] Covers all statement nodes (SelectStmt, InsertStmt, UpdateStmt, DeleteStmt, MergeStmt, DropStmt, CreateTableStmt, AlterTableStmt, etc.)
- [ ] Covers all expression nodes (BinaryExpr, UnaryExpr, BoolExpr, FuncCallExpr, CaseExpr, CastExpr, etc.)
- [ ] Covers all table expression nodes (TableRef, SubqueryRef, LateralRef, JoinClause, etc.)
- [ ] Covers all clause/definition nodes (ColumnDef, TableConstraint, ColumnConstraint, PartitionClause, etc.)
- [ ] Covers all PL/SQL nodes (PLSQLBlock, PLSQLIf, PLSQLLoop, PLSQLAssign, etc.)
- [ ] Covers RawStmt (returns Loc field)
- [ ] Returns `NoLoc()` for nil input
- [ ] Returns `NoLoc()` for unknown/uncovered types (default case)
- [ ] Test: `NodeLoc(nil)` returns `NoLoc()`
- [ ] Test: `NodeLoc(&SelectStmt{Loc: Loc{Start: 0, End: 42}})` returns `Loc{0, 42}`
- [ ] Test: `NodeLoc(&String{})` returns `NoLoc()` (no Loc field)

### 1.6 ListSpan Utility

Current: no way to compute byte range of a list of nodes.
Target: `ListSpan(list *List) Loc` matching PG behavior.

- [ ] `ListSpan(list *List) Loc` function added to `oracle/ast/loc.go`
- [ ] Returns `NoLoc()` for nil list
- [ ] Returns `NoLoc()` for empty list (len=0)
- [ ] Returns `NoLoc()` if first item has `Start == -1`
- [ ] Returns `NoLoc()` if last item has `End == -1`
- [ ] Returns `Loc{first.Start, last.End}` for valid list
- [ ] Test: `ListSpan(nil)` returns `NoLoc()`
- [ ] Test: `ListSpan(&List{})` returns `NoLoc()`
- [ ] Test: single-item list returns that item's Loc
- [ ] Test: multi-item list returns `{first.Start, last.End}`
- [ ] Test: list with unknown Loc items returns `NoLoc()`

---

## Phase 2: ParseError Enhancement

### 2.1 ParseError Struct Alignment

Current Oracle ParseError: `{Message string, Position int}`.
Target: `{Severity string, Code string, Message string, Position int}` matching PG.

- [ ] ParseError gains `Severity string` field
- [ ] ParseError gains `Code string` field
- [ ] `Error()` method returns `"ERROR: <message> (SQLSTATE 42601)"` format matching PG
- [ ] `Error()` defaults Severity to "ERROR" when empty
- [ ] `Error()` defaults Code to "42601" when empty
- [ ] Existing error creation sites still compile (Message and Position are unchanged)
- [ ] `go test ./oracle/...` passes

### 2.2 Error Helper Methods

Current: no `syntaxErrorAtCur`, `syntaxErrorAtTok`, `tokenText`, `lexerError`.
Target: full set of error helpers matching PG parser.

- [ ] `tokenText(tok Token) string` method added to Parser — extracts source text for token using `p.source[tok.Loc:tok.End]`
- [ ] `tokenText` returns empty string for EOF token
- [ ] `tokenText` falls back to `tok.Str` if Loc/End out of bounds
- [ ] `tokenText` returns `string(rune(tok.Type))` for single-char tokens as last fallback
- [ ] `syntaxErrorAtTok(tok Token) *ParseError` — returns `"syntax error at or near \"<text>\""` with Position
- [ ] `syntaxErrorAtTok` returns `"syntax error at end of input"` for empty token text
- [ ] `syntaxErrorAtCur() *ParseError` — delegates to `syntaxErrorAtTok(p.cur)`
- [ ] `lexerError() *ParseError` — returns `"<lexer error> at or near \"<text>\""` with Position
- [ ] `lexerError` returns plain lexer error message when token text is empty
- [ ] Test: `syntaxErrorAtCur()` on EOF produces "syntax error at end of input"
- [ ] Test: `syntaxErrorAtCur()` on keyword token produces `"syntax error at or near \"SELECT\""`
- [ ] Test: `syntaxErrorAtTok()` with specific token produces correct message and position
- [ ] `go test ./oracle/...` passes

### 2.3 expect() Error Message Alignment

Current: `expect()` returns `"expected token type %d, got %d (%q)"` — raw token IDs, not human-readable.
Target: use `syntaxErrorAtCur()` for PG-compatible error messages.

- [ ] `expect()` uses `syntaxErrorAtCur()` instead of manually constructing ParseError
- [ ] Error message now reads `"syntax error at or near \"<token>\""` instead of token type numbers
- [ ] Error Position is set from current token's Loc
- [ ] All existing callers of `expect()` continue to work (signature unchanged: `(Token, error)`)
- [ ] `go test ./oracle/...` passes

### 2.4 Parse() Top-Level Error Alignment

Current Parse(): creates `ParseError{Message: "unexpected token...", Position: ...}` inline.
Target: use `syntaxErrorAtCur()` and add lexer error checking matching PG.

- [ ] Parse() checks `p.lexer.Err != nil` before each statement parse, returns `p.lexerError()`
- [ ] Parse() uses `syntaxErrorAtCur()` when `parseStmt()` returns nil
- [ ] Parse() checks `p.lexer.Err != nil` after each statement parse, returns `p.lexerError()`
- [ ] Parse() creates RawStmt with `Loc{Start: stmtLoc, End: p.prev.End}` (depends on Phase 1)
- [ ] Error messages now use PG format: `"syntax error at or near \"<token>\""` instead of `"unexpected token %q at position %d"`
- [ ] `go test ./oracle/...` passes

---

## Phase 3: Loc.End Consistency Across Parser Files

### 3.0 Loc.End Pattern Migration

Current Oracle pattern: `node.Loc.End = p.pos()` (uses current token start = next token's position).
Target PG pattern: `node.Loc.End = p.prev.End` (uses previous token's end = precise end of last consumed token).

After Phase 1.2 adds Token.End, the Loc.End pattern must change across all parser files.

- [ ] All existing `node.Loc.End = p.pos()` assignments migrated to `node.Loc.End = p.prev.End`
- [ ] All `Loc{Start: start, End: p.pos()}` changed to `Loc{Start: start, End: p.prev.End}`
- [ ] No remaining uses of `p.pos()` for Loc.End (grep verification: zero matches for `Loc.End = p.pos()` and `End: p.pos()}`)
- [ ] `go test ./oracle/...` passes after pattern migration

### 3.1 select.go — Loc.End Audit

File: `oracle/parser/select.go` (2727 lines, 45 parse functions, 50 Loc.End assignments).
Target: every parse function that creates an AST node with Loc sets Loc.End = p.prev.End before return.

- [ ] Audit all 45 parse functions — identify any missing Loc.End assignments
- [ ] All SelectStmt nodes have Loc.End set
- [ ] All CTE/WithClause nodes have Loc.End set
- [ ] All JoinClause nodes have Loc.End set
- [ ] All ResTarget nodes have Loc.End set
- [ ] All SortBy nodes have Loc.End set
- [ ] All WindowSpec/WindowDef nodes have Loc.End set
- [ ] All subquery-related nodes have Loc.End set
- [ ] Pattern: `node.Loc.End = p.prev.End` (using Token.End from Phase 1.2)
- [ ] `go test ./oracle/parser/...` passes

### 3.2 expr.go — Loc.End Audit

File: `oracle/parser/expr.go` (1617 lines, 37 parse functions, 18 Loc.End assignments).
Target: close the gap — 37 functions but only 18 End assignments.

- [ ] Audit all 37 parse functions — identify missing Loc.End assignments
- [ ] All BinaryExpr/UnaryExpr/BoolExpr nodes have Loc.End set
- [ ] All FuncCallExpr nodes have Loc.End set
- [ ] All CaseExpr/CaseWhen nodes have Loc.End set
- [ ] All CastExpr/TreatExpr nodes have Loc.End set
- [ ] All BetweenExpr/InExpr/LikeExpr/IsExpr nodes have Loc.End set
- [ ] All SubqueryExpr/ExistsExpr/ParenExpr nodes have Loc.End set
- [ ] All IntervalExpr nodes have Loc.End set
- [ ] All BindVariable/PseudoColumn nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.3 create_table.go — Loc.End Audit

File: `oracle/parser/create_table.go` (1867 lines, 20 parse functions, 8 Loc.End assignments).
Target: 20 functions but only 8 End assignments — significant gap.

- [ ] Audit all 20 parse functions
- [ ] CreateTableStmt has Loc.End set
- [ ] All ColumnDef nodes have Loc.End set
- [ ] All ColumnConstraint nodes have Loc.End set
- [ ] All TableConstraint nodes have Loc.End set
- [ ] All PartitionClause/PartitionDef nodes have Loc.End set
- [ ] All StorageClause nodes have Loc.End set
- [ ] All IdentityClause nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.4 alter_table.go — Loc.End Audit

File: `oracle/parser/alter_table.go` (2080 lines, 29 parse functions, 1 Loc.End assignment).
Target: 29 functions but only 1 End assignment — largest gap.

- [ ] Audit all 29 parse functions
- [ ] AlterTableStmt has Loc.End set
- [ ] All AlterTableCmd nodes have Loc.End set
- [ ] All inline ColumnDef/Constraint nodes have Loc.End set
- [ ] All partition-related ALTER nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.5 alter_misc.go — Loc.End Audit

File: `oracle/parser/alter_misc.go` (3567 lines, 32 parse functions, 15 Loc.End assignments).
Target: close gap from 15 to 32.

- [ ] Audit all 32 parse functions
- [ ] AlterIndexStmt has Loc.End set
- [ ] AlterViewStmt has Loc.End set
- [ ] AlterSequenceStmt has Loc.End set
- [ ] AlterProcedureStmt/AlterFunctionStmt/AlterPackageStmt have Loc.End set
- [ ] AlterTriggerStmt/AlterTypeStmt have Loc.End set
- [ ] AlterSynonymStmt/AlterDatabaseLinkStmt have Loc.End set
- [ ] All remaining ALTER statement nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.6 create_admin.go — Loc.End Audit

File: `oracle/parser/create_admin.go` (6437 lines, 81 parse functions, 73 Loc.End assignments).
Target: close small gap from 73 to 81.

- [ ] Audit all 81 parse functions — identify 8 missing Loc.End assignments
- [ ] All CREATE TABLESPACE/CLUSTER/DIMENSION nodes have Loc.End set
- [ ] All CREATE MATERIALIZED VIEW/LOG nodes have Loc.End set
- [ ] All analytic view/attribute dimension/hierarchy nodes have Loc.End set
- [ ] All remaining admin DDL nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.7 database.go — Loc.End Audit

File: `oracle/parser/database.go` (6480 lines, 34 parse functions, 10 Loc.End assignments).
Target: 34 functions but only 10 End assignments — large gap.

- [ ] Audit all 34 parse functions
- [ ] All CREATE/ALTER DATABASE nodes have Loc.End set
- [ ] All CREATE/ALTER PLUGGABLE DATABASE nodes have Loc.End set
- [ ] All CREATE/ALTER DISKGROUP nodes have Loc.End set
- [ ] All remaining database-level nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.8 security.go — Loc.End Audit

File: `oracle/parser/security.go` (1520 lines, 23 parse functions, 20 Loc.End assignments).
Target: close small gap from 20 to 23.

- [ ] Audit all 23 parse functions
- [ ] All CreateUserStmt/AlterUserStmt nodes have Loc.End set
- [ ] All CreateRoleStmt/AlterRoleStmt nodes have Loc.End set
- [ ] All CreateProfileStmt/AlterProfileStmt nodes have Loc.End set
- [ ] All audit policy nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.9 plsql_block.go — Loc.End Audit

File: `oracle/parser/plsql_block.go` (1238 lines, 31 parse functions, 25 Loc.End assignments).
Target: close gap from 25 to 31.

- [ ] Audit all 31 parse functions
- [ ] PLSQLBlock has Loc.End set
- [ ] All PLSQLIf/PLSQLElsIf nodes have Loc.End set
- [ ] All PLSQLLoop/PLSQLForall nodes have Loc.End set
- [ ] All PLSQLAssign/PLSQLReturn/PLSQLRaise nodes have Loc.End set
- [ ] All PLSQLExecImmediate/PLSQLOpen/PLSQLFetch/PLSQLClose nodes have Loc.End set
- [ ] All PLSQLVarDecl/PLSQLCursorDecl nodes have Loc.End set
- [ ] `go test ./oracle/parser/...` passes

### 3.10 Remaining Files — Loc.End Audit

Files: `create_index.go` (15 fn, 10 End), `create_proc.go` (10 fn, 6 End), `create_view.go` (8 fn, 7 End), `create_trigger.go` (2 fn, 2 End), `create_type.go` (6 fn, 4 End), `insert.go` (7 fn, 5 End), `update_delete.go` (5 fn, 4 End), `merge.go` (5 fn, 3 End), `drop.go` (3 fn, 4 End), `grant.go` (6 fn, 2 End), `name.go` (5 fn, 7 End), `type.go` (5 fn, 1 End), `utility.go` (18 fn, 18 End), `transaction.go` (4 fn, 4 End), `session.go` (2 fn, 4 End), `sequence_synonym.go` (4 fn, 3 End), `comment.go` (1 fn, 1 End).

- [ ] create_index.go: audit 15 functions, close gap from 10 to 15
- [ ] create_proc.go: audit 10 functions, close gap from 6 to 10
- [ ] create_view.go: audit 8 functions, close gap from 7 to 8
- [ ] create_type.go: audit 6 functions, close gap from 4 to 6
- [ ] insert.go: audit 7 functions, close gap from 5 to 7
- [ ] update_delete.go: audit 5 functions, close gap from 4 to 5
- [ ] merge.go: audit 5 functions, close gap from 3 to 5
- [ ] grant.go: audit 6 functions, close gap from 2 to 6
- [ ] type.go: audit 5 functions, close gap from 1 to 5
- [ ] sequence_synonym.go: audit 4 functions, close gap from 3 to 4
- [ ] create_trigger.go, drop.go, utility.go, transaction.go, session.go, comment.go: verify already complete
- [ ] `go test ./oracle/...` passes — all files audited

---

## Phase 4: Verification & Integration Tests

### 4.1 Loc Accuracy Integration Tests

Verify Loc tracking produces correct byte ranges for representative SQL across all statement categories.

- [ ] Test: `SELECT * FROM t` — SelectStmt.Loc spans entire statement
- [ ] Test: `SELECT a, b FROM t WHERE x = 1 ORDER BY a` — all sub-nodes have correct Loc ranges
- [ ] Test: `INSERT INTO t (a) VALUES (1)` — InsertStmt.Loc spans entire statement
- [ ] Test: `UPDATE t SET a = 1 WHERE b = 2` — UpdateStmt.Loc spans entire statement
- [ ] Test: `DELETE FROM t WHERE a = 1` — DeleteStmt.Loc spans entire statement
- [ ] Test: `CREATE TABLE t (id NUMBER, name VARCHAR2(100))` — all ColumnDef Loc ranges correct
- [ ] Test: `ALTER TABLE t ADD (col NUMBER)` — AlterTableStmt and AlterTableCmd Loc correct
- [ ] Test: `CREATE OR REPLACE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;` — function Loc spans full
- [ ] Test: `DECLARE x NUMBER; BEGIN x := 1; END;` — PLSQLBlock and child node Loc correct
- [ ] Test: Multi-statement input — each RawStmt.Loc covers only its statement, not adjacent ones
- [ ] Test: RawStmt.Loc.End for last statement in multi-statement input is correct
- [ ] Test: NodeLoc extracts correct Loc for each statement type in a multi-statement parse

### 4.2 ParseError Format Tests

Verify error messages match PG format.

- [ ] Test: parse `SELEC` produces `"ERROR: syntax error at or near \"SELEC\" (SQLSTATE 42601)"`
- [ ] Test: parse empty string produces no error (empty list)
- [ ] Test: parse `SELECT` (incomplete) produces `"syntax error at end of input"` in message
- [ ] Test: parse `SELECT FROM` produces error at "FROM" with correct Position
- [ ] Test: ParseError.Severity defaults to "ERROR"
- [ ] Test: ParseError.Code defaults to "42601"
- [ ] Test: lexer error on unterminated string produces error with position

### 4.3 Backward Compatibility

Verify all existing functionality still works after foundation changes.

- [ ] All existing oracle/parser/ test files pass without modification (except Loc field renames)
- [ ] All existing oracle/ast/ test files pass
- [ ] outfuncs.go serialization produces valid output for all node types
- [ ] Parse → outfuncs round-trip produces equivalent SQL for all existing test cases
- [ ] No external packages importing oracle/ast break (check with `go build ./...`)
