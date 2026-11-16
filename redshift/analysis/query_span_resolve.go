package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/redshift"
	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/catalog"
)

const redshiftGetDatabaseMetadataError = "failed to get database metadata for database: %s"

type querySpanResolver struct {
	defaultDatabase string
	searchPath      []string
	opts            ResolveOptions
}

func newQuerySpanResolver(defaultDatabase string, searchPath []string, opts ResolveOptions) *querySpanResolver {
	return &querySpanResolver{
		defaultDatabase: defaultDatabase,
		searchPath:      normalizeSearchPath(searchPath),
		opts:            opts,
	}
}

func normalizeSearchPath(searchPath []string) []string {
	result := make([]string, 0, len(searchPath))
	for _, schema := range searchPath {
		if schema != "" {
			result = append(result, schema)
		}
	}
	if len(result) == 0 {
		return []string{"public"}
	}
	return result
}

func (q *querySpanResolver) getQuerySpan(statement string) (*ResolvedQuerySpan, error) {
	return q.getQuerySpanWithFunctionStack(statement, map[string]bool{})
}

func (q *querySpanResolver) getQuerySpanWithFunctionStack(statement string, functionStack map[string]bool) (*ResolvedQuerySpan, error) {
	stmts, err := parseRedshiftSpan(statement)
	if err != nil {
		return nil, err
	}
	if len(stmts) == 0 || stmts[0].Empty() {
		return emptyQuerySpan(Select), nil
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}

	node := stmts[0].AST
	queryType := querySpanType(node, false /* allSystems */)
	if explain, ok := node.(*ast.ExplainStmt); ok && hasExplainAnalyze(explain) {
		return q.getExplainAnalyzeQuerySpan(explain.Query, queryType)
	}
	if unload, ok := node.(*ast.UnloadStmt); ok {
		return q.getUnloadQuerySpan(unload, functionStack)
	}
	if _, ok := node.(*ast.VariableSetStmt); ok {
		return emptyQuerySpan(queryType), nil
	}
	if queryType != Select {
		return emptyQuerySpan(queryType), nil
	}

	accessTables, earlySpan, err := q.collectSelectAccessTables(node)
	if err != nil {
		return nil, err
	}
	if earlySpan != nil {
		return earlySpan, nil
	}

	cat, err := q.newQuerySpanCatalog()
	if err != nil {
		return nil, err
	}
	selectStmt, _ := node.(*ast.SelectStmt)
	span, err := getQuerySpanFromSelect(cat, statement, selectStmt)
	if err != nil {
		if isQuerySpanNotFound(err) {
			return q.notFoundQuerySpan(err), nil
		}
		return nil, err
	}

	results := q.convertResults(span.Results)
	functionAccessTables, functionResultSources, functionIssues, err := q.collectFunctionBodySources(node, functionStack, make(map[string]functionBodySources))
	if err != nil {
		return nil, err
	}
	for source := range functionAccessTables {
		accessTables[source] = true
	}
	if functionIssues.functionNotSupportedError != nil && len(accessTables) == 0 {
		accessTables[ColumnResource{Database: q.defaultDatabase}] = true
	}
	mergeFunctionResultSources(results, functionResultSources)
	if hasRangeVar(node) && hasResultWithoutSource(results) && !q.allAccessTablesExist(accessTables) {
		return q.notFoundQuerySpan(&ResourceNotFoundError{Database: &q.defaultDatabase}), nil
	}

	return &ResolvedQuerySpan{
		Type:                      Select,
		SourceColumns:             accessTables,
		PredicateColumns:          q.convertColumnList(span.PredicateColumns),
		Results:                   results,
		NotFoundError:             functionIssues.notFoundError,
		FunctionNotSupportedError: functionIssues.functionNotSupportedError,
	}, nil
}

func (q *querySpanResolver) getExplainAnalyzeQuerySpan(query ast.Node, queryType QueryType) (*ResolvedQuerySpan, error) {
	if queryType != Select {
		return emptyQuerySpan(queryType), nil
	}

	accessTables, earlySpan, err := q.collectSelectAccessTables(query)
	if err != nil {
		return nil, err
	}
	if earlySpan != nil {
		return earlySpan, nil
	}
	return &ResolvedQuerySpan{
		Type:             Select,
		SourceColumns:    accessTables,
		PredicateColumns: ColumnSet{},
		Results:          []ResolvedQuerySpanResult{},
	}, nil
}

func (q *querySpanResolver) getUnloadQuerySpan(unload *ast.UnloadStmt, functionStack map[string]bool) (*ResolvedQuerySpan, error) {
	span, err := q.getQuerySpanWithFunctionStack(unload.Query, functionStack)
	if err != nil {
		return nil, err
	}
	return &ResolvedQuerySpan{
		Type:                      DML,
		SourceColumns:             span.SourceColumns,
		PredicateColumns:          span.PredicateColumns,
		Results:                   []ResolvedQuerySpanResult{},
		NotFoundError:             span.NotFoundError,
		FunctionNotSupportedError: span.FunctionNotSupportedError,
	}, nil
}

func (q *querySpanResolver) collectSelectAccessTables(node ast.Node) (ColumnSet, *ResolvedQuerySpan, error) {
	accessTables, err := q.collectAccessTables(node)
	if err != nil {
		var resourceNotFound *ResourceNotFoundError
		if errors.As(err, &resourceNotFound) {
			return nil, q.notFoundQuerySpan(resourceNotFound), nil
		}
		return nil, nil, err
	}
	allSystems, mixed := isMixedQuery(accessTables)
	if mixed {
		return nil, nil, ErrMixUserSystemTables
	}
	if allSystems {
		return nil, emptyQuerySpan(SelectInfoSchema), nil
	}
	return accessTables, nil, nil
}

func (q *querySpanResolver) notFoundQuerySpan(err error) *ResolvedQuerySpan {
	return &ResolvedQuerySpan{
		Type: Select,
		SourceColumns: ColumnSet{
			{Database: q.defaultDatabase}: true,
		},
		PredicateColumns: ColumnSet{},
		Results:          []ResolvedQuerySpanResult{},
		NotFoundError:    err,
	}
}

func emptyQuerySpan(queryType QueryType) *ResolvedQuerySpan {
	return &ResolvedQuerySpan{
		Type:          queryType,
		SourceColumns: ColumnSet{},
		Results:       []ResolvedQuerySpanResult{},
	}
}

func (q *querySpanResolver) newQuerySpanCatalog() (*catalog.Catalog, error) {
	if q.opts.GetDatabaseMetadata == nil {
		return nil, errors.New("GetDatabaseMetadataFunc is not set in GetQuerySpanContext")
	}
	dbMeta, err := q.getDatabaseMetadata(q.defaultDatabase)
	if err != nil {
		return nil, err
	}
	if dbMeta == nil {
		return nil, &ResourceNotFoundError{Database: &q.defaultDatabase}
	}
	cat := catalog.New()
	cat.SetSearchPath(q.searchPath)
	cat.UseMetadata(dbMeta)
	return cat, nil
}

type functionBodySources struct {
	accessTables ColumnSet
	results      []ResolvedQuerySpanResult
	issues       functionBodyIssues
}

type functionBodyIssues struct {
	notFoundError             error
	functionNotSupportedError error
}

func (i *functionBodyIssues) merge(other functionBodyIssues) {
	if i.notFoundError == nil {
		i.notFoundError = other.notFoundError
	}
	if i.functionNotSupportedError == nil {
		i.functionNotSupportedError = other.functionNotSupportedError
	}
}

