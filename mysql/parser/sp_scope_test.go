package parser

import (
	"strings"
	"testing"
)

// TestStoredProgramLabels holds LEAVE, ITERATE and a block's or loop's label to the labels in
// reach, as MySQL's grammar actions do (sp_pcontext::find_label). Each refusal names the
// error mysql:8.0.46 answers to PREPARE of the statement; wantAt is the text the error points
// at.
func TestStoredProgramLabels(t *testing.T) {
	cases := []struct {
		sql    string
		code   int // 0: the server takes it
		wantAt string
	}{
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN LEAVE a; END; END LOOP"},
		{sql: "CREATE PROCEDURE p() a: BEGIN b: BEGIN LEAVE a; END b; END a"},
		{sql: "CREATE PROCEDURE p() a: WHILE 1 DO ITERATE a; END WHILE"},
		{sql: "CREATE PROCEDURE p() a: REPEAT ITERATE A; UNTIL 1 END REPEAT"},
		{sql: "CREATE PROCEDURE p() BEGIN a: LOOP LEAVE a; END LOOP; a: LOOP LEAVE a; END LOOP; END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN DECLARE EXIT HANDLER FOR 1 a: LOOP LEAVE a; END LOOP; END; END LOOP"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 b: BEGIN LEAVE b; END; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT a: BEGIN LEAVE a; RETURN 1; END"},
		{sql: "CREATE TRIGGER t BEFORE INSERT ON t FOR EACH ROW a: BEGIN LEAVE a; END"},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO a: LOOP LEAVE a; END LOOP"},
		{sql: "ALTER EVENT e DO a: LOOP LEAVE a; END LOOP"},

		{sql: "CREATE PROCEDURE p() LEAVE a", code: 1308, wantAt: "a"},
		{sql: "CREATE PROCEDURE p() é: LOOP LEAVE e; END LOOP", code: 1308, wantAt: "e; END"},
		{sql: "ALTER EVENT e DO ITERATE a", code: 1308, wantAt: "a"},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN LEAVE b; END", code: 1308, wantAt: "b; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN ITERATE a; END", code: 1308, wantAt: "a; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN b: BEGIN END b; LEAVE b; END", code: 1308, wantAt: "b; END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN DECLARE EXIT HANDLER FOR 1 LEAVE a; END; END LOOP", code: 1308, wantAt: "a; END; END"},
		{sql: "CREATE PROCEDURE p() a: BEGIN a: BEGIN END; END", code: 1309, wantAt: "a: BEGIN END"},
		{sql: "CREATE PROCEDURE p() a: LOOP BEGIN A: LOOP LEAVE a; END LOOP; END; END LOOP", code: 1309, wantAt: "A: LOOP"},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 {
				if err != nil {
					t.Fatalf("Parse refused what MySQL takes: %v", err)
				}
				return
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("Parse took what MySQL refuses with %d (err %v)", c.code, err)
			}
			if !strings.HasPrefix(c.sql[pe.Position:], c.wantAt) {
				t.Errorf("error at %q, want at %q: %v", c.sql[pe.Position:], c.wantAt, pe)
			}
		})
	}
}

