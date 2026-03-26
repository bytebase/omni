// Package deparse converts MySQL AST nodes back to SQL text,
// matching MySQL 8.0's SHOW CREATE VIEW formatting.
package deparse

import (
	"fmt"
	"math/big"
	"strings"

	ast "github.com/bytebase/omni/mysql/ast"
)

// Deparse converts an expression AST node to its SQL text representation,
// matching MySQL 8.0's canonical formatting (as seen in SHOW CREATE VIEW).
func Deparse(node ast.ExprNode) string {
	if node == nil {
		return ""
	}
	return deparseExpr(node)
}

// DeparseSelect converts a SelectStmt AST node to its SQL text representation,
// matching MySQL 8.0's SHOW CREATE VIEW formatting.
func DeparseSelect(stmt *ast.SelectStmt) string {
	if stmt == nil {
		return ""
	}
	return deparseSelectStmt(stmt)
}

func deparseSelectStmt(stmt *ast.SelectStmt) string {
	return deparseSelectStmtCtx(stmt, false)
}

// deparseSelectStmtNoAlias formats a SELECT without target list aliases.
// Used for subquery contexts (IN, EXISTS) where MySQL omits AS alias.
func deparseSelectStmtNoAlias(stmt *ast.SelectStmt) string {
	return deparseSelectStmtCtx(stmt, true)
}

func deparseSelectStmtCtx(stmt *ast.SelectStmt, suppressAlias bool) string {
	// Handle set operations: UNION / UNION ALL / INTERSECT / EXCEPT
	if stmt.SetOp != ast.SetOpNone {
		return deparseSetOperation(stmt)
	}

	var b strings.Builder

	// CTE (WITH clause) — emit before the SELECT keyword
	if len(stmt.CTEs) > 0 {
		b.WriteString(deparseCTEs(stmt.CTEs))
		b.WriteString(" ")
	}

	b.WriteString("select ")

	// DISTINCT
	if stmt.DistinctKind == ast.DistinctOn {
		b.WriteString("distinct ")
	}

	// Target list
	for i, target := range stmt.TargetList {
		if i > 0 {
			b.WriteString(",")
		}
		if suppressAlias {
			b.WriteString(deparseResTargetNoAlias(target))
		} else {
			b.WriteString(deparseResTarget(target, i+1))
		}
	}

	// FROM clause
	if len(stmt.From) > 0 {
		b.WriteString(" from ")
		if len(stmt.From) == 1 {
			b.WriteString(deparseTableExpr(stmt.From[0]))
		} else {
			// Multiple tables (implicit cross join) → normalized to explicit join with parens
			// e.g., FROM t1, t2 → from (`t1` join `t2`)
			// For 3+ tables: FROM t1, t2, t3 → from ((`t1` join `t2`) join `t3`)
			b.WriteString(deparseImplicitCrossJoin(stmt.From))
		}
	}

	// WHERE clause
	if stmt.Where != nil {
		b.WriteString(" where ")
		b.WriteString(deparseExpr(stmt.Where))
	}

	// GROUP BY clause
	if len(stmt.GroupBy) > 0 {
		b.WriteString(" group by ")
		for i, expr := range stmt.GroupBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(expr))
		}
		if stmt.WithRollup {
			b.WriteString(" with rollup")
		}
	}

	// HAVING clause
	if stmt.Having != nil {
		b.WriteString(" having ")
		b.WriteString(deparseExpr(stmt.Having))
	}

	// ORDER BY clause
	if len(stmt.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range stmt.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(item.Expr))
			if item.Desc {
				b.WriteString(" desc")
			}
		}
	}

	// LIMIT clause
	if stmt.Limit != nil {
		b.WriteString(" limit ")
		if stmt.Limit.Offset != nil {
			// MySQL comma syntax: LIMIT offset,count
			b.WriteString(deparseExpr(stmt.Limit.Offset))
			b.WriteString(",")
		}
		b.WriteString(deparseExpr(stmt.Limit.Count))
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if stmt.ForUpdate != nil {
		b.WriteString(" ")
		b.WriteString(deparseForUpdate(stmt.ForUpdate))
	}

	return b.String()
}

// deparseForUpdate formats a FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE clause.
// MySQL 8.0 format:
//   - for update
//   - for share
//   - lock in share mode (legacy syntax)
//   - for update of `t`
//   - for update nowait
//   - for update skip locked
func deparseForUpdate(fu *ast.ForUpdate) string {
	if fu.LockInShareMode {
		return "lock in share mode"
	}

	var b strings.Builder
	if fu.Share {
		b.WriteString("for share")
	} else {
		b.WriteString("for update")
	}

	// OF table list
	if len(fu.Tables) > 0 {
		b.WriteString(" of ")
		for i, tbl := range fu.Tables {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString("`")
			b.WriteString(tbl.Name)
			b.WriteString("`")
		}
	}

	// NOWAIT / SKIP LOCKED
	if fu.NoWait {
		b.WriteString(" nowait")
	} else if fu.SkipLocked {
		b.WriteString(" skip locked")
	}

	return b.String()
}

// deparseSetOperation formats a set operation (UNION, INTERSECT, EXCEPT).
// MySQL 8.0 format: select ... union [all] select ... (flat, no parens around sub-selects)
// CTEs from the leftmost child are hoisted and emitted before the entire set operation.
func deparseSetOperation(stmt *ast.SelectStmt) string {
	// Hoist CTEs from the leftmost descendant
	var ctePrefix string
	if ctes := extractCTEs(stmt); len(ctes) > 0 {
		ctePrefix = deparseCTEs(ctes) + " "
	}

	left := deparseSelectStmt(stmt.Left)
	right := deparseSelectStmt(stmt.Right)

	var op string
	switch stmt.SetOp {
	case ast.SetOpUnion:
		if stmt.SetAll {
			op = "union all"
		} else {
			op = "union"
		}
	case ast.SetOpIntersect:
		if stmt.SetAll {
			op = "intersect all"
		} else {
			op = "intersect"
		}
	case ast.SetOpExcept:
		if stmt.SetAll {
			op = "except all"
		} else {
			op = "except"
		}
	}

	return ctePrefix + left + " " + op + " " + right
}

