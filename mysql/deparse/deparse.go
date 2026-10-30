// Package deparse converts MySQL AST nodes back to SQL text,
// matching MySQL 8.0's SHOW CREATE VIEW formatting.
//
// Every node is written into the one strings.Builder of the call: a node's
// text contains its children's, so building each level as a string of its own
// and concatenating it into the parent's copied the text of a long chain
// (1 + 1 + ..., a UNION ALL chain, a join chain, an AND chain) once per level,
// and deparsing took time quadratic in the chain's length.
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
	var b strings.Builder
	writeSelectStmt(&b, stmt, false)
	return b.String()
}

// deparseExpr returns the text of an expression, for the callers that need it
// as a string of its own; everything inside this file writes with writeExpr.
func deparseExpr(node ast.ExprNode) string {
	var b strings.Builder
	writeExpr(&b, node)
	return b.String()
}

// writeSelectStmt writes a query expression. suppressAlias omits the target
// list aliases, for subquery contexts (IN, EXISTS) where MySQL omits AS alias.
func writeSelectStmt(b *strings.Builder, stmt *ast.SelectStmt, suppressAlias bool) {
	// Parenthesized query expression: ( inner ) with optional OUTER ORDER BY /
	// LIMIT held on the wrapper. Materialize the parens — a dropped ParenSource
	// silently loses the inner query.
	if stmt.ParenSource != nil {
		writeParenSource(b, stmt, suppressAlias)
		return
	}

	if stmt.TableSource != nil {
		b.WriteString(deparseTableQueryPrimary(stmt))
		return
	}

	if stmt.ValuesSource != nil {
		b.WriteString(deparseValuesQueryPrimary(stmt))
		return
	}

	// Handle set operations: UNION / UNION ALL / INTERSECT / EXCEPT
	if stmt.SetOp != ast.SetOpNone {
		writeSetOperation(b, stmt, suppressAlias)
		return
	}

	// CTE (WITH clause) — emit before the SELECT keyword
	if len(stmt.CTEs) > 0 {
		writeCTEs(b, stmt.CTEs)
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
			writeResTargetNoAlias(b, target)
		} else {
			writeResTarget(b, target, i+1)
		}
	}

	// FROM clause
	if len(stmt.From) > 0 {
		b.WriteString(" from ")
		if len(stmt.From) == 1 {
			writeTableExpr(b, stmt.From[0])
		} else {
			// Multiple tables (implicit cross join) → normalized to explicit join with parens
			// e.g., FROM t1, t2 → from (`t1` join `t2`)
			// For 3+ tables: FROM t1, t2, t3 → from ((`t1` join `t2`) join `t3`)
			writeImplicitCrossJoin(b, stmt.From)
		}
	}

	// WHERE clause
	if stmt.Where != nil {
		b.WriteString(" where ")
		writeExpr(b, stmt.Where)
	}

	// GROUP BY clause
	if len(stmt.GroupBy) > 0 {
		b.WriteString(" group by ")
		for i, expr := range stmt.GroupBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, expr)
		}
		if stmt.WithRollup {
			b.WriteString(" with rollup")
		}
	}

	// HAVING clause
	if stmt.Having != nil {
		b.WriteString(" having ")
		writeExpr(b, stmt.Having)
	}

	// WINDOW clause: WINDOW `w` AS (window_spec) [, ...]
	// MySQL 8.0 appends a trailing space after the window clause.
	if len(stmt.WindowClause) > 0 {
		b.WriteString(" window ")
		for i, wd := range stmt.WindowClause {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(quoteIdent(wd.Name))
			b.WriteString(" AS ")
			writeWindowBody(b, wd)
		}
		b.WriteString(" ")
	}

	// ORDER BY clause
	if len(stmt.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range stmt.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, item.Expr)
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
			writeExpr(b, stmt.Limit.Offset)
			b.WriteString(",")
		}
		writeExpr(b, stmt.Limit.Count)
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if stmt.ForUpdate != nil {
		b.WriteString(" ")
		b.WriteString(deparseForUpdate(stmt.ForUpdate))
	}
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
			b.WriteString(quoteIdent(tbl.Name))
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

// writeParenSource formats a parenthesized query expression: an optional WITH,
// '(' inner ')', and the OUTER trailing ORDER BY / LIMIT held on the wrapper.
// The inner query keeps its own clauses, so both scopes round-trip (e.g.
// "(select 1 limit 5) limit 2").
func writeParenSource(b *strings.Builder, stmt *ast.SelectStmt, suppressAlias bool) {
	if len(stmt.CTEs) > 0 {
		writeCTEs(b, stmt.CTEs)
		b.WriteString(" ")
	}

	b.WriteString("(")
	// Propagate the caller's alias context: top-level parenthesized query
	// expressions (set-op ARMs, stacked LIMIT) keep select-item aliases like
	// every stored form, while wrappers reached through expression-subquery
	// rendering suppress them like every sibling subquery path.
	writeSelectStmt(b, stmt.ParenSource, suppressAlias)
	b.WriteString(")")

	// Outer ORDER BY (applies to the parenthesized result).
	if len(stmt.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range stmt.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, item.Expr)
			if item.Desc {
				b.WriteString(" desc")
			}
		}
	}

	// Outer LIMIT (applies to the parenthesized result).
	if stmt.Limit != nil {
		b.WriteString(" limit ")
		if stmt.Limit.Offset != nil {
			writeExpr(b, stmt.Limit.Offset)
			b.WriteString(",")
		}
		writeExpr(b, stmt.Limit.Count)
	}
}

func deparseTableQueryPrimary(stmt *ast.SelectStmt) string {
	var b strings.Builder

	if len(stmt.CTEs) > 0 {
		b.WriteString(deparseCTEs(stmt.CTEs))
		b.WriteString(" ")
	}

	b.WriteString("table ")
	if stmt.TableSource.Table != nil {
		b.WriteString(deparseTableRef(stmt.TableSource.Table))
	}
	appendOrderByLimit(&b, stmt.OrderBy, stmt.Limit)
	return b.String()
}

func deparseValuesQueryPrimary(stmt *ast.SelectStmt) string {
	var b strings.Builder

	if len(stmt.CTEs) > 0 {
		b.WriteString(deparseCTEs(stmt.CTEs))
		b.WriteString(" ")
	}

	b.WriteString("values ")
	for i, row := range stmt.ValuesSource.Rows {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("row(")
		for j, expr := range row {
			if j > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(expr))
		}
		b.WriteString(")")
	}
	appendOrderByLimit(&b, stmt.OrderBy, stmt.Limit)
	return b.String()
}

func appendOrderByLimit(b *strings.Builder, orderBy []*ast.OrderByItem, limit *ast.Limit) {
	if len(orderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range orderBy {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(item.Expr))
			if item.Desc {
				b.WriteString(" desc")
			}
		}
	}
	if limit != nil {
		b.WriteString(" limit ")
		if limit.Offset != nil {
			b.WriteString(deparseExpr(limit.Offset))
			b.WriteString(",")
		}
		b.WriteString(deparseExpr(limit.Count))
	}
}

