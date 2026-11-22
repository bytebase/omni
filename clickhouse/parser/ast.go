package parser

import "fmt"

// ParseError represents a ClickHouse lexical or syntax error with a byte offset.
type ParseError struct {
	Message  string
	Position int // 0-based byte offset in the input SQL
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("clickhouse syntax error at byte %d: %s", e.Position, e.Message)
}

// Node is the base interface implemented by all ClickHouse AST nodes.
type Node interface {
	StartPos() int
	EndPos() int
}

// Statement is implemented by all top-level ClickHouse SQL statements.
type Statement interface {
	Node
	stmtNode()
}

// Expr is implemented by all ClickHouse SQL expressions.
type Expr interface {
	Node
	exprNode()
}

type span struct {
	Pos int
	End int
}

func (s span) StartPos() int { return s.Pos }
func (s span) EndPos() int   { return s.End }

// OutputClause captures optional INTO OUTFILE ... and FORMAT ... clauses
// attached to top-level ClickHouse queries.
type OutputClause struct {
	span
	Outfile     string
	Compression string
	Format      string
}

// --- Statements ---

type SelectUnionStmt struct {
	span
	Items  []*SelectStmt
	Ops    []string // "UNION ALL", "UNION DISTINCT", "UNION", "INTERSECT ...", "EXCEPT ..."
	Output *OutputClause
}

func (*SelectUnionStmt) stmtNode() {}

type SelectStmt struct {
	span
	With        []*WithExpr
	Distinct    bool
	DistinctOn  []Expr
	Top         *TopClause
	Columns     []Expr
	From        *JoinExpr
	ArrayJoin   *ArrayJoinClause
	Windows     []*WindowDef
	Prewhere    Expr
	Where       Expr
	GroupBy     *GroupByClause
	Having      Expr
	Qualify     Expr
	OrderBy     []*OrderByItem
	Interpolate []Expr
	LimitBy     *LimitByClause
	Limit       *LimitClause
	Settings    []*SettingExpr
	// Nested holds a parenthesized SelectUnionStmt when this item is `(SELECT ...)`.
	Nested *SelectUnionStmt
}

func (*SelectStmt) stmtNode() {}

type WithExpr struct {
	span
	Recursive bool
	Alias     string
	Expr      Expr
	Subquery  *SelectUnionStmt
}

type TopClause struct {
	span
	Count    string
	WithTies bool
}

type ArrayJoinClause struct {
	span
	Left  bool
	Inner bool
	Exprs []Expr
}

type WindowDef struct {
	span
	Name string
	Spec *WindowSpec
}

type WindowSpec struct {
	span
	BaseWindow  string // optional inherited window name in OVER (base_win ...)
	PartitionBy []Expr
	OrderBy     []*OrderByItem
	FrameUnits  string // "ROWS" or "RANGE"
	FrameStart  string
	FrameEnd    string
}

type GroupByClause struct {
	span
	All        bool
	Modifier   string // "", "CUBE", "ROLLUP", "GROUPING SETS"
	Exprs      []Expr
	WithCube   bool
	WithRollup bool
	WithTotals bool
}

type OrderByItem struct {
	span
	Expr     Expr
	Desc     bool
	NullsPos string // "", "FIRST", "LAST"
	Collate  string
	WithFill bool
	FillFrom Expr
	FillTo   Expr
	FillStep Expr
}

type LimitByClause struct {
	span
	Limit  Expr
	Offset Expr
	By     []Expr
}

type LimitClause struct {
	span
	Limit    Expr
	Offset   Expr
	WithTies bool
}

type SettingExpr struct {
	span
	Name  string
	Value Expr
}

type JoinExpr struct {
	span
	Left       *JoinExpr
	Right      *JoinExpr
	Table      *TableExpr
	Global     bool
	Local      bool
	Kind       string // "INNER", "LEFT", "RIGHT", "FULL", "CROSS", "COMMA"
	Strictness string // "", "ALL", "ANY", "ASOF", "SEMI", "ANTI"
	On         []Expr
	Using      []Expr
	Final      bool
	Sample     string
	SampleOff  string
}

type TableExpr struct {
	span
	Database string
	Table    string
	Func     *FuncCallExpr
	Subquery *SelectUnionStmt
	Alias    string
}

// DML Statements

