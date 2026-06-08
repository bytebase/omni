package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-query-clauses node: PIVOT / UNPIVOT / TABLESAMPLE
// table-source operators, the AT SYSTEM TIME (FOR SYSTEM_TIME AS OF) suffix, and
// the SELECT-level differential-privacy clause. Accept/reject is also proven
// against the live Cloud Spanner emulator in pivot_unpivot_oracle_test.go (build
// tag googlesql_oracle) for the oracle-authoritative forms; these hand-written
// tests assert the AST STRUCTURE and cover the BigQuery-only differential-
// privacy clause that is non-authoritative against that emulator.

// selFrom0 parses one SELECT statement and returns its first FROM item.
func selFrom0(t *testing.T, sql string) ast.Node {
	t.Helper()
	q := parseQ(t, sql)
	sel := selectOf(t, q)
	if len(sel.From) == 0 {
		t.Fatalf("Parse(%q): SELECT has no FROM items", sql)
	}
	return sel.From[0]
}

// tableExpr0 parses one SELECT and returns its first FROM item as a *TableExpr.
func tableExpr0(t *testing.T, sql string) *ast.TableExpr {
	t.Helper()
	te, ok := selFrom0(t, sql).(*ast.TableExpr)
	if !ok {
		t.Fatalf("Parse(%q): FROM[0] = %T, want *ast.TableExpr", sql, selFrom0(t, sql))
	}
	return te
}

// --- PIVOT structure ---------------------------------------------------------

func TestPivot_Structure(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM Produce PIVOT(SUM(sales) FOR quarter IN ('Q1', 'Q2', 'Q3', 'Q4'))")
	if te.Pivot == nil {
		t.Fatalf("PIVOT not parsed: te.Pivot is nil")
	}
	if te.Unpivot != nil || te.Sample != nil {
		t.Errorf("only Pivot should be set; got Unpivot=%v Sample=%v", te.Unpivot, te.Sample)
	}
	if got := len(te.Pivot.Aggregates); got != 1 {
		t.Errorf("aggregates: got %d, want 1", got)
	}
	if te.Pivot.For == nil {
		t.Errorf("FOR expression is nil")
	}
	if got := len(te.Pivot.Values); got != 4 {
		t.Errorf("IN values: got %d, want 4", got)
	}
	if te.Pivot.Alias != "" {
		t.Errorf("alias: got %q, want empty", te.Pivot.Alias)
	}
}

func TestPivot_MultiAggregateWithAliases(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM (SELECT product, sales, quarter FROM Produce) PIVOT(SUM(sales) AS total_sales, COUNT(*) AS num_records FOR quarter IN ('Q1', 'Q2'))")
	if te.Pivot == nil {
		t.Fatalf("PIVOT not parsed")
	}
	if got := len(te.Pivot.Aggregates); got != 2 {
		t.Fatalf("aggregates: got %d, want 2", got)
	}
	if te.Pivot.Aggregates[0].Alias != "total_sales" {
		t.Errorf("agg[0] alias: got %q, want total_sales", te.Pivot.Aggregates[0].Alias)
	}
	if te.Pivot.Aggregates[1].Alias != "num_records" {
		t.Errorf("agg[1] alias: got %q, want num_records", te.Pivot.Aggregates[1].Alias)
	}
	// The source is a parenthesized subquery.
	if te.Subquery == nil {
		t.Errorf("expected a subquery source under PIVOT")
	}
}

func TestPivot_ValueAliases(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1' AS first, 'Q2'))")
	if te.Pivot == nil {
		t.Fatalf("PIVOT not parsed")
	}
	if got := len(te.Pivot.Values); got != 2 {
		t.Fatalf("values: got %d, want 2", got)
	}
	if te.Pivot.Values[0].Alias != "first" {
		t.Errorf("value[0] alias: got %q, want first", te.Pivot.Values[0].Alias)
	}
	if te.Pivot.Values[1].Alias != "" {
		t.Errorf("value[1] alias: got %q, want empty", te.Pivot.Values[1].Alias)
	}
}

func TestPivot_TableAndPivotAlias(t *testing.T) {
	// `t AS x PIVOT(...) AS p`: x is the table alias, p is the pivot alias.
	te := tableExpr0(t, "SELECT * FROM t AS x PIVOT(SUM(s) FOR q IN ('a')) AS p")
	if te.Alias != "x" {
		t.Errorf("table alias: got %q, want x", te.Alias)
	}
	if te.Pivot == nil {
		t.Fatalf("PIVOT not parsed")
	}
	if te.Pivot.Alias != "p" {
		t.Errorf("pivot alias: got %q, want p", te.Pivot.Alias)
	}
}

