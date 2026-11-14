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
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

// Token classes of Snowflake's own placeholders, beside the ones every
// docs grammar has (docs.ClassName and the rest).
const (
	tName        = docs.ClassName
	tString      = docs.ClassString
	tNumber      = docs.ClassNumber
	tExpr        = docs.ClassExpr
	tQuery       = docs.ClassQuery
	tType        = docs.ClassDataType
	tQName       = docs.ClassQName
	tQName2      = docs.ClassQName2
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
// Keywords are literals. ID, QNAME, QNAME2, STRING, NUMBER, EXPR, QUERY,
// COPY_STATEMENT, DATA_TYPE, STAGE, FILE, POSITION, VARIABLE and OBJECT_TYPES
// are the values the docs write as <placeholders>; tokens.tsv
// samples them. A
// form a page writes with ... for text it leaves out, or with a word it
// never defines, is not converted and is listed under its page.
`,
	Site:           docsSite,
	Index:          indexPath,
	Pages:          commandPages,
	Blocks:         syntaxBlocks,
	Prepare:        prepare,
	Classify:       placeholderClass,
	Value:          value,
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
	syntaxSecs = regexp.MustCompile(`data-section="(?:variant-)?syntax[a-z0-9-]*"`)
)

// syntaxBlocks returns the text of every code block in a page's Syntax
// sections: Syntax, or Syntax for ... where a page has several, and Variant
// syntax, as CREATE OR ALTER TABLE has, each from
// its heading to the next second-level heading, so the subsections some
// pages split their syntax into are included. A page with none, such as an
// overview, returns none.
func syntaxBlocks(doc string) []string {
	var out []string
	end := 0
	for _, loc := range syntaxSecs.FindAllStringIndex(doc, -1) {
		if loc[0] < end && !strings.HasPrefix(doc[loc[1]:], "><h2") {
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
	numberName     = regexp.MustCompile(`(^|_)(num|number|integer|int|n|rows|seconds|minutes|hours|days|bytes|count|percent|ptus|buckets|interval|ms|threshold|width|gb|hour)(_?[0-9]+|_|$)|^(initial|start|min|max|step)_value$|^increment$`)
	queryName      = regexp.MustCompile(`^(query|subquery|select_statement|statement|dml_statement|sql|sql_statement)$`)
	typeName       = regexp.MustCompile(`(^|_)(data_type|col[0-9]*_type|arg_type|structured_type)(_?[0-9]+|_|$)|^type$|^arg_data_type`)
	runtimeVersion = regexp.MustCompile(`(python|java|scala|jdk|runtime)[a-z_]*_version$|^runtime_version$`)
	stringName     = regexp.MustCompile(`(^|_|-)(string|literal|text)(_?[0-9]+|_|-|$)`)
)

// qualifiedName is a placeholder for a schema object's name, which
// db.schema.name may qualify: <table_name>, <masking_policy_name>.
var qualifiedName = regexp.MustCompile(`(^|_)(table|view|stage|stream|task|pipe|policy|tag|function|procedure|seq|sequence|model|secret|service|alert|format|repository|notebook|streamlit|dataset|experiment|gateway|bundle|project|semantic)(_|$)`)

// objectWord is a placeholder that ends in a schema object's kind, without
// saying name: <target_table>, <additional_tables>, <source_schema>.
var objectWord = regexp.MustCompile(`(^|_)(table|view|stage|stream|task|alert|pipe|sequence)s?$`)

// databaseScoped is a placeholder for an object a database holds, whose
// name db.name may qualify: <schema_name>, <database_role_name>.
var databaseScoped = regexp.MustCompile(`(^|_)(schema|database_role)(_name)?$`)

// databaseScopedPage is a page for such an object, whose <name> is QNAME2.
var databaseScopedPage = regexp.MustCompile(`^[a-z]+-(schema|database-role)$`)

// singleName is a placeholder for a name of one part, as a column's is,
// whatever object its other words name.
var singleName = regexp.MustCompile(`(^|_)(col[0-9]*|column|arg|alias|param|property)(_|$)`)

// nameClass is a name placeholder's class: QNAME for a schema object's
// name, and for <name> and <object_name>, which the converter narrows to
// ID on an account-level object's page; ID for the rest.
func nameClass(n string) string {
	switch {
	case singleName.MatchString(n):
		return tName
	case databaseScoped.MatchString(n):
		return tQName2
	case n == "name" || n == "new_name" || n == "object_name" || qualifiedName.MatchString(n):
		return tQName
	}
	return tName
}

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
	case strings.HasPrefix(n, "json_"):
		// SAML_IDENTITY_PROVIDER = <json_object>: JSON is written quoted.
		return tString
	case strings.Contains(n, "public_key") || strings.HasPrefix(n, "path_to_"):
		// RSA_PUBLIC_KEY = '<key>'; USING '@stage/path': quoted, as the
		// examples write them.
		return tString
	case strings.HasSuffix(n, "password"):
		// ADMIN_PASSWORD = '<user_password>': a password is a string.
		return tString
	case n == "arg":
		// CALL p( [ <arg_name> => ] <arg> ): an argument is an expression.
		return tExpr
	case n == "body":
		// A policy's -> <body> is a SQL expression.
		return tExpr
	case n == "col_default" || n == "col_write_default" || n == "default_value":
		// DEFAULT 100: a default is an expression.
		return tExpr
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
	case strings.HasSuffix(n, "_name"):
		// <text_column_name> names a column, whatever its first words say.
		return nameClass(n)
	case stringName.MatchString(n):
		return tString
	case strings.Contains(n, "name") || strings.Contains(n, "identifier"):
		return nameClass(n)
	case numberName.MatchString(n):
		return tNumber
	}
	if objectWord.MatchString(n) && !singleName.MatchString(n) {
		// <target_table>, <source_stream>: a schema object, a QNAME,
		// though it says no name.
		return tQName
	}
	return tName
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

// timeTravelOffset is Time Travel's OFFSET, a negative number of seconds,
// and timeTravelStatement its STATEMENT, a query ID, which is a string.
var (
	timeTravelOffset    = regexp.MustCompile(`OFFSET => <time_difference>`)
	timeTravelStatement = regexp.MustCompile(`STATEMENT => <(id|query_id)>`)
)

// nameArray is a property set to an array of column names, as a model
// monitor's ID_COLUMNS = <id_column_name_array> is: ( 'a' [ , ... ] ).
var nameArray = regexp.MustCompile(`= <([a-z_-]+_array)>`)

// constraintProperties is what [ <constraint_properties> ] stands for, from
// https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint.
const constraintProperties = "[ [ NOT ] ENFORCED ] [ [ NOT ] DEFERRABLE ] [ INITIALLY { DEFERRED | IMMEDIATE } ] [ { ENABLE | DISABLE } ] [ { VALIDATE | NOVALIDATE } ] [ { RELY | NORELY } ]"

// moreColumns finds a table's column definition and the line after it
// that abbreviates the next columns, [ , <col_name> <col_type> [ ... ] ]:
// every column takes the whole definition.
var moreColumns = regexp.MustCompile(`(?s)(--\s*Column definition[^\n]*\n)(.*?)(\n[ \t]*\n[ \t]*-- Additional column definitions[^\n]*\n[ \t]*)\[ , <col_name> <col_type>[^\n]*\]`)

// outOfLineKeys is the column list of an out-of-line UNIQUE, PRIMARY KEY
// or FOREIGN KEY constraint, which the pages write optional because their
// inline constraints share the syntax; out of line, it is required, as
// https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint
// says.
var outOfLineKeys = regexp.MustCompile(`\[ \( <col_name> \[ , <col_name> , \.\.\. \] \) \]`)

// typeSection starts a section of formatTypeOptions for one file format
// type: -- If TYPE = CSV, or -- If FILE_FORMAT = ( TYPE = CSV ... ).
var typeSection = regexp.MustCompile(`^\s*--\s*If\s+(?:FILE_FORMAT\s*=\s*\(\s*)?TYPE\s*=\s*([A-Z0-9_]+)`)

// typeOptionsRef is a choice of TYPE followed by the options, which the
// pages write once for every type.
var typeOptionsRef = regexp.MustCompile(`TYPE = \{ ([A-Z0-9_ |]+) \} \[ formatTypeOptions \]`)

// typeOptionsName is formatTypeOptions named anywhere else, as ALTER FILE
// FORMAT ... SET [ formatTypeOptions ] does.
var typeOptionsName = regexp.MustCompile(`formatTypeOptions([^_A-Za-z0-9]|$)`)

// typeOptions splits formatTypeOptions ::= into one definition for each of
// its -- If TYPE = X sections, formatTypeOptions_X ::=, and has each
// reference choose a TYPE with its own options:
//
//	TYPE = { CSV | JSON } [ formatTypeOptions ]
//
// reads as { TYPE = CSV [ formatTypeOptions_CSV ] | TYPE = JSON
// [ formatTypeOptions_JSON ] }, since the pages list each option under
// the one type it applies to.
func typeOptions(blocks []string) ([]string, bool) {
	at := -1
	for i, b := range blocks {
		if strings.Contains(b, "formatTypeOptions ::=") {
			at = i
			break
		}
	}
	if at < 0 {
		return blocks, false
	}
	lines := strings.Split(blocks[at], "\n")
	var out []string
	var types []string
	in, done := false, false
	for _, l := range lines {
		switch {
		case done:
			out = append(out, l)
		case strings.Contains(l, "formatTypeOptions ::="):
			in = true
		case !in:
			out = append(out, l)
		case typeSection.MatchString(l):
			t := typeSection.FindStringSubmatch(l)[1]
			types = append(types, t)
			out = append(out, "formatTypeOptions_"+t+" ::=")
		case docs.DefHead.MatchString(l):
			done = true
			out = append(out, l)
		case len(types) == 0:
			if strings.TrimSpace(l) != "" {
				// Options before any section apply to every type: not
				// this shape.
				return blocks, false
			}
		default:
			out = append(out, l)
		}
	}
	if len(types) == 0 {
		return blocks, false
	}
	has := map[string]bool{}
	for _, t := range types {
		has[t] = true
	}
	each := "{ formatTypeOptions_" + strings.Join(types, " | formatTypeOptions_") + " }"
	res := append([]string(nil), blocks...)
	res[at] = strings.Join(out, "\n")
	for i, b := range res {
		if i == at {
			continue
		}
		b = typeOptionsRef.ReplaceAllStringFunc(b, func(m string) string {
			var alts []string
			for _, t := range strings.Split(typeOptionsRef.FindStringSubmatch(m)[1], "|") {
				t = strings.TrimSpace(t)
				if has[t] {
					alts = append(alts, "TYPE = "+t+" [ formatTypeOptions_"+t+" ]")
				} else {
					alts = append(alts, "TYPE = "+t)
				}
			}
			return "{ " + strings.Join(alts, " | ") + " }"
		})
		res[i] = typeOptionsName.ReplaceAllString(b, each+"$1")
	}
	return res, true
}

// privilegeSection starts a section of a privileges definition for the
// objects it names: -- For WAREHOUSE, -- For FILE FORMAT, FUNCTION (UDF or
// external function), or TYPE, -- For { MASKING | ROW ACCESS } POLICY or TAG.
var privilegeSection = regexp.MustCompile(`^\s*--\s*For\s+(.+?)\s*$`)

var (
	parenthetical = regexp.MustCompile(`\([^)]*\)`)
	braceObjects  = regexp.MustCompile(`\{([^}]*)\}\s*([A-Z][A-Z ]*[A-Z])`)
	lowerWord     = regexp.MustCompile(`\b[a-z]+\b`)
)

// sectionObjects lists the object types a section header names.
func sectionObjects(h string) []string {
	h = parenthetical.ReplaceAllString(h, "")
	h = braceObjects.ReplaceAllStringFunc(h, func(m string) string {
		sub := braceObjects.FindStringSubmatch(m)
		var out []string
		for _, a := range strings.Split(sub[1], "|") {
			out = append(out, strings.TrimSpace(a)+" "+sub[2])
		}
		return strings.Join(out, ", ")
	})
	h = strings.ReplaceAll(h, " or ", ", ")
	h = lowerWord.ReplaceAllString(h, "")
	var out []string
	for _, o := range strings.Split(h, ",") {
		o = strings.Join(strings.Fields(o), " ")
		if o != "" {
			out = append(out, o)
		}
	}
	return out
}

// plural is an object type's plural, as ON ALL <object_type_plural> writes
// it: TABLE, TABLES; POLICY, POLICIES.
func plural(o string) string {
	switch {
	case strings.HasSuffix(o, "Y") && !strings.ContainsRune("AEIOU", rune(o[len(o)-2])):
		return o[:len(o)-1] + "IES"
	case strings.HasSuffix(o, "S") || strings.HasSuffix(o, "X"):
		return o + "ES"
	}
	return o + "S"
}

// privilegeSections splits name ::= into one definition for each of its
// -- For X sections, name_1 ::= and on, and returns which of them each
// object type takes.
func privilegeSections(blocks []string, name string) ([]string, map[string][]string, bool) {
	head := name + " ::="
	at := -1
	for i, b := range blocks {
		if strings.Contains(b, head) {
			at = i
			break
		}
	}
	if at < 0 {
		return blocks, nil, false
	}
	takes := map[string][]string{}
	var out []string
	n := 0
	in, done := false, false
	for _, l := range strings.Split(blocks[at], "\n") {
		switch {
		case done:
			out = append(out, l)
		case strings.Contains(l, head):
			in = true
		case !in:
			out = append(out, l)
		case privilegeSection.MatchString(l):
			n++
			sec := fmt.Sprintf("%s_%d", name, n)
			for _, o := range sectionObjects(privilegeSection.FindStringSubmatch(l)[1]) {
				takes[o] = append(takes[o], sec)
			}
			out = append(out, sec+" ::=")
		case docs.DefHead.MatchString(l):
			done = true
			out = append(out, l)
		case n == 0:
			if strings.TrimSpace(l) != "" {
				return blocks, nil, false
			}
		default:
			out = append(out, l)
		}
	}
	if n == 0 {
		return blocks, nil, false
	}
	res := append([]string(nil), blocks...)
	res[at] = strings.Join(out, "\n")
	return res, takes, true
}

// privilegesFor is the privileges an object type takes, or ALL.
func privilegesFor(secs []string) string {
	if len(secs) == 0 {
		return "ALL [ PRIVILEGES ]"
	}
	return "{ " + strings.Join(secs, " | ") + " | ALL [ PRIVILEGES ] }"
}

// sortedObjects lists the object types in a stable order.
func sortedObjects(takes map[string][]string) []string {
	var out []string
	for o := range takes {
		out = append(out, o)
	}
	sort.Strings(out)
	return out
}

// routine holds the object types granted on with their argument types.
var routine = map[string]bool{"FUNCTION": true, "PROCEDURE": true}

// noFuture holds the object types future grants cannot name, from
// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#restrictions-and-limitations.
var noFuture = map[string]bool{
	"COMPUTE POOL": true, "IMAGE REPOSITORY": true, "SNAPSHOT": true, "TAG": true,
	"AGGREGATION POLICY": true, "JOIN POLICY": true, "MAINTENANCE POLICY": true,
	"MASKING POLICY": true, "PACKAGES POLICY": true, "PROJECTION POLICY": true,
	"ROW ACCESS POLICY": true, "SESSION POLICY": true, "STORAGE LIFECYCLE POLICY": true,
}

// noAll holds the object types ON ALL cannot name, from the GRANT page's
// object_type_plural parameter and usage notes.
var noAll = map[string]bool{"MAINTENANCE POLICY": true, "PIPE": true, "TAG": true, "MASKING POLICY": true}

// objectPrivilegesRef grants account-object privileges on a choice of
// object types, as GRANT and REVOKE write it.
var objectPrivilegesRef = regexp.MustCompile(`\{ accountObjectPrivileges\s+\| ALL \[ PRIVILEGES \] \} ON \{ ([A-Z |]+) \} (<[a-z_]+>)`)

// schemaObjectRef grants schema-object privileges on an object or on all
// of a type in a schema; schemaFutureRef on future ones. Each leaves the
// brace that closes it where it is, on its line or the next.
var (
	schemaObjectRef = regexp.MustCompile(`\{ schemaObjectPrivileges\s+\| ALL \[ PRIVILEGES \] \} ON \{ <object_type> <object_name> \| ALL <object_type_plural> IN (\{ DATABASE <db_name> \| SCHEMA <schema_name> \}|SCHEMA <schema_name>)`)
	schemaFutureRef = regexp.MustCompile(`\{ schemaObjectPrivileges\s+\| ALL \[ PRIVILEGES \] \} ON FUTURE <object_type_plural> IN (\{ DATABASE <db_name> \| SCHEMA <schema_name> \}|SCHEMA <schema_name>)`)
)

// objectPrivileges splits accountObjectPrivileges ::= and
// schemaObjectPrivileges ::= into one definition for each of their -- For X
// sections, and has each grant choose an object type with its own
// privileges:
//
//	{ accountObjectPrivileges | ALL [ PRIVILEGES ] } ON { USER | WAREHOUSE } <name>
//
// reads as { { accountObjectPrivileges_1 | ALL [ PRIVILEGES ] } ON USER
// <name> | ... }, and ON <object_type> <object_name> as ON TABLE <name>
// and every other type a section names, since the pages list each
// privilege under the objects it applies to.
func objectPrivileges(blocks []string) ([]string, bool) {
	changed := false
	if res, takes, ok := privilegeSections(blocks, "accountObjectPrivileges"); ok {
		for i, b := range res {
			res[i] = objectPrivilegesRef.ReplaceAllStringFunc(b, func(m string) string {
				sub := objectPrivilegesRef.FindStringSubmatch(m)
				var alts []string
				for _, o := range strings.Split(sub[1], "|") {
					o = strings.TrimSpace(o)
					// An account-level object's name is one part.
					alts = append(alts, privilegesFor(takes[o])+" ON "+o+" <account_object_name>")
				}
				return "{ " + strings.Join(alts, " | ") + " }"
			})
		}
		if !mentions(res, "accountObjectPrivileges") {
			blocks, changed = res, true
		}
	}
	if res, takes, ok := privilegeSections(blocks, "schemaObjectPrivileges"); ok {
		objects := sortedObjects(takes)
		for i, b := range res {
			b = schemaObjectRef.ReplaceAllStringFunc(b, func(m string) string {
				in := schemaObjectRef.FindStringSubmatch(m)[1]
				var alts []string
				for _, o := range objects {
					target := o + " <object_name>"
					if routine[o] {
						// A function or procedure is named with its argument
						// types, which tell overloads apart.
						target += " ( [ <arg_data_type> [ , ... ] ] )"
					}
					if noAll[o] {
						alts = append(alts, privilegesFor(takes[o])+" ON "+target)
						continue
					}
					alts = append(alts, privilegesFor(takes[o])+" ON { "+target+" | ALL "+plural(o)+" IN "+in+" }")
				}
				// The brace the page closes the ON choice with closes this.
				return "{ " + strings.Join(alts, " | ")
			})
			res[i] = schemaFutureRef.ReplaceAllStringFunc(b, func(m string) string {
				in := schemaFutureRef.FindStringSubmatch(m)[1]
				var alts []string
				for _, o := range objects {
					if noFuture[o] {
						continue
					}
					alts = append(alts, privilegesFor(takes[o])+" ON FUTURE "+plural(o)+" IN "+in)
				}
				return "{ " + strings.Join(alts, " | ") + " }"
			})
		}
		if !mentions(res, "schemaObjectPrivileges") {
			blocks, changed = res, true
		}
	}
	return blocks, changed
}

// mentions reports whether a block still names the definition name, which a
// split has replaced: a reference the rewrite did not read.
func mentions(blocks []string, name string) bool {
	for _, b := range blocks {
		for _, l := range strings.Split(b, "\n") {
			if i := strings.Index(l, name); i >= 0 && !strings.HasPrefix(l[i+len(name):], "_") {
				return true
			}
		}
	}
	return false
}

// foreignKeyRef is a foreign key's REFERENCES clause, which takes the
// properties only foreign keys have, from
// https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint.
var foreignKeyRef = regexp.MustCompile(`(REFERENCES <ref_table_name> \[ \( <ref_col_name>[^()]*\) \])`)

const foreignKeyProperties = " [ MATCH { FULL | SIMPLE | PARTIAL } ] [ ON [ UPDATE { CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION } ] [ DELETE { CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION } ] ]"

// prefixList is API_ALLOWED_PREFIXES or API_BLOCKED_PREFIXES written with
// one elided URL, ('<...>'): a list of URLs.
var prefixList = regexp.MustCompile(`(_PREFIXES = )\(\s*'[^'\n]*<\.\.\.>'\s*\)`)

// commaKeywords is a list of keywords written with commas where a choice
// is meant: { READ, WRITE } [ , ... ], or a line MONITOR, OPERATE, SELECT
// [ , ... ]; the [ , ... ] already makes the list.
var (
	commaKeywordBrace = regexp.MustCompile(`\{ ((?:[A-Z][A-Z_ ]*[A-Z_], )+[A-Z][A-Z_ ]*[A-Z_]) \}`)
	commaSep          = regexp.MustCompile(`,[ \t]*`)
	commaKeywordLine  = regexp.MustCompile(`(?m)^([ \t]*)((?:[A-Z][A-Z_ ]*[A-Z_],[ \t]*)+[A-Z][A-Z_ ]*[A-Z_])[ \t]*(\[ , \.\.\. \])[ \t]*$`)
	// temporary is CREATE TABLE's [ { LOCAL | GLOBAL } ] TEMP | TEMPORARY,
	// where the prefix takes either spelling.
	temporary = regexp.MustCompile(`\[ \{ LOCAL \| GLOBAL \} \] TEMP \| TEMPORARY`)
)

// propertyRun is the constraint properties as the pages list them, one
// optional clause after another; propertySet is what they mean, any of
// them in any order, as the constraint page says.
var (
	propertyRun = regexp.MustCompile(`\[ \[ NOT \] ENFORCED \]\s*\[ \[ NOT \] DEFERRABLE \]\s*\[ INITIALLY \{ DEFERRED \| IMMEDIATE \} \]\s*\[ \{ ENABLE \| DISABLE \} \]\s*\[ \{ VALIDATE \| NOVALIDATE \} \]\s*\[ \{ RELY \| NORELY \} \]`)
	propertySet = "[ { [ NOT ] ENFORCED | [ NOT ] DEFERRABLE | INITIALLY { DEFERRED | IMMEDIATE } | { ENABLE | DISABLE } | { VALIDATE | NOVALIDATE } | { RELY | NORELY } } [ ... ] ]"
)

// objectTypeChoice is <object_type>: the account-level objects GRANT names
// and the schema-level ones its object_type parameter lists, from
// https://docs.snowflake.com/en/sql-reference/sql/grant-privilege.
const objectTypeChoice = "('ACCOUNT' | 'APPLICATION' | 'APPLICATION' 'PACKAGE' | 'COMPUTE' 'POOL' | 'CONNECTION' | 'DATA' 'EXCHANGE' 'LISTING' | 'DATABASE' | 'DATABASE' 'ROLE' | 'EXTERNAL' 'VOLUME' | 'FAILOVER' 'GROUP' | 'INTEGRATION' | 'REPLICATION' 'GROUP' | 'RESOURCE' 'MONITOR' | 'ROLE' | 'SCHEMA' | 'SHARE' | 'USER' | 'WAREHOUSE' | 'AGENT' | 'AGGREGATION' 'POLICY' | 'ALERT' | 'AUTHENTICATION' 'POLICY' | 'CORTEX' 'EXTENSION' | 'CORTEX' 'SEARCH' 'SERVICE' | 'DATA' 'METRIC' 'FUNCTION' | 'DATASET' | 'DBT' 'PROJECT' | 'DCM' 'PROJECT' | 'DYNAMIC' 'TABLE' | 'EVENT' 'TABLE' | 'EXPERIMENT' | 'EXTERNAL' 'TABLE' | 'FILE' 'FORMAT' | 'FUNCTION' | 'GATEWAY' | 'GIT' 'REPOSITORY' | 'ICEBERG' 'TABLE' | 'IMAGE' 'REPOSITORY' | 'INTERACTIVE' 'TABLE' | 'JOIN' 'POLICY' | 'MAINTENANCE' 'POLICY' | 'MASKING' 'POLICY' | 'MATERIALIZED' 'VIEW' | 'MCP' 'SERVER' | 'MODEL' | 'MODEL' 'MONITOR' | 'NETWORK' 'RULE' | 'NOTEBOOK' | 'ONLINE' 'FEATURE' 'TABLE' | 'PACKAGES' 'POLICY' | 'PASSWORD' 'POLICY' | 'PIPE' | 'PRIVACY' 'POLICY' | 'PROCEDURE' | 'PROJECTION' 'POLICY' | 'RESTRICTED' 'SESSION' 'SCOPE' | 'ROW' 'ACCESS' 'POLICY' | 'SECRET' | 'SEMANTIC' 'VIEW' | 'SEQUENCE' | 'SERVICE' | 'SESSION' 'POLICY' | 'SNAPSHOT' | 'SNAPSHOT' 'POLICY' | 'SNAPSHOT' 'SET' | 'STAGE' | 'STORAGE' 'LIFECYCLE' 'POLICY' | 'STREAM' | 'STREAMLIT' | 'TABLE' | 'TAG' | 'TASK' | 'TYPE' | 'VIEW' | 'WORKSPACE')"

// templateValue is a template argument, USING ( <key> => <value> ).
var templateValue = regexp.MustCompile(`=> <value>`)

// missingBar finds the typo many SHOW pages share: the IN list's
// <schema_name> line has no | before the APPLICATION lines after it.
var missingBar = regexp.MustCompile(`(<schema_name>)[ \t]*\n((?:[ \t]*\n)*[ \t]*APPLICATION )`)

// prepare reads what Snowflake's pages write in a way of their own, page
// after page, after docs-fixes.tsv's rows.
func prepare(_ string, blocks []string) ([]string, []string) {
	var notes []string
	offset := false
	for i, b := range blocks {
		fixed := timeTravelOffset.ReplaceAllString(b, "OFFSET => -<seconds>")
		fixed = timeTravelStatement.ReplaceAllString(fixed, "STATEMENT => '<query_id>'")
		if fixed != b {
			blocks[i] = fixed
			offset = true
		}
	}
	if offset {
		notes = append(notes, "read: Time Travel's OFFSET is a negative number of seconds, -N, and its STATEMENT a quoted query ID, as the Time Travel parameters say")
	}
	for i, b := range blocks {
		if fixed := strings.ReplaceAll(b, "[ <constraint_properties> ]", constraintProperties); fixed != b {
			blocks[i] = fixed
			notes = append(notes, "read: <constraint_properties> is the properties CREATE TABLE's constraint page lists")
		}
		if fixed := templateValue.ReplaceAllString(blocks[i], "=> <expr>"); fixed != blocks[i] {
			blocks[i] = fixed
			notes = append(notes, "read: a template argument's <value> is a constant, a string or JSON value, as the USING parameters say")
		}
	}
	outOfLine := false
	for i, b := range blocks {
		at := strings.Index(b, "outoflineConstraint ::=")
		if at < 0 {
			continue
		}
		end := len(b)
		if next := strings.Index(b[at+len("outoflineConstraint ::="):], "::="); next >= 0 {
			end = at + len("outoflineConstraint ::=") + next
		}
		def := outOfLineKeys.ReplaceAllString(b[at:end], "( <col_name> [ , <col_name> , ... ] )")
		if def != b[at:end] {
			blocks[i] = b[:at] + def + b[end:]
			outOfLine = true
		}
	}
	if outOfLine {
		notes = append(notes, "read: an out-of-line key constraint names its columns, as the constraint page says")
	}
	if split, ok := typeOptions(blocks); ok {
		blocks = split
		notes = append(notes, "read: each -- If TYPE = X section of formatTypeOptions is the options of TYPE = X alone, and TYPE = { ... } [ formatTypeOptions ] chooses a TYPE with its own options")
	}
	if split, ok := objectPrivileges(blocks); ok {
		blocks = split
		notes = append(notes, "read: each -- For X section of accountObjectPrivileges and schemaObjectPrivileges is the privileges on X alone, and the ON clause chooses an object type with its own privileges")
	}
	fk, prefixes := false, false
	for i, b := range blocks {
		if !strings.Contains(b, "MATCH {") {
			if fixed := foreignKeyRef.ReplaceAllString(b, "${1}"+foreignKeyProperties); fixed != b {
				blocks[i] = fixed
				fk = true
			}
		}
		if fixed := prefixList.ReplaceAllString(blocks[i], "${1}( '<prefix>' [ , '<prefix>' ... ] )"); fixed != blocks[i] {
			blocks[i] = fixed
			prefixes = true
		}
	}
	if fk {
		notes = append(notes, "read: a foreign key's REFERENCES takes MATCH and ON UPDATE/DELETE, the properties the constraint page gives foreign keys")
	}
	if prefixes {
		notes = append(notes, "read: ('<...>') after _PREFIXES = is a list of URLs, as the parameters say")
	}
	for i, b := range blocks {
		blocks[i] = namespaceRef.ReplaceAllString(b, "${1}[ <ns1> . [ <ns2> . ] ]")
	}
	choices := false
	for i, b := range blocks {
		fixed := commaKeywordBrace.ReplaceAllStringFunc(b, func(m string) string {
			return "{ " + strings.ReplaceAll(commaKeywordBrace.FindStringSubmatch(m)[1], ", ", " | ") + " }"
		})
		fixed = commaKeywordLine.ReplaceAllStringFunc(fixed, func(m string) string {
			sub := commaKeywordLine.FindStringSubmatch(m)
			words := commaSep.Split(sub[2], -1)
			return sub[1] + "{ " + strings.Join(words, " | ") + " } " + sub[3]
		})
		fixed = temporary.ReplaceAllString(fixed, "[ { LOCAL | GLOBAL } ] { TEMP | TEMPORARY }")
		if fixed != b {
			blocks[i] = fixed
			choices = true
		}
	}
	if choices {
		notes = append(notes, "read: keywords written with commas before [ , ... ] are a choice, the list the [ , ... ] makes")
	}
	for i, b := range blocks {
		blocks[i] = propertyRun.ReplaceAllString(b, propertySet)
	}
	columns := false
	for i, b := range blocks {
		if fixed := moreColumns.ReplaceAllString(b, "${1}{ ${2} }${3}[ , ... ]"); fixed != b {
			blocks[i] = fixed
			columns = true
		}
	}
	if columns {
		notes = append(notes, "read: the additional columns take the whole column definition, which their line abbreviates")
	}
	for i, b := range blocks {
		if fixed := nameArray.ReplaceAllString(b, "= ( '<$1>' [ , ... ] )"); fixed != b {
			blocks[i] = fixed
			notes = append(notes, "read: a <..._array> of column names is ( '<name>' [ , ... ] ), as the examples write it")
		}
	}
	for i, b := range blocks {
		if fixed := missingBar.ReplaceAllString(b, "$1 |\n$2"); fixed != b {
			blocks[i] = fixed
			notes = append(notes, "read: the IN list's <schema_name> line has no | before APPLICATION; one is read there")
		}
	}
	return blocks, notes
}

// accountLevelPage is a page for an account-level object, whose <name> is
// one part: a warehouse, user, role, database, share or integration.
var accountLevelPage = regexp.MustCompile(`^[a-z]+-(warehouse|user|role|database|share|connection|resource-monitor|compute-pool|failover-group|replication-group|external-volume|application|application-package|account|organization-[a-z-]+|network-policy|listing|[a-z-]*integration[a-z-]*)$`)

// namespaceRef is [<namespace>.] before a name outside a stage or a
// string: a database, or a database and a schema.
var namespaceRef = regexp.MustCompile(`(^|[^@'])\[\s*<namespace>\.\s*\]`)

// value writes a value as the grammar reads it on a page: <object_type>
// as the object types GRANT lists, and a <name> of an account-level
// object as one part, of a database object as two at most.
func value(slug, placeholder, cls string) string {
	if cls == tObjectType {
		// SHOW GRANTS ON <object_type> <name>: one of the object types the
		// GRANT page lists.
		return objectTypeChoice
	}
	if cls == tQName && (placeholder == "name" || placeholder == "new_name" || placeholder == "object_name") {
		switch {
		case accountLevelPage.MatchString(slug):
			// ALTER WAREHOUSE <name>: an account-level object's name is
			// one part.
			return tName
		case databaseScopedPage.MatchString(slug):
			// CREATE SCHEMA <name>: db.name at most.
			return tQName2
		}
	}
	return cls
}
