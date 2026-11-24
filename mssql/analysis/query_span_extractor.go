package analysis

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"unicode"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
)

type querySpanExtractor struct {
	defaultDatabase     string
	defaultSchema       string
	ignoreCaseSensitive bool

	catalog Catalog
	// tempTables is the batch's temp tables and table variables, shared with
	// the caller so later statements see what earlier ones declared.
	tempTables map[string]*TempTable
	// ctes records the common table expressions visible in the current scope.
	// Pushed when entering a WITH clause; truncated on exit.
	ctes []*pseudoTable

	// outerTableSources lets correlated subqueries resolve columns against the
	// enclosing query's FROM clause.
	outerTableSources []tableSource

	// tableSourcesFrom is this scope's FROM clause.
	tableSourcesFrom []tableSource

	predicateColumns ColumnSet

	viewResolutionStack map[string]bool

	// source is the original SQL text, used to slice expression names via Loc.
	source string
	// setOption is the statement when it is a SET of session options.
	setOption *ast.SetOptionStmt
}

func newQuerySpanExtractor(defaultDatabase string, defaultSchema string, catalog Catalog, tempTables map[string]*TempTable, ignoreCaseSensitive bool) *querySpanExtractor {
	if defaultSchema == "" {
		// Fall back to the default schema `dbo`.
		// Reference: https://learn.microsoft.com/en-us/sql/relational-databases/security/authentication-access/ownership-and-user-schema-separation#the-dbo-schema
		defaultSchema = "dbo"
	}
	return &querySpanExtractor{
		defaultDatabase:     defaultDatabase,
		defaultSchema:       defaultSchema,
		catalog:             catalog,
		tempTables:          tempTables,
		ignoreCaseSensitive: ignoreCaseSensitive,
		predicateColumns:    make(ColumnSet),
	}
}

func (q *querySpanExtractor) findTempTable(rawName string) (*TempTable, bool) {
	for name, tempTable := range q.tempTables {
		if q.isIdentifierEqual(rawName, name) {
			return tempTable, true
		}
	}
	return nil, false
}

// untracedError reports the values a statement that answers with rows reads
// and this call cannot trace: a #temp table or @table variable, whose rows the
// batch's own statements wrote, and a local variable, whose value a statement
// of the batch assigned. A write and a SELECT ... INTO answer with no rows, so
// they read nothing a caller has to withhold.
//
// The scan is syntactic, so it reports these whether or not the span resolved a
// table source: a query of only the server's metadata returns before resolution
// and still reads whatever its FROM clause names. A variable method call (the
// XML receiver of @x.nodes('/r')) names its variable in the node, not in a
// child.
func (q *querySpanExtractor) untracedError(root ast.Node) *UntracedError {
	if !returnsRows(root) {
		return nil
	}
	seenRelations := map[ColumnResource]bool{}
	seenVariables := map[string]bool{}
	var relations []ColumnResource
	var variables []string
	addRelation := func(name string) {
		// Name the relation as the batch declared it, whatever case the
		// statement wrote.
		if tempTable, ok := q.findTempTable(name); ok {
			name = tempTable.Name
		}
		relation := ColumnResource{Table: name}
		if name == "" || seenRelations[relation] {
			return
		}
		seenRelations[relation] = true
		relations = append(relations, relation)
	}
	addVariable := func(name string) {
		// A @@ name is a system function, whose value is the server's state
		// rather than a column.
		if name == "" || strings.HasPrefix(name, "@@") || seenVariables[name] {
			return
		}
		seenVariables[name] = true
		variables = append(variables, name)
	}
	ast.Inspect(root, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.VariableRef:
			addVariable(n.Name)
		case *ast.TableVarMethodCallRef:
			addVariable(n.Var)
		case *ast.TableRef:
			if strings.HasPrefix(n.Object, "#") {
				addRelation(n.Object)
			}
		case *ast.TableVarRef:
			addRelation(n.Name)
		default:
		}
		return true
	})
	if len(relations) == 0 && len(variables) == 0 {
		return nil
	}
	slices.SortFunc(relations, func(a, b ColumnResource) int {
		return cmp.Or(cmp.Compare(a.Server, b.Server), cmp.Compare(a.Database, b.Database), cmp.Compare(a.Schema, b.Schema), cmp.Compare(a.Table, b.Table))
	})
	slices.Sort(variables)
	return &UntracedError{Relations: relations, Variables: variables}
}

// tsqlFindTableSchemaByParts resolves a (linkedServer, database, schema, table)
// tuple to a TableSource. Empty strings for database/schema mean "not specified
// by the user"; defaults are applied here when looking up metadata. CTEs
// registered in q.ctes shadow physical tables when no database/schema was
// specified.
func (q *querySpanExtractor) tsqlFindTableSchemaByParts(linkedServer, rawDatabase, rawSchema, rawTable string) (tableSource, error) {
	if linkedServer != "" {
		// The server resolves the relation itself, through sys.servers. No
		// object of the connected instance describes the table, so neither its
		// columns nor a masking policy for them exist here.
		parts := make([]string, 0, 3)
		for _, part := range []string{rawDatabase, rawSchema, rawTable} {
			if part != "" {
				parts = append(parts, part)
			}
		}
		return nil, &LinkedServerError{Server: linkedServer, Table: strings.Join(parts, ".")}
	}
	if strings.HasPrefix(rawTable, "#") {
		// A temp table's rows come from the batch's own statements, which are
		// other calls to ResolveQuerySpan, so the span marks the relation
		// untraced rather than tracing through them. That holds for one the
		// batch declared and for one it did not: the caller's map carries the
		// columns, not what filled them.
		if tempTable, ok := q.findTempTable(rawTable); ok {
			return (*physicalTable)(tempTable), nil
		}
		return &pseudoTable{}, nil
	}

	// SQL Server: CTEs shadow physical tables — check first, nearest match wins.
	if rawDatabase == "" && rawSchema == "" {
		for _, cte := range q.ctes {
			if q.isIdentifierEqual(rawTable, cte.Name) {
				return cte, nil
			}
		}
	}

	database := q.defaultDatabase
	if rawDatabase != "" {
		database = rawDatabase
	}
	schema := q.defaultSchema
	if rawSchema != "" {
		schema = rawSchema
	}

	rel, err := q.findRelation(database, schema, rawTable)
	if err != nil {
		return nil, err
	}
	switch {
	case rel == nil:
	case rel.view != nil:
		viewColumns, err := q.getColumnsFromCreateView(rel.view.Definition, rel.database, rel.schema, rel.name)
		if err != nil {
			return nil, fmt.Errorf("failed to get columns for view %s.%s.%s: %w", rel.database, rel.schema, rel.name, err)
		}
		return &physicalView{
			Database: rel.database,
			Schema:   rel.schema,
			Name:     rel.name,
			Columns:  viewColumns,
		}, nil
	default:
		columns := make([]string, 0, len(rel.table.GetColumns()))
		for _, c := range rel.table.GetColumns() {
			columns = append(columns, c.Name)
		}
		return &physicalTable{
			Database: rel.database,
			Schema:   rel.schema,
			Name:     rel.name,
			Columns:  columns,
		}, nil
	}
	return nil, &NotFoundError{
		Database: &database,
		Schema:   &schema,
		Table:    &rawTable,
	}
}

// catalogRelation is a table or a view of the catalog, named as the catalog
// names it.
type catalogRelation struct {
	database, schema, name string
	table                  *metadata.TableMetadata // nil for a view
	view                   *metadata.ViewMetadata  // nil for a table
}

// findRelation looks a table or view up in the catalog. An empty database or
// schema matches any. It returns nil for a name the catalog does not have, and
// an error for a database the catalog lists but has not synced.
func (q *querySpanExtractor) findRelation(database, schema, name string) (*catalogRelation, error) {
	allDatabases, err := q.catalog.ListDatabaseNames()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %v", err)
	}
	for _, databaseName := range allDatabases {
		if database != "" && !q.isIdentifierEqual(database, databaseName) {
			continue
		}
		databaseMeta, err := q.catalog.GetDatabase(databaseName)
		if err != nil {
			return nil, fmt.Errorf("failed to get database %s metadata: %w", databaseName, err)
		}
		if databaseMeta == nil {
			return nil, fmt.Errorf("database metadata for database %q not found (database not synced)", databaseName)
		}
		for _, schemaName := range databaseMeta.ListSchemaNames() {
			if schema != "" && !q.isIdentifierEqual(schema, schemaName) {
				continue
			}
			schemaSchema := databaseMeta.GetSchema(schemaName)
			for _, tableName := range schemaSchema.ListTableNames() {
				if q.isIdentifierEqual(name, tableName) {
					table := schemaSchema.GetTable(tableName)
					return &catalogRelation{database: databaseName, schema: schemaName, name: table.Name, table: table}, nil
				}
			}
			for _, viewName := range schemaSchema.ListViewNames() {
				if q.isIdentifierEqual(name, viewName) {
					view := schemaSchema.GetView(viewName)
					return &catalogRelation{database: databaseName, schema: schemaName, name: view.Name, view: view}, nil
				}
			}
		}
	}
	return nil, nil
}

// resolveAccessTables names each table and view the statement accesses as
// the catalog names it, which is the name a grant or a policy on it holds:
// SQL Server reads dbo.SALARIES as dbo.salaries. A name the catalog does not
// have, or that cannot be looked up, stays as written; resolving the
// statement's columns reports the failure. A view is not expanded into the
// tables of its own database: SQL Server checks a reader's permission on the
// view alone when its owner owns the tables too. It is expanded into what it
// reads in another database or on another server, since ownership chains do
// not cross databases by default (the `cross db ownership chaining` server
// option is off), so SQL Server checks the reader's permission there.
func (q *querySpanExtractor) resolveAccessTables(tables ColumnSet) ColumnSet {
	out := make(ColumnSet, len(tables))
	seen := map[string]bool{}
	for r := range tables {
		r, rel := q.resolveAccessTable(r)
		out[r] = true
		q.addViewCrossDatabaseReads(rel, out, seen)
	}
	return out
}

// resolveAccessTable names r as the catalog does and returns the catalog's
// table or view, or nil when the catalog does not have it.
func (q *querySpanExtractor) resolveAccessTable(r ColumnResource) (ColumnResource, *catalogRelation) {
	if r.Server != "" || strings.HasPrefix(r.Table, "#") {
		return r, nil
	}
	rel, _ := q.findRelation(r.Database, r.Schema, r.Table)
	if rel != nil {
		r.Database, r.Schema, r.Table = rel.database, rel.schema, rel.name
	}
	return r, rel
}

