package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// MATCH_RECOGNIZE
// ---------------------------------------------------------------------------

func TestMatchRecognize_Full(t *testing.T) {
	// corpus example_03 — all sub-clauses present.
	in := `SELECT * FROM stock_price_history
  MATCH_RECOGNIZE(
    PARTITION BY company
    ORDER BY price_date
    MEASURES
      MATCH_NUMBER() AS match_number,
      FIRST(price_date) AS start_date,
      LAST(price_date) AS end_date,
      COUNT(*) AS rows_in_sequence,
      COUNT(row_with_price_decrease.*) AS num_decreases,
      COUNT(row_with_price_increase.*) AS num_increases
    ONE ROW PER MATCH
    AFTER MATCH SKIP TO LAST row_with_price_increase
    PATTERN(row_before_decrease row_with_price_decrease+ row_with_price_increase+)
    DEFINE
      row_with_price_decrease AS price < LAG(price),
      row_with_price_increase AS price > LAG(price)
  )
ORDER BY company, match_number`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr == nil {
		t.Fatal("expected MatchRecognize clause")
	}
	if len(mr.PartitionBy) != 1 {
		t.Errorf("partition by = %d, want 1", len(mr.PartitionBy))
	}
	if len(mr.OrderBy) != 1 {
		t.Errorf("order by = %d, want 1", len(mr.OrderBy))
	}
	if len(mr.Measures) != 6 {
		t.Fatalf("measures = %d, want 6", len(mr.Measures))
	}
	if mr.Measures[0].Alias.Name != "match_number" {
		t.Errorf("measure[0] alias = %q, want match_number", mr.Measures[0].Alias.Name)
	}
	if mr.RowsPerMatch == nil || mr.RowsPerMatch.Kind != ast.OneRowPerMatch {
		t.Errorf("rows per match = %+v, want ONE ROW", mr.RowsPerMatch)
	}
	if mr.AfterMatch == nil || mr.AfterMatch.Kind != ast.AfterMatchSkipToLast {
		t.Fatalf("after match = %+v, want SKIP TO LAST", mr.AfterMatch)
	}
	if mr.AfterMatch.Symbol.Name != "row_with_price_increase" {
		t.Errorf("skip-to symbol = %q", mr.AfterMatch.Symbol.Name)
	}
	if mr.Pattern == nil {
		t.Fatal("expected PATTERN")
	}
	wantPat := "row_before_decrease row_with_price_decrease+ row_with_price_increase+"
	if mr.Pattern.Raw != wantPat {
		t.Errorf("pattern raw = %q, want %q", mr.Pattern.Raw, wantPat)
	}
	if len(mr.Define) != 2 {
		t.Fatalf("define = %d, want 2", len(mr.Define))
	}
	if mr.Define[0].Symbol.Name != "row_with_price_decrease" {
		t.Errorf("define[0] symbol = %q", mr.Define[0].Symbol.Name)
	}
	if mr.Define[0].Cond == nil {
		t.Error("define[0] cond is nil")
	}
}

func TestMatchRecognize_LowercaseOnSubquery(t *testing.T) {
	// corpus example_04 — lowercase, over a parenthesized subquery, ALL ROWS,
	// quoted measure alias, classifier()/match_sequence_number().
	in := `select price_date, match_number from
  (select * from stock_price_history where company='ABCD') match_recognize(
    order by price_date
    measures
        match_number() as "MATCH_NUMBER",
        match_sequence_number() as msq,
        classifier() as cl
    all rows per match
    pattern(ANY_ROW UP+)
    define
        ANY_ROW AS TRUE,
        UP as price > lag(price)
)
order by match_number, msq`
	sel := firstSelect(t, in)
	ref := firstTableRef(t, sel)
	if ref.Subquery == nil {
		t.Fatal("expected subquery source")
	}
	mr := ref.MatchRecognize
	if mr == nil {
		t.Fatal("expected MatchRecognize on subquery source")
	}
	if mr.RowsPerMatch == nil || mr.RowsPerMatch.Kind != ast.AllRowsPerMatch {
		t.Errorf("rows per match = %+v, want ALL ROWS", mr.RowsPerMatch)
	}
	if len(mr.Measures) != 3 {
		t.Fatalf("measures = %d, want 3", len(mr.Measures))
	}
	if !mr.Measures[0].Alias.Quoted || mr.Measures[0].Alias.Name != "MATCH_NUMBER" {
		t.Errorf("measure[0] alias = %+v, want quoted MATCH_NUMBER", mr.Measures[0].Alias)
	}
	if mr.Pattern.Raw != "ANY_ROW UP+" {
		t.Errorf("pattern = %q, want ANY_ROW UP+", mr.Pattern.Raw)
	}
}

