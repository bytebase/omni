package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// The bodies of STORAGE LIFECYCLE POLICY and DATA MOVEMENT RULE
// ---------------------------------------------------------------------------
//
// Both objects are registered in data_protection.go. A body is one clause:
// its signature in Params, its RETURNS type in Type, and the expression after
// -> in Body.

// parseCreateStorageLifecyclePolicy reads what follows the name of
//
//	CREATE [ OR REPLACE ] STORAGE LIFECYCLE POLICY [ IF NOT EXISTS ] <name>
//	  AS ( <arg_name> <arg_type> [ , ... ] )
//	  RETURNS BOOLEAN -> <body>
//	  [ ARCHIVE_TIER = { COOL | COLD } ]
//	  [ ARCHIVE_FOR_DAYS = <number_of_days> ]
//	  [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-storage-lifecycle-policy
func parseCreateStorageLifecyclePolicy(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	params, err := p.parsePolicyArgs()
	if err != nil {
		return err
	}
	ret, body, err := p.parseBodyReturns("BOOLEAN", false)
	if err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, "AS", start)
	c.Params, c.Type, c.Body = params, ret, body
	_, err = p.parseObjectProperties(s, "", false)
	return err
}

// parseCreateDataMovementRule reads what follows the name of
//
//	CREATE [ OR REPLACE ] DATA MOVEMENT RULE [ IF NOT EXISTS ] <name>
//	  TYPE = '<movement_type>'
//	  MAX_ROWS AS () RETURNS INTEGER
//	  -> ( <expression> )
//	  [ COMMENT = '<string_literal>' ]
//
// TYPE is a property clause of its own, MAX_ROWS AS the body's clause, and
// COMMENT a property clause after it.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-data-movement-rule
func parseCreateDataMovementRule(p *Parser, s *ast.ObjectStmt) error {
	if !p.atWords("TYPE") || !p.nextIs('=') {
		return p.syntaxErrorAtCur()
	}
	start := p.cur.Loc.Start
	typ, err := p.parseObjectProperty()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "", start).Options = []*ast.CopyOption{typ}
	if err := p.parseMaxRowsBody(s, "MAX_ROWS AS", p.cur.Loc.Start); err != nil {
		return err
	}
	if p.atWords("COMMENT") && p.nextIs('=') {
		start := p.cur.Loc.Start
		comment, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "", start).Options = []*ast.CopyOption{comment}
	}
	return nil
}

// parseMaxRowsBody reads MAX_ROWS AS ( ) RETURNS INTEGER -> ( <expression> )
// into a clause of keywords starting at start, whose Params is empty. cur is
// MAX_ROWS.
func (p *Parser) parseMaxRowsBody(s *ast.ObjectStmt, keywords string, start int) error {
	if err := p.expectWords("MAX_ROWS"); err != nil {
		return err
	}
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	ret, body, err := p.parseBodyReturns("INTEGER", true)
	if err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Params, c.Type, c.Body = []*ast.PolicyArg{}, ret, body
	return nil
}

// parseBodyReturns reads RETURNS <typ> -> <body>, typ being the one type the
// page writes. paren says the page writes the body in parentheses,
// -> ( <expression> ); it is then a ParenExpr.
func (p *Parser) parseBodyReturns(typ string, paren bool) (*ast.TypeName, ast.Node, error) {
	if _, err := p.expect(kwRETURNS); err != nil {
		return nil, nil, err
	}
	if !p.curWord(typ) {
		return nil, nil, p.syntaxErrorAtCur()
	}
	ret, err := p.parseDataType()
	if err != nil {
		return nil, nil, err
	}
	body, err := p.parseArrowBody(paren)
	if err != nil {
		return nil, nil, err
	}
	return ret, body, nil
}

// parseArrowBody reads -> <expression>, or -> ( <expression> ) into a
// ParenExpr when paren is set.
func (p *Parser) parseArrowBody(paren bool) (ast.Node, error) {
	if _, err := p.expect(tokArrow); err != nil {
		return nil, err
	}
	if !paren {
		return p.parseExpr()
	}
	open, err := p.expect('(')
	if err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &ast.ParenExpr{Expr: expr, Loc: ast.Loc{Start: open.Loc.Start, End: p.prev.Loc.End}}, nil
}
