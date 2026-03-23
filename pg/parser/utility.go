package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseExplainStmt parses an EXPLAIN statement.
// The EXPLAIN keyword has already been consumed.
func (p *Parser) parseExplainStmt() (nodes.Node, error) {
	if p.cur.Type == '(' {
		p.advance()
		opts := p.parseUtilityOptionList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		query, err := p.parseExplainableStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.ExplainStmt{Query: query, Options: opts}, nil
	}
	if p.cur.Type == ANALYZE || p.cur.Type == ANALYSE {
		p.advance()
		if p.cur.Type == VERBOSE {
			p.advance()
			query, err := p.parseExplainableStmt()
			if err != nil {
				return nil, err
			}
			return &nodes.ExplainStmt{
				Query: query,
				Options: &nodes.List{Items: []nodes.Node{
					&nodes.DefElem{Defname: "analyze", Loc: nodes.NoLoc()},
					&nodes.DefElem{Defname: "verbose", Loc: nodes.NoLoc()},
				}},
			}, nil
		}
		query, err := p.parseExplainableStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.ExplainStmt{
			Query:   query,
			Options: &nodes.List{Items: []nodes.Node{&nodes.DefElem{Defname: "analyze", Loc: nodes.NoLoc()}}},
		}, nil
	}
	if p.cur.Type == VERBOSE {
		p.advance()
		query, err := p.parseExplainableStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.ExplainStmt{
			Query:   query,
			Options: &nodes.List{Items: []nodes.Node{&nodes.DefElem{Defname: "verbose", Loc: nodes.NoLoc()}}},
		}, nil
	}
	query, err := p.parseExplainableStmt()
	if err != nil {
		return nil, err
	}
	return &nodes.ExplainStmt{Query: query}, nil
}

// parseExplainableStmt parses the statement that can follow EXPLAIN.
func (p *Parser) parseExplainableStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case SELECT, VALUES, TABLE, '(':
		return p.parseSelectNoParens()
	case WITH:
		return p.parseWithStmt()
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	case MERGE:
		return p.parseMergeStmt(nil)
	case EXECUTE:
		return p.parseExecuteStmt()
	case CREATE:
		return p.parseCreateDispatch()
	case DECLARE:
		return p.parseDeclareCursorStmt()
	case REFRESH:
		loc := p.pos()
		p.advance()
		stmt, err := p.parseRefreshMatViewStmt()
		if stmt != nil { stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End} }
		return stmt, err
	default:
		return nil, nil
	}
}

// parseDoStmt parses a DO statement. The DO keyword has already been consumed.
func (p *Parser) parseDoStmt() (nodes.Node, error) {
	items := []nodes.Node{p.parseDostmtOptItem()}
	for p.cur.Type == SCONST || p.cur.Type == LANGUAGE {
		items = append(items, p.parseDostmtOptItem())
	}
	return &nodes.DoStmt{Args: &nodes.List{Items: items}}, nil
}

// parseDostmtOptItem parses dostmt_opt_item.
func (p *Parser) parseDostmtOptItem() *nodes.DefElem {
	loc := p.pos()
	if p.cur.Type == LANGUAGE {
		p.advance()
		lang := p.parseNonReservedWordOrSconst()
		return &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: lang}, Loc: nodes.Loc{Start: loc, End: p.pos()}}
	}
	s := p.cur.Str
	p.expect(SCONST)
	return &nodes.DefElem{Defname: "as", Arg: &nodes.String{Str: s}, Loc: nodes.Loc{Start: loc, End: p.pos()}}
}

// parseCheckPointStmt parses a CHECKPOINT statement.
func (p *Parser) parseCheckPointStmt() (nodes.Node, error) {
	return &nodes.CheckPointStmt{}, nil
}

// parseDiscardStmt parses a DISCARD statement.
func (p *Parser) parseDiscardStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case ALL:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_ALL}, nil
	case TEMP:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}, nil
	case TEMPORARY:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}, nil
	case PLANS:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_PLANS}, nil
	case SEQUENCES:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_SEQUENCES}, nil
	default:
		return nil, nil
	}
}

// parseListenStmt parses a LISTEN statement.
func (p *Parser) parseListenStmt() (nodes.Node, error) {
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.ListenStmt{Conditionname: name}, nil
}

// parseUnlistenStmt parses an UNLISTEN statement.
func (p *Parser) parseUnlistenStmt() (nodes.Node, error) {
	if p.cur.Type == '*' {
		p.advance()
		return &nodes.UnlistenStmt{Conditionname: ""}, nil
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.UnlistenStmt{Conditionname: name}, nil
}

// parseNotifyStmt parses a NOTIFY statement.
func (p *Parser) parseNotifyStmt() (nodes.Node, error) {
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	payload := ""
	if p.cur.Type == ',' {
		p.advance()
		payload = p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
	}
	return &nodes.NotifyStmt{Conditionname: name, Payload: payload}, nil
}

// parseLoadStmt parses a LOAD statement.
func (p *Parser) parseLoadStmt() (nodes.Node, error) {
	filename := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	return &nodes.LoadStmt{Filename: filename}, nil
}

// parseCallStmt parses a CALL statement.
// The CALL keyword has already been consumed.
//
// Ref: https://www.postgresql.org/docs/17/sql-call.html
//
//	CALL name ( [ argument ] [, ...] )
func (p *Parser) parseCallStmt() (nodes.Node, error) {
	funcName, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	loc := p.pos()
	fc, err := p.parseFuncApplication(funcName, loc)
	if err != nil {
		return nil, err
	}
	if fc == nil {
		return nil, nil
	}
	return &nodes.CallStmt{
		Funccall: fc.(*nodes.FuncCall),
	}, nil
}

// parseReassignOwnedStmt parses a REASSIGN OWNED BY statement.
func (p *Parser) parseReassignOwnedStmt() (nodes.Node, error) {
	if _, err := p.expect(OWNED); err != nil {
		return nil, err
	}
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}
	roles := p.parseRoleList()
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	newrole := p.parseRoleSpec()
	return &nodes.ReassignOwnedStmt{Roles: roles, Newrole: newrole}, nil
}
