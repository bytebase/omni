package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under trino/ast must declare a unique
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

	// T_Identifier is the tag for *Identifier, a single Trino identifier
	// (unquoted, "double-quoted", `backtick-quoted`, or digit-leading).
	T_Identifier

	// T_QualifiedName is the tag for *QualifiedName, a dot-separated chain
	// of identifiers (e.g., catalog.schema.table).
	T_QualifiedName

	// --- statement nodes (parser-utility DAG node) ---
	//
	// These tags are returned by statement structs that live in package
	// parser, not here: a statement embeds parser-local Expr / DataType
	// children (which the ast package cannot import without a cycle), so the
	// structs stay in parser and satisfy ast.Node by returning the tag
	// constant declared below. Mirrors how expr.go / datatypes.go keep their
	// node types in package parser; only the closed tag set lives here.

	// T_ShowStmt is the tag for the SHOW family and DESCRIBE / DESC
	// (parser.ShowStmt): SHOW {TABLES|SCHEMAS|CATALOGS|COLUMNS|FUNCTIONS|
	// SESSION|STATS|GRANTS|ROLES|ROLE GRANTS|CREATE …} and DESCRIBE/DESC name.
	T_ShowStmt

	// T_UseStmt is the tag for parser.UseStmt: USE [catalog.]schema.
	T_UseStmt

	// T_SetSessionStmt is the tag for parser.SetSessionStmt:
	// SET SESSION name = expression.
	T_SetSessionStmt

	// T_ResetSessionStmt is the tag for parser.ResetSessionStmt:
	// RESET SESSION name.
	T_ResetSessionStmt

	// T_SetSessionAuthorizationStmt is the tag for
	// parser.SetSessionAuthorizationStmt: SET SESSION AUTHORIZATION user.
	T_SetSessionAuthorizationStmt

	// T_ResetSessionAuthorizationStmt is the tag for
	// parser.ResetSessionAuthorizationStmt: RESET SESSION AUTHORIZATION.
	T_ResetSessionAuthorizationStmt

	// T_SetRoleStmt is the tag for parser.SetRoleStmt:
	// SET ROLE (role | ALL | NONE) [IN catalog].
	T_SetRoleStmt

	// T_SetPathStmt is the tag for parser.SetPathStmt:
	// SET PATH pathElement (, pathElement)*.
	T_SetPathStmt

	// T_SetTimeZoneStmt is the tag for parser.SetTimeZoneStmt:
	// SET TIME ZONE (LOCAL | expression).
	T_SetTimeZoneStmt

	// T_ExplainStmt is the tag for parser.ExplainStmt:
	// EXPLAIN [ (option, …) ] statement and EXPLAIN ANALYZE [VERBOSE] statement.
	T_ExplainStmt

	// T_CallStmt is the tag for parser.CallStmt:
	// CALL name ( [name =>] expression, … ).
	T_CallStmt
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_Identifier:
		return "Identifier"
	case T_QualifiedName:
		return "QualifiedName"
	case T_ShowStmt:
		return "ShowStmt"
	case T_UseStmt:
		return "UseStmt"
	case T_SetSessionStmt:
		return "SetSessionStmt"
	case T_ResetSessionStmt:
		return "ResetSessionStmt"
	case T_SetSessionAuthorizationStmt:
		return "SetSessionAuthorizationStmt"
	case T_ResetSessionAuthorizationStmt:
		return "ResetSessionAuthorizationStmt"
	case T_SetRoleStmt:
		return "SetRoleStmt"
	case T_SetPathStmt:
		return "SetPathStmt"
	case T_SetTimeZoneStmt:
		return "SetTimeZoneStmt"
	case T_ExplainStmt:
		return "ExplainStmt"
	case T_CallStmt:
		return "CallStmt"
	default:
		return "Unknown"
	}
}
