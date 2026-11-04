package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// exprTableNames renders a table-expression list the way tableNames renders a
// TableRef list, following a join to both sides, so a list that lost or
// reordered an entry reads differently.
func exprTableNames(exprs []ast.TableExpr) []string {
	var out []string
	var walk func(ast.TableExpr)
	walk = func(e ast.TableExpr) {
		switch v := e.(type) {
		case *ast.TableRef:
			out = append(out, tableNames([]*ast.TableRef{v})...)
		case *ast.JoinClause:
			walk(v.Left)
			walk(v.Right)
		}
	}
	for _, e := range exprs {
		walk(e)
	}
	return out
}

// litInt reads an integer literal's value, failing on anything else: a test
// that compared the node instead would pass on a literal of the wrong type.
func litInt(t *testing.T, e ast.ExprNode) int64 {
	t.Helper()
	lit, ok := e.(*ast.IntLit)
	if !ok {
		t.Fatalf("expression is %T, want *ast.IntLit", e)
	}
	return lit.Value
}

func nodeTypeName(n any) string { return fmt.Sprintf("%T", n) }

// What TRACE leaves in the tree.
//
// The corpus records only whether a statement parses, which is why a family
// can close every one of its rows and still ship a defect: a modifier that
// parses into the wrong field, a choice that reads as its neighbour, a
// wrapper whose inner statement comes out empty or wrong. For TRACE the
// field that matters most is TracePlan, because it is what says whether the
// server runs the statement: a TRACE PLAN read as a plain TRACE would have
// the catalog apply a SET the server never performed, and a plain TRACE read
// as a TRACE PLAN would have review stop checking a statement that does run.
// Nothing in the corpus says either must not happen.

func parseOneTrace(t *testing.T, sql string) *ast.TraceStmt {
	t.Helper()
	stmt, ok := parseOneStmt(t, sql).(*ast.TraceStmt)
	if !ok {
		t.Fatalf("Parse(%q) gave a statement of another kind", sql)
	}
	return stmt
}

// TestTraceModifierFields pins each arm's fields by value. The four arms
// differ only in these three fields, so a reading that confused two of them
// would still parse every corpus row.
func TestTraceModifierFields(t *testing.T) {
	tests := []struct {
		sql       string
		format    string
		tracePlan bool
		target    string
	}{
		// A bare TRACE sets no modifier. Upstream's own node defaults
		// Format to "row" here; this one leaves it empty, so a consumer can
		// tell TRACE from TRACE FORMAT = 'row'.
		{sql: "trace select 1"},
		{sql: "trace format = 'row' select 1", format: "row"},
		{sql: "trace format = 'json' select 1", format: "json"},
		// The name is kept as written, in either quote, and is not checked.
		{sql: "trace format = 'ROW' select 1", format: "ROW"},
		{sql: `trace format = "json" select 1`, format: "json"},
		{sql: "trace format = 'nosuchformat' select 1", format: "nosuchformat"},
		{sql: "trace format = '' select 1", format: ""},
		// PLAN sets TracePlan and no format.
		{sql: "trace plan select 1", tracePlan: true},
		{sql: "trace plan target = 'estimation' select 1", tracePlan: true, target: "estimation"},
		{sql: "trace plan target = 'arandomstring' select 1", tracePlan: true, target: "arandomstring"},
		{sql: "trace plan target = '' select 1", tracePlan: true, target: ""},
	}
	for _, c := range tests {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseOneTrace(t, c.sql)
			if stmt.Format != c.format {
				t.Errorf("Format = %q, want %q", stmt.Format, c.format)
			}
			if stmt.TracePlan != c.tracePlan {
				t.Errorf("TracePlan = %v, want %v", stmt.TracePlan, c.tracePlan)
			}
			if stmt.Target != c.target {
				t.Errorf("Target = %q, want %q", stmt.Target, c.target)
			}
		})
	}
}