type functionSourceProjection struct {
	names   []string
	sources []ColumnSet
	// opaque marks a projection with an unresolved star, whose output width is
	// unknown; expanding a star over it would misalign the following columns.
	opaque bool
}

// functionResultSource is one output column's name and SQL-function sources.
type functionResultSource struct {
	Name    string
	Sources ColumnSet
}

// functionProjection is a FROM item's output, kept as a concatenation tree for
// ordinary joins so a left-deep chain does not copy its accumulated left side
// at every node. It is flattened once for a star or a coalescing join.
type functionProjection struct {
	flat        *functionSourceProjection
	left, right *functionProjection
}

func flatFunctionProjection(p functionSourceProjection) *functionProjection {
	return &functionProjection{flat: &p}
}

func (p *functionProjection) flatten() (functionSourceProjection, bool) {
	if p == nil {
		return functionSourceProjection{}, false
	}
	if p.flat != nil {
		return *p.flat, true
	}
	var result functionSourceProjection
	var walk func(*functionProjection) bool
	walk = func(n *functionProjection) bool {
		if n == nil {
			return false
		}
		if n.flat != nil {
			result.names = append(result.names, n.flat.names...)
			result.sources = append(result.sources, n.flat.sources...)
			result.opaque = result.opaque || n.flat.opaque
			return true
		}
		return walk(n.left) && walk(n.right)
	}
	if !walk(p) {
		return functionSourceProjection{}, false
	}
	p.flat = &result
	p.left, p.right = nil, nil
	return result, true
}

// functionProjectionSource computes a query's whole FROM projection on first
// use, so a query with neither a star nor a function does not build one.
type functionProjectionSource struct {
	done       bool
	projection functionSourceProjection
	ok         bool
	compute    func() (functionSourceProjection, bool)
}

func (s *functionProjectionSource) get() (functionSourceProjection, bool) {
	if s == nil {
		return functionSourceProjection{}, false
	}
	if !s.done {
		s.projection, s.ok = s.compute()
		s.done = true
	}
	return s.projection, s.ok
}

// functionSourceScope separates the CTE definitions a FROM item may reference
// (defs) from the relations actually bound in the current FROM clause
// (bindings). Keeping them apart lets `FROM cte AS x` expose only x without also
// leaving the CTE name visible, which would make an unqualified column look
// ambiguous.
type functionSourceScope struct {
	defs     map[string]functionSourceProjection
	bindings map[string]functionSourceProjection
	// projection is this query's whole FROM output, including a USING/NATURAL
	// join's coalesced columns, built lazily when an unqualified column needs
	// it.
	projection *functionProjectionSource
	// outerLevels are the enclosing queries' scopes, nearest last. Each level is
	// searched in turn so a correlated reference uses the nearest relation that
	// exposes the column, and ambiguity is decided within one level only.
	outerLevels []functionScopeLevel
}

type functionScopeLevel struct {
	bindings   map[string]functionSourceProjection
	projection *functionProjectionSource
}

