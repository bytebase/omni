package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// A terminal is a token the lexer section of googlesql.tm declares.
type terminal struct {
	name    string // the token's id, as KW_SELECT
	display string // its quoted name, as SELECT; "" when it has none
	spell   string // the text it matches, when the rule matches one text
}

// tokenLine reads one token declaration of the lexer section:
// "TEXT" (ID) {type}: /regex/, or ID {type}: /regex/, the regex optional.
var tokenLine = regexp.MustCompile(`^(?:"([^"]*)"\s*\(([A-Za-z_][A-Za-z_0-9]*)\)|([A-Za-z_][A-Za-z_0-9]*))\s*(?:\{[^}]*\})?\s*:\s*(/(?:\\.|[^/\\])*/)?`)

// readLexer returns the terminals the lexer section declares, by id and by
// quoted name. A lowercase rule (whitespace, eoi, invalid_token) is the
// lexer's own and no terminal.
func readLexer(src string) (byName, byDisplay map[string]*terminal) {
	byName, byDisplay = map[string]*terminal{}, map[string]*terminal{}
	src = blockComments.ReplaceAllStringFunc(src, func(c string) string {
		return strings.Repeat("\n", strings.Count(c, "\n"))
	})
	for _, line := range strings.Split(src, "\n") {
		m := tokenLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		t := &terminal{display: m[1], name: m[2]}
		if m[3] != "" {
			t.name = m[3]
			if t.name != strings.ToUpper(t.name) {
				continue
			}
		}
		if m[4] != "" {
			t.spell = literal(m[4][1 : len(m[4])-1])
		}
		if t.spell == "" && t.display != "" && !isWord(t.display) {
			// Punctuation the tokenizer makes from smaller tokens: "<>", ">>".
			t.spell = t.display
		}
		if _, dup := byName[t.name]; dup {
			continue
		}
		byName[t.name] = t
		if t.display != "" {
			byDisplay[t.display] = t
		}
	}
	return byName, byDisplay
}

var blockComments = regexp.MustCompile(`(?s)/\*.*?\*/`)

// literal returns the one text a token regex matches, upper-cased as
// googlesql.tm's caseInsensitive lexer reads keywords, or "" when it matches
// a class of texts.
func literal(re string) string {
	var b strings.Builder
	for i := 0; i < len(re); i++ {
		c := re[i]
		switch {
		case c == '\\' && i+1 < len(re):
			i++
			if strings.IndexByte(`.*+?()[]{}|/\^$-'"!`, re[i]) < 0 {
				return ""
			}
			b.WriteByte(re[i])
		case strings.IndexByte(`.*+?()[]{}|^$`, c) >= 0:
			return ""
		default:
			b.WriteByte(c)
		}
	}
	return strings.ToUpper(b.String())
}

func isWord(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c == '_' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
			return false
		}
	}
	return s != ""
}

// A ptok is a token of the parser section, with the line it starts on.
type ptok struct {
	text string
	line int
}

// lexParser splits the parser section into names, quoted names and
// punctuation, dropping comments (//, # and /* */) and the C++ of actions
// and rule types, which are brace blocks.
func lexParser(src string, firstLine int) ([]ptok, error) {
	var out []ptok
	line := firstLine
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == '\n':
			line++
			i++
		case c == ' ' || c == '\t' || c == '\r' || c == '\f':
			i++
		case strings.HasPrefix(src[i:], "//") || c == '#':
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				end = len(src) - i
			}
			i += end
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i+2:], "*/")
			if end < 0 {
				return nil, fmt.Errorf("googlesql.tm:%d: unterminated comment", line)
			}
			line += strings.Count(src[i:i+2+end+2], "\n")
			i += 2 + end + 2
		case c == '{':
			j, err := skipBlock(src, i)
			if err != nil {
				return nil, fmt.Errorf("googlesql.tm:%d: %v", line, err)
			}
			line += strings.Count(src[i:j], "\n")
			i = j
		case c == '"':
			end := strings.IndexByte(src[i+1:], '"')
			if end < 0 {
				return nil, fmt.Errorf("googlesql.tm:%d: unterminated quoted name", line)
			}
			out = append(out, ptok{src[i : i+1+end+1], line})
			i += end + 2
		case c == '%' || c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i + 1
			for j < len(src) && (src[j] == '_' || src[j] >= 'a' && src[j] <= 'z' || src[j] >= 'A' && src[j] <= 'Z' || src[j] >= '0' && src[j] <= '9') {
				j++
			}
			out = append(out, ptok{src[i:j], line})
			i = j
		default:
			out = append(out, ptok{string(c), line})
			i++
		}
	}
	return out, nil
}

// skipBlock returns the end of the brace block of C++ at src[i], skipping
// the string and character literals and comments inside it.
func skipBlock(src string, i int) (int, error) {
	depth := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '"' || c == '\'':
			j := i + 1
			for j < len(src) && src[j] != c {
				if src[j] == '\\' {
					j++
				}
				j++
			}
			i = j + 1
			continue
		case strings.HasPrefix(src[i:], "//"):
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				return 0, fmt.Errorf("unterminated action")
			}
			i += end
			continue
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i+2:], "*/")
			if end < 0 {
				return 0, fmt.Errorf("unterminated comment in an action")
			}
			i += end + 4
			continue
		case c == '{':
			depth++
		case c == '}':
			depth--
			if depth == 0 {
				return i + 1, nil
			}
		}
		i++
	}
	return 0, fmt.Errorf("unterminated action")
}

