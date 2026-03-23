package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseExplainStmt parses an EXPLAIN statement.
// The EXPLAIN keyword has already been consumed.
func (p *Parser) parseExplainStmt() (nodes.Node, error) {
	loc := p.prev.Loc // EXPLAIN already consumed
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
		return &nodes.ExplainStmt{Query: query, Options: opts, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
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
				Loc: nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		query, err := p.parseExplainableStmt()
		if err != nil {
			return nil, err
		}
		return &nodes.ExplainStmt{
			Query:   query,
			Options: &nodes.List{Items: []nodes.Node{&nodes.DefElem{Defname: "analyze", Loc: nodes.NoLoc()}}},
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
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
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
	query, err := p.parseExplainableStmt()
	if err != nil {
		return nil, err
	}
	return &nodes.ExplainStmt{Query: query, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
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
	loc := p.prev.Loc // DO already consumed
	items := []nodes.Node{p.parseDostmtOptItem()}
	for p.cur.Type == SCONST || p.cur.Type == LANGUAGE {
		items = append(items, p.parseDostmtOptItem())
	}
	return &nodes.DoStmt{Args: &nodes.List{Items: items}, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
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
	return &nodes.CheckPointStmt{Loc: nodes.Loc{Start: p.prev.Loc, End: p.prev.End}}, nil
}

// parseDiscardStmt parses a DISCARD statement.
func (p *Parser) parseDiscardStmt() (nodes.Node, error) {
	loc := p.prev.Loc // DISCARD already consumed
	switch p.cur.Type {
	case ALL:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_ALL, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case TEMP:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case TEMPORARY:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case PLANS:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_PLANS, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	case SEQUENCES:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_SEQUENCES, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	default:
		return nil, nil
	}
}

// parseListenStmt parses a LISTEN statement.
func (p *Parser) parseListenStmt() (nodes.Node, error) {
	loc := p.prev.Loc // LISTEN already consumed
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.ListenStmt{Conditionname: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseUnlistenStmt parses an UNLISTEN statement.
func (p *Parser) parseUnlistenStmt() (nodes.Node, error) {
	loc := p.prev.Loc // UNLISTEN already consumed
	if p.cur.Type == '*' {
		p.advance()
		return &nodes.UnlistenStmt{Conditionname: "", Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
	}
	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	return &nodes.UnlistenStmt{Conditionname: name, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseNotifyStmt parses a NOTIFY statement.
func (p *Parser) parseNotifyStmt() (nodes.Node, error) {
	loc := p.prev.Loc // NOTIFY already consumed
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
	return &nodes.NotifyStmt{Conditionname: name, Payload: payload, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseLoadStmt parses a LOAD statement.
func (p *Parser) parseLoadStmt() (nodes.Node, error) {
	loc := p.prev.Loc // LOAD already consumed
	filename := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	return &nodes.LoadStmt{Filename: filename, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}

// parseCallStmt parses a CALL statement.
// The CALL keyword has already been consumed.
//
// Ref: https://www.postgresql.org/docs/17/sql-call.html
//
//	CALL name ( [ argument ] [, ...] )
func (p *Parser) parseCallStmt() (nodes.Node, error) {
	stmtLoc := p.prev.Loc // CALL already consumed
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
		Loc:      nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseReassignOwnedStmt parses a REASSIGN OWNED BY statement.
func (p *Parser) parseReassignOwnedStmt() (nodes.Node, error) {
	loc := p.prev.Loc // REASSIGN already consumed
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
	return &nodes.ReassignOwnedStmt{Roles: roles, Newrole: newrole, Loc: nodes.Loc{Start: loc, End: p.prev.End}}, nil
}
