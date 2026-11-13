package analysis

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// ColumnResource identifies a column, or a whole table when Column is empty,
// by its folded names.
type ColumnResource struct {
	Database string
	Schema   string
	Table    string
	Column   string
}

// ColumnSet is a set of columns.
type ColumnSet map[ColumnResource]bool

// QuerySpanResult is one result column of a query and the columns it reads.
type QuerySpanResult struct {
	Name          string
	SourceColumns ColumnSet
	// IsPlainField is true when the result is a bare column reference rather
	// than an expression.
	IsPlainField bool
}

// ResolvedSpan is the lineage ResolveQuerySpan returns.
type ResolvedSpan struct {
	Type QueryType
	// SourceColumns holds every table the statement names, as table-level
	// resources (Column empty), including CTE references.
	SourceColumns ColumnSet
	Results       []QuerySpanResult
}

// Table is a table as the catalog describes it.
type Table struct {
	// Name is the table's name as the catalog spells it.
	Name string
	// Columns are the column names in table order.
	Columns []string
}

// TableResolver looks up a table by its folded names, with the defaults
// already applied to an unqualified name. Its errors are returned to the
// caller of ResolveQuerySpan, wrapped in a *ResolveError.
type TableResolver func(database, schema, table string) (*Table, error)

// ResolveOptions configures ResolveQuerySpan.
type ResolveOptions struct {
	// DefaultDatabase and DefaultSchema qualify unqualified names. An empty
	// DefaultSchema means PUBLIC.
	DefaultDatabase string
	DefaultSchema   string
	// IgnoreCaseSensitive is passed to the system-table check.
	IgnoreCaseSensitive bool
	ResolveTable        TableResolver
}

// ErrMixUserSystemTables reports a statement that reads user and system
// tables together.
var ErrMixUserSystemTables = errors.New("cannot access user and system tables at the same time")

// ResolveError is an error raised while resolving a query's lineage. The
// tables the statement accesses are known by then and carried with it.
type ResolveError struct {
	SourceColumns ColumnSet
	Err           error
}

func (e *ResolveError) Error() string { return e.Err.Error() }

func (e *ResolveError) Unwrap() error { return e.Err }

// mergeColumnSet returns the union of m and n, and whether it is larger than m.
func mergeColumnSet(m, n ColumnSet) (ColumnSet, bool) {
	r := make(ColumnSet)
	for k := range m {
		r[k] = true
	}
	for k := range n {
		if _, ok := r[k]; !ok {
			r[k] = true
		}
	}
	return r, len(r) != len(m)
}

// tableSource is a relation in scope: a catalog table or a derived one.
type tableSource interface {
	tableName() string
	schemaName() string
	databaseName() string
	// querySpanResult returns the relation's columns. A derivedTable returns
	// its own slice, which callers must copy before modifying.
	querySpanResult() []QuerySpanResult
}

// derivedTable is a derived relation: a CTE, subquery, or join.
type derivedTable struct {
	Name    string
	Columns []QuerySpanResult
}

func (p *derivedTable) tableName() string                  { return p.Name }
func (*derivedTable) schemaName() string                   { return "" }
func (*derivedTable) databaseName() string                 { return "" }
func (p *derivedTable) querySpanResult() []QuerySpanResult { return p.Columns }

// catalogTable is a catalog table; each column is its own plain source.
type catalogTable struct {
	Database string
	Schema   string
	Name     string
	Columns  []string
}

func (p *catalogTable) tableName() string    { return p.Name }
func (p *catalogTable) schemaName() string   { return p.Schema }
func (p *catalogTable) databaseName() string { return p.Database }

func (p *catalogTable) querySpanResult() []QuerySpanResult {
	result := make([]QuerySpanResult, 0, len(p.Columns))
	for _, column := range p.Columns {
		sourceColumnSet := make(ColumnSet, 1)
		sourceColumnSet[ColumnResource{
			Database: p.Database,
			Schema:   p.Schema,
			Table:    p.Name,
			Column:   column,
		}] = true
		result = append(result, QuerySpanResult{
			Name:          column,
			SourceColumns: sourceColumnSet,
			IsPlainField:  true,
		})
	}
	return result
}

// ResolveQuerySpan returns the column lineage of one Snowflake statement,
// resolving every table it names through opts.ResolveTable.
//
// Unlike Extract, which works without a catalog, ResolveQuerySpan expands
// stars and checks column references against the tables the resolver returns,
// and reports the error when one cannot be resolved. It returns nil, nil when
// the statement holds no SQL, and ErrMixUserSystemTables when the statement
// reads user and system tables together. An error raised while resolving
// lineage, after the accessed tables are known, is a *ResolveError.
//
// Statements that are not queries return their QueryType and the tables they
// access, with empty Results.
func ResolveQuerySpan(statement string, opts ResolveOptions) (*ResolvedSpan, error) {
	q := newQuerySpanExtractor(opts.DefaultDatabase, opts.DefaultSchema, opts.ResolveTable, opts.IgnoreCaseSensitive)
	return q.getQuerySpan(statement)
}

// querySpanExtractor walks an omni Snowflake AST and produces a ResolvedSpan
// describing the column-level lineage of a single SELECT/set-operation query.
type querySpanExtractor struct {
	defaultDatbase string
	defaultSchema  string
	// https://docs.snowflake.com/en/sql-reference/identifiers-syntax
	ignoreCaseSensitive bool

	resolveTable TableResolver
	// Private fields.
	// ctes is used to record the common table expressions (CTEs) in the query.
	// It should be shrunk to 0 after each query span extraction.
	ctes []*derivedTable

	// tableSourcesFrom is used to record the table sources from the query.
	tableSourcesFrom []tableSource

	// statementText is the single statement being extracted; expression result
	// names slice it by node Loc.
	statementText string

	// joinMemberSources records the individual relations inside JOIN trees of
	// the current FROM scope. tableSourcesFrom holds one (unnamed, merged)
	// pseudo-table per FROM item, so qualified references like T2.B or T2.*
	// inside a join could not resolve by table name; lookups fall back to this
	// list. Scoped/truncated together with tableSourcesFrom.
	joinMemberSources []tableSource
}

func newQuerySpanExtractor(defaultDatabase, defaultSchema string, resolveTable TableResolver, ignoreCaseSensitive bool) *querySpanExtractor {
	if defaultSchema == "" {
		// Fall back to the default schema `PUBLIC`.
		// Reference: https://docs.snowflake.com/en/sql-reference/name-resolution#name-resolution-in-queries
		defaultSchema = "PUBLIC"
	}
	return &querySpanExtractor{
		defaultDatbase:      defaultDatabase,
		defaultSchema:       defaultSchema,
		ignoreCaseSensitive: ignoreCaseSensitive,
		resolveTable:        resolveTable,
	}
}

