package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/parser"
	"github.com/bytebase/omni/metadata"
)

// Engine is Cassandra's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter          = Engine
	_ analysis.Parser            = Engine
	_ analysis.QueryValidator    = Engine
	_ analysis.QuerySpanResolver = Engine
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
			return nil, syntaxError(segment.Text, err, analysis.PositionAt(script, segment.Start))
		}
		for _, statement := range statements {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: statement.AST})
		}
	}
	return parsed, nil
}

// syntaxError positions a parse error of text in the script that holds text
// at start; any other error is returned as is.
func syntaxError(text string, err error, start analysis.Position) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(text, parseErr.Position).Shift(start)
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

// QuerySpans resolves each statement in the default keyspace, expanding a star
// from the keyspace's synced metadata as stored. A syntax error is positioned
// in its statement's own text.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	opts := ResolveOptions{DefaultKeyspace: request.DefaultDatabase}
	if load := request.Catalog.MetadataLoader(); load != nil {
		opts.GetDatabaseMetadata = func(keyspace string) (*metadata.DatabaseSchemaMetadata, error) {
			return load(ctx, keyspace)
		}
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, opts)
		if err != nil {
			return nil, syntaxError(statement.Text, err, analysis.Position{Line: 1, Column: 1})
		}
		spans = append(spans, span)
	}
	return spans, nil
}
