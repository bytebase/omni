package ddl

// legacy_test.go contains legacy types and inference functions that are only needed by tests.
// Production code uses the new single-pipeline (analyzeSelectStmt → Query → exprType/exprTypMod/exprCollation).
// These types exist here to support existing test helpers (reverseConvertSelectStmt, makeViewStmt, etc.)
// that construct legacy SelectStmt objects and convert them to pgparser AST for ProcessUtility.

import "fmt"

// --- Test helpers ---

// newTestCatalogWithTable creates a catalog with a single test table "t" (id int4, name text, val int8).
func newTestCatalogWithTable() *Catalog {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
		{Name: "val", Type: TypeName{Name: "int8", TypeMod: -1}},
	}, nil, false), 'r')
	return c
}

// --- Legacy types (were in select.go and expr.go) ---

// CTEDef represents a CTE (WITH clause) definition for the legacy inference path.
type CTEDef struct {
	Name    string
	Query   *SelectStmt
	Aliases []string // column aliases (if specified)
}

// SelectStmt represents a SELECT statement or a set operation.
type SelectStmt struct {
	// Simple SELECT
	TargetList []*ResTarget
	From       []*FromItem
	CTEs       []*CTEDef // WITH clause CTEs

	// Set operation (UNION/INTERSECT/EXCEPT)
	Op    SetOpType // SetOpNone for simple SELECT
	Left  *SelectStmt
	Right *SelectStmt
}

// ResTarget represents a single entry in the SELECT target list.
type ResTarget struct {
	Name string // alias; empty = infer from expression
	Val  *Expr
}

// FromKind identifies the kind of FROM clause item.
type FromKind int

const (
	FromTable    FromKind = iota // plain table reference
	FromJoin                     // joined tables
	FromSubquery                 // subquery in FROM
	FromCTE                      // common table expression reference
)

// FromItem represents an item in the FROM clause.
type FromItem struct {
	Kind FromKind

	// FromTable
	Schema string
	Table  string
	Alias  string

	// FromJoin
	JoinType  JoinType
	JoinLeft  *FromItem
	JoinRight *FromItem

	// FromSubquery
	Subquery *SelectStmt
	SubAlias string

	// FromCTE
	CTEIndex int // index into SelectStmt.CTEs
}

// ExprKind identifies the kind of expression node.
type ExprKind int

const (
	ExprColumnRef    ExprKind = iota // table.column or column
	ExprStar                         // * or table.*
	ExprLiteral                      // typed literal
	ExprFuncCall                     // function(args...) including aggregates
	ExprOpExpr                       // binary/unary operator
	ExprTypeCast                     // expr::type or CAST(expr AS type)
	ExprNullConst                    // NULL
	ExprCaseExpr                     // CASE WHEN ... THEN ... ELSE ... END
	ExprCoalesceExpr                 // COALESCE(a, b, ...)
	ExprSubquery                     // scalar (SELECT ...)
	ExprBoolExpr                     // AND/OR/NOT
)

// Expr represents an expression node in the AST.
// It uses a flat struct with a Kind discriminator.
type Expr struct {
	Kind ExprKind

	// ExprColumnRef
	TableName  string // optional qualifier
	ColumnName string

	// ExprStar
	StarTable string // empty = all tables

	// ExprLiteral
	LiteralType uint32 // UNKNOWNOID for untyped string, INT4OID for int, etc.

	// ExprFuncCall
	FuncName string
	FuncArgs []*Expr
	FuncStar bool // for count(*)

	// ExprOpExpr
	OpName string
	Left   *Expr // nil for prefix operator
	Right  *Expr

	// ExprTypeCast
	CastArg  *Expr
	CastType TypeName

	// ExprCaseExpr
	CaseResults []*Expr // THEN result expressions
	CaseElse    *Expr   // may be nil (→ NULL)

	// ExprCoalesceExpr
	CoalesceArgs []*Expr

	// ExprSubquery
	Subquery *SelectStmt

	// ExprBoolExpr
	BoolOp   BoolOpType
	BoolArgs []*Expr
}

// --- Legacy inference engine (was in infer.go) ---

// scope tracks available columns from the FROM clause.
type scope struct {
	tables []scopeTable
}

type scopeTable struct {
	alias   string
	columns []*Column
}

// buildScope constructs a scope from a list of FROM items.
func (c *Catalog) buildScope(from []*FromItem) (*scope, error) {
	return c.buildScopeWithCTEs(from, nil)
}

