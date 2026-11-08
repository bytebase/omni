package catalog

import (
	"strings"
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
		// `?` in a general SimpleExpr position (parser.y:7886) is accepted by
		// the parser but must still be refused here.
		"SELECT ? FROM t",
		"SELECT a + ? FROM t",
		"SELECT a FROM t WHERE a = ?",
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

// TestCreateViewWithParamMarkerRejected checks the DDL entry point: CREATE VIEW
// with a `?` in its body must fail Exec and must not install a view, since the
// analyzer error alone is swallowed by createView.
func TestCreateViewWithParamMarkerRejected(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT)")

	for _, sql := range []string{
		"CREATE VIEW v AS SELECT a FROM t LIMIT ?",
		"CREATE VIEW v AS SELECT a FROM t LIMIT ? OFFSET ?",
		"CREATE VIEW v AS SELECT (SELECT 1 FROM t LIMIT ?) AS c FROM t",
		// `?` as a general SimpleExpr operand (parser.y:7886).
		"CREATE VIEW v AS SELECT ? AS c FROM t",
		"CREATE VIEW v AS SELECT a + ? AS c FROM t",
	} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Errorf("Exec(%q) parse error: %v", sql, err)
			continue
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Errorf("Exec(%q) = %+v, want a rejection of the ? marker", sql, results)
			continue
		}
		if msg := results[0].Error.Error(); !strings.Contains(msg, "?") {
			t.Errorf("Exec(%q) error = %q, want it to mention the ? marker", sql, msg)
		}
		if _, ok := c.GetDatabase("testdb").Views["v"]; ok {
			t.Errorf("Exec(%q) installed view v despite the rejection", sql)
		}
	}
}

// TestDDLParamMarkerRejected checks that a `?` parameter marker in a table-DDL
// expression — a column default, generated column or CHECK — is refused by
// Exec and never stored. The parser accepts `?` as a general operand
// (parser.y:7886); the catalog boundary rejects any DDL that carries one.
func TestDDLParamMarkerRejected(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT)")

	for _, sql := range []string{
		// tidb rejects an expression DEFAULT at parse (it has no DEFAULT (expr)
		// form), so a marker there never reaches the catalog; AS and CHECK do.
		"CREATE TABLE u (a INT, b INT AS (a + ?))",
		"CREATE TABLE u (a INT, CHECK (a > ?))",
	} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			// The parser refused it outright; that is also a rejection and
			// stores nothing.
			continue
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Errorf("Exec(%q) = %+v, want a rejection of the ? marker", sql, results)
			continue
		}
		if msg := results[0].Error.Error(); !strings.Contains(msg, "?") {
			t.Errorf("Exec(%q) error = %q, want it to mention the ? marker", sql, msg)
		}
	}
}

// TestRoutineBodyParamMarkerRejected checks that a `?` marker inside the rows of
// an INSERT ... VALUES or a VALUES statement — which nodes.Inspect's walker
// skips — is still refused when it rides in a stored routine, trigger, or view
// body, so no unbound definition enters the catalog.
func TestRoutineBodyParamMarkerRejected(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, "CREATE TABLE t (a INT)")
	for _, sql := range []string{
		"CREATE PROCEDURE p() INSERT INTO t VALUES (?)",
		"CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW INSERT INTO t VALUES (?)",
		"CREATE VIEW v AS SELECT * FROM (VALUES ROW(?)) AS x",
	} {
		results, err := c.Exec(sql, nil)
		if err != nil {
			continue
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Errorf("Exec(%q) = %+v, want a rejection of the ? marker", sql, results)
			continue
		}
		if msg := results[0].Error.Error(); !strings.Contains(msg, "?") {
			t.Errorf("Exec(%q) error = %q, want it to mention the ? marker", sql, msg)
		}
	}
}
