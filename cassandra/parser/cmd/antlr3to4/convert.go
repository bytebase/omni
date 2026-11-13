package main

import (
	"fmt"
	"regexp"
	"strings"
)

// A source is one ANTLR 3 grammar file: the name its rules are cited by,
// and its text.
type source struct {
	name string
	text string
}

// convert writes the rules of ANTLR 3 grammars as one ANTLR 4 combined
// grammar: the parser rules of every source in the order given, then the
// lexer rules. The first parser rule is the start rule, so the grammar that
// holds it comes first.
//
// The output keeps what a rule derives and drops what it runs:
//
//   - arguments, return values, labels, scopes, tree operators, and @init
//     and @after blocks;
//   - actions. An alternative whose action does nothing but refuse the
//     input, by addRecognitionError or by throwing SyntaxException, is
//     dropped with it, since the parser accepts no statement through it. An
//     action that refuses under a condition is noted beside its rule, and
//     the condition is lost;
//   - semantic predicates, {...}?, each noted beside its rule, and an action
//     that makes a lexer rule emit another token type, $type = X;
//   - syntactic predicates, (...)=>, which choose between alternatives and
//     derive nothing.
//
// ANTLR 3 reads .* up to what follows it, which ANTLR 4 writes .*?, and an
// action that sets $channel to HIDDEN becomes -> channel(HIDDEN).
//
// Every rule is written under the file and line it came from.
func convert(grammar, header string, sources ...source) (string, error) {
	var parser, lexer []*rule
	for _, src := range sources {
		toks, err := lex(src.text)
		if err != nil {
			return "", fmt.Errorf("%s: %v", src.name, err)
		}
		p := &reader{file: src.name, toks: toks}
		rules, err := p.grammar()
		if err != nil {
			return "", fmt.Errorf("%s: %v", src.name, err)
		}
		for _, r := range rules {
			if r.lexer() {
				lexer = append(lexer, r)
			} else {
				parser = append(parser, r)
			}
		}
	}
	if len(parser) == 0 {
		return "", fmt.Errorf("no parser rule found")
	}
	var b strings.Builder
	b.WriteString(header)
	fmt.Fprintf(&b, "grammar %s;\n", grammar)
	for _, r := range append(parser, lexer...) {
		b.WriteString("\n")
		r.write(&b)
	}
	return b.String(), nil
}

// --- the grammar's tokens ---

type kind int

const (
	tIdent  kind = iota
	tString      // a '...' literal; text is its content, unescaped
	tAction      // a { ... } block; text is the code between the braces
	tArg         // a [ ... ] block
	tPunct
)

type token struct {
	kind kind
	text string
	line int
}

func lex(src string) ([]token, error) {
	var out []token
	r := []rune(src)
	n := len(r)
	line := 1
	for i := 0; i < n; {
		c := r[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r':
			i++
		case c == '/' && i+1 < n && r[i+1] == '/':
			for i < n && r[i] != '\n' {
				i++
			}
		case c == '/' && i+1 < n && r[i+1] == '*':
			j := i + 2
			for j+1 < n && !(r[j] == '*' && r[j+1] == '/') {
				j++
			}
			if j+1 >= n {
				return nil, fmt.Errorf("line %d: comment never closed", line)
			}
			line += strings.Count(string(r[i:j+2]), "\n")
			i = j + 2
		case c == '{' || c == '[':
			end, err := block(r, i)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", line, err)
			}
			k := tAction
			if c == '[' {
				k = tArg
			}
			out = append(out, token{k, string(r[i+1 : end-1]), line})
			line += strings.Count(string(r[i:end]), "\n")
			i = end
		case c == '\'':
			text, end, err := literal(r, i)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", line, err)
			}
			out = append(out, token{tString, text, line})
			i = end
		case c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i
			for j < n && (r[j] == '_' || r[j] >= 'a' && r[j] <= 'z' || r[j] >= 'A' && r[j] <= 'Z' || r[j] >= '0' && r[j] <= '9') {
				j++
			}
			out = append(out, token{tIdent, string(r[i:j]), line})
			i = j
		default:
			op := string(c)
			if i+1 < n {
				switch two := string(r[i : i+2]); two {
				case "..", "+=", "=>", "->", "::":
					op = two
				}
			}
			if len(op) == 1 && !strings.Contains(":;|()?*+~.=^!@,", op) {
				return nil, fmt.Errorf("line %d: unexpected %q", line, c)
			}
			out = append(out, token{tPunct, op, line})
			i += len(op)
		}
	}
	return out, nil
}

