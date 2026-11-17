package catalog

import (
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/analysis"
	"github.com/bytebase/omni/mssql/ast"
)

// getViewDependencies returns the tables a view definition reads, as
// "schema.table" identifiers. It mirrors Bytebase's schema/mssql implementation:
// parse the CREATE VIEW, resolve the select body's source columns against an
// empty catalog, and keep the table-level resources.
func getViewDependencies(viewDef, schemaName string) ([]string, error) {
	stmts, err := mssql.Parse(viewDef)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return nil, nil
	}
	var query *ast.SelectStmt
	for i := range stmts {
		createView, ok := stmts[i].AST.(*ast.CreateViewStmt)
		if !ok || createView.Query == nil {
			continue
		}
		query, _ = createView.Query.(*ast.SelectStmt)
		break
	}
	if query == nil {
		return nil, nil
	}
	start := query.Loc.Start
	end := selectStmtSpanEnd(query)
	if start < 0 || end < 0 || start >= end || end > len(viewDef) {
		return nil, nil
	}
	span, err := analysis.ResolveQuerySpan(viewDef[start:end], analysis.ResolveOptions{
		DefaultSchema:       schemaName,
		IgnoreCaseSensitive: false,
		Catalog: analysis.Catalog{
			ListDatabaseNames: func() ([]string, error) { return nil, nil },
			GetDatabase:       func(string) (analysis.Database, error) { return emptyDatabase{schemaName: schemaName}, nil },
		},
	})
	if err != nil {
		return nil, nil //nolint:nilerr // a view whose lineage is unmodeled contributes no edges
	}
	seen := make(map[string]bool)
	for resource := range span.SourceColumns {
		seen[objectID(resource.Schema, resource.Table)] = true
	}
	var dependencies []string
	for dep := range seen {
		dependencies = append(dependencies, dep)
	}
	return dependencies, nil
}

// selectStmtSpanEnd returns the largest End across a set-operation chain. A
// combined SelectStmt inherits End=-1, so the right-most leaf holds the end.
func selectStmtSpanEnd(s *ast.SelectStmt) int {
	if s == nil {
		return -1
	}
	end := s.Loc.End
	if e := selectStmtSpanEnd(s.Larg); e > end {
		end = e
	}
	if e := selectStmtSpanEnd(s.Rarg); e > end {
		end = e
	}
	return end
}

// emptyDatabase answers every lookup with an empty schema, so lineage records
// each table a view names without needing the table's columns.

type emptyDatabase struct {
	schemaName string
}

func (d emptyDatabase) ListSchemaNames() []string { return []string{d.schemaName} }
func (d emptyDatabase) GetSchema(string) analysis.Schema {
	return emptySchema{}
}

type emptySchema struct{}

func (emptySchema) ListTableNames() []string                { return nil }
func (emptySchema) GetTable(string) *metadata.TableMetadata { return nil }
func (emptySchema) ListViewNames() []string                 { return nil }
func (emptySchema) GetView(string) *metadata.ViewMetadata   { return nil }
