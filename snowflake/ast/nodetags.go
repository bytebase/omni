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

	T_Literal
	T_ColumnRef
	T_StarExpr
	T_BinaryExpr
	T_UnaryExpr
	T_ParenExpr
	T_CastExpr
	T_CaseExpr
	T_FuncCallExpr
	T_IffExpr
	T_CollateExpr
	T_IsExpr
	T_BetweenExpr
	T_InExpr
	T_LikeExpr
	T_AccessExpr
	T_ArrayLiteralExpr
	T_JsonLiteralExpr
	T_LambdaExpr
	T_SubqueryExpr
	T_ExistsExpr
	T_SelectStmt
	T_SetOperationStmt
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
	case T_Literal:
		return "Literal"
	case T_ColumnRef:
		return "ColumnRef"
	case T_StarExpr:
		return "StarExpr"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_ParenExpr:
		return "ParenExpr"
	case T_CastExpr:
		return "CastExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_FuncCallExpr:
		return "FuncCallExpr"
	case T_IffExpr:
		return "IffExpr"
	case T_CollateExpr:
		return "CollateExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_InExpr:
		return "InExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_AccessExpr:
		return "AccessExpr"
	case T_ArrayLiteralExpr:
		return "ArrayLiteralExpr"
	case T_JsonLiteralExpr:
		return "JsonLiteralExpr"
	case T_LambdaExpr:
		return "LambdaExpr"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
	case T_SelectStmt:
		return "SelectStmt"
	case T_SetOperationStmt:
		return "SetOperationStmt"
	default:
		return "Unknown"
	}
}