// buildScopeWithCTEs constructs a scope from a list of FROM items, with CTE context.
func (c *Catalog) buildScopeWithCTEs(from []*FromItem, ctes []*CTEDef) (*scope, error) {
	s := &scope{}
	for _, fi := range from {
		if err := c.addFromItemWithCTEs(s, fi, ctes); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (c *Catalog) addFromItem(s *scope, fi *FromItem) error {
	return c.addFromItemWithCTEs(s, fi, nil)
}

func (c *Catalog) addFromItemWithCTEs(s *scope, fi *FromItem, ctes []*CTEDef) error {
	switch fi.Kind {
	case FromTable:
		// Check if this references a CTE (unqualified name only).
		if fi.Schema == "" && ctes != nil {
			for idx, cteDef := range ctes {
				if cteDef.Name == fi.Table {
					// Redirect to CTE handling.
					cteFI := &FromItem{Kind: FromCTE, Table: fi.Table, Alias: fi.Alias, CTEIndex: idx}
					return c.addFromItemWithCTEs(s, cteFI, ctes)
				}
			}
		}
		_, rel, err := c.findRelation(fi.Schema, fi.Table)
		if err != nil {
			return err
		}
		alias := fi.Alias
		if alias == "" {
			alias = fi.Table
		}
		s.tables = append(s.tables, scopeTable{alias: alias, columns: rel.Columns})
		return nil

	case FromJoin:
		if err := c.addFromItemWithCTEs(s, fi.JoinLeft, ctes); err != nil {
			return err
		}
		return c.addFromItemWithCTEs(s, fi.JoinRight, ctes)

	case FromSubquery:
		cols, err := c.inferSelect(fi.Subquery)
		if err != nil {
			return err
		}
		alias := fi.SubAlias
		synthCols := make([]*Column, len(cols))
		for i, rc := range cols {
			synthCols[i] = &Column{
				AttNum:    int16(i + 1),
				Name:      rc.Name,
				TypeOID:   rc.TypeOID,
				TypeMod:   rc.TypeMod,
				Collation: rc.Collation,
			}
		}
		s.tables = append(s.tables, scopeTable{alias: alias, columns: synthCols})
		return nil

	case FromCTE:
		if ctes == nil || fi.CTEIndex >= len(ctes) {
			return fmt.Errorf("CTE index %d out of range", fi.CTEIndex)
		}
		cteDef := ctes[fi.CTEIndex]
		cols, err := c.inferSelect(cteDef.Query)
		if err != nil {
			return err
		}
		// Apply column aliases if provided.
		if len(cteDef.Aliases) > 0 {
			for i := range cols {
				if i < len(cteDef.Aliases) {
					cols[i].Name = cteDef.Aliases[i]
				}
			}
		}
		alias := fi.Alias
		if alias == "" {
			alias = cteDef.Name
		}
		synthCols := make([]*Column, len(cols))
		for i, rc := range cols {
			synthCols[i] = &Column{
				AttNum:    int16(i + 1),
				Name:      rc.Name,
				TypeOID:   rc.TypeOID,
				TypeMod:   rc.TypeMod,
				Collation: rc.Collation,
			}
		}
		s.tables = append(s.tables, scopeTable{alias: alias, columns: synthCols})
		return nil

	default:
		return fmt.Errorf("unsupported FROM kind %d", fi.Kind)
	}
}

// resolveColumn resolves a column reference in the scope.
func (s *scope) resolveColumn(table, column string) (*Column, error) {
	if table != "" {
		for i := range s.tables {
			if s.tables[i].alias != table {
				continue
			}
			for _, col := range s.tables[i].columns {
				if col.Name == column {
					return col, nil
				}
			}
			return nil, errUndefinedColumn(column)
		}
		return nil, errUndefinedTable(table)
	}

	var found *Column
	for i := range s.tables {
		for _, col := range s.tables[i].columns {
			if col.Name == column {
				if found != nil {
					return nil, errAmbiguousColumn(column)
				}
				found = col
			}
		}
	}
	if found == nil {
		return nil, errUndefinedColumn(column)
	}
	return found, nil
}

// expandStar expands a star expression into result columns.
func (s *scope) expandStar(table string) ([]*ResultColumn, error) {
	if table != "" {
		for i := range s.tables {
			if s.tables[i].alias != table {
				continue
			}
			var out []*ResultColumn
			for _, col := range s.tables[i].columns {
				out = append(out, &ResultColumn{Name: col.Name, TypeOID: col.TypeOID, TypeMod: col.TypeMod, Collation: col.Collation})
			}
			return out, nil
		}
		return nil, errUndefinedTable(table)
	}

	var out []*ResultColumn
	for i := range s.tables {
		for _, col := range s.tables[i].columns {
			out = append(out, &ResultColumn{Name: col.Name, TypeOID: col.TypeOID, TypeMod: col.TypeMod, Collation: col.Collation})
		}
	}
	return out, nil
}

// inferExprType infers the result type of an expression.
// Returns (typeOID, typmod, collation, error).
func (c *Catalog) inferExprType(s *scope, expr *Expr) (uint32, int32, uint32, error) {
	switch expr.Kind {
	case ExprColumnRef:
		col, err := s.resolveColumn(expr.TableName, expr.ColumnName)
		if err != nil {
			return 0, -1, 0, err
		}
		return col.TypeOID, col.TypeMod, col.Collation, nil

	case ExprStar:
		return 0, -1, 0, &Error{Code: CodeFeatureNotSupported, Message: "star not allowed in expression context"}

	case ExprLiteral:
		return expr.LiteralType, -1, c.typeCollation(expr.LiteralType), nil

	case ExprNullConst:
		return UNKNOWNOID, -1, 0, nil

	case ExprBoolExpr:
		for _, arg := range expr.BoolArgs {
			oid, _, _, err := c.inferExprType(s, arg)
			if err != nil {
				return 0, -1, 0, err
			}
			if oid != BOOLOID && oid != UNKNOWNOID {
				if !c.CanCoerce(oid, BOOLOID, 'i') {
					return 0, -1, 0, errDatatypeMismatch(fmt.Sprintf(
						"argument of %s must be type boolean, not type %s",
						boolOpName(expr.BoolOp), c.typeName(oid),
					))
				}
			}
		}
		return BOOLOID, -1, 0, nil

	case ExprOpExpr:
		return c.inferOpExpr(s, expr)

	case ExprFuncCall:
		return c.inferFuncCall(s, expr)

	case ExprTypeCast:
		_, _, _, err := c.inferExprType(s, expr.CastArg)
		if err != nil {
			return 0, -1, 0, err
		}
		oid, typmod, err := c.ResolveType(expr.CastType)
		if err != nil {
			return 0, -1, 0, err
		}
		return oid, typmod, c.typeCollation(oid), nil

	case ExprCaseExpr:
		return c.inferCaseExpr(s, expr)

	case ExprCoalesceExpr:
		return c.inferCoalesceExpr(s, expr)

	case ExprSubquery:
		cols, err := c.inferSelect(expr.Subquery)
		if err != nil {
			return 0, -1, 0, err
		}
		if len(cols) != 1 {
			return 0, -1, 0, &Error{Code: CodeTooManyColumns, Message: "subquery must return only one column"}
		}
		return cols[0].TypeOID, cols[0].TypeMod, cols[0].Collation, nil

	default:
		return 0, -1, 0, fmt.Errorf("unsupported expression kind %d", expr.Kind)
	}
}

func (c *Catalog) inferOpExpr(s *scope, expr *Expr) (uint32, int32, uint32, error) {
	var leftOID uint32
	if expr.Left != nil {
		oid, _, _, err := c.inferExprType(s, expr.Left)
		if err != nil {
			return 0, -1, 0, err
		}
		leftOID = oid
	}

	rightOID, _, _, err := c.inferExprType(s, expr.Right)
	if err != nil {
		return 0, -1, 0, err
	}

	isPrefix := expr.Left == nil
	resultOID, err := c.resolveOp(expr.OpName, leftOID, rightOID, isPrefix)
	if err != nil {
		return 0, -1, 0, err
	}
	return resultOID, -1, c.typeCollation(resultOID), nil
}

func (c *Catalog) inferFuncCall(s *scope, expr *Expr) (uint32, int32, uint32, error) {
	argTypes := make([]uint32, len(expr.FuncArgs))
	for i, arg := range expr.FuncArgs {
		oid, _, _, err := c.inferExprType(s, arg)
		if err != nil {
			return 0, -1, 0, err
		}
		argTypes[i] = oid
	}

	retType, err := c.resolveFunc(expr.FuncName, argTypes, expr.FuncStar)
	if err != nil {
		return 0, -1, 0, err
	}
	return retType, -1, c.typeCollation(retType), nil
}

func (c *Catalog) inferCaseExpr(s *scope, expr *Expr) (uint32, int32, uint32, error) {
	var typeOIDs []uint32
	for _, result := range expr.CaseResults {
		oid, _, _, err := c.inferExprType(s, result)
		if err != nil {
			return 0, -1, 0, err
		}
		typeOIDs = append(typeOIDs, oid)
	}
	if expr.CaseElse != nil {
		oid, _, _, err := c.inferExprType(s, expr.CaseElse)
		if err != nil {
			return 0, -1, 0, err
		}
		typeOIDs = append(typeOIDs, oid)
	} else {
		typeOIDs = append(typeOIDs, UNKNOWNOID)
	}

	common, err := c.selectCommonTypeFromOIDs(typeOIDs, false)
	if err != nil {
		return 0, -1, 0, err
	}
	return common, -1, c.typeCollation(common), nil
}

func (c *Catalog) inferCoalesceExpr(s *scope, expr *Expr) (uint32, int32, uint32, error) {
	typeOIDs := make([]uint32, len(expr.CoalesceArgs))
	for i, arg := range expr.CoalesceArgs {
		oid, _, _, err := c.inferExprType(s, arg)
		if err != nil {
			return 0, -1, 0, err
		}
		typeOIDs[i] = oid
	}

	common, err := c.selectCommonTypeFromOIDs(typeOIDs, false)
	if err != nil {
		return 0, -1, 0, err
	}
	return common, -1, c.typeCollation(common), nil
}

// resolveOp resolves an operator to its result type.
func (c *Catalog) resolveOp(name string, leftType, rightType uint32, isPrefix bool) (uint32, error) {
	if isPrefix {
		// Prefix operator: left=0.
		if ops := c.operByKey[operKey{name: name, left: 0, right: rightType}]; len(ops) > 0 {
			return ops[0].Result, nil
		}
		// Try with UNKNOWN → common types.
		if rightType == UNKNOWNOID {
			for _, commonType := range []uint32{TEXTOID, INT4OID, FLOAT8OID} {
				if ops := c.operByKey[operKey{name: name, left: 0, right: commonType}]; len(ops) > 0 {
					return ops[0].Result, nil
				}
			}
		}
		return 0, errUndefinedFunction(name, []uint32{rightType})
	}

	// Binary operator: exact match.
	if ops := c.operByKey[operKey{name: name, left: leftType, right: rightType}]; len(ops) > 0 {
		return ops[0].Result, nil
	}

	// If one side is UNKNOWN, try substituting the other side's type.
	if leftType == UNKNOWNOID && rightType != UNKNOWNOID {
		if ops := c.operByKey[operKey{name: name, left: rightType, right: rightType}]; len(ops) > 0 {
			return ops[0].Result, nil
		}
	}
	if rightType == UNKNOWNOID && leftType != UNKNOWNOID {
		if ops := c.operByKey[operKey{name: name, left: leftType, right: leftType}]; len(ops) > 0 {
			return ops[0].Result, nil
		}
	}
	if leftType == UNKNOWNOID && rightType == UNKNOWNOID {
		// Both unknown: try text.
		if ops := c.operByKey[operKey{name: name, left: TEXTOID, right: TEXTOID}]; len(ops) > 0 {
			return ops[0].Result, nil
		}
	}

	// Try implicit coercion: find an operator with matching name where both args can be coerced.
	return c.resolveOpWithCoercion(name, leftType, rightType)
}

func (c *Catalog) resolveOpWithCoercion(name string, leftType, rightType uint32) (uint32, error) {
	// Collect all operators with this name.
	type candidate struct {
		op    *BuiltinOperator
		score int
	}
	var candidates []candidate

	for key, ops := range c.operByKey {
		if key.name != name {
			continue
		}
		for _, op := range ops {
			if op.Kind == 'l' {
				continue // skip prefix operators for binary resolution
			}
			leftOK := leftType == op.Left || c.CanCoerce(leftType, op.Left, 'i')
			rightOK := rightType == op.Right || c.CanCoerce(rightType, op.Right, 'i')
			if leftOK && rightOK {
				score := 0
				if leftType != op.Left {
					score++
				}
				if rightType != op.Right {
					score++
				}
				candidates = append(candidates, candidate{op: op, score: score})
			}
		}
	}

	if len(candidates) == 0 {
		return 0, errUndefinedFunction(name, []uint32{leftType, rightType})
	}

	// Pick best (lowest score).
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.score < best.score {
			best = c
		}
	}
	return best.op.Result, nil
}

// resolveFunc resolves a function call to its return type.
func (c *Catalog) resolveFunc(name string, argTypes []uint32, hasStar bool) (uint32, error) {
	// Special case: count(*).
	if hasStar && name == "count" {
		return INT8OID, nil
	}

	procs := c.procByName[name]
	if len(procs) == 0 {
		return 0, errUndefinedFunction(name, argTypes)
	}

	type candidate struct {
		proc  *BuiltinProc
		score int
	}
	var candidates []candidate

	for _, p := range procs {
		if int(p.NArgs) != len(argTypes) {
			continue
		}

		score := 0
		match := true
		for i, argOID := range argTypes {
			paramOID := p.ArgTypes[i]
			if argOID == paramOID {
				continue // exact match
			}
			if isPolymorphic(paramOID) {
				score++
				continue
			}
			if argOID == UNKNOWNOID {
				score++
				continue
			}
			if c.CanCoerce(argOID, paramOID, 'i') {
				score += 2
				continue
			}
			match = false
			break
		}
		if match {
			candidates = append(candidates, candidate{proc: p, score: score})
		}
	}

	if len(candidates) == 0 {
		return 0, errUndefinedFunction(name, argTypes)
	}

	// Pick best (lowest score).
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.score < best.score {
			best = c
		}
	}

	return c.resolveReturnType(best.proc, argTypes), nil
}

