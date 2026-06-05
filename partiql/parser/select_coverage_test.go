package parser

import (
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// select_coverage_test.go is a TEST-ONLY negative/coverage backfill for the
// SELECT query forms in select.go / from.go / expr.go (parseBagOp). It pins
// implemented-but-previously-untested syntax: the SELECT ALL quantifier, the
// RIGHT/FULL/FULL OUTER/bare-OUTER/bare-JOIN join types, GROUP PARTIAL BY,
// GROUP AS, group-key AS, the set-op modifiers (UNION DISTINCT/ALL, EXCEPT
// ALL, INTERSECT DISTINCT, OUTER UNION/EXCEPT/INTERSECT), bare projection
// aliases, and HAVING — plus malformed variants the grammar rejects.
//
// Oracle (every expected accept/reject below is derived from it, not from the
// author's judgment):
//   - truth2: the EXECUTABLE generated ANTLR PartiQL parser at
//     github.com/bytebase/parser/partiql. "accept" == 0 lexer+parser syntax
//     errors on Script(); "reject" == >0. Verdicts and the exact accepted
//     AST shapes were taken from a throwaway differential probe that ran both
//     the ANTLR parser and omni Parse() over each input.
//   - truth1: the PartiQL spec / AWS DynamoDB PartiQL reference, cross-checked
//     against PartiQLParser.g4 (cited inline by rule name + line).
//
// Accept cases assert the full AST via ast.NodeToString (golden copied
// verbatim from the probe). Reject cases assert only that ParseStatement
// errors (the oracle decides accept/reject; the message is implementation-
// defined). ParseStatement asserts EOF, so a dangling/unconsumed token fails.

// TestSelectCoverage_SetQuantifier covers setQuantifierStrategy
// (PartiQLParser.g4:229-232 `DISTINCT | ALL`) on the SELECT clause, including
// the previously-untested `SELECT ALL` form for both projection lists and
// `SELECT ALL *`. Oracle: ANTLR accepts ALL and DISTINCT; it rejects stacking
// two quantifiers (the rule allows at most one `setQuantifierStrategy?`).
func TestSelectCoverage_SetQuantifier(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "all_items",
			input: "SELECT ALL x FROM t",
			want:  "SelectStmt{Quantifier:ALL Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}",
		},
		{
			name:  "all_star",
			input: "SELECT ALL * FROM t",
			want:  "SelectStmt{Quantifier:ALL Star:true From:VarRef{Name:t}}",
		},
		{
			// Positive control: DISTINCT, the sibling quantifier.
			name:  "distinct_items",
			input: "SELECT DISTINCT x FROM t",
			want:  "SelectStmt{Quantifier:DISTINCT Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_SetQuantifierReject covers the rejects for
// setQuantifierStrategy: at most one quantifier may appear. Oracle: ANTLR
// rejects both stacking orders.
func TestSelectCoverage_SetQuantifierReject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"all_distinct", "SELECT ALL DISTINCT x FROM t"},
		{"distinct_all", "SELECT DISTINCT ALL x FROM t"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runReject(t, tc.input)
		})
	}
}

