package analysis

import (
	"strings"
	"testing"
)

// resultNames returns the result column names in order.
func resultNames(span *QuerySpan) []string {
	out := make([]string, 0, len(span.Results))
	for _, r := range span.Results {
		out = append(out, r.Name)
	}
	return out
}

func sameNames(got []string, want ...string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestGetQuerySpan_StarOverDerivedExpands covers BYT-9678: SELECT * over a
// derived table must expand to the derived relation's exact projection (width
// and order), each column carrying its base lineage — not remain an opaque "*"
// that the consumer expands to the base table's full column set (which
// misaligns the positional masker).
func TestGetQuerySpan_StarOverDerivedExpands(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT phone, name FROM customer) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phone", "name") {
		t.Fatalf("Results = %v, want exactly [phone name] in order", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[0].SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}
	if !hasSource(span.Results[1].SourceColumns, ColumnRef{Column: "name"}) {
		t.Errorf("Results[1].SourceColumns = %+v, want {Column:name}", span.Results[1].SourceColumns)
	}
}

// TestGetQuerySpan_StarOverCTEExpands expands SELECT * over a CTE reference.
func TestGetQuerySpan_StarOverCTEExpands(t *testing.T) {
	span, err := GetQuerySpan("WITH w AS (SELECT phone, name FROM customer) SELECT * FROM w")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phone", "name") {
		t.Fatalf("Results = %v, want [phone name]", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[0].SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}
}

// TestGetQuerySpan_QualifiedStarOverDerivedExpands expands d.* over a derived
// table, including a renamed projection.
func TestGetQuerySpan_QualifiedStarOverDerivedExpands(t *testing.T) {
	span, err := GetQuerySpan("SELECT d.* FROM (SELECT phone AS x FROM customer) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "x") {
		t.Fatalf("Results = %v, want [x]", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("x.SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}
}

// TestGetQuerySpan_StarOverBaseStaysOpaque guards parity: a star over a base
// table cannot be expanded without metadata and must stay a single "*" result
// for the consumer's metadata-based expansion.
func TestGetQuerySpan_StarOverBaseStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM customer")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*]", got)
	}
}

// TestGetQuerySpan_StarOverMixedBaseDerivedStaysOpaque guards parity: a star
// covering a base table joined with a derived relation has unknown total width
// (the base side needs metadata) and must stay opaque.
func TestGetQuerySpan_StarOverMixedBaseDerivedStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM customer JOIN (SELECT x FROM t) d ON true")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (mixed base+derived must not expand)", got)
	}
}

// TestGetQuerySpan_StarOverUsingJoinStaysOpaque guards that a USING join blocks
// expansion: the coalesced join column changes the output width/order.
func TestGetQuerySpan_StarOverUsingJoinStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("WITH a AS (SELECT k, p FROM t1), b AS (SELECT k, q FROM t2) SELECT * FROM a JOIN b USING (k)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (USING join coalesces; must not expand)", got)
	}
}

// TestGetQuerySpan_StarOverDerivedContainingBaseStarStaysOpaque guards that a
// derived relation whose own projection is an unexpandable star (over a base
// table) makes the outer star opaque too — its width is unknown.
func TestGetQuerySpan_StarOverDerivedContainingBaseStarStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT * FROM customer) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (inner base star is opaque)", got)
	}
}