// addViewCrossDatabaseReads adds to out what the view rel reads on another
// server or in a user database other than its own, through any view in
// between. A view in seen was expanded already, so a cycle ends; resolving
// the view's columns reports it, and a definition that does not parse.
func (q *querySpanExtractor) addViewCrossDatabaseReads(rel *catalogRelation, out ColumnSet, seen map[string]bool) {
	if rel == nil || rel.view == nil {
		return
	}
	key := tsqlViewResolutionKey(rel.database, rel.schema, rel.name)
	if seen[key] {
		return
	}
	seen[key] = true
	stmts, err := mssql.Parse(rel.view.Definition)
	if err != nil {
		return
	}
	for _, s := range stmts {
		cv, ok := s.AST.(*ast.CreateViewStmt)
		if !ok {
			continue
		}
		viewQ := newQuerySpanExtractor(rel.database, rel.schema, q.catalog, q.tempTables, q.ignoreCaseSensitive)
		for t := range viewQ.collectAccessTables(cv.Query) {
			t, tRel := viewQ.resolveAccessTable(t)
			if t.Server != "" || (!q.isIdentifierEqual(t.Database, rel.database) && !isSystemResource(t, q.ignoreCaseSensitive)) {
				out[t] = true
			}
			viewQ.addViewCrossDatabaseReads(tRel, out, seen)
		}
	}
}

// getColumnsFromCreateView parses a CREATE VIEW definition with omni, extracts
// the body's result columns, and applies the optional column alias list.
func (q *querySpanExtractor) getColumnsFromCreateView(definition string, viewDatabaseName string, viewSchemaName string, viewName string) ([]QuerySpanResult, error) {
	key := tsqlViewResolutionKey(viewDatabaseName, viewSchemaName, viewName)
	if q.viewResolutionStack[key] {
		return nil, fmt.Errorf("cyclic view reference detected while resolving %q", viewName)
	}

	stmts, err := mssql.Parse(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CREATE VIEW definition: %w", err)
	}
	var createView *ast.CreateViewStmt
	for _, s := range stmts {
		if cv, ok := s.AST.(*ast.CreateViewStmt); ok {
			createView = cv
			break
		}
	}
	if createView == nil {
		return nil, fmt.Errorf("no CREATE VIEW statement found in definition")
	}
	body, ok := createView.Query.(*ast.SelectStmt)
	if !ok || body == nil {
		return nil, fmt.Errorf("CREATE VIEW body is not a SELECT")
	}
	var columnAliases []string
	if createView.Columns != nil {
		for _, it := range createView.Columns.Items {
			if s, ok := it.(*ast.String); ok {
				columnAliases = append(columnAliases, s.Str)
			}
		}
	}

	newQ := newQuerySpanExtractor(viewDatabaseName, viewSchemaName, q.catalog, q.tempTables, q.ignoreCaseSensitive)
	newQ.source = definition
	newQ.viewResolutionStack = cloneViewResolutionStack(q.viewResolutionStack)
	newQ.viewResolutionStack[key] = true
	pseudo, err := newQ.extractFromSelectStmt(body)
	if err != nil {
		var resourceNotFound *NotFoundError
		if errors.As(err, &resourceNotFound) {
			return nil, resourceNotFound
		}
		return nil, fmt.Errorf("failed to extract view body: %w", err)
	}
	results := pseudo.GetQuerySpanResult()

	if len(columnAliases) > 0 && len(columnAliases) == len(results) {
		for i, alias := range columnAliases {
			results[i].Name = alias
		}
	}
	// The view's own conditions decide which of its rows the caller reads.
	q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, newQ.predicateColumns)
	return results, nil
}

func tsqlViewResolutionKey(databaseName, schemaName, viewName string) string {
	return databaseName + "\x00" + schemaName + "\x00" + viewName
}

