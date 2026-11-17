package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// legacyLexerFile is the legacy SnowflakeLexer.g4 grammar, the source of
// truth for keyword drift detection. It is mirrored under
// testdata/legacy_grammar (see the README there for provenance).
const legacyLexerFile = "SnowflakeLexer.g4"

// legacyReservedScriptFile is the legacy build_id_contains_non_reserved_keywords.py
// script, whose snowflake_reserved_keyword dict defines the reserved keyword
// set. Also mirrored under testdata/legacy_grammar.
const legacyReservedScriptFile = "build_id_contains_non_reserved_keywords.py"

// legacyGrammarPath returns the absolute path of one of the mirrored legacy
// files. The files are part of the repo, so a missing one is a real error,
// never a skip.
func legacyGrammarPath(t *testing.T, name string) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed; cannot locate legacy grammar")
	}
	path := filepath.Join(filepath.Dir(thisFile), "testdata", "legacy_grammar", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("legacy grammar not found at %s: %v", path, err)
	}
	return path
}

// TestKeywordCompleteness asserts that every single-word keyword token
// defined in SnowflakeLexer.g4 has a corresponding entry in keywordMap.
func TestKeywordCompleteness(t *testing.T) {
	lexerPath := legacyGrammarPath(t, legacyLexerFile)
	data, err := os.ReadFile(lexerPath)
	if err != nil {
		t.Fatalf("reading legacy lexer grammar %s: %v", lexerPath, err)
	}

	// Match lines like:  ABORT : 'ABORT';
	// Captures the literal between the single quotes.
	re := regexp.MustCompile(`(?m)^[A-Z_][A-Z0-9_]*\s*:\s*'([A-Z][A-Z0-9_]*)'\s*;`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatalf("no keyword tokens extracted from %s — regex broken?", lexerPath)
	}

	missing := []string{}
	for _, m := range matches {
		literal := m[1] // the captured keyword text
		lower := strings.ToLower(literal)
		if _, ok := keywordMap[lower]; !ok {
			missing = append(missing, literal)
		}
	}

	if len(missing) > 0 {
		t.Errorf("%d keywords missing from keywordMap (sample of first 10): %v",
			len(missing), missing[:min(len(missing), 10)])
	}

	t.Logf("checked %d keywords from legacy grammar; %d missing", len(matches), len(missing))
}

// TestReservedKeywordsCompleteness asserts that every keyword in the legacy
// build_id_contains_non_reserved_keywords.py snowflake_reserved_keyword dict
// is present in keywordReserved, EXCEPT for a known set of phantoms that
// the legacy Python script reserves but the legacy SnowflakeLexer.g4 does
// not actually emit as token rules. The omni lexer cannot mark a keyword
// as reserved if it has no kw* token, so these phantoms are deliberately
// excluded from keywordReserved (see the comment in keywords.go).
func TestReservedKeywordsCompleteness(t *testing.T) {
	scriptPath := legacyGrammarPath(t, legacyReservedScriptFile)
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("reading legacy reserved-keyword script %s: %v", scriptPath, err)
	}

	// knownPhantoms are keyword names that the Python reserved-set file
	// declares but the legacy SnowflakeLexer.g4 does not emit as a token
	// rule (two are commented out, seven are entirely absent). Verified
	// 2026-04-07. If a future grammar update adds any of these, drop the
	// entry from this set and add it to keywordReserved in keywords.go.
	knownPhantoms := map[string]bool{
		"CURRENT_USER":   true, // commented out in lexer grammar
		"FOLLOWING":      true, // not in lexer grammar
		"GSCLUSTER":      true, // not in lexer grammar
		"ISSUE":          true, // not in lexer grammar
		"LOCALTIME":      true, // not in lexer grammar
		"LOCALTIMESTAMP": true, // not in lexer grammar
		"REGEXP":         true, // not in lexer grammar
		"TRIGGER":        true, // commented out in lexer grammar
		"WHENEVER":       true, // not in lexer grammar
	}

	// Match lines like:    "ACCOUNT": True,
	// inside the snowflake_reserved_keyword dict.
	re := regexp.MustCompile(`(?m)^\s*"([A-Z_][A-Z0-9_]*)"\s*:\s*True\s*,?`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatalf("no reserved keywords extracted from %s — regex broken?", scriptPath)
	}

	missing := []string{}
	skipped := 0
	for _, m := range matches {
		kw := m[1]
		if knownPhantoms[kw] {
			skipped++
			continue
		}
		if !IsReservedKeyword(kw) {
			missing = append(missing, kw)
		}
	}

	if len(missing) > 0 {
		t.Errorf("%d reserved keywords missing from keywordReserved (not in known-phantom set): %v",
			len(missing), missing)
	}

	t.Logf("checked %d reserved keywords from legacy script; %d skipped as known phantoms; %d missing",
		len(matches), skipped, len(missing))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
