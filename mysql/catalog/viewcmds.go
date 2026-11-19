package catalog

import (
	"errors"
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/deparse"
)

// createView applies CREATE VIEW. source is the text the statement was parsed
// from, or "" when there is none (see deparseViewSelect). checkBody refuses a
// view body the engine rejects (a `?` parameter marker); it is true on the
// CREATE VIEW path (Exec) and false on a DefineView sync, whose body the engine
// already accepted.
func (c *Catalog) createView(stmt *nodes.CreateViewStmt, source string, checkBody bool) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	// A view body is judged before the view's name is looked up on the Exec
	// path: mysql:8.0.46 answers 1054 for a SELECT naming a column no relation
	// of the body has, even when the name already belongs to a table (1050) or
	// a view. A synced view's body is the engine's already, so it is not judged
	// here.
	analyzedQuery, bodyErr := c.analyzeViewBody(stmt.Select)
	if checkBody && bodyErr != nil {
		return bodyErr
	}
	key := toLower(stmt.Name.Name)
	// Tables and views share the same namespace in MySQL.
	if _, exists := db.Tables[key]; exists {
		return errDupTable(stmt.Name.Name)
	}
	if _, exists := db.Views[key]; exists {
		if !stmt.OrReplace {
			return errDupTable(stmt.Name.Name)
		}
	}

	// A `?` parameter marker is not allowed in a view body; MySQL rejects
	// CREATE VIEW with one. The grammar permits it (it is a LIMIT operand), so
	// the rejection lives here, not in the parser, and must fail Exec rather
	// than be swallowed as the analyzer's other failures are.
	if checkBody && stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	columnMetadata, err := inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}

	algorithm := stmt.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	sqlSecurity := stmt.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}

	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Algorithm:       algorithm,
		Definer:         definer,
		SqlSecurity:     sqlSecurity,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

// analyzeViewBody analyzes a view body for the semantic IR lineage and SDL
// diff use, and reports the failures MySQL's CREATE VIEW reports too: a star
// with no table to expand (1096), and a target-list column in a SELECT with no
// FROM, which names no column at all (1054). The analyzer models less than the
// engine here — it cannot resolve a select alias a HAVING or ORDER BY reads, for
// one — so every other failure leaves a view without AnalyzedQuery rather than
// refusing a view the engine accepts.
func (c *Catalog) analyzeViewBody(sel *nodes.SelectStmt) (*Query, error) {
	if sel == nil {
		return nil, nil
	}
	if err := checkViewBodyStar(sel); err != nil {
		return nil, err
	}
	q, err := c.AnalyzeSelectStmt(sel)
	if err == nil {
		return q, nil
	}
	var refusal *Error
	if errors.As(err, &refusal) && refusal.Code == ErrNoSuchColumn && viewBodyMayNameNoColumn(sel) {
		return nil, refusal
	}
	return nil, nil
}

// checkViewBodyStar refuses a view body whose star has no table to expand. Every
// query block MySQL resolves is judged: the ones a set operation, a parenthesis,
// a derived table, a subquery expression and a CTE body reach. MySQL resolves a
// CTE only where the body reads it — `WITH c AS (SELECT *) SELECT 1` is a view it
// takes — so a CTE the body never names is left alone. In a block with no FROM, a
// bare star is 1096 "No tables used" and a qualified one is 1051 "Unknown table
// 't'", whose table the block reads nowhere. Both answer before the view's own
// name is looked up, so `CREATE VIEW a AS SELECT *` is 1096 even over an existing
// table a (1050). A FROM-less block naming no star is left alone: `CREATE VIEW v
// AS SELECT 1` and `... SELECT count(*)` are both views MySQL takes.
func checkViewBodyStar(sel *nodes.SelectStmt) error {
	if sel == nil {
		return nil
	}
	named := bodyTableNames(sel)
	var refusal error
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if refusal != nil {
			return false
		}
		switch e := n.(type) {
		case *nodes.CommonTableExpr:
			return e.Name != "" && named[toLower(e.Name)]
		case *nodes.SelectStmt:
			refusal = viewBlockStar(e)
		}
		return refusal == nil
	})
	return refusal
}

