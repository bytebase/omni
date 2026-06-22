// Package ast defines parse-tree node types for the omni Cassandra CQL parser.
package ast

// Loc represents a source location range (byte offsets).
type Loc struct {
	Start int // inclusive byte offset, -1 = unknown
	End   int // exclusive byte offset, -1 = unknown
}

func NoLoc() Loc { return Loc{Start: -1, End: -1} }

// Node is the interface implemented by all parse-tree nodes.
type Node interface {
	nodeTag()
	GetLoc() Loc
}

// StmtNode marks top-level statement nodes.
type StmtNode interface {
	Node
	stmtNode()
}

// ExprNode marks expression nodes.
type ExprNode interface {
	Node
	exprNode()
}

// List is a generic ordered collection of nodes.
type List struct {
	Items []Node
	Loc   Loc
}

func (*List) nodeTag()      {}
func (l *List) GetLoc() Loc { return l.Loc }
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}

// RawStmt wraps a top-level statement with its location in the input.
type RawStmt struct {
	Stmt         Node
	StmtLocation int
	StmtLen      int
}

func (*RawStmt) nodeTag()      {}
func (r *RawStmt) GetLoc() Loc { return Loc{Start: r.StmtLocation, End: r.StmtLocation + r.StmtLen} }

// ---------------------------------------------------------------------------
// Identifier / Name nodes
// ---------------------------------------------------------------------------

// Identifier represents a simple name (column, table, keyspace, etc.).
type Identifier struct {
	Name   string
	Quoted bool
	Loc    Loc
}

func (*Identifier) nodeTag()      {}
func (n *Identifier) GetLoc() Loc { return n.Loc }
func (*Identifier) exprNode()     {}

// QualifiedName represents a dotted name like keyspace.table.
type QualifiedName struct {
	Parts []*Identifier
	Loc   Loc
}

func (*QualifiedName) nodeTag()      {}
func (n *QualifiedName) GetLoc() Loc { return n.Loc }
func (*QualifiedName) exprNode()     {}

// ---------------------------------------------------------------------------
// Literal nodes
// ---------------------------------------------------------------------------

type StringLit struct {
	Val string
	Loc Loc
}

func (*StringLit) nodeTag()      {}
func (n *StringLit) GetLoc() Loc { return n.Loc }
func (*StringLit) exprNode()     {}

type IntegerLit struct {
	Val string
	Loc Loc
}

func (*IntegerLit) nodeTag()      {}
func (n *IntegerLit) GetLoc() Loc { return n.Loc }
func (*IntegerLit) exprNode()     {}

type FloatLit struct {
	Val string
	Loc Loc
}

func (*FloatLit) nodeTag()      {}
func (n *FloatLit) GetLoc() Loc { return n.Loc }
func (*FloatLit) exprNode()     {}

type BoolLit struct {
	Val bool
	Loc Loc
}

func (*BoolLit) nodeTag()      {}
func (n *BoolLit) GetLoc() Loc { return n.Loc }
func (*BoolLit) exprNode()     {}

type NullLit struct {
	Loc Loc
}

func (*NullLit) nodeTag()      {}
func (n *NullLit) GetLoc() Loc { return n.Loc }
func (*NullLit) exprNode()     {}

type UUIDLit struct {
	Val string
	Loc Loc
}

func (*UUIDLit) nodeTag()      {}
func (n *UUIDLit) GetLoc() Loc { return n.Loc }
func (*UUIDLit) exprNode()     {}

type HexLit struct {
	Val string
	Loc Loc
}

func (*HexLit) nodeTag()      {}
func (n *HexLit) GetLoc() Loc { return n.Loc }
func (*HexLit) exprNode()     {}

type CodeBlock struct {
	Val string
	Loc Loc
}

func (*CodeBlock) nodeTag()      {}
func (n *CodeBlock) GetLoc() Loc { return n.Loc }
func (*CodeBlock) exprNode()     {}

// ---------------------------------------------------------------------------
// Collection literals
// ---------------------------------------------------------------------------

