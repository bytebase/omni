package v1review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereDisallowFunctionsAndCalculations reports each WHERE-clause
// function call or calculation applied to an indexed column, judged against
// the pre-change schema.
//
// Judged: the WHERE of top-level SELECT, UPDATE and DELETE, of the SELECT in
// INSERT ... SELECT, CREATE TABLE ... SELECT, CREATE VIEW and ALTER VIEW,
// and of every SELECT nested anywhere in them (CTEs, set operations,
// parenthesized and derived tables, scalar and correlated subqueries,
// INSERT VALUES / SET / ON DUPLICATE KEY UPDATE values, UPDATE SET values).
// A function call is reported when one of its arguments, after peeling
// parentheses and CASTs, is an indexed column; an arithmetic or bitwise
// binary expression when its left, else its right, operand is; a unary minus
// or bitwise NOT when its operand is. Name is the column as written; Value
// is the function name upper-cased.
//
// A column is indexed when it appears, lower-cased, among the expressions of
// any index of its table. Tables resolve in the snapshot's schema "" by
// name; a table qualified with another database is never resolved. A
// qualified column resolves against its alias in the local FROM, else the
// enclosing queries'. An unqualified column resolves locally when any local
// table declares it, and is reported only when such a table indexes it;
// otherwise it falls back to the enclosing queries unless a local source's
// columns are unknown (a derived table, or a table missing from the
// snapshot). Two unaliased tables of one name from different databases make
// that name unresolvable. Aliases and table names compare case-insensitively
// unless the snapshot is present and ObjectCaseSensitive; column names
// always compare lower-cased.
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
	case *ast.CreateTableStmt:
		r.onCreateTable(n)
	case *ast.CreateViewStmt:
		r.onCreateView(n)
	case *ast.AlterViewStmt:
		r.onAlterView(n)
	case *ast.UpdateStmt:
		r.onUpdate(n)
	case *ast.DeleteStmt:
		r.onDelete(n)
	default:
	}
}

// onInsert handles INSERT [... SELECT | ... VALUES (...) | ... SET ... | ... ON DUPLICATE KEY UPDATE].
// Each of SELECT/VALUES/SET/ON DUPLICATE KEY can embed scalar subqueries with their own WHERE.
func (r *whereFunctionRule) onInsert(n *ast.InsertStmt) {
	if n.Select != nil {
		r.checkSelect(n.Select)
	}
	for _, row := range n.Values {
		for _, v := range row {
			r.recurseIntoNestedSelectsInNode(v, nil)
		}
	}
	r.recurseAssignmentValues(n.SetList, nil)
	r.recurseAssignmentValues(n.OnDuplicateKey, nil)
}

func (r *whereFunctionRule) onCreateTable(n *ast.CreateTableStmt) {
	if n.Select != nil {
		r.checkSelect(n.Select)
	}
}

func (r *whereFunctionRule) onCreateView(n *ast.CreateViewStmt) {
	if n.Select != nil {
		r.checkSelect(n.Select)
	}
}

func (r *whereFunctionRule) onAlterView(n *ast.AlterViewStmt) {
	if n.Select != nil {
		r.checkSelect(n.Select)
	}
}

// onUpdate handles UPDATE WHERE and any scalar subqueries hidden in SET expressions
// or derived tables in the table list (e.g. JOIN (SELECT ... WHERE ...) s).
func (r *whereFunctionRule) onUpdate(n *ast.UpdateStmt) {
	tables := r.collectTableNames(n.Tables)
	r.checkWhere(n.Where, tables)
	for _, te := range n.Tables {
		r.recurseIntoNestedSelectsInNode(te, tables)
	}
	r.recurseAssignmentValues(n.SetList, tables)
}

// onDelete handles DELETE WHERE, including multi-table DELETE ... USING and derived
// tables that carry their own WHERE clauses.
func (r *whereFunctionRule) onDelete(n *ast.DeleteStmt) {
	tables := r.collectTableNames(n.Tables)
	for alias, ref := range r.collectTableNames(n.Using) {
		tables[alias] = ref
	}
	r.checkWhere(n.Where, tables)
	for _, te := range n.Tables {
		r.recurseIntoNestedSelectsInNode(te, tables)
	}
	for _, te := range n.Using {
		r.recurseIntoNestedSelectsInNode(te, tables)
	}
}

// recurseAssignmentValues walks the RHS of each assignment (UPDATE SET, INSERT SET,
// INSERT ... ON DUPLICATE KEY UPDATE) for scalar subqueries.
func (r *whereFunctionRule) recurseAssignmentValues(assignments []*ast.Assignment, outer whereTablesByAlias) {
	for _, a := range assignments {
		if a != nil && a.Value != nil {
			r.recurseIntoNestedSelectsInNode(a.Value, outer)
		}
	}
}

