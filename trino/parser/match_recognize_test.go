package parser

import (
	"strconv"
	"testing"
)

// This file is the parser-match-recognize node's structural slice: the oracle
// differential (oracle_match_recognize_test.go) pins accept/reject, while these
// tests pin the AST shapes accept/reject cannot see — row-pattern precedence and
// associativity (R1), the at-most-one-quantifier model (R2), the quantifier
// bounds, the double-alias split (M3), and the window-frame pattern attachment
// (M6). Helpers parseOneQuery / querySpec live in select_test.go.

// patternRecognitionFrom parses a single `SELECT * FROM <fromText>` and returns
// the sole FROM relation as a *PatternRecognition, failing if the shape differs.
func patternRecognitionFrom(t *testing.T, fromText string) *PatternRecognition {
	t.Helper()
	spec := querySpec(t, parseOneQuery(t, "SELECT * FROM "+fromText))
	if len(spec.From) != 1 {
		t.Fatalf("want 1 FROM relation, got %d", len(spec.From))
	}
	pr, ok := spec.From[0].(*PatternRecognition)
	if !ok {
		t.Fatalf("FROM relation is %T, want *PatternRecognition", spec.From[0])
	}
	return pr
}

// rowPatternOf parses a `MATCH_RECOGNIZE (PATTERN ( <patternText> ) DEFINE A AS
// true)` and returns the parsed PATTERN tree.
func rowPatternOf(t *testing.T, patternText string) RowPattern {
	t.Helper()
	pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE (PATTERN ("+patternText+") DEFINE A AS true)")
	if pr.Pattern == nil {
		t.Fatal("PATTERN tree is nil")
	}
	return pr.Pattern
}

func TestMatchRecognize_StructureMandatoryFields(t *testing.T) {
	pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE (PATTERN (A B) DEFINE A AS true, B AS false)")
	// Pattern must be present, Definitions must list both A and B.
	if pr.Pattern == nil {
		t.Fatal("Pattern is nil")
	}
	if len(pr.Definitions) != 2 {
		t.Fatalf("want 2 DEFINE entries, got %d", len(pr.Definitions))
	}
	if pr.Definitions[0].Name.Value != "A" || pr.Definitions[1].Name.Value != "B" {
		t.Errorf("DEFINE names = %q,%q, want A,B", pr.Definitions[0].Name.Value, pr.Definitions[1].Name.Value)
	}
	// The wrapped input relation is the table t.
	ar, ok := pr.Input.(*AliasedRelation)
	if !ok {
		t.Fatalf("Input is %T, want *AliasedRelation", pr.Input)
	}
	if _, ok := ar.Inner.(*TableRelation); !ok {
		t.Errorf("Input.Inner is %T, want *TableRelation", ar.Inner)
	}
}

func TestMatchRecognize_StructureFullClauses(t *testing.T) {
	pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE ("+
		"PARTITION BY a, b ORDER BY c MEASURES x AS m, y AS n "+
		"ALL ROWS PER MATCH WITH UNMATCHED ROWS "+
		"AFTER MATCH SKIP TO FIRST B INITIAL "+
		"PATTERN (A B) SUBSET U = (A, B) DEFINE A AS true, B AS false)")

	if len(pr.PartitionBy) != 2 {
		t.Errorf("PARTITION BY count = %d, want 2", len(pr.PartitionBy))
	}
	if len(pr.OrderBy) != 1 {
		t.Errorf("ORDER BY count = %d, want 1", len(pr.OrderBy))
	}
	if len(pr.Measures) != 2 || pr.Measures[0].Name.Value != "m" || pr.Measures[1].Name.Value != "n" {
		t.Errorf("MEASURES = %+v, want m,n", pr.Measures)
	}
	if pr.RowsPerMatch == nil || pr.RowsPerMatch.Kind != AllRowsPerMatch ||
		pr.RowsPerMatch.EmptyHandling != "WITH UNMATCHED ROWS" {
		t.Errorf("RowsPerMatch = %+v, want ALL ROWS / WITH UNMATCHED ROWS", pr.RowsPerMatch)
	}
	if pr.AfterMatchSkip == nil || pr.AfterMatchSkip.Kind != SkipToFirst ||
		pr.AfterMatchSkip.Variable == nil || pr.AfterMatchSkip.Variable.Value != "B" {
		t.Errorf("AfterMatchSkip = %+v, want SKIP TO FIRST B", pr.AfterMatchSkip)
	}
	if pr.SearchMode != "INITIAL" {
		t.Errorf("SearchMode = %q, want INITIAL", pr.SearchMode)
	}
	if len(pr.Subsets) != 1 || pr.Subsets[0].Name.Value != "U" || len(pr.Subsets[0].Of) != 2 {
		t.Errorf("Subsets = %+v, want U = (A, B)", pr.Subsets)
	}
}

