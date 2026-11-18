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
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// relational_schema.sql_identifier
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// AS Expression, both required: CREATE XML SCHEMA COLLECTION AS 1 is
	// Msg 156 and ... a AS is Msg 102 (measured).
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	stmt.XmlSchemaNamespaces, err = p.parseExpr()
	if err != nil {
		return nil, err
	}
	if stmt.XmlSchemaNamespaces == nil {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prevEnd()
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
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// relational_schema.sql_identifier
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// ADD Expression, both required: ALTER XML SCHEMA COLLECTION ADD 1 is
	// Msg 156 and ... a ADD is Msg 102 (measured).
	if _, err := p.expect(kwADD); err != nil {
		return nil, err
	}
	stmt.XmlSchemaNamespaces, err = p.parseExpr()
	if err != nil {
		return nil, err
	}
	if stmt.XmlSchemaNamespaces == nil {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
