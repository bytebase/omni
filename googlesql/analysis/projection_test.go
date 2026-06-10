package analysis

import (
	"sort"
	"strings"
	"testing"
)

// These tests pin the relation-projection resolution contract (the masking-grade
// core): a `SELECT *` / `rel.*` / column reference over a CTE, a derived
// subquery, or a join resolves to that relation's RESOLVED output projection —
// concrete columns with base-column lineage, or a base-table-star segment the
// catalog-aware consumer expands — and a set operation position-merges resolved
// projections. They assert the omni-side ColumnInfo wire form (the bytebase
// extractor's behaviour on top of it is covered by the bigquery query-span
// corpus, recorded from the legacy resolver).

// resolvedSources renders one resolved ColumnInfo's source lineage as sorted
// "table.column" (or "column" when unqualified) strings.
func resolvedSources(info ColumnInfo) []string {
	var out []string
	for _, sc := range info.SourceColumns {
		if sc.Table != "" {
			out = append(out, sc.Table+"."+sc.Column)
		} else {
			out = append(out, sc.Column)
		}
	}
	sort.Strings(out)
	return out
}

// TestProjection_CTESubsetThenStar: a CTE projecting a SUBSET of a wider base
// table, read by `SELECT *`, resolves the star to ONLY the CTE's projected column
// (with base lineage) — not the base table's other columns. This is the core
// under-attribution guard: a `*` over `WITH c AS (SELECT email FROM users)` must
// surface only email, never ssn/name.
func TestProjection_CTESubsetThenStar(t *testing.T) {
	span, err := GetQuerySpan("WITH c AS (SELECT email FROM users) SELECT * FROM c", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("got %d results, want 1 (the CTE projects exactly email)", len(span.Results))
	}
	r := span.Results[0]
	if len(r.StarSegments) != 0 {
		t.Errorf("result should be a resolved concrete column, not a star-segment item; got %d segments", len(r.StarSegments))
	}
	if !strings.EqualFold(r.Name, "email") {
		t.Errorf("name = %q, want email", r.Name)
	}
	if r.IsPlain {
		t.Errorf("IsPlain = true, want false (a CTE column from an explicit select is not a plain field)")
	}
	if got := resolvedSources(r); !eqStrings(got, []string{"users.email"}) {
		t.Errorf("sources = %v, want [users.email] (the sole-base-relation attribution ties the bare column to its FROM table)", got)
	}
}

// TestProjection_DerivedStar: a `SELECT *` over a derived subquery reproduces the
// subquery's resolved projection (its explicit columns), not the base table's
// columns.
func TestProjection_DerivedStar(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT a, b FROM t) d", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"a", "b"}) {
		t.Errorf("results = %v, want [a b] (the derived projection, not t's full columns)", got)
	}
	for i, r := range span.Results {
		if len(r.StarSegments) != 0 {
			t.Errorf("result[%d] should be concrete, got a star-segment item", i)
		}
	}
}

// TestProjection_DerivedQualifiedColumnAndStar pins the "Table subquery" shape: a
// qualified reference through a derived relation that wraps a base-table star
// (`result.ID` over `(SELECT * FROM A)` where A = `SELECT * FROM people`)
// resolves to the base column people.ID, while a sibling `*` reproduces the whole
// base-table star (one StarSegment over people).
func TestProjection_DerivedQualifiedColumnAndStar(t *testing.T) {
	span, err := GetQuerySpan("WITH A AS (SELECT * FROM people) SELECT result.ID, * FROM (SELECT * FROM A) result", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("got %d results, want 2 (result.ID + the star)", len(span.Results))
	}
	// result.ID resolves through the derived relation and CTE A to people.ID.
	idCol := span.Results[0]
	if !strings.EqualFold(idCol.Name, "ID") {
		t.Errorf("result[0] name = %q, want ID", idCol.Name)
	}
	if got := resolvedSources(idCol); !eqStrings(got, []string{"people.ID"}) {
		t.Errorf("result.ID sources = %v, want [people.ID]", got)
	}
	// The sibling * is a single base-table-star segment over people (the consumer
	// enumerates its columns).
	star := span.Results[1]
	if len(star.StarSegments) != 1 || star.StarSegments[0].BaseTable == nil {
		t.Fatalf("result[1] should be a single base-table-star segment, got %+v", star.StarSegments)
	}
	if !strings.EqualFold(star.StarSegments[0].BaseTable.Table, "people") {
		t.Errorf("star segment base table = %q, want people", star.StarSegments[0].BaseTable.Table)
	}
	if !star.StarSegments[0].Plain {
		t.Errorf("a base-table * passthrough column should be a plain field")
	}
}