type MapLit struct {
	Keys   []ExprNode
	Values []ExprNode
	Loc    Loc
}

func (*MapLit) nodeTag()      {}
func (n *MapLit) GetLoc() Loc { return n.Loc }
func (*MapLit) exprNode()     {}

type SetLit struct {
	Elements []ExprNode
	Loc      Loc
}

func (*SetLit) nodeTag()      {}
func (n *SetLit) GetLoc() Loc { return n.Loc }
func (*SetLit) exprNode()     {}

type ListLit struct {
	Elements []ExprNode
	Loc      Loc
}

func (*ListLit) nodeTag()      {}
func (n *ListLit) GetLoc() Loc { return n.Loc }
func (*ListLit) exprNode()     {}

type TupleLit struct {
	Elements []ExprNode
	Loc      Loc
}

func (*TupleLit) nodeTag()      {}
func (n *TupleLit) GetLoc() Loc { return n.Loc }
func (*TupleLit) exprNode()     {}

type VectorLit struct {
	Elements []ExprNode
	Loc      Loc
}

func (*VectorLit) nodeTag()      {}
func (n *VectorLit) GetLoc() Loc { return n.Loc }
func (*VectorLit) exprNode()     {}

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

// FunctionCall represents a function invocation like token(col) or now().
type FunctionCall struct {
	Name *Identifier
	Args []ExprNode
	Star bool // e.g. count(*)
	Loc  Loc
}

func (*FunctionCall) nodeTag()      {}
func (n *FunctionCall) GetLoc() Loc { return n.Loc }
func (*FunctionCall) exprNode()     {}

// BinaryExpr represents a binary operation (used in WHERE conditions).
type BinaryExpr struct {
	Left  ExprNode
	Op    string // =, <, >, <=, >=, +, -
	Right ExprNode
	Loc   Loc
}

func (*BinaryExpr) nodeTag()      {}
func (n *BinaryExpr) GetLoc() Loc { return n.Loc }
func (*BinaryExpr) exprNode()     {}

// InExpr represents col IN (val1, val2, ...).
type InExpr struct {
	Column ExprNode
	Values []ExprNode
	Loc    Loc
}

func (*InExpr) nodeTag()      {}
func (n *InExpr) GetLoc() Loc { return n.Loc }
func (*InExpr) exprNode()     {}

// ContainsExpr represents col CONTAINS value or col CONTAINS KEY value.
type ContainsExpr struct {
	Column ExprNode
	Value  ExprNode
	IsKey  bool
	Loc    Loc
}

func (*ContainsExpr) nodeTag()      {}
func (n *ContainsExpr) GetLoc() Loc { return n.Loc }
func (*ContainsExpr) exprNode()     {}

// TupleCompareExpr represents (col1, col2) op (val1, val2).
type TupleCompareExpr struct {
	Columns []ExprNode
	Op      string
	Values  []ExprNode
	Loc     Loc
}

func (*TupleCompareExpr) nodeTag()      {}
func (n *TupleCompareExpr) GetLoc() Loc { return n.Loc }
func (*TupleCompareExpr) exprNode()     {}

// TupleInExpr represents (col1, col2) IN ((v1, v2), (v3, v4)).
type TupleInExpr struct {
	Columns []ExprNode
	Tuples  []*TupleLit
	Loc     Loc
}

func (*TupleInExpr) nodeTag()      {}
func (n *TupleInExpr) GetLoc() Loc { return n.Loc }
func (*TupleInExpr) exprNode()     {}

// IndexAccess represents col[index] (map/list element access).
type IndexAccess struct {
	Collection ExprNode
	Index      ExprNode
	Loc        Loc
}

func (*IndexAccess) nodeTag()      {}
func (n *IndexAccess) GetLoc() Loc { return n.Loc }
func (*IndexAccess) exprNode()     {}

// DotAccess represents keyspace.table or similar dotted access in FROM.
type DotAccess struct {
	Object ExprNode
	Field  *Identifier
	Loc    Loc
}

