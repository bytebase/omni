package deparse

// restore.go writes queries and expressions the way TiDB's parser restores
// them. RestoreFlags picks the context: a view body (the default), or the
// generated column, column default, and partition expressions of table DDL,
// which TiDB restores with other flags (RestoreForGeneratedColumn and its
// siblings). CREATE VIEW stores the body restored with
// RestoreStringSingleQuotes | RestoreKeyWordUppercase | RestoreNameBackQuotes
// (pingcap/tidb pkg/ddl/create_table.go, BuildViewInfo), after the planner
// has given every select field a name and expanded every * (Resolver with
// Restore set does that part). Each writer below follows the Restore method
// of the TiDB node the parser builds for the same text (pingcap/tidb
// pkg/parser/ast, pkg/types/parser_driver/value_expr.go); where TiDB's parser
// rewrites a form -- a + INTERVAL into DATE_ADD, COUNT(*) into COUNT(1), MOD()
// into % -- the writer prints the rewritten form.
//
// Verified against pingcap/tidb:v8.5.5 with SHOW CREATE VIEW: the
// TestDeparse* and TestDeparseContainer* tests in tidb/catalog compare every
// form here with the server's.

import (
	"fmt"
	"strconv"
	"strings"

	ast "github.com/bytebase/omni/tidb/ast"
)

// RestoreFlags selects how Restore writes, as TiDB's format.RestoreFlags
// does. Every form writes names in backquotes and strings in single quotes;
// the zero value is what CREATE VIEW stores: keywords and function names in
// upper case, symbol operators with no spaces, string literals with their
// character set.
type RestoreFlags uint

const (
	// RestoreKeyWordLowercase writes keywords and function names in lower
	// case (format.RestoreKeyWordLowercase).
	RestoreKeyWordLowercase RestoreFlags = 1 << iota
	// RestoreSpacesAroundBinaryOperation writes a space on each side of a
	// symbol operator: a + 1 (format.RestoreSpacesAroundBinaryOperation).
	RestoreSpacesAroundBinaryOperation
	// RestoreBracketAroundBinaryOperation parenthesizes every binary
	// operation: (a+(b*2)) (format.RestoreBracketAroundBinaryOperation).
	RestoreBracketAroundBinaryOperation
	// RestoreWithoutSchemaName drops the schema from column and table
	// names (format.RestoreWithoutSchemaName).
	RestoreWithoutSchemaName
	// RestoreWithoutTableName drops the table from column names
	// (format.RestoreWithoutTableName).
	RestoreWithoutTableName
	// RestoreStringWithoutCharset writes a string literal without its
	// character set introducer (format.RestoreStringWithoutCharset).
	RestoreStringWithoutCharset
)

// The flag sets TiDB's DDL restores each kind of expression with, as v8.5.5
// prints them in SHOW CREATE TABLE.
const (
	// RestoreForView is a view body (pkg/ddl BuildViewInfo).
	RestoreForView RestoreFlags = 0
	// RestoreForGeneratedColumn is a generated column, an expression index
	// part, and a CHECK constraint: `a` + `b` * 2, concat(`c`,
	// _utf8mb4'x'), lower(`c`) (pkg/ddl create_table.go, constraint.go).
	RestoreForGeneratedColumn = RestoreKeyWordLowercase | RestoreSpacesAroundBinaryOperation |
		RestoreWithoutSchemaName | RestoreWithoutTableName
	// RestoreForColumnDefault is an expression default: uuid(),
	// replace(upper(uuid()), _utf8mb4'-', _utf8mb4'') (pkg/ddl
	// add_column.go, restoreFuncCall). TiDB writes some defaults its own
	// way before restoring (NOW() as CURRENT_TIMESTAMP, CURDATE() as
	// CURRENT_DATE); those are the caller's.
	RestoreForColumnDefault = RestoreKeyWordLowercase | RestoreSpacesAroundBinaryOperation
	// RestoreForPartition is a partitioning expression: YEAR(`created`),
	// (`a`+(`b`*2)), ((`a` DIV 2)+(`b`%3)) (pkg/ddl partition.go).
	RestoreForPartition = RestoreBracketAroundBinaryOperation | RestoreWithoutSchemaName | RestoreWithoutTableName
)

// restorer is the text being written and the flags writing it.
type restorer struct {
	strings.Builder
	flags RestoreFlags
}

func (b *restorer) has(f RestoreFlags) bool { return b.flags&f != 0 }

// kw writes a keyword or function name in the case the flags ask for
// (format.RestoreCtx.WriteKeyWord).
func (b *restorer) kw(s string) {
	if b.has(RestoreKeyWordLowercase) {
		b.WriteString(strings.ToLower(s))
		return
	}
	b.WriteString(strings.ToUpper(s))
}

// Restore writes an expression as TiDB restores it in a view definition.
func Restore(node ast.ExprNode) string {
	return RestoreWith(node, RestoreForView)
}

// RestoreWith writes an expression as TiDB restores it with the given flags.
func RestoreWith(node ast.ExprNode, flags RestoreFlags) string {
	b := &restorer{flags: flags}
	restoreExpr(b, node)
	return b.String()
}

// RestoreSelect writes a query as TiDB restores it in a view definition. Its
// select fields carry the names the view's planner gives them only when the
// query went through a Resolver with Restore set.
func RestoreSelect(stmt *ast.SelectStmt) string {
	return RestoreSelectWith(stmt, RestoreForView)
}

// RestoreSelectWith writes a query as TiDB restores it with the given flags.
func RestoreSelectWith(stmt *ast.SelectStmt, flags RestoreFlags) string {
	b := &restorer{flags: flags}
	restoreQuery(b, stmt)
	return b.String()
}

// restoreName writes an identifier in backquotes, doubling a backquote in it
// (format.RestoreCtx.WriteName with RestoreNameBackQuotes).
func restoreName(b *restorer, name string) {
	b.WriteByte('`')
	b.WriteString(strings.ReplaceAll(name, "`", "``"))
	b.WriteByte('`')
}

// restoreString writes a string in single quotes, doubling a quote in it
// (format.RestoreCtx.WriteString with RestoreStringSingleQuotes).
func restoreString(b *restorer, s string) {
	b.WriteByte('\'')
	b.WriteString(strings.ReplaceAll(s, "'", "''"))
	b.WriteByte('\'')
}

