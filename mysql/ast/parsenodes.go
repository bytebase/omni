package ast

// -----------------------------------------------------------------------
// Statement nodes
// -----------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	Loc           Loc
	CTEs          []*CommonTableExpr // WITH clause (CTEs)
	DistinctKind  DistinctKind       // DISTINCT, ALL, or none
	CalcFoundRows bool               // SQL_CALC_FOUND_ROWS
	HighPriority  bool               // HIGH_PRIORITY
	StraightJoin  bool               // STRAIGHT_JOIN hint
	TargetList    []ExprNode         // select expressions (ResTarget)
	From          []TableExpr        // FROM clause (TableRef / JoinClause)
	Where         ExprNode           // WHERE condition
	GroupBy       []ExprNode         // GROUP BY expressions
	Having        ExprNode           // HAVING condition
	OrderBy       []*OrderByItem     // ORDER BY clause
	Limit         *Limit             // LIMIT / OFFSET
	ForUpdate     *ForUpdate         // FOR UPDATE / SHARE
	WindowClause  []*WindowDef       // WINDOW clause
	Into          *IntoClause        // INTO clause
	SetOp         SetOperation       // UNION / INTERSECT / EXCEPT
	SetAll        bool               // ALL modifier for set operations
	Left          *SelectStmt        // left side of set operation
	Right         *SelectStmt        // right side of set operation
}

// CommonTableExpr represents a single CTE in a WITH clause.
type CommonTableExpr struct {
	Loc       Loc
	Name      string
	Columns   []string
	Select    *SelectStmt
	Recursive bool
}

func (c *CommonTableExpr) nodeTag() {}

func (s *SelectStmt) nodeTag()  {}
func (s *SelectStmt) stmtNode() {}

// DistinctKind enumerates DISTINCT modes.
type DistinctKind int

const (
	DistinctNone DistinctKind = iota
	DistinctAll
	DistinctOn
)

// SetOperation enumerates set operations.
type SetOperation int

const (
	SetOpNone SetOperation = iota
	SetOpUnion
	SetOpIntersect
	SetOpExcept
)

// InsertStmt represents an INSERT or REPLACE statement.
type InsertStmt struct {
	Loc            Loc
	IsReplace      bool // REPLACE instead of INSERT
	Priority       InsertPriority
	Ignore         bool          // INSERT IGNORE
	Table          *TableRef     // target table
	Columns        []*ColumnRef  // explicit column list
	Values         [][]ExprNode  // VALUES rows
	Select         *SelectStmt   // INSERT ... SELECT
	SetList        []*Assignment // INSERT ... SET col=val (MySQL extension)
	OnDuplicateKey []*Assignment // ON DUPLICATE KEY UPDATE
	Returning      []Node        // not standard MySQL but included for completeness
}

func (s *InsertStmt) nodeTag()  {}
func (s *InsertStmt) stmtNode() {}

// InsertPriority enumerates INSERT priority modifiers.
type InsertPriority int

const (
	InsertPriorityNone InsertPriority = iota
	InsertPriorityLow
	InsertPriorityDelayed
	InsertPriorityHigh
)

// UpdateStmt represents an UPDATE statement.
type UpdateStmt struct {
	Loc         Loc
	LowPriority bool
	Ignore      bool
	Tables      []TableExpr    // table references (multi-table update)
	SetList     []*Assignment  // SET clause
	Where       ExprNode       // WHERE condition
	OrderBy     []*OrderByItem // ORDER BY
	Limit       *Limit         // LIMIT
}

func (s *UpdateStmt) nodeTag()  {}
func (s *UpdateStmt) stmtNode() {}

// DeleteStmt represents a DELETE statement.
type DeleteStmt struct {
	Loc         Loc
	LowPriority bool
	Quick       bool
	Ignore      bool
	Tables      []TableExpr    // tables to delete from
	Using       []TableExpr    // USING clause (multi-table delete)
	Where       ExprNode       // WHERE condition
	OrderBy     []*OrderByItem // ORDER BY
	Limit       *Limit         // LIMIT
}

func (s *DeleteStmt) nodeTag()  {}
func (s *DeleteStmt) stmtNode() {}

// CreateTableStmt represents a CREATE TABLE statement.
type CreateTableStmt struct {
	Loc         Loc
	IfNotExists bool
	Temporary   bool
	Table       *TableRef
	Columns     []*ColumnDef
	Constraints []*Constraint
	Options     []*TableOption
	Partitions  *PartitionClause
	Like        *TableRef   // CREATE TABLE ... LIKE
	Select      *SelectStmt // CREATE TABLE ... AS SELECT
}

func (s *CreateTableStmt) nodeTag()  {}
func (s *CreateTableStmt) stmtNode() {}

// AlterTableStmt represents an ALTER TABLE statement.
type AlterTableStmt struct {
	Loc      Loc
	Table    *TableRef
	Commands []*AlterTableCmd
}

func (s *AlterTableStmt) nodeTag()  {}
func (s *AlterTableStmt) stmtNode() {}

// AlterTableCmdType enumerates ALTER TABLE operations.
type AlterTableCmdType int

const (
	ATAddColumn AlterTableCmdType = iota
	ATDropColumn
	ATModifyColumn
	ATChangeColumn
	ATAlterColumnDefault
	ATAlterColumnSetNotNull
	ATAlterColumnDropNotNull
	ATAddConstraint
	ATDropConstraint
	ATAddIndex
	ATDropIndex
	ATRenameTable
	ATRenameColumn
	ATRenameIndex
	ATAddPartition
	ATDropPartition
	ATTableOption
	ATConvertCharset
	ATAlgorithm
	ATLock
)

