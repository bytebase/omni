# MSSQL Parser Error Alignment Scenarios

> Goal: Align the MSSQL parser with PG parser improvements — dual return (Node, error), soft-fail error detection, and error message quality
> Verification: `go build ./mssql/parser/ && go test ./mssql/parser/ -count=1`; new error_test.go table-driven tests per scenario
> Reference sources: PG parser (pg/parser/) as pattern reference; truncated SQL must return error, not silent OK

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Dual Return Migration

Migrate all parse functions from returning single values to `(T, error)` tuples. Start with leaf files, work up through callers. Each scenario = one file's functions migrated + all cross-file callers updated + builds + tests pass.

### 1.1 Leaf: type.go, name.go

- [x] type.go — `parseDataType()` returns `(*nodes.DataType, error)`, all callers updated
- [x] name.go — `parseIdentifier()` already returns `(string, bool)`, keep; `parseTableRef()` returns `(*nodes.TableRef, error)`; `parseIdentExpr()` returns `(nodes.ExprNode, error)`; `parseQualifiedRef()` returns `(nodes.ExprNode, error)`; `parseFuncCallWithSchema()` returns `(nodes.ExprNode, error)` — all callers updated
- [x] Build passes after 1.1

### 1.2 Expression: expr.go

- [x] Binary operators — `parseExpr`, `parseOr`, `parseAnd`, `parseNot`, `parseComparison`, `parseAddition`, `parseMultiplication`, `parseUnary` return `(nodes.ExprNode, error)`
- [x] Primaries — `parsePrimary`, `parseCast`, `parseTryCast`, `parseConvert`, `parseTryConvert` return `(nodes.ExprNode, error)`
- [x] Special expressions — `parseCaseExpr`, `parseCoalesce`, `parseNullif`, `parseIif`, `parseExists` return `(nodes.ExprNode, error)`
- [x] Function/window — `parseFuncCall`, `parseOverClause`, `parseWindowFrame`, `parseWindowFrameBound`, `parseWithinGroupClause` return dual
- [x] Collate/ATZ — `parseCollateExpr`, `parseAtTimeZoneExpr` return `(nodes.ExprNode, error)`
- [x] All cross-file callers of expr.go functions updated (select.go, insert.go, update_delete.go, merge.go, create_table.go, alter_table.go, control_flow.go, declare_set.go, etc.)
- [x] Build passes after 1.2

### 1.3 DML: select.go

- [x] Top-level — `parseSelectStmt`, `parseSetOperation`, `parseWithClause`, `parseCTE` return dual
- [x] Clauses — `parseTopClause`, `parseTargetList`, `parseFromClause`, `parseForClause` return dual
- [x] Table sources — `parseTableSource`, `parsePrimaryTableSource`, `parsePivotUnpivot`, `parsePivotExpr`, `parseUnpivotExpr`, `parseTableSampleClause`, `parseTableValuedFunction` return dual
- [x] Lists — `parseExprList`, `parseGroupByList`, `parseGroupingSet`, `parseWindowClause`, `parseOrderByList` return dual
- [x] Hints/options — `parseTableHints`, `parseTableHint`, `parseIndexValue`, `parseOptionClause`, `parseQueryHint`, `parseOptimizeForParam` return dual
- [x] All cross-file callers updated (parser.go, insert.go, merge.go, create_view.go, cursor.go, etc.)
- [x] Build passes after 1.3

### 1.4 DML: insert.go, update_delete.go, merge.go

- [x] insert.go — `parseInsertStmt`, `parseValuesClause`, `parseOutputClause` return dual
- [x] update_delete.go — `parseUpdateStmt`, `parseDeleteStmt`, `parseWhereClauseBody`, `parseSetClauseList`, `parseSetClause`, `parseSetTarget` return dual
- [x] merge.go — `parseMergeStmt`, `parseMergeWhenClause`, `parseMergeInsertAction` return dual
- [x] All cross-file callers updated
- [x] Build passes after 1.4

### 1.5 DDL: create_table.go, create_index.go

