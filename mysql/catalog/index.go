package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

type Index struct {
	Name         string
	Table        *Table
	Columns      []*IndexColumn
	Unique       bool
	Primary      bool
	Fulltext     bool
	Spatial      bool
	IndexType    string // BTREE, HASH, FULLTEXT, SPATIAL
	Comment      string
	Visible      bool
	KeyBlockSize int
	// Generated is true when the index was implicitly synthesized by InnoDB to
	// back a foreign key constraint (MySQL's HA_GENERATED_KEY). MySQL drops a
	// generated FK backing index when an ALTER TABLE or CREATE INDEX adds a new
	// index, and re-creates one only if no remaining index covers the FK.
	Generated bool
}

type IndexColumn struct {
	Name       string
	Expr       string
	ExprNode   nodes.ExprNode
	Length     int
	Descending bool
}
