package analysis

import (
	"testing"

	"github.com/bytebase/omni/mssql"
)

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		statement string
		valid     bool
		allQuery  bool
	}{
		{"SELECT id FROM t", true, true},
		{"SELECT 1; SELECT 2;", true, true},
		{"SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED\nSELECT id FROM t", true, false},
		{"set nocount on; select id from t", true, false},
		{"SET NOCOUNT OFF; SET TRANSACTION ISOLATION LEVEL READ COMMITTED", true, false},
		{"SET NOCOUNT, FMTONLY ON; SELECT id FROM t", false, false},
		{"SET STATISTICS XML ON; SELECT id FROM t", false, false},
		{"SET SHOWPLAN_XML ON", false, false},
		{"DECLARE @n int = 1; SELECT @n", false, false},
		{"SET @n = 1", false, false},
		{"SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED; DELETE FROM t", false, false},
		{"SELECT id INTO #t FROM t", false, false},
		{"SELECT id INTO t2 FROM t UNION SELECT id FROM t1", false, false},
		{"", false, false},
		{";", false, false},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		valid, allQuery := ValidateQuery(stmts)
		if valid != tc.valid || allQuery != tc.allQuery {
			t.Errorf("%q: got (%v, %v), want (%v, %v)", tc.statement, valid, allQuery, tc.valid, tc.allQuery)
		}
	}
}

func TestStatementTypeOf(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementType
	}{
		{"ALTER DATABASE d SET RECOVERY FULL", AlterDatabase},
		{"ALTER INDEX i ON t REBUILD", AlterIndex},
		{"ALTER TABLE t ADD b INT", AlterTable},
		{"CREATE DATABASE d", CreateDatabase},
		{"CREATE INDEX i ON t (a)", CreateIndex},
		{"CREATE SCHEMA s", CreateSchema},
		{"CREATE TABLE t (a INT)", CreateTable},
		{"CREATE VIEW v AS SELECT 1 AS a", CreateView},
		{"DROP DATABASE d", DropDatabase},
		{"DROP INDEX i ON t", DropIndex},
		{"DROP SCHEMA s", DropSchema},
		{"DROP TABLE t", DropTable},
		{"DROP VIEW v", DropView},
		{"DROP PROCEDURE p", StatementTypeUnspecified},
		{"TRUNCATE TABLE t", Truncate},
		{"DELETE FROM t", Delete},
		{"INSERT INTO t VALUES (1)", Insert},
		{"UPDATE t SET a = 1", Update},
		{"EXEC p", Execute},
		{"SELECT 1", StatementTypeUnspecified},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		if got := StatementTypeOf(stmts[0].AST); got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.statement, got, tc.want)
		}
	}
}
