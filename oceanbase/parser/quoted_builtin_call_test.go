package parser

import (
	"context"
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// quotedBuiltinCallCases are calls whose name is a backtick-quoted builtin with a
// keyword argument. The grammar reads TIMESTAMPDIFF, TIMESTAMPADD, GET_FORMAT and
// WEIGHT_STRING from their keywords, so a quoted name is an identifier, the call
// is an ordinary one, and MySQL 8.0 parses each accepted case and answers 1630
// for the stored function it then cannot find. A rejected case is one an ordinary
// argument list cannot hold. 5.7.25-OceanBase_CE-v4.3.5.6 gives the same verdict
// on every case (TestOceanBaseQuotedBuiltinCall). tu has DATETIME columns a and b.
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

// TestQuotedBuiltinCallOracle asserts omni agrees with a live mysql:8.0 on
// quotedBuiltinCallCases (container test).
func TestQuotedBuiltinCallOracle(t *testing.T) {
	o := startParserOracle(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE tu (a DATETIME, b DATETIME)"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range quotedBuiltinCallCases {
		if got := o.checkSyntax(t, tc.sql); got != tc.accept {
			t.Errorf("MySQL accepts %q = %v, want %v", tc.sql, got, tc.accept)
		}
		assertOracleMatch(t, o, tc.sql)
	}
}

// TestOceanBaseQuotedBuiltinCall asks OceanBase about quotedBuiltinCallCases. It
// reads OCEANBASE_WALK_DSN and skips without it, like the rest of the OceanBase
// oracle:
//
//	scripts/oceanbase-oracle.sh go test ./oceanbase/parser/ -run OceanBase -count=1
func TestOceanBaseQuotedBuiltinCall(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS obprobe_qb",
		"USE obprobe_qb",
		"DROP TABLE IF EXISTS tu",
		"CREATE TABLE tu (a DATETIME, b DATETIME)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obprobe_qb")
	})
	for _, tc := range quotedBuiltinCallCases {
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
