// Command docsgrammar converts the syntax of Oracle's SQL Language
// Reference, as oracle/parser/bnf/ quotes it, to oracle/testdata/OracleDocs.g4,
// the grammar oracle/parser's compat.tsv is generated from, and pins what it
// read in the docs row of oracle/testdata/pins.tsv: a SHA-256 of every file
// it read, and the date they were quoted. Run it from the repo root after a
// change to oracle/parser/bnf/ or to docs-fixes.tsv:
//
//	go run ./oracle/parser/cmd/docsgrammar -date 2026-03-15
//
// The converter is enginetest/grammar/docs; this file reads the BNF files
// and rewrites their notation into the one it reads. The pinned Oracle
// judges the rows, not the grammar: the grammar only says what to ask.
package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

// bnfDir holds one file per statement page, each starting with the page's
// URL on a -- Source: line.
var bnfDir = filepath.Join("oracle", "parser", "bnf")

// site and index are where the pages are: a page's path is relative to
// site, and index, the reference's own directory, lists them.
const (
	site  = "https://docs.oracle.com/en/database/oracle/oracle-database/"
	index = "26/sqlrf/"
)

var dialect = &docs.Dialect{
	Grammar:  "OracleDocs",
	Command:  "oracle/parser/cmd/docsgrammar",
	Testdata: filepath.Join("oracle", "testdata"),
	Header: `// Oracle's SQL Language Reference as an ANTLR 4 grammar: the syntax of the
// statement pages oracle/parser/bnf/ quotes, converted by
// oracle/parser/cmd/docsgrammar. Do not edit it: correct a page's text in
// docs-fixes.tsv, or the quote in oracle/parser/bnf/, and run the converter
// again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, CONDITION, QUERY,
// DATA_TYPE and HINT are the values the pages write in lower case and
// define nowhere; tokens.tsv samples them. A fragment another page defines
// is copied under the page that uses it, except a query, which is QUERY. A
// clause no quoted page defines is left out where it is optional, and the
// form or definition that requires it is not converted. Those, and a form
// that does not read as notation, are listed under their page.
`,
	Site:           site,
	Index:          index,
	Pages:          func(index string) []string { return strings.Fields(index) },
	Slug:           func(p string) string { return strings.TrimSuffix(path.Base(p), ".html") },
	RuleName:       func(slug string) string { return docs.RuleName(strings.ToLower(slug)) },
	Blocks:         blocks,
	Classify:       placeholderClass,
	StatementWords: statementWords,
	Read:           read,
	Prepare:        prepare,
	PinRole:        "docs",
}

func main() {
	if err := load(); err != nil {
		fmt.Fprintln(os.Stderr, "docsgrammar:", err)
		os.Exit(1)
	}
	docs.Run(dialect)
}

// A page is one BNF file: its statement form and its definitions.
type page struct {
	slug, path string
	form       string            // the statement's syntax
	defs       map[string]string // name, lower case, to its body
	order      []string          // defs in the file's order
}

var (
	pages   = map[string]*page{} // by path relative to Site
	byPath  []string
	defined = map[string]string{} // a definition's name to the first page, in file order, that defines it
)

