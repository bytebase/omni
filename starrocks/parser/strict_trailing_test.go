package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestStrictParseRejectsTrailingTokens covers the hand-picked trailing-junk
// shapes the old parser silently swallowed (BYT-10085).
func TestStrictParseRejectsTrailingTokens(t *testing.T) {
	for _, sql := range []string{
		"SELECT 1 )))",
		"SELECT 1 */ 2",
		"SELECT j->'$.a' FROM t",
		"SELECT a FROM t WHERE x = 1 ) ) ) DROP",
		// Engine-rejected StarRocks forms that used to ride the swallow.
		"SELECT * FROM t TABLESAMPLE(1000 ROWS)",
		"SELECT a, SUM(b) FROM t GROUP BY a WITH ROLLUP",
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
	for _, sql := range []string{"EXPLAIN SELECT * FROM", "EXPLAIN SELECT (", "EXPLAIN BOGUS x"} {
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