func cloneViewResolutionStack(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// tsqlGetAllFieldsOfTableInFromOrOuterCTE expands `table.*` to the column set
// of the matching table source in scope.
func (q *querySpanExtractor) tsqlGetAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName string) ([]QuerySpanResult, error) {
	type maskType = uint8
	const (
		maskNone         maskType = 0
		maskDatabaseName maskType = 1 << iota
		maskSchemaName
		maskTableName
	)
	mask := maskNone
	if normalizedTableName != "" {
		mask |= maskTableName
	}
	if normalizedSchemaName != "" {
		if mask&maskTableName == 0 {
			return nil, fmt.Errorf(`table name %s is specified without column name`, normalizedTableName)
		}
		mask |= maskSchemaName
	}
	if normalizedDatabaseName != "" {
		if mask&maskSchemaName == 0 {
			return nil, fmt.Errorf(`database name %s is specified without schema name`, normalizedDatabaseName)
		}
		mask |= maskDatabaseName
	}

	for _, tableSource := range q.tableSourcesFrom {
		if mask&maskDatabaseName != 0 && !q.isIdentifierEqual(normalizedDatabaseName, tableSource.GetDatabaseName()) {
			continue
		}
		if mask&maskSchemaName != 0 && !q.isIdentifierEqual(normalizedSchemaName, tableSource.GetSchemaName()) {
			continue
		}
		if mask&maskTableName != 0 && !q.isIdentifierEqual(normalizedTableName, tableSource.GetTableName()) {
			continue
		}
		return tableSource.GetQuerySpanResult(), nil
	}
	return nil, fmt.Errorf(`no matching table %q.%q.%q`, normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
}

// tsqlIsFieldSensitive resolves a (database, schema, table, column) tuple to a
// QuerySpanResult by scanning tableSourcesFrom first, then outerTableSources
// (nearest-first) for correlated-subquery resolution.
func (q *querySpanExtractor) tsqlIsFieldSensitive(normalizedDatabaseName string, normalizedSchemaName string, normalizedTableName string, normalizedColumnName string) (QuerySpanResult, error) {
	type maskType = uint8
	const (
		maskNone         maskType = 0
		maskDatabaseName maskType = 1 << iota
		maskSchemaName
		maskTableName
		maskColumnName
	)
	mask := maskNone
	if normalizedColumnName != "" {
		mask |= maskColumnName
	}
	if normalizedTableName != "" {
		if mask&maskColumnName == 0 {
			return QuerySpanResult{}, fmt.Errorf(`table name %s is specified without column name`, normalizedTableName)
		}
		mask |= maskTableName
	}
	if normalizedSchemaName != "" {
		if mask&maskTableName == 0 {
			return QuerySpanResult{}, fmt.Errorf(`schema name %s is specified without table name`, normalizedSchemaName)
		}
		mask |= maskSchemaName
	}
	if normalizedDatabaseName != "" {
		if mask&maskSchemaName == 0 {
			return QuerySpanResult{}, fmt.Errorf(`database name %s is specified without schema name`, normalizedDatabaseName)
		}
		mask |= maskDatabaseName
	}

	if mask == maskNone {
		return QuerySpanResult{}, fmt.Errorf(`no object name is specified`)
	}

	// The first table source of the FROM clause that has the column is the
	// one it reads. A column that names no table and that two of them have is
	// ambiguous, which SQL Server refuses (error 209):
	//   SELECT C1 FROM T1, T2;
	// Users can specify the table name to disambiguate:
	//   SELECT T1.C1 FROM T1, T2;
	// If an alias is set the original table name is shadowed:
	//   SELECT T1.C1 FROM T1 AS T3, T2;  -- invalid
	var (
		found     QuerySpanResult
		foundFrom tableSource
	)
	for _, tableSource := range q.tableSourcesFrom {
		if mask&maskDatabaseName != 0 && !q.isIdentifierEqual(normalizedDatabaseName, tableSource.GetDatabaseName()) {
			continue
		}
		if mask&maskSchemaName != 0 && !q.isIdentifierEqual(normalizedSchemaName, tableSource.GetSchemaName()) {
			continue
		}
		if mask&maskTableName != 0 && !q.isIdentifierEqual(normalizedTableName, tableSource.GetTableName()) {
			continue
		}
		for _, column := range tableSource.GetQuerySpanResult() {
			if mask&maskColumnName != 0 && !q.isIdentifierEqual(normalizedColumnName, column.Name) {
				continue
			}
			switch {
			case foundFrom == nil:
				found, foundFrom = column, tableSource
			case foundFrom != tableSource && mask&maskTableName == 0:
				return QuerySpanResult{}, &ambiguousColumnError{column: normalizedColumnName}
			}
		}
	}
	if foundFrom != nil {
		return found, nil
	}

	for i := len(q.outerTableSources) - 1; i >= 0; i-- {
		tableSource := q.outerTableSources[i]
		if mask&maskDatabaseName != 0 && !q.isIdentifierEqual(normalizedDatabaseName, tableSource.GetDatabaseName()) {
			continue
		}
		if mask&maskSchemaName != 0 && !q.isIdentifierEqual(normalizedSchemaName, tableSource.GetSchemaName()) {
			continue
		}
		if mask&maskTableName != 0 && !q.isIdentifierEqual(normalizedTableName, tableSource.GetTableName()) {
			continue
		}
		for _, column := range tableSource.GetQuerySpanResult() {
			if mask&maskColumnName != 0 && !q.isIdentifierEqual(normalizedColumnName, column.Name) {
				continue
			}
			return column, nil
		}
	}
	return QuerySpanResult{}, &unknownColumnError{database: normalizedDatabaseName, schema: normalizedSchemaName, table: normalizedTableName, column: normalizedColumnName}
}

// unknownColumnError is a column reference no table source in scope has,
// which SQL Server refuses (error 207, or 4104 for a qualified one).
type unknownColumnError struct{ database, schema, table, column string }

func (e *unknownColumnError) Error() string {
	return fmt.Sprintf(`no matching column %q.%q.%q.%q`, e.database, e.schema, e.table, e.column)
}

// ambiguousColumnError is a column reference that names no table and that two
// table sources of the FROM clause have, which SQL Server refuses (error 209).
type ambiguousColumnError struct{ column string }

func (e *ambiguousColumnError) Error() string {
	return fmt.Sprintf("ambiguous column name %q", e.column)
}

// isIdentifierEqual compares identifiers using the extractor's case sensitivity flag.
func (q *querySpanExtractor) isIdentifierEqual(a, b string) bool {
	if !q.ignoreCaseSensitive {
		return a == b
	}
	if len(a) != len(b) {
		return false
	}
	runeA, runeB := []rune(a), []rune(b)
	for i := 0; i < len(runeA); i++ {
		if unicode.ToLower(runeA[i]) != unicode.ToLower(runeB[i]) {
			return false
		}
	}
	return true
}

// unionTableSources unions two or more table sources column-wise.
func unionTableSources(tableSources ...tableSource) ([]QuerySpanResult, error) {
	if len(tableSources) == 0 {
		return nil, errors.New("no table source to union")
	}

	anchor := tableSources[0].GetQuerySpanResult()
	for i := 1; i < len(tableSources); i++ {
		current := tableSources[i].GetQuerySpanResult()
		if len(current) != len(anchor) {
			return nil, fmt.Errorf("the %dth table source has different column number with previous anchors, previous: %d, current %d", i+1, len(anchor), len(current))
		}
		for j := range anchor {
			anchor[j].SourceColumns, _ = mergeColumnSet(anchor[j].SourceColumns, current[j].SourceColumns)
		}
	}
	return anchor, nil
}

// isMixedQuery reports (allSystems, mixed): mixed=true means the query touches
// both user tables and system tables simultaneously, which we treat as an
// error. allSystems=true means every accessed table is a system table.
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

func isSystemResource(resource ColumnResource, ignoreCaseSensitive bool) bool {
	if IsSystemDatabase(resource.Database, !ignoreCaseSensitive) {
		return true
	}
	if IsSystemSchema(resource.Schema, !ignoreCaseSensitive) {
		return true
	}
	return false
}

func (q *querySpanExtractor) getQuerySpan(statement string) (*ResolvedSpan, error) {
	q.source = statement

	stmts, err := mssql.Parse(statement)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 {
		return &ResolvedSpan{
			Type:          Select,
			SourceColumns: make(ColumnSet),
			Results:       []QuerySpanResult{},
		}, nil
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}

	root := stmts[0].AST
	if root == nil {
		// An empty statement.
		return &ResolvedSpan{
			Type:          Select,
			SourceColumns: make(ColumnSet),
			Results:       []QuerySpanResult{},
		}, nil
	}
	if _, isGo := root.(*ast.GoStmt); isGo {
		// A bare "GO" batch separator.
		return &ResolvedSpan{
			Type:          Select,
			SourceColumns: make(ColumnSet),
			Results:       []QuerySpanResult{},
		}, nil
	}

	// DECLARE @t TABLE(...) registers @t in the temp tables so later
	// statements of the batch can resolve it. A scalar DECLARE is a no-op.
	if d, ok := root.(*ast.DeclareStmt); ok {
		populateTempTables(d, q.tempTables)
	}
	if set, ok := root.(*ast.SetOptionStmt); ok {
		q.setOption = set
	}
	if c, ok := root.(*ast.CreateTableStmt); ok {
		populateCreateTempTable(c, q.tempTables)
	}

	accessTables := q.collectAccessTables(root)
	allSystems, mixed := isMixedQuery(accessTables, q.ignoreCaseSensitive)
	if mixed {
		return nil, ErrMixUserSystemTables
	}
	if !allSystems {
		accessTables = q.resolveAccessTables(accessTables)
	}

	queryType := classifyQueryType(root, allSystems)
	sel, isSelect := root.(*ast.SelectStmt)
	// classifyQueryType maps SET/DECLARE/SetOption to Select for historical
	// reasons; those have no result columns. Return early if the root isn't a
	// real SELECT.
	if queryType != Select || !isSelect {
		return &ResolvedSpan{
			Type:          queryType,
			SourceColumns: accessTables,
			Results:       []QuerySpanResult{},
			Untraced:      q.untracedError(root),
		}, nil
	}

	pseudo, err := q.extractFromSelectStmt(sel)
	if err != nil {
		var resourceNotFound *NotFoundError
		if errors.As(err, &resourceNotFound) {
			return &ResolvedSpan{
				Type:          Select,
				SourceColumns: accessTables,
				Results:       []QuerySpanResult{},
				NotFoundError: resourceNotFound,
			}, nil
		}
		return nil, err
	}

	return &ResolvedSpan{
		Type:             Select,
		SourceColumns:    accessTables,
		Results:          pseudo.GetQuerySpanResult(),
		PredicateColumns: q.predicateColumns,
		Untraced:         q.untracedError(root),
	}, nil
}

// returnsRows reports whether the statement answers with a result set: a
// SELECT that writes no table of its own, whether it reads user tables or only
// the system's.
func returnsRows(root ast.Node) bool {
	sel, ok := root.(*ast.SelectStmt)
	return ok && selectIntoTarget(sel) == nil
}

// -------------------- SELECT processing --------------------

// extractFromSelectStmt processes a full SELECT statement (including set-ops,
// WITH clause) and returns the result column set as a PseudoTable.
func (q *querySpanExtractor) extractFromSelectStmt(sel *ast.SelectStmt) (*pseudoTable, error) {
	if f := documentClause(sel); f != nil {
		// FOR XML and FOR JSON return the rows as one document in one column,
		// which SQL Server names for the mode, so the column holds every
		// column of the select list.
		rows, err := q.extractFromSelectStmt(withoutDocumentClause(sel))
		if err != nil {
			return nil, err
		}
		doc := QuerySpanResult{Name: forClauseColumnName(f), SourceColumns: make(ColumnSet)}
		for _, c := range rows.Columns {
			doc.SourceColumns, _ = mergeColumnSet(doc.SourceColumns, c.SourceColumns)
		}
		return &pseudoTable{Columns: []QuerySpanResult{doc}}, nil
	}

	// WITH clause: populate CTE list BEFORE dispatching set-ops so that both
	// arms of a `WITH cte AS (...) SELECT ... FROM cte UNION SELECT ... FROM cte`
	// see the CTE in scope. The defer restores q.ctes on exit.
	ctesSnapshot := len(q.ctes)
	defer func() { q.ctes = q.ctes[:ctesSnapshot] }()
	if sel.WithClause != nil {
		if err := q.processWithClause(sel.WithClause); err != nil {
			return nil, err
		}
	}

	// Set operations: recurse into arms and union. The WITH clause above is
	// already in q.ctes, so both arms resolve CTE references correctly.
	if sel.Op != ast.SetOpNone && (sel.Larg != nil || sel.Rarg != nil) {
		return q.extractFromSetOp(sel)
	}

	// FROM clause — a SELECT without FROM is legal (e.g. SELECT 1).
	tableSourcesSnapshot := len(q.tableSourcesFrom)
	defer func() { q.tableSourcesFrom = q.tableSourcesFrom[:tableSourcesSnapshot] }()
	if sel.FromClause != nil {
		for _, item := range sel.FromClause.Items {
			ts, err := q.extractTableSource(item)
			if err != nil {
				return nil, err
			}
			q.tableSourcesFrom = append(q.tableSourcesFrom, ts...)
		}
	}

	// WHERE / HAVING predicate columns — resolved in THIS select's scope.
	// Subquery expressions encountered during the walk get cloned extractors;
	// their results merge back into q.predicateColumns. Errors (including
	// ResourceNotFoundError from a subquery referencing a missing table)
	// propagate up to getQuerySpan which converts them into a span with
	// NotFoundError populated.
	if err := q.collectPredicatesInScope(sel.WhereClause); err != nil {
		return nil, err
	}
	if err := q.collectPredicatesInScope(sel.HavingClause); err != nil {
		return nil, err
	}
	if err := q.collectRowCountPredicates(sel); err != nil {
		return nil, err
	}

	// Target list.
	results, err := q.extractTargetList(sel.TargetList)
	if err != nil {
		return nil, err
	}
	q.populateSelectIntoTempTable(sel, results)

	return &pseudoTable{Columns: results}, nil
}

// documentClause returns the FOR XML or FOR JSON clause of a query: its own,
// or, for a set operation, its last operand's, where the parser keeps the
// clause that serializes the whole result. FOR BROWSE returns the rows as
// they are, and is not one.
func documentClause(sel *ast.SelectStmt) *ast.ForClause {
	last := sel
	for last.Op != ast.SetOpNone && last.Rarg != nil {
		last = last.Rarg
	}
	if f := last.ForClause; f != nil && (f.Mode == ast.ForXML || f.Mode == ast.ForJSON) {
		return f
	}
	return nil
}

// withoutDocumentClause copies sel, down to the operand documentClause found
// the clause on, without it.
func withoutDocumentClause(sel *ast.SelectStmt) *ast.SelectStmt {
	cp := *sel
	if cp.Op != ast.SetOpNone && cp.Rarg != nil {
		cp.Rarg = withoutDocumentClause(cp.Rarg)
	} else {
		cp.ForClause = nil
	}
	return &cp
}

// forClauseColumnName is the name SQL Server gives the one column of a FOR XML
// or FOR JSON query.
func forClauseColumnName(f *ast.ForClause) string {
	if f.Mode == ast.ForJSON {
		return "JSON_F52E2B61-18A1-11d1-B105-00805F49916B"
	}
	return "XML_F52E2B61-18A1-11d1-B105-00805F49916B"
}

// collectPredicatesInScope walks a boolean expression in the extractor's
// current scope: plain ColumnRefs resolve against q.tableSourcesFrom +
// q.outerTableSources; subquery nodes trigger a scope clone, and the
// subquery's output-column sources plus its own predicate-column set merge
// back into q.predicateColumns.
//
// Errors from nested subquery extraction (ResourceNotFoundError or real
// parser errors) propagate so the caller can surface them via the top-level
// NotFoundError path instead of producing a silently-partial span.
func (q *querySpanExtractor) collectPredicatesInScope(expr ast.Node) error {
	if expr == nil {
		return nil
	}
	switch v := expr.(type) {
	case *ast.ColumnRef:
		if r, handled, err := q.resolvePseudoColumnRef(v); handled {
			if err != nil {
				return fmt.Errorf("failed to resolve predicate column %q: %w", v.Column, err)
			}
			q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, r.SourceColumns)
			return nil
		}
		r, err := q.tsqlIsFieldSensitive(v.Database, v.Schema, v.Table, v.Column)
		if err != nil {
			return fmt.Errorf("failed to resolve predicate column %q: %w", v.Column, err)
		}
		q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, r.SourceColumns)
		return nil
	case *ast.FullTextPredicate:
		if v.Columns != nil {
			for _, it := range v.Columns.Items {
				if cr, ok := it.(*ast.ColumnRef); ok {
					if err := q.collectPredicatesInScope(cr); err != nil {
						return err
					}
				}
			}
		}
		if err := q.collectPredicatesInScope(v.Value); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.LanguageTerm)
	case *ast.GraphMatchPredicate:
		// A graph pattern names node and edge aliases, whose columns are the
		// predicate's. Each name is read as a column, the way a graph
		// pseudocolumn is (query_span_pseudo_column.go).
		return q.collectGraphMatchColumns(v.Expression)
	case *ast.SubqueryExpr:
		return q.mergeSubqueryIntoPredicates(v.Query)
	case *ast.ExistsExpr:
		return q.mergeSubqueryIntoPredicates(v.Subquery)
	case *ast.SubqueryComparisonExpr:
		if err := q.collectPredicatesInScope(v.Left); err != nil {
			return err
		}
		return q.mergeSubqueryIntoPredicates(v.Subquery)
	case *ast.InExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		if v.List != nil {
			for _, it := range v.List.Items {
				if err := q.collectPredicatesInScope(it); err != nil {
					return err
				}
			}
		}
		return q.collectPredicatesInScope(v.Subquery)
	case *ast.BinaryExpr:
		if err := q.collectPredicatesInScope(v.Left); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.Right)
	case *ast.UnaryExpr:
		return q.collectPredicatesInScope(v.Operand)
	case *ast.BetweenExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		if err := q.collectPredicatesInScope(v.Low); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.High)
	case *ast.LikeExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		if err := q.collectPredicatesInScope(v.Pattern); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.Escape)
	case *ast.IsExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		// IS [NOT] DISTINCT FROM keeps its right operand in Value.
		return q.collectPredicatesInScope(v.Value)
	case *ast.CaseExpr:
		if err := q.collectPredicatesInScope(v.Arg); err != nil {
			return err
		}
		if err := q.collectPredicatesInScope(v.ElseExpr); err != nil {
			return err
		}
		if v.WhenList != nil {
			for _, it := range v.WhenList.Items {
				if cw, ok := it.(*ast.CaseWhen); ok {
					if err := q.collectPredicatesInScope(cw.Condition); err != nil {
						return err
					}
					if err := q.collectPredicatesInScope(cw.Result); err != nil {
						return err
					}
				}
			}
		}
		return nil
	case *ast.IifExpr:
		if err := q.collectPredicatesInScope(v.Condition); err != nil {
			return err
		}
		if err := q.collectPredicatesInScope(v.TrueVal); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.FalseVal)
	case *ast.CoalesceExpr:
		if v.Args != nil {
			for _, it := range v.Args.Items {
				if err := q.collectPredicatesInScope(it); err != nil {
					return err
				}
			}
		}
		return nil
	case *ast.NullifExpr:
		if err := q.collectPredicatesInScope(v.Left); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.Right)
	case *ast.FuncCallExpr:
		if v.Args != nil {
			for _, it := range v.Args.Items {
				if err := q.collectPredicatesInScope(it); err != nil {
					return err
				}
			}
		}
		return nil
	case *ast.MethodCallExpr:
		// An instance method reads its receiver too.
		if v.Expr != nil {
			if err := q.collectPredicatesInScope(v.Expr); err != nil {
				return err
			}
		}
		if v.Args != nil {
			for _, it := range v.Args.Items {
				if err := q.collectPredicatesInScope(it); err != nil {
					return err
				}
			}
		}
		return nil
	case *ast.CastExpr:
		return q.collectPredicatesInScope(v.Expr)
	case *ast.ConvertExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.Style)
	case *ast.TryCastExpr:
		return q.collectPredicatesInScope(v.Expr)
	case *ast.TryConvertExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.Style)
	case *ast.ParenExpr:
		return q.collectPredicatesInScope(v.Expr)
	case *ast.CollateExpr:
		return q.collectPredicatesInScope(v.Expr)
	case *ast.AtTimeZoneExpr:
		if err := q.collectPredicatesInScope(v.Expr); err != nil {
			return err
		}
		return q.collectPredicatesInScope(v.TimeZone)
	default:
		r, err := q.resolveExpressionNode(expr)
		if err != nil {
			return err
		}
		q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, r.SourceColumns)
		return nil
	}
}

