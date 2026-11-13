package grammar

import (
	"strconv"
	"strings"
	"unicode"
)

// maxMatchDepth bounds how deep lexerMatches follows lexer rules into one
// another; a lexer rule that recurses deeper is taken not to match.
const maxMatchDepth = 64

// lexerMatches reports whether the lexer rule name matches all of text, as
// its rules are written: literals, sets, ranges, ~, ., groups, the rules it
// names, and ?, * and +. Actions and predicates are not run.
func (p *antlrParser) lexerMatches(name, text string) bool {
	alts, ok := p.lexer[name]
	if !ok {
		return false
	}
	s := []rune(text)
	for _, end := range p.matchAlts(alts, s, 0, 0) {
		if end == len(s) {
			return true
		}
	}
	return false
}

// matchAlts returns every position one of alts can end at, begun at pos.
func (p *antlrParser) matchAlts(alts []lexAlt, s []rune, pos, depth int) []int {
	if depth > maxMatchDepth {
		return nil
	}
	seen := map[int]bool{}
	var out []int
	for _, alt := range alts {
		at := []int{pos}
		for _, e := range alt {
			at = p.matchRepeated(e, s, at, depth)
			if len(at) == 0 {
				break
			}
		}
		for _, end := range at {
			if !seen[end] {
				seen[end] = true
				out = append(out, end)
			}
		}
	}
	return out
}

// matchRepeated applies e, with its operator, from every position in at.
func (p *antlrParser) matchRepeated(e lexElem, s []rune, at []int, depth int) []int {
	once := func(from []int) []int {
		seen := map[int]bool{}
		var out []int
		for _, pos := range from {
			for _, end := range p.matchOne(e, s, pos, depth) {
				if !seen[end] {
					seen[end] = true
					out = append(out, end)
				}
			}
		}
		return out
	}
	switch e.op {
	case '?':
		return union(at, once(at))
	case '*', '+':
		all := map[int]bool{}
		var out []int
		if e.op == '*' {
			out = union(nil, at)
			for _, pos := range at {
				all[pos] = true
			}
		}
		for frontier := once(at); len(frontier) > 0; {
			var fresh []int
			for _, pos := range frontier {
				if !all[pos] {
					all[pos] = true
					out = append(out, pos)
					fresh = append(fresh, pos)
				}
			}
			frontier = once(fresh)
		}
		return out
	}
	return once(at)
}

func union(a, b []int) []int {
	seen := map[int]bool{}
	var out []int
	for _, x := range append(append([]int(nil), a...), b...) {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

// matchOne returns where e, without its operator, can end, begun at pos.
func (p *antlrParser) matchOne(e lexElem, s []rune, pos, depth int) []int {
	switch {
	case e.neg != nil:
		if pos < len(s) && !p.matchesChar(*e.neg, s[pos], depth) {
			return []int{pos + 1}
		}
	case e.any:
		if pos < len(s) {
			return []int{pos + 1}
		}
	case e.ref != "":
		if alts, ok := p.lexer[e.ref]; ok {
			return p.matchAlts(alts, s, pos, depth+1)
		}
	case e.group != nil:
		return p.matchAlts(e.group, s, pos, depth+1)
	case e.ranged[0] != "" || e.class != "":
		if pos < len(s) && p.matchesChar(e, s[pos], depth) {
			return []int{pos + 1}
		}
	default:
		lit := []rune(e.literal)
		if pos+len(lit) > len(s) {
			return nil
		}
		for i, c := range lit {
			if !p.sameRune(c, s[pos+i]) {
				return nil
			}
		}
		return []int{pos + len(lit)}
	}
	return nil
}

// matchesChar reports whether e, a one-character element, matches c.
func (p *antlrParser) matchesChar(e lexElem, c rune, depth int) bool {
	switch {
	case e.class != "":
		return p.inClass(e.class, c)
	case e.ranged[0] != "":
		lo, hi := []rune(e.ranged[0]), []rune(e.ranged[1])
		if len(lo) != 1 || len(hi) != 1 {
			return false
		}
		return p.inRange(lo[0], hi[0], c)
	}
	for _, end := range p.matchOne(lexElem{literal: e.literal, ref: e.ref, group: e.group}, []rune{c}, 0, depth) {
		if end == 1 {
			return true
		}
	}
	return false
}

func (p *antlrParser) sameRune(a, b rune) bool {
	return a == b || p.caseInsensitive && unicode.ToLower(a) == unicode.ToLower(b)
}

func (p *antlrParser) inRange(lo, hi, c rune) bool {
	if lo <= c && c <= hi {
		return true
	}
	if p.caseInsensitive {
		l, u := unicode.ToLower(c), unicode.ToUpper(c)
		return lo <= l && l <= hi || lo <= u && u <= hi
	}
	return false
}

// inClass reports whether c is in a [...] set, written as the grammar has
// it: characters, ranges a-z, and escapes. A Unicode property, \p{...},
// is taken to hold every letter.
func (p *antlrParser) inClass(class string, c rune) bool {
	r := []rune(class)
	read := func(i int) (rune, int, bool) {
		if r[i] != '\\' || i+1 >= len(r) {
			return r[i], i + 1, false
		}
		switch e := r[i+1]; e {
		case 'n':
			return '\n', i + 2, false
		case 't':
			return '\t', i + 2, false
		case 'r':
			return '\r', i + 2, false
		case 'b':
			return '\b', i + 2, false
		case 'f':
			return '\f', i + 2, false
		case 'u':
			if i+6 <= len(r) {
				if v, err := strconv.ParseUint(string(r[i+2:i+6]), 16, 32); err == nil {
					return rune(v), i + 6, false
				}
			}
			return 'u', i + 2, false
		case 'p', 'P':
			end := strings.IndexRune(string(r[i:]), '}')
			if end < 0 {
				return 0, len(r), true
			}
			return 0, i + len([]rune(string(r[i:])[:end+1])), true
		default:
			return e, i + 2, false
		}
	}
	for i := 0; i < len(r); {
		lo, next, property := read(i)
		if property {
			if unicode.IsLetter(c) {
				return true
			}
			i = next
			continue
		}
		hi := lo
		if next+1 < len(r) && r[next] == '-' {
			hi, next, _ = read(next + 1)
		}
		if p.inRange(lo, hi, c) {
			return true
		}
		i = next
	}
	return false
}
