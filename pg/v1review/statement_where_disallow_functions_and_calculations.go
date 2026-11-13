package v1review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereDisallowFunctionsAndCalculations reports each WHERE-clause
// function call or calculation applied to an indexed column, judged against
// the pre-change schema.
//
// Judged: the WHERE and JOIN ON conditions of top-level SELECT, UPDATE (with
// its FROM) and DELETE (with its USING); of INSERT's ON CONFLICT WHERE and
// its index-inference WHERE; of MERGE's join condition, source JOIN ON and
// WHEN conditions; of the queries of INSERT ... SELECT, CREATE VIEW and
// CREATE TABLE AS / CREATE MATERIALIZED VIEW; and of every SELECT nested in
// them (CTEs, data-modifying CTEs, set operations, derived tables, scalar
// and correlated subqueries in select lists, VALUES, HAVING, GROUP BY,
// ORDER BY, SET values and RETURNING). HAVING itself is never judged, nor is
// anything under EXPLAIN or in a routine body. A function call is reported
// when one of its arguments, after peeling casts, is an indexed column; a
// binary + - * / % ^ & | # << >> when its left, else its right, operand is;
// a unary - or ~ when its operand is. Name is the column, folded as names
// are; Detail is the function's last name part upper-cased.
//
// A column is indexed when it appears among the expressions of any index of
// its table. A table qualified with a schema resolves there; unqualified, in
// the first schema of searchPath holding it, unless a CTE of that name is in
// scope (a WITH RECURSIVE body sees every sibling, itself included; a plain
// WITH body only earlier ones). Tables resolve through Schema.Table, partitions
// included. A qualified column resolves against its alias (the last
// qualifier) in the local FROM, else the enclosing queries'. An unqualified
// column resolves locally when any local table declares it, and is reported
// only when such a table indexes it; otherwise it falls back to the
// enclosing queries unless a local source's columns are unknown (a derived
// table, function, CTE, aliased join, or a table missing from the snapshot).
// Two FROM entries under one alias or name that name different tables make
// it unresolvable. Names compare as written when the snapshot is present and
// ObjectCaseSensitive, lower-cased otherwise.
//
// Each finding anchors the reported function call or expression. Ported
// from Bytebase's V1 rule STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS.

func StatementWhereDisallowFunctionsAndCalculations(stmts []ast.Node, schema v1review.Schema, searchPath []string) []v1review.Finding {
	r := &whereFunctionRule{schema: schema, searchPath: searchPath}
	for _, stmt := range stmts {
		r.onStatement(stmt)
	}
	return r.findings
}

type whereFunctionRule struct {
	schema   v1review.Schema
	findings []v1review.Finding
	// searchPath is the list of schemas to consult for unqualified table refs,
	// in priority order, as the caller resolved PG's `search_path`.
	searchPath []string
	// cteStack is a stack of CTE-name sets visible at the current nesting
	// level. Pushed on entry to a query that has a WITH clause; popped on
	// exit. A FROM-clause name that matches any stacked entry is treated
	// as opaque (`whereTableRef{}`) so it never picks up real-table indexes.
	// Inner WITHs shadow outer per normal SQL scoping.
	cteStack []map[string]bool
}

// pushCTEs pushes a new CTE-name set onto the visibility stack.
func (r *whereFunctionRule) pushCTEs(names map[string]bool) {
	r.cteStack = append(r.cteStack, names)
}

// popCTEs pops the most-recently pushed CTE-name set.
func (r *whereFunctionRule) popCTEs() {
	if n := len(r.cteStack); n > 0 {
		r.cteStack = r.cteStack[:n-1]
	}
}

// cteVisible reports whether `name` (already normalized) is in scope as a CTE.
func (r *whereFunctionRule) cteVisible(name string) bool {
	for _, frame := range r.cteStack {
		if frame[name] {
			return true
		}
	}
	return false
}

// whereTableRef captures a resolved table reference.
type whereTableRef struct {
	schema string // empty means resolve via search_path
	name   string
}

// whereTablesByAlias maps alias (or name if no alias) to the underlying whereTableRef.
type whereTablesByAlias = map[string]whereTableRef

// whereIndexedByAlias maps alias to a set of column names for that alias.
type whereIndexedByAlias = map[string]map[string]bool

// whereScopeIndexed carries per-alias column metadata for the local (inner) scope
// and any outer scope inherited from enclosing queries.
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

