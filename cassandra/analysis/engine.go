package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/parser"
)

// Engine is Cassandra's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Parser         = Engine
	_ analysis.QueryValidator = Engine
)

// Split cuts the script lexically; a segment's text excludes its semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := cassandra.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  s.Text,
			Empty: s.Empty,
			Start: s.ByteStart,
			End:   s.ByteEnd,
		})
	}
	return statements
}

// Parse parses each non-empty segment and returns one statement per node it
// holds, each carrying its segment. A syntax error is positioned in the script.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		statements, err := cassandra.Parse(segment.Text)
		if err != nil {
			return nil, segmentSyntaxError(script, segment, err)
		}
		for _, statement := range statements {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: statement.AST})
		}
	}
	return parsed, nil
}

func segmentSyntaxError(script string, segment analysis.Statement, err error) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(segment.Text, parseErr.Position).Shift(analysis.PositionAt(script, segment.Start))
	return &analysis.SyntaxError{
		Position:   position,
		Message:    fmt.Sprintf("Syntax error at line %d: %s", position.Line, parseErr.Message),
		RawMessage: parseErr.Message,
	}
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}
