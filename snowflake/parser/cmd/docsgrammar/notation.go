package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// The Syntax blocks of Snowflake's reference are near-BNF:
//
//	ALTER WAREHOUSE [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , <tag_name> = '<tag_value>' ... ]
//
// Uppercase words are keywords, <placeholders> are values the reader
// supplies, [ ] is optional, { a | b } is a choice, ... repeats, and a
// lowercase word names a fragment the page defines below with ::=. Other
// punctuation, ( ) included, is literal. This file reads that notation.

// A ntok is one token of the notation.
type ntok struct {
	kind ntkind
	text string
}

type ntkind int

const (
	nKeyword     ntkind = iota // an uppercase word: written as is
	nLiteral                   // punctuation, a number or a quoted string: written as is
	nPlaceholder               // <name>: text is the name
	nTerminal                  // a token class the notation spells otherwise: '<x>' is STRING
	nFragment                  // a lowercase word: a fragment reference, or an undefined word
	nOpen                      // [ or {
	nClose                     // ] or }
	nBar                       // |
	nEllipsis                  // ... or <...>
)

var (
	placeholderRE = regexp.MustCompile(`^<[A-Za-z_][A-Za-z0-9_ .-]*>`)
	// quotedChoice is a choice written inside quotes: '{ S3 | S3CHINA }',
	// 'LOGIN_NAME | EMAIL_ADDRESS'. Each choice is a string of its own.
	quotedChoice = regexp.MustCompile(`^\s*\{?\s*([A-Za-z0-9_-]+(?:\s*\|\s*[A-Za-z0-9_-]+)+)\s*\}?\s*$`)
	// keywordSuffix is a keyword's optional ending, DESC[RIBE].
	keywordSuffix = regexp.MustCompile(`^\[([A-Z]+)\]`)
)

