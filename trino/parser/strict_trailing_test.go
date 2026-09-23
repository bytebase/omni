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
		file, errs := Parse(sql)
		if len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want trailing-token error", sql)
		}
		if len(file.Stmts) != 0 {
			t.Errorf("Parse(%q) returned %d statements, want none: the truncated prefix must not leak", sql, len(file.Stmts))
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
	file, errs := Parse("SELECT 1; SELECT 2 ))); SELECT 3")
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
	_, errs := Parse("SELECT 1 ))) 'unterminated")
	if len(errs) < 2 {
		t.Fatalf("errs = %v, want the trailing-token error AND the lex error", errs)
	}
}

func TestStrictParseReportsLexErrorInDroppedSegment(t *testing.T) {
	// Split drops a segment that lexes to nothing, and with it the lex error
	// that made it empty; strict Parse promotes it from a whole-input pass.
	_, errs := Parse("/* unterminated")
	if len(errs) == 0 {
		t.Fatal("Parse(\"/* unterminated\") reported no error")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Msg, "comment") {
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
			clean, errs := Parse(sql)
			if len(errs) > 0 || len(clean.Stmts) != 1 {
				continue // not one cleanly parseable statement; nothing to prove
			}
			canary := sql + "\n )))"
			file, errs := Parse(canary)
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
