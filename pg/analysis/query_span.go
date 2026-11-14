package analysis

import (
	"context"
	"fmt"
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
// Accesses are conservative: CTE names can match physical tables, so this may
// refuse a CTE-only query. Reads missing from the access set remain unchecked.
func (e *querySpanExtractor) unresolvedColumnsError(accesses ColumnSet) *UnresolvedColumnsError {
	seen := make(map[ColumnResource]bool, len(accesses))
	var unresolved []ColumnResource
	for access := range accesses {
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
	accessTables, err := e.extractAccessTables(stmt)
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
	return e.cat.AnalyzeSelectStmtWithOptions(sel, catalog.AnalyzeOptions{OpaqueUnknownFunctions: true})
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
	return e.extractLineageWithVisited(q, selStmt, make(map[*catalog.Query]bool))
}

func (e *querySpanExtractor) extractLineageWithVisited(q *catalog.Query, selStmt *ast.SelectStmt, visited map[*catalog.Query]bool) []QuerySpanResult {
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
		return e.extractSetOpLineageWithVisited(q, selStmt, visited)
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
		e.walkExpr(q, te.Expr, sourceColSet)
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
		if rte.CTEIndex >= 0 && rte.CTEIndex < len(q.CTEList) {
			cte := q.CTEList[rte.CTEIndex]
			if te := outputColumn(cte.Query, int(v.AttNum)); te != nil {
				return isUltimatelyPlainColumnWithVisited(cte.Query, te.Expr, visited)
			}
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
		// Qualified star: count consecutive entries from the named RTE.
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
			if rte.ERef != tableName && rte.Alias != tableName {
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
		if rte.Kind == catalog.RTEJoin {
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

func (e *querySpanExtractor) extractSetOpLineageWithVisited(q *catalog.Query, selStmt *ast.SelectStmt, visited map[*catalog.Query]bool) []QuerySpanResult {
	if q.LArg == nil || q.RArg == nil {
		return nil
	}
	var lSelStmt, rSelStmt *ast.SelectStmt
	if selStmt != nil {
		lSelStmt = selStmt.Larg
		rSelStmt = selStmt.Rarg
	}
	lResults := e.extractLineageWithVisited(q.LArg, lSelStmt, visited)
	rResults := e.extractLineageWithVisited(q.RArg, rSelStmt, visited)

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
		colName := ""
		if colIdx >= 0 && colIdx < len(rte.ColNames) {
			colName = rte.ColNames[colIdx]
		}
		result[ColumnResource{
			Database: e.defaultDatabase,
			Schema:   rel.Schema.Name,
			Table:    rel.Name,
			Column:   colName,
		}] = true

	case catalog.RTESubquery:
		if te := outputColumn(rte.Subquery, int(v.AttNum)); te != nil {
			e.walkExprWithVisited(appendQuery(effectiveStack, rte.Subquery), te.Expr, result, visited)
		}

	case catalog.RTECTE:
		if cte := findCTEForRangeTableEntry(effectiveStack, rte); cte != nil {
			if te := outputColumn(cte.Query, int(v.AttNum)); te != nil {
				e.walkExprWithVisited(appendQuery(effectiveStack, cte.Query), te.Expr, result, visited)
			}
		}

	case catalog.RTEJoin:
		// Only a column that FULL JOIN USING merges is read through the join.
		if colIdx >= 0 && colIdx < len(rte.JoinAliasVars) {
			e.walkExprWithVisited(effectiveStack, rte.JoinAliasVars[colIdx], result, visited)
		}

	case catalog.RTEFunction:
		// Trace lineage through user-defined function body.
		if len(rte.FuncExprs) > 0 {
			if fc, ok := rte.FuncExprs[0].(*catalog.FuncCallExpr); ok {
				if procs := e.userProcsCalled(fc); len(procs) > 0 {
					for _, proc := range procs {
						bodySets := e.analyzeFunctionBody(proc)
						if colIdx >= 0 && colIdx < len(bodySets) {
							for k, val := range bodySets[colIdx] {
								result[k] = val
							}
						}
					}
				} else {
					// Built-in function (no user proc) — trace lineage through
					// the function's arguments. E.g., jsonb_each(a) depends on column a.
					for _, arg := range fc.Args {
						e.walkExprWithVisited(effectiveStack, arg, result, visited)
					}
				}
			}
		}

	default:
		// Unknown RTE kind — skip.
	}
}

func findCTEForRangeTableEntry(queryStack []*catalog.Query, rte *catalog.RangeTableEntry) *catalog.CommonTableExprQ {
	if len(queryStack) == 0 || rte == nil {
		return nil
	}

	current := currentQuery(queryStack)
	if current != nil && rte.CTEIndex >= 0 && rte.CTEIndex < len(current.CTEList) {
		return current.CTEList[rte.CTEIndex]
	}

	if rte.CTEName == "" {
		return nil
	}
	for i := len(queryStack) - 1; i >= 0; i-- {
		for _, cte := range queryStack[i].CTEList {
			if strings.EqualFold(cte.Name, rte.CTEName) {
				return cte
			}
		}
	}
	return nil
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

	// Check if SELECT has explicit target list (not *).
	if selStmt.TargetList != nil {
		var results []QuerySpanResult
		allStar := true
		for _, item := range selStmt.TargetList.Items {
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				continue
			}
			if !isStarTarget(rt) {
				allStar = false
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
		}
		if !allStar && len(results) > 0 {
			// Also extract WHERE clause columns for predicate tracking.
			if selStmt.WhereClause != nil {
				whereColSet := make(ColumnSet)
				analyzer.extractColumnRefsFromExpr(selStmt.WhereClause, fromTables, whereColSet)
				for k, v := range whereColSet {
					e.predicateColumns[k] = v
				}
			}
			return results
		}
	}

	// For SELECT *, combine columns from ALL FROM clause items.
	if selStmt.FromClause != nil {
		var allCols []QuerySpanResult
		for _, item := range selStmt.FromClause.Items {
			allCols = append(allCols, e.extractColumnsFromFromItem(item, analyzer, fromTables)...)
		}
		if len(allCols) > 0 {
			return allCols
		}
	}

	return nil
}

// figureResTargetName extracts a column name from a ResTarget's expression.
func figureResTargetName(rt *ast.ResTarget) string {
	if rt == nil || rt.Val == nil {
		return ""
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
	case *ast.RangeVar:
		return e.extractColumnsFromRangeVar(v)
	case *ast.RangeFunction:
		return e.extractColumnsFromRangeFunction(v, analyzer, fromTables)
	case *ast.RangeSubselect:
		if v.Alias != nil && v.Alias.Colnames != nil && len(v.Alias.Colnames.Items) > 0 {
			return extractColnamesFromAlias(v.Alias)
		}
		// No explicit column aliases — extract column names from subquery's target list
		// (only when the target list has explicit columns, not SELECT *).
		if subSel, ok := v.Subquery.(*ast.SelectStmt); ok && subSel.TargetList != nil {
			var results []QuerySpanResult
			hasStar := false
			for _, item := range subSel.TargetList.Items {
				rt, ok := item.(*ast.ResTarget)
				if !ok {
					continue
				}
				if isStarTarget(rt) {
					hasStar = true
					break
				}
				name := rt.Name
				if name == "" {
					name = figureResTargetName(rt)
				}
				if name == "" {
					name = "?column?"
				}
				results = append(results, QuerySpanResult{
					Name:          name,
					SourceColumns: ColumnSet{},
				})
			}
			if !hasStar && len(results) > 0 {
				return results
			}
			// For SELECT * in subquery, try to resolve columns from the subquery's FROM tables.
			if hasStar && subSel.FromClause != nil {
				var starResults []QuerySpanResult
				for _, fromItem := range subSel.FromClause.Items {
					starResults = append(starResults, e.extractColumnsFromFromItem(fromItem, analyzer, fromTables)...)
				}
				if len(starResults) > 0 {
					return starResults
				}
			}
		}
		return nil
	case *ast.JoinExpr:
		var results []QuerySpanResult
		results = append(results, e.extractColumnsFromFromItem(v.Larg, analyzer, fromTables)...)
		results = append(results, e.extractColumnsFromFromItem(v.Rarg, analyzer, fromTables)...)
		return results
	default:
		return nil
	}
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
				return results
			}
		}
		// Catalog lookup failed — try metadata (for RETURNS TABLE functions).
		dummyAccesses := make(ColumnSet)
		if metaResults := e.tryMetadataFuncLookup(funcName, dummyAccesses); metaResults != nil {
			return metaResults
		}
	}

	if rf.Alias != nil && rf.Alias.Colnames != nil && len(rf.Alias.Colnames.Items) > 0 {
		// Even with alias column names, resolve function arg column refs for lineage.
		results := extractColnamesFromAlias(rf.Alias)
		e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
		return results
	}
	results := extractFuncColumnNames(rf, rf.Alias)
	e.enrichFuncResultsWithArgLineage(rf, analyzer, fromTables, results)
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

// extractFuncColumnNames extracts column names from a RangeFunction.
// For functions like unnest(ARRAY, ARRAY), each argument produces a column
// named after the function.
func extractFuncColumnNames(rf *ast.RangeFunction, _ *ast.Alias) []QuerySpanResult {
	if rf.Functions == nil {
		return nil
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
		// Determine number of columns: for set-returning functions like unnest,
		// each argument produces a column.
		nCols := 1
		if fc.Args != nil && len(fc.Args.Items) > 1 {
			nCols = len(fc.Args.Items)
		}
		for range nCols {
			results = append(results, QuerySpanResult{
				Name:          funcName,
				SourceColumns: ColumnSet{},
			})
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
			return e.buildFuncResults(outNames, bodySets, accessesMap)
		}

		// Catalog lookup failed — try metadata. This handles cases where
		// the metadata function name differs from the name in the DDL definition.
		if results := e.tryMetadataFuncLookup(funcName, accessesMap); results != nil {
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
