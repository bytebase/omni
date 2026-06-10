package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// firstSelect parses input and returns the first *ast.SelectStmt, failing the
// test on any parse error or a non-SELECT result. It tolerates input wrapped
// in a set operation by unwrapping the leftmost SELECT.
func firstSelect(t *testing.T, input string) *ast.SelectStmt {
	t.Helper()
	result := ParseBestEffort(input)
	if len(result.Errors) > 0 {
		t.Fatalf("parse %q: unexpected errors: %v", input, result.Errors)
	}
	if len(result.File.Stmts) == 0 {
		t.Fatalf("parse %q: no statements", input)
	}
	return unwrapSelect(t, result.File.Stmts[0])
}

func unwrapSelect(t *testing.T, n ast.Node) *ast.SelectStmt {
	t.Helper()
	switch s := n.(type) {
	case *ast.SelectStmt:
		return s
	case *ast.SetOperationStmt:
		return unwrapSelect(t, s.Left)
	default:
		t.Fatalf("expected SelectStmt, got %T", n)
		return nil
	}
}

// firstTableRef returns the first FROM item of a SELECT as a *ast.TableRef.
func firstTableRef(t *testing.T, sel *ast.SelectStmt) *ast.TableRef {
	t.Helper()
	if len(sel.From) == 0 {
		t.Fatal("SELECT has no FROM items")
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("FROM[0] is %T, want *ast.TableRef", sel.From[0])
	}
	return ref
}

// ---------------------------------------------------------------------------
// PIVOT
// ---------------------------------------------------------------------------

func TestPivot_AnyOrderBy(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter)) ORDER BY empid")
	ref := firstTableRef(t, sel)
	if ref.Pivot == nil {
		t.Fatal("expected Pivot clause")
	}
	pv := ref.Pivot
	if pv.Agg == nil || strings.ToUpper(pv.Agg.Name.Name.Name) != "SUM" {
		t.Fatalf("agg = %+v, want SUM(...)", pv.Agg)
	}
	if pv.ForColumn == nil || len(pv.ForColumn.Parts) != 1 || pv.ForColumn.Parts[0].Name != "quarter" {
		t.Fatalf("for column = %+v, want quarter", pv.ForColumn)
	}
	if pv.In == nil || pv.In.Kind != ast.PivotInAny {
		t.Fatalf("IN kind = %v, want PivotInAny", pv.In.Kind)
	}
	if len(pv.In.OrderBy) != 1 {
		t.Fatalf("ANY ORDER BY items = %d, want 1", len(pv.In.OrderBy))
	}
	// Loc sanity: the clause must span from PIVOT to its closing paren.
	if pv.Loc.Start <= 0 || pv.Loc.End <= pv.Loc.Start {
		t.Errorf("pivot Loc invalid: %+v", pv.Loc)
	}
}

func TestPivot_AggAlias(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) AS total FOR quarter IN (ANY ORDER BY quarter))")
	pv := firstTableRef(t, sel).Pivot
	if pv == nil {
		t.Fatal("expected Pivot")
	}
	if pv.AggAlias.Name != "total" {
		t.Errorf("agg alias = %q, want total", pv.AggAlias.Name)
	}
}

func TestPivot_Subquery(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN (SELECT DISTINCT quarter FROM ad_campaign_types_by_quarter WHERE television = TRUE ORDER BY quarter))")
	pv := firstTableRef(t, sel).Pivot
	if pv.In.Kind != ast.PivotInSubquery {
		t.Fatalf("IN kind = %v, want PivotInSubquery", pv.In.Kind)
	}
	if pv.In.Subquery == nil {
		t.Error("expected IN subquery")
	}
}

func TestPivot_ValueList(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN ('2023_Q1', '2023_Q2', '2023_Q3'))")
	pv := firstTableRef(t, sel).Pivot
	if pv.In.Kind != ast.PivotInValues {
		t.Fatalf("IN kind = %v, want PivotInValues", pv.In.Kind)
	}
	if len(pv.In.Values) != 3 {
		t.Fatalf("values = %d, want 3", len(pv.In.Values))
	}
	for i, v := range pv.In.Values {
		lit, ok := v.Value.(*ast.Literal)
		if !ok || lit.Kind != ast.LitString {
			t.Errorf("value[%d] = %T, want string literal", i, v.Value)
		}
	}
}

