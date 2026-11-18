package ast

// This file holds AST node types for StarRocks's data cache management
// statements (StarRocks.g4:1845-1869).

// DataCacheSelectStmt represents dataCacheSelectStatement (StarRocks.g4:1867):
//
//	CACHE SELECT selectItem (',' selectItem)* FROM qualifiedName
//	    (WHERE expression)? properties?
type DataCacheSelectStmt struct {
	Items      []*SelectItem
	From       *ObjectName
	Where      Node
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *DataCacheSelectStmt) Tag() NodeTag { return T_DataCacheSelectStmt }

var _ Node = (*DataCacheSelectStmt)(nil)
