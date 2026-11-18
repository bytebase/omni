package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/parser"
)

type querySpanExtractor struct {
	defaultDatabase string

	catalog Catalog

	ignoreCaseSensitive bool

	// ctes records common table expressions visible to the current select scope.
	ctes []*pseudoTable

	// outerTableSources resolves correlated subquery column references.
	outerTableSources []tableSource

	// tableSourceFrom is the set of table sources visible from the current FROM clause.
	tableSourceFrom []tableSource

	// priorTableInFrom resolves JSON_TABLE document expressions against preceding FROM items.
	priorTableInFrom []tableSource

	viewResolutionStack map[string]bool
	// source is the original SQL text, used to slice expression names via Loc.
	source string
}

func newQuerySpanExtractor(defaultDatabase string, catalog Catalog, ignoreCaseSensitive bool) *querySpanExtractor {
	return &querySpanExtractor{
		defaultDatabase:     defaultDatabase,
		catalog:             catalog,
		ignoreCaseSensitive: ignoreCaseSensitive,
	}
}

type joinType int

const (
	Join joinType = iota
	InnerJoin
	CrossJoin
	StraightJoin
	LeftOuterJoin
	RightOuterJoin
	NaturalInnerJoin
	NaturalLeftOuterJoin
	NaturalRightOuterJoin
	FullOuterJoin
	NaturalFullOuterJoin
)

// joinTableSources returns the columns a join exposes, in the order MySQL
// returns them. A column merged by USING or NATURAL is COALESCE(left, right),
// and in matched rows both sides hold its value, so it carries both sides'
// sources. Merged columns come first, in the order of the join's first table
// (the right one for RIGHT JOIN), then that table's other columns, then the
// other table's.
func joinTableSources(l, r tableSource, tp joinType, using []string) tableSource {
	switch tp {
	case Join, InnerJoin, CrossJoin, StraightJoin, LeftOuterJoin, RightOuterJoin,
		NaturalInnerJoin, NaturalLeftOuterJoin, NaturalRightOuterJoin,
		FullOuterJoin, NaturalFullOuterJoin:
	default:
		return nil
	}
	first, second := l.GetQuerySpanResult(), r.GetQuerySpanResult()
	names := joinedColumnNames(l, r, tp, using)
	if len(names) == 0 {
		return &pseudoTable{Columns: append(append([]QuerySpanResult{}, first...), second...)}
	}
	if tp == RightOuterJoin || tp == NaturalRightOuterJoin {
		first, second = second, first
	}
	joined := make(map[string]bool)
	for _, name := range names {
		joined[strings.ToLower(name)] = true
	}
	var columns []QuerySpanResult
	for _, field := range first {
		if !joined[strings.ToLower(field.Name)] {
			continue
		}
		if other, ok := querySpanResultNamed(second, field.Name); ok {
			field.SourceColumns, _ = mergeColumnSet(field.SourceColumns, other.SourceColumns)
		}
		columns = append(columns, field)
	}
	for _, side := range [][]QuerySpanResult{first, second} {
		for _, field := range side {
			if !joined[strings.ToLower(field.Name)] {
				columns = append(columns, field)
			}
		}
	}
	return &pseudoTable{
		Columns: columns,
	}
}

// joinedColumnNames returns the columns a join merges: the USING list, or for
// a NATURAL join every column name both sides share.
func joinedColumnNames(l, r tableSource, tp joinType, using []string) []string {
	switch tp {
	case NaturalInnerJoin, NaturalLeftOuterJoin, NaturalRightOuterJoin, NaturalFullOuterJoin:
		var names []string
		rightColumns := r.GetQuerySpanResult()
		for _, field := range l.GetQuerySpanResult() {
			if _, ok := querySpanResultNamed(rightColumns, field.Name); ok {
				names = append(names, field.Name)
			}
		}
		return names
	default:
		return using
	}
}

func querySpanResultNamed(results []QuerySpanResult, name string) (QuerySpanResult, bool) {
	for _, result := range results {
		if strings.EqualFold(result.Name, name) {
			return result, true
		}
	}
	return QuerySpanResult{}, false
}

func (q *querySpanExtractor) getAllTableColumnSources(databaseName, tableName string) ([]QuerySpanResult, bool) {
	findInTableSource := func(tableSource tableSource) ([]QuerySpanResult, bool) {
		if q.ignoreCaseSensitive {
			if databaseName != "" && !strings.EqualFold(databaseName, tableSource.GetDatabaseName()) {
				return nil, false
			}
			if tableName != "" && !strings.EqualFold(tableName, tableSource.GetTableName()) {
				return nil, false
			}
		} else {
			if databaseName != "" && databaseName != tableSource.GetDatabaseName() {
				return nil, false
			}
			if tableName != "" && tableName != tableSource.GetTableName() {
				return nil, false
			}
		}
		return tableSource.GetQuerySpanResult(), true
	}

	// The subquery's own FROM first, then enclosing queries, nearest first.
	for i := len(q.tableSourceFrom) - 1; i >= 0; i-- {
		if querySpanResult, ok := findInTableSource(q.tableSourceFrom[i]); ok {
			return querySpanResult, true
		}
	}
	for i := len(q.outerTableSources) - 1; i >= 0; i-- {
		if querySpanResult, ok := findInTableSource(q.outerTableSources[i]); ok {
			return querySpanResult, true
		}
	}
	return nil, false
}

func (q *querySpanExtractor) getFieldColumnSource(databaseName, tableName, fieldName string) (QuerySpanResult, error) {
	findInTableSource := func(tableSource tableSource) (QuerySpanResult, bool) {
		if q.ignoreCaseSensitive {
			if databaseName != "" && !strings.EqualFold(databaseName, tableSource.GetDatabaseName()) {
				return QuerySpanResult{}, false
			}
			if tableName != "" && !strings.EqualFold(tableName, tableSource.GetTableName()) {
				return QuerySpanResult{}, false
			}
		} else {
			if databaseName != "" && databaseName != tableSource.GetDatabaseName() {
				return QuerySpanResult{}, false
			}
			if tableName != "" && tableName != tableSource.GetTableName() {
				return QuerySpanResult{}, false
			}
		}
		if _, ok := tableSource.(*unknownColumnsTable); ok {
			return QuerySpanResult{Name: fieldName}, true
		}
		for _, column := range tableSource.GetQuerySpanResult() {
			if strings.EqualFold(column.Name, fieldName) {
				return column, true
			}
		}
		return QuerySpanResult{}, false
	}

	// A column resolves in the subquery's own FROM first, then in enclosing
	// queries, nearest first, as a correlated reference.
	for i := len(q.tableSourceFrom) - 1; i >= 0; i-- {
		if source, ok := findInTableSource(q.tableSourceFrom[i]); ok {
			return source, nil
		}
	}
	for i := len(q.priorTableInFrom) - 1; i >= 0; i-- {
		if source, ok := findInTableSource(q.priorTableInFrom[i]); ok {
			return source, nil
		}
	}
	for i := len(q.outerTableSources) - 1; i >= 0; i-- {
		if source, ok := findInTableSource(q.outerTableSources[i]); ok {
			return source, nil
		}
	}

	return QuerySpanResult{}, &NotFoundError{
		Database: &databaseName,
		Table:    &tableName,
		Column:   &fieldName,
	}
}

