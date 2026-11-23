package analysis

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	omnipg "github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/catalog"
)

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase is the database the statement runs in. Its metadata
	// loads the catalog the statement is analyzed against.
	DefaultDatabase string
	// SearchPath resolves unqualified names. Empty means public.
	SearchPath []string
	// Catalog reads the instance's synced metadata. GetDatabase must be set.
	Catalog Catalog
}

// ResolveQuerySpan returns the column lineage of one PostgreSQL statement.
//
// A statement that does not parse returns a *SyntaxError. A statement that
// reads user and system tables together returns ErrMixUserSystemTables.
// Statements that are not queries return their QueryType and the tables they
// access, with no Results.
func ResolveQuerySpan(ctx context.Context, statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	return newQuerySpanExtractor(opts.DefaultDatabase, opts.SearchPath, opts.Catalog).getQuerySpan(ctx, statement)
}

// querySpanExtractor extracts query span information using the omni parser
// and catalog.
type querySpanExtractor struct {
	ctx             context.Context
	catalog         Catalog
	defaultDatabase string
	searchPath      []string
	// metaCache is a lazy-load cache for database metadata.
	// Use getDatabaseMetadata() instead of accessing directly.
	metaCache     map[string]Database
	cat           *catalog.Catalog
	funcBodyCache map[uint32]*funcBodyAnalysis
	// snapshotFuncs holds the snapshot's definition of every overload, keyed by
	// lowercase "schema.name". Used by analyzeFunctionBody for a function the
	// catalog installed as a stand-in, listed in standInFuncs by the same key;
	// the stand-in may stand for any overload of the name.
	snapshotFuncs map[string][]snapshotFunc
	standInFuncs  map[string]bool
	// funcSourceColumns accumulates table-level access (column="") discovered inside
	// function bodies. Merged into the top-level span's SourceColumns.
	funcSourceColumns ColumnSet
	// predicateColumns accumulates the columns that decide which rows the
	// statement returns; it becomes the span's PredicateColumns. Function body
	// analysis collects into a set of its own, cached per function, and merges
	// it into the caller's.
	predicateColumns ColumnSet
	// fallbackCTEMap holds CTE definitions during fallback column extraction.
	// Set by extractFallbackColumns and used by extractColumnsFromRangeVar
	// to resolve CTE references.
	fallbackCTEMap map[string]*ast.SelectStmt
	// fallbackCTEAliases holds CTE column aliases (e.g., WITH t1(cc1, cc2) AS (...)).
	// Keyed by lowercase CTE name.
	fallbackCTEAliases map[string][]string
	// fromColumns caches the fallback's fromItemColumns by FROM item, for
	// every analyzer of the statement; nil marks an item being resolved, so
	// an item reached again through itself, as a recursive CTE's reference
	// is, reads nothing.
	fromColumns map[ast.Node][]QuerySpanResult
}

func newQuerySpanExtractor(defaultDatabase string, searchPath []string, cat Catalog) *querySpanExtractor {
	if len(searchPath) == 0 {
		searchPath = []string{"public"}
	}
	return &querySpanExtractor{
		defaultDatabase:   defaultDatabase,
		searchPath:        searchPath,
		catalog:           cat,
		metaCache:         make(map[string]Database),
		funcBodyCache:     make(map[uint32]*funcBodyAnalysis),
		snapshotFuncs:     make(map[string][]snapshotFunc),
		standInFuncs:      make(map[string]bool),
		funcSourceColumns: make(ColumnSet),
		predicateColumns:  make(ColumnSet),
	}
}

// getDatabaseMetadata returns cached database metadata, fetching it if not yet cached.
func (e *querySpanExtractor) getDatabaseMetadata(database string) (Database, error) {
	if meta, ok := e.metaCache[database]; ok {
		return meta, nil
	}
	meta, err := e.catalog.GetDatabase(database)
	if err != nil {
		return nil, fmt.Errorf("failed to get database metadata for database: %s: %w", database, err)
	}
	e.metaCache[database] = meta
	return meta, nil
}

// initCatalog loads the database metadata into a new omni catalog. An object
// that fails to install gets a text-typed stand-in or is left out; either way
// the rest of the catalog still loads.
func (e *querySpanExtractor) initCatalog() error {
	meta, err := e.getDatabaseMetadata(e.defaultDatabase)
	if err != nil {
		return fmt.Errorf("failed to get database metadata for catalog init: %w", err)
	}

	e.cat = catalog.New()
	e.cat.SetSearchPath(e.searchPath)

	if meta == nil {
		return nil
	}

	report, err := e.cat.LoadMetadata(e.ctx, meta.Metadata(), catalog.LoadMetadataOptions{})
	if err != nil {
		return fmt.Errorf("failed to load the catalog for %q: %w", e.defaultDatabase, err)
	}
	for key := range report.Degraded {
		if name, ok := strings.CutPrefix(key, "func:"); ok {
			name, _, _ = strings.Cut(name, "|")
			e.standInFuncs[strings.ToLower(name)] = true
		}
	}

	for _, s := range meta.Metadata().GetSchemas() {
		for _, f := range s.GetFunctions() {
			if f.Definition == "" {
				continue
			}
			_, body, language := parseFuncDefForLineage(f.Definition)
			if body != "" {
				key := strings.ToLower(s.Name + "." + f.Name)
				e.snapshotFuncs[key] = append(e.snapshotFuncs[key], snapshotFunc{language: language, body: body})
			}
		}
	}

	return nil
}

// extractFuncBodyFromDef extracts the dollar-quoted function body from a
// CREATE FUNCTION definition string.
func extractFuncBodyFromDef(definition string) string {
	// Find dollar-quoted body.
	tag, bodyStart, bodyEnd := findDollarQuotedBody(definition)
	if tag == "" {
		return ""
	}
	return definition[bodyStart : bodyEnd-len(tag)]
}

// findDollarQuotedBody locates the dollar-quoted body in a function definition.
// Returns the tag, the start of body content (after opening tag), and the end
// position (after closing tag).
func findDollarQuotedBody(definition string) (tag string, bodyStart int, bodyEnd int) {
	for i := 0; i < len(definition); i++ {
		if definition[i] != '$' {
			continue
		}
		// Find end of tag.
		tagEnd := i + 1
		for tagEnd < len(definition) && definition[tagEnd] != '$' {
			tagEnd++
		}
		if tagEnd >= len(definition) {
			continue
		}
		tag := definition[i : tagEnd+1]
		contentStart := tagEnd + 1
		// Find closing tag.
		closeIdx := strings.Index(definition[contentStart:], tag)
		if closeIdx >= 0 {
			return tag, contentStart, contentStart + closeIdx + len(tag)
		}
	}
	return "", 0, 0
}

// unresolvedColumnsError checks stored metadata independently of column lineage.
// A name that a CTE in scope shadows reads the CTE and is not an access.
// Reads missing from the access set remain unchecked.
func (e *querySpanExtractor) unresolvedColumnsError(accesses ColumnSet) *UnresolvedColumnsError {
	seen := make(map[ColumnResource]bool, len(accesses))
	var unresolved []ColumnResource
	for _, access := range sortedColumnSet(accesses) {
		relation := ColumnResource{
			Database: access.Database,
			Schema:   access.Schema,
			Table:    access.Table,
		}
		if relation.Table == "" || seen[relation] {
			continue
		}
		seen[relation] = true
		if e.relationHasNoSyncedColumns(relation) {
			unresolved = append(unresolved, relation)
		}
	}
	if len(unresolved) == 0 {
		return nil
	}
	return &UnresolvedColumnsError{Relations: unresolved}
}

// relationHasNoSyncedColumns checks tables only: masking policies attach to
// table columns, and GetTable resolves partitions to their parent.
// A legal zero-column table is indistinguishable from degraded metadata and
// remains blocked even after re-syncing.
func (e *querySpanExtractor) relationHasNoSyncedColumns(relation ColumnResource) bool {
	meta, err := e.getDatabaseMetadata(relation.Database)
	if err != nil || meta == nil {
		return false
	}
	table := meta.GetTable(relation.Schema, relation.Table)
	return table != nil && len(table.GetColumns()) == 0
}