// An elem is one element of an alternative: a symbol, or a group of
// alternatives, under a quantifier ("", "?", "*" or "+").
type elem struct {
	sym  string
	alts [][]elem
	sep  []elem // the separator of a (x separator ",")+ list
	q    string
}

// A rule is one nonterminal: its name, the line it starts on, and its
// alternatives.
type rule struct {
	name string
	line int
	alts [][]elem
}

type parser struct {
	toks []ptok
	pos  int
}

func (p *parser) peek() string {
	if p.pos < len(p.toks) {
		return p.toks[p.pos].text
	}
	return ""
}

func (p *parser) line() int {
	if p.pos < len(p.toks) {
		return p.toks[p.pos].line
	}
	return 0
}

// readRules reads the rules of the parser section, skipping its directives:
// %input, precedence and %expect.
func readRules(toks []ptok) ([]rule, error) {
	p := &parser{toks: toks}
	var rules []rule
	for p.pos < len(p.toks) {
		t := p.peek()
		if strings.HasPrefix(t, "%") {
			for p.pos < len(p.toks) && p.peek() != ";" {
				p.pos++
			}
			p.pos++
			continue
		}
		if p.pos+1 >= len(p.toks) || p.toks[p.pos+1].text != ":" {
			return nil, fmt.Errorf("googlesql.tm:%d: want a rule, found %q", p.line(), t)
		}
		r := rule{name: t, line: p.line()}
		p.pos += 2
		alts, err := p.alternatives(";")
		if err != nil {
			return nil, err
		}
		p.pos++
		r.alts = alts
		rules = append(rules, r)
	}
	return rules, nil
}

// alternatives reads alternatives separated by | up to end, which it leaves.
func (p *parser) alternatives(end string) ([][]elem, error) {
	var alts [][]elem
	var cur []elem
	for {
		t := p.peek()
		switch {
		case t == "":
			return nil, fmt.Errorf("googlesql.tm: unexpected end of the parser section")
		case t == end:
			return append(alts, cur), nil
		case t == "|":
			alts = append(alts, cur)
			cur = nil
			p.pos++
		case t == "%empty":
			p.pos++
		case t == "%prec":
			p.pos += 2
		case t == "[":
			// An alias: x[name].
			for p.pos < len(p.toks) && p.peek() != "]" {
				p.pos++
			}
			p.pos++
		case t == "?" || t == "*" || t == "+":
			if len(cur) == 0 || cur[len(cur)-1].q != "" {
				return nil, fmt.Errorf("googlesql.tm:%d: quantifier %s with nothing before it", p.line(), t)
			}
			cur[len(cur)-1].q = t
			p.pos++
		case t == "(":
			p.pos++
			inner, err := p.alternatives(")")
			if err != nil {
				return nil, err
			}
			p.pos++
			e := elem{alts: inner}
			// (x separator ",")+: the last alternative ends with separator s.
			last := inner[len(inner)-1]
			for k, x := range last {
				if len(inner) == 1 && x.sym == "separator" && x.alts == nil && k > 0 && k < len(last)-1 {
					e.sep = last[k+1:]
					e.alts = [][]elem{last[:k]}
					break
				}
			}
			cur = append(cur, e)
		case strings.HasPrefix(t, "\"") || isWord(t):
			cur = append(cur, elem{sym: t})
			p.pos++
		default:
			return nil, fmt.Errorf("googlesql.tm:%d: unexpected %q", p.line(), t)
		}
	}
}

// g4Reserved are the words ANTLR 4 does not take as a rule name.
var g4Reserved = map[string]bool{
	"options": true, "tokens": true, "channels": true, "import": true, "fragment": true,
	"lexer": true, "parser": true, "grammar": true, "returns": true, "locals": true,
	"throws": true, "catch": true, "finally": true, "mode": true,
}

// A converter writes the rules as ANTLR 4.
type converter struct {
	byName, byDisplay map[string]*terminal
	rules             map[string]*rule
	// classes are the terminals written as a token name, which tokens.tsv
	// samples: a class of texts, or a token the tokenizer makes in code.
	classes map[string]bool
}

// symbol returns the ANTLR 4 text of one symbol of googlesql.tm.
func (c *converter) symbol(s string) (string, error) {
	if strings.HasPrefix(s, "\"") {
		d := s[1 : len(s)-1]
		t, ok := c.byDisplay[d]
		if !ok {
			return "", fmt.Errorf("the lexer declares no token %s", s)
		}
		return c.terminal(t), nil
	}
	if t, ok := c.byName[s]; ok {
		return c.terminal(t), nil
	}
	if _, ok := c.rules[s]; ok {
		return c.ruleName(s), nil
	}
	if base, ok := strings.CutSuffix(s, "_opt"); ok {
		if _, isRule := c.rules[base]; isRule {
			return c.ruleName(base) + "?", nil
		}
		if t, isTerm := c.byName[base]; isTerm {
			return c.terminal(t) + "?", nil
		}
	}
	return "", fmt.Errorf("unknown symbol %s", s)
}

