package completion

import (
	"context"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
)

// identProbe is a bare identifier no keyword spells.
const identProbe = "omni_completion_ident"

// lookaheadContinuations lists, for each keyword CockroachDB's lexer rewrites
// by the tokens after it (lexer.go Lex: NOT before IN becomes NOT_LA, and so
// on), the continuations that select the rewritten token. A keyword probed at
// the end of the text is never rewritten, so these are probed too.
var lookaheadContinuations = map[string][]string{
	"as":        {"of system"},
	"not":       {"between", "in", "like", "ilike", "similar"},
	"generated": {"always", "by"},
	"with":      {"time", "ordinality", "bucket_count"},
	"nulls":     {"first", "last"},
	"reset":     {"all"},
	"role":      {"all"},
	"user":      {"all"},
	"on":        {"delete", "update no"},
	"tenant":    {"all"},
	"cluster":   {"all"},
	"set":       {"tracing"},
	"index":     {"(", "omni_completion_ident ("},
}

// lexerRewrites are KeywordNames entries that name tokens CockroachDB's lexer
// produces from INDEX and NOTHING by lookahead (lexer.go Lex); no one types
// them.
var lexerRewrites = map[string]bool{
	"index_after_order_by_before_at": true,
	"index_before_name_then_paren":   true,
	"index_before_paren":             true,
	"nothing_after_returning":        true,
}

// punctuationFollowers are the tokens probed after a column-name or
// type/function-name keyword and after an identifier to tell a keyword used
// as a keyword from one the grammar only accepts as a name. Every keyword's
// lookahead continuations and helpFollowers, the keywords that follow it in
// CockroachDB's own statement help, are probed too.
var punctuationFollowers = []string{
	identProbe, ".", "(", ")", ",", "=", "1", "'s'", "[", "@", "*", "::",
}

// keywordBudget caps the time keyword probing takes when ctx allows longer.
// Probing stops there and the keywords found so far are returned, so a long
// statement never blocks a keystroke.
const keywordBudget = 200 * time.Millisecond

// keywordCandidates returns the keywords CockroachDB's grammar accepts after
// prefix, by asking the vendor parser about prefix followed by each keyword.
//
// An unreserved keyword is also an identifier, so where an identifier fits
// every unreserved keyword parses. There a non-reserved keyword is offered
// only when it reads as a keyword: some follower token parses after it and
// not after an identifier, or the reverse.
//
// Reserved keywords are probed first, then column-name and type/function-name
// keywords, then unreserved ones, so a probe cut short by the budget keeps
// the keywords most statements need.
func keywordCandidates(ctx context.Context, prefix string) []string {
	prefix = probePrefix(prefix)
	if cached, ok := keywordCache.get(prefix); ok {
		return cached
	}
	ctx, cancel := context.WithTimeout(ctx, keywordBudget)
	defer cancel()
	result := probeKeywords(ctx, prefix)
	if ctx.Err() == nil {
		keywordCache.put(prefix, result)
	}
	return result
}

// probeKeywords returns the keywords found to fit after prefix before ctx
// is done.
func probeKeywords(ctx context.Context, prefix string) []string {
	base := prefix + " "
	ident := &identFollowers{text: base + identProbe, fits: map[string]bool{}}
	identFits := reachesPast(ident.text, len(base))

	order := keywordOrder()
	accepted := make([]bool, len(lexbase.KeywordNames))
	parallel(ctx, len(order), func(j int) {
		i := order[j]
		accepted[i] = keywordFits(ctx, base, lexbase.KeywordNames[i], identFits, ident)
	})

	var result []string
	for i, ok := range accepted {
		if ok {
			result = append(result, strings.ToUpper(lexbase.KeywordNames[i]))
		}
	}
	return result
}

// keywordFits reports whether kw is offered after base.
func keywordFits(ctx context.Context, base, kw string, identFits bool, ident *identFollowers) bool {
	start := len(base)
	text := base + kw
	conts, rewritten := lookaheadContinuations[kw]
	category := lexbase.KeywordsCategories[kw]
	if !identFits || category == "R" || rewritten {
		fits := reachesPast(text, start)
		for _, cont := range conts {
			fits = fits || reachesPast(text+" "+cont, start)
		}
		if !fits || !identFits || category == "R" {
			return fits
		}
	}
	followers := slices.Concat(conts, helpFollowers()[kw])
	if category != "U" {
		// Column-name and type/function-name keywords are names only in
		// some positions; punctuation and completeness tell them apart.
		followers = slices.Concat([]string{""}, punctuationFollowers, followers)
	}
	for _, f := range followers {
		if ctx.Err() != nil {
			return false
		}
		kwFits, followerFits := probeFollower(text, start, f)
		if !kwFits && !rewritten {
			// The lexer reads kw the same before any follower, so the
			// parser rejected kw itself.
			return false
		}
		if followerFits != ident.follows(f) {
			return true
		}
	}
	return false
}

