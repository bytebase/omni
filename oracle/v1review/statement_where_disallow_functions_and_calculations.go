package v1review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/v1review"
)

// StatementWhereDisallowFunctionsAndCalculations reports each WHERE-clause
// function call or calculation applied to an indexed column, judged against
// the pre-change schema.
//
// Judged: the WHERE of SELECT, UPDATE and DELETE, MERGE's ON and WHEN
// conditions, INSERT ALL/FIRST's WHEN conditions, START WITH, CONNECT BY and
// JOIN ON, in top-level statements, in CREATE TABLE AS SELECT and CREATE
// VIEW, in the query statements of an anonymous PL/SQL block, and in every
// SELECT nested in them (CTEs, set operations, subqueries, derived tables).
// Routine, trigger and package bodies and EXPLAIN PLAN are not judged. A
// function call is reported when an argument is an indexed column, bare or
// under PRIOR; any binary expression whose operator is not a comparison
// (+, -, *, /, ||, ...) and a unary minus or plus when any column beneath
// it is. Every such node beneath a judged expression is judged, including
// those in the select list of a subquery, against the scope of the
// expression's query. Name is the column as folded below; Value is the
// function's name without its qualifier, as the parser folded it.
//
// Names are compared as the parser folded them (unquoted upper-cased,
// quoted as written) for the table lookup, and upper-cased again for
// everything else (lower-cased when ObjectCaseSensitive is false). A table
// resolves by name in currentSchema of the snapshot, whatever its
// qualifier, unless the qualifier and currentDatabase are both set and
// differ so folded; a database link is ignored. A CTE name visible in scope
// hides the table of that name. A column is indexed when its folded name is
// among its table's folded index expressions, so a quoted mixed-case column
// matches its upper-cased spelling; an expression in double quotes is
// compared by its inner text, unfolded when ObjectCaseSensitive.
//
// The same node is reported once for each path that reaches it, so a
// finding may repeat; the V1 advisor printed one advice per text and line.
// Each finding anchors the reported function call or expression. Ported
// from Bytebase's V1 rule STATEMENT_WHERE_DISALLOW_FUNCTIONS_AND_CALCULATIONS.
func StatementWhereDisallowFunctionsAndCalculations(stmts []ast.Node, schema v1review.Schema, currentSchema, currentDatabase string) []v1review.Finding {
	r := &whereFunctionRule{
		currentSchema:         currentSchema,
		currentDatabase:       currentDatabase,
		schema:                schema,
		isObjectCaseSensitive: schema.ObjectCaseSensitive,
	}
	for _, stmt := range stmts {
		r.onStatement(stmt)
	}
	return r.findings
}

type whereTableRef struct {
	schema string
	name   string
}

type whereTablesByAlias = map[string]whereTableRef
type whereIndexedByAlias = map[string]map[string]bool

type whereScopeIndexed struct {
	local        whereIndexedByAlias
	localAll     whereIndexedByAlias
	localAliases map[string]bool
}

type whereFunctionRule struct {
	currentSchema         string
	currentDatabase       string
	schema                v1review.Schema
	isObjectCaseSensitive bool
	// cteStack holds the CTE names of each enclosing WITH. A FROM table of a
	// stacked name is recorded as opaque (whereTableRef{}).
	cteStack []map[string]bool
	findings []v1review.Finding
}

func (r *whereFunctionRule) pushCTEs(names map[string]bool) {
	r.cteStack = append(r.cteStack, names)
}

func (r *whereFunctionRule) popCTEs() {
	if n := len(r.cteStack); n > 0 {
		r.cteStack = r.cteStack[:n-1]
	}
}

func (r *whereFunctionRule) cteVisible(name string) bool {
	for _, frame := range r.cteStack {
		if frame[name] {
			return true
		}
	}
	return false
}

