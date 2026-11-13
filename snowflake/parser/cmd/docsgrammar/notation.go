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

var placeholderRE = regexp.MustCompile(`^<[A-Za-z_][A-Za-z0-9_ .-]*>`)

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
}

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
		switch t.kind {
		case nBar, nClose:
			return out, nil
		case nOpen:
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
	return out, nil
}

func resolveSeq(seq []*node) ([]*node, error) {
	out, err := resolveSeqEllipses(seq)
	if err != nil {
		return nil, err
	}
	return atLeastOneAfterSet(out), nil
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
		run := seq[i+1:]
		if len(run) == 1 && run[0].tok == nil && run[0].brace && run[0].op == 0 && len(run[0].alts) == 1 {
			// SET { [ A ] [ B ] }: the same, in braces.
			inner := run[0].alts[0]
			wrapped := atLeastOneAfterSet(append([]*node{n}, inner...))
			if len(wrapped) == 2 {
				return append(append([]*node(nil), seq[:i]...), wrapped...)
			}
			return seq
		}
		if len(run) == 0 {
			return seq
		}
		for _, r := range run {
			if r.tok != nil || r.brace || (r.op != '?' && r.op != '*') {
				return seq
			}
		}
		choice := &node{brace: true}
		for j, r := range run {
			req := *r
			if req.op == '?' {
				req.op = 0
				req.brace = true
			} else {
				req.op = '+'
			}
			alt := append([]*node{&req}, run[j+1:]...)
			choice.alts = append(choice.alts, alt)
		}
		return append(append([]*node(nil), seq[:i+1]...), choice)
	}
	return seq
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
				if comma {
					out = append(out, repeatAfterComma(prev))
				} else {
					out[len(out)-1] = repeat(prev, '+')
				}
				continue
			}
			alts, err := resolveGroup(n)
			if err != nil {
				return nil, err
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
		if len(body) >= 2 && isComma(body[len(body)-1]) {
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

// isRepeatable reports whether X ... can repeat X: a value or a group, not
// punctuation or a keyword, which ... after reads as text left out.
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
