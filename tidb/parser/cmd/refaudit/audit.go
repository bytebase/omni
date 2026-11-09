package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// A Rule is one left-hand side of the grammar and the line span it occupies in
// parser.y, 1-based and inclusive. A citation's line number resolves to the
// rule whose span contains it.
type Rule struct {
	Name  string
	Start int
	End   int
}

// ruleHead matches a goyacc rule head: a left-hand-side name alone on a column-0
// line, ending in a colon. Every one of v8.5.5's 690 rule heads has this shape;
// action code inside a rule body is tab-indented, so nothing else matches.
var ruleHead = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*):[ \t]*$`)

// parseGrammarRules returns the rules of parser.y in source order, each with the
// line span it covers, plus the inclusive [start, end] line bounds of the
// grammar section (between the two "%%" markers). Line numbers are 1-based so
// they match the parser.y:N citations verbatim.
func parseGrammarRules(src string) (rules []Rule, secStart, secEnd int, err error) {
	lines := strings.Split(src, "\n")
	var markers []int
	for i, line := range lines {
		if strings.TrimSpace(line) == "%%" {
			markers = append(markers, i+1) // 1-based
		}
	}
	if len(markers) < 2 {
		return nil, 0, 0, fmt.Errorf("want two %%%% section markers, found %d", len(markers))
	}
	secStart, secEnd = markers[0]+1, markers[1]-1

	type head struct {
		name string
		line int
	}
	var heads []head
	for ln := secStart; ln <= secEnd; ln++ {
		if m := ruleHead.FindStringSubmatch(lines[ln-1]); m != nil {
			heads = append(heads, head{name: m[1], line: ln})
		}
	}
	if len(heads) == 0 {
		return nil, 0, 0, fmt.Errorf("no rule heads found in the grammar section")
	}
	for i, h := range heads {
		end := secEnd
		if i+1 < len(heads) {
			end = heads[i+1].line - 1
		}
		rules = append(rules, Rule{Name: h.name, Start: h.line, End: end})
	}
	return rules, secStart, secEnd, nil
}

// A Citation is one parser.y line range cited by a Ref comment in omni's
// sources, with the file and line it was written on so a broken one can be
// found.
type Citation struct {
	File     string
	Line     int
	Start    int // first parser.y line of the range
	End      int // last; equals Start for a single-line citation
	RawRange string
}

// citationRange matches one "N" or "N-N" range after a "parser.y:" anchor.
var citationRange = regexp.MustCompile(`^(\d+)(?:-(\d+))?`)

// ruleWord matches a CamelCase identifier that could name a grammar rule, used
// only to credit coverage for citations that name a rule without a line number.
var ruleWord = regexp.MustCompile(`[A-Z][A-Za-z0-9_]*`)

const parserYAnchor = "parser.y:"

// commentBody returns the text after the "//" of a comment line, and whether
// the line is a comment at all. Leading indentation is ignored, so a doc comment
// and a tab-indented inline comment read the same.
func commentBody(line string) (string, bool) {
	t := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(t, "//") {
		return "", false
	}
	return strings.TrimPrefix(t, "//"), true
}

// extractCitations pulls the line-anchored parser.y citations from one source
// line, treating it as a one-line comment block. It is the single-line entry
// the unit tests use; scanFile reads real multi-line blocks.
func extractCitations(file string, lineNo int, text string) []Citation {
	if !strings.Contains(text, "Ref:") {
		return nil
	}
	cits, _ := extractBlock(file, []blockLine{{no: lineNo, text: text}})
	return cits
}

// A blockLine is one comment line of a block: its 1-based source line and the
// text after the "//".
type blockLine struct {
	no   int
	text string
}

// combine joins a block's comment bodies into one string with single-space
// separators, and returns a function mapping an offset in that string back to
// the source line it falls on. A Ref comment that wraps across several lines is
// then read as one run of text, so a parser.y: range, a continuation range, or a
// rule name on a later line is not lost, while a broken citation still reports
// the line it was written on.
func combine(block []blockLine) (string, func(int) int) {
	var b strings.Builder
	starts := make([]int, len(block))
	for i, bl := range block {
		if i > 0 {
			b.WriteByte(' ')
		}
		starts[i] = b.Len()
		b.WriteString(bl.text)
	}
	lineAt := func(off int) int {
		ln := block[0].no
		for i, s := range starts {
			if s <= off {
				ln = block[i].no
			} else {
				break
			}
		}
		return ln
	}
	return b.String(), lineAt
}

// extractBlock reads one comment block as a unit: every parser.y: line range in
// it (across continuation lines), each attributed to the line it appears on, and
// every CamelCase rule-like word it names. A range list under one anchor may be
// separated by commas or by "and", and each range may trail a "(annotation)".
func extractBlock(file string, block []blockLine) (cits []Citation, named []string) {
	s, lineAt := combine(block)
	named = ruleWord.FindAllString(s, -1)

	pos := 0
	for {
		i := strings.Index(s[pos:], parserYAnchor)
		if i < 0 {
			break
		}
		p := pos + i + len(parserYAnchor)
		for {
			for p < len(s) && s[p] == ' ' {
				p++
			}
			m := citationRange.FindStringSubmatch(s[p:])
			if m == nil {
				break
			}
			start, _ := strconv.Atoi(m[1])
			end := start
			if m[2] != "" {
				end, _ = strconv.Atoi(m[2])
			}
			cits = append(cits, Citation{File: file, Line: lineAt(p), Start: start, End: end, RawRange: m[0]})
			p += len(m[0])
			// Skip a trailing " (annotation)".
			for p < len(s) && s[p] == ' ' {
				p++
			}
			if p < len(s) && s[p] == '(' {
				if j := strings.IndexByte(s[p:], ')'); j >= 0 {
					p += j + 1
				} else {
					p = len(s)
				}
			}
			for p < len(s) && s[p] == ' ' {
				p++
			}
			// A comma or an "and" before another range continues the list.
			if p < len(s) && s[p] == ',' {
				p++
				continue
			}
			if rest := s[p:]; strings.HasPrefix(rest, "and ") {
				q := p + len("and ")
				for q < len(s) && s[q] == ' ' {
					q++
				}
				if citationRange.MatchString(s[q:]) {
					p = q
					continue
				}
			}
			break
		}
		pos = p
	}
	return cits, named
}

// scanFile returns the citations and named rules in one file's Ref comment
// blocks. Consecutive comment lines form a block; a block containing "Ref:" is
// read whole, so continuation lines take part in both validation and coverage.
func scanFile(rel, src string) (cits []Citation, named []string) {
	lines := strings.Split(src, "\n")
	for i := 0; i < len(lines); {
		if _, ok := commentBody(lines[i]); !ok {
			i++
			continue
		}
		var block []blockLine
		hasRef := false
		for i < len(lines) {
			body, ok := commentBody(lines[i])
			if !ok {
				break
			}
			if strings.Contains(body, "Ref:") {
				hasRef = true
			}
			block = append(block, blockLine{no: i + 1, text: body})
			i++
		}
		if !hasRef {
			continue
		}
		c, n := extractBlock(rel, block)
		cits = append(cits, c...)
		named = append(named, n...)
	}
	return cits, named
}

// scanTree walks the Go sources under root (excluding _test.go files, which are
// not productions) and returns every line-anchored parser.y citation and the
// set of rule-like names mentioned in any Ref comment block. Paths in the
// returned citations are relative to root.
func scanTree(root string) (cits []Citation, named map[string]bool, err error) {
	named = map[string]bool{}
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		c, n := scanFile(rel, string(b))
		cits = append(cits, c...)
		for _, name := range n {
			named[name] = true
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return cits, named, nil
}

// An Audit is the result of resolving every citation against the grammar: which
// rules are cited, which citations do not land on a real rule, and the rule set
// itself.
type Audit struct {
	Rules   []Rule
	SecS    int
	SecE    int
	Cited   map[string]bool // rule name -> cited by at least one production
	Broken  []BrokenCitation
	NumCits int
}

// A BrokenCitation is a citation whose line range does not resolve to a rule:
// it points outside the grammar section or into the gap between rules.
type BrokenCitation struct {
	Citation
	Reason string
}

// resolveLine returns the rule whose span contains the 1-based parser.y line, or
// nil when the line falls outside every rule.
func resolveLine(rules []Rule, line int) *Rule {
	// rules are sorted by Start; binary search for the last rule starting at or
	// before line.
	i := sort.Search(len(rules), func(i int) bool { return rules[i].Start > line }) - 1
	if i < 0 {
		return nil
	}
	if line <= rules[i].End {
		return &rules[i]
	}
	return nil
}

// runAudit resolves citations against the grammar rules. A line-range citation
// credits every rule its range intersects; one whose start line is outside the
// grammar section or lands between rules is broken. A rule named outright by any
// Ref comment is credited too (named), so a name-only Ref — common for the Opt
// sub-rules, which cite "AttributesOpt in ... parser.y." with no line — does not
// read as uncovered.
func runAudit(rules []Rule, cits []Citation, named map[string]bool) Audit {
	byName := make(map[string]*Rule, len(rules))
	for i := range rules {
		byName[rules[i].Name] = &rules[i]
	}
	a := Audit{Rules: rules, Cited: map[string]bool{}, NumCits: len(cits)}
	if len(rules) > 0 {
		a.SecS, a.SecE = rules[0].Start, rules[len(rules)-1].End
	}
	for _, c := range cits {
		if c.Start < a.SecS || c.End > a.SecE {
			a.Broken = append(a.Broken, BrokenCitation{Citation: c,
				Reason: fmt.Sprintf("line %s is outside the grammar section (%d-%d)", c.RawRange, a.SecS, a.SecE)})
			continue
		}
		r := resolveLine(rules, c.Start)
		if r == nil {
			a.Broken = append(a.Broken, BrokenCitation{Citation: c,
				Reason: fmt.Sprintf("line %s does not fall inside any rule", c.RawRange)})
			continue
		}
		// Credit every rule the range spans.
		for i := range rules {
			if rules[i].End >= c.Start && rules[i].Start <= c.End {
				a.Cited[rules[i].Name] = true
			}
		}
	}
	// Credit rules named outright by a Ref comment, whether or not it also gives
	// a line range.
	for n := range named {
		if _, ok := byName[n]; ok {
			a.Cited[n] = true
		}
	}
	return a
}

// uncitedRules returns the names of rules no production cites, in source order.
func (a Audit) uncitedRules() []string {
	var out []string
	for _, r := range a.Rules {
		if !a.Cited[r.Name] {
			out = append(out, r.Name)
		}
	}
	return out
}
