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

	ast "github.com/bytebase/omni/mariadb/ast"
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
		writeParenSource(b, stmt)
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

	// CTE (WITH clause) — emit before the SELECT keyword. MariaDB writes no
	// separator between the clause and the query it heads:
	// `with cte as (...)select ...`. Measured on mariadb:11.8.8.
	if len(stmt.CTEs) > 0 {
		writeCTEs(b, stmt.CTEs)
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
		for i, item := range stmt.GroupBy {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, item.Expr)
			if item.Desc {
				b.WriteString(" desc")
			}
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
		writeLimit(b, stmt.Limit)
	}

	// FOR UPDATE / FOR SHARE / LOCK IN SHARE MODE
	if stmt.ForUpdate != nil {
		b.WriteString(" ")
		b.WriteString(deparseForUpdate(stmt.ForUpdate))
	}
}

// deparseForUpdate formats a locking clause the way MariaDB's SHOW CREATE
// VIEW does. MariaDB has no FOR SHARE, so a shared lock is always written
// LOCK IN SHARE MODE:
//   - for update
//   - lock in share mode
//   - for update nowait
//   - for update wait 5
//   - for update skip locked
func deparseForUpdate(fu *ast.ForUpdate) string {
	var b strings.Builder
	if fu.LockInShareMode {
		b.WriteString("lock in share mode")
	} else {
		b.WriteString("for update")
	}

	// NOWAIT / WAIT n / SKIP LOCKED
	switch {
	case fu.NoWait:
		b.WriteString(" nowait")
	case fu.HasWait:
		b.WriteString(" wait ")
		b.WriteString(fu.WaitSeconds)
	case fu.SkipLocked:
		b.WriteString(" skip locked")
	}

	return b.String()
}

// writeParenSource formats a parenthesized query expression: an optional WITH,
// '(' inner ')', and the OUTER trailing ORDER BY / LIMIT held on the wrapper.
// The inner query keeps its own clauses, so both scopes round-trip (e.g.
// "(select 1 limit 5) limit 2").
func writeParenSource(b *strings.Builder, stmt *ast.SelectStmt) {
	if len(stmt.CTEs) > 0 {
		writeCTEs(b, stmt.CTEs)
	}

	b.WriteString("(")
	writeSelectStmt(b, stmt.ParenSource, false)
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
		writeLimit(b, stmt.Limit)
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
	// MariaDB has no TABLE statement, so this writer is unreachable; it still
	// writes the locking clause rather than disagree with the VALUES writer.
	if stmt.ForUpdate != nil {
		b.WriteString(" ")
		b.WriteString(deparseForUpdate(stmt.ForUpdate))
	}
	return b.String()
}

// deparseValuesQueryPrimary formats a table value constructor and the
// clauses after it as MariaDB's SHOW CREATE VIEW does:
// values (1,2),(3,4) order by 1 limit 1 for update.
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
		b.WriteString("(")
		for j, expr := range row {
			if j > 0 {
				b.WriteString(",")
			}
			b.WriteString(deparseExpr(expr))
		}
		b.WriteString(")")
	}
	appendOrderByLimit(&b, stmt.OrderBy, stmt.Limit)
	if stmt.ForUpdate != nil {
		b.WriteString(" ")
		b.WriteString(deparseForUpdate(stmt.ForUpdate))
	}
	return b.String()
}

// limitHasParamMarker reports whether the LIMIT's count or offset is a `?`
// parameter marker.
func limitHasParamMarker(limit *ast.Limit) bool {
	if _, ok := limit.Count.(*ast.ParamMarkerExpr); ok {
		return true
	}
	_, ok := limit.Offset.(*ast.ParamMarkerExpr)
	return ok
}

// writeLimit emits a LIMIT clause (lowercase). The MySQL comma form
// `limit offset,count` swaps the two operands, which is harmless for literals
// but reorders the positional bindings of `?` parameter markers; when a marker
// is present and the source used the OFFSET keyword, the original
// `limit count offset offset` form is reproduced so the binding order is kept.
func writeLimit(b *strings.Builder, limit *ast.Limit) {
	b.WriteString(" limit ")
	if limit.Offset != nil {
		if limit.OffsetKeyword && limitHasParamMarker(limit) {
			writeExpr(b, limit.Count)
			b.WriteString(" offset ")
			writeExpr(b, limit.Offset)
			return
		}
		writeExpr(b, limit.Offset)
		b.WriteString(",")
	}
	writeExpr(b, limit.Count)
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
			if limit.OffsetKeyword && limitHasParamMarker(limit) {
				b.WriteString(deparseExpr(limit.Count))
				b.WriteString(" offset ")
				b.WriteString(deparseExpr(limit.Offset))
				return
			}
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
		writeLimit(b, stmt.Limit)
	}
}

