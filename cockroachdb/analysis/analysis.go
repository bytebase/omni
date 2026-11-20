// Package analysis is CockroachDB's statement judgment: the statement types a
// script changes, the tables it writes, and its system catalogs. The wire
// protocol is PostgreSQL's, but CockroachDB's grammar is not: UPSERT,
// statement sources, regional DDL, zone configuration and IMPORT INTO have no
// PostgreSQL node, so the judgment walks CockroachDB's own AST.
package analysis

import (
	"strings"
	"unicode"

	crrawparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	// Registers parser.ParseDoBlockFn; without it the SQL parser panics on a DO block.
	_ "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
)

// ParsedStatement is an opaque parsed CockroachDB statement that keeps
// cockroachdb-parser AST types encapsulated inside this package.
type ParsedStatement struct {
	stmt statements.Statement[tree.Statement]
}

// Text returns the statement SQL as CockroachDB's parser read it.
func (s ParsedStatement) Text() string {
	return s.stmt.SQL
}

// ParseScript parses script, the whole script or one statement, into opaque
// ParsedStatement values.
func ParseScript(script string) ([]ParsedStatement, error) {
	stmts, err := Parse(script)
	if err != nil {
		return nil, err
	}
	out := make([]ParsedStatement, len(stmts))
	for i, stmt := range stmts {
		out[i] = ParsedStatement{stmt: stmt}
	}
	return out, nil
}

// SplitScript splits script lexically and returns each non-empty statement
// terminated with a semicolon. Unlike Split, which keeps empty segments as "",
// SplitScript omits empty segments; unlike ParseScript, it succeeds on
// syntactically invalid scripts as long as statement boundaries can be scanned.
func SplitScript(script string) []string {
	segments := Split(script)
	result := make([]string, 0, len(segments))
	for _, seg := range segments {
		if seg == "" {
			continue
		}
		if !strings.HasSuffix(seg, ";") {
			seg += ";"
		}
		result = append(result, seg)
	}
	return result
}

func unwrapStatements(stmts []ParsedStatement) statements.Statements {
	out := make(statements.Statements, len(stmts))
	for i, stmt := range stmts {
		out[i] = stmt.stmt
	}
	return out
}

// Parse parses script, the whole script or one statement, into statements
// carrying the AST and the text CockroachDB's parser read. A trailing
// semicolon is written so a script without one parses as CockroachDB's
// parser expects.
func Parse(script string) (statements.Statements, error) {
	script = strings.TrimRightFunc(script, isSpaceOrSemicolon) + ";"
	return crrawparser.Parse(script)
}

// Split returns each statement of script with its terminating semicolon,
// lexically, so a script the parser rejects still splits. A segment that holds
// no statement, such as a lone ";", is "" so a caller can account for every
// boundary.
func Split(script string) []string {
	var result []string
	for {
		pos, ok := crrawparser.SplitFirstStatement(script)
		if !ok {
			break
		}
		segment := strings.TrimLeftFunc(script[:pos], unicode.IsSpace)
		if !holdsStatement(segment) {
			segment = ""
		}
		result = append(result, segment)
		script = script[pos:]
	}
	if segment := strings.TrimLeftFunc(script, unicode.IsSpace); holdsStatement(segment) {
		result = append(result, segment)
	}
	return result
}

func holdsStatement(segment string) bool {
	if strings.TrimSpace(segment) == "" {
		return false
	}
	tokens, ok := crrawparser.Tokens(segment)
	if !ok {
		// A lexical error still leaves text worth reporting as one segment.
		return true
	}
	for _, token := range tokens {
		if token.Str != ";" {
			return true
		}
	}
	return false
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
