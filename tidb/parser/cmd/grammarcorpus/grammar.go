package main

import (
	"fmt"
	"strings"
)

// A Grammar is parser.y reduced to what sentence generation needs: the rules
// of the grammar section, and a sample of text for every terminal. Everything
// goyacc adds for the Go parser it builds -- the prologue and epilogue code,
// the %union, the semantic actions, the precedence declarations -- is dropped,
// because none of it changes which strings the grammar derives.
type Grammar struct {
	// rules in the order parser.y lists them; the first is the start symbol.
	rules []*Rule
	// byName indexes rules by their left-hand side.
	byName map[string]*Rule
	// termText gives the sample text each terminal lexes from, keyed by the
	// terminal's symbol text (a token name, or a quoted/char literal's content).
	termText map[string]string
}

// A Rule is one left-hand side and its alternatives, in source order.
type Rule struct {
	name string
	alts [][]Symbol
}

// A Symbol is one element of an alternative. A nonterminal names another rule;
// a terminal is a token the lexer produces, identified by text: a token name
// (identifier, intLit, a keyword), or the content of a "..." or '.' literal.
type Symbol struct {
	text    string
	nonterm bool
	// literal is set when the symbol came from a "..." or '.' grammar literal,
	// whose content is the text to emit directly, never a token name to look up.
	literal bool
}

// specialTermText overrides the sample text for the tokens whose declared text
// in parser.y describes the token rather than spelling it: the literal classes
// and the few named operators the lexer forms from more than their word. Every
// other terminal lexes from its own text, so the generator uses that directly.
var specialTermText = map[string]string{
	"identifier":         "ident",
	"stringLit":          "'s'",
	"floatLit":           "1e0", // an exponent makes it a float; without one TiDB lexes a decimal
	"decLit":             "1.5",
	"intLit":             "1",
	"hexLit":             "0x1",
	"bitLit":             "b'1'",
	"singleAtIdentifier": "@v",
	"doubleAtIdentifier": "@@v",
	"underscoreCS":       "_utf8",
	"paramMarker":        "?",
	"not2":               "NOT",
	"invalid":            "", // a lexer-only error token; never reached in a sentence
	"hintComment":        "", // optimizer hints never change acceptance (design non-goal)

	// pipesAsOr is how the lexer reads "||" when it means OR rather than
	// concatenation; the text is the same either way.
	"pipesAsOr": "||",

	// The builtinX tokens are declared without text: the lexer forms each one
	// when a function's name is read immediately before '('. The sample is that
	// name; the generator keeps the name and '(' tight, so it lexes back.
	"builtinCast":                "CAST",
	"builtinCount":               "COUNT",
	"builtinCurDate":             "CURDATE",
	"builtinCurTime":             "CURTIME",
	"builtinNow":                 "NOW",
	"builtinSysDate":             "SYSDATE",
	"builtinUser":                "USER",
	"builtinExtract":             "EXTRACT",
	"builtinPosition":            "POSITION",
	"builtinSubstring":           "SUBSTRING",
	"builtinTrim":                "TRIM",
	"builtinTranslate":           "TRANSLATE",
	"builtinDateAdd":             "DATE_ADD",
	"builtinDateSub":             "DATE_SUB",
	"builtinApproxCountDistinct": "APPROX_COUNT_DISTINCT",
	"builtinApproxPercentile":    "APPROX_PERCENTILE",
	"builtinBitAnd":              "BIT_AND",
	"builtinBitOr":               "BIT_OR",
	"builtinBitXor":              "BIT_XOR",
	"builtinGroupConcat":         "GROUP_CONCAT",
	"builtinMax":                 "MAX",
	"builtinMin":                 "MIN",
	"builtinSum":                 "SUM",
	"builtinStddevPop":           "STDDEV_POP",
	"builtinStddevSamp":          "STDDEV_SAMP",
	"builtinVarPop":              "VAR_POP",
	"builtinVarSamp":             "VAR_SAMP",
}

