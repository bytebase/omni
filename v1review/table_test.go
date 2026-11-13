package v1review

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func tableSchema(objectCaseSensitive bool) Schema {
	pk := func(name string, columns ...string) *metadata.IndexMetadata {
		return &metadata.IndexMetadata{Name: name, Primary: true, Expressions: columns}
	}
	return Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
			{Name: "old", Tables: []*metadata.TableMetadata{{Name: "gone"}}},
			{
				Name: "public",
				Tables: []*metadata.TableMetadata{
					{
						Name:    "Book",
						Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "title"}},
						Indexes: []*metadata.IndexMetadata{index("idx_a", false, "a"), pk("PRIMARY", "id")},
						Partitions: []*metadata.TablePartitionMetadata{
							{Name: "p0", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "sp0"}}},
						},
					},
					{Name: "author", Columns: []*metadata.ColumnMetadata{{Name: "name"}}, Partitions: []*metadata.TablePartitionMetadata{{Name: "p0"}}},
					// The later "primary" replaces "PRIMARY" when index names
					// are compared case-insensitively.
					{Name: "two", Indexes: []*metadata.IndexMetadata{pk("PRIMARY", "a"), {Name: "primary", Expressions: []string{"b"}}}},
					{Name: "dup", Partitions: []*metadata.TablePartitionMetadata{{Name: "dp"}}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{pk("PRIMARY", "new")}},
				},
			},
			// A later schema of the same name replaces the earlier one.
			{Name: "OLD", Tables: []*metadata.TableMetadata{{Name: "kept"}}},
		}},
	}
}

func TestSchemaTables(t *testing.T) {
	ci, cs := tableSchema(false), tableSchema(true)
	for _, tt := range []struct {
		name          string
		schema        Schema
		in, table     string
		indexCS       bool
		has           bool
		columns       []string
		primaryKeyCol string
	}{
		{name: "table", schema: cs, in: "public", table: "Book", has: true, columns: []string{"id", "title"}, primaryKeyCol: "id"},
		{name: "case-insensitive table", schema: ci, in: "PUBLIC", table: "book", has: true, columns: []string{"id", "title"}, primaryKeyCol: "id"},
		{name: "case-sensitive miss", schema: cs, in: "public", table: "book"},
		{name: "partition has its table's columns and no primary key", schema: cs, in: "public", table: "sp0", has: true, columns: []string{"id", "title"}},
		{name: "last table holding a partition name wins", schema: cs, in: "public", table: "p0", has: true, columns: []string{"name"}},
		{name: "replaced table still lends its partitions", schema: cs, in: "public", table: "dp", has: true},
		{name: "later table of a name replaces the earlier", schema: cs, in: "public", table: "dup", has: true, primaryKeyCol: "new"},
		{name: "index replaced by a later one of the same name", schema: cs, in: "public", table: "two", indexCS: false, has: true},
		{name: "index names kept apart when case sensitive", schema: cs, in: "public", table: "two", indexCS: true, has: true, primaryKeyCol: "a"},
		{name: "replaced schema's table is gone", schema: ci, in: "old", table: "gone"},
		{name: "missing table", schema: cs, in: "public", table: "nope"},
		{name: "missing schema", schema: cs, in: "nope", table: "Book"},
		{name: "nil metadata", schema: Schema{}, in: "", table: "Book"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.schema.HasTable(tt.in, tt.table); got != tt.has {
				t.Errorf("HasTable = %v, want %v", got, tt.has)
			}
			columns, found := tt.schema.TableColumns(tt.in, tt.table)
			var names []string
			for _, c := range columns {
				names = append(names, c.GetName())
			}
			if found != tt.has || !slices.Equal(names, tt.columns) {
				t.Errorf("TableColumns = %v, %v; want %v, %v", names, found, tt.columns, tt.has)
			}
			pk := tt.schema.TablePrimaryKey(tt.in, tt.table, tt.indexCS)
			var got string
			if pk != nil {
				got = pk.GetExpressions()[0]
			}
			if got != tt.primaryKeyCol {
				t.Errorf("TablePrimaryKey column = %q, want %q", got, tt.primaryKeyCol)
			}
		})
	}
}

func TestSchemaNames(t *testing.T) {
	if got, want := tableSchema(false).SchemaNames(), []string{"OLD", "public"}; !slices.Equal(got, want) {
		t.Errorf("case-insensitive SchemaNames = %v, want %v", got, want)
	}
	if got, want := tableSchema(true).SchemaNames(), []string{"OLD", "old", "public"}; !slices.Equal(got, want) {
		t.Errorf("case-sensitive SchemaNames = %v, want %v", got, want)
	}
	if got, want := tableSchema(true).TableNames("public"), []string{"Book", "author", "dup", "two"}; !slices.Equal(got, want) {
		t.Errorf("TableNames = %v, want %v", got, want)
	}
	if got := tableSchema(true).TableNames("nope"); got != nil {
		t.Errorf("TableNames of a missing schema = %v", got)
	}
}