// AlterTableCmd represents a single ALTER TABLE operation.
type AlterTableCmd struct {
	Loc        Loc
	Type       AlterTableCmdType
	Name       string       // column/constraint/index name
	NewName    string       // for RENAME operations
	Column     *ColumnDef   // for ADD/MODIFY/CHANGE COLUMN
	Constraint *Constraint  // for ADD CONSTRAINT
	Option     *TableOption // for table options
	IfExists   bool
	First      bool   // FIRST positioning
	After      string // AFTER column positioning
}

func (c *AlterTableCmd) nodeTag() {}

// DropTableStmt represents a DROP TABLE statement.
type DropTableStmt struct {
	Loc       Loc
	IfExists  bool
	Tables    []*TableRef
	Temporary bool
	Cascade   bool
	Restrict  bool
}

func (s *DropTableStmt) nodeTag()  {}
func (s *DropTableStmt) stmtNode() {}

// CreateIndexStmt represents a CREATE INDEX statement.
type CreateIndexStmt struct {
	Loc         Loc
	Unique      bool
	Fulltext    bool
	Spatial     bool
	IfNotExists bool
	IndexName   string
	Table       *TableRef
	IndexType   string // BTREE, HASH
	Columns     []*IndexColumn
	Options     []*IndexOption
	Algorithm   string
	Lock        string
}

func (s *CreateIndexStmt) nodeTag()  {}
func (s *CreateIndexStmt) stmtNode() {}

// IndexColumn represents a column in an index definition.
type IndexColumn struct {
	Loc    Loc
	Expr   ExprNode // column name or expression
	Length int      // prefix length
	Desc   bool     // DESC ordering
}

func (c *IndexColumn) nodeTag() {}

// IndexOption represents an index option.
type IndexOption struct {
	Loc   Loc
	Name  string
	Value Node // heterogeneous: may be *String, *Integer, or other literal nodes
}

func (o *IndexOption) nodeTag() {}

// CreateViewStmt represents a CREATE VIEW statement.
type CreateViewStmt struct {
	Loc         Loc
	OrReplace   bool
	Algorithm   string // UNDEFINED, MERGE, TEMPTABLE
	Definer     string
	SqlSecurity string // DEFINER, INVOKER
	Name        *TableRef
	Columns     []string
	Select      *SelectStmt
	CheckOption string // CASCADED, LOCAL
}

func (s *CreateViewStmt) nodeTag()  {}
func (s *CreateViewStmt) stmtNode() {}

// CreateDatabaseStmt represents a CREATE DATABASE statement.
type CreateDatabaseStmt struct {
	Loc         Loc
	IfNotExists bool
	Name        string
	Options     []*DatabaseOption
}

func (s *CreateDatabaseStmt) nodeTag()  {}
func (s *CreateDatabaseStmt) stmtNode() {}

// DatabaseOption represents a database option (CHARACTER SET, COLLATE, etc.).
type DatabaseOption struct {
	Loc   Loc
	Name  string
	Value string
}

func (o *DatabaseOption) nodeTag() {}

// TruncateStmt represents a TRUNCATE TABLE statement.
type TruncateStmt struct {
	Loc    Loc
	Tables []*TableRef
}

func (s *TruncateStmt) nodeTag()  {}
func (s *TruncateStmt) stmtNode() {}

// RenameTableStmt represents a RENAME TABLE statement.
type RenameTableStmt struct {
	Loc   Loc
	Pairs []*RenameTablePair
}

func (s *RenameTableStmt) nodeTag()  {}
func (s *RenameTableStmt) stmtNode() {}

// RenameTablePair is a single rename pair (old -> new).
type RenameTablePair struct {
	Loc Loc
	Old *TableRef
	New *TableRef
}

func (p *RenameTablePair) nodeTag() {}

// -----------------------------------------------------------------------
// Table components
// -----------------------------------------------------------------------

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Loc           Loc
	Name          string
	TypeName      *DataType
	Constraints   []*ColumnConstraint
	DefaultValue  ExprNode
	Comment       string
	AutoIncrement bool
	OnUpdate      ExprNode // ON UPDATE CURRENT_TIMESTAMP
	Generated     *GeneratedColumn
}

func (d *ColumnDef) nodeTag() {}

// GeneratedColumn represents a generated column specification.
type GeneratedColumn struct {
	Loc    Loc
	Expr   ExprNode
	Stored bool // STORED vs VIRTUAL
}

func (g *GeneratedColumn) nodeTag() {}

// ColumnConstraintType enumerates column constraint types.
type ColumnConstraintType int

const (
	ColConstrNotNull ColumnConstraintType = iota
	ColConstrNull
	ColConstrDefault
	ColConstrPrimaryKey
	ColConstrUnique
	ColConstrCheck
	ColConstrReferences
	ColConstrComment
	ColConstrCollate
	ColConstrAutoIncrement
)

// ColumnConstraint represents a column-level constraint.
type ColumnConstraint struct {
	Loc        Loc
	Type       ColumnConstraintType
	Name       string   // constraint name
	Expr       ExprNode // for CHECK, DEFAULT, etc.
	RefTable   *TableRef
	RefColumns []string
	OnDelete   ReferenceAction
	OnUpdate   ReferenceAction
}

func (c *ColumnConstraint) nodeTag() {}

// ConstraintType enumerates table-level constraint types.
type ConstraintType int

const (
	ConstrPrimaryKey ConstraintType = iota
	ConstrUnique
	ConstrForeignKey
	ConstrCheck
	ConstrIndex
	ConstrFulltextIndex
	ConstrSpatialIndex
)

