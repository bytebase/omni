package parser

import (
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// The Loc contract (AGENTS.md): every node's Loc is the byte range of its
// source text, from the start of its first token to the end of its last,
// exclusive, never covering the whitespace or comments around it. A
// location that is not known is ast.NoLoc().

// locSpan is a node that carries a Loc, found by walking a parse result.
type locSpan struct {
	node   string  // the node's Go type name
	loc    ast.Loc // its location
	parent ast.Loc // the location of the nearest enclosing node that has one
}

// locSpans walks nodes depth first in field order and returns every node
// that carries a Loc.
func locSpans(nodes ...ast.Node) []locSpan {
	var spans []locSpan
	locType := reflect.TypeOf(ast.Loc{})
	var walk func(v reflect.Value, parent ast.Loc)
	walk = func(v reflect.Value, parent ast.Loc) {
		for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			if v.Type() == locType {
				return
			}
			if f := v.FieldByName("Loc"); f.IsValid() && f.Type() == locType {
				loc := f.Interface().(ast.Loc)
				spans = append(spans, locSpan{node: v.Type().Name(), loc: loc, parent: parent})
				if loc != ast.NoLoc() {
					parent = loc
				}
			}
			for i := 0; i < v.NumField(); i++ {
				if v.Type().Field(i).IsExported() {
					walk(v.Field(i), parent)
				}
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				walk(v.Index(i), parent)
			}
		}
	}
	for _, n := range nodes {
		walk(reflect.ValueOf(n), ast.NoLoc())
	}
	return spans
}

