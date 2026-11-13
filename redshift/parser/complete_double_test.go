package parser

import (
	"slices"
	"testing"
)

// TestCollectDoublePrecision: double alone is a GenericType, but with the
// cursor after it completion still offers PRECISION.
func TestCollectDoublePrecision(t *testing.T) {
	for _, q := range []string{"CREATE TABLE t (c double ", "SELECT CAST(1 AS double "} {
		if cs := Collect(q, len(q)); !slices.Contains(cs.Tokens, PRECISION) {
			t.Errorf("Collect(%q) does not offer PRECISION", q)
		}
	}
}
