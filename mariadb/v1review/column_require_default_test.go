package v1review

import "testing"

func TestColumnRequireDefault(t *testing.T) {
	runColumnCases(t, ColumnRequireDefault, []columnCase{
		{name: "defaults and exempt types", sql: "CREATE TABLE t(a INT DEFAULT 1, b TEXT, c JSON, d POINT, e GEOMETRYCOLLECTION);"},
		{
			name: "create skips primary key columns, auto-increment is not exempt",
			sql:  "CREATE TABLE db.t(id INT PRIMARY KEY, k INT, a INT AUTO_INCREMENT, PRIMARY KEY (k));",
			want: []string{"column-no-default t.a  @0: a INT AUTO_INCREMENT"},
		},
		{
			name: "alter judges primary keys too",
			sql:  "ALTER TABLE t ADD COLUMN id INT PRIMARY KEY, MODIFY b INT, CHANGE c d INT DEFAULT 0;",
			want: []string{
				"column-no-default t.id  @0: id INT PRIMARY KEY",
				"column-no-default t.b  @0: b INT",
			},
		},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t(a INT); END;"},
	})
}