func (q *querySpanResolver) collectAccessTables(node ast.Node) (ColumnSet, error) {
	result := make(ColumnSet)
	if err := q.collectAccessTablesFromNode(node, map[string]bool{}, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (q *querySpanResolver) collectAccessTablesFromNode(node ast.Node, ctes map[string]bool, result ColumnSet) error {
	if node == nil {
		return nil
	}
	if selectStmt, ok := node.(*ast.SelectStmt); ok {
		return q.collectAccessTablesFromSelect(selectStmt, ctes, result)
	}

	var err error
	ast.Inspect(node, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		if n != node {
			if selectStmt, ok := n.(*ast.SelectStmt); ok {
				err = q.collectAccessTablesFromSelect(selectStmt, ctes, result)
				return false
			}
		}
		rangeVar, ok := n.(*ast.RangeVar)
		if !ok || rangeVar == nil || rangeVar.Relname == "" {
			return true
		}
		if isCTEReference(rangeVar, ctes) {
			return true
		}

		resource, resolveErr := q.resolveRangeVar(rangeVar)
		if resolveErr != nil {
			err = resolveErr
			return false
		}
		if resource.Table != "" {
			result[resource] = true
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}

func (q *querySpanResolver) collectAccessTablesFromSelect(selectStmt *ast.SelectStmt, ctes map[string]bool, result ColumnSet) error {
	localCTEs := cloneCTENameSet(ctes)
	if selectStmt.WithClause != nil && selectStmt.WithClause.Ctes != nil {
		cteScope := cloneCTENameSet(ctes)
		if selectStmt.WithClause.Recursive {
			for _, item := range selectStmt.WithClause.Ctes.Items {
				cte, ok := item.(*ast.CommonTableExpr)
				if ok && cte.Ctename != "" {
					cteScope[strings.ToLower(cte.Ctename)] = true
				}
			}
		}
		for _, item := range selectStmt.WithClause.Ctes.Items {
			cte, ok := item.(*ast.CommonTableExpr)
			if !ok || cte.Ctename == "" {
				continue
			}
			if err := q.collectAccessTablesFromNode(cte.Ctequery, cteScope, result); err != nil {
				return err
			}
			name := strings.ToLower(cte.Ctename)
			cteScope[name] = true
			localCTEs[name] = true
		}
	}

	for _, node := range []ast.Node{
		selectStmt.WhereClause,
		selectStmt.HavingClause,
		selectStmt.QualifyClause,
		selectStmt.LimitOffset,
		selectStmt.LimitCount,
	} {
		if err := q.collectAccessTablesFromNode(node, localCTEs, result); err != nil {
			return err
		}
	}
	if selectStmt.IntoClause != nil {
		if err := q.collectAccessTablesFromNode(selectStmt.IntoClause, localCTEs, result); err != nil {
			return err
		}
	}
	if selectStmt.Larg != nil {
		if err := q.collectAccessTablesFromNode(selectStmt.Larg, localCTEs, result); err != nil {
			return err
		}
	}
	if selectStmt.Rarg != nil {
		if err := q.collectAccessTablesFromNode(selectStmt.Rarg, localCTEs, result); err != nil {
			return err
		}
	}
	for _, list := range []*ast.List{
		selectStmt.DistinctClause,
		selectStmt.TargetList,
		selectStmt.FromClause,
		selectStmt.GroupClause,
		selectStmt.WindowClause,
		selectStmt.ExcludeList,
		selectStmt.ValuesLists,
		selectStmt.SortClause,
		selectStmt.LockingClause,
	} {
		if err := q.collectAccessTablesFromList(list, localCTEs, result); err != nil {
			return err
		}
	}
	return nil
}

func (q *querySpanResolver) collectAccessTablesFromList(list *ast.List, ctes map[string]bool, result ColumnSet) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		if err := q.collectAccessTablesFromNode(item, ctes, result); err != nil {
			return err
		}
	}
	return nil
}

func (q *querySpanResolver) collectFunctionBodySources(node ast.Node, functionStack map[string]bool, cache map[string]functionBodySources) (ColumnSet, []functionResultSource, functionBodyIssues, error) {
	accessTables, issues, err := q.collectFunctionBodyAccessTables(node, functionStack, cache)
	if err != nil {
		return nil, nil, functionBodyIssues{}, err
	}

	var resultSources []functionResultSource
	if selectStmt, ok := node.(*ast.SelectStmt); ok {
		projection, sourceIssues, err := q.collectSelectFunctionProjection(selectStmt, functionStack, cache, functionSourceScope{})
		if err != nil {
			return nil, nil, functionBodyIssues{}, err
		}
		issues.merge(sourceIssues)
		resultSources = functionSourcesFromProjection(projection)
	}
	return accessTables, resultSources, issues, nil
}

// functionSourcesFromProjection pairs each source set with its output name.
func functionSourcesFromProjection(projection functionSourceProjection) []functionResultSource {
	result := make([]functionResultSource, len(projection.sources))
	for i, sources := range projection.sources {
		name := ""
		if i < len(projection.names) {
			name = projection.names[i]
		}
		result[i] = functionResultSource{Name: name, Sources: sources}
	}
	return result
}

func (q *querySpanResolver) collectFunctionBodyAccessTables(node ast.Node, functionStack map[string]bool, cache map[string]functionBodySources) (ColumnSet, functionBodyIssues, error) {
	result := make(ColumnSet)
	if node == nil {
		return result, functionBodyIssues{}, nil
	}

	var err error
	issues := functionBodyIssues{}
	ast.Inspect(node, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		call, ok := n.(*ast.FuncCall)
		if !ok {
			return true
		}
		sources, sourceErr := q.functionBodySourcesForCall(call, functionStack, cache)
		if sourceErr != nil {
			err = sourceErr
			return false
		}
		issues.merge(sources.issues)
		for source := range sources.accessTables {
			result[source] = true
		}
		return true
	})
	return result, issues, err
}

// collectSelectFunctionProjection returns a select's output names and their
// SQL-function sources together. A star expands both at once, so a following
// function's name and sources stay aligned instead of indexing past the star's
// single select item.
func (q *querySpanResolver) collectSelectFunctionProjection(selectStmt *ast.SelectStmt, functionStack map[string]bool, cache map[string]functionBodySources, outerScope functionSourceScope) (functionSourceProjection, functionBodyIssues, error) {
	if selectStmt == nil {
		return functionSourceProjection{}, functionBodyIssues{}, nil
	}

	if selectStmt.Op != ast.SETOP_NONE {
		// A set operation's output column i is fed by column i of every arm, so
		// its SQL-function lineage is the union across arms; names come from the
		// left arm, as in SQL. The upper statement may define CTEs the arms
		// reference.
		scope, issues, err := q.functionSourceScopeForSelect(selectStmt, functionStack, cache, outerScope)
		if err != nil {
			return functionSourceProjection{}, functionBodyIssues{}, err
		}
		left, leftIssues, err := q.collectSelectFunctionProjection(selectStmt.Larg, functionStack, cache, scope)
		if err != nil {
			return functionSourceProjection{}, functionBodyIssues{}, err
		}
		issues.merge(leftIssues)
		right, rightIssues, err := q.collectSelectFunctionProjection(selectStmt.Rarg, functionStack, cache, scope)
		if err != nil {
			return functionSourceProjection{}, functionBodyIssues{}, err
		}
		issues.merge(rightIssues)
		if left.opaque || right.opaque {
			// An unresolved star makes an arm's width unknown; positional
			// merging would pair the wrong columns. Keep the left arm's shape as
			// the output shape and leave the projection opaque.
			result := left
			result.opaque = true
			return result, issues, nil
		}
		return functionSourceProjection{
			names:   left.names,
			sources: mergeFunctionArmResultSources(left.sources, right.sources),
		}, issues, nil
	}

	if selectStmt.TargetList == nil {
		return functionSourceProjection{}, functionBodyIssues{}, nil
	}

	scope, issues, err := q.functionSourceScopeForSelect(selectStmt, functionStack, cache, outerScope)
	if err != nil {
		return functionSourceProjection{}, functionBodyIssues{}, err
	}

	result := functionSourceProjection{
		names:   make([]string, 0, selectStmt.TargetList.Len()),
		sources: make([]ColumnSet, 0, selectStmt.TargetList.Len()),
	}
	for _, item := range selectStmt.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok {
			result.names = append(result.names, "")
			result.sources = append(result.sources, nil)
			continue
		}
		if star, ok := q.functionStarProjection(target, selectStmt.FromClause, scope); ok {
			result.names = append(result.names, star.names...)
			result.sources = append(result.sources, star.sources...)
			continue
		}
		if _, isStar := redshiftStarTarget(target); isStar {
			// The star covers a relation with no function projection. Its true
			// width is unknown, so keep a single placeholder and mark the
			// projection opaque; a star over this relation later stays
			// unexpanded rather than misaligning the following columns.
			result.opaque = true
		}
		sources, sourceIssues, err := q.collectFunctionBodyResultSources(target.Val, functionStack, cache, scope)
		if err != nil {
			return functionSourceProjection{}, functionBodyIssues{}, err
		}
		issues.merge(sourceIssues)
		result.names = append(result.names, redshiftTargetOutputName(target))
		result.sources = append(result.sources, sources)
	}
	return result, issues, nil
}

// functionStarProjection expands a star select item into the output names and
// source sets of its columns. Names and sources grow together, so a projection
// built from `SELECT *, fn() AS y` keeps y aligned with the function's sources
// rather than with the star's first column.
func (q *querySpanResolver) functionStarProjection(target *ast.ResTarget, from *ast.List, scope functionSourceScope) (functionSourceProjection, bool) {
	qualifier, isStar := redshiftStarTarget(target)
	if !isStar {
		return functionSourceProjection{}, false
	}
	if qualifier != "" {
		projection, ok := scope.bindings[strings.ToLower(qualifier)]
		if ok && !projection.opaque {
			return projection, true
		}
		if from != nil {
			for _, item := range from.Items {
				if projection, ok := q.functionProjectionForQualifier(item, scope, qualifier); ok {
					return projection, true
				}
			}
		}
		return functionSourceProjection{}, false
	}
	return q.functionProjectionForList(from, scope)
}

// functionProjectionForQualifier finds the FROM item a qualified star names.
func (q *querySpanResolver) functionProjectionForQualifier(node ast.Node, scope functionSourceScope, qualifier string) (functionSourceProjection, bool) {
	switch n := node.(type) {
	case *ast.RangeVar:
		if n == nil {
			return functionSourceProjection{}, false
		}
		name := n.Relname
		if n.Alias != nil && n.Alias.Aliasname != "" {
			name = n.Alias.Aliasname
		}
		if !strings.EqualFold(name, qualifier) {
			return functionSourceProjection{}, false
		}
		projection, ok := q.functionProjectionForNode(n, scope)
		if !ok {
			return functionSourceProjection{}, false
		}
		return projection.flatten()
	case *ast.RangeSubselect:
		if n == nil || n.Alias == nil || !strings.EqualFold(n.Alias.Aliasname, qualifier) {
			return functionSourceProjection{}, false
		}
		projection, ok := q.functionProjectionForNode(n, scope)
		if !ok {
			return functionSourceProjection{}, false
		}
		return projection.flatten()
	case *ast.JoinExpr:
		if n == nil {
			return functionSourceProjection{}, false
		}
		if projection, ok := q.functionProjectionForQualifier(n.Larg, scope, qualifier); ok {
			return projection, true
		}
		return q.functionProjectionForQualifier(n.Rarg, scope, qualifier)
	default:
		return functionSourceProjection{}, false
	}
}

