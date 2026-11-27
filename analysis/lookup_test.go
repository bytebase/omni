package analysis

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/metadata"
)

func lookupFixture() *metadata.DatabaseSchemaMetadata {
	return &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{
		{
			Name: "Public",
			Tables: []*metadata.TableMetadata{
				{Name: "Orders", Partitions: []*metadata.TablePartitionMetadata{{Name: "orders_2025", Subpartitions: []*metadata.TablePartitionMetadata{{Name: "orders_2025_q1"}}}}},
				{Name: "orders_2025"},
			},
			Views:     []*metadata.ViewMetadata{{Name: "V"}},
			Sequences: []*metadata.SequenceMetadata{{Name: "Seq"}},
			Functions: []*metadata.FunctionMetadata{{Name: "Fn", Signature: "a"}, {Name: "fn", Signature: "b"}},
		},
		{Name: "public", Tables: []*metadata.TableMetadata{{Name: "dup"}}},
		{Name: "s1", MaterializedViews: []*metadata.MaterializedViewMetadata{{Name: "mv"}}, ExternalTables: []*metadata.ExternalTableMetadata{{Name: "ext"}}},
	}}
}

func TestLookupFoldsCaseInsensitiveNamesAndTheLaterDuplicateWins(t *testing.T) {
	l := NewLookup(lookupFixture(), false, false)
	if got := l.Schema("PUBLIC").Proto().GetName(); got != "public" {
		t.Fatalf("Schema(PUBLIC) = %q, want the later duplicate", got)
	}
	if got := l.ListSchemaNames(); !slices.Equal(got, []string{"public", "s1"}) {
		t.Fatalf("ListSchemaNames = %v", got)
	}
	if l.Schema("public").Table("Orders") != nil {
		t.Fatal("the earlier duplicate schema's table is still found")
	}
}

func TestLookupMatchesCaseSensitiveNamesExactly(t *testing.T) {
	l := NewLookup(lookupFixture(), true, true)
	s := l.Schema("Public")
	if s == nil || l.Schema("PUBLIC") != nil {
		t.Fatal("schema names must match exactly")
	}
	if s.Table("Orders") == nil || s.Table("orders") != nil || s.View("V") == nil || s.View("v") != nil {
		t.Fatal("object names must match exactly")
	}
	if s.Sequence("Seq") == nil || s.Sequence("seq") != nil {
		t.Fatal("sequence names follow detail case sensitivity")
	}
	if got := s.Function("fn").GetSignature(); got != "b" {
		t.Fatalf("Function(fn) = %q", got)
	}
	if got := s.ListTableNames(); !slices.Equal(got, []string{"Orders", "orders_2025"}) {
		t.Fatalf("ListTableNames = %v", got)
	}
}

func TestLookupDetailCaseSensitivityIsSeparate(t *testing.T) {
	s := NewLookup(lookupFixture(), true, false).Schema("Public")
	if s.Sequence("SEQ") == nil {
		t.Fatal("a case-insensitive detail matches any case")
	}
	if got := s.Function("FN").GetSignature(); got != "a" {
		t.Fatalf("Function(FN) = %q, want the first match", got)
	}
	if s.Table("orders") != nil {
		t.Fatal("object names stay case-sensitive")
	}
}

func TestLookupPartitionsResolveToTheirRootAfterRealTables(t *testing.T) {
	s := NewLookup(lookupFixture(), true, true).Schema("Public")
	if got := s.Table("orders_2025_q1").GetName(); got != "Orders" {
		t.Fatalf("Table(subpartition) = %q", got)
	}
	if got := s.Table("orders_2025").GetName(); got != "orders_2025" {
		t.Fatalf("Table(orders_2025) = %q, want the real table", got)
	}
}

func TestLookupSearchesAlongThePath(t *testing.T) {
	l := NewLookup(lookupFixture(), true, true)
	if schema, name := l.SearchRelation([]string{"missing", "s1", "Public"}, "ext"); schema != "s1" || name != "ext" {
		t.Fatalf("SearchRelation(ext) = %q, %q", schema, name)
	}
	if schema, _ := l.SearchRelation([]string{"Public"}, "Seq"); schema != "Public" {
		t.Fatal("a sequence is a relation")
	}
	if schema, name := l.SearchRelation([]string{"Public"}, "Fn"); schema != "" || name != "" {
		t.Fatal("a function is not a relation")
	}
	schemas, functions := l.SearchFunctions([]string{"s1", "Public"}, "Fn")
	if !slices.Equal(schemas, []string{"Public"}) || len(functions) != 1 || functions[0].GetSignature() != "a" {
		t.Fatalf("SearchFunctions = %v, %v", schemas, functions)
	}
	if s, f := NewLookup(nil, true, true).SearchFunctions([]string{"x"}, "f"); s != nil || f != nil {
		t.Fatal("a nil database finds nothing")
	}
}

func TestDatabaseLookupIsBuiltOnce(t *testing.T) {
	d := &Database{Metadata: lookupFixture(), ObjectCaseSensitive: true}
	if d.Lookup(true) != d.Lookup(false) {
		t.Fatal("Lookup must be built once")
	}
}