// getQuerySpan extracts the query span for the given SQL statement.
func (e *querySpanExtractor) getQuerySpan(ctx context.Context, stmt string) (*ResolvedSpan, error) {
	e.ctx = ctx

	// Step 1: Parse with omni.
	stmts, err := omnipg.Parse(stmt)
	if err != nil {
		return nil, &SyntaxError{Err: err}
	}
	if len(stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(stmts))
	}

	// Step 2: Extract access tables (non-fatal on error — some statement
	// types like SET don't have table references).
	accessesMap := make(ColumnSet)
	accessTables, err := e.accessTablesOfNode(stmts[0].AST)
	if err == nil {
		for _, resource := range accessTables {
			accessesMap[resource] = true
		}
	}

	// Step 3: Check for mixed system/user tables.
	allSystems, mixed := isMixedQuery(accessesMap)
	if mixed {
		return nil, ErrMixUserSystemTables
	}

	// Step 4: Classify query type.
	queryType, isExplainAnalyze := classifyQueryType(stmts[0].AST, allSystems)

	// Step 5: Return early for non-SELECT queries, keeping the tables the
	// statement touches and no result columns. EXPLAIN ANALYZE runs its query
	// and reports the rows each condition kept, so its conditions still count
	// as predicates.
	if queryType != Select || isExplainAnalyze {
		var unresolved *UnresolvedColumnsError
		if isExplainAnalyze {
			if err := e.collectExplainedPredicateColumns(stmts[0].AST, accessesMap); err != nil {
				return nil, err
			}
			e.mergeFunctionAccesses(accessesMap)
			// Columns a stale snapshot lacks can hide a condition, so the
			// caller syncs before running it, as for the SELECT itself.
			unresolved = e.unresolvedColumnsError(accessesMap)
		}
		return &ResolvedSpan{
			Type:                   queryType,
			ExplainAnalyze:         isExplainAnalyze,
			SourceColumns:          accessesMap,
			PredicateColumns:       e.predicateColumns,
			Results:                []QuerySpanResult{},
			UnresolvedColumnsError: unresolved,
		}, nil
	}

	// Step 6: Cast to SelectStmt, init catalog, analyze.
	selStmt, ok := stmts[0].AST.(*ast.SelectStmt)
	if !ok {
		// Not a SelectStmt (e.g., SET, SHOW classified as Select type).
		return &ResolvedSpan{
			Type:          Select,
			SourceColumns: accessesMap,
			Results:       []QuerySpanResult{},
		}, nil
	}

	if err := e.initCatalog(); err != nil {
		return nil, fmt.Errorf("failed to initialize catalog: %w", err)
	}

	query, err := e.analyzeSelect(selStmt)
	if err != nil {
		e.collectFallbackPredicateColumns(selStmt, accessesMap)
		// Before falling back, try to handle user-defined table-returning functions
		// that omni's AnalyzeSelectStmt can't resolve (e.g., RETURNS TABLE functions
		// used as table sources: SELECT * FROM func()).
		results := e.tryUserFuncTableSource(selStmt, accessesMap)
		if results == nil {
			// Fail-open: return access tables with best-effort column names and lineage
			// when analysis fails (e.g., unsupported built-in functions).
			results = e.extractFallbackColumns(selStmt)
		}
		e.mergeFunctionAccesses(accessesMap)
		return &ResolvedSpan{
			Type:                   Select,
			SourceColumns:          accessesMap,
			Results:                results,
			PredicateColumns:       e.predicateColumns,
			UnresolvedColumnsError: e.unresolvedColumnsError(accessesMap),
		}, nil
	}

	// Step 7: Extract lineage from analyzed query.
	// SourceColumns (table-level access) is already captured in accessesMap
	// from ExtractAccessTables above — no need to re-extract from the query.
	results := e.extractLineage(query, selStmt)
	e.collectStatementPredicateColumns(query, selStmt, accessesMap)

	e.mergeFunctionAccesses(accessesMap)

	return &ResolvedSpan{
		Type:                   Select,
		SourceColumns:          accessesMap,
		PredicateColumns:       e.predicateColumns,
		Results:                results,
		UnresolvedColumnsError: e.unresolvedColumnsError(accessesMap),
	}, nil
}

func (e *querySpanExtractor) collectExplainedPredicateColumns(stmt ast.Node, accesses ColumnSet) error {
	explain, ok := stmt.(*ast.ExplainStmt)
	if !ok {
		return nil
	}
	explained, ok := explain.Query.(*ast.SelectStmt)
	if !ok {
		return nil
	}
	if err := e.initCatalog(); err != nil {
		return fmt.Errorf("failed to initialize catalog: %w", err)
	}
	if query, err := e.analyzeSelect(explained); err == nil {
		e.collectStatementPredicateColumns(query, explained, accesses)
	} else {
		e.collectFallbackPredicateColumns(explained, accesses)
	}
	return nil
}

// mergeFunctionAccesses adds the tables the analyzed function bodies read to
// accesses, so access checks and the unresolved-columns signal see them.
func (e *querySpanExtractor) mergeFunctionAccesses(accesses ColumnSet) {
	for col := range e.funcSourceColumns {
		accesses[col] = true
	}
}

// analyzeSelect analyzes sel for lineage. A function the synced schema does not
// define, such as one from an extension, stays an opaque call over its
// arguments instead of failing the statement into the parse-tree fallback.
func (e *querySpanExtractor) analyzeSelect(sel *ast.SelectStmt) (*catalog.Query, error) {
	rewriteValuesLists(sel)
	return e.cat.AnalyzeSelectStmtWithOptions(sel, catalog.AnalyzeOptions{OpaqueUnknownFunctions: true})
}

// rewriteValuesLists rewrites, in place, every VALUES list in node into a
// one-row SELECT with PostgreSQL's column names (column1, ...), whose column
// is the row's expression, or for several rows the COALESCE of every row's
// expression at that position. COALESCE resolves its type over all the
// expressions at once, as VALUES does, so the rewrite analyzes whenever the
// VALUES list does. The catalog analyzes a VALUES list as a query with no
// output, which would leave its expressions without lineage and its
// subqueries' conditions unread. An INSERT's VALUES is left alone, since it
// may hold DEFAULT.
func rewriteValuesLists(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.InsertStmt:
			return false
		case *ast.SelectStmt:
			if v.ValuesLists == nil {
				return true
			}
			var columns [][]ast.Node
			for _, item := range v.ValuesLists.Items {
				row, ok := item.(*ast.List)
				if !ok {
					continue
				}
				for i, expr := range row.Items {
					if i == len(columns) {
						columns = append(columns, nil)
					}
					columns[i] = append(columns[i], expr)
				}
			}
			v.ValuesLists = nil
			targets := &ast.List{}
			for i, exprs := range columns {
				val := exprs[0]
				if len(exprs) > 1 {
					val = &ast.CoalesceExpr{Args: &ast.List{Items: exprs}, Loc: ast.NoLoc()}
				}
				targets.Items = append(targets.Items, &ast.ResTarget{Name: fmt.Sprintf("column%d", i+1), Val: val, Loc: ast.NoLoc()})
			}
			v.TargetList = targets
		default:
		}
		return true
	})
}

// outputColumn returns the target entry of q's attNum-th output column. Junk
// entries follow the output columns, so the position indexes the target list.
func outputColumn(q *catalog.Query, attNum int) *catalog.TargetEntry {
	if q == nil || attNum < 1 || attNum > len(q.TargetList) {
		return nil
	}
	return q.TargetList[attNum-1]
}

// extractLineage extracts the result column lineage from an analyzed query.
// selStmt is the original parse tree, used to determine IsPlainField
// (which is true only for columns that came from * expansion).
func (e *querySpanExtractor) extractLineage(q *catalog.Query, selStmt *ast.SelectStmt) []QuerySpanResult {
	return e.extractLineageWithVisited(nil, q, selStmt, make(map[*catalog.Query]bool))
}

// extractLineageWithVisited extracts the lineage of q, nested in the queries
// of parents, the set operations whose branch it is, which may define the
// CTEs it reads.
func (e *querySpanExtractor) extractLineageWithVisited(parents []*catalog.Query, q *catalog.Query, selStmt *ast.SelectStmt, visited map[*catalog.Query]bool) []QuerySpanResult {
	if q == nil {
		return nil
	}
	if visited[q] {
		return nil
	}
	visited[q] = true
	defer delete(visited, q)

	// Handle set operations (UNION/INTERSECT/EXCEPT).
	if q.SetOp != catalog.SetOpNone {
		return e.extractSetOpLineageWithVisited(parents, q, selStmt, visited)
	}

	// Build a plain-field mask from the parse tree's target list.
	// IsPlainField = true only for columns expanded from SELECT * or SELECT t.*.
	plainMask := buildPlainFieldMask(selStmt, q)

	// Collect parse tree ResTargets for name fixups.
	var parseTargets []*ast.ResTarget
	if selStmt != nil && selStmt.TargetList != nil {
		for _, item := range selStmt.TargetList.Items {
			if rt, ok := item.(*ast.ResTarget); ok {
				parseTargets = append(parseTargets, rt)
			}
		}
	}

	// Build a helper analyzer for parse-tree fallback extraction.
	var analyzer *plpgsqlAnalyzer
	var fromTables map[string]ColumnResource
	if selStmt != nil {
		analyzer = &plpgsqlAnalyzer{extractor: e, scope: newVariableScope(nil)}
		fromTables = analyzer.collectFromTables(selStmt)
	}

	var results []QuerySpanResult
	idx := 0
	for _, te := range q.TargetList {
		if te.ResJunk {
			continue
		}
		sourceColSet := make(ColumnSet)
		e.walkExprWithVisited(pushQuery(parents, q), te.Expr, sourceColSet, make(map[queryExprKey]bool))
		isPlain := idx < len(plainMask) && plainMask[idx]
		// Even if the parse tree says this came from *, the underlying
		// expression through CTEs/subqueries may not be a simple column ref.
		if isPlain {
			isPlain = isUltimatelyPlainColumn(q, te.Expr)
		}
		name := te.ResName
		// When the catalog reduces an expression to ConstExpr (losing column refs,
		// e.g., json_object('id': a)), fall back to the parse tree to extract
		// column references and derive a better name.
		if idx < len(parseTargets) {
			if name == "?column?" {
				if better := figureResTargetName(parseTargets[idx]); better != "" {
					name = better
				}
			}
			if len(sourceColSet) == 0 && analyzer != nil {
				analyzer.extractColumnRefsFromExpr(parseTargets[idx].Val, fromTables, sourceColSet)
			}
		}
		results = append(results, QuerySpanResult{
			Name:          name,
			SourceColumns: sourceColSet,
			IsPlainField:  isPlain,
		})
		idx++
	}
	return results
}

// isUltimatelyPlainColumn checks whether an expression, after resolving through
// CTEs and subqueries, is a simple column reference (VarExpr pointing to a
// physical table). Aggregates, functions, and other expressions return false.
func isUltimatelyPlainColumn(q *catalog.Query, expr catalog.AnalyzedExpr) bool {
	return isUltimatelyPlainColumnWithVisited(q, expr, make(map[queryExprKey]bool))
}

type queryExprKey struct {
	query *catalog.Query
	expr  catalog.AnalyzedExpr
}

