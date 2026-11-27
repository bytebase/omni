package analysis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/bigquery/diagnostics"
	"github.com/bytebase/omni/bigquery/parser"
	"github.com/bytebase/omni/metadata"
)

// Engine is BigQuery's implementation of the omni/analysis capability
// interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Parser         = Engine
	_ analysis.Diagnoser      = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.ResultLimiter  = Engine
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

// Parse returns one statement per node of each non-empty segment. It never
// fails: a segment the parser rejects is one statement with a nil Node.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		file, err := parser.Parse(segment.Text)
		if err != nil || file == nil || len(file.Stmts) == 0 {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment})
			continue
		}
		for _, node := range file.Stmts {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: node})
		}
	}
	return parsed, nil
}

// Diagnose reports every syntax error except the parser's "statement parsing is
// not yet supported" stubs, which mark valid statements it does not model.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	positions := analysis.NewPositions(script)
	result := []analysis.Diagnostic{}
	for _, d := range diagnostics.Analyze(script) {
		if strings.HasSuffix(d.Message, "statement parsing is not yet supported") {
			continue
		}
		result = append(result, analysis.Diagnostic{
			Start:    positions.At(d.Range.Start.Offset),
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

// ResultLimit wraps the statement in a CTE, adding no LIMIT for a limit of zero
// or less. The caller has already confirmed the statement is a query.
func (engine) ResultLimit(statement string, limit int, _ string) string {
	limitPart := ""
	if limit > 0 {
		limitPart = fmt.Sprintf(" LIMIT %d", limit)
	}
	return fmt.Sprintf("WITH result AS (%s) SELECT * FROM result%s;", analysis.TrimStatement(statement), limitPart)
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves each statement with its written dataset as the
// database, the default one filling in, and a SET statement as a select. The
// default schema and IgnoreCaseSensitive are not read, and predicate columns
// are not reported. A syntax error is the parser's own.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	load := request.Catalog.MetadataLoader()
	options := ResolveOptions{
		DefaultDatabase:      request.DefaultDatabase,
		SetStatementIsSelect: true,
		GetDatabaseMetadata: func(database string) (*metadata.DatabaseSchemaMetadata, error) {
			if load == nil {
				return nil, &ResourceNotFoundError{Database: &database}
			}
			return load(ctx, database)
		},
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, options)
		if err != nil {
			return nil, err
		}
		resolved := &analysis.ResolvedSpan{Type: span.Type, SourceColumns: span.SourceColumns, Results: span.Results}
		errors.As(span.NotFoundError, &resolved.NotFoundError)
		spans = append(spans, resolved)
	}
	return spans, nil
}
