// Package analysis provides CQL query analysis for Cassandra: the query span
// Bytebase consumes for field-level lineage (masking) and SQL-editor access
// tracking, plus the statement classification it derives from the same walk.
//
// The walk is parser-driven over the cassandra/ast tree. A `SELECT *` is the
// only shape that needs catalog metadata: the caller's GetDatabaseMetadata
// supplies the keyspace's synced relations, and a star over one it does not hold
// stays a SelectAsterisk marker with no columns — Bytebase's masker then fails
// closed on that result rather than guessing a projection. A metadata lookup
// that FAILS is returned as an error: a star resolved from nothing would read as
// a statement that touches no column, which the masking and access checks cannot
// tell from a statement that genuinely reads nothing.
package analysis

import (
	"fmt"

	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/metadata"
)

// QueryType classifies a CQL statement. The values mirror Bytebase's
// base.QueryType so the adapter maps them 1:1.
type QueryType int

const (
	// QueryTypeUnknown is a statement node this package does not classify.
	QueryTypeUnknown QueryType = iota
	// Select is a read-only SELECT.
	Select
	// Explain is reserved: CQL has no EXPLAIN statement.
	Explain
	// SelectInfoSchema is reserved: CQL has no information schema in the
	// statement grammar.
	SelectInfoSchema
	// DDL changes schema.
	DDL
	// DML changes table data (INSERT, UPDATE, DELETE, BATCH, TRUNCATE).
	DML
)

// ColumnResource identifies a keyspace/table/column resource.
type ColumnResource struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

// ColumnSet is a set of column resources.
type ColumnSet map[ColumnResource]bool

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult struct {
	Name string
	// SourceColumns is the lineage of this output column. Empty for an
	// expression whose source this package cannot name.
	SourceColumns ColumnSet
	// IsPlainField marks a direct column passthrough.
	IsPlainField bool
	// SelectAsterisk marks a `SELECT *` that could not be expanded because the
	// caller supplied no metadata for the table (or supplied none for it).
	SelectAsterisk bool
}

// ResolvedQuerySpan is a query span resolved against the caller's metadata.
// SourceColumns holds table-level resources (Column empty) for every table the
// statement touches; PredicateColumns holds the columns its WHERE clause reads.
type ResolvedQuerySpan struct {
	Type             QueryType
	SourceColumns    ColumnSet
	PredicateColumns ColumnSet
	Results          []ResolvedQuerySpanResult
}

// ResolveOptions carries the facts the caller resolves from its metadata store.
// DefaultKeyspace is the session keyspace used for a table reference without a
// keyspace qualifier. GetDatabaseMetadata is keyed by keyspace and may be nil,
// in which case a star cannot be expanded.
type ResolveOptions struct {
	DefaultKeyspace     string
	GetDatabaseMetadata func(keyspace string) (*metadata.DatabaseSchemaMetadata, error)
}

// ResolveQuerySpan analyzes a single CQL statement and returns the tables it
// touches, the columns its WHERE clause reads, its output columns with their
// lineage, and its classification.
//
// A statement text with no statement yields an empty span; more than one is an
// error, matching the one-statement contract every engine's query span keeps.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedQuerySpan, error) {
	stmts, err := cassandra.Parse(statement)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return &ResolvedQuerySpan{
			SourceColumns: ColumnSet{},
			Results:       []ResolvedQuerySpanResult{},
		}, nil
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expecting only one statement to get query span, but got %d", len(stmts))
	}

	e := &querySpanExtractor{
		opts: opts,
		querySpan: &ResolvedQuerySpan{
			Type:             QueryTypeUnknown,
			Results:          []ResolvedQuerySpanResult{},
			SourceColumns:    ColumnSet{},
			PredicateColumns: ColumnSet{},
		},
	}
	return e.extract(stmts[0].AST)
}

type querySpanExtractor struct {
	opts      ResolveOptions
	querySpan *ResolvedQuerySpan
}

