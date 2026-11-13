package v1review

import (
	"testing"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnAutoIncrementInitialValue(t *testing.T) {
	rule := func(value int) func([]ast.Node) []v1review.Finding {
		return func(stmts []ast.Node) []v1review.Finding { return ColumnAutoIncrementInitialValue(stmts, value) }
	}
	runColumnCases(t, rule(20), []columnCase{
		{name: "matching value", sql: "CREATE TABLE t(a INT) AUTO_INCREMENT = 20;"},
		{name: "no option", sql: "CREATE TABLE t(a INT);"},
		{
			name: "create anchors the statement, qualifier dropped",
			sql:  "CREATE TABLE db.t(a INT) auto_increment = 5;",
			want: []string{"column-auto-increment-initial-value t. 5 @0: CREATE TABLE db.t(a INT) auto_increment = 5"},
		},
		{
			name: "alter anchors the statement",
			sql:  "ALTER TABLE t ADD COLUMN b INT, AUTO_INCREMENT = 7;",
			want: []string{"column-auto-increment-initial-value t. 7 @0: ALTER TABLE t ADD COLUMN b INT, AUTO_INCREMENT = 7"},
		},
		{name: "other options not judged", sql: "ALTER TABLE t ENGINE = InnoDB;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE t AUTO_INCREMENT = 7; END;"},
	})
	runColumnCases(t, rule(-1), []columnCase{
		{
			name: "a negative value is compared as uint64",
			sql:  "CREATE TABLE t(a INT) AUTO_INCREMENT = 18446744073709551614;",
			want: []string{"column-auto-increment-initial-value t. 18446744073709551614 @0: CREATE TABLE t(a INT) AUTO_INCREMENT = 18446744073709551614"},
		},
		{
			name: "the uint64 wrap of a negative value matches",
			sql:  "CREATE TABLE t(a INT) AUTO_INCREMENT = 18446744073709551615;",
		},
	})
}