func isUltimatelyPlainColumnWithVisited(q *catalog.Query, expr catalog.AnalyzedExpr, visited map[queryExprKey]bool) bool {
	if expr == nil {
		return false
	}
	key := queryExprKey{query: q, expr: expr}
	if visited[key] {
		return false
	}
	visited[key] = true

	v, ok := expr.(*catalog.VarExpr)
	if !ok {
		return false
	}
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return false
	}
	rte := q.RangeTable[v.RangeIdx]

	switch rte.Kind {
	case catalog.RTERelation:
		return true
	case catalog.RTESubquery:
		if te := outputColumn(rte.Subquery, int(v.AttNum)); te != nil {
			return isUltimatelyPlainColumnWithVisited(rte.Subquery, te.Expr, visited)
		}
	case catalog.RTECTE:
		// Without the enclosing queries a CTE defined outside q is unknown,
		// so its columns are not plain.
		if rte.CTEIndex >= 0 && rte.CTEIndex < len(q.CTEList) && q.CTEList[rte.CTEIndex].Name == rte.CTEName {
			cte := q.CTEList[rte.CTEIndex]
			if te := outputColumn(cte.Query, int(v.AttNum)); te != nil {
				return isUltimatelyPlainColumnWithVisited(cte.Query, te.Expr, visited)
			}
		}
	case catalog.RTEJoin:
		colIdx := int(v.AttNum - 1)
		if colIdx >= 0 && colIdx < len(rte.JoinAliasVars) {
			return isUltimatelyPlainColumnWithVisited(q, rte.JoinAliasVars[colIdx], visited)
		}
	default:
		return false
	}
	return false
}

// buildPlainFieldMask returns a boolean slice aligned with the analyzed
// query's non-junk TargetList. An entry is true if the column was produced
// by * expansion (SELECT * or SELECT t.*), false for explicit column refs.
func buildPlainFieldMask(selStmt *ast.SelectStmt, q *catalog.Query) []bool {
	// Collect non-junk target entries.
	var targets []*catalog.TargetEntry
	for _, te := range q.TargetList {
		if !te.ResJunk {
			targets = append(targets, te)
		}
	}
	nResults := len(targets)

	if selStmt == nil || selStmt.TargetList == nil {
		return make([]bool, nResults)
	}

	// Count parse-tree items (stars and non-stars).
	var parseTargets []*ast.ResTarget
	for _, item := range selStmt.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok {
			parseTargets = append(parseTargets, rt)
		}
	}

	// For each star, determine how many columns it expands to by
	// looking at the table it references and counting columns in that RTE.
	mask := make([]bool, nResults)
	pos := 0

	for _, rt := range parseTargets {
		if pos >= nResults {
			break
		}
		if isStarTarget(rt) {
			tableName := starTableName(rt)
			starCount := countStarColumnsFromRTE(q, targets, pos, tableName)
			for i := 0; i < starCount && pos < nResults; i++ {
				mask[pos] = true
				pos++
			}
		} else {
			mask[pos] = false
			pos++
		}
	}

	return mask
}

// isStarTarget checks if a ResTarget in the parse tree is a star expression
// (SELECT * or SELECT t.*).
func isStarTarget(rt *ast.ResTarget) bool {
	if rt == nil || rt.Val == nil {
		return false
	}
	cr, ok := rt.Val.(*ast.ColumnRef)
	if !ok || cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return false
	}
	last := cr.Fields.Items[len(cr.Fields.Items)-1]
	_, isStar := last.(*ast.A_Star)
	return isStar
}

// starTableName extracts the table qualifier from a star target.
// For "SELECT *" returns "", for "SELECT t.*" returns "t".
func starTableName(rt *ast.ResTarget) string {
	cr, ok := rt.Val.(*ast.ColumnRef)
	if !ok || cr.Fields == nil {
		return ""
	}
	// For t.*, Fields = [String("t"), A_Star{}]
	if len(cr.Fields.Items) >= 2 {
		if s, ok := cr.Fields.Items[0].(*ast.String); ok {
			return s.Str
		}
	}
	return ""
}

// countStarColumnsFromRTE counts how many analyzed target entries starting
// at pos belong to a star expansion. For unqualified * (tableName=""), it
// counts entries that share the same contiguous RangeIdx sequence covering
// all non-join RTEs. For qualified t.* (tableName="t"), it counts entries
// whose RTE ERef matches the table name.
func countStarColumnsFromRTE(q *catalog.Query, targets []*catalog.TargetEntry, startPos int, tableName string) int {
	if startPos >= len(targets) {
		return 0
	}

	if tableName != "" {
		// Qualified star: count the named RTE's columns in order, 1, 2, ...,
		// so an explicit reference to it that follows is not counted.
		count, first := 0, 0
		for i := startPos; i < len(targets); i++ {
			v, ok := targets[i].Expr.(*catalog.VarExpr)
			if !ok {
				break
			}
			if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
				break
			}
			rte := q.RangeTable[v.RangeIdx]
			if rte.ERef != tableName && rte.Alias != tableName {
				break
			}
			if count == 0 {
				first = v.RangeIdx
			}
			if v.RangeIdx != first || int(v.AttNum) != count+1 {
				break
			}
			count++
		}
		if count == 0 {
			return 1
		}
		return count
	}

	// Unqualified star: expands all non-join RTEs in FROM order.
	// Count how many consecutive VarExprs follow the expected star expansion pattern.
	count := 0
	for i := startPos; i < len(targets); i++ {
		v, ok := targets[i].Expr.(*catalog.VarExpr)
		if !ok {
			break
		}
		if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
			break
		}
		rte := q.RangeTable[v.RangeIdx]
		if rte.Kind == catalog.RTEJoin && rte.Alias == "" {
			break
		}
		// Check if this column's AttNum is sequential within its RTE.
		// Star expansion produces columns in order: attnum 1, 2, 3, ...
		if count > 0 {
			if prev, ok := targets[i-1].Expr.(*catalog.VarExpr); ok && v.RangeIdx == prev.RangeIdx && v.AttNum != prev.AttNum+1 {
				break
			}
		}
		count++
	}
	if count == 0 {
		return 1
	}
	return count
}

func (e *querySpanExtractor) extractSetOpLineageWithVisited(parents []*catalog.Query, q *catalog.Query, selStmt *ast.SelectStmt, visited map[*catalog.Query]bool) []QuerySpanResult {
	if q.LArg == nil || q.RArg == nil {
		return nil
	}
	var lSelStmt, rSelStmt *ast.SelectStmt
	if selStmt != nil {
		lSelStmt = selStmt.Larg
		rSelStmt = selStmt.Rarg
	}
	branchParents := pushQuery(parents, q)
	lResults := e.extractLineageWithVisited(branchParents, q.LArg, lSelStmt, visited)
	rResults := e.extractLineageWithVisited(branchParents, q.RArg, rSelStmt, visited)

	// For EXCEPT, output rows come only from the left branch — do not
	// merge right-side lineage into the result columns.
	includeRight := q.SetOp != catalog.SetOpExcept && q.SetOp != catalog.SetOpExceptAll

	var results []QuerySpanResult
	for i, lr := range lResults {
		merged := make(ColumnSet)
		for col := range lr.SourceColumns {
			merged[col] = true
		}
		if includeRight && i < len(rResults) {
			for col := range rResults[i].SourceColumns {
				merged[col] = true
			}
		}
		results = append(results, QuerySpanResult{
			Name:          lr.Name,
			SourceColumns: merged,
			IsPlainField:  false, // set-op results are never plain fields
		})
	}
	return results
}

// walkExpr recursively walks an analyzed expression tree and collects all
// source column references into the result set.
func (e *querySpanExtractor) walkExpr(q *catalog.Query, expr catalog.AnalyzedExpr, result ColumnSet) {
	e.walkExprWithVisited([]*catalog.Query{q}, expr, result, make(map[queryExprKey]bool))
}

func (e *querySpanExtractor) walkExprWithVisited(queryStack []*catalog.Query, expr catalog.AnalyzedExpr, result ColumnSet, visited map[queryExprKey]bool) {
	if expr == nil {
		return
	}
	q := currentQuery(queryStack)
	if q == nil {
		return
	}
	key := queryExprKey{query: q, expr: expr}
	if visited[key] {
		return
	}
	visited[key] = true

	switch v := expr.(type) {
	case *catalog.VarExpr:
		e.resolveVar(queryStack, v, result, visited)
	case *catalog.FuncCallExpr:
		// Walk arguments first — they contribute to access tables.
		for _, arg := range v.Args {
			e.walkExprWithVisited(queryStack, arg, result, visited)
		}
		// Trace lineage through user-defined function body.
		for _, proc := range e.userProcsCalled(v) {
			// Expression context: merge all output columns into one set.
			for _, colSet := range e.analyzeFunctionBody(proc) {
				for k, val := range colSet {
					result[k] = val
				}
			}
		}
	case *catalog.SubLinkExpr:
		e.walkExprWithVisited(queryStack, v.TestExpr, result, visited)
		if v.SubQuery != nil {
			subQueryStack := appendQuery(queryStack, v.SubQuery)
			for _, te := range v.SubQuery.TargetList {
				if !te.ResJunk {
					e.walkExprWithVisited(subQueryStack, te.Expr, result, visited)
				}
			}
		}
	case *catalog.JsonConstructorExprQ:
		for _, child := range catalog.AnalyzedExprChildren(expr) {
			e.walkExprWithVisited(queryStack, child, result, visited)
		}
		if v.Subquery != nil {
			subQueryStack := appendQuery(queryStack, v.Subquery)
			for _, te := range v.Subquery.TargetList {
				if !te.ResJunk {
					e.walkExprWithVisited(subQueryStack, te.Expr, result, visited)
				}
			}
		}
	case *catalog.SetOpColumnExpr:
		// q is the set operation. EXCEPT returns only rows of its left branch.
		branches := []*catalog.Query{q.LArg, q.RArg}
		if q.SetOp == catalog.SetOpExcept || q.SetOp == catalog.SetOpExceptAll {
			branches = branches[:1]
		}
		for _, branch := range branches {
			if te := outputColumn(branch, int(v.ResNo)); te != nil {
				e.walkExprWithVisited(appendQuery(queryStack, branch), te.Expr, result, visited)
			}
		}
	default:
		for _, child := range catalog.AnalyzedExprChildren(expr) {
			e.walkExprWithVisited(queryStack, child, result, visited)
		}
	}
}