// functionProjectionForList returns the output of a whole FROM list, or false
// when an item has no function projection or an unresolved star width.
func (q *querySpanResolver) functionProjectionForList(from *ast.List, scope functionSourceScope) (functionSourceProjection, bool) {
	var tree *functionProjection
	if from != nil {
		for _, item := range from.Items {
			projection, ok := q.functionProjectionForNode(item, scope)
			if !ok {
				return functionSourceProjection{}, false
			}
			if tree == nil {
				tree = projection
			} else {
				tree = &functionProjection{left: tree, right: projection}
			}
		}
	}
	if tree == nil {
		return functionSourceProjection{}, true
	}
	result, ok := tree.flatten()
	if !ok || result.opaque {
		return functionSourceProjection{}, false
	}
	return result, true
}

// functionProjectionForNode returns one FROM item's output names and function
// sources. A USING/NATURAL join coalesces its common columns into one output
// column, so expanding `*` over it must not list both inputs' copies.
func (q *querySpanResolver) functionProjectionForNode(node ast.Node, scope functionSourceScope) (*functionProjection, bool) {
	switch n := node.(type) {
	case *ast.RangeVar:
		if n == nil {
			return nil, false
		}
		name := n.Relname
		if n.Alias != nil && n.Alias.Aliasname != "" {
			name = n.Alias.Aliasname
		}
		projection, ok := scope.bindings[strings.ToLower(name)]
		if !ok {
			// A base table has no function lineage of its own, but its column
			// count is still needed to expand a star to the right width.
			projection, ok = q.functionBaseTableProjection(n)
		}
		if !ok {
			return nil, false
		}
		return flatFunctionProjection(projection), true
	case *ast.RangeSubselect:
		if n == nil || n.Alias == nil || n.Alias.Aliasname == "" {
			return nil, false
		}
		projection, ok := scope.bindings[strings.ToLower(n.Alias.Aliasname)]
		if !ok {
			return nil, false
		}
		return flatFunctionProjection(projection), true
	case *ast.JoinExpr:
		if n == nil {
			return nil, false
		}
		left, ok := q.functionProjectionForNode(n.Larg, scope)
		if !ok {
			return nil, false
		}
		right, ok := q.functionProjectionForNode(n.Rarg, scope)
		if !ok {
			return nil, false
		}
		using := redshiftStringList(n.UsingClause)
		if len(using) == 0 && !n.IsNatural {
			return &functionProjection{left: left, right: right}, true
		}
		leftFlat, ok := left.flatten()
		if !ok {
			return nil, false
		}
		rightFlat, ok := right.flatten()
		if !ok {
			return nil, false
		}
		merged, ok := mergeFunctionJoinProjection(leftFlat, rightFlat, n)
		if !ok {
			return nil, false
		}
		return flatFunctionProjection(merged), true
	default:
		return nil, false
	}
}

// functionBaseTableProjection builds a base relation's output names from
// metadata, with no lineage: it exists so a star over it has the right width.
func (q *querySpanResolver) functionBaseTableProjection(n *ast.RangeVar) (functionSourceProjection, bool) {
	if n == nil || q.opts.GetDatabaseMetadata == nil {
		return functionSourceProjection{}, false
	}
	database := q.defaultDatabase
	if n.Catalogname != "" {
		database = n.Catalogname
	}
	meta, err := q.getDatabaseMetadata(database)
	if err != nil || meta == nil {
		return functionSourceProjection{}, false
	}
	schemaName := n.Schemaname
	if schemaName == "" {
		schemaName, _ = searchRelation(meta, q.searchPath, n.Relname)
	}
	schema := schemaMetadata(meta, schemaName)
	if schema == nil {
		return functionSourceProjection{}, false
	}
	var names []string
	switch {
	case findFunctionTable(schema, n.Relname) != nil:
		for _, column := range findFunctionTable(schema, n.Relname).GetColumns() {
			names = append(names, column.GetName())
		}
	case findFunctionView(schema, n.Relname) != nil:
		for _, column := range findFunctionView(schema, n.Relname).GetColumns() {
			names = append(names, column.GetName())
		}
	case findFunctionExternalTable(schema, n.Relname) != nil:
		for _, column := range findFunctionExternalTable(schema, n.Relname).GetColumns() {
			names = append(names, column.GetName())
		}
	default:
		return functionSourceProjection{}, false
	}
	if n.Alias != nil {
		names = applyFunctionColumnAliases(names, redshiftStringList(n.Alias.Colnames))
	}
	return functionSourceProjection{names: names, sources: make([]ColumnSet, len(names))}, true
}

func findFunctionTable(schema *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	for _, table := range schema.GetTables() {
		if table.GetName() == name {
			return table
		}
	}
	return nil
}

func findFunctionView(schema *metadata.SchemaMetadata, name string) *metadata.ViewMetadata {
	for _, view := range schema.GetViews() {
		if view.GetName() == name {
			return view
		}
	}
	return nil
}

func findFunctionExternalTable(schema *metadata.SchemaMetadata, name string) *metadata.ExternalTableMetadata {
	for _, table := range schema.GetExternalTables() {
		if table.GetName() == name {
			return table
		}
	}
	return nil
}

// appendFunctionProjection concatenates one projection onto another.
func appendFunctionProjection(dst *functionSourceProjection, src functionSourceProjection) {
	dst.names = append(dst.names, src.names...)
	dst.sources = append(dst.sources, src.sources...)
	dst.opaque = dst.opaque || src.opaque
}

// mergeFunctionJoinProjection combines the two sides of a join. A plain join
// concatenates them; a USING/NATURAL join lists each common column once, first,
// unioning its sources, then the remaining left and right columns.
func mergeFunctionJoinProjection(left, right functionSourceProjection, n *ast.JoinExpr) (functionSourceProjection, bool) {
	if left.opaque || right.opaque {
		return functionSourceProjection{}, false
	}
	using := redshiftStringList(n.UsingClause)
	if len(using) == 0 && !n.IsNatural {
		var result functionSourceProjection
		appendFunctionProjection(&result, left)
		appendFunctionProjection(&result, right)
		return result, true
	}
	var names []string
	if n.IsNatural {
		seen := make(map[string]bool, len(left.names))
		for _, name := range left.names {
			key := strings.ToLower(name)
			if seen[key] {
				continue
			}
			if _, ok := functionColumnIndex(right, name); !ok {
				continue
			}
			seen[key] = true
			names = append(names, name)
		}
	} else {
		names = using
	}
	usedLeft := make([]bool, len(left.names))
	usedRight := make([]bool, len(right.names))
	var result functionSourceProjection
	for _, name := range names {
		li, ok := functionColumnIndex(left, name)
		if !ok {
			return functionSourceProjection{}, false
		}
		ri, ok := functionColumnIndex(right, name)
		if !ok {
			return functionSourceProjection{}, false
		}
		merged := make(ColumnSet)
		mergeSourceColumnSet(merged, left.sources[li])
		mergeSourceColumnSet(merged, right.sources[ri])
		result.names = append(result.names, left.names[li])
		result.sources = append(result.sources, merged)
		usedLeft[li] = true
		usedRight[ri] = true
	}
	for i, name := range left.names {
		if !usedLeft[i] {
			result.names = append(result.names, name)
			result.sources = append(result.sources, left.sources[i])
		}
	}
	for i, name := range right.names {
		if !usedRight[i] {
			result.names = append(result.names, name)
			result.sources = append(result.sources, right.sources[i])
		}
	}
	return result, true
}

func functionColumnIndex(p functionSourceProjection, name string) (int, bool) {
	for i, candidate := range p.names {
		if strings.EqualFold(candidate, name) && i < len(p.sources) {
			return i, true
		}
	}
	return 0, false
}

