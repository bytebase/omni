package ast

// ---------------------------------------------------------------------------
// ADMIN statements (T7.4)
// ---------------------------------------------------------------------------

// AdminStmt is a generic node for ADMIN ... variants.
// Doris has many ADMIN commands; this node captures the leading verb+object
// and stores remaining variant-specific state in structured fields or as raw
// text in Args.
type AdminStmt struct {
	// Verb is the first keyword after ADMIN: "SHOW", "REBALANCE", "CANCEL",
	// "DIAGNOSE", "COMPACT", "CHECK", "REPAIR", "SET", "CLEAN", "COPY",
	// "DECOMMISSION", etc.
	Verb string

	// Object is the second keyword: "REPLICA", "TABLET", "FRONTEND", "BACKEND",
	// "CONFIG", "TRASH", "DATA", "PARTITION", "TABLE", "DISK", etc.
	Object string

	// Target is the optional table / tablet target (e.g., FROM table).
	Target *ObjectName

	// Properties holds any PROPERTIES(...) clause.
	Properties []*Property

	// Args holds raw text for variant-specific tokens that are not individually
	// structured (best-effort capture).
	Args string

	Loc Loc
}

// Tag implements Node.
func (n *AdminStmt) Tag() NodeTag { return T_AdminStmt }

// Compile-time assertion that *AdminStmt satisfies Node.
var _ Node = (*AdminStmt)(nil)

// ---------------------------------------------------------------------------
// ALTER SYSTEM statements (T7.4)
// ---------------------------------------------------------------------------

// SystemAlterStmt represents ALTER SYSTEM ADD/DROP/DECOMMISSION/MODIFY BACKEND/
// FRONTEND/OBSERVER/FOLLOWER/BROKER/... and SET LOAD ERROR HUB.
type SystemAlterStmt struct {
	// Verb is the operation keyword: "ADD", "DROP", "DECOMMISSION", "MODIFY",
	// "SET", etc.
	Verb string

	// Object is the target object type: "BACKEND", "FRONTEND", "OBSERVER",
	// "FOLLOWER", "BROKER", "LOAD ERROR HUB", etc.
	Object string

	// Hosts is the list of host:port strings (for BACKEND / FRONTEND variants).
	Hosts []string

	// BrokerName holds the broker name for BROKER operations.
	BrokerName string

	// DropAll is true for "DROP ALL BROKER broker_name".
	DropAll bool

	// SetClause holds key=value pairs passed after the SET keyword in
	// "MODIFY BACKEND ... SET (...)".
	SetClause []*Property

	// Properties holds any PROPERTIES(...) clause.
	Properties []*Property

	// Args holds raw text fallback for unstructured trailing tokens.
	Args string

	Loc Loc
}

// Tag implements Node.
func (n *SystemAlterStmt) Tag() NodeTag { return T_SystemAlterStmt }

// Compile-time assertion that *SystemAlterStmt satisfies Node.
var _ Node = (*SystemAlterStmt)(nil)

// ---------------------------------------------------------------------------
// CANCEL DECOMMISSION BACKEND (T7.4)
// ---------------------------------------------------------------------------

// CancelDecommissionStmt represents CANCEL DECOMMISSION BACKEND 'host:port' [, ...].
type CancelDecommissionStmt struct {
	// Object is always "BACKEND" for the currently supported form.
	Object string

	// Hosts is the list of host:port strings.
	Hosts []string

	Loc Loc
}

// Tag implements Node.
func (n *CancelDecommissionStmt) Tag() NodeTag { return T_CancelDecommissionStmt }

// Compile-time assertion that *CancelDecommissionStmt satisfies Node.
var _ Node = (*CancelDecommissionStmt)(nil)
