package analysis

import (
	"cmp"
	"context"
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// Engine is TiDB's implementation of the omni/analysis capability interfaces.
var Engine engine

type engine struct{}

var (
	_ analysis.Splitter        = Engine
	_ analysis.StatementRanger = Engine
	_ analysis.Diagnoser       = Engine
	_ analysis.QueryValidator  = Engine
	_ analysis.Explainer       = Engine
	_ analysis.ResultLimiter   = Engine
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

// syntaxError converts a parse error in text, which starts on base's line in
// the script, into a *analysis.SyntaxError. Any other error is returned as is.
func syntaxError(err error, text string, base analysis.Position) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	return newSyntaxError(text, base, parseErr.Position, parseErr.Message, parseErr.RelatedText)
}

// newSyntaxError positions an error at offset in text on the script's lines;
// its column is the text's own.
func newSyntaxError(text string, base analysis.Position, offset int, message, relatedText string) *analysis.SyntaxError {
	position := analysis.PositionAt(text, offset)
	position.Line += base.Line - 1
	formatted := fmt.Sprintf("Syntax error at line %d:%d: %s", position.Line, position.Column, message)
	if relatedText != "" {
		formatted += "\nrelated text: " + relatedText
	}
	return &analysis.SyntaxError{Position: position, Message: formatted, RawMessage: message}
}

func (engine) StatementRanges(script string) []analysis.LSPRange {
	return analysis.LSPRanges(script, StatementRanges(script))
}

// Diagnose reports the script's first syntax error, positioned in the text of
// its statement through the ";" that follows. An error that is not a syntax
// error is marked at the script's start.
func (engine) Diagnose(script string) []analysis.Diagnostic {
	diagnostics := []analysis.Diagnostic{}
	diagnoseErr, err := Diagnose(script)
	if err != nil {
		return append(diagnostics, analysis.Diagnostic{
			Start:    analysis.Position{Line: 1, Column: 1},
			Severity: analysis.SeverityError,
			Message:  err.Error(),
		})
	}
	if diagnoseErr == nil {
		return diagnostics
	}
	end := diagnoseErr.ByteEnd
	if end < len(script) && script[end] == ';' {
		end++
	}
	syntaxErr := newSyntaxError(script[diagnoseErr.ByteStart:end], analysis.PositionAt(script, diagnoseErr.ByteStart),
		diagnoseErr.Offset, diagnoseErr.Message, diagnoseErr.RelatedText)
	return append(diagnostics, analysis.Diagnostic{
		Start:    syntaxErr.Position,
		Severity: analysis.SeverityError,
		Message:  cmp.Or(syntaxErr.RawMessage, syntaxErr.Message),
	})
}

// ValidateQuery parses each non-empty segment. A syntax error is positioned on
// the script's lines; its column is the segment's own.
func (e engine) ValidateQuery(script string) (bool, bool, error) {
	var nodes []ast.Node
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		list, err := parser.Parse(segment.Text)
		if err != nil {
			return false, false, syntaxError(err, segment.Text, analysis.PositionAt(script, segment.Start))
		}
		nodes = append(nodes, list.Items...)
	}
	readOnly, returnsData := ValidateQuery(nodes)
	return readOnly, returnsData, nil
}

// Explain writes EXPLAIN FORMAT=tidb_json for the JSON format and a plain
// EXPLAIN for every other one.
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

// QuerySpans resolves each statement in the default database, matching names
// as TiDB does whatever IgnoreCaseSensitive asks. A parse error is the
// parser's own, unpositioned in the script.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	var catalog Catalog
	if get := request.Catalog.GetDatabase; get != nil {
		catalog.GetDatabase = func(name string) (Database, error) {
			database, err := get(ctx, name)
			if err != nil || database == nil {
				return nil, err
			}
			return lookupDatabase{database}, nil
		}
	}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, ResolveOptions{DefaultDatabase: request.DefaultDatabase, Catalog: catalog})
		if err != nil {
			return nil, fmt.Errorf("failed to get query span from statement: %s: %w", statement.Text, err)
		}
		spans = append(spans, span)
	}
	return spans, nil
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

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, resolving unqualified
// names in the default database. A statement without a node writes nothing.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	var changes analysis.Changes
	for _, statement := range statements {
		node, ok := statement.Node.(ast.Node)
		if !ok && statement.Node != nil {
			return nil, errors.New("expected TiDB AST")
		}
		c := ExtractChangedResources(node, request.DefaultDatabase, statement.Text)
		changes.Resources = append(changes.Resources, c.Resources...)
		for _, text := range c.DMLStatements {
			changes.DMLStatements = append(changes.DMLStatements, analysis.DMLStatement{Text: text})
		}
		changes.DMLCount += c.DMLCount
		changes.InsertCount += c.InsertCount
	}
	return &changes, nil
}
