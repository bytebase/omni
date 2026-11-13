package grammar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ReadANTLR reads an ANTLR 4 grammar: a combined grammar, or a lexer and a
// parser grammar given as one text. Parser rules (a lowercase name) become
// the rules rows are generated from, the first of them the start symbol.
// Lexer rules (an uppercase name) are tokens: one that spells a single text,
// such as SELECT: 'SELECT', or [sS] [eE] ... through fragments, is sampled
// by that text, and tokens supplies the rest, such as identifiers and
// numbers. A token table form overrides a lexer rule's own text.
//
// EBNF is lowered to plain alternatives the generator reads: a group or an
// element under ?, * or + becomes a rule of its own, named after the rule
// that holds it and the group's text, as in alterWarehouse(IF_EXISTS)?, so
// a row's label stays stable while the group does. Actions, semantic
// predicates, element labels and alternative labels carry no symbol and are
// dropped; a predicate's restriction is lost with it. A parser rule may not
// use ~ or the . wildcard, whose shortest text is not defined.
func ReadANTLR(src string, tokens *Tokens) (*Grammar, error) {
	// A file saved with a byte order mark, as Cosmos DB's sql.g4 is.
	src = strings.TrimPrefix(src, "\ufeff")
	toks, err := lexANTLR(src)
	if err != nil {
		return nil, err
	}
	g := newGrammar(tokens)
	g.antlr = true
	p := &antlrParser{toks: toks, g: g, lexer: map[string][]lexAlt{}, fragments: map[string]bool{}}
	if err := p.file(); err != nil {
		return nil, err
	}
	if len(g.rules) == 0 {
		return nil, fmt.Errorf("no parser rules found")
	}
	if err := g.sampleRules(tokens); err != nil {
		return nil, err
	}
	if _, ok := g.termText["EOF"]; !ok && !g.unreachable["EOF"] {
		g.termText["EOF"] = ""
	}
	// The tokens that match a class of texts, such as an identifier.
	// A fragment is no token, and a token the token table marks
	// unreachable is never produced: neither takes a word from another.
	inactive := func(name string) bool { return p.fragments[name] || g.unreachable[name] }
	var classes []int
	for i, name := range p.lexerOrder {
		if inactive(name) {
			continue
		}
		if _, finite := p.lexerTexts(name, map[string]bool{}); !finite {
			classes = append(classes, i)
		}
	}
	// takenBefore reports whether a class token written before position i
	// matches the word: ANTLR gives a text two lexer rules match to the one
	// written first.
	takenBefore := func(i int, word string) bool {
		for _, c := range classes {
			if c < i && p.lexerMatches(p.lexerOrder[c], word) {
				return true
			}
		}
		return false
	}
	for i, name := range p.lexerOrder {
		if inactive(name) {
			continue
		}
		text, _, sampled := p.lexerSample(name, map[string]bool{})
		texts, finite := p.lexerTexts(name, map[string]bool{})
		if finite && len(classes) > 0 {
			// A word a token spells, which no class token written before it
			// matches, is that token and never a name. The parser rules say
			// where it may stand for one, as Cosmos DB's identifier :
			// LEX_IDENTIFIER | K_ALL does.
			// The recognizer reserves a word in every case, so it does only
			// where no case the token is spelled in goes to an earlier class.
			for _, t := range texts {
				if !wordText.MatchString(t) {
					continue
				}
				taken := false
				for _, v := range p.caseVariants(name, t) {
					if takenBefore(i, v) {
						taken = true
						break
					}
				}
				if !taken {
					g.reserved[strings.ToUpper(t)] = true
				}
			}
		}
		// One text, in whatever case: a keyword. A token of a few texts is
		// sampled by its first, in capitals where it is a word.
		fixed := finite && len(texts) == 1
		if finite && len(texts) > 0 {
			text, sampled = texts[0], true
		}
		if _, ok := g.termText[name]; ok && finite && len(texts) > 1 {
			// The token table picks the spelling rows write, as TABLE for
			// Cassandra's COLUMNFAMILY | TABLE; the token still matches each.
			g.spellings[name] = texts
		}
		if _, ok := g.termText[name]; ok || g.unreachable[name] || !sampled {
			continue
		}
		g.termText[name] = text
		switch {
		case fixed:
			g.fixed[name] = true
		case finite:
			// DIRECTION : 'ASC' | 'DESC' : those words, and no other.
			g.spellings[name] = texts
		default:
			g.patterned[name] = true
		}
	}
	if len(classes) > 0 {
		// A literal in a parser rule is the lexer rule that spells it alone,
		// LATE : 'late', where there is one, with that rule's precedence.
		// Otherwise it is a token of its own, defined before every lexer
		// rule, so a word written as one is reserved.
		spelledBy := map[string]int{}
		for i, name := range p.lexerOrder {
			if inactive(name) {
				continue
			}
			if text, ok := p.literalRule(name); ok {
				if _, seen := spelledBy[text]; !seen {
					spelledBy[text] = i
				}
			}
		}
		for _, r := range g.rules {
			for _, alt := range r.alts {
				for _, sym := range alt {
					if !sym.literal || !wordText.MatchString(sym.text) {
						continue
					}
					if i, ok := spelledBy[sym.text]; ok && takenBefore(i, sym.text) {
						continue
					}
					g.reserved[strings.ToUpper(sym.text)] = true
				}
			}
		}
	}
	g.start = g.rules[0].name
	if err := g.resolveSymbols(); err != nil {
		return nil, err
	}
	return g, nil
}

