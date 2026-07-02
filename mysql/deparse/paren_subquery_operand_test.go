package deparse

import (
	"testing"

	ast "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
)

// A parenthesized subquery as the LEFT operand of a binary operator —
// `((SELECT ...) = 0)`, the stock sys.metrics shape — must deparse to a
// parse→deparse FIXED POINT. The continuation parse consumes every paren
// layer through parseSelectStmt (wrapper nodes) while SubqueryExpr rendering
// adds one pair of its own; without folding the bare wrappers each cycle
// would grow the nesting by a level, so a catalog dump would never converge.
func TestDeparse_ParenSubqueryOperandFixedPoint(t *testing.T) {
	for _, sql := range []string{
		"SELECT ((SELECT count(0) FROM t) = 0)",
		"SELECT (((SELECT count(0) FROM t)) = 0)",
		"SELECT ((((SELECT count(0) FROM t)) = 0))",
		"SELECT if((((SELECT count(0) FROM t)) = 0),'NO','YES')",
		"SELECT ((SELECT max(a) FROM t) + 1)",
		"SELECT ((SELECT 1), 2) = ROW(1,2)",
		"SELECT 1 IN ((SELECT 1) = 1)",
		"SELECT 1 IN ((SELECT 1), 2)",
		"SELECT 1 IN (((SELECT 1) = 1))",
		// Wrapper with outer clauses keeps its scoping parens.
		"SELECT (((SELECT 1) LIMIT 1) = 1)",
		"SELECT (((SELECT 1) UNION (SELECT 1)) = 1)",
		// Pure forms must keep their pre-fix depth-preserving rendering.
		"SELECT ((SELECT 1))",
		"SELECT ((SELECT 1) UNION (SELECT 2))",
		"SELECT 1 IN ((SELECT 1))",
	} {
		t.Run(sql, func(t *testing.T) {
			d1 := DeparseSelect(parseSelect(t, sql))
			list, err := parser.Parse(d1)
			if err != nil {
				t.Fatalf("deparsed SQL does not reparse: %q: %v", d1, err)
			}
			sel, ok := list.Items[0].(*ast.SelectStmt)
			if !ok {
				t.Fatalf("reparse of %q: expected SelectStmt, got %T", d1, list.Items[0])
			}
			if d2 := DeparseSelect(sel); d2 != d1 {
				t.Errorf("not a fixed point:\n  cycle 1: %q\n  cycle 2: %q", d1, d2)
			}
		})
	}
}

// Subquery parens render as exactly ONE canonical pair — the engine collapses
// any extra depth in stored view bodies (8.0.32 + 5.7.25), so canonical text
// must too — and the derived alias text stays valid inside its backtick
// quoting.
func TestDeparse_ParenSubqueryOperandForms(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single_paren_subquery_compare",
			input:    "SELECT ((SELECT count(0) FROM t) = 0)",
			expected: "select ((select count(0) from `t`) = 0) AS `((SELECT COUNT(0) FROM t) = 0)`",
		},
		{
			// Redundant wrapper parens collapse to the canonical single pair,
			// matching the engine's stored form (8.0.32 + 5.7.25 both store
			// `((SELECT ...)) = 0` as `((select ...) = 0)`).
			name:     "double_paren_subquery_compare",
			input:    "SELECT (((SELECT count(0) FROM t)) = 0)",
			expected: "select ((select count(0) from `t`) = 0) AS `((SELECT COUNT(0) FROM t) = 0)`",
		},
		{
			name:     "row_with_subquery_head",
			input:    "SELECT ((SELECT 1), 2) = ROW(1,2)",
			expected: "select (row((select 1),2) = row(1,2)) AS `((SELECT 1), 2) = (1, 2)`",
		},
		{
			name:     "in_list_subquery_compare",
			input:    "SELECT 1 IN ((SELECT 1) = 1)",
			expected: "select (1 in (((select 1) = 1))) AS `1 IN ((SELECT 1) = 1)`",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DeparseSelect(parseSelect(t, tc.input)); got != tc.expected {
				t.Errorf("DeparseSelect(%q) =\n  %q\nwant:\n  %q", tc.input, got, tc.expected)
			}
		})
	}
}