func (q *querySpanExtractor) findTableSchema(databaseName, tableName string) (tableSource, error) {
	if databaseName == "" {
		for i := len(q.ctes) - 1; i >= 0; i-- {
			table := q.ctes[i]
			if table.Name == tableName {
				return table, nil
			}
		}
	}

	if databaseName == "" {
		databaseName = q.defaultDatabase
	}

	var dbMetadata Database
	allDatabaseNames, err := q.catalog.ListDatabaseNames()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}
	if q.ignoreCaseSensitive {
		for _, db := range allDatabaseNames {
			if strings.EqualFold(db, databaseName) {
				dbMetadata, err = q.catalog.GetDatabase(db)
				if err != nil {
					return nil, fmt.Errorf("failed to get database metadata for database %q: %w", db, err)
				}
				break
			}
		}
	} else {
		for _, db := range allDatabaseNames {
			if db == databaseName {
				dbMetadata, err = q.catalog.GetDatabase(db)
				if err != nil {
					return nil, fmt.Errorf("failed to get database metadata for database %q: %w", db, err)
				}
				break
			}
		}
	}
	if dbMetadata == nil {
		return nil, &NotFoundError{
			Database: &databaseName,
		}
	}

	emptySchema := ""
	if !dbMetadata.HasSchema() {
		return nil, &NotFoundError{
			Database: &databaseName,
			Schema:   &emptySchema,
		}
	}

	var tableSchema *metadata.TableMetadata
	if q.ignoreCaseSensitive {
		for _, table := range dbMetadata.ListTableNames() {
			if strings.EqualFold(table, tableName) {
				tableSchema = dbMetadata.GetTable(table)
				break
			}
		}
	} else {
		tableSchema = dbMetadata.GetTable(tableName)
	}
	if tableSchema != nil {
		columnNames := make([]string, 0, len(tableSchema.GetColumns()))
		for _, column := range tableSchema.GetColumns() {
			columnNames = append(columnNames, column.Name)
		}
		return &physicalTable{
			Name:     tableSchema.Name,
			Schema:   emptySchema,
			Database: dbMetadata.Name(),
			Columns:  columnNames,
		}, nil
	}

	var viewSchema *metadata.ViewMetadata
	if q.ignoreCaseSensitive {
		for _, view := range dbMetadata.ListViewNames() {
			if strings.EqualFold(view, tableName) {
				viewSchema = dbMetadata.GetView(view)
				break
			}
		}
	} else {
		viewSchema = dbMetadata.GetView(tableName)
	}
	if viewSchema != nil {
		columns, err := q.getColumnsForView(dbMetadata.Name(), viewSchema.Name, viewSchema.Definition)
		if err != nil {
			return nil, fmt.Errorf("failed to get columns for view %q: %w", tableName, err)
		}
		return &physicalView{
			Name:     viewSchema.Name,
			Schema:   emptySchema,
			Database: dbMetadata.Name(),
			Columns:  columns,
		}, nil
	}

	return nil, &NotFoundError{
		Database: &databaseName,
		Schema:   &emptySchema,
		Table:    &tableName,
	}
}

func (q *querySpanExtractor) getColumnsForView(databaseName, viewName, definition string) ([]QuerySpanResult, error) {
	key := mysqlViewResolutionKey(databaseName, viewName)
	if q.viewResolutionStack[key] {
		return nil, fmt.Errorf("cyclic view reference detected while resolving %q", viewName)
	}

	newQ := newQuerySpanExtractor(q.defaultDatabase, q.catalog, q.ignoreCaseSensitive)
	newQ.viewResolutionStack = cloneViewResolutionStack(q.viewResolutionStack)
	newQ.viewResolutionStack[key] = true
	span, err := newQ.getQuerySpan(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to get query span for view: %w", err)
	}
	if span.NotFoundError != nil {
		return nil, span.NotFoundError
	}
	return span.Results, nil
}

func mysqlViewResolutionKey(databaseName, viewName string) string {
	return databaseName + "\x00" + viewName
}