// TestProjection_SetOpMergesResolvedCTEs: a set operation of two subset-projecting
// CTEs position-merges their resolved projections — output i draws from BOTH arms
// at position i — with names from the left arm.
func TestProjection_SetOpMergesResolvedCTEs(t *testing.T) {
	sql := "WITH a AS (SELECT id, name FROM users), b AS (SELECT pid, label FROM members) " +
		"SELECT * FROM a UNION ALL SELECT * FROM b"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"id", "name"}) {
		t.Fatalf("names = %v, want [id name] (left arm names)", got)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"members.pid", "users.id"}) {
		t.Errorf("result[0] sources = %v, want [members.pid users.id] (both arms position 0, table-attributed)", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"members.label", "users.name"}) {
		t.Errorf("result[1] sources = %v, want [members.label users.name] (both arms position 1, table-attributed)", got)
	}
}

// TestProjection_SetOpConcreteAndBaseStarArm pins the StarMerge contract: a
// `concrete UNION ALL (SELECT * FROM base)` keeps the concrete arm's per-position
// lineage AND attaches a StarMerge so the consumer expands the base table and
// position-merges its i-th column (omni cannot know the base table's arity).
func TestProjection_SetOpConcreteAndBaseStarArm(t *testing.T) {
	span, err := GetQuerySpan("SELECT a, b FROM t1 UNION ALL SELECT * FROM t2", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 2 {
		t.Fatalf("got %d results, want 2 (left arm arity)", len(span.Results))
	}
	for i, r := range span.Results {
		if r.StarMerge == nil {
			t.Fatalf("result[%d] should carry a StarMerge for the base-star arm", i)
		}
		if !strings.EqualFold(r.StarMerge.Table.Table, "t2") {
			t.Errorf("result[%d] StarMerge table = %q, want t2", i, r.StarMerge.Table.Table)
		}
		if r.StarMerge.Index != i {
			t.Errorf("result[%d] StarMerge index = %d, want %d", i, r.StarMerge.Index, i)
		}
		if r.StarMerge.LeftStar {
			t.Errorf("result[%d] StarMerge.LeftStar = true, want false (the star arm is the RIGHT one)", i)
		}
	}
}

// TestProjection_BaseTableStarSegment: a bare `SELECT *` over a single base table
// is a single base-table-star segment (the consumer enumerates the columns), not
// a flat set of concrete columns omni cannot know.
func TestProjection_BaseTableStarSegment(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM people", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 1 || len(span.Results[0].StarSegments) != 1 {
		t.Fatalf("want 1 star-segment result, got %+v", span.Results)
	}
	seg := span.Results[0].StarSegments[0]
	if seg.BaseTable == nil || !strings.EqualFold(seg.BaseTable.Table, "people") {
		t.Errorf("segment = %+v, want a base-table segment over people", seg)
	}
}

// TestProjection_StarOverJoinSegmentsAndPlainness: a bare `*` over a JOIN of base
// tables yields one base-table-star item per side, in FROM order (a modifier-free
// star inlines its segments), with the legacy join plain-field rule — the LEFT/
// anchor side's columns are NOT plain, the right side's are.
func TestProjection_StarOverJoinSegmentsAndPlainness(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM a JOIN b ON a.id = b.aid", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	// A modifier-free `*` over a join inlines one base-table-star item per side.
	if len(span.Results) != 2 {
		t.Fatalf("got %d results, want 2 (one base-table star per join side)", len(span.Results))
	}
	left := span.Results[0]
	right := span.Results[1]
	if len(left.StarSegments) != 1 || left.StarSegments[0].BaseTable == nil ||
		!strings.EqualFold(left.StarSegments[0].BaseTable.Table, "a") || left.StarSegments[0].Plain {
		t.Errorf("result[0] = %+v, want a base table a star NOT plain (join anchor/left side is rewrapped)", left.StarSegments)
	}
	if len(right.StarSegments) != 1 || right.StarSegments[0].BaseTable == nil ||
		!strings.EqualFold(right.StarSegments[0].BaseTable.Table, "b") || !right.StarSegments[0].Plain {
		t.Errorf("result[1] = %+v, want a base table b star plain (right side keeps plainness)", right.StarSegments)
	}
}

