package ast

import "strconv"

// -----------------------------------------------------------------------
// Statement nodes
// -----------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	Loc               Loc
	CTEs              []*CommonTableExpr // WITH clause (CTEs)
	DistinctKind      DistinctKind       // DISTINCT, ALL, or none
	CalcFoundRows     bool               // SQL_CALC_FOUND_ROWS
	Priority          Priority           // LOW_PRIORITY, HIGH_PRIORITY or DELAYED; the last one written
	StraightJoin      bool               // STRAIGHT_JOIN hint
	SmallResult       bool               // SQL_SMALL_RESULT
	BigResult         bool               // SQL_BIG_RESULT
	BufferResult      bool               // SQL_BUFFER_RESULT
	NoCache           bool               // SQL_NO_CACHE
	TargetList        []ExprNode         // select expressions (ResTarget)
	From              []TableExpr        // FROM clause (TableRef / JoinClause)
	Where             ExprNode           // WHERE condition
	GroupBy           []*OrderByItem     // GROUP BY items (TiDB's ByItem: expr with optional ASC/DESC)
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
	// TableSource and ValuesSource are the TABLE or VALUES query primary this
	// node stands for. The source holds only its table or rows; the ORDER BY,
	// LIMIT, locking clause, and INTO after it are this node's, except that an
	// ORDER BY and LIMIT after the last operand of a set operation are the set
	// operation's.
	TableSource  *TableStmt  // TABLE query primary
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
	Loc            Loc
	IsReplace      bool // REPLACE instead of INSERT
	Priority       Priority
	Ignore         bool          // INSERT IGNORE
	Table          *TableRef     // target table
	Partitions     []string      // PARTITION (p0, p1, ...)
	Columns        []*ColumnRef  // explicit column list; nil without one, non-nil and empty for ()
	Values         [][]ExprNode  // VALUES rows; a row of () has no elements
	Select         *SelectStmt   // INSERT ... SELECT, ... TABLE t, ... (query)
	SetList        []*Assignment // INSERT ... SET col=val (MySQL extension)
	RowAlias       string        // AS row_alias (MySQL 8.0.19+)
	ColAliases     []string      // (col_alias, ...) after row_alias
	OnDuplicateKey []*Assignment // ON DUPLICATE KEY UPDATE
	Returning      []Node        // not standard MySQL but included for completeness
}

func (s *InsertStmt) nodeTag()  {}
func (s *InsertStmt) stmtNode() {}

// Priority is the LOW_PRIORITY, HIGH_PRIORITY or DELAYED a SELECT, INSERT,
// REPLACE, UPDATE or DELETE carries: TiDB's Priority (parser.y:9046, v8.5.5),
// which every one of those statements reads.
type Priority int

const (
	PriorityNone Priority = iota
	PriorityLow
	PriorityDelayed
	PriorityHigh
)

// UpdateStmt represents an UPDATE statement.
type UpdateStmt struct {
	Loc      Loc
	Priority Priority
	Ignore   bool
	Tables   []TableExpr    // table references (multi-table update)
	SetList  []*Assignment  // SET clause
	Where    ExprNode       // WHERE condition
	OrderBy  []*OrderByItem // ORDER BY
	Limit    *Limit         // LIMIT
}

func (s *UpdateStmt) nodeTag()  {}
func (s *UpdateStmt) stmtNode() {}

// DeleteStmt represents a DELETE statement.
type DeleteStmt struct {
	Loc      Loc
	Priority Priority
	Quick    bool
	Ignore   bool
	Tables   []TableExpr    // tables to delete from
	Using    []TableExpr    // USING clause (multi-table delete)
	Where    ExprNode       // WHERE condition
	OrderBy  []*OrderByItem // ORDER BY
	Limit    *Limit         // LIMIT
}

func (s *DeleteStmt) nodeTag()  {}
func (s *DeleteStmt) stmtNode() {}

// BatchDryRunMode enumerates TiDB BATCH non-transactional DML dry-run modes.
// Numeric values mirror pingcap ast (dml.go: NoDryRun=0, DryRunQuery=1,
// DryRunSplitDml=2). DRY RUN → SplitDML; DRY RUN QUERY → Query.
type BatchDryRunMode int

const (
	BatchDryRunNone     BatchDryRunMode = iota // 0: execute the split DMLs
	BatchDryRunQuery                           // 1: DRY RUN QUERY — show the SELECT that splits
	BatchDryRunSplitDML                        // 2: DRY RUN — show the split DML jobs
)

// BatchStmt represents a TiDB non-transactional DML statement:
//
//	BATCH [ON <col>] LIMIT <n> [DRY RUN [QUERY]] {DELETE | UPDATE | INSERT | REPLACE}
//
// Ref: pingcap parser.y NonTransactionalDMLStmt (production "BATCH"
// OptionalShardColumn "LIMIT" NUM DryRunOptions ShardableStmt). REPLACE is
// surfaced as *InsertStmt with IsReplace=true (omni unifies INSERT/REPLACE).
type BatchStmt struct {
	Loc         Loc
	ShardColumn *ColumnRef // ON <col>; nil → TiDB auto-selects the handle column
	Limit       uint64     // LIMIT <n>; grammar is NUM (integer literal), not an expression
	DryRun      BatchDryRunMode
	DML         StmtNode // *DeleteStmt | *UpdateStmt | *InsertStmt
}

func (s *BatchStmt) nodeTag()  {}
func (s *BatchStmt) stmtNode() {}

// CreateTableStmt represents a CREATE TABLE statement.
type CreateTableStmt struct {
	Loc         Loc
	IfNotExists bool
	Temporary   bool
	// Global marks CREATE GLOBAL TEMPORARY TABLE, TiDB's own kind of temporary
	// table, and OnCommit holds the DELETE ROWS or PRESERVE ROWS that spelling
	// must end with.
	//
	// Temporary stays false for it, because a global temporary table is not in
	// the session's temporary namespace: v8.5.5 lists it in SHOW TABLES and in
	// information_schema.tables as a BASE TABLE, refuses a base table of the
	// same name with 1050, drops it with a plain DROP TABLE, and answers 1051
	// to DROP TEMPORARY TABLE for it. Only its rows are per-session.
	Global      bool
	OnCommit    string
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
	ATAlterColumnSetNotNull  // never built: TiDB's ALTER COLUMN has no SET NOT NULL
	ATAlterColumnDropNotNull // never built: TiDB's ALTER COLUMN has no DROP NOT NULL
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
	ATAlterColumnVisible   // never built: TiDB has no column visibility
	ATAlterColumnInvisible // never built: TiDB has no column visibility
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
	ATPartitionBy
	// TiDB-specific ALTER TABLE commands.
	ATSetTiFlashReplica // SET TIFLASH REPLICA n
	ATRemoveTTL         // REMOVE TTL
	ATCache             // CACHE
	ATNoCache           // NOCACHE
	ATCompact           // COMPACT [PARTITION names] [TIFLASH REPLICA]
	ATWriteable         // READ ONLY | READ WRITE
	ATAttributes        // ATTRIBUTES [=] {'str' | DEFAULT}
	ATStatsOptions      // STATS_OPTIONS [=] {'str' | DEFAULT}
	// ATPartitionAttributes and ATPartitionOptions are the two readings of
	// `PARTITION name ...`: an attributes clause, or a list of the options a
	// partition definition takes.
	ATPartitionAttributes // PARTITION name ATTRIBUTES [=] {'str' | DEFAULT}
	ATPartitionOptions    // PARTITION name partition_definition_option...
	// The four interval-partition actions, which name a partition by its
	// bound rather than by its name. LAST and FIRST are actions in the
	// comma-separated list; SPLIT MAXVALUE and MERGE FIRST close the
	// statement on their own.
	ATAddLastPartition         // LAST PARTITION LESS THAN (expr)
	ATDropFirstPartition       // FIRST PARTITION LESS THAN (expr)
	ATReorganizeLastPartition  // SPLIT MAXVALUE PARTITION LESS THAN (expr)
	ATReorganizeFirstPartition // MERGE FIRST PARTITION LESS THAN (expr)
	// Extended-statistics actions, TiDB's own, over one or more columns.
	ATAddStatistics  // ADD STATS_EXTENDED name type (col, ...)
	ATDropStatistics // DROP STATS_EXTENDED name
)

// AlterTableCmd represents a single ALTER TABLE operation.
type AlterTableCmd struct {
	OldColumnSchema string // CHANGE: the schema written before the old name
	OldColumnTable  string // CHANGE: the table written before the old name
	Loc             Loc
	Type            AlterTableCmdType
	Name            string       // column/constraint/index name
	NewName         string       // for RENAME operations
	Column          *ColumnDef   // for ADD/MODIFY/CHANGE COLUMN
	Columns         []*ColumnDef // for ADD (col1, col2, ...) multi-column form
	Constraint      *Constraint  // for ADD CONSTRAINT
	Option          *TableOption // for table options
	DefaultExpr     ExprNode     // for ALTER COLUMN SET DEFAULT expr
	IfExists        bool
	First           bool            // FIRST positioning
	After           string          // AFTER column positioning
	PartitionNames  []string        // for partition operations (DROP/TRUNCATE/ANALYZE/CHECK/OPTIMIZE/REBUILD/REPAIR/REORGANIZE/DISCARD/IMPORT)
	AllPartitions   bool            // for partition operations using ALL
	PartitionDefs   []*PartitionDef // for ADD PARTITION / REORGANIZE PARTITION ... INTO
	Number          int             // for COALESCE PARTITION number / ADD PARTITION PARTITIONS number
	// NoWriteToBinlog records MySQL's opt_no_write_to_binlog, the
	// NO_WRITE_TO_BINLOG or LOCAL modifier a partition-maintenance action
	// takes after the PARTITION keyword. The grammar folds both spellings to
	// one value, so the spelling used is not kept.
	NoWriteToBinlog bool
	// NewTable is the new name of a RENAME [TO | AS | =] action, which may
	// carry a database qualifier. NewName holds the same name as written,
	// qualifier and dot included, for the consumers that read a plain string.
	NewTable              *TableRef
	ExchangeTable         *TableRef        // for EXCHANGE PARTITION ... WITH TABLE
	WithValidation        *bool            // for EXCHANGE PARTITION: true=WITH VALIDATION, false=WITHOUT VALIDATION, nil=not specified
	OrderByItems          []*OrderByItem   // for ORDER BY operation
	PartitionBy           *PartitionClause // for PARTITION BY (repartition)
	TiFlashReplica        int              // TiDB: replica count for SET TIFLASH REPLICA
	TiFlashLocationLabels []string         // TiDB: optional LOCATION LABELS clause (empty when unset)
	// ForeignKey marks an ATDropConstraint written DROP FOREIGN KEY, which
	// shares the command type with DROP CHECK and DROP CONSTRAINT but not
	// their handling: with tidb_enable_check_constraint off, TiDB lets the
	// last two name nothing and answers 1091 for a missing foreign key.
	ForeignKey bool
	// IfNotExists records the IF NOT EXISTS the ADD actions TiDB and MariaDB
	// allow it on carry: ADD [COLUMN], ADD PARTITION, ADD {INDEX | KEY} and
	// ADD FOREIGN KEY. MySQL 8.0.46 has none of them.
	IfNotExists bool
	// Constraints holds the table-level constraints an
	// `ADD [COLUMN] (table_element, ...)` action carries beside Columns: the
	// paren form reads a table element list, not a column list, so
	// `ADD COLUMN (a INT, PRIMARY KEY (a))` is one action with both.
	Constraints []*Constraint
	// Attributes carries the value of an ATTRIBUTES or STATS_OPTIONS clause,
	// on the table (ATAttributes, ATStatsOptions) or on one partition
	// (ATPartitionAttributes). The command type says which clause was written.
	Attributes *AttributesSpec
	// Writeable is the state a READ ONLY or READ WRITE action asks for: true
	// for READ WRITE, false for READ ONLY.
	Writeable bool
	// PartitionOptions holds the partition-definition options a
	// `PARTITION name ...` action carries (ATPartitionOptions). The list is
	// empty for `ALTER TABLE t PARTITION p`, which TiDB accepts with no
	// options at all.
	PartitionOptions []*TableOption
	// CompactTiFlashReplica records the TIFLASH REPLICA that narrows a COMPACT
	// action to the TiFlash replicas; false compacts every replica kind.
	CompactTiFlashReplica bool
	// PartitionBound is the bound an interval-partition action names a
	// partition by: the expression inside LESS THAN (...).
	PartitionBound ExprNode
	// AnalyzeIndexFlag records the INDEX keyword of
	// `ANALYZE PARTITION names INDEX [index_names]`, which TiDB takes with an
	// empty name list, so the keyword and the names are kept apart.
	AnalyzeIndexFlag bool
	AnalyzeIndexes   []string
	// AnalyzeOptions is the WITH list, the same one ANALYZE TABLE takes.
	AnalyzeOptions []*AnalyzeOption
	// Hypo records the HYPO modifier of `SET HYPO TIFLASH REPLICA n`, which
	// pins a hypothetical TiFlash replica for the optimizer rather than a real
	// one. It is meaningful only when Type is ATSetTiFlashReplica.
	Hypo bool
	// StatsName, StatsType and StatsColumns carry the extended-statistics
	// actions: ADD STATS_EXTENDED (ATAddStatistics) uses all three, DROP
	// STATS_EXTENDED (ATDropStatistics) only the name.
	StatsName    string
	StatsType    StatisticsType
	StatsColumns []*ColumnRef
}