func TestMatchRecognize_OmitEmptyMatches(t *testing.T) {
	// corpus example_05 — ALL ROWS PER MATCH OMIT EMPTY MATCHES, window in DEFINE.
	in := `select * from stock_price_history match_recognize(
    partition by company
    order by price_date
    measures match_number() as "MATCH_NUMBER"
    all rows per match omit empty matches
    pattern(OVERAVG*)
    define
        OVERAVG as price > avg(price) over (rows between unbounded preceding and unbounded following)
)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.RowsPerMatch.Opt != ast.RowsPerMatchOmitEmpty {
		t.Errorf("opt = %v, want OmitEmpty", mr.RowsPerMatch.Opt)
	}
	if mr.Pattern.Raw != "OVERAVG*" {
		t.Errorf("pattern = %q, want OVERAVG*", mr.Pattern.Raw)
	}
}

func TestMatchRecognize_WithUnmatchedRows(t *testing.T) {
	// corpus example_06.
	in := `select * from stock_price_history match_recognize(
    partition by company order by price_date
    measures match_number() as "MATCH_NUMBER", classifier() as cl
    all rows per match with unmatched rows
    pattern(OVERAVG+)
    define OVERAVG as price > 1
)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.RowsPerMatch.Opt != ast.RowsPerMatchWithUnmatched {
		t.Errorf("opt = %v, want WithUnmatched", mr.RowsPerMatch.Opt)
	}
}

func TestMatchRecognize_FinalMeasures_SkipPastLastRow(t *testing.T) {
	// corpus example_07 — FINAL FIRST/LAST measures, SKIP PAST LAST ROW.
	in := `SELECT company FROM stock_price_history
       MATCH_RECOGNIZE (
            PARTITION BY company
            ORDER BY price_date
            MEASURES
                FINAL FIRST(LT45.price) AS "FINAL FIRST(LT45.price)",
                FINAL LAST(LT45.price)  AS "FINAL LAST(LT45.price)"
            ALL ROWS PER MATCH
            AFTER MATCH SKIP PAST LAST ROW
            PATTERN (LT45 LT45)
            DEFINE LT45 AS price < 45.00
       )
    WHERE company = 'ABCD'`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if len(mr.Measures) != 2 {
		t.Fatalf("measures = %d, want 2", len(mr.Measures))
	}
	if mr.Measures[0].Semantics != ast.MatchSemanticsFinal {
		t.Errorf("measure[0] semantics = %v, want FINAL", mr.Measures[0].Semantics)
	}
	if mr.Measures[1].Semantics != ast.MatchSemanticsFinal {
		t.Errorf("measure[1] semantics = %v, want FINAL", mr.Measures[1].Semantics)
	}
	if mr.AfterMatch == nil || mr.AfterMatch.Kind != ast.AfterMatchSkipPastLastRow {
		t.Errorf("after match = %+v, want SKIP PAST LAST ROW", mr.AfterMatch)
	}
	if mr.Pattern.Raw != "LT45 LT45" {
		t.Errorf("pattern = %q, want LT45 LT45", mr.Pattern.Raw)
	}
	// WHERE applies to the SELECT, not consumed by MATCH_RECOGNIZE.
	if sel.Where == nil {
		t.Error("expected WHERE on the outer SELECT")
	}
}

func TestMatchRecognize_SkipToNextRow(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		ORDER BY a
		AFTER MATCH SKIP TO NEXT ROW
		PATTERN(A B)
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.AfterMatch.Kind != ast.AfterMatchSkipToNextRow {
		t.Errorf("after match = %+v, want SKIP TO NEXT ROW", mr.AfterMatch)
	}
}

func TestMatchRecognize_SkipToFirst(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		AFTER MATCH SKIP TO FIRST x
		PATTERN(x+)
		DEFINE x AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.AfterMatch.Kind != ast.AfterMatchSkipToFirst || mr.AfterMatch.Symbol.Name != "x" {
		t.Errorf("after match = %+v, want SKIP TO FIRST x", mr.AfterMatch)
	}
}

func TestMatchRecognize_SkipToBareVar(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		AFTER MATCH SKIP TO myvar
		PATTERN(myvar+)
		DEFINE myvar AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.AfterMatch.Kind != ast.AfterMatchSkipToVar || mr.AfterMatch.Symbol.Name != "myvar" {
		t.Errorf("after match = %+v, want SKIP TO myvar", mr.AfterMatch)
	}
}

