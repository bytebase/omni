package analysis

import (
	"fmt"
	"slices"

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
	return e.accessTablesOfNode(stmts[0].AST)
}

// accessTablesOfNode is extractAccessTables for a statement already parsed.
func (e *querySpanExtractor) accessTablesOfNode(node ast.Node) ([]ColumnResource, error) {
	extractor := &accessTableExtractor{
		defaultDatabase: e.defaultDatabase,
		searchPath:      e.searchPath,
		getDatabase:     e.catalog.GetDatabase,
	}

	extractor.walk(node, nil, nil)

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
	// targets holds the target relations of the DML statements walked.
	targets map[*ast.RangeVar]bool
}

// walk finds every table reference (RangeVar) under n. ctes holds the names
// of the common table expressions in scope: an unqualified reference to one
// reads the CTE, not a relation, as PostgreSQL's parser resolves it. skip is
// the WITH clause of n, which walkWith has already visited.
func (a *accessTableExtractor) walk(n ast.Node, ctes []string, skip ast.Node) {
	if w := withClauseOf(n); w != nil && skip == nil {
		a.walkWith(n, w, ctes)
		return
	}
	ast.Inspect(n, func(c ast.Node) bool {
		if a.err != nil || (skip != nil && c == skip) {
			return false
		}
		if c != n {
			if w := withClauseOf(c); w != nil {
				a.walkWith(c, w, ctes)
				return false
			}
		}
		if rv := dmlTarget(c); rv != nil {
			// A statement's target is a relation whatever CTE shares its
			// name.
			// pg: src/backend/parser/parse_clause.c — setTargetTable
			if a.targets == nil {
				a.targets = map[*ast.RangeVar]bool{}
			}
			a.targets[rv] = true
		}
		if rv, ok := c.(*ast.RangeVar); ok && (a.targets[rv] || !(rv.Schemaname == "" && rv.Catalogname == "" && slices.Contains(ctes, rv.Relname))) {
			a.processRangeVar(rv)
		}
		return true
	})
}

// walkWith walks a statement with a WITH clause. A CTE's own query sees the
// CTEs before it, or, in WITH RECURSIVE, every CTE of the clause; the rest of
// the statement sees them all.
func (a *accessTableExtractor) walkWith(n ast.Node, w *ast.WithClause, outer []string) {
	var names []string
	var ctes []*ast.CommonTableExpr
	if w.Ctes != nil {
		for _, item := range w.Ctes.Items {
			if cte, ok := item.(*ast.CommonTableExpr); ok {
				ctes = append(ctes, cte)
				names = append(names, cte.Ctename)
			}
		}
	}
	for i, cte := range ctes {
		visible := names[:i]
		if w.Recursive {
			visible = names
		}
		a.walk(cte, append(slices.Clone(outer), visible...), nil)
	}
	a.walk(n, append(slices.Clone(outer), names...), w)
}

// dmlTarget returns the target relation of an INSERT, UPDATE, DELETE or
// MERGE, or nil.
func dmlTarget(n ast.Node) *ast.RangeVar {
	switch n := n.(type) {
	case *ast.InsertStmt:
		return n.Relation
	case *ast.UpdateStmt:
		return n.Relation
	case *ast.DeleteStmt:
		return n.Relation
	case *ast.MergeStmt:
		return n.Relation
	}
	return nil
}

// withClauseOf returns the WITH clause of a statement that can have one.
func withClauseOf(n ast.Node) *ast.WithClause {
	switch n := n.(type) {
	case *ast.SelectStmt:
		return n.WithClause
	case *ast.InsertStmt:
		return n.WithClause
	case *ast.UpdateStmt:
		return n.WithClause
	case *ast.DeleteStmt:
		return n.WithClause
	case *ast.MergeStmt:
		return n.WithClause
	}
	return nil
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