func TestPivot_BarePivotAlias(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t PIVOT(SUM(s) FOR q IN ('a')) p")
	if te.Alias != "" {
		t.Errorf("table alias: got %q, want empty", te.Alias)
	}
	if te.Pivot == nil || te.Pivot.Alias != "p" {
		t.Errorf("pivot alias: want p, got %+v", te.Pivot)
	}
}

// --- UNPIVOT structure -------------------------------------------------------

func TestUnpivot_SingleColumn(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM sales_table UNPIVOT(total_sales FOR quarter IN (Q1, Q2, Q3, Q4))")
	if te.Unpivot == nil {
		t.Fatalf("UNPIVOT not parsed")
	}
	if te.Pivot != nil || te.Sample != nil {
		t.Errorf("only Unpivot should be set")
	}
	if got := len(te.Unpivot.ValueColumns); got != 1 {
		t.Errorf("value columns: got %d, want 1", got)
	}
	if te.Unpivot.NameColumn == nil {
		t.Errorf("name column is nil")
	}
	if got := len(te.Unpivot.Items); got != 4 {
		t.Errorf("IN items: got %d, want 4", got)
	}
	if te.Unpivot.NullsMode != ast.UnpivotNullsUnspecified {
		t.Errorf("nulls mode: got %v, want unspecified", te.Unpivot.NullsMode)
	}
}

func TestUnpivot_IncludeNulls(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t UNPIVOT INCLUDE NULLS (total_sales FOR quarter IN (Q1, Q2))")
	if te.Unpivot == nil {
		t.Fatalf("UNPIVOT not parsed")
	}
	if te.Unpivot.NullsMode != ast.UnpivotIncludeNulls {
		t.Errorf("nulls mode: got %v, want IncludeNulls", te.Unpivot.NullsMode)
	}
}

func TestUnpivot_ExcludeNulls(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t UNPIVOT EXCLUDE NULLS (s FOR q IN (Q1, Q2))")
	if te.Unpivot == nil || te.Unpivot.NullsMode != ast.UnpivotExcludeNulls {
		t.Errorf("want ExcludeNulls, got %+v", te.Unpivot)
	}
}

func TestUnpivot_MultiColumn(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t UNPIVOT ((s1, s2) FOR q IN ((Q1, Q2), (Q3, Q4)))")
	if te.Unpivot == nil {
		t.Fatalf("UNPIVOT not parsed")
	}
	if got := len(te.Unpivot.ValueColumns); got != 2 {
		t.Errorf("value columns: got %d, want 2 (multi-column)", got)
	}
	if got := len(te.Unpivot.Items); got != 2 {
		t.Fatalf("IN items: got %d, want 2", got)
	}
	if got := len(te.Unpivot.Items[0].Columns); got != 2 {
		t.Errorf("item[0] columns: got %d, want 2 (parenthesized group)", got)
	}
}

func TestUnpivot_InItemAliases(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t UNPIVOT (s FOR q IN (Q1 AS 'a', Q2 AS 2)) AS u")
	if te.Unpivot == nil {
		t.Fatalf("UNPIVOT not parsed")
	}
	if te.Unpivot.Alias != "u" {
		t.Errorf("unpivot alias: got %q, want u", te.Unpivot.Alias)
	}
	if got := len(te.Unpivot.Items); got != 2 {
		t.Fatalf("IN items: got %d, want 2", got)
	}
	it0 := te.Unpivot.Items[0]
	if !it0.HasAlias || it0.AliasIsInt || it0.AliasString != "a" {
		t.Errorf("item[0]: want string alias 'a', got %+v", it0)
	}
	it1 := te.Unpivot.Items[1]
	if !it1.HasAlias || !it1.AliasIsInt || it1.AliasInt != "2" {
		t.Errorf("item[1]: want int alias 2, got %+v", it1)
	}
}

// --- TABLESAMPLE structure ---------------------------------------------------

