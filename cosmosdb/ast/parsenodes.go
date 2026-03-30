package ast

// ---------------------------------------------------------------------------
// Wrapper / top-level
// ---------------------------------------------------------------------------

// RawStmt wraps a parsed statement with position info.
type RawStmt struct {
	Stmt         StmtNode
	StmtLocation int // start byte offset
	StmtLen      int // length in bytes; 0 means "rest of string"
}

func (*RawStmt) nodeTag() {}

// ---------------------------------------------------------------------------
// Statement nodes
// ---------------------------------------------------------------------------

// SelectStmt represents a complete SELECT query.
type SelectStmt struct {
	Top         *int               // TOP n (nil if not specified)
	Distinct    bool               // DISTINCT modifier
	Value       bool               // VALUE modifier
	Star        bool               // SELECT *
	Targets     []*TargetEntry     // select list (nil when Star is true)
	From        TableExpr          // FROM clause (nil if omitted)
	Joins       []*JoinExpr        // JOIN clauses
	Where       ExprNode           // WHERE clause (nil if omitted)
	GroupBy     []ExprNode         // GROUP BY expressions
	Having      ExprNode           // HAVING clause (nil if omitted)
	OrderBy     []*SortExpr        // ORDER BY expressions
	OffsetLimit *OffsetLimitClause // OFFSET ... LIMIT ... (nil if omitted)
	Loc         Loc
}

func (*SelectStmt) nodeTag()  {}
func (*SelectStmt) stmtNode() {}

// TargetEntry represents one item in a SELECT list: expr [AS alias].
type TargetEntry struct {
	Expr  ExprNode
	Alias *string // nil if no alias
	Loc   Loc
}

func (*TargetEntry) nodeTag() {}

// SortExpr represents one item in an ORDER BY list.
type SortExpr struct {
	Expr ExprNode
	Desc bool // true for DESC, false for ASC (default)
	Rank bool // true for RANK expr
	Loc  Loc
}

func (*SortExpr) nodeTag() {}

// OffsetLimitClause represents OFFSET n LIMIT m.
type OffsetLimitClause struct {
	Offset ExprNode
	Limit  ExprNode
	Loc    Loc
}

func (*OffsetLimitClause) nodeTag() {}

// JoinExpr represents a JOIN clause: JOIN from_source.
type JoinExpr struct {
	Source TableExpr
	Loc    Loc
}

func (*JoinExpr) nodeTag() {}

// ---------------------------------------------------------------------------
// Table expression nodes — implement TableExpr
// ---------------------------------------------------------------------------

// ContainerRef references a container by name, or ROOT.
type ContainerRef struct {
	Name string // empty when Root is true
	Root bool
	Loc  Loc
}

func (*ContainerRef) nodeTag()   {}
func (*ContainerRef) tableExpr() {}

// AliasedTableExpr wraps a table expression with an alias: source [AS] alias.
type AliasedTableExpr struct {
	Source TableExpr
	Alias  string
	Loc    Loc
}

func (*AliasedTableExpr) nodeTag()   {}
func (*AliasedTableExpr) tableExpr() {}

// ArrayIterationExpr represents: alias IN container_expression.
type ArrayIterationExpr struct {
	Alias  string
	Source TableExpr
	Loc    Loc
}

func (*ArrayIterationExpr) nodeTag()   {}
func (*ArrayIterationExpr) tableExpr() {}

// SubqueryExpr represents a parenthesized SELECT in FROM position.
type SubqueryExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*SubqueryExpr) nodeTag()   {}
func (*SubqueryExpr) tableExpr() {}

// ---------------------------------------------------------------------------
// Expression nodes — implement ExprNode
//
// DotAccessExpr and BracketAccessExpr also implement TableExpr since they
// appear in both FROM container_expression and scalar_expression contexts.
// ---------------------------------------------------------------------------

// ColumnRef represents an identifier reference (alias, property name).
type ColumnRef struct {
	Name string
	Loc  Loc
}

func (*ColumnRef) nodeTag()  {}
func (*ColumnRef) exprNode() {}

// DotAccessExpr represents property access: expr.property.
type DotAccessExpr struct {
	Expr     ExprNode
	Property string
	Loc      Loc
}

