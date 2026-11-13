// Command docsgrammar converts the Syntax sections of Snowflake's SQL
// command reference to snowflake/testdata/SnowflakeDocs.g4, the grammar
// snowflake/parser's compat.tsv is generated from and judged by, and pins
// what it read in snowflake/testdata/pins.tsv: the index page, a SHA-256 of
// every Syntax block it read, and the date it read them. Run it from the
// repo root to bump the pin, or after a change to docs-fixes.tsv:
//
//	go run ./snowflake/parser/cmd/docsgrammar -cache /tmp/snowflake-docs
//
// With -cache it reads pages from the directory, and keeps the ones it
// fetches there. The date in the pin moves only when the blocks do.
//
// The converter is enginetest/grammar/docs; this file says where
// Snowflake's pages are and what its placeholders stand for.
package main

import (
	"html"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

// Token classes of Snowflake's own placeholders, beside the ones every
// docs grammar has (docs.ClassName and the rest).
const (
	tObjectType  = "OBJECT_TYPE"
	tObjectTypes = "OBJECT_TYPES"
	tCopy        = "COPY_STATEMENT"
)

var dialect = &docs.Dialect{
	Grammar:  "SnowflakeDocs",
	Command:  "snowflake/parser/cmd/docsgrammar",
	Testdata: filepath.Join("snowflake", "testdata"),
	Header: `// Snowflake's SQL command reference as an ANTLR 4 grammar: the Syntax
// sections of the pages https://docs.snowflake.com/en/sql-reference/sql-all
// lists, converted by snowflake/parser/cmd/docsgrammar. Do not edit it:
// correct a page's text in docs-fixes.tsv and run the converter again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, QUERY, COPY_STATEMENT,
// DATA_TYPE, STAGE, FILE, POSITION, OBJECT_TYPE and OBJECT_TYPES are the
// values the docs write as <placeholders>; tokens.tsv samples them. A
// form a page writes with ... for text it leaves out, or with a word it
// never defines, is not converted and is listed under its page.
`,
	Site:           docsSite,
	Index:          indexPath,
	Pages:          commandPages,
	Blocks:         syntaxBlocks,
	Prepare:        prepare,
	Classify:       placeholderClass,
	StatementWords: statementWords,
}

func main() { docs.Run(dialect) }

// docsSite is where the command reference lives. The index lists every
// command page by a link relative to the site root.
const (
	docsSite  = "https://docs.snowflake.com/en"
	indexPath = "/sql-reference/sql-all"
)

var commandLink = regexp.MustCompile(`href="(/(?:en/)?sql-reference/sql/[a-z0-9_-]+)"`)

// commandPages returns the command pages the index links to, in the
// index's order, each once.
func commandPages(index string) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range commandLink.FindAllStringSubmatch(index, -1) {
		p := strings.TrimPrefix(m[1], "/en")
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

var (
	preBlock   = regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`)
	anyTag     = regexp.MustCompile(`(?s)<[^>]+>`)
	nextH2     = regexp.MustCompile(`<h2[ >]`)
	syntaxSecs = regexp.MustCompile(`data-section="syntax[a-z0-9-]*"`)
)

// syntaxBlocks returns the text of every code block in a page's Syntax
// sections: Syntax, or Syntax for ... where a page has several, each from
// its heading to the next second-level heading, so the subsections some
// pages split their syntax into are included. A page with none, such as an
// overview, returns none.
func syntaxBlocks(doc string) []string {
	var out []string
	end := 0
	for _, loc := range syntaxSecs.FindAllStringIndex(doc, -1) {
		if loc[0] < end {
			continue // a subsection of the section just read
		}
		sec := doc[loc[0]:]
		// Skip the section's own <h2>, then stop at the next one.
		if h := strings.Index(sec, "</h2>"); h >= 0 {
			if j := nextH2.FindStringIndex(sec[h:]); j != nil {
				sec = sec[:h+j[0]]
			}
		}
		end = loc[0] + len(sec)
		for _, m := range preBlock.FindAllStringSubmatch(sec, -1) {
			text := html.UnescapeString(anyTag.ReplaceAllString(m[1], ""))
			text = strings.ReplaceAll(text, "\u00a0", " ")
			out = append(out, strings.TrimRight(text, " \n\t"))
		}
	}
	return out
}

var (
	exprName       = regexp.MustCompile(`(^|_)(expr|exprs|expression|condition|predicate)(_?[0-9]+|_|$)`)
	numberName     = regexp.MustCompile(`(^|_)(num|number|integer|int|n|rows|seconds|minutes|hours|days|bytes|count|percent|ptus|buckets|interval|ms|threshold|width)(_?[0-9]+|_|$)|^(initial|start|min|max|step)_value$|^increment$`)
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
		return docs.ClassString
	case runtimeVersion.MatchString(n):
		// RUNTIME_VERSION = '3.10': a runtime's version is a string.
		return docs.ClassString
	case queryName.MatchString(n):
		return docs.ClassQuery
	case typeName.MatchString(n):
		return docs.ClassDataType
	case exprName.MatchString(n):
		return docs.ClassExpr
	case stringName.MatchString(n):
		return docs.ClassString
	case strings.Contains(n, "name") || strings.Contains(n, "identifier"):
		return docs.ClassName
	case numberName.MatchString(n):
		return docs.ClassNumber
	}
	return docs.ClassName
}

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

// missingBar finds the typo many SHOW pages share: the IN list's
// <schema_name> line has no | before the APPLICATION lines after it.
var missingBar = regexp.MustCompile(`(<schema_name>)[ \t]*\n((?:[ \t]*\n)*[ \t]*APPLICATION )`)

// prepare corrects the typo many SHOW pages share.
func prepare(_ string, blocks []string) []string {
	var notes []string
	for i, b := range blocks {
		if fixed := missingBar.ReplaceAllString(b, "$1 |\n$2"); fixed != b {
			blocks[i] = fixed
			notes = append(notes, "read: the IN list's <schema_name> line has no | before APPLICATION; one is read there")
		}
	}
	return notes
}