// TestStoredProgramDeclarations holds parameters, variables and cursors to a stored program's
// scope (sp_pcontext::find_variable, find_cursor), and DECLARE to the head of a BEGIN ... END
// block. code is what mysql:8.0.46 answers to PREPARE of the statement, 0 when the compat
// judge reads that as parsed (compat.PrepareVerdict).
func TestStoredProgramDeclarations(t *testing.T) {
	cases := []struct {
		sql  string
		code int
	}{
		{sql: "CREATE PROCEDURE p(a INT, b INT) BEGIN DECLARE a INT; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; BEGIN DECLARE x INT; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE x CURSOR FOR SELECT 1; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; BEGIN DECLARE c CURSOR FOR SELECT 2; OPEN c; END; CLOSE c; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE EXIT HANDLER FOR 1 CLOSE c; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CURSOR FOR SELECT 1; FETCH NEXT FROM c INTO x; END; END"},
		{sql: "CREATE PROCEDURE p(n INT) BEGIN DECLARE c CURSOR FOR SELECT 1 LIMIT n; SELECT 1 INTO n; GET DIAGNOSTICS n = NUMBER; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE next CURSOR FOR SELECT 1; FETCH NEXT FROM next INTO x; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE next CURSOR FOR SELECT 1; FETCH next INTO x; END", code: 1064},
		{sql: "CREATE FUNCTION f(a INT) RETURNS INT BEGIN DECLARE b INT; SELECT 1 INTO a; SELECT 2 INTO b LIMIT a; RETURN b; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE é INT; GET DIAGNOSTICS e = NUMBER; END"},
		{sql: "GET DIAGNOSTICS @n = NUMBER, @r = ROW_COUNT"},
		{sql: "SELECT 1 INTO a"},
		{sql: "SELECT 1 LIMIT a"},
		{sql: "CREATE PROCEDURE p(a INT, A INT) BEGIN END", code: 1330},
		{sql: "CREATE PROCEDURE p(IN a INT, INOUT a INT) BEGIN END", code: 1330},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x, y, X INT; END", code: 1331},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; DECLARE C CURSOR FOR SELECT 2; END", code: 1333},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1 INTO @a; END", code: 1323},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1 UNION SELECT 2 INTO @a; END", code: 1323},
		{sql: "CREATE PROCEDURE p() OPEN c", code: 1324},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE c CURSOR FOR SELECT 1; END; OPEN c; END", code: 1324},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE c CURSOR FOR SELECT 1; FETCH c INTO y; END", code: 1327},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH c INTO @x; END", code: 1064},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE c CURSOR FOR SELECT 1; FETCH NEXT c INTO x; END", code: 1064},
		{sql: "CREATE PROCEDURE p() GET DIAGNOSTICS a = NUMBER", code: 1327},
		{sql: "GET DIAGNOSTICS a = NUMBER", code: 1327},
		{sql: "GET DIAGNOSTICS CONDITION 1 a = MESSAGE_TEXT", code: 1327},
		{sql: "CREATE PROCEDURE p() SELECT 1 INTO a", code: 1327},
		{sql: "CREATE PROCEDURE p() BEGIN BEGIN DECLARE a INT; END; SELECT 1 INTO a; END", code: 1327},
		{sql: "CREATE TRIGGER t BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE a INT; SELECT 1 INTO a, b; END", code: 1327},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1 LIMIT n", code: 1327},
		{sql: "CREATE PROCEDURE p() DECLARE x INT", code: 1064},
		{sql: "CREATE PROCEDURE p() a: LOOP DECLARE x INT; END LOOP", code: 1064},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 DECLARE y INT; END", code: 1064},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 && err != nil {
				t.Fatalf("Parse refused what MySQL takes: %v", err)
			}
			if c.code != 0 && err == nil {
				t.Fatalf("Parse took what MySQL refuses with %d", c.code)
			}
		})
	}
}

// TestStoredProgramConditions holds conditions, handlers, SIGNAL and RESIGNAL to a stored program's
// scope (sp_pcontext::find_condition, check_duplicate_handler). code is what mysql:8.0.46 answers to PREPARE of the
// statement, 0 when the compat judge reads that as parsed (compat.PrepareVerdict).
func TestStoredProgramConditions(t *testing.T) {
	cases := []struct {
		sql  string
		code int
	}{
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR SQLSTATE '45000'; SIGNAL x; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR SQLSTATE '45000'; BEGIN DECLARE EXIT HANDLER FOR x BEGIN END; RESIGNAL x; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1; BEGIN DECLARE x CONDITION FOR 2; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1; DECLARE EXIT HANDLER FOR x BEGIN END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x INT; DECLARE x CONDITION FOR 1; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN END; BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN END; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLSTATE '23000' BEGIN END; DECLARE EXIT HANDLER FOR 23000 BEGIN END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLWARNING, NOT FOUND, SQLEXCEPTION, 1, SQLSTATE '42000' BEGIN END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN DECLARE EXIT HANDLER FOR 1 BEGIN END; END; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR SQLSTATE '45000'; DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN DECLARE y CONDITION FOR SQLSTATE '45001'; SIGNAL y; END; SIGNAL x; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE é CONDITION FOR SQLSTATE '45000'; SIGNAL e; END"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE e CONDITION FOR 1; BEGIN DECLARE é CONDITION FOR SQLSTATE '45000'; SIGNAL e; END; END"},
		{sql: "SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'm'"},
		{sql: "RESIGNAL SET MESSAGE_TEXT = 'm'"},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 01062 BEGIN END; DECLARE EXIT HANDLER FOR 1062 BEGIN END; END", code: 1413},
		{sql: "SIGNAL x", code: 1319},
		{sql: "RESIGNAL x", code: 1319},
		{sql: "CREATE PROCEDURE p() SIGNAL x", code: 1319},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1; SIGNAL x; END", code: 1646},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1; DECLARE EXIT HANDLER FOR SQLEXCEPTION RESIGNAL x; END", code: 1646},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1; DECLARE X CONDITION FOR 2; END", code: 1332},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR y BEGIN END; END", code: 1319},
		{sql: "CREATE PROCEDURE p() BEGIN BEGIN DECLARE y CONDITION FOR 1; END; BEGIN DECLARE EXIT HANDLER FOR y BEGIN END; END; END", code: 1319},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR 1, 2, 1 BEGIN END; END", code: 1413},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE EXIT HANDLER FOR SQLWARNING BEGIN END; DECLARE CONTINUE HANDLER FOR SQLWARNING BEGIN END; END", code: 1413},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR SQLSTATE '23000'; DECLARE EXIT HANDLER FOR SQLSTATE '23000' BEGIN END; DECLARE EXIT HANDLER FOR x BEGIN END; END", code: 1413},
		{sql: "CREATE PROCEDURE p() BEGIN DECLARE x CONDITION FOR 1062; DECLARE y CONDITION FOR 1062; DECLARE EXIT HANDLER FOR x, y BEGIN END; END", code: 1413},
		{sql: "CREATE TRIGGER t BEFORE INSERT ON t FOR EACH ROW SIGNAL x", code: 1319},
		{sql: "CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN DECLARE EXIT HANDLER FOR NOT FOUND, NOT FOUND BEGIN END; END", code: 1413},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 && err != nil {
				t.Fatalf("Parse refused what MySQL takes: %v", err)
			}
			if c.code != 0 && err == nil {
				t.Fatalf("Parse took what MySQL refuses with %d", c.code)
			}
		})
	}
}