func (*DotAccess) nodeTag()      {}
func (n *DotAccess) GetLoc() Loc { return n.Loc }
func (*DotAccess) exprNode()     {}

// StarExpr represents * in SELECT *.
type StarExpr struct {
	Loc Loc
}

func (*StarExpr) nodeTag()      {}
func (n *StarExpr) GetLoc() Loc { return n.Loc }
func (*StarExpr) exprNode()     {}

// ---------------------------------------------------------------------------
// Type nodes
// ---------------------------------------------------------------------------

// DataType represents a CQL data type like MAP<text, int>, FROZEN<udt>, VECTOR<float, 3>.
type DataType struct {
	Name       *Identifier
	TypeParams []*DataType
	Dimension  *IntegerLit // for VECTOR<type, dimension>
	Loc        Loc
}

func (*DataType) nodeTag()      {}
func (n *DataType) GetLoc() Loc { return n.Loc }

// ---------------------------------------------------------------------------
// Clause / helper nodes
// ---------------------------------------------------------------------------

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Name       *Identifier
	Type       *DataType
	PrimaryKey bool
	Static     bool
	Loc        Loc
}

func (*ColumnDef) nodeTag()      {}
func (n *ColumnDef) GetLoc() Loc { return n.Loc }

// PrimaryKeyDef represents a PRIMARY KEY definition.
type PrimaryKeyDef struct {
	PartitionKeys  []*Identifier
	ClusteringKeys []*Identifier
	Loc            Loc
}

func (*PrimaryKeyDef) nodeTag()      {}
func (n *PrimaryKeyDef) GetLoc() Loc { return n.Loc }

// ClusteringOrder represents CLUSTERING ORDER BY (col ASC/DESC).
type ClusteringOrder struct {
	Column    *Identifier
	Direction string // "ASC" or "DESC"
	Loc       Loc
}

func (*ClusteringOrder) nodeTag()      {}
func (n *ClusteringOrder) GetLoc() Loc { return n.Loc }

// TableOption represents a table option like key = value.
type TableOption struct {
	Name  *Identifier
	Value ExprNode // StringLit, FloatLit, or OptionHash
	Loc   Loc
}

func (*TableOption) nodeTag()      {}
func (n *TableOption) GetLoc() Loc { return n.Loc }

// OptionHash represents { 'key': 'value', ... }.
type OptionHash struct {
	Items []*OptionHashItem
	Loc   Loc
}

func (*OptionHash) nodeTag()      {}
func (n *OptionHash) GetLoc() Loc { return n.Loc }
func (*OptionHash) exprNode()     {}

type OptionHashItem struct {
	Key   ExprNode
	Value ExprNode
	Loc   Loc
}

func (*OptionHashItem) nodeTag()      {}
func (n *OptionHashItem) GetLoc() Loc { return n.Loc }

// SelectElement represents a single item in a SELECT clause.
type SelectElement struct {
	Expr  ExprNode
	Alias *Identifier
	Loc   Loc
}

func (*SelectElement) nodeTag()      {}
func (n *SelectElement) GetLoc() Loc { return n.Loc }

// AssignmentElement represents col = expr in UPDATE SET.
type AssignmentElement struct {
	Target   ExprNode
	Value    ExprNode
	Operator string // "=", "+=", "-=" (desugared from col = col + val)
	Loc      Loc
}

func (*AssignmentElement) nodeTag()      {}
func (n *AssignmentElement) GetLoc() Loc { return n.Loc }

// IfCondition represents a LWT condition: col = value.
type IfCondition struct {
	Column *Identifier
	Op     string
	Value  ExprNode
	Loc    Loc
}

func (*IfCondition) nodeTag()      {}
func (n *IfCondition) GetLoc() Loc { return n.Loc }

// UsingClause represents USING TTL n AND TIMESTAMP m.
type UsingClause struct {
	TTL       ExprNode
	Timestamp ExprNode
	Loc       Loc
}

