package ast

// -----------------------------------------------------------------------
// Statement nodes
// -----------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	Loc               Loc
	CTEs              []*CommonTableExpr // WITH clause (CTEs)
	DistinctKind      DistinctKind       // DISTINCT, ALL, or none
	CalcFoundRows     bool               // SQL_CALC_FOUND_ROWS
	HighPriority      bool               // HIGH_PRIORITY
	StraightJoin      bool               // STRAIGHT_JOIN hint
	SmallResult       bool               // SQL_SMALL_RESULT
	BigResult         bool               // SQL_BIG_RESULT
	BufferResult      bool               // SQL_BUFFER_RESULT
	NoCache           bool               // SQL_NO_CACHE
	TargetList        []ExprNode         // select expressions (ResTarget)
	From              []TableExpr        // FROM clause (TableRef / JoinClause)
	Where             ExprNode           // WHERE condition
	GroupBy           []ExprNode         // GROUP BY expressions
	WithRollup        bool               // GROUP BY ... WITH ROLLUP
	Having            ExprNode           // HAVING condition
	OrderBy           []*OrderByItem     // ORDER BY clause
	OrderByWithRollup bool               // ORDER BY ... WITH ROLLUP
	Limit             *Limit             // LIMIT / OFFSET
	ForUpdate         *ForUpdate         // FOR UPDATE / SHARE
	WindowClause      []*WindowDef       // WINDOW clause
	Into              *IntoClause        // INTO clause
	SetOp             SetOperation       // UNION / INTERSECT / EXCEPT
	SetAll            bool               // ALL modifier for set operations
	Left              *SelectStmt        // left side of set operation
	Right             *SelectStmt        // right side of set operation
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
	Partitions     []string      // PARTITION (p0, p1, ...)
	Columns        []*ColumnRef  // explicit column list
	Values         [][]ExprNode  // VALUES rows
	Select         *SelectStmt   // INSERT ... SELECT
	TableSource    *TableStmt    // INSERT ... TABLE table_name (MySQL 8.0.19+)
	SetList        []*Assignment // INSERT ... SET col=val (MySQL extension)
	RowAlias       string        // AS row_alias (MySQL 8.0.19+)
	ColAliases     []string      // (col_alias, ...) after row_alias
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
	Ignore      bool        // IGNORE before SELECT in CTAS
	Replace     bool        // REPLACE before SELECT in CTAS
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
	ATCoalescePartition
	ATReorganizePartition
	ATExchangePartition
	ATTruncatePartition
	ATAnalyzePartition
	ATCheckPartition
	ATOptimizePartition
	ATRebuildPartition
	ATRepairPartition
	ATDiscardPartitionTablespace
	ATImportPartitionTablespace
	ATRemovePartitioning
	ATAlterColumnVisible
	ATAlterColumnInvisible
	ATAlterIndexVisible
	ATAlterIndexInvisible
	ATForce
	ATOrderBy
	ATEnableKeys
	ATDisableKeys
	ATDiscardTablespace
	ATImportTablespace
	ATWithValidation
	ATWithoutValidation
	ATAlterCheckEnforced
	ATSecondaryLoad
	ATSecondaryUnload
)

// AlterTableCmd represents a single ALTER TABLE operation.
type AlterTableCmd struct {
	Loc            Loc
	Type           AlterTableCmdType
	Name           string       // column/constraint/index name
	NewName        string       // for RENAME operations
	Column         *ColumnDef   // for ADD/MODIFY/CHANGE COLUMN
	Columns        []*ColumnDef // for ADD (col1, col2, ...) multi-column form
	Constraint     *Constraint  // for ADD CONSTRAINT
	Option         *TableOption // for table options
	DefaultExpr    ExprNode     // for ALTER COLUMN SET DEFAULT expr
	IfExists       bool
	First          bool            // FIRST positioning
	After          string          // AFTER column positioning
	PartitionNames []string        // for partition operations (DROP/TRUNCATE/ANALYZE/CHECK/OPTIMIZE/REBUILD/REPAIR/REORGANIZE/DISCARD/IMPORT)
	AllPartitions  bool            // for partition operations using ALL
	PartitionDefs  []*PartitionDef // for ADD PARTITION / REORGANIZE PARTITION ... INTO
	Number         int             // for COALESCE PARTITION number
	ExchangeTable  *TableRef       // for EXCHANGE PARTITION ... WITH TABLE
	WithValidation *bool           // for EXCHANGE PARTITION: true=WITH VALIDATION, false=WITHOUT VALIDATION, nil=not specified
	OrderByItems   []*OrderByItem  // for ORDER BY operation
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
	SelectText  string // raw SQL text of the AS select_statement portion
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
	Loc                      Loc
	Name                     string
	TypeName                 *DataType
	Constraints              []*ColumnConstraint
	DefaultValue             ExprNode
	Comment                  string
	AutoIncrement            bool
	OnUpdate                 ExprNode // ON UPDATE CURRENT_TIMESTAMP
	Generated                *GeneratedColumn
	ColumnFormat             string // COLUMN_FORMAT {FIXED | DYNAMIC | DEFAULT}
	Storage                  string // STORAGE {DISK | MEMORY}
	EngineAttribute          string // ENGINE_ATTRIBUTE [=] 'string'
	SecondaryEngineAttribute string // SECONDARY_ENGINE_ATTRIBUTE [=] 'string'
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
	ColConstrVisible
	ColConstrInvisible
	ColConstrColumnFormat
	ColConstrStorage
	ColConstrEngineAttribute
	ColConstrSecondaryEngineAttribute
)

// ColumnConstraint represents a column-level constraint.
type ColumnConstraint struct {
	Loc         Loc
	Type        ColumnConstraintType
	Name        string   // constraint name
	Expr        ExprNode // for CHECK, DEFAULT, etc.
	RefTable    *TableRef
	RefColumns  []string
	Match       string // FULL, PARTIAL, SIMPLE
	OnDelete    ReferenceAction
	OnUpdate    ReferenceAction
	NotEnforced bool // for CHECK ... NOT ENFORCED
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
	Loc          Loc
	Type         ConstraintType
	Name         string
	Columns      []string       // simple column names
	IndexColumns []*IndexColumn // key parts with optional expressions (for functional indexes)
	IndexType    string         // BTREE, HASH
	IndexOptions []*IndexOption // index_option list (KEY_BLOCK_SIZE, COMMENT, VISIBLE/INVISIBLE, etc.)
	Expr         ExprNode       // for CHECK
	RefTable     *TableRef      // for FOREIGN KEY
	RefColumns   []string       // for FOREIGN KEY
	OnDelete     ReferenceAction
	OnUpdate     ReferenceAction
	Match        string // FULL, PARTIAL, SIMPLE
	NotEnforced  bool   // for CHECK ... NOT ENFORCED
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
	Loc     Loc
	Name    string
	Value   string
	Storage string // TABLESPACE ... STORAGE {DISK|MEMORY}
}

func (o *TableOption) nodeTag() {}