func (q *querySpanExtractor) getQuerySpan(statement string) (*ResolvedSpan, error) {
	q.statementText = statement

	file, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	if file == nil || len(file.Stmts) == 0 {
		return nil, nil
	}
	if len(file.Stmts) != 1 {
		return nil, fmt.Errorf("expected exactly 1 statement, got %d", len(file.Stmts))
	}

	node := file.Stmts[0]

	accessTables := getAccessTables(q.defaultDatbase, q.defaultSchema, node)
	// We do not support simultaneous access to the system table and the user table
	// because we do not synchronize the schema of the system table.
	// As a result, we exclude getting query span results for accessing only the system table.
	allSystems, mixed := isMixedQuery(accessTables, q.ignoreCaseSensitive)
	if mixed {
		return nil, ErrMixUserSystemTables
	}

	queryType := QueryTypeOf(node)
	// A SELECT/set-op that touches only system tables is reclassified as
	// SelectInfoSchema.
	if queryType == Select && allSystems {
		queryType = SelectInfoSchema
	}

	if queryType != Select {
		return &ResolvedSpan{
			Type:          queryType,
			SourceColumns: accessTables,
			Results:       []QuerySpanResult{},
		}, nil
	}

	// USE and SET are session commands that QueryTypeOf classifies as Select,
	// but they are not query expressions and produce no result columns.
	switch node.(type) {
	case *ast.UseStmt, *ast.SetStmt:
		return &ResolvedSpan{
			Type:          queryType,
			SourceColumns: accessTables,
			Results:       []QuerySpanResult{},
		}, nil
	default:
	}

	result, err := q.extractPseudoTableFromQueryNode(node)
	if err != nil {
		return nil, &ResolveError{SourceColumns: accessTables, Err: err}
	}

	return &ResolvedSpan{
		Type:          queryType,
		SourceColumns: accessTables,
		Results:       result.querySpanResult(),
	}, nil
}

// extractPseudoTableFromQueryNode extracts the result columns of a query
// expression node: a *ast.SelectStmt (optionally carrying a WITH clause) or a
// *ast.SetOperationStmt.
func (q *querySpanExtractor) extractPseudoTableFromQueryNode(node ast.Node) (*derivedTable, error) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		// A WITH clause on this SELECT introduces CTEs visible to the SELECT body
		// and to later CTEs; push them, then pop on the way out so sibling scopes
		// don't see them.
		if len(n.With) > 0 {
			originalCTECount := len(q.ctes)
			if err := q.extractCTEFromWith(n.With, isRecursiveWith(n.With)); err != nil {
				return nil, err
			}
			defer func() {
				q.ctes = q.ctes[:originalCTECount]
			}()
		}
		return q.extractPseudoTableFromSelectStmt(n)
	case *ast.SetOperationStmt:
		// A WITH clause preceding a set operation is attached by omni to the
		// LEFTMOST SelectStmt, but its CTEs are visible to EVERY operand. Hoist them over
		// the whole set-op extraction; the leftmost SelectStmt pushes them
		// again during its own extraction, which is harmless shadowing.
		if leftmost := leftmostSelect(n); leftmost != nil && len(leftmost.With) > 0 {
			originalCTECount := len(q.ctes)
			if err := q.extractCTEFromWith(leftmost.With, isRecursiveWith(leftmost.With)); err != nil {
				return nil, err
			}
			defer func() {
				q.ctes = q.ctes[:originalCTECount]
			}()
		}
		return q.extractPseudoTableFromSetOperation(n)
	case *ast.ResultScanStmt:
		// FAIL CLOSED: a result-pipe (->>) query reads the previous statement's
		// result set, which has no resolvable schema here; resolving the trailing
		// SELECT against $1 would produce wrong lineage. Same posture as PIVOT.
		return nil, errors.New("result-pipe (->>) statements are not supported for query span extraction yet")
	case *ast.ShowStmt:
		// Only reachable for SHOW ... ->> <query> (a plain SHOW classifies
		// SelectInfoSchema and returns before extraction). Same fail-closed
		// posture as ResultScanStmt: $1's schema is not resolvable.
		return nil, errors.New("result-pipe (->>) statements are not supported for query span extraction yet")
	default:
		return nil, fmt.Errorf("unsupported query node type %T", node)
	}
}

// leftmostSelect descends a set-operation chain's left spine to the SelectStmt
// that textually leads the statement (where omni attaches a leading WITH).
func leftmostSelect(node ast.Node) *ast.SelectStmt {
	switch n := node.(type) {
	case *ast.SelectStmt:
		return n
	case *ast.SetOperationStmt:
		return leftmostSelect(n.Left)
	default:
		return nil
	}
}

// isRecursiveWith reports whether any CTE in the WITH list carries the RECURSIVE
// flag. Snowflake applies RECURSIVE to the whole WITH list, so a single
// recursive CTE makes the list recursive.
func isRecursiveWith(ctes []*ast.CTE) bool {
	for _, cte := range ctes {
		if cte.Recursive {
			return true
		}
	}
	return false
}

func (q *querySpanExtractor) extractCTEFromWith(ctes []*ast.CTE, recursiveWith bool) error {
	for _, cte := range ctes {
		normalizedCTEName := normalizeIdentifier(cte.Name)
		isRecursiveCTE := recursiveWith || cte.Recursive || queryNodeHasSetOperation(cte.Query)
		// UNION [ALL] BY NAME merges columns by name, which the recursive
		// fixed-point merge below does not model (it merges positionally and
		// would attribute A's lineage to B for swapped columns). Route BY NAME
		// set-op bodies through the plain extraction path, which handles BY
		// NAME correctly; a genuinely self-referencing BY NAME CTE then fails
		// resolution explicitly instead of producing wrong lineage.
		if isRecursiveCTE && setOpChainHasByName(cte.Query) {
			isRecursiveCTE = false
		}
		if isRecursiveCTE {
			if err := q.extractRecursiveCTE(normalizedCTEName, cte); err != nil {
				return err
			}
			continue
		}

		pseudoTable, err := q.extractPseudoTableFromQueryNode(cte.Query)
		if err != nil {
			return fmt.Errorf("failed to extract sensitive fields of the CTE %q: %w", normalizedCTEName, err)
		}
		pseudoTable, err = applyCTEColumnList(normalizedCTEName, pseudoTable, cte.Columns)
		if err != nil {
			return fmt.Errorf("failed to extract sensitive fields of the CTE %q: %w", normalizedCTEName, err)
		}
		q.ctes = append(q.ctes, pseudoTable)
	}
	return nil
}

// queryNodeHasSetOperation reports whether a CTE body is (or wraps) a
// set-operation, which Snowflake treats as recursive-capable. omni nests
// set-ops as *ast.SetOperationStmt; a leading WITH-bearing SELECT whose body is
// a set-op is not representable here (omni hangs WITH on the SelectStmt), so a
// plain check of the top node suffices.
func queryNodeHasSetOperation(node ast.Node) bool {
	_, ok := node.(*ast.SetOperationStmt)
	return ok
}

// setOpChainHasByName reports whether any set operation in a (possibly nested)
// set-operation chain uses the BY NAME column-matching mode.
func setOpChainHasByName(node ast.Node) bool {
	setOp, ok := node.(*ast.SetOperationStmt)
	if !ok {
		return false
	}
	if setOp.ByName {
		return true
	}
	return setOpChainHasByName(setOp.Left) || setOpChainHasByName(setOp.Right)
}

