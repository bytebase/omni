package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
)

// TestParamMarkerRejectedByAnalyzer checks that the view analyzer refuses a `?`
// parameter marker. Markers cannot appear in a view definition (the only thing
// this analyzer processes), and the normalizing deparse cannot round-trip their
// positional binding, so they are rejected rather than carried through the IR.
func TestParamMarkerRejectedByAnalyzer(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t (a INT)")

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
		if _, err := c.AnalyzeSelectStmt(list.Items[0].(*nodes.SelectStmt)); err == nil {
			t.Errorf("AnalyzeSelectStmt(%q) = nil error, want a rejection of the ? marker", sql)
		}
	}
}

// TestCreateViewWithParamMarkerRejected checks the DDL entry point: CREATE VIEW
// with a `?` in its body must fail Exec and must not install a view, since the
// analyzer error alone is swallowed by createView.
func TestCreateViewWithParamMarkerRejected(t *testing.T) {
	for _, sql := range []string{
		"CREATE VIEW v AS SELECT a FROM t LIMIT ?",
		"CREATE VIEW v AS SELECT a FROM t LIMIT ? OFFSET ?",
		"CREATE VIEW v AS SELECT (SELECT 1 FROM t LIMIT ?) AS c FROM t",
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE t (a INT)")

		results, _ := c.Exec(sql, &ExecOptions{ContinueOnError: true})
		if len(results) != 1 || results[0].Error == nil {
			t.Errorf("Exec(%q) = %+v, want a rejection of the ? marker", sql, results)
			continue
		}
		if msg := results[0].Error.Error(); !strings.Contains(msg, "?") {
			t.Errorf("Exec(%q) error = %q, want it to mention the ? marker", sql, msg)
		}
		if _, ok := c.GetDatabase("test").Views["v"]; ok {
			t.Errorf("Exec(%q) installed view v despite the rejection", sql)
		}
	}
}
