package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementMergeAlterTable(t *testing.T) {
	redundant := func(table string, count, stmt int) stmtFinding {
		return stmtFinding{Code: v1review.RedundantAlterTable, Table: table, Count: count, Statement: stmt}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []stmtFinding
	}{
		{name: "two alters", sql: "ALTER TABLE t ADD a INT;\nALTER TABLE t ADD b INT;", want: []stmtFinding{redundant("t", 2, 1)}},
		{name: "create then alter", sql: "CREATE TABLE t (a INT);\nSELECT 1;\nALTER TABLE t ADD b INT;", want: []stmtFinding{redundant("t", 2, 2)}},
		{name: "create resets count", sql: "ALTER TABLE t ADD a INT;\nALTER TABLE t ADD b INT;\nCREATE TABLE t (x INT);"},
		{name: "bare name keys tables", sql: "ALTER TABLE db1.t ADD a INT;\nALTER TABLE db2.t ADD b INT;", want: []stmtFinding{redundant("t", 2, 1)}},
		{name: "case-sensitive", sql: "ALTER TABLE T ADD a INT;\nALTER TABLE t ADD b INT;"},
		{
			name: "ordered by last statement",
			sql:  "ALTER TABLE u ADD a INT;\nALTER TABLE t ADD a INT;\nALTER TABLE t ADD b INT;\nALTER TABLE u ADD b INT;",
			want: []stmtFinding{redundant("t", 2, 2), redundant("u", 2, 3)},
		},
		{name: "routine body not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t ADD a INT; END;\nALTER TABLE t ADD b INT;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.want, runStatement(t, tt.sql, StatementMergeAlterTable)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}
