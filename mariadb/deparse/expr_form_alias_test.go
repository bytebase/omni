package deparse

import (
	"strings"
	"testing"
)

// TestExprFormAliasWithoutSource: with no source text, an unaliased item's
// name is rebuilt from the AST, so it has to carry what the text did. A
// one-argument directional TRIM is named TRIM(LEADING FROM ...), not
// TRIM_LEADING(...), and CHAR keeps its USING clause, so two calls that
// differ only in charset get two names.
func TestExprFormAliasWithoutSource(t *testing.T) {
	for _, tc := range []struct{ expr, want string }{
		{"TRIM(LEADING FROM c)", "TRIM(LEADING FROM c)"},
		{"TRIM(BOTH FROM ' x')", "TRIM(BOTH FROM ' x')"},
		{"CHAR(65 USING utf8mb4)", "CHAR(65 USING utf8mb4)"},
		{"CHAR(65, 66 USING latin1)", "CHAR(65, 66 USING latin1)"},
		{"CHAR(65)", "CHAR(65)"},
	} {
		var b strings.Builder
		writeExprAlias(&b, parseExpr(t, tc.expr))
		if got := b.String(); got != tc.want {
			t.Errorf("name of %s = %q, want %q", tc.expr, got, tc.want)
		}
	}
}
