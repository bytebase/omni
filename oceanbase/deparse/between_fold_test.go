package deparse

import (
	"strings"
	"testing"
)

// TestFoldBetweenIsBounded: a BETWEEN over constant bounds is two comparisons, which
// writes its operand twice, and BETWEENs nested in the operand would double the text at
// every level. Past maxFoldedOperand the BETWEEN is kept, so the text a statement
// deparses to stays within a small multiple of the text it was written in.
func TestFoldBetweenIsBounded(t *testing.T) {
	folded := func(sql string) string {
		sel := parseSelect(t, sql)
		RewriteSelectStmt(sel)
		return DeparseSelect(sel)
	}
	if got, want := folded("SELECT a + 1 BETWEEN 1 AND 10 AS x FROM t"), "select (((`a` + 1) >= 1) and ((`a` + 1) <= 10)) AS `x` from `t`"; got != want {
		t.Errorf("a small operand is folded:\n got: %s\nwant: %s", got, want)
	}
	for _, depth := range []int{4, 20, 200} {
		nested := "a"
		for i := 0; i < depth; i++ {
			nested = "(" + nested + " BETWEEN 0 AND 1)"
		}
		sql := "SELECT " + nested + " AS x FROM t"
		got := folded(sql)
		if len(got) > 64*len(sql) {
			t.Errorf("depth %d: %d bytes of SQL deparse to %d", depth, len(sql), len(got))
		}
		if depth > 4 && !strings.Contains(got, " between ") {
			t.Errorf("depth %d: no BETWEEN is kept: %.200s", depth, got)
		}
	}
}