func (q *querySpanExtractor) extractRecursiveCTE(cteName string, cte *ast.CTE) error {
	setOp, ok := cte.Query.(*ast.SetOperationStmt)
	if !ok {
		// A recursive CTE whose body is not a set-operation has no recursive term
		// to iterate; treat it as a plain CTE (extract its single SELECT body).
		pseudoTable, err := q.extractPseudoTableFromQueryNode(cte.Query)
		if err != nil {
			return fmt.Errorf("failed to extract sensitive fields of recursive CTE %q: %w", cteName, err)
		}
		pseudoTable, err = applyCTEColumnList(cteName, pseudoTable, cte.Columns)
		if err != nil {
			return err
		}
		q.ctes = append(q.ctes, pseudoTable)
		return nil
	}

	// Anchor term: the left branch of the set-operation. It must be resolvable
	// without the CTE itself being visible yet.
	anchorTableSource, err := q.extractPseudoTableFromQueryNode(setOp.Left)
	if err != nil {
		return fmt.Errorf("failed to extract sensitive fields of the anchor clause of recursive CTE %q: %w", cteName, err)
	}
	anchorTableSource, err = applyCTEColumnList(cteName, anchorTableSource, cte.Columns)
	if err != nil {
		return fmt.Errorf("failed to extract sensitive fields of the anchor clause of recursive CTE %q: %w", cteName, err)
	}

	q.ctes = append(q.ctes, anchorTableSource)
	// Iterate the recursive term (the right branch) until the lineage of the CTE
	// columns reaches a fixed point.
	for {
		recursivePartTableSource, err := q.extractPseudoTableFromQueryNode(setOp.Right)
		if err != nil {
			return fmt.Errorf("failed to extract sensitive fields of the recursive clause of recursive CTE %q: %w", cteName, err)
		}
		recursivePartTableSource, err = applyCTEColumnList(cteName, recursivePartTableSource, cte.Columns)
		if err != nil {
			return fmt.Errorf("failed to extract sensitive fields of the recursive clause of recursive CTE %q: %w", cteName, err)
		}

		currentCTE := q.ctes[len(q.ctes)-1]
		mergedColumns, changed, err := mergeQuerySpanResults(currentCTE.querySpanResult(), recursivePartTableSource.querySpanResult(), cteName)
		if err != nil {
			return err
		}
		q.ctes[len(q.ctes)-1] = &derivedTable{
			Name:    cteName,
			Columns: mergedColumns,
		}
		if !changed {
			break
		}
	}
	return nil
}

func applyCTEColumnList(cteName string, pseudoTable *derivedTable, columns []ast.Ident) (*derivedTable, error) {
	if pseudoTable == nil {
		pseudoTable = &derivedTable{
			Name:    cteName,
			Columns: []QuerySpanResult{},
		}
	}

	result := &derivedTable{
		Name:    cteName,
		Columns: make([]QuerySpanResult, len(pseudoTable.querySpanResult())),
	}
	copy(result.Columns, pseudoTable.querySpanResult())

	if len(columns) == 0 {
		return result, nil
	}
	if len(columns) != len(result.querySpanResult()) {
		return nil, fmt.Errorf("the number of columns in the CTE %q returns %d fields, but the column list returns %d fields", cteName, len(result.querySpanResult()), len(columns))
	}
	for i, columnName := range columns {
		result.Columns[i].Name = normalizeIdentifier(columnName)
	}
	return result, nil
}

func mergeQuerySpanResults(currentColumns, newColumns []QuerySpanResult, cteName string) ([]QuerySpanResult, bool, error) {
	if len(currentColumns) != len(newColumns) {
		return nil, false, fmt.Errorf("recursive clause returns %d fields, but anchor clause returns %d fields in recursive CTE %q", len(newColumns), len(currentColumns), cteName)
	}

	mergedColumns := make([]QuerySpanResult, len(currentColumns))
	copy(mergedColumns, currentColumns)
	changed := false
	for i := range mergedColumns {
		var hasChange bool
		mergedColumns[i].SourceColumns, hasChange = mergeColumnSet(mergedColumns[i].SourceColumns, newColumns[i].SourceColumns)
		changed = changed || hasChange
	}
	return mergedColumns, changed, nil
}

// extractPseudoTableFromSetOperation extracts result columns from a
// set-operation query (UNION/INTERSECT/EXCEPT). omni nests chained set-ops
// left-associatively, so the left branch may itself be a *ast.SetOperationStmt.
// The left branch's column names are kept; each right branch is merged into the
// left's columns — positionally for plain set operators (the executor requires
// matching column counts), or by column name for UNION [ALL] BY NAME.
func (q *querySpanExtractor) extractPseudoTableFromSetOperation(setOp *ast.SetOperationStmt) (*derivedTable, error) {
	left, err := q.extractPseudoTableFromQueryNode(setOp.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to extract the left part of the set operation: %w", err)
	}
	right, err := q.extractPseudoTableFromQueryNode(setOp.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to extract the right part of the set operation: %w", err)
	}
	if left == nil {
		return right, nil
	}
	if right == nil {
		return left, nil
	}
	if err := mergeSetOperatorColumns(left, right, setOp.ByName); err != nil {
		return nil, err
	}
	return left, nil
}

func mergeSetOperatorColumns(left, right *derivedTable, byName bool) error {
	if left == nil || right == nil {
		return nil
	}
	leftColumns := left.querySpanResult()
	rightColumns := right.querySpanResult()

	if byName {
		// UNION [ALL] BY NAME (Snowflake-specific) aligns the branches by
		// case-insensitive column name instead of by position, and the column
		// counts may differ (columns missing on one side are NULL-filled by the
		// engine). Mirror omni's own snowflake/analysis by-name merge: keep the
		// left columns in order, merging in the lineage of the same-named right
		// column, then append the right-only columns.
		rightByName := make(map[string]QuerySpanResult, len(rightColumns))
		for _, rightColumn := range rightColumns {
			rightByName[strings.ToUpper(rightColumn.Name)] = rightColumn
		}
		seen := make(map[string]bool, len(leftColumns))
		merged := make([]QuerySpanResult, 0, len(leftColumns))
		for _, leftColumn := range leftColumns {
			key := strings.ToUpper(leftColumn.Name)
			seen[key] = true
			if rightColumn, ok := rightByName[key]; ok {
				leftColumn.SourceColumns, _ = mergeColumnSet(leftColumn.SourceColumns, rightColumn.SourceColumns)
			}
			merged = append(merged, leftColumn)
		}
		for _, rightColumn := range rightColumns {
			if !seen[strings.ToUpper(rightColumn.Name)] {
				merged = append(merged, rightColumn)
			}
		}
		left.Columns = merged
		return nil
	}

	if len(leftColumns) != len(rightColumns) {
		return fmt.Errorf("the number of columns in the left part of the set operation returns %d fields, but the right part returns %d fields", len(leftColumns), len(rightColumns))
	}
	for i := range rightColumns {
		left.Columns[i].SourceColumns, _ = mergeColumnSet(leftColumns[i].SourceColumns, rightColumns[i].SourceColumns)
	}
	return nil
}

