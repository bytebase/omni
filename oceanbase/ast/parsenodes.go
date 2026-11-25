package ast

// -----------------------------------------------------------------------
// Statement nodes
// -----------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	Loc               Loc
	CTEs              []*CommonTableExpr // WITH clause (CTEs)
	OptimizerHints    []string           // /*+ ... */ optimizer hint bodies immediately after SELECT
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
	GroupBy           []*OrderByItem     // GROUP BY items (expr with optional ASC/DESC; OceanBase keeps ORDER BY's ByItem)
	WithRollup        bool               // GROUP BY ... WITH ROLLUP
	Having            ExprNode           // HAVING condition
	OrderBy           []*OrderByItem     // ORDER BY clause
	OrderByWithRollup bool               // ORDER BY ... WITH ROLLUP
	Limit             *Limit             // LIMIT / OFFSET
	// Approx is opt_approx, the APPROX or APPROXIMATE word of a vector search:
	// ORDER BY expr APPROX LIMIT n. OceanBase takes it only before a LIMIT.
	Approx bool
	// VectorIndexParams is opt_with_vector_index_parameters:
	// PARAMETERS(name = value, ...) after the LIMIT of a vector search.
	VectorIndexParams []*VectorIndexParam
	ForUpdate         *ForUpdate   // FOR UPDATE / SHARE
	WindowClause      []*WindowDef // WINDOW clause
	Into              *IntoClause  // INTO clause
	SetOp             SetOperation // UNION / INTERSECT / EXCEPT
	SetAll            bool         // ALL modifier for set operations
	Left              *SelectStmt  // left side of set operation
	Right             *SelectStmt  // right side of set operation
	// TableSource and ValuesSource are the TABLE or VALUES query primary this
	// node stands for. The source holds only its table or rows; the ORDER BY,
	// LIMIT, and INTO after it are this node's, except that an ORDER BY and
	// LIMIT after the last operand of a set operation are the set operation's.
	// OceanBase has no TABLE query, so the parser never sets TableSource.
	TableSource  *TableStmt  // TABLE query primary; never set by the parser
	ValuesSource *ValuesStmt // VALUES query primary
	// ParenSource is the inner query of a parenthesized query expression
	// ( SELECT ... ). When non-nil this node is a parentheses wrapper: its own
	// OrderBy/Limit hold the OUTER trailing clauses, while ParenSource carries
	// the inner query (which may have its own OrderBy/Limit). This keeps the two
	// scopes of "(SELECT 1 LIMIT 5) LIMIT 2" distinct.
	ParenSource *SelectStmt
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

// UnwrapParenSource unwraps parenthesized-query wrappers (ParenSource) and
// returns the underlying query node — a SELECT leaf or a set-operation node.
// Unlike walking to the leftmost leaf, it stops at a set-op so callers can
// inspect SetOp / Left / Right (e.g. recursive-CTE detection). Returns the
// input unchanged when it is not a wrapper.
func UnwrapParenSource(s *SelectStmt) *SelectStmt {
	for s != nil && s.ParenSource != nil {
		s = s.ParenSource
	}
	return s
}

// LeftmostQueryLeaf returns the leftmost output SELECT of a query expression: it
// unwraps ParenSource wrappers and walks set-operation left arms until it
// reaches the SELECT whose target list defines the result columns. For a bare
// SELECT it returns the input unchanged. Use this wherever a leftmost-leaf walk
// reads a target list (view columns, CTE / derived virtual tables, ORDER BY
// ordinals); UnwrapParenSource is the variant that stops at a set-op.
func LeftmostQueryLeaf(s *SelectStmt) *SelectStmt {
	for s != nil {
		switch {
		case s.ParenSource != nil:
			s = s.ParenSource
		case s.SetOp != SetOpNone && s.Left != nil:
			s = s.Left
		default:
			return s
		}
	}
	return s
}

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
	Loc       Loc
	IsReplace bool // REPLACE instead of INSERT
	Priority  InsertPriority
	Ignore    bool // INSERT IGNORE
	// Overwrite is OceanBase's INSERT OVERWRITE, whose target is a single table
	// and whose source is a query: insert_stmt's OVERWRITE arm writes no VALUES.
	Overwrite      bool
	Table          *TableRef     // target table
	Partitions     []string      // PARTITION (p0, p1, ...)
	Columns        []*ColumnRef  // explicit column list; nil without one, non-nil and empty for ()
	Values         [][]ExprNode  // VALUES rows; a row of () has no elements
	Select         *SelectStmt   // INSERT ... SELECT
	TableSource    *TableStmt    // INSERT ... TABLE table_name (MySQL 8.0.19+); never set by the parser
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
	// External is set by CREATE EXTERNAL TABLE's special_table_type, which is
	// one of TEMPORARY and EXTERNAL and never both.
	External    bool
	Table       *TableRef
	Columns     []*ColumnDef
	Constraints []*Constraint
	Options     []*TableOption
	Partitions  *PartitionClause
	Like        *TableRef   // CREATE TABLE ... LIKE
	Select      *SelectStmt // CREATE TABLE ... AS SELECT
	Ignore      bool        // IGNORE before SELECT in CTAS
	Replace     bool        // REPLACE before SELECT in CTAS
	ColumnGroup string      // WITH COLUMN GROUP (...): the list as written
}

func (s *CreateTableStmt) nodeTag()  {}
func (s *CreateTableStmt) stmtNode() {}

// AlterTableStmt represents an ALTER TABLE statement.
type AlterTableStmt struct {
	Loc   Loc
	Table *TableRef
	// External is set by ALTER EXTERNAL TABLE, which takes the same actions and
	// two of its own: ADD PARTITION (...) LOCATION 's' and DROP PARTITION LOCATION 's'.
	External bool
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
	ATAlterColumnSetNotNull  // never built: OceanBase's ALTER COLUMN has no SET NOT NULL
	ATAlterColumnDropNotNull // never built: OceanBase's ALTER COLUMN has no DROP NOT NULL
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
	ATAlterColumnVisible   // never built: OceanBase's ALTER COLUMN has no SET VISIBLE
	ATAlterColumnInvisible // never built: OceanBase's ALTER COLUMN has no SET INVISIBLE
	ATAlterIndexVisible
	ATAlterIndexInvisible
	// ATAlterIndexStorageCachePolicy is ALTER INDEX name STORAGE_CACHE_POLICY (...).
	ATAlterIndexStorageCachePolicy
	// ATAlterIndexParallel is ALTER INDEX name PARALLEL n | NOPARALLEL.
	ATAlterIndexParallel
	// ATRefresh is alter_table_action's REFRESH, which refreshes an external
	// table's metadata.
	ATRefresh
	// ATRemoveTTL is alter_table_action's REMOVE TTL.
	ATRemoveTTL
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
	ATPartitionBy
	ATDropTablegroup // DROP TABLEGROUP: the table leaves its table group
	// DROP SUBPARTITION and TRUNCATE SUBPARTITION, with PartitionNames.
	ATDropSubpartition
	ATTruncateSubpartition
	// {ALTER PARTITION | ALTER SUBPARTITION} names STORAGE_CACHE_POLICY = 'policy':
	// Name is PARTITION or SUBPARTITION, PartitionNames the names and
	// Option.Value the policy.
	ATAlterPartitionStorageCachePolicy
	// MODIFY PARTITION name ADD SUBPARTITION (...): Name is the partition and
	// PartitionDefs holds one PartitionDef whose SubPartitions are added.
	ATAddSubpartition
	// {ADD | DROP} COLUMN GROUP (...) [DELAYED]: Name is the list as written,
	// and Option, where DELAYED stands, names it.
	ATAddColumnGroup
	ATDropColumnGroup
	// ALTER EXTERNAL TABLE's own actions: ADD PARTITION (...) LOCATION 's'
	// holds the added partitions in PartitionValues and the path in Location;
	// DROP PARTITION LOCATION 's' holds only the path.
	ATAddExternalPartition
	ATDropExternalPartition
)

// AlterTableCmd represents a single ALTER TABLE operation.
// ConstraintKeyword is the keyword an ALTER TABLE action writes before the
// constraint it names. The grammar spells one action two ways -- DROP CHECK /
// ALTER CHECK beside DROP CONSTRAINT / ALTER CONSTRAINT -- and OceanBase answers
// a different error for each when no such constraint is there: 3821 "Check
// constraint '%s' is not found in the table." for CHECK and 3940 "Constraint
// '%s' does not exist." for CONSTRAINT. The keyword is kept, not folded, so the
// catalog answers what the statement wrote; it also tells DROP CHECK from DROP
// CONSTRAINT, which name a check alone and any constraint of that name.
//
// Ported from mysql, which measured the same two errors on mysql:8.0.46.
type ConstraintKeyword int

const (
	// ConstraintKeywordNone is an action that wrote no keyword of its own:
	// DROP COLUMN, DROP INDEX and the rest.
	ConstraintKeywordNone ConstraintKeyword = iota
	// ConstraintKeywordPrimaryKey is DROP PRIMARY KEY.
	ConstraintKeywordPrimaryKey
	// ConstraintKeywordForeignKey is DROP FOREIGN KEY.
	ConstraintKeywordForeignKey
	// ConstraintKeywordCheck is DROP CHECK or ALTER CHECK.
	ConstraintKeywordCheck
	// ConstraintKeywordConstraint is DROP CONSTRAINT or ALTER CONSTRAINT.
	ConstraintKeywordConstraint
)

// String is the keyword as the statement wrote it, upper-cased, "" for
// ConstraintKeywordNone.
func (k ConstraintKeyword) String() string {
	switch k {
	case ConstraintKeywordPrimaryKey:
		return "PRIMARY KEY"
	case ConstraintKeywordForeignKey:
		return "FOREIGN KEY"
	case ConstraintKeywordCheck:
		return "CHECK"
	case ConstraintKeywordConstraint:
		return "CONSTRAINT"
	}
	return ""
}

