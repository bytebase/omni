package ast

import "strings"

// This file holds the concrete Doris parse-tree node types. F1 ships only
// the File root container and ObjectName; later migration nodes (T1.1+)
// populate the rest.

// File is the root node of a parsed Doris source file. It holds the
// top-level statement list and the byte range covering the entire file.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (n *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// Identifier types
// ---------------------------------------------------------------------------

// ObjectName represents a qualified multi-part identifier as used by Doris's
// multipartIdentifier grammar rule. Parts are stored in order: for a 3-part
// name like catalog.db.table, Parts = ["catalog", "db", "table"].
//
// Parts stores the raw text from source: for backtick-quoted identifiers,
// the content between the backticks; for unquoted identifiers, the source
// bytes with case preserved.
//
// ObjectName is a Node and is visited by the AST walker, but has no Node
// children to descend into.
type ObjectName struct {
	Parts []string
	Loc   Loc
}

// Tag implements Node.
func (n *ObjectName) Tag() NodeTag { return T_ObjectName }

// Compile-time assertion that *ObjectName satisfies Node.
var _ Node = (*ObjectName)(nil)

// String returns the dotted form of the name (e.g., "catalog.db.table").
func (n *ObjectName) String() string {
	return strings.Join(n.Parts, ".")
}
