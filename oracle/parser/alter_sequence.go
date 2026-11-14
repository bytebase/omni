package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAlterSequenceStmt parses an ALTER SEQUENCE statement. It is called
// with SEQUENCE as the current token.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-SEQUENCE.html
//
//	ALTER SEQUENCE [ IF EXISTS ] [ schema. ] sequence
//	    { INCREMENT BY integer
//	    | { MAXVALUE integer | NOMAXVALUE }
//	    | { MINVALUE integer | NOMINVALUE }
//	    | { CYCLE | NOCYCLE }
//	    | { CACHE integer | NOCACHE }
//	    | { ORDER | NOORDER }
//	    | { KEEP | NOKEEP }
//	    | { SCALE [ EXTEND | NOEXTEND ] | NOSCALE }
//	    | { SHARD [ EXTEND | NOEXTEND ] | NOSHARD }
//	    | { GLOBAL | SESSION }
//	    | RESTART [ START WITH integer ]
//	    } ...
//
// The page writes RESTART [ WITH integer ]; Oracle 23ai refuses that WITH
// (ORA-03048) and takes START WITH after RESTART. It takes the options in
// any order, each group once (ORA-02280 to ORA-02284 at a repeated or
// conflicting one), refuses START WITH without RESTART (ORA-02283), and
// needs one option at least (ORA-02286). The sequence takes no database link
// (ORA-02286 at the '@'). SHARD reaches ORA-02511 on a database that is not
// sharded, past the parser.
func (p *Parser) parseAlterSequenceStmt(start int) (nodes.StmtNode, error) {
	p.advance() // consume SEQUENCE
	stmt := &nodes.AlterSequenceStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance() // consume IF
		p.advance() // consume EXISTS
	}
	name, err := p.parsePlainObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	seen := map[string]bool{}
	extendSeen := false
	for {
		group := p.sequenceOptionGroup()
		if p.isKeywordStr("RESTART") {
			group = "RESTART"
		}
		if group == "" || group == "START" {
			if len(seen) == 0 || group == "START" {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		if seen[group] {
			return nil, p.syntaxErrorAtCur()
		}
		seen[group] = true

		switch {
		case p.cur.Type == kwINCREMENT:
			p.advance() // consume INCREMENT
			if err := p.expectCreateTableToken(kwBY); err != nil {
				return nil, err
			}
			if stmt.IncrementBy, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwMAXVALUE:
			p.advance() // consume MAXVALUE
			if stmt.MaxValue, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwNOMAXVALUE:
			stmt.NoMaxValue = true
			p.advance()
		case p.cur.Type == kwMINVALUE:
			p.advance() // consume MINVALUE
			if stmt.MinValue, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwNOMINVALUE:
			stmt.NoMinValue = true
			p.advance()
		case p.cur.Type == kwCYCLE:
			stmt.Cycle = true
			p.advance()
		case p.cur.Type == kwNOCYCLE:
			stmt.NoCycle = true
			p.advance()
		case p.cur.Type == kwCACHE:
			p.advance() // consume CACHE
			if stmt.Cache, err = p.parseRequiredExpr(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwNOCACHE:
			stmt.NoCache = true
			p.advance()
		case p.cur.Type == kwORDER:
			stmt.Order = true
			p.advance()
		case p.cur.Type == kwNOORDER:
			stmt.NoOrder = true
			p.advance()
		case p.cur.Type == kwKEEP:
			stmt.Keep = true
			p.advance()
		case p.isKeywordStr("NOKEEP"):
			stmt.NoKeep = true
			p.advance()
		case p.isKeywordStr("RESTART"):
			stmt.Restart = true
			p.advance() // consume RESTART
			if p.cur.Type == kwSTART {
				p.advance() // consume START
				if err := p.expectCreateTableToken(kwWITH); err != nil {
					return nil, err
				}
				if stmt.RestartWith, err = p.parseRequiredExpr(); err != nil {
					return nil, err
				}
			}
		case p.isKeywordStr("SCALE"):
			stmt.Scale = true
			p.advance() // consume SCALE
			if stmt.ScaleExtend, stmt.ScaleNoExtend, err = p.parseSequenceExtendModifier(&extendSeen); err != nil {
				return nil, err
			}
		case p.isKeywordStr("NOSCALE"):
			stmt.NoScale = true
			p.advance()
		case p.isKeywordStr("SHARD"):
			stmt.Shard = true
			p.advance() // consume SHARD
			if stmt.ShardExtend, stmt.ShardNoExtend, err = p.parseSequenceExtendModifier(&extendSeen); err != nil {
				return nil, err
			}
		case p.isKeywordStr("NOSHARD"):
			stmt.NoShard = true
			p.advance()
		case p.cur.Type == kwGLOBAL:
			stmt.Global = true
			p.advance()
		case p.cur.Type == kwSESSION:
			stmt.Session = true
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