// inferSelect infers result column types from a SELECT statement.
func (c *Catalog) inferSelect(stmt *SelectStmt) ([]*ResultColumn, error) {
	if stmt.Op != SetOpNone {
		return c.inferSetOp(stmt)
	}
	return c.inferSimpleSelect(stmt)
}

func (c *Catalog) inferSetOp(stmt *SelectStmt) ([]*ResultColumn, error) {
	// Propagate CTEs from the outer set-op to child branches.
	if len(stmt.CTEs) > 0 {
		if stmt.Left != nil && len(stmt.Left.CTEs) == 0 {
			stmt.Left.CTEs = stmt.CTEs
		}
		if stmt.Right != nil && len(stmt.Right.CTEs) == 0 {
			stmt.Right.CTEs = stmt.CTEs
		}
	}
	leftCols, err := c.inferSelect(stmt.Left)
	if err != nil {
		return nil, err
	}
	rightCols, err := c.inferSelect(stmt.Right)
	if err != nil {
		return nil, err
	}

	if len(leftCols) != len(rightCols) {
		return nil, &Error{
			Code:    CodeDatatypeMismatch,
			Message: fmt.Sprintf("each %s query must have the same number of columns", setOpName(stmt.Op)),
		}
	}

	result := make([]*ResultColumn, len(leftCols))
	for i := range leftCols {
		common, err := c.selectCommonTypeFromOIDs([]uint32{leftCols[i].TypeOID, rightCols[i].TypeOID}, false)
		if err != nil {
			return nil, err
		}
		coll := selectCommonCollation(leftCols[i].Collation, rightCols[i].Collation)
		result[i] = &ResultColumn{
			Name:      leftCols[i].Name,
			TypeOID:   common,
			TypeMod:   -1,
			Collation: coll,
		}
	}
	return result, nil
}

