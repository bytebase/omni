package v1review

import "testing"

func TestColumnAutoIncrementMustInteger(t *testing.T) {
	runColumnCases(t, ColumnAutoIncrementMustInteger, []columnCase{
		{name: "integer types", sql: "CREATE TABLE t(a INT AUTO_INCREMENT, b bigint AUTO_INCREMENT, c INTEGER AUTO_INCREMENT, d TINYINT UNSIGNED AUTO_INCREMENT);"},
		{name: "not auto-increment", sql: "CREATE TABLE t(a VARCHAR(5));"},
		{
			name: "create anchors the column",
			sql:  "CREATE TABLE db.t(\n  a INT,\n  b DECIMAL(10) AUTO_INCREMENT);",
			want: []string{"column-auto-increment-not-integer t.b  @0: b DECIMAL(10) AUTO_INCREMENT"},
		},
		{
			name: "add, modify and change",
			sql:  "ALTER TABLE t ADD COLUMN (a FLOAT AUTO_INCREMENT), MODIFY b DOUBLE AUTO_INCREMENT, CHANGE c d FLOAT AUTO_INCREMENT;",
			want: []string{
				"column-auto-increment-not-integer t.a  @0: a FLOAT AUTO_INCREMENT",
				"column-auto-increment-not-integer t.b  @0: b DOUBLE AUTO_INCREMENT",
				"column-auto-increment-not-integer t.d  @0: d FLOAT AUTO_INCREMENT",
			},
		},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t(a FLOAT AUTO_INCREMENT); END;"},
	})
}
