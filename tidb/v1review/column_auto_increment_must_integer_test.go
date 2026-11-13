package v1review

import (
	"testing"

	"github.com/bytebase/omni/v1review"
)

func TestColumnAutoIncrementMustInteger(t *testing.T) {
	bad := func(table, name, text string) placed {
		return placed{Code: v1review.ColumnAutoIncrementNotInteger, Name: name, Table: table, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "integer types pass", sql: "CREATE TABLE t (a TINYINT AUTO_INCREMENT, b SMALLINT AUTO_INCREMENT, c MEDIUMINT AUTO_INCREMENT, d INT AUTO_INCREMENT, e INTEGER AUTO_INCREMENT, f BIGINT AUTO_INCREMENT, g BOOL AUTO_INCREMENT, h boolean AUTO_INCREMENT);"},
		{name: "not auto increment passes", sql: "CREATE TABLE t (a VARCHAR(10), b DOUBLE);"},
		{
			name: "create anchors the column",
			sql:  "CREATE TABLE db.t (\n  id DOUBLE AUTO_INCREMENT,\n  n DECIMAL(10,2) auto_increment\n);",
			want: []placed{bad("t", "id", "id DOUBLE AUTO_INCREMENT"), bad("t", "n", "n DECIMAL(10,2) auto_increment")},
		},
		{
			name: "alter anchors the statement",
			sql:  "ALTER TABLE t ADD COLUMN a FLOAT AUTO_INCREMENT, ADD (b VARCHAR(5) AUTO_INCREMENT, c INT AUTO_INCREMENT), CHANGE d e BIT AUTO_INCREMENT, MODIFY f YEAR AUTO_INCREMENT;",
			want: []placed{
				bad("t", "a", "ALTER TABLE t ADD COLUMN a FLOAT AUTO_INCREMENT, ADD (b VARCHAR(5) AUTO_INCREMENT, c INT AUTO_INCREMENT), CHANGE d e BIT AUTO_INCREMENT, MODIFY f YEAR AUTO_INCREMENT"),
				bad("t", "b", "ALTER TABLE t ADD COLUMN a FLOAT AUTO_INCREMENT, ADD (b VARCHAR(5) AUTO_INCREMENT, c INT AUTO_INCREMENT), CHANGE d e BIT AUTO_INCREMENT, MODIFY f YEAR AUTO_INCREMENT"),
				bad("t", "e", "ALTER TABLE t ADD COLUMN a FLOAT AUTO_INCREMENT, ADD (b VARCHAR(5) AUTO_INCREMENT, c INT AUTO_INCREMENT), CHANGE d e BIT AUTO_INCREMENT, MODIFY f YEAR AUTO_INCREMENT"),
				bad("t", "f", "ALTER TABLE t ADD COLUMN a FLOAT AUTO_INCREMENT, ADD (b VARCHAR(5) AUTO_INCREMENT, c INT AUTO_INCREMENT), CHANGE d e BIT AUTO_INCREMENT, MODIFY f YEAR AUTO_INCREMENT"),
			},
		},
		{name: "routine bodies are not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE t (a DOUBLE AUTO_INCREMENT); END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, ColumnAutoIncrementMustInteger, tt.want)
		})
	}
}