// block returns the index past the { ... } or [ ... ] block opened at i.
// The block holds target-language code, so braces inside its strings,
// character literals and comments do not count.
func block(r []rune, i int) (int, error) {
	open, shut := r[i], '}'
	if open == '[' {
		shut = ']'
	}
	depth := 0
	for j := i; j < len(r); j++ {
		switch c := r[j]; {
		case c == '"' || c == '\'':
			for j++; j < len(r) && r[j] != c; j++ {
				if r[j] == '\\' {
					j++
				}
			}
		case c == '/' && j+1 < len(r) && r[j+1] == '/':
			for j < len(r) && r[j] != '\n' {
				j++
			}
		case c == '/' && j+1 < len(r) && r[j+1] == '*':
			for j += 2; j+1 < len(r) && !(r[j] == '*' && r[j+1] == '/'); j++ {
			}
			j++
		case c == open:
			depth++
		case c == shut:
			if depth--; depth == 0 {
				return j + 1, nil
			}
		}
	}
	return 0, fmt.Errorf("%c never closed", open)
}

// literal reads a '...' literal with ANTLR 3's escapes. A backslash before
// any other character stands for that character, as Cassandra's '\*' does.
func literal(r []rune, i int) (string, int, error) {
	var b strings.Builder
	for j := i + 1; j < len(r); j++ {
		switch c := r[j]; {
		case c == '\'':
			return b.String(), j + 1, nil
		case c == '\n':
			return "", 0, fmt.Errorf("newline in a literal")
		case c == '\\' && j+1 < len(r):
			j++
			switch e := r[j]; e {
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
				if j+4 >= len(r) {
					return "", 0, fmt.Errorf("bad \\u escape")
				}
				var v rune
				if _, err := fmt.Sscanf(string(r[j+1:j+5]), "%04x", &v); err != nil {
					return "", 0, fmt.Errorf("bad \\u escape")
				}
				b.WriteRune(v)
				j += 4
			default:
				b.WriteRune(e)
			}
		default:
			b.WriteRune(c)
		}
	}
	return "", 0, fmt.Errorf("literal never closed")
}

// --- rules ---

type rule struct {
	file     string
	line     int
	name     string
	fragment bool
	alts     []alt
	notes    []string // what the conversion lost, each citing its line
	hidden   bool     // an action sent the token to the hidden channel
}

func (r *rule) lexer() bool { return r.name[0] >= 'A' && r.name[0] <= 'Z' }

type alt []elem

// An elem is a rule or token reference, a literal, a range of characters,
// the . wildcard, or a group, with the operator after it and a ~ before.
type elem struct {
	ref    string
	lit    string
	isLit  bool
	hi     string // the end of a range lit..hi
	ranged bool
	any    bool
	group  []alt
	not    bool
	op     string // "", "?", "*" or "+"
}

type reader struct {
	file string
	toks []token
	i    int
	rule *rule
}

func (p *reader) peek() token {
	if p.i < len(p.toks) {
		return p.toks[p.i]
	}
	return token{kind: tPunct, text: "<end>", line: -1}
}

func (p *reader) is(k kind, text string) bool {
	t := p.peek()
	return t.kind == k && t.text == text
}