// extractCTEs walks down the left spine of a set operation tree and extracts
// CTEs from the leftmost leaf SelectStmt, clearing them so they aren't emitted
// again by deparseSelectStmt.
func extractCTEs(stmt *ast.SelectStmt) []*ast.CommonTableExpr {
	// Walk to the leftmost leaf
	cur := stmt
	for cur.SetOp != ast.SetOpNone && cur.Left != nil {
		cur = cur.Left
	}
	if len(cur.CTEs) > 0 {
		ctes := cur.CTEs
		cur.CTEs = nil // prevent double emission
		return ctes
	}
	return nil
}

// deparseCTEs formats a WITH clause (one or more CTEs).
// MySQL 8.0 format: with [recursive] `name` [(`col`, ...)] as (select ...) [, ...]
func deparseCTEs(ctes []*ast.CommonTableExpr) string {
	var b strings.Builder
	b.WriteString("with ")

	// Check if any CTE is recursive (the flag is per-CTE in the AST
	// but WITH RECURSIVE applies to the whole clause in SQL)
	recursive := false
	for _, cte := range ctes {
		if cte.Recursive {
			recursive = true
			break
		}
	}
	if recursive {
		b.WriteString("recursive ")
	}

	for i, cte := range ctes {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("`")
		b.WriteString(cte.Name)
		b.WriteString("`")

		// Column list
		if len(cte.Columns) > 0 {
			b.WriteString(" (")
			for j, col := range cte.Columns {
				if j > 0 {
					b.WriteString(",")
				}
				b.WriteString("`")
				b.WriteString(col)
				b.WriteString("`")
			}
			b.WriteString(")")
		}

		b.WriteString(" as (")
		if cte.Select != nil {
			b.WriteString(deparseSelectStmt(cte.Select))
		}
		b.WriteString(")")
	}

	return b.String()
}

// deparseResTargetNoAlias formats a target list entry without the AS alias.
// Used for subquery contexts (IN, EXISTS) where MySQL omits aliases.
func deparseResTargetNoAlias(node ast.ExprNode) string {
	if rt, ok := node.(*ast.ResTarget); ok {
		return deparseExpr(rt.Val)
	}
	return deparseExpr(node)
}

// deparseResTarget formats a single result target in the SELECT list.
// MySQL 8.0 SHOW CREATE VIEW format: expr AS `alias`
// - Always uses AS keyword
// - Alias is always backtick-quoted
// - Auto-alias: column ref → column name; literal → literal text; expression → expression text
func deparseResTarget(node ast.ExprNode, position int) string {
	rt, isRT := node.(*ast.ResTarget)

	var expr ast.ExprNode
	var explicitAlias string
	if isRT {
		expr = rt.Val
		explicitAlias = rt.Name
	} else {
		expr = node
	}

	exprStr := deparseExpr(expr)

	// Determine alias
	alias := explicitAlias
	if alias == "" {
		alias = autoAlias(expr, exprStr, position)
	}

	// MySQL 8.0 uses double-space before AS for window function expressions.
	// The OVER clause already ends with " )", and MySQL adds an extra space.
	if hasWindowFunction(expr) {
		return exprStr + "  AS `" + alias + "`"
	}
	return exprStr + " AS `" + alias + "`"
}

// hasWindowFunction checks if an expression is or contains a window function (has OVER clause).
func hasWindowFunction(node ast.ExprNode) bool {
	if fc, ok := node.(*ast.FuncCallExpr); ok && fc.Over != nil {
		return true
	}
	return false
}

// autoAlias generates an automatic alias for a SELECT target expression.
// MySQL 8.0 rules:
// - Column ref → column name (unqualified)
// - Literal → literal text representation
// - Short expression → expression text (without backtick quoting)
// - Long/complex expression → Name_exp_N
func autoAlias(expr ast.ExprNode, exprStr string, position int) string {
	switch n := expr.(type) {
	case *ast.ColumnRef:
		return n.Column
	case *ast.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *ast.FloatLit:
		return n.Value
	case *ast.StringLit:
		if n.Value == "" {
			return fmt.Sprintf("Name_exp_%d", position)
		}
		return n.Value
	case *ast.NullLit:
		return "NULL"
	case *ast.BoolLit:
		if n.Value {
			return "TRUE"
		}
		return "FALSE"
	case *ast.HexLit:
		// MySQL 8.0 preserves original literal form in auto-alias.
		// 0xFF stays as "0xFF"; X'FF' form stored as "FF" → "X'FF'"
		val := n.Value
		if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
			return val // preserve original case: 0xFF
		}
		return "X'" + val + "'"
	case *ast.BitLit:
		// MySQL 8.0 preserves original literal form in auto-alias.
		// 0b1010 stays as "0b1010"; b'1010' form stored as "1010" → "b'1010'"
		val := n.Value
		if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
			return val
		}
		return "b'" + val + "'"
	case *ast.TemporalLit:
		return n.Type + " '" + n.Value + "'"
	default:
		// For expressions: generate a human-readable alias text without backtick quoting.
		// MySQL 8.0 uses the original expression text for the alias.
		aliasText := deparseExprAlias(expr)
		if len(aliasText) > 64 {
			return fmt.Sprintf("Name_exp_%d", position)
		}
		return aliasText
	}
}

