package ast

// This file holds DDL AST node types for CATALOG statements (T5.2).
//
// Supported forms:
//   CREATE [EXTERNAL] CATALOG [IF NOT EXISTS] name [COMMENT '...'] [PROPERTIES(...)] [WITH RESOURCE name]
//   ALTER CATALOG name { RENAME new_name | SET PROPERTIES (...) | MODIFY COMMENT 'text' | PROPERTY (...) }
//   DROP CATALOG [IF EXISTS] name
//   REFRESH CATALOG name [PROPERTIES(...)]

// CreateCatalogStmt represents:
//
//	CREATE [EXTERNAL] CATALOG [IF NOT EXISTS] catalog_name
//	    [COMMENT 'comment']
//	    [PROPERTIES ("key"="value", ...)]
//	    [WITH RESOURCE resource_name]
type CreateCatalogStmt struct {
	Name         string
	External     bool
	IfNotExists  bool
	Comment      string
	Properties   []*Property
	WithResource string // non-empty when WITH RESOURCE resource_name is present
	Loc          Loc
}

// Tag implements Node.
func (n *CreateCatalogStmt) Tag() NodeTag { return T_CreateCatalogStmt }

var _ Node = (*CreateCatalogStmt)(nil)

// AlterCatalogAction identifies the kind of action in an ALTER CATALOG statement.
type AlterCatalogAction int

const (
	AlterCatalogRename         AlterCatalogAction = iota // RENAME new_name
	AlterCatalogSetProperties                            // SET PROPERTIES (...)
	AlterCatalogModifyComment                            // MODIFY COMMENT 'text'
	AlterCatalogSetProperty                              // PROPERTY ("key"="value")
)

// AlterCatalogStmt represents:
//
//	ALTER CATALOG name
//	    { RENAME new_name
//	    | SET PROPERTIES ("key"="value", ...)
//	    | MODIFY COMMENT 'text'
//	    | PROPERTY ("key"="value") }
type AlterCatalogStmt struct {
	Name       string
	Action     AlterCatalogAction
	NewName    string      // for RENAME
	Properties []*Property // for SET PROPERTIES / PROPERTY
	Comment    string      // for MODIFY COMMENT
	Loc        Loc
}

// Tag implements Node.
func (n *AlterCatalogStmt) Tag() NodeTag { return T_AlterCatalogStmt }

var _ Node = (*AlterCatalogStmt)(nil)

// DropCatalogStmt represents:
//
//	DROP CATALOG [IF EXISTS] catalog_name
type DropCatalogStmt struct {
	Name     string
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropCatalogStmt) Tag() NodeTag { return T_DropCatalogStmt }

var _ Node = (*DropCatalogStmt)(nil)

// RefreshCatalogStmt represents:
//
//	REFRESH CATALOG catalog_name [PROPERTIES(...)]
type RefreshCatalogStmt struct {
	Name       string
	Properties []*Property
	Loc        Loc
}

// Tag implements Node.
func (n *RefreshCatalogStmt) Tag() NodeTag { return T_RefreshCatalogStmt }

var _ Node = (*RefreshCatalogStmt)(nil)
