package v1review

import "testing"

func TestColumnRequireCharset(t *testing.T) {
	runColumnCases(t, ColumnRequireCharset, []columnCase{
		{name: "charset named, other types", sql: "CREATE TABLE t(a VARCHAR(5) CHARSET utf8mb4, b INT, c BLOB, d ENUM('x'));"},
		{
			name: "create anchors the column",
			sql:  "CREATE TABLE t(a CHAR(1), b text COLLATE utf8mb4_bin);",
			want: []string{
				"column-no-charset .a  @0: a CHAR(1)",
				"column-no-charset .b  @0: b text COLLATE utf8mb4_bin",
			},
		},
		{
			name: "add anchors the action, modify and change not judged",
			sql:  "ALTER TABLE t ADD COLUMN (a LONGTEXT, b TINYTEXT), MODIFY c CHAR(1), CHANGE d e CHAR(1);",
			want: []string{
				"column-no-charset .a  @0: ADD COLUMN (a LONGTEXT, b TINYTEXT)",
				"column-no-charset .b  @0: ADD COLUMN (a LONGTEXT, b TINYTEXT)",
			},
		},
	})
}