// onStatement dispatches each top-level AST node to the detector.
func (r *whereFunctionRule) onStatement(node ast.Node) {
	switch n := unwrapRaw(node).(type) {
	case *ast.SelectStmt:
		r.checkSelectInScope(n, nil)
	case *ast.InsertStmt:
		r.onInsert(n)
	case *ast.UpdateStmt:
		r.onUpdate(n)
	case *ast.DeleteStmt:
		r.onDelete(n)
	case *ast.ViewStmt:
		r.onViewStmt(n)
	case *ast.CreateTableAsStmt:
		r.onCreateTableAs(n)
	case *ast.MergeStmt:
		r.onMerge(n)
	default:
	}
}

// ---- Identifier / metadata helpers ---------------------------------------

// normalizeIdent returns the canonical lookup key for an identifier. PG is
// case-sensitive for unquoted identifiers in our fixtures (all lowercase).
func (r *whereFunctionRule) normalizeIdent(name string) string {
	if r.schema.Metadata != nil && r.schema.ObjectCaseSensitive {
		return name
	}
	return strings.ToLower(name)
}

// tableMetadata fetches the snapshot's table for the given ref. PG resolves
// unqualified table refs through the `search_path` GUC: each schema in
// priority order is consulted until a match is found.
func (r *whereFunctionRule) tableMetadata(ref whereTableRef) *metadata.TableMetadata {
	if r.schema.Metadata == nil {
		return nil
	}
	if ref.schema != "" {
		return r.schema.Table(ref.schema, ref.name)
	}
	for _, schema := range r.searchPath {
		if t := r.schema.Table(schema, ref.name); t != nil {
			return t
		}
	}
	return nil
}

// indexedColumns returns indexed column names for the given table reference.
// Keys go through normalizeIdent so case-sensitive schemas distinguish
// quoted mixed-case identifiers.
func (r *whereFunctionRule) indexedColumns(ref whereTableRef) map[string]bool {
	table := r.tableMetadata(ref)
	if table == nil {
		return nil
	}
	cols := make(map[string]bool)
	for _, idx := range table.GetIndexes() {
		for _, expr := range idx.GetExpressions() {
			cols[r.normalizeIdent(expr)] = true
		}
	}
	if len(cols) == 0 {
		return nil
	}
	return cols
}

// allColumns returns every column name defined for the given table reference.
// Keys use normalizeIdent for the same reason as indexedColumns.
func (r *whereFunctionRule) allColumns(ref whereTableRef) map[string]bool {
	table := r.tableMetadata(ref)
	if table == nil {
		return nil
	}
	cols := make(map[string]bool)
	for _, c := range table.GetColumns() {
		cols[r.normalizeIdent(c.GetName())] = true
	}
	if len(cols) == 0 {
		return nil
	}
	return cols
}

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

// ---- FROM-clause collection ----------------------------------------------

// collectTableNames walks a FROM-clause list and extracts table aliases/refs.
// Descent stops at subquery boundaries so inner SELECTs do not leak their
// sources into the outer scope; derived-table aliases are recorded as
// opaque whereTableRef{} entries so qualifier-based resolution can still tell
// them apart from an outer alias of the same name.
func (r *whereFunctionRule) collectTableNames(from *ast.List) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	if from == nil {
		return tables
	}
	for _, item := range from.Items {
		r.collectTableNamesFromNode(item, tables)
	}
	return tables
}

func (r *whereFunctionRule) collectTableNamesFromNode(node ast.Node, tables whereTablesByAlias) {
	if node == nil {
		return
	}
	switch n := node.(type) {
	case *ast.RangeVar:
		key := ""
		if n.Alias != nil && n.Alias.Aliasname != "" {
			key = n.Alias.Aliasname
		} else {
			key = n.Relname
		}
		normKey := r.normalizeIdent(key)
		// CTE shadowing: if this RangeVar's underlying name (NOT the alias)
		// matches a CTE in scope, the FROM-ref binds to the CTE result, not
		// the physical table. Record opaque so we don't attach base-table
		// indexes to CTE columns. Schema-qualified refs (`s.t`) bypass this
		// since CTEs are unqualified.
		if n.Schemaname == "" && r.cteVisible(r.normalizeIdent(n.Relname)) {
			tables[normKey] = whereTableRef{}
			return
		}
		newRef := whereTableRef{schema: n.Schemaname, name: n.Relname}
		if existing, ok := tables[normKey]; ok && (existing.schema != newRef.schema || existing.name != newRef.name) {
			tables[normKey] = whereTableRef{}
		} else {
			tables[normKey] = newRef
		}
	case *ast.RangeSubselect:
		if n.Alias != nil && n.Alias.Aliasname != "" {
			key := r.normalizeIdent(n.Alias.Aliasname)
			if _, exists := tables[key]; !exists {
				tables[key] = whereTableRef{}
			}
		}
	case *ast.RangeFunction:
		if n.Alias != nil && n.Alias.Aliasname != "" {
			key := r.normalizeIdent(n.Alias.Aliasname)
			if _, exists := tables[key]; !exists {
				tables[key] = whereTableRef{}
			}
		}
	case *ast.JoinExpr:
		r.collectTableNamesFromNode(n.Larg, tables)
		r.collectTableNamesFromNode(n.Rarg, tables)
		if n.Alias != nil && n.Alias.Aliasname != "" {
			key := r.normalizeIdent(n.Alias.Aliasname)
			if _, exists := tables[key]; !exists {
				tables[key] = whereTableRef{}
			}
		}
	default:
	}
}