// writeSetOperation formats a set operation (UNION, INTERSECT, EXCEPT).
// MySQL 8.0 format: select ... union [all] select ... (flat, no parens around sub-selects)
// The WITH clause, held on the set-operation root, comes before the whole set
// operation. suppressAlias omits the operands' select-item aliases, as for a
// set operation inside a subquery expression: (select 1 union select 2).
func writeSetOperation(b *strings.Builder, stmt *ast.SelectStmt, suppressAlias bool) {
	if len(stmt.CTEs) > 0 {
		writeCTEs(b, stmt.CTEs)
		b.WriteString(" ")
	}

	writeSelectStmt(b, stmt.Left, suppressAlias)

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

	b.WriteString(" ")
	b.WriteString(op)
	b.WriteString(" ")
	writeSelectStmt(b, stmt.Right, suppressAlias)

	// ORDER BY clause (applies to the entire set operation)
	if len(stmt.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range stmt.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, item.Expr)
			if item.Desc {
				b.WriteString(" desc")
			}
		}
	}

	// LIMIT clause (applies to the entire set operation)
	if stmt.Limit != nil {
		b.WriteString(" limit ")
		if stmt.Limit.Offset != nil {
			writeExpr(b, stmt.Limit.Offset)
			b.WriteString(",")
		}
		writeExpr(b, stmt.Limit.Count)
	}
}

func deparseCTEs(ctes []*ast.CommonTableExpr) string {
	var b strings.Builder
	writeCTEs(&b, ctes)
	return b.String()
}

// writeCTEs formats a WITH clause (one or more CTEs).
// MySQL 8.0 format: with [recursive] `name` [(`col`, ...)] as (select ...) [, ...]
func writeCTEs(b *strings.Builder, ctes []*ast.CommonTableExpr) {
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
			b.WriteString(", ")
		}
		b.WriteString(quoteIdent(cte.Name))

		// Column list
		if len(cte.Columns) > 0 {
			b.WriteString(" (")
			for j, col := range cte.Columns {
				if j > 0 {
					b.WriteString(",")
				}
				b.WriteString(quoteIdent(col))
			}
			b.WriteString(")")
		}

		b.WriteString(" as (")
		if cte.Select != nil {
			writeSelectStmt(b, cte.Select, false)
		}
		b.WriteString(")")
	}
}

// writeResTargetNoAlias formats a target list entry without the AS alias.
// Used for subquery contexts (IN, EXISTS) where MySQL omits aliases.
func writeResTargetNoAlias(b *strings.Builder, node ast.ExprNode) {
	if rt, ok := node.(*ast.ResTarget); ok {
		writeExpr(b, rt.Val)
		return
	}
	writeExpr(b, node)
}

// writeResTarget formats a single result target in the SELECT list.
// MySQL 8.0 SHOW CREATE VIEW format: expr AS `alias`
// - Always uses AS keyword
// - Alias is always backtick-quoted
// - Auto-alias: column ref → column name; literal → literal text; expression → expression text
func writeResTarget(b *strings.Builder, node ast.ExprNode, position int) {
	rt, isRT := node.(*ast.ResTarget)

	var expr ast.ExprNode
	var explicitAlias string
	if isRT {
		expr = rt.Val
		explicitAlias = rt.Name
	} else {
		expr = node
	}

	exprStart := b.Len()
	writeExpr(b, expr)

	// Determine alias
	alias := explicitAlias
	if alias == "" {
		alias = autoAlias(expr, b.String()[exprStart:], position)
	}

	// MySQL 8.0 uses double-space before AS for window function expressions.
	// The OVER clause already ends with " )", and MySQL adds an extra space.
	if hasWindowFunction(expr) {
		b.WriteString("  AS ")
	} else {
		b.WriteString(" AS ")
	}
	b.WriteString(quoteIdent(alias))
}

