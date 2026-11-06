package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The priority and query-cache modifiers as TiDB v8.5.5 reads them, every verdict
// measured with PREPARE on pingcap/tidb:v8.5.5 on 2026-11-06.
//
// Priority (parser.y:9046) is LOW_PRIORITY, HIGH_PRIORITY or DELAYED, all three
// reserved words. SELECT reads it as one SelectStmtOpt among the others, in any
// order and any number of times, and the last one wins (SelectStmtOptsList,
// parser.y:10240). UPDATE (UpdateStmtNoWith, parser.y:13626) and DELETE
// (DeleteWithoutUsingStmt and DeleteWithUsingStmt, parser.y:5231 and :5281) read
// at most one PriorityOpt, in a fixed place: UPDATE before IGNORE, DELETE before
// QUICK and IGNORE.
//
// SELECT also reads SQL_CACHE beside SQL_NO_CACHE (SelectStmtSQLCache,
// parser.y:10306); the statement skips the cache if any SQL_NO_CACHE is present.
// SQL_CACHE is unreserved (parser.y:623), so beside a dot it names a table.

type priorityCase struct {
	sql    string
	accept bool
}

var priorityCases = []priorityCase{
	// SELECT
	{"select low_priority * from t", true},
	{"select delayed * from t", true},
	{"select high_priority * from t", true},
	{"select low_priority high_priority delayed a from t", true},
	{"select distinct low_priority straight_join a from t", true},
	{"(select low_priority 1) union (select delayed 2)", true},
	{"select low_priority from t", false},
	{"select a from t where low_priority = 1", false},
	{"select SQL_CACHE * from t", true},
	{"select sql_cache sql_no_cache a from t", true},
	{"select sql_no_cache sql_cache a from t", true},
	{"select sql_cache.a from sql_cache", true},
	{"select sql_cache from t", false},
	{"select `sql_cache` from t", true},

	// UPDATE
	{"update low_priority t set a = 2", true},
	{"update high_priority t set a = 2", true},
	{"update delayed t set a = 2", true},
	{"update high_priority ignore t set a = 2", true},
	{"update high_priority t1, t2 set t1.a = 2", true},
	{"update ignore high_priority t set a = 2", false},
	{"update low_priority high_priority t set a = 2", false},

	// DELETE
	{"delete high_priority from t where a = 2", true},
	{"delete delayed from t where a = 2", true},
	{"delete low_priority quick ignore from t", true},
	{"delete delayed quick ignore t1 from t1, t2 where t1.a = t2.a", true},
	{"delete high_priority from t1 using t1, t2", true},
	{"with c as (select 1) delete high_priority from t", true},
	{"delete quick high_priority from t", false},
	{"delete high_priority low_priority from t", false},

	// INSERT and REPLACE already read Priority; these hold that in place.
	{"insert delayed ignore into t values (1)", true},
	{"replace high_priority into t values (1)", true},
	{"replace delayed into t values (1)", true},
}

// TestPriorityForms holds omni to the recorded verdicts without a server.
func TestPriorityForms(t *testing.T) {
	for _, c := range priorityCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbCorpusVersion, err)
			}
		})
	}
}

// TestPriorityTrees checks what the modifiers leave on the tree: the last
// priority of a SELECT wins, and SQL_NO_CACHE anywhere beats SQL_CACHE.
func TestPriorityTrees(t *testing.T) {
	first := func(t *testing.T, sql string) nodes.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		if len(list.Items) != 1 {
			t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
		}
		return list.Items[0]
	}
	selects := []struct {
		sql      string
		priority nodes.Priority
		noCache  bool
	}{
		{"select a from t", nodes.PriorityNone, false},
		{"select low_priority a from t", nodes.PriorityLow, false},
		{"select delayed a from t", nodes.PriorityDelayed, false},
		{"select low_priority high_priority a from t", nodes.PriorityHigh, false},
		{"select high_priority delayed a from t", nodes.PriorityDelayed, false},
		{"select sql_cache a from t", nodes.PriorityNone, false},
		{"select sql_cache sql_no_cache a from t", nodes.PriorityNone, true},
		{"select sql_no_cache sql_cache a from t", nodes.PriorityNone, true},
	}
	for _, c := range selects {
		t.Run(c.sql, func(t *testing.T) {
			s, ok := first(t, c.sql).(*nodes.SelectStmt)
			if !ok {
				t.Fatalf("Parse(%q) is not a SelectStmt", c.sql)
			}
			if s.Priority != c.priority || s.NoCache != c.noCache {
				t.Errorf("Parse(%q): Priority=%v NoCache=%v, want %v %v", c.sql, s.Priority, s.NoCache, c.priority, c.noCache)
			}
		})
	}
	others := []struct {
		sql      string
		priority nodes.Priority
	}{
		{"update t set a = 1", nodes.PriorityNone},
		{"update low_priority t set a = 1", nodes.PriorityLow},
		{"update high_priority t set a = 1", nodes.PriorityHigh},
		{"update delayed t set a = 1", nodes.PriorityDelayed},
		{"delete from t", nodes.PriorityNone},
		{"delete low_priority from t", nodes.PriorityLow},
		{"delete high_priority from t", nodes.PriorityHigh},
		{"delete delayed t1 from t1, t2", nodes.PriorityDelayed},
		{"replace high_priority into t values (1)", nodes.PriorityHigh},
	}
	for _, c := range others {
		t.Run(c.sql, func(t *testing.T) {
			var got nodes.Priority
			switch n := first(t, c.sql).(type) {
			case *nodes.UpdateStmt:
				got = n.Priority
			case *nodes.DeleteStmt:
				got = n.Priority
			case *nodes.InsertStmt:
				got = n.Priority
			default:
				t.Fatalf("Parse(%q) is a %T", c.sql, n)
			}
			if got != c.priority {
				t.Errorf("Parse(%q): Priority=%v, want %v", c.sql, got, c.priority)
			}
		})
	}
}