// lexNotation reads notation text into tokens. Comments (-- to the end of a
// line, /* */) are dropped. A quoted string that holds a placeholder is a
// STRING; $$ ... $$ is a STRING whatever it holds.
func lexNotation(s string) ([]ntok, error) {
	var out []ntok
	rs := []rune(s)
	n := len(rs)
	for i := 0; i < n; {
		c := rs[i]
		rest := string(rs[i:])
		switch {
		case unicode.IsSpace(c):
			i++
		case strings.HasPrefix(rest, "file://"):
			// A local file, file://<path>/<file>: one token to its end, as
			// a stage location is, so no space is written inside it.
			j := i
			for j < n && !unicode.IsSpace(rs[j]) && !strings.ContainsRune(",)}|;", rs[j]) {
				j++
			}
			out = append(out, ntok{nTerminal, tFile})
			i = j
		case strings.HasPrefix(rest, "--"):
			for i < n && rs[i] != '\n' {
				i++
			}
		case strings.HasPrefix(rest, "/*"):
			j := strings.Index(rest[2:], "*/")
			if j < 0 {
				return nil, fmt.Errorf("unterminated /*")
			}
			i += len([]rune(rest[:j+4]))
		case strings.HasPrefix(rest, "$$"):
			j := strings.Index(rest[2:], "$$")
			if j < 0 {
				return nil, fmt.Errorf("unterminated $$")
			}
			out = append(out, ntok{nTerminal, tString})
			i += len([]rune(rest[:j+4]))
		case strings.HasPrefix(rest, "..."):
			out = append(out, ntok{nEllipsis, "..."})
			i += 3
		case strings.HasPrefix(rest, "<...>"):
			out = append(out, ntok{nEllipsis, "..."})
			i += 5
		case c == '\'' || c == '"':
			j := i + 1
			for j < n && rs[j] != c && rs[j] != '\n' {
				j++
			}
			if j >= n || rs[j] != c {
				return nil, fmt.Errorf("unterminated %c at %q", c, abbreviate(rest))
			}
			body := string(rs[i+1 : j])
			switch {
			case c == '"':
				// A quoted identifier: a name, whatever the docs show in it.
				out = append(out, ntok{nTerminal, tName})
			case quotedChoice.MatchString(body):
				out = append(out, ntok{nOpen, "{"})
				for k, v := range strings.Split(quotedChoice.FindStringSubmatch(body)[1], "|") {
					if k > 0 {
						out = append(out, ntok{nBar, "|"})
					}
					out = append(out, ntok{nLiteral, "'" + strings.TrimSpace(v) + "'"})
				}
				out = append(out, ntok{nClose, "}"})
			case strings.ContainsAny(body, "<>") || body == "":
				out = append(out, ntok{nTerminal, tString})
			default:
				out = append(out, ntok{nLiteral, string(rs[i : j+1])})
			}
			i = j + 1
		case c == '<' && placeholderRE.MatchString(rest):
			m := placeholderRE.FindString(rest)
			out = append(out, ntok{nPlaceholder, strings.Trim(m, "<>")})
			i += len([]rune(m))
		case c == '[' || c == '{':
			out = append(out, ntok{nOpen, string(c)})
			i++
		case c == ']' || c == '}':
			out = append(out, ntok{nClose, string(c)})
			i++
		case c == '|' && !strings.HasPrefix(rest, "||"):
			out = append(out, ntok{nBar, "|"})
			i++
		case c == '_' || unicode.IsLetter(c):
			j := i
			for j < n && (rs[j] == '_' || rs[j] == '$' || unicode.IsLetter(rs[j]) || unicode.IsDigit(rs[j])) {
				j++
			}
			w := string(rs[i:j])
			// SKIP_FILE_<num>, METADATA$<field>, RETURN_<n>_ROWS: a value
			// written inside a word is one word, spelled with a sample.
			for j < n && rs[j] == '<' && unicode.IsUpper([]rune(w)[0]) {
				m := placeholderRE.FindString(string(rs[j:]))
				if m == "" {
					break
				}
				sample := "a"
				if placeholderClass(strings.Trim(m, "<>")) == tNumber {
					sample = "1"
				}
				k := j + len([]rune(m))
				for k < n && (rs[k] == '_' || rs[k] == '$' || unicode.IsLetter(rs[k]) || unicode.IsDigit(rs[k])) {
					k++
				}
				w += sample + string(rs[j+len([]rune(m)):k])
				j = k
			}
			if m := keywordSuffix.FindString(string(rs[j:])); m != "" && unicode.IsUpper([]rune(w)[0]) {
				// DESC[RIBE] is DESC or DESCRIBE.
				full := w + keywordSuffix.FindStringSubmatch(m)[1]
				out = append(out, ntok{nOpen, "{"}, ntok{nKeyword, w}, ntok{nBar, "|"}, ntok{nKeyword, full}, ntok{nClose, "}"})
				i = j + len([]rune(m))
				continue
			}
			if unicode.IsUpper([]rune(w)[0]) {
				out = append(out, ntok{nKeyword, w})
			} else {
				out = append(out, ntok{nFragment, w})
			}
			i = j
		case unicode.IsDigit(c):
			j := i
			for j < n && (unicode.IsDigit(rs[j]) || rs[j] == '.') {
				j++
			}
			out = append(out, ntok{nLiteral, string(rs[i:j])})
			i = j
		case c == '@' && i+1 < n && strings.ContainsRune("<[%~", rs[i+1]):
			// A stage location, @<stage> or @[<namespace>.]<stage>[/<path>/],
			// @%<table> or @~: one token, through its brackets.
			j, depth := i+1, 0
			for j < n && (depth > 0 || !unicode.IsSpace(rs[j]) && !strings.ContainsRune(",)}|;", rs[j])) {
				switch rs[j] {
				case '[':
					depth++
				case ']':
					depth--
				}
				j++
			}
			out = append(out, ntok{nTerminal, tStage})
			i = j
		case c == '$' && i+1 < n && rs[i+1] == '<':
			// $<col_position>: a column by its position.
			m := placeholderRE.FindString(string(rs[i+1:]))
			if m == "" {
				return nil, fmt.Errorf("$ before %q", abbreviate(rest))
			}
			out = append(out, ntok{nTerminal, tPosition})
			i += 1 + len([]rune(m))
		default:
			op := string(c)
			for _, two := range []string{"=>", "->", "::", "||", "!=", "<=", ">=", "<>", "@~"} {
				if strings.HasPrefix(rest, two) {
					op = two
					break
				}
			}
			if !strings.ContainsRune("()[]{},.;:=*+-/%@~!<>|$&^?", c) {
				return nil, fmt.Errorf("unexpected %q", c)
			}
			out = append(out, ntok{nLiteral, op})
			i += len([]rune(op))
		}
	}
	return out, nil
}

