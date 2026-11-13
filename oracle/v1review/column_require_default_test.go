package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestColumnRequireDefault(t *testing.T) {
	none := func(name string, statement int, text string) columnFinding {
		return columnFinding{Code: v1review.ColumnNoDefault, Name: name, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []columnFinding
	}{
		{
			name: "create columns in key order",
			sql:  "CREATE TABLE t (b int, a int DEFAULT 1, c int GENERATED ALWAYS AS IDENTITY);",
			want: []columnFinding{none("B", 0, "b int"), none("C", 0, "c int GENERATED ALWAYS AS IDENTITY")},
		},
		{name: "modify with default forgets", sql: "ALTER TABLE t ADD a int;\nALTER TABLE t MODIFY a DEFAULT 0;"},
		{
			name: "modify without default leaves it",
			sql:  "ALTER TABLE t ADD a int;\nALTER TABLE t MODIFY a varchar2(10);\nALTER TABLE t DROP COLUMN a;",
			want: []columnFinding{none("A", 0, "a int")},
		},
		{name: "modify alone not judged", sql: "ALTER TABLE t MODIFY a int;"},
		{
			name: "re-added column anchors the later statement",
			sql:  "CREATE TABLE t (a int);\nALTER TABLE t ADD (a int, \"x.y\" int);",
			want: []columnFinding{none("A", 1, "a int"), none("y", 1, `"x.y" int`)},
		},
		{
			name: "current database qualifies keys",
			sql:  "CREATE TABLE t (a int);\nALTER TABLE db.t MODIFY a DEFAULT 1;\nCREATE TABLE s.t (a int);",
			want: []columnFinding{none("A", 2, "a int")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := columnFindings(tt.sql, ColumnRequireDefault(parseStmts(t, tt.sql), "DB"))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