type AlterTableCmd struct {
	Loc            Loc
	Type           AlterTableCmdType
	Name           string       // column/constraint/index name
	NewName        string       // for RENAME operations
	NewTable       *TableRef    // for RENAME TABLE
	Column         *ColumnDef   // for ADD/MODIFY/CHANGE COLUMN
	Columns        []*ColumnDef // for ADD (col1, col2, ...) multi-column form
	Constraint     *Constraint  // for ADD CONSTRAINT
	Option         *TableOption // for table options
	DefaultExpr    ExprNode     // for ALTER COLUMN SET DEFAULT expr
	IfExists       bool
	First          bool             // FIRST positioning
	After          string           // AFTER column positioning
	Before         string           // BEFORE column positioning (opt_position_column)
	PartitionNames []string         // for partition operations (DROP/TRUNCATE/ANALYZE/CHECK/OPTIMIZE/REBUILD/REPAIR/REORGANIZE/DISCARD/IMPORT)
	AllPartitions  bool             // for partition operations using ALL
	PartitionDefs  []*PartitionDef  // for ADD PARTITION / REORGANIZE PARTITION ... INTO
	Number         int              // for COALESCE PARTITION number
	ExchangeTable  *TableRef        // for EXCHANGE PARTITION ... WITH TABLE
	WithValidation *bool            // for EXCHANGE PARTITION: true=WITH VALIDATION, false=WITHOUT VALIDATION, nil=not specified
	OrderByItems   []*OrderByItem   // for ORDER BY operation
	PartitionBy    *PartitionClause // for PARTITION BY (repartition)
	// PartitionValues holds an external table's added partitions, each a column
	// and its value, and Location the path that follows them (or the one a DROP
	// PARTITION names).
	PartitionValues []*Assignment
	Location        string
	// IndexOptions holds ALTER INDEX's STORAGE_CACHE_POLICY attributes.
	IndexOptions []*IndexOption
	// Parallel is ALTER INDEX's PARALLEL n or NOPARALLEL.
	Parallel *ParallelOption
	// ConstraintNames is DROP CONSTRAINT's and DROP CHECK's '(' name_list ')'.
	ConstraintNames []string
	// ConstraintKeyword is the word the action wrote before the constraint it
	// names, which tells the spellings apart: DROP/ALTER CHECK answer 3821 and
	// DROP/ALTER CONSTRAINT 3940, and a CHECK drop takes a check alone.
	ConstraintKeyword ConstraintKeyword
	// LobStorage is lob_storage_clause, which follows an ALTER TABLE ADD's column
	// list: JSON (column) STORE AS (CHUNK n).
	LobStorage *LobStorageClause
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
	Vector      bool // CREATE VECTOR INDEX
	IfNotExists bool
	IndexName   string
	// IndexSchema and IndexCatalog are the qualifiers CREATE INDEX lets its name carry, as
	// in d.i and c.d.i, and IndexDBLink the name after @. OceanBase requires the schema to
	// be the table's.
	IndexSchema  string
	IndexCatalog string
	IndexDBLink  string
	Table        *TableRef
	IndexType    string // BTREE, HASH
	Columns      []*IndexColumn
	Options      []*IndexOption
	Algorithm    string
	Lock         string
	ColumnGroup  string // WITH COLUMN GROUP (...): the list as written
}

func (s *CreateIndexStmt) nodeTag()  {}
func (s *CreateIndexStmt) stmtNode() {}

// LobStorageClause is lob_storage_clause, which OceanBase writes after an ALTER
// TABLE ADD's column list: `JSON (column) STORE AS (CHUNK n)`.
//
// Ref: lob_storage_clause, lob_storage_parameters, lob_storage_parameter and
// lob_chunk_size in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 25494-25534).
type LobStorageClause struct {
	Loc    Loc
	Column string // column_name
	Chunk  string // lob_chunk_size's last CHUNK value, as written
	// Chunks holds every lob_storage_parameters entry in the order written: the
	// grammar lets a clause repeat `CHUNK n`, and each occurrence is kept.
	Chunks []string
}

func (c *LobStorageClause) nodeTag() {}

// IndexColumn represents a column in an index definition.
type IndexColumn struct {
	Loc        Loc
	Expr       ExprNode // column name or expression
	Functional bool     // true when parsed from functional key-part syntax: (expr)
	Length     int      // prefix length
	Desc       bool     // DESC ordering
	// ColumnID and HasColumnID are opt_column_id's `ID n`, which pins the index
	// column's storage id.
	ColumnID    int64
	HasColumnID bool
}

func (c *IndexColumn) nodeTag() {}

// IndexOption represents an index option.
type IndexOption struct {
	Loc   Loc
	Name  string
	Value Node // heterogeneous: may be *String, *Integer, or other literal nodes
	// Properties holds a key = value list: PARSER_PROPERTIES on a fulltext index, the
	// WITH (...) a vector index takes its parameters from, or STORAGE_CACHE_POLICY. nil
	// for every other option.
	Properties []*ParserProperty
	// Columns holds the names in a STORING (...), as written.
	Columns []string
}

func (o *IndexOption) nodeTag() {}

// ParserProperty is one key = value of an OceanBase fulltext index's
// PARSER_PROPERTIES, such as min_token_size = 3 or ik_mode = 'smart', or of an index's
// WITH (...), such as distance = l2 or m = 16.
type ParserProperty struct {
	Loc   Loc
	Name  string // the key, lower-cased: min_token_size, ngram_token_size, ik_mode, distance, ...
	Value Node   // *IntLit for a size, *StringLit for a table name, ik_mode, or a WITH word or decimal
}

func (p *ParserProperty) nodeTag() {}

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
	// TableID and HasTableID are opt_table_id's `TABLE_ID = n`, which pins the
	// view's id: the server parses it and refuses to run it.
	TableID    int64
	HasTableID bool
}

func (s *CreateViewStmt) nodeTag()  {}
func (s *CreateViewStmt) stmtNode() {}

// CreateMaterializedViewStmt represents a CREATE MATERIALIZED VIEW statement.
//
// Ref: create_mview_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 9766).
type CreateMaterializedViewStmt struct {
	Loc          Loc
	Name         *TableRef
	Columns      []string                    // opt_mv_column_list's column_name_list
	PrimaryKey   *MaterializedViewPrimaryKey // PRIMARY KEY (...) beside the columns, if any
	Options      []*TableOption              // opt_table_option_list
	Partitions   *PartitionClause            // opt_partition_option
	ColumnGroup  string                      // WITH COLUMN GROUP (...), as written
	MviewOptions *MviewOptions               // create_mview_opts
	Select       *SelectStmt
	SelectText   string // raw SQL text of the AS view_select_stmt portion
	CheckOption  string // CASCADED, LOCAL
}

func (s *CreateMaterializedViewStmt) nodeTag()  {}
func (s *CreateMaterializedViewStmt) stmtNode() {}

// MaterializedViewPrimaryKey is the PRIMARY KEY (...) clause of a materialized
// view's column list.
type MaterializedViewPrimaryKey struct {
	Loc       Loc
	Columns   []string
	IndexType string // BTREE or HASH, from USING; empty when unspecified
	Comment   string
}

func (p *MaterializedViewPrimaryKey) nodeTag() {}

// MviewOptions holds create_mview_opts: the refresh clause and the two
// ENABLE/DISABLE switches that may follow it in either order.
type MviewOptions struct {
	Loc                Loc
	Refresh            *MviewRefreshInfo // nil when no refresh clause is written
	OnQueryComputation string            // ENABLE or DISABLE; empty when absent
	QueryRewrite       string            // ENABLE or DISABLE; empty when absent
}

func (o *MviewOptions) nodeTag() {}

// MviewRefreshInfo is one materialized view refresh clause. Which fields are
// set depends on the form it came from: NEVER REFRESH sets Never; REFRESH
// <method> sets Method and the optional clauses after it; an ALTER MATERIALIZED
// VIEW action that names an interval sets StartWith and Next; one that names a
// nested refresh mode sets NestedRefresh.
type MviewRefreshInfo struct {
	Loc           Loc
	Never         bool     // NEVER REFRESH
	Method        string   // FAST, COMPLETE, FORCE
	NestedRefresh string   // INDIVIDUAL, INCONSISTENT, CONSISTENT
	On            string   // DEMAND, COMMIT, STATEMENT
	StartWith     ExprNode // START WITH <bit_expr>
	Next          ExprNode // NEXT <bit_expr>
}

func (r *MviewRefreshInfo) nodeTag() {}

// ParallelOption is parallel_option: PARALLEL [=] n or NOPARALLEL.
type ParallelOption struct {
	Loc        Loc
	NoParallel bool
	Degree     int64 // PARALLEL's degree; 0 for NOPARALLEL
}

func (o *ParallelOption) nodeTag() {}

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
	// IsNull marks an option written with the keyword NULL, such as
	// `TABLEGROUP = NULL`, as opposed to a quoted identifier named `NULL`.
	IsNull bool
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
	ColumnFormat  string // COLUMN_FORMAT {FIXED | DYNAMIC | DEFAULT}
	// ID and HasID are column_attribute's `ID INTNUM`, which pins the column's
	// storage id. OceanBase admits the same attribute on a generated column.
	ID    int64
	HasID bool
	// Chunk and HasChunk are column_attribute's lob_chunk_size, `CHUNK n` or
	// `CHUNK 'size'`, which the server parses and refuses to run.
	Chunk    string
	HasChunk bool
	// OrigDefault is column_attribute's `ORIG_DEFAULT now_or_signed_literal`,
	// which the server parses and refuses to run.
	OrigDefault              ExprNode
	HasOrigDefault           bool
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
	// ID and HasID are generated_column_attribute's `ID INTNUM`, which pins the
	// generated column's storage id.
	ID    int64
	HasID bool
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
	// ConstrVector is OceanBase's VECTOR {INDEX | KEY}, a table element only: 4.3.5.6
	// answers 1064 for ADD VECTOR INDEX and for a CONSTRAINT symbol in front of it.
	ConstrVector
)

