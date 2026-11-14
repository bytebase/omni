package analysis

import (
	"fmt"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
)

// extractAccessTables returns every table and view statement names, resolving
// unqualified names through the search path. A name the synced metadata does
// not have, such as a CTE or a pseudo table, is left out; system relations
// are kept as written.
func (e *querySpanExtractor) extractAccessTables(statement string) ([]ColumnResource, error) {
	stmts, err := omnipg.Parse(statement)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL statement: %w", err)
	}

	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}

	extractor := &accessTableExtractor{
		defaultDatabase: e.defaultDatabase,
		searchPath:      e.searchPath,
		getDatabase:     e.catalog.GetDatabase,
	}

	// Walk the omni AST to find all RangeVar nodes (table references)
	ast.Inspect(stmts[0].AST, func(n ast.Node) bool {
		if extractor.err != nil {
			return false
		}
		rv, ok := n.(*ast.RangeVar)
		if !ok {
			return true
		}
		extractor.processRangeVar(rv)
		return true
	})

	if extractor.err != nil {
		return nil, fmt.Errorf("failed to extract access tables: %w", extractor.err)
	}

	return extractor.accessTables, nil
}

type accessTableExtractor struct {
	err             error
	defaultDatabase string
	searchPath      []string
	accessTables    []ColumnResource
	getDatabase     func(name string) (Database, error)
}

func (a *accessTableExtractor) processRangeVar(rv *ast.RangeVar) {
	resource := ColumnResource{
		Database: a.defaultDatabase,
	}

	if rv.Catalogname != "" {
		resource.Database = rv.Catalogname
		resource.Schema = rv.Schemaname
		resource.Table = rv.Relname
	} else if rv.Schemaname != "" {
		resource.Schema = rv.Schemaname
		resource.Table = rv.Relname
	} else {
		resource.Table = rv.Relname
	}

	if !isSystemResource(resource) {
		searchPath := a.searchPath
		if resource.Schema != "" {
			searchPath = []string{resource.Schema}
		}

		databaseMetadata, err := a.getDatabase(a.defaultDatabase)
		if err != nil {
			a.err = fmt.Errorf("failed to get database metadata for database: %s: %w", a.defaultDatabase, err)
			return
		}
		// Access pseudo table or table/view we do not sync, return directly.
		if databaseMetadata == nil {
			return
		}
		schemaName, name := databaseMetadata.SearchRelation(searchPath, resource.Table)
		if schemaName == "" && name == "" {
			return
		}
		resource.Schema = schemaName
	}

	a.accessTables = append(a.accessTables, resource)
}

// isMixedQuery checks whether the query accesses the user table and system table at the same time.
func isMixedQuery(m ColumnSet) (bool, bool) {
	hasSystem, hasUser := false, false
	for table := range m {
		if isSystemResource(table) {
			hasSystem = true
		} else {
			hasUser = true
		}
	}

	if hasSystem && hasUser {
		return false, true
	}

	return !hasUser && hasSystem, false
}

func isSystemResource(resource ColumnResource) bool {
	// User can access the system table/view by name directly without database/schema name.
	// For example: `SELECT * FROM pg_database`, which will access the system table `pg_database`.
	// Additionally, user can create a table/view with the same name with system table/view and access them
	// by specify the schema name, for example:
	// `CREATE TABLE pg_database(id INT); SELECT * FROM public.pg_database;` which will access the user table `pg_database`.
	if IsSystemSchema(resource.Schema) {
		return true
	}
	if resource.Schema == "" && IsSystemView(resource.Table) {
		return true
	}
	if resource.Schema == "" && IsSystemTable(resource.Table) {
		return true
	}
	return false
}
