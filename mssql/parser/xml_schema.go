// Package parser - xml_schema.go implements T-SQL XML SCHEMA COLLECTION parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateXmlSchemaCollectionStmt parses a CREATE XML SCHEMA COLLECTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-xml-schema-collection-transact-sql
//
//	CREATE XML SCHEMA COLLECTION [ <relational_schema>. ]sql_identifier AS Expression
func (p *Parser) parseCreateXmlSchemaCollectionStmt() *nodes.CreateXmlSchemaCollectionStmt {
	loc := p.pos()
	// COLLECTION keyword already consumed by caller

	stmt := &nodes.CreateXmlSchemaCollectionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// relational_schema.sql_identifier
	stmt.Name = p.parseTableRef()

	// AS Expression
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
		p.advance()
	}
	stmt.XmlSchemaNamespaces = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAlterXmlSchemaCollectionStmt parses an ALTER XML SCHEMA COLLECTION statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/alter-xml-schema-collection-transact-sql
//
//	ALTER XML SCHEMA COLLECTION [ relational_schema. ] sql_identifier ADD 'Schema Component'
func (p *Parser) parseAlterXmlSchemaCollectionStmt() *nodes.AlterXmlSchemaCollectionStmt {
	loc := p.pos()
	// COLLECTION keyword already consumed by caller

	stmt := &nodes.AlterXmlSchemaCollectionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// relational_schema.sql_identifier
	stmt.Name = p.parseTableRef()

	// ADD Expression
	if _, ok := p.match(kwADD); ok {
		stmt.XmlSchemaNamespaces = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
