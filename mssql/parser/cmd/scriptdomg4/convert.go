package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// A source is one vendor file: the name it is cited by, and its text.
type source struct {
	name string
	text string
}

// startRule is ScriptDOM's entry point for a script, TSqlParser.Parse's: its
// batches, between GO lines. The generator starts from the first rule, so it
// is written first.
const startRule = "script"

// convert writes ScriptDOM's ANTLR 2 grammar as one ANTLR 4 combined
// grammar: the parser rules, start rule first, then a lexer rule for each
// keyword, then the lexer rules.
//
// The output keeps what a rule derives and drops what it runs:
//
//   - arguments, return values, labels, init blocks and exception handlers;
//   - actions, but for one thing they decide. ScriptDOM reads most of its
//     words as an Identifier and checks the text in an action, Match(t,
//     CodeGenerationSupporter.X), so the grammar alone would derive any name
//     there. A token whose action matches it against one or two words, at
//     the action's top level, is written as those words, spelled as
//     CodeGenerationSupporter.cs spells them. A match under a condition, and
//     a word read by an option helper (ParseOption and the like), stay an
//     Identifier;
//   - an alternative whose action does nothing but throw a parse error is
//     dropped, since the parser accepts no statement through it. An action
//     that throws under a condition is noted beside its rule;
//   - semantic predicates, {...}?, each noted beside its rule;
//   - syntactic predicates, (...)=>, which choose between alternatives and
//     derive nothing.
//
// Keywords are the tokens the two tokens { } blocks give a text: the
// TSqlTokenTypes.g block every version shares, and TSql160.g's lexer's
// own. ScriptDOM's lexer reads an identifier and turns it into a keyword
// when its text is one (testLiterals), case-insensitively, so each is
// written as a lexer rule ahead of the others under caseInsensitive. A lexer
// rule that only names a type another rule's action sets ($setType), such as
// Integer, has no text of its own: it is listed under the rule that sets it,
// and tokens.tsv samples it.
//
// Every rule is written under the file and line it came from.
func convert(g, tokenTypes, supporter source) (string, error) {
	words := supporterWords(supporter.text)
	keywords, err := tokenTexts(tokenTypes)
	if err != nil {
		return "", err
	}
	toks, err := lex(g.text)
	if err != nil {
		return "", fmt.Errorf("%s: %v", g.name, err)
	}
	p := &reader{file: g.name, toks: toks, words: words}
	if err := p.grammar(); err != nil {
		return "", fmt.Errorf("%s: %v", g.name, err)
	}
	keywords = append(keywords, p.keywords...)
	p.parser = reachable(specialize(prune(p.parser)))
	var start *rule
	var parser []*rule
	for _, r := range p.parser {
		if r.name == startRule {
			start = r
		} else {
			parser = append(parser, r)
		}
	}
	if start == nil {
		return "", fmt.Errorf("%s: no rule %s", g.name, startRule)
	}

	var b strings.Builder
	b.WriteString(header)
	b.WriteString("grammar TSql160;\n\noptions { caseInsensitive = true; }\n")
	for _, r := range append([]*rule{start}, parser...) {
		b.WriteString("\n")
		r.write(&b)
	}
	b.WriteString("\n// Keywords: the tokens { } blocks, which ScriptDOM's lexer matches an\n// identifier's text against.\n")
	seen := map[string]bool{}
	for _, k := range keywords {
		if seen[k.name] {
			continue
		}
		seen[k.name] = true
		fmt.Fprintf(&b, "\n// %s:%d\n%s\n    : %s\n    ;\n", k.file, k.line, k.name, quote(strings.ToUpper(k.text)))
	}
	for _, r := range p.lexer {
		if len(r.alts) == 0 {
			continue
		}
		b.WriteString("\n")
		r.write(&b)
	}
	if len(p.retyped) > 0 {
		b.WriteString("\n// Tokens with no rule of their own: an action of the rule named sets\n// their type ($setType), and tokens.tsv samples them.\n")
		for _, n := range p.retyped {
			fmt.Fprintf(&b, "// %s\n", n)
		}
	}
	return b.String(), nil
}

