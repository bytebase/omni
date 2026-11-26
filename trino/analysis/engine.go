package analysis

import (
	"errors"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/trino/parser"
)

// Engine is Trino's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Parser         = Engine
	_ analysis.Diagnoser      = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.Explainer      = Engine
)

// Split keeps empty segments; a segment's text ends at the semicolon that
// follows it.
func (engine) Split(script string) []analysis.Statement {
	segments := parser.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		end := s.ByteEnd
		if end < len(script) && script[end] == ';' {
			end++
		}
		statements = append(statements, analysis.Statement{
			Text:  script[s.ByteStart:end],
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   end,
		})
	}
	return statements
}

// Parse keeps every segment, an empty one with no node, and parses each other
// one into its first node. A syntax error is positioned in the script.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	parsed := []analysis.ParsedStatement{}
	for _, segment := range e.Split(script) {
		statement := analysis.ParsedStatement{Statement: segment}
		if !segment.Empty {
			file, err := parser.Parse(segment.Text)
			if err != nil {
				return nil, syntaxError(segment.Text, err, analysis.PositionAt(script, segment.Start))
			}
			if file != nil && len(file.Stmts) > 0 {
				statement.Node = file.Stmts[0]
			}
		}
		parsed = append(parsed, statement)
	}
	return parsed, nil
}

// syntaxError converts a *parser.ParseError in text into a *SyntaxError
// shifted to base; any other error is returned as is.
func syntaxError(text string, err error, base analysis.Position) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	return &analysis.SyntaxError{
		Position:   analysis.PositionAt(text, parseErr.Position).Shift(base),
		Message:    parseErr.Error(),
		RawMessage: parseErr.Message,
	}
}

func (engine) Diagnose(script string) []analysis.Diagnostic {
	found := Diagnose(script)
	result := make([]analysis.Diagnostic, 0, len(found))
	for _, d := range found {
		result = append(result, analysis.Diagnostic{
			Start:    analysis.PositionAt(script, d.Loc.Start),
			Severity: analysis.SeverityError,
			Message:  d.Msg,
		})
	}
	return result
}

func (engine) ValidateQuery(script string) (bool, bool, error) {
	readOnly, returnsData, err := ValidateQuery(script)
	if err != nil {
		return false, false, syntaxError(script, err, analysis.Position{Line: 1, Column: 1})
	}
	return readOnly, returnsData, nil
}

// Explain writes EXPLAIN (FORMAT JSON) for the JSON format and the text EXPLAIN
// for any other.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	explain, err := ExplainStatement(statement, format == analysis.ExplainFormatJSON)
	if err != nil && !errors.Is(err, ErrExplainAnalyzeExecutes) {
		return "", syntaxError(statement, err, analysis.Position{Line: 1, Column: 1})
	}
	return explain, err
}
