package catalog

import (
	"reflect"
	"testing"
)

// TestAddTableDuplicateColumnReplaced pins addColumn's replace-in-place
// behavior: a duplicate (already-normalized) column name must update the
// existing column rather than append a phantom second entry, so the ordinal
// column list never contains duplicates.
func TestAddTableDuplicateColumnReplaced(t *testing.T) {
	s := newSchema("s")
	tbl := s.AddTable("t",
		NewColumn("a", "bigint", false),
		NewColumn("b", "varchar", true),
		NewColumn("a", "double", true), // duplicate of "a"
	)

	if got, want := tbl.ColumnNames(), []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ColumnNames() = %v, want %v (no duplicate, position preserved)", got, want)
	}
	if got := len(tbl.Columns()); got != 2 {
		t.Errorf("len(Columns()) = %d, want 2", got)
	}
	// Last definition wins for the column's metadata.
	if col := tbl.GetColumn("a"); col == nil || col.Type != "double" || !col.Nullable {
		t.Errorf("GetColumn(\"a\") = %+v, want last definition {double, nullable}", col)
	}
}

// TestAddViewNilColumnsFiltered ensures AddView drops nil column arguments (as
// AddTable does via addColumn), so ColumnNames never dereferences a nil column.
func TestAddViewNilColumnsFiltered(t *testing.T) {
	s := newSchema("s")
	v := s.AddView("v", nil, NewColumn("x", "integer", false), nil)

	if got, want := v.ColumnNames(), []string{"x"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ColumnNames() = %v, want %v (nils filtered)", got, want)
	}
	// Must not panic on an all-nil view either.
	empty := s.AddView("empty", nil, nil)
	if got := empty.ColumnNames(); len(got) != 0 {
		t.Errorf("ColumnNames() = %v, want empty", got)
	}
}

// TestAddViewReplaces pins last-definition-wins for AddView, mirroring the
// documented AddTable replacement semantics.
func TestAddViewReplaces(t *testing.T) {
	s := newSchema("s")
	s.AddView("v", NewColumn("old", "integer", false))
	v := s.AddView("v", NewColumn("new1", "varchar", true), NewColumn("new2", "bigint", false))

	if got, want := v.ColumnNames(), []string{"new1", "new2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ColumnNames() = %v, want %v (last definition wins)", got, want)
	}
	if got := s.GetView("v"); got == nil || !reflect.DeepEqual(got.ColumnNames(), []string{"new1", "new2"}) {
		t.Errorf("GetView(\"v\") did not reflect the replacement: %+v", got)
	}
}

// TestResolveOnePartMissingCatalog covers the 1-part resolution branch when the
// current schema is set but the current catalog is not — it must fail closed
// (Found == false) rather than resolve against an empty catalog key.
func TestResolveOnePartMissingCatalog(t *testing.T) {
	c := New()
	c.EnsureCatalog("tpch").EnsureSchema("sf1").AddTable("nation", NewColumn("n", "bigint", false))
	c.SetCurrentSchema("sf1") // current catalog deliberately left unset

	if r := c.ResolveRelation([]string{"nation"}); r.Found {
		t.Errorf("ResolveRelation([nation]) with no current catalog = %+v, want Found=false", r)
	}
}
