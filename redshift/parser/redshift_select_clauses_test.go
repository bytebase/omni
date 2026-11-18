package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

// TestRedshiftSelectExcludeParse pins EXCLUDE, which r_SELECT_synopsis and
// r_SELECT_INTO write as one name or as a parenthesized list of them. Both
// forms are recorded as the result names the SELECT drops, in the order
// written, and SELECT INTO takes the same clause.
func TestRedshiftSelectExcludeParse(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		names []string
	}{
		{"SELECT * EXCLUDE a FROM t;", []string{"a"}},
		{"SELECT * EXCLUDE (a, b) FROM t;", []string{"a", "b"}},
		{"SELECT * EXCLUDE a INTO t2 FROM t;", []string{"a"}},
		{"SELECT * EXCLUDE (a) INTO t2 FROM t;", []string{"a"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := singleStmt(t, tc.sql)
			sel, ok := stmt.(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("expected SelectStmt, got %T", stmt)
			}
			if sel.ExcludeList == nil || len(sel.ExcludeList.Items) != len(tc.names) {
				t.Fatalf("ExcludeList = %#v, want %d names", sel.ExcludeList, len(tc.names))
			}
			for i, want := range tc.names {
				name, ok := sel.ExcludeList.Items[i].(*nodes.String)
				if !ok || name.Str != want {
					t.Fatalf("ExcludeList[%d] = %#v, want %q", i, sel.ExcludeList.Items[i], want)
				}
			}
		})
	}

	// The comma list is the parenthesized form's, so a bare pair is no EXCLUDE.
	if _, err := Parse("SELECT * EXCLUDE a, b FROM t;"); err == nil {
		t.Error("EXCLUDE a, b parsed, want a syntax error")
	}
	for _, sql := range []string{
		"SELECT * EXCLUDE a,",
		"SELECT * EXCLUDE (a,)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftSelectConnectByParse pins the hierarchical-query clauses
// r_SELECT_synopsis writes between WHERE and GROUP BY, START WITH optional
// before the CONNECT BY that carries the hierarchy. A relation named start or
// connect stays an alias: only START WITH and CONNECT BY begin the clause.
func TestRedshiftSelectConnectByParse(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		startWith bool
	}{
		{"SELECT * CONNECT BY 1;", false},
		{"SELECT * START WITH 1 CONNECT BY 1;", true},
		{"SELECT * FROM t START WITH 1 CONNECT BY 1;", true},
		{"SELECT * FROM t WHERE a = 1 START WITH 1 CONNECT BY 1 GROUP BY ALL;", true},
		{"SELECT * INTO t2 CONNECT BY 1;", false},
		{"SELECT * INTO t2 START WITH 1 CONNECT BY 1;", true},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := singleStmt(t, tc.sql)
			sel, ok := stmt.(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("expected SelectStmt, got %T", stmt)
			}
			if sel.ConnectBy == nil {
				t.Fatal("expected a CONNECT BY condition")
			}
			if got := sel.StartWith != nil; got != tc.startWith {
				t.Fatalf("StartWith present = %v, want %v", got, tc.startWith)
			}
			// A relation names the table it is, not the clause: these keep the
			// name as the relation's table.
			if tc.sql == "SELECT * FROM t START WITH 1 CONNECT BY 1;" {
				if sel.FromClause == nil || len(sel.FromClause.Items) != 1 {
					t.Fatalf("FromClause = %#v, want one relation", sel.FromClause)
				}
				rel, ok := sel.FromClause.Items[0].(*nodes.RangeVar)
				if !ok || rel.Relname != "t" {
					t.Fatalf("FromClause[0] = %#v, want relation t", sel.FromClause.Items[0])
				}
				if rel.Alias != nil {
					t.Fatalf("relation alias = %#v, want none: START WITH is a clause", rel.Alias)
				}
			}
		})
	}
}