// header opens TSql160.g4: the notice of the files it is derived from, and
// what was changed.
const header = `// MIT License
//
// Copyright (c) Microsoft Corporation.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This file is derived from Microsoft ScriptDOM's TSql160.g, TSqlTokenTypes.g
// and CodeGenerationSupporter.cs, which are beside it as published, with their
// license (ScriptDom.LICENSE). mssql/parser/cmd/scriptdomg4 wrote it from
// them, converting ANTLR 2 to ANTLR 4 and keeping only what each rule derives:
// actions, arguments, labels and predicates are gone, a word an action matches
// an identifier against is written as that word, and what else was lost is
// noted beside its rule. It generates test statements; it is not a parser.

`

// supporterWords reads CodeGenerationSupporter.cs: each constant's name and
// its text.
func supporterWords(cs string) map[string]string {
	words := map[string]string{}
	for _, m := range constString.FindAllStringSubmatch(cs, -1) {
		words[m[1]] = m[2]
	}
	return words
}

var constString = regexp.MustCompile(`const string (\w+)\s*=\s*"([^"\\]*)"\s*;`)

// A keyword is a token a tokens { } block gives a text.
type keyword struct {
	file string
	line int
	name string
	text string
}

var tokenText = regexp.MustCompile(`^\s*(\w+)\s*=\s*"([^"]*)"\s*;`)

// tokenTexts reads the tokens { } block of a file such as TSqlTokenTypes.g.
func tokenTexts(src source) ([]keyword, error) {
	open := regexp.MustCompile(`(?m)^tokens\s*\{`).FindStringIndex(src.text)
	if open == nil {
		return nil, fmt.Errorf("%s: no tokens block", src.name)
	}
	return blockTexts(src.name, src.text, open[1]), nil
}

// blockTexts reads Name = "text"; lines from the tokens block whose body
// starts at offset at of text, up to its closing brace.
func blockTexts(file, text string, at int) []keyword {
	line := 1 + strings.Count(text[:at], "\n")
	var out []keyword
	for _, l := range strings.Split(text[at:], "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "}") {
			break
		}
		if m := tokenText.FindStringSubmatch(l); m != nil {
			out = append(out, keyword{file, line, m[1], m[2]})
		}
		line++
	}
	return out
}

// --- the grammar's tokens ---

type kind int

const (
	tIdent  kind = iota
	tString      // a "..." literal or a '...' character; text is its content, unescaped
	tAction      // a { ... } block; text is the code between the braces
	tArg         // a [ ... ] block
	tPunct
)

type token struct {
	kind kind
	text string
	line int
	off  int // byte offset of the token in the file
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
		case c == ' ' || c == '\t' || c == '\r' || c == '\f':
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
			out = append(out, token{k, string(r[i+1 : end-1]), line, i})
			line += strings.Count(string(r[i:end]), "\n")
			i = end
		case c == '\'' || c == '"':
			text, end, err := literal(r, i)
			if err != nil {
				return nil, fmt.Errorf("line %d: %v", line, err)
			}
			out = append(out, token{tString, text, line, i})
			i = end
		case c == '_' || c == '$' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i
			for j < n && (r[j] == '_' || r[j] == '$' || r[j] >= 'a' && r[j] <= 'z' || r[j] >= 'A' && r[j] <= 'Z' || r[j] >= '0' && r[j] <= '9') {
				j++
			}
			out = append(out, token{tIdent, string(r[i:j]), line, i})
			i = j
		default:
			op := string(c)
			if i+1 < n {
				switch two := string(r[i : i+2]); two {
				case "..", "=>":
					op = two
				}
			}
			if len(op) == 1 && !strings.Contains(":;|()?*+~.=^!,", op) {
				return nil, fmt.Errorf("line %d: unexpected %q", line, c)
			}
			out = append(out, token{tPunct, op, line, i})
			i += len([]rune(op))
		}
	}
	return out, nil
}

