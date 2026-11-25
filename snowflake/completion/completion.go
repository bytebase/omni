// Package completion completes Snowflake SQL with Snowflake's own grammar,
// read from omni's Snowflake parser.
//
//   - Keywords: the parser reads the statement up to the caret word in
//     collection mode, recording every keyword it tests the end of that
//     text against (parser.ExpectedKeywords). One parse answers the whole
//     keyword set.
//   - Names: the caret word is replaced by a placeholder identifier and the
//     statement is parsed, closed by a few tails when it is incomplete.
//     Where the placeholder lands in the AST says what it names (a table, a
//     column, a function) and which tables, aliases and CTEs are in scope.
//   - Built-in functions come from functions.tsv, recorded from Snowflake's
//     function index.
package completion

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType = analysis.CandidateType

const (
	CandidateKeyword          = analysis.CandidateKeyword          // SQL keyword
	CandidateSchema           = analysis.CandidateSchema           // schema name
	CandidateTable            = analysis.CandidateTable            // table name
	CandidateView             = analysis.CandidateView             // view name
	CandidateMaterializedView = analysis.CandidateMaterializedView // materialized view name
	CandidateColumn           = analysis.CandidateColumn           // column name
	CandidateFunction         = analysis.CandidateFunction         // function name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// caretWord is the word being typed at the caret: the bytes the placeholder
// replaces and the text typed before the caret, unquoted.
type caretWord struct {
	start, end int
	typed      string
	quoted     bool
}

// name is a candidate whose text is a name still to be quoted or, when
// written is set, a built-in function spelled as the docs spell it.
type name struct {
	Candidate
	written bool
	builtin bool
	// qualifier, when set, is the identifier a qualifier candidate offers;
	// Text is derived from it once the caret's quoting is known.
	qualifier *ast.Ident
}

// complete returns the candidates for the statement stmt at byte offset
// caret. cat may be nil.
func complete(stmt string, caret int, cat *catalog) []Candidate {
	caret = max(0, min(caret, len(stmt)))
	word, ok := caretWordAt(stmt, caret)
	if !ok {
		return nil
	}
	prefix := stmt[:word.start]
	var out []Candidate
	switch before := lastTokenType(prefix); before {
	case '@', '$':
		// A stage or a variable: nothing in the catalog names it.
		return nil
	case '.', ':':
		// A qualified name or a semi-structured path step.
	default:
		if !word.quoted {
			for _, kw := range parser.ExpectedKeywords(prefix) {
				out = append(out, Candidate{Text: kw, Type: CandidateKeyword})
			}
		}
	}
	out = filterByPrefix(out, word.typed)
	keywords := make(map[string]bool, len(out))
	for _, c := range out {
		keywords[c.Text] = true
	}
	nc := nameContextAt(stmt, word, cat)
	call := !followedByParen(stmt[word.end:])
	for _, n := range nameCandidates(nc, cat, word.typed) {
		c := n.Candidate
		if n.builtin && keywords[c.Text] {
			// A built-in spelled by a keyword (CAST, CASE) is offered once,
			// as the keyword the grammar reads.
			continue
		}
		switch {
		case n.qualifier != nil && word.quoted:
			c.Text = QuoteIdentifierIfNeeded(n.qualifier.Normalize(), true)
		case n.qualifier != nil:
			c.Text = n.qualifier.String()
		case !n.written:
			c.Text = QuoteIdentifierIfNeeded(c.Text, word.quoted)
		}
		if c.Type == CandidateFunction && call && !word.quoted {
			c.Text += "()"
		}
		out = append(out, c)
	}
	return out
}

// caretWordAt finds the word at caret. It fails when the caret is inside a
// string, a comment, or a variable, where nothing completes.
func caretWordAt(stmt string, caret int) (caretWord, bool) {
	prefix := stmt[:caret]
	toks, _ := parser.Tokenize(prefix)
	toks = toks[:len(toks)-1] // EOF
	if n := len(toks); n > 0 && toks[n-1].Loc.End == caret {
		last := toks[n-1]
		switch {
		case isWordToken(last):
			end := caret
			for end < len(stmt) && isIdentCont(stmt[end]) {
				end++
			}
			return caretWord{start: last.Loc.Start, end: end, typed: prefix[last.Loc.Start:]}, true
		case parser.TokenName(last.Type) == "INVALID" && prefix[last.Loc.Start] == '"':
			typed := strings.ReplaceAll(prefix[last.Loc.Start+1:], `""`, `"`)
			return caretWord{start: last.Loc.Start, end: quotedIdentEnd(stmt, caret), typed: typed, quoted: true}, true
		case parser.TokenName(last.Type) == "INVALID", parser.TokenName(last.Type) == "VARIABLE":
			// An unterminated string, or a $name variable.
			return caretWord{}, false
		}
	}
	// A word typed inside a comment adds no token.
	if more, _ := parser.Tokenize(prefix + " x"); len(more) == len(toks)+1 {
		return caretWord{}, false
	}
	end := caret
	for end < len(stmt) && isIdentCont(stmt[end]) {
		end++
	}
	return caretWord{start: caret, end: end}, true
}

// isWordToken reports whether t is an unquoted identifier or a keyword.
func isWordToken(t parser.Token) bool {
	if parser.TokenName(t.Type) == "IDENT" {
		return true
	}
	kw, ok := parser.KeywordToken(t.Str)
	return ok && kw == t.Type
}

// isIdentCont follows the lexer: an unquoted identifier continues with
// letters, digits, _, $ and @.
func isIdentCont(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '$' || c == '@'
}

// quotedIdentEnd returns the end of the double-quoted identifier the caret
// is inside: just past its closing quote, or the end of its line.
func quotedIdentEnd(stmt string, caret int) int {
	for i := caret; i < len(stmt); i++ {
		switch stmt[i] {
		case '\n':
			return i
		case '"':
			if i+1 < len(stmt) && stmt[i+1] == '"' {
				i++
				continue
			}
			return i + 1
		}
	}
	return len(stmt)
}

// lastTokenType returns the type of the last token of text, 0 when none. A
// single-character token's type is the character.
func lastTokenType(text string) int {
	toks, _ := parser.Tokenize(text)
	if len(toks) < 2 {
		return 0
	}
	last := toks[len(toks)-2]
	if parser.TokenName(last.Type) == "VARIABLE" {
		return '$'
	}
	return last.Type
}

func followedByParen(rest string) bool {
	return strings.HasPrefix(strings.TrimLeft(rest, " \t\r\n"), "(")
}

// filterByPrefix keeps the candidates whose text starts with typed, ignoring
// case.
func filterByPrefix(candidates []Candidate, typed string) []Candidate {
	if typed == "" {
		return candidates
	}
	upper := strings.ToUpper(typed)
	return slices.DeleteFunc(candidates, func(c Candidate) bool {
		return !strings.HasPrefix(strings.ToUpper(c.Text), upper)
	})
}

// nameCandidates returns the names nc offers that start with typed.
func nameCandidates(nc *nameContext, cat *catalog, typed string) []name {
	if nc == nil {
		return nil
	}
	upper := strings.ToUpper(typed)
	var out []name
	add := func(n name) {
		match := n.Text
		if n.qualifier != nil {
			match = n.qualifier.Name
		}
		if strings.HasPrefix(strings.ToUpper(match), upper) {
			out = append(out, n)
		}
	}
	addColumns := func(items []*rangeItem) {
		for _, item := range items {
			for _, c := range item.columns() {
				add(name{Candidate: Candidate{Text: c.name, Type: CandidateColumn, Definition: c.definition, Comment: c.comment}})
			}
		}
	}
	addRelations := func(s *schema) {
		if s == nil {
			return
		}
		for _, r := range s.relations {
			if len(nc.kinds) > 0 && !slices.Contains(nc.kinds, r.kind) {
				continue
			}
			add(name{Candidate: Candidate{Text: r.name, Type: r.kind, Definition: r.qualifiedName(), Comment: r.comment}})
		}
	}
	addSchemas := func(db *database) {
		if db == nil {
			return
		}
		for _, s := range db.schemas {
			add(name{Candidate: Candidate{Text: s.name, Type: CandidateSchema}})
		}
	}
	addFunctions := func(s *schema, builtins bool) {
		if builtins {
			for _, f := range builtinFunctions() {
				add(name{Candidate: Candidate{Text: f.name, Type: CandidateFunction, Comment: f.comment}, written: true, builtin: true})
			}
		}
		if s != nil {
			for _, f := range s.functions {
				add(name{Candidate: Candidate{Text: f.name, Type: CandidateFunction, Definition: f.definition, Comment: f.comment}})
			}
		}
	}
	// searchPath returns the schemas unqualified names are offered from.
	searchPath := func() []*schema {
		if cat == nil {
			return nil
		}
		var out []*schema
		for _, name := range cat.searchPath(nc.query) {
			if s := cat.database(cat.current).schema(name); s != nil {
				out = append(out, s)
			}
		}
		return out
	}
	// addUnqualified adds what each schema of the search path offers, a
	// name an earlier schema offers hiding a later one's, as it does when
	// Snowflake resolves the name.
	addUnqualified := func(add func(*schema)) {
		before := len(out)
		for _, s := range searchPath() {
			add(s)
		}
		seen := map[string]bool{}
		kept := out[:before]
		for _, n := range out[before:] {
			if !seen[n.Text] {
				seen[n.Text] = true
				kept = append(kept, n)
			}
		}
		out = kept
	}

	switch nc.kind {
	case kindColumnName:
		addColumns(nc.targets)
	case kindRelation:
		switch len(nc.qualifier) {
		case 0:
			for _, cte := range nc.scope.ctes {
				add(name{Candidate: Candidate{Text: cte, Type: CandidateTable}})
			}
			addUnqualified(addRelations)
			if cat != nil {
				addSchemas(cat.database(cat.current))
			}
		case 1:
			// x. is a schema of the current database or another database;
			// only a name that is no schema is loaded as a database.
			if s := cat.schemaAt(nc.qualifier); s != nil {
				addRelations(s)
			} else if cat != nil {
				addSchemas(cat.database(nc.qualifier[0]))
			}
		case 2:
			addRelations(cat.schemaAt(nc.qualifier))
		}
	case kindFunction:
		if len(nc.qualifier) == 0 {
			addFunctions(nil, true)
			addUnqualified(func(s *schema) { addFunctions(s, false) })
		} else if s := cat.schemaAt(nc.qualifier); s != nil {
			addFunctions(s, false)
		} else if len(nc.qualifier) == 1 && cat != nil {
			addSchemas(cat.database(nc.qualifier[0]))
		}
	case kindExpr:
		if len(nc.qualifier) == 0 {
			for _, level := range nc.scope.levels {
				addColumns(level)
				for _, item := range level {
					if !item.qualifier.IsEmpty() {
						add(name{Candidate: Candidate{Type: CandidateTable}, qualifier: &item.qualifier})
					}
				}
			}
			for _, alias := range nc.scope.aliases {
				add(name{Candidate: Candidate{Text: alias, Type: CandidateColumn}})
			}
			addFunctions(nil, true)
			addUnqualified(func(s *schema) { addFunctions(s, false) })
			return out
		}
		// The innermost query whose relation the qualifier names hides the
		// outer ones.
		for _, level := range nc.scope.levels {
			var hit []*rangeItem
			for _, item := range level {
				if item.addressedBy(nc.qualifier) {
					hit = append(hit, item)
				}
			}
			if len(hit) > 0 {
				addColumns(hit)
				return out
			}
		}
		if s := cat.schemaAt(nc.qualifier); s != nil {
			addRelations(s)
			addFunctions(s, false)
		} else if len(nc.qualifier) == 1 && cat != nil {
			// db. in db.schema.f() or db.schema.t.c.
			addSchemas(cat.database(nc.qualifier[0]))
		}
	}
	return out
}