func (q *querySpanExtractor) extractPseudoTableFromSelectStmt(ctx *ast.SelectStmt) (*derivedTable, error) {
	if ctx == nil {
		return nil, nil
	}

	if len(ctx.From) > 0 {
		// Snapshot BEFORE extracting the FROM scope: extractTableSourceFromFrom
		// appends JOIN member relations to q.joinMemberSources as it walks join
		// trees, so a later snapshot would bake those into the restore point and
		// leak them into sibling/outer scopes.
		originalFromFieldsLength := len(q.tableSourcesFrom)
		originalJoinMemberLength := len(q.joinMemberSources)
		defer func() {
			q.tableSourcesFrom = q.tableSourcesFrom[:originalFromFieldsLength]
			q.joinMemberSources = q.joinMemberSources[:originalJoinMemberLength]
		}()
		tableSourceFrom, err := q.extractTableSourceFromFrom(ctx.From)
		if err != nil {
			return nil, err
		}
		if tableSourceFrom != nil {
			q.tableSourcesFrom = append(q.tableSourcesFrom, tableSourceFrom)
		}
	}

	result := &derivedTable{
		Name:    "",
		Columns: make([]QuerySpanResult, 0),
	}

	for _, target := range ctx.Targets {
		if target == nil {
			continue
		}
		if target.Star {
			// SELECT * or SELECT t.* [EXCLUDE (cols)].
			var normalizedDatabaseName, normalizedSchemaName, normalizedTableName string
			if star, ok := target.Expr.(*ast.StarExpr); ok && star != nil && star.Qualifier != nil {
				normalizedDatabaseName, normalizedSchemaName, normalizedTableName = normalizeObjectName(star.Qualifier, "", "")
			}
			left, err := q.getAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
			if err != nil {
				return nil, fmt.Errorf("failed to extract sensitive fields of the query statement: %w", err)
			}
			left = filterIlikeColumns(left, target.Ilike)
			left = filterExcludedColumns(left, target.Exclude)
			left, err = q.applyStarReplaces(left, target.Replace)
			if err != nil {
				return nil, err
			}
			left = applyStarRenames(left, target.Rename)
			result.Columns = append(result.Columns, left...)
			continue
		}

		// Non-star target: a column reference resolves to a single source column;
		// any other expression collects every column it references.
		columnName, querySpanResult, err := q.extractQuerySpanResultFromTargetExpr(target.Expr)
		if err != nil {
			return nil, err
		}
		if !target.Alias.IsEmpty() {
			columnName = normalizeIdentifier(target.Alias)
		}
		result.Columns = append(result.Columns, QuerySpanResult{
			Name:          columnName,
			SourceColumns: querySpanResult.SourceColumns,
			IsPlainField:  querySpanResult.IsPlainField,
		})
	}

	return result, nil
}

// extractQuerySpanResultFromTargetExpr resolves one non-star SELECT-list
// expression into its output name and source columns. A bare column reference
// (the IsPlainField case) is resolved against the FROM scope back to its
// physical source column; any compound expression instead unions the source
// columns of every column reference inside it (IsPlainField=false).
func (q *querySpanExtractor) extractQuerySpanResultFromTargetExpr(expr ast.Node) (string, QuerySpanResult, error) {
	if colRef, ok := expr.(*ast.ColumnRef); ok {
		normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName := normalizeColumnRef(colRef)
		querySpanResult, err := q.getField(normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName)
		if err != nil {
			return "", QuerySpanResult{}, fmt.Errorf("failed to check whether the column %q is sensitive: %w", normalizedColumnName, err)
		}
		return querySpanResult.Name, querySpanResult, nil
	}

	// Positional column reference $N (optionally qualified, d.$1). Session
	// variables ($var, Positional=false) fall through to the generic path like
	// any other opaque expression.
	if dollar, ok := expr.(*ast.DollarRef); ok && dollar.Positional {
		querySpanResult, err := q.resolvePositionalRef(dollar)
		if err != nil {
			return "", QuerySpanResult{}, err
		}
		return querySpanResult.Name, querySpanResult, nil
	}

	name := q.exprDisplayName(expr)
	sourceColumns, err := q.collectSourceColumnsFromExpr(expr)
	if err != nil {
		return "", QuerySpanResult{}, err
	}
	return name, QuerySpanResult{
		Name:          name,
		SourceColumns: sourceColumns,
		IsPlainField:  false,
	}, nil
}

// collectSourceColumnsFromExpr walks an arbitrary SELECT-list / predicate
// expression and unions the source columns of every column reference it
// contains. Column references are resolved against the current FROM scope;
// subqueries nested in the expression contribute their own source columns.
func (q *querySpanExtractor) collectSourceColumnsFromExpr(expr ast.Node) (ColumnSet, error) {
	result := make(ColumnSet)
	if expr == nil {
		return result, nil
	}

	var walkErr error
	// walk is the ast.Inspect callback. It is a named closure so the manual
	// recursions below — for sub-trees that omni's GENERATED walker does not
	// descend into because they are non-Node embedded structs — can re-enter
	// the same traversal. Collecting into a set keeps the manual recursion
	// idempotent: it stays correct (merely redundant) if the upstream walker
	// later learns to visit these children itself.
	var walk func(node ast.Node) bool
	walk = func(node ast.Node) bool {
		if node == nil || walkErr != nil {
			return false
		}
		switch n := node.(type) {
		case *ast.DollarRef:
			// Positional $N inside a compound expression (e.g. TO_DATE($1),
			// $1 + 1) contributes the N-th FROM-scope column's lineage. Session
			// variables ($var) carry no column lineage.
			if !n.Positional {
				return true
			}
			querySpanResult, err := q.resolvePositionalRef(n)
			if err != nil {
				walkErr = err
				return false
			}
			result, _ = mergeColumnSet(result, querySpanResult.SourceColumns)
		case *ast.ColumnRef:
			normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName := normalizeColumnRef(n)
			querySpanResult, err := q.getField(normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName)
			if err != nil {
				walkErr = fmt.Errorf("failed to check whether the column %q is sensitive: %w", normalizedColumnName, err)
				return false
			}
			result, _ = mergeColumnSet(result, querySpanResult.SourceColumns)
			return false
		case *ast.SubqueryExpr:
			sub, err := q.extractPseudoTableFromQueryNode(n.Query)
			if err != nil {
				walkErr = err
				return false
			}
			for _, field := range sub.querySpanResult() {
				result, _ = mergeColumnSet(result, field.SourceColumns)
			}
			return false
		case *ast.ExistsExpr:
			sub, err := q.extractPseudoTableFromQueryNode(n.Query)
			if err != nil {
				walkErr = err
				return false
			}
			for _, field := range sub.querySpanResult() {
				result, _ = mergeColumnSet(result, field.SourceColumns)
			}
			return false
		case *ast.StarExpr:
			// COUNT(*) and similar: a star with no resolvable column contributes no
			// specific source column.
			return false
		case *ast.CaseExpr:
			// omni's generated walker only visits CaseExpr.Operand and
			// CaseExpr.Else; the WHEN branches are []*WhenClause — a non-Node
			// embedded struct — so without this manual recursion
			// `CASE WHEN A > 0 THEN B END` would contribute no lineage at all.
			for _, when := range n.Whens {
				if when == nil {
					continue
				}
				ast.Inspect(when.Cond, walk)
				ast.Inspect(when.Result, walk)
			}
			// Continue the normal traversal for Operand and Else.
			return true
		case *ast.FuncCallExpr:
			// omni's generated walker only visits FuncCallExpr.Args; the
			// WITHIN GROUP (ORDER BY ...) items and the OVER (...) window
			// specification are non-Node embedded structs, so
			// `ROW_NUMBER() OVER (ORDER BY A)` would otherwise contribute no
			// lineage at all.
			for _, item := range n.OrderBy {
				if item != nil {
					ast.Inspect(item.Expr, walk)
				}
			}
			if n.Over != nil {
				for _, partition := range n.Over.PartitionBy {
					ast.Inspect(partition, walk)
				}
				for _, item := range n.Over.OrderBy {
					if item != nil {
						ast.Inspect(item.Expr, walk)
					}
				}
				if n.Over.Frame != nil {
					ast.Inspect(n.Over.Frame.Start.Offset, walk)
					ast.Inspect(n.Over.Frame.End.Offset, walk)
				}
			}
			// Continue the normal traversal for Args.
			return true
		default:
		}
		return true
	}
	ast.Inspect(expr, walk)
	if walkErr != nil {
		return nil, walkErr
	}
	return result, nil
}

