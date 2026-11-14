package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseInsert parses an INSERT statement (Parser.g:485):
//
//	INSERT INTO [keyspace.]table insertValuesSpec [IF NOT EXISTS] [USING TTL n [AND TIMESTAMP m]]
//	insertValuesSpec: '(' columnList ')' VALUES '(' term (',' term)* ')'
//	                | JSON (string | bind marker) [DEFAULT (NULL | UNSET)]
func (p *Parser) parseInsert() (*ast.InsertStmt, error) {
	start := p.curLoc()
	if err := p.expectKeyword(tokINSERT); err != nil {
		return nil, err
	}
	if err := p.expectKeyword(tokINTO); err != nil {
		return nil, err
	}

	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}

	stmt := &ast.InsertStmt{Table: table}

	// The column list, '(' cident (',' cident)* ')', comes with VALUES and
	// only with VALUES (Parser.g:498).
	if p.cur.Type != tokJSON {
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		for {
			col, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, col)
			if !p.match(tokCOMMA) {
				break
			}
		}
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	}

	// Parse insert values spec: VALUES (...) or JSON constant
	switch {
	case p.cur.Type == tokVALUES && stmt.Columns != nil:
		p.advance() // VALUES
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		values, err := p.parseTermList()
		if err != nil {
			return nil, err
		}
		stmt.Values = values
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	case p.cur.Type == tokJSON && stmt.Columns == nil:
		p.advance() // JSON
		stmt.IsJSON = true
		// jsonValue (Parser.g:523): a string or a bind marker.
		var jsonVal ast.ExprNode
		var err error
		if p.cur.Type == tokSTRING {
			jsonVal, err = p.parseConstant()
		} else if p.cur.Type == tokQMARK || p.cur.Type == tokCOLON {
			jsonVal, err = p.parseBindMarker()
		} else {
			err = p.errorf("expected string or bind marker after JSON, got %s", p.tokenDesc())
		}
		if err != nil {
			return nil, err
		}
		stmt.JSONValue = jsonVal
		// Optional DEFAULT UNSET | DEFAULT NULL
		if p.cur.Type == tokDEFAULT {
			p.advance() // DEFAULT
			if p.cur.Type == tokNULL {
				p.advance()
				stmt.DefaultNull = true
			} else if err := p.expectKeyword(tokUNSET); err != nil {
				return nil, err
			} else {
				stmt.DefaultUnset = true
			}
		}
	default:
		return nil, p.errorf("expected VALUES, got %s", p.tokenDesc())
	}

	// Optional IF NOT EXISTS
	ifne, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifne

	// Optional USING clause
	using, err := p.parseUsingClause()
	if err != nil {
		return nil, err
	}
	stmt.Using = using

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}
