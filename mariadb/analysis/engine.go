package analysis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/metadata"
)

// Engine is MariaDB's implementation of the omni/analysis capability interfaces.
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
		end := withSemicolon(script, s.ByteEnd)
		statements = append(statements, analysis.Statement{
			Text:  script[s.ByteStart:end],
			Empty: s.Empty(),
			Start: s.ByteStart,
			End:   end,
		})
	}
	return statements
}

// withSemicolon extends a split segment's end over the ";" that follows it.
func withSemicolon(script string, end int) int {
	if end < len(script) && script[end] == ';' {
		end++
	}
	return end
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
	catalog := unnamedSchemaCatalog(ctx, request.Catalog)
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

// unnamedSchemaCatalog reads the request's catalog, each database through its
// unnamed schema.
func unnamedSchemaCatalog(ctx context.Context, catalog analysis.Catalog) Catalog {
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

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, resolving unqualified
// names in the default database. A statement without a node writes nothing.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	var changes analysis.Changes
	for _, statement := range statements {
		node, ok := statement.Node.(ast.Node)
		if !ok && statement.Node != nil {
			return nil, errors.New("expected AST for MariaDB")
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

var (
	_ analysis.BackupGenerator  = Engine
	_ analysis.RestoreGenerator = Engine
)

// Backups backs up the script's UPDATE and DELETE statements, an unqualified
// table being in the source database. A backup starts where its first
// statement's segment starts and ends where its last one's ends, as Split
// places them.
func (engine) Backups(ctx context.Context, request analysis.BackupRequest, script string) ([]analysis.Backup, error) {
	catalog := unnamedSchemaCatalog(ctx, request.Catalog)
	backups, err := TransformDMLToSelect(script, BackupOptions{
		DefaultDatabase:      request.SourceDatabase,
		BackupDatabase:       request.BackupDatabase,
		TablePrefix:          request.TablePrefix,
		CaseSensitive:        request.CaseSensitive,
		Catalog:              catalog,
		ReferencingDatabases: request.ReferencingDatabases,
		Columns:              columnClassifier(catalog, request.CaseSensitive),
	})
	if err != nil || backups == nil {
		return nil, err
	}
	positions := analysis.NewPositions(script)
	result := make([]analysis.Backup, 0, len(backups))
	for _, backup := range backups {
		result = append(result, analysis.Backup{
			Statement:   backup.Statement,
			SourceTable: backup.SourceTable,
			BackupTable: backup.BackupTable,
			Start:       positions.At(backup.ByteStart),
			End:         positions.At(withSemicolon(script, backup.ByteEnd)),
		})
	}
	return result, nil
}

// columnClassifier is ClassifyColumns over catalog, naming the database the
// catalog does not list.
func columnClassifier(catalog Catalog, caseSensitive bool) func(database, table string) (TableColumns, error) {
	return func(database, table string) (TableColumns, error) {
		columns, err := ClassifyColumns(catalog, caseSensitive, database, table)
		var notFound *NotFoundError
		if errors.As(err, &notFound) {
			return TableColumns{}, fmt.Errorf("failed to get database metadata for Database %q", database)
		}
		return columns, err
	}
}

// RestoreSQL restores the rows the statements in the backup's range that name
// the source table changed.
func (engine) RestoreSQL(ctx context.Context, request analysis.RestoreRequest) (string, error) {
	var originalSQL strings.Builder
	for _, statement := range statementsInRange(request.Script, request.Start, request.End) {
		matched, err := ReferencesTable(statement.Text, request.SourceDatabase, request.SourceTable, request.CaseSensitive)
		if err != nil {
			return "", fmt.Errorf("failed to extract single SQL: failed to parse sql: %v", err)
		}
		if matched {
			originalSQL.WriteString(statement.Text)
		}
	}
	nodes, err := DMLs(originalSQL.String())
	if err != nil {
		return "", fmt.Errorf("failed to parse statement: %w", err)
	}
	if len(nodes) == 0 {
		return "", errors.New("no DML statement found in extracted SQL")
	}
	catalog := unnamedSchemaCatalog(ctx, request.Catalog)
	columns, err := columnClassifier(catalog, request.CaseSensitive)(request.SourceDatabase, request.SourceTable)
	if err != nil {
		return "", fmt.Errorf("failed to classify columns for %s.%s: %w", request.SourceDatabase, request.SourceTable, err)
	}
	return RestoreSQL(nodes, RestoreOptions{
		Statement:      originalSQL.String(),
		Database:       request.SourceDatabase,
		Table:          request.SourceTable,
		BackupDatabase: request.BackupDatabase,
		BackupTable:    request.BackupTable,
		Columns:        columns.Normal,
		Catalog:        catalog,
		CaseSensitive:  request.CaseSensitive,
	})
}

// statementsInRange returns Split's statements from the last one starting at
// or before start through the first one starting at or after end.
func statementsInRange(script string, start, end analysis.Position) []analysis.Statement {
	statements := Engine.Split(script)
	positions := analysis.NewPositions(script)
	first, last, ended := 0, len(statements)-1, false
	for i, statement := range statements {
		at := positions.At(statement.Start)
		if !before(start, at) {
			first = i
		}
		if !ended && !before(at, end) {
			last, ended = i, true
		}
	}
	if first > last {
		return nil
	}
	return statements[first : last+1]
}

func before(a, b analysis.Position) bool {
	return a.Line < b.Line || a.Line == b.Line && a.Column < b.Column
}