// Constraint represents a table-level constraint.
type Constraint struct {
	Loc        Loc
	Type       ConstraintType
	Name       string
	Columns    []string
	IndexType  string    // BTREE, HASH
	Expr       ExprNode  // for CHECK
	RefTable   *TableRef // for FOREIGN KEY
	RefColumns []string  // for FOREIGN KEY
	OnDelete   ReferenceAction
	OnUpdate   ReferenceAction
	Match      string // FULL, PARTIAL, SIMPLE
}

func (c *Constraint) nodeTag() {}

// ReferenceAction enumerates foreign key actions.
type ReferenceAction int

const (
	RefActNone ReferenceAction = iota
	RefActRestrict
	RefActCascade
	RefActSetNull
	RefActSetDefault
	RefActNoAction
)

// TableOption represents a table option (ENGINE, CHARSET, etc.).
type TableOption struct {
	Loc   Loc
	Name  string
	Value string
}

func (o *TableOption) nodeTag() {}

// PartitionClause represents a PARTITION BY clause.
type PartitionClause struct {
	Loc        Loc
	Type       PartitionType
	Expr       ExprNode
	Columns    []string
	NumParts   int
	Partitions []*PartitionDef
}

func (p *PartitionClause) nodeTag() {}

// PartitionType enumerates partition types.
type PartitionType int

const (
	PartitionRange PartitionType = iota
	PartitionList
	PartitionHash
	PartitionKey
)

// PartitionDef represents a single partition definition.
type PartitionDef struct {
	Loc     Loc
	Name    string
	Values  Node // VALUES LESS THAN (expr) or IN (list) — kept as Node to support heterogeneous partition value types
	Options []*TableOption
}

func (d *PartitionDef) nodeTag() {}

// DataType represents a MySQL data type.
type DataType struct {
	Loc        Loc
	Name       string // INT, VARCHAR, DECIMAL, etc.
	Length     int    // type length (e.g., VARCHAR(255))
	Scale      int    // decimal scale (e.g., DECIMAL(10,2))
	Unsigned   bool
	Zerofill   bool
	Charset    string   // CHARACTER SET
	Collate    string   // COLLATE
	EnumValues []string // for ENUM and SET types
	ArrayDim   int      // not used in MySQL, but here for consistency
}

func (t *DataType) nodeTag() {}

// -----------------------------------------------------------------------
// Expression nodes
// -----------------------------------------------------------------------

// BinaryOp enumerates binary operator types.
type BinaryOp int

const (
	BinOpAdd BinaryOp = iota
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpMod
	BinOpEq
	BinOpNe
	BinOpLt
	BinOpGt
	BinOpLe
	BinOpGe
	BinOpAnd
	BinOpOr
	BinOpBitAnd
	BinOpBitOr
	BinOpBitXor
	BinOpShiftLeft
	BinOpShiftRight
	BinOpDivInt // DIV
	BinOpRegexp
	BinOpLikeEscape
	BinOpNullSafeEq // <=>
	BinOpAssign        // :=
	BinOpJsonExtract   // ->
	BinOpJsonUnquote   // ->>
)

// BinaryExpr represents a binary expression (a op b).
type BinaryExpr struct {
	Loc   Loc
	Op    BinaryOp
	Left  ExprNode
	Right ExprNode
}

func (e *BinaryExpr) nodeTag()  {}
func (e *BinaryExpr) exprNode() {}

// UnaryOp enumerates unary operator types.
type UnaryOp int

const (
	UnaryMinus UnaryOp = iota
	UnaryPlus
	UnaryNot
	UnaryBitNot
)

// UnaryExpr represents a unary expression (op expr).
type UnaryExpr struct {
	Loc     Loc
	Op      UnaryOp
	Operand ExprNode
}

func (e *UnaryExpr) nodeTag()  {}
func (e *UnaryExpr) exprNode() {}

// FuncCallExpr represents a function call.
type FuncCallExpr struct {
	Loc      Loc
	Name     string
	Schema   string // schema-qualified name
	Args     []ExprNode
	Distinct bool // COUNT(DISTINCT ...)
	Star     bool // COUNT(*)
	OrderBy  []*OrderByItem
	Over     *WindowDef // OVER clause for window functions
}

func (e *FuncCallExpr) nodeTag()  {}
func (e *FuncCallExpr) exprNode() {}

// SubqueryExpr represents a subquery expression.
type SubqueryExpr struct {
	Loc    Loc
	Select *SelectStmt
	Exists bool // EXISTS (SELECT ...)
}

func (e *SubqueryExpr) nodeTag()   {}
func (e *SubqueryExpr) exprNode()  {}
func (e *SubqueryExpr) tableExpr() {}

// CaseExpr represents a CASE expression.
type CaseExpr struct {
	Loc     Loc
	Operand ExprNode // simple CASE operand (nil for searched CASE)
	Whens   []*CaseWhen
	Default ExprNode // ELSE clause
}

func (e *CaseExpr) nodeTag()  {}
func (e *CaseExpr) exprNode() {}

// CaseWhen represents a WHEN clause in a CASE expression.
type CaseWhen struct {
	Loc    Loc
	Cond   ExprNode
	Result ExprNode
}

func (w *CaseWhen) nodeTag() {}

// BetweenExpr represents a BETWEEN expression.
type BetweenExpr struct {
	Loc  Loc
	Not  bool
	Expr ExprNode
	Low  ExprNode
	High ExprNode
}

func (e *BetweenExpr) nodeTag()  {}
func (e *BetweenExpr) exprNode() {}

