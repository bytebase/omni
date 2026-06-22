package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseDelete parses a DELETE statement:
//
//	DELETE [deleteColumnList] FROM [keyspace.]table [USING TIMESTAMP n] WHERE relationElements [IF EXISTS | IF ifConditionList]
//	deleteColumnList: deleteColumnItem (',' deleteColumnItem)*
//	deleteColumnItem: IDENT | IDENT '[' (string|decimal) ']'
func (p *Parser) parseDelete() (*ast.DeleteStmt, error) {
	start := p.curLoc()
	if err := p.expectKeyword(tokDELETE); err != nil {
		return nil, err
	}

	stmt := &ast.DeleteStmt{}

	// Parse optional delete column list (appears before FROM).
	// If the next token is not FROM, we have a column list.
	if p.cur.Type != tokFROM {
		cols, err := p.parseDeleteColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// FROM table
	if err := p.expectKeyword(tokFROM); err != nil {
		return nil, err
	}

	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.From = table

	// Optional USING TIMESTAMP clause
	using, err := p.parseUsingClause()
	if err != nil {
		return nil, err
	}
	stmt.Using = using

	// WHERE clause
	where, err := p.parseWhereClause()
	if err != nil {
		return nil, err
	}
	stmt.Where = where

	// Optional IF EXISTS or IF conditions
	if p.cur.Type == tokIF {
		if p.peekNext().Type == tokEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		} else {
			conds, err := p.parseIfConditions()
			if err != nil {
				return nil, err
			}
			stmt.IfConditions = conds
		}
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseDeleteColumnList parses: deleteColumnItem (',' deleteColumnItem)*
func (p *Parser) parseDeleteColumnList() ([]ast.ExprNode, error) {
	var cols []ast.ExprNode
	first, err := p.parseDeleteColumnItem()
	if err != nil {
		return nil, err
	}
	cols = append(cols, first)
	for p.match(tokCOMMA) {
		item, err := p.parseDeleteColumnItem()
		if err != nil {
			return nil, err
		}
		cols = append(cols, item)
	}
	return cols, nil
}

// parseDeleteColumnItem parses: IDENT | IDENT '[' expression ']'
func (p *Parser) parseDeleteColumnItem() (ast.ExprNode, error) {
	start := p.curLoc()

	ident, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for index access: IDENT '[' expression ']'
	if p.cur.Type == tokLBRACK {
		p.advance() // [
		idx, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(tokRBRACK); err != nil {
			return nil, err
		}
		return &ast.IndexAccess{
			Collection: ident,
			Index:      idx,
			Loc:        p.makeLoc(start),
		}, nil
	}

	return ident, nil
}
