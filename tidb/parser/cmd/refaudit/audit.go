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
	File      string
	Line      int
	Start     int // first parser.y line of the range
	End       int // last; equals Start for a single-line citation
	RawRange  string
	RuleNames []string // CamelCase rule-like words named in the same comment
}

// citationRange matches one "N" or "N-N" range after a "parser.y:" anchor.
var citationRange = regexp.MustCompile(`^(\d+)(?:-(\d+))?`)

// ruleWord matches a CamelCase identifier that could name a grammar rule, used
// only to credit coverage for citations that name a rule without a line number.
var ruleWord = regexp.MustCompile(`[A-Z][A-Za-z0-9_]*`)

const parserYAnchor = "parser.y:"

// extractCitations pulls every line-anchored parser.y citation from one Go
// source line. A citation is a comment that mentions Ref: and carries a
// "parser.y:" anchor followed by one or more comma-separated N or N-N ranges,
// each optionally trailed by a "(annotation)". Multiple ranges under one anchor
// (parser.y:13764-13769, 2068-2084) are all returned; so are the rule-like
// words named anywhere in the comment, for coverage crediting.
func extractCitations(file string, lineNo int, text string) []Citation {
	if !strings.Contains(text, "Ref:") || !strings.Contains(text, parserYAnchor) {
		return nil
	}
	names := ruleWord.FindAllString(text, -1)

	var out []Citation
	rest := text
	for {
		i := strings.Index(rest, parserYAnchor)
		if i < 0 {
			break
		}
		rest = rest[i+len(parserYAnchor):]
		// Read a comma-separated list of ranges immediately after the anchor.
		for {
			rest = strings.TrimLeft(rest, " ")
			m := citationRange.FindStringSubmatch(rest)
			if m == nil {
				break
			}
			start, _ := strconv.Atoi(m[1])
			end := start
			if m[2] != "" {
				end, _ = strconv.Atoi(m[2])
			}
			out = append(out, Citation{
				File: file, Line: lineNo, Start: start, End: end, RawRange: m[0], RuleNames: names,
			})
			rest = rest[len(m[0]):]
			// Skip a trailing " (annotation)" before looking for a comma.
			rest = strings.TrimLeft(rest, " ")
			if strings.HasPrefix(rest, "(") {
				if j := strings.Index(rest, ")"); j >= 0 {
					rest = rest[j+1:]
				} else {
					rest = "" // unbalanced paren continues on the next line; stop
				}
			}
			rest = strings.TrimLeft(rest, " ")
			if strings.HasPrefix(rest, ",") {
				rest = rest[1:]
				continue
			}
			break
		}
	}
	return out
}

// scanTree walks the Go sources under root (excluding _test.go files, which are
// not productions) and returns every line-anchored parser.y citation in them.
// Paths in the returned citations are relative to root.
func scanTree(root string) ([]Citation, error) {
	var cits []Citation
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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
		for i, line := range strings.Split(string(b), "\n") {
			cits = append(cits, extractCitations(rel, i+1, line)...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cits, nil
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

// runAudit resolves citations against the grammar rules. A citation credits
// every rule its range intersects. A citation whose start line is outside the
// grammar section or lands between rules is broken. A rule named by a citation
// that carries no usable line range is credited too, so a name-only Ref does not
// read as an uncovered rule.
func runAudit(rules []Rule, cits []Citation) Audit {
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
	// Credit rules named outright by any Ref comment, line range or not.
	for _, c := range cits {
		for _, n := range c.RuleNames {
			if _, ok := byName[n]; ok {
				a.Cited[n] = true
			}
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