// exprDisplayName picks the output column name for an unaliased non-star
// target: the trailing part of a column reference, or else the expression's
// source text by its Loc, which for a binary expression starts at the
// operator. It is empty when the Loc is unusable.
func (q *querySpanExtractor) exprDisplayName(expr ast.Node) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if len(e.Parts) > 0 {
			return normalizeIdentifier(e.Parts[len(e.Parts)-1])
		}
	case *ast.ParenExpr:
		return q.exprDisplayName(e.Expr)
	case *ast.CastExpr:
		return q.exprDisplayName(e.Expr)
	case *ast.CollateExpr:
		return q.exprDisplayName(e.Expr)
	default:
	}
	if loc := ast.NodeLoc(expr); loc.End > loc.Start && loc.Start >= 0 && loc.End <= len(q.statementText) {
		if name := strings.TrimSpace(q.statementText[loc.Start:loc.End]); utf8.ValidString(name) {
			return name
		}
	}
	return ""
}

// normalizeColumnRef splits an omni ColumnRef (1-4 parts) into its normalized
// (database, schema, table, column) components. A 1-part ref is a bare column;
// 2-part is table.column; 3-part is schema.table.column; 4-part is
// database.schema.table.column.
func normalizeColumnRef(ref *ast.ColumnRef) (database, schema, table, column string) {
	if ref == nil {
		return "", "", "", ""
	}
	parts := ref.Parts
	switch len(parts) {
	case 1:
		return "", "", "", normalizeIdentifier(parts[0])
	case 2:
		return "", "", normalizeIdentifier(parts[0]), normalizeIdentifier(parts[1])
	case 3:
		return "", normalizeIdentifier(parts[0]), normalizeIdentifier(parts[1]), normalizeIdentifier(parts[2])
	case 4:
		return normalizeIdentifier(parts[0]), normalizeIdentifier(parts[1]), normalizeIdentifier(parts[2]), normalizeIdentifier(parts[3])
	default:
		if len(parts) == 0 {
			return "", "", "", ""
		}
		return "", "", "", normalizeIdentifier(parts[len(parts)-1])
	}
}

func filterExcludedColumns(columns []QuerySpanResult, exclude []ast.Ident) []QuerySpanResult {
	if len(exclude) == 0 {
		return columns
	}

	excludedColumns := make(map[string]bool, len(exclude))
	for _, columnName := range exclude {
		excludedColumns[normalizeIdentifier(columnName)] = true
	}

	var result []QuerySpanResult
	for _, column := range columns {
		if !excludedColumns[column.Name] {
			result = append(result, column)
		}
	}
	return result
}

// extractTableSourceFromFrom resolves a FROM clause (a list of *ast.TableRef and
// *ast.JoinExpr items) into a single combined table source exposing every column
// it provides. Multiple top-level FROM items are an implicit CROSS JOIN, so
// their columns are concatenated into one pseudo-table. A single FROM item is returned as-is so
// its alias/name is preserved for qualified column resolution; two or more items
// collapse into an unnamed pseudo-table.
func (q *querySpanExtractor) extractTableSourceFromFrom(froms []ast.Node) (tableSource, error) {
	var result tableSource
	merged := false
	for _, item := range froms {
		tableSource, err := q.extractTableSourceFromItem(item)
		if err != nil {
			return nil, err
		}
		if tableSource == nil {
			continue
		}
		if result == nil {
			result = tableSource
			continue
		}
		// Comma-list FROM items merge into one unnamed scope for positional/star
		// purposes; record each named item in joinMemberSources so qualified
		// references (alias.col, table.col) keep resolving — the same fallback
		// JOIN trees use.
		if !merged {
			q.joinMemberSources = append(q.joinMemberSources, result)
			merged = true
		}
		q.joinMemberSources = append(q.joinMemberSources, tableSource)
		result = &derivedTable{
			Name:    "",
			Columns: append(result.querySpanResult(), tableSource.querySpanResult()...),
		}
	}
	return result, nil
}

func (q *querySpanExtractor) extractTableSourceFromItem(item ast.Node) (tableSource, error) {
	switch n := item.(type) {
	case *ast.TableRef:
		return q.extractTableSourceFromTableRef(n)
	case *ast.JoinExpr:
		return q.extractTableSourceFromJoin(n)
	default:
		return nil, nil
	}
}