func TestMatchRecognize_StructureRowsPerMatchVariants(t *testing.T) {
	cases := []struct {
		text     string
		wantKind RowsPerMatchKind
		wantEmpt string
	}{
		{"ONE ROW PER MATCH", OneRowPerMatch, ""},
		{"ALL ROWS PER MATCH", AllRowsPerMatch, ""},
		{"ALL ROWS PER MATCH SHOW EMPTY MATCHES", AllRowsPerMatch, "SHOW EMPTY MATCHES"},
		{"ALL ROWS PER MATCH OMIT EMPTY MATCHES", AllRowsPerMatch, "OMIT EMPTY MATCHES"},
		{"ALL ROWS PER MATCH WITH UNMATCHED ROWS", AllRowsPerMatch, "WITH UNMATCHED ROWS"},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE ("+c.text+" PATTERN (A) DEFINE A AS true)")
			if pr.RowsPerMatch == nil {
				t.Fatal("RowsPerMatch is nil")
			}
			if pr.RowsPerMatch.Kind != c.wantKind || pr.RowsPerMatch.EmptyHandling != c.wantEmpt {
				t.Errorf("got kind=%v empty=%q, want kind=%v empty=%q",
					pr.RowsPerMatch.Kind, pr.RowsPerMatch.EmptyHandling, c.wantKind, c.wantEmpt)
			}
		})
	}
}

func TestMatchRecognize_StructureSkipToVariants(t *testing.T) {
	cases := []struct {
		text     string
		wantKind SkipToKind
		wantVar  string
	}{
		{"AFTER MATCH SKIP PAST LAST ROW", SkipPastLastRow, ""},
		{"AFTER MATCH SKIP TO NEXT ROW", SkipToNextRow, ""},
		{"AFTER MATCH SKIP TO FIRST A", SkipToFirst, "A"},
		{"AFTER MATCH SKIP TO LAST A", SkipToLast, "A"},
		{"AFTER MATCH SKIP TO A", SkipToVariable, "A"},
		// non-reserved NEXT/FIRST/LAST used as the variable name itself (the special
		// form does not apply because the disambiguating token does not follow).
		{"AFTER MATCH SKIP TO NEXT", SkipToVariable, "NEXT"},   // NEXT not followed by ROW
		{"AFTER MATCH SKIP TO FIRST", SkipToVariable, "FIRST"}, // FIRST not followed by ident+successor
		{"AFTER MATCH SKIP TO LAST", SkipToVariable, "LAST"},   // LAST not followed by ident+successor
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE ("+c.text+" PATTERN (A) DEFINE A AS true)")
			if pr.AfterMatchSkip == nil {
				t.Fatal("AfterMatchSkip is nil")
			}
			if pr.AfterMatchSkip.Kind != c.wantKind {
				t.Errorf("kind = %v, want %v", pr.AfterMatchSkip.Kind, c.wantKind)
			}
			gotVar := ""
			if pr.AfterMatchSkip.Variable != nil {
				gotVar = pr.AfterMatchSkip.Variable.Value
			}
			if gotVar != c.wantVar {
				t.Errorf("var = %q, want %q", gotVar, c.wantVar)
			}
		})
	}
}