func currentQuery(queryStack []*catalog.Query) *catalog.Query {
	if len(queryStack) == 0 {
		return nil
	}
	return queryStack[len(queryStack)-1]
}

// pushQuery returns queryStack with q appended, reusing its array. The caller
// passes its whole stack and is done with the result before it pushes
// another query, so no stack in use shares the slot q takes.
func pushQuery(queryStack []*catalog.Query, q *catalog.Query) []*catalog.Query {
	return append(queryStack, q)
}

func appendQuery(queryStack []*catalog.Query, q *catalog.Query) []*catalog.Query {
	nextStack := make([]*catalog.Query, 0, len(queryStack)+1)
	nextStack = append(nextStack, queryStack...)
	if q == nil {
		return nextStack
	}
	return append(nextStack, q)
}

// resolveVar resolves a VarExpr to its ultimate source column(s).
func (e *querySpanExtractor) resolveVar(queryStack []*catalog.Query, v *catalog.VarExpr, result ColumnSet, visited map[queryExprKey]bool) {
	q := currentQuery(queryStack)
	effectiveStack := queryStack
	if v.LevelsUp > 0 {
		ancestorIdx := len(queryStack) - 1 - v.LevelsUp
		if ancestorIdx < 0 {
			return
		}
		q = queryStack[ancestorIdx]
		effectiveStack = queryStack[:ancestorIdx+1]
	}
	if q == nil {
		return
	}
	if v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colIdx := int(v.AttNum - 1) // AttNum is 1-based

	switch rte.Kind {
	case catalog.RTERelation:
		rel := e.cat.GetRelationByOID(rte.RelOID)
		if rel == nil || rel.Schema == nil {
			return
		}
		// Standalone composite types are not FROM-able; a valid analysis
		// never produces such an RTE, so treat it as unresolved.
		if rel.RelKind == 'c' {
			return
		}
		// A whole-row reference (AttNum == 0, e.g. to_jsonb(t1)) reads every
		// column of the relation.
		if v.AttNum == 0 {
			if len(rel.Columns) == 0 {
				result[ColumnResource{
					Database: e.defaultDatabase,
					Schema:   rel.Schema.Name,
					Table:    rel.Name,
				}] = true
				return
			}
			for _, col := range rel.Columns {
				result[ColumnResource{
					Database: e.defaultDatabase,
					Schema:   rel.Schema.Name,
					Table:    rel.Name,
					Column:   col.Name,
				}] = true
			}
			return
		}
		// The table's own column name, which an alias column list may hide.
		colName := catalog.SystemColumnName(v.AttNum)
		if colIdx >= 0 && colIdx < len(rel.Columns) {
			colName = rel.Columns[colIdx].Name
		}
		result[ColumnResource{
			Database: e.defaultDatabase,
			Schema:   rel.Schema.Name,
			Table:    rel.Name,
			Column:   colName,
		}] = true

	case catalog.RTESubquery:
		if v.AttNum == 0 && rte.Subquery != nil {
			subStack := appendQuery(effectiveStack, rte.Subquery)
			for _, te := range rte.Subquery.TargetList {
				if !te.ResJunk {
					e.walkExprWithVisited(subStack, te.Expr, result, visited)
				}
			}
			return
		}
		if te := outputColumn(rte.Subquery, int(v.AttNum)); te != nil {
			e.walkExprWithVisited(appendQuery(effectiveStack, rte.Subquery), te.Expr, result, visited)
		}

	case catalog.RTECTE:
		if cte, level := findCTEForRangeTableEntry(effectiveStack, rte); cte != nil {
			cteStack := appendQuery(effectiveStack[:level+1], cte.Query)
			if v.AttNum == 0 {
				for _, te := range cte.Query.TargetList {
					if !te.ResJunk {
						e.walkExprWithVisited(cteStack, te.Expr, result, visited)
					}
				}
				return
			}
			if te := outputColumn(cte.Query, int(v.AttNum)); te != nil {
				// The CTE's query sees the scopes of the query that defines it.
				e.walkExprWithVisited(cteStack, te.Expr, result, visited)
			}
		}

	case catalog.RTEJoin:
		if v.AttNum == 0 {
			for _, aliasVar := range rte.JoinAliasVars {
				e.walkExprWithVisited(effectiveStack, aliasVar, result, visited)
			}
			return
		}
		// Only a column that FULL JOIN USING merges is read through the join.
		if colIdx >= 0 && colIdx < len(rte.JoinAliasVars) {
			e.walkExprWithVisited(effectiveStack, rte.JoinAliasVars[colIdx], result, visited)
		}

	case catalog.RTETableFunc:
		if rte.TableFunc != nil {
			e.walkExprWithVisited(effectiveStack, rte.TableFunc.DocExpr, result, visited)
			e.walkExprWithVisited(effectiveStack, rte.TableFunc.RowExpr, result, visited)
			for _, ns := range rte.TableFunc.Namespaces {
				e.walkExprWithVisited(effectiveStack, ns.Expr, result, visited)
			}
			if v.AttNum == 0 {
				for _, col := range rte.TableFunc.Columns {
					e.walkExprWithVisited(effectiveStack, col.ColExpr, result, visited)
					e.walkExprWithVisited(effectiveStack, col.DefaultExpr, result, visited)
				}
			} else if colIdx >= 0 && colIdx < len(rte.TableFunc.Columns) {
				col := rte.TableFunc.Columns[colIdx]
				e.walkExprWithVisited(effectiveStack, col.ColExpr, result, visited)
				e.walkExprWithVisited(effectiveStack, col.DefaultExpr, result, visited)
			}
		}

	case catalog.RTEFunction:
		if v.AttNum == 0 {
			for _, expr := range rte.FuncExprs {
				fc, ok := expr.(*catalog.FuncCallExpr)
				if !ok {
					continue
				}
				for _, arg := range fc.Args {
					e.walkExprWithVisited(effectiveStack, arg, result, visited)
				}
				for _, proc := range e.userProcsCalled(fc) {
					for _, bodySet := range e.analyzeFunctionBody(proc) {
						for k, val := range bodySet {
							result[k] = val
						}
					}
				}
			}
			return
		}
		// An output reads the arguments of the call that produces it, as
		// jsonb_each(a) reads a, and for a user function also what its body
		// returns there. The body's parameters carry no lineage, so the
		// arguments stand in for them.
		if fc, localIdx := functionColumn(rte, colIdx); fc != nil {
			for _, arg := range fc.Args {
				e.walkExprWithVisited(effectiveStack, arg, result, visited)
			}
			for _, proc := range e.userProcsCalled(fc) {
				bodySets := e.analyzeFunctionBody(proc)
				if localIdx >= 0 && localIdx < len(bodySets) {
					for k, val := range bodySets[localIdx] {
						result[k] = val
					}
				}
			}
		}

	default:
		// Unknown RTE kind — skip.
	}
}

// functionColumn returns the call of a function RTE that produces its
// colIdx-th column (0-based) and the column's position among that call's
// columns. ROWS FROM and a multi-argument unnest place each call's columns
// after the previous call's; ORDINALITY's column belongs to none.
func functionColumn(rte *catalog.RangeTableEntry, colIdx int) (*catalog.FuncCallExpr, int) {
	if rte.FuncColCounts == nil {
		if len(rte.FuncExprs) != 1 {
			return nil, -1
		}
		fc, _ := rte.FuncExprs[0].(*catalog.FuncCallExpr)
		return fc, colIdx
	}
	for i, count := range rte.FuncColCounts {
		if colIdx < count {
			if i >= len(rte.FuncExprs) {
				return nil, -1
			}
			fc, _ := rte.FuncExprs[i].(*catalog.FuncCallExpr)
			return fc, colIdx
		}
		colIdx -= count
	}
	return nil, -1
}

// findCTEForRangeTableEntry returns the CTE that rte, an entry of the last
// query of queryStack, references, and the stack index of the query that
// defines it. rte.CTEIndex indexes the defining query's list, which may be an
// enclosing query's, so the name is matched from the innermost query out, as
// PostgreSQL scopes WITH names.
func findCTEForRangeTableEntry(queryStack []*catalog.Query, rte *catalog.RangeTableEntry) (*catalog.CommonTableExprQ, int) {
	if rte == nil || rte.CTEName == "" {
		return nil, -1
	}
	for i := len(queryStack) - 1; i >= 0; i-- {
		if queryStack[i] == nil {
			continue
		}
		for _, cte := range queryStack[i].CTEList {
			if cte != nil && cte.Name == rte.CTEName {
				return cte, i
			}
		}
	}
	return nil, -1
}

