package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

func TestSchemaTable(t *testing.T) {
	orders := &metadata.TableMetadata{Name: "Orders", Partitions: []*metadata.TablePartitionMetadata{
		{Name: "p2024", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "p2024_east"}}},
	}}
	ordersAgain := &metadata.TableMetadata{Name: "orders"}
	p2024 := &metadata.TableMetadata{Name: "P2024"}
	db := &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{Name: "public", Tables: []*metadata.TableMetadata{orders, ordersAgain, p2024}},
	}}
	for _, tt := range []struct {
		name          string
		caseSensitive bool
		schema, table string
		want          *metadata.TableMetadata
	}{
		{name: "exact", caseSensitive: true, schema: "public", table: "Orders", want: orders},
		{name: "folded, last wins", schema: "PUBLIC", table: "ORDERS", want: ordersAgain},
		{name: "case sensitive misses", caseSensitive: true, schema: "public", table: "ORDERS"},
		{name: "subpartition finds its table", caseSensitive: true, schema: "public", table: "p2024_east", want: orders},
		{name: "partition found by name", caseSensitive: true, schema: "public", table: "p2024", want: orders},
		{name: "a table outranks a partition", schema: "public", table: "p2024", want: p2024},
		{name: "missing schema", schema: "other", table: "Orders"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Schema{Metadata: db, ObjectCaseSensitive: tt.caseSensitive}.Table(tt.schema, tt.table)
			if got != tt.want {
				t.Errorf("Table(%q, %q) = %v, want %v", tt.schema, tt.table, got, tt.want)
			}
		})
	}
	if got := (Schema{}).Table("public", "Orders"); got != nil {
		t.Errorf("nil snapshot: got %v", got)
	}
}
