package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestColumnSetDefaultForNotNull(t *testing.T) {
	bad := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnNotNullNoDefault, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "nullable or defaulted passes", sql: "CREATE TABLE t (a INT, b INT NOT NULL DEFAULT 1, c TEXT NOT NULL, d INT NOT NULL AUTO_INCREMENT, e INT PRIMARY KEY);"},
		{
			name: "table primary key matched exactly",
			sql:  "CREATE TABLE t (ID INT NOT NULL, id INT NOT NULL, PRIMARY KEY (id));",
			want: []placed{bad("t", "ID", "ID INT NOT NULL")},
		},
		{
			name: "alter anchors the statement",
			sql:  "ALTER TABLE t ADD a INT NOT NULL, MODIFY b INT NULL, CHANGE c d INT NOT NULL;",
			want: []placed{
				bad("t", "a", "ALTER TABLE t ADD a INT NOT NULL, MODIFY b INT NULL, CHANGE c d INT NOT NULL"),
				bad("t", "d", "ALTER TABLE t ADD a INT NOT NULL, MODIFY b INT NULL, CHANGE c d INT NOT NULL"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnSetDefaultForNotNull, tt.want)
		})
	}
}
