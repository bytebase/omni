package parser

import "testing"

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