// TestRedshiftSelectConnectByPriorParse pins PRIOR, the prefix operator that
// marks the parent row's side of a CONNECT BY condition. PRIOR is unreserved,
// so outside the condition it still names a column.
func TestRedshiftSelectConnectByPriorParse(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t CONNECT BY PRIOR id = parent_id;",
		"SELECT * FROM t START WITH parent_id IS NULL CONNECT BY PRIOR id = parent_id;",
		"SELECT * FROM t CONNECT BY PRIOR id > 0 AND PRIOR depth < 5;",
		"SELECT prior FROM t;",
		"SELECT * FROM prior;",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
		}
	}

	sel, ok := singleStmt(t, "SELECT * FROM t CONNECT BY PRIOR id = parent_id;").(*nodes.SelectStmt)
	if !ok {
		t.Fatal("expected SelectStmt")
	}
	eq, ok := sel.ConnectBy.(*nodes.A_Expr)
	if !ok || eq.Lexpr == nil {
		t.Fatalf("ConnectBy = %#v, want the comparison over PRIOR", sel.ConnectBy)
	}
	prior, ok := eq.Lexpr.(*nodes.A_Expr)
	if !ok || prior.Lexpr != nil || prior.Rexpr == nil {
		t.Fatalf("left side = %#v, want the prefix operator over id", eq.Lexpr)
	}
	if name := prior.Name.Items[0].(*nodes.String).Str; name != "prior" {
		t.Fatalf("prefix operator = %q, want prior", name)
	}

	// Outside a CONNECT BY condition PRIOR is a name and no operator.
	if _, err := Parse("SELECT * FROM t WHERE PRIOR id = 1;"); err == nil {
		t.Error("PRIOR as a prefix operator outside CONNECT BY parsed, want a syntax error")
	}
}

// TestRedshiftSelectConnectByNeedsClauseKeywords holds the two ways in: a bare
// START or CONNECT is a name, and the clause needs its own continuation.
func TestRedshiftSelectConnectByNeedsClauseKeywords(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t start;",
		"SELECT * FROM connect;",
	} {
		if _, err := Parse(sql); err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
	}
	for _, sql := range []string{
		"SELECT * START 1 CONNECT BY 1;",
		"SELECT * CONNECT 1;",
		"SELECT * START WITH 1;",
	} {
		if _, err := Parse(sql); err == nil {
			t.Fatalf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestRedshiftSelectGroupByAllParse pins GROUP BY ALL, which the docs write
// as a choice against the group list rather than before it, in SELECT and
// SELECT INTO.
func TestRedshiftSelectGroupByAllParse(t *testing.T) {
	for _, sql := range []string{
		"SELECT * GROUP BY ALL;",
		"SELECT * FROM t GROUP BY ALL;",
		"SELECT * INTO t2 FROM t GROUP BY ALL;",
		"SELECT * FROM t GROUP BY ALL HAVING COUNT(*) > 1;",
	} {
		t.Run(sql, func(t *testing.T) {
			stmt := singleStmt(t, sql)
			sel, ok := stmt.(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("expected SelectStmt, got %T", stmt)
			}
			if !sel.GroupAll {
				t.Fatal("expected GroupAll")
			}
			if sel.GroupClause != nil {
				t.Fatalf("GroupClause = %#v, want none beside ALL", sel.GroupClause)
			}
		})
	}

	// ALL also modifies a list, which is how PostgreSQL writes its
	// set_quantifier, so GROUP BY ALL rollup(a, b) keeps both readings.
	qualified, ok := singleStmt(t, "SELECT * FROM t GROUP BY ALL rollup(a, b), rollup(a, c);").(*nodes.SelectStmt)
	if !ok {
		t.Fatal("expected SelectStmt")
	}
	if !qualified.GroupAll || qualified.GroupClause == nil || len(qualified.GroupClause.Items) != 2 {
		t.Fatalf("GROUP BY ALL rollup(a, b), rollup(a, c) = %#v, want GroupAll over two items", qualified)
	}

	// DISTINCT still takes its list, as PostgreSQL writes it, and the two markers
	// are apart.
	for _, tc := range []struct {
		sql  string
		name string
	}{
		{"SELECT * FROM t GROUP BY DISTINCT a;", "postgresql's GROUP BY DISTINCT"},
		{"SELECT * FROM t GROUP BY a, b;", "a group list"},
	} {
		sel, ok := singleStmt(t, tc.sql).(*nodes.SelectStmt)
		if !ok {
			t.Fatalf("Parse(%q): expected SelectStmt", tc.sql)
		}
		if sel.GroupAll {
			t.Fatalf("Parse(%q): GroupAll = true, want false for %s", tc.sql, tc.name)
		}
		if sel.GroupClause == nil {
			t.Fatalf("Parse(%q): GroupClause = nil, want %s", tc.sql, tc.name)
		}
	}
}