// localHasUnresolvedSources reports whether any FROM entry hides columns we
// cannot enumerate. This walks SOURCE POSITIONS only — top-level FROM items
// and JoinExpr.Larg/Rarg — looking for opaque producers (derived tables,
// RangeSubselect/SubLink/SelectStmt, RangeFunction). Subqueries appearing
// inside JoinExpr.Quals are VALUE expressions, not source producers, and
// must NOT count toward unresolved: treating them as opaque suppresses the
// outer-scope fallback in isIndexed and produces false negatives for
// correlated refs.
func (r *whereFunctionRule) localHasUnresolvedSources(from *ast.List, local whereTablesByAlias) bool {
	if from != nil {
		for _, item := range from.Items {
			if whereHasOpaqueSource(item) {
				return true
			}
		}
	}
	for _, ref := range local {
		if ref == (whereTableRef{}) {
			// Opaque alias (CTE-shadowed RangeVar, JoinExpr, RangeSubselect,
			// or RangeFunction) — by definition we can't enumerate its
			// columns, so the local scope is unresolved.
			return true
		}
		if r.tableMetadata(ref) == nil {
			return true
		}
	}
	return false
}

// whereHasOpaqueSource reports whether a FROM item (or its JOIN sub-tree) contains
// an opaque source node. Only SOURCE positions are inspected:
//   - the item itself
//   - for JoinExpr, Larg and Rarg (NOT Quals — those are predicates)
func whereHasOpaqueSource(n ast.Node) bool {
	if n == nil {
		return false
	}
	switch x := n.(type) {
	case *ast.SubLink, *ast.RangeSubselect, *ast.SelectStmt, *ast.RangeFunction:
		return true
	case *ast.JoinExpr:
		return whereHasOpaqueSource(x.Larg) || whereHasOpaqueSource(x.Rarg)
	}
	return false
}

// ---- Column reference resolution -----------------------------------------

// extractColumnRefParts returns (qualifier, column) for a ColumnRef.
// Returns ("", "") when the ref is not a plain column.
func (*whereFunctionRule) extractColumnRefParts(ref *ast.ColumnRef) (string, string) {
	if ref == nil || ref.Fields == nil || len(ref.Fields.Items) == 0 {
		return "", ""
	}
	parts := make([]string, 0, len(ref.Fields.Items))
	for _, f := range ref.Fields.Items {
		s, ok := f.(*ast.String)
		if !ok {
			return "", ""
		}
		parts = append(parts, s.Str)
	}
	switch len(parts) {
	case 1:
		return "", parts[0]
	default:
		// a.b -> qualifier=a, column=b. a.b.c -> qualifier=b, column=c
		// (schema 'a' ignored; alias lookup uses the last qualifier).
		return parts[len(parts)-2], parts[len(parts)-1]
	}
}

