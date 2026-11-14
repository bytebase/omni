package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
)

// parseOne parses sql as one statement and checks that every node's Loc lies
// within the text, Start before End.
func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	stmt := list.Items[0].(*ast.RawStmt).Stmt
	ast.Inspect(stmt, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		loc := n.GetLoc()
		if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End {
			t.Errorf("Parse(%q): %T has Loc %+v", sql, n, loc)
		}
		return true
	})
	return stmt
}

// TestGrammarForms holds the parser to Parser.g on selectors, terms,
// relations, IF conditions and update operations.
func TestGrammarForms(t *testing.T) {
	accept := []string{
		// Selectors (Parser.g:304-445).
		"SELECT m['k'], m['a'..'c'], m['a'..], m[..'c'], u.f.g, l[0].f FROM t",
		"SELECT a + b * c - d / e % f, - a, -1, -NaN, -Infinity FROM t",
		"SELECT 1, 'a', 1.5, 0x1, 1h30m, 01234567-0123-0123-0123-0123456789ab, true, NULL, ?, :b FROM t",
		"SELECT (int) ?, (text) 'a', (frozen<list<int>>) [1], (a) b FROM t",
		"SELECT (a), (a, b), [], [a, 1], {}, {a}, {a: 1, 'b': c} FROM t",
		"SELECT count(*), count(a), writetime(a), maxwritetime(a.b), ttl(a[0]) FROM t",
		"SELECT CAST(a AS text), ks.f(a, 1), f() AS x FROM t",
		"SELECT DISTINCT FROM t",
		"SELECT DISTINCT a FROM t",
		"SELECT JSON FROM t",
		"SELECT JSON DISTINCT a FROM t",
		"SELECT ttl, count, json, distinct FROM t",
		// Relations (Parser.g:1774).
		"SELECT * FROM t WHERE a IN ? AND b IN :m AND c IN () AND d IN (1, ?)",
		"SELECT * FROM t WHERE (a) = ? AND (a, b) IN ? AND (a) IN () AND (a) IN (?, :m) AND (a, b) > (1, 2)",
		"SELECT * FROM t WHERE (a = 1) AND ((b > 2))",
		"SELECT * FROM t WHERE a IS NOT NULL AND b LIKE 'x%' AND m['k'] = 1 AND s CONTAINS KEY 'k'",
		"SELECT * FROM t WHERE expr(idx, 'q') AND token(a, b) > token(1, 2)",
		"SELECT * FROM t WHERE a = (int) ? AND b = [1, 2] + [3] AND c = {f: 1}",
		"SELECT * FROM t ORDER BY v ANN OF [1, 2] LIMIT 5",
		"SELECT * FROM t ORDER BY v ANN OF ? DESC",
		// A sign written against a number is part of it (Lexer.g:337).
		"INSERT INTO t (a, b, c) VALUES (CAST(-1 AS int), CAST(-1.5 AS double), (int) -1)",
		"UPDATE t SET a = CAST(-1 AS int), b = CAST(-NaN AS double) WHERE k = 1",
		// IF conditions (Parser.g:1728).
		"UPDATE t SET a = 1 WHERE k = 1 IF a.f = 1 AND m['k'] IN ? AND b IN (1) AND s CONTAINS 1 AND u.f IN :m",
		"DELETE FROM t WHERE k = 1 IF a IN ?",
		// Update operations (Parser.g:1668).
		"UPDATE t SET a += 1, b -= {'x'}, c = c + 1, d = d - 1, e = [1] + e, f = f 3 WHERE k = 1",
		"UPDATE t SET m['k'] = 1, m['j'] = 2, u.f = 1, u.g = 2, n = n + {'k': 1}, n['j'] = 3 WHERE k = 1",
		`UPDATE t SET A = a + 1, "b" = "b" - 1 WHERE k = 1`,
		"INSERT INTO t (a) VALUES (1 + 2 * 3)",
		"INSERT INTO t JSON ?",
		"INSERT INTO t JSON :j DEFAULT UNSET",
		"DELETE m[1 + 1], u.f FROM t WHERE k = 1",
	}
	for _, sql := range accept {
		parseOne(t, sql)
	}

	reject := []string{
		"SELECT ks.* FROM t",                          // no qualified star
		"INSERT INTO t (a) VALUES (CAST(- 1 AS int))", // a sign apart from its number is no constant
		"SELECT counter(*) FROM t",                    // only COUNT takes *
		"SELECT ttl() FROM t",                         // TTL takes a column
		"SELECT writetime(1) FROM t",                  // WRITETIME takes a column
		"SELECT a.keyspace FROM t",                    // a reserved word is no field name
		"SELECT : keyspace FROM t",                    // nor a bind marker name
		"SELECT CAST(a AS list<int>) FROM t",          // CAST takes a native type
		"SELECT * FROM t WHERE token() = 1",
		"SELECT * FROM t WHERE trigger(a) = 1",
		"SELECT * FROM t WHERE a IN (b)", // a name is no term
		"SELECT * FROM t WHERE a.b = 1",
		"SELECT * FROM t WHERE (a) IN ((1), ?)",
		"SELECT * FROM t WHERE expr (idx, 'q')", // 'expr(' is one token
		"SELECT * FROM t WHERE EXPR(idx, 'q')",
		"UPDATE t SET a = a WHERE k = 1",
		"UPDATE t SET a = b + 1 WHERE k = 1", // X = X + <value> only
		"UPDATE t SET a = 1 + b WHERE k = 1", // X = <value> + X only
		"UPDATE t SET a = b 3 WHERE k = 1",
		"UPDATE t SET a = 1, a = 2 WHERE k = 1",     // set twice
		"UPDATE t SET a += 1, a = 2 WHERE k = 1",    // then set
		"UPDATE t SET u.f = 1, u.f = 2 WHERE k = 1", // a field set twice
		"UPDATE t SET A = 1, a += 1 WHERE k = 1",    // A and a are one column
		"INSERT INTO t (a) VALUES (b)",
		"INSERT INTO t VALUES (1)", // the column list is not optional
		"INSERT INTO t (a) JSON '{}'",
		"INSERT INTO t JSON 1",
		"UPDATE t SET a = 1 WHERE k = 1 IF a LIKE 'x'",
	}
	for _, sql := range reject {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want an error", sql)
		}
	}
}

