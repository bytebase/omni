package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// FETCH FIRST as TiDB v8.5.5 reads it, every verdict measured with PREPARE on
// pingcap/tidb:v8.5.5 on 2026-11-06.
//
// "FETCH" FirstOrNext FetchFirstOpt RowOrRows "ONLY" (parser.y:10151) is one of
// the four spellings of SelectStmtLimit, so it stands wherever LIMIT does in a
// query and nowhere else: not after LIMIT, not in UPDATE or DELETE, whose
// LimitClause is LIMIT alone. The count is a LimitOption, as LIMIT's is, and one
// when left out. There is no OFFSET form.

type fetchFirstCase struct {
	sql    string
	accept bool
}

var fetchFirstCases = []fetchFirstCase{
	{"select * from t fetch first 5 rows only", true},
	{"select * from t fetch next 5 row only", true},
	{"select * from t fetch first row only", true},
	{"select * from t fetch next rows only", true},
	{"select * from t order by a fetch first 2 rows only", true},
	{"select * from t fetch first 1 rows only for update", true},
	{"select * from (select * from t fetch first 1 rows only) x", true},
	{"(select * from t fetch first 1 rows only) union (select * from t)", true},
	{"select * from t union select * from t fetch first 1 rows only", true},
	{"(select * from t) fetch first 1 rows only", true},
	{"table t fetch first 1 rows only", true},
	{"values row(1) fetch first 1 rows only", true},
	{"select 1 from dual fetch first 1 rows only", true},
	{"select * from ((select * from t) fetch first 1 rows only) x", true},
	{"select 1 in ((select 1) fetch first 1 rows only)", true},
	{"select (select 1 fetch first 1 rows only)", true},
	{"((select 1) fetch first 1 rows only)", true},
	{"select * from t where a in (select a from t fetch first 1 rows only)", true},
	{"select * from t limit 1 fetch first 1 rows only", false},
	{"select * from t fetch first 1 rows only limit 1", false},
	{"select * from t fetch first 5 rows only offset 1", false},
	{"select * from t offset 1 rows fetch first 5 rows only", false},
	{"select * from t fetch first 5 only", false},
	{"select * from t fetch first 5 rows", false},
	{"select * from t fetch 5 rows only", false},
	{"select * from t fetch first 1.5 rows only", false},
	{"select * from t fetch first -1 rows only", false},
	{"select * from t fetch first (1) rows only", false},
	{"select * from t fetch first a rows only", false},
	{"update t set a = 1 fetch first 1 rows only", false},
	{"delete from t fetch first 1 rows only", false},
}

// TestFetchFirstForms holds omni to the recorded verdicts without a server.
func TestFetchFirstForms(t *testing.T) {
	for _, c := range fetchFirstCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbRelease, err)
			}
		})
	}
}

// TestFetchFirstTrees checks FETCH FIRST builds the Limit LIMIT would, with a
// count of one when it is left out, and spans the clause.
func TestFetchFirstTrees(t *testing.T) {
	for _, c := range []struct {
		sql   string
		count int64
		text  string
	}{
		{"select * from t fetch first 5 rows only", 5, "fetch first 5 rows only"},
		{"select * from t fetch next row only", 1, "fetch next row only"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			lim := list.Items[0].(*nodes.SelectStmt).Limit
			if lim == nil || lim.Offset != nil {
				t.Fatalf("Parse(%q).Limit = %#v, want a count and no offset", c.sql, lim)
			}
			n, ok := lim.Count.(*nodes.IntLit)
			if !ok || n.Value != c.count {
				t.Errorf("Parse(%q) count = %#v, want %d", c.sql, lim.Count, c.count)
			}
			if got := c.sql[lim.Loc.Start:lim.Loc.End]; got != c.text {
				t.Errorf("Loc covers %q, want %q", got, c.text)
			}
		})
	}
}