func TestTableSample_Percent(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM dataset.my_table TABLESAMPLE SYSTEM (10 PERCENT)")
	if te.Sample == nil {
		t.Fatalf("TABLESAMPLE not parsed")
	}
	if te.Pivot != nil || te.Unpivot != nil {
		t.Errorf("only Sample should be set")
	}
	if te.Sample.Method != "SYSTEM" {
		t.Errorf("method: got %q, want SYSTEM", te.Sample.Method)
	}
	if te.Sample.Unit != ast.SampleUnitPercent {
		t.Errorf("unit: got %v, want PERCENT", te.Sample.Unit)
	}
	if te.Sample.Size == nil {
		t.Errorf("size is nil")
	}
}

func TestTableSample_Rows(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM Albums TABLESAMPLE RESERVOIR (100 ROWS)")
	if te.Sample == nil {
		t.Fatalf("TABLESAMPLE not parsed")
	}
	if te.Sample.Method != "RESERVOIR" {
		t.Errorf("method: got %q, want RESERVOIR", te.Sample.Method)
	}
	if te.Sample.Unit != ast.SampleUnitRows {
		t.Errorf("unit: got %v, want ROWS", te.Sample.Unit)
	}
}

func TestTableSample_Bernoulli(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM Singers TABLESAMPLE BERNOULLI (10 PERCENT)")
	if te.Sample == nil || te.Sample.Method != "BERNOULLI" {
		t.Errorf("want BERNOULLI sample, got %+v", te.Sample)
	}
}

func TestTableSample_Repeatable(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t TABLESAMPLE BERNOULLI (5 PERCENT) REPEATABLE (42)")
	if te.Sample == nil {
		t.Fatalf("TABLESAMPLE not parsed")
	}
	if te.Sample.Repeatable == nil {
		t.Errorf("REPEATABLE seed not captured")
	}
	if te.Sample.WithWeight {
		t.Errorf("WithWeight should be false")
	}
}

func TestTableSample_WithWeight(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t TABLESAMPLE RESERVOIR (100 ROWS) WITH WEIGHT")
	if te.Sample == nil {
		t.Fatalf("TABLESAMPLE not parsed")
	}
	if !te.Sample.WithWeight {
		t.Errorf("WITH WEIGHT not captured")
	}
}

func TestTableSample_WithWeightAlias(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t TABLESAMPLE RESERVOIR (100 ROWS) WITH WEIGHT AS w")
	if te.Sample == nil || !te.Sample.WithWeight {
		t.Fatalf("WITH WEIGHT not parsed")
	}
	if te.Sample.WeightAlias != "w" {
		t.Errorf("weight alias: got %q, want w", te.Sample.WeightAlias)
	}
}

func TestTableSample_PartitionBy(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t TABLESAMPLE SYSTEM (10 PERCENT PARTITION BY x, y)")
	if te.Sample == nil {
		t.Fatalf("TABLESAMPLE not parsed")
	}
	if got := len(te.Sample.PartitionBy); got != 2 {
		t.Errorf("PARTITION BY exprs: got %d, want 2", got)
	}
}

func TestTableSample_WithWeightAliasRepeatable(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) WITH WEIGHT w REPEATABLE (1)")
	if te.Sample == nil || !te.Sample.WithWeight {
		t.Fatalf("WITH WEIGHT not parsed")
	}
	if te.Sample.WeightAlias != "w" {
		t.Errorf("weight alias: got %q, want w", te.Sample.WeightAlias)
	}
	if te.Sample.Repeatable == nil {
		t.Errorf("trailing REPEATABLE seed not captured")
	}
}

func TestTableSample_TableAliasBefore(t *testing.T) {
	// `t AS x TABLESAMPLE(...)`: x is the table alias, sample wraps it.
	te := tableExpr0(t, "SELECT * FROM t AS x TABLESAMPLE BERNOULLI (10 PERCENT)")
	if te.Alias != "x" {
		t.Errorf("table alias: got %q, want x", te.Alias)
	}
	if te.Sample == nil {
		t.Errorf("TABLESAMPLE not parsed")
	}
}

// --- combinations ------------------------------------------------------------

func TestPivot_ThenTableSample(t *testing.T) {
	// `t PIVOT(...) TABLESAMPLE(...)`: PIVOT binds to the source, then TABLESAMPLE
	// wraps the whole thing (oracle-accepted). Both fields set on the same source.
	te := tableExpr0(t, "SELECT * FROM t PIVOT(SUM(s) FOR q IN ('a')) TABLESAMPLE BERNOULLI (10 PERCENT)")
	if te.Pivot == nil {
		t.Errorf("PIVOT not parsed")
	}
	if te.Sample == nil {
		t.Errorf("TABLESAMPLE not parsed")
	}
}

