package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// TestParseInlineAnalyticView reads the derived analytic view the page draws
// and the pinned Oracle takes: ANALYTIC VIEW ( sub_av_clause ) [ alias ] in
// a FROM clause, and subav_name ANALYTIC VIEW AS ( sub_av_clause ) in a WITH
// clause. The page's own example is the control:
//
//	SELECT ... FROM ANALYTIC VIEW ( USING sales_av
//	  HIERARCHIES (time_hier)
//	  ADD MEASURES ( lag_sales AS (LAG(sales) OVER (HIERARCHY time_hier OFFSET 1)) ) )
func TestParseInlineAnalyticView(t *testing.T) {
	result := ParseAndCheck(t, `
SELECT time_hier.member_name, sales, lag_sales
FROM ANALYTIC VIEW (
  USING sales_av
  HIERARCHIES (time_hier)
  FILTER FACT (time_hier TO quarter_of_year IN (1, 2))
  ADD MEASURES (lag_sales AS (LAG(sales) OVER (HIERARCHY time_hier OFFSET 1)))
) av`)

	stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
	if stmt.FromClause == nil || stmt.FromClause.Len() != 1 {
		t.Fatalf("FromClause len = %v, want 1", stmt.FromClause)
	}
	ref, ok := stmt.FromClause.Items[0].(*ast.InlineAnalyticView)
	if !ok {
		t.Fatalf("FROM item = %T, want *ast.InlineAnalyticView", stmt.FromClause.Items[0])
	}
	if ref.Clause == nil || ref.Clause.Using == nil || ref.Clause.Using.Name != "SALES_AV" {
		t.Fatalf("Clause.Using = %#v, want SALES_AV", ref.Clause)
	}
	if ref.Clause.Using.Schema != "" {
		t.Fatalf("Clause.Using.Schema = %q, want none", ref.Clause.Using.Schema)
	}
	if ref.Clause.Options == nil || ref.Clause.Options.Len() != 3 {
		t.Fatalf("Clause.Options len = %v, want 3", ref.Clause.Options)
	}
	for i, want := range []string{"HIERARCHIES", "FILTER FACT", "ADD MEASURES"} {
		opt, ok := ref.Clause.Options.Items[i].(*ast.DDLOption)
		if !ok || opt.Key != want {
			t.Fatalf("Options[%d] = %#v, want %s", i, ref.Clause.Options.Items[i], want)
		}
	}
	// The filter's condition is parsed, not just recorded as text.
	filter, ok := ref.Clause.Options.Items[1].(*ast.DDLOption)
	if !ok || filter.Items == nil || filter.Items.Len() != 1 || filter.Items.Items[0] == nil {
		t.Fatalf("filter clause = %#v, want its parsed condition", ref.Clause.Options.Items[1])
	}
	// A calculated measure names a column the derived view adds.
	measures, ok := ref.Clause.Options.Items[2].(*ast.DDLOption)
	if !ok || measures.Items == nil || measures.Items.Len() != 1 {
		t.Fatalf("add measures clause = %#v, want its measure", ref.Clause.Options.Items[2])
	}
	if measure, ok := measures.Items.Items[0].(*ast.ResTarget); !ok || measure.Name != "LAG_SALES" {
		t.Fatalf("measure = %#v, want LAG_SALES", measures.Items.Items[0])
	}
	if ref.Alias == nil || ref.Alias.Name != "AV" {
		t.Fatalf("Alias = %#v, want AV", ref.Alias)
	}
}