// ---------------------------------------------------------------------------
// Queries
// ---------------------------------------------------------------------------

// restoreQuery writes a query expression: a SELECT, a set operation, a
// parenthesized query, or a TABLE or VALUES primary.
func restoreQuery(b *restorer, stmt *ast.SelectStmt) {
	if stmt == nil {
		return
	}
	switch {
	case stmt.SetOp != ast.SetOpNone:
		restoreSetOp(b, stmt)
	case stmt.ParenSource != nil:
		restoreWith(b, stmt.CTEs)
		b.WriteByte('(')
		restoreQuery(b, stmt.ParenSource)
		b.WriteByte(')')
		restoreOrderLimit(b, stmt.OrderBy, stmt.Limit)
	default:
		restoreWith(b, stmt.CTEs)
		restoreSelectLeaf(b, stmt)
	}
}

// restoreWith writes a WITH clause and the space after it
// (ast.WithClause.Restore). TiDB marks the whole clause RECURSIVE.
func restoreWith(b *restorer, ctes []*ast.CommonTableExpr) {
	if len(ctes) == 0 {
		return
	}
	b.kw("WITH ")
	for _, cte := range ctes {
		if cte.Recursive {
			b.kw("RECURSIVE ")
			break
		}
	}
	for i, cte := range ctes {
		if i > 0 {
			b.WriteString(", ")
		}
		restoreName(b, cte.Name)
		if len(cte.Columns) > 0 {
			b.WriteString(" (")
			for j, c := range cte.Columns {
				if j > 0 {
					b.WriteString(", ")
				}
				restoreName(b, c)
			}
			b.WriteByte(')')
		}
		b.kw(" AS (")
		restoreQuery(b, cte.Select)
		b.WriteByte(')')
	}
	b.WriteByte(' ')
}

// restoreSetOp writes a set operation (ast.SetOprStmt.Restore): TiDB keeps a
// chain of operations as one flat list, so the left-deep tree of the omni AST
// is written without parentheses, and an operand is parenthesized only when
// it was written so.
func restoreSetOp(b *restorer, stmt *ast.SelectStmt) {
	restoreWith(b, stmt.CTEs)
	restoreSetOperand(b, stmt.Left, true)
	b.WriteByte(' ')
	b.kw(setOpKeyword(stmt))
	b.WriteByte(' ')
	restoreSetOperand(b, stmt.Right, false)
	restoreOrderLimit(b, stmt.OrderBy, stmt.Limit)
}

func setOpKeyword(stmt *ast.SelectStmt) string {
	var kw string
	switch stmt.SetOp {
	case ast.SetOpUnion:
		kw = "UNION"
	case ast.SetOpIntersect:
		kw = "INTERSECT"
	case ast.SetOpExcept:
		kw = "EXCEPT"
	}
	if stmt.SetAll {
		kw += " ALL"
	}
	return kw
}

// restoreSetOperand writes one operand of a set operation. A set operation
// on the left continues the flat list; one on the right was parenthesized.
func restoreSetOperand(b *restorer, stmt *ast.SelectStmt, left bool) {
	if stmt == nil {
		return
	}
	if stmt.SetOp != ast.SetOpNone && len(stmt.CTEs) == 0 && stmt.OrderBy == nil && stmt.Limit == nil {
		if left {
			restoreSetOperand(b, stmt.Left, true)
			b.WriteByte(' ')
			b.kw(setOpKeyword(stmt))
			b.WriteByte(' ')
			restoreSetOperand(b, stmt.Right, false)
			return
		}
		b.WriteByte('(')
		restoreSetOp(b, stmt)
		b.WriteByte(')')
		return
	}
	restoreQuery(b, stmt)
}

// restoreSelectLeaf writes one SELECT, TABLE, or VALUES query block
// (ast.SelectStmt.Restore).
func restoreSelectLeaf(b *restorer, stmt *ast.SelectStmt) {
	switch {
	case stmt.TableSource != nil:
		b.kw("TABLE ")
		if stmt.TableSource.Table != nil {
			restoreTableRef(b, stmt.TableSource.Table)
		}
	case stmt.ValuesSource != nil:
		b.kw("VALUES ")
		for i, row := range stmt.ValuesSource.Rows {
			if i > 0 {
				b.WriteString(", ")
			}
			b.kw("ROW(")
			for j, v := range row {
				if j > 0 {
					b.WriteByte(',')
				}
				restoreExpr(b, v)
			}
			b.WriteByte(')')
		}
	default:
		restoreSelectBody(b, stmt)
	}
	restoreOrderLimit(b, stmt.OrderBy, stmt.Limit)
	if stmt.ForUpdate != nil {
		// A view's planner drops the locking clause and leaves its lock
		// type at none, which TiDB's Restore writes as the space before it
		// and nothing else: SELECT a FROM t FOR UPDATE is stored as
		// "SELECT `a` AS `a` FROM `test`.`t` " (verified on v8.5.5).
		b.WriteByte(' ')
	}
}

func restoreSelectBody(b *restorer, stmt *ast.SelectStmt) {
	b.kw("SELECT ")
	if stmt.HighPriority {
		b.kw("HIGH_PRIORITY ")
	}
	if stmt.SmallResult {
		b.kw("SQL_SMALL_RESULT ")
	}
	if stmt.BigResult {
		b.kw("SQL_BIG_RESULT ")
	}
	if stmt.BufferResult {
		b.kw("SQL_BUFFER_RESULT ")
	}
	if stmt.NoCache {
		b.kw("SQL_NO_CACHE ")
	}
	if stmt.CalcFoundRows {
		b.kw("SQL_CALC_FOUND_ROWS ")
	}
	switch stmt.DistinctKind {
	case ast.DistinctOn:
		b.kw("DISTINCT ")
	case ast.DistinctAll:
		b.kw("ALL ")
	}
	if stmt.StraightJoin {
		b.kw("STRAIGHT_JOIN ")
	}
	for i, target := range stmt.TargetList {
		if i > 0 {
			b.WriteByte(',')
		}
		restoreField(b, target)
	}
	if len(stmt.From) > 0 {
		b.kw(" FROM ")
		restoreFrom(b, stmt.From)
	} else if stmt.Where != nil {
		b.kw(" FROM DUAL")
	}
	if stmt.Where != nil {
		b.kw(" WHERE ")
		restoreExpr(b, stmt.Where)
	}
	if len(stmt.GroupBy) > 0 {
		b.kw(" GROUP BY ")
		for i, g := range stmt.GroupBy {
			if i > 0 {
				b.WriteByte(',')
			}
			restoreExpr(b, g)
		}
		if stmt.WithRollup {
			b.kw(" WITH ROLLUP")
		}
	}
	if stmt.Having != nil {
		b.kw(" HAVING ")
		restoreExpr(b, stmt.Having)
	}
	if len(stmt.WindowClause) > 0 {
		b.kw(" WINDOW ")
		for i, wd := range stmt.WindowClause {
			if i > 0 {
				b.WriteByte(',')
			}
			restoreName(b, wd.Name)
			b.kw(" AS ")
			restoreWindowSpec(b, wd)
		}
	}
}