type InsertStmt struct {
	span
	Database string
	Table    string
	Func     *FuncCallExpr
	Columns  []string
	Format   string
	Values   [][]Expr
	Data     string // raw inline payload after FORMAT <name> or opaque VALUES data
	Select   *SelectUnionStmt
	Settings []*SettingExpr
}

func (*InsertStmt) stmtNode() {}

type DeleteStmt struct {
	span
	Table     string
	Cluster   string
	Partition Expr
	Where     Expr
	Settings  []*SettingExpr
}

func (*DeleteStmt) stmtNode() {}

type UpdateStmt struct {
	span
	Table       string
	Assignments []*AssignmentExpr
	Cluster     string
	Partition   Expr
	Where       Expr
	Settings    []*SettingExpr
}

func (*UpdateStmt) stmtNode() {}

type AssignmentExpr struct {
	span
	Column string
	Value  Expr
}

// DDL Statements

type CreateDatabaseStmt struct {
	span
	Attach      bool
	IfNotExists bool
	Name        string
	Cluster     string
	Engine      *EngineExpr
	Comment     string
	Output      *OutputClause
}

func (*CreateDatabaseStmt) stmtNode() {}

type CreateTableStmt struct {
	span
	Action      string // "CREATE", "ATTACH", "REPLACE"
	OrReplace   bool
	Temporary   bool
	IfNotExists bool
	Database    string
	Name        string
	UUID        string
	Cluster     string
	Columns     []*TableElement
	AsTable     string
	AsFunc      *FuncCallExpr
	Engine      *EngineClause
	Subquery    *SelectUnionStmt
	Comment     string
	Output      *OutputClause
}

func (*CreateTableStmt) stmtNode() {}

type CreateViewStmt struct {
	span
	Action       string // "CREATE", "ATTACH"
	OrReplace    bool
	Materialized bool
	IfNotExists  bool
	Database     string
	Name         string
	UUID         string
	Cluster      string
	Columns      []*TableElement
	ToTable      string
	Engine       *EngineClause
	Populate     bool
	Subquery     *SelectUnionStmt
	Output       *OutputClause
}

func (*CreateViewStmt) stmtNode() {}

type CreateDictionaryStmt struct {
	span
	Action      string
	OrReplace   bool
	IfNotExists bool
	Database    string
	Name        string
	UUID        string
	Cluster     string
	Attributes  []*DictionaryAttr
	PrimaryKey  []Expr
	Source      *DictionaryConfig
	Layout      *DictionaryConfig
	Lifetime    *DictionaryLifetime
	Range       *DictionaryRange
	Settings    []*SettingExpr
	Output      *OutputClause
}

func (*CreateDictionaryStmt) stmtNode() {}

type DictionaryAttr struct {
	span
	Name         string
	Type         *DataType
	Default      *LiteralExpr
	Expression   Expr
	Hierarchical bool
	Injective    bool
	IsObjectID   bool
}

type DictionaryConfig struct {
	span
	Name   string
	Params []*SettingExpr
}

type DictionaryLifetime struct {
	span
	Min string
	Max string
}

type DictionaryRange struct {
	span
	Min string
	Max string
}

type CreateFunctionStmt struct {
	span
	OrReplace   bool
	IfNotExists bool
	Name        string
	Cluster     string
	Body        Expr
	Output      *OutputClause
}

func (*CreateFunctionStmt) stmtNode() {}

type EngineExpr struct {
	span
	Name string
	Args []Expr
}

type EngineClause struct {
	span
	Engine      *EngineExpr
	OrderBy     []*OrderByItem
	PartitionBy Expr
	PrimaryKey  Expr
	SampleBy    Expr
	TTL         []*TTLExpr
	Settings    []*SettingExpr
}

type TTLExpr struct {
	span
	Expr        Expr              // expiration expression
	Action      string            // "", "DELETE", "TO DISK", "TO VOLUME", "GROUP BY"
	Where       Expr              // optional DELETE WHERE predicate
	Destination string            // target disk or volume name for TO DISK / TO VOLUME
	GroupBy     *GroupByClause    // GROUP BY clause for rollup TTL
	Set         []*AssignmentExpr // SET assignments for GROUP BY rollup TTL
}

type CodecExpr struct {
	span
	Codecs []*CodecArg
}

type CodecArg struct {
	span
	Name string
	Args []Expr
}