// InExpr represents an IN expression.
type InExpr struct {
	Loc    Loc
	Not    bool
	Expr   ExprNode
	List   []ExprNode  // value list
	Select *SelectStmt // subquery
}

func (e *InExpr) nodeTag()  {}
func (e *InExpr) exprNode() {}

// LikeExpr represents a LIKE expression.
type LikeExpr struct {
	Loc     Loc
	Not     bool
	Expr    ExprNode
	Pattern ExprNode
	Escape  ExprNode // ESCAPE clause
}

func (e *LikeExpr) nodeTag()  {}
func (e *LikeExpr) exprNode() {}

// IsExpr represents an IS expression (IS NULL, IS TRUE, etc.).
type IsExpr struct {
	Loc  Loc
	Not  bool
	Expr ExprNode
	Test IsTestType
}

func (e *IsExpr) nodeTag()  {}
func (e *IsExpr) exprNode() {}

// IsTestType enumerates IS test types.
type IsTestType int

const (
	IsNull IsTestType = iota
	IsTrue
	IsFalse
	IsUnknown
)

// ExistsExpr represents an EXISTS expression.
type ExistsExpr struct {
	Loc    Loc
	Select *SelectStmt
}

func (e *ExistsExpr) nodeTag()  {}
func (e *ExistsExpr) exprNode() {}

// CastExpr represents a CAST expression.
type CastExpr struct {
	Loc      Loc
	Expr     ExprNode
	TypeName *DataType
}

func (e *CastExpr) nodeTag()  {}
func (e *CastExpr) exprNode() {}

// -----------------------------------------------------------------------
// Reference and literal nodes
// -----------------------------------------------------------------------

// ColumnRef represents a column reference (possibly qualified).
type ColumnRef struct {
	Loc    Loc
	Table  string // table name or alias (may be empty)
	Schema string // schema/database name (may be empty)
	Column string // column name
	Star   bool   // table.*
}

func (r *ColumnRef) nodeTag()  {}
func (r *ColumnRef) exprNode() {}

// TableRef represents a table reference (possibly qualified, with alias).
type TableRef struct {
	Loc        Loc
	Schema     string // database name
	Name       string // table name
	Alias      string // AS alias
	IndexHints []*IndexHint
}

func (r *TableRef) nodeTag()   {}
func (r *TableRef) tableExpr() {}

// IndexHint represents an index hint (USE/FORCE/IGNORE INDEX).
type IndexHint struct {
	Loc     Loc
	Type    IndexHintType
	Scope   IndexHintScope
	Indexes []string
}

func (h *IndexHint) nodeTag() {}

// IndexHintType enumerates index hint types.
type IndexHintType int

const (
	HintUse IndexHintType = iota
	HintForce
	HintIgnore
)

// IndexHintScope enumerates index hint scopes.
type IndexHintScope int

const (
	HintScopeAll IndexHintScope = iota
	HintScopeJoin
	HintScopeOrderBy
	HintScopeGroupBy
)

// IntLit represents an integer literal.
type IntLit struct {
	Loc   Loc
	Value int64
}

func (l *IntLit) nodeTag()  {}
func (l *IntLit) exprNode() {}

// FloatLit represents a floating-point literal.
type FloatLit struct {
	Loc   Loc
	Value string // stored as string to preserve precision
}

func (l *FloatLit) nodeTag()  {}
func (l *FloatLit) exprNode() {}

// StringLit represents a string literal.
type StringLit struct {
	Loc     Loc
	Value   string
	Charset string // optional charset introducer (_utf8, etc.)
}

func (l *StringLit) nodeTag()  {}
func (l *StringLit) exprNode() {}

// BoolLit represents a boolean literal (TRUE/FALSE).
type BoolLit struct {
	Loc   Loc
	Value bool
}

func (l *BoolLit) nodeTag()  {}
func (l *BoolLit) exprNode() {}

// NullLit represents a NULL literal.
type NullLit struct {
	Loc Loc
}

func (l *NullLit) nodeTag()  {}
func (l *NullLit) exprNode() {}

// HexLit represents a hexadecimal literal (0xFF or X'FF').
type HexLit struct {
	Loc   Loc
	Value string
}

func (l *HexLit) nodeTag()  {}
func (l *HexLit) exprNode() {}

// BitLit represents a bit literal (0b101 or b'101').
type BitLit struct {
	Loc   Loc
	Value string
}

func (l *BitLit) nodeTag()  {}
func (l *BitLit) exprNode() {}

// -----------------------------------------------------------------------
// Clause and helper nodes
// -----------------------------------------------------------------------

// Assignment represents a SET assignment (col = value).
type Assignment struct {
	Loc    Loc
	Column *ColumnRef
	Value  ExprNode
}

func (a *Assignment) nodeTag() {}

// OrderByItem represents an ORDER BY item.
type OrderByItem struct {
	Loc        Loc
	Expr       ExprNode
	Desc       bool
	NullsFirst *bool // nil means default
}

func (o *OrderByItem) nodeTag() {}

// Limit represents LIMIT and OFFSET.
type Limit struct {
	Loc    Loc
	Count  ExprNode // LIMIT value
	Offset ExprNode // OFFSET value
}

func (l *Limit) nodeTag() {}

// JoinType enumerates join types.
type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinCross
	JoinNatural
	JoinStraight // STRAIGHT_JOIN (MySQL-specific)
)

// JoinClause represents a JOIN clause.
type JoinClause struct {
	Loc       Loc
	Type      JoinType
	Left      TableExpr // table or join
	Right     TableExpr // table or join
	Condition Node      // ON (*OnCondition) or USING (*UsingCondition) — kept as Node because these are structurally different types
}

