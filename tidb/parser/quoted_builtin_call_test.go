package parser

import (
	"context"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// quotedBuiltinCallCases are calls whose name is a backtick-quoted or
// schema-qualified builtin with special argument syntax: COUNT(*), DISTINCT,
// TRIM's and SUBSTRING's FROM, GROUP_CONCAT's ORDER BY and SEPARATOR, a bare
// INTERVAL, OVER. The grammar reads that syntax from the name's keyword, and a
// quoted name is an identifier and a qualified one a stored function, so the
// call is an ordinary one: TiDB v8.5.5 parses each accepted case (runs it, or
// answers 1305 for the missing stored function, or 1582 or 1105 when it
// resolves the name to its builtin after all) and answers 1064 for every
// rejected one. tq has a INT, b VARCHAR(20) and d DATE.
var quotedBuiltinCallCases = []struct {
	sql    string
	accept bool
}{
	// A quoted name takes ordinary arguments, any number of them.
	{"SELECT `count`(a) FROM tq", true},
	{"SELECT `count`()", true},
	{"SELECT `sum`(a) FROM tq", true},
	{"SELECT `trim`(b) FROM tq", true},
	{"SELECT `trim`(b, 'x') FROM tq", true},
	{"SELECT `trim`() FROM tq", true},
	{"SELECT `substring`(b, 1, 2) FROM tq", true},
	{"SELECT `substring`(b) FROM tq", true},
	{"SELECT `substring`(b, 1, 2, 3) FROM tq", true},
	{"SELECT `substr`(b, 1) FROM tq", true},
	{"SELECT `group_concat`(b) FROM tq", true},
	{"SELECT `group_concat`(a, b) FROM tq", true},
	{"SELECT `rank`() FROM tq", true},
	{"SELECT `row_number`() FROM tq", true},
	{"SELECT `cast`(a) FROM tq", true},
	{"SELECT `extract`(d) FROM tq", true},
	{"SELECT `position`(b) FROM tq", true},
	{"SELECT `char`(65)", true},
	{"SELECT `date_add`(d, 1) FROM tq", true},
	{"SELECT `f`(d + INTERVAL 1 DAY) FROM tq", true},
	{"SELECT test.count(a) FROM tq", true},
	{"SELECT test.trim(b) FROM tq", true},
	{"SELECT test.trim(b, 'x') FROM tq", true},
	{"SELECT test.substring(b, 1, 2) FROM tq", true},
	{"SELECT test.substr(b, 1) FROM tq", true},
	// ... and none of the builtin's own syntax.
	{"SELECT `count`(*) FROM tq", false},
	{"SELECT `COUNT`(*) FROM tq", false},
	{"SELECT `count`(*, a) FROM tq", false},
	{"SELECT `count`(DISTINCT a) FROM tq", false},
	{"SELECT `count`(DISTINCT a, b) FROM tq", false},
	{"SELECT `sum`(DISTINCT a) FROM tq", false},
	{"SELECT `avg`(DISTINCT a) FROM tq", false},
	{"SELECT `min`(DISTINCT a) FROM tq", false},
	{"SELECT `max`(DISTINCT a) FROM tq", false},
	{"SELECT `json_arrayagg`(DISTINCT a) FROM tq", false},
	{"SELECT `group_concat`(DISTINCT b) FROM tq", false},
	{"SELECT `group_concat`(b SEPARATOR ',') FROM tq", false},
	{"SELECT `group_concat`(b ORDER BY b) FROM tq", false},
	{"SELECT `group_concat`(b ORDER BY b SEPARATOR ',') FROM tq", false},
	{"SELECT `trim`(LEADING 'x' FROM b) FROM tq", false},
	{"SELECT `trim`(BOTH FROM b) FROM tq", false},
	{"SELECT `trim`('x' FROM b) FROM tq", false},
	{"SELECT `substring`(b FROM 1) FROM tq", false},
	{"SELECT `substring`(b FROM 1 FOR 2) FROM tq", false},
	{"SELECT `substr`(b FROM 1) FROM tq", false},
	{"SELECT `date_add`(d, INTERVAL 1 DAY) FROM tq", false},
	{"SELECT `adddate`(d, INTERVAL 1 DAY) FROM tq", false},
	{"SELECT `sum`(a) OVER () FROM tq", false},
	{"SELECT `rank`() OVER () FROM tq", false},
	{"SELECT `row_number`() OVER () FROM tq", false},
	{"SELECT `lag`(a) OVER () FROM tq", false},
	{"SELECT `json_arrayagg`(a) OVER () FROM tq", false},
	{"SELECT `json_objectagg`(a, b) OVER () FROM tq", false},
	{"SELECT `count`(a) OVER w FROM tq WINDOW w AS ()", false},
	{"SELECT `count`(ALL a) FROM tq", false},
	{"SELECT `cast`(a AS CHAR) FROM tq", false},
	{"SELECT `convert`(b USING utf8mb4) FROM tq", false},
	{"SELECT `convert`(b, CHAR) FROM tq", false},
	{"SELECT `extract`(YEAR FROM d) FROM tq", false},
	{"SELECT `position`('a' IN b) FROM tq", false},
	{"SELECT `char`(65 USING utf8mb4)", false},
	{"SELECT * FROM tq WHERE `match`(b) AGAINST ('x')", false},
	{"SELECT * FROM `json_table`('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt", false},
	{"SELECT test.count(*) FROM tq", false},
	{"SELECT test.count(DISTINCT a) FROM tq", false},
	{"SELECT test.sum(DISTINCT a) FROM tq", false},
	{"SELECT test.trim(LEADING 'x' FROM b) FROM tq", false},
	{"SELECT test.trim(TRAILING 'x' FROM b) FROM tq", false},
	{"SELECT test.trim('x' FROM b) FROM tq", false},
	{"SELECT `test`.trim(LEADING 'x' FROM b) FROM tq", false},
	{"SELECT test.`trim`(LEADING 'x' FROM b) FROM tq", false},
	{"SELECT test.substring(b FROM 1) FROM tq", false},
	{"SELECT test.substring(b FROM 1 FOR 2) FROM tq", false},
	{"SELECT test.substr(b FROM 1) FROM tq", false},
	{"SELECT test.group_concat(b SEPARATOR ',') FROM tq", false},
	{"SELECT test.group_concat(b ORDER BY b) FROM tq", false},
	{"SELECT test.group_concat(DISTINCT b) FROM tq", false},
	{"SELECT test.date_add(d, INTERVAL 1 DAY) FROM tq", false},
	{"SELECT test.sum(a) OVER () FROM tq", false},
	{"SELECT test.count(a) OVER () FROM tq", false},
	// The unquoted builtin keeps its own grammar.
	{"SELECT COUNT(*) FROM tq", true},
	{"SELECT COUNT(DISTINCT a) FROM tq", true},
	{"SELECT TRIM(LEADING 'x' FROM b) FROM tq", true},
	{"SELECT SUBSTRING(b FROM 1 FOR 2) FROM tq", true},
	{"SELECT GROUP_CONCAT(b ORDER BY b SEPARATOR ',') FROM tq", true},
	{"SELECT SUM(a) OVER () FROM tq", true},
	{"SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM tq", true},
	{"SELECT JSON_ARRAYAGG(a) OVER () FROM tq", true},
}

// TestQuotedBuiltinCall holds the parser to quotedBuiltinCallCases.
func TestQuotedBuiltinCall(t *testing.T) {
	for _, tc := range quotedBuiltinCallCases {
		if _, err := Parse(tc.sql); (err == nil) != tc.accept {
			t.Errorf("Parse(%q) error = %v, want accepted %v", tc.sql, err, tc.accept)
		}
	}
}

// TestQuotedSpecialFormCallAST pins the tree of a quoted or qualified call
// whose name is a builtin with special syntax: an ordinary call that keeps
// every argument, with no star, DISTINCT or window.
func TestQuotedSpecialFormCallAST(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		schema string
		name   string
		args   int
	}{
		{"SELECT `count`(a) FROM tq", "", "COUNT", 1},
		{"SELECT `trim`(b, 'x') FROM tq", "", "TRIM", 2},
		{"SELECT `group_concat`(a, b) FROM tq", "", "GROUP_CONCAT", 2},
		{"SELECT test.substring(b, 1, 2, 3) FROM tq", "test", "SUBSTRING", 4},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		var call *ast.FuncCallExpr
		ast.Inspect(list.Items[0], func(n ast.Node) bool {
			if fc, ok := n.(*ast.FuncCallExpr); ok && call == nil {
				call = fc
			}
			return true
		})
		if call == nil || call.Schema != tc.schema || call.Name != tc.name || len(call.Args) != tc.args {
			t.Fatalf("%q: call = %+v, want %s.%s with %d arguments", tc.sql, call, tc.schema, tc.name, tc.args)
		}
		if call.Star || call.Quantifier != ast.AggQuantNone || call.Over != nil {
			t.Errorf("%q: call = %+v, want no star, DISTINCT or window", tc.sql, call)
		}
		if col, ok := call.Args[0].(*ast.ColumnRef); !ok || col.Column == "" {
			t.Errorf("%q: first argument = %#v, want a column", tc.sql, call.Args[0])
		}
	}
}

// TestQuotedBuiltinCallTiDBOracle asserts omni agrees with a live TiDB v8.5.5
// on quotedBuiltinCallCases (container test). The statements run in the
// DSN's database, test, on one connection.
func TestQuotedBuiltinCallTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, ddl := range []string{
		"USE test",
		"DROP TABLE IF EXISTS tq",
		"CREATE TABLE tq (a INT, b VARCHAR(20), d DATE)",
	} {
		if _, err := conn.ExecContext(tc.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	defer conn.ExecContext(context.Background(), "DROP TABLE IF EXISTS tq")
	for _, c := range quotedBuiltinCallCases {
		_, execErr := conn.ExecContext(tc.ctx, c.sql)
		tidbAccepts := !tidbRejectedSyntax(t, execErr)
		if tidbAccepts != c.accept {
			t.Errorf("TiDB accepts %q = %v (%v), want %v", c.sql, tidbAccepts, execErr, c.accept)
		}
		if _, err := Parse(c.sql); (err == nil) != tidbAccepts {
			t.Errorf("omni accepts %q = %v (err %v), TiDB = %v", c.sql, err == nil, err, tidbAccepts)
		}
	}
}