func (e *querySpanExtractor) extract(node ast.Node) (*ResolvedQuerySpan, error) {
	switch stmt := node.(type) {
	case *ast.SelectStmt:
		if err := e.extractSelect(stmt); err != nil {
			return nil, err
		}
	case *ast.InsertStmt:
		e.extractInsert(stmt)
	case *ast.UpdateStmt:
		e.extractUpdate(stmt)
	case *ast.DeleteStmt:
		e.extractDelete(stmt)
	case *ast.BatchStmt:
		e.querySpan.Type = DML
		for _, child := range stmt.Statements {
			e.extract(child)
		}
	case *ast.TruncateStmt:
		e.querySpan.Type = DML
	case *ast.CreateTableStmt, *ast.AlterTableStmt, *ast.DropTableStmt,
		*ast.CreateKeyspaceStmt, *ast.AlterKeyspaceStmt, *ast.DropKeyspaceStmt,
		*ast.CreateIndexStmt, *ast.DropIndexStmt,
		*ast.CreateMVStmt, *ast.AlterMVStmt, *ast.DropMVStmt,
		*ast.CreateTypeStmt, *ast.AlterTypeStmt, *ast.DropTypeStmt,
		*ast.CreateFunctionStmt, *ast.DropFunctionStmt,
		*ast.CreateAggregateStmt, *ast.DropAggregateStmt,
		*ast.CreateTriggerStmt, *ast.DropTriggerStmt:
		e.querySpan.Type = DDL
	default:
	}
	return e.querySpan, nil
}

func (e *querySpanExtractor) extractSelect(stmt *ast.SelectStmt) error {
	e.querySpan.Type = Select

	keyspace, table := e.qualifiedNameParts(stmt.From)

	if table != "" {
		e.querySpan.SourceColumns[ColumnResource{
			Database: keyspace,
			Table:    table,
		}] = true
	}

	results, err := e.selectElements(stmt.Elements, keyspace, table)
	if err != nil {
		return err
	}
	e.querySpan.Results = results

	e.extractWhereColumns(stmt.Where, keyspace, table)
	return nil
}

func (e *querySpanExtractor) extractInsert(stmt *ast.InsertStmt) {
	e.querySpan.Type = DML
	keyspace, table := e.qualifiedNameParts(stmt.Table)
	if table != "" {
		e.querySpan.SourceColumns[ColumnResource{
			Database: keyspace,
			Table:    table,
		}] = true
	}
}

func (e *querySpanExtractor) extractUpdate(stmt *ast.UpdateStmt) {
	e.querySpan.Type = DML
	keyspace, table := e.qualifiedNameParts(stmt.Table)
	if table != "" {
		e.querySpan.SourceColumns[ColumnResource{
			Database: keyspace,
			Table:    table,
		}] = true
	}
	e.extractWhereColumns(stmt.Where, keyspace, table)
}

func (e *querySpanExtractor) extractDelete(stmt *ast.DeleteStmt) {
	e.querySpan.Type = DML
	keyspace, table := e.qualifiedNameParts(stmt.From)
	if table != "" {
		e.querySpan.SourceColumns[ColumnResource{
			Database: keyspace,
			Table:    table,
		}] = true
	}
	e.extractWhereColumns(stmt.Where, keyspace, table)
}

func (e *querySpanExtractor) qualifiedNameParts(qn *ast.QualifiedName) (keyspace, table string) {
	if qn == nil {
		return e.opts.DefaultKeyspace, ""
	}
	switch len(qn.Parts) {
	case 2:
		return qn.Parts[0].Name, qn.Parts[1].Name
	case 1:
		return e.opts.DefaultKeyspace, qn.Parts[0].Name
	default:
		return e.opts.DefaultKeyspace, ""
	}
}

func (e *querySpanExtractor) selectElements(elems []*ast.SelectElement, keyspace, table string) ([]ResolvedQuerySpanResult, error) {
	if len(elems) == 0 {
		return nil, nil
	}

	var results []ResolvedQuerySpanResult
	for _, elem := range elems {
		if _, isStar := elem.Expr.(*ast.StarExpr); isStar {
			return e.expandSelectAsterisk(keyspace, table)
		}

		sourceName := exprColumnName(elem.Expr)
		resultName := sourceName
		if elem.Alias != nil {
			resultName = elem.Alias.Name
		}
		if resultName == "" && sourceName == "" {
			results = append(results, ResolvedQuerySpanResult{
				Name:          "",
				SourceColumns: ColumnSet{},
			})
			continue
		}

		sc := ColumnSet{}
		sc[ColumnResource{
			Database: keyspace,
			Table:    table,
			Column:   sourceName,
		}] = true

		results = append(results, ResolvedQuerySpanResult{
			Name:          resultName,
			SourceColumns: sc,
			IsPlainField:  true,
		})
	}
	return results, nil
}