func (q *querySpanExtractor) extractTableSourceFromJoin(join *ast.JoinExpr) (tableSource, error) {
	if join == nil {
		return nil, nil
	}
	left, err := q.extractTableSourceFromItem(join.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sensitive fields of the left part of the join: %w", err)
	}
	if left != nil {
		q.joinMemberSources = append(q.joinMemberSources, left)
	}
	right, err := q.extractTableSourceFromItem(join.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sensitive fields of the right part of the join: %w", err)
	}
	if right != nil {
		q.joinMemberSources = append(q.joinMemberSources, right)
	}

	var leftColumns, rightColumns []QuerySpanResult
	if left != nil {
		leftColumns = left.querySpanResult()
	}
	if right != nil {
		rightColumns = right.querySpanResult()
	}

	// Snowflake has 6 join types: INNER, LEFT OUTER, RIGHT OUTER, FULL OUTER,
	// CROSS, and NATURAL. NATURAL JOIN and JOIN ... USING collapse the shared
	// join keys to a single output column; getting the collapse wrong is worse
	// than wrong lineage, because a duplicated shared column shifts every later
	// positional mask by one.

	if join.Natural {
		// NATURAL JOIN: each shared column appears once, carrying the lineage of
		// BOTH sides. Left columns keep their order, then the right-only
		// columns follow.
		rightIndexByName := make(map[string]int, len(rightColumns))
		for i, rightColumn := range rightColumns {
			if _, ok := rightIndexByName[rightColumn.Name]; !ok {
				rightIndexByName[rightColumn.Name] = i
			}
		}
		sharedNames := make(map[string]bool)
		var result []QuerySpanResult
		for _, leftColumn := range leftColumns {
			if i, ok := rightIndexByName[leftColumn.Name]; ok {
				sharedNames[leftColumn.Name] = true
				leftColumn.SourceColumns, _ = mergeColumnSet(leftColumn.SourceColumns, rightColumns[i].SourceColumns)
			}
			result = append(result, leftColumn)
		}
		for _, rightColumn := range rightColumns {
			if !sharedNames[rightColumn.Name] {
				result = append(result, rightColumn)
			}
		}
		return &derivedTable{
			Name:    "",
			Columns: result,
		}, nil
	}

	if len(join.Using) > 0 {
		// JOIN ... USING (cols): per
		// https://docs.snowflake.com/en/sql-reference/constructs/join the output
		// contains the USING columns first (in the order specified, each ONCE,
		// carrying the merged lineage of both sides), then the left table
		// columns not in USING, then the right table columns not in USING.
		usingNames := make([]string, 0, len(join.Using))
		usingSet := make(map[string]bool, len(join.Using))
		for _, col := range join.Using {
			name := normalizeIdentifier(col)
			if usingSet[name] {
				// Defensive: a duplicated USING column is invalid SQL; keep one.
				continue
			}
			usingSet[name] = true
			usingNames = append(usingNames, name)
		}

		var result []QuerySpanResult
		for _, name := range usingNames {
			var merged *QuerySpanResult
			for _, leftColumn := range leftColumns {
				if leftColumn.Name == name {
					columnCopy := leftColumn
					merged = &columnCopy
					break
				}
			}
			for _, rightColumn := range rightColumns {
				if rightColumn.Name == name {
					if merged == nil {
						columnCopy := rightColumn
						merged = &columnCopy
					} else {
						merged.SourceColumns, _ = mergeColumnSet(merged.SourceColumns, rightColumn.SourceColumns)
					}
					break
				}
			}
			if merged == nil {
				// The USING column resolves on neither side (invalid SQL or an
				// unresolvable source); emit an empty-lineage column to keep the
				// positional column count stable.
				merged = &QuerySpanResult{
					Name:          name,
					SourceColumns: make(ColumnSet),
					IsPlainField:  false,
				}
			}
			result = append(result, *merged)
		}
		for _, leftColumn := range leftColumns {
			if !usingSet[leftColumn.Name] {
				result = append(result, leftColumn)
			}
		}
		for _, rightColumn := range rightColumns {
			if !usingSet[rightColumn.Name] {
				result = append(result, rightColumn)
			}
		}
		return &derivedTable{
			Name:    "",
			Columns: result,
		}, nil
	}

	var result []QuerySpanResult
	result = append(result, leftColumns...)
	result = append(result, rightColumns...)
	return &derivedTable{
		Name:    "",
		Columns: result,
	}, nil
}

func (q *querySpanExtractor) extractTableSourceFromTableRef(ref *ast.TableRef) (tableSource, error) {
	if ref == nil {
		return nil, nil
	}

	// FAIL CLOSED on PIVOT/UNPIVOT: the pivoted relation's projection (dynamic
	// pivot columns) is not modeled here yet, so resolving the bare base table
	// would silently produce wrong lineage/positions for masking. Returning an
	// error makes the caller degrade explicitly instead.
	if ref.Pivot != nil || ref.Unpivot != nil || ref.Nested != nil {
		return nil, errors.New("PIVOT/UNPIVOT table sources are not supported for query span extraction yet")
	}

	var result []QuerySpanResult
	// underlying keeps the resolved (named) table source so an unaliased
	// reference retains its identity — qualified references (T2.B, T2.*) must
	// resolve by the original table name.
	var underlying tableSource

	switch {
	case ref.Name != nil:
		_, tableSource, err := q.findTableSchema(ref.Name, q.defaultDatbase, q.defaultSchema)
		if err != nil {
			return nil, err
		}
		underlying = tableSource
		result = append(result, tableSource.querySpanResult()...)
	case ref.Subquery != nil:
		tableSource, err := q.extractPseudoTableFromQueryNode(ref.Subquery)
		if err != nil {
			return nil, fmt.Errorf("failed to extract sensitive fields of subquery in FROM: %w", err)
		}
		result = append(result, tableSource.querySpanResult()...)
	case ref.FuncCall != nil:
		// NOTE(zp): In data-warehouse, defining a table function that returns
		// multiple rows is widespread; we should parse the function definition to
		// extract the sensitive fields. For now, a table function exposes no known
		// columns.
		result = nil
	default:
		result = nil
	}

	// If the table reference carries an alias, expose its columns under the alias
	// name so qualified column references (alias.column) resolve to it.
	if !ref.Alias.IsEmpty() {
		aliasName := normalizeIdentifier(ref.Alias)
		return &derivedTable{
			Name:    aliasName,
			Columns: result,
		}, nil
	}

	// No alias: keep the resolved table source (with its real name) so
	// qualified lookups by the original table name still match.
	if underlying != nil {
		return underlying, nil
	}

	return &derivedTable{
		Name:    "",
		Columns: result,
	}, nil
}