// TestProjection_QualifiedColumnThroughCTE: a qualified `alias.col` over a CTE
// reference resolves to the CTE column's lineage, even when the CTE is self-joined
// under two aliases. The CTE column's lineage is the body's resolved source — here
// a bare `id` / `parent_id` (the CTE body selected them over the base table
// `nodes`, whose columns omni does not enumerate), which the catalog-aware
// consumer matches back to nodes.id / nodes.parent_id (verified in the bigquery
// corpus). The point pinned here is that `x.id` resolves THROUGH the CTE relation
// (not to an empty/aliased ref).
func TestProjection_QualifiedColumnThroughCTE(t *testing.T) {
	sql := "WITH c AS (SELECT id, parent_id FROM nodes) " +
		"SELECT x.id, y.parent_id FROM c x JOIN c y ON x.parent_id = y.id"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"nodes.id"}) {
		t.Errorf("x.id sources = %v, want [nodes.id] (CTE column attributed to its base table)", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"nodes.parent_id"}) {
		t.Errorf("y.parent_id sources = %v, want [nodes.parent_id] (CTE column attributed to its base table)", got)
	}
}

// TestProjection_RecursiveCTEFixpoint: a recursive CTE's projection resolves the
// recursive arm's self-references against the anchor and iterates to a fixpoint,
// so a column whose recursive expression reads another column accumulates that
// column's lineage (c3 = c3*c2 picks up c2's {a,b}).
func TestProjection_RecursiveCTEFixpoint(t *testing.T) {
	sql := "WITH RECURSIVE c AS (" +
		"(SELECT a AS c1, b AS c2, c AS c3 FROM t) " +
		"UNION ALL SELECT c1*c2, c2+c1, c3*c2 FROM c" +
		") SELECT c1, c2, c3 FROM c"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 3 {
		t.Fatalf("got %d results, want 3", len(span.Results))
	}
	// c1 = c1*c2 → {a,b}; c2 = c2+c1 → {a,b}; c3 = c3*c2 → {a,b,c} via the fixpoint
	// (c2 grew to {a,b} in an earlier pass, then c3 reads c2 and picks up a).
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t.a", "t.b"}) {
		t.Errorf("c1 sources = %v, want [a b]", got)
	}
	if got := resolvedSources(span.Results[2]); !eqStrings(got, []string{"t.a", "t.b", "t.c"}) {
		t.Errorf("c3 sources = %v, want [a b c] (fixpoint propagates c2's growth)", got)
	}
}

// TestProjection_StructFieldRootIsColumn: a `s.field` / `s.a.b` where `s` is not a
// FROM relation resolves `s` as a (struct) column of the base table, dropping the
// trailing field path — the legacy a-as-column field-path behaviour.
func TestProjection_StructFieldRootIsColumn(t *testing.T) {
	for _, sql := range []string{"SELECT s.field FROM t", "SELECT s.a.b FROM t"} {
		span, err := GetQuerySpan(sql, DialectBigQuery)
		if err != nil {
			t.Fatalf("%q error: %v", sql, err)
		}
		if len(span.Results) != 1 {
			t.Fatalf("%q got %d results, want 1", sql, len(span.Results))
		}
		if !strings.EqualFold(span.Results[0].Name, "s") {
			t.Errorf("%q name = %q, want s (the struct column root)", sql, span.Results[0].Name)
		}
		if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t.s"}) {
			t.Errorf("%q sources = %v, want [t.s] (struct root attributed to its sole base relation)", sql, got)
		}
	}
}

