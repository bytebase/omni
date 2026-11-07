package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Window null treatment and FROM FIRST/LAST as TiDB v8.5.5 reads them, every
// verdict measured with PREPARE on pingcap/tidb:v8.5.5 on 2026-11-07.
//
// OptNullTreatment (RESPECT|IGNORE NULLS, parser.y:9845) follows the argument list
// of LEAD, LAG, FIRST_VALUE, LAST_VALUE and NTH_VALUE, before OVER; OptFromFirstLast
// (FROM FIRST|LAST, parser.y:9858) follows only NTH_VALUE's, before the null
// treatment. Both belong to those WindowFuncCall arms, so each requires OVER and no
// other function takes them.

type windowNullCase struct {
	sql    string
	accept bool
}

var windowNullCases = []windowNullCase{
	{"SELECT FIRST_VALUE(val) RESPECT NULLS OVER w FROM t", true},
	{"SELECT FIRST_VALUE(val) IGNORE NULLS OVER w FROM t", true},
	{"SELECT LAST_VALUE(val) IGNORE NULLS OVER w FROM t", true},
	{"SELECT LEAD(val) IGNORE NULLS OVER w FROM t", true},
	{"SELECT LEAD(val, 2) RESPECT NULLS OVER w FROM t", true},
	{"SELECT LAG(val, 2, 0) IGNORE NULLS OVER w FROM t", true},
	{"SELECT NTH_VALUE(val, 233) FROM FIRST OVER w FROM t", true},
	{"SELECT NTH_VALUE(val, 233) FROM LAST OVER w FROM t", true},
	{"SELECT NTH_VALUE(val, 233) FROM LAST IGNORE NULLS OVER w FROM t", true},
	{"SELECT NTH_VALUE(val, 233) FROM FIRST RESPECT NULLS OVER w FROM t", true},
	{"SELECT NTH_VALUE(val, 1) IGNORE NULLS OVER w FROM t", true},
	{"SELECT FIRST_VALUE(val) IGNORE NULLS OVER () FROM t", true},
	// Null treatment and FROM FIRST/LAST require OVER.
	{"SELECT FIRST_VALUE(val) IGNORE NULLS FROM t", false},
	{"SELECT NTH_VALUE(val, 1) FROM LAST FROM t", false},
	// Wrong order: null treatment before FROM FIRST/LAST.
	{"SELECT NTH_VALUE(val, 1) IGNORE NULLS FROM LAST OVER w FROM t", false},
	// FROM FIRST/LAST is NTH_VALUE's alone.
	{"SELECT FIRST_VALUE(val) FROM LAST OVER w FROM t", false},
	{"SELECT LEAD(val) FROM FIRST OVER w FROM t", false},
	// Null treatment is these five functions' alone.
	{"SELECT COUNT(val) IGNORE NULLS OVER w FROM t", false},
	{"SELECT SUM(val) RESPECT NULLS OVER w FROM t", false},
	{"SELECT ROW_NUMBER() IGNORE NULLS OVER w FROM t", false},
	// Only NULLS, and only one word before it.
	{"SELECT FIRST_VALUE(val) IGNORE OVER w FROM t", false},
	{"SELECT FIRST_VALUE(val) RESPECT OVER w FROM t", false},
	{"SELECT FIRST_VALUE(val) IGNORE NULL OVER w FROM t", false},
	// An ordinary (quoted) function name takes none of this.
	{"SELECT `first_value`(val) IGNORE NULLS OVER w FROM t", false},
}

// TestWindowNullForms holds omni to the recorded verdicts without a server.
func TestWindowNullForms(t *testing.T) {
	for _, c := range windowNullCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbCorpusVersion, err)
			}
		})
	}
}

// TestWindowNullTrees checks the flags each clause sets, and that the defaults
// RESPECT NULLS and FROM FIRST leave them false.
func TestWindowNullTrees(t *testing.T) {
	item := func(t *testing.T, sql string) *nodes.FuncCallExpr {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		fc, ok := list.Items[0].(*nodes.SelectStmt).TargetList[0].(*nodes.FuncCallExpr)
		if !ok {
			t.Fatalf("Parse(%q): item is not a FuncCallExpr", sql)
		}
		return fc
	}
	for _, c := range []struct {
		sql              string
		ignore, fromLast bool
	}{
		{"SELECT FIRST_VALUE(v) OVER w FROM t", false, false},
		{"SELECT FIRST_VALUE(v) RESPECT NULLS OVER w FROM t", false, false},
		{"SELECT FIRST_VALUE(v) IGNORE NULLS OVER w FROM t", true, false},
		{"SELECT LEAD(v) IGNORE NULLS OVER w FROM t", true, false},
		{"SELECT NTH_VALUE(v, 1) FROM FIRST OVER w FROM t", false, false},
		{"SELECT NTH_VALUE(v, 1) FROM LAST OVER w FROM t", false, true},
		{"SELECT NTH_VALUE(v, 1) FROM LAST IGNORE NULLS OVER w FROM t", true, true},
	} {
		t.Run(c.sql, func(t *testing.T) {
			fc := item(t, c.sql)
			if fc.IgnoreNull != c.ignore || fc.FromLast != c.fromLast {
				t.Errorf("Parse(%q): IgnoreNull=%v FromLast=%v, want %v %v", c.sql, fc.IgnoreNull, fc.FromLast, c.ignore, c.fromLast)
			}
		})
	}
}
