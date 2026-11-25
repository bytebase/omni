package parser

import (
	"context"
	"errors"
	"maps"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
)

// ExpectedKeywords returns the keywords ClickHouse's grammar accepts at the
// end of prefix, the start of a statement cut where a caret is. Each is upper
// case, and the list is sorted. It stops when ctx is done and returns the
// keywords found by then.
//
// ClickHouse reads almost every keyword as an identifier too, so a keyword
// fits at the caret only where the parser takes it as that keyword. Each
// keyword is tried after prefix with the parser told not to read it as an
// identifier; it fits when the parse gets past it. A keyword the parser
// takes only before particular tokens, as GROUP before BY, CAST before ( or
// IF before NOT EXISTS, fails alone, but the parser then looked past it; up
// to two followers are tried after it. A word the parser matches by text,
// such as GROUPS or OPTION, is not in the keyword table: the parser records
// the words it tests at the end of prefix.
func ExpectedKeywords(ctx context.Context, prefix string) []string {
	pr, ok := newProber(prefix)
	if !ok || ctx.Err() != nil {
		return nil
	}
	p := pr.parser("", false)
	p.wants = map[string]bool{}
	_, _ = p.parseStatement() // the error is expected: prefix is cut short
	found := map[string]bool{}
	for word := range p.wants {
		found[word] = true
	}

	var words []string
	for word := range keywords {
		if !found[word] {
			words = append(words, word)
		}
	}
	slices.Sort(words)
	followers := append(slices.Clone(words), nonKeywordFollowers...)
	// Each level holds the word sequences the parser looked past without
	// taking; a keyword fits once one sequence starting with it does.
	var level []sequence
	generic := pr.fits(sentinel).lookedPast
	for i, r := range pr.each(ctx, len(words), func(i int) probeResult { return pr.fits(words[i]) }) {
		switch {
		case r.fits:
			found[words[i]] = true
		case r.lookedPast && !generic:
			level = append(level, sequence{head: len(level), words: words[i], extra: r.words})
		}
	}
	heads := make([]string, len(level))
	for i, seq := range level {
		heads[i] = seq.words
	}
	headFits := make([]atomic.Bool, len(heads))
	for range maxFollowers {
		var mu sync.Mutex
		var next []sequence
		generic := make([]bool, len(level))
		for i, r := range pr.each(ctx, len(level), func(i int) probeResult { return pr.fits(level[i].words + " " + sentinel) }) {
			generic[i] = r.lookedPast
		}
		// Each sequence is tried with every keyword and generic token, and
		// with the words the parser matched by text after it, as OPTION
		// after WITH GRANT.
		type try struct {
			seq      int
			follower string
		}
		var tries []try
		for i, seq := range level {
			for _, f := range followers {
				tries = append(tries, try{i, f})
			}
			for _, f := range seq.extra {
				tries = append(tries, try{i, f})
			}
		}
		pr.each(ctx, len(tries), func(i int) probeResult {
			seq := level[tries[i].seq]
			if headFits[seq.head].Load() {
				return probeResult{}
			}
			longer := sequence{head: seq.head, words: seq.words + " " + tries[i].follower}
			switch r := pr.fits(longer.words); {
			case r.fits:
				headFits[seq.head].Store(true)
			case r.lookedPast && !generic[tries[i].seq]:
				longer.extra = r.words
				mu.Lock()
				next = append(next, longer)
				mu.Unlock()
			}
			return probeResult{}
		})
		level = next
	}
	for h, head := range heads {
		if headFits[h].Load() {
			found[head] = true
		}
	}

	out := make([]string, 0, len(found))
	for word := range found {
		out = append(out, word)
	}
	slices.Sort(out)
	return out
}

// sentinel is a token the grammar takes after no keyword. When the parser
// looks past it too, looking past a keyword at the same place says nothing
// about that keyword.
const sentinel = "]"

// maxFollowers is how many tokens are tried after a keyword the parser looks
// past: IF NOT EXISTS needs two.
const maxFollowers = 2

