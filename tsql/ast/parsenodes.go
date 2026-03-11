package ast

// This file contains T-SQL parse tree node types.
// Reference: https://learn.microsoft.com/en-us/sql/t-sql/statements/

// ---------- Statement nodes ----------

// SelectStmt represents a SELECT statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/select-transact-sql
type SelectStmt struct {
	// CTE (WITH clause)
	WithClause *WithClause

	// DISTINCT or ALL
	Distinct bool
	All      bool

	// TOP clause
	Top *TopClause

	// Target list (result columns)
	TargetList *List

	// INTO clause (SELECT INTO)
	IntoTable *TableRef

	// FROM clause
	FromClause *List

	// WHERE clause
	WhereClause ExprNode

	// GROUP BY clause
	GroupByClause *List

	// HAVING clause
	HavingClause ExprNode

	// ORDER BY clause
	OrderByClause *List

	// OFFSET / FETCH clause
	OffsetClause ExprNode
	FetchClause  *FetchClause

	// FOR clause (FOR XML, FOR JSON)
	ForClause *ForClause

	// OPTION clause (query hints)
	OptionClause *List

	// Set operations (UNION, INTERSECT, EXCEPT)
	Op   SetOperation
	Larg *SelectStmt
	Rarg *SelectStmt

	Loc Loc
}

func (n *SelectStmt) nodeTag()  {}
func (n *SelectStmt) stmtNode() {}

// SetOperation enumerates set operations for SELECT.
type SetOperation int

const (
	SetOpNone      SetOperation = iota
	SetOpUnion                  // UNION
	SetOpIntersect              // INTERSECT
	SetOpExcept                 // EXCEPT
)

// TopClause represents a TOP clause.
type TopClause struct {
	Count    ExprNode // expression for TOP count
	Percent  bool     // TOP ... PERCENT
	WithTies bool     // WITH TIES
	Loc      Loc
}

func (n *TopClause) nodeTag() {}

// FetchClause represents OFFSET...FETCH in T-SQL.
type FetchClause struct {
	Count ExprNode // FETCH NEXT n ROWS ONLY
	Loc   Loc
}

func (n *FetchClause) nodeTag() {}

// ForClause represents FOR XML or FOR JSON.
type ForClause struct {
	Mode    ForMode // XML or JSON
	SubMode string  // RAW, AUTO, PATH, EXPLICIT (XML) or AUTO, PATH (JSON)
	Options *List   // FOR XML options or FOR JSON options
	Loc     Loc
}

func (n *ForClause) nodeTag() {}

// ForMode enumerates FOR clause modes.
type ForMode int

const (
	ForXML ForMode = iota
	ForJSON
)

// WithClause represents a WITH (CTE) clause.
type WithClause struct {
	CTEs *List // list of CommonTableExpr
	Loc  Loc
}

func (n *WithClause) nodeTag() {}

// CommonTableExpr represents a single CTE definition.
type CommonTableExpr struct {
	Name    string
	Columns *List    // optional column name list
	Query   StmtNode // the CTE query (SelectStmt)
	Loc     Loc
}

func (n *CommonTableExpr) nodeTag() {}

// InsertStmt represents an INSERT statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/insert-transact-sql
type InsertStmt struct {
	WithClause   *WithClause
	Top          *TopClause
	Relation     *TableRef
	Cols         *List // column name list
	Source       Node  // SELECT, VALUES, EXEC, or DEFAULT VALUES
	OutputClause *OutputClause
	Loc          Loc
}

func (n *InsertStmt) nodeTag()  {}
func (n *InsertStmt) stmtNode() {}

// UpdateStmt represents an UPDATE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/update-transact-sql
type UpdateStmt struct {
	WithClause   *WithClause
	Top          *TopClause
	Relation     *TableRef
	SetClause    *List // list of SetExpr
	OutputClause *OutputClause
	FromClause   *List
	WhereClause  ExprNode
	Loc          Loc
}

func (n *UpdateStmt) nodeTag()  {}
func (n *UpdateStmt) stmtNode() {}

// DeleteStmt represents a DELETE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/delete-transact-sql
type DeleteStmt struct {
	WithClause   *WithClause
	Top          *TopClause
	Relation     *TableRef
	OutputClause *OutputClause
	FromClause   *List
	WhereClause  ExprNode
	Loc          Loc
}

func (n *DeleteStmt) nodeTag()  {}
func (n *DeleteStmt) stmtNode() {}

// MergeStmt represents a MERGE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/merge-transact-sql
type MergeStmt struct {
	WithClause   *WithClause
	Target       *TableRef
	Source       TableExpr // table source
	SourceAlias  string
	OnCondition  ExprNode
	WhenClauses  *List // list of MergeWhenClause
	OutputClause *OutputClause
	Loc          Loc
}