func TestMatchRecognize_ShowEmptyMatches(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		ALL ROWS PER MATCH SHOW EMPTY MATCHES
		PATTERN(A+)
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.RowsPerMatch.Opt != ast.RowsPerMatchShowEmpty {
		t.Errorf("opt = %v, want ShowEmpty", mr.RowsPerMatch.Opt)
	}
}

func TestMatchRecognize_NestedPatternParens(t *testing.T) {
	// PATTERN with a grouped sub-pattern and alternation must be captured whole.
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		PATTERN( A (B | C)+ D? )
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.Pattern.Raw != "A (B | C)+ D?" {
		t.Errorf("pattern = %q, want %q", mr.Pattern.Raw, "A (B | C)+ D?")
	}
}

func TestMatchRecognize_TrailingAlias(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(PATTERN(A) DEFINE A AS TRUE) AS mr`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.Alias.Name != "mr" {
		t.Errorf("alias = %q, want mr", mr.Alias.Name)
	}
}

func TestMatchRecognize_MinimalPatternOnly(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(PATTERN(A))`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.Pattern.Raw != "A" {
		t.Errorf("pattern = %q, want A", mr.Pattern.Raw)
	}
	if mr.PartitionBy != nil || mr.OrderBy != nil || mr.Measures != nil ||
		mr.RowsPerMatch != nil || mr.AfterMatch != nil || mr.Define != nil {
		t.Error("expected all optional sub-clauses absent")
	}
}

func TestMatchRecognize_Negatives(t *testing.T) {
	cases := []string{
		`SELECT * FROM t MATCH_RECOGNIZE(PATTERN A)`,                // PATTERN without (
		`SELECT * FROM t MATCH_RECOGNIZE(PATTERN(A`,                 // unterminated PATTERN
		`SELECT * FROM t MATCH_RECOGNIZE(AFTER MATCH PATTERN(A))`,   // AFTER MATCH without SKIP
		`SELECT * FROM t MATCH_RECOGNIZE(ALL PER MATCH PATTERN(A))`, // ALL without ROWS
		`SELECT * FROM t MATCH_RECOGNIZE(DEFINE x PATTERN(A))`,      // DEFINE without AS
		`SELECT * FROM t MATCH_RECOGNIZE(PATTERN(A)`,                // unterminated MR (no closing paren)
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

// PATTERN raw-text Loc must point at the actual inner span.
func TestMatchRecognize_PatternLoc(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(PATTERN(ABC) DEFINE ABC AS TRUE)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	got := in[mr.Pattern.RawLoc.Start:mr.Pattern.RawLoc.End]
	if !strings.Contains(got, "ABC") {
		t.Errorf("RawLoc slice = %q, want to contain ABC", got)
	}
}

func TestMatchRecognize_RunningSemantics(t *testing.T) {
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		MEASURES RUNNING COUNT(*) AS c
		PATTERN(A+)
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if len(mr.Measures) != 1 || mr.Measures[0].Semantics != ast.MatchSemanticsRunning {
		t.Fatalf("measures = %+v, want one RUNNING measure", mr.Measures)
	}
}

func TestMatchRecognize_FinalAsFunctionName(t *testing.T) {
	// FINAL immediately followed by '(' is a function call, not the semantics
	// keyword — the guard must not consume it as a prefix.
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		MEASURES final(price) AS f
		PATTERN(A+)
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.Measures[0].Semantics != ast.MatchSemanticsUnspecified {
		t.Errorf("semantics = %v, want unspecified (final() is a call)", mr.Measures[0].Semantics)
	}
	if _, ok := mr.Measures[0].Expr.(*ast.FuncCallExpr); !ok {
		t.Errorf("measure expr = %T, want FuncCallExpr", mr.Measures[0].Expr)
	}
}

func TestMatchRecognize_RunningAsColumnName(t *testing.T) {
	// "running" used as a bare measure operand (alias follows) must NOT be
	// consumed as the RUNNING semantics keyword.
	in := `SELECT * FROM t MATCH_RECOGNIZE(
		MEASURES running AS r
		PATTERN(A+)
		DEFINE A AS TRUE
	)`
	sel := firstSelect(t, in)
	mr := firstTableRef(t, sel).MatchRecognize
	if mr.Measures[0].Semantics != ast.MatchSemanticsUnspecified {
		t.Errorf("semantics = %v, want unspecified", mr.Measures[0].Semantics)
	}
	cr, ok := mr.Measures[0].Expr.(*ast.ColumnRef)
	if !ok || cr.Parts[0].Name != "running" {
		t.Errorf("measure expr = %T, want ColumnRef{running}", mr.Measures[0].Expr)
	}
	if mr.Measures[0].Alias.Name != "r" {
		t.Errorf("alias = %q, want r", mr.Measures[0].Alias.Name)
	}
}