// Constraint represents a table-level constraint.
//
// Name is the constraint's name: the CONSTRAINT symbol of a FOREIGN KEY or
// CHECK, and the index name of a key, which for UNIQUE is the name after
// UNIQUE [KEY|INDEX] when one is written and else the symbol. IndexName is
// the name written after FOREIGN KEY, kept apart from the symbol as
// sql_yacc.yy's PT_foreign_key_definition keeps them: the foreign key is
// named after the symbol, or generated without one, never after this
// name; its index is named after the symbol, or this name without one
// (verified on 8.0.46).
type Constraint struct {
	Loc          Loc
	Type         ConstraintType
	Name         string
	IndexName    string         // for FOREIGN KEY: the name after FOREIGN KEY [KEY]
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
	ColumnGroup  string // WITH COLUMN GROUP (...) of a key: the list as written
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
	// TTLColumns is the column each entry of a TTL (col + INTERVAL n unit, ...)
	// names, as the parser read them, so a consumer needs no second reading of
	// Value. Nil for every other option.
	TTLColumns []string
}

func (o *TableOption) nodeTag() {}

// PartitionClause represents a PARTITION BY clause.
type PartitionClause struct {
	Loc        Loc
	Type       PartitionType
	Linear     bool // LINEAR HASH or LINEAR KEY
	Expr       ExprNode
	Columns    []string
	Algorithm  int // ALGORITHM={1|2} for KEY partitioning (0 if not specified)
	NumParts   int
	Partitions []*PartitionDef
	// SubPartType is the SUBPARTITION BY type. Its zero value is PartitionRange, so a
	// clause is subpartitioned when SubPartExpr or SubPartColumns is set, as OceanBase's
	// RANGE and LIST subpartitioning always sets one of them.
	SubPartType    PartitionType
	SubPartExpr    ExprNode // SUBPARTITION BY {HASH | RANGE | LIST} (expr)
	SubPartColumns []string // SUBPARTITION BY {KEY | RANGE COLUMNS | LIST COLUMNS} (columns)
	SubPartAlgo    int      // ALGORITHM={1|2} for SUBPARTITION BY KEY (0 if not specified)
	NumSubParts    int      // SUBPARTITIONS num
	// SubPartTemplate holds OceanBase's SUBPARTITION TEMPLATE (...) list, which
	// gives every partition the same subpartitions. MySQL has no template and
	// repeats them per partition instead.
	SubPartTemplate []*SubPartitionDef
	// SplitSize is a RANGE partitioning's SIZE('size'), the tablet size it splits
	// at (auto_split_tablet_size_option).
	SplitSize string
	// AutoSize is auto_partition_option's PARTITION SIZE value, as written, AUTO
	// included; empty when the table is not partitioned automatically.
	AutoSize string
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
	ValuesIn      bool // true when Values came from VALUES IN rather than VALUES LESS THAN
	Options       []*TableOption
	SubPartitions []*SubPartitionDef
}

// SubPartitionDef represents a subpartition definition.
type SubPartitionDef struct {
	Loc  Loc
	Name string
	// Values is an OceanBase RANGE or LIST subpartition's VALUES LESS THAN or VALUES
	// IN, kept as PartitionDef.Values keeps a partition's; nil for HASH and KEY.
	Values Node
	// ValuesIn is set when Values came from VALUES IN rather than VALUES LESS THAN,
	// which a multi-column bound's list does not tell apart.
	ValuesIn bool
	Options  []*TableOption
}

func (d *SubPartitionDef) nodeTag() {}

func (d *PartitionDef) nodeTag() {}

// DataType represents a MySQL data type.
type DataType struct {
	Loc    Loc
	Name   string // INT, VARCHAR, DECIMAL, etc.
	Length int    // type length (e.g., VARCHAR(255))
	// HasLength records that a length was written, since CHAR(0) is no CHAR: Length is 0 for both.
	HasLength  bool
	Scale      int // decimal scale (e.g., DECIMAL(10,2))
	Unsigned   bool
	Zerofill   bool
	Charset    string    // CHARACTER SET
	Collate    string    // COLLATE
	Binary     bool      // standalone BINARY modifier
	EnumValues []string  // for ENUM and SET types
	ArrayDim   int       // not used in MySQL, but here for consistency
	ArrayElem  *DataType // the element of an ARRAY, itself an ARRAY when they nest; a MAP's value
	MapKey     *DataType // a MAP's key
	SRID       int       // Spatial Reference ID for spatial types (valid only when HasSRID)
	// HasSRID records whether the column declared an explicit SRID clause at all. MySQL
	// treats `SRID 0` as present-and-distinct from a column with no SRID clause (SHOW
	// CREATE echoes `/*!80003 SRID 0 */`; information_schema reports SRS_ID=0 vs NULL), so
	// presence cannot be inferred from SRID != 0 — a `SRID 0` column has SRID==0 yet is set.
	HasSRID bool
	// FromStringValue is data_type's STRING_VALUE arm, a quoted literal in the
	// type's place: the parser takes it and the resolver refuses it with 5874.
	FromStringValue bool
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
	Loc        Loc
	Op         UnaryOp
	Operand    ExprNode
	OriginalOp string // Original operator text (e.g., "!" vs "NOT") for alias generation.
}

func (e *UnaryExpr) nodeTag()  {}
func (e *UnaryExpr) exprNode() {}

// FuncCallExpr represents a function call.
type FuncCallExpr struct {
	Loc        Loc
	Name       string
	Schema     string // schema-qualified name
	Args       []ExprNode
	Quantifier AggQuantifier // COUNT(DISTINCT ...), SUM(ALL ...)
	Star       bool          // COUNT(*)
	HasParens  bool          // true if explicit parentheses were written (e.g., NOW() vs CURRENT_TIMESTAMP)
	OrderBy    []*OrderByItem
	Over       *WindowDef        // OVER clause for window functions
	Separator  ExprNode          // GROUP_CONCAT SEPARATOR value (nil = not specified)
	JSONValue  *JSONValueOptions // JSON_VALUE's clauses after its path (nil when written with none)
	// JSONValueArm marks a call read through JSON_VALUE's own arm of the grammar, which has a
	// stored form of its own even with no clause written. Quoted or qualified, the name is an
	// ordinary call to whatever function has that name.
	JSONValueArm bool
	JSONQuery    *JSONQueryOptions // JSON_QUERY's clauses after its path (nil when written with none)
	// JSONQueryArm marks a call read through JSON_QUERY's own arm of the grammar, as
	// JSONValueArm does for JSON_VALUE.
	JSONQueryArm bool
	Using        string // CHAR(n, ... USING charset_name): the character set as written ("" = no USING)
	UsingQuoted  bool   // Using was written as a string, as ConvertExpr.CharsetQuoted
	// NullTreatment is a window function's RESPECT NULLS or IGNORE NULLS, as
	// "RESPECT" or "IGNORE", whether written inside the argument list or after
	// it. Only FIRST_VALUE, LAST_VALUE, LEAD, LAG and NTH_VALUE take it.
	NullTreatment string
	// NullTreatmentInArgs records which of the two places it was written in:
	// LEAD's and LAG's parameter list (and FIRST_VALUE's, LAST_VALUE's on some
	// spellings) end with it, and the trailing form follows the ')'.
	NullTreatmentInArgs bool
	// FromFirstLast is NTH_VALUE's FROM FIRST or FROM LAST, as "FIRST" or
	// "LAST". It follows the argument list and precedes NullTreatment.
	FromFirstLast string
}

// JSONValueOptions is what JSON_VALUE takes after its path. The call's document and path are
// its Args.
//
// Ref: JSON_VALUE_SYM in sql_parser_mysql_mode.y (v4.3.5_CE_BP6).
type JSONValueOptions struct {
	Loc       Loc
	Returning *DataType          // RETURNING type (nil = not specified)
	Truncate  bool               // TRUNCATE
	ASCII     bool               // ASCII
	OnEmpty   *JSONValueResponse // ... ON EMPTY (nil = not specified)
	OnError   *JSONValueResponse // ... ON ERROR (nil = not specified)
}

// JSONValueResponse is one ON EMPTY or ON ERROR: NULL, ERROR, or DEFAULT and a literal.
type JSONValueResponse struct {
	Loc     Loc
	Error   bool     // ERROR
	Default ExprNode // DEFAULT literal (nil with Error unset = NULL)
}

func (n *JSONValueOptions) nodeTag()  {}
func (n *JSONValueResponse) nodeTag() {}

// JSONQueryOptions is what JSON_QUERY takes after its path, in the order its grammar has
// them. The call's document and path are its Args. A clause that is a choice of words holds
// the words as OceanBase prints them, in capitals.
//
// Ref: json_query_expr in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, line 24917).
type JSONQueryOptions struct {
	Loc        Loc
	Returning  *DataType // RETURNING type (nil = not specified)
	Truncate   bool      // TRUNCATE
	Scalars    string    // ALLOW or DISALLOW, of ALLOW SCALARS and DISALLOW SCALARS ("" = not specified)
	Pretty     bool      // PRETTY
	ASCII      bool      // ASCII
	Wrapper    string    // the whole clause, such as WITHOUT ARRAY WRAPPER ("" = not specified)
	Asis       bool      // ASIS
	OnError    string    // ERROR, NULL, EMPTY, EMPTY ARRAY or EMPTY OBJECT ("" = not specified)
	OnEmpty    string    // as OnError
	OnMismatch string    // ERROR, NULL or DOT ("" = not specified)
	Multivalue bool      // MULTIVALUE
}

func (n *JSONQueryOptions) nodeTag() {}

// AggQuantifier is the prefix an aggregate's argument list carries. It is the
// only record of what was written, so a consumer asking whether duplicates are
// removed asks RemovesDuplicates rather than comparing to AggQuantDistinct:
// DISTINCTROW is a synonym, and the DISTINCT ALL forms de-duplicate too.
type AggQuantifier uint8

const (
	AggQuantNone AggQuantifier = iota
	AggQuantDistinct
	AggQuantDistinctRow
	AggQuantAll
	AggQuantDistinctAll
	AggQuantDistinctRowAll
	// AggQuantUnique is COUNT's UNIQUE spelling. It de-duplicates like
	// DISTINCT, but the written word is kept because an automatic column
	// name is the text as written.
	AggQuantUnique
)

