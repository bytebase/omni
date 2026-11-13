package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// JSON_QUERY's own arm of the grammar, as OceanBase 4.3.5 reads it, every verdict measured
// with PREPARE on oceanbase-ce:4.3.5.6 on 2026-11-13. A statement that parses may still be
// refused: the server answers 5443 for most of these, since JSON_QUERY must return JSON or
// be written ASIS, and that is no verdict on the syntax.
//
// Ref: json_query_expr in OceanBase's sql_parser_mysql_mode.y (v4.3.5_CE_BP6, line 24917).
var jsonQueryCases = []parserFormCase{
	// The document and the path alone, and each clause on its own.
	{"SELECT JSON_QUERY(a, '$')", true},
	{"SELECT JSON_QUERY(a, '$' RETURNING JSON)", true},
	{"SELECT JSON_QUERY(a, '$' RETURNING CHAR(10) CHARSET ascii)", true},
	{"SELECT JSON_QUERY(a, '$' RETURNING UNSIGNED)", true},
	{"SELECT JSON_QUERY(a, '$' RETURNING NUMBER(10,2))", true},
	{"SELECT JSON_QUERY(a, '$' TRUNCATE)", true},
	{"SELECT JSON_QUERY(a, '$' ALLOW SCALARS)", true},
	{"SELECT JSON_QUERY(a, '$' DISALLOW SCALARS)", true},
	{"SELECT JSON_QUERY(a, '$' PRETTY)", true},
	{"SELECT JSON_QUERY(a, '$' ASCII)", true},
	{"SELECT JSON_QUERY(a, '$' WITHOUT WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITHOUT ARRAY WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH ARRAY WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH UNCONDITIONAL WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH CONDITIONAL WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH UNCONDITIONAL ARRAY WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' WITH CONDITIONAL ARRAY WRAPPER)", true},
	{"SELECT JSON_QUERY(a, '$' ASIS)", true},
	{"SELECT JSON_QUERY(a, '$' ERROR ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' NULL ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' ERROR ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' NULL ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' ERROR ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' NULL ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' DOT ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' MULTIVALUE)", true},
	// ON ERROR and ON EMPTY in either order, and ON MISMATCH after them.
	{"SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON EMPTY)", true},
	{"SELECT JSON_QUERY(a, '$' NULL ON EMPTY ERROR ON ERROR)", true},
	{"SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' NULL ON EMPTY NULL ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH)", true},
	{"SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON EMPTY EMPTY OBJECT ON ERROR DOT ON MISMATCH)", true},
	// Every clause at once, and the text a multi-valued index reads back as.
	{"SELECT JSON_QUERY(a, '$' RETURNING CHAR(10) TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER ASIS NULL ON EMPTY NULL ON ERROR NULL ON MISMATCH MULTIVALUE)", true},
	{"SELECT json_query(`a`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)", true},
	{"SELECT json_query(`a`,'$.codes' RETURNING char(40) charset utf8mb4 WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)", true},
	{"SELECT json_query(`a`,'$.amts' RETURNING number(10,2) WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)", true},
	{"CREATE TABLE t2 (id INT PRIMARY KEY, j JSON, INDEX mv ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))))", true},
	// The document is a simple_expr and the path a string literal.
	{"SELECT JSON_QUERY((a), '$')", true},
	{"SELECT JSON_QUERY(-a, '$')", true},
	{"SELECT JSON_QUERY(a->'$.x', '$')", true},
	{"SELECT JSON_QUERY(a COLLATE utf8mb4_bin, '$')", true},
	{"SELECT JSON_QUERY(a, N'$')", true},
	{`SELECT JSON_QUERY(a, "$")`, true},
	{"SELECT JSON_QUERY(a, '$' '.x')", true},
	{"SELECT JSON_QUERY(a, _utf8mb4'$')", true},
	{"SELECT JSON_QUERY(a, _latin1 X'24')", true},
	// Quoted or qualified, the name is an ordinary call, and unquoted it is still a name.
	{"SELECT `json_query`(a, '$')", true},
	{"SELECT `json_query`(1, 2)", true},
	{"SELECT d.json_query(1, 2)", true},
	{"SELECT json_query FROM (SELECT 1 AS json_query) x", true},

	// The arguments.
	{"SELECT JSON_QUERY()", false},
	{"SELECT JSON_QUERY(*)", false},
	{"SELECT JSON_QUERY(a)", false},
	{"SELECT JSON_QUERY(a, '$', '$')", false},
	{"SELECT JSON_QUERY(1, 2)", false},
	{"SELECT JSON_QUERY(a, a)", false},
	{"SELECT JSON_QUERY(a, @p)", false},
	{"SELECT JSON_QUERY(a, CONCAT('$'))", false},
	{"SELECT JSON_QUERY(a + 1, '$')", false},
	{"SELECT JSON_QUERY(a IS NULL, '$')", false},
	// RETURNING takes a cast's type, once, and comes first.
	{"SELECT JSON_QUERY(a, '$' RETURNING)", false},
	{"SELECT JSON_QUERY(a, '$' RETURNING VARCHAR(10))", false},
	{"SELECT JSON_QUERY(a, '$' RETURNING TEXT)", false},
	{"SELECT JSON_QUERY(a, '$' RETURNING JSON RETURNING JSON)", false},
	{"SELECT JSON_QUERY(a, '$' TRUNCATE RETURNING JSON)", false},
	// The clauses in another order, or twice.
	{"SELECT JSON_QUERY(a, '$' PRETTY ALLOW SCALARS)", false},
	{"SELECT JSON_QUERY(a, '$' ASCII PRETTY)", false},
	{"SELECT JSON_QUERY(a, '$' WITH WRAPPER ASCII)", false},
	{"SELECT JSON_QUERY(a, '$' ASIS WITH WRAPPER)", false},
	{"SELECT JSON_QUERY(a, '$' MULTIVALUE ASIS)", false},
	{"SELECT JSON_QUERY(a, '$' MULTIVALUE MULTIVALUE)", false},
	{"SELECT JSON_QUERY(a, '$' ALLOW SCALARS DISALLOW SCALARS)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON ERROR NULL ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON EMPTY NULL ON EMPTY)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON MISMATCH NULL ON MISMATCH)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON MISMATCH NULL ON EMPTY)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON MISMATCH ERROR ON ERROR)", false},
	// A clause cut short, or with a word it does not take.
	{"SELECT JSON_QUERY(a, '$' ALLOW)", false},
	{"SELECT JSON_QUERY(a, '$' SCALARS)", false},
	{"SELECT JSON_QUERY(a, '$' WITHOUT)", false},
	{"SELECT JSON_QUERY(a, '$' WITH)", false},
	{"SELECT JSON_QUERY(a, '$' WITHOUT ARRAY)", false},
	{"SELECT JSON_QUERY(a, '$' WITH CONDITIONAL)", false},
	{"SELECT JSON_QUERY(a, '$' WITHOUT CONDITIONAL WRAPPER)", false},
	{"SELECT JSON_QUERY(a, '$' WITH ARRAY CONDITIONAL WRAPPER)", false},
	{"SELECT JSON_QUERY(a, '$' ERROR ON)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' EMPTY ON MISMATCH)", false},
	{"SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON MISMATCH)", false},
	{"SELECT JSON_QUERY(a, '$' DOT ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' DOT ON EMPTY)", false},
	{"SELECT JSON_QUERY(a, '$' DEFAULT 'x' ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' DEFAULT 'x' ON EMPTY)", false},
	{"SELECT JSON_QUERY(a, '$' RETURNING JSON,)", false},
	// A clause's words are keywords, so a quoted one is no clause.
	{"SELECT JSON_QUERY(a, '$' `asis`)", false},
	{"SELECT JSON_QUERY(a, '$' `pretty`)", false},
	{"SELECT JSON_QUERY(a, '$' `multivalue`)", false},
	{"SELECT JSON_QUERY(a, '$' ALLOW `scalars`)", false},
	{"SELECT JSON_QUERY(a, '$' ERROR ON `error`)", false},
	{"SELECT JSON_QUERY(a, '$' NULL ON `mismatch`)", false},
	{"SELECT JSON_QUERY(a, '$' `dot` ON MISMATCH)", false},
	{"SELECT JSON_QUERY(a, '$' `empty` ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' EMPTY `object` ON ERROR)", false},
	{"SELECT JSON_QUERY(a, '$' WITH `conditional` WRAPPER)", false},
	// The clauses belong to the arm: an ordinary call has none, and the arm takes no OVER.
	{"SELECT `json_query`(a, '$' RETURNING JSON)", false},
	{"SELECT d.json_query(a, '$' RETURNING JSON)", false},
	{"SELECT JSON_QUERY(a, '$' RETURNING JSON) OVER ()", false},
}