// collectRowCountPredicates records what TOP, OFFSET, and FETCH read, since
// they decide how many rows come back, and the conditions of subqueries in
// ORDER BY. The ORDER BY keys themselves only order the rows.
func (q *querySpanExtractor) collectRowCountPredicates(sel *ast.SelectStmt) error {
	if sel.Top != nil {
		if err := q.collectPredicatesInScope(sel.Top.Count); err != nil {
			return err
		}
	}
	if err := q.collectPredicatesInScope(sel.OffsetClause); err != nil {
		return err
	}
	if sel.FetchClause != nil {
		if err := q.collectPredicatesInScope(sel.FetchClause.Count); err != nil {
			return err
		}
	}
	if sel.OrderByClause == nil {
		return nil
	}
	var err error
	ast.Inspect(sel.OrderByClause, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		sub, ok := n.(*ast.SelectStmt)
		if !ok {
			return true
		}
		clone := q.cloneForSubquery()
		if _, err = clone.extractFromSelectStmt(sub); err == nil {
			q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
		}
		return false
	})
	return err
}

func (q *querySpanExtractor) mergeSubqueryIntoPredicates(body ast.Node) error {
	sel, ok := body.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	clone := q.cloneForSubquery()
	pseudo, err := clone.extractFromSelectStmt(sel)
	if err != nil {
		return err
	}
	// Subquery output columns become outer-scope predicate columns.
	for _, c := range pseudo.GetQuerySpanResult() {
		q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, c.SourceColumns)
	}
	// And the clone's own predicate columns (from its WHERE etc.) propagate up.
	q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
	return nil
}

func (q *querySpanExtractor) extractFromSetOp(sel *ast.SelectStmt) (*pseudoTable, error) {
	// Omni puts the WITH clause on the leftmost set-op arm (the parser attaches
	// the WITH to the first SELECT expression, not to the outer Op node). Lift
	// it into the current scope so the CTE is visible in BOTH arms, then
	// recurse into a copy of Larg with WithClause cleared to avoid double
	// processing.
	larg := sel.Larg
	if larg != nil && larg.WithClause != nil {
		if err := q.processWithClause(larg.WithClause); err != nil {
			return nil, err
		}
		cp := *larg
		cp.WithClause = nil
		larg = &cp
	}
	left, err := q.extractFromSelectStmt(larg)
	if err != nil {
		return nil, err
	}
	right, err := q.extractFromSelectStmt(sel.Rarg)
	if err != nil {
		return nil, err
	}
	if sel.Op == ast.SetOpIntersect || sel.Op == ast.SetOpExcept {
		// INTERSECT and EXCEPT keep a row by whether the other arm has it, so
		// both arms' output columns decide which rows come back.
		for _, arm := range []tableSource{left, right} {
			for _, c := range arm.GetQuerySpanResult() {
				q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, c.SourceColumns)
			}
		}
	}
	merged, err := unionTableSources(left, right)
	if err != nil {
		return nil, err
	}
	return &pseudoTable{Columns: merged}, nil
}

func (q *querySpanExtractor) processWithClause(w *ast.WithClause) error {
	if w.CTEs == nil {
		return nil
	}
	for _, item := range w.CTEs.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		body, ok := cte.Query.(*ast.SelectStmt)
		if !ok {
			return unsupportedNodeError("only SELECT CTE bodies are supported", cte.Query)
		}

		var columnAliases []string
		if cte.Columns != nil {
			for _, it := range cte.Columns.Items {
				if s, ok := it.(*ast.String); ok {
					columnAliases = append(columnAliases, s.Str)
				}
			}
		}
		cols, predicates, err := q.extractCTEBody(cte.Name, body, columnAliases)
		if err != nil {
			return err
		}
		q.ctes = append(q.ctes, &pseudoTable{Name: cte.Name, Columns: cols})
		q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, predicates)
	}
	return nil
}

// extractCTEBody extracts the column set of a CTE body, named by its column
// list when it has one. For recursive CTEs (body is a set-op SelectStmt whose
// recursive arm references the CTE name), the extractor iterates until the
// column sources stabilize; the recursive arm reads the CTE by its column
// list's names, so the list applies to the anchor first.
func (q *querySpanExtractor) extractCTEBody(cteName string, body *ast.SelectStmt, columnAliases []string) ([]QuerySpanResult, ColumnSet, error) {
	newClone := func() *querySpanExtractor {
		return &querySpanExtractor{
			defaultDatabase:     q.defaultDatabase,
			defaultSchema:       q.defaultSchema,
			ignoreCaseSensitive: q.ignoreCaseSensitive,
			catalog:             q.catalog,
			tempTables:          q.tempTables,
			ctes:                append([]*pseudoTable{}, q.ctes...),
			outerTableSources:   nil,
			predicateColumns:    make(ColumnSet),
			viewResolutionStack: cloneViewResolutionStack(q.viewResolutionStack),
			source:              q.source,
		}
	}

	// Only UNION can be recursive; extract every other CTE directly.
	if body.Op != ast.SetOpUnion || (body.Larg == nil && body.Rarg == nil) {
		clone := newClone()
		pseudo, err := clone.extractFromSelectStmt(body)
		if err != nil {
			return nil, nil, err
		}
		cols := pseudo.GetQuerySpanResult()
		if err := applyColumnAliases(cols, columnAliases); err != nil {
			return nil, nil, fmt.Errorf("CTE %q: %w", cteName, err)
		}
		return cols, clone.predicateColumns, nil
	}

	// Set-op CTE: extract anchor (Larg), register placeholder, then iterate Rarg
	// until the source column sets for the CTE columns stabilize.
	anchorClone := newClone()
	anchorPseudo, err := anchorClone.extractFromSelectStmt(body.Larg)
	if err != nil {
		return nil, nil, err
	}
	cols := append([]QuerySpanResult{}, anchorPseudo.GetQuerySpanResult()...)
	if err := applyColumnAliases(cols, columnAliases); err != nil {
		return nil, nil, fmt.Errorf("CTE %q: %w", cteName, err)
	}
	predicates := ColumnSet{}
	predicates, _ = mergeColumnSet(predicates, anchorClone.predicateColumns)

	for iter := 0; iter < 16; iter++ {
		placeholder := &pseudoTable{Name: cteName, Columns: cols}
		clone := newClone()
		clone.ctes = append(clone.ctes, placeholder)
		rarg, rerr := clone.extractFromSelectStmt(body.Rarg)
		if rerr != nil {
			// Recursive arm failed to resolve (e.g. ResourceNotFoundError from
			// a missing table). Propagate — the top-level getQuerySpan
			// converts ResourceNotFoundError into a span with NotFoundError
			// populated so callers don't see a silently-partial lineage.
			return nil, nil, rerr
		}
		rcols := rarg.GetQuerySpanResult()
		if len(rcols) != len(cols) {
			return nil, nil, fmt.Errorf("CTE %q recursive arm returns %d columns, anchor returns %d", cteName, len(rcols), len(cols))
		}
		changed := false
		for i := range rcols {
			merged, anyChange := mergeColumnSet(cols[i].SourceColumns, rcols[i].SourceColumns)
			cols[i].SourceColumns = merged
			if anyChange {
				changed = true
			}
		}
		predicates, _ = mergeColumnSet(predicates, clone.predicateColumns)
		if !changed {
			break
		}
	}
	return cols, predicates, nil
}

// -------------------- FROM / table source --------------------

