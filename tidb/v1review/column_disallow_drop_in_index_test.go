package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnDisallowDropInIndex(t *testing.T) {
	schema := columnSchema([][2]string{{"id", "int"}}, []string{"id"}, []string{"a", "LOWER(b)"})
	dropped := func(table, name string, statement int, text string) placed {
		return placed{Code: v1review.ColumnInIndexDropped, Name: name, Table: table, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "not indexed", sql: "ALTER TABLE book DROP COLUMN x; ALTER TABLE shelf DROP COLUMN id;"},
		{
			name: "indexes in the snapshot, every action",
			sql:  "ALTER TABLE book DROP COLUMN id, DROP COLUMN a, DROP COLUMN b;\nALTER TABLE db.BOOK DROP COLUMN `LOWER(b)`;",
			want: []placed{
				dropped("book", "id", 0, "ALTER TABLE book DROP COLUMN id, DROP COLUMN a, DROP COLUMN b"),
				dropped("book", "a", 0, "ALTER TABLE book DROP COLUMN id, DROP COLUMN a, DROP COLUMN b"),
				dropped("BOOK", "LOWER(b)", 1, "ALTER TABLE db.BOOK DROP COLUMN `LOWER(b)`"),
			},
		},
		{name: "column names compared bytewise", sql: "ALTER TABLE book DROP COLUMN ID;"},
		{
			name: "plain indexes of a create table",
			sql:  "CREATE TABLE t (a INT, b INT, c INT, d INT, INDEX (a), KEY k (b, (c + 1)), UNIQUE (c), PRIMARY KEY (d));\nDROP TABLE t;\nALTER TABLE t DROP COLUMN a, DROP COLUMN b, DROP COLUMN c, DROP COLUMN d;\nALTER TABLE T DROP COLUMN a;",
			want: []placed{
				dropped("t", "a", 2, "ALTER TABLE t DROP COLUMN a, DROP COLUMN b, DROP COLUMN c, DROP COLUMN d"),
				dropped("t", "b", 2, "ALTER TABLE t DROP COLUMN a, DROP COLUMN b, DROP COLUMN c, DROP COLUMN d"),
			},
		},
		{name: "a partition has no indexes", sql: "ALTER TABLE book_p0 DROP COLUMN id;"},
		{name: "routine bodies are not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE book DROP COLUMN id; END;"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnDisallowDropInIndex(stmts, schema)
			}, tt.want)
		})
	}
	t.Run("object case sensitive", func(t *testing.T) {
		schema := schema
		schema.ObjectCaseSensitive = true
		checkPlaced(t, "ALTER TABLE BOOK DROP COLUMN id;", func(stmts []ast.Node) []v1review.Finding {
			return ColumnDisallowDropInIndex(stmts, schema)
		}, nil)
	})
}
