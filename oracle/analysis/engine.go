package analysis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

// Engine is Oracle's implementation of the omni/analysis capability
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

func (engine) Split(script string) []analysis.Statement {
	return SplitStatements(script)
}

// Parse parses each non-empty segment and returns one statement per node it
// holds, each carrying its segment; a segment holding no statement is dropped.
// A syntax error is positioned in the script.
func (e engine) Parse(script string) ([]analysis.ParsedStatement, error) {
	var parsed []analysis.ParsedStatement
	for _, segment := range e.Split(script) {
		if segment.Empty {
			continue
		}
		list, err := Parse(segment.Text)
		if err != nil {
			return nil, syntaxError(segment.Text, analysis.PositionAt(script, segment.Start), err)
		}
		for _, item := range list.Items {
			raw, ok := item.(*ast.RawStmt)
			if !ok {
				continue
			}
			parsed = append(parsed, analysis.ParsedStatement{Statement: segment, Node: raw.Stmt})
		}
	}
	return parsed, nil
}

// syntaxError positions a parse error of text, which starts at base in the
// script. Any other error is returned as is.
func syntaxError(text string, base analysis.Position, err error) error {
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) {
		return err
	}
	position := analysis.PositionAt(text, parseErr.Position).Shift(base)
	return &analysis.SyntaxError{
		Position:   position,
		Message:    fmt.Sprintf("Syntax error at line %d:%d: %s", position.Line, position.Column, parseErr.Message),
		RawMessage: parseErr.Message,
	}
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
	readOnly, returnsData, err := ValidateQuery(script)
	if err != nil {
		return false, false, syntaxError(script, analysis.Position{Line: 1, Column: 1}, err)
	}
	return readOnly, returnsData, nil
}

// ResultLimit falls back to ROWNUM when the FETCH NEXT rewrite fails.
func (engine) ResultLimit(statement string, limit int, engineVersion string) string {
	rewritten, err := RewriteResultLimit(statement, limit, engineVersion)
	if err != nil {
		return RowNumResultLimit(statement, limit)
	}
	return rewritten
}

var _ analysis.QuerySpanResolver = Engine

// QuerySpans resolves each statement on the request's instance, the default
// database being the schema of unqualified names. A parse error is wrapped
// without a position.
func (engine) QuerySpans(ctx context.Context, request analysis.QuerySpanRequest, statements []analysis.QuerySpanStatement) ([]*analysis.ResolvedSpan, error) {
	opts := ResolveOptions{Instance: request.Instance, DefaultDatabase: request.DefaultDatabase, Catalog: spanCatalog(ctx, request.Catalog)}
	spans := make([]*analysis.ResolvedSpan, 0, len(statements))
	for _, statement := range statements {
		span, err := ResolveQuerySpan(statement.Text, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get query span from statement: %s: %w", statement.Text, err)
		}
		spans = append(spans, span)
	}
	return spans, nil
}

// spanCatalog reads each database of an instance once per script. A nil
// lookup stays nil.
func spanCatalog(ctx context.Context, catalog analysis.Catalog) Catalog {
	var c Catalog
	if get := catalog.GetInstanceDatabase; get != nil {
		databases := map[[2]string]Database{}
		c.GetDatabase = func(instance, name string) (Database, error) {
			key := [2]string{instance, name}
			if database, ok := databases[key]; ok {
				return database, nil
			}
			database, err := get(ctx, instance, name)
			if err != nil {
				return nil, err
			}
			databases[key] = catalogDatabase(database)
			return databases[key], nil
		}
	}
	if getLinked := catalog.GetLinkedDatabase; getLinked != nil {
		c.GetLinkedDatabase = func(instance, link, schema string) (string, string, Database, error) {
			database, err := getLinked(ctx, instance, link, schema)
			if err != nil || database == nil {
				return "", "", nil, err
			}
			return database.Instance, database.Name, catalogDatabase(database), nil
		}
	}
	return c
}

// catalogDatabase is a synced database's Database, through the caller's name
// matching, or nil for a database the caller has not synced.
func catalogDatabase(database *analysis.Database) Database {
	if database == nil {
		return nil
	}
	return lookupDatabase{database.Lookup(true).Schema("")}
}

// lookupDatabase is a database's unnamed schema, where Oracle metadata keeps
// its objects.
type lookupDatabase struct{ schema *analysis.SchemaLookup }

func (d lookupDatabase) HasSchema() bool { return d.schema != nil }

func (d lookupDatabase) GetTable(name string) *metadata.TableMetadata { return d.schema.Table(name) }

func (d lookupDatabase) GetView(name string) *metadata.ViewMetadata { return d.schema.View(name) }

func (d lookupDatabase) GetMaterializedView(name string) *metadata.MaterializedViewMetadata {
	return d.schema.MaterializedView(name)
}

func (d lookupDatabase) GetExternalTable(name string) *metadata.ExternalTableMetadata {
	return d.schema.ExternalTable(name)
}

