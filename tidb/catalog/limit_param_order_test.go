package catalog

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestParamMarkerRejectedByAnalyzer checks that the view analyzer refuses a `?`
// parameter marker. Markers cannot appear in a view definition (the only thing
// this analyzer processes), and the normalizing deparse cannot round-trip their
// positional binding, so they are rejected rather than carried through the IR.
func TestParamMarkerRejectedByAnalyzer(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT)")

	for _, sql := range []string{
		"SELECT a FROM t LIMIT ?",
		"SELECT a FROM t LIMIT ?, ?",
		"SELECT a FROM t LIMIT ? OFFSET ?",
		"SELECT (SELECT 1 FROM t LIMIT ?) FROM t",
	} {
		list, err := parser.Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", sql, err)
			continue
		}
		if _, err := c.AnalyzeSelectStmt(list.Items[0].(*ast.SelectStmt)); err == nil {
			t.Errorf("AnalyzeSelectStmt(%q) = nil error, want a rejection of the ? marker", sql)
		}
	}
}