func abbreviate(s string) string {
	if len(s) > 40 {
		return s[:40] + "..."
	}
	return s
}

// --- the tree ---

// A node is one item of a form: a token, or a group.
type node struct {
	tok   *ntok
	alts  [][]*node // a group's alternatives
	op    byte      // 0, '?', '*' or '+'
	brace bool      // a { } group, as opposed to [ ]
}

// parseNotation reads tokens into alternatives of sequences.
func parseNotation(toks []ntok) ([][]*node, error) {
	p := &nparser{toks: toks}
	alts, err := p.alts()
	if err != nil {
		return nil, err
	}
	if p.i < len(toks) {
		return nil, fmt.Errorf("unbalanced %q", toks[p.i].text)
	}
	return alts, nil
}

type nparser struct {
	toks []ntok
	i    int
	// parens counts the literal ( the parser is inside, where a ) ends a
	// sequence as ] and } do.
	parens int
}

func isLiteral(t ntok, text string) bool { return t.kind == nLiteral && t.text == text }

func (p *nparser) alts() ([][]*node, error) {
	var out [][]*node
	for {
		seq, err := p.seq()
		if err != nil {
			return nil, err
		}
		out = append(out, seq)
		if p.i < len(p.toks) && p.toks[p.i].kind == nBar {
			p.i++
			continue
		}
		return out, nil
	}
}

func (p *nparser) seq() ([]*node, error) {
	var out []*node
	for p.i < len(p.toks) {
		t := p.toks[p.i]
		switch {
		case t.kind == nBar || t.kind == nClose || p.parens > 0 && isLiteral(t, ")"):
			return out, nil
		case isLiteral(t, "("):
			// ( ... ) is SQL's own, but a | inside it chooses within it:
			// ENCRYPTION = ( TYPE = 'A' | TYPE = 'B' ). Where no ) closes
			// it at this level, as in [ ( ], it is a plain literal.
			start := p.i
			p.i++
			p.parens++
			alts, err := p.alts()
			p.parens--
			if err != nil || p.i >= len(p.toks) || !isLiteral(p.toks[p.i], ")") {
				p.i = start + 1
				tok := t
				out = append(out, &node{tok: &tok})
				continue
			}
			p.i++
			open, closeTok := t, p.toks[p.i-1]
			out = append(out, &node{tok: &open})
			if len(alts) == 1 {
				out = append(out, alts[0]...)
			} else {
				out = append(out, &node{alts: alts, brace: true})
			}
			out = append(out, &node{tok: &closeTok})
		case t.kind == nOpen:
			p.i++
			alts, err := p.alts()
			if err != nil {
				return nil, err
			}
			if p.i >= len(p.toks) || p.toks[p.i].kind != nClose {
				return nil, fmt.Errorf("unclosed %s", t.text)
			}
			want := map[string]string{"[": "]", "{": "}"}[t.text]
			if p.toks[p.i].text != want {
				return nil, fmt.Errorf("%s closed by %s", t.text, p.toks[p.i].text)
			}
			p.i++
			g := &node{alts: alts, brace: t.text == "{"}
			if !g.brace {
				g.op = '?'
			}
			out = append(out, g)
		default:
			tok := t
			out = append(out, &node{tok: &tok})
			p.i++
		}
	}
	return out, nil
}

// --- ellipses ---

// errElided marks a form that writes ... for text it leaves out, such as
// SELECT * ILIKE ... REPLACE ..., rather than for a repetition.
var errElided = fmt.Errorf("... stands for text the docs leave out")

func isEllipsis(n *node) bool { return n.tok != nil && n.tok.kind == nEllipsis }
func isComma(n *node) bool    { return n.tok != nil && n.tok.kind == nLiteral && n.tok.text == "," }

// onlyEllipsis reports whether a [ ] group holds just ... or just , ...
func onlyEllipsis(n *node) (comma, ok bool) {
	if n.tok != nil || n.brace || len(n.alts) != 1 {
		return false, false
	}
	seq := n.alts[0]
	switch {
	case len(seq) == 1 && isEllipsis(seq[0]):
		return false, true
	case len(seq) == 2 && isComma(seq[0]) && isEllipsis(seq[1]):
		return true, true
	}
	return false, false
}

