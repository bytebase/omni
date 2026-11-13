package v1review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereDisallowFunctionsAndCalculations reports each WHERE-clause
// function call or calculation applied to an indexed column, judged against
// the pre-change schema.
//
// Judged: the WHERE of top-level SELECT, UPDATE and DELETE and of MERGE's
// search conditions as the V1 advisor read them, the queries of INSERT,
// CREATE TABLE AS SELECT, CREATE VIEW, CREATE MATERIALIZED VIEW, CREATE
// EXTERNAL TABLE AS SELECT and CREATE REMOTE TABLE AS SELECT, and every
// SELECT nested in them (CTEs, set operations, derived tables, scalar and
// correlated subqueries, SET values). A function call is reported when an
// argument, after peeling wrappers, is an indexed column; an arithmetic
// expression when an operand is. Name is the column as written; Detail is
// the function name as the V1 advisor printed it.
//
// A table resolves in its named schema; unqualified, in "dbo", then in every
// other schema in sorted name order. A table qualified with another database
// never resolves. Column and alias resolution through scopes follows the V1
// advisor exactly, case folding included: aliases fold unless the snapshot
// is present and ObjectCaseSensitive; columns always compare lower-cased.
//
// Each finding anchors the reported function call or expression. Ported
// from Bytebase's V1 rule STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS.
func StatementWhereDisallowFunctionsAndCalculations(stmts []ast.Node, schema v1review.Schema) []v1review.Finding {
	r := &whereFunctionRule{schema: schema}
	for _, stmt := range stmts {
		r.onStatement(stmt)
	}
	return r.findings
}

type whereFunctionRule struct {
	schema   v1review.Schema
	findings []v1review.Finding
}

func (r *whereFunctionRule) onStatement(node ast.Node) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		r.checkSelect(n)
	case *ast.InsertStmt:
		r.onInsert(n)
	case *ast.CreateTableAsSelectStmt:
		r.checkSelectNode(n.Query)
	case *ast.CreateViewStmt:
		r.checkSelectNode(n.Query)
	case *ast.CreateMaterializedViewStmt:
		r.checkSelectNode(n.Query)
	case *ast.CreateExternalTableAsSelectStmt:
		r.checkSelectNode(n.Query)
	case *ast.CreateRemoteTableAsSelectStmt:
		r.checkSelectNode(n.Query)
	case *ast.UpdateStmt:
		r.onUpdate(n)
	case *ast.DeleteStmt:
		r.onDelete(n)
	case *ast.MergeStmt:
		r.checkCTEs(n.WithClause)
		r.checkMerge(n)
	default:
	}
}

// checkSelectNode runs checkSelect if node is a *ast.SelectStmt; no-ops otherwise.
func (r *whereFunctionRule) checkSelectNode(node ast.Node) {
	if sel, ok := node.(*ast.SelectStmt); ok {
		r.checkSelect(sel)
	}
}

// onInsert handles T-SQL INSERT. Source can be a SELECT or a VALUES list; both
// may embed scalar subqueries with their own WHERE clauses.
func (r *whereFunctionRule) onInsert(n *ast.InsertStmt) {
	r.checkCTEs(n.WithClause)
	r.recurseIntoNestedSelectsInNode(n.Source, nil)
}

// onUpdate handles T-SQL UPDATE. It covers the WHERE, derived tables in
// Relation/FromClause (each with its own WHERE), and scalar subqueries hidden
// in SET assignments.
func (r *whereFunctionRule) onUpdate(n *ast.UpdateStmt) {
	r.checkCTEs(n.WithClause)
	tables := r.collectUpdateDeleteTables(n.Relation, n.FromClause)
	r.checkWhere(n.WhereClause, tables)
	r.recurseDerivedTables(n.Relation, n.FromClause, tables)
	r.recurseSetClauseValues(n.SetClause, tables)
}

// onDelete handles T-SQL DELETE, including any JOINed sources in the FROM clause
// that carry their own WHERE clauses inside derived tables.
func (r *whereFunctionRule) onDelete(n *ast.DeleteStmt) {
	r.checkCTEs(n.WithClause)
	tables := r.collectUpdateDeleteTables(n.Relation, n.FromClause)
	r.checkWhere(n.WhereClause, tables)
	r.recurseDerivedTables(n.Relation, n.FromClause, tables)
}

