package v1review

import "testing"

func TestColumnAutoIncrementMustUnsigned(t *testing.T) {
	runColumnCases(t, ColumnAutoIncrementMustUnsigned, []columnCase{
		{name: "unsigned and zerofill", sql: "CREATE TABLE t(a INT UNSIGNED AUTO_INCREMENT, b INT ZEROFILL AUTO_INCREMENT);"},
		{name: "not auto-increment", sql: "CREATE TABLE t(a INT);"},
		{
			name: "signed columns",
			sql:  "CREATE TABLE t(a INT AUTO_INCREMENT);\nALTER TABLE u MODIFY b BIGINT AUTO_INCREMENT;",
			want: []string{
				"column-auto-increment-not-unsigned t.a  @0: a INT AUTO_INCREMENT",
				"column-auto-increment-not-unsigned u.b  @0: b BIGINT AUTO_INCREMENT",
			},
		},
		{
			name: "a non-integer type is reported",
			sql:  "ALTER TABLE t ADD COLUMN a DOUBLE AUTO_INCREMENT;",
			want: []string{"column-auto-increment-not-unsigned t.a  @0: a DOUBLE AUTO_INCREMENT"},
		},
		{name: "drop is not judged", sql: "ALTER TABLE t DROP COLUMN a;"},
	})
}