// isIndexed reports whether a column reference is indexed within scope s.
func (r *whereFunctionRule) isIndexed(ref *ast.ColumnRef, s whereScopeIndexed) bool {
	if ref == nil || s.empty() {
		return false
	}
	qualifier, column := r.extractColumnRefParts(ref)
	if column == "" {
		return false
	}
	col := r.normalizeIdent(column)
	if qualifier != "" {
		alias := r.normalizeIdent(qualifier)
		if s.localAliases[alias] {
			return s.local[alias][col]
		}
		return s.outer[alias][col]
	}
	// Unqualified: does any local table declare this column?
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

// ---- Detection ------------------------------------------------------------

// checkSelectInScope is the main recursive entry point for SELECT queries,
// threading an outer scope for correlated subqueries.
func (r *whereFunctionRule) checkSelectInScope(sel *ast.SelectStmt, outer whereTablesByAlias) {
	if sel == nil {
		return
	}
	// CTE visibility per SQL: non-recursive WITH bodies see only earlier
	// siblings; WITH RECURSIVE bodies see all siblings (including self).
	// handleWithClause walks each body in declaration order and pushes
	// names with the right ordering; cleanup pops at our exit.
	popCTEs := r.handleWithClause(sel.WithClause)
	defer popCTEs()
	if sel.Op != ast.SETOP_NONE {
		r.checkSelectInScope(sel.Larg, outer)
		r.checkSelectInScope(sel.Rarg, outer)
		return
	}
	local := r.collectTableNames(sel.FromClause)
	unresolved := r.localHasUnresolvedSources(sel.FromClause, local)
	// Build combined scope so inner subqueries see current local+outer as outer.
	combined := make(whereTablesByAlias, len(local)+len(outer))
	for alias, ref := range outer {
		combined[alias] = ref
	}
	for alias, ref := range local {
		combined[alias] = ref
	}
	r.checkWhereWithScopes(sel.WhereClause, local, outer, unresolved)
	// JOIN ON predicates are bound to the same local scope as WHERE and must
	// be checked for function/calculation violations on indexed columns.
	r.checkJoinConditions(sel.FromClause, local, outer, unresolved)
	// Recurse into nested subqueries elsewhere (target list, FROM, HAVING,
	// GROUP BY, ORDER BY) passing current combined scope as outer. HAVING
	// predicates are never flagged directly; nested subqueries inside HAVING
	// are still recursed into.
	r.recurseIntoNestedSelects(sel, combined)
}

// checkJoinConditions walks a FROM clause and runs WHERE-detection on each
// JoinExpr.Quals encountered. JOIN ON predicates are semantically WHERE-like
// — a function/calculation on an indexed column still defeats the index.
func (r *whereFunctionRule) checkJoinConditions(from *ast.List, local, outer whereTablesByAlias, unresolved bool) {
	if from == nil {
		return
	}
	for _, item := range from.Items {
		r.checkJoinConditionsInNode(item, local, outer, unresolved)
	}
}

func (r *whereFunctionRule) checkJoinConditionsInNode(node ast.Node, local, outer whereTablesByAlias, unresolved bool) {
	if node == nil {
		return
	}
	j, ok := node.(*ast.JoinExpr)
	if !ok {
		return
	}
	r.checkWhereWithScopes(j.Quals, local, outer, unresolved)
	r.checkJoinConditionsInNode(j.Larg, local, outer, unresolved)
	r.checkJoinConditionsInNode(j.Rarg, local, outer, unresolved)
}

// walkCTEBody dispatches a single CTE body to the appropriate handler.
// The CTE body is checked with fresh scope (CTEs are not correlated with
// their referencing SELECT). Data-modifying CTEs (INSERT/UPDATE/DELETE/
// MERGE) dispatch to their own handlers.
func (r *whereFunctionRule) walkCTEBody(query ast.Node) {
	switch q := query.(type) {
	case *ast.SelectStmt:
		r.checkSelectInScope(q, nil)
	case *ast.UpdateStmt:
		r.onUpdate(q)
	case *ast.DeleteStmt:
		r.onDelete(q)
	case *ast.InsertStmt:
		r.onInsert(q)
	case *ast.MergeStmt:
		r.onMerge(q)
	default:
	}
}

// handleWithClause walks the CTE bodies AND maintains the cteStack with
// SQL-correct visibility, returning a cleanup func the caller must defer.
//
// CTE visibility matrix (SQL:2011):
//
//	| Mode                  | body sees self | preceding siblings | later siblings |
//	|-----------------------|----------------|--------------------|----------------|
//	| Non-recursive WITH    |       No       |         Yes        |        No      |
//	| WITH RECURSIVE        |       Yes      |         Yes        |  Yes (mutual)  |
//
// The distinction matters for this rule because when a CTE name shadows a
// real base table (e.g. `WITH orders AS (SELECT * FROM orders WHERE …) …`):
//   - In NON-recursive WITH, inner `FROM orders` resolves to the REAL base
//     table (per SQL spec). We MUST resolve indexes against it, otherwise
//     we silently miss `UPPER(indexed_col)` violations — a false negative.
//   - In WITH RECURSIVE, inner `FROM orders` resolves to the CTE itself
//     (self-reference). We MUST treat it as opaque to avoid attaching the
//     same-named base table's indexes to the CTE — a false positive.
//
// Strategy: push-timing depends on `with.Recursive`.
//   - Recursive: collect all CTE names first, push the whole frame ONCE,
//     then walk every body. Every body sees every name (self + others).
//   - Non-recursive: push an EMPTY frame, then for each CTE: walk its body
//     first (current name invisible), then add its name to the frame so
//     subsequent siblings see it.
//
// On exit, the frame contains every CTE name in both modes, so the OUTER
// query (walked by the caller after handleWithClause returns) always sees
// the full set of CTE names — preserving the existing opaque-on-outer
// behavior.
func (r *whereFunctionRule) handleWithClause(with *ast.WithClause) func() {
	noop := func() {
		// Intentionally empty: no CTE frame was pushed, so nothing to pop.
	}
	if with == nil || with.Ctes == nil {
		return noop
	}
	frame := make(map[string]bool)
	r.pushCTEs(frame)

	if with.Recursive {
		for _, item := range with.Ctes.Items {
			cte, ok := item.(*ast.CommonTableExpr)
			if !ok || cte.Ctename == "" {
				continue
			}
			frame[r.normalizeIdent(cte.Ctename)] = true
		}
		for _, item := range with.Ctes.Items {
			cte, ok := item.(*ast.CommonTableExpr)
			if !ok || cte.Ctequery == nil {
				continue
			}
			r.walkCTEBody(cte.Ctequery)
		}
		return r.popCTEs
	}

	for _, item := range with.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok || cte.Ctequery == nil {
			continue
		}
		r.walkCTEBody(cte.Ctequery)
		if cte.Ctename != "" {
			frame[r.normalizeIdent(cte.Ctename)] = true
		}
	}
	return r.popCTEs
}

