package parser

import (
	"context"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// quotedBuiltinCallCases are calls whose name is a backtick-quoted or
// schema-qualified builtin with special argument syntax. The grammar reads
// TIMESTAMPDIFF, TIMESTAMPADD, GET_FORMAT, WEIGHT_STRING, COUNT, TRIM, SUBSTRING,
// GROUP_CONCAT and the aggregates and window functions that take DISTINCT or
// OVER from their keywords, so a quoted name is an identifier, a qualified one a
// stored function, the call is an ordinary one, and MySQL 8.0 parses each
// accepted case and answers 1630 (or 1305) for the stored function it then
// cannot find. A rejected case is one an ordinary argument list cannot hold:
// MySQL answers 1064. 5.7.25-OceanBase_CE-v4.3.5.6 gives the same verdict on
// every case it is asked about (TestOceanBaseQuotedBuiltinCall). tu has DATETIME
// columns a and b; tq has a INT, b VARCHAR(20) and d DATE.
var quotedBuiltinCallCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT `timestampdiff`(1, 2, 3)", true},
	{"SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu", true},
	{"SELECT `timestampadd`(a, 2, 3) FROM tu", true},
	{"SELECT `timestampadd`(DAY, 1, a) FROM tu", true},
	{"SELECT `get_format`(1, 2)", true},
	{"SELECT `get_format`(DATE, 'USA')", true},
	{"SELECT `timestampdiff`()", true},
	{"SELECT `weight_string`('a')", true},
	// DAY_HOUR is reserved, so it cannot be a column, and an ordinary argument
	// list has no AS clause.
	{"SELECT `timestampdiff`(DAY_HOUR, a, b) FROM tu", false},
	{"SELECT `weight_string`('a' AS CHAR(2))", false},
	// The unquoted builtin keeps its own grammar.
	{"SELECT TIMESTAMPDIFF(SECOND, a, b) FROM tu", true},
	{"SELECT timestampdiff(1, 2, 3)", false},
	{"SELECT GET_FORMAT(1, 2)", false},
	{"SELECT WEIGHT_STRING('a' AS CHAR(2))", true},

	// A quoted name takes ordinary arguments, so a stored function by that name
	// can be called with any number of them.
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
	// ... and none of the builtin's own syntax: no *, DISTINCT, TRIM's or
	// SUBSTRING's FROM, GROUP_CONCAT's ORDER BY or SEPARATOR, bare INTERVAL
	// or OVER.
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

// TestQuotedBuiltinCallAST pins the tree a quoted name gives: an ordinary call
// whose first argument is the column it names, which is what MySQL resolves.
func TestQuotedBuiltinCallAST(t *testing.T) {
	list, err := Parse("SELECT `timestampdiff`(SECOND, a, b) FROM tu")
	if err != nil {
		t.Fatal(err)
	}
	var call *ast.FuncCallExpr
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		if fc, ok := n.(*ast.FuncCallExpr); ok && call == nil {
			call = fc
		}
		_, isKeyword := n.(*ast.KeywordArg)
		if isKeyword {
			t.Errorf("a quoted name built a KeywordArg")
		}
		return true
	})
	if call == nil || call.Name != "TIMESTAMPDIFF" || len(call.Args) != 3 {
		t.Fatalf("call = %+v, want TIMESTAMPDIFF with three arguments", call)
	}
	if col, ok := call.Args[0].(*ast.ColumnRef); !ok || col.Column != "SECOND" {
		t.Errorf("first argument = %#v, want the column SECOND", call.Args[0])
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

// oceanBaseCrashCases are the quotedBuiltinCallCases kept off the OceanBase
// oracle. OceanBase resolves a quoted name to its own function where it has
// one, answering 1582 for `trim`() or `extract`(d), but SELECT `cast`(a)
// FROM t took the 4.3.5.6 observer down with SIGSEGV rather than answer,
// which would leave every later case without a verdict. The parser test still
// holds it.
var oceanBaseCrashCases = map[string]bool{
	"SELECT `cast`(a) FROM tq": true,
}

// TestOceanBaseQuotedBuiltinCall asks OceanBase about quotedBuiltinCallCases. It
// runs on the package's own OceanBase, like the rest of the OceanBase oracle:
//
//	go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseQuotedBuiltinCall(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obprobe_qb",
		"USE obprobe_qb",
		"DROP TABLE IF EXISTS tu",
		"CREATE TABLE tu (a DATETIME, b DATETIME)",
		"DROP TABLE IF EXISTS tq",
		"CREATE TABLE tq (a INT, b VARCHAR(20), d DATE)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obprobe_qb")
	})
	for _, tc := range quotedBuiltinCallCases {
		if oceanBaseCrashCases[tc.sql] {
			continue
		}
		got, err := engineAccepts(ctx, db, tc.sql)
		if err != nil {
			t.Fatalf("the oracle gave no verdict on %q: %v", tc.sql, err)
		}
		if got != tc.accept {
			t.Errorf("OceanBase accepts %q = %v, want %v", tc.sql, got, tc.accept)
		}
		if _, perr := Parse(tc.sql); (perr == nil) != got {
			t.Errorf("DIVERGENCE omni accepts=%v, OceanBase accepts=%v for %q: %v", perr == nil, got, tc.sql, perr)
		}
	}
}