// wordText matches a text an identifier rule could also match: a keyword.
var wordText = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_$]*$`)

// wordPart matches a piece of a word, such as one letter of a keyword.
var wordPart = regexp.MustCompile(`^[A-Za-z0-9_$]+$`)

// caseInsensitiveOption matches an options block that sets caseInsensitive.
var caseInsensitiveOption = regexp.MustCompile(`\bcaseInsensitive\s*=\s*true\b`)

// --- lexing the grammar file ---

type akind int

const (
	aIdent  akind = iota
	aString       // a '...' literal; text is its content
	aPunct        // an operator: : ; | ( ) ? * + ~ . .. = += # -> , < > @ ::
	aAction       // a { ... } action or option block
	aArg          // a [ ... ] argument or character set; text is its content
)

type atok struct {
	kind akind
	text string
	pos  int // byte offset, for errors
}

func lexANTLR(src string) ([]atok, error) {
	var out []atok
	r := []rune(src)
	// The byte offset of each rune, for errors.
	offsets := make([]int, len(r)+1)
	for i, c := range r {
		offsets[i+1] = offsets[i] + utf8.RuneLen(c)
	}
	byteAt := func(i int) int { return offsets[min(i, len(r))] }
	n := len(r)
	for i := 0; i < n; {
		c := r[i]
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '/' && i+1 < n && r[i+1] == '*':
			i = skipBlockComment(r, i)
		case c == '/' && i+1 < n && r[i+1] == '/':
			i = skipToEOL(r, i)
		case c == '{':
			start := i
			i = skipAction(r, i)
			out = append(out, atok{kind: aAction, text: string(r[start:i]), pos: byteAt(start)})
		case c == '[':
			start := i
			i++
			for i < n && r[i] != ']' {
				if r[i] == '\\' {
					i++
				}
				i++
			}
			if i >= n {
				return nil, fmt.Errorf("byte %d: unterminated [", byteAt(start))
			}
			out = append(out, atok{kind: aArg, text: string(r[start+1 : i]), pos: byteAt(start)})
			i++
		case c == '\'':
			start := i
			text, next, err := readANTLRString(r, i)
			if err != nil {
				return nil, fmt.Errorf("byte %d: %v", byteAt(start), err)
			}
			out = append(out, atok{kind: aString, text: text, pos: byteAt(start)})
			i = next
		case isIdentStart(c):
			start := i
			word, next := readWord(r, i)
			out = append(out, atok{kind: aIdent, text: word, pos: byteAt(start)})
			i = next
		default:
			start := i
			op := string(c)
			if i+1 < n {
				switch two := string(r[i : i+2]); two {
				case "..", "+=", "->", "::":
					op = two
				}
			}
			if len(op) == 1 && !strings.Contains(":;|()?*+~.=#,<>@", op) {
				return nil, fmt.Errorf("byte %d: unexpected %q", byteAt(start), c)
			}
			out = append(out, atok{kind: aPunct, text: op, pos: byteAt(start)})
			i += utf8.RuneCountInString(op)
		}
	}
	return out, nil
}

// readANTLRString reads a '...' literal with ANTLR's escapes.
func readANTLRString(r []rune, i int) (string, int, error) {
	i++
	var b strings.Builder
	for i < len(r) {
		c := r[i]
		switch {
		case c == '\'':
			return b.String(), i + 1, nil
		case c == '\\' && i+1 < len(r):
			i++
			switch e := r[i]; e {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case 'u':
				if i+4 < len(r) {
					if v, err := strconv.ParseUint(string(r[i+1:i+5]), 16, 32); err == nil {
						b.WriteRune(rune(v))
						i += 4
						break
					}
				}
				return "", i, fmt.Errorf("bad \\u escape")
			default:
				b.WriteRune(e)
			}
			i++
		case c == '\n':
			return "", i, fmt.Errorf("newline in a literal")
		default:
			b.WriteRune(c)
			i++
		}
	}
	return "", i, fmt.Errorf("unterminated literal")
}

// --- parsing the grammar file ---

// An element of an EBNF alternative: a symbol or a parenthesised group,
// with the operator that follows it (0, '?', '*' or '+').
type element struct {
	sym   *symbol
	group [][]element
	op    byte
}

// A lexAlt is one alternative of a lexer rule, kept only to find the one
// text the rule spells, if it spells one.
type lexAlt []lexElem

type lexElem struct {
	literal string    // a '...' literal
	ref     string    // a lexer rule
	set     string    // a [...] set: its first character
	oneChar bool      // the set holds one character, or one letter in both cases
	group   []lexAlt  // a parenthesised group
	op      byte      // 0, '?', '*' or '+'
	opaque  bool      // ~x or the . wildcard: no single text
	ranged  [2]string // 'a'..'z'
	class   string    // a [...] set as written, for matching
	any     bool      // the . wildcard
	neg     *lexElem  // ~x: one character x does not match
}

type antlrParser struct {
	toks  []atok
	i     int
	g     *Grammar
	lexer map[string][]lexAlt
	// lexerOrder lists the lexer rules as the grammar writes them: where two
	// match a text, the earlier one takes it.
	lexerOrder []string
	// caseInsensitive is the grammar's option of that name: its lexer rules
	// match letters in either case.
	caseInsensitive bool
	// fragments holds the lexer rules that are parts of tokens, not tokens.
	fragments map[string]bool
	// helpers maps, per rule, each lowered group's name to its structure:
	// group written twice in one rule shares its rule.
	helpers map[string]string
}

func (p *antlrParser) peek() (atok, bool) {
	if p.i < len(p.toks) {
		return p.toks[p.i], true
	}
	return atok{}, false
}

func (p *antlrParser) is(kind akind, text string) bool {
	t, ok := p.peek()
	return ok && t.kind == kind && t.text == text
}

func (p *antlrParser) errorf(format string, args ...any) error {
	pos := -1
	if t, ok := p.peek(); ok {
		pos = t.pos
	}
	return fmt.Errorf("byte %d: %s", pos, fmt.Sprintf(format, args...))
}

func (p *antlrParser) expect(kind akind, text string) error {
	if !p.is(kind, text) {
		return p.errorf("want %q", text)
	}
	p.i++
	return nil
}

func (p *antlrParser) file() error {
	for {
		t, ok := p.peek()
		if !ok {
			return nil
		}
		switch {
		case t.kind == aIdent && (t.text == "grammar" || t.text == "lexer" || t.text == "parser" || t.text == "mode"):
			// A header (lexer grammar X;) or a lexer mode (mode X;): no rule.
			if t.text == "mode" || (p.i+1 < len(p.toks) && p.toks[p.i+1].kind == aIdent) {
				if err := p.skipTo(";"); err != nil {
					return err
				}
				continue
			}
			if err := p.rule(); err != nil {
				return err
			}
		case t.kind == aIdent && (t.text == "options" || t.text == "tokens" || t.text == "channels") &&
			p.i+1 < len(p.toks) && p.toks[p.i+1].kind == aAction:
			if t.text == "options" && caseInsensitiveOption.MatchString(p.toks[p.i+1].text) {
				p.caseInsensitive = true
			}
			p.i += 2
		case t.kind == aIdent && t.text == "import":
			return p.errorf("import is not supported: give the imported grammar's text along with this one")
		case t.kind == aPunct && t.text == "@":
			// A named action: @header {...}, @parser::members {...}.
			p.i++
			for p.i < len(p.toks) && p.toks[p.i].kind != aAction {
				p.i++
			}
			p.i++
		case t.kind == aIdent:
			if err := p.rule(); err != nil {
				return err
			}
		default:
			return p.errorf("expected a rule")
		}
	}
}

func (p *antlrParser) skipTo(text string) error {
	for p.i < len(p.toks) {
		t := p.toks[p.i]
		p.i++
		if t.kind == aPunct && t.text == text {
			return nil
		}
	}
	return fmt.Errorf("unexpected end of grammar, want %q", text)
}

func (p *antlrParser) rule() error {
	fragment := p.is(aIdent, "fragment")
	if fragment {
		p.i++
	}
	t, ok := p.peek()
	if !ok || t.kind != aIdent {
		return p.errorf("want a rule name")
	}
	name := t.text
	if fragment {
		p.fragments[name] = true
	}
	p.i++
	// Everything between the name and the colon describes the generated
	// parser: arguments, returns, locals, throws, options, @init.
	for !p.is(aPunct, ":") {
		if _, ok := p.peek(); !ok {
			return fmt.Errorf("rule %s: no colon", name)
		}
		p.i++
	}
	p.i++
	if unicode.IsUpper([]rune(name)[0]) {
		alts, err := p.lexAlts()
		if err != nil {
			return fmt.Errorf("lexer rule %s: %v", name, err)
		}
		if _, seen := p.lexer[name]; !seen {
			p.lexerOrder = append(p.lexerOrder, name)
		}
		p.lexer[name] = append(p.lexer[name], alts...)
	} else {
		alts, err := p.alts()
		if err != nil {
			return fmt.Errorf("rule %s: %v", name, err)
		}
		r := &rule{name: name}
		if existing, ok := p.g.byName[name]; ok {
			r = existing
		} else {
			p.g.rules = append(p.g.rules, r)
			p.g.byName[name] = r
		}
		p.helpers = map[string]string{}
		for _, alt := range alts {
			r.alts = append(r.alts, p.lower(name, alt))
		}
	}
	if err := p.expect(aPunct, ";"); err != nil {
		return fmt.Errorf("rule %s: %v", name, err)
	}
	// Exception handlers follow a rule: catch [...] {...}, finally {...}.
	for p.is(aIdent, "catch") || p.is(aIdent, "finally") {
		for p.i < len(p.toks) && p.toks[p.i].kind != aAction {
			p.i++
		}
		p.i++
	}
	return nil
}

// alts reads alternatives up to a ')' or ';', which it leaves unread.
func (p *antlrParser) alts() ([][]element, error) {
	var out [][]element
	for {
		alt, err := p.alt()
		if err != nil {
			return nil, err
		}
		out = append(out, alt)
		if !p.is(aPunct, "|") {
			return out, nil
		}
		p.i++
	}
}

func (p *antlrParser) alt() ([]element, error) {
	var out []element
	for {
		t, ok := p.peek()
		if !ok {
			return nil, fmt.Errorf("unexpected end of grammar")
		}
		switch {
		case t.kind == aPunct && (t.text == "|" || t.text == ")" || t.text == ";"):
			return out, nil
		case t.kind == aPunct && t.text == "#":
			// An alternative label names the generated context class.
			p.i += 2
			continue
		case t.kind == aPunct && t.text == "<":
			// Element options: <assoc=right>.
			if err := p.skipTo(">"); err != nil {
				return nil, err
			}
			continue
		case t.kind == aAction:
			// An action, or a semantic predicate when ? follows.
			p.i++
			if p.is(aPunct, "?") {
				p.i++
			}
			continue
		case t.kind == aPunct && t.text == "->":
			return nil, p.errorf("a lexer command in a parser rule")
		}
		e, err := p.element()
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
}

func (p *antlrParser) element() (element, error) {
	// A label: name=atom or name+=atom.
	if t, ok := p.peek(); ok && t.kind == aIdent && p.i+1 < len(p.toks) {
		if next := p.toks[p.i+1]; next.kind == aPunct && (next.text == "=" || next.text == "+=") {
			p.i += 2
		}
	}
	t, _ := p.peek()
	var e element
	switch {
	case t.kind == aIdent:
		e.sym = &symbol{text: t.text}
		p.i++
	case t.kind == aString:
		e.sym = &symbol{text: t.text, literal: true}
		if l := labelText(t.text); l != t.text {
			e.sym.label = l
		}
		p.i++
	case t.kind == aPunct && t.text == "(":
		p.i++
		alts, err := p.alts()
		if err != nil {
			return e, err
		}
		if err := p.expect(aPunct, ")"); err != nil {
			return e, err
		}
		e.group = alts
	case t.kind == aPunct && (t.text == "~" || t.text == "."):
		return e, p.errorf("%s in a parser rule: no shortest text to write for it", t.text)
	default:
		return e, p.errorf("unexpected %q", t.text)
	}
	if t, ok := p.peek(); ok && t.kind == aPunct && (t.text == "?" || t.text == "*" || t.text == "+") {
		e.op = t.text[0]
		p.i++
		if p.is(aPunct, "?") { // non-greedy
			p.i++
		}
	}
	return e, nil
}

// --- lowering EBNF to the generator's alternatives ---

// maxGroupName is the longest a lowered group's text may be in its rule's
// name, and maxAltLabel the longest an alternative's symbols may be in its
// label; a longer one is cut and ends in a hash of the whole.
const (
	maxGroupName = 48
	maxAltLabel  = 64
)

// lower turns one EBNF alternative of rule into plain symbols, adding a
// rule for each group or operator it holds.
func (p *antlrParser) lower(owner string, alt []element) []symbol {
	var out []symbol
	for _, e := range alt {
		if e.op == 0 && e.group != nil && len(e.group) == 1 {
			// (a b) with no operator is just its symbols.
			out = append(out, p.lower(owner, e.group[0])...)
			continue
		}
		if e.op == 0 && e.group == nil {
			out = append(out, *e.sym)
			continue
		}
		out = append(out, p.helper(owner, e))
	}
	return out
}

// helper returns the symbol of the rule a group or an operator is lowered
// to, adding the rule the first time owner names it:
//
//	x?  h: | x        x*  h: | x h        x+  h: x | x h
//
// where x is the group's one alternative, or a rule for its choice.
func (p *antlrParser) helper(owner string, e element) symbol {
	body := e.group
	if body == nil {
		body = [][]element{{{sym: e.sym}}}
	}
	text := "(" + groupText(body) + ")"
	if e.op != 0 {
		text += string(e.op)
	}
	if len(text) > maxGroupName {
		// Cut inside the parentheses, so the operator still ends the name.
		inner := strings.TrimSuffix(text, string(e.op))
		inner = capLabel(strings.TrimSuffix(inner, ")"), maxGroupName-2) + ")"
		if e.op != 0 {
			inner += string(e.op)
		}
		text = inner
	}
	// The same group written twice in one rule shares its rule. Two groups
	// whose labels read alike, ('A_B')? and ('A' 'B')?, do not: the later
	// one is numbered.
	structure := string(e.op) + groupStructure(body)
	name := owner + text
	for n := 2; ; n++ {
		seen, ok := p.helpers[name]
		if !ok || seen == structure {
			break
		}
		name = fmt.Sprintf("%s%s#%d", owner, text, n)
	}
	sym := symbol{text: name, nonterm: true, label: strings.TrimPrefix(name, owner)}
	if _, ok := p.helpers[name]; ok {
		return sym
	}
	p.helpers[name] = structure
	h := &rule{name: name}
	// Added before its alternatives are lowered, so a group nested in it is
	// listed after it, as the grammar writes them.
	p.g.rules = append(p.g.rules, h)
	p.g.byName[name] = h
	var alts [][]symbol
	if (e.op == '*' || e.op == '+') && len(body) > 1 {
		// (a | b)+ repeats a choice, which is a rule of its own: each
		// choice has a row, and repeating takes two rows, not two a choice.
		alts = [][]symbol{{p.helper(owner, element{group: body})}}
	} else {
		for _, a := range body {
			alts = append(alts, p.lower(owner, a))
		}
	}
	switch e.op {
	case 0:
		h.alts = alts
	case '?':
		h.alts = append([][]symbol{{}}, alts...)
	case '*':
		h.alts = [][]symbol{{}}
		for _, a := range alts {
			h.alts = append(h.alts, append(append([]symbol(nil), a...), sym))
		}
	case '+':
		h.alts = append(h.alts, alts...)
		for _, a := range alts {
			h.alts = append(h.alts, append(append([]symbol(nil), a...), sym))
		}
	}
	return sym
}

// groupText writes a group as a label does: symbols joined by _, its
// alternatives by |, a nested group in parentheses with its operator. A
// label holds no space, and no colon or slash, which mark a hand-added row
// and a rule's end.
func groupText(alts [][]element) string {
	parts := make([]string, len(alts))
	for i, alt := range alts {
		var syms []string
		for _, e := range alt {
			s := ""
			if e.group != nil {
				s = "(" + groupText(e.group) + ")"
			} else {
				s = labelText(e.sym.text)
			}
			if e.op != 0 {
				s += string(e.op)
			}
			syms = append(syms, s)
		}
		parts[i] = strings.Join(syms, "_")
	}
	return strings.Join(parts, "|")
}

// groupStructure writes a group so that two groups differ in it exactly
// when they differ: a literal quoted, separators no symbol holds.
func groupStructure(alts [][]element) string {
	var b strings.Builder
	for i, alt := range alts {
		if i > 0 {
			b.WriteString("\x00|")
		}
		for _, e := range alt {
			switch {
			case e.group != nil:
				b.WriteString("\x00(" + groupStructure(e.group) + "\x00)")
			case e.sym.literal:
				b.WriteString("\x00'" + e.sym.text)
			default:
				b.WriteString("\x00 " + e.sym.text)
			}
			if e.op != 0 {
				b.WriteByte(e.op)
			}
		}
	}
	return b.String()
}

func labelText(s string) string {
	return strings.NewReplacer(" ", "_", "\t", "_", "\n", "_", ":", "COLON", "/", "SLASH").Replace(s)
}

// --- lexer rules: the one text a token spells ---

func (p *antlrParser) lexAlts() ([]lexAlt, error) {
	var out []lexAlt
	for {
		var alt lexAlt
	elements:
		for {
			t, ok := p.peek()
			if !ok {
				return nil, fmt.Errorf("unexpected end of grammar")
			}
			var e lexElem
			switch {
			case t.kind == aPunct && (t.text == "|" || t.text == ")" || t.text == ";"):
				break elements
			case t.kind == aPunct && t.text == "->":
				// Lexer commands: -> skip, -> channel(HIDDEN), -> type(X).
				for p.i < len(p.toks) && !p.is(aPunct, ";") && !p.is(aPunct, "|") && !p.is(aPunct, ")") {
					if p.is(aPunct, "(") {
						if err := p.skipTo(")"); err != nil {
							return nil, err
						}
						continue
					}
					p.i++
				}
				continue
			case t.kind == aAction:
				p.i++
				if p.is(aPunct, "?") {
					p.i++
				}
				continue
			case t.kind == aString:
				p.i++
				e.literal = t.text
				if p.is(aPunct, "..") {
					p.i++
					hi, _ := p.peek()
					p.i++
					e = lexElem{ranged: [2]string{t.text, hi.text}}
				}
			case t.kind == aIdent:
				p.i++
				e.ref = t.text
			case t.kind == aArg:
				p.i++
				e.set = firstOfSet(t.text)
				e.oneChar = oneCharSet(t.text)
				e.opaque = e.set == ""
				e.class = t.text
			case t.kind == aPunct && t.text == "~":
				p.i++
				operand, err := p.lexAtom()
				if err != nil {
					return nil, err
				}
				e.opaque, e.neg = true, &operand
			case t.kind == aPunct && t.text == ".":
				p.i++
				e.opaque, e.any = true, true
			case t.kind == aPunct && t.text == "(":
				p.i++
				g, err := p.lexAlts()
				if err != nil {
					return nil, err
				}
				if err := p.expect(aPunct, ")"); err != nil {
					return nil, err
				}
				e.group = g
			default:
				return nil, p.errorf("unexpected %q", t.text)
			}
			if t, ok := p.peek(); ok && t.kind == aPunct && (t.text == "?" || t.text == "*" || t.text == "+") {
				e.op = t.text[0]
				p.i++
				if p.is(aPunct, "?") {
					p.i++
				}
			}
			alt = append(alt, e)
		}
		out = append(out, alt)
		if !p.is(aPunct, "|") {
			return out, nil
		}
		p.i++
	}
}

// caseVariants returns the spellings of folded, one of name's texts, that
// differ only in case: every one where they are few, and otherwise, for a
// keyword spelled letter by letter in either case, its lower, upper and
// capitalized forms.
func (p *antlrParser) caseVariants(name, folded string) []string {
	if texts, ok := p.ruleTexts(name, map[string]bool{}, false); ok {
		var out []string
		for _, t := range texts {
			if strings.EqualFold(t, folded) {
				out = append(out, t)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	lower := strings.ToLower(folded)
	return []string{folded, lower, strings.ToUpper(folded), strings.ToUpper(lower[:1]) + lower[1:]}
}

// literalRule returns the text of a lexer rule that is one literal and
// nothing else, X : 'x', which ANTLR takes as the token a parser rule's 'x'
// names.
func (p *antlrParser) literalRule(name string) (string, bool) {
	alts := p.lexer[name]
	if p.fragments[name] || len(alts) != 1 || len(alts[0]) != 1 {
		return "", false
	}
	e := alts[0][0]
	if e.op != 0 || e.opaque || e.ref != "" || e.group != nil || e.class != "" || e.ranged[0] != "" || e.literal == "" {
		return "", false
	}
	return e.literal, true
}

// lexAtom reads the operand of ~: a set, a literal, a rule, or a group.
func (p *antlrParser) lexAtom() (lexElem, error) {
	t, ok := p.peek()
	if !ok {
		return lexElem{}, fmt.Errorf("unexpected end of grammar")
	}
	p.i++
	switch {
	case t.kind == aPunct && t.text == "(":
		g, err := p.lexAlts()
		if err != nil {
			return lexElem{}, err
		}
		return lexElem{group: g}, p.expect(aPunct, ")")
	case t.kind == aArg:
		return lexElem{class: t.text}, nil
	case t.kind == aString:
		return lexElem{literal: t.text}, nil
	case t.kind == aIdent:
		return lexElem{ref: t.text}, nil
	}
	return lexElem{}, nil
}

// firstOfSet returns the first character a [...] set holds, or "" for an
// empty set. A letter in both cases, [sS], is the uppercase one: it is how
// a case-insensitive grammar spells a keyword, and a keyword is sampled in
// capitals.
func firstOfSet(set string) string {
	r := []rune(set)
	if len(r) == 0 {
		return ""
	}
	if len(r) == 2 && r[0] != r[1] && unicode.ToUpper(r[0]) == unicode.ToUpper(r[1]) {
		return string(unicode.ToUpper(r[0]))
	}
	if r[0] == '\\' && len(r) > 1 {
		switch r[1] {
		case 'n':
			return "\n"
		case 't':
			return "\t"
		case 'r':
			return "\r"
		case 'u':
			if len(r) >= 6 {
				if v, err := strconv.ParseUint(string(r[2:6]), 16, 32); err == nil {
					return string(rune(v))
				}
			}
			return ""
		case 'p', 'P':
			return ""
		}
		return string(r[1])
	}
	return string(r[0])
}

// lexerSample returns the shortest text a lexer rule matches when it can
// be told without running the lexer: each literal, set and range at its
// first character, an optional or repeated part left out, the first
// alternative that has a text. A rule that needs ~ or . has none.
//
// fixed reports whether the rule matches only that text, as a keyword
// does ('select', or S E L E C T through case-insensitive fragments), not a
// class of texts, as an identifier does.
func (p *antlrParser) lexerSample(name string, visiting map[string]bool) (text string, fixed, ok bool) {
	alts, ok := p.lexer[name]
	if !ok || visiting[name] {
		return "", false, false
	}
	visiting[name] = true
	defer delete(visiting, name)
	return p.altsSample(alts, visiting)
}

func (p *antlrParser) altsSample(alts []lexAlt, visiting map[string]bool) (string, bool, bool) {
	for _, alt := range alts {
		if s, fixed, ok := p.altSample(alt, visiting); ok {
			return s, fixed && len(alts) == 1, true
		}
	}
	return "", false, false
}

// maxSpellings bounds the texts lexerTexts lists for one token.
const maxSpellings = 64

// lexerTexts returns every text a lexer rule matches when they are few and
// fixed, as a rule choosing among keywords ('ASC' | 'DESC') matches, and
// false for a rule that matches a class of texts, such as identifiers. A
// word spelled in either case is one text in capitals (see foldCase).
func (p *antlrParser) lexerTexts(name string, visiting map[string]bool) ([]string, bool) {
	return p.ruleTexts(name, visiting, true)
}

func (p *antlrParser) ruleTexts(name string, visiting map[string]bool, fold bool) ([]string, bool) {
	alts, ok := p.lexer[name]
	if !ok || visiting[name] {
		return nil, false
	}
	visiting[name] = true
	defer delete(visiting, name)
	return p.altsTexts(alts, visiting, fold)
}

// altsTexts returns the texts of alts. With fold, each piece's spellings
// that differ only in case are folded as they are read, so a long keyword
// stays one text. The recognizer matches a word in any case but a quoted
// string or name exactly, so an alternative whose texts hold a quote is read
// again without folding, and keeps its case.
func (p *antlrParser) altsTexts(alts []lexAlt, visiting map[string]bool, fold bool) ([]string, bool) {
	var out []string
	for _, alt := range alts {
		texts, ok := p.altTexts(alt, visiting, fold)
		if ok && fold && holdsQuote(texts) {
			texts, ok = p.altTexts(alt, visiting, false)
		}
		if !ok {
			return nil, false
		}
		out = append(out, texts...)
	}
	if fold {
		out = foldCase(out)
	}
	return out, len(out) <= maxSpellings
}

func (p *antlrParser) altTexts(alt lexAlt, visiting map[string]bool, fold bool) ([]string, bool) {
	texts := []string{""}
	for _, e := range alt {
		var part []string
		switch {
		case e.op != 0, e.opaque, e.ranged[0] != "" && e.ranged[0] != e.ranged[1]:
			return nil, false
		case e.ranged[0] != "":
			part = []string{e.ranged[0]}
		case e.set != "":
			if !e.oneChar {
				return nil, false
			}
			// [sS] is two texts, which folding makes one where it can.
			part = setTexts(e.class)
		case e.ref != "":
			var ok bool
			if part, ok = p.ruleTexts(e.ref, visiting, fold); !ok {
				return nil, false
			}
		case e.group != nil:
			var ok bool
			if part, ok = p.altsTexts(e.group, visiting, fold); !ok {
				return nil, false
			}
		default:
			part = []string{e.literal}
		}
		if fold {
			part = foldCase(part)
		}
		var next []string
		for _, a := range texts {
			for _, b := range part {
				next = append(next, a+b)
			}
		}
		if len(next) > maxSpellings {
			return nil, false
		}
		texts = next
	}
	return texts, true
}

// setTexts returns the texts of a one-character [...] set: its character,
// or a letter in each of its cases.
func setTexts(class string) []string {
	r := []rune(class)
	if len(r) == 2 && r[0] != r[1] && unicode.ToUpper(r[0]) == unicode.ToUpper(r[1]) {
		return []string{string(r[0]), string(r[1])}
	}
	return []string{firstOfSet(class)}
}

// holdsQuote reports whether a text holds a quote, so that its case counts.
func holdsQuote(texts []string) bool {
	for _, t := range texts {
		if strings.ContainsAny(t, "'\"`") {
			return true
		}
	}
	return false
}

