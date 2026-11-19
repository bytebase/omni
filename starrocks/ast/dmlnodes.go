package ast

// This file holds AST node types for DML statements (T4.x).

// ---------------------------------------------------------------------------
// INSERT statement (T4.1)
// ---------------------------------------------------------------------------

// InsertStmt represents an INSERT INTO or INSERT OVERWRITE TABLE statement.
//
//	INSERT [INTO | OVERWRITE TABLE] [TEMPORARY PARTITION] table_name
//	    [PARTITION(p1, p2, ...) | PARTITION(*)]
//	    [WITH LABEL label_name]
//	    [(col1, col2, ...)]
//	    { VALUES (expr, ...) [, (...)] | SELECT ... | WITH ... SELECT ... }
type InsertStmt struct {
	Overwrite     bool        // true for INSERT OVERWRITE TABLE; false for INSERT INTO
	Target        *ObjectName // target table name
	Label         string      // optional WITH LABEL name; empty if absent
	Branch        string      // optional write branch (VERSION AS OF); empty if absent
	Partition     []string    // optional PARTITION(p1, p2, ...) names; nil if absent
	KeyPartitions []*Property // optional PARTITION(k = v, ...) key selectors
	TempPartition bool        // true if TEMPORARY PARTITION was specified
	PartitionStar bool        // true if PARTITION(*) was used
	Columns       []string    // optional column list; nil if absent
	ByName        bool        // true for INSERT ... BY NAME (column-name matching)
	Blackhole     bool        // INSERT INTO BLACKHOLE() sink
	FileTarget    []*Property // FILES(propertyList) target; nil for a table target (mutually exclusive with Target)
	Properties    []*Property // optional PROPERTIES(...) before the source
	// Source: exactly one of Values or Query is non-nil.
	Values [][]Node // VALUES rows: each inner slice is one row's expressions
	Query  Node     // SELECT or WITH…SELECT (*SelectStmt); nil if VALUES form
	Loc    Loc
}

// Tag implements Node.
func (n *InsertStmt) Tag() NodeTag { return T_InsertStmt }

// Compile-time assertion that *InsertStmt satisfies Node.
var _ Node = (*InsertStmt)(nil)

// ---------------------------------------------------------------------------
// Assignment: col = expr (used in UPDATE SET) (T4.2)
// ---------------------------------------------------------------------------

// Assignment represents a single SET assignment in UPDATE:
//
//	col = expr
//	t.col = expr
type Assignment struct {
	Column *ObjectName // column name; may be qualified (t.col)
	Value  Node        // right-hand expression
	Loc    Loc
}

// Tag implements Node.
func (n *Assignment) Tag() NodeTag { return T_Assignment }

// Compile-time assertion that *Assignment satisfies Node.
var _ Node = (*Assignment)(nil)

// ---------------------------------------------------------------------------
// UpdateStmt (T4.2)
// ---------------------------------------------------------------------------

// UpdateStmt represents an UPDATE statement:
//
//	UPDATE table [AS alias]
//	    SET col1 = expr1 [, col2 = expr2 ...]
//	    [FROM table_refs]
//	    [WHERE condition]
type UpdateStmt struct {
	With        *WithClause   // optional leading CTE clause (StarRocks WITH … UPDATE); nil if absent
	Target      *ObjectName   // table to update
	TargetAlias string        // optional alias (AS alias or bare alias)
	Assignments []*Assignment // SET clause assignments
	From        []Node        // optional FROM clause table references
	Where       Node          // optional WHERE expression
	Loc         Loc
}

// Tag implements Node.
func (n *UpdateStmt) Tag() NodeTag { return T_UpdateStmt }

// Compile-time assertion that *UpdateStmt satisfies Node.
var _ Node = (*UpdateStmt)(nil)

// ---------------------------------------------------------------------------
// DeleteStmt (T4.2)
// ---------------------------------------------------------------------------