func (r *whereFunctionRule) onStatement(node ast.Node) {
	switch n := node.(type) {
	case *ast.SelectStmt:
		r.checkSelect(n, nil)
	case *ast.UpdateStmt:
		local := r.collectDMLTargetTables(n.Target)
		r.checkWhere(n.WhereClause, local)
		r.checkNestedSubqueries(n, local)
	case *ast.DeleteStmt:
		local := r.collectDMLTargetTables(n.Target)
		r.checkWhere(n.WhereClause, local)
		r.checkNestedSubqueries(n, local)
	case *ast.InsertStmt:
		if n.Select != nil {
			r.checkSelect(n.Select, nil)
		}
		if subquery, ok := n.Subquery.(*ast.SelectStmt); ok {
			r.checkSelect(subquery, nil)
			sourceTables := r.collectTables(subquery.FromClause)
			for _, item := range whereListItems(n.MultiTable) {
				clause, ok := item.(*ast.InsertIntoClause)
				if !ok {
					continue
				}
				r.checkWhere(clause.When, sourceTables)
			}
		}
		r.checkNestedSubqueries(n, nil)
	case *ast.MergeStmt:
		local := whereTablesByAlias{}
		if n.Target != nil {
			ref := whereTableRefFromObjectName(n.Target)
			local[""] = ref
			local[r.normalizeIdent(n.Target.Name)] = ref
			if n.TargetAlias != nil && n.TargetAlias.Name != "" {
				local[r.normalizeIdent(n.TargetAlias.Name)] = ref
			}
		}
		switch source := n.Source.(type) {
		case *ast.TableRef:
			if source.Name == nil {
				break
			}
			ref := whereTableRefFromObjectName(source.Name)
			local[r.normalizeIdent(ref.name)] = ref
			if source.Alias != nil && source.Alias.Name != "" {
				local[r.normalizeIdent(source.Alias.Name)] = ref
			}
			if n.SourceAlias != nil && n.SourceAlias.Name != "" {
				local[r.normalizeIdent(n.SourceAlias.Name)] = ref
			}
		case *ast.SubqueryRef:
			if selectStmt, ok := source.Subquery.(*ast.SelectStmt); ok {
				r.checkSelect(selectStmt, nil)
			}
			if source.Alias != nil && source.Alias.Name != "" {
				local[r.normalizeIdent(source.Alias.Name)] = whereTableRef{}
			}
			if n.SourceAlias != nil && n.SourceAlias.Name != "" {
				local[r.normalizeIdent(n.SourceAlias.Name)] = whereTableRef{}
			}
		default:
		}
		r.checkWhere(n.On, local)
		for _, item := range whereListItems(n.Clauses) {
			clause, ok := item.(*ast.MergeClause)
			if !ok {
				continue
			}
			r.checkWhere(clause.Condition, local)
			r.checkWhere(clause.UpdateWhere, local)
			r.checkWhere(clause.DeleteWhere, local)
			r.checkWhere(clause.InsertWhere, local)
			r.checkNestedSubqueries(clause, local)
		}
	case *ast.CreateTableStmt:
		if selectStmt, ok := n.AsQuery.(*ast.SelectStmt); ok {
			r.checkSelect(selectStmt, nil)
		}
	case *ast.CreateViewStmt:
		if selectStmt, ok := n.Query.(*ast.SelectStmt); ok {
			r.checkSelect(selectStmt, nil)
		}
	case *ast.PLSQLBlock:
		r.checkPLSQLBlock(n)
	default:
	}
}

func (r *whereFunctionRule) checkPLSQLBlock(block *ast.PLSQLBlock) {
	plsqlBlockStatements(block, func(stmt ast.StmtNode) bool {
		return !r.checkPLSQLQueryStatement(stmt)
	})
}

func (r *whereFunctionRule) checkPLSQLQueryStatement(stmt ast.StmtNode) bool {
	switch stmt.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.MergeStmt, *ast.CreateTableStmt, *ast.CreateViewStmt:
		r.onStatement(stmt)
		return true
	default:
		return false
	}
}

func (r *whereFunctionRule) collectDMLTargetTables(target ast.TableExpr) whereTablesByAlias {
	switch n := target.(type) {
	case *ast.TableRef:
		ref := whereTableRefFromObjectName(n.Name)
		tables := whereTablesByAlias{"": ref}
		if n.Name != nil {
			tables[r.normalizeIdent(n.Name.Name)] = ref
		}
		if n.Alias != nil && n.Alias.Name != "" {
			tables[r.normalizeIdent(n.Alias.Name)] = ref
		}
		return tables
	case *ast.SubqueryRef:
		if selectStmt, ok := n.Subquery.(*ast.SelectStmt); ok {
			r.checkSelect(selectStmt, nil)
		}
		if n.Alias != nil && n.Alias.Name != "" {
			return whereTablesByAlias{r.normalizeIdent(n.Alias.Name): whereTableRef{}}
		}
	default:
	}
	return nil
}

