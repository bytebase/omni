package ast

// This file holds AST node types for StarRocks's scheduled-task statements
// (StarRocks.g4:625-642): SUBMIT TASK and DROP TASK. The Doris JOB statements
// in job_nodes.go are a different family, and starrocks/parser rejects them.

// TaskSchedule represents the SCHEDULE clause of SUBMIT TASK:
//
//	SCHEDULE (START '(' string ')')? EVERY '(' INTERVAL integer unit ')'
type TaskSchedule struct {
	Start    string // optional START('...') value; empty if absent
	Interval string // INTERVAL value, as written (an integer literal)
	Unit     string // DAY, HOUR, MINUTE or SECOND, upper-cased
	Loc      Loc
}

// Tag implements Node.
func (n *TaskSchedule) Tag() NodeTag { return T_TaskSchedule }

var _ Node = (*TaskSchedule)(nil)

// SubmitTaskStmt represents submitTaskStatement (StarRocks.g4:625):
//
//	SUBMIT TASK qualifiedName? taskClause*
//	    AS (createTableAsSelectStatement | insertStatement | dataCacheSelectStatement)
//
// taskClause is PROPERTIES(...) or taskScheduleDesc.
type SubmitTaskStmt struct {
	Name       *ObjectName   // optional task name; nil if absent
	Properties []*Property   // PROPERTIES(...) clauses, in source order
	Schedule   *TaskSchedule // optional SCHEDULE clause; nil if absent
	Statement  Node          // the submitted statement: CTAS, INSERT or CACHE SELECT
	Loc        Loc
}

// Tag implements Node.
func (n *SubmitTaskStmt) Tag() NodeTag { return T_SubmitTaskStmt }

var _ Node = (*SubmitTaskStmt)(nil)

// DropTaskStmt represents dropTaskStatement (StarRocks.g4:636):
//
//	DROP TASK qualifiedName FORCE?
type DropTaskStmt struct {
	Name  *ObjectName
	Force bool
	Loc   Loc
}

// Tag implements Node.
func (n *DropTaskStmt) Tag() NodeTag { return T_DropTaskStmt }

var _ Node = (*DropTaskStmt)(nil)