// PartitionClause represents a PARTITION BY clause.
type PartitionClause struct {
	Loc            Loc
	Type           PartitionType
	Linear         bool // LINEAR HASH or LINEAR KEY
	Expr           ExprNode
	Columns        []string
	Algorithm      int // ALGORITHM={1|2} for KEY partitioning (0 if not specified)
	NumParts       int
	Partitions     []*PartitionDef
	SubPartType    PartitionType // SUBPARTITION BY type (0 if no subpartitioning)
	SubPartExpr    ExprNode      // SUBPARTITION BY ... (expr)
	SubPartColumns []string      // SUBPARTITION BY KEY (columns)
	SubPartAlgo    int           // ALGORITHM={1|2} for SUBPARTITION BY KEY (0 if not specified)
	NumSubParts    int           // SUBPARTITIONS num
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
	Loc           Loc
	Name          string
	Values        Node // VALUES LESS THAN (expr) or IN (list) — kept as Node to support heterogeneous partition value types
	Options       []*TableOption
	SubPartitions []*SubPartitionDef
}

// SubPartitionDef represents a subpartition definition.
type SubPartitionDef struct {
	Loc     Loc
	Name    string
	Options []*TableOption
}

func (d *SubPartitionDef) nodeTag() {}

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
	SRID       int      // Spatial Reference ID for spatial types (0 = not set)
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
	BinOpXor // logical XOR
	BinOpBitAnd
	BinOpBitOr
	BinOpBitXor
	BinOpShiftLeft
	BinOpShiftRight
	BinOpDivInt // DIV
	BinOpRegexp
	BinOpLikeEscape
	BinOpNullSafeEq  // <=>
	BinOpAssign      // :=
	BinOpJsonExtract // ->
	BinOpJsonUnquote // ->>
	BinOpSoundsLike  // SOUNDS LIKE
)

// BinaryExpr represents a binary expression (a op b).
type BinaryExpr struct {
	Loc        Loc
	Op         BinaryOp
	Left       ExprNode
	Right      ExprNode
	OriginalOp string // Original operator text (e.g., "MOD", "!=") when different from canonical form; used for auto-alias generation
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
	UnaryBinary // BINARY expr (cast to binary string)
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
	Loc       Loc
	Name      string
	Schema    string // schema-qualified name
	Args      []ExprNode
	Distinct  bool // COUNT(DISTINCT ...)
	Star      bool // COUNT(*)
	HasParens bool // true if explicit parentheses were written (e.g., NOW() vs CURRENT_TIMESTAMP)
	OrderBy   []*OrderByItem
	Over      *WindowDef // OVER clause for window functions
	Separator ExprNode   // GROUP_CONCAT SEPARATOR value (nil = not specified)
}

func (e *FuncCallExpr) nodeTag()  {}
func (e *FuncCallExpr) exprNode() {}

// SubqueryExpr represents a subquery expression.
type SubqueryExpr struct {
	Loc     Loc
	Select  *SelectStmt
	Exists  bool   // EXISTS (SELECT ...)
	Lateral bool   // LATERAL derived table
	Alias   string // AS alias (for derived tables)
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
	Schema     string   // database name
	Name       string   // table name
	Alias      string   // AS alias
	Partitions []string // PARTITION (p0, p1, ...)
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

// TemporalLit represents a temporal literal (DATE '...', TIME '...', TIMESTAMP '...').
type TemporalLit struct {
	Loc   Loc
	Type  string // "DATE", "TIME", "TIMESTAMP"
	Value string
}

func (l *TemporalLit) nodeTag()  {}
func (l *TemporalLit) exprNode() {}

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
	JoinStraight     // STRAIGHT_JOIN (MySQL-specific)
	JoinNaturalLeft  // NATURAL LEFT [OUTER] JOIN
	JoinNaturalRight // NATURAL RIGHT [OUTER] JOIN
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
	Loc                Loc
	Outfile            string
	Dumpfile           string
	Vars               []*VariableRef
	Charset            string // CHARACTER SET charset_name
	FieldsTerminatedBy string // FIELDS TERMINATED BY 'string'
	FieldsEnclosedBy   string // [OPTIONALLY] ENCLOSED BY 'char'
	FieldsOptionalEncl bool   // OPTIONALLY keyword before ENCLOSED BY
	FieldsEscapedBy    string // ESCAPED BY 'char'
	LinesStartingBy    string // LINES STARTING BY 'string'
	LinesTerminatedBy  string // LINES TERMINATED BY 'string'
	HasFieldsClause    bool   // true if FIELDS/COLUMNS clause present
	HasLinesClause     bool   // true if LINES clause present
}

func (c *IntoClause) nodeTag() {}