// String returns the quantifier as written, empty for AggQuantNone.
func (q AggQuantifier) String() string {
	switch q {
	case AggQuantDistinct:
		return "DISTINCT"
	case AggQuantDistinctRow:
		return "DISTINCTROW"
	case AggQuantAll:
		return "ALL"
	case AggQuantDistinctAll:
		return "DISTINCT ALL"
	case AggQuantDistinctRowAll:
		return "DISTINCTROW ALL"
	case AggQuantUnique:
		return "UNIQUE"
	}
	return ""
}

// RemovesDuplicates reports whether the quantifier makes the aggregate read
// each distinct value once. ALL is the default and changes nothing; the
// DISTINCT ALL forms de-duplicate, which the grammar reads as a DISTINCT over
// an in_sum_expr whose own opt_all follows it.
func (q AggQuantifier) RemovesDuplicates() bool {
	switch q {
	case AggQuantDistinct, AggQuantDistinctRow, AggQuantDistinctAll, AggQuantDistinctRowAll, AggQuantUnique:
		return true
	}
	return false
}

func (e *FuncCallExpr) nodeTag()  {}
func (e *FuncCallExpr) exprNode() {}

// SubqueryExpr represents a subquery expression.
type SubqueryExpr struct {
	Loc        Loc
	Select     *SelectStmt
	Exists     bool     // EXISTS (SELECT ...)
	Lateral    bool     // LATERAL derived table
	Alias      string   // AS alias (for derived tables)
	Columns    []string // derived table column alias list: (SELECT ...) AS t(c1, c2)
	Quantifier string   // "ANY" / "SOME" / "ALL" for `expr op {ANY|SOME|ALL} (subquery)`; empty for scalar subquery
	// Flashback is the use_flashback a derived table may carry before its alias,
	// as in `FROM (SELECT ...) AS OF SNAPSHOT 1 t` (table_factor's
	// select_with_parens use_flashback alternative).
	Flashback *TableFlashback
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
	// Array is set for the multi-valued-index form CAST(expr AS type ARRAY)
	// (MySQL 8.0.17+). It maps to the grammar's opt_array_cast, an optional
	// trailing ARRAY keyword after the cast target type. ARRAY is a
	// non-reserved (contextual) keyword; it only attaches to CAST, never CONVERT.
	Array bool
	// Ignore is sys_view_cast_opt's IGNORE, which makes the cast a value that
	// cannot be converted the default instead of an error.
	Ignore bool
}

func (e *CastExpr) nodeTag()  {}
func (e *CastExpr) exprNode() {}

// ExtractExpr represents EXTRACT(unit FROM expr).
type ExtractExpr struct {
	Loc  Loc
	Unit string // DAY, HOUR, MINUTE, SECOND, MONTH, YEAR, etc.
	Expr ExprNode
}

func (e *ExtractExpr) nodeTag()  {}
func (e *ExtractExpr) exprNode() {}

// PositionExpr represents POSITION(substr IN str). The keyword separates the
// arguments and the comma form is a syntax error, so the call cannot be a
// FuncCallExpr: deparse would write POSITION(a, b), which no engine parses.
type PositionExpr struct {
	Loc    Loc
	Substr ExprNode
	Str    ExprNode
}

func (e *PositionExpr) nodeTag()  {}
func (e *PositionExpr) exprNode() {}

// -----------------------------------------------------------------------
// Reference and literal nodes
// -----------------------------------------------------------------------

// ColumnRef represents a column reference (possibly qualified).
type ColumnRef struct {
	Loc    Loc
	Table  string // table name or alias (may be empty)
	Schema string // schema/database name (may be empty)
	// Catalog is the first part of a four-part name, `catalog.db.tbl.col`
	// (column_ref's id_dot_id_dot_id '.' column_name arm). It is only written
	// with Schema and Table.
	Catalog string
	Column  string // column name
	Star    bool   // table.*
}

func (r *ColumnRef) nodeTag()  {}
func (r *ColumnRef) exprNode() {}

// TableRef represents a table reference (possibly qualified, with alias).
type TableRef struct {
	Loc    Loc
	Schema string // database name
	// Catalog is the first part of a three-part name, `catalog.db.tbl`
	// (normal_relation_factor's id_dot_id_dot_id arm). It is only written with
	// Schema.
	Catalog string
	Name    string // table name
	// Dblink is opt_dblink's `@name`, which reads the table through a database
	// link.
	Dblink     string
	Alias      string   // AS alias
	Partitions []string // PARTITION (p0, p1, ...)
	// PartitionValues is an external table's PARTITION (col = value, ...).
	PartitionValues []*Assignment
	Sample          *TableSample    // SAMPLE ... (percent) [SEED (n)]
	Flashback       *TableFlashback // AS OF SNAPSHOT expr | AS OF PROCTIME()
	IndexHints      []*IndexHint
}

func (r *TableRef) nodeTag()   {}
func (r *TableRef) tableExpr() {}

// RemoteSequenceExpr is a sequence value read through a database link:
// `t.col@link`, `db.t.col@link` or `catalog.db.t.col@link`. The name it reads
// is a column reference, and Dblink the `@link` after it.
//
// Ref: simple_expr's USER_VARIABLE arms in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, lines 1938-1990), which build T_REMOTE_SEQUENCE.
type RemoteSequenceExpr struct {
	Loc    Loc
	Ref    *ColumnRef // the qualified name of the sequence
	Dblink string     // the @name after it
}

func (e *RemoteSequenceExpr) nodeTag()  {}
func (e *RemoteSequenceExpr) exprNode() {}

// QuantifiedExpr is any_expr's parenthesized-expression arm: the right side of
// `expr comp_op {ANY|SOME|ALL} (expr)`, which holds one expression and no
// query. `x = ANY(a)` is the array_contains form the server builds, and
// `x = ALL(a)` compares against the one value.
//
// Ref: any_expr and bool_pri's `COMP_EQ sub_query_flag any_expr` arm in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 1514-1527 and 2155-2169).
type QuantifiedExpr struct {
	Loc        Loc
	Quantifier string // "ANY", "SOME" or "ALL"
	Expr       ExprNode
}

func (e *QuantifiedExpr) nodeTag()  {}
func (e *QuantifiedExpr) exprNode() {}

// OdbcEscapeExpr is the ODBC escape sequence simple_expr carries as
// `'{' relation_name expr '}'`, such as `{fn CONCAT(a, b)}`.
//
// Ref: simple_expr's '{' relation_name expr '}' arm in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 1948), which builds T_ODBC_ESCAPE_SEQUENCES.
type OdbcEscapeExpr struct {
	Loc  Loc
	Name string   // the word after '{'
	Expr ExprNode // the expression after it
}

func (e *OdbcEscapeExpr) nodeTag()  {}
func (e *OdbcEscapeExpr) exprNode() {}

// LambdaExpr is lambda_expr's `params -> expr`, the first argument of
// ARRAY_MAP, ARRAY_FIRST, ARRAY_SORTBY and ARRAY_FILTER.
//
// Ref: lambda_expr and lambda_expr_params in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, lines 3600-3625).
type LambdaExpr struct {
	Loc    Loc
	Params []string
	Body   ExprNode
}

func (e *LambdaExpr) nodeTag()  {}
func (e *LambdaExpr) exprNode() {}

// ExternalTableRef is a table read from a location rather than from a stored
// table. tbl_name takes it as url_expr: a quoted path with its FORMAT and
// PATTERN, or OceanBase's SOURCE(...) and FILES(...) table functions, with
// opt_as_alias after it.
//
// Ref: url_expr and url_table_function_expr in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, lines 14499-14549).
type ExternalTableRef struct {
	Loc Loc
	// Function is "SOURCE" or "FILES" for a table function, and "" for a quoted
	// path. A path names its location itself, and SOURCE(...) carries properties
	// rather than a path; FILES(...) names its location with LOCATION.
	Function string
	Path     string              // the location: the quoted path, or FILES' LOCATION 's'
	Source   []*ExternalProperty // SOURCE(...): external_properties_list
	Format   []*ExternalProperty // FORMAT(...): external_file_format_list
	Pattern  string              // PATTERN 's'
	// PatternSet marks the clause present, since PATTERN '' is a pattern and not
	// an absent one.
	PatternSet bool
	Alias      string // opt_as_alias: AS relation_name
}

func (r *ExternalTableRef) nodeTag()   {}
func (r *ExternalTableRef) tableExpr() {}

// ExternalProperty is one entry of a SOURCE(...) or FORMAT(...) list: a word
// and the value the grammar gives it, as TYPE = 'csv' or SKIP_HEADER = 1.
// NULL_IF and COLUMN_BLOOM_FILTER keep their parenthesized list as a RowExpr.
type ExternalProperty struct {
	Loc   Loc
	Name  string   // the property's word, in upper case
	Value ExprNode // its value
}

func (p *ExternalProperty) nodeTag() {}

// TableSample is the sample_clause of a table reference, with its seed:
// SAMPLE [BLOCK | HYBRID] [ALL | BASE | INCR] (percent) [SEED (n)].
type TableSample struct {
	Loc     Loc
	Method  string // "", "BLOCK" or "HYBRID"
	Scope   string // "", "ALL", "BASE" or "INCR"
	Percent string // the number as written
	Seed    string // the SEED number as written, "" without one
}

func (s *TableSample) nodeTag() {}

// TableFlashback is the use_flashback of a table reference: the table as it
// stood at a snapshot, or as it stands when the row is processed.
type TableFlashback struct {
	Loc      Loc
	Snapshot ExprNode // AS OF SNAPSHOT expr
	Proctime bool     // AS OF PROCTIME()
}

