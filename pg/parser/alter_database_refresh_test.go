package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestAlterDatabaseRefreshCollation(t *testing.T) {
	t.Run("REFRESH COLLATION VERSION", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE postgres REFRESH COLLATION VERSION")
		rc, ok := stmt.(*nodes.AlterDatabaseRefreshCollStmt)
		if !ok {
			t.Fatalf("expected AlterDatabaseRefreshCollStmt, got %T", stmt)
		}
		if rc.Dbname != "postgres" {
			t.Fatalf("expected Dbname=postgres, got %q", rc.Dbname)
		}
	})

	// Regression-sanity.
	t.Run("SET timezone (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE d SET timezone = 'UTC'")
		if _, ok := stmt.(*nodes.AlterDatabaseSetStmt); !ok {
			t.Fatalf("expected AlterDatabaseSetStmt, got %T", stmt)
		}
	})

	t.Run("RENAME TO (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER DATABASE d RENAME TO d2")
		if _, ok := stmt.(*nodes.RenameStmt); !ok {
			t.Fatalf("expected RenameStmt, got %T", stmt)
		}
	})
}