// extractFallbackColumns attempts to extract column names and lineage from the parse tree
// when AnalyzeSelectStmt fails. It walks the AST to find column references and resolves
// them against tables in the FROM clause.
func (e *querySpanExtractor) extractFallbackColumns(selStmt *ast.SelectStmt) []QuerySpanResult {
	if selStmt == nil {
		return nil
	}

	// Build a helper analyzer to reuse column reference extraction.
	analyzer := &plpgsqlAnalyzer{
		extractor: e,
		scope:     newVariableScope(nil),
	}
	fromTables := analyzer.collectFromTables(selStmt)

	// Build CTE map so we can trace through CTEs to real tables.
	cteMap := collectCTEDefinitions(selStmt)
	analyzer.cteMap = cteMap
	// Store cteMap on the extractor so extractColumnsFromRangeVar can resolve CTEs.
	e.fallbackCTEMap = cteMap
	e.fallbackCTEAliases = collectCTEAliases(selStmt)

	// Each target, a star among them, gives its columns in place, so the
	// results line up with the query's output by position.
	if selStmt.TargetList != nil {
		var results []QuerySpanResult
		hasExpr := false
		for _, item := range selStmt.TargetList.Items {
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				continue
			}
			if isStarTarget(rt) {
				results = append(results, e.expandStar(selStmt, rt, analyzer, fromTables)...)
				continue
			}
			hasExpr = true
			name := rt.Name
			if name == "" {
				name = figureResTargetName(rt)
			}
			// PostgreSQL uses "?column?" for unnamed result columns.
			if name == "" {
				name = "?column?"
			}
			// Extract column references from the expression for lineage.
			colSet := make(ColumnSet)
			analyzer.extractColumnRefsFromExpr(rt.Val, fromTables, colSet)
			results = append(results, QuerySpanResult{
				Name:          name,
				SourceColumns: colSet,
			})
		}
		// Also extract WHERE clause columns for predicate tracking.
		if hasExpr && selStmt.WhereClause != nil {
			whereColSet := make(ColumnSet)
			analyzer.extractColumnRefsFromExpr(selStmt.WhereClause, fromTables, whereColSet)
			for k, v := range whereColSet {
				e.predicateColumns[k] = v
			}
		}
		if len(results) > 0 {
			return results
		}
	}

	// For SELECT *, combine columns from ALL FROM clause items.
	if cols := e.fromClauseColumns(selStmt, analyzer, fromTables); len(cols) > 0 {
		return cols
	}

	return nil
}

// expandStar returns the columns the star target rt of sel expands to. A
// star whose FROM item was not found takes the columns of every FROM item,
// rather than none, which would shift the columns after it.
func (e *querySpanExtractor) expandStar(sel *ast.SelectStmt, rt *ast.ResTarget, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) []QuerySpanResult {
	if columns, found := e.starColumns(sel, rt, analyzer, fromTables); found {
		return columns
	}
	return e.fromClauseColumns(sel, analyzer, fromTables)
}

// fromClauseColumns returns the columns of every FROM item of sel, in order.
func (e *querySpanExtractor) fromClauseColumns(sel *ast.SelectStmt, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) []QuerySpanResult {
	if sel.FromClause == nil {
		return nil
	}
	var columns []QuerySpanResult
	for _, item := range sel.FromClause.Items {
		columns = append(columns, e.extractColumnsFromFromItem(item, analyzer, fromTables)...)
	}
	return columns
}

// starColumns returns the columns a star target of sel expands to: those of
// every FROM item for *, and those of the FROM item the qualifier names for
// name.*. It reports false when it found no FROM item the qualifier names.
//
// pg: src/backend/parser/parse_target.c — ExpandColumnRefStar
func (e *querySpanExtractor) starColumns(sel *ast.SelectStmt, rt *ast.ResTarget, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) ([]QuerySpanResult, bool) {
	if sel.FromClause == nil {
		return nil, false
	}
	// The names before the star: [table], [schema, table] or
	// [database, schema, table].
	var qualifiers []string
	if cr, ok := rt.Val.(*ast.ColumnRef); ok && cr.Fields != nil {
		for _, f := range cr.Fields.Items {
			if s, ok := f.(*ast.String); ok {
				qualifiers = append(qualifiers, s.Str)
			}
		}
	}
	qualifier, schema := "", ""
	if n := len(qualifiers); n > 0 {
		qualifier = qualifiers[n-1]
		if n > 1 {
			schema = qualifiers[n-2]
		}
	}
	var columns []QuerySpanResult
	for _, item := range sel.FromClause.Items {
		if qualifier == "" {
			columns = append(columns, e.extractColumnsFromFromItem(item, analyzer, fromTables)...)
			continue
		}
		if named := e.fromItemNamed(item, schema, qualifier); named != nil {
			return e.extractColumnsFromFromItem(named, analyzer, fromTables), true
		}
	}
	return columns, qualifier == ""
}

// fromItemNamed returns the FROM item, item or one an alias-less join of it
// holds, whose reference name is name, or nil. A schema, as in
// schema.table.*, matches only an unaliased table in that schema, written
// with it or found there through the search path.
//
// pg: src/backend/parser/parse_relation.c — refnameNamespaceItem
func (e *querySpanExtractor) fromItemNamed(item ast.Node, schema, name string) ast.Node {
	aliasOf := func(a *ast.Alias) string {
		if a == nil {
			return ""
		}
		return a.Aliasname
	}
	if schema != "" {
		if _, ok := item.(*ast.RangeVar); !ok {
			if j, ok := item.(*ast.JoinExpr); !ok || j.Alias != nil {
				return nil
			}
		}
	}
	switch v := item.(type) {
	case *ast.RangeVar:
		ref := aliasOf(v.Alias)
		if schema != "" {
			if ref == "" && v.Relname == name && e.relationSchema(v) == schema {
				return v
			}
			return nil
		}
		if ref == name || ref == "" && v.Relname == name {
			return v
		}
	case *ast.RangeSubselect:
		if aliasOf(v.Alias) == name {
			return v
		}
	case *ast.RangeFunction:
		// An unaliased function is named for its first function.
		// pg: src/backend/parser/parse_relation.c — addRangeTableEntryForFunction
		if ref := aliasOf(v.Alias); ref == name || ref == "" && rangeFunctionName(v) == name {
			return v
		}
	case *ast.JsonTable:
		if aliasOf(v.Alias) == name {
			return v
		}
	case *ast.JoinExpr:
		if v.Alias != nil {
			if v.Alias.Aliasname == name {
				return v
			}
			return nil
		}
		if l := e.fromItemNamed(v.Larg, schema, name); l != nil {
			return l
		}
		return e.fromItemNamed(v.Rarg, schema, name)
	}
	return nil
}

// relationSchema is the schema of the relation rv names: its own, or the
// first schema of the search path that has it.
func (e *querySpanExtractor) relationSchema(rv *ast.RangeVar) string {
	if rv.Schemaname != "" {
		return rv.Schemaname
	}
	db, err := e.getDatabaseMetadata(e.defaultDatabase)
	if err != nil || db == nil {
		return ""
	}
	schema, _ := db.SearchRelation(e.searchPath, rv.Relname)
	return schema
}

// rangeFunctionName is the name of the first function of a FROM item, the
// name an unaliased one is referred to by.
func rangeFunctionName(rf *ast.RangeFunction) string {
	if rf.Functions == nil || len(rf.Functions.Items) == 0 {
		return ""
	}
	list, ok := rf.Functions.Items[0].(*ast.List)
	if !ok || len(list.Items) == 0 {
		return ""
	}
	fc, ok := list.Items[0].(*ast.FuncCall)
	if !ok || fc.Funcname == nil || len(fc.Funcname.Items) == 0 {
		return ""
	}
	name, _ := fc.Funcname.Items[len(fc.Funcname.Items)-1].(*ast.String)
	if name == nil {
		return ""
	}
	return name.Str
}

// jsonTableColumns returns the columns of a JSON_TABLE: each column of its
// COLUMNS clause, NESTED ones in place, reading the context item, the
// PASSING arguments and its own DEFAULT values; an ORDINALITY column reads
// nothing.
//
// pg: src/backend/parser/parse_jsontable.c — transformJsonTable
func jsonTableColumns(jt *ast.JsonTable, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) []QuerySpanResult {
	sources := make(ColumnSet)
	if analyzer != nil {
		if jt.ContextItem != nil {
			analyzer.extractColumnRefsFromExpr(jt.ContextItem.RawExpr, fromTables, sources)
		}
		if jt.Passing != nil {
			for _, arg := range jt.Passing.Items {
				analyzer.extractColumnRefsFromExpr(arg, fromTables, sources)
			}
		}
	}
	var columns []QuerySpanResult
	var add func(list *ast.List)
	add = func(list *ast.List) {
		if list == nil {
			return
		}
		for _, item := range list.Items {
			col, ok := item.(*ast.JsonTableColumn)
			if !ok {
				continue
			}
			switch col.Coltype {
			case ast.JTC_NESTED:
				add(col.Columns)
			case ast.JTC_FOR_ORDINALITY:
				columns = append(columns, QuerySpanResult{Name: col.Name, SourceColumns: make(ColumnSet)})
			default:
				// A DEFAULT ... ON EMPTY or ON ERROR value is the column's
				// value too.
				colSources := maps.Clone(sources)
				if analyzer != nil {
					for _, b := range []*ast.JsonBehavior{col.OnEmpty, col.OnError} {
						if b != nil && b.Expr != nil {
							analyzer.extractColumnRefsFromExpr(b.Expr, fromTables, colSources)
						}
					}
				}
				columns = append(columns, QuerySpanResult{Name: col.Name, SourceColumns: colSources})
			}
		}
	}
	add(jt.Columns)
	if jt.Alias != nil && jt.Alias.Colnames != nil {
		for i, name := range astStrings(jt.Alias.Colnames) {
			if i < len(columns) {
				columns[i].Name = name
			}
		}
	}
	return columns
}

