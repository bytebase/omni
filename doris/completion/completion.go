// Package completion completes Doris SQL with Doris's own grammar, read from
// omni's Doris parser.
//
//   - Keywords: the text before the caret word is parsed once per keyword
//     with that keyword appended (parser.ExpectedKeywords). A keyword the
//     parser takes there is offered.
//   - Names: the caret word is replaced by a placeholder identifier and the
//     statement is parsed, closed by a few tails when it is incomplete.
//     Where the placeholder lands in the AST says whether a table or a
//     column is named, and the query's FROM items say which tables are in
//     scope. The names come from the synced metadata.
//
// The parser keeps a parenthesized query as raw text, so a caret inside one
// is completed on that query alone.
package completion

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/doris/parser"
)

// CandidateType classifies a completion candidate.
type CandidateType = analysis.CandidateType

const (
	CandidateKeyword          = analysis.CandidateKeyword          // SQL keyword
	CandidateDatabase         = analysis.CandidateDatabase         // database name
	CandidateTable            = analysis.CandidateTable            // table name
	CandidateView             = analysis.CandidateView             // view name
	CandidateMaterializedView = analysis.CandidateMaterializedView // materialized view name
	CandidateColumn           = analysis.CandidateColumn           // column name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition: a column's type, a table's qualified name
	Comment    string        // optional comment
}

// caretWord is the word being typed at the caret: the bytes the placeholder
// replaces and the text typed before the caret, without its backtick.
type caretWord struct {
	start, end int
	typed      string
	quoted     bool
}

// complete returns the candidates for the statement stmt at byte offset
// caret. cat may be nil.
func complete(stmt string, caret int, cat *catalog) []Candidate {
	caret = max(0, min(caret, len(stmt)))
	stmt, caret = innermostQuery(stmt, caret)
	word, ok := caretWordAt(stmt, caret)
	if !ok {
		return nil
	}
	prefix := stmt[:word.start]
	var out []Candidate
	if !word.quoted && lastTokenKind(prefix) != '.' {
		for _, kw := range parser.ExpectedKeywords(prefix) {
			out = append(out, Candidate{Text: kw, Type: CandidateKeyword})
		}
	}
	out = filterByPrefix(out, word.typed)
	for _, c := range filterByPrefix(nameCandidates(nameContextAt(stmt, word), cat), word.typed) {
		c.Text = QuoteIdentifierIfNeeded(c.Text, word.quoted)
		out = append(out, c)
	}
	return out
}

// innermostQuery narrows stmt to the innermost query the parser keeps as raw
// text that holds caret: a parenthesized query, the text after its '(' up
// to its ')' or the end of stmt, where SELECT or WITH follows the '('; or
// the query of CREATE ... AS SELECT.
func innermostQuery(stmt string, caret int) (string, int) {
	toks, _ := parser.Tokenize(stmt)
	type group struct {
		bodyStart int
		query     bool
	}
	var open []group
	// queryStart is where the query of CREATE ... AS SELECT starts, which
	// the parser keeps as raw text too.
	queryStart := -1
	createStmt := len(toks) > 0 && isKeyword(toks[0]) && strings.EqualFold(toks[0].Str, "CREATE")
	for i, t := range toks {
		if createStmt && queryStart < 0 && len(open) == 0 && isKeyword(t) && strings.EqualFold(t.Str, "AS") {
			if name := parser.TokenName(toks[i+1].Kind); (name == "SELECT" || name == "WITH") && toks[i+1].Loc.Start < caret {
				queryStart = toks[i+1].Loc.Start
			}
		}
		if t.Loc.Start >= caret || parser.TokenName(t.Kind) == "EOF" {
			break
		}
		switch t.Kind {
		case '(':
			next := toks[i+1]
			name := parser.TokenName(next.Kind)
			open = append(open, group{bodyStart: t.Loc.End, query: name == "SELECT" || name == "WITH"})
		case ')':
			if len(open) > 0 {
				open = open[:len(open)-1]
			}
		}
	}
	for i := len(open) - 1; i >= 0; i-- {
		if !open[i].query {
			continue
		}
		if open[i].bodyStart < queryStart {
			break // the CTAS query holds this group
		}
		start := open[i].bodyStart
		end := len(stmt)
		depth := len(open) - 1 - i // the groups inside the query still open at the caret
		for _, t := range toks {
			if t.Loc.Start < caret {
				continue
			}
			if t.Kind == '(' {
				depth++
			} else if t.Kind == ')' {
				if depth == 0 {
					end = t.Loc.Start
					break
				}
				depth--
			}
		}
		return stmt[start:end], caret - start
	}
	if queryStart >= 0 {
		return stmt[queryStart:], caret - queryStart
	}
	return stmt, caret
}