// redshiftStarTarget reports the relation qualifier of a star select item and
// whether the target is a star at all. `*` yields "", true; `rel.*` yields
// "rel", true; every other expression yields false.
func redshiftStarTarget(target *ast.ResTarget) (string, bool) {
	if target == nil {
		return "", false
	}
	ref, ok := target.Val.(*ast.ColumnRef)
	if !ok || ref.Fields == nil || ref.Fields.Len() == 0 {
		return "", false
	}
	items := ref.Fields.Items
	if _, ok := items[len(items)-1].(*ast.A_Star); !ok {
		return "", false
	}
	if len(items) == 1 {
		return "", true
	}
	name, ok := items[len(items)-2].(*ast.String)
	if !ok {
		return "", false
	}
	return name.Str, true
}

// mergeFunctionArmResultSources merges set-operation arms positionally: output
// column i carries the union of the arms' column i sources.
func mergeFunctionArmResultSources(left, right []ColumnSet) []ColumnSet {
	n := len(left)
	if len(right) > n {
		n = len(right)
	}
	result := make([]ColumnSet, n)
	for i := 0; i < n; i++ {
		merged := make(ColumnSet)
		if i < len(left) {
			mergeSourceColumnSet(merged, left[i])
		}
		if i < len(right) {
			mergeSourceColumnSet(merged, right[i])
		}
		result[i] = merged
	}
	return result
}

func (q *querySpanResolver) functionSourceScopeForSelect(selectStmt *ast.SelectStmt, functionStack map[string]bool, cache map[string]functionBodySources, outerScope functionSourceScope) (functionSourceScope, functionBodyIssues, error) {
	scope := enterFunctionSourceScope(outerScope)
	issues := functionBodyIssues{}

	if selectStmt.WithClause != nil && selectStmt.WithClause.Ctes != nil {
		cteScope := cloneFunctionSourceScope(scope)
		for _, item := range selectStmt.WithClause.Ctes.Items {
			cte, ok := item.(*ast.CommonTableExpr)
			if !ok || cte.Ctename == "" {
				continue
			}
			cteSelect, ok := cte.Ctequery.(*ast.SelectStmt)
			if !ok {
				continue
			}
			projection, sourceIssues, err := q.collectSelectFunctionProjection(cteSelect, functionStack, cache, cteScope)
			if err != nil {
				return functionSourceScope{}, functionBodyIssues{}, err
			}
			issues.merge(sourceIssues)
			projection.names = applyFunctionColumnAliases(projection.names, redshiftStringList(cte.Aliascolnames))
			name := strings.ToLower(cte.Ctename)
			scope.defs[name] = projection
			cteScope.defs[name] = projection
		}
	}

	sourceIssues, err := q.collectFunctionSourceScopeFromList(selectStmt.FromClause, functionStack, cache, scope)
	if err != nil {
		return functionSourceScope{}, functionBodyIssues{}, err
	}
	issues.merge(sourceIssues)
	if selectStmt.FromClause != nil {
		from := selectStmt.FromClause
		projectionScope := scope
		scope.projection = &functionProjectionSource{
			compute: func() (functionSourceProjection, bool) {
				return q.functionProjectionForList(from, projectionScope)
			},
		}
	}
	return scope, issues, nil
}

func (q *querySpanResolver) collectFunctionSourceScopeFromList(list *ast.List, functionStack map[string]bool, cache map[string]functionBodySources, scope functionSourceScope) (functionBodyIssues, error) {
	if list == nil {
		return functionBodyIssues{}, nil
	}
	issues := functionBodyIssues{}
	for _, item := range list.Items {
		sourceIssues, err := q.collectFunctionSourceScopeFromNode(item, functionStack, cache, scope)
		if err != nil {
			return functionBodyIssues{}, err
		}
		issues.merge(sourceIssues)
	}
	return issues, nil
}

func (q *querySpanResolver) collectFunctionSourceScopeFromNode(node ast.Node, functionStack map[string]bool, cache map[string]functionBodySources, scope functionSourceScope) (functionBodyIssues, error) {
	switch n := node.(type) {
	case *ast.RangeVar:
		if n == nil {
			return functionBodyIssues{}, nil
		}
		projection, ok := scope.defs[strings.ToLower(n.Relname)]
		if !ok {
			return functionBodyIssues{}, nil
		}
		if n.Alias != nil && n.Alias.Aliasname != "" {
			projection.names = applyFunctionColumnAliases(projection.names, redshiftStringList(n.Alias.Colnames))
			scope.bindings[strings.ToLower(n.Alias.Aliasname)] = projection
			return functionBodyIssues{}, nil
		}
		scope.bindings[strings.ToLower(n.Relname)] = projection
		return functionBodyIssues{}, nil
	case *ast.RangeSubselect:
		if n == nil || n.Alias == nil || n.Alias.Aliasname == "" {
			return functionBodyIssues{}, nil
		}
		selectStmt, ok := n.Subquery.(*ast.SelectStmt)
		if !ok {
			return functionBodyIssues{}, nil
		}
		projection, issues, err := q.collectSelectFunctionProjection(selectStmt, functionStack, cache, scope)
		if err != nil {
			return functionBodyIssues{}, err
		}
		projection.names = applyFunctionColumnAliases(projection.names, redshiftStringList(n.Alias.Colnames))
		scope.bindings[strings.ToLower(n.Alias.Aliasname)] = projection
		return issues, nil
	case *ast.JoinExpr:
		if n == nil {
			return functionBodyIssues{}, nil
		}
		leftIssues, err := q.collectFunctionSourceScopeFromNode(n.Larg, functionStack, cache, scope)
		if err != nil {
			return functionBodyIssues{}, err
		}
		rightIssues, err := q.collectFunctionSourceScopeFromNode(n.Rarg, functionStack, cache, scope)
		if err != nil {
			return functionBodyIssues{}, err
		}
		leftIssues.merge(rightIssues)
		return leftIssues, nil
	default:
		return functionBodyIssues{}, nil
	}
}

func (q *querySpanResolver) collectFunctionBodyResultSources(node ast.Node, functionStack map[string]bool, cache map[string]functionBodySources, scope functionSourceScope) (ColumnSet, functionBodyIssues, error) {
	result := make(ColumnSet)
	if node == nil {
		return result, functionBodyIssues{}, nil
	}

	var err error
	issues := functionBodyIssues{}
	ast.Inspect(node, func(n ast.Node) bool {
		if err != nil {
			return false
		}
		if sub, ok := n.(*ast.SubLink); ok {
			if sub.Testexpr != nil {
				sources, sourceIssues, sourceErr := q.collectFunctionBodyResultSources(sub.Testexpr, functionStack, cache, scope)
				if sourceErr != nil {
					err = sourceErr
					return false
				}
				issues.merge(sourceIssues)
				mergeSourceColumnSet(result, sources)
			}
			selectStmt, ok := sub.Subselect.(*ast.SelectStmt)
			if !ok {
				return false
			}
			sources, sourceIssues, sourceErr := q.collectSubqueryFunctionSources(selectStmt, functionStack, cache, scope)
			if sourceErr != nil {
				err = sourceErr
				return false
			}
			issues.merge(sourceIssues)
			mergeSourceColumnSet(result, sources)
			return false
		}
		if columnRef, ok := n.(*ast.ColumnRef); ok {
			mergeSourceColumnSet(result, resolveFunctionColumnRef(scope, columnRef))
			return true
		}
		call, ok := n.(*ast.FuncCall)
		if !ok {
			return true
		}
		sources, sourceErr := q.functionBodySourcesForCall(call, functionStack, cache)
		if sourceErr != nil {
			err = sourceErr
			return false
		}
		issues.merge(sources.issues)
		for _, resultColumn := range sources.results {
			for source := range resultColumn.SourceColumns {
				result[source] = true
			}
		}
		return true
	})
	return result, issues, err
}