// bodyTableNames is every name the body reads in a FROM clause, at any depth,
// which is the set of CTEs MySQL resolves.
func bodyTableNames(sel *nodes.SelectStmt) map[string]bool {
	names := make(map[string]bool)
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.TableRef); ok && ref.Name != "" {
			names[toLower(ref.Name)] = true
		}
		return true
	})
	return names
}

// viewBlockStar reports the star of one query block that reads no table.
func viewBlockStar(sel *nodes.SelectStmt) error {
	if len(sel.From) != 0 {
		return nil
	}
	for _, item := range sel.TargetList {
		expr := item
		if target, ok := item.(*nodes.ResTarget); ok {
			expr = target.Val
		}
		switch e := expr.(type) {
		case *nodes.StarExpr:
			return errNoTablesUsed()
		case *nodes.ColumnRef:
			if e.Star {
				return errUnknownTableName(e.Table)
			}
		}
	}
	return nil
}

// viewBodyMayNameNoColumn reports whether a no-FROM body reads a column that
// therefore names nothing: a target expression, or a WHERE, outside a nested
// query's own scope. The alias-bearing clauses (HAVING, GROUP BY, ORDER BY) are
// left alone because the analyzer cannot resolve a select alias there, and a
// schema-qualified call is left alone because its argument may name a stored
// function's column. Measured on mysql:8.0.46.
func viewBodyMayNameNoColumn(sel *nodes.SelectStmt) bool {
	if sel == nil {
		return false
	}
	// A set operation's root carries no FROM or target list of its own; each
	// arm is a scope of its own.
	if sel.Left != nil || sel.Right != nil {
		return viewBodyMayNameNoColumn(sel.Left) || viewBodyMayNameNoColumn(sel.Right)
	}
	if len(sel.From) != 0 {
		return false
	}
	if exprNamesColumn(sel.Where) {
		return true
	}
	for _, item := range sel.TargetList {
		var val nodes.ExprNode
		switch target := item.(type) {
		case *nodes.ResTarget:
			val = target.Val
		case nodes.ExprNode:
			val = target
		default:
			continue
		}
		if exprNamesColumn(val) {
			return true
		}
	}
	return false
}

// exprNamesColumn reports whether expr reads a column of this scope. A nested
// query has a scope of its own, and a schema-qualified call may name a stored
// function's column.
func exprNamesColumn(expr nodes.ExprNode) bool {
	if expr == nil {
		return false
	}
	found := false
	nodes.Inspect(expr, func(n nodes.Node) bool {
		switch e := n.(type) {
		case *nodes.ColumnRef:
			if !e.Star {
				found = true
			}
			return false
		case *nodes.SubqueryExpr, *nodes.ExistsExpr:
			return false
		case *nodes.FuncCallExpr:
			if e.Schema != "" {
				return false
			}
		}
		return !found
	})
	return found
}

// containsParamMarker reports whether the subtree holds a `?` parameter
// marker, at any depth including nested subqueries. It also inspects the rows
// of INSERT ... VALUES and of a VALUES statement, which nodes.Inspect's walker
// skips, so a marker in a stored routine or trigger body cannot slip past.
func containsParamMarker(n nodes.Node) bool {
	found := false
	nodes.Inspect(n, func(x nodes.Node) bool {
		if found {
			return false
		}
		switch e := x.(type) {
		case *nodes.ParamMarkerExpr:
			found = true
		case *nodes.InsertStmt:
			found = rowsHaveParamMarker(e.Values)
		case *nodes.ValuesStmt:
			found = rowsHaveParamMarker(e.Rows)
		}
		return !found
	})
	return found
}

// rowsHaveParamMarker reports whether any expression in the rows of an
// INSERT ... VALUES or a VALUES statement holds a `?` marker. Those rows are
// not reached by nodes.Inspect, so containsParamMarker checks them here.
func rowsHaveParamMarker(rows [][]nodes.ExprNode) bool {
	for _, row := range rows {
		for _, v := range row {
			if v != nil && containsParamMarker(v) {
				return true
			}
		}
	}
	return false
}

