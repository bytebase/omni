package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func singleStmt(t *testing.T, sql string) nodes.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q): expected 1 stmt, got %d", sql, len(list.Items))
	}
	raw := list.Items[0].(*nodes.RawStmt)
	return raw.Stmt
}

// findColumnConstraint walks the first CreateStmt column looking for a
// Constraint of type CONSTR_UNIQUE and returns it.
func firstColUniqueConstraint(t *testing.T, stmt nodes.Node) *nodes.Constraint {
	t.Helper()
	cs, ok := stmt.(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected CreateStmt, got %T", stmt)
	}
	if cs.TableElts == nil || len(cs.TableElts.Items) == 0 {
		t.Fatalf("no table elements")
	}
	col, ok := cs.TableElts.Items[0].(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("elem[0] not ColumnDef: %T", cs.TableElts.Items[0])
	}
	if col.Constraints == nil {
		t.Fatalf("no constraints on column")
	}
	for _, n := range col.Constraints.Items {
		c, ok := n.(*nodes.Constraint)
		if ok && c.Contype == nodes.CONSTR_UNIQUE {
			return c
		}
	}
	t.Fatalf("no UNIQUE constraint found")
	return nil
}

func TestUniqueNullsNotDistinct(t *testing.T) {
	t.Run("column UNIQUE NULLS NOT DISTINCT", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE TABLE t (i int UNIQUE NULLS NOT DISTINCT, x text)")
		c := firstColUniqueConstraint(t, stmt)
		if !c.NullsNotDistinct {
			t.Fatalf("expected NullsNotDistinct=true")
		}
	})

	t.Run("column UNIQUE NULLS DISTINCT", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE TABLE t (i int UNIQUE NULLS DISTINCT)")
		c := firstColUniqueConstraint(t, stmt)
		if c.NullsNotDistinct {
			t.Fatalf("expected NullsNotDistinct=false")
		}
	})

	t.Run("column UNIQUE (baseline, no NULLS clause)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE TABLE t (i int UNIQUE)")
		c := firstColUniqueConstraint(t, stmt)
		if c.NullsNotDistinct {
			t.Fatalf("expected default NullsNotDistinct=false")
		}
	})

	t.Run("table-level UNIQUE NULLS NOT DISTINCT (c)", func(t *testing.T) {
		parseOK(t, "CREATE TABLE t (c int, UNIQUE NULLS NOT DISTINCT (c))")
	})

	t.Run("CREATE UNIQUE INDEX NULLS NOT DISTINCT", func(t *testing.T) {
		parseOK(t, "CREATE UNIQUE INDEX i ON t (c) NULLS NOT DISTINCT")
	})

	t.Run("CREATE UNIQUE INDEX NULLS DISTINCT", func(t *testing.T) {
		parseOK(t, "CREATE UNIQUE INDEX i ON t (c) NULLS DISTINCT")
	})

	// Regression-sanity: NULLS FIRST/LAST in ORDER BY must still work
	// (exercises the unchanged NULLS_LA reclassification path).
	t.Run("ORDER BY x NULLS FIRST (baseline)", func(t *testing.T) {
		parseOK(t, "SELECT x FROM t ORDER BY x NULLS FIRST")
	})

	t.Run("ORDER BY x NULLS LAST (baseline)", func(t *testing.T) {
		parseOK(t, "SELECT x FROM t ORDER BY x NULLS LAST")
	})
}