func (j *JoinClause) nodeTag()   {}
func (j *JoinClause) tableExpr() {}

// OnCondition represents an ON condition in a JOIN.
type OnCondition struct {
	Loc  Loc
	Expr ExprNode
}

func (c *OnCondition) nodeTag() {}

// UsingCondition represents a USING condition in a JOIN.
type UsingCondition struct {
	Loc     Loc
	Columns []string
}

func (c *UsingCondition) nodeTag() {}

// ResTarget represents a result target in a SELECT list.
type ResTarget struct {
	Loc  Loc
	Name string   // AS alias
	Val  ExprNode // expression
}

func (r *ResTarget) nodeTag()  {}
func (r *ResTarget) exprNode() {}

// WindowDef represents a window definition.
type WindowDef struct {
	Loc         Loc
	Name        string
	RefName     string // reference to existing window
	PartitionBy []ExprNode
	OrderBy     []*OrderByItem
	Frame       *WindowFrame
}

func (w *WindowDef) nodeTag() {}

// WindowFrame represents a window frame specification.
type WindowFrame struct {
	Loc   Loc
	Type  WindowFrameType
	Start *WindowFrameBound
	End   *WindowFrameBound
}

func (f *WindowFrame) nodeTag() {}

// WindowFrameType enumerates window frame types.
type WindowFrameType int

const (
	FrameRows WindowFrameType = iota
	FrameRange
	FrameGroups
)

// WindowFrameBound represents a window frame bound.
type WindowFrameBound struct {
	Loc    Loc
	Type   WindowBoundType
	Offset ExprNode
}

func (b *WindowFrameBound) nodeTag() {}

// WindowBoundType enumerates window frame bound types.
type WindowBoundType int

const (
	BoundUnboundedPreceding WindowBoundType = iota
	BoundPreceding
	BoundCurrentRow
	BoundFollowing
	BoundUnboundedFollowing
)

// IntoClause represents a SELECT ... INTO clause.
type IntoClause struct {
	Loc      Loc
	Outfile  string
	Dumpfile string
	Vars     []*VariableRef
}

func (c *IntoClause) nodeTag() {}

// ForUpdate represents FOR UPDATE / FOR SHARE.
type ForUpdate struct {
	Loc        Loc
	Share      bool // FOR SHARE instead of FOR UPDATE
	Tables     []*TableRef
	NoWait     bool
	SkipLocked bool
}

func (f *ForUpdate) nodeTag() {}

// VariableRef represents a variable reference (@var or @@var).
type VariableRef struct {
	Loc    Loc
	Name   string
	System bool   // @@ (system variable) vs @ (user variable)
	Scope  string // GLOBAL, SESSION, LOCAL for system variables
}

func (v *VariableRef) nodeTag()  {}
func (v *VariableRef) exprNode() {}

// -----------------------------------------------------------------------
// Additional statement nodes
// -----------------------------------------------------------------------

// SetStmt represents a SET statement.
type SetStmt struct {
	Loc         Loc
	Scope       string // GLOBAL, SESSION, LOCAL
	Assignments []*Assignment
}

func (s *SetStmt) nodeTag()  {}
func (s *SetStmt) stmtNode() {}

// ShowStmt represents a SHOW statement.
type ShowStmt struct {
	Loc   Loc
	Type  string // DATABASES, TABLES, COLUMNS, etc.
	From  *TableRef
	Like  ExprNode
	Where ExprNode
}

func (s *ShowStmt) nodeTag()  {}
func (s *ShowStmt) stmtNode() {}

// UseStmt represents a USE statement.
type UseStmt struct {
	Loc      Loc
	Database string
}

func (s *UseStmt) nodeTag()  {}
func (s *UseStmt) stmtNode() {}

// ExplainStmt represents an EXPLAIN/DESCRIBE statement.
type ExplainStmt struct {
	Loc    Loc
	Format string // TRADITIONAL, JSON, TREE
	Stmt   StmtNode
}

func (s *ExplainStmt) nodeTag()  {}
func (s *ExplainStmt) stmtNode() {}

// BeginStmt represents BEGIN / START TRANSACTION.
type BeginStmt struct {
	Loc                    Loc
	ReadOnly               bool
	ReadWrite              bool
	WithConsistentSnapshot bool
}

func (s *BeginStmt) nodeTag()  {}
func (s *BeginStmt) stmtNode() {}

// CommitStmt represents a COMMIT statement.
type CommitStmt struct {
	Loc     Loc
	Chain   bool
	Release bool
}

func (s *CommitStmt) nodeTag()  {}
func (s *CommitStmt) stmtNode() {}

// RollbackStmt represents a ROLLBACK statement.
type RollbackStmt struct {
	Loc       Loc
	Savepoint string
	Chain     bool
	Release   bool
}

func (s *RollbackStmt) nodeTag()  {}
func (s *RollbackStmt) stmtNode() {}

// SavepointStmt represents a SAVEPOINT statement.
type SavepointStmt struct {
	Loc  Loc
	Name string
}

func (s *SavepointStmt) nodeTag()  {}
func (s *SavepointStmt) stmtNode() {}

// LockTablesStmt represents a LOCK TABLES statement.
type LockTablesStmt struct {
	Loc    Loc
	Tables []*LockTable
}

func (s *LockTablesStmt) nodeTag()  {}
func (s *LockTablesStmt) stmtNode() {}

// LockTable represents a table lock specification.
type LockTable struct {
	Loc      Loc
	Table    *TableRef
	LockType string // READ, WRITE, READ LOCAL, LOW_PRIORITY WRITE
}

func (l *LockTable) nodeTag() {}

