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

func TestSchemaTableColumnAsV1Read(t *testing.T) {
	first := &metadata.ColumnMetadata{Name: "id", Type: "int"}
	last := &metadata.ColumnMetadata{Name: "ID", Type: "bigint"}
	other := &metadata.ColumnMetadata{Name: "x"}
	db := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{
			{Name: "t", Columns: []*metadata.ColumnMetadata{first, last}},
			{Name: "u", Columns: []*metadata.ColumnMetadata{other}, Partitions: []*metadata.TablePartitionMetadata{
				{Name: "p0", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "sp"}}},
				{Name: "t"},
			}},
		},
	}}}
	s := Schema{Metadata: db}
	for _, tt := range []struct {
		name                string
		table, column       string
		columnCaseSensitive bool
		objectCaseSensitive bool
		want                *metadata.ColumnMetadata
	}{
		{name: "case-insensitive column, last wins", table: "T", column: "Id", want: last},
		{name: "case-sensitive column", table: "t", column: "id", columnCaseSensitive: true, want: first},
		{name: "case-sensitive column missing", table: "t", column: "Id", columnCaseSensitive: true},
		{name: "partition finds its table's columns", table: "p0", column: "x", want: other},
		{name: "subpartition finds its table's columns", table: "sp", column: "x", want: other},
		{name: "a table wins over a partition of the same name", table: "t", column: "x"},
		{name: "case-sensitive table", table: "T", column: "id", objectCaseSensitive: true},
		{name: "missing table", table: "v", column: "id"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s.ObjectCaseSensitive = tt.objectCaseSensitive
			if got := s.TableColumn("", tt.table, tt.column, tt.columnCaseSensitive); got != tt.want {
				t.Errorf("TableColumn = %v, want %v", got, tt.want)
			}
		})
	}
	if got := (Schema{}).TableColumn("", "t", "id", false); got != nil {
		t.Errorf("nil metadata: got %v", got)
	}
}
