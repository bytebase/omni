package parser

import (
	"slices"
	"strings"
)

// ExpectedKeywords returns the keywords and words Snowflake's grammar accepts
// at the end of prefix, the start of a statement cut where a caret is. Each is
// upper case, and the list is sorted.
//
// The parser reads prefix in collection mode: every token test it makes while
// the current token is the end of input (curIs, nextIs and the word tests
// below) records what it tested for. That is the set of tokens the parser
// would take there, read from the parser itself rather than from a probe per
// keyword. A word the parser matches by text, such as EXECUTE, is recorded as
// that word.
func ExpectedKeywords(prefix string) []string {
	w := &expectation{tokens: map[int]bool{}, words: map[string]bool{}}
	p := &Parser{lexer: NewLexer(prefix), input: prefix, wants: w}
	p.advance()
	if p.cur.Type == tokEOF {
		w.atEOF = true
	}
	_, _ = p.parseTopStmt() // the error is expected: prefix is cut short
	if !w.atEOF {
		return nil
	}
	seen := map[string]bool{}
	var out []string
	add := func(word string) {
		if word != "" && !seen[word] {
			seen[word] = true
			out = append(out, word)
		}
	}
	for t := range w.tokens {
		if name, ok := keywordName(t); ok {
			add(name)
		}
	}
	for word := range w.words {
		add(word)
	}
	slices.Sort(out)
	return out
}

// expectation collects the tokens and words tested at the end of input.
// atEOF records that the parser reached the end of input at all: a prefix
// that fails earlier says nothing about what follows it.
type expectation struct {
	tokens map[int]bool
	words  map[string]bool
	atEOF  bool
}

// collecting reports whether tok is the end of input of a collection-mode
// parse, the place whose expectations are recorded.
func (p *Parser) collecting(tok Token) bool {
	if p.wants == nil || tok.Type != tokEOF {
		return false
	}
	p.wants.atEOF = true
	return true
}

// curIs reports whether the current token has type t. Parser code tests the
// current token through it, never by comparing p.cur.Type, so collection mode
// sees every test (expected_test.go enforces this).
func (p *Parser) curIs(t int) bool {
	if p.cur.Type == t {
		return true
	}
	if p.collecting(p.cur) {
		p.wants.tokens[t] = true
	}
	return false
}

// nextIs reports whether the token after the current one has type t. The
// test is recorded only when the current token is real, so the token after
// the end of input is never mistaken for one at it.
func (p *Parser) nextIs(t int) bool {
	next := p.peekNext()
	if next.Type == t {
		return true
	}
	if p.cur.Type != tokEOF && p.collecting(next) {
		p.wants.tokens[t] = true
	}
	return false
}

// wantWords records words as tested at the current token, for code that
// dispatches on a word's text through a map or a switch on that text.
func (p *Parser) wantWords(words ...string) {
	if p.collecting(p.cur) {
		for _, w := range words {
			p.wants.words[strings.ToUpper(w)] = true
		}
	}
}

// wantNextWords is wantWords for the token after the current one.
func (p *Parser) wantNextWords(words ...string) {
	if p.cur.Type != tokEOF && p.collecting(p.peekNext()) {
		for _, w := range words {
			p.wants.words[strings.ToUpper(w)] = true
		}
	}
}

// curWord is isWordToken(p.cur, upper), recorded in collection mode.
func (p *Parser) curWord(upper string) bool {
	if isWordToken(p.cur, upper) {
		return true
	}
	p.wantWords(upper)
	return false
}

// nextWord is isWordToken(p.peekNext(), upper), recorded in collection mode.
func (p *Parser) nextWord(upper string) bool {
	if isWordToken(p.peekNext(), upper) {
		return true
	}
	p.wantNextWords(upper)
	return false
}

// curNameWord is tokIsWord(p.cur, upper), recorded in collection mode.
func (p *Parser) curNameWord(upper string) bool {
	if tokIsWord(p.cur, upper) {
		return true
	}
	p.wantWords(upper)
	return false
}

// nextNameWord is tokIsWord(p.peekNext(), upper), recorded in collection
// mode.
func (p *Parser) nextNameWord(upper string) bool {
	if tokIsWord(p.peekNext(), upper) {
		return true
	}
	p.wantNextWords(upper)
	return false
}