// extractTableSource returns the TableSource(s) produced by a single FROM item.
// A JoinClause contributes multiple sources (left + right flattened).
func (q *querySpanExtractor) extractTableSource(node ast.Node) ([]tableSource, error) {
	switch v := node.(type) {
	case *ast.JoinClause:
		return q.extractFromJoin(v)
	case *ast.AdHocTableRef:
		// OPENDATASOURCE: a provider reads the relation, not the connected
		// instance.
		return nil, &LinkedServerError{}
	case *ast.TableRef:
		ts, err := q.resolveTableRef(v)
		if err != nil {
			return nil, err
		}
		return []tableSource{ts}, nil
	case *ast.AliasedTableRef:
		ts, err := q.resolveAliasedTableRef(v)
		if err != nil {
			return nil, err
		}
		return []tableSource{ts}, nil
	case *ast.SubqueryExpr:
		ts, err := q.resolveSubqueryAsTable(v, "")
		if err != nil {
			return nil, err
		}
		return []tableSource{ts}, nil
	case *ast.ValuesClause:
		ts, err := q.resolveValuesClause(v, "", nil)
		if err != nil {
			return nil, err
		}
		return []tableSource{ts}, nil
	case *ast.TableVarRef:
		return q.resolveTableVar(v, v.Alias)
	case *ast.TableVarMethodCallRef:
		return []tableSource{&pseudoTable{Name: v.Alias, Columns: xmlNodesColumns(v.Columns)}}, nil
	case *ast.PivotExpr:
		return nil, unsupportedTableSourceError(node)
	case *ast.UnpivotExpr:
		return nil, unsupportedTableSourceError(node)
	case *ast.FuncCallExpr:
		ts, err := q.resolveTableValuedFunction(v, "", nil)
		if err != nil {
			return nil, err
		}
		return []tableSource{ts}, nil
	default:
		return nil, unsupportedTableSourceError(node)
	}
}

func unsupportedTableSourceError(node ast.Node) error {
	return unsupportedNodeError("only full table name, supported TVF, derived table, values clause, and temp table in table source item are supported", node)
}

func unsupportedNodeError(message string, node ast.Node) error {
	typeName := "<nil>"
	if typ := reflect.TypeOf(node); typ != nil {
		typeName = typ.String()
	}
	return &TypeNotSupportedError{
		Err:  errors.New(message),
		Type: typeName,
	}
}

func (q *querySpanExtractor) extractFromJoin(j *ast.JoinClause) ([]tableSource, error) {
	tableSourcesSnapshot := len(q.tableSourcesFrom)
	defer func() { q.tableSourcesFrom = q.tableSourcesFrom[:tableSourcesSnapshot] }()
	left, err := q.extractTableSource(j.Left)
	if err != nil {
		return nil, err
	}
	if j.Type == ast.JoinCrossApply || j.Type == ast.JoinOuterApply {
		q.tableSourcesFrom = append(q.tableSourcesFrom, left...)
	}
	right, err := q.extractTableSource(j.Right)
	if err != nil {
		return nil, err
	}
	if j.Condition != nil {
		q.tableSourcesFrom = append(q.tableSourcesFrom[:tableSourcesSnapshot], left...)
		q.tableSourcesFrom = append(q.tableSourcesFrom, right...)
		if err := q.collectPredicatesInScope(j.Condition); err != nil {
			return nil, err
		}
	}
	return append(left, right...), nil
}

func (q *querySpanExtractor) resolveTableRef(t *ast.TableRef) (tableSource, error) {
	ts, err := q.tsqlFindTableSchemaByParts(t.Server, t.Database, t.Schema, t.Object)
	if err != nil {
		return nil, err
	}
	if t.Alias != "" {
		ts = &pseudoTable{
			Name:    t.Alias,
			Columns: ts.GetQuerySpanResult(),
		}
	}
	return ts, nil
}

func (q *querySpanExtractor) resolveAliasedTableRef(at *ast.AliasedTableRef) (tableSource, error) {
	var columnAliases []string
	if at.Columns != nil {
		for _, it := range at.Columns.Items {
			if s, ok := it.(*ast.String); ok {
				columnAliases = append(columnAliases, s.Str)
			}
		}
	}
	switch inner := at.Table.(type) {
	case *ast.TableRef:
		ts, err := q.tsqlFindTableSchemaByParts(inner.Server, inner.Database, inner.Schema, inner.Object)
		if err != nil {
			return nil, err
		}
		cols := append([]QuerySpanResult{}, ts.GetQuerySpanResult()...)
		if err := applyColumnAliases(cols, columnAliases); err != nil {
			return nil, err
		}
		name := inner.Object
		if at.Alias != "" {
			name = at.Alias
		}
		return &pseudoTable{Name: name, Columns: cols}, nil
	case *ast.SubqueryExpr:
		return q.resolveSubqueryAsTable(inner, at.Alias, columnAliases...)
	case *ast.ValuesClause:
		return q.resolveValuesClause(inner, at.Alias, columnAliases)
	case *ast.TableVarRef:
		inners, err := q.resolveTableVar(inner, at.Alias)
		if err != nil {
			return nil, err
		}
		if len(inners) == 1 {
			return inners[0], nil
		}
		return nil, unsupportedNodeError("aliased table variable must produce a single table source", at.Table)
	case *ast.FuncCallExpr:
		return q.resolveTableValuedFunction(inner, at.Alias, columnAliases)
	case *ast.TableVarMethodCallRef:
		cols := xmlNodesColumns(inner.Columns)
		if err := applyColumnAliases(cols, columnAliases); err != nil {
			return nil, err
		}
		return &pseudoTable{Name: at.Alias, Columns: cols}, nil
	default:
		return nil, unsupportedNodeError("unsupported aliased table source item", at.Table)
	}
}

func (q *querySpanExtractor) resolveTableValuedFunction(fn *ast.FuncCallExpr, alias string, columnAliases []string) (tableSource, error) {
	if fn == nil || fn.Name == nil {
		return nil, unsupportedTableSourceError(fn)
	}
	var columnNames []string
	switch strings.ToUpper(fn.Name.Object) {
	case "OPENQUERY":
		// A pass-through query: the columns are the remote server's, and the
		// statement spells none of them.
		return nil, &LinkedServerError{Server: rowsetServerArgument(fn), PassThrough: true}
	case "OPENROWSET":
		// An ad-hoc provider read, with no linked server to name.
		return nil, &LinkedServerError{}
	case "STRING_SPLIT":
		columnNames = []string{"value"}
		hasOrdinal, err := stringSplitHasOrdinalColumn(fn.Args)
		if err != nil {
			return nil, err
		}
		if hasOrdinal {
			columnNames = append(columnNames, "ordinal")
		}
	case "OPENJSON":
		columnNames = []string{"key", "value", "type"}
	default:
		return nil, unsupportedTableSourceError(fn)
	}

	argSources, err := q.mergeListSources(fn.Args)
	if err != nil {
		return nil, err
	}
	// The arguments decide how many rows the function returns.
	q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, argSources.SourceColumns)
	columns := make([]QuerySpanResult, 0, len(columnNames))
	for _, name := range columnNames {
		columns = append(columns, QuerySpanResult{
			Name:          name,
			SourceColumns: argSources.SourceColumns,
		})
	}
	if strings.EqualFold(fn.Name.Object, "STRING_SPLIT") && len(columns) > 1 {
		columns[1].SourceColumns = make(ColumnSet)
	}
	if err := applyColumnAliases(columns, columnAliases); err != nil {
		return nil, err
	}
	return &pseudoTable{Name: alias, Columns: columns}, nil
}

// rowsetServerArgument returns the linked server an OPENQUERY call names: its
// first argument, which the parser reads as a table reference.
// collectGraphMatchColumns adds the tables a graph pattern names to the
// search condition's predicate columns.
func (q *querySpanExtractor) collectGraphMatchColumns(expr ast.ExprNode) error {
	return q.walkGraphMatchRefs(expr, func(name string) error {
		r, err := q.resolveGraphMatchRef(name)
		if err != nil {
			return err
		}
		q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, r.SourceColumns)
		return nil
	})
}

// walkGraphMatchRefs calls f for each node and edge alias of a graph pattern.
func (q *querySpanExtractor) walkGraphMatchRefs(expr ast.ExprNode, f func(name string) error) error {
	switch v := expr.(type) {
	case nil:
		return nil
	case *ast.GraphMatchNodeExpr:
		return f(v.Node)
	case *ast.GraphMatchCompositeExpr:
		if err := q.walkGraphMatchRefs(v.Left, f); err != nil {
			return err
		}
		if err := f(v.Edge); err != nil {
			return err
		}
		return q.walkGraphMatchRefs(v.Right, f)
	case *ast.GraphMatchLastNodePredicate:
		if err := q.walkGraphMatchRefs(v.Left, f); err != nil {
			return err
		}
		return q.walkGraphMatchRefs(v.Right, f)
	case *ast.GraphMatchRecursivePredicate:
		if err := q.walkGraphMatchRefs(v.Outer, f); err != nil {
			return err
		}
		if v.Expression != nil {
			for _, it := range v.Expression.Items {
				step, ok := it.(ast.ExprNode)
				if !ok {
					continue
				}
				if err := q.walkGraphMatchRefs(step, f); err != nil {
					return err
				}
			}
		}
		return nil
	case *ast.BinaryExpr:
		if err := q.walkGraphMatchRefs(v.Left, f); err != nil {
			return err
		}
		return q.walkGraphMatchRefs(v.Right, f)
	}
	return nil
}

// resolveGraphMatchRef resolves a node or edge alias of a graph MATCH pattern
// to the table source it names. A graph pattern names tables, so the predicate
// filters rows of those tables and its lineage is the union of their columns,
// as a graph pseudo-column's is (resolveGraphPseudoColumn).
func (q *querySpanExtractor) resolveGraphMatchRef(name string) (QuerySpanResult, error) {
	r := QuerySpanResult{Name: name, SourceColumns: make(ColumnSet), IsPlainField: false}
	found := false
	add := func(ts tableSource) error {
		found = true
		for _, column := range ts.GetQuerySpanResult() {
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, column.SourceColumns)
		}
		return nil
	}
	for _, ts := range q.tableSourcesFrom {
		if q.isIdentifierEqual(name, ts.GetTableName()) {
			if err := add(ts); err != nil {
				return QuerySpanResult{}, err
			}
			break
		}
	}
	if !found {
		for i := len(q.outerTableSources) - 1; i >= 0; i-- {
			ts := q.outerTableSources[i]
			if q.isIdentifierEqual(name, ts.GetTableName()) {
				if err := add(ts); err != nil {
					return QuerySpanResult{}, err
				}
				break
			}
		}
	}
	if !found {
		return QuerySpanResult{}, fmt.Errorf("graph node %q names no table in the FROM clause", name)
	}
	// A graph table with no user-defined columns contributes none, and empty
	// SourceColumns downstream means NoneMasker, the fail-open this
	// table-level lineage is meant to prevent: fail closed.
	if len(r.SourceColumns) == 0 {
		return QuerySpanResult{}, fmt.Errorf("cannot derive lineage for graph node %q: table has no columns in the catalog", name)
	}
	return r, nil
}