// TestStoredProgramStatements holds a stored program's statements to what each kind of
// program may hold (spBadStatement, spNotInFunction, refuseResultSet). Each statement is
// written as the body of a function, a trigger, a procedure and an event; the codes are what
// mysql:8.0.46 answers to PREPARE of each, 0 when the compat judge reads that as parsed.
func TestStoredProgramStatements(t *testing.T) {
	cases := []struct {
		stmt               string
		fn, trg, proc, evt int
	}{
		{"COMMIT", 1422, 1422, 0, 0},
		{"ROLLBACK", 1422, 1422, 0, 0},
		{"ROLLBACK TO SAVEPOINT a", 0, 0, 0, 0},
		{"SAVEPOINT a", 0, 0, 0, 0},
		{"RELEASE SAVEPOINT a", 0, 0, 0, 0},
		{"START TRANSACTION", 1422, 1422, 0, 0},
		{"XA START 'a'", 0, 0, 0, 0},
		{"SET autocommit = 1", 1445, 1445, 0, 0},
		{"SET GLOBAL autocommit = 1", 1445, 1445, 0, 0},
		{"SET @@session.autocommit = 1", 1445, 1445, 0, 0},
		{"SET sql_mode = '', autocommit = 0", 1445, 1445, 0, 0},
		{"SET @autocommit = 1", 0, 0, 0, 0},
		{"SET PASSWORD = 'x'", 1445, 1445, 0, 0},
		{"SET TRANSACTION READ ONLY", 0, 0, 0, 0},
		{"PREPARE s FROM 'SELECT 1'", 1336, 1336, 0, 0},
		{"EXECUTE s", 1336, 1336, 0, 0},
		{"DEALLOCATE PREPARE s", 1336, 1336, 0, 0},
		{"FLUSH TABLES", 1336, 1336, 0, 0},
		{"RESET MASTER", 1336, 1336, 0, 0},
		{"RESET SLAVE", 1336, 1336, 0, 0},
		{"RESET PERSIST", 1336, 1336, 0, 0},
		{"CREATE TABLE t2 (a INT)", 1422, 1422, 0, 0},
		{"CREATE TEMPORARY TABLE t2 (a INT)", 0, 0, 0, 0},
		{"DROP TABLE t2", 1422, 1422, 0, 0},
		{"DROP TEMPORARY TABLE t2", 0, 0, 0, 0},
		{"ALTER TABLE t ADD c INT", 1422, 1422, 0, 0},
		{"TRUNCATE t", 1422, 1422, 0, 0},
		{"CREATE INDEX i ON t (a)", 1422, 1422, 0, 0},
		{"DROP INDEX i ON t", 1422, 1422, 0, 0},
		{"RENAME TABLE t TO t3", 1422, 1422, 0, 0},
		{"CREATE VIEW v AS SELECT 1", 1422, 1422, 0, 0},
		{"DROP VIEW v", 1422, 1422, 0, 0},
		{"ALTER VIEW v AS SELECT 1", 1314, 1314, 1314, 1314},
		{"CREATE DATABASE d2", 1422, 1422, 0, 0},
		{"ALTER DATABASE d2 CHARACTER SET utf8mb4", 1422, 1422, 0, 0},
		{"DROP DATABASE d2", 1422, 1422, 0, 0},
		{"CREATE USER u", 1422, 1422, 0, 0},
		{"ALTER USER u ACCOUNT LOCK", 1422, 1422, 0, 0},
		{"DROP USER u", 1422, 1422, 0, 0},
		{"RENAME USER u TO v", 1422, 1422, 0, 0},
		{"GRANT SELECT ON *.* TO u", 1422, 1422, 0, 0},
		{"REVOKE SELECT ON *.* FROM u", 1422, 1422, 0, 0},
		{"CREATE ROLE r", 1422, 1422, 0, 0},
		{"DROP ROLE r", 1422, 1422, 0, 0},
		{"GRANT r TO u", 1422, 1422, 0, 0},
		{"SET DEFAULT ROLE r TO u", 0, 0, 0, 0},
		{"SET ROLE r", 0, 0, 0, 0},
		{"INSTALL PLUGIN p SONAME 'p.so'", 1422, 1422, 0, 0},
		{"UNINSTALL PLUGIN p", 1422, 1422, 0, 0},
		{"INSTALL COMPONENT 'file://c'", 0, 0, 0, 0},
		{"CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (USER 'u')", 1422, 1422, 0, 0},
		{"DROP SERVER s", 1422, 1422, 0, 0},
		{"CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'h'", 1422, 1422, 0, 0},
		{"START REPLICA", 1422, 1422, 0, 0},
		{"STOP REPLICA", 1422, 1422, 0, 0},
		{"START GROUP_REPLICATION", 0, 0, 0, 0},
		{"PURGE BINARY LOGS TO 'b'", 0, 0, 0, 0},
		{"ALTER INSTANCE ROTATE INNODB MASTER KEY", 1422, 1422, 0, 0},
		{"LOCK INSTANCE FOR BACKUP", 0, 0, 0, 0},
		{"CREATE TABLESPACE ts ADD DATAFILE 'a.ibd'", 1422, 1422, 0, 0},
		{"CREATE RESOURCE GROUP g TYPE = USER", 1422, 1422, 0, 0},
		{"CREATE SPATIAL REFERENCE SYSTEM 1000000 NAME 'n' DEFINITION 'd'", 1422, 1422, 0, 0},
		{"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1", 1576, 1576, 1576, 1576},
		{"ALTER EVENT e DO SELECT 1", 1576, 1576, 1576, 1576},
		{"ALTER EVENT e ENABLE", 1422, 1422, 0, 0},
		{"DROP EVENT e", 1422, 1422, 0, 0},
		{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @a = 1", 1303, 1303, 1303, 1303},
		{"DROP TRIGGER tr", 1422, 1422, 0, 0},
		{"CREATE PROCEDURE q() SELECT 1", 1303, 1303, 1303, 1303},
		{"CREATE FUNCTION g() RETURNS INT RETURN 1", 1303, 1303, 1303, 1303},
		{"CREATE FUNCTION g RETURNS STRING SONAME 'x.so'", 0, 0, 0, 0},
		{"DROP PROCEDURE q", 1357, 1357, 1357, 1357},
		{"DROP FUNCTION g", 1357, 1357, 1357, 1357},
		{"ALTER PROCEDURE q COMMENT 'c'", 1357, 1357, 1357, 1357},
		{"USE d", 1314, 1314, 1314, 1314},
		{"LOCK TABLES t READ", 1314, 1314, 1314, 1314},
		{"UNLOCK TABLES", 1314, 1314, 1314, 1314},
		{"CHECK TABLE t", 1314, 1314, 1314, 1314},
		{"LOAD DATA INFILE 'x' INTO TABLE t", 1314, 1314, 1314, 1314},
		{"LOAD XML INFILE 'x' INTO TABLE t", 1314, 1314, 1314, 1314},
		{"HELP 'x'", 1314, 1314, 1314, 1314},
		{"ANALYZE TABLE t", 1415, 1415, 0, 0},
		{"OPTIMIZE TABLE t", 1415, 1415, 0, 0},
		{"REPAIR TABLE t", 1415, 1415, 0, 0},
		{"CHECKSUM TABLE t", 1415, 1415, 0, 0},
		{"CACHE INDEX t IN k", 1415, 1415, 0, 0},
		{"LOAD INDEX INTO CACHE t", 1415, 1415, 0, 0},
		{"HANDLER t OPEN", 1314, 1314, 1314, 1314},
		{"HANDLER t READ FIRST", 1314, 1314, 1314, 1314},
		{"HANDLER t CLOSE", 1314, 1314, 1314, 1314},
		{"SHOW TABLES", 1415, 1415, 0, 0},
		{"EXPLAIN SELECT 1", 1415, 1415, 0, 0},
		{"TABLE t", 1415, 1415, 0, 0},
		{"VALUES ROW(1)", 1415, 1415, 0, 0},
		{"TABLE t INTO @a", 0, 0, 0, 0},
		{"SELECT 1 INTO @a", 0, 0, 0, 0},
		{"SELECT 1 INTO OUTFILE 'x'", 0, 0, 0, 0},
		{"KILL 1", 0, 0, 0, 0},
		{"SHUTDOWN", 0, 0, 0, 0},
		{"DO 1", 0, 0, 0, 0},
		{"CALL q()", 0, 0, 0, 0},
		{"IMPORT TABLE FROM 'a.sdi'", 0, 0, 0, 0},
		{"CLONE LOCAL DATA DIRECTORY 'a'", 0, 0, 0, 0},
		{"SET RESOURCE GROUP g", 0, 0, 0, 0},
		{"BINLOG 'x'", 1064, 1064, 1064, 1064},
	}
	for _, c := range cases {
		for i, sql := range []string{
			"CREATE FUNCTION f() RETURNS INT BEGIN " + c.stmt + "; RETURN 1; END",
			"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN " + c.stmt + "; END",
			"CREATE PROCEDURE p() BEGIN " + c.stmt + "; END",
			"CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO BEGIN " + c.stmt + "; END",
		} {
			code := []int{c.fn, c.trg, c.proc, c.evt}[i]
			_, err := Parse(sql)
			if code == 0 && err != nil {
				t.Errorf("Parse(%q) refused what MySQL takes: %v", sql, err)
			}
			if code != 0 && err == nil {
				t.Errorf("Parse(%q) took what MySQL refuses with %d", sql, code)
			}
		}
	}
}

// TestStoredProgramFunctionBodies holds a function's body to a RETURN somewhere in it (1320) and
// SET autocommit to the system variable. code is what mysql:8.0.46 answers to PREPARE of the
// statement, 0 when the compat judge reads that as parsed (compat.PrepareVerdict).
func TestStoredProgramFunctionBodies(t *testing.T) {
	cases := []struct {
		sql  string
		code int
	}{
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN DECLARE autocommit INT; SET autocommit = 1; RETURN 1; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN DECLARE autocommit INT; SET @@autocommit = 1; RETURN 1; END", code: 1445},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET GLOBAL sql_mode = '', autocommit = 1; RETURN 1; END", code: 1445},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN DECLARE EXIT HANDLER FOR 1 RETURN 1; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN END", code: 1320},
		{sql: "CREATE FUNCTION f() RETURNS INT SELECT 1 INTO @a", code: 1320},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SIGNAL SQLSTATE '45000'; END", code: 1320},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN CALL p(); RETURN 1; END"},
		{sql: "CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW CALL p()"},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET @'@autocommit' = 1; RETURN 1; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET @\"@@autocommit\" = 1; RETURN 1; END"},
		{sql: "CREATE FUNCTION f(`@@autocommit` INT) RETURNS INT BEGIN SET `@@autocommit` = 1; RETURN 1; END"},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET @@autocommit = 1; RETURN 1; END", code: 1445},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET @@session.autocommit = 1; RETURN 1; END", code: 1445},
		{sql: "CREATE FUNCTION f() RETURNS INT BEGIN SET @a = 1, @@autocommit = 1; RETURN 1; END", code: 1445},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 && err != nil {
				t.Fatalf("Parse refused what MySQL takes: %v", err)
			}
			if c.code != 0 && err == nil {
				t.Fatalf("Parse took what MySQL refuses with %d", c.code)
			}
		})
	}
}
