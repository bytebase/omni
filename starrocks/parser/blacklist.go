package parser

import (
	"github.com/bytebase/omni/starrocks/ast"
)

// The SQL and backend blacklists (StarRocks.g4:1815-1842). Both open with ADD
// or DELETE, which are reserved, and the word after tells them apart from
// DELETE's own statement.

// atAddBlacklist reports whether an ADD at cur opens a blacklist statement.
func (p *Parser) atAddBlacklist() bool {
	if p.peekNext().Kind == kwSQLBLACKLIST {
		return true
	}
	if p.peekNext().Kind != kwBACKEND {
		return false
	}
	mark := p.save()
	p.advance() // consume ADD
	p.advance() // consume BACKEND
	ok := p.cur.Kind == kwBLACKLIST
	p.restore(mark)
	return ok
}

// atDeleteBlacklist reports whether a DELETE at cur opens a blacklist
// statement, rather than DELETE FROM.
func (p *Parser) atDeleteBlacklist() bool {
	if p.peekNext().Kind == kwSQLBLACKLIST {
		return true
	}
	if p.peekNext().Kind != kwBACKEND {
		return false
	}
	mark := p.save()
	p.advance() // consume DELETE
	p.advance() // consume BACKEND
	ok := p.cur.Kind == kwBLACKLIST
	p.restore(mark)
	return ok
}

// parseAddBlacklist parses addSqlBlackListStatement (StarRocks.g4:1815) and
// addBackendBlackListStatement (:1833):
//
//	ADD SQLBLACKLIST string
//	ADD BACKEND BLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
//
// On entry ADD has been consumed; cur is SQLBLACKLIST or BACKEND.
func (p *Parser) parseAddBlacklist(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.BlacklistStmt{Add: true}
	return p.parseBlacklistTail(startLoc, stmt)
}

// parseDeleteBlacklist parses delSqlBlackListStatement (StarRocks.g4:1819) and
// delBackendBlackListStatement (:1837):
//
//	DELETE SQLBLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
//	DELETE BACKEND BLACKLIST INTEGER_VALUE [',' INTEGER_VALUE]*
//
// On entry DELETE has been consumed; cur is SQLBLACKLIST or BACKEND.
func (p *Parser) parseDeleteBlacklist(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.BlacklistStmt{}
	return p.parseBlacklistTail(startLoc, stmt)
}

func (p *Parser) parseBlacklistTail(startLoc ast.Loc, stmt *ast.BlacklistStmt) (ast.Node, error) {
	if p.cur.Kind == kwSQLBLACKLIST {
		p.advance() // consume SQLBLACKLIST
		if stmt.Add {
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.SQL = p.cur.Str
			endLoc := p.cur.Loc
			p.advance()
			stmt.Loc = startLoc.Merge(endLoc)
			return stmt, nil
		}
	} else if p.cur.Kind == kwBACKEND {
		p.advance() // consume BACKEND
		if _, err := p.expect(kwBLACKLIST); err != nil {
			return nil, err
		}
		stmt.Backend = true
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	var endLoc ast.Loc
	for {
		if p.cur.Kind != tokInt {
			return nil, p.syntaxErrorAtCur()
		}
		endLoc = p.cur.Loc
		stmt.IDs = append(stmt.IDs, p.advance().Str)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// atPlanAdvisor reports whether PLAN ADVISOR follows the word at cur.
func (p *Parser) atPlanAdvisor() bool {
	if p.peekNext().Kind != kwPLAN {
		return false
	}
	mark := p.save()
	p.advance() // consume the leading word
	p.advance() // consume PLAN
	ok := p.cur.Kind == kwADVISOR
	p.restore(mark)
	return ok
}

// parseAlterPlanAdvisor parses alterPlanAdvisorAddStatement (StarRocks.g4:2024)
// and alterPlanAdvisorDropStatement (:2030):
//
//	ALTER PLAN ADVISOR ADD queryStatement
//	ALTER PLAN ADVISOR DROP string
//
// On entry ALTER has been consumed; cur is PLAN.
func (p *Parser) parseAlterPlanAdvisor(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PLAN
	if _, err := p.expect(kwADVISOR); err != nil {
		return nil, err
	}
	stmt := &ast.AlterPlanAdvisorStmt{}
	switch p.cur.Kind {
	case kwADD:
		p.advance() // consume ADD
		stmt.Action = "ADD"
		query, err := p.parseQueryRelation()
		if err != nil {
			return nil, err
		}
		stmt.Query = query
		stmt.Loc = startLoc.Merge(ast.NodeLoc(query))
	case kwDROP:
		p.advance() // consume DROP
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = "DROP"
		stmt.ID = p.cur.Str
		stmt.Loc = startLoc.Merge(p.cur.Loc)
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// parseTruncatePlanAdvisor parses truncatePlanAdvisorStatement
// (StarRocks.g4:2027): TRUNCATE PLAN ADVISOR.
//
// On entry TRUNCATE has NOT been consumed; cur is TRUNCATE.
func (p *Parser) parseTruncatePlanAdvisor() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume TRUNCATE
	p.advance() // consume PLAN
	advisorTok, err := p.expect(kwADVISOR)
	if err != nil {
		return nil, err
	}
	return &ast.TruncatePlanAdvisorStmt{Loc: startLoc.Merge(advisorTok.Loc)}, nil
}