func (r *whereFunctionRule) checkNestedSubqueries(node ast.Node, outer whereTablesByAlias) {
	whereWalk(node, func(n ast.Node) {
		switch subquery := n.(type) {
		case *ast.SubqueryExpr:
			if selectStmt, ok := subquery.Subquery.(*ast.SelectStmt); ok {
				r.checkSelect(selectStmt, outer)
			}
		case *ast.ExistsExpr:
			if selectStmt, ok := subquery.Subquery.(*ast.SelectStmt); ok {
				r.checkSelect(selectStmt, outer)
			}
		default:
		}
	})
}

func whereTableRefFromObjectName(name *ast.ObjectName) whereTableRef {
	if name == nil {
		return whereTableRef{}
	}
	return whereTableRef{schema: name.Schema, name: name.Name}
}

func (r *whereFunctionRule) checkSelect(stmt *ast.SelectStmt, outer whereTablesByAlias) {
	if stmt == nil {
		return
	}
	if stmt.WithClause != nil {
		names := make(map[string]bool)
		r.pushCTEs(names)
		defer r.popCTEs()
		for _, item := range whereListItems(stmt.WithClause.CTEs) {
			cte, ok := item.(*ast.CTE)
			if !ok {
				continue
			}
			cteName := r.normalizeIdent(cte.Name)
			if r.queryIsRecursiveCTE(cte.Query, cteName) {
				names[cteName] = true
			}
			if query, ok := cte.Query.(*ast.SelectStmt); ok {
				r.checkSelect(query, outer)
			}
			names[cteName] = true
		}
	}
	if stmt.Larg != nil {
		r.checkSelect(stmt.Larg, outer)
	}
	if stmt.Rarg != nil {
		r.checkSelect(stmt.Rarg, outer)
	}
	local := r.collectTables(stmt.FromClause)
	tables := whereMergeTableAliases(outer, local)
	r.checkWhere(stmt.WhereClause, tables)
	if stmt.Hierarchical != nil {
		r.checkWhere(stmt.Hierarchical.StartWith, tables)
		r.checkWhere(stmt.Hierarchical.ConnectBy, tables)
	}
	for _, item := range whereListItems(stmt.FromClause) {
		r.checkJoinPredicates(item, tables)
	}
	whereWalk(stmt, func(n ast.Node) {
		subquery, ok := n.(*ast.SubqueryExpr)
		if ok {
			if nested, ok := subquery.Subquery.(*ast.SelectStmt); ok {
				r.checkSelect(nested, whereMergeTableAliases(outer, local))
			}
		}
	})
}

func (r *whereFunctionRule) queryIsRecursiveCTE(node ast.Node, tableName string) bool {
	selectStmt, ok := node.(*ast.SelectStmt)
	if !ok || (selectStmt.Larg == nil && selectStmt.Rarg == nil) {
		return false
	}
	found := false
	whereWalk(node, func(n ast.Node) {
		if found {
			return
		}
		table, ok := n.(*ast.TableRef)
		if !ok || table.Name == nil {
			return
		}
		found = r.normalizeIdent(table.Name.Name) == tableName
	})
	return found
}

func (r *whereFunctionRule) checkJoinPredicates(node ast.Node, tables whereTablesByAlias) {
	join, ok := node.(*ast.JoinClause)
	if !ok {
		return
	}
	r.checkWhere(join.On, tables)
	if left, ok := join.Left.(ast.Node); ok {
		r.checkJoinPredicates(left, tables)
	}
	if right, ok := join.Right.(ast.Node); ok {
		r.checkJoinPredicates(right, tables)
	}
}