// TestProjection_QualifiedColumnKeepsDialectBuckets: a fully-qualified column whose
// qualifier names a FROM base table (`proj.ds.t.c`) keeps its dialect-bucketed
// qualifier so it lines up with that table's access, rather than being misread as
// a struct-column root.
func TestProjection_QualifiedColumnKeepsDialectBuckets(t *testing.T) {
	span, err := GetQuerySpan("SELECT proj.ds.t.c FROM proj.ds.t", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 1 || len(span.Results[0].SourceColumns) != 1 {
		t.Fatalf("want 1 result with 1 source, got %+v", span.Results)
	}
	c := span.Results[0].SourceColumns[0]
	if c.Catalog != "proj" || c.Database != "ds" || c.Table != "t" || c.Column != "c" {
		t.Errorf("ref = {Catalog:%q Database:%q Schema:%q Table:%q Column:%q}, want {proj ds \"\" t c}",
			c.Catalog, c.Database, c.Schema, c.Table, c.Column)
	}
}

// TestProjection_CTEReferenceRecorded: a CTE reference (`FROM c`) is recorded in
// CTEReferences (not AccessTables), so the consumer can surface it in the
// table-level source set while never column-expanding it.
func TestProjection_CTEReferenceRecorded(t *testing.T) {
	span, err := GetQuerySpan("WITH c AS (SELECT id FROM users) SELECT id FROM c", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := tableNames(span); !eqStrings(got, []string{"users"}) {
		t.Errorf("AccessTables = %v, want [users] (the CTE reference c is NOT a base table)", got)
	}
	if len(span.CTEReferences) != 1 || !strings.EqualFold(span.CTEReferences[0].Table, "c") {
		t.Errorf("CTEReferences = %+v, want one reference to c", span.CTEReferences)
	}
}

// TestProjection_LambdaParamExcludedFromResolvedLineage guards the resolver path
// of the lambda-parameter exclusion: a bound lambda parameter must not appear in
// a resolved select item's source lineage (only the free columns do).
func TestProjection_LambdaParamExcludedFromResolvedLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT ARRAY_TRANSFORM(arr, e -> e + threshold) AS r FROM t", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t.arr", "t.threshold"}) {
		t.Errorf("sources = %v, want [t.arr t.threshold] (lambda param e excluded; table-attributed)", got)
	}
}

// ---------------------------------------------------------------------------
// Fix F1: SELECT * over JOIN ... USING coalesces the key columns.
//
// GoogleSQL's `SELECT *` over `lt JOIN rt USING (k)` returns the key column
// ONCE (first, coalesced from both sides), then the left side's non-key
// columns, then the right side's. Concatenating both sides' stars instead
// shifts every position after the key — the positional masker then applies
// rt.k's policy to the column that actually holds rt's first NON-key (secret)
// column: a masking under-attribution leak. The legacy resolver coalesced only
// upper-case-written keys (a case-fold bug — `USING (k)` leaked); omni
// coalesces case-insensitively (GoogleSQL identifiers are case-insensitive),
// which is strictly safer, so the lowercase shape is pinned HERE (it cannot be
// corpus-recorded from legacy without recording the leak).
// ---------------------------------------------------------------------------

// TestProjection_JoinUsingStarCoalescesLowercase pins the F1 wire shape over
// base tables with a lowercase key: one concrete coalesced key column (lineage
// from BOTH sides), then one base-table star per side carrying the key in its
// ExceptColumns so the consumer's expansion skips it.
func TestProjection_JoinUsingStarCoalescesLowercase(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM lt JOIN rt USING (k)", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 3 {
		t.Fatalf("got %d results, want 3 (coalesced k + lt-star + rt-star), results=%+v", len(span.Results), span.Results)
	}
	key := span.Results[0]
	if key.StarSegments != nil || !strings.EqualFold(key.Name, "k") {
		t.Fatalf("result[0] = %+v, want a concrete coalesced column k", key)
	}
	if got := resolvedSources(key); !eqStrings(got, []string{"lt.k", "rt.k"}) {
		t.Errorf("coalesced k sources = %v, want [lt.k rt.k] (BOTH sides)", got)
	}
	if key.IsPlain {
		t.Errorf("coalesced key must not be a plain field")
	}
	left, right := span.Results[1], span.Results[2]
	if len(left.StarSegments) != 1 || left.StarSegments[0].BaseTable == nil ||
		!strings.EqualFold(left.StarSegments[0].BaseTable.Table, "lt") {
		t.Fatalf("result[1] = %+v, want a base-table star over lt", left.StarSegments)
	}
	if got := left.StarSegments[0].ExceptColumns; len(got) != 1 || !strings.EqualFold(got[0], "k") {
		t.Errorf("lt star ExceptColumns = %v, want [k] (the coalesced key is excluded)", got)
	}
	if left.StarSegments[0].Plain {
		t.Errorf("join-left star must be rewrapped non-plain")
	}
	if len(right.StarSegments) != 1 || right.StarSegments[0].BaseTable == nil ||
		!strings.EqualFold(right.StarSegments[0].BaseTable.Table, "rt") {
		t.Fatalf("result[2] = %+v, want a base-table star over rt", right.StarSegments)
	}
	if got := right.StarSegments[0].ExceptColumns; len(got) != 1 || !strings.EqualFold(got[0], "k") {
		t.Errorf("rt star ExceptColumns = %v, want [k]", got)
	}
	if !right.StarSegments[0].Plain {
		t.Errorf("join-right star keeps plainness")
	}
}

