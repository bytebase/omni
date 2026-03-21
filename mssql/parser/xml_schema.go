// Package parser - xml_schema.go implements T-SQL XML SCHEMA COLLECTION parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateXmlSchemaCollectionStmt parses a CREATE XML SCHEMA COLLECTION statement.
//
// BNF: mssql/parser/bnf/create-xml-schema-collection-transact-sql.bnf
//
//	CREATE XML SCHEMA COLLECTION [ <relational_schema>. ] sql_identifier AS Expression
func (p *Parser) parseCreateXmlSchemaCollectionStmt() (*nodes.CreateXmlSchemaCollectionStmt, error) {
	loc := p.pos()
	// COLLECTION keyword already consumed by caller

	stmt := &nodes.CreateXmlSchemaCollectionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// relational_schema.sql_identifier
	stmt.Name , _ = p.parseTableRef()

	// AS Expression
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
		p.advance()
	}
	stmt.XmlSchemaNamespaces, _ = p.parseExpr()

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterXmlSchemaCollectionStmt parses an ALTER XML SCHEMA COLLECTION statement.
//
// BNF: mssql/parser/bnf/alter-xml-schema-collection-transact-sql.bnf
//
//	ALTER XML SCHEMA COLLECTION [ relational_schema. ] sql_identifier ADD 'Schema Component'
func (p *Parser) parseAlterXmlSchemaCollectionStmt() (*nodes.AlterXmlSchemaCollectionStmt, error) {
	loc := p.pos()
	// COLLECTION keyword already consumed by caller

	stmt := &nodes.AlterXmlSchemaCollectionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// relational_schema.sql_identifier
	stmt.Name , _ = p.parseTableRef()

	// ADD Expression
	if _, ok := p.match(kwADD); ok {
		stmt.XmlSchemaNamespaces, _ = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