// checkWhereWithScopes builds the whereScopeIndexed and inspects the WHERE body,
// dispatching to the detector for each non-subquery node. Nested subqueries
// are recursed into with the combined (local+outer) scope as their outer.
func (r *whereFunctionRule) checkWhereWithScopes(where ast.Node, local, outer whereTablesByAlias, localUnresolved bool) {
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
		switch x := n.(type) {
		case *ast.SubLink:
			if sub, ok := x.Subselect.(*ast.SelectStmt); ok {
				r.checkSelectInScope(sub, combined)
			}
			// Still visit Testexpr for the outer side of ANY/ALL.
			if x.Testexpr != nil {
				ast.Inspect(x.Testexpr, func(m ast.Node) bool {
					switch y := m.(type) {
					case *ast.SelectStmt:
						r.checkSelectInScope(y, combined)
						return false
					case *ast.SubLink:
						if sub2, ok := y.Subselect.(*ast.SelectStmt); ok {
							r.checkSelectInScope(sub2, combined)
						}
						return false
					case *ast.RangeSubselect:
						if sub2, ok := y.Subquery.(*ast.SelectStmt); ok {
							r.checkSelectInScope(sub2, combined)
						}
						return false
					}
					if !s.empty() {
						r.checkWhereNode(m, s)
					}
					return true
				})
			}
			return false
		case *ast.RangeSubselect:
			if sub, ok := x.Subquery.(*ast.SelectStmt); ok {
				r.checkSelectInScope(sub, combined)
			}
			return false
		case *ast.SelectStmt:
			r.checkSelectInScope(x, combined)
			return false
		}
		if !s.empty() {
			r.checkWhereNode(n, s)
		}
		return true
	})
}

// checkWhereNode runs the detector cells on a single node.
func (r *whereFunctionRule) checkWhereNode(n ast.Node, indexed whereScopeIndexed) {
	switch expr := n.(type) {
	case *ast.FuncCall:
		if col := r.findIndexedColumnInArgs(expr.Args, indexed); col != "" {
			r.addFunction(r.funcCallName(expr), col, expr.Loc)
		}
	case *ast.A_Expr:
		if expr.Kind != ast.AEXPR_OP {
			return
		}
		op := whereAExprOpSymbol(expr)
		if op == "" {
			return
		}
		if expr.Lexpr == nil {
			// Unary op.
			if op == "-" || op == "~" {
				if col := r.unwrapAndFindIndexedColumn(expr.Rexpr, indexed); col != "" {
					r.addCalculation(col, expr.Loc)
				}
			}
			return
		}
		if whereIsArithmeticOp(op) {
			if col := r.unwrapAndFindIndexedColumn(expr.Lexpr, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
				return
			}
			if col := r.unwrapAndFindIndexedColumn(expr.Rexpr, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
			}
		}
	default:
	}
}