// caretWordAt finds the word at caret. It fails when the caret is inside a
// string or a comment, where nothing completes.
func caretWordAt(stmt string, caret int) (caretWord, bool) {
	prefix := stmt[:caret]
	toks, _ := parser.Tokenize(prefix)
	toks = toks[:len(toks)-1] // EOF
	if n := len(toks); n > 0 && toks[n-1].Loc.End == caret {
		last := toks[n-1]
		switch name := parser.TokenName(last.Kind); {
		case name == "IDENT" || isKeyword(last):
			return caretWord{start: last.Loc.Start, end: wordEnd(stmt, caret), typed: prefix[last.Loc.Start:]}, true
		case name == "INVALID" && prefix[last.Loc.Start] == '`':
			typed := strings.ReplaceAll(prefix[last.Loc.Start+1:], "``", "`")
			return caretWord{start: last.Loc.Start, end: quotedEnd(stmt, caret), typed: typed, quoted: true}, true
		case name == "INVALID":
			// An unterminated string.
			return caretWord{}, false
		}
	}
	// A word typed inside a comment adds no token.
	if more, _ := parser.Tokenize(prefix + " x"); len(more) == len(toks)+1 {
		return caretWord{}, false
	}
	return caretWord{start: caret, end: wordEnd(stmt, caret)}, true
}

// isKeyword reports whether t is a keyword written as a word.
func isKeyword(t parser.Token) bool {
	kind, ok := parser.KeywordToken(t.Str)
	return ok && kind == t.Kind
}

// wordEnd returns the end of the unquoted word that continues at i.
func wordEnd(stmt string, i int) int {
	for i < len(stmt) && isIdentCont(stmt[i]) {
		i++
	}
	return i
}

// isIdentCont follows the lexer: an unquoted identifier continues with
// letters, digits, _, $ and any non-ASCII byte.
func isIdentCont(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '$' || c > 127
}

// quotedEnd returns the end of the backtick-quoted identifier the caret is
// inside: just past its closing backtick, or the end of its line.
func quotedEnd(stmt string, caret int) int {
	for i := caret; i < len(stmt); i++ {
		switch stmt[i] {
		case '\n':
			return i
		case '`':
			if i+1 < len(stmt) && stmt[i+1] == '`' {
				i++
				continue
			}
			return i + 1
		}
	}
	return len(stmt)
}

// lastTokenKind returns the kind of the last token of text, 0 when none.
func lastTokenKind(text string) int {
	toks, _ := parser.Tokenize(text)
	if len(toks) < 2 {
		return 0
	}
	return toks[len(toks)-2].Kind
}

// filterByPrefix keeps the candidates whose text starts with typed,
// ignoring case.
func filterByPrefix(candidates []Candidate, typed string) []Candidate {
	if typed == "" {
		return candidates
	}
	upper := strings.ToUpper(typed)
	return slices.DeleteFunc(candidates, func(c Candidate) bool {
		return !strings.HasPrefix(strings.ToUpper(c.Text), upper)
	})
}
