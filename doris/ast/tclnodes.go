package ast

// TCL (Transaction Control Language) statement nodes.

// BeginStmt represents a BEGIN [WITH LABEL label_name] statement.
type BeginStmt struct {
	Label string // optional label name from WITH LABEL clause; empty if absent
	Loc   Loc
}

func (n *BeginStmt) Tag() NodeTag { return T_BeginStmt }

// CommitStmt represents a COMMIT [WORK] [AND [NO] CHAIN] [[NO] RELEASE] statement.
type CommitStmt struct {
	Work    bool   // true if WORK keyword was present
	Chain   string // "AND CHAIN", "AND NO CHAIN", or "" if absent
	Release string // "RELEASE", "NO RELEASE", or "" if absent
	Loc     Loc
}

func (n *CommitStmt) Tag() NodeTag { return T_CommitStmt }

// RollbackStmt represents a ROLLBACK [WORK] [AND [NO] CHAIN] [[NO] RELEASE] statement.
type RollbackStmt struct {
	Work    bool   // true if WORK keyword was present
	Chain   string // "AND CHAIN", "AND NO CHAIN", or "" if absent
	Release string // "RELEASE", "NO RELEASE", or "" if absent
	Loc     Loc
}

func (n *RollbackStmt) Tag() NodeTag { return T_RollbackStmt }

// StartTransactionStmt represents START TRANSACTION [WITH CONSISTENT SNAPSHOT].
type StartTransactionStmt struct {
	ConsistentSnapshot bool // WITH CONSISTENT SNAPSHOT
	Loc                Loc
}

// Tag implements Node.
func (n *StartTransactionStmt) Tag() NodeTag { return T_StartTransactionStmt }

var _ Node = (*StartTransactionStmt)(nil)
