package v1review

import "testing"

func TestColumnSetDefaultForNotNull(t *testing.T) {
	runColumnCases(t, ColumnSetDefaultForNotNull, []columnCase{
		{name: "nullable or defaulted", sql: "CREATE TABLE t(a INT, b INT NOT NULL DEFAULT 1);"},
		{name: "exempt kinds", sql: "CREATE TABLE t(a INT NOT NULL AUTO_INCREMENT, b TEXT NOT NULL, c JSON NOT NULL, d GEOMETRY NOT NULL);"},
		{
			name: "create skips primary keys, spatial subtypes are not exempt",
			sql:  "CREATE TABLE t(id INT NOT NULL, k INT NOT NULL, p POINT NOT NULL, PRIMARY KEY (id));",
			want: []string{
				"column-not-null-no-default t.k  @0: k INT NOT NULL",
				"column-not-null-no-default t.p  @0: p POINT NOT NULL",
			},
		},
		{
			name: "alter exempts inline primary keys",
			sql:  "ALTER TABLE t ADD COLUMN id INT PRIMARY KEY, MODIFY b INT NOT NULL, CHANGE c d INT NOT NULL DEFAULT 0;",
			want: []string{"column-not-null-no-default t.b  @0: b INT NOT NULL"},
		},
	})
}
