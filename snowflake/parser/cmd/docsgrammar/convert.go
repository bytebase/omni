package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Token classes a placeholder stands for. tokens.tsv samples each one.
const (
	tName        = "ID"
	tString      = "STRING"
	tNumber      = "NUMBER"
	tExpr        = "EXPR"
	tQuery       = "QUERY"
	tType        = "DATA_TYPE"
	tStage       = "STAGE"
	tPosition    = "POSITION"
	tObjectType  = "OBJECT_TYPE"
	tObjectTypes = "OBJECT_TYPES"
	tFile        = "FILE"
	tCopy        = "COPY_STATEMENT"
)

var (
	exprName       = regexp.MustCompile(`(^|_)(expr|exprs|expression|condition|predicate)(_?[0-9]+|_|$)`)
	numberName     = regexp.MustCompile(`(^|_)(num|number|integer|int|n|rows|seconds|minutes|hours|days|bytes|count|percent|ptus|buckets|interval|ms|threshold)(_?[0-9]+|_|$)|^(initial|start|min|max|step)_value$|^increment$`)
	queryName      = regexp.MustCompile(`^(query|subquery|select_statement|statement|dml_statement|sql|sql_statement)$`)
	typeName       = regexp.MustCompile(`(^|_)(data_type|col[0-9]*_type|arg_type|structured_type)(_?[0-9]+|_|$)|^type$|^arg_data_type`)
	runtimeVersion = regexp.MustCompile(`(python|java|scala|jdk|runtime)[a-z_]*_version$|^runtime_version$`)
	stringName     = regexp.MustCompile(`(^|_|-)(string|literal|text)(_?[0-9]+|_|-|$)`)
)

// placeholderClass returns the token class a <placeholder> stands for, by
// its name. Anything not recognised as an expression, number, query, data
// type or string is a name.
func placeholderClass(name string) string {
	n := strings.ToLower(name)
	switch {
	case n == "object_type":
		return tObjectType
	case n == "object_type_plural":
		return tObjectTypes
	case n == "copy_statement":
		return tCopy
	case n == "timestamp" || n == "function_definition" || n == "procedure_definition" || n == "yaml_definition":
		// START_TIMESTAMP = '2026-11-13'; AS '<body>' or $$<body>$$.
		return tString
	case runtimeVersion.MatchString(n):
		// RUNTIME_VERSION = '3.10': a runtime's version is a string.
		return tString
	case queryName.MatchString(n):
		return tQuery
	case typeName.MatchString(n):
		return tType
	case exprName.MatchString(n):
		return tExpr
	case stringName.MatchString(n):
		return tString
	case strings.Contains(n, "name") || strings.Contains(n, "identifier"):
		return tName
	case numberName.MatchString(n):
		return tNumber
	}
	return tName
}

// A fix is one hand correction of a page's syntax text, from
// docs-fixes.tsv: the docs are prose, and some blocks do not read as
// notation as written.
type fix struct {
	page, find, replace, reason string
	used                        bool
}

// A unit is one statement form or one ::= definition of a page.
type unit struct {
	def string // the defined fragment, or "" for a statement form
	// variant is set for a definition written for one variant, as
	// restAuthenticationParams (for Bearer token) ::= is.
	variant bool
	lines   []codeLine // the text, without the definition's head
}

type codeLine struct {
	text    string
	comment bool // a line that is only a -- comment
}

