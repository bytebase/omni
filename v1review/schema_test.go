package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func index(name string, unique bool, columns ...string) *metadata.IndexMetadata {
	return &metadata.IndexMetadata{Name: name, Unique: unique, Expressions: columns}
}

func testSchema(objectCaseSensitive bool) Schema {
	return Schema{
		ObjectCaseSensitive: objectCaseSensitive,
		Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadata.TableMetadata{
					{Name: "Book", Indexes: []*metadata.IndexMetadata{index("idx_a", false, "a"), index("IDX_A", true, "b")}},
					{Name: "author", Indexes: []*metadata.IndexMetadata{index("idx_shared", false, "x")},
						Partitions: []*metadata.TablePartitionMetadata{{Name: "author_p0"}}},
					{Name: "review", Indexes: []*metadata.IndexMetadata{index("idx_shared", true, "y")}},
					// A later table of the same name replaces the earlier one.
					{Name: "dup", Indexes: []*metadata.IndexMetadata{index("idx_old", false, "old")}},
					{Name: "dup", Indexes: []*metadata.IndexMetadata{index("idx_new", false, "new")}},
				},
				MaterializedViews: []*metadata.MaterializedViewMetadata{
					{Name: "mv", Indexes: []*metadata.IndexMetadata{index("MV_idx", false, "m")}},
				},
			},
			{Name: "other", Tables: []*metadata.TableMetadata{{Name: "Book", Indexes: []*metadata.IndexMetadata{index("idx_o", false, "o")}}}},
		}},
	}
}

func TestSchemaTableIndex(t *testing.T) {
	for _, tt := range []struct {
		name                 string
		objectCS, indexCS    bool
		schema, table, index string
		want                 string // first expression of the found index, "" when none
	}{
		{name: "exact", objectCS: true, indexCS: true, schema: "public", table: "Book", index: "idx_a", want: "a"},
		{name: "index case sensitive keeps names apart", objectCS: true, indexCS: true, schema: "public", table: "Book", index: "IDX_A", want: "b"},
		{name: "index case insensitive takes the last", objectCS: true, indexCS: false, schema: "public", table: "Book", index: "idx_a", want: "b"},
		{name: "table case sensitive misses", objectCS: true, indexCS: true, schema: "public", table: "book", index: "idx_a"},
		{name: "table case insensitive finds", objectCS: false, indexCS: true, schema: "PUBLIC", table: "book", index: "idx_a", want: "a"},
		{name: "schema case sensitive misses", objectCS: true, indexCS: true, schema: "Public", table: "Book", index: "idx_a"},
		{name: "other schema", objectCS: true, indexCS: true, schema: "other", table: "Book", index: "idx_o", want: "o"},
		{name: "missing schema", objectCS: true, indexCS: true, schema: "", table: "Book", index: "idx_a"},
		{name: "missing index", objectCS: true, indexCS: true, schema: "public", table: "Book", index: "nope"},
		{name: "partition name finds no index", objectCS: true, indexCS: true, schema: "public", table: "author_p0", index: "idx_shared"},
		{name: "later table of a name replaces the earlier", objectCS: true, indexCS: true, schema: "public", table: "dup", index: "idx_new", want: "new"},
		{name: "replaced table's index is gone", objectCS: true, indexCS: true, schema: "public", table: "dup", index: "idx_old"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := testSchema(tt.objectCS).TableIndex(tt.schema, tt.table, tt.index, tt.indexCS)
			if tt.want == "" {
				if got != nil {
					t.Fatalf("found %v, want none", got)
				}
				return
			}
			if got == nil || got.Expressions[0] != tt.want {
				t.Fatalf("got %v, want the index on %q", got, tt.want)
			}
		})
	}
}

func TestSchemaSchemaIndex(t *testing.T) {
	for _, tt := range []struct {
		name              string
		objectCS, indexCS bool
		schema, index     string
		wantTable, want   string
	}{
		{name: "in a table", objectCS: true, indexCS: true, schema: "public", index: "idx_a", wantTable: "Book", want: "a"},
		{name: "first table in snapshot order wins", objectCS: true, indexCS: true, schema: "public", index: "idx_shared", wantTable: "author", want: "x"},
		{name: "replaced table is not searched", objectCS: true, indexCS: true, schema: "public", index: "idx_old"},
		{name: "replacing table is searched", objectCS: true, indexCS: true, schema: "public", index: "idx_new", wantTable: "dup", want: "new"},
		{name: "materialized view index", objectCS: true, indexCS: true, schema: "public", index: "MV_idx", wantTable: "mv", want: "m"},
		{name: "materialized view index compares as an object name", objectCS: false, indexCS: true, schema: "public", index: "mv_IDX", wantTable: "mv", want: "m"},
		{name: "materialized view index, object names case sensitive", objectCS: true, indexCS: false, schema: "public", index: "mv_idx"},
		{name: "missing schema", objectCS: true, indexCS: true, schema: "nope", index: "idx_a"},
		{name: "missing index", objectCS: true, indexCS: true, schema: "public", index: "nope"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			table, got := testSchema(tt.objectCS).SchemaIndex(tt.schema, tt.index, tt.indexCS)
			if tt.want == "" {
				if got != nil || table != "" {
					t.Fatalf("found %v on %q, want none", got, table)
				}
				return
			}
			if got == nil || got.Expressions[0] != tt.want || table != tt.wantTable {
				t.Fatalf("got %v on %q, want the index on %q of %q", got, table, tt.want, tt.wantTable)
			}
		})
	}
}

func TestSchemaEmpty(t *testing.T) {
	var s Schema
	if got := s.TableIndex("public", "t", "i", true); got != nil {
		t.Errorf("TableIndex on an empty schema = %v", got)
	}
	if table, got := s.SchemaIndex("public", "i", true); got != nil || table != "" {
		t.Errorf("SchemaIndex on an empty schema = %v, %q", got, table)
	}
}
