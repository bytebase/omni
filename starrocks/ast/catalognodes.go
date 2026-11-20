package ast

// This file holds DDL AST node types for CATALOG statements (T5.2).
//
// Supported forms:
//   CREATE EXTERNAL CATALOG [IF NOT EXISTS] name [COMMENT '...'] PROPERTIES(...)
//   ALTER CATALOG name SET (property, ...)
//   DROP CATALOG [IF EXISTS] name
//   REFRESH CATALOG name [PROPERTIES(...)]

// CreateCatalogStmt represents createExternalCatalogStatement
// (StarRocks.g4:760):
//
//	CREATE EXTERNAL CATALOG [IF NOT EXISTS] catalog_name
//	    [COMMENT 'comment']
//	    PROPERTIES ("key"="value", ...)
type CreateCatalogStmt struct {
	Name        string
	External    bool
	IfNotExists bool
	Comment     string
	Properties  []*Property
	Loc         Loc
}

// Tag implements Node.
func (n *CreateCatalogStmt) Tag() NodeTag { return T_CreateCatalogStmt }

var _ Node = (*CreateCatalogStmt)(nil)

// AlterCatalogStmt represents alterCatalogStatement (StarRocks.g4:2152):
//
//	ALTER CATALOG catalogName=identifierOrString SET '(' property [',' property]* ')'
//
// The grammar's modifyPropertiesClause is the statement's only clause: the
// pinned server refuses RENAME, SET PROPERTIES, MODIFY COMMENT and PROPERTY,
// which the fork inherited from Doris.
type AlterCatalogStmt struct {
	Name       string
	Properties []*Property
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
