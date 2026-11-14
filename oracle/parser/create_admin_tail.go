package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateAdminTail dispatches CREATE CONTEXT, CREATE [ CLEAN ] RESTORE
// POINT, CREATE [ PUBLIC ] ROLLBACK SEGMENT, and CREATE INMEMORY JOIN GROUP,
// from the word after the CREATE dispatcher's modifiers. ok is false when the
// current token starts none of them. ifTok is the IF of an IF NOT EXISTS the
// dispatcher consumed, and public its PUBLIC, or nil.
//
// Oracle 23ai takes OR REPLACE only before CONTEXT (ORA-00922 at the object
// word otherwise), PUBLIC only before ROLLBACK (ORA-00901 at the object word),
// and IF NOT EXISTS before none of them (ORA-00901 at the IF).
func (p *Parser) parseCreateAdminTail(start int, orReplace, ifNotExists bool, ifTok Token, public *Token) (stmt nodes.StmtNode, ok bool, err error) {
	var orReplaceOK, publicOK bool
	switch {
	case p.cur.Type == kwCONTEXT:
		orReplaceOK = true
	case p.isKeywordStr("RESTORE"), p.isKeywordStr("CLEAN"), p.isKeywordStr("INMEMORY"):
	case p.isKeywordStr("ROLLBACK"):
		publicOK = true
	default:
		return nil, false, nil
	}
	switch {
	case ifNotExists:
		return nil, true, p.syntaxErrorAtTok(ifTok)
	case orReplace && !orReplaceOK, public != nil && !publicOK:
		return nil, true, p.syntaxErrorAtCur()
	}

	switch {
	case p.cur.Type == kwCONTEXT:
		p.advance() // consume CONTEXT
		stmt, err = p.parseCreateContextStmt(start, orReplace)
	case p.isKeywordStr("CLEAN"), p.isKeywordStr("RESTORE"):
		var clean *Token
		if p.isKeywordStr("CLEAN") {
			cleanTok := p.cur
			clean = &cleanTok
			p.advance() // consume CLEAN
			if !p.isKeywordStr("RESTORE") {
				return nil, true, p.syntaxErrorAtCur()
			}
		}
		p.advance() // consume RESTORE
		if err := p.expectCreateTableWord("POINT"); err != nil {
			return nil, true, err
		}
		stmt, err = p.parseCreateRestorePointStmt(start, clean)
	case p.isKeywordStr("ROLLBACK"):
		p.advance() // consume ROLLBACK
		if err := p.expectCreateTableWord("SEGMENT"); err != nil {
			return nil, true, err
		}
		stmt, err = p.parseCreateRollbackSegmentStmt(start, public)
	default: // INMEMORY
		p.advance() // consume INMEMORY
		if err := p.expectCreateTableWord("JOIN"); err != nil {
			return nil, true, err
		}
		if err := p.expectCreateTableWord("GROUP"); err != nil {
			return nil, true, err
		}
		stmt, err = p.parseCreateInmemoryJoinGroupStmt(start)
	}
	if err != nil {
		return nil, true, err
	}
	return stmt, true, nil
}

// parseUnqualifiedName parses a name that takes no schema: a quoted
// identifier or a word that is not reserved. CREATE CONTEXT (ORA-02012 at a
// '.'), CREATE RESTORE POINT (ORA-03048), and the rollback segment statements
// (ORA-02176, ORA-02244) read their names so.
func (p *Parser) parseUnqualifiedName() (*nodes.ObjectName, error) {
	if !p.isCreateTableNameStart() {
		return nil, p.syntaxErrorAtCur()
	}
	name := &nodes.ObjectName{Name: p.cur.Str, Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}}
	p.advance()
	return name, nil
}
