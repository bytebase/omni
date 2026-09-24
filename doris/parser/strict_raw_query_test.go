package parser

import (
	"strings"
	"testing"
)

// TestStrictParseValidatesRawQueryBodies: subquery placeholders and the CTAS
// query are captured as raw text; strict Parse must reject a body the engine
// rejects, in outer-statement coordinates, and keep accepting valid shapes.
func TestStrictParseValidatesRawQueryBodies(t *testing.T) {
	for _, c := range []struct{ sql, at string }{
		{"SELECT (SELECT 1 */ 2 FROM secret) FROM public", "*/"},
		{"SELECT (SELECT 1 FROM t WHERE x = ) FROM u", " ) FROM"},
		{"SELECT EXISTS () FROM public", ")"},
		{"SELECT EXISTS (/* c */) FROM public", ")"},
		{"SELECT EXISTS (DELETE FROM secret) FROM public", "DELETE"},
		{"SELECT a FROM t WHERE b IN (SELECT (SELECT 1 FROM v 1 2) FROM u)", "1 2"},
		{"CREATE TABLE dest AS SELECT 1 */ 2 FROM secret", "*/"},
	} {
		file, errs := parseForTest(c.sql)
		if len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want a subquery error", c.sql)
			continue
		}
		if len(file.Stmts) != 0 {
			t.Errorf("Parse(%q) kept %d statements, want none", c.sql, len(file.Stmts))
		}
		if want := strings.Index(c.sql, c.at); errs[0].Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%q)", c.sql, errs[0].Position, want, c.at)
		}
	}
	for _, sql := range []string{
		"SELECT (SELECT 1 FROM t) FROM u",
		"SELECT x FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id = t.id)",
		"SELECT x FROM t WHERE y IN (SELECT y FROM u)",
		"SELECT (SELECT (SELECT 1 FROM w) FROM v) FROM u",
		"SELECT a FROM (SELECT a FROM t) s",
		"CREATE TABLE dest AS SELECT a FROM src",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
	if r := ParseBestEffort("SELECT (SELECT 1 FROM t WHERE x = ) FROM u"); len(r.Errors) != 0 || len(r.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort: errors=%v stmts=%d, want tolerance", r.Errors, len(r.File.Stmts))
	}
}

// TestParseEmptyCTASBodyDoesNotPanic: `CREATE TABLE t AS ` once sliced the
// raw query backwards (EOF starts after the trailing space, prev ends at AS)
// and panicked; it is a positioned parse error in both modes.
func TestParseEmptyCTASBodyDoesNotPanic(t *testing.T) {
	for _, sql := range []string{"CREATE TABLE t AS ", "CREATE TABLE t AS", "CREATE TABLE t AS;", "CREATE TABLE t AS\n"} {
		_, errs := parseForTest(sql)
		if len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want an error", sql)
		}
		if r := ParseBestEffort(sql); len(r.Errors) == 0 {
			t.Errorf("ParseBestEffort(%q) reported no error", sql)
		}
	}
}

// TestStrictParseReportsNestedLexErrorOnce: a lexical error inside a raw
// query body is seen by the nested strict parse and by the outer lexer;
// the diagnostic appears once.
func TestStrictParseReportsNestedLexErrorOnce(t *testing.T) {
	for _, sql := range []string{"SELECT (SELECT 'x) FROM t", "SELECT (SELECT (SELECT 'x)) FROM t", "CREATE TABLE d AS SELECT 'x FROM t"} {
		_, errs := parseForTest(sql)
		if len(errs) == 0 {
			t.Errorf("Parse(%q) reported no error", sql)
			continue
		}
		type key struct {
			pos int
			msg string
		}
		dup := map[key]int{}
		for _, e := range errs {
			dup[key{e.Position, e.Message}]++
		}
		for k, n := range dup {
			if n > 1 {
				t.Errorf("Parse(%q) reported %q at %d %d times", sql, k.msg, k.pos, n)
			}
		}
	}
}

// TestStrictParseDropsNodeOnLexError: a lex error surfaces after the node was
// built; strict Parse returns the error and no node, best-effort keeps it.
func TestStrictParseDropsNodeOnLexError(t *testing.T) {
	for _, sql := range []string{"SELECT 1 /* unterminated", "SELECT 1 'unterminated"} {
		file, errs := parseForTest(sql)
		if len(errs) == 0 || len(file.Stmts) != 0 {
			t.Errorf("Parse(%q): errs=%d stmts=%d, want an error and no node", sql, len(errs), len(file.Stmts))
		}
		if r := ParseBestEffort(sql); len(r.File.Stmts) != 1 {
			t.Errorf("ParseBestEffort(%q) stmts=%d, want the parsed prefix", sql, len(r.File.Stmts))
		}
	}
}