// hasWindowFunction checks if an expression is or contains a window function
// with an inline OVER (...) clause (not a named window reference like OVER w).
// MySQL 8.0 uses double-space before AS only for inline window definitions.
func hasWindowFunction(node ast.ExprNode) bool {
	if fc, ok := node.(*ast.FuncCallExpr); ok && fc.Over != nil {
		// Named window reference (OVER w) — no double space
		if fc.Over.RefName != "" && len(fc.Over.PartitionBy) == 0 && len(fc.Over.OrderBy) == 0 && fc.Over.Frame == nil {
			return false
		}
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
		// An adjacent-literal run is named after its FIRST segment, not the
		// folded value; an empty first segment falls back to Name_exp_N
		// (oracle: SELECT '' 'b' → value 'b' named Name_exp_1).
		name := n.Value
		if n.Concatenated {
			name = n.FirstSegment
		}
		if name == "" {
			return fmt.Sprintf("Name_exp_%d", position)
		}
		return name
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
	var b strings.Builder
	writeExprAlias(&b, node)
	return b.String()
}

// writeExprAlias writes the auto-alias text of an expression (see deparseExprAlias).
func writeExprAlias(b *strings.Builder, node ast.ExprNode) {
	switch n := node.(type) {
	case *ast.ColumnRef:
		if n.Table != "" {
			b.WriteString(n.Table)
			b.WriteString(".")
		}
		b.WriteString(n.Column)
	case *ast.IntLit:
		fmt.Fprintf(b, "%d", n.Value)
	case *ast.FloatLit:
		b.WriteString(n.Value)
	case *ast.StringLit:
		// When embedded in an expression, include quotes to match MySQL 8.0's behavior.
		// The top-level autoAlias handles standalone StringLit without quotes.
		// Escape backslashes and single quotes like deparseStringLit does.
		escaped := strings.ReplaceAll(n.Value, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `'`, `\'`)
		b.WriteString("'")
		b.WriteString(escaped)
		b.WriteString("'")
	case *ast.NullLit:
		b.WriteString("NULL")
	case *ast.BoolLit:
		if n.Value {
			b.WriteString("TRUE")
		} else {
			b.WriteString("FALSE")
		}
	case *ast.HexLit:
		val := n.Value
		if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
			b.WriteString(val)
		} else {
			b.WriteString("X'")
			b.WriteString(val)
			b.WriteString("'")
		}
	case *ast.BitLit:
		val := n.Value
		if strings.HasPrefix(val, "0b") || strings.HasPrefix(val, "0B") {
			b.WriteString(val)
		} else {
			b.WriteString("b'")
			b.WriteString(val)
			b.WriteString("'")
		}
	case *ast.TemporalLit:
		b.WriteString(n.Type)
		b.WriteString(" '")
		b.WriteString(n.Value)
		b.WriteString("'")
	case *ast.BinaryExpr:
		// Special operator aliases for REGEXP, ->, ->>
		var op string
		switch n.Op {
		case ast.BinOpRegexp:
			op = " REGEXP "
		case ast.BinOpJsonExtract:
			op = "->"
		case ast.BinOpJsonUnquote:
			op = "->>"
		case ast.BinOpSoundsLike:
			op = " SOUNDS LIKE "
		default:
			op = binaryOpToStringAlias(n.Op)
			// Use original operator text for alias when available (e.g., "MOD" instead of "%", "!=" instead of "<>")
			if n.OriginalOp != "" {
				op = n.OriginalOp
			}
			op = " " + op + " "
		}
		writeExprAlias(b, n.Left)
		b.WriteString(op)
		writeExprAlias(b, n.Right)
	case *ast.UnaryExpr:
		switch n.Op {
		case ast.UnaryMinus:
			b.WriteString("-")
		case ast.UnaryNot:
			// NOT REGEXP → "a NOT REGEXP 'pattern'" (NOT between left and REGEXP)
			if binExpr, ok := unwrapParen(n.Operand).(*ast.BinaryExpr); ok && binExpr.Op == ast.BinOpRegexp {
				writeExprAlias(b, binExpr.Left)
				b.WriteString(" NOT REGEXP ")
				writeExprAlias(b, binExpr.Right)
				return
			}
			// MySQL 8.0 preserves ! vs NOT in auto-alias.
			if n.OriginalOp == "!" {
				b.WriteString("!")
			} else {
				b.WriteString("NOT ")
			}
		case ast.UnaryBitNot:
			b.WriteString("~")
		}
		writeExprAlias(b, n.Operand)
	case *ast.FuncCallExpr:
		// MySQL 8.0 auto-alias preserves original function name case (uppercase from parser).
		name := n.Name
		upperName := strings.ToUpper(name)

		// Handle TRIM directional forms: TRIM_LEADING → TRIM(LEADING 'x' FROM a)
		trimHead := ""
		switch upperName {
		case "TRIM_LEADING":
			trimHead = "TRIM(LEADING "
		case "TRIM_TRAILING":
			trimHead = "TRIM(TRAILING "
		case "TRIM_BOTH":
			trimHead = "TRIM(BOTH "
		case "TRIM_FROM":
			trimHead = "TRIM("
		}
		if trimHead != "" && len(n.Args) == 2 {
			b.WriteString(trimHead)
			writeExprAlias(b, n.Args[0])
			b.WriteString(" FROM ")
			writeExprAlias(b, n.Args[1])
			b.WriteString(")")
			return
		}

		// Handle GROUP_CONCAT: alias includes ORDER BY and SEPARATOR
		if upperName == "GROUP_CONCAT" {
			writeGroupConcatAlias(b, n)
			return
		}

		if n.Star {
			// COUNT(*) → alias "COUNT(*)" — keep *, not 0.
			b.WriteString(name)
			b.WriteString("(*)")
			if n.Over != nil {
				b.WriteString(" ")
				writeWindowDefAlias(b, n.Over)
			}
			return
		}
		// Zero-arg keyword functions without explicit parens: alias is just the keyword name.
		// e.g., CURRENT_TIMESTAMP → alias "CURRENT_TIMESTAMP" (no parens).
		// With parens: CURRENT_TIMESTAMP() → alias "CURRENT_TIMESTAMP()".
		if len(n.Args) == 0 && !n.HasParens {
			b.WriteString(name)
			return
		}
		b.WriteString(name)
		if n.Distinct {
			b.WriteString("(DISTINCT ")
		} else {
			b.WriteString("(")
		}
		for i, arg := range n.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, arg)
		}
		b.WriteString(")")
		if n.Over != nil {
			b.WriteString(" ")
			writeWindowDefAlias(b, n.Over)
		}
	case *ast.ParenExpr:
		b.WriteString("(")
		writeExprAlias(b, n.Expr)
		b.WriteString(")")
	case *ast.CastExpr:
		// MySQL 8.0 auto-alias: "CAST(a AS CHAR)" — uppercase keywords, no charset.
		typeName := deparseDataTypeAlias(n.TypeName)
		if n.Array {
			typeName += " ARRAY"
		}
		b.WriteString("CAST(")
		writeExprAlias(b, n.Expr)
		b.WriteString(" AS ")
		b.WriteString(typeName)
		b.WriteString(")")
	case *ast.ConvertExpr:
		b.WriteString("CONVERT(")
		writeExprAlias(b, n.Expr)
		if n.Charset != "" {
			b.WriteString(" USING ")
			b.WriteString(strings.ToLower(n.Charset))
		} else {
			// MySQL 8.0 auto-alias preserves "CONVERT(a, CHAR)" form (comma-separated).
			b.WriteString(", ")
			b.WriteString(deparseDataTypeAlias(n.TypeName))
		}
		b.WriteString(")")
	case *ast.CaseExpr:
		// MySQL 8.0 auto-alias: "CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END" — uppercase keywords.
		b.WriteString("CASE")
		if n.Operand != nil {
			b.WriteString(" ")
			writeExprAlias(b, n.Operand)
		}
		for _, w := range n.Whens {
			b.WriteString(" WHEN ")
			writeExprAlias(b, w.Cond)
			b.WriteString(" THEN ")
			writeExprAlias(b, w.Result)
		}
		if n.Default != nil {
			b.WriteString(" ELSE ")
			writeExprAlias(b, n.Default)
		}
		b.WriteString(" END")
	case *ast.IsExpr:
		var test string
		switch n.Test {
		case ast.IsNull:
			test = " IS NULL"
			if n.Not {
				test = " IS NOT NULL"
			}
		case ast.IsTrue:
			test = " IS TRUE"
			if n.Not {
				test = " IS NOT TRUE"
			}
		case ast.IsFalse:
			test = " IS FALSE"
			if n.Not {
				test = " IS NOT FALSE"
			}
		case ast.IsUnknown:
			test = " IS UNKNOWN"
			if n.Not {
				test = " IS NOT UNKNOWN"
			}
		default:
			writeExpr(b, node)
			return
		}
		writeExprAlias(b, n.Expr)
		b.WriteString(test)
	case *ast.InExpr:
		writeExprAlias(b, n.Expr)
		if n.Not {
			b.WriteString(" NOT IN (")
		} else {
			b.WriteString(" IN (")
		}
		if n.Select != nil {
			b.WriteString("...)")
			return
		}
		for i, item := range n.List {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, item)
		}
		b.WriteString(")")
	case *ast.BetweenExpr:
		writeExprAlias(b, n.Expr)
		if n.Not {
			b.WriteString(" NOT BETWEEN ")
		} else {
			b.WriteString(" BETWEEN ")
		}
		writeExprAlias(b, n.Low)
		b.WriteString(" AND ")
		writeExprAlias(b, n.High)
	case *ast.LikeExpr:
		writeExprAlias(b, n.Expr)
		if n.Not {
			b.WriteString(" NOT LIKE ")
		} else {
			b.WriteString(" LIKE ")
		}
		writeExprAlias(b, n.Pattern)
		if n.Escape != nil {
			b.WriteString(" ESCAPE ")
			writeExprAlias(b, n.Escape)
		}
	case *ast.ExistsExpr:
		// MySQL 8.0 auto-alias for EXISTS: "EXISTS(SELECT 1 FROM t WHERE a > 0)" — uppercase keywords,
		// unqualified column names, no backtick quoting, no column aliases.
		if n.Select != nil {
			b.WriteString("EXISTS(")
			writeSelectStmtAlias(b, n.Select)
			b.WriteString(")")
			return
		}
		b.WriteString("EXISTS(/* subquery */)")
	case *ast.SubqueryExpr:
		// MySQL 8.0 auto-alias for subquery: "(SELECT MAX(a) FROM t)" — uppercase keywords,
		// unqualified column names, no backtick quoting, no column aliases.
		if n.Select != nil {
			b.WriteString("(")
			writeSelectStmtAlias(b, n.Select)
			b.WriteString(")")
			return
		}
		b.WriteString("(/* subquery */)")
	case *ast.RowExpr:
		// Row constructor: "(a, b)" — alias-style items (no backtick quoting,
		// no inner column aliases). Reachable unaliased since a parenthesized
		// subquery may head a row: ((SELECT 1), 2) = ROW(1,2).
		b.WriteString("(")
		for i, item := range n.Items {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, item)
		}
		b.WriteString(")")
	case *ast.CollateExpr:
		// MySQL 8.0 auto-alias: "a COLLATE utf8mb4_unicode_ci" — uppercase COLLATE, no parens.
		writeExprAlias(b, n.Expr)
		b.WriteString(" COLLATE ")
		b.WriteString(n.Collation)
	case *ast.IntervalExpr:
		// MySQL 8.0 auto-alias: "INTERVAL 1 DAY" — uppercase keywords.
		b.WriteString("INTERVAL ")
		writeExprAlias(b, n.Value)
		b.WriteString(" ")
		b.WriteString(strings.ToUpper(n.Unit))
	case *ast.KeywordArg:
		b.WriteString(n.Keyword)
	case *ast.ExtractExpr:
		// Alias approximation with the uppercase-keyword convention:
		// "EXTRACT(DAY_HOUR FROM a)".
		b.WriteString("EXTRACT(")
		b.WriteString(n.Unit)
		b.WriteString(" FROM ")
		writeExprAlias(b, n.Expr)
		b.WriteString(")")
	case *ast.WeightStringExpr:
		b.WriteString("WEIGHT_STRING(")
		writeExprAlias(b, n.Expr)
		if n.AsChar != nil {
			fmt.Fprintf(b, " AS CHAR(%d)", n.AsChar.Length)
		}
		for i, lv := range n.Levels {
			if i == 0 {
				b.WriteString(" LEVEL ")
			} else {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "%d", lv.Level)
			if lv.Desc {
				b.WriteString(" DESC")
			}
			if lv.Reverse {
				b.WriteString(" REVERSE")
			}
		}
		b.WriteString(")")
	default:
		// Fallback: use the regular deparsed text
		writeExpr(b, node)
	}
}

