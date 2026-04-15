package ast

// This file holds DDL AST node types for CREATE TABLE (T2.1) and ALTER TABLE (T2.2).

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

// ALTER TABLE statement and supporting types (T2.2)
// ---------------------------------------------------------------------------

// AlterActionType identifies the kind of action in an ALTER TABLE statement.
type AlterActionType int

const (
	AlterActionUnknown      AlterActionType = iota
	AlterAddColumn                          // ADD COLUMN col_def [AFTER col | FIRST]
	AlterDropColumn                         // DROP COLUMN col
	AlterModifyColumn                       // MODIFY COLUMN col_def [AFTER col | FIRST]
	AlterRenameColumn                       // RENAME COLUMN old TO new (or old new)
	AlterRenameTable                        // RENAME TO new_name (or RENAME new_name)
	AlterRenameRollup                       // RENAME ROLLUP old new
	AlterRenamePartition                    // RENAME PARTITION old new
	AlterAddPartition                       // ADD PARTITION ...
	AlterDropPartition                      // DROP PARTITION name
	AlterTruncatePartition                  // TRUNCATE PARTITION name
	AlterReplacePartition                   // REPLACE PARTITION ... WITH TEMPORARY ...
	AlterAddRollup                          // ADD ROLLUP name (cols) [FROM base] [PROPERTIES(...)]
	AlterDropRollup                         // DROP ROLLUP name
	AlterSetProperties                      // SET ("key"="value", ...)
	AlterModifyPartition                    // MODIFY PARTITION p SET ("key"="val")
	AlterModifyDistribution                 // MODIFY DISTRIBUTION DISTRIBUTED BY ...
	AlterModifyComment                      // MODIFY COMMENT 'text'
	AlterModifyEngine                       // MODIFY ENGINE TO engine
	AlterEnableFeature                      // ENABLE FEATURE 'name' [WITH PROPERTIES (...)]
	AlterOrderBy                            // ORDER BY (col1, col2, ...)
	AlterRaw                                // fallback: raw text for unparsed actions
)

// AlterTableAction represents a single action in an ALTER TABLE statement.
// Doris supports many action types; we use a generic struct with Type
// and typed fields for common actions.
type AlterTableAction struct {
	Type AlterActionType
	// Column actions (ADD/MODIFY COLUMN):
	Column     *ColumnDef // for ADD/MODIFY COLUMN
	ColumnName string     // for DROP/RENAME COLUMN
	NewName    string     // for RENAME COLUMN (new name) or RENAME ROLLUP/PARTITION
	After      string     // optional AFTER column name
	First      bool       // FIRST position flag
	// Rename table:
	NewTableName *ObjectName // for RENAME TO new_name
	// Partition actions:
	Partition     *PartitionItem // for ADD PARTITION
	PartitionName string         // for DROP/TRUNCATE/MODIFY PARTITION
	PartitionList []string       // for MODIFY PARTITION (p1, p2) or (*)
	PartitionStar bool           // MODIFY PARTITION (*)
	// ADD PARTITION extras:
	PartitionDist *DistributionDesc // optional DISTRIBUTED BY after ADD PARTITION
	PartitionProps []*Property      // optional properties after ADD PARTITION
	// Rollup actions:
	Rollup     *RollupDef // for ADD ROLLUP
	RollupName string     // for DROP ROLLUP / RENAME ROLLUP (old name)
	// Properties (SET / MODIFY PARTITION SET):
	Properties []*Property
	// MODIFY DISTRIBUTION:
	Distribution *DistributionDesc
	// MODIFY COMMENT:
	Comment string
	// MODIFY ENGINE:
	Engine string
	// ENABLE FEATURE:
	FeatureName string
	// ORDER BY columns:
	OrderByColumns []string
	// Raw text for unsupported/complex actions:
	RawText string
	Loc     Loc
}

// Tag implements Node.
func (n *AlterTableAction) Tag() NodeTag { return T_AlterTableAction }

