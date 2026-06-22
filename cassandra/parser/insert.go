package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseInsert parses an INSERT statement:
//
//	INSERT INTO [keyspace.]table [(columnList)] insertValuesSpec [IF NOT EXISTS] [USING TTL n [AND TIMESTAMP m]]
//	insertValuesSpec: VALUES '(' expressionList ')' | JSON constant [DEFAULT UNSET]
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

	// Parse optional column list: '(' identifier (',' identifier)* ')'
	if p.cur.Type == tokLPAREN {
		p.advance() // (
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
	switch p.cur.Type {
	case tokVALUES:
		p.advance() // VALUES
		if _, err := p.expect(tokLPAREN); err != nil {
			return nil, err
		}
		values, err := p.parseExpressionList()
		if err != nil {
			return nil, err
		}
		stmt.Values = values
		if _, err := p.expect(tokRPAREN); err != nil {
			return nil, err
		}
	case tokJSON:
		p.advance() // JSON
		stmt.IsJSON = true
		jsonVal, err := p.parseConstant()
		if err != nil {
			return nil, err
		}
		stmt.JSONValue = jsonVal
		// Optional DEFAULT UNSET
		if p.cur.Type == tokDEFAULT {
			p.advance() // DEFAULT
			if err := p.expectKeyword(tokUNSET); err != nil {
				return nil, err
			}
			stmt.DefaultUnset = true
		}
	default:
		return nil, p.errorf("expected VALUES or JSON, got %s", p.tokenDesc())
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