func TestPivot_ValueListWithAliases(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN ('2023_Q1' AS q1, '2023_Q2' AS q2))")
	pv := firstTableRef(t, sel).Pivot
	if len(pv.In.Values) != 2 {
		t.Fatalf("values = %d, want 2", len(pv.In.Values))
	}
	if pv.In.Values[0].Alias.Name != "q1" || pv.In.Values[1].Alias.Name != "q2" {
		t.Errorf("value aliases = %q,%q want q1,q2", pv.In.Values[0].Alias.Name, pv.In.Values[1].Alias.Name)
	}
}

func TestPivot_DefaultOnNull(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN (ANY ORDER BY quarter) DEFAULT ON NULL (0))")
	pv := firstTableRef(t, sel).Pivot
	if pv.DefaultVal == nil {
		t.Fatal("expected DEFAULT ON NULL value")
	}
	lit, ok := pv.DefaultVal.(*ast.Literal)
	if !ok || lit.Ival != 0 {
		t.Errorf("default = %+v, want literal 0", pv.DefaultVal)
	}
}

func TestPivot_DefaultOnNullValueList(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN ('2023_Q1', '2023_Q2') DEFAULT ON NULL (0))")
	pv := firstTableRef(t, sel).Pivot
	if pv.In.Kind != ast.PivotInValues || len(pv.In.Values) != 2 {
		t.Fatalf("IN = %+v, want 2 values", pv.In)
	}
	if pv.DefaultVal == nil {
		t.Error("expected DEFAULT ON NULL value")
	}
}

func TestPivot_TrailingAlias(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM quarterly_sales PIVOT(SUM(amount) FOR quarter IN (ANY)) AS p")
	pv := firstTableRef(t, sel).Pivot
	if pv.Alias.Name != "p" {
		t.Errorf("pivot alias = %q, want p", pv.Alias.Name)
	}
}

func TestPivot_Chained(t *testing.T) {
	// Two PIVOTs in a row (corpus example_20) over a subquery source.
	in := `SELECT q1, q2 FROM
	  (SELECT amount, quarter AS quarter_amount, quarter AS quarter_discount, discount_percent FROM quarterly_sales)
	  PIVOT (SUM(amount) FOR quarter_amount IN ('2023_Q1', '2023_Q2'))
	  PIVOT (MAX(discount_percent) FOR quarter_discount IN ('2023_Q1', '2023_Q2'))`
	sel := firstSelect(t, in)
	ref := firstTableRef(t, sel)
	// The outer source must itself carry a PIVOT, and its inner source the other.
	if ref.Pivot == nil {
		t.Fatal("outer table ref has no PIVOT")
	}
	// The outer clause is the SECOND pivot (MAX); the keyword PIVOT must not
	// have been eaten as the first clause's implicit alias.
	if got := ref.Pivot.Agg.Name.Name.Normalize(); got != "MAX" {
		t.Errorf("outer pivot agg = %q, want MAX", got)
	}
	if ref.Nested == nil {
		t.Fatal("outer table ref has no Nested source for the chain")
	}
	inner := ref.Nested
	if inner.Pivot == nil {
		t.Fatal("nested table ref has no PIVOT")
	}
	if got := inner.Pivot.Agg.Name.Name.Normalize(); got != "SUM" {
		t.Errorf("inner pivot agg = %q, want SUM", got)
	}
	if inner.Pivot.Alias.Name != "" {
		t.Errorf("inner pivot alias = %q, want empty (PIVOT keyword must not become an alias)", inner.Pivot.Alias.Name)
	}
	if inner.Subquery == nil {
		t.Error("nested table ref lost its subquery source")
	}
	// Loc sanity: the outer ref spans from the subquery start through the
	// second pivot's closing paren.
	if ref.Loc.Start != inner.Loc.Start {
		t.Errorf("outer Loc.Start = %d, want %d (inner start)", ref.Loc.Start, inner.Loc.Start)
	}
	if ref.Loc.End != ref.Pivot.Loc.End {
		t.Errorf("outer Loc.End = %d, want %d (outer pivot end)", ref.Loc.End, ref.Pivot.Loc.End)
	}
}

