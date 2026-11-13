package v1review

import "testing"

func TestColumnRequireCollation(t *testing.T) {
	runColumnCases(t, ColumnRequireCollation, []columnCase{
		{name: "collation named, other types", sql: "CREATE TABLE t(a VARCHAR(5) COLLATE utf8mb4_bin, b INT);"},
		{
			name: "create anchors the column",
			sql:  "CREATE TABLE t(a CHAR(1) CHARSET utf8mb4, b MEDIUMTEXT);",
			want: []string{
				"column-no-collation .a  @0: a CHAR(1) CHARSET utf8mb4",
				"column-no-collation .b  @0: b MEDIUMTEXT",
			},
		},
		{
			name: "add anchors the action, modify not judged",
			sql:  "ALTER TABLE t ADD a VARCHAR(5), MODIFY c CHAR(1);",
			want: []string{"column-no-collation .a  @0: ADD a VARCHAR(5)"},
		},
	})
}