func (n *MergeStmt) nodeTag()  {}
func (n *MergeStmt) stmtNode() {}

// MergeWhenClause represents a WHEN clause in MERGE.
type MergeWhenClause struct {
	Matched   bool     // WHEN MATCHED vs WHEN NOT MATCHED
	ByTarget  bool     // BY TARGET (default) vs BY SOURCE
	Condition ExprNode // AND condition
	Action    Node     // UpdateAction, DeleteAction, or InsertAction
	Loc       Loc
}

func (n *MergeWhenClause) nodeTag() {}

// MergeUpdateAction represents UPDATE SET in a MERGE WHEN clause.
type MergeUpdateAction struct {
	SetClause *List // list of SetExpr
	Loc       Loc
}

func (n *MergeUpdateAction) nodeTag() {}

// MergeDeleteAction represents DELETE in a MERGE WHEN clause.
type MergeDeleteAction struct {
	Loc Loc
}

func (n *MergeDeleteAction) nodeTag() {}

// MergeInsertAction represents INSERT in a MERGE WHEN clause.
type MergeInsertAction struct {
	Cols   *List // column list
	Values *List // VALUES list
	Loc    Loc
}

func (n *MergeInsertAction) nodeTag() {}

// OutputClause represents an OUTPUT clause.
type OutputClause struct {
	Targets   *List // output expressions
	IntoTable *TableRef
	IntoCols  *List // column list for INTO
	Loc       Loc
}

func (n *OutputClause) nodeTag() {}

// CreateTableStmt represents a CREATE TABLE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
type CreateTableStmt struct {
	Name        *TableRef
	Columns     *List // list of ColumnDef
	Constraints *List // table-level constraints
	IfNotExists bool
	Loc         Loc
}

func (n *CreateTableStmt) nodeTag()  {}
func (n *CreateTableStmt) stmtNode() {}

// ColumnDef represents a column definition in CREATE TABLE.
type ColumnDef struct {
	Name        string
	DataType    *DataType
	Identity    *IdentitySpec
	Computed    *ComputedColumnDef
	DefaultExpr ExprNode
	Collation   string
	Constraints *List // column-level constraints
	Nullable    *NullableSpec
	Loc         Loc
}

func (n *ColumnDef) nodeTag() {}

// NullableSpec indicates NULL / NOT NULL for a column.
type NullableSpec struct {
	NotNull bool
	Loc     Loc
}

func (n *NullableSpec) nodeTag() {}

// IdentitySpec represents IDENTITY(seed, increment).
type IdentitySpec struct {
	Seed      int64
	Increment int64
	Loc       Loc
}

func (n *IdentitySpec) nodeTag() {}

// ComputedColumnDef represents AS (expression) for computed columns.
type ComputedColumnDef struct {
	Expr      ExprNode
	Persisted bool
	Loc       Loc
}

func (n *ComputedColumnDef) nodeTag() {}

// ConstraintDef represents a table or column constraint.
type ConstraintDef struct {
	Type       ConstraintType
	Name       string   // constraint name (optional)
	Columns    *List    // columns for PK, UNIQUE, INDEX
	Expr       ExprNode // CHECK expression, default expression
	RefTable   *TableRef
	RefColumns *List // FK referenced columns
	OnDelete   ReferentialAction
	OnUpdate   ReferentialAction
	Clustered  *bool // true=CLUSTERED, false=NONCLUSTERED, nil=unspecified
	Loc        Loc
}

func (n *ConstraintDef) nodeTag() {}

// ConstraintType enumerates constraint types.
type ConstraintType int

const (
	ConstraintPrimaryKey ConstraintType = iota
	ConstraintUnique
	ConstraintCheck
	ConstraintDefault
	ConstraintForeignKey
	ConstraintNotNull
)

// ReferentialAction enumerates FK actions.
type ReferentialAction int

const (
	RefActNone ReferentialAction = iota
	RefActCascade
	RefActSetNull
	RefActSetDefault
	RefActNoAction
)

// AlterTableStmt represents an ALTER TABLE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-table-transact-sql
type AlterTableStmt struct {
	Name    *TableRef
	Actions *List // list of AlterTableAction
	Loc     Loc
}

func (n *AlterTableStmt) nodeTag()  {}
func (n *AlterTableStmt) stmtNode() {}

// AlterTableAction represents a single ALTER TABLE action.
type AlterTableAction struct {
	Type       AlterTableActionType
	Column     *ColumnDef
	ColName    string
	Constraint *ConstraintDef
	DataType   *DataType
	Loc        Loc
}

func (n *AlterTableAction) nodeTag() {}

// AlterTableActionType enumerates ALTER TABLE actions.
type AlterTableActionType int

const (
	ATAddColumn AlterTableActionType = iota
	ATDropColumn
	ATAlterColumn
	ATAddConstraint
	ATDropConstraint
)