// resolveEllipses rewrites every ... of a form into a repetition:
//
//	X [ ... ]          X+            X [ , ... ]       X ( , X )*
//	[ A ... ]          ( A )*        [ , A , ... ]     ( , A )*
//	X , ...            X ( , X )*    X ...             X+
//
// Any other ... stands for text left out, and the form is not converted.
func resolveEllipses(alts [][]*node) ([][]*node, error) {
	out := make([][]*node, len(alts))
	for i, seq := range alts {
		r, err := resolveSeq(seq)
		if err != nil {
			return nil, err
		}
		out[i] = r
	}
	return scopeValueChoice(out), nil
}

// scopeValueChoice reads [ KEY = A | B | C ] as KEY = one of A, B or C:
// where the first alternative assigns a property and the others are bare
// values, each other takes the first one's KEY =.
func scopeValueChoice(alts [][]*node) [][]*node {
	if len(alts) < 2 || len(alts[0]) < 3 {
		return alts
	}
	head := alts[0]
	if head[0].tok == nil || head[0].tok.kind != nKeyword || head[1].tok == nil || !isLiteral(*head[1].tok, "=") {
		return alts
	}
	for _, alt := range alts[1:] {
		if len(alt) == 0 {
			return alts
		}
		for _, n := range alt {
			if n.tok != nil && isLiteral(*n.tok, "=") {
				return alts
			}
		}
		if alt[0].tok != nil && alt[0].tok.text == head[0].tok.text {
			return alts
		}
	}
	out := [][]*node{head}
	for _, alt := range alts[1:] {
		out = append(out, append(append([]*node(nil), head[:2]...), alt...))
	}
	return out
}

func resolveSeq(seq []*node) ([]*node, error) {
	out, err := resolveSeqEllipses(seq)
	if err != nil {
		return nil, err
	}
	return atLeastOneAfterSet(commaOptions(unsetList(unsetNamesOnly(pairParens(out))))), nil
}

// pairParens reads [ ( ] X [ ) ] as ( X ) or X: SQL's parentheses come
// in pairs, though the docs mark each optional.
func pairParens(seq []*node) []*node {
	only := func(n *node, text string) bool {
		return n.tok == nil && !n.brace && n.op == '?' && len(n.alts) == 1 && len(n.alts[0]) == 1 &&
			n.alts[0][0].tok != nil && isLiteral(*n.alts[0][0].tok, text)
	}
	for i := range seq {
		if !only(seq[i], "(") {
			continue
		}
		for j := i + 1; j < len(seq); j++ {
			if !only(seq[j], ")") {
				continue
			}
			mid := seq[i+1 : j]
			open, closeTok := seq[i].alts[0][0], seq[j].alts[0][0]
			withParens := append(append([]*node{open}, mid...), closeTok)
			choice := &node{alts: [][]*node{withParens, append([]*node(nil), mid...)}, brace: true}
			out := append(append([]*node(nil), seq[:i]...), choice)
			return pairParens(append(out, seq[j+1:]...))
		}
	}
	return seq
}

// unsetList reads property names UNSET lists one per bracket,
//
//	UNSET [ PASSWORD_MIN_LENGTH ] [ PASSWORD_MAX_LENGTH ] ...
//
// as one or more of them with commas between, as UNSET takes them.
func unsetList(seq []*node) []*node {
	name := func(n *node) bool {
		return n.tok == nil && !n.brace && n.op == '?' && len(n.alts) == 1 && len(n.alts[0]) == 1 &&
			n.alts[0][0].tok != nil && n.alts[0][0].tok.kind == nKeyword
	}
	for i, n := range seq {
		if n.tok == nil || n.tok.kind != nKeyword || n.tok.text != "UNSET" {
			continue
		}
		j := i + 1
		for j < len(seq) && name(seq[j]) {
			j++
		}
		if j-i-1 < 2 {
			return seq
		}
		choice := &node{brace: true}
		for _, r := range seq[i+1 : j] {
			choice.alts = append(choice.alts, r.alts[0])
		}
		out := append(append([]*node(nil), seq[:i+1]...), choice, repeatAfterComma(choice))
		return append(out, seq[j:]...)
	}
	return seq
}

