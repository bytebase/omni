package ast

// This file holds the concrete GoogleSQL parse-tree node types. The ast-core
// foundation ships only the File root container; later migration nodes
// (identifiers, types, expressions, SELECT core, joins, set-ops, DML, DDL,
// etc.) populate the rest, following the ZetaSQL-shaped tree.
//
// The cmd/genwalker code generator scans this file together with node.go to
// produce walk_generated.go.

// File is the root node of a parsed GoogleSQL source file. It holds the
// top-level statement list and the byte range covering the entire file.
// The parser entry point returns *File from Parse.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (f *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// DCL — GRANT / REVOKE (parser-dcl node)
// ---------------------------------------------------------------------------
//
// These node types model the legacy ANTLR grant_statement / revoke_statement
// family verbatim (GoogleSQLParser.g4 §2.8), which is itself a hand-port of
// Google's open-source ZetaSQL reference and the grammar bytebase consumes
// today. The grammar is:
//
//	grant_statement:  GRANT  privileges ON (identifier identifier?)? path_expression TO   grantee_list
//	revoke_statement: REVOKE privileges ON (identifier identifier?)? path_expression FROM grantee_list
//	privileges:       ALL PRIVILEGES? | privilege_list
//	privilege_list:   privilege (',' privilege)*
//	privilege:        privilege_name ('(' path_expression_list ')')?
//	privilege_name:   identifier | SELECT
//	grantee_list:     string_literal_or_parameter (',' string_literal_or_parameter)*
//
// The canonical ZetaSQL corpus (parser/googlesql/examples/.../grant_and_revoke.sql)
// is the accept/reject oracle for these forms. NOTE — DCL is a dialect-divergent
// zone: the Spanner emulator speaks a DIFFERENT GRANT dialect (`GRANT priv ON
// TABLE x TO ROLE r`, role-name grantees) and is therefore NON-AUTHORITATIVE
// here; see the parser-dcl divergence ledger entry. These nodes follow the
// ZetaSQL/.g4 grammar, not the Spanner emulator.
//
// Object names / privilege-column paths are kept as the lightweight NamePath
// value type (below) rather than a Node so they do NOT pre-empt the canonical
// path-expression node the expressions DAG node will own; DCL is self-contained
// (it depends only on the parser foundation).

// NamePath is a dotted identifier path (the grammar's path_expression:
// identifier ('.' identifier)*), used for a GRANT/REVOKE object name and for the
// optional per-privilege column list. It is a plain value type, NOT a Node: the
// walker does not descend into it and it intentionally does not collide with the
// canonical path-expression node owned by the expressions node.
//
// Parts holds each name component as the lexer surfaced it: a `backtick`-quoted
// identifier already has its backticks stripped, an unquoted identifier or
// keyword-as-identifier is its source word. (The lexer normalizes backtick
// identifiers the same way the legacy bigquery/spanner query-span extractors
// do.) Callers that need the original source slice can use Loc; DCL is not a
// query-span consumer, so no raw spelling is retained.
type NamePath struct {
	Parts []string // component names (>= 1)
	Loc   Loc
}

// String renders the path by joining its normalized parts with '.'. It does NOT
// re-quote; it is a debug/representation helper, not a deparser.
func (p NamePath) String() string {
	out := ""
	for i, part := range p.Parts {
		if i > 0 {
			out += "."
		}
		out += part
	}
	return out
}

// GranteeKind discriminates the four shapes of the grammar's
// string_literal_or_parameter (the entries of a grantee_list).
type GranteeKind int

const (
	// GranteeString is a string-literal grantee, e.g. 'user@google.com'.
	// Value holds the unquoted string content.
	GranteeString GranteeKind = iota
	// GranteeNamedParameter is a named query parameter, e.g. @user1 (the
	// grammar's named_parameter_expression: '@' identifier). Value holds the
	// parameter name without the leading '@'.
	GranteeNamedParameter
	// GranteePositionalParameter is a '?' positional query parameter (the
	// grammar's QUESTION). Value is empty.
	GranteePositionalParameter
	// GranteeSystemVariable is a system variable, e.g. @@user2 (the grammar's
	// system_variable_expression: '@@' path_expression). Value holds the
	// dotted variable path without the leading '@@'.
	GranteeSystemVariable
)