var (
	sourceLine = regexp.MustCompile(`(?m)^-- Source: (\S+)`)
	defHead    = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\s*(?:::=|:)\s*(.*)$`)
)

// load reads every BNF file. A line at the start of a line that is name ::=
// or name: (a name with a lower-case letter) starts a definition; the first
// unit of a file that is not one is the statement form, and where the file
// starts with a definition, as select: does, that definition's body is the
// form too. Every other line continues the unit before it.
func load() error {
	files, err := filepath.Glob(filepath.Join(bnfDir, "*.bnf"))
	if err != nil {
		return err
	}
	fixes, err := readBNFFixes(filepath.Join(dialect.Testdata, "bnf-fixes.tsv"))
	if err != nil {
		return err
	}
	sort.Strings(files)
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		text := string(raw)
		slug := strings.TrimSuffix(filepath.Base(f), ".bnf")
		for _, fx := range fixes {
			if fx.page != slug {
				continue
			}
			if !strings.Contains(text, fx.find) {
				return fmt.Errorf("bnf-fixes.tsv: %s: %q is not in the file", fx.page, fx.find)
			}
			text = strings.ReplaceAll(text, fx.find, fx.replace)
			fx.used = true
		}
		b := []byte(text)
		m := sourceLine.FindStringSubmatch(string(b))
		if m == nil || !strings.HasPrefix(m[1], site) {
			return fmt.Errorf("%s: no -- Source: line under %s", f, site)
		}
		p := &page{slug: slug, path: strings.TrimPrefix(m[1], site), defs: map[string]string{}}
		var name string
		var body []string
		formDone := false
		flush := func() {
			text := strings.TrimRight(strings.Join(body, "\n"), "\n ")
			switch {
			case name != "":
				if _, ok := p.defs[name]; !ok {
					p.order = append(p.order, name)
				}
				p.defs[name] = text
				if p.form == "" && !formDone {
					p.form = text
				}
			case strings.TrimSpace(text) != "":
				p.form = text
			}
			formDone = formDone || strings.TrimSpace(text) != ""
			name, body = "", nil
		}
		for _, line := range strings.Split(string(b), "\n") {
			t := strings.TrimSpace(line)
			if t == "" || strings.HasPrefix(t, "--") {
				continue
			}
			if m := defHead.FindStringSubmatch(line); m != nil && strings.ToUpper(m[1]) != m[1] {
				flush()
				name = strings.ToLower(m[1])
				body = []string{m[2]}
				continue
			}
			body = append(body, line)
		}
		flush()
		if p.form == "" {
			return fmt.Errorf("%s: no statement form", f)
		}
		if _, dup := pages[p.path]; dup {
			return fmt.Errorf("%s: a second file for %s", f, p.path)
		}
		pages[p.path] = p
		byPath = append(byPath, p.path)
		for _, n := range p.order {
			if _, ok := defined[n]; !ok {
				defined[n] = p.path
			}
		}
	}
	for _, fx := range fixes {
		if !fx.used {
			return fmt.Errorf("bnf-fixes.tsv: no file %s.bnf", fx.page)
		}
	}
	return nil
}

type bnfFix struct {
	page, find, replace string
	used                bool
}

// readBNFFixes reads bnf-fixes.tsv: page, text, replacement, reason.
func readBNFFixes(file string) ([]*bnfFix, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	unescape := strings.NewReplacer(`\n`, "\n", `\t`, "\t")
	var out []*bnfFix
	for i, line := range strings.Split(string(b), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) != 4 {
			return nil, fmt.Errorf("%s:%d: want page, text, replacement and reason", file, i+1)
		}
		out = append(out, &bnfFix{page: cols[0], find: unescape.Replace(cols[1]), replace: unescape.Replace(cols[2])})
	}
	return out, nil
}

// read answers the converter's reads from the loaded files: the index lists
// the pages, and a page is its file's path, which blocks finds again.
func read(p string) (string, error) {
	if p == index {
		return strings.Join(byPath, "\n"), nil
	}
	if _, ok := pages[p]; !ok {
		return "", fmt.Errorf("no BNF file for %s", p)
	}
	return p, nil
}

// queryFragments are definitions of SELECT's page another page names, which
// stand for a whole query there: QUERY.
var queryFragments = map[string]bool{"subquery": true, "query_block": true}

// clauseName is a name that stands for syntax, not a value: one no quoted
// page defines is left out of the grammar rather than sampled as a name.
var clauseName = regexp.MustCompile(`(_clauses?|_parameters|_params|_spec|_specification|_props|_properties|_declaration|_definition|_source|_statement|_storage|_virtual_columns|_attributes|_shortcut)$|^(statement|segment_attributes|sql_statement)$`)

var word = regexp.MustCompile(`'[^'\n]*'|[A-Za-z_][A-Za-z0-9_$#]*`)

// blocks returns a page's syntax in the converter's notation: the statement
// form, then name ::= for each definition it reaches, its own page's first
// and then the first page in file order that defines it.
func blocks(p string) []string {
	pg := pages[p]
	leftOut = map[string]bool{}
	defer func() { notes[pg.slug] = sortedKeys(leftOut) }()
	type def struct{ name, body string }
	var defs []def
	seen := map[string]bool{}
	var queue []string
	reach := func(text string) {
		for _, w := range word.FindAllString(text, -1) {
			n := strings.ToLower(w)
			if w[0] == '\'' || !hasLower(w) || seen[n] {
				continue
			}
			if _, own := pg.defs[n]; !own {
				if _, ok := defined[n]; !ok || queryFragments[n] {
					continue
				}
			}
			seen[n] = true
			queue = append(queue, n)
		}
	}
	reach(pg.form)
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		body, ok := pg.defs[n]
		if !ok {
			body = pages[defined[n]].defs[n]
		}
		defs = append(defs, def{n, body})
		reach(body)
	}
	// A definition that still needs a clause no page defines is itself
	// left out where it is optional, as the clause is: drop them until none
	// is left that needs one.
	usable := map[string]bool{}
	for n := range seen {
		usable[n] = true
	}
	rewritten := map[string]string{}
	for changed := true; changed; {
		changed = false
		for _, d := range defs {
			if !usable[d.name] {
				continue
			}
			rewritten[d.name] = rewrite(d.body, usable)
			if refused.MatchString(rewritten[d.name]) {
				usable[d.name] = false
				changed = true
			}
		}
	}
	out := []string{rewrite(pg.form, usable)}
	for _, d := range defs {
		if usable[d.name] {
			out = append(out, d.name+" ::=\n"+rewritten[d.name])
		}
	}
	return out
}

