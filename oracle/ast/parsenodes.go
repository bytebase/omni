package ast

// This file contains parse tree node types for Oracle SQL and PL/SQL.
// Reference: Oracle Database 23c SQL Language Reference
// https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/

// ---------------------------------------------------------------------------
// Enums
// ---------------------------------------------------------------------------

// SetOperation represents types of set operations (UNION, INTERSECT, MINUS).
type SetOperation int

const (
	SETOP_NONE SetOperation = iota
	SETOP_UNION
	SETOP_INTERSECT
	SETOP_MINUS
)

// JoinType represents types of joins.
type JoinType int

const (
	JOIN_INNER JoinType = iota
	JOIN_LEFT
	JOIN_FULL
	JOIN_RIGHT
	JOIN_CROSS
	JOIN_NATURAL_INNER
	JOIN_NATURAL_LEFT
	JOIN_NATURAL_RIGHT
	JOIN_NATURAL_FULL
	JOIN_CROSS_APPLY
	JOIN_OUTER_APPLY
)

// SortByDir represents sort ordering direction.
type SortByDir int

const (
	SORTBY_DEFAULT SortByDir = iota
	SORTBY_ASC
	SORTBY_DESC
)

// SortByNulls represents NULLS FIRST/LAST option.
type SortByNulls int

const (
	SORTBY_NULLS_DEFAULT SortByNulls = iota
	SORTBY_NULLS_FIRST
	SORTBY_NULLS_LAST
)

// DropBehavior represents DROP behavior.
type DropBehavior int

const (
	DROP_RESTRICT DropBehavior = iota
	DROP_CASCADE
	DROP_PURGE
)

// ObjectType represents the type of database object.
type ObjectType int

const (
	OBJECT_TABLE ObjectType = iota
	OBJECT_INDEX
	OBJECT_VIEW
	OBJECT_MATERIALIZED_VIEW
	OBJECT_SEQUENCE
	OBJECT_SYNONYM
	OBJECT_PACKAGE
	OBJECT_PACKAGE_BODY
	OBJECT_PROCEDURE
	OBJECT_FUNCTION
	OBJECT_TRIGGER
	OBJECT_TYPE
	OBJECT_TYPE_BODY
	OBJECT_DATABASE_LINK
	OBJECT_USER
	OBJECT_ROLE
	OBJECT_TABLESPACE
	OBJECT_PROFILE
	OBJECT_DIRECTORY
	OBJECT_CONTEXT
	OBJECT_CLUSTER
	OBJECT_DIMENSION
	OBJECT_FLASHBACK_ARCHIVE
	OBJECT_JAVA
	OBJECT_LIBRARY
	OBJECT_MATERIALIZED_VIEW_LOG
	OBJECT_DATABASE
	OBJECT_CONTROLFILE
	OBJECT_DATABASE_DICTIONARY
	OBJECT_DISKGROUP
	OBJECT_PLUGGABLE_DATABASE
	OBJECT_KEY_MANAGEMENT
	OBJECT_AUDIT_POLICY
	OBJECT_ANALYTIC_VIEW
	OBJECT_ATTRIBUTE_DIMENSION
	OBJECT_HIERARCHY
	OBJECT_DOMAIN
	OBJECT_INDEXTYPE
	OBJECT_OPERATOR
	OBJECT_LOCKDOWN_PROFILE
	OBJECT_OUTLINE
	OBJECT_MATERIALIZED_ZONEMAP
	OBJECT_INMEMORY_JOIN_GROUP
	OBJECT_JSON_DUALITY_VIEW
	OBJECT_ROLLBACK_SEGMENT
	OBJECT_RESOURCE_COST
	OBJECT_EDITION
	OBJECT_TABLESPACE_SET
	OBJECT_MLE_ENV
	OBJECT_MLE_MODULE
	OBJECT_PFILE
	OBJECT_SPFILE
	OBJECT_PROPERTY_GRAPH
	OBJECT_VECTOR_INDEX
	OBJECT_RESTORE_POINT
	OBJECT_LOGICAL_PARTITION_TRACKING
	OBJECT_PMEM_FILESTORE
)

// ConstraintType represents constraint types.
type ConstraintType int

const (
	CONSTRAINT_NULL ConstraintType = iota
	CONSTRAINT_NOT_NULL
	CONSTRAINT_DEFAULT
	CONSTRAINT_CHECK
	CONSTRAINT_PRIMARY
	CONSTRAINT_UNIQUE
	CONSTRAINT_FOREIGN
)

// PartitionType represents partition types.
type PartitionType int

const (
	PARTITION_RANGE PartitionType = iota
	PARTITION_LIST
	PARTITION_HASH
	PARTITION_COMPOSITE
)

// TriggerTiming represents trigger timing.
type TriggerTiming int

const (
	TRIGGER_BEFORE TriggerTiming = iota
	TRIGGER_AFTER
	TRIGGER_INSTEAD_OF
)

// TriggerEvent represents trigger events.
type TriggerEvent int

const (
	TRIGGER_INSERT TriggerEvent = iota
	TRIGGER_UPDATE
	TRIGGER_DELETE
)

// BoolExprType represents types of boolean expressions.
type BoolExprType int

const (
	BOOL_AND BoolExprType = iota
	BOOL_OR
	BOOL_NOT
)

// LikeType represents the LIKE variant.
type LikeType int

const (
	LIKE_LIKE LikeType = iota
	LIKE_LIKEC
	LIKE_LIKE2
	LIKE_LIKE4
)

// WindowBoundType represents window frame bound types.
type WindowBoundType int

const (
	WINDOW_UNBOUNDED_PRECEDING WindowBoundType = iota
	WINDOW_CURRENT_ROW
	WINDOW_UNBOUNDED_FOLLOWING
	WINDOW_VALUE_PRECEDING
	WINDOW_VALUE_FOLLOWING
)

// WindowFrameType represents window frame types.
type WindowFrameType int

const (
	WINDOW_ROWS WindowFrameType = iota
	WINDOW_RANGE
	WINDOW_GROUPS
)

// PseudoColumnType represents Oracle pseudo-column types.
type PseudoColumnType int

const (
	PSEUDO_ROWID PseudoColumnType = iota
	PSEUDO_ROWNUM
	PSEUDO_LEVEL
	PSEUDO_SYSDATE
	PSEUDO_SYSTIMESTAMP
	PSEUDO_USER
)

// InsertType represents the type of INSERT statement.
type InsertType int

const (
	INSERT_SINGLE InsertType = iota
	INSERT_ALL
	INSERT_FIRST
)

// AlterTableAction represents ALTER TABLE subcommand types.
type AlterTableAction int

const (
	AT_ADD_COLUMN AlterTableAction = iota
	AT_MODIFY_COLUMN
	AT_DROP_COLUMN
	AT_RENAME_COLUMN
	AT_ADD_CONSTRAINT
	AT_DROP_CONSTRAINT
	AT_MODIFY_CONSTRAINT
	AT_ADD_PARTITION
	AT_DROP_PARTITION
	AT_TRUNCATE_PARTITION
	AT_RENAME
	AT_RENAME_CONSTRAINT
	AT_ENABLE_DISABLE           // ENABLE/DISABLE constraint
	AT_ENABLE_DISABLE_TABLE_LOCK // ENABLE/DISABLE TABLE LOCK
	AT_ENABLE_DISABLE_TRIGGERS  // ENABLE/DISABLE ALL TRIGGERS
	AT_SET_UNUSED               // SET UNUSED COLUMN
	AT_DROP_UNUSED_COLUMNS      // DROP UNUSED COLUMNS / COLUMNS CONTINUE
	AT_MOVE                     // MOVE [ONLINE]
	AT_SPLIT_PARTITION          // SPLIT PARTITION
	AT_SPLIT_SUBPARTITION       // SPLIT SUBPARTITION
	AT_MERGE_PARTITIONS         // MERGE PARTITIONS
	AT_MERGE_SUBPARTITIONS      // MERGE SUBPARTITIONS
	AT_EXCHANGE_PARTITION       // EXCHANGE PARTITION/SUBPARTITION
	AT_COALESCE_PARTITION       // COALESCE PARTITION
	AT_COALESCE_SUBPARTITION    // COALESCE SUBPARTITION
	AT_MODIFY_PARTITION         // MODIFY PARTITION
	AT_MODIFY_SUBPARTITION      // MODIFY SUBPARTITION
	AT_MODIFY_DEFAULT_ATTRS     // MODIFY DEFAULT ATTRIBUTES
	AT_DROP_SUBPARTITION        // DROP SUBPARTITION
	AT_RENAME_PARTITION         // RENAME PARTITION/SUBPARTITION
	AT_SET_INTERVAL             // SET INTERVAL
	AT_SET_PARTITIONING         // SET PARTITIONING
	AT_SET_SUBPARTITION_TEMPLATE // SET SUBPARTITION TEMPLATE
	AT_SHRINK_SPACE             // SHRINK SPACE
	AT_ALTER_PROPERTY           // generic alter property (LOGGING, COMPRESS, etc.)
	AT_ADD_PERIOD               // ADD PERIOD FOR
	AT_DROP_PERIOD              // DROP PERIOD FOR
	AT_MODIFY_LOB               // MODIFY LOB
	AT_MODIFY_VARRAY            // MODIFY VARRAY
	AT_MODIFY_NESTED_TABLE      // MODIFY NESTED TABLE
	AT_MODIFY_OPAQUE_TYPE       // MODIFY OPAQUE TYPE
	AT_MODIFY_TO_PARTITIONED    // MODIFY ... partitioning
	AT_IMMUTABLE_TABLE          // NO DROP / NO DELETE
	AT_BLOCKCHAIN_TABLE         // blockchain clauses
	AT_DUPLICATED_REFRESH       // REFRESH / NO REFRESH
	AT_MODIFY_PARTITIONSET      // MODIFY PARTITIONSET
)

// PLSQLLoopType represents PL/SQL loop types.
type PLSQLLoopType int

const (
	LOOP_BASIC PLSQLLoopType = iota
	LOOP_WHILE
	LOOP_FOR
	LOOP_CURSOR_FOR
)

// ---------------------------------------------------------------------------
// Wrapper / top-level
// ---------------------------------------------------------------------------

// RawStmt wraps a raw (unparsed) statement with position info.
type RawStmt struct {
	Stmt StmtNode // raw statement
	Loc  Loc      // source location range (Start = start offset, End = start + length)
}

func (n *RawStmt) nodeTag() {}

// ---------------------------------------------------------------------------
// Name / reference nodes
// ---------------------------------------------------------------------------

// ObjectName represents a possibly schema-qualified object name.
// Examples: my_table, hr.employees, hr.employees@dblink
type ObjectName struct {
	Schema string // optional schema
	Name   string // object name
	DBLink string // optional @dblink
	Loc    Loc
}

func (n *ObjectName) nodeTag() {}

// ColumnRef represents a column reference (possibly qualified).
type ColumnRef struct {
	Table     string // optional table/alias
	Column    string // column name, or "*"
	Schema    string // optional schema prefix
	OuterJoin bool   // true if followed by (+) — Oracle legacy outer join syntax
	Loc       Loc
}

func (n *ColumnRef) nodeTag()  {}
func (n *ColumnRef) exprNode() {}

// Alias represents an alias (AS clause).
type Alias struct {
	Name string // alias name
	Cols *List  // optional column aliases (list of *String)
	Loc  Loc
}

func (n *Alias) nodeTag() {}

// BindVariable represents a bind variable (:name or :1).
type BindVariable struct {
	Name string // variable name (without colon)
	Loc  Loc
}

func (n *BindVariable) nodeTag()  {}
func (n *BindVariable) exprNode() {}

// PseudoColumn represents an Oracle pseudo-column.
type PseudoColumn struct {
	Type PseudoColumnType
	Loc  Loc
}

func (n *PseudoColumn) nodeTag()  {}
func (n *PseudoColumn) exprNode() {}

// Hint represents an Oracle optimizer hint (/*+ ... */).
type Hint struct {
	Text string // raw hint text
	Loc  Loc
}

func (n *Hint) nodeTag() {}

// ---------------------------------------------------------------------------
// Type nodes
// ---------------------------------------------------------------------------

// TypeName represents a data type reference.
type TypeName struct {
	Names         *List // type name parts (list of *String)
	TypeMods      *List // type modifiers (precision, scale)
	IsPercType    bool  // %TYPE
	IsPercRowtype bool  // %ROWTYPE
	ArrayBounds   *List // array bounds if any
	Loc           Loc   // start location
}

func (n *TypeName) nodeTag() {}

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

// BinaryExpr represents a binary expression (e.g., a + b, a = b).
type BinaryExpr struct {
	Op    string   // operator: +, -, *, /, ||, =, !=, <, >, <=, >=, **, etc.
	Left  ExprNode // left operand
	Right ExprNode // right operand
	Loc   Loc
}

func (n *BinaryExpr) nodeTag()  {}
func (n *BinaryExpr) exprNode() {}

// UnaryExpr represents a unary expression (e.g., -a, NOT a).
type UnaryExpr struct {
	Op      string   // operator: -, +, NOT, PRIOR, CONNECT_BY_ROOT
	Operand ExprNode // operand
	Loc     Loc
}

func (n *UnaryExpr) nodeTag()  {}
func (n *UnaryExpr) exprNode() {}

// BoolExpr represents a boolean expression (AND, OR, NOT).
type BoolExpr struct {
	Boolop BoolExprType
	Args   *List // list of arguments
	Loc    Loc
}

func (n *BoolExpr) nodeTag()  {}
func (n *BoolExpr) exprNode() {}

// FuncCallExpr represents a function call.
type FuncCallExpr struct {
	FuncName   *ObjectName // function name
	Args       *List       // list of argument expressions
	Distinct   bool        // DISTINCT specified
	Star       bool        // * specified (e.g., COUNT(*))
	OrderBy    *List       // ORDER BY within aggregate
	KeepClause *KeepClause // KEEP (DENSE_RANK ...)
	Over       *WindowSpec // OVER clause (analytic function)
	Loc        Loc         // start location
}

func (n *FuncCallExpr) nodeTag()  {}
func (n *FuncCallExpr) exprNode() {}

// CaseExpr represents a CASE expression.
type CaseExpr struct {
	Arg     ExprNode // test expression for simple CASE (nil for searched CASE)
	Whens   *List    // list of *CaseWhen
	Default ExprNode // ELSE expression (nil if absent)
	Loc     Loc
}

func (n *CaseExpr) nodeTag()  {}
func (n *CaseExpr) exprNode() {}

// CaseWhen represents a WHEN clause in a CASE expression.
type CaseWhen struct {
	Condition ExprNode // WHEN condition
	Result    ExprNode // THEN result
	Loc       Loc      // start location
}

func (n *CaseWhen) nodeTag()  {}
func (n *CaseWhen) exprNode() {}

// DecodeExpr represents Oracle's DECODE function (treated as expression).
type DecodeExpr struct {
	Arg     ExprNode // expression to decode
	Pairs   *List    // list of *DecodePair (search, result pairs)
	Default ExprNode // default expression (nil if absent)
	Loc     Loc
}

func (n *DecodeExpr) nodeTag()  {}
func (n *DecodeExpr) exprNode() {}

// DecodePair represents a search-result pair in DECODE.
type DecodePair struct {
	Search ExprNode // search value
	Result ExprNode // result value
	Loc    Loc
}

func (n *DecodePair) nodeTag()  {}
func (n *DecodePair) exprNode() {}

// BetweenExpr represents a BETWEEN expression.
type BetweenExpr struct {
	Expr ExprNode // expression being tested
	Low  ExprNode // lower bound
	High ExprNode // upper bound
	Not  bool     // NOT BETWEEN
	Loc  Loc
}

func (n *BetweenExpr) nodeTag()  {}
func (n *BetweenExpr) exprNode() {}

// InExpr represents an IN expression.
type InExpr struct {
	Expr ExprNode // expression being tested
	List *List    // list of values or subquery
	Not  bool     // NOT IN
	Loc  Loc
}

func (n *InExpr) nodeTag()  {}
func (n *InExpr) exprNode() {}

// LikeExpr represents a LIKE / LIKEC / LIKE2 / LIKE4 expression.
type LikeExpr struct {
	Expr    ExprNode // expression being tested
	Pattern ExprNode // pattern expression
	Escape  ExprNode // ESCAPE expression (nil if absent)
	Not     bool     // NOT LIKE
	Type    LikeType // LIKE variant
	Loc     Loc
}

func (n *LikeExpr) nodeTag()  {}
func (n *LikeExpr) exprNode() {}

// IsExpr represents IS [NOT] {NULL|NAN|INFINITE|EMPTY|OF|A SET} expression.
type IsExpr struct {
	Expr     ExprNode // expression being tested
	Test     string   // "NULL", "NAN", "INFINITE", "EMPTY", "A SET", "OF", "JSON"
	Not      bool     // IS NOT
	TypeList *List    // for IS OF (type list)
	Loc      Loc
}

func (n *IsExpr) nodeTag()  {}
func (n *IsExpr) exprNode() {}

// ExistsExpr represents an EXISTS subquery expression.
type ExistsExpr struct {
	Subquery StmtNode // the subquery
	Loc      Loc
}

func (n *ExistsExpr) nodeTag()  {}
func (n *ExistsExpr) exprNode() {}

// CastExpr represents a CAST expression.
type CastExpr struct {
	Arg      ExprNode  // expression to cast
	TypeName *TypeName // target type
	Loc      Loc       // start location
}

func (n *CastExpr) nodeTag()  {}
func (n *CastExpr) exprNode() {}

// MultisetExpr represents a MULTISET expression.
type MultisetExpr struct {
	Op    string   // UNION, INTERSECT, EXCEPT
	Left  ExprNode // left operand
	Right ExprNode // right operand
	All   bool     // ALL or DISTINCT
	Loc   Loc
}

func (n *MultisetExpr) nodeTag()  {}
func (n *MultisetExpr) exprNode() {}

// CursorExpr represents a CURSOR(subquery) expression.
type CursorExpr struct {
	Subquery StmtNode // the subquery
	Loc      Loc
}

func (n *CursorExpr) nodeTag()  {}
func (n *CursorExpr) exprNode() {}

// TreatExpr represents a TREAT(expr AS type) expression.
type TreatExpr struct {
	Expr     ExprNode  // expression to treat
	TypeName *TypeName // target type
	Loc      Loc
}

func (n *TreatExpr) nodeTag()  {}
func (n *TreatExpr) exprNode() {}

// SubqueryExpr represents a scalar subquery expression.
type SubqueryExpr struct {
	Subquery StmtNode // the subquery (a *SelectStmt)
	Loc      Loc
}

func (n *SubqueryExpr) nodeTag()   {}
func (n *SubqueryExpr) exprNode()  {}
func (n *SubqueryExpr) tableExpr() {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Expr ExprNode // inner expression
	Loc  Loc
}

func (n *ParenExpr) nodeTag()  {}
func (n *ParenExpr) exprNode() {}

// NullLiteral represents a NULL literal.
type NullLiteral struct {
	Loc Loc
}

func (n *NullLiteral) nodeTag()  {}
func (n *NullLiteral) exprNode() {}

// StringLiteral represents a string literal.
type StringLiteral struct {
	Val     string // string value
	IsNChar bool   // N'...' national character literal
	Loc     Loc
}

func (n *StringLiteral) nodeTag()  {}
func (n *StringLiteral) exprNode() {}

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Val     string // raw numeric text
	Ival    int64  // parsed integer value (if integer)
	IsFloat bool   // true if floating point
	Loc     Loc
}

func (n *NumberLiteral) nodeTag()  {}
func (n *NumberLiteral) exprNode() {}

// IntervalExpr represents an INTERVAL expression.
type IntervalExpr struct {
	Value ExprNode // expression
	From  string   // starting field (YEAR, MONTH, DAY, HOUR, MINUTE, SECOND)
	To    string   // ending field (optional)
	Loc   Loc
}

func (n *IntervalExpr) nodeTag()  {}
func (n *IntervalExpr) exprNode() {}

// ---------------------------------------------------------------------------
// Analytic / window function nodes
// ---------------------------------------------------------------------------

// WindowSpec represents an analytic function's OVER clause.
type WindowSpec struct {
	PartitionBy *List        // PARTITION BY expressions
	OrderBy     *List        // ORDER BY clause (list of *SortBy)
	Frame       *WindowFrame // window frame specification
	WindowName  string       // reference to named window (optional)
	Loc         Loc          // start location
}

func (n *WindowSpec) nodeTag() {}

// WindowFrame represents the frame specification of a window.
type WindowFrame struct {
	Type  WindowFrameType // ROWS, RANGE, GROUPS
	Start *WindowBound    // start bound
	End   *WindowBound    // end bound (nil for single bound)
	Loc   Loc             // start location
}

func (n *WindowFrame) nodeTag() {}

// WindowBound represents a window frame boundary.
type WindowBound struct {
	Type  WindowBoundType // bound type
	Value ExprNode        // offset value (for VALUE_PRECEDING / VALUE_FOLLOWING)
	Loc   Loc             // start location
}

func (n *WindowBound) nodeTag() {}

// KeepClause represents the KEEP (DENSE_RANK FIRST/LAST ORDER BY ...) clause.
type KeepClause struct {
	IsFirst bool  // FIRST or LAST
	OrderBy *List // ORDER BY clause
	Loc     Loc
}

func (n *KeepClause) nodeTag() {}

// SortBy represents an ORDER BY element.
type SortBy struct {
	Expr      ExprNode    // sort key expression
	Dir       SortByDir   // ASC/DESC
	NullOrder SortByNulls // NULLS FIRST/LAST
	Loc       Loc         // start location
}

func (n *SortBy) nodeTag() {}

// ---------------------------------------------------------------------------
// Hierarchical query nodes
// ---------------------------------------------------------------------------

// HierarchicalClause represents CONNECT BY / START WITH.
type HierarchicalClause struct {
	ConnectBy ExprNode // CONNECT BY condition
	StartWith ExprNode // START WITH condition (nil if absent)
	IsNocycle bool     // NOCYCLE specified
	Loc       Loc      // start location
}

