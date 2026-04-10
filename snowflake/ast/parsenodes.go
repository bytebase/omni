package ast

import "strings"

// This file holds the concrete snowflake parse-tree node types. F1 ships
// only the File root container; Tier 1+ migration nodes (identifiers,
// types, expressions, SELECT core, DDL, etc.) populate the rest.
//
// The cmd/genwalker code generator scans this file together with node.go
// to produce walk_generated.go.

// File is the root node of a parsed Snowflake source file. It holds the
// top-level statement list and the byte range covering the entire file.
// F4 (parser-entry) returns *File from Parse.
type File struct {
	Stmts []Node
	Loc   Loc
}

// Tag implements Node.
func (f *File) Tag() NodeTag { return T_File }

// Compile-time assertion that *File satisfies Node.
var _ Node = (*File)(nil)

// ---------------------------------------------------------------------------
// Identifier types
// ---------------------------------------------------------------------------

// Ident represents a single identifier — a name used to reference a database
// object (table, column, schema, etc.).
//
// Name is the raw text from source: for quoted identifiers, the content
// between the double-quotes with "" un-escaped; for unquoted identifiers,
// the source bytes with case preserved.
//
// Quoted reports whether the source used "..." quoting. This matters because
// Snowflake case-folds unquoted identifiers to uppercase at resolution time,
// while quoted identifiers preserve case.
//
// Ident is a value struct, NOT a Node. It is embedded by value in parent
// nodes (e.g. ObjectName) and is not visited by the AST walker.
//
// The zero value (Name == "" && Quoted == false) represents an absent
// identifier — used by ObjectName for unused parts (e.g. a 1-part name
// has zero Database and Schema).
type Ident struct {
	Name   string
	Quoted bool
	Loc    Loc
}

// Normalize returns the canonical form of the identifier per Snowflake
// resolution rules:
//   - Quoted identifiers: returned as-is (case-sensitive)
//   - Unquoted identifiers: uppercased
func (i Ident) Normalize() string {
	if i.Quoted {
		return i.Name
	}
	return strings.ToUpper(i.Name)
}

// String returns the source form of the identifier, re-quoting if it was
// originally quoted. Inner " characters are escaped as "". Useful for
// deparse and error messages.
func (i Ident) String() string {
	if !i.Quoted {
		return i.Name
	}
	return `"` + strings.ReplaceAll(i.Name, `"`, `""`) + `"`
}

// IsEmpty reports whether the identifier is the zero value (absent).
// Used by ObjectName to check whether a part (Database, Schema) is present.
func (i Ident) IsEmpty() bool {
	return i.Name == "" && !i.Quoted
}

// ObjectName represents a qualified object name (1/2/3-part) like
// `table`, `schema.table`, or `database.schema.table`.
//
// For 1-part names, only Name is set (Database and Schema are zero Idents).
// For 2-part names, Schema and Name are set.
// For 3-part names, all three are set.
//
// ObjectName is a Node and can be used as a child in the AST tree. The
// walker visits *ObjectName but does NOT descend into the embedded Ident
// fields (they are value structs, not Nodes).
type ObjectName struct {
	Database Ident // may be zero (IsEmpty)
	Schema   Ident // may be zero (IsEmpty)
	Name     Ident // always present for a valid ObjectName
	Loc      Loc
}

// Tag implements Node.
func (n *ObjectName) Tag() NodeTag { return T_ObjectName }

// Compile-time assertion that *ObjectName satisfies Node.
var _ Node = (*ObjectName)(nil)

// Normalize returns the canonical dotted form with each non-empty part
// normalized per Snowflake resolution rules.
//
// Examples:
//   - 1-part: "TABLE"
//   - 2-part: "SCHEMA.TABLE"
//   - 3-part: "DB.SCHEMA.TABLE"
func (n ObjectName) Normalize() string {
	parts := n.Parts()
	normalized := make([]string, len(parts))
	for i, p := range parts {
		normalized[i] = p.Normalize()
	}
	return strings.Join(normalized, ".")
}

// String returns the source form with dots. Each part is re-quoted if it
// was originally quoted.
func (n ObjectName) String() string {
	parts := n.Parts()
	strs := make([]string, len(parts))
	for i, p := range parts {
		strs[i] = p.String()
	}
	return strings.Join(strs, ".")
}

// Parts returns the non-empty parts in order. Length is 1, 2, or 3.
func (n ObjectName) Parts() []Ident {
	switch {
	case !n.Database.IsEmpty():
		return []Ident{n.Database, n.Schema, n.Name}
	case !n.Schema.IsEmpty():
		return []Ident{n.Schema, n.Name}
	default:
		return []Ident{n.Name}
	}
}

