package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestTableFunctionInFrom covers the #tableFunction relation primary
// (unnest(...) AS u) parsed without a LATERAL prefix — a previously dropped
// construct (the args used to be silently discarded).
func TestTableFunctionInFrom(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM t, unnest(t.arr) AS u")
	if len(stmt.From) != 2 {
		t.Fatalf("from = %d, want 2", len(stmt.From))
	}
	tf, ok := stmt.From[1].(*ast.TableFunctionRef)
	if !ok {
		t.Fatalf("from[1] = %T, want *ast.TableFunctionRef", stmt.From[1])
	}
	if tf.Lateral {
		t.Error("Lateral = true, want false (no LATERAL prefix)")
	}
	if tf.Alias != "u" {
		t.Errorf("alias = %q, want u", tf.Alias)
	}
	if tf.Call == nil || tf.Call.Name == nil ||
		tf.Call.Name.Parts[len(tf.Call.Name.Parts)-1] != "unnest" {
		t.Fatalf("call = %+v, want unnest(...)", tf.Call)
	}
	if len(tf.Call.Args) != 1 {
		t.Errorf("call args = %d, want 1 (t.arr)", len(tf.Call.Args))
	}
}

// TestTableFunctionColumnAliases covers the AS u(col) column-alias list form.
func TestTableFunctionColumnAliases(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM t, unnest(t.arr) AS u(elem)")
	tf := stmt.From[1].(*ast.TableFunctionRef)
	if got := tf.ColumnAliases; len(got) != 1 || got[0] != "elem" {
		t.Errorf("columnAliases = %v, want [elem]", got)
	}
}

// TestLateralUnnestComma covers LATERAL in the comma relation list.
func TestLateralUnnestComma(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT t.id, u.unnest FROM t, LATERAL unnest(t.arr) AS u")
	if len(stmt.From) != 2 {
		t.Fatalf("from = %d, want 2", len(stmt.From))
	}
	tf, ok := stmt.From[1].(*ast.TableFunctionRef)
	if !ok {
		t.Fatalf("from[1] = %T, want *ast.TableFunctionRef", stmt.From[1])
	}
	if !tf.Lateral {
		t.Error("Lateral = false, want true")
	}
}

// TestLateralUnnestInParenList covers LATERAL after a comma inside a
// parenthesized relation list: (t, LATERAL unnest(t.arr) AS u). The comma list
// folds into a cross JoinClause whose right side is the lateral table function.
func TestLateralUnnestInParenList(t *testing.T) {
	stmt := mustParseSelect(t, "SELECT * FROM (t, LATERAL unnest(t.arr) AS u)")
	if len(stmt.From) != 1 {
		t.Fatalf("from = %d, want 1", len(stmt.From))
	}
	join, ok := stmt.From[0].(*ast.JoinClause)
	if !ok {
		t.Fatalf("from[0] = %T, want *ast.JoinClause", stmt.From[0])
	}
	tf, ok := join.Right.(*ast.TableFunctionRef)
	if !ok {
		t.Fatalf("join.Right = %T, want *ast.TableFunctionRef", join.Right)
	}
	if !tf.Lateral {
		t.Error("Lateral = false, want true")
	}
}

// TestLateralUnnestJoin covers LATERAL on the right side of a JOIN.
func TestLateralUnnestJoin(t *testing.T) {
	join := mustParseJoin(t, "SELECT * FROM t INNER JOIN LATERAL unnest(t.tags) AS g ON TRUE")
	if join.Type != ast.JoinInner {
		t.Errorf("join.Type = %v, want JoinInner", join.Type)
	}
	tf, ok := join.Right.(*ast.TableFunctionRef)
	if !ok {
		t.Fatalf("join.Right = %T, want *ast.TableFunctionRef", join.Right)
	}
	if !tf.Lateral {
		t.Error("Lateral = false, want true")
	}
	if tf.Alias != "g" {
		t.Errorf("alias = %q, want g", tf.Alias)
	}
}