// unsetNamesOnly reads a property UNSET lists with a value,
//
//	ALTER VIEW <name> UNSET [ SECURE ] [ COMMENT = '<string_literal>' ]
//
// as its name alone: UNSET takes property names, and a value is an error.
func unsetNamesOnly(seq []*node) []*node {
	for i, n := range seq {
		if n.tok == nil || n.tok.kind != nKeyword || n.tok.text != "UNSET" {
			continue
		}
		out := append([]*node(nil), seq[:i+1]...)
		for _, r := range seq[i+1:] {
			if r.tok == nil && !r.brace && r.op == '?' && len(r.alts) == 1 && len(r.alts[0]) >= 3 {
				alt := r.alts[0]
				if alt[0].tok != nil && alt[0].tok.kind == nKeyword && alt[1].tok != nil && isLiteral(*alt[1].tok, "=") {
					c := *r
					c.alts = [][]*node{{alt[0]}}
					r = &c
				}
			}
			out = append(out, r)
		}
		return out
	}
	return seq
}

// commaOptions reads optional clauses separated by commas,
//
//	SET [ A = 1 ], [ B = 2 ], [ C = 3 ]
//
// as a list of one or more of them with a comma between each two, which
// is what the commas say: a comma may stand only between clauses that are
// there. A comma the docs leave out between two clauses changes nothing.
func commaOptions(seq []*node) []*node {
	isOpt := func(n *node) bool { return n.tok == nil && !n.brace && n.op == '?' }
	for i := 0; i < len(seq); i++ {
		if !isOpt(seq[i]) {
			continue
		}
		j, commas := i, 0
		var opts []*node
		for j < len(seq) {
			if isOpt(seq[j]) {
				opts = append(opts, seq[j])
				j++
				continue
			}
			if isComma(seq[j]) && j+1 < len(seq) && isOpt(seq[j+1]) {
				commas++
				j++
				continue
			}
			if isComma(seq[j]) && j+1 == len(seq) {
				j++ // a comma left after the last clause
			}
			break
		}
		if commas < 2 {
			continue
		}
		choice := &node{brace: true}
		for _, o := range opts {
			choice.alts = append(choice.alts, o.alts...)
		}
		out := append(append([]*node(nil), seq[:i]...), choice, repeatAfterComma(choice))
		return append(out, seq[j:]...)
	}
	return seq
}

// atLeastOneAfterSet reads SET or UNSET followed only by optional parts as
// setting at least one of them: the ALTER pages say SET sets one or more
// properties, and list each as optional.
//
//	SET [ A ] [ B ]   SET { A [ B ] | B }
func atLeastOneAfterSet(seq []*node) []*node {
	for i, n := range seq {
		if n.tok == nil || n.tok.kind != nKeyword || (n.tok.text != "SET" && n.tok.text != "UNSET") {
			continue
		}
		run := requireOne(seq[i+1:])
		return append(append([]*node(nil), seq[:i+1]...), run...)
	}
	return seq
}

// requireOne returns a sequence that derives what seq does, less the
// empty text: a run of optional parts becomes a choice of which comes
// first, and a choice in braces needs one of its branches non-empty.
//
//	[ A ] [ B ]          { A [ B ] | B }
//	{ [ A ] | [ B ] }    { A | B }
func requireOne(seq []*node) []*node {
	if len(seq) == 0 {
		return seq
	}
	if len(seq) == 1 && seq[0].tok == nil && seq[0].brace && seq[0].op == 0 {
		choice := &node{brace: true}
		for _, alt := range seq[0].alts {
			if len(alt) == 0 {
				continue
			}
			choice.alts = append(choice.alts, requireOne(alt))
		}
		return []*node{choice}
	}
	for _, r := range seq {
		if r.tok != nil || r.brace || (r.op != '?' && r.op != '*') {
			return seq
		}
	}
	choice := &node{brace: true}
	for j, r := range seq {
		req := *r
		if req.op == '?' {
			req.op = 0
			req.brace = true
		} else {
			req.op = '+'
		}
		alt := append([]*node{&req}, seq[j+1:]...)
		choice.alts = append(choice.alts, alt)
	}
	return []*node{choice}
}

