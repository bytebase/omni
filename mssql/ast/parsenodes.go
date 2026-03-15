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
	GroupByAll    bool // GROUP BY ALL (deprecated)

	// HAVING clause
	HavingClause ExprNode

	// WINDOW clause (named window definitions)
	WindowClause *List // list of WindowDef

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

	// Element name for RAW('name') or PATH('name')
	ElementName string

	// Common directives
	BinaryBase64 bool   // BINARY BASE64
	Type         bool   // TYPE
	Root         bool   // ROOT [('name')]
	RootName     string // optional name for ROOT('name')

	// XML-specific options
	Elements       bool   // ELEMENTS
	ElementsMode   string // "" (default), "XSINIL", "ABSENT"
	XmlData        bool   // XMLDATA
	XmlSchema      bool   // XMLSCHEMA [('TargetNameSpaceURI')]
	XmlSchemaURI   string // optional target namespace URI

	// JSON-specific options
	IncludeNullValues   bool // INCLUDE_NULL_VALUES
	WithoutArrayWrapper bool // WITHOUT_ARRAY_WRAPPER

	Loc Loc
}

func (n *ForClause) nodeTag() {}

// ForMode enumerates FOR clause modes.
type ForMode int

const (
	ForXML ForMode = iota
	ForJSON
	ForBrowse
)

// WithClause represents a WITH (CTE) clause.
type WithClause struct {
	XmlNamespaces *List // optional XMLNAMESPACES declarations
	CTEs          *List // list of CommonTableExpr
	Loc           Loc
}

// XmlNamespaceDecl represents a single XMLNAMESPACES declaration.
//
//	uri AS prefix | DEFAULT uri
type XmlNamespaceDecl struct {
	URI       string // the namespace URI
	Prefix    string // AS prefix (empty for DEFAULT)
	IsDefault bool   // DEFAULT uri
	Loc       Loc
}

func (n *XmlNamespaceDecl) nodeTag() {}

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
	OptionClause *List // OPTION ( <query_hint> [,...n] )
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
	OptionClause *List // OPTION ( <query_hint> [,...n] )
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
	OptionClause *List // OPTION ( <query_hint> [,...n] )
	Loc          Loc
}

func (n *DeleteStmt) nodeTag()  {}
func (n *DeleteStmt) stmtNode() {}

// MergeStmt represents a MERGE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/merge-transact-sql
type MergeStmt struct {
	WithClause   *WithClause
	Top          *TopClause
	Target       *TableRef
	Source       TableExpr // table source
	SourceAlias  string
	OnCondition  ExprNode
	WhenClauses  *List // list of MergeWhenClause
	OutputClause *OutputClause
	OptionClause *List // OPTION ( <query_hint> [,...n] )
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
	Cols          *List // column list
	Values        *List // VALUES list
	DefaultValues bool  // DEFAULT VALUES
	Loc           Loc
}

func (n *MergeInsertAction) nodeTag() {}

// CurrentOfExpr represents WHERE CURRENT OF cursor_name in UPDATE/DELETE.
type CurrentOfExpr struct {
	CursorName string // cursor name or @cursor_variable
	Global     bool   // GLOBAL cursor
	Loc        Loc
}

func (n *CurrentOfExpr) nodeTag()  {}
func (n *CurrentOfExpr) exprNode() {}

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
	IsNode      bool // AS NODE (graph table)
	IsEdge      bool // AS EDGE (graph table)
	IsFileTable bool // AS FILETABLE

	// PERIOD FOR SYSTEM_TIME (start_col, end_col)
	PeriodStartCol string
	PeriodEndCol   string

	// Table storage options: ON filegroup, TEXTIMAGE_ON, FILESTREAM_ON
	OnFilegroup        string
	TextImageOn        string
	FilestreamOn       string

	// Inline indexes (INDEX ... inside CREATE TABLE body)
	Indexes *List // list of *InlineIndexDef

	// WITH table options
	TableOptions *List // list of *TableOption

	Loc Loc
}

func (n *CreateTableStmt) nodeTag()  {}
func (n *CreateTableStmt) stmtNode() {}

// InlineIndexDef represents an inline INDEX definition inside a CREATE TABLE body.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
//
//	INDEX index_name [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ]
//	    ( column_name [ ASC | DESC ] [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
//	    [ WHERE <filter_predicate> ]
//	    [ WITH ( <index_option> [ ,...n ] ) ]
type InlineIndexDef struct {
	Name         string
	Unique       bool
	Clustered    *bool    // true=CLUSTERED, false=NONCLUSTERED, nil=unspecified
	Columnstore  bool     // COLUMNSTORE index
	Columns      *List    // list of *IndexColumn
	IncludeCols  *List    // INCLUDE columns
	WhereClause  ExprNode // filtered index
	Options      *List    // WITH options
	OnFilegroup  string   // ON filegroup or partition scheme
	FilestreamOn string   // FILESTREAM_ON filegroup
	Loc          Loc
}

func (n *InlineIndexDef) nodeTag() {}

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

	// Advanced column options
	Sparse      bool // SPARSE
	Filestream  bool // FILESTREAM
	Rowguidcol  bool // ROWGUIDCOL
	Hidden      bool // HIDDEN
	IsColumnSet bool // XML COLUMN_SET FOR ALL_SPARSE_COLUMNS

	// MASKED WITH (FUNCTION = 'mask_function')
	MaskFunction string

	// ENCRYPTED WITH (COLUMN_ENCRYPTION_KEY = key, ENCRYPTION_TYPE = type, ALGORITHM = alg)
	EncryptedWith *EncryptedWithSpec

	// GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END } [ HIDDEN ]
	GeneratedAlways *GeneratedAlwaysSpec

	// NOT FOR REPLICATION
	NotForReplication bool

	Loc Loc
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
	NotNull   bool // NOT NULL after PERSISTED
	Loc       Loc
}

func (n *ComputedColumnDef) nodeTag() {}

// EncryptedWithSpec represents ENCRYPTED WITH options for Always Encrypted.
type EncryptedWithSpec struct {
	ColumnEncryptionKey string // COLUMN_ENCRYPTION_KEY = key_name
	EncryptionType      string // DETERMINISTIC | RANDOMIZED
	Algorithm           string // AEAD_AES_256_CBC_HMAC_SHA_256
	Loc                 Loc
}

func (n *EncryptedWithSpec) nodeTag() {}

// GeneratedAlwaysSpec represents GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END }.
type GeneratedAlwaysSpec struct {
	Kind    string // "ROW", "TRANSACTION_ID", "SEQUENCE_NUMBER"
	StartEnd string // "START" or "END"
	Loc     Loc
}

func (n *GeneratedAlwaysSpec) nodeTag() {}

// TableOption represents a single option in WITH (...) for CREATE TABLE.
type TableOption struct {
	Name  string // e.g. MEMORY_OPTIMIZED, DURABILITY, SYSTEM_VERSIONING, DATA_COMPRESSION
	Value string // e.g. ON, OFF, SCHEMA_ONLY, ROW, PAGE

	// For SYSTEM_VERSIONING = ON (HISTORY_TABLE = schema.table, ...)
	HistoryTable          string
	DataConsistencyCheck  string // ON or OFF
	HistoryRetentionPeriod string

	Loc Loc
}

func (n *TableOption) nodeTag() {}

// ConstraintDef represents a table or column constraint.
type ConstraintDef struct {
	Type       ConstraintType
	Name       string   // constraint name (optional)
	Columns    *List    // columns for PK, UNIQUE, INDEX
	Expr       ExprNode // CHECK expression, default expression
	RefTable   *TableRef
	RefColumns *List // FK referenced columns
	OnDelete          ReferentialAction
	OnUpdate          ReferentialAction
	Clustered         *bool  // true=CLUSTERED, false=NONCLUSTERED, nil=unspecified
	NotForReplication bool   // NOT FOR REPLICATION (CHECK, FK)
	IndexOptions      *List  // WITH ( index_option [,...n] )
	Fillfactor        int    // WITH FILLFACTOR = N
	OnFilegroup       string // ON filegroup or partition_scheme(col)
	EdgeConnections   *List  // list of *EdgeConnectionDef for EDGE CONSTRAINT
	Loc               Loc
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
	ConstraintEdge // EDGE CONSTRAINT CONNECTION (...)
)

// EdgeConnectionDef represents a single edge constraint connection clause (FromTable TO ToTable).
type EdgeConnectionDef struct {
	FromTable *TableRef
	ToTable   *TableRef
	Loc       Loc
}

func (n *EdgeConnectionDef) nodeTag() {}

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
	Collation  string   // ALTER COLUMN ... COLLATE collation_name
	Options    *List    // generic options list (e.g., SET options, REBUILD options, WITH options)
	TargetName *TableRef // SWITCH ... TO target_table
	Names      *List    // constraint/trigger name list for CHECK/NOCHECK/ENABLE/DISABLE
	Partition  ExprNode // SWITCH PARTITION n / REBUILD PARTITION = n
	TargetPart ExprNode // SWITCH TO ... PARTITION n
	WithCheck  string   // "CHECK" or "NOCHECK" prefix for WITH { CHECK | NOCHECK }
	IfExists   bool     // DROP COLUMN IF EXISTS / DROP CONSTRAINT IF EXISTS
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
	ATCheckConstraint       // [WITH CHECK|NOCHECK] CHECK CONSTRAINT { ALL | name [,...n] }
	ATNocheckConstraint     // [WITH CHECK|NOCHECK] NOCHECK CONSTRAINT { ALL | name [,...n] }
	ATEnableTrigger         // ENABLE TRIGGER { ALL | name [,...n] }
	ATDisableTrigger        // DISABLE TRIGGER { ALL | name [,...n] }
	ATEnableChangeTracking  // ENABLE CHANGE_TRACKING [WITH (...)]
	ATDisableChangeTracking // DISABLE CHANGE_TRACKING
	ATSwitchPartition       // SWITCH [PARTITION n] TO target [PARTITION n]
	ATRebuild               // REBUILD [PARTITION = ALL|n] [WITH (...)]
	ATSet                   // SET (LOCK_ESCALATION = ..., FILESTREAM_ON = ..., SYSTEM_VERSIONING = ...)
	ATAlterColumnAddDrop          // ALTER COLUMN col {ADD|DROP} {ROWGUIDCOL|PERSISTED|...}
	ATEnableFiletableNamespace    // ENABLE FILETABLE_NAMESPACE
	ATDisableFiletableNamespace   // DISABLE FILETABLE_NAMESPACE
	ATAddPeriod                   // ADD PERIOD FOR SYSTEM_TIME (start_col, end_col)
	ATDropPeriod                  // DROP PERIOD FOR SYSTEM_TIME
	ATSplitRange                  // SPLIT RANGE (boundary_value)
	ATMergeRange                  // MERGE RANGE (boundary_value)
)

// AlterColumnOption represents a typed ADD/DROP option in ALTER COLUMN.
// Replaces nodes.String concatenations in parseAlterColumnAddDrop.
//
// Examples:
//   - Action="ADD", Option="ROWGUIDCOL"
//   - Action="DROP", Option="PERSISTED"
//   - Action="ADD", Option="MASKED", MaskFunction="default()"
//   - Action="DROP", Option="NOT FOR REPLICATION"
type AlterColumnOption struct {
	Action       string // "ADD" or "DROP"
	Option       string // e.g., "ROWGUIDCOL", "PERSISTED", "SPARSE", "HIDDEN", "MASKED", "NOT FOR REPLICATION"
	MaskFunction string // only for MASKED WITH (FUNCTION = '...')
	Loc          Loc
}

func (n *AlterColumnOption) nodeTag() {}

// DropStmt represents a DROP statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-table-transact-sql
type DropStmt struct {
	ObjectType   DropObjectType
	Names        *List // list of TableRef
	IfExists     bool
	Options      *List // WITH options (e.g., DROP INDEX ... WITH (MAXDOP=1, ONLINE=ON))
	OnDatabase   bool  // ON DATABASE (DROP TRIGGER for DDL triggers)
	OnAllServer  bool  // ON ALL SERVER (DROP TRIGGER for DDL/logon triggers)
	NoDependents bool  // WITH NO DEPENDENTS (DROP ASSEMBLY)
	Loc          Loc
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
	DropSequence
	DropSynonym
	DropAssembly
	DropPartitionFunction
	DropPartitionScheme
	DropStatistics
	DropDefault
	DropRule
	DropXmlSchemaCollection
	DropFulltextIndex
	DropFulltextCatalog
	DropMaterializedView
)

