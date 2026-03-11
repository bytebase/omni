package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestIntegrationAllStatementTypes verifies every statement type routes through
// the top-level Parse() entry point and produces the expected AST node.
func TestIntegrationAllStatementTypes(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantTyp string // Go type name we expect for the top-level stmt
	}{
		// DML
		{"SELECT", "SELECT 1 FROM dual", "*ast.SelectStmt"},
		{"WITH_SELECT", "WITH cte AS (SELECT 1 FROM dual) SELECT * FROM cte", "*ast.SelectStmt"},
		{"INSERT", "INSERT INTO t (a) VALUES (1)", "*ast.InsertStmt"},
		{"UPDATE", "UPDATE t SET a = 1 WHERE b = 2", "*ast.UpdateStmt"},
		{"DELETE", "DELETE FROM t WHERE a = 1", "*ast.DeleteStmt"},
		{"MERGE", "MERGE INTO t USING s ON (t.id = s.id) WHEN MATCHED THEN UPDATE SET t.x = s.x", "*ast.MergeStmt"},

		// DDL - CREATE
		{"CREATE_TABLE", "CREATE TABLE t (id NUMBER)", "*ast.CreateTableStmt"},
		{"CREATE_INDEX", "CREATE INDEX idx ON t (a)", "*ast.CreateIndexStmt"},
		{"CREATE_VIEW", "CREATE VIEW v AS SELECT 1 FROM dual", "*ast.CreateViewStmt"},
		{"CREATE_SEQUENCE", "CREATE SEQUENCE seq START WITH 1", "*ast.CreateSequenceStmt"},
		{"CREATE_SYNONYM", "CREATE SYNONYM s FOR t", "*ast.CreateSynonymStmt"},
		{"CREATE_TYPE", "CREATE TYPE my_type AS OBJECT (id NUMBER)", "*ast.CreateTypeStmt"},
		{"CREATE_PROCEDURE", "CREATE PROCEDURE p IS BEGIN NULL; END;", "*ast.CreateProcedureStmt"},
		{"CREATE_FUNCTION", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;", "*ast.CreateFunctionStmt"},
		{"CREATE_PACKAGE", "CREATE PACKAGE pkg IS END pkg;", "*ast.CreatePackageStmt"},
		{"CREATE_TRIGGER", "CREATE TRIGGER trg BEFORE INSERT ON t BEGIN NULL; END;", "*ast.CreateTriggerStmt"},

		// DDL - ALTER
		{"ALTER_TABLE", "ALTER TABLE t ADD (col NUMBER)", "*ast.AlterTableStmt"},
		{"ALTER_INDEX", "ALTER INDEX idx REBUILD", "*ast.AdminDDLStmt"},
		{"ALTER_SEQUENCE", "ALTER SEQUENCE seq INCREMENT BY 2", "*ast.AdminDDLStmt"},
		{"ALTER_SESSION", "ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY-MM-DD'", "*ast.AlterSessionStmt"},

		// DDL - DROP
		{"DROP_TABLE", "DROP TABLE t", "*ast.DropStmt"},
		{"DROP_INDEX", "DROP INDEX idx", "*ast.DropStmt"},
		{"DROP_VIEW", "DROP VIEW v", "*ast.DropStmt"},

		// Transaction
		{"COMMIT", "COMMIT", "*ast.CommitStmt"},
		{"ROLLBACK", "ROLLBACK", "*ast.RollbackStmt"},
		{"SAVEPOINT", "SAVEPOINT sp1", "*ast.SavepointStmt"},
		{"SET_TRANSACTION", "SET TRANSACTION READ ONLY", "*ast.SetTransactionStmt"},

		// Privilege
		{"GRANT", "GRANT SELECT ON t TO u", "*ast.GrantStmt"},
		{"REVOKE", "REVOKE SELECT ON t FROM u", "*ast.RevokeStmt"},

		// Utility
		{"TRUNCATE", "TRUNCATE TABLE t", "*ast.TruncateStmt"},
		{"COMMENT", "COMMENT ON TABLE t IS 'test'", "*ast.CommentStmt"},
		{"EXPLAIN_PLAN", "EXPLAIN PLAN FOR SELECT 1 FROM dual", "*ast.ExplainPlanStmt"},

		// PL/SQL
		{"PLSQL_BLOCK", "BEGIN NULL; END;", "*ast.PLSQLBlock"},
		{"PLSQL_DECLARE", "DECLARE v NUMBER; BEGIN v := 1; END;", "*ast.PLSQLBlock"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseAndCheck(t, tc.sql)
			if result.Len() < 1 {
				t.Fatal("expected at least 1 statement")
			}
			raw := result.Items[0].(*ast.RawStmt)
			if raw.Stmt == nil {
				t.Fatal("expected non-nil Stmt")
			}
			got := typeName(raw.Stmt)
			if got != tc.wantTyp {
				t.Errorf("expected %s, got %s", tc.wantTyp, got)
			}
		})
	}
}

// TestIntegrationMultiStatement tests parsing multiple semicolon-separated statements.
func TestIntegrationMultiStatement(t *testing.T) {
	sql := "SELECT 1 FROM dual; INSERT INTO t (a) VALUES (1); COMMIT"
	result := ParseAndCheck(t, sql)
	if result.Len() != 3 {
		t.Fatalf("expected 3 statements, got %d", result.Len())
	}

	// Verify types
	types := []string{"*ast.SelectStmt", "*ast.InsertStmt", "*ast.CommitStmt"}
	for i, want := range types {
		raw := result.Items[i].(*ast.RawStmt)
		got := typeName(raw.Stmt)
		if got != want {
			t.Errorf("stmt %d: expected %s, got %s", i, want, got)
		}
	}
}

func typeName(n ast.Node) string {
	return fmt.Sprintf("%T", n)
}