// DeleteStmt represents a DELETE statement:
//
//	DELETE FROM table [AS alias]
//	    [PARTITION(p1 [, p2 ...])]
//	    [USING table_refs]
//	    [WHERE condition]
type DeleteStmt struct {
	With        *WithClause // optional leading CTE clause (StarRocks WITH … DELETE); nil if absent
	Target      *ObjectName // table to delete from
	TargetAlias string      // optional alias
	Partition   []string    // optional PARTITION(p1, p2, ...) names
	Using       []Node      // optional USING clause table references
	Where       Node        // optional WHERE expression
	Loc         Loc
}

// Tag implements Node.
func (n *DeleteStmt) Tag() NodeTag { return T_DeleteStmt }

// Compile-time assertion that *DeleteStmt satisfies Node.
var _ Node = (*DeleteStmt)(nil)

// ---------------------------------------------------------------------------
// TRUNCATE TABLE statement (T6.1)
// ---------------------------------------------------------------------------

// TruncateTableStmt represents a TRUNCATE TABLE statement:
//
//	TRUNCATE TABLE [db.]name [PARTITION(p1, p2, ...)] [FORCE]
type TruncateTableStmt struct {
	Target    *ObjectName // table name
	Partition []string    // optional PARTITION(p1, p2, ...) names; nil if absent
	Force     bool        // optional FORCE flag
	Loc       Loc
}

// Tag implements Node.
func (n *TruncateTableStmt) Tag() NodeTag { return T_TruncateTableStmt }

// Compile-time assertion that *TruncateTableStmt satisfies Node.
var _ Node = (*TruncateTableStmt)(nil)

// ---------------------------------------------------------------------------
// LOAD DATA statement (T6.1)
// ---------------------------------------------------------------------------

// LoadDataDesc describes one dataDesc of a LOAD statement
// (StarRocks.g4:1447):
//
//	DATA INFILE ('path', ...) [NEGATIVE] INTO TABLE t [partitionNames]
//	    [COLUMNS TERMINATED BY 'c'] [ROWS TERMINATED BY 'r'] [FORMAT AS f]
//	    [(SKIP_HEADER = n TRIM_SPACE = b ENCLOSE = 'e' ESCAPE = 'e')]
//	    [(c1, ...)] [COLUMNS FROM PATH AS (c, ...)] [SET (expr, ...)] [WHERE expr]
//	DATA FROM TABLE src [NEGATIVE] INTO TABLE t [partitionNames]
//	    [SET (expr, ...)] [WHERE expr]
type LoadDataDesc struct {
	SourceFiles     []string    // DATA INFILE ('f1', 'f2', ...)
	SourceTable     *ObjectName // DATA FROM TABLE src
	Negative        bool        // NEGATIVE modifier
	Target          *ObjectName // INTO TABLE target
	TemporaryPart   bool        // TEMPORARY PARTITION
	Partition       []string    // optional partition names, or "col = value" key partitions
	ColumnSeparator string      // COLUMNS TERMINATED BY
	RowDelimiter    string      // ROWS TERMINATED BY
	Format          string      // FORMAT AS csv/parquet/...
	FormatProps     []*Property // (SKIP_HEADER = 1 ...), values as written
	ColumnList      []string    // optional (c1, c2, ...)
	ColumnsFromPath []string    // optional COLUMNS FROM PATH AS (...)
	ColumnMappings  []Node      // optional SET (expr, ...)
	Where           Node        // optional WHERE expr
	Loc             Loc
}

// Tag implements Node.
func (n *LoadDataDesc) Tag() NodeTag { return T_LoadDataDesc }

// Compile-time assertion that *LoadDataDesc satisfies Node.
var _ Node = (*LoadDataDesc)(nil)

// LoadDataStmt represents a LOAD statement (StarRocks.g4:1427):
//
//	LOAD LABEL [db.]label [(data_desc, ...)]
//	    [WITH BROKER [name] [(...)] | WITH RESOURCE name [(...)]]
//	    [BY system] [PROPERTIES (...)]
type LoadDataStmt struct {
	Label            string          // load label (db_name.label_name or label_name)
	LabelParts       []string        // label parts as written: [db_name, label_name] or [label_name]
	DataDescs        []*LoadDataDesc // data descriptions, if any
	StorageType      string          // BROKER or RESOURCE, or ""
	BrokerName       string          // WITH BROKER or WITH RESOURCE name
	BrokerProperties []*Property     // the broker's or resource's (...) list
	System           string          // BY system
	Properties       []*Property     // optional PROPERTIES(...)
	Loc              Loc
}

