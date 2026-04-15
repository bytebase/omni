package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestAlterViewRenameColumn(t *testing.T) {
	t.Run("RENAME COLUMN a TO b", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER VIEW v RENAME COLUMN a TO b")
		rn, ok := stmt.(*nodes.RenameStmt)
		if !ok {
			t.Fatalf("expected RenameStmt, got %T", stmt)
		}
		if rn.RenameType != nodes.OBJECT_COLUMN {
			t.Fatalf("expected RenameType=OBJECT_COLUMN, got %v", rn.RenameType)
		}
		if rn.RelationType != nodes.OBJECT_VIEW {
			t.Fatalf("expected RelationType=OBJECT_VIEW, got %v", rn.RelationType)
		}
		if rn.Subname != "a" || rn.Newname != "b" {
			t.Fatalf("expected a->b, got %q->%q", rn.Subname, rn.Newname)
		}
		if rn.Relation == nil || rn.Relation.Relname != "v" {
			t.Fatalf("expected relation v, got %+v", rn.Relation)
		}
	})

	// Regression-sanity: RENAME TO (the view itself) still works.
	t.Run("RENAME TO (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER VIEW v RENAME TO v2")
		rn, ok := stmt.(*nodes.RenameStmt)
		if !ok {
			t.Fatalf("expected RenameStmt, got %T", stmt)
		}
		if rn.RenameType != nodes.OBJECT_VIEW {
			t.Fatalf("expected OBJECT_VIEW, got %v", rn.RenameType)
		}
		if rn.Newname != "v2" {
			t.Fatalf("expected v2, got %q", rn.Newname)
		}
	})

	// Regression-sanity: ALTER COLUMN SET DEFAULT still works.
	t.Run("ALTER COLUMN SET DEFAULT (baseline)", func(t *testing.T) {
		parseOK(t, "ALTER VIEW v ALTER COLUMN a SET DEFAULT 1")
	})
}
