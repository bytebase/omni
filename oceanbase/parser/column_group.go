package parser

// parseWithColumnGroup parses with_column_group where one stands, and returns
// its list as written:
//
//	WITH COLUMN GROUP ( element [, element] ... )
//
//	element: ALL COLUMNS | EACH COLUMN | name ( column [, column] ... )
//
// OceanBase's lexer reads WITH COLUMN GROUP as one token, WITH_COLUMN_GROUP, so
// the three words stand together. COLUMN is reserved, so WITH COLUMN opens no
// common table expression.
//
// Ref: with_column_group, column_group_list and column_group_element in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, line 8507).
func (p *Parser) parseWithColumnGroup() (string, bool, error) {
	if p.cur.Type != kwWITH || p.peekNext().Type != kwCOLUMN {
		return "", false, nil
	}
	p.advance() // consume WITH
	p.advance() // consume COLUMN
	if _, err := p.expect(kwGROUP); err != nil {
		return "", true, err
	}
	list, err := p.parseColumnGroupList()
	return list, true, err
}

// parseColumnGroupList parses '(' column_group_list ')' and returns the list
// inside the parentheses as written.
func (p *Parser) parseColumnGroupList() (string, error) {
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	start := p.pos()
	for {
		switch {
		case p.cur.Type == kwALL:
			p.advance()
			if !p.isBareWord("COLUMNS") {
				return "", p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isBareWord("EACH") && p.peekNext().Type == kwCOLUMN:
			p.advance()
			p.advance()
		default:
			if _, _, err := p.parseIdent(); err != nil {
				return "", err
			}
			if _, err := p.expect('('); err != nil {
				return "", err
			}
			for {
				if _, _, err := p.parseIdent(); err != nil {
					return "", err
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			if _, err := p.expect(')'); err != nil {
				return "", err
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	text := p.inputText(start, p.prev.End)
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return text, nil
}