// deparseExprAlias generates a human-readable expression text for use as an auto-alias.
// Unlike deparseExpr, this preserves the original expression text style matching MySQL 8.0's
// auto-alias behavior: function names stay uppercase, spaces after commas, CASE/CAST
// keywords uppercase, no charset addition in CAST, COUNT(*) keeps *.
func deparseExprAlias(node ast.ExprNode) string {
	switch n := node.(type) {
	case *ast.ColumnRef:
		if n.Table != "" {
			return n.Table + "." + n.Column
		}
		return n.Column
	case *ast.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *ast.FloatLit:
		return n.Value
	case *ast.StringLit:
		// When embedded in an expression, include quotes to match MySQL 8.0's behavior.
		// The top-level autoAlias handles standalone StringLit without quotes.
		// Escape backslashes and single quotes like deparseStringLit does.
		escaped := strings.ReplaceAll(n.Value, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `'`, `\'`)
		return "'" + escaped + "'"
	case *ast.NullLit:
		return "NULL"
	case *ast.BoolLit:
		if n.Value {
			return "TRUE"
		}
		return "FALSE"
	case *ast.HexLit:
		val := n.Value
		if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
			return val
		}
		return "X'" + val + "'"
	case *ast.BitLit:
		val := n.Value
		if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
			return val
		}
		return "b'" + val + "'"
	case *ast.TemporalLit:
		return n.Type + " '" + n.Value + "'"
	case *ast.BinaryExpr:
		left := deparseExprAlias(n.Left)
		right := deparseExprAlias(n.Right)
		op := binaryOpToStringAlias(n.Op)
		// Use original operator text for alias when available (e.g., "MOD" instead of "%", "!=" instead of "<>")
		if n.OriginalOp != "" {
			op = n.OriginalOp
		}
		return left + " " + op + " " + right
	case *ast.UnaryExpr:
		operand := deparseExprAlias(n.Operand)
		switch n.Op {
		case ast.UnaryMinus:
			return "-" + operand
		case ast.UnaryPlus:
			return operand
		case ast.UnaryNot:
			return "NOT " + operand
		case ast.UnaryBitNot:
			return "~" + operand
		}
		return operand
	case *ast.FuncCallExpr:
		// MySQL 8.0 auto-alias preserves original function name case (uppercase from parser).
		name := n.Name
		if n.Star {
			// COUNT(*) → alias "COUNT(*)" — keep *, not 0.
			result := name + "(*)"
			if n.Over != nil {
				result += " " + deparseWindowDefAlias(n.Over)
			}
			return result
		}
		args := make([]string, len(n.Args))
		for i, arg := range n.Args {
			args[i] = deparseExprAlias(arg)
		}
		var result string
		if n.Distinct {
			result = name + "(distinct " + strings.Join(args, ", ") + ")"
		} else {
			result = name + "(" + strings.Join(args, ", ") + ")"
		}
		if n.Over != nil {
			result += " " + deparseWindowDefAlias(n.Over)
		}
		return result
	case *ast.ParenExpr:
		return deparseExprAlias(n.Expr)
	case *ast.CastExpr:
		// MySQL 8.0 auto-alias: "CAST(a AS CHAR)" — uppercase keywords, no charset.
		return "CAST(" + deparseExprAlias(n.Expr) + " AS " + deparseDataTypeAlias(n.TypeName) + ")"
	case *ast.ConvertExpr:
		if n.Charset != "" {
			return "CONVERT(" + deparseExprAlias(n.Expr) + " USING " + strings.ToLower(n.Charset) + ")"
		}
		return "CAST(" + deparseExprAlias(n.Expr) + " AS " + deparseDataTypeAlias(n.TypeName) + ")"
	case *ast.CaseExpr:
		// MySQL 8.0 auto-alias: "CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END" — uppercase keywords.
		var b strings.Builder
		b.WriteString("CASE")
		if n.Operand != nil {
			b.WriteString(" ")
			b.WriteString(deparseExprAlias(n.Operand))
		}
		for _, w := range n.Whens {
			b.WriteString(" WHEN ")
			b.WriteString(deparseExprAlias(w.Cond))
			b.WriteString(" THEN ")
			b.WriteString(deparseExprAlias(w.Result))
		}
		if n.Default != nil {
			b.WriteString(" ELSE ")
			b.WriteString(deparseExprAlias(n.Default))
		}
		b.WriteString(" END")
		return b.String()
	case *ast.IsExpr:
		expr := deparseExprAlias(n.Expr)
		switch n.Test {
		case ast.IsNull:
			if n.Not {
				return expr + " IS NOT NULL"
			}
			return expr + " IS NULL"
		case ast.IsTrue:
			if n.Not {
				return expr + " IS NOT TRUE"
			}
			return expr + " IS TRUE"
		case ast.IsFalse:
			if n.Not {
				return expr + " IS NOT FALSE"
			}
			return expr + " IS FALSE"
		case ast.IsUnknown:
			if n.Not {
				return expr + " IS NOT UNKNOWN"
			}
			return expr + " IS UNKNOWN"
		}
		return deparseExpr(node)
	case *ast.InExpr:
		expr := deparseExprAlias(n.Expr)
		keyword := "IN"
		if n.Not {
			keyword = "NOT IN"
		}
		if n.Select != nil {
			return expr + " " + keyword + " (...)"
		}
		items := make([]string, len(n.List))
		for i, item := range n.List {
			items[i] = deparseExprAlias(item)
		}
		return expr + " " + keyword + " (" + strings.Join(items, ", ") + ")"
	case *ast.BetweenExpr:
		expr := deparseExprAlias(n.Expr)
		low := deparseExprAlias(n.Low)
		high := deparseExprAlias(n.High)
		keyword := "BETWEEN"
		if n.Not {
			keyword = "NOT BETWEEN"
		}
		return expr + " " + keyword + " " + low + " AND " + high
	case *ast.LikeExpr:
		expr := deparseExprAlias(n.Expr)
		pattern := deparseExprAlias(n.Pattern)
		keyword := "LIKE"
		if n.Not {
			keyword = "NOT LIKE"
		}
		result := expr + " " + keyword + " " + pattern
		if n.Escape != nil {
			result += " ESCAPE " + deparseExprAlias(n.Escape)
		}
		return result
	default:
		// Fallback: use the regular deparsed text
		return deparseExpr(node)
	}
}