// collectUpdateDeleteTables merges tables referenced by the target Relation
// with any JOINed tables in the extra FROM clause (T-SQL UPDATE/DELETE shape).
func (r *whereFunctionRule) collectUpdateDeleteTables(relation ast.TableExpr, from *ast.List) whereTablesByAlias {
	tables := r.collectTablesFromRelation(relation)
	for alias, ref := range r.collectTableNamesFromList(from) {
		tables[alias] = ref
	}
	return tables
}

// recurseDerivedTables walks the relation and every FROM item for nested SELECTs
// whose own WHERE clause must be checked in its own scope.
func (r *whereFunctionRule) recurseDerivedTables(relation ast.TableExpr, from *ast.List, outer whereTablesByAlias) {
	r.recurseIntoNestedSelectsInNode(relation, outer)
	if from == nil {
		return
	}
	for _, item := range from.Items {
		r.recurseIntoNestedSelectsInNode(item, outer)
	}
}

// recurseSetClauseValues walks each SET assignment's value for scalar subqueries.
func (r *whereFunctionRule) recurseSetClauseValues(setClause *ast.List, outer whereTablesByAlias) {
	if setClause == nil {
		return
	}
	for _, item := range setClause.Items {
		if se, ok := item.(*ast.SetExpr); ok && se.Value != nil {
			r.recurseIntoNestedSelectsInNode(se.Value, outer)
		}
	}
}

// checkMerge handles T-SQL MERGE statements. The ON condition and each WHEN-clause
// condition are treated like a WHERE for index-aware checking, scoped to the target
// plus source tables. Nested SELECTs inside the source or conditions are recursed.
func (r *whereFunctionRule) checkMerge(n *ast.MergeStmt) {
	tables := make(whereTablesByAlias)
	if ref, ok := n.Target.(*ast.TableRef); ok {
		addTableRef(r, tables, ref)
	}
	// The source is one FROM item carrying its own alias: a table, a derived
	// table, or a join. Its nested SELECTs are checked in their own scope.
	if n.Source != nil {
		for alias, ref := range r.collectTableNamesFromList(&ast.List{Items: []ast.Node{n.Source}}) {
			tables[alias] = ref
		}
		r.recurseIntoNestedSelectsInNode(n.Source, nil)
	}
	r.checkWhere(n.OnCondition, tables)
	if n.WhenClauses != nil {
		for _, item := range n.WhenClauses.Items {
			when, ok := item.(*ast.MergeWhenClause)
			if !ok {
				continue
			}
			r.checkWhere(when.Condition, tables)
			// MERGE actions can embed scalar subqueries with their own WHERE clauses,
			// e.g. WHEN MATCHED THEN UPDATE SET x = (SELECT ... WHERE ABS(col) > 0).
			switch a := when.Action.(type) {
			case *ast.MergeUpdateAction:
				if a != nil {
					r.recurseSetClauseValues(a.SetClause, tables)
				}
			case *ast.MergeInsertAction:
				if a != nil && a.Values != nil {
					for _, v := range a.Values.Items {
						r.recurseIntoNestedSelectsInNode(v, tables)
					}
				}
			default:
			}
		}
	}
}

func (r *whereFunctionRule) checkCTEs(wc *ast.WithClause) {
	if wc == nil || wc.CTEs == nil {
		return
	}
	for _, item := range wc.CTEs.Items {
		if cte, ok := item.(*ast.CommonTableExpr); ok {
			if sel, ok := cte.Query.(*ast.SelectStmt); ok {
				r.checkSelect(sel)
			}
		}
	}
}

func (r *whereFunctionRule) checkSelect(sel *ast.SelectStmt) {
	r.checkSelectInScope(sel, nil)
}

