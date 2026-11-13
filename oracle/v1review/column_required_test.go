package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnRequired(t *testing.T) {
	missing := func(table, name, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnMissing, Name: name, Table: table, Text: text}
	}
	required := []string{"ID", "CREATED_TS", "a", "a!"}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{
			name: "create lacks columns, quoted order",
			sql:  "CREATE TABLE s.t (id int, \"a\" int);",
			want: []columnFinding{missing("T", "CREATED_TS", `CREATE TABLE s.t (id int, "a" int)`), missing("T", "a!", `CREATE TABLE s.t (id int, "a" int)`)},
		},
		{
			name: "unquoted names upper-cased",
			sql:  "CREATE TABLE t (id int, created_ts int, a int, \"a!\" int);",
			want: []columnFinding{missing("T", "a", `CREATE TABLE t (id int, created_ts int, a int, "a!" int)`)},
		},
		{
			name: "create as select lacks all",
			sql:  "CREATE TABLE t AS SELECT * FROM u;",
			want: []columnFinding{
				missing("T", "CREATED_TS", "CREATE TABLE t AS SELECT * FROM u"),
				missing("T", "ID", "CREATE TABLE t AS SELECT * FROM u"),
				missing("T", "a!", "CREATE TABLE t AS SELECT * FROM u"),
				missing("T", "a", "CREATE TABLE t AS SELECT * FROM u"),
			},
		},
		{
			name: "drop and rename",
			sql:  "ALTER TABLE t DROP COLUMN id;\nALTER TABLE t RENAME COLUMN created_ts TO x;\nALTER TABLE t RENAME COLUMN id TO id;",
			want: []columnFinding{missing("T", "ID", "ALTER TABLE t DROP COLUMN id"), missing("T", "CREATED_TS", "ALTER TABLE t RENAME COLUMN created_ts TO x")},
		},
		{name: "drop of other column", sql: "ALTER TABLE t DROP COLUMN x;"},
		{name: "nothing carried across statements", sql: "ALTER TABLE t DROP COLUMN x;\nALTER TABLE t ADD id int;"},
		{name: "other statements", sql: "DROP TABLE t;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnRequired(parseStmts(t, tt.sql), required))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
