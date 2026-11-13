package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestColumnRequireDefault(t *testing.T) {
	bad := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnNoDefault, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "exempt columns", sql: "CREATE TABLE t (a INT DEFAULT 1, b INT AUTO_INCREMENT, c INT AS (a + 1), d INT PRIMARY KEY, e BLOB, f TEXT, g JSON, j LONGTEXT);"},
		{
			name: "table primary key matched without case",
			sql:  "CREATE TABLE t (ID INT, b INT, PRIMARY KEY (id));",
			want: []placed{bad("t", "b", "b INT")},
		},
		{
			name: "alter anchors the statement",
			sql:  "ALTER TABLE t ADD a INT, ADD b INT DEFAULT 1, CHANGE c d VARCHAR(5), MODIFY e INT AUTO_INCREMENT, ADD PRIMARY KEY (a);",
			want: []placed{
				bad("t", "a", "ALTER TABLE t ADD a INT, ADD b INT DEFAULT 1, CHANGE c d VARCHAR(5), MODIFY e INT AUTO_INCREMENT, ADD PRIMARY KEY (a)"),
				bad("t", "d", "ALTER TABLE t ADD a INT, ADD b INT DEFAULT 1, CHANGE c d VARCHAR(5), MODIFY e INT AUTO_INCREMENT, ADD PRIMARY KEY (a)"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnRequireDefault, tt.want)
		})
	}
}