func TestSubquery_TableSample(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM (SELECT 1 AS s, 2 AS q) TABLESAMPLE BERNOULLI (10 PERCENT)")
	if te.Subquery == nil {
		t.Errorf("expected subquery source")
	}
	if te.Sample == nil {
		t.Errorf("TABLESAMPLE not parsed on subquery")
	}
}

func TestSubquery_Pivot(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM (SELECT 1 AS s, 2 AS q) AS sub PIVOT(SUM(s) FOR q IN (2)) AS p")
	if te.Subquery == nil {
		t.Errorf("expected subquery source")
	}
	if te.Alias != "sub" {
		t.Errorf("subquery alias: got %q, want sub", te.Alias)
	}
	if te.Pivot == nil || te.Pivot.Alias != "p" {
		t.Errorf("pivot/alias not parsed: %+v", te.Pivot)
	}
}

// --- AT SYSTEM TIME (FOR SYSTEM_TIME AS OF) ----------------------------------

func TestForSystemTime_Captured(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t FOR SYSTEM_TIME AS OF '2017-01-01 10:00:00-07:00'")
	if te.SystemTime == nil {
		t.Errorf("FOR SYSTEM_TIME AS OF not captured")
	}
}

func TestForSystemTime_Expression(t *testing.T) {
	te := tableExpr0(t, "SELECT * FROM t FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR)")
	if te.SystemTime == nil {
		t.Errorf("FOR SYSTEM_TIME AS OF expression not captured")
	}
}

// --- differential-privacy SELECT clause (BigQuery-only) ----------------------

func TestSelectWith_DifferentialPrivacy(t *testing.T) {
	// DIFFERENTIAL_PRIVACY is over-reserved by the omni lexer; the select-with
	// clause must still parse it (parseSelectWith admits that one keyword). The
	// mechanism name is source-preserved (not folded).
	q := parseQ(t, "SELECT WITH DIFFERENTIAL_PRIVACY OPTIONS(epsilon=1, delta=1e-10) a FROM t")
	sel := selectOf(t, q)
	if !strings.EqualFold(sel.SelectWith, "differential_privacy") {
		t.Errorf("SelectWith: got %q, want differential_privacy (any case)", sel.SelectWith)
	}
	if sel.SelectWith == "" {
		t.Errorf("SelectWith is empty; clause not captured")
	}
	if len(sel.Items) != 1 {
		t.Errorf("select items: got %d, want 1", len(sel.Items))
	}
}

func TestSelectWith_Anonymization(t *testing.T) {
	q := parseQ(t, "SELECT WITH ANONYMIZATION OPTIONS(epsilon=1) a FROM t")
	sel := selectOf(t, q)
	if !strings.EqualFold(sel.SelectWith, "anonymization") {
		t.Errorf("SelectWith: got %q, want anonymization (any case)", sel.SelectWith)
	}
}

func TestSelectWith_AggregationThreshold(t *testing.T) {
	q := parseQ(t, "SELECT WITH AGGREGATION_THRESHOLD OPTIONS(threshold=10, privacy_unit_column=id) COUNT(*) FROM t")
	sel := selectOf(t, q)
	if !strings.EqualFold(sel.SelectWith, "aggregation_threshold") {
		t.Errorf("SelectWith: got %q, want aggregation_threshold (any case)", sel.SelectWith)
	}
}

func TestSelectWith_NoOptions(t *testing.T) {
	q := parseQ(t, "SELECT WITH DIFFERENTIAL_PRIVACY a FROM t")
	sel := selectOf(t, q)
	if !strings.EqualFold(sel.SelectWith, "differential_privacy") {
		t.Errorf("SelectWith: got %q, want differential_privacy (any case)", sel.SelectWith)
	}
}

// WITH (...) inline expression must NOT be read as a select-with clause.
func TestSelectWith_NotInlineWithExpr(t *testing.T) {
	// `SELECT WITH(...)` is an inline WITH expression column, not opt_select_with.
	// (The disambiguator is `WITH (` vs `WITH <identifier>`.) Just assert it does
	// not get mis-captured as a differential-privacy clause name.
	file, errs := Parse("SELECT WITH(x AS 1, x + 1) FROM t")
	if len(errs) == 0 {
		q := file.Stmts[0].(*ast.QueryStmt)
		sel := q.Body.(*ast.SelectStmt)
		if sel.SelectWith != "" {
			t.Errorf("inline WITH() mis-captured as select-with %q", sel.SelectWith)
		}
	}
}