func (f *TableFlashback) nodeTag() {}

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
//
// Value is the int64 form used by consumers that need an integer (constant folding,
// deparse). Text preserves the EXACT source digits, which matters for an unsigned literal
// larger than math.MaxInt64 (e.g. a BIGINT UNSIGNED DEFAULT 18446744073709551615): Value
// would be clamped to math.MaxInt64 by the lexer, so anything that must round-trip the
// literal exactly (a column default) reads Text. Text may be empty for synthesized nodes;
// callers fall back to Value when it is.
type IntLit struct {
	Loc   Loc
	Value int64
	Text  string
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
	// National marks N'...', a string in the national character set (utf8mb3).
	National bool
	// Concatenated marks a literal folded from ADJACENT quoted segments per
	// MySQL's adjacent string-literal concatenation ('a' 'b' → "ab"). Value
	// holds the folded result; FirstSegment holds the first segment's value,
	// because MySQL derives the implicit output column name from the first
	// segment only (SELECT 'a' 'b' → column named "a" with value "ab", and an
	// empty first segment falls back to the positional Name_exp_N form).
	Concatenated bool
	FirstSegment string
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
	Scope  string // SET only: the scope word written on this assignment
	Column *ColumnRef
	Value  ExprNode
	// CharsetForm is set for SET NAMES and the charset forms, where Column holds the form's word
	// (NAMES, COLLATE or CHARACTER SET) and names no variable.
	CharsetForm bool
}

func (a *Assignment) nodeTag() {}

// OrderDirection preserves whether ORDER BY direction was omitted or explicit.
type OrderDirection int

const (
	OrderDirectionDefault OrderDirection = iota
	OrderDirectionAsc
	OrderDirectionDesc
)

// OrderByItem represents an ORDER BY item.
type OrderByItem struct {
	Loc        Loc
	Expr       ExprNode
	Direction  OrderDirection
	Desc       bool
	NullsFirst *bool // nil means default
}

func (o *OrderByItem) nodeTag() {}

// GroupingSet is one item of a GROUP BY list that is a grouping set rather than
// an expression: ROLLUP (list), CUBE (list), GROUPING SETS (list, ...) and the
// empty set `()`. An item of ROLLUP, CUBE or GROUPING SETS is an expression or
// the empty set, except in a GROUPING SETS list, where an item may be a ROLLUP
// or a CUBE of its own.
//
// Ref: sort_key_for_group_by, rollup_clause, cube_clause, grouping_sets_clause,
// grouping_sets and gs_group_by_expr in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, lines 13531-13623).
type GroupingSet struct {
	Loc Loc
	// Kind is "ROLLUP", "CUBE", "GROUPING SETS", or "" for the empty set.
	Kind  string
	Items []ExprNode
}

func (e *GroupingSet) nodeTag()  {}
func (e *GroupingSet) exprNode() {}

// VectorIndexParam is one vec_index_param: the name of a vector index parameter
// and its value, as written.
//
// Ref: vec_index_params, vec_index_param and vec_index_param_value in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 24877-24910).
type VectorIndexParam struct {
	Loc   Loc
	Name  string
	Value string
	// IsName is set when the value is a relation_name rather than a number or a
	// boolean, so deparse quotes it as the identifier it is.
	IsName bool
}

func (v *VectorIndexParam) nodeTag() {}

// FileOption is one file_option of INTO OUTFILE: SINGLE [=] bool,
// MAX_FILE_SIZE [=] size or BUFFER_SIZE [=] size. Value holds the value as
// written.
//
// Ref: file_option, file_opt, file_option_list and file_size_const in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 12042-12102).
type FileOption struct {
	Loc   Loc
	Name  string // "SINGLE", "MAX_FILE_SIZE" or "BUFFER_SIZE"
	Value string
	// ValueIsString is set when the value was a string literal, which a size
	// takes beside an integer and a stored statement keeps quoted.
	ValueIsString bool
}

func (o *FileOption) nodeTag() {}

// Limit represents LIMIT and OFFSET.
type Limit struct {
	Loc    Loc
	Count  ExprNode // LIMIT value
	Offset ExprNode // OFFSET value
	// OffsetKeyword is true when the offset was written as
	// "LIMIT count OFFSET offset" rather than the comma form
	// "LIMIT offset, count". The two forms place a `?` parameter marker in
	// opposite positions, so deparse must reproduce the original form to keep
	// the positional binding order. False when there is no offset.
	OffsetKeyword bool
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
	JoinFull         // FULL [OUTER] JOIN
	JoinNaturalFull  // NATURAL FULL [OUTER] JOIN
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
	// Each option written, even as '': an empty string turns a default off, so it
	// is not the same as no option.
	HasOutfile            bool
	HasDumpfile           bool
	HasFieldsTerminatedBy bool
	HasFieldsEnclosedBy   bool
	HasFieldsEscapedBy    bool
	HasLinesStartingBy    bool
	HasLinesTerminatedBy  bool
	// PartitionBy is OceanBase's file_partition_opt: PARTITION BY bit_expr.
	PartitionBy ExprNode
	// ExternalFormat is OceanBase's FORMAT opt_equal_mark '('
	// external_file_format_list ')', the format the file is written in.
	ExternalFormat []*ExternalProperty
	// FileOptions is file_opt: SINGLE, MAX_FILE_SIZE and BUFFER_SIZE, which
	// follow the fields and lines clauses.
	FileOptions []*FileOption

	// Aurora MySQL: SELECT ... INTO OUTFILE S3 's3-uri' export options.
	// Ref: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Integrating.SaveIntoS3.html
	OutfileS3     bool   // INTO OUTFILE S3 (Outfile then holds the S3 URI)
	Format        string // FORMAT {CSV | TEXT}; empty when omitted
	Header        bool   // FORMAT ... HEADER
	Manifest      string // MANIFEST {ON | OFF}; empty when omitted
	Overwrite     string // OVERWRITE {ON | OFF}; empty when omitted
	Encryption    string // ENCRYPTION {ON | OFF | SSE_S3 | SSE_KMS}; empty when omitted
	EncryptionKey string // ENCRYPTION SSE_KMS 'cmk_id'
}

func (c *IntoClause) nodeTag() {}

// ForUpdate represents FOR UPDATE / FOR SHARE and the wait modifier after it.
//
// NoWait, HasWait and SkipLocked are the three wait modifiers, of which a
// locking clause carries at most one. HasWait, not a zero WaitSeconds, is
// what says WAIT was written, so "FOR UPDATE WAIT 0" is not mistaken for no
// WAIT at all. WaitSeconds keeps the literal's own text, since OceanBase
// takes any numeric literal there.
type ForUpdate struct {
	Loc             Loc
	Share           bool        // set with LockInShareMode: OceanBase has no FOR SHARE
	LockInShareMode bool        // LOCK IN SHARE MODE (legacy syntax)
	Tables          []*TableRef // never set by the parser: OceanBase has no FOR UPDATE OF
	NoWait          bool
	HasWait         bool   // WAIT n present
	WaitSeconds     string // the literal after WAIT, as written; only when HasWait
	SkipLocked      bool
}

func (f *ForUpdate) nodeTag() {}

// VariableRef represents a variable reference (@var or @@var).
type VariableRef struct {
	Loc    Loc
	Name   string
	System bool   // @@ (system variable) vs @ (user variable)
	Scope  string // GLOBAL, SESSION, LOCAL for system variables
	Bare   bool   // a bare name, the stored-program variable of SELECT ... INTO
}

func (v *VariableRef) nodeTag()  {}
func (v *VariableRef) exprNode() {}

// -----------------------------------------------------------------------
// Additional statement nodes
// -----------------------------------------------------------------------

// SetStmt represents a SET statement.
type SetStmt struct {
	Loc         Loc
	Scope       string // the first assignment's scope; each assignment holds its own
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
	// Catalog is database_with_catalog_factor's first name when the statement
	// writes one: `USE catalog.db`.
	Catalog string
}

func (s *UseStmt) nodeTag()  {}
func (s *UseStmt) stmtNode() {}

// ExplainStmt represents an EXPLAIN/DESCRIBE statement.
//
// Which fields a statement can set is OceanBase's modifier list: EXTENDED and
// PARTITIONS are still here, as in MySQL 5.7, and there is neither an EXPLAIN
// ANALYZE nor a FOR CONNECTION, so no statement sets Analyze or ForConnection
// (oceanbase/parser/explain_forms_test.go).
type ExplainStmt struct {
	Loc           Loc
	Format        string // TRADITIONAL or JSON, bare
	Analyze       bool   // EXPLAIN ANALYZE: no OceanBase syntax sets this
	Extended      bool   // EXPLAIN EXTENDED
	Partitions    bool   // EXPLAIN PARTITIONS
	ForConnection int64  // EXPLAIN FOR CONNECTION: no OceanBase syntax sets this
	// Type is the plan type word: BASIC, OUTLINE, EXTENDED, EXTENDED_NOADDR,
	// PLANREGRESS or PARTITIONS. Extended and Partitions repeat two of them.
	Type string
	// Display is PRETTY or PRETTY_COLOR.
	Display string
	// Into is the plan table of EXPLAIN INTO, and StatementID the literal of
	// SET STATEMENT_ID = literal.
	Into        string
	StatementID ExprNode
	Stmt        StmtNode
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
	AsUser       string          // AS user, MySQL 8.0's; never set, OceanBase has no such clause
	WithRoleType string          // AS user's WITH ROLE; never set
	WithRoles    []string        // AS user's WITH ROLE names; never set
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
	// AllExceptAll and AllExceptAllPrivileges are set_role_clause's
	// `ALL EXCEPT grant_privileges`, whose exemption is every privilege rather
	// than a role.
	AllExceptAll           bool
	AllExceptAllPrivileges bool
	Roles                  []string // specific roles
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
	Body            Node   // parsed sp_proc_stmt (compound or simple). nil for loadable UDFs.
	BodyText        string // raw source bytes of the body (byte-preserved from the input segment).
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
	Schema      string // database the trigger's name was qualified with
	Name        string
	Timing      string // BEFORE, AFTER
	Event       string // INSERT, UPDATE, DELETE
	Table       *TableRef
	Order       *TriggerOrder
	Body        Node   // parsed sp_proc_stmt
	BodyText    string // raw source bytes of the body
}

func (s *CreateTriggerStmt) nodeTag()  {}
func (s *CreateTriggerStmt) stmtNode() {}