// String returns a human-readable name for the grantee kind.
func (k GranteeKind) String() string {
	switch k {
	case GranteeString:
		return "STRING"
	case GranteeNamedParameter:
		return "NAMED_PARAMETER"
	case GranteePositionalParameter:
		return "POSITIONAL_PARAMETER"
	case GranteeSystemVariable:
		return "SYSTEM_VARIABLE"
	default:
		return "UNKNOWN"
	}
}

// Grantee is one entry of a GRANT's TO list or a REVOKE's FROM list — the
// grammar's string_literal_or_parameter. Kind selects the shape; Value carries
// the payload (string content / parameter name / system-variable path), empty
// for a positional '?' parameter.
type Grantee struct {
	Kind  GranteeKind
	Value string
	Loc   Loc
}

// Tag implements Node.
func (n *Grantee) Tag() NodeTag { return T_Grantee }

// Privilege is one entry of a privilege_list: a privilege name with an optional
// parenthesized column list (the grammar's privilege:
// privilege_name path_expression_list_with_parens?). It is used only when a
// GRANT/REVOKE names explicit privileges; an `ALL PRIVILEGES` grant carries no
// Privilege nodes (see GrantStmt.AllPrivileges).
//
// Name is the privilege name (privilege_name: identifier | SELECT), as the lexer
// surfaced it (backtick identifiers unquoted). Columns is the optional
// per-privilege column list (e.g. the `(col1, col2)` in `insert(col1, col2)`);
// it is nil when no parentheses are present, and otherwise has >= 1 entry (the
// grammar's path_expression_list requires at least one path).
type Privilege struct {
	Name    string
	Columns []NamePath // nil when no '(' ... ')' column list was given
	Loc     Loc
}

// Tag implements Node.
func (n *Privilege) Tag() NodeTag { return T_Privilege }

// GrantStmt is a GRANT statement (grammar: grant_statement). It is the union of
// the documented forms; AllPrivileges selects ALL [PRIVILEGES] vs an explicit
// Privileges list.
//
//	GRANT { ALL [PRIVILEGES] | priv [, priv ...] }
//	  ON [ <object_type> [ <object_subtype> ] ] <path>
//	  TO <grantee> [, <grantee> ...]
//
// ObjectType holds the optional 0/1/2 object-type words (the grammar's
// (identifier identifier?)?), e.g. [] / ["table"] / ["materialized","view"].
// Path is the object name (a dotted path_expression). Grantees is the non-empty
// TO list.
type GrantStmt struct {
	AllPrivileges bool         // ALL [PRIVILEGES]
	Privileges    []*Privilege // explicit privilege list; nil when AllPrivileges
	ObjectType    []string     // 0, 1, or 2 object-type words
	Path          NamePath     // the ON object name
	Grantees      []*Grantee   // the TO recipients (>= 1)
	Loc           Loc
}

// Tag implements Node.
func (n *GrantStmt) Tag() NodeTag { return T_GrantStmt }

// RevokeStmt is a REVOKE statement (grammar: revoke_statement). It mirrors
// GrantStmt with a FROM grantee list instead of TO:
//
//	REVOKE { ALL [PRIVILEGES] | priv [, priv ...] }
//	  ON [ <object_type> [ <object_subtype> ] ] <path>
//	  FROM <grantee> [, <grantee> ...]
type RevokeStmt struct {
	AllPrivileges bool         // ALL [PRIVILEGES]
	Privileges    []*Privilege // explicit privilege list; nil when AllPrivileges
	ObjectType    []string     // 0, 1, or 2 object-type words
	Path          NamePath     // the ON object name
	Grantees      []*Grantee   // the FROM subjects (>= 1)
	Loc           Loc
}

// Tag implements Node.
func (n *RevokeStmt) Tag() NodeTag { return T_RevokeStmt }

// Compile-time assertions that the DCL node types satisfy Node.
var (
	_ Node = (*GrantStmt)(nil)
	_ Node = (*RevokeStmt)(nil)
	_ Node = (*Privilege)(nil)
	_ Node = (*Grantee)(nil)
)
