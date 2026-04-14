package ast

// This file holds DDL AST node types for CREATE TABLE (T2.1).

// ---------------------------------------------------------------------------
// CREATE TABLE statement and supporting types
// ---------------------------------------------------------------------------

// ConstraintType classifies a table-level constraint.
type ConstraintType int

const (
	ConstraintPrimaryKey ConstraintType = iota
	ConstraintUnique
)

// CreateTableStmt represents:
//
//	CREATE [EXTERNAL] [TEMPORARY] TABLE [IF NOT EXISTS] name
//	    (column_def, ... [, index_def, ...])
//	    [ENGINE = name]
//	    [AGGREGATE|UNIQUE|DUPLICATE KEY (cols)]
//	    [COMMENT 'text']
//	    [PARTITION BY RANGE|LIST (cols) (...)]
//	    [DISTRIBUTED BY HASH(cols)|RANDOM [BUCKETS n|AUTO]]
//	    [ROLLUP (...)]
//	    [PROPERTIES (...)]
//	    [AS query]
type CreateTableStmt struct {
	Name          *ObjectName
	IfNotExists   bool
	External      bool
	Temporary     bool
	Columns       []*ColumnDef
	Indexes       []*IndexDef       // inline INDEX definitions
	Constraints   []*TableConstraint // table-level constraints
	KeyDesc       *KeyDesc          // AGGREGATE KEY / UNIQUE KEY / DUPLICATE KEY
	PartitionBy   *PartitionDesc    // PARTITION BY RANGE/LIST
	DistributedBy *DistributionDesc // DISTRIBUTED BY HASH/RANDOM
	Rollup        []*RollupDef      // ROLLUP (...)
	Properties    []*Property       // PROPERTIES (...)
	Engine        string            // ENGINE = xxx
	Comment       string            // COMMENT 'xxx'
	Like          *ObjectName       // CREATE TABLE ... LIKE other_table
	AsSelect      *RawQuery         // CREATE TABLE ... AS SELECT ...
	Loc           Loc
}

// Tag implements Node.
func (n *CreateTableStmt) Tag() NodeTag { return T_CreateTableStmt }

var _ Node = (*CreateTableStmt)(nil)

// ColumnDef represents a column definition in a CREATE TABLE statement.
//
//	col_name data_type [KEY] [agg_type]
//	    [GENERATED ALWAYS AS (expr)]
//	    [NOT NULL | NULL]
//	    [AUTO_INCREMENT]
//	    [DEFAULT value]
//	    [ON UPDATE CURRENT_TIMESTAMP]
//	    [COMMENT 'str']
type ColumnDef struct {
	Name       string
	Type       *TypeName
	Nullable   *bool  // nil = not specified, true = NULL, false = NOT NULL
	Default    Node   // DEFAULT expression (nil if absent)
	Comment    string // COMMENT string
	IsKey      bool   // KEY keyword
	AggType    string // aggregate type: SUM, MAX, MIN, REPLACE, REPLACE_IF_NOT_NULL, etc.
	AutoInc    bool   // AUTO_INCREMENT
	Generated  Node   // GENERATED ALWAYS AS (expr) -- the expression
	OnUpdate   string // ON UPDATE CURRENT_TIMESTAMP
	Loc        Loc
}

// Tag implements Node.
func (n *ColumnDef) Tag() NodeTag { return T_ColumnDef }

var _ Node = (*ColumnDef)(nil)

// IndexDef represents an inline index definition in CREATE TABLE.
//
//	INDEX [IF NOT EXISTS] name (col1, col2, ...) [USING type] [PROPERTIES(...)] [COMMENT 'str']
type IndexDef struct {
	Name        string
	Columns     []string
	IfNotExists bool
	IndexType   string      // BITMAP, INVERTED, NGRAM_BF, ANN, or ""
	Properties  []*Property // inline PROPERTIES(...)
	Comment     string
	Loc         Loc
}

// Tag implements Node.
func (n *IndexDef) Tag() NodeTag { return T_IndexDef }

var _ Node = (*IndexDef)(nil)