// TriggerOrder represents FOLLOWS/PRECEDES in trigger definition.
type TriggerOrder struct {
	Loc         Loc
	Follows     bool
	Schema      string // database the other trigger's name was qualified with
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
	Body         Node   // parsed sp_proc_stmt
	BodyText     string // raw source bytes of the body
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
//
// Columns holds the column-and-user-variable list: each element is a
// *ColumnRef for a column name and a *VariableRef for a user variable, so
// `(a, @v)` names a column and a variable rather than two columns, one of
// them called "@v".
type LoadDataStmt struct {
	Loc         Loc
	IsXML       bool // true for LOAD XML, false for LOAD DATA
	LowPriority bool // LOW_PRIORITY modifier
	Concurrent  bool // CONCURRENT modifier
	Local       bool
	Infile      string // INFILE 'file_name' (empty for the Aurora S3 form)
	FromS3      bool   // Aurora MySQL: LOAD DATA [FROM] S3 ... / LOAD XML FROM S3 ...
	S3Kind      string // Aurora MySQL: FILE | PREFIX | MANIFEST (empty when omitted; FILE is the server default)
	S3URI       string // Aurora MySQL: the 's3[-region]://bucket/key' literal after S3 [FILE | PREFIX | MANIFEST]
	// Source is OceanBase's LOAD DATA ... FROM url_spec: a location the rows come
	// from, instead of an INFILE. A url_spec that is a query over one is in
	// SourceQuery.
	Source      *ExternalTableRef
	SourceQuery *SelectStmt
	// OnError is opt_on_error, as written: "LOG ERRORS",
	// "LOG ERRORS REJECT LIMIT 1" or "LOG ERRORS REJECT LIMIT UNLIMITED".
	OnError string
	// RemoteOSS is opt_load_local's REMOTE_OSS: the file is read from an OSS
	// endpoint rather than the server's own filesystem.
	RemoteOSS  bool
	Replace    bool
	Ignore     bool
	Table      *TableRef
	Partitions []string // PARTITION (p0, p1, ...)
	// Compression is opt_compression's compression_name.
	Compression  string
	CharacterSet string // CHARACTER SET charset_name
	// IgnoreGenerated is opt_load_ignore_rows' GENERATED form: the server skips
	// n generated rows instead of n leading ones.
	IgnoreGenerated bool
	// The load_data_extended_option_list, which OceanBase adds after SET: the
	// error log file, the rejected-rows limit and the bad-rows file. A limit of
	// 0 with HasRejectLimit set is REJECT LIMIT 0, and an empty path with
	// HasLogFile set is LOGFILE ''.
	LogFile        string
	HasLogFile     bool
	HasRejectLimit bool
	RejectLimit    int
	BadFile        string
	HasBadFile     bool
	// ExtendedOptions holds load_data_extended_option_list in the order it was
	// written. A statement may repeat an option, and the list is what keeps
	// every occurrence; the scalar fields above keep the last one.
	ExtendedOptions    []*LoadDataExtendedOption
	RowsIdentifiedBy   string // ROWS IDENTIFIED BY '<tagname>' (XML only)
	Columns            []ExprNode
	SetList            []*Assignment
	LinesStartingBy    string
	LinesTerminatedBy  string
	FieldsTerminatedBy string
	FieldsEnclosedBy   string
	FieldsOptionalEncl bool // OPTIONALLY ENCLOSED BY
	FieldsEscapedBy    string
	IgnoreRows         int
}

// LoadDataExtendedOption is one load_data_extended_option of LOAD DATA, in the
// order written: LOGFILE [=] 'file', BADFILE [=] 'file', or REJECT LIMIT n.
//
// Ref: load_data_extended_option in sql_parser_mysql_mode.y (v4.3.5_CE_BP6,
// lines 5455-5472).
type LoadDataExtendedOption struct {
	Loc   Loc
	Kind  string // "LOGFILE", "BADFILE" or "REJECT LIMIT"
	Value string // the path, or the limit's number as written
}

func (o *LoadDataExtendedOption) nodeTag() {}

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
	Loc Loc
	// NoWriteToBinlog is never set: OceanBase's ANALYZE takes neither LOCAL nor
	// NO_WRITE_TO_BINLOG. The field stays because Bytebase builds against it.
	NoWriteToBinlog  bool
	Tables           []*TableRef
	HistogramOp      string   // "UPDATE" or "DROP"
	HistogramColumns []string // column names
	Buckets          int      // WITH N BUCKETS
	// UsingData is never set: OceanBase's UPDATE HISTOGRAM takes no USING DATA.
	UsingData string
	// Statistics is "COMPUTE" or "ESTIMATE" for ANALYZE TABLE t [PARTITION (p)]
	// {COMPUTE | ESTIMATE} STATISTICS ..., whose clauses are parsed and not kept.
	Statistics string
}

func (s *AnalyzeTableStmt) nodeTag()  {}
func (s *AnalyzeTableStmt) stmtNode() {}

// RestorePointStmt is OceanBase's CREATE RESTORE POINT and DROP RESTORE POINT.
//
// Ref: create_restore_point_stmt and drop_restore_point_stmt in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 4906-4918).
type RestorePointStmt struct {
	Loc  Loc
	Drop bool   // DROP RESTORE POINT, where false is CREATE RESTORE POINT
	Name string // relation_name
}

func (s *RestorePointStmt) nodeTag()  {}
func (s *RestorePointStmt) stmtNode() {}

// CatalogStmt is OceanBase's CREATE EXTERNAL CATALOG, DROP CATALOG and
// SET CATALOG: the catalogs a session reads external tables through.
//
// Ref: create_catalog_stmt, drop_catalog_stmt and set_catalog_stmt in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 4926-4956).
type CatalogStmt struct {
	Loc    Loc
	Action string // "CREATE", "DROP" or "SET"
	Name   string // relation_name
	// Properties is CREATE EXTERNAL CATALOG's external_properties_list.
	Properties  []*ExternalProperty
	IfNotExists bool // CREATE ... IF NOT EXISTS
	IfExists    bool // DROP ... IF EXISTS
}

func (s *CatalogStmt) nodeTag()  {}
func (s *CatalogStmt) stmtNode() {}

// DatabaseLinkStmt is OceanBase's CREATE DATABASE LINK and DROP DATABASE LINK.
//
// Ref: create_dblink_stmt and drop_dblink_stmt in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, lines 17690-17730).
type DatabaseLinkStmt struct {
	Loc      Loc
	Drop     bool   // DROP DATABASE LINK, where false is CREATE DATABASE LINK
	Name     string // the link's name
	User     string // CONNECT TO user
	Tenant   string // the @tenant after it, with its @
	Database string // DATABASE name
	// DatabaseCatalog is the catalog a `catalog.database` target names, apart
	// from the database's own name: `DATABASE `cat.db`` has none, and its
	// database is called cat.db.
	DatabaseCatalog string
	Password        string // IDENTIFIED BY 'password'
	Server          string // SERVER [=] 'host:port' or HOST 'host:port'
	Cluster         string // CLUSTER 'name', or ""
	IfNotExists     bool
	IfExists        bool
}

func (s *DatabaseLinkStmt) nodeTag()  {}
func (s *DatabaseLinkStmt) stmtNode() {}

// MethodOptStmt is method_opt, which the grammar reaches as a statement of its
// own (stmt: method_opt): the statistics method of DBMS_STATS, which the server
// parses and runs nothing for.
//
// Ref: method_opt, method_list, method, for_all, size_clause, for_columns and
// for_columns_item in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, lines
// 23802-23941).
type MethodOptStmt struct {
	Loc     Loc
	Methods []*StatsMethod
}

func (s *MethodOptStmt) nodeTag()  {}
func (s *MethodOptStmt) stmtNode() {}

// StatsMethod is one method of method_opt: FOR ALL [INDEXED | HIDDEN] COLUMNS
// [size_clause], or FOR COLUMNS with its for_columns_list.
type StatsMethod struct {
	Loc     Loc
	ForAll  bool                   // FOR ALL, where false is FOR COLUMNS
	Indexed string                 // "INDEXED" or "HIDDEN", and "" without one
	Size    string                 // FOR ALL's size_clause: "AUTO", "REPEAT", "SKEWONLY" or a number
	Items   []*StatsForColumnsItem // for_columns_list, for FOR COLUMNS
}

func (m *StatsMethod) nodeTag() {}

// StatsForColumnsItem is one for_columns_item: a column_clause, a size_clause,
// or a column with the size that follows it.
type StatsForColumnsItem struct {
	Loc     Loc
	Columns []string // column_clause: one column name, or an extension's '(' column_name_list ')'
	Size    string   // size_clause: "AUTO", "REPEAT", "SKEWONLY" or the number as written
}

func (i *StatsForColumnsItem) nodeTag() {}