func (r *whereFunctionRule) checkSelect(sel *ast.SelectStmt) {
	r.checkSelectInScope(sel, nil)
}

// checkSelectInScope is like checkSelect but also threads an outer table scope down
// so that correlated subqueries can resolve references to outer aliases. Local
// (inner) tables take precedence for unqualified columns per SQL name resolution;
// outer tables are consulted only for qualified references (e.g. outer_alias.col).
func (r *whereFunctionRule) checkSelectInScope(sel *ast.SelectStmt, outer whereTablesByAlias) {
	if sel == nil {
		return
	}
	// CTEs are not correlated with their referencing SELECT, so check them without outer scope.
	for _, cte := range sel.CTEs {
		if cte.Select != nil {
			r.checkSelect(cte.Select)
		}
	}
	if inner := sel.ParenSource; inner != nil {
		r.checkSelectInScope(inner, outer)
		return
	}
	if sel.SetOp != ast.SetOpNone {
		r.checkSelectInScope(sel.Left, outer)
		r.checkSelectInScope(sel.Right, outer)
		return
	}
	local := r.collectTableNames(sel.From)
	// Build a combined map for downstream recursion into nested subqueries, where
	// local aliases are "outer" relative to that deeper scope.
	combined := make(whereTablesByAlias, len(local)+len(outer))
	for alias, ref := range outer {
		combined[alias] = ref
	}
	for alias, ref := range local {
		combined[alias] = ref
	}
	r.checkWhereWithScopes(sel.Where, local, outer, r.localHasUnresolvedSources(sel.From, local))
	// Recurse into nested SelectStmts elsewhere in this SELECT, passing current combined
	// scope as outer so correlated references in scalar/derived subqueries resolve.
	r.recurseIntoNestedSelects(sel, combined)
}