// locSpanCases pins, for every node kind that carries a Loc, the text its
// Loc covers. node names the Go type; nth picks among the nodes of that
// type in walk order.
var locSpanCases = []struct {
	node string
	sql  string
	want string
	nth  int
}{
	{"AlterDatabaseStmt", "ALTER DATABASE mydb READ ONLY = 1", "ALTER DATABASE mydb READ ONLY = 1", 0},
	{"AlterEventStmt", "ALTER EVENT ev1 ENABLE", "ALTER EVENT ev1 ENABLE", 0},
	{"AlterInstanceStmt", "ALTER INSTANCE RELOAD TLS", "ALTER INSTANCE RELOAD TLS", 0},
	{"AlterLogfileGroupStmt", "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' ENGINE = NDB", "ALTER LOGFILE GROUP lg1 ADD UNDOFILE 'undo2.dat' ENGINE = NDB", 0},
	{"AlterResourceGroupStmt", "ALTER RESOURCE GROUP rg1 VCPU = 0-3", "ALTER RESOURCE GROUP rg1 VCPU = 0-3", 0},
	{"AlterRoutineStmt", "ALTER FUNCTION f1 COMMENT 'test'", "ALTER FUNCTION f1 COMMENT 'test'", 0},
	{"AlterServerStmt", "ALTER SERVER s1 OPTIONS (USER 'admin')", "ALTER SERVER s1 OPTIONS (USER 'admin')", 0},
	{"AlterTableCmd", "ALTER TABLE t FORCE", "FORCE", 0},
	{"AlterTableStmt", "ALTER TABLE t ADD COLUMN c INT, DROP COLUMN d", "ALTER TABLE t ADD COLUMN c INT, DROP COLUMN d", 0},
	{"AlterTablespaceStmt", "ALTER TABLESPACE ts1 RENAME TO ts2", "ALTER TABLESPACE ts1 RENAME TO ts2", 0},
	{"AlterUserStmt", "ALTER USER USER() 2 FACTOR UNREGISTER", "ALTER USER USER() 2 FACTOR UNREGISTER", 0},
	{"AlterViewStmt", "ALTER VIEW v AS SELECT 1", "ALTER VIEW v AS SELECT 1", 0},
	{"AnalyzeTableStmt", "ANALYZE TABLE t1", "ANALYZE TABLE t1", 0},
	{"Assignment", "SET x = 1", "x = 1", 0},
	{"Assignment", "UPDATE t SET a = 1 WHERE id = 2", "a = 1", 0},
	{"Assignment", "INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1\n", "x = 1", 0},
	{"BeginEndBlock", "CREATE PROCEDURE p() BEGIN END", "BEGIN END", 0},
	{"BeginStmt", "BEGIN", "BEGIN", 0},
	{"BetweenExpr", "SELECT a BETWEEN 1 AND 10 FROM t", "a BETWEEN 1 AND 10", 0},
	{"BinaryExpr", "SELECT * FROM t WHERE n = NULL ORDER BY n", "n = NULL", 0},
	{"BinaryExpr", "DO 1 + 2", "1 + 2", 0},
	{"BinaryExpr", "SELECT +a = 1 FROM t", "+a = 1", 0},
	{"BinaryExpr", "SELECT (a) = 1 FROM t", "(a) = 1", 0},
	{"BinaryExpr", "SELECT ((SELECT 1) = 0)", "(SELECT 1) = 0", 0},
	{"BinaryExpr", "SELECT * FROM t WHERE a = 1 AND b = 2 ", "a = 1 AND b = 2", 0},
	{"BinaryExpr", "SELECT * FROM t WHERE a = 1 AND b = 2 ", "a = 1", 1},
	{"BinaryExpr", "SELECT a REGEXP 'x' FROM t", "a REGEXP 'x'", 0},
	{"BinaryExpr", "SELECT a SOUNDS LIKE b FROM t", "a SOUNDS LIKE b", 0},
	{"BinaryExpr", "SELECT a = ANY (SELECT 1) FROM t", "a = ANY (SELECT 1)", 0},
	{"BinlogStmt", "BINLOG 'base64'", "BINLOG 'base64'", 0},
	{"BitLit", "SELECT b'01000001' 'c'", "b'01000001'", 0},
	{"BoolLit", "UPDATE t SET active = TRUE WHERE id = 1", "TRUE", 0},
	{"CacheIndexStmt", "CACHE INDEX t1 IN hot_cache", "CACHE INDEX t1 IN hot_cache", 0},
	{"CallStmt", "CALL my_proc(1, 'a')", "CALL my_proc(1, 'a')", 0},
	{"CaseExpr", "SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t", "CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END", 0},
	{"CaseStmtNode", "CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SELECT 1; END CASE; END", "CASE x WHEN 1 THEN SELECT 1; END CASE", 0},
	{"CaseStmtWhen", "CREATE PROCEDURE p() BEGIN CASE x WHEN 1 THEN SELECT 1; END CASE; END", "WHEN 1 THEN SELECT 1;", 0},
	{"CaseWhen", "SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t", "WHEN a > 0 THEN 'pos'", 0},
	{"CastExpr", "SELECT CAST(a AS CHAR) FROM t", "CAST(a AS CHAR)", 0},
	{"ChangeReplicationFilterStmt", "CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1)", "CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1)", 0},
	{"ChangeReplicationSourceStmt", "CHANGE MASTER TO MASTER_SSL = 0", "CHANGE MASTER TO MASTER_SSL = 0", 0},
	{"CheckTableStmt", "CHECK TABLE t1", "CHECK TABLE t1", 0},
	{"ChecksumTableStmt", "CHECKSUM TABLE t1", "CHECKSUM TABLE t1", 0},
	{"CloneStmt", "CLONE LOCAL DATA DIRECTORY '/tmp/clone'", "CLONE LOCAL DATA DIRECTORY '/tmp/clone'", 0},
	{"CloseCursorStmt", "CREATE PROCEDURE p_close_undecl() BEGIN CLOSE nope; END", "CLOSE nope", 0},
	{"CollateExpr", "SELECT a COLLATE utf8mb4_bin FROM t", "a COLLATE utf8mb4_bin", 0},
	{"CollateExpr", "SELECT ((SELECT 'a') COLLATE utf8mb4_bin = 'a')", "(SELECT 'a') COLLATE utf8mb4_bin", 0},
	{"ColumnConstraint", "CREATE TABLE t (id INT NULL)", "NULL", 0},
	{"ColumnDef", "CREATE TABLE t (id INT)", "id INT", 0},
	{"ColumnRef", "SET x = 1", "x", 0},
	{"CommitStmt", "COMMIT", "COMMIT", 0},
	{"CommonTableExpr", "WITH cte AS (SELECT 1) SELECT * FROM cte", "cte AS (SELECT 1)", 0},
	{"Constraint", "ALTER TABLE t ADD INDEX idx_c (c)", "INDEX idx_c (c)", 0},
	{"ConvertExpr", "SELECT CONVERT(a, CHAR) FROM t", "CONVERT(a, CHAR)", 0},
	{"CreateDatabaseStmt", "CREATE SCHEMA mydb", "CREATE SCHEMA mydb", 0},
	{"CreateEventStmt", "CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO SELECT 1", "CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO SELECT 1", 0},
	{"CreateFunctionStmt", "CREATE PROCEDURE p() BEGIN END", "CREATE PROCEDURE p() BEGIN END", 0},
	{"CreateFunctionStmt", "CREATE DEFINER=root PROCEDURE p() BEGIN END", "CREATE DEFINER=root PROCEDURE p() BEGIN END", 0},
	{"CreateFunctionStmt", "CREATE AGGREGATE FUNCTION f RETURNS STRING SONAME 'f.so'", "CREATE AGGREGATE FUNCTION f RETURNS STRING SONAME 'f.so'", 0},
	{"CreateIndexStmt", "CREATE INDEX idx ON t (a)", "CREATE INDEX idx ON t (a)", 0},
	{"CreateIndexStmt", "CREATE UNIQUE INDEX idx ON t (a)", "CREATE UNIQUE INDEX idx ON t (a)", 0},
	{"CreateLogfileGroupStmt", "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' ENGINE = NDB", "CREATE LOGFILE GROUP lg1 ADD UNDOFILE 'undo1.dat' ENGINE = NDB", 0},
	{"CreateResourceGroupStmt", "CREATE RESOURCE GROUP rg1 TYPE = USER", "CREATE RESOURCE GROUP rg1 TYPE = USER", 0},
	{"CreateRoleStmt", "CREATE ROLE r1", "CREATE ROLE r1", 0},
	{"CreateServerStmt", "CREATE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", "CREATE SERVER s1 FOREIGN DATA WRAPPER mysql OPTIONS (USER 'root')", 0},
	{"CreateSpatialRefSysStmt", "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[]'", "CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 4120 NAME 'Greek' DEFINITION 'GEOGCS[]'", 0},
	{"CreateTableStmt", "CREATE TABLE a (id INT)", "CREATE TABLE a (id INT)", 0},
	{"CreateTableStmt", "CREATE TEMPORARY TABLE IF NOT EXISTS a (id INT) ENGINE=InnoDB", "CREATE TEMPORARY TABLE IF NOT EXISTS a (id INT) ENGINE=InnoDB", 0},
	{"CreateTablespaceStmt", "CREATE TABLESPACE ts1", "CREATE TABLESPACE ts1", 0},
	{"CreateTriggerStmt", "CREATE TRIGGER trg AFTER UPDATE ON t FOR EACH ROW SET @a = 1", "CREATE TRIGGER trg AFTER UPDATE ON t FOR EACH ROW SET @a = 1", 0},
	{"CreateUserStmt", "CREATE USER testuser", "CREATE USER testuser", 0},
	{"CreateViewStmt", "CREATE VIEW v AS SELECT 1", "CREATE VIEW v AS SELECT 1", 0},
	{"CreateViewStmt", "CREATE OR REPLACE ALGORITHM=MERGE DEFINER=root SQL SECURITY INVOKER VIEW v AS SELECT 1", "CREATE OR REPLACE ALGORITHM=MERGE DEFINER=root SQL SECURITY INVOKER VIEW v AS SELECT 1", 0},
	{"DataType", "CREATE TABLE t (id INT)", "INT", 0},
	{"DataType", "CREATE AGGREGATE FUNCTION f RETURNS STRING SONAME 'f.so'", "STRING", 0},
	{"DatabaseOption", "ALTER DATABASE mydb READ ONLY = 1", "READ ONLY = 1", 0},
	{"DeallocateStmt", "DEALLOCATE PREPARE stmt1", "DEALLOCATE PREPARE stmt1", 0},
	{"DeallocateStmt", "DROP PREPARE stmt1", "DROP PREPARE stmt1", 0},
	{"DeclareConditionStmt", "CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR SQLSTATE '45000'; END", "DECLARE c CONDITION FOR SQLSTATE '45000'", 0},
	{"DeclareCursorStmt", "CREATE PROCEDURE p() BEGIN DECLARE cur CURSOR FOR SELECT 1; END", "DECLARE cur CURSOR FOR SELECT 1", 0},
	{"DeclareHandlerStmt", "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1062 SET @e=1; SELECT 1; END", "DECLARE EXIT HANDLER FOR 1062 SET @e=1", 0},
	{"DeclareVarStmt", "CREATE PROCEDURE p() BEGIN DECLARE x INT; END", "DECLARE x INT", 0},
	{"DefaultExpr", "SELECT DEFAULT(name) FROM t", "DEFAULT(name)", 0},
	{"DeleteStmt", "DELETE FROM t", "DELETE FROM t", 0},
	{"DeleteStmt", "WITH c AS (SELECT 1) DELETE FROM t", "WITH c AS (SELECT 1) DELETE FROM t", 0},
	{"DiagnosticsItem", "GET DIAGNOSTICS cnt = NUMBER", "cnt = NUMBER", 0},
	{"DoStmt", "DO 1", "DO 1", 0},
	{"DropDatabaseStmt", "DROP DATABASE db1", "DROP DATABASE db1", 0},
	{"DropEventStmt", "DROP EVENT ev1", "DROP EVENT ev1", 0},
	{"DropIndexStmt", "DROP INDEX idx ON t", "DROP INDEX idx ON t", 0},
	{"DropLogfileGroupStmt", "DROP LOGFILE GROUP lg1 ENGINE = NDB", "DROP LOGFILE GROUP lg1 ENGINE = NDB", 0},
	{"DropResourceGroupStmt", "DROP RESOURCE GROUP rg1", "DROP RESOURCE GROUP rg1", 0},
	{"DropRoleStmt", "DROP ROLE r1", "DROP ROLE r1", 0},
	{"DropRoutineStmt", "DROP FUNCTION f1", "DROP FUNCTION f1", 0},
	{"DropRoutineStmt", "DROP PROCEDURE IF EXISTS p", "DROP PROCEDURE IF EXISTS p", 0},
	{"DropServerStmt", "DROP SERVER s1", "DROP SERVER s1", 0},
	{"DropSpatialRefSysStmt", "DROP SPATIAL REFERENCE SYSTEM 4120", "DROP SPATIAL REFERENCE SYSTEM 4120", 0},
	{"DropTableStmt", "DROP TABLE IF EXISTS t, p", "DROP TABLE IF EXISTS t, p", 0},
	{"DropTableStmt", "SELECT 1;\nDROP TEMPORARY TABLE t ;", "DROP TEMPORARY TABLE t", 0},
	{"DropTablespaceStmt", "DROP TABLESPACE ts1", "DROP TABLESPACE ts1", 0},
	{"DropTriggerStmt", "DROP TRIGGER tr1", "DROP TRIGGER tr1", 0},
	{"DropUserStmt", "DROP USER IF EXISTS u1, u2", "DROP USER IF EXISTS u1, u2", 0},
	{"DropViewStmt", "DROP VIEW v", "DROP VIEW v", 0},
	{"ElseIf", "CREATE PROCEDURE p() BEGIN IF x THEN SELECT 1; ELSEIF y THEN SELECT 2; END IF; END", "ELSEIF y THEN SELECT 2;", 0},
	{"EventSchedule", "CREATE EVENT ev ON SCHEDULE EVERY 1 HOUR DO SELECT 1", "EVERY 1 HOUR", 0},
	{"ExecuteStmt", "EXECUTE stmt", "EXECUTE stmt", 0},
	{"ExistsExpr", "SELECT EXISTS (SELECT 1 FROM t)", "EXISTS (SELECT 1 FROM t)", 0},
	{"ExplainStmt", "EXPLAIN SELECT * FROM t", "EXPLAIN SELECT * FROM t", 0},
	{"ExtractExpr", "SELECT EXTRACT(YEAR FROM ts)", "EXTRACT(YEAR FROM ts)", 0},
	{"FactorOp", "ALTER USER 'admin'@'localhost' DROP 2 FACTOR", "DROP 2 FACTOR", 0},
	{"FetchCursorStmt", "CREATE PROCEDURE p_fetch_undecl() BEGIN FETCH nope INTO @x; END", "FETCH nope INTO @x", 0},
	{"FloatLit", "SELECT TRUNCATE(1.234, 2)", "1.234", 0},
	{"FlushStmt", "FLUSH TABLES", "FLUSH TABLES", 0},
	{"ForUpdate", "SELECT * FROM t FOR SHARE", "FOR SHARE", 0},
	{"FuncCallExpr", "DO SLEEP(1)", "SLEEP(1)", 0},
	{"FuncParam", "CREATE FUNCTION myfunc(a INT) RETURNS INT RETURN a + 1", "a INT", 0},
	{"GeneratedColumn", "CREATE TABLE t (a INT, b INT, c INT GENERATED ALWAYS AS (a + b) STORED)", "GENERATED ALWAYS AS (a + b) STORED", 0},
	{"GetDiagnosticsStmt", "GET DIAGNOSTICS cnt = NUMBER", "GET DIAGNOSTICS cnt = NUMBER", 0},
	{"GrantRoleStmt", "GRANT r1 TO u1", "GRANT r1 TO u1", 0},
	{"GrantStmt", "GRANT ALL PRIVILEGES ON *.* TO root", "GRANT ALL PRIVILEGES ON *.* TO root", 0},
	{"GrantTarget", "GRANT ALL PRIVILEGES ON *.* TO root", "*.*", 0},
	{"HandlerCloseStmt", "HANDLER t CLOSE", "HANDLER t CLOSE", 0},
	{"HandlerOpenStmt", "HANDLER t OPEN", "HANDLER t OPEN", 0},
	{"HandlerReadStmt", "HANDLER t1 READ NEXT", "HANDLER t1 READ NEXT", 0},
	{"HelpStmt", "HELP 'SELECT'", "HELP 'SELECT'", 0},
	{"HexLit", "SELECT x'41' 'b'", "x'41'", 0},
	{"IfStmt", "CREATE PROCEDURE p() BEGIN IF 1 THEN SELECT 1; END IF; END", "IF 1 THEN SELECT 1; END IF", 0},
	{"ImportTableStmt", "IMPORT TABLE FROM '/tmp/t1.sdi'", "IMPORT TABLE FROM '/tmp/t1.sdi'", 0},
	{"InExpr", "SELECT 1 IN ((SELECT 1))", "1 IN ((SELECT 1))", 0},
	{"InExpr", "SELECT * FROM t WHERE a NOT IN (1, 2) ORDER BY a", "a NOT IN (1, 2)", 0},
	{"IndexColumn", "CREATE INDEX idx ON t (a)", "a", 0},
	{"IndexHint", "SELECT * FROM t1 USE INDEX ()", "USE INDEX ()", 0},
	{"IndexOption", "CREATE TABLE t (id INT, PRIMARY KEY (id) VISIBLE)", "VISIBLE", 0},
	{"InsertStmt", "INSERT INTO t1 TABLE t2", "INSERT INTO t1 TABLE t2", 0},
	{"InstallComponentStmt", "INSTALL COMPONENT 'file://c1'", "INSTALL COMPONENT 'file://c1'", 0},
	{"InstallPluginStmt", "INSTALL PLUGIN p SONAME 'p.so'", "INSTALL PLUGIN p SONAME 'p.so'", 0},
	{"IntLit", "DO 1", "1", 0},
	{"IntLit", "CREATE INDEX idx ON t (a) KEY_BLOCK_SIZE = 8", "8", 0},
	{"IntervalExpr", "SELECT NOW() + INTERVAL 1 DAY", "INTERVAL 1 DAY", 0},
	{"IntoClause", "TABLE t1 INTO @v1, @v2", "@v1, @v2", 0},
	{"IsExpr", "SELECT 1 IS NULL", "1 IS NULL", 0},
	{"IsExpr", "SELECT * FROM t WHERE a IS NOT NULL ORDER BY a", "a IS NOT NULL", 0},
	{"IterateStmt", "CREATE PROCEDURE p_it_undecl() BEGIN ITERATE nowhere; END", "ITERATE nowhere", 0},
	{"JoinClause", "SELECT * FROM (t1, t2)", "t1, t2", 0},
	{"JoinClause", "SELECT * FROM t0 LEFT JOIN (t1, t2) ON t0.id = t1.id", "t0 LEFT JOIN (t1, t2) ON t0.id = t1.id", 0},
	{"JoinClause", "SELECT * FROM t0 LEFT JOIN (t1, t2) ON t0.id = t1.id", "t1, t2", 1},
	{"JsonTableColumn", "SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (v INT PATH '$')) AS j", "v INT PATH '$'", 0},
	{"JsonTableExpr", "SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (v INT PATH '$')) AS j", "JSON_TABLE('[]', '$[*]' COLUMNS (v INT PATH '$')) AS j", 0},
	{"KeywordArg", "SELECT TIMESTAMPADD(MINUTE, 30, NOW())", "MINUTE", 0},
	{"KillStmt", "KILL 123", "KILL 123", 0},
	{"LeaveStmt", "CREATE PROCEDURE p_lv_undecl() BEGIN LEAVE nowhere; END", "LEAVE nowhere", 0},
	{"LikeExpr", "SELECT a LIKE 'foo%' FROM t", "a LIKE 'foo%'", 0},
	{"LikeExpr", "SELECT a NOT LIKE 'f!%' ESCAPE '!' FROM t", "a NOT LIKE 'f!%' ESCAPE '!'", 0},
	{"Limit", "TABLE t1 LIMIT 10", "10", 0},
	{"LoadDataStmt", "LOAD DATA INFILE '/tmp/x.csv' INTO TABLE t1", "LOAD DATA INFILE '/tmp/x.csv' INTO TABLE t1", 0},
	{"LoadIndexIntoCacheStmt", "LOAD INDEX INTO CACHE t1", "LOAD INDEX INTO CACHE t1", 0},
	{"LoadIndexTable", "LOAD INDEX INTO CACHE t1", "t1", 0},
	{"LockInstanceStmt", "LOCK INSTANCE FOR BACKUP", "LOCK INSTANCE FOR BACKUP", 0},
	{"LockTable", "LOCK TABLE t READ", "t READ", 0},
	{"LockTablesStmt", "LOCK TABLE t READ", "LOCK TABLE t READ", 0},
	{"LoopStmt", "CREATE PROCEDURE p() BEGIN LOOP SELECT 1; END LOOP; END", "LOOP SELECT 1; END LOOP", 0},
	{"MatchExpr", "SELECT * FROM articles WHERE MATCH (title, body) AGAINST ('database')", "MATCH (title, body) AGAINST ('database')", 0},
	{"MemberOfExpr", "SELECT 17 MEMBER OF('[23, 17, 10]')", "17 MEMBER OF('[23, 17, 10]')", 0},
	{"NullLit", "SELECT * FROM t WHERE n = NULL ORDER BY n", "NULL", 0},
	{"OnCondition", "SELECT * FROM a JOIN b ON a.id = b.id", "a.id = b.id", 0},
	{"OpenCursorStmt", "CREATE PROCEDURE p_open_undecl() BEGIN OPEN nope; END", "OPEN nope", 0},
	{"OptimizeTableStmt", "OPTIMIZE TABLE t1", "OPTIMIZE TABLE t1", 0},
	{"OrderByItem", "TABLE t1 ORDER BY id", "id", 0},
	{"ParenExpr", "SELECT ((1))", "((1))", 0},
	{"PartitionClause", "ALTER TABLE t PARTITION BY HASH(id) PARTITIONS 4", "PARTITION BY HASH(id) PARTITIONS 4", 0},
	{"PartitionDef", "ALTER TABLE t ADD PARTITION (PARTITION p4 VALUES IN (10, 20, 30))", "PARTITION p4 VALUES IN (10, 20, 30)", 0},
	{"PrepareStmt", "PREPARE stmt FROM 'SELECT ?'", "PREPARE stmt FROM 'SELECT ?'", 0},
	{"PurgeBinaryLogsStmt", "PURGE BINARY LOGS TO 'mysql-bin.010'", "PURGE BINARY LOGS TO 'mysql-bin.010'", 0},
	{"RegistrationOp", "ALTER USER USER() 2 FACTOR UNREGISTER", "2 FACTOR UNREGISTER", 0},
	{"RenameTablePair", "RENAME TABLE t1 TO t2", "t1 TO t2", 0},
	{"RenameTableStmt", "RENAME TABLE t1 TO t2", "RENAME TABLE t1 TO t2", 0},
	{"RenameUserPair", "RENAME USER foo TO bar", "foo TO bar", 0},
	{"RenameUserStmt", "RENAME USER foo TO bar", "RENAME USER foo TO bar", 0},
	{"RepairTableStmt", "REPAIR TABLE t1", "REPAIR TABLE t1", 0},
	{"RepeatStmt", "CREATE PROCEDURE p() BEGIN REPEAT SET x = x + 1; UNTIL x > 5 END REPEAT; END", "REPEAT SET x = x + 1; UNTIL x > 5 END REPEAT", 0},
	{"ReplicationFilter", "CHANGE REPLICATION FILTER REPLICATE_DO_DB = (db1)", "REPLICATE_DO_DB = (db1)", 0},
	{"ReplicationOption", "CHANGE MASTER TO MASTER_SSL = 0", "MASTER_SSL = 0", 0},
	{"RequireClause", "CREATE USER 'jeffrey'@'localhost' REQUIRE X509", "REQUIRE X509", 0},
	{"ResTarget", "SELECT 1 'x'", "1 'x'", 0},
	{"ResetMasterStmt", "RESET MASTER", "RESET MASTER", 0},
	{"ResetPersistStmt", "RESET PERSIST", "RESET PERSIST", 0},
	{"ResetReplicaStmt", "RESET SLAVE", "RESET SLAVE", 0},
	{"ResignalStmt", "RESIGNAL", "RESIGNAL", 0},
	{"ResourceOption", "GRANT ALL ON *.* TO 'jeffrey'@'localhost' WITH MAX_USER_CONNECTIONS 10", "MAX_USER_CONNECTIONS 10", 0},
	{"RestartStmt", "RESTART", "RESTART", 0},
	{"ReturnStmt", "CREATE FUNCTION myfunc() RETURNS INT RETURN 1", "RETURN 1", 0},
	{"RevokeRoleStmt", "REVOKE myrole FROM myuser", "REVOKE myrole FROM myuser", 0},
	{"RevokeStmt", "REVOKE SELECT ON mydb.* FROM app_user", "REVOKE SELECT ON mydb.* FROM app_user", 0},
	{"RollbackStmt", "ROLLBACK", "ROLLBACK", 0},
	{"RoutineCharacteristic", "ALTER FUNCTION f1 COMMENT 'test'", "COMMENT 'test'", 0},
	{"RowExpr", "SELECT ROW(1)", "ROW(1)", 0},
	{"SavepointStmt", "SAVEPOINT sp1", "SAVEPOINT sp1", 0},
	{"SelectStmt", "SELECT 1", "SELECT 1", 0},
	{"SelectStmt", "SELECT (TABLE t1 ORDER BY a LIMIT 1)", "TABLE t1 ORDER BY a LIMIT 1", 1},
	{"SetDefaultRoleStmt", "SET DEFAULT ROLE ALL TO root", "SET DEFAULT ROLE ALL TO root", 0},
	{"SetPasswordStmt", "SET PASSWORD TO RANDOM", "SET PASSWORD TO RANDOM", 0},
	{"SetResourceGroupStmt", "SET RESOURCE GROUP rg1", "SET RESOURCE GROUP rg1", 0},
	{"SetRoleStmt", "SET ROLE ALL", "SET ROLE ALL", 0},
	{"SetStmt", "SET x = 1", "SET x = 1", 0},
	{"SetTransactionStmt", "SET TRANSACTION READ ONLY", "SET TRANSACTION READ ONLY", 0},
	{"ShowStmt", "SHOW TABLES", "SHOW TABLES", 0},
	{"ShutdownStmt", "SHUTDOWN", "SHUTDOWN", 0},
	{"SignalInfoItem", "RESIGNAL SET MESSAGE_TEXT = 'Error handled'", "MESSAGE_TEXT = 'Error handled'", 0},
	{"SignalStmt", "SIGNAL my_error", "SIGNAL my_error", 0},
	{"StarExpr", "SELECT *", "*", 0},
	{"StartGroupReplicationStmt", "START GROUP_REPLICATION", "START GROUP_REPLICATION", 0},
	{"StartReplicaStmt", "START SLAVE", "START SLAVE", 0},
	{"StopGroupReplicationStmt", "STOP GROUP_REPLICATION", "STOP GROUP_REPLICATION", 0},
	{"StopReplicaStmt", "STOP SLAVE", "STOP SLAVE", 0},
	{"StringLit", "SELECT 'a'", "'a'", 0},
	{"StringLit", "SELECT 'a' 'b' FROM t", "'a' 'b'", 0},
	{"StringLit", "SELECT _utf8mb4'it''s' FROM t", "_utf8mb4'it''s'", 0},
	{"StringLit", "CREATE INDEX idx ON t (a) COMMENT 'c' USING BTREE", "'c'", 0},
	{"StringLit", "CREATE INDEX idx ON t (a) COMMENT 'c' USING BTREE", "BTREE", 1},
	{"SubPartitionDef", "CREATE TABLE t (id INT, d DATE) PARTITION BY RANGE (YEAR(d)) SUBPARTITION BY HASH (id) (PARTITION p0 VALUES LESS THAN (2000) (SUBPARTITION s0, SUBPARTITION s1))", "SUBPARTITION s0", 0},
	{"SubqueryExpr", "SELECT (SELECT 1)", "(SELECT 1)", 0},
	{"TableOption", "ALTER TABLE t ENGINE=InnoDB", "ENGINE=InnoDB", 0},
	{"TableRef", "TABLE t1", "t1", 0},
	{"TableStmt", "TABLE t1", "TABLE t1", 0},
	{"TableStmt", "TABLE t1 ORDER BY a LIMIT 1", "TABLE t1 ORDER BY a LIMIT 1", 0},
	{"TableStmt", "SELECT (TABLE t1 ORDER BY a LIMIT 1)", "TABLE t1", 0},
	{"TemporalLit", "SELECT DATE '2024-01-01' FROM t", "DATE '2024-01-01'", 0},
	{"TriggerOrder", "CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW FOLLOWS other SET @a = 1", "FOLLOWS other", 0},
	{"TruncateStmt", "TRUNCATE t", "TRUNCATE t", 0},
	{"UnaryExpr", "SELECT -1", "-1", 0},
	{"UninstallComponentStmt", "UNINSTALL COMPONENT 'file://c1'", "UNINSTALL COMPONENT 'file://c1'", 0},
	{"UninstallPluginStmt", "UNINSTALL PLUGIN p", "UNINSTALL PLUGIN p", 0},
	{"UnlockInstanceStmt", "UNLOCK INSTANCE", "UNLOCK INSTANCE", 0},
	{"UnlockTablesStmt", "UNLOCK TABLES", "UNLOCK TABLES", 0},
	{"UpdateStmt", "UPDATE t SET a = 1", "UPDATE t SET a = 1", 0},
	{"UpdateStmt", "WITH c AS (SELECT 1) UPDATE t SET a = 1", "WITH c AS (SELECT 1) UPDATE t SET a = 1", 0},
	{"UseStmt", "USE db1", "USE db1", 0},
	{"UserSpec", "CREATE USER testuser", "testuser", 0},
	{"UserSpec", "ALTER USER USER() IDENTIFIED BY 'x'", "USER() IDENTIFIED BY 'x'", 0},
	{"UsingCondition", "SELECT a FROM t1 JOIN t2 USING (id)", "USING (id)", 0},
	{"ValuesStmt", "VALUES ROW(1), ROW(2)", "VALUES ROW(1), ROW(2)", 0},
	{"ValuesStmt", "SELECT * FROM (VALUES ROW(1), ROW(2) LIMIT 1) x", "VALUES ROW(1), ROW(2)", 0},
	{"VariableRef", "SELECT @var", "@var", 0},
	{"WeightStringExpr", "SELECT WEIGHT_STRING(a AS CHAR(3)) FROM t", "WEIGHT_STRING(a AS CHAR(3))", 0},
	{"WhileStmt", "CREATE PROCEDURE p() BEGIN WHILE x > 0 DO SET x = x - 1; END WHILE; END", "WHILE x > 0 DO SET x = x - 1; END WHILE", 0},
	{"WindowDef", "SELECT ROW_NUMBER() OVER ()", "()", 0},
	{"WindowFrame", "SELECT SUM(x) OVER (ORDER BY id GROUPS CURRENT ROW) FROM t", "GROUPS CURRENT ROW", 0},
	{"WindowFrameBound", "SELECT SUM(x) OVER (ORDER BY id GROUPS CURRENT ROW) FROM t", "CURRENT ROW", 0},
	{"XAStmt", "XA RECOVER", "XA RECOVER", 0},
}