// collectSubqueryFunctionSources resolves everything a subquery reads in its
// own scopes: each set-operation arm gets its own FROM bindings, and a simple
// arm collects its target and predicate functions together.
func (q *querySpanResolver) collectSubqueryFunctionSources(selectStmt *ast.SelectStmt, functionStack map[string]bool, cache map[string]functionBodySources, outerScope functionSourceScope) (ColumnSet, functionBodyIssues, error) {
	result := make(ColumnSet)
	issues := functionBodyIssues{}
	scope, sourceIssues, err := q.functionSourceScopeForSelect(selectStmt, functionStack, cache, outerScope)
	if err != nil {
		return nil, functionBodyIssues{}, err
	}
	issues.merge(sourceIssues)
	if selectStmt.Op != ast.SETOP_NONE {
		left, leftIssues, err := q.collectSubqueryFunctionSources(selectStmt.Larg, functionStack, cache, scope)
		if err != nil {
			return nil, functionBodyIssues{}, err
		}
		issues.merge(leftIssues)
		mergeSourceColumnSet(result, left)
		right, rightIssues, err := q.collectSubqueryFunctionSources(selectStmt.Rarg, functionStack, cache, scope)
		if err != nil {
			return nil, functionBodyIssues{}, err
		}
		issues.merge(rightIssues)
		mergeSourceColumnSet(result, right)
		return result, issues, nil
	}
	sources, subIssues, sourceErr := q.collectFunctionBodyResultSources(selectStmt, functionStack, cache, scope)
	if sourceErr != nil {
		return nil, functionBodyIssues{}, sourceErr
	}
	issues.merge(subIssues)
	mergeSourceColumnSet(result, sources)
	return result, issues, nil
}

func (q *querySpanResolver) functionBodySourcesForCall(call *ast.FuncCall, functionStack map[string]bool, cache map[string]functionBodySources) (functionBodySources, error) {
	schemaName, functionName := functionCallName(call.Funcname)
	if functionName == "" || IsSystemFunction(functionName, "") {
		return functionBodySources{accessTables: ColumnSet{}}, nil
	}

	functionSchema, function, findErr := q.findFunctionMetadata(schemaName, functionName, call.Args.Len())
	if findErr != nil || function == nil {
		return functionBodySources{accessTables: ColumnSet{}}, findErr
	}

	cacheKey := strings.ToLower(functionSchema + "." + function.GetName() + "/" + function.GetSignature())
	if sources, ok := cache[cacheKey]; ok {
		return sources, nil
	}
	if functionStack[cacheKey] {
		return functionBodySources{accessTables: ColumnSet{}}, nil
	}

	body, ok, bodyErr := extractSQLFunctionBody(function.GetDefinition())
	if bodyErr != nil {
		// Body extraction failures should fail closed through QuerySpan.FunctionNotSupportedError.
		//nolint:nilerr
		return functionBodySources{
			accessTables: ColumnSet{},
			issues: functionBodyIssues{
				functionNotSupportedError: &FunctionNotSupportedError{
					Err:        bodyErr,
					Function:   functionSchema + "." + function.GetName(),
					Definition: function.GetDefinition(),
				},
			},
		}, nil
	}
	if !ok {
		return functionBodySources{
			accessTables: ColumnSet{},
			issues: functionBodyIssues{
				functionNotSupportedError: &FunctionNotSupportedError{
					Err:        errors.New("unsupported or empty SQL function body"),
					Function:   functionSchema + "." + function.GetName(),
					Definition: function.GetDefinition(),
				},
			},
		}, nil
	}

	functionStack[cacheKey] = true
	span, querySpanErr := q.getQuerySpanWithFunctionStack(body, functionStack)
	delete(functionStack, cacheKey)
	if querySpanErr != nil {
		return functionBodySources{}, fmt.Errorf("failed to get query span for function: %s.%s: %w", functionSchema, functionName, querySpanErr)
	}
	sources := functionBodySources{
		accessTables: span.SourceColumns,
		results:      span.Results,
		issues: functionBodyIssues{
			notFoundError:             span.NotFoundError,
			functionNotSupportedError: span.FunctionNotSupportedError,
		},
	}
	if sources.issues.notFoundError == nil && sources.issues.functionNotSupportedError == nil {
		cache[cacheKey] = sources
	}
	return sources, nil
}

func (q *querySpanResolver) findFunctionMetadata(schemaName, functionName string, argCount int) (string, *metadata.FunctionMetadata, error) {
	if q.opts.GetDatabaseMetadata == nil {
		return "", nil, nil
	}
	metadata, err := q.getDatabaseMetadata(q.defaultDatabase)
	if err != nil {
		return "", nil, err
	}
	if metadata == nil {
		return "", nil, nil
	}

	searchPath := q.searchPath
	if schemaName != "" {
		searchPath = []string{schemaName}
	}
	schemas, functions := searchFunctions(metadata, searchPath, functionName)
	if len(functions) == 0 {
		return "", nil, nil
	}
	if len(functions) == 1 {
		return schemas[0], functions[0], nil
	}
	for i, function := range functions {
		count, ok := functionParameterCount(function.GetDefinition())
		if ok && count == argCount {
			return schemas[i], function, nil
		}
	}
	return "", nil, nil
}

func functionCallName(list *ast.List) (string, string) {
	if list == nil {
		return "", ""
	}
	parts := make([]string, 0, list.Len())
	for _, item := range list.Items {
		if value, ok := item.(*ast.String); ok {
			parts = append(parts, value.Str)
		}
	}
	switch len(parts) {
	case 1:
		return "", parts[0]
	case 2:
		return parts[0], parts[1]
	case 3:
		return parts[1], parts[2]
	default:
		return "", ""
	}
}

func functionParameterCount(definition string) (int, bool) {
	stmts, err := redshift.Parse(definition)
	if err != nil || len(stmts) != 1 {
		return 0, false
	}
	stmt, ok := stmts[0].AST.(*ast.CreateFunctionStmt)
	if !ok || stmt.Parameters == nil {
		return 0, ok
	}
	return stmt.Parameters.Len(), true
}

func extractSQLFunctionBody(definition string) (string, bool, error) {
	stmts, err := redshift.Parse(definition)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse function definition: %w", err)
	}
	if len(stmts) != 1 {
		return "", false, nil
	}
	stmt, ok := stmts[0].AST.(*ast.CreateFunctionStmt)
	if !ok {
		return "", false, nil
	}

	language := ""
	body := ""
	if stmt.Options != nil {
		for _, item := range stmt.Options.Items {
			elem, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(elem.Defname) {
			case "language":
				language = defElemString(elem.Arg)
			case "as":
				body = defElemString(elem.Arg)
			default:
			}
		}
	}
	if language != "" && !strings.EqualFold(language, "sql") {
		return "", false, nil
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return "", false, nil
	}
	return body, true, nil
}

func defElemString(node ast.Node) string {
	switch value := node.(type) {
	case *ast.String:
		return value.Str
	case *ast.List:
		if len(value.Items) == 0 {
			return ""
		}
		return defElemString(value.Items[0])
	default:
		return ""
	}
}