func (p *reader) errorf(format string, args ...any) error {
	return fmt.Errorf("line %d: %s", p.peek().line, fmt.Sprintf(format, args...))
}

func (p *reader) expect(text string) error {
	if !p.is(tPunct, text) {
		return p.errorf("want %q, got %q", text, p.peek().text)
	}
	p.i++
	return nil
}

func (p *reader) skipTo(text string) error {
	for ; p.i < len(p.toks); p.i++ {
		if p.is(tPunct, text) {
			p.i++
			return nil
		}
	}
	return fmt.Errorf("want %q before the end", text)
}

func (p *reader) grammar() ([]*rule, error) {
	var rules []*rule
	for p.i < len(p.toks) {
		t := p.peek()
		next := token{}
		if p.i+1 < len(p.toks) {
			next = p.toks[p.i+1]
		}
		switch {
		case t.kind == tIdent && (t.text == "grammar" || t.text == "import" ||
			(t.text == "lexer" || t.text == "parser") && next.kind == tIdent && next.text == "grammar"):
			if err := p.skipTo(";"); err != nil {
				return nil, err
			}
		case t.kind == tIdent && (t.text == "options" || t.text == "tokens") && next.kind == tAction:
			p.i += 2
		case t.kind == tPunct && t.text == "@":
			// @header {...}, @lexer::members {...}
			for p.i++; p.i < len(p.toks) && p.peek().kind != tAction; p.i++ {
			}
			p.i++
		case t.kind == tIdent:
			r, err := p.readRule()
			if err != nil {
				return nil, err
			}
			rules = append(rules, r)
		default:
			return nil, p.errorf("want a rule, got %q", t.text)
		}
	}
	return rules, nil
}

func (p *reader) readRule() (*rule, error) {
	r := &rule{file: p.file}
	if p.is(tIdent, "fragment") {
		r.fragment = true
		p.i++
	}
	t := p.peek()
	if t.kind != tIdent {
		return nil, p.errorf("want a rule name, got %q", t.text)
	}
	r.name, r.line = t.text, t.line
	p.i++
	// Arguments, returns, throws, options, scope, @init and @after: what the
	// generated parser runs, up to the colon that opens the alternatives.
	for !p.is(tPunct, ":") {
		if p.i >= len(p.toks) {
			return nil, fmt.Errorf("rule %s: no colon", r.name)
		}
		p.i++
	}
	p.i++
	p.rule = r
	alts, err := p.alts()
	if err != nil {
		return nil, fmt.Errorf("rule %s: %v", r.name, err)
	}
	if len(alts) == 0 {
		return nil, fmt.Errorf("rule %s: every alternative refuses its input", r.name)
	}
	r.alts = alts
	if err := p.expect(";"); err != nil {
		return nil, fmt.Errorf("rule %s: %v", r.name, err)
	}
	return r, nil
}

// refusal matches an action that does nothing but refuse the input.
var refusal = regexp.MustCompile(`(?s)^\s*(addRecognitionError\(.*\)|throw new SyntaxException\(.*\))\s*;?\s*$`)

// refuses matches an action that may refuse the input.
var refuses = regexp.MustCompile(`addRecognitionError|SyntaxException`)

// retypes matches an action that makes a lexer rule emit another token.
var retypes = regexp.MustCompile(`\$type\s*=`)

// hides matches an action that sends a token to the hidden channel.
var hides = regexp.MustCompile(`\$channel\s*=\s*HIDDEN`)

// alts reads alternatives up to a ')' or ';', which it leaves unread. An
// alternative its action refuses is left out.
func (p *reader) alts() ([]alt, error) {
	var out []alt
	for {
		start := p.i
		line := p.peek().line
		a, refused, err := p.alt()
		if err != nil {
			return nil, err
		}
		if refused {
			p.rule.notes = append(p.rule.notes, fmt.Sprintf("%s:%d: alternative dropped, its action refuses every input: %s", p.file, line, sourceText(p.toks[start:p.i])))
		} else {
			out = append(out, a)
		}
		if !p.is(tPunct, "|") {
			return out, nil
		}
		p.i++
	}
}