func rowsetServerArgument(fn *ast.FuncCallExpr) string {
	if fn.Args == nil || fn.Args.Len() == 0 {
		return ""
	}
	if ref, ok := fn.Args.Items[0].(*ast.TableRef); ok {
		return ref.Object
	}
	return ""
}

func stringSplitHasOrdinalColumn(args *ast.List) (bool, error) {
	if args == nil || args.Len() < 3 {
		return false, nil
	}
	arg := args.Items[2]
	for {
		paren, ok := arg.(*ast.ParenExpr)
		if !ok {
			break
		}
		arg = paren.Expr
	}
	literal, ok := arg.(*ast.Literal)
	if !ok {
		return false, unsupportedNodeError("STRING_SPLIT enable_ordinal must be a constant 0, 1, or NULL", arg)
	}
	switch literal.Type {
	case ast.LitInteger:
		switch literal.Ival {
		case 0:
			return false, nil
		case 1:
			return true, nil
		default:
			return false, unsupportedNodeError("STRING_SPLIT enable_ordinal must be 0 or 1", arg)
		}
	case ast.LitNull:
		return false, nil
	default:
		return false, unsupportedNodeError("STRING_SPLIT enable_ordinal must be a constant 0, 1, or NULL", arg)
	}
}

func (q *querySpanExtractor) resolveSubqueryAsTable(sq *ast.SubqueryExpr, alias string, columnAliases ...string) (tableSource, error) {
	body, ok := sq.Query.(*ast.SelectStmt)
	if !ok {
		return nil, unsupportedNodeError("only SELECT subqueries in table source are supported", sq.Query)
	}
	clone := q.cloneForSubquery()
	pseudo, err := clone.extractFromSelectStmt(body)
	if err != nil {
		return nil, err
	}
	cols := append([]QuerySpanResult{}, pseudo.GetQuerySpanResult()...)
	if err := applyColumnAliases(cols, columnAliases); err != nil {
		return nil, err
	}
	q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
	return &pseudoTable{Name: alias, Columns: cols}, nil
}

func (q *querySpanExtractor) resolveValuesClause(vc *ast.ValuesClause, alias string, columnAliases []string) (tableSource, error) {
	name := alias
	if vc.Rows == nil || vc.Rows.Len() == 0 {
		return &pseudoTable{Name: name}, nil
	}
	firstRow, ok := vc.Rows.Items[0].(*ast.List)
	if !ok {
		return &pseudoTable{Name: name}, nil
	}
	cols := make([]QuerySpanResult, 0, firstRow.Len())
	for rowIndex, item := range vc.Rows.Items {
		row, ok := item.(*ast.List)
		if !ok {
			continue
		}
		if row.Len() != firstRow.Len() {
			return nil, fmt.Errorf("VALUES row %d has %d columns, first row has %d", rowIndex+1, row.Len(), firstRow.Len())
		}
		for i, it := range row.Items {
			expr, ok := it.(ast.ExprNode)
			if !ok {
				if rowIndex == 0 {
					cols = append(cols, QuerySpanResult{})
				}
				continue
			}
			r, err := q.resolveExpression(expr)
			if err != nil {
				return nil, err
			}
			if rowIndex == 0 {
				cols = append(cols, r)
				continue
			}
			cols[i].SourceColumns, _ = mergeColumnSet(cols[i].SourceColumns, r.SourceColumns)
		}
	}
	if err := applyColumnAliases(cols, columnAliases); err != nil {
		return nil, err
	}
	return &pseudoTable{Name: name, Columns: cols}, nil
}

func (q *querySpanExtractor) resolveTableVar(tv *ast.TableVarRef, alias string) ([]tableSource, error) {
	if tv == nil {
		return nil, nil
	}
	name := tv.Name
	if temp, ok := q.findTempTable(name); ok {
		tableName := name
		if alias != "" {
			tableName = alias
		}
		cols := make([]QuerySpanResult, 0, len(temp.Columns))
		for _, c := range temp.Columns {
			cols = append(cols, QuerySpanResult{
				Name: c,
				SourceColumns: ColumnSet{
					ColumnResource{
						Server:   temp.Server,
						Database: temp.Database,
						Schema:   temp.Schema,
						Table:    temp.Name,
						Column:   c,
					}: true,
				},
				IsPlainField: true,
			})
		}
		return []tableSource{&pseudoTable{Name: tableName, Columns: cols}}, nil
	}
	return nil, &NotFoundError{
		Table: &name,
		Err:   fmt.Errorf("temp table %s not found", name),
	}
}

// populateTempTables registers each table variable d declares, keyed by the
// variable name including its leading '@'.
func populateTempTables(d *ast.DeclareStmt, out map[string]*TempTable) {
	if d == nil || d.Variables == nil || out == nil {
		return
	}
	for _, it := range d.Variables.Items {
		v, ok := it.(*ast.VariableDecl)
		if !ok || !v.IsTable || v.TableDef == nil {
			continue
		}
		columns := make([]string, 0, v.TableDef.Len())
		for _, cdef := range v.TableDef.Items {
			if cd, ok := cdef.(*ast.ColumnDef); ok {
				columns = append(columns, cd.Name)
			}
		}
		out[v.Name] = &TempTable{
			Name:    v.Name,
			Columns: columns,
		}
	}
}

func populateCreateTempTable(c *ast.CreateTableStmt, out map[string]*TempTable) {
	if c == nil || c.Name == nil || out == nil || !strings.HasPrefix(c.Name.Object, "#") {
		return
	}
	var columns []string
	if c.Columns != nil {
		columns = make([]string, 0, c.Columns.Len())
		for _, item := range c.Columns.Items {
			if cd, ok := item.(*ast.ColumnDef); ok {
				columns = append(columns, cd.Name)
			}
		}
	}
	out[c.Name.Object] = &TempTable{
		Name:    c.Name.Object,
		Columns: columns,
	}
}

func (q *querySpanExtractor) populateSelectIntoTempTable(sel *ast.SelectStmt, results []QuerySpanResult) {
	if sel == nil || sel.IntoTable == nil || q.tempTables == nil || !strings.HasPrefix(sel.IntoTable.Object, "#") {
		return
	}
	columns := make([]string, 0, len(results))
	for i, result := range results {
		name := result.Name
		if name == "" && sel.TargetList != nil && i < len(sel.TargetList.Items) {
			name = q.sliceName(sel.TargetList.Items[i])
		}
		columns = append(columns, name)
	}
	q.tempTables[sel.IntoTable.Object] = &TempTable{
		Name:    sel.IntoTable.Object,
		Columns: columns,
	}
}

func xmlNodesColumns(cols []string) []QuerySpanResult {
	out := make([]QuerySpanResult, 0, len(cols))
	for _, c := range cols {
		out = append(out, QuerySpanResult{Name: c})
	}
	return out
}

func applyColumnAliases(cols []QuerySpanResult, aliases []string) error {
	if len(aliases) == 0 {
		return nil
	}
	if len(aliases) != len(cols) {
		return fmt.Errorf("column alias count %d does not match column count %d", len(aliases), len(cols))
	}
	for i, name := range aliases {
		cols[i].Name = name
	}
	return nil
}

func (q *querySpanExtractor) cloneForSubquery() *querySpanExtractor {
	clone := &querySpanExtractor{
		defaultDatabase:     q.defaultDatabase,
		defaultSchema:       q.defaultSchema,
		ignoreCaseSensitive: q.ignoreCaseSensitive,
		catalog:             q.catalog,
		tempTables:          q.tempTables,
		ctes:                append([]*pseudoTable{}, q.ctes...),
		outerTableSources:   slices.Concat(q.outerTableSources, q.tableSourcesFrom),
		predicateColumns:    make(ColumnSet),
		viewResolutionStack: cloneViewResolutionStack(q.viewResolutionStack),
		source:              q.source,
	}
	return clone
}

// -------------------- Target list --------------------

func (q *querySpanExtractor) extractTargetList(list *ast.List) ([]QuerySpanResult, error) {
	if list == nil {
		return nil, nil
	}
	var results []QuerySpanResult
	for _, item := range list.Items {
		switch v := item.(type) {
		case *ast.ResTarget:
			if star, ok := v.Val.(*ast.StarExpr); ok {
				fields, err := q.expandStar(star)
				if err != nil {
					return nil, err
				}
				results = append(results, fields...)
				continue
			}
			if assign, ok := v.Val.(*ast.SelectAssign); ok {
				// SELECT @v = expr assigns a variable and produces no result
				// column. SQL Server returns no result set for it, and rejects
				// mixing it with a retrieving target (error 141). Resolve the
				// value so the tables it reads stay in the span.
				if _, err := q.resolveExpression(assign.Value); err != nil {
					return nil, err
				}
				continue
			}
			r, err := q.resolveExpression(v.Val)
			if err != nil {
				return nil, err
			}
			r.Name = q.resultName(v.Name, v.Val, r.Name)
			results = append(results, r)
		case *ast.StarExpr:
			fields, err := q.expandStar(v)
			if err != nil {
				return nil, err
			}
			results = append(results, fields...)
		case *ast.SelectAssign:
			// SELECT @v = expr as a bare target: no result column, and the
			// value's tables are already in the span's access set.
			continue
		default:
			if expr, ok := item.(ast.ExprNode); ok {
				r, err := q.resolveExpression(expr)
				if err != nil {
					return nil, err
				}
				r.Name = q.resultName("", expr, r.Name)
				results = append(results, r)
				continue
			}
			return nil, unsupportedNodeError("unsupported SELECT target item", item)
		}
	}
	return results, nil
}

// resultName is the name SQL Server gives a select-list item: its alias; for
// a column reference, the column as the query spells it, not as the catalog
// does, through any parentheses; for a pseudo-column ($IDENTITY,
// IDENTITYCOL, $node_id), the column it resolved to, whose name resolved
// holds; anything else has no name.
func (q *querySpanExtractor) resultName(alias string, val ast.ExprNode, resolved string) string {
	if alias != "" {
		return alias
	}
	for {
		paren, ok := val.(*ast.ParenExpr)
		if !ok {
			break
		}
		val = paren.Expr
	}
	if ref, ok := val.(*ast.ColumnRef); ok {
		if classifyPseudoColumn(ref.Column) != pseudoColumnNone && !q.columnSegmentIsDelimited(ref) {
			return resolved
		}
		return ref.Column
	}
	return ""
}

func (q *querySpanExtractor) expandStar(s *ast.StarExpr) ([]QuerySpanResult, error) {
	if s.Qualifier == "" {
		var out []QuerySpanResult
		for _, ts := range q.tableSourcesFrom {
			out = append(out, ts.GetQuerySpanResult()...)
		}
		return out, nil
	}
	return q.tsqlGetAllFieldsOfTableInFromOrOuterCTE("", "", s.Qualifier)
}

// -------------------- Expression resolver --------------------

