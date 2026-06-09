package analysis

import (
	"slices"
	"testing"
)

// hasSource reports whether a result column's resolved sources contain want.
func hasSource(srcs []ColumnRef, want ColumnRef) bool {
	return slices.Contains(srcs, want)
}

func resultByName(span *QuerySpan, name string) (ColumnInfo, bool) {
	for _, r := range span.Results {
		if r.Name == name {
			return r, true
		}
	}
	return ColumnInfo{}, false
}

// TestGetQuerySpan_DerivedTableColumnLineage covers BYT-9674: a column projected
// through a derived table (subquery in FROM) and re-aliased must resolve to the
// underlying base column, not the derived relation's alias, so masking sees it.
func TestGetQuerySpan_DerivedTableColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT d.x FROM (SELECT phone AS x FROM customer) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "x")
	if !ok {
		t.Fatalf("Results = %+v, want a column named x", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("x.SourceColumns = %+v, want to contain {Column:phone} (resolved through derived table d)", r.SourceColumns)
	}
}

// TestGetQuerySpan_CTEColumnLineage covers BYT-9675: a column projected through a
// CTE and re-aliased must resolve to the underlying base column. A CTE reference
// is the same derived-relation mechanism as a subquery in FROM.
func TestGetQuerySpan_CTEColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("WITH w AS (SELECT phone AS pp FROM customer) SELECT pp FROM w")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "pp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named pp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("pp.SourceColumns = %+v, want to contain {Column:phone} (resolved through CTE w)", r.SourceColumns)
	}
}

// TestGetQuerySpan_CTEColumnAliasLineage exercises the CTE column-alias form
// `WITH w (pp) AS (...)`, where the projected column is renamed by the CTE's
// column-alias list rather than a select-item alias.
func TestGetQuerySpan_CTEColumnAliasLineage(t *testing.T) {
	span, err := GetQuerySpan("WITH w (pp) AS (SELECT phone FROM customer) SELECT pp FROM w")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "pp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named pp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("pp.SourceColumns = %+v, want to contain {Column:phone} (resolved through CTE column alias)", r.SourceColumns)
	}
}

// TestGetQuerySpan_NestedDerivedColumnLineage resolves a column through two
// nested derived tables back to the base column.
func TestGetQuerySpan_NestedDerivedColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT o.c FROM (SELECT i.b AS c FROM (SELECT phone AS b FROM customer) i) o")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "c")
	if !ok {
		t.Fatalf("Results = %+v, want a column named c", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("c.SourceColumns = %+v, want to contain {Column:phone} (resolved through two derived tables)", r.SourceColumns)
	}
}

// TestGetQuerySpan_DerivedExpressionColumnLineage resolves a derived column that
// is an expression: both base inputs must surface.
func TestGetQuerySpan_DerivedExpressionColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT d.s FROM (SELECT a + b AS s FROM t) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "s")
	if !ok {
		t.Fatalf("Results = %+v, want a column named s", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "a"}) || !hasSource(r.SourceColumns, ColumnRef{Column: "b"}) {
		t.Errorf("s.SourceColumns = %+v, want to contain {Column:a} and {Column:b}", r.SourceColumns)
	}
}

// TestGetQuerySpan_DerivedSetOpDoesNotDropLineage guards against an under-masking
// regression: a derived relation containing a set operation has only its left
// arm's lineage computed, but a bare outer ref must still retain its original
// base ref so a sensitive column in the right arm is not dropped (additive
// resolution).
func TestGetQuerySpan_DerivedSetOpDoesNotDropLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT phone FROM (SELECT CAST(NULL AS varchar) AS phone UNION ALL SELECT phone FROM customer) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want to retain {Column:phone} (set-op right arm must not be dropped)", r.SourceColumns)
	}
}

// TestGetQuerySpan_DerivedStarAliasDoesNotDropLineage guards against an
// under-masking regression: a derived column sourced from SELECT * has unknown
// (empty) lineage; renaming it via a relation column alias must not turn the
// outer ref into a known-empty source and drop the base ref.
func TestGetQuerySpan_DerivedStarAliasDoesNotDropLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT phone FROM (SELECT * FROM (SELECT phone FROM customer) c) AS d(phone)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want to retain {Column:phone} (star-sourced derived column must not drop lineage)", r.SourceColumns)
	}
}

// TestGetQuerySpan_JoinUsingDerivedDoesNotDropLineage guards against an
// under-masking regression: with JOIN ... USING (phone), bare phone is a valid
// output coalescing a base table and a derived relation; the base ref must be
// retained even though the derived relation also exposes phone.
func TestGetQuerySpan_JoinUsingDerivedDoesNotDropLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT phone FROM customer JOIN (SELECT public_id AS phone FROM public_ids) p USING (phone)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want to retain {Column:phone} (JOIN USING base side must not be dropped)", r.SourceColumns)
	}
}

// TestGetQuerySpan_DirectColumnLineageUnchanged guards that a direct base-table
// column (no derived relation in scope) is left exactly as the primary walk
// produced it — the resolver must only deepen lineage through indirection.
func TestGetQuerySpan_DirectColumnLineageUnchanged(t *testing.T) {
	span, err := GetQuerySpan("SELECT phone FROM customer")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if len(r.SourceColumns) != 1 || r.SourceColumns[0] != (ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want exactly [{Column:phone}]", r.SourceColumns)
	}
}