// deparseDataTypeAlias formats a data type for auto-alias purposes.
// Unlike deparseDataType, this does NOT add charset (matching MySQL 8.0's alias behavior).
func deparseDataTypeAlias(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	name := strings.ToUpper(dt.Name)
	switch strings.ToLower(dt.Name) {
	case "char":
		if dt.Length > 0 {
			return fmt.Sprintf("%s(%d)", name, dt.Length)
		}
		return name
	case "binary":
		if dt.Length > 0 {
			return fmt.Sprintf("%s(%d)", name, dt.Length)
		}
		return name
	case "signed", "signed integer":
		return "SIGNED"
	case "unsigned", "unsigned integer":
		return "UNSIGNED"
	case "decimal":
		if dt.Scale > 0 {
			return fmt.Sprintf("DECIMAL(%d,%d)", dt.Length, dt.Scale)
		}
		if dt.Length > 0 {
			return fmt.Sprintf("DECIMAL(%d)", dt.Length)
		}
		return "DECIMAL"
	default:
		return name
	}
}

// deparseImplicitCrossJoin normalizes multiple FROM tables (implicit cross join)
// into explicit join syntax with parentheses.
// e.g., FROM t1, t2, t3 → ((`t1` join `t2`) join `t3`)
func deparseImplicitCrossJoin(tables []ast.TableExpr) string {
	if len(tables) == 0 {
		return ""
	}
	result := deparseTableExpr(tables[0])
	for i := 1; i < len(tables); i++ {
		result = "(" + result + " join " + deparseTableExpr(tables[i]) + ")"
	}
	return result
}

// deparseTableExpr formats a table expression in the FROM clause.
func deparseTableExpr(tbl ast.TableExpr) string {
	switch t := tbl.(type) {
	case *ast.TableRef:
		return deparseTableRef(t)
	case *ast.JoinClause:
		return deparseJoinClause(t)
	case *ast.SubqueryExpr:
		return deparseSubqueryTableExpr(t)
	default:
		return fmt.Sprintf("/* unsupported table expr: %T */", tbl)
	}
}

// deparseJoinClause formats a JOIN clause.
// MySQL 8.0 format: (`t1` join `t2` on((...)))
func deparseJoinClause(j *ast.JoinClause) string {
	left := deparseTableExpr(j.Left)
	right := deparseTableExpr(j.Right)

	var joinType string
	switch j.Type {
	case ast.JoinInner:
		joinType = "join"
	case ast.JoinLeft:
		joinType = "left join"
	case ast.JoinRight:
		// RIGHT JOIN → LEFT JOIN with table swap
		joinType = "left join"
		left, right = right, left
	case ast.JoinCross:
		// CROSS JOIN → plain join
		joinType = "join"
	case ast.JoinStraight:
		joinType = "straight_join"
	case ast.JoinNatural:
		joinType = "join"
	case ast.JoinNaturalLeft:
		joinType = "left join"
	case ast.JoinNaturalRight:
		joinType = "left join"
		left, right = right, left
	default:
		joinType = "join"
	}

	var b strings.Builder
	b.WriteString("(")
	b.WriteString(left)
	b.WriteString(" ")
	b.WriteString(joinType)
	b.WriteString(" ")
	b.WriteString(right)

	// ON condition
	if j.Condition != nil {
		switch cond := j.Condition.(type) {
		case *ast.OnCondition:
			b.WriteString(" on(")
			b.WriteString(deparseExpr(cond.Expr))
			b.WriteString(")")
		case *ast.UsingCondition:
			// USING (col1, col2) → on((`left`.`col1` = `right`.`col1`) and (...))
			// Requires resolving table names from Left/Right table expressions.
			// For RIGHT JOIN, left/right are already swapped above.
			leftName := tableExprName(j.Left)
			rightName := tableExprName(j.Right)
			if j.Type == ast.JoinRight {
				// Tables were swapped above, so swap names to match original SQL
				leftName, rightName = rightName, leftName
			}
			b.WriteString(" on(")
			b.WriteString(deparseUsingAsOn(cond.Columns, leftName, rightName))
			b.WriteString(")")
		}
	}

	b.WriteString(")")
	return b.String()
}

// tableExprName extracts the effective name (alias or table name) from a table expression.
// Used for USING → ON expansion to qualify column references.
func tableExprName(tbl ast.TableExpr) string {
	switch t := tbl.(type) {
	case *ast.TableRef:
		if t.Alias != "" {
			return t.Alias
		}
		return t.Name
	default:
		return ""
	}
}

// deparseUsingAsOn expands USING columns into ON condition format.
// e.g., USING (a, b) with left=t1, right=t2 → (`t1`.`a` = `t2`.`a`) and (`t1`.`b` = `t2`.`b`)
// MySQL 8.0 format: on((`t1`.`a` = `t2`.`a`))
func deparseUsingAsOn(columns []string, leftName, rightName string) string {
	if len(columns) == 0 {
		return ""
	}
	parts := make([]string, len(columns))
	for i, col := range columns {
		parts[i] = "(`" + leftName + "`.`" + col + "` = `" + rightName + "`.`" + col + "`)"
	}
	if len(parts) == 1 {
		return parts[0]
	}
	// Multiple columns: chain with "and"
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result = "(" + result + " and " + parts[i] + ")"
	}
	return result
}

