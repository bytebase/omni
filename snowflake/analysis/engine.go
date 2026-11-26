package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/snowflake/diagnostics"
	"github.com/bytebase/omni/snowflake/parser"
)

// Engine is Snowflake's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter        = Engine
	_ analysis.Parser          = Engine
	_ analysis.StatementRanger = Engine
	_ analysis.Diagnoser       = Engine
	_ analysis.QueryValidator  = Engine
	_ analysis.Explainer       = Engine
	_ analysis.ResultLimiter   = Engine
)

// Split returns SplitStatements: empty segments are dropped, and a segment's
// text ends at the semicolon that follows it.
func (engine) Split(script string) []analysis.Statement {
	return SplitStatements(script)
}

// Parse parses each segment Split returns, which holds no empty one, into the
// segment's first node. A syntax error is positioned in the script, and its
// message names that position.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	parsed := []analysis.ParsedStatement{}
	for _, segment := range e.Split(script) {
		file, err := parser.Parse(segment.Text)
		if err != nil {
			return nil, segmentSyntaxError(script, segment, err)
		}
		if file == nil || len(file.Stmts) == 0 {
			return nil, errors.New("statement yielded no parse tree")
		}
		parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: file.Stmts[0]})
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
		Message:    fmt.Sprintf("Syntax error at line %d:%d \n%s", position.Line, position.Column, parseErr.Message),
		RawMessage: parseErr.Message,
	}
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	return analysis.LSPRanges(script, StatementRanges(script))
}

func (engine) Diagnose(script string) []analysis.Diagnostic {
	found := diagnostics.Analyze(script)
	result := make([]analysis.Diagnostic, 0, len(found))
	for _, d := range found {
		result = append(result, analysis.Diagnostic{
			Start:    analysis.PositionAt(script, d.Range.Start.Offset),
			Severity: analysis.SeverityError,
			Message:  d.Message,
		})
	}
	return result
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	return ValidateQuery(script)
}

// Explain writes EXPLAIN USING JSON for the JSON format and the tabular EXPLAIN
// for any other. Its errors are the parser's own.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	return ExplainStatement(statement, format == analysis.ExplainFormatJSON)
}

// ResultLimit rewrites the statement's own LIMIT, or wraps a statement it
// cannot rewrite in a subquery.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := RewriteResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("SELECT * FROM (%s) LIMIT %d", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}
