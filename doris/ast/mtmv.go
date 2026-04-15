package ast

// This file holds AST node types for Doris Materialized View (MTMV) DDL (T5.1).
//
// Doris supports two kinds of materialized views:
//   - Sync materialized views (CREATE MATERIALIZED VIEW ... AS SELECT ...)
//     that refresh synchronously with DML and are attached to a base table.
//   - Async (MTMV) materialized views that have explicit refresh scheduling.
//
// The nodes here cover both variants.

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// MTMVRefreshTrigger represents the ON SCHEDULE / ON COMMIT / ON MANUAL
// trigger clause in CREATE MATERIALIZED VIEW.
type MTMVRefreshTrigger struct {
	OnManual   bool
	OnCommit   bool
	OnSchedule bool
	Interval   string // raw interval text, e.g., "1 DAY"
	StartsAt   string // optional STARTS timestamp string literal value
	Loc        Loc
}

// Tag implements Node.
func (n *MTMVRefreshTrigger) Tag() NodeTag { return T_MTMVRefreshTrigger }

var _ Node = (*MTMVRefreshTrigger)(nil)

// CreateMTMVStmt represents:
//
//	CREATE MATERIALIZED VIEW [IF NOT EXISTS] mv_name
//	    [BUILD IMMEDIATE | DEFERRED]
//	    [REFRESH COMPLETE | AUTO | NEVER | INCREMENTAL]
//	    [(col_def, ...)]
//	    [COMMENT 'text']
//	    [PARTITION BY ...]
//	    [DISTRIBUTED BY ...]
//	    [ON SCHEDULE EVERY interval [STARTS ts] | ON COMMIT | ON MANUAL]
//	    [PROPERTIES (...)]
//	    AS query
type CreateMTMVStmt struct {
	Name           *ObjectName
	IfNotExists    bool
	BuildMode      string             // "IMMEDIATE", "DEFERRED", or ""
	RefreshMethod  string             // "COMPLETE", "AUTO", "NEVER", "INCREMENTAL", or ""
	RefreshTrigger *MTMVRefreshTrigger // ON SCHEDULE / ON COMMIT / ON MANUAL
	Columns        []*ViewColumn      // optional column list
	Comment        string
	PartitionBy    *PartitionDesc
	DistributedBy  *DistributionDesc
	Properties     []*Property
	Query          Node // *SelectStmt or *SetOpStmt
	Loc            Loc
}

// Tag implements Node.
func (n *CreateMTMVStmt) Tag() NodeTag { return T_CreateMTMVStmt }

var _ Node = (*CreateMTMVStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// AlterMTMVStmt represents:
//
//	ALTER MATERIALIZED VIEW mv_name
//	    { RENAME new_name
//	    | REFRESH method
//	    | REPLACE WITH MATERIALIZED VIEW other_mv
//	    | SET PROPERTIES (...) }
type AlterMTMVStmt struct {
	Name          *ObjectName
	NewName       *ObjectName // for RENAME
	RefreshMethod string      // for REFRESH method change
	Replace       bool        // REPLACE WITH MATERIALIZED VIEW
	ReplaceTarget *ObjectName // the other MV for REPLACE WITH
	Properties    []*Property // for SET PROPERTIES
	Loc           Loc
}

// Tag implements Node.
func (n *AlterMTMVStmt) Tag() NodeTag { return T_AlterMTMVStmt }

var _ Node = (*AlterMTMVStmt)(nil)

// ---------------------------------------------------------------------------
// DROP MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// DropMTMVStmt represents:
//
//	DROP MATERIALIZED VIEW [IF EXISTS] mv_name [ON base_table]
type DropMTMVStmt struct {
	Name     *ObjectName
	IfExists bool
	OnBase   *ObjectName // for sync MV: ON base_table
	Loc      Loc
}

// Tag implements Node.
func (n *DropMTMVStmt) Tag() NodeTag { return T_DropMTMVStmt }

var _ Node = (*DropMTMVStmt)(nil)

// ---------------------------------------------------------------------------
// REFRESH MATERIALIZED VIEW
// ---------------------------------------------------------------------------

// RefreshMTMVStmt represents:
//
//	REFRESH MATERIALIZED VIEW mv_name [COMPLETE | AUTO | PARTITIONS(p1, p2, ...)]
type RefreshMTMVStmt struct {
	Name       *ObjectName
	Mode       string   // "COMPLETE", "AUTO", or ""
	Partitions []string // partition names when PARTITIONS(...) specified
	Loc        Loc
}

// Tag implements Node.
func (n *RefreshMTMVStmt) Tag() NodeTag { return T_RefreshMTMVStmt }

var _ Node = (*RefreshMTMVStmt)(nil)

// ---------------------------------------------------------------------------
// PAUSE / RESUME MATERIALIZED VIEW JOB
// ---------------------------------------------------------------------------

// PauseMTMVJobStmt represents:
//
//	PAUSE MATERIALIZED VIEW JOB ON mv_name
type PauseMTMVJobStmt struct {
	Name *ObjectName // the MV name after ON
	Loc  Loc
}

// Tag implements Node.
func (n *PauseMTMVJobStmt) Tag() NodeTag { return T_PauseMTMVJobStmt }

var _ Node = (*PauseMTMVJobStmt)(nil)

// ResumeMTMVJobStmt represents:
//
//	RESUME MATERIALIZED VIEW JOB ON mv_name
type ResumeMTMVJobStmt struct {
	Name *ObjectName // the MV name after ON
	Loc  Loc
}

// Tag implements Node.
func (n *ResumeMTMVJobStmt) Tag() NodeTag { return T_ResumeMTMVJobStmt }

var _ Node = (*ResumeMTMVJobStmt)(nil)

// ---------------------------------------------------------------------------
// CANCEL MATERIALIZED VIEW TASK
// ---------------------------------------------------------------------------

// CancelMTMVTaskStmt represents:
//
//	CANCEL MATERIALIZED VIEW TASK task_id ON mv_name
type CancelMTMVTaskStmt struct {
	TaskID int64
	Name   *ObjectName // the MV name after ON
	Loc    Loc
}

// Tag implements Node.
func (n *CancelMTMVTaskStmt) Tag() NodeTag { return T_CancelMTMVTaskStmt }

var _ Node = (*CancelMTMVTaskStmt)(nil)
