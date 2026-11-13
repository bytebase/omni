package v1review

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnNoNull(t *testing.T) {
	final := bookSchema([][2]string{{"id", "int"}, {"a", "?int"}, {"B", "?int"}, {"c", "?int"}, {"b.c", "?int"}})
	runColumnCases(t, func(stmts []ast.Node) []v1review.Finding { return ColumnNoNull(stmts, final) }, []columnCase{
		{name: "not null or unknown", sql: "CREATE TABLE book (id INT, x INT); CREATE TABLE shelf (a INT);"},
		{
			name: "create anchors columns, ties by name",
			sql:  "CREATE TABLE db.book (c INT, a INT, id INT, b INT);",
			want: []string{
				"column-nullable book.c  @0: c INT",
				"column-nullable book.a  @0: a INT",
				"column-nullable book.b  @0: b INT",
			},
		},
		{
			name: "alter anchors the action, add and change only",
			sql:  "ALTER TABLE book ADD c INT, ADD (b INT, x INT), CHANGE x a INT, MODIFY id INT;",
			want: []string{
				"column-nullable book.c  @0: ADD c INT",
				"column-nullable book.b  @0: ADD (b INT, x INT)",
				"column-nullable book.a  @0: CHANGE x a INT",
			},
		},
		{
			name: "first definition wins",
			sql:  "ALTER TABLE book ADD a INT;\nCREATE TABLE book (a INT);\nALTER TABLE BOOK ADD a INT;",
			want: []string{
				"column-nullable book.a  @0: ADD a INT",
				"column-nullable BOOK.a  @2: ADD a INT",
			},
		},
		{name: "table and column joined with a dot", sql: "CREATE TABLE `book.b` (c INT); CREATE TABLE book (`b.c` INT);"},
		{name: "a partition finds its table", sql: "ALTER TABLE book_p0 ADD a INT;", want: []string{"column-nullable book_p0.a  @0: ADD a INT"}},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN CREATE TABLE book (a INT); END;"},
	})
}