// alterView applies ALTER VIEW; source is as for createView.
func (c *Catalog) alterView(stmt *nodes.AlterViewStmt, source string) error {
	db, err := c.resolveDatabase(stmt.Name.Schema)
	if err != nil {
		return err
	}
	key := toLower(stmt.Name.Name)
	// ALTER VIEW requires the view to exist, and answers by what holds the name:
	// mysql:8.0.46 says 1347 "'db.a' is not VIEW" for a table and 1146 "Table
	// 'db.a' doesn't exist" for nothing.
	if _, exists := db.Views[key]; !exists {
		if db.Tables[key] != nil {
			return errWrongObject(db.Name, stmt.Name.Name, "VIEW")
		}
		return errNoSuchTable(db.Name, stmt.Name.Name)
	}

	// A `?` parameter marker is not allowed in a view body; the engine rejects
	// ALTER VIEW with one just as it rejects CREATE VIEW. ALTER VIEW comes only
	// from Exec, so its body is always checked.
	if stmt.Select != nil && containsParamMarker(stmt.Select) {
		return &Error{
			Code:     1351,
			SQLState: "HY000",
			Message:  "View's SELECT contains a '?' (parameter marker)",
		}
	}

	// MySQL always sets a definer. Default to `root`@`%` when not specified.
	definer := stmt.Definer
	if definer == "" {
		definer = "`root`@`%`"
	}

	// The body is analyzed for semantic IR (used by lineage, SDL diff) and for
	// the one failure MySQL's ALTER VIEW reports too.
	analyzedQuery, err := c.analyzeViewBody(stmt.Select)
	if err != nil {
		return err
	}

	// Resolve, rewrite, and deparse the SELECT to produce canonical definition.
	definition, derivedCols := c.deparseViewSelect(stmt.Select, stmt.SelectText, source, db)

	// Use explicit column list if provided, otherwise derive from SELECT target list.
	hasExplicit := len(stmt.Columns) > 0
	viewCols := stmt.Columns
	if !hasExplicit {
		viewCols = derivedCols
	}
	columnMetadata, err := inferViewColumnMetadata(stmt.Select, viewCols, db)
	if err != nil {
		return err
	}

	algorithm := stmt.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	sqlSecurity := stmt.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}

	db.Views[key] = &View{
		Name:            stmt.Name.Name,
		Database:        db,
		Definition:      definition,
		Algorithm:       algorithm,
		Definer:         definer,
		SqlSecurity:     sqlSecurity,
		CheckOption:     stmt.CheckOption,
		Columns:         viewCols,
		ColumnMetadata:  columnMetadata,
		ExplicitColumns: hasExplicit,
		AnalyzedQuery:   analyzedQuery,
		IsUpdatable:     inferViewUpdatable(stmt.Select, algorithm),
	}
	return nil
}

type viewRelationInfo struct {
	table    *Table
	optional bool
}

type viewCollationDerivation int

const (
	viewDerivationExplicit viewCollationDerivation = iota
	viewDerivationImplicit
	viewDerivationCoercible
)

type viewCollationInfo struct {
	charset    string
	collation  string
	derivation viewCollationDerivation
	valid      bool
}

func inferViewColumnMetadata(sel *nodes.SelectStmt, names []string, db *Database) ([]ViewColumn, error) {
	cols := make([]ViewColumn, len(names))
	for i, name := range names {
		cols[i].Name = name
	}
	if sel == nil {
		return cols, nil
	}
	relations := map[string]viewRelationInfo{}
	for _, from := range sel.From {
		collectViewRelations(from, db, false, relations)
	}
	for i, target := range sel.TargetList {
		if i >= len(cols) {
			break
		}
		rt, ok := target.(*nodes.ResTarget)
		if !ok {
			continue
		}
		cols[i].Nullable = inferViewExprNullable(rt.Val, relations)
		if ci, err := inferViewExprCollation(rt.Val, relations); err != nil {
			return nil, err
		} else if ci.valid {
			cols[i].Charset = ci.charset
			cols[i].Collation = ci.collation
		}
	}
	if err := validateViewSelectCollations(sel, relations); err != nil {
		return nil, err
	}
	return cols, nil
}