func (n *HierarchicalClause) nodeTag() {}

// ---------------------------------------------------------------------------
// Join nodes
// ---------------------------------------------------------------------------

// JoinClause represents a JOIN.
type JoinClause struct {
	Type       JoinType  // type of join
	Left       TableExpr // left table reference
	Right      TableExpr // right table reference
	On         ExprNode  // ON condition (nil for CROSS JOIN, USING)
	Using      *List     // USING columns (list of *String)
	OracleJoin bool      // Oracle (+) outer join syntax detected
	Loc        Loc       // start location
}

func (n *JoinClause) nodeTag()   {}
func (n *JoinClause) tableExpr() {}

// TableRef represents a table reference in a FROM clause.
type TableRef struct {
	Name         *ObjectName      // table name
	Alias        *Alias           // optional alias
	Sample       *SampleClause    // SAMPLE clause
	Flashback    *FlashbackClause // AS OF / VERSIONS BETWEEN
	PartitionExt *PartitionExtClause // PARTITION/SUBPARTITION clause
	Dblink       string           // @dblink name
	Loc          Loc              // start location
}

func (n *TableRef) nodeTag()   {}
func (n *TableRef) tableExpr() {}

// SubqueryRef represents a subquery in a FROM clause.
type SubqueryRef struct {
	Subquery StmtNode // the subquery
	Alias    *Alias   // alias
	Loc      Loc
}

func (n *SubqueryRef) nodeTag()   {}
func (n *SubqueryRef) tableExpr() {}

// LateralRef represents a LATERAL inline view in FROM.
//
//	LATERAL ( subquery ) [ alias ]
type LateralRef struct {
	Subquery StmtNode // the subquery
	Alias    *Alias   // alias
	Loc      Loc
}

func (n *LateralRef) nodeTag()   {}
func (n *LateralRef) tableExpr() {}

// XmlTableRef represents an XMLTABLE expression in FROM.
//
//	XMLTABLE ( xpath_expr PASSING xml_expr COLUMNS column_def [, ...] ) [ alias ]
type XmlTableRef struct {
	XPath   ExprNode // XPath string expression
	Passing ExprNode // PASSING xml expression
	Columns *List    // list of *XmlTableColumn
	Alias   *Alias   // alias
	Loc     Loc
}

func (n *XmlTableRef) nodeTag()   {}
func (n *XmlTableRef) tableExpr() {}

// XmlTableColumn represents a column definition in XMLTABLE.
type XmlTableColumn struct {
	Name     string    // column name
	TypeName *TypeName // data type
	Path     ExprNode  // PATH string (nil if default)
	Default  ExprNode  // DEFAULT value
	ForOrdinality bool // FOR ORDINALITY
	Loc      Loc
}

func (n *XmlTableColumn) nodeTag() {}

// JsonTableRef represents a JSON_TABLE expression in FROM.
//
//	JSON_TABLE ( expr, path_expr COLUMNS ( column_def [, ...] ) ) [ alias ]
type JsonTableRef struct {
	Expr    ExprNode // JSON expression
	Path    ExprNode // JSON path string
	Columns *List    // list of *JsonTableColumn
	Alias   *Alias   // alias
	Loc     Loc
}

func (n *JsonTableRef) nodeTag()   {}
func (n *JsonTableRef) tableExpr() {}

// JsonTableColumn represents a column definition in JSON_TABLE.
type JsonTableColumn struct {
	Name     string    // column name
	TypeName *TypeName // data type
	Path     ExprNode  // PATH string
	ForOrdinality bool // FOR ORDINALITY
	Exists   bool      // EXISTS
	Nested   *JsonTableRef // NESTED [PATH] for nested columns
	Loc      Loc
}

func (n *JsonTableColumn) nodeTag() {}

// SampleClause represents a SAMPLE clause.
type SampleClause struct {
	Percent ExprNode // sample percentage
	Seed    ExprNode // SEED value (nil if absent)
	Block   bool     // BLOCK keyword
	Loc     Loc
}

func (n *SampleClause) nodeTag() {}

// ---------------------------------------------------------------------------
// Pivot / Unpivot
// ---------------------------------------------------------------------------

// PivotClause represents a PIVOT clause.
type PivotClause struct {
	XML      bool     // PIVOT XML
	AggFuncs *List    // aggregate functions
	ForCol   ExprNode // FOR column (single column)
	ForCols  *List    // FOR (col1, col2, ...) multi-column
	InList   *List    // IN list
	Alias    *Alias   // optional alias
	Loc      Loc
}

func (n *PivotClause) nodeTag()   {}
func (n *PivotClause) tableExpr() {}

// UnpivotClause represents an UNPIVOT clause.
type UnpivotClause struct {
	ValueCol     ExprNode // value column
	PivotCol     ExprNode // pivot column
	InList       *List    // IN list
	IncludeNulls bool     // INCLUDE NULLS
	Alias        *Alias   // optional alias
	Loc          Loc      // start location
}

func (n *UnpivotClause) nodeTag()   {}
func (n *UnpivotClause) tableExpr() {}

// ---------------------------------------------------------------------------
// WITH clause (subquery factoring)
// ---------------------------------------------------------------------------

// WithClause represents a WITH clause (common table expressions).
type WithClause struct {
	CTEs      *List // list of *CTE
	Recursive bool  // WITH RECURSIVE
	Loc       Loc   // start location
}

func (n *WithClause) nodeTag() {}

// CTE represents a common table expression.
type CTE struct {
	Name    string   // CTE name
	Columns *List    // optional column list (list of *String)
	Query   StmtNode // the CTE query
	Search  *CTESearchClause // SEARCH clause
	Cycle   *CTECycleClause  // CYCLE clause
	Loc     Loc
}

func (n *CTE) nodeTag() {}

// ---------------------------------------------------------------------------
// SELECT statement
// ---------------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
type SelectStmt struct {
	WithClause   *WithClause         // WITH clause
	Distinct     bool                // DISTINCT
	UniqueKw     bool                // UNIQUE (Oracle synonym for DISTINCT)
	All          bool                // ALL
	TargetList   *List               // select expressions (list of *ResTarget)
	Into         *ObjectName         // INTO (PL/SQL)
	FromClause   *List               // FROM clause
	WhereClause  ExprNode            // WHERE condition
	Hierarchical *HierarchicalClause // CONNECT BY / START WITH
	GroupClause  *List               // GROUP BY
	HavingClause ExprNode            // HAVING condition
	ModelClause   *ModelClause        // MODEL clause
	WindowDefs    []*WindowDef        // WINDOW clause (named window definitions)
	QualifyClause ExprNode            // QUALIFY condition
	OrderBy       *List               // ORDER BY (list of *SortBy)
	SiblingsOrder bool                // ORDER SIBLINGS BY
	ForUpdate     *ForUpdateClause    // FOR UPDATE
	FetchFirst    *FetchFirstClause   // FETCH FIRST / OFFSET
	Pivot         *PivotClause        // PIVOT clause
	Unpivot       *UnpivotClause      // UNPIVOT clause
	Hints         *List               // optimizer hints (list of *Hint)

	// Set operations
	Op     SetOperation // UNION, INTERSECT, MINUS
	SetAll bool         // ALL specified with set op
	Larg   *SelectStmt  // left child
	Rarg   *SelectStmt  // right child

	Loc Loc
}

func (n *SelectStmt) nodeTag()  {}
func (n *SelectStmt) stmtNode() {}

// ResTarget represents a target in the select list.
type ResTarget struct {
	Name string   // column alias (AS name)
	Expr ExprNode // expression
	Loc  Loc
}

func (n *ResTarget) nodeTag() {}

// ForUpdateClause represents a FOR UPDATE clause.
type ForUpdateClause struct {
	Tables     *List    // OF table list (nil for all)
	NoWait     bool     // NOWAIT
	Wait       ExprNode // WAIT seconds expression
	SkipLocked bool     // SKIP LOCKED
	Loc        Loc      // start location
}

func (n *ForUpdateClause) nodeTag() {}

// FetchFirstClause represents FETCH FIRST / OFFSET / ROWNUM limiting.
type FetchFirstClause struct {
	Offset   ExprNode // OFFSET expression
	Count    ExprNode // FETCH FIRST count expression
	Percent  bool     // PERCENT specified
	WithTies bool     // WITH TIES
	Loc      Loc
}

func (n *FetchFirstClause) nodeTag() {}

// ModelClause represents an Oracle MODEL clause.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	model_clause ::=
//	    MODEL
//	      [ cell_reference_options ]
//	      [ return_rows_clause ]
//	      [ reference_model ]...
//	      main_model
type ModelClause struct {
	CellRefOptions *ModelCellRefOptions // IGNORE NAV / KEEP NAV, UNIQUE DIMENSION / UNIQUE SINGLE REFERENCE
	ReturnRows     string               // "" (default), "UPDATED", "ALL"
	RefModels      []*ModelRefModel     // REFERENCE models
	MainModel      *ModelMainModel      // MAIN model
	Loc            Loc
}

func (n *ModelClause) nodeTag() {}

// ModelCellRefOptions represents cell_reference_options.
//
//	cell_reference_options ::=
//	    { IGNORE NAV | KEEP NAV }
//	    { UNIQUE DIMENSION | UNIQUE SINGLE REFERENCE }
type ModelCellRefOptions struct {
	IgnoreNav       bool // IGNORE NAV (vs KEEP NAV)
	KeepNav         bool // KEEP NAV
	UniqueDimension bool // UNIQUE DIMENSION
	UniqueSingleRef bool // UNIQUE SINGLE REFERENCE
	Loc             Loc
}

func (n *ModelCellRefOptions) nodeTag() {}

// ModelRefModel represents a REFERENCE model.
//
//	reference_model ::=
//	    REFERENCE reference_model_name ON ( subquery )
//	        model_column_clauses
//	        [ cell_reference_options ]
type ModelRefModel struct {
	Name           string               // reference model name
	Subquery       *SelectStmt          // ON ( subquery )
	ColumnClauses  *ModelColumnClauses   // PARTITION BY, DIMENSION BY, MEASURES
	CellRefOptions *ModelCellRefOptions  // optional cell ref options
	Loc            Loc
}

func (n *ModelRefModel) nodeTag() {}

// ModelMainModel represents the main_model.
//
//	main_model ::=
//	    [ MAIN main_model_name ]
//	        model_column_clauses
//	        [ cell_reference_options ]
//	        model_rules_clause
type ModelMainModel struct {
	Name           string               // optional main model name
	ColumnClauses  *ModelColumnClauses   // PARTITION BY, DIMENSION BY, MEASURES
	CellRefOptions *ModelCellRefOptions  // optional cell ref options
	RulesClause    *ModelRulesClause     // RULES clause
	Loc            Loc
}

func (n *ModelMainModel) nodeTag() {}

// ModelColumnClauses represents model_column_clauses.
//
//	model_column_clauses ::=
//	    [ PARTITION BY ( expr [ [ AS ] c_alias ] [, ...] ) ]
//	    DIMENSION BY ( expr [ [ AS ] c_alias ] [, ...] )
//	    MEASURES ( expr [ [ AS ] c_alias ] [, ...] )
type ModelColumnClauses struct {
	PartitionBy *List // list of *ResTarget (may be nil)
	DimensionBy *List // list of *ResTarget
	Measures    *List // list of *ResTarget
	Loc         Loc
}

func (n *ModelColumnClauses) nodeTag() {}

// ModelRulesClause represents model_rules_clause.
//
//	model_rules_clause ::=
//	    [ RULES ]
//	    [ { UPDATE | UPSERT [ ALL ] } ]
//	    [ { AUTOMATIC | SEQUENTIAL } ORDER ]
//	    [ model_iterate_clause ]
//	    ( cell_assignment [, ...] )
type ModelRulesClause struct {
	UpdateMode string // "", "UPDATE", "UPSERT", "UPSERT ALL"
	OrderMode  string // "", "AUTOMATIC", "SEQUENTIAL"
	Iterate    ExprNode // iteration count expression
	Until      ExprNode // UNTIL condition
	Rules      *List    // list of *ModelRule
	Loc        Loc
}

func (n *ModelRulesClause) nodeTag() {}

// ModelRule represents a single cell_assignment rule.
//
//	cell_assignment ::=
//	    measure_column [ dimension_conditions ] = expr
type ModelRule struct {
	CellRef ExprNode // left side (cell reference with dimension subscripts)
	Expr    ExprNode // right side value expression
	Loc     Loc
}

func (n *ModelRule) nodeTag() {}

// ModelForLoop represents FOR loop in model cell assignment.
//
//	single_column_for_loop ::=
//	    FOR dimension_column
//	      { IN ( { literal [, ...] | subquery } )
//	      | [ LIKE pattern ] FROM literal TO literal { INCREMENT | DECREMENT } literal
//	      }
type ModelForLoop struct {
	Column    string      // dimension column name
	InList    *List       // IN ( values ) - list of ExprNode
	Subquery  *SelectStmt // IN ( subquery )
	LikePattern ExprNode  // LIKE pattern (for FROM..TO)
	FromExpr  ExprNode    // FROM literal
	ToExpr    ExprNode    // TO literal
	Increment bool        // true = INCREMENT, false = DECREMENT
	IncrExpr  ExprNode    // increment/decrement value
	Loc       Loc
}

func (n *ModelForLoop) nodeTag()  {}
func (n *ModelForLoop) exprNode() {}

// FlashbackClause represents a flashback query (AS OF / VERSIONS BETWEEN).
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/SELECT.html
//
//	flashback_query_clause ::=
//	    { VERSIONS BETWEEN { SCN | TIMESTAMP } expr AND expr
//	    | AS OF { SCN | TIMESTAMP } expr
//	    }
type FlashbackClause struct {
	Type         string   // "SCN" or "TIMESTAMP"
	Expr         ExprNode // AS OF expression
	IsVersions   bool     // VERSIONS BETWEEN (vs AS OF)
	IsPeriodFor  bool     // VERSIONS PERIOD FOR / AS OF PERIOD FOR
	PeriodColumn string   // valid_time_column for PERIOD FOR
	VersionsLow  ExprNode // VERSIONS BETWEEN low expr
	VersionsHigh ExprNode // VERSIONS BETWEEN high expr (AND expr)
	Loc          Loc
}

func (n *FlashbackClause) nodeTag() {}

// ---------------------------------------------------------------------------
// Window clause
// ---------------------------------------------------------------------------

// WindowDef represents a named window definition in the WINDOW clause.
//
//	WINDOW window_name AS ( window_specification )
type WindowDef struct {
	Name string      // window name
	Spec *WindowSpec // window specification
	Loc  Loc
}

func (n *WindowDef) nodeTag() {}

// ---------------------------------------------------------------------------
// CTE SEARCH / CYCLE clauses
// ---------------------------------------------------------------------------

// CTESearchClause represents a SEARCH clause on a CTE.
//
//	SEARCH { BREADTH FIRST | DEPTH FIRST } BY col [ASC|DESC] [NULLS FIRST|LAST] [, ...] SET ordering_column
type CTESearchClause struct {
	BreadthFirst bool   // true = BREADTH FIRST, false = DEPTH FIRST
	Columns      *List  // list of *SortBy
	SetColumn    string // SET ordering_column
	Loc          Loc
}

func (n *CTESearchClause) nodeTag() {}

// CTECycleClause represents a CYCLE clause on a CTE.
//
//	CYCLE col [, ...] SET cycle_mark_c_alias TO cycle_value DEFAULT no_cycle_value
type CTECycleClause struct {
	Columns        *List    // list of *String (column names)
	SetColumn      string   // SET cycle_mark_c_alias
	CycleValue     ExprNode // TO cycle_value
	NoCycleValue   ExprNode // DEFAULT no_cycle_value
	Loc            Loc
}

func (n *CTECycleClause) nodeTag() {}

// ---------------------------------------------------------------------------
// Partition extension clause
// ---------------------------------------------------------------------------

// PartitionExtClause represents a partition extension clause on a table ref.
//
//	PARTITION (name) | PARTITION FOR (key, ...) | SUBPARTITION (name) | SUBPARTITION FOR (key, ...)
type PartitionExtClause struct {
	IsSubpartition bool     // SUBPARTITION vs PARTITION
	IsFor          bool     // FOR (key values) vs (name)
	Name           string   // partition/subpartition name (when IsFor is false)
	Keys           *List    // key values (when IsFor is true)
	Loc            Loc
}

func (n *PartitionExtClause) nodeTag() {}

// ---------------------------------------------------------------------------
// Table collection expression
// ---------------------------------------------------------------------------

// TableCollectionExpr represents TABLE(collection_expression) [(+)] in FROM.
type TableCollectionExpr struct {
	Expr      ExprNode // collection expression
	OuterJoin bool     // (+) specified
	Alias     *Alias   // optional alias
	Loc       Loc
}

func (n *TableCollectionExpr) nodeTag()   {}
func (n *TableCollectionExpr) tableExpr() {}

// ---------------------------------------------------------------------------
// MATCH_RECOGNIZE clause
// ---------------------------------------------------------------------------

// MatchRecognizeClause represents a MATCH_RECOGNIZE clause on a table reference.
type MatchRecognizeClause struct {
	PartitionBy *List    // PARTITION BY columns
	OrderBy     *List    // ORDER BY columns (list of *SortBy)
	Measures    *List    // MEASURES (list of *ResTarget)
	RowsPerMatch string  // "ONE ROW PER MATCH" or "ALL ROWS PER MATCH" (+ SHOW/OMIT EMPTY MATCHES)
	AfterMatch  string   // AFTER MATCH SKIP action
	Pattern     string   // PATTERN as raw text
	Subsets     *List    // SUBSET items (list of *ResTarget)
	Definitions *List    // DEFINE items (list of *ResTarget)
	Alias       *Alias   // optional alias
	Loc         Loc
}

func (n *MatchRecognizeClause) nodeTag()   {}
func (n *MatchRecognizeClause) tableExpr() {}

// ---------------------------------------------------------------------------
// Containers / Shards clause
// ---------------------------------------------------------------------------

// ContainersExpr represents CONTAINERS(schema.table) or SHARDS(schema.table) in FROM.
type ContainersExpr struct {
	IsShards bool        // true = SHARDS, false = CONTAINERS
	Name     *ObjectName // table/view name
	Alias    *Alias      // optional alias
	Loc      Loc
}

func (n *ContainersExpr) nodeTag()   {}
func (n *ContainersExpr) tableExpr() {}

// ---------------------------------------------------------------------------
// Inline external table
// ---------------------------------------------------------------------------

// InlineExternalTable represents an inline external table in FROM.
type InlineExternalTable struct {
	Columns    *List    // column definitions (list of *ColumnDef)
	Type       string   // external table type
	Directory  string   // DEFAULT DIRECTORY
	AccessParams string // ACCESS PARAMETERS text
	Location   string   // LOCATION
	RejectLimit ExprNode // REJECT LIMIT
	Alias      *Alias   // optional alias
	Loc        Loc
}

func (n *InlineExternalTable) nodeTag()   {}
func (n *InlineExternalTable) tableExpr() {}

// GroupingSetsClause represents a GROUPING SETS clause in GROUP BY.
//
//	GROUPING SETS ( grouping_set [, ...] )
type GroupingSetsClause struct {
	Sets *List // list of ExprNode or *List (composite grouping sets)
	Loc  Loc
}

func (n *GroupingSetsClause) nodeTag()  {}
func (n *GroupingSetsClause) exprNode() {}

// CubeClause represents a CUBE clause in GROUP BY.
//
//	CUBE ( expr [, ...] )
type CubeClause struct {
	Args *List // list of ExprNode
	Loc  Loc
}

func (n *CubeClause) nodeTag()  {}
func (n *CubeClause) exprNode() {}

// RollupClause represents a ROLLUP clause in GROUP BY.
//
//	ROLLUP ( expr [, ...] )
type RollupClause struct {
	Args *List // list of ExprNode
	Loc  Loc
}

func (n *RollupClause) nodeTag()  {}
func (n *RollupClause) exprNode() {}

// ---------------------------------------------------------------------------
// INSERT statement
// ---------------------------------------------------------------------------

// InsertStmt represents an INSERT statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/INSERT.html
type InsertStmt struct {
	InsertType   InsertType          // single, ALL, FIRST
	Table        *ObjectName         // target table (for single insert)
	PartitionExt *PartitionExtClause // PARTITION/SUBPARTITION extension
	Dblink       string              // @dblink name
	Alias        *Alias              // table alias
	Columns      *List               // column list (list of *ColumnRef)
	Values       *List               // VALUES list (list of expressions)
	SetClauses   *List               // INSERT SET col=expr form (list of *SetClause)
	ByName       bool                // BY NAME subquery
	ByPosition   bool                // BY POSITION subquery
	Select       *SelectStmt         // subquery source
	MultiTable   *List               // for INSERT ALL/FIRST: list of *InsertIntoClause
	Subquery     StmtNode            // source subquery for multi-table insert
	Returning    *List               // RETURNING clause
	ErrorLog     *ErrorLogClause     // LOG ERRORS
	Hints        *List               // optimizer hints
	Loc          Loc                 // start location
}

func (n *InsertStmt) nodeTag()  {}
func (n *InsertStmt) stmtNode() {}

// InsertIntoClause represents a single INTO clause in multi-table INSERT.
type InsertIntoClause struct {
	Table   *ObjectName // target table
	Columns *List       // column list
	Values  *List       // VALUES list
	When    ExprNode    // WHEN condition (for conditional insert)
	Loc     Loc         // start location
}

func (n *InsertIntoClause) nodeTag() {}

// ErrorLogClause represents a LOG ERRORS clause.
type ErrorLogClause struct {
	Into   *ObjectName // error log table
	Tag    ExprNode    // simple expression tag
	Reject ExprNode    // REJECT LIMIT
	Loc    Loc         // start location
}

func (n *ErrorLogClause) nodeTag() {}

// ---------------------------------------------------------------------------
// UPDATE statement
// ---------------------------------------------------------------------------