// CreateIndexStmt represents a CREATE INDEX statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-index-transact-sql
type CreateIndexStmt struct {
	Unique       bool
	Clustered    *bool // true=CLUSTERED, false=NONCLUSTERED
	Columnstore  bool
	Name         string
	Table        *TableRef
	Columns      *List    // index columns
	IncludeCols  *List    // INCLUDE columns
	OrderCols    *List    // ORDER columns (columnstore only)
	WhereClause  ExprNode // filtered index
	Options      *List    // WITH options
	OnFileGroup  string
	FilestreamOn string // FILESTREAM_ON filegroup
	Loc          Loc
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
	Options       *List    // WITH options (SCHEMABINDING, VIEW_METADATA, ENCRYPTION)
	Loc           Loc
}

func (n *CreateViewStmt) nodeTag()  {}
func (n *CreateViewStmt) stmtNode() {}

// CreateTriggerStmt represents a CREATE TRIGGER statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-trigger-transact-sql
//
// DML trigger:
//
//	CREATE [ OR ALTER ] TRIGGER [ schema_name . ] trigger_name
//	ON { table | view }
//	[ WITH <dml_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER | INSTEAD OF }
//	{ [ INSERT ] [ , ] [ UPDATE ] [ , ] [ DELETE ] }
//	[ WITH APPEND ]
//	[ NOT FOR REPLICATION ]
//	AS { sql_statement [ ; ] [ , ...n ] }
//
// DDL trigger:
//
//	CREATE [ OR ALTER ] TRIGGER trigger_name
//	ON { ALL SERVER | DATABASE }
//	[ WITH <ddl_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER } { event_type | event_group } [ , ...n ]
//	AS { sql_statement [ ; ] [ , ...n ] }
//
// Logon trigger:
//
//	CREATE [ OR ALTER ] TRIGGER trigger_name
//	ON ALL SERVER
//	[ WITH <logon_trigger_option> [ , ...n ] ]
//	{ FOR | AFTER } LOGON
//	AS { sql_statement [ ; ] [ , ...n ] }
type CreateTriggerStmt struct {
	OrAlter           bool
	Name              *TableRef // trigger name (possibly schema-qualified)
	Table             *TableRef // ON table/view (DML trigger, nil for DDL/Logon)
	OnDatabase        bool      // ON DATABASE (DDL trigger)
	OnAllServer       bool      // ON ALL SERVER (DDL/Logon trigger)
	TriggerOptions    *List     // WITH options: ENCRYPTION, EXECUTE AS, NATIVE_COMPILATION, SCHEMABINDING
	TriggerType       string    // "FOR", "AFTER", "INSTEAD OF"
	Events            *List     // list of String: INSERT/UPDATE/DELETE (DML) or event types (DDL)
	WithAppend        bool      // WITH APPEND
	NotForReplication bool      // NOT FOR REPLICATION
	ExternalName      string    // EXTERNAL NAME assembly_name.class_name.method_name (CLR trigger)
	Body              Node      // statement body (BeginEndStmt or single stmt)
	Loc               Loc
}

func (n *CreateTriggerStmt) nodeTag()  {}
func (n *CreateTriggerStmt) stmtNode() {}

// TriggerEvent represents a typed event in a trigger event list.
// Replaces nodes.String for INSERT, UPDATE, DELETE, LOGON, and DDL event types.
type TriggerEvent struct {
	Name string // event name: "INSERT", "UPDATE", "DELETE", "LOGON", or DDL event type
	Loc  Loc
}

func (n *TriggerEvent) nodeTag() {}

// TriggerOption represents a typed WITH option in a CREATE TRIGGER statement.
// Replaces nodes.String for ENCRYPTION, EXECUTE AS, NATIVE_COMPILATION, SCHEMABINDING.
//
// Examples:
//   - Name="ENCRYPTION"
//   - Name="EXECUTE AS", Value="CALLER"
//   - Name="NATIVE_COMPILATION"
//   - Name="SCHEMABINDING"
type TriggerOption struct {
	Name  string // option name
	Value string // option value (for EXECUTE AS)
	Loc   Loc
}

func (n *TriggerOption) nodeTag() {}

// EnableDisableTriggerStmt represents ENABLE TRIGGER or DISABLE TRIGGER.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/enable-trigger-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/disable-trigger-transact-sql
//
//	{ ENABLE | DISABLE } TRIGGER { [ schema_name . ] trigger_name [ , ...n ] | ALL }
//	    ON { object_name | DATABASE | ALL SERVER }
type EnableDisableTriggerStmt struct {
	Enable      bool      // true for ENABLE, false for DISABLE
	TriggerAll  bool      // ALL triggers
	Triggers    *List     // list of trigger names (as String nodes)
	OnObject    *TableRef // ON table/view name (nil for DATABASE/ALL SERVER)
	OnDatabase  bool      // ON DATABASE
	OnAllServer bool      // ON ALL SERVER
	Loc         Loc
}

func (n *EnableDisableTriggerStmt) nodeTag()  {}
func (n *EnableDisableTriggerStmt) stmtNode() {}

// BulkInsertStmt represents a BULK INSERT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/bulk-insert-transact-sql
//
//	BULK INSERT [ database_name . [ schema_name ] . | schema_name . ] table_name
//	    FROM 'data_file'
//	    [ WITH ( option [,...n] ) ]
type BulkInsertStmt struct {
	Table    *TableRef
	DataFile string // 'file_path'
	Options  *List  // WITH options as key=value or flag strings
	Loc      Loc
}

func (n *BulkInsertStmt) nodeTag()  {}
func (n *BulkInsertStmt) stmtNode() {}

// CreateFunctionStmt represents a CREATE FUNCTION statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-function-transact-sql
type CreateFunctionStmt struct {
	OrAlter      bool
	Name         *TableRef
	Params       *List // list of ParamDef
	ReturnType   *DataType
	ReturnsTable *ReturnsTableDef
	Body         Node   // BeginEndStmt or single expression
	Options      *List  // WITH options
	ExternalName string // EXTERNAL NAME assembly.class.method (CLR)
	OrderClause  *List  // ORDER (...) for CLR table-valued functions
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
	Varying  bool     // VARYING keyword (for cursor parameters)
	Null     bool     // NULL keyword
	Default  ExprNode // = default_value
	Output   bool     // OUTPUT keyword
	ReadOnly bool     // READONLY keyword
	Loc      Loc
}

func (n *ParamDef) nodeTag() {}

// CreateProcedureStmt represents a CREATE PROCEDURE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-procedure-transact-sql
type CreateProcedureStmt struct {
	OrAlter        bool
	Name           *TableRef
	Number         int    // procedure number (;number)
	Params         *List  // list of ParamDef
	Body           Node   // BeginEndStmt or EXTERNAL NAME
	Options        *List  // WITH options
	ForReplication bool   // FOR REPLICATION
	ExternalName   string // EXTERNAL NAME assembly.class.method (CLR)
	Loc            Loc
}

func (n *CreateProcedureStmt) nodeTag()  {}
func (n *CreateProcedureStmt) stmtNode() {}

// CreateDatabaseStmt represents a CREATE DATABASE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-database-transact-sql
type CreateDatabaseStmt struct {
	Name          string
	Containment   string          // NONE | PARTIAL (empty if not specified)
	OnPrimary     *List           // list of *DatabaseFileSpec for PRIMARY filegroup
	Filegroups    *List           // list of *DatabaseFilegroup for additional filegroups
	LogOn         *List           // list of *DatabaseFileSpec for LOG ON
	Collation     string          // COLLATE collation_name
	WithOptions   *List           // list of *DatabaseOption for WITH options
	ForAttach     bool            // FOR ATTACH
	AttachOptions *List           // WITH options for ATTACH (as *DatabaseOption nodes)
	ForAttachRebuildLog bool      // FOR ATTACH_REBUILD_LOG
	SnapshotOf    string          // AS SNAPSHOT OF source_database
	Options       *List           // legacy options field (kept for backward compat)
	Loc           Loc
}

func (n *CreateDatabaseStmt) nodeTag()  {}
func (n *CreateDatabaseStmt) stmtNode() {}

// DatabaseOption represents a structured option in CREATE DATABASE WITH or FOR ATTACH WITH.
//
// Simple options: DB_CHAINING ON/OFF, TRUSTWORTHY ON/OFF, LEDGER ON/OFF,
//
//	NESTED_TRIGGERS ON/OFF, TRANSFORM_NOISE_WORDS ON/OFF
//
// Key=value options: DEFAULT_FULLTEXT_LANGUAGE, DEFAULT_LANGUAGE, TWO_DIGIT_YEAR_CUTOFF,
//
//	CATALOG_COLLATION
//
// Flag options (attach): ENABLE_BROKER, NEW_BROKER, ERROR_BROKER_CONVERSATIONS, RESTRICTED_USER
//
// FILESTREAM: sub-options NON_TRANSACTED_ACCESS = OFF|READ_ONLY|FULL, DIRECTORY_NAME = 'name'
//
// PERSISTENT_LOG_BUFFER: = ON with DIRECTORY_NAME sub-option
type DatabaseOption struct {
	Name  string // option name
	Value string // option value (ON, OFF, identifier, number, string literal)

	// For FILESTREAM option
	FilestreamAccess   string // NON_TRANSACTED_ACCESS value: OFF, READ_ONLY, FULL
	FilestreamDirName  string // DIRECTORY_NAME value

	// For PERSISTENT_LOG_BUFFER = ON (DIRECTORY_NAME = 'path')
	PersistentLogDir string // DIRECTORY_NAME value

	Loc Loc
}

func (n *DatabaseOption) nodeTag() {}

// SizeValue represents a structured size value with numeric part and optional unit.
//
//	size [KB|MB|GB|TB|%]
type SizeValue struct {
	Value string // numeric value (e.g. "10", "100", "5")
	Unit  string // unit suffix: "KB", "MB", "GB", "TB", "%", or "" for bare number
	Loc   Loc
}

func (n *SizeValue) nodeTag() {}

// DatabaseFileSpec represents a file specification in CREATE DATABASE.
//
//	( NAME = logical_file_name, FILENAME = 'os_file_name'
//	  [, SIZE = size [KB|MB|GB|TB]]
//	  [, MAXSIZE = { max_size [KB|MB|GB|TB] | UNLIMITED }]
//	  [, FILEGROWTH = growth_increment [KB|MB|GB|TB|%]] )
type DatabaseFileSpec struct {
	Name       string     // logical file name
	NewName    string     // new logical name (NEWNAME, for ALTER DATABASE MODIFY FILE)
	Filename   string     // OS file path
	Size       *SizeValue // e.g. 10 MB
	MaxSize    *SizeValue // e.g. 100 MB or nil (check MaxSizeUnlimited)
	MaxSizeUnlimited bool // MAXSIZE = UNLIMITED
	FileGrowth *SizeValue // e.g. 5 MB or 10 %
	Offline    bool       // OFFLINE flag (for ALTER DATABASE MODIFY FILE)
	Loc        Loc
}

func (n *DatabaseFileSpec) nodeTag() {}

// DatabaseFilegroup represents a FILEGROUP clause in CREATE DATABASE.
//
//	FILEGROUP filegroup_name [CONTAINS FILESTREAM] [DEFAULT]
//	    | CONTAINS MEMORY_OPTIMIZED_DATA
//	    <filespec> [, ...n]
type DatabaseFilegroup struct {
	Name              string
	ContainsFilestream bool
	ContainsMemoryOptimized bool
	IsDefault         bool
	Files             *List // list of *DatabaseFileSpec
	Loc               Loc
}

func (n *DatabaseFilegroup) nodeTag() {}

