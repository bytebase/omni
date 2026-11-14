package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateContextStmt parses a CREATE CONTEXT statement. It is called
// after CREATE [ OR REPLACE ] CONTEXT.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-CONTEXT.html
//
//	CREATE [ OR REPLACE ] CONTEXT namespace
//	    USING [ schema. ] package
//	    [ INITIALIZED { EXTERNALLY | GLOBALLY } | ACCESSED GLOBALLY ]
//
// The page lists SHARING and lets INITIALIZED and ACCESSED GLOBALLY stand
// together. Oracle 23ai rejects SHARING and takes one of the other two only,
// with nothing after it (ORA-03049). The namespace takes no schema
// (ORA-02012); the package takes no database link (ORA-03048).
func (p *Parser) parseCreateContextStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_CONTEXT,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	name, err := p.parseUnqualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// USING [ schema. ] package
	opts := &nodes.List{}
	optStart := p.pos()
	if p.cur.Type != kwUSING {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume USING
	pkg, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "USING",
		Value: pkg.Name,
		Items: &nodes.List{Items: []nodes.Node{pkg}},
		Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
	})

	// [ INITIALIZED { EXTERNALLY | GLOBALLY } | ACCESSED GLOBALLY ]
	optStart = p.pos()
	switch {
	case p.isKeywordStr("INITIALIZED"):
		p.advance() // consume INITIALIZED
		if !p.isKeywordStr("EXTERNALLY") && !p.isKeywordStr("GLOBALLY") {
			return nil, p.syntaxErrorAtCur()
		}
		value := p.cur.Str
		p.advance()
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "INITIALIZED", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	case p.isKeywordStr("ACCESSED"):
		p.advance() // consume ACCESSED
		if err := p.expectCreateTableWord("GLOBALLY"); err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "ACCESSED GLOBALLY", Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