// deparseSubqueryTableExpr formats a derived table (subquery as table expression).
// MySQL 8.0 format: (select ...) `alias` — no AS keyword for alias
func deparseSubqueryTableExpr(s *ast.SubqueryExpr) string {
	var b strings.Builder
	b.WriteString("(")
	if s.Select != nil {
		b.WriteString(deparseSelectStmt(s.Select))
	}
	b.WriteString(")")
	if s.Alias != "" {
		b.WriteString(" `")
		b.WriteString(s.Alias)
		b.WriteString("`")
	}
	return b.String()
}

// deparseTableRef formats a simple table reference.
func deparseTableRef(t *ast.TableRef) string {
	var b strings.Builder
	if t.Schema != "" {
		b.WriteString("`")
		b.WriteString(t.Schema)
		b.WriteString("`.")
	}
	b.WriteString("`")
	b.WriteString(t.Name)
	b.WriteString("`")
	if t.Alias != "" {
		b.WriteString(" `")
		b.WriteString(t.Alias)
		b.WriteString("`")
	}
	return b.String()
}

func deparseExpr(node ast.ExprNode) string {
	switch n := node.(type) {
	case *ast.IntLit:
		return fmt.Sprintf("%d", n.Value)
	case *ast.FloatLit:
		return n.Value
	case *ast.BoolLit:
		if n.Value {
			return "true"
		}
		return "false"
	case *ast.StringLit:
		return deparseStringLit(n)
	case *ast.NullLit:
		return "NULL"
	case *ast.HexLit:
		return deparseHexLit(n)
	case *ast.BitLit:
		return deparseBitLit(n)
	case *ast.TemporalLit:
		return n.Type + "'" + n.Value + "'"
	case *ast.BinaryExpr:
		return deparseBinaryExpr(n)
	case *ast.ColumnRef:
		return deparseColumnRef(n)
	case *ast.UnaryExpr:
		return deparseUnaryExpr(n)
	case *ast.ParenExpr:
		return deparseExpr(n.Expr)
	case *ast.InExpr:
		return deparseInExpr(n)
	case *ast.BetweenExpr:
		return deparseBetweenExpr(n)
	case *ast.LikeExpr:
		return deparseLikeExpr(n)
	case *ast.IsExpr:
		return deparseIsExpr(n)
	case *ast.RowExpr:
		return deparseRowExpr(n)
	case *ast.CaseExpr:
		return deparseCaseExpr(n)
	case *ast.CastExpr:
		return deparseCastExpr(n)
	case *ast.ConvertExpr:
		return deparseConvertExpr(n)
	case *ast.IntervalExpr:
		return deparseIntervalExpr(n)
	case *ast.CollateExpr:
		return deparseCollateExpr(n)
	case *ast.FuncCallExpr:
		return deparseFuncCallExpr(n)
	case *ast.ExistsExpr:
		return deparseExistsExpr(n)
	case *ast.SubqueryExpr:
		return deparseSubqueryExpr(n)
	default:
		return fmt.Sprintf("/* unsupported: %T */", node)
	}
}

func deparseStringLit(n *ast.StringLit) string {
	// MySQL 8.0 uses backslash escaping for single quotes: '' → \'
	// and preserves backslashes as-is.
	escaped := strings.ReplaceAll(n.Value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
	if n.Charset != "" {
		return n.Charset + "'" + escaped + "'"
	}
	return "'" + escaped + "'"
}

func deparseHexLit(n *ast.HexLit) string {
	// MySQL 8.0 normalizes all hex literals to 0x lowercase form.
	// HexLit.Value is either "0xFF" (0x prefix form) or "FF" (X'' form).
	val := n.Value
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		// Already has 0x prefix — just lowercase
		return "0x" + strings.ToLower(val[2:])
	}
	// X'FF' form — value is just the hex digits
	return "0x" + strings.ToLower(val)
}

func deparseBitLit(n *ast.BitLit) string {
	// MySQL 8.0 converts all bit literals to hex form.
	// BitLit.Value is either "0b1010" (0b prefix form) or "1010" (b'' form).
	val := n.Value
	if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
		val = val[2:]
	}
	// Parse binary string to integer, then format as hex
	i := new(big.Int)
	i.SetString(val, 2)
	return "0x" + fmt.Sprintf("%02x", i)
}

func deparseBinaryExpr(n *ast.BinaryExpr) string {
	// Operator-to-function rewrites:
	// REGEXP → regexp_like(left, right)
	// -> → json_extract(left, right)
	// ->> → json_unquote(json_extract(left, right))
	switch n.Op {
	case ast.BinOpRegexp:
		return "regexp_like(" + deparseExpr(n.Left) + "," + deparseExpr(n.Right) + ")"
	case ast.BinOpJsonExtract:
		return "json_extract(" + deparseExpr(n.Left) + "," + deparseExpr(n.Right) + ")"
	case ast.BinOpJsonUnquote:
		return "json_unquote(json_extract(" + deparseExpr(n.Left) + "," + deparseExpr(n.Right) + "))"
	}

	left := n.Left
	right := n.Right
	// MySQL normalizes INTERVAL + expr to expr + INTERVAL (interval on the right)
	if _, ok := left.(*ast.IntervalExpr); ok {
		if _, ok2 := right.(*ast.IntervalExpr); !ok2 {
			left, right = right, left
		}
	}
	leftStr := deparseExpr(left)
	rightStr := deparseExpr(right)
	op := binaryOpToString(n.Op)
	return "(" + leftStr + " " + op + " " + rightStr + ")"
}