func TestLocSpans(t *testing.T) {
	for _, tc := range locSpanCases {
		t.Run(tc.node+"/"+tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.sql, err)
			}
			var found []ast.Loc
			for _, s := range locSpans(list.Items...) {
				if s.node == tc.node {
					found = append(found, s.loc)
				}
			}
			if tc.nth >= len(found) {
				t.Fatalf("Parse(%q): %d %s nodes, want at least %d", tc.sql, len(found), tc.node, tc.nth+1)
			}
			loc := found[tc.nth]
			if loc.Start < 0 || loc.End <= loc.Start || loc.End > len(tc.sql) {
				t.Fatalf("Parse(%q): %s Loc = %+v, want the range of %q", tc.sql, tc.node, loc, tc.want)
			}
			if got := tc.sql[loc.Start:loc.End]; got != tc.want {
				t.Errorf("Parse(%q): %s Loc = %+v covers %q, want %q", tc.sql, tc.node, loc, got, tc.want)
			}
		})
	}
}

// locKindsNotBuilt lists the node kinds with a Loc that the parser never
// builds, so TestLocSpans cannot pin them.
var locKindsNotBuilt = map[string]string{
	"RawStmt": "Parse returns statement nodes unwrapped",
}

// TestLocSpansCoverEveryNodeKind requires a TestLocSpans case for every
// struct in oceanbase/ast that carries a Loc, so a new node kind gets its
// location pinned when it is added.
func TestLocSpansCoverEveryNodeKind(t *testing.T) {
	file, err := goparser.ParseFile(token.NewFileSet(), filepath.Join("..", "ast", "parsenodes.go"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	pinned := map[string]bool{}
	for _, tc := range locSpanCases {
		pinned[tc.node] = true
	}
	kinds := map[string]bool{}
	for _, decl := range file.Decls {
		gd, ok := decl.(*goast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts := spec.(*goast.TypeSpec)
			st, ok := ts.Type.(*goast.StructType)
			if !ok || !hasLocField(st) {
				continue
			}
			name := ts.Name.Name
			kinds[name] = true
			if _, ok := locKindsNotBuilt[name]; ok {
				if pinned[name] {
					t.Errorf("%s is pinned in locSpanCases but listed in locKindsNotBuilt", name)
				}
				continue
			}
			if !pinned[name] {
				t.Errorf("ast.%s carries a Loc but has no locSpanCases entry", name)
			}
		}
	}
	for name := range pinned {
		if !kinds[name] {
			t.Errorf("locSpanCases pins %s, which is not a struct with a Loc in ast", name)
		}
	}
}

func hasLocField(st *goast.StructType) bool {
	for _, f := range st.Fields.List {
		for _, n := range f.Names {
			if n.Name == "Loc" {
				return true
			}
		}
	}
	return false
}

// TestLocContract checks every location in the parse of every SQL text this
// package's tests use, plus the quality corpus: a Loc is NoLoc or a
// non-empty range that starts at a token's start, ends at a token's end,
// and lies inside the nearest enclosing node's Loc; and the statements of a
// segment span it from its first token to its last. The tokens come from
// lexing each segment the way Parse does.
func TestLocContract(t *testing.T) {
	corpus := locContractCorpus(t)
	checked := 0
	var failures []string
	for _, sql := range corpus {
		var segFailures []string
		ok := true
		for _, seg := range Split(sql) {
			if seg.Empty() {
				continue
			}
			list, err := parseSingle(seg.Text, seg.ByteStart)
			if err != nil {
				ok = false
				break
			}
			segFailures = append(segFailures, checkSegmentLocs(sql, seg, list.Items)...)
		}
		if !ok {
			continue // only text the parser accepts has locations to check
		}
		checked++
		failures = append(failures, segFailures...)
	}
	if checked < 1000 {
		t.Fatalf("checked %d parseable texts out of %d; the corpus harvest is broken", checked, len(corpus))
	}
	for i, f := range failures {
		if i == 50 {
			t.Errorf("... and %d more", len(failures)-50)
			break
		}
		t.Error(f)
	}
}

func checkSegmentLocs(sql string, seg Segment, items []ast.Node) []string {
	var failures []string
	starts, ends := map[int]bool{}, map[int]bool{}
	first, last := -1, -1
	lex := NewLexerWithOffset(seg.Text, seg.ByteStart)
	for {
		tok := lex.NextToken()
		if tok.Type == tokEOF {
			break
		}
		starts[tok.Loc], ends[tok.End] = true, true
		if tok.Type == ';' {
			continue
		}
		if first < 0 {
			first = tok.Loc
		}
		last = tok.End
	}
	fail := func(s locSpan, format string, args ...any) {
		text := ""
		if s.loc.Start >= 0 && s.loc.Start <= s.loc.End && s.loc.End <= len(sql) {
			text = sql[s.loc.Start:s.loc.End]
		}
		failures = append(failures, fmt.Sprintf("Parse(%q): %s Loc %+v (%q): %s", sql, s.node, s.loc, text, fmt.Sprintf(format, args...)))
	}
	for _, s := range locSpans(items...) {
		loc := s.loc
		switch {
		case loc == ast.NoLoc():
			continue
		case loc.Start < seg.ByteStart || loc.End > seg.ByteEnd || loc.End <= loc.Start:
			fail(s, "not a non-empty range inside its segment [%d,%d)", seg.ByteStart, seg.ByteEnd)
			continue
		}
		if !starts[loc.Start] {
			fail(s, "Start is not the start of a token")
		}
		if !ends[loc.End] {
			fail(s, "End is not the end of a token")
		}
		if p := s.parent; p != ast.NoLoc() && (loc.Start < p.Start || loc.End > p.End) {
			fail(s, "outside its parent's Loc %+v", p)
		}
	}
	if len(items) > 0 {
		head, tail := locSpans(items[0])[0], locSpans(items[len(items)-1])[0]
		if head.loc.Start != first {
			fail(head, "statement starts at %d, want its first token at %d", head.loc.Start, first)
		}
		if tail.loc.End != last {
			fail(tail, "statement ends at %d, want its last token's end at %d", tail.loc.End, last)
		}
	}
	return failures
}

// locContractCorpus returns every string literal in this package's tests
// and every blank-line-separated entry of the quality corpus. Most are SQL;
// the rest fail to parse and are skipped by the caller.
func locContractCorpus(t *testing.T) []string {
	t.Helper()
	seen := map[string]bool{}
	var corpus []string
	add := func(s string) {
		if strings.TrimSpace(s) != "" && !seen[s] {
			seen[s] = true
			corpus = append(corpus, s)
		}
	}
	files, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	for _, name := range files {
		file, err := goparser.ParseFile(fset, name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		goast.Inspect(file, func(n goast.Node) bool {
			if lit, ok := n.(*goast.BasicLit); ok && lit.Kind == token.STRING {
				if s, err := strconv.Unquote(lit.Value); err == nil {
					add(s)
				}
			}
			return true
		})
	}
	sqlFiles, err := filepath.Glob(filepath.Join("..", "quality", "corpus", "*.sql"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range sqlFiles {
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range strings.Split(string(data), "\n\n") {
			add(entry)
		}
	}
	sort.Strings(corpus)
	return corpus
}