// checkSelectInScope is like checkSelect but also threads an outer table scope down
// so that correlated subqueries can resolve references to outer aliases. Local
// (inner) tables take precedence for unqualified columns per SQL name resolution;
// outer tables are consulted only for qualified references.
func (r *whereFunctionRule) checkSelectInScope(sel *ast.SelectStmt, outer whereTablesByAlias) {
	if sel == nil {
		return
	}
	// CTEs run before the referencing SELECT and aren't correlated — check without outer scope.
	r.checkCTEs(sel.WithClause)
	if sel.Op != ast.SetOpNone {
		r.checkSelectInScope(sel.Larg, outer)
		r.checkSelectInScope(sel.Rarg, outer)
		return
	}
	local := r.collectTableNamesFromList(sel.FromClause)
	combined := make(whereTablesByAlias, len(local)+len(outer))
	for alias, ref := range outer {
		combined[alias] = ref
	}
	for alias, ref := range local {
		combined[alias] = ref
	}
	r.checkWhereWithScopes(sel.WhereClause, local, outer, r.localHasUnresolvedSources(sel.FromClause, local))
	r.recurseIntoNestedSelects(sel, combined)
}

// localHasUnresolvedSources reports whether the given FROM list and local alias
// map contain sources whose columns we cannot enumerate — CTE references,
// derived tables (subqueries in FROM), or plain tables absent from metadata.
// When true, isIndexed must not fall back to outer scope on unqualified refs.
func (r *whereFunctionRule) localHasUnresolvedSources(from *ast.List, local whereTablesByAlias) bool {
	if from != nil {
		for _, item := range from.Items {
			derivedFound := false
			ast.Inspect(item, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.SubqueryExpr, *ast.SelectStmt:
					derivedFound = true
					return false
				}
				return !derivedFound
			})
			if derivedFound {
				return true
			}
		}
	}
	for _, ref := range local {
		if ref == (whereTableRef{}) {
			continue
		}
		if r.tableMetadata(ref) == nil {
			return true
		}
	}
	return false
}

// recurseIntoNestedSelectsInNode walks the given AST node and recurses into any
// nested *ast.SelectStmt, threading the given outer scope.
func (r *whereFunctionRule) recurseIntoNestedSelectsInNode(n ast.Node, outer whereTablesByAlias) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(m ast.Node) bool {
		if sub, ok := m.(*ast.SelectStmt); ok {
			r.checkSelectInScope(sub, outer)
			return false
		}
		return true
	})
}

// recurseIntoNestedSelects finds nested SelectStmts in the non-WHERE parts of sel
// (FROM derived tables, SELECT list, HAVING, GROUP BY, ORDER BY) and invokes
// checkSelectInScope on them with the given outer scope so correlated subqueries resolve.
func (r *whereFunctionRule) recurseIntoNestedSelects(sel *ast.SelectStmt, outer whereTablesByAlias) {
	visitor := func(m ast.Node) bool {
		if sub, ok := m.(*ast.SelectStmt); ok {
			r.checkSelectInScope(sub, outer)
			return false
		}
		return true
	}
	visitList := func(l *ast.List) {
		if l == nil {
			return
		}
		for _, item := range l.Items {
			if item != nil {
				ast.Inspect(item, visitor)
			}
		}
	}
	// WhereClause is handled by checkWhere; nested selects there are handled by
	// checkWhere's own ast.Inspect. No need to revisit it here.
	visitList(sel.FromClause)
	visitList(sel.TargetList)
	if sel.HavingClause != nil {
		ast.Inspect(sel.HavingClause, visitor)
	}
	visitList(sel.GroupByClause)
	visitList(sel.OrderByClause)
}

// whereTableRef captures the real database, schema, and object name for an aliased table reference.
type whereTableRef struct {
	database string // empty means current database
	schema   string // empty means default (dbo)
	object   string
}

// whereTablesByAlias maps alias (or object name if no alias) to the underlying whereTableRef.
type whereTablesByAlias = map[string]whereTableRef