// Tag implements Node.
func (n *LoadDataStmt) Tag() NodeTag { return T_LoadDataStmt }

// Compile-time assertion that *LoadDataStmt satisfies Node.
var _ Node = (*LoadDataStmt)(nil)

// AlterLoadStmt represents alterLoadStatement (StarRocks.g4:1499):
//
//	ALTER LOAD FOR (db=qualifiedName '.')? name=identifier jobProperties?
type AlterLoadStmt struct {
	Name       *ObjectName // the load's [db.]name
	Properties []*Property // optional PROPERTIES(...)
	Loc        Loc
}

// Tag implements Node.
func (n *AlterLoadStmt) Tag() NodeTag { return T_AlterLoadStmt }

var _ Node = (*AlterLoadStmt)(nil)

// ---------------------------------------------------------------------------
// EXPORT statement (T6.1)
// ---------------------------------------------------------------------------

// ExportStmt represents an EXPORT TABLE statement:
//
//	EXPORT TABLE name [PARTITION(p1, p2, ...)] [WHERE expr]
//	    TO 'path'
//	    [PROPERTIES(...)]
//	    [WITH BROKER name (...)]
type ExportStmt struct {
	Target           *ObjectName // table to export
	Partition        []string    // optional partitionNames
	ColumnAliases    []string    // optional columnAliases
	Path             string      // TO 'path'
	Mode             string      // WITH SYNC MODE / WITH ASYNC MODE
	Properties       []*Property // optional properties
	Broker           bool        // WITH BROKER seen
	BrokerName       string      // optional WITH BROKER name
	BrokerProperties []*Property // optional WITH BROKER name (...) list
	Loc              Loc
}

// Tag implements Node.
func (n *ExportStmt) Tag() NodeTag { return T_ExportStmt }

// Compile-time assertion that *ExportStmt satisfies Node.
var _ Node = (*ExportStmt)(nil)

// ROUTINE LOAD statements (T6.2)
// ---------------------------------------------------------------------------

// CreateRoutineLoadStmt represents a CREATE ROUTINE LOAD statement:
//
//	CREATE ROUTINE LOAD [db.]job_name ON table_name
//	    [load_properties]
//	    [PROPERTIES (...)]
//	    FROM {KAFKA | S3 | HDFS} (data_source_properties)
//	    [COMMENT 'text']
type CreateRoutineLoadStmt struct {
	Name                 *ObjectName // job name, optionally qualified with db
	OnTable              *ObjectName // ON table name
	LoadProps            string      // the load properties' source text, as written
	JobProperties        []*Property // PROPERTIES (...) key=value pairs
	DataSourceType       string      // "KAFKA", "S3", "HDFS", etc.
	DataSourceProperties []*Property // FROM ... (...) key=value pairs
	Comment              string      // optional COMMENT 'text'
	Loc                  Loc
}

// Tag implements Node.
func (n *CreateRoutineLoadStmt) Tag() NodeTag { return T_CreateRoutineLoadStmt }

// Compile-time assertion that *CreateRoutineLoadStmt satisfies Node.
var _ Node = (*CreateRoutineLoadStmt)(nil)

// AlterRoutineLoadStmt represents an ALTER ROUTINE LOAD statement:
//
//	ALTER ROUTINE LOAD FOR [db.]job_name
//	    [PROPERTIES (...)]
//	    [FROM KAFKA (...)]
type AlterRoutineLoadStmt struct {
	Name                 *ObjectName
	LoadProps            string // the load properties' source text, as written
	Properties           []*Property
	DataSourceType       string
	DataSourceProperties []*Property
	Loc                  Loc
}

// Tag implements Node.
func (n *AlterRoutineLoadStmt) Tag() NodeTag { return T_AlterRoutineLoadStmt }

// Compile-time assertion that *AlterRoutineLoadStmt satisfies Node.
var _ Node = (*AlterRoutineLoadStmt)(nil)