func (p *reader) alt() (a alt, refused bool, err error) {
	for {
		t := p.peek()
		switch {
		case t.kind == tPunct && (t.text == "|" || t.text == ")" || t.text == ";"):
			return a, refused, nil
		case t.kind == tPunct && t.text == "<end>":
			return nil, false, fmt.Errorf("unexpected end of grammar")
		case t.kind == tAction:
			p.i++
			if p.is(tPunct, "?") {
				p.i++
				if p.is(tPunct, "=>") {
					p.i++
				}
				p.rule.notes = append(p.rule.notes, fmt.Sprintf("%s:%d: predicate dropped: {%s}?", p.file, t.line, oneLine(t.text)))
				continue
			}
			switch {
			case refusal.MatchString(t.text) && !strings.Contains(t.text, "if"):
				refused = true
			case refuses.MatchString(t.text):
				p.rule.notes = append(p.rule.notes, fmt.Sprintf("%s:%d: action may refuse the input: {%s}", p.file, t.line, oneLine(t.text)))
			case hides.MatchString(t.text):
				p.rule.hidden = true
			case retypes.MatchString(t.text):
				p.rule.notes = append(p.rule.notes, fmt.Sprintf("%s:%d: action dropped, it emits another token type: {%s}", p.file, t.line, oneLine(t.text)))
			}
			continue
		case t.kind == tPunct && t.text == "(":
			// A syntactic predicate, ( ... ) =>, derives nothing.
			if end := p.closing(); end+1 < len(p.toks) && p.toks[end+1].kind == tPunct && p.toks[end+1].text == "=>" {
				p.i = end + 2
				continue
			}
		}
		e, ok, err := p.elem()
		if err != nil {
			return nil, false, err
		}
		if !ok {
			// A group every alternative of which refuses.
			refused = true
			continue
		}
		a = append(a, e)
	}
}

// closing returns the index of the ')' that closes the '(' at p.i.
func (p *reader) closing() int {
	depth := 0
	for j := p.i; j < len(p.toks); j++ {
		if p.toks[j].kind != tPunct {
			continue
		}
		switch p.toks[j].text {
		case "(":
			depth++
		case ")":
			if depth--; depth == 0 {
				return j
			}
		}
	}
	return len(p.toks)
}

// elem reads one element. ok is false for a group left with no alternative.
func (p *reader) elem() (e elem, ok bool, err error) {
	// A label: x=atom or x+=atom.
	if t := p.peek(); t.kind == tIdent && p.i+1 < len(p.toks) {
		if next := p.toks[p.i+1]; next.kind == tPunct && (next.text == "=" || next.text == "+=") {
			p.i += 2
		}
	}
	if p.is(tPunct, "~") {
		e.not = true
		p.i++
	}
	t := p.peek()
	switch {
	case t.kind == tIdent:
		e.ref = t.text
		p.i++
		if p.peek().kind == tArg {
			p.i++
		}
	case t.kind == tString:
		e.lit, e.isLit = t.text, true
		p.i++
		if p.is(tPunct, "..") {
			p.i++
			hi := p.peek()
			if hi.kind != tString {
				return e, false, p.errorf("want a literal after ..")
			}
			e.hi, e.ranged = hi.text, true
			p.i++
		}
	case t.kind == tPunct && t.text == ".":
		e.any = true
		p.i++
	case t.kind == tPunct && t.text == "(":
		p.i++
		// ( options { greedy=false; } : ... )
		if p.is(tIdent, "options") {
			if err := p.skipTo(":"); err != nil {
				return e, false, err
			}
		}
		group, err := p.alts()
		if err != nil {
			return e, false, err
		}
		if err := p.expect(")"); err != nil {
			return e, false, err
		}
		e.group = group
	default:
		return e, false, p.errorf("unexpected %q", t.text)
	}
	if t := p.peek(); t.kind == tPunct && (t.text == "?" || t.text == "*" || t.text == "+") {
		e.op = t.text
		p.i++
	}
	// Tree operators build an AST and derive nothing.
	for p.is(tPunct, "^") || p.is(tPunct, "!") {
		p.i++
	}
	if e.ref == "" && !e.isLit && !e.any && len(e.group) == 0 {
		return e, e.op == "?" || e.op == "*", nil
	}
	return e, true, nil
}