// TestTraceCarriesInnerStatement pins the statement each TRACE holds, by the
// node it is and by what that node says. A wrapper that dropped its inner
// statement, or kept the right type with the wrong contents, parses every
// corpus row either way; the second is the failure a nil check does not see.
func TestTraceCarriesInnerStatement(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace select c1 from t1")
		sel, ok := stmt.Stmt.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.SelectStmt", stmt.Stmt)
		}
		if len(sel.TargetList) != 1 {
			t.Fatalf("the traced SELECT has %d select expressions, want 1", len(sel.TargetList))
		}
		col, ok := sel.TargetList[0].(*ast.ColumnRef)
		if !ok || col.Column != "c1" {
			t.Errorf("the traced SELECT selects %#v, want column c1", sel.TargetList[0])
		}
		if got := exprTableNames(sel.From); join(got) != join([]string{"t1"}) {
			t.Errorf("the traced SELECT reads %s, want [t1]", join(got))
		}
	})

	t.Run("union with limit", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace select c1 from t1 union (select c2 from t2) limit 1, 1")
		sel, ok := stmt.Stmt.(*ast.SelectStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.SelectStmt", stmt.Stmt)
		}
		if sel.SetOp != ast.SetOpUnion {
			t.Errorf("SetOp = %v, want UNION", sel.SetOp)
		}
		// The LIMIT belongs to the set operation, not to either side, which
		// is the boundary a wrapper that re-parsed the tail would lose.
		if sel.Limit == nil {
			t.Fatal("the traced UNION has no LIMIT")
		}
		if got := litInt(t, sel.Limit.Count); got != 1 {
			t.Errorf("LIMIT count = %d, want 1", got)
		}
		if got := litInt(t, sel.Limit.Offset); got != 1 {
			t.Errorf("LIMIT offset = %d, want 1", got)
		}
	})

	t.Run("update", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace format = 'json' update t set id = id + 1 order by id desc")
		upd, ok := stmt.Stmt.(*ast.UpdateStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.UpdateStmt", stmt.Stmt)
		}
		if len(upd.SetList) != 1 || upd.SetList[0].Column.Column != "id" {
			t.Errorf("the traced UPDATE assigns %#v, want one assignment to id", upd.SetList)
		}
		if len(upd.OrderBy) != 1 || !upd.OrderBy[0].Desc {
			t.Errorf("the traced UPDATE orders by %#v, want one descending item", upd.OrderBy)
		}
		// The modifier is the TRACE's, not the UPDATE's.
		if stmt.Format != "json" {
			t.Errorf("Format = %q, want json", stmt.Format)
		}
	})

	t.Run("multi-table delete", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace delete t1, t2 from t1 inner join t2 where t1.id=t2.id")
		del, ok := stmt.Stmt.(*ast.DeleteStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.DeleteStmt", stmt.Stmt)
		}
		if got := exprTableNames(del.Tables); join(got) != join([]string{"t1", "t2"}) {
			t.Errorf("the traced DELETE targets %s, want [t1|t2]", join(got))
		}
		if del.Where == nil {
			t.Error("the traced DELETE lost its WHERE clause")
		}
	})

	t.Run("insert", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace insert into t values (1), (2), (3)")
		ins, ok := stmt.Stmt.(*ast.InsertStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.InsertStmt", stmt.Stmt)
		}
		if ins.IsReplace {
			t.Error("the traced INSERT reads as a REPLACE")
		}
		if len(ins.Values) != 3 {
			t.Errorf("the traced INSERT has %d row lists, want 3", len(ins.Values))
		}
	})

	t.Run("replace", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace replace into foo values (1)")
		ins, ok := stmt.Stmt.(*ast.InsertStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.InsertStmt", stmt.Stmt)
		}
		// REPLACE and INSERT share a node, and the flag is the only thing
		// that tells them apart -- which prior backup's rule turns on.
		if !ins.IsReplace {
			t.Error("the traced REPLACE reads as an INSERT")
		}
	})

	t.Run("set", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace set foreign_key_checks = 0")
		set, ok := stmt.Stmt.(*ast.SetStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.SetStmt", stmt.Stmt)
		}
		if len(set.Assignments) != 1 {
			t.Fatalf("the traced SET has %d assignments, want 1", len(set.Assignments))
		}
		if got := set.Assignments[0].Column.Column; got != "foreign_key_checks" {
			t.Errorf("the traced SET assigns %q, want foreign_key_checks", got)
		}
	})

	t.Run("transaction statements", func(t *testing.T) {
		for _, c := range []struct {
			sql  string
			want any
		}{
			{"trace begin", &ast.BeginStmt{}},
			{"trace start transaction", &ast.BeginStmt{}},
			{"trace commit", &ast.CommitStmt{}},
			{"trace rollback", &ast.RollbackStmt{}},
			{"trace savepoint sp1", &ast.SavepointStmt{}},
			{"trace release savepoint sp1", &ast.SavepointStmt{}},
		} {
			stmt := parseOneTrace(t, c.sql)
			if gotType, wantType := nodeTypeName(stmt.Stmt), nodeTypeName(c.want); gotType != wantType {
				t.Errorf("%q carries %s, want %s", c.sql, gotType, wantType)
			}
		}
	})

	t.Run("load data", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace load data local infile '/tmp/a.csv' into table t")
		load, ok := stmt.Stmt.(*ast.LoadDataStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.LoadDataStmt", stmt.Stmt)
		}
		if load.IsXML {
			t.Error("the traced LOAD DATA reads as a LOAD XML")
		}
		if !load.Local {
			t.Error("the traced LOAD DATA lost its LOCAL modifier")
		}
		if load.Infile != "/tmp/a.csv" {
			t.Errorf("Infile = %q, want /tmp/a.csv", load.Infile)
		}
	})

	t.Run("analyze table", func(t *testing.T) {
		stmt := parseOneTrace(t, "trace analyze table t index idx")
		an, ok := stmt.Stmt.(*ast.AnalyzeTableStmt)
		if !ok {
			t.Fatalf("Stmt is %T, want *ast.AnalyzeTableStmt", stmt.Stmt)
		}
		if got := tableNames(an.Tables); join(got) != join([]string{"t"}) {
			t.Errorf("the traced ANALYZE reads %s, want [t]", join(got))
		}
	})
}

// TestTraceLoc pins the statement's own byte range: the whole TRACE, from its
// first word to the last token of what it carries. A wrapper whose Loc held
// only its keyword, or only the inner statement, would anchor every finding
// on it in the wrong place.
func TestTraceLoc(t *testing.T) {
	const sql = "trace plan target = 'estimation' select 1"
	stmt := parseOneTrace(t, sql)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(sql) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(sql))
	}
	// The inner statement keeps its own range, inside the wrapper's.
	inner := stmt.Stmt.(*ast.SelectStmt)
	if got := sql[inner.Loc.Start:inner.Loc.End]; got != "select 1" {
		t.Errorf("the traced statement spans %q, want %q", got, "select 1")
	}
}