// AttributesSpec is the value an ATTRIBUTES or a STATS_OPTIONS clause carries.
// Both clauses take a string literal or the word DEFAULT, and nothing else:
// TiDB v8.5.5 answers 1064 for `ATTRIBUTES`, `ATTRIBUTES = 1` and
// `ATTRIBUTES 1`.
//
// Ref: AttributesOpt and StatsOptionsOpt in TiDB v8.5.5 pkg/parser/parser.y.
type AttributesSpec struct {
	Loc     Loc
	Default bool   // the clause was written as DEFAULT
	Value   string // the string literal, when Default is false
}

func (s *AttributesSpec) nodeTag() {}

func (c *AlterTableCmd) nodeTag() {}

// DropTableStmt represents a DROP TABLE statement.
type DropTableStmt struct {
	Loc       Loc
	IfExists  bool
	Tables    []*TableRef
	Temporary bool
	// Global marks DROP GLOBAL TEMPORARY TABLE. It is another way to name a
	// table CREATE GLOBAL TEMPORARY TABLE made, not the only one -- a plain
	// DROP TABLE drops it too -- and Temporary stays false for it, since the
	// table is not in the session's temporary namespace.
	Global   bool
	Cascade  bool
	Restrict bool
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
	Loc        Loc
	Expr       ExprNode // column name or expression
	Length     int      // prefix length
	Desc       bool     // DESC ordering
	Functional bool     // true when parsed from functional key-part syntax: (expr)
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

	// TiDB: when Name=="TIFLASH REPLICA", the numeric count and optional
	// label list live here (Value is ignored for that case).
	// Ref: parser.y:4482 DatabaseOption arm reuses LocationLabelList.
	TiFlashReplica        int
	TiFlashLocationLabels []string
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
	Schema                   string // the schema written before the name, in `d.t.a INT`
	Table                    string // the table written before the name, in `t.a INT`
	Name                     string
	TypeName                 *DataType
	Constraints              []*ColumnConstraint
	DefaultValue             ExprNode
	Comment                  string
	AutoIncrement            bool
	AutoRandom               bool     // TiDB: AUTO_RANDOM attribute
	AutoRandomShardBits      int      // TiDB: AUTO_RANDOM(shard_bits[, range])
	AutoRandomRangeBits      int      // TiDB: AUTO_RANDOM(shard_bits, range_bits)
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
	NotEnforced bool  // for CHECK ... NOT ENFORCED
	Clustered   *bool // TiDB: for PRIMARY KEY column constraint
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
	// ConstrVector is TiDB's VECTOR INDEX, which has no CONSTRAINT symbol
	// form: v8.5.5 answers 1064 for `ADD CONSTRAINT c VECTOR INDEX (...)`.
	ConstrVector
)

// Constraint represents a table-level constraint.
//
// Name is the constraint's name: the CONSTRAINT symbol of a FOREIGN KEY or
// CHECK, and the index name of a key, which for UNIQUE is the CONSTRAINT
// symbol when one is written and else the name after UNIQUE [KEY|INDEX]:
// TiDB names the key after the symbol (pkg/parser/parser.y, rule
// Constraint; verified on v8.5.5). IndexName is the name written after
// FOREIGN KEY, kept apart from the symbol: the foreign key and its index
// take the symbol when one is written, and this name only when none is
// (verified on v8.5.5).
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
	Clustered    *bool  // TiDB: nil=unset, true=CLUSTERED, false=NONCLUSTERED
	// IfNotExists records the IF NOT EXISTS that TiDB and MariaDB allow after
	// {INDEX | KEY}, after FOREIGN KEY and after VECTOR INDEX, and nowhere
	// else: TiDB v8.5.5 answers 1064 for it after UNIQUE, PRIMARY KEY,
	// FULLTEXT and CHECK.
	IfNotExists bool
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
	// Force records the FORCE that TiDB allows before AUTO_INCREMENT and
	// AUTO_RANDOM_BASE, and before no other option. It changes what the server
	// does -- the value is applied even when it moves the counter backwards --
	// so it is part of the option rather than noise before it.
	Force bool
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
	// Interval is RANGE's INTERVAL clause, from which TiDB generates the partitions.
	Interval *PartitionInterval
	// SystemTimeInterval and SystemTimeUnit are SYSTEM_TIME INTERVAL expr unit, and
	// HasSystemTimeLimit and SystemTimeLimit are SYSTEM_TIME LIMIT n.
	SystemTimeInterval ExprNode
	SystemTimeUnit     string
	HasSystemTimeLimit bool
	SystemTimeLimit    uint64
	// UpdateIndexes holds the UPDATE INDEXES clause TiDB allows after the
	// partition definitions, which says whether each named index stays global
	// or becomes local once the table is repartitioned.
	UpdateIndexes []*UpdateIndexSpec
}

// PartitionInterval is TiDB's INTERVAL (expr [unit]) clause on RANGE partitioning.
//
//	INTERVAL (expr [unit]) [FIRST PARTITION LESS THAN (expr) LAST PARTITION LESS THAN (expr)]
//	    [NULL PARTITION] [MAXVALUE PARTITION]
//
// Ref: PartitionIntervalOpt in TiDB v8.5.5 pkg/parser/parser.y.
type PartitionInterval struct {
	Loc        Loc
	Expr       ExprNode
	Unit       string   // normalized time unit, "" when not written
	First      ExprNode // FIRST PARTITION LESS THAN, nil when not written
	Last       ExprNode // LAST PARTITION LESS THAN, nil when not written
	NullPart   bool     // NULL PARTITION
	MaxValPart bool     // MAXVALUE PARTITION
}

func (n *PartitionInterval) nodeTag() {}

// UpdateIndexSpec is one entry of a TiDB UPDATE INDEXES clause. Global is true
// for GLOBAL and false for LOCAL; one of the two is required, so there is no
// unset state: TiDB v8.5.5 answers 1064 for `UPDATE INDEXES (i)`.
//
// Ref: UpdateIndexesList in TiDB v8.5.5 pkg/parser/parser.y.
type UpdateIndexSpec struct {
	Loc    Loc
	Name   string
	Global bool
}

func (s *UpdateIndexSpec) nodeTag() {}

func (p *PartitionClause) nodeTag() {}

// PartitionType enumerates partition types.
type PartitionType int

const (
	PartitionRange PartitionType = iota
	PartitionList
	PartitionHash
	PartitionKey
	PartitionSystemTime // SYSTEM_TIME, which TiDB reads and does not apply
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
	Loc    Loc
	Name   string // INT, VARCHAR, DECIMAL, etc.
	Length int    // type length (e.g., VARCHAR(255))
	Scale  int    // decimal scale (e.g., DECIMAL(10,2))
	// HasScale records that a scale was written, since FLOAT(30,0) is no FLOAT(30): Scale is 0 for both.
	HasScale   bool
	Unsigned   bool
	Zerofill   bool
	Charset    string   // CHARACTER SET
	Collate    string   // COLLATE
	Binary     bool     // standalone BINARY modifier, e.g. CHAR(10) BINARY
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
	Over       *WindowDef // OVER clause for window functions
	IgnoreNull bool       // IGNORE NULLS (LEAD/LAG/FIRST_VALUE/LAST_VALUE/NTH_VALUE), vs the default RESPECT NULLS
	FromLast   bool       // NTH_VALUE ... FROM LAST, vs the default FROM FIRST
	Separator  ExprNode   // GROUP_CONCAT SEPARATOR value (nil = not specified)
	Using      string     // CHAR(n, ... USING charset_name): the character set as written ("" = no USING)
}

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
	}
	return ""
}

// RemovesDuplicates reports whether the quantifier makes the aggregate read
// each distinct value once. ALL is the default and changes nothing; the
// DISTINCT ALL forms de-duplicate, which is what the engine's grammar calls a
// buggy distinct option (pingcap/tidb pkg/parser/parser.y,
// BuggyDefaultFalseDistinctOpt).
func (q AggQuantifier) RemovesDuplicates() bool {
	switch q {
	case AggQuantDistinct, AggQuantDistinctRow, AggQuantDistinctAll, AggQuantDistinctRowAll:
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
	ILike   bool // ILIKE, TiDB's case-insensitive LIKE
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
	// Array is CAST(expr AS type ARRAY), the multi-valued index form. It
	// belongs to the cast and not to DataType: no other position takes it,
	// and CONVERT(expr, type ARRAY) is a syntax error.
	Array bool
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
	Column string // column name
	Star   bool   // table.*
}

func (r *ColumnRef) nodeTag()  {}
func (r *ColumnRef) exprNode() {}

// TableRef represents a table reference (possibly qualified, with alias).
type TableRef struct {
	Loc        Loc
	Schema     string      // database name
	Name       string      // table name
	Alias      string      // AS alias
	Partitions []string    // PARTITION (p0, p1, ...)
	AsOf       *AsOfClause // AS OF TIMESTAMP clause; nil when there is none
	IndexHints []*IndexHint
	Sample     *TableSample // TABLESAMPLE clause; nil when there is none
}

func (r *TableRef) nodeTag()   {}
func (r *TableRef) tableExpr() {}

// AsOfClause is the stale-read clause of a table reference:
//
//	AS OF TIMESTAMP expr
//
// Loc runs from AS to the last token of the expression, and the TableRef's Loc
// includes it. Timestamp is never nil.
type AsOfClause struct {
	Loc       Loc
	Timestamp ExprNode
}

func (c *AsOfClause) nodeTag() {}

// TableSample is the TABLESAMPLE clause of a table reference:
//
//	TABLESAMPLE [SYSTEM | BERNOULLI | REGIONS] ([expr [ROWS | PERCENT]]) [REPEATABLE (expr)]
//
// Loc runs from TABLESAMPLE to the closing parenthesis of the last group, and
// the TableRef's Loc includes it.
type TableSample struct {
	Loc        Loc
	Method     TableSampleMethod
	Size       ExprNode        // nil for empty parentheses, which is not the same as a size of 0
	Unit       TableSampleUnit // SampleUnitDefault when the size names no unit, or there is no size
	Repeatable ExprNode        // the REPEATABLE seed; nil when the clause has none
}

func (s *TableSample) nodeTag() {}

// TableSampleMethod enumerates the sampling methods TABLESAMPLE names.
type TableSampleMethod int

const (
	SampleMethodDefault   TableSampleMethod = iota // no method written
	SampleMethodSystem                             // SYSTEM
	SampleMethodBernoulli                          // BERNOULLI
	SampleMethodRegions                            // REGIONS
)

// TableSampleUnit enumerates the units a TABLESAMPLE size names.
type TableSampleUnit int

const (
	SampleUnitDefault TableSampleUnit = iota // no unit written
	SampleUnitRows                           // ROWS
	SampleUnitPercent                        // PERCENT
)

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
	// Concatenated marks a literal folded from adjacent quoted strings ('a' 'b'
	// is "ab"), TiDB's StringLiteral (parser.y:7660, v8.5.5). Value holds the
	// folded result and FirstSegment the first string's value, because TiDB
	// names the column after the first string only: SELECT 'a' 'b' is column a,
	// value ab.
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
	Loc     Loc
	Value   string
	Charset string // optional charset introducer, as in _utf8 0x61 ("" = none)
}

func (l *HexLit) nodeTag()  {}
func (l *HexLit) exprNode() {}

// BitLit represents a bit literal (0b101 or b'101').
type BitLit struct {
	Loc     Loc
	Value   string
	Charset string // optional charset introducer, as in _utf8 0b1 ("" = none)
}

func (l *BitLit) nodeTag()  {}
func (l *BitLit) exprNode() {}

