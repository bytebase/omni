package ast

import (
	"strings"
	"testing"
)

// TestIndexColumnZeroPrefixSerializes keeps a written zero prefix, KEY (s(0)),
// apart from no prefix in NodeToString, as the MariaDB writer does: TiDB
// answers 1391 for the one and takes the other.
func TestIndexColumnZeroPrefixSerializes(t *testing.T) {
	zero := NodeToString(&IndexColumn{Expr: &ColumnRef{Column: "s"}, HasPrefix: true})
	none := NodeToString(&IndexColumn{Expr: &ColumnRef{Column: "s"}})
	if !strings.Contains(zero, ":len 0") || strings.Contains(none, ":len") {
		t.Fatalf("zero prefix %q and no prefix %q serialize alike", zero, none)
	}
}