// TestProjection_JoinUsingCaseInsensitive: the key match is case-insensitive
// (GoogleSQL identifier rule): an uppercase-written `USING (K)` coalesces the
// same way (this is also the legacy-recordable corpus shape).
func TestProjection_JoinUsingCaseInsensitive(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM lt JOIN rt USING (K)", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 3 {
		t.Fatalf("got %d results, want 3, results=%+v", len(span.Results), span.Results)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"lt.K", "rt.K"}) {
		t.Errorf("coalesced K sources = %v, want [lt.K rt.K]", got)
	}
}

// TestProjection_JoinUsingMultiKeyAndChain: multi-key USING coalesces each key
// (omni parses `USING (a, b)`; legacy could not), and chained joins compose —
// the outer join's keys come first, then the left side's remaining columns
// (which start with the inner join's coalesced keys), then the right side's.
func TestProjection_JoinUsingMultiKeyAndChain(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM a JOIN b USING (k, j) JOIN c USING (m)", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 6 {
		t.Fatalf("got %d results, want 6 (m, k, j, a*, b*, c*), results=%+v", len(span.Results), span.Results)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"a.m", "b.m", "c.m"}) {
		t.Errorf("outer key m sources = %v, want [a.m b.m c.m] (left side's stars may each own m; over-attribution is safe)", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"a.k", "b.k"}) {
		t.Errorf("inner key k sources = %v, want [a.k b.k]", got)
	}
	if got := resolvedSources(span.Results[2]); !eqStrings(got, []string{"a.j", "b.j"}) {
		t.Errorf("inner key j sources = %v, want [a.j b.j]", got)
	}
	wantExcept := map[int][]string{3: {"k", "j", "m"}, 4: {"k", "j", "m"}, 5: {"m"}}
	wantTable := map[int]string{3: "a", 4: "b", 5: "c"}
	for i := 3; i <= 5; i++ {
		r := span.Results[i]
		if len(r.StarSegments) != 1 || r.StarSegments[0].BaseTable == nil ||
			!strings.EqualFold(r.StarSegments[0].BaseTable.Table, wantTable[i]) {
			t.Fatalf("result[%d] = %+v, want base star over %s", i, r.StarSegments, wantTable[i])
		}
		got := append([]string{}, r.StarSegments[0].ExceptColumns...)
		want := append([]string{}, wantExcept[i]...)
		for i := range got {
			got[i] = strings.ToLower(got[i])
		}
		sort.Strings(got)
		sort.Strings(want)
		if !eqStrings(got, want) {
			t.Errorf("result[%d] ExceptColumns = %v, want %v", i, got, want)
		}
	}
}

// TestProjection_JoinUsingConcreteSides: USING over two CONCRETE (CTE) sides
// name-filters inline — the coalesced key unions both CTE columns' base
// lineage and each side's non-key columns follow.
func TestProjection_JoinUsingConcreteSides(t *testing.T) {
	sql := "WITH lc AS (SELECT lk AS k, note FROM t1), rc AS (SELECT rk AS k, sec FROM t2) " +
		"SELECT * FROM lc JOIN rc USING (k)"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"k", "note", "sec"}) {
		t.Fatalf("names = %v, want [k note sec]", got)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t1.lk", "t2.rk"}) {
		t.Errorf("k sources = %v, want [lk rk]", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"t1.note"}) {
		t.Errorf("note sources = %v, want [note]", got)
	}
	if got := resolvedSources(span.Results[2]); !eqStrings(got, []string{"t2.sec"}) {
		t.Errorf("sec sources = %v, want [sec]", got)
	}
}