// refused is a clause dropUndefined could not leave out.
var refused = regexp.MustCompile(`<[a-z]+[A-Z][A-Za-z0-9]*>`)

// seenDef reports whether n is a definition the page reaches that was left
// out.
func seenDef(defs map[string]bool, n string) bool {
	usable, ok := defs[n]
	return ok && !usable
}

func hasLower(w string) bool { return strings.ToUpper(w) != w }

// rewrite turns BNF text into the notation: a defined name in lower case, a
// name no page defines into a <placeholder>, a quoted value into '<value>',
// @ apart from what follows it. A clause no page defines is dropped with
// the innermost [ ] around it, and where none is, becomes a reference the
// converter refuses.
func rewrite(text string, defs map[string]bool) string {
	var b strings.Builder
	last := 0
	for _, loc := range word.FindAllStringIndex(text, -1) {
		w := text[loc[0]:loc[1]]
		b.WriteString(text[last:loc[0]])
		last = loc[1]
		n := strings.ToLower(w)
		switch {
		case w[0] == '\'':
			if hasLower(w) || strings.ContainsAny(w, "[]") {
				b.WriteString("'<" + strings.Trim(strings.NewReplacer("[", "", "]", "", " ", "").Replace(w), "'") + ">'")
			} else {
				b.WriteString(w)
			}
		case !hasLower(w):
			b.WriteString(w)
		case !strings.Contains(w, "_") && w[0] >= 'A' && w[0] <= 'Z' && !defs[n]:
			// XMLType, XMLSchema: a keyword the docs write in mixed case.
			b.WriteString(strings.ToUpper(w))
		case defs[n]:
			b.WriteString(n)
		case seenDef(defs, n):
			// Defined, and left out: it needs a clause no page defines.
			leftOut[n] = true
			b.WriteString("<" + undefinedMark + n + ">")
		case queryFragments[n]:
			b.WriteString("<subquery>")
		case plsqlSources[n] != "":
			b.WriteString("<" + n + ">")
		case clauseName.MatchString(n):
			leftOut[n] = true
			b.WriteString("<" + undefinedMark + n + ">")
		default:
			b.WriteString("<" + w + ">")
		}
	}
	b.WriteString(text[last:])
	s := strings.ReplaceAll(b.String(), "@", "@ ")
	return dropUndefined(s)
}

// leftOut collects the clauses blocks leaves out of the page it is
// converting, and notes keeps them by page for prepare to report.
var (
	leftOut = map[string]bool{}
	notes   = map[string][]string{}
)

func sortedKeys(m map[string]bool) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// prepare reports, under each page, the clauses its syntax left out.
func prepare(slug string, blocks []string) ([]string, []string) {
	if len(notes[slug]) == 0 {
		return blocks, nil
	}
	return blocks, []string{"left out where optional, since no quoted page defines it or what it needs: " + strings.Join(notes[slug], ", ")}
}

// undefinedMark prefixes a clause no page defines until dropUndefined has
// placed it.
const undefinedMark = "undefined~"

var markRE = regexp.MustCompile(`<` + regexp.QuoteMeta(undefinedMark) + `([a-z0-9_]+)>`)

// dropUndefined leaves out each clause no page defines: the alternative
// that holds it, where the innermost group around it has others; else the
// whole [ ] around it, with a ... after it; else, inside { } or ( ) with no
// other alternative, the same one level out. Where nothing optional holds
// it, it is written <aCamelName>, which the converter refuses as a fragment
// the page does not define, with the form or definition around it.
func dropUndefined(s string) string {
	for {
		loc := markRE.FindStringIndex(s)
		if loc == nil {
			return s
		}
		start, end := loc[0], loc[1] // the span to place
		for {
			open := enclosing(s, start)
			if open < 0 {
				name := markRE.FindStringSubmatch(s[loc[0]:loc[1]])[1]
				delete(leftOut, name)
				s = s[:loc[0]] + "<" + camel(name) + ">" + s[loc[1]:]
				break
			}
			close := closing(s, open)
			if close < 0 {
				return s // unbalanced: the converter says so
			}
			bars := topBars(s, open, close)
			if len(bars) > 0 {
				// Drop the alternative around the span, and one | beside it.
				from, to := open+1, close
				for _, b := range bars {
					if b < start {
						from = b
					} else if b >= end {
						to = b
						break
					}
				}
				if from == open+1 {
					to++ // the | after it goes
				}
				s = s[:from] + s[to:]
				break
			}
			if s[open] == '[' {
				after := close + 1
				if rest := strings.TrimLeft(s[after:], " \t\n"); strings.HasPrefix(rest, "...") {
					after = len(s) - len(rest) + 3
				}
				s = s[:open] + s[after:]
				break
			}
			start, end = open, close+1
		}
	}
}