func (q *querySpanExtractor) findTableSchema(objectName *ast.ObjectName, normalizedFallbackDatabaseName, normalizedFallbackSchemaName string) (string, tableSource, error) {
	normalizedDatabaseName, normalizedSchemaName, normalizedTableName := normalizeObjectName(objectName, "", "")
	// Scan CTEs newest-first so an inner scope's CTE shadows an outer one with
	// the same name (q.ctes is a stack of scopes; duplicate names within one
	// WITH list are a Snowflake error, so this only affects cross-scope shadowing).
	if normalizedDatabaseName == "" && normalizedSchemaName == "" {
		for i := len(q.ctes) - 1; i >= 0; i-- {
			cte := q.ctes[i]
			if normalizedTableName == cte.Name {
				return normalizedDatabaseName, cte, nil
			}
		}
	}
	normalizedDatabaseName, normalizedSchemaName, normalizedTableName = normalizeObjectName(objectName, normalizedFallbackDatabaseName, normalizedFallbackSchemaName)
	table, err := q.resolveTable(normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
	if err != nil {
		return "", nil, err
	}
	return normalizedDatabaseName, &catalogTable{
		Name:     table.Name,
		Database: normalizedDatabaseName,
		Schema:   normalizedSchemaName,
		Columns:  table.Columns,
	}, nil
}

func (q *querySpanExtractor) getAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName string) ([]QuerySpanResult, error) {
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

	var result []QuerySpanResult

	if mask&maskDatabaseName == 0 && mask&maskSchemaName == 0 && mask&maskTableName != 0 {
		// Newest-first: inner-scope CTEs shadow outer ones (see findTableSchema).
		for i := len(q.ctes) - 1; i >= 0; i-- {
			cte := q.ctes[i]
			if normalizedTableName == cte.Name {
				result = append(result, cte.querySpanResult()...)
				return result, nil
			}
		}
	}

	for _, tableSource := range q.tableSourcesFrom {
		if mask&maskDatabaseName != 0 && normalizedDatabaseName != tableSource.databaseName() {
			continue
		}
		if mask&maskSchemaName != 0 && normalizedSchemaName != tableSource.schemaName() {
			continue
		}
		if mask&maskTableName != 0 && normalizedTableName != tableSource.tableName() {
			continue
		}
		result = append(result, tableSource.querySpanResult()...)
		return result, nil
	}

	// Qualified reference that did not match a FROM item: it may name a
	// relation inside a JOIN tree (tableSourcesFrom holds only the merged,
	// unnamed join pseudo-table).
	if mask&maskTableName != 0 {
		for _, tableSource := range q.joinMemberSources {
			if mask&maskDatabaseName != 0 && normalizedDatabaseName != tableSource.databaseName() {
				continue
			}
			if mask&maskSchemaName != 0 && normalizedSchemaName != tableSource.schemaName() {
				continue
			}
			if normalizedTableName != tableSource.tableName() {
				continue
			}
			result = append(result, tableSource.querySpanResult()...)
			return result, nil
		}
	}

	return nil, fmt.Errorf(`no matching table %q.%q.%q`, normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
}

// getField iterates through the tableSourcesFrom sequentially until we find the first matching object and return the column name, and returns the fieldInfo.
func (q *querySpanExtractor) getField(normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName string) (QuerySpanResult, error) {
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

	// We just need to iterate through the tableSourcesFrom sequentially until we find the first matching object.

	// It is safe if there are two or more objects in the tableSourcesFrom have the same column name, because the executor
	// will throw a compilation error if the column name is ambiguous.
	// For example, there are two tables T1 and T2, and both of them have a column named "C1". The following query will throw
	// a compilation error:
	// SELECT C1 FROM T1, T2;
	//
	// But users can specify the table name to avoid the compilation error:
	// SELECT T1.C1 FROM T1, T2;
	//
	// Further more, users can not use the original table name if they specify the alias name:
	// SELECT T1.C1 FROM T1 AS T3, T2; -- invalid identifier 'ADDRESS.ID'

	for _, tableSource := range q.tableSourcesFrom {
		if mask&maskDatabaseName != 0 && normalizedDatabaseName != tableSource.databaseName() {
			continue
		}
		if mask&maskSchemaName != 0 && normalizedSchemaName != tableSource.schemaName() {
			continue
		}
		if mask&maskTableName != 0 && normalizedTableName != tableSource.tableName() {
			continue
		}
		for _, field := range tableSource.querySpanResult() {
			if mask&maskColumnName != 0 && normalizedColumnName != field.Name {
				continue
			}
			return field, nil
		}
	}
	// Qualified reference into a JOIN tree: tableSourcesFrom holds only the
	// merged, unnamed join pseudo-table, so match the member relations.
	if mask&maskTableName != 0 {
		for _, tableSource := range q.joinMemberSources {
			if mask&maskDatabaseName != 0 && normalizedDatabaseName != tableSource.databaseName() {
				continue
			}
			if mask&maskSchemaName != 0 && normalizedSchemaName != tableSource.schemaName() {
				continue
			}
			if normalizedTableName != tableSource.tableName() {
				continue
			}
			for _, field := range tableSource.querySpanResult() {
				if mask&maskColumnName != 0 && normalizedColumnName != field.Name {
					continue
				}
				return field, nil
			}
		}
	}

	return QuerySpanResult{}, fmt.Errorf(`no matching column %q.%q.%q.%q`, normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName)
}

// getAccessTables extracts the list of resources from the SELECT statement, and
// normalizes the object names with the NON-EMPTY currentNormalizedDatabase and
// currentNormalizedSchema. It collects every *ast.TableRef with a table name —
// including CTE references and subquery tables.
//
// It descends into a SelectStmt's parts by hand rather than trusting
// ast.Inspect alone, so it does not depend on the generated walker's coverage.
func getAccessTables(currentNormalizedDatabase, currentNormalizedSchema string, node ast.Node) ColumnSet {
	resourceMap := make(ColumnSet)
	collectAccessTables(node, currentNormalizedDatabase, currentNormalizedSchema, resourceMap)
	return resourceMap
}

// collectAccessTables recursively records every table reference reachable from
// node into resourceMap. It explicitly descends into the SelectStmt sub-parts
// (WITH/targets/GROUP BY/ORDER BY/FETCH; redundant with the fixed omni walker,
// see getAccessTables); for everything else it falls back to ast.Inspect.
func collectAccessTables(node ast.Node, defaultDatabase, defaultSchema string, resourceMap ColumnSet) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.File:
		for _, stmt := range n.Stmts {
			collectAccessTables(stmt, defaultDatabase, defaultSchema, resourceMap)
		}
		return
	case *ast.SetOperationStmt:
		collectAccessTables(n.Left, defaultDatabase, defaultSchema, resourceMap)
		collectAccessTables(n.Right, defaultDatabase, defaultSchema, resourceMap)
		return
	case *ast.SelectStmt:
		for _, cte := range n.With {
			collectAccessTables(cte.Query, defaultDatabase, defaultSchema, resourceMap)
		}
		for _, item := range n.From {
			collectAccessTables(item, defaultDatabase, defaultSchema, resourceMap)
		}
		for _, target := range n.Targets {
			if target != nil {
				collectAccessTables(target.Expr, defaultDatabase, defaultSchema, resourceMap)
			}
		}
		collectAccessTables(n.Where, defaultDatabase, defaultSchema, resourceMap)
		if n.GroupBy != nil {
			for _, item := range n.GroupBy.Items {
				collectAccessTables(item, defaultDatabase, defaultSchema, resourceMap)
			}
		}
		collectAccessTables(n.Having, defaultDatabase, defaultSchema, resourceMap)
		collectAccessTables(n.Qualify, defaultDatabase, defaultSchema, resourceMap)
		for _, item := range n.OrderBy {
			if item != nil {
				collectAccessTables(item.Expr, defaultDatabase, defaultSchema, resourceMap)
			}
		}
		return
	case *ast.JoinExpr:
		collectAccessTables(n.Left, defaultDatabase, defaultSchema, resourceMap)
		collectAccessTables(n.Right, defaultDatabase, defaultSchema, resourceMap)
		collectAccessTables(n.On, defaultDatabase, defaultSchema, resourceMap)
		return
	case *ast.TableRef:
		if n.Name != nil {
			database, schema, table := normalizeObjectName(n.Name, defaultDatabase, defaultSchema)
			resourceMap[ColumnResource{
				Database: database,
				Schema:   schema,
				Table:    table,
			}] = true
		}
		if n.Subquery != nil {
			collectAccessTables(n.Subquery, defaultDatabase, defaultSchema, resourceMap)
		}
		if n.FuncCall != nil {
			collectAccessTables(n.FuncCall, defaultDatabase, defaultSchema, resourceMap)
		}
		return
	case *ast.SubqueryExpr:
		collectAccessTables(n.Query, defaultDatabase, defaultSchema, resourceMap)
		return
	case *ast.ExistsExpr:
		collectAccessTables(n.Query, defaultDatabase, defaultSchema, resourceMap)
		return
	}

	// For any other expression node, walk its children with ast.Inspect (which
	// fully traverses non-SelectStmt sub-trees) and recurse into the
	// query-bearing nodes it surfaces. This catches table refs nested inside
	// arbitrary expressions (e.g. a scalar subquery inside a CASE in WHERE).
	ast.Inspect(node, func(child ast.Node) bool {
		if child == nil || child == node {
			return true
		}
		switch child.(type) {
		case *ast.SelectStmt, *ast.SetOperationStmt, *ast.TableRef, *ast.JoinExpr, *ast.SubqueryExpr, *ast.ExistsExpr:
			collectAccessTables(child, defaultDatabase, defaultSchema, resourceMap)
			return false
		}
		return true
	})
}

