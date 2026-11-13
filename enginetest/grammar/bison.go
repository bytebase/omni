package grammar

import (
	"fmt"
	"strings"
)

// A Grammar is a vendor grammar reduced to what row generation needs: the
// rules, the start symbol, and a sample of text for every terminal.
// Everything a parser generator adds for the parser it builds (prologue and
// epilogue code, %union, semantic actions, precedence declarations) is
// dropped, because none of it changes which strings the grammar derives.
type Grammar struct {
	// rules in source order.
	rules []*rule
	// byName indexes rules by their left-hand side.
	byName map[string]*rule
	start  string
	// termText gives the sample text each terminal lexes from, keyed by the
	// terminal's symbol text (a token name, or a quoted literal's content).
	termText map[string]string
	// unreachable holds the tokens the lexer never produces in the mode the
	// judge runs in. An alternative that needs one gets no generated row.
	unreachable map[string]bool
	// spaced holds the token texts written with a space on each side.
	spaced map[string]bool
	// ruleText gives the text written where a rule stands in another rule's
	// row, for the rules the token table samples (see Tokens).
	ruleText map[string]string
}

// A rule is one left-hand side and its alternatives, in source order.
type rule struct {
	name string
	alts [][]symbol
}

// A symbol is one element of an alternative. A nonterminal names another
// rule; a terminal is a token the lexer produces, identified by text: a
// token name (identifier, intLit, a keyword), or the content of a "..." or
// '.' literal.
type symbol struct {
	text    string
	nonterm bool
	// literal is set when the symbol came from a "..." or '.' grammar
	// literal, whose content is the text to emit, never a token name.
	literal bool
}

// ReadBison reads a bison or goyacc grammar. tokens supplies the sample of
// every terminal whose declaration does not spell it: literal classes such
// as identifiers and numbers, and tokens the lexer forms from other text.
// Every other terminal is sampled by its %token text.
//
// defined names the conditions the grammar is built with, for a file that
// holds more than one grammar between %ifdef lines (see Conditional). A
// file with no such line needs none.
func ReadBison(src string, tokens *Tokens, defined ...string) (*Grammar, error) {
	src, err := Conditional(src, defined...)
	if err != nil {
		return nil, err
	}
	decls, gram, err := splitSections(src)
	if err != nil {
		return nil, err
	}
	g := &Grammar{byName: map[string]*rule{}, termText: map[string]string{}, unreachable: map[string]bool{}, spaced: map[string]bool{}, ruleText: map[string]string{}}
	for name, form := range tokens.forms {
		switch form {
		case Spaced:
			g.spaced[name] = true
		case Unreachable:
			g.unreachable[name] = true
		case Empty:
			g.termText[name] = ""
		default:
			g.termText[name] = form
		}
	}
	g.readDeclarations(decls)

	if err := g.readRules(lexGrammar(gram)); err != nil {
		return nil, err
	}
	for _, r := range g.rules {
		form, ok := tokens.forms[r.name]
		if !ok {
			continue
		}
		switch form {
		case Spaced, Unreachable, Empty:
			return nil, fmt.Errorf("tokens.tsv gives the rule %s the form %s, which is a token's: a rule's row holds the text written for it", r.name, form)
		}
		g.ruleText[r.name] = form
		delete(g.termText, r.name)
	}
	if g.start == "" {
		g.start = g.rules[0].name
	}
	if _, ok := g.byName[g.start]; !ok {
		return nil, fmt.Errorf("start symbol %q has no rule", g.start)
	}
	if err := g.resolveSymbols(); err != nil {
		return nil, err
	}
	return g, nil
}

// splitSections returns the declaration section and the grammar section,
// the two halves either side of the first "%%". A second "%%" ends the
// grammar section: what follows it is code and is dropped. A grammar with
// no code after its rules, as MariaDB's sql_yacc.yy, has none.
func splitSections(src string) (decls, gram string, err error) {
	first := indexOfLine(src, "%%", 0)
	if first < 0 {
		return "", "", fmt.Errorf("no %%%% section marker found")
	}
	lines := strings.Split(src, "\n")
	second := indexOfLine(src, "%%", first+1)
	if second < 0 {
		second = len(lines)
	}
	return strings.Join(lines[:first], "\n"), strings.Join(lines[first+1:second], "\n"), nil
}

// indexOfLine returns the index, in lines, of the first line at or after
// start whose trimmed text equals want, or -1.
func indexOfLine(src, want string, start int) int {
	for i, line := range strings.Split(src, "\n") {
		if i >= start && strings.TrimSpace(line) == want {
			return i
		}
	}
	return -1
}