// figureResTargetName extracts a column name from a ResTarget's expression.
func figureResTargetName(rt *ast.ResTarget) string {
	if rt == nil || rt.Val == nil {
		return ""
	}
	// A cast or collation keeps its operand's name.
	// pg: src/backend/parser/parse_target.c — FigureColnameInternal
	switch v := rt.Val.(type) {
	case *ast.TypeCast:
		if name := figureResTargetName(&ast.ResTarget{Val: v.Arg}); name != "" {
			return name
		}
		if v.TypeName != nil && v.TypeName.Names != nil && len(v.TypeName.Names.Items) > 0 {
			if s, ok := v.TypeName.Names.Items[len(v.TypeName.Names.Items)-1].(*ast.String); ok {
				return s.Str
			}
		}
		return ""
	case *ast.CollateClause:
		return figureResTargetName(&ast.ResTarget{Val: v.Arg})
	default:
	}
	if cr, ok := rt.Val.(*ast.ColumnRef); ok && cr.Fields != nil {
		// Use the last String field — for "ia.approver_emails", return "approver_emails".
		// For unqualified "col", return "col". Skip A_Star nodes.
		name := ""
		for _, f := range cr.Fields.Items {
			if s, ok := f.(*ast.String); ok {
				name = s.Str
			}
		}
		return name
	}
	if fc, ok := rt.Val.(*ast.FuncCall); ok && fc.Funcname != nil {
		// Use the last name part — for "pg_catalog.func", return "func".
		name := ""
		for _, f := range fc.Funcname.Items {
			if s, ok := f.(*ast.String); ok {
				name = s.Str
			}
		}
		return name
	}
	// SQL/JSON constructors.
	if _, ok := rt.Val.(*ast.JsonObjectConstructor); ok {
		return "json_object"
	}
	if _, ok := rt.Val.(*ast.JsonArrayConstructor); ok {
		return "json_array"
	}
	return ""
}

// extractColumnsFromFromItem extracts column names from a FROM clause item.
// For RangeVar (tables), it looks up the catalog for column names and lineage.
// For RangeFunction, it uses function name or alias column names.
// analyzer and fromTables are optional (may be nil) and used for resolving
// function argument column refs in the fallback path.
func (e *querySpanExtractor) extractColumnsFromFromItem(item ast.Node, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) []QuerySpanResult {
	if item == nil {
		return nil
	}

	switch v := item.(type) {
	case *ast.RangeVar, *ast.RangeSubselect:
		if analyzer == nil {
			analyzer = &plpgsqlAnalyzer{extractor: e, scope: newVariableScope(nil)}
		}
		columns := analyzer.fromItemColumns(item, fromTables)
		// A subquery's column also reads its tables' columns of that name,
		// the failsafe fromItemColumn applies to a column named outright, so
		// a star never reads less than naming the column would.
		if rs, ok := v.(*ast.RangeSubselect); ok {
			if sel, ok := rs.Subquery.(*ast.SelectStmt); ok {
				for i := range columns {
					if columns[i].SourceColumns == nil {
						columns[i].SourceColumns = make(ColumnSet)
					}
					analyzer.innerColumnsNamed(sel, columns[i].Name, columns[i].SourceColumns)
				}
			}
		}
		return columns
	case *ast.RangeFunction:
		return e.extractColumnsFromRangeFunction(v, analyzer, fromTables)
	case *ast.JsonTable:
		return jsonTableColumns(v, analyzer, fromTables)
	case *ast.JoinExpr:
		columns := joinColumns(v, e.extractColumnsFromFromItem(v.Larg, analyzer, fromTables), e.extractColumnsFromFromItem(v.Rarg, analyzer, fromTables))
		// An alias column list renames the join's first columns.
		// pg: src/backend/parser/parse_clause.c — transformFromClauseItem (JoinExpr)
		if v.Alias != nil {
			for i, name := range astStrings(v.Alias.Colnames) {
				if i < len(columns) {
					columns[i].Name = name
				}
			}
		}
		return columns
	default:
		return nil
	}
}

// joinColumns returns the columns of a join of left and right: each USING
// (or NATURAL) column once, reading both sides, then the other columns of
// the left and of the right.
//
// pg: src/backend/parser/parse_clause.c — transformFromClauseItem (JoinExpr)
func joinColumns(join *ast.JoinExpr, left, right []QuerySpanResult) []QuerySpanResult {
	var using []string
	if join.IsNatural {
		for _, l := range left {
			if slices.ContainsFunc(right, func(r QuerySpanResult) bool { return r.Name == l.Name }) {
				using = append(using, l.Name)
			}
		}
	} else {
		using = astStrings(join.UsingClause)
	}
	if len(using) == 0 {
		return append(left, right...)
	}
	leftUsed, rightUsed := make([]bool, len(left)), make([]bool, len(right))
	take := func(columns []QuerySpanResult, used []bool, name string) *QuerySpanResult {
		for i := range columns {
			if !used[i] && columns[i].Name == name {
				used[i] = true
				return &columns[i]
			}
		}
		return nil
	}
	var columns []QuerySpanResult
	for _, name := range using {
		merged := QuerySpanResult{Name: name, SourceColumns: make(ColumnSet)}
		if l := take(left, leftUsed, name); l != nil {
			maps.Copy(merged.SourceColumns, l.SourceColumns)
			merged.IsPlainField = l.IsPlainField
		}
		if r := take(right, rightUsed, name); r != nil {
			maps.Copy(merged.SourceColumns, r.SourceColumns)
		}
		columns = append(columns, merged)
	}
	for i, column := range left {
		if !leftUsed[i] {
			columns = append(columns, column)
		}
	}
	for i, column := range right {
		if !rightUsed[i] {
			columns = append(columns, column)
		}
	}
	return columns
}

// extractColumnsFromRangeVar looks up a table in the catalog and returns its columns with lineage.
// If the table name matches a CTE in fallbackCTEMap, it extracts columns from the CTE definition.
func (e *querySpanExtractor) extractColumnsFromRangeVar(rv *ast.RangeVar) []QuerySpanResult {
	// Check if this is a CTE reference.
	if rv.Schemaname == "" && e.fallbackCTEMap != nil {
		if cteSel, ok := e.fallbackCTEMap[strings.ToLower(rv.Relname)]; ok {
			return e.extractColumnsFromCTE(rv.Relname, cteSel)
		}
	}

	var rel *catalog.Relation
	if rv.Schemaname != "" {
		rel = e.cat.GetRelation(rv.Schemaname, rv.Relname)
	} else {
		for _, schema := range e.searchPath {
			rel = e.cat.GetRelation(schema, rv.Relname)
			if rel != nil {
				break
			}
		}
	}
	if rel == nil {
		return nil
	}
	// Standalone composite types share the relation namespace but are not
	// FROM-able; PostgreSQL rejects them as table sources.
	if rel.RelKind == 'c' {
		return nil
	}
	var results []QuerySpanResult
	for _, col := range rel.Columns {
		colSet := make(ColumnSet)
		colSet[ColumnResource{
			Database: e.defaultDatabase,
			Schema:   rel.Schema.Name,
			Table:    rel.Name,
			Column:   col.Name,
		}] = true
		results = append(results, QuerySpanResult{
			Name:          col.Name,
			SourceColumns: colSet,
			IsPlainField:  true,
		})
	}
	return results
}

// extractColumnsFromCTE extracts columns from a CTE definition.
// For UNION CTEs, it uses the left branch (the non-recursive term) for column names
// and merges lineage from both branches.
func (e *querySpanExtractor) extractColumnsFromCTE(cteName string, cteSel *ast.SelectStmt) []QuerySpanResult {
	if cteSel == nil {
		return nil
	}

	// Build a helper analyzer.
	analyzer := &plpgsqlAnalyzer{
		extractor: e,
		scope:     newVariableScope(nil),
		cteMap:    e.fallbackCTEMap,
	}

	// For UNION (recursive or not), use the left branch for column names
	// and extract lineage from both branches.
	if cteSel.Larg != nil {
		leftResults := e.extractCTEBranchColumns(analyzer, cteSel.Larg)
		rightResults := e.extractCTEBranchColumns(analyzer, cteSel.Rarg)

		// Merge right-branch lineage into left-branch results.
		for i := range leftResults {
			if i < len(rightResults) {
				for k, v := range rightResults[i].SourceColumns {
					leftResults[i].SourceColumns[k] = v
				}
			}
		}

		// For recursive CTEs, the recursive term references CTE columns
		// (e.g., cc1*cc2). Track which CTE columns each recursive expression
		// depends on, then resolve those to source columns from the non-recursive term.
		if selectReferencesTable(cteSel.Rarg, cteName) {
			aliases := e.getCTEColumnAliases(cteName)
			if len(aliases) == 0 {
				// No explicit aliases — use the column names from leftResults.
				for _, r := range leftResults {
					aliases = append(aliases, r.Name)
				}
			}
			e.mergeRecursiveCTELineage(leftResults, cteSel.Rarg, cteName, aliases)
		}

		// Look up CTE column aliases from the WithClause.
		if aliases := e.getCTEColumnAliases(cteName); len(aliases) > 0 {
			for i := range leftResults {
				if i < len(aliases) {
					leftResults[i].Name = aliases[i]
				}
			}
		}

		return leftResults
	}

	// Simple CTE (no UNION) — extract directly.
	return e.extractCTEBranchColumns(analyzer, cteSel)
}

// extractCTEBranchColumns extracts column names and lineage from a single CTE branch.
func (*querySpanExtractor) extractCTEBranchColumns(analyzer *plpgsqlAnalyzer, sel *ast.SelectStmt) []QuerySpanResult {
	if sel == nil || sel.TargetList == nil {
		return nil
	}

	fromTables := analyzer.collectFromTables(sel)

	var results []QuerySpanResult
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		name := rt.Name
		if name == "" {
			name = figureResTargetName(rt)
		}
		colSet := make(ColumnSet)
		analyzer.extractColumnRefsFromExpr(rt.Val, fromTables, colSet)
		results = append(results, QuerySpanResult{
			Name:          name,
			SourceColumns: colSet,
		})
	}
	return results
}