- [x] create_table.go — `parseCreateTableStmt`, `parseColumnDef`, `parseTableConstraint`, `parseInlineConstraint`, `parseIdentitySpec`, `parseEncryptedWith`, `parseGeneratedAlways`, `parseTableOptions`, `parseOneTableOption`, `parseEdgeConstraintConnections`, `parseInlineTableIndex`, `parseCreateTableAsSelectStmt`, `parseCTASWithOptions`, `parseCTASOption`, `parseParenIdentList`, `parseColumnConstraint` return dual (note: `parseReferentialActions` and `parseRefAction` kept as-is since they are void/value-type helpers not returning nodes)
- [x] create_index.go — `parseCreateIndexStmt`, `parseIndexColumnList`, `parseCreateXmlIndexStmt`, `parseCreateSelectiveXmlIndexStmt`, `parseCreateSpatialIndexStmt`, `parseCreateAggregateStmt`, `parseDropAggregateStmt`, `parseCreateJsonIndexStmt`, `parseCreateVectorIndexStmt`, `parseOptionList` return dual
- [x] All cross-file callers updated
- [x] Build passes after 1.5

### 1.6 DDL: create_view.go, create_proc.go, create_trigger.go

- [x] create_view.go — `parseCreateViewStmt`, `parseCreateMaterializedViewStmt`, `parseAlterMaterializedViewStmt` return dual
- [x] create_proc.go — `parseCreateProcedureStmt`, `parseCreateFunctionStmt`, `parseParamDefList`, `parseParamDef`, `parseRoutineOptionList`, `parseRoutineOption`, `parseMethodSpecifier` return dual
- [x] create_trigger.go — `parseCreateTriggerStmt`, `parseTriggerWithOptions` return dual
- [x] All cross-file callers updated
- [x] Build passes after 1.6

### 1.7 DDL: create_database.go, create_type.go, create_sequence.go, create_schema.go, create_statistics.go, create_synonym.go

- [ ] create_database.go — all 9 functions return dual
- [ ] create_type.go — `parseCreateTypeStmt`, `parseTableTypeIndex` return dual
- [ ] create_sequence.go — `parseCreateSequenceStmt`, `parseAlterSequenceStmt` return dual
- [ ] create_schema.go — `parseCreateSchemaStmt`, `parseAlterSchemaStmt` return dual
- [ ] create_statistics.go — all 4 functions return dual
- [ ] create_synonym.go — `parseCreateSynonymStmt` returns dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.7

### 1.8 ALTER: alter_table.go, alter_objects.go

- [ ] alter_table.go — `parseAlterTableStmt`, `parseAlterTableAdd`, `parseAlterTableDrop`, `parseAlterTableAlterColumn`, `parseAlterColumnAddDrop`, `parseAlterTableCheckConstraint`, `parseAlterTableEnableDisable`, `parseAlterTableEnableDisableTrigger`, `parseAlterTableChangeTracking`, `parseAlterTableSwitch`, `parseAlterTableRebuild`, `parseAlterTableSplitMergeRange`, `parseAlterTableSet`, `parseAlterTableFiletableNamespace`, `parseAlterTableAddPeriod`, `parseKeyValueOptionList` return dual
- [ ] alter_objects.go — `parseAlterDatabaseStmt`, `parseAlterIndexStmt`, `parseAlterIndexOptions` return dual; void helpers update to return error
- [ ] All cross-file callers updated
- [ ] Build passes after 1.8

### 1.9 Security: grant.go, security_principals.go, security_misc.go, security_keys.go, security_audit.go

- [ ] grant.go — `parseGrantStmt`, `parseRevokeStmt`, `parseDenyStmt`, `parsePrivilegeList`, `parsePrincipalList` return dual
- [ ] security_principals.go — all 9 functions return dual
- [ ] security_misc.go — all 10 functions return dual
- [ ] security_keys.go — all 14 functions return dual
- [ ] security_audit.go — all 12 functions return dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.9

### 1.10 Control Flow: control_flow.go, cursor.go, transaction.go, declare_set.go, execute.go

