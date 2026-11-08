package deparse

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestIntervalValueBareDeparse pins how an INTERVAL value deparses. TiDB reads
// `INTERVAL (x) unit` as the INTERVAL function, never the literal, so a bit_expr
// value must come back out bare — with only the parentheses its own precedence
// requires, never the full wrapping writeBinaryExpr gives a binary expression
// elsewhere. Each want re-parses to the same interval literal (checked below).
func TestIntervalValueBareDeparse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"shift", "d + INTERVAL 1 << 37 SECOND", "(`d` + interval 1 << 37 second)"},
		{"shift_minus", "d - INTERVAL 1 << 39 SECOND", "(`d` - interval 1 << 39 second)"},
		{"add", "d + INTERVAL 1 + 2 DAY", "(`d` + interval 1 + 2 day)"},
		{"mul_add_precedence", "d + INTERVAL 1 * 2 + 3 DAY", "(`d` + interval 1 * 2 + 3 day)"},
		{"shift_chain", "d + INTERVAL 1 << 2 << 3 SECOND", "(`d` + interval 1 << 2 << 3 second)"},
		{"bit_or_shift", "d + INTERVAL 1 | 2 << 3 DAY", "(`d` + interval 1 | 2 << 3 day)"},
		{"inner_parens_kept", "d + INTERVAL 1 * (2 + 3) DAY", "(`d` + interval 1 * (2 + 3) day)"},
		// A grouped predicate keeps its parentheses, which it needs to reparse:
		// IN binds looser than +, and writeExpr writes it without its own pair.
		{"grouped_predicate", "d + INTERVAL 1 + (a IN (SELECT 1)) DAY", "(`d` + interval 1 + (`a` in (select 1)) day)"},
		{"column_value", "d + INTERVAL a << 2 SECOND", "(`d` + interval `a` << 2 second)"},
		{"signed_literal", "d + INTERVAL -1 DAY", "(`d` + interval -1 day)"},
		{"plain_literal", "d + INTERVAL 1 DAY", "(`d` + interval 1 day)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := parseExpr(t, tc.input)
			got := Deparse(node)
			if got != tc.want {
				t.Errorf("Deparse(%q) = %q, want %q", tc.input, got, tc.want)
			}
			// The value must not come back parenthesized at the top, or a
			// re-parse reads it as the INTERVAL function.
			seg := got[strings.Index(got, "interval ")+len("interval "):]
			if strings.HasPrefix(seg, "(") {
				t.Errorf("value starts with '(' (reads as the INTERVAL function): %q", got)
			}
			// It re-parses, to the same interval literal.
			list, err := parser.Parse("SELECT " + got + " FROM t")
			if err != nil {
				t.Fatalf("re-parse %q: %v", got, err)
			}
			if again := Deparse(list.Items[0].(*ast.SelectStmt).TargetList[0]); again != got {
				t.Errorf("not stable: %q then %q", got, again)
			}
		})
	}
}