// writeSelectStmtAlias generates a human-readable query expression for
// auto-alias purposes: uppercase keywords, unqualified column names, and each
// item's own alias. MySQL names an item after its source text, so this
// matches it only for a query written this way; Resolver.Source gives the
// source text itself.
func writeSelectStmtAlias(b *strings.Builder, stmt *ast.SelectStmt) {
	if len(stmt.CTEs) > 0 {
		writeCTEsAlias(b, stmt.CTEs)
	}
	if stmt.ParenSource != nil {
		b.WriteString("(")
		writeSelectStmtAlias(b, stmt.ParenSource)
		b.WriteString(")")
		writeOrderByLimitAlias(b, stmt.OrderBy, stmt.Limit)
		return
	}
	if stmt.SetOp != ast.SetOpNone {
		writeSelectStmtAlias(b, stmt.Left)
		switch stmt.SetOp {
		case ast.SetOpUnion:
			b.WriteString(" UNION ")
		case ast.SetOpIntersect:
			b.WriteString(" INTERSECT ")
		case ast.SetOpExcept:
			b.WriteString(" EXCEPT ")
		}
		if stmt.SetAll {
			b.WriteString("ALL ")
		}
		writeSelectStmtAlias(b, stmt.Right)
		writeOrderByLimitAlias(b, stmt.OrderBy, stmt.Limit)
		return
	}
	if stmt.TableSource != nil {
		b.WriteString("TABLE ")
		if t := stmt.TableSource.Table; t != nil {
			if t.Schema != "" {
				b.WriteString(t.Schema)
				b.WriteString(".")
			}
			b.WriteString(t.Name)
		}
		writeOrderByLimitAlias(b, stmt.OrderBy, stmt.Limit)
		return
	}
	if stmt.ValuesSource != nil {
		b.WriteString("VALUES ")
		for i, row := range stmt.ValuesSource.Rows {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("ROW(")
			for j, expr := range row {
				if j > 0 {
					b.WriteString(", ")
				}
				writeExprAlias(b, expr)
			}
			b.WriteString(")")
		}
		writeOrderByLimitAlias(b, stmt.OrderBy, stmt.Limit)
		return
	}

	b.WriteString("SELECT ")
	if stmt.DistinctKind == ast.DistinctOn {
		b.WriteString("DISTINCT ")
	}

	// Target list: an item's own alias is part of its query's text.
	for i, target := range stmt.TargetList {
		if i > 0 {
			b.WriteString(", ")
		}
		if rt, ok := target.(*ast.ResTarget); ok {
			writeExprAlias(b, rt.Val)
			if rt.Name != "" {
				b.WriteString(" AS ")
				b.WriteString(rt.Name)
			}
		} else {
			writeExprAlias(b, target)
		}
	}

	// FROM clause
	if len(stmt.From) > 0 {
		b.WriteString(" FROM ")
		for i, tbl := range stmt.From {
			if i > 0 {
				b.WriteString(", ")
			}
			writeTableExprAlias(b, tbl)
		}
	}

	// WHERE clause
	if stmt.Where != nil {
		b.WriteString(" WHERE ")
		writeExprAlias(b, stmt.Where)
	}

	if len(stmt.GroupBy) > 0 {
		b.WriteString(" GROUP BY ")
		for i, expr := range stmt.GroupBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, expr)
		}
		if stmt.WithRollup {
			b.WriteString(" WITH ROLLUP")
		}
	}
	if stmt.Having != nil {
		b.WriteString(" HAVING ")
		writeExprAlias(b, stmt.Having)
	}
	writeOrderByLimitAlias(b, stmt.OrderBy, stmt.Limit)
}

