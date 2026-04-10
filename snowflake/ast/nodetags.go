package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under snowflake/ast must declare a unique
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
	// returned by F4 (parser-entry).
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified 1/2/3-part name.
	T_ObjectName

	// T_TypeName is the tag for *TypeName, a data type reference.
	T_TypeName
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
	default:
		return "Unknown"
	}
}