// isMixedQuery checks whether the query accesses the user table and system table at the same time.
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

// isSystemResource recognizes no system table yet, so ErrMixUserSystemTables
// never fires and no query is reclassified as SelectInfoSchema.
func isSystemResource(ColumnResource, bool) bool {
	return false
}

// resolvePositionalRef resolves a positional column reference $N (optionally
// qualified, d.$1) to the N-th column of the FROM scope (or of the qualified
// relation): errors on a non-numeric, < 1, or out-of-range position.
func (q *querySpanExtractor) resolvePositionalRef(dollar *ast.DollarRef) (QuerySpanResult, error) {
	position, err := strconv.Atoi(dollar.Name)
	if err != nil {
		return QuerySpanResult{}, fmt.Errorf("failed to parse column position %q to integer: %w", dollar.Name, err)
	}
	if position < 1 {
		return QuerySpanResult{}, fmt.Errorf("column position %d is invalid because it is less than 1", position)
	}
	var normalizedDatabaseName, normalizedSchemaName, normalizedTableName string
	if dollar.Qualifier != nil {
		normalizedDatabaseName, normalizedSchemaName, normalizedTableName = normalizeObjectName(dollar.Qualifier, "", "")
	}
	left, err := q.getAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
	if err != nil {
		return QuerySpanResult{}, fmt.Errorf("failed to resolve positional column reference $%s: %w", dollar.Name, err)
	}
	if position > len(left) {
		return QuerySpanResult{}, fmt.Errorf("column position $%d is invalid: the FROM clause only returns %d columns", position, len(left))
	}
	return left[position-1], nil
}

// applyStarRenames applies the `SELECT * RENAME (col AS alias, ...)` transform
// to a star expansion: the matching result columns keep their lineage but take
// the alias as their output name (Snowflake returns the renamed column).
// Matching uses normalized identifiers, like the rest of the resolver.
func applyStarRenames(columns []QuerySpanResult, renames []ast.StarRename) []QuerySpanResult {
	if len(renames) == 0 {
		return columns
	}
	aliasByColumn := make(map[string]string, len(renames))
	for _, r := range renames {
		aliasByColumn[normalizeIdentifier(r.Col)] = normalizeIdentifier(r.Alias)
	}
	// Copy before renaming: columns may be the FROM-scope table source's
	// internal slice (getAllFieldsOfTableInFromOrOuterCTE returns it directly,
	// and pseudoTable documents that callers must copy before modifying);
	// mutating it in place would corrupt the scope for later references.
	renamed := make([]QuerySpanResult, len(columns))
	copy(renamed, columns)
	for i := range renamed {
		if alias, ok := aliasByColumn[renamed[i].Name]; ok {
			renamed[i].Name = alias
		}
	}
	return renamed
}

// filterIlikeColumns applies the `SELECT * ILIKE '<pattern>'` transform: only
// columns whose name matches the pattern (case-insensitive; % = any run,
// _ = any single character) survive the star expansion. Applied FIRST, before
// EXCLUDE/REPLACE/RENAME, per the documented transform order. The positional
// masker depends on the span matching the actual result-set shape, so this
// must genuinely filter (not over-attribute).
func filterIlikeColumns(columns []QuerySpanResult, pattern *ast.Literal) []QuerySpanResult {
	if pattern == nil {
		return columns
	}
	matcher := ilikeMatcher(pattern.Value)
	filtered := make([]QuerySpanResult, 0, len(columns))
	for _, column := range columns {
		if matcher(column.Name) {
			filtered = append(filtered, column)
		}
	}
	return filtered
}

// ilikeMatcher compiles a Snowflake ILIKE pattern into a case-insensitive
// matcher over the whole name (% -> any run, _ -> any single rune). Go regexp
// (?i) covers ASCII case-insensitivity; it does not replicate Snowflake SQL
// collation for non-ASCII identifiers (low risk: column names are near-always
// ASCII).
func ilikeMatcher(pattern string) func(string) bool {
	var sb strings.Builder
	sb.WriteString(`(?is)\A`)
	for _, r := range pattern {
		switch r {
		case '%':
			sb.WriteString(`.*`)
		case '_':
			sb.WriteString(`.`)
		default:
			sb.WriteString(regexp.QuoteMeta(string(r)))
		}
	}
	sb.WriteString(`\z`)
	re, err := regexp.Compile(sb.String())
	if err != nil {
		// A pattern that fails to compile matches nothing (fail closed).
		return func(string) bool { return false }
	}
	return re.MatchString
}

// applyStarReplaces applies the `SELECT * REPLACE (expr AS col, ...)` transform
// to a star expansion: the matching result column keeps its position and name
// but its lineage becomes the replacement EXPRESSION's source columns
// (Snowflake returns expr in place of the original column value). Copy-on-write
// like applyStarRenames — the input may be a table source's internal slice.
func (q *querySpanExtractor) applyStarReplaces(columns []QuerySpanResult, replaces []ast.StarReplace) ([]QuerySpanResult, error) {
	if len(replaces) == 0 {
		return columns, nil
	}
	exprByColumn := make(map[string]ast.Node, len(replaces))
	for _, r := range replaces {
		exprByColumn[normalizeIdentifier(r.Col)] = r.Expr
	}
	replaced := make([]QuerySpanResult, len(columns))
	copy(replaced, columns)
	for i := range replaced {
		expr, ok := exprByColumn[replaced[i].Name]
		if !ok {
			continue
		}
		sourceColumns, err := q.collectSourceColumnsFromExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve star REPLACE expression for column %q: %w", replaced[i].Name, err)
		}
		replaced[i].SourceColumns = sourceColumns
		replaced[i].IsPlainField = false
	}
	return replaced, nil
}
