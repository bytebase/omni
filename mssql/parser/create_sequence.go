// Package parser - create_sequence.go implements CREATE/ALTER SEQUENCE parsing.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreateSequenceStmt parses a CREATE SEQUENCE statement.
//
// BNF: mssql/parser/bnf/create-sequence-transact-sql.bnf
//
//	CREATE SEQUENCE [ schema_name . ] sequence_name
//	    [ AS [ built_in_integer_type | user-defined_integer_type ] ]
//	    [ START WITH <constant> ]
//	    [ INCREMENT BY <constant> ]
//	    [ { MINVALUE [ <constant> ] } | { NO MINVALUE } ]
//	    [ { MAXVALUE [ <constant> ] } | { NO MAXVALUE } ]
//	    [ CYCLE | { NO CYCLE } ]
//	    [ { CACHE [ <constant> ] } | { NO CACHE } ]
func (p *Parser) parseCreateSequenceStmt() (*nodes.CreateSequenceStmt, error) {
	stmt := &nodes.CreateSequenceStmt{}
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == kwAS {
		p.advance()
		stmt.DataType, err = p.parseDataType()
		if err != nil {
			return nil, err
		}
	}
	if err := p.parseSequenceOptions(stmt, false); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterSequenceStmt parses an ALTER SEQUENCE statement.
//
// BNF: mssql/parser/bnf/alter-sequence-transact-sql.bnf
//
//	ALTER SEQUENCE [ schema_name. ] sequence_name
//	    [ RESTART [ WITH <constant> ] ]
//	    [ INCREMENT BY <constant> ]
//	    [ { MINVALUE <constant> } | { NO MINVALUE } ]
//	    [ { MAXVALUE <constant> } | { NO MAXVALUE } ]
//	    [ CYCLE | { NO CYCLE } ]
//	    [ { CACHE [ <constant> ] } | { NO CACHE } ]
func (p *Parser) parseAlterSequenceStmt() (*nodes.AlterSequenceStmt, error) {
	stmt := &nodes.AlterSequenceStmt{}
	if p.collectMode() {
		p.addRuleCandidate("sequence_ref")
		return nil, errCollecting
	}
	var err error
	stmt.Name, err = p.parseSequenceRef()
	if err != nil {
		return nil, err
	}
	tmp := &nodes.CreateSequenceStmt{}
	if err := p.parseSequenceOptions(tmp, true); err != nil {
		return nil, err
	}
	stmt.Restart = tmp.Restart
	stmt.RestartWith = tmp.RestartWith
	stmt.Increment = tmp.Increment
	stmt.MinValue = tmp.MinValue
	stmt.MaxValue = tmp.MaxValue
	stmt.NoMinVal = tmp.NoMinVal
	stmt.NoMaxVal = tmp.NoMaxVal
	stmt.Cycle = tmp.Cycle
	stmt.Cache = tmp.Cache
	stmt.NoCache = tmp.NoCache
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseSequenceOptions parses the common options for CREATE/ALTER SEQUENCE.
// Options can appear in any order per the BNF grammar.
func (p *Parser) parseSequenceOptions(stmt *nodes.CreateSequenceStmt, isAlter bool) error {
	for {
		if p.cur.Type == kwNO {
			// NO { MINVALUE | MAXVALUE | CYCLE | CACHE }
			if p.isNextIdentCI("MINVALUE") {
				p.advance()
				p.advance()
				stmt.NoMinVal = true
				continue
			} else if p.isNextIdentCI("MAXVALUE") {
				p.advance()
				p.advance()
				stmt.NoMaxVal = true
				continue
			} else if p.isNextIdentCI("CYCLE") {
				p.advance()
				p.advance()
				b := false
				stmt.Cycle = &b
				continue
			} else if p.isNextIdentCI("CACHE") {
				p.advance()
				p.advance()
				stmt.NoCache = true
				continue
			}
		}

		var err error
		if !isAlter && p.matchIdentCI("START") {
			if p.cur.Type == kwWITH {
				p.advance()
			}
			stmt.Start, err = p.parseExpr()
			if err != nil {
				return err
			}
		} else if isAlter && p.matchIdentCI("RESTART") {
			stmt.Restart = true
			if p.cur.Type == kwWITH {
				p.advance()
				stmt.RestartWith, err = p.parseExpr()
				if err != nil {
					return err
				}
			}
		} else if p.matchIdentCI("INCREMENT") {
			p.matchIdentCI("BY")
			stmt.Increment, err = p.parseExpr()
			if err != nil {
				return err
			}
		} else if p.matchIdentCI("MINVALUE") {
			stmt.MinValue, err = p.parseExpr()
			if err != nil {
				return err
			}
		} else if p.matchIdentCI("MAXVALUE") {
			stmt.MaxValue, err = p.parseExpr()
			if err != nil {
				return err
			}
		} else if p.matchIdentCI("CYCLE") {
			b := true
			stmt.Cycle = &b
		} else if p.matchIdentCI("CACHE") {
			if p.cur.Type == tokICONST {
				stmt.Cache, err = p.parseExpr()
				if err != nil {
					return err
				}
			}
		} else {
			break
		}
	}
	return nil
}

// isNextIdentCI checks if the next token (via peekNext) is an identifier-like
// token matching the given string case-insensitively, without consuming tokens.
func (p *Parser) isNextIdentCI(s string) bool {
	next := p.peekNext()
	return next.Type == lookupKeyword(s)
}