- [ ] control_flow.go — `parseIfStmt`, `parseWhileStmt`, `parseBeginStmt`, `parseTryCatchStmt`, `parseReturnStmt`, `parseBreakStmt`, `parseContinueStmt`, `parseGotoStmt`, `parseWaitForStmt` return dual
- [ ] cursor.go — all 6 functions return dual
- [ ] transaction.go — all 5 functions return dual
- [ ] declare_set.go — `parseDeclareStmt`, `parseVariableDecl`, `parseSetStmt`, `parseSetOptionStmt` return dual
- [ ] execute.go — all 7 functions return dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.10

### 1.11 Server/Admin: server.go, backup_restore.go

- [ ] server.go — all 16 functions return dual
- [ ] backup_restore.go — all 11 functions return dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.11

### 1.12 Server/Admin: service_broker.go, resource_governor.go, availability.go, endpoint.go

- [ ] service_broker.go — all 26 functions return dual
- [ ] resource_governor.go — all 19 functions return dual
- [ ] availability.go — all 6 functions return dual
- [ ] endpoint.go — all 11 functions return dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.12

### 1.13 Other: drop.go, utility.go, fulltext.go, external.go, assembly.go, partition.go, xml_schema.go, dbcc.go, bulk_insert.go, go_batch.go, rowset_functions.go, event.go

- [ ] drop.go — `parseDropStmt` returns dual
- [ ] utility.go — all 34 functions return dual
- [ ] fulltext.go — all 11 functions return dual
- [ ] external.go — all 21 functions return dual
- [ ] assembly.go — 2 functions return dual
- [ ] partition.go — 4 functions return dual
- [ ] xml_schema.go — 2 functions return dual
- [ ] dbcc.go — `parseDbccStmt` returns dual
- [ ] bulk_insert.go — 2 functions return dual
- [ ] go_batch.go — `parseGoStmt` returns dual
- [ ] rowset_functions.go — `parseRowsetFunction`, `parseRowsetWithClause` return dual
- [ ] event.go — all 13 functions return dual
- [ ] All cross-file callers updated
- [ ] Build passes after 1.13

### 1.14 Top-level: parser.go

- [ ] `parseStmt` returns `(nodes.StmtNode, error)` — main dispatch updated
- [ ] `parseCreateStmt` returns `(nodes.StmtNode, error)`
- [ ] `parseAlterStmt` returns `(nodes.StmtNode, error)`
- [ ] `parseDropOrSecurityStmt` returns `(nodes.StmtNode, error)`
- [ ] `parseAddStmt` returns `(nodes.StmtNode, error)`
- [ ] `parseLabelStmt` returns dual
- [ ] `Parse()` top-level loop updated to propagate errors from `parseStmt`
- [ ] All existing tests pass
- [ ] Build passes after 1.14

---

## Phase 2: Soft-Fail Error Alignment

After dual return is in place, add nil checks after advance+parse patterns. Each scenario = one truncated SQL that must return an error instead of silent OK. Tested via error_test.go table-driven tests.

### 2.1 expr.go — Binary Operators

- [ ] `SELECT 1 OR` — OR consumed, no right operand → error
- [ ] `SELECT 1 AND` — AND consumed, no right operand → error
- [ ] `SELECT 1 <` — less-than consumed, no right operand → error
- [ ] `SELECT 1 >` — greater-than consumed, no right operand → error
- [ ] `SELECT 1 =` — equals consumed, no right operand → error
- [ ] `SELECT 1 <=` — LESS_EQUALS consumed, no right → error
- [ ] `SELECT 1 >=` — GREATER_EQUALS consumed, no right → error
- [ ] `SELECT 1 <>` — NOT_EQUALS consumed, no right → error
- [ ] `SELECT 1 +` — plus consumed, no right → error
- [ ] `SELECT 1 -` — minus consumed, no right → error
- [ ] `SELECT 1 *` — multiply consumed, no right → error
- [ ] `SELECT 1 /` — divide consumed, no right → error
- [ ] `SELECT 1 %` — modulo consumed, no right → error

### 2.2 expr.go — Pattern Matching & Special Expressions