func (c *converter) ruleName(s string) string {
	if g4Reserved[s] {
		return s + "_"
	}
	return s
}

// terminal writes a token as the text it spells, or by its id where it
// spells a class of texts or the tokenizer makes it in code.
func (c *converter) terminal(t *terminal) string {
	if t.spell != "" {
		return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(t.spell) + "'"
	}
	c.classes[t.name] = true
	return t.name
}

// alts writes alternatives, each on its own line after the first.
func (c *converter) alt(a []elem) (string, error) {
	var parts []string
	for _, e := range a {
		s, err := c.elem(e)
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, " "), nil
}

func (c *converter) elem(e elem) (string, error) {
	if e.alts == nil {
		s, err := c.symbol(e.sym)
		if err != nil {
			return "", err
		}
		if e.q != "" && strings.HasSuffix(s, "?") {
			s = "(" + s + ")"
		}
		return s + e.q, nil
	}
	var inner []string
	for _, a := range e.alts {
		s, err := c.alt(a)
		if err != nil {
			return "", err
		}
		inner = append(inner, s)
	}
	body := strings.Join(inner, " | ")
	if e.sep != nil {
		sep, err := c.alt(e.sep)
		if err != nil {
			return "", err
		}
		list := "(" + body + ") (" + sep + " (" + body + "))*"
		switch e.q {
		case "+":
			return "(" + list + ")", nil
		case "*", "?":
			return "(" + list + ")?", nil
		}
		return "", fmt.Errorf("a separator list under %q", e.q)
	}
	return "(" + body + ")" + e.q, nil
}

// convert writes googlesql.tm as an ANTLR 4 grammar whose first rule is
// start: each rule start reaches, under the line of googlesql.tm it came
// from. It returns the grammar and the class tokens it writes by name.
func convert(header, tm, start string) (string, []string, error) {
	lexStart := strings.Index(tm, "\n:: lexer")
	parseStart := strings.Index(tm, "\n:: parser")
	parseEnd := strings.Index(tm, "\n%%")
	if lexStart < 0 || parseStart < lexStart || parseEnd < parseStart {
		return "", nil, fmt.Errorf("googlesql.tm has no :: lexer, :: parser and %%%% sections")
	}
	byName, byDisplay := readLexer(tm[lexStart:parseStart])
	body := tm[parseStart+len("\n:: parser") : parseEnd]
	toks, err := lexParser(body, strings.Count(tm[:parseStart+1], "\n")+1)
	if err != nil {
		return "", nil, err
	}
	rules, err := readRules(toks)
	if err != nil {
		return "", nil, err
	}
	c := &converter{byName: byName, byDisplay: byDisplay, rules: map[string]*rule{}, classes: map[string]bool{}}
	for i := range rules {
		c.rules[rules[i].name] = &rules[i]
	}
	if _, ok := c.rules[start]; !ok {
		return "", nil, fmt.Errorf("googlesql.tm has no rule %s", start)
	}
	// The rules start reaches, in the order googlesql.tm writes them, start
	// first.
	reached := map[string]bool{start: true}
	queue := []string{start}
	for len(queue) > 0 {
		r := c.rules[queue[0]]
		queue = queue[1:]
		var walk func([][]elem)
		walk = func(alts [][]elem) {
			for _, a := range alts {
				for _, e := range a {
					if e.alts != nil {
						walk(e.alts)
						continue
					}
					name := e.sym
					if _, ok := c.rules[name]; !ok {
						base, cut := strings.CutSuffix(name, "_opt")
						if _, ok := c.rules[base]; !cut || !ok {
							continue
						}
						name = base
					}
					if !reached[name] {
						reached[name] = true
						queue = append(queue, name)
					}
				}
			}
		}
		walk(r.alts)
	}
	order := []*rule{c.rules[start]}
	for i := range rules {
		if rules[i].name != start && reached[rules[i].name] {
			order = append(order, &rules[i])
		}
	}
	var b strings.Builder
	b.WriteString(header)
	b.WriteString("grammar GoogleSQL;\n")
	for _, r := range order {
		fmt.Fprintf(&b, "\n// googlesql.tm:%d\n%s\n", r.line, c.ruleName(r.name))
		for i, a := range r.alts {
			s, err := c.alt(a)
			if err != nil {
				return "", nil, fmt.Errorf("googlesql.tm:%d: %s: %v", r.line, r.name, err)
			}
			lead := "    |"
			if i == 0 {
				lead = "    :"
			}
			if s == "" {
				s = "/* empty */"
			}
			fmt.Fprintf(&b, "%s %s\n", lead, s)
		}
		b.WriteString("    ;\n")
	}
	classes := make([]string, 0, len(c.classes))
	for n := range c.classes {
		classes = append(classes, n)
	}
	sort.Strings(classes)
	return b.String(), classes, nil
}