// DropStmt represents a DROP statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-table-transact-sql
type DropStmt struct {
	ObjectType DropObjectType
	Names      *List // list of TableRef
	IfExists   bool
	Loc        Loc
}

func (n *DropStmt) nodeTag()  {}
func (n *DropStmt) stmtNode() {}

// DropObjectType enumerates droppable object types.
type DropObjectType int

const (
	DropTable DropObjectType = iota
	DropView
	DropIndex
	DropProcedure
	DropFunction
	DropDatabase
	DropSchema
	DropTrigger
	DropType
)

// CreateIndexStmt represents a CREATE INDEX statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql
type CreateIndexStmt struct {
	Unique      bool
	Clustered   *bool // true=CLUSTERED, false=NONCLUSTERED
	Columnstore bool
	Name        string
	Table       *TableRef
	Columns     *List    // index columns
	IncludeCols *List    // INCLUDE columns
	WhereClause ExprNode // filtered index
	Options     *List    // WITH options
	OnFileGroup string
	Loc         Loc
}

func (n *CreateIndexStmt) nodeTag()  {}
func (n *CreateIndexStmt) stmtNode() {}

// IndexColumn represents a column in an index definition.
type IndexColumn struct {
	Name    string
	Expr    ExprNode // expression for computed indexes
	SortDir SortDirection
	Loc     Loc
}

func (n *IndexColumn) nodeTag() {}

// SortDirection enumerates sort directions.
type SortDirection int

const (
	SortDefault SortDirection = iota
	SortAsc
	SortDesc
)

// CreateViewStmt represents a CREATE VIEW statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-view-transact-sql
type CreateViewStmt struct {
	OrAlter       bool
	Name          *TableRef
	Columns       *List    // column name list
	Query         StmtNode // SelectStmt
	WithCheck     bool     // WITH CHECK OPTION
	SchemaBinding bool     // WITH SCHEMABINDING
	Loc           Loc
}

func (n *CreateViewStmt) nodeTag()  {}
func (n *CreateViewStmt) stmtNode() {}

// CreateFunctionStmt represents a CREATE FUNCTION statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-function-transact-sql
type CreateFunctionStmt struct {
	OrAlter      bool
	Name         *TableRef
	Params       *List // list of ParamDef
	ReturnType   *DataType
	ReturnsTable *ReturnsTableDef
	Body         Node // BeginEndStmt or single expression
	Options      *List
	Loc          Loc
}

func (n *CreateFunctionStmt) nodeTag()  {}
func (n *CreateFunctionStmt) stmtNode() {}

// ReturnsTableDef represents RETURNS TABLE (...) in a function definition.
type ReturnsTableDef struct {
	Columns  *List  // list of ColumnDef for inline table-valued function
	Variable string // @variable for multi-statement TVF
	Loc      Loc
}

func (n *ReturnsTableDef) nodeTag() {}

// ParamDef represents a parameter definition.
type ParamDef struct {
	Name     string // @param
	DataType *DataType
	Default  ExprNode
	Output   bool // OUTPUT keyword
	ReadOnly bool // READONLY keyword
	Loc      Loc
}

func (n *ParamDef) nodeTag() {}

// CreateProcedureStmt represents a CREATE PROCEDURE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-procedure-transact-sql
type CreateProcedureStmt struct {
	OrAlter bool
	Name    *TableRef
	Params  *List // list of ParamDef
	Body    Node  // BeginEndStmt
	Options *List
	Loc     Loc
}

func (n *CreateProcedureStmt) nodeTag()  {}
func (n *CreateProcedureStmt) stmtNode() {}

// CreateDatabaseStmt represents a CREATE DATABASE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-transact-sql
type CreateDatabaseStmt struct {
	Name    string
	Options *List
	Loc     Loc
}

func (n *CreateDatabaseStmt) nodeTag()  {}
func (n *CreateDatabaseStmt) stmtNode() {}

// TruncateStmt represents a TRUNCATE TABLE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/truncate-table-transact-sql
type TruncateStmt struct {
	Table *TableRef
	Loc   Loc
}

func (n *TruncateStmt) nodeTag()  {}
func (n *TruncateStmt) stmtNode() {}

// ---------- Control flow statements ----------

// DeclareStmt represents a DECLARE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/declare-local-variable-transact-sql
type DeclareStmt struct {
	Variables *List // list of VariableDecl
	Loc       Loc
}

func (n *DeclareStmt) nodeTag()  {}
func (n *DeclareStmt) stmtNode() {}

// VariableDecl represents a variable declaration in DECLARE.
type VariableDecl struct {
	Name     string    // @varname
	DataType *DataType // type
	Default  ExprNode  // = expression
	IsTable  bool      // TABLE type
	TableDef *List     // column defs for table variable
	IsCursor bool      // CURSOR
	Loc      Loc
}

