package ast

// This file holds AST node types for statistics, analyze, and constraint
// management statements (T8.3).

// ---------------------------------------------------------------------------
// ADD / DROP CONSTRAINT
// ---------------------------------------------------------------------------

// AddConstraintStmt represents:
//
//	ALTER TABLE name ADD CONSTRAINT cname PRIMARY KEY (cols)
//	ALTER TABLE name ADD CONSTRAINT cname UNIQUE (cols)
//	ALTER TABLE name ADD CONSTRAINT cname FOREIGN KEY (cols) REFERENCES ref_table (cols)
type AddConstraintStmt struct {
	Table      *ObjectName
	Name       string      // constraint name
	Type       string      // "PRIMARY KEY", "UNIQUE", "FOREIGN KEY"
	Columns    []string
	RefTable   *ObjectName // for FOREIGN KEY
	RefColumns []string    // for FOREIGN KEY
	Loc        Loc
}

// Tag implements Node.
func (n *AddConstraintStmt) Tag() NodeTag { return T_AddConstraintStmt }

var _ Node = (*AddConstraintStmt)(nil)

// DropConstraintStmt represents:
//
//	ALTER TABLE name DROP CONSTRAINT cname
type DropConstraintStmt struct {
	Table *ObjectName
	Name  string
	Loc   Loc
}

// Tag implements Node.
func (n *DropConstraintStmt) Tag() NodeTag { return T_DropConstraintStmt }

var _ Node = (*DropConstraintStmt)(nil)

// ShowConstraintsStmt represents:
//
//	SHOW CONSTRAINTS FROM table
type ShowConstraintsStmt struct {
	Table *ObjectName
	Loc   Loc
}

// Tag implements Node.
func (n *ShowConstraintsStmt) Tag() NodeTag { return T_ShowConstraintsStmt }

var _ Node = (*ShowConstraintsStmt)(nil)

// ---------------------------------------------------------------------------
// ANALYZE
// ---------------------------------------------------------------------------

// AnalyzeStmt represents:
//
//	ANALYZE DATABASE name
//	ANALYZE TABLE name [(col1, col2, ...)]
//	ANALYZE PROFILE
//
// followed by optional modifiers:
//
//	WITH SAMPLE PERCENT n
//	WITH SAMPLE ROWS n
//	WITH SYNC
//	WITH INCREMENTAL
//	PROPERTIES(...)
type AnalyzeStmt struct {
	TargetType string      // "DATABASE", "TABLE", "PROFILE", or ""
	Target     *ObjectName // nil for PROFILE
	Columns    []string    // ANALYZE TABLE t(col1, col2) — optional column list
	Properties []*Property // WITH ... or PROPERTIES(...)
	Loc        Loc
}

// Tag implements Node.
func (n *AnalyzeStmt) Tag() NodeTag { return T_AnalyzeStmt }

var _ Node = (*AnalyzeStmt)(nil)

// ---------------------------------------------------------------------------
// SHOW ANALYZE / SHOW STATS
// ---------------------------------------------------------------------------

// ShowAnalyzeStmt represents:
//
//	SHOW [ALL | QUEUED] ANALYZE [JOB] [job_id | FOR table | LIKE pat]
//	SHOW ANALYZE TASK STATUS job_id
type ShowAnalyzeStmt struct {
	All    bool        // SHOW ALL ANALYZE
	Queued bool        // SHOW QUEUED ANALYZE JOBS
	IsTask bool        // SHOW ANALYZE TASK STATUS
	JobID  int64       // numeric job id (0 if absent)
	For    *ObjectName // FOR table
	Like   string      // LIKE 'pat'
	Where  Node        // WHERE expr
	Loc    Loc
}

// Tag implements Node.
func (n *ShowAnalyzeStmt) Tag() NodeTag { return T_ShowAnalyzeStmt }

var _ Node = (*ShowAnalyzeStmt)(nil)

// ShowStatsStmt represents:
//
//	SHOW [COLUMN | TABLE | INDEX | PARTITION] STATS [target] [args...]
type ShowStatsStmt struct {
	Type    string      // "COLUMN", "TABLE", "INDEX", "PARTITION", or ""
	Target  *ObjectName // optional table reference
	Where   Node        // optional WHERE clause
	Loc     Loc
}

// Tag implements Node.
func (n *ShowStatsStmt) Tag() NodeTag { return T_ShowStatsStmt }

var _ Node = (*ShowStatsStmt)(nil)

// ---------------------------------------------------------------------------
// DROP STATS
// ---------------------------------------------------------------------------

// DropStatsStmt represents:
//
//	DROP STATS target [(col1, col2)]
//	DROP EXPIRED STATS target
//	DROP CACHED STATS target
type DropStatsStmt struct {
	Variant string      // "", "EXPIRED", "CACHED"
	Target  *ObjectName
	Columns []string    // optional column list for plain DROP STATS
	Loc     Loc
}

// Tag implements Node.
func (n *DropStatsStmt) Tag() NodeTag { return T_DropStatsStmt }

var _ Node = (*DropStatsStmt)(nil)

// ---------------------------------------------------------------------------
// KILL ANALYZE
// ---------------------------------------------------------------------------

// KillAnalyzeStmt represents:
//
//	KILL ANALYZE job_id
type KillAnalyzeStmt struct {
	JobID int64
	Loc   Loc
}

// Tag implements Node.
func (n *KillAnalyzeStmt) Tag() NodeTag { return T_KillAnalyzeStmt }

var _ Node = (*KillAnalyzeStmt)(nil)