// UnlockTablesStmt represents an UNLOCK TABLES statement.
type UnlockTablesStmt struct {
	Loc Loc
}

func (s *UnlockTablesStmt) nodeTag()  {}
func (s *UnlockTablesStmt) stmtNode() {}

// GrantStmt represents a GRANT statement.
type GrantStmt struct {
	Loc        Loc
	Privileges []string
	AllPriv    bool
	On         *GrantTarget
	To         []string
	WithGrant  bool
}

func (s *GrantStmt) nodeTag()  {}
func (s *GrantStmt) stmtNode() {}

// RevokeStmt represents a REVOKE statement.
type RevokeStmt struct {
	Loc        Loc
	Privileges []string
	AllPriv    bool
	On         *GrantTarget
	From       []string
}

func (s *RevokeStmt) nodeTag()  {}
func (s *RevokeStmt) stmtNode() {}

// GrantTarget represents the target of a GRANT/REVOKE.
type GrantTarget struct {
	Loc  Loc
	Type string // TABLE, DATABASE, PROCEDURE, FUNCTION, etc.
	Name *TableRef
}

func (t *GrantTarget) nodeTag() {}

// CreateUserStmt represents a CREATE USER statement.
type CreateUserStmt struct {
	Loc         Loc
	IfNotExists bool
	Users       []*UserSpec
}

func (s *CreateUserStmt) nodeTag()  {}
func (s *CreateUserStmt) stmtNode() {}

// DropUserStmt represents a DROP USER statement.
type DropUserStmt struct {
	Loc      Loc
	IfExists bool
	Users    []string
}

func (s *DropUserStmt) nodeTag()  {}
func (s *DropUserStmt) stmtNode() {}

// AlterUserStmt represents an ALTER USER statement.
type AlterUserStmt struct {
	Loc   Loc
	Users []*UserSpec
}

func (s *AlterUserStmt) nodeTag()  {}
func (s *AlterUserStmt) stmtNode() {}

// UserSpec represents a user specification.
type UserSpec struct {
	Loc        Loc
	Name       string
	Host       string
	AuthPlugin string
	Password   string
}

func (u *UserSpec) nodeTag() {}

// CreateFunctionStmt represents a CREATE FUNCTION or CREATE PROCEDURE.
type CreateFunctionStmt struct {
	Loc             Loc
	OrReplace       bool
	IsProcedure     bool
	Definer         string
	Name            *TableRef
	Params          []*FuncParam
	Returns         *DataType
	Body            string
	Characteristics []*RoutineCharacteristic
}

func (s *CreateFunctionStmt) nodeTag()  {}
func (s *CreateFunctionStmt) stmtNode() {}

// FuncParam represents a function/procedure parameter.
type FuncParam struct {
	Loc       Loc
	Direction string // IN, OUT, INOUT
	Name      string
	TypeName  *DataType
}

func (p *FuncParam) nodeTag() {}

// RoutineCharacteristic represents a routine characteristic (DETERMINISTIC, etc.).
type RoutineCharacteristic struct {
	Loc   Loc
	Name  string
	Value string
}

func (c *RoutineCharacteristic) nodeTag() {}

// CreateTriggerStmt represents a CREATE TRIGGER statement.
type CreateTriggerStmt struct {
	Loc     Loc
	Definer string
	Name    string
	Timing  string // BEFORE, AFTER
	Event   string // INSERT, UPDATE, DELETE
	Table   *TableRef
	Order   *TriggerOrder
	Body    string
}

func (s *CreateTriggerStmt) nodeTag()  {}
func (s *CreateTriggerStmt) stmtNode() {}

// TriggerOrder represents FOLLOWS/PRECEDES in trigger definition.
type TriggerOrder struct {
	Loc         Loc
	Follows     bool
	TriggerName string
}

func (o *TriggerOrder) nodeTag() {}

// CreateEventStmt represents a CREATE EVENT statement.
type CreateEventStmt struct {
	Loc          Loc
	IfNotExists  bool
	Definer      string
	Name         string
	Schedule     *EventSchedule
	OnCompletion string
	Enable       string // ENABLE, DISABLE, DISABLE ON SLAVE
	Comment      string
	Body         string
}

func (s *CreateEventStmt) nodeTag()  {}
func (s *CreateEventStmt) stmtNode() {}

// EventSchedule represents an event schedule.
type EventSchedule struct {
	Loc    Loc
	At     ExprNode
	Every  ExprNode
	Starts ExprNode
	Ends   ExprNode
}

func (s *EventSchedule) nodeTag() {}

// LoadDataStmt represents a LOAD DATA statement.
type LoadDataStmt struct {
	Loc                Loc
	Local              bool
	Infile             string
	Replace            bool
	Ignore             bool
	Table              *TableRef
	Columns            []*ColumnRef
	SetList            []*Assignment
	LinesTerminatedBy  string
	FieldsTerminatedBy string
	FieldsEnclosedBy   string
	FieldsEscapedBy    string
	IgnoreRows         int
}

func (s *LoadDataStmt) nodeTag()  {}
func (s *LoadDataStmt) stmtNode() {}

// PrepareStmt represents a PREPARE statement.
type PrepareStmt struct {
	Loc  Loc
	Name string
	Stmt string // SQL text
}

func (s *PrepareStmt) nodeTag()  {}
func (s *PrepareStmt) stmtNode() {}

// ExecuteStmt represents an EXECUTE statement.
type ExecuteStmt struct {
	Loc    Loc
	Name   string
	Params []ExprNode
}

func (s *ExecuteStmt) nodeTag()  {}
func (s *ExecuteStmt) stmtNode() {}

