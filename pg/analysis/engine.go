package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// Engine is PostgreSQL's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter            = Engine
	_ analysis.Parser              = Engine
	_ analysis.StatementClassifier = Engine
	_ analysis.StatementRanger     = Engine
	_ analysis.Diagnoser           = Engine
	_ analysis.QueryValidator      = Engine
	_ analysis.Explainer           = Engine
	_ analysis.ResultLimiter       = Engine
)

// Split cuts the script lexically; a segment's text ends at its own
// semicolon.
func (engine) Split(script string) []analysis.Statement {
	segments := pg.Split(script)
	if len(segments) == 0 {
		return nil
	}
	statements := make([]analysis.Statement, 0, len(segments))
	for _, s := range segments {
		statements = append(statements, analysis.Statement{
			Text:  script[s.ByteStart:s.ByteEnd],
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   s.ByteEnd,
		})
	}
	return statements
}

// Parse parses each non-empty segment and returns one statement per node it
// holds, each carrying its segment. A syntax error is positioned on the
// segment's lines; its column is the segment's own.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		statements, err := pg.Parse(segment.Text)
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
	position := analysis.PositionAt(segment.Text, parseErr.Position)
	position.Line += analysis.PositionAt(script, segment.Start).Line - 1
	return &analysis.SyntaxError{
		Position:   position,
		Message:    fmt.Sprintf("Syntax error at line %d: %s", position.Line, parseErr.Message),
		RawMessage: parseErr.Message,
	}
}

func (engine) StatementTypes(statement analysis.ParsedStatement) []analysis.StatementType {
	node, ok := statement.Node.(ast.Node)
	if !ok || node == nil {
		return nil
	}
	return StatementTypes(node)
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	ranges := StatementRanges(script)
	if ranges == nil {
		return nil
	}
	converted := make([]analysis.LSPRange, 0, len(ranges))
	for _, r := range ranges {
		converted = append(converted, analysis.LSPRange{
			Start: analysis.LSPPosition{Line: r.Start.Line, Character: r.Start.Character},
			End:   analysis.LSPPosition{Line: r.End.Line, Character: r.End.Character},
		})
	}
	return converted
}

// Diagnose reports the script's first syntax error. An error the parser did
// not locate is marked at the script's start.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	diagnostics := []analysis.Diagnostic{}
	err := Diagnose(script)
	if err == nil {
		return diagnostics
	}
	position := analysis.Position{Line: 1, Column: 1}
	if err.ParseError && err.Offset >= 0 {
		position = analysis.PositionAt(err.Text, err.Offset)
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    position,
		Severity: analysis.SeverityError,
		Message:  err.Message,
	})
}

// ValidateQuery returns the parser's own error for a script it rejects.
func (engine) ValidateQuery(script string) (bool, bool, error) {
	statements, err := pg.Parse(script)
	if err != nil {
		return false, false, err
	}
	nodes := make([]ast.Node, 0, len(statements))
	for _, statement := range statements {
		nodes = append(nodes, statement.AST)
	}
	readOnly, returnsData := ValidateQuery(nodes)
	return readOnly, returnsData, nil
}

func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	var name string
	switch format {
	case analysis.ExplainFormatJSON:
		name = "json"
	case analysis.ExplainFormatXML:
		name = "xml"
	case analysis.ExplainFormatYAML:
		name = "yaml"
	default:
	}
	return ExplainStatement(statement, name)
}

// ResultLimit rewrites the statement's own LIMIT, or wraps a statement it
// cannot rewrite in a CTE.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	rewritten, err := RewriteResultLimit(statement, limit)
	if err != nil {
		return fmt.Sprintf("WITH result AS (\n%s\n) SELECT * FROM result LIMIT %d;", analysis.TrimStatement(statement), limit)
	}
	return rewritten
}