func TestMatchRecognize_StructureDoubleAlias(t *testing.T) {
	// `t AS r MATCH_RECOGNIZE (…) AS m`: the inner alias r is on the wrapped
	// relation; the trailing alias m is on the PatternRecognition node (M3).
	pr := patternRecognitionFrom(t, "t AS r MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) AS m (x, y)")
	if pr.Alias == nil || pr.Alias.Value != "m" {
		t.Errorf("trailing alias = %v, want m", pr.Alias)
	}
	if len(pr.ColumnAliases) != 2 || pr.ColumnAliases[0].Value != "x" || pr.ColumnAliases[1].Value != "y" {
		t.Errorf("trailing column aliases = %+v, want x,y", pr.ColumnAliases)
	}
	ar, ok := pr.Input.(*AliasedRelation)
	if !ok {
		t.Fatalf("Input is %T, want *AliasedRelation", pr.Input)
	}
	if ar.Alias == nil || ar.Alias.Value != "r" {
		t.Errorf("inner alias = %v, want r", ar.Alias)
	}
}

// TestMatchRecognize_TrailingAliasNamedMatchRecognize pins the edge where the
// trailing alias is itself the non-reserved keyword MATCH_RECOGNIZE with a column
// list (no second clause can follow the MATCH_RECOGNIZE clause, so it is an alias).
func TestMatchRecognize_TrailingAliasNamedMatchRecognize(t *testing.T) {
	pr := patternRecognitionFrom(t, "t MATCH_RECOGNIZE (PATTERN (A) DEFINE A AS true) MATCH_RECOGNIZE (x)")
	if pr.Alias == nil || pr.Alias.Value != "MATCH_RECOGNIZE" {
		t.Errorf("trailing alias = %v, want MATCH_RECOGNIZE", pr.Alias)
	}
	if len(pr.ColumnAliases) != 1 || pr.ColumnAliases[0].Value != "x" {
		t.Errorf("trailing column aliases = %+v, want [x]", pr.ColumnAliases)
	}
}

// TestMatchRecognize_WindowFrameSpan verifies the WindowFrame source span covers
// the whole frame including any trailing row-pattern clauses, not just the
// frameExtent (so downstream span consumers see the full range).
func TestMatchRecognize_WindowFrameSpan(t *testing.T) {
	const src = "SELECT count(*) OVER (ROWS CURRENT ROW PATTERN (A) DEFINE A AS true) FROM t"
	spec := querySpec(t, parseOneQuery(t, src))
	frame := spec.Items[0].Expr.(*FuncCall).Over.Frame
	if frame == nil {
		t.Fatal("frame is nil")
	}
	covered := src[frame.Loc.Start:frame.Loc.End]
	if covered != "ROWS CURRENT ROW PATTERN (A) DEFINE A AS true" {
		t.Errorf("frame span = %q, want it to cover through the DEFINE clause", covered)
	}
}

// TestMatchRecognize_RowPatternConcatLeftAssoc verifies `A B C` nests
// left-associatively: PatternConcat{PatternConcat{A,B}, C} (R1).
func TestMatchRecognize_RowPatternConcatLeftAssoc(t *testing.T) {
	root, ok := rowPatternOf(t, "A B C").(*PatternConcat)
	if !ok {
		t.Fatalf("root is %T, want *PatternConcat", rowPatternOf(t, "A B C"))
	}
	if patternVarName(t, root.Right) != "C" {
		t.Errorf("root.Right = %v, want variable C", root.Right)
	}
	left, ok := root.Left.(*PatternConcat)
	if !ok {
		t.Fatalf("root.Left is %T, want *PatternConcat", root.Left)
	}
	if patternVarName(t, left.Left) != "A" || patternVarName(t, left.Right) != "B" {
		t.Errorf("left = (%v, %v), want (A, B)", left.Left, left.Right)
	}
}

// TestMatchRecognize_RowPatternAltBindsLooser verifies concatenation binds
// tighter than `|`: `A B | C D` is PatternAlternation{(A B), (C D)} (R1).
func TestMatchRecognize_RowPatternAltBindsLooser(t *testing.T) {
	root, ok := rowPatternOf(t, "A B | C D").(*PatternAlternation)
	if !ok {
		t.Fatalf("root is %T, want *PatternAlternation", rowPatternOf(t, "A B | C D"))
	}
	if _, ok := root.Left.(*PatternConcat); !ok {
		t.Errorf("root.Left is %T, want *PatternConcat (A B)", root.Left)
	}
	if _, ok := root.Right.(*PatternConcat); !ok {
		t.Errorf("root.Right is %T, want *PatternConcat (C D)", root.Right)
	}
}