var defHead = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*(\([^)]*\))?\s*::=\s*`)

// statementWords returns the words a statement form on a page can start
// with: its slug's first word, in each spelling Snowflake takes, WITH, and
// the DDL verbs, which a page about one statement may show another of (the
// constraint page's CREATE TABLE and ALTER TABLE). None of them starts a
// clause.
func statementWords(slug string) map[string]bool {
	first, _, _ := strings.Cut(slug, "-")
	words := map[string]bool{strings.ToUpper(first): true, "WITH": true}
	for _, w := range []string{"CREATE", "ALTER", "DROP", "UNDROP", "SHOW", "DESC", "DESCRIBE", "GRANT", "REVOKE"} {
		words[w] = true
	}
	for _, same := range [][]string{{"DESC", "DESCRIBE"}, {"LS", "LIST"}, {"RM", "REMOVE"}, {"BEGIN", "START"}} {
		for _, w := range same {
			if words[w] {
				for _, v := range same {
					words[v] = true
				}
				break
			}
		}
	}
	return words
}

// splitUnits splits a page's blocks into statement forms and definitions.
// A form starts on a line outside any bracket that is not indented and
// starts with one of the page's statement words, or that starts with the
// first two words of the form before it, as ALTER TABLE's indented forms
// do; a definition starts at name ::= and runs to the next unit. Any other
// line continues the unit before it: COPY INTO writes its clauses
// unindented.
func splitUnits(blocks []string, starts map[string]bool) []unit {
	var out []unit
	for _, blk := range blocks {
		depth := 0
		var cur *unit
		flush := func() {
			if cur != nil && hasCode(cur.lines) {
				out = append(out, *cur)
			}
			cur = nil
		}
		for _, line := range strings.Split(blk, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "/*") && strings.HasSuffix(trimmed, "*/") {
				if cur != nil {
					cur.lines = append(cur.lines, codeLine{text: trimmed, comment: true})
				}
				continue
			}
			if m := defHead.FindStringSubmatch(line); m != nil {
				flush()
				depth = 0
				cur = &unit{def: m[1], variant: m[2] != ""}
				line = line[len(m[0]):]
				if strings.TrimSpace(line) == "" {
					continue
				}
			} else if depth <= 0 && (cur == nil || cur.def == "" &&
				(!unicode.IsSpace([]rune(line)[0]) && startsForm(trimmed, starts) || repeatsHead(trimmed, cur.lines))) &&
				!strings.ContainsAny(trimmed[:1], "|)]}") {
				flush()
				depth = 0
				cur = &unit{}
			}
			if cur == nil {
				cur = &unit{}
			}
			cur.lines = append(cur.lines, codeLine{text: line})
			depth += bracketDepth(line)
		}
		flush()
	}
	return out
}

// repeatsHead reports whether a line starts with the first two words of
// a unit's first line.
func repeatsHead(line string, lines []codeLine) bool {
	for _, l := range lines {
		if l.comment {
			continue
		}
		head := strings.Fields(l.text)
		if len(head) < 2 {
			return false
		}
		words := strings.Fields(line)
		return len(words) >= 2 && words[0] == head[0] && words[1] == head[1]
	}
	return false
}

// startsForm reports whether a line starts with one of a page's statement
// words.
func startsForm(line string, starts map[string]bool) bool {
	word := strings.FieldsFunc(line, func(r rune) bool { return !unicode.IsLetter(r) && r != '_' })
	if len(word) == 0 || !starts[word[0]] || !strings.HasPrefix(line, word[0]) {
		return false
	}
	// COMMENT = '...' sets a property; it is not a statement.
	return !strings.HasPrefix(strings.TrimSpace(line[len(word[0]):]), "=")
}

func hasCode(lines []codeLine) bool {
	for _, l := range lines {
		if !l.comment {
			return true
		}
	}
	return false
}

// bracketDepth returns how far a line opens brackets, outside quotes and
// comments.
func bracketDepth(line string) int {
	d := 0
	var quote rune
	rs := []rune(line)
	for i := 0; i < len(rs); i++ {
		c := rs[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '-' && i+1 < len(rs) && rs[i+1] == '-':
			return d
		case c == '[' || c == '{' || c == '(':
			d++
		case c == ']' || c == '}' || c == ')':
			d--
		}
	}
	return d
}

// A pageGrammar is what one page converts to.
type pageGrammar struct {
	slug, url, rule string
	fixes           []string
	forms           []string // converted alternatives of rule, as ANTLR text
	frags           []fragRule
	failed          []failure // forms and definitions not converted
}

type fragRule struct {
	name string
	alts []string
}

type failure struct {
	what, text, reason string
}

var errUndefined = errors.New("undefined")

// valueChoiceLine is a line that is KEY = A | B, with no bracket on it.
var valueChoiceLine = regexp.MustCompile(`(?m)^(\s*)([A-Z][A-Z0-9_]*\s*=\s*)([^\[\]{}\n]*\|[^\[\]{}\n]*?)\s*$`)

// countKey is a property that holds a count: of seconds, days, or things.
var countKey = regexp.MustCompile(`_(SECONDS|MINUTES|HOURS|DAYS|COUNT)$`)

// listName is the name of a list of settings, which a fragment holds.
var listName = regexp.MustCompile(`(parameters|params|options|properties)$`)

// missingBar finds the typo many SHOW pages share: the IN list's
// <schema_name> line has no | before the APPLICATION lines after it.
var missingBar = regexp.MustCompile(`(<schema_name>)[ \t]*\n((?:[ \t]*\n)*[ \t]*APPLICATION )`)

// convertPage converts one page.
func convertPage(slug string, blocks []string, fixes []*fix) pageGrammar {
	pg := pageGrammar{slug: slug, url: docsSite + "/sql-reference/sql/" + slug, rule: ruleName(slug)}
	for _, f := range fixes {
		if f.page != slug {
			continue
		}
		for i, b := range blocks {
			if strings.Contains(b, f.find) {
				blocks[i] = strings.ReplaceAll(b, f.find, f.replace)
				f.used = true
			}
		}
		if f.used {
			pg.fixes = append(pg.fixes, f.reason)
		}
	}
	for i, b := range blocks {
		// A line that sets a property to one of several values, with no
		// bracket to say so: OAUTH_CLIENT_TYPE = 'CONFIDENTIAL' | 'PUBLIC'.
		blocks[i] = valueChoiceLine.ReplaceAllStringFunc(b, func(line string) string {
			m := valueChoiceLine.FindStringSubmatch(line)
			if strings.HasSuffix(strings.TrimSpace(m[3]), "|") {
				// EVENT_TABLE = <name> | in a list of names: the | is the
				// list's, not a choice of values.
				return line
			}
			return m[1] + m[2] + "{ " + m[3] + " }"
		})
	}
	for i, b := range blocks {
		if fixed := missingBar.ReplaceAllString(b, "$1 |\n$2"); fixed != b {
			blocks[i] = fixed
			pg.fixes = append(pg.fixes, "read: the IN list's <schema_name> line has no | before APPLICATION; one is read there")
		}
	}
	units := splitUnits(blocks, statementWords(slug))

	// Definitions first: forms refer to them by name. A name defined twice,
	// as cloudProviderParams (for Amazon S3) and (for Google Cloud Storage)
	// are, has both bodies as alternatives.
	defs := map[string][]unit{}
	var defOrder []string
	for _, u := range units {
		if u.def == "" {
			continue
		}
		if _, ok := defs[u.def]; !ok {
			defOrder = append(defOrder, u.def)
		}
		defs[u.def] = append(defs[u.def], u)
	}
	c := &converter{page: pg.rule, defs: map[string]bool{}}
	for _, name := range defOrder {
		c.defs[name] = true
	}
	// A definition that fails to convert is left out, and so is every form
	// that uses it: drop them until none fails.
	converted := map[string][]string{}
	for changed := true; changed; {
		changed = false
		for _, name := range defOrder {
			if !c.defs[name] || converted[name] != nil {
				continue
			}
			var alts []string
			var err error
			for _, u := range defs[name] {
				var a []string
				a, err = c.definition(u.lines, u.variant)
				if err != nil {
					break
				}
				alts = append(alts, a...)
			}
			if err != nil {
				c.defs[name] = false
				pg.failed = append(pg.failed, failure{what: name + " ::=", text: unitText(defs[name][0].lines), reason: err.Error()})
				converted = map[string][]string{}
				changed = true
				break
			}
			converted[name] = dedupe(alts)
		}
	}
	for _, name := range defOrder {
		if c.defs[name] {
			pg.frags = append(pg.frags, fragRule{name: c.fragName(name), alts: converted[name]})
		}
	}
	for _, u := range units {
		if u.def != "" {
			continue
		}
		alts, err := c.statement(unitText(u.lines))
		if err != nil {
			pg.failed = append(pg.failed, failure{what: "form", text: unitText(u.lines), reason: err.Error()})
			continue
		}
		pg.forms = append(pg.forms, alts...)
	}
	pg.forms = dedupe(pg.forms)
	pg.frags = usedFrags(pg)
	return pg
}

// usedFrags drops the fragments no converted form reaches.
func usedFrags(pg pageGrammar) []fragRule {
	byName := map[string]fragRule{}
	for _, f := range pg.frags {
		byName[f.name] = f
	}
	used := map[string]bool{}
	var visit func(text string)
	visit = func(text string) {
		for _, f := range pg.frags {
			if !used[f.name] && containsWord(text, f.name) {
				used[f.name] = true
				for _, a := range f.alts {
					visit(a)
				}
			}
		}
	}
	for _, a := range pg.forms {
		visit(a)
	}
	var out []fragRule
	for _, f := range pg.frags {
		if used[f.name] {
			out = append(out, f)
		}
	}
	return out
}

func containsWord(text, word string) bool {
	for i := 0; ; {
		j := strings.Index(text[i:], word)
		if j < 0 {
			return false
		}
		j += i
		end := j + len(word)
		before := j == 0 || !isWordByte(text[j-1])
		after := end == len(text) || !isWordByte(text[end])
		if before && after {
			return true
		}
		i = j + 1
	}
}

func isWordByte(b byte) bool {
	return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

func unitText(lines []codeLine) string {
	var b strings.Builder
	for _, l := range lines {
		if l.comment {
			continue
		}
		b.WriteString(l.text + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func dedupe(xs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

// ruleName turns a page's slug into a parser rule name: alter-warehouse is
// alterWarehouse.
func ruleName(slug string) string {
	parts := strings.FieldsFunc(slug, func(r rune) bool { return r == '-' || r == '_' })
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

type converter struct {
	page string
	defs map[string]bool // fragments the page defines and that converted
}

func (c *converter) fragName(def string) string { return c.page + "_" + def }

var propertyLine = regexp.MustCompile(`^\s*[A-Z][A-Z0-9_]*\s*=`)

// definition converts a ::= body. A body that lists one property per line
// (NAME = value), or variants between -- comments, is a choice among its
// lines: any number of properties, one variant. A definition written for
// one variant, name (for X) ::=, lists that variant's parameters, which
// it needs in order. Anything else is one sequence.
func (c *converter) definition(lines []codeLine, variant bool) ([]string, error) {
	type item struct {
		text    string
		comment bool
	}
	var items []item
	depth := 0
	for _, l := range lines {
		if l.comment {
			// A comment between lines, -- For X, separates variants; one
			// inside brackets only explains.
			if depth <= 0 {
				items = append(items, item{comment: true})
			}
			continue
		}
		t := strings.TrimSpace(l.text)
		joins := len(items) > 0 && !items[len(items)-1].comment &&
			(depth > 0 || strings.HasPrefix(t, "|") || strings.HasSuffix(strings.TrimSpace(items[len(items)-1].text), "|") ||
				strings.ContainsAny(t[:1], ")]}"))
		if joins {
			items[len(items)-1].text += "\n" + l.text
		} else {
			items = append(items, item{text: l.text})
		}
		depth += bracketDepth(l.text)
	}
	var code []string
	commented := false
	for _, it := range items {
		if it.comment {
			commented = true
			continue
		}
		code = append(code, it.text)
	}
	properties := len(code) >= 2 && !variant
	for _, t := range code {
		if !propertyLine.MatchString(t) {
			properties = false
		}
	}
	switch {
	case properties:
		var alts []string
		for _, t := range code {
			// NAME = A | B: the choice is the value's.
			k := strings.Index(t, "=")
			t = t[:k+1] + " { " + t[k+1:] + " }"
			a, err := c.form(t)
			if err != nil {
				return nil, err
			}
			alts = append(alts, a...)
		}
		return []string{"(" + strings.Join(dedupe(alts), " | ") + ")+"}, nil
	case commented && len(code) >= 2:
		var alts []string
		for _, t := range code {
			a, err := c.form(t)
			if err != nil {
				return nil, err
			}
			alts = append(alts, a...)
		}
		return alts, nil
	}
	return c.form(strings.Join(code, "\n"))
}

// errScope marks a statement form with a | outside any bracket whose
// scope the docs do not show.
var errScope = errors.New("a | outside brackets, whose scope the docs do not show")

// statement converts a statement form. A | outside brackets is how the
// docs write a choice of what follows the object's name, or of the last
// word when the choices are single words:
//
//	ALTER TASK [ IF EXISTS ] <name> RESUME | SUSPEND
//	ALTER INTEGRATION <name> UNSET ENABLED | COMMENT
//
// so an alternative that does not start with the statement's first word
// takes the first one's words before the choice.
func (c *converter) statement(text string) ([]string, error) {
	return c.convert(text, true)
}

// form converts notation to ANTLR alternatives.
func (c *converter) form(text string) ([]string, error) {
	return c.convert(text, false)
}

// scopeChoices gives each statement alternative that does not start with
// the statement's first word the first alternative's words through the
// object's name.
func scopeChoices(tree [][]*node) ([][]*node, error) {
	if len(tree) < 2 || len(tree[0]) == 0 || tree[0][0].tok == nil || tree[0][0].tok.kind != nKeyword {
		return tree, nil
	}
	head := tree[0]
	verb := head[0].tok.text
	name := -1
	for i, n := range head {
		if isName(n) {
			name = i
			break
		}
	}
	// The choice starts after the name, or is the head's last word when
	// each other choice is one word and the head has more after the name.
	from := name + 1
	// <name> SET A | B | C: each choice is set, so SET stays with each.
	if from < len(head) && head[from].tok != nil && head[from].tok.kind == nKeyword && (head[from].tok.text == "SET" || head[from].tok.text == "UNSET") {
		shared := true
		for _, alt := range tree[1:] {
			if len(alt) > 0 && alt[0].tok != nil && alt[0].tok.text == head[from].tok.text {
				shared = false
			}
		}
		if shared {
			from++
		}
	}
	words := true
	for _, alt := range tree[1:] {
		if len(alt) != 1 || alt[0].tok == nil || alt[0].tok.kind != nKeyword {
			words = false
		}
	}
	if words && name >= 0 && len(head)-name-1 > 1 {
		from = len(head) - 1
	}
	out := [][]*node{head}
	for _, alt := range tree[1:] {
		if len(alt) > 0 && alt[0].tok != nil && alt[0].tok.text == verb {
			out = append(out, alt)
			continue
		}
		if name < 0 || name == len(head)-1 {
			return nil, errScope
		}
		out = append(out, append(append([]*node(nil), head[:from]...), alt...))
	}
	return out, nil
}

// isName reports whether a node is a value, such as an object's <name>,
// or an optional one.
func isName(n *node) bool {
	if n.tok != nil {
		return n.tok.kind == nPlaceholder || n.tok.kind == nTerminal
	}
	return n.op == '?' && len(n.alts) == 1 && len(n.alts[0]) == 1 && isName(n.alts[0][0])
}

func (c *converter) convert(text string, statement bool) ([]string, error) {
	toks, err := lexNotation(text)
	if err != nil {
		return nil, err
	}
	for len(toks) > 0 && toks[len(toks)-1].kind == nLiteral && toks[len(toks)-1].text == ";" {
		toks = toks[:len(toks)-1]
	}
	if len(toks) == 0 {
		return nil, fmt.Errorf("no syntax")
	}
	if err := c.checkWords(toks); err != nil {
		return nil, err
	}
	// CREATE <integration_type> INTEGRATION, ALTER [ <integration_type> ]
	// INTEGRATION: a placeholder standing for the keyword of a type, whose
	// own page gives its syntax. Not a value.
	for i, t := range toks {
		next := i + 1
		if next < len(toks) && toks[next].kind == nClose && i > 0 && toks[i-1].kind == nOpen {
			next++ // [ <integration_type> ] INTEGRATION
		}
		if t.kind == nPlaceholder && strings.HasSuffix(t.text, "_type") && placeholderClass(t.text) == tName &&
			i > 0 && toks[i-1].kind != nKeyword && !isLiteral(toks[i-1], "=") &&
			next < len(toks) && toks[next].kind == nKeyword {
			return nil, fmt.Errorf("<%s> stands for a type's keyword, whose page gives its syntax", t.text)
		}
	}
	// REFRESH_INTERVAL_SECONDS = <value>: a value a key names as a count of
	// time units, or a count, is a number, whatever the placeholder says.
	for i, t := range toks {
		if t.kind == nPlaceholder && placeholderClass(t.text) == tName && i >= 2 && isLiteral(toks[i-1], "=") &&
			toks[i-2].kind == nKeyword && countKey.MatchString(toks[i-2].text) {
			toks[i].text = "num"
		}
	}
	// KEY = with nothing after it is a value the docs left out; an = alone
	// in brackets, INCREMENT [ BY ] [ = ] <n>, is an optional sign.
	for i, t := range toks {
		if isLiteral(t, "=") && i > 0 && toks[i-1].kind != nOpen && (i+1 == len(toks) || toks[i+1].kind == nClose || toks[i+1].kind == nBar) {
			return nil, fmt.Errorf("an = with no value after it")
		}
	}
	// AS SELECT <query>: the query's sample starts with its own SELECT.
	for i := len(toks) - 2; i >= 0; i-- {
		if toks[i].kind == nKeyword && toks[i].text == "SELECT" && toks[i+1].kind == nPlaceholder && placeholderClass(toks[i+1].text) == tQuery {
			toks = append(toks[:i], toks[i+1:]...)
		}
	}
	tree, err := parseNotation(toks)
	if err != nil {
		return nil, err
	}
	tree, err = resolveEllipses(tree)
	if err != nil {
		return nil, err
	}
	if statement {
		if tree, err = scopeChoices(tree); err != nil {
			return nil, err
		}
	}
	var alts []string
	for _, seq := range tree {
		s, err := c.seq(seq)
		if err != nil {
			return nil, err
		}
		if s == "" {
			continue
		}
		alts = append(alts, s)
	}
	if len(alts) == 0 {
		return nil, fmt.Errorf("no syntax")
	}
	return alts, nil
}

func (c *converter) seq(seq []*node) (string, error) {
	var parts []string
	for _, n := range seq {
		s, err := c.node(n)
		if err != nil {
			return "", err
		}
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, " "), nil
}

func (c *converter) node(n *node) (string, error) {
	if n.tok != nil {
		s, err := c.token(*n.tok)
		if err != nil || n.op == 0 || s == "" {
			return s, err
		}
		// intoClause [ ... ]: a repeated token keeps its operator.
		if strings.ContainsAny(s, " ()") {
			s = "(" + s + ")"
		}
		return s + string(n.op), nil
	}
	// An empty alternative, { A | B | }, is a | left before the brace, as
	// the SHOW pages' IN lists have: it adds nothing to choose.
	var alts []string
	for _, seq := range n.alts {
		s, err := c.seq(seq)
		if err != nil {
			return "", err
		}
		if s != "" {
			alts = append(alts, s)
		}
	}
	alts = dedupe(alts)
	if len(alts) == 0 {
		return "", nil
	}
	body := strings.Join(alts, " | ")
	if n.op == 0 {
		if len(alts) == 1 {
			return body, nil
		}
		return "(" + body + ")", nil
	}
	if len(alts) == 1 && !strings.ContainsAny(body, " ()") {
		return body + string(n.op), nil
	}
	return "(" + body + ")" + string(n.op), nil
}

func (c *converter) token(t ntok) (string, error) {
	switch t.kind {
	case nKeyword:
		return quote(t.text), nil
	case nLiteral:
		return quote(t.text), nil
	case nTerminal:
		return t.text, nil
	case nPlaceholder:
		if c.defs[t.text] {
			return c.fragName(t.text), nil
		}
		return placeholderClass(t.text), nil
	case nFragment:
		if c.defs[t.text] {
			return c.fragName(t.text), nil
		}
		if strings.ToLower(t.text) == t.text {
			// A lowercase value, as in API_PROVIDER = { aws_api_gateway | ... }:
			// checkWords let it through only where a value stands.
			return quote(t.text), nil
		}
		return "", fmt.Errorf("%w fragment %q: the page does not define it, or it did not convert", errUndefined, t.text)
	case nEllipsis:
		return "", errElided
	}
	return "", fmt.Errorf("token %q", t.text)
}

// checkWords refuses a lowercase word the page does not define unless it
// stands where a value does, first in its alternative (after =, |, {, [
// or (), or where a property's name does, before =. Anywhere else, as in
// TYPE = OAUTH2 oauth_parameters, it names a fragment the page never
// wrote down.
func (c *converter) checkWords(toks []ntok) error {
	for i, t := range toks {
		if t.kind != nFragment || c.defs[t.text] || strings.ToLower(t.text) != t.text {
			continue
		}
		valueBefore := i > 0 && (toks[i-1].kind == nBar || toks[i-1].kind == nOpen || isLiteral(toks[i-1], "=") || isLiteral(toks[i-1], "("))
		valueAfter := i+1 < len(toks) && isLiteral(toks[i+1], "=")
		// A word alone in braces is text to fill in, as <name> is:
		// APPLICATION {app_name}. One that names a list of parameters,
		// {oauth_parameters}, is a fragment the page never wrote down.
		if i > 0 && toks[i-1].kind == nOpen && toks[i-1].text == "{" && i+1 < len(toks) && toks[i+1].kind == nClose {
			if listName.MatchString(t.text) {
				return fmt.Errorf("%w fragment %q: the page does not define it", errUndefined, t.text)
			}
			toks[i].kind = nPlaceholder
			continue
		}
		if !valueBefore && !valueAfter {
			return fmt.Errorf("%w fragment %q: the page does not define it", errUndefined, t.text)
		}
	}
	return nil
}

// quote writes text as an ANTLR literal.
func quote(s string) string {
	return "'" + strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s) + "'"
}

// writeG4 writes the grammar: a statement rule naming every page's rule,
// then each page's rules under its URL, and what was not converted.
func writeG4(pin string, pages []pageGrammar) string {
	var b strings.Builder
	b.WriteString(`// Snowflake's SQL command reference as an ANTLR 4 grammar: the Syntax