func collectViewRelations(te nodes.TableExpr, db *Database, optional bool, out map[string]viewRelationInfo) {
	switch n := te.(type) {
	case *nodes.TableRef:
		if db == nil {
			return
		}
		tbl := db.GetTable(n.Name)
		if tbl == nil {
			return
		}
		key := n.Name
		if n.Alias != "" {
			key = n.Alias
		}
		out[toLower(key)] = viewRelationInfo{table: tbl, optional: optional}
	case *nodes.JoinClause:
		leftOptional := optional
		rightOptional := optional
		switch n.Type {
		case nodes.JoinLeft, nodes.JoinNaturalLeft:
			rightOptional = true
		case nodes.JoinRight, nodes.JoinNaturalRight:
			leftOptional = true
		}
		collectViewRelations(n.Left, db, leftOptional, out)
		collectViewRelations(n.Right, db, rightOptional, out)
	}
}

func inferViewExprNullable(expr nodes.ExprNode, relations map[string]viewRelationInfo) bool {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefNullable(e, relations)
	case *nodes.FuncCallExpr:
		if isViewStringMetadataNullableFunc(e.Name) {
			return true
		}
		for _, arg := range e.Args {
			if inferViewExprNullable(arg, relations) {
				return true
			}
		}
		return false
	case *nodes.ParenExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.BinaryExpr:
		return inferViewExprNullable(e.Left, relations) || inferViewExprNullable(e.Right, relations)
	case *nodes.UnaryExpr:
		return inferViewExprNullable(e.Operand, relations)
	case *nodes.WeightStringExpr:
		// Propagate like the plain WEIGHT_STRING(expr) function-call path.
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.ExtractExpr:
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.IntervalExpr:
		return inferViewExprNullable(e.Value, relations)
	case *nodes.CastExpr:
		// Covers plain casts and the WEIGHT_STRING(... AS BINARY(n)) desugar,
		// whose WeightStringExpr.Expr is a CastExpr (oracle 8.0.32:
		// IS_NULLABLE=YES for both over a nullable column).
		return inferViewExprNullable(e.Expr, relations)
	case *nodes.PositionExpr:
		// mysql:8.0: IS_NULLABLE is YES for position('x' in s) over a
		// nullable s and for position('x' in NULL), NO over a NOT NULL column.
		return inferViewExprNullable(e.Substr, relations) ||
			inferViewExprNullable(e.Str, relations)
	case *nodes.NullLit:
		return true
	default:
		return false
	}
}

func inferViewColumnRefNullable(cr *nodes.ColumnRef, relations map[string]viewRelationInfo) bool {
	if cr == nil || cr.Star {
		return true
	}
	if cr.Table != "" {
		if rel, ok := relations[toLower(cr.Table)]; ok {
			if col := rel.table.GetColumn(cr.Column); col != nil {
				return rel.optional || col.Nullable
			}
		}
		return true
	}
	found := false
	nullable := true
	for _, rel := range relations {
		col := rel.table.GetColumn(cr.Column)
		if col == nil {
			continue
		}
		if found {
			return true
		}
		found = true
		nullable = rel.optional || col.Nullable
	}
	return nullable
}

func validateViewSelectCollations(sel *nodes.SelectStmt, relations map[string]viewRelationInfo) error {
	if sel == nil {
		return nil
	}
	for _, expr := range []nodes.ExprNode{sel.Where, sel.Having} {
		if expr == nil {
			continue
		}
		if _, err := inferViewExprCollation(expr, relations); err != nil {
			return err
		}
	}
	for _, expr := range sel.GroupBy {
		if _, err := inferViewExprCollation(expr, relations); err != nil {
			return err
		}
	}
	for _, item := range sel.OrderBy {
		if item == nil {
			continue
		}
		if _, err := inferViewExprCollation(item.Expr, relations); err != nil {
			return err
		}
	}
	return nil
}