func (r *whereFunctionRule) collectTables(from *ast.List) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	for _, item := range whereListItems(from) {
		switch n := item.(type) {
		case *ast.TableRef:
			ref := whereTableRefFromObjectName(n.Name)
			if r.cteVisible(r.normalizeIdent(ref.name)) {
				tables[r.normalizeIdent(ref.name)] = whereTableRef{}
				if n.Alias != nil && n.Alias.Name != "" {
					tables[r.normalizeIdent(n.Alias.Name)] = whereTableRef{}
				}
				continue
			}
			if ref.name == "" {
				continue
			}
			tables[r.normalizeIdent(ref.name)] = ref
			if n.Alias != nil && n.Alias.Name != "" {
				tables[r.normalizeIdent(n.Alias.Name)] = ref
			}
			if len(tables) == 1 {
				tables[""] = ref
			}
		case *ast.JoinClause:
			for alias, ref := range r.collectJoinTables(n) {
				tables[alias] = ref
			}
		default:
		}
	}
	return tables
}

func (r *whereFunctionRule) collectJoinTables(join *ast.JoinClause) whereTablesByAlias {
	tables := make(whereTablesByAlias)
	for _, item := range []ast.TableExpr{join.Left, join.Right} {
		switch n := item.(type) {
		case *ast.TableRef:
			ref := whereTableRefFromObjectName(n.Name)
			if r.cteVisible(r.normalizeIdent(ref.name)) {
				tables[r.normalizeIdent(ref.name)] = whereTableRef{}
				if n.Alias != nil && n.Alias.Name != "" {
					tables[r.normalizeIdent(n.Alias.Name)] = whereTableRef{}
				}
				continue
			}
			tables[r.normalizeIdent(ref.name)] = ref
			if n.Alias != nil && n.Alias.Name != "" {
				tables[r.normalizeIdent(n.Alias.Name)] = ref
			}
		case *ast.JoinClause:
			for alias, ref := range r.collectJoinTables(n) {
				tables[alias] = ref
			}
		default:
		}
	}
	return tables
}

func whereMergeTableAliases(a, b whereTablesByAlias) whereTablesByAlias {
	merged := make(whereTablesByAlias)
	for k, v := range a {
		merged[k] = v
	}
	for k, v := range b {
		merged[k] = v
	}
	return merged
}

func (r *whereFunctionRule) checkWhere(expr ast.ExprNode, tables whereTablesByAlias) {
	if expr == nil || len(tables) == 0 {
		return
	}
	indexed := r.resolveIndexedColumns(tables)
	allCols := r.resolveAllColumns(tables)
	if len(indexed) == 0 {
		return
	}
	r.walkWhere(expr, whereScopeIndexed{local: indexed, localAll: allCols, localAliases: map[string]bool{}}, tables)
}

func (r *whereFunctionRule) walkWhere(expr ast.ExprNode, scope whereScopeIndexed, tables whereTablesByAlias) {
	switch n := expr.(type) {
	case *ast.FuncCallExpr:
		if col := r.findIndexedColumnInExprList(n.Args, scope); col != "" {
			r.addFunctionAdvice(n.FuncName.Name, col, n.Loc)
		}
	case *ast.BinaryExpr:
		if n.Op != "=" && n.Op != "<>" && n.Op != "!=" && n.Op != "<" && n.Op != ">" && n.Op != "<=" && n.Op != ">=" {
			if col := r.findIndexedColumnInExpr(n, scope); col != "" {
				r.addCalculationAdvice(col, n.Loc)
			}
		}
	case *ast.UnaryExpr:
		if n.Op == "-" || n.Op == "+" {
			if col := r.findIndexedColumnInExpr(n, scope); col != "" {
				r.addCalculationAdvice(col, n.Loc)
			}
		}
	case *ast.SubqueryExpr:
		if selectStmt, ok := n.Subquery.(*ast.SelectStmt); ok {
			r.checkSelect(selectStmt, tables)
		}
	case *ast.ExistsExpr:
		if selectStmt, ok := n.Subquery.(*ast.SelectStmt); ok {
			r.checkSelect(selectStmt, tables)
		}
	default:
	}
	whereWalk(expr, func(child ast.Node) {
		if child == expr {
			return
		}
		if e, ok := child.(ast.ExprNode); ok {
			switch e.(type) {
			case *ast.FuncCallExpr, *ast.BinaryExpr, *ast.UnaryExpr, *ast.SubqueryExpr, *ast.ExistsExpr:
				r.walkWhere(e, scope, tables)
			default:
			}
		}
	})
}