func deparseColumnRef(n *ast.ColumnRef) string {
	if n.Schema != "" {
		return "`" + n.Schema + "`.`" + n.Table + "`.`" + n.Column + "`"
	}
	if n.Table != "" {
		return "`" + n.Table + "`.`" + n.Column + "`"
	}
	return "`" + n.Column + "`"
}

func binaryOpToString(op ast.BinaryOp) string {
	switch op {
	case ast.BinOpAdd:
		return "+"
	case ast.BinOpSub:
		return "-"
	case ast.BinOpMul:
		return "*"
	case ast.BinOpDiv:
		return "/"
	case ast.BinOpMod:
		return "%"
	case ast.BinOpDivInt:
		return "DIV"
	case ast.BinOpEq:
		return "="
	case ast.BinOpNe:
		return "<>"
	case ast.BinOpLt:
		return "<"
	case ast.BinOpGt:
		return ">"
	case ast.BinOpLe:
		return "<="
	case ast.BinOpGe:
		return ">="
	case ast.BinOpNullSafeEq:
		return "<=>"
	case ast.BinOpAnd:
		return "and"
	case ast.BinOpOr:
		return "or"
	case ast.BinOpXor:
		return "xor"
	case ast.BinOpBitAnd:
		return "&"
	case ast.BinOpBitOr:
		return "|"
	case ast.BinOpBitXor:
		return "^"
	case ast.BinOpShiftLeft:
		return "<<"
	case ast.BinOpShiftRight:
		return ">>"
	case ast.BinOpSoundsLike:
		return "sounds like"
	default:
		return "?"
	}
}

// binaryOpToStringAlias returns the operator string for auto-alias purposes.
// MySQL 8.0 uses uppercase AND/OR/XOR in auto-aliases.
func binaryOpToStringAlias(op ast.BinaryOp) string {
	switch op {
	case ast.BinOpAnd:
		return "AND"
	case ast.BinOpOr:
		return "OR"
	case ast.BinOpXor:
		return "XOR"
	default:
		return binaryOpToString(op)
	}
}

// deparseWindowDefAlias formats a window definition for auto-alias purposes.
// MySQL 8.0 uses: OVER (ORDER BY col) in the alias text.
func deparseWindowDefAlias(wd *ast.WindowDef) string {
	var b strings.Builder
	b.WriteString("OVER (")

	needSpace := false
	if len(wd.PartitionBy) > 0 {
		b.WriteString("PARTITION BY ")
		for i, expr := range wd.PartitionBy {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(deparseExprAlias(expr))
		}
		needSpace = true
	}
	if len(wd.OrderBy) > 0 {
		if needSpace {
			b.WriteString(" ")
		}
		b.WriteString("ORDER BY ")
		for i, item := range wd.OrderBy {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(deparseExprAlias(item.Expr))
			if item.Desc {
				b.WriteString(" desc")
			}
		}
	}
	b.WriteString(")")
	return b.String()
}

func deparseUnaryExpr(n *ast.UnaryExpr) string {
	operand := deparseExpr(n.Operand)
	switch n.Op {
	case ast.UnaryMinus:
		return "-" + operand
	case ast.UnaryPlus:
		// MySQL drops unary plus entirely
		return operand
	case ast.UnaryNot:
		return "(not(" + operand + "))"
	case ast.UnaryBitNot:
		return "~(" + operand + ")"
	default:
		return operand
	}
}

func deparseInExpr(n *ast.InExpr) string {
	expr := deparseExpr(n.Expr)
	keyword := "in"
	if n.Not {
		keyword = "not in"
	}

	// IN subquery: a IN (SELECT ...)
	// MySQL 8.0 does NOT wrap IN subquery in outer parens (unlike IN value list).
	// MySQL 8.0 also omits AS aliases in the subquery's target list.
	if n.Select != nil {
		return expr + " " + keyword + " (" + deparseSelectStmtNoAlias(n.Select) + ")"
	}

	// Build the value list with no spaces after commas
	items := make([]string, len(n.List))
	for i, item := range n.List {
		items[i] = deparseExpr(item)
	}
	return "(" + expr + " " + keyword + " (" + strings.Join(items, ",") + "))"
}

func deparseBetweenExpr(n *ast.BetweenExpr) string {
	expr := deparseExpr(n.Expr)
	low := deparseExpr(n.Low)
	high := deparseExpr(n.High)
	keyword := "between"
	if n.Not {
		keyword = "not between"
	}
	return "(" + expr + " " + keyword + " " + low + " and " + high + ")"
}

func deparseLikeExpr(n *ast.LikeExpr) string {
	expr := deparseExpr(n.Expr)
	pattern := deparseExpr(n.Pattern)
	likeClause := "(" + expr + " like " + pattern
	if n.Escape != nil {
		likeClause += " escape " + deparseExpr(n.Escape)
	}
	likeClause += ")"
	if n.Not {
		return "(not(" + likeClause + "))"
	}
	return likeClause
}

func deparseIsExpr(n *ast.IsExpr) string {
	expr := deparseExpr(n.Expr)
	var test string
	switch n.Test {
	case ast.IsNull:
		if n.Not {
			test = "is not null"
		} else {
			test = "is null"
		}
	case ast.IsTrue:
		if n.Not {
			test = "is not true"
		} else {
			test = "is true"
		}
	case ast.IsFalse:
		if n.Not {
			test = "is not false"
		} else {
			test = "is false"
		}
	case ast.IsUnknown:
		if n.Not {
			test = "is not unknown"
		} else {
			test = "is unknown"
		}
	default:
		test = "is ?"
	}
	return "(" + expr + " " + test + ")"
}

