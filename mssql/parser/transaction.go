// Package parser - transaction.go implements T-SQL transaction statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBeginTransStmt parses BEGIN TRAN/TRANSACTION.
//
// BNF: mssql/parser/bnf/begin-transaction-transact-sql.bnf
//
//	BEGIN { TRAN | TRANSACTION }
//	    [ { transaction_name | @tran_name_variable }
//	      [ WITH MARK [ 'description' ] ]
//	    ]
//
// beginTransactionStatement (TSql160.g:15961): the name may be left out
// before WITH MARK, and may be the old form [-]n : name.name
// (weirdTransactionName, TSql160.g:16022).
func (p *Parser) parseBeginTransStmt() (*nodes.BeginTransStmt, error) {
	loc := p.pos()
	p.advance() // consume BEGIN

	if _, ok := p.match(kwTRAN, kwTRANSACTION); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	// Completion: after BEGIN TRAN → identifier
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.BeginTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	switch {
	case p.cur.Type == tokVARIABLE:
		stmt.Name = p.cur.Str
		p.advance()
	case p.cur.Type == '-' || procIsInteger(p.cur):
		name, err := p.parseTransWeirdName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	case p.isIdentLike() && p.cur.Type != kwWITH:
		stmt.Name = p.cur.Str
		p.advance()
	}

	if p.cur.Type == kwWITH {
		p.advance()
		if err := p.procExpectWord("MARK"); err != nil {
			return nil, err
		}
		stmt.WithMark = true
		if procIsString(p.cur) || p.cur.Type == tokVARIABLE {
			stmt.MarkDescription = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseBeginDistributedTransStmt parses BEGIN DISTRIBUTED TRAN/TRANSACTION.
//
// BNF: mssql/parser/bnf/begin-transaction-transact-sql.bnf
//
//	BEGIN DISTRIBUTED { TRAN | TRANSACTION }
//	    [ transaction_name | @tran_name_variable ]
//
// TRAN or TRANSACTION is required (beginTransactionStatement,
// TSql160.g:15961).
func (p *Parser) parseBeginDistributedTransStmt() (*nodes.BeginDistributedTransStmt, error) {
	loc := p.pos()
	p.advance() // consume BEGIN
	p.advance() // consume DISTRIBUTED

	if _, ok := p.match(kwTRAN, kwTRANSACTION); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	stmt := &nodes.BeginDistributedTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if p.isIdentLike() || p.cur.Type == tokVARIABLE {
		stmt.Name = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCommitStmt parses a COMMIT statement.
//
// BNF: mssql/parser/bnf/commit-transaction-transact-sql.bnf
//
//	COMMIT [ { TRAN | TRANSACTION } [ transaction_name | @tran_name_variable ] ]
//	    [ WITH ( DELAYED_DURABILITY = { OFF | ON } ) ]
//	COMMIT [ WORK ]
func (p *Parser) parseCommitStmt() (*nodes.CommitTransStmt, error) {
	loc := p.pos()
	p.advance() // consume COMMIT

	// Completion: after COMMIT → TRANSACTION
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		return nil, errCollecting
	}

	stmt := &nodes.CommitTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}
	stmt.Name = p.parseTransEnd()

	// WITH ( DELAYED_DURABILITY = { OFF | ON } )
	if _, ok := p.match(kwWITH); ok {
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if err := p.procExpectWord("DELAYED_DURABILITY"); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type != kwON && p.cur.Type != kwOFF {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.DelayedDurability = strings.ToUpper(p.cur.Str)
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseRollbackStmt parses a ROLLBACK statement.
//
// BNF: mssql/parser/bnf/rollback-transaction-transact-sql.bnf
//
//	ROLLBACK { TRAN | TRANSACTION }
//	    [ transaction_name | @tran_name_variable
//	    | savepoint_name | @savepoint_variable ]
//	ROLLBACK [ WORK ]
func (p *Parser) parseRollbackStmt() (*nodes.RollbackTransStmt, error) {
	loc := p.pos()
	p.advance() // consume ROLLBACK

	// Completion: after ROLLBACK → TRANSACTION
	if p.collectMode() {
		p.addTokenCandidate(kwTRANSACTION)
		return nil, errCollecting
	}

	stmt := &nodes.RollbackTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}
	stmt.Name = p.parseTransEnd()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSaveTransStmt parses a SAVE TRANSACTION statement.
//
// BNF: mssql/parser/bnf/save-transaction-transact-sql.bnf
//
//	SAVE { TRAN | TRANSACTION } { savepoint_name | @savepoint_variable }
//
// saveTransactionStatement (TSql160.g:15873): both parts are required.
func (p *Parser) parseSaveTransStmt() (*nodes.SaveTransStmt, error) {
	loc := p.pos()
	p.advance() // consume SAVE

	if _, ok := p.match(kwTRAN, kwTRANSACTION); !ok {
		return nil, p.syntaxErrorAtCur()
	}

	// Completion: after SAVE TRAN → identifier
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.SaveTransStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	if !p.isIdentLike() && p.cur.Type != tokVARIABLE {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = p.cur.Str
	p.advance()

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// transSyntaxErrorAtPrev returns a syntax error at the token just consumed.
func (p *Parser) transSyntaxErrorAtPrev() *ParseError {
	return p.newParseError(p.prev.Loc, "syntax error at or near \""+p.tokenText(p.prev)+"\"")
}

// parseTransWeirdName parses weirdTransactionName (TSql160.g:16022):
// [ - ] integer : name . name.
func (p *Parser) parseTransWeirdName() (string, error) {
	start := p.pos()
	p.match('-')
	if !procIsInteger(p.cur) {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(':'); err != nil {
		return "", err
	}
	if _, ok := p.parseIdentifier(); !ok {
		return "", p.syntaxErrorAtCur()
	}
	if _, err := p.expect('.'); err != nil {
		return "", err
	}
	if _, ok := p.parseIdentifier(); !ok {
		return "", p.syntaxErrorAtCur()
	}
	return p.source[start:p.prevEnd()], nil
}

// parseTransEnd parses what COMMIT and ROLLBACK share: WORK, or TRAN or
// TRANSACTION and an optional name (commitTransactionStatement,
// TSql160.g:15917; rollbackTransactionStatement, TSql160.g:15883). A name
// needs TRAN or TRANSACTION before it: SQL Server refuses COMMIT a and
// COMMIT WORK a (Msg 102).
func (p *Parser) parseTransEnd() string {
	if p.procIsWord("WORK") {
		p.advance()
		return ""
	}
	if _, ok := p.match(kwTRAN, kwTRANSACTION); !ok {
		return ""
	}
	if p.cur.Type == tokVARIABLE || (p.isIdentLike() && p.cur.Type != kwWITH) {
		name := p.cur.Str
		p.advance()
		return name
	}
	return ""
}
