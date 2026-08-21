package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// TestStrictParseRejectsTrailingTokens covers the hand-picked trailing-junk
// shapes the old parser silently swallowed (BYT-10085).
func TestStrictParseRejectsTrailingTokens(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 )))",
		"SELECT 1 */ 2",
		"SELECT j->'$.a' FROM t",
		"SELECT a FROM t WHERE x = 1 ) ) ) DROP",
	} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want trailing-token error", sql)
		}
	}
}

// canaryExempt reports whether a statement's parser deliberately captures
// its tail as raw text, which absorbs the canary junk before the strict
// trailing-token check can see it: ShowStmt variant args
// (collectRemainingRaw), CREATE FUNCTION property/body capture, and the
// CTAS raw SELECT. These raw-capture parsers are a swallow mechanism of
// their own — tracked separately as the raw-capture cleanup — and are
// exempted here rather than papered over with a weaker assertion.
func canaryExempt(n ast.Node) bool {
	switch stmt := n.(type) {
	case *ast.ShowStmt, *ast.CreateFunctionStmt:
		return true
	case *ast.CreateTableStmt:
		return stmt.AsSelect != nil
	}
	return false
}

// TestStrictParseCorpusCanary is the categorical guarantee: every statement
// the legacy corpus parses cleanly must FAIL once trailing junk is appended.
// If any construct regresses into a swallow again, its corpus entry trips
// this canary. The junk starts on a new line so a trailing line comment in a
// corpus statement cannot absorb it.
func TestStrictParseCorpusCanary(t *testing.T) {
	files, err := filepath.Glob("testdata/legacy/*.sql")
	if err != nil || len(files) == 0 {
		t.Fatalf("no corpus files: %v", err)
	}
	regression, _ := filepath.Glob("testdata/legacy/regression/*.sql")
	files = append(files, regression...)

	checked := 0
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for _, seg := range Split(string(data)) {
			clean, errs := Parse(seg.Text)
			if len(errs) > 0 {
				continue // not cleanly parseable on its own; nothing to prove
			}
			if len(clean.Stmts) == 1 && canaryExempt(clean.Stmts[0]) {
				continue
			}
			canary := seg.Text + "\n )))"
			if _, errs := Parse(canary); len(errs) == 0 {
				stmt := strings.Join(strings.Fields(seg.Text), " ")
				if len(stmt) > 90 {
					stmt = stmt[:90] + "..."
				}
				t.Errorf("%s: trailing junk swallowed after: %s", filepath.Base(f), stmt)
			}
			checked++
		}
	}
	if checked == 0 {
		t.Fatal("corpus canary checked zero statements")
	}
	t.Logf("canary verified %d corpus statements", checked)
}

// TestStrictParseMultiStatementIsolation: a bad segment errors without
// taking the good segments around it down.
func TestStrictParseMultiStatementIsolation(t *testing.T) {
	file, errs := Parse("SELECT 1; SELECT 2 ))); SELECT 3")
	if len(errs) == 0 {
		t.Fatal("middle segment junk not reported")
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("Stmts = %d, want 2 (segments 1 and 3)", len(file.Stmts))
	}
}

// TestParseBestEffortToleratesTrailingTokens: the tolerance contract for
// partial-input consumers stays — only strict Parse rejects.
func TestParseBestEffortToleratesTrailingTokens(t *testing.T) {
	result := ParseBestEffort("SELECT 1 )))")
	if len(result.Errors) != 0 {
		t.Errorf("ParseBestEffort errors = %v, want none", result.Errors)
	}
	if len(result.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort stmts = %d, want the parsed prefix", len(result.File.Stmts))
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

func TestStrictParseRejectsMalformedExplain(t *testing.T) {
	// The RawQuery fallback is best-effort recovery only; on the strict path
	// it must not swallow a nested parse failure.
	for _, sql := range []string{"EXPLAIN SELECT * FROM", "EXPLAIN SELECT (", "EXPLAIN BOGUS x", "EXPLAIN", "EXPLAIN /*comment*/"} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want nested error", sql)
		}
		if r := ParseBestEffort(sql); len(r.Errors) != 0 {
			t.Errorf("ParseBestEffort(%q) errors = %v, want fallback recovery", sql, r.Errors)
		}
	}
	// A well-formed explained statement still parses strictly.
	if _, errs := Parse("EXPLAIN SELECT 1"); len(errs) != 0 {
		t.Errorf("EXPLAIN SELECT 1 errors: %v", errs)
	}
}