// TestGetQuerySpan_NestedDerivedStarExpands expands a star whose derived
// relation itself contains a resolved star over another derived relation.
func TestGetQuerySpan_NestedDerivedStarExpands(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT * FROM (SELECT phone AS x FROM customer) i) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "x") {
		t.Fatalf("Results = %v, want [x] (nested resolved stars expand transitively)", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("x.SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}
}

// TestGetQuerySpan_MixedItemsStarSplice verifies the splice keeps non-star items
// aligned: SELECT a, * over a derived relation yields [a, b, c].
func TestGetQuerySpan_MixedItemsStarSplice(t *testing.T) {
	span, err := GetQuerySpan("SELECT a, * FROM (SELECT b, c FROM t) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "a", "b", "c") {
		t.Fatalf("Results = %v, want [a b c]", got)
	}
}

// TestGetQuerySpan_SetOpResolvedStarArmMerges verifies that a set operation
// whose left arm is a RESOLVED star merges the right arm positionally at the
// expanded width: the sensitive right-arm column lands on the expanded column.
func TestGetQuerySpan_SetOpResolvedStarArmMerges(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT name FROM customer) d UNION SELECT phone FROM customer")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "name") {
		t.Fatalf("Results = %v, want [name] (expanded star arm)", got)
	}
	srcs := span.Results[0].SourceColumns
	if !hasSource(srcs, ColumnRef{Column: "name"}) || !hasSource(srcs, ColumnRef{Column: "phone"}) {
		t.Errorf("name.SourceColumns = %+v, want both {Column:name} and {Column:phone} (arms merged at expanded width)", srcs)
	}
}

// TestGetQuerySpan_StarOverUnaliasedUnnestStaysOpaque guards that UNNEST without
// column aliases (output width unknown without type metadata) blocks expansion.
func TestGetQuerySpan_StarOverUnaliasedUnnestStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT phones FROM customer) d CROSS JOIN UNNEST(d.phones)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "*") {
		t.Fatalf("Results = %v, want exactly [*] (unaliased UNNEST width unknown)", got)
	}
}

// TestGetQuerySpan_StarOverDerivedWithUnnestExpands expands a star over a
// derived relation joined with an ALIASED UNNEST (both widths known).
func TestGetQuerySpan_StarOverDerivedWithUnnestExpands(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT phones FROM customer) d CROSS JOIN UNNEST(d.phones) AS t(p)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phones", "p") {
		t.Fatalf("Results = %v, want [phones p]", got)
	}
	if !hasSource(span.Results[1].SourceColumns, ColumnRef{Column: "phones"}) {
		t.Errorf("p.SourceColumns = %+v, want {Column:phones}", span.Results[1].SourceColumns)
	}
}

// TestGetQuerySpan_QualifiedStarThroughInnerStarResolves covers the transitive
// fix from inline star expansion: a qualified outer ref through a derived
// relation whose projection is a resolved star now finds the column.
func TestGetQuerySpan_QualifiedStarThroughInnerStarResolves(t *testing.T) {
	span, err := GetQuerySpan("SELECT d.phone FROM (SELECT * FROM (SELECT phone FROM customer) c) d")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "phone")
	if !ok {
		t.Fatalf("Results = %+v, want a column named phone", span.Results)
	}
	if !hasSource(r.SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("phone.SourceColumns = %+v, want {Column:phone} (through inner resolved star)", r.SourceColumns)
	}
}

// TestGetQuerySpan_AliasedJoinStarExpands covers an alias on a parenthesized
// join: ((…) a JOIN (…) b) AS j is one relation whose projection is the join's
// column concatenation (verified against Trino 481), so j.* and * expand.
func TestGetQuerySpan_AliasedJoinStarExpands(t *testing.T) {
	span, err := GetQuerySpan("SELECT j.* FROM ((SELECT phone FROM customer) a JOIN (SELECT name FROM customer) b ON true) AS j")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phone", "name") {
		t.Fatalf("Results = %v, want [phone name]", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[0].SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}

	span, err = GetQuerySpan("SELECT * FROM ((SELECT phone FROM customer) a JOIN (SELECT name FROM customer) b ON true) AS j")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phone", "name") {
		t.Fatalf("unqualified Results = %v, want [phone name] (the aliased join is ONE relation; its parts must not double-count)", got)
	}
}

// TestGetQuerySpan_AliasedJoinWithBaseStaysOpaque guards that an aliased join
// containing a base table is bound as an unresolvable relation: stars through
// it stay opaque.
func TestGetQuerySpan_AliasedJoinWithBaseStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("SELECT j.* FROM (customer c JOIN (SELECT x FROM t) d ON true) AS j")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "j.*") {
		t.Fatalf("Results = %v, want exactly [j.*] (base table inside the aliased join)", got)
	}
}

