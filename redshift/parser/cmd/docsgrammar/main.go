// Command docsgrammar converts the Syntax sections of Amazon Redshift's SQL
// command reference to redshift/testdata/RedshiftDocs.g4, the grammar
// redshift/parser's compat.tsv is generated from and judged by, and pins
// what it read in redshift/testdata/pins.tsv: the index page, a SHA-256 of
// every Syntax block it read, and the date it read them. Run it from the
// repo root to bump the pin, or after a change to docs-fixes.tsv:
//
//	go run ./redshift/parser/cmd/docsgrammar -cache /tmp/redshift-docs
//
// The converter is enginetest/grammar/docs; this file says where AWS's
// pages are and rewrites their notation into the one it reads.
package main

import (
	"html"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

var dialect = &docs.Dialect{
	Grammar:  "RedshiftDocs",
	Command:  "redshift/parser/cmd/docsgrammar",
	Testdata: filepath.Join("redshift", "testdata"),
	Header: `// Amazon Redshift's SQL command reference as an ANTLR 4 grammar: the
// Syntax sections of the pages
// https://docs.aws.amazon.com/redshift/latest/dg/c_SQL_commands.html lists,
// converted by redshift/parser/cmd/docsgrammar. Do not edit it: correct a
// page's text in docs-fixes.tsv and run the converter again.
//
// Keywords are literals. ID, STRING, NUMBER, EXPR, QUERY, DATA_TYPE, CTE,
// ARGMODE and SCHEMA_ELEMENT are the values the docs write in italics or in lower case;
// tokens.tsv samples them. A form a page writes with ... for text it leaves out, or that does
// not read as notation, is not converted and is listed under its page. The
// pages a command page lists under Topics, such as SELECT's clauses, are
// not read: their syntax is a clause's, not a statement's.
`,
	Site:           "https://docs.aws.amazon.com",
	Index:          "/redshift/latest/dg/c_SQL_commands.html",
	Pages:          commandPages,
	Slug:           func(p string) string { return strings.TrimSuffix(path.Base(p), ".html") },
	RuleName:       ruleName,
	Blocks:         syntaxBlocks,
	Classify:       placeholderClass,
	StatementWords: statementWords,
}

func main() { docs.Run(dialect) }

const pageDir = "/redshift/latest/dg/"

var (
	topics    = regexp.MustCompile(`(?s)<div class="highlights">.*?</ul>`)
	topicLink = regexp.MustCompile(`href="\./([A-Za-z0-9_-]+\.html)"`)
)

// commandPages returns the pages the index's Topics list links to, in its
// order, each once: one per command.
func commandPages(index string) []string {
	list := topics.FindString(index)
	var out []string
	seen := map[string]bool{}
	for _, m := range topicLink.FindAllStringSubmatch(list, -1) {
		if p := pageDir + m[1]; !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}
	return out
}

// ruleName turns a page's slug into a parser rule name:
// r_CREATE_TABLE_NEW is createTableNew, materialized-view-create-sql-command
// materializedViewCreateSqlCommand.
func ruleName(slug string) string {
	return docs.RuleName(strings.ToLower(strings.TrimPrefix(slug, "r_")))
}

var (
	heading    = regexp.MustCompile(`(?s)<h([23])[^>]*>(.*?)</h[23]>`)
	syntaxHead = regexp.MustCompile(`(?i)syntax|synopsis`)
	preBlock   = regexp.MustCompile(`(?s)<pre class="(?:programlisting|screen)"[^>]*>(.*?)</pre>`)
	emSpan     = regexp.MustCompile(`(?s)<em>(.*?)</em>`)
	anyTag     = regexp.MustCompile(`(?s)<[^>]+>`)
)

// syntaxBlocks returns the code blocks of a page's Syntax sections, in the
// notation docs reads. A section starts at a heading that says Syntax (or
// whose id says synopsis), as "Syntax", "COPY syntax" and "Full CREATE
// MODEL syntax" do, and runs to the next heading of its level or above. A
// page with none, such as START TRANSACTION's, which points at BEGIN,
// returns none.
func syntaxBlocks(doc string) []string {
	hs := heading.FindAllStringSubmatchIndex(doc, -1)
	var out []string
	for i, h := range hs {
		if !syntaxHead.MatchString(doc[h[0]:h[1]]) || strings.Contains(doc[h[0]:h[1]], "awsdocs-topiclist-title") {
			continue
		}
		level := doc[h[2]:h[3]]
		end := len(doc)
		for _, next := range hs[i+1:] {
			if doc[next[2]:next[3]] <= level {
				end = next[0]
				break
			}
		}
		for _, m := range preBlock.FindAllStringSubmatch(doc[h[1]:end], -1) {
			out = append(out, notation(m[1]))
		}
	}
	return out
}

// lowerKeywords are keywords the pages write in lower case: IAM_ROLE {
// default | ... }, AUTHENTICATION [ none | iam | mtls ], LANGUAGE {
// plpythonu | sql }, LANGUAGE plpgsql.
var lowerKeywords = map[string]bool{"default": true, "none": true, "iam": true, "mtls": true, "plpythonu": true, "plpgsql": true, "sql": true}

// defHead is how a page defines a fragment, the name in italics: where
// option is, and column_constraints are:, where grant_or_revoke_clause is
// one of:, where option can be:.
var defHead = regexp.MustCompile(`(?m)^[ \t]*(?:[Ww]here|and)[ \t]+\x01([A-Za-z_][A-Za-z0-9_]*)\x02[ \t]+(?:is|are|can be)(?:[ \t]+one of)?[ \t]*:?[ \t]*$`)

// notation rewrites a code block of AWS's into the notation docs reads.
// AWS writes a value the reader supplies in italics, or in lower case
// without them, where docs reads <value>; a fragment is defined by a
// "where name is:" line, where docs reads name ::=. Values in italics are
// marked \x01name\x02 while the block's tags are taken out.
func notation(block string) string {
	block = emSpan.ReplaceAllStringFunc(block, func(em string) string {
		return placeholders(html.UnescapeString(anyTag.ReplaceAllString(emSpan.FindStringSubmatch(em)[1], "")), true)
	})
	text := html.UnescapeString(anyTag.ReplaceAllString(block, ""))
	text = strings.NewReplacer("\u00a0", " ", "‘", "'", "’", "'", "“", `"`, "”", `"`, "…", "...").Replace(text)
	text = defHead.ReplaceAllString(text, "$1 ::=")
	text = variants(text)
	text = placeholders(text, false)
	text = strings.NewReplacer("\x01", "<", "\x02", ">").Replace(text)
	return strings.TrimRight(text, " \n\t")
}

// variants separates the paragraphs of a definition with a -- line, which
// docs reads as a choice of them: "where grant_or_revoke_clause is one of:"
// lists whole GRANT and REVOKE forms, a blank line between each two.
func variants(text string) string {
	lines := strings.Split(text, "\n")
	var out []string
	inDef, blank := false, false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		switch {
		case strings.HasSuffix(t, "::="):
			inDef, blank = true, false
		case t == "":
			blank = blank || inDef && len(out) > 0 && !strings.HasSuffix(strings.TrimSpace(out[len(out)-1]), "::=")
			out = append(out, l)
			continue
		case blank && inDef:
			out = append(out, "--")
		}
		blank = false
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

// placeholders marks each value in text: in italics, every word; elsewhere,
// every word in lower case that is not a keyword or a definition's name. A
// quoted string that holds lower case is a value as a whole, '<text>', and
// one in upper case ('SESSION') a literal.
func placeholders(text string, italic bool) string {
	var b strings.Builder
	rs := []rune(text)
	for i := 0; i < len(rs); {
		c := rs[i]
		switch {
		case c == '<' || c == '\x01':
			// A value already: one the docs write <account-id>, or in italics.
			end := map[rune]rune{'<': '>', '\x01': '\x02'}[c]
			j := i + 1
			for j < len(rs) && rs[j] != end && rs[j] != '\n' {
				j++
			}
			if j < len(rs) && rs[j] == end {
				// Marked, so taking the block's tags out leaves it.
				b.WriteString("\x01" + string(rs[i+1:j]) + "\x02")
				i = j + 1
				continue
			}
			b.WriteRune(c)
			i++
		case c == '\'' && !italic:
			j := i + 1
			for j < len(rs) && rs[j] != '\'' && rs[j] != '\n' {
				j++
			}
			if j >= len(rs) || rs[j] != '\'' {
				b.WriteRune(c)
				i++
				continue
			}
			body := string(rs[i+1 : j])
			if !strings.ContainsAny(body, "<\x01") && strings.ToUpper(body) != body {
				body = "\x01" + body + "\x02"
			}
			b.WriteString("'" + body + "'")
			i = j + 1
		case c == '_' || unicode.IsLetter(c):
			j := i
			for j < len(rs) && (rs[j] == '_' || rs[j] == '-' && italic || unicode.IsLetter(rs[j]) || unicode.IsDigit(rs[j]) ||
				rs[j] == '.' && j+1 < len(rs) && (unicode.IsLetter(rs[j+1]) || rs[j+1] == '_')) {
				j++
			}
			w := string(rs[i:j])
			rest := strings.TrimLeft(string(rs[j:min(j+8, len(rs))]), " \t")
			switch {
			case lowerKeywords[w]:
				b.WriteString(strings.ToUpper(w))
			case strings.HasPrefix(rest, "::="):
				b.WriteString(w)
			case italic || unicode.IsLower([]rune(w)[0]):
				b.WriteString("\x01" + w + "\x02")
			default:
				b.WriteString(w)
			}
			i = j
		default:
			b.WriteRune(c)
			i++
		}
	}
	return b.String()
}

// Token classes of Redshift's own values, beside the ones every docs
// grammar has (docs.ClassName and the rest).
const (
	classCTE     = "CTE"            // a WITH clause's common table expression
	classArgMode = "ARGMODE"        // a procedure argument's IN, OUT or INOUT
	classElement = "SCHEMA_ELEMENT" // a statement CREATE SCHEMA runs in the schema
)

var (
	cteName    = regexp.MustCompile(`^(with_subquery|common_table_expression)$`)
	queryName  = regexp.MustCompile(`(^|_)(query|subquery|select_statement|query_definition|select_clause|statement)$`)
	exprName   = regexp.MustCompile(`(^|_)(expr|exp|expression|condition|predicate)(_?[0-9]+|_|$)`)
	typeName   = regexp.MustCompile(`(^|_)(data_type|type|argtype|datatype)(_?[0-9]+|_|$)|data_type_size$`)
	numberName = regexp.MustCompile(`(^|_)(seed|step|limit|quota|seconds|milliseconds|count|size|priority|numrows|number|integer|interval|threshold|start|n)(_?[0-9]+|_|$)|^process_id$`)
)

// placeholderClass returns the token class a value stands for, by its
// name. Anything not recognised as a query, expression, data type or
// number is a name; an encoding (ENCODE <encode_type>) is a name too.
func placeholderClass(name string) string {
	n := strings.ReplaceAll(strings.ToLower(name), "-", "_")
	switch {
	case cteName.MatchString(n):
		return classCTE
	case n == "argmode":
		return classArgMode
	case n == "schema_element":
		return classElement
	case queryName.MatchString(n):
		return docs.ClassQuery
	case strings.Contains(n, "encod"):
		return docs.ClassName
	case typeName.MatchString(n) && !strings.HasSuffix(n, "_name"):
		return docs.ClassDataType
	case exprName.MatchString(n):
		return docs.ClassExpr
	case strings.Contains(n, "name"):
		return docs.ClassName
	case numberName.MatchString(n):
		return docs.ClassNumber
	}
	return docs.ClassName
}

// statementWords returns the words a statement form on a page can start
// with: the first word of its slug, the DDL verbs, WITH, and START and
// REFRESH, which BEGIN's and the materialized view pages' forms start
// with.
func statementWords(slug string) map[string]bool {
	first := strings.FieldsFunc(strings.TrimPrefix(slug, "r_"), func(r rune) bool { return r == '_' || r == '-' })
	words := map[string]bool{"WITH": true, "START": true, "REFRESH": true}
	if len(first) > 0 {
		words[strings.ToUpper(first[0])] = true
	}
	for _, w := range []string{"CREATE", "ALTER", "DROP", "SHOW", "DESC", "GRANT", "REVOKE", "ATTACH", "DETACH"} {
		words[w] = true
	}
	return words
}
