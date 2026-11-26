package analysis

import (
	"errors"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/dynamodb"
	"github.com/bytebase/omni/dynamodb/parser"
)

// Engine is DynamoDB PartiQL's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.QueryValidator = Engine
)

// Split cuts the script lexically; a segment's text ends at its own
// semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := dynamodb.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  s.Text,
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   s.ByteEnd,
		})
	}
	return statements
}

// ValidateQuery reports a script of DQL only as read-only and returning rows,
// and anything else it parses as neither. A syntax error's Message is the
// parser's, which names the error's byte offset.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	err := ValidateQuery(script)
	if err == nil {
		return true, true, nil
	}
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return false, false, nil
	}
	return false, false, &analysis.SyntaxError{
		Position:   analysis.PositionAt(script, parseErr.Position),
		Message:    parseErr.Error(),
		RawMessage: parseErr.Message,
	}
}