func cloneViewResolutionStack(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func isMixedQuery(m ColumnSet, ignoreCaseSensitive bool) (bool, bool) {
	hasSystem, hasUser := false, false
	for table := range m {
		if isSystemResource(table, ignoreCaseSensitive) {
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

var reservedSystemDatabases = map[string]bool{
	"information_schema": true,
	"performance_schema": true,
}

var onDiskSystemDatabases = map[string]bool{
	"mysql": true,
}

func isSystemResource(resource ColumnResource, ignoreCaseSensitive bool) bool {
	if reservedSystemDatabases[strings.ToLower(resource.Database)] {
		return true
	}
	database := resource.Database
	if ignoreCaseSensitive {
		database = strings.ToLower(database)
	}
	return onDiskSystemDatabases[database]
}

// getQuerySpan resolves one statement; view definitions recurse through it.
func (q *querySpanExtractor) getQuerySpan(statement string) (*ResolvedSpan, error) {
	q.source = statement

	stmts, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	if stmts == nil || len(stmts.Items) == 0 {
		return &ResolvedSpan{
			Type:             Select,
			SourceColumns:    ColumnSet{},
			Results:          []QuerySpanResult{},
			PredicateColumns: ColumnSet{},
		}, nil
	}
	if len(stmts.Items) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts.Items))
	}

	root := stmts.Items[0]
	accessTables := collectAccessTables(root, q.defaultDatabase)
	allSystems, mixed := isMixedQuery(accessTables, q.ignoreCaseSensitive)
	if mixed {
		return nil, ErrMixUserSystemTables
	}

	queryType := classifyQueryType(root, allSystems, q.source)
	if allSystems {
		accessTables = ColumnSet{}
	}
	explain, isExplain := root.(*ast.ExplainStmt)
	explainAnalyze := isExplain && explain.Analyze

	var results []QuerySpanResult
	if queryType == Select {
		tableSource, err := q.extractSelectRoot(root)
		if err != nil {
			var resourceNotFound *NotFoundError
			if errors.As(err, &resourceNotFound) {
				return &ResolvedSpan{
					Type:             Select,
					SourceColumns:    accessTables,
					Results:          []QuerySpanResult{},
					PredicateColumns: ColumnSet{},
					NotFoundError:    resourceNotFound,
				}, nil
			}
			return nil, err
		}
		results = tableSource.GetQuerySpanResult()
	}
	return &ResolvedSpan{
		Type:             queryType,
		ExplainAnalyze:   explainAnalyze,
		SourceColumns:    accessTables,
		Results:          results,
		PredicateColumns: ColumnSet{},
	}, nil
}

func (q *querySpanExtractor) extractSelectRoot(root ast.Node) (*pseudoTable, error) {
	switch n := root.(type) {
	case *ast.SelectStmt:
		return q.extractSelectStmt(n)
	case *ast.ExplainStmt:
		return &pseudoTable{}, nil
	case *ast.TableStmt:
		return q.extractTableStmt(n)
	case *ast.ValuesStmt:
		return q.extractValuesStmt(n)
	default:
		return &pseudoTable{}, nil
	}
}

func (q *querySpanExtractor) extractSelectStmt(stmt *ast.SelectStmt) (*pseudoTable, error) {
	if stmt == nil {
		return &pseudoTable{}, nil
	}
	if stmt.Into != nil {
		return nil, errors.New("meet unsupported select statement with into")
	}

	originalTableSourceFrom := q.tableSourceFrom
	originalPriorLength := len(q.priorTableInFrom)
	originalCTELength := len(q.ctes)
	defer func() {
		q.tableSourceFrom = originalTableSourceFrom
		q.priorTableInFrom = q.priorTableInFrom[:originalPriorLength]
		q.ctes = q.ctes[:originalCTELength]
	}()

	if err := q.processCTEs(stmt.CTEs); err != nil {
		return nil, err
	}
	if inner := stmt.ParenSource; inner != nil {
		return q.extractSelectStmt(inner)
	}
	if tableStmt := stmt.TableSource; tableStmt != nil {
		return q.extractTableStmt(tableStmt)
	}
	if valuesStmt := stmt.ValuesSource; valuesStmt != nil {
		return q.extractValuesStmt(valuesStmt)
	}
	if stmt.SetOp != ast.SetOpNone {
		return q.extractSetOp(stmt)
	}

	fromSources, err := q.extractTableSources(stmt.From)
	if err != nil {
		return nil, err
	}
	q.tableSourceFrom = append(q.tableSourceFrom, fromSources...)

	results, err := q.extractTargetList(stmt.TargetList, fromSources)
	if err != nil {
		return nil, err
	}
	return &pseudoTable{
		Columns: results,
	}, nil
}

func (q *querySpanExtractor) extractTableStmt(stmt *ast.TableStmt) (*pseudoTable, error) {
	if stmt == nil || stmt.Table == nil {
		return &pseudoTable{}, nil
	}
	if stmt.Into != nil {
		return nil, errors.New("meet unsupported table statement with into")
	}
	tableSource, err := q.findTableSchema(stmt.Table.Schema, stmt.Table.Name)
	if err != nil {
		return nil, err
	}
	return &pseudoTable{Columns: tableSource.GetQuerySpanResult()}, nil
}

func (q *querySpanExtractor) extractValuesStmt(stmt *ast.ValuesStmt) (*pseudoTable, error) {
	if stmt == nil || len(stmt.Rows) == 0 {
		return &pseudoTable{}, nil
	}
	var columns []QuerySpanResult
	for _, expr := range stmt.Rows[0] {
		field, err := q.extractExpr(expr)
		if err != nil {
			return nil, err
		}
		columns = append(columns, field)
	}
	return &pseudoTable{Columns: columns}, nil
}

func (q *querySpanExtractor) extractSetOp(stmt *ast.SelectStmt) (*pseudoTable, error) {
	if stmt == nil || stmt.SetOp == ast.SetOpNone {
		return q.extractSelectStmt(stmt)
	}
	left, err := q.extractSelectStmt(stmt.Left)
	if err != nil {
		return nil, err
	}
	right, err := q.extractSelectStmt(stmt.Right)
	if err != nil {
		return nil, err
	}
	leftResults := left.GetQuerySpanResult()
	rightResults := right.GetQuerySpanResult()
	if len(leftResults) != len(rightResults) {
		return nil, fmt.Errorf("MySQL UNION operator left has %d fields, right has %d fields", len(leftResults), len(rightResults))
	}
	results := make([]QuerySpanResult, 0, len(leftResults))
	for i := range leftResults {
		sourceColumns, _ := mergeColumnSet(leftResults[i].SourceColumns, rightResults[i].SourceColumns)
		results = append(results, QuerySpanResult{
			Name:          leftResults[i].Name,
			SourceColumns: sourceColumns,
			IsPlainField:  false,
		})
	}
	return &pseudoTable{Columns: results}, nil
}

func (q *querySpanExtractor) processCTEs(ctes []*ast.CommonTableExpr) error {
	for _, cte := range ctes {
		tableSource, err := q.extractCTE(cte)
		if err != nil {
			return err
		}
		if tableSource != nil {
			q.ctes = append(q.ctes, tableSource)
		}
	}
	return nil
}

func (q *querySpanExtractor) extractCTE(cte *ast.CommonTableExpr) (*pseudoTable, error) {
	if cte == nil {
		return nil, nil
	}
	if cte.Recursive {
		return q.extractRecursiveCTE(cte)
	}
	return q.extractNonRecursiveCTE(cte)
}

func (q *querySpanExtractor) extractNonRecursiveCTE(cte *ast.CommonTableExpr) (*pseudoTable, error) {
	tableSource, err := q.extractSelectStmt(cte.Select)
	if err != nil {
		return nil, err
	}
	results := cloneQuerySpanResults(tableSource.GetQuerySpanResult())
	if len(cte.Columns) > 0 {
		if len(cte.Columns) != len(results) {
			return nil, fmt.Errorf("MySQL CTE column list should have the same length, but got %d and %d", len(cte.Columns), len(results))
		}
		for i := range cte.Columns {
			results[i].Name = cte.Columns[i]
		}
	}
	return &pseudoTable{Name: cte.Name, Columns: results}, nil
}

func (q *querySpanExtractor) extractRecursiveCTE(cte *ast.CommonTableExpr) (*pseudoTable, error) {
	selectStmt := ast.UnwrapParenSource(cte.Select)
	if selectStmt == nil || selectStmt.SetOp == ast.SetOpNone {
		return q.extractNonRecursiveCTE(cte)
	}
	initialTable, err := q.extractRecursiveCTEAnchor(selectStmt)
	if err != nil {
		return nil, err
	}
	results := cloneQuerySpanResults(initialTable.GetQuerySpanResult())
	if len(cte.Columns) > 0 {
		if len(cte.Columns) != len(results) {
			return nil, fmt.Errorf("The common table expression and column names list have different column counts")
		}
		for i := range cte.Columns {
			results[i].Name = cte.Columns[i]
		}
	}
	cteInfo := &pseudoTable{Name: cte.Name, Columns: results}
	q.ctes = append(q.ctes, cteInfo)
	defer func() {
		q.ctes = q.ctes[:len(q.ctes)-1]
	}()

	for {
		recursiveTable, err := q.extractSelectStmt(cte.Select)
		if err != nil {
			return nil, err
		}
		recursiveResults := recursiveTable.GetQuerySpanResult()
		if len(recursiveResults) != len(cteInfo.Columns) {
			return nil, fmt.Errorf("The common table expression and column names list have different column counts")
		}
		changed := false
		for i := range recursiveResults {
			var ok bool
			cteInfo.Columns[i].SourceColumns, ok = mergeColumnSet(cteInfo.Columns[i].SourceColumns, recursiveResults[i].SourceColumns)
			cteInfo.Columns[i].IsPlainField = false
			changed = changed || ok
		}
		if !changed {
			break
		}
	}
	return cteInfo, nil
}

func (q *querySpanExtractor) extractRecursiveCTEAnchor(stmt *ast.SelectStmt) (*pseudoTable, error) {
	stmt = ast.UnwrapParenSource(stmt)
	if stmt == nil || stmt.SetOp == ast.SetOpNone {
		return q.extractSelectStmt(stmt)
	}
	if stmt.Left != nil {
		return q.extractRecursiveCTEAnchor(stmt.Left)
	}
	return q.extractRecursiveCTEAnchor(stmt.Right)
}

func cloneQuerySpanResults(results []QuerySpanResult) []QuerySpanResult {
	clone := make([]QuerySpanResult, len(results))
	copy(clone, results)
	return clone
}

func (q *querySpanExtractor) extractTableSources(tableExprs []ast.TableExpr) ([]tableSource, error) {
	var result []tableSource
	for _, tableExpr := range tableExprs {
		tableSource, err := q.extractTableSource(tableExpr)
		if err != nil {
			return nil, err
		}
		if tableSource == nil {
			continue
		}
		q.priorTableInFrom = append(q.priorTableInFrom, tableSource)
		result = append(result, tableSource)
	}
	return result, nil
}

func (q *querySpanExtractor) extractTableSource(tableExpr ast.TableExpr) (tableSource, error) {
	switch t := tableExpr.(type) {
	case nil:
		return nil, nil
	case *ast.TableRef:
		if isDualTable(t) {
			return &pseudoTable{Name: t.Name}, nil
		}
		// A table read through a database link or named in another catalog is
		// remote: its columns are not the local catalog's, so nothing of a
		// same-named local table's is.
		if t.Dblink != "" || t.Catalog != "" {
			name := t.Name
			if t.Alias != "" {
				name = t.Alias
			}
			return &pseudoTable{Name: name}, nil
		}
		tableSource, err := q.findTableSchema(t.Schema, t.Name)
		if err != nil {
			return nil, err
		}
		if t.Alias == "" {
			return tableSource, nil
		}
		return &pseudoTable{
			Name:    t.Alias,
			Columns: tableSource.GetQuerySpanResult(),
		}, nil
	case *ast.SubqueryExpr:
		return q.extractSubquery(t)
	case *ast.JoinClause:
		return q.extractJoin(t)
	case *ast.JsonTableExpr:
		return q.extractJSONTable(t)
	case *ast.UnnestExpr:
		return q.extractTableFunction(t.Alias, t.Columns, t.Args)
	case *ast.RbIterateExpr:
		var columns []string
		if t.Column != "" {
			columns = []string{t.Column}
		}
		return q.extractTableFunction(t.Alias, columns, []ast.ExprNode{t.Arg})
	case *ast.ExternalTableRef:
		// The file's columns are known only when the query runs, so the source
		// has a name and no columns; any column reference to it resolves.
		name := t.Alias
		if name == "" {
			name = t.Path
		}
		return &unknownColumnsTable{pseudoTable{Name: name}}, nil
	default:
		return nil, fmt.Errorf("unsupported omni MySQL table source %T", tableExpr)
	}
}

// extractTableFunction reads a table function in FROM — UNNEST or RB_ITERATE —
// as a pseudo table named by its alias (or its argument), whose columns are the
// names it declares and whose lineage is the arguments' columns.
func (q *querySpanExtractor) extractTableFunction(alias string, columns []string, args []ast.ExprNode) (tableSource, error) {
	name := alias
	if name == "" && len(args) > 0 {
		name = q.exprName(args[0])
	}
	fields := make([]QuerySpanResult, len(args))
	all := ColumnSet{}
	for i, arg := range args {
		field, err := q.extractExpr(arg)
		if err != nil {
			return nil, err
		}
		fields[i] = field
		all, _ = mergeColumnSet(all, field.SourceColumns)
	}
	var result []QuerySpanResult
	// Without a column alias list the function's outputs are its arguments.
	if len(columns) == 0 {
		for i, field := range fields {
			name := field.Name
			if name == "" {
				name = q.exprName(args[i])
			}
			result = append(result, QuerySpanResult{Name: name, SourceColumns: field.SourceColumns, IsPlainField: field.IsPlainField})
		}
		return &pseudoTable{Name: name, Columns: result}, nil
	}
	for i, column := range columns {
		// UNNEST's outputs follow its arguments: the i-th output is the i-th
		// argument's elements. With no such pairing, each reads them all.
		sources := all
		if len(columns) == len(fields) {
			sources = fields[i].SourceColumns
		}
		result = append(result, QuerySpanResult{Name: column, SourceColumns: sources})
	}
	return &pseudoTable{Name: name, Columns: result}, nil
}

func (q *querySpanExtractor) extractJSONTable(expr *ast.JsonTableExpr) (tableSource, error) {
	if expr == nil {
		return nil, nil
	}
	name := expr.Alias
	if name == "" {
		name = q.exprName(expr.Expr)
	}
	field, err := q.extractExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	var columns []QuerySpanResult
	for _, column := range flattenJSONTableColumns(expr.Columns) {
		columns = append(columns, QuerySpanResult{
			Name:          column.Name,
			SourceColumns: field.SourceColumns,
			IsPlainField:  field.IsPlainField,
		})
	}
	return &pseudoTable{Name: name, Columns: columns}, nil
}

func flattenJSONTableColumns(columns []*ast.JsonTableColumn) []*ast.JsonTableColumn {
	var result []*ast.JsonTableColumn
	for _, column := range columns {
		if column == nil {
			continue
		}
		if column.Nested {
			result = append(result, flattenJSONTableColumns(column.NestedCols)...)
			continue
		}
		result = append(result, column)
	}
	return result
}

func (q *querySpanExtractor) extractJoin(join *ast.JoinClause) (tableSource, error) {
	if join == nil {
		return nil, nil
	}
	left, err := q.extractTableSource(join.Left)
	if err != nil {
		return nil, err
	}
	if left != nil {
		q.tableSourceFrom = append(q.tableSourceFrom, left)
	}
	right, err := q.extractTableSource(join.Right)
	if err != nil {
		return nil, err
	}
	if right != nil {
		q.tableSourceFrom = append(q.tableSourceFrom, right)
	}

	var using []string
	if condition, ok := join.Condition.(*ast.UsingCondition); ok && condition != nil {
		using = condition.Columns
	}
	return joinTableSources(left, right, toJoinType(join.Type), using), nil
}

func toJoinType(joinType ast.JoinType) joinType {
	switch joinType {
	case ast.JoinInner:
		return InnerJoin
	case ast.JoinLeft:
		return LeftOuterJoin
	case ast.JoinRight:
		return RightOuterJoin
	case ast.JoinCross:
		return CrossJoin
	case ast.JoinNatural:
		return NaturalInnerJoin
	case ast.JoinStraight:
		return StraightJoin
	case ast.JoinNaturalLeft:
		return NaturalLeftOuterJoin
	case ast.JoinNaturalRight:
		return NaturalRightOuterJoin
	case ast.JoinFull:
		return FullOuterJoin
	case ast.JoinNaturalFull:
		return NaturalFullOuterJoin
	default:
		return Join
	}
}

func (q *querySpanExtractor) extractSubquery(expr *ast.SubqueryExpr) (*pseudoTable, error) {
	if expr == nil {
		return &pseudoTable{}, nil
	}
	subqueryExtractor := q.cloneForSubquery()
	if expr.Lateral {
		subqueryExtractor.outerTableSources = append(subqueryExtractor.outerTableSources, q.priorTableInFrom...)
	}
	tableSource, err := subqueryExtractor.extractSelectStmt(expr.Select)
	if err != nil {
		return nil, err
	}
	results := cloneQuerySpanResults(tableSource.GetQuerySpanResult())
	if len(expr.Columns) > 0 {
		if len(expr.Columns) != len(results) {
			return nil, fmt.Errorf("derived table column list length %d doesn't match result column length %d", len(expr.Columns), len(results))
		}
		for i := range expr.Columns {
			results[i].Name = expr.Columns[i]
		}
	}
	if expr.Alias == "" {
		return &pseudoTable{Name: tableSource.Name, Columns: results}, nil
	}
	return &pseudoTable{
		Name:    expr.Alias,
		Columns: results,
	}, nil
}

func (q *querySpanExtractor) cloneForSubquery() *querySpanExtractor {
	return &querySpanExtractor{
		defaultDatabase:     q.defaultDatabase,
		catalog:             q.catalog,
		ctes:                q.ctes,
		outerTableSources:   append(slices.Clone(q.outerTableSources), q.tableSourceFrom...),
		tableSourceFrom:     []tableSource{},
		ignoreCaseSensitive: q.ignoreCaseSensitive,
		viewResolutionStack: cloneViewResolutionStack(q.viewResolutionStack),
		source:              q.source,
	}
}

func (q *querySpanExtractor) extractTargetList(targets []ast.ExprNode, fromSources []tableSource) ([]QuerySpanResult, error) {
	var result []QuerySpanResult
	for _, target := range targets {
		switch t := target.(type) {
		case *ast.ResTarget:
			if _, ok := t.Val.(*ast.StarExpr); ok {
				fields, err := q.extractTableWild("", "", fromSources)
				if err != nil {
					return nil, err
				}
				result = append(result, fields...)
				continue
			}
			if columnRef, ok := t.Val.(*ast.ColumnRef); ok && columnRef.Star {
				fields, err := q.extractColumnWild(columnRef, fromSources)
				if err != nil {
					return nil, err
				}
				result = append(result, fields...)
				continue
			}
		case *ast.StarExpr:
			fields, err := q.extractTableWild("", "", fromSources)
			if err != nil {
				return nil, err
			}
			result = append(result, fields...)
			continue
		case *ast.ColumnRef:
			if t.Star {
				fields, err := q.extractColumnWild(t, fromSources)
				if err != nil {
					return nil, err
				}
				result = append(result, fields...)
				continue
			}
		default:
		}

		field, err := q.extractTarget(target)
		if err != nil {
			return nil, err
		}
		result = append(result, field)
	}
	return result, nil
}

func (q *querySpanExtractor) extractTarget(target ast.ExprNode) (QuerySpanResult, error) {
	switch t := target.(type) {
	case *ast.ResTarget:
		field, err := q.extractExpr(t.Val)
		if err != nil {
			return QuerySpanResult{}, err
		}
		if t.Name != "" {
			field.Name = t.Name
		}
		return field, nil
	default:
		return q.extractExpr(target)
	}
}

func (q *querySpanExtractor) extractExpr(expr ast.ExprNode) (QuerySpanResult, error) {
	switch e := expr.(type) {
	case nil:
		return QuerySpanResult{
			SourceColumns: ColumnSet{},
			IsPlainField:  true,
		}, nil
	case *ast.ColumnRef:
		if e.Star {
			return QuerySpanResult{}, fmt.Errorf("asterisk is not a scalar expression")
		}
		// A column of a view, derived table or CTE is plain only if that
		// query's column is.
		source, err := q.getFieldColumnSource(e.Schema, e.Table, e.Column)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{
			Name:          e.Column,
			SourceColumns: source.SourceColumns,
			IsPlainField:  source.IsPlainField,
		}, nil
	case *ast.StarExpr:
		return QuerySpanResult{}, fmt.Errorf("asterisk is not a scalar expression")
	case *ast.IntLit, *ast.FloatLit, *ast.StringLit, *ast.BoolLit, *ast.NullLit, *ast.HexLit, *ast.BitLit, *ast.TemporalLit:
		return QuerySpanResult{
			Name:          q.exprName(e),
			SourceColumns: ColumnSet{},
			IsPlainField:  true,
		}, nil
	case *ast.VariableRef:
		return QuerySpanResult{
			Name:          q.exprName(e),
			SourceColumns: ColumnSet{},
			IsPlainField:  true,
		}, nil
	case *ast.ParenExpr:
		field, err := q.extractExpr(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		field.Name = q.exprName(e)
		field.IsPlainField = false
		return field, nil
	case *ast.BinaryExpr:
		sourceColumns, err := q.mergeExprSources(e.Left, e.Right)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.UnaryExpr:
		sourceColumns, err := q.mergeExprSources(e.Operand)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.FuncCallExpr:
		exprs := append([]ast.ExprNode{}, e.Args...)
		exprs = append(exprs, e.Separator)
		for _, item := range e.OrderBy {
			if item != nil {
				exprs = append(exprs, item.Expr)
			}
		}
		if e.Over != nil {
			exprs = append(exprs, e.Over.PartitionBy...)
			for _, item := range e.Over.OrderBy {
				if item != nil {
					exprs = append(exprs, item.Expr)
				}
			}
		}
		sourceColumns, err := q.mergeExprSources(exprs...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.SubqueryExpr:
		sourceColumns, isPlain, err := q.extractSelectSourceColumns(e.Select)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: isPlain}, nil
	case *ast.ExistsExpr:
		sourceColumns, _, err := q.extractSelectSourceColumns(e.Select)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.CaseExpr:
		var exprs []ast.ExprNode
		exprs = append(exprs, e.Operand)
		for _, when := range e.Whens {
			if when == nil {
				continue
			}
			exprs = append(exprs, when.Cond, when.Result)
		}
		exprs = append(exprs, e.Default)
		sourceColumns, err := q.mergeExprSources(exprs...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.BetweenExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr, e.Low, e.High)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.InExpr:
		exprs := append([]ast.ExprNode{e.Expr}, e.List...)
		sourceColumns, err := q.mergeExprSources(exprs...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		if e.Select != nil {
			selectSourceColumns, _, err := q.extractSelectSourceColumns(e.Select)
			if err != nil {
				return QuerySpanResult{}, err
			}
			sourceColumns, _ = mergeColumnSet(sourceColumns, selectSourceColumns)
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.LikeExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr, e.Pattern, e.Escape)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.IsExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.CastExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.ExtractExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.IntervalExpr:
		sourceColumns, err := q.mergeExprSources(e.Value)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.CollateExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.MatchExpr:
		exprs := make([]ast.ExprNode, 0, len(e.Columns)+1)
		for _, column := range e.Columns {
			exprs = append(exprs, column)
		}
		exprs = append(exprs, e.Against)
		sourceColumns, err := q.mergeExprSources(exprs...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.ConvertExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.DefaultExpr:
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: ColumnSet{}, IsPlainField: true}, nil
	case *ast.RowExpr:
		sourceColumns, err := q.mergeExprSources(e.Items...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.MemberOfExpr:
		sourceColumns, err := q.mergeExprSources(e.Value, e.Array)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.KeywordArg:
		// Keyword arguments (TIMESTAMPDIFF(SECOND, ...), GET_FORMAT(DATE, ...))
		// are unit selectors, not column references — no lineage.
		return QuerySpanResult{Name: e.Keyword, SourceColumns: ColumnSet{}, IsPlainField: false}, nil
	case *ast.WeightStringExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.QuantifiedExpr:
		// any_expr's `ALL (expr)` is one expression, whose columns are the
		// comparison's.
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.OdbcEscapeExpr:
		sourceColumns, err := q.mergeExprSources(e.Expr)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	case *ast.RemoteSequenceExpr:
		// A sequence read through a database link is no column of a table here.
		return QuerySpanResult{Name: q.exprName(e)}, nil
	case *ast.GroupingSet:
		sourceColumns, err := q.mergeExprSources(e.Items...)
		if err != nil {
			return QuerySpanResult{}, err
		}
		return QuerySpanResult{Name: q.exprName(e), SourceColumns: sourceColumns, IsPlainField: false}, nil
	default:
		return QuerySpanResult{}, fmt.Errorf("unsupported omni MySQL expression %T", expr)
	}
}

func (q *querySpanExtractor) extractSelectSourceColumns(selectStmt *ast.SelectStmt) (ColumnSet, bool, error) {
	subqueryExtractor := q.cloneForSubquery()
	tableSource, err := subqueryExtractor.extractSelectStmt(selectStmt)
	if err != nil {
		return nil, false, err
	}
	spanResults := tableSource.GetQuerySpanResult()
	sourceColumns := make(ColumnSet)
	for _, field := range spanResults {
		sourceColumns, _ = mergeColumnSet(sourceColumns, field.SourceColumns)
	}
	isPlain := false
	if len(spanResults) == 1 {
		isPlain = spanResults[0].IsPlainField
	}
	return sourceColumns, isPlain, nil
}

func (q *querySpanExtractor) mergeExprSources(exprs ...ast.ExprNode) (ColumnSet, error) {
	result := make(ColumnSet)
	for _, expr := range exprs {
		if expr == nil {
			continue
		}
		field, err := q.extractExpr(expr)
		if err != nil {
			return nil, err
		}
		result, _ = mergeColumnSet(result, field.SourceColumns)
	}
	return result, nil
}

func (q *querySpanExtractor) extractColumnWild(columnRef *ast.ColumnRef, fromSources []tableSource) ([]QuerySpanResult, error) {
	if columnRef == nil {
		return nil, nil
	}
	return q.extractTableWild(columnRef.Schema, columnRef.Table, fromSources)
}

func (q *querySpanExtractor) extractTableWild(databaseName, tableName string, fromSources []tableSource) ([]QuerySpanResult, error) {
	if databaseName == "" && tableName == "" {
		var result []QuerySpanResult
		for _, fromSource := range fromSources {
			result = append(result, fromSource.GetQuerySpanResult()...)
		}
		return result, nil
	}

	result, ok := q.getAllTableColumnSources(databaseName, tableName)
	if !ok {
		return nil, &NotFoundError{
			Database: &databaseName,
			Table:    &tableName,
		}
	}
	return result, nil
}

func (q *querySpanExtractor) exprName(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		var parts []string
		if e.Schema != "" {
			parts = append(parts, e.Schema)
		}
		if e.Table != "" {
			parts = append(parts, e.Table)
		}
		if e.Star {
			parts = append(parts, "*")
		} else {
			parts = append(parts, e.Column)
		}
		return strings.Join(parts, ".")
	case *ast.BinaryExpr:
		return q.exprName(e.Left) + binaryOpName(e) + q.exprName(e.Right)
	case *ast.InExpr:
		var list []string
		for _, expr := range e.List {
			list = append(list, q.exprName(expr))
		}
		operator := " in "
		if e.Not {
			operator = " not in "
		}
		return q.exprName(e.Expr) + operator + "(" + strings.Join(list, ", ") + ")"
	case *ast.IntLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return strconv.FormatInt(e.Value, 10)
	case *ast.FloatLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return e.Value
	case *ast.StringLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return e.Value
	case *ast.BoolLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return strconv.FormatBool(e.Value)
	case *ast.NullLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return "NULL"
	case *ast.HexLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return e.Value
	case *ast.BitLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return e.Value
	case *ast.TemporalLit:
		if name := q.slice(e.Loc); name != "" {
			return name
		}
		return e.Value
	default:
		if loc, ok := nodeLoc(e); ok {
			return q.slice(loc)
		}
		return ""
	}
}

func binaryOpName(expr *ast.BinaryExpr) string {
	if expr == nil {
		return ""
	}
	if expr.OriginalOp != "" {
		return expr.OriginalOp
	}
	return binaryOpNames[expr.Op]
}

var binaryOpNames = map[ast.BinaryOp]string{
	ast.BinOpAdd:         "+",
	ast.BinOpSub:         "-",
	ast.BinOpMul:         "*",
	ast.BinOpDiv:         "/",
	ast.BinOpMod:         "%",
	ast.BinOpEq:          "=",
	ast.BinOpNe:          "!=",
	ast.BinOpLt:          "<",
	ast.BinOpGt:          ">",
	ast.BinOpLe:          "<=",
	ast.BinOpGe:          ">=",
	ast.BinOpAnd:         " and ",
	ast.BinOpOr:          " or ",
	ast.BinOpXor:         " xor ",
	ast.BinOpBitAnd:      "&",
	ast.BinOpBitOr:       "|",
	ast.BinOpBitXor:      "^",
	ast.BinOpShiftLeft:   "<<",
	ast.BinOpShiftRight:  ">>",
	ast.BinOpDivInt:      " div ",
	ast.BinOpRegexp:      " regexp ",
	ast.BinOpLikeEscape:  " like ",
	ast.BinOpNullSafeEq:  "<=>",
	ast.BinOpAssign:      ":=",
	ast.BinOpJsonExtract: "->",
	ast.BinOpJsonUnquote: "->>",
	ast.BinOpSoundsLike:  " sounds like ",
}

func (q *querySpanExtractor) slice(loc ast.Loc) string {
	if loc.Start < 0 || loc.End <= loc.Start || loc.End > len(q.source) {
		return ""
	}
	return strings.TrimSpace(q.source[loc.Start:loc.End])
}

func nodeLoc(node ast.Node) (ast.Loc, bool) {
	if loc, ok := exprNodeLoc(node); ok {
		return loc, true
	}
	return literalNodeLoc(node)
}

func exprNodeLoc(node ast.Node) (ast.Loc, bool) {
	switch n := node.(type) {
	case *ast.ResTarget:
		return n.Loc, true
	case *ast.ColumnRef:
		return n.Loc, true
	case *ast.BinaryExpr:
		return n.Loc, true
	case *ast.UnaryExpr:
		return n.Loc, true
	case *ast.FuncCallExpr:
		return n.Loc, true
	case *ast.SubqueryExpr:
		return n.Loc, true
	case *ast.CaseExpr:
		return n.Loc, true
	case *ast.BetweenExpr:
		return n.Loc, true
	case *ast.InExpr:
		return n.Loc, true
	case *ast.LikeExpr:
		return n.Loc, true
	case *ast.IsExpr:
		return n.Loc, true
	case *ast.ExistsExpr:
		return n.Loc, true
	case *ast.CastExpr:
		return n.Loc, true
	case *ast.ExtractExpr:
		return n.Loc, true
	case *ast.IntervalExpr:
		return n.Loc, true
	case *ast.CollateExpr:
		return n.Loc, true
	case *ast.MatchExpr:
		return n.Loc, true
	case *ast.ConvertExpr:
		return n.Loc, true
	case *ast.DefaultExpr:
		return n.Loc, true
	case *ast.RowExpr:
		return n.Loc, true
	case *ast.MemberOfExpr:
		return n.Loc, true
	case *ast.ParenExpr:
		return n.Loc, true
	case *ast.StarExpr:
		return n.Loc, true
	case *ast.QuantifiedExpr:
		return n.Loc, true
	case *ast.OdbcEscapeExpr:
		return n.Loc, true
	case *ast.RemoteSequenceExpr:
		return n.Loc, true
	case *ast.GroupingSet:
		return n.Loc, true
	default:
		return ast.Loc{}, false
	}
}

func literalNodeLoc(node ast.Node) (ast.Loc, bool) {
	switch n := node.(type) {
	case *ast.IntLit:
		return n.Loc, true
	case *ast.FloatLit:
		return n.Loc, true
	case *ast.StringLit:
		return n.Loc, true
	case *ast.BoolLit:
		return n.Loc, true
	case *ast.NullLit:
		return n.Loc, true
	case *ast.HexLit:
		return n.Loc, true
	case *ast.BitLit:
		return n.Loc, true
	case *ast.TemporalLit:
		return n.Loc, true
	case *ast.VariableRef:
		return n.Loc, true
	default:
		return ast.Loc{}, false
	}
}

func classifyQueryType(node ast.Node, allSystems bool, source string) QueryType {
	switch n := node.(type) {
	case *ast.SelectStmt, *ast.TableStmt, *ast.ValuesStmt:
		if allSystems {
			return SelectInfoSchema
		}
		return Select
	case *ast.ExplainStmt:
		if isDescribeStatement(source, n.Loc) {
			return SelectInfoSchema
		}
		// EXPLAIN INTO stores the plan in the table it names.
		if n.Into != "" {
			return DML
		}
		if n.Analyze {
			switch n.Stmt.(type) {
			case *ast.SelectStmt, *ast.TableStmt, *ast.ValuesStmt:
				return Select
			default:
				return DML
			}
		}
		return Explain
	case *ast.ShowStmt:
		return SelectInfoSchema
	case *ast.SetStmt, *ast.SetDefaultRoleStmt, *ast.SetResourceGroupStmt, *ast.SetRoleStmt, *ast.SetTransactionStmt:
		return Select
	case *ast.SetPasswordStmt, *ast.HelpStmt:
		return QueryTypeUnknown
	case *ast.CreateDatabaseStmt, *ast.CreateTableStmt, *ast.CreateIndexStmt, *ast.CreateViewStmt,
		*ast.CreateMaterializedViewStmt, *ast.CreateMaterializedViewLogStmt,
		*ast.CreateEventStmt, *ast.CreateTriggerStmt, *ast.CreateFunctionStmt,
		*ast.CreateUserStmt, *ast.CreateRoleStmt, *ast.CreateTablespaceStmt, *ast.CreateServerStmt,
		*ast.CreateLogfileGroupStmt, *ast.CreateSpatialRefSysStmt, *ast.CreateResourceGroupStmt,
		*ast.AlterDatabaseStmt, *ast.AlterTableStmt, *ast.AlterViewStmt, *ast.AlterEventStmt,
		*ast.AlterMaterializedViewStmt, *ast.AlterMaterializedViewLogStmt,
		*ast.AlterUserStmt, *ast.AlterRoutineStmt, *ast.AlterTablespaceStmt, *ast.AlterServerStmt,
		*ast.AlterLogfileGroupStmt, *ast.AlterResourceGroupStmt, *ast.AlterInstanceStmt,
		*ast.DropDatabaseStmt, *ast.DropTableStmt, *ast.DropIndexStmt, *ast.DropViewStmt,
		*ast.DropMaterializedViewLogStmt,
		*ast.DropEventStmt, *ast.DropTriggerStmt, *ast.DropRoutineStmt,
		*ast.DropUserStmt, *ast.DropRoleStmt, *ast.DropTablespaceStmt, *ast.DropServerStmt,
		*ast.DropLogfileGroupStmt, *ast.DropSpatialRefSysStmt, *ast.DropResourceGroupStmt,
		*ast.RenameTableStmt, *ast.RenameUserStmt, *ast.TruncateStmt, *ast.ImportTableStmt,
		*ast.GrantStmt, *ast.RevokeStmt, *ast.GrantRoleStmt, *ast.RevokeRoleStmt,
		*ast.InstallPluginStmt, *ast.UninstallPluginStmt, *ast.InstallComponentStmt,
		*ast.UninstallComponentStmt, *ast.CloneStmt,
		// purge_stmt and flashback_stmt destroy and restore a table (or a
		// database, a tenant or the recycle bin) as DDL does.
		*ast.PurgeStmt, *ast.FlashbackStmt:
		return DDL
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.BeginStmt, *ast.CommitStmt,
		*ast.RollbackStmt, *ast.SavepointStmt, *ast.LockTablesStmt, *ast.UnlockTablesStmt,
		*ast.LoadDataStmt, *ast.PrepareStmt, *ast.ExecuteStmt, *ast.DeallocateStmt,
		*ast.CallStmt, *ast.DoStmt, *ast.HandlerOpenStmt, *ast.HandlerReadStmt, *ast.HandlerCloseStmt,
		*ast.AnalyzeTableStmt, *ast.OptimizeTableStmt, *ast.CheckTableStmt, *ast.RepairTableStmt,
		*ast.FlushStmt, *ast.KillStmt, *ast.ShutdownStmt, *ast.RestartStmt, *ast.XAStmt,
		// method_opt is DBMS_STATS' statistics method and runs nothing; it keeps
		// ANALYZE's class rather than staying unknown, which Bytebase denies.
		*ast.MethodOptStmt,
		*ast.SignalStmt, *ast.ResignalStmt, *ast.GetDiagnosticsStmt, *ast.DeclareVarStmt,
		*ast.DeclareConditionStmt, *ast.DeclareHandlerStmt, *ast.DeclareCursorStmt,
		*ast.IfStmt, *ast.WhileStmt, *ast.RepeatStmt, *ast.LoopStmt, *ast.LeaveStmt,
		*ast.IterateStmt, *ast.ReturnStmt, *ast.OpenCursorStmt, *ast.FetchCursorStmt,
		*ast.CloseCursorStmt, *ast.LockInstanceStmt, *ast.UnlockInstanceStmt, *ast.BinlogStmt,
		*ast.CacheIndexStmt, *ast.LoadIndexIntoCacheStmt, *ast.ResetPersistStmt,
		*ast.ChangeReplicationSourceStmt, *ast.ChangeReplicationFilterStmt,
		*ast.StartReplicaStmt, *ast.StopReplicaStmt, *ast.ResetReplicaStmt,
		*ast.PurgeBinaryLogsStmt, *ast.ResetMasterStmt,
		*ast.StartGroupReplicationStmt, *ast.StopGroupReplicationStmt:
		return DML
	default:
		return QueryTypeUnknown
	}
}

func isDescribeStatement(source string, loc ast.Loc) bool {
	statement := strings.TrimSpace(source)
	if loc.Start >= 0 && loc.End > loc.Start && loc.End <= len(source) {
		statement = strings.TrimSpace(source[loc.Start:loc.End])
	}
	upper := strings.ToUpper(statement)
	fields := strings.Fields(upper)
	return len(fields) > 0 && (fields[0] == "DESCRIBE" || fields[0] == "DESC")
}

// collectAccessTablesFromMviewRefresh visits the expressions a materialized
// view refresh clause holds. START WITH and NEXT are bit_exprs, so either may
// carry a subquery that reads a table.
func collectAccessTablesFromMviewRefresh(result ColumnSet, refresh *ast.MviewRefreshInfo, defaultDatabase string) {
	if refresh == nil {
		return
	}
	collectAccessTablesFromExpr(result, refresh.StartWith, defaultDatabase)
	collectAccessTablesFromExpr(result, refresh.Next, defaultDatabase)
}

func collectAccessTables(root ast.Node, defaultDatabase string) ColumnSet {
	result := make(ColumnSet)
	collectAccessTablesFromNode(result, root, defaultDatabase)
	return result
}

func collectAccessTablesFromNode(result ColumnSet, node ast.Node, defaultDatabase string) {
	switch n := node.(type) {
	case nil:
	case *ast.SelectStmt:
		if n == nil {
			return
		}
		for _, cte := range n.CTEs {
			if cte != nil {
				collectAccessTablesFromNode(result, cte.Select, defaultDatabase)
			}
		}
		if n.Left != nil {
			collectAccessTablesFromNode(result, n.Left, defaultDatabase)
		}
		if n.Right != nil {
			collectAccessTablesFromNode(result, n.Right, defaultDatabase)
		}
		if inner := n.ParenSource; inner != nil {
			collectAccessTablesFromNode(result, inner, defaultDatabase)
		}
		if tableStmt := n.TableSource; tableStmt != nil {
			collectAccessTablesFromNode(result, tableStmt, defaultDatabase)
		}
		if valuesStmt := n.ValuesSource; valuesStmt != nil {
			collectAccessTablesFromNode(result, valuesStmt, defaultDatabase)
		}
		for _, target := range n.TargetList {
			collectAccessTablesFromExpr(result, target, defaultDatabase)
		}
		for _, tableExpr := range n.From {
			collectAccessTablesFromTableExpr(result, tableExpr, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, n.Where, defaultDatabase)
		for _, groupBy := range n.GroupBy {
			if groupBy != nil {
				collectAccessTablesFromExpr(result, groupBy.Expr, defaultDatabase)
			}
		}
		collectAccessTablesFromExpr(result, n.Having, defaultDatabase)
		for _, orderBy := range n.OrderBy {
			if orderBy != nil {
				collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
			}
		}
		// INTO OUTFILE's partition key and its FORMAT list may hold a subquery.
		if n.Into != nil {
			collectAccessTablesFromExpr(result, n.Into.PartitionBy, defaultDatabase)
			for _, prop := range n.Into.ExternalFormat {
				if prop != nil {
					collectAccessTablesFromExpr(result, prop.Value, defaultDatabase)
				}
			}
		}
	case *ast.ExplainStmt:
		collectAccessTablesFromNode(result, n.Stmt, defaultDatabase)
	case *ast.InsertStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		for _, row := range n.Values {
			for _, expr := range row {
				collectAccessTablesFromExpr(result, expr, defaultDatabase)
			}
		}
		collectAccessTablesFromNode(result, n.Select, defaultDatabase)
		collectAccessTablesFromNode(result, n.TableSource, defaultDatabase)
		for _, assignment := range n.SetList {
			if assignment != nil {
				collectAccessTablesFromExpr(result, assignment.Value, defaultDatabase)
			}
		}
		for _, assignment := range n.OnDuplicateKey {
			if assignment != nil {
				collectAccessTablesFromExpr(result, assignment.Value, defaultDatabase)
			}
		}
	case *ast.LoadDataStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		if n.Source != nil {
			collectAccessTablesFromTableExpr(result, n.Source, defaultDatabase)
		}
		collectAccessTablesFromNode(result, n.SourceQuery, defaultDatabase)
		for _, column := range n.Columns {
			collectAccessTablesFromExpr(result, column, defaultDatabase)
		}
		for _, assignment := range n.SetList {
			if assignment != nil {
				collectAccessTablesFromExpr(result, assignment.Value, defaultDatabase)
			}
		}
	case *ast.UpdateStmt:
		for _, tableExpr := range n.Tables {
			collectAccessTablesFromTableExpr(result, tableExpr, defaultDatabase)
		}
		for _, assignment := range n.SetList {
			if assignment != nil {
				collectAccessTablesFromExpr(result, assignment.Value, defaultDatabase)
			}
		}
		collectAccessTablesFromExpr(result, n.Where, defaultDatabase)
		for _, orderBy := range n.OrderBy {
			if orderBy != nil {
				collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
			}
		}
	case *ast.DeleteStmt:
		for _, tableExpr := range n.Tables {
			collectAccessTablesFromTableExpr(result, tableExpr, defaultDatabase)
		}
		for _, tableExpr := range n.Using {
			collectAccessTablesFromTableExpr(result, tableExpr, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, n.Where, defaultDatabase)
		for _, orderBy := range n.OrderBy {
			if orderBy != nil {
				collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
			}
		}
	case *ast.CreateTableStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		collectAccessTablesFromTableRef(result, n.Like, defaultDatabase)
		collectAccessTablesFromNode(result, n.Select, defaultDatabase)
	case *ast.CreateIndexStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		for _, column := range n.Columns {
			if column != nil {
				collectAccessTablesFromExpr(result, column.Expr, defaultDatabase)
			}
		}
	case *ast.CreateViewStmt:
		collectAccessTablesFromTableRef(result, n.Name, defaultDatabase)
		collectAccessTablesFromNode(result, n.Select, defaultDatabase)
	case *ast.CreateMaterializedViewStmt:
		collectAccessTablesFromTableRef(result, n.Name, defaultDatabase)
		collectAccessTablesFromNode(result, n.Select, defaultDatabase)
		if n.MviewOptions != nil {
			collectAccessTablesFromMviewRefresh(result, n.MviewOptions.Refresh, defaultDatabase)
		}
	case *ast.CreateMaterializedViewLogStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		if n.Purge != nil {
			collectAccessTablesFromExpr(result, n.Purge.StartWith, defaultDatabase)
			collectAccessTablesFromExpr(result, n.Purge.Next, defaultDatabase)
		}
	case *ast.AlterMaterializedViewStmt:
		collectAccessTablesFromTableRef(result, n.Name, defaultDatabase)
		for _, action := range n.Actions {
			if action != nil {
				collectAccessTablesFromMviewRefresh(result, action.Refresh, defaultDatabase)
			}
		}
	case *ast.AlterMaterializedViewLogStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		for _, action := range n.Actions {
			if action != nil {
				collectAccessTablesFromExpr(result, action.StartWith, defaultDatabase)
				collectAccessTablesFromExpr(result, action.Next, defaultDatabase)
			}
		}
	case *ast.DropMaterializedViewLogStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
	case *ast.AlterTableStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		for _, command := range n.Commands {
			if command == nil {
				continue
			}
			collectAccessTablesFromExpr(result, command.DefaultExpr, defaultDatabase)
			collectAccessTablesFromTableRef(result, command.ExchangeTable, defaultDatabase)
			collectAccessTablesFromOrderBy(result, command.OrderByItems, defaultDatabase)
		}
	case *ast.DropTableStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.DropIndexStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
	case *ast.DropViewStmt:
		collectAccessTablesFromTableRefs(result, n.Views, defaultDatabase)
	case *ast.RenameTableStmt:
		for _, pair := range n.Pairs {
			if pair == nil {
				continue
			}
			collectAccessTablesFromTableRef(result, pair.Old, defaultDatabase)
			collectAccessTablesFromTableRef(result, pair.New, defaultDatabase)
		}
	case *ast.TruncateStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.AnalyzeTableStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.OptimizeTableStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.CheckTableStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.RepairTableStmt:
		collectAccessTablesFromTableRefs(result, n.Tables, defaultDatabase)
	case *ast.TableStmt:
		if n == nil {
			return
		}
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		collectAccessTablesFromOrderBy(result, n.OrderBy, defaultDatabase)
		collectAccessTablesFromLimit(result, n.Limit, defaultDatabase)
	case *ast.ValuesStmt:
		for _, row := range n.Rows {
			for _, expr := range row {
				collectAccessTablesFromExpr(result, expr, defaultDatabase)
			}
		}
		collectAccessTablesFromOrderBy(result, n.OrderBy, defaultDatabase)
		collectAccessTablesFromLimit(result, n.Limit, defaultDatabase)
	case *ast.CallStmt:
		for _, arg := range n.Args {
			collectAccessTablesFromExpr(result, arg, defaultDatabase)
		}
	case *ast.SetStmt:
		for _, assignment := range n.Assignments {
			if assignment != nil {
				collectAccessTablesFromExpr(result, assignment.Value, defaultDatabase)
			}
		}
	case *ast.DoStmt:
		for _, expr := range n.Exprs {
			collectAccessTablesFromExpr(result, expr, defaultDatabase)
		}
	case *ast.HandlerOpenStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
	case *ast.HandlerReadStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
		// The key's values are read as well: one may hold a subquery.
		for _, value := range n.Values {
			collectAccessTablesFromExpr(result, value, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, n.Where, defaultDatabase)
	case *ast.HandlerCloseStmt:
		collectAccessTablesFromTableRef(result, n.Table, defaultDatabase)
	case *ast.LockTablesStmt:
		for _, table := range n.Tables {
			if table != nil {
				collectAccessTablesFromTableRef(result, table.Table, defaultDatabase)
			}
		}
	default:
	}
}

func collectAccessTablesFromOrderBy(result ColumnSet, orderByItems []*ast.OrderByItem, defaultDatabase string) {
	for _, orderBy := range orderByItems {
		if orderBy != nil {
			collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
		}
	}
}

func collectAccessTablesFromLimit(result ColumnSet, limit *ast.Limit, defaultDatabase string) {
	if limit == nil {
		return
	}
	collectAccessTablesFromExpr(result, limit.Count, defaultDatabase)
	collectAccessTablesFromExpr(result, limit.Offset, defaultDatabase)
}

func collectAccessTablesFromTableRefs(result ColumnSet, tableRefs []*ast.TableRef, defaultDatabase string) {
	for _, tableRef := range tableRefs {
		collectAccessTablesFromTableRef(result, tableRef, defaultDatabase)
	}
}

func collectAccessTablesFromTableRef(result ColumnSet, tableRef *ast.TableRef, defaultDatabase string) {
	if tableRef == nil {
		return
	}
	if isDualTable(tableRef) {
		return
	}
	// A database link's table, and one a catalog qualifies, is not the local
	// catalog's resource.
	if tableRef.Dblink != "" || tableRef.Catalog != "" {
		return
	}
	database := tableRef.Schema
	if database == "" {
		database = defaultDatabase
	}
	result[ColumnResource{
		Database: database,
		Table:    tableRef.Name,
	}] = true
}

func isDualTable(tableRef *ast.TableRef) bool {
	return tableRef != nil && tableRef.Schema == "" && strings.EqualFold(tableRef.Name, "dual")
}

func collectAccessTablesFromTableExpr(result ColumnSet, tableExpr ast.TableExpr, defaultDatabase string) {
	switch n := tableExpr.(type) {
	case *ast.TableRef:
		collectAccessTablesFromTableRef(result, n, defaultDatabase)
		if n.Flashback != nil {
			collectAccessTablesFromExpr(result, n.Flashback.Snapshot, defaultDatabase)
		}
	case *ast.JoinClause:
		collectAccessTablesFromTableExpr(result, n.Left, defaultDatabase)
		collectAccessTablesFromTableExpr(result, n.Right, defaultDatabase)
		switch condition := n.Condition.(type) {
		case *ast.OnCondition:
			collectAccessTablesFromExpr(result, condition.Expr, defaultDatabase)
		default:
		}
	case *ast.SubqueryExpr:
		collectAccessTablesFromNode(result, n.Select, defaultDatabase)
		if n.Flashback != nil {
			collectAccessTablesFromExpr(result, n.Flashback.Snapshot, defaultDatabase)
		}
	case *ast.JsonTableExpr:
		collectAccessTablesFromExpr(result, n.Expr, defaultDatabase)
		collectAccessTablesFromExpr(result, n.Path, defaultDatabase)
	case *ast.UnnestExpr:
		for _, arg := range n.Args {
			collectAccessTablesFromExpr(result, arg, defaultDatabase)
		}
	case *ast.RbIterateExpr:
		collectAccessTablesFromExpr(result, n.Arg, defaultDatabase)
	case *ast.ExternalTableRef:
		// A format or property value may hold an expression, as FIELD_DELIMITER
		// = (SELECT a FROM t1) does.
		for _, prop := range n.Source {
			if prop != nil {
				collectAccessTablesFromExpr(result, prop.Value, defaultDatabase)
			}
		}
		for _, prop := range n.Format {
			if prop != nil {
				collectAccessTablesFromExpr(result, prop.Value, defaultDatabase)
			}
		}
	default:
	}
}

func collectAccessTablesFromExpr(result ColumnSet, expr ast.ExprNode, defaultDatabase string) {
	switch e := expr.(type) {
	case nil:
	case *ast.ResTarget:
		collectAccessTablesFromExpr(result, e.Val, defaultDatabase)
	case *ast.BinaryExpr:
		collectAccessTablesFromExpr(result, e.Left, defaultDatabase)
		collectAccessTablesFromExpr(result, e.Right, defaultDatabase)
	case *ast.UnaryExpr:
		collectAccessTablesFromExpr(result, e.Operand, defaultDatabase)
	case *ast.QuantifiedExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.OdbcEscapeExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.GroupingSet:
		for _, item := range e.Items {
			collectAccessTablesFromExpr(result, item, defaultDatabase)
		}
	case *ast.FuncCallExpr:
		for _, arg := range e.Args {
			collectAccessTablesFromExpr(result, arg, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, e.Separator, defaultDatabase)
		for _, orderBy := range e.OrderBy {
			if orderBy != nil {
				collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
			}
		}
		if e.Over != nil {
			for _, partitionBy := range e.Over.PartitionBy {
				collectAccessTablesFromExpr(result, partitionBy, defaultDatabase)
			}
			for _, orderBy := range e.Over.OrderBy {
				if orderBy != nil {
					collectAccessTablesFromExpr(result, orderBy.Expr, defaultDatabase)
				}
			}
		}
	case *ast.SubqueryExpr:
		collectAccessTablesFromNode(result, e.Select, defaultDatabase)
	case *ast.CaseExpr:
		collectAccessTablesFromExpr(result, e.Operand, defaultDatabase)
		for _, when := range e.Whens {
			if when == nil {
				continue
			}
			collectAccessTablesFromExpr(result, when.Cond, defaultDatabase)
			collectAccessTablesFromExpr(result, when.Result, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, e.Default, defaultDatabase)
	case *ast.BetweenExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
		collectAccessTablesFromExpr(result, e.Low, defaultDatabase)
		collectAccessTablesFromExpr(result, e.High, defaultDatabase)
	case *ast.InExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
		for _, item := range e.List {
			collectAccessTablesFromExpr(result, item, defaultDatabase)
		}
		collectAccessTablesFromNode(result, e.Select, defaultDatabase)
	case *ast.LikeExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
		collectAccessTablesFromExpr(result, e.Pattern, defaultDatabase)
		collectAccessTablesFromExpr(result, e.Escape, defaultDatabase)
	case *ast.IsExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.ExistsExpr:
		collectAccessTablesFromNode(result, e.Select, defaultDatabase)
	case *ast.CastExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.ExtractExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.IntervalExpr:
		collectAccessTablesFromExpr(result, e.Value, defaultDatabase)
	case *ast.CollateExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.MatchExpr:
		for _, column := range e.Columns {
			collectAccessTablesFromExpr(result, column, defaultDatabase)
		}
		collectAccessTablesFromExpr(result, e.Against, defaultDatabase)
	case *ast.ConvertExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.RowExpr:
		for _, item := range e.Items {
			collectAccessTablesFromExpr(result, item, defaultDatabase)
		}
	case *ast.MemberOfExpr:
		collectAccessTablesFromExpr(result, e.Value, defaultDatabase)
		collectAccessTablesFromExpr(result, e.Array, defaultDatabase)
	case *ast.ParenExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	case *ast.WeightStringExpr:
		collectAccessTablesFromExpr(result, e.Expr, defaultDatabase)
	default:
	}
}
