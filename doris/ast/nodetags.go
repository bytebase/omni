package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under doris/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable -- do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified multi-part name
	// (e.g., catalog.db.table).
	T_ObjectName

	// T_TypeName is the tag for *TypeName, a Doris data type as it appears in
	// SQL source (e.g., INT, VARCHAR(255), ARRAY<INT>, MAP<STRING,INT>).
	T_TypeName

	// T_CreateIndexStmt is the tag for *CreateIndexStmt.
	T_CreateIndexStmt

	// T_DropIndexStmt is the tag for *DropIndexStmt.
	T_DropIndexStmt

	// T_BuildIndexStmt is the tag for *BuildIndexStmt.
	T_BuildIndexStmt

	// T_CreateDatabaseStmt is the tag for *CreateDatabaseStmt.
	T_CreateDatabaseStmt

	// T_AlterDatabaseStmt is the tag for *AlterDatabaseStmt.
	T_AlterDatabaseStmt

	// T_DropDatabaseStmt is the tag for *DropDatabaseStmt.
	T_DropDatabaseStmt

	// T_Property is the tag for *Property, a key=value pair in PROPERTIES clauses.
	T_Property
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_ObjectName:
		return "ObjectName"
	case T_TypeName:
		return "TypeName"
	case T_CreateIndexStmt:
		return "CreateIndexStmt"
	case T_DropIndexStmt:
		return "DropIndexStmt"
	case T_BuildIndexStmt:
		return "BuildIndexStmt"
	case T_CreateDatabaseStmt:
		return "CreateDatabaseStmt"
	case T_AlterDatabaseStmt:
		return "AlterDatabaseStmt"
	case T_DropDatabaseStmt:
		return "DropDatabaseStmt"
	case T_Property:
		return "Property"
	default:
		return "Unknown"
	}
}
