// Package completion completes CockroachDB SQL with CockroachDB's own
// grammar, read from the vendor goyacc parser rather than PostgreSQL's.
//
// The vendor parser keeps its LALR tables private, so completion asks it
// questions through its public Parse instead:
//
//   - Keywords: the text before the caret is parsed with each CockroachDB
//     keyword appended. A keyword is offered when the parser shifts it, which
//     an LALR(1) parser does only for a token that can follow the prefix.
//     The prefix is first reduced to a shorter one that leaves the parser in
//     the same state (probePrefix), and probing stops at a time budget.
//   - Names: the word at the caret is replaced by a placeholder identifier
//     and the statement is parsed, patched with a few closing tails when it
//     is incomplete. Where the placeholder lands in CockroachDB's AST says
//     what it names (a table, a column, an index after @) and which tables,
//     aliases and CTEs are in scope. When no patch parses, a lexical scan of
//     the statement supplies the scope.
//   - Built-in functions and type names come from builtins.go, recorded from
//     the pinned CockroachDB server.
package completion

import (
	"context"
	"slices"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
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
	CandidateSequence         = analysis.CandidateSequence         // sequence name
	CandidateIndex            = analysis.CandidateIndex            // index name
	CandidateType_            = analysis.CandidateType_            // type name
)

// Candidate is a single completion suggestion.
type Candidate struct {
	Text       string        // the completion text
	Type       CandidateType // what kind of object this is
	Definition string        // optional definition/signature
	Comment    string        // optional doc comment
}

// caretWord is the word being typed at the caret: the bytes the placeholder
// replaces and the text typed before the caret.
type caretWord struct {
	start, end int
	typed      string
	quoted     bool
}

// completion is what complete found: the candidates, with function and type
// texts final and other names still to quote, and whether the caret is
// inside double quotes.
type completion struct {
	candidates []Candidate
	quoted     bool // the caret is inside a double-quoted identifier
}

// complete returns the candidates for the statement stmt at byte offset
// caret. cat may be nil.
func complete(ctx context.Context, stmt string, caret int, cat *catalog) completion {
	caret = max(0, min(caret, len(stmt)))
	word, ok := caretWordAt(stmt, caret)
	if !ok {
		return completion{}
	}
	prefix := stmt[:word.start]

	nc := nameContextAt(stmt, word, cat)
	var out []Candidate
	if !word.quoted && !afterQualifier(prefix) {
		objectName := nc != nil && nc.kind == kindObject
		for _, kw := range keywordCandidates(ctx, prefix, objectName) {
			out = append(out, Candidate{Text: kw, Type: CandidateKeyword})
		}
	}
	names := nameCandidates(nc, cat, word.typed != "")
	if nc != nil && nc.kind == kindType {
		// A type keyword such as INT8 is offered once, as a type.
		out = slices.DeleteFunc(out, func(c Candidate) bool {
			return slices.ContainsFunc(names, func(n Candidate) bool { return strings.EqualFold(n.Text, c.Text) })
		})
	}
	out = filterByPrefix(append(out, names...), word.typed)
	call := nc != nil && nc.kind == kindExpr && !followedByParen(stmt[word.end:])
	return completion{
		candidates: formatNames(prefix, out, word.quoted, call),
		quoted:     word.quoted,
	}
}

// caretWordAt finds the word at caret. It fails when the caret is inside a
// string, a comment, or after a lexical error, where nothing completes.
func caretWordAt(stmt string, caret int) (caretWord, bool) {
	prefix := stmt[:caret]
	toks := tokenize(prefix)
	if n := len(toks); n > 0 {
		last := toks[n-1]
		if last.id == lexbase.ERROR {
			if !last.quotedIdent {
				return caretWord{}, false
			}
			return caretWord{start: last.start, end: quotedIdentEnd(stmt, caret), typed: last.str, quoted: true}, true
		}
		if last.end == caret && isIdentStart(stmt[last.start]) {
			end := caret
			for end < len(stmt) && isIdentMiddle(stmt[end]) {
				end++
			}
			return caretWord{start: last.start, end: end, typed: prefix[last.start:]}, true
		}
	}
	// A word typed inside a -- comment adds no token.
	if len(tokenize(prefix+" x")) == len(toks) {
		return caretWord{}, false
	}
	end := caret
	for end < len(stmt) && isIdentMiddle(stmt[end]) {
		end++
	}
	return caretWord{start: caret, end: end}, true
}

// quotedIdentEnd returns the end of the double-quoted identifier the caret
// is inside: just past its closing quote, or the end of stmt.
func quotedIdentEnd(stmt string, caret int) int {
	for i := caret; i < len(stmt); i++ {
		if stmt[i] != '"' {
			continue
		}
		if i+1 < len(stmt) && stmt[i+1] == '"' {
			i++
			continue
		}
		return i + 1
	}
	return len(stmt)
}

func afterQualifier(prefix string) bool {
	toks := tokenize(prefix)
	if len(toks) == 0 {
		return false
	}
	last := toks[len(toks)-1]
	return last.id == '.' || last.id == '@'
}

func followedByParen(rest string) bool {
	rest = strings.TrimLeft(rest, " \t\r\n")
	return strings.HasPrefix(rest, "(")
}

// isIdentStart and isIdentMiddle follow CockroachDB's scanner, which takes
// every non-ASCII byte as part of an identifier.
func isIdentStart(c byte) bool {
	return lexbase.IsIdentStart(int(c)) || c >= 0x80
}

func isIdentMiddle(c byte) bool {
	return lexbase.IsIdentMiddle(int(c)) || c >= 0x80
}

// filterByPrefix keeps the candidates whose text starts with typed, ignoring
// case.
func filterByPrefix(candidates []Candidate, typed string) []Candidate {
	if typed == "" {
		return candidates
	}
	upper := strings.ToUpper(typed)
	var out []Candidate
	for _, c := range candidates {
		if strings.HasPrefix(strings.ToUpper(c.Text), upper) {
			out = append(out, c)
		}
	}
	return out
}