// nonKeywordFollowers are the tokens other than keywords tried after a
// keyword: a parenthesis, a number, a string and a name, as after CAST,
// INTERVAL, TIMESTAMP and CONSTRAINT.
var nonKeywordFollowers = []string{"(", "1", "''", "omni_name"}

// sequence is a keyword followed by the tokens tried after it.
type sequence struct {
	head  int // the keyword's index in heads
	words string
	extra []string // the words the parser matched by text after it
}

// prober tries words after a prefix, sharing the prefix's tokens.
type prober struct {
	prefix string
	tokens []Token // the prefix's tokens, without its EOF
}

// newProber lexes prefix. It fails when a word typed after prefix would not
// start a token of its own: prefix does not lex, ends inside a comment, or
// ends in a character a word continues.
func newProber(prefix string) (*prober, bool) {
	if prefix != "" && isIdentPart(prefix[len(prefix)-1]) {
		return nil, false
	}
	toks, err := Tokenize(prefix + " x")
	if err != nil || len(toks) < 2 {
		return nil, false
	}
	if last := toks[len(toks)-2]; last.Type != TokIdent || last.Pos != len(prefix)+1 {
		return nil, false
	}
	return &prober{prefix: prefix, tokens: toks[:len(toks)-2]}, true
}

// parser returns a parser over prefix followed by words, its tokens up to
// the end of prefix already lexed. probing marks the words keyword-only.
func (pr *prober) parser(words string, probing bool) *parser {
	src := pr.prefix + words
	lex := &Lexer{src: src, pos: len(pr.prefix)}
	if n := len(pr.tokens); n > 0 {
		lex.prev = pr.tokens[n-1]
	}
	return &parser{
		src:     src,
		lex:     lex,
		tokens:  slices.Clone(pr.tokens), // consumeRawTail rewrites tokens in place
		probing: probing,
		probeAt: len(pr.prefix),
	}
}

type probeResult struct {
	fits       bool
	lookedPast bool
	words      []string // the words the parser matched by text at the end
}

// fits reports whether the parser reads words, a keyword and the tokens
// tried after it, right after the prefix, the keywords as keywords, and goes
// on past them. lookedPast reports that the parser read the end of input
// after them, which a test of the token after a keyword does.
func (pr *prober) fits(words string) probeResult {
	p := pr.parser(words, true)
	p.wants = map[string]bool{}
	_, err := p.parseStatement()
	if len(p.tokens) <= len(pr.tokens) {
		// Nothing after the prefix was lexed, or an INSERT's FORMAT data
		// took the rest as raw text.
		return probeResult{}
	}
	lexed := p.tokens[len(pr.tokens):]
	var pe *ParseError
	last := lexed[len(lexed)-1]
	return probeResult{
		// The first word must lex as a keyword after the prefix.
		// A speculative parse that got past the keyword and failed later,
		// as an INDEX definition in a column list does, still read it.
		fits:       lexed[0].IsKeyword() && (err == nil || errors.As(err, &pe) && pe.Position > len(pr.prefix) || p.furthest > len(pr.prefix)),
		lookedPast: last.Type == TokEOF && last.Pos == len(p.src),
		words:      slices.Sorted(maps.Keys(p.wants)),
	}
}

// each runs probe for 0..n-1 on every CPU and returns the results in
// order. A probe not started when ctx is done is left zero.
func (pr *prober) each(ctx context.Context, n int, probe func(int) probeResult) []probeResult {
	results := make([]probeResult, n)
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), n) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				results[i] = probe(i)
			}
		}()
	}
	for i := range n {
		if ctx.Err() != nil {
			break
		}
		next <- i
	}
	close(next)
	wg.Wait()
	return results
}

// parseStatement parses one statement as Parse does, without the checks
// that follow it.
func (p *parser) parseStatement() (Statement, error) {
	if p.cur().Type == TokEOF || p.cur().Type == TokIllegal {
		return nil, p.errAt(p.cur(), "empty query")
	}
	stmt, err := p.parseTopLevel()
	if err != nil {
		return nil, err
	}
	if p.cur().Type == TokSemicolon {
		p.advance()
	}
	if p.cur().Type != TokEOF {
		return nil, p.errAt(p.cur(), "unexpected token")
	}
	return stmt, nil
}