func TestPivot_ChainedTrailingAlias(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t PIVOT (SUM(a) FOR b IN ('x')) PIVOT (MAX(c) FOR d IN ('y')) AS pp")
	ref := firstTableRef(t, sel)
	if ref.Pivot == nil || ref.Nested == nil || ref.Nested.Pivot == nil {
		t.Fatalf("chain shape: outer pivot=%v nested=%v", ref.Pivot != nil, ref.Nested != nil)
	}
	if ref.Pivot.Alias.Name != "pp" {
		t.Errorf("outer pivot alias = %q, want pp", ref.Pivot.Alias.Name)
	}
	if ref.Nested.Name == nil {
		t.Error("nested ref lost its table name")
	}
}

func TestPivot_ChainedMixedUnpivotPivot(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t UNPIVOT (v FOR n IN (a, b)) PIVOT (SUM(v) FOR n IN ('a'))")
	ref := firstTableRef(t, sel)
	if ref.Pivot == nil {
		t.Fatal("outer ref must carry the PIVOT")
	}
	if ref.Nested == nil || ref.Nested.Unpivot == nil {
		t.Fatal("nested ref must carry the UNPIVOT")
	}
	if ref.Unpivot != nil {
		t.Error("outer ref must not also carry the UNPIVOT")
	}
}

func TestPivot_ChainedThenSample(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM t PIVOT (SUM(a) FOR b IN ('x')) PIVOT (MAX(c) FOR d IN ('y')) SAMPLE (10)")
	ref := firstTableRef(t, sel)
	if ref.Pivot == nil || ref.Nested == nil {
		t.Fatal("expected chained pivot shape")
	}
	if ref.Sample == nil {
		t.Error("trailing SAMPLE must attach to the outer pivoted ref")
	}
	if ref.Nested.Sample != nil {
		t.Error("SAMPLE must not leak onto the nested ref")
	}
}

func TestPivot_OverSubquerySource(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM (SELECT * FROM t) PIVOT(SUM(a) FOR b IN (ANY))")
	ref := firstTableRef(t, sel)
	if ref.Subquery == nil {
		t.Fatal("expected subquery source")
	}
	if ref.Pivot == nil {
		t.Fatal("expected PIVOT on subquery source")
	}
}

// ---------------------------------------------------------------------------
// UNPIVOT
// ---------------------------------------------------------------------------

func TestUnpivot_Basic(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM monthly_sales UNPIVOT (sales FOR month IN (jan, feb, mar, apr)) ORDER BY empid")
	uv := firstTableRef(t, sel).Unpivot
	if uv == nil {
		t.Fatal("expected Unpivot clause")
	}
	if uv.ValueColumn.Name != "sales" {
		t.Errorf("value column = %q, want sales", uv.ValueColumn.Name)
	}
	if uv.NameColumn.Name != "month" {
		t.Errorf("name column = %q, want month", uv.NameColumn.Name)
	}
	if len(uv.Columns) != 4 {
		t.Fatalf("columns = %d, want 4", len(uv.Columns))
	}
	if uv.NullsMode != ast.UnpivotNullsUnspecified {
		t.Errorf("nulls mode = %v, want unspecified", uv.NullsMode)
	}
}

func TestUnpivot_ColumnAliases(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM monthly_sales UNPIVOT (sales FOR month IN (jan AS january, feb AS february))")
	uv := firstTableRef(t, sel).Unpivot
	if len(uv.Columns) != 2 {
		t.Fatalf("columns = %d, want 2", len(uv.Columns))
	}
	if uv.Columns[0].Alias.Name != "january" || uv.Columns[1].Alias.Name != "february" {
		t.Errorf("aliases = %q,%q want january,february", uv.Columns[0].Alias.Name, uv.Columns[1].Alias.Name)
	}
}

func TestUnpivot_IncludeNulls(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM monthly_sales UNPIVOT INCLUDE NULLS (sales FOR month IN (jan, feb, mar, apr))")
	uv := firstTableRef(t, sel).Unpivot
	if uv.NullsMode != ast.UnpivotIncludeNulls {
		t.Errorf("nulls mode = %v, want IncludeNulls", uv.NullsMode)
	}
}

