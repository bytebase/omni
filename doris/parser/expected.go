package parser

import (
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
)

// probeName is the identifier ExpectedKeywords probes a name position with.
const probeName = "omni_probe_name"

// probeFollowers are the punctuation a keyword the parser decides on by the
// token after it is retried with, beside every keyword.
var probeFollowers = []string{"(", ")", ",", ".", "=", "*", "<", "[", "'s'", "1"}

// ExpectedKeywords returns the keywords Doris's grammar accepts at the end of
// prefix, the start of a statement cut where a caret is. Each is upper case,
// and the list is sorted.
//
// It asks the parser itself, by trial: prefix is parsed once with each
// keyword appended, and a keyword fits when the parse succeeds or fails past
// it. A recursive-descent parser fails at the first token no production
// takes, so a failure at the keyword itself means the grammar has no place
// for it there. Two cases need more than that one parse:
//
//   - The parser takes some keywords only by the token after them (LEFT
//     before JOIN or OUTER). When a keyword fails at itself after the parser
//     looked past it, it is retried with each keyword and punctuation after
//     it, and fits when one of them parses.
//   - Where a name fits, every non-reserved keyword parses as that name.
//     Such a keyword is offered only when the parser reads it otherwise than
//     a name: its parse ends differently from the same prefix with a plain
//     identifier, or, when the parser looked past it, some follower parses
//     after it and not after the identifier (TABLET before '(').
//
// The parser reads a parenthesized query as raw text and parses it only
// once the statement around it is complete, so a prefix that ends inside
// one is judged by its closing parenthesis. Pass the query's own text
// instead: everything after the '(' that opens it.
func ExpectedKeywords(prefix string) []string {
	sep := ""
	if prefix != "" && !strings.ContainsAny(prefix[len(prefix)-1:], " \t\r\n") {
		sep = " "
	}
	words := make([]string, 0, len(keywordMap))
	for word := range keywordMap {
		words = append(words, word)
	}
	slices.Sort(words)

	// Each probe parses the whole prefix, so the probes run on every CPU:
	// first one per keyword, then one per follower of each keyword the
	// parser looked past.
	name := probe(prefix, sep, probeName, "")
	fits := make([]bool, len(words))
	asName := make([]bool, len(words))
	var retry []int
	var mu sync.Mutex
	parallel(len(words), func(i int) {
		r := probe(prefix, sep, words[i], "")
		asName[i] = name.fits && !IsReserved(keywordMap[words[i]])
		switch {
		case r.fits && !(asName[i] && r.sameEnd(name)):
			fits[i] = true
		case r.peeked:
			mu.Lock()
			retry = append(retry, i)
			mu.Unlock()
		}
	})

	// A word the parser decided on by the token after it fits when a
	// follower parses after it. Where a name fits too, the follower must not
	// parse after a name, and only punctuation is tried: after a name most
	// keywords parse as well, and the parser looks past most words there.
	var nameTakes []bool
	if slices.ContainsFunc(retry, func(i int) bool { return asName[i] }) {
		nameTakes = make([]bool, len(probeFollowers))
		parallel(len(probeFollowers), func(j int) {
			nameTakes[j] = probe(prefix, sep, probeName, probeFollowers[j]).pastFollower
		})
	}
	type pair struct{ word, follower int }
	var pairs []pair
	for _, i := range retry {
		n := len(followerWords)
		if asName[i] {
			n = len(probeFollowers)
		}
		for j := range n {
			pairs = append(pairs, pair{i, j})
		}
	}
	found := make([]atomic.Bool, len(words))
	parallel(len(pairs), func(k int) {
		i, j := pairs[k].word, pairs[k].follower
		if found[i].Load() {
			return
		}
		r := probe(prefix, sep, words[i], followerWords[j])
		if asName[i] && r.pastFollower && !nameTakes[j] || !asName[i] && r.fits {
			found[i].Store(true)
		}
	})

	var out []string
	for i, word := range words {
		if fits[i] || found[i].Load() {
			out = append(out, strings.ToUpper(word))
		}
	}
	return out
}

// parallel runs f(0), ..., f(n-1) on every CPU.
func parallel(n int, f func(int)) {
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), n) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				f(i)
			}
		}()
	}
	for i := range n {
		next <- i
	}
	close(next)
	wg.Wait()
}

// followerWords are the followers a keyword is retried with: punctuation,
// then every keyword.
var followerWords = func() []string {
	out := slices.Clone(probeFollowers)
	for word := range keywordMap {
		out = append(out, word)
	}
	slices.Sort(out[len(probeFollowers):])
	return out
}()

// wordProbe is the state a probe parse records: where the appended word
// starts, and whether the parser looked at the token after it while that
// token was the end of input.
type wordProbe struct {
	start  int
	peeked bool
}

// probeResult is how the parse of a prefix with one word appended ended.
type probeResult struct {
	fits         bool // the parse got past the word
	pastFollower bool // the parse got past the follower
	peeked       bool
	// ok is a parse with no error; otherwise fromEnd is how far before the
	// end of the text the first error is, and message its message.
	ok      bool
	fromEnd int
	message string
}

// sameEnd reports whether two probes of different words ended the same way.
func (r probeResult) sameEnd(o probeResult) bool {
	return r.ok == o.ok && r.fromEnd == o.fromEnd && r.message == o.message
}

// probe parses prefix with word, and follower when not empty, appended.
func probe(prefix, sep, word, follower string) probeResult {
	text := prefix + sep + word
	followerStart := len(text) + 1
	if follower != "" {
		text += " " + follower
	}
	w := &wordProbe{start: len(prefix) + len(sep)}
	p := &Parser{lexer: NewLexer(text), input: text, strictTrailing: true, probe: w}
	_, errs := p.parseSegment()
	if len(errs) == 0 {
		return probeResult{fits: true, pastFollower: true, peeked: w.peeked, ok: true}
	}
	first := errs[0]
	for _, e := range errs[1:] {
		if e.Position < first.Position {
			first = e
		}
	}
	return probeResult{
		fits:         first.Position > w.start,
		pastFollower: first.Position > followerStart,
		peeked:       w.peeked,
		fromEnd:      len(text) - first.Position,
		message:      first.Message,
	}
}