// AlterDatabaseStmt represents an ALTER DATABASE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-database-transact-sql
//
//	ALTER DATABASE { database_name | CURRENT }
//	{
//	    SET <option_spec> [ ,...n ] [ WITH <termination> ]
//	  | ADD FILE <filespec> [ ,...n ] [ TO FILEGROUP filegroup_name ]
//	  | ADD LOG FILE <filespec> [ ,...n ]
//	  | REMOVE FILE logical_file_name
//	  | MODIFY FILE <filespec>
//	  | ADD FILEGROUP filegroup_name [ CONTAINS FILESTREAM | CONTAINS MEMORY_OPTIMIZED_DATA ]
//	  | REMOVE FILEGROUP filegroup_name
//	  | MODIFY FILEGROUP filegroup_name { <filegroup_updatability_option> | DEFAULT | NAME = new_name | AUTOGROW_SINGLE_FILE | AUTOGROW_ALL_FILES }
//	  | COLLATE collation_name
//	  | MODIFY NAME = new_database_name
//	}
type AlterDatabaseStmt struct {
	Name        string // database name or "CURRENT"
	Action      string // SET, ADD, REMOVE, MODIFY, COLLATE
	SubAction   string // FILE, LOG FILE, FILEGROUP, NAME (qualifier for ADD/REMOVE/MODIFY)
	Options     *List  // SET options as String nodes; or MODIFY FILEGROUP options
	FileSpecs   *List  // file specifications (DatabaseFileSpec nodes)
	TargetName  string // file/filegroup/collation name (context-dependent)
	NewName     string // for MODIFY NAME or MODIFY FILEGROUP NAME
	Termination string // WITH termination clause (e.g. "ROLLBACK IMMEDIATE")
	Loc         Loc
}

func (n *AlterDatabaseStmt) nodeTag()  {}
func (n *AlterDatabaseStmt) stmtNode() {}

// AlterIndexStmt represents an ALTER INDEX statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-index-transact-sql
//
//	ALTER INDEX { index_name | ALL } ON <object>
//	    { REBUILD | REORGANIZE | DISABLE | SET ( ... ) }
type AlterIndexStmt struct {
	IndexName string    // index name or "ALL"
	Table     *TableRef // ON table_name
	Action    string    // REBUILD, REORGANIZE, DISABLE, SET, RESUME, PAUSE, ABORT
	Partition string    // partition number or "ALL" (empty if not specified)
	Options   *List     // WITH options as String nodes (key=value or key)
	Loc       Loc
}

func (n *AlterIndexStmt) nodeTag()  {}
func (n *AlterIndexStmt) stmtNode() {}