func TestSelectorAST(t *testing.T) {
	sel := parseOne(t, "SELECT m['a'..], -1, - a, (int) ?, writetime(u.f) FROM t").(*ast.SelectStmt)
	slice, ok := sel.Elements[0].Expr.(*ast.SliceAccess)
	if !ok || slice.From == nil || slice.To != nil {
		t.Errorf("m['a'..] = %#v, want a SliceAccess from 'a'", sel.Elements[0].Expr)
	}
	if lit, ok := sel.Elements[1].Expr.(*ast.IntegerLit); !ok || lit.Val != "-1" {
		t.Errorf("-1 = %#v, want IntegerLit -1", sel.Elements[1].Expr)
	}
	if u, ok := sel.Elements[2].Expr.(*ast.UnaryExpr); !ok || u.Op != "-" {
		t.Errorf("- a = %#v, want a UnaryExpr", sel.Elements[2].Expr)
	}
	if h, ok := sel.Elements[3].Expr.(*ast.TypeHintExpr); !ok || h.Type.Name.Name != "int" {
		t.Errorf("(int) ? = %#v, want a TypeHintExpr", sel.Elements[3].Expr)
	}
	if f, ok := sel.Elements[4].Expr.(*ast.FunctionCall); !ok || len(f.Args) != 1 {
		t.Errorf("writetime(u.f) = %#v, want a call of one argument", sel.Elements[4].Expr)
	} else if _, ok := f.Args[0].(*ast.DotAccess); !ok {
		t.Errorf("writetime(u.f) argument = %#v, want a DotAccess", f.Args[0])
	}

	// A sign written against a number is part of it (Lexer.g:337), so
	// (int) -1 is a type hint on -1, not (int) minus 1.
	for i, want := range []string{"-1", "-1.5", "-Infinity"} {
		hinted := parseOne(t, "SELECT (int) -1, (double) -1.5, (double) -Infinity FROM t").(*ast.SelectStmt)
		h, ok := hinted.Elements[i].Expr.(*ast.TypeHintExpr)
		if !ok {
			t.Errorf("element %d = %#v, want a TypeHintExpr on %s", i, hinted.Elements[i].Expr, want)
			continue
		}
		if got := h.Expr.GetLoc(); got.End-got.Start != len(want) {
			t.Errorf("element %d hints %#v, want %s", i, h.Expr, want)
		}
	}

	d := parseOne(t, "SELECT DISTINCT FROM t").(*ast.SelectStmt)
	if d.Distinct || len(d.Elements) != 1 || d.Elements[0].Expr.(*ast.Identifier).Name != "DISTINCT" {
		t.Errorf("SELECT DISTINCT FROM t: %+v, want the column distinct", d)
	}
	j := parseOne(t, "SELECT JSON DISTINCT a FROM t").(*ast.SelectStmt)
	if !j.JSON || !j.Distinct {
		t.Errorf("SELECT JSON DISTINCT a: JSON %v, Distinct %v", j.JSON, j.Distinct)
	}

	ann := parseOne(t, "SELECT * FROM t ORDER BY v ANN OF [1, 2] LIMIT 5").(*ast.SelectStmt)
	if _, ok := ann.OrderBy[0].AnnVector.(*ast.VectorLit); !ok || ann.Limit == nil {
		t.Errorf("ANN OF [1, 2] LIMIT 5: vector %#v, limit %#v", ann.OrderBy[0].AnnVector, ann.Limit)
	}
}

