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
	default:
		return "Unknown"
	}
}