// writeCTEsAlias writes a WITH clause for auto-alias purposes, followed by a
// space: WITH c AS (SELECT 1).
func writeCTEsAlias(b *strings.Builder, ctes []*ast.CommonTableExpr) {
	b.WriteString("WITH ")
	for _, cte := range ctes {
		if cte.Recursive {
			b.WriteString("RECURSIVE ")
			break
		}
	}
	for i, cte := range ctes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(cte.Name)
		if len(cte.Columns) > 0 {
			b.WriteString("(")
			b.WriteString(strings.Join(cte.Columns, ", "))
			b.WriteString(")")
		}
		b.WriteString(" AS (")
		if cte.Select != nil {
			writeSelectStmtAlias(b, cte.Select)
		}
		b.WriteString(")")
	}
	b.WriteString(" ")
}

// writeOrderByLimitAlias writes ORDER BY and LIMIT clauses for auto-alias
// purposes.
func writeOrderByLimitAlias(b *strings.Builder, orderBy []*ast.OrderByItem, limit *ast.Limit) {
	if len(orderBy) > 0 {
		b.WriteString(" ORDER BY ")
		for i, item := range orderBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, item.Expr)
			if item.Desc {
				b.WriteString(" DESC")
			}
		}
	}
	if limit != nil {
		b.WriteString(" LIMIT ")
		if limit.Offset != nil {
			writeExprAlias(b, limit.Offset)
			b.WriteString(", ")
		}
		writeExprAlias(b, limit.Count)
	}
}