// TestGetQuerySpan_TableQueryOverCTEExpands covers TABLE <cte>: equivalent to
// SELECT * FROM cte, whose projection is the CTE's resolved columns (verified
// against Trino 481).
func TestGetQuerySpan_TableQueryOverCTEExpands(t *testing.T) {
	span, err := GetQuerySpan("WITH w AS (SELECT phone, name FROM customer) TABLE w")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "phone", "name") {
		t.Fatalf("Results = %v, want [phone name]", got)
	}
	if !hasSource(span.Results[0].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[0].SourceColumns = %+v, want {Column:phone}", span.Results[0].SourceColumns)
	}
}

// TestGetQuerySpan_TopLevelValuesResults covers a top-level VALUES: the primary
// walk emits no Results (no select items), which left the consumer with zero
// positional maskers; the resolver must synthesize the projection so a
// sensitive value in a VALUES row is maskable.
func TestGetQuerySpan_TopLevelValuesResults(t *testing.T) {
	span, err := GetQuerySpan("VALUES ('x', (SELECT phone FROM customer LIMIT 1))")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("Results = %+v, want exactly 2 columns", span.Results)
	}
	if !hasSource(span.Results[1].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("Results[1].SourceColumns = %+v, want {Column:phone}", span.Results[1].SourceColumns)
	}
}

// TestGetQuerySpan_ValuesRowConstructorArity locks the oracle-verified VALUES
// arity: ROW('a', expr) in VALUES is unpacked into 2 output columns exactly
// like ('a', expr) (verified against Trino 481, where AS v(r) over it errors
// with "alias list has 1 entries but 2 columns").
func TestGetQuerySpan_ValuesRowConstructorArity(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (VALUES ROW(1, (SELECT phone FROM customer LIMIT 1))) AS v(a, b)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "a", "b") {
		t.Fatalf("Results = %v, want [a b] (ROW unpacks to 2 columns)", got)
	}
	if !hasSource(span.Results[1].SourceColumns, ColumnRef{Column: "phone"}) {
		t.Errorf("b.SourceColumns = %+v, want {Column:phone}", span.Results[1].SourceColumns)
	}
}

// TestGetQuerySpan_QualifiedStarBesideUsingJoinStaysOpaque locks the
// oracle-verified USING semantics: with a(k, phone), `SELECT a.* FROM a JOIN b
// USING (k)` returns only [phone] — Trino EXCLUDES the using columns from a
// qualified star — so expanding a.* to a's full projection would be
// width-wrong. The conservative scope-wide block is correct.
func TestGetQuerySpan_QualifiedStarBesideUsingJoinStaysOpaque(t *testing.T) {
	span, err := GetQuerySpan("WITH a AS (SELECT k, phone FROM customer), b AS (SELECT k FROM orders) SELECT a.* FROM a JOIN b USING (k)")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	if got := resultNames(span); !sameNames(got, "a.*") {
		t.Fatalf("Results = %v, want exactly [a.*] (USING strips join columns from a qualified star; expansion would misalign)", got)
	}
}

// TestGetQuerySpan_QualifiedStarOverBaseShapePreserved guards the consumer
// contract for an OPAQUE qualified star: its result keeps exactly the walk's
// single relation-name source ref (consumers detect a qualified star by that
// shape: one ref whose Column equals the qualifier), even when a derived
// relation in scope exposes a column of the same name.
func TestGetQuerySpan_QualifiedStarOverBaseShapePreserved(t *testing.T) {
	span, err := GetQuerySpan("SELECT u.*, d.x FROM users u JOIN (SELECT phone AS u, phone AS x FROM customer) d ON true")
	if err != nil {
		t.Fatalf("GetQuerySpan returned error: %v", err)
	}
	r, ok := resultByName(span, "u.*")
	if !ok {
		t.Fatalf("Results = %+v, want a column named u.*", span.Results)
	}
	if len(r.SourceColumns) != 1 || !strings.EqualFold(r.SourceColumns[0].Column, "u") {
		t.Errorf("u.*.SourceColumns = %+v, want exactly the single relation ref {Column:u} (shape consumers key on)", r.SourceColumns)
	}
}