// readDeclarations records %start and the sample text of every terminal
// declared with one. A %token declaration names tokens, each optionally
// followed by a number and its text in quotes: the literal the lexer reads
// for a keyword or operator, or a description for a literal class, which
// the token table overrides. It may start on the %token line or continue
// on the lines after it.
func (g *Grammar) readDeclarations(decls string) {
	inToken := false
	for _, line := range strings.Split(decls, "\n") {
		trimmed := strings.TrimSpace(stripLineComment(line))
		switch {
		case strings.HasPrefix(trimmed, "%start"):
			g.start = strings.TrimSpace(strings.TrimPrefix(trimmed, "%start"))
			inToken = false
			continue
		case strings.HasPrefix(trimmed, "%token"):
			inToken = true
			trimmed = strings.TrimPrefix(trimmed, "%token")
		case strings.HasPrefix(trimmed, "%"):
			// Any other directive (%type, %left, %union, ...) ends a token
			// block; %type lists nonterminals, never terminals.
			inToken = false
			continue
		}
		if !inToken {
			continue
		}
		for _, d := range tokenDecls(trimmed) {
			g.declare(d[0], d[1])
		}
	}
}

// tokenDecls reads the name and quoted text of each token one declaration
// line declares, skipping <type> tags and token numbers.
func tokenDecls(line string) [][2]string {
	var out [][2]string
	r := []rune(line)
	for i := 0; i < len(r); {
		switch c := r[i]; {
		case c == '<':
			for i < len(r) && r[i] != '>' {
				i++
			}
			i++
		case c == '"':
			text, next := readDoubleLiteral(r, i)
			if n := len(out); n > 0 && out[n-1][1] == "" {
				out[n-1][1] = text
			}
			i = next
		case isIdentStart(c):
			word, next := readWord(r, i)
			out = append(out, [2]string{word, ""})
			i = next
		default:
			i++
		}
	}
	return out
}

// declare records a token's declared text as its sample, unless the token
// table gives one, and registers the text as an alias a rule body may name
// the token by ("UNDERSCORE_CHARSET"), with the same sample so a table
// override (underscoreCS -> _utf8) is not lost.
func (g *Grammar) declare(name, text string) {
	if text == "" {
		return
	}
	sample, fromTable := g.termText[name]
	if !fromTable && !g.unreachable[name] {
		sample = text
		g.termText[name] = text
	}
	if _, exists := g.termText[text]; !exists {
		if g.unreachable[name] {
			g.unreachable[text] = true
		} else {
			g.termText[text] = sample
		}
	}
}

// grammar-section token kinds.
type gtokKind int

const (
	gIdent   gtokKind = iota // a bareword: a rule name or a token name
	gColon                   // ':' separating a left-hand side from its body
	gPipe                    // '|' separating alternatives
	gLiteral                 // a "..." or '.' terminal; text is its content
)

type gtok struct {
	kind gtokKind
	text string
}

// lexGrammar turns the grammar section into the tokens readRules consumes.
// Semantic actions, %prec directives, and comments carry no symbol and are
// dropped here.
func lexGrammar(gram string) []gtok {
	var out []gtok
	r := []rune(gram)
	i, n := 0, len(r)
	for i < n {
		c := r[i]
		switch {
		case c == ' ' || c == '\t' || c == '\r' || c == '\n':
			i++
		case c == '/' && i+1 < n && r[i+1] == '*':
			i = skipBlockComment(r, i)
		case c == '/' && i+1 < n && r[i+1] == '/':
			i = skipToEOL(r, i)
		case c == '{':
			i = skipAction(r, i)
		case c == '%':
			// A directive. %prec names a precedence symbol the alternative
			// does not hold, so drop it: the symbol is a bareword token name
			// or a quoted operator literal ('+', "||"), so skip whichever
			// follows. %empty marks an empty alternative and needs nothing.
			word, next := readWord(r, i+1)
			i = next
			if word == "prec" {
				for i < n && (r[i] == ' ' || r[i] == '\t') {
					i++
				}
				if i < n {
					switch r[i] {
					case '\'':
						_, i = readCharLiteral(r, i)
					case '"':
						_, i = readDoubleLiteral(r, i)
					default:
						_, i = readWord(r, i)
					}
				}
			}
		case c == '"':
			text, next := readDoubleLiteral(r, i)
			out = append(out, gtok{kind: gLiteral, text: text})
			i = next
		case c == '\'':
			text, next := readCharLiteral(r, i)
			out = append(out, gtok{kind: gLiteral, text: text})
			i = next
		case c == ':':
			out = append(out, gtok{kind: gColon})
			i++
		case c == '|':
			out = append(out, gtok{kind: gPipe})
			i++
		case isIdentStart(c):
			word, next := readWord(r, i)
			out = append(out, gtok{kind: gIdent, text: word})
			i = next
		default:
			// '$' and other action-only characters never appear outside an
			// action; anything else is punctuation the generator would
			// reject. Skip.
			i++
		}
	}
	return out
}