// ForUpdate represents FOR UPDATE / FOR SHARE.
type ForUpdate struct {
	Loc             Loc
	Share           bool // FOR SHARE instead of FOR UPDATE
	LockInShareMode bool // LOCK IN SHARE MODE (legacy syntax)
	Tables          []*TableRef
	NoWait          bool
	SkipLocked      bool
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

// SetPasswordStmt represents a SET PASSWORD statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/set-password.html
//
//	SET PASSWORD [FOR user] = 'auth_string'
type SetPasswordStmt struct {
	Loc                   Loc
	User                  *UserSpec // optional FOR user@host
	Password              string    // the password string
	ToRandom              bool      // SET PASSWORD ... TO RANDOM
	Replace               string    // [REPLACE 'current_auth_string']
	RetainCurrentPassword bool      // [RETAIN CURRENT PASSWORD]
}

func (s *SetPasswordStmt) nodeTag()  {}
func (s *SetPasswordStmt) stmtNode() {}

// ShowStmt represents a SHOW statement.
type ShowStmt struct {
	Loc          Loc
	Type         string // DATABASES, TABLES, COLUMNS, etc.
	From         *TableRef
	Like         ExprNode
	Where        ExprNode
	IfNotExists  bool        // SHOW CREATE DATABASE IF NOT EXISTS
	ProfileTypes []string    // SHOW PROFILE type list (CPU, BLOCK IO, etc.)
	ForQuery     ExprNode    // SHOW PROFILE FOR QUERY n
	ForUser      *UserSpec   // SHOW GRANTS FOR user
	Using        []*UserSpec // SHOW GRANTS ... USING role [, role]
	FromPos      ExprNode    // SHOW BINLOG/RELAYLOG EVENTS FROM pos
	LimitCount   ExprNode    // SHOW ... LIMIT count
	LimitOffset  ExprNode    // SHOW ... LIMIT offset, count
	Channel      string      // FOR CHANNEL channel (RELAYLOG EVENTS, REPLICA STATUS)
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
	Loc           Loc
	Format        string // TRADITIONAL, JSON, TREE
	Analyze       bool   // EXPLAIN ANALYZE
	Extended      bool   // EXPLAIN EXTENDED (deprecated but still parsed)
	Partitions    bool   // EXPLAIN PARTITIONS (deprecated but still parsed)
	ForConnection int64  // EXPLAIN FOR CONNECTION connection_id
	Stmt          StmtNode
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

// SavepointStmt represents a SAVEPOINT or RELEASE SAVEPOINT statement.
type SavepointStmt struct {
	Loc     Loc
	Name    string
	Release bool // true for RELEASE SAVEPOINT
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

// RequireClause represents a REQUIRE clause for TLS options.
type RequireClause struct {
	Loc     Loc
	None    bool   // REQUIRE NONE
	SSL     bool   // REQUIRE SSL
	X509    bool   // REQUIRE X509
	Cipher  string // REQUIRE CIPHER 'cipher'
	Issuer  string // REQUIRE ISSUER 'issuer'
	Subject string // REQUIRE SUBJECT 'subject'
}

func (r *RequireClause) nodeTag() {}

// ResourceOption represents a WITH resource limit option.
type ResourceOption struct {
	Loc                   Loc
	MaxQueriesPerHour     int
	MaxUpdatesPerHour     int
	MaxConnectionsPerHour int
	MaxUserConnections    int
	HasMaxQueries         bool
	HasMaxUpdates         bool
	HasMaxConnections     bool
	HasMaxUserConn        bool
}

func (r *ResourceOption) nodeTag() {}

// GrantStmt represents a GRANT statement.
type GrantStmt struct {
	Loc          Loc
	Privileges   []string
	AllPriv      bool
	ProxyUser    string // GRANT PROXY ON user
	On           *GrantTarget
	To           []string
	WithGrant    bool
	AsUser       string          // AS user (MySQL 8.0+ role context)
	WithRoleType string          // DEFAULT, NONE, ALL, ALL EXCEPT, or "" (roles list)
	WithRoles    []string        // role names for WITH ROLE role, role, ...
	Require      *RequireClause  // REQUIRE tls_option
	Resource     *ResourceOption // WITH resource_option
}

func (s *GrantStmt) nodeTag()  {}
func (s *GrantStmt) stmtNode() {}

// RevokeStmt represents a REVOKE statement.
type RevokeStmt struct {
	Loc               Loc
	IfExists          bool
	Privileges        []string
	AllPriv           bool
	GrantOption       bool // REVOKE ALL PRIVILEGES, GRANT OPTION FROM user
	On                *GrantTarget
	From              []string
	IgnoreUnknownUser bool // [IGNORE UNKNOWN USER]
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
	Loc          Loc
	IfNotExists  bool
	Users        []*UserSpec
	DefaultRoles []string        // DEFAULT ROLE role [, role] ...
	Require      *RequireClause  // REQUIRE tls_option
	Resource     *ResourceOption // WITH resource_option
	// Password options
	PasswordExpire         string // "", "DEFAULT", "NEVER", "INTERVAL N DAY"
	PasswordHistory        string // "", "DEFAULT", "N"
	PasswordReuseInterval  string // "", "DEFAULT", "N DAY"
	PasswordRequireCurrent string // "", "DEFAULT", "OPTIONAL"
	FailedLoginAttempts    int
	HasFailedLogin         bool
	PasswordLockTime       string // "", "N", "UNBOUNDED"
	// Account lock
	AccountLock   bool
	AccountUnlock bool
	// Comment / Attribute
	Comment   string
	Attribute string
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
	Loc             Loc
	IfExists        bool
	IsUserFunc      bool // ALTER USER USER() ...
	Users           []*UserSpec
	DefaultRoleUser string // ALTER USER user DEFAULT ROLE ...
	DefaultRoleType string // NONE, ALL, or "" (specific roles)
	DefaultRoles    []string
	Require         *RequireClause  // REQUIRE tls_option
	Resource        *ResourceOption // WITH resource_option
	// Password options
	PasswordExpire         string // "", "DEFAULT", "NEVER", "INTERVAL N DAY"
	PasswordHistory        string // "", "DEFAULT", "N"
	PasswordReuseInterval  string // "", "DEFAULT", "N DAY"
	PasswordRequireCurrent string // "", "DEFAULT", "OPTIONAL"
	FailedLoginAttempts    int
	HasFailedLogin         bool
	PasswordLockTime       string // "", "N", "UNBOUNDED"
	// Account lock
	AccountLock   bool
	AccountUnlock bool
	// Comment / Attribute
	Comment   string
	Attribute string
	// Factor operations (ADD/MODIFY/DROP factor)
	FactorOps []*FactorOp
	// Registration operation
	RegistrationOp *RegistrationOp
}

// FactorOp represents an ADD/MODIFY/DROP factor operation in ALTER USER.
//
//	ADD factor factor_auth_option
//	MODIFY factor factor_auth_option
//	DROP factor
type FactorOp struct {
	Loc            Loc
	Action         string // "ADD", "MODIFY", "DROP"
	Factor         int    // 2 or 3
	AuthPlugin     string // IDENTIFIED WITH auth_plugin
	Password       string // IDENTIFIED BY 'auth_string'
	AuthHash       string // IDENTIFIED WITH plugin AS 'auth_string'
	PasswordRandom bool   // IDENTIFIED BY RANDOM PASSWORD
}

func (f *FactorOp) nodeTag() {}

// RegistrationOp represents a registration_option in ALTER USER.
//
//	factor INITIATE REGISTRATION
//	factor FINISH REGISTRATION SET CHALLENGE_RESPONSE AS 'auth_string'
//	factor UNREGISTER
type RegistrationOp struct {
	Loc               Loc
	Factor            int    // 2 or 3
	Action            string // "INITIATE", "FINISH", "UNREGISTER"
	ChallengeResponse string // for FINISH REGISTRATION
}

func (r *RegistrationOp) nodeTag() {}

func (s *AlterUserStmt) nodeTag()  {}
func (s *AlterUserStmt) stmtNode() {}

// UserSpec represents a user specification.
type UserSpec struct {
	Loc                   Loc
	Name                  string
	Host                  string
	AuthPlugin            string
	Password              string
	AuthHash              string      // IDENTIFIED WITH plugin AS 'hash_string'
	PasswordRandom        bool        // IDENTIFIED BY RANDOM PASSWORD
	RetainCurrentPassword bool        // RETAIN CURRENT PASSWORD
	DiscardOldPassword    bool        // DISCARD OLD PASSWORD
	Replace               string      // [REPLACE 'current_auth_string'] (ALTER USER)
	AuthFactors           []*UserSpec // AND 2fa/3fa auth chaining
	InitialAuthPlugin     string      // INITIAL AUTHENTICATION IDENTIFIED WITH plugin
	InitialAuthString     string      // INITIAL AUTHENTICATION ... AS/BY 'auth_string'
	InitialAuthRandom     bool        // INITIAL AUTHENTICATION IDENTIFIED BY RANDOM PASSWORD
}

func (u *UserSpec) nodeTag() {}

// CreateRoleStmt represents a CREATE ROLE statement.
type CreateRoleStmt struct {
	Loc         Loc
	IfNotExists bool
	Roles       []string
}

func (s *CreateRoleStmt) nodeTag()  {}
func (s *CreateRoleStmt) stmtNode() {}

// DropRoleStmt represents a DROP ROLE statement.
type DropRoleStmt struct {
	Loc      Loc
	IfExists bool
	Roles    []string
}

func (s *DropRoleStmt) nodeTag()  {}
func (s *DropRoleStmt) stmtNode() {}

// SetDefaultRoleStmt represents a SET DEFAULT ROLE statement.
type SetDefaultRoleStmt struct {
	Loc   Loc
	None  bool     // NONE
	All   bool     // ALL
	Roles []string // specific roles
	To    []string // target users
}

func (s *SetDefaultRoleStmt) nodeTag()  {}
func (s *SetDefaultRoleStmt) stmtNode() {}

// SetRoleStmt represents a SET ROLE statement.
type SetRoleStmt struct {
	Loc       Loc
	Default   bool     // DEFAULT
	None      bool     // NONE
	All       bool     // ALL
	AllExcept []string // ALL EXCEPT roles
	Roles     []string // specific roles
}

func (s *SetRoleStmt) nodeTag()  {}
func (s *SetRoleStmt) stmtNode() {}

// GrantRoleStmt represents a GRANT role TO user statement.
type GrantRoleStmt struct {
	Loc       Loc
	Roles     []string
	To        []string
	WithAdmin bool
}

func (s *GrantRoleStmt) nodeTag()  {}
func (s *GrantRoleStmt) stmtNode() {}

// RevokeRoleStmt represents a REVOKE role FROM user statement.
type RevokeRoleStmt struct {
	Loc               Loc
	IfExists          bool
	Roles             []string
	From              []string
	IgnoreUnknownUser bool
}

func (s *RevokeRoleStmt) nodeTag()  {}
func (s *RevokeRoleStmt) stmtNode() {}

// CreateFunctionStmt represents a CREATE FUNCTION or CREATE PROCEDURE.
type CreateFunctionStmt struct {
	Loc             Loc
	OrReplace       bool
	IfNotExists     bool
	IsProcedure     bool
	IsAggregate     bool // CREATE AGGREGATE FUNCTION (loadable UDF)
	Definer         string
	Name            *TableRef
	Params          []*FuncParam
	Returns         *DataType
	Soname          string // SONAME 'shared_library' (loadable UDF)
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
	Loc         Loc
	Definer     string
	IfNotExists bool
	Name        string
	Timing      string // BEFORE, AFTER
	Event       string // INSERT, UPDATE, DELETE
	Table       *TableRef
	Order       *TriggerOrder
	Body        string
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
	Loc     Loc
	At      ExprNode
	Every   ExprNode
	Starts  ExprNode
	Ends    ExprNode
	RawText string // raw schedule text from input (e.g., "EVERY 1 HOUR" or "AT CURRENT_TIMESTAMP + INTERVAL 1 HOUR")
}

func (s *EventSchedule) nodeTag() {}

// LoadDataStmt represents a LOAD DATA statement.
type LoadDataStmt struct {
	Loc                Loc
	IsXML              bool // true for LOAD XML, false for LOAD DATA
	LowPriority        bool // LOW_PRIORITY modifier
	Concurrent         bool // CONCURRENT modifier
	Local              bool
	Infile             string
	Replace            bool
	Ignore             bool
	Table              *TableRef
	Partitions         []string // PARTITION (p0, p1, ...)
	CharacterSet       string   // CHARACTER SET charset_name
	RowsIdentifiedBy   string   // ROWS IDENTIFIED BY '<tagname>' (XML only)
	Columns            []*ColumnRef
	SetList            []*Assignment
	LinesStartingBy    string
	LinesTerminatedBy  string
	FieldsTerminatedBy string
	FieldsEnclosedBy   string
	FieldsOptionalEncl bool // OPTIONALLY ENCLOSED BY
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
	Loc              Loc
	NoWriteToBinlog  bool
	Tables           []*TableRef
	HistogramOp      string   // "UPDATE" or "DROP"
	HistogramColumns []string // column names
	Buckets          int      // WITH N BUCKETS
	UsingData        string   // USING DATA 'json_data'
}

func (s *AnalyzeTableStmt) nodeTag()  {}
func (s *AnalyzeTableStmt) stmtNode() {}

// OptimizeTableStmt represents an OPTIMIZE TABLE statement.
type OptimizeTableStmt struct {
	Loc             Loc
	NoWriteToBinlog bool
	Tables          []*TableRef
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
	Loc             Loc
	NoWriteToBinlog bool
	Tables          []*TableRef
	Quick           bool
	Extended        bool
	UseFrm          bool
}

func (s *RepairTableStmt) nodeTag()  {}
func (s *RepairTableStmt) stmtNode() {}

// FlushStmt represents a FLUSH statement.
type FlushStmt struct {
	Loc             Loc
	NoWriteToBinlog bool
	Options         []string
	Tables          []*TableRef // FLUSH TABLES tbl_list
	WithReadLock    bool        // WITH READ LOCK
	ForExport       bool        // FOR EXPORT
	RelayChannel    string      // RELAY LOGS FOR CHANNEL channel
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

// ChecksumTableStmt represents a CHECKSUM TABLE statement.
type ChecksumTableStmt struct {
	Loc      Loc
	Tables   []*TableRef
	Quick    bool
	Extended bool
}

func (s *ChecksumTableStmt) nodeTag()  {}
func (s *ChecksumTableStmt) stmtNode() {}

// ShutdownStmt represents a SHUTDOWN statement.
type ShutdownStmt struct {
	Loc Loc
}

func (s *ShutdownStmt) nodeTag()  {}
func (s *ShutdownStmt) stmtNode() {}

// RestartStmt represents a RESTART statement.
type RestartStmt struct {
	Loc Loc
}

func (s *RestartStmt) nodeTag()  {}
func (s *RestartStmt) stmtNode() {}

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

// DefaultExpr represents the DEFAULT keyword or DEFAULT(col_name) function used as an expression.
type DefaultExpr struct {
	Loc    Loc
	Column string // column name for DEFAULT(col_name) form; empty for bare DEFAULT
}

func (e *DefaultExpr) nodeTag()  {}
func (e *DefaultExpr) exprNode() {}

// RowExpr represents a ROW(expr, expr, ...) constructor expression.
type RowExpr struct {
	Loc   Loc
	Items []ExprNode
}

func (e *RowExpr) nodeTag()  {}
func (e *RowExpr) exprNode() {}

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
	Loc         Loc
	Name        string
	TypeName    *DataType
	Path        string
	Exists      bool   // FOR ORDINALITY or EXISTS PATH
	ErrorOption string // ERROR ON ERROR, NULL ON ERROR, DEFAULT val ON ERROR
	EmptyOption string // ERROR ON EMPTY, NULL ON EMPTY, DEFAULT val ON EMPTY
	Ordinality  bool   // FOR ORDINALITY
	Nested      bool   // NESTED PATH
	NestedPath  string
	NestedCols  []*JsonTableColumn
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
	Target Node   // target: *VariableRef or *ColumnRef (local variable)
	Name   string // item name (e.g. NUMBER, ROW_COUNT, MESSAGE_TEXT, etc.)
}

func (d *DiagnosticsItem) nodeTag() {}

// BeginEndBlock represents a BEGIN...END compound statement block.
type BeginEndBlock struct {
	Loc      Loc
	Label    string // optional begin label name
	EndLabel string // optional end label name
	Stmts    []Node // statements inside the block
}

func (s *BeginEndBlock) nodeTag()  {}
func (s *BeginEndBlock) stmtNode() {}

// DeclareVarStmt represents a DECLARE variable statement inside a compound block.
type DeclareVarStmt struct {
	Loc      Loc
	Names    []string  // variable names
	TypeName *DataType // data type
	Default  ExprNode  // DEFAULT value (may be nil)
}

func (s *DeclareVarStmt) nodeTag()  {}
func (s *DeclareVarStmt) stmtNode() {}

// DeclareConditionStmt represents a DECLARE ... CONDITION FOR statement.
type DeclareConditionStmt struct {
	Loc            Loc
	Name           string // condition name
	ConditionValue string // SQLSTATE value or mysql_error_code
}

func (s *DeclareConditionStmt) nodeTag()  {}
func (s *DeclareConditionStmt) stmtNode() {}

// DeclareHandlerStmt represents a DECLARE ... HANDLER FOR statement.
type DeclareHandlerStmt struct {
	Loc        Loc
	Action     string   // CONTINUE, EXIT, or UNDO
	Conditions []string // condition values
	Stmt       Node     // handler statement body
}

func (s *DeclareHandlerStmt) nodeTag()  {}
func (s *DeclareHandlerStmt) stmtNode() {}

// DeclareCursorStmt represents a DECLARE ... CURSOR FOR statement.
type DeclareCursorStmt struct {
	Loc    Loc
	Name   string      // cursor name
	Select *SelectStmt // select statement
}

func (s *DeclareCursorStmt) nodeTag()  {}
func (s *DeclareCursorStmt) stmtNode() {}

// TableStmt represents a TABLE statement (MySQL 8.0.19+).
// TABLE t is equivalent to SELECT * FROM t.
type TableStmt struct {
	Loc     Loc
	Table   *TableRef      // table name
	OrderBy []*OrderByItem // optional ORDER BY
	Limit   *Limit         // optional LIMIT/OFFSET
	Into    *IntoClause    // optional INTO OUTFILE/DUMPFILE/var
}

func (s *TableStmt) nodeTag()  {}
func (s *TableStmt) stmtNode() {}

// ValuesStmt represents a VALUES statement (MySQL 8.0.19+).
// VALUES ROW(1,2), ROW(3,4) returns rows of literal values.
type ValuesStmt struct {
	Loc     Loc
	Rows    [][]ExprNode   // list of ROW(...) value lists
	OrderBy []*OrderByItem // optional ORDER BY
	Limit   *Limit         // optional LIMIT/OFFSET
}

func (s *ValuesStmt) nodeTag()  {}
func (s *ValuesStmt) stmtNode() {}

// IfStmt represents an IF/ELSEIF/ELSE compound statement.
type IfStmt struct {
	Loc      Loc
	Cond     ExprNode  // IF condition
	ThenList []Node    // THEN statement list
	ElseIfs  []*ElseIf // ELSEIF clauses
	ElseList []Node    // ELSE statement list (may be nil)
}

func (s *IfStmt) nodeTag()  {}
func (s *IfStmt) stmtNode() {}

// ElseIf represents an ELSEIF clause within an IF statement.
type ElseIf struct {
	Loc      Loc
	Cond     ExprNode // ELSEIF condition
	ThenList []Node   // THEN statement list
}

func (e *ElseIf) nodeTag() {}

// CaseStmtNode represents a CASE statement (not expression) in stored programs.
type CaseStmtNode struct {
	Loc      Loc
	Operand  ExprNode        // simple CASE operand (nil for searched CASE)
	Whens    []*CaseStmtWhen // WHEN clauses
	ElseList []Node          // ELSE statement list (may be nil)
}

func (s *CaseStmtNode) nodeTag()  {}
func (s *CaseStmtNode) stmtNode() {}

// CaseStmtWhen represents a WHEN clause in a CASE statement.
type CaseStmtWhen struct {
	Loc      Loc
	Cond     ExprNode // WHEN condition/value
	ThenList []Node   // THEN statement list
}

func (w *CaseStmtWhen) nodeTag() {}

// WhileStmt represents a WHILE loop compound statement.
type WhileStmt struct {
	Loc      Loc
	Label    string   // optional begin label
	EndLabel string   // optional end label
	Cond     ExprNode // loop condition
	Stmts    []Node   // statement list
}

func (s *WhileStmt) nodeTag()  {}
func (s *WhileStmt) stmtNode() {}

// RepeatStmt represents a REPEAT loop compound statement.
type RepeatStmt struct {
	Loc      Loc
	Label    string   // optional begin label
	EndLabel string   // optional end label
	Stmts    []Node   // statement list
	Cond     ExprNode // UNTIL condition
}

func (s *RepeatStmt) nodeTag()  {}
func (s *RepeatStmt) stmtNode() {}

// LoopStmt represents a LOOP compound statement.
type LoopStmt struct {
	Loc      Loc
	Label    string // optional begin label
	EndLabel string // optional end label
	Stmts    []Node // statement list
}

func (s *LoopStmt) nodeTag()  {}
func (s *LoopStmt) stmtNode() {}

// LeaveStmt represents a LEAVE statement.
type LeaveStmt struct {
	Loc   Loc
	Label string
}

func (s *LeaveStmt) nodeTag()  {}
func (s *LeaveStmt) stmtNode() {}

// IterateStmt represents an ITERATE statement.
type IterateStmt struct {
	Loc   Loc
	Label string
}

func (s *IterateStmt) nodeTag()  {}
func (s *IterateStmt) stmtNode() {}

// ReturnStmt represents a RETURN statement in stored functions.
type ReturnStmt struct {
	Loc  Loc
	Expr ExprNode
}

func (s *ReturnStmt) nodeTag()  {}
func (s *ReturnStmt) stmtNode() {}

// OpenCursorStmt represents an OPEN cursor_name statement.
type OpenCursorStmt struct {
	Loc  Loc
	Name string
}

func (s *OpenCursorStmt) nodeTag()  {}
func (s *OpenCursorStmt) stmtNode() {}

// FetchCursorStmt represents a FETCH cursor_name INTO var_name [, var_name] ... statement.
type FetchCursorStmt struct {
	Loc  Loc
	Name string   // cursor name
	Into []string // variable names
}

func (s *FetchCursorStmt) nodeTag()  {}
func (s *FetchCursorStmt) stmtNode() {}

// CloseCursorStmt represents a CLOSE cursor_name statement.
type CloseCursorStmt struct {
	Loc  Loc
	Name string
}

func (s *CloseCursorStmt) nodeTag()  {}
func (s *CloseCursorStmt) stmtNode() {}

// CloneStmt represents a CLONE statement (MySQL 8.0.17+).
type CloneStmt struct {
	Loc        Loc
	Local      bool   // true = CLONE LOCAL, false = CLONE INSTANCE FROM
	Directory  string // DATA DIRECTORY path
	User       string // remote user (CLONE INSTANCE only)
	Host       string // remote host (CLONE INSTANCE only)
	Port       int64  // remote port (CLONE INSTANCE only)
	Password   string // IDENTIFIED BY password (CLONE INSTANCE only)
	RequireSSL *bool  // nil = not specified, true = REQUIRE SSL, false = REQUIRE NO SSL
}

func (s *CloneStmt) nodeTag()  {}
func (s *CloneStmt) stmtNode() {}

// InstallPluginStmt represents an INSTALL PLUGIN statement.
type InstallPluginStmt struct {
	Loc        Loc
	PluginName string // plugin name
	Soname     string // shared library name
}

func (s *InstallPluginStmt) nodeTag()  {}
func (s *InstallPluginStmt) stmtNode() {}

// UninstallPluginStmt represents an UNINSTALL PLUGIN statement.
type UninstallPluginStmt struct {
	Loc        Loc
	PluginName string // plugin name
}

func (s *UninstallPluginStmt) nodeTag()  {}
func (s *UninstallPluginStmt) stmtNode() {}

// InstallComponentStmt represents an INSTALL COMPONENT statement.
type InstallComponentStmt struct {
	Loc        Loc
	Components []string      // component names (string literals)
	SetVars    []*Assignment // optional SET variable = expr list
}

func (s *InstallComponentStmt) nodeTag()  {}
func (s *InstallComponentStmt) stmtNode() {}

// UninstallComponentStmt represents an UNINSTALL COMPONENT statement.
type UninstallComponentStmt struct {
	Loc        Loc
	Components []string // component names (string literals)
}

func (s *UninstallComponentStmt) nodeTag()  {}
func (s *UninstallComponentStmt) stmtNode() {}

// CreateTablespaceStmt represents a CREATE TABLESPACE statement.
type CreateTablespaceStmt struct {
	Loc             Loc
	Undo            bool   // UNDO tablespace
	Name            string // tablespace name
	DataFile        string // ADD DATAFILE 'file_name'
	AutoextendSize  string // AUTOEXTEND_SIZE value
	FileBlockSize   string // FILE_BLOCK_SIZE value
	Encryption      string // ENCRYPTION value ('Y' or 'N')
	UseLogfileGroup string // USE LOGFILE GROUP name
	ExtentSize      string // EXTENT_SIZE value
	InitialSize     string // INITIAL_SIZE value
	MaxSize         string // MAX_SIZE value
	NodeGroup       string // NODEGROUP value
	Wait            bool   // WAIT
	Comment         string // COMMENT value
	Engine          string // ENGINE name
	EngineAttribute string // ENGINE_ATTRIBUTE value
}

func (s *CreateTablespaceStmt) nodeTag()  {}
func (s *CreateTablespaceStmt) stmtNode() {}

// AlterTablespaceStmt represents an ALTER TABLESPACE statement.
type AlterTablespaceStmt struct {
	Loc             Loc
	Undo            bool   // UNDO tablespace
	Name            string // tablespace name
	AddDataFile     string // ADD DATAFILE 'file_name'
	DropDataFile    string // DROP DATAFILE 'file_name'
	InitialSize     string // INITIAL_SIZE value
	Wait            bool   // WAIT
	RenameTo        string // RENAME TO tablespace_name
	AutoextendSize  string // AUTOEXTEND_SIZE value
	Encryption      string // ENCRYPTION value
	Engine          string // ENGINE name
	EngineAttribute string // ENGINE_ATTRIBUTE value
	SetActive       bool   // SET ACTIVE
	SetInactive     bool   // SET INACTIVE
}

func (s *AlterTablespaceStmt) nodeTag()  {}
func (s *AlterTablespaceStmt) stmtNode() {}

// DropTablespaceStmt represents a DROP TABLESPACE statement.
type DropTablespaceStmt struct {
	Loc    Loc
	Undo   bool   // UNDO tablespace
	Name   string // tablespace name
	Engine string // ENGINE name
}

func (s *DropTablespaceStmt) nodeTag()  {}
func (s *DropTablespaceStmt) stmtNode() {}

// CreateServerStmt represents a CREATE SERVER statement.
type CreateServerStmt struct {
	Loc         Loc
	Name        string   // server name
	WrapperName string   // FOREIGN DATA WRAPPER name
	Options     []string // OPTIONS values
}

func (s *CreateServerStmt) nodeTag()  {}
func (s *CreateServerStmt) stmtNode() {}

// AlterServerStmt represents an ALTER SERVER statement.
type AlterServerStmt struct {
	Loc     Loc
	Name    string   // server name
	Options []string // OPTIONS values
}

func (s *AlterServerStmt) nodeTag()  {}
func (s *AlterServerStmt) stmtNode() {}

// DropServerStmt represents a DROP SERVER statement.
type DropServerStmt struct {
	Loc      Loc
	IfExists bool
	Name     string // server name
}

func (s *DropServerStmt) nodeTag()  {}
func (s *DropServerStmt) stmtNode() {}

// CreateLogfileGroupStmt represents a CREATE LOGFILE GROUP statement (NDB).
type CreateLogfileGroupStmt struct {
	Loc            Loc
	Name           string // logfile group name
	UndoFile       string // ADD UNDOFILE 'file_name'
	InitialSize    string // INITIAL_SIZE value
	UndoBufferSize string // UNDO_BUFFER_SIZE value
	RedoBufferSize string // REDO_BUFFER_SIZE value
	NodeGroup      string // NODEGROUP value
	Wait           bool   // WAIT
	Comment        string // COMMENT 'string'
	Engine         string // ENGINE name
}

func (s *CreateLogfileGroupStmt) nodeTag()  {}
func (s *CreateLogfileGroupStmt) stmtNode() {}

// AlterLogfileGroupStmt represents an ALTER LOGFILE GROUP statement (NDB).
type AlterLogfileGroupStmt struct {
	Loc         Loc
	Name        string // logfile group name
	UndoFile    string // ADD UNDOFILE 'file_name'
	InitialSize string // INITIAL_SIZE value
	Wait        bool   // WAIT
	Engine      string // ENGINE name
}

func (s *AlterLogfileGroupStmt) nodeTag()  {}
func (s *AlterLogfileGroupStmt) stmtNode() {}

// DropLogfileGroupStmt represents a DROP LOGFILE GROUP statement (NDB).
type DropLogfileGroupStmt struct {
	Loc    Loc
	Name   string // logfile group name
	Engine string // ENGINE name
}

func (s *DropLogfileGroupStmt) nodeTag()  {}
func (s *DropLogfileGroupStmt) stmtNode() {}

// CreateSpatialRefSysStmt represents a CREATE SPATIAL REFERENCE SYSTEM statement.
type CreateSpatialRefSysStmt struct {
	Loc          Loc
	OrReplace    bool   // OR REPLACE
	IfNotExists  bool   // IF NOT EXISTS
	SRID         int64  // spatial reference system ID
	Name         string // NAME 'srs_name'
	Definition   string // DEFINITION 'definition'
	Organization string // ORGANIZATION 'org_name'
	OrgSRID      int64  // IDENTIFIED BY srid (under ORGANIZATION)
	Description  string // DESCRIPTION 'description'
}

func (s *CreateSpatialRefSysStmt) nodeTag()  {}
func (s *CreateSpatialRefSysStmt) stmtNode() {}

// DropSpatialRefSysStmt represents a DROP SPATIAL REFERENCE SYSTEM statement.
type DropSpatialRefSysStmt struct {
	Loc      Loc
	IfExists bool
	SRID     int64
}

func (s *DropSpatialRefSysStmt) nodeTag()  {}
func (s *DropSpatialRefSysStmt) stmtNode() {}

// VCPUSpec represents a VCPU specification in a resource group (e.g., 0-3 or 5).
type VCPUSpec struct {
	Start int64
	End   int64 // -1 means single CPU (no range)
}

func (v *VCPUSpec) nodeTag() {}

// CreateResourceGroupStmt represents a CREATE RESOURCE GROUP statement.
type CreateResourceGroupStmt struct {
	Loc            Loc
	Name           string     // group name
	Type           string     // SYSTEM or USER
	VCPUs          []VCPUSpec // VCPU specs
	ThreadPriority *int64     // THREAD_PRIORITY value
	Enable         *bool      // ENABLE or DISABLE (nil = not specified)
}

func (s *CreateResourceGroupStmt) nodeTag()  {}
func (s *CreateResourceGroupStmt) stmtNode() {}

// AlterResourceGroupStmt represents an ALTER RESOURCE GROUP statement.
type AlterResourceGroupStmt struct {
	Loc            Loc
	Name           string     // group name
	VCPUs          []VCPUSpec // VCPU specs
	ThreadPriority *int64     // THREAD_PRIORITY value
	Enable         *bool      // ENABLE or DISABLE (nil = not specified)
	Force          bool       // FORCE
}

func (s *AlterResourceGroupStmt) nodeTag()  {}
func (s *AlterResourceGroupStmt) stmtNode() {}

// DropResourceGroupStmt represents a DROP RESOURCE GROUP statement.
type DropResourceGroupStmt struct {
	Loc   Loc
	Name  string // group name
	Force bool   // FORCE
}

func (s *DropResourceGroupStmt) nodeTag()  {}
func (s *DropResourceGroupStmt) stmtNode() {}

// AlterViewStmt represents an ALTER VIEW statement.
type AlterViewStmt struct {
	Loc         Loc
	Algorithm   string // UNDEFINED, MERGE, TEMPTABLE
	Definer     string
	SqlSecurity string // DEFINER, INVOKER
	Name        *TableRef
	Columns     []string
	Select      *SelectStmt
	SelectText  string // raw SQL text of the AS select_statement portion
	CheckOption string // CASCADED, LOCAL
}

func (s *AlterViewStmt) nodeTag()  {}
func (s *AlterViewStmt) stmtNode() {}

// AlterEventStmt represents an ALTER EVENT statement.
type AlterEventStmt struct {
	Loc          Loc
	Definer      string
	Name         string
	Schedule     *EventSchedule
	OnCompletion string // PRESERVE, NOT PRESERVE
	RenameTo     string
	Enable       string // ENABLE, DISABLE, DISABLE ON SLAVE
	Comment      string
	Body         string
}

func (s *AlterEventStmt) nodeTag()  {}
func (s *AlterEventStmt) stmtNode() {}

// AlterRoutineStmt represents an ALTER FUNCTION or ALTER PROCEDURE statement.
type AlterRoutineStmt struct {
	Loc             Loc
	IsProcedure     bool
	Name            *TableRef
	Characteristics []*RoutineCharacteristic
}

func (s *AlterRoutineStmt) nodeTag()  {}
func (s *AlterRoutineStmt) stmtNode() {}

// DropRoutineStmt represents a DROP FUNCTION or DROP PROCEDURE statement.
type DropRoutineStmt struct {
	Loc         Loc
	IsProcedure bool
	IfExists    bool
	Name        *TableRef
}

func (s *DropRoutineStmt) nodeTag()  {}
func (s *DropRoutineStmt) stmtNode() {}

// DropTriggerStmt represents a DROP TRIGGER statement.
type DropTriggerStmt struct {
	Loc      Loc
	IfExists bool
	Name     *TableRef
}

func (s *DropTriggerStmt) nodeTag()  {}
func (s *DropTriggerStmt) stmtNode() {}

// DropEventStmt represents a DROP EVENT statement.
type DropEventStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
}

func (s *DropEventStmt) nodeTag()  {}
func (s *DropEventStmt) stmtNode() {}

// ChangeReplicationSourceStmt represents a CHANGE REPLICATION SOURCE TO statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-source-to.html
//
//	CHANGE REPLICATION SOURCE TO option [, option] ... [ channel_option ]
type ChangeReplicationSourceStmt struct {
	Loc     Loc
	Options []*ReplicationOption // SOURCE_HOST = '...', etc.
	Channel string               // FOR CHANNEL channel
	Legacy  bool                 // true for CHANGE MASTER TO (legacy alias)
}

func (s *ChangeReplicationSourceStmt) nodeTag()  {}
func (s *ChangeReplicationSourceStmt) stmtNode() {}

// ReplicationOption represents a single key = value option in CHANGE REPLICATION SOURCE TO.
type ReplicationOption struct {
	Loc   Loc
	Name  string  // e.g. SOURCE_HOST, SOURCE_PORT, SOURCE_LOG_FILE, etc.
	Value string  // the value as string (number, string literal, or identifier)
	IDs   []int64 // for IGNORE_SERVER_IDS = (id, id, ...)
}

func (o *ReplicationOption) nodeTag() {}

// ChangeReplicationFilterStmt represents a CHANGE REPLICATION FILTER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/change-replication-filter.html
//
//	CHANGE REPLICATION FILTER filter[, filter] [FOR CHANNEL channel]
type ChangeReplicationFilterStmt struct {
	Loc     Loc
	Filters []*ReplicationFilter
	Channel string // FOR CHANNEL channel
}

func (s *ChangeReplicationFilterStmt) nodeTag()  {}
func (s *ChangeReplicationFilterStmt) stmtNode() {}

// ReplicationFilter represents a single filter in CHANGE REPLICATION FILTER.
type ReplicationFilter struct {
	Loc    Loc
	Type   string   // REPLICATE_DO_DB, REPLICATE_IGNORE_DB, etc.
	Values []string // database or table names, or patterns
}

func (f *ReplicationFilter) nodeTag() {}

// StartReplicaStmt represents a START REPLICA or START SLAVE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/start-replica.html
//
//	START REPLICA [thread_types] [until_option] [connection_options] [channel_option]
type StartReplicaStmt struct {
	Loc         Loc
	IOThread    bool   // IO_THREAD specified
	SQLThread   bool   // SQL_THREAD specified
	UntilType   string // SQL_BEFORE_GTIDS, SQL_AFTER_GTIDS, SOURCE_LOG_FILE, RELAY_LOG_FILE, SQL_AFTER_MTS_GAPS
	UntilValue  string // gtid_set or log file name
	UntilPos    int64  // log position (for log file modes)
	User        string // USER = 'user'
	Password    string // PASSWORD = 'pass'
	DefaultAuth string // DEFAULT_AUTH = 'plugin'
	PluginDir   string // PLUGIN_DIR = 'dir'
	Channel     string // FOR CHANNEL channel
}

func (s *StartReplicaStmt) nodeTag()  {}
func (s *StartReplicaStmt) stmtNode() {}

// StopReplicaStmt represents a STOP REPLICA or STOP SLAVE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/stop-replica.html
//
//	STOP REPLICA [thread_types] [channel_option]
type StopReplicaStmt struct {
	Loc       Loc
	IOThread  bool   // IO_THREAD specified
	SQLThread bool   // SQL_THREAD specified
	Channel   string // FOR CHANNEL channel
}

func (s *StopReplicaStmt) nodeTag()  {}
func (s *StopReplicaStmt) stmtNode() {}

// ResetReplicaStmt represents a RESET REPLICA or RESET SLAVE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/reset-replica.html
//
//	RESET REPLICA [ALL] [channel_option]
type ResetReplicaStmt struct {
	Loc     Loc
	All     bool   // ALL modifier
	Channel string // FOR CHANNEL channel
}

func (s *ResetReplicaStmt) nodeTag()  {}
func (s *ResetReplicaStmt) stmtNode() {}

// PurgeBinaryLogsStmt represents a PURGE BINARY LOGS statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/purge-binary-logs.html
//
//	PURGE { BINARY | MASTER } LOGS { TO 'log_name' | BEFORE datetime_expr }
type PurgeBinaryLogsStmt struct {
	Loc        Loc
	To         string   // TO 'log_name'
	BeforeExpr ExprNode // BEFORE datetime_expr
}

func (s *PurgeBinaryLogsStmt) nodeTag()  {}
func (s *PurgeBinaryLogsStmt) stmtNode() {}

// ResetMasterStmt represents a RESET MASTER statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/reset-master.html
//
//	RESET MASTER [TO binary_log_file_index_number]
type ResetMasterStmt struct {
	Loc Loc
	To  int64 // TO index number; 0 = not specified
}

func (s *ResetMasterStmt) nodeTag()  {}
func (s *ResetMasterStmt) stmtNode() {}

// StartGroupReplicationStmt represents a START GROUP_REPLICATION statement.
type StartGroupReplicationStmt struct {
	Loc         Loc
	User        string // USER = 'user'
	Password    string // PASSWORD = 'pass'
	DefaultAuth string // DEFAULT_AUTH = 'plugin'
}

func (s *StartGroupReplicationStmt) nodeTag()  {}
func (s *StartGroupReplicationStmt) stmtNode() {}

// StopGroupReplicationStmt represents a STOP GROUP_REPLICATION statement.
type StopGroupReplicationStmt struct {
	Loc Loc
}

func (s *StopGroupReplicationStmt) nodeTag()  {}
func (s *StopGroupReplicationStmt) stmtNode() {}

// AlterInstanceStmt represents an ALTER INSTANCE statement.
type AlterInstanceStmt struct {
	Loc               Loc
	Action            string // ROTATE INNODB MASTER KEY, ROTATE BINLOG MASTER KEY, ENABLE INNODB REDO_LOG, etc.
	Channel           string // FOR CHANNEL value (RELOAD TLS)
	NoRollbackOnError bool
}

func (s *AlterInstanceStmt) nodeTag()  {}
func (s *AlterInstanceStmt) stmtNode() {}

// LockInstanceStmt represents a LOCK INSTANCE FOR BACKUP statement.
type LockInstanceStmt struct {
	Loc Loc
}

func (s *LockInstanceStmt) nodeTag()  {}
func (s *LockInstanceStmt) stmtNode() {}

// UnlockInstanceStmt represents an UNLOCK INSTANCE statement.
type UnlockInstanceStmt struct {
	Loc Loc
}

func (s *UnlockInstanceStmt) nodeTag()  {}
func (s *UnlockInstanceStmt) stmtNode() {}

// ImportTableStmt represents an IMPORT TABLE FROM statement.
type ImportTableStmt struct {
	Loc   Loc
	Files []string // .sdi file paths
}

func (s *ImportTableStmt) nodeTag()  {}
func (s *ImportTableStmt) stmtNode() {}

// BinlogStmt represents a BINLOG 'str' statement.
type BinlogStmt struct {
	Loc Loc
	Str string // base64-encoded string
}

func (s *BinlogStmt) nodeTag()  {}
func (s *BinlogStmt) stmtNode() {}

// CacheIndexStmt represents a CACHE INDEX statement.
type CacheIndexStmt struct {
	Loc        Loc
	Tables     []*TableRef
	Partitions []string // PARTITION (p1, p2, ...) or PARTITION ALL
	Indexes    []string // INDEX (idx1, idx2, ...)
	CacheName  string   // IN cache_name
}

func (s *CacheIndexStmt) nodeTag()  {}
func (s *CacheIndexStmt) stmtNode() {}

// LoadIndexIntoCacheStmt represents a LOAD INDEX INTO CACHE statement.
type LoadIndexIntoCacheStmt struct {
	Loc    Loc
	Tables []*LoadIndexTable
}

// LoadIndexTable represents a table entry in LOAD INDEX INTO CACHE.
type LoadIndexTable struct {
	Loc          Loc
	Table        *TableRef
	Partitions   []string // PARTITION (p0, p1, ...) or ALL
	Indexes      []string // INDEX|KEY (idx1, idx2, ...)
	IgnoreLeaves bool     // IGNORE LEAVES
}

func (n *LoadIndexTable) nodeTag() {}

func (s *LoadIndexIntoCacheStmt) nodeTag()  {}
func (s *LoadIndexIntoCacheStmt) stmtNode() {}

// ResetPersistStmt represents a RESET PERSIST statement.
type ResetPersistStmt struct {
	Loc      Loc
	IfExists bool
	Variable string // variable name, or "" for all
}

func (s *ResetPersistStmt) nodeTag()  {}
func (s *ResetPersistStmt) stmtNode() {}

// RenameUserStmt represents a RENAME USER statement.
type RenameUserStmt struct {
	Loc   Loc
	Pairs []*RenameUserPair
}

func (s *RenameUserStmt) nodeTag()  {}
func (s *RenameUserStmt) stmtNode() {}

// RenameUserPair represents a single old TO new pair in RENAME USER.
type RenameUserPair struct {
	Loc     Loc
	OldUser string
	OldHost string
	NewUser string
	NewHost string
}

func (p *RenameUserPair) nodeTag() {}

// SetResourceGroupStmt represents a SET RESOURCE GROUP statement.
type SetResourceGroupStmt struct {
	Loc       Loc
	Name      string  // group name
	ThreadIDs []int64 // FOR thread_id [, thread_id] ...
}

func (s *SetResourceGroupStmt) nodeTag()  {}
func (s *SetResourceGroupStmt) stmtNode() {}

// HelpStmt represents a HELP 'topic' statement.
type HelpStmt struct {
	Loc   Loc
	Topic string
}

func (s *HelpStmt) nodeTag()  {}
func (s *HelpStmt) stmtNode() {}

// RawStmt wraps a statement node with its location span.
type RawStmt struct {
	Loc     Loc
	StmtLen int
	Stmt    StmtNode
}

func (s *RawStmt) nodeTag()  {}
func (s *RawStmt) stmtNode() {}