// foldCase returns texts that are one word in several cases as that word in
// capitals, and any other texts as they are. A grammar spells a keyword of
// any case letter by letter, each [sS] or ('s'|'S'), and a keyword is
// sampled in capitals; folding each letter as it is read keeps a long
// keyword from counting as 2^n texts. Only a word is folded, since the
// recognizer matches words without regard to case and anything else, such
// as a string, exactly.
func foldCase(texts []string) []string {
	if len(texts) < 2 || !wordPart.MatchString(texts[0]) {
		return texts
	}
	for _, t := range texts[1:] {
		if !strings.EqualFold(t, texts[0]) {
			return texts
		}
	}
	return []string{strings.ToUpper(texts[0])}
}

func (p *antlrParser) altSample(alt lexAlt, visiting map[string]bool) (string, bool, bool) {
	var b strings.Builder
	fixed := true
	for _, e := range alt {
		if e.op != 0 {
			fixed = false
		}
		if e.op == '?' || e.op == '*' {
			continue
		}
		var s string
		ok, one := true, true
		switch {
		case e.opaque:
			return "", false, false
		case e.ranged[0] != "":
			s, one = e.ranged[0], e.ranged[0] == e.ranged[1]
		case e.set != "":
			s, one = e.set, e.oneChar
		case e.ref != "":
			s, one, ok = p.lexerSample(e.ref, visiting)
		case e.group != nil:
			s, one, ok = p.altsSample(e.group, visiting)
		default:
			s = e.literal
		}
		if !ok {
			return "", false, false
		}
		fixed = fixed && one
		b.WriteString(s)
	}
	return b.String(), fixed, true
}

// oneCharSet reports whether a [...] set holds one character, or one
// letter in both cases.
func oneCharSet(set string) bool {
	r := []rune(set)
	switch {
	case len(r) == 1:
		return true
	case len(r) == 2 && r[0] == '\\':
		return true
	case len(r) == 2:
		return r[0] != r[1] && unicode.ToUpper(r[0]) == unicode.ToUpper(r[1])
	}
	return false
}