// TableConstraint represents a table-level constraint: PRIMARY KEY or UNIQUE.
type TableConstraint struct {
	Type    ConstraintType
	Name    string   // constraint name (may be empty)
	Columns []string
	Loc     Loc
}

// Tag implements Node.
func (n *TableConstraint) Tag() NodeTag { return T_TableConstraint }

var _ Node = (*TableConstraint)(nil)

// KeyDesc represents the table-level key declaration:
//
//	AGGREGATE KEY(cols) | UNIQUE KEY(cols) | DUPLICATE KEY(cols)
type KeyDesc struct {
	Type    string   // "AGGREGATE", "UNIQUE", "DUPLICATE"
	Columns []string
	Loc     Loc
}

// Tag implements Node.
func (n *KeyDesc) Tag() NodeTag { return T_KeyDesc }

var _ Node = (*KeyDesc)(nil)

// PartitionDesc represents:
//
//	[AUTO] PARTITION BY RANGE|LIST (cols/functions) (partition_defs...)
type PartitionDesc struct {
	Type       string           // "RANGE" or "LIST"
	Auto       bool             // AUTO PARTITION
	Columns    []string         // column names or function expressions (stored as raw strings)
	FuncExprs  []string         // if partition columns include function expressions
	Partitions []*PartitionItem // individual partition definitions
	Loc        Loc
}

// Tag implements Node.
func (n *PartitionDesc) Tag() NodeTag { return T_PartitionDesc }

var _ Node = (*PartitionDesc)(nil)

// PartitionItem represents one partition definition within PARTITION BY.
//
// For LESS THAN: PARTITION name VALUES LESS THAN (values) or MAXVALUE
// For IN: PARTITION name VALUES IN ((values), (values))
// For step: FROM (values) TO (values) INTERVAL n [unit]
type PartitionItem struct {
	Name       string   // partition name (empty for step partitions)
	IsMaxValue bool     // VALUES LESS THAN MAXVALUE
	Values     []string // raw value strings for LESS THAN / fixed
	InValues   [][]string // for IN partitions: list of value lists
	// Step partition fields
	IsStep     bool
	FromValues []string
	ToValues   []string
	Interval   string // interval amount
	IntervalUnit string // interval unit (e.g., "DAY")
	Loc        Loc
}

// Tag implements Node.
func (n *PartitionItem) Tag() NodeTag { return T_PartitionItem }

var _ Node = (*PartitionItem)(nil)

// DistributionDesc represents:
//
//	DISTRIBUTED BY HASH(cols) [BUCKETS n | BUCKETS AUTO]
//	DISTRIBUTED BY RANDOM [BUCKETS n | BUCKETS AUTO]
type DistributionDesc struct {
	Type    string   // "HASH" or "RANDOM"
	Columns []string // for HASH; empty for RANDOM
	Buckets int      // bucket count; 0 if not specified
	Auto    bool     // BUCKETS AUTO
	Loc     Loc
}

// Tag implements Node.
func (n *DistributionDesc) Tag() NodeTag { return T_DistributionDesc }

var _ Node = (*DistributionDesc)(nil)

// RollupDef represents one rollup definition in ROLLUP (...).
//
//	rollup_name (col1, col2, ...) [DUPLICATE KEY (key_cols)] [PROPERTIES(...)]
type RollupDef struct {
	Name       string
	Columns    []string
	DupKeys    []string    // optional DUPLICATE KEY columns
	Properties []*Property // optional PROPERTIES
	Loc        Loc
}

// Tag implements Node.
func (n *RollupDef) Tag() NodeTag { return T_RollupDef }

var _ Node = (*RollupDef)(nil)

// RawQuery is a placeholder for a query (SELECT statement) that has not been
// fully parsed. Used for CTAS (CREATE TABLE ... AS SELECT ...).
type RawQuery struct {
	RawText string
	Loc     Loc
}

// Tag implements Node.
func (n *RawQuery) Tag() NodeTag { return T_RawQuery }

var _ Node = (*RawQuery)(nil)