type TableElement struct {
	span
	Kind        string // "COLUMN", "CONSTRAINT", "INDEX", "PROJECTION"
	Name        string
	Type        *DataType
	DefaultKind string // "", "DEFAULT", "MATERIALIZED", "ALIAS"
	Expr        Expr
	Codec       *CodecExpr  // column CODEC(...)
	TTLExpr     Expr        // column TTL expression
	Granularity string      // for INDEX definitions
	Projection  *SelectStmt // for PROJECTION definitions
	Comment     string
}

type AlterTableStmt struct {
	span
	Database string
	Table    string
	Cluster  string
	Clauses  []*AlterClause
	Output   *OutputClause
}

func (*AlterTableStmt) stmtNode() {}

type AlterClause struct {
	span
	Action      string
	Name        string
	NewName     string            // for RENAME COLUMN
	Element     *TableElement     // column, index, projection, or constraint definition
	After       string            // for ADD ... AFTER col
	First       bool              // for ADD ... FIRST
	Assignments []*AssignmentExpr // for UPDATE ... WHERE ...
	Where       Expr              // for DELETE WHERE ... and UPDATE ... WHERE ...
	Partition   Expr              // for partition clauses
	Target      string            // for MOVE PARTITION TO ..., ATTACH/REPLACE ... FROM ...
	Comment     string            // for COMMENT COLUMN / MODIFY COLUMN COMMENT
	Codec       *CodecExpr        // for MODIFY COLUMN ... CODEC(...)
	RemoveProp  string            // for MODIFY COLUMN REMOVE <prop>
	OrderBy     Expr              // for MODIFY ORDER BY
	TTL         []*TTLExpr        // for MODIFY TTL
	Settings    []*SettingExpr    // for MODIFY SETTING
}

type DropStmt struct {
	span
	Action      string // "DROP", "DETACH"
	ObjectType  string // "DATABASE", "TABLE", "TEMPORARY TABLE", "VIEW", "DICTIONARY", "FUNCTION", ...
	IfExists    bool
	Database    string
	Name        string
	Cluster     string
	NoDelay     bool
	Permanently bool
	Sync        bool
	Output      *OutputClause
}

func (*DropStmt) stmtNode() {}

type AttachStmt struct {
	span
	ObjectType string
	Database   string
	Name       string
	Cluster    string
	Output     *OutputClause
}

func (*AttachStmt) stmtNode() {}

type RenameStmt struct {
	span
	ObjectType string // "TABLE", "DATABASE", "DICTIONARY"
	Pairs      [][2]string
	Cluster    string
	Output     *OutputClause
}

func (*RenameStmt) stmtNode() {}

type ExchangeStmt struct {
	span
	ObjectType string // "TABLES", "DICTIONARIES"
	First      string
	Second     string
	Cluster    string
	Output     *OutputClause
}

func (*ExchangeStmt) stmtNode() {}

type TruncateStmt struct {
	span
	Temporary bool
	IfExists  bool
	Database  string
	Table     string
	Cluster   string
	Output    *OutputClause
}

func (*TruncateStmt) stmtNode() {}

type UndropStmt struct {
	span
	Database string
	Table    string
	UUID     string
	Cluster  string
	Output   *OutputClause
}

func (*UndropStmt) stmtNode() {}

// Admin / Utility Statements

type ShowStmt struct {
	span
	Kind         string
	Target       string
	Targets      []string // comma-separated targets (e.g. SHOW CREATE USER u1, u2; SHOW GRANTS FOR u1, u2)
	From         string   // FROM/IN database or table
	Full         bool
	Extended     bool
	Temporary    bool
	NotLike      bool
	Like         string
	ILike        bool
	Where        Expr
	Limit        Expr
	WithImplicit bool
	Final        bool
	Output       *OutputClause
}

func (*ShowStmt) stmtNode() {}

type DescribeStmt struct {
	span
	Target *TableExpr
	Output *OutputClause
}

func (*DescribeStmt) stmtNode() {}

type ExplainStmt struct {
	span
	Kind     string // "", "AST", "SYNTAX", "QUERY TREE", "PLAN", "PIPELINE", "ESTIMATE", "TABLE OVERRIDE"
	Settings []*SettingExpr
	Table    string // target table for EXPLAIN TABLE OVERRIDE
	Inner    Statement
	Output   *OutputClause
}

