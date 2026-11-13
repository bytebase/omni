package v1review

import "testing"

func TestColumnDisallowChange(t *testing.T) {
	runColumnCases(t, ColumnDisallowChange, []columnCase{
		{name: "modify and rename are not judged", sql: "ALTER TABLE t MODIFY a INT, RENAME COLUMN b TO c;"},
		{
			name: "one finding per action, anchored at the statement",
			sql:  "ALTER TABLE t CHANGE a b INT, CHANGE COLUMN c d INT;",
			want: []string{
				"column-changed .a  @0: ALTER TABLE t CHANGE a b INT, CHANGE COLUMN c d INT",
				"column-changed .c  @0: ALTER TABLE t CHANGE a b INT, CHANGE COLUMN c d INT",
			},
		},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t CHANGE a b INT; END;"},
	})
}