func mergeFunctionResultSources(results []ResolvedQuerySpanResult, functionResultSources []functionResultSource) {
	if len(functionResultSources) == len(results) {
		for i, source := range functionResultSources {
			for column := range source.Sources {
				results[i].SourceColumns[column] = true
			}
		}
		return
	}
	// A partially expanded star makes the source list narrower than the result
	// list. Match by output name rather than filling the first empty result,
	// which would attach a later function's sources to the star's column.
	for _, source := range functionResultSources {
		if len(source.Sources) == 0 || source.Name == "" {
			continue
		}
		for i := range results {
			if len(results[i].SourceColumns) != 0 || !strings.EqualFold(results[i].Name, source.Name) {
				continue
			}
			for column := range source.Sources {
				results[i].SourceColumns[column] = true
			}
			break
		}
	}
}

func cloneFunctionSourceScope(scope functionSourceScope) functionSourceScope {
	result := functionSourceScope{
		defs:        make(map[string]functionSourceProjection, len(scope.defs)),
		bindings:    make(map[string]functionSourceProjection, len(scope.bindings)),
		projection:  scope.projection,
		outerLevels: append([]functionScopeLevel(nil), scope.outerLevels...),
	}
	for name, projection := range scope.defs {
		result.defs[name] = projection
	}
	for name, projection := range scope.bindings {
		result.bindings[name] = projection
	}
	return result
}

// enterFunctionSourceScope starts a nested query's scope. It inherits the CTE
// definitions and makes the enclosing query's bindings fallbacks, not peers of
// the nested FROM bindings.
func enterFunctionSourceScope(outer functionSourceScope) functionSourceScope {
	result := functionSourceScope{
		defs:        make(map[string]functionSourceProjection, len(outer.defs)),
		bindings:    make(map[string]functionSourceProjection),
		outerLevels: append([]functionScopeLevel(nil), outer.outerLevels...),
	}
	for name, projection := range outer.defs {
		result.defs[name] = projection
	}
	if len(outer.bindings) > 0 || outer.projection != nil {
		result.outerLevels = append(result.outerLevels, functionScopeLevel{
			bindings:   outer.bindings,
			projection: outer.projection,
		})
	}
	return result
}

func mergeSourceColumnSet(dst, src ColumnSet) {
	for source := range src {
		dst[source] = true
	}
}

func resolveFunctionColumnRef(scope functionSourceScope, columnRef *ast.ColumnRef) ColumnSet {
	parts := redshiftColumnRefParts(columnRef)
	switch len(parts) {
	case 1:
		return resolveUnqualifiedFunctionColumn(scope, parts[0])
	case 2:
		key := strings.ToLower(parts[0])
		if projection, ok := scope.bindings[key]; ok {
			return projection.functionSourcesForColumn(parts[1])
		}
		for i := len(scope.outerLevels) - 1; i >= 0; i-- {
			if projection, ok := scope.outerLevels[i].bindings[key]; ok {
				return projection.functionSourcesForColumn(parts[1])
			}
		}
		return ColumnSet{}
	default:
		return ColumnSet{}
	}
}

func resolveUnqualifiedFunctionColumn(scope functionSourceScope, columnName string) ColumnSet {
	if sources, handled := resolveUnqualifiedScopeLevel(scope.bindings, scope.projection, columnName); handled {
		return sources
	}
	for i := len(scope.outerLevels) - 1; i >= 0; i-- {
		level := scope.outerLevels[i]
		if sources, handled := resolveUnqualifiedScopeLevel(level.bindings, level.projection, columnName); handled {
			return sources
		}
	}
	return ColumnSet{}
}

// resolveUnqualifiedScopeLevel searches one query scope. handled is true when
// the level has the column (sources may be empty) or reports it ambiguous.
func resolveUnqualifiedScopeLevel(bindings map[string]functionSourceProjection, projection *functionProjectionSource, columnName string) (ColumnSet, bool) {
	if p, ok := projection.get(); ok {
		if index, ok := functionColumnIndex(p, columnName); ok {
			return p.sources[index], true
		}
		return nil, false
	}
	sources, found, ambiguous := lookupUnqualifiedFunctionColumn(bindings, columnName)
	if ambiguous {
		return ColumnSet{}, true
	}
	return sources, found
}

// lookupUnqualifiedFunctionColumn searches one scope level. found reports a
// single visible relation exposing the name (its sources may be empty);
// ambiguous reports that several do.
func lookupUnqualifiedFunctionColumn(bindings map[string]functionSourceProjection, columnName string) (sources ColumnSet, found, ambiguous bool) {
	for _, projection := range bindings {
		index, ok := functionColumnIndex(projection, columnName)
		if !ok {
			continue
		}
		if found {
			return nil, false, true
		}
		found = true
		sources = projection.sources[index]
	}
	return sources, found, false
}

// applyFunctionColumnAliases renames the prefix an alias list names and keeps
// the remaining output columns, matching the catalog's applyColumnAliases.
func applyFunctionColumnAliases(names, aliases []string) []string {
	if len(aliases) == 0 {
		return names
	}
	out := append([]string(nil), names...)
	for i, alias := range aliases {
		if i >= len(out) {
			break
		}
		if alias != "" {
			out[i] = alias
		}
	}
	return out
}

func (p functionSourceProjection) functionSourcesForColumn(columnName string) ColumnSet {
	for i, name := range p.names {
		if !strings.EqualFold(name, columnName) || i >= len(p.sources) {
			continue
		}
		return p.sources[i]
	}
	return ColumnSet{}
}

func redshiftTargetOutputName(target *ast.ResTarget) string {
	if target == nil {
		return ""
	}
	if target.Name != "" {
		return target.Name
	}
	if parts := redshiftColumnRefParts(asColumnRef(target.Val)); len(parts) > 0 {
		return parts[len(parts)-1]
	}
	if call, ok := target.Val.(*ast.FuncCall); ok {
		_, functionName := functionCallName(call.Funcname)
		return functionName
	}
	return ""
}

func asColumnRef(node ast.Node) *ast.ColumnRef {
	columnRef, ok := node.(*ast.ColumnRef)
	if !ok {
		return nil
	}
	return columnRef
}

func redshiftColumnRefParts(columnRef *ast.ColumnRef) []string {
	if columnRef == nil || columnRef.Fields == nil {
		return nil
	}
	result := make([]string, 0, columnRef.Fields.Len())
	for _, item := range columnRef.Fields.Items {
		value, ok := item.(*ast.String)
		if !ok {
			return nil
		}
		result = append(result, value.Str)
	}
	return result
}

func redshiftStringList(list *ast.List) []string {
	if list == nil {
		return nil
	}
	result := make([]string, 0, list.Len())
	for _, item := range list.Items {
		value, ok := item.(*ast.String)
		if !ok {
			continue
		}
		result = append(result, value.Str)
	}
	return result
}

func cloneCTENameSet(ctes map[string]bool) map[string]bool {
	result := make(map[string]bool, len(ctes))
	for name, ok := range ctes {
		result[name] = ok
	}
	return result
}

func isCTEReference(rangeVar *ast.RangeVar, ctes map[string]bool) bool {
	if rangeVar == nil || rangeVar.Schemaname != "" || rangeVar.Catalogname != "" {
		return false
	}
	return ctes[strings.ToLower(rangeVar.Relname)]
}

func hasResultWithoutSource(results []ResolvedQuerySpanResult) bool {
	for _, result := range results {
		if len(result.SourceColumns) == 0 {
			return true
		}
	}
	return false
}

