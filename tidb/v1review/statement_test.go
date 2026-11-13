package v1review

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// stmtFinding is a statement rule's finding as the tests compare it: the
// finding with its range replaced by the text the range anchors.
type stmtFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
	Text      string
}

// stmtCase is one statement rule test: sql and the findings it yields.
type stmtCase struct {
	name string
	sql  string
	want []stmtFinding
}

// runStatementCases runs a statement rule over each case's statements.
func runStatementCases(t *testing.T, rule func([]ast.Node) []v1review.Finding, cases []stmtCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var got []stmtFinding
			for _, f := range rule(parse(t, tt.sql)) {
				got = append(got, stmtFinding{f.Code, f.Name, f.Table, f.Value, f.Statement, tt.sql[f.Start:f.End]})
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStmtAnchor(t *testing.T) {
	for _, tt := range []struct {
		loc        ast.Loc
		start, end int
	}{
		{ast.Loc{Start: 3, End: 9}, 3, 9},
		{ast.Loc{Start: 3, End: -1}, 3, 3},
		{ast.Loc{Start: 0, End: 0}, 0, 0},
		{ast.Loc{Start: -1, End: 5}, 0, 0},
	} {
		if start, end := stmtAnchor(tt.loc); start != tt.start || end != tt.end {
			t.Errorf("stmtAnchor(%+v) = %d,%d, want %d,%d", tt.loc, start, end, tt.start, tt.end)
		}
	}
}

func TestStatementDisallowCommit(t *testing.T) {
	commit := stmtFinding{Code: v1review.Commit}
	runStatementCases(t, StatementDisallowCommit, []stmtCase{
		{name: "commit", sql: "COMMIT", want: []stmtFinding{commit}},
		{name: "two commits", sql: "BEGIN; COMMIT; COMMIT", want: []stmtFinding{commit, commit}},
		{name: "rollback not judged", sql: "ROLLBACK"},
		{name: "in a routine body not judged", sql: "CREATE PROCEDURE p() BEGIN COMMIT; END"},
	})
}

func TestStatementInsertMustSpecifyColumn(t *testing.T) {
	missing := func(text string) stmtFinding {
		return stmtFinding{Code: v1review.InsertColumnsUnspecified, Text: text}
	}
	runStatementCases(t, StatementInsertMustSpecifyColumn, []stmtCase{
		{name: "values", sql: "INSERT INTO t VALUES (1)", want: []stmtFinding{missing("INSERT INTO t VALUES (1)")}},
		{name: "empty column list", sql: "INSERT INTO t () VALUES ()", want: []stmtFinding{missing("INSERT INTO t () VALUES ()")}},
		{name: "select", sql: "REPLACE INTO t SELECT * FROM u", want: []stmtFinding{missing("REPLACE INTO t SELECT * FROM u")}},
		{name: "columns", sql: "INSERT INTO t (a) VALUES (1)"},
		{name: "set form not judged", sql: "INSERT INTO t SET a = 1"},
	})
}

func TestStatementInsertDisallowOrderByRand(t *testing.T) {
	rand := stmtFinding{Code: v1review.InsertOrderByRand}
	runStatementCases(t, StatementInsertDisallowOrderByRand, []stmtCase{
		{name: "order by rand", sql: "INSERT INTO t SELECT a FROM u ORDER BY rand()", want: []stmtFinding{rand}},
		{name: "once per statement", sql: "INSERT INTO t SELECT a FROM u ORDER BY a, RAND(), RAND(1)", want: []stmtFinding{rand}},
		{name: "rand in an expression not judged", sql: "INSERT INTO t SELECT a FROM u ORDER BY RAND() + 1"},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT a FROM (SELECT a FROM u ORDER BY RAND()) x"},
		{name: "select not judged", sql: "SELECT a FROM u ORDER BY RAND()"},
	})
}

func TestStatementMergeAlterTable(t *testing.T) {
	redundant := func(table string, count, statement int) stmtFinding {
		return stmtFinding{Code: v1review.RedundantAlterTable, Table: table, Value: strconv.Itoa(count), Statement: statement}
	}
	runStatementCases(t, StatementMergeAlterTable, []stmtCase{
		{
			name: "create then alters",
			sql:  "CREATE TABLE t (a INT); ALTER TABLE t ADD b INT; ALTER TABLE u ADD c INT; ALTER TABLE db.t ADD d INT",
			want: []stmtFinding{redundant("t", 3, 3)},
		},
		{
			name: "create resets",
			sql:  "ALTER TABLE t ADD a INT; ALTER TABLE t ADD b INT; DROP TABLE t; CREATE TABLE t (a INT); ALTER TABLE t ADD c INT",
			want: []stmtFinding{redundant("t", 2, 4)},
		},
		{
			name: "names are case-sensitive, findings in last-statement order",
			sql:  "ALTER TABLE u ADD a INT; ALTER TABLE t ADD a INT; ALTER TABLE T ADD b INT; ALTER TABLE t ADD c INT; ALTER TABLE u ADD d INT",
			want: []stmtFinding{redundant("t", 2, 3), redundant("u", 2, 4)},
		},
		{name: "one alter", sql: "ALTER TABLE t ADD a INT"},
	})
}

func TestStatementDisallowLimit(t *testing.T) {
	limit := func(kind string) stmtFinding { return stmtFinding{Code: v1review.DMLLimit, Name: kind} }
	runStatementCases(t, StatementDisallowLimit, []stmtCase{
		{name: "update", sql: "UPDATE t SET a = 1 LIMIT 1", want: []stmtFinding{limit("UPDATE")}},
		{name: "delete", sql: "DELETE FROM t LIMIT 1", want: []stmtFinding{limit("DELETE")}},
		{name: "insert select", sql: "INSERT INTO t SELECT a FROM u LIMIT 1", want: []stmtFinding{limit("INSERT")}},
		{name: "subquery not judged", sql: "INSERT INTO t SELECT a FROM (SELECT a FROM u LIMIT 1) x"},
		{name: "select not judged", sql: "SELECT a FROM u LIMIT 1"},
	})
}

func TestStatementDisallowOrderBy(t *testing.T) {
	orderBy := func(kind string) stmtFinding { return stmtFinding{Code: v1review.DMLOrderBy, Name: kind} }
	runStatementCases(t, StatementDisallowOrderBy, []stmtCase{
		{name: "update", sql: "UPDATE t SET a = 1 ORDER BY a", want: []stmtFinding{orderBy("UPDATE")}},
		{name: "delete", sql: "DELETE FROM t ORDER BY a", want: []stmtFinding{orderBy("DELETE")}},
		{name: "insert select not judged", sql: "INSERT INTO t SELECT a FROM u ORDER BY a"},
	})
}

func TestStatementWhereNoLeadingWildcardLike(t *testing.T) {
	like := stmtFinding{Code: v1review.LeadingWildcardLike}
	runStatementCases(t, StatementWhereNoLeadingWildcardLike, []stmtCase{
		{name: "leading wildcard", sql: "SELECT a FROM t WHERE a LIKE '%x'", want: []stmtFinding{like}},
		{name: "once per statement", sql: "SELECT a FROM t WHERE a LIKE '%x' OR b NOT LIKE '%y'", want: []stmtFinding{like}},
		{name: "collate unwrapped", sql: "SELECT a FROM t WHERE a LIKE '%x' COLLATE utf8mb4_bin", want: []stmtFinding{like}},
		{name: "trailing wildcard", sql: "SELECT a FROM t WHERE a LIKE 'x%'"},
		{name: "non-literal not judged", sql: "SELECT a FROM t WHERE a LIKE CONCAT('%', b)"},
	})
}

func TestStatementSelectNoSelectAll(t *testing.T) {
	star := stmtFinding{Code: v1review.SelectAll}
	runStatementCases(t, StatementSelectNoSelectAll, []stmtCase{
		{name: "star", sql: "SELECT * FROM t", want: []stmtFinding{star}},
		{name: "once per select", sql: "SELECT *, t.* FROM t", want: []stmtFinding{star}},
		{name: "each select", sql: "SELECT * FROM (SELECT t.* FROM t) x", want: []stmtFinding{star, star}},
		{name: "count star not judged", sql: "SELECT COUNT(*) FROM t"},
	})
}

func TestStatementWhereRequireSelect(t *testing.T) {
	missing := stmtFinding{Code: v1review.WhereMissing}
	runStatementCases(t, StatementWhereRequireSelect, []stmtCase{
		{name: "no where", sql: "SELECT a FROM t", want: []stmtFinding{missing}},
		{name: "subquery", sql: "SELECT a FROM t WHERE a IN (SELECT a FROM u)", want: []stmtFinding{missing}},
		{name: "no from not judged", sql: "SELECT 1"},
		{name: "explain", sql: "EXPLAIN SELECT a FROM t", want: []stmtFinding{missing}},
	})
}

func TestStatementWhereRequireUpdateDelete(t *testing.T) {
	missing := stmtFinding{Code: v1review.WhereMissing}
	runStatementCases(t, StatementWhereRequireUpdateDelete, []stmtCase{
		{name: "update", sql: "UPDATE t SET a = 1", want: []stmtFinding{missing}},
		{name: "delete", sql: "DELETE FROM t", want: []stmtFinding{missing}},
		{name: "with where", sql: "DELETE FROM t WHERE a = 1"},
		{name: "explain", sql: "EXPLAIN DELETE FROM t", want: []stmtFinding{missing}},
	})
}