// resolveExpression maps an ExprNode to a QuerySpanResult.
func (q *querySpanExtractor) resolveExpression(expr ast.ExprNode) (QuerySpanResult, error) {
	return q.resolveExpressionNode(expr)
}

func (q *querySpanExtractor) resolveExpressionNode(n ast.Node) (QuerySpanResult, error) {
	if n == nil {
		return QuerySpanResult{SourceColumns: make(ColumnSet)}, nil
	}
	if r, handled, err := q.resolveSimpleExpressionNode(n); handled {
		return r, err
	}
	switch v := n.(type) {
	case *ast.ColumnRef:
		if r, handled, err := q.resolvePseudoColumnRef(v); handled {
			if err != nil {
				return QuerySpanResult{}, fmt.Errorf("failed to resolve column %q: %w", v.Column, err)
			}
			return r, nil
		}
		r, err := q.tsqlIsFieldSensitive(v.Database, v.Schema, v.Table, v.Column)
		if err != nil {
			return QuerySpanResult{}, fmt.Errorf("failed to resolve column %q: %w", v.Column, err)
		}
		r.IsPlainField = true
		return r, nil
	case *ast.AtTimeZoneExpr:
		r, err := q.mergeExprSources(v.Expr, v.TimeZone)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.BinaryExpr:
		r, err := q.mergeExprSources(v.Left, v.Right)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.UnaryExpr:
		r, err := q.resolveExpressionNode(v.Operand)
		if err != nil {
			return r, err
		}
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, nil
	case *ast.BetweenExpr:
		r, err := q.mergeExprSources(v.Expr, v.Low, v.High)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.LikeExpr:
		r, err := q.mergeExprSources(v.Expr, v.Pattern, v.Escape)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.IsExpr:
		r, err := q.mergeExprSources(v.Expr, v.Value)
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, err
	case *ast.InExpr:
		exprs := []ast.ExprNode{v.Expr}
		if v.List != nil {
			for _, it := range v.List.Items {
				if e, ok := it.(ast.ExprNode); ok {
					exprs = append(exprs, e)
				}
			}
		}
		if v.Subquery != nil {
			exprs = append(exprs, v.Subquery)
		}
		r, err := q.mergeExprSources(exprs...)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.CaseExpr:
		exprs := []ast.ExprNode{v.Arg, v.ElseExpr}
		if v.WhenList != nil {
			for _, it := range v.WhenList.Items {
				if cw, ok := it.(*ast.CaseWhen); ok {
					exprs = append(exprs, cw.Condition, cw.Result)
				}
			}
		}
		r, err := q.mergeExprSources(exprs...)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.IifExpr:
		r, err := q.mergeExprSources(v.Condition, v.TrueVal, v.FalseVal)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.CoalesceExpr:
		r, err := q.mergeListSources(v.Args)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.NullifExpr:
		r, err := q.mergeExprSources(v.Left, v.Right)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.FuncCallExpr:
		r, err := q.mergeListSources(v.Args)
		if err != nil {
			return r, err
		}
		if err := q.mergeOverSources(&r, v.Over); err != nil {
			return r, err
		}
		if v.Within != nil {
			sub, err := q.mergeOrderBySources(v.Within)
			if err != nil {
				return r, err
			}
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
		}
		r.Name = q.sliceName(n)
		return r, nil
	case *ast.MethodCallExpr:
		r, err := q.mergeListSources(v.Args)
		if err != nil {
			return r, err
		}
		// An instance method reads its receiver: CAST(shape AS
		// geometry).STArea() reads shape.
		if v.Expr != nil {
			recv, err := q.resolveExpressionNode(v.Expr)
			if err != nil {
				return r, err
			}
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, recv.SourceColumns)
		}
		r.Name = q.sliceName(n)
		return r, nil
	case *ast.CastExpr:
		r, err := q.resolveExpressionNode(v.Expr)
		if err != nil {
			return r, err
		}
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, nil
	case *ast.ConvertExpr:
		r, err := q.mergeExprSources(v.Expr, v.Style)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.TryCastExpr:
		r, err := q.resolveExpressionNode(v.Expr)
		if err != nil {
			return r, err
		}
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, nil
	case *ast.TryConvertExpr:
		r, err := q.mergeExprSources(v.Expr, v.Style)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.OdbcFunctionCallExpr:
		var args []ast.ExprNode
		if v.Args != nil {
			for _, it := range v.Args.Items {
				if arg, ok := it.(ast.ExprNode); ok {
					args = append(args, arg)
				}
			}
		}
		r, err := q.mergeExprSources(args...)
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, err
	case *ast.OdbcLiteralExpr:
		return QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}, nil
	case *ast.SubqueryExpr:
		return q.resolveSubqueryExpr(v)
	case *ast.SubqueryComparisonExpr:
		r, err := q.resolveExpressionNode(v.Left)
		if err != nil {
			return r, err
		}
		if body, ok := v.Subquery.(*ast.SelectStmt); ok {
			clone := q.cloneForSubquery()
			sub, subErr := clone.extractFromSelectStmt(body)
			if subErr != nil {
				return r, subErr
			}
			for _, sc := range sub.GetQuerySpanResult() {
				r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sc.SourceColumns)
			}
			q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
		}
		r.Name = q.sliceName(n)
		r.IsPlainField = false
		return r, nil
	case *ast.ExistsExpr:
		r := QuerySpanResult{
			Name:          q.sliceName(n),
			SourceColumns: make(ColumnSet),
		}
		if body, ok := v.Subquery.(*ast.SelectStmt); ok {
			clone := q.cloneForSubquery()
			if _, err := clone.extractFromSelectStmt(body); err != nil {
				return r, err
			}
			q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
		}
		return r, nil
	case *ast.FullTextPredicate:
		r := QuerySpanResult{SourceColumns: make(ColumnSet)}
		if v.Columns != nil {
			for _, it := range v.Columns.Items {
				if cr, ok := it.(*ast.ColumnRef); ok {
					sub, err := q.tsqlIsFieldSensitive(cr.Database, cr.Schema, cr.Table, cr.Column)
					if err != nil {
						return r, fmt.Errorf("failed to resolve full-text predicate column %q: %w", cr.Column, err)
					}
					r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
				}
			}
		}
		if v.Value != nil {
			sub, err := q.resolveExpressionNode(v.Value)
			if err != nil {
				return r, err
			}
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
		}
		r.Name = q.sliceName(n)
		return r, nil
	case *ast.GraphMatchPredicate:
		r := QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}
		if err := q.walkGraphMatchRefs(v.Expression, func(name string) error {
			sub, err := q.resolveGraphMatchRef(name)
			if err != nil {
				return err
			}
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
			return nil
		}); err != nil {
			return r, err
		}
		r.IsPlainField = false
		return r, nil
	case *ast.StarExpr:
		// Star in scalar context: treat as empty (already handled at target-list level).
		return QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}, nil
	case *ast.GroupingSetsExpr, *ast.RollupExpr, *ast.CubeExpr:
		// Appear in GROUP BY; not typical targets. Fallback to empty.
		return QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}, nil
	case *ast.ResTarget:
		r, err := q.resolveExpressionNode(v.Val)
		if err != nil {
			return r, err
		}
		if v.Name != "" {
			r.Name = v.Name
		}
		return r, nil
	case *ast.SelectAssign:
		return q.resolveExpressionNode(v.Value)
	case *ast.CaseWhen:
		return q.mergeExprSources(v.Condition, v.Result)
	case *ast.ParseExpr:
		r, err := q.mergeExprSources(v.Expr, v.Culture)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.JsonKeyValueExpr:
		r, err := q.mergeExprSources(v.Key, v.Value)
		r.Name = q.sliceName(n)
		return r, err
	case *ast.NextValueForExpr:
		r := QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}
		err := q.mergeOverSources(&r, v.Over)
		return r, err
	case *ast.CurrentOfExpr, *ast.DatePart:
		return QuerySpanResult{Name: q.sliceName(n), SourceColumns: make(ColumnSet)}, nil
	default:
		// An expression whose lineage is not modeled would otherwise lose the
		// columns it reads, and with them their masking.
		return QuerySpanResult{}, unsupportedNodeError("unsupported expression", n)
	}
}

// mergeOverSources adds the columns a window's PARTITION BY and ORDER BY
// read to r.
func (q *querySpanExtractor) mergeOverSources(r *QuerySpanResult, over *ast.OverClause) error {
	if over == nil {
		return nil
	}
	sub, err := q.mergeListSources(over.PartitionBy)
	if err != nil {
		return err
	}
	r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
	sub, err = q.mergeOrderBySources(over.OrderBy)
	if err != nil {
		return err
	}
	r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
	return nil
}

func (q *querySpanExtractor) resolveSimpleExpressionNode(n ast.Node) (QuerySpanResult, bool, error) {
	switch v := n.(type) {
	case *ast.Literal, *ast.VariableRef, *ast.Boolean:
		return QuerySpanResult{Name: q.sliceName(v), SourceColumns: make(ColumnSet)}, true, nil
	case *ast.ParenExpr:
		r, err := q.resolveExpressionNode(v.Expr)
		return r, true, err
	case *ast.CollateExpr:
		r, err := q.resolveExpressionNode(v.Expr)
		return r, true, err
	default:
		return QuerySpanResult{}, false, nil
	}
}

func (q *querySpanExtractor) resolveSubqueryExpr(sq *ast.SubqueryExpr) (QuerySpanResult, error) {
	r := QuerySpanResult{
		Name:          q.sliceName(sq),
		SourceColumns: make(ColumnSet),
	}
	body, ok := sq.Query.(*ast.SelectStmt)
	if !ok {
		return r, nil
	}
	clone := q.cloneForSubquery()
	pseudo, err := clone.extractFromSelectStmt(body)
	if err != nil {
		return r, err
	}
	// Scalar subquery contributes its first result column's source columns.
	results := pseudo.GetQuerySpanResult()
	for _, c := range results {
		r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, c.SourceColumns)
	}
	q.predicateColumns, _ = mergeColumnSet(q.predicateColumns, clone.predicateColumns)
	return r, nil
}

func (q *querySpanExtractor) mergeExprSources(exprs ...ast.ExprNode) (QuerySpanResult, error) {
	r := QuerySpanResult{SourceColumns: make(ColumnSet)}
	for _, e := range exprs {
		if e == nil {
			continue
		}
		sub, err := q.resolveExpressionNode(e)
		if err != nil {
			return r, err
		}
		r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
	}
	return r, nil
}

func (q *querySpanExtractor) mergeListSources(list *ast.List) (QuerySpanResult, error) {
	r := QuerySpanResult{SourceColumns: make(ColumnSet)}
	if list == nil {
		return r, nil
	}
	for _, it := range list.Items {
		if e, ok := it.(ast.ExprNode); ok {
			sub, err := q.resolveExpressionNode(e)
			if err != nil {
				return r, err
			}
			r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
		}
	}
	return r, nil
}