func TestUnpivot_ExcludeNulls(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM monthly_sales UNPIVOT EXCLUDE NULLS (sales FOR month IN (jan, feb))")
	uv := firstTableRef(t, sel).Unpivot
	if uv.NullsMode != ast.UnpivotExcludeNulls {
		t.Errorf("nulls mode = %v, want ExcludeNulls", uv.NullsMode)
	}
}

func TestUnpivot_TrailingAliasThenJoin(t *testing.T) {
	// corpus example_06: UNPIVOT (...) unpvt JOIN LATERAL (...)
	sel := firstSelect(t, "SELECT * FROM monthly_sales UNPIVOT (sales FOR month IN (jan, feb)) unpvt JOIN LATERAL (SELECT unpvt.sales AS sales_value) jl")
	// The FROM item is a JoinExpr whose Left carries the UNPIVOT + alias.
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want *ast.JoinExpr", sel.From[0])
	}
	left, ok := join.Left.(*ast.TableRef)
	if !ok || left.Unpivot == nil {
		t.Fatalf("join.Left = %T (unpivot=%v), want TableRef with UNPIVOT", join.Left, left)
	}
	if left.Unpivot.Alias.Name != "unpvt" {
		t.Errorf("unpivot alias = %q, want unpvt", left.Unpivot.Alias.Name)
	}
}

// ---------------------------------------------------------------------------
// SAMPLE / TABLESAMPLE
// ---------------------------------------------------------------------------

func TestSample_Probability(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable SAMPLE (10)")
	s := firstTableRef(t, sel).Sample
	if s == nil {
		t.Fatal("expected Sample clause")
	}
	if s.Keyword != ast.SampleKwSample {
		t.Errorf("keyword = %v, want SAMPLE", s.Keyword)
	}
	if s.Method != ast.SampleMethodUnspecified {
		t.Errorf("method = %v, want unspecified", s.Method)
	}
	if s.Rows {
		t.Error("rows should be false")
	}
	lit, ok := s.Quantity.(*ast.Literal)
	if !ok || lit.Ival != 10 {
		t.Errorf("quantity = %+v, want 10", s.Quantity)
	}
}

func TestSample_TablesampleBernoulli(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable TABLESAMPLE BERNOULLI (20.3)")
	s := firstTableRef(t, sel).Sample
	if s.Keyword != ast.SampleKwTablesample {
		t.Errorf("keyword = %v, want TABLESAMPLE", s.Keyword)
	}
	if s.Method != ast.SampleMethodBernoulli {
		t.Errorf("method = %v, want BERNOULLI", s.Method)
	}
}

func TestSample_Row(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable SAMPLE ROW (0)")
	s := firstTableRef(t, sel).Sample
	if s.Method != ast.SampleMethodRow {
		t.Errorf("method = %v, want ROW", s.Method)
	}
}

func TestSample_SystemSeed(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable SAMPLE SYSTEM (3) SEED (82)")
	s := firstTableRef(t, sel).Sample
	if s.Method != ast.SampleMethodSystem {
		t.Errorf("method = %v, want SYSTEM", s.Method)
	}
	if s.SeedKind != ast.SampleSeedSeed {
		t.Errorf("seed kind = %v, want SEED", s.SeedKind)
	}
	lit, ok := s.Seed.(*ast.Literal)
	if !ok || lit.Ival != 82 {
		t.Errorf("seed = %+v, want 82", s.Seed)
	}
}

func TestSample_BlockRepeatable(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable SAMPLE BLOCK (0.012) REPEATABLE (99992)")
	s := firstTableRef(t, sel).Sample
	if s.Method != ast.SampleMethodBlock {
		t.Errorf("method = %v, want BLOCK", s.Method)
	}
	if s.SeedKind != ast.SampleSeedRepeatable {
		t.Errorf("seed kind = %v, want REPEATABLE", s.SeedKind)
	}
}