// collectTableNamesFromList extracts table refs from a FROM clause List, preserving schema.
// Descent is stopped at subquery boundaries so that derived tables don't leak their
// inner TableRefs into the outer scope, but a derived table's own alias (captured
// via AliasedTableRef wrapping a SubqueryExpr) is recorded as an opaque entry so
// qualifier-based resolution can tell it apart from an outer alias of the same name.
func (r *whereFunctionRule) collectTableNamesFromList(from *ast.List) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	if from == nil {
		return tables
	}
	for _, item := range from.Items {
		ast.Inspect(item, func(n ast.Node) bool {
			// AliasedTableRef wraps a real TableRef or a SubqueryExpr with an alias.
			// For subquery sources, record the alias before we stop descent so that
			// inner qualified refs resolve locally rather than to an outer table.
			if atr, ok := n.(*ast.AliasedTableRef); ok {
				if _, isSub := atr.Table.(*ast.SubqueryExpr); isSub && atr.Alias != "" {
					key := r.normalizeIdent(atr.Alias)
					if _, exists := tables[key]; !exists {
						tables[key] = whereTableRef{}
					}
					return false
				}
				return true
			}
			if _, ok := n.(*ast.SubqueryExpr); ok {
				return false
			}
			if _, ok := n.(*ast.SelectStmt); ok {
				return false
			}
			t, ok := n.(*ast.TableRef)
			if !ok {
				return true
			}
			addTableRef(r, tables, t)
			return false
		})
	}
	return tables
}

// collectTablesFromRelation builds a table map from a single table reference (UPDATE/DELETE).
func (r *whereFunctionRule) collectTablesFromRelation(te ast.TableExpr) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	if te == nil {
		return tables
	}
	if ref, ok := te.(*ast.TableRef); ok {
		addTableRef(r, tables, ref)
	}
	return tables
}

// normalizeIdent lowercases identifiers when the DB is configured case-insensitive
// (respecting IsObjectCaseSensitive). MSSQL default collation is CI, but a CS DB
// must be honored too.
func (r *whereFunctionRule) normalizeIdent(name string) string {
	if r.schema.Metadata != nil && r.schema.ObjectCaseSensitive {
		return name
	}
	return strings.ToLower(name)
}

func addTableRef(r *whereFunctionRule, tables whereTablesByAlias, t *ast.TableRef) {
	key := t.Alias
	if key == "" {
		key = t.Object
	}
	normKey := r.normalizeIdent(key)
	newRef := whereTableRef{database: t.Database, schema: t.Schema, object: t.Object}
	// If an unaliased same-name from a different schema/database already occupies this
	// key (rare: `FROM dbo.orders, sales.orders` without aliases), mark the entry as
	// ambiguous so isIndexed returns false rather than silently resolving to the wrong
	// table. Users should alias these tables to disambiguate.
	if existing, ok := tables[normKey]; ok && (existing.database != newRef.database || existing.schema != newRef.schema || existing.object != newRef.object) {
		tables[normKey] = whereTableRef{}
		return
	}
	tables[normKey] = newRef
}

// tableMetadata fetches the store-model table for the given ref, applying the
// cross-DB guard and MSSQL's dbo-first-then-sorted-others schema resolution.
// Returns nil when not resolvable.
func (r *whereFunctionRule) tableMetadata(ref whereTableRef) *metadata.TableMetadata {
	if r.schema.Metadata == nil {
		return nil
	}
	if ref.database != "" && !strings.EqualFold(ref.database, r.schema.Metadata.GetName()) {
		return nil
	}
	if ref.schema != "" {
		return r.schema.Table(ref.schema, ref.object)
	}
	if t := r.schema.Table("dbo", ref.object); t != nil {
		return t
	}
	// The V1 advisor sorted the schema names for a deterministic fallback.
	names := r.schema.SchemaNames()
	slices.Sort(names)
	for _, name := range names {
		if strings.EqualFold(name, "dbo") {
			continue
		}
		if t := r.schema.Table(name, ref.object); t != nil {
			return t
		}
	}
	return nil
}

// indexedColumns returns indexed column names for the table at the given database+schema+object.
func (r *whereFunctionRule) indexedColumns(ref whereTableRef) map[string]bool {
	table := r.tableMetadata(ref)
	if table == nil {
		return nil
	}
	cols := make(map[string]bool)
	for _, idx := range table.GetIndexes() {
		for _, expr := range idx.GetExpressions() {
			cols[strings.ToLower(expr)] = true
		}
	}
	if len(cols) == 0 {
		return nil
	}
	return cols
}

// allColumns returns every column name defined for the given table reference.
func (r *whereFunctionRule) allColumns(ref whereTableRef) map[string]bool {
	table := r.tableMetadata(ref)
	if table == nil {
		return nil
	}
	cols := make(map[string]bool)
	for _, c := range table.GetColumns() {
		cols[strings.ToLower(c.GetName())] = true
	}
	if len(cols) == 0 {
		return nil
	}
	return cols
}

