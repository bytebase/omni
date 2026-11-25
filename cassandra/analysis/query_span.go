// Package analysis provides CQL query analysis for Cassandra: the query span
// Bytebase consumes for field-level lineage (masking) and SQL-editor access
// tracking, plus the statement classification it derives from the same walk.
//
// The walk is parser-driven over the cassandra/ast tree, and its names and
// lineage are held to Cassandra's result metadata by testdata/span.tsv. Catalog
// metadata, from the caller's GetDatabaseMetadata, serves two shapes. A `SELECT
// *` expands against the keyspace's synced relations, and a star over one it
// does not hold stays a SelectAsterisk marker with no columns — Bytebase's
// masker then fails closed on that result rather than guessing a projection. A
// metadata lookup for a star that FAILS is returned as an error: a star resolved
// from nothing would read as a statement that touches no column, which the
// masking and access checks cannot tell from a statement that genuinely reads
// nothing. And a materialized view's columns also read the base-table columns
// the view copies, so a masking policy on the base table applies to the view.
package analysis

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cassandra"
	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/metadata"
)

// QueryType classifies a CQL statement. The values mirror Bytebase's
// base.QueryType so the adapter maps them 1:1.
type QueryType = analysis.QueryType

const (
	// QueryTypeUnknown is a statement node this package does not classify.
	QueryTypeUnknown = analysis.QueryTypeUnknown
	// Select is a read-only SELECT.
	Select = analysis.Select
	// Explain is reserved: CQL has no EXPLAIN statement.
	Explain = analysis.Explain
	// SelectInfoSchema is reserved: CQL has no information schema in the
	// statement grammar.
	SelectInfoSchema = analysis.SelectInfoSchema
	// DDL changes schema.
	DDL = analysis.DDL
	// DML changes table data (INSERT, UPDATE, DELETE, BATCH, TRUNCATE).
	DML = analysis.DML
)

// ColumnResource identifies a keyspace/table/column resource.
type ColumnResource = analysis.ColumnResource

// ColumnSet is a set of column resources.
type ColumnSet = analysis.ColumnSet

// ResolvedQuerySpanResult describes one output column of a resolved span.
type ResolvedQuerySpanResult = analysis.QuerySpanResult

// ResolvedQuerySpan is a query span resolved against the caller's metadata.
// SourceColumns holds table-level resources (Column empty) for every table the
// statement touches; PredicateColumns holds the columns its WHERE clause reads.
type ResolvedQuerySpan = analysis.ResolvedSpan

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
	if stmt.JSON {
		results = jsonResult(results)
	}
	e.querySpan.Results = results

	e.extractWhereColumns(stmt.Where, keyspace, table)
	return nil
}

// jsonResult is the one column SELECT JSON returns, named [json] as Cassandra
// names it: a JSON object of every selected column, so it reads what they
// read. A star nothing could expand stays a SelectAsterisk marker, which the
// masking consumer rejects.
func jsonResult(results []ResolvedQuerySpanResult) []ResolvedQuerySpanResult {
	out := ResolvedQuerySpanResult{Name: "[json]", SourceColumns: ColumnSet{}}
	for _, r := range results {
		out.SelectAsterisk = out.SelectAsterisk || r.SelectAsterisk
		for c := range r.SourceColumns {
			out.SourceColumns[c] = true
		}
	}
	return []ResolvedQuerySpanResult{out}
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
		return identName(qn.Parts[0]), identName(qn.Parts[1])
	case 1:
		return e.opts.DefaultKeyspace, identName(qn.Parts[0])
	default:
		return e.opts.DefaultKeyspace, ""
	}
}

// identName is the name an identifier stands for: CQL folds an unquoted name
// to lower case, and a quoted one keeps its case.
func identName(id *ast.Identifier) string {
	if id == nil {
		return ""
	}
	if id.Quoted {
		return id.Name
	}
	return strings.ToLower(id.Name)
}

func (e *querySpanExtractor) selectElements(elems []*ast.SelectElement, keyspace, table string) ([]ResolvedQuerySpanResult, error) {
	if len(elems) == 0 {
		return nil, nil
	}
	for _, elem := range elems {
		if _, isStar := elem.Expr.(*ast.StarExpr); isStar {
			return e.expandSelectAsterisk(keyspace, table)
		}
	}

	rel := e.relation(keyspace, table)
	var results []ResolvedQuerySpanResult
	for _, elem := range elems {
		name := selectorName(elem.Expr)
		if elem.Alias != nil {
			name = identName(elem.Alias)
		}
		sc := ColumnSet{}
		for _, column := range collectColumnNames(elem.Expr) {
			rel.addSources(sc, keyspace, table, column)
		}
		_, plain := elem.Expr.(*ast.Identifier)
		results = append(results, ResolvedQuerySpanResult{
			Name:          name,
			SourceColumns: sc,
			IsPlainField:  plain,
		})
	}
	return results, nil
}