// DeallocateStmt represents a DEALLOCATE PREPARE statement.
type DeallocateStmt struct {
	Loc  Loc
	Name string
}

func (s *DeallocateStmt) nodeTag()  {}
func (s *DeallocateStmt) stmtNode() {}

// AnalyzeTableStmt represents an ANALYZE TABLE statement.
type AnalyzeTableStmt struct {
	Loc    Loc
	Tables []*TableRef
}

func (s *AnalyzeTableStmt) nodeTag()  {}
func (s *AnalyzeTableStmt) stmtNode() {}

// OptimizeTableStmt represents an OPTIMIZE TABLE statement.
type OptimizeTableStmt struct {
	Loc    Loc
	Tables []*TableRef
}

func (s *OptimizeTableStmt) nodeTag()  {}
func (s *OptimizeTableStmt) stmtNode() {}

// CheckTableStmt represents a CHECK TABLE statement.
type CheckTableStmt struct {
	Loc     Loc
	Tables  []*TableRef
	Options []string
}

func (s *CheckTableStmt) nodeTag()  {}
func (s *CheckTableStmt) stmtNode() {}

// RepairTableStmt represents a REPAIR TABLE statement.
type RepairTableStmt struct {
	Loc      Loc
	Tables   []*TableRef
	Quick    bool
	Extended bool
}

func (s *RepairTableStmt) nodeTag()  {}
func (s *RepairTableStmt) stmtNode() {}

// FlushStmt represents a FLUSH statement.
type FlushStmt struct {
	Loc     Loc
	Options []string
}

func (s *FlushStmt) nodeTag()  {}
func (s *FlushStmt) stmtNode() {}

// KillStmt represents a KILL statement.
type KillStmt struct {
	Loc          Loc
	ConnectionID ExprNode
	Query        bool // KILL QUERY vs KILL CONNECTION
}

func (s *KillStmt) nodeTag()  {}
func (s *KillStmt) stmtNode() {}

// DoStmt represents a DO statement.
type DoStmt struct {
	Loc   Loc
	Exprs []ExprNode
}

func (s *DoStmt) nodeTag()  {}
func (s *DoStmt) stmtNode() {}

// AlterDatabaseStmt represents an ALTER DATABASE statement.
type AlterDatabaseStmt struct {
	Loc     Loc
	Name    string
	Options []*DatabaseOption
}

func (s *AlterDatabaseStmt) nodeTag()  {}
func (s *AlterDatabaseStmt) stmtNode() {}

// DropDatabaseStmt represents a DROP DATABASE statement.
type DropDatabaseStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
}

func (s *DropDatabaseStmt) nodeTag()  {}
func (s *DropDatabaseStmt) stmtNode() {}

// DropIndexStmt represents a DROP INDEX statement.
type DropIndexStmt struct {
	Loc       Loc
	Name      string
	Table     *TableRef
	Algorithm string
	Lock      string
}

func (s *DropIndexStmt) nodeTag()  {}
func (s *DropIndexStmt) stmtNode() {}

// DropViewStmt represents a DROP VIEW statement.
type DropViewStmt struct {
	Loc      Loc
	IfExists bool
	Views    []*TableRef
	Cascade  bool
	Restrict bool
}

func (s *DropViewStmt) nodeTag()  {}
func (s *DropViewStmt) stmtNode() {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Loc  Loc
	Expr ExprNode
}

func (e *ParenExpr) nodeTag()  {}
func (e *ParenExpr) exprNode() {}

// StarExpr represents a * (all columns) expression.
type StarExpr struct {
	Loc Loc
}

func (e *StarExpr) nodeTag()  {}
func (e *StarExpr) exprNode() {}

// IntervalExpr represents an INTERVAL expression.
type IntervalExpr struct {
	Loc   Loc
	Value ExprNode
	Unit  string
}

func (e *IntervalExpr) nodeTag()  {}
func (e *IntervalExpr) exprNode() {}

// CollateExpr represents a COLLATE expression.
type CollateExpr struct {
	Loc       Loc
	Expr      ExprNode
	Collation string
}

func (e *CollateExpr) nodeTag()  {}
func (e *CollateExpr) exprNode() {}

// MatchExpr represents a MATCH ... AGAINST expression (fulltext search).
type MatchExpr struct {
	Loc      Loc
	Columns  []*ColumnRef
	Against  ExprNode
	Modifier string // IN NATURAL LANGUAGE MODE, IN BOOLEAN MODE, etc.
}

func (e *MatchExpr) nodeTag()  {}
func (e *MatchExpr) exprNode() {}

// ConvertExpr represents a CONVERT expression.
type ConvertExpr struct {
	Loc      Loc
	Expr     ExprNode
	TypeName *DataType
	Charset  string // CONVERT(expr USING charset)
}

func (e *ConvertExpr) nodeTag()  {}
func (e *ConvertExpr) exprNode() {}

// DefaultExpr represents the DEFAULT keyword used as an expression.
type DefaultExpr struct {
	Loc Loc
}

func (e *DefaultExpr) nodeTag()  {}
func (e *DefaultExpr) exprNode() {}

// MemberOfExpr represents a value MEMBER OF(json_array) expression.
type MemberOfExpr struct {
	Loc   Loc
	Value ExprNode
	Array ExprNode
}

func (e *MemberOfExpr) nodeTag()  {}
func (e *MemberOfExpr) exprNode() {}

// JsonTableExpr represents a JSON_TABLE() table function.
type JsonTableExpr struct {
	Loc     Loc
	Expr    ExprNode // JSON document expression
	Path    ExprNode // JSON path expression
	Columns []*JsonTableColumn
	Alias   string
}