// TestMatchRecognize_RowPatternAltLeftAssoc verifies `A | B | C` nests
// left-associatively: PatternAlternation{PatternAlternation{A,B}, C} (R1).
func TestMatchRecognize_RowPatternAltLeftAssoc(t *testing.T) {
	root, ok := rowPatternOf(t, "A | B | C").(*PatternAlternation)
	if !ok {
		t.Fatalf("root is %T, want *PatternAlternation", rowPatternOf(t, "A | B | C"))
	}
	if patternVarName(t, root.Right) != "C" {
		t.Errorf("root.Right = %v, want C", root.Right)
	}
	if _, ok := root.Left.(*PatternAlternation); !ok {
		t.Errorf("root.Left is %T, want *PatternAlternation (A | B)", root.Left)
	}
}

func TestMatchRecognize_RowPatternQuantifiers(t *testing.T) {
	mk := func(v int64) *int64 { return &v }
	cases := []struct {
		text      string
		wantKind  PatternQuantifierKind
		reluctant bool
		exactly   *int64
		atLeast   *int64
		atMost    *int64
	}{
		{"A*", QuantZeroOrMore, false, nil, nil, nil},
		{"A+", QuantOneOrMore, false, nil, nil, nil},
		{"A?", QuantZeroOrOne, false, nil, nil, nil},
		{"A*?", QuantZeroOrMore, true, nil, nil, nil},
		{"A+?", QuantOneOrMore, true, nil, nil, nil},
		{"A??", QuantZeroOrOne, true, nil, nil, nil},
		{"A{3}", QuantRange, false, mk(3), nil, nil},
		{"A{1,3}", QuantRange, false, nil, mk(1), mk(3)},
		{"A{,3}", QuantRange, false, nil, nil, mk(3)},
		{"A{2,}", QuantRange, false, nil, mk(2), nil},
		{"A{2,}?", QuantRange, true, nil, mk(2), nil},
		{"A{,}", QuantRange, false, nil, nil, nil},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			qp, ok := rowPatternOf(t, c.text).(*QuantifiedPattern)
			if !ok {
				t.Fatalf("root is %T, want *QuantifiedPattern", rowPatternOf(t, c.text))
			}
			q := qp.Quantifier
			if q == nil {
				t.Fatal("Quantifier is nil")
			}
			if q.Kind != c.wantKind || q.Reluctant != c.reluctant {
				t.Errorf("kind=%v reluctant=%v, want kind=%v reluctant=%v", q.Kind, q.Reluctant, c.wantKind, c.reluctant)
			}
			if !eqIntPtr(q.Exactly, c.exactly) || !eqIntPtr(q.AtLeast, c.atLeast) || !eqIntPtr(q.AtMost, c.atMost) {
				t.Errorf("bounds exactly=%v atLeast=%v atMost=%v, want exactly=%v atLeast=%v atMost=%v",
					derefStr(q.Exactly), derefStr(q.AtLeast), derefStr(q.AtMost),
					derefStr(c.exactly), derefStr(c.atLeast), derefStr(c.atMost))
			}
		})
	}
}

