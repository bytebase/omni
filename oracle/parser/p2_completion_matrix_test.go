package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

var p2CompletionPositiveSQL = []string{
	"ALTER AUDIT POLICY pol ADD PRIVILEGES CREATE SESSION",
	"ALTER CLUSTER c SIZE 1024",
	"ALTER FLASHBACK ARCHIVE fba MODIFY RETENTION 1 DAY",
	"ALTER FUNCTION f COMPILE",
	"ALTER INDEX idx REBUILD",
	"ALTER MATERIALIZED VIEW LOG ON t ADD ROWID",
	"ALTER MATERIALIZED VIEW mv COMPILE",
	"ALTER PACKAGE pkg COMPILE",
	"ALTER PROCEDURE p COMPILE",
	"ALTER PROFILE prof LIMIT SESSIONS_PER_USER 2",
	"ALTER ROLE r IDENTIFIED BY pwd",
	"ALTER ROLLBACK SEGMENT rbs ONLINE",
	"ALTER SEQUENCE seq INCREMENT BY 1",
	"ALTER SESSION SET nls_date_format = 'YYYY-MM-DD'",
	"ALTER SYNONYM syn COMPILE",
	"ALTER SYSTEM SWITCH LOGFILE",
	"ALTER TABLE t ADD (c NUMBER)",
	"ALTER TRIGGER trg COMPILE",
	"ALTER TYPE typ COMPILE",
	"ALTER USER u IDENTIFIED BY p",
	"ALTER VIEW v COMPILE",
	"ANALYZE TABLE t COMPUTE STATISTICS",
	"ASSOCIATE STATISTICS WITH COLUMNS t.c USING stat_type",
	"AUDIT SELECT ON t",
	"AUDIT POLICY pol",
	"CALL proc()",
	"COMMENT ON TABLE t IS 'x'",
	"COMMIT",
	"CREATE AUDIT POLICY pol PRIVILEGES CREATE SESSION",
	"CREATE CLUSTER c (id NUMBER)",
	"CREATE CONTEXT ctx USING pkg",
	"CREATE DIRECTORY dir AS '/tmp'",
	"CREATE EDITION ed",
	"CREATE FLASHBACK ARCHIVE fba TABLESPACE ts QUOTA 1G RETENTION 1 DAY",
	"CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;",
	"CREATE INDEX idx ON t (c)",
	"CREATE INDEXTYPE my_itype FOR my_op(NUMBER) USING my_type WITH ARRAY DML (NUMBER)",
	"CREATE MATERIALIZED VIEW LOG ON t",
	"CREATE MATERIALIZED VIEW mv AS SELECT c FROM t",
	"CREATE PACKAGE pkg AS PROCEDURE p; END;",
	"CREATE PACKAGE BODY pkg AS PROCEDURE p IS BEGIN NULL; END; END;",
	"CREATE PROCEDURE p IS BEGIN NULL; END;",
	"CREATE PROFILE prof LIMIT SESSIONS_PER_USER 2",
	"CREATE ROLE r",
	"CREATE ROLLBACK SEGMENT rbs",
	"CREATE SCHEMA AUTHORIZATION s CREATE TABLE t (c NUMBER)",
	"CREATE SEQUENCE seq",
	"CREATE SYNONYM syn FOR t",
	"CREATE TABLE t (c NUMBER)",
	"CREATE TRIGGER trg BEFORE INSERT ON t BEGIN NULL; END;",
	"CREATE TYPE typ AS OBJECT (c NUMBER)",
	"CREATE USER u IDENTIFIED BY p",
	"CREATE VIEW v AS SELECT c FROM t",
	"DELETE FROM t WHERE c = 1",
	"DISASSOCIATE STATISTICS FROM COLUMNS t.c",
	"DROP AUDIT POLICY pol",
	"DROP CLUSTER c",
	"DROP CONTEXT ctx",
	"DROP DIRECTORY dir",
	"DROP EDITION ed",
	"DROP FLASHBACK ARCHIVE fba",
	"DROP FUNCTION f",
	"DROP INDEX idx",
	"DROP INDEXTYPE my_itype",
	"DROP MATERIALIZED VIEW LOG ON t",
	"DROP MATERIALIZED VIEW mv",
	"DROP PACKAGE pkg",
	"DROP PROCEDURE p",
	"DROP PROFILE prof",
	"DROP ROLE r",
	"DROP ROLLBACK SEGMENT rbs",
	"DROP SEQUENCE seq",
	"DROP SYNONYM syn",
	"DROP TABLE t",
	"DROP TRIGGER trg",
	"DROP TYPE BODY typ",
	"DROP TYPE typ",
	"DROP USER u",
	"DROP VIEW v",
	"EXPLAIN PLAN FOR SELECT 1 FROM dual",
	"FLASHBACK TABLE t TO BEFORE DROP",
	"GRANT SELECT ON t TO u",
	"INSERT INTO t (c) VALUES (1)",
	"LOCK TABLE t IN EXCLUSIVE MODE",
	"MERGE INTO t USING s ON (t.id = s.id) WHEN MATCHED THEN UPDATE SET c = s.c",
	"NOAUDIT SELECT ON t",
	"NOAUDIT POLICY pol",
	"PURGE TABLE t",
	"RENAME old_name TO new_name",
	"REVOKE SELECT ON t FROM u",
	"ROLLBACK",
	"SAVEPOINT sp",
	"SELECT 1 FROM dual",
	"SET CONSTRAINTS ALL IMMEDIATE",
	"SET ROLE ALL",
	"SET TRANSACTION READ ONLY",
	"TRUNCATE CLUSTER c",
	"TRUNCATE TABLE t",
	"UPDATE t SET c = 1",
}

func TestP2CompletionPositive(t *testing.T) {
	for _, sql := range p2CompletionPositiveSQL {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() == 0 {
				t.Fatalf("expected at least one statement")
			}
		})
	}
}

func TestP2CompletionNegative(t *testing.T) {
	tests := []string{
		"CREATE TABLE bad_t (c NUMBER",
		"CREATE INDEX bad_idx ON",
		"DROP TABLE",
		"SELECT 1 +",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			ParseShouldFail(t, sql)
		})
	}
}

func TestP2CompletionLoc(t *testing.T) {
	for _, sql := range p2CompletionPositiveSQL {
		t.Run(sql, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			for _, item := range result.Items {
				raw := item.(*ast.RawStmt)
				loc := ast.NodeLoc(raw.Stmt)
				if loc.Start < 0 || loc.End <= loc.Start {
					t.Fatalf("invalid statement Loc=%+v", loc)
				}
			}
		})
	}
}