// funcCallName returns the function's simple name, upper-cased.
func (*whereFunctionRule) funcCallName(fc *ast.FuncCall) string {
	if fc == nil || fc.Funcname == nil || len(fc.Funcname.Items) == 0 {
		return ""
	}
	last, ok := fc.Funcname.Items[len(fc.Funcname.Items)-1].(*ast.String)
	if !ok {
		return ""
	}
	return strings.ToUpper(last.Str)
}

// whereAExprOpSymbol extracts the operator symbol from an A_Expr.Name list.
func whereAExprOpSymbol(expr *ast.A_Expr) string {
	if expr == nil || expr.Name == nil || len(expr.Name.Items) == 0 {
		return ""
	}
	last, ok := expr.Name.Items[len(expr.Name.Items)-1].(*ast.String)
	if !ok {
		return ""
	}
	return last.Str
}

// whereIsArithmeticOp reports whether a binary operator symbol defeats index use
// on its operand column. The inventory is derived from PostgreSQL's operator
// set (not mirrored from MySQL, whose set differs). Covered:
//   - arithmetic: + - * / % ^ (^ is PG exponent, not XOR)
//   - bitwise: & | # << >>
func whereIsArithmeticOp(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%", "^", "&", "|", "#", "<<", ">>":
		return true
	}
	return false
}

// findIndexedColumnInArgs searches function arguments for an indexed ColumnRef,
// unwrapping passthrough wrappers (TypeCast) so CAST(name AS TEXT) still flags.
// Arithmetic (A_Expr) is NOT crossed — those are handled by the calculation
// path so `ABS(id + 1)` is reported as a calculation.
func (r *whereFunctionRule) findIndexedColumnInArgs(args *ast.List, indexed whereScopeIndexed) string {
	if args == nil {
		return ""
	}
	for _, arg := range args.Items {
		if col := r.unwrapAndFindIndexedColumn(arg, indexed); col != "" {
			return col
		}
	}
	return ""
}

// unwrapAndFindIndexedColumn peels passthrough wrappers (TypeCast) and, if the
// inner expression is an indexed ColumnRef, returns its normalized column
// name. Returns "" otherwise.
func (r *whereFunctionRule) unwrapAndFindIndexedColumn(expr ast.Node, indexed whereScopeIndexed) string {
	for expr != nil {
		switch e := expr.(type) {
		case *ast.ColumnRef:
			if r.isIndexed(e, indexed) {
				_, col := r.extractColumnRefParts(e)
				return r.normalizeIdent(col)
			}
			return ""
		case *ast.TypeCast:
			expr = e.Arg
		default:
			return ""
		}
	}
	return ""
}

// recurseIntoNestedSelects finds nested *SelectStmt / *SubLink / *RangeSubselect
// nodes in the non-WHERE parts of the given SELECT and invokes
// checkSelectInScope on each with the given outer scope. HAVING predicates
// are NEVER flagged themselves, but nested subqueries inside HAVING are still
// recursed into.
func (r *whereFunctionRule) recurseIntoNestedSelects(sel *ast.SelectStmt, outer whereTablesByAlias) {
	if sel.FromClause != nil {
		for _, f := range sel.FromClause.Items {
			r.recurseIntoNestedSelectsInFromItem(f, outer)
		}
	}
	if sel.TargetList != nil {
		for _, t := range sel.TargetList.Items {
			r.recurseIntoNestedSelectsInNode(t, outer)
		}
	}
	// VALUES ((SELECT ... WHERE ...)) — each row is a *List whose Items hold
	// column expressions. ast.Inspect does not descend into *List itself, so
	// iterate the inner items explicitly.
	if sel.ValuesLists != nil {
		for _, row := range sel.ValuesLists.Items {
			if rl, ok := row.(*ast.List); ok {
				for _, item := range rl.Items {
					r.recurseIntoNestedSelectsInNode(item, outer)
				}
			} else {
				r.recurseIntoNestedSelectsInNode(row, outer)
			}
		}
	}
	r.recurseIntoNestedSelectsInNode(sel.HavingClause, outer)
	if sel.GroupClause != nil {
		for _, g := range sel.GroupClause.Items {
			r.recurseIntoNestedSelectsInNode(g, outer)
		}
	}
	if sel.SortClause != nil {
		for _, s := range sel.SortClause.Items {
			if sb, ok := s.(*ast.SortBy); ok {
				r.recurseIntoNestedSelectsInNode(sb.Node, outer)
			} else {
				r.recurseIntoNestedSelectsInNode(s, outer)
			}
		}
	}
}