func deparseCTEs(ctes []*ast.CommonTableExpr) string {
	var b strings.Builder
	writeCTEs(&b, ctes)
	return b.String()
}

// writeCTEs formats a WITH clause (one or more CTEs).
// MariaDB format: with [recursive] name [(`col`, ...)] as (select ...) [, ...]
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
		// MariaDB writes a CTE's name bare, where MySQL 8.0 backticks it.
		// Measured on mariadb:11.8.8 through information_schema.VIEWS.
		b.WriteString(cte.Name)

		// Column list, written against the name with no space:
		// `with cte(`x`,`y`) as (...)`. Measured on mariadb:11.8.8.
		if len(cte.Columns) > 0 {
			b.WriteString("(")
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

	// MariaDB writes one space before AS whatever the expression is; MySQL 8.0
	// writes two after a window function, whose OVER clause it ends with
	// " )". Measured on mariadb:11.8.8.
	b.WriteString(" AS ")
	b.WriteString(quoteIdent(alias))
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
	case *ast.ParamMarkerExpr:
		b.WriteByte('?')
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
		}
		if trimHead != "" && len(n.Args) == 2 {
			b.WriteString(trimHead)
			writeExprAlias(b, n.Args[0])
			b.WriteString(" FROM ")
			writeExprAlias(b, n.Args[1])
			b.WriteString(")")
			return
		}
		if trimHead != "" && len(n.Args) == 1 {
			// TRIM(LEADING FROM str): the name is the text as written.
			b.WriteString(trimHead)
			b.WriteString("FROM ")
			writeExprAlias(b, n.Args[0])
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
		b.WriteString("(")
		// The engine names an unaliased item after its own text, so the
		// quantifier appears as written: SUM(DISTINCTROW ALL c1) is the column
		// name mysql:8.0, mariadb:11.8.8 and TiDB v8.5.5 all give.
		if q := n.Quantifier.String(); q != "" {
			b.WriteString(q)
			b.WriteString(" ")
		}
		for i, arg := range n.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, arg)
		}
		if n.Using != "" {
			// The name keeps CHAR's USING clause as written, whatever the
			// stored body does with it.
			b.WriteString(" USING ")
			b.WriteString(n.Using)
		}
		if n.WeightAs != "" {
			fmt.Fprintf(b, " AS %s(%d)", n.WeightAs, n.WeightAsLen)
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
		b.WriteString("CAST(")
		writeExprAlias(b, n.Expr)
		b.WriteString(" AS ")
		b.WriteString(deparseDataTypeAlias(n.TypeName))
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
			if n.Quantifier != "" {
				b.WriteString(n.Quantifier)
				b.WriteString(" ")
			}
			b.WriteString("(")
			writeSelectStmtAlias(b, n.Select)
			b.WriteString(")")
			return
		}
		b.WriteString("(/* subquery */)")
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
	case *ast.ExtractExpr:
		// Alias approximation with the uppercase-keyword convention:
		// "EXTRACT(DAY_HOUR FROM a)".
		b.WriteString("EXTRACT(")
		b.WriteString(n.Unit)
		b.WriteString(" FROM ")
		writeExprAlias(b, n.Expr)
		b.WriteString(")")
	case *ast.PositionExpr:
		b.WriteString("POSITION(")
		writeExprAlias(b, n.Substr)
		b.WriteString(" IN ")
		writeExprAlias(b, n.Str)
		b.WriteString(")")
	case *ast.KeywordArg:
		b.WriteString(n.Keyword)
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
			b.WriteString("(")
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
		for i, item := range stmt.GroupBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeExprAlias(b, item.Expr)
			if item.Desc {
				b.WriteString(" DESC")
			}
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
			if limit.OffsetKeyword && limitHasParamMarker(limit) {
				writeExprAlias(b, limit.Count)
				b.WriteString(" OFFSET ")
				writeExprAlias(b, limit.Offset)
				return
			}
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
	// A flat `and` chain with no parentheses of its own, as the server writes
	// it: on(`t1`.`a` = `t3`.`a` and `t1`.`b` = `t3`.`b`). Measured on
	// mariadb:11.8.8 over `t1 JOIN t3 USING (a,b)`.
	for i, col := range columns {
		if i > 0 {
			b.WriteString(" and ")
		}
		b.WriteString(quoteIdent(leftName) + "." + quoteIdent(col) + " = " + quoteIdent(rightName) + "." + quoteIdent(col))
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
	if t.SystemTime != nil {
		b.WriteString(" ")
		b.WriteString(deparseSystemTime(t.SystemTime))
	}
	if t.Alias != "" {
		b.WriteString(" " + quoteIdent(t.Alias))
	}
	return b.String()
}