func resolveSeqEllipses(seq []*node) ([]*node, error) {
	var out []*node
	for _, n := range seq {
		if n.tok == nil {
			if comma, ok := onlyEllipsis(n); ok {
				if len(out) == 0 || isComma(out[len(out)-1]) {
					return nil, errElided
				}
				prev := out[len(out)-1]
				switch {
				case comma:
					// X [ , ... ] repeats the list element before it (see
					// listElement). An element of optional parts only,
					// ( [ <arg> <type> ] [ , ... ] ), makes the list optional,
					// with an element after every comma.
					k, required := listElement(out)
					if k == 0 {
						out = append(out, repeatAfterComma(prev))
						break
					}
					elem := append([]*node(nil), out[len(out)-k:]...)
					out = out[:len(out)-k]
					if !required {
						first := *elem[0]
						first.op, first.brace = 0, true
						elem[0] = &first
					}
					x := elem[0]
					if len(elem) > 1 {
						x = &node{alts: [][]*node{elem}, brace: true}
					}
					list := []*node{x, repeatAfterComma(x)}
					if required {
						out = append(out, list...)
					} else {
						out = append(out, &node{alts: [][]*node{list}, op: '?'})
					}
				default:
					out[len(out)-1] = repeat(prev, '+')
				}
				continue
			}
			if seq := n.alts; !n.brace && len(seq) == 1 && len(seq[0]) >= 3 && isComma(seq[0][0]) {
				if _, only := onlyEllipsis(seq[0][len(seq[0])-1]); only {
					// [ , <col> <type> ... [ ... ] ]: the trailing [ ... ]
					// repeats the whole comma-led group.
					body, err := resolveSeq(seq[0][:len(seq[0])-1])
					if err != nil {
						return nil, err
					}
					out = append(out, &node{alts: [][]*node{body}, op: '*'})
					continue
				}
			}
			alts, err := resolveGroup(n)
			if err != nil {
				return nil, err
			}
			if n.op == '?' && optionalList(n) {
				// [ X , ... ]: ( X ( , X )* )?
				x := &node{alts: alts, brace: true}
				out = append(out, &node{alts: [][]*node{{x, repeatAfterComma(x)}}, op: '?'})
				continue
			}
			out = append(out, &node{alts: alts, op: n.op, brace: n.brace})
			if n.op == '?' && len(alts) == 1 && trailingRepeat(n) {
				out[len(out)-1].op = '*'
			}
			continue
		}
		if !isEllipsis(n) {
			out = append(out, n)
			continue
		}
		switch {
		case len(out) >= 2 && isComma(out[len(out)-1]) && !isComma(out[len(out)-2]) && !isEllipsis(out[len(out)-2]):
			// X , ..., or E , E , ... where the docs spell a list element
			// twice, as <col> AS <alias>, <col> AS <alias>, ... does: a list
			// of E. Without a second spelling, X is the one node before.
			body := out[:len(out)-1]
			if k := repeatedElement(body); k > 0 {
				e := &node{alts: [][]*node{append([]*node(nil), body[len(body)-k:]...)}, brace: true}
				out = append(append([]*node(nil), body[:len(body)-2*k-1]...), e, repeatAfterComma(e))
				break
			}
			x := body[len(body)-1]
			out = append(body, repeatAfterComma(x))
		case len(out) >= 1 && !isComma(out[len(out)-1]) && isRepeatable(out[len(out)-1]):
			// X ...
			out[len(out)-1] = repeat(out[len(out)-1], '+')
		default:
			return nil, errElided
		}
	}
	return out, nil
}

// trailingRepeat reports whether a [ ] group's one sequence ended in ...,
// which resolveGroup dropped: the group repeats.
func trailingRepeat(n *node) bool {
	seq := n.alts[0]
	return len(seq) >= 2 && isEllipsis(seq[len(seq)-1])
}

// resolveGroup resolves a group's alternatives. A [ ] group whose one
// sequence ends in ..., after something other than the start, repeats as a
// whole: the ... and a comma before it are dropped here, and the caller
// makes the group *.
func resolveGroup(n *node) ([][]*node, error) {
	if !n.brace && len(n.alts) == 1 && trailingRepeat(n) {
		seq := n.alts[0]
		body := seq[:len(seq)-1]
		comma := len(body) >= 2 && isComma(body[len(body)-1])
		if comma {
			body = body[:len(body)-1]
		}
		if len(body) == 0 || (len(body) == 1 && isComma(body[0])) {
			return nil, errElided
		}
		r, err := resolveSeq(body)
		if err != nil {
			return nil, err
		}
		return [][]*node{r}, nil
	}
	return resolveEllipses(n.alts)
}

