// Package parser parses CockroachDB SQL statements and normalizes parse
// failures into the repository's ParseError contract.
//
// It delegates to cockroachdb/analysis, which is the single non-driver package
// that links github.com/cockroachdb/cockroachdb-parser, keeping vendor AST and
// error types encapsulated inside cockroachdb/analysis.
package parser

import (
	"fmt"

	"github.com/bytebase/omni/cockroachdb/analysis"
)

// Statement is a parsed CockroachDB SQL statement. It keeps the underlying
// parser AST encapsulated inside cockroachdb/analysis.
type Statement struct {
	stmt analysis.ParsedStatement
}

// Text returns the statement SQL as CockroachDB's parser read it.
func (s Statement) Text() string {
	return s.stmt.Text()
}

// ParseError represents a syntax or lexical error from CockroachDB's parser.
// Position is the byte offset into the parsed SQL, or -1 when CockroachDB
// reports no source location.
type ParseError struct {
	Message  string
	Position int
	Code     string
}

// Error returns the formatted parse error message.
func (e *ParseError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s (SQLSTATE %s)", e.Message, e.Code)
	}
	return e.Message
}

// Parse parses one or more CockroachDB SQL statements. On failure it returns a
// *ParseError carrying Message, byte-offset Position (-1 if unavailable), and
// SQLSTATE Code.
func Parse(sql string) ([]Statement, error) {
	parsed, info, err := analysis.ParseRaw(sql)
	if err != nil {
		if info != nil {
			return nil, &ParseError{
				Message:  info.Message,
				Position: info.Position,
				Code:     info.Code,
			}
		}
		return nil, &ParseError{
			Message:  err.Error(),
			Position: -1,
		}
	}
	out := make([]Statement, len(parsed))
	for i, stmt := range parsed {
		out[i] = Statement{stmt: stmt}
	}
	return out, nil
}