// deparseSystemTime formats the FOR SYSTEM_TIME temporal clause. Dropping it
// would silently turn a time-travel query into a current-data query.
func deparseSystemTime(st *ast.SystemTime) string {
	switch st.Kind {
	case ast.SystemTimeAsOf:
		return "for system_time as of " + deparseExpr(st.From)
	case ast.SystemTimeAsOfTransaction:
		return "for system_time as of transaction " + deparseExpr(st.From)
	case ast.SystemTimeBetween:
		return "for system_time between " + systemTimeBound(st.From, st.FromTransaction) +
			" and " + systemTimeBound(st.To, st.ToTransaction)
	case ast.SystemTimeFromTo:
		return "for system_time from " + systemTimeBound(st.From, st.FromTransaction) +
			" to " + systemTimeBound(st.To, st.ToTransaction)
	case ast.SystemTimeAll:
		return "for system_time all"
	default:
		return ""
	}
}

// systemTimeBound renders a FOR SYSTEM_TIME range bound, prefixing TRANSACTION
// when the bound carries the transaction-id qualifier.
func systemTimeBound(expr ast.ExprNode, transaction bool) string {
	if transaction {
		return "transaction " + deparseExpr(expr)
	}
	return deparseExpr(expr)
}

func writeExpr(b *strings.Builder, node ast.ExprNode) {
	switch n := node.(type) {
	case *ast.IntLit:
		fmt.Fprintf(b, "%d", n.Value)
	case *ast.FloatLit:
		b.WriteString(n.Value)
	case *ast.BoolLit:
		// MariaDB stores TRUE and FALSE as 1 and 0. Measured on mariadb:11.8.8.
		if n.Value {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	case *ast.StringLit:
		b.WriteString(deparseStringLit(n))
	case *ast.NullLit:
		b.WriteString("NULL")
	case *ast.ParamMarkerExpr:
		b.WriteByte('?')
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
	case *ast.ExtractExpr:
		writeExtractExpr(b, n)
	case *ast.PositionExpr:
		writePositionExpr(b, n)
	case *ast.KeywordArg:
		// Bare keyword argument (timestampdiff unit, get_format type) —
		// never quoted: the engine stores timestampdiff(SECOND,...) and
		// rejects a backtick-quoted unit.
		b.WriteString(n.Keyword)
	case *ast.CollateExpr:
		writeCollateExpr(b, n)
	case *ast.FuncCallExpr:
		writeFuncCallExpr(b, n)
	case *ast.ExistsExpr:
		writeExistsExpr(b, n)
	case *ast.SubqueryExpr:
		writeSubqueryExpr(b, n)
	case *ast.IgnoreExpr:
		// MariaDB's value for "leave this column alone", in a VALUES row or an
		// assignment.
		b.WriteString("ignore")
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
	// MariaDB keeps the spelling the literal was written in: `0xFF` comes back
	// as 0xff and `X'FF'` as X'ff', where MySQL 8.0 normalizes both to the 0x
	// form. Measured on mariadb:11.8.8. HexLit.Value is either "0xFF" (0x
	// prefix form) or "FF" (X'' form).
	val := n.Value
	if strings.HasPrefix(val, "0x") || strings.HasPrefix(val, "0X") {
		return "0x" + strings.ToLower(val[2:])
	}
	return "X'" + strings.ToLower(val) + "'"
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

// A quantified comparison, and an IN over a subquery, print with no parentheses of
// their own, so whether one needs them depends on where it stands. The grammar reads
// `a > ALL (q)` as a boolean primary and `a IN (q)` as a predicate, and neither is a
// bit_expr: `(a > ALL (q)) + 1` printed bare is a syntax error, and `1 = (a > ALL (q))`
// printed bare parses as `(1 = a) > ALL (q)`.
//
// MariaDB 11.8 parenthesizes where the grammar needs it, as this writer does.
const (
	printsDelimited = iota // the writer already delimits it
	printsPredicate        // `a in (q)`
	printsBoolPri          // `a > all (q)`
)

// bareLevel reports what node prints as, looking through what the writer drops.
func bareLevel(node ast.ExprNode) int {
	for {
		switch n := node.(type) {
		case *ast.ParenExpr:
			node = n.Expr
			continue
		case *ast.UnaryExpr:
			if n.Op == ast.UnaryPlus {
				node = n.Operand
				continue
			}
		case *ast.InExpr:
			if n.Select != nil {
				return printsPredicate
			}
		case *ast.BinaryExpr:
			sub, ok := n.Right.(*ast.SubqueryExpr)
			if !ok || sub.Quantifier == "" || sub.Select == nil {
				return printsDelimited
			}
			all := sub.Quantifier == "ALL"
			switch {
			case n.Op == ast.BinOpEq && !all:
				return printsPredicate
			case n.Op == ast.BinOpNe && all:
				return printsDelimited
			}
			return printsBoolPri
		}
		return printsDelimited
	}
}

// MariaDB's operator precedence, which its own printer uses to decide where a
// parenthesis is needed; higher binds tighter. MySQL 8.0 instead parenthesizes
// every operator it prints, which is why this writer used to.
//
// Measured on mariadb:11.8.8 through information_schema.VIEWS: `(a + b) * 2`
// keeps the parentheses the multiplication needs, `a + b * 2` needs none,
// `a = 1 and b = 2 and c = 'y'` prints flat, `(a > 0 or b < 1) and c = 'x'`
// keeps only the one the OR needs, and `a - (b - 1)` keeps the one the
// left-associative minus needs on its right.
//
// The predicates bind tighter than a comparison, and IN and LIKE tighter than
// BETWEEN, which is what `(a = 1) in (1,2)` keeping its parentheses and
// `a between 1 and 2 = 1` losing them say.
// https://mariadb.com/docs/server/reference/sql-statements/operators/operator-precedence
const (
	precLowest = iota
	precOr
	precXor
	precAnd
	precNot
	precComparison
	precBetween
	precIn
	precBitOr
	precBitAnd
	precShift
	precAdd
	precMul
	precBitXor
	precUnary
	precBang
	precCollate
	// precAtom is what the writer delimits itself: a literal, a column, a
	// function call, a parenthesized subquery.
	precAtom
)

// binOpPrec returns op's precedence. An operator this writer rewrites into a
// function call (REGEXP, ->, ->>) prints delimited, so it is an atom.
func binOpPrec(op ast.BinaryOp) int {
	switch op {
	case ast.BinOpOr:
		return precOr
	case ast.BinOpXor:
		return precXor
	case ast.BinOpAnd:
		return precAnd
	case ast.BinOpEq, ast.BinOpNe, ast.BinOpLt, ast.BinOpGt, ast.BinOpLe, ast.BinOpGe,
		ast.BinOpNullSafeEq, ast.BinOpSoundsLike:
		return precComparison
	case ast.BinOpBitOr:
		return precBitOr
	case ast.BinOpBitAnd:
		return precBitAnd
	case ast.BinOpShiftLeft, ast.BinOpShiftRight:
		return precShift
	case ast.BinOpAdd, ast.BinOpSub:
		return precAdd
	case ast.BinOpMul, ast.BinOpDiv, ast.BinOpMod, ast.BinOpDivInt:
		return precMul
	case ast.BinOpBitXor:
		return precBitXor
	}
	return precAtom
}

// exprPrec returns the precedence node prints at: how tightly its outermost
// operator binds, or precAtom for one the writer delimits itself. A ParenExpr
// is looked through, since writeExpr drops it and this writer puts a
// parenthesis back where the precedence needs one.
func exprPrec(node ast.ExprNode) int {
	switch n := node.(type) {
	case *ast.ParenExpr:
		return exprPrec(n.Expr)
	case *ast.BinaryExpr:
		if sub, ok := n.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" && sub.Select != nil {
			// The server stores `a = ANY (q)` as `a in (q)` and `a <> ALL (q)`
			// as `!(a in (q))`, so each prints at its own form's precedence:
			// `1 = (a = ANY (q))` comes back as `1 = a in (q)`. Measured on
			// mariadb:11.8.8.
			switch {
			case n.Op == ast.BinOpEq && sub.Quantifier != "ALL":
				return precIn
			case n.Op == ast.BinOpNe && sub.Quantifier == "ALL":
				return precBang
			}
			return precComparison
		}
		switch n.Op {
		case ast.BinOpRegexp:
			return precComparison
		case ast.BinOpJsonExtract, ast.BinOpJsonUnquote:
			return precAtom
		}
		return binOpPrec(n.Op)
	case *ast.UnaryExpr:
		switch n.Op {
		case ast.UnaryMinus, ast.UnaryBitNot:
			return precUnary
		case ast.UnaryNot:
			// Printed as `!`, which binds tighter than the unary minus.
			return precBang
		case ast.UnaryPlus:
			// Dropped, so the operand is what prints.
			return exprPrec(n.Operand)
		}
	case *ast.IsExpr:
		return precComparison
	case *ast.InExpr, *ast.LikeExpr:
		return precIn
	case *ast.BetweenExpr:
		return precBetween
	case *ast.CollateExpr:
		return precCollate
	}
	return precAtom
}

// writeOperand writes node as the operand of an operator whose operand slot
// reaches precedence min, parenthesizing it where MariaDB would.
func writeOperand(b *strings.Builder, node ast.ExprNode, min int) {
	if exprPrec(node) < min {
		b.WriteString("(")
		writeExpr(b, node)
		b.WriteString(")")
		return
	}
	writeExpr(b, node)
}

// writeBitOperand writes an operand the grammar reads as a bit_expr.
func writeBitOperand(b *strings.Builder, node ast.ExprNode) {
	writeBitOperandAt(b, node, precLowest)
}

// writeBitOperandAt writes a bit_expr operand of an operator reaching min.
func writeBitOperandAt(b *strings.Builder, node ast.ExprNode, min int) {
	if bareLevel(node) == printsDelimited && exprPrec(node) >= min {
		writeExpr(b, node)
		return
	}
	b.WriteString("(")
	writeExpr(b, node)
	b.WriteString(")")
}

// writeComparisonRight writes the right side of a comparison, which the grammar reads
// as a predicate.
func writeComparisonRight(b *strings.Builder, node ast.ExprNode) {
	writeComparisonRightAt(b, node, precLowest)
}

// writeComparisonRightAt writes a predicate operand of an operator reaching min.
func writeComparisonRightAt(b *strings.Builder, node ast.ExprNode, min int) {
	if bareLevel(node) != printsBoolPri && exprPrec(node) >= min {
		writeExpr(b, node)
		return
	}
	b.WriteString("(")
	writeExpr(b, node)
	b.WriteString(")")
}

// isComparison reports whether op is one of the grammar's comp_op operators.
func isComparison(op ast.BinaryOp) bool {
	switch op {
	case ast.BinOpEq, ast.BinOpNe, ast.BinOpLt, ast.BinOpGt, ast.BinOpLe, ast.BinOpGe, ast.BinOpNullSafeEq:
		return true
	}
	return false
}

// isLogical reports whether op takes whole expressions as operands.
func isLogical(op ast.BinaryOp) bool {
	switch op {
	case ast.BinOpAnd, ast.BinOpOr, ast.BinOpXor:
		return true
	}
	return false
}

// writeQuantifiedComparison formats `expr op {ANY|SOME|ALL} (subquery)`. The
// quantifier is the meaning of the comparison, so dropping it would turn `a > ALL
// (q)` into a scalar comparison against q.
//
// MariaDB 11.8 stores `a = ANY (q)` as `a in (q)` and `a <> ALL (q)` as `!(a in (q))`,
// SOME as any, and leaves the rest as written, with no outer parentheses.
func writeQuantifiedComparison(b *strings.Builder, n *ast.BinaryExpr, sub *ast.SubqueryExpr) {
	all := sub.Quantifier == "ALL"
	word := "any"
	if all {
		word = "all"
	}
	switch {
	case n.Op == ast.BinOpEq && !all:
		writeBitOperand(b, n.Left)
		b.WriteString(" in (")
		writeSelectStmt(b, sub.Select, true)
		b.WriteString(")")
	case n.Op == ast.BinOpNe && all:
		b.WriteString("!(")
		writeBitOperand(b, n.Left)
		b.WriteString(" in (")
		writeSelectStmt(b, sub.Select, true)
		b.WriteString("))")
	default:
		writeExpr(b, n.Left)
		b.WriteString(" ")
		b.WriteString(binaryOpToString(n.Op))
		b.WriteString(" ")
		b.WriteString(word)
		b.WriteString(" (")
		writeSelectStmt(b, sub.Select, true)
		b.WriteString(")")
	}
}

func writeBinaryExpr(b *strings.Builder, n *ast.BinaryExpr) {
	if sub, ok := n.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" && sub.Select != nil {
		writeQuantifiedComparison(b, n, sub)
		return
	}

	// MariaDB keeps REGEXP an operator — `a regexp 'x'`, and `!(a regexp 'y')`
	// for NOT REGEXP — where MySQL 8.0 rewrites it to regexp_like(). Measured
	// on mariadb:11.8.8.
	if n.Op == ast.BinOpRegexp {
		writeOperand(b, n.Left, precComparison)
		b.WriteString(" regexp ")
		writeOperand(b, n.Right, precComparison+1)
		return
	}

	// Operator-to-function rewrites:
	// -> → json_extract(left, right)
	// ->> → json_unquote(json_extract(left, right))
	switch n.Op {
	case ast.BinOpJsonExtract, ast.BinOpJsonUnquote:
		closing := ")"
		switch n.Op {
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
	// Every operator here is left-associative, so an operand of the same
	// precedence needs a parenthesis on the right and none on the left:
	// `a = 1 = 1` is what `(a = 1) = 1` prints as, and `a - (b - 1)` keeps its
	// own. Measured on mariadb:11.8.8.
	own := binOpPrec(n.Op)
	if isLogical(n.Op) || isComparison(n.Op) {
		writeOperand(b, left, own)
	} else {
		writeBitOperandAt(b, left, own)
	}
	b.WriteString(" ")
	b.WriteString(binaryOpToString(n.Op))
	b.WriteString(" ")
	switch {
	case isLogical(n.Op):
		writeOperand(b, right, own+1)
	case isComparison(n.Op):
		writeComparisonRightAt(b, right, own+1)
	default:
		writeBitOperandAt(b, right, own+1)
	}
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
		// MariaDB writes MOD for both spellings: `a % b` comes back as
		// `a MOD b`. Measured on mariadb:11.8.8.
		return "MOD"
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
		// MariaDB writes the operand bare where it binds at least as tightly as
		// the unary minus: `-a`, `-5`, and `-(a + b)`. MySQL 8.0 parenthesizes
		// every non-literal operand. Measured on mariadb:11.8.8.
		b.WriteString("-")
		writeOperand(b, n.Operand, precUnary)
	case ast.UnaryNot:
		// MariaDB writes NOT as `!`, parenthesizing the operand only where the
		// precedence needs it: `!(a + 1)`, `!(a in (q))`, `!(a regexp 'y')`,
		// but `!ifnull(a,0)`. MySQL 8.0 writes `(not(...))`. Measured on
		// mariadb:11.8.8.
		b.WriteString("!")
		writeOperand(b, n.Operand, precBang)
	case ast.UnaryBitNot:
		b.WriteString("~")
		writeOperand(b, n.Operand, precUnary)
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

	// IN subquery: a IN (SELECT ...), with no outer parenthesis and no AS
	// aliases in the subquery's target list.
	if n.Select != nil {
		writeBitOperandAt(b, n.Expr, precIn+1)
		b.WriteString(keyword)
		writeSelectStmt(b, n.Select, true)
		b.WriteString(")")
		return
	}

	// Value list, with no spaces after the commas and no parenthesis around
	// the whole: `a in (1,2,3)`. MySQL 8.0 writes `(a in (1,2,3))`. Measured on
	// mariadb:11.8.8.
	writeBitOperandAt(b, n.Expr, precIn+1)
	b.WriteString(keyword)
	for i, item := range n.List {
		if i > 0 {
			b.WriteString(",")
		}
		writeExpr(b, item)
	}
	b.WriteString(")")
}

func writeBetweenExpr(b *strings.Builder, n *ast.BetweenExpr) {
	writeBitOperandAt(b, n.Expr, precBetween+1)
	if n.Not {
		b.WriteString(" not between ")
	} else {
		b.WriteString(" between ")
	}
	writeBitOperandAt(b, n.Low, precBetween+1)
	b.WriteString(" and ")
	writeBitOperandAt(b, n.High, precBetween+1)
}

// writeLikeExpr writes `a like 'x'`, and the negated form as the operator
// MariaDB stores, with the extra space it writes before the keyword:
// `a  not like 'x'`. MySQL 8.0 writes `(not((a like 'x')))`. Measured on
// mariadb:11.8.8.
func writeLikeExpr(b *strings.Builder, n *ast.LikeExpr) {
	writeBitOperandAt(b, n.Expr, precIn+1)
	if n.Not {
		b.WriteString("  not like ")
	} else {
		b.WriteString(" like ")
	}
	writeBitOperandAt(b, n.Pattern, precIn+1)
	if n.Escape != nil {
		b.WriteString(" escape ")
		writeBitOperandAt(b, n.Escape, precIn+1)
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
		// MariaDB stores IS UNKNOWN as IS NULL, which is what it means for a
		// three-valued operand. Measured on mariadb:11.8.8.
		if n.Not {
			test = "is not null"
		} else {
			test = "is null"
		}
	default:
		test = "is ?"
	}
	writeOperand(b, n.Expr, precComparison+1)
	b.WriteString(" ")
	b.WriteString(test)
}

// writeRowExpr writes a row constructor as the parenthesized list MariaDB
// stores: `ROW(a,b) = ROW(1,2)` comes back as `(`a`,`b`) = (1,2)`. Measured on
// mariadb:11.8.8.
func writeRowExpr(b *strings.Builder, n *ast.RowExpr) {
	b.WriteString("(")
	for i, item := range n.Items {
		if i > 0 {
			b.WriteString(",")
		}
		writeExpr(b, item)
	}
	b.WriteString(")")
}

// writeCaseExpr writes `case ... end`, which delimits itself: MariaDB adds no
// parentheses around it, where MySQL 8.0 does. Measured on mariadb:11.8.8.
func writeCaseExpr(b *strings.Builder, n *ast.CaseExpr) {
	b.WriteString("case")
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
	b.WriteString(" end")
}

func writeCastExpr(b *strings.Builder, n *ast.CastExpr) {
	b.WriteString("cast(")
	writeExpr(b, n.Expr)
	b.WriteString(" as ")
	b.WriteString(deparseDataType(n.TypeName))
	b.WriteString(")")
}

// writeExtractExpr renders EXTRACT in the engine-stored form:
// extract(day_hour from `t`.`a`) — lowercase unit, never quoted (oracle
// mariadb:11.8.8, all simple and compound units).
func writeExtractExpr(b *strings.Builder, n *ast.ExtractExpr) {
	b.WriteString("extract(")
	b.WriteString(strings.ToLower(n.Unit))
	b.WriteString(" from ")
	writeExpr(b, n.Expr)
	b.WriteString(")")
}

// writePositionExpr writes POSITION(substr IN str). MySQL 8.0 and MariaDB
// 11.8.8 store the call as locate(substr, str), the same function under its
// other name, and TiDB v8.5.5 stores POSITION. Keeping POSITION makes the
// text parse back to this node on every engine.
func writePositionExpr(b *strings.Builder, n *ast.PositionExpr) {
	b.WriteString("position(")
	writeExpr(b, n.Substr)
	b.WriteString(" in ")
	writeExpr(b, n.Str)
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

func writeCollateExpr(b *strings.Builder, n *ast.CollateExpr) {
	writeBitOperandAt(b, n.Expr, precCollate)
	b.WriteString(" collate ")
	b.WriteString(n.Collation)
}

// funcNameRewrites maps uppercase function names to the canonical form MariaDB
// stores. Measured on mariadb:11.8.8 through information_schema.VIEWS:
// NOW() and CURRENT_TIMESTAMP both come back as current_timestamp(), UPPER as
// ucase, LOWER as lcase, SUBSTRING as substr, CHARACTER_LENGTH as char_length
// and STDDEV as std. MySQL 8.0 picks now(), upper() and lower() instead.
var funcNameRewrites = map[string]string{
	"SUBSTRING":         "substr",
	"SUBSTR":            "substr",
	"CURRENT_TIMESTAMP": "current_timestamp",
	"CURRENT_DATE":      "curdate",
	"CURRENT_TIME":      "curtime",
	"CURRENT_USER":      "current_user",
	"NOW":               "current_timestamp",
	"LOCALTIME":         "current_timestamp",
	"LOCALTIMESTAMP":    "current_timestamp",
	"UPPER":             "ucase",
	"UCASE":             "ucase",
	"LOWER":             "lcase",
	"LCASE":             "lcase",
	"CHARACTER_LENGTH":  "char_length",
	"CHAR_LENGTH":       "char_length",
	"STDDEV":            "std",
	"STDDEV_POP":        "std",
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
	// TRIM(LEADING FROM str) trims spaces, which MariaDB stores as the
	// one-sided function: ltrim(str), rtrim(str), trim(str) (mariadb:11.8.8,
	// SHOW CREATE VIEW).
	if len(args) == 1 {
		switch direction {
		case "leading":
			b.WriteString("ltrim(")
		case "trailing":
			b.WriteString("rtrim(")
		default:
			b.WriteString("trim(")
		}
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
		// The engine keeps only whether duplicates are removed: it stores
		// sum(distinct `c1`) for every DISTINCT, DISTINCTROW and DISTINCT ALL
		// spelling and drops a bare ALL, which is the default (oracle
		// mysql:8.0, mariadb:11.8.8, TiDB v8.5.5: SHOW CREATE VIEW).
		if n.Quantifier.RemovesDuplicates() {
			b.WriteString("distinct ")
		}
		for i, arg := range n.Args {
			if i > 0 {
				b.WriteString(",")
			}
			writeExpr(b, arg)
		}
		writeCharUsing(b, n.Using)
		if n.WeightAs != "" {
			fmt.Fprintf(b, " as %s(%d)", strings.ToLower(n.WeightAs), n.WeightAsLen)
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
	// A named window reference that the rewrite could not resolve. MariaDB
	// inlines every one it stores, so this is a reference to a window the
	// statement does not define.
	if wd.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil {
		b.WriteString("over " + quoteIdent(wd.RefName))
		return
	}

	b.WriteString("over ")
	writeWindowBody(b, wd)
}

// writeWindowBody formats the body of a window definition: (PARTITION BY ... ORDER BY ... frame_clause )
// Used by both OVER clauses and WINDOW named window definitions.
// writeWindowBody formats a window body as MariaDB stores it: a space after
// the opening parenthesis, lowercase keywords, and no space before the closing
// one except the trailing one the frame clause writes. Measured on
// mariadb:11.8.8: `over ( order by \`t\`.\`a\`)` and
// `over ( order by \`t\`.\`a\` rows between  unbounded  preceding  and  current row )`.
// MySQL 8.0 writes uppercase keywords and the space before the closing
// parenthesis instead.
func writeWindowBody(b *strings.Builder, wd *ast.WindowDef) {
	b.WriteString("( ")

	needSpace := false

	// PARTITION BY
	if len(wd.PartitionBy) > 0 {
		b.WriteString("partition by ")
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
		b.WriteString("order by ")
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

	// Frame clause. It writes its own trailing space, so the closing
	// parenthesis follows whatever came last.
	if wd.Frame != nil {
		if needSpace {
			b.WriteString(" ")
		}
		writeWindowFrame(b, wd.Frame)
	}
	b.WriteString(")")
}

// writeWindowFrame formats a window frame as MariaDB stores it: lowercase, and
// always the BETWEEN form — a single bound is the start of a frame that ends at
// CURRENT ROW, and the server writes it back that way. The spacing is the
// server's own, measured on mariadb:11.8.8:
//
//	rows between  unbounded  preceding  and  current row
//	rows between 2 preceding  and  current row
//	range between  unbounded  preceding  and  unbounded  following
//
// which is two spaces before `and`, a bound that carries its own leading space
// where it has no offset, and a trailing space the caller writes nothing after.
func writeWindowFrame(b *strings.Builder, f *ast.WindowFrame) {
	switch f.Type {
	case ast.FrameRows:
		b.WriteString("rows")
	case ast.FrameRange:
		b.WriteString("range")
	case ast.FrameGroups:
		b.WriteString("groups")
	}

	b.WriteString(" between ")
	writeWindowFrameBound(b, f.Start)
	b.WriteString("  and ")
	if f.End != nil {
		writeWindowFrameBound(b, f.End)
	} else {
		b.WriteString(" current row")
	}
	b.WriteString(" ")
}

// writeWindowFrameBound formats a window frame bound.
func writeWindowFrameBound(b *strings.Builder, fb *ast.WindowFrameBound) {
	switch fb.Type {
	case ast.BoundUnboundedPreceding:
		b.WriteString(" unbounded  preceding")
	case ast.BoundPreceding:
		writeExpr(b, fb.Offset)
		b.WriteString(" preceding")
	case ast.BoundCurrentRow:
		b.WriteString(" current row")
	case ast.BoundFollowing:
		writeExpr(b, fb.Offset)
		b.WriteString(" following")
	case ast.BoundUnboundedFollowing:
		b.WriteString(" unbounded  following")
	default:
		b.WriteString("/* unknown bound */")
	}
}

// writeExistsExpr formats an EXISTS expression.
// MySQL 8.0 format: exists(select ...) — column aliases are omitted inside EXISTS.
// writeExistsExpr writes `exists(...)`. MariaDB stores the subquery with
// `limit 1`, which it adds whatever the subquery asked for — `EXISTS(SELECT 1
// FROM t LIMIT 2)` comes back as `exists(select 1 from t limit 1)`, since one
// row settles the predicate. Measured on mariadb:11.8.8.
func writeExistsExpr(b *strings.Builder, n *ast.ExistsExpr) {
	if n.Select != nil {
		b.WriteString("exists(")
		inner := *n.Select
		inner.Limit = &ast.Limit{Count: &ast.IntLit{Value: 1}}
		writeSelectStmt(b, &inner, true)
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

	if q := n.Quantifier.String(); q != "" {
		b.WriteString(q)
		b.WriteString(" ")
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
	if n.Over != nil {
		b.WriteString(" ")
		writeWindowDefAlias(b, n.Over)
	}
}

// writeGroupConcat handles GROUP_CONCAT with its special syntax:
// group_concat([distinct] expr [order by expr ASC|DESC] separator 'str')
// MySQL 8.0 always shows the separator (default ',') and explicit ASC in ORDER BY.
func writeGroupConcat(b *strings.Builder, n *ast.FuncCallExpr) {
	b.WriteString("group_concat(")

	if n.Quantifier.RemovesDuplicates() {
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
	// The call takes a window, like any other aggregate's.
	if n.Over != nil {
		b.WriteString(" ")
		writeWindowDef(b, n.Over)
	}
}

// quoteIdent quotes an identifier in backticks, doubling any backtick inside
// it, as MySQL writes identifiers in SHOW CREATE output. A name can hold one:
// a quoted identifier written with a doubled backtick, or an unaliased item
// named after its source text, such as `a`+1.
func quoteIdent(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// writeCharUsing writes CHAR's USING clause as MariaDB stores it: the charset
// in lower case with utf8 as utf8mb3, and nothing for binary, the default
// (mariadb:11.8.8, SHOW CREATE VIEW: CHAR(65 USING utf8) is char(65 using
// utf8mb3), CHAR(65 USING binary) is char(65)).
func writeCharUsing(b *strings.Builder, cs string) {
	cs = strings.ToLower(cs)
	switch cs {
	case "", "binary":
		return
	case "utf8":
		cs = "utf8mb3"
	}
	b.WriteString(" using ")
	b.WriteString(cs)
}