// getCTEColumnAliases returns the column aliases defined on a CTE (e.g., t1(cc1, cc2, cc3)).
func (e *querySpanExtractor) getCTEColumnAliases(cteName string) []string {
	if e.fallbackCTEAliases == nil {
		return nil
	}
	return e.fallbackCTEAliases[strings.ToLower(cteName)]
}

// selectReferencesTable checks if a SELECT statement references a table name
// in its FROM clause. Used to detect recursive CTEs.
func selectReferencesTable(sel *ast.SelectStmt, tableName string) bool {
	if sel == nil || sel.FromClause == nil {
		return false
	}
	lower := strings.ToLower(tableName)
	for _, item := range sel.FromClause.Items {
		if nodeReferencesTable(item, lower) {
			return true
		}
	}
	return false
}

func nodeReferencesTable(node ast.Node, tableName string) bool {
	if node == nil {
		return false
	}
	switch v := node.(type) {
	case *ast.RangeVar:
		return strings.ToLower(v.Relname) == tableName
	case *ast.JoinExpr:
		return nodeReferencesTable(v.Larg, tableName) || nodeReferencesTable(v.Rarg, tableName)
	default:
		return false
	}
}

// mergeRecursiveCTELineage tracks which CTE columns each expression in the
// recursive term depends on, then resolves those to source columns from the
// non-recursive term results.
func (*querySpanExtractor) mergeRecursiveCTELineage(results []QuerySpanResult, recursiveSel *ast.SelectStmt, _ string, cteColNames []string) {
	if recursiveSel == nil || recursiveSel.TargetList == nil {
		return
	}

	// Build a map of CTE column name → index (for resolving references).
	cteColIndex := make(map[string]int)
	for i, name := range cteColNames {
		cteColIndex[strings.ToLower(name)] = i
	}

	// For each expression in the recursive term, find which CTE columns it references.
	for i, item := range recursiveSel.TargetList.Items {
		if i >= len(results) {
			break
		}
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}

		// Collect all column references from this expression.
		refs := collectColumnRefNames(rt.Val)

		// For each referenced column, if it's a CTE column, merge that column's
		// source set into this result's source set.
		for _, ref := range refs {
			if idx, ok := cteColIndex[strings.ToLower(ref)]; ok && idx < len(results) {
				for k, v := range results[idx].SourceColumns {
					results[i].SourceColumns[k] = v
				}
			}
		}
	}
}

// collectColumnRefNames extracts all unqualified column reference names from an expression.
func collectColumnRefNames(node ast.Node) []string {
	if node == nil {
		return nil
	}
	var names []string
	switch v := node.(type) {
	case *ast.ColumnRef:
		if v.Fields != nil && len(v.Fields.Items) == 1 {
			if s, ok := v.Fields.Items[0].(*ast.String); ok {
				names = append(names, s.Str)
			}
		}
	case *ast.FuncCall:
		if v.Args != nil {
			for _, arg := range v.Args.Items {
				names = append(names, collectColumnRefNames(arg)...)
			}
		}
	case *ast.A_Expr:
		names = append(names, collectColumnRefNames(v.Lexpr)...)
		names = append(names, collectColumnRefNames(v.Rexpr)...)
	case *ast.TypeCast:
		names = append(names, collectColumnRefNames(v.Arg)...)
	default:
	}
	return names
}

// collectCTEAliases extracts CTE column aliases from the WITH clause.
func collectCTEAliases(selStmt *ast.SelectStmt) map[string][]string {
	if selStmt == nil || selStmt.WithClause == nil {
		return nil
	}

	aliasMap := make(map[string][]string)
	for _, item := range selStmt.WithClause.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		if cte.Aliascolnames == nil || len(cte.Aliascolnames.Items) == 0 {
			continue
		}
		var aliases []string
		for _, nameItem := range cte.Aliascolnames.Items {
			if s, ok := nameItem.(*ast.String); ok {
				aliases = append(aliases, s.Str)
			}
		}
		if len(aliases) > 0 {
			aliasMap[strings.ToLower(cte.Ctename)] = aliases
		}
	}
	return aliasMap
}

// extractColumnsFromRangeFunction extracts columns from a function in FROM clause.
// It first tries to look up user-defined functions in the catalog for lineage tracing.
// analyzer and fromTables are optional and used for resolving function argument column refs.
func (e *querySpanExtractor) extractColumnsFromRangeFunction(rf *ast.RangeFunction, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource) []QuerySpanResult {
	// Try to resolve user-defined function for lineage tracing.
	funcName := extractFuncNameFromRange(rf)
	if funcName != "" && e.cat != nil {
		// First try catalog lookup.
		if proc := e.lookupUserProcByName(funcName); proc != nil {
			outNames := getOutputParamNames(proc)
			if len(outNames) > 0 {
				bodySets := e.analyzeFunctionBody(proc)
				var results []QuerySpanResult
				for i, name := range outNames {
					colSet := make(ColumnSet)
					if i < len(bodySets) {
						for k, v := range bodySets[i] {
							colSet[k] = v
						}
					}
					results = append(results, QuerySpanResult{
						Name:          name,
						SourceColumns: colSet,
					})
				}
				e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
				return results
			}
		}
		// Catalog lookup failed — try metadata (for RETURNS TABLE functions).
		dummyAccesses := make(ColumnSet)
		if metaResults := e.tryMetadataFuncLookup(funcName, dummyAccesses); metaResults != nil {
			e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, metaResults)
			return metaResults
		}
	}

	if rf.Alias != nil && rf.Alias.Colnames != nil && len(rf.Alias.Colnames.Items) > 0 {
		// Even with alias column names, resolve function arg column refs for lineage.
		results := extractColnamesFromAlias(rf.Alias)
		e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
		return results
	}
	results := extractFuncColumnNames(rf)
	e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
	// The ordinality column counts rows and reads no argument.
	if rf.Ordinality {
		results = append(results, QuerySpanResult{Name: "ordinality", SourceColumns: ColumnSet{}})
	}
	return results
}

// enrichFuncResultsWithArgLineage resolves column references from function arguments
// and adds them as source columns to all function result columns.
// For built-in functions like jsonb_each(a), the output depends on the input column a.
func (*querySpanExtractor) enrichFuncResultsWithArgLineage(rf *ast.RangeFunction, analyzer *plpgsqlAnalyzer, fromTables map[string]ColumnResource, results []QuerySpanResult) {
	if analyzer == nil || fromTables == nil || rf.Functions == nil || len(results) == 0 {
		return
	}

	// Collect source columns from all function arguments.
	argSources := make(ColumnSet)
	for _, funcItem := range rf.Functions.Items {
		funcList, ok := funcItem.(*ast.List)
		if !ok || len(funcList.Items) < 1 {
			continue
		}
		fc, ok := funcList.Items[0].(*ast.FuncCall)
		if !ok || fc.Args == nil {
			continue
		}
		for _, arg := range fc.Args.Items {
			analyzer.extractColumnRefsFromExpr(arg, fromTables, argSources)
		}
	}

	// Apply argument source columns to all result columns.
	for i := range results {
		for k, v := range argSources {
			if results[i].SourceColumns == nil {
				results[i].SourceColumns = make(ColumnSet)
			}
			results[i].SourceColumns[k] = v
		}
	}
}

// extractColnamesFromAlias extracts column names from an alias's column name list.
func extractColnamesFromAlias(alias *ast.Alias) []QuerySpanResult {
	if alias == nil || alias.Colnames == nil {
		return nil
	}
	var results []QuerySpanResult
	for _, item := range alias.Colnames.Items {
		if s, ok := item.(*ast.String); ok {
			results = append(results, QuerySpanResult{
				Name:          s.Str,
				SourceColumns: ColumnSet{},
			})
		}
	}
	return results
}

// extractFuncColumnNames names the columns of a RangeFunction whose
// functions are not user functions omni resolved: those of a function's
// column definition list, or else one column named for the function, or one
// per argument of unnest. A star expands to these, and the ordinality column
// WITH ORDINALITY adds, in place, so their number keeps the columns after it
// in line.
//
// pg: src/backend/parser/parse_clause.c — transformRangeFunction
// pg: src/backend/parser/parse_relation.c — addRangeTableEntryForFunction
func extractFuncColumnNames(rf *ast.RangeFunction) []QuerySpanResult {
	if rf.Functions == nil {
		return nil
	}

	column := func(name string) QuerySpanResult {
		return QuerySpanResult{Name: name, SourceColumns: ColumnSet{}}
	}
	var results []QuerySpanResult
	for _, funcItem := range rf.Functions.Items {
		funcList, ok := funcItem.(*ast.List)
		if !ok || len(funcList.Items) < 1 {
			continue
		}
		fc, ok := funcList.Items[0].(*ast.FuncCall)
		if !ok || fc.Funcname == nil {
			continue
		}
		// A column definition list: on the function inside ROWS FROM, or
		// on the RangeFunction otherwise.
		coldefs := rf.Coldeflist
		if len(funcList.Items) > 1 {
			coldefs, _ = funcList.Items[1].(*ast.List)
		}
		if coldefs != nil && len(coldefs.Items) > 0 {
			for _, item := range coldefs.Items {
				if def, ok := item.(*ast.ColumnDef); ok {
					results = append(results, column(def.Colname))
				}
			}
			continue
		}
		// Get function name.
		funcName := ""
		for _, nameItem := range fc.Funcname.Items {
			if s, ok := nameItem.(*ast.String); ok {
				funcName = s.Str
			}
		}
		if funcName == "" {
			continue
		}
		// unnest with several arguments is a function per argument.
		nCols := 1
		if funcName == "unnest" && fc.Args != nil && len(fc.Args.Items) > 1 {
			nCols = len(fc.Args.Items)
		}
		for range nCols {
			results = append(results, column(funcName))
		}
	}
	return results
}