// writeTableExprAlias formats a table expression for auto-alias purposes.
func writeTableExprAlias(b *strings.Builder, tbl ast.TableExpr) {
	switch t := tbl.(type) {
	case *ast.TableRef:
		b.WriteString(t.Name)
		if t.Alias != "" {
			b.WriteString(" ")
			b.WriteString(t.Alias)
		}
	default:
		writeExprAlias(b, tbl.(ast.ExprNode))
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

// writeImplicitCrossJoin normalizes multiple FROM tables (implicit cross join)
// into explicit join syntax with parentheses.
// e.g., FROM t1, t2, t3 → ((`t1` join `t2`) join `t3`)
func writeImplicitCrossJoin(b *strings.Builder, tables []ast.TableExpr) {
	if len(tables) == 0 {
		return
	}
	for i := 1; i < len(tables); i++ {
		b.WriteString("(")
	}
	writeTableExpr(b, tables[0])
	for i := 1; i < len(tables); i++ {
		b.WriteString(" join ")
		writeTableExpr(b, tables[i])
		b.WriteString(")")
	}
}

// writeTableExpr formats a table expression in the FROM clause.
func writeTableExpr(b *strings.Builder, tbl ast.TableExpr) {
	switch t := tbl.(type) {
	case *ast.TableRef:
		b.WriteString(deparseTableRef(t))
	case *ast.JoinClause:
		writeJoinClause(b, t)
	case *ast.SubqueryExpr:
		writeSubqueryTableExpr(b, t)
	default:
		fmt.Fprintf(b, "/* unsupported table expr: %T */", tbl)
	}
}

// writeJoinClause formats a JOIN clause.
// MySQL 8.0 format: (`t1` join `t2` on((...)))
func writeJoinClause(b *strings.Builder, j *ast.JoinClause) {
	left, right := j.Left, j.Right

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

	b.WriteString("(")
	writeTableExpr(b, left)
	b.WriteString(" ")
	b.WriteString(joinType)
	b.WriteString(" ")
	writeTableExpr(b, right)

	// ON condition
	if j.Condition != nil {
		switch cond := j.Condition.(type) {
		case *ast.OnCondition:
			b.WriteString(" on(")
			writeExpr(b, cond.Expr)
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
			writeUsingAsOn(b, cond.Columns, leftName, rightName)
			b.WriteString(")")
		}
	}

	b.WriteString(")")
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
	case *ast.SubqueryExpr:
		return t.Alias
	default:
		return ""
	}
}

// writeUsingAsOn expands USING columns into ON condition format.
// e.g., USING (a, b) with left=t1, right=t2 → ((`t1`.`a` = `t2`.`a`) and (`t1`.`b` = `t2`.`b`))
// MySQL 8.0 format: on((`t1`.`a` = `t2`.`a`))
func writeUsingAsOn(b *strings.Builder, columns []string, leftName, rightName string) {
	// Multiple columns: chain with "and", left-deep.
	for i := 1; i < len(columns); i++ {
		b.WriteString("(")
	}
	for i, col := range columns {
		if i > 0 {
			b.WriteString(" and ")
		}
		b.WriteString("(" + quoteIdent(leftName) + "." + quoteIdent(col) + " = " + quoteIdent(rightName) + "." + quoteIdent(col) + ")")
		if i > 0 {
			b.WriteString(")")
		}
	}
}

// writeSubqueryTableExpr formats a derived table (subquery as table expression).
// MySQL 8.0 format: (select ...) `alias` — no AS keyword for alias
func writeSubqueryTableExpr(b *strings.Builder, s *ast.SubqueryExpr) {
	b.WriteString("(")
	if s.Select != nil {
		writeSelectStmt(b, s.Select, false)
	}
	b.WriteString(")")
	if s.Alias != "" {
		b.WriteString(" " + quoteIdent(s.Alias))
	}
}

// deparseTableRef formats a simple table reference.
func deparseTableRef(t *ast.TableRef) string {
	var b strings.Builder
	if t.Schema != "" {
		b.WriteString(quoteIdent(t.Schema))
		b.WriteString(".")
	}
	b.WriteString(quoteIdent(t.Name))
	if t.Alias != "" {
		b.WriteString(" " + quoteIdent(t.Alias))
	}
	return b.String()
}

func writeExpr(b *strings.Builder, node ast.ExprNode) {
	switch n := node.(type) {
	case *ast.IntLit:
		fmt.Fprintf(b, "%d", n.Value)
	case *ast.FloatLit:
		b.WriteString(n.Value)
	case *ast.BoolLit:
		if n.Value {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case *ast.StringLit:
		b.WriteString(deparseStringLit(n))
	case *ast.NullLit:
		b.WriteString("NULL")
	case *ast.HexLit:
		b.WriteString(deparseHexLit(n))
	case *ast.BitLit:
		b.WriteString(deparseBitLit(n))
	case *ast.TemporalLit:
		b.WriteString(n.Type)
		b.WriteString("'")
		b.WriteString(n.Value)
		b.WriteString("'")
	case *ast.BinaryExpr:
		writeBinaryExpr(b, n)
	case *ast.ColumnRef:
		b.WriteString(deparseColumnRef(n))
	case *ast.UnaryExpr:
		writeUnaryExpr(b, n)
	case *ast.ParenExpr:
		writeExpr(b, n.Expr)
	case *ast.InExpr:
		writeInExpr(b, n)
	case *ast.BetweenExpr:
		writeBetweenExpr(b, n)
	case *ast.LikeExpr:
		writeLikeExpr(b, n)
	case *ast.IsExpr:
		writeIsExpr(b, n)
	case *ast.RowExpr:
		writeRowExpr(b, n)
	case *ast.CaseExpr:
		writeCaseExpr(b, n)
	case *ast.CastExpr:
		writeCastExpr(b, n)
	case *ast.ConvertExpr:
		writeConvertExpr(b, n)
	case *ast.IntervalExpr:
		writeIntervalExpr(b, n)
	case *ast.KeywordArg:
		// Bare keyword argument (timestampdiff unit, get_format type) —
		// never quoted: the engine stores timestampdiff(SECOND,...) and
		// rejects a backtick-quoted unit (oracle 8.0.32 + 5.7.25).
		b.WriteString(n.Keyword)
	case *ast.ExtractExpr:
		writeExtractExpr(b, n)
	case *ast.WeightStringExpr:
		writeWeightStringExpr(b, n)
	case *ast.CollateExpr:
		writeCollateExpr(b, n)
	case *ast.FuncCallExpr:
		writeFuncCallExpr(b, n)
	case *ast.ExistsExpr:
		writeExistsExpr(b, n)
	case *ast.SubqueryExpr:
		writeSubqueryExpr(b, n)
	default:
		fmt.Fprintf(b, "/* unsupported: %T */", node)
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

func writeBinaryExpr(b *strings.Builder, n *ast.BinaryExpr) {
	// Operator-to-function rewrites:
	// REGEXP → regexp_like(left, right)
	// -> → json_extract(left, right)
	// ->> → json_unquote(json_extract(left, right))
	switch n.Op {
	case ast.BinOpRegexp, ast.BinOpJsonExtract, ast.BinOpJsonUnquote:
		closing := ")"
		switch n.Op {
		case ast.BinOpRegexp:
			b.WriteString("regexp_like(")
		case ast.BinOpJsonExtract:
			b.WriteString("json_extract(")
		default:
			b.WriteString("json_unquote(json_extract(")
			closing = "))"
		}
		writeExpr(b, n.Left)
		b.WriteString(",")
		writeExpr(b, n.Right)
		b.WriteString(closing)
		return
	}

	left := n.Left
	right := n.Right
	// MySQL normalizes INTERVAL + expr to expr + INTERVAL (interval on the right)
	if _, ok := left.(*ast.IntervalExpr); ok {
		if _, ok2 := right.(*ast.IntervalExpr); !ok2 {
			left, right = right, left
		}
	}
	b.WriteString("(")
	writeExpr(b, left)
	b.WriteString(" ")
	b.WriteString(binaryOpToString(n.Op))
	b.WriteString(" ")
	writeExpr(b, right)
	b.WriteString(")")
}

func deparseColumnRef(n *ast.ColumnRef) string {
	if n.Schema != "" {
		return quoteIdent(n.Schema) + "." + quoteIdent(n.Table) + "." + quoteIdent(n.Column)
	}
	if n.Table != "" {
		return quoteIdent(n.Table) + "." + quoteIdent(n.Column)
	}
	return quoteIdent(n.Column)
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

// writeWindowDefAlias formats a window definition for auto-alias purposes.
// MySQL 8.0 uses: OVER (ORDER BY col) in the alias text, or OVER w for named windows.
func writeWindowDefAlias(b *strings.Builder, wd *ast.WindowDef) {
	// Named window reference: OVER window_name
	if wd.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil {
		b.WriteString("OVER ")
		b.WriteString(wd.RefName)
		return
	}

	b.WriteString("OVER (")

	needSpace := false
	if len(wd.PartitionBy) > 0 {
		b.WriteString("PARTITION BY ")
		for i, expr := range wd.PartitionBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, expr)
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
			writeExprAlias(b, item.Expr)
			if item.Desc {
				b.WriteString(" desc")
			}
		}
		needSpace = true
	}
	if wd.Frame != nil {
		if needSpace {
			b.WriteString(" ")
		}
		writeWindowFrame(b, wd.Frame)
	}
	b.WriteString(")")
}

func writeUnaryExpr(b *strings.Builder, n *ast.UnaryExpr) {
	switch n.Op {
	case ast.UnaryMinus:
		// MySQL parenthesizes the operand of unary minus, literals included:
		// SELECT -5 stores as -(5) (oracle 8.0.32 + 5.7.25 — views, generated
		// columns, DEFAULT expressions, CHECK constraints, and the 8.0
		// partition-function printer all agree).
		b.WriteString("-(")
		writeExpr(b, n.Operand)
		b.WriteString(")")
	case ast.UnaryNot:
		b.WriteString("(not(")
		writeExpr(b, n.Operand)
		b.WriteString("))")
	case ast.UnaryBitNot:
		b.WriteString("~(")
		writeExpr(b, n.Operand)
		b.WriteString(")")
	default:
		// MySQL drops unary plus entirely
		writeExpr(b, n.Operand)
	}
}

func writeInExpr(b *strings.Builder, n *ast.InExpr) {
	keyword := " in ("
	if n.Not {
		keyword = " not in ("
	}

	// IN subquery: a IN (SELECT ...)
	// MySQL 8.0 does NOT wrap IN subquery in outer parens (unlike IN value list).
	// MySQL 8.0 also omits AS aliases in the subquery's target list.
	if n.Select != nil {
		writeExpr(b, n.Expr)
		b.WriteString(keyword)
		writeSelectStmt(b, n.Select, true)
		b.WriteString(")")
		return
	}

	// Build the value list with no spaces after commas
	b.WriteString("(")
	writeExpr(b, n.Expr)
	b.WriteString(keyword)
	for i, item := range n.List {
		if i > 0 {
			b.WriteString(",")
		}
		writeExpr(b, item)
	}
	b.WriteString("))")
}

func writeBetweenExpr(b *strings.Builder, n *ast.BetweenExpr) {
	b.WriteString("(")
	writeExpr(b, n.Expr)
	if n.Not {
		b.WriteString(" not between ")
	} else {
		b.WriteString(" between ")
	}
	writeExpr(b, n.Low)
	b.WriteString(" and ")
	writeExpr(b, n.High)
	b.WriteString(")")
}

func writeLikeExpr(b *strings.Builder, n *ast.LikeExpr) {
	if n.Not {
		b.WriteString("(not(")
	}
	b.WriteString("(")
	writeExpr(b, n.Expr)
	b.WriteString(" like ")
	writeExpr(b, n.Pattern)
	if n.Escape != nil {
		b.WriteString(" escape ")
		writeExpr(b, n.Escape)
	}
	b.WriteString(")")
	if n.Not {
		b.WriteString("))")
	}
}

func writeIsExpr(b *strings.Builder, n *ast.IsExpr) {
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
	b.WriteString("(")
	writeExpr(b, n.Expr)
	b.WriteString(" ")
	b.WriteString(test)
	b.WriteString(")")
}

func writeRowExpr(b *strings.Builder, n *ast.RowExpr) {
	b.WriteString("row(")
	for i, item := range n.Items {
		if i > 0 {
			b.WriteString(",")
		}
		writeExpr(b, item)
	}
	b.WriteString(")")
}

func writeCaseExpr(b *strings.Builder, n *ast.CaseExpr) {
	b.WriteString("(case")
	if n.Operand != nil {
		b.WriteString(" ")
		writeExpr(b, n.Operand)
	}
	for _, w := range n.Whens {
		b.WriteString(" when ")
		writeExpr(b, w.Cond)
		b.WriteString(" then ")
		writeExpr(b, w.Result)
	}
	if n.Default != nil {
		b.WriteString(" else ")
		writeExpr(b, n.Default)
	}
	b.WriteString(" end)")
}

func writeCastExpr(b *strings.Builder, n *ast.CastExpr) {
	b.WriteString("cast(")
	writeExpr(b, n.Expr)
	b.WriteString(" as ")
	b.WriteString(deparseDataType(n.TypeName))
	// Multi-valued-index form CAST(expr AS type ARRAY); MySQL stores the trailing
	// ARRAY in lowercase in SHOW CREATE TABLE, so emit it to keep the round-trip stable.
	if n.Array {
		b.WriteString(" array")
	}
	b.WriteString(")")
}

func writeConvertExpr(b *strings.Builder, n *ast.ConvertExpr) {
	// CONVERT(expr USING charset) form
	if n.Charset != "" {
		b.WriteString("convert(")
		writeExpr(b, n.Expr)
		b.WriteString(" using ")
		b.WriteString(strings.ToLower(n.Charset))
		b.WriteString(")")
		return
	}
	// CONVERT(expr, type) form — MySQL rewrites to CAST
	b.WriteString("cast(")
	writeExpr(b, n.Expr)
	b.WriteString(" as ")
	b.WriteString(deparseDataType(n.TypeName))
	b.WriteString(")")
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

func writeIntervalExpr(b *strings.Builder, n *ast.IntervalExpr) {
	b.WriteString("interval ")
	writeExpr(b, n.Value)
	b.WriteString(" ")
	b.WriteString(strings.ToLower(n.Unit))
}

// writeExtractExpr renders EXTRACT in the engine-stored form:
// extract(day_hour from `t`.`a`) — lowercase unit, never quoted
// (oracle 8.0.32 + 5.7.25, all simple and compound units).
func writeExtractExpr(b *strings.Builder, n *ast.ExtractExpr) {
	b.WriteString("extract(")
	b.WriteString(strings.ToLower(n.Unit))
	b.WriteString(" from ")
	writeExpr(b, n.Expr)
	b.WriteString(")")
}

// writeWeightStringExpr renders WEIGHT_STRING in the engine-stored form:
// weight_string('ab' as char(4)) on 8.0, with 5.7 readbacks additionally
// carrying a level list — weight_string('ab' as char(4) level 1)
// (oracle 8.0.32 + 5.7.25).
func writeWeightStringExpr(b *strings.Builder, n *ast.WeightStringExpr) {
	b.WriteString("weight_string(")
	writeExpr(b, n.Expr)
	if n.AsChar != nil {
		fmt.Fprintf(b, " as char(%d)", n.AsChar.Length)
	}
	for i, lv := range n.Levels {
		if i == 0 {
			b.WriteString(" level ")
		} else {
			b.WriteString(",")
		}
		fmt.Fprintf(b, "%d", lv.Level)
		if lv.Desc {
			b.WriteString(" desc")
		}
		if lv.Reverse {
			b.WriteString(" reverse")
		}
	}
	b.WriteString(")")
}

func writeCollateExpr(b *strings.Builder, n *ast.CollateExpr) {
	b.WriteString("(")
	writeExpr(b, n.Expr)
	b.WriteString(" collate ")
	b.WriteString(n.Collation)
	b.WriteString(")")
}

// funcNameRewrites maps uppercase function names to their MySQL 8.0 canonical forms.
// These rewrites are applied by SHOW CREATE VIEW in MySQL 8.0.
var funcNameRewrites = map[string]string{
	"SUBSTRING":         "substr",
	"CURRENT_TIMESTAMP": "now",
	"CURRENT_DATE":      "curdate",
	"CURRENT_TIME":      "curtime",
	"CURRENT_USER":      "current_user",
	"NOW":               "now",
	"LOCALTIME":         "now",
	"LOCALTIMESTAMP":    "now",
}

// writeTrimDirectional handles TRIM(LEADING|TRAILING|BOTH remstr FROM str).
// MySQL 8.0 SHOW CREATE VIEW format: trim(leading 'x' from `a`)
func writeTrimDirectional(b *strings.Builder, direction string, args []ast.ExprNode) {
	if len(args) == 2 {
		b.WriteString("trim(")
		b.WriteString(direction)
		b.WriteString(" ")
		writeExpr(b, args[0])
		b.WriteString(" from ")
		writeExpr(b, args[1])
		b.WriteString(")")
		return
	}
	// Fallback: single arg (shouldn't happen for directional, but be safe)
	if len(args) == 1 {
		b.WriteString("trim(")
		b.WriteString(direction)
		b.WriteString(" ")
		writeExpr(b, args[0])
		b.WriteString(")")
		return
	}
	b.WriteString("trim()")
}

func writeFuncCallExpr(b *strings.Builder, n *ast.FuncCallExpr) {
	// Handle TRIM special forms: TRIM_LEADING, TRIM_TRAILING, TRIM_BOTH
	// Parser encodes these as FuncCallExpr with Name="TRIM_LEADING" etc.
	// Args: [remstr, str] for directional forms
	name := strings.ToUpper(n.Name)
	switch name {
	case "TRIM_LEADING":
		writeTrimDirectional(b, "leading", n.Args)
		return
	case "TRIM_TRAILING":
		writeTrimDirectional(b, "trailing", n.Args)
		return
	case "TRIM_BOTH":
		writeTrimDirectional(b, "both", n.Args)
		return
	case "TRIM_FROM":
		// TRIM(remstr FROM str) without a direction keyword — the engine
		// stores this form verbatim, not as comma arguments
		// (oracle 8.0.32 + 5.7.25: trim(' x' from 'axa')).
		if len(n.Args) == 2 {
			b.WriteString("trim(")
			writeExpr(b, n.Args[0])
			b.WriteString(" from ")
			writeExpr(b, n.Args[1])
			b.WriteString(")")
			return
		}
	case "GET_FORMAT":
		// Stored with a space after the keyword argument's comma:
		// get_format(DATE, 'USA') (oracle 8.0.32 + 5.7.25). Only the bare
		// builtin — a schema-qualified name is a stored function and takes
		// the generic path, keeping its qualifier.
		if n.Schema == "" && len(n.Args) == 2 {
			b.WriteString("get_format(")
			writeExpr(b, n.Args[0])
			b.WriteString(", ")
			writeExpr(b, n.Args[1])
			b.WriteString(")")
			return
		}
	}

	// GROUP_CONCAT has special formatting
	if name == "GROUP_CONCAT" {
		writeGroupConcat(b, n)
		return
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

	b.WriteString(canonical)
	switch {
	case len(n.Args) == 0 && !n.Star:
		// Zero-arg functions (CURRENT_TIMESTAMP, NOW(), etc.) — always emit parens
		b.WriteString("()")
	case n.Star:
		// COUNT(*) — MySQL 8.0 rewrites COUNT(*) to count(0)
		b.WriteString("(0)")
	default:
		// Build argument list with no spaces after commas
		b.WriteString("(")
		if n.Distinct {
			b.WriteString("distinct ")
		}
		for i, arg := range n.Args {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, arg)
		}
		b.WriteString(")")
	}

	// Append OVER clause for window functions
	if n.Over != nil {
		b.WriteString(" ")
		writeWindowDef(b, n.Over)
	}
}

// writeWindowDef formats a window definition.
// MySQL 8.0 format: OVER (PARTITION BY ... ORDER BY ... frame_clause )
// Note: trailing space before closing paren, uppercase keywords.
func writeWindowDef(b *strings.Builder, wd *ast.WindowDef) {
	// Named window reference: OVER `window_name`
	if wd.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil {
		b.WriteString("OVER " + quoteIdent(wd.RefName))
		return
	}

	b.WriteString("OVER ")
	writeWindowBody(b, wd)
}

// writeWindowBody formats the body of a window definition: (PARTITION BY ... ORDER BY ... frame_clause )
// Used by both OVER clauses and WINDOW named window definitions.
func writeWindowBody(b *strings.Builder, wd *ast.WindowDef) {
	b.WriteString("(")

	needSpace := false

	// PARTITION BY
	if len(wd.PartitionBy) > 0 {
		b.WriteString("PARTITION BY ")
		for i, expr := range wd.PartitionBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, expr)
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
			writeExpr(b, item.Expr)
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
		writeWindowFrame(b, wd.Frame)
		// When frame clause is present, MySQL 8.0 has no trailing space before )
		b.WriteString(")")
	} else {
		// No frame clause — MySQL 8.0 puts a trailing space before )
		b.WriteString(" )")
	}
}

// writeWindowFrame formats a window frame specification.
// MySQL 8.0 format: ROWS/RANGE/GROUPS BETWEEN start AND end (all uppercase).
func writeWindowFrame(b *strings.Builder, f *ast.WindowFrame) {
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
		writeWindowFrameBound(b, f.Start)
		b.WriteString(" AND ")
		writeWindowFrameBound(b, f.End)
	} else {
		// Single bound form
		b.WriteString(" ")
		writeWindowFrameBound(b, f.Start)
	}
}

// writeWindowFrameBound formats a window frame bound.
func writeWindowFrameBound(b *strings.Builder, fb *ast.WindowFrameBound) {
	switch fb.Type {
	case ast.BoundUnboundedPreceding:
		b.WriteString("UNBOUNDED PRECEDING")
	case ast.BoundPreceding:
		writeExpr(b, fb.Offset)
		b.WriteString(" PRECEDING")
	case ast.BoundCurrentRow:
		b.WriteString("CURRENT ROW")
	case ast.BoundFollowing:
		writeExpr(b, fb.Offset)
		b.WriteString(" FOLLOWING")
	case ast.BoundUnboundedFollowing:
		b.WriteString("UNBOUNDED FOLLOWING")
	default:
		b.WriteString("/* unknown bound */")
	}
}

// writeExistsExpr formats an EXISTS expression.
// MySQL 8.0 format: exists(select ...) — column aliases are omitted inside EXISTS.
func writeExistsExpr(b *strings.Builder, n *ast.ExistsExpr) {
	if n.Select != nil {
		b.WriteString("exists(")
		writeSelectStmt(b, n.Select, true)
		b.WriteString(")")
		return
	}
	b.WriteString("exists(/* subquery */)")
}

// writeSubqueryExpr formats a subquery expression.
// MySQL 8.0 format: (select ...) — scalar subqueries in expression context
// omit column aliases (AS `...`).
func writeSubqueryExpr(b *strings.Builder, n *ast.SubqueryExpr) {
	if n.Select != nil {
		b.WriteString("(")
		writeSelectStmt(b, n.Select, true)
		b.WriteString(")")
		return
	}
	b.WriteString("(/* subquery */)")
}

// writeGroupConcatAlias generates the MySQL 8.0 auto-alias for GROUP_CONCAT.
// MySQL 8.0 alias format: GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';')
// Uses uppercase keywords and original column names (not table-qualified).
func writeGroupConcatAlias(b *strings.Builder, n *ast.FuncCallExpr) {
	b.WriteString("GROUP_CONCAT(")

	// DISTINCT
	if n.Distinct {
		b.WriteString("DISTINCT ")
	}

	// Arguments
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		writeExprAlias(b, arg)
	}

	// ORDER BY — MySQL 8.0 alias omits ASC (default), only shows DESC
	if len(n.OrderBy) > 0 {
		b.WriteString(" ORDER BY ")
		for i, item := range n.OrderBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, item.Expr)
			if item.Desc {
				b.WriteString(" DESC")
			}
		}
	}

	// SEPARATOR
	b.WriteString(" SEPARATOR ")
	if n.Separator != nil {
		writeExprAlias(b, n.Separator)
	} else {
		b.WriteString("','")
	}

	b.WriteString(")")
}

// writeGroupConcat handles GROUP_CONCAT with its special syntax:
// group_concat([distinct] expr [order by expr ASC|DESC] separator 'str')
// MySQL 8.0 always shows the separator (default ',') and explicit ASC in ORDER BY.
func writeGroupConcat(b *strings.Builder, n *ast.FuncCallExpr) {
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
		writeExpr(b, arg)
	}

	// ORDER BY — MySQL 8.0 always shows explicit ASC/DESC
	if len(n.OrderBy) > 0 {
		b.WriteString(" order by ")
		for i, item := range n.OrderBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, item.Expr)
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
		writeExpr(b, n.Separator)
	} else {
		b.WriteString("','")
	}

	b.WriteString(")")
}

// quoteIdent quotes an identifier in backticks, doubling any backtick inside
// it, as MySQL writes identifiers in SHOW CREATE output. A name can hold one:
// a quoted identifier written with “, or an unaliased item named after its
// source text, such as `a`+1.
func quoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
