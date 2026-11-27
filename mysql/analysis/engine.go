package analysis

import (
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/parser"
)

// Engine is MySQL's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter       = Engine
	_ analysis.Parser         = Engine
	_ analysis.Diagnoser      = Engine
	_ analysis.QueryValidator = Engine
	_ analysis.Explainer      = Engine
	_ analysis.ResultLimiter  = Engine
)

// Split cuts the script lexically; a segment's text runs through the ";" that
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

// Parse parses each non-empty segment and returns one statement per node it
// holds, each carrying its segment; a segment holding no node is dropped. A
// syntax error is positioned in the script.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		list, err := parser.Parse(segment.Text)
		if err != nil {
			return nil, syntaxError(err, segment.Text, analysis.PositionAt(script, segment.Start))
		}
		for _, node := range list.Items {
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: node})
		}
	}
	return parsed, nil
}

// syntaxError converts a parse error in text, which starts at base in the
// script, into a *analysis.SyntaxError. Any other error is returned as is.
func syntaxError(err error, text string, base analysis.Position) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(text, parseErr.Position).Shift(base)
	message := fmt.Sprintf("Syntax error at line %d:%d: %s", position.Line, position.Column, parseErr.Message)
	if parseErr.RelatedText != "" {
		message += "\nrelated text: " + parseErr.RelatedText
	}
	return &analysis.SyntaxError{Position: position, Message: message, RawMessage: parseErr.Message}
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
	if err.Offset >= 0 {
		position = analysis.PositionAt(err.Text, err.Offset)
	}
	return append(diagnostics, analysis.Diagnostic{
		Start:    position,
		Severity: analysis.SeverityError,
		Message:  err.Message,
	})
}

func (engine) ValidateQuery(script string) (bool, bool, error) {
	list, err := parser.Parse(script)
	if err != nil {
		return false, false, syntaxError(err, script, analysis.Position{Line: 1, Column: 1})
	}
	readOnly, returnsData := ValidateQuery(list.Items)
	return readOnly, returnsData, nil
}

// Explain writes EXPLAIN FORMAT=JSON for the JSON format and a plain EXPLAIN
// for every other one.
func (engine) Explain(statement string, format analysis.ExplainFormat) (string, error) {
	explain, err := ExplainStatement(statement, format == analysis.ExplainFormatJSON)
	if err != nil {
		return "", syntaxError(err, statement, analysis.Position{Line: 1, Column: 1})
	}
	return explain, nil
}

func (engine) ResultLimit(statement string, limit int, _ string) string {
	return StatementWithResultLimit(statement, limit)
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves each statement in the default database. A syntax error
// is positioned in the script; a statement given without its start, in its
// own text.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	catalog := querySpanCatalog(ctx, request.Catalog)
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, ResolveOptions{
			DefaultDatabase:     request.DefaultDatabase,
			IgnoreCaseSensitive: request.IgnoreCaseSensitive,
			Catalog:             catalog,
		})
		if err != nil {
			return nil, syntaxError(err, statement.Text, statement.Start)
		}
		spans = append(spans, span)
	}
	return spans, nil
}

// querySpanCatalog reads the request's catalog, each database through its
// unnamed schema.
func querySpanCatalog(ctx context.Context, catalog analysis.Catalog) Catalog {
	var c Catalog
	if list := catalog.ListDatabaseNames; list != nil {
		c.ListDatabaseNames = func() ([]string, error) { return list(ctx) }
	}
	if get := catalog.GetDatabase; get != nil {
		c.GetDatabase = func(name string) (Database, error) {
			database, err := get(ctx, name)
			if err != nil || database == nil {
				return nil, err
			}
			return lookupDatabase{database}, nil
		}
	}
	return c
}

// lookupDatabase is a synced database's Database, through the caller's name
// matching.
type lookupDatabase struct{ database *analysis.Database }

func (d lookupDatabase) schema() *analysis.SchemaLookup { return d.database.Lookup(false).Schema("") }

// Name is the name the metadata holds, not the one the catalog resolved.
func (d lookupDatabase) Name() string             { return d.database.Metadata.GetName() }
func (d lookupDatabase) HasSchema() bool          { return d.schema() != nil }
func (d lookupDatabase) ListTableNames() []string { return d.schema().ListTableNames() }
func (d lookupDatabase) ListViewNames() []string  { return d.schema().ListViewNames() }

func (d lookupDatabase) GetTable(name string) *metadata.TableMetadata {
	return d.schema().Table(name)
}

func (d lookupDatabase) GetView(name string) *metadata.ViewMetadata {
	return d.schema().View(name)
}