// TruncateStmt represents a TRUNCATE TABLE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/truncate-table-transact-sql
type TruncateStmt struct {
	Table      *TableRef
	Partitions *List // WITH (PARTITIONS (range [,...n]))
	Loc        Loc
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
	Operator string   // "=" (default), "+=", "-=", "*=", "/=", "%=", "&=", "^=", "|="
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
//
// BNF: mssql/parser/bnf/waitfor-transact-sql.bnf
//
//	WAITFOR
//	{
//	    DELAY 'time_to_pass'
//	  | TIME 'time_to_execute'
//	  | [ ( receive_statement ) | ( get_conversation_group_statement ) ]
//	    [ , TIMEOUT timeout ]
//	}
type WaitForStmt struct {
	WaitType  string   // "DELAY" or "TIME" (empty for parenthesized statement form)
	Value     ExprNode // time expression for DELAY/TIME
	InnerStmt StmtNode // for WAITFOR ( receive_statement | get_conversation_group_statement )
	Timeout   ExprNode // TIMEOUT value
	Loc       Loc
}

func (n *WaitForStmt) nodeTag()  {}
func (n *WaitForStmt) stmtNode() {}

// ---------- Execution / utility statements ----------

// ExecStmt represents an EXEC/EXECUTE statement.
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/execute-transact-sql
type ExecStmt struct {
	Name         *TableRef
	Args         *List    // list of ExecArg
	ReturnVar    string   // @var = EXEC ...
	ExecString   ExprNode // EXEC ('string') dynamic SQL
	AsLogin      string   // AS LOGIN = 'name'
	AsUser       string   // AS USER = 'name'
	AtServer     string   // AT linked_server_name
	AtDataSource string   // AT DATA_SOURCE data_source_name
	WithOptions  *List    // WITH RECOMPILE, RESULT SETS ...
	Loc          Loc
}

func (n *ExecStmt) nodeTag()  {}
func (n *ExecStmt) stmtNode() {}

// ExecArg represents an argument in EXEC.
type ExecArg struct {
	Name      string // @param = value (named) or empty (positional)
	Value     ExprNode
	Output    bool // OUTPUT keyword
	IsDefault bool // DEFAULT keyword
	Loc       Loc
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
	Name            string
	WithMark        bool   // WITH MARK
	MarkDescription string // optional description after WITH MARK
	Loc             Loc
}

func (n *BeginTransStmt) nodeTag()  {}
func (n *BeginTransStmt) stmtNode() {}

// CommitTransStmt represents COMMIT TRANSACTION.
type CommitTransStmt struct {
	Name             string
	DelayedDurability string // OFF or ON, empty if not specified
	Loc              Loc
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

// SecurityStmt represents CREATE/ALTER/DROP USER, LOGIN, ROLE, APPLICATION ROLE,
// and ADD/DROP ROLE MEMBER statements.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-user-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-login-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-application-role-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-role-transact-sql
type SecurityStmt struct {
	Action      string   // CREATE, ALTER, DROP, ADD
	ObjectType  string   // USER, LOGIN, ROLE, APPLICATION ROLE
	Name        string   // principal name
	Options     *List    // list of *SecurityPrincipalOption nodes
	WhereClause ExprNode // WHERE predicate (for SERVER AUDIT)
	Loc         Loc
}

func (n *SecurityStmt) nodeTag()  {}
func (n *SecurityStmt) stmtNode() {}

// EndpointOption represents a typed key=value option in a CREATE/ALTER ENDPOINT statement.
// Replaces nodes.String concatenations with structured Name + Value pairs.
//
// Examples:
//   - Protocol: Name="AS", Value="TCP"
//   - State: Name="STATE", Value="STARTED"
//   - Payload: Name="FOR", Value="SERVICE_BROKER"
//   - Auth: Name="AUTHENTICATION", Value="WINDOWS KERBEROS CERTIFICATE MyCert"
//   - Encryption: Name="ENCRYPTION", Value="REQUIRED ALGORITHM AES"
//   - Role: Name="ROLE", Value="PARTNER"
//   - Listener: Name="LISTENER_PORT", Value="4022"
//   - Forwarding: Name="MESSAGE_FORWARDING", Value="ENABLED"
//   - Forward size: Name="MESSAGE_FORWARD_SIZE", Value="10"
//   - Authorization: Name="AUTHORIZATION", Value="sa"
//   - Flag options: Name="VERBOSE", Value=""
type EndpointOption struct {
	Name  string // option name (e.g., "AS", "FOR", "STATE", "AUTHENTICATION", "ENCRYPTION", "ROLE")
	Value string // option value (may be empty for flag-like options)
	Loc   Loc
}

func (n *EndpointOption) nodeTag() {}

// AvailabilityGroupOption represents a typed option in a CREATE/ALTER AVAILABILITY GROUP statement.
// Replaces nodes.String concatenations with structured Name + Value pairs.
//
// Examples:
//   - Action keywords: Name="WITH", Value=""
//   - Key=value pairs: Name="ENDPOINT_URL", Value="'TCP://server1:5022'"
//   - Actions with args: Name="ADD DATABASE", Value="MyDB"
//   - Nested options: Name="SECONDARY_ROLE", Value="(ALLOW_CONNECTIONS=READ_ONLY)"
//   - IP tuples: Name="", Value="('10.0.0.1', '255.255.255.0')"
type AvailabilityGroupOption struct {
	Name  string // option name (e.g., WITH, SET, REPLICA ON, ENDPOINT_URL)
	Value string // option value (may be empty for flag-like options)
	Loc   Loc
}

func (n *AvailabilityGroupOption) nodeTag() {}

// AuditSpecAction represents a single ADD or DROP action in an AUDIT SPECIFICATION.
//
// For server audit specs: ADD ( audit_action_group_name )
// For database audit specs: ADD ( action [,...n] ON [class ::] securable BY principal [,...n] )
type AuditSpecAction struct {
	Action     string   // "ADD" or "DROP"
	GroupName  string   // audit_action_group_name (e.g., "FAILED_LOGIN_GROUP") — used when no ON/BY
	Actions    []string // action names (e.g., ["SELECT", "INSERT"]) for database audit specs
	ClassName  string   // "OBJECT", "SCHEMA", "DATABASE" etc. (before ::)
	Securable  string   // fully qualified securable name (after :: or after ON without class)
	Principals []string // principal names (after BY)
	Loc        Loc
}

func (n *AuditSpecAction) nodeTag() {}

// SecurityPrincipalOption represents a single structured option for security principal statements.
//
// Flag options: MUST_CHANGE, HASHED, ENABLE, DISABLE, NO REVERT,
//
//	CHECK_EXPIRATION, CHECK_POLICY, ALLOW_ENCRYPTED_VALUE_MODIFICATIONS
//
// Key=value options: PASSWORD, DEFAULT_SCHEMA, DEFAULT_LANGUAGE, DEFAULT_DATABASE,
//
//	CREDENTIAL, SID, NAME, LOGIN, FROM, AUTHORIZATION,
//	ADD MEMBER, DROP MEMBER, COOKIE, OLD_PASSWORD
//
// ON/OFF options: CHECK_EXPIRATION, CHECK_POLICY
type SecurityPrincipalOption struct {
	Name  string // option name
	Value string // option value (may be empty for flags)

	// PASSWORD sub-options
	MustChange  bool   // MUST_CHANGE
	Hashed      bool   // HASHED
	Unlock      bool   // UNLOCK (ALTER LOGIN)
	OldPassword string // OLD_PASSWORD = 'old' (ALTER LOGIN)

	Loc Loc
}

func (n *SecurityPrincipalOption) nodeTag() {}

// EventNotificationOption represents a structured option for EVENT NOTIFICATION statements.
//
// Scope: SERVER, DATABASE, or QUEUE (with QueueName).
// FanIn: true if WITH FAN_IN was specified.
// Events: list of event type/group names from FOR clause.
// ServiceName: broker service name from TO SERVICE clause.
// BrokerInstance: broker instance specifier from TO SERVICE clause.
// ExtraNames: additional notification names for DROP ... name1, name2.
type EventNotificationOption struct {
	Scope          string   // "SERVER", "DATABASE", "QUEUE"
	QueueName      string   // queue name (only when Scope == "QUEUE")
	FanIn          bool     // WITH FAN_IN
	Events         []string // FOR event_type/event_group names
	ServiceName    string   // TO SERVICE 'broker_service'
	BrokerInstance string   // broker_instance_specifier
	ExtraNames     []string // additional names (for DROP with multiple names)
	Loc            Loc
}

func (n *EventNotificationOption) nodeTag() {}

// ResourceGovernorOption represents a structured non-WITH option for Resource Governor statements.
//
// Used for outer-loop options like:
//   - USING pool_name
//   - EXTERNAL ext_pool_name
//   - RECONFIGURE, DISABLE, RESET STATISTICS
//   - key = value pairs outside WITH
type ResourceGovernorOption struct {
	Name  string // option keyword (e.g., "USING", "EXTERNAL", "RECONFIGURE", "DISABLE", "RESET")
	Value string // option value (e.g., pool name, "STATISTICS", qualified name)
	Loc   Loc
}

func (n *ResourceGovernorOption) nodeTag() {}

// ExternalOption represents a structured key=value option for external object statements.
//
// Used by ALTER EXTERNAL DATA SOURCE SET options and similar.
type ExternalOption struct {
	Key   string // option key (e.g., "LOCATION", "CREDENTIAL")
	Value string // option value
	Loc   Loc
}

func (n *ExternalOption) nodeTag() {}

// CreateSchemaStmt represents CREATE SCHEMA.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-schema-transact-sql
type CreateSchemaStmt struct {
	Name          string // schema name (may be empty if AUTHORIZATION only)
	Authorization string // AUTHORIZATION owner_name (may be empty)
	Elements      *List  // optional schema elements (CREATE/GRANT/REVOKE/DENY)
	Loc           Loc
}

func (n *CreateSchemaStmt) nodeTag()  {}
func (n *CreateSchemaStmt) stmtNode() {}

// AlterSchemaStmt represents ALTER SCHEMA.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-schema-transact-sql
type AlterSchemaStmt struct {
	Name           string // schema name
	TransferType   string // OBJECT, TYPE, XML SCHEMA COLLECTION, or ""
	TransferEntity string // dot-qualified entity name
	Loc            Loc
}

func (n *AlterSchemaStmt) nodeTag()  {}
func (n *AlterSchemaStmt) stmtNode() {}

// CreateTypeStmt represents CREATE TYPE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-type-transact-sql
type CreateTypeStmt struct {
	Name            *TableRef // [schema.]type_name
	BaseType        *DataType // FROM base_type (alias type)
	Nullable        *bool     // NULL / NOT NULL for alias type
	ExternalName    string    // EXTERNAL NAME assembly.class
	TableDef        *List     // AS TABLE (...) column/constraint definitions
	MemoryOptimized bool      // WITH (MEMORY_OPTIMIZED = ON)
	Loc             Loc
}

func (n *CreateTypeStmt) nodeTag()  {}
func (n *CreateTypeStmt) stmtNode() {}

// TableTypeIndex represents an INDEX clause within CREATE TYPE AS TABLE.
//
// BNF:
//
//	INDEX index_name [ CLUSTERED | NONCLUSTERED ]
//	    [ HASH WITH ( BUCKET_COUNT = count ) ]
//	    ( column_name [ ASC | DESC ] [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
type TableTypeIndex struct {
	Name        string // index name
	Clustered   *bool  // true=CLUSTERED, false=NONCLUSTERED, nil=unspecified
	Hash        bool   // HASH index (memory-optimized)
	BucketCount ExprNode // BUCKET_COUNT for hash index
	Columns     *List  // IndexColumn list
	IncludeCols *List  // INCLUDE columns (string names)
	Loc         Loc
}

func (n *TableTypeIndex) nodeTag() {}

// CreateSequenceStmt represents CREATE SEQUENCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-sequence-transact-sql
type CreateSequenceStmt struct {
	Name        *TableRef // [schema.]sequence_name
	DataType    *DataType // AS integer_type (optional)
	Start       ExprNode  // START WITH constant
	Restart     bool      // RESTART (used internally for ALTER SEQUENCE option parsing)
	RestartWith ExprNode  // RESTART WITH constant (used internally for ALTER SEQUENCE)
	Increment   ExprNode  // INCREMENT BY constant
	MinValue    ExprNode  // MINVALUE constant
	MaxValue    ExprNode  // MAXVALUE constant
	NoMinVal    bool      // NO MINVALUE
	NoMaxVal    bool      // NO MAXVALUE
	Cycle       *bool     // CYCLE (true) / NO CYCLE (false) / nil (unset)
	Cache       ExprNode  // CACHE n
	NoCache     bool      // NO CACHE
	Loc         Loc
}

func (n *CreateSequenceStmt) nodeTag()  {}
func (n *CreateSequenceStmt) stmtNode() {}

// AlterSequenceStmt represents ALTER SEQUENCE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-sequence-transact-sql
type AlterSequenceStmt struct {
	Name        *TableRef // [schema.]sequence_name
	Restart     bool      // RESTART
	RestartWith ExprNode  // RESTART WITH constant
	Increment   ExprNode  // INCREMENT BY constant
	MinValue    ExprNode  // MINVALUE constant
	MaxValue    ExprNode  // MAXVALUE constant
	NoMinVal    bool      // NO MINVALUE
	NoMaxVal    bool      // NO MAXVALUE
	Cycle       *bool     // CYCLE / NO CYCLE
	Cache       ExprNode  // CACHE n
	NoCache     bool      // NO CACHE
	Loc         Loc
}

func (n *AlterSequenceStmt) nodeTag()  {}
func (n *AlterSequenceStmt) stmtNode() {}

// CreateSynonymStmt represents CREATE SYNONYM.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-synonym-transact-sql
type CreateSynonymStmt struct {
	Name   *TableRef // [schema.]synonym_name
	Target *TableRef // FOR [server.][database.][schema.]object
	Loc    Loc
}

func (n *CreateSynonymStmt) nodeTag()  {}
func (n *CreateSynonymStmt) stmtNode() {}

// GrantStmt represents GRANT/REVOKE/DENY.
type GrantStmt struct {
	StmtType       GrantType // GRANT, REVOKE, DENY
	GrantOptionFor bool      // REVOKE GRANT OPTION FOR ...
	Privileges     *List
	OnType         string    // securable class: SCHEMA, OBJECT, DATABASE, LOGIN, etc.
	OnName         *TableRef
	Principals     *List  // TO/FROM principals
	WithGrant      bool   // WITH GRANT OPTION
	AsPrincipal    string // AS principal
	CascadeOpt     bool   // CASCADE
	Loc            Loc
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

// SubqueryComparisonExpr represents expr comparison_op { ANY | SOME | ALL } (subquery).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/some-any-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/all-transact-sql
//
//	scalar_expression { = | <> | != | > | >= | !> | < | <= | !< }
//	    { ALL | SOME | ANY } ( subquery )
type SubqueryComparisonExpr struct {
	Left     ExprNode // left-hand scalar expression
	Op       BinaryOp // comparison operator
	Quantifier string // "ALL", "SOME", or "ANY"
	Subquery StmtNode // subquery (SelectStmt)
	Loc      Loc
}

func (n *SubqueryComparisonExpr) nodeTag()  {}
func (n *SubqueryComparisonExpr) exprNode() {}

// CollateExpr represents a postfix COLLATE operator on an expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/collations
//
//	expr COLLATE { collation_name | database_default }
type CollateExpr struct {
	Expr      ExprNode // the expression being collated
	Collation string   // collation name or "database_default"
	Loc       Loc
}

func (n *CollateExpr) nodeTag()  {}
func (n *CollateExpr) exprNode() {}

// AtTimeZoneExpr represents the AT TIME ZONE postfix expression.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/at-time-zone-transact-sql
//
//	inputdate AT TIME ZONE timezone
type AtTimeZoneExpr struct {
	Expr     ExprNode // the datetime expression
	TimeZone ExprNode // the timezone expression (usually a string literal)
	Loc      Loc
}

func (n *AtTimeZoneExpr) nodeTag()  {}
func (n *AtTimeZoneExpr) exprNode() {}

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
	Hints    *List // table hints: WITH (NOLOCK), WITH (INDEX(...)), etc.
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
	Table       TableExpr          // TableRef, SubqueryExpr, etc.
	Alias       string
	Columns     *List              // alias column list
	TableSample *TableSampleClause // optional TABLESAMPLE clause
	Hints       *List              // table hints: WITH (NOLOCK), etc.
	Loc         Loc
}

func (n *AliasedTableRef) nodeTag()   {}
func (n *AliasedTableRef) tableExpr() {}

// ---------- Window / OVER clause ----------

// WindowDef represents a named window definition in a WINDOW clause.
//
//	window_name AS ( window_specification )
type WindowDef struct {
	Name        string     // window name
	PartitionBy *List      // PARTITION BY expressions
	OrderBy     *List      // ORDER BY items
	Frame       *WindowFrame
	RefName     string // reference to another named window
	Loc         Loc
}

func (n *WindowDef) nodeTag() {}

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
	Column      *ColumnRef
	Variable    string // @var = expr
	VarColumn   *ColumnRef // @variable = column = expression (dual assignment)
	Operator    string     // "=" (default), "+=", "-=", "*=", "/=", "%=", "&=", "^=", "|="
	WriteMethod bool       // column.WRITE(expression, @Offset, @Length) form
	Value       ExprNode
	Loc         Loc
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

// ---------- DBCC ----------

// DbccStmt represents a DBCC (Database Console Command) statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/database-console-commands/dbcc-transact-sql
//
//	DBCC command_name [ ( arg [, ...] ) ] [ WITH option [, ...] ]
type DbccStmt struct {
	Command string // e.g. CHECKDB, SHRINKDATABASE, FREEPROCCACHE, etc.
	Args    *List  // optional arguments inside parentheses
	Options *List  // optional WITH options (list of *DbccOption nodes)
	Loc     Loc
}

func (n *DbccStmt) nodeTag()  {}
func (n *DbccStmt) stmtNode() {}

// DbccOption represents a single structured DBCC WITH option.
//
// Common options: NO_INFOMSGS, ALL_ERRORMSGS, PHYSICAL_ONLY,
//
//	EXTENDED_LOGICAL_CHECKS, DATA_PURITY, TABLOCK,
//	ESTIMATEONLY, COUNT_ROWS, TABLERESULTS
type DbccOption struct {
	Name string // option name (always uppercase)
	Loc  Loc
}

func (n *DbccOption) nodeTag() {}

// ---------- Backup / Restore ----------

// BackupStmt represents a BACKUP DATABASE or BACKUP LOG statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/backup-transact-sql
//
//	BACKUP { DATABASE | LOG } { database_name | @database_name_var }
//	  [ <file_or_filegroup> [ ,...n ] ]
//	  TO <backup_device> [ ,...n ]
//	  [ <MIRROR TO clause> ] [ next-mirror-to ]
//	  [ WITH { DIFFERENTIAL | <general_WITH_options> } [ ,...n ] ]
type BackupStmt struct {
	Type         string // "DATABASE", "LOG", or "CERTIFICATE"
	Database     string // database name
	Target       string // first TO device path value (backward compat)
	FileSpecs    *List  // FILE/FILEGROUP/READ_WRITE_FILEGROUPS specs (String nodes)
	Devices      *List  // all TO devices (String nodes: "TYPE=path" or "logical_name")
	MirrorTo     bool   // true if MIRROR TO clause present
	MirrorDevice string // first MIRROR TO device path
	Options      *List  // WITH options (as BackupRestoreOption nodes)
	Loc          Loc
}

func (n *BackupStmt) nodeTag()  {}
func (n *BackupStmt) stmtNode() {}

// RestoreStmt represents a RESTORE DATABASE / LOG / HEADERONLY / FILELISTONLY statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/restore-statements-transact-sql
//
//	RESTORE { DATABASE | LOG } { database_name | @database_name_var }
//	  [ <file_or_filegroup> [ ,...n ] ]
//	  [ FROM <backup_device> [ ,...n ] ]
//	  [ WITH options ]
//
//	RESTORE { HEADERONLY | FILELISTONLY | VERIFYONLY | LABELONLY | REWINDONLY }
//	  FROM <backup_device> [ ,...n ]
//	  [ WITH options ]
//
//	RESTORE DATABASE { database_name } FROM DATABASE_SNAPSHOT = snapshot_name
type RestoreStmt struct {
	Type         string // "DATABASE", "LOG", "HEADERONLY", "FILELISTONLY", "VERIFYONLY", "LABELONLY", "REWINDONLY"
	Database     string // database name (may be empty for HEADERONLY/FILELISTONLY)
	Source       string // first FROM device path value (backward compat)
	FileSpecs    *List  // FILE/FILEGROUP/READ_WRITE_FILEGROUPS/PAGE specs (String nodes)
	Devices      *List  // all FROM devices (String nodes)
	SnapshotName string // DATABASE_SNAPSHOT = name
	Options      *List  // WITH options (as BackupRestoreOption nodes)
	Loc          Loc
}

func (n *RestoreStmt) nodeTag()  {}
func (n *RestoreStmt) stmtNode() {}

// BackupRestoreOption represents a single structured option in a BACKUP/RESTORE WITH clause.
//
// Flag options: COMPRESSION, NO_COMPRESSION, DIFFERENTIAL, COPY_ONLY, INIT, NOINIT,
//
//	NOSKIP, SKIP, FORMAT, NOFORMAT, NO_CHECKSUM, CHECKSUM,
//	STOP_ON_ERROR, CONTINUE_AFTER_ERROR, RESTART, REPLACE,
//	RECOVERY, NORECOVERY, NO_TRUNCATE, FILE_SNAPSHOT,
//	ENABLE_BROKER, NEW_BROKER, ERROR_BROKER_CONVERSATIONS,
//	REWIND, NOREWIND, UNLOAD, NOUNLOAD,
//	RESTRICTED_USER, KEEP_REPLICATION, KEEP_CDC,
//	PARTIAL, CREDENTIAL, METADATA_ONLY, SNAPSHOT
//
// Key=value options: NAME, DESCRIPTION, EXPIREDATE, RETAINDAYS, STATS,
//
//	BLOCKSIZE, BUFFERCOUNT, MAXTRANSFERSIZE, MEDIADESCRIPTION,
//	MEDIANAME, MEDIAPASSWORD, STANDBY, STOPAT, STOPATMARK,
//	STOPBEFOREMARK, FILE, PASSWORD, DBNAME
//
// ENCRYPTION: ALGORITHM = alg, SERVER CERTIFICATE name | ASYMMETRIC KEY name
// MOVE: MOVE 'logical' TO 'physical'
// FILESTREAM: FILESTREAM ( DIRECTORY_NAME = directory_name )
type BackupRestoreOption struct {
	Name  string // option name
	Value string // value for key=value options

	// ENCRYPTION sub-options
	Algorithm     string // AES_128, AES_192, AES_256, TRIPLE_DES_3KEY
	EncryptorType string // SERVER CERTIFICATE, ASYMMETRIC KEY
	EncryptorName string // certificate/key name

	// MOVE sub-options (RESTORE)
	MoveFrom string // logical file name
	MoveTo   string // OS file name

	Loc Loc
}

func (n *BackupRestoreOption) nodeTag() {}

// ---------- Security keys/certs ----------

// SecurityKeyStmt represents a CREATE/ALTER/DROP/OPEN/CLOSE/BACKUP statement
// for security objects: MASTER KEY, SYMMETRIC KEY, ASYMMETRIC KEY, CERTIFICATE, CREDENTIAL.
type SecurityKeyStmt struct {
	Action     string // CREATE, ALTER, DROP, OPEN, CLOSE, BACKUP
	ObjectType string // MASTER KEY, SYMMETRIC KEY, ASYMMETRIC KEY, CERTIFICATE, CREDENTIAL
	Name       string // object name (may be empty for MASTER KEY)
	Options    *List  // generic list of String nodes for options/clauses
	Loc        Loc
}

func (n *SecurityKeyStmt) nodeTag()  {}
func (n *SecurityKeyStmt) stmtNode() {}

// ---------- Batch 39: BEGIN DISTRIBUTED TRANSACTION ----------

// BeginDistributedTransStmt represents BEGIN DISTRIBUTED TRAN[SACTION] [name].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/begin-distributed-transaction-transact-sql
type BeginDistributedTransStmt struct {
	Name string // optional transaction name
	Loc  Loc
}

func (n *BeginDistributedTransStmt) nodeTag()  {}
func (n *BeginDistributedTransStmt) stmtNode() {}

// ---------- Batch 40: CREATE/UPDATE/DROP STATISTICS ----------

// CreateStatisticsStmt represents CREATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-statistics-transact-sql
type CreateStatisticsStmt struct {
	Name    string     // statistics name
	Table   *TableRef  // table or indexed view
	Columns *List      // column name list
	Where   ExprNode   // WHERE filter predicate (optional)
	Options *List      // WITH options as String nodes
	Loc     Loc
}

func (n *CreateStatisticsStmt) nodeTag()  {}
func (n *CreateStatisticsStmt) stmtNode() {}

// UpdateStatisticsStmt represents UPDATE STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/update-statistics-transact-sql
type UpdateStatisticsStmt struct {
	Table   *TableRef // table or indexed view
	Name    string    // statistics name (optional, single name)
	Names   *List     // statistics/index names list (optional, parenthesized list)
	Options *List     // WITH options as String nodes
	Loc     Loc
}

func (n *UpdateStatisticsStmt) nodeTag()  {}
func (n *UpdateStatisticsStmt) stmtNode() {}

// DropStatisticsStmt represents DROP STATISTICS statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-statistics-transact-sql
type DropStatisticsStmt struct {
	// Each item is "table.stats_name" as a String node
	Names *List
	Loc   Loc
}

func (n *DropStatisticsStmt) nodeTag()  {}
func (n *DropStatisticsStmt) stmtNode() {}

// ---------- Batch 41: SET session options ----------

// SetOptionStmt represents SET session options like SET NOCOUNT ON/OFF,
// SET ANSI_NULLS ON/OFF, SET TRANSACTION ISOLATION LEVEL, etc.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/set-statements-transact-sql
type SetOptionStmt struct {
	Option string   // option name (e.g., "NOCOUNT", "ANSI_NULLS", "TRANSACTION ISOLATION LEVEL READ COMMITTED")
	Value  ExprNode // ON/OFF or other value; may be a ColumnRef("ON"/"OFF") or literal
	Loc    Loc
}

func (n *SetOptionStmt) nodeTag()  {}
func (n *SetOptionStmt) stmtNode() {}

// ---------- Batch 42: PARTITION FUNCTION/SCHEME ----------

// CreatePartitionFunctionStmt represents CREATE PARTITION FUNCTION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-partition-function-transact-sql
type CreatePartitionFunctionStmt struct {
	Name      string    // partition function name
	InputType *DataType // input parameter type
	Range     string    // LEFT or RIGHT
	Values    *List     // boundary values as expressions
	Loc       Loc
}

func (n *CreatePartitionFunctionStmt) nodeTag()  {}
func (n *CreatePartitionFunctionStmt) stmtNode() {}

// AlterPartitionFunctionStmt represents ALTER PARTITION FUNCTION ... SPLIT/MERGE RANGE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-partition-function-transact-sql
type AlterPartitionFunctionStmt struct {
	Name      string   // partition function name
	Action    string   // SPLIT or MERGE
	BoundaryValue ExprNode // boundary value
	Loc       Loc
}

func (n *AlterPartitionFunctionStmt) nodeTag()  {}
func (n *AlterPartitionFunctionStmt) stmtNode() {}

// CreatePartitionSchemeStmt represents CREATE PARTITION SCHEME.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-partition-scheme-transact-sql
type CreatePartitionSchemeStmt struct {
	Name            string // partition scheme name
	FunctionName    string // partition function name
	FileGroups      *List  // file group names as String nodes; "ALL" if single [ALL TO]
	AllToFileGroup  string // if ALL TO filegroup, stores the filegroup name
	Loc             Loc
}

func (n *CreatePartitionSchemeStmt) nodeTag()  {}
func (n *CreatePartitionSchemeStmt) stmtNode() {}

// AlterPartitionSchemeStmt represents ALTER PARTITION SCHEME ... NEXT USED filegroup.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-partition-scheme-transact-sql
type AlterPartitionSchemeStmt struct {
	Name      string // partition scheme name
	FileGroup string // NEXT USED filegroup name
	Loc       Loc
}

func (n *AlterPartitionSchemeStmt) nodeTag()  {}
func (n *AlterPartitionSchemeStmt) stmtNode() {}

// ---------- Batch 43: FULLTEXT ----------

// CreateFulltextIndexStmt represents CREATE FULLTEXT INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-fulltext-index-transact-sql
type CreateFulltextIndexStmt struct {
	Columns      *List  // fulltext index columns
	Table        *TableRef
	KeyIndex     string // unique index name
	CatalogName  string // fulltext catalog (optional)
	Options      *List  // WITH options as String nodes
	Loc          Loc
}

func (n *CreateFulltextIndexStmt) nodeTag()  {}
func (n *CreateFulltextIndexStmt) stmtNode() {}

// AlterFulltextIndexStmt represents ALTER FULLTEXT INDEX ON table action.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-fulltext-index-transact-sql
type AlterFulltextIndexStmt struct {
	Table            *TableRef
	Action           string // ENABLE, DISABLE, SET, ADD, ALTER, DROP, START, STOP, PAUSE, RESUME
	ChangeTracking   string // MANUAL, AUTO, OFF (for SET CHANGE_TRACKING)
	Columns          *List  // column names for ADD/DROP (String nodes)
	ColumnName       string // column name for ALTER COLUMN
	ColumnAction     string // ADD or DROP (for ALTER COLUMN ... STATISTICAL_SEMANTICS)
	PopulationType   string // FULL, INCREMENTAL, UPDATE (for START ... POPULATION)
	WithNoPopulation bool   // WITH NO POPULATION (for ADD/DROP/ALTER COLUMN)
	Options          *List  // additional options as String nodes
	Loc              Loc
}

func (n *AlterFulltextIndexStmt) nodeTag()  {}
func (n *AlterFulltextIndexStmt) stmtNode() {}

// CreateFulltextCatalogStmt represents CREATE FULLTEXT CATALOG.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-fulltext-catalog-transact-sql
type CreateFulltextCatalogStmt struct {
	Name    string
	Options *List // WITH options as String nodes
	Loc     Loc
}

func (n *CreateFulltextCatalogStmt) nodeTag()  {}
func (n *CreateFulltextCatalogStmt) stmtNode() {}

// AlterFulltextCatalogStmt represents ALTER FULLTEXT CATALOG name action.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-fulltext-catalog-transact-sql
type AlterFulltextCatalogStmt struct {
	Name    string
	Action  string // REBUILD, REORGANIZE, AS DEFAULT
	Options *List  // WITH options (e.g., ACCENT_SENSITIVITY=ON)
	Loc     Loc
}

func (n *AlterFulltextCatalogStmt) nodeTag()  {}
func (n *AlterFulltextCatalogStmt) stmtNode() {}

// ---------- Batch 44: XML SCHEMA COLLECTION ----------

// CreateXmlSchemaCollectionStmt represents CREATE XML SCHEMA COLLECTION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-xml-schema-collection-transact-sql
type CreateXmlSchemaCollectionStmt struct {
	Name           *TableRef // relational_schema.sql_identifier
	XmlSchemaNamespaces ExprNode  // xml_Schema_namespace expression
	Loc            Loc
}

func (n *CreateXmlSchemaCollectionStmt) nodeTag()  {}
func (n *CreateXmlSchemaCollectionStmt) stmtNode() {}

// AlterXmlSchemaCollectionStmt represents ALTER XML SCHEMA COLLECTION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-xml-schema-collection-transact-sql
type AlterXmlSchemaCollectionStmt struct {
	Name                *TableRef
	XmlSchemaNamespaces ExprNode
	Loc                 Loc
}

func (n *AlterXmlSchemaCollectionStmt) nodeTag()  {}
func (n *AlterXmlSchemaCollectionStmt) stmtNode() {}

// ---------- Batch 45: ASSEMBLY ----------

// CreateAssemblyStmt represents CREATE ASSEMBLY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-assembly-transact-sql
type CreateAssemblyStmt struct {
	Name          string
	Authorization string   // AUTHORIZATION owner_name
	FromFiles     *List    // file paths as String nodes
	PermissionSet string   // SAFE, EXTERNAL_ACCESS, UNSAFE
	Loc           Loc
}

func (n *CreateAssemblyStmt) nodeTag()  {}
func (n *CreateAssemblyStmt) stmtNode() {}

// AlterAssemblyStmt represents ALTER ASSEMBLY.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-assembly-transact-sql
type AlterAssemblyStmt struct {
	Name    string
	Actions *List // list of actions as String nodes
	Loc     Loc
}

func (n *AlterAssemblyStmt) nodeTag()  {}
func (n *AlterAssemblyStmt) stmtNode() {}

// AssemblyAction represents a typed action in an ALTER ASSEMBLY statement.
// Replaces nodes.String concatenations for FROM, PERMISSION_SET, DROP FILE, ADD FILE.
//
// Examples:
//   - Name="FROM", Value="C:\\path\\to\\assembly.dll"
//   - Name="PERMISSION_SET", Value="SAFE"
//   - Name="VISIBILITY", Value="ON"
//   - Name="DROP FILE"
//   - Name="ADD FILE"
type AssemblyAction struct {
	Name  string // action name: "FROM", "PERMISSION_SET", "VISIBILITY", "DROP FILE", "ADD FILE"
	Value string // action value (may be empty)
	Loc   Loc
}

func (n *AssemblyAction) nodeTag() {}

// AssemblyFile represents a file path in a CREATE ASSEMBLY FROM clause.
// Replaces nodes.String for file paths.
type AssemblyFile struct {
	Path string // file path
	Loc  Loc
}

func (n *AssemblyFile) nodeTag() {}

// ---------- Batch 46: SERVICE BROKER ----------

// ServiceBrokerStmt is a generic Service Broker statement node.
// It covers CREATE/ALTER/DROP MESSAGE TYPE, CONTRACT, QUEUE, SERVICE, ROUTE,
// SEND, RECEIVE, BEGIN/END CONVERSATION, GET CONVERSATION GROUP.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/service-broker-statements
type ServiceBrokerStmt struct {
	Action     string   // CREATE, ALTER, DROP, SEND, RECEIVE, BEGIN, END, GET
	ObjectType string   // MESSAGE TYPE, CONTRACT, QUEUE, SERVICE, ROUTE, CONVERSATION, etc.
	Name       string   // object name (may be empty for some forms)
	Options    *List    // options as String nodes
	Loc        Loc
}

func (n *ServiceBrokerStmt) nodeTag()  {}
func (n *ServiceBrokerStmt) stmtNode() {}

// ServiceBrokerOption represents a typed key=value option in a Service Broker statement.
type ServiceBrokerOption struct {
	Name  string // option name (e.g., STATUS, RETENTION, FROM SERVICE, ENCRYPTION)
	Value string // option value (may be empty for flag-like options such as CLEANUP, DROP)
	Loc   Loc
}

func (n *ServiceBrokerOption) nodeTag() {}

// ReceiveStmt represents a RECEIVE statement with structured fields.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/receive-transact-sql
//
//	RECEIVE [ TOP ( n ) ]
//	    <column_specifier> [ ,...n ]
//	    FROM <queue>
//	    [ INTO table_variable ]
//	    [ WHERE { conversation_handle = @handle | conversation_group_id = @group_id } ]
type ReceiveStmt struct {
	Top         ExprNode   // optional TOP (n) expression
	Columns     *List      // list of ReceiveColumn nodes, or nil for *
	AllColumns  bool       // true if RECEIVE *
	Queue       *TableRef  // FROM queue
	IntoVar     string     // INTO @table_variable
	WhereClause ExprNode   // WHERE condition
	Loc         Loc
}

func (n *ReceiveStmt) nodeTag()  {}
func (n *ReceiveStmt) stmtNode() {}

// ReceiveColumn represents a column entry in a RECEIVE statement.
//
//	{ column_name | expression } [ [ AS ] column_alias ]
type ReceiveColumn struct {
	Expr  ExprNode // column expression
	Alias string   // optional alias (with or without AS)
	Loc   Loc
}

func (n *ReceiveColumn) nodeTag() {}

// ---------- Batch 47: MISC UTILITY ----------

// CheckpointStmt represents CHECKPOINT [checkpoint_duration].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/checkpoint-transact-sql
type CheckpointStmt struct {
	Duration ExprNode // optional checkpoint duration
	Loc      Loc
}

func (n *CheckpointStmt) nodeTag()  {}
func (n *CheckpointStmt) stmtNode() {}

// ReconfigureStmt represents RECONFIGURE [WITH OVERRIDE].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/reconfigure-transact-sql
type ReconfigureStmt struct {
	WithOverride bool
	Loc          Loc
}

func (n *ReconfigureStmt) nodeTag()  {}
func (n *ReconfigureStmt) stmtNode() {}

// ShutdownStmt represents SHUTDOWN [WITH NOWAIT].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/shutdown-transact-sql
type ShutdownStmt struct {
	WithNoWait bool
	Loc        Loc
}

func (n *ShutdownStmt) nodeTag()  {}
func (n *ShutdownStmt) stmtNode() {}

// KillStmt represents KILL { session_id | UOW } [WITH STATUSONLY | COMMIT | ROLLBACK].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/kill-transact-sql
type KillStmt struct {
	SessionID  ExprNode // session ID or UOW
	StatusOnly bool
	WithAction string // "COMMIT" or "ROLLBACK" (for UOW)
	Loc        Loc
}

func (n *KillStmt) nodeTag()  {}
func (n *KillStmt) stmtNode() {}

// KillQueryNotificationStmt represents KILL QUERY NOTIFICATION SUBSCRIPTION { ALL | subscription_id }.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/language-elements/kill-query-notification-subscription-transact-sql
type KillQueryNotificationStmt struct {
	All            bool     // true if ALL is specified
	SubscriptionID ExprNode // subscription_id expression (when not ALL)
	Loc            Loc
}

func (n *KillQueryNotificationStmt) nodeTag()  {}
func (n *KillQueryNotificationStmt) stmtNode() {}

// ReadtextStmt represents READTEXT table.column textpointer offset size [HOLDLOCK].
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/readtext-transact-sql
type ReadtextStmt struct {
	Column    *ColumnRef
	TextPtr   ExprNode
	Offset    ExprNode
	Size      ExprNode
	HoldLock  bool
	Loc       Loc
}

func (n *ReadtextStmt) nodeTag()  {}
func (n *ReadtextStmt) stmtNode() {}

// WritetextStmt represents WRITETEXT table.column textpointer [WITH LOG] data.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/writetext-transact-sql
type WritetextStmt struct {
	Column  *ColumnRef
	TextPtr ExprNode
	WithLog bool
	Data    ExprNode
	Loc     Loc
}

func (n *WritetextStmt) nodeTag()  {}
func (n *WritetextStmt) stmtNode() {}

// UpdatetextStmt represents UPDATETEXT table.column textpointer deleteoffset deletelength [WITH LOG] inserteddata.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/updatetext-transact-sql
type UpdatetextStmt struct {
	DestColumn   *ColumnRef
	DestTextPtr  ExprNode
	InsertOffset ExprNode
	DeleteLength ExprNode
	WithLog      bool
	InsertedData ExprNode
	Loc          Loc
}

func (n *UpdatetextStmt) nodeTag()  {}
func (n *UpdatetextStmt) stmtNode() {}

// ---------- Batch 48: DROP extended ----------
// (Uses existing DropStmt with extended DropObjectType constants below)

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

// ---------- Batch 35: PIVOT/UNPIVOT ----------

// PivotExpr represents a PIVOT operation.
type PivotExpr struct {
	Source   TableExpr // source table
	AggFunc  ExprNode  // aggregate function call (e.g., SUM(Amount))
	ForCol   string    // FOR column name
	InValues *List     // IN list of column values (as String nodes)
	Alias    string
	Loc      Loc
}

func (n *PivotExpr) nodeTag()   {}
func (n *PivotExpr) tableExpr() {}

// UnpivotExpr represents an UNPIVOT operation.
type UnpivotExpr struct {
	Source   TableExpr // source table
	ValueCol string    // value column name
	ForCol   string    // FOR column name
	InCols   *List     // IN list of column names (as String nodes)
	Alias    string
	Loc      Loc
}

func (n *UnpivotExpr) nodeTag()   {}
func (n *UnpivotExpr) tableExpr() {}

// ---------- Batch 36: TABLESAMPLE ----------

// TableSampleClause represents TABLESAMPLE on a table reference.
type TableSampleClause struct {
	Size       ExprNode // sample size expression
	Unit       string   // "PERCENT" or "ROWS"
	Repeatable ExprNode // REPEATABLE seed (optional)
	Loc        Loc
}

func (n *TableSampleClause) nodeTag() {}

// ---------- Batch 83: Table Hints ----------

// TableHint represents a single table hint in WITH (hint, hint, ...).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-table
//
//	<table_hint> ::=
//	{ NOEXPAND
//	  | INDEX ( <index_value> [ , ...n ] ) | INDEX = ( <index_value> )
//	  | FORCESEEK [ ( <index_value> ( <index_column_name> [ , ... ] ) ) ]
//	  | FORCESCAN
//	  | HOLDLOCK
//	  | NOLOCK
//	  | NOWAIT
//	  | PAGLOCK
//	  | READCOMMITTED
//	  | READCOMMITTEDLOCK
//	  | READPAST
//	  | READUNCOMMITTED
//	  | REPEATABLEREAD
//	  | ROWLOCK
//	  | SERIALIZABLE
//	  | SNAPSHOT
//	  | SPATIAL_WINDOW_MAX_CELLS = <integer_value>
//	  | TABLOCK
//	  | TABLOCKX
//	  | UPDLOCK
//	  | XLOCK
//	}
type TableHint struct {
	Name             string   // hint name: NOLOCK, INDEX, FORCESEEK, etc.
	IndexValues      *List    // INDEX(val, ...) or FORCESEEK index value
	ForceSeekColumns *List    // FORCESEEK(idx(col, col, ...)) column names
	IntValue         ExprNode // SPATIAL_WINDOW_MAX_CELLS = N
	Loc              Loc
}

func (n *TableHint) nodeTag() {}

// ---------- Batch 38: GROUPING SETS/CUBE/ROLLUP ----------

// GroupingSetsExpr represents GROUPING SETS (...) in GROUP BY.
type GroupingSetsExpr struct {
	Sets *List // list of *List (each set is a list of expressions)
	Loc  Loc
}

func (n *GroupingSetsExpr) nodeTag()  {}
func (n *GroupingSetsExpr) exprNode() {}

// RollupExpr represents ROLLUP (...) in GROUP BY.
type RollupExpr struct {
	Args *List
	Loc  Loc
}

func (n *RollupExpr) nodeTag()  {}
func (n *RollupExpr) exprNode() {}

// CubeExpr represents CUBE (...) in GROUP BY.
type CubeExpr struct {
	Args *List
	Loc  Loc
}

func (n *CubeExpr) nodeTag()  {}
func (n *CubeExpr) exprNode() {}

// ---------- Batch 60: Server-level objects ----------

// AlterServerConfigurationStmt represents ALTER SERVER CONFIGURATION SET <option>.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-server-configuration-transact-sql
type AlterServerConfigurationStmt struct {
	OptionType string // e.g. "PROCESS AFFINITY", "DIAGNOSTICS LOG", "BUFFER POOL EXTENSION", etc.
	Options    *List  // list of *ServerConfigOption nodes
	Loc        Loc
}

func (n *AlterServerConfigurationStmt) nodeTag()  {}
func (n *AlterServerConfigurationStmt) stmtNode() {}

// ServerConfigOption represents a typed option in server-level statements.
// Used by ALTER SERVER CONFIGURATION, CREATE/ALTER SERVER ROLE.
// Replaces nodes.String key=value concatenations with structured Name + Value pairs.
//
// Examples:
//   - ON/OFF flags: Name="ON", Value=""
//   - Key=value: Name="CPU", Value="AUTO"
//   - Key=range: Name="CPU", Value="0 TO 3, 8 TO 11"
//   - File paths: Name="FILENAME", Value="'/path/to/file'"
//   - Size specs: Name="SIZE", Value="10 GB"
//   - Actions: Name="AUTHORIZATION", Value="securityadmin"
//   - Sub-options: Name="RESOURCE_POOL", Value="'mypool'"
type ServerConfigOption struct {
	Name  string // option name (e.g., CPU, FILENAME, SIZE, ON, OFF, AUTHORIZATION)
	Value string // option value (may be empty for flag-like options)
	Loc   Loc
}

func (n *ServerConfigOption) nodeTag() {}

// ---------- Batch 61: FULLTEXT STOPLIST / SEARCH PROPERTY LIST ----------

// CreateFulltextStoplistStmt represents CREATE FULLTEXT STOPLIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-fulltext-stoplist-transact-sql
type CreateFulltextStoplistStmt struct {
	Name          string // stoplist name
	SourceDB      string // optional source database name
	SourceList    string // optional source stoplist name
	SystemStoplist bool  // FROM SYSTEM STOPLIST
	Authorization string // AUTHORIZATION owner_name
	Loc           Loc
}

func (n *CreateFulltextStoplistStmt) nodeTag()  {}
func (n *CreateFulltextStoplistStmt) stmtNode() {}

// AlterFulltextStoplistStmt represents ALTER FULLTEXT STOPLIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-fulltext-stoplist-transact-sql
type AlterFulltextStoplistStmt struct {
	Name     string // stoplist name
	Action   string // ADD or DROP
	Stopword string // the stopword (for ADD/DROP single)
	IsNStr   bool   // N prefix on stopword string
	Language string // LANGUAGE term
	DropAll  bool   // DROP ALL (all stopwords) or DROP ALL LANGUAGE
	Loc      Loc
}

func (n *AlterFulltextStoplistStmt) nodeTag()  {}
func (n *AlterFulltextStoplistStmt) stmtNode() {}

// DropFulltextStoplistStmt represents DROP FULLTEXT STOPLIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-fulltext-stoplist-transact-sql
type DropFulltextStoplistStmt struct {
	Name string
	Loc  Loc
}

func (n *DropFulltextStoplistStmt) nodeTag()  {}
func (n *DropFulltextStoplistStmt) stmtNode() {}

// CreateSearchPropertyListStmt represents CREATE SEARCH PROPERTY LIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-search-property-list-transact-sql
type CreateSearchPropertyListStmt struct {
	Name          string // property list name
	SourceDB      string // optional source database name
	SourceList    string // optional source property list name
	Authorization string // AUTHORIZATION owner_name
	Loc           Loc
}

func (n *CreateSearchPropertyListStmt) nodeTag()  {}
func (n *CreateSearchPropertyListStmt) stmtNode() {}

// AlterSearchPropertyListStmt represents ALTER SEARCH PROPERTY LIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-search-property-list-transact-sql
type AlterSearchPropertyListStmt struct {
	Name            string // property list name
	Action          string // ADD or DROP
	PropertyName    string // the search property name
	PropertySetGUID string // PROPERTY_SET_GUID (for ADD)
	PropertyIntID   string // PROPERTY_INT_ID (for ADD)
	PropertyDesc    string // PROPERTY_DESCRIPTION (for ADD, optional)
	Loc             Loc
}

func (n *AlterSearchPropertyListStmt) nodeTag()  {}
func (n *AlterSearchPropertyListStmt) stmtNode() {}

// DropSearchPropertyListStmt represents DROP SEARCH PROPERTY LIST.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-search-property-list-transact-sql
type DropSearchPropertyListStmt struct {
	Name string
	Loc  Loc
}

func (n *DropSearchPropertyListStmt) nodeTag()  {}
func (n *DropSearchPropertyListStmt) stmtNode() {}

// ---------- Batch 62: Security Policy / Classification / Signature ----------

// SecurityPolicyStmt represents CREATE/ALTER/DROP SECURITY POLICY (row-level security).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-security-policy-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-security-policy-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-security-policy-transact-sql
type SecurityPolicyStmt struct {
	Action          string   // CREATE, ALTER, DROP
	Name            *TableRef // [schema.]policy_name
	IfExists        bool     // DROP IF EXISTS
	Predicates      *List    // list of *SecurityPredicate (for CREATE/ALTER)
	StateOn         *bool    // WITH (STATE = ON|OFF); nil if unspecified
	SchemaBinding   *bool    // WITH (SCHEMABINDING = ON|OFF); nil if unspecified
	NotForReplication bool
	Loc             Loc
}

func (n *SecurityPolicyStmt) nodeTag()  {}
func (n *SecurityPolicyStmt) stmtNode() {}

// SecurityPredicate represents a FILTER/BLOCK PREDICATE in a security policy.
type SecurityPredicate struct {
	Action        string   // ADD, ALTER, DROP
	PredicateType string   // FILTER or BLOCK
	Function      *TableRef // schema.function_name
	Args          *List    // function arguments (column names / expressions)
	Table         *TableRef // ON schema.table
	BlockDMLOp    string   // AFTER INSERT, AFTER UPDATE, BEFORE UPDATE, BEFORE DELETE, or ""
	Loc           Loc
}

func (n *SecurityPredicate) nodeTag() {}

// SensitivityClassificationStmt represents ADD/DROP SENSITIVITY CLASSIFICATION.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/add-sensitivity-classification-transact-sql
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/drop-sensitivity-classification-transact-sql
type SensitivityClassificationStmt struct {
	Action  string // ADD or DROP
	Columns *List  // list of *TableRef (schema.table.column references)
	Options *List  // list of *SensitivityOption (for ADD: LABEL, LABEL_ID, INFORMATION_TYPE, etc.)
	Loc     Loc
}

func (n *SensitivityClassificationStmt) nodeTag()  {}
func (n *SensitivityClassificationStmt) stmtNode() {}

// SignatureStmt represents ADD/DROP [COUNTER] SIGNATURE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/add-signature-transact-sql
type SignatureStmt struct {
	Action      string // ADD or DROP
	IsCounter   bool   // COUNTER SIGNATURE
	ModuleClass string // OBJECT, ASSEMBLY, etc. (default OBJECT)
	ModuleName  *TableRef // module name
	CryptoList  *List  // list of *CryptoItem with certificate/key references
	Loc         Loc
}

func (n *SignatureStmt) nodeTag()  {}
func (n *SignatureStmt) stmtNode() {}

// ---------- Batch 63: Specialized Indexes / Aggregate ----------

// CreateXmlIndexStmt represents CREATE [PRIMARY] XML INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-xml-index-transact-sql
type CreateXmlIndexStmt struct {
	Primary      bool      // PRIMARY XML INDEX
	Name         string    // index name
	Table        *TableRef // ON table
	XmlColumn    string    // (xml_column)
	UsingIndex   string    // USING XML INDEX parent_index_name (secondary only)
	SecondaryFor string    // FOR VALUE|PATH|PROPERTY (secondary only)
	Options      *List     // WITH (options)
	Loc          Loc
}

func (n *CreateXmlIndexStmt) nodeTag()  {}
func (n *CreateXmlIndexStmt) stmtNode() {}

// CreateSelectiveXmlIndexStmt represents CREATE SELECTIVE XML INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-selective-xml-index-transact-sql
type CreateSelectiveXmlIndexStmt struct {
	Name       string    // index name
	Table      *TableRef // ON table
	XmlColumn  string    // (xml_column)
	Namespaces *List     // WITH XMLNAMESPACES(...) promoted paths
	Paths      *List     // FOR (path_list) as *String items
	Options    *List     // WITH (options)
	Loc        Loc
}

func (n *CreateSelectiveXmlIndexStmt) nodeTag()  {}
func (n *CreateSelectiveXmlIndexStmt) stmtNode() {}

// CreateSpatialIndexStmt represents CREATE SPATIAL INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-spatial-index-transact-sql
type CreateSpatialIndexStmt struct {
	Name          string    // index name
	Table         *TableRef // ON table
	SpatialColumn string    // (spatial_column)
	Using         string    // USING tessellation type (GEOMETRY_GRID, etc.)
	Options       *List     // WITH (options)
	OnFileGroup   string    // ON filegroup
	Loc           Loc
}

func (n *CreateSpatialIndexStmt) nodeTag()  {}
func (n *CreateSpatialIndexStmt) stmtNode() {}

// CreateAggregateStmt represents CREATE AGGREGATE.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-aggregate-transact-sql
type CreateAggregateStmt struct {
	Name         *TableRef  // [schema.]aggregate_name
	Params       *List      // parameters as *ParamDef
	ReturnType   *DataType  // RETURNS type
	ExternalName string     // EXTERNAL NAME assembly[.class]
	Loc          Loc
}

func (n *CreateAggregateStmt) nodeTag()  {}
func (n *CreateAggregateStmt) stmtNode() {}

// DropAggregateStmt represents DROP AGGREGATE.
type DropAggregateStmt struct {
	Name     *TableRef // [schema.]aggregate_name
	IfExists bool
	Loc      Loc
}

func (n *DropAggregateStmt) nodeTag()  {}
func (n *DropAggregateStmt) stmtNode() {}

// ---------- Batch 109: JSON Index / Vector Index ----------

// CreateJsonIndexStmt represents CREATE JSON INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-json-index-transact-sql
type CreateJsonIndexStmt struct {
	Name        string    // index name
	Table       *TableRef // ON table
	JsonColumn  string    // (json_column_name)
	ForPaths    *List     // FOR ('$.path1', ...) as *String items
	Options     *List     // WITH (options)
	OnFileGroup string    // ON filegroup
	Loc         Loc
}

func (n *CreateJsonIndexStmt) nodeTag()  {}
func (n *CreateJsonIndexStmt) stmtNode() {}

// CreateVectorIndexStmt represents CREATE VECTOR INDEX.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-vector-index-transact-sql
type CreateVectorIndexStmt struct {
	Name        string    // index name
	Table       *TableRef // ON table
	VectorCol   string    // (vector_column)
	Options     *List     // WITH (options)
	OnFileGroup string    // ON filegroup
	Loc         Loc
}

func (n *CreateVectorIndexStmt) nodeTag()  {}
func (n *CreateVectorIndexStmt) stmtNode() {}

// CreateMaterializedViewStmt represents a CREATE MATERIALIZED VIEW AS SELECT statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-materialized-view-as-select-transact-sql
//
//	CREATE MATERIALIZED VIEW [ schema_name. ] materialized_view_name
//	    WITH (
//	      <distribution_option>
//	      [, FOR_APPEND ]
//	    )
//	    AS <select_statement>
//
//	<distribution_option> ::=
//	    {
//	        DISTRIBUTION = HASH ( distribution_column_name [, ...n] )
//	      | DISTRIBUTION = ROUND_ROBIN
//	    }
type CreateMaterializedViewStmt struct {
	Name         *TableRef // view name (may include schema)
	Distribution string    // "HASH" or "ROUND_ROBIN"
	HashColumns  *List     // columns for HASH distribution
	ForAppend    bool      // FOR_APPEND option
	Query        StmtNode  // the SELECT statement
	Loc          Loc
}

func (n *CreateMaterializedViewStmt) nodeTag()  {}
func (n *CreateMaterializedViewStmt) stmtNode() {}

// AlterMaterializedViewStmt represents an ALTER MATERIALIZED VIEW statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-materialized-view-transact-sql
//
//	ALTER MATERIALIZED VIEW [ schema_name. ] view_name
//	{
//	    REBUILD | DISABLE
//	}
type AlterMaterializedViewStmt struct {
	Name   *TableRef // view name (may include schema)
	Action string    // "REBUILD" or "DISABLE"
	Loc    Loc
}

func (n *AlterMaterializedViewStmt) nodeTag()  {}
func (n *AlterMaterializedViewStmt) stmtNode() {}

// CopyIntoStmt represents a COPY INTO statement (Azure Synapse / Fabric).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/copy-into-transact-sql
//
//	COPY INTO [ schema. ] table_name
//	[ (Column_list) ]
//	FROM '<external_location>' [ , ...n ]
//	WITH
//	(
//	  [ FILE_TYPE = { 'CSV' | 'PARQUET' | 'ORC' } ]
//	  [ , FILE_FORMAT = EXTERNAL FILE FORMAT OBJECT ]
//	  [ , CREDENTIAL = (AZURE CREDENTIAL) ]
//	  [ , ERRORFILE = '[http(s)://storageaccount/container]/errorfile_directory[/]' ]
//	  [ , ERRORFILE_CREDENTIAL = (AZURE CREDENTIAL) ]
//	  [ , MAXERRORS = max_errors ]
//	  [ , COMPRESSION = { 'Gzip' | 'DefaultCodec' | 'Snappy' } ]
//	  [ , FIELDQUOTE = 'string_delimiter' ]
//	  [ , FIELDTERMINATOR = 'field_terminator' ]
//	  [ , ROWTERMINATOR = 'row_terminator' ]
//	  [ , FIRSTROW = first_row ]
//	  [ , DATEFORMAT = 'date_format' ]
//	  [ , ENCODING = { 'UTF8' | 'UTF16' } ]
//	  [ , IDENTITY_INSERT = { 'ON' | 'OFF' } ]
//	  [ , AUTO_CREATE_TABLE = { 'ON' | 'OFF' } ]
//	)
type CopyIntoStmt struct {
	Table      *TableRef // target table
	ColumnList *List     // optional column list with default/field_number
	Sources    *List     // list of String (external location URLs)
	Options    *List     // WITH options as key=value or flag strings
	Loc        Loc
}

func (n *CopyIntoStmt) nodeTag()  {}
func (n *CopyIntoStmt) stmtNode() {}

// CopyIntoColumn represents a column entry in COPY INTO column list.
//
//	Column_name [ DEFAULT value ] [ field_number ]
type CopyIntoColumn struct {
	Name         string   // column name
	DefaultValue ExprNode // optional DEFAULT value
	FieldNumber  int      // optional field_number (0 means not specified)
	Loc          Loc
}

func (n *CopyIntoColumn) nodeTag() {}

// RenameStmt represents a RENAME statement (Azure Synapse / PDW).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/rename-transact-sql
//
//	RENAME OBJECT [::] [ [ database_name . [ schema_name ] . ] | [ schema_name . ] ] table_name TO new_table_name
//	RENAME DATABASE [::] database_name TO new_database_name
//	RENAME OBJECT [::] [ [ database_name . [ schema_name ] . ] | [ schema_name . ] ] table_name COLUMN column_name TO new_column_name
type RenameStmt struct {
	ObjectType    string    // "OBJECT" or "DATABASE"
	Name          *TableRef // object being renamed
	NewName       string    // new name
	ColumnName    string    // column being renamed (for RENAME OBJECT ... COLUMN)
	NewColumnName string    // new column name
	Loc           Loc
}

func (n *RenameStmt) nodeTag()  {}
func (n *RenameStmt) stmtNode() {}

// CreateExternalTableAsSelectStmt represents a CREATE EXTERNAL TABLE AS SELECT (CETAS) statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-external-table-as-select-transact-sql
//
//	CREATE EXTERNAL TABLE { [ [ database_name . [ schema_name ] . ] | schema_name . ] table_name }
//	    [ (column_name [ , ...n ] ) ]
//	    WITH (
//	        LOCATION = 'hdfs_folder' | '<prefix>://<path>[:<port>]' ,
//	        DATA_SOURCE = external_data_source_name ,
//	        FILE_FORMAT = external_file_format_name
//	        [ , <reject_options> [ , ...n ] ]
//	    )
//	    AS <select_statement>
type CreateExternalTableAsSelectStmt struct {
	Name    *TableRef // table name
	Columns *List     // optional column list (String nodes)
	Options *List     // WITH options
	Query   Node      // the SELECT statement
	Loc     Loc
}

func (n *CreateExternalTableAsSelectStmt) nodeTag()  {}
func (n *CreateExternalTableAsSelectStmt) stmtNode() {}

// CreateTableCloneStmt represents a CREATE TABLE AS CLONE OF statement (Fabric).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-as-clone-of-transact-sql
//
//	CREATE TABLE
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	AS CLONE OF
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ AT { point_in_time } ]
type CreateTableCloneStmt struct {
	Name       *TableRef // new table name
	SourceName *TableRef // source table name
	AtTime     string    // optional point-in-time (datetime string)
	Loc        Loc
}

func (n *CreateTableCloneStmt) nodeTag()  {}
func (n *CreateTableCloneStmt) stmtNode() {}

// CreateTableAsSelectStmt represents a CREATE TABLE AS SELECT (CTAS) statement
// for Azure Synapse Analytics / Analytics Platform System.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-as-select-azure-sql-data-warehouse
//
//	CREATE TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ ( column_name [ ,...n ] ) ]
//	    WITH (
//	      <distribution_option>
//	      [ , <table_option> [ ,...n ] ]
//	    )
//	    AS <select_statement>
//	    OPTION <query_hint>
//
//	<distribution_option> ::=
//	    {
//	        DISTRIBUTION = HASH ( distribution_column_name [, ...n] )
//	      | DISTRIBUTION = ROUND_ROBIN
//	      | DISTRIBUTION = REPLICATE
//	    }
//
//	<table_option> ::=
//	    {
//	        CLUSTERED COLUMNSTORE INDEX
//	      | CLUSTERED COLUMNSTORE INDEX ORDER ( column [,...n] )
//	      | HEAP
//	      | CLUSTERED INDEX ( { index_column_name [ ASC | DESC ] } [ ,...n ] )
//	    }
//	    | PARTITION ( partition_column_name RANGE [ LEFT | RIGHT ]
//	        FOR VALUES ( [ boundary_value [,...n] ] ) )
type CreateTableAsSelectStmt struct {
	Name    *TableRef // table name
	Columns *List     // optional column list (String nodes)
	Options *List     // WITH options (TableOption nodes for DISTRIBUTION, index type, PARTITION, etc.)
	Query   Node      // the SELECT statement
	Loc     Loc
}

func (n *CreateTableAsSelectStmt) nodeTag()  {}
func (n *CreateTableAsSelectStmt) stmtNode() {}

// CreateRemoteTableAsSelectStmt represents a CREATE REMOTE TABLE AS SELECT (CRTAS) statement
// for Analytics Platform System (PDW).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-remote-table-as-select-parallel-data-warehouse
//
//	CREATE REMOTE TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    AT ('<connection_string>')
//	    [ WITH ( BATCH_SIZE = batch_size ) ]
//	    AS <select_statement>
type CreateRemoteTableAsSelectStmt struct {
	Name             *TableRef // table name
	ConnectionString string    // AT connection string
	Options          *List     // WITH options (e.g., BATCH_SIZE)
	Query            Node      // the SELECT statement
	Loc              Loc
}

func (n *CreateRemoteTableAsSelectStmt) nodeTag()  {}
func (n *CreateRemoteTableAsSelectStmt) stmtNode() {}

// PredictStmt represents a PREDICT statement (SQL Server 2022+).
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/predict-transact-sql
//
//	PREDICT (
//	  MODEL = @model | model_literal,
//	  DATA = object AS <table_alias>
//	  [, RUNTIME = ONNX ]
//	)
//	WITH ( <result_set_definition> )
//
//	<result_set_definition> ::=
//	  { column_name data_type [ COLLATE collation_name ] [ NULL | NOT NULL ] } [,...n]
type PredictStmt struct {
	Model       ExprNode // MODEL = @var | 'literal' | (subquery)
	Data        ExprNode // DATA = table_source
	DataAlias   string   // AS alias for DATA
	Runtime     string   // RUNTIME = ONNX (optional)
	WithColumns *List    // WITH ( column_def [,...n] )
	Loc         Loc
}

func (n *PredictStmt) nodeTag()  {}
func (n *PredictStmt) stmtNode() {}

// ---------- Batch 127: Structured Query Hints ----------

// QueryHint represents a single query hint in an OPTION clause.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/queries/hints-transact-sql-query
//
//	query_hint ::=
//	    { HASH | ORDER } GROUP
//	  | { CONCAT | HASH | MERGE } UNION
//	  | { LOOP | MERGE | HASH } JOIN
//	  | EXPAND VIEWS
//	  | FAST number_rows
//	  | FORCE ORDER
//	  | { FORCE | DISABLE } EXTERNALPUSHDOWN
//	  | { FORCE | DISABLE } SCALEOUTEXECUTION
//	  | IGNORE_NONCLUSTERED_COLUMNSTORE_INDEX
//	  | KEEP PLAN
//	  | KEEPFIXED PLAN
//	  | MAX_GRANT_PERCENT = percent
//	  | MIN_GRANT_PERCENT = percent
//	  | MAXDOP number_of_processors
//	  | MAXRECURSION number
//	  | NO_PERFORMANCE_SPOOL
//	  | OPTIMIZE FOR ( @variable_name { UNKNOWN | = literal } [ , ...n ] )
//	  | OPTIMIZE FOR UNKNOWN
//	  | PARAMETERIZATION { SIMPLE | FORCED }
//	  | QUERYTRACEON trace_flag
//	  | RECOMPILE
//	  | ROBUST PLAN
//	  | USE HINT ( 'hint_name' [ , ...n ] )
//	  | USE PLAN N'xml_plan'
//	  | TABLE HINT ( exposed_object_name [ , <table_hint> [ , ...n ] ] )
type QueryHint struct {
	Kind       string   // Hint kind: "RECOMPILE", "HASH JOIN", "OPTIMIZE FOR", "TABLE HINT", etc.
	Value      ExprNode // Numeric value for MAXDOP, MAXRECURSION, FAST, QUERYTRACEON, MAX_GRANT_PERCENT, MIN_GRANT_PERCENT
	StrValue   string   // String value for PARAMETERIZATION mode, USE PLAN xml string
	Params     *List    // OPTIMIZE FOR params (*OptimizeForParam list), USE HINT string values
	TableName  *TableRef // exposed_object_name for TABLE HINT
	TableHints *List    // list of *TableHint for TABLE HINT
	Loc        Loc
}

func (n *QueryHint) nodeTag() {}

// OptimizeForParam represents a single parameter in OPTIMIZE FOR (@var = val | UNKNOWN).
type OptimizeForParam struct {
	Variable string   // @variable_name
	Unknown  bool     // true if UNKNOWN
	Value    ExprNode // literal value when not UNKNOWN
	Loc      Loc
}

func (n *OptimizeForParam) nodeTag() {}

// ---------- Batch 130: Security Misc Remaining Depth ----------

// CryptoItem represents a single item in a SIGNATURE BY crypto_list.
//
// CERTIFICATE cert_name [ WITH PASSWORD = 'password' | WITH SIGNATURE = signed_blob ]
// ASYMMETRIC KEY key_name [ WITH PASSWORD = 'password' | WITH SIGNATURE = signed_blob ]
type CryptoItem struct {
	Mechanism string    // "CERTIFICATE" or "ASYMMETRIC KEY"
	Name      string    // certificate or key name
	WithType  string    // "PASSWORD" or "SIGNATURE" (empty if none)
	WithValue string    // password or signature blob value
	Loc       Loc
}

func (n *CryptoItem) nodeTag() {}

// SensitivityOption represents a single key=value option in ADD SENSITIVITY CLASSIFICATION WITH clause.
type SensitivityOption struct {
	Key   string // LABEL, LABEL_ID, INFORMATION_TYPE, INFORMATION_TYPE_ID, RANK
	Value string // option value
	Loc   Loc
}

func (n *SensitivityOption) nodeTag() {}

// ---------- Batch 174: Federation Statements ----------

// CreateFederationStmt represents a CREATE FEDERATION statement (retired Azure SQL DB v1).
//
//	CREATE FEDERATION federation_name ( distribution_name data_type RANGE )
type CreateFederationStmt struct {
	Name             string // federation name
	DistributionName string // distribution column name
	DataType         *DataType // scalar data type
	Loc              Loc
}

func (n *CreateFederationStmt) nodeTag()  {}
func (n *CreateFederationStmt) stmtNode() {}

// AlterFederationStmt represents an ALTER FEDERATION statement (retired Azure SQL DB v1).
//
//	ALTER FEDERATION federation_name
//	  { SPLIT AT ( distribution_name = value )
//	  | DROP AT ( LOW | HIGH distribution_name = value ) }
type AlterFederationStmt struct {
	Name             string   // federation name
	Kind             string   // "SPLIT", "DROP LOW", "DROP HIGH"
	DistributionName string   // distribution column name
	Boundary         ExprNode // boundary value expression
	Loc              Loc
}

func (n *AlterFederationStmt) nodeTag()  {}
func (n *AlterFederationStmt) stmtNode() {}

// DropFederationStmt represents a DROP FEDERATION statement (retired Azure SQL DB v1).
//
//	DROP FEDERATION federation_name
type DropFederationStmt struct {
	Name string
	Loc  Loc
}

func (n *DropFederationStmt) nodeTag()  {}
func (n *DropFederationStmt) stmtNode() {}

// UseFederationStmt represents a USE FEDERATION statement (retired Azure SQL DB v1).
//
//	USE FEDERATION { ROOT WITH | federation_name ( distribution_name = value )
//	  WITH FILTERING = { ON | OFF } , } RESET
type UseFederationStmt struct {
	FederationName   string   // federation name (empty for ROOT)
	IsRoot           bool     // true if USE FEDERATION ROOT
	DistributionName string   // distribution column name
	Value            ExprNode // distribution value expression
	Filtering        bool     // FILTERING = ON (true) or OFF (false)
	Loc              Loc
}

func (n *UseFederationStmt) nodeTag()  {}
func (n *UseFederationStmt) stmtNode() {}
