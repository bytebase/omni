package v1review

import "testing"

func TestColumnDisallowDrop(t *testing.T) {
	runColumnCases(t, ColumnDisallowDrop, []columnCase{
		{name: "other actions", sql: "ALTER TABLE t ADD COLUMN a INT, MODIFY b INT, DROP INDEX idx;"},
		{
			name: "each drop anchors its action, qualifier dropped",
			sql:  "ALTER TABLE db.t\n  DROP COLUMN a,\n  DROP b;",
			want: []string{
				"column-dropped t.a  @0: DROP COLUMN a",
				"column-dropped t.b  @0: DROP b",
			},
		},
		{name: "drop table is not judged", sql: "DROP TABLE t;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t DROP COLUMN a; END;"},
	})
}