// Matches reports whether this ObjectName suffix-matches other using
// normalized (case-folded) comparison. A 1-part name matches any other
// with the same normalized Name. A 2-part name matches any other with
// the same normalized Schema + Name. A 3-part name requires all three
// parts to match.
func (n ObjectName) Matches(other ObjectName) bool {
	if n.Name.Normalize() != other.Name.Normalize() {
		return false
	}
	if n.Schema.IsEmpty() {
		return true // 1-part match
	}
	if n.Schema.Normalize() != other.Schema.Normalize() {
		return false
	}
	if n.Database.IsEmpty() {
		return true // 2-part match
	}
	return n.Database.Normalize() == other.Database.Normalize()
}

// ---------------------------------------------------------------------------
// Data type types
// ---------------------------------------------------------------------------

// TypeKind classifies Snowflake data types into categories for fast
// switch dispatch by downstream consumers. The Name field on TypeName
// carries the exact source text for round-tripping; Kind carries the
// category for semantic analysis.
type TypeKind int

const (
	TypeUnknown      TypeKind = iota
	TypeInt                   // INT, INTEGER, SMALLINT, TINYINT, BYTEINT, BIGINT
	TypeNumber                // NUMBER, NUMERIC, DECIMAL — may have (precision, scale)
	TypeFloat                 // FLOAT, FLOAT4, FLOAT8, DOUBLE, DOUBLE PRECISION, REAL
	TypeBoolean               // BOOLEAN
	TypeDate                  // DATE
	TypeDateTime              // DATETIME — may have (precision)
	TypeTime                  // TIME — may have (precision)
	TypeTimestamp             // TIMESTAMP — may have (precision)
	TypeTimestampLTZ          // TIMESTAMP_LTZ — may have (precision)
	TypeTimestampNTZ          // TIMESTAMP_NTZ — may have (precision)
	TypeTimestampTZ           // TIMESTAMP_TZ — may have (precision)
	TypeChar                  // CHAR, NCHAR, CHARACTER — may have (length)
	TypeVarchar               // VARCHAR, CHAR VARYING, NCHAR VARYING, NVARCHAR, NVARCHAR2, STRING, TEXT
	TypeBinary                // BINARY — may have (length)
	TypeVarbinary             // VARBINARY — may have (length)
	TypeVariant               // VARIANT
	TypeObject                // OBJECT
	TypeArray                 // ARRAY — may have ElementType
	TypeGeography             // GEOGRAPHY
	TypeGeometry              // GEOMETRY
	TypeVector                // VECTOR — has ElementType + VectorDim
)

// String returns the human-readable name of the TypeKind.
func (k TypeKind) String() string {
	switch k {
	case TypeUnknown:
		return "Unknown"
	case TypeInt:
		return "Int"
	case TypeNumber:
		return "Number"
	case TypeFloat:
		return "Float"
	case TypeBoolean:
		return "Boolean"
	case TypeDate:
		return "Date"
	case TypeDateTime:
		return "DateTime"
	case TypeTime:
		return "Time"
	case TypeTimestamp:
		return "Timestamp"
	case TypeTimestampLTZ:
		return "TimestampLTZ"
	case TypeTimestampNTZ:
		return "TimestampNTZ"
	case TypeTimestampTZ:
		return "TimestampTZ"
	case TypeChar:
		return "Char"
	case TypeVarchar:
		return "Varchar"
	case TypeBinary:
		return "Binary"
	case TypeVarbinary:
		return "Varbinary"
	case TypeVariant:
		return "Variant"
	case TypeObject:
		return "Object"
	case TypeArray:
		return "Array"
	case TypeGeography:
		return "Geography"
	case TypeGeometry:
		return "Geometry"
	case TypeVector:
		return "Vector"
	default:
		return "Unknown"
	}
}

// TypeName represents a Snowflake data type as it appears in SQL source.
//
// Examples:
//
//	INT                  → Kind=TypeInt, Name="INT", Params=nil
//	NUMBER(38, 0)        → Kind=TypeNumber, Name="NUMBER", Params=[38, 0]
//	VARCHAR(100)         → Kind=TypeVarchar, Name="VARCHAR", Params=[100]
//	TIMESTAMP_LTZ(9)     → Kind=TypeTimestampLTZ, Name="TIMESTAMP_LTZ", Params=[9]
//	DOUBLE PRECISION     → Kind=TypeFloat, Name="DOUBLE PRECISION", Params=nil
//	ARRAY(VARCHAR)       → Kind=TypeArray, Name="ARRAY", ElementType=&TypeName{...}
//	VECTOR(INT, 256)     → Kind=TypeVector, Name="VECTOR", ElementType=&TypeName{...}, VectorDim=256
//
// TypeName is a Node. The walker descends into ElementType when non-nil.
type TypeName struct {
	Kind        TypeKind  // classified type category
	Name        string    // source text of the type name for round-tripping
	Params      []int     // numeric type parameters; nil if absent
	ElementType *TypeName // element type for ARRAY and VECTOR; nil otherwise
	VectorDim   int       // dimension for VECTOR(type, dim); -1 if not VECTOR
	Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion that *TypeName satisfies Node.
var _ Node = (*TypeName)(nil)
