package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
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
		// `?` in a general SimpleExpr position is accepted by the parser but
		// must still be refused here.
		"SELECT ? FROM t",
		"SELECT a + ? FROM t",
		"SELECT a FROM t WHERE a = ?",
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
		// `?` as a general SimpleExpr operand.
		"CREATE VIEW v AS SELECT ? AS c FROM t",
		"CREATE VIEW v AS SELECT a + ? AS c FROM t",
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

// TestAlterViewWithParamMarkerRejected checks that ALTER VIEW with a `?` in its
// body fails Exec and leaves the existing view's definition unchanged; the
// ALTER VIEW path does not route through createView.
func TestAlterViewWithParamMarkerRejected(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE test")
	c.SetCurrentDatabase("test")
	mustExec(t, c, "CREATE TABLE t (a INT)")
	mustExec(t, c, "CREATE VIEW v AS SELECT a FROM t")
	before := c.GetDatabase("test").Views["v"].Definition

	results, _ := c.Exec("ALTER VIEW v AS SELECT a FROM t LIMIT ?", &ExecOptions{ContinueOnError: true})
	if len(results) != 1 || results[0].Error == nil {
		t.Fatalf("Exec(ALTER VIEW ... LIMIT ?) = %+v, want a rejection", results)
	}
	if msg := results[0].Error.Error(); !strings.Contains(msg, "?") {
		t.Errorf("error = %q, want it to mention the ? marker", msg)
	}
	if after := c.GetDatabase("test").Views["v"].Definition; after != before {
		t.Errorf("ALTER VIEW mutated the view despite the rejection: %q -> %q", before, after)
	}
}

// TestDDLParamMarkerRejected checks that a `?` parameter marker in a table-DDL
// expression — a column default, generated column or CHECK — is refused by
// Exec and never stored. The parser accepts `?` as a general operand; the
// catalog boundary rejects any DDL that carries one.
func TestDDLParamMarkerRejected(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE u (a INT DEFAULT (?))",
		"CREATE TABLE u (a INT, b INT AS (a + ?))",
		"CREATE TABLE u (a INT, CHECK (a > ?))",
		"ALTER TABLE t ADD COLUMN c INT DEFAULT (?)",
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE t (a INT)")
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
		if _, ok := c.GetDatabase("test").Tables["u"]; ok {
			t.Errorf("Exec(%q) created table u despite the rejection", sql)
		}
	}
}

// TestRoutineBodyParamMarkerRejected checks that a `?` marker inside the rows of
// an INSERT ... VALUES or a VALUES statement — which nodes.Inspect's walker
// skips — is still refused when it rides in a stored routine, trigger, or view
// body, so no unbound definition enters the catalog.
func TestRoutineBodyParamMarkerRejected(t *testing.T) {
	for _, sql := range []string{
		"CREATE PROCEDURE p() INSERT INTO t VALUES (?)",
		"CREATE TRIGGER trg BEFORE INSERT ON t FOR EACH ROW INSERT INTO t VALUES (?)",
		"CREATE VIEW v AS SELECT * FROM (VALUES ROW(?)) AS x",
	} {
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE t (a INT)")
		results, err := c.Exec(sql, nil)
		if err != nil {
			// A parser refusal is also a rejection and stores nothing.
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
