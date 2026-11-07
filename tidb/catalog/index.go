package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

type Index struct {
	Name     string
	Table    *Table
	Columns  []*IndexColumn
	Unique   bool
	Primary  bool
	Fulltext bool
	Spatial  bool
	// Vector marks TiDB's VECTOR INDEX, which SHOW CREATE TABLE writes with
	// that word in place of KEY (ConstraintVector in v8.5.5
	// pkg/parser/ast/ddl.go restores `VECTOR INDEX`).
	Vector bool
	// Global marks TiDB's GLOBAL index option, which exempts a unique key on
	// a partitioned table from covering the partitioning columns. v8.5.5
	// answers 8264 for a unique key that covers none of them without it, and
	// writes the one that has it as
	// "UNIQUE KEY `i` (`b`) /*T![global_index] GLOBAL */" (read on
	// 2026-11-04). LOCAL is the default and is written as nothing.
	Global       bool
	IndexType    string // BTREE, HASH, FULLTEXT, SPATIAL
	Comment      string
	Visible      bool
	KeyBlockSize int
	// addedByAlter marks a primary key ALTER TABLE added, which TiDB builds
	// NONCLUSTERED unless told CLUSTERED.
	addedByAlter bool
}

type IndexColumn struct {
	Name string
	Expr string
	// ExprNode retains the parsed key part for validation without reparsing.
	ExprNode   nodes.ExprNode
	ExprType   *ResolvedType
	Length     int
	Descending bool
}
