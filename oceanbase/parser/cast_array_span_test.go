package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestNestedCastArraySpans: a nested ARRAY that closes with the first half of a >> ends after
// that half, so each element's span is its own text and not its parent's bracket too.
func TestNestedCastArraySpans(t *testing.T) {
	for sql, want := range map[string][]string{
		"SELECT CAST(a AS ARRAY<ARRAY<INT>>)":               {"ARRAY<ARRAY<INT>>", "ARRAY<INT>"},
		"SELECT CAST(a AS ARRAY<ARRAY<ARRAY<INT>>>)":        {"ARRAY<ARRAY<ARRAY<INT>>>", "ARRAY<ARRAY<INT>>", "ARRAY<INT>"},
		"SELECT CAST(a AS ARRAY<ARRAY<ARRAY<ARRAY<INT>>>>)": {"ARRAY<ARRAY<ARRAY<ARRAY<INT>>>>", "ARRAY<ARRAY<ARRAY<INT>>>", "ARRAY<ARRAY<INT>>", "ARRAY<INT>"},
		"SELECT CAST(a AS ARRAY(ARRAY(INT)))":               {"ARRAY(ARRAY(INT))", "ARRAY(INT)"},
	} {
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		var cast *nodes.CastExpr
		nodes.Inspect(list, func(n nodes.Node) bool {
			if c, ok := n.(*nodes.CastExpr); ok && cast == nil {
				cast = c
			}
			return true
		})
		if cast == nil {
			t.Fatalf("%q: no CAST", sql)
		}
		var got []string
		for dt := cast.TypeName; dt != nil && dt.Name == "ARRAY"; dt = dt.ArrayElem {
			got = append(got, sql[dt.Loc.Start:dt.Loc.End])
		}
		if len(got) != len(want) {
			t.Fatalf("%q: spans %q, want %q", sql, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%q: span %d is %q, want %q", sql, i, got[i], want[i])
			}
		}
	}
}
