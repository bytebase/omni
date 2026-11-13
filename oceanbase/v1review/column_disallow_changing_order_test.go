package v1review

import "testing"

func TestColumnDisallowChangingOrder(t *testing.T) {
	runColumnCases(t, ColumnDisallowChangingOrder, []columnCase{
		{name: "no position", sql: "ALTER TABLE t MODIFY a INT, CHANGE b c INT;"},
		{name: "add with a position is not judged", sql: "ALTER TABLE t ADD COLUMN a INT FIRST, ADD b INT AFTER a;"},
		{
			name: "first and after, anchored at the statement",
			sql:  "ALTER TABLE t MODIFY a INT FIRST, CHANGE b c INT AFTER a;",
			want: []string{
				"column-order-changed .a  @0: ALTER TABLE t MODIFY a INT FIRST, CHANGE b c INT AFTER a",
				"column-order-changed .b  @0: ALTER TABLE t MODIFY a INT FIRST, CHANGE b c INT AFTER a",
			},
		},
	})
}
