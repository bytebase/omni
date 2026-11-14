package ast

import "strings"

// This file holds the concrete Doris parse-tree node types. F1 ships only
// the File root container and ObjectName; later migration nodes (T1.1+)
// populate the rest.

// ---------------------------------------------------------------------------
// DDL — Database / Schema statements
// ---------------------------------------------------------------------------

// CreateDatabaseStmt represents CREATE DATABASE [IF NOT EXISTS] name [PROPERTIES(...)].
// CREATE SCHEMA is an alias and produces the same node.
type CreateDatabaseStmt struct {
	Name        *ObjectName
	IfNotExists bool
	Properties  []*Property // optional PROPERTIES clause
	Loc         Loc
}

// Tag implements Node.
func (n *CreateDatabaseStmt) Tag() NodeTag { return T_CreateDatabaseStmt }

// Compile-time assertion that *CreateDatabaseStmt satisfies Node.
var _ Node = (*CreateDatabaseStmt)(nil)

// AlterDatabaseStmt represents ALTER DATABASE name SET PROPERTIES(...) | RENAME name | SET QUOTA ...
type AlterDatabaseStmt struct {
	Name       *ObjectName
	NewName    *ObjectName // for RENAME
	Properties []*Property // for SET PROPERTIES
	QuotaType  string      // DATA, REPLICA or TRANSACTION for SET ... QUOTA
	Quota      string      // the quota of SET ... QUOTA
	Loc        Loc
}

// Tag implements Node.
func (n *AlterDatabaseStmt) Tag() NodeTag { return T_AlterDatabaseStmt }

// Compile-time assertion that *AlterDatabaseStmt satisfies Node.
var _ Node = (*AlterDatabaseStmt)(nil)

// DropDatabaseStmt represents DROP DATABASE [IF EXISTS] name [FORCE].
// DROP SCHEMA is an alias and produces the same node.
type DropDatabaseStmt struct {
	Name     *ObjectName
	IfExists bool
	Force    bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropDatabaseStmt) Tag() NodeTag { return T_DropDatabaseStmt }

// Compile-time assertion that *DropDatabaseStmt satisfies Node.
var _ Node = (*DropDatabaseStmt)(nil)

// DropTableStmt represents DROP [TEMPORARY] TABLE [IF EXISTS] name [FORCE].
type DropTableStmt struct {
	Name      *ObjectName
	IfExists  bool
	Force     bool
	Temporary bool
	Loc       Loc
}

// Tag implements Node.
func (n *DropTableStmt) Tag() NodeTag { return T_DropTableStmt }

// Compile-time assertion that *DropTableStmt satisfies Node.
var _ Node = (*DropTableStmt)(nil)

// CreateFunctionStmt represents
//
//	CREATE [GLOBAL | SESSION | LOCAL] [TABLES | AGGREGATE] FUNCTION [IF NOT EXISTS]
//	    name(arg_types) RETURNS type [INTERMEDIATE type] [PROPERTIES (...)]
//	CREATE [GLOBAL | SESSION | LOCAL] ALIAS FUNCTION [IF NOT EXISTS]
//	    name(arg_types) WITH PARAMETER (params) AS expr
//
// RawSignature holds the text after CREATE, verbatim.
type CreateFunctionStmt struct {
	Global       bool
	Aggregate    bool
	Tables       bool
	Alias        bool
	IfNotExists  bool
	Name         *ObjectName
	Args         []*TypeName // argument types; a trailing "..." sets Variadic
	Variadic     bool
	Returns      *TypeName   // nil for an ALIAS function
	Intermediate *TypeName   // INTERMEDIATE type, or nil
	Properties   []*Property // PROPERTIES (...)
	Parameters   []string    // ALIAS: WITH PARAMETER (...)
	Body         Node        // ALIAS: the expression after AS
	RawSignature string
	Loc          Loc
}

// Tag implements Node.
func (n *CreateFunctionStmt) Tag() NodeTag { return T_CreateFunctionStmt }

// Compile-time assertion that *CreateFunctionStmt satisfies Node.
var _ Node = (*CreateFunctionStmt)(nil)

// DropFunctionStmt represents
//
//	DROP [GLOBAL | SESSION | LOCAL] FUNCTION [IF EXISTS] name(arg_types).
type DropFunctionStmt struct {
	Name         *ObjectName
	Global       bool
	IfExists     bool
	Args         []*TypeName // argument types; a trailing "..." sets Variadic
	Variadic     bool
	RawSignature string // function name + argument-type list (verbatim)
	Loc          Loc
}

// Tag implements Node.
func (n *DropFunctionStmt) Tag() NodeTag { return T_DropFunctionStmt }

// Compile-time assertion that *DropFunctionStmt satisfies Node.
var _ Node = (*DropFunctionStmt)(nil)

// Property is a key=value pair used in PROPERTIES clauses.
type Property struct {
	Key   string
	Value string
	Loc   Loc
}

// Tag implements Node.
func (n *Property) Tag() NodeTag { return T_Property }

// Compile-time assertion that *Property satisfies Node.
var _ Node = (*Property)(nil)

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
//	STRUCT<name:VARCHAR(50)>    → Name="STRUCT", Fields=[{Name:"name", Type:...}]
//	VARIANT<'a':INT>            → Name="VARIANT", Fields=[{Name:"a", Type:...}]
//	AGG_STATE<sum(INT NULL)>    → Name="AGG_STATE", Func="sum", Args=[{Name:"INT", Nullable:"NULL"}]
type TypeName struct {
	Name        string         // source text of type name: "INT", "VARCHAR", "ARRAY", etc.
	Params      []int          // numeric params like (10) or (10,2); nil if absent
	Star        bool           // VARCHAR(*): the first parameter is '*'
	ElementType *TypeName      // for ARRAY<T>: element type; for MAP<K,V>: key type
	ValueType   *TypeName      // for MAP<K,V>: value type
	Fields      []*StructField // for STRUCT<name:type, ...> and VARIANT<'path':type, ...>
	Properties  []*Property    // for VARIANT<..., PROPERTIES(...)>
	Func        string         // for AGG_STATE<func(...)>: the function name
	Args        []*TypeName    // for AGG_STATE<func(...)>: the argument types
	Nullable    string         // an AGG_STATE argument's "NULL" or "NOT NULL", else ""
	Loc         Loc
}

// Tag implements Node.
func (n *TypeName) Tag() NodeTag { return T_TypeName }

// Compile-time assertion that *TypeName satisfies Node.
var _ Node = (*TypeName)(nil)

// StructField represents one named field in a STRUCT type, `name:dataType`,
// or a predefined sub-column of a VARIANT type, `'path':dataType`.
type StructField struct {
	Name      string
	Type      *TypeName
	Comment   string // COMMENT 'text', or ""
	MatchType string // VARIANT only: MATCH_NAME or MATCH_NAME_GLOB, or ""
	Loc       Loc
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
	Columns     []string // column names (bare, stripped of backticks)
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