func (*UsingClause) nodeTag()      {}
func (n *UsingClause) GetLoc() Loc { return n.Loc }

// OrderByElement represents a single ORDER BY column direction.
type OrderByElement struct {
	Column    ExprNode
	Direction string // "ASC", "DESC", or ""
	// ANN OF support
	IsANN       bool
	AnnVector   ExprNode
	AnnLimit    ExprNode
	Loc         Loc
}

func (*OrderByElement) nodeTag()      {}
func (n *OrderByElement) GetLoc() Loc { return n.Loc }

// ---------------------------------------------------------------------------
// DML statement nodes
// ---------------------------------------------------------------------------

type SelectStmt struct {
	Distinct       bool
	JSON           bool
	Elements       []*SelectElement
	From           *QualifiedName
	Where          []ExprNode // relation elements connected by AND
	OrderBy        []*OrderByElement
	Limit          ExprNode
	AllowFiltering bool
	Loc            Loc
}

func (*SelectStmt) nodeTag()      {}
func (n *SelectStmt) GetLoc() Loc { return n.Loc }
func (*SelectStmt) stmtNode()     {}

type InsertStmt struct {
	Table      *QualifiedName
	Columns    []*Identifier
	Values     []ExprNode
	IsJSON     bool
	JSONValue  ExprNode
	DefaultUnset bool
	IfNotExists bool
	Using      *UsingClause
	Loc        Loc
}

func (*InsertStmt) nodeTag()      {}
func (n *InsertStmt) GetLoc() Loc { return n.Loc }
func (*InsertStmt) stmtNode()     {}

type UpdateStmt struct {
	Table       *QualifiedName
	Using       *UsingClause
	Assignments []*AssignmentElement
	Where       []ExprNode
	IfExists    bool
	IfConditions []*IfCondition
	Loc         Loc
}

func (*UpdateStmt) nodeTag()      {}
func (n *UpdateStmt) GetLoc() Loc { return n.Loc }
func (*UpdateStmt) stmtNode()     {}

type DeleteStmt struct {
	Columns  []ExprNode
	From     *QualifiedName
	Using    *UsingClause
	Where    []ExprNode
	IfExists bool
	IfConditions []*IfCondition
	Loc      Loc
}

func (*DeleteStmt) nodeTag()      {}
func (n *DeleteStmt) GetLoc() Loc { return n.Loc }
func (*DeleteStmt) stmtNode()     {}

// BatchType enumerates BATCH types.
type BatchType int

const (
	BatchDefault  BatchType = iota
	BatchLogged
	BatchUnlogged
)

type BatchStmt struct {
	Type       BatchType
	Using      *UsingClause
	Statements []StmtNode
	Loc        Loc
}

func (*BatchStmt) nodeTag()      {}
func (n *BatchStmt) GetLoc() Loc { return n.Loc }
func (*BatchStmt) stmtNode()     {}

type TruncateStmt struct {
	Table *QualifiedName
	Loc   Loc
}

func (*TruncateStmt) nodeTag()      {}
func (n *TruncateStmt) GetLoc() Loc { return n.Loc }
func (*TruncateStmt) stmtNode()     {}

type UseStmt struct {
	Keyspace *Identifier
	Loc      Loc
}

func (*UseStmt) nodeTag()      {}
func (n *UseStmt) GetLoc() Loc { return n.Loc }
func (*UseStmt) stmtNode()     {}

// ---------------------------------------------------------------------------
// DDL statement nodes
// ---------------------------------------------------------------------------

type CreateKeyspaceStmt struct {
	IfNotExists bool
	Name        *Identifier
	Replication *OptionHash
	DurableWrites *BoolLit
	Loc         Loc
}

func (*CreateKeyspaceStmt) nodeTag()      {}
func (n *CreateKeyspaceStmt) GetLoc() Loc { return n.Loc }
func (*CreateKeyspaceStmt) stmtNode()     {}

type AlterKeyspaceStmt struct {
	Name        *Identifier
	Replication *OptionHash
	DurableWrites *BoolLit
	Loc         Loc
}

