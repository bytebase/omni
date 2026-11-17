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

	// Registers parser.ParseDoBlockFn; without it the SQL parser panics on a DO block.
	_ "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
)

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
