package v1review

import "testing"

func TestColumnDisallowSetCharset(t *testing.T) {
	runColumnCases(t, ColumnDisallowSetCharset, []columnCase{
		{name: "no charset and binary", sql: "CREATE TABLE t(a VARCHAR(5), b VARCHAR(5) CHARACTER SET binary);"},
		{
			name: "create, anchored at the statement",
			sql:  "CREATE TABLE t(a VARCHAR(5) CHARSET utf8mb4, b TEXT CHARACTER SET latin1);",
			want: []string{
				"column-charset-set .a utf8mb4 @0: CREATE TABLE t(a VARCHAR(5) CHARSET utf8mb4, b TEXT CHARACTER SET latin1)",
				"column-charset-set .b latin1 @0: CREATE TABLE t(a VARCHAR(5) CHARSET utf8mb4, b TEXT CHARACTER SET latin1)",
			},
		},
		{
			name: "add, change and modify",
			sql:  "ALTER TABLE t ADD COLUMN (a CHAR(1) CHARSET ascii), CHANGE b c CHAR(1) CHARSET ascii, MODIFY d CHAR(1) CHARSET ascii, ALTER COLUMN e SET DEFAULT 1;",
			want: []string{
				"column-charset-set .a ascii @0: ALTER TABLE t ADD COLUMN (a CHAR(1) CHARSET ascii), CHANGE b c CHAR(1) CHARSET ascii, MODIFY d CHAR(1) CHARSET ascii, ALTER COLUMN e SET DEFAULT 1",
				"column-charset-set .c ascii @0: ALTER TABLE t ADD COLUMN (a CHAR(1) CHARSET ascii), CHANGE b c CHAR(1) CHARSET ascii, MODIFY d CHAR(1) CHARSET ascii, ALTER COLUMN e SET DEFAULT 1",
				"column-charset-set .d ascii @0: ALTER TABLE t ADD COLUMN (a CHAR(1) CHARSET ascii), CHANGE b c CHAR(1) CHARSET ascii, MODIFY d CHAR(1) CHARSET ascii, ALTER COLUMN e SET DEFAULT 1",
			},
		},
		{name: "table charset is not judged", sql: "ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4;"},
	})
}
