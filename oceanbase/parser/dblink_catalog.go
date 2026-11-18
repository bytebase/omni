package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// This file reads the OceanBase statements that name a database link, a catalog
// or a restore point, none of which MySQL has.

// parseCreateRestorePointStmt parses `CREATE RESTORE POINT relation_name`.
// p.cur is RESTORE.
//
// Ref: create_restore_point_stmt in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, lines 4906-4912).
func (p *Parser) parseCreateRestorePointStmt(start int) (nodes.Node, error) {
	p.advance() // consume RESTORE
	if _, err := p.expect(kwPOINT); err != nil {
		return nil, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.RestorePointStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Name: name,
	}, nil
}

// parseDropRestorePointStmt parses `DROP RESTORE POINT relation_name`. DROP is
// consumed and p.cur is RESTORE.
//
// Ref: drop_restore_point_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6,
// lines 4913-4918).
func (p *Parser) parseDropRestorePointStmt(start int) (nodes.Node, error) {
	p.advance() // consume RESTORE
	if _, err := p.expect(kwPOINT); err != nil {
		return nil, err
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.RestorePointStmt{
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
		Drop: true,
		Name: name,
	}, nil
}

// parseCreateCatalogStmt parses
// `CREATE EXTERNAL CATALOG [IF NOT EXISTS] relation_name PROPERTIES [=] '(' external_properties_list ')'`.
// CREATE and EXTERNAL are consumed, and p.cur is CATALOG.
//
// Ref: create_catalog_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 4926-4936).
func (p *Parser) parseCreateCatalogStmt(start int) (nodes.Node, error) {
	p.advance() // consume the CATALOG word

	stmt := &nodes.CatalogStmt{Loc: nodes.Loc{Start: start}, Action: "CREATE"}
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	if !p.isBareWord("PROPERTIES") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	p.match('=')
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	props, err := p.parseExternalProperties()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	stmt.Properties = props
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropCatalogStmt parses `DROP CATALOG [IF EXISTS] relation_name`. DROP is
// consumed and p.cur is CATALOG.
//
// Ref: drop_catalog_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 4944-4948).
func (p *Parser) parseDropCatalogStmt(start int) (nodes.Node, error) {
	p.advance() // consume the CATALOG word

	stmt := &nodes.CatalogStmt{Loc: nodes.Loc{Start: start}, Action: "DROP"}
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetCatalogStmt parses `SET CATALOG relation_name`, which set_catalog_stmt
// is: a statement of its own and no assignment of a variable called CATALOG.
// SET is consumed and p.cur is CATALOG.
//
// Ref: set_catalog_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 4950-4956).
func (p *Parser) parseSetCatalogStmt(start int) (nodes.Node, error) {
	p.advance() // consume the CATALOG word
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.CatalogStmt{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Action: "SET",
		Name:   name,
	}, nil
}

// atSetCatalog reports whether `SET CATALOG name` is this statement rather than
// an assignment: CATALOG is an ordinary name too, so `SET CATALOG = 1` assigns
// to a variable called CATALOG and `SET CATALOG a` switches catalogs.
func (p *Parser) atSetCatalog() bool {
	if !p.isBareWord("CATALOG") {
		return false
	}
	next := p.peekNext()
	return next.Type != '=' && next.Type != tokAssign && next.Type != kwTO &&
		next.Type != tokEOF && next.Type != ';'
}

// parseCreateDblinkStmt parses
// `CREATE DATABASE LINK [IF NOT EXISTS] dblink CONNECT TO user @tenant DATABASE db IDENTIFIED BY 'password' {SERVER [=] 'host:port' | HOST 'host:port'} [CLUSTER 'name']`.
// CREATE and DATABASE are consumed, and p.cur is LINK.
//
// Ref: create_dblink_stmt, dblink, tenant, ip_port and opt_cluster in
// sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 17690-17747 and 22283-22294).
func (p *Parser) parseCreateDblinkStmt(start int) (nodes.Node, error) {
	p.advance() // consume DATABASE
	p.advance() // consume the LINK word

	stmt := &nodes.DatabaseLinkStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CONNECT is no keyword of this lexer's, so it is a bare word here.
	if !p.isBareWord("CONNECT") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	// user: a string, a name or an unreserved keyword.
	user, _, err := p.parseIdentOrText()
	if err != nil {
		return nil, err
	}
	stmt.User = user

	// tenant: USER_VARIABLE, its @ written as it is.
	if !p.isUserVariable() {
		return nil, p.syntaxErrorAtCur()
	}
	vref, err := p.parseVariableRef()
	if err != nil {
		return nil, err
	}
	stmt.Tenant = "@" + vref.Name

	if _, err := p.expect(kwDATABASE); err != nil {
		return nil, err
	}
	// database_factor: a database's name, or `catalog.database`. Neither a
	// three-part name nor a database link is one.
	db, err := p.parseRecyclebinName("DATABASE")
	if err != nil {
		return nil, err
	}
	stmt.Database = db.Name
	if db.Catalog != "" {
		stmt.Database = db.Catalog + "." + db.Name
	}

	if _, err := p.expect(kwIDENTIFIED); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Password = p.cur.Str
	p.advance()

	// ip_port: SERVER [=] 'host:port' | HOST 'host:port'.
	if p.cur.Type == kwSERVER {
		p.advance()
		p.match('=')
	} else if p.cur.Type == kwHOST {
		p.advance()
	} else {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Server = p.cur.Str
	p.advance()

	// opt_cluster: CLUSTER 'name'.
	if p.isBareWord("CLUSTER") {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Cluster = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropDblinkStmt parses `DROP DATABASE LINK [IF EXISTS] dblink`. DROP and
// DATABASE are consumed, and p.cur is LINK.
//
// Ref: drop_dblink_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 17725-17730).
func (p *Parser) parseDropDblinkStmt(start int) (nodes.Node, error) {
	p.advance() // consume DATABASE
	p.advance() // consume the LINK word

	stmt := &nodes.DatabaseLinkStmt{Loc: nodes.Loc{Start: start}, Drop: true}
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// atDatabaseLink reports whether the tokens are DATABASE LINK, which switches
// CREATE/DROP DATABASE from the database statement to the database link's.
func (p *Parser) atDatabaseLink() bool {
	if p.cur.Type != kwDATABASE {
		return false
	}
	next := p.peekNext()
	if !((next.Type == tokIDENT || next.Type >= 700) && strings.EqualFold(next.Str, "LINK")) {
		return false
	}
	// LINK is an ordinary name too, so `DROP DATABASE LINK` is DROP DATABASE of
	// the database called LINK, and only a word after LINK makes a link.
	look := *p.lexer
	tok := look.NextToken()
	return tok.Type != tokEOF && tok.Type != ';'
}

// atExternalCatalog reports whether the tokens are CATALOG after CREATE
// EXTERNAL, which is CREATE EXTERNAL CATALOG and no table.
func (p *Parser) atExternalCatalog() bool {
	return p.isBareWord("CATALOG")
}