// keywordOrder lists the indexes of lexbase.KeywordNames in probing order:
// reserved, then column-name and type/function-name, then unreserved.
var keywordOrder = sync.OnceValue(func() []int {
	rank := map[string]int{"R": 0, "C": 1, "T": 1, "U": 2}
	order := make([]int, len(lexbase.KeywordNames))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		return rank[lexbase.KeywordsCategories[lexbase.KeywordNames[a]]] - rank[lexbase.KeywordsCategories[lexbase.KeywordNames[b]]]
	})
	return slices.DeleteFunc(order, func(i int) bool { return lexerRewrites[lexbase.KeywordNames[i]] })
})

// identFollowers memoizes which followers parse after an identifier at the
// caret.
type identFollowers struct {
	text string
	mu   sync.Mutex
	fits map[string]bool
}

func (f *identFollowers) follows(follower string) bool {
	f.mu.Lock()
	ok, done := f.fits[follower]
	f.mu.Unlock()
	if done {
		return ok
	}
	ok = followerFits(f.text, follower)
	f.mu.Lock()
	f.fits[follower] = ok
	f.mu.Unlock()
	return ok
}

// followerFits reports whether follower parses after text; the empty follower
// asks whether text is a complete statement.
func followerFits(text, follower string) bool {
	_, ok := probeFollower(text, len(text), follower)
	return ok
}

// probeFollower parses text followed by follower, and reports whether the
// parser shifted the token at byte offset tokStart of text and whether it
// shifted follower; the empty follower asks whether text is complete.
func probeFollower(text string, tokStart int, follower string) (tokFits, followerFits bool) {
	full := text
	if follower != "" {
		full += " " + follower
	}
	_, err := parseOne(full)
	if err == nil {
		return true, true
	}
	pos := errorOffset(full, err)
	return pos > tokStart, follower != "" && pos > len(text)+1
}

// helpFollowers maps each keyword to the keywords that follow it in
// CockroachDB's statement help (parser.HelpMessages, generated from sql.y),
// such as IF to NOT in "CREATE TABLE [IF NOT EXISTS]".
var helpFollowers = sync.OnceValue(func() map[string][]string {
	out := make(map[string][]string)
	add := func(prev, next string) {
		if !slices.Contains(out[prev], next) {
			out[prev] = append(out[prev], next)
		}
	}
	for _, msg := range crparser.HelpMessages {
		prev := ""
		for _, word := range helpWord.FindAllString(msg.Text, -1) {
			switch {
			case strings.ContainsAny(word, "[]{}|"):
				// Optional and alternative groups do not break a keyword run.
			case isHelpKeyword(word):
				kw := strings.ToLower(word)
				if prev != "" {
					add(prev, kw)
				}
				prev = kw
			default:
				prev = ""
			}
		}
	}
	return out
})

// helpWord splits CockroachDB's help text into words, <placeholders>, and
// single punctuation characters.
var helpWord = regexp.MustCompile(`[A-Za-z_][A-Za-z_0-9]*|<[^>]*>|\S`)

func isHelpKeyword(word string) bool {
	if word != strings.ToUpper(word) {
		return false
	}
	_, ok := lexbase.KeywordsCategories[strings.ToLower(word)]
	return ok
}

// parallel runs fn(0..n-1) on a bounded worker pool, stopping early once ctx
// is done.
func parallel(ctx context.Context, n int, fn func(i int)) {
	workers := min(runtime.GOMAXPROCS(0), n)
	if workers <= 1 {
		for i := range n {
			if ctx.Err() != nil {
				return
			}
			fn(i)
		}
		return
	}
	next := make(chan int)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				fn(i)
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
}

// keywordCache keeps the keywords of the last few prefixes. While a word is
// typed the text before it does not change, so each keystroke after the
// first is answered without parsing.
var keywordCache = &prefixCache{limit: 32}

type prefixCache struct {
	mu      sync.Mutex
	limit   int
	order   []string
	entries map[string][]string
}

func (c *prefixCache) get(prefix string) ([]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.entries[prefix]
	return v, ok
}

func (c *prefixCache) put(prefix string, keywords []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.entries == nil {
		c.entries = make(map[string][]string)
	}
	if _, ok := c.entries[prefix]; ok {
		return
	}
	if len(c.order) == c.limit {
		delete(c.entries, c.order[0])
		c.order = c.order[1:]
	}
	c.order = append(c.order, prefix)
	c.entries[prefix] = keywords
}

func (c *prefixCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.order, c.entries = nil, nil
}
