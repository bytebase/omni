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
//
// The grammar is createSequenceStatement (TSql160.g:22948): AS may come
// anywhere among the options, each option at most once (Msg 11712).
func (p *Parser) parseCreateSequenceStmt() (*nodes.CreateSequenceStmt, error) {
	stmt := &nodes.CreateSequenceStmt{}
	var err error
	stmt.Name, err = p.sequenceParseName(p.parseTableRef)
	if err != nil {
		return nil, err
	}
	if _, err := p.parseSequenceOptions(stmt, false); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// sequenceParseName parses the name of a sequence, whose parts are not
// empty (schemaObjectTwoPartName, TSql160.g:30390).
func (p *Parser) sequenceParseName(parse func() (*nodes.TableRef, error)) (*nodes.TableRef, error) {
	name, err := parse()
	if err != nil {
		return nil, err
	}
	if name == nil || name.Object == "" {
		return nil, p.syntaxErrorAtCur()
	}
	return name, nil
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
	stmt.Name, err = p.sequenceParseName(p.parseSequenceRef)
	if err != nil {
		return nil, err
	}
	// alterSequenceOptionList (TSql160.g:12705): SQL Server wants at least
	// one option.
	tmp := &nodes.CreateSequenceStmt{}
	n, err := p.parseSequenceOptions(tmp, true)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, p.syntaxErrorAtCur()
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

// parseSequenceOptions parses the options of CREATE or ALTER SEQUENCE, in
// any order, and returns how many it read. Each option comes at most once,
// counting MINVALUE with NO MINVALUE and so on (Msg 11712). START WITH, AS
// and RESTART belong to one statement each (Msg 11709, 11710, 11711).
// MINVALUE, MAXVALUE, START WITH, INCREMENT BY and RESTART WITH take a
// signed decimal number, CACHE an unsigned integer
// (commonSequenceOptionListElement, TSql160.g:22997).
func (p *Parser) parseSequenceOptions(stmt *nodes.CreateSequenceStmt, isAlter bool) (int, error) {
	seen := map[string]bool{}
	n := 0
	for {
		optLoc := p.pos()
		var opt string
		var err error
		switch {
		case p.cur.Type == kwNO:
			p.advance()
			switch {
			case p.ddlIsWord("MINVALUE"):
				opt = "MINVALUE"
				stmt.NoMinVal = true
			case p.ddlIsWord("MAXVALUE"):
				opt = "MAXVALUE"
				stmt.NoMaxVal = true
			case p.ddlIsWord("CYCLE"):
				opt = "CYCLE"
				b := false
				stmt.Cycle = &b
			case p.ddlIsWord("CACHE"):
				opt = "CACHE"
				stmt.NoCache = true
			default:
				return n, p.syntaxErrorAtCur()
			}
			p.advance()
		case !isAlter && p.cur.Type == kwAS:
			p.advance()
			opt = "AS"
			stmt.DataType, err = p.parseDataType()
			if err == nil && stmt.DataType == nil {
				err = p.syntaxErrorAtCur()
			}
		case !isAlter && p.ddlIsWord("START"):
			p.advance()
			opt = "START"
			if _, err := p.expect(kwWITH); err != nil {
				return n, err
			}
			stmt.Start, err = p.sequenceParseNumber()
		case isAlter && p.ddlIsWord("RESTART"):
			p.advance()
			opt = "RESTART"
			stmt.Restart = true
			if p.cur.Type == kwWITH {
				p.advance()
				stmt.RestartWith, err = p.sequenceParseNumber()
			}
		case p.ddlIsWord("INCREMENT"):
			p.advance()
			opt = "INCREMENT"
			if err := p.ddlExpectWord("BY"); err != nil {
				return n, err
			}
			stmt.Increment, err = p.sequenceParseNumber()
		case p.ddlIsWord("MINVALUE"):
			p.advance()
			opt = "MINVALUE"
			stmt.MinValue, err = p.sequenceParseNumber()
		case p.ddlIsWord("MAXVALUE"):
			p.advance()
			opt = "MAXVALUE"
			stmt.MaxValue, err = p.sequenceParseNumber()
		case p.ddlIsWord("CYCLE"):
			p.advance()
			opt = "CYCLE"
			b := true
			stmt.Cycle = &b
		case p.ddlIsWord("CACHE"):
			p.advance()
			opt = "CACHE"
			if p.cur.Type == tokICONST && sequenceIsDecimal(p.cur.Str, false) {
				stmt.Cache, err = p.parsePrimary()
			}
		default:
			return n, nil
		}
		if err != nil {
			return n, err
		}
		if seen[opt] {
			return n, p.newParseError(optLoc, "argument "+opt+" can not be specified more than once")
		}
		seen[opt] = true
		n++
	}
}

// sequenceParseNumber parses a sequence option value: [ + | - ] a decimal
// number, integer or not. SQL Server refuses a string, a real (1e1), money,
// a binary literal, a variable and an expression (Msg 102).
func (p *Parser) sequenceParseNumber() (nodes.ExprNode, error) {
	loc := p.pos()
	var op nodes.UnaryOp
	signed := false
	switch p.cur.Type {
	case '-':
		op, signed = nodes.UnaryMinus, true
	case '+':
		op, signed = nodes.UnaryPlus, true
	}
	if signed {
		p.advance()
	}
	if (p.cur.Type != tokICONST && p.cur.Type != tokFCONST) || !sequenceIsDecimal(p.cur.Str, true) {
		return nil, p.syntaxErrorAtCur()
	}
	v, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	if !signed {
		return v, nil
	}
	return &nodes.UnaryExpr{Op: op, Operand: v, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}, nil
}

// sequenceIsDecimal reports whether s is written as a decimal number: digits,
// with one point when point is true.
func sequenceIsDecimal(s string, point bool) bool {
	digits, points := 0, 0
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			digits++
		case c == '.' && point:
			points++
		default:
			return false
		}
	}
	return digits > 0 && points <= 1
}