func inferViewExprCollation(expr nodes.ExprNode, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	switch e := expr.(type) {
	case *nodes.ColumnRef:
		return inferViewColumnRefCollation(e, relations), nil
	case *nodes.StringLit:
		charset := normalizeCharsetName(strings.TrimPrefix(toLower(e.Charset), "_"))
		if charset == "" {
			charset = "utf8mb4"
		}
		return newViewCollation(charset, "", viewDerivationCoercible), nil
	case *nodes.CollateExpr:
		child, err := inferViewExprCollation(e.Expr, relations)
		if err != nil {
			return viewCollationInfo{}, err
		}
		collation := toLower(e.Collation)
		charset := normalizeCharsetName(charsetForCollation(collation))
		if charset == "" && child.valid {
			charset = child.charset
		}
		return newViewCollation(charset, collation, viewDerivationExplicit), nil
	case *nodes.ConvertExpr:
		if e.Charset != "" {
			return newViewCollation(normalizeCharsetName(e.Charset), "", viewDerivationImplicit), nil
		}
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.FuncCallExpr:
		return inferViewFuncCollation(e, relations)
	case *nodes.ParenExpr:
		return inferViewExprCollation(e.Expr, relations)
	case *nodes.BinaryExpr:
		return inferViewBinaryCollation(e, relations)
	case *nodes.UnaryExpr:
		return inferViewExprCollation(e.Operand, relations)
	default:
		return viewCollationInfo{}, nil
	}
}

func inferViewFuncCollation(fn *nodes.FuncCallExpr, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	name := strings.ToLower(fn.Name)
	switch name {
	case "concat", "concat_ws":
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	case "repeat", "lpad", "rpad":
		if len(fn.Args) == 0 {
			return viewCollationInfo{}, nil
		}
		return inferViewExprCollation(fn.Args[0], relations)
	default:
		var out viewCollationInfo
		for _, arg := range fn.Args {
			ci, err := inferViewExprCollation(arg, relations)
			if err != nil {
				return viewCollationInfo{}, err
			}
			agg, err := aggregateViewCollations(out, ci, name)
			if err != nil {
				return viewCollationInfo{}, err
			}
			out = agg
		}
		return out, nil
	}
}