// TestParseInlineAnalyticViewSuffixes reads the clauses a
// query_table_expression names after its source: a flashback query and
// SAMPLE, before the reference's alias.
func TestParseInlineAnalyticViewSuffixes(t *testing.T) {
	ref := func(sql string) *ast.InlineAnalyticView {
		t.Helper()
		stmt := ParseAndCheck(t, sql).Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
		got, ok := stmt.FromClause.Items[0].(*ast.InlineAnalyticView)
		if !ok {
			t.Fatalf("FROM item = %T, want *ast.InlineAnalyticView", stmt.FromClause.Items[0])
		}
		return got
	}

	if got := ref(`SELECT * FROM ANALYTIC VIEW(USING a) SAMPLE (10)`); got.Sample == nil || !strings.Contains(ast.NodeToString(got), ":sample") {
		t.Fatalf("SAMPLE after an inline analytic view = %s", ast.NodeToString(got))
	}
	if got := ref(`SELECT * FROM ANALYTIC VIEW(USING a) av SAMPLE (10)`); got.Sample == nil || got.Alias == nil || got.Alias.Name != "AV" {
		t.Fatalf("inline_av_alias before SAMPLE: sample %#v, alias %#v", got.Sample, got.Alias)
	}
	if got := ref(`SELECT * FROM ANALYTIC VIEW(USING a) SAMPLE`); got.Sample != nil || got.Alias == nil || got.Alias.Name != "SAMPLE" {
		t.Fatalf("bare SAMPLE after an inline analytic view: sample %#v, alias %#v", got.Sample, got.Alias)
	}
	if got := ref(`SELECT * FROM ANALYTIC VIEW(USING a) AS OF SCN 1 av`); got.Flashback == nil || got.Alias == nil || got.Alias.Name != "AV" || !strings.Contains(ast.NodeToString(got), ":flashback") {
		t.Fatalf("AS OF after an inline analytic view: %s", ast.NodeToString(got))
	}
	mrStmt := ParseAndCheck(t, `SELECT * FROM ANALYTIC VIEW(USING a) MATCH_RECOGNIZE (ORDER BY id PATTERN (p) DEFINE p AS p.n > 0) mr`).Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
	mr, ok := mrStmt.FromClause.Items[0].(*ast.MatchRecognizeClause)
	if !ok || mr.Alias == nil || mr.Alias.Name != "MR" {
		t.Fatalf("MATCH_RECOGNIZE item = %#v", mrStmt.FromClause.Items[0])
	}
	if _, ok := mr.Source.(*ast.InlineAnalyticView); !ok {
		t.Fatalf("MATCH_RECOGNIZE source = %T, want *ast.InlineAnalyticView", mr.Source)
	}
	if got := ref(`SELECT * FROM (ANALYTIC VIEW(USING a)) av`); got.Alias == nil || got.Alias.Name != "AV" {
		t.Fatalf("parenthesized inline analytic view alias = %#v, want AV", got.Alias)
	}
	if got := ref(`SELECT * FROM ((ANALYTIC VIEW(USING a))) AS av`); got.Alias == nil || got.Alias.Name != "AV" {
		t.Fatalf("nested parenthesized inline analytic view alias = %#v, want AV", got.Alias)
	}
	ParseShouldFail(t, `DELETE a FROM ANALYTIC VIEW(USING a) SAMPLE (10)`)
}

// TestParseInlineAnalyticViewDelete reads the same reference where DELETE
// takes its table: dml_table_expression_clause's FROM names a
// table_reference, which holds the inline analytic view.
func TestParseInlineAnalyticViewDelete(t *testing.T) {
	result := ParseAndCheck(t, `DELETE a FROM ANALYTIC VIEW(USING a) v`)

	stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.DeleteStmt)
	if stmt.FromClause == nil || stmt.FromClause.Len() != 1 {
		t.Fatalf("FromClause len = %v, want 1", stmt.FromClause)
	}
	ref, ok := stmt.FromClause.Items[0].(*ast.InlineAnalyticView)
	if !ok {
		t.Fatalf("FROM item = %T, want *ast.InlineAnalyticView", stmt.FromClause.Items[0])
	}
	if ref.Clause == nil || ref.Clause.Using == nil || ref.Clause.Using.Name != "A" {
		t.Fatalf("Clause.Using = %#v, want A", ref.Clause)
	}
	if ref.Alias == nil || ref.Alias.Name != "V" {
		t.Fatalf("Alias = %#v, want V", ref.Alias)
	}
}