// restoreField writes one select field and its name (ast.SelectField.Restore).
func restoreField(b *restorer, target ast.ExprNode) {
	switch t := target.(type) {
	case *ast.ResTarget:
		restoreExpr(b, t.Val)
		if t.Name != "" {
			b.kw(" AS ")
			restoreName(b, t.Name)
		}
	case *ast.StarExpr:
		b.WriteByte('*')
	default:
		restoreExpr(b, target)
	}
}

func restoreOrderLimit(b *restorer, orderBy []*ast.OrderByItem, limit *ast.Limit) {
	if len(orderBy) > 0 {
		b.kw(" ORDER BY ")
		restoreByItems(b, orderBy, ",")
	}
	if limit != nil && limit.Count != nil {
		b.kw(" LIMIT ")
		if limit.Offset != nil {
			restoreExpr(b, limit.Offset)
			b.WriteByte(',')
		}
		restoreExpr(b, limit.Count)
	}
}

func restoreByItems(b *restorer, items []*ast.OrderByItem, sep string) {
	for i, item := range items {
		if i > 0 {
			b.WriteString(sep)
		}
		restoreExpr(b, item.Expr)
		if item.Desc {
			b.kw(" DESC")
		}
	}
}

// ---------------------------------------------------------------------------
// FROM clause
// ---------------------------------------------------------------------------

// joinNode is a FROM item as TiDB's parser builds it: a table, a derived
// table, or a join of two items. A comma list a, b, c is the left-deep
// cross join ((a) JOIN b) JOIN c whose innermost join holds a alone
// (parser.y, TableRefs), which is why TiDB writes the first item of a comma
// list in parentheses.
type joinNode struct {
	leaf        ast.TableExpr   // a table, derived table, or omni join
	left, right *joinNode       // a TiDB join: right is nil for the lone first item
	join        *ast.JoinClause // the omni join a TiDB join stands for; nil for a comma
}

func (n *joinNode) isJoin() bool {
	if n.leaf != nil {
		_, ok := n.leaf.(*ast.JoinClause)
		return ok
	}
	return true
}

func restoreFrom(b *restorer, from []ast.TableExpr) {
	var cur *joinNode
	for i, te := range from {
		item := &joinNode{leaf: te}
		if i == 0 {
			if _, ok := te.(*ast.JoinClause); ok {
				cur = item
			} else {
				cur = &joinNode{left: item}
			}
			continue
		}
		cur = &joinNode{left: cur, right: item}
	}
	restoreJoinNode(b, cur)
}

func restoreJoinNode(b *restorer, n *joinNode) {
	if n.leaf != nil {
		if j, ok := n.leaf.(*ast.JoinClause); ok {
			restoreJoin(b, j)
			return
		}
		restoreTableExpr(b, n.leaf)
		return
	}
	// A comma join, or the lone first item of a comma list.
	useComma := false
	leftIsJoin := n.left.isJoin()
	if n.left.leaf == nil && n.left.right == nil && n.left.left != nil {
		if _, ok := n.left.left.leaf.(*ast.SubqueryExpr); ok {
			useComma = true
		}
	}
	if leftIsJoin && !useComma {
		b.WriteByte('(')
	}
	restoreJoinNode(b, n.left)
	if leftIsJoin && !useComma {
		b.WriteByte(')')
	}
	if n.right == nil {
		return
	}
	if useComma {
		b.WriteString(", ")
	} else {
		b.kw(" JOIN ")
	}
	rightIsJoin := n.right.isJoin()
	if rightIsJoin {
		b.WriteByte('(')
	}
	restoreJoinNode(b, n.right)
	if rightIsJoin {
		b.WriteByte(')')
	}
}

// restoreJoin writes a join (ast.Join.Restore): a join on either side is
// parenthesized, CROSS and INNER joins are both JOIN, and USING and NATURAL
// stay as written.
func restoreJoin(b *restorer, j *ast.JoinClause) {
	_, leftIsJoin := j.Left.(*ast.JoinClause)
	if leftIsJoin {
		b.WriteByte('(')
	}
	restoreTableExpr(b, j.Left)
	if leftIsJoin {
		b.WriteByte(')')
	}
	switch j.Type {
	case ast.JoinNatural:
		b.kw(" NATURAL")
	case ast.JoinNaturalLeft:
		b.kw(" NATURAL LEFT")
	case ast.JoinNaturalRight:
		b.kw(" NATURAL RIGHT")
	case ast.JoinLeft:
		b.kw(" LEFT")
	case ast.JoinRight:
		b.kw(" RIGHT")
	}
	if j.Type == ast.JoinStraight {
		b.kw(" STRAIGHT_JOIN ")
	} else {
		b.kw(" JOIN ")
	}
	_, rightIsJoin := j.Right.(*ast.JoinClause)
	if rightIsJoin {
		b.WriteByte('(')
	}
	restoreTableExpr(b, j.Right)
	if rightIsJoin {
		b.WriteByte(')')
	}
	switch c := j.Condition.(type) {
	case *ast.OnCondition:
		b.kw(" ON ")
		restoreExpr(b, c.Expr)
	case *ast.UsingCondition:
		b.kw(" USING (")
		for i, col := range c.Columns {
			if i > 0 {
				b.WriteByte(',')
			}
			restoreName(b, col)
		}
		b.WriteByte(')')
	}
}