func (*AlterKeyspaceStmt) nodeTag()      {}
func (n *AlterKeyspaceStmt) GetLoc() Loc { return n.Loc }
func (*AlterKeyspaceStmt) stmtNode()     {}

type DropKeyspaceStmt struct {
	IfExists bool
	Name     *Identifier
	Loc      Loc
}

func (*DropKeyspaceStmt) nodeTag()      {}
func (n *DropKeyspaceStmt) GetLoc() Loc { return n.Loc }
func (*DropKeyspaceStmt) stmtNode()     {}

type CreateTableStmt struct {
	IfNotExists     bool
	Name            *QualifiedName
	Columns         []*ColumnDef
	PrimaryKey      *PrimaryKeyDef
	Options         []*TableOption
	ClusteringOrders []*ClusteringOrder
	CompactStorage  bool
	Loc             Loc
}

func (*CreateTableStmt) nodeTag()      {}
func (n *CreateTableStmt) GetLoc() Loc { return n.Loc }
func (*CreateTableStmt) stmtNode()     {}

// AlterTableOp enumerates ALTER TABLE operations.
type AlterTableOp int

const (
	AlterTableAdd AlterTableOp = iota
	AlterTableDrop
	AlterTableRename
	AlterTableWith
	AlterTableDropCompactStorage
)

type AlterTableStmt struct {
	Name    *QualifiedName
	Op      AlterTableOp
	AddColumns   []*ColumnDef
	DropColumns  []*Identifier
	RenameFrom   *Identifier
	RenameTo     *Identifier
	Options      []*TableOption
	Loc          Loc
}

func (*AlterTableStmt) nodeTag()      {}
func (n *AlterTableStmt) GetLoc() Loc { return n.Loc }
func (*AlterTableStmt) stmtNode()     {}

type DropTableStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropTableStmt) nodeTag()      {}
func (n *DropTableStmt) GetLoc() Loc { return n.Loc }
func (*DropTableStmt) stmtNode()     {}

type CreateIndexStmt struct {
	IsCustom    bool
	IfNotExists bool
	IndexName   *Identifier
	Table       *QualifiedName
	Column      ExprNode // Identifier or FunctionCall (e.g., FULL(col))
	UsingClass  ExprNode
	Options     *OptionHash
	Loc         Loc
}

func (*CreateIndexStmt) nodeTag()      {}
func (n *CreateIndexStmt) GetLoc() Loc { return n.Loc }
func (*CreateIndexStmt) stmtNode()     {}

type DropIndexStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropIndexStmt) nodeTag()      {}
func (n *DropIndexStmt) GetLoc() Loc { return n.Loc }
func (*DropIndexStmt) stmtNode()     {}

type CreateTypeStmt struct {
	IfNotExists bool
	Name        *QualifiedName
	Fields      []*ColumnDef
	Loc         Loc
}

func (*CreateTypeStmt) nodeTag()      {}
func (n *CreateTypeStmt) GetLoc() Loc { return n.Loc }
func (*CreateTypeStmt) stmtNode()     {}

// AlterTypeOp enumerates ALTER TYPE operations.
type AlterTypeOp int

const (
	AlterTypeAlter AlterTypeOp = iota
	AlterTypeAdd
	AlterTypeRename
)

type AlterTypeStmt struct {
	Name        *QualifiedName
	Op          AlterTypeOp
	AlterColumn *Identifier
	AlterType   *DataType
	AddFields   []*ColumnDef
	Renames     []*AlterTypeRenameItem
	Loc         Loc
}

func (*AlterTypeStmt) nodeTag()      {}
func (n *AlterTypeStmt) GetLoc() Loc { return n.Loc }
func (*AlterTypeStmt) stmtNode()     {}

type AlterTypeRenameItem struct {
	From *Identifier
	To   *Identifier
	Loc  Loc
}

func (*AlterTypeRenameItem) nodeTag()      {}
func (n *AlterTypeRenameItem) GetLoc() Loc { return n.Loc }

type DropTypeStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropTypeStmt) nodeTag()      {}
func (n *DropTypeStmt) GetLoc() Loc { return n.Loc }
func (*DropTypeStmt) stmtNode()     {}

type CreateMVStmt struct {
	IfNotExists      bool
	Name             *QualifiedName
	SelectAll        bool
	SelectColumns    []*Identifier
	FromTable        *QualifiedName
	WhereNotNull     []*Identifier
	WhereRelations   []ExprNode
	PrimaryKey       *PrimaryKeyDef
	Options          []*TableOption
	ClusteringOrders []*ClusteringOrder
	Loc              Loc
}

func (*CreateMVStmt) nodeTag()      {}
func (n *CreateMVStmt) GetLoc() Loc { return n.Loc }
func (*CreateMVStmt) stmtNode()     {}

type AlterMVStmt struct {
	Name    *QualifiedName
	Options []*TableOption
	Loc     Loc
}

func (*AlterMVStmt) nodeTag()      {}
func (n *AlterMVStmt) GetLoc() Loc { return n.Loc }
func (*AlterMVStmt) stmtNode()     {}

type DropMVStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropMVStmt) nodeTag()      {}
func (n *DropMVStmt) GetLoc() Loc { return n.Loc }
func (*DropMVStmt) stmtNode()     {}

// ReturnMode for CREATE FUNCTION.
type ReturnMode int

const (
	ReturnCalledOnNull ReturnMode = iota
	ReturnNullOnNull
)

type CreateFunctionStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *QualifiedName
	Params      []*FunctionParam
	ReturnMode  ReturnMode
	ReturnType  *DataType
	Language    *Identifier
	Body        ExprNode // CodeBlock or StringLit
	Loc         Loc
}

func (*CreateFunctionStmt) nodeTag()      {}
func (n *CreateFunctionStmt) GetLoc() Loc { return n.Loc }
func (*CreateFunctionStmt) stmtNode()     {}

type FunctionParam struct {
	Name *Identifier
	Type *DataType
	Loc  Loc
}

func (*FunctionParam) nodeTag()      {}
func (n *FunctionParam) GetLoc() Loc { return n.Loc }

type DropFunctionStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropFunctionStmt) nodeTag()      {}
func (n *DropFunctionStmt) GetLoc() Loc { return n.Loc }
func (*DropFunctionStmt) stmtNode()     {}

type CreateAggregateStmt struct {
	OrReplace   bool
	IfNotExists bool
	Name        *QualifiedName
	ParamType   *DataType
	SFunc       *Identifier
	SType       *DataType
	FinalFunc   *Identifier
	InitCond    ExprNode
	Loc         Loc
}

func (*CreateAggregateStmt) nodeTag()      {}
func (n *CreateAggregateStmt) GetLoc() Loc { return n.Loc }
func (*CreateAggregateStmt) stmtNode()     {}

type DropAggregateStmt struct {
	IfExists bool
	Name     *QualifiedName
	Loc      Loc
}

func (*DropAggregateStmt) nodeTag()      {}
func (n *DropAggregateStmt) GetLoc() Loc { return n.Loc }
func (*DropAggregateStmt) stmtNode()     {}

type CreateTriggerStmt struct {
	IfNotExists bool
	Name        *Identifier
	Table       *QualifiedName
	UsingClass  ExprNode
	Loc         Loc
}

func (*CreateTriggerStmt) nodeTag()      {}
func (n *CreateTriggerStmt) GetLoc() Loc { return n.Loc }
func (*CreateTriggerStmt) stmtNode()     {}

type DropTriggerStmt struct {
	IfExists bool
	Name     *Identifier
	Table    *QualifiedName
	Loc      Loc
}

func (*DropTriggerStmt) nodeTag()      {}
func (n *DropTriggerStmt) GetLoc() Loc { return n.Loc }
func (*DropTriggerStmt) stmtNode()     {}