// enclosing returns the index of the innermost [, { or ( open around i, or -1.
func enclosing(s string, i int) int {
	depth := 0
	for j := i - 1; j >= 0; j-- {
		switch s[j] {
		case ']', '}', ')':
			depth++
		case '[', '{', '(':
			if depth == 0 {
				return j
			}
			depth--
		}
	}
	return -1
}

// closing returns the index of what closes the bracket at open, or -1.
func closing(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '[', '{', '(':
			depth++
		case ']', '}', ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// topBars returns the indexes of the | between open and close outside any
// nested bracket.
func topBars(s string, open, close int) []int {
	var out []int
	depth := 0
	for i := open + 1; i < close; i++ {
		switch s[i] {
		case '[', '{', '(':
			depth++
		case ']', '}', ')':
			depth--
		case '|':
			if depth == 0 && (i+1 >= len(s) || s[i+1] != '|') && (i == 0 || s[i-1] != '|') {
				out = append(out, i)
			}
		}
	}
	return out
}

// camel writes storage_clause as storageClause, the converter's spelling of
// a fragment name.
func camel(n string) string {
	parts := strings.Split(n, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// statementWords is the first word of the page's form.
func statementWords(slug string) map[string]bool {
	for _, p := range pages {
		if p.slug == slug {
			if f := strings.Fields(p.form); len(f) > 0 {
				return map[string]bool{strings.ToUpper(strings.Trim(f[0], "[{(")): true}
			}
		}
	}
	return map[string]bool{}
}

var (
	exprName   = regexp.MustCompile(`^(expr|expression|.*_expr|.*_expression|literal|value|date|timestamp|.*_value)$`)
	numberName = regexp.MustCompile(`^(integer|number|n|m|.*_number|.*number|quantity|.*_quantity|percent|.*_percent|rowcount|.*_cost|default_selectivity|compression_ratio|file_number|filenumber|instance_id|offset|position|size_value|blocksize_value|scn_value|seed_value)$`)
	stringName = regexp.MustCompile(`(?i)^(.*filename|.*_location\d*|keystore\d_location|library_path|path_name|match_string|search_string|string|version_string|.*_version|opaque_format_spec|language_options_string|hsm_auth_string|xmlschema_url|source_char|server_file_name|text)$`)
)

// plsqlSources are what another reference defines and the SQL pages name,
// as the PL/SQL Language Reference does a unit's source after its header:
// a sample of each, in tokens.tsv, stands for it.
var plsqlSources = map[string]string{
	"plsql_function_source":     "FUNCTION_SOURCE",
	"plsql_procedure_source":    "PROCEDURE_SOURCE",
	"plsql_package_source":      "PACKAGE_SOURCE",
	"plsql_package_body_source": "PACKAGE_BODY_SOURCE",
	"plsql_trigger_source":      "TRIGGER_SOURCE",
	"plsql_type_source":         "TYPE_SOURCE",
	"plsql_type_body_source":    "TYPE_BODY_SOURCE",
	"call_spec":                 "CALL_SPEC",
	// A statement or declaration of its own, which a page names inside
	// another: CREATE SCHEMA's statements, a type's method specifications,
	// WITH's PL/SQL declarations, a redo log or data file named by a string.
	"create_table_statement":  "CREATE_TABLE",
	"create_view_statement":   "CREATE_VIEW",
	"grant_statement":         "GRANT",
	"subprogram_spec":         "SUBPROGRAM_SPEC",
	"map_order_function_spec": "MAP_ORDER_SPEC",
	"function_declaration":    "FUNCTION_DECLARATION",
	"procedure_declaration":   "PROCEDURE_DECLARATION",
	"redo_log_file_spec":      docs.ClassString,
	"datafile_tempfile_spec":  docs.ClassString,
	"sql_statement":           docs.ClassString,
}

// placeholderClass returns the token class a placeholder stands for.
func placeholderClass(name string) string {
	n := strings.ToLower(name)
	if c := plsqlSources[n]; c != "" {
		return c
	}
	switch {
	case n == "condition" || strings.HasSuffix(n, "_condition"):
		return "CONDITION"
	case n == "subquery" || n == "inline_view":
		return docs.ClassQuery
	case n == "datatype" || n == "return_type" || n == "parameter_type":
		return docs.ClassDataType
	case n == "hint":
		return "HINT"
	case numberName.MatchString(n):
		return docs.ClassNumber
	case stringName.MatchString(n):
		return docs.ClassString
	case exprName.MatchString(n):
		return docs.ClassExpr
	}
	return docs.ClassName
}