func (n *VariableDecl) nodeTag() {}

// SetStmt represents a SET statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/set-local-variable-transact-sql
type SetStmt struct {
	Variable string   // @variable or SET option name
	Value    ExprNode // expression
	Loc      Loc
}

func (n *SetStmt) nodeTag()  {}
func (n *SetStmt) stmtNode() {}

// IfStmt represents an IF...ELSE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/if-else-transact-sql
type IfStmt struct {
	Condition ExprNode
	Then      StmtNode // single statement or BeginEndStmt
	Else      StmtNode // optional ELSE branch
	Loc       Loc
}

func (n *IfStmt) nodeTag()  {}
func (n *IfStmt) stmtNode() {}

// WhileStmt represents a WHILE loop.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/while-transact-sql
type WhileStmt struct {
	Condition ExprNode
	Body      StmtNode
	Loc       Loc
}

func (n *WhileStmt) nodeTag()  {}
func (n *WhileStmt) stmtNode() {}

// BeginEndStmt represents a BEGIN...END block.
type BeginEndStmt struct {
	Stmts *List
	Loc   Loc
}

func (n *BeginEndStmt) nodeTag()  {}
func (n *BeginEndStmt) stmtNode() {}

// TryCatchStmt represents BEGIN TRY...END TRY BEGIN CATCH...END CATCH.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/try-catch-transact-sql
type TryCatchStmt struct {
	TryBlock   *List
	CatchBlock *List
	Loc        Loc
}

func (n *TryCatchStmt) nodeTag()  {}
func (n *TryCatchStmt) stmtNode() {}

// ReturnStmt represents a RETURN statement.
type ReturnStmt struct {
	Value ExprNode // optional return expression
	Loc   Loc
}

func (n *ReturnStmt) nodeTag()  {}
func (n *ReturnStmt) stmtNode() {}

// BreakStmt represents a BREAK statement.
type BreakStmt struct {
	Loc Loc
}

func (n *BreakStmt) nodeTag()  {}
func (n *BreakStmt) stmtNode() {}

// ContinueStmt represents a CONTINUE statement.
type ContinueStmt struct {
	Loc Loc
}

func (n *ContinueStmt) nodeTag()  {}
func (n *ContinueStmt) stmtNode() {}

// GotoStmt represents a GOTO label statement.
type GotoStmt struct {
	Label string
	Loc   Loc
}

func (n *GotoStmt) nodeTag()  {}
func (n *GotoStmt) stmtNode() {}

// LabelStmt represents a label: statement.
type LabelStmt struct {
	Label string
	Loc   Loc
}

func (n *LabelStmt) nodeTag()  {}
func (n *LabelStmt) stmtNode() {}

// WaitForStmt represents a WAITFOR statement.
type WaitForStmt struct {
	WaitType string // DELAY or TIME
	Value    ExprNode
	Loc      Loc
}

func (n *WaitForStmt) nodeTag()  {}
func (n *WaitForStmt) stmtNode() {}

// ---------- Execution / utility statements ----------

// ExecStmt represents an EXEC/EXECUTE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/execute-transact-sql
type ExecStmt struct {
	Name      *TableRef
	Args      *List  // list of ExecArg
	ReturnVar string // @var = EXEC ...
	Loc       Loc
}

func (n *ExecStmt) nodeTag()  {}
func (n *ExecStmt) stmtNode() {}

// ExecArg represents an argument in EXEC.
type ExecArg struct {
	Name   string // @param = value (named) or empty (positional)
	Value  ExprNode
	Output bool // OUTPUT keyword
	Loc    Loc
}

func (n *ExecArg) nodeTag() {}

// PrintStmt represents a PRINT statement.
type PrintStmt struct {
	Expr ExprNode
	Loc  Loc
}

func (n *PrintStmt) nodeTag()  {}
func (n *PrintStmt) stmtNode() {}

// RaiseErrorStmt represents a RAISERROR statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/raiserror-transact-sql
type RaiseErrorStmt struct {
	Message  ExprNode // message string or number
	Severity ExprNode
	State    ExprNode
	Args     *List // optional formatting args
	Options  *List // WITH options (LOG, NOWAIT, SETERROR)
	Loc      Loc
}

func (n *RaiseErrorStmt) nodeTag()  {}
func (n *RaiseErrorStmt) stmtNode() {}

// ThrowStmt represents a THROW statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/throw-transact-sql
type ThrowStmt struct {
	ErrorNumber ExprNode
	Message     ExprNode
	State       ExprNode
	Loc         Loc
}

func (n *ThrowStmt) nodeTag()  {}
func (n *ThrowStmt) stmtNode() {}

// UseStmt represents a USE database statement.
type UseStmt struct {
	Database string
	Loc      Loc
}