// optionalList reports whether a [ ] group is a list of one or more
// elements with commas between, [ X , ... ], which is optional as a whole
// rather than repeated as written.
func optionalList(n *node) bool {
	if n.brace || len(n.alts) != 1 || !trailingRepeat(n) {
		return false
	}
	seq := n.alts[0]
	return len(seq) >= 3 && isComma(seq[len(seq)-2]) && !isComma(seq[0])
}

// isRepeatable reports whether X ... can repeat X: a value or a group, not
// punctuation or a keyword, which ... after reads as text left out.
// listElement returns how many nodes at the end of seq are the list
// element a [ , ... ] after them repeats, and whether it has a required
// part. The element is any optional parts at the end, the required part
// before them (a value, a choice in braces, or a parenthesised run), and
// the values before that, back to a keyword, a comma, an unmatched ( or
// another group:
//
//	{ <expr> | <col> } [ [ AS ] <alias> ]    <col_name> = <value>
//	<arg_name> <arg_data_type>               ( <v> )
func listElement(seq []*node) (int, bool) {
	i := len(seq) - 1
	for i >= 0 && seq[i].tok == nil && !seq[i].brace && (seq[i].op == '?' || seq[i].op == '*') {
		i--
	}
	optional := len(seq) - 1 - i
	if i >= 0 && seq[i].tok == nil && seq[i].brace && seq[i].op == 0 {
		// A choice in braces is the required part; what comes before it is
		// not the element's.
		return optional + 1, true
	}
	if i >= 0 && seq[i].tok != nil && seq[i].tok.kind == nKeyword {
		// READ [ , WRITE ] [ , ... ]: a keyword is the required part.
		return optional + 1, true
	}
	depth, k := 0, 0
	for ; i >= 0; i-- {
		n := seq[i]
		if n.tok != nil && isLiteral(*n.tok, ")") {
			depth++
		} else if n.tok != nil && isLiteral(*n.tok, "(") {
			if depth == 0 {
				break
			}
			depth--
		} else if depth == 0 && (n.tok == nil || n.tok.kind == nKeyword || n.tok.kind == nFragment || isComma(n)) {
			break
		}
		k++
	}
	if depth > 0 {
		k = 0
	}
	return optional + k, k > 0
}

// repeatedElement returns the length of the list element seq ends with
// when the docs spell it twice, E , E, the longest such E, or 0.
func repeatedElement(seq []*node) int {
	n := len(seq)
	for k := (n - 1) / 2; k >= 1; k-- {
		if isComma(seq[n-k-1]) && nodesText(seq[n-2*k-1:n-k-1]) == nodesText(seq[n-k:]) {
			return k
		}
	}
	return 0
}

// nodesText writes nodes so that two runs are equal in it exactly when
// they are written alike.
func nodesText(ns []*node) string {
	var b strings.Builder
	for _, n := range ns {
		if n.tok != nil {
			fmt.Fprintf(&b, "\x00%d:%s", n.tok.kind, n.tok.text)
			continue
		}
		fmt.Fprintf(&b, "\x00%c%v(", n.op, n.brace)
		for _, alt := range n.alts {
			b.WriteString(nodesText(alt) + "|")
		}
		b.WriteString(")")
	}
	return b.String()
}

func isRepeatable(n *node) bool {
	if n.tok == nil {
		return true
	}
	switch n.tok.kind {
	case nPlaceholder, nTerminal, nFragment:
		return true
	}
	return false
}

func repeat(n *node, op byte) *node {
	c := *n
	if c.op == '?' || c.op == '*' {
		c.op = '*'
		return &c
	}
	if c.op == 0 {
		c.op = op
		return &c
	}
	return &node{alts: [][]*node{{n}}, op: op, brace: true}
}

func repeatAfterComma(x *node) *node {
	comma := &node{tok: &ntok{nLiteral, ","}}
	return &node{alts: [][]*node{{comma, x}}, op: '*'}
}
