package ast

// This file holds AST node types for SELECT statement parsing (T1.4).

// ---------------------------------------------------------------------------
// SELECT statement
// ---------------------------------------------------------------------------

// WithClause represents a WITH clause preceding a SELECT statement.
//
//	WITH [RECURSIVE] cte_name [(col1, col2, ...)] AS (query)
//	  [, cte_name2 AS (query2)]
type WithClause struct {
	Recursive bool   // RECURSIVE keyword present
	CTEs      []*CTE // one or more CTE definitions
	Loc       Loc
}

// Tag implements Node.
func (n *WithClause) Tag() NodeTag { return T_WithClause }

// Compile-time assertion that *WithClause satisfies Node.
var _ Node = (*WithClause)(nil)

// CTE is one named entry in a WITH clause:
//
//	cte_name [(col1, col2, ...)] AS (query)
type CTE struct {
	Name    string   // the CTE alias name
	Columns []string // optional column aliases; nil/empty if not specified
	Query   Node     // the inner SELECT statement (*SelectStmt)
	Loc     Loc
}

// Tag implements Node.
func (n *CTE) Tag() NodeTag { return T_CTE }

// Compile-time assertion that *CTE satisfies Node.
var _ Node = (*CTE)(nil)

// SelectStmt represents a full SELECT statement.
//
//	[WITH [RECURSIVE] cte ...]
//	SELECT [DISTINCT|ALL] select_list
//	  [FROM table_references]
//	  [WHERE condition]
//	  [GROUP BY expr, ...]
//	  [HAVING condition]
//	  [QUALIFY condition]
//	  [ORDER BY expr [ASC|DESC] [NULLS FIRST|LAST], ...]
//	  [LIMIT count [OFFSET offset]]
type SelectStmt struct {
	With     *WithClause   // optional WITH clause (nil if absent)
	Distinct bool          // DISTINCT keyword present
	All      bool          // ALL keyword present (explicit, rarely used)
	Items    []*SelectItem // SELECT list
	From     []Node        // FROM clause table references (TableRef or JoinClause)
	Where    Node          // WHERE expression (nil if absent)
	GroupBy  []Node        // GROUP BY expressions
	Having   Node          // HAVING expression (nil if absent)
	Qualify  Node          // QUALIFY expression (nil if absent)
	OrderBy  []*OrderByItem // ORDER BY items
	Limit    Node          // LIMIT expression (nil if absent)
	Offset   Node          // OFFSET expression (nil if absent)
	Loc      Loc
}

// Tag implements Node.
func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }

// Compile-time assertion that *SelectStmt satisfies Node.
var _ Node = (*SelectStmt)(nil)

// ---------------------------------------------------------------------------
// SELECT list item
// ---------------------------------------------------------------------------

// SelectItem represents one item in the SELECT list:
//   - expr [AS alias]    — a computed expression with optional alias
//   - *                  — all columns
//   - table.*            — all columns from a specific table
type SelectItem struct {
	Expr  Node   // the expression; for *, this is nil
	Alias string // optional alias name; empty if absent
	Star  bool   // true for * or table.*
	// For table.*, TableName holds the qualifier ObjectName.
	TableName *ObjectName
	Loc       Loc
}

// Tag implements Node.
func (n *SelectItem) Tag() NodeTag { return T_SelectItem }

// Compile-time assertion that *SelectItem satisfies Node.
var _ Node = (*SelectItem)(nil)

// ---------------------------------------------------------------------------
// Table reference
// ---------------------------------------------------------------------------

// TableRef represents a table reference in the FROM clause.
//
//	table_name [AS alias]
//	db.table_name [AS alias]
//	catalog.db.table_name [AS alias]
type TableRef struct {
	Name  *ObjectName // table name (may be qualified: db.table or catalog.db.table)
	Alias string      // optional alias; empty if absent
	Loc   Loc
}

// Tag implements Node.
func (n *TableRef) Tag() NodeTag { return T_TableRef }

// Compile-time assertion that *TableRef satisfies Node.
var _ Node = (*TableRef)(nil)

// ---------------------------------------------------------------------------
// JOIN clause (basic structure for T1.4; T1.5 will flesh out)
// ---------------------------------------------------------------------------

// JoinType classifies the kind of join.
type JoinType int

const (
	JoinInner     JoinType = iota // [INNER] JOIN
	JoinLeft                      // LEFT [OUTER] JOIN
	JoinRight                     // RIGHT [OUTER] JOIN
	JoinFull                      // FULL [OUTER] JOIN
	JoinCross                     // CROSS JOIN
	JoinLeftSemi                  // LEFT SEMI JOIN
	JoinRightSemi                 // RIGHT SEMI JOIN
	JoinLeftAnti                  // LEFT ANTI JOIN
	JoinRightAnti                 // RIGHT ANTI JOIN
)

// JoinClause represents a JOIN expression in the FROM clause.
type JoinClause struct {
	Type    JoinType
	Left    Node     // left side of the join
	Right   Node     // right side of the join
	Natural bool     // NATURAL join modifier
	On      Node     // ON condition (nil if absent)
	Using   []string // USING (col1, col2, ...) column names
	Hints   []string // execution hints, e.g. [shuffle], [broadcast]
	Loc     Loc
}

// Tag implements Node.
func (n *JoinClause) Tag() NodeTag { return T_JoinClause }

// Compile-time assertion that *JoinClause satisfies Node.
var _ Node = (*JoinClause)(nil)
