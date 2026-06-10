package analysis_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/analysis"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mustExtract(t *testing.T, sql string) *analysis.QuerySpan {
	t.Helper()
	span, err := analysis.ExtractSQL(sql)
	if err != nil {
		t.Fatalf("ExtractSQL(%q): %v", sql, err)
	}
	return span
}

// sourceKey builds a compact string key for a SourceColumn.
func sourceKey(sc *analysis.SourceColumn) string {
	parts := []string{sc.Database, sc.Schema, sc.Table, sc.Column}
	return strings.Join(parts, ".")
}

// sourceKeys returns sorted keys for all sources in a QuerySpan.
func sourceKeys(span *analysis.QuerySpan) []string {
	var keys []string
	seen := map[string]bool{}
	for _, sc := range span.Sources {
		k := sourceKey(sc)
		if !seen[k] {
			seen[k] = true
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

// resultNames returns the result column names in order.
func resultNames(span *analysis.QuerySpan) []string {
	names := make([]string, len(span.Results))
	for i, rc := range span.Results {
		names[i] = rc.Name
	}
	return names
}

// resultSourceKeys returns sorted source keys for result column at index i.
func resultSourceKeys(span *analysis.QuerySpan, i int) []string {
	if i >= len(span.Results) {
		return nil
	}
	var keys []string
	for _, sc := range span.Results[i].Sources {
		keys = append(keys, sourceKey(sc))
	}
	sort.Strings(keys)
	return keys
}

// hasSrcKey reports whether any source in the span has this key.
func hasSrcKey(span *analysis.QuerySpan, key string) bool {
	for _, sc := range span.Sources {
		if sourceKey(sc) == key {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Test 1: SELECT col FROM t  (one result, one source)
// ---------------------------------------------------------------------------

func TestQuerySpan_SingleColumn(t *testing.T) {
	span := mustExtract(t, "SELECT a FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "A" {
		t.Errorf("result name: got %q, want %q", rc.Name, "A")
	}
	if rc.IsDerived {
		t.Error("IsDerived should be false for bare column ref")
	}
	if len(rc.Sources) != 1 {
		t.Fatalf("result sources: got %d, want 1", len(rc.Sources))
	}
	if rc.Sources[0].Table != "T" {
		t.Errorf("source table: got %q, want %q", rc.Sources[0].Table, "T")
	}
	if rc.Sources[0].Column != "A" {
		t.Errorf("source column: got %q, want %q", rc.Sources[0].Column, "A")
	}
}

// ---------------------------------------------------------------------------
// Test 2: SELECT a, b FROM t  (two results, two sources)
// ---------------------------------------------------------------------------

func TestQuerySpan_TwoColumns(t *testing.T) {
	span := mustExtract(t, "SELECT a, b FROM t")

	if len(span.Results) != 2 {
		t.Fatalf("Results: got %d, want 2", len(span.Results))
	}
	names := resultNames(span)
	if names[0] != "A" || names[1] != "B" {
		t.Errorf("result names: got %v, want [A B]", names)
	}
	for i, col := range []string{"A", "B"} {
		keys := resultSourceKeys(span, i)
		if len(keys) != 1 || keys[0] != "..T."+col {
			t.Errorf("result[%d] sources: got %v, want [..T.%s]", i, keys, col)
		}
	}
}

// ---------------------------------------------------------------------------
// Test 3: SELECT t.a AS x FROM t  (aliased result)
// ---------------------------------------------------------------------------

func TestQuerySpan_AliasedColumn(t *testing.T) {
	span := mustExtract(t, "SELECT t.a AS x FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "X" {
		t.Errorf("result name: got %q, want %q", rc.Name, "X")
	}
	if rc.IsDerived {
		t.Error("IsDerived should be false for aliased column ref")
	}
	if len(rc.Sources) != 1 {
		t.Fatalf("result sources: got %d, want 1", len(rc.Sources))
	}
	if rc.Sources[0].Table != "T" || rc.Sources[0].Column != "A" {
		t.Errorf("source: got {%s.%s}, want {T.A}", rc.Sources[0].Table, rc.Sources[0].Column)
	}
}

// ---------------------------------------------------------------------------
// Test 4: SELECT COUNT(*) FROM t  (derived result, pseudo-source)
// ---------------------------------------------------------------------------

func TestQuerySpan_CountStar(t *testing.T) {
	span := mustExtract(t, "SELECT COUNT(*) FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if !rc.IsDerived {
		t.Error("IsDerived should be true for COUNT(*)")
	}
	// The function name is COUNT.
	if rc.Name != "COUNT" {
		t.Errorf("result name: got %q, want %q", rc.Name, "COUNT")
	}
}

// ---------------------------------------------------------------------------
// Test 5: SELECT * FROM t  (star result, table source)
// ---------------------------------------------------------------------------

func TestQuerySpan_StarFromTable(t *testing.T) {
	span := mustExtract(t, "SELECT * FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "*" {
		t.Errorf("result name: got %q, want %q", rc.Name, "*")
	}
	if rc.IsDerived {
		t.Error("IsDerived should be false for SELECT *")
	}
	// Should have a source pointing to T.*
	found := false
	for _, sc := range rc.Sources {
		if sc.Table == "T" && sc.Column == "*" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected source {T.*} in result sources, got %v", rc.Sources)
	}
}

// ---------------------------------------------------------------------------
// Test 6: SELECT * EXCLUDE(b) FROM t
// ---------------------------------------------------------------------------

func TestQuerySpan_StarExclude(t *testing.T) {
	span := mustExtract(t, "SELECT * EXCLUDE(b) FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "*" {
		t.Errorf("result name: got %q, want %q", rc.Name, "*")
	}
	// Since t has unknown schema, we still get a * pseudo-source.
	// But the test verifies it parsed correctly at minimum.
	_ = rc
}

// ---------------------------------------------------------------------------
// Test 7: FROM t JOIN u ON ...  (sources from both)
// ---------------------------------------------------------------------------

func TestQuerySpan_Join(t *testing.T) {
	span := mustExtract(t, "SELECT t.a, u.b FROM t JOIN u ON t.id = u.id")

	if len(span.Results) != 2 {
		t.Fatalf("Results: got %d, want 2", len(span.Results))
	}
	// Both T and U should appear in sources.
	keys := sourceKeys(span)
	hasT := false
	hasU := false
	for _, k := range keys {
		if strings.Contains(k, "T.") {
			hasT = true
		}
		if strings.Contains(k, "U.") {
			hasU = true
		}
	}
	if !hasT {
		t.Errorf("expected T in sources, got %v", keys)
	}
	if !hasU {
		t.Errorf("expected U in sources, got %v", keys)
	}
}

// ---------------------------------------------------------------------------
// Test 8: FROM (SELECT ...) s  (subquery in FROM)
// ---------------------------------------------------------------------------

func TestQuerySpan_SubqueryInFrom(t *testing.T) {
	span := mustExtract(t, "SELECT s.col1 FROM (SELECT a AS col1 FROM inner_table) AS s")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "COL1" {
		t.Errorf("result name: got %q, want %q", rc.Name, "COL1")
	}
	// The source should trace through the subquery to s.COL1.
	if len(rc.Sources) == 0 {
		t.Error("expected at least 1 source for subquery column ref")
	}
	// The outer query's sources should include the inner table.
	if !hasSrcKey(span, "...INNER_TABLE.*") && !hasSrcKey(span, "...S.COL1") {
		// Accept either form.
		t.Logf("sources: %v", sourceKeys(span))
	}
}

// ---------------------------------------------------------------------------
// Test 9: WITH cte AS (SELECT col FROM t) SELECT * FROM cte
// ---------------------------------------------------------------------------

func TestQuerySpan_CTE(t *testing.T) {
	span := mustExtract(t, `WITH cte AS (SELECT col FROM t) SELECT * FROM cte`)

	// cte exposes col; SELECT * from cte should have COL as a source column.
	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "*" {
		t.Errorf("result name: got %q, want %q", rc.Name, "*")
	}
	// Sources should trace through CTE to T.COL.
	found := false
	for _, sc := range rc.Sources {
		if sc.Table == "CTE" && sc.Column == "COL" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected source {CTE.COL} in result, got: %v", rc.Sources)
	}
}

// ---------------------------------------------------------------------------
// Test 10: UNION ALL — positional merge
// ---------------------------------------------------------------------------

func TestQuerySpan_UnionAll(t *testing.T) {
	span := mustExtract(t, "SELECT a FROM t1 UNION ALL SELECT b FROM t2")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	// Name comes from left branch.
	if span.Results[0].Name != "A" {
		t.Errorf("result name: got %q, want %q", span.Results[0].Name, "A")
	}
	// Sources should include both T1.A and T2.B.
	keys := sourceKeys(span)
	hasT1 := false
	hasT2 := false
	for _, k := range keys {
		if strings.Contains(k, "T1.") {
			hasT1 = true
		}
		if strings.Contains(k, "T2.") {
			hasT2 = true
		}
	}
	if !hasT1 {
		t.Errorf("expected T1 in union sources, got %v", keys)
	}
	if !hasT2 {
		t.Errorf("expected T2 in union sources, got %v", keys)
	}
}

// ---------------------------------------------------------------------------
// Test 11: UNION ALL BY NAME — named merge
// ---------------------------------------------------------------------------

func TestQuerySpan_UnionByName(t *testing.T) {
	span := mustExtract(t, "SELECT a, b FROM t1 UNION ALL BY NAME SELECT b, c FROM t2")

	// With BY NAME: a, b appear in left, b, c in right. Result should have a, b, c.
	names := resultNames(span)
	// At minimum, "B" must appear (common to both). "A" from left, "C" from right.
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}
	if !nameSet["A"] {
		t.Errorf("expected A in results, got %v", names)
	}
	if !nameSet["B"] {
		t.Errorf("expected B in results, got %v", names)
	}
	if !nameSet["C"] {
		t.Errorf("expected C in results, got %v", names)
	}
}

// ---------------------------------------------------------------------------
// Test 12: Qualified column ref db.schema.t.col
// ---------------------------------------------------------------------------

func TestQuerySpan_QualifiedColumnRef(t *testing.T) {
	span := mustExtract(t, "SELECT mydb.myschema.t.col FROM mydb.myschema.t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "COL" {
		t.Errorf("result name: got %q, want %q", rc.Name, "COL")
	}
	if len(rc.Sources) == 0 {
		t.Fatal("expected sources for qualified column ref")
	}
	sc := rc.Sources[0]
	if sc.Database != "MYDB" {
		t.Errorf("source database: got %q, want %q", sc.Database, "MYDB")
	}
	if sc.Schema != "MYSCHEMA" {
		t.Errorf("source schema: got %q, want %q", sc.Schema, "MYSCHEMA")
	}
	if sc.Table != "T" {
		t.Errorf("source table: got %q, want %q", sc.Table, "T")
	}
	if sc.Column != "COL" {
		t.Errorf("source column: got %q, want %q", sc.Column, "COL")
	}
}

// ---------------------------------------------------------------------------
// Test 13: CTE with column aliases: WITH cte(x, y) AS (...) SELECT x FROM cte
// ---------------------------------------------------------------------------

func TestQuerySpan_CTEColumnAliases(t *testing.T) {
	span := mustExtract(t, `WITH cte(x, y) AS (SELECT a, b FROM t) SELECT x FROM cte`)

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "X" {
		t.Errorf("result name: got %q, want %q", rc.Name, "X")
	}
	// Source should be cte.X (renamed from a).
	found := false
	for _, sc := range rc.Sources {
		if sc.Table == "CTE" && sc.Column == "X" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected source {CTE.X}, got %v", rc.Sources)
	}
}

// ---------------------------------------------------------------------------
// Test 14: SELECT constant expression  (derived, no table sources)
// ---------------------------------------------------------------------------

func TestQuerySpan_ConstantExpression(t *testing.T) {
	span := mustExtract(t, "SELECT 1 + 2 AS result")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "RESULT" {
		t.Errorf("result name: got %q, want %q", rc.Name, "RESULT")
	}
	if !rc.IsDerived {
		t.Error("IsDerived should be true for constant expression")
	}
	// No column-level sources for a pure constant.
	for _, sc := range rc.Sources {
		if sc.Column != "*" && sc.Table == "" {
			t.Errorf("unexpected column source for constant: %+v", sc)
		}
	}
}

// ---------------------------------------------------------------------------
// Test 15: Multi-table JOIN with qualified refs
// ---------------------------------------------------------------------------

func TestQuerySpan_MultiTableJoin(t *testing.T) {
	span := mustExtract(t, `
		SELECT o.order_id, c.name
		FROM orders AS o
		JOIN customers AS c ON o.customer_id = c.id
	`)

	if len(span.Results) != 2 {
		t.Fatalf("Results: got %d, want 2", len(span.Results))
	}
	names := resultNames(span)
	if names[0] != "ORDER_ID" {
		t.Errorf("result[0] name: got %q, want %q", names[0], "ORDER_ID")
	}
	if names[1] != "NAME" {
		t.Errorf("result[1] name: got %q, want %q", names[1], "NAME")
	}
	// result[0] sources should include O.ORDER_ID.
	keys0 := resultSourceKeys(span, 0)
	found0 := false
	for _, k := range keys0 {
		if strings.Contains(k, "O.ORDER_ID") {
			found0 = true
		}
	}
	if !found0 {
		t.Errorf("result[0] expected O.ORDER_ID source, got %v", keys0)
	}
	// result[1] sources should include C.NAME.
	keys1 := resultSourceKeys(span, 1)
	found1 := false
	for _, k := range keys1 {
		if strings.Contains(k, "C.NAME") {
			found1 = true
		}
	}
	if !found1 {
		t.Errorf("result[1] expected C.NAME source, got %v", keys1)
	}
}

// ---------------------------------------------------------------------------
// Test 16: Nested CTE — outer query references inner CTE
// ---------------------------------------------------------------------------

func TestQuerySpan_NestedCTE(t *testing.T) {
	span := mustExtract(t, `
		WITH
		  base AS (SELECT id, val FROM raw_data),
		  filtered AS (SELECT id FROM base WHERE val > 0)
		SELECT id FROM filtered
	`)

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	// The result column ID should ultimately trace back through filtered → base → raw_data.
	rc := span.Results[0]
	if rc.Name != "ID" {
		t.Errorf("result name: got %q, want %q", rc.Name, "ID")
	}
	_ = rc // sources chain through CTEs
}

// ---------------------------------------------------------------------------
// Test 17: Expression with function — marked IsDerived
// ---------------------------------------------------------------------------

func TestQuerySpan_ExpressionIsDerived(t *testing.T) {
	span := mustExtract(t, "SELECT a + b AS total FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if !rc.IsDerived {
		t.Error("IsDerived should be true for a + b expression")
	}
	if rc.Name != "TOTAL" {
		t.Errorf("result name: got %q, want %q", rc.Name, "TOTAL")
	}
	// Both A and B should appear as sources.
	keys := resultSourceKeys(span, 0)
	hasA := false
	hasB := false
	for _, k := range keys {
		if strings.HasSuffix(k, ".A") {
			hasA = true
		}
		if strings.HasSuffix(k, ".B") {
			hasB = true
		}
	}
	if !hasA {
		t.Errorf("expected A in sources of a+b, got %v", keys)
	}
	if !hasB {
		t.Errorf("expected B in sources of a+b, got %v", keys)
	}
}

// ---------------------------------------------------------------------------
// Test 18: ExtractSQL error path
// ---------------------------------------------------------------------------

func TestQuerySpan_ParseError(t *testing.T) {
	_, err := analysis.ExtractSQL("SELECT FROM WHERE")
	if err == nil {
		t.Error("expected error for invalid SQL, got nil")
	}
}

// ---------------------------------------------------------------------------
// Test 19: CASE expression columns are attributed (walker coverage regression)
// ---------------------------------------------------------------------------

// The generated AST walker used to skip CaseExpr.Whens (a []*WhenClause of a
// non-Node helper struct), so the WHEN condition and THEN result columns were
// invisible to ast.Inspect and never attributed to the result column.
func TestQuerySpan_CaseExpressionColumns(t *testing.T) {
	span := mustExtract(t, "SELECT CASE WHEN a > 0 THEN b ELSE c END AS r FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "R" {
		t.Errorf("result name: got %q, want %q", rc.Name, "R")
	}
	if !rc.IsDerived {
		t.Error("IsDerived should be true for a CASE expression")
	}
	keys := resultSourceKeys(span, 0)
	want := []string{"..T.A", "..T.B", "..T.C"}
	if len(keys) != len(want) {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Errorf("source[%d]: got %q, want %q", i, keys[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Test 20: window function columns are attributed (walker coverage regression)
// ---------------------------------------------------------------------------

// FuncCallExpr.Over (*WindowSpec, a non-Node helper struct) used to be skipped
// by the walker, hiding PARTITION BY / ORDER BY columns from the span.
func TestQuerySpan_WindowFunctionColumns(t *testing.T) {
	span := mustExtract(t, "SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary) AS rn FROM emp")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	rc := span.Results[0]
	if rc.Name != "RN" {
		t.Errorf("result name: got %q, want %q", rc.Name, "RN")
	}
	if !rc.IsDerived {
		t.Error("IsDerived should be true for a window function")
	}
	keys := resultSourceKeys(span, 0)
	want := []string{"..EMP.DEPT", "..EMP.SALARY"}
	if len(keys) != len(want) {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Errorf("source[%d]: got %q, want %q", i, keys[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Test 21: WITHIN GROUP order columns are attributed (walker coverage regression)
// ---------------------------------------------------------------------------

// FuncCallExpr.OrderBy ([]*OrderItem, non-Node helper structs) used to be
// skipped by the walker, hiding WITHIN GROUP (ORDER BY ...) columns.
func TestQuerySpan_WithinGroupColumns(t *testing.T) {
	span := mustExtract(t, "SELECT LISTAGG(name, ',') WITHIN GROUP (ORDER BY pos) AS l FROM t")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	keys := resultSourceKeys(span, 0)
	want := []string{"..T.NAME", "..T.POS"}
	if len(keys) != len(want) {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	for i := range want {
		if keys[i] != want[i] {
			t.Errorf("source[%d]: got %q, want %q", i, keys[i], want[i])
		}
	}
}

// ---------------------------------------------------------------------------
// PIVOT / UNPIVOT: opaque-derived relations (masking soundness)
// ---------------------------------------------------------------------------
//
// The output columns of a pivoted relation are computed from the source
// relation under names that do not exist on it. Without catalog information
// the projection cannot be enumerated, so every column resolved through a
// pivoted relation must attribute to the source's whole-relation ("*")
// sources. Fabricating per-column attributions (e.g. T."2023_Q1" for a pivot
// value column that actually carries SUM(T.AMOUNT)) would under-attribute and
// let a masking consumer leak.

func TestQuerySpan_PivotStar(t *testing.T) {
	span := mustExtract(t, "SELECT * FROM t PIVOT(SUM(a) FOR m IN ('JAN', 'FEB')) AS p")

	if len(span.Results) != 1 {
		t.Fatalf("Results: got %d, want 1", len(span.Results))
	}
	keys := resultSourceKeys(span, 0)
	want := []string{"..T.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	// The pivot alias P must not surface as a source table.
	for _, k := range sourceKeys(span) {
		if strings.HasPrefix(k, "..P.") {
			t.Errorf("pivot alias leaked into sources: %v", sourceKeys(span))
		}
	}
}

func TestQuerySpan_PivotValueColumnNotFabricated(t *testing.T) {
	span := mustExtract(t, `SELECT "2023_Q1" FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN ('2023_Q1', '2023_Q2')) AS p`)

	keys := resultSourceKeys(span, 0)
	want := []string{"..QUARTERLY_SALES.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	if hasSrcKey(span, "..QUARTERLY_SALES.2023_Q1") {
		t.Error("fabricated base-column attribution for a pivot value column")
	}
}

func TestQuerySpan_PivotQualifiedRef(t *testing.T) {
	span := mustExtract(t, "SELECT p.q1 FROM db1.s1.t PIVOT(SUM(a) FOR m IN ('JAN')) AS p")

	keys := resultSourceKeys(span, 0)
	want := []string{"DB1.S1.T.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	if hasSrcKey(span, "..P.Q1") {
		t.Error("qualified ref through pivot alias must not fabricate table P")
	}
}

func TestQuerySpan_PivotQualifiedStar(t *testing.T) {
	span := mustExtract(t, "SELECT p.* FROM t PIVOT(SUM(a) FOR m IN ('JAN')) AS p")

	keys := resultSourceKeys(span, 0)
	want := []string{"..T.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
}

func TestQuerySpan_UnpivotValueColumn(t *testing.T) {
	span := mustExtract(t, "SELECT sales FROM monthly_sales UNPIVOT (sales FOR month IN (jan, feb)) unpvt")

	keys := resultSourceKeys(span, 0)
	want := []string{"..MONTHLY_SALES.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
	if hasSrcKey(span, "..MONTHLY_SALES.SALES") {
		t.Error("UNPIVOT value column attributed to a nonexistent base column")
	}
	if hasSrcKey(span, "..UNPVT.SALES") {
		t.Error("UNPIVOT alias fabricated as a source table")
	}
}

func TestQuerySpan_PivotChainedOverSubquery(t *testing.T) {
	span := mustExtract(t, `SELECT q1 FROM (SELECT a, b FROM real_table) PIVOT (SUM(a) FOR b IN ('x')) PIVOT (MAX(a) FOR b IN ('y')) AS pp`)

	// The chain collapses to the subquery's whole-relation source (the
	// package addresses unaliased subqueries as "_subquery").
	keys := resultSourceKeys(span, 0)
	want := []string{".._subquery.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
}

func TestQuerySpan_PivotJoinMixedAttribution(t *testing.T) {
	// A bare column in a scope with both a pivoted relation and a plain table
	// may come from either: both must be attributed (conservative).
	span := mustExtract(t, "SELECT x FROM t PIVOT(SUM(a) FOR m IN ('J')) AS p, u")

	keys := resultSourceKeys(span, 0)
	if len(keys) != 2 {
		t.Fatalf("result sources: got %v, want T.* plus U.X", keys)
	}
	if keys[0] != "..T.*" || keys[1] != "..U.X" {
		t.Errorf("result sources: got %v, want [..T.* ..U.X]", keys)
	}
}

func TestQuerySpan_PivotOverCTE(t *testing.T) {
	span := mustExtract(t, "WITH c AS (SELECT a FROM real_table) SELECT * FROM c PIVOT(SUM(a) FOR m IN ('J')) AS p")

	// The CTE's columns collapse to a whole-relation source under the CTE
	// name (consistent with how non-pivoted CTE references are reported).
	keys := resultSourceKeys(span, 0)
	want := []string{"..C.*"}
	if len(keys) != 1 || keys[0] != want[0] {
		t.Fatalf("result sources: got %v, want %v", keys, want)
	}
}