func (c *Catalog) inferSimpleSelect(stmt *SelectStmt) ([]*ResultColumn, error) {
	s, err := c.buildScopeWithCTEs(stmt.From, stmt.CTEs)
	if err != nil {
		return nil, err
	}

	var result []*ResultColumn
	for _, rt := range stmt.TargetList {
		if rt.Val.Kind == ExprStar {
			expanded, err := s.expandStar(rt.Val.StarTable)
			if err != nil {
				return nil, err
			}
			result = append(result, expanded...)
			continue
		}

		oid, typmod, coll, err := c.inferExprType(s, rt.Val)
		if err != nil {
			return nil, err
		}

		name := rt.Name
		if name == "" {
			name = inferColumnName(rt.Val)
		}

		result = append(result, &ResultColumn{Name: name, TypeOID: oid, TypeMod: typmod, Collation: coll})
	}
	return result, nil
}

// inferColumnName infers a column name from an expression.
func inferColumnName(expr *Expr) string {
	switch expr.Kind {
	case ExprColumnRef:
		return expr.ColumnName
	case ExprFuncCall:
		return expr.FuncName
	default:
		return "?column?"
	}
}

func boolOpName(op BoolOpType) string {
	switch op {
	case BoolAnd:
		return "AND"
	case BoolOr:
		return "OR"
	case BoolNot:
		return "NOT"
	default:
		return "BOOL"
	}
}

func setOpName(op SetOpType) string {
	switch op {
	case SetOpUnion, SetOpUnionAll:
		return "UNION"
	case SetOpIntersect, SetOpIntersectAll:
		return "INTERSECT"
	case SetOpExcept, SetOpExceptAll:
		return "EXCEPT"
	default:
		return "SET"
	}
}
