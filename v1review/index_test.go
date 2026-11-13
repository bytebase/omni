package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func indexTableSchema(objectCaseSensitive bool) Schema {
	column := func(name, columnType string) *metadata.ColumnMetadata {
		return &metadata.ColumnMetadata{Name: name, Type: columnType}
	}
	return Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadata.TableMetadata{
					{Name: "Book", Indexes: []*metadata.IndexMetadata{{Name: "book"}}, Columns: []*metadata.ColumnMetadata{column("id", "int"), column("ID", "bigint"), column("Title", "text")},
						Partitions: []*metadata.TablePartitionMetadata{{Name: "p0", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "sp0"}}}}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{{Name: "first"}}, Partitions: []*metadata.TablePartitionMetadata{{Name: "pd"}}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{{Name: "second"}}},
					{Name: "author", Indexes: []*metadata.IndexMetadata{{Name: "author"}}, Partitions: []*metadata.TablePartitionMetadata{{Name: "Book"}, {Name: "pd"}}},
				},
			},
		}},
	}
}

func TestSchemaTableIndexes(t *testing.T) {
	for _, tt := range []struct {
		name          string
		caseSensitive bool
		table         string
		want          string // the found table's one index name, "" when none
	}{
		{name: "exact", caseSensitive: true, table: "Book", want: "book"},
		{name: "case sensitive misses", caseSensitive: true, table: "book"},
		{name: "case insensitive finds", table: "BOOK", want: "book"},
		{name: "a table wins over a partition of its name", caseSensitive: true, table: "Book", want: "book"},
		{name: "partition", caseSensitive: true, table: "p0", want: "book"},
		{name: "subpartition", caseSensitive: true, table: "sp0", want: "book"},
		{name: "partition case insensitive", table: "SP0", want: "book"},
		{name: "partition case sensitive misses", caseSensitive: true, table: "SP0"},
		{name: "later table of a name replaces the earlier", caseSensitive: true, table: "dup", want: "second"},
		{name: "last table holding a partition name wins", caseSensitive: true, table: "pd", want: "author"},
		{name: "missing", caseSensitive: true, table: "nope"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			indexes, found := indexTableSchema(tt.caseSensitive).TableIndexes("public", tt.table)
			got := ""
			if len(indexes) > 0 {
				got = indexes[0].GetName()
			}
			if got != tt.want || found != (tt.want != "") {
				t.Errorf("TableIndexes(%q) = %q, %v, want %q", tt.table, got, found, tt.want)
			}
		})
	}
	if got, found := (Schema{}).TableIndexes("public", "Book"); got != nil || found {
		t.Errorf("an empty schema found %v", got)
	}
}

func TestSchemaTableColumn(t *testing.T) {
	for _, tt := range []struct {
		name                string
		columnCaseSensitive bool
		table, column, want string
	}{
		{name: "exact", columnCaseSensitive: true, table: "Book", column: "id", want: "int"},
		{name: "case sensitive keeps names apart", columnCaseSensitive: true, table: "Book", column: "ID", want: "bigint"},
		{name: "case insensitive takes the last", table: "Book", column: "id", want: "bigint"},
		{name: "through a partition", table: "p0", column: "title", want: "text"},
		{name: "missing column", table: "Book", column: "nope"},
		{name: "missing table", table: "nope", column: "id"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := indexTableSchema(true).TableColumn("public", tt.table, tt.column, tt.columnCaseSensitive)
			if got.GetType() != tt.want {
				t.Errorf("TableColumn(%q, %q) = %q, want %q", tt.table, tt.column, got.GetType(), tt.want)
			}
		})
	}
}
