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
	const create = "CREATE TABLE t (a INT)"
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "inline primary key", sql: "CREATE TABLE t (a INT PRIMARY KEY);"},
		{name: "table primary key", sql: "CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a));"},
		{name: "missing anchors the statement", sql: create + ";", want: []stmtFinding{missing("t", 0, create)}},
		{name: "as select, like and clone not read", sql: "CREATE TABLE t LIKE u; CREATE TABLE t CLONE u; CREATE TABLE t AS SELECT 1 AS a;"},
		{name: "added constraint", sql: create + "; ALTER TABLE PUBLIC.T ADD PRIMARY KEY (a);"},
		{name: "added column's primary key not read", sql: create + "; ALTER TABLE t ADD COLUMN b INT PRIMARY KEY;", want: []stmtFinding{missing("t", 0, create)}},
		{
			name: "drop primary key anchors the action and renames",
			sql:  "CREATE TABLE t (a INT PRIMARY KEY); ALTER TABLE \"T\" DROP PRIMARY KEY;",
			want: []stmtFinding{missing(`"T"`, 1, "DROP PRIMARY KEY")},
		},
		{name: "drop constraint by name keeps it", sql: "CREATE TABLE t (a INT PRIMARY KEY); ALTER TABLE t DROP CONSTRAINT pk;"},
		{name: "drop primary key of an unread table", sql: "ALTER TABLE t DROP PRIMARY KEY;"},
		{
			name: "added then dropped on an uncreated table",
			sql:  "ALTER TABLE t ADD PRIMARY KEY (a); ALTER TABLE t DROP PRIMARY KEY;",
			want: []stmtFinding{missing("t", 1, "DROP PRIMARY KEY")},
		},
		{name: "drop table forgets", sql: create + "; DROP TABLE T;"},
		{name: "drop view does not", sql: create + "; DROP VIEW t;", want: []stmtFinding{missing("t", 0, create)}},
		{
			name: "a table read again after a drop table is new",
			sql:  "ALTER TABLE \"T\" RENAME CONSTRAINT x TO y; DROP TABLE t; ALTER TABLE t ADD PRIMARY KEY (a); ALTER TABLE t DROP PRIMARY KEY;",
			want: []stmtFinding{missing("t", 3, "DROP PRIMARY KEY")},
		},
		{
			name: "tables in the order first read, a dropped table read again as new",
			sql:  "CREATE TABLE b (a INT); CREATE TABLE a (a INT); CREATE TABLE c (a INT); DROP TABLE b; CREATE TABLE b (a INT);",
			want: []stmtFinding{missing("a", 1, "CREATE TABLE a (a INT)"), missing("c", 2, "CREATE TABLE c (a INT)"), missing("b", 4, "CREATE TABLE b (a INT)")},
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
