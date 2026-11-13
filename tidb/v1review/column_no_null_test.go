package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnNoNull(t *testing.T) {
	final := columnSchema([][2]string{{"id", "int"}, {"a", "?int"}, {"B", "?int"}, {"c", "?int"}, {"b.c", "?int"}})
	nullable := func(table, name string, statement int, text string) placed {
		return placed{Code: v1review.ColumnNullable, Name: name, Table: table, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "not null or unknown", sql: "CREATE TABLE book (id INT, x INT); CREATE TABLE shelf (a INT);"},
		{
			name: "create anchors columns, ties by name",
			sql:  "CREATE TABLE db.book (c INT, a INT, id INT, b INT);",
			want: []placed{
				nullable("book", "c", 0, "c INT"),
				nullable("book", "a", 0, "a INT"),
				nullable("book", "b", 0, "b INT"),
			},
		},
		{
			name: "alter anchors the statement, add and change only",
			sql:  "ALTER TABLE book ADD c INT, ADD (b INT, x INT), CHANGE x a INT, MODIFY id INT;\nALTER TABLE book MODIFY c INT NULL;",
			want: []placed{
				nullable("book", "a", 0, "ALTER TABLE book ADD c INT, ADD (b INT, x INT), CHANGE x a INT, MODIFY id INT"),
				nullable("book", "b", 0, "ALTER TABLE book ADD c INT, ADD (b INT, x INT), CHANGE x a INT, MODIFY id INT"),
				nullable("book", "c", 0, "ALTER TABLE book ADD c INT, ADD (b INT, x INT), CHANGE x a INT, MODIFY id INT"),
			},
		},
		{
			name: "first definition wins",
			sql:  "ALTER TABLE book ADD a INT;\nCREATE TABLE book (a INT);\nALTER TABLE BOOK ADD a INT;",
			want: []placed{
				nullable("book", "a", 0, "ALTER TABLE book ADD a INT"),
				nullable("BOOK", "a", 2, "ALTER TABLE BOOK ADD a INT"),
			},
		},
		{
			name: "table and column joined with a dot",
			sql:  "CREATE TABLE `book.b` (c INT); CREATE TABLE book (`b.c` INT);",
		},
		{name: "a partition finds its table", sql: "ALTER TABLE book_p0 ADD a INT;", want: []placed{nullable("book_p0", "a", 0, "ALTER TABLE book_p0 ADD a INT")}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnNoNull(stmts, final)
			}, tt.want)
		})
	}
	t.Run("dot key", func(t *testing.T) {
		// `book.b`.c comes first and takes the key, so book.`b.c`, which
		// final holds as nullable, is never looked up.
		final := final
		final.Metadata.Schemas[0].Tables = append(final.Metadata.Schemas[0].Tables, &metadata.TableMetadata{Name: "book.b", Columns: []*metadata.ColumnMetadata{{Name: "c", Nullable: true}}})
		checkPlaced(t, "CREATE TABLE `book.b` (c INT); CREATE TABLE book (`b.c` INT);", func(stmts []ast.Node) []v1review.Finding {
			return ColumnNoNull(stmts, final)
		}, []placed{nullable("book.b", "c", 0, "c INT")})
	})
	t.Run("no snapshot", func(t *testing.T) {
		checkPlaced(t, "CREATE TABLE book (a INT);", func(stmts []ast.Node) []v1review.Finding {
			return ColumnNoNull(stmts, v1review.Schema{})
		}, nil)
	})
}