// TestSelectCoverage_JoinTypes covers the join types in joinType
// (PartiQLParser.g4:419-425) that the existing from_join_test.go suite does
// NOT already pin: RIGHT [OUTER] JOIN, FULL [OUTER] JOIN, bare OUTER JOIN, and
// bare JOIN. Oracle: ANTLR accepts all; omni maps the modifier to the matching
// ast.JoinKind (RIGHT/FULL/OUTER/INNER). A qualified JOIN requires the ON
// joinSpec (PartiQLParser.g4:392,416), which these inputs supply.
func TestSelectCoverage_JoinTypes(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "right_join",
			input: "SELECT * FROM t1 RIGHT JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:RIGHT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "right_outer_join",
			input: "SELECT * FROM t1 RIGHT OUTER JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:RIGHT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "full_join",
			input: "SELECT * FROM t1 FULL JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:FULL Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "full_outer_join",
			input: "SELECT * FROM t1 FULL OUTER JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:FULL Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			// joinType#mod=OUTER (g4:424): a bare OUTER JOIN. PartiQL-specific
			// natural-outer form; rendered with ast.JoinKindOuter.
			name:  "bare_outer_join",
			input: "SELECT * FROM t1 OUTER JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:OUTER Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			// tableReference#TableQualifiedJoin with no joinType (g4:392):
			// a bare JOIN defaults to INNER.
			name:  "bare_join",
			input: "SELECT * FROM t1 JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:INNER Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_ModifierCrossJoin pins the tableReference#TableCrossJoin
// form `lhs joinType? CROSS JOIN rhs` (PartiQLParser.g4:390) where a join
// modifier PRECEDES CROSS JOIN, e.g. "FULL CROSS JOIN".
//
// DIVERGENCE (residual omni bug — see new_findings): the ANTLR oracle ACCEPTS
// every one of these (the parse tree shows a `(joinType FULL) CROSS JOIN`
// node with zero errors), but omni Parse() REJECTS them all. omni's
// tryParseJoinType (from.go:438-509) only recognizes a bare `CROSS JOIN`; a
// modifier followed by CROSS is rewound by the dangling-modifier guard, the
// join loop breaks, and the leftover `FULL CROSS JOIN t2` fails at the EOF
// check ("unexpected token after statement"). Per the TEST-ONLY honesty rule
// these cases are skipped (not weakened) so the suite stays green and the
// divergence is reported, not silently encoded as correct-reject.
func TestSelectCoverage_ModifierCrossJoin(t *testing.T) {
	// Oracle verdict for every case below is ACCEPT (ANTLR, 0 errors).
	cases := []struct {
		name  string
		input string
	}{
		{"full_cross", "SELECT a FROM t FULL CROSS JOIN t2"},
		{"right_cross", "SELECT a FROM t RIGHT CROSS JOIN t2"},
		{"left_cross", "SELECT a FROM t LEFT CROSS JOIN t2"},
		{"inner_cross", "SELECT a FROM t INNER CROSS JOIN t2"},
		{"outer_cross", "SELECT a FROM t OUTER CROSS JOIN t2"},
		{"left_outer_cross", "SELECT a FROM t LEFT OUTER CROSS JOIN t2"},
		{"right_outer_cross", "SELECT a FROM t RIGHT OUTER CROSS JOIN t2"},
		{"full_outer_cross", "SELECT a FROM t FULL OUTER CROSS JOIN t2"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Skip("DIVERGENCE: omni rejects `<joinType> CROSS JOIN` " +
				"(tryParseJoinType in from.go handles only bare CROSS JOIN); " +
				"ANTLR/PartiQLParser.g4:390 TableCrossJoin `lhs joinType? CROSS JOIN rhs` accepts it. " +
				"Residual omni bug; input: " + tc.input)
		})
	}
}

// TestSelectCoverage_CrossJoinSpec pins the boundary of the CROSS JOIN form:
// tableReference#TableCrossJoin (PartiQLParser.g4:390) has NO joinSpec, so a
// trailing ON is illegal. Oracle: ANTLR accepts bare `CROSS JOIN` and rejects
// `CROSS JOIN ... ON ...`.
func TestSelectCoverage_CrossJoinSpec(t *testing.T) {
	t.Run("cross_join_accept", func(t *testing.T) {
		runAcceptGolden(t, "SELECT a FROM t CROSS JOIN t2",
			"SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:JoinExpr{Kind:CROSS Left:VarRef{Name:t} Right:VarRef{Name:t2}}}")
	})
	t.Run("cross_join_on_reject", func(t *testing.T) {
		runReject(t, "SELECT a FROM t CROSS JOIN t2 ON a=b")
	})
}

