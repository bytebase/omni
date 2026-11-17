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

// Split returns the text of each statement of script, with a trailing
// semicolon on each.
func Split(script string) ([]string, error) {
	parsed, err := Parse(script)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(parsed))
	for _, stmt := range parsed {
		text := stmt.SQL
		if !strings.HasSuffix(text, ";") {
			text += ";"
		}
		result = append(result, text)
	}
	return result, nil
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}