// --- regression: Codex review findings ---------------------------------------

// PIVOT / UNPIVOT are non-reserved, so a bare `FROM t pivot` is an implicit table
// alias (oracle-confirmed accept), NOT a malformed operator. Regression for the
// atTableAliasStop over-reject (Codex finding).
func TestPivot_NonReservedImplicitAlias(t *testing.T) {
	// Simple cases where FROM[0] is the aliased TableExpr itself.
	for _, sql := range []string{
		"SELECT * FROM t pivot",
		"SELECT * FROM t unpivot",
		"SELECT * FROM t AS pivot",
		"SELECT * FROM t pivot WHERE x > 1",
	} {
		te := tableExpr0(t, sql)
		if te.Pivot != nil || te.Unpivot != nil {
			t.Errorf("%q: bare keyword wrongly parsed as operator (Pivot=%v Unpivot=%v)", sql, te.Pivot, te.Unpivot)
		}
		if !strings.EqualFold(te.Alias, "pivot") && !strings.EqualFold(te.Alias, "unpivot") {
			t.Errorf("%q: alias = %q, want pivot/unpivot", sql, te.Alias)
		}
	}
	// `t pivot, s` and `t pivot JOIN s ...` must also parse cleanly (the bare
	// keyword is the alias of t, then a comma/explicit join follows).
	for _, sql := range []string{
		"SELECT * FROM t pivot, s",
		"SELECT * FROM t pivot JOIN s ON s.x = t.x",
	} {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("%q: expected accept (pivot is an alias, then a join), got %v", sql, errs)
		}
	}
	// But `t PIVOT(...)` (operator-start) is still the operator, not an alias.
	te := tableExpr0(t, "SELECT * FROM t PIVOT(SUM(s) FOR q IN ('a'))")
	if te.Pivot == nil {
		t.Errorf("`t PIVOT(...)` should be an operator")
	}
	if te.Alias != "" {
		t.Errorf("`t PIVOT(...)`: table alias should be empty, got %q", te.Alias)
	}
}

// TABLESAMPLE size is restricted to possibly_cast_int_literal_or_parameter |
// floating_point_literal — a bare identifier / arithmetic / function call before
// the unit must reject (oracle-confirmed). Regression for the size over-accept
// (Codex finding).
func TestTableSample_SizeRestricted(t *testing.T) {
	reject := []string{
		"SELECT * FROM t TABLESAMPLE BERNOULLI (1 + 1 ROWS)",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (foo() ROWS)",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (x ROWS)",
	}
	for _, sql := range reject {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("%q: expected reject (non-literal sample size), got accept", sql)
		}
	}
	accept := []string{
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT)",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10.5 PERCENT)",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (CAST(10 AS INT64) PERCENT)",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (@p PERCENT)",
	}
	for _, sql := range accept {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("%q: expected accept (literal/cast/param size), got %v", sql, errs)
		}
	}
}

// The PIVOT FOR input is expression_higher_prec_than_and: arithmetic / field
// access accept, but a top-level comparison / NOT / AND in FOR rejects
// (oracle-confirmed). DEFEND of the FOR-precedence choice (Codex finding 1 was a
// grammar misread — the emulator rejects `FOR a = b IN` and `FOR NOT f IN`).
func TestPivot_ForPrecedence(t *testing.T) {
	accept := []string{
		"SELECT * FROM t PIVOT(SUM(s) FOR a + b IN ('x'))",
		"SELECT * FROM t PIVOT(SUM(s) FOR q.x IN ('x'))",
	}
	for _, sql := range accept {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("%q: expected accept, got %v", sql, errs)
		}
	}
	reject := []string{
		"SELECT * FROM t PIVOT(SUM(s) FOR a = b IN ('x'))",
		"SELECT * FROM t PIVOT(SUM(s) FOR NOT flag IN ('x'))",
		"SELECT * FROM t PIVOT(SUM(s) FOR a AND b IN ('x'))",
	}
	for _, sql := range reject {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("%q: expected reject (FOR is higher-prec-than-AND), got accept", sql)
		}
	}
}

