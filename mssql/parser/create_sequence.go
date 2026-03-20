// Package parser - create_sequence.go implements CREATE/ALTER SEQUENCE parsing.
package parser

import (
	"strings"

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
func (p *Parser) parseCreateSequenceStmt() *nodes.CreateSequenceStmt {
	stmt := &nodes.CreateSequenceStmt{}
	stmt.Name , _ = p.parseTableRef()
	if p.cur.Type == kwAS {
		p.advance()
		stmt.DataType , _ = p.parseDataType()
	}
	p.parseSequenceOptions(stmt, false)
	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseAlterSequenceStmt() *nodes.AlterSequenceStmt {
	stmt := &nodes.AlterSequenceStmt{}
	stmt.Name , _ = p.parseTableRef()
	tmp := &nodes.CreateSequenceStmt{}
	p.parseSequenceOptions(tmp, true)
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
	stmt.Loc.End = p.pos()
	return stmt
}

// parseSequenceOptions parses the common options for CREATE/ALTER SEQUENCE.
// Options can appear in any order per the BNF grammar.
func (p *Parser) parseSequenceOptions(stmt *nodes.CreateSequenceStmt, isAlter bool) {
	for {
		if p.isIdentLike() && strings.EqualFold(p.cur.Str, "NO") {
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

		if !isAlter && p.matchIdentCI("START") {
			if p.cur.Type == kwWITH {
				p.advance()
			}
			stmt.Start, _ = p.parseExpr()
		} else if isAlter && p.matchIdentCI("RESTART") {
			stmt.Restart = true
			if p.cur.Type == kwWITH {
				p.advance()
				stmt.RestartWith, _ = p.parseExpr()
			}
		} else if p.matchIdentCI("INCREMENT") {
			p.matchIdentCI("BY")
			stmt.Increment, _ = p.parseExpr()
		} else if p.matchIdentCI("MINVALUE") {
			stmt.MinValue, _ = p.parseExpr()
		} else if p.matchIdentCI("MAXVALUE") {
			stmt.MaxValue, _ = p.parseExpr()
		} else if p.matchIdentCI("CYCLE") {
			b := true
			stmt.Cycle = &b
		} else if p.matchIdentCI("CACHE") {
			if p.cur.Type == tokICONST {
				stmt.Cache, _ = p.parseExpr()
			}
		} else {
			break
		}
	}
}

// isNextIdentCI checks if the next token (via peekNext) is an identifier-like
// token matching the given string case-insensitively, without consuming tokens.
func (p *Parser) isNextIdentCI(s string) bool {
	next := p.peekNext()
	return (next.Type == tokIDENT || (next.Type >= kwADD && next.Str != "")) &&
		strings.EqualFold(next.Str, s)
}