func deparseRowExpr(n *ast.RowExpr) string {
	items := make([]string, len(n.Items))
	for i, item := range n.Items {
		items[i] = deparseExpr(item)
	}
	return "row(" + strings.Join(items, ",") + ")"
}

func deparseCaseExpr(n *ast.CaseExpr) string {
	var b strings.Builder
	b.WriteString("(case")
	if n.Operand != nil {
		b.WriteString(" ")
		b.WriteString(deparseExpr(n.Operand))
	}
	for _, w := range n.Whens {
		b.WriteString(" when ")
		b.WriteString(deparseExpr(w.Cond))
		b.WriteString(" then ")
		b.WriteString(deparseExpr(w.Result))
	}
	if n.Default != nil {
		b.WriteString(" else ")
		b.WriteString(deparseExpr(n.Default))
	}
	b.WriteString(" end)")
	return b.String()
}

func deparseCastExpr(n *ast.CastExpr) string {
	expr := deparseExpr(n.Expr)
	typeName := deparseDataType(n.TypeName)
	return "cast(" + expr + " as " + typeName + ")"
}

func deparseConvertExpr(n *ast.ConvertExpr) string {
	expr := deparseExpr(n.Expr)
	// CONVERT(expr USING charset) form
	if n.Charset != "" {
		return "convert(" + expr + " using " + strings.ToLower(n.Charset) + ")"
	}
	// CONVERT(expr, type) form — MySQL rewrites to CAST
	typeName := deparseDataType(n.TypeName)
	return "cast(" + expr + " as " + typeName + ")"
}

func deparseDataType(dt *ast.DataType) string {
	if dt == nil {
		return ""
	}
	name := strings.ToLower(dt.Name)
	switch name {
	case "char":
		result := "char"
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
		// MySQL adds charset for CHAR in CAST
		charset := dt.Charset
		if charset == "" {
			charset = "utf8mb4"
		}
		result += " charset " + strings.ToLower(charset)
		return result
	case "binary":
		// CAST to BINARY becomes cast(x as char charset binary)
		result := "char"
		if dt.Length > 0 {
			result += fmt.Sprintf("(%d)", dt.Length)
		}
		result += " charset binary"
		return result
	case "signed", "signed integer":
		return "signed"
	case "unsigned", "unsigned integer":
		return "unsigned"
	case "decimal":
		if dt.Scale > 0 {
			return fmt.Sprintf("decimal(%d,%d)", dt.Length, dt.Scale)
		}
		if dt.Length > 0 {
			return fmt.Sprintf("decimal(%d)", dt.Length)
		}
		return "decimal"
	case "date":
		return "date"
	case "datetime":
		if dt.Length > 0 {
			return fmt.Sprintf("datetime(%d)", dt.Length)
		}
		return "datetime"
	case "time":
		if dt.Length > 0 {
			return fmt.Sprintf("time(%d)", dt.Length)
		}
		return "time"
	case "json":
		return "json"
	case "float":
		return "float"
	case "double":
		return "double"
	default:
		return name
	}
}

func deparseIntervalExpr(n *ast.IntervalExpr) string {
	val := deparseExpr(n.Value)
	return "interval " + val + " " + strings.ToLower(n.Unit)
}

func deparseCollateExpr(n *ast.CollateExpr) string {
	expr := deparseExpr(n.Expr)
	return "(" + expr + " collate " + n.Collation + ")"
}

// funcNameRewrites maps uppercase function names to their MySQL 8.0 canonical forms.
// These rewrites are applied by SHOW CREATE VIEW in MySQL 8.0.
var funcNameRewrites = map[string]string{
	"SUBSTRING":         "substr",
	"CURRENT_TIMESTAMP":  "now",
	"CURRENT_DATE":       "curdate",
	"CURRENT_TIME":       "curtime",
	"CURRENT_USER":       "current_user",
	"NOW":                "now",
	"LOCALTIME":          "now",
	"LOCALTIMESTAMP":     "now",
}

// deparseTrimDirectional handles TRIM(LEADING|TRAILING|BOTH remstr FROM str).
// MySQL 8.0 SHOW CREATE VIEW format: trim(leading 'x' from `a`)
func deparseTrimDirectional(direction string, args []ast.ExprNode) string {
	if len(args) == 2 {
		remstr := deparseExpr(args[0])
		str := deparseExpr(args[1])
		return "trim(" + direction + " " + remstr + " from " + str + ")"
	}
	// Fallback: single arg (shouldn't happen for directional, but be safe)
	if len(args) == 1 {
		return "trim(" + direction + " " + deparseExpr(args[0]) + ")"
	}
	return "trim()"
}