// recurseIntoNestedSelectsInFromItem walks a FROM-clause item for nested
// subqueries — derived tables (RangeSubselect) and bare SelectStmts.
// JoinExpr.Quals are NOT descended into here, because checkJoinConditions
// already recurses into subqueries inside ON predicates; descending here too
// would double-flag violations inside correlated subqueries.
func (r *whereFunctionRule) recurseIntoNestedSelectsInFromItem(n ast.Node, outer whereTablesByAlias) {
	if n == nil {
		return
	}
	switch x := n.(type) {
	case *ast.JoinExpr:
		r.recurseIntoNestedSelectsInFromItem(x.Larg, outer)
		r.recurseIntoNestedSelectsInFromItem(x.Rarg, outer)
	case *ast.RangeSubselect:
		if sub, ok := x.Subquery.(*ast.SelectStmt); ok {
			r.checkSelectInScope(sub, outer)
		}
	default:
		r.recurseIntoNestedSelectsInNode(n, outer)
	}
}

// recurseIntoNestedSelectsInNode walks an arbitrary expression tree for
// nested subqueries, threading the given outer scope. Does not run the
// WHERE-detector on the expression itself.
func (r *whereFunctionRule) recurseIntoNestedSelectsInNode(n ast.Node, outer whereTablesByAlias) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(m ast.Node) bool {
		switch x := m.(type) {
		case *ast.SelectStmt:
			r.checkSelectInScope(x, outer)
			return false
		case *ast.SubLink:
			if sub, ok := x.Subselect.(*ast.SelectStmt); ok {
				r.checkSelectInScope(sub, outer)
			}
			if x.Testexpr != nil {
				r.recurseIntoNestedSelectsInNode(x.Testexpr, outer)
			}
			return false
		case *ast.RangeSubselect:
			if sub, ok := x.Subquery.(*ast.SelectStmt); ok {
				r.checkSelectInScope(sub, outer)
			}
			return false
		}
		return true
	})
}

// recurseIntoNestedSelectsInList walks each item of a List for nested
// subqueries, passing the given outer scope.
func (r *whereFunctionRule) recurseIntoNestedSelectsInList(list *ast.List, outer whereTablesByAlias) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		r.recurseIntoNestedSelectsInNode(item, outer)
	}
}

// ---- UPDATE / DELETE / INSERT / MERGE / CTAS / VIEW dispatch -------------

// collectTargetAndExtra collects the target relation plus the items of an
// extra clause (FROM for UPDATE, USING for DELETE).
func (r *whereFunctionRule) collectTargetAndExtra(relation ast.Node, extra *ast.List) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	if relation != nil {
		r.collectTableNamesFromNode(relation, tables)
	}
	if extra != nil {
		for _, item := range extra.Items {
			r.collectTableNamesFromNode(item, tables)
		}
	}
	return tables
}

// onUpdate handles UPDATE [WITH] … FROM … WHERE …
func (r *whereFunctionRule) onUpdate(n *ast.UpdateStmt) {
	if n == nil {
		return
	}
	popCTEs := r.handleWithClause(n.WithClause)
	defer popCTEs()
	tables := r.collectTargetAndExtra(n.Relation, n.FromClause)
	unresolved := r.localHasUnresolvedSources(n.FromClause, tables)
	r.checkWhereWithScopes(n.WhereClause, tables, nil, unresolved)
	r.checkJoinConditions(n.FromClause, tables, nil, unresolved)
	if n.FromClause != nil {
		for _, f := range n.FromClause.Items {
			r.recurseIntoNestedSelectsInFromItem(f, tables)
		}
	}
	r.recurseIntoNestedSelectsInList(n.TargetList, tables)
	r.recurseIntoNestedSelectsInList(n.ReturningList, tables)
}

// onDelete handles DELETE [WITH] … [USING …] WHERE …
func (r *whereFunctionRule) onDelete(n *ast.DeleteStmt) {
	if n == nil {
		return
	}
	popCTEs := r.handleWithClause(n.WithClause)
	defer popCTEs()
	tables := r.collectTargetAndExtra(n.Relation, n.UsingClause)
	unresolved := r.localHasUnresolvedSources(n.UsingClause, tables)
	r.checkWhereWithScopes(n.WhereClause, tables, nil, unresolved)
	r.checkJoinConditions(n.UsingClause, tables, nil, unresolved)
	if n.UsingClause != nil {
		for _, u := range n.UsingClause.Items {
			r.recurseIntoNestedSelectsInFromItem(u, tables)
		}
	}
	r.recurseIntoNestedSelectsInList(n.ReturningList, tables)
}

