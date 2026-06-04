package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under googlesql/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable — do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// DCL — GRANT / REVOKE (parser-dcl node).
	//
	// These mirror the legacy ANTLR grant_statement / revoke_statement family
	// (GoogleSQLParser.g4): one statement node per verb plus the shared
	// privilege and grantee leaf nodes.
	T_GrantStmt
	T_RevokeStmt
	T_Privilege
	T_Grantee
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_GrantStmt:
		return "GrantStmt"
	case T_RevokeStmt:
		return "RevokeStmt"
	case T_Privilege:
		return "Privilege"
	case T_Grantee:
		return "Grantee"
	default:
		return "Unknown"
	}
}