// OptimizeTableStmt represents an OPTIMIZE TABLE statement, and OceanBase's
// OPTIMIZE TENANT beside it.
type OptimizeTableStmt struct {
	Loc             Loc
	NoWriteToBinlog bool
	Tables          []*TableRef
	// Tenant is optimize_stmt's TENANT form: the tenant's name, and "" for the
	// TABLE form, which is the one MySQL has.
	Tenant string
	// AllTenants is OPTIMIZE TENANT ALL, which is a mode of its own: a tenant
	// may be called ALL, and one written in backquotes is that tenant.
	AllTenants bool
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

// AlterSystemStmt is one of the cluster statements of OceanBase's
// alter_system_stmt: ALTER SYSTEM, ALTER DISKGROUP, SET ENCRYPTION and SET
// DECRYPTION. Action is the words that name the statement after ALTER SYSTEM,
// as the grammar writes them: "MAJOR FREEZE", "FLUSH", "SET" for a parameter,
// "DISKGROUP ADD DISK", "SET ENCRYPTION". The clauses are parsed and not kept.
type AlterSystemStmt struct {
	Loc    Loc
	Action string
}

func (s *AlterSystemStmt) nodeTag()  {}
func (s *AlterSystemStmt) stmtNode() {}

// ObjectDDLStmt is a CREATE, ALTER or DROP of an OceanBase object this parser
// keeps no tree for: a tenant, a resource unit or pool, a tablegroup, a
// sequence, an outline, a tenant snapshot, a concurrent limiting rule or a
// sensitive rule. Action is the statement and its object as the grammar writes
// them ("CREATE TENANT", "ALTER RESOURCE POOL", "CREATE OR REPLACE FORMAT
// OUTLINE"); Schema and Name are the object's name, Name empty where the
// statement names none (ALTER TENANT ALL, CREATE SNAPSHOT). IfExists is IF
// EXISTS or IF NOT EXISTS before the name. A sequence's options are kept in
// Options; every other object's options, and an outline's statements, are
// parsed and not kept.
type ObjectDDLStmt struct {
	Loc      Loc
	Action   string
	Schema   string
	Name     string
	IfExists bool
	Options  []*SequenceOption
}

func (s *ObjectDDLStmt) nodeTag()  {}
func (s *ObjectDDLStmt) stmtNode() {}

// SequenceOption is one option of CREATE or ALTER SEQUENCE, in written order.
// Name is the option's words: "INCREMENT BY", "START WITH", "MAXVALUE",
// "MINVALUE", "CACHE", "NOMAXVALUE", "NOMINVALUE", "CYCLE", "NOCYCLE",
// "NOCACHE", "ORDER", "NOORDER" or "RESTART". Value is the number of the first
// five as written, its sign joined to it ("-1", "1.5"), and empty for the rest.
//
// Ref: sequence_option and simple_num in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6).
type SequenceOption struct {
	Loc   Loc
	Name  string
	Value string
}

func (o *SequenceOption) nodeTag() {}

// AuditStmt is OceanBase's AUDIT or NOAUDIT statement. What it audits is
// parsed and not kept.
type AuditStmt struct {
	Loc     Loc
	NoAudit bool
}

func (s *AuditStmt) nodeTag()  {}
func (s *AuditStmt) stmtNode() {}

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

// DropViewStmt represents a DROP VIEW statement, and DROP MATERIALIZED VIEW
// when Materialized is set.
type DropViewStmt struct {
	Loc          Loc
	IfExists     bool
	Views        []*TableRef
	Cascade      bool
	Restrict     bool
	Materialized bool // DROP MATERIALIZED VIEW rather than DROP VIEW
}

func (s *DropViewStmt) nodeTag()  {}
func (s *DropViewStmt) stmtNode() {}

// DropMaterializedViewLogStmt represents DROP MATERIALIZED VIEW LOG ON <table>.
//
// Ref: drop_mlog_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 10926).
type DropMaterializedViewLogStmt struct {
	Loc   Loc
	Table *TableRef
}

func (s *DropMaterializedViewLogStmt) nodeTag()  {}
func (s *DropMaterializedViewLogStmt) stmtNode() {}

// CreateMaterializedViewLogStmt represents CREATE [OR REPLACE] MATERIALIZED
// VIEW LOG ON <table>.
//
// Ref: create_mlog_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 10724).
type CreateMaterializedViewLogStmt struct {
	Loc       Loc
	OrReplace bool
	Table     *TableRef
	Options   []*ParallelOption // opt_mlog_option_list
	With      *MlogWith         // nil when no WITH clause is written
	NewValues string            // INCLUDING or EXCLUDING; empty when absent
	Purge     *MlogPurge        // nil when no PURGE clause is written
}

func (s *CreateMaterializedViewLogStmt) nodeTag()  {}
func (s *CreateMaterializedViewLogStmt) stmtNode() {}

// MlogWith is the WITH clause of a materialized view log: any special columns
// (ROWID, PRIMARY KEY, SEQUENCE), then an optional parenthesized list of
// reference columns. HasReferenceList tells an absent list from WITH ().
type MlogWith struct {
	Loc              Loc
	SpecialColumns   []string // ROWID, PRIMARY KEY, SEQUENCE, in order
	ReferenceColumns []string
	HasReferenceList bool // the (...) after the special columns was written
}

func (w *MlogWith) nodeTag() {}

// MlogPurge is the PURGE clause of a materialized view log.
type MlogPurge struct {
	Loc       Loc
	Immediate bool   // PURGE IMMEDIATE
	Sync      string // SYNCHRONOUS or ASYNCHRONOUS; empty when absent
	StartWith ExprNode
	Next      ExprNode
}

func (p *MlogPurge) nodeTag() {}

// AlterMaterializedViewLogStmt represents ALTER MATERIALIZED VIEW LOG ON
// <table> with one or more actions.
//
// Ref: alter_mlog_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 10075).
type AlterMaterializedViewLogStmt struct {
	Loc     Loc
	Table   *TableRef
	Actions []*AlterMlogAction
}

func (s *AlterMaterializedViewLogStmt) nodeTag()  {}
func (s *AlterMaterializedViewLogStmt) stmtNode() {}

// AlterMlogActionType enumerates the actions alter_mlog_action takes.
type AlterMlogActionType int

const (
	AMLogParallel          AlterMlogActionType = iota // parallel_option
	AMLogPurge                                        // PURGE mv_refresh_interval
	AMLogLobInRowThreshold                            // LOB_INROW_THRESHOLD [=] n
)

// AlterMlogAction is one action of ALTER MATERIALIZED VIEW LOG.
type AlterMlogAction struct {
	Loc       Loc
	Type      AlterMlogActionType
	Parallel  *ParallelOption // AMLogParallel
	StartWith ExprNode        // AMLogPurge
	Next      ExprNode        // AMLogPurge
	Threshold int64           // AMLogLobInRowThreshold
}

func (a *AlterMlogAction) nodeTag() {}

// ParenExpr represents a parenthesized expression.
type ParenExpr struct {
	Loc  Loc
	Expr ExprNode
}

func (e *ParenExpr) nodeTag()  {}
func (e *ParenExpr) exprNode() {}

// ParamMarkerExpr represents a `?` parameter marker (paramMarker in the
// grammar), the placeholder a prepared statement binds a value to. It is a leaf
// with no value of its own. The parser accepts it only as a LIMIT/OFFSET
// operand; `?` in general expression position stays rejected.
type ParamMarkerExpr struct {
	Loc Loc
}

func (e *ParamMarkerExpr) nodeTag()  {}
func (e *ParamMarkerExpr) exprNode() {}

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

// KeywordArg represents a bare keyword in a function argument position: the
// temporal unit of TIMESTAMPDIFF/TIMESTAMPADD (timestampdiff(SECOND,a,b)) or
// the format type of GET_FORMAT (get_format(DATE,'USA')). MySQL lexes these
// positions as keywords, never identifiers — a quoted `SECOND` is a syntax
// error (oracle 8.0.32 + 5.7.25) — so deparse must render Keyword unquoted.
// Keyword is stored normalized: uppercase, and GET_FORMAT's TIMESTAMP folded
// to DATETIME, matching SHOW CREATE VIEW. OceanBase takes no SQL_TSI_ unit, so
// there is no synonym left to fold.
type KeywordArg struct {
	Loc     Loc
	Keyword string
}

func (e *KeywordArg) nodeTag()  {}
func (e *KeywordArg) exprNode() {}

// WeightStringExpr represents WEIGHT_STRING(expr [AS {CHAR|BINARY}(n)] [LEVEL ...]).
// OceanBase 4.3.5.6 stores every WEIGHT_STRING form as its internal five-argument
// call weight_string(expr, 0, nweights, flags, is_binary).
type WeightStringExpr struct {
	Loc      Loc
	Expr     ExprNode
	AsChar   *DataType // optional AS CHAR(n) or AS BINARY(n); nil when absent
	AsBinary bool      // true for AS BINARY(n)
	Levels   []WeightStringLevel
}

func (e *WeightStringExpr) nodeTag()  {}
func (e *WeightStringExpr) exprNode() {}

// WeightStringLevel is one element of a WEIGHT_STRING LEVEL list. The 5.7
// grammar allows a direction and REVERSE together, in that order — LEVEL 1
// DESC REVERSE stores as "level 1 desc reverse", while REVERSE DESC is a
// syntax error and ASC normalizes away (oracle 5.7.25).
type WeightStringLevel struct {
	Level   int
	Desc    bool
	Reverse bool
}

// CollateExpr represents a COLLATE expression.
type CollateExpr struct {
	Loc       Loc
	Expr      ExprNode
	Collation string
	// Quoted is set when the collation was written as a string literal, which
	// collation_name allows: `a COLLATE 'x y'` keeps its quotes in deparse.
	Quoted bool
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
	// CharsetQuoted marks a USING name written as a string ('utf8'). OceanBase reads it
	// as a parameter (charset_name's STRING_VALUE arm) and checks it after the argument.
	CharsetQuoted bool
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

// RbIterateExpr is RB_ITERATE(expr) [AS] alias [(column)], a table function of
// OceanBase's own: the rows of an array produced by a row-based expression.
//
// Ref: rb_iterate_expr in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6,
// line 25536).
type RbIterateExpr struct {
	Loc    Loc
	Arg    ExprNode // the simple_expr inside the parentheses
	Alias  string   // [AS] relation_name
	Column string   // the relation_name in parentheses after the alias, and "" without one
}

func (e *RbIterateExpr) nodeTag()   {}
func (e *RbIterateExpr) tableExpr() {}

// UnnestExpr is UNNEST(expr_list) [AS] alias [(column list)], which reads the
// rows of the arrays its arguments produce.
//
// Ref: unnest_expr in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, line
// 25559).
type UnnestExpr struct {
	Loc     Loc
	Args    []ExprNode // simple_expr_list
	Alias   string     // [AS] relation_name
	Columns []string   // relation_name_list in parentheses after the alias
}

func (e *UnnestExpr) nodeTag()   {}
func (e *UnnestExpr) tableExpr() {}

// TableCollectionExpr is tbl_name's `TABLE '(' simple_expr ')' [[AS] relation_name]`,
// which reads a collection expression as a table in FROM or DELETE FROM.
//
// Ref: tbl_name in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, lines
// 14476-14487).
type TableCollectionExpr struct {
	Loc   Loc
	Expr  ExprNode
	Alias string
}

func (e *TableCollectionExpr) nodeTag()   {}
func (e *TableCollectionExpr) tableExpr() {}

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
	Direction string // FIRST, NEXT, PREV, LAST
	Index     string // index name for index-based reads
	// Values is the expr_list of an index read with a comparison operator:
	// `HANDLER t READ i >= (1, 2)`, and Operator that operator as written.
	Values   []ExprNode
	Operator string
	Where    ExprNode // WHERE condition
	Limit    *Limit   // LIMIT clause
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

// HandlerCondKind classifies a handler condition value so that semantic
// analysis can distinguish SQLSTATE literals, error-code literals, built-in
// condition categories, and user-declared condition names. The prior `string`
// form collapsed all kinds into one field, making e.g. `SQLSTATE '23000'`
// indistinguishable from a user condition named `23000`.
type HandlerCondKind int

const (
	// HandlerCondSQLState — SQLSTATE '23000' / SQLSTATE VALUE '23000'.
	HandlerCondSQLState HandlerCondKind = iota
	// HandlerCondErrorCode — a bare integer mysql_error_code (e.g., 1062).
	HandlerCondErrorCode
	// HandlerCondSQLWarning — the built-in category.
	HandlerCondSQLWarning
	// HandlerCondNotFound — the built-in category.
	HandlerCondNotFound
	// HandlerCondSQLException — the built-in category.
	HandlerCondSQLException
	// HandlerCondName — a user-declared condition name (resolved via DECLARE
	// CONDITION in an enclosing scope).
	HandlerCondName
)

// HandlerCondValue is one handler condition value (one item in a DECLARE
// HANDLER FOR list, or the RHS of DECLARE CONDITION FOR).
type HandlerCondValue struct {
	Kind  HandlerCondKind
	Value string // SQLSTATE literal / error code / condition name; empty for SQLWARNING/NOT FOUND/SQLEXCEPTION
}

// DeclareConditionStmt represents a DECLARE ... CONDITION FOR statement.
type DeclareConditionStmt struct {
	Loc            Loc
	Name           string // condition name
	ConditionValue HandlerCondValue
}

func (s *DeclareConditionStmt) nodeTag()  {}
func (s *DeclareConditionStmt) stmtNode() {}

// DeclareHandlerStmt represents a DECLARE ... HANDLER FOR statement.
type DeclareHandlerStmt struct {
	Loc        Loc
	Action     string // CONTINUE, EXIT, or UNDO
	Conditions []HandlerCondValue
	Stmt       Node // handler statement body
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
//
// OceanBase has no TABLE query, so the parser never builds a TableStmt: "TABLE
// t" is 1064 there as a statement and wherever a query goes. The type stays so
// code written against the MySQL-shaped tree keeps compiling.
//
// A TableStmt that stands alone (a statement, an EXPLAIN target, the source
// of INSERT ... TABLE) holds its own ORDER BY, LIMIT, and INTO. As the
// TableSource of a SelectStmt (a subquery, a set-operation operand, a CTE or
// view body) it holds only the table, and the SelectStmt holds the clauses.
type TableStmt struct {
	Loc     Loc
	Table   *TableRef      // table name
	OrderBy []*OrderByItem // optional ORDER BY
	Limit   *Limit         // optional LIMIT/OFFSET
	Into    *IntoClause    // optional INTO OUTFILE/DUMPFILE/var
}

func (s *TableStmt) nodeTag()  {}
func (s *TableStmt) stmtNode() {}

// ValuesStmt represents a VALUES query primary: VALUES ROW(1,2), ROW(3,4).
//
// The parser builds one only as the ValuesSource of a SelectStmt, which holds its
// clauses: a VALUES statement, in any position, is a SelectStmt over a ValuesStmt.
// OrderBy and Limit are never set by the parser, and stay because Bytebase still
// references this node.
type ValuesStmt struct {
	Loc          Loc
	Rows         [][]ExprNode   // list of ROW(...) value lists
	LeadingComma bool           // true if any ROW(...) started with an empty leading comma (e.g. ROW(, a))
	OrderBy      []*OrderByItem // never set by the parser
	Limit        *Limit         // never set by the parser
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

// The resource group nodes below came from the MySQL tree this engine was
// forked from. OceanBase has no resource group statements at all: OceanBase 4.3.5-LTS in MySQL mode
// answers 1064 for CREATE, ALTER and DROP RESOURCE GROUP and for SET RESOURCE
// GROUP, so the parser reads none of them and produces none of these nodes.
// The types stay because callers name them in type switches.

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
	// TableID and HasTableID are opt_table_id's `TABLE_ID = n`, which pins the
	// view's id: the server parses it and refuses to run it.
	TableID    int64
	HasTableID bool
}

func (s *AlterViewStmt) nodeTag()  {}
func (s *AlterViewStmt) stmtNode() {}

// AlterMaterializedViewStmt represents ALTER MATERIALIZED VIEW with one or
// more actions.
//
// Ref: alter_mview_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 9985).
type AlterMaterializedViewStmt struct {
	Loc     Loc
	Name    *TableRef
	Actions []*AlterMviewAction
}

func (s *AlterMaterializedViewStmt) nodeTag()  {}
func (s *AlterMaterializedViewStmt) stmtNode() {}

// AlterMviewActionType enumerates the actions alter_mview_action takes.
type AlterMviewActionType int

const (
	AMViewEnableOnQueryComputation AlterMviewActionType = iota // [ENABLE|DISABLE] ON QUERY COMPUTATION
	AMViewEnableQueryRewrite                                   // [ENABLE|DISABLE] QUERY REWRITE
	AMViewParallel                                             // opt_set parallel_option
	AMViewRefreshMethod                                        // REFRESH {FAST|COMPLETE|FORCE}
	AMViewNeverRefresh                                         // NEVER REFRESH
	AMViewRefreshInterval                                      // REFRESH [START WITH e] [NEXT e]
	AMViewNestedRefresh                                        // REFRESH {INDIVIDUAL|INCONSISTENT|CONSISTENT}
)

// AlterMviewAction is one action of ALTER MATERIALIZED VIEW.
type AlterMviewAction struct {
	Loc      Loc
	Type     AlterMviewActionType
	Enable   bool            // the ENABLE/DISABLE switches
	Parallel *ParallelOption // AMViewParallel
	Refresh  *MviewRefreshInfo
}

func (a *AlterMviewAction) nodeTag() {}

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
	Body         Node   // parsed sp_proc_stmt (nil when no DO clause)
	BodyText     string // raw source bytes of the body
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

// PurgeStmt is OceanBase's purge_stmt, which empties an object out of the
// recycle bin. The binary-log form is PurgeBinaryLogsStmt beside it.
//
// Ref: purge_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 20190-20215).
//
//	PURGE {TABLE | INDEX} relation_factor
//	PURGE {DATABASE | SCHEMA} database_factor
//	PURGE TENANT relation_name
//	PURGE RECYCLEBIN
type PurgeStmt struct {
	Loc Loc
	// Kind is the word after PURGE: "TABLE", "INDEX", "DATABASE",
	// "TENANT" or "RECYCLEBIN".
	Kind string
	// Name is the object to purge. It is nil for PURGE RECYCLEBIN.
	Name *TableRef
}

func (s *PurgeStmt) nodeTag()  {}
func (s *PurgeStmt) stmtNode() {}

// FlashbackStmt is OceanBase's flashback_stmt, which restores an object from
// the recycle bin.
//
// Ref: flashback_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 20152-20188).
//
//	FLASHBACK TABLE relation_factor TO BEFORE DROP [RENAME TO relation_factor]
//	FLASHBACK {DATABASE | SCHEMA} database_factor TO BEFORE DROP [RENAME TO database_factor]
//	FLASHBACK TENANT relation_name TO BEFORE DROP [RENAME TO relation_name]
type FlashbackStmt struct {
	Loc Loc
	// Kind is the word after FLASHBACK: "TABLE", "DATABASE" or "TENANT".
	Kind string
	// Name is the object to restore.
	Name *TableRef
	// RenameTo is opt_flashback_rename_*'s RENAME TO name, or nil.
	RenameTo *TableRef
}

func (s *FlashbackStmt) nodeTag()  {}
func (s *FlashbackStmt) stmtNode() {}

// DumpMemoryStmt is OceanBase's dump_memory_stmt: the observer's memory
// diagnostics, which the server parses and answers to a client that asks for
// them.
//
// Ref: dump_memory_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 20236-20291).
type DumpMemoryStmt struct {
	Loc Loc
	// Kind names the alternative as written: "ENTITY_ALL", "ENTITY",
	// "CHUNK_ALL", "CHUNK", "CHUNK_NAME", "LEAK_MOD", "LEAK_RATE" or
	// "MEMORY_LEAK".
	Kind string
	// Args holds the alternative's values in the order it writes them, each as
	// written: the entity name and slot index of DUMP ENTITY, the tenant id and
	// context of DUMP CHUNK, the chunk name of DUMP CHUNK P_CHUNK, and the mode
	// or rate of SET OPTION.
	Args []string
}

func (s *DumpMemoryStmt) nodeTag()  {}
func (s *DumpMemoryStmt) stmtNode() {}

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
	Loc    Loc
	Tables []*LoadIndexTable // tbl_index_list, or the one table of the PARTITION form
	// CacheName is IN key_cache_name, with the catalog, database and database
	// link it may name beside it.
	CacheName    string
	CacheSchema  string
	CacheCatalog string
	CacheDblink  string
	// IgnoreLeaves is opt_ignore_leaves, which applies to the whole table list.
	IgnoreLeaves bool
}

func (s *CacheIndexStmt) nodeTag()  {}
func (s *CacheIndexStmt) stmtNode() {}

// LoadIndexIntoCacheStmt represents a LOAD INDEX INTO CACHE statement.
type LoadIndexIntoCacheStmt struct {
	Loc    Loc
	Tables []*LoadIndexTable
	// IgnoreLeaves is opt_ignore_leaves, which follows the whole table list and
	// so applies to every table of it.
	IgnoreLeaves bool
}

// LoadIndexTable represents one table of LOAD INDEX INTO CACHE, and one
// tbl_index of CACHE INDEX, which the same shapes serve.
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

// SetResourceGroupStmt represents a SET RESOURCE GROUP statement. This engine
// has no such statement and the parser produces none of these; see the note on
// VCPUSpec.
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