// TestProjection_JoinUsingStarGroupSide: a USING side that is an EXCEPT-star
// derived table keeps its star-group (the key joins its EXCEPT set; the key's
// lineage is read from the group's segments).
func TestProjection_JoinUsingStarGroupSide(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM (SELECT * EXCEPT (x) FROM a) d JOIN b USING (k)", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 3 {
		t.Fatalf("got %d results, want 3 (k + d's star-group + b's star), results=%+v", len(span.Results), span.Results)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"a.k", "b.k"}) {
		t.Errorf("k sources = %v, want [a.k b.k]", got)
	}
	grp := span.Results[1]
	if len(grp.StarSegments) != 1 || grp.StarSegments[0].BaseTable == nil ||
		!strings.EqualFold(grp.StarSegments[0].BaseTable.Table, "a") {
		t.Fatalf("result[1] = %+v, want d's star-group over a", grp)
	}
	gotExcept := append([]string{}, grp.StarExcept...)
	for i := range gotExcept {
		gotExcept[i] = strings.ToLower(gotExcept[i])
	}
	sort.Strings(gotExcept)
	if !eqStrings(gotExcept, []string{"k", "x"}) {
		t.Errorf("d star-group EXCEPT = %v, want [k x] (original EXCEPT + the coalesced key)", gotExcept)
	}
}

// TestProjection_NaturalJoinFailsClosed: both BigQuery and Spanner reject
// NATURAL JOIN at analysis, and omni cannot know the shared columns without a
// catalog — silent concatenation would misalign positions, so the analysis
// fails closed (structural rule: correct lineage or an error, never silent
// misalignment for shapes legacy could not parse).
func TestProjection_NaturalJoinFailsClosed(t *testing.T) {
	if _, err := GetQuerySpan("SELECT * FROM a NATURAL JOIN b", DialectBigQuery); err == nil {
		t.Fatalf("NATURAL JOIN: want a fail-closed error, got nil")
	}
}

// TestProjection_JoinUsingDeferredSideFailsClosed: a USING side whose
// projection carries a deferred set-operation marker (its arity is known only
// to the metadata-aware consumer) cannot be name-partitioned here; the
// analysis fails closed rather than emitting misaligned lineage.
func TestProjection_JoinUsingDeferredSideFailsClosed(t *testing.T) {
	sql := "SELECT * FROM (SELECT a, b FROM t UNION ALL SELECT * FROM u) d JOIN r USING (a)"
	if _, err := GetQuerySpan(sql, DialectBigQuery); err == nil {
		t.Fatalf("USING over a deferred set-op side: want a fail-closed error, got nil")
	}
}

// ---------------------------------------------------------------------------
// Fix F2: UNNEST table sources produce a value relation with element lineage.
//
// `FROM victim, UNNEST(victim.secret_tokens) AS elem` projects elem as a real
// output column whose data IS victim.secret_tokens' elements. Dropping the
// relation made `SELECT elem` resolve to an empty/unmatched lineage — the
// fail-open masker then returned the secret unmasked. (The legacy extractor
// ERRORED on UNNEST sources = fail-closed; omni resolves the lineage.)
// ---------------------------------------------------------------------------

// TestProjection_UnnestElementLineage pins the core F2 contract.
func TestProjection_UnnestElementLineage(t *testing.T) {
	span, err := GetQuerySpan("SELECT elem FROM victim, UNNEST(victim.secret_tokens) AS elem", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 1 {
		t.Fatalf("got %d results, want 1, results=%+v", len(span.Results), span.Results)
	}
	r := span.Results[0]
	if !strings.EqualFold(r.Name, "elem") {
		t.Errorf("name = %q, want elem", r.Name)
	}
	if got := resolvedSources(r); !eqStrings(got, []string{"victim.secret_tokens"}) {
		t.Errorf("elem sources = %v, want [victim.secret_tokens]", got)
	}
	if got := tableNames(span); !eqStrings(got, []string{"victim"}) {
		t.Errorf("tables = %v, want [victim]", got)
	}
}

// TestProjection_UnnestStarAndOffset: a bare `*` includes the unnest element
// column (after the base table's star) and the WITH OFFSET companion column;
// the offset is positional, not data, so it has NO lineage.
func TestProjection_UnnestStarAndOffset(t *testing.T) {
	span, err := GetQuerySpan("SELECT * FROM victim, UNNEST(victim.secret_tokens) AS elem WITH OFFSET AS pos", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 3 {
		t.Fatalf("got %d results, want 3 (victim* + elem + pos), results=%+v", len(span.Results), span.Results)
	}
	if len(span.Results[0].StarSegments) != 1 || span.Results[0].StarSegments[0].BaseTable == nil {
		t.Fatalf("result[0] = %+v, want victim's base star", span.Results[0])
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"victim.secret_tokens"}) {
		t.Errorf("elem sources = %v, want [victim.secret_tokens]", got)
	}
	if !strings.EqualFold(span.Results[2].Name, "pos") || len(span.Results[2].SourceColumns) != 0 {
		t.Errorf("result[2] = %+v, want offset column pos with NO lineage", span.Results[2])
	}
}

// TestProjection_UnnestImplicitAliasAndLiteral: an unaliased UNNEST over a
// column path takes the path's last component as the implicit element name
// (the GoogleSQL implicit-alias rule), and an UNNEST over a literal array has
// a RESOLVED-empty lineage (not an unresolved name to be matched fail-open).
func TestProjection_UnnestImplicitAliasAndLiteral(t *testing.T) {
	span, err := GetQuerySpan("SELECT secret_tokens FROM victim, UNNEST(victim.secret_tokens)", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"victim.secret_tokens"}) {
		t.Errorf("implicit-alias element sources = %v, want [victim.secret_tokens]", got)
	}

	span2, err := GetQuerySpan("SELECT e FROM UNNEST([1, 2, 3]) AS e", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span2.Results) != 1 || len(span2.Results[0].SourceColumns) != 0 {
		t.Errorf("literal-array element = %+v, want one column with empty (resolved) lineage", span2.Results)
	}
	if len(span2.AccessTables) != 0 {
		t.Errorf("tables = %v, want [] (a literal array reads no table)", tableNames(span2))
	}
}