func restoreTableExpr(b *restorer, te ast.TableExpr) {
	switch t := te.(type) {
	case *ast.TableRef:
		restoreTableRef(b, t)
	case *ast.JoinClause:
		restoreJoin(b, t)
	case *ast.SubqueryExpr:
		b.WriteByte('(')
		restoreQuery(b, t.Select)
		b.WriteByte(')')
		if t.Alias != "" {
			b.kw(" AS ")
			restoreName(b, t.Alias)
		}
	default:
		fmt.Fprintf(b, "/* unsupported: %T */", te)
	}
}

// restoreTableRef writes a table (ast.TableSource.Restore over a TableName):
// its schema, which the resolver fills in for every table that is not a CTE,
// partitions, alias, and index hints.
func restoreTableRef(b *restorer, t *ast.TableRef) {
	if t.Schema != "" && !b.has(RestoreWithoutSchemaName) {
		restoreName(b, t.Schema)
		b.WriteByte('.')
	}
	restoreName(b, t.Name)
	if len(t.Partitions) > 0 {
		b.kw(" PARTITION(")
		for i, p := range t.Partitions {
			if i > 0 {
				b.WriteString(", ")
			}
			restoreName(b, p)
		}
		b.WriteByte(')')
	}
	if t.Alias != "" {
		b.kw(" AS ")
		restoreName(b, t.Alias)
	}
	if t.AsOf != nil {
		b.kw(" AS OF TIMESTAMP ")
		restoreExpr(b, t.AsOf.Timestamp)
	}
	for _, h := range t.IndexHints {
		b.WriteByte(' ')
		switch h.Type {
		case ast.HintUse:
			b.kw("USE")
		case ast.HintForce:
			b.kw("FORCE")
		case ast.HintIgnore:
			b.kw("IGNORE")
		}
		b.kw(" INDEX")
		switch h.Scope {
		case ast.HintScopeJoin:
			b.kw(" FOR JOIN")
		case ast.HintScopeOrderBy:
			b.kw(" FOR ORDER BY")
		case ast.HintScopeGroupBy:
			b.kw(" FOR GROUP BY")
		}
		b.WriteString(" (")
		for i, idx := range h.Indexes {
			if i > 0 {
				b.WriteString(", ")
			}
			restoreName(b, idx)
		}
		b.WriteByte(')')
	}
	if t.Sample != nil {
		restoreTableSample(b, t.Sample)
	}
}

