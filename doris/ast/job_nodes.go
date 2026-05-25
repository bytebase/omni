package ast

// This file holds AST node types for Doris JOB DDL (T8.1).
//
// Doris supports scheduled jobs via the JOB framework:
//   - CREATE JOB ... ON SCHEDULE ... DO statement
//   - ALTER JOB
//   - DROP JOB
//   - PAUSE JOB
//   - RESUME JOB
//   - CANCEL TASK FOR job_name task_id
//   - SHOW JOB / SHOW JOB TASK FOR job_name

// ---------------------------------------------------------------------------
// JobSchedule — ON SCHEDULE clause
// ---------------------------------------------------------------------------

// JobSchedule represents the ON SCHEDULE clause in CREATE JOB:
//
//	ON SCHEDULE { EVERY interval [STARTS timestamp] [ENDS timestamp] | AT timestamp }
type JobSchedule struct {
	Every  string // EVERY interval text (e.g. "1 HOUR"), empty when AT is used
	At     string // AT timestamp string literal value, empty when EVERY is used
	Starts string // optional STARTS timestamp string literal value
	Ends   string // optional ENDS timestamp string literal value
	Loc    Loc
}

// Tag implements Node.
func (n *JobSchedule) Tag() NodeTag { return T_JobSchedule }

var _ Node = (*JobSchedule)(nil)

// ---------------------------------------------------------------------------
// CREATE JOB
// ---------------------------------------------------------------------------

// CreateJobStmt represents:
//
//	CREATE JOB [IF NOT EXISTS] job_name
//	    ON SCHEDULE { EVERY interval [STARTS ts] [ENDS ts] | AT timestamp }
//	    [COMMENT 'text']
//	    DO statement
//
// or the streaming variant:
//
//	CREATE JOB [IF NOT EXISTS] job_name AS STREAMING
//	    [COMMENT 'text']
//	    DO statement
type CreateJobStmt struct {
	Name        *ObjectName
	IfNotExists bool
	JobType     string       // "STREAMING" or "SCHEDULE" (default when ON SCHEDULE used)
	Schedule    *JobSchedule // non-nil for SCHEDULE type
	Comment     string
	DoStmt      Node // the DO statement (nested parsed statement)
	Loc         Loc
}

// Tag implements Node.
func (n *CreateJobStmt) Tag() NodeTag { return T_CreateJobStmt }

var _ Node = (*CreateJobStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER JOB
// ---------------------------------------------------------------------------

// AlterJobStmt represents:
//
//	ALTER JOB job_name { PROPERTIES(...) | DO statement }
type AlterJobStmt struct {
	Name       *ObjectName
	Properties []*Property
	NewStmt    Node // optional new DO statement
	Loc        Loc
}

// Tag implements Node.
func (n *AlterJobStmt) Tag() NodeTag { return T_AlterJobStmt }

var _ Node = (*AlterJobStmt)(nil)

// ---------------------------------------------------------------------------
// DROP JOB
// ---------------------------------------------------------------------------

// DropJobStmt represents:
//
//	DROP JOB [IF EXISTS] job_name
//	DROP JOB WHERE expr
type DropJobStmt struct {
	Name     *ObjectName // non-nil when dropping by name
	IfExists bool
	Where    Node // non-nil when dropping via WHERE clause
	Loc      Loc
}

// Tag implements Node.
func (n *DropJobStmt) Tag() NodeTag { return T_DropJobStmt }

var _ Node = (*DropJobStmt)(nil)

// ---------------------------------------------------------------------------
// PAUSE JOB
// ---------------------------------------------------------------------------

// PauseJobStmt represents:
//
//	PAUSE JOB job_name
//	PAUSE JOB WHERE expr
type PauseJobStmt struct {
	Name  *ObjectName // non-nil when pausing by name
	Where Node        // non-nil when pausing via WHERE clause
	Loc   Loc
}

// Tag implements Node.
func (n *PauseJobStmt) Tag() NodeTag { return T_PauseJobStmt }

var _ Node = (*PauseJobStmt)(nil)

// ---------------------------------------------------------------------------
// RESUME JOB
// ---------------------------------------------------------------------------

// ResumeJobStmt represents:
//
//	RESUME JOB job_name
//	RESUME JOB WHERE expr
type ResumeJobStmt struct {
	Name  *ObjectName // non-nil when resuming by name
	Where Node        // non-nil when resuming via WHERE clause
	Loc   Loc
}

// Tag implements Node.
func (n *ResumeJobStmt) Tag() NodeTag { return T_ResumeJobStmt }

var _ Node = (*ResumeJobStmt)(nil)

// ---------------------------------------------------------------------------
// CANCEL TASK
// ---------------------------------------------------------------------------

// CancelTaskStmt represents:
//
//	CANCEL TASK FOR job_name task_id
type CancelTaskStmt struct {
	For    *ObjectName // FOR job_name
	TaskID int64       // task id
	Loc    Loc
}

// Tag implements Node.
func (n *CancelTaskStmt) Tag() NodeTag { return T_CancelTaskStmt }

var _ Node = (*CancelTaskStmt)(nil)

// ---------------------------------------------------------------------------
// SHOW JOB / SHOW JOB TASK
// ---------------------------------------------------------------------------

// ShowJobStmt represents:
//
//	SHOW JOB [LIKE 'pat' | WHERE expr]
type ShowJobStmt struct {
	Like  string // LIKE pattern, empty if not specified
	Where Node   // WHERE expr, nil if not specified
	Loc   Loc
}

// Tag implements Node.
func (n *ShowJobStmt) Tag() NodeTag { return T_ShowJobStmt }

var _ Node = (*ShowJobStmt)(nil)

// ShowJobTaskStmt represents:
//
//	SHOW JOB TASK FOR job_name
type ShowJobTaskStmt struct {
	For *ObjectName // FOR job_name
	Loc Loc
}

// Tag implements Node.
func (n *ShowJobTaskStmt) Tag() NodeTag { return T_ShowJobTaskStmt }

var _ Node = (*ShowJobTaskStmt)(nil)
