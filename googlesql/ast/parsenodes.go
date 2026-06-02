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