// tryUserFuncTableSource handles the case where AnalyzeSelectStmt fails because
// the FROM clause contains a user-defined RETURNS TABLE function. The omni catalog
// cannot resolve these as table sources, so we look up the function manually,
// analyze its body for lineage, and construct the result.
// Returns nil if this case doesn't apply.
func (e *querySpanExtractor) tryUserFuncTableSource(selStmt *ast.SelectStmt, accessesMap ColumnSet) []QuerySpanResult {
	if selStmt == nil || selStmt.FromClause == nil {
		return nil
	}

	analyzer := &plpgsqlAnalyzer{extractor: e, scope: newVariableScope(nil)}
	fromTables := analyzer.collectFromTables(selStmt)

	// Look for a single RangeFunction in the FROM clause.
	for _, item := range selStmt.FromClause.Items {
		rf, ok := item.(*ast.RangeFunction)
		if !ok || rf.Functions == nil {
			continue
		}

		// Extract function name from the RangeFunction.
		funcName := extractFuncNameFromRange(rf)
		if funcName == "" {
			continue
		}

		// Look up the function in the catalog.
		proc := e.lookupUserProcByName(funcName)
		if proc != nil {
			// Found in catalog — use it for lineage analysis.
			outNames := getOutputParamNames(proc)
			if len(outNames) == 0 {
				continue
			}
			bodySets := e.analyzeFunctionBody(proc)
			results := e.buildFuncResults(outNames, bodySets, accessesMap)
			e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
			return results
		}

		// Catalog lookup failed — try metadata. This handles cases where
		// the metadata function name differs from the name in the DDL definition.
		if results := e.tryMetadataFuncLookup(funcName, accessesMap); results != nil {
			e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
			return results
		}
	}

	return nil
}

// extractFuncNameFromRange extracts the function name from a RangeFunction node.
func extractFuncNameFromRange(rf *ast.RangeFunction) string {
	if rf.Functions == nil {
		return ""
	}
	for _, funcItem := range rf.Functions.Items {
		funcList, ok := funcItem.(*ast.List)
		if !ok || len(funcList.Items) < 1 {
			continue
		}
		fc, ok := funcList.Items[0].(*ast.FuncCall)
		if !ok || fc.Funcname == nil {
			continue
		}
		// Use the last name part — for "schema.func", return "func".
		name := ""
		for _, nameItem := range fc.Funcname.Items {
			if s, ok := nameItem.(*ast.String); ok {
				name = s.Str
			}
		}
		return name
	}
	return ""
}

// lookupUserProcByName finds a user-defined function by name in the catalog,
// searching through schemas in the search path.
func (e *querySpanExtractor) lookupUserProcByName(name string) *catalog.UserProc {
	lowerName := strings.ToLower(name)
	for _, schemaName := range e.searchPath {
		for _, row := range e.cat.QueryPgProc(schemaName) {
			if strings.ToLower(row.ProName) == lowerName {
				return e.cat.GetUserProcByOID(row.OID)
			}
		}
	}
	return nil
}

// getOutputParamNames returns the names of output parameters (TABLE, OUT, or INOUT mode)
// for a user-defined function. Returns nil if the function has no output params.
func getOutputParamNames(proc *catalog.UserProc) []string {
	var names []string
	for i, mode := range proc.ArgModes {
		if mode == 't' || mode == 'o' || mode == 'b' {
			if i < len(proc.ArgNames) {
				names = append(names, proc.ArgNames[i])
			}
		}
	}
	return names
}

// buildFuncResults constructs QuerySpanResults from output parameter names and
// function body lineage analysis, and populates table-level access in accessesMap.
func (*querySpanExtractor) buildFuncResults(outNames []string, bodySets []ColumnSet, accessesMap ColumnSet) []QuerySpanResult {
	for _, colSet := range bodySets {
		for k, v := range colSet {
			accessesMap[ColumnResource{
				Database: k.Database,
				Schema:   k.Schema,
				Table:    k.Table,
			}] = v
		}
	}

	var results []QuerySpanResult
	for i, name := range outNames {
		sourceColumns := make(ColumnSet)
		if i < len(bodySets) {
			for k, v := range bodySets[i] {
				sourceColumns[k] = v
			}
		}
		results = append(results, QuerySpanResult{
			Name:          name,
			SourceColumns: sourceColumns,
		})
	}
	return results
}

// tryMetadataFuncLookup searches for a function by name in the database metadata
// (not the catalog). This handles cases where the catalog can't resolve the
// function — either because the metadata function name differs from the DDL name,
// or because RETURNS TABLE functions can't be loaded into the catalog.
func (e *querySpanExtractor) tryMetadataFuncLookup(funcName string, accessesMap ColumnSet) []QuerySpanResult {
	meta, err := e.getDatabaseMetadata(e.defaultDatabase)
	if err != nil || meta == nil {
		return nil
	}

	funcs := meta.SearchFunctions(e.searchPath, funcName)
	if len(funcs) == 0 {
		return nil
	}

	// Use the first matching function's definition.
	funcDef := funcs[0].Definition
	if funcDef == "" {
		return nil
	}

	// Parse the function definition to extract output param names and body.
	// We can't rely on the catalog for RETURNS TABLE functions because the
	// catalog may reject them due to type validation.
	outNames, body, lang := parseFuncDefForLineage(funcDef)
	if len(outNames) == 0 || body == "" {
		return nil
	}

	// Construct a synthetic UserProc for body analysis.
	syntheticProc := &catalog.UserProc{
		Name:     funcName,
		Language: lang,
		Body:     body,
	}
	// Set ArgModes and ArgNames for output params.
	for _, name := range outNames {
		syntheticProc.ArgModes = append(syntheticProc.ArgModes, 't')
		syntheticProc.ArgNames = append(syntheticProc.ArgNames, name)
	}
	// Also extract input params for variable substitution.
	inputNames := parseInputParamNames(funcDef)
	for _, name := range inputNames {
		syntheticProc.ArgModes = append([]byte{'i'}, syntheticProc.ArgModes...)
		syntheticProc.ArgNames = append([]string{name}, syntheticProc.ArgNames...)
	}

	// Use a temporary OID for caching.
	syntheticProc.OID = uint32(0xFFFF0000 + len(e.funcBodyCache))

	bodySets := e.analyzeFunctionBody(syntheticProc)
	return e.buildFuncResults(outNames, bodySets, accessesMap)
}

// parseFuncDefForLineage extracts output parameter names, function body, and
// language from a CREATE FUNCTION definition string.
func parseFuncDefForLineage(definition string) (outNames []string, body string, lang string) {
	upper := strings.ToUpper(definition)

	// Extract language.
	if langIdx := strings.Index(upper, "LANGUAGE"); langIdx >= 0 {
		rest := strings.TrimSpace(definition[langIdx+8:])
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			lang = strings.ToLower(strings.TrimRight(parts[0], ";"))
		}
	}

	// Extract body.
	body = extractFuncBodyFromDef(definition)

	// Extract output param names from RETURNS TABLE(name type, ...).
	tableIdx := strings.Index(upper, "RETURNS TABLE")
	if tableIdx < 0 {
		tableIdx = strings.Index(upper, "RETURNS\n    TABLE")
	}
	if tableIdx >= 0 {
		rest := definition[tableIdx:]
		parenStart := strings.Index(rest, "(")
		if parenStart >= 0 {
			depth := 0
			parenEnd := -1
			for i := parenStart; i < len(rest); i++ {
				switch rest[i] {
				case '(':
					depth++
				case ')':
					depth--
					if depth == 0 {
						parenEnd = i
					}
				default:
				}
				if parenEnd >= 0 {
					break
				}
			}
			if parenEnd > parenStart {
				paramList := rest[parenStart+1 : parenEnd]
				for _, param := range strings.Split(paramList, ",") {
					parts := strings.Fields(strings.TrimSpace(param))
					if len(parts) >= 2 {
						outNames = append(outNames, parts[0])
					}
				}
			}
		}
	}

	return outNames, body, lang
}

// parseInputParamNames extracts input parameter names from a CREATE FUNCTION
// definition. It parses the parameter list before RETURNS.
func parseInputParamNames(definition string) []string {
	upper := strings.ToUpper(definition)

	// Find the parameter list: between FUNCTION name( and ) RETURNS.
	funcIdx := strings.Index(upper, "FUNCTION")
	if funcIdx < 0 {
		return nil
	}

	rest := definition[funcIdx:]
	parenStart := strings.Index(rest, "(")
	if parenStart < 0 {
		return nil
	}

	// Find the closing paren that ends the input params.
	returnsIdx := strings.Index(upper[funcIdx+parenStart:], "RETURNS")
	if returnsIdx < 0 {
		return nil
	}

	paramSection := rest[parenStart+1 : parenStart+returnsIdx]
	// Find the last ) before RETURNS.
	lastParen := strings.LastIndex(paramSection, ")")
	if lastParen >= 0 {
		paramSection = paramSection[:lastParen]
	}

	modeKeywords := map[string]bool{
		"IN": true, "OUT": true, "INOUT": true, "VARIADIC": true,
	}
	var names []string
	for _, param := range strings.Split(paramSection, ",") {
		parts := strings.Fields(strings.TrimSpace(param))
		if len(parts) < 2 {
			continue
		}
		// If the first word is a mode keyword (IN, OUT, INOUT, VARIADIC),
		// the parameter name is the second word. Otherwise the first word is the name.
		nameIdx := 0
		if modeKeywords[strings.ToUpper(parts[0])] {
			nameIdx = 1
		}
		if nameIdx < len(parts) {
			// Skip OUT/INOUT params — we only want input params for variable scope.
			upperFirst := strings.ToUpper(parts[0])
			if upperFirst == "OUT" {
				continue
			}
			names = append(names, parts[nameIdx])
		}
	}
	return names
}