// ---------------------------------------------------------------------------
// Auth / Role / User statement nodes
// ---------------------------------------------------------------------------

type CreateRoleStmt struct {
	IfNotExists bool
	Name        *Identifier
	Options     []*RoleOption
	Loc         Loc
}

func (*CreateRoleStmt) nodeTag()      {}
func (n *CreateRoleStmt) GetLoc() Loc { return n.Loc }
func (*CreateRoleStmt) stmtNode()     {}

type RoleOption struct {
	Key   string // "PASSWORD", "LOGIN", "SUPERUSER", "OPTIONS"
	Value ExprNode
	Loc   Loc
}

func (*RoleOption) nodeTag()      {}
func (n *RoleOption) GetLoc() Loc { return n.Loc }

type AlterRoleStmt struct {
	Name    *Identifier
	Options []*RoleOption
	Loc     Loc
}

func (*AlterRoleStmt) nodeTag()      {}
func (n *AlterRoleStmt) GetLoc() Loc { return n.Loc }
func (*AlterRoleStmt) stmtNode()     {}

type DropRoleStmt struct {
	IfExists bool
	Name     *Identifier
	Loc      Loc
}

func (*DropRoleStmt) nodeTag()      {}
func (n *DropRoleStmt) GetLoc() Loc { return n.Loc }
func (*DropRoleStmt) stmtNode()     {}

type CreateUserStmt struct {
	IfNotExists bool
	Name        *Identifier
	Password    ExprNode
	Superuser   *bool
	Loc         Loc
}

func (*CreateUserStmt) nodeTag()      {}
func (n *CreateUserStmt) GetLoc() Loc { return n.Loc }
func (*CreateUserStmt) stmtNode()     {}

type AlterUserStmt struct {
	Name      *Identifier
	Password  ExprNode
	Superuser *bool
	Loc       Loc
}

func (*AlterUserStmt) nodeTag()      {}
func (n *AlterUserStmt) GetLoc() Loc { return n.Loc }
func (*AlterUserStmt) stmtNode()     {}

type DropUserStmt struct {
	IfExists bool
	Name     *Identifier
	Loc      Loc
}

func (*DropUserStmt) nodeTag()      {}
func (n *DropUserStmt) GetLoc() Loc { return n.Loc }
func (*DropUserStmt) stmtNode()     {}

type GrantStmt struct {
	Privilege string
	Resource  *Resource
	Role      *Identifier
	Loc       Loc
}

func (*GrantStmt) nodeTag()      {}
func (n *GrantStmt) GetLoc() Loc { return n.Loc }
func (*GrantStmt) stmtNode()     {}

type RevokeStmt struct {
	Privilege string
	Resource  *Resource
	Role      *Identifier
	Loc       Loc
}

func (*RevokeStmt) nodeTag()      {}
func (n *RevokeStmt) GetLoc() Loc { return n.Loc }
func (*RevokeStmt) stmtNode()     {}

// Resource represents a CQL resource (ALL KEYSPACES, KEYSPACE ks, TABLE ks.t, etc.)
type Resource struct {
	Type     string // "ALL KEYSPACES", "KEYSPACE", "TABLE", "ALL FUNCTIONS", "FUNCTION", "ALL ROLES", "ROLE"
	Name     *QualifiedName
	Loc      Loc
}

func (*Resource) nodeTag()      {}
func (n *Resource) GetLoc() Loc { return n.Loc }

type ListPermissionsStmt struct {
	Privilege string
	Resource  *Resource
	Role      *Identifier
	Loc       Loc
}

func (*ListPermissionsStmt) nodeTag()      {}
func (n *ListPermissionsStmt) GetLoc() Loc { return n.Loc }
func (*ListPermissionsStmt) stmtNode()     {}

type ListRolesStmt struct {
	Of          *Identifier
	NoRecursive bool
	Loc         Loc
}

func (*ListRolesStmt) nodeTag()      {}
func (n *ListRolesStmt) GetLoc() Loc { return n.Loc }
func (*ListRolesStmt) stmtNode()     {}