// parseGrammar reads the whole text of parser.y and returns its Grammar.
func parseGrammar(src string) (*Grammar, error) {
	decls, gram, err := splitSections(src)
	if err != nil {
		return nil, err
	}
	g := &Grammar{byName: map[string]*Rule{}, termText: map[string]string{}}
	g.readTokenText(decls)

	toks := lexGrammar(gram)
	if err := g.readRules(toks); err != nil {
		return nil, err
	}
	if err := g.resolveSymbols(); err != nil {
		return nil, err
	}
	return g, nil
}

// splitSections returns the declaration section and the grammar section of a
// goyacc file, the two halves either side of the first "%%". The third section
// after the second "%%", Go epilogue code, is dropped.
func splitSections(src string) (decls, gram string, err error) {
	first := indexOfLine(src, "%%", 0)
	if first < 0 {
		return "", "", fmt.Errorf("no %%%% section marker found")
	}
	second := indexOfLine(src, "%%", first+1)
	if second < 0 {
		return "", "", fmt.Errorf("only one %%%% section marker found")
	}
	lines := strings.Split(src, "\n")
	return strings.Join(lines[:first], "\n"), strings.Join(lines[first+1:second], "\n"), nil
}

// indexOfLine returns the index, in lines, of the first line at or after start
// whose trimmed text equals want, or -1.
func indexOfLine(src, want string, start int) int {
	for i, line := range strings.Split(src, "\n") {
		if i >= start && strings.TrimSpace(line) == want {
			return i
		}
	}
	return -1
}

// readTokenText records the sample text of every terminal declared in the
// declaration section. A %token line names a token and gives its text in
// quotes; the text is the literal the lexer reads for a keyword or operator,
// or a description for a literal class, which specialTermText then overrides.
func (g *Grammar) readTokenText(decls string) {
	for name, text := range specialTermText {
		g.termText[name] = text
	}
	inToken := false
	for _, line := range strings.Split(decls, "\n") {
		trimmed := strings.TrimSpace(stripLineComment(line))
		switch {
		case strings.HasPrefix(trimmed, "%token"):
			inToken = true
			continue
		case strings.HasPrefix(trimmed, "%"):
			// Any other directive (%type, %left, %precedence, %union, ...)
			// ends a token block; %type lists nonterminals, never terminals.
			inToken = false
			continue
		}
		if !inToken || trimmed == "" {
			continue
		}
		// A declaration line is: tokenName "declared text".
		name, rest, ok := strings.Cut(trimmed, " ")
		if !ok {
			continue
		}
		q := strings.TrimSpace(rest)
		if !strings.HasPrefix(q, `"`) {
			continue
		}
		text, ok := unquoteDouble(q)
		if !ok {
			continue
		}
		sample := text
		if s, special := specialTermText[name]; special {
			sample = s
		} else {
			g.termText[name] = text
		}
		// A body may reference a token by its quoted alias ("UNDERSCORE_CHARSET")
		// rather than its name; register the alias with the same sample so a
		// special token's override (underscoreCS -> _utf8) is not lost.
		if _, exists := g.termText[text]; !exists {
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
			// A directive. %prec names a precedence symbol the alternative does
			// not want, so drop it: the symbol is either a bareword token name or
			// a quoted operator literal ('+', "||"), so skip whichever follows.
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
			// action; anything else is punctuation goyacc would reject. Skip.
			i++
		}
	}
	return out
}