var _ Node = (*AlterTableAction)(nil)

// AlterTableStmt represents:
//
//	ALTER TABLE name action [, action ...]
type AlterTableStmt struct {
	Name    *ObjectName
	Actions []*AlterTableAction
	Loc     Loc
}

// Tag implements Node.
func (n *AlterTableStmt) Tag() NodeTag { return T_AlterTableStmt }

var _ Node = (*AlterTableStmt)(nil)

// ---------------------------------------------------------------------------
// VIEW DDL nodes (T2.4)
// ---------------------------------------------------------------------------

// ViewColumn is one entry in a CREATE/ALTER VIEW column list:
//
//	name [COMMENT 'text']
type ViewColumn struct {
	Name    string
	Comment string
	Loc     Loc
}

// Tag implements Node.
func (n *ViewColumn) Tag() NodeTag { return T_ViewColumn }

var _ Node = (*ViewColumn)(nil)

// CreateViewStmt represents:
//
//	CREATE [OR REPLACE] VIEW [IF NOT EXISTS] view_name
//	    [(col1 [COMMENT 'text'], ...)]
//	    [COMMENT 'view comment']
//	    AS query
type CreateViewStmt struct {
	Name        *ObjectName
	OrReplace   bool
	IfNotExists bool
	Columns     []*ViewColumn // optional column list
	Comment     string        // view-level COMMENT
	Query       Node          // SELECT statement (*SelectStmt)
	Loc         Loc
}

// Tag implements Node.
func (n *CreateViewStmt) Tag() NodeTag { return T_CreateViewStmt }

var _ Node = (*CreateViewStmt)(nil)

// AlterViewStmt represents:
//
//	ALTER VIEW view_name [(col1 [COMMENT 'text'], ...)] AS query
type AlterViewStmt struct {
	Name    *ObjectName
	Columns []*ViewColumn
	Query   Node
	Loc     Loc
}

// Tag implements Node.
func (n *AlterViewStmt) Tag() NodeTag { return T_AlterViewStmt }

var _ Node = (*AlterViewStmt)(nil)

// DropViewStmt represents:
//
//	DROP VIEW [IF EXISTS] view_name
type DropViewStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropViewStmt) Tag() NodeTag { return T_DropViewStmt }

var _ Node = (*DropViewStmt)(nil)

// ---------------------------------------------------------------------------
// WORKLOAD GROUP DDL nodes (T5.4)
// ---------------------------------------------------------------------------

// CreateWorkloadGroupStmt represents:
//
//	CREATE WORKLOAD GROUP [IF NOT EXISTS] name PROPERTIES(...)
type CreateWorkloadGroupStmt struct {
	Name        string
	IfNotExists bool
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateWorkloadGroupStmt) Tag() NodeTag { return T_CreateWorkloadGroupStmt }

var _ Node = (*CreateWorkloadGroupStmt)(nil)

// AlterWorkloadGroupStmt represents:
//
//	ALTER WORKLOAD GROUP name PROPERTIES(...)
type AlterWorkloadGroupStmt struct {
	Name       string
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterWorkloadGroupStmt) Tag() NodeTag { return T_AlterWorkloadGroupStmt }

var _ Node = (*AlterWorkloadGroupStmt)(nil)

// DropWorkloadGroupStmt represents:
//
//	DROP WORKLOAD GROUP [IF EXISTS] name
type DropWorkloadGroupStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropWorkloadGroupStmt) Tag() NodeTag { return T_DropWorkloadGroupStmt }

var _ Node = (*DropWorkloadGroupStmt)(nil)

// ---------------------------------------------------------------------------
// WORKLOAD POLICY DDL nodes (T5.4)
// ---------------------------------------------------------------------------

// WorkloadPolicyItem holds a raw-text capture of one condition or action item
// from a WORKLOAD POLICY CONDITIONS/ACTIONS clause.
type WorkloadPolicyItem struct {
	RawText string
	Loc     Loc
}