// onInsert handles INSERT [WITH] … SELECT / VALUES … [ON CONFLICT …].
func (r *whereFunctionRule) onInsert(n *ast.InsertStmt) {
	if n == nil {
		return
	}
	popCTEs := r.handleWithClause(n.WithClause)
	defer popCTEs()
	if sel, ok := n.SelectStmt.(*ast.SelectStmt); ok {
		r.checkSelectInScope(sel, nil)
	}
	if n.OnConflictClause != nil {
		target := make(whereTablesByAlias)
		if n.Relation != nil {
			r.collectTableNamesFromNode(n.Relation, target)
		}
		if n.OnConflictClause.WhereClause != nil {
			r.checkWhereWithScopes(n.OnConflictClause.WhereClause, target, nil, false)
		}
		// ON CONFLICT (col) WHERE <pred> — the Infer.WhereClause is the
		// partial-index inference predicate, a filter over the target table.
		if n.OnConflictClause.Infer != nil && n.OnConflictClause.Infer.WhereClause != nil {
			r.checkWhereWithScopes(n.OnConflictClause.Infer.WhereClause, target, nil, false)
		}
		r.recurseIntoNestedSelectsInList(n.OnConflictClause.TargetList, target)
	}
	target := make(whereTablesByAlias)
	if n.Relation != nil {
		r.collectTableNamesFromNode(n.Relation, target)
	}
	r.recurseIntoNestedSelectsInList(n.ReturningList, target)
}

func (r *whereFunctionRule) onViewStmt(n *ast.ViewStmt) {
	if n == nil {
		return
	}
	if sel, ok := n.Query.(*ast.SelectStmt); ok {
		r.checkSelectInScope(sel, nil)
	}
}

func (r *whereFunctionRule) onCreateTableAs(n *ast.CreateTableAsStmt) {
	if n == nil {
		return
	}
	if sel, ok := n.Query.(*ast.SelectStmt); ok {
		r.checkSelectInScope(sel, nil)
	}
}

// onMerge handles MERGE … USING … ON … WHEN ….
func (r *whereFunctionRule) onMerge(n *ast.MergeStmt) {
	if n == nil {
		return
	}
	popCTEs := r.handleWithClause(n.WithClause)
	defer popCTEs()
	tables := make(whereTablesByAlias)
	if n.Relation != nil {
		r.collectTableNamesFromNode(n.Relation, tables)
	}
	if n.SourceRelation != nil {
		r.collectTableNamesFromNode(n.SourceRelation, tables)
	}
	unresolved := false
	if n.SourceRelation != nil {
		ast.Inspect(n.SourceRelation, func(m ast.Node) bool {
			switch m.(type) {
			case *ast.SubLink, *ast.RangeSubselect, *ast.SelectStmt:
				unresolved = true
				return false
			}
			return !unresolved
		})
	}
	for _, ref := range tables {
		if ref == (whereTableRef{}) {
			continue
		}
		if r.tableMetadata(ref) == nil {
			unresolved = true
			break
		}
	}
	r.checkWhereWithScopes(n.JoinCondition, tables, nil, unresolved)
	// Source-side JOIN ON predicates (`USING a JOIN b ON UPPER(a.idx)=…`)
	// are bound to the MERGE's combined target+source scope.
	if n.SourceRelation != nil {
		r.checkJoinConditionsInNode(n.SourceRelation, tables, nil, unresolved)
	}
	// MERGE source may itself be a subquery (`USING (SELECT … WHERE …)`).
	// Use the FromItem walker so JoinExpr.Quals are skipped — those were
	// already handled by checkJoinConditionsInNode above.
	r.recurseIntoNestedSelectsInFromItem(n.SourceRelation, tables)
	if n.MergeWhenClauses != nil {
		for _, item := range n.MergeWhenClauses.Items {
			c, ok := item.(*ast.MergeWhenClause)
			if !ok {
				continue
			}
			r.checkWhereWithScopes(c.Condition, tables, nil, unresolved)
			r.recurseIntoNestedSelectsInList(c.TargetList, tables)
			r.recurseIntoNestedSelectsInList(c.Values, tables)
		}
	}
	r.recurseIntoNestedSelectsInList(n.ReturningList, tables)
}

// ---- Findings -------------------------------------------------------------

func (r *whereFunctionRule) addFunction(funcName, col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.FunctionOnIndexedColumn, Name: col, Detail: funcName, Start: start, End: end})
}

func (r *whereFunctionRule) addCalculation(col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.CalculationOnIndexedColumn, Name: col, Start: start, End: end})
}