// sections of the pages https://docs.snowflake.com/en/sql-reference/sql-all
// lists, converted by snowflake/parser/cmd/docsgrammar. Do not edit it:
// correct a page's text in docs-fixes.tsv and run the converter again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, QUERY, COPY_STATEMENT,
// DATA_TYPE, STAGE, FILE, POSITION, OBJECT_TYPE and OBJECT_TYPES are the
// values the docs write as <placeholders>; tokens.tsv samples them. A
// form a page writes with ... for text it leaves out, or with a word it
// never defines, is not converted and is listed under its page.
`)
	fmt.Fprintf(&b, "//\n// pin: %s\n\ngrammar SnowflakeDocs;\n\nstatement\n", pin)
	first := true
	for _, pg := range pages {
		if len(pg.forms) == 0 {
			continue
		}
		if first {
			fmt.Fprintf(&b, "    : %s\n", pg.rule)
			first = false
		} else {
			fmt.Fprintf(&b, "    | %s\n", pg.rule)
		}
	}
	b.WriteString("    ;\n")
	for _, pg := range pages {
		fmt.Fprintf(&b, "\n// %s\n", pg.url)
		for _, f := range pg.fixes {
			fmt.Fprintf(&b, "// fix: %s\n", f)
		}
		for _, f := range pg.failed {
			fmt.Fprintf(&b, "// not converted (%s): %s\n", f.what, f.reason)
			for _, l := range strings.Split(f.text, "\n") {
				fmt.Fprintf(&b, "//   %s\n", strings.TrimRight(l, " "))
			}
		}
		if len(pg.forms) == 0 {
			if len(pg.failed) == 0 {
				b.WriteString("// no syntax\n")
			}
			continue
		}
		writeRule(&b, pg.rule, pg.forms)
		for _, f := range pg.frags {
			writeRule(&b, f.name, f.alts)
		}
	}
	return b.String()
}

func writeRule(b *strings.Builder, name string, alts []string) {
	fmt.Fprintf(b, "%s\n", name)
	for i, a := range alts {
		lead := "    |"
		if i == 0 {
			lead = "    :"
		}
		fmt.Fprintf(b, "%s %s\n", lead, a)
	}
	b.WriteString("    ;\n")
}

// classes lists the token classes a grammar uses, for the report.
func classes(g4 string) []string {
	var out []string
	for _, c := range []string{tName, tString, tNumber, tExpr, tQuery, tType, tStage, tPosition, tObjectType, tObjectTypes, tFile, tCopy} {
		if containsWord(g4, c) {
			out = append(out, c)
		}
	}
	sort.Strings(out)
	return out
}