// selectorName is the name Cassandra gives an output column without an alias
// (Selectable's toString in its result metadata): the column's own name, the
// selector's text with function and type names in lower case, and a native
// function named in the system keyspace. COUNT(*) and COUNT(1) are count.
func selectorName(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return identName(e)
	case *ast.FunctionCall:
		name := strings.ToLower(e.Name.Name)
		if name == "count" && (e.Star || len(e.Args) == 1 && isLiteral(e.Args[0], "1")) {
			return "count"
		}
		args := make([]string, len(e.Args))
		for i, arg := range e.Args {
			args[i] = selectorName(arg)
		}
		switch {
		case e.Star:
			args = []string{"*"}
		case name == "writetime" || name == "maxwritetime" || name == "ttl":
		case !strings.Contains(name, "."):
			name = "system." + name
		}
		return name + "(" + strings.Join(args, ", ") + ")"
	case *ast.CastExpr:
		return "cast(" + selectorName(e.Expr) + " as " + typeName(e.Type) + ")"
	case *ast.TypeHintExpr:
		return "(" + typeName(e.Type) + ")" + selectorName(e.Expr)
	case *ast.BinaryExpr:
		return selectorName(e.Left) + " " + e.Op + " " + selectorName(e.Right)
	case *ast.UnaryExpr:
		return e.Op + selectorName(e.Operand)
	case *ast.DotAccess:
		return selectorName(e.Object) + "." + identName(e.Field)
	case *ast.IndexAccess:
		return selectorName(e.Collection) + "[" + selectorName(e.Index) + "]"
	case *ast.SliceAccess:
		var from, to string
		if e.From != nil {
			from = selectorName(e.From)
		}
		if e.To != nil {
			to = selectorName(e.To)
		}
		return selectorName(e.Collection) + "[" + from + ".." + to + "]"
	case *ast.StringLit:
		return "'" + strings.ReplaceAll(e.Val, "'", "''") + "'"
	case *ast.IntegerLit:
		return e.Val
	case *ast.FloatLit:
		return e.Val
	case *ast.BoolLit:
		if e.Val {
			return "true"
		}
		return "false"
	case *ast.NullLit:
		return "null"
	case *ast.UUIDLit:
		return e.Val
	case *ast.HexLit:
		return e.Val
	case *ast.DurationLit:
		return e.Val
	case *ast.ListLit:
		return "[" + selectorNames(e.Elements) + "]"
	case *ast.VectorLit:
		return "[" + selectorNames(e.Elements) + "]"
	case *ast.SetLit:
		return "{" + selectorNames(e.Elements) + "}"
	case *ast.TupleLit:
		return "(" + selectorNames(e.Elements) + ")"
	case *ast.MapLit:
		entries := make([]string, len(e.Keys))
		for i := range e.Keys {
			entries[i] = selectorName(e.Keys[i]) + ": " + selectorName(e.Values[i])
		}
		return "{" + strings.Join(entries, ", ") + "}"
	case *ast.BindMarker:
		if e.Name != "" {
			return ":" + e.Name
		}
		return "?"
	default:
		return ""
	}
}

func selectorNames(exprs []ast.ExprNode) string {
	names := make([]string, len(exprs))
	for i, expr := range exprs {
		names[i] = selectorName(expr)
	}
	return strings.Join(names, ", ")
}

func isLiteral(expr ast.ExprNode, val string) bool {
	lit, ok := expr.(*ast.IntegerLit)
	return ok && lit.Val == val
}

// typeName renders a CQL type in lower case, as Cassandra names it in a
// selector's name.
func typeName(t *ast.DataType) string {
	if t == nil {
		return ""
	}
	name := strings.ToLower(t.Name.Name)
	if t.Keyspace != nil {
		name = identName(t.Keyspace) + "." + name
	}
	var params []string
	for _, p := range t.TypeParams {
		params = append(params, typeName(p))
	}
	if t.Dimension != nil {
		params = append(params, t.Dimension.Val)
	}
	if len(params) > 0 {
		name += "<" + strings.Join(params, ", ") + ">"
	}
	return name
}

