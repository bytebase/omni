package parser

import (
	"strings"
	"testing"
)

// TestStrictParseRejectsTrailingTokens: a statement prefix followed by junk
// is a syntax error, and the truncated node is not returned. Before this,
// Parse returned the prefix next to the error, and a consumer reading the
// File alone analyzed less than Trino would execute.
func TestStrictParseRejectsTrailingTokens(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 )))",
		"SELECT a FROM t 1 2",
		"SELECT a FROM t WHERE x = 1 ) ) ) DROP",
		"INSERT INTO t VALUES (1) extra",
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

// TestStrictParseValidatesSubqueryBodies: the placeholder scan captures an
// expression subquery's body as raw text; strict Parse must still reject a
// body Trino rejects, in outer-statement coordinates, and keep accepting the
// valid shapes.
func TestStrictParseValidatesSubqueryBodies(t *testing.T) {
	for _, c := range []struct{ sql, at string }{
		{"SELECT (SELECT 1 FROM t a b)", "b)"},
		{"SELECT (SELECT 1 FROM t WHERE x = ) FROM u", " ) FROM"},
		{"SELECT EXISTS () FROM t", ")"},
		{"SELECT EXISTS (/* c */) FROM t", ")"},
		{"SELECT EXISTS (DELETE FROM t) FROM t", "DELETE"},
		{"SELECT (SELECT 1 FROM t; SELECT 2) FROM t", ";"},
		{"SELECT x FROM t WHERE y IN (SELECT (SELECT 1 FROM v 1 2) FROM u)", "1 2"},
		// SHOW STATS FOR (query) captures its body the same way.
		{"SHOW STATS FOR (SELECT 1 FROM t a b)", "b)"},
		{"SHOW STATS FOR (SELECT 1 FROM t; SELECT 2)", ";"},
	} {
		file, err := parseForTest(c.sql)
		if len(err) == 0 {
			t.Errorf("Parse(%q) succeeded, want a subquery error", c.sql)
			continue
		}
		if len(file.Stmts) != 0 {
			t.Errorf("Parse(%q) kept %d statements, want none", c.sql, len(file.Stmts))
		}
		if want := strings.Index(c.sql, c.at); err[0].Position != want {
			t.Errorf("Parse(%q) error at %d, want %d (%q)", c.sql, err[0].Position, want, c.at)
		}
	}
	for _, sql := range []string{
		"SELECT (SELECT 1 FROM t) FROM u",
		"SELECT x FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.id = t.id)",
		"SELECT x FROM t WHERE y IN (SELECT y FROM u) AND z > ANY (SELECT z FROM v)",
		"SELECT (SELECT (SELECT 1 FROM w) FROM v) FROM u",
		"SELECT (WITH c AS (SELECT 1 AS a) SELECT a FROM c) FROM u",
		"SHOW STATS FOR (SELECT * FROM t WHERE x > 1)",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
	// Best-effort stays tolerant of an unfinished body.
	if r := ParseBestEffort("SELECT (SELECT 1 FROM t WHERE x = ) FROM u"); len(r.Errors) != 0 || len(r.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort: errors=%v stmts=%d, want tolerance", r.Errors, len(r.File.Stmts))
	}
}

// TestStrictParseReportsNestedLexErrorOnce: a lexical error inside a raw
// query body is seen by the nested strict parse and by the outer lexer;
// the diagnostic appears once.
func TestStrictParseReportsNestedLexErrorOnce(t *testing.T) {
	for _, sql := range []string{"SELECT (SELECT \x00)", "SELECT (SELECT 'x) FROM t", "SELECT (SELECT (SELECT 'x)) FROM t"} {
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

func TestStrictParseDrainsLexErrorsAfterTrailingJunk(t *testing.T) {
	// The lexer is lazy: without draining the segment after the first
	// trailing-token error, the unterminated string would never be lexed and
	// its diagnostic would be lost.
	_, errs := parseForTest("SELECT 1 ))) 'unterminated")
	if len(errs) < 2 {
		t.Fatalf("errs = %v, want the trailing-token error AND the lex error", errs)
	}
}

// TestStrictParseDropsNodeOnLexError: a lex error surfaces after the node was
// built (an unterminated comment after a complete statement); strict Parse
// returns the error and no node, best-effort keeps the prefix.
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

func TestStrictParseReportsLexErrorInDroppedSegment(t *testing.T) {
	// Split drops a segment that lexes to nothing, and with it the lex error
	// that made it empty; strict Parse promotes it from a whole-input pass.
	_, errs := parseForTest("/* unterminated")
	if len(errs) == 0 {
		t.Fatal("parseForTest(\"/* unterminated\") reported no error")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Message, "comment") {
			found = true
		}
	}
	if !found {
		t.Errorf("errs = %v, want one about the unterminated comment", errs)
	}
}

// TestStrictParseCorpusCanary is the categorical guarantee: every statement
// the accept corpora parse cleanly must FAIL once trailing junk is appended,
// and contribute no node. If any construct regresses into a swallow again,
// its corpus entry trips this canary. The junk starts on a new line so a
// trailing line comment in a corpus statement cannot absorb it.
func TestStrictParseCorpusCanary(t *testing.T) {
	corpora := [][]string{
		selectAcceptCorpus, relationAcceptCorpus, cteAcceptCorpus,
		setopAcceptCorpus, windowAcceptCorpus, exprAcceptCorpus,
		typeAcceptCorpus, oracleCorpus, dmlOracleCorpus, ddlOracleCorpus,
		dclTclOracleCorpus, routinesOracleCorpus, matchRecognizeOracleCorpus,
	}
	checked := 0
	for _, corpus := range corpora {
		for _, sql := range corpus {
			clean, errs := parseForTest(sql)
			if len(errs) > 0 || len(clean.Stmts) != 1 {
				continue // not one cleanly parseable statement; nothing to prove
			}
			canary := sql + "\n )))"
			file, errs := parseForTest(canary)
			if len(errs) == 0 || len(file.Stmts) != 0 {
				stmt := strings.Join(strings.Fields(sql), " ")
				if len(stmt) > 90 {
					stmt = stmt[:90] + "..."
				}
				t.Errorf("trailing junk swallowed after: %s (errs=%d, stmts=%d)", stmt, len(errs), len(file.Stmts))
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("corpus canary checked zero statements")
	}
	t.Logf("canary verified %d corpus statements", checked)
}