// TestParseSubavFactoringClause reads the transitory analytic view a WITH
// clause declares, alone and beside a query element:
//
//	WITH my_av ANALYTIC VIEW AS ( USING sales_av HIERARCHIES (time_hier) ) ...
func TestParseSubavFactoringClause(t *testing.T) {
	result := ParseAndCheck(t, `
WITH my_av ANALYTIC VIEW AS ( USING sales_av HIERARCHIES (time_hier) ),
     other AS ( SELECT 1 FROM dual )
SELECT time_hier.member_name FROM my_av`)

	stmt := result.Items[0].(*ast.RawStmt).Stmt.(*ast.SelectStmt)
	if stmt.WithClause == nil || stmt.WithClause.CTEs == nil || stmt.WithClause.CTEs.Len() != 2 {
		t.Fatalf("CTEs = %#v, want 2", stmt.WithClause)
	}
	first, ok := stmt.WithClause.CTEs.Items[0].(*ast.CTE)
	if !ok {
		t.Fatalf("first CTE = %T", stmt.WithClause.CTEs.Items[0])
	}
	if first.Name != "MY_AV" || first.Query != nil || first.SubAV == nil {
		t.Fatalf("first CTE = %#v, want a derived analytic view named MY_AV", first)
	}
	if first.SubAV.Using == nil || first.SubAV.Using.Name != "SALES_AV" {
		t.Fatalf("SubAV.Using = %#v, want SALES_AV", first.SubAV.Using)
	}
	second, ok := stmt.WithClause.CTEs.Items[1].(*ast.CTE)
	if !ok || second.Query == nil || second.SubAV != nil {
		t.Fatalf("second CTE = %#v, want a query named OTHER", second)
	}
}

// TestDerivedAnalyticViewRefusals pins where the parser stops in the forms
// the pinned Oracle refuses, at the byte it reports (compat.tsv). It wants
// USING at the first token inside an inline analytic view's parentheses, and
// ANALYTIC VIEW at what a WITH element opens with, so the page's older
// reading of the clause — a nested ANALYTIC VIEW — is refused exactly where
// Oracle refuses it.
func TestDerivedAnalyticViewRefusals(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// ORA-02000, "missing USING keyword", at the token inside the parens.
		{`SELECT * FROM ANALYTIC VIEW(ANALYTIC VIEW USING a)`, 28},
		{`SELECT * FROM ANALYTIC VIEW(VIEW USING a)`, 28},
		{`SELECT * FROM ANALYTIC VIEW(ANALYTIC USING a)`, 28},
		{`SELECT * FROM ANALYTIC VIEW(a)`, 28},
		{`DELETE a FROM ANALYTIC VIEW(ANALYTIC VIEW USING a)`, 28},
		// The same token: the WITH element's clause wants USING at its '(' + 1.
		{`WITH a ANALYTIC VIEW AS(ANALYTIC VIEW USING a) SELECT *`, 24},
		// ORA-00906, "missing left parenthesis", where the parens belong.
		{`SELECT * FROM ANALYTIC VIEW USING a`, 28},
		// Pinned Oracle 23.26.3 refuses multiple hierarchy IDs in one filter
		// clause and added base measures without AGGREGATE BY.
		{`SELECT * FROM ANALYTIC VIEW(USING a FILTER FACT(a, a TO 1 = 1))`, 49},
		{`SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(a))`, 50},
		{`SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(a FACT a))`, 57},
		{`SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(a SUM))`, 51},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// The controls, which Oracle reads to a name lookup: the clause starts at
	// USING, in a schema-qualified name, with an empty HIERARCHIES list, and
	// with the WITH element's ANALYTIC VIEW outside the clause.
	for _, sql := range []string{
		`SELECT * FROM ANALYTIC VIEW(USING a)`,
		`SELECT * FROM ANALYTIC VIEW(USING s.a) a`,
		`SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES())`,
		`SELECT * FROM ANALYTIC VIEW(USING a HIERARCHIES(d.h, h))`,
		`SELECT * FROM ANALYTIC VIEW(USING a FILTER FACT(h TO 1 = 1, d.h TO 2 = 2))`,
		`SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(m AS(1)))`,
		`SELECT * FROM ANALYTIC VIEW(USING a ADD MEASURES(m FACT c AGGREGATE BY SUM))`,
		`DELETE a FROM ANALYTIC VIEW(USING a)`,
		`WITH a ANALYTIC VIEW AS(USING a) SELECT *`,
		`WITH a ANALYTIC VIEW AS(USING a), b ANALYTIC VIEW AS(USING a) SELECT *`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want it to parse", sql, err)
		}
	}
}