// --- writing ANTLR 4 ---

func (r *rule) write(b *strings.Builder) {
	fmt.Fprintf(b, "// %s:%d\n", r.file, r.line)
	for _, n := range r.notes {
		fmt.Fprintf(b, "// %s\n", n)
	}
	if r.fragment {
		b.WriteString("fragment ")
	}
	b.WriteString(r.name + "\n")
	for i, a := range r.alts {
		sep := "    : "
		if i > 0 {
			sep = "    | "
		}
		text := altText(a)
		if text == "" {
			text = "/* empty */"
		}
		b.WriteString(sep + text + "\n")
	}
	if r.hidden {
		b.WriteString("      -> channel(HIDDEN)\n")
	}
	b.WriteString("    ;\n")
}

func altText(a alt) string {
	parts := make([]string, len(a))
	for i, e := range a {
		parts[i] = e.text()
	}
	return strings.Join(parts, " ")
}

func (e elem) text() string {
	var s string
	switch {
	case e.any && e.op != "":
		// ANTLR 3 stops .* and .+ at what follows; ANTLR 4 writes that .*?
		return "." + e.op + "?"
	case e.any:
		s = "."
	case e.ranged:
		s = quote(e.lit) + ".." + quote(e.hi)
	case e.isLit:
		s = quote(e.lit)
	case e.ref != "":
		s = e.ref
	default:
		if len(e.group) == 1 && len(e.group[0]) == 1 && e.group[0][0].any && e.group[0][0].op == "" && e.op != "" {
			return "." + e.op + "?"
		}
		alts := make([]string, len(e.group))
		for i, a := range e.group {
			alts[i] = altText(a)
		}
		s = "( " + strings.Join(alts, " | ") + " )"
	}
	if e.not {
		s = "~" + s
	}
	return s + e.op
}

// quote writes a literal's content as ANTLR 4 does.
func quote(s string) string {
	var b strings.Builder
	b.WriteByte('\'')
	for _, c := range s {
		switch {
		case c == '\'' || c == '\\':
			b.WriteByte('\\')
			b.WriteRune(c)
		case c == '\n':
			b.WriteString(`\n`)
		case c == '\r':
			b.WriteString(`\r`)
		case c == '\t':
			b.WriteString(`\t`)
		case c == '\b':
			b.WriteString(`\b`)
		case c == '\f':
			b.WriteString(`\f`)
		case c < ' ' || c > '~':
			fmt.Fprintf(&b, `\u%04X`, c)
		default:
			b.WriteRune(c)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

// sourceText writes tokens back roughly as the grammar had them, for a note.
func sourceText(toks []token) string {
	parts := make([]string, len(toks))
	for i, t := range toks {
		switch t.kind {
		case tString:
			parts[i] = quote(t.text)
		case tAction:
			parts[i] = "{" + oneLine(t.text) + "}"
		case tArg:
			parts[i] = "[" + oneLine(t.text) + "]"
		default:
			parts[i] = t.text
		}
	}
	return strings.Join(parts, " ")
}

var spaces = regexp.MustCompile(`\s+`)

// oneLine writes code on one line, so a note stays one comment.
func oneLine(s string) string { return strings.TrimSpace(spaces.ReplaceAllString(s, " ")) }