// A bare (unparenthesized) multi-column UNPIVOT value list rejects; only the
// parenthesized group is valid (oracle-confirmed). DEFEND of the
// parsePathListWithOptParens single-vs-parenthesized handling (Codex finding 2).
func TestUnpivot_BareMultiColumnRejects(t *testing.T) {
	if _, errs := Parse("SELECT * FROM t UNPIVOT(s1, s2 FOR q IN ((Q1, Q2)))"); len(errs) == 0 {
		t.Errorf("bare multi-column value list should reject (must be parenthesized)")
	}
	if _, errs := Parse("SELECT * FROM t UNPIVOT((s1, s2) FOR q IN ((Q1, Q2)))"); len(errs) != 0 {
		t.Errorf("parenthesized multi-column value list should accept, got %v", errs)
	}
}

// The AST walker must descend into PIVOT/UNPIVOT/TABLESAMPLE child expressions —
// including on an UNNEST source (UnnestExpr.Pivot/Unpivot/Sample). Regression for
// the generated-walker hole on UNNEST operators (Codex finding 5): a visitor
// must observe the aggregate expression inside `UNNEST(...) PIVOT(SUM(x) ...)`.
func TestWalker_DescendsIntoOperators(t *testing.T) {
	// A FuncCall named SUM_SENTINEL under a PIVOT on a TableExpr source.
	if !walkFindsFuncCall(t, "SELECT * FROM t PIVOT(SUM_SENTINEL(s) FOR q IN ('a'))", "SUM_SENTINEL") {
		t.Errorf("walker did not reach the PIVOT aggregate on a TableExpr source")
	}
	// The same operator on an UNNEST source must also be walked (the UnnestExpr
	// case must visit Pivot/Unpivot/Sample, not just Array).
	if !walkFindsFuncCall(t, "SELECT * FROM UNNEST([1, 2]) PIVOT(SUM_SENTINEL(s) FOR q IN ('a'))", "SUM_SENTINEL") {
		t.Errorf("walker did not reach the PIVOT aggregate on an UNNEST source (walker hole)")
	}
	// A TABLESAMPLE size CAST under an UNNEST source must be reachable too.
	if !walkFindsFuncCall(t, "SELECT * FROM UNNEST([1, 2]) AS x PIVOT(SUM_SENTINEL(s) FOR q IN ('a')) TABLESAMPLE BERNOULLI (10 PERCENT)", "SUM_SENTINEL") {
		t.Errorf("walker did not reach a PIVOT aggregate under an UNNEST+TABLESAMPLE source")
	}
}

// walkFindsFuncCall parses sql and reports whether ast.Inspect reaches a
// *FuncCall whose name path matches name (case-insensitive last component).
func walkFindsFuncCall(t *testing.T, sql, name string) bool {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): %v", sql, errs)
	}
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if fc, ok := n.(*ast.FuncCall); ok && fc.Name != nil {
			parts := fc.Name.Parts
			if len(parts) > 0 && strings.EqualFold(parts[len(parts)-1], name) {
				found = true
			}
		}
		return true
	})
	return found
}

// --- reject cases ------------------------------------------------------------

func TestQueryClauses_Reject(t *testing.T) {
	reject := []string{
		// PIVOT requires an aggregate, a FOR, and a non-empty IN list.
		"SELECT * FROM t PIVOT(FOR q IN ('Q1'))",
		"SELECT * FROM t PIVOT(SUM(s) FOR q IN ())",
		"SELECT * FROM t PIVOT(SUM(s) q IN ('Q1'))", // missing FOR
		// UNPIVOT requires value col / FOR / non-empty IN.
		"SELECT * FROM t UNPIVOT()",
		"SELECT * FROM t UNPIVOT(s FOR q IN ())",
		// At most one pivot/unpivot per source.
		"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) UNPIVOT(x FOR y IN (a))",
		// TABLESAMPLE requires a method, a size, and a unit.
		"SELECT * FROM t TABLESAMPLE (10 PERCENT)",       // missing method
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10)",     // missing unit
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10 FOO)", // bad unit
		// TABLESAMPLE takes no trailing alias.
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) AS x",
		// TABLESAMPLE before PIVOT is rejected (sample is outermost).
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10 PERCENT) PIVOT(SUM(s) FOR q IN ('a'))",
		// PIVOT alias cannot carry a column list.
		"SELECT * FROM t PIVOT(SUM(s) FOR q IN ('Q1')) AS p (x, y)",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			_, errs := Parse(sql)
			if len(errs) == 0 {
				t.Errorf("Parse(%q): expected a parse error, got none", sql)
			}
		})
	}
}
