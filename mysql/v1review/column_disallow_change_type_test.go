package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// bookSchema returns a schema holding table book, with partition book_p0,
// whose columns are given as name and type, a type starting with "?" being
// nullable, and whose indexes list the given expressions.
func bookSchema(columns [][2]string, indexes ...[]string) v1review.Schema {
	table := &metadata.TableMetadata{Name: "book", Partitions: []*metadata.TablePartitionMetadata{{Name: "book_p0"}}}
	for _, c := range columns {
		tp, nullable := c[1], false
		if len(tp) > 0 && tp[0] == '?' {
			tp, nullable = tp[1:], true
		}
		table.Columns = append(table.Columns, &metadata.ColumnMetadata{Name: c[0], Type: tp, Nullable: nullable})
	}
	for i, expressions := range indexes {
		table.Indexes = append(table.Indexes, &metadata.IndexMetadata{Name: string(rune('a' + i)), Expressions: expressions})
	}
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{table}}}}}
}

func TestColumnDisallowChangeType(t *testing.T) {
	schema := bookSchema([][2]string{{"id", "int(11)"}, {"Name", "varchar(20)"}, {"flag", "tinyint(1)"}, {"amount", "decimal(10,2) unsigned"}, {"big", "BIGINT(20) UNSIGNED"}, {"n", "int"}})
	runColumnCases(t, func(stmts []ast.Node) []v1review.Finding { return ColumnDisallowChangeType(stmts, schema) }, []columnCase{
		{name: "same types", sql: "ALTER TABLE book MODIFY id INT, MODIFY name varchar(20), CHANGE flag f TINYINT(1), MODIFY amount DECIMAL(10,2) UNSIGNED ZEROFILL, MODIFY big BIGINT UNSIGNED, MODIFY n INT(11);"},
		{name: "unknown column or table", sql: "ALTER TABLE book MODIFY other BIGINT; ALTER TABLE shelf MODIFY id BIGINT;"},
		{
			name: "each action anchors itself, case-insensitive lookup",
			sql:  "ALTER TABLE book MODIFY id BIGINT,\n  CHANGE NAME n VARCHAR(30);\nALTER TABLE db.BOOK CHANGE id id INT UNSIGNED;",
			want: []string{
				"column-type-changed book.id  @0: MODIFY id BIGINT",
				"column-type-changed book.NAME  @0: CHANGE NAME n VARCHAR(30)",
				"column-type-changed BOOK.id  @0: CHANGE id id INT UNSIGNED",
			},
		},
		{
			name: "bool is not normalized, INTEGER parses as INT",
			sql:  "ALTER TABLE book MODIFY n INTEGER, MODIFY flag BOOL;",
			want: []string{"column-type-changed book.flag  @0: MODIFY flag BOOL"},
		},
		{name: "a partition finds its table", sql: "ALTER TABLE book_p0 MODIFY id BIGINT;", want: []string{"column-type-changed book_p0.id  @0: MODIFY id BIGINT"}},
		{name: "other actions are not judged", sql: "ALTER TABLE book ADD id BIGINT, RENAME COLUMN id TO x, ALTER COLUMN id SET DEFAULT 1;"},
		{name: "a routine body is not judged", sql: "CREATE PROCEDURE p() BEGIN ALTER TABLE book MODIFY id BIGINT; END;"},
	})
	runColumnCases(t, func(stmts []ast.Node) []v1review.Finding { return ColumnDisallowChangeType(stmts, v1review.Schema{}) }, []columnCase{
		{name: "no schema", sql: "ALTER TABLE book MODIFY id BIGINT;"},
	})
}