// TestJSONQueryForms holds omni to the recorded verdicts without a server.
func TestJSONQueryForms(t *testing.T) {
	for _, c := range jsonQueryCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestJSONQueryTree holds what the clauses leave in the tree.
func TestJSONQueryTree(t *testing.T) {
	call := func(t *testing.T, sql string) *nodes.FuncCallExpr {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		target := list.Items[0].(*nodes.SelectStmt).TargetList[0]
		if rt, ok := target.(*nodes.ResTarget); ok {
			target = rt.Val
		}
		fc, ok := target.(*nodes.FuncCallExpr)
		if !ok {
			t.Fatalf("%q reads as %T, want a call", sql, target)
		}
		return fc
	}

	t.Run("every clause", func(t *testing.T) {
		sql := "SELECT JSON_QUERY(a, '$.x' returning char(10) truncate allow scalars pretty ascii with conditional array wrapper asis empty array on empty empty object on error dot on mismatch multivalue)"
		fc := call(t, sql)
		if !fc.JSONQueryArm || fc.JSONQuery == nil || len(fc.Args) != 2 {
			t.Fatalf("arm %v, options %v, %d arguments", fc.JSONQueryArm, fc.JSONQuery, len(fc.Args))
		}
		got := *fc.JSONQuery
		if got.Returning == nil || got.Returning.Name != "CHAR" || got.Returning.Length != 10 {
			t.Errorf("RETURNING is %+v, want CHAR(10)", got.Returning)
		}
		want := nodes.JSONQueryOptions{
			Loc: got.Loc, Returning: got.Returning,
			Truncate: true, Scalars: "ALLOW", Pretty: true, ASCII: true,
			Wrapper: "WITH CONDITIONAL ARRAY WRAPPER", Asis: true,
			OnError: "EMPTY OBJECT", OnEmpty: "EMPTY ARRAY", OnMismatch: "DOT", Multivalue: true,
		}
		if got != want {
			t.Errorf("options are %+v, want %+v", got, want)
		}
		if text := sql[got.Loc.Start:got.Loc.End]; text != sql[len("SELECT JSON_QUERY(a, '$.x' "):len(sql)-1] {
			t.Errorf("options span %q", text)
		}
	})

	t.Run("no clause", func(t *testing.T) {
		fc := call(t, "SELECT JSON_QUERY(a, '$.x')")
		if !fc.JSONQueryArm || fc.JSONQuery != nil {
			t.Errorf("arm %v, options %v, want the arm with no options", fc.JSONQueryArm, fc.JSONQuery)
		}
	})

	t.Run("a quoted name is an ordinary call", func(t *testing.T) {
		fc := call(t, "SELECT `json_query`(a, '$.x')")
		if fc.JSONQueryArm || fc.JSONQuery != nil {
			t.Errorf("arm %v, options %v, want neither", fc.JSONQueryArm, fc.JSONQuery)
		}
	})
}