func TestSetTransactionCharacteristics(t *testing.T) {
	// SET TRANSACTION used to raw-capture to EOF, hiding any junk from the
	// strict check. The characteristics are now parsed for real.
	for _, sql := range []string{
		"SET TRANSACTION READ ONLY",
		"SET TRANSACTION READ WRITE",
		"SET TRANSACTION ISOLATION LEVEL READ COMMITTED",
		"SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED",
		"SET TRANSACTION ISOLATION LEVEL REPEATABLE READ",
		"SET TRANSACTION ISOLATION LEVEL SERIALIZABLE, READ ONLY",
	} {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
	for _, sql := range []string{
		"SET TRANSACTION",
		"SET TRANSACTION BOGUS",
		"SET TRANSACTION READ ONLY )))",
		"SET TRANSACTION ISOLATION LEVEL READ",
	} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
}

func TestStrictParseRejectsMalformedSetExpressions(t *testing.T) {
	// parseSetItem's raw fallback consumed to the comma or EOF on an
	// expression error, hiding the failure from the strict check.
	for _, sql := range []string{"SET x = (", "SET x = 1 +", "SET x = EXISTS ("} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
		if r := ParseBestEffort(sql); len(r.Errors) != 0 {
			t.Errorf("ParseBestEffort(%q) errors = %v, want fallback recovery", sql, r.Errors)
		}
	}
	if _, errs := Parse("SET x = 1, y = 'a'"); len(errs) != 0 {
		t.Errorf("well-formed SET errors: %v", errs)
	}
}

func TestStrictParseRejectsValuelessSetAssignments(t *testing.T) {
	// The bare-name forms reached EOF or the comma before the strict check
	// could see anything; the engine requires the separator and a value.
	for _, sql := range []string{"SET x", "SET x =", "SET x, y", "SET x 1"} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
	// Engine-valid forms keep parsing, including the scoped TRANSACTION
	// spelling that must not fall into the generic assignment path.
	for _, sql := range []string{
		"SET NAMES utf8",
		"SET PASSWORD = PASSWORD('123456')",
		"SET GLOBAL exec_mem_limit = 137438953472",
		"SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ",
		"SET GLOBAL TRANSACTION READ ONLY",
	} {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}

	// The scoped spelling keeps its qualifier on the AST — SET SESSION
	// TRANSACTION must stay distinguishable from SET TRANSACTION.
	file, _ := Parse("SET GLOBAL TRANSACTION READ ONLY")
	item := file.Stmts[0].(*ast.SetStmt).Items[0]
	if item.Scope != "GLOBAL" || item.Raw != "READ ONLY" {
		t.Errorf("scoped transaction item = %+v, want Scope GLOBAL Raw READ ONLY", item)
	}
	file, _ = Parse("SET TRANSACTION READ WRITE")
	if got := file.Stmts[0].(*ast.SetStmt).Items[0].Scope; got != "" {
		t.Errorf("unqualified transaction Scope = %q, want empty", got)
	}
}

func TestStrictParseRejectsMalformedShowClauses(t *testing.T) {
	// parseShowLikeWhere consumed to EOF before failing, so the leftover
	// check alone never saw SHOW TABLES WHERE ( go wrong.
	for _, sql := range []string{"SHOW TABLES WHERE (", "SHOW DATABASES WHERE 1 +", "SHOW TABLES FROM", "SHOW TABLES IN"} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
	if _, errs := Parse("SHOW TABLES WHERE table_name LIKE 'a%'"); len(errs) != 0 {
		t.Errorf("well-formed SHOW WHERE errors: %v", errs)
	}
}

func TestStrictParsePromotesFilteredSegmentLexErrors(t *testing.T) {
	// Split drops comment-only segments, and their lex errors with them; the
	// strict full-input sweep has to bring those diagnostics back.
	for _, sql := range []string{"/* unterminated", "SELECT 1; /* unterminated"} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want unterminated-comment error", sql)
		}
	}
	// The full-input sweep must not duplicate an error a segment already
	// reported: the segment pairs a syntax error at the invalid token with
	// the promoted lex error, and the sweep adds nothing on top.
	_, errs := Parse("SELECT 'unterminated")
	if len(errs) != 2 {
		t.Errorf("Parse(SELECT 'unterminated) errs = %v, want the segment's two", errs)
	}
}

func TestStrictParseRejectsIncompleteSetNamesCharset(t *testing.T) {
	// The specialized SET forms finished at EOF with their operands missing,
	// invisible to the leftover-token check. Engine-verified rejects on both
	// engines.
	for _, sql := range []string{"SET NAMES", "SET CHARSET", "SET NAMES utf8 COLLATE"} {
		if _, errs := Parse(sql); len(errs) == 0 {
			t.Errorf("Parse(%q) succeeded, want error", sql)
		}
	}
	for _, sql := range []string{"SET NAMES utf8", "SET NAMES utf8 COLLATE utf8_general_ci", "SET CHARSET utf8"} {
		if _, errs := Parse(sql); len(errs) != 0 {
			t.Errorf("Parse(%q) errors: %v", sql, errs)
		}
	}
}

func TestParseBestEffortKeepsSetTransactionRecovery(t *testing.T) {
	// Completion-style partial input must still yield a statement in
	// best-effort mode; only strict Parse rejects it.
	r := ParseBestEffort("SET TRANSACTION ISOLATION LEVEL READ")
	if len(r.Errors) != 0 || len(r.File.Stmts) != 1 {
		t.Fatalf("ParseBestEffort: stmts=%d errs=%v, want recovered statement", len(r.File.Stmts), r.Errors)
	}
	if _, errs := Parse("SET TRANSACTION ISOLATION LEVEL READ"); len(errs) == 0 {
		t.Error("strict Parse accepted the incomplete characteristic")
	}
}
