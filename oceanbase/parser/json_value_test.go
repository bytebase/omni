package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestQuotedIntroducerIsAName: a backtick-quoted charset name is an identifier, so `_latin1`'x'
// is the column _latin1 with the alias x, not a string of that charset.
func TestQuotedIntroducerIsAName(t *testing.T) {
	list, err := Parse("SELECT `_latin1`'x'")
	if err != nil {
		t.Fatal(err)
	}
	target := list.Items[0].(*nodes.SelectStmt).TargetList[0]
	rt, ok := target.(*nodes.ResTarget)
	if !ok {
		t.Fatalf("target is %T, want a column with an alias", target)
	}
	if _, ok := rt.Val.(*nodes.ColumnRef); !ok || rt.Name != "x" {
		t.Errorf("read %T named %q, want the column _latin1 named x", rt.Val, rt.Name)
	}
}