// whereIndexedByAlias maps alias (or object if no alias) to the set of indexed columns on that table.
type whereIndexedByAlias = map[string]map[string]bool

// whereScopeIndexed carries per-alias column metadata for the local (inner) and outer scopes.
// Unqualified refs resolve against local first (any table that declares the column);
// only when no local table declares the column — and the local scope is fully
// resolvable — do we fall back to outer. If any local source has unknown columns
// (CTEs, derived tables, tables without metadata), we stay silent to avoid
// misattributing an inner column to an outer indexed column of the same name.
type whereScopeIndexed struct {
	local           whereIndexedByAlias
	localAll        whereIndexedByAlias
	localAliases    map[string]bool
	outer           whereIndexedByAlias
	localUnresolved bool
}

func (s whereScopeIndexed) empty() bool {
	return len(s.local) == 0 && len(s.outer) == 0
}

// resolveIndexedColumns builds a per-alias map of indexed columns.
func (r *whereFunctionRule) resolveIndexedColumns(tables whereTablesByAlias) whereIndexedByAlias {
	all := make(whereIndexedByAlias)
	for alias, ref := range tables {
		if cols := r.indexedColumns(ref); cols != nil {
			all[alias] = cols
		}
	}
	if len(all) == 0 {
		return nil
	}
	return all
}

// resolveAllColumns builds a per-alias map of every known column for each table.
func (r *whereFunctionRule) resolveAllColumns(tables whereTablesByAlias) whereIndexedByAlias {
	all := make(whereIndexedByAlias)
	for alias, ref := range tables {
		if cols := r.allColumns(ref); cols != nil {
			all[alias] = cols
		}
	}
	if len(all) == 0 {
		return nil
	}
	return all
}

// isIndexed reports whether a column reference is indexed, respecting its table qualifier
// and SQL name resolution for correlated subqueries.
//   - Qualified (ref.Table set): check local first, then outer.
//   - Unqualified: if any local table declares the column, local wins. Only if no
//     local table declares the column do we fall back to outer indexed columns.
func (r *whereFunctionRule) isIndexed(ref *ast.ColumnRef, s whereScopeIndexed) bool {
	if ref == nil || s.empty() {
		return false
	}
	col := strings.ToLower(ref.Column)
	if ref.Table != "" {
		alias := r.normalizeIdent(ref.Table)
		// Local scope takes precedence: if the alias exists locally (even as a
		// derived table or CTE without index metadata), the ref resolves there.
		// Otherwise the qualified alias is definitively outer, even when other
		// local sources have unresolved metadata — collectTableNamesFromList
		// records every local alias (including derived tables) so this check
		// is sound.
		if s.localAliases[alias] {
			return s.local[alias][col]
		}
		return s.outer[alias][col]
	}
	localClaims := false
	for _, cols := range s.localAll {
		if cols[col] {
			localClaims = true
			break
		}
	}
	if localClaims {
		for _, cols := range s.local {
			if cols[col] {
				return true
			}
		}
		return false
	}
	// No local alias declared the column. Only fall back to outer when the local
	// scope is fully resolvable; otherwise a CTE/derived-table column could be
	// misattributed to an outer indexed column of the same name.
	if s.localUnresolved {
		return false
	}
	for _, cols := range s.outer {
		if cols[col] {
			return true
		}
	}
	return false
}

func (r *whereFunctionRule) checkWhere(where ast.ExprNode, tables whereTablesByAlias) {
	r.checkWhereWithScopes(where, tables, nil, false)
}

func (r *whereFunctionRule) checkWhereWithScopes(where ast.ExprNode, local, outer whereTablesByAlias, localUnresolved bool) {
	if where == nil {
		return
	}
	localAliases := make(map[string]bool, len(local))
	for alias := range local {
		localAliases[alias] = true
	}
	s := whereScopeIndexed{
		local:           r.resolveIndexedColumns(local),
		localAll:        r.resolveAllColumns(local),
		localAliases:    localAliases,
		outer:           r.resolveIndexedColumns(outer),
		localUnresolved: localUnresolved,
	}
	combined := make(whereTablesByAlias, len(local)+len(outer))
	for alias, ref := range outer {
		combined[alias] = ref
	}
	for alias, ref := range local {
		combined[alias] = ref
	}
	ast.Inspect(where, func(n ast.Node) bool {
		if sub, ok := n.(*ast.SelectStmt); ok {
			r.checkSelectInScope(sub, combined)
			return false
		}
		if !s.empty() {
			r.checkWhereNode(n, s)
		}
		return true
	})
}