- [ ] `SELECT 1 BETWEEN` — BETWEEN consumed, no lower bound → error
- [ ] `SELECT 1 BETWEEN 0 AND` — lower bound parsed, AND consumed, no upper → error
- [ ] `SELECT 1 NOT BETWEEN` — NOT BETWEEN consumed, no lower → error
- [ ] `SELECT 'a' LIKE` — LIKE consumed, no pattern → error
- [ ] `SELECT 'a' NOT LIKE` — NOT LIKE consumed, no pattern → error
- [ ] `SELECT 1 IN (` — IN ( consumed, no expression list → error
- [ ] `SELECT 1 NOT IN (` — NOT IN ( consumed, no expression list → error
- [ ] `SELECT NOT` — NOT consumed, no operand → error

### 2.3 expr.go — CAST, CONVERT & Function Calls

- [ ] `SELECT CAST(` — CAST( consumed, no expression → error
- [ ] `SELECT CAST(1 AS` — AS consumed, no type → error
- [ ] `SELECT TRY_CAST(` — TRY_CAST( consumed, no expression → error
- [ ] `SELECT CONVERT(` — CONVERT( consumed, no type → error
- [ ] `SELECT TRY_CONVERT(` — TRY_CONVERT( consumed, no type → error
- [ ] `SELECT COALESCE(` — COALESCE( consumed, no expression → error
- [ ] `SELECT NULLIF(` — NULLIF( consumed, no expression → error
- [ ] `SELECT IIF(` — IIF( consumed, no condition → error
- [ ] `SELECT CASE WHEN` — CASE WHEN consumed, no condition → error
- [ ] `SELECT CASE WHEN 1=1 THEN` — THEN consumed, no result → error
- [ ] `SELECT EXISTS (` — EXISTS ( consumed, no subquery → error

### 2.4 expr.go — Collate, AT TIME ZONE, Unary

- [ ] `SELECT 'a' COLLATE` — COLLATE consumed, no collation → error
- [ ] `SELECT GETDATE() AT TIME ZONE` — AT TIME ZONE consumed, no zone → error
- [ ] `SELECT -` — unary minus, no operand → error
- [ ] `SELECT +` — unary plus, no operand → error
- [ ] `SELECT ~` — bitwise NOT, no operand → error

### 2.5 select.go — FROM, JOIN, WHERE, GROUP BY

- [ ] `SELECT * FROM` — FROM consumed, no table source → error
- [ ] `SELECT * FROM t JOIN` — JOIN consumed, no right table → error
- [ ] `SELECT * FROM t LEFT JOIN` — LEFT JOIN consumed, no right table → error
- [ ] `SELECT * FROM t CROSS JOIN` — CROSS JOIN consumed, no right table → error
- [ ] `SELECT * FROM t JOIN t2 ON` — ON consumed, no join condition → error
- [ ] `SELECT * FROM t WHERE` — WHERE consumed, no condition → error
- [ ] `SELECT * FROM t GROUP BY` — GROUP BY consumed, no grouping expressions → error
- [ ] `SELECT * FROM t HAVING` — HAVING consumed, no condition → error
- [ ] `SELECT * FROM t ORDER BY` — ORDER BY consumed, no sort expressions → error
- [ ] `SELECT TOP` — TOP consumed, no count → error
- [ ] `SELECT TOP (` — TOP ( consumed, no expression → error

### 2.6 select.go — Subquery, CTE, UNION

- [ ] `SELECT * FROM (` — subquery FROM consumed, no content → error
- [ ] `WITH cte AS (` — CTE AS consumed, no query → error
- [ ] `WITH cte AS (SELECT 1) SELECT 1 UNION` — UNION consumed, no right query → error
- [ ] `WITH cte AS (SELECT 1) SELECT 1 EXCEPT` — EXCEPT consumed, no right query → error
- [ ] `WITH cte AS (SELECT 1) SELECT 1 INTERSECT` — INTERSECT consumed, no right query → error
- [ ] `SELECT 1 UNION ALL` — UNION ALL consumed, no right query → error

### 2.7 DML: INSERT, UPDATE, DELETE, MERGE Truncations

- [ ] `INSERT INTO` — INTO consumed, no table name → error
- [ ] `INSERT INTO t VALUES (` — VALUES ( consumed, no values → error
- [ ] `INSERT INTO t (` — column list ( consumed, no columns → error
- [ ] `UPDATE` — UPDATE with no table → error
- [ ] `UPDATE t SET` — SET consumed, no assignments → error
- [ ] `UPDATE t SET a =` — = consumed, no value → error
- [ ] `DELETE FROM` — DELETE FROM consumed, no table → error
- [ ] `MERGE INTO` — MERGE INTO consumed, no target table → error
- [ ] `MERGE INTO t USING` — USING consumed, no source → error
- [ ] `MERGE INTO t USING s ON` — ON consumed, no condition → error

### 2.8 DDL: CREATE TABLE Truncations

- [ ] `CREATE TABLE` — CREATE TABLE with no name → error
- [ ] `CREATE TABLE t (` — open paren, no columns → error
- [ ] `CREATE TABLE t (a` — column name, no type → error
- [ ] `CREATE TABLE t (a INT,` — trailing comma, no next column → error
- [ ] `CREATE TABLE t (a INT CONSTRAINT` — CONSTRAINT consumed, no name → error
- [ ] `CREATE TABLE t (a INT PRIMARY` — PRIMARY consumed, no KEY → error
- [ ] `CREATE TABLE t (a INT REFERENCES` — REFERENCES consumed, no table → error
- [ ] `CREATE TABLE t (a INT DEFAULT` — DEFAULT consumed, no value → error
- [ ] `CREATE TABLE t (a INT CHECK (` — CHECK ( consumed, no condition → error
- [ ] `CREATE TABLE t (CONSTRAINT pk PRIMARY KEY (` — PK ( consumed, no columns → error

### 2.9 DDL: ALTER TABLE, CREATE INDEX Truncations

- [ ] `ALTER TABLE` — ALTER TABLE with no name → error
- [ ] `ALTER TABLE t ADD` — ADD consumed, no column/constraint → error
- [ ] `ALTER TABLE t DROP COLUMN` — DROP COLUMN consumed, no name → error
- [ ] `ALTER TABLE t ALTER COLUMN` — ALTER COLUMN consumed, no name → error
- [ ] `ALTER TABLE t ALTER COLUMN a` — column name, no type → error
- [ ] `CREATE INDEX` — CREATE INDEX with no name → error
- [ ] `CREATE INDEX ix ON` — ON consumed, no table → error
- [ ] `CREATE INDEX ix ON t (` — ( consumed, no columns → error
- [ ] `CREATE UNIQUE INDEX` — no name → error

### 2.10 Control Flow & Declaration Truncations

- [ ] `IF` — IF with no condition → error
- [ ] `IF 1=1` — no BEGIN or statement body → error (or accepted as valid depending on T-SQL rules)
- [ ] `WHILE` — WHILE with no condition → error
- [ ] `DECLARE @v` — variable declared, no type → error
- [ ] `SET @v =` — = consumed, no value → error
- [ ] `RETURN` — valid (RETURN with no value is legal in T-SQL procedures)
- [ ] `EXEC` — EXEC with no procedure name → error
- [ ] `BEGIN TRANSACTION` — valid (no name required)
- [ ] `PRINT` — PRINT with no expression → error
- [ ] `THROW` — valid (re-throw with no args is legal in T-SQL CATCH blocks)

---

## Phase 3: Error Message Enhancement

### 3.1 Parser Infrastructure

- [ ] Add `source string` field to Parser struct, populated in `Parse()`
- [ ] Add `lexerError()` method that wraps `p.lexer.Err` with position
- [ ] Add lexer error check in main parse loop (before and after each statement)
- [ ] `syntaxErrorAtCur()` helper returns `ParseError` with "at or near" context
- [ ] `syntaxErrorAtEnd()` helper returns `ParseError` for end-of-input
- [ ] Build passes after 3.1

### 3.2 ParseError Enhancement

- [ ] `ParseError` struct gains `Source string` field for original SQL
- [ ] `ParseError.Error()` includes "at or near" text extracted from source
- [ ] `ParseError.Error()` distinguishes end-of-input ("at end of input") vs mid-input ("at or near 'X'")
- [ ] Position information is accurate for multi-line SQL
- [ ] `expect()` uses enhanced `ParseError` with context
- [ ] Error messages for truncated `SELECT 1 OR` say "at end of input"
- [ ] Error messages for invalid token mid-statement say "at or near" with the token text
- [ ] All existing tests still pass with new error format
- [ ] Build passes after 3.2