func inferViewBinaryCollation(expr *nodes.BinaryExpr, relations map[string]viewRelationInfo) (viewCollationInfo, error) {
	left, err := inferViewExprCollation(expr.Left, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	right, err := inferViewExprCollation(expr.Right, relations)
	if err != nil {
		return viewCollationInfo{}, err
	}
	op := binaryOpToString(expr.Op)
	if isViewCollationComparisonOp(expr.Op) {
		if _, err := aggregateViewCollationsForComparison(left, right, op); err != nil {
			return viewCollationInfo{}, err
		}
		return viewCollationInfo{}, nil
	}
	return aggregateViewCollations(left, right, op)
}

func inferViewColumnRefCollation(cr *nodes.ColumnRef, relations map[string]viewRelationInfo) viewCollationInfo {
	if cr == nil || cr.Star {
		return viewCollationInfo{}
	}
	if cr.Table != "" {
		if rel, ok := relations[toLower(cr.Table)]; ok {
			return viewColumnCollation(rel.table.GetColumn(cr.Column))
		}
		return viewCollationInfo{}
	}
	var out viewCollationInfo
	found := false
	for _, rel := range relations {
		ci := viewColumnCollation(rel.table.GetColumn(cr.Column))
		if !ci.valid {
			continue
		}
		if found {
			return viewCollationInfo{}
		}
		found = true
		out = ci
	}
	return out
}

func viewColumnCollation(col *Column) viewCollationInfo {
	if col == nil || col.Charset == "" || !isStringType(col.DataType) {
		return viewCollationInfo{}
	}
	return newViewCollation(col.Charset, col.Collation, viewDerivationImplicit)
}

func newViewCollation(charset, collation string, derivation viewCollationDerivation) viewCollationInfo {
	charset = normalizeCharsetName(charset)
	collation = toLower(collation)
	if charset == "" && collation != "" {
		charset = normalizeCharsetName(charsetForCollation(collation))
	}
	if collation == "" && charset != "" {
		collation = defaultCollationForCharset[toLower(charset)]
	}
	if charset == "" || collation == "" {
		return viewCollationInfo{}
	}
	return viewCollationInfo{
		charset:    charset,
		collation:  collation,
		derivation: derivation,
		valid:      true,
	}
}

func aggregateViewCollations(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid {
		return right, nil
	}
	if !right.valid {
		return left, nil
	}
	if strings.EqualFold(left.collation, right.collation) {
		return left, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation != right.derivation {
		if left.derivation < right.derivation {
			return left, nil
		}
		return right, nil
	}
	if isCharsetSuperset(left.charset, right.charset) {
		return left, nil
	}
	if isCharsetSuperset(right.charset, left.charset) {
		return right, nil
	}
	// MySQL 8.0 permits several same-charset string function mixes that the
	// comparison path rejects. Preserve the left side for view metadata.
	if strings.EqualFold(left.charset, right.charset) {
		return left, nil
	}
	return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
}

func aggregateViewCollationsForComparison(left, right viewCollationInfo, op string) (viewCollationInfo, error) {
	if !left.valid || !right.valid || strings.EqualFold(left.collation, right.collation) {
		if left.valid {
			return left, nil
		}
		return right, nil
	}
	if left.derivation == viewDerivationExplicit && right.derivation == viewDerivationExplicit {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	if left.derivation == right.derivation && strings.EqualFold(left.charset, right.charset) {
		return viewCollationInfo{}, errIllegalMixCollations(left, right, op)
	}
	return aggregateViewCollations(left, right, op)
}

func isViewCollationComparisonOp(op nodes.BinaryOp) bool {
	switch op {
	case nodes.BinOpEq, nodes.BinOpNe, nodes.BinOpLt, nodes.BinOpGt, nodes.BinOpLe, nodes.BinOpGe, nodes.BinOpNullSafeEq:
		return true
	default:
		return false
	}
}

func isCharsetSuperset(charset, other string) bool {
	charset = normalizeCharsetName(charset)
	other = normalizeCharsetName(other)
	if charset == other {
		return true
	}
	switch charset {
	case "utf8mb4":
		return other == "latin1" || other == "ascii" || other == "utf8mb3" || other == "utf8"
	case "utf8mb3":
		return other == "ascii"
	case "latin1":
		return other == "ascii"
	default:
		return false
	}
}

func errIllegalMixCollations(left, right viewCollationInfo, op string) error {
	return &Error{
		Code:     1267,
		SQLState: "HY000",
		Message: fmt.Sprintf("Illegal mix of collations (%s,%s) and (%s,%s) for operation '%s'",
			left.collation, viewDerivationName(left.derivation),
			right.collation, viewDerivationName(right.derivation), op),
	}
}

func viewDerivationName(derivation viewCollationDerivation) string {
	switch derivation {
	case viewDerivationExplicit:
		return "EXPLICIT"
	case viewDerivationImplicit:
		return "IMPLICIT"
	case viewDerivationCoercible:
		return "COERCIBLE"
	default:
		return "NONE"
	}
}

func isViewStringMetadataNullableFunc(name string) bool {
	switch strings.ToLower(name) {
	case "concat", "concat_ws", "ifnull", "coalesce":
		return true
	default:
		return false
	}
}

func inferViewUpdatable(sel *nodes.SelectStmt, algorithm string) bool {
	if sel == nil || strings.EqualFold(algorithm, "TEMPTABLE") {
		return false
	}
	if sel.DistinctKind != nodes.DistinctNone && sel.DistinctKind != nodes.DistinctAll {
		return false
	}
	if sel.SetOp != nodes.SetOpNone || len(sel.GroupBy) > 0 || sel.Having != nil {
		return false
	}
	hasAgg := false
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if hasAgg {
			return false
		}
		// A spaced gated keyword is a generic scalar call, not the aggregate, so
		// it does not make the view non-updatable (see analyzeFuncCall).
		if fn, ok := n.(*nodes.FuncCallExpr); ok && isAggregateFunc(fn.Name) && !fn.SpacedParen {
			hasAgg = true
			return false
		}
		return true
	})
	return !hasAgg
}

func (c *Catalog) dropView(stmt *nodes.DropViewStmt) error {
	for _, ref := range stmt.Views {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		if _, exists := db.Views[key]; !exists {
			if stmt.IfExists {
				continue
			}
			// A name no view holds is 1347 when a table holds it, and 1051
			// "Unknown table" otherwise (mysql:8.0.46).
			if db.Tables[key] != nil {
				return errWrongObject(db.Name, ref.Name, "VIEW")
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		delete(db.Views, key)
	}
	return nil
}

// deparseViewSelect resolves, rewrites, and deparses the SELECT AST for a view.
// If the AST is nil (parser didn't produce one), falls back to the raw SelectText.
// Returns the deparsed definition and the derived column names from the resolved
// SELECT target list (used when no explicit column list is specified).
//
// source is the text the statement was parsed from, which the AST's Loc
// offsets index into. An unaliased select item is named after its text in
// source, as the engine names it; with no source ("" from DefineView), the
// name is rebuilt from the AST.
func (c *Catalog) deparseViewSelect(sel *nodes.SelectStmt, rawText, source string, db *Database) (string, []string) {
	if sel == nil {
		return rawText, nil
	}

	// Build a TableLookup that resolves table names from this database.
	lookup := tableLookupForDB(db)

	// Determine the database charset for CAST resolution.
	charset := db.Charset
	if charset == "" {
		charset = c.defaultCharset
	}

	// Resolve: qualify columns, expand *, normalize JOINs.
	resolver := &deparse.Resolver{
		Lookup:         lookup,
		DefaultCharset: charset,
		Source:         source,
	}
	resolver.Resolve(sel)

	// Extract column names from the resolved target list.
	derivedCols := extractViewColumns(sel, lookup)

	// Rewrite: NOT folding, boolean context wrapping.
	deparse.RewriteSelectStmt(sel)

	// Deparse: AST → canonical SQL text.
	return deparse.DeparseSelect(sel), derivedCols
}

// extractViewColumns extracts column names from a resolved SELECT target list.
// This produces the column list that MySQL would derive for a view. lookup
// resolves the table of a TABLE query primary, which has no target list to
// read.
func extractViewColumns(sel *nodes.SelectStmt, lookup deparse.TableLookup) []string {
	if sel == nil {
		return nil
	}
	// The result columns live on the leftmost output SELECT. Unwrap any
	// parenthesized-query wrappers and walk set-operation left arms to reach it,
	// so a parenthesized or set-op view body (CREATE VIEW v AS (SELECT a, b) or
	// (SELECT 1 AS a) UNION (SELECT 2 AS a)) derives columns from that SELECT.
	// A TABLE query primary is SELECT * FROM t and a VALUES primary a row
	// constructor, so neither has a target list to read. Their columns are
	// derived by the same code the resolver uses for a CTE or derived table
	// over one, which searches the query's own CTE scope before lookup's
	// tables (container-verified: CREATE VIEW v AS TABLE t declares the
	// columns of t, CREATE VIEW v AS VALUES ROW(1,2) declares column_0 and
	// column_1, and WITH x AS (SELECT 1 AS a) TABLE x declares a).
	if primary := deparse.PrimaryOutputColumns(sel, lookup); primary != nil {
		cols := make([]string, 0, len(primary))
		for _, c := range primary {
			cols = append(cols, c.Name)
		}
		return cols
	}
	sel = nodes.LeftmostQueryLeaf(sel)
	var cols []string
	for _, target := range sel.TargetList {
		if rt, ok := target.(*nodes.ResTarget); ok {
			if rt.Name != "" {
				cols = append(cols, rt.Name)
			} else if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
				cols = append(cols, cr.Column)
			} else {
				cols = append(cols, nodeToSQL(rt.Val))
			}
		} else if cr, ok := target.(*nodes.ColumnRef); ok {
			cols = append(cols, cr.Column)
		} else {
			cols = append(cols, nodeToSQL(target))
		}
	}
	return cols
}

// tableLookupForDB returns a deparse.TableLookup function that resolves table
// and view names from the given database's Tables and Views maps.
func tableLookupForDB(db *Database) deparse.TableLookup {
	return func(tableName string) *deparse.ResolverTable {
		key := toLower(tableName)
		// Try tables first.
		tbl := db.Tables[key]
		if tbl != nil {
			cols := make([]deparse.ResolverColumn, len(tbl.Columns))
			for i, c := range tbl.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     c.Name,
					Position: c.Position,
				}
			}
			return &deparse.ResolverTable{
				Name:    tbl.Name,
				Columns: cols,
			}
		}
		// Fall back to views.
		v := db.Views[key]
		if v != nil {
			cols := make([]deparse.ResolverColumn, len(v.Columns))
			for i, colName := range v.Columns {
				cols[i] = deparse.ResolverColumn{
					Name:     colName,
					Position: i + 1,
				}
			}
			return &deparse.ResolverTable{
				Name:    v.Name,
				Columns: cols,
			}
		}
		return nil
	}
}

