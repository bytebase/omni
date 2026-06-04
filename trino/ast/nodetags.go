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
	// --- DCL / TCL / prepared statement nodes (parser-dcl-tcl node) ---
	//
	// These concrete statement types live in package parser (grant_revoke.go,
	// transaction.go, prepared.go) — matching the Trino convention that
	// parser-package node types (Expr, DataType) are defined where they are
	// built — but they are returned from parseStmt as ast.Node, so they need a
	// tag here. The ast-core comment notes the tag set is grown by the parser
	// DAG nodes as they add statement types; these are the parser-dcl-tcl
	// node's additions.

	// T_GrantPrivStmt tags a GRANT <privileges> ON <target> TO <grantee>
	// statement (privilege grant).
	T_GrantPrivStmt
	// T_RevokePrivStmt tags a REVOKE [GRANT OPTION FOR] <privileges> ON
	// <target> FROM <grantee> statement (privilege revoke).
	T_RevokePrivStmt
	// T_DenyStmt tags a DENY <privileges> ON <target> TO <grantee> statement.
	T_DenyStmt
	// T_GrantRolesStmt tags a GRANT <roles> TO <principal> ... statement
	// (role grant).
	T_GrantRolesStmt
	// T_RevokeRolesStmt tags a REVOKE [ADMIN OPTION FOR] <roles> FROM
	// <principal> ... statement (role revoke).
	T_RevokeRolesStmt
	// T_CreateRoleStmt tags a CREATE ROLE statement.
	T_CreateRoleStmt
	// T_DropRoleStmt tags a DROP ROLE statement.
	T_DropRoleStmt

	// T_StartTransactionStmt tags a START TRANSACTION [mode, ...] statement.
	T_StartTransactionStmt
	// T_CommitStmt tags a COMMIT [WORK] statement.
	T_CommitStmt
	// T_RollbackStmt tags a ROLLBACK [WORK] statement.
	T_RollbackStmt

	// T_PrepareStmt tags a PREPARE <name> FROM <statement> statement.
	T_PrepareStmt
	// T_DeallocateStmt tags a DEALLOCATE PREPARE <name> statement.
	T_DeallocateStmt
	// T_ExecuteStmt tags an EXECUTE <name> [USING ...] statement.
	T_ExecuteStmt
	// T_ExecuteImmediateStmt tags an EXECUTE IMMEDIATE '<sql>' [USING ...]
	// statement.
	T_ExecuteImmediateStmt
	// T_DescribeInputStmt tags a DESCRIBE INPUT <name> statement.
	T_DescribeInputStmt
	// T_DescribeOutputStmt tags a DESCRIBE OUTPUT <name> statement.
	T_DescribeOutputStmt
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
	case T_GrantPrivStmt:
		return "GrantPrivStmt"
	case T_RevokePrivStmt:
		return "RevokePrivStmt"
	case T_DenyStmt:
		return "DenyStmt"
	case T_GrantRolesStmt:
		return "GrantRolesStmt"
	case T_RevokeRolesStmt:
		return "RevokeRolesStmt"
	case T_CreateRoleStmt:
		return "CreateRoleStmt"
	case T_DropRoleStmt:
		return "DropRoleStmt"
	case T_StartTransactionStmt:
		return "StartTransactionStmt"
	case T_CommitStmt:
		return "CommitStmt"
	case T_RollbackStmt:
		return "RollbackStmt"
	case T_PrepareStmt:
		return "PrepareStmt"
	case T_DeallocateStmt:
		return "DeallocateStmt"
	case T_ExecuteStmt:
		return "ExecuteStmt"
	case T_ExecuteImmediateStmt:
		return "ExecuteImmediateStmt"
	case T_DescribeInputStmt:
		return "DescribeInputStmt"
	case T_DescribeOutputStmt:
		return "DescribeOutputStmt"
	default:
		return "Unknown"
	}
}