// readRules builds the rules from the token stream. A rule head is an ident
// immediately followed by a colon; every token until the next head belongs
// to the current rule, with '|' starting a fresh alternative.
func (g *Grammar) readRules(toks []gtok) error {
	i, n := 0, len(toks)
	isHead := func(j int) bool {
		return j+1 < n && toks[j].kind == gIdent && toks[j+1].kind == gColon
	}
	for i < n {
		if !isHead(i) {
			return fmt.Errorf("expected a rule head at token %d (%v %q)", i, toks[i].kind, toks[i].text)
		}
		r := &rule{name: toks[i].text, alts: [][]symbol{{}}}
		i += 2
		for i < n && !isHead(i) {
			last := len(r.alts) - 1
			switch toks[i].kind {
			case gPipe:
				r.alts = append(r.alts, []symbol{})
			case gIdent:
				r.alts[last] = append(r.alts[last], symbol{text: toks[i].text})
			case gLiteral:
				r.alts[last] = append(r.alts[last], symbol{text: toks[i].text, literal: true})
			case gColon:
				return fmt.Errorf("unexpected colon in body of rule %q", r.name)
			}
			i++
		}
		if existing, ok := g.byName[r.name]; ok {
			// A repeated left-hand side adds alternatives to the first.
			existing.alts = append(existing.alts, r.alts...)
			continue
		}
		g.rules = append(g.rules, r)
		g.byName[r.name] = r
	}
	if len(g.rules) == 0 {
		return fmt.Errorf("no rules found")
	}
	return nil
}

// resolveSymbols marks every symbol as a nonterminal or a terminal. A
// bareword is a nonterminal when a rule defines it; otherwise it is a token
// name, and it must have a sample. A literal is always a terminal whose text
// is its own content.
func (g *Grammar) resolveSymbols() error {
	var missing []string
	seen := map[string]bool{}
	for _, r := range g.rules {
		for _, alt := range r.alts {
			for k := range alt {
				sym := &alt[k]
				if sym.literal {
					continue
				}
				if _, ok := g.byName[sym.text]; ok {
					sym.nonterm = true
					continue
				}
				if _, ok := g.termText[sym.text]; ok || g.unreachable[sym.text] {
					continue
				}
				if !seen[sym.text] {
					seen[sym.text] = true
					missing = append(missing, sym.text)
				}
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%d terminal(s) have no sample text; add them to tokens.tsv: %s", len(missing), strings.Join(missing, ", "))
	}
	return nil
}

// sample returns the text a terminal symbol lexes from.
func (g *Grammar) sample(sym symbol) string {
	if t, ok := g.termText[sym.text]; ok {
		return t
	}
	return sym.text
}

// --- small lexer helpers ---

func isIdentStart(c rune) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPart(c rune) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9')
}

func readWord(r []rune, i int) (string, int) {
	start := i
	for i < len(r) && isIdentPart(r[i]) {
		i++
	}
	return string(r[start:i]), i
}

func skipBlockComment(r []rune, i int) int {
	i += 2
	for i+1 < len(r) && !(r[i] == '*' && r[i+1] == '/') {
		i++
	}
	return i + 2
}

func skipToEOL(r []rune, i int) int {
	for i < len(r) && r[i] != '\n' {
		i++
	}
	return i
}

// skipAction consumes a balanced { ... } semantic action, respecting the
// strings, characters, and comments inside it so a brace in a literal does
// not end the action early. Go and C actions lex the same way here.
func skipAction(r []rune, i int) int {
	depth := 0
	n := len(r)
	for i < n {
		c := r[i]
		switch {
		case c == '{':
			depth++
			i++
		case c == '}':
			depth--
			i++
			if depth == 0 {
				return i
			}
		case c == '/' && i+1 < n && r[i+1] == '*':
			i = skipBlockComment(r, i)
		case c == '/' && i+1 < n && r[i+1] == '/':
			i = skipToEOL(r, i)
		case c == '"':
			i = skipString(r, i, '"')
		case c == '`':
			i = skipString(r, i, '`')
		case c == '\'':
			i = skipString(r, i, '\'')
		default:
			i++
		}
	}
	return i
}

// skipString consumes a string, raw string, or character literal opened by
// quote at i, returning the index past its closing quote.
func skipString(r []rune, i int, quote rune) int {
	i++ // opening quote
	for i < len(r) {
		c := r[i]
		if c == '\\' && quote != '`' {
			i += 2
			continue
		}
		if c == quote {
			return i + 1
		}
		i++
	}
	return i
}

// readDoubleLiteral reads a "..." grammar literal and returns its content.
func readDoubleLiteral(r []rune, i int) (string, int) {
	return readQuoted(r, i, '"')
}

// readCharLiteral reads a '.' grammar literal and returns its content.
func readCharLiteral(r []rune, i int) (string, int) {
	return readQuoted(r, i, '\'')
}

func readQuoted(r []rune, i int, quote rune) (string, int) {
	i++ // opening quote
	var b strings.Builder
	for i < len(r) {
		c := r[i]
		if c == '\\' && i+1 < len(r) {
			b.WriteRune(r[i+1])
			i += 2
			continue
		}
		if c == quote {
			return b.String(), i + 1
		}
		b.WriteRune(c)
		i++
	}
	return b.String(), i
}

// stripLineComment removes a // comment from a declaration line.
func stripLineComment(line string) string {
	if idx := strings.Index(line, "//"); idx >= 0 {
		return line[:idx]
	}
	return line
}
