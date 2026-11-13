package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestNamingTableNoKeyword(t *testing.T) {
	// A finding has no range, so its anchored text is empty.
	keyword := func(names ...string) []anchored {
		var want []anchored
		for _, name := range names {
			want = append(want, anchored{Code: v1review.NameIsKeyword, Name: name})
		}
		return want
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []anchored
	}{
		{name: "no keyword", sql: `CREATE TABLE "TABLE_NAME" (id INT);`},
		{name: "create", sql: `CREATE TABLE "CONNECT" (id INT);`, want: keyword("CONNECT")},
		{
			name: "a quoted name matches whatever its case and keeps it",
			sql:  "CREATE TABLE \"connect\" (id INT);\nCREATE TABLE \"Select\" (id INT);",
			want: keyword("connect", "Select"),
		},
		{
			name: "an unquoted name is reported upper-cased",
			sql:  "CREATE TABLE issue (id INT);\nCREATE TABLE Trigger (id INT);",
			want: keyword("ISSUE", "TRIGGER"),
		},
		{
			name: "the qualifier is not judged",
			sql:  "CREATE TABLE \"SELECT\".\"FROM\".t (id INT);\nCREATE TABLE db.sch.\"WHERE\" (id INT);",
			want: keyword("WHERE"),
		},
		{
			name: "near misses",
			sql:  "CREATE TABLE \"CONNECTS\" (id INT);\nCREATE TABLE \"SELECT \" (id INT);\nCREATE TABLE \"LIMIT\" (id INT);",
		},
		{
			name: "temporary, transient, hybrid and or replace tables",
			sql:  "CREATE TEMPORARY TABLE \"TO\" (id INT);\nCREATE TRANSIENT TABLE \"OF\" (id INT);\nCREATE HYBRID TABLE \"ON\" (id INT PRIMARY KEY);\nCREATE OR REPLACE TABLE \"OR\" (id INT);",
			want: keyword("TO", "OF", "ON", "OR"),
		},
		{
			name: "create as select, like and clone are not judged",
			sql:  "CREATE TABLE \"SELECT\" AS SELECT 1 AS a;\nCREATE TABLE \"SELECT\" (a INT) AS SELECT 1;\nCREATE TABLE \"SELECT\" LIKE t;\nCREATE TABLE \"SELECT\" CLONE t;",
		},
		{
			name: "alter rename judges only the new name",
			sql:  "ALTER TABLE \"SELECT\" RENAME TO t;\nALTER TABLE t RENAME TO db.sch.\"connect\";",
			want: keyword("connect"),
		},
		{
			name: "each rename of one statement is judged",
			sql:  `ALTER TABLE t RENAME TO "SELECT", RENAME COLUMN a TO "from", RENAME TO "x", RENAME TO "where";`,
			want: keyword("SELECT", "where"),
		},
		{
			name: "a swap, a column rename, a dynamic table and a view are not judged",
			sql:  "ALTER TABLE t SWAP WITH \"SELECT\";\nALTER TABLE t RENAME COLUMN a TO \"SELECT\";\nALTER DYNAMIC TABLE dt RENAME TO \"SELECT\";\nALTER VIEW v RENAME TO \"SELECT\";\nCREATE VIEW \"SELECT\" AS SELECT 1;",
		},
		{
			name: "an identifier literal is judged as written and never matches",
			sql:  "CREATE TABLE IDENTIFIER('SELECT') (id INT);\nALTER TABLE t RENAME TO IDENTIFIER('SELECT');",
		},
		{name: "a statement under EXPLAIN is not judged", sql: `EXPLAIN CREATE TABLE "SELECT" (id INT);`},
		{name: "a scripting block is not judged", sql: "BEGIN\n  CREATE TABLE \"SELECT\" (id INT);\nEND"},
		{
			name: "a task body is not judged",
			sql:  `CREATE TASK tk WAREHOUSE = wh SCHEDULE = '1 MINUTE' AS CREATE TABLE "SELECT" (id INT);`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := anchor(tt.sql, NamingTableNoKeyword(parse(t, tt.sql)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
