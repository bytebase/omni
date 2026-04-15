package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestAlterPolicyRename(t *testing.T) {
	t.Run("RENAME TO", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER POLICY p1 ON t RENAME TO p2")
		rn, ok := stmt.(*nodes.RenameStmt)
		if !ok {
			t.Fatalf("expected RenameStmt, got %T", stmt)
		}
		if rn.RenameType != nodes.OBJECT_POLICY {
			t.Fatalf("expected OBJECT_POLICY, got %v", rn.RenameType)
		}
		if rn.Subname != "p1" {
			t.Fatalf("expected Subname=p1, got %q", rn.Subname)
		}
		if rn.Newname != "p2" {
			t.Fatalf("expected Newname=p2, got %q", rn.Newname)
		}
		if rn.Relation == nil || rn.Relation.Relname != "t" {
			t.Fatalf("expected relation t, got %+v", rn.Relation)
		}
	})

	t.Run("RENAME TO on schema.table", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER POLICY p1 ON s1.t RENAME TO p2")
		rn, ok := stmt.(*nodes.RenameStmt)
		if !ok {
			t.Fatalf("expected RenameStmt, got %T", stmt)
		}
		if rn.Relation == nil || rn.Relation.Schemaname != "s1" || rn.Relation.Relname != "t" {
			t.Fatalf("expected s1.t, got %+v", rn.Relation)
		}
	})

	// Regression-sanity: existing ALTER POLICY forms still produce AlterPolicyStmt.
	t.Run("USING (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER POLICY p1 ON t USING (a > 0)")
		if _, ok := stmt.(*nodes.AlterPolicyStmt); !ok {
			t.Fatalf("expected AlterPolicyStmt, got %T", stmt)
		}
	})

	t.Run("WITH CHECK (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "ALTER POLICY p1 ON t WITH CHECK (a > 0)")
		if _, ok := stmt.(*nodes.AlterPolicyStmt); !ok {
			t.Fatalf("expected AlterPolicyStmt, got %T", stmt)
		}
	})
}
