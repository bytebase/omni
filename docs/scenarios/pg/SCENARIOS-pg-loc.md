# PG AST Loc Field Coverage Scenarios

> Goal: Add `Loc` field to all 139 AST node structs in pg/ast that are currently missing it, and set location tracking in the parser
> Verification: Parse SQL → navigate to node → verify `sql[node.Loc.Start:node.Loc.End]` matches expected source text
> Reference sources: pg/ast/parsenodes.go (struct definitions), pg/parser/*.go (node creation sites), existing Loc patterns (Alias, RangeVar, SelectStmt)
> Proof: `go test ./pg/parser/... ./pg/parsertest/...` | same | `go test ./pg/...`

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Core Query Nodes

### 1.1 FROM Clause & Join Nodes (select.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/select.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go (shared with all sections)
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] JoinExpr — `SELECT * FROM t1 JOIN t2 ON t1.id = t2.id` → Loc covers `t1 JOIN t2 ON t1.id = t2.id`
- [x] JoinExpr LEFT — `SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id` → Loc covers full join expr
- [x] JoinExpr CROSS — `SELECT * FROM t1 CROSS JOIN t2` → Loc covers cross join expr
- [x] JoinExpr NATURAL — `SELECT * FROM t1 NATURAL JOIN t2` → Loc covers natural join expr
- [x] RangeSubselect — `SELECT * FROM (SELECT 1) AS sub` → Loc covers `(SELECT 1) AS sub`
- [x] RangeSubselect LATERAL — `SELECT * FROM t1, LATERAL (SELECT * FROM t2) AS sub` → Loc covers lateral subselect
- [x] RangeFunction — `SELECT * FROM generate_series(1, 10)` → Loc covers function call in FROM
- [x] RangeFunction WITH ORDINALITY — `SELECT * FROM generate_series(1, 10) WITH ORDINALITY` → Loc includes WITH ORDINALITY
- [x] RangeFunction ROWS FROM — `SELECT * FROM ROWS FROM(generate_series(1,3), generate_series(1,4))` → Loc covers ROWS FROM
- [x] CurrentOfExpr — `DELETE FROM t1 WHERE CURRENT OF cursor1` → Loc covers `CURRENT OF cursor1`
- [x] LockingClause — `SELECT * FROM t1 FOR UPDATE` → Loc covers `FOR UPDATE`
- [x] IntoClause — `SELECT * INTO newtable FROM t1` → Loc covers `INTO newtable`

### 1.2 Expression Helper Nodes (expr.go, name.go, update.go, insert.go, create_table.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/expr.go, pg/parser/name.go, pg/parser/select.go, pg/parser/copy.go, pg/parser/maintenance.go, pg/parser/create_view.go, pg/parser/update.go, pg/parser/insert.go, pg/parser/create_table.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`
> Note: A_Star is currently `struct{}` (empty). Adding Loc changes its memory layout — verify all creation sites still compile.

- [x] A_Indices single — `SELECT a[1] FROM t` → Loc covers `[1]`
- [x] A_Indices slice — `SELECT a[1:3] FROM t` → Loc covers `[1:3]`
- [x] A_Indirection — `SELECT (row).field FROM t` → Loc covers `.field` indirection
- [x] A_Star — `SELECT * FROM t` → Loc covers `*` (created in expr.go, name.go, select.go, copy.go, maintenance.go, create_view.go)
- [x] MergeWhenClause MATCHED — `MERGE INTO t USING s ON t.id=s.id WHEN MATCHED THEN UPDATE SET x=1` → Loc covers `WHEN MATCHED THEN UPDATE SET x=1`
- [x] MergeWhenClause NOT MATCHED — `MERGE INTO t USING s ON t.id=s.id WHEN NOT MATCHED THEN INSERT VALUES(1)` → Loc covers WHEN NOT MATCHED clause
- [x] MultiAssignRef — `INSERT INTO t(a,b) SELECT 1,2` with multi-assign → Loc covers ref
- [x] TableLikeClause — `CREATE TABLE t2 (LIKE t1 INCLUDING ALL)` → Loc covers `LIKE t1 INCLUDING ALL`

### 1.3 JSON Nodes (json.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/json.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] JsonArgument — `SELECT json_object('key': 'value')` → Loc covers argument
- [x] JsonKeyValue — `SELECT JSON_OBJECT('key' VALUE 'val')` → Loc covers key-value pair
- [x] JsonObjectAgg — `SELECT JSON_OBJECTAGG(k VALUE v) FROM t` → Loc covers aggregate
- [x] JsonArrayAgg — `SELECT JSON_ARRAYAGG(v) FROM t` → Loc covers aggregate
- [x] JsonOutput — `SELECT JSON_OBJECT('k': 'v' RETURNING text)` → Loc covers RETURNING clause
- [x] JsonValueExpr — JSON value expression in JSON functions → Loc covers value expr

---

## Phase 2: DDL Definition Nodes

### 2.1 Type & Operator Definitions (define.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/define.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] DefineStmt TYPE — `CREATE TYPE mytype AS (a int, b text)` → Loc covers full statement
- [x] DefineStmt OPERATOR — `CREATE OPERATOR === (LEFTARG=int, RIGHTARG=int, FUNCTION=int4eq)` → Loc covers full statement
- [x] DefineStmt AGGREGATE — `CREATE AGGREGATE myagg(int) (SFUNC=int4pl, STYPE=int)` → Loc covers full statement
- [x] CompositeTypeStmt — `CREATE TYPE comptype AS (x int, y int)` → Loc covers full statement
- [x] CreateEnumStmt — `CREATE TYPE mood AS ENUM ('sad', 'ok', 'happy')` → Loc covers full statement
- [x] CreateRangeStmt — `CREATE TYPE floatrange AS RANGE (SUBTYPE = float8)` → Loc covers full statement
- [x] CreateOpClassStmt — `CREATE OPERATOR CLASS myclass FOR TYPE int USING btree AS OPERATOR 1 <` → Loc covers full statement
- [x] CreateOpFamilyStmt — `CREATE OPERATOR FAMILY myfam USING btree` → Loc covers full statement
- [x] CreateOpClassItem — individual item in operator class definition → Loc covers item
- [x] CreateConversionStmt — `CREATE CONVERSION myconv FOR 'UTF8' TO 'LATIN1' FROM utf8_to_iso8859_1` → Loc covers full statement
- [x] CreateStatsStmt — `CREATE STATISTICS mystats (dependencies) ON a, b FROM t` → Loc covers full statement
- [x] AlterDefaultPrivilegesStmt — `ALTER DEFAULT PRIVILEGES GRANT SELECT ON TABLES TO public` → Loc covers full statement
- [x] AlterOpFamilyStmt — `ALTER OPERATOR FAMILY myfam USING btree ADD OPERATOR 1 <(int,int)` → Loc covers full statement
- [x] AlterOperatorStmt — `ALTER OPERATOR ===(int,int) SET (RESTRICT=eqsel)` → Loc covers full statement
- [x] AlterStatsStmt — `ALTER STATISTICS mystats SET STATISTICS 100` → Loc covers full statement
- [x] StatsElem — individual element in CREATE STATISTICS → Loc covers element

### 2.2 Extension Nodes (extension.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/extension.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreateExtensionStmt — `CREATE EXTENSION hstore` → Loc covers full statement
- [x] CreateAmStmt — `CREATE ACCESS METHOD myam TYPE INDEX HANDLER myhandler` → Loc covers full statement
- [x] CreateCastStmt — `CREATE CAST (int AS text) WITH FUNCTION int4_to_text(int)` → Loc covers full statement
- [x] CreateTransformStmt — `CREATE TRANSFORM FOR int LANGUAGE plpgsql (FROM SQL WITH FUNCTION ..., TO SQL WITH FUNCTION ...)` → Loc covers full statement
- [x] AlterExtensionStmt — `ALTER EXTENSION hstore UPDATE TO '2.0'` → Loc covers full statement
- [x] AlterExtensionContentsStmt — `ALTER EXTENSION hstore ADD TABLE mytable` → Loc covers full statement

### 2.3 Foreign Data Wrapper Nodes (fdw.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/fdw.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreateFdwStmt — `CREATE FOREIGN DATA WRAPPER myfdw` → Loc covers full statement
- [x] CreateForeignServerStmt — `CREATE SERVER myserver FOREIGN DATA WRAPPER myfdw` → Loc covers full statement
- [x] CreateForeignTableStmt — `CREATE FOREIGN TABLE ft (a int) SERVER myserver` → Loc covers full statement
- [x] CreatePLangStmt — `CREATE LANGUAGE plmylang` → Loc covers full statement
- [x] CreateUserMappingStmt — `CREATE USER MAPPING FOR current_user SERVER myserver` → Loc covers full statement
- [x] AlterFdwStmt — `ALTER FOREIGN DATA WRAPPER myfdw OPTIONS (ADD host 'foo')` → Loc covers full statement
- [x] AlterForeignServerStmt — `ALTER SERVER myserver OPTIONS (SET host 'bar')` → Loc covers full statement
- [x] AlterUserMappingStmt — `ALTER USER MAPPING FOR current_user SERVER myserver OPTIONS (SET password 'x')` → Loc covers full statement
- [x] DropUserMappingStmt — `DROP USER MAPPING FOR current_user SERVER myserver` → Loc covers full statement
- [x] ImportForeignSchemaStmt — `IMPORT FOREIGN SCHEMA public FROM SERVER myserver INTO local_schema` → Loc covers full statement

---

## Phase 3: ALTER & Access Control Nodes

### 3.1 General ALTER Nodes (alter_misc.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/alter_misc.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] AlterCollationStmt — `ALTER COLLATION "C" REFRESH VERSION` → Loc covers full statement
- [x] AlterDomainStmt — `ALTER DOMAIN mydom SET NOT NULL` → Loc covers full statement
- [x] AlterEnumStmt — `ALTER TYPE mood ADD VALUE 'great'` → Loc covers full statement
- [x] AlterEventTrigStmt — `ALTER EVENT TRIGGER mytrig DISABLE` → Loc covers full statement
- [x] AlterFunctionStmt — `ALTER FUNCTION myfunc(int) STABLE` → Loc covers full statement
- [x] AlterObjectDependsStmt — `ALTER FUNCTION myfunc DEPENDS ON EXTENSION hstore` → Loc covers full statement
- [x] AlterObjectSchemaStmt — `ALTER TABLE t SET SCHEMA newschema` → Loc covers full statement
- [x] AlterOwnerStmt — `ALTER TABLE t OWNER TO newowner` → Loc covers full statement
- [x] AlterTableSpaceOptionsStmt — `ALTER TABLESPACE myts SET (seq_page_cost=2)` → Loc covers full statement
- [x] AlterTSConfigurationStmt — `ALTER TEXT SEARCH CONFIGURATION myconfig ADD MAPPING FOR word WITH simple` → Loc covers full statement
- [x] AlterTSDictionaryStmt — `ALTER TEXT SEARCH DICTIONARY mydict (STOPWORDS = 'english')` → Loc covers full statement
- [x] AlterTypeStmt — `ALTER TYPE comptype ADD ATTRIBUTE c int` → Loc covers full statement

### 3.2 ALTER TABLE Nodes (alter_table.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/alter_table.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] AlterTableStmt — `ALTER TABLE t ADD COLUMN x int` → Loc covers full statement
- [x] AlterTableCmd — individual ALTER TABLE sub-command → Loc covers the sub-command
- [x] AlterSeqStmt — `ALTER SEQUENCE myseq RESTART WITH 1` → Loc covers full statement
- [x] AlterTableMoveAllStmt — `ALTER TABLE ALL IN TABLESPACE ts1 SET TABLESPACE ts2` → Loc covers full statement
- [x] PartitionCmd — `ALTER TABLE t ATTACH PARTITION p FOR VALUES FROM (1) TO (10)` → Loc covers partition command
- [x] RenameStmt — `ALTER TABLE t RENAME COLUMN old TO new` → Loc covers full statement

### 3.3 Grant & Role Nodes (grant.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/grant.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] GrantStmt — `GRANT SELECT ON t TO myrole` → Loc covers full statement
- [x] GrantRoleStmt — `GRANT myrole TO otheruser` → Loc covers full statement
- [x] CreateRoleStmt — `CREATE ROLE myrole WITH LOGIN` → Loc covers full statement
- [x] AlterRoleStmt — `ALTER ROLE myrole WITH SUPERUSER` → Loc covers full statement
- [x] AlterRoleSetStmt — `ALTER ROLE myrole SET search_path TO public` → Loc covers full statement
- [x] AlterPolicyStmt — `ALTER POLICY mypol ON t USING (true)` → Loc covers full statement
- [x] CreatePolicyStmt — `CREATE POLICY mypol ON t USING (true)` → Loc covers full statement
- [x] DropRoleStmt — `DROP ROLE myrole` → Loc covers full statement
- [x] AccessPriv — individual privilege in GRANT → Loc covers privilege entry

### 3.4 Publication & Subscription Nodes (publication.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/publication.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreatePublicationStmt — `CREATE PUBLICATION mypub FOR ALL TABLES` → Loc covers full statement
- [x] AlterPublicationStmt — `ALTER PUBLICATION mypub ADD TABLE t` → Loc covers full statement
- [x] CreateSubscriptionStmt — `CREATE SUBSCRIPTION mysub CONNECTION 'conninfo' PUBLICATION mypub` → Loc covers full statement
- [x] AlterSubscriptionStmt — `ALTER SUBSCRIPTION mysub SET PUBLICATION mypub` → Loc covers full statement
- [x] PublicationTable — individual table in publication → Loc covers table reference
- [x] RuleStmt — `CREATE RULE myrule AS ON INSERT TO t DO INSTEAD NOTHING` → Loc covers full statement

---

## Phase 4: Remaining Statement Nodes

### 4.1 Database & Schema Nodes (database.go, schema.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/database.go, pg/parser/schema.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreatedbStmt — `CREATE DATABASE mydb` → Loc covers full statement
- [x] AlterDatabaseStmt — `ALTER DATABASE mydb SET TABLESPACE newtbs` → Loc covers full statement
- [x] AlterDatabaseSetStmt — `ALTER DATABASE mydb SET timezone TO 'UTC'` → Loc covers full statement
- [x] DropdbStmt — `DROP DATABASE mydb` → Loc covers full statement
- [x] CreateSchemaStmt — `CREATE SCHEMA myschema` → Loc covers full statement
- [ ] CreateTableSpaceStmt — `CREATE TABLESPACE mytbs LOCATION '/data'` → Loc covers full statement
- [x] CommentStmt — `COMMENT ON TABLE t IS 'my comment'` → Loc covers full statement
- [x] SecLabelStmt — `SECURITY LABEL ON TABLE t IS 'classified'` → Loc covers full statement
- [ ] ObjectWithArgs — object reference with argument types (e.g., in DROP FUNCTION) → Loc covers object+args

### 4.2 Sequence, Function & Domain Nodes (sequence.go, create_function.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/sequence.go, pg/parser/create_function.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreateSeqStmt — `CREATE SEQUENCE myseq START 1` → Loc covers full statement
- [x] CreateDomainStmt — `CREATE DOMAIN posint AS int CHECK (VALUE > 0)` → Loc covers full statement
- [x] CreateFunctionStmt — `CREATE FUNCTION myfunc(int) RETURNS int AS $$ SELECT $1 $$ LANGUAGE sql` → Loc covers full statement
- [x] FunctionParameter — individual parameter in function definition → Loc covers parameter
- [x] ReturnStmt — `RETURN expr` in function body → Loc covers RETURN statement

### 4.3 Trigger, Index & View Nodes (trigger.go, create_index.go, create_view.go, create_table.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/trigger.go, pg/parser/create_index.go, pg/parser/create_view.go, pg/parser/create_table.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] CreateTrigStmt — `CREATE TRIGGER mytrig BEFORE INSERT ON t EXECUTE FUNCTION myfunc()` → Loc covers full statement
- [x] CreateEventTrigStmt — `CREATE EVENT TRIGGER myevt ON ddl_command_start EXECUTE FUNCTION myfunc()` → Loc covers full statement
- [x] TriggerTransition — `REFERENCING NEW TABLE AS newtab` → Loc covers transition clause
- [x] IndexStmt — `CREATE INDEX myidx ON t (col)` → Loc covers full statement
- [x] IndexElem — individual index element → Loc covers element
- [x] ViewStmt — `CREATE VIEW myview AS SELECT * FROM t` → Loc covers full statement
- [x] CreateTableAsStmt — `CREATE TABLE t AS SELECT 1` → Loc covers full statement
- [x] RefreshMatViewStmt — `REFRESH MATERIALIZED VIEW myview` → Loc covers full statement
- [x] CreateStmt — `CREATE TABLE t (a int, b text)` → Loc covers full statement

### 4.4 Utility Statement Nodes (utility.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/utility.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] ExplainStmt — `EXPLAIN SELECT 1` → Loc covers full statement
- [x] CallStmt — `CALL myproc()` → Loc covers full statement
- [x] DoStmt — `DO $$ BEGIN RAISE NOTICE 'hi'; END $$` → Loc covers full statement
- [x] CheckPointStmt — `CHECKPOINT` → Loc covers full statement
- [x] DiscardStmt — `DISCARD ALL` → Loc covers full statement
- [x] ListenStmt — `LISTEN mychannel` → Loc covers full statement
- [x] NotifyStmt — `NOTIFY mychannel, 'payload'` → Loc covers full statement
- [x] UnlistenStmt — `UNLISTEN mychannel` → Loc covers full statement
- [x] LoadStmt — `LOAD 'mylib'` → Loc covers full statement
- [x] ReassignOwnedStmt — `REASSIGN OWNED BY oldrole TO newrole` → Loc covers full statement

### 4.5 Cursor, Prepare & IO Nodes (cursor.go, prepare.go, copy.go, lock.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/cursor.go, pg/parser/prepare.go, pg/parser/copy.go, pg/parser/lock.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] DeclareCursorStmt — `DECLARE mycur CURSOR FOR SELECT 1` → Loc covers full statement
- [x] FetchStmt — `FETCH NEXT FROM mycur` → Loc covers full statement
- [x] ClosePortalStmt — `CLOSE mycur` → Loc covers full statement
- [x] PrepareStmt — `PREPARE myplan AS SELECT 1` → Loc covers full statement
- [x] ExecuteStmt — `EXECUTE myplan` → Loc covers full statement
- [x] CopyStmt — `COPY t FROM STDIN` → Loc covers full statement
- [x] LockStmt — `LOCK TABLE t IN ACCESS EXCLUSIVE MODE` → Loc covers full statement

### 4.6 Maintenance & SET Nodes (maintenance.go, set.go, drop.go)

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, pg/parser/maintenance.go, pg/parser/set.go, pg/parser/drop.go, pg/parsertest/loc_test.go
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] VacuumStmt — `VACUUM t` → Loc covers full statement
- [x] VacuumRelation — individual relation in VACUUM → Loc covers relation reference
- [x] ClusterStmt — `CLUSTER t USING myidx` → Loc covers full statement
- [x] ReindexStmt — `REINDEX TABLE t` → Loc covers full statement
- [x] VariableSetStmt — `SET search_path TO public` → Loc covers full statement
- [x] VariableShowStmt — `SHOW search_path` → Loc covers full statement
- [x] AlterSystemStmt — `ALTER SYSTEM SET max_connections = 200` → Loc covers full statement
- [x] ConstraintsSetStmt — `SET CONSTRAINTS ALL DEFERRED` → Loc covers full statement
- [x] DropStmt — `DROP TABLE t` → Loc covers full statement
- [x] DropOwnedStmt — `DROP OWNED BY myrole` → Loc covers full statement
- [x] DropSubscriptionStmt — `DROP SUBSCRIPTION mysub` → Loc covers full statement
- [x] DropTableSpaceStmt — `DROP TABLESPACE mytbs` → Loc covers full statement
- [x] TruncateStmt — `TRUNCATE t` → Loc covers full statement

---

## Phase 5: Investigation Nodes

### 5.1 Nodes Without Parser Creation Sites

> Targets: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go, investigation required
> Shared: pg/ast/parsenodes.go, pg/ast/loc.go, pg/ast/outfuncs.go
> Proof: `go test ./pg/parser/... ./pg/parsertest/... -run Loc`

- [x] FromExpr — N/A: planner-internal only; not created in parser (only in pg/catalog/query.go)
- [x] SetOperationStmt — N/A: planner-internal only; parser creates SelectStmt, analyze.go transforms to SetOperationStmt
- [x] WindowClause — N/A: planner-internal only; parser creates WindowDef, analyze.go transforms to WindowClauseQ
- [x] SortGroupClause — N/A: planner-internal only; created exclusively in pg/catalog/analyze.go
- [x] JsonReturning — N/A: struct exists but is never instantiated; parser uses JsonOutput instead