// TestProjection_CorrelatedArrayPathRelation: the implicit-UNNEST comma form
// (`FROM t, t.arr AS a`) builds the same value relation, and a struct-field
// reference through the element alias (`a.child`) resolves to the element's
// array lineage.
func TestProjection_CorrelatedArrayPathRelation(t *testing.T) {
	span, err := GetQuerySpan("SELECT a FROM t, t.arr AS a", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t.arr"}) {
		t.Errorf("element sources = %v, want [t.arr]", got)
	}

	span2, err := GetQuerySpan("SELECT elem.child FROM victim, UNNEST(victim.kids) AS elem", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span2.Results[0]); !eqStrings(got, []string{"victim.kids"}) {
		t.Errorf("value-field sources = %v, want [victim.kids] (field access reads the element)", got)
	}
}

// ---------------------------------------------------------------------------
// Fix F3: BY NAME / CORRESPONDING set operations merge arms by column NAME.
//
// `UNION ALL BY NAME` matches arm columns by (case-insensitive) name, not by
// ordinal. Merging by ordinal attributed output `id` to the right arm's
// same-POSITION column (a secret) and vice versa — under-attribution both
// ways. Legacy could not parse BY NAME at all (fail-closed); omni must merge
// by name (or defer to the metadata-aware consumer when an arm carries an
// un-enumerable star).
// ---------------------------------------------------------------------------

// TestProjection_SetOpByNameConcrete: both arms concrete → inline name merge;
// output order = left arm's order; right-only names append (FULL BY NAME).
func TestProjection_SetOpByNameConcrete(t *testing.T) {
	sql := "WITH l AS (SELECT a AS id, b AS sec FROM t1), r AS (SELECT c AS sec, d AS id FROM t2) " +
		"SELECT * FROM l UNION ALL BY NAME SELECT * FROM r"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"id", "sec"}) {
		t.Fatalf("names = %v, want [id sec] (left arm order)", got)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t1.a", "t2.d"}) {
		t.Errorf("id sources = %v, want [t1.a t2.d] (NAME-matched, not ordinal)", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"t1.b", "t2.c"}) {
		t.Errorf("sec sources = %v, want [t1.b t2.c]", got)
	}
}

// TestProjection_SetOpByNameRightOnlyAppends: BY NAME appends right-only names
// after the left arm's columns (over-attribution-safe: a trailing extra never
// shifts earlier positions; with mismatched name sets BigQuery errors at
// execution anyway, so the extra entry is conservative). The FULL/LEFT prefix
// spelling (`FULL UNION ALL BY NAME`) is NOT parsed by the omni parser — it
// fail-closes at parse time — so only the bare BY NAME form reaches here.
func TestProjection_SetOpByNameRightOnlyAppends(t *testing.T) {
	sql := "WITH l AS (SELECT a AS id FROM t1), r AS (SELECT d AS id, e AS extra FROM t2) " +
		"SELECT * FROM l UNION ALL BY NAME SELECT * FROM r"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"id", "extra"}) {
		t.Fatalf("names = %v, want [id extra]", got)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t1.a", "t2.d"}) {
		t.Errorf("id sources = %v, want [t1.a t2.d]", got)
	}
	if got := resolvedSources(span.Results[1]); !eqStrings(got, []string{"t2.e"}) {
		t.Errorf("extra sources = %v, want [t2.e]", got)
	}
}