// ShowCreateView produces MySQL 8.0-compatible SHOW CREATE VIEW output.
// Returns "" if the database or view does not exist.
func (c *Catalog) ShowCreateView(database, name string) string {
	db := c.GetDatabase(database)
	if db == nil {
		return ""
	}
	v := db.Views[toLower(name)]
	if v == nil {
		return ""
	}
	return showCreateView(v)
}

// formatDefiner ensures the definer string is backtick-quoted per MySQL 8.0 format.
// Input can be: `root`@`%`, root@%, 'root'@'%', etc.
// Output: `root`@`%`
func formatDefiner(definer string) string {
	// If already formatted with backticks, return as-is.
	if strings.HasPrefix(definer, "`") && strings.Contains(definer, "@") {
		return definer
	}
	// Split on @
	parts := strings.SplitN(definer, "@", 2)
	if len(parts) == 1 {
		// No @ — just backtick-quote the whole thing.
		return "`" + strings.Trim(parts[0], "`'") + "`"
	}
	user := strings.Trim(parts[0], "`'")
	host := strings.Trim(parts[1], "`'")
	return fmt.Sprintf("`%s`@`%s`", user, host)
}

// showCreateView produces the SHOW CREATE VIEW output for a view.
// MySQL 8.0 format:
//
//	CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`%` SQL SECURITY DEFINER VIEW `view_name` AS select_statement
//	WITH CASCADED CHECK OPTION
func showCreateView(v *View) string {
	var b strings.Builder

	b.WriteString("CREATE")

	// ALGORITHM — MySQL 8.0 always shows ALGORITHM, defaults to UNDEFINED.
	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	b.WriteString(fmt.Sprintf(" ALGORITHM=%s", strings.ToUpper(algorithm)))

	// DEFINER — MySQL 8.0 always shows DEFINER with backtick-quoted user@host.
	if v.Definer != "" {
		b.WriteString(fmt.Sprintf(" DEFINER=%s", formatDefiner(v.Definer)))
	}

	// SQL SECURITY — MySQL 8.0 always shows SQL SECURITY, defaults to DEFINER.
	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	b.WriteString(fmt.Sprintf(" SQL SECURITY %s", strings.ToUpper(sqlSecurity)))

	// VIEW name
	b.WriteString(" VIEW " + backtickIDQ(v.Name))

	// Column list (only if explicitly specified by user in CREATE VIEW).
	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = backtickIDQ(c)
		}
		b.WriteString(fmt.Sprintf(" (%s)", strings.Join(cols, ",")))
	}

	// AS select_statement
	b.WriteString(" AS ")
	b.WriteString(v.Definition)

	// WITH CHECK OPTION
	if v.CheckOption != "" {
		b.WriteString(fmt.Sprintf(" WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption)))
	}

	return b.String()
}