// localHasUnresolvedSources reports whether the given FROM and local alias map
// contain sources whose columns we cannot enumerate from metadata — CTE refs,
// derived tables (subqueries in FROM), and plain tables absent from metadata.
// When true, we must not fall back to outer scope for unqualified columns,
// because the inner source likely owns the name.
func (r *whereFunctionRule) localHasUnresolvedSources(from []ast.TableExpr, local whereTablesByAlias) bool {
	for _, te := range from {
		derivedFound := false
		ast.Inspect(te, func(n ast.Node) bool {
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

// recurseIntoNestedSelects finds nested *ast.SelectStmt nodes in the non-WHERE parts of
// the given SELECT and recursively invokes checkSelectInScope on them with outer scope.
func (r *whereFunctionRule) recurseIntoNestedSelects(sel *ast.SelectStmt, outer whereTablesByAlias) {
	visit := func(n ast.Node) {
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
	for _, from := range sel.From {
		visit(from)
	}
	for _, t := range sel.TargetList {
		visit(t)
	}
	visit(sel.Having)
	for _, g := range sel.GroupBy {
		visit(g)
	}
	for _, o := range sel.OrderBy {
		if o != nil {
			visit(o.Expr)
		}
	}
}

// whereTableRef captures the real database and table name for an aliased table reference.
// In MySQL terminology, TableRef.Schema is the database name.
type whereTableRef struct {
	database string // empty means current database
	name     string
}

// whereTablesByAlias maps alias (or name if no alias) to the underlying whereTableRef.
type whereTablesByAlias = map[string]whereTableRef

// collectTableNames extracts table refs (with database qualifier) from FROM/table expressions.
// Descent is stopped at subquery boundaries so that derived tables don't leak their
// inner TableRefs into the outer scope, but the derived table's own alias is captured
// as an opaque whereTableRef{} entry so qualifier-based resolution can still tell it apart
// from an outer alias of the same name.
// Alias keys are normalized per IsObjectCaseSensitive.
func (r *whereFunctionRule) collectTableNames(tableExprs []ast.TableExpr) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	for _, from := range tableExprs {
		ast.Inspect(from, func(n ast.Node) bool {
			// Stop at subquery boundaries — inner SELECTs are a separate scope — but
			// record the derived table's alias so qualifier lookups recognize it as local.
			if sq, ok := n.(*ast.SubqueryExpr); ok {
				if sq.Alias != "" {
					if _, exists := tables[r.normalizeIdent(sq.Alias)]; !exists {
						tables[r.normalizeIdent(sq.Alias)] = whereTableRef{}
					}
				}
				return false
			}
			if _, ok := n.(*ast.SelectStmt); ok {
				return false
			}
			t, ok := n.(*ast.TableRef)
			if !ok {
				return true
			}
			key := t.Alias
			if key == "" {
				key = t.Name
			}
			normKey := r.normalizeIdent(key)
			newRef := whereTableRef{database: t.Schema, name: t.Name}
			// If an unaliased same-name from a different database already occupies this
			// key (rare: `FROM db1.orders, db2.orders` without aliases), mark as ambiguous
			// so isIndexed returns false rather than resolving to the wrong table.
			// Users should alias these tables to disambiguate.
			if existing, ok := tables[normKey]; ok && (existing.database != newRef.database || existing.name != newRef.name) {
				tables[normKey] = whereTableRef{}
			} else {
				tables[normKey] = newRef
			}
			return false
		})
	}
	return tables
}

// normalizeIdent lowercases identifiers when the DB is configured case-insensitive.
func (r *whereFunctionRule) normalizeIdent(name string) string {
	if r.schema.Metadata != nil && r.schema.ObjectCaseSensitive {
		return name
	}
	return strings.ToLower(name)
}

// tableMetadata fetches the store-model table for the given ref, applying the
// cross-DB guard. Returns nil when not resolvable.
func (r *whereFunctionRule) tableMetadata(ref whereTableRef) *metadata.TableMetadata {
	if r.schema.Metadata == nil {
		return nil
	}
	// Skip cross-database references; the snapshot holds only the current database.
	if ref.database != "" && !strings.EqualFold(ref.database, r.schema.Metadata.GetName()) {
		return nil
	}
	return r.schema.Table("", ref.name)
}

// indexedColumns returns indexed column names for the given table reference.
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

// whereIndexedByAlias maps alias (or table name if no alias) to the set of indexed columns on that table.
type whereIndexedByAlias = map[string]map[string]bool

// whereScopeIndexed carries per-alias column metadata for the local (inner) scope
// and any outer scope inherited from enclosing queries. Unqualified refs first
// resolve against local (any table that declares the column); only when no
// local table declares the column — and the full local scope is resolvable —
// do we fall back to outer. If any local source has unknown columns (CTEs,
// derived tables, tables without metadata), we do NOT fall back: inner scope
// still owns the name, we just can't prove indexability, so we stay silent
// rather than misattribute to an outer indexed column.
type whereScopeIndexed struct {
	local           whereIndexedByAlias // per-alias indexed columns (local scope)
	localAll        whereIndexedByAlias // per-alias full column sets (local scope)
	localAliases    map[string]bool     // every alias present in the local scope (including unresolved/derived)
	outer           whereIndexedByAlias // per-alias indexed columns (outer scope)
	localUnresolved bool                // true when the local scope has sources whose columns we can't enumerate
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
//   - Unqualified: if any local table declares the column, local wins (flag only if
//     that local table has it indexed). If no local table declares the column, fall
//     back to outer indexed columns (supports correlated subqueries).
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
		// local sources have unresolved metadata — collectTableNames records
		// every local alias (including derived tables) so this check is sound.
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
	// Combined table set for recursing into nested subqueries (each is a new scope
	// where current local+outer both become "outer").
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
		if col := r.findIndexedColumnInArgs(expr.Args, indexed); col != "" {
			r.addFunction(strings.ToUpper(expr.Name), col, expr.Loc)
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
		if expr.Op == ast.UnaryMinus || expr.Op == ast.UnaryBitNot {
			if col := r.unwrapAndFindIndexedColumn(expr.Operand, indexed); col != "" {
				r.addCalculation(col, expr.Loc)
			}
		}
	default:
	}
}

func (r *whereFunctionRule) addFunction(funcName, col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.FunctionOnIndexedColumn, Name: col, Value: funcName, Start: start, End: end})
}

func (r *whereFunctionRule) addCalculation(col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.CalculationOnIndexedColumn, Name: col, Start: start, End: end})
}

// findIndexedColumnInArgs searches function arguments for an indexed ColumnRef,
// descending through passthrough wrappers (ParenExpr, CastExpr) so that
// `UPPER((name))` and `UPPER(CAST(name AS CHAR))` are still flagged.
// BinaryExpr/UnaryExpr (arithmetic) are NOT crossed — those are handled by the
// calculation advice path, so `ABS(id + 1)` is reported once as a calculation.
func (r *whereFunctionRule) findIndexedColumnInArgs(args []ast.ExprNode, indexed whereScopeIndexed) string {
	for _, arg := range args {
		if col := r.unwrapAndFindIndexedColumn(arg, indexed); col != "" {
			return col
		}
	}
	return ""
}

// unwrapAndFindIndexedColumn peels ParenExpr/CastExpr wrappers off expr and, if the
// inner expression is an indexed ColumnRef, returns its name. Returns "" otherwise.
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
		default:
			return ""
		}
	}
	return ""
}

func isWhereArithmeticOp(op ast.BinaryOp) bool {
	switch op {
	case ast.BinOpAdd, ast.BinOpSub, ast.BinOpMul, ast.BinOpDiv, ast.BinOpMod,
		ast.BinOpDivInt, ast.BinOpBitAnd, ast.BinOpBitOr, ast.BinOpBitXor,
		ast.BinOpShiftLeft, ast.BinOpShiftRight:
		return true
	default:
		return false
	}
}
