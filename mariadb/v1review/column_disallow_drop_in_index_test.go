package v1review

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowDropInIndex(t *testing.T) {
	schema := bookSchema([][2]string{{"id", "int"}, {"a", "int"}, {"b", "int"}}, []string{"a"}, []string{"id", "B"})
	runColumnCases(t, func(stmts []ast.Node) []v1review.Finding { return ColumnDisallowDropInIndex(stmts, schema) }, []columnCase{
		{name: "not indexed", sql: "ALTER TABLE book DROP COLUMN b; ALTER TABLE shelf DROP COLUMN a;"},
		{
			name: "schema indexes, each drop anchored, statement recorded",
			sql:  "ALTER TABLE book ADD x INT;\nALTER TABLE db.book\n  DROP COLUMN a,\n  DROP B, DROP COLUMN id;",
			want: []string{
				"column-in-index-dropped book.a  @1: DROP COLUMN a",
				"column-in-index-dropped book.B  @1: DROP B",
				"column-in-index-dropped book.id  @1: DROP COLUMN id",
			},
		},
		{
			name: "plain indexes of a created table, compared as written",
			sql:  "CREATE TABLE u(a INT, b INT, c INT, KEY (a), UNIQUE KEY (b), PRIMARY KEY (c));\nALTER TABLE u DROP COLUMN a, DROP COLUMN b, DROP COLUMN c, DROP COLUMN A;\nALTER TABLE U DROP COLUMN a;",
			want: []string{"column-in-index-dropped u.a  @1: DROP COLUMN a"},
		},
		{
			name: "drop table does not forget",
			sql:  "CREATE TABLE u(a INT, INDEX (a)); DROP TABLE u; ALTER TABLE u DROP COLUMN a;",
			want: []string{"column-in-index-dropped u.a  @2: DROP COLUMN a"},
		},
		{name: "a partition has no indexes", sql: "ALTER TABLE book_p0 DROP COLUMN a;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE book DROP COLUMN a; END;"},
	})
}