func (n *UseStmt) nodeTag()  {}
func (n *UseStmt) stmtNode() {}

// GoStmt represents a GO batch separator.
type GoStmt struct {
	Count int // optional repeat count
	Loc   Loc
}

func (n *GoStmt) nodeTag()  {}
func (n *GoStmt) stmtNode() {}

// ---------- Transaction statements ----------

// BeginTransStmt represents BEGIN TRANSACTION.
type BeginTransStmt struct {
	Name string
	Loc  Loc
}

func (n *BeginTransStmt) nodeTag()  {}
func (n *BeginTransStmt) stmtNode() {}

// CommitTransStmt represents COMMIT TRANSACTION.
type CommitTransStmt struct {
	Name string
	Loc  Loc
}

func (n *CommitTransStmt) nodeTag()  {}
func (n *CommitTransStmt) stmtNode() {}

// RollbackTransStmt represents ROLLBACK TRANSACTION.
type RollbackTransStmt struct {
	Name      string
	Savepoint string
	Loc       Loc
}

func (n *RollbackTransStmt) nodeTag()  {}
func (n *RollbackTransStmt) stmtNode() {}

// SaveTransStmt represents SAVE TRANSACTION.
type SaveTransStmt struct {
	Name string
	Loc  Loc
}

func (n *SaveTransStmt) nodeTag()  {}
func (n *SaveTransStmt) stmtNode() {}

// ---------- Security statements ----------

// GrantStmt represents GRANT/REVOKE/DENY.
type GrantStmt struct {
	StmtType   GrantType // GRANT, REVOKE, DENY
	Privileges *List
	OnType     string // TABLE, VIEW, PROCEDURE, etc.
	OnName     *TableRef
	Principals *List // TO/FROM principals
	WithGrant  bool  // WITH GRANT OPTION
	CascadeOpt bool  // CASCADE
	Loc        Loc
}

func (n *GrantStmt) nodeTag()  {}
func (n *GrantStmt) stmtNode() {}

// GrantType enumerates grant statement types.
type GrantType int

const (
	GrantTypeGrant GrantType = iota
	GrantTypeRevoke
	GrantTypeDeny
)

// ---------- Expression nodes ----------