func hasRangeVar(node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		if _, ok := n.(*ast.RangeVar); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

func (q *querySpanResolver) allAccessTablesExist(accessTables ColumnSet) bool {
	if len(accessTables) == 0 || q.opts.GetDatabaseMetadata == nil {
		return false
	}
	for resource := range accessTables {
		if !q.accessTableExists(resource) {
			return false
		}
	}
	return true
}

func (q *querySpanResolver) accessTableExists(resource ColumnResource) bool {
	if isSystemResource(resource) {
		return true
	}
	database := resource.Database
	if database == "" {
		database = q.defaultDatabase
	}
	metadata, err := q.getDatabaseMetadata(database)
	if err != nil || metadata == nil {
		return false
	}
	if resource.Schema == "" {
		schemaName, objectName := searchRelation(metadata, q.searchPath, resource.Table)
		return schemaName != "" || objectName != ""
	}
	schema := schemaMetadata(metadata, resource.Schema)
	if schema == nil {
		return false
	}
	// Match the relation kinds used by searchRelation above.
	return schemaHasRelation(schema, resource.Table)
}

func (q *querySpanResolver) resolveRangeVar(rangeVar *ast.RangeVar) (ColumnResource, error) {
	resource := ColumnResource{
		Database: q.defaultDatabase,
		Schema:   rangeVar.Schemaname,
		Table:    rangeVar.Relname,
	}
	if rangeVar.Catalogname != "" {
		resource.Database = rangeVar.Catalogname
	}
	if resource.Schema != "" || isSystemResource(resource) || q.opts.GetDatabaseMetadata == nil {
		return resource, nil
	}

	metadata, err := q.getDatabaseMetadata(resource.Database)
	if err != nil {
		return ColumnResource{}, err
	}
	if metadata == nil {
		return resource, nil
	}
	schemaName, objectName := searchRelation(metadata, q.searchPath, resource.Table)
	if schemaName == "" && objectName == "" {
		return ColumnResource{}, &ResourceNotFoundError{
			Database: &resource.Database,
			Table:    &resource.Table,
		}
	}
	resource.Schema = schemaName
	resource.Table = objectName
	return resource, nil
}

func (q *querySpanResolver) convertResults(results []QuerySpanResult) []ResolvedQuerySpanResult {
	converted := make([]ResolvedQuerySpanResult, 0, len(results))
	for _, result := range results {
		converted = append(converted, ResolvedQuerySpanResult{
			Name:          result.Name,
			SourceColumns: q.convertColumnList(result.SourceColumns),
			IsPlainField:  result.IsPlainField,
		})
	}
	return converted
}

func (q *querySpanResolver) convertColumnList(columns []ColumnResource) ColumnSet {
	result := make(ColumnSet)
	for _, column := range columns {
		database := column.Database
		if database == "" {
			database = q.defaultDatabase
		}
		result[ColumnResource{
			Database: database,
			Schema:   column.Schema,
			Table:    column.Table,
			Column:   column.Column,
		}] = true
	}
	return result
}

func querySpanType(node ast.Node, allSystems bool) QueryType {
	switch n := node.(type) {
	case *ast.SelectStmt:
		if hasSelectIntoClause(n) {
			return DDL
		}
		if allSystems {
			return SelectInfoSchema
		}
		return Select
	case *ast.ExplainStmt:
		if hasExplainAnalyze(n) {
			return querySpanType(n.Query, allSystems)
		}
		return Explain
	case *ast.RedshiftShowStmt, *ast.RedshiftDescStmt, *ast.VariableShowStmt:
		return SelectInfoSchema
	case *ast.VariableSetStmt:
		return Select
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt,
		*ast.CopyStmt, *ast.UnloadStmt, *ast.RefreshMatViewStmt, *ast.CallStmt:
		return DML
	case *ast.CreateStmt, *ast.CreateTableAsStmt, *ast.ViewStmt, *ast.IndexStmt, *ast.CreateSeqStmt,
		*ast.CreateSchemaStmt, *ast.CreateFunctionStmt, *ast.CreatedbStmt, *ast.CreateTrigStmt,
		*ast.CreateEnumStmt, *ast.CreateDomainStmt, *ast.CreateEventTrigStmt, *ast.CreatePLangStmt,
		*ast.CreateFdwStmt, *ast.CreateForeignServerStmt, *ast.CreateForeignTableStmt,
		*ast.CreateUserMappingStmt, *ast.CreateExtensionStmt, *ast.CreateTableSpaceStmt,
		*ast.CreateAmStmt, *ast.CreatePolicyStmt, *ast.CreatePublicationStmt,
		*ast.CreateSubscriptionStmt, *ast.CreateStatsStmt, *ast.CreateOpClassStmt,
		*ast.CreateOpFamilyStmt, *ast.CreateCastStmt, *ast.CreateTransformStmt,
		*ast.CreateConversionStmt, *ast.CreateRangeStmt, *ast.CreateRoleStmt,
		*ast.DropStmt, *ast.DropdbStmt, *ast.DropRoleStmt, *ast.DropUserMappingStmt,
		*ast.DropSubscriptionStmt, *ast.DropTableSpaceStmt, *ast.DropOwnedStmt,
		*ast.AlterTableStmt, *ast.AlterTableMoveAllStmt, *ast.AlterSeqStmt,
		*ast.AlterEnumStmt, *ast.AlterDomainStmt, *ast.AlterObjectSchemaStmt,
		*ast.AlterOwnerStmt, *ast.AlterDatabaseStmt, *ast.AlterDatabaseSetStmt,
		*ast.AlterDatabaseRefreshCollStmt, *ast.AlterSystemStmt, *ast.AlterCollationStmt,
		*ast.AlterFunctionStmt, *ast.AlterEventTrigStmt, *ast.AlterFdwStmt,
		*ast.AlterForeignServerStmt, *ast.AlterUserMappingStmt, *ast.AlterExtensionStmt,
		*ast.AlterExtensionContentsStmt, *ast.AlterTableSpaceOptionsStmt, *ast.AlterPolicyStmt,
		*ast.AlterPublicationStmt, *ast.AlterSubscriptionStmt, *ast.AlterObjectDependsStmt,
		*ast.AlterOperatorStmt, *ast.AlterTypeStmt, *ast.AlterDefaultPrivilegesStmt,
		*ast.AlterTSDictionaryStmt, *ast.AlterTSConfigurationStmt, *ast.AlterStatsStmt,
		*ast.AlterOpFamilyStmt, *ast.AlterRoleStmt, *ast.AlterRoleSetStmt,
		*ast.RenameStmt, *ast.TruncateStmt, *ast.VacuumStmt, *ast.GrantStmt,
		*ast.GrantRoleStmt, *ast.CommentStmt, *ast.LockStmt, *ast.ClusterStmt,
		*ast.ReindexStmt, *ast.RuleStmt, *ast.CheckPointStmt, *ast.DiscardStmt,
		*ast.LoadStmt, *ast.ConstraintsSetStmt, *ast.FetchStmt, *ast.SecLabelStmt,
		*ast.DoStmt, *ast.ImportForeignSchemaStmt, *ast.ReassignOwnedStmt,
		*ast.RedshiftObjectStmt:
		return DDL
	default:
		return QueryTypeUnknown
	}
}

func hasSelectIntoClause(n *ast.SelectStmt) bool {
	if n == nil {
		return false
	}
	if n.IntoClause != nil {
		return true
	}
	if hasSelectIntoClause(n.Larg) {
		return true
	}
	if hasSelectIntoClause(n.Rarg) {
		return true
	}
	return false
}

func isQuerySpanNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "does not exist") || strings.Contains(msg, "not found")
}