// PauseRoutineLoadStmt represents a PAUSE ROUTINE LOAD statement:
//
//	PAUSE ROUTINE LOAD FOR [db.]job_name
//	PAUSE ALL ROUTINE LOAD [FOR db]
type PauseRoutineLoadStmt struct {
	Name *ObjectName // nil when All=true
	All  bool        // true for PAUSE ALL ROUTINE LOAD
	For  string      // optional FOR db (used with All=true)
	Loc  Loc
}

// Tag implements Node.
func (n *PauseRoutineLoadStmt) Tag() NodeTag { return T_PauseRoutineLoadStmt }

// Compile-time assertion that *PauseRoutineLoadStmt satisfies Node.
var _ Node = (*PauseRoutineLoadStmt)(nil)

// ResumeRoutineLoadStmt represents a RESUME ROUTINE LOAD statement:
//
//	RESUME ROUTINE LOAD FOR [db.]job_name
//	RESUME ALL ROUTINE LOAD [FOR db]
type ResumeRoutineLoadStmt struct {
	Name *ObjectName // nil when All=true
	All  bool        // true for RESUME ALL ROUTINE LOAD
	For  string      // optional FOR db (used with All=true)
	Loc  Loc
}

// Tag implements Node.
func (n *ResumeRoutineLoadStmt) Tag() NodeTag { return T_ResumeRoutineLoadStmt }

// Compile-time assertion that *ResumeRoutineLoadStmt satisfies Node.
var _ Node = (*ResumeRoutineLoadStmt)(nil)

// StopRoutineLoadStmt represents a STOP ROUTINE LOAD statement:
//
//	STOP ROUTINE LOAD FOR [db.]job_name
type StopRoutineLoadStmt struct {
	Name *ObjectName
	Loc  Loc
}

// Tag implements Node.
func (n *StopRoutineLoadStmt) Tag() NodeTag { return T_StopRoutineLoadStmt }

// Compile-time assertion that *StopRoutineLoadStmt satisfies Node.
var _ Node = (*StopRoutineLoadStmt)(nil)

// ShowRoutineLoadStmt represents a SHOW ROUTINE LOAD statement:
//
//	SHOW [ALL] ROUTINE LOAD [FOR [db.]name | LIKE 'pattern' | FROM db]
type ShowRoutineLoadStmt struct {
	Name *ObjectName // optional FOR [db.]name
	All  bool        // true for SHOW ALL ROUTINE LOAD
	Like string      // optional LIKE 'pattern'
	From string      // optional FROM db

	Where   Node           // optional WHERE expression
	OrderBy []*OrderByItem // optional ORDER BY items
	Limit   Node           // optional LIMIT count
	Offset  Node           // optional LIMIT offset
	Create  bool           // SHOW CREATE ROUTINE LOAD name
	Loc     Loc
}

// Tag implements Node.
func (n *ShowRoutineLoadStmt) Tag() NodeTag { return T_ShowRoutineLoadStmt }

// Compile-time assertion that *ShowRoutineLoadStmt satisfies Node.
var _ Node = (*ShowRoutineLoadStmt)(nil)

// ShowRoutineLoadTaskStmt represents a SHOW ROUTINE LOAD TASK statement:
//
//	SHOW ROUTINE LOAD TASK FROM db WHERE ...
type ShowRoutineLoadTaskStmt struct {
	Name  *ObjectName // optional job name
	From  string      // optional FROM db
	Where Node        // optional WHERE condition
	Loc   Loc
}

// Tag implements Node.
func (n *ShowRoutineLoadTaskStmt) Tag() NodeTag { return T_ShowRoutineLoadTaskStmt }

// Compile-time assertion that *ShowRoutineLoadTaskStmt satisfies Node.
var _ Node = (*ShowRoutineLoadTaskStmt)(nil)

// SyncStmt represents a SYNC statement.
type SyncStmt struct {
	Loc Loc
}

// Tag implements Node.
func (n *SyncStmt) Tag() NodeTag { return T_SyncStmt }

// Compile-time assertion that *SyncStmt satisfies Node.
var _ Node = (*SyncStmt)(nil)