func TestRelationAST(t *testing.T) {
	sel := parseOne(t, "SELECT * FROM t WHERE a IN ? AND (b, c) IN (?, ?) AND (b) = :m AND (d IS NOT NULL) AND expr(i, 'q')").(*ast.SelectStmt)
	if in, ok := sel.Where[0].(*ast.InExpr); !ok || in.Marker == nil {
		t.Errorf("a IN ? = %#v", sel.Where[0])
	}
	if in, ok := sel.Where[1].(*ast.TupleInExpr); !ok || len(in.Markers) != 2 {
		t.Errorf("(b, c) IN (?, ?) = %#v", sel.Where[1])
	}
	if cmp, ok := sel.Where[2].(*ast.TupleCompareExpr); !ok || cmp.Marker == nil || cmp.Marker.Name != "m" {
		t.Errorf("(b) = :m = %#v", sel.Where[2])
	}
	if _, ok := sel.Where[3].(*ast.IsNotNullExpr); !ok {
		t.Errorf("(d IS NOT NULL) = %#v", sel.Where[3])
	}
	if c, ok := sel.Where[4].(*ast.CustomIndexExpr); !ok || c.Index.Name != "i" {
		t.Errorf("expr(i, 'q') = %#v", sel.Where[4])
	}
}

func TestAssignmentAST(t *testing.T) {
	up := parseOne(t, "UPDATE t SET a += 1, b = b - 2, c = [1] + c, d = d 3, m['k'] = 1 WHERE k = 1 IF m['k'] IN ? AND u.f = 1").(*ast.UpdateStmt)
	for i, want := range []string{"+", "-", "+", "+", "="} {
		if got := up.Assignments[i].Operator; got != want {
			t.Errorf("assignment %d: Operator %q, want %q", i, got, want)
		}
	}
	if c := up.IfConditions[0]; c.Element == nil || c.InMarker == nil || c.Op != "IN" {
		t.Errorf("IF m['k'] IN ? = %+v", c)
	}
	if c := up.IfConditions[1]; c.Field == nil || c.Field.Name != "f" {
		t.Errorf("IF u.f = 1 = %+v", c)
	}

	_, err := Parse("UPDATE t SET a = 1, a = 2 WHERE k = 1")
	if err == nil || !strings.Contains(err.Error(), "multiple incompatible setting of column a") {
		t.Errorf("a set twice: err = %v", err)
	}
	if pe, ok := err.(*ParseError); !ok || pe.Position != strings.Index("UPDATE t SET a = 1, a = 2", "a = 2") {
		t.Errorf("a set twice: err = %#v, want it at the second assignment", err)
	}
}

func TestDurationLexing(t *testing.T) {
	for _, d := range []string{"1d", "1h30m", "2mo", "3ms", "4us", "5µs", "6ns", "1y2mo3w4d5h6m7s8ms9us10ns"} {
		ins := parseOne(t, "INSERT INTO t (a) VALUES ("+d+")").(*ast.InsertStmt)
		if lit, ok := ins.Values[0].(*ast.DurationLit); !ok || lit.Val != d {
			t.Errorf("%s: %#v, want a DurationLit", d, ins.Values[0])
		}
	}
	ins := parseOne(t, "INSERT INTO t (a) VALUES (-1d)").(*ast.InsertStmt)
	if lit, ok := ins.Values[0].(*ast.DurationLit); !ok || lit.Val != "-1d" {
		t.Errorf("-1d: %#v, want DurationLit -1d", ins.Values[0])
	}
}
