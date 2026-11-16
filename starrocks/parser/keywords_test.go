package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// keywordRule matches one keyword rule of StarRocksLex.g4, capturing its name
// and the word it matches: `ACCOUNT_LOCK : 'ACCOUNT_LOCK';`. A rule that
// matches anything but a word (an operator, a quoted literal) is not a
// keyword and does not match.
var keywordRule = regexp.MustCompile(`(?m)^([A-Z][A-Z0-9_]*)\s*:\s*'([A-Za-z_][A-Za-z0-9_]*)'\s*;\s*$`)

// nonReservedRule matches StarRocks.g4's nonReserved rule and captures its
// body, a `|`-separated list of token names.
var nonReservedRule = regexp.MustCompile(`(?ms)^nonReserved\s*\n(.*?)\n\s*;`)

// wordToken is a token of StarRocks.g4 that nonReserved names and that no
// keyword rule matches: [*] and ..., whose literals are not words.
var wordlessTokens = []string{"ARRAY_ELEMENT", "DOTDOTDOT"}

// TestKeywordTableMatchesGrammar derives both keyword tables from the pinned
// release's own grammars in ../testdata and holds them to it. keywordMap is
// StarRocksLex.g4's keyword rules, one token per word; nonReservedKeywords is
// StarRocks.g4's nonReserved rule. A keyword the engine adds, drops or
// un-reserves is a failure here rather than an over-acceptance nobody sees.
//
// dorisOnlyKeywords is the exception the test states instead of hiding: every
// word in it is absent from the grammars and non-reserved, so an identifier
// use of it reads as the engine reads it, and each statement behind it is
// pinned in testdata/compat.tsv as a row only omni accepts.
func TestKeywordTableMatchesGrammar(t *testing.T) {
	rules, nonReserved := readGrammars(t)

	if len(keywordMap) != len(rules) {
		t.Errorf("keywordMap holds %d keywords; StarRocksLex.g4 has %d", len(keywordMap), len(rules))
	}
	kinds := make(map[TokenKind]string, len(rules))
	for word, rule := range rules {
		kind, ok := keywordMap[strings.ToLower(word)]
		if !ok {
			t.Errorf("%s: StarRocksLex.g4 has the keyword, keywordMap does not", word)
			continue
		}
		if prev, dup := kinds[kind]; dup {
			t.Errorf("%s and %s share one token: %s is a keyword of its own", prev, word, rule)
		}
		kinds[kind] = word
		if name := TokenName(kind); name != word {
			t.Errorf("%s: TokenName = %q, want the word the grammar matches", word, name)
		}
		wantReserved := !nonReserved[rule]
		if got := IsReserved(kind); got != wantReserved {
			t.Errorf("%s: IsReserved = %t; StarRocks.g4's nonReserved says %t", word, got, !wantReserved)
		}
	}
	for word := range keywordMap {
		if _, ok := rules[strings.ToUpper(word)]; !ok {
			t.Errorf("keywordMap holds %q, which StarRocksLex.g4 has no keyword for", word)
		}
	}

	for rule := range nonReserved {
		if _, ok := rules[strings.ToUpper(rule)]; ok {
			continue
		}
		if !slices.Contains(wordlessTokens, rule) {
			t.Errorf("StarRocks.g4's nonReserved names %s, and no keyword matches it", rule)
		}
	}

	words := make([]string, 0, len(dorisOnlyKeywords))
	for word := range dorisOnlyKeywords {
		words = append(words, word)
	}
	slices.Sort(words)
	for _, word := range words {
		kind := dorisOnlyKeywords[word]
		if _, ok := rules[strings.ToUpper(word)]; ok {
			t.Errorf("dorisOnlyKeywords holds %q, which StarRocksLex.g4 makes a keyword", word)
		}
		if IsReserved(kind) {
			t.Errorf("%q is not a keyword of the engine, so it must be non-reserved", word)
		}
		if _, ok := kinds[kind]; ok {
			t.Errorf("%q shares a token with a StarRocks keyword: %s", word, kinds[kind])
		}
		if _, dup := keywordMap[word]; dup {
			t.Errorf("%q is in keywordMap and dorisOnlyKeywords", word)
		}
	}
}

// readGrammars reads the pinned grammars and returns StarRocksLex.g4's keyword
// rules by word, and StarRocks.g4's nonReserved names.
func readGrammars(t *testing.T) (rules map[string]string, nonReserved map[string]bool) {
	t.Helper()
	read := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join("..", "testdata", name))
		if err != nil {
			t.Fatalf("%v", err)
		}
		return string(b)
	}
	rules = map[string]string{}
	for _, m := range keywordRule.FindAllStringSubmatch(read("StarRocksLex.g4"), -1) {
		if _, dup := rules[m[2]]; dup {
			t.Fatalf("StarRocksLex.g4 matches the word %s twice", m[2])
		}
		rules[m[2]] = m[1]
	}
	if len(rules) < 400 {
		t.Fatalf("read %d keyword rules; StarRocksLex.g4 has hundreds", len(rules))
	}
	body := nonReservedRule.FindStringSubmatch(read("StarRocks.g4"))
	if body == nil {
		t.Fatal("StarRocks.g4 has no nonReserved rule")
	}
	nonReserved = map[string]bool{}
	for _, name := range regexp.MustCompile(`[A-Z][A-Z0-9_]*`).FindAllString(body[1], -1) {
		nonReserved[name] = true
	}
	if len(nonReserved) < 300 {
		t.Fatalf("read %d nonReserved names; StarRocks.g4 lists hundreds", len(nonReserved))
	}
	return rules, nonReserved
}
