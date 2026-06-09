package ast

// This file holds utility statement AST node types (T7.3):
//   - SHOW (40+ variants, generic + well-known forms)
//   - DESCRIBE / DESC
//   - EXPLAIN
//   - USE
//   - SET (generic variable assignment)
//   - UNSET
//   - HELP

// ---------------------------------------------------------------------------
// SHOW
// ---------------------------------------------------------------------------

// ShowStmt is a generic node for the many SHOW variants in Doris.
// Type identifies which SHOW variant (TABLES, DATABASES, COLUMNS, etc.).
// Target/From/Like/Where capture common modifiers.
// Args carries raw remaining tokens for variant-specific options.
type ShowStmt struct {
	// Type is the SHOW variant keyword(s), e.g. "TABLES", "DATABASES",
	// "COLUMNS", "CREATE TABLE", "GRANTS", "VARIABLES", "PARTITIONS", etc.
	Type     string
	Target   *ObjectName // optional target (e.g. SHOW COLUMNS FROM target)
	From     string      // optional FROM db name
	Like     string      // optional LIKE 'pattern'
	Where    Node        // optional WHERE expression
	Args     string      // raw text for variant-specific args (ORDER BY, LIMIT, etc.)
	Full     bool        // SHOW FULL TABLES / SHOW FULL COLUMNS
	Extended bool        // SHOW EXTENDED ...
	Loc      Loc
}

// Tag implements Node.
func (n *ShowStmt) Tag() NodeTag { return T_ShowStmt }

var _ Node = (*ShowStmt)(nil)

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

// DescribeStmt represents DESCRIBE / DESC statements.
//
//	DESC [FULL] table_name
//	DESCRIBE [FULL] table_name [ALL VERBOSE]
//	DESCRIBE FUNCTION name
type DescribeStmt struct {
	Target     *ObjectName
	Full       bool // DESCRIBE FULL ...
	AllVerbose bool // ALL VERBOSE suffix
	Loc        Loc
}

// Tag implements Node.
func (n *DescribeStmt) Tag() NodeTag { return T_DescribeStmt }

var _ Node = (*DescribeStmt)(nil)

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

// ExplainStmt represents EXPLAIN [modifier] query.
//
//	EXPLAIN [VERBOSE | GRAPH | PARSED | ANALYZED | REWRITTEN | PLAN |
//	         PHYSICAL | MEMO | SHAPE | DUMP | OPTIMIZED | PLAN PROCESS] query
type ExplainStmt struct {
	// Type is the optional explain modifier, e.g. "VERBOSE", "GRAPH", "PLAN", "".
	Type  string
	Query Node // the explained query (SelectStmt, InsertStmt, etc.)
	Loc   Loc
}

// Tag implements Node.
func (n *ExplainStmt) Tag() NodeTag { return T_ExplainStmt }

var _ Node = (*ExplainStmt)(nil)

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

// UseStmt represents a USE statement.
//
//	USE db_name
//	USE catalog_name.db_name
//	USE db_name@cluster_name
type UseStmt struct {
	Database string // database (right-hand of catalog, or sole name)
	Catalog  string // catalog prefix before '.' (empty if absent)
	Cluster  string // cluster suffix after '@' (empty if absent)
	Loc      Loc
}

// Tag implements Node.
func (n *UseStmt) Tag() NodeTag { return T_UseStmt }

var _ Node = (*UseStmt)(nil)

// ---------------------------------------------------------------------------
// SET (generic variable assignment)
// ---------------------------------------------------------------------------

// SetStmt represents a generic SET statement (variable assignments, NAMES,
// CHARSET, TRANSACTION). SET PASSWORD and SET DEFAULT STORAGE VAULT are
// handled by their own dedicated nodes.
//
//	SET [GLOBAL|SESSION|LOCAL] var = expr [, ...]
//	SET @@[GLOBAL.|SESSION.]var = expr
//	SET NAMES 'charset' [COLLATE 'collation']
//	SET CHARSET 'charset'
//	SET TRANSACTION { READ ONLY | READ WRITE | ISOLATION LEVEL ... }
type SetStmt struct {
	// Type identifies special SET forms: "VARIABLE" (default), "NAMES",
	// "CHARSET", "TRANSACTION".
	Type  string
	Items []*SetItem
	Loc   Loc
}

// Tag implements Node.
func (n *SetStmt) Tag() NodeTag { return T_SetStmt }

var _ Node = (*SetStmt)(nil)

// SetItem is one assignment in a SET statement.
type SetItem struct {
	// Scope is "GLOBAL", "SESSION", "LOCAL", or "" (session default).
	Scope string
	// Name is the variable name (without @@ prefix).
	Name  string
	Value Node   // parsed expression value; may be nil when Raw is set
	Raw   string // raw text for complex RHS
	Loc   Loc
}

// Tag implements Node.
func (n *SetItem) Tag() NodeTag { return T_SetItem }

var _ Node = (*SetItem)(nil)

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

// UnsetStmt represents an UNSET statement (variable clearing).
//
//	UNSET [GLOBAL|SESSION] VARIABLE name [, name ...]
//	UNSET [GLOBAL|SESSION] VARIABLE ALL
type UnsetStmt struct {
	// Type is "VARIABLE" (the only supported form so far).
	Type  string
	Scope string   // "GLOBAL", "SESSION", or ""
	Names []string // variable names; empty when All is true
	All   bool     // UNSET VARIABLE ALL (or UNSET VARIABLE *)
	Loc   Loc
}

// Tag implements Node.
func (n *UnsetStmt) Tag() NodeTag { return T_UnsetStmt }

var _ Node = (*UnsetStmt)(nil)

// ---------------------------------------------------------------------------
// HELP
// ---------------------------------------------------------------------------

// HelpStmt represents HELP 'mask'.
type HelpStmt struct {
	Mask string // the help topic / mask string
	Loc  Loc
}

// Tag implements Node.
func (n *HelpStmt) Tag() NodeTag { return T_HelpStmt }

var _ Node = (*HelpStmt)(nil)