func TestMatchRecognize_RowPatternPrimaries(t *testing.T) {
	// Empty pattern ().
	if _, ok := rowPatternOf(t, "()").(*EmptyPattern); !ok {
		t.Errorf("() is %T, want *EmptyPattern", rowPatternOf(t, "()"))
	}
	// Grouped pattern (A B).
	if g, ok := rowPatternOf(t, "(A B)").(*GroupedPattern); !ok {
		t.Errorf("(A B) is %T, want *GroupedPattern", rowPatternOf(t, "(A B)"))
	} else if _, ok := g.Inner.(*PatternConcat); !ok {
		t.Errorf("grouped inner is %T, want *PatternConcat", g.Inner)
	}
	// Start / end anchors.
	if a, ok := rowPatternOf(t, "^").(*AnchorPattern); !ok || !a.Start {
		t.Errorf("^ is %T (start=%v), want *AnchorPattern start=true", rowPatternOf(t, "^"), ok && a.Start)
	}
	if a, ok := rowPatternOf(t, "$").(*AnchorPattern); !ok || a.Start {
		t.Errorf("$ is %T, want *AnchorPattern start=false", rowPatternOf(t, "$"))
	}
	// Excluded pattern {- A -}.
	if e, ok := rowPatternOf(t, "{- A -}").(*ExcludedPattern); !ok {
		t.Errorf("{- A -} is %T, want *ExcludedPattern", rowPatternOf(t, "{- A -}"))
	} else if patternVarName(t, e.Inner) != "A" {
		t.Errorf("excluded inner = %v, want A", e.Inner)
	}
	// PERMUTE(A, B).
	if perm, ok := rowPatternOf(t, "PERMUTE(A, B)").(*PatternPermutation); !ok {
		t.Errorf("PERMUTE(A, B) is %T, want *PatternPermutation", rowPatternOf(t, "PERMUTE(A, B)"))
	} else if len(perm.Patterns) != 2 {
		t.Errorf("PERMUTE operand count = %d, want 2", len(perm.Patterns))
	}
	// Empty PERMUTE() — accepted (divergence #92), zero operands.
	if perm, ok := rowPatternOf(t, "PERMUTE()").(*PatternPermutation); !ok {
		t.Errorf("PERMUTE() is %T, want *PatternPermutation", rowPatternOf(t, "PERMUTE()"))
	} else if len(perm.Patterns) != 0 {
		t.Errorf("PERMUTE() operand count = %d, want 0", len(perm.Patterns))
	}
	// PERMUTE is non-reserved: a bare `permute` (no '(') is a pattern variable.
	if pv, ok := rowPatternOf(t, "permute").(*PatternVariable); !ok || pv.Name.Value != "permute" {
		t.Errorf("bare permute is %T, want *PatternVariable(permute)", rowPatternOf(t, "permute"))
	}
}

// TestMatchRecognize_WindowMeasuresVsName pins the MEASURES-clause vs
// existing-window-name disambiguation in an OVER frame (the >1-token lookahead
// case Codex flagged): MEASURES followed by `expr AS id` is the clause; MEASURES
// followed by a window-spec part (here a ROWS frame) is the window name.
func TestMatchRecognize_WindowMeasuresVsName(t *testing.T) {
	// MEASURES clause whose first measure column is the non-reserved keyword "rows".
	spec := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (MEASURES rows AS m ROWS CURRENT ROW PATTERN (A) DEFINE A AS true) FROM t"))
	fc := spec.Items[0].Expr.(*FuncCall)
	if fc.Over == nil || fc.Over.ExistingName != nil {
		t.Errorf("MEASURES clause: ExistingName = %v, want nil", fc.Over.ExistingName)
	}
	if fc.Over.Frame == nil || fc.Over.Frame.Pattern == nil || len(fc.Over.Frame.Pattern.Measures) != 1 {
		t.Fatalf("MEASURES clause: frame.Pattern.Measures wrong: %+v", fc.Over.Frame)
	}
	if fc.Over.Frame.Pattern.Measures[0].Name.Value != "m" {
		t.Errorf("measure name = %v, want m", fc.Over.Frame.Pattern.Measures[0].Name.Value)
	}

	// A window literally named "MEASURES" + a ROWS frame: MEASURES is the existing
	// window name, the frame is ordinary (no pattern additions).
	spec2 := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (MEASURES ROWS CURRENT ROW) FROM t WINDOW MEASURES AS (PARTITION BY a)"))
	fc2 := spec2.Items[0].Expr.(*FuncCall)
	if fc2.Over == nil || fc2.Over.ExistingName == nil || fc2.Over.ExistingName.Value != "MEASURES" {
		t.Errorf("window named MEASURES: ExistingName = %v, want MEASURES", fc2.Over.ExistingName)
	}
	if fc2.Over.Frame == nil || fc2.Over.Frame.Pattern != nil {
		t.Errorf("window named MEASURES: frame should be ordinary (Pattern nil), got %+v", fc2.Over.Frame)
	}
}