func (*ExplainStmt) stmtNode() {}

type ExistsStmt struct {
	span
	ObjectType string
	Database   string
	Name       string
	Output     *OutputClause
}

func (*ExistsStmt) stmtNode() {}

type CheckStmt struct {
	span
	Database  string
	Table     string
	Partition Expr
	Output    *OutputClause
}

func (*CheckStmt) stmtNode() {}

type OptimizeStmt struct {
	span
	Database      string
	Table         string
	Cluster       string
	Partition     Expr
	Final         bool
	Deduplicate   bool
	DeduplicateBy []Expr
	Output        *OutputClause
}

func (*OptimizeStmt) stmtNode() {}

type KillStmt struct {
	span
	Kind    string // "MUTATION", "QUERY"
	Cluster string
	Where   Expr
	Mode    string // "", "SYNC", "ASYNC", "TEST"
	Output  *OutputClause
}

func (*KillStmt) stmtNode() {}

type SystemStmt struct {
	span
	Action string
	Target string
	Output *OutputClause
}

func (*SystemStmt) stmtNode() {}

type SetStmt struct {
	span
	Settings []*SettingExpr
	Output   *OutputClause
}

func (*SetStmt) stmtNode() {}

type UseStmt struct {
	span
	Database string
	Output   *OutputClause
}

func (*UseStmt) stmtNode() {}

type WatchStmt struct {
	span
	Database string
	Table    string
	Events   bool
	Limit    string
	Output   *OutputClause
}

func (*WatchStmt) stmtNode() {}

type PrivilegeSpec struct {
	span
	Name    string
	Columns []string
}

type AccessControlStmt struct {
	span
	Kind           string           // "GRANT", "REVOKE", "CREATE USER", "ALTER USER", "CREATE ROLE", ...
	OrReplace      bool             // CREATE OR REPLACE
	IfNotExists    bool             // CREATE ... IF NOT EXISTS
	IfExists       bool             // ALTER ... IF EXISTS
	Cluster        string           // optional ON CLUSTER
	Names          []string         // entity names for CREATE/ALTER USER/ROLE/QUOTA/PROFILE/POLICY
	Tables         []string         // ON table list for ROW POLICY
	RevokeOption   string           // "", "ADMIN OPTION", "GRANT OPTION" for REVOKE ... OPTION FOR
	Privileges     []*PrivilegeSpec // privileges or roles for GRANT/REVOKE
	Target         string           // ON target (*.*, db.*, db.tbl, tbl) for GRANT/REVOKE
	Grantees       []string         // TO / FROM grantees
	ExceptGrantees []string         // EXCEPT grantees
	Options        []string         // WITH GRANT OPTION / WITH ADMIN OPTION / WITH REPLACE OPTION
	IdentifiedWith string           // USER IDENTIFIED WITH <type>
	IdentifiedBy   string           // USER IDENTIFIED [WITH ...] BY '...'
	DefaultRoles   []string         // USER DEFAULT ROLE ...
	Hosts          []string         // USER HOST ...
	ForSelect      bool             // ROW POLICY FOR SELECT
	Using          Expr             // ROW POLICY USING expr
	Settings       []*SettingExpr   // SETTINGS ...
	Output         *OutputClause
}

func (*AccessControlStmt) stmtNode() {}

type BackupElement struct {
	span
	Kind       string // "ALL", "TABLE", "TEMPORARY TABLE", "DATABASE", "DICTIONARY", "VIEW"
	Name       string
	Alias      string
	ExceptKind string // "", "TABLES", "DATABASES", "TABLE", "DATABASE"
	Except     []string
}

type BackupRestoreStmt struct {
	span
	Kind            string           // "BACKUP", "RESTORE"
	Items           []*BackupElement // elements being backed up or restored
	Cluster         string           // optional ON CLUSTER
	Destination     string           // "Null" or identifier destination
	DestinationFunc *FuncCallExpr    // Disk(...), S3(...), File(...)
	Settings        []*SettingExpr
	Async           bool
	Output          *OutputClause
}

func (*BackupRestoreStmt) stmtNode() {}

// --- Expressions ---

type IdentExpr struct {
	span
	Parts []string
}

func (*IdentExpr) exprNode() {}

type LiteralExpr struct {
	span
	Kind string // "INT", "FLOAT", "STRING", "HEX_STRING", "BIN_STRING", "BOOL", "NULL", "INF", "NAN"
	Raw  string
	Val  string
}

