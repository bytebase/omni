package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseExplainStmt parses an EXPLAIN statement.
// The EXPLAIN keyword has already been consumed.
func (p *Parser) parseExplainStmt() nodes.Node {
	if p.cur.Type == '(' {
		p.advance()
		opts := p.parseUtilityOptionList()
		p.expect(')')
		query := p.parseExplainableStmt()
		return &nodes.ExplainStmt{Query: query, Options: opts}
	}
	if p.cur.Type == ANALYZE || p.cur.Type == ANALYSE {
		p.advance()
		if p.cur.Type == VERBOSE {
			p.advance()
			query := p.parseExplainableStmt()
			return &nodes.ExplainStmt{
				Query: query,
				Options: &nodes.List{Items: []nodes.Node{
					&nodes.DefElem{Defname: "analyze"},
					&nodes.DefElem{Defname: "verbose"},
				}},
			}
		}
		query := p.parseExplainableStmt()
		return &nodes.ExplainStmt{
			Query:   query,
			Options: &nodes.List{Items: []nodes.Node{&nodes.DefElem{Defname: "analyze"}}},
		}
	}
	if p.cur.Type == VERBOSE {
		p.advance()
		query := p.parseExplainableStmt()
		return &nodes.ExplainStmt{
			Query:   query,
			Options: &nodes.List{Items: []nodes.Node{&nodes.DefElem{Defname: "verbose"}}},
		}
	}
	query := p.parseExplainableStmt()
	return &nodes.ExplainStmt{Query: query}
}

// parseExplainableStmt parses the statement that can follow EXPLAIN.
func (p *Parser) parseExplainableStmt() nodes.Node {
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
		p.advance()
		return p.parseRefreshMatViewStmt()
	default:
		return nil
	}
}

// parseDoStmt parses a DO statement. The DO keyword has already been consumed.
func (p *Parser) parseDoStmt() nodes.Node {
	items := []nodes.Node{p.parseDostmtOptItem()}
	for p.cur.Type == SCONST || p.cur.Type == LANGUAGE {
		items = append(items, p.parseDostmtOptItem())
	}
	return &nodes.DoStmt{Args: &nodes.List{Items: items}}
}

// parseDostmtOptItem parses dostmt_opt_item.
func (p *Parser) parseDostmtOptItem() *nodes.DefElem {
	if p.cur.Type == LANGUAGE {
		p.advance()
		lang := p.parseNonReservedWordOrSconst()
		return &nodes.DefElem{Defname: "language", Arg: &nodes.String{Str: lang}}
	}
	s := p.cur.Str
	p.expect(SCONST)
	return &nodes.DefElem{Defname: "as", Arg: &nodes.String{Str: s}}
}

// parseCheckPointStmt parses a CHECKPOINT statement.
func (p *Parser) parseCheckPointStmt() nodes.Node {
	return &nodes.CheckPointStmt{}
}

// parseDiscardStmt parses a DISCARD statement.
func (p *Parser) parseDiscardStmt() nodes.Node {
	switch p.cur.Type {
	case ALL:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_ALL}
	case TEMP:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}
	case TEMPORARY:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_TEMP}
	case PLANS:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_PLANS}
	case SEQUENCES:
		p.advance()
		return &nodes.DiscardStmt{Target: nodes.DISCARD_SEQUENCES}
	default:
		return nil
	}
}

// parseListenStmt parses a LISTEN statement.
func (p *Parser) parseListenStmt() nodes.Node {
	name, _ := p.parseColId()
	return &nodes.ListenStmt{Conditionname: name}
}

// parseUnlistenStmt parses an UNLISTEN statement.
func (p *Parser) parseUnlistenStmt() nodes.Node {
	if p.cur.Type == '*' {
		p.advance()
		return &nodes.UnlistenStmt{Conditionname: ""}
	}
	name, _ := p.parseColId()
	return &nodes.UnlistenStmt{Conditionname: name}
}

// parseNotifyStmt parses a NOTIFY statement.
func (p *Parser) parseNotifyStmt() nodes.Node {
	name, _ := p.parseColId()
	payload := ""
	if p.cur.Type == ',' {
		p.advance()
		payload = p.cur.Str
		p.expect(SCONST)
	}
	return &nodes.NotifyStmt{Conditionname: name, Payload: payload}
}

// parseLoadStmt parses a LOAD statement.
func (p *Parser) parseLoadStmt() nodes.Node {
	filename := p.cur.Str
	p.expect(SCONST)
	return &nodes.LoadStmt{Filename: filename}
}

// parseReassignOwnedStmt parses a REASSIGN OWNED BY statement.
func (p *Parser) parseReassignOwnedStmt() nodes.Node {
	p.expect(OWNED)
	p.expect(BY)
	roles := p.parseRoleList()
	p.expect(TO)
	newrole := p.parseRoleSpec()
	return &nodes.ReassignOwnedStmt{Roles: roles, Newrole: newrole}
}
