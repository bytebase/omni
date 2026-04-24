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
}

type IndexColumn struct {
	Name       string
	Expr       string
	ExprNode   nodes.ExprNode
	Length     int
	Descending bool
}