var _ analysis.ChangeExtractor = Engine

// ExtractChanges reports what the statements write, the default database
// being the schema of unqualified names, as Oracle has it. Without synced
// metadata a dropped index changes nothing; a foreign node is skipped.
func (engine) ExtractChanges(request analysis.ChangeRequest, statements []analysis.ParsedStatement) (*analysis.Changes, error) {
	var indexTable IndexTableFunc
	if request.Database != nil {
		lookup := request.Database.Lookup(true)
		indexTable = func(schema, index string) (string, bool) { return lookup.Schema(schema).IndexTable(index) }
	}
	changes := &analysis.Changes{}
	for _, statement := range statements {
		node, ok := statement.Node.(ast.Node)
		if !ok {
			continue
		}
		appendChanges(changes, ExtractChangedResources(node, request.DefaultDatabase, statement.Text, indexTable))
	}
	return changes, nil
}

func appendChanges(changes *analysis.Changes, statement StatementChanges) {
	changes.Resources = append(changes.Resources, statement.Resources...)
	for _, text := range statement.DMLStatements {
		changes.DMLStatements = append(changes.DMLStatements, analysis.DMLStatement{Text: text})
	}
	changes.DMLCount += statement.DMLCount
	changes.InsertCount += statement.InsertCount
}

var (
	_ analysis.BackupGenerator  = Engine
	_ analysis.RestoreGenerator = Engine
)

// Backups backs up the script's UPDATE and DELETE statements, whose
// unqualified tables are in the source database, into the backup database, a
// schema. A *Version request version caps the backup table names. A backup
// starts at its first character's 0-based column and ends at its last's, or
// one column later when that character is multi-byte.
func (engine) Backups(_ context.Context, request analysis.BackupRequest, script string) ([]analysis.Backup, error) {
	version, _ := request.Version.(*Version)
	backups, err := TransformDMLToSelect(script, BackupOptions{
		SourceDatabase: request.SourceDatabase,
		TargetDatabase: request.BackupDatabase,
		TablePrefix:    request.TablePrefix,
		Version:        version,
	})
	if err != nil {
		return nil, err
	}
	positions := analysis.NewPositions(script)
	result := make([]analysis.Backup, 0, len(backups))
	for _, backup := range backups {
		start := positions.At(backup.Start)
		start.Column--
		// End is the last character's last byte, which puts a multi-byte
		// character before it.
		endOffset := backup.End
		for endOffset < len(script) && !utf8.RuneStart(script[endOffset]) {
			endOffset++
		}
		end := positions.At(endOffset)
		end.Column--
		result = append(result, analysis.Backup{
			Statement:    backup.Statement,
			SourceSchema: backup.SourceSchema,
			SourceTable:  backup.SourceTableName,
			BackupTable:  backup.TargetTableName,
			Start:        start,
			End:          end,
		})
	}
	return result, nil
}

// RestoreSQL restores the rows the statements in the backup's range that
// change the source table changed. A statement starting at or before the
// range's start begins it, and one starting at or after its end ends it.
func (e engine) RestoreSQL(ctx context.Context, request analysis.RestoreRequest) (string, error) {
	statements := e.Split(request.Script)
	positions := analysis.NewPositions(request.Script)
	starts := make([]analysis.Position, 0, len(statements))
	for _, statement := range statements {
		starts = append(starts, positions.At(statement.Start))
	}
	first, last := 0, len(statements)-1
	for i, start := range starts {
		if !before(request.Start, start) {
			first = i
		}
	}
	for i := len(starts) - 1; i >= 0; i-- {
		if !before(starts[i], request.End) {
			last = i
		}
	}
	var sqls []string
	for i := first; i <= last; i++ {
		matched, err := ReferencesTable(statements[i].Text, request.SourceDatabase, request.SourceTable)
		if err != nil {
			return "", fmt.Errorf("failed to extract single SQL: failed to prepare transformation: %v", err)
		}
		if matched {
			sqls = append(sqls, NormalizeRestoreSQL(statements[i].Text))
		}
	}
	originalSQL := strings.Join(sqls, "\n")
	node, err := FindFirstDML(originalSQL)
	if err != nil {
		return "", err
	}
	if node == nil {
		return "", errors.New("no DML statement found in extracted SQL")
	}
	table, err := analysis.LookUpRestoreTable(ctx, request, "", true, nil)
	if err != nil {
		return "", err
	}
	return RestoreSQL(originalSQL, node, RestoreOptions{
		SourceDatabase: request.SourceDatabase,
		SourceTable:    request.SourceTable,
		BackupDatabase: request.BackupDatabase,
		BackupTable:    request.BackupTable,
		Table:          table.Table,
		PrimaryKey:     table.PrimaryKey,
	})
}

// before reports whether a is before b, comparing lines, then columns.
func before(a, b analysis.Position) bool {
	return a.Line < b.Line || a.Line == b.Line && a.Column < b.Column
}
