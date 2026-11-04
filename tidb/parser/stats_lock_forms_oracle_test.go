package parser

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// LOCK STATS and UNLOCK STATS, which freeze and release a table's statistics.
// They are TiDB's own statements: MySQL, MariaDB and OceanBase answer 1064 to
// every shape, and each engine's own tests pin that refusal.
//
//	{LOCK | UNLOCK} STATS tbl_name [, tbl_name] ...
//	{LOCK | UNLOCK} STATS tbl_name PARTITION [(] p_name [, p_name] ... [)]
//
// Two things about the shape needed asking rather than reading: the partition
// arm takes one table and not a list, and its parentheses are optional. The
// corpus writes six of these statements and nothing malformed, so the
// rejections are the point of this file. Every case was measured against
// pingcap/tidb:v8.5.5 on 2026-11-04, and the expectation is what that server
// answered, so an image upgrade that changes a verdict fails here instead of
// passing silently.
//
// Ref: LockStatsStmt and UnlockStatsStmt in pingcap/tidb pkg/parser/parser.y.
// https://docs.pingcap.com/tidb/stable/sql-statement-lock-stats/

// statsLockCase is a statement and whether TiDB's parser accepts it.
type statsLockCase struct {
	sql    string
	accept bool
}

var statsLockCases = []statsLockCase{
	// A list of tables, qualified or not, quoted or not.
	{"LOCK STATS t", true},
	{"UNLOCK STATS t", true},
	{"lock stats test.t", true},
	{"unlock stats test.t", true},
	{"LOCK STATS t, t2", true},
	{"UNLOCK STATS t, t2", true},
	{"LOCK STATS test.t, test.t2", true},
	{"LOCK STATS `test`.`t`, t2", true},
	// One table and its partitions, with the parentheses and without them.
	{"LOCK STATS t PARTITION (p0)", true},
	{"LOCK STATS t PARTITION (p0, p1)", true},
	{"LOCK STATS t PARTITION p0", true},
	{"LOCK STATS t PARTITION p0, p1", true},
	{"UNLOCK STATS t PARTITION (p0, p1)", true},
	{"UNLOCK STATS t PARTITION p0", true},
	{"LOCK STATS test.t PARTITION p0", true},
	{"LOCK STATS t PARTITION `p0`", true},
	// The partition arm reads one table, so a list in front of PARTITION is a
	// syntax error rather than a partition list that applies to all of them.
	{"LOCK STATS t, t2 PARTITION (p0)", false},
	{"LOCK STATS t, t2 PARTITION p0", false},
	{"LOCK STATS t PARTITION (p0), t2", false},
	{"LOCK STATS t PARTITION (p0), t2 PARTITION (p1)", false},
	// Neither list is optional, and neither takes a trailing comma, an empty
	// list, a doubled parenthesis or a second PARTITION.
	{"LOCK STATS", false},
	{"UNLOCK STATS", false},
	{"LOCK STATS PARTITION (p0)", false},
	{"LOCK STATS t,", false},
	{"LOCK STATS t PARTITION", false},
	{"LOCK STATS t PARTITION ()", false},
	{"LOCK STATS t PARTITION (p0,)", false},
	{"LOCK STATS t PARTITION ((p0))", false},
	{"LOCK STATS t PARTITION (p0) (p1)", false},
	{"LOCK STATS t PARTITION (p0) PARTITION (p1)", false},
	{"LOCK STATS t PARTITION p0 p1", false},
	// A table name is an identifier: no number, no string, no reserved word.
	// A partition name is one too, which PARTITION itself is not.
	{"LOCK STATS = t", false},
	{"UNLOCK STATS = t", false},
	{"LOCK STATS 1", false},
	{"LOCK STATS 't'", false},
	{"LOCK STATS select", false},
	{"LOCK STATS t PARTITION partition", false},
	// Nothing follows the targets, and the word is STATS and only STATS.
	{"LOCK STATS t READ", false},
	{"LOCK STAT t", false},
	{"LOCK STATISTICS t", false},
	{"UNLOCK STAT t", false},
	// STATS is a token in TiDB's lexer, so a backquoted one is an ordinary
	// table name and the statement reads as a LOCK TABLES that never names
	// what it locks.
	{"LOCK STATS `t`", true},
	{"LOCK `STATS` t", false},
	{"UNLOCK `STATS` t", false},
}