// BinaryExpr represents a binary operation.
type BinaryExpr struct {
	Op    BinaryOp
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (n *BinaryExpr) nodeTag()  {}
func (n *BinaryExpr) exprNode() {}

// BinaryOp enumerates binary operators.
type BinaryOp int

const (
	BinOpAdd BinaryOp = iota
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpMod
	BinOpBitAnd
	BinOpBitOr
	BinOpBitXor
	BinOpEq
	BinOpNeq
	BinOpLt
	BinOpGt
	BinOpLte
	BinOpGte
	BinOpNotLt // !<
	BinOpNotGt // !>
	BinOpAnd
	BinOpOr
)

// UnaryExpr represents a unary operation.
type UnaryExpr struct {
	Op      UnaryOp
	Operand ExprNode
	Loc     Loc
}

func (n *UnaryExpr) nodeTag()  {}
func (n *UnaryExpr) exprNode() {}

// UnaryOp enumerates unary operators.
type UnaryOp int

const (
	UnaryPlus UnaryOp = iota
	UnaryMinus
	UnaryBitNot // ~
	UnaryNot    // NOT
)

// FuncCallExpr represents a function call.
type FuncCallExpr struct {
	Name     *TableRef // potentially schema-qualified
	Args     *List
	Distinct bool
	Star     bool // e.g., COUNT(*)
	Over     *OverClause
	Within   *List // WITHIN GROUP (ORDER BY ...)
	Loc      Loc
}

func (n *FuncCallExpr) nodeTag()   {}
func (n *FuncCallExpr) exprNode()  {}
func (n *FuncCallExpr) tableExpr() {} // table-valued functions can appear in FROM

// CaseExpr represents a CASE expression.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/case-transact-sql
type CaseExpr struct {
	Arg      ExprNode // simple CASE argument, nil for searched CASE
	WhenList *List    // list of CaseWhen
	ElseExpr ExprNode
	Loc      Loc
}

func (n *CaseExpr) nodeTag()  {}
func (n *CaseExpr) exprNode() {}

// CaseWhen represents a WHEN clause in CASE.
type CaseWhen struct {
	Condition ExprNode
	Result    ExprNode
	Loc       Loc
}

func (n *CaseWhen) nodeTag()  {}
func (n *CaseWhen) exprNode() {}

// BetweenExpr represents expr BETWEEN low AND high.
type BetweenExpr struct {
	Expr ExprNode
	Low  ExprNode
	High ExprNode
	Not  bool
	Loc  Loc
}

func (n *BetweenExpr) nodeTag()  {}
func (n *BetweenExpr) exprNode() {}

// InExpr represents expr IN (values...) or expr IN (subquery).
type InExpr struct {
	Expr     ExprNode
	List     *List    // value list
	Subquery ExprNode // subquery (SubqueryExpr)
	Not      bool
	Loc      Loc
}

func (n *InExpr) nodeTag()  {}
func (n *InExpr) exprNode() {}

// LikeExpr represents expr LIKE pattern.
type LikeExpr struct {
	Expr    ExprNode
	Pattern ExprNode
	Escape  ExprNode
	Not     bool
	Loc     Loc
}

func (n *LikeExpr) nodeTag()  {}
func (n *LikeExpr) exprNode() {}

// IsExpr represents IS NULL, IS NOT NULL, etc.
type IsExpr struct {
	Expr     ExprNode
	TestType IsTestType
	Loc      Loc
}

func (n *IsExpr) nodeTag()  {}
func (n *IsExpr) exprNode() {}

// IsTestType enumerates IS test types.
type IsTestType int

const (
	IsNull IsTestType = iota
	IsNotNull
	IsTrue
	IsNotTrue
	IsFalse
	IsNotFalse
)

// ExistsExpr represents EXISTS (subquery).
type ExistsExpr struct {
	Subquery StmtNode
	Loc      Loc
}

func (n *ExistsExpr) nodeTag()  {}
func (n *ExistsExpr) exprNode() {}

// CastExpr represents CAST(expr AS type).
type CastExpr struct {
	Expr     ExprNode
	DataType *DataType
	Loc      Loc
}

func (n *CastExpr) nodeTag()  {}
func (n *CastExpr) exprNode() {}

// ConvertExpr represents CONVERT(type, expr [, style]).
type ConvertExpr struct {
	DataType *DataType
	Expr     ExprNode
	Style    ExprNode
	Loc      Loc
}

func (n *ConvertExpr) nodeTag()  {}
func (n *ConvertExpr) exprNode() {}

// TryCastExpr represents TRY_CAST(expr AS type).
type TryCastExpr struct {
	Expr     ExprNode
	DataType *DataType
	Loc      Loc
}

func (n *TryCastExpr) nodeTag()  {}
func (n *TryCastExpr) exprNode() {}

// TryConvertExpr represents TRY_CONVERT(type, expr [, style]).
type TryConvertExpr struct {
	DataType *DataType
	Expr     ExprNode
	Style    ExprNode
	Loc      Loc
}

func (n *TryConvertExpr) nodeTag()  {}
func (n *TryConvertExpr) exprNode() {}

// CoalesceExpr represents COALESCE(expr, expr, ...).
type CoalesceExpr struct {
	Args *List
	Loc  Loc
}

func (n *CoalesceExpr) nodeTag()  {}
func (n *CoalesceExpr) exprNode() {}

// NullifExpr represents NULLIF(expr1, expr2).
type NullifExpr struct {
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (n *NullifExpr) nodeTag()  {}
func (n *NullifExpr) exprNode() {}

// IifExpr represents IIF(condition, true_val, false_val).
type IifExpr struct {
	Condition ExprNode
	TrueVal   ExprNode
	FalseVal  ExprNode
	Loc       Loc
}

func (n *IifExpr) nodeTag()  {}
func (n *IifExpr) exprNode() {}

// ColumnRef represents a column reference (possibly qualified).
type ColumnRef struct {
	Server   string // linked server
	Database string
	Schema   string
	Table    string
	Column   string
	Loc      Loc
}

func (n *ColumnRef) nodeTag()  {}
func (n *ColumnRef) exprNode() {}

// VariableRef represents @variable or @@systemvariable.
type VariableRef struct {
	Name string // includes @ or @@ prefix
	Loc  Loc
}

func (n *VariableRef) nodeTag()  {}
func (n *VariableRef) exprNode() {}

// StarExpr represents * in SELECT.
type StarExpr struct {
	Qualifier string // optional table qualifier
	Loc       Loc
}

func (n *StarExpr) nodeTag()  {}
func (n *StarExpr) exprNode() {}

// Literal represents a literal value.
type Literal struct {
	Type    LiteralType
	Str     string
	Ival    int64
	IsNChar bool // N'...' string
	Loc     Loc
}

func (n *Literal) nodeTag()  {}
func (n *Literal) exprNode() {}

// LiteralType enumerates literal types.
type LiteralType int

const (
	LitString LiteralType = iota
	LitInteger
	LitFloat
	LitNull
	LitDefault
)

// SubqueryExpr represents a scalar subquery in an expression.
type SubqueryExpr struct {
	Query StmtNode // SelectStmt
	Loc   Loc
}

func (n *SubqueryExpr) nodeTag()   {}
func (n *SubqueryExpr) exprNode()  {}
func (n *SubqueryExpr) tableExpr() {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Expr ExprNode
	Loc  Loc
}

func (n *ParenExpr) nodeTag()  {}
func (n *ParenExpr) exprNode() {}

// ---------- Table/Name reference nodes ----------

// TableRef represents a qualified table/object name.
// server.database.schema.object
type TableRef struct {
	Server   string
	Database string
	Schema   string
	Object   string
	Alias    string
	Loc      Loc
}

func (n *TableRef) nodeTag()   {}
func (n *TableRef) tableExpr() {}

// DataType represents a T-SQL data type reference.
type DataType struct {
	Name      string   // INT, VARCHAR, NVARCHAR, DECIMAL, etc.
	Schema    string   // optional schema qualifier
	Precision ExprNode // e.g., DECIMAL(10, 2) -> 10
	Scale     ExprNode // e.g., DECIMAL(10, 2) -> 2
	Length    ExprNode // e.g., VARCHAR(100) -> 100
	MaxLength bool     // VARCHAR(MAX)
	Loc       Loc
}

func (n *DataType) nodeTag() {}

// ---------- FROM / JOIN nodes ----------

// JoinClause represents a JOIN expression.
type JoinClause struct {
	Type      JoinType
	Left      TableExpr
	Right     TableExpr
	Condition ExprNode // ON condition
	Using     *List    // USING column list
	Loc       Loc
}

func (n *JoinClause) nodeTag()   {}
func (n *JoinClause) tableExpr() {}

// JoinType enumerates join types.
type JoinType int

const (
	JoinInner JoinType = iota
	JoinLeft
	JoinRight
	JoinFull
	JoinCross
	JoinCrossApply // CROSS APPLY
	JoinOuterApply // OUTER APPLY
)

// AliasedTableRef represents a table reference with optional alias.
type AliasedTableRef struct {
	Table   TableExpr // TableRef, SubqueryExpr, etc.
	Alias   string
	Columns *List // alias column list
	Loc     Loc
}

func (n *AliasedTableRef) nodeTag()   {}
func (n *AliasedTableRef) tableExpr() {}

// ---------- Window / OVER clause ----------

// OverClause represents an OVER clause for window functions.
type OverClause struct {
	PartitionBy *List // PARTITION BY expressions
	OrderBy     *List // ORDER BY items
	WindowFrame *WindowFrame
	WindowName  string // reference to named window
	Loc         Loc
}

func (n *OverClause) nodeTag() {}

// WindowFrame represents the frame specification of a window.
type WindowFrame struct {
	Type  WindowFrameType
	Start *WindowBound
	End   *WindowBound
	Loc   Loc
}

func (n *WindowFrame) nodeTag() {}

// WindowFrameType enumerates window frame types.
type WindowFrameType int

const (
	FrameRows WindowFrameType = iota
	FrameRange
	FrameGroups
)

// WindowBound represents a window frame bound.
type WindowBound struct {
	Type   WindowBoundType
	Offset ExprNode // expression for N PRECEDING / N FOLLOWING
	Loc    Loc
}

func (n *WindowBound) nodeTag() {}

// WindowBoundType enumerates window frame bound types.
type WindowBoundType int

const (
	BoundUnboundedPreceding WindowBoundType = iota
	BoundPreceding
	BoundCurrentRow
	BoundFollowing
	BoundUnboundedFollowing
)

// OrderByItem represents an ORDER BY item.
type OrderByItem struct {
	Expr       ExprNode
	SortDir    SortDirection
	NullsOrder NullsOrder
	Loc        Loc
}

func (n *OrderByItem) nodeTag() {}

// NullsOrder enumerates NULLS FIRST / NULLS LAST.
type NullsOrder int

const (
	NullsDefault NullsOrder = iota
	NullsFirst
	NullsLast
)

// ---------- SET clause ----------

// SetExpr represents SET column = expr in UPDATE.
type SetExpr struct {
	Column   *ColumnRef
	Variable string // @var = expr
	Value    ExprNode
	Loc      Loc
}

func (n *SetExpr) nodeTag() {}

// ---------- VALUES clause ----------

// ValuesClause represents VALUES (...), (...).
type ValuesClause struct {
	Rows *List // list of List (each row is a list of expressions)
	Loc  Loc
}

func (n *ValuesClause) nodeTag() {}

// ---------- Result target ----------

// ResTarget represents a result column in SELECT.
type ResTarget struct {
	Name string   // alias (AS name)
	Val  ExprNode // expression
	Loc  Loc
}

func (n *ResTarget) nodeTag()  {}
func (n *ResTarget) exprNode() {}

// ---------- Assignment in SELECT ----------

// SelectAssign represents SELECT @var = expr.
type SelectAssign struct {
	Variable string
	Value    ExprNode
	Loc      Loc
}

func (n *SelectAssign) nodeTag()  {}
func (n *SelectAssign) exprNode() {}

// ---------- Cursor statements ----------

// DeclareCursorStmt represents a DECLARE cursor_name CURSOR statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/declare-cursor-transact-sql
//
// ISO syntax:
//
//	DECLARE cursor_name [ INSENSITIVE ] [ SCROLL ] CURSOR
//	    FOR select_statement
//	    [ FOR { READ_ONLY | UPDATE [ OF column_name [ , ...n ] ] } ]
//
// Transact-SQL extended syntax:
//
//	DECLARE cursor_name CURSOR [ LOCAL | GLOBAL ]
//	    [ FORWARD_ONLY | SCROLL ]
//	    [ STATIC | KEYSET | DYNAMIC | FAST_FORWARD ]
//	    [ READ_ONLY | SCROLL_LOCKS | OPTIMISTIC ]
//	    [ TYPE_WARNING ]
//	    FOR select_statement
//	    [ FOR UPDATE [ OF column_name [ , ...n ] ] ]
type DeclareCursorStmt struct {
	Name        string   // cursor name
	Insensitive bool     // INSENSITIVE (ISO)
	Scroll      bool     // SCROLL
	Scope       string   // "LOCAL" or "GLOBAL" (T-SQL extended, empty = default)
	ForwardOnly bool     // FORWARD_ONLY (T-SQL extended)
	CursorType  string   // "STATIC", "KEYSET", "DYNAMIC", "FAST_FORWARD" (empty = default)
	Concurrency string   // "READ_ONLY", "SCROLL_LOCKS", "OPTIMISTIC" (empty = default)
	TypeWarning bool     // TYPE_WARNING
	Query       StmtNode // SELECT statement
	ForUpdate   bool     // FOR UPDATE
	UpdateCols  *List    // OF column_name [,...n] (nil if no column list)
	Loc         Loc
}

func (n *DeclareCursorStmt) nodeTag()  {}
func (n *DeclareCursorStmt) stmtNode() {}

// OpenCursorStmt represents an OPEN cursor statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/open-transact-sql
//
//	OPEN { { [ GLOBAL ] cursor_name } | cursor_variable_name }
type OpenCursorStmt struct {
	Name   string // cursor name or @cursor_variable
	Global bool   // GLOBAL keyword specified
	Loc    Loc
}

func (n *OpenCursorStmt) nodeTag()  {}
func (n *OpenCursorStmt) stmtNode() {}

// FetchCursorStmt represents a FETCH cursor statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/fetch-transact-sql
//
//	FETCH
//	    [ [ NEXT | PRIOR | FIRST | LAST
//	            | ABSOLUTE { n | @nvar }
//	            | RELATIVE { n | @nvar }
//	       ]
//	       FROM
//	    ]
//	{ { [ GLOBAL ] cursor_name } | @cursor_variable_name }
//	[ INTO @variable_name [ ,...n ] ]
type FetchCursorStmt struct {
	Orientation string   // "NEXT", "PRIOR", "FIRST", "LAST", "ABSOLUTE", "RELATIVE" (empty = default NEXT)
	FetchOffset ExprNode // offset expression for ABSOLUTE/RELATIVE
	Name        string   // cursor name or @cursor_variable
	Global      bool     // GLOBAL keyword specified
	IntoVars    *List    // INTO @var1, @var2, ... (list of String nodes)
	Loc         Loc
}

func (n *FetchCursorStmt) nodeTag()  {}
func (n *FetchCursorStmt) stmtNode() {}

// CloseCursorStmt represents a CLOSE cursor statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/close-transact-sql
//
//	CLOSE { { [ GLOBAL ] cursor_name } | cursor_variable_name }
type CloseCursorStmt struct {
	Name   string // cursor name or @cursor_variable
	Global bool   // GLOBAL keyword specified
	Loc    Loc
}

func (n *CloseCursorStmt) nodeTag()  {}
func (n *CloseCursorStmt) stmtNode() {}

// DeallocateCursorStmt represents a DEALLOCATE cursor statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/deallocate-transact-sql
//
//	DEALLOCATE { { [ GLOBAL ] cursor_name } | @cursor_variable_name }
type DeallocateCursorStmt struct {
	Name   string // cursor name or @cursor_variable
	Global bool   // GLOBAL keyword specified
	Loc    Loc
}

func (n *DeallocateCursorStmt) nodeTag()  {}
func (n *DeallocateCursorStmt) stmtNode() {}

// ---------- Method call ----------

// MethodCallExpr represents a :: static method call.
type MethodCallExpr struct {
	Type   *DataType
	Method string
	Args   *List
	Loc    Loc
}

func (n *MethodCallExpr) nodeTag()  {}
func (n *MethodCallExpr) exprNode() {}