// TestMatchRecognize_WindowFramePattern verifies the OVER-clause window frame
// carries the row-pattern additions on WindowFrame.Pattern (M6).
func TestMatchRecognize_WindowFramePattern(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (MEASURES x AS m ROWS CURRENT ROW AFTER MATCH SKIP PAST LAST ROW "+
			"PATTERN (A B) SUBSET U = (A, B) DEFINE A AS true) FROM t"))
	fc, ok := spec.Items[0].Expr.(*FuncCall)
	if !ok {
		t.Fatalf("select item is %T, want *FuncCall", spec.Items[0].Expr)
	}
	if fc.Over == nil || fc.Over.Frame == nil {
		t.Fatal("OVER frame is nil")
	}
	mrf := fc.Over.Frame.Pattern
	if mrf == nil {
		t.Fatal("frame.Pattern (MatchRecognizeFrame) is nil")
	}
	if len(mrf.Measures) != 1 || mrf.Measures[0].Name.Value != "m" {
		t.Errorf("frame MEASURES = %+v, want m", mrf.Measures)
	}
	if mrf.AfterMatchSkip == nil || mrf.AfterMatchSkip.Kind != SkipPastLastRow {
		t.Errorf("frame AFTER MATCH = %+v, want SKIP PAST LAST ROW", mrf.AfterMatchSkip)
	}
	if mrf.Pattern == nil {
		t.Error("frame PATTERN is nil")
	}
	if len(mrf.Subsets) != 1 {
		t.Errorf("frame SUBSET count = %d, want 1", len(mrf.Subsets))
	}
	if len(mrf.Definitions) != 1 {
		t.Errorf("frame DEFINE count = %d, want 1", len(mrf.Definitions))
	}
}

// TestMatchRecognize_WindowFrameSkipFollowers verifies the FIRST/LAST skip
// keyword disambiguation in a window frame, where the skipTo target may be
// followed only by the frame-closing ')' or by an optional SUBSET/DEFINE (no
// PATTERN). FIRST/LAST is still the keyword (with the bare label as its variable).
func TestMatchRecognize_WindowFrameSkipFollowers(t *testing.T) {
	// skipTo target ends the frame (next token is ')').
	spec := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP TO FIRST A) FROM t"))
	skip := spec.Items[0].Expr.(*FuncCall).Over.Frame.Pattern.AfterMatchSkip
	if skip == nil || skip.Kind != SkipToFirst || skip.Variable == nil || skip.Variable.Value != "A" {
		t.Errorf("window SKIP TO FIRST A): got %+v, want SkipToFirst(A)", skip)
	}
	// skipTo target followed by SUBSET/DEFINE (no PATTERN) in the frame.
	spec2 := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (ROWS CURRENT ROW AFTER MATCH SKIP TO LAST A SUBSET U = (A) DEFINE A AS true) FROM t"))
	mrf := spec2.Items[0].Expr.(*FuncCall).Over.Frame.Pattern
	if mrf.AfterMatchSkip == nil || mrf.AfterMatchSkip.Kind != SkipToLast || mrf.AfterMatchSkip.Variable.Value != "A" {
		t.Errorf("window SKIP TO LAST A SUBSET: got %+v, want SkipToLast(A)", mrf.AfterMatchSkip)
	}
	if mrf.Pattern != nil {
		t.Errorf("window frame with no PATTERN: mrf.Pattern = %+v, want nil", mrf.Pattern)
	}
	if len(mrf.Subsets) != 1 || len(mrf.Definitions) != 1 {
		t.Errorf("frame SUBSET/DEFINE counts = %d/%d, want 1/1", len(mrf.Subsets), len(mrf.Definitions))
	}
}

// TestMatchRecognize_OrdinaryFrameNoPattern verifies an ordinary window frame
// (no pattern recognition) leaves WindowFrame.Pattern nil.
func TestMatchRecognize_OrdinaryFrameNoPattern(t *testing.T) {
	spec := querySpec(t, parseOneQuery(t,
		"SELECT count(*) OVER (ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t"))
	fc := spec.Items[0].Expr.(*FuncCall)
	if fc.Over == nil || fc.Over.Frame == nil {
		t.Fatal("OVER frame is nil")
	}
	if fc.Over.Frame.Pattern != nil {
		t.Errorf("ordinary frame Pattern = %+v, want nil", fc.Over.Frame.Pattern)
	}
}

// --- small helpers ---

func patternVarName(t *testing.T, p RowPattern) string {
	t.Helper()
	pv, ok := p.(*PatternVariable)
	if !ok {
		t.Fatalf("pattern is %T, want *PatternVariable", p)
	}
	return pv.Name.Value
}

func eqIntPtr(a, b *int64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func derefStr(p *int64) string {
	if p == nil {
		return "<nil>"
	}
	return strconv.FormatInt(*p, 10)
}
