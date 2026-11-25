package parser

import (
	"strings"
	"testing"
)

// TestStrictParseRejectsTrailingTokens: a statement prefix followed by junk
// is a syntax error and contributes no node.
func TestStrictParseRejectsTrailingTokens(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 )))",
		"SELECT a FROM t 1 2",
		"SELECT a FROM t WHERE x = 1 ) ) ) DROP",
		"INSERT INTO t (a) VALUES (1) extra",
		// The embedded subquery is re-parsed by fillSubqueries; its failure
		// must drop the outer statement too, not just add an error.
		"SELECT (SELECT 1 FROM t a b) AS x FROM u",
		"SELECT a FROM t WHERE EXISTS (SELECT (SELECT 1 FROM w 1 2) FROM v)",
	} {
		file, errs := parseForTest(sql)
		if len(errs) == 0 {
			t.Errorf("parseForTest(%q) succeeded, want trailing-token error", sql)
		}
		if len(file.Stmts) != 0 {
			t.Errorf("parseForTest(%q) returned %d statements, want none: the truncated prefix must not leak", sql, len(file.Stmts))
		}
	}
}

// TestStrictParseSubqueryErrorPositions: a re-parse error inside an embedded
// subquery lands on the offending token in the outer text, whatever
// whitespace follows the opening delimiter.
func TestStrictParseSubqueryErrorPositions(t *testing.T) {
	for _, c := range []struct{ sql, at string }{
		{"SELECT EXISTS(  SELECT 1 FROM t a b)", "b)"},
		{"SELECT (\n  SELECT 1 FROM t a b) AS x FROM u", "b)"},
		{"SELECT ARRAY(   SELECT x FROM t a b)", "b)"},
	} {
		_, errs := parseForTest(c.sql)
		if len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want a subquery error", c.sql)
			continue
		}
		if want := strings.Index(c.sql, c.at); errs[0].Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%q)", c.sql, errs[0].Position, want, c.at)
		}
	}
}

// TestStrictParseRejectsEmptySubqueryBodies: an empty or comment-only
// EXISTS/ARRAY/scalar body holds no query; strict Parse errors and drops the
// statement instead of skipping the re-parse.
func TestStrictParseRejectsEmptySubqueryBodies(t *testing.T) {
	for _, sql := range []string{"SELECT EXISTS()", "SELECT ARRAY(/* c */)", "SELECT ( ) AS x"} {
		file, errs := parseForTest(sql)
		if len(errs) == 0 || len(file.Stmts) != 0 {
			t.Errorf("Parse(%q): errs=%d stmts=%d, want an error and no node", sql, len(errs), len(file.Stmts))
		}
	}
}

// TestStrictParseDropsNodeOnHintError: a malformed statement hint is recorded
// before parseStmt runs; the node it still builds must not survive strict
// mode.
func TestStrictParseDropsNodeOnHintError(t *testing.T) {
	file, errs := parseForTest("@[5@] SELECT 1")
	if len(errs) == 0 || len(file.Stmts) != 0 {
		t.Errorf("errs=%d stmts=%d, want an error and no node", len(errs), len(file.Stmts))
	}
}

// TestStrictParseErrorsInSourceOrder: lex errors are merged into the parse
// errors by position rather than appended after every segment.
func TestStrictParseErrorsInSourceOrder(t *testing.T) {
	_, errs := parseForTest("\x00; SELECT FROM")
	if len(errs) < 2 {
		t.Fatalf("errs = %v, want the invalid byte and the syntax error", errs)
	}
	for i := 1; i < len(errs); i++ {
		if errs[i].Position < errs[i-1].Position {
			t.Fatalf("errors out of source order: %v", errs)
		}
	}
}

// TestStrictParseDropsNodeOnLexError: a lex error is collected by the
// whole-input pass after the node was built (an unterminated comment after a
// complete statement); strict Parse returns the error and no node for that
// segment while the other segments survive, and best-effort keeps the prefix.
func TestStrictParseDropsNodeOnLexError(t *testing.T) {
	file, errs := parseForTest("SELECT 1 /* unterminated")
	if len(errs) == 0 || len(file.Stmts) != 0 {
		t.Errorf("errs=%d stmts=%d, want an error and no node", len(errs), len(file.Stmts))
	}
	file, errs = parseForTest("SELECT 1; SELECT 2 'unterminated")
	if len(errs) == 0 || len(file.Stmts) != 1 {
		t.Errorf("two segments: errs=%d stmts=%d, want an error and only the clean segment", len(errs), len(file.Stmts))
	}
	if r := ParseBestEffort("SELECT 1 /* unterminated"); len(r.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort stmts=%d, want the parsed prefix", len(r.File.Stmts))
	}
}

// TestParseBestEffortToleratesTrailingTokens: the tolerant entry keeps the
// parsed prefix for partial-input consumers; only strict Parse rejects.
func TestParseBestEffortToleratesTrailingTokens(t *testing.T) {
	result := ParseBestEffort("SELECT a FROM t )))")
	if len(result.Errors) != 0 {
		t.Errorf("ParseBestEffort errors = %v, want none", result.Errors)
	}
	if len(result.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort stmts = %d, want the parsed prefix", len(result.File.Stmts))
	}
}

// TestStrictParseMultiStatementIsolation: a bad segment errors without taking
// the good segments around it down.
func TestStrictParseMultiStatementIsolation(t *testing.T) {
	file, errs := parseForTest("SELECT 1; SELECT 2 ))); SELECT 3")
	if len(errs) == 0 {
		t.Fatal("middle segment junk not reported")
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("Stmts = %d, want 2 (segments 1 and 3)", len(file.Stmts))
	}
}

// TestStrictParseCorpusCanary is the categorical guarantee: every official
// corpus block that parses cleanly as one statement must FAIL once trailing
// junk is appended, and contribute no node. If any construct regresses into a
// swallow again, its corpus entry trips this canary. The junk starts on a new
// line so a trailing line comment in a block cannot absorb it.
func TestStrictParseCorpusCanary(t *testing.T) {
	checked := 0
	for _, b := range collectOfficialCorpus(t) {
		if _, skipped := officialCorpusSkips[b.Key()]; skipped {
			continue
		}
		clean, errs := parseForTest(b.Text)
		if len(errs) > 0 || len(clean.Stmts) != 1 {
			continue // not one cleanly parseable statement; nothing to prove
		}
		// Append the junk to the statement's own segment: a block may carry a
		// line comment after its ';', and junk placed after that comment would
		// be a separate (correctly rejected) segment, proving nothing.
		segs := Split(b.Text)
		if len(segs) != 1 {
			continue
		}
		canary := segs[0].Text + "\n )))"
		file, errs := parseForTest(canary)
		if len(errs) == 0 || len(file.Stmts) != 0 {
			t.Errorf("%s: trailing junk swallowed (errs=%d, stmts=%d):\n%s", b.Label(), len(errs), len(file.Stmts), indentBlock(b.Text))
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("corpus canary checked zero statements")
	}
	t.Logf("canary verified %d corpus statements", checked)
}