func (e *JsonTableExpr) nodeTag()   {}
func (e *JsonTableExpr) tableExpr() {}

// JsonTableColumn represents a column definition in JSON_TABLE.
type JsonTableColumn struct {
	Loc          Loc
	Name         string
	TypeName     *DataType
	Path         string
	Exists       bool   // FOR ORDINALITY or EXISTS PATH
	ErrorOption  string // ERROR ON ERROR, NULL ON ERROR, DEFAULT val ON ERROR
	EmptyOption  string // ERROR ON EMPTY, NULL ON EMPTY, DEFAULT val ON EMPTY
	Ordinality   bool   // FOR ORDINALITY
	Nested       bool   // NESTED PATH
	NestedPath   string
	NestedCols   []*JsonTableColumn
}

func (c *JsonTableColumn) nodeTag() {}

// SetTransactionStmt represents a SET TRANSACTION statement.
type SetTransactionStmt struct {
	Loc            Loc
	Scope          string // GLOBAL, SESSION, or ""
	IsolationLevel string // READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE
	AccessMode     string // READ ONLY, READ WRITE
}

func (s *SetTransactionStmt) nodeTag()  {}
func (s *SetTransactionStmt) stmtNode() {}

// XAStmtType enumerates XA statement types.
type XAStmtType int

const (
	XAStart XAStmtType = iota
	XAEnd
	XAPrepare
	XACommit
	XARollback
	XARecover
)

// XAStmt represents an XA distributed transaction statement.
type XAStmt struct {
	Loc      Loc
	Type     XAStmtType
	Xid      []ExprNode // gtrid [, bqual [, formatID]]
	Join     bool       // JOIN option for XA START
	Resume   bool       // RESUME option for XA START
	Suspend  bool       // SUSPEND option for XA END
	Migrate  bool       // FOR MIGRATE option for XA END
	OnePhase bool       // ONE PHASE option for XA COMMIT
	Convert  bool       // CONVERT XID option for XA RECOVER
}

func (s *XAStmt) nodeTag()  {}
func (s *XAStmt) stmtNode() {}

// CallStmt represents a CALL procedure_name([arg[,arg]...]) statement.
type CallStmt struct {
	Loc  Loc
	Name *TableRef
	Args []ExprNode
}

func (s *CallStmt) nodeTag()  {}
func (s *CallStmt) stmtNode() {}

// HandlerOpenStmt represents a HANDLER tbl_name OPEN [AS alias] statement.
type HandlerOpenStmt struct {
	Loc   Loc
	Table *TableRef
	Alias string
}

func (s *HandlerOpenStmt) nodeTag()  {}
func (s *HandlerOpenStmt) stmtNode() {}

// HandlerReadStmt represents a HANDLER tbl_name READ ... statement.
type HandlerReadStmt struct {
	Loc       Loc
	Table     *TableRef
	Direction string   // FIRST, NEXT, PREV, LAST
	Index     string   // index name for index-based reads
	Where     ExprNode // WHERE condition
	Limit     *Limit   // LIMIT clause
}

func (s *HandlerReadStmt) nodeTag()  {}
func (s *HandlerReadStmt) stmtNode() {}

// HandlerCloseStmt represents a HANDLER tbl_name CLOSE statement.
type HandlerCloseStmt struct {
	Loc   Loc
	Table *TableRef
}

func (s *HandlerCloseStmt) nodeTag()  {}
func (s *HandlerCloseStmt) stmtNode() {}

// SignalStmt represents a SIGNAL statement.
type SignalStmt struct {
	Loc            Loc
	ConditionValue string            // SQLSTATE value or condition name
	SetItems       []*SignalInfoItem // SET clause items
}

func (s *SignalStmt) nodeTag()  {}
func (s *SignalStmt) stmtNode() {}

// ResignalStmt represents a RESIGNAL statement.
type ResignalStmt struct {
	Loc            Loc
	ConditionValue string            // SQLSTATE value or condition name (may be empty)
	SetItems       []*SignalInfoItem // SET clause items
}

func (s *ResignalStmt) nodeTag()  {}
func (s *ResignalStmt) stmtNode() {}

// SignalInfoItem represents a signal information item (condition_info_name = value).
type SignalInfoItem struct {
	Loc   Loc
	Name  string   // e.g. MESSAGE_TEXT, MYSQL_ERRNO, CLASS_ORIGIN, etc.
	Value ExprNode // simple_value_specification
}

func (s *SignalInfoItem) nodeTag() {}

// GetDiagnosticsStmt represents a GET DIAGNOSTICS statement.
type GetDiagnosticsStmt struct {
	Loc             Loc
	Stacked         bool               // STACKED (vs CURRENT, which is default)
	StatementInfo   bool               // true = statement info items, false = condition info items
	ConditionNumber ExprNode           // condition number (only when StatementInfo is false)
	Items           []*DiagnosticsItem // target = item_name pairs
}

func (s *GetDiagnosticsStmt) nodeTag()  {}
func (s *GetDiagnosticsStmt) stmtNode() {}

// DiagnosticsItem represents a diagnostics item (target = item_name).
type DiagnosticsItem struct {
	Loc    Loc
	Target *VariableRef // target variable
	Name   string       // item name (e.g. NUMBER, ROW_COUNT, MESSAGE_TEXT, etc.)
}

func (d *DiagnosticsItem) nodeTag() {}

// RawStmt wraps a statement node with its location span.
type RawStmt struct {
	Loc     Loc
	StmtLen int
	Stmt    StmtNode
}

func (s *RawStmt) nodeTag()  {}
func (s *RawStmt) stmtNode() {}