func exprColumnName(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return e.Name
	case *ast.DotAccess:
		if e.Field != nil {
			return e.Field.Name
		}
	default:
	}
	return ""
}

// expandSelectAsterisk expands `SELECT *` against the keyspace's synced
// relations. Without metadata — or when no schema of the keyspace holds the
// relation — it answers one SelectAsterisk result with no columns, which the
// masking consumer rejects rather than masks by guess. A table and a view (a
// materialized view is one: Cassandra keeps its definition in
// system_schema.views) are expanded the same way, from the relation's own ordered
// column list.
func (e *querySpanExtractor) expandSelectAsterisk(keyspace, table string) ([]ResolvedQuerySpanResult, error) {
	unexpanded := []ResolvedQuerySpanResult{{
		SourceColumns:  ColumnSet{},
		SelectAsterisk: true,
	}}
	if e.opts.GetDatabaseMetadata == nil || table == "" {
		return unexpanded, nil
	}

	dbMetadata, err := e.opts.GetDatabaseMetadata(keyspace)
	if err != nil {
		return nil, fmt.Errorf("failed to get database metadata for keyspace %q: %w", keyspace, err)
	}
	if dbMetadata == nil {
		return unexpanded, nil
	}

	for _, schema := range dbMetadata.GetSchemas() {
		columns, found := relationColumns(schema, table)
		if !found {
			continue
		}
		var results []ResolvedQuerySpanResult
		for _, col := range columns {
			sc := ColumnSet{}
			sc[ColumnResource{
				Database: keyspace,
				Table:    table,
				Column:   col.GetName(),
			}] = true
			results = append(results, ResolvedQuerySpanResult{
				Name:          col.GetName(),
				SourceColumns: sc,
				IsPlainField:  true,
			})
		}
		return results, nil
	}

	return unexpanded, nil
}

func (e *querySpanExtractor) extractWhereColumns(where []ast.ExprNode, keyspace, table string) {
	for _, expr := range where {
		e.extractColumnRefsFromExpr(expr, keyspace, table)
	}
}

func (e *querySpanExtractor) extractColumnRefsFromExpr(expr ast.ExprNode, keyspace, table string) {
	switch ex := expr.(type) {
	case *ast.BinaryExpr:
		e.addPredicateColumnsFromExpr(ex.Left, keyspace, table)
	case *ast.InExpr:
		e.addPredicateColumnsFromExpr(ex.Column, keyspace, table)
	case *ast.ContainsExpr:
		e.addPredicateColumnsFromExpr(ex.Column, keyspace, table)
	case *ast.TupleInExpr:
		for _, col := range ex.Columns {
			e.addPredicateColumnsFromExpr(col, keyspace, table)
		}
	case *ast.TupleCompareExpr:
		for _, col := range ex.Columns {
			e.addPredicateColumnsFromExpr(col, keyspace, table)
		}
	default:
	}
}

func (e *querySpanExtractor) addPredicateColumnsFromExpr(expr ast.ExprNode, keyspace, table string) {
	for _, name := range collectColumnNames(expr) {
		e.querySpan.PredicateColumns[ColumnResource{
			Database: keyspace,
			Table:    table,
			Column:   name,
		}] = true
	}
}

func collectColumnNames(expr ast.ExprNode) []string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return []string{e.Name}
	case *ast.DotAccess:
		var names []string
		names = append(names, collectColumnNames(e.Object)...)
		if e.Field != nil {
			names = append(names, e.Field.Name)
		}
		return names
	case *ast.FunctionCall:
		var names []string
		for _, arg := range e.Args {
			names = append(names, collectColumnNames(arg)...)
		}
		return names
	default:
		return nil
	}
}

// relationColumns returns the ordered column list of the relation with exactly
// this name — a table first, then a view (which is how a materialized view
// reaches a schema snapshot: Cassandra lists it in system_schema.views, and its
// columns in system_schema.columns under the view's own name, never in
// system_schema.tables). CQL identifiers arrive already normalized by the parser
// (quoted ones keep their case), so an additional case-folded match would resolve
// a differently-cased relation the engine itself treats as distinct.
func relationColumns(schema *metadata.SchemaMetadata, name string) ([]*metadata.ColumnMetadata, bool) {
	for _, table := range schema.GetTables() {
		if table.GetName() == name {
			return table.GetColumns(), true
		}
	}
	for _, view := range schema.GetViews() {
		if view.GetName() == name {
			return view.GetColumns(), true
		}
	}
	return nil, false
}