// block returns the index past the { ... } or [ ... ] block opened at i.
// The block holds C#, so braces inside its strings, character literals and
// comments do not count.
func block(r []rune, i int) (int, error) {
	open, shut := r[i], '}'
	if open == '[' {
		shut = ']'
	}
	depth := 0
	for j := i; j < len(r); j++ {
		switch c := r[j]; {
		case c == '@' && j+1 < len(r) && r[j+1] == '"':
			// A verbatim string: "" is a quote in it.
			for j += 2; j < len(r); j++ {
				if r[j] == '"' {
					if j+1 < len(r) && r[j+1] == '"' {
						j++
						continue
					}
					break
				}
			}
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

// literal reads a "..." or '...' literal with ANTLR 2's escapes.
func literal(r []rune, i int) (string, int, error) {
	q := r[i]
	var b strings.Builder
	for j := i + 1; j < len(r); j++ {
		switch c := r[j]; {
		case c == q:
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
	lexer    bool
	fragment bool
	alts     []alt
	notes    []string        // what the conversion lost, each citing its line
	hidden   bool            // whitespace and comments, which the parser never sees
	refers   map[string]bool // dropped rules an alternative of this one needed
}

type alt []elem

// An elem is a rule or token reference, a literal, a range of characters,
// the . wildcard, or a group, with the operator after it and a ~ before.
type elem struct {
	label  string
	ref    string
	lit    string
	isLit  bool
	hi     string // the end of a range lit..hi
	ranged bool
	any    bool
	group  []alt
	not    bool
	op     string // "", "?", "*" or "+"

	// next holds the words a predicate before the element requires the
	// next token to spell, {NextTokenMatches(CodeGenerationSupporter.X)}?,
	// and pred the predicate's text and line.
	next []string
	pred token
}

type reader struct {
	file  string
	toks  []token
	i     int
	words map[string]string
	rule  *rule
	inLex bool // reading the lexer class

	parser, lexer []*rule
	keywords      []keyword
	retyped       []string
}

func (p *reader) peek() token {
	if p.i < len(p.toks) {
		return p.toks[p.i]
	}
	return token{kind: tPunct, text: "<end>", line: -1}
}

func (p *reader) at(k int) token {
	if p.i+k < len(p.toks) {
		return p.toks[p.i+k]
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

func (p *reader) grammar() error {
	for p.i < len(p.toks) {
		t := p.peek()
		switch {
		case t.kind == tAction:
			// The class's members, or the file's usings.
			p.i++
		case t.kind == tIdent && t.text == "class":
			p.inLex = p.at(3).kind == tIdent && p.at(3).text == "Lexer"
			if err := p.skipTo(";"); err != nil {
				return err
			}
		case t.kind == tIdent && (t.text == "options" || t.text == "header") && p.at(1).kind == tAction:
			p.i += 2
		case t.kind == tIdent && t.text == "tokens" && p.at(1).kind == tAction:
			// The block's body starts after its brace.
			for _, k := range blockTexts(p.file, p.at(1).text, 0) {
				k.line += p.at(1).line - 1
				p.keywords = append(p.keywords, k)
			}
			p.i += 2
		case t.kind == tIdent:
			r, err := p.readRule()
			if err != nil {
				return err
			}
			if r.lexer {
				p.lexer = append(p.lexer, r)
			} else {
				p.parser = append(p.parser, r)
			}
		default:
			return p.errorf("want a rule, got %q", t.text)
		}
	}
	return nil
}

func (p *reader) readRule() (*rule, error) {
	r := &rule{file: p.file, lexer: p.inLex}
	for p.is(tIdent, "protected") || p.is(tIdent, "public") || p.is(tIdent, "private") {
		r.fragment = r.fragment || p.peek().text == "protected"
		p.i++
	}
	t := p.peek()
	if t.kind != tIdent {
		return nil, p.errorf("want a rule name, got %q", t.text)
	}
	r.name, r.line = t.text, t.line
	p.i++
	// Arguments, returns, throws, options and the init block: what the
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
	if err := p.expect(";"); err != nil {
		return nil, fmt.Errorf("rule %s: %v", r.name, err)
	}
	p.skipHandlers()
	empty := true
	for _, a := range alts {
		if len(a) > 0 {
			empty = false
		}
	}
	switch {
	case r.lexer && empty:
		// Integer :; and the like: a type another rule sets.
		r.alts = nil
		p.retyped = append(p.retyped, fmt.Sprintf("%s: %s:%d", r.name, p.file, r.line))
		return r, nil
	case len(alts) == 0:
		// odbcInitiator: a rule the parser enters only to refuse. A
		// reference to it refuses too (prune).
		p.note(r.line, "rule dropped, every alternative refuses its input")
	}
	r.alts = alts
	switch r.name {
	case "WhiteSpace", "SingleLineComment", "MultilineComment":
		r.hidden = true
	}
	return r, nil
}

// skipHandlers skips exception handlers: exception [label] catch [...] {...}.
func (p *reader) skipHandlers() {
	for p.is(tIdent, "exception") {
		p.i++
		if p.peek().kind == tArg {
			p.i++
		}
		for p.is(tIdent, "catch") {
			p.i++
			if p.peek().kind == tArg {
				p.i++
			}
			if p.peek().kind == tAction {
				p.i++
			}
		}
	}
}

// refusal matches an action that does nothing but throw a parse error.
var refusal = regexp.MustCompile(`(?s)^\s*(Throw\w*Exception\([^;]*\)|throw [^;]*)\s*;\s*$`)

// refuses matches an action that may throw a parse error.
var refuses = regexp.MustCompile(`Throw\w*Exception\(|throw `)

// retypes matches a lexer action that makes its rule emit another token.
var retypes = regexp.MustCompile(`\$setType\(\s*(\w+)\s*\)`)

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
			p.note(line, "alternative dropped, its action refuses every input: %s", sourceText(p.toks[start:p.i]))
		} else {
			out = append(out, a)
		}
		if !p.is(tPunct, "|") {
			return out, nil
		}
		p.i++
	}
}

func (p *reader) note(line int, format string, args ...any) {
	n := fmt.Sprintf("%s:%d: %s", p.file, line, fmt.Sprintf(format, args...))
	if len(n) > 300 {
		n = n[:297] + "..."
	}
	p.rule.notes = append(p.rule.notes, n)
}

func (p *reader) alt() (a alt, refused bool, err error) {
	var actions []token
	var next []string // the words of a predicate before the next element
	var pred token
	for {
		t := p.peek()
		switch {
		case t.kind == tPunct && (t.text == "|" || t.text == ")" || t.text == ";"):
			if !p.rule.lexer {
				a = p.matchWords(a, actions)
			}
			return a, refused, nil
		case t.kind == tPunct && t.text == "<end>":
			return nil, false, fmt.Errorf("unexpected end of grammar")
		case t.kind == tIdent && t.text == "exception":
			p.skipHandlers()
			continue
		case t.kind == tAction:
			p.i++
			if p.is(tPunct, "?") {
				p.i++
				if words := p.nextWords(t.text); words != nil && !p.rule.lexer {
					next, pred = words, t
					continue
				}
				p.note(t.line, "predicate dropped: {%s}?", oneLine(t.text))
				continue
			}
			actions = append(actions, t)
			switch {
			case refusal.MatchString(t.text):
				refused = true
			case refuses.MatchString(t.text):
				p.note(t.line, "action may refuse the input: {%s}", oneLine(t.text))
			case p.rule.lexer && retypes.MatchString(t.text):
				for _, m := range retypes.FindAllStringSubmatch(t.text, -1) {
					p.note(t.line, "this alternative may emit %s instead ($setType)", m[1])
				}
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
		if next != nil {
			e.next, e.pred = next, pred
			next = nil
		}
		a = append(a, e)
	}
}

// predWord matches one test of a predicate on the next token's text,
// NextTokenMatches(CodeGenerationSupporter.X), with what it is and-ed with
// after it, which the conversion loses.
var predWord = regexp.MustCompile(`^\s*\(?\s*NextTokenMatches\(\s*CodeGenerationSupporter\.(\w+)\s*\)\s*\)?\s*(&&.*)?$`)

// nextWords returns the words a predicate requires the next token to
// spell, one of which it must: NextTokenMatches tests or-ed together. It
// returns nil for any other predicate.
func (p *reader) nextWords(code string) []string {
	if strings.Contains(strings.ReplaceAll(code, "!=", ""), "!") {
		return nil
	}
	var words []string
	for _, part := range strings.Split(code, "||") {
		m := predWord.FindStringSubmatch(part)
		if m == nil {
			return nil
		}
		w, ok := p.words[m[1]]
		if !ok || !wordShape.MatchString(w) {
			return nil
		}
		words = append(words, strings.ToUpper(w))
	}
	return words
}

// matchCall finds Match(label, CodeGenerationSupporter.X[, CodeGenerationSupporter.Y][, token]).
var matchCall = regexp.MustCompile(`\bMatch\(\s*(\w+)\s*,\s*CodeGenerationSupporter\.(\w+)\s*(?:,\s*CodeGenerationSupporter\.(\w+)\s*)?(?:,\s*\w+\s*)?\)`)

// wordShape is a word a T-SQL identifier spells.
var wordShape = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// matchWords writes, as the words themselves, the elements of an
// alternative whose text an action checks against words of
// CodeGenerationSupporter. A Match at an action's top level, which every
// parse through the alternative runs, gives the one or two words the element
// is. Otherwise the element is the words every Match and TryMatch of it
// names, under whatever condition: ScriptDOM reads a choice of words as
// if (TryMatch(t, A)) ... else Match(t, B). A branch that takes any other
// text, where there is one, is lost.
func (p *reader) matchWords(a alt, actions []token) alt {
	sure := map[string][]string{}  // label: the words of an unconditional Match
	maybe := map[string][]string{} // label: the words of every Match and TryMatch
	for _, act := range actions {
		// Every Match and TryMatch, in the order the action has them.
		type call struct {
			m   []int
			try bool
		}
		var calls []call
		for _, m := range matchCall.FindAllStringSubmatchIndex(act.text, -1) {
			calls = append(calls, call{m, false})
		}
		for _, m := range tryMatchCall.FindAllStringSubmatchIndex(act.text, -1) {
			calls = append(calls, call{m, true})
		}
		slices.SortFunc(calls, func(x, y call) int { return x.m[0] - y.m[0] })
		for _, c := range calls {
			label := act.text[c.m[2]:c.m[3]]
			idx := c.m[4:]
			if c.try {
				idx = c.m[4:6]
			}
			texts, ok := p.matchTexts(act.text, idx)
			if !ok {
				continue
			}
			if !c.try && unconditional(act.text, c.m[0]) && sure[label] == nil {
				sure[label] = texts
			}
			maybe[label] = appendNew(maybe[label], texts...)
		}
	}
	for k := range a {
		e := &a[k]
		if e.label == "" || e.ref == "" || e.op != "" || e.not {
			continue
		}
		texts := sure[e.label]
		if texts == nil {
			texts = maybe[e.label]
		}
		switch len(texts) {
		case 0:
		case 1:
			*e = elem{lit: texts[0], isLit: true}
		default:
			var g []alt
			for _, w := range texts {
				g = append(g, alt{{lit: w, isLit: true}})
			}
			*e = elem{group: g}
		}
	}
	return a
}

// tryMatchCall finds TryMatch(label, CodeGenerationSupporter.X).
var tryMatchCall = regexp.MustCompile(`\bTryMatch\(\s*(\w+)\s*,\s*CodeGenerationSupporter\.(\w+)\s*\)`)

// matchTexts returns the words of the CodeGenerationSupporter names at the
// submatch index pairs idx of code, capitals, skipping a pair that did not
// match; false when a name has no word a T-SQL identifier spells.
func (p *reader) matchTexts(code string, idx []int) ([]string, bool) {
	var texts []string
	for i := 0; i+1 < len(idx); i += 2 {
		if idx[i] < 0 {
			continue
		}
		w, found := p.words[code[idx[i]:idx[i+1]]]
		if !found || !wordShape.MatchString(w) {
			return nil, false
		}
		texts = append(texts, strings.ToUpper(w))
	}
	return texts, len(texts) > 0
}

// appendNew appends the words list does not hold yet.
func appendNew(list []string, words ...string) []string {
	for _, w := range words {
		if !slices.Contains(list, w) {
			list = append(list, w)
		}
	}
	return list
}

// unconditional reports whether the call at offset at of an action runs
// whenever the action does: at the action's top level, as a statement of its
// own, not under an if, else, loop or switch.
func unconditional(code string, at int) bool {
	depth := 0
	stmt := 0 // where the statement holding at starts
	for i := 0; i < at; i++ {
		switch code[i] {
		case '{':
			depth++
			stmt = i + 1
		case '}':
			depth--
			stmt = i + 1
		case ';':
			if depth == 0 {
				stmt = i + 1
			}
		case '"':
			for i++; i < at && code[i] != '"'; i++ {
				if code[i] == '\\' {
					i++
				}
			}
		case '/':
			if i+1 < at && code[i+1] == '/' {
				for i < at && code[i] != '\n' {
					i++
				}
				stmt = i
			}
		}
	}
	if depth != 0 {
		return false
	}
	return strings.TrimSpace(code[stmt:at]) == ""
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
	// A label, x:atom, or a return value's assignment, v=rule.
	if t := p.peek(); t.kind == tIdent {
		if next := p.at(1); next.kind == tPunct && (next.text == ":" || next.text == "=") {
			e.label = t.text
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
		// ( { init action } : ... ), which runs before the subrule.
		if p.peek().kind == tAction && p.at(1).kind == tPunct && p.at(1).text == ":" {
			p.i += 2
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

// A specializer writes the element after a predicate on the next token's
// text as what the predicate lets through: its first identifier spelled as
// the predicate's word. ScriptDOM reads most options as an Identifier and
// tells them apart in a predicate before the option's rule, as
// {NextTokenMatches(CodeGenerationSupporter.Edition)}? editionOption, so
// without this every option would be written as the sample identifier.
type specializer struct {
	byName  map[string]*rule
	made    map[string]*rule // rule_WORD, made for a predicate; nil where none can be
	pending map[string]bool  // being made, so a cycle fails rather than loops
	order   []*rule
}

// specialize applies every predicate kept on an element of rules, and
// returns the rules with the ones it made after them.
func specialize(rules []*rule) []*rule {
	s := &specializer{byName: map[string]*rule{}, made: map[string]*rule{}, pending: map[string]bool{}}
	for _, r := range rules {
		s.byName[r.name] = r
	}
	for _, r := range rules {
		for i := range r.alts {
			r.alts[i] = s.applyAlt(r, r.alts[i])
		}
	}
	for k := 0; k < len(s.order); k++ {
		r := s.order[k]
		for i := range r.alts {
			r.alts[i] = s.applyAlt(r, r.alts[i])
		}
	}
	return append(rules, s.order...)
}

func (s *specializer) applyAlt(r *rule, a alt) alt {
	for k := range a {
		e := &a[k]
		for g := range e.group {
			e.group[g] = s.applyAlt(r, e.group[g])
		}
		if e.next == nil {
			continue
		}
		words, pred := e.next, e.pred
		e.next = nil
		var group []alt
		for _, w := range words {
			plain := *e
			plain.op = ""
			se, ok := s.first(plain, w)
			if !ok {
				group = nil
				break
			}
			group = append(group, alt{se})
		}
		if group == nil {
			r.notes = append(r.notes, fmt.Sprintf("%s:%d: predicate dropped: {%s}?", r.file, pred.line, oneLine(pred.text)))
			continue
		}
		op := e.op
		if len(group) == 1 {
			*e = group[0][0]
		} else {
			*e = elem{group: group}
		}
		if e.op == "" {
			e.op = op
		} else if op != "" {
			*e = elem{group: []alt{{*e}}, op: op}
		}
		r.notes = append(r.notes, fmt.Sprintf("%s:%d: predicate written as its words: {%s}?", r.file, pred.line, oneLine(pred.text)))
	}
	return a
}

// first returns e, which has no operator, written to start with the word w
// where it would start with an identifier, and false where it cannot start
// with w.
func (s *specializer) first(e elem, w string) (elem, bool) {
	switch {
	case e.not || e.any || e.ranged:
		return e, false
	case e.isLit:
		return e, strings.EqualFold(e.lit, w)
	case e.ref == "Identifier":
		return elem{lit: w, isLit: true}, true
	case e.ref != "":
		r := s.byName[e.ref]
		if r == nil || r.lexer {
			return e, false
		}
		sr, ok := s.rule(r, w)
		if !ok {
			return e, false
		}
		return elem{ref: sr.name}, true
	case len(e.group) > 0:
		var g []alt
		for _, a := range e.group {
			if fa, ok := s.firstAlt(a, w); ok {
				g = append(g, fa)
			}
		}
		if len(g) == 0 {
			return e, false
		}
		return elem{group: g}, true
	}
	return e, false
}

// firstAlt returns a written to start with w, or false. Only a first
// element that must appear can be written so.
func (s *specializer) firstAlt(a alt, w string) (alt, bool) {
	if len(a) == 0 {
		return nil, false
	}
	head := a[0]
	switch head.op {
	case "":
	case "+":
		// The first repetition starts with w; the rest are as before.
		rest := head
		rest.op = "*"
		head.op = ""
		fe, ok := s.first(head, w)
		if !ok {
			return nil, false
		}
		return append(alt{fe, rest}, a[1:]...), true
	default:
		return nil, false
	}
	fe, ok := s.first(head, w)
	if !ok {
		return nil, false
	}
	return append(alt{fe}, a[1:]...), true
}

// rule returns r written to start with w, made once, named r_WORD.
func (s *specializer) rule(r *rule, w string) (*rule, bool) {
	name := r.name + "_" + w
	if sr, ok := s.made[name]; ok {
		return sr, sr != nil
	}
	if s.pending[name] {
		return nil, false
	}
	s.pending[name] = true
	var alts []alt
	for _, a := range r.alts {
		if fa, ok := s.firstAlt(cloneAlt(a), w); ok {
			alts = append(alts, fa)
		}
	}
	delete(s.pending, name)
	if len(alts) == 0 {
		s.made[name] = nil
		return nil, false
	}
	sr := &rule{file: r.file, line: r.line, name: name, alts: alts,
		notes: []string{fmt.Sprintf("%s where a predicate requires the next word to be %s", r.name, w)}}
	s.made[name] = sr
	s.order = append(s.order, sr)
	return sr, true
}

// reachable returns the rules the start rule reaches, in their order. A
// rule every reference to which a predicate specialized is reached only as
// its copies, and the entry points other than script (entryPointExpression
// and the like, which ScriptDOM's other Parse methods start from) are not
// reached at all.
func reachable(rules []*rule) []*rule {
	byName := map[string]*rule{}
	for _, r := range rules {
		byName[r.name] = r
	}
	seen := map[string]bool{}
	var visit func(a alt)
	var enter func(name string)
	enter = func(name string) {
		r := byName[name]
		if r == nil || seen[name] {
			return
		}
		seen[name] = true
		for _, a := range r.alts {
			visit(a)
		}
	}
	visit = func(a alt) {
		for _, e := range a {
			if e.ref != "" {
				enter(e.ref)
			}
			for _, g := range e.group {
				visit(g)
			}
		}
	}
	enter(startRule)
	var out []*rule
	for _, r := range rules {
		if seen[r.name] {
			out = append(out, r)
		}
	}
	return out
}

// cloneAlt copies a, so a rule made from another shares no group with it.
func cloneAlt(a alt) alt {
	out := make(alt, len(a))
	for i, e := range a {
		if e.group != nil {
			g := make([]alt, len(e.group))
			for j, ga := range e.group {
				g[j] = cloneAlt(ga)
			}
			e.group = g
		}
		out[i] = e
	}
	return out
}

// prune drops the rules left with no alternative, and with them every
// alternative that needs one, until none is left: an optional element that
// needs one is dropped alone. The note goes to the rule that referred.
func prune(rules []*rule) []*rule {
	for {
		gone := map[string]*rule{}
		var kept []*rule
		for _, r := range rules {
			if len(r.alts) == 0 {
				gone[r.name] = r
			} else {
				kept = append(kept, r)
			}
		}
		if len(gone) == 0 {
			return rules
		}
		for _, r := range kept {
			var alts []alt
			for _, a := range r.alts {
				if a, ok := pruneAlt(a, gone, r); ok {
					alts = append(alts, a)
				}
			}
			r.alts = alts
		}
		for _, g := range gone {
			for _, r := range kept {
				if r.refers[g.name] {
					r.notes = append(r.notes, g.notes...)
				}
			}
		}
		rules = kept
	}
}

// pruneAlt drops from a the elements that need a rule gone; ok is false
// when a needs one.
func pruneAlt(a alt, gone map[string]*rule, r *rule) (alt, bool) {
	var out alt
	for _, e := range a {
		ok := true
		switch {
		case e.ref != "" && gone[e.ref] != nil:
			ok = false
			if r.refers == nil {
				r.refers = map[string]bool{}
			}
			r.refers[e.ref] = true
		case len(e.group) > 0:
			var g []alt
			for _, ga := range e.group {
				if ga, ok := pruneAlt(ga, gone, r); ok {
					g = append(g, ga)
				}
			}
			e.group = g
			ok = len(g) > 0
		}
		if ok {
			out = append(out, e)
			continue
		}
		if e.op != "?" && e.op != "*" {
			return nil, false
		}
	}
	return out, true
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
	case e.any:
		s = "."
	case e.ranged:
		s = quote(e.lit) + ".." + quote(e.hi)
	case e.isLit:
		s = quote(e.lit)
	case e.ref != "":
		s = e.ref
	default:
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
