package analysis

// TableRef is an engine-neutral schema-qualified table or object reference.
type TableRef struct {
	Schema string
	Name   string
}

// RenamePair captures the source and destination tables of a rename statement.
type RenamePair struct {
	Old TableRef
	New TableRef
}

// WalkThroughStatementKind classifies an AST statement for DDL walk-through
// precheck and cross-database rename target discovery.
type WalkThroughStatementKind int

const (
	// WalkThroughStatementOther is any statement not needing walk-through precheck.
	WalkThroughStatementOther WalkThroughStatementKind = iota
	// WalkThroughStatementCreateDatabase is a CREATE DATABASE / CREATE SCHEMA statement.
	WalkThroughStatementCreateDatabase
	// WalkThroughStatementAlterDatabase is an ALTER DATABASE / ALTER SCHEMA statement.
	WalkThroughStatementAlterDatabase
	// WalkThroughStatementDropDatabase is a DROP DATABASE / DROP SCHEMA statement.
	WalkThroughStatementDropDatabase
	// WalkThroughStatementCreateTable is a CREATE TABLE statement.
	WalkThroughStatementCreateTable
	// WalkThroughStatementDropTable is a DROP TABLE statement.
	WalkThroughStatementDropTable
	// WalkThroughStatementTableRefCheck is a DDL statement whose table or object
	// references must belong to the current database.
	WalkThroughStatementTableRefCheck
	// WalkThroughStatementRenameTable is a RENAME TABLE statement.
	WalkThroughStatementRenameTable
)

// String returns the canonical name of the walk-through statement kind.
func (k WalkThroughStatementKind) String() string {
	switch k {
	case WalkThroughStatementCreateDatabase:
		return "CREATE_DATABASE"
	case WalkThroughStatementAlterDatabase:
		return "ALTER_DATABASE"
	case WalkThroughStatementDropDatabase:
		return "DROP_DATABASE"
	case WalkThroughStatementCreateTable:
		return "CREATE_TABLE"
	case WalkThroughStatementDropTable:
		return "DROP_TABLE"
	case WalkThroughStatementTableRefCheck:
		return "TABLE_REF_CHECK"
	case WalkThroughStatementRenameTable:
		return "RENAME_TABLE"
	default:
		return "OTHER"
	}
}

// WalkThroughStatement holds the engine-neutral fields needed for walk-through
// precheck and cross-database rename target discovery. Source is optional
// caller-supplied SQL text for diagnostic formatting.
type WalkThroughStatement struct {
	Kind        WalkThroughStatementKind
	Source      string
	Database    string
	IfNotExists bool
	IsCTAS      bool
	TableRefs   []TableRef
	RenamePairs []RenamePair
}