func (r *whereFunctionRule) checkWhereNode(n ast.Node, indexed whereScopeIndexed) {
	switch expr := n.(type) {
	case *ast.FuncCallExpr:
		if col := r.findIndexedColumnInFuncArgs(expr.Args, indexed); col != "" {
			funcName := ""
			if expr.Name != nil {
				funcName = strings.ToUpper(whereFunctionNameText(expr.Name))
			}
			r.addFunction(funcName, col, expr.Loc)
		}
	case *ast.BinaryExpr:
		if isWhereArithmeticOp(expr.Op) {
			if col := r.unwrapAndFindIndexedColumn(expr.Left, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
			} else if col := r.unwrapAndFindIndexedColumn(expr.Right, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
			}
		}
	case *ast.UnaryExpr:
		// UnaryPlus is a no-op and does not prevent index usage; skip it.
		if expr.Op == ast.UnaryMinus || expr.Op == ast.UnaryBitNot {
			if col := r.unwrapAndFindIndexedColumn(expr.Operand, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
			}
		}
	default:
	}
}

func (r *whereFunctionRule) addFunction(funcName, col string, loc ast.Loc) {
	start, end := stmtAnchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.FunctionOnIndexedColumn, Name: col, Detail: funcName, Start: start, End: end})
}

func (r *whereFunctionRule) addCalculation(col string, loc ast.Loc) {
	start, end := stmtAnchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.CalculationOnIndexedColumn, Name: col, Start: start, End: end})
}

// findIndexedColumnInFuncArgs searches function arguments for an indexed ColumnRef,
// descending through passthrough wrappers (ParenExpr, CastExpr, ConvertExpr) so that
// `ABS((c1))` and `ABS(CAST(c1 AS int))` are still flagged.
// BinaryExpr/UnaryExpr (arithmetic) are NOT crossed — those are handled by the
// calculation advice path, so `ABS(c1 + 1)` is reported once as a calculation.
func (r *whereFunctionRule) findIndexedColumnInFuncArgs(args *ast.List, indexed whereScopeIndexed) string {
	if args == nil {
		return ""
	}
	for _, item := range args.Items {
		arg, ok := item.(ast.ExprNode)
		if !ok {
			continue
		}
		if col := r.unwrapAndFindIndexedColumn(arg, indexed); col != "" {
			return col
		}
	}
	return ""
}

// unwrapAndFindIndexedColumn peels ParenExpr/CastExpr/ConvertExpr wrappers off expr
// and returns the column name if the inner expression is an indexed ColumnRef.
func (r *whereFunctionRule) unwrapAndFindIndexedColumn(expr ast.ExprNode, indexed whereScopeIndexed) string {
	for expr != nil {
		switch e := expr.(type) {
		case *ast.ColumnRef:
			if r.isIndexed(e, indexed) {
				return e.Column
			}
			return ""
		case *ast.ParenExpr:
			expr = e.Expr
		case *ast.CastExpr:
			expr = e.Expr
		case *ast.ConvertExpr:
			expr = e.Expr
		default:
			return ""
		}
	}
	return ""
}

func isWhereArithmeticOp(op ast.BinaryOp) bool {
	switch op {
	case ast.BinOpAdd, ast.BinOpSub, ast.BinOpMul, ast.BinOpDiv, ast.BinOpMod,
		ast.BinOpBitAnd, ast.BinOpBitOr, ast.BinOpBitXor:
		return true
	default:
		return false
	}
}

// whereFunctionNameText returns a function's name as the V1 advisor printed
// it: database.schema.object, each part present only when written.
func whereFunctionNameText(ref *ast.TableRef) string {
	if ref == nil {
		return ""
	}
	result := ref.Object
	if ref.Schema != "" {
		result = ref.Schema + "." + result
	}
	if ref.Database != "" {
		result = ref.Database + "." + result
	}
	return result
}