// TestSelectCoverage_GroupBy covers the groupClause forms
// (PartiQLParser.g4:262-269): GROUP PARTIAL BY, the GROUP AS alias
// (groupAlias g4:265-266), and a per-key AS (groupKey g4:268-269). Oracle:
// ANTLR accepts all; omni records Partial, GroupAs, and the per-item Alias.
func TestSelectCoverage_GroupBy(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "group_partial_by",
			input: "SELECT a FROM t GROUP PARTIAL BY a",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Partial:true Items:[GroupByItem{Expr:VarRef{Name:a}}]}}",
		},
		{
			name:  "group_as",
			input: "SELECT a FROM t GROUP BY a GROUP AS g",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Items:[GroupByItem{Expr:VarRef{Name:a}}] GroupAs:g}}",
		},
		{
			name:  "group_key_as",
			input: "SELECT a FROM t GROUP BY a AS k",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Items:[GroupByItem{Expr:VarRef{Name:a} Alias:k}]}}",
		},
		{
			// All three features at once, plus a multi-key list.
			name:  "group_partial_key_as_multi_group_as",
			input: "SELECT a FROM t GROUP PARTIAL BY a AS k, b GROUP AS g",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Partial:true Items:[GroupByItem{Expr:VarRef{Name:a} Alias:k} GroupByItem{Expr:VarRef{Name:b}}] GroupAs:g}}",
		},
		{
			// Positive control: a plain multi-key GROUP BY.
			name:  "group_by_multi",
			input: "SELECT a FROM t GROUP BY a, b",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Items:[GroupByItem{Expr:VarRef{Name:a}} GroupByItem{Expr:VarRef{Name:b}}]}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_GroupByReject covers malformed groupClause variants.
// Oracle: ANTLR rejects each — GROUP requires PARTIAL? BY (g4:263), GROUP BY
// requires at least one groupKey, GROUP AS requires both AS and the alias
// symbol (groupAlias g4:266), and a per-key AS requires the alias symbol.
func TestSelectCoverage_GroupByReject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"group_no_by", "SELECT a FROM t GROUP a"},
		{"group_by_no_key", "SELECT a FROM t GROUP BY"},
		{"group_partial_no_by", "SELECT a FROM t GROUP PARTIAL a"},
		{"group_as_no_as_kw", "SELECT a FROM t GROUP BY a GROUP g"},
		{"group_as_no_symbol", "SELECT a FROM t GROUP BY a GROUP AS"},
		{"group_key_as_no_symbol", "SELECT a FROM t GROUP BY a AS"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runReject(t, tc.input)
		})
	}
}

// TestSelectCoverage_SetOps covers exprBagOp (PartiQLParser.g4:449-453) set-op
// modifiers that were previously untested: the (DISTINCT|ALL)? quantifier on
// UNION/EXCEPT/INTERSECT and the OUTER? prefix (OUTER UNION/EXCEPT/INTERSECT).
// Oracle: ANTLR accepts all; omni builds a SetOpStmt carrying Op, Quantifier,
// and Outer. Note the operands render as bare SelectStmt (the SetOpStmt's Left
// /Right are StmtNodes), matching exprToStmt's unwrapping in expr.go.
func TestSelectCoverage_SetOps(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "union_distinct",
			input: "SELECT a FROM t UNION DISTINCT SELECT a FROM s",
			want:  "SetOpStmt{Op:UNION Quantifier:DISTINCT Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "union_all",
			input: "SELECT a FROM t UNION ALL SELECT a FROM s",
			want:  "SetOpStmt{Op:UNION Quantifier:ALL Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "except_all",
			input: "SELECT a FROM t EXCEPT ALL SELECT a FROM s",
			want:  "SetOpStmt{Op:EXCEPT Quantifier:ALL Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "intersect_distinct",
			input: "SELECT a FROM t INTERSECT DISTINCT SELECT a FROM s",
			want:  "SetOpStmt{Op:INTERSECT Quantifier:DISTINCT Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "outer_union",
			input: "SELECT a FROM t OUTER UNION SELECT a FROM s",
			want:  "SetOpStmt{Op:UNION Outer:true Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "outer_except",
			input: "SELECT a FROM t OUTER EXCEPT SELECT a FROM s",
			want:  "SetOpStmt{Op:EXCEPT Outer:true Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			name:  "outer_intersect",
			input: "SELECT a FROM t OUTER INTERSECT SELECT a FROM s",
			want:  "SetOpStmt{Op:INTERSECT Outer:true Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
		{
			// Positive control: a plain UNION (no quantifier, no OUTER).
			name:  "union_plain",
			input: "SELECT a FROM t UNION SELECT a FROM s",
			want:  "SetOpStmt{Op:UNION Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:s}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_SetOpReject covers malformed set-op modifiers. Oracle:
// ANTLR rejects stacking DISTINCT and ALL in either order — exprBagOp allows
// at most one `(DISTINCT|ALL)?` (g4:450-452) — and rejects a bare OUTER prefix
// that is not followed by a set-op keyword.
func TestSelectCoverage_SetOpReject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"union_distinct_all", "SELECT a FROM t UNION DISTINCT ALL SELECT a FROM s"},
		{"union_all_distinct", "SELECT a FROM t UNION ALL DISTINCT SELECT a FROM s"},
		// OUTER with no following UNION/EXCEPT/INTERSECT is not a set-op; the
		// trailing SELECT cannot attach, so the statement fails at EOF.
		{"outer_no_setop", "SELECT a FROM t OUTER SELECT a FROM s"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runReject(t, tc.input)
		})
	}
}

