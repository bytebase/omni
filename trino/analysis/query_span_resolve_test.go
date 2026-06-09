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

// TestGetQuerySpan_UnnestColumnLineage covers BYT-9680: a column produced by
// UNNEST of a sensitive array column must resolve to that array column so the
// unnested values are masked.
func TestGetQuerySpan_UnnestColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT t.p FROM customer CROSS JOIN UNNEST(phones) AS t(p)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phones"}) {
		t.Errorf("p.SourceColumns = %+v, want to contain {Column:phones} (resolved through UNNEST)", r.SourceColumns)
	}
}

// TestGetQuerySpan_UnnestBareColumnLineage resolves an UNNEST output referenced
// by its bare column alias (no relation qualifier).
func TestGetQuerySpan_UnnestBareColumnLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT p FROM customer CROSS JOIN UNNEST(phones) AS t(p)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phones"}) {
		t.Errorf("p.SourceColumns = %+v, want to contain {Column:phones}", r.SourceColumns)
	}
}

// TestGetQuerySpan_UnnestWithOrdinality resolves the element column's lineage
// while leaving the trailing WITH ORDINALITY ordinal column without lineage.
func TestGetQuerySpan_UnnestWithOrdinality(t *testing.T) {
	span, err := GetQuerySpan("SELECT t.p, t.ord FROM customer CROSS JOIN UNNEST(phones) WITH ORDINALITY AS t(p, ord)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	p, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(p.SourceColumns, ColumnRef{Column: "phones"}) {
		t.Errorf("p.SourceColumns = %+v, want to contain {Column:phones}", p.SourceColumns)
	}
	ord, ok := resultByName(span, "ord")
	if !ok {
		t.Fatalf("Results = %+v, want a column named ord", span.Results)
	}
	// The ordinal column derives from no base column; it must not pick up phones.
	if hasSource(ord.SourceColumns, ColumnRef{Column: "phones"}) {
		t.Errorf("ord.SourceColumns = %+v, must not contain {Column:phones} (ordinal has no lineage)", ord.SourceColumns)
	}
}

// TestGetQuerySpan_ScalarSubqueryLineage covers BYT-9676: a scalar subquery used
// as a SELECT value must resolve to its inner output column's lineage.
func TestGetQuerySpan_ScalarSubqueryLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT (SELECT phone FROM customer LIMIT 1) AS sp FROM customer")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "sp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named sp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("sp.SourceColumns = %+v, want to contain {Column:phone} (resolved through scalar subquery)", r.SourceColumns)
	}
}

// TestGetQuerySpan_ScalarSubqueryInExprLineage resolves a scalar subquery nested
// inside a larger select expression.
func TestGetQuerySpan_ScalarSubqueryInExprLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT id, (SELECT max(phone) FROM customer) AS mp FROM orders")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "mp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named mp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("mp.SourceColumns = %+v, want to contain {Column:phone}", r.SourceColumns)
	}
}

// TestGetQuerySpan_ScalarSubqueryThroughDerived resolves a scalar subquery that
// is projected by a derived table (composition of both mechanisms).
func TestGetQuerySpan_ScalarSubqueryThroughDerived(t *testing.T) {
	span, err := GetQuerySpan("SELECT d.sp FROM (SELECT (SELECT phone FROM customer LIMIT 1) AS sp FROM x) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "sp")
	if !ok {
		t.Fatalf("Results = %+v, want a column named sp", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("sp.SourceColumns = %+v, want to contain {Column:phone} (scalar subquery through derived table)", r.SourceColumns)
	}
}

// TestGetQuerySpan_UnnestOfScalarSubqueryLineage covers the intersection of
// UNNEST and a scalar-subquery argument, flagged during review: UNNEST of a
// subquery-returned array must resolve to the subquery's sensitive column.
func TestGetQuerySpan_UnnestOfScalarSubqueryLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT p FROM UNNEST((SELECT array_agg(phone) FROM customer)) AS t(p)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "p")
	if !ok {
		t.Fatalf("Results = %+v, want a column named p", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("p.SourceColumns = %+v, want to contain {Column:phone} (UNNEST of scalar subquery)", r.SourceColumns)
	}
}

// TestGetQuerySpan_ExistsSubqueryNotResolved guards that an EXISTS subquery (a
// boolean, not a value) does not contribute its inner columns as lineage.
func TestGetQuerySpan_ExistsSubqueryNotResolved(t *testing.T) {
	span, err := GetQuerySpan("SELECT EXISTS (SELECT 1 FROM customer WHERE phone IS NOT NULL) AS e FROM orders")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "e")
	if !ok {
		t.Fatalf("Results = %+v, want a column named e", span.Results)
	}
	if hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("e.SourceColumns = %+v, must not contain {Column:phone} (EXISTS yields a boolean)", r.SourceColumns)
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