// expandSelectAsterisk expands `SELECT *` against the keyspace's synced
// relations. Without metadata — or when no schema of the keyspace holds the
// relation — it answers one SelectAsterisk result with no columns, which the
// masking consumer rejects rather than masks by guess. A table and a view are
// expanded the same way, from the relation's own ordered column list; a
// materialized view's columns also read the base-table columns it copies.
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
	rel, found := findRelation(dbMetadata, table)
	if !found {
		return unexpanded, nil
	}
	e.addBaseTables(rel, keyspace)
	var results []ResolvedQuerySpanResult
	for _, col := range rel.columns {
		sc := ColumnSet{}
		rel.addSources(sc, keyspace, table, col.GetName())
		results = append(results, ResolvedQuerySpanResult{
			Name:          col.GetName(),
			SourceColumns: sc,
			IsPlainField:  true,
		})
	}
	return results, nil
}

// relation looks up the relation a selection other than a star reads, for
// what a materialized view copies from its base table. Its columns are not
// needed, so a keyspace the caller holds no metadata for — or cannot read —
// resolves each column to the relation it names, as it does without metadata.
func (e *querySpanExtractor) relation(keyspace, table string) relationInfo {
	if e.opts.GetDatabaseMetadata == nil || table == "" {
		return relationInfo{}
	}
	dbMetadata, err := e.opts.GetDatabaseMetadata(keyspace)
	if err != nil {
		return relationInfo{}
	}
	rel, _ := findRelation(dbMetadata, table)
	e.addBaseTables(rel, keyspace)
	return rel
}

// addBaseTables adds the base tables a materialized view copies to the tables
// the statement reads.
func (e *querySpanExtractor) addBaseTables(rel relationInfo, keyspace string) {
	for _, dep := range rel.dependencies {
		e.querySpan.SourceColumns[ColumnResource{Database: keyspace, Table: dep.GetTable()}] = true
	}
}

// relationInfo is a synced relation: its ordered columns and, for a
// materialized view, the base-table columns it copies.
type relationInfo struct {
	columns      []*metadata.ColumnMetadata
	dependencies []*metadata.DependencyColumn
}

// addSources adds column of the relation, and the base-table column a
// materialized view copies into it, to sc.
func (r relationInfo) addSources(sc ColumnSet, keyspace, table, column string) {
	sc[ColumnResource{Database: keyspace, Table: table, Column: column}] = true
	for _, dep := range r.dependencies {
		if dep.GetColumn() == column {
			sc[ColumnResource{Database: keyspace, Table: dep.GetTable(), Column: column}] = true
		}
	}
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

// collectColumnNames returns the columns an expression reads: a field or an
// element of a column reads the column.
func collectColumnNames(expr ast.ExprNode) []string {
	switch e := expr.(type) {
	case *ast.Identifier:
		return []string{identName(e)}
	case *ast.DotAccess:
		return collectColumnNames(e.Object)
	case *ast.IndexAccess:
		return append(collectColumnNames(e.Collection), collectColumnNames(e.Index)...)
	case *ast.SliceAccess:
		names := collectColumnNames(e.Collection)
		if e.From != nil {
			names = append(names, collectColumnNames(e.From)...)
		}
		if e.To != nil {
			names = append(names, collectColumnNames(e.To)...)
		}
		return names
	case *ast.FunctionCall:
		var names []string
		for _, arg := range e.Args {
			names = append(names, collectColumnNames(arg)...)
		}
		return names
	case *ast.CastExpr:
		return collectColumnNames(e.Expr)
	case *ast.TypeHintExpr:
		return collectColumnNames(e.Expr)
	case *ast.BinaryExpr:
		return append(collectColumnNames(e.Left), collectColumnNames(e.Right)...)
	case *ast.UnaryExpr:
		return collectColumnNames(e.Operand)
	default:
		return nil
	}
}

// findRelation returns the synced relation with exactly this name — a table
// first, then a materialized view (cassandra/driver's sync lists one from
// system_schema.views, its columns from system_schema.columns under the view's
// own name, and the base-table columns it copies as its dependency columns),
// then a view. CQL identifiers arrive folded by identName (quoted ones keep
// their case), so a case-folded match would resolve a differently-cased
// relation the engine itself treats as distinct.
func findRelation(dbMetadata *metadata.DatabaseSchemaMetadata, name string) (relationInfo, bool) {
	for _, schema := range dbMetadata.GetSchemas() {
		for _, table := range schema.GetTables() {
			if table.GetName() == name {
				return relationInfo{columns: table.GetColumns()}, true
			}
		}
		for _, view := range schema.GetMaterializedViews() {
			if view.GetName() == name {
				return relationInfo{columns: view.GetColumns(), dependencies: view.GetDependencyColumns()}, true
			}
		}
		for _, view := range schema.GetViews() {
			if view.GetName() == name {
				return relationInfo{columns: view.GetColumns()}, true
			}
		}
	}
	return relationInfo{}, false
}