// UpdateStmt represents an UPDATE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/UPDATE.html
type UpdateStmt struct {
	Table        *ObjectName         // table to update
	PartitionExt *PartitionExtClause // PARTITION/SUBPARTITION extension
	Dblink       string              // @dblink name
	Alias        *Alias              // table alias
	SetClauses   *List               // list of *SetClause
	FromClause   *List               // FROM clause (list of table references)
	WhereClause  ExprNode            // WHERE condition
	Returning    *List               // RETURNING INTO
	ErrorLog     *ErrorLogClause     // LOG ERRORS
	Hints        *List               // optimizer hints
	Loc          Loc                 // start location
}

func (n *UpdateStmt) nodeTag()  {}
func (n *UpdateStmt) stmtNode() {}

// SetClause represents a SET clause in UPDATE.
type SetClause struct {
	Column  *ColumnRef // column to set
	Columns *List      // multi-column set: (col1, col2) = (subquery)
	Value   ExprNode   // value expression or subquery
	Loc     Loc        // start location
}

func (n *SetClause) nodeTag() {}

// ---------------------------------------------------------------------------
// DELETE statement
// ---------------------------------------------------------------------------

// DeleteStmt represents a DELETE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/DELETE.html
type DeleteStmt struct {
	Table        *ObjectName         // table to delete from
	PartitionExt *PartitionExtClause // PARTITION/SUBPARTITION extension
	Dblink       string              // @dblink name
	Alias        *Alias              // table alias
	WhereClause  ExprNode            // WHERE condition
	Returning    *List               // RETURNING INTO
	ErrorLog     *ErrorLogClause     // LOG ERRORS
	Hints        *List               // optimizer hints
	Loc          Loc                 // start location
}