// TemporalLit represents a temporal literal: DATE '...', TIME '...', TIMESTAMP
// '...', or the ODBC escape {d ...}, {t ...} or {ts ...} that means the same.
type TemporalLit struct {
	Loc   Loc
	Type  string // "DATE", "TIME", "TIMESTAMP"
	Value string
	// Expr is the argument of an ODBC escape that is not a string, as in
	// {d 123}: TiDB parses it and refuses it only when it evaluates it. Value
	// is then empty.
	Expr ExprNode
}

func (l *TemporalLit) nodeTag()  {}
func (l *TemporalLit) exprNode() {}

// -----------------------------------------------------------------------
// Clause and helper nodes
// -----------------------------------------------------------------------

// Assignment represents a SET assignment (col = value).
type Assignment struct {
	Loc    Loc
	Scope  string // SET only: the GLOBAL, SESSION or LOCAL written on this assignment
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
//
// NoWait, HasWait and SkipLocked are the three wait modifiers, of which a
// locking clause carries at most one. HasWait, not a zero WaitSeconds, is
// what says WAIT was written, so "FOR UPDATE WAIT 0" is not mistaken for no
// WAIT at all. WaitSeconds keeps the literal's own text, since what each
// engine accepts there differs (TiDB takes an integer that fits a uint64;
// MariaDB and OceanBase take any numeric literal).
type ForUpdate struct {
	Loc             Loc
	Share           bool // FOR SHARE instead of FOR UPDATE
	LockInShareMode bool // LOCK IN SHARE MODE (legacy syntax)
	Tables          []*TableRef
	NoWait          bool
	HasWait         bool   // WAIT n present (TiDB's own modifier; MySQL has none)
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
}

func (v *VariableRef) nodeTag()  {}
func (v *VariableRef) exprNode() {}

// -----------------------------------------------------------------------
// Additional statement nodes
// -----------------------------------------------------------------------

// SetStmt represents a SET statement.
type SetStmt struct {
	Loc         Loc
	Scope       string // the first assignment's scope; each one holds its own, and none carries to the next
	Assignments []*Assignment
	// SessionStates holds the string of SET SESSION_STATES '...', a standalone
	// form with no assignments; HasSessionStates tells it from SET with none.
	SessionStates    string
	HasSessionStates bool
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

	// TiDB's own SHOW statements. The partition list of SHOW TABLE ...
	// REGIONS and SHOW TABLE ... DISTRIBUTIONS rides on From.Partitions,
	// where TiDB's own parser keeps it (ShowStmt.Table.PartitionNames in
	// parser.y:11818); Partition below is the single name of SHOW PLACEMENT
	// FOR TABLE ... PARTITION, which TiDB keeps apart for the same reason
	// (ShowStmt.Partition, parser.y:11977).
	IndexName         string   // INDEX name of SHOW TABLE ... INDEX idx REGIONS
	Partition         string   // PARTITION name of SHOW PLACEMENT FOR TABLE t PARTITION p
	PolicyName        string   // policy name of SHOW CREATE PLACEMENT POLICY
	ResourceGroupName string   // group name of SHOW CREATE RESOURCE GROUP
	JobID             ExprNode // job id of SHOW DISTRIBUTION JOB n
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
//
// Which fields a statement can set is TiDB's modifier list, not MySQL's:
// Extended and Partitions are MySQL 5.7 modifiers TiDB never had, so no
// statement sets either (tidb/parser/explain_forms_test.go). Stmt is nil for
// the FOR CONNECTION form, which is how a consumer tells it from a connection
// id that happens to be 0.
type ExplainStmt struct {
	Loc           Loc
	Format        string // TRADITIONAL, JSON, ROW, DOT, BRIEF, VERBOSE, TRUE_CARD_COST, TIDB_JSON, or any quoted name
	Analyze       bool   // EXPLAIN ANALYZE, which runs the statement it carries
	Extended      bool   // EXPLAIN EXTENDED: no TiDB syntax sets this
	Partitions    bool   // EXPLAIN PARTITIONS: no TiDB syntax sets this
	ForConnection int64  // EXPLAIN FOR CONNECTION connection_id
	Stmt          StmtNode
}

func (s *ExplainStmt) nodeTag()  {}
func (s *ExplainStmt) stmtNode() {}

// TraceStmt represents TRACE, TiDB's execution and optimizer tracer.
//
// TRACE wraps a statement and reports how the server ran it. The wrapped
// statement is the engine's TraceableStmt, which holds DML, queries, LOAD
// DATA, ANALYZE TABLE, the transaction statements and SET -- and no DDL at
// all. TRACE ALTER TABLE t ADD COLUMN c is a 1064 on v8.5.5, which is why
// this node carries no DDL arm (tidb/parser/trace_forms_test.go measures the
// whole traceable set and the DDL that is not in it).
//
// TracePlan says whether the statement runs. Plain TRACE and TRACE FORMAT =
// 'json' run it -- TRACE SET foreign_key_checks = 0 leaves the session
// variable at 0 -- while TRACE PLAN and TRACE PLAN TARGET = 'estimation'
// trace the optimizer and run nothing, leaving it at 1. Both were measured on
// v8.5.5 and are what the catalog and review's walk follow.
//
// Format is the name FORMAT = gave, empty when the statement names none. The
// engine's own node defaults it to "row" instead; this one keeps ExplainStmt's
// convention, so a consumer can tell TRACE from TRACE FORMAT = 'row'.
type TraceStmt struct {
	Loc Loc
	// Format is FORMAT = 'name'. It is a string literal in the grammar and
	// nothing else: TRACE FORMAT = row SELECT 1 is a 1064, where EXPLAIN
	// takes the bare name.
	Format string
	// TracePlan is TRACE PLAN, which traces the optimizer rather than
	// running the statement.
	TracePlan bool
	// Target is TRACE PLAN TARGET = 'name'. The engine checks the name when
	// the statement runs, not when it parses, so any text parses --
	// TARGET = 'arandomstring' included.
	Target string
	Stmt   StmtNode
}

func (s *TraceStmt) nodeTag()  {}
func (s *TraceStmt) stmtNode() {}

// PlanReplayerStmt represents PLAN REPLAYER, TiDB's plan-reproduction dump.
//
// The statement has four shapes, and the flags say which one:
//
//   - DUMP EXPLAIN [ANALYZE] <stmt>, which carries Stmt
//   - DUMP EXPLAIN [ANALYZE] SLOW QUERY ..., which sets SlowQuery
//   - DUMP EXPLAIN [ANALYZE] 'file' and LOAD 'file', which carry File
//   - CAPTURE [REMOVE] 'sqlDigest' 'planDigest', which carry both digests
//
// Analyze is the one flag with an effect outside the dump: like EXPLAIN's, it
// runs the statement it carries, and the statement it carries may be DDL.
// PLAN REPLAYER DUMP EXPLAIN ANALYZE ALTER TABLE t ADD COLUMN c adds the
// column on v8.5.5, while PLAN REPLAYER DUMP EXPLAIN ALTER TABLE t ADD COLUMN
// c writes a dump and leaves the table alone (both measured; pinned by
// tidb/parser/plan_replayer_forms_test.go and applied by
// tidb/catalog/exec.go).
type PlanReplayerStmt struct {
	Loc Loc
	// Analyze is the ANALYZE after EXPLAIN, which runs the statement.
	Analyze bool
	// SlowQuery is the SLOW QUERY form, which dumps the slow-query log
	// rather than one statement. It is a flag of its own because every
	// other field of that form is optional: PLAN REPLAYER DUMP EXPLAIN
	// SLOW QUERY sets nothing else.
	SlowQuery bool
	// Load is PLAN REPLAYER LOAD 'file', which reads a dump back.
	Load bool
	// Capture is PLAN REPLAYER CAPTURE, and Remove its REMOVE form. Both
	// carry SQLDigest and PlanDigest.
	Capture bool
	Remove  bool
	// File is the file LOAD reads, or the one the DUMP EXPLAIN 'file' form
	// reads statements from.
	File       string
	SQLDigest  string
	PlanDigest string
	// HistoricalStats is the expression in DUMP WITH STATS AS OF TIMESTAMP
	// <expr>, which dumps the statistics as of that time. WITH STATS is not
	// optional in itself: PLAN REPLAYER DUMP WITH STATS EXPLAIN SELECT 1 is
	// a 1064.
	HistoricalStats ExprNode
	// Stmt is the dumped statement, nil for every form that carries none.
	Stmt    StmtNode
	Where   ExprNode
	OrderBy []*OrderByItem
	Limit   *Limit
}

func (s *PlanReplayerStmt) nodeTag()  {}
func (s *PlanReplayerStmt) stmtNode() {}

// BeginStmt represents BEGIN / START TRANSACTION.
type BeginStmt struct {
	Loc Loc
	// Mode is the transaction mode BEGIN names: "PESSIMISTIC", "OPTIMISTIC",
	// or "" for none.
	Mode                   string
	ReadOnly               bool
	ReadWrite              bool
	WithConsistentSnapshot bool
	// CausalConsistencyOnly is START TRANSACTION WITH CAUSAL CONSISTENCY ONLY.
	CausalConsistencyOnly bool
	// AsOf is the expression of START TRANSACTION READ ONLY AS OF TIMESTAMP
	// expr, the time a stale read reads at; nil when there is none.
	AsOf ExprNode
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

// LockStatsStmt represents a LOCK STATS statement, which freezes the
// statistics of the tables or partitions it names. TiDB has it and MySQL has
// not.
//
// Partitions is set only by the arm that names one table and its partitions;
// the arm that names a list of tables leaves it empty.
type LockStatsStmt struct {
	Loc        Loc
	Tables     []*TableRef
	Partitions []string
}

func (s *LockStatsStmt) nodeTag()  {}
func (s *LockStatsStmt) stmtNode() {}

// UnlockStatsStmt represents an UNLOCK STATS statement, the other half of
// LockStatsStmt, which takes the same targets.
type UnlockStatsStmt struct {
	Loc        Loc
	Tables     []*TableRef
	Partitions []string
}

func (s *UnlockStatsStmt) nodeTag()  {}
func (s *UnlockStatsStmt) stmtNode() {}

// RequireClause represents a REQUIRE clause for TLS options.
type RequireClause struct {
	Loc     Loc
	None    bool   // REQUIRE NONE
	SSL     bool   // REQUIRE SSL
	X509    bool   // REQUIRE X509
	Cipher  string // REQUIRE CIPHER 'cipher'
	Issuer  string // REQUIRE ISSUER 'issuer'
	Subject string // REQUIRE SUBJECT 'subject'
	// SAN and TokenIssuer are TiDB's own REQUIRE options; MySQL 8.0 has
	// neither.
	SAN         string // REQUIRE SAN 'san'
	TokenIssuer string // REQUIRE TOKEN_ISSUER 'issuer'
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
	// RESOURCE GROUP group_name, the last clause of the statement
	// (ResourceGroupNameOption, TiDB v8.5.5 parser.y:14213).
	ResourceGroup string
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
	// RESOURCE GROUP group_name, the last clause of the statement
	// (ResourceGroupNameOption, TiDB v8.5.5 parser.y:14213).
	ResourceGroup string
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
	Body            Node   // parsed sp_proc_stmt (compound or simple); nil for loadable UDFs.
	BodyText        string // body text as tokenized (the lexer expands conditional comments).
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
	Body        Node   // parsed sp_proc_stmt (compound or simple).
	BodyText    string // body text as tokenized (the lexer expands conditional comments).
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
	Body         Node   // parsed sp_proc_stmt (compound or simple); nil when no DO body.
	BodyText     string // body text as tokenized (the lexer expands conditional comments).
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
	Loc                Loc
	IsXML              bool // true for LOAD XML, false for LOAD DATA
	LowPriority        bool // LOW_PRIORITY modifier
	Concurrent         bool // CONCURRENT modifier
	Local              bool
	Infile             string
	Format             string // TiDB: FORMAT 'delimited data' | 'sql file' | 'parquet'
	HasFormat          bool   // FORMAT written, which FORMAT '' does not show
	Replace            bool
	Ignore             bool
	Table              *TableRef
	Partitions         []string // PARTITION (p0, p1, ...)
	CharacterSet       string   // CHARACTER SET charset_name
	RowsIdentifiedBy   string   // ROWS IDENTIFIED BY '<tagname>' (XML only)
	Columns            []ExprNode
	SetList            []*Assignment
	LinesStartingBy    string
	LinesTerminatedBy  string
	FieldsTerminatedBy string
	FieldsEnclosedBy   string
	FieldsOptionalEncl bool // OPTIONALLY ENCLOSED BY
	FieldsEscapedBy    string
	// TiDB: FIELDS DEFINED NULL BY 'x' [OPTIONALLY ENCLOSED]. HasDefinedNullBy
	// is what tells DEFINED NULL BY '' from no such item.
	FieldsDefinedNullBy    string
	HasDefinedNullBy       bool
	DefinedNullOptEnclosed bool
	IgnoreRows             int
	Options                []*LoadDataOption // TiDB: WITH detached, threads=10
}

func (s *LoadDataStmt) nodeTag()  {}
func (s *LoadDataStmt) stmtNode() {}

// LoadDataOption is one entry of TiDB's WITH list on LOAD DATA and IMPORT
// INTO. Value is nil for a bare name such as `detached`.
type LoadDataOption struct {
	Loc   Loc
	Name  string
	Value ExprNode
}

func (o *LoadDataOption) nodeTag() {}

// ImportIntoStmt represents TiDB's IMPORT INTO statement. Exactly one of
// Path, Select and TableSource carries the source the rows come from: a file
// path, a query, or the TABLE statement `FROM TABLE s` is. Columns holds the
// same column-and-user-variable list LoadDataStmt documents.
//
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-import-into/
type ImportIntoStmt struct {
	Loc     Loc
	Table   *TableRef
	Columns []ExprNode
	SetList []*Assignment
	Path    string
	// Select is the query source: FROM select_statement, and also FROM
	// TABLE / FROM VALUES, which are query primaries the engine reads as a
	// SelectStmt with its own TableSource or ValuesSource.
	Select    *SelectStmt
	Format    string // FORMAT 'csv' | 'sql file' | 'parquet'
	HasFormat bool   // FORMAT written, which FORMAT '' does not show
	Options   []*LoadDataOption
}

func (s *ImportIntoStmt) nodeTag()  {}
func (s *ImportIntoStmt) stmtNode() {}

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

// AnalyzeColumnChoice names which columns an ANALYZE TABLE statement collects
// statistics on, which TiDB spells four ways: nothing at all, ALL COLUMNS,
// PREDICATE COLUMNS, or COLUMNS with a list.
type AnalyzeColumnChoice int

const (
	// AnalyzeColumnsUnset is a statement that named no column choice.
	AnalyzeColumnsUnset AnalyzeColumnChoice = iota
	// AnalyzeAllColumns is ALL COLUMNS.
	AnalyzeAllColumns
	// AnalyzePredicateColumns is PREDICATE COLUMNS.
	AnalyzePredicateColumns
	// AnalyzeColumnList is COLUMNS with the list in AnalyzeTableStmt.Columns.
	AnalyzeColumnList
)

// AnalyzeOptionType names the unit of one option in ANALYZE TABLE's WITH
// list. The list is closed: TiDB's AnalyzeOptionType has these six and no
// other, so a word outside it is a syntax error rather than an option whose
// value is ignored.
//
// Ref: AnalyzeOption in pingcap/tidb pkg/parser/parser.y.
type AnalyzeOptionType int

const (
	// AnalyzeOptBuckets is `n BUCKETS`.
	AnalyzeOptBuckets AnalyzeOptionType = iota
	// AnalyzeOptTopN is `n TOPN`.
	AnalyzeOptTopN
	// AnalyzeOptCMSketchWidth is `n CMSKETCH WIDTH`.
	AnalyzeOptCMSketchWidth
	// AnalyzeOptCMSketchDepth is `n CMSKETCH DEPTH`.
	AnalyzeOptCMSketchDepth
	// AnalyzeOptSamples is `n SAMPLES`.
	AnalyzeOptSamples
	// AnalyzeOptSampleRate is `n SAMPLERATE`, the one unit that takes a
	// number rather than an integer.
	AnalyzeOptSampleRate
)

// String names the option as TiDB spells it, so a tree dump and a message
// tell CMSKETCH WIDTH from CMSKETCH DEPTH.
func (t AnalyzeOptionType) String() string {
	switch t {
	case AnalyzeOptBuckets:
		return "BUCKETS"
	case AnalyzeOptTopN:
		return "TOPN"
	case AnalyzeOptCMSketchWidth:
		return "CMSKETCH WIDTH"
	case AnalyzeOptCMSketchDepth:
		return "CMSKETCH DEPTH"
	case AnalyzeOptSamples:
		return "SAMPLES"
	case AnalyzeOptSampleRate:
		return "SAMPLERATE"
	}
	return "UNKNOWN"
}

// AnalyzeOption is one `n UNIT` entry of ANALYZE TABLE's WITH list.
//
// Value is the number as the statement writes it. Every unit but SAMPLERATE
// takes a decimal integer that fits in uint64 -- v8.5.5 answers 1064 for
// `WITH 1.5 BUCKETS`, `WITH 0x10 BUCKETS` and for an integer past that range
// -- so a consumer reads those with strconv.ParseUint. SAMPLERATE takes any
// unsigned number, and 4, 1.0, .5 and 1e3 all reach here as written, for
// strconv.ParseFloat.
type AnalyzeOption struct {
	Loc   Loc
	Type  AnalyzeOptionType
	Value string
}

func (o *AnalyzeOption) nodeTag() {}

// AnalyzeTableStmt represents an ANALYZE TABLE statement.
//
// TiDB's ANALYZE is far wider than MySQL's: besides MySQL's table list and
// its UPDATE/DROP HISTOGRAM arms it takes one table's partitions, an index
// list, a column choice, and a WITH list of tuning options. The arms are
// mutually exclusive in ways only the server settles, and the fields below
// record which one a statement took:
//
//	Partitions   PARTITION p, ... -- one table only
//	Indexes      INDEX arm; HasIndex is what tells `INDEX` with an empty
//	             list from no INDEX at all
//	ColumnChoice ALL COLUMNS, PREDICATE COLUMNS, or COLUMNS with Columns
//	Incremental  INCREMENTAL, which TiDB takes only with the INDEX arm
//	Options      the WITH list, in the order written
type AnalyzeTableStmt struct {
	Loc              Loc
	NoWriteToBinlog  bool
	Incremental      bool
	Tables           []*TableRef
	Partitions       []string
	Indexes          []string
	HasIndex         bool
	ColumnChoice     AnalyzeColumnChoice
	Columns          []string
	HistogramOp      string   // "UPDATE" or "DROP"
	HistogramColumns []string // column names
	Options          []*AnalyzeOption
}

func (s *AnalyzeTableStmt) nodeTag()  {}
func (s *AnalyzeTableStmt) stmtNode() {}

// DropStatsStmt represents TiDB's DROP STATS, which discards the statistics
// of the tables or partitions it names. MySQL has no such statement.
//
// The three arms are a list of tables, one table with a list of its
// partitions, and one table with GLOBAL; Partitions and Global are set only
// by their own arm, and v8.5.5 answers 1064 for either beside a table list.
type DropStatsStmt struct {
	Loc        Loc
	Tables     []*TableRef
	Partitions []string
	Global     bool
}

func (s *DropStatsStmt) nodeTag()  {}
func (s *DropStatsStmt) stmtNode() {}

// LoadStatsStmt represents TiDB's LOAD STATS, which reads a statistics dump
// back from the file it names.
type LoadStatsStmt struct {
	Loc  Loc
	Path string
}

func (s *LoadStatsStmt) nodeTag()  {}
func (s *LoadStatsStmt) stmtNode() {}

// StatisticsType names the kind of extended statistics CREATE STATISTICS
// builds. TiDB requires exactly one of the three.
type StatisticsType int

const (
	// StatisticsCardinality is (CARDINALITY).
	StatisticsCardinality StatisticsType = iota
	// StatisticsDependency is (DEPENDENCY).
	StatisticsDependency
	// StatisticsCorrelation is (CORRELATION).
	StatisticsCorrelation
)

// String names the type as TiDB spells it.
func (t StatisticsType) String() string {
	switch t {
	case StatisticsCardinality:
		return "CARDINALITY"
	case StatisticsDependency:
		return "DEPENDENCY"
	case StatisticsCorrelation:
		return "CORRELATION"
	}
	return "UNKNOWN"
}

// CreateStatisticsStmt represents TiDB's CREATE STATISTICS, which builds one
// extended statistic over a table's columns.
//
//	CREATE STATISTICS [IF NOT EXISTS] name (type) ON tbl (col, ...)
type CreateStatisticsStmt struct {
	Loc         Loc
	IfNotExists bool
	Name        string
	Type        StatisticsType
	Table       *TableRef
	Columns     []*ColumnRef
}

func (s *CreateStatisticsStmt) nodeTag()  {}
func (s *CreateStatisticsStmt) stmtNode() {}

// DropStatisticsStmt represents TiDB's DROP STATISTICS, which drops one
// extended statistic by name. It takes one name and no IF EXISTS.
type DropStatisticsStmt struct {
	Loc  Loc
	Name string
}

func (s *DropStatisticsStmt) nodeTag()  {}
func (s *DropStatisticsStmt) stmtNode() {}

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
	ForExport       bool        // Never set: TiDB has no FOR EXPORT
	RelayChannel    string      // Never set: TiDB has no FLUSH RELAY LOGS
	Plugins         []string    // FLUSH TIDB PLUGINS name, ...
}

func (s *FlushStmt) nodeTag()  {}
func (s *FlushStmt) stmtNode() {}

// KillStmt represents a KILL statement.
type KillStmt struct {
	Loc Loc
	// ConnectionID is what the statement names: an *IntLit for a number, or
	// a *FuncCallExpr for KILL CONNECTION_ID(), the one other form the
	// grammar has. An IntLit holds an int64 and saturates there, so the
	// number itself is in ID.
	ConnectionID ExprNode
	// ID is the exact value of the number, which runs to the whole uint64
	// range. It is 0 for the function form.
	ID    uint64
	Query bool // KILL QUERY vs KILL CONNECTION
	// TiDB is KILL TIDB, the form TiDB documents as usable only with a client
	// connected to the server directly.
	TiDB bool
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

// CreatePlacementPolicyStmt represents a TiDB CREATE [OR REPLACE]
// PLACEMENT POLICY [IF NOT EXISTS] name option_list statement.
// Ref: TiDB v8.5.5 parser.y:15427-15436.
type CreatePlacementPolicyStmt struct {
	Loc         Loc
	OrReplace   bool
	IfNotExists bool
	Name        string
	Options     []*PlacementPolicyOption
}

func (s *CreatePlacementPolicyStmt) nodeTag()  {}
func (s *CreatePlacementPolicyStmt) stmtNode() {}

// AlterPlacementPolicyStmt represents a TiDB ALTER PLACEMENT POLICY
// [IF EXISTS] name option_list statement.
// Ref: TiDB v8.5.5 parser.y:15438-15446.
type AlterPlacementPolicyStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
	Options  []*PlacementPolicyOption
}

func (s *AlterPlacementPolicyStmt) nodeTag()  {}
func (s *AlterPlacementPolicyStmt) stmtNode() {}

// DropPlacementPolicyStmt represents a TiDB DROP PLACEMENT POLICY
// [IF EXISTS] name statement. Policies cannot be comma-listed.
// Ref: TiDB v8.5.5 parser.y:15389-15396.
type DropPlacementPolicyStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
}

func (s *DropPlacementPolicyStmt) nodeTag()  {}
func (s *DropPlacementPolicyStmt) stmtNode() {}

// BindingScope is the scope word a SQL plan binding statement carries.
// TiDB takes GLOBAL or SESSION, or neither, in which case the server
// picks the default for the statement; the parser records which was
// written rather than resolving it.
type BindingScope int

const (
	// BindingScopeUnset is a statement that names no scope.
	BindingScopeUnset BindingScope = iota
	// BindingScopeGlobal is GLOBAL.
	BindingScopeGlobal
	// BindingScopeSession is SESSION.
	BindingScopeSession
)

func (s BindingScope) String() string {
	switch s {
	case BindingScopeGlobal:
		return "GLOBAL"
	case BindingScopeSession:
		return "SESSION"
	}
	return ""
}

// CreateBindingStmt represents TiDB's CREATE BINDING, which records a plan
// for a statement shape. It has three forms:
//
//	CREATE [GLOBAL | SESSION] BINDING FOR stmt USING stmt
//	CREATE [GLOBAL | SESSION] BINDING USING stmt
//	CREATE [GLOBAL | SESSION] BINDING FROM HISTORY USING PLAN DIGEST digest [, digest] ...
//
// OriginStmt is the statement after FOR, the shape the binding matches,
// and is nil in the second and third forms, which have none. HintedStmt is
// the statement after USING, the plan to bind, and is nil in the third
// form, which names plans already in the statement history by digest in
// PlanDigests. A digest is a string literal or a user variable.
//
// Neither inner statement is run: a binding records a plan for later
// statements, so a consumer that reasons about what a change does stops at
// this node rather than reading through it.
//
// Ref: TiDB's SQL Plan Management,
// https://docs.pingcap.com/tidb/v8.5/sql-plan-management/ (every form
// measured against v8.5.5 in tidb/parser/binding_oracle_test.go).
type CreateBindingStmt struct {
	Loc         Loc
	Scope       BindingScope
	OriginStmt  StmtNode
	HintedStmt  StmtNode
	PlanDigests []ExprNode
}

func (s *CreateBindingStmt) nodeTag()  {}
func (s *CreateBindingStmt) stmtNode() {}

// DropBindingStmt represents TiDB's DROP BINDING:
//
//	DROP [GLOBAL | SESSION] BINDING FOR stmt [USING stmt]
//	DROP [GLOBAL | SESSION] BINDING FOR SQL DIGEST digest [, digest] ...
//
// OriginStmt is the statement after FOR and is nil in the SQL DIGEST form;
// HintedStmt is the statement after USING, nil when none is written.
// SQLDigests is non-empty only in the SQL DIGEST form, where a digest is a
// string literal or a user variable. As on CreateBindingStmt, neither
// inner statement is run.
//
// Ref: TiDB's SQL Plan Management,
// https://docs.pingcap.com/tidb/v8.5/sql-plan-management/.
type DropBindingStmt struct {
	Loc        Loc
	Scope      BindingScope
	OriginStmt StmtNode
	HintedStmt StmtNode
	SQLDigests []ExprNode
}

func (s *DropBindingStmt) nodeTag()  {}
func (s *DropBindingStmt) stmtNode() {}

// SetBindingStmt represents TiDB's SET BINDING:
//
//	SET BINDING {ENABLED | DISABLED} FOR stmt [USING stmt]
//	SET BINDING {ENABLED | DISABLED} FOR SQL DIGEST digest
//
// Enabled is ENABLED, and false is DISABLED. The statement takes no
// GLOBAL or SESSION scope word: v8.5.5 answers 1064 for `SET GLOBAL
// BINDING ENABLED FOR ...`. SQLDigest is the one digest of the second
// form, which takes a single string literal -- not a list, and not a user
// variable, unlike DROP BINDING FOR SQL DIGEST. As on CreateBindingStmt,
// neither inner statement is run.
//
// Ref: TiDB's SQL Plan Management,
// https://docs.pingcap.com/tidb/v8.5/sql-plan-management/.
type SetBindingStmt struct {
	Loc        Loc
	Enabled    bool
	OriginStmt StmtNode
	HintedStmt StmtNode
	SQLDigest  ExprNode
}

func (s *SetBindingStmt) nodeTag()  {}
func (s *SetBindingStmt) stmtNode() {}

// PlacementPolicyOption is one entry in a CREATE/ALTER PLACEMENT POLICY
// option list (DirectPlacementOption in TiDB v8.5.5 parser.y:2013-2066).
// The grammar accepts these as comma-separated OR whitespace-separated,
// duplicates are accepted at parse time (last-wins is a semantic concern).
//
// Value representation:
//   - String-shaped options (PRIMARY_REGION, REGIONS, SCHEDULE,
//     CONSTRAINTS and its leader/follower/voter/learner variants,
//     SURVIVAL_PREFERENCES) carry their raw quoted-content in Value.
//   - Integer-shaped options (FOLLOWERS, VOTERS, LEARNERS) use IntValue,
//     and Value holds the literal numeric text for round-trip fidelity.
//
// Name is the uppercased option keyword (e.g. "PRIMARY_REGION").
type PlacementPolicyOption struct {
	Loc      Loc
	Name     string
	Value    string // string-shaped options OR raw numeric text
	IntValue uint64 // numeric-shaped options (FOLLOWERS/VOTERS/LEARNERS)
	IsInt    bool   // true when option was a numeric literal
}

func (o *PlacementPolicyOption) nodeTag() {}

// CreateSequenceStmt represents a TiDB CREATE SEQUENCE statement:
//
//	CREATE SEQUENCE [IF NOT EXISTS] [schema.]name
//	    [sequence_option ...] [table_option [[,] table_option] ...]
//
// The grammar takes any table option. TiDB's DDL keeps COMMENT, ignores
// ENGINE, and rejects the rest with error 8227 (pkg/ddl/sequence.go
// buildSequenceInfo), which is not a syntax error.
// Ref: TiDB v8.5.5 parser.y:15462-15471.
type CreateSequenceStmt struct {
	Loc          Loc
	IfNotExists  bool
	Name         *TableRef
	Options      []*SequenceOption
	TableOptions []*TableOption
}

func (s *CreateSequenceStmt) nodeTag()  {}
func (s *CreateSequenceStmt) stmtNode() {}

// AlterSequenceStmt represents a TiDB ALTER SEQUENCE statement:
//
//	ALTER SEQUENCE [IF EXISTS] [schema.]name alter_sequence_option ...
//
// At least one option is required, RESTART is ALTER-only, and ALTER takes
// no table options.
// Ref: TiDB v8.5.5 parser.y:15598-15632.
type AlterSequenceStmt struct {
	Loc      Loc
	IfExists bool
	Name     *TableRef
	Options  []*SequenceOption
}

func (s *AlterSequenceStmt) nodeTag()  {}
func (s *AlterSequenceStmt) stmtNode() {}

// DropSequenceStmt represents a TiDB DROP SEQUENCE [IF EXISTS] name
// [, name] ... statement.
// Ref: TiDB v8.5.5 parser.y:15575-15582.
type DropSequenceStmt struct {
	Loc       Loc
	IfExists  bool
	Sequences []*TableRef
}

func (s *DropSequenceStmt) nodeTag()  {}
func (s *DropSequenceStmt) stmtNode() {}

// SequenceOptionType enumerates the options of SequenceOption and
// AlterSequenceOption (TiDB v8.5.5 parser.y:15489-15553, 15618-15632).
type SequenceOptionType int

const (
	SequenceOptionIncrementBy SequenceOptionType = iota // INCREMENT [BY | =] n
	SequenceOptionStartWith                             // START [WITH | =] n
	SequenceOptionMinValue                              // MINVALUE [=] n
	SequenceOptionNoMinValue                            // NOMINVALUE | NO MINVALUE
	SequenceOptionMaxValue                              // MAXVALUE [=] n
	SequenceOptionNoMaxValue                            // NOMAXVALUE | NO MAXVALUE
	SequenceOptionCache                                 // CACHE [=] n
	SequenceOptionNoCache                               // NOCACHE | NO CACHE
	SequenceOptionCycle                                 // CYCLE
	SequenceOptionNoCycle                               // NOCYCLE | NO CYCLE
	SequenceOptionRestart                               // RESTART (ALTER only)
	SequenceOptionRestartWith                           // RESTART [WITH | =] n (ALTER only)
)

// SequenceOption is one CREATE or ALTER SEQUENCE option. The list keeps
// source order and repeats because TiDB applies options in order, so a
// later option overrides an earlier one (pkg/ddl/sequence.go
// handleSequenceOptions, alterSequenceOptions).
type SequenceOption struct {
	Loc  Loc
	Type SequenceOptionType
	// Value is the signed integer of INCREMENT, START, MINVALUE, MAXVALUE,
	// CACHE, and RESTART WITH; it is 0 for the other options.
	Value int64
}

func (o *SequenceOption) nodeTag() {}

// DropIndexStmt represents a DROP INDEX statement.
type DropIndexStmt struct {
	Loc       Loc
	IfExists  bool
	Name      string
	Table     *TableRef
	Algorithm string
	Lock      string
	// IsHypo marks DROP HYPO INDEX, TiDB's drop of a hypothetical index. That
	// form takes no algorithm or lock options.
	IsHypo bool
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

// ParamMarkerExpr represents a `?` parameter marker (paramMarker in parser.y),
// the placeholder a prepared statement binds a value to. It is a leaf with no
// value of its own.
type ParamMarkerExpr struct {
	Loc Loc
}

func (e *ParamMarkerExpr) nodeTag()  {}
func (e *ParamMarkerExpr) exprNode() {}

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
// error (oracle mysql 8.0) — so deparse must render Keyword unquoted.
// Keyword is stored normalized: uppercase, SQL_TSI_ prefix stripped, and
// GET_FORMAT's TIMESTAMP folded to DATETIME, matching SHOW CREATE VIEW.
type KeywordArg struct {
	Loc     Loc
	Keyword string
}

func (e *KeywordArg) nodeTag()  {}
func (e *KeywordArg) exprNode() {}

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

// SetTransactionStmt represents a SET TRANSACTION statement.
type SetTransactionStmt struct {
	Loc            Loc
	Scope          string // GLOBAL, SESSION, or ""
	IsolationLevel string // READ UNCOMMITTED, READ COMMITTED, REPEATABLE READ, SERIALIZABLE
	AccessMode     string // READ ONLY, READ WRITE
	// AsOf is the expression of READ ONLY AS OF TIMESTAMP expr, the time a stale
	// read reads at; nil when there is none.
	AsOf ExprNode
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

// TableStmt represents a TABLE query primary (MySQL 8.0.19+).
// TABLE t is equivalent to SELECT * FROM t.
//
// The parser builds one only as the TableSource of a SelectStmt, which holds
// the clauses after it. A TABLE statement, and a TABLE in any other position,
// is a SelectStmt over a TableStmt.
type TableStmt struct {
	Loc   Loc
	Table *TableRef // table name
}

func (s *TableStmt) nodeTag()  {}
func (s *TableStmt) stmtNode() {}

// ValuesStmt represents a VALUES query primary (MySQL 8.0.19+).
// VALUES ROW(1,2), ROW(3,4) returns rows of literal values.
//
// The parser builds one only as the ValuesSource of a SelectStmt, which holds
// the clauses after it. A VALUES statement, and a VALUES in any other
// position, is a SelectStmt over a ValuesStmt.
type ValuesStmt struct {
	Loc  Loc
	Rows [][]ExprNode // list of ROW(...) value lists
}

func (s *ValuesStmt) nodeTag()  {}
func (s *ValuesStmt) stmtNode() {}

// ValuesColumnName returns the name of the i'th (0-based) output column of a
// VALUES query. A VALUES row carries no names of its own, so MySQL 8.0 calls
// the columns column_0, column_1, ... -- in a result header, in a derived
// table or CTE over one, and in the column list of a view whose body is one
// (container-verified).
// TiDB's own executor refuses a VALUES query outright (error 1051, Unknown
// table ”), so there is no TiDB result header to compare against; the names
// are MySQL's, whose VALUES statement TiDB's grammar implements.
func ValuesColumnName(i int) string {
	return "column_" + strconv.Itoa(i)
}

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

// CreateSpatialRefSysStmt represents a CREATE SPATIAL REFERENCE SYSTEM statement.
// The parser never builds one: TiDB has no spatial reference systems and answers
// 1064 for the statement. The type stays because callers build against it.
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
// The parser never builds one, for the reason CreateSpatialRefSysStmt gives.
type DropSpatialRefSysStmt struct {
	Loc      Loc
	IfExists bool
	SRID     int64
}

func (s *DropSpatialRefSysStmt) nodeTag()  {}
func (s *DropSpatialRefSysStmt) stmtNode() {}

// ResourceGroupOption is one option of CREATE or ALTER RESOURCE GROUP
// (DirectResourceGroupOption, TiDB v8.5.5 parser.y:1946). The options are
// separated by commas or by nothing at all, in any order, and each may be
// written with or without its `=`; a repeat of one the group already carries
// is a syntax error.
//
// Name is the option keyword, uppercased. The field that carries the value
// goes with the name:
//
//	RU_PER_SEC    IntValue, or Unlimited for RU_PER_SEC = UNLIMITED
//	PRIORITY      Value, one of LOW, MEDIUM, HIGH
//	BURSTABLE     Burstable; bare BURSTABLE is true, and a value is true
//	              when it is not zero, which is what v8.5.5 stores:
//	              BURSTABLE = 2 reads back as YES and BURSTABLE = 0 as NO
//	QUERY_LIMIT   Runaway; nil for QUERY_LIMIT = NULL and QUERY_LIMIT = ()
//	BACKGROUND    Background; nil for BACKGROUND = NULL and BACKGROUND = ()
//
// Burstable is only what the BURSTABLE option said. A group written with
// RU_PER_SEC = UNLIMITED and no BURSTABLE at all reads back from
// information_schema.resource_groups as burstable, but that is the server's
// rule about an unlimited rate, not something the option list carries.
type ResourceGroupOption struct {
	Loc        Loc
	Name       string
	Value      string
	IntValue   uint64
	Unlimited  bool
	Burstable  bool
	Runaway    []*ResourceGroupRunawayOption
	Background []*ResourceGroupBackgroundOption
}

func (o *ResourceGroupOption) nodeTag() {}

// ResourceGroupRunawayOption is one rule inside QUERY_LIMIT
// (DirectResourceGroupRunawayOption, TiDB v8.5.5 parser.y:1877).
//
//	EXEC_ELAPSED    Value, a duration time.ParseDuration reads
//	PROCESSED_KEYS  IntValue, unsigned: the grammar's intLit carries no sign
//	RU              IntValue, the same
//	ACTION          Action, one of DRYRUN, COOLDOWN, KILL, SWITCH_GROUP;
//	                SwitchGroup carries the group SWITCH_GROUP(name) names
//	WATCH           Watch, one of EXACT, SIMILAR, PLAN, with the optional
//	                DURATION in Duration, or DurationUnlimited for the bare
//	                keyword DURATION = UNLIMITED
//
// A repeated ACTION or WATCH is a syntax error; a repeated EXEC_ELAPSED,
// PROCESSED_KEYS or RU is not, and v8.5.5 keeps the last of them.
type ResourceGroupRunawayOption struct {
	Loc               Loc
	Name              string
	Value             string
	IntValue          uint64
	Action            string
	SwitchGroup       string
	Watch             string
	Duration          string
	DurationUnlimited bool
}

func (o *ResourceGroupRunawayOption) nodeTag() {}

// ResourceGroupBackgroundOption is one option inside BACKGROUND
// (DirectResourceGroupBackgroundOption, TiDB v8.5.5 parser.y:2019).
//
//	TASK_TYPES         Value, a comma-separated list inside one string
//	UTILIZATION_LIMIT  IntValue
type ResourceGroupBackgroundOption struct {
	Loc      Loc
	Name     string
	Value    string
	IntValue uint64
}

func (o *ResourceGroupBackgroundOption) nodeTag() {}

// CreateResourceGroupStmt represents a CREATE RESOURCE GROUP statement.
//
// Ref: TiDB v8.5.5 parser.y:15733
//
//	CREATE RESOURCE GROUP [IF NOT EXISTS] group_name option [[,] option] ...
//
// At least one option is required: TiDB answers 1064 for a group with none.
type CreateResourceGroupStmt struct {
	Loc         Loc
	IfNotExists bool
	Name        string
	Options     []*ResourceGroupOption
}

func (s *CreateResourceGroupStmt) nodeTag()  {}
func (s *CreateResourceGroupStmt) stmtNode() {}

// AlterResourceGroupStmt represents an ALTER RESOURCE GROUP statement.
//
// Ref: TiDB v8.5.5 parser.y:15743
//
//	ALTER RESOURCE GROUP [IF EXISTS] group_name option [[,] option] ...
type AlterResourceGroupStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
	Options  []*ResourceGroupOption
}

func (s *AlterResourceGroupStmt) nodeTag()  {}
func (s *AlterResourceGroupStmt) stmtNode() {}

// DropResourceGroupStmt represents a DROP RESOURCE GROUP statement.
//
// Ref: TiDB v8.5.5 parser.y:15753
//
//	DROP RESOURCE GROUP [IF EXISTS] group_name
//
// One group per statement: TiDB has no list here.
type DropResourceGroupStmt struct {
	Loc      Loc
	IfExists bool
	Name     string
}

func (s *DropResourceGroupStmt) nodeTag()  {}
func (s *DropResourceGroupStmt) stmtNode() {}

// CalibrateResourceOption is one option of CALIBRATE RESOURCE's dynamic form
// (DynamicCalibrateResourceOption, TiDB v8.5.5 parser.y:16929).
//
//	START_TIME  Ts, any expression
//	END_TIME    Ts, any expression
//	DURATION    Value, a duration TiDB's own parser reads (so '20d' counts),
//	            or, for DURATION INTERVAL expr unit, Ts and Unit
type CalibrateResourceOption struct {
	Loc   Loc
	Name  string
	Ts    ExprNode
	Value string
	Unit  string
}

func (o *CalibrateResourceOption) nodeTag() {}

// CalibrateResourceStmt represents a CALIBRATE RESOURCE statement.
//
// Ref: TiDB v8.5.5 parser.y:16882
//
//	CALIBRATE RESOURCE [WORKLOAD workload | option [[,] option] ...]
//
// The two forms do not mix: a statement carries a workload or dynamic
// options, never both, and may carry neither.
type CalibrateResourceStmt struct {
	Loc      Loc
	Workload string // TPCC, OLTP_READ_WRITE, OLTP_READ_ONLY, OLTP_WRITE_ONLY, TPCH_10
	Options  []*CalibrateResourceOption
}

func (s *CalibrateResourceStmt) nodeTag()  {}
func (s *CalibrateResourceStmt) stmtNode() {}

// QueryWatchOption is one option of QUERY WATCH ADD (QueryWatchOption,
// TiDB v8.5.5 parser.y:17010).
//
//	RESOURCE_GROUP  ResourceGroup, or ResourceGroupExpr for a user variable
//	ACTION          Action, as on a runaway rule, with SwitchGroup
//	SQL_DIGEST      Pattern
//	PLAN_DIGEST     Pattern
//	SQL_TEXT        Watch (EXACT, SIMILAR or PLAN) and Pattern
//
// One of each kind at most, and the three text kinds share one slot: a
// statement carries a single SQL_DIGEST, PLAN_DIGEST or SQL_TEXT.
type QueryWatchOption struct {
	Loc               Loc
	Name              string
	ResourceGroup     string
	ResourceGroupExpr ExprNode
	Action            string
	SwitchGroup       string
	Watch             string
	Pattern           ExprNode
}

func (o *QueryWatchOption) nodeTag() {}

// AddQueryWatchStmt represents a QUERY WATCH ADD statement.
//
// Ref: TiDB v8.5.5 parser.y:16980
//
//	QUERY WATCH ADD option [[,] option] ...
type AddQueryWatchStmt struct {
	Loc     Loc
	Options []*QueryWatchOption
}

func (s *AddQueryWatchStmt) nodeTag()  {}
func (s *AddQueryWatchStmt) stmtNode() {}

// DropQueryWatchStmt represents a QUERY WATCH REMOVE statement.
//
// Ref: TiDB v8.5.5 parser.y:17068
//
//	QUERY WATCH REMOVE watch_id
//
// v8.5.5 takes a watch id and nothing else: QUERY WATCH REMOVE RESOURCE
// GROUP rg, which pingcap's parser on master reads, is 1064 there, and so
// is a negative id.
//
// The id is unsigned because the grammar's NUM is: the lexer scans it with
// no sign and up to 2^64-1, which an int64 cannot hold without wrapping.
type DropQueryWatchStmt struct {
	Loc Loc
	ID  uint64
}

func (s *DropQueryWatchStmt) nodeTag()  {}
func (s *DropQueryWatchStmt) stmtNode() {}

// AlterViewStmt represents an ALTER VIEW statement. The parser never builds it:
// TiDB has no ALTER VIEW, and a view is redefined with CREATE OR REPLACE VIEW.
// It stays because Bytebase's adapters for the MySQL family switch on it.
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
	Body         Node   // parsed sp_proc_stmt (compound or simple); nil when no DO body.
	BodyText     string // body text as tokenized (the lexer expands conditional comments).
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
//
// Ref: TiDB v8.5.5 parser.y:11017
//
//	SET RESOURCE GROUP group_name
//
// MySQL's FOR thread_id list is not TiDB's: v8.5.5 answers 1064 for it.
type SetResourceGroupStmt struct {
	Loc  Loc
	Name string // group name
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

// SplitRegionStmt represents a TiDB SPLIT ... TABLE statement, which
// pre-splits the regions of a table or of one of its indexes.
//
//	SPLIT [REGION FOR] [PARTITION] TABLE tbl [PARTITION (p, ...)]
//	    [INDEX idx] split_option
//
// REGION FOR and PARTITION are spelling variants that TiDB's executor
// ignores, so RegionFor and Partition only record what was written.
//
// Ref: TiDB v8.5.5 parser.y:3239-3258 (SplitRegionStmt) and
// parser.y:3276-3299 (SplitSyntaxOption).
type SplitRegionStmt struct {
	Loc Loc
	// RegionFor is whether the statement spelled REGION FOR.
	RegionFor bool
	// Partition is whether the statement spelled the PARTITION keyword
	// before TABLE. It is not the PARTITION (p, ...) list, which lives on
	// Table.Partitions.
	Partition bool
	// Table is the one table being split, with its PARTITION (p, ...)
	// list in Table.Partitions.
	Table *TableRef
	// Index is the INDEX name, empty when the statement names none.
	Index string
	// IndexLoc is the INDEX name's own range, NoLoc() when there is none.
	IndexLoc Loc
	// Opt is the split specification, never nil: the grammar has no form
	// without one.
	Opt *SplitOption
}

func (s *SplitRegionStmt) nodeTag()  {}
func (s *SplitRegionStmt) stmtNode() {}

// SplitOption is a SPLIT statement's one split specification, in either the
// enumerated or the ranged form. Ranged says which, because both forms can
// leave every value row empty: TiDB's parser takes `BY ()` and
// `BETWEEN () AND () REGIONS 10` and rejects the empty row only when the
// statement runs (1136 and 1105 on v8.5.5).
//
//	BY (v, ...), (v, ...), ...
//	BETWEEN (v, ...) AND (v, ...) REGIONS n
//
// Ref: TiDB v8.5.5 parser.y:3260-3274 (SplitOption).
type SplitOption struct {
	Loc Loc
	// Ranged is whether this is the BETWEEN form.
	Ranged bool
	// ValueLists holds the BY form's rows, one slice per parenthesized
	// row, in the order written. A row written `()` is an empty,
	// non-nil slice. Nil for the ranged form.
	ValueLists [][]ExprNode
	// Lower and Upper are the ranged form's bound rows, each possibly
	// empty. Nil for the enumerated form.
	Lower []ExprNode
	Upper []ExprNode
	// Regions is the ranged form's REGIONS count, 0 for the enumerated
	// form. The grammar is Int64Num, so a sign, a decimal, a string, a
	// placeholder, and an out-of-range value are all syntax errors;
	// REGIONS 0 is a syntax-valid statement that fails when it runs.
	Regions int64
}

func (o *SplitOption) nodeTag() {}

// DistributeTableStmt represents a TiDB DISTRIBUTE TABLE statement, which
// schedules a table's or partition's replicas across an engine's stores.
//
//	DISTRIBUTE TABLE tbl [PARTITION (p, ...)]
//	    RULE [=] 'rule' ENGINE [=] 'engine' [TIMEOUT [=] 'timeout']
//
// RULE and ENGINE are both mandatory and fixed in that order; each value is
// a string literal and the `=` is optional. TIMEOUT, when present, comes
// last. TiDB rejects any other order, a bare identifier value, and a comma
// between the options (all 1064 on v8.5.5).
//
// Ref: TiDB v8.5.5 parser.y:3202-3221 (DistributeTableStmt).
type DistributeTableStmt struct {
	Loc Loc
	// Table is the one table being distributed, with its
	// PARTITION (p, ...) list in Table.Partitions.
	Table *TableRef
	// Rule and Engine are the mandatory RULE and ENGINE values.
	Rule   string
	Engine string
	// Timeout is the TIMEOUT value, empty when the statement has no
	// TIMEOUT clause. TIMEOUT = '' is written the same way, as upstream
	// records it.
	Timeout string
}

func (s *DistributeTableStmt) nodeTag()  {}
func (s *DistributeTableStmt) stmtNode() {}

// RawStmt wraps a statement node with its location span.
type RawStmt struct {
	Loc     Loc
	StmtLen int
	Stmt    StmtNode
}

func (s *RawStmt) nodeTag()  {}
func (s *RawStmt) stmtNode() {}

// AdminKind names the subcommand an AdminStmt carries. TiDB gives each ADMIN
// subcommand its own production and its own operands rather than one
// permissive tail, so the kind, not the shape, says what a statement does.
//
// Ref: TiDB v8.5.5 parser.y:11204-11425 (AdminStmt).
type AdminKind int

// The zero AdminKind is no kind, so a node built without one never reads as a
// subcommand.
const (
	AdminShowDDL AdminKind = iota + 1
	AdminShowDDLJobs
	AdminShowDDLJobQueries
	AdminShowDDLJobQueriesWithRange
	AdminShowNextRowID
	AdminShowSlow
	AdminShowBDRRole
	AdminCheckTable
	AdminCheckIndex
	AdminCheckIndexRange
	AdminChecksumTable
	AdminRecoverIndex
	AdminCleanupIndex
	AdminCleanupTableLock
	AdminRepairTable
	AdminCancelDDLJobs
	AdminPauseDDLJobs
	AdminResumeDDLJobs
	AdminAlterDDLJobs
	AdminReloadExprPushdownBlacklist
	AdminReloadOptRuleBlacklist
	AdminReloadStatistics
	AdminReloadBindings
	AdminFlushBindings
	AdminCaptureBindings
	AdminEvolveBindings
	AdminFlushPlanCache
	AdminPluginEnable
	AdminPluginDisable
	AdminSetBDRRole
	AdminUnsetBDRRole
)

// AdminStatementScope is the scope of ADMIN FLUSH ... PLAN_CACHE. TiDB's
// StatementScope defaults to SESSION when the scope is left out.
//
// Ref: TiDB v8.5.5 parser.y:16553-16568 (StatementScope).
type AdminStatementScope int

const (
	AdminScopeSession AdminStatementScope = iota
	AdminScopeGlobal
	AdminScopeInstance
)

// AdminSlowKind is which statements ADMIN SHOW SLOW TOP reports on.
type AdminSlowKind int

const (
	AdminSlowKindDefault AdminSlowKind = iota
	AdminSlowKindInternal
	AdminSlowKindAll
)

// AdminSlowSpec is the tail of ADMIN SHOW SLOW: RECENT n, or TOP [INTERNAL |
// ALL] n. Recent and the kind do not combine -- RECENT takes no kind -- so
// Recent reading true leaves Kind at its zero value.
//
// Ref: TiDB v8.5.5 parser.y:11460-11487 (AdminShowSlow).
type AdminSlowSpec struct {
	Loc    Loc
	Recent bool
	Kind   AdminSlowKind
	Count  uint64
}

func (s *AdminSlowSpec) nodeTag() {}

// AdminHandleRange is one (begin, end) pair of ADMIN CHECK INDEX's handle
// range list. Both bounds are unsigned integer literals in the grammar, so
// neither carries a sign.
//
// Ref: TiDB v8.5.5 parser.y:11489-11501 (HandleRangeList, HandleRange).
type AdminHandleRange struct {
	Loc   Loc
	Begin int64
	End   int64
}

func (r *AdminHandleRange) nodeTag() {}

// AdminLimit is the LIMIT tail of ADMIN SHOW DDL JOB QUERIES: LIMIT count,
// LIMIT offset, count, or LIMIT count OFFSET offset.
//
// Ref: TiDB v8.5.5 parser.y:11427-11440 (AdminStmtLimitOpt).
type AdminLimit struct {
	Loc    Loc
	Count  uint64
	Offset uint64
}

func (l *AdminLimit) nodeTag() {}

// AdminJobOption is one `name = value` of ADMIN ALTER DDL JOBS. The grammar
// takes any identifier for the name and any signed literal for the value;
// TiDB's executor keeps THREAD, BATCH_SIZE, and MAX_WRITE_SPEED and answers
// 1105 for anything else, which is a runtime answer, not a parse error.
//
// Ref: TiDB v8.5.5 parser.y:11442-11456 (AlterJobOptionList, AlterJobOption).
type AdminJobOption struct {
	Loc   Loc
	Name  string
	Value ExprNode
}

func (o *AdminJobOption) nodeTag() {}

// AdminStmt represents a TiDB ADMIN statement. Kind says which subcommand it
// is, and each subcommand fills only its own operands: a field a kind does not
// take is left at its zero value.
//
//	Tables      CHECK/CHECKSUM TABLE, CLEANUP TABLE LOCK, SHOW ... NEXT_ROW_ID,
//	            CHECK/RECOVER/CLEANUP INDEX, REPAIR TABLE
//	Index       CHECK/RECOVER/CLEANUP INDEX
//	Ranges      CHECK INDEX ... (begin, end), ...
//	Where       SHOW DDL JOBS [n] WHERE ...
//	JobNumber   SHOW DDL JOBS n, ALTER DDL JOBS n
//	JobIDs      CANCEL/PAUSE/RESUME DDL JOBS, SHOW DDL JOB QUERIES
//	Limit       SHOW DDL JOB QUERIES LIMIT ...
//	Slow        SHOW SLOW
//	Plugins     PLUGINS ENABLE/DISABLE
//	Scope       FLUSH [SESSION|GLOBAL|INSTANCE] PLAN_CACHE
//	BDRRole     SET BDR ROLE
//	JobOptions  ALTER DDL JOBS
//	CreateTable REPAIR TABLE
//
// Ref: TiDB v8.5.5 parser.y:11204-11425 (AdminStmt).
type AdminStmt struct {
	Loc  Loc
	Kind AdminKind

	Tables []*TableRef
	Index  string
	Ranges []*AdminHandleRange

	Where ExprNode

	// JobNumber is the job count of SHOW DDL JOBS n and the job id of ALTER
	// DDL JOBS n. HasJobNumber separates SHOW DDL JOBS 0 from SHOW DDL JOBS,
	// which takes no number.
	JobNumber    int64
	HasJobNumber bool
	JobIDs       []int64
	Limit        *AdminLimit

	Slow    *AdminSlowSpec
	Plugins []string
	Scope   AdminStatementScope
	// BDRRole is PRIMARY or SECONDARY, uppercased; the grammar takes no
	// other word and no quoted name.
	BDRRole     string
	JobOptions  []*AdminJobOption
	CreateTable *CreateTableStmt
}

func (s *AdminStmt) nodeTag()  {}
func (s *AdminStmt) stmtNode() {}

// BRIEKind says which of TiDB's backup and restore statements a BRIEStmt is.
// Upstream calls the family BRIE, for backup, restore, import and export.
//
// Ref: TiDB v8.5.5 parser.y:5623-5721 (BRIEStmt).
type BRIEKind int

const (
	// BRIEBackup is BACKUP DATABASE|SCHEMA|TABLE ... TO 'storage'.
	BRIEBackup BRIEKind = iota
	// BRIERestore is RESTORE DATABASE|SCHEMA|TABLE ... FROM 'storage'.
	BRIERestore
	// BRIEBackupLogs is BACKUP LOGS TO 'storage', which starts a log
	// backup task rather than taking one snapshot.
	BRIEBackupLogs
	// BRIERestorePoint is RESTORE POINT FROM 'storage', the point-in-time
	// restore.
	BRIERestorePoint
	// BRIEStopLogs is STOP BACKUP LOGS.
	BRIEStopLogs
	// BRIEPauseLogs is PAUSE BACKUP LOGS [option ...].
	BRIEPauseLogs
	// BRIEResumeLogs is RESUME BACKUP LOGS.
	BRIEResumeLogs
	// BRIEPurgeLogs is PURGE BACKUP LOGS FROM 'storage' [option ...].
	BRIEPurgeLogs
	// BRIEShowLogsStatus is SHOW BACKUP LOGS STATUS.
	BRIEShowLogsStatus
	// BRIEShowLogsMetadata is SHOW BACKUP LOGS METADATA FROM 'storage'.
	BRIEShowLogsMetadata
	// BRIEShowBackupMetadata is SHOW BACKUP METADATA FROM 'storage'.
	BRIEShowBackupMetadata
	// BRIEShowJob is SHOW BR JOB job_id.
	BRIEShowJob
	// BRIEShowJobQuery is SHOW BR JOB QUERY job_id.
	BRIEShowJobQuery
	// BRIECancelJob is CANCEL BR JOB job_id.
	BRIECancelJob
)

// IsLogTaskOrJob reports whether the statement manages a log backup task, reads
// backup metadata, or addresses a BR job, as opposed to taking or restoring a
// backup: everything from BRIEStopLogs on. None of them names a database or a
// table, and none changes a schema.
func (k BRIEKind) IsLogTaskOrJob() bool {
	return k >= BRIEStopLogs
}

// BRIETarget says what a BACKUP or RESTORE names, which the grammar's
// BRIETables decides. BACKUP LOGS and RESTORE POINT name nothing.
//
// Ref: TiDB v8.5.5 parser.y:5725-5737 (BRIETables).
type BRIETarget int

const (
	// BRIETargetNone is BACKUP LOGS, RESTORE POINT and the log task and job
	// statements, which name no database and no table.
	BRIETargetNone BRIETarget = iota
	// BRIETargetAllSchemas is DATABASE *, every database in the cluster.
	// The asterisk is the bare token: `DATABASE `*`` is a one-element
	// database list naming a database called `*`.
	BRIETargetAllSchemas
	// BRIETargetSchemas is DATABASE a, b, c, with the names in Schemas.
	BRIETargetSchemas
	// BRIETargetTables is TABLE a.b, c, with the tables in Tables.
	BRIETargetTables
)

// BRIEStmt represents one of TiDB's backup and restore statements.
//
//	BACKUP  DATABASE {* | db, ...} TO   'storage' [option ...]
//	BACKUP  TABLE    tbl, ...      TO   'storage' [option ...]
//	BACKUP  LOGS                   TO   'storage' [option ...]
//	RESTORE DATABASE {* | db, ...} FROM 'storage' [option ...]
//	RESTORE TABLE    tbl, ...      FROM 'storage' [option ...]
//	RESTORE POINT                  FROM 'storage' [option ...]
//
// SCHEMA is DATABASE: TiDB's lexer maps the two words to one token, so the
// node cannot say which was written, and neither can upstream's.
//
// Ref: TiDB v8.5.5 parser.y:5623-5721 (BRIEStmt), 5725-5747 (BRIETables,
// DBNameList).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-backup/
type BRIEStmt struct {
	Loc  Loc
	Kind BRIEKind
	// Target says which of Schemas and Tables is filled, and separates
	// DATABASE * from both.
	Target BRIETarget
	// Schemas are the database names of the DATABASE a, b, c form.
	Schemas []string
	// Tables are the tables of the TABLE form, each a bare or qualified
	// name with no alias, partition list, or index hint.
	Tables []*TableRef
	// Storage is the destination or source URL, the one string literal
	// after TO or FROM. The grammar's stringLit is a single literal here:
	// `TO 'a' 'b'` is a syntax error, where the same two literals
	// concatenate inside an expression.
	Storage string
	// StorageLoc is the string literal's own range.
	StorageLoc Loc
	// Options is the option list, in written order, empty when the
	// statement gave none. A repeated option is not an error and is kept
	// twice; TiDB's executor takes the last.
	Options []*BRIEOption
	// JobID is the job of SHOW BR JOB, SHOW BR JOB QUERY and CANCEL BR JOB.
	JobID int64
}

func (s *BRIEStmt) nodeTag()  {}
func (s *BRIEStmt) stmtNode() {}

// BRIEOptionValue says which of a BRIEOption's value fields carries its
// value. The option's name does not decide this on its own: SNAPSHOT and
// LAST_BACKUP each take a string or a number, SNAPSHOT also takes
// `n <unit> AGO`, and CHECKSUM and ANALYZE each take a boolean or a level.
//
// Ref: TiDB v8.5.5 parser.y:5871-6039 (BRIEOption).
type BRIEOptionValue int

const (
	// BRIEValueInt is a value in IntValue: the integer options, and a
	// boolean option, where TRUE and FALSE arrive as 1 and 0.
	BRIEValueInt BRIEOptionValue = iota
	// BRIEValueString is a value in StrValue, from a string literal or,
	// for BACKEND and ON_DUPLICATE, from an identifier.
	BRIEValueString
	// BRIEValueLevel is CHECKSUM or ANALYZE written as a level, with OFF,
	// OPTIONAL or REQUIRED uppercased in StrValue.
	BRIEValueLevel
	// BRIEValueTimeAgo is SNAPSHOT n <unit> AGO, with n in IntValue and
	// the unit uppercased in Unit.
	BRIEValueTimeAgo
	// BRIEValueRateLimit is RATE_LIMIT n MB/SECOND, with n in IntValue.
	// MB per second is the only rate the grammar spells.
	BRIEValueRateLimit
	// BRIEValueCSVHeaderWords is CSV_HEADER = FIELDS or COLUMNS, with the
	// word uppercased in StrValue.
	BRIEValueCSVHeaderWords
)

// BRIEOption is one option of a BACKUP or RESTORE statement.
//
// Name is the option's own word, uppercased. `ON DUPLICATE` and
// `ON_DUPLICATE` are the same option and both arrive as ON_DUPLICATE.
// Value says which field holds the value.
//
// Every option name is read as a bare word, which is what upstream's
// keyword tokens come to: v8.5.5 answers 1064 for every backquoted
// spelling, `CONCURRENCY` and `OFF` included, and takes a backquoted
// identifier only where the grammar names one, which is BACKEND's and
// ON_DUPLICATE's value.
//
// Ref: TiDB v8.5.5 parser.y:5759-5869 (the option name rules),
// 5871-6039 (BRIEOption).
type BRIEOption struct {
	Loc   Loc
	Name  string
	Value BRIEOptionValue
	// StrValue carries the value of BRIEValueString, BRIEValueLevel and
	// BRIEValueCSVHeaderWords.
	StrValue string
	// IntValue carries the value of BRIEValueInt, BRIEValueTimeAgo and
	// BRIEValueRateLimit, unscaled: `RATE_LIMIT 500 MB/SECOND` is 500 and
	// `SNAPSHOT 5 MINUTE AGO` is 5, where upstream multiplies each out.
	IntValue uint64
	// Unit is the time unit of BRIEValueTimeAgo, uppercased, one of the
	// grammar's TimestampUnit words.
	Unit string
}

func (o *BRIEOption) nodeTag() {}

// RecoverTableStmt represents TiDB's RECOVER TABLE, which brings back a
// table dropped within the GC lifetime.
//
//	RECOVER TABLE BY JOB job_id
//	RECOVER TABLE tbl [job_num]
//
// The two forms are different productions, not one with an optional BY JOB:
// the first names a DDL job and no table, and the second names one table
// and an optional number. ByJob says which was written, so
// `RECOVER TABLE BY JOB 11` and `RECOVER TABLE t1 100` never read alike.
//
// Each number is the grammar's Int64Num, an unsigned literal checked
// against int64: v8.5.5 takes `BY JOB 0` and `t1 0` and answers 1064 for
// `BY JOB -1`, `t1 -1` and 9223372036854775808.
//
// Ref: TiDB v8.5.5 parser.y:3039-3059 (RecoverTableStmt).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-recover-table/
type RecoverTableStmt struct {
	Loc Loc
	// ByJob is whether the statement named a DDL job rather than a table.
	ByJob bool
	// JobID is the job id of the BY JOB form.
	JobID int64
	// Table is the one table of the table form, nil in the BY JOB form.
	Table *TableRef
	// JobNum is the optional number after the table, and HasJobNum
	// separates `RECOVER TABLE t1 0` from `RECOVER TABLE t1`.
	JobNum    int64
	HasJobNum bool
}

func (s *RecoverTableStmt) nodeTag()  {}
func (s *RecoverTableStmt) stmtNode() {}

// FlashbackTableStmt represents FLASHBACK TABLE, which brings back a table
// dropped or truncated within the GC lifetime, optionally under a new name.
//
//	FLASHBACK TABLE tbl [TO new_name]
//
// This is not a rename of a live table: the table named must already be
// gone, and the statement creates the new name. Its sibling makes that
// plain -- FLASHBACK DATABASE on a database that is still there answers
// 1007, "Can't create database 'd'; database exists", on v8.5.5.
//
// It is a different statement from FLASHBACK TABLE ... TO TIMESTAMP, which
// is FlashbackToTimestampStmt. Upstream separates the two in the lexer:
// `TO TIMESTAMP` becomes one token only when a string literal follows it,
// so `FLASHBACK TABLE t TO timestamp` renames to a table called
// `timestamp` and `FLASHBACK TABLE t TO TIMESTAMP '...'` does not.
//
// Ref: TiDB v8.5.5 parser.y:3134-3156 (FlashbackTableStmt,
// FlashbackToNewName), lexer.go:272-280 (the toTimestamp token).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-flashback-table/
type FlashbackTableStmt struct {
	Loc Loc
	// Table is the table being brought back, bare or qualified.
	Table *TableRef
	// NewName is the TO name, empty when the statement gave none. It is a
	// bare identifier, never qualified: `TO c.d` is a syntax error.
	NewName string
	// NewNameLoc is the new name's own range, NoLoc() when there is none.
	NewNameLoc Loc
}

func (s *FlashbackTableStmt) nodeTag()  {}
func (s *FlashbackTableStmt) stmtNode() {}

// FlashbackDatabaseStmt represents FLASHBACK DATABASE, which brings back a
// dropped database, optionally under a new name.
//
//	FLASHBACK {DATABASE | SCHEMA} db [TO new_db]
//
// Like FlashbackTableStmt this creates the new name rather than renaming a
// live database: v8.5.5 answers 1007, "Can't create database 'd'; database
// exists", for FLASHBACK DATABASE d when d is still there. It is a
// different statement from FLASHBACK DATABASE ... TO TIMESTAMP.
//
// Ref: TiDB v8.5.5 parser.y:3160-3166 (FlashbackDatabaseStmt).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-flashback-database/
type FlashbackDatabaseStmt struct {
	Loc Loc
	// Name is the database being brought back, a bare identifier:
	// `FLASHBACK DATABASE a.b` is a syntax error.
	Name string
	// NameLoc is that name's own range.
	NameLoc Loc
	// NewName is the TO name, empty when the statement gave none.
	NewName string
	// NewNameLoc is the new name's own range, NoLoc() when there is none.
	NewNameLoc Loc
}

func (s *FlashbackDatabaseStmt) nodeTag()  {}
func (s *FlashbackDatabaseStmt) stmtNode() {}

// FlashbackScope says what a FLASHBACK ... TO TIMESTAMP or TO TSO rolls
// back.
type FlashbackScope int

const (
	// FlashbackCluster is FLASHBACK CLUSTER, the whole cluster.
	FlashbackCluster FlashbackScope = iota
	// FlashbackTables is FLASHBACK TABLE t, ..., one or more tables.
	FlashbackTables
	// FlashbackDatabase is FLASHBACK DATABASE|SCHEMA db.
	FlashbackDatabase
)

// FlashbackToTimestampStmt represents FLASHBACK ... TO TIMESTAMP and
// FLASHBACK ... TO TSO, which roll data back to a point in time.
//
//	FLASHBACK CLUSTER              TO TIMESTAMP 'ts'
//	FLASHBACK TABLE tbl, ...       TO TIMESTAMP 'ts'
//	FLASHBACK {DATABASE|SCHEMA} db TO TIMESTAMP 'ts'
//	FLASHBACK CLUSTER              TO TSO n
//	FLASHBACK TABLE tbl, ...       TO TSO n
//	FLASHBACK {DATABASE|SCHEMA} db TO TSO n
//
// The timestamp is one string literal, unvalidated by the parser: v8.5.5
// takes `TO TIMESTAMP 'x'` and fails it when the statement runs. The TSO
// is the grammar's LengthNum and must be above zero: v8.5.5 answers 1064
// for `TO TSO 0`, `TO TSO -100` and `TO TSO 18446744073709551616`.
//
// Ref: TiDB v8.5.5 parser.y:3067-3124 (FlashbackToTimestampStmt).
// Ref: https://docs.pingcap.com/tidb/stable/sql-statement-flashback-cluster/
type FlashbackToTimestampStmt struct {
	Loc Loc
	// Scope says which of Tables and Database is filled, and separates the
	// cluster form, which fills neither.
	Scope FlashbackScope
	// Tables are the tables of the TABLE form.
	Tables []*TableRef
	// Database is the database of the DATABASE form.
	Database string
	// DatabaseLoc is that name's own range, NoLoc() in the other forms.
	DatabaseLoc Loc
	// IsTSO says whether TO TSO was written rather than TO TIMESTAMP, and
	// so which of TSO and Timestamp holds the point rolled back to.
	IsTSO bool
	// Timestamp is the TO TIMESTAMP string, as written without its quotes.
	Timestamp string
	// TSO is the TO TSO value, always above zero.
	TSO uint64
}

func (s *FlashbackToTimestampStmt) nodeTag()  {}
func (s *FlashbackToTimestampStmt) stmtNode() {}

// RecommendIndexAction says which RECOMMEND INDEX statement a node is.
type RecommendIndexAction int

const (
	// RecommendIndexRun is RECOMMEND INDEX RUN [FOR 'sql'] [WITH option, ...].
	RecommendIndexRun RecommendIndexAction = iota
	// RecommendIndexShow is RECOMMEND INDEX SHOW.
	RecommendIndexShow
	// RecommendIndexApply is RECOMMEND INDEX APPLY id. v8.5.5 parses it and
	// answers 1105 "unsupported action apply" when it runs.
	RecommendIndexApply
	// RecommendIndexIgnore is RECOMMEND INDEX IGNORE id, answered the same way.
	RecommendIndexIgnore
	// RecommendIndexSet is RECOMMEND INDEX SET option, ....
	RecommendIndexSet
)

// RecommendIndexStmt represents RECOMMEND INDEX, TiDB's index advisor.
//
// Ref: TiDB v8.5.5 parser.y:14346-14426 (RecommendIndexStmt).
// Ref: https://docs.pingcap.com/tidb/stable/index-advisor/
type RecommendIndexStmt struct {
	Loc    Loc
	Action RecommendIndexAction
	// HasSQL says FOR was written, and SQL is its one string literal, which
	// the parser does not read: `RECOMMEND INDEX RUN FOR 'x'` parses and fails
	// when it runs. The two are apart because FOR '' is a statement of its own.
	HasSQL bool
	SQL    string
	// ID is the recommendation APPLY and IGNORE name.
	ID int64
	// Options holds RUN's WITH list and SET's list, in written order.
	Options []*RecommendIndexOption
}

func (s *RecommendIndexStmt) nodeTag()  {}
func (s *RecommendIndexStmt) stmtNode() {}

// RecommendIndexOption is one name = literal of a RECOMMEND INDEX statement.
// The grammar takes a Literal, not an expression: -1, (1) and 1 + 1 are 1064s.
type RecommendIndexOption struct {
	Loc   Loc
	Name  string
	Value ExprNode
}

func (o *RecommendIndexOption) nodeTag() {}

// SetConfigStmt represents SET CONFIG, which changes a cluster component's
// configuration while it runs. The statement names either a component type or
// one instance's address, and HasInstance says which.
//
// Ref: TiDB v8.5.5 parser.y:10781-10788 (SetStmt), 10931-10940 (ConfigItemName).
// Ref: https://docs.pingcap.com/tidb/stable/dynamic-config/
type SetConfigStmt struct {
	Loc Loc
	// Type is the component word, lowercased as upstream records it. The
	// parser takes any name: v8.5.5 answers 1105 "unknown type" for the ones
	// it does not know when the statement runs.
	Type string
	// HasInstance says the statement named an address, the string in Instance.
	HasInstance bool
	Instance    string
	// Name is the configuration item, its parts joined by the `.` and `-` they
	// were written with, unquoted: log.level, auto-compaction-mode.
	Name string
	// Value is the SetExpr, so ON and BINARY arrive as *StringLit and DEFAULT
	// as *DefaultExpr, the way a variable assignment's do.
	Value ExprNode
}

func (s *SetConfigStmt) nodeTag()  {}
func (s *SetConfigStmt) stmtNode() {}

// AlterRangeStmt represents ALTER RANGE name PLACEMENT POLICY ..., which sets
// the placement of a named key range.
//
// The parser takes any name. v8.5.5 answers 1105 "range name is not supported"
// for one other than global and meta when the statement runs, and 8239 for a
// policy that does not exist.
//
// Ref: TiDB v8.5.5 parser.y:13764-13769 (AlterRangeStmt), 2068-2084
// (PlacementPolicyOption).
type AlterRangeStmt struct {
	Loc  Loc
	Name string
	// Policy is the policy as written, by name or as a string, and "DEFAULT" for
	// the keyword forms. v8.5.5 clears the policy for every spelling of default,
	// so a consumer compares it without case: `default` and 'default' clear it
	// just as DEFAULT and SET DEFAULT do.
	Policy string
}

func (s *AlterRangeStmt) nodeTag()  {}
func (s *AlterRangeStmt) stmtNode() {}

// CancelDistributionJobStmt represents CANCEL DISTRIBUTION JOB job_id, which
// cancels a DISTRIBUTE TABLE job.
//
// Ref: TiDB v8.5.5 parser.y:3198-3204 (CancelDistributionJobStmt).
type CancelDistributionJobStmt struct {
	Loc   Loc
	JobID int64
}

func (s *CancelDistributionJobStmt) nodeTag()  {}
func (s *CancelDistributionJobStmt) stmtNode() {}

// CancelImportJobStmt represents CANCEL IMPORT JOB job_id, which cancels an
// IMPORT INTO job.
//
// Ref: TiDB v8.5.5 parser.y:6089-6096 (CancelImportStmt).
type CancelImportJobStmt struct {
	Loc   Loc
	JobID int64
}

func (s *CancelImportJobStmt) nodeTag()  {}
func (s *CancelImportJobStmt) stmtNode() {}