// restoreTableSample writes a TABLESAMPLE clause (ast.TableSample.Restore).
// TiDB writes the REGIONS method as REGION, which its own parser refuses, so
// a view it stores that way cannot be read (error 1105 on v8.5.5); this
// writes REGIONS, the word the grammar has, and is otherwise TiDB's text:
// TABLESAMPLE REGIONS () REPEATABLE(1).
func restoreTableSample(b *restorer, s *ast.TableSample) {
	b.kw(" TABLESAMPLE ")
	switch s.Method {
	case ast.SampleMethodSystem:
		b.kw("SYSTEM ")
	case ast.SampleMethodBernoulli:
		b.kw("BERNOULLI ")
	case ast.SampleMethodRegions:
		b.kw("REGIONS ")
	}
	b.WriteByte('(')
	if s.Size != nil {
		restoreExpr(b, s.Size)
	}
	switch s.Unit {
	case ast.SampleUnitPercent:
		b.kw(" PERCENT")
	case ast.SampleUnitRows:
		b.kw(" ROWS")
	}
	b.WriteByte(')')
	if s.Repeatable != nil {
		b.kw(" REPEATABLE(")
		restoreExpr(b, s.Repeatable)
		b.WriteByte(')')
	}
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------

func restoreExpr(b *restorer, node ast.ExprNode) {
	switch n := node.(type) {
	case nil:
	case *ast.IntLit:
		b.WriteString(strconv.FormatInt(n.Value, 10))
	case *ast.FloatLit:
		b.WriteString(restoreFloat(n.Value))
	case *ast.StringLit:
		restoreStringLit(b, n)
	case *ast.BoolLit:
		if n.Value {
			b.kw("TRUE")
		} else {
			b.kw("FALSE")
		}
	case *ast.NullLit:
		b.kw("NULL")
	case *ast.HexLit:
		b.WriteString(restoreHex(n.Value))
	case *ast.BitLit:
		b.WriteString(restoreBits(n.Value))
	case *ast.TemporalLit:
		// ast.FuncCallExpr.customRestore: the keyword, then the string as
		// the lexer read it, with no character set.
		b.WriteString(strings.ToUpper(n.Type))
		b.WriteByte(' ')
		restoreString(b, n.Value)
	case *ast.ColumnRef:
		restoreColumnRef(b, n)
	case *ast.StarExpr:
		b.WriteByte('*')
	case *ast.ParenExpr:
		b.WriteByte('(')
		restoreExpr(b, n.Expr)
		b.WriteByte(')')
	case *ast.BinaryExpr:
		restoreBinary(b, n)
	case *ast.UnaryExpr:
		restoreUnary(b, n)
	case *ast.InExpr:
		restoreExpr(b, n.Expr)
		if n.Not {
			b.kw(" NOT IN ")
		} else {
			b.kw(" IN ")
		}
		if n.Select != nil {
			b.WriteByte('(')
			restoreQuery(b, n.Select)
			b.WriteByte(')')
		} else {
			b.WriteByte('(')
			for i, item := range n.List {
				if i > 0 {
					b.WriteByte(',')
				}
				restoreExpr(b, item)
			}
			b.WriteByte(')')
		}
	case *ast.BetweenExpr:
		restoreExpr(b, n.Expr)
		if n.Not {
			b.kw(" NOT BETWEEN ")
		} else {
			b.kw(" BETWEEN ")
		}
		restoreExpr(b, n.Low)
		b.kw(" AND ")
		restoreExpr(b, n.High)
	case *ast.LikeExpr:
		restoreExpr(b, n.Expr)
		// ast.PatternLikeOrIlikeExpr.Restore: ILIKE as written.
		op := " LIKE "
		if n.ILike {
			op = " ILIKE "
		}
		if n.Not {
			b.kw(" NOT" + op)
		} else {
			b.kw(op)
		}
		restoreExpr(b, n.Pattern)
		if n.Escape != nil {
			// The escape is kept as a byte and written back as a plain
			// string (ast.PatternLikeOrIlikeExpr.Restore); the default
			// backslash is not written.
			if s, ok := n.Escape.(*ast.StringLit); ok {
				if s.Value != `\` {
					b.kw(" ESCAPE ")
					restoreString(b, s.Value)
				}
			} else {
				b.kw(" ESCAPE ")
				restoreExpr(b, n.Escape)
			}
		}
	case *ast.IsExpr:
		restoreExpr(b, n.Expr)
		b.kw(" IS")
		if n.Not {
			b.kw(" NOT")
		}
		switch n.Test {
		case ast.IsNull, ast.IsUnknown:
			// IS UNKNOWN parses as IS NULL.
			b.kw(" NULL")
		case ast.IsTrue:
			b.kw(" TRUE")
		case ast.IsFalse:
			b.kw(" FALSE")
		}
	case *ast.RowExpr:
		b.kw("ROW(")
		for i, item := range n.Items {
			if i > 0 {
				b.WriteByte(',')
			}
			restoreExpr(b, item)
		}
		b.WriteByte(')')
	case *ast.CaseExpr:
		b.kw("CASE")
		if n.Operand != nil {
			b.WriteByte(' ')
			restoreExpr(b, n.Operand)
		}
		for _, w := range n.Whens {
			b.kw(" WHEN ")
			restoreExpr(b, w.Cond)
			b.kw(" THEN ")
			restoreExpr(b, w.Result)
		}
		if n.Default != nil {
			b.kw(" ELSE ")
			restoreExpr(b, n.Default)
		}
		b.kw(" END")
	case *ast.CastExpr:
		b.kw("CAST(")
		restoreExpr(b, n.Expr)
		b.kw(" AS ")
		restoreCastType(b, n.TypeName)
		if n.Array {
			b.kw(" ARRAY")
		}
		b.WriteByte(')')
	case *ast.ConvertExpr:
		b.kw("CONVERT(")
		restoreExpr(b, n.Expr)
		if n.TypeName != nil {
			b.WriteString(", ")
			restoreCastType(b, n.TypeName)
		} else {
			// The character set is a string argument with no introducer
			// (ast.FuncCallExpr.Restore, case "convert"), named as the
			// charset registry names it: USING UTF8MB4 is 'utf8mb4'.
			b.kw(" USING ")
			restoreString(b, charsetName(n.Charset))
		}
		b.WriteByte(')')
	case *ast.CollateExpr:
		restoreExpr(b, n.Expr)
		b.kw(" COLLATE ")
		b.WriteString(n.Collation)
	case *ast.IntervalExpr:
		restoreInterval(b, n)
	case *ast.ExtractExpr:
		b.kw("EXTRACT(")
		b.kw(n.Unit)
		b.kw(" FROM ")
		restoreExpr(b, n.Expr)
		b.WriteByte(')')
	case *ast.PositionExpr:
		b.kw("POSITION(")
		restoreExpr(b, n.Substr)
		b.kw(" IN ")
		restoreExpr(b, n.Str)
		b.WriteByte(')')
	case *ast.KeywordArg:
		b.kw(n.Keyword)
	case *ast.FuncCallExpr:
		restoreFuncCall(b, n)
	case *ast.ExistsExpr:
		b.kw("EXISTS (")
		restoreQuery(b, n.Select)
		b.WriteByte(')')
	case *ast.SubqueryExpr:
		if n.Exists {
			b.kw("EXISTS ")
		}
		b.WriteByte('(')
		restoreQuery(b, n.Select)
		b.WriteByte(')')
	case *ast.MemberOfExpr:
		restoreExpr(b, n.Value)
		b.kw(" MEMBER OF (")
		restoreExpr(b, n.Array)
		b.WriteByte(')')
	case *ast.DefaultExpr:
		b.kw("DEFAULT")
		if n.Column != "" {
			b.WriteByte('(')
			restoreName(b, n.Column)
			b.WriteByte(')')
		}
	case *ast.MatchExpr:
		b.kw("MATCH (")
		for i, c := range n.Columns {
			if i > 0 {
				b.WriteByte(',')
			}
			restoreColumnRef(b, c)
		}
		b.kw(") AGAINST (")
		restoreExpr(b, n.Against)
		if m := strings.ToUpper(strings.TrimSpace(n.Modifier)); m == "IN BOOLEAN MODE" || m == "WITH QUERY EXPANSION" {
			b.WriteByte(' ')
			b.WriteString(m)
		}
		b.WriteByte(')')
	case *ast.VariableRef:
		restoreVariable(b, n)
	case *ast.ResTarget:
		restoreExpr(b, n.Val)
	default:
		fmt.Fprintf(b, "/* unsupported: %T */", node)
	}
}

// restoreVariable writes a user or system variable (ast.VariableExpr.Restore):
// @`a`, @@`var`, @@GLOBAL.`b`, @@SESSION.`c`. TiDB's parser lowers a system
// variable's name and reads @@LOCAL as SESSION; a user variable keeps its case.
func restoreVariable(b *restorer, n *ast.VariableRef) {
	if !n.System {
		b.WriteByte('@')
		restoreName(b, n.Name)
		return
	}
	b.WriteString("@@")
	switch strings.ToUpper(n.Scope) {
	case "GLOBAL":
		b.kw("GLOBAL")
		b.WriteByte('.')
	case "SESSION", "LOCAL":
		b.kw("SESSION")
		b.WriteByte('.')
	}
	restoreName(b, strings.ToLower(n.Name))
}

func restoreColumnRef(b *restorer, n *ast.ColumnRef) {
	if n.Schema != "" && !b.has(RestoreWithoutSchemaName) {
		restoreName(b, n.Schema)
		b.WriteByte('.')
	}
	if n.Table != "" && !b.has(RestoreWithoutTableName) {
		restoreName(b, n.Table)
		b.WriteByte('.')
	}
	if n.Star {
		b.WriteByte('*')
		return
	}
	restoreName(b, n.Column)
}

// charsetName is a character set as TiDB's CharsetName action names it, the
// registry's name (charset.GetCharsetInfo): lower case, with utf8mb3 as its
// alias utf8. CHAR(65 USING UTF8MB3) is stored as CHAR_FUNC(65, 'utf8')
// (verified on v8.5.5).
func charsetName(cs string) string {
	cs = strings.ToLower(cs)
	if cs == "utf8mb3" {
		return "utf8"
	}
	return cs
}

// restoreStringLit writes a string literal with the character set TiDB gives
// it: the introducer written, or the connection's, utf8mb4
// (parser_driver.ValueExpr.Restore). A backslash is doubled.
func restoreStringLit(b *restorer, n *ast.StringLit) {
	// N'...' is read as _utf8'...', as TiDB's lexer reads it.
	cs := strings.TrimPrefix(n.Charset, "_")
	if cs == "" {
		cs = "utf8mb4"
	}
	if !b.has(RestoreStringWithoutCharset) {
		b.WriteByte('_')
		b.kw(cs)
	}
	restoreString(b, strings.ReplaceAll(n.Value, `\`, `\\`))
}

// restoreFloat writes a number with a fraction or an exponent. With an
// exponent it is a double, which TiDB writes in Go's shortest e form (1e3
// is 1e+03); without one it is a decimal, which loses its leading zeros
// and gains a 0 before a bare point (.5 is 0.5).
func restoreFloat(s string) string {
	if strings.ContainsAny(s, "eE") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return strconv.FormatFloat(f, 'e', -1, 64)
		}
		return s
	}
	intPart, frac, hasPoint := strings.Cut(s, ".")
	intPart = strings.TrimLeft(intPart, "0")
	if intPart == "" {
		intPart = "0"
	}
	if !hasPoint || frac == "" {
		return intPart
	}
	return intPart + "." + frac
}

// restoreHex writes a hexadecimal literal as TiDB does: x'..' in lower case,
// padded to whole bytes.
func restoreHex(v string) string {
	h := v
	switch {
	case strings.HasPrefix(h, "0x"), strings.HasPrefix(h, "0X"):
		h = h[2:]
	case len(h) >= 3 && (h[0] == 'x' || h[0] == 'X') && h[1] == '\'':
		h = strings.TrimSuffix(h[2:], "'")
	}
	h = strings.ToLower(h)
	if len(h)%2 == 1 {
		h = "0" + h
	}
	return "x'" + h + "'"
}

// restoreBits writes a bit literal as TiDB does: b'..' without leading zeros
// (types.BinaryLiteral.ToBitLiteralString).
func restoreBits(v string) string {
	bits := v
	switch {
	case strings.HasPrefix(bits, "0b"), strings.HasPrefix(bits, "0B"):
		bits = bits[2:]
	case len(bits) >= 3 && (bits[0] == 'b' || bits[0] == 'B') && bits[1] == '\'':
		bits = strings.TrimSuffix(bits[2:], "'")
	}
	bits = strings.TrimLeft(bits, "0")
	if bits == "" {
		return "b''"
	}
	return "b'" + bits + "'"
}

// binaryOps maps an operator to its TiDB spelling (opcode.Ops); a keyword
// operator is written with a space on each side, a symbol with none.
var binaryOps = map[ast.BinaryOp]struct {
	text    string
	keyword bool
}{
	ast.BinOpAdd:        {"+", false},
	ast.BinOpSub:        {"-", false},
	ast.BinOpMul:        {"*", false},
	ast.BinOpDiv:        {"/", false},
	ast.BinOpMod:        {"%", false},
	ast.BinOpEq:         {"=", false},
	ast.BinOpNe:         {"!=", false},
	ast.BinOpLt:         {"<", false},
	ast.BinOpGt:         {">", false},
	ast.BinOpLe:         {"<=", false},
	ast.BinOpGe:         {">=", false},
	ast.BinOpAnd:        {"AND", true},
	ast.BinOpOr:         {"OR", true},
	ast.BinOpXor:        {"XOR", true},
	ast.BinOpBitAnd:     {"&", false},
	ast.BinOpBitOr:      {"|", false},
	ast.BinOpBitXor:     {"^", false},
	ast.BinOpShiftLeft:  {"<<", false},
	ast.BinOpShiftRight: {">>", false},
	ast.BinOpDivInt:     {"DIV", true},
	ast.BinOpRegexp:     {"REGEXP", true},
	ast.BinOpNullSafeEq: {"<=>", false},
	ast.BinOpAssign:     {":=", false},
}

func writeBinaryOp(b *restorer, op ast.BinaryOp) {
	o, ok := binaryOps[op]
	if !ok {
		fmt.Fprintf(b, " /* unsupported op %d */ ", op)
		return
	}
	if o.keyword {
		b.WriteByte(' ')
		b.kw(o.text)
		b.WriteByte(' ')
		return
	}
	if b.has(RestoreSpacesAroundBinaryOperation) {
		b.WriteByte(' ')
		b.WriteString(o.text)
		b.WriteByte(' ')
		return
	}
	b.WriteString(o.text)
}

func restoreBinary(b *restorer, n *ast.BinaryExpr) {
	switch n.Op {
	case ast.BinOpJsonExtract:
		// col->path is JSON_EXTRACT(col, path) once parsed.
		b.kw("JSON_EXTRACT(")
		restoreExpr(b, n.Left)
		b.WriteString(", ")
		restoreExpr(b, n.Right)
		b.WriteByte(')')
		return
	case ast.BinOpJsonUnquote:
		b.kw("JSON_UNQUOTE(JSON_EXTRACT(")
		restoreExpr(b, n.Left)
		b.WriteString(", ")
		restoreExpr(b, n.Right)
		b.WriteString("))")
		return
	case ast.BinOpAdd, ast.BinOpSub:
		// expr + INTERVAL n unit is DATE_ADD(expr, INTERVAL n unit), and
		// INTERVAL n unit + expr the same; expr - INTERVAL is DATE_SUB.
		if iv, ok := n.Right.(*ast.IntervalExpr); ok {
			if n.Op == ast.BinOpAdd {
				b.kw("DATE_ADD(")
			} else {
				b.kw("DATE_SUB(")
			}
			restoreExpr(b, n.Left)
			b.WriteString(", ")
			restoreInterval(b, iv)
			b.WriteByte(')')
			return
		}
		if iv, ok := n.Left.(*ast.IntervalExpr); ok && n.Op == ast.BinOpAdd {
			b.kw("DATE_ADD(")
			restoreExpr(b, n.Right)
			b.WriteString(", ")
			restoreInterval(b, iv)
			b.WriteByte(')')
			return
		}
	}
	if n.Op == ast.BinOpRegexp {
		// A PatternRegexpExpr in TiDB, never bracketed.
		restoreExpr(b, n.Left)
		writeBinaryOp(b, n.Op)
		restoreExpr(b, n.Right)
		return
	}
	if sub, ok := n.Right.(*ast.SubqueryExpr); ok && sub.Quantifier != "" {
		restoreExpr(b, n.Left)
		writeBinaryOp(b, n.Op)
		// ast.CompareSubqueryExpr.Restore: SOME is ANY.
		if strings.EqualFold(sub.Quantifier, "ALL") {
			b.kw("ALL ")
		} else {
			b.kw("ANY ")
		}
		b.WriteByte('(')
		restoreQuery(b, sub.Select)
		b.WriteByte(')')
		return
	}
	// ast.BinaryOperationExpr.Restore.
	bracket := b.has(RestoreBracketAroundBinaryOperation)
	if bracket {
		b.WriteByte('(')
	}
	restoreExpr(b, n.Left)
	writeBinaryOp(b, n.Op)
	restoreExpr(b, n.Right)
	if bracket {
		b.WriteByte(')')
	}
}

func restoreUnary(b *restorer, n *ast.UnaryExpr) {
	switch n.Op {
	case ast.UnaryMinus:
		b.WriteByte('-')
	case ast.UnaryPlus:
		b.WriteByte('+')
	case ast.UnaryBitNot:
		b.WriteByte('~')
	case ast.UnaryBinary:
		b.kw("BINARY ")
	case ast.UnaryNot:
		if n.OriginalOp == "!" {
			b.WriteByte('!')
		} else {
			// NOT REGEXP, NOT EXISTS, and NOT LIKE fold into the operator.
			switch inner := n.Operand.(type) {
			case *ast.BinaryExpr:
				if inner.Op == ast.BinOpRegexp {
					restoreExpr(b, inner.Left)
					b.kw(" NOT REGEXP ")
					restoreExpr(b, inner.Right)
					return
				}
			case *ast.ExistsExpr:
				b.kw("NOT EXISTS (")
				restoreQuery(b, inner.Select)
				b.WriteByte(')')
				return
			}
			b.kw("NOT ")
		}
	}
	restoreExpr(b, n.Operand)
}

func restoreInterval(b *restorer, n *ast.IntervalExpr) {
	b.kw("INTERVAL ")
	restoreExpr(b, n.Value)
	b.WriteByte(' ')
	b.kw(n.Unit)
}

// restoreCastType writes the target type of CAST or CONVERT
// (types.FieldType.RestoreAsCastType).
func restoreCastType(b *restorer, dt *ast.DataType) {
	if dt == nil {
		return
	}
	name := strings.ToUpper(dt.Name)
	switch name {
	case "CHAR", "NCHAR", "VARCHAR", "BINARY":
		binary := name == "BINARY" || strings.EqualFold(dt.Charset, "binary")
		if binary {
			b.kw("BINARY")
		} else {
			b.kw("CHAR")
		}
		if dt.Length > 0 {
			fmt.Fprintf(b, "(%d)", dt.Length)
		}
		if binary {
			return
		}
		if dt.Binary {
			b.kw(" BINARY")
		}
		if cs := strings.ToLower(dt.Charset); cs != "" && cs != "utf8mb4" {
			b.kw(" CHARSET ")
			b.kw(cs)
		}
	case "DATE", "JSON", "DOUBLE", "FLOAT", "YEAR", "VECTOR":
		b.kw(name)
	case "REAL":
		b.kw("DOUBLE")
	case "DATETIME", "TIME":
		b.kw(name)
		if dt.Length > 0 {
			fmt.Fprintf(b, "(%d)", dt.Length)
		}
	case "DECIMAL", "DEC", "NUMERIC", "FIXED":
		b.kw("DECIMAL")
		switch {
		case dt.Length > 0 && dt.Scale > 0:
			fmt.Fprintf(b, "(%d, %d)", dt.Length, dt.Scale)
		case dt.Length > 0:
			fmt.Fprintf(b, "(%d)", dt.Length)
		default:
			// DECIMAL with no precision is DECIMAL(10).
			b.WriteString("(10)")
		}
	case "SIGNED", "SIGNED INTEGER", "INT", "INTEGER":
		if dt.Unsigned {
			b.kw("UNSIGNED")
		} else {
			b.kw("SIGNED")
		}
	case "UNSIGNED", "UNSIGNED INTEGER":
		b.kw("UNSIGNED")
	default:
		b.kw(name)
	}
}

// ---------------------------------------------------------------------------
// Function calls
// ---------------------------------------------------------------------------

// funcNames maps a function to the name TiDB's parser gives it, where that
// differs from the name written.
var funcNames = map[string]string{
	"STD":      "STDDEV_POP",
	"STDDEV":   "STDDEV_POP",
	"VARIANCE": "VAR_POP",
	"CHAR":     "CHAR_FUNC",
}

func restoreFuncCall(b *restorer, n *ast.FuncCallExpr) {
	name := strings.ToUpper(n.Name)
	switch name {
	case "MOD":
		if len(n.Args) == 2 && n.Over == nil {
			// MOD(a, b) parses as a % b, and TiDB writes it with no
			// parentheses, which regroups it: 2*MOD(5,3) is stored as
			// 2*5%3, and the view then answers 1, not 4 (verified on
			// v8.5.5). The stored text is what TiDB runs, so it is kept.
			restoreBinary(b, &ast.BinaryExpr{Op: ast.BinOpMod, Left: n.Args[0], Right: n.Args[1]})
			return
		}
	case "TRIM", "TRIM_LEADING", "TRIM_TRAILING", "TRIM_BOTH":
		restoreTrim(b, name, n.Args)
		return
	case "ADDDATE", "SUBDATE", "DATE_ADD", "DATE_SUB":
		if len(n.Args) == 2 {
			restoreFuncName(b, n, name)
			b.WriteByte('(')
			restoreExpr(b, n.Args[0])
			b.WriteString(", ")
			if iv, ok := n.Args[1].(*ast.IntervalExpr); ok {
				restoreInterval(b, iv)
			} else {
				// ADDDATE(d, n) is ADDDATE(d, INTERVAL n DAY).
				b.kw("INTERVAL ")
				restoreExpr(b, n.Args[1])
				b.kw(" DAY")
			}
			b.WriteByte(')')
			return
		}
	case "CHAR":
		// CHAR(n, ... [USING cs]) is CHAR_FUNC(n, ..., charset)
		// (parser.y:8184-8201): the charset is NULL when no USING clause is
		// written, and otherwise a string with no introducer, named as the
		// charset registry names it: USING UTF8MB4 is 'utf8mb4'.
		b.kw("CHAR_FUNC(")
		for _, arg := range n.Args {
			restoreExpr(b, arg)
			b.WriteString(", ")
		}
		if n.Using == "" {
			b.kw("NULL")
		} else {
			restoreString(b, charsetName(n.Using))
		}
		b.WriteByte(')')
		return
	case "GROUP_CONCAT":
		restoreGroupConcat(b, n)
		return
	}

	restoreFuncName(b, n, name)
	b.WriteByte('(')
	if n.Quantifier.RemovesDuplicates() {
		b.kw("DISTINCT ")
	}
	if n.Star {
		// COUNT(*) parses as COUNT(1).
		b.WriteByte('1')
	}
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		restoreExpr(b, arg)
	}
	b.WriteByte(')')
	if n.Over != nil {
		b.kw(" OVER ")
		restoreWindowSpec(b, n.Over)
	}
}

func restoreFuncName(b *restorer, n *ast.FuncCallExpr, name string) {
	if n.Schema != "" {
		restoreName(b, n.Schema)
		b.WriteByte('.')
		restoreName(b, n.Name)
		return
	}
	if mapped, ok := funcNames[name]; ok {
		name = mapped
	}
	b.kw(name)
}

// restoreTrim writes TRIM (ast.FuncCallExpr.Restore, case "trim"): the
// direction, then the string to remove, then FROM and the string.
func restoreTrim(b *restorer, name string, args []ast.ExprNode) {
	b.kw("TRIM(")
	direction := strings.TrimPrefix(name, "TRIM")
	direction = strings.TrimPrefix(direction, "_")
	switch {
	case len(args) == 1 && direction == "":
		restoreExpr(b, args[0])
	case len(args) == 1:
		// TRIM(LEADING FROM s) removes spaces.
		b.kw(direction)
		b.WriteByte(' ')
		restoreStringLit(b, &ast.StringLit{Value: " "})
		b.kw(" FROM ")
		restoreExpr(b, args[0])
	case len(args) >= 2:
		if direction != "" {
			b.kw(direction)
			b.WriteByte(' ')
		}
		restoreExpr(b, args[0])
		b.kw(" FROM ")
		restoreExpr(b, args[1])
	}
	b.WriteByte(')')
}

// restoreGroupConcat writes GROUP_CONCAT with its separator, which TiDB
// always writes, as a string with no introducer and its backslashes doubled
// (ast.AggregateFuncExpr.Restore over a ValueExpr): SEPARATOR '\\' is stored
// as SEPARATOR '\\' (verified on v8.5.5).
func restoreGroupConcat(b *restorer, n *ast.FuncCallExpr) {
	b.kw("GROUP_CONCAT(")
	if n.Quantifier.RemovesDuplicates() {
		b.kw("DISTINCT ")
	}
	for i, arg := range n.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		restoreExpr(b, arg)
	}
	if len(n.OrderBy) > 0 {
		b.kw(" ORDER BY ")
		restoreByItems(b, n.OrderBy, ",")
	}
	b.kw(" SEPARATOR ")
	sep := ","
	if s, ok := n.Separator.(*ast.StringLit); ok {
		sep = s.Value
	}
	restoreString(b, strings.ReplaceAll(sep, `\`, `\\`))
	b.WriteByte(')')
	if n.Over != nil {
		b.kw(" OVER ")
		restoreWindowSpec(b, n.Over)
	}
}

// restoreWindowSpec writes a window (ast.WindowSpec.Restore): a bare name
// for OVER w, and otherwise the parenthesized specification.
func restoreWindowSpec(b *restorer, wd *ast.WindowDef) {
	if wd.RefName != "" && len(wd.PartitionBy) == 0 && len(wd.OrderBy) == 0 && wd.Frame == nil && wd.Name == "" {
		restoreName(b, wd.RefName)
		return
	}
	b.WriteByte('(')
	sep := ""
	if wd.RefName != "" {
		restoreName(b, wd.RefName)
		sep = " "
	}
	if len(wd.PartitionBy) > 0 {
		b.WriteString(sep)
		b.kw("PARTITION BY ")
		for i, e := range wd.PartitionBy {
			if i > 0 {
				b.WriteString(", ")
			}
			restoreExpr(b, e)
		}
		sep = " "
	}
	if len(wd.OrderBy) > 0 {
		b.WriteString(sep)
		b.kw("ORDER BY ")
		restoreByItems(b, wd.OrderBy, ",")
		sep = " "
	}
	if f := wd.Frame; f != nil {
		b.WriteString(sep)
		switch f.Type {
		case ast.FrameRange:
			b.kw("RANGE")
		case ast.FrameGroups:
			b.kw("GROUPS")
		default:
			b.kw("ROWS")
		}
		b.kw(" BETWEEN ")
		restoreFrameBound(b, f.Start)
		b.kw(" AND ")
		if f.End != nil {
			restoreFrameBound(b, f.End)
		} else {
			// A frame with one bound ends at the current row.
			b.kw("CURRENT ROW")
		}
	}
	b.WriteByte(')')
}

func restoreFrameBound(b *restorer, fb *ast.WindowFrameBound) {
	if fb == nil {
		return
	}
	switch fb.Type {
	case ast.BoundUnboundedPreceding:
		b.kw("UNBOUNDED PRECEDING")
	case ast.BoundUnboundedFollowing:
		b.kw("UNBOUNDED FOLLOWING")
	case ast.BoundCurrentRow:
		b.kw("CURRENT ROW")
	case ast.BoundPreceding, ast.BoundFollowing:
		restoreExpr(b, fb.Offset)
		if fb.Type == ast.BoundPreceding {
			b.kw(" PRECEDING")
		} else {
			b.kw(" FOLLOWING")
		}
	}
}