func (q *querySpanExtractor) mergeOrderBySources(list *ast.List) (QuerySpanResult, error) {
	r := QuerySpanResult{SourceColumns: make(ColumnSet)}
	if list == nil {
		return r, nil
	}
	for _, it := range list.Items {
		var expr ast.Node
		switch v := it.(type) {
		case *ast.OrderByItem:
			expr = v.Expr
		case ast.ExprNode:
			expr = v
		default:
			return r, unsupportedNodeError("unsupported ORDER BY item", it)
		}
		sub, err := q.resolveExpressionNode(expr)
		if err != nil {
			return r, err
		}
		r.SourceColumns, _ = mergeColumnSet(r.SourceColumns, sub.SourceColumns)
	}
	return r, nil
}

// -------------------- Helpers --------------------

// sliceName returns the whitespace-stripped source text covering a node's Loc.
// For *SubqueryExpr (and similar paren-wrapped nodes) the Loc spans the
// surrounding parens; a matching leading `(` / trailing `)` pair is stripped.
func (q *querySpanExtractor) sliceName(n ast.Node) string {
	loc := nodeLoc(n)
	if loc.Start < 0 || loc.End < 0 || loc.Start >= len(q.source) {
		return ""
	}
	end := loc.End
	if end > len(q.source) {
		end = len(q.source)
	}
	text := stripSpaces(q.source[loc.Start:end])
	switch n.(type) {
	case *ast.SubqueryExpr, *ast.ExistsExpr, *ast.SubqueryComparisonExpr:
		if len(text) >= 2 && text[0] == '(' && text[len(text)-1] == ')' {
			text = text[1 : len(text)-1]
		}
	default:
	}
	return text
}

func stripSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// nodeLoc extracts the Loc field from an ast.Node via type switching.
// Not every node has a Loc (the Node interface doesn't expose it), so we
// enumerate the types we care about. Returns zero Loc otherwise.
func nodeLoc(n ast.Node) ast.Loc {
	if n == nil {
		return ast.NoLoc()
	}
	if loc, ok := expressionNodeLoc(n); ok {
		return loc
	}
	switch v := n.(type) {
	case *ast.SelectStmt:
		return v.Loc
	case *ast.TableRef:
		return v.Loc
	case *ast.ResTarget:
		return v.Loc
	case *ast.SelectAssign:
		return v.Loc
	default:
		return ast.NoLoc()
	}
}

func expressionNodeLoc(n ast.Node) (ast.Loc, bool) {
	switch v := n.(type) {
	case *ast.ColumnRef:
		return v.Loc, true
	case *ast.Literal:
		return v.Loc, true
	case *ast.VariableRef:
		return v.Loc, true
	case *ast.StarExpr:
		return v.Loc, true
	case *ast.BinaryExpr:
		return v.Loc, true
	case *ast.UnaryExpr:
		return v.Loc, true
	case *ast.BetweenExpr:
		return v.Loc, true
	case *ast.LikeExpr:
		return v.Loc, true
	case *ast.IsExpr:
		return v.Loc, true
	case *ast.InExpr:
		return v.Loc, true
	case *ast.CaseExpr:
		return v.Loc, true
	case *ast.CaseWhen:
		return v.Loc, true
	case *ast.IifExpr:
		return v.Loc, true
	case *ast.CoalesceExpr:
		return v.Loc, true
	case *ast.NullifExpr:
		return v.Loc, true
	case *ast.FuncCallExpr:
		return v.Loc, true
	case *ast.MethodCallExpr:
		return v.Loc, true
	case *ast.CastExpr:
		return v.Loc, true
	case *ast.ConvertExpr:
		return v.Loc, true
	case *ast.TryCastExpr:
		return v.Loc, true
	case *ast.TryConvertExpr:
		return v.Loc, true
	case *ast.SubqueryExpr:
		return v.Loc, true
	case *ast.SubqueryComparisonExpr:
		return v.Loc, true
	case *ast.ExistsExpr:
		return v.Loc, true
	case *ast.FullTextPredicate:
		return v.Loc, true
	case *ast.OdbcFunctionCallExpr:
		return v.Loc, true
	case *ast.OdbcLiteralExpr:
		return v.Loc, true
	case *ast.GraphMatchPredicate:
		return v.Loc, true
	case *ast.GraphMatchNodeExpr:
		return v.Loc, true
	case *ast.GraphMatchCompositeExpr:
		return v.Loc, true
	case *ast.GraphMatchLastNodePredicate:
		return v.Loc, true
	case *ast.GraphMatchRecursivePredicate:
		return v.Loc, true
	case *ast.ParenExpr:
		return v.Loc, true
	case *ast.CollateExpr:
		return v.Loc, true
	case *ast.AtTimeZoneExpr:
		return v.Loc, true
	case *ast.ParseExpr:
		return v.Loc, true
	case *ast.JsonKeyValueExpr:
		return v.Loc, true
	case *ast.NextValueForExpr:
		return v.Loc, true
	case *ast.DatePart:
		return v.Loc, true
	case *ast.CurrentOfExpr:
		return v.Loc, true
	default:
		return ast.NoLoc(), false
	}
}

// -------------------- Access tables --------------------

// collectAccessTables records every table the statement reads or writes in a
// table-source position, including those inside CTE bodies and subqueries. A
// reference to a CTE in scope names no physical table and is skipped.
func (q *querySpanExtractor) collectAccessTables(root ast.Node) ColumnSet {
	w := &accessTableWalker{q: q, lifted: map[*ast.WithClause]bool{}, out: make(ColumnSet)}
	w.walk(root)
	return w.out
}

type accessTableWalker struct {
	q *querySpanExtractor
	// ctes holds the CTE names in scope.
	ctes []string
	// lifted holds the WITH clauses already in scope. The parser attaches a
	// set operation's WITH clause to its leftmost arm, so the clause is
	// brought into scope at the set operation and skipped at the arm.
	lifted map[*ast.WithClause]bool
	out    ColumnSet
}

func (w *accessTableWalker) walk(n ast.Node) {
	switch v := n.(type) {
	case nil:
	case *ast.SelectStmt:
		defer w.restoreScope(len(w.ctes))
		w.with(leadingWithClause(v))
		w.walkList(v.FromClause)
		if v.Larg != nil {
			w.walk(v.Larg)
		}
		if v.Rarg != nil {
			w.walk(v.Rarg)
		}
		w.subqueries(v.Top, v.TargetList, v.WhereClause, v.GroupByClause, v.HavingClause, v.WindowClause, v.OrderByClause, v.OffsetClause, v.FetchClause)
	case *ast.InsertStmt:
		defer w.restoreScope(len(w.ctes))
		w.with(v.WithClause)
		w.walk(v.Relation)
		w.walk(v.Source)
		w.subqueries(v.Top)
	case *ast.UpdateStmt:
		defer w.restoreScope(len(w.ctes))
		w.with(v.WithClause)
		w.walk(v.Relation)
		w.walkList(v.FromClause)
		w.subqueries(v.Top, v.SetClause, v.WhereClause)
	case *ast.DeleteStmt:
		defer w.restoreScope(len(w.ctes))
		w.with(v.WithClause)
		w.walk(v.Relation)
		w.walkList(v.FromClause)
		w.subqueries(v.Top, v.WhereClause)
	case *ast.MergeStmt:
		defer w.restoreScope(len(w.ctes))
		w.with(v.WithClause)
		w.walk(v.Target)
		w.walk(v.Source)
		w.subqueries(v.Top, v.OnCondition, v.WhenClauses)
	case *ast.TableRef:
		if !w.isCTE(v) {
			addAccessTable(v, w.q.defaultDatabase, w.q.defaultSchema, w.out)
		}
	case *ast.AliasedTableRef:
		w.walk(v.Table)
	case *ast.JoinClause:
		w.walk(v.Left)
		w.walk(v.Right)
		w.subqueries(v.Condition)
	case *ast.PivotExpr:
		w.walk(v.Source)
	case *ast.UnpivotExpr:
		w.walk(v.Source)
	case *ast.ValuesClause:
		w.subqueries(v)
	case *ast.TableVarRef, *ast.TableVarMethodCallRef:
		// Table variables are not physical tables.
	case ast.ExprNode:
		// A derived table, a table-valued function, or an expression.
		w.subqueries(v)
	default:
	}
}

func (w *accessTableWalker) walkList(list *ast.List) {
	if list == nil {
		return
	}
	for _, it := range list.Items {
		w.walk(it)
	}
}

// subqueries walks the queries nested anywhere in nodes, which are not table
// sources themselves.
func (w *accessTableWalker) subqueries(nodes ...ast.Node) {
	for _, n := range nodes {
		ast.Inspect(n, func(c ast.Node) bool {
			if sel, ok := c.(*ast.SelectStmt); ok {
				w.walk(sel)
				return false
			}
			return true
		})
	}
}

// with brings a WITH clause's CTEs into scope. Each CTE's body sees itself,
// for recursion, and the CTEs before it.
func (w *accessTableWalker) with(clause *ast.WithClause) {
	if clause == nil || w.lifted[clause] || clause.CTEs == nil {
		return
	}
	w.lifted[clause] = true
	for _, it := range clause.CTEs.Items {
		cte, ok := it.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		w.ctes = append(w.ctes, cte.Name)
		w.walk(cte.Query)
	}
}

func (w *accessTableWalker) restoreScope(n int) {
	w.ctes = w.ctes[:n]
}

// isCTE reports whether t names a CTE in scope. As in name resolution, a CTE
// shadows a table of the same name only when no schema or database is given.
func (w *accessTableWalker) isCTE(t *ast.TableRef) bool {
	if t.Server != "" || t.Database != "" || t.Schema != "" {
		return false
	}
	return slices.ContainsFunc(w.ctes, func(name string) bool {
		return w.q.isIdentifierEqual(t.Object, name)
	})
}

// leadingWithClause returns the WITH clause of sel, which the parser attaches
// to the leftmost arm when sel is a set operation.
func leadingWithClause(sel *ast.SelectStmt) *ast.WithClause {
	for sel != nil {
		if sel.WithClause != nil || sel.Op == ast.SetOpNone {
			return sel.WithClause
		}
		sel = sel.Larg
	}
	return nil
}

func addAccessTable(t *ast.TableRef, defDB, defSchema string, out ColumnSet) {
	if t == nil || strings.HasPrefix(t.Object, "#") {
		return
	}
	database := defDB
	if t.Database != "" {
		database = t.Database
	}
	schema := defSchema
	if t.Schema != "" {
		schema = t.Schema
	}
	out[ColumnResource{
		Server:   t.Server,
		Database: database,
		Schema:   schema,
		Table:    t.Object,
	}] = true
}