func (*DotAccessExpr) nodeTag()   {}
func (*DotAccessExpr) exprNode()  {}
func (*DotAccessExpr) tableExpr() {}

// BracketAccessExpr represents bracket access: expr[index].
type BracketAccessExpr struct {
	Expr  ExprNode
	Index ExprNode
	Loc   Loc
}

func (*BracketAccessExpr) nodeTag()   {}
func (*BracketAccessExpr) exprNode()  {}
func (*BracketAccessExpr) tableExpr() {}

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	Op    string
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (*BinaryExpr) nodeTag()  {}
func (*BinaryExpr) exprNode() {}

// UnaryExpr represents a unary operation: op operand.
type UnaryExpr struct {
	Op      string
	Operand ExprNode
	Loc     Loc
}

func (*UnaryExpr) nodeTag()  {}
func (*UnaryExpr) exprNode() {}

// TernaryExpr represents: cond ? then : else.
type TernaryExpr struct {
	Cond ExprNode
	Then ExprNode
	Else ExprNode
	Loc  Loc
}

func (*TernaryExpr) nodeTag()  {}
func (*TernaryExpr) exprNode() {}

// InExpr represents: expr [NOT] IN (list).
type InExpr struct {
	Expr ExprNode
	List []ExprNode
	Not  bool
	Loc  Loc
}

func (*InExpr) nodeTag()  {}
func (*InExpr) exprNode() {}

// BetweenExpr represents: expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Expr ExprNode
	Low  ExprNode
	High ExprNode
	Not  bool
	Loc  Loc
}

func (*BetweenExpr) nodeTag()  {}
func (*BetweenExpr) exprNode() {}

// LikeExpr represents: expr [NOT] LIKE pattern [ESCAPE escape].
type LikeExpr struct {
	Expr    ExprNode
	Pattern ExprNode
	Escape  ExprNode // nil if no ESCAPE clause
	Not     bool
	Loc     Loc
}

func (*LikeExpr) nodeTag()  {}
func (*LikeExpr) exprNode() {}

// FuncCall represents a built-in function call: name([*|args]).
type FuncCall struct {
	Name string
	Args []ExprNode
	Star bool // true for COUNT(*)
	Loc  Loc
}

func (*FuncCall) nodeTag()  {}
func (*FuncCall) exprNode() {}

// UDFCall represents a user-defined function call: udf.name(args).
type UDFCall struct {
	Name string
	Args []ExprNode
	Loc  Loc
}

func (*UDFCall) nodeTag()  {}
func (*UDFCall) exprNode() {}

// ExistsExpr represents EXISTS(SELECT ...).
type ExistsExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*ExistsExpr) nodeTag()  {}
func (*ExistsExpr) exprNode() {}

// ArrayExpr represents ARRAY(SELECT ...).
type ArrayExpr struct {
	Select *SelectStmt
	Loc    Loc
}

func (*ArrayExpr) nodeTag()  {}
func (*ArrayExpr) exprNode() {}

// CreateArrayExpr represents an array literal: [elem, ...].
type CreateArrayExpr struct {
	Elements []ExprNode
	Loc      Loc
}

func (*CreateArrayExpr) nodeTag()  {}
func (*CreateArrayExpr) exprNode() {}

// CreateObjectExpr represents an object literal: {key: val, ...}.
type CreateObjectExpr struct {
	Fields []*ObjectFieldPair
	Loc    Loc
}

func (*CreateObjectExpr) nodeTag()  {}
func (*CreateObjectExpr) exprNode() {}

// ObjectFieldPair represents one field in an object literal: key: value.
type ObjectFieldPair struct {
	Key   string
	Value ExprNode
	Loc   Loc
}

func (*ObjectFieldPair) nodeTag() {}

// ParamRef represents a parameter reference: @name.
type ParamRef struct {
	Name string
	Loc  Loc
}

func (*ParamRef) nodeTag()  {}
func (*ParamRef) exprNode() {}

// SubLink represents a parenthesized subquery in expression position: (SELECT ...).
type SubLink struct {
	Select *SelectStmt
	Loc    Loc
}

func (*SubLink) nodeTag()  {}
func (*SubLink) exprNode() {}
