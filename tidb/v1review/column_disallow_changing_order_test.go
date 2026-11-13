package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowChangingOrder(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "no position passes", sql: "ALTER TABLE t MODIFY a INT, CHANGE b c INT;"},
		{name: "add with a position is not judged", sql: "ALTER TABLE t ADD a INT FIRST, ADD b INT AFTER a;"},
		{
			name: "first and after, once per statement",
			sql:  "ALTER TABLE t MODIFY a INT FIRST, CHANGE b c INT AFTER a;\nALTER TABLE u CHANGE b c INT AFTER a;",
			want: []placed{
				{Code: v1review.ColumnOrderChanged, Text: "ALTER TABLE t MODIFY a INT FIRST, CHANGE b c INT AFTER a"},
				{Code: v1review.ColumnOrderChanged, Text: "ALTER TABLE u CHANGE b c INT AFTER a"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnDisallowChangingOrder, tt.want)
		})
	}
}