// Tag implements Node.
func (n *WorkloadPolicyItem) Tag() NodeTag { return T_WorkloadPolicyItem }

var _ Node = (*WorkloadPolicyItem)(nil)

// CreateWorkloadPolicyStmt represents:
//
//	CREATE WORKLOAD POLICY [IF NOT EXISTS] name
//	    CONDITIONS(condition_list)
//	    ACTIONS(action_list)
//	    [PROPERTIES(...)]
type CreateWorkloadPolicyStmt struct {
	Name        string
	IfNotExists bool
	Conditions  []*WorkloadPolicyItem
	Actions     []*WorkloadPolicyItem
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateWorkloadPolicyStmt) Tag() NodeTag { return T_CreateWorkloadPolicyStmt }

var _ Node = (*CreateWorkloadPolicyStmt)(nil)

// AlterWorkloadPolicyStmt represents:
//
//	ALTER WORKLOAD POLICY name PROPERTIES(...)
type AlterWorkloadPolicyStmt struct {
	Name       string
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterWorkloadPolicyStmt) Tag() NodeTag { return T_AlterWorkloadPolicyStmt }

var _ Node = (*AlterWorkloadPolicyStmt)(nil)

// DropWorkloadPolicyStmt represents:
//
//	DROP WORKLOAD POLICY [IF EXISTS] name
type DropWorkloadPolicyStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropWorkloadPolicyStmt) Tag() NodeTag { return T_DropWorkloadPolicyStmt }

var _ Node = (*DropWorkloadPolicyStmt)(nil)

// ---------------------------------------------------------------------------
// RESOURCE DDL nodes (T5.4)
// ---------------------------------------------------------------------------

// CreateResourceStmt represents:
//
//	CREATE [EXTERNAL] RESOURCE [IF NOT EXISTS] name PROPERTIES(...)
type CreateResourceStmt struct {
	Name        string
	IfNotExists bool
	External    bool
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateResourceStmt) Tag() NodeTag { return T_CreateResourceStmt }

var _ Node = (*CreateResourceStmt)(nil)

// AlterResourceStmt represents:
//
//	ALTER RESOURCE name PROPERTIES(...)
type AlterResourceStmt struct {
	Name       string
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterResourceStmt) Tag() NodeTag { return T_AlterResourceStmt }

var _ Node = (*AlterResourceStmt)(nil)

// DropResourceStmt represents:
//
//	DROP RESOURCE [IF EXISTS] name
type DropResourceStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropResourceStmt) Tag() NodeTag { return T_DropResourceStmt }

var _ Node = (*DropResourceStmt)(nil)

// ---------------------------------------------------------------------------
// SQL BLOCK RULE DDL nodes (T5.4)
// ---------------------------------------------------------------------------

// CreateSQLBlockRuleStmt represents:
//
//	CREATE SQL_BLOCK_RULE [IF NOT EXISTS] name PROPERTIES(...)
type CreateSQLBlockRuleStmt struct {
	Name        string
	IfNotExists bool
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateSQLBlockRuleStmt) Tag() NodeTag { return T_CreateSQLBlockRuleStmt }

var _ Node = (*CreateSQLBlockRuleStmt)(nil)

// AlterSQLBlockRuleStmt represents:
//
//	ALTER SQL_BLOCK_RULE name PROPERTIES(...)
type AlterSQLBlockRuleStmt struct {
	Name       string
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *AlterSQLBlockRuleStmt) Tag() NodeTag { return T_AlterSQLBlockRuleStmt }

var _ Node = (*AlterSQLBlockRuleStmt)(nil)

// DropSQLBlockRuleStmt represents:
//
//	DROP SQL_BLOCK_RULE [IF EXISTS] name
type DropSQLBlockRuleStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropSQLBlockRuleStmt) Tag() NodeTag { return T_DropSQLBlockRuleStmt }

var _ Node = (*DropSQLBlockRuleStmt)(nil)
