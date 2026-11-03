package catalog

import nodes "github.com/bytebase/omni/oceanbase/ast"

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
	// Scope is OceanBase's LOCAL or GLOBAL on an index of a partitioned table, "" when
	// the DDL states neither. A global index is a different physical object from a local
	// one, so it belongs in the diff rather than being dropped after parsing.
	Scope string
	// Parser is a fulltext index's WITH PARSER name, lower-cased, "" when the DDL names
	// none and OceanBase uses its default, space.
	Parser string
	// ParserProperties is a fulltext index's PARSER_PROPERTIES, key to value, nil when
	// the DDL gives none. A later PARSER_PROPERTIES replaces an earlier one, as on
	// OceanBase (fulltext_parser.go).
	ParserProperties map[string]string
}

type IndexColumn struct {
	Name       string
	Expr       string
	ExprNode   nodes.ExprNode
	Length     int
	Descending bool
}