// readRules builds the rules from the token stream. A rule head is an ident
// immediately followed by a colon; every token until the next head belongs to
// the current rule, with '|' starting a fresh alternative.
func (g *Grammar) readRules(toks []gtok) error {
	i, n := 0, len(toks)
	isHead := func(j int) bool {
		return j+1 < n && toks[j].kind == gIdent && toks[j+1].kind == gColon
	}
	for i < n {
		if !isHead(i) {
			return fmt.Errorf("expected a rule head at token %d (%v %q)", i, toks[i].kind, toks[i].text)
		}
		rule := &Rule{name: toks[i].text, alts: [][]Symbol{{}}}
		i += 2
		for i < n && !isHead(i) {
			switch toks[i].kind {
			case gPipe:
				rule.alts = append(rule.alts, []Symbol{})
			case gIdent:
				last := len(rule.alts) - 1
				rule.alts[last] = append(rule.alts[last], Symbol{text: toks[i].text})
			case gLiteral:
				last := len(rule.alts) - 1
				rule.alts[last] = append(rule.alts[last], Symbol{text: toks[i].text, literal: true})
			case gColon:
				return fmt.Errorf("unexpected colon in body of rule %q", rule.name)
			}
			i++
		}
		if existing, ok := g.byName[rule.name]; ok {
			// goyacc forbids a repeated left-hand side; if one appears, merge.
			existing.alts = append(existing.alts, rule.alts...)
			continue
		}
		g.rules = append(g.rules, rule)
		g.byName[rule.name] = rule
	}
	if len(g.rules) == 0 {
		return fmt.Errorf("no rules found")
	}
	return nil
}

// resolveSymbols marks every symbol as a nonterminal or a terminal, and fills
// in each terminal's sample text. A bareword is a nonterminal when a rule
// defines it; otherwise it is a token name. A quoted or char literal is always
// a terminal whose text is its own content.
func (g *Grammar) resolveSymbols() error {
	var missing []string
	seen := map[string]bool{}
	for _, rule := range g.rules {
		for _, alt := range rule.alts {
			for k := range alt {
				sym := &alt[k]
				// A literal's content is the text to emit; it is never a rule
				// name or a token name.
				if sym.literal {
					continue
				}
				if _, ok := g.byName[sym.text]; ok {
					sym.nonterm = true
					continue
				}
				// A bareword terminal resolves through the token table.
				if _, ok := g.termText[sym.text]; ok {
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
		return fmt.Errorf("%d terminal(s) have no sample text: %s", len(missing), strings.Join(missing, ", "))
	}
	return nil
}

// sample returns the text a terminal symbol lexes from.
func (g *Grammar) sample(sym Symbol) string {
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

// skipAction consumes a balanced { ... } semantic action, respecting the Go
// strings, runes, and comments inside it so a brace in a literal does not end
// the action early.
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
			i = skipGoString(r, i, '"')
		case c == '`':
			i = skipGoString(r, i, '`')
		case c == '\'':
			i = skipGoString(r, i, '\'')
		default:
			i++
		}
	}
	return i
}

// skipGoString consumes a Go string, raw string, or rune literal opened by
// quote at i, returning the index past its closing quote.
func skipGoString(r []rune, i int, quote rune) int {
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
	i++ // opening quote
	var b strings.Builder
	for i < len(r) {
		c := r[i]
		if c == '\\' && i+1 < len(r) {
			b.WriteRune(r[i+1])
			i += 2
			continue
		}
		if c == '"' {
			return b.String(), i + 1
		}
		b.WriteRune(c)
		i++
	}
	return b.String(), i
}

// readCharLiteral reads a '.' grammar literal and returns its single character.
func readCharLiteral(r []rune, i int) (string, int) {
	i++ // opening quote
	var b strings.Builder
	for i < len(r) {
		c := r[i]
		if c == '\\' && i+1 < len(r) {
			b.WriteRune(r[i+1])
			i += 2
			continue
		}
		if c == '\'' {
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

// unquoteDouble returns the content of a leading "..." in s.
func unquoteDouble(s string) (string, bool) {
	if !strings.HasPrefix(s, `"`) {
		return "", false
	}
	r := []rune(s)
	text, _ := readDoubleLiteral(r, 0)
	return text, true
}