func (n *DeleteStmt) nodeTag()  {}
func (n *DeleteStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// MERGE statement
// ---------------------------------------------------------------------------

// MergeStmt represents a MERGE INTO statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/MERGE.html
type MergeStmt struct {
	Target      *ObjectName     // target table
	TargetAlias *Alias          // target alias
	Source      TableExpr       // source table or subquery
	SourceAlias *Alias          // source alias
	On          ExprNode        // ON condition
	Clauses     *List           // list of *MergeClause (WHEN MATCHED / NOT MATCHED)
	ErrorLog    *ErrorLogClause // LOG ERRORS
	Hints       *List           // optimizer hints
	Loc         Loc             // start location
}

func (n *MergeStmt) nodeTag()  {}
func (n *MergeStmt) stmtNode() {}

// MergeClause represents a WHEN MATCHED / WHEN NOT MATCHED clause.
type MergeClause struct {
	Matched     bool     // true for WHEN MATCHED, false for WHEN NOT MATCHED
	Condition   ExprNode // AND condition (nil if absent)
	UpdateSet   *List    // SET clauses for UPDATE (list of *SetClause)
	UpdateWhere ExprNode // WHERE condition after UPDATE SET
	DeleteWhere ExprNode // DELETE WHERE condition (inside merge_update_clause)
	InsertCols  *List    // INSERT column list
	InsertVals  *List    // INSERT values
	InsertWhere ExprNode // WHERE condition after INSERT VALUES
	IsDelete    bool     // DELETE (no SET)
	Loc         Loc      // start location
}

func (n *MergeClause) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE TABLE statement
// ---------------------------------------------------------------------------

// CreateTableStmt represents a CREATE TABLE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLE.html
type CreateTableStmt struct {
	OrReplace        bool             // OR REPLACE (23c)
	Global           bool             // GLOBAL TEMPORARY
	Private          bool             // PRIVATE TEMPORARY
	Sharded          bool             // SHARDED
	Duplicated       bool             // DUPLICATED
	Name             *ObjectName      // table name
	Sharing          string           // SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
	Columns          *List            // column definitions (list of *ColumnDef)
	Constraints      *List            // table-level constraints (list of *TableConstraint)
	AsQuery          StmtNode         // AS subquery (CTAS)
	Tablespace       string           // TABLESPACE name
	Storage          *StorageClause   // storage parameters
	Partition        *PartitionClause // partitioning
	OnCommit         string           // ON COMMIT (PRESERVE/DELETE ROWS)
	Parallel         string           // PARALLEL/NOPARALLEL
	Compress         string           // COMPRESS/NOCOMPRESS
	IfNotExists      bool             // IF NOT EXISTS (23c)
	Hints            *List            // optimizer hints
	Collation        string           // DEFAULT COLLATION collation_name
	Organization     string           // ORGANIZATION HEAP/INDEX/EXTERNAL
	SegmentCreation  string           // SEGMENT CREATION IMMEDIATE/DEFERRED
	Logging          string           // LOGGING/NOLOGGING
	Cache            string           // CACHE/NOCACHE
	RowMovement      string           // ENABLE/DISABLE ROW MOVEMENT
	ReadOnly         string           // READ ONLY/READ WRITE
	Indexing         string           // INDEXING ON/OFF
	ResultCache      string           // RESULT_CACHE MODE DEFAULT/FORCE
	RowDependencies  string           // ROWDEPENDENCIES/NOROWDEPENDENCIES
	FlashbackArchive string           // FLASHBACK ARCHIVE name / NO FLASHBACK ARCHIVE
	ImmutableNoDrop  string           // NO DROP [UNTIL N DAYS IDLE]
	ImmutableNoDel   string           // NO DELETE [LOCKED | UNTIL N DAYS AFTER INSERT]
	BlockchainHash   string           // HASHING USING 'hash_algorithm'
	BlockchainVer    string           // VERSION 'version_string'
	MemoptimizeRead  bool             // MEMOPTIMIZE FOR READ
	MemoptimizeWrite bool             // MEMOPTIMIZE FOR WRITE
	Parent           *ObjectName      // PARENT table
	Loc              Loc              // start location
}

func (n *CreateTableStmt) nodeTag()  {}
func (n *CreateTableStmt) stmtNode() {}

// ColumnDef represents a column definition.
type ColumnDef struct {
	Name                   string          // column name
	TypeName               *TypeName       // data type
	Domain                 *ObjectName     // DOMAIN domain_name (instead of type)
	Sort                   bool            // SORT
	Visible                bool            // VISIBLE (explicit)
	Default                ExprNode        // DEFAULT expression
	DefaultOnNull          bool            // DEFAULT ON NULL
	DefaultOnNullInsertOnly bool           // DEFAULT ON NULL FOR INSERT ONLY
	Identity               *IdentityClause // GENERATED ... AS IDENTITY
	Virtual                ExprNode        // GENERATED ALWAYS AS (expr) VIRTUAL
	Invisible              bool            // INVISIBLE
	NotNull                bool            // NOT NULL
	Null                   bool            // NULL (explicit)
	Encrypt                string          // ENCRYPT encryption_spec / DECRYPT
	DropIdentity           bool            // DROP IDENTITY (ALTER TABLE MODIFY)
	Constraints            *List           // column constraints (list of *ColumnConstraint)
	Collation              string          // COLLATE
	Loc                    Loc             // start location
}

func (n *ColumnDef) nodeTag() {}

// IdentityClause represents a GENERATED ... AS IDENTITY clause.
type IdentityClause struct {
	Always  bool  // GENERATED ALWAYS vs BY DEFAULT
	Options *List // sequence options
	Loc     Loc
}

func (n *IdentityClause) nodeTag() {}

// ColumnConstraint represents a column-level constraint.
type ColumnConstraint struct {
	Name       string         // constraint name (nil if unnamed)
	Type       ConstraintType // constraint type
	Expr       ExprNode       // CHECK expression
	RefTable   *ObjectName    // REFERENCES table
	RefColumns *List          // REFERENCES columns
	OnDelete   string         // ON DELETE action
	Deferrable bool           // DEFERRABLE
	Initially  string         // INITIALLY DEFERRED/IMMEDIATE
	Loc        Loc            // start location
}

func (n *ColumnConstraint) nodeTag() {}

// TableConstraint represents a table-level constraint.
type TableConstraint struct {
	Name       string         // constraint name
	Type       ConstraintType // constraint type
	Columns    *List          // constraint columns (list of *String)
	Expr       ExprNode       // CHECK expression
	RefTable   *ObjectName    // REFERENCES table
	RefColumns *List          // REFERENCES columns
	OnDelete   string         // ON DELETE action
	Deferrable bool           // DEFERRABLE
	Initially  string         // INITIALLY DEFERRED/IMMEDIATE
	Tablespace string         // USING INDEX TABLESPACE
	Loc        Loc            // start location
}

func (n *TableConstraint) nodeTag() {}

// StorageClause represents Oracle storage parameters.
type StorageClause struct {
	Initial     string // INITIAL
	Next        string // NEXT
	PctIncrease string // PCTINCREASE
	MinExtents  string // MINEXTENTS
	MaxExtents  string // MAXEXTENTS
	PctFree     string // PCTFREE
	PctUsed     string // PCTUSED
	Logging     bool   // LOGGING/NOLOGGING
	Loc         Loc    // start location
}

func (n *StorageClause) nodeTag() {}

// PartitionClause represents a partitioning specification.
type PartitionClause struct {
	Type         PartitionType    // RANGE/LIST/HASH
	Columns      *List            // partition columns
	Partitions   *List            // list of *PartitionDef
	Subpartition *PartitionClause // subpartition template
	Loc          Loc              // start location
}

func (n *PartitionClause) nodeTag() {}

// PartitionDef represents a single partition definition.
type PartitionDef struct {
	Name       string // partition name
	Values     *List  // VALUES LESS THAN / VALUES (for list)
	Tablespace string // TABLESPACE
	Loc        Loc    // start location
}

func (n *PartitionDef) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER TABLE statement
// ---------------------------------------------------------------------------

// AlterTableStmt represents an ALTER TABLE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-TABLE.html
type AlterTableStmt struct {
	Name    *ObjectName // table name
	Actions *List       // list of *AlterTableCmd
	Loc     Loc         // start location
}

func (n *AlterTableStmt) nodeTag()  {}
func (n *AlterTableStmt) stmtNode() {}

// AlterTableCmd represents a single ALTER TABLE action.
type AlterTableCmd struct {
	Action     AlterTableAction // action type
	ColumnDef  *ColumnDef       // for ADD/MODIFY COLUMN
	ColumnDefs *List            // for ADD/MODIFY multiple columns
	ColumnName string           // for DROP/RENAME COLUMN, partition name, property name
	NewName    string           // for RENAME, property value
	Constraint *TableConstraint // for ADD/DROP CONSTRAINT
	Subtype    string           // sub-action qualifier (e.g. "ENABLE"/"DISABLE", "PARTITION"/"SUBPARTITION", property name)
	Loc        Loc              // start location
}

func (n *AlterTableCmd) nodeTag() {}

// ---------------------------------------------------------------------------
// DROP statement
// ---------------------------------------------------------------------------

// DropStmt represents a DROP statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/DROP-TABLE.html
type DropStmt struct {
	ObjectType   ObjectType // what to drop
	Names        *List      // object names (list of *ObjectName)
	IfExists     bool       // IF EXISTS
	Cascade      bool       // CASCADE CONSTRAINTS
	Purge        bool       // PURGE (for tables)
	Online       bool       // ONLINE (for DROP INDEX)
	Force        bool       // FORCE (for DROP INDEX/INDEXTYPE/OPERATOR)
	Validate     bool       // VALIDATE (for DROP TYPE)
	Invalidation string     // "DEFERRED" or "IMMEDIATE" (for DROP INDEX)
	Loc          Loc        // start location
}

func (n *DropStmt) nodeTag()  {}
func (n *DropStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE INDEX statement
// ---------------------------------------------------------------------------

// CreateIndexStmt represents a CREATE INDEX statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEX.html
//
//	CREATE [ UNIQUE ] [ BITMAP ] [ MULTIVALUE ] INDEX [ IF NOT EXISTS ]
//	    [ schema. ] index_name
//	    [ index_ilm_clause ]
//	    { cluster_index_clause | table_index_clause | bitmap_join_index_clause }
//	    [ { DEFERRED | IMMEDIATE } INVALIDATION ]
type CreateIndexStmt struct {
	Unique        bool        // UNIQUE
	Bitmap        bool        // BITMAP
	Multivalue    bool        // MULTIVALUE
	Name          *ObjectName // index name
	Table         *ObjectName // table name (ON table)
	Cluster       *ObjectName // ON CLUSTER name (cluster_index_clause)
	Alias         string      // table alias
	Columns       *List       // index columns (list of *IndexColumn)
	FunctionBased bool        // function-based index
	Reverse       bool        // REVERSE
	Sort          bool        // SORT
	NoSort        bool        // NOSORT
	Local         bool        // LOCAL partitioned
	Global        bool        // GLOBAL partitioned
	Tablespace    string      // TABLESPACE
	Parallel      string      // PARALLEL/NOPARALLEL
	Compress      string      // COMPRESS N
	Online        bool        // ONLINE
	Logging       bool        // LOGGING
	NoLogging     bool        // NOLOGGING
	Visible       bool        // VISIBLE
	Invisible     bool        // INVISIBLE
	IndexType     *ObjectName // INDEXTYPE IS indextype (domain_index_clause)
	Parameters    string      // PARAMETERS ('...')
	Invalidation  string      // "DEFERRED" or "IMMEDIATE" INVALIDATION
	IfNotExists   bool        // IF NOT EXISTS
	// bitmap_join_index_clause
	FromTables  *List    // FROM table_list
	Where       ExprNode // WHERE join_condition
	// Physical attributes
	PctFree     string // PCTFREE integer
	InitTrans   string // INITRANS integer
	MaxTrans    string // MAXTRANS integer
	// Indexing clause
	IndexingFull    bool // INDEXING FULL
	IndexingPartial bool // INDEXING PARTIAL
	Loc             Loc  // start location
}

func (n *CreateIndexStmt) nodeTag()  {}
func (n *CreateIndexStmt) stmtNode() {}

// IndexColumn represents a column in a CREATE INDEX.
type IndexColumn struct {
	Expr      ExprNode    // column reference or expression
	Dir       SortByDir   // ASC/DESC
	NullOrder SortByNulls // NULLS FIRST/LAST
	Loc       Loc         // start location
}

func (n *IndexColumn) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE VIEW / MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// CreateViewStmt represents a CREATE VIEW statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-VIEW.html
type CreateViewStmt struct {
	OrReplace    bool        // OR REPLACE
	IfNotExists  bool        // IF NOT EXISTS
	Force        bool        // FORCE
	NoForce      bool        // NO FORCE
	Editioning   bool        // EDITIONING
	Editionable  string      // "EDITIONABLE" or "NONEDITIONABLE"
	Sharing      string      // METADATA, DATA, EXTENDED DATA, NONE
	Materialized bool        // MATERIALIZED VIEW
	Name         *ObjectName // view name
	Columns      *List       // column aliases
	Query        StmtNode    // AS SELECT ...
	WithCheckOpt bool        // WITH CHECK OPTION
	WithReadOnly bool        // WITH READ ONLY
	ConstraintName string    // CONSTRAINT name for WITH READ ONLY/CHECK OPTION
	Bequeath     string      // "CURRENT_USER" or "DEFINER"
	DefaultCollation string  // DEFAULT COLLATION name
	ContainerMap bool        // CONTAINER_MAP
	ContainersDefault bool   // CONTAINERS_DEFAULT
	// Materialized view specific
	BuildMode     string
	RefreshMode   string
	RefreshMethod string    // FAST, COMPLETE, FORCE
	OnPrebuilt    bool       // ON PREBUILT TABLE
	ReducedPrec   string     // "WITH_REDUCED", "WITHOUT_REDUCED"
	NeverRefresh  bool       // NEVER REFRESH
	WithPK        bool       // WITH PRIMARY KEY
	WithRowID     bool       // WITH ROWID
	RefreshOnStmt bool       // ON STATEMENT
	StartWith     ExprNode   // START WITH
	Next          ExprNode   // NEXT
	EnableQuery   bool
	DisableQuery  bool       // DISABLE QUERY REWRITE
	EnableOnQueryComputation bool
	DisableOnQueryComputation bool
	CacheMode     string     // "CACHE" or "NOCACHE"
	ParallelMode  string     // "PARALLEL" or "NOPARALLEL"
	ParallelDegree string    // degree number
	Options       *List      // remaining generic options
	Loc           Loc        // start location
}

func (n *CreateViewStmt) nodeTag()  {}
func (n *CreateViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE SEQUENCE
// ---------------------------------------------------------------------------

// CreateSequenceStmt represents a CREATE SEQUENCE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-SEQUENCE.html
type CreateSequenceStmt struct {
	Name        *ObjectName // sequence name
	IncrementBy ExprNode    // INCREMENT BY
	StartWith   ExprNode    // START WITH
	MaxValue    ExprNode    // MAXVALUE
	MinValue    ExprNode    // MINVALUE
	NoMaxValue  bool        // NOMAXVALUE
	NoMinValue  bool        // NOMINVALUE
	Cycle       bool        // CYCLE
	NoCycle     bool        // NOCYCLE
	Cache       ExprNode    // CACHE n
	NoCache     bool        // NOCACHE
	Order       bool        // ORDER
	NoOrder     bool        // NOORDER
	IfNotExists bool        // IF NOT EXISTS
	Loc         Loc         // start location
}

func (n *CreateSequenceStmt) nodeTag()  {}
func (n *CreateSequenceStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE SYNONYM
// ---------------------------------------------------------------------------

// CreateSynonymStmt represents a CREATE SYNONYM statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-SYNONYM.html
type CreateSynonymStmt struct {
	OrReplace bool        // OR REPLACE
	Public    bool        // PUBLIC
	Name      *ObjectName // synonym name
	Target    *ObjectName // target object
	Loc       Loc         // start location
}

func (n *CreateSynonymStmt) nodeTag()  {}
func (n *CreateSynonymStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE DATABASE LINK
// ---------------------------------------------------------------------------

// CreateDatabaseLinkStmt represents a CREATE DATABASE LINK statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-DATABASE-LINK.html
type CreateDatabaseLinkStmt struct {
	Public     bool   // PUBLIC
	Name       string // link name
	ConnectTo  string // CONNECT TO user
	Identified string // IDENTIFIED BY password
	Using      string // USING connect string
	Loc        Loc    // start location
}

func (n *CreateDatabaseLinkStmt) nodeTag()  {}
func (n *CreateDatabaseLinkStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE TYPE
// ---------------------------------------------------------------------------

// CreateTypeStmt represents a CREATE TYPE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TYPE.html
type CreateTypeStmt struct {
	OrReplace      bool        // OR REPLACE
	IfNotExists    bool        // IF NOT EXISTS
	Editionable    bool        // EDITIONABLE
	NonEditionable bool        // NONEDITIONABLE
	Name           *ObjectName // type name
	Attributes     *List       // list of *ColumnDef (for object types)
	AsTable        *TypeName   // TABLE OF type (nested table)
	AsVarray       *TypeName   // VARRAY(n) OF type
	VarraySize     ExprNode    // varray size limit
	IsBody         bool        // CREATE TYPE BODY
	Body           *List       // type body members (list of *TypeBodyMember)
	Loc            Loc         // start location
}

func (n *CreateTypeStmt) nodeTag()  {}
func (n *CreateTypeStmt) stmtNode() {}

// TypeBodyMemberKind represents the kind of a type body member.
type TypeBodyMemberKind int

const (
	TYPE_BODY_MEMBER      TypeBodyMemberKind = iota // MEMBER
	TYPE_BODY_STATIC                                // STATIC
	TYPE_BODY_MAP                                   // MAP MEMBER
	TYPE_BODY_ORDER                                 // ORDER MEMBER
	TYPE_BODY_CONSTRUCTOR                           // CONSTRUCTOR
)

// TypeBodyMember represents a single member definition inside a CREATE TYPE BODY.
// It wraps a procedure or function with its kind qualifier.
type TypeBodyMember struct {
	Kind    TypeBodyMemberKind // MEMBER, STATIC, MAP, ORDER, CONSTRUCTOR
	Subprog Node               // *CreateProcedureStmt or *CreateFunctionStmt
	Loc     Loc
}

func (n *TypeBodyMember) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE PACKAGE / PACKAGE BODY
// ---------------------------------------------------------------------------

// CreatePackageStmt represents a CREATE PACKAGE statement.
type CreatePackageStmt struct {
	OrReplace      bool        // OR REPLACE
	IfNotExists    bool        // IF NOT EXISTS
	Editionable    bool        // EDITIONABLE
	NonEditionable bool        // NONEDITIONABLE
	Sharing        string      // SHARING = METADATA | NONE
	Name           *ObjectName // package name
	IsBody         bool        // PACKAGE BODY
	Body           *List       // package body declarations
	Loc            Loc         // start location
}

func (n *CreatePackageStmt) nodeTag()  {}
func (n *CreatePackageStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE PROCEDURE / FUNCTION
// ---------------------------------------------------------------------------

// CreateProcedureStmt represents a CREATE PROCEDURE statement.
type CreateProcedureStmt struct {
	OrReplace      bool        // OR REPLACE
	IfNotExists    bool        // IF NOT EXISTS
	Editionable    bool        // EDITIONABLE
	NonEditionable bool        // NONEDITIONABLE
	Sharing        string      // SHARING = METADATA | NONE
	Name           *ObjectName // procedure name
	Parameters     *List       // parameter list (list of *Parameter)
	Body           StmtNode    // procedure body (PL/SQL block)
	Loc            Loc         // start location
}

func (n *CreateProcedureStmt) nodeTag()  {}
func (n *CreateProcedureStmt) stmtNode() {}

// CreateFunctionStmt represents a CREATE FUNCTION statement.
type CreateFunctionStmt struct {
	OrReplace      bool        // OR REPLACE
	IfNotExists    bool        // IF NOT EXISTS
	Editionable    bool        // EDITIONABLE
	NonEditionable bool        // NONEDITIONABLE
	Sharing        string      // SHARING = METADATA | NONE
	Name           *ObjectName // function name
	Parameters     *List       // parameter list (list of *Parameter)
	ReturnType     *TypeName   // RETURN type
	Deterministic  bool        // DETERMINISTIC
	Pipelined      bool        // PIPELINED
	Parallel       bool        // PARALLEL_ENABLE
	ResultCache    bool        // RESULT_CACHE
	Aggregate      bool        // AGGREGATE USING
	SqlMacro       bool        // SQL_MACRO
	Body           StmtNode    // function body (PL/SQL block)
	Loc            Loc         // start location
}

func (n *CreateFunctionStmt) nodeTag()  {}
func (n *CreateFunctionStmt) stmtNode() {}

// Parameter represents a procedure/function parameter.
type Parameter struct {
	Name     string    // parameter name
	Mode     string    // IN, OUT, IN OUT, NOCOPY
	TypeName *TypeName // parameter type
	Default  ExprNode  // default value
	Loc      Loc       // start location
}

func (n *Parameter) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE TRIGGER
// ---------------------------------------------------------------------------

// CreateTriggerStmt represents a CREATE TRIGGER statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TRIGGER.html
type CreateTriggerStmt struct {
	OrReplace      bool          // OR REPLACE
	IfNotExists    bool          // IF NOT EXISTS
	Editionable    bool          // EDITIONABLE
	NonEditionable bool          // NONEDITIONABLE
	Name           *ObjectName   // trigger name
	Timing         TriggerTiming // BEFORE/AFTER/INSTEAD OF
	Events         *List         // trigger events (list of TriggerEvent)
	Table          *ObjectName   // ON table
	ForEachRow     bool          // FOR EACH ROW
	When           ExprNode      // WHEN condition
	Body           StmtNode      // trigger body
	Compound       bool          // COMPOUND TRIGGER
	Enable         bool          // ENABLE (default true)
	Loc            Loc           // start location
}

func (n *CreateTriggerStmt) nodeTag()  {}
func (n *CreateTriggerStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// TRUNCATE statement
// ---------------------------------------------------------------------------

// TruncateStmt represents a TRUNCATE TABLE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/TRUNCATE-TABLE.html
type TruncateStmt struct {
	Table      *ObjectName // table name
	Cluster    bool        // TRUNCATE CLUSTER
	PurgeMVLog bool        // PURGE MATERIALIZED VIEW LOG
	Cascade    bool        // CASCADE
	Loc        Loc         // start location
}

func (n *TruncateStmt) nodeTag()  {}
func (n *TruncateStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// GRANT / REVOKE
// ---------------------------------------------------------------------------

// GrantStmt represents a GRANT statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/GRANT.html
type GrantStmt struct {
	Privileges *List       // list of *String (privilege names)
	AllPriv    bool        // ALL PRIVILEGES
	OnObject   *ObjectName // ON object (nil for system privileges)
	OnType     ObjectType  // object type
	Grantees   *List       // list of *String (grantee names)
	WithAdmin  bool        // WITH ADMIN OPTION
	WithGrant  bool        // WITH GRANT OPTION
	Loc        Loc         // start location
}

func (n *GrantStmt) nodeTag()  {}
func (n *GrantStmt) stmtNode() {}

// RevokeStmt represents a REVOKE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/REVOKE.html
type RevokeStmt struct {
	Privileges *List       // list of *String (privilege names)
	AllPriv    bool        // ALL PRIVILEGES
	OnObject   *ObjectName // ON object (nil for system privileges)
	OnType     ObjectType  // object type
	Grantees   *List       // list of *String (grantee names)
	Loc        Loc         // start location
}

func (n *RevokeStmt) nodeTag()  {}
func (n *RevokeStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// COMMENT statement
// ---------------------------------------------------------------------------

// CommentStmt represents a COMMENT ON statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/COMMENT.html
type CommentStmt struct {
	ObjectType ObjectType  // type of object
	Object     *ObjectName // object name
	Column     string      // column name (for COMMENT ON COLUMN)
	Comment    string      // comment text
	Loc        Loc         // start location
}

func (n *CommentStmt) nodeTag()  {}
func (n *CommentStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER SESSION / ALTER SYSTEM
// ---------------------------------------------------------------------------

// AlterSessionStmt represents an ALTER SESSION statement.
type AlterSessionStmt struct {
	Action         string // "SET", "ADVISE", "CLOSE_DATABASE_LINK", "ENABLE", "DISABLE", "FORCE_PARALLEL", "SYNC_WITH_PRIMARY"
	SetParams      *List  // list of *SetParam (for SET action)
	AdviseAction   string // "COMMIT", "ROLLBACK", "NOTHING" (for ADVISE action)
	DBLink         string // database link name (for CLOSE DATABASE LINK)
	Feature        string // "COMMIT_IN_PROCEDURE", "GUARD", "PARALLEL_DML", "PARALLEL_DDL", "PARALLEL_QUERY", "RESUMABLE", "SHARD_DDL"
	ParallelDegree int    // PARALLEL integer (for FORCE PARALLEL)
	Timeout        int    // TIMEOUT integer (for RESUMABLE)
	ResumableName  string // NAME 'string' (for RESUMABLE)
	Loc            Loc
}

func (n *AlterSessionStmt) nodeTag()  {}
func (n *AlterSessionStmt) stmtNode() {}

// AlterSystemStmt represents an ALTER SYSTEM statement.
type AlterSystemStmt struct {
	Action          string // "SET", "RESET", "KILL_SESSION", "DISCONNECT_SESSION", "FLUSH", "CHECKPOINT", "SWITCH_LOGFILE", "ARCHIVE_LOG", "REGISTER", "SUSPEND", "RESUME", "QUIESCE", "UNQUIESCE", "CANCEL_SQL", "SHUTDOWN", "RELOCATE_CLIENT", "CHECK_DATAFILES", "ENABLE", "DISABLE", "START_ROLLING_MIGRATION", "STOP_ROLLING_MIGRATION", "START_ROLLING_PATCH", "STOP_ROLLING_PATCH", "SET_ENCRYPTION"
	SetParams       *List  // list of *SetParam (for SET action)
	ResetParam      string // parameter name (for RESET action)
	SessionID       string // 'sid,serial#' (for KILL/DISCONNECT/CANCEL)
	InstanceID      string // @instance_id (for KILL/CANCEL)
	Immediate       bool   // IMMEDIATE flag
	Force           bool   // FORCE flag (for KILL SESSION)
	PostTransaction bool   // POST_TRANSACTION (for DISCONNECT SESSION)
	NoReplay        bool   // NOREPLAY (for KILL SESSION)
	Timeout         int    // TIMEOUT integer (for KILL SESSION)
	FlushTarget     string // "SHARED_POOL", "GLOBAL_CONTEXT", "BUFFER_CACHE", "FLASH_CACHE", "REDO", "PASSWORDFILE_METADATA_CACHE"
	FlushScope      string // "GLOBAL", "LOCAL" (for BUFFER_CACHE/FLASH_CACHE)
	FlushRedoDB     string // target_db_name (for REDO flush)
	FlushRedoConfirm string // "NO_CONFIRM_APPLY", "CONFIRM_APPLY" (for REDO flush)
	ArchiveLogSpec  string // "SEQUENCE", "CHANGE", "CURRENT", "GROUP", "LOGFILE", "NEXT", "ALL"
	ArchiveLogValue string // integer/filename for archive log
	ArchiveInstance string // INSTANCE 'instance_name'
	ArchiveThread   int    // THREAD integer
	ArchiveNoSwitch bool   // NOSWITCH (for CURRENT)
	ArchiveBackupCF bool   // USING BACKUP CONTROLFILE (for LOGFILE)
	ArchiveTo       string // TO 'location'
	CheckScope      string // "GLOBAL", "LOCAL" (for CHECKPOINT/CHECK DATAFILES)
	Feature         string // "DISTRIBUTED_RECOVERY", "RESTRICTED_SESSION" (for ENABLE/DISABLE)
	ShutdownDisp    string // dispatcher name (for SHUTDOWN)
	RelocateClient  string // client_id (for RELOCATE CLIENT)
	RollingVersion  string // ASM_version (for START ROLLING MIGRATION)
	EncryptionAction string // "OPEN", "CLOSE", "SET_KEY" (for SET ENCRYPTION WALLET/KEY)
	Comment         string // COMMENT = 'text' (for SET)
	Deferred        bool   // DEFERRED (for SET)
	Scope           string // "MEMORY", "SPFILE", "BOTH" (for SET/RESET)
	SID             string // SID = 'sid' | '*' (for SET/RESET)
	Container       string // "ALL", "CURRENT" (for SET)
	SqlID           string // SQL_ID 'sql_id' (for CANCEL SQL)
	Loc             Loc
}

func (n *AlterSystemStmt) nodeTag()  {}
func (n *AlterSystemStmt) stmtNode() {}

// SetParam represents a parameter setting (name = value).
type SetParam struct {
	Name  string   // parameter name
	Value ExprNode // parameter value
	Loc   Loc
}

func (n *SetParam) nodeTag() {}

// ---------------------------------------------------------------------------
// Transaction statements
// ---------------------------------------------------------------------------

// CommitStmt represents a COMMIT statement.
type CommitStmt struct {
	Work         bool   // WORK keyword
	Comment      string // COMMENT 'text'
	WriteWait    string // "WAIT" or "NOWAIT"
	WriteBatch   string // "IMMEDIATE" or "BATCH"
	Force        string // FORCE 'text'
	ForceInteger int    // optional integer after FORCE string
	HasForceInt  bool   // whether ForceInteger was specified
	Loc          Loc
}

func (n *CommitStmt) nodeTag()  {}
func (n *CommitStmt) stmtNode() {}

// RollbackStmt represents a ROLLBACK statement.
type RollbackStmt struct {
	Work        bool   // WORK keyword
	ToSavepoint string // TO SAVEPOINT name
	Force       string // FORCE 'text'
	Loc         Loc    // start location
}

func (n *RollbackStmt) nodeTag()  {}
func (n *RollbackStmt) stmtNode() {}

// SavepointStmt represents a SAVEPOINT statement.
type SavepointStmt struct {
	Name string // savepoint name
	Loc  Loc
}

func (n *SavepointStmt) nodeTag()  {}
func (n *SavepointStmt) stmtNode() {}

// SetTransactionStmt represents a SET TRANSACTION statement.
type SetTransactionStmt struct {
	ReadOnly            bool   // READ ONLY
	ReadWrite           bool   // READ WRITE
	IsolLevel           string // ISOLATION LEVEL
	UseRollbackSegment  string // USE ROLLBACK SEGMENT segment_name
	Name                string // NAME 'text'
	Loc                 Loc    // start location
}

func (n *SetTransactionStmt) nodeTag()  {}
func (n *SetTransactionStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// PL/SQL block statements
// ---------------------------------------------------------------------------

// PLSQLBlock represents a PL/SQL block (DECLARE ... BEGIN ... END).
type PLSQLBlock struct {
	Declarations *List  // DECLARE section
	Statements   *List  // BEGIN ... END statements
	Exceptions   *List  // EXCEPTION section (list of *ExceptionHandler)
	Label        string // block label
	Loc          Loc    // start location
}

func (n *PLSQLBlock) nodeTag()  {}
func (n *PLSQLBlock) stmtNode() {}

// ExceptionHandler represents an exception handler in a PL/SQL block.
type ExceptionHandler struct {
	Exceptions *List // exception names (list of *String)
	Statements *List // handler statements
	Loc        Loc   // start location
}

func (n *ExceptionHandler) nodeTag() {}

// PLSQLIf represents an IF/ELSIF/ELSE statement.
type PLSQLIf struct {
	Condition ExprNode // IF condition
	Then      *List    // THEN statements
	ElsIfs    *List    // list of *PLSQLElsIf
	Else      *List    // ELSE statements
	Loc       Loc      // start location
}

func (n *PLSQLIf) nodeTag()  {}
func (n *PLSQLIf) stmtNode() {}

// PLSQLElsIf represents an ELSIF clause.
type PLSQLElsIf struct {
	Condition ExprNode // ELSIF condition
	Then      *List    // THEN statements
	Loc       Loc      // start location
}

func (n *PLSQLElsIf) nodeTag() {}

// PLSQLLoop represents a LOOP/WHILE/FOR statement.
type PLSQLLoop struct {
	Type       PLSQLLoopType // loop type
	Label      string        // loop label
	Condition  ExprNode      // WHILE condition
	Iterator   string        // FOR iterator variable
	LowerBound ExprNode      // FOR lower bound
	UpperBound ExprNode      // FOR upper bound
	Reverse    bool          // REVERSE
	CursorName string        // cursor name (for cursor FOR loop)
	CursorArgs *List         // cursor arguments
	Statements *List         // loop body
	Loc        Loc           // start location
}

func (n *PLSQLLoop) nodeTag()  {}
func (n *PLSQLLoop) stmtNode() {}

// PLSQLReturn represents a RETURN statement.
type PLSQLReturn struct {
	Expr ExprNode // return expression (nil for procedures)
	Loc  Loc
}

func (n *PLSQLReturn) nodeTag()  {}
func (n *PLSQLReturn) stmtNode() {}

// PLSQLGoto represents a GOTO statement.
type PLSQLGoto struct {
	Label string // target label
	Loc   Loc
}

func (n *PLSQLGoto) nodeTag()  {}
func (n *PLSQLGoto) stmtNode() {}

// PLSQLAssign represents a PL/SQL assignment statement.
type PLSQLAssign struct {
	Target ExprNode // assignment target
	Value  ExprNode // assignment value
	Loc    Loc
}

func (n *PLSQLAssign) nodeTag()  {}
func (n *PLSQLAssign) stmtNode() {}

// PLSQLRaise represents a RAISE statement.
type PLSQLRaise struct {
	Exception string // exception name (empty for re-raise)
	Loc       Loc    // start location
}

func (n *PLSQLRaise) nodeTag()  {}
func (n *PLSQLRaise) stmtNode() {}

// PLSQLNull represents a NULL statement (no-op).
type PLSQLNull struct {
	Loc Loc
}

func (n *PLSQLNull) nodeTag()  {}
func (n *PLSQLNull) stmtNode() {}

// PLSQLVarDecl represents a variable declaration.
type PLSQLVarDecl struct {
	Name     string    // variable name
	TypeName *TypeName // variable type
	Constant bool      // CONSTANT
	NotNull  bool      // NOT NULL
	Default  ExprNode  // DEFAULT / := value
	Loc      Loc       // start location
}

func (n *PLSQLVarDecl) nodeTag() {}

// PLSQLCursorDecl represents a cursor declaration.
type PLSQLCursorDecl struct {
	Name       string   // cursor name
	Parameters *List    // parameter list
	Query      StmtNode // SELECT statement
	Loc        Loc      // start location
}

func (n *PLSQLCursorDecl) nodeTag() {}

// PLSQLExecImmediate represents an EXECUTE IMMEDIATE statement.
type PLSQLExecImmediate struct {
	SQL   ExprNode // SQL string expression
	Into  *List    // INTO variable list
	Using *List    // USING bind variable list
	Loc   Loc
}

func (n *PLSQLExecImmediate) nodeTag()  {}
func (n *PLSQLExecImmediate) stmtNode() {}

// PLSQLOpen represents an OPEN cursor statement.
type PLSQLOpen struct {
	Cursor   string   // cursor name
	Args     *List    // cursor arguments
	ForQuery StmtNode // FOR select statement (ref cursor)
	Loc      Loc
}

func (n *PLSQLOpen) nodeTag()  {}
func (n *PLSQLOpen) stmtNode() {}

// PLSQLFetch represents a FETCH cursor statement.
type PLSQLFetch struct {
	Cursor string   // cursor name
	Into   *List    // INTO variable list
	Bulk   bool     // BULK COLLECT
	Limit  ExprNode // LIMIT expression
	Loc    Loc
}

func (n *PLSQLFetch) nodeTag()  {}
func (n *PLSQLFetch) stmtNode() {}

// PLSQLClose represents a CLOSE cursor statement.
type PLSQLClose struct {
	Cursor string // cursor name
	Loc    Loc
}

func (n *PLSQLClose) nodeTag()  {}
func (n *PLSQLClose) stmtNode() {}

// PLSQLForall represents a FORALL statement.
//
//	FORALL index IN lower..upper [SAVE EXCEPTIONS] dml_statement
type PLSQLForall struct {
	Index string   // index variable
	Lower ExprNode // lower bound
	Upper ExprNode // upper bound
	Body  StmtNode // DML statement
	Loc   Loc
}

func (n *PLSQLForall) nodeTag()  {}
func (n *PLSQLForall) stmtNode() {}

// PLSQLExit represents an EXIT [label] [WHEN condition] statement.
type PLSQLExit struct {
	Label     string   // optional label
	Condition ExprNode // WHEN condition
	Loc       Loc
}

func (n *PLSQLExit) nodeTag()  {}
func (n *PLSQLExit) stmtNode() {}

// PLSQLContinue represents a CONTINUE [label] [WHEN condition] statement.
type PLSQLContinue struct {
	Label     string   // optional label
	Condition ExprNode // WHEN condition
	Loc       Loc
}

func (n *PLSQLContinue) nodeTag()  {}
func (n *PLSQLContinue) stmtNode() {}

// PLSQLPipeRow represents a PIPE ROW statement for pipelined functions.
type PLSQLPipeRow struct {
	Row ExprNode // the row expression
	Loc Loc
}

func (n *PLSQLPipeRow) nodeTag()  {}
func (n *PLSQLPipeRow) stmtNode() {}

// PLSQLPragma represents a PRAGMA directive.
type PLSQLPragma struct {
	Name string // AUTONOMOUS_TRANSACTION, EXCEPTION_INIT, etc.
	Args *List  // optional arguments
	Loc  Loc
}

func (n *PLSQLPragma) nodeTag() {}

// PLSQLCase represents a PL/SQL CASE statement (distinct from CASE expression).
type PLSQLCase struct {
	Expr  ExprNode    // search expression (nil for searched CASE)
	Whens []*PLSQLWhen // WHEN clauses
	Else  []StmtNode  // ELSE statements
	Loc   Loc
}

func (n *PLSQLCase) nodeTag()  {}
func (n *PLSQLCase) stmtNode() {}

// PLSQLWhen represents a WHEN clause in a PL/SQL CASE statement.
type PLSQLWhen struct {
	Expr  ExprNode   // WHEN expression
	Stmts []StmtNode // THEN statements
	Loc   Loc
}

func (n *PLSQLWhen) nodeTag() {}

// PLSQLTypeDecl represents a PL/SQL TYPE declaration.
//
//	TYPE name IS TABLE OF type [INDEX BY type]
//	TYPE name IS VARRAY(n) OF type
//	TYPE name IS RECORD (field type [,...])
//	TYPE name IS REF CURSOR [RETURN type]
type PLSQLTypeDecl struct {
	Name      string    // type name
	Kind      string    // TABLE, VARRAY, RECORD, REF_CURSOR
	ElementType *TypeName // element type (TABLE OF/VARRAY OF)
	IndexBy   *TypeName // INDEX BY type (associative arrays)
	Limit     ExprNode  // VARRAY limit
	Fields    *List     // RECORD fields
	ReturnType *TypeName // REF CURSOR RETURN type
	Loc       Loc
}

func (n *PLSQLTypeDecl) nodeTag() {}

// PLSQLCall represents a PL/SQL standalone procedure call statement.
type PLSQLCall struct {
	Name ExprNode // procedure/function reference (could be dotted name)
	Loc  Loc
}

func (n *PLSQLCall) nodeTag()  {}
func (n *PLSQLCall) stmtNode() {}

// ---------------------------------------------------------------------------
// Utility statements
// ---------------------------------------------------------------------------

// AnalyzeStmt represents an ANALYZE statement.
type AnalyzeStmt struct {
	Table          *ObjectName // table/index/cluster to analyze
	ObjectType     ObjectType  // TABLE, INDEX, CLUSTER
	Action         string      // COMPUTE STATISTICS, ESTIMATE STATISTICS, DELETE STATISTICS, VALIDATE STRUCTURE, VALIDATE REF UPDATE, LIST CHAINED ROWS
	DeleteSystem   bool        // DELETE SYSTEM STATISTICS
	SampleValue    int         // ESTIMATE STATISTICS SAMPLE n {ROWS|PERCENT}
	SampleUnit     string      // ROWS or PERCENT
	CascadeFast    bool        // VALIDATE STRUCTURE CASCADE FAST
	IntoTable      *ObjectName // INTO table (VALIDATE STRUCTURE / LIST CHAINED ROWS)
	SetDanglingNull bool       // VALIDATE REF UPDATE SET DANGLING TO NULL
	Online         bool        // VALIDATE STRUCTURE ONLINE
	Offline        bool        // VALIDATE STRUCTURE OFFLINE
	Loc            Loc         // start location
}

func (n *AnalyzeStmt) nodeTag()  {}
func (n *AnalyzeStmt) stmtNode() {}

// ExplainPlanStmt represents an EXPLAIN PLAN statement.
type ExplainPlanStmt struct {
	StatementID string      // SET STATEMENT_ID
	Into        *ObjectName // INTO table
	Statement   StmtNode    // the statement to explain
	Loc         Loc         // start location
}

func (n *ExplainPlanStmt) nodeTag()  {}
func (n *ExplainPlanStmt) stmtNode() {}

// FlashbackTableStmt represents a FLASHBACK TABLE statement.
type FlashbackTableStmt struct {
	Tables         []*ObjectName // one or more table names
	Table          *ObjectName   // table name (kept for backward compat; first entry)
	ToSCN          ExprNode      // TO SCN expression
	ToTimestamp    ExprNode      // TO TIMESTAMP expression
	ToRestorePoint string        // TO RESTORE POINT name
	Before         bool          // true if TO BEFORE { SCN | TIMESTAMP }
	ToBeforeDrop   bool          // TO BEFORE DROP
	Rename         string        // RENAME TO name
	EnableTriggers *bool         // nil = not specified, true = ENABLE, false = DISABLE
	Loc            Loc           // start location
}

func (n *FlashbackTableStmt) nodeTag()  {}
func (n *FlashbackTableStmt) stmtNode() {}

// FlashbackDatabaseStmt represents a FLASHBACK DATABASE statement.
//
// BNF: oracle/parser/bnf/FLASHBACK-DATABASE.bnf
//
//	FLASHBACK [ STANDBY | PLUGGABLE ] DATABASE [ database ]
//	    { TO SCN scn_number
//	    | TO BEFORE SCN scn_number
//	    | TO TIMESTAMP timestamp_expression
//	    | TO BEFORE TIMESTAMP timestamp_expression
//	    | TO RESTORE POINT restore_point_name
//	    | TO BEFORE RESETLOGS
//	    } ;
type FlashbackDatabaseStmt struct {
	Modifier       string      // "STANDBY", "PLUGGABLE", or ""
	DatabaseName   *ObjectName // optional database name
	Before         bool        // true if TO BEFORE variant
	ToSCN          ExprNode    // SCN expression
	ToTimestamp    ExprNode    // TIMESTAMP expression
	ToRestorePoint string      // RESTORE POINT name
	ToResetlogs    bool        // TO BEFORE RESETLOGS
	Loc            Loc
}

func (n *FlashbackDatabaseStmt) nodeTag()  {}
func (n *FlashbackDatabaseStmt) stmtNode() {}

// PurgeStmt represents a PURGE statement.
type PurgeStmt struct {
	ObjectType ObjectType  // TABLE, INDEX, RECYCLEBIN, DBA_RECYCLEBIN
	Name       *ObjectName // object name
	Loc        Loc         // start location
}

func (n *PurgeStmt) nodeTag()  {}
func (n *PurgeStmt) stmtNode() {}

// LockTableItem represents a single table entry in a LOCK TABLE statement,
// including optional partition_extension_clause.
type LockTableItem struct {
	Table         *ObjectName // table/view name (with optional schema and dblink)
	PartitionType string      // "PARTITION" or "SUBPARTITION" (empty if none)
	PartitionFor  bool        // true if PARTITION FOR / SUBPARTITION FOR
	PartitionName string      // partition or subpartition name / key value
	Loc           Loc
}

func (n *LockTableItem) nodeTag() {}

// LockTableStmt represents a LOCK TABLE statement.
//
//	LOCK TABLE [schema.]table [partition_extension_clause]
//	    [, [schema.]table [partition_extension_clause]]...
//	    IN lock_mode MODE [NOWAIT | WAIT n]
type LockTableStmt struct {
	Tables   []*LockTableItem // one or more table entries
	LockMode string           // ROW SHARE, ROW EXCLUSIVE, SHARE, etc.
	Nowait   bool             // NOWAIT
	Wait     ExprNode         // WAIT n
	Loc      Loc
}

func (n *LockTableStmt) nodeTag()  {}
func (n *LockTableStmt) stmtNode() {}

// CallStmt represents a CALL statement.
//
//	CALL [schema.]routine_name ( [args] ) [INTO :bind_variable]
type CallStmt struct {
	Name *ObjectName // routine name
	Args *List       // argument list
	Into ExprNode    // INTO bind variable
	Loc  Loc
}

func (n *CallStmt) nodeTag()  {}
func (n *CallStmt) stmtNode() {}

// RenameStmt represents a RENAME statement.
//
//	RENAME old_name TO new_name
type RenameStmt struct {
	OldName *ObjectName // old name
	NewName *ObjectName // new name
	Loc     Loc
}

func (n *RenameStmt) nodeTag()  {}
func (n *RenameStmt) stmtNode() {}

// IdentifiedClauseType represents the type of IDENTIFIED clause.
type IdentifiedClauseType int

const (
	IDENTIFIED_BY        IdentifiedClauseType = iota // IDENTIFIED BY password
	IDENTIFIED_EXTERNALLY                            // IDENTIFIED EXTERNALLY [AS '...']
	IDENTIFIED_GLOBALLY                              // IDENTIFIED GLOBALLY [AS '...']
	IDENTIFIED_NO_AUTH                               // NO AUTHENTICATION
)

// IdentifiedClause represents an IDENTIFIED BY/EXTERNALLY/GLOBALLY clause.
type IdentifiedClause struct {
	Type       IdentifiedClauseType
	Password   string // password for IDENTIFIED BY
	OldPass    string // REPLACE old_password (ALTER USER only)
	ExternalAs string // AS 'certificate_DN' etc. for EXTERNALLY/GLOBALLY
	Loc        Loc
}

func (n *IdentifiedClause) nodeTag() {}

// UserQuotaClause represents a QUOTA clause for a user.
type UserQuotaClause struct {
	Size       string      // size value (e.g. "10M") or "UNLIMITED"
	Tablespace *ObjectName // ON tablespace
	Loc        Loc
}

func (n *UserQuotaClause) nodeTag() {}

// DefaultRoleClause represents a DEFAULT ROLE clause for ALTER USER.
type DefaultRoleClause struct {
	AllRoles  bool          // ALL
	NoneRole  bool          // NONE
	Roles     []*ObjectName // specific role list
	ExceptAll bool          // ALL EXCEPT ...
	Loc       Loc
}

func (n *DefaultRoleClause) nodeTag() {}

// CreateUserStmt represents a CREATE USER statement.
//
//	CREATE USER [IF NOT EXISTS] user
//	    IDENTIFIED { BY password | EXTERNALLY [AS '...'] | GLOBALLY [AS '...'] | NO AUTHENTICATION }
//	    [DEFAULT COLLATION collation_name]
//	    [DEFAULT TABLESPACE tablespace]
//	    [[LOCAL] TEMPORARY TABLESPACE { tablespace | tablespace_group_name }]
//	    [QUOTA { size_clause | UNLIMITED } ON tablespace] ...
//	    [PROFILE profile_name]
//	    [PASSWORD EXPIRE]
//	    [ACCOUNT { LOCK | UNLOCK }]
//	    [ENABLE EDITIONS]
//	    [CONTAINER = { ALL | CURRENT }]
type CreateUserStmt struct {
	Name              *ObjectName      // user name
	IfNotExists       bool             // IF NOT EXISTS
	Identified        *IdentifiedClause // IDENTIFIED clause
	DefaultTablespace string            // DEFAULT TABLESPACE
	TempTablespace    string            // [LOCAL] TEMPORARY TABLESPACE
	LocalTemp         bool              // LOCAL keyword present
	Quotas            []*UserQuotaClause // QUOTA clauses
	Profile           string            // PROFILE name
	PasswordExpire    bool              // PASSWORD EXPIRE
	AccountLock       *bool             // ACCOUNT LOCK (true) / UNLOCK (false) / nil (not specified)
	EnableEditions    bool              // ENABLE EDITIONS
	DefaultCollation  string            // DEFAULT COLLATION
	ContainerAll      *bool             // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc               Loc
}

func (n *CreateUserStmt) nodeTag()  {}
func (n *CreateUserStmt) stmtNode() {}

// AlterUserStmt represents an ALTER USER statement.
//
//	ALTER USER [IF EXISTS] user
//	    [IDENTIFIED { BY password [REPLACE old] | EXTERNALLY | GLOBALLY AS '...' | NO AUTHENTICATION }]
//	    [DEFAULT COLLATION collation_name]
//	    [DEFAULT TABLESPACE tablespace]
//	    [[LOCAL] TEMPORARY TABLESPACE { tablespace | tablespace_group_name }]
//	    [QUOTA { size_clause | UNLIMITED } ON tablespace] ...
//	    [PROFILE profile_name]
//	    [DEFAULT ROLE { role [, role]... | ALL [EXCEPT role [, role]...] | NONE }]
//	    [PASSWORD EXPIRE]
//	    [ACCOUNT { LOCK | UNLOCK }]
//	    [ENABLE EDITIONS]
//	    [CONTAINER = { ALL | CURRENT }]
type AlterUserStmt struct {
	Name              *ObjectName        // user name
	IfExists          bool               // IF EXISTS
	Identified        *IdentifiedClause  // IDENTIFIED clause
	DefaultTablespace string             // DEFAULT TABLESPACE
	TempTablespace    string             // [LOCAL] TEMPORARY TABLESPACE
	LocalTemp         bool               // LOCAL keyword present
	Quotas            []*UserQuotaClause // QUOTA clauses
	Profile           string             // PROFILE name
	DefaultRole       *DefaultRoleClause // DEFAULT ROLE clause
	PasswordExpire    bool               // PASSWORD EXPIRE
	AccountLock       *bool              // ACCOUNT LOCK (true) / UNLOCK (false) / nil (not specified)
	EnableEditions    bool               // ENABLE EDITIONS
	DefaultCollation  string             // DEFAULT COLLATION
	ContainerAll      *bool              // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc               Loc
}

func (n *AlterUserStmt) nodeTag()  {}
func (n *AlterUserStmt) stmtNode() {}

// RoleIdentifiedType represents the type of role identification.
type RoleIdentifiedType int

const (
	ROLE_NOT_IDENTIFIED  RoleIdentifiedType = iota // NOT IDENTIFIED
	ROLE_IDENTIFIED_BY                             // IDENTIFIED BY password
	ROLE_IDENTIFIED_USING                          // IDENTIFIED USING [schema.]package
	ROLE_IDENTIFIED_EXTERNALLY                     // IDENTIFIED EXTERNALLY
	ROLE_IDENTIFIED_GLOBALLY                       // IDENTIFIED GLOBALLY [AS 'directory_name']
)

// CreateRoleStmt represents a CREATE ROLE statement.
//
//	CREATE ROLE role
//	    [ NOT IDENTIFIED
//	    | IDENTIFIED { BY password
//	                 | USING [ schema. ] package
//	                 | EXTERNALLY
//	                 | GLOBALLY [ AS 'directory_name' ]
//	                 }
//	    ]
//	    [ CONTAINER = { ALL | CURRENT } ]
type CreateRoleStmt struct {
	Name          *ObjectName        // role name
	IdentifiedType RoleIdentifiedType // identification type
	IdentifyBy    string             // password, package name, or global AS value
	IdentifySchema string            // schema for USING clause
	HasIdentified bool               // true if any IDENTIFIED/NOT IDENTIFIED clause present
	ContainerAll  *bool              // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc           Loc
}

func (n *CreateRoleStmt) nodeTag()  {}
func (n *CreateRoleStmt) stmtNode() {}

// AlterRoleStmt represents an ALTER ROLE statement.
//
//	ALTER ROLE role
//	    { NOT IDENTIFIED
//	    | IDENTIFIED BY password
//	    | IDENTIFIED EXTERNALLY
//	    | IDENTIFIED GLOBALLY AS 'domain_name'
//	    | IDENTIFIED USING [ schema. ] package_name
//	    }
//	    [ CONTAINER = { ALL | CURRENT } ]
type AlterRoleStmt struct {
	Name           *ObjectName        // role name
	IdentifiedType RoleIdentifiedType // identification type
	IdentifyBy     string             // password, package name, or global AS value
	IdentifySchema string             // schema for USING clause
	ContainerAll   *bool              // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc            Loc
}

func (n *AlterRoleStmt) nodeTag()  {}
func (n *AlterRoleStmt) stmtNode() {}

// ProfileLimitType represents the type of profile limit.
type ProfileLimitType int

const (
	PROFILE_RESOURCE ProfileLimitType = iota // resource parameter
	PROFILE_PASSWORD                         // password parameter
)

// ProfileLimit represents a single profile limit entry.
type ProfileLimit struct {
	Name  string // parameter name (e.g. "SESSIONS_PER_USER", "FAILED_LOGIN_ATTEMPTS")
	Value string // value: integer, size clause, "UNLIMITED", "DEFAULT", "NULL", or function name
	Loc   Loc
}

func (n *ProfileLimit) nodeTag() {}

// CreateProfileStmt represents a CREATE PROFILE statement.
//
//	CREATE [MANDATORY] PROFILE profile_name
//	    LIMIT { resource_parameters | password_parameters } ...
//	    [CONTAINER = { ALL | CURRENT }]
type CreateProfileStmt struct {
	Name         *ObjectName     // profile name
	Mandatory    bool            // MANDATORY keyword
	Limits       []*ProfileLimit // parsed limit entries
	ContainerAll *bool           // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc          Loc
}

func (n *CreateProfileStmt) nodeTag()  {}
func (n *CreateProfileStmt) stmtNode() {}

// AlterProfileStmt represents an ALTER PROFILE statement.
//
//	ALTER PROFILE profile_name
//	    LIMIT { resource_parameters | password_parameters }
//	    [ CONTAINER = { ALL | CURRENT } ]
type AlterProfileStmt struct {
	Name         *ObjectName     // profile name
	Limits       []*ProfileLimit // parsed limit entries
	ContainerAll *bool           // CONTAINER = ALL (true) / CURRENT (false) / nil (not specified)
	Loc          Loc
}

func (n *AlterProfileStmt) nodeTag()  {}
func (n *AlterProfileStmt) stmtNode() {}

// ResourceCostEntry represents a single resource cost entry in ALTER RESOURCE COST.
type ResourceCostEntry struct {
	Name  string // CPU_PER_SESSION, CONNECT_TIME, LOGICAL_READS_PER_SESSION, PRIVATE_SGA
	Value string // integer value
	Loc   Loc
}

func (n *ResourceCostEntry) nodeTag() {}

// AlterResourceCostStmt represents an ALTER RESOURCE COST statement.
//
//	ALTER RESOURCE COST
//	    { CPU_PER_SESSION integer
//	    | CONNECT_TIME integer
//	    | LOGICAL_READS_PER_SESSION integer
//	    | PRIVATE_SGA integer
//	    } ...
type AlterResourceCostStmt struct {
	Costs []*ResourceCostEntry
	Loc   Loc
}

func (n *AlterResourceCostStmt) nodeTag()  {}
func (n *AlterResourceCostStmt) stmtNode() {}

// AdminDDLStmt is a generic statement node for administrative DDL statements
// (CREATE/ALTER/DROP TABLESPACE, DIRECTORY, CONTEXT, CLUSTER, DIMENSION,
// FLASHBACK ARCHIVE, JAVA, LIBRARY) that captures the DDL action, object type,
// and name without detailed parsing of all options.
type AdminDDLStmt struct {
	Action     string      // CREATE, ALTER, DROP
	ObjectType ObjectType  // OBJECT_TABLESPACE, OBJECT_DIRECTORY, etc.
	Name       *ObjectName // object name
	OrReplace  bool        // OR REPLACE
	IfExists   bool        // IF EXISTS (for DROP)
	Options    *List       // parsed DDL options (list of *DDLOption)
	Loc        Loc
}

func (n *AdminDDLStmt) nodeTag()  {}
func (n *AdminDDLStmt) stmtNode() {}

// DDLOption represents a parsed key-value option in a DDL statement.
type DDLOption struct {
	Key   string // option name (e.g. "MAXLOGFILES", "CHARACTER_SET", "USER_SYS")
	Value string // option value (e.g. "16", "AL32UTF8", "password")
	Items *List  // sub-items (e.g. list of DatafileClause for LOGFILE)
	Loc   Loc
}

func (n *DDLOption) nodeTag() {}

// CreateSchemaStmt represents a CREATE SCHEMA statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-SCHEMA.html
//
//	CREATE SCHEMA AUTHORIZATION schema_name
//	  { create_table_statement
//	  | create_view_statement
//	  | grant_statement
//	  } ...
type CreateSchemaStmt struct {
	SchemaName string    // AUTHORIZATION schema_name
	Stmts      *List     // nested CREATE TABLE / CREATE VIEW / GRANT statements
	Loc        Loc
}

func (n *CreateSchemaStmt) nodeTag()  {}
func (n *CreateSchemaStmt) stmtNode() {}

// AlterDatabaseLinkStmt represents an ALTER DATABASE LINK statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-DATABASE-LINK.html
//
//	ALTER [ SHARED ] [ PUBLIC ] DATABASE LINK dblink_name
//	  CONNECT TO user IDENTIFIED BY password
//	  [ AUTHENTICATED BY user IDENTIFIED BY password ]
type AlterDatabaseLinkStmt struct {
	Name              *ObjectName // database link name
	Shared            bool        // SHARED keyword present
	Public            bool        // PUBLIC keyword present
	ConnectUser       string      // CONNECT TO user
	ConnectPassword   string      // IDENTIFIED BY password
	AuthenticatedUser string      // AUTHENTICATED BY user (optional)
	AuthenticatedPass string      // IDENTIFIED BY password (optional, for AUTHENTICATED)
	Loc               Loc
}

func (n *AlterDatabaseLinkStmt) nodeTag()  {}
func (n *AlterDatabaseLinkStmt) stmtNode() {}

// AlterSynonymStmt represents an ALTER SYNONYM statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SYNONYM.html
//
//	ALTER [ PUBLIC ] SYNONYM [ IF EXISTS ] [ schema. ] synonym
//	  { EDITIONABLE | NONEDITIONABLE | COMPILE }
type AlterSynonymStmt struct {
	Name     *ObjectName // synonym name
	Public   bool        // PUBLIC keyword present
	IfExists bool        // IF EXISTS
	Action   string      // "COMPILE", "EDITIONABLE", "NONEDITIONABLE"
	Loc      Loc
}

func (n *AlterSynonymStmt) nodeTag()  {}
func (n *AlterSynonymStmt) stmtNode() {}

// AlterMaterializedViewStmt represents an ALTER MATERIALIZED VIEW statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-MATERIALIZED-VIEW.html
//
//	ALTER MATERIALIZED VIEW [IF EXISTS] [schema.]materialized_view
//	  { alter_mv_refresh
//	  | ENABLE QUERY REWRITE
//	  | DISABLE QUERY REWRITE
//	  | COMPILE
//	  | CONSIDER FRESH
//	  | { ENABLE | DISABLE } CONCURRENT REFRESH
//	  | shrink_clause
//	  | { CACHE | NOCACHE }
//	  | parallel_clause
//	  | logging_clause
//	  | ... }
//
//	alter_mv_refresh:
//	  REFRESH
//	    [ FAST | COMPLETE | FORCE ]
//	    [ ON { COMMIT | DEMAND } ]
//	    [ START WITH date ]
//	    [ NEXT date ]
//	    [ WITH PRIMARY KEY ]
//	    [ USING ROLLBACK SEGMENT rollback_segment ]
//	    [ USING { ENFORCED | TRUSTED } CONSTRAINTS ]
//	    [ { ENABLE | DISABLE } ON QUERY COMPUTATION ]
//
//	shrink_clause:
//	  SHRINK SPACE [ COMPACT | CASCADE ]
type AlterMaterializedViewStmt struct {
	Name     *ObjectName // materialized view name
	IfExists bool        // IF EXISTS

	// Action type
	Action string // "COMPILE", "CONSIDER_FRESH", "REFRESH", "ENABLE_QUERY_REWRITE", "DISABLE_QUERY_REWRITE", "ENABLE_CONCURRENT_REFRESH", "DISABLE_CONCURRENT_REFRESH", "ENABLE_ON_QUERY_COMPUTATION", "DISABLE_ON_QUERY_COMPUTATION", "SHRINK", "CACHE", "NOCACHE", "PARALLEL", "NOPARALLEL", "LOGGING", "NOLOGGING", "ALLOCATE_EXTENT", "DEALLOCATE_UNUSED", "EVALUATE_USING_EDITION", "PCTFREE", "PCTUSED", "INITRANS", "COMPRESS", "NOCOMPRESS", "INMEMORY", "NO_INMEMORY", "MODIFY"

	// REFRESH clause options
	RefreshMethod string   // "FAST", "COMPLETE", "FORCE"
	RefreshMode   string   // "ON_COMMIT", "ON_DEMAND"
	StartWith     ExprNode // START WITH date
	Next          ExprNode // NEXT date
	WithPrimaryKey bool    // WITH PRIMARY KEY
	UsingRollbackSegment string // USING ROLLBACK SEGMENT name
	UsingConstraints     string // "ENFORCED" or "TRUSTED"
	EnableOnQueryComputation  bool // ENABLE ON QUERY COMPUTATION
	DisableOnQueryComputation bool // DISABLE ON QUERY COMPUTATION

	// SHRINK SPACE options
	Compact bool // COMPACT
	Cascade bool // CASCADE

	// PARALLEL options
	ParallelDegree string // parallel degree (number as string)

	// Physical attributes
	PctFree  int // PCTFREE integer
	PctUsed  int // PCTUSED integer
	IniTrans int // INITRANS integer

	// ALLOCATE EXTENT / DEALLOCATE UNUSED
	AllocateSize     string // SIZE value
	AllocateDatafile string // DATAFILE 'filename'
	AllocateInstance int    // INSTANCE integer
	DeallocateKeep   string // KEEP size

	// EVALUATE USING EDITION
	EditionName string // edition name

	// MODIFY SCOPE FOR
	ScopeColumn string      // ref column name
	ScopeTable  *ObjectName // target table name

	// Options (generic list for clauses parsed as options)
	Options *List

	Loc Loc
}

func (n *AlterMaterializedViewStmt) nodeTag()  {}
func (n *AlterMaterializedViewStmt) stmtNode() {}

// SetRoleStmt represents a SET ROLE statement.
//
//	SET ROLE { role [,...] | ALL [EXCEPT role [,...]] | NONE }
type SetRoleStmt struct {
	Roles  []*ObjectName // role names
	All    bool          // ALL
	Except []*ObjectName // EXCEPT role list (when ALL EXCEPT)
	None   bool          // NONE
	Loc    Loc
}

func (n *SetRoleStmt) nodeTag()  {}
func (n *SetRoleStmt) stmtNode() {}

// SetConstraintsStmt represents a SET CONSTRAINT(S) statement.
//
//	SET { CONSTRAINT | CONSTRAINTS } { ALL | constraint [,...] } { IMMEDIATE | DEFERRED }
type SetConstraintsStmt struct {
	All         bool          // ALL constraints
	Constraints []*ObjectName // specific constraint names
	Deferred    bool          // DEFERRED (false = IMMEDIATE)
	Loc         Loc
}

func (n *SetConstraintsStmt) nodeTag()  {}
func (n *SetConstraintsStmt) stmtNode() {}

// CreateAuditPolicyStmt represents a CREATE AUDIT POLICY statement (Unified Auditing).
type CreateAuditPolicyStmt struct {
	Name          string              // policy name
	Privileges    []string            // PRIVILEGES system_privilege [, ...]
	Actions       []*AuditActionEntry // ACTIONS entries
	ComponentActions []*AuditComponentAction // COMPONENT = ... actions
	Roles         []string            // ROLES role [, ...]
	WhenCondition string              // WHEN 'condition'
	EvaluatePer   string              // EVALUATE PER { STATEMENT | SESSION | INSTANCE }
	OnlyToplevel  bool                // ONLY TOPLEVEL
	ContainerAll  *bool               // CONTAINER = { ALL | CURRENT }
	Loc           Loc
}

func (n *CreateAuditPolicyStmt) nodeTag()  {}
func (n *CreateAuditPolicyStmt) stmtNode() {}

// AuditActionEntry represents a single action entry in audit policy ACTIONS clause.
type AuditActionEntry struct {
	Action    string      // action name (SELECT, INSERT, etc.) or ALL
	Columns   []string    // optional column list
	Object    *ObjectName // ON [schema.]object
	Directory string      // ON DIRECTORY dir_name
	MiningModel bool      // ON MINING MODEL
	Loc       Loc
}

func (n *AuditActionEntry) nodeTag() {}

// AuditComponentAction represents a COMPONENT = ... action in audit policy.
type AuditComponentAction struct {
	Component string   // DATAPUMP, DIRECT_LOAD, OLS, XS, DV, SQL_FIREWALL, PROTOCOL
	Actions   []string // component actions or ALL
	Object    string   // optional ON object_name (DV only)
	Loc       Loc
}

func (n *AuditComponentAction) nodeTag() {}

// AlterAuditPolicyStmt represents an ALTER AUDIT POLICY statement (Unified Auditing).
type AlterAuditPolicyStmt struct {
	Name             string              // policy name
	AddDrop          string              // ADD or DROP
	Privileges       []string            // PRIVILEGES privilege [, ...]
	Actions          []*AuditActionEntry // ACTIONS entries
	ComponentActions []*AuditComponentAction // COMPONENT = ... actions
	Roles            []string            // ROLES role [, ...]
	ConditionDrop    bool                // CONDITION DROP
	Condition        string              // CONDITION 'audit_condition'
	EvaluatePer      string              // EVALUATE PER { STATEMENT | SESSION | INSTANCE }
	AddToplevel      bool                // ADD ONLY TOPLEVEL
	DropToplevel     bool                // DROP ONLY TOPLEVEL
	Loc              Loc
}

func (n *AlterAuditPolicyStmt) nodeTag()  {}
func (n *AlterAuditPolicyStmt) stmtNode() {}

// DropAuditPolicyStmt represents a DROP AUDIT POLICY statement (Unified Auditing).
type DropAuditPolicyStmt struct {
	Name string // policy name
	Loc  Loc
}

func (n *DropAuditPolicyStmt) nodeTag()  {}
func (n *DropAuditPolicyStmt) stmtNode() {}

// AuditStmt represents an AUDIT statement (Traditional + Unified Auditing).
//
// Traditional:
//
//	AUDIT { audit_operation_clause [ auditing_by_clause | auditing_on_clause ]
//	      | audit_schema_object_clause }
//	    [ BY { SESSION | ACCESS } ] [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ]
//
// Unified:
//
//	AUDIT POLICY policy_name
//	    [ { BY | EXCEPT } user [, user ]... ]
//	    [ { BY | EXCEPT } USERS WITH ROLE role [, role ]... ]
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	| AUDIT CONTEXT NAMESPACE namespace ATTRIBUTES attribute [, attribute]...
//	    [ BY user [, user ]... ]
type AuditStmt struct {
	// Unified auditing fields
	Policy        string      // POLICY policy_name (unified)
	ByUsers       []string    // BY user [, user]... (unified)
	ExceptUsers   []string    // EXCEPT user [, user]... (unified)
	WithRoles     []string    // BY/EXCEPT USERS WITH ROLE role [, ...]
	WithRoleExcept bool       // true if EXCEPT USERS WITH ROLE (vs BY)
	ContextNS     string      // CONTEXT NAMESPACE namespace (unified)
	ContextAttrs  []string    // ATTRIBUTES attr [, attr]... (unified)
	// Traditional auditing fields
	Actions       []string    // audit actions (traditional)
	Object        *ObjectName // optional object being audited (traditional)
	By            string      // BY SESSION or BY ACCESS (traditional)
	ByUsers2      []string    // auditing_by_clause: BY user [, user]... (traditional)
	OnDefault     bool        // ON DEFAULT (traditional)
	OnDirectory   string      // ON DIRECTORY directory_name (traditional)
	OnNetwork     bool        // NETWORK (traditional)
	OnDirectPath  bool        // DIRECT_PATH LOAD (traditional)
	// Shared
	When          string      // WHENEVER clause text
	ContainerAll  *bool       // CONTAINER = { ALL | CURRENT }
	Loc           Loc
}

func (n *AuditStmt) nodeTag()  {}
func (n *AuditStmt) stmtNode() {}

// NoauditStmt represents a NOAUDIT statement (Traditional + Unified Auditing).
//
// Traditional:
//
//	NOAUDIT { audit_operation_clause [ auditing_by_clause ]
//	        | audit_schema_object_clause }
//	    [ WHENEVER [ NOT ] SUCCESSFUL ] [ CONTAINER = { CURRENT | ALL } ]
//
// Unified:
//
//	NOAUDIT POLICY policy [ BY user [, user]... [ WITH ROLE role [, role]... ] ]
//	| NOAUDIT CONTEXT NAMESPACE namespace ATTRIBUTES attr [, attr]...
//	    [ BY user [, user]... [ WITH ROLE role [, role]... ] ]
type NoauditStmt struct {
	// Unified auditing fields
	Policy       string   // POLICY policy (unified)
	ByUsers      []string // BY user [, user]... (unified)
	WithRoles    []string // WITH ROLE role [, role]... (unified)
	ContextNS    string   // CONTEXT NAMESPACE namespace (unified)
	ContextAttrs []string // ATTRIBUTES attr [, attr]... (unified)
	// Traditional auditing fields
	Actions      []string    // noaudit actions (traditional)
	Object       *ObjectName // optional object (traditional)
	ByUsers2     []string    // auditing_by_clause: BY user [, user]... (traditional)
	OnDefault    bool        // ON DEFAULT (traditional)
	OnDirectory  string      // ON DIRECTORY directory_name (traditional)
	OnNetwork    bool        // NETWORK (traditional)
	OnDirectPath bool        // DIRECT_PATH LOAD (traditional)
	// Shared
	When         string      // WHENEVER clause text
	ContainerAll *bool       // CONTAINER = { CURRENT | ALL }
	Loc          Loc
}

func (n *NoauditStmt) nodeTag()  {}
func (n *NoauditStmt) stmtNode() {}

// AssociateStatisticsStmt represents an ASSOCIATE STATISTICS statement.
//
//	ASSOCIATE STATISTICS WITH { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES }
//	    object [,...] USING [schema.]statistics_type
type AssociateStatisticsStmt struct {
	ObjectType string        // COLUMNS, FUNCTIONS, etc.
	Objects    []*ObjectName // objects to associate
	Using      *ObjectName   // statistics type
	Loc        Loc
}

func (n *AssociateStatisticsStmt) nodeTag()  {}
func (n *AssociateStatisticsStmt) stmtNode() {}

// DisassociateStatisticsStmt represents a DISASSOCIATE STATISTICS statement.
//
//	DISASSOCIATE STATISTICS FROM { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES }
//	    object [,...] [FORCE]
type DisassociateStatisticsStmt struct {
	ObjectType string        // COLUMNS, FUNCTIONS, etc.
	Objects    []*ObjectName // objects to disassociate
	Force      bool          // FORCE
	Loc        Loc
}

func (n *DisassociateStatisticsStmt) nodeTag()  {}
func (n *DisassociateStatisticsStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE TABLESPACE
// ---------------------------------------------------------------------------

// CreateTablespaceStmt represents a CREATE TABLESPACE statement.
//
//	CREATE [ BIGFILE | SMALLFILE ]
//	    { permanent_tablespace_clause
//	    | temporary_tablespace_clause
//	    | undo_tablespace_clause }
type CreateTablespaceStmt struct {
	Name                *ObjectName       // tablespace name
	Bigfile             bool              // BIGFILE
	Smallfile           bool              // SMALLFILE
	Temporary           bool              // TEMPORARY
	Local               bool              // LOCAL (for LOCAL TEMPORARY)
	Undo                bool              // UNDO
	IfNotExists         bool              // IF NOT EXISTS
	Datafiles           []*DatafileClause // DATAFILE / TEMPFILE clauses
	Size                string            // SIZE value (e.g. "100M")
	Autoextend          *AutoextendClause // AUTOEXTEND ON/OFF
	Logging             string            // LOGGING / NOLOGGING / FORCE LOGGING / FILESYSTEM_LIKE_LOGGING
	Online              bool              // ONLINE
	Offline             bool              // OFFLINE
	Extent              string            // EXTENT MANAGEMENT LOCAL/DICTIONARY (AUTOALLOCATE/UNIFORM)
	Segment             string            // SEGMENT SPACE MANAGEMENT AUTO/MANUAL
	Blocksize           string            // BLOCKSIZE value
	Retention           string            // RETENTION GUARANTEE / NOGUARANTEE
	Encryption          string            // ENCRYPTION clause text
	EncryptionAlgorithm string            // USING 'algorithm'
	Compress            string            // DEFAULT COMPRESS / NOCOMPRESS
	MaxSize             string            // MAXSIZE value
	TablespaceGroup     string            // TABLESPACE GROUP name
	Flashback           string            // FLASHBACK ON / OFF
	LostWriteProtection string            // ENABLE/DISABLE/SUSPEND/REMOVE LOST WRITE PROTECTION
	MinimumExtent       string            // MINIMUM EXTENT size
	DefaultParams       string            // DEFAULT tablespace params (summary)
	ForLeaf             string            // FOR ALL / FOR LEAF (temporary tablespace)
	Options             []string          // remaining unparsed option tokens
	Loc                 Loc
}

func (n *CreateTablespaceStmt) nodeTag()  {}
func (n *CreateTablespaceStmt) stmtNode() {}

// AlterTablespaceStmt represents an ALTER TABLESPACE or ALTER TABLESPACE SET statement.
type AlterTablespaceStmt struct {
	Name                *ObjectName       // tablespace name
	IsSet               bool              // true for ALTER TABLESPACE SET
	IfExists            bool              // IF EXISTS
	DefaultParams       string            // DEFAULT params summary
	MinimumExtent       string            // MINIMUM EXTENT size
	Resize              string            // RESIZE size
	Coalesce            bool              // COALESCE
	ShrinkSpace         bool              // SHRINK SPACE
	ShrinkKeep          string            // SHRINK SPACE KEEP size
	ShrinkTempfile      string            // SHRINK TEMPFILE filename/number
	ShrinkTempfileKeep  string            // SHRINK TEMPFILE KEEP size
	RenameTo            string            // RENAME TO new_name
	BeginBackup         bool              // BEGIN BACKUP
	EndBackup           bool              // END BACKUP
	AddDatafile         bool              // ADD DATAFILE
	AddTempfile         bool              // ADD TEMPFILE
	Datafiles           []*DatafileClause // file specifications for ADD
	DropDatafile        bool              // DROP DATAFILE
	DropTempfile        bool              // DROP TEMPFILE
	DropFileRef         string            // filename or file number for DROP
	RenameDatafile      bool              // RENAME DATAFILE
	RenameFrom          []string          // old filenames
	RenameTo2           []string          // new filenames (for RENAME DATAFILE ... TO ...)
	DatafileOnline      bool              // DATAFILE ONLINE
	DatafileOffline     bool              // DATAFILE OFFLINE
	TempfileOnline      bool              // TEMPFILE ONLINE
	TempfileOffline     bool              // TEMPFILE OFFLINE
	Logging             string            // LOGGING / NOLOGGING / FILESYSTEM_LIKE_LOGGING
	ForceLogging        string            // FORCE LOGGING / NO FORCE LOGGING
	Online              bool              // ONLINE
	Offline             bool              // OFFLINE
	OfflineMode         string            // NORMAL / TEMPORARY / IMMEDIATE
	ReadOnly            bool              // READ ONLY
	ReadWrite           bool              // READ WRITE
	Permanent           bool              // PERMANENT
	TempState           bool              // TEMPORARY (state)
	Autoextend          *AutoextendClause // AUTOEXTEND
	Flashback           string            // FLASHBACK ON / OFF
	Retention           string            // RETENTION GUARANTEE / NOGUARANTEE
	Encryption          string            // ENCRYPTION clause summary
	LostWriteProtection string            // ENABLE/DISABLE/SUSPEND/REMOVE LOST WRITE PROTECTION
	TablespaceGroup     string            // TABLESPACE GROUP name
	Loc                 Loc
}

func (n *AlterTablespaceStmt) nodeTag()  {}
func (n *AlterTablespaceStmt) stmtNode() {}

// CreateTablespaceSetStmt represents a CREATE TABLESPACE SET statement.
type CreateTablespaceSetStmt struct {
	Name        *ObjectName       // tablespace set name
	Shardspace  string            // IN SHARDSPACE name
	Datafiles   []*DatafileClause // DATAFILE clauses in USING TEMPLATE
	Logging     string            // LOGGING / NOLOGGING / FILESYSTEM_LIKE_LOGGING
	Encryption  string            // ENCRYPTION clause
	DefaultParams string          // DEFAULT params summary
	Extent      string            // EXTENT MANAGEMENT
	Segment     string            // SEGMENT SPACE MANAGEMENT
	Flashback   string            // FLASHBACK ON / OFF
	Loc         Loc
}

func (n *CreateTablespaceSetStmt) nodeTag()  {}
func (n *CreateTablespaceSetStmt) stmtNode() {}

// DropTablespaceStmt represents a DROP TABLESPACE or DROP TABLESPACE SET statement.
type DropTablespaceStmt struct {
	Name              *ObjectName // tablespace name
	IsSet             bool        // true for DROP TABLESPACE SET
	IfExists          bool        // IF EXISTS
	DropQuota         bool        // DROP QUOTA
	KeepQuota         bool        // KEEP QUOTA
	IncludingContents bool        // INCLUDING CONTENTS
	AndDatafiles      bool        // AND DATAFILES
	KeepDatafiles     bool        // KEEP DATAFILES
	CascadeConstraints bool       // CASCADE CONSTRAINTS
	Loc               Loc
}

func (n *DropTablespaceStmt) nodeTag()  {}
func (n *DropTablespaceStmt) stmtNode() {}

// DatafileClause represents a DATAFILE or TEMPFILE specification.
type DatafileClause struct {
	Filename   string // file path string
	Size       string // SIZE value
	Reuse      bool   // REUSE
	Autoextend *AutoextendClause
	Loc        Loc
}

func (n *DatafileClause) nodeTag() {}

// AutoextendClause represents AUTOEXTEND ON/OFF.
type AutoextendClause struct {
	On      bool   // ON vs OFF
	Next    string // NEXT size
	MaxSize string // MAXSIZE value or UNLIMITED
	Loc     Loc
}

func (n *AutoextendClause) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE CLUSTER
// ---------------------------------------------------------------------------

// CreateClusterStmt represents a CREATE CLUSTER statement.
//
//	CREATE CLUSTER [ schema. ] cluster
//	    ( column datatype [ SORT ] [, ...] )
//	    [ physical_attributes_clause ]
//	    [ SIZE size_clause ]
//	    [ TABLESPACE tablespace ]
//	    [ { INDEX | [ SINGLE TABLE ] HASHKEYS integer [ HASH IS expr ] } ]
//	    [ parallel_clause ]
//	    [ NOROWDEPENDENCIES | ROWDEPENDENCIES ]
//	    [ CACHE | NOCACHE ]
type CreateClusterStmt struct {
	Name        *ObjectName          // cluster name
	Columns     []*ClusterColumn     // cluster key columns
	PctFree     *int                 // PCTFREE value
	PctUsed     *int                 // PCTUSED value
	InitTrans   *int                 // INITRANS value
	Size        string               // SIZE value
	Tablespace  string               // TABLESPACE name
	IsIndex     bool                 // INDEX (explicit)
	IsHash      bool                 // HASHKEYS specified
	HashKeys    string               // HASHKEYS integer value
	SingleTable bool                 // SINGLE TABLE
	HashExpr    ExprNode             // HASH IS expr
	Cache       bool                 // CACHE
	NoCache     bool                 // NOCACHE
	Parallel    string               // PARALLEL / NOPARALLEL / PARALLEL n
	RowDep      bool                 // ROWDEPENDENCIES
	NoRowDep    bool                 // NOROWDEPENDENCIES
	Storage     []string             // STORAGE clause tokens
	Loc         Loc
}

func (n *CreateClusterStmt) nodeTag()  {}
func (n *CreateClusterStmt) stmtNode() {}

// ClusterColumn represents a column in a cluster key.
type ClusterColumn struct {
	Name     string   // column name
	DataType *TypeName // column data type
	Sort     bool      // SORT keyword
	Loc      Loc
}

func (n *ClusterColumn) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE DIMENSION
// ---------------------------------------------------------------------------

// CreateDimensionStmt represents a CREATE DIMENSION statement.
//
//	CREATE DIMENSION [ schema. ] dimension
//	    level_clause ...
//	    { hierarchy_clause | attribute_clause | extended_attribute_clause } ...
type CreateDimensionStmt struct {
	Name        *ObjectName            // dimension name
	Levels      []*DimensionLevel      // LEVEL clauses
	Hierarchies []*DimensionHierarchy  // HIERARCHY clauses
	Attributes  []*DimensionAttribute  // ATTRIBUTE clauses
	Loc         Loc
}

func (n *CreateDimensionStmt) nodeTag()  {}
func (n *CreateDimensionStmt) stmtNode() {}

// DimensionLevel represents a LEVEL clause in CREATE DIMENSION.
//
//	LEVEL level IS ( level_table.level_column [, ...] ) [ SKIP WHEN NULL ]
type DimensionLevel struct {
	Name         string        // level name
	Columns      []*ObjectName // table.column references
	SkipWhenNull bool          // SKIP WHEN NULL
	Loc          Loc
}

func (n *DimensionLevel) nodeTag() {}

// DimensionHierarchy represents a HIERARCHY clause in CREATE DIMENSION.
//
//	HIERARCHY hierarchy_name ( child_level CHILD OF parent_level [CHILD OF ...] )
//	    [ JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level ] ...
type DimensionHierarchy struct {
	Name      string                  // hierarchy name
	Levels    []string                // level names from child to top parent
	JoinKeys  []*DimensionJoinKey     // JOIN KEY clauses
	Loc       Loc
}

func (n *DimensionHierarchy) nodeTag() {}

// DimensionJoinKey represents a JOIN KEY clause in a dimension hierarchy.
//
//	JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level
type DimensionJoinKey struct {
	ChildKeys   []*ObjectName // child key columns
	ParentLevel string        // parent level name
	Loc         Loc
}

func (n *DimensionJoinKey) nodeTag() {}

// DimensionAttribute represents an ATTRIBUTE clause in CREATE DIMENSION.
//
//	ATTRIBUTE level DETERMINES ( dependent_column [, ...] )
//	-- or extended form:
//	ATTRIBUTE attr_name LEVEL level DETERMINES ( dependent_column [, ...] )
type DimensionAttribute struct {
	AttrName    string        // attribute name (may be same as level)
	LevelName   string        // level name (for extended form)
	Columns     []*ObjectName // dependent columns
	Loc         Loc
}

func (n *DimensionAttribute) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER CLUSTER
// ---------------------------------------------------------------------------

// AlterClusterStmt represents an ALTER CLUSTER statement.
//
// BNF: oracle/parser/bnf/ALTER-CLUSTER.bnf
//
//	ALTER CLUSTER [ IF EXISTS ] [ schema . ] cluster
//	    { physical_attributes_clause
//	    | SIZE integer
//	    | MODIFY PARTITION partition allocate_extent_clause
//	    | allocate_extent_clause
//	    | deallocate_unused_clause
//	    | parallel_clause
//	    }
type AlterClusterStmt struct {
	Name     *ObjectName // cluster name
	IfExists bool        // IF EXISTS
	Action   string      // "SIZE", "ALLOCATE_EXTENT", "DEALLOCATE_UNUSED", "MODIFY_PARTITION", "PHYSICAL_ATTRS", "PARALLEL"
	// Physical attributes
	PctFree   *int
	PctUsed   *int
	InitTrans *int
	// SIZE
	Size string
	// MODIFY PARTITION
	ModifyPartition string
	// Parallel
	Parallel string
	Loc      Loc
}

func (n *AlterClusterStmt) nodeTag()  {}
func (n *AlterClusterStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER DIMENSION
// ---------------------------------------------------------------------------

// AlterDimensionStmt represents an ALTER DIMENSION statement.
//
// BNF: oracle/parser/bnf/ALTER-DIMENSION.bnf
//
//	ALTER DIMENSION [ schema . ] dimension
//	    { ADD level_clause
//	    | ADD hierarchy_clause
//	    | ADD attribute_clause
//	    | ADD extended_attribute_clause
//	    | DROP level_clause
//	    | DROP hierarchy_clause
//	    | DROP attribute_clause
//	    | DROP extended_attribute_clause
//	    | COMPILE
//	    }
type AlterDimensionStmt struct {
	Name    *ObjectName // dimension name
	Compile bool        // COMPILE
	// ADD / DROP items
	AddLevels      []*DimensionLevel
	AddHierarchies []*DimensionHierarchy
	AddAttributes  []*DimensionAttribute
	DropLevels     []string // level names to drop
	DropHierarchies []string // hierarchy names to drop
	DropAttributes  []string // attribute names to drop
	Loc            Loc
}

func (n *AlterDimensionStmt) nodeTag()  {}
func (n *AlterDimensionStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED ZONEMAP
// ---------------------------------------------------------------------------

// CreateMaterializedZonemapStmt represents a CREATE MATERIALIZED ZONEMAP statement.
//
// BNF: oracle/parser/bnf/CREATE-MATERIALIZED-ZONEMAP.bnf
//
//	CREATE MATERIALIZED ZONEMAP [ IF NOT EXISTS ]
//	    [ schema. ] zonemap_name
//	    [ zonemap_attributes ]
//	    [ zonemap_refresh_clause ]
//	    [ { ENABLE | DISABLE } PRUNING ]
//	    { create_zonemap_on_table | create_zonemap_as_subquery }
type CreateMaterializedZonemapStmt struct {
	IfNotExists bool        // IF NOT EXISTS
	Name        *ObjectName // zonemap name
	// zonemap_attributes
	Tablespace string
	Scale      *int
	PctFree    *int
	PctUsed    *int
	Cache      bool
	NoCache    bool
	// zonemap_refresh_clause
	RefreshMethod string // FAST, COMPLETE, FORCE
	RefreshOn     string // "ON DEMAND", "ON COMMIT", "ON LOAD", "ON DATA MOVEMENT", "ON LOAD DATA MOVEMENT"
	// ENABLE/DISABLE PRUNING
	EnablePruning  bool
	DisablePruning bool
	// ON table ( columns... )
	OnTable   *ObjectName
	OnColumns []string
	// AS subquery
	ColumnAliases []string
	AsQuery       *SelectStmt
	Loc           Loc
}

func (n *CreateMaterializedZonemapStmt) nodeTag()  {}
func (n *CreateMaterializedZonemapStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED ZONEMAP
// ---------------------------------------------------------------------------

// AlterMaterializedZonemapStmt represents an ALTER MATERIALIZED ZONEMAP statement.
//
// BNF: oracle/parser/bnf/ALTER-MATERIALIZED-ZONEMAP.bnf
//
//	ALTER MATERIALIZED ZONEMAP [ IF EXISTS ] [ schema. ] zonemap_name
//	    { alter_zonemap_attributes
//	    | zonemap_refresh_clause
//	    | { ENABLE | DISABLE } PRUNING
//	    | COMPILE
//	    | REBUILD
//	    | UNUSABLE
//	    } ;
type AlterMaterializedZonemapStmt struct {
	IfExists bool        // IF EXISTS
	Name     *ObjectName // zonemap name
	Action   string      // "ATTRS", "REFRESH", "ENABLE_PRUNING", "DISABLE_PRUNING", "COMPILE", "REBUILD", "UNUSABLE"
	// alter_zonemap_attributes
	PctFree *int
	PctUsed *int
	Cache   bool
	NoCache bool
	// refresh
	RefreshMethod string
	RefreshOn     string
	Loc           Loc
}

func (n *AlterMaterializedZonemapStmt) nodeTag()  {}
func (n *AlterMaterializedZonemapStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// CreateInmemoryJoinGroupStmt represents a CREATE INMEMORY JOIN GROUP statement.
//
// BNF: oracle/parser/bnf/CREATE-INMEMORY-JOIN-GROUP.bnf
//
//	CREATE INMEMORY JOIN GROUP [ IF NOT EXISTS ] [ schema. ] join_group
//	    ( [ schema. ] table ( column )
//	      [, [ schema. ] table ( column ) ]... ) ;
type CreateInmemoryJoinGroupStmt struct {
	IfNotExists bool        // IF NOT EXISTS
	Name        *ObjectName // join group name
	Members     []*JoinGroupMember
	Loc         Loc
}

func (n *CreateInmemoryJoinGroupStmt) nodeTag()  {}
func (n *CreateInmemoryJoinGroupStmt) stmtNode() {}

// JoinGroupMember represents a table(column) in an INMEMORY JOIN GROUP.
type JoinGroupMember struct {
	Table  *ObjectName
	Column string
	Loc    Loc
}

func (n *JoinGroupMember) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// AlterInmemoryJoinGroupStmt represents an ALTER INMEMORY JOIN GROUP statement.
//
// BNF: oracle/parser/bnf/ALTER-INMEMORY-JOIN-GROUP.bnf
//
//	ALTER INMEMORY JOIN GROUP [ IF EXISTS ] [ schema. ] join_group
//	    { ADD | REMOVE } ( [ schema. ] table ( column ) ) ;
type AlterInmemoryJoinGroupStmt struct {
	IfExists bool        // IF EXISTS
	Name     *ObjectName // join group name
	Action   string      // "ADD" or "REMOVE"
	Member   *JoinGroupMember
	Loc      Loc
}

func (n *AlterInmemoryJoinGroupStmt) nodeTag()  {}
func (n *AlterInmemoryJoinGroupStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER INDEX statement
// ---------------------------------------------------------------------------

// AlterIndexStmt represents an ALTER INDEX statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEX.html
//
//	ALTER INDEX [ IF EXISTS ] [ schema. ] index_name
//	    { deallocate_unused_clause
//	    | allocate_extent_clause
//	    | shrink_clause
//	    | parallel_clause
//	    | physical_attributes_clause
//	    | logging_clause
//	    | partial_index_clause
//	    | rebuild_clause
//	    | alter_index_partitioning
//	    | PARAMETERS ( 'odci_parameters' )
//	    | { DEFERRED | IMMEDIATE } INVALIDATION
//	    | COMPILE
//	    | ENABLE
//	    | DISABLE
//	    | { USABLE | UNUSABLE } [ ONLINE ]
//	    | { VISIBLE | INVISIBLE }
//	    | RENAME TO new_index_name
//	    | COALESCE [ CLEANUP [ ONLY ] ] [ parallel_clause ]
//	    | MONITORING USAGE
//	    | NOMONITORING USAGE
//	    | UPDATE BLOCK REFERENCES
//	    | annotations_clause
//	    }
type AlterIndexStmt struct {
	Name     *ObjectName // index name
	IfExists bool        // IF EXISTS
	Action   string      // action keyword: "REBUILD", "RENAME", "COALESCE", etc.
	// REBUILD options
	Partition    string // REBUILD PARTITION name
	Subpartition string // REBUILD SUBPARTITION name
	Tablespace   string // TABLESPACE name
	Online       bool   // ONLINE
	Reverse      bool   // REVERSE
	NoReverse    bool   // NOREVERSE
	// RENAME
	NewName string // RENAME TO new_name (also for RENAME PARTITION/SUBPARTITION)
	// COALESCE
	Cleanup     bool // CLEANUP
	CleanupOnly bool // CLEANUP ONLY
	// SHRINK SPACE
	Compact  bool // COMPACT
	Cascade  bool // CASCADE
	// PARALLEL
	Parallel   string // parallel degree or ""
	NoParallel bool   // NOPARALLEL
	// LOGGING
	Logging   bool // LOGGING
	NoLogging bool // NOLOGGING
	// COMPRESS
	Compress   string // COMPRESS [n]
	NoCompress bool   // NOCOMPRESS
	// INDEXING
	IndexingFull    bool // INDEXING FULL
	IndexingPartial bool // INDEXING PARTIAL
	// Physical attributes
	PctFree   string // PCTFREE integer
	PctUsed   string // PCTUSED integer
	InitTrans string // INITRANS integer
	MaxTrans  string // MAXTRANS integer
	// PARAMETERS
	Parameters string // PARAMETERS ('...')
	// INVALIDATION
	Invalidation string // "DEFERRED" or "IMMEDIATE"
	// DEALLOCATE UNUSED
	DeallocateKeep string // KEEP size
	// ALTER INDEX PARTITIONING sub-actions
	// MODIFY DEFAULT ATTRIBUTES
	ModifyDefaultFor string // FOR PARTITION name (optional)
	// ADD PARTITION
	AddPartitionName string // partition name (optional)
	// MODIFY PARTITION / MODIFY SUBPARTITION sub-action
	ModifyPartAction string // sub-action within MODIFY PARTITION
	// SPLIT PARTITION
	SplitPartition string // old partition name
	SplitValues    *List  // AT (literal, ...)
	Loc            Loc
}

func (n *AlterIndexStmt) nodeTag()  {}
func (n *AlterIndexStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER VIEW statement
// ---------------------------------------------------------------------------

// AlterViewStmt represents an ALTER VIEW statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-VIEW.html
//
//	ALTER VIEW [IF EXISTS] [schema.]view
//	{   COMPILE
//	  | ADD out_of_line_constraint
//	  | MODIFY CONSTRAINT constraint_name { RELY | NORELY }
//	  | DROP CONSTRAINT constraint_name
//	  | { READ ONLY | READ WRITE }
//	  | { EDITIONABLE | NONEDITIONABLE }
//	}
type AlterViewStmt struct {
	Name           *ObjectName      // view name
	IfExists       bool             // IF EXISTS
	Action         string           // "COMPILE", "ADD_CONSTRAINT", "MODIFY_CONSTRAINT", "DROP_CONSTRAINT", "READ_ONLY", "READ_WRITE", "EDITIONABLE", "NONEDITIONABLE", "ANNOTATIONS"
	Constraint     *TableConstraint // for ADD constraint
	ConstraintName string           // for MODIFY/DROP CONSTRAINT
	Rely           bool             // MODIFY CONSTRAINT ... RELY
	NoRely         bool             // MODIFY CONSTRAINT ... NORELY
	Annotations    *List            // for ANNOTATIONS action
	Loc            Loc
}

func (n *AlterViewStmt) nodeTag()  {}
func (n *AlterViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER SEQUENCE statement
// ---------------------------------------------------------------------------

// AlterSequenceStmt represents an ALTER SEQUENCE statement.
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/ALTER-SEQUENCE.html
//
//	ALTER SEQUENCE [IF EXISTS] [schema.]sequence_name
//	  [ INCREMENT BY integer ]
//	  [ MAXVALUE integer | NOMAXVALUE ]
//	  [ MINVALUE integer | NOMINVALUE ]
//	  [ CYCLE | NOCYCLE ]
//	  [ CACHE integer | NOCACHE ]
//	  [ ORDER | NOORDER ]
//	  [ KEEP | NOKEEP ]
//	  [ RESTART [ WITH integer ] ]
//	  [ SCALE [ EXTEND | NOEXTEND ] | NOSCALE ]
//	  [ SHARD [ EXTEND | NOEXTEND ] | NOSHARD ]
//	  [ GLOBAL | SESSION ]
type AlterSequenceStmt struct {
	Name          *ObjectName // sequence name
	IfExists      bool        // IF EXISTS
	IncrementBy   ExprNode    // INCREMENT BY value
	MaxValue      ExprNode    // MAXVALUE value
	MinValue      ExprNode    // MINVALUE value
	NoMaxValue    bool        // NOMAXVALUE
	NoMinValue    bool        // NOMINVALUE
	Cycle         bool        // CYCLE
	NoCycle       bool        // NOCYCLE
	Cache         ExprNode    // CACHE n
	NoCache       bool        // NOCACHE
	Order         bool        // ORDER
	NoOrder       bool        // NOORDER
	Keep          bool        // KEEP
	NoKeep        bool        // NOKEEP
	Restart       bool        // RESTART
	RestartWith   ExprNode    // RESTART WITH value
	Scale         bool        // SCALE
	NoScale       bool        // NOSCALE
	ScaleExtend   bool        // SCALE EXTEND
	ScaleNoExtend bool        // SCALE NOEXTEND
	Shard         bool        // SHARD
	NoShard       bool        // NOSHARD
	ShardExtend   bool        // SHARD EXTEND
	ShardNoExtend bool        // SHARD NOEXTEND
	Global        bool        // GLOBAL
	Session       bool        // SESSION
	Loc           Loc
}

func (n *AlterSequenceStmt) nodeTag()  {}
func (n *AlterSequenceStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER PROCEDURE / FUNCTION / PACKAGE / TRIGGER statements
// ---------------------------------------------------------------------------

// AlterProcedureStmt represents an ALTER PROCEDURE statement.
//
//	ALTER PROCEDURE [IF EXISTS] [schema.]procedure_name
//	    { COMPILE [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS] }
//	    [ EDITIONABLE | NONEDITIONABLE ]
type AlterProcedureStmt struct {
	Name            *ObjectName  // procedure name
	IfExists        bool         // IF EXISTS
	Compile         bool         // COMPILE
	Debug           bool         // COMPILE DEBUG
	ReuseSettings   bool         // REUSE SETTINGS
	CompilerParams  []*SetParam  // compiler_parameters_clause (name=value pairs)
	Editionable     bool         // EDITIONABLE
	NonEditionable  bool         // NONEDITIONABLE
	Loc             Loc
}

func (n *AlterProcedureStmt) nodeTag()  {}
func (n *AlterProcedureStmt) stmtNode() {}

// AlterFunctionStmt represents an ALTER FUNCTION statement.
//
//	ALTER FUNCTION [IF EXISTS] [schema.]function_name
//	    { COMPILE [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS] }
//	    [ EDITIONABLE | NONEDITIONABLE ]
type AlterFunctionStmt struct {
	Name            *ObjectName  // function name
	IfExists        bool         // IF EXISTS
	Compile         bool         // COMPILE
	Debug           bool         // COMPILE DEBUG
	ReuseSettings   bool         // REUSE SETTINGS
	CompilerParams  []*SetParam  // compiler_parameters_clause (name=value pairs)
	Editionable     bool         // EDITIONABLE
	NonEditionable  bool         // NONEDITIONABLE
	Loc             Loc
}

func (n *AlterFunctionStmt) nodeTag()  {}
func (n *AlterFunctionStmt) stmtNode() {}

// AlterPackageStmt represents an ALTER PACKAGE statement.
//
//	ALTER PACKAGE [schema.]package_name
//	    { COMPILE [ PACKAGE | BODY | SPECIFICATION ] [DEBUG]
//	      [compiler_parameters_clause ...] [REUSE SETTINGS] }
//	    [ EDITIONABLE | NONEDITIONABLE ]
type AlterPackageStmt struct {
	Name            *ObjectName  // package name
	IfExists        bool         // IF EXISTS
	Compile         bool         // COMPILE
	CompileTarget   string       // "PACKAGE", "BODY", "SPECIFICATION", or "" (default)
	Debug           bool         // COMPILE DEBUG
	ReuseSettings   bool         // REUSE SETTINGS
	CompilerParams  []*SetParam  // compiler_parameters_clause (name=value pairs)
	Editionable     bool         // EDITIONABLE
	NonEditionable  bool         // NONEDITIONABLE
	Loc             Loc
}

func (n *AlterPackageStmt) nodeTag()  {}
func (n *AlterPackageStmt) stmtNode() {}

// AlterTriggerStmt represents an ALTER TRIGGER statement.
//
//	ALTER TRIGGER [IF EXISTS] [schema.]trigger_name
//	  { ENABLE
//	  | DISABLE
//	  | COMPILE [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS]
//	  | RENAME TO new_name
//	  | EDITIONABLE
//	  | NONEDITIONABLE
//	  }
type AlterTriggerStmt struct {
	Name            *ObjectName  // trigger name
	IfExists        bool         // IF EXISTS
	Action          string       // "ENABLE", "DISABLE", "COMPILE", "RENAME", "EDITIONABLE", "NONEDITIONABLE"
	Debug           bool         // COMPILE DEBUG
	ReuseSettings   bool         // REUSE SETTINGS
	CompilerParams  []*SetParam  // compiler_parameters_clause (name=value pairs)
	NewName         string       // RENAME TO new_name
	Loc             Loc
}

func (n *AlterTriggerStmt) nodeTag()  {}
func (n *AlterTriggerStmt) stmtNode() {}

// AlterTypeStmt represents an ALTER TYPE statement.
//
//	ALTER TYPE [IF EXISTS] [schema.]type_name
//	  { alter_type_clause | type_compile_clause }
//	  [ EDITIONABLE | NONEDITIONABLE ]
//
//	alter_type_clause:
//	    RESET
//	  | [NOT] INSTANTIABLE
//	  | [NOT] FINAL
//	  | ADD ATTRIBUTE ( attribute datatype [, ...] )
//	  | DROP ATTRIBUTE ( attribute [, ...] )
//	  | MODIFY ATTRIBUTE ( attribute datatype [, ...] )
//	  | ADD { MAP | ORDER } MEMBER FUNCTION ...
//	  | ADD { MEMBER | STATIC } { FUNCTION | PROCEDURE } ...
//	  | ADD CONSTRUCTOR FUNCTION ...
//	  | DROP { MAP | ORDER } MEMBER FUNCTION ...
//	  | DROP { MEMBER | STATIC } { FUNCTION | PROCEDURE } ...
//	  | MODIFY LIMIT integer
//	  | MODIFY ELEMENT TYPE datatype
//	  | dependent_handling_clause
//
//	type_compile_clause:
//	    COMPILE [SPECIFICATION | BODY] [DEBUG] [compiler_parameters_clause ...] [REUSE SETTINGS]
//
//	dependent_handling_clause:
//	    INVALIDATE
//	  | CASCADE [INCLUDING TABLE DATA | NOT INCLUDING TABLE DATA | CONVERT TO SUBSTITUTABLE]
//	    [FORCE]
type AlterTypeStmt struct {
	Name            *ObjectName       // type name
	IfExists        bool              // IF EXISTS
	Action          string            // "COMPILE", "ADD_ATTRIBUTE", "DROP_ATTRIBUTE", "MODIFY_ATTRIBUTE",
	                                  // "ADD_METHOD", "DROP_METHOD", "NOT_INSTANTIABLE", "INSTANTIABLE",
	                                  // "NOT_FINAL", "FINAL", "MODIFY_LIMIT", "MODIFY_ELEMENT_TYPE",
	                                  // "RESET", "EDITIONABLE", "NONEDITIONABLE"
	CompileTarget   string            // "SPECIFICATION", "BODY", "" (for COMPILE)
	Debug           bool              // COMPILE DEBUG
	ReuseSettings   bool              // REUSE SETTINGS
	CompilerParams  []*SetParam       // compiler_parameters_clause
	Attributes      []*TypeAttribute  // for ADD/DROP/MODIFY ATTRIBUTE
	MethodSpec      string            // raw method signature text for ADD/DROP method
	MethodKind      string            // "MEMBER", "STATIC", "MAP MEMBER", "ORDER MEMBER", "CONSTRUCTOR"
	MethodType      string            // "FUNCTION" or "PROCEDURE"
	MethodName      string            // method name
	MethodParams    []*Parameter      // method parameters
	MethodReturn    *TypeName         // method return type
	LimitValue      ExprNode          // for MODIFY LIMIT
	ElementType     *TypeName         // for MODIFY ELEMENT TYPE
	Invalidate      bool              // INVALIDATE
	Cascade         bool              // CASCADE
	IncludeData     *bool             // CASCADE INCLUDING TABLE DATA (true) / NOT INCLUDING TABLE DATA (false) / nil (not specified)
	ConvertToSubst  bool              // CASCADE CONVERT TO SUBSTITUTABLE
	Force           bool              // FORCE (exceptions_clause)
	Editionable     bool              // EDITIONABLE
	NonEditionable  bool              // NONEDITIONABLE
	Loc             Loc
}

// TypeAttribute represents an attribute name with optional datatype in ALTER TYPE.
type TypeAttribute struct {
	Name     string    // attribute name
	DataType *TypeName // datatype (nil for DROP ATTRIBUTE)
	Loc      Loc
}

func (n *TypeAttribute) nodeTag() {}

func (n *AlterTypeStmt) nodeTag()  {}
func (n *AlterTypeStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE INDEXTYPE statement
// ---------------------------------------------------------------------------

// CreateIndextypeStmt represents a CREATE INDEXTYPE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEXTYPE.html
//
//	CREATE [ OR REPLACE ] INDEXTYPE [ IF NOT EXISTS ] [ schema. ] indextype
//	    [ SHARING = { METADATA | NONE } ]
//	    FOR [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	    using_type_clause
//	    [ WITH LOCAL [ RANGE ] PARTITION ]
//	    [ storage_table_clause ]
//	    [ array_DML_clause ]
type CreateIndextypeStmt struct {
	OrReplace    bool            // OR REPLACE
	IfNotExists  bool            // IF NOT EXISTS
	Name         *ObjectName     // indextype name
	Sharing      string          // METADATA or NONE
	Operators    []*IndextypeOp  // FOR operator_list
	UsingType    *ObjectName     // USING implementation_type
	WithLocal    bool            // WITH LOCAL PARTITION
	WithRange    bool            // WITH LOCAL RANGE PARTITION
	StorageTable string          // "SYSTEM" or "USER"
	ArrayDML     bool            // WITH ARRAY DML
	Loc          Loc
}

func (n *CreateIndextypeStmt) nodeTag()  {}
func (n *CreateIndextypeStmt) stmtNode() {}

// IndextypeOp represents an operator binding in CREATE/ALTER INDEXTYPE.
//
//	[ schema. ] operator ( parameter_type [, parameter_type ]... )
type IndextypeOp struct {
	Name       *ObjectName // operator name
	ParamTypes []string    // parameter types
	Loc        Loc
}

func (n *IndextypeOp) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER INDEXTYPE statement
// ---------------------------------------------------------------------------

// AlterIndextypeStmt represents an ALTER INDEXTYPE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEXTYPE.html
//
//	ALTER INDEXTYPE [ IF EXISTS ] [ schema. ] indextype
//	    { { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [ , { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	        [ using_type_clause ]
//	    | COMPILE
//	    }
//	    [ WITH LOCAL PARTITION ]
//	    [ storage_table_clause ]
type AlterIndextypeStmt struct {
	IfExists     bool                // IF EXISTS
	Name         *ObjectName         // indextype name
	Action       string              // "ADD_DROP" or "COMPILE"
	Modifications []*IndextypeModOp  // ADD/DROP operators
	UsingType    *ObjectName         // USING implementation_type
	ArrayDML     bool                // WITH ARRAY DML
	WithLocal    bool                // WITH LOCAL PARTITION
	StorageTable string              // "SYSTEM" or "USER"
	Loc          Loc
}

func (n *AlterIndextypeStmt) nodeTag()  {}
func (n *AlterIndextypeStmt) stmtNode() {}

// IndextypeModOp represents an ADD/DROP operator in ALTER INDEXTYPE.
type IndextypeModOp struct {
	Add        bool        // true=ADD, false=DROP
	Name       *ObjectName // operator name
	ParamTypes []string    // parameter types
	Loc        Loc
}

func (n *IndextypeModOp) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE OPERATOR statement
// ---------------------------------------------------------------------------

// CreateOperatorStmt represents a CREATE OPERATOR statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OPERATOR.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] OPERATOR
//	    [ schema. ] operator
//	    binding_clause
//	    [ SHARING = { METADATA | NONE } ]
type CreateOperatorStmt struct {
	OrReplace   bool              // OR REPLACE
	IfNotExists bool              // IF NOT EXISTS
	Name        *ObjectName       // operator name
	Bindings    []*OperatorBinding // BINDING clauses
	Sharing     string            // METADATA or NONE
	Loc         Loc
}

func (n *CreateOperatorStmt) nodeTag()  {}
func (n *CreateOperatorStmt) stmtNode() {}

// OperatorBinding represents a BINDING clause in CREATE/ALTER OPERATOR.
//
//	BINDING ( [ parameter_type [, parameter_type ]... ] )
//	    RETURN return_type
//	    [ implementation_clause ]
//	    using_function_clause
type OperatorBinding struct {
	ParamTypes   []string    // parameter types
	ReturnType   string      // RETURN type
	UsingFunc    *ObjectName // USING function
	AncillaryTo  *ObjectName // ANCILLARY TO primary_operator
	AncillaryParams []string // ANCILLARY TO parameter types
	WithIndexCtx bool        // WITH INDEX CONTEXT
	ScanCtxType  string      // SCAN CONTEXT implementation_type
	ComputeAnc   bool        // COMPUTE ANCILLARY DATA
	WithColumnCtx bool       // WITH COLUMN CONTEXT
	Loc          Loc
}

func (n *OperatorBinding) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER OPERATOR statement
// ---------------------------------------------------------------------------

// AlterOperatorStmt represents an ALTER OPERATOR statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OPERATOR.html
//
//	ALTER OPERATOR [ IF EXISTS ] [ schema. ] operator
//	    { add_binding_clause
//	    | drop_binding_clause
//	    | COMPILE
//	    }
type AlterOperatorStmt struct {
	IfExists  bool              // IF EXISTS
	Name      *ObjectName       // operator name
	Action    string            // "ADD_BINDING", "DROP_BINDING", "COMPILE"
	Binding   *OperatorBinding  // for ADD BINDING
	DropTypes []string          // DROP BINDING parameter types
	DropForce bool              // DROP BINDING ... FORCE
	Loc       Loc
}

func (n *AlterOperatorStmt) nodeTag()  {}
func (n *AlterOperatorStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW LOG
// ---------------------------------------------------------------------------

// CreateMviewLogStmt represents a CREATE MATERIALIZED VIEW LOG statement.
type CreateMviewLogStmt struct {
	OnTable     *ObjectName // ON table
	IfNotExists bool
	Sharing     string   // METADATA, NONE
	WithPK      bool     // WITH PRIMARY KEY
	WithRowID   bool     // WITH ROWID
	WithOID     bool     // WITH OBJECT ID
	WithSeq     bool     // WITH SEQUENCE
	CommitSCN   bool     // COMMIT SCN
	Columns     *List    // column list
	Including   bool     // INCLUDING NEW VALUES
	Excluding   bool     // EXCLUDING NEW VALUES
	PurgeMode   string   // "IMMEDIATE_SYNC", "IMMEDIATE_ASYNC", "START_WITH"
	PurgeStart  ExprNode // START WITH datetime
	PurgeNext   ExprNode // NEXT datetime
	ForRefresh  string   // "FAST", "SYNCHRONOUS"
	StagingLog  string   // staging log name for SYNCHRONOUS REFRESH
	Options     *List    // generic options (physical attrs, TABLESPACE, etc.)
	Loc         Loc
}

func (n *CreateMviewLogStmt) nodeTag()  {}
func (n *CreateMviewLogStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW LOG
// ---------------------------------------------------------------------------

// AlterMviewLogStmt represents an ALTER MATERIALIZED VIEW LOG statement.
type AlterMviewLogStmt struct {
	OnTable  *ObjectName // ON table
	IfExists bool
	Force    bool
	Action   string // "ADD", "ADD_PRIMARY_KEY", "ADD_ROWID", "ADD_OBJECT_ID", "ADD_SEQUENCE", "ADD_COLUMNS", "SHRINK", "MOVE", "PURGE", "FOR_REFRESH", "PARALLEL", "NOPARALLEL", "LOGGING", "NOLOGGING", "ALLOCATE_EXTENT"
	Columns  *List  // columns for ADD
	Options  *List  // generic options
	Loc      Loc
}

func (n *AlterMviewLogStmt) nodeTag()  {}
func (n *AlterMviewLogStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE ANALYTIC VIEW
// ---------------------------------------------------------------------------

// CreateAnalyticViewStmt represents a CREATE ANALYTIC VIEW statement.
type CreateAnalyticViewStmt struct {
	Name             *ObjectName // view name
	OrReplace        bool
	Force            bool
	NoForce          bool
	IfNotExists      bool
	Sharing          string      // METADATA, NONE
	UsingTable       *ObjectName // USING table/view
	UsingAlias       string
	DimBy            *List  // DIMENSION BY list
	Measures         *List  // MEASURES list
	DefaultMeasure   string // DEFAULT MEASURE name
	DefaultAggregate string // DEFAULT AGGREGATE BY function
	Options          *List  // remaining options (cache, fact columns, etc.)
	Loc              Loc
}

func (n *CreateAnalyticViewStmt) nodeTag()  {}
func (n *CreateAnalyticViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER ANALYTIC VIEW
// ---------------------------------------------------------------------------

// AlterAnalyticViewStmt represents an ALTER ANALYTIC VIEW statement.
type AlterAnalyticViewStmt struct {
	Name     *ObjectName
	IfExists bool
	Action   string      // "RENAME", "COMPILE", "ADD_CACHE", "DROP_CACHE"
	NewName  *ObjectName // for RENAME
	Options  *List       // for cache clauses
	Loc      Loc
}

func (n *AlterAnalyticViewStmt) nodeTag()  {}
func (n *AlterAnalyticViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE JSON RELATIONAL DUALITY VIEW
// ---------------------------------------------------------------------------

// CreateJsonDualityViewStmt represents a CREATE JSON RELATIONAL DUALITY VIEW statement.
type CreateJsonDualityViewStmt struct {
	Name                      *ObjectName
	OrReplace                 bool
	IfNotExists               bool
	EnableLogicalReplication  bool
	DisableLogicalReplication bool
	Query                     StmtNode // AS subquery
	Options                   *List    // remaining options
	Loc                       Loc
}

func (n *CreateJsonDualityViewStmt) nodeTag()  {}
func (n *CreateJsonDualityViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// ALTER JSON RELATIONAL DUALITY VIEW
// ---------------------------------------------------------------------------

// AlterJsonDualityViewStmt represents an ALTER JSON RELATIONAL DUALITY VIEW statement.
type AlterJsonDualityViewStmt struct {
	Name   *ObjectName
	Action string // "ENABLE_LOGICAL_REPLICATION", "DISABLE_LOGICAL_REPLICATION"
	Loc    Loc
}

func (n *AlterJsonDualityViewStmt) nodeTag()  {}
func (n *AlterJsonDualityViewStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE ATTRIBUTE DIMENSION
// ---------------------------------------------------------------------------

// CreateAttributeDimensionStmt represents a CREATE ATTRIBUTE DIMENSION statement.
type CreateAttributeDimensionStmt struct {
	Name          *ObjectName
	OrReplace     bool
	Force         bool
	NoForce       bool
	IfNotExists   bool
	Sharing       string // METADATA, NONE
	DimensionType string // STANDARD, TIME
	Classifications *List // classification_clause items
	Sources       *List   // attr_dim_using_clause: source_clause list
	Attributes    *List   // attributes_clause: attr_dim_attribute_clause list
	Levels        *List   // attr_dim_level_clause list
	AllClause     *AttrDimAllClause
	Loc           Loc
}

func (n *CreateAttributeDimensionStmt) nodeTag()  {}
func (n *CreateAttributeDimensionStmt) stmtNode() {}

// AttrDimSourceClause represents a source clause in USING.
type AttrDimSourceClause struct {
	Remote  bool
	Name    *ObjectName
	Alias   string
	IsJoinPath    bool
	JoinPathName  string
	JoinCondition *List // list of *AttrDimJoinCondElem
	Loc     Loc
}

func (n *AttrDimSourceClause) nodeTag() {}

// AttrDimJoinCondElem represents a join condition element: col = col.
type AttrDimJoinCondElem struct {
	LeftTable  string
	LeftCol    string
	RightTable string
	RightCol   string
	Loc        Loc
}

func (n *AttrDimJoinCondElem) nodeTag() {}

// AttrDimAttribute represents an attribute in the ATTRIBUTES clause.
type AttrDimAttribute struct {
	Column string
	Alias  string
	Classifications *List
	Loc    Loc
}

func (n *AttrDimAttribute) nodeTag() {}

// AttrDimLevel represents a LEVEL clause in CREATE ATTRIBUTE DIMENSION.
type AttrDimLevel struct {
	Name            string
	LevelType       string // STANDARD, YEARS, HALF_YEARS, QUARTERS, MONTHS, WEEKS, DAYS, HOURS, MINUTES, SECONDS
	Classifications *List
	KeyAttrs        *List  // key attributes (list of *String)
	KeyNotNull      bool
	KeySkipWhenNull bool
	AltKeyAttrs     *List  // alternate key attributes
	MemberName      ExprNode
	MemberCaption   ExprNode
	MemberDesc      ExprNode
	OrderByAttrs    *List  // list of *AttrDimOrderByItem
	Determines      *List  // list of *String
	Loc             Loc
}

func (n *AttrDimLevel) nodeTag() {}

// AttrDimOrderByItem represents an ORDER BY item in an attribute dimension level.
type AttrDimOrderByItem struct {
	Attribute string
	Direction string // ASC, DESC
	Loc       Loc
}

func (n *AttrDimOrderByItem) nodeTag() {}

// AttrDimAllClause represents the ALL clause at the end of CREATE ATTRIBUTE DIMENSION.
type AttrDimAllClause struct {
	MemberName    ExprNode
	MemberCaption ExprNode
	MemberDesc    ExprNode
	Loc           Loc
}

func (n *AttrDimAllClause) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER ATTRIBUTE DIMENSION
// ---------------------------------------------------------------------------

// AlterAttributeDimensionStmt represents an ALTER ATTRIBUTE DIMENSION statement.
type AlterAttributeDimensionStmt struct {
	Name     *ObjectName
	IfExists bool
	Action   string      // RENAME, COMPILE
	NewName  *ObjectName // for RENAME TO
	Loc      Loc
}

func (n *AlterAttributeDimensionStmt) nodeTag()  {}
func (n *AlterAttributeDimensionStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE HIERARCHY
// ---------------------------------------------------------------------------

// CreateHierarchyStmt represents a CREATE HIERARCHY statement.
type CreateHierarchyStmt struct {
	Name            *ObjectName
	OrReplace       bool
	Force           bool
	NoForce         bool
	IfNotExists     bool
	Sharing         string // METADATA, NONE
	Classifications *List
	UsingAttrDim    *ObjectName // USING attr_dimension
	LevelHier       *HierLevelClause
	HierAttrs       *List // HIERARCHICAL ATTRIBUTES list
	Loc             Loc
}

func (n *CreateHierarchyStmt) nodeTag()  {}
func (n *CreateHierarchyStmt) stmtNode() {}

// HierLevelClause represents a level in a hierarchy definition.
type HierLevelClause struct {
	Name            string
	Classifications *List
	ChildOf         *HierLevelClause // recursive
	Loc             Loc
}

func (n *HierLevelClause) nodeTag() {}

// HierAttr represents a hierarchical attribute.
type HierAttr struct {
	Name            string // HIER_ORDER, DEPTH, IS_LEAF, etc.
	Classifications *List
	Loc             Loc
}

func (n *HierAttr) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER HIERARCHY
// ---------------------------------------------------------------------------

// AlterHierarchyStmt represents an ALTER HIERARCHY statement.
type AlterHierarchyStmt struct {
	Name     *ObjectName
	IfExists bool
	Action   string      // RENAME, COMPILE
	NewName  *ObjectName // for RENAME TO
	Loc      Loc
}

func (n *AlterHierarchyStmt) nodeTag()  {}
func (n *AlterHierarchyStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE DOMAIN
// ---------------------------------------------------------------------------

// CreateDomainStmt represents a CREATE DOMAIN statement.
type CreateDomainStmt struct {
	Name        *ObjectName
	OrReplace   bool
	Usecase     bool
	IfNotExists bool
	DomainType  string // SINGLE, MULTI, FLEXIBLE, ENUM
	// Single column domain fields
	DataType    *TypeName
	Strict      bool
	Default     ExprNode
	DefaultOnNull bool
	Constraints *List  // constraint_clause list
	Collation   string
	Display     ExprNode
	Order       ExprNode
	Annotations *List
	// Enum fields
	EnumItems   *List // list of *DomainEnumItem
	// Multi column domain fields
	Columns     *List // list of *DomainColumn
	// Flexible domain fields
	FlexDomainName *ObjectName
	FlexColumns    *List // column names
	ChooseUsing    *List // discriminant columns (list of *DomainColumn)
	ChooseExpr     ExprNode // DECODE or CASE expression
	Loc         Loc
}

func (n *CreateDomainStmt) nodeTag()  {}
func (n *CreateDomainStmt) stmtNode() {}

// DomainColumn represents a column in a multi-column domain.
type DomainColumn struct {
	Name        string
	DataType    *TypeName
	Annotations *List
	Loc         Loc
}

func (n *DomainColumn) nodeTag() {}

// DomainEnumItem represents an item in an ENUM domain.
type DomainEnumItem struct {
	Name    string
	Aliases []string
	Value   ExprNode
	Loc     Loc
}

func (n *DomainEnumItem) nodeTag() {}

// DomainConstraint represents a constraint in a domain definition.
type DomainConstraint struct {
	Name       string // constraint name (optional)
	Type       string // NOT_NULL, NULL, CHECK
	CheckExpr  ExprNode
	State      *List // constraint state options
	Loc        Loc
}

func (n *DomainConstraint) nodeTag() {}

// ---------------------------------------------------------------------------
// ALTER DOMAIN
// ---------------------------------------------------------------------------

// AlterDomainStmt represents an ALTER DOMAIN statement.
type AlterDomainStmt struct {
	Name        *ObjectName
	Usecase     bool
	IfExists    bool
	Action      string   // ADD_DISPLAY, MODIFY_DISPLAY, DROP_DISPLAY, ADD_ORDER, MODIFY_ORDER, DROP_ORDER, ANNOTATIONS
	Display     ExprNode // for ADD/MODIFY DISPLAY
	Order       ExprNode // for ADD/MODIFY ORDER
	Annotations *List    // for annotations_clause
	Loc         Loc
}

func (n *AlterDomainStmt) nodeTag()  {}
func (n *AlterDomainStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE PROPERTY GRAPH
// ---------------------------------------------------------------------------

// CreatePropertyGraphStmt represents a CREATE PROPERTY GRAPH statement.
type CreatePropertyGraphStmt struct {
	Name         *ObjectName
	OrReplace    bool
	IfNotExists  bool
	VertexTables []*GraphTableDef  // VERTEX TABLES (...)
	EdgeTables   []*GraphEdgeDef   // EDGE TABLES (...)
	Options      *GraphOptions     // OPTIONS (...)
	Loc          Loc
}

func (n *CreatePropertyGraphStmt) nodeTag()  {}
func (n *CreatePropertyGraphStmt) stmtNode() {}

// GraphTableDef represents a vertex table definition in CREATE PROPERTY GRAPH.
type GraphTableDef struct {
	Name         *ObjectName // [schema.]table_name
	Alias        string      // AS graph_element_name
	KeyColumns   []string    // KEY (col, ...)
	Labels       []string    // LABEL label_name ...
	Properties   string      // ALL, ALL_EXCEPT, LIST, NO
	PropColumns  []string    // column names for PROPERTIES (...)
	PropAliases  []string    // AS property_name aliases
	DefaultLabel bool        // DEFAULT LABEL
	Loc          Loc
}

func (n *GraphTableDef) nodeTag() {}

// GraphEdgeDef represents an edge table definition in CREATE PROPERTY GRAPH.
type GraphEdgeDef struct {
	Name              *ObjectName // [schema.]table_name
	Alias             string      // AS graph_element_name
	KeyColumns        []string    // KEY (col, ...)
	SourceKeyColumns  []string    // SOURCE KEY (col, ...)
	SourceRef         string      // REFERENCES vertex_table_name
	SourceRefColumns  []string    // vertex_table_name (col, ...)
	DestKeyColumns    []string    // DESTINATION KEY (col, ...)
	DestRef           string      // REFERENCES vertex_table_name
	DestRefColumns    []string    // vertex_table_name (col, ...)
	Labels            []string    // LABEL label_name ...
	Properties        string      // ALL, ALL_EXCEPT, LIST, NO
	PropColumns       []string    // column names for PROPERTIES (...)
	PropAliases       []string    // AS property_name aliases
	DefaultLabel      bool        // DEFAULT LABEL
	Loc               Loc
}

func (n *GraphEdgeDef) nodeTag() {}

// GraphOptions represents OPTIONS clause in CREATE PROPERTY GRAPH.
type GraphOptions struct {
	Mode           string // ENFORCED, TRUSTED
	MixedPropTypes string // ALLOW, DISALLOW
	Loc            Loc
}

func (n *GraphOptions) nodeTag() {}

// ---------------------------------------------------------------------------
// CREATE VECTOR INDEX
// ---------------------------------------------------------------------------

// CreateVectorIndexStmt represents a CREATE VECTOR INDEX statement.
type CreateVectorIndexStmt struct {
	Name            *ObjectName // index name
	IfNotExists     bool
	TableName       *ObjectName // ON table
	Column          string      // (column_name)
	IncludeColumns  []string    // INCLUDE (col, ...)
	Organization    string      // INMEMORY_NEIGHBOR_GRAPH, NEIGHBOR_PARTITIONS
	Distance        string      // metric_name
	TargetAccuracy  int         // WITH TARGET ACCURACY integer
	ParameterType   string      // HNSW, IVF
	Neighbors       int         // NEIGHBORS / M value
	EfConstruction  int         // EFCONSTRUCTION value
	RescoreFactor   int         // RESCORE_FACTOR value
	Quantization    string      // quantization algorithm
	NeighborParts   int         // NEIGHBOR PARTITIONS value (IVF)
	SamplesPerPart  int         // SAMPLES_PER_PARTITION (IVF)
	MinVecsPerPart  int         // MIN_VECTORS_PER_PARTITION (IVF)
	Replication     string      // DUPLICATE_ALL, DISTRIBUTE, DISTRIBUTE_AUTO, etc.
	Online          bool
	Parallel        int
	Loc             Loc
}

func (n *CreateVectorIndexStmt) nodeTag()  {}
func (n *CreateVectorIndexStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE / ALTER LOCKDOWN PROFILE
// ---------------------------------------------------------------------------

// CreateLockdownProfileStmt represents a CREATE LOCKDOWN PROFILE statement.
type CreateLockdownProfileStmt struct {
	Name        string // profile_name
	Using       string // USING base_profile_name (static)
	Including   string // INCLUDING base_profile_name (dynamic)
	Loc         Loc
}

func (n *CreateLockdownProfileStmt) nodeTag()  {}
func (n *CreateLockdownProfileStmt) stmtNode() {}

// AlterLockdownProfileStmt represents an ALTER LOCKDOWN PROFILE statement.
type AlterLockdownProfileStmt struct {
	Name       string // profile_name
	Action     string // ENABLE, DISABLE
	RuleType   string // FEATURE, OPTION, STATEMENT
	RuleItems  []string // feature/option/statement names
	AllItems   bool     // ALL
	ExceptItems []string // ALL EXCEPT (...)
	Clauses    []string // CLAUSE = (...)
	ClauseAll  bool     // CLAUSE = ALL
	ClauseExceptItems []string // CLAUSE ALL EXCEPT (...)
	ClauseOptions []string // OPTION = (...)
	ValueItems []string // VALUE = (...)
	MinValue   string   // MINVALUE =
	MaxValue   string   // MAXVALUE =
	Users      string   // ALL, COMMON, LOCAL
	Loc        Loc
}

func (n *AlterLockdownProfileStmt) nodeTag()  {}
func (n *AlterLockdownProfileStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// CREATE / ALTER OUTLINE
// ---------------------------------------------------------------------------

// CreateOutlineStmt represents a CREATE OUTLINE statement.
type CreateOutlineStmt struct {
	Name       string // outline name (optional)
	OrReplace  bool
	Public     bool   // PUBLIC
	Private    bool   // PRIVATE
	Category   string // FOR CATEGORY category
	FromSource string // FROM [PRIVATE] source_outline
	FromPrivate bool  // FROM PRIVATE
	Loc        Loc
}

func (n *CreateOutlineStmt) nodeTag()  {}
func (n *CreateOutlineStmt) stmtNode() {}

// AlterOutlineStmt represents an ALTER OUTLINE statement.
type AlterOutlineStmt struct {
	Name      string // outline name
	Public    bool   // PUBLIC
	Private   bool   // PRIVATE
	Action    string // REBUILD, RENAME, CHANGE_CATEGORY, ENABLE, DISABLE
	NewName   string // RENAME TO new_name
	Category  string // CHANGE CATEGORY TO new_category
	Loc       Loc
}

func (n *AlterOutlineStmt) nodeTag()  {}
func (n *AlterOutlineStmt) stmtNode() {}

// ---------------------------------------------------------------------------
// Star expression (SELECT *)
// ---------------------------------------------------------------------------

// Star represents a * in a select list.
type Star struct {
	Loc Loc
}

func (n *Star) nodeTag()  {}
func (n *Star) exprNode() {}
