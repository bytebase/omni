package parser

import (
	"math"
	"strconv"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseCreateSequenceStmt parses
//
//	CREATE SEQUENCE [IF NOT EXISTS] [schema.]name
//	    [sequence_option ...] [table_option [[,] table_option] ...]
//
// Sequence options all precede the table options: TiDB rejects
// `CREATE SEQUENCE s COMMENT 'x' START 1`. The grammar has no OR REPLACE
// and no TEMPORARY; parseCreateDispatch rejects both before calling here.
//
// Called from parseCreateDispatch after CREATE; p.cur is kwSEQUENCE.
//
// Ref: TiDB v8.5.5 parser.y:15462-15471 (CreateSequenceStmt).
func (p *Parser) parseCreateSequenceStmt(start int) (nodes.Node, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.CreateSequenceStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	// Completion: the user is naming a new sequence — no candidates.
	p.checkCursor()
	if p.collectMode() {
		return nil, &ParseError{Message: "collecting"}
	}

	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseSequenceOptions(false)
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	// CreateTableOptionListOpt (parser.y:12804-12823): options separated
	// by whitespace or a comma, and a comma must be followed by another
	// option.
	afterComma := false
	for {
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			if afterComma {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		stmt.TableOptions = append(stmt.TableOptions, opt)
		_, afterComma = p.match(',')
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSequenceStmt parses
//
//	ALTER SEQUENCE [IF EXISTS] [schema.]name alter_sequence_option ...
//
// The option list is required and takes no table options.
//
// Called from parseAlterDispatch after ALTER; p.cur is kwSEQUENCE.
//
// Ref: TiDB v8.5.5 parser.y:15598-15616 (AlterSequenceStmt).
func (p *Parser) parseAlterSequenceStmt(start int) (nodes.Node, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.AlterSequenceStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	name, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseSequenceOptions(true)
	if err != nil {
		return nil, err
	}
	if len(opts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropSequenceStmt parses DROP SEQUENCE [IF EXISTS] name [, name] ...
// The grammar has no TEMPORARY and no CASCADE/RESTRICT.
//
// Called from parseDropDispatch after DROP; p.cur is kwSEQUENCE.
//
// Ref: TiDB v8.5.5 parser.y:15575-15582 (DropSequenceStmt).
func (p *Parser) parseDropSequenceStmt(start int) (nodes.Node, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.DropSequenceStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}

	for {
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, &ParseError{Message: "collecting"}
		}

		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Sequences = append(stmt.Sequences, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// sequenceOptionTokens is the set of keywords that start a SequenceOption,
// offered as completion candidates. ALTER adds RESTART; CREATE adds COMMENT
// and ENGINE, the table options TiDB's DDL accepts on a sequence.
var sequenceOptionTokens = []int{
	kwINCREMENT, kwSTART, kwMINVALUE, kwNOMINVALUE, kwMAXVALUE, kwNOMAXVALUE,
	kwCACHE, kwNOCACHE, kwCYCLE, kwNOCYCLE, kwNO,
}

// parseSequenceOptions parses zero or more sequence options, in any order
// and with repeats allowed. RESTART is accepted only when isAlter is true.
//
//	INCREMENT [BY | =] n
//	START [WITH | =] n
//	MINVALUE [=] n | NOMINVALUE | NO MINVALUE
//	MAXVALUE [=] n | NOMAXVALUE | NO MAXVALUE
//	CACHE [=] n | NOCACHE | NO CACHE
//	CYCLE | NOCYCLE | NO CYCLE
//	RESTART [[WITH | =] n]                      -- ALTER only
//
// Ref: TiDB v8.5.5 parser.y:15489-15553 (SequenceOption) and
// parser.y:15618-15632 (AlterSequenceOption).
func (p *Parser) parseSequenceOptions(isAlter bool) ([]*nodes.SequenceOption, error) {
	var opts []*nodes.SequenceOption
	for {
		p.checkCursor()
		if p.collectMode() {
			for _, t := range sequenceOptionTokens {
				p.addTokenCandidate(t)
			}
			if isAlter {
				p.addTokenCandidate(kwRESTART)
			} else {
				p.addTokenCandidate(kwCOMMENT)
				p.addTokenCandidate(kwENGINE)
			}
			return nil, &ParseError{Message: "collecting"}
		}

		start := p.pos()
		opt := &nodes.SequenceOption{Loc: nodes.Loc{Start: start}}
		withValue := false

		switch p.cur.Type {
		case kwINCREMENT:
			p.advance()
			if _, ok := p.match(kwBY); !ok {
				p.match('=')
			}
			opt.Type = nodes.SequenceOptionIncrementBy
			withValue = true
		case kwSTART:
			p.advance()
			if _, ok := p.match(kwWITH); !ok {
				p.match('=')
			}
			opt.Type = nodes.SequenceOptionStartWith
			withValue = true
		case kwMINVALUE:
			p.advance()
			p.match('=')
			opt.Type = nodes.SequenceOptionMinValue
			withValue = true
		case kwMAXVALUE:
			p.advance()
			p.match('=')
			opt.Type = nodes.SequenceOptionMaxValue
			withValue = true
		case kwCACHE:
			p.advance()
			p.match('=')
			opt.Type = nodes.SequenceOptionCache
			withValue = true
		case kwNOMINVALUE:
			p.advance()
			opt.Type = nodes.SequenceOptionNoMinValue
		case kwNOMAXVALUE:
			p.advance()
			opt.Type = nodes.SequenceOptionNoMaxValue
		case kwNOCACHE:
			p.advance()
			opt.Type = nodes.SequenceOptionNoCache
		case kwCYCLE:
			p.advance()
			opt.Type = nodes.SequenceOptionCycle
		case kwNOCYCLE:
			p.advance()
			opt.Type = nodes.SequenceOptionNoCycle
		case kwNO:
			p.advance() // consume NO
			p.checkCursor()
			if p.collectMode() {
				for _, t := range []int{kwMINVALUE, kwMAXVALUE, kwCACHE, kwCYCLE} {
					p.addTokenCandidate(t)
				}
				return nil, &ParseError{Message: "collecting"}
			}
			switch p.cur.Type {
			case kwMINVALUE:
				opt.Type = nodes.SequenceOptionNoMinValue
			case kwMAXVALUE:
				opt.Type = nodes.SequenceOptionNoMaxValue
			case kwCACHE:
				opt.Type = nodes.SequenceOptionNoCache
			case kwCYCLE:
				opt.Type = nodes.SequenceOptionNoCycle
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		case kwRESTART:
			if !isAlter {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			// RESTART alone, or RESTART [WITH | =] n. Nothing else in the
			// statement can begin with WITH, '=', a sign, or a number, so
			// one token decides.
			switch p.cur.Type {
			case kwWITH, '=':
				p.advance()
				opt.Type = nodes.SequenceOptionRestartWith
				withValue = true
			case '+', '-', tokICONST:
				opt.Type = nodes.SequenceOptionRestartWith
				withValue = true
			default:
				opt.Type = nodes.SequenceOptionRestart
			}
		default:
			return opts, nil
		}

		if withValue {
			v, err := p.parseSequenceSignedNum()
			if err != nil {
				return nil, err
			}
			opt.Value = v
		}
		opt.Loc.End = p.prev.End
		opts = append(opts, opt)
	}
}

// parseSequenceSignedNum parses SignedNum: an integer literal with an
// optional sign, not an expression. Unsigned and '+' values must fit int64;
// a '-' value may reach -9223372036854775808. TiDB rejects an out-of-range
// value as a syntax error, as it does a decimal, hex, string, or
// parenthesized value.
//
// Ref: TiDB v8.5.5 parser.y:15555-15573 (SignedNum) and
// parser.y:6047-6056 (Int64Num).
func (p *Parser) parseSequenceSignedNum() (int64, error) {
	neg := false
	switch p.cur.Type {
	case '-':
		neg = true
		p.advance()
	case '+':
		p.advance()
	}
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	tok := p.cur
	u, err := strconv.ParseUint(tok.Str, 10, 64)
	limit := uint64(math.MaxInt64)
	if neg {
		limit++ // -2^63 is representable
	}
	if err != nil || u > limit {
		return 0, &ParseError{
			Message:  "sequence value out of range [-9223372036854775808, 9223372036854775807]",
			Position: tok.Loc,
		}
	}
	p.advance()
	if neg {
		// For u == 2^63, -int64(u) wraps to math.MinInt64, which is the
		// intended value.
		return -int64(u), nil
	}
	return int64(u), nil
}