// TestSelectCoverage_ProjectionAlias covers projectionItem aliasing
// (PartiQLParser.g4:226-227 `expr (AS? symbolPrimitive)?`): both the explicit
// `AS b` and the bare `a b` form must produce the same aliased TargetEntry.
// Oracle: ANTLR accepts both; omni records Alias either way.
func TestSelectCoverage_ProjectionAlias(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "bare_alias",
			input: "SELECT a b FROM t",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a} Alias:b}] From:VarRef{Name:t}}",
		},
		{
			name:  "as_alias",
			input: "SELECT a AS b FROM t",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a} Alias:b}] From:VarRef{Name:t}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_Having covers havingClause (PartiQLParser.g4:294-295
// `HAVING arg=exprSelect`), both with a preceding GROUP BY and standalone
// (the grammar makes groupClause and havingClause independently optional in
// exprSelect, g4:461-462). Oracle: ANTLR accepts both; omni stores Having.
func TestSelectCoverage_Having(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "group_by_having",
			input: "SELECT a FROM t GROUP BY a HAVING a > 1",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} GroupBy:GroupByClause{Items:[GroupByItem{Expr:VarRef{Name:a}}]} Having:BinaryExpr{Op:> Left:VarRef{Name:a} Right:NumberLit{Val:1}}}",
		},
		{
			name:  "having_no_group",
			input: "SELECT a FROM t HAVING a > 1",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:VarRef{Name:t} Having:BinaryExpr{Op:> Left:VarRef{Name:a} Right:NumberLit{Val:1}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runAcceptGolden(t, tc.input, tc.want)
		})
	}
}

// TestSelectCoverage_HavingReject covers a HAVING with no argument expression.
// Oracle: ANTLR rejects (havingClause requires arg=exprSelect, g4:295).
func TestSelectCoverage_HavingReject(t *testing.T) {
	runReject(t, "SELECT a FROM t HAVING")
}

// runAcceptGolden parses input via ParseStatement (which asserts EOF) and
// requires it to succeed with ast.NodeToString(stmt) == want.
func runAcceptGolden(t *testing.T, input, want string) {
	t.Helper()
	p := NewParser(input)
	stmt, err := p.ParseStatement()
	if err != nil {
		t.Fatalf("expected accept, got error: %v\ninput: %s", err, input)
	}
	got := ast.NodeToString(stmt)
	if got != want {
		t.Errorf("AST mismatch\ninput: %s\ngot:  %s\nwant: %s", input, got, want)
	}
}

// runReject parses input via ParseStatement and requires it to error (the
// oracle decides accept/reject; the exact message is implementation-defined).
func runReject(t *testing.T, input string) {
	t.Helper()
	p := NewParser(input)
	stmt, err := p.ParseStatement()
	if err == nil {
		t.Fatalf("expected rejection, but parsed OK: %s\ninput: %s", ast.NodeToString(stmt), input)
	}
}
