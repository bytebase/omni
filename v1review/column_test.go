package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestSchemaTableColumn(t *testing.T) {
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

func TestSchemaTableIndexes(t *testing.T) {
	a := &metadata.IndexMetadata{Name: "a"}
	b := &metadata.IndexMetadata{Name: "b"}
	upperA := &metadata.IndexMetadata{Name: "A"}
	db := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{
			Name:       "t",
			Indexes:    []*metadata.IndexMetadata{a, b, upperA},
			Partitions: []*metadata.TablePartitionMetadata{{Name: "p"}},
		}},
	}}}
	s := Schema{Metadata: db}
	if got := s.TableIndexes("", "t", false); len(got) != 2 || got[0] != upperA || got[1] != b {
		t.Errorf("case-insensitive: got %v", got)
	}
	if got := s.TableIndexes("", "t", true); len(got) != 3 {
		t.Errorf("case-sensitive: got %v", got)
	}
	if got := s.TableIndexes("", "p", false); got != nil {
		t.Errorf("partition: got %v", got)
	}
	if got := s.TableIndexes("", "missing", false); got != nil {
		t.Errorf("missing: got %v", got)
	}
}
