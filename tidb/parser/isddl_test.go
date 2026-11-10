package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TestIsDDLOnParsed checks ast.IsDDL end to end from the parser, especially the
// statements whose omni AST shape differs from pingcap's (DROP VIEW, ADMIN
// CLEANUP TABLE LOCK, ADMIN REPAIR TABLE), so the prior-backup gate's DDL/DML
// split is preserved on the actual parsed nodes.
func TestIsDDLOnParsed(t *testing.T) {
	ddl := []string{
		"DROP VIEW v",
		"DROP VIEW IF EXISTS v",
		"ADMIN CLEANUP TABLE LOCK t",
		"ADMIN REPAIR TABLE t CREATE TABLE t (a int)",
		"REPAIR TABLE t",
		"OPTIMIZE TABLE t",
		"LOCK TABLES t READ",
		"UNLOCK TABLES",
		"CREATE TABLE t (a int)",
		"CREATE VIEW v AS SELECT 1",
		"DROP TABLE t",
		"RENAME TABLE a TO b",
		"TRUNCATE TABLE t",
	}
	for _, sql := range ddl {
		n := parseOneStmt(t, sql)
		if !ast.IsDDL(n) {
			t.Errorf("IsDDL(%q) [%T] = false, want true", sql, n)
		}
	}

	notDDL := []string{
		"ADMIN CHECK TABLE t",
		"ADMIN SHOW DDL",
		"SELECT 1",
		"INSERT INTO t VALUES (1)",
		"UPDATE t SET a = 1",
		"DELETE FROM t",
		"SET @x = 1",
		"SHOW TABLES",
		"GRANT SELECT ON `d`.* TO 'u'@'%'",
	}
	for _, sql := range notDDL {
		n := parseOneStmt(t, sql)
		if ast.IsDDL(n) {
			t.Errorf("IsDDL(%q) [%T] = true, want false", sql, n)
		}
	}
}
