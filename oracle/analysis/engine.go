package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/analysis"
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