// TestStatsLockFormsOracle asserts omni agrees with the pinned TiDB on
// statsLockCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection, and
// any other server error means the text parsed and failed later.
func TestStatsLockFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, c := range statsLockCases {
		t.Run(c.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = statsLockCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, c.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != c.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, c.accept, c.sql)
			}
			_, perr := Parse(c.sql)
			if omniAccepts := perr == nil; omniAccepts != c.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, c.accept, c.sql, perr)
			}
		})
	}
}

// TestStatsLockInTree asserts what the statement leaves in the tree, which is
// what a parse-only verdict cannot see: before this grammar landed,
// `LOCK STATS t` parsed as a LOCK TABLES of a table called `t` and the corpus
// recorded it as an acceptance.
func TestStatsLockInTree(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		unlock     bool
		tables     []string
		partitions []string
	}{
		{sql: "LOCK STATS t", tables: []string{"t"}},
		{sql: "UNLOCK STATS t", unlock: true, tables: []string{"t"}},
		{sql: "lock stats test.t", tables: []string{"test.t"}},
		{sql: "LOCK STATS t, test.t2, `t3`", tables: []string{"t", "test.t2", "t3"}},
		{sql: "UNLOCK STATS t, t2", unlock: true, tables: []string{"t", "t2"}},
		{sql: "LOCK STATS t PARTITION (p0)", tables: []string{"t"}, partitions: []string{"p0"}},
		{sql: "LOCK STATS t PARTITION (p0, p1)", tables: []string{"t"}, partitions: []string{"p0", "p1"}},
		{sql: "LOCK STATS t PARTITION p0, p1", tables: []string{"t"}, partitions: []string{"p0", "p1"}},
		{sql: "UNLOCK STATS test.t PARTITION `p0`", unlock: true, tables: []string{"test.t"}, partitions: []string{"p0"}},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) read %d statements, want 1", tc.sql, len(list.Items))
			}
			var tables []*ast.TableRef
			var partitions []string
			switch n := list.Items[0].(type) {
			case *ast.LockStatsStmt:
				if tc.unlock {
					t.Fatalf("Parse(%q) read a LockStatsStmt, want an UnlockStatsStmt", tc.sql)
				}
				tables, partitions = n.Tables, n.Partitions
			case *ast.UnlockStatsStmt:
				if !tc.unlock {
					t.Fatalf("Parse(%q) read an UnlockStatsStmt, want a LockStatsStmt", tc.sql)
				}
				tables, partitions = n.Tables, n.Partitions
			default:
				t.Fatalf("Parse(%q) read a %T, want a stats lock statement", tc.sql, list.Items[0])
			}
			got := make([]string, 0, len(tables))
			for _, ref := range tables {
				if ref.Schema != "" {
					got = append(got, ref.Schema+"."+ref.Name)
					continue
				}
				got = append(got, ref.Name)
			}
			if len(got) != len(tc.tables) {
				t.Fatalf("tables = %q, want %q", got, tc.tables)
			}
			for i := range got {
				if got[i] != tc.tables[i] {
					t.Errorf("table %d = %q, want %q", i, got[i], tc.tables[i])
				}
			}
			if len(partitions) != len(tc.partitions) {
				t.Fatalf("partitions = %q, want %q", partitions, tc.partitions)
			}
			for i := range partitions {
				if partitions[i] != tc.partitions[i] {
					t.Errorf("partition %d = %q, want %q", i, partitions[i], tc.partitions[i])
				}
			}
		})
	}
}
