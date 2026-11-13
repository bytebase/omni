package grammar

import (
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// A Recognizer decides whether a grammar derives a statement: the judge of
// an engine whose answers come from its grammar or its docs (§2.3 of
// docs/specs/2026-11-11-engine-compatibility-without-servers-design.md).
//
// It reads the statement with a generic SQL lexer and matches the grammar's
// terminals by their samples. A literal, or a token sampled by a keyword
// (an all-caps word) or by several tokens, matches only that text, words
// without regard to case. A token sampled by one name, string, number or
// variable stands for its class: it matches any token of that kind, and a
// name class takes no word the token table marks reserved. A name sampled
// in quotes is the class of names in those quotes, and no bare word. So it knows no
// more of the engine's lexer than the samples show: a hint or a comment the
// lexer reads as a token is skipped as a comment.
type Recognizer struct {
	g      *Grammar
	rules  []crule
	terms  []cterm
	start  int
	punct  [][]rune              // the grammar's punctuation, longest first
	exact  map[string][]spelling // exact terminal texts by the key of their first token
	class  map[string][]int      // class terminals by the class they match (classOf)
	null   []bool                // per rule: derives the empty string
	first  []bitset              // per rule: the terminals it can start with
	altFst [][]bitset            // per rule and alternative

	mu        sync.Mutex
	sets      []*earleySet
	predicted []int64 // per rule: the stamp of the set it was last predicted in
	stamp     int64
}

type crule struct {
	name string
	alts [][]csym
}

// A csym is a nonterminal (rule >= 0) or a terminal.
type csym struct {
	rule int
	term int
}

type cterm struct {
	toks  []ltok   // the text it matches; empty for a token written as nothing
	class string   // when set, it matches one token of this class (classOf)
	more  [][]ltok // other texts it matches exactly, for a token of several spellings
}

// spelling is one exact text of a terminal: its toks, or one of more.
type spelling struct{ term, alt int }

func (t cterm) spelling(alt int) []ltok {
	if alt == 0 {
		return t.toks
	}
	return t.more[alt-1]
}

type bitset []uint64

func (b bitset) has(i int) bool { return b[i/64]&(1<<(uint(i)%64)) != 0 }
func (b bitset) set(i int) bool {
	w, m := i/64, uint64(1)<<(uint(i)%64)
	if b[w]&m != 0 {
		return false
	}
	b[w] |= m
	return true
}

func (b bitset) union(o bitset) bool {
	changed := false
	for i := range b {
		if n := b[i] | o[i]; n != b[i] {
			b[i] = n
			changed = true
		}
	}
	return changed
}

// Recognizer returns a recognizer for the grammar.
func (g *Grammar) Recognizer() *Recognizer {
	r := &Recognizer{g: g, exact: map[string][]spelling{}, class: map[string][]int{}}
	ruleID := map[string]int{}
	for i, gr := range g.rules {
		ruleID[gr.name] = i
	}
	r.start = ruleID[g.start]

	// The punctuation the grammar spells, so the lexer reads => as one
	// token where the grammar has it, and as = and > where it does not.
	puncts := map[string]bool{}
	for _, gr := range g.rules {
		for _, alt := range gr.alts {
			for _, sym := range alt {
				if !sym.nonterm {
					for _, piece := range punctPieces(unquoted(g.sample(sym))) {
						puncts[piece] = true
					}
				}
			}
		}
	}
	for p := range puncts {
		r.punct = append(r.punct, []rune(p))
	}
	sort.Slice(r.punct, func(i, j int) bool {
		if len(r.punct[i]) != len(r.punct[j]) {
			return len(r.punct[i]) > len(r.punct[j])
		}
		return string(r.punct[i]) < string(r.punct[j])
	})

	termID := map[string]int{}
	term := func(sym symbol) int {
		key := sym.text
		if sym.literal {
			key = "'" + key
		}
		if id, ok := termID[key]; ok {
			return id
		}
		id := len(r.terms)
		termID[key] = id
		if !sym.literal && g.unreachable[sym.text] {
			// A token the lexer never produces matches nothing: one bad
			// token, indexed nowhere, so it is neither read nor empty.
			r.terms = append(r.terms, cterm{toks: []ltok{{kind: kBad}}})
			return id
		}
		toks := r.lex(g.sample(sym))
		t := cterm{toks: toks}
		spellings := g.spellings[sym.text]
		if !sym.literal && len(spellings) > 0 {
			// A token of a few fixed texts matches each of them, and nothing
			// else: not the class its sample would stand for.
			t.toks = r.lex(spellings[0])
			for _, s := range spellings[1:] {
				t.more = append(t.more, r.lex(s))
			}
		} else if !sym.literal && !g.fixed[sym.text] && len(toks) == 1 {
			switch tk := toks[0]; tk.kind {
			case kName:
				if tk.quoted || g.patterned[sym.text] || !keywordText.MatchString(tk.text) {
					t.class = classOf(tk)
				}
			case kString, kNumber, kVar:
				t.class = classOf(tk)
			}
		}
		r.terms = append(r.terms, t)
		switch {
		case t.class != "":
			r.class[t.class] = append(r.class[t.class], id)
		default:
			for a := 0; a <= len(t.more); a++ {
				if s := t.spelling(a); len(s) > 0 {
					r.exact[s[0].key()] = append(r.exact[s[0].key()], spelling{id, a})
				}
			}
		}
		return id
	}
	r.rules = make([]crule, len(g.rules))
	for i, gr := range g.rules {
		cr := crule{name: gr.name}
		for _, alt := range gr.alts {
			var ca []csym
			for _, sym := range alt {
				if sym.nonterm {
					ca = append(ca, csym{rule: ruleID[sym.text], term: -1})
				} else {
					ca = append(ca, csym{rule: -1, term: term(sym)})
				}
			}
			cr.alts = append(cr.alts, ca)
		}
		r.rules[i] = cr
	}
	r.analyze()
	r.predicted = make([]int64, len(r.rules))
	return r
}

// analyze computes which rules derive the empty string and which
// terminals each rule and alternative can start with.
func (r *Recognizer) analyze() {
	n := len(r.rules)
	words := (len(r.terms) + 63) / 64
	r.null = make([]bool, n)
	symNull := func(s csym) bool {
		if s.rule >= 0 {
			return r.null[s.rule]
		}
		return len(r.terms[s.term].toks) == 0
	}
	for changed := true; changed; {
		changed = false
		for i, cr := range r.rules {
			if r.null[i] {
				continue
			}
			for _, alt := range cr.alts {
				all := true
				for _, s := range alt {
					if !symNull(s) {
						all = false
						break
					}
				}
				if all {
					r.null[i], changed = true, true
					break
				}
			}
		}
	}
	r.first = make([]bitset, n)
	for i := range r.first {
		r.first[i] = make(bitset, words)
	}
	// altFirst adds an alternative's first terminals to dst.
	altFirst := func(dst bitset, alt []csym) bool {
		changed := false
		for _, s := range alt {
			if s.rule >= 0 {
				if dst.union(r.first[s.rule]) {
					changed = true
				}
			} else if len(r.terms[s.term].toks) > 0 && dst.set(s.term) {
				changed = true
			}
			if !symNull(s) {
				break
			}
		}
		return changed
	}
	for changed := true; changed; {
		changed = false
		for i, cr := range r.rules {
			for _, alt := range cr.alts {
				if altFirst(r.first[i], alt) {
					changed = true
				}
			}
		}
	}
	r.altFst = make([][]bitset, n)
	for i, cr := range r.rules {
		for _, alt := range cr.alts {
			b := make(bitset, words)
			altFirst(b, alt)
			r.altFst[i] = append(r.altFst[i], b)
		}
	}
}

// An item is an Earley item: alternative alt of rule, read up to dot,
// begun at origin.
type item struct {
	rule, alt, dot, origin int32
}

// An earleySet is the items at one position. Sets are kept and reused
// from one statement to the next, which is most of what makes the
// recognizer fast enough to judge every row on every run.
type earleySet struct {
	items []item
	seen  map[item]struct{}
	// head and chain hold, per rule, the items whose next symbol is that
	// rule: head[rule] is the last one's index in chain, and chain[i].next
	// the one before it, or -1.
	head  map[int32]int32
	chain []waiter
}

type waiter struct {
	it   item
	next int32
}

// Accepts reports whether the grammar derives the statement from its
// start symbol. A Recognizer reads one statement at a time.
func (r *Recognizer) Accepts(sql string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	in := r.lex(sql)
	n := len(in)
	for len(r.sets) < n+1 {
		r.sets = append(r.sets, &earleySet{seen: map[item]struct{}{}, head: map[int32]int32{}})
	}
	sets := r.sets[:n+1]
	for _, s := range sets {
		s.items, s.chain = s.items[:0], s.chain[:0]
		clear(s.seen)
		clear(s.head)
	}
	// A rule predicted in set k of this call holds base+k+1 in predicted.
	base := r.stamp
	r.stamp += int64(n + 1)
	matchCache := make([][]termMatch, n+1)
	matched := make([]bool, n+1)
	matches := func(k int) []termMatch {
		if !matched[k] {
			matchCache[k] = r.matchesAt(in, k)
			matched[k] = true
		}
		return matchCache[k]
	}
	canStart := func(b bitset, k int) bool {
		for _, m := range matches(k) {
			if b.has(m.term) {
				return true
			}
		}
		return false
	}
	add := func(k int, it item) {
		s := sets[k]
		if _, ok := s.seen[it]; ok {
			return
		}
		s.seen[it] = struct{}{}
		s.items = append(s.items, it)
		alt := r.rules[it.rule].alts[it.alt]
		if int(it.dot) < len(alt) && alt[it.dot].rule >= 0 {
			nt := int32(alt[it.dot].rule)
			prev, ok := s.head[nt]
			if !ok {
				prev = -1
			}
			s.head[nt] = int32(len(s.chain))
			s.chain = append(s.chain, waiter{it: it, next: prev})
		}
	}
	predict := func(k, nt int) {
		stamp := base + int64(k) + 1
		if r.predicted[nt] == stamp {
			return
		}
		r.predicted[nt] = stamp
		// An alternative that cannot start with the next token is not
		// predicted: one that derives nothing here is passed over where
		// it is used, which is what a nullable rule's prediction is for.
		for a := range r.rules[nt].alts {
			if canStart(r.altFst[nt][a], k) {
				add(k, item{int32(nt), int32(a), 0, int32(k)})
			}
		}
	}
	predict(0, r.start)
	if n == 0 {
		return r.null[r.start]
	}
	for k := 0; k <= n; k++ {
		s := sets[k]
		for q := 0; q < len(s.items); q++ {
			it := s.items[q]
			alt := r.rules[it.rule].alts[it.alt]
			if int(it.dot) == len(alt) {
				o := sets[it.origin]
				if i, ok := o.head[it.rule]; ok {
					for ; i >= 0; i = o.chain[i].next {
						w := o.chain[i].it
						add(k, item{w.rule, w.alt, w.dot + 1, w.origin})
					}
				}
				continue
			}
			sym := alt[it.dot]
			next := item{it.rule, it.alt, it.dot + 1, it.origin}
			if sym.rule >= 0 {
				if canStart(r.first[sym.rule], k) {
					predict(k, sym.rule)
				}
				if r.null[sym.rule] {
					add(k, next)
				}
				continue
			}
			t := r.terms[sym.term]
			if len(t.toks) == 0 {
				add(k, next)
				continue
			}
			for _, m := range matches(k) {
				if m.term == sym.term {
					add(k+m.width, next)
				}
			}
		}
	}
	for _, it := range sets[n].items {
		if int(it.rule) == r.start && it.origin == 0 && int(it.dot) == len(r.rules[it.rule].alts[it.alt]) {
			return true
		}
	}
	return false
}

type termMatch struct{ term, width int }

// matchesAt returns every terminal that matches the input at k, and how
// many tokens each takes.
func (r *Recognizer) matchesAt(in []ltok, k int) []termMatch {
	if k >= len(in) {
		return nil
	}
	var out []termMatch
	tk := in[k]
	for _, sp := range r.exact[tk.key()] {
		toks := r.terms[sp.term].spelling(sp.alt)
		if k+len(toks) > len(in) {
			continue
		}
		ok := true
		for j, want := range toks {
			if in[k+j].key() != want.key() {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, termMatch{sp.term, len(toks)})
		}
	}
	if tk.kind == kName && !tk.quoted && r.g.reserved[strings.ToUpper(tk.text)] {
		return out
	}
	for _, id := range r.class[classOf(tk)] {
		out = append(out, termMatch{id, 1})
	}
	if tk.kind == kName && tk.quoted {
		for _, id := range r.class[nameClass] {
			out = append(out, termMatch{id, 1})
		}
	}
	return out
}

// --- a generic SQL lexer ---

type lkind int

const (
	kNone   lkind = iota
	kName         // a word, or a "quoted" or `quoted` identifier
	kString       // '...', $$...$$, or a prefixed string such as N'...'
	kNumber       // 1, 1.5, 1e3, .5, 0x1F
	kVar          // @v, @@v, $1, $v
	kPunct        // an operator or punctuation mark
	kBad          // a quote never closed; it matches no terminal
)

type ltok struct {
	kind   lkind
	text   string
	quoted bool
}

// classOf is the class a token stands for when a sample: its kind, for a
// variable its sigil, so a stage @a and a column $1 are told apart, and for
// a quoted name its quote, so PartiQL's "a", a name, and `a`, an Ion value,
// are.
func classOf(t ltok) string {
	c := string(rune('0' + t.kind))
	switch {
	case t.kind == kVar:
		c += t.text[:len(t.text)-len(strings.TrimLeft(t.text, "@$"))]
	case t.kind == kName && t.quoted:
		c += t.text[:1]
	}
	return c
}

// nameClass is the class of a token sampled by an unquoted name. It takes
// a quoted name too, in any quotes: a grammar from docs samples one name
// for both.
var nameClass = classOf(ltok{kind: kName})

// key is what an exact match compares: a word without regard to case,
// anything else as written.
func (t ltok) key() string {
	switch {
	case t.kind == kName && !t.quoted:
		return "w" + strings.ToUpper(t.text)
	case t.kind == kName:
		return "q" + t.text
	}
	return string(rune('0'+t.kind)) + t.text
}

// hasRunePrefix reports whether rs starts with p.
func hasRunePrefix(rs, p []rune) bool {
	if len(rs) < len(p) {
		return false
	}
	for i, c := range p {
		if rs[i] != c {
			return false
		}
	}
	return true
}

// unquoted returns a sample without what its quotes hold, so the
// punctuation inside a string sample, '()', does not become a token of the
// grammar's.
func unquoted(s string) string {
	var b strings.Builder
	var quote rune
	for _, c := range s {
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"' || c == '`':
			quote = c
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
}

// punctPieces returns the runs of punctuation in a sample, which the lexer
// must read whole.
func punctPieces(s string) []string {
	var out []string
	var cur strings.Builder
	flush := func() {
		if cur.Len() > 1 {
			out = append(out, cur.String())
		}
		cur.Reset()
	}
	for _, c := range s {
		if isPunct(c) {
			cur.WriteRune(c)
		} else {
			flush()
		}
	}
	flush()
	return out
}

func isPunct(c rune) bool {
	return c < utf8.RuneSelf && strings.ContainsRune("!%&*+-/<=>?^|~:;,.()[]{}", c)
}

func isWordStart(c rune) bool { return c == '_' || unicode.IsLetter(c) }
func isWordPart(c rune) bool {
	return c == '_' || c == '$' || unicode.IsLetter(c) || unicode.IsDigit(c)
}

// lex reads text into tokens. It never fails: a character it has no rule
// for is a token of its own.
func (r *Recognizer) lex(s string) []ltok {
	var out []ltok
	rs := []rune(s)
	n := len(rs)
	for i := 0; i < n; {
		c := rs[i]
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '-' && i+1 < n && rs[i+1] == '-':
			for i < n && rs[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < n && rs[i+1] == '*':
			j := skipBlockComment(rs, i)
			if j > n {
				// A comment never closed: what follows is not a statement.
				out = append(out, ltok{kind: kBad, text: string(rs[i:])})
			}
			i = j
		case c == '\'':
			j, closed := quoted(rs, i, '\'')
			out = append(out, ltok{kind: closedKind(closed, kString), text: string(rs[i:j])})
			i = j
		case c == '$' && i+1 < n && rs[i+1] == '$':
			j := i + 2
			for j+1 < n && !(rs[j] == '$' && rs[j+1] == '$') {
				j++
			}
			closed := j+1 < n
			j = min(j+2, n)
			out = append(out, ltok{kind: closedKind(closed, kString), text: string(rs[i:j])})
			i = j
		case c == '"' || c == '`':
			j, closed := quoted(rs, i, c)
			out = append(out, ltok{kind: closedKind(closed, kName), text: string(rs[i:j]), quoted: true})
			i = j
		case unicode.IsDigit(c) || (c == '.' && i+1 < n && unicode.IsDigit(rs[i+1])):
			j, ok := number(rs, i)
			out = append(out, ltok{kind: closedKind(ok, kNumber), text: string(rs[i:j])})
			i = j
		case (c == '@' || c == '$') && i+1 < n && (isWordPart(rs[i+1]) || rs[i+1] == '@'):
			j := i + 1
			for j < n && rs[j] == '@' {
				j++
			}
			for j < n && isWordPart(rs[j]) {
				j++
			}
			out = append(out, ltok{kind: kVar, text: string(rs[i:j])})
			i = j
		case isWordStart(c):
			j := i
			for j < n && isWordPart(rs[j]) {
				j++
			}
			if j-i == 1 && j < n && rs[j] == '\'' {
				// A prefixed string: N'...', X'...', B'...', E'...'.
				k, closed := quoted(rs, j, '\'')
				out = append(out, ltok{kind: closedKind(closed, kString), text: string(rs[i:k])})
				i = k
				continue
			}
			out = append(out, ltok{kind: kName, text: string(rs[i:j])})
			i = j
		case isPunct(c):
			width := 1
			for _, p := range r.punct {
				if hasRunePrefix(rs[i:], p) {
					width = len(p)
					break
				}
			}
			out = append(out, ltok{kind: kPunct, text: string(rs[i : i+width])})
			i += width
		default:
			out = append(out, ltok{kind: kPunct, text: string(c)})
			i++
		}
	}
	return out
}

// quoted returns the index past a quoted run opened at i, where a doubled
// quote or a backslash escapes one, and whether the run was closed.
func quoted(rs []rune, i int, q rune) (int, bool) {
	j := i + 1
	for j < len(rs) {
		switch {
		case rs[j] == '\\' && q == '\'':
			j += 2
		case rs[j] == q && j+1 < len(rs) && rs[j+1] == q:
			j += 2
		case rs[j] == q:
			return j + 1, true
		default:
			j++
		}
	}
	return len(rs), false
}

// closedKind is kind for a quoted run that was closed, and kBad for one
// that ran to the end of the text.
func closedKind(closed bool, kind lkind) lkind {
	if closed {
		return kind
	}
	return kBad
}

// number returns the index past a number at i, and false for a hex prefix
// with no digit after it. A decimal point ends a number, 2., unless a dot
// or a name follows, as in a range 1..2 or a qualified name.
func number(rs []rune, i int) (int, bool) {
	n := len(rs)
	j := i
	if rs[j] == '0' && j+1 < n && (rs[j+1] == 'x' || rs[j+1] == 'X') {
		j += 2
		start := j
		for j < n && strings.ContainsRune("0123456789abcdefABCDEF", rs[j]) {
			j++
		}
		return j, j > start
	}
	for j < n && unicode.IsDigit(rs[j]) {
		j++
	}
	if j < n && rs[j] == '.' && (j+1 >= n || rs[j+1] != '.' && !isWordStart(rs[j+1])) {
		j++
		for j < n && unicode.IsDigit(rs[j]) {
			j++
		}
	}
	if j < n && (rs[j] == 'e' || rs[j] == 'E') {
		k := j + 1
		if k < n && (rs[k] == '+' || rs[k] == '-') {
			k++
		}
		if k < n && unicode.IsDigit(rs[k]) {
			for k < n && unicode.IsDigit(rs[k]) {
				k++
			}
			j = k
		}
	}
	return j, true
}