func TestSample_NRows(t *testing.T) {
	sel := firstSelect(t, "SELECT * FROM testtable SAMPLE (10 ROWS)")
	s := firstTableRef(t, sel).Sample
	if !s.Rows {
		t.Error("rows should be true")
	}
	lit, ok := s.Quantity.(*ast.Literal)
	if !ok || lit.Ival != 10 {
		t.Errorf("quantity = %+v, want 10", s.Quantity)
	}
}

func TestSample_OnJoinedTable(t *testing.T) {
	// corpus example_06: only the second table is sampled.
	sel := firstSelect(t, "SELECT i, j FROM table1 AS t1 INNER JOIN table2 AS t2 SAMPLE (50) WHERE t2.j = t1.i")
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("FROM[0] = %T, want JoinExpr", sel.From[0])
	}
	right, ok := join.Right.(*ast.TableRef)
	if !ok || right.Sample == nil {
		t.Fatalf("join.Right = %T, want TableRef with SAMPLE", join.Right)
	}
	left, ok := join.Left.(*ast.TableRef)
	if !ok || left.Sample != nil {
		t.Fatalf("join.Left should not carry a SAMPLE")
	}
}

func TestSample_OnSubquery(t *testing.T) {
	// corpus example_07: (subquery) SAMPLE (1).
	sel := firstSelect(t, "SELECT * FROM (SELECT * FROM t1 JOIN t2 ON t1.a = t2.c) SAMPLE (1)")
	ref := firstTableRef(t, sel)
	if ref.Subquery == nil || ref.Sample == nil {
		t.Fatalf("expected subquery source with SAMPLE; got subquery=%v sample=%v", ref.Subquery != nil, ref.Sample != nil)
	}
}

// ---------------------------------------------------------------------------
// Negative cases
// ---------------------------------------------------------------------------

func TestPivotSample_Negatives(t *testing.T) {
	cases := []string{
		"SELECT * FROM t PIVOT(SUM(a) quarter IN (ANY))",   // missing FOR
		"SELECT * FROM t PIVOT(SUM(a) FOR b (ANY))",        // missing IN
		"SELECT * FROM t PIVOT(amount FOR b IN (ANY))",     // aggregate not a func call
		"SELECT * FROM t UNPIVOT (v month IN (a, b))",      // missing FOR
		"SELECT * FROM t UNPIVOT INCLUDE (v FOR m IN (a))", // INCLUDE without NULLS
		"SELECT * FROM t SAMPLE",                           // missing ( … )
		"SELECT * FROM t SAMPLE BERNOULLI",                 // method without ( … )
		"SELECT * FROM t SAMPLE (10) SEED",                 // SEED without ( … )
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			result := ParseBestEffort(in)
			if len(result.Errors) == 0 {
				t.Errorf("expected a parse error for %q, got none", in)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Corpus differential — every SELECT statement in the owned corpus dirs must
// parse with zero errors. CREATE/INSERT/UPDATE/ALTER setup statements and the
// $-positional examples are filtered (they belong to other nodes).
// ---------------------------------------------------------------------------

// selectCorpusDirs are the official corpus directories owned by this node whose
// SELECT statements exercise the T5.3 clauses.
var selectCorpusDirs = []string{
	"testdata/official/pivot",
	"testdata/official/unpivot",
	"testdata/official/sample",
	"testdata/official/match-recognize",
	"testdata/official/connect-by",
	"testdata/official/at-before",
}

func TestT53_OwnedCorpusParses(t *testing.T) {
	for _, dir := range selectCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				for _, seg := range Split(string(data)) {
					text := strings.TrimSpace(seg.Text)
					if text == "" {
						continue
					}
					upper := strings.ToUpper(text)
					// Only exercise query statements; skip DDL/DML setup.
					if !(strings.HasPrefix(upper, "SELECT") || strings.HasPrefix(upper, "WITH")) {
						continue
					}
					// The $-positional pivot example (example_20) selects $1..$8
					// AFTER the pivot; those positional refs belong to the
					// expression node, not this clause node. Skip it.
					if strings.Contains(text, "$1") || strings.Contains(text, "$2") {
						continue
					}
					node, errs := parseSingle(seg.Text, seg.ByteStart)
					if len(errs) > 0 {
						t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
						continue
					}
					if node == nil {
						t.Errorf("statement %q parsed to nil", text)
					}
				}
			})
		}
	}
}
