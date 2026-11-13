package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestTableRequirePK(t *testing.T) {
	missing := func(name string, stmt int, text string) stmtFinding {
		return stmtFinding{v1review.TablePrimaryKeyMissing, name, stmt, text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "column primary key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{name: "table primary key", sql: "CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a));"},
		{name: "missing anchors the table name", sql: "CREATE TABLE [dbo].[t] (a INT UNIQUE);", want: []stmtFinding{missing("dbo.t", 0, "[dbo].[t]")}},
		{name: "added constraint", sql: "CREATE TABLE t (a INT); ALTER TABLE dbo.t ADD CONSTRAINT pk PRIMARY KEY (a);"},
		{name: "added column's constraint", sql: "CREATE TABLE t (a INT); ALTER TABLE t ADD b INT PRIMARY KEY;"},
		{
			name: "a later create starts the table again",
			sql:  "ALTER TABLE t ADD CONSTRAINT pk PRIMARY KEY (a); CREATE TABLE t (a INT);",
			want: []stmtFinding{missing("t", 1, "t")},
		},
		{name: "dropping the constraint does not remove it", sql: "CREATE TABLE t (a INT PRIMARY KEY); ALTER TABLE t DROP CONSTRAINT pk;"},
		{name: "dropping the table does not forget it", sql: "CREATE TABLE t (a INT); DROP TABLE t;", want: []stmtFinding{missing("t", 0, "t")}},
		{name: "uncreated table not judged", sql: "ALTER TABLE t ADD b INT;"},
		{name: "select into not judged", sql: "SELECT * INTO t FROM u;"},
		{name: "nested not read", sql: "IF 1 = 1 CREATE TABLE t (a INT);"},
		{
			name: "last create anchors",
			sql:  "CREATE TABLE t (a INT); CREATE TABLE dbo.t (b INT);",
			want: []stmtFinding{missing("dbo.t", 1, "dbo.t")},
		},
		{
			name: "tables in the order first tracked",
			sql:  "ALTER TABLE c ADD CONSTRAINT pk PRIMARY KEY (a); CREATE TABLE b (a INT); CREATE TABLE a (a INT); CREATE TABLE c (a INT);",
			want: []stmtFinding{missing("c", 3, "c"), missing("b", 1, "b"), missing("a", 2, "a")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := runStmtFindings(t, tt.sql, TableRequirePK)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
