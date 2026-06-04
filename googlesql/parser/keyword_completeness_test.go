package parser

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// legacyLexerPath is the absolute path to the legacy GoogleSQLLexer.g4
// grammar — the source of truth for keyword drift detection. The legacy
// grammar is a hand-port of ZetaSQL's flex_tokenizer.l. If the file is
// missing (e.g. CI without the legacy checkout), the tests skip.
const legacyLexerPath = "/Users/h3n4l/OpenSource/parser/googlesql/GoogleSQLLexer.g4"

// legacyParserPath is the legacy GoogleSQLParser.g4 — its
// common_keyword_as_identifier rule is the source of truth for the
// reserved/non-reserved keyword split (a word-keyword is non-reserved iff it
// appears in that rule, plus SIMPLE via keyword_as_identifier).
const legacyParserPath = "/Users/h3n4l/OpenSource/parser/googlesql/GoogleSQLParser.g4"

// keywordRuleRE matches a single-word keyword token rule in the lexer
// grammar, e.g.  SELECT_SYMBOL: 'SELECT';  capturing the literal text. It
// tolerates a newline between the ':' and the literal (the
// KW_MATCH_RECOGNIZE_NONRESERVED_SYMBOL rule is split across two lines).
var keywordRuleRE = regexp.MustCompile(`(?s)\b[A-Z_]+_SYMBOL\s*:\s*'([A-Z][A-Z0-9_]*)'\s*;`)

// extractLegacyKeywords returns the set of all word-keyword literals declared
// in GoogleSQLLexer.g4 (lower-cased).
func extractLegacyKeywords(t *testing.T) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(legacyLexerPath)
	if err != nil {
		t.Skipf("legacy lexer grammar not available at %s: %v", legacyLexerPath, err)
	}
	matches := keywordRuleRE.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatalf("no keyword tokens extracted from %s — regex broken?", legacyLexerPath)
	}
	set := make(map[string]bool, len(matches))
	for _, m := range matches {
		set[strings.ToLower(m[1])] = true
	}
	return set
}

// TestKeywordCompleteness asserts that every word-keyword token in
// GoogleSQLLexer.g4 has a corresponding entry in keywordMap, and that
// keywordMap has no extra entries the grammar does not declare.
func TestKeywordCompleteness(t *testing.T) {
	legacy := extractLegacyKeywords(t)

	missing := []string{}
	for kw := range legacy {
		if _, ok := keywordMap[kw]; !ok {
			missing = append(missing, kw)
		}
	}
	if len(missing) > 0 {
		t.Errorf("%d keywords in grammar missing from keywordMap: %v", len(missing), missing)
	}

	extra := []string{}
	for kw := range keywordMap {
		if !legacy[kw] {
			extra = append(extra, kw)
		}
	}
	if len(extra) > 0 {
		t.Errorf("%d keywords in keywordMap not declared in grammar: %v", len(extra), extra)
	}

	if len(keywordMap) != len(legacy) {
		t.Errorf("keywordMap size = %d, grammar word-keyword count = %d", len(keywordMap), len(legacy))
	}
	t.Logf("checked %d keywords from grammar; %d missing, %d extra", len(legacy), len(missing), len(extra))
}

// nonReservedRuleRE captures the body of the common_keyword_as_identifier
// rule from the parser grammar.
var commonKeywordRuleRE = regexp.MustCompile(`(?s)common_keyword_as_identifier:(.*?);`)

// extractLegacyNonReserved returns the set of non-reserved word-keyword
// literals (lower-cased): the common_keyword_as_identifier members plus
// SIMPLE (added by keyword_as_identifier).
func extractLegacyNonReserved(t *testing.T) map[string]bool {
	t.Helper()
	data, err := os.ReadFile(legacyParserPath)
	if err != nil {
		t.Skipf("legacy parser grammar not available at %s: %v", legacyParserPath, err)
	}
	m := commonKeywordRuleRE.FindStringSubmatch(string(data))
	if m == nil {
		t.Fatalf("common_keyword_as_identifier rule not found in %s", legacyParserPath)
	}
	tokenRE := regexp.MustCompile(`([A-Z_]+)_SYMBOL`)
	set := make(map[string]bool)
	for _, tm := range tokenRE.FindAllStringSubmatch(m[1], -1) {
		set[strings.ToLower(tm[1])] = true
	}
	// keyword_as_identifier adds SIMPLE on top of common_keyword_as_identifier.
	set["simple"] = true
	return set
}

// TestReservedKeywordSplit asserts that the lexer's reserved/non-reserved
// classification matches the legacy grammar exactly: a word-keyword is
// reserved iff it is NOT in common_keyword_as_identifier (+ SIMPLE).
func TestReservedKeywordSplit(t *testing.T) {
	allKw := extractLegacyKeywords(t)
	nonReserved := extractLegacyNonReserved(t)

	// Every non-reserved grammar keyword must be a known keyword and NOT
	// classified reserved.
	for kw := range nonReserved {
		if !allKw[kw] {
			// KW_MATCH_RECOGNIZE_NONRESERVED is a member; it is a real token.
			t.Errorf("non-reserved keyword %q is not a declared word-keyword", kw)
			continue
		}
		if IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = true, want false (it is in common_keyword_as_identifier)", kw)
		}
	}

	// Every reserved grammar keyword (allKw - nonReserved) must be classified
	// reserved.
	reservedCount := 0
	for kw := range allKw {
		if nonReserved[kw] {
			continue
		}
		reservedCount++
		if !IsReservedKeyword(kw) {
			t.Errorf("IsReservedKeyword(%q) = false, want true (not in common_keyword_as_identifier)", kw)
		}
	}

	t.Logf("grammar: %d word-keywords = %d non-reserved + %d reserved",
		len(allKw), len(nonReserved), reservedCount)
}