func (r *whereFunctionRule) findIndexedColumnInExprList(list *ast.List, scope whereScopeIndexed) string {
	for _, item := range whereListItems(list) {
		switch expr := item.(type) {
		case *ast.ColumnRef:
			qualifier := r.normalizeIdent(expr.Table)
			column := r.normalizeIdent(expr.Column)
			if r.isIndexedColumn(qualifier, column, scope) {
				return column
			}
		case *ast.UnaryExpr:
			if expr.Op == "PRIOR" {
				if col := r.findDirectIndexedColumnInExpr(expr.Operand, scope); col != "" {
					return col
				}
			}
		default:
		}
	}
	return ""
}

func (r *whereFunctionRule) findDirectIndexedColumnInExpr(expr ast.ExprNode, scope whereScopeIndexed) string {
	col, ok := expr.(*ast.ColumnRef)
	if !ok {
		return ""
	}
	qualifier := r.normalizeIdent(col.Table)
	column := r.normalizeIdent(col.Column)
	if r.isIndexedColumn(qualifier, column, scope) {
		return column
	}
	return ""
}

func (r *whereFunctionRule) findIndexedColumnInExpr(expr ast.ExprNode, scope whereScopeIndexed) string {
	var found string
	whereWalk(expr, func(n ast.Node) {
		if found != "" {
			return
		}
		col, ok := n.(*ast.ColumnRef)
		if !ok {
			return
		}
		qualifier := r.normalizeIdent(col.Table)
		column := r.normalizeIdent(col.Column)
		if r.isIndexedColumn(qualifier, column, scope) {
			found = column
		}
	})
	return found
}

func (*whereFunctionRule) isIndexedColumn(qualifier, column string, scope whereScopeIndexed) bool {
	if qualifier != "" {
		return scope.local[qualifier][column]
	}
	for alias, cols := range scope.local {
		if alias == "" {
			continue
		}
		if cols[column] {
			return true
		}
	}
	return scope.local[""][column]
}

func (r *whereFunctionRule) addFunctionAdvice(funcName, col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.FunctionOnIndexedColumn, Name: col, Value: funcName, Start: start, End: end})
}

func (r *whereFunctionRule) addCalculationAdvice(col string, loc ast.Loc) {
	start, end := anchor(loc)
	r.findings = append(r.findings, v1review.Finding{Code: v1review.CalculationOnIndexedColumn, Name: col, Start: start, End: end})
}

// normalizeIdent folds a name for comparison. The parser has already
// stripped the quotes of identifiers; index expressions may keep theirs.
func (r *whereFunctionRule) normalizeIdent(name string) string {
	if name == "" {
		return name
	}
	if strings.HasPrefix(name, "\"") && strings.HasSuffix(name, "\"") && len(name) >= 2 {
		inner := name[1 : len(name)-1]
		if r.isObjectCaseSensitive {
			return inner
		}
		return strings.ToLower(inner)
	}
	if r.isObjectCaseSensitive {
		return strings.ToUpper(name)
	}
	return strings.ToLower(name)
}

// tableMetadata returns the table ref names. The snapshot holds only the
// connected owner's objects, in currentSchema, so a ref qualified with
// another owner resolves to nothing.
func (r *whereFunctionRule) tableMetadata(ref whereTableRef) *metadata.TableMetadata {
	if r.schema.Metadata == nil {
		return nil
	}
	if ref.schema != "" && r.currentDatabase != "" &&
		r.normalizeIdent(ref.schema) != r.normalizeIdent(r.currentDatabase) {
		return nil
	}
	return r.schema.Table(r.currentSchema, ref.name)
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
			cols[r.normalizeIdent(expr)] = true
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

func whereListItems(list *ast.List) []ast.Node {
	if list == nil {
		return nil
	}
	return list.Items
}

func whereWalk(node ast.Node, visit func(ast.Node)) {
	ast.Inspect(node, func(n ast.Node) bool {
		visit(n)
		return true
	})
}