func deparseFuncCallExpr(n *ast.FuncCallExpr) string {
	// Handle TRIM special forms: TRIM_LEADING, TRIM_TRAILING, TRIM_BOTH
	// Parser encodes these as FuncCallExpr with Name="TRIM_LEADING" etc.
	// Args: [remstr, str] for directional forms
	name := strings.ToUpper(n.Name)
	switch name {
	case "TRIM_LEADING":
		return deparseTrimDirectional("leading", n.Args)
	case "TRIM_TRAILING":
		return deparseTrimDirectional("trailing", n.Args)
	case "TRIM_BOTH":
		return deparseTrimDirectional("both", n.Args)
	}

	// GROUP_CONCAT has special formatting
	if name == "GROUP_CONCAT" {
		return deparseGroupConcat(n)
	}

	// Determine the canonical function name
	canonical, ok := funcNameRewrites[name]
	if !ok {
		canonical = strings.ToLower(n.Name)
	}

	// Schema-qualified name
	if n.Schema != "" {
		canonical = strings.ToLower(n.Schema) + "." + canonical
	}

	// Zero-arg functions (CURRENT_TIMESTAMP, NOW(), etc.) — always emit parens
	if len(n.Args) == 0 && !n.Star {
		result := canonical + "()"
		if n.Over != nil {
			result += " " + deparseWindowDef(n.Over)
		}
		return result
	}

	// COUNT(*) — MySQL 8.0 rewrites COUNT(*) to count(0)
	if n.Star {
		result := canonical + "(0)"
		if n.Over != nil {
			result += " " + deparseWindowDef(n.Over)
		}
		return result
	}

	// Build argument list with no spaces after commas
	args := make([]string, len(n.Args))
	for i, arg := range n.Args {
		args[i] = deparseExpr(arg)
	}

	var argStr string
	if n.Distinct {
		argStr = "distinct " + strings.Join(args, ",")
	} else {
		argStr = strings.Join(args, ",")
	}

	result := canonical + "(" + argStr + ")"

	// Append OVER clause for window functions
	if n.Over != nil {
		result += " " + deparseWindowDef(n.Over)
	}

	return result
}

// deparseWindowDef formats a window definition.
// MySQL 8.0 format: OVER (PARTITION BY ... ORDER BY ... frame_clause )
// Note: trailing space before closing paren, uppercase keywords.
func deparseWindowDef(wd *ast.WindowDef) string {
	// Named window reference: OVER window_name
	if wd.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil {
		return "OVER " + wd.RefName
	}

	var b strings.Builder
	b.WriteString("OVER (")

	needSpace := false

	// PARTITION BY
	if len(wd.PartitionBy) > 0 {
		b.WriteString("PARTITION BY ")
		for i, expr := range wd.PartitionBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(expr))
		}
		needSpace = true
	}

	// ORDER BY
	if len(wd.OrderBy) > 0 {
		if needSpace {
			b.WriteString(" ")
		}
		b.WriteString("ORDER BY ")
		for i, item := range wd.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(item.Expr))
			if item.Desc {
				b.WriteString(" desc")
			}
		}
		needSpace = true
	}

	// Frame clause
	if wd.Frame != nil {
		if needSpace {
			b.WriteString(" ")
		}
		b.WriteString(deparseWindowFrame(wd.Frame))
		needSpace = true
	}

	// Trailing space before closing paren (MySQL 8.0 format)
	b.WriteString(" )")

	return b.String()
}

// deparseWindowFrame formats a window frame specification.
// MySQL 8.0 format: ROWS/RANGE/GROUPS BETWEEN start AND end (all uppercase).
func deparseWindowFrame(f *ast.WindowFrame) string {
	var b strings.Builder

	// Frame type
	switch f.Type {
	case ast.FrameRows:
		b.WriteString("ROWS")
	case ast.FrameRange:
		b.WriteString("RANGE")
	case ast.FrameGroups:
		b.WriteString("GROUPS")
	}

	if f.End != nil {
		// BETWEEN ... AND ... form
		b.WriteString(" BETWEEN ")
		b.WriteString(deparseWindowFrameBound(f.Start))
		b.WriteString(" AND ")
		b.WriteString(deparseWindowFrameBound(f.End))
	} else {
		// Single bound form
		b.WriteString(" ")
		b.WriteString(deparseWindowFrameBound(f.Start))
	}

	return b.String()
}

// deparseWindowFrameBound formats a window frame bound.
func deparseWindowFrameBound(fb *ast.WindowFrameBound) string {
	switch fb.Type {
	case ast.BoundUnboundedPreceding:
		return "UNBOUNDED PRECEDING"
	case ast.BoundPreceding:
		return deparseExpr(fb.Offset) + " PRECEDING"
	case ast.BoundCurrentRow:
		return "CURRENT ROW"
	case ast.BoundFollowing:
		return deparseExpr(fb.Offset) + " FOLLOWING"
	case ast.BoundUnboundedFollowing:
		return "UNBOUNDED FOLLOWING"
	default:
		return "/* unknown bound */"
	}
}

// deparseExistsExpr formats an EXISTS expression.
// MySQL 8.0 format: exists(select ...)
func deparseExistsExpr(n *ast.ExistsExpr) string {
	if n.Select != nil {
		return "exists(" + deparseSelectStmt(n.Select) + ")"
	}
	return "exists(/* subquery */)"
}

// deparseSubqueryExpr formats a subquery expression.
// MySQL 8.0 format: (select ...)
func deparseSubqueryExpr(n *ast.SubqueryExpr) string {
	if n.Select != nil {
		return "(" + deparseSelectStmt(n.Select) + ")"
	}
	return "(/* subquery */)"
}

// deparseGroupConcat handles GROUP_CONCAT with its special syntax:
// group_concat([distinct] expr [order by expr ASC|DESC] separator 'str')
// MySQL 8.0 always shows the separator (default ',') and explicit ASC in ORDER BY.
func deparseGroupConcat(n *ast.FuncCallExpr) string {
	var b strings.Builder
	b.WriteString("group_concat(")

	// DISTINCT
	if n.Distinct {
		b.WriteString("distinct ")
	}

	// Arguments
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(deparseExpr(arg))
	}

	// ORDER BY — MySQL 8.0 always shows explicit ASC/DESC
	if len(n.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range n.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(item.Expr))
			if item.Desc {
				b.WriteString(" DESC")
			} else {
				b.WriteString(" ASC")
			}
		}
	}

	// SEPARATOR — always shown; default is ','
	b.WriteString(" separator ")
	if n.Separator != nil {
		b.WriteString(deparseExpr(n.Separator))
	} else {
		b.WriteString("','")
	}

	b.WriteString(")")
	return b.String()
}