func (*LiteralExpr) exprNode() {}

type UnaryExpr struct {
	span
	Op   string
	Expr Expr
}

func (*UnaryExpr) exprNode() {}

type BinaryExpr struct {
	span
	Op    string
	Left  Expr
	Right Expr
}

func (*BinaryExpr) exprNode() {}

type TernaryExpr struct {
	span
	Cond Expr
	Then Expr
	Else Expr
}

func (*TernaryExpr) exprNode() {}

type AliasExpr struct {
	span
	Expr  Expr
	Alias string
}

func (*AliasExpr) exprNode() {}

type FuncCallExpr struct {
	span
	Name     string
	Params   []Expr // parametric aggregate params: f(p1)(a1)
	Distinct bool
	Args     []Expr
	Over     *WindowSpec
	OverName string
}

func (*FuncCallExpr) exprNode() {}

type CastExpr struct {
	span
	Expr Expr
	Type *DataType
}

func (*CastExpr) exprNode() {}

type CaseExpr struct {
	span
	Operand Expr
	Whens   []CaseWhen
	Else    Expr
}

type CaseWhen struct {
	Cond Expr
	Then Expr
}

func (*CaseExpr) exprNode() {}

type BetweenExpr struct {
	span
	Not  bool
	Expr Expr
	Low  Expr
	High Expr
}

func (*BetweenExpr) exprNode() {}

type InExpr struct {
	span
	Global bool
	Not    bool
	Left   Expr
	Right  Expr
}

func (*InExpr) exprNode() {}

type LikeExpr struct {
	span
	Not   bool
	Ilike bool
	Left  Expr
	Right Expr
}

func (*LikeExpr) exprNode() {}

type IsNullExpr struct {
	span
	Not  bool
	Expr Expr
}

func (*IsNullExpr) exprNode() {}

type ArrayExpr struct {
	span
	Items []Expr
}

func (*ArrayExpr) exprNode() {}

type TupleExpr struct {
	span
	Items []Expr
}

func (*TupleExpr) exprNode() {}

type ParenExpr struct {
	span
	Expr Expr
}

func (*ParenExpr) exprNode() {}

type IndexExpr struct {
	span
	Expr  Expr
	Index Expr
}

func (*IndexExpr) exprNode() {}

type TupleAccessExpr struct {
	span
	Expr   Expr
	Member string
}

func (*TupleAccessExpr) exprNode() {}

type SubqueryExpr struct {
	span
	Query *SelectUnionStmt
}

func (*SubqueryExpr) exprNode() {}

type StarExpr struct {
	span
	Qualifier []string
	Except    []string
	ExceptRe  string
}

func (*StarExpr) exprNode() {}

type LambdaExpr struct {
	span
	Params []string
	Body   Expr
}

func (*LambdaExpr) exprNode() {}

type IntervalExpr struct {
	span
	Value Expr
	Unit  string
}

func (*IntervalExpr) exprNode() {}

type ExtractExpr struct {
	span
	Part string
	From Expr
}

func (*ExtractExpr) exprNode() {}

type TrimExpr struct {
	span
	Mode  string // "BOTH", "LEADING", "TRAILING"
	Chars Expr
	From  Expr
}

func (*TrimExpr) exprNode() {}

type SubstringExpr struct {
	span
	Expr Expr
	From Expr
	For  Expr
}

func (*SubstringExpr) exprNode() {}

type ParamExpr struct {
	span
	Name string
	Type *DataType
}

func (*ParamExpr) exprNode() {}

// TypeExpr wraps a nested *DataType inside a type argument list so callers can
// inspect the full recursive DataType tree (e.g. Array(Nullable(UInt64))).
type TypeExpr struct {
	span
	Type *DataType
}

func (*TypeExpr) exprNode() {}

// DataType represents a ClickHouse type expression such as UInt64,
// Nullable(String), Decimal(18, 4), Enum8('a' = 1, 'b' = 2), or Nested(...).
type DataType struct {
	span
	Name       string
	Args       []Expr
	TypeArgs   []*DataType
	NestedCols []NestedCol
	EnumValues []EnumValue
}

type NestedCol struct {
	Name string
	Type *DataType
}

type EnumValue struct {
	Name  string
	Value string
}
