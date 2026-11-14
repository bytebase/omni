package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/cosmosdb/ast"
)

func TestParseSelectStarFromC(t *testing.T) {
	result, err := Parse("SELECT * FROM c")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if result.Len() != 1 {
		t.Fatalf("expected 1 item, got %d", result.Len())
	}
	raw, ok := result.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected *RawStmt, got %T", result.Items[0])
	}
	sel, ok := raw.Stmt.(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("expected *SelectStmt, got %T", raw.Stmt)
	}
	if !sel.Star {
		t.Error("expected Star=true")
	}
	if sel.From == nil {
		t.Fatal("expected From to be set")
	}
	cr, ok := sel.From.(*nodes.ContainerRef)
	if !ok {
		t.Fatalf("expected *ContainerRef, got %T", sel.From)
	}
	if cr.Name != "c" {
		t.Errorf("expected container name 'c', got %q", cr.Name)
	}
}

func TestParseSelectWithWhere(t *testing.T) {
	result, err := Parse("SELECT c.id, c.name FROM c WHERE c.age > 18")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	raw := result.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if len(sel.Targets) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(sel.Targets))
	}
	if sel.Where == nil {
		t.Fatal("expected WHERE clause")
	}
}

func TestParseSelectTopDistinctValue(t *testing.T) {
	result, err := Parse("SELECT TOP 10 DISTINCT VALUE c.name FROM c")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	raw := result.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if sel.Top == nil || *sel.Top != 10 {
		t.Errorf("expected TOP 10")
	}
	if !sel.Distinct {
		t.Error("expected DISTINCT")
	}
	if !sel.Value {
		t.Error("expected VALUE")
	}
}

func TestParseJoin(t *testing.T) {
	result, err := Parse("SELECT * FROM c JOIN d IN c.children")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	raw := result.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if len(sel.Joins) != 1 {
		t.Fatalf("expected 1 join, got %d", len(sel.Joins))
	}
	iter, ok := sel.Joins[0].Source.(*nodes.ArrayIterationExpr)
	if !ok {
		t.Fatalf("expected ArrayIterationExpr, got %T", sel.Joins[0].Source)
	}
	if iter.Alias != "d" {
		t.Errorf("expected alias 'd', got %q", iter.Alias)
	}
}

func TestParseOrderByOffsetLimit(t *testing.T) {
	result, err := Parse("SELECT * FROM c ORDER BY c.name ASC OFFSET 0 LIMIT 10")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	raw := result.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if len(sel.OrderBy) != 1 {
		t.Fatalf("expected 1 order by, got %d", len(sel.OrderBy))
	}
	if sel.OffsetLimit == nil {
		t.Fatal("expected offset limit clause")
	}
}

func TestParseExpressions(t *testing.T) {
	tests := []string{
		"SELECT c.a + c.b FROM c",
		"SELECT c.a BETWEEN 1 AND 10 FROM c",
		"SELECT c.a NOT IN (1, 2, 3) FROM c",
		"SELECT c.a LIKE 'foo%' ESCAPE '\\\\' FROM c",
		"SELECT c.a ? 'yes' : 'no' FROM c",
		"SELECT c.a ?? 'default' FROM c",
		"SELECT EXISTS(SELECT * FROM c) FROM c",
		"SELECT ARRAY(SELECT VALUE t FROM t IN c.tags) FROM c",
		"SELECT udf.myFunc(c.x) FROM c",
		"SELECT [1, 2, 3] FROM c",
		"SELECT {'a': 1, 'b': 2} FROM c",
		"SELECT (c.a + c.b) * c.d FROM c",
	}
	for _, sql := range tests {
		_, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q) failed: %v", sql, err)
		}
	}
}

func TestParseGroupByHaving(t *testing.T) {
	result, err := Parse("SELECT c.city, COUNT(1) FROM c GROUP BY c.city HAVING COUNT(1) > 5")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	raw := result.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	if len(sel.GroupBy) != 1 {
		t.Fatalf("expected 1 group by expr, got %d", len(sel.GroupBy))
	}
	if sel.Having == nil {
		t.Fatal("expected HAVING clause")
	}
}

// TestParseGrammarForms holds the forms sql.g4 derives that omni once read
// differently: each case parses, or fails, as the grammar says.
func TestParseGrammarForms(t *testing.T) {
	for _, tc := range []struct {
		sql string
		ok  bool
	}{
		{"SELECT TOP @n * FROM c", true},                    // sql.g4:13
		{"SELECT * FROM c ORDER BY RANK f(c.x) DESC", true}, // sql.g4:67
		{"SELECT ALL(SELECT VALUE 1)", true},                // sql.g4:183
		{"SELECT first(SELECT VALUE 1)", true},              // sql.g4:184
		{"SELECT ALL(c) FROM c", true},                      // ALL is also a name, sql.g4:203
		{"SELECT LEFT(c.s, 1) FROM c", true},                // sql.g4:191
		{"SELECT c.left FROM c", true},                      // a property name after a dot
		{"SELECT COUNT(*) FROM c", false},                   // an argument is an expression, sql.g4:190
		{"SELECT c IN () FROM c", false},                    // sql.g4:105
		{"SELECT LEFT FROM c", false},                       // a keyword is no name, sql.g4:201
		{"SELECT * FROM IN c", false},
		{"SELECT value FROM c", false},
		{"SELECT * FROM c JOIN JOIN IN c.a", false},
	} {
		_, err := Parse(tc.sql)
		if (err == nil) != tc.ok {
			t.Errorf("Parse(%q): err = %v, want ok = %v", tc.sql, err, tc.ok)
		}
	}
}

func TestParseSubqueryFuncAndTopParam(t *testing.T) {
	result, err := Parse("SELECT TOP @n LAST(SELECT VALUE c.a FROM c) AS l")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	sel := result.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt)
	if sel.Top != nil || sel.TopParam == nil || sel.TopParam.Name != "n" {
		t.Fatalf("TOP: got Top=%v TopParam=%+v", sel.Top, sel.TopParam)
	}
	fn, ok := sel.Targets[0].Expr.(*nodes.SubqueryFuncExpr)
	if !ok || fn.Func != "LAST" || fn.Select == nil {
		t.Fatalf("target: got %#v", sel.Targets[0].Expr)
	}
	if fn.Loc.Start != 14 || fn.Loc.End != 43 {
		t.Errorf("loc: got %+v", fn.Loc)
	}
}
