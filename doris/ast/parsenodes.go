package ast

import "strings"

// Property represents a key=value pair used in PROPERTIES(...) clauses.
type Property struct {
	Key   string
	Value string
	Loc   Loc
}

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
// Parts stores the normalized identifier text: the parser strips backtick
// quoting at parse time, so all parts are bare names. Quoting information
// is not preserved because Doris identifier resolution is case-insensitive
// regardless of quoting (unlike Snowflake where quoting affects case folding).
// If a future consumer needs to distinguish quoted vs unquoted identifiers,
// a Quoted []bool field can be added alongside Parts.
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

// ---------------------------------------------------------------------------
// Data type nodes
// ---------------------------------------------------------------------------

// TypeName represents a Doris data type as it appears in SQL source.
//
// Examples:
//
//	INT                         → Name="INT", Params=nil
//	VARCHAR(255)                → Name="VARCHAR", Params=[255]
//	DECIMAL(10,2)               → Name="DECIMAL", Params=[10,2]
//	ARRAY<INT>                  → Name="ARRAY", ElementType=&TypeName{Name:"INT"}
//	MAP<STRING,INT>             → Name="MAP", ElementType=&TypeName{Name:"STRING"}, ValueType=&TypeName{Name:"INT"}
//	STRUCT<name VARCHAR(50)>    → Name="STRUCT", Fields=[{Name:"name", Type:...}]
type TypeName struct {
	Name        string         // source text of type name: "INT", "VARCHAR", "ARRAY", etc.
	Params      []int          // numeric params like (10) or (10,2); nil if absent
	ElementType *TypeName      // for ARRAY<T>: element type; for MAP<K,V>: key type
	ValueType   *TypeName      // for MAP<K,V>: value type
	Fields      []*StructField // for STRUCT<name type, ...>: field list
	Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion that *TypeName satisfies Node.
var _ Node = (*TypeName)(nil)

// StructField represents one named field in a STRUCT type: `name dataType`.
type StructField struct {
	Name string
	Type *TypeName
	Loc  Loc
}

// ---------------------------------------------------------------------------
// Index DDL nodes (T2.5)
// ---------------------------------------------------------------------------

// CreateIndexStmt represents:
//
//	CREATE INDEX [IF NOT EXISTS] name ON table (col1, col2, ...)
//	  [USING index_type]
//	  [COMMENT 'comment']
//	  [PROPERTIES("key"="value", ...)]
//
// index_type is one of: BITMAP, NGRAM_BF, INVERTED, ANN, or empty string.
type CreateIndexStmt struct {
	Name        *ObjectName
	Table       *ObjectName
	Columns     []string    // column names (bare, stripped of backticks)
	IfNotExists bool
	IndexType   string      // BITMAP, NGRAM_BF, INVERTED, ANN, or ""
	Properties  []*Property // optional PROPERTIES(...)
	Comment     string      // optional COMMENT 'text'
	Loc         Loc
}

// Tag implements Node.
func (n *CreateIndexStmt) Tag() NodeTag { return T_CreateIndexStmt }

// Compile-time assertion that *CreateIndexStmt satisfies Node.
var _ Node = (*CreateIndexStmt)(nil)

// DropIndexStmt represents:
//
//	DROP INDEX [IF EXISTS] name ON table
type DropIndexStmt struct {
	Name     *ObjectName
	Table    *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropIndexStmt) Tag() NodeTag { return T_DropIndexStmt }

// Compile-time assertion that *DropIndexStmt satisfies Node.
var _ Node = (*DropIndexStmt)(nil)

// BuildIndexStmt represents:
//
//	BUILD INDEX name ON table [PARTITIONS(p1, p2, ...)]
type BuildIndexStmt struct {
	Name       *ObjectName
	Table      *ObjectName
	Partitions []string // optional partition names
	Loc        Loc
}

// Tag implements Node.
func (n *BuildIndexStmt) Tag() NodeTag { return T_BuildIndexStmt }

// Compile-time assertion that *BuildIndexStmt satisfies Node.
var _ Node = (*BuildIndexStmt)(nil)