// TestProjection_SetOpByNameMatchColumns: `BY NAME ON (cols)` restricts the
// output to the listed columns, in list order (each unioning both arms'
// same-named lineage).
func TestProjection_SetOpByNameMatchColumns(t *testing.T) {
	sql := "WITH l AS (SELECT a AS id, b AS sec FROM t1), r AS (SELECT c AS sec, d AS id FROM t2) " +
		"SELECT * FROM l UNION ALL BY NAME ON (sec) SELECT * FROM r"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resultNames(span); !eqStrings(got, []string{"sec"}) {
		t.Fatalf("names = %v, want [sec]", got)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t1.b", "t2.c"}) {
		t.Errorf("sec sources = %v, want [t1.b t2.c]", got)
	}
}

// TestProjection_SetOpCorrespondingMergesByName: CORRESPONDING is the synonym
// family of BY NAME — same name-merge semantics.
func TestProjection_SetOpCorrespondingMergesByName(t *testing.T) {
	sql := "WITH l AS (SELECT a AS id, b AS sec FROM t1), r AS (SELECT c AS sec, d AS id FROM t2) " +
		"SELECT * FROM l UNION ALL CORRESPONDING SELECT * FROM r"
	span, err := GetQuerySpan(sql, DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if got := resolvedSources(span.Results[0]); !eqStrings(got, []string{"t1.a", "t2.d"}) {
		t.Errorf("id sources = %v, want [t1.a t2.d]", got)
	}
}

// TestProjection_SetOpByNameDeferredStarArm: an arm carrying an un-enumerable
// base-table star cannot be name-merged metadata-free — the whole merge defers
// to the consumer with the ByName flag (NEVER a silent ordinal merge).
func TestProjection_SetOpByNameDeferredStarArm(t *testing.T) {
	span, err := GetQuerySpan("SELECT id, sec FROM lt UNION ALL BY NAME SELECT * FROM rt", DialectBigQuery)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(span.Results) != 1 || span.Results[0].SetOpMerge == nil {
		t.Fatalf("results = %+v, want a single deferred SetOpMerge item", span.Results)
	}
	m := span.Results[0].SetOpMerge
	if !m.ByName {
		t.Errorf("SetOpMerge.ByName = false, want true (the consumer must name-merge after expansion)")
	}
	if len(m.Left) != 2 || len(m.Right) != 1 {
		t.Errorf("arm widths = %d/%d, want 2 concrete / 1 star", len(m.Left), len(m.Right))
	}
}

// TestProjection_StructFieldStarFailsClosed pins the gate-round P0: a STRUCT-
// field star through a relation (`d.s.*` / `t.s.*`) cannot be enumerated
// metadata-free — returning the relation's whole projection would misalign the
// positional masker (the first result's masker lands on the struct's first
// sub-column → a sensitive column is returned unmasked). It must FAIL CLOSED
// (the legacy resolver errored on multi-part wild paths too). A plain `rel.*`
// and a schema-qualified `schema.table.*` still resolve.
func TestProjection_StructFieldStarFailsClosed(t *testing.T) {
	for _, sql := range []string{
		"SELECT d.s.* FROM (SELECT public, STRUCT(secret AS ssn) AS s FROM t) AS d",
		"SELECT t.s.* FROM t",
	} {
		if _, err := GetQuerySpan(sql, DialectSpanner); err == nil {
			t.Errorf("GetQuerySpan(%q) = nil error, want fail-closed (struct-field star is not enumerable)", sql)
		}
	}
	// Single-part rel.* still resolves (no error).
	if _, err := GetQuerySpan("SELECT d.* FROM (SELECT a, b FROM t) AS d", DialectSpanner); err != nil {
		t.Errorf("rel.* should still resolve, got error: %v", err)
	}
	// Schema-qualified table star still resolves (the head is not a relation).
	if _, err := GetQuerySpan("SELECT analytics.events.* FROM analytics.events", DialectSpanner); err != nil {
		t.Errorf("schema-qualified star should still resolve, got error: %v", err)
	}
}
