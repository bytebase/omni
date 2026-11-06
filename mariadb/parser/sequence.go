// Package parser - sequence.go implements CREATE/ALTER/DROP SEQUENCE parsing
// for MariaDB 11.8 (BYT-9135).
//
// The option grammar is shared between CREATE and ALTER (container-verified):
// both accept the same any-order options including AS int_type and START; the
// only difference is RESTART, which is ALTER-only (CREATE ... RESTART is 1064).
// Option values are integer literals (not expressions), and each value is
// mandatory — bare CACHE/INCREMENT/START/MINVALUE/MAXVALUE are 1064.
package parser

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// sequenceASIntTypes is the set of integer type names MariaDB accepts after the
// sequence AS clause. Non-integer types and any display width (INT(11)) are 1064.
var sequenceASIntTypes = map[string]bool{
	"TINYINT":   true,
	"SMALLINT":  true,
	"MEDIUMINT": true,
	"INT":       true,
	"INTEGER":   true,
	"BIGINT":    true,
}

// sequenceOptions accumulates the any-order CREATE/ALTER SEQUENCE options before
// they are copied onto the concrete statement node.
type sequenceOptions struct {
	dataType    *nodes.DataType
	start       nodes.ExprNode
	restart     bool
	restartWith nodes.ExprNode
	increment   nodes.ExprNode
	minValue    nodes.ExprNode
	maxValue    nodes.ExprNode
	noMinValue  bool
	noMaxValue  bool
	cache       nodes.ExprNode
	noCache     bool
	cycle       *bool
}

// isEmpty reports whether no option was parsed. ALTER SEQUENCE requires at least
// one option — MariaDB rejects a bare `ALTER SEQUENCE name` with 1064 — whereas
// CREATE SEQUENCE permits no options.
func (o *sequenceOptions) isEmpty() bool {
	return o.dataType == nil && o.start == nil && !o.restart && o.restartWith == nil &&
		o.increment == nil && o.minValue == nil && o.maxValue == nil &&
		!o.noMinValue && !o.noMaxValue && o.cache == nil && !o.noCache && o.cycle == nil
}

// parseCreateSequenceStmt parses CREATE [OR REPLACE] SEQUENCE [IF NOT EXISTS].
// The CREATE keyword (and any OR REPLACE / TEMPORARY) is already consumed by the
// dispatcher; orReplace is propagated in. p.cur is the SEQUENCE keyword.
func (p *Parser) parseCreateSequenceStmt(start int, orReplace, temporary bool) (*nodes.CreateSequenceStmt, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.CreateSequenceStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
		Temporary: temporary,
	}

	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// Completion: sequence name is a new identifier (no candidates).
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
	stmt.DataType = opts.dataType
	stmt.Start = opts.start
	stmt.Increment = opts.increment
	stmt.MinValue = opts.minValue
	stmt.MaxValue = opts.maxValue
	stmt.NoMinValue = opts.noMinValue
	stmt.NoMaxValue = opts.noMaxValue
	stmt.Cache = opts.cache
	stmt.NoCache = opts.noCache
	stmt.Cycle = opts.cycle
	if err := p.refuseSequenceRange(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterSequenceStmt parses ALTER SEQUENCE [IF EXISTS] name with any-order
// options. p.cur is the SEQUENCE keyword (ALTER already consumed).
func (p *Parser) parseAlterSequenceStmt(start int) (*nodes.AlterSequenceStmt, error) {
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
	if opts.isEmpty() {
		// MariaDB rejects a bare ALTER SEQUENCE with no options (1064).
		return nil, p.syntaxErrorAtCur()
	}
	stmt.DataType = opts.dataType
	stmt.Start = opts.start
	stmt.Restart = opts.restart
	stmt.RestartWith = opts.restartWith
	stmt.Increment = opts.increment
	stmt.MinValue = opts.minValue
	stmt.MaxValue = opts.maxValue
	stmt.NoMinValue = opts.noMinValue
	stmt.NoMaxValue = opts.noMaxValue
	stmt.Cache = opts.cache
	stmt.NoCache = opts.noCache
	stmt.Cycle = opts.cycle

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropSequenceStmt parses DROP SEQUENCE [IF EXISTS] name [, name ...].
// p.cur is the SEQUENCE keyword (DROP and any TEMPORARY already consumed).
func (p *Parser) parseDropSequenceStmt(start int, temporary bool) (*nodes.DropSequenceStmt, error) {
	p.advance() // consume SEQUENCE

	stmt := &nodes.DropSequenceStmt{Loc: nodes.Loc{Start: start}, Temporary: temporary}

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

	for {
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		stmt.Sequences = append(stmt.Sequences, ref)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSequenceOptions parses the any-order CREATE/ALTER SEQUENCE option list.
// RESTART is accepted only when isAlter is true (CREATE ... RESTART is 1064).
func (p *Parser) parseSequenceOptions(isAlter bool) (*sequenceOptions, error) {
	opts := &sequenceOptions{}
	// Each option is once, a word and its NO form being one option: the action answers 1225,
	// ER_DUP_ARGUMENT, for the second.
	seen := map[int]bool{}
	for {
		slot := p.cur.Type
		switch slot {
		case kwNOCACHE:
			slot = kwCACHE
		case kwNOCYCLE:
			slot = kwCYCLE
		case kwNOMINVALUE:
			slot = kwMINVALUE
		case kwNOMAXVALUE:
			slot = kwMAXVALUE
		case kwNO:
			slot = p.peekNext().Type
		}
		if seen[slot] {
			return nil, &ParseError{Message: "Option used twice in statement", Position: p.cur.Loc}
		}
		seen[slot] = true
		switch p.cur.Type {
		case kwAS:
			p.advance()
			dt, err := p.parseSequenceASType()
			if err != nil {
				return nil, err
			}
			opts.dataType = dt
		case kwSTART:
			p.advance()
			if p.cur.Type == kwWITH {
				p.advance()
			} else {
				p.matchEqual()
			}
			v, err := p.parseSequenceValue()
			if err != nil {
				return nil, err
			}
			opts.start = v
		case kwINCREMENT:
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			} else {
				p.matchEqual()
			}
			v, err := p.parseSequenceValue()
			if err != nil {
				return nil, err
			}
			if lit, ok := v.(*nodes.IntLit); ok && lit.Value == math.MaxInt64 && p.prev.Str != "9223372036854775807" {
				return nil, p.syntaxErrorAtTok(p.prev)
			}
			opts.increment = v
		case kwMINVALUE:
			p.advance()
			p.matchEqual()
			v, err := p.parseTruncatedSequenceValue()
			if err != nil {
				return nil, err
			}
			opts.minValue = v
		case kwMAXVALUE:
			p.advance()
			p.matchEqual()
			v, err := p.parseTruncatedSequenceValue()
			if err != nil {
				return nil, err
			}
			opts.maxValue = v
		case kwCACHE:
			p.advance()
			p.matchEqual()
			v, err := p.parseSequenceValue()
			if err != nil {
				return nil, err
			}
			opts.cache = v
		case kwNOCACHE:
			p.advance()
			opts.noCache = true
		case kwCYCLE:
			p.advance()
			cycle := true
			opts.cycle = &cycle
		case kwNOCYCLE:
			p.advance()
			cycle := false
			opts.cycle = &cycle
		case kwNOMINVALUE:
			p.advance()
			opts.noMinValue = true
		case kwNOMAXVALUE:
			p.advance()
			opts.noMaxValue = true
		case kwNO:
			// MariaDB accepts spaced NO MINVALUE / NO MAXVALUE, but NO CACHE /
			// NO CYCLE are 1064 (one-word NOCACHE / NOCYCLE only).
			p.advance()
			switch p.cur.Type {
			case kwMINVALUE:
				p.advance()
				opts.noMinValue = true
			case kwMAXVALUE:
				p.advance()
				opts.noMaxValue = true
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case kwRESTART:
			if !isAlter {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			opts.restart = true
			if p.cur.Type == kwWITH || p.cur.Type == '=' || p.cur.Type == tokAssign {
				p.advance()
				v, err := p.parseSequenceValue()
				if err != nil {
					return nil, err
				}
				opts.restartWith = v
			} else if p.isSequenceValueStart() {
				v, err := p.parseSequenceValue()
				if err != nil {
					return nil, err
				}
				opts.restartWith = v
			}
		default:
			return opts, nil
		}
	}
}

// parseSequenceValue parses a sequence option value: an optionally-negated
// integer literal. MariaDB rejects expressions/parentheses here (1064), so this
// deliberately does not call parseExpr (which would also swallow a following
// non-reserved option keyword such as CYCLE).
func (p *Parser) parseSequenceValue() (nodes.ExprNode, error) {
	return p.parseSequenceValueRange(true)
}

// parseTruncatedSequenceValue parses MINVALUE's and MAXVALUE's value, which the server
// truncates to the type where the others answer 1064 past 64 bits.
func (p *Parser) parseTruncatedSequenceValue() (nodes.ExprNode, error) {
	return p.parseSequenceValueRange(false)
}

func (p *Parser) parseSequenceValueRange(bounded bool) (nodes.ExprNode, error) {
	start := p.pos()
	neg := false
	switch p.cur.Type {
	case '-':
		neg = true
		p.advance()
	case '+':
		p.advance()
	}
	if p.cur.Type != tokICONST {
		return nil, p.syntaxErrorAtCur()
	}
	// longlong_num on the minus side and ulonglong_num on the plus: what is past those is
	// a 1064.
	if neg && bounded {
		if _, err := strconv.ParseInt("-"+p.cur.Str, 10, 64); err != nil {
			return nil, p.syntaxErrorAtCur()
		}
	} else if bounded {
		if _, err := strconv.ParseUint(p.cur.Str, 10, 64); err != nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	tok := p.advance()
	lit := &nodes.IntLit{Loc: nodes.Loc{Start: tok.Loc, End: p.prev.End}, Value: tok.Ival}
	if neg {
		return &nodes.UnaryExpr{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			Op:         nodes.UnaryMinus,
			Operand:    lit,
			OriginalOp: "-",
		}, nil
	}
	return lit, nil
}

// isSequenceValueStart reports whether the current token can begin a sequence
// option value, used to decide if RESTART has an (optional) value.
func (p *Parser) isSequenceValueStart() bool {
	return p.cur.Type == '-' || p.cur.Type == '+' || p.cur.Type == tokICONST
}

// parseSequenceASType parses and validates the AS clause type. MariaDB restricts
// it to integer types without display width; everything else is 1064.
func (p *Parser) parseSequenceASType() (*nodes.DataType, error) {
	loc := p.cur.Loc
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if !sequenceASIntTypes[strings.ToUpper(dt.Name)] || dt.Length != 0 || dt.Scale != 0 {
		return nil, &ParseError{
			Message:  "sequence AS type must be an integer type without display width",
			Position: loc,
		}
	}
	return dt, nil
}

// parseNextValueFor parses NEXT VALUE FOR sequence_name as a primary expression.
// MariaDB admits it in any expression position (the CHECK / generated-column /
// partition restrictions are semantic, not 1064), but rejects a trailing OVER
// clause (1064) — so, unlike the mssql donor, no OVER is parsed. The caller has
// already confirmed p.cur is NEXT and the next token is VALUE.
func (p *Parser) parseNextValueFor() (nodes.ExprNode, error) {
	if err := p.refuseStoredExprSubquery(p.pos(), true); err != nil {
		return nil, err
	}
	start := p.pos()
	p.advance() // consume NEXT
	if _, err := p.expect(kwVALUE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	seq, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	return &nodes.NextValueForExpr{
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
		Sequence: seq,
	}, nil
}

// parsePreviousValueFor parses PREVIOUS VALUE FOR sequence_name (the SQL-standard
// spelling alongside the generic LASTVAL() function). Mirrors parseNextValueFor.
func (p *Parser) parsePreviousValueFor() (nodes.ExprNode, error) {
	start := p.pos()
	p.advance() // consume PREVIOUS
	if _, err := p.expect(kwVALUE); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	seq, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	return &nodes.PreviousValueForExpr{
		Loc:      nodes.Loc{Start: start, End: p.prev.End},
		Sequence: seq,
	}, nil
}

// sequenceBound holds a type's range as check_and_adjust reads it: the type's own maximum less
// one, and its minimum plus one when it is signed.
func sequenceBound(dt *nodes.DataType) (lo, hi *big.Int) {
	bits, unsigned := uint(64), false
	if dt != nil {
		switch dt.Name {
		case "TINYINT":
			bits = 8
		case "SMALLINT":
			bits = 16
		case "MEDIUMINT":
			bits = 24
		case "INT":
			bits = 32
		}
		unsigned = dt.Unsigned
	}
	one := big.NewInt(1)
	if unsigned {
		top := new(big.Int).Lsh(one, bits)
		return big.NewInt(0), top.Sub(top, big.NewInt(2))
	}
	half := new(big.Int).Lsh(one, bits-1)
	return new(big.Int).Add(new(big.Int).Neg(half), one), new(big.Int).Sub(half, big.NewInt(2))
}

// sequenceValue reads what parseSequenceValue built, from the text: a value may not fit a
// signed 64 bits, which is what the literal holds.
func (p *Parser) sequenceValue(v nodes.ExprNode) *big.Int {
	neg := false
	if unary, ok := v.(*nodes.UnaryExpr); ok {
		neg = true
		v = unary.Operand
	}
	lit, ok := v.(*nodes.IntLit)
	if !ok {
		return new(big.Int)
	}
	n, _ := new(big.Int).SetString(p.inputText(lit.Loc.Start, lit.Loc.End), 10)
	if n == nil {
		return new(big.Int)
	}
	if neg {
		n.Neg(n)
	}
	return n
}

// refuseSequenceRange is sequence_definition::check_and_adjust, which CREATE SEQUENCE runs while
// it parses and answers 4085, ER_SEQUENCE_INVALID_DATA, for. MINVALUE and MAXVALUE are truncated
// to the type's range first, and START is not. ALTER SEQUENCE checks against what the sequence
// already holds, so it is not checked here.
//
// Ref: sequence_definition::check_and_adjust in MariaDB 11.8 sql/sql_sequence.cc.
func (p *Parser) refuseSequenceRange(stmt *nodes.CreateSequenceStmt) error {
	lo, hi := sequenceBound(stmt.DataType)
	clamp := func(v *big.Int) *big.Int {
		if v.Cmp(lo) < 0 {
			return lo
		}
		if v.Cmp(hi) > 0 {
			return hi
		}
		return v
	}
	increment := big.NewInt(1)
	if stmt.Increment != nil {
		increment = p.sequenceValue(stmt.Increment)
	}
	descending := increment.Sign() < 0
	minimum, maximum := big.NewInt(1), hi
	if descending {
		minimum, maximum = lo, big.NewInt(-1)
		if stmt.DataType != nil && stmt.DataType.Unsigned {
			maximum = hi
		}
	}
	if stmt.MinValue != nil {
		minimum = clamp(p.sequenceValue(stmt.MinValue))
	}
	if stmt.MaxValue != nil {
		maximum = clamp(p.sequenceValue(stmt.MaxValue))
	}
	start := minimum
	if descending {
		start = maximum
	}
	if stmt.Start != nil {
		start = p.sequenceValue(stmt.Start)
	}
	// cache < (LONGLONG_MAX - max_increment) / max_increment, where an increment of 0 counts as
	// MAX_AUTO_INCREMENT_VALUE and the default cache is 1000.
	magnitude := new(big.Int).Abs(increment)
	if magnitude.Sign() == 0 {
		magnitude = big.NewInt(65535)
	}
	cache := big.NewInt(1000)
	if stmt.NoCache {
		cache = big.NewInt(0)
	}
	if stmt.Cache != nil {
		cache = p.sequenceValue(stmt.Cache)
	}
	longMax := big.NewInt(math.MaxInt64)
	limit := new(big.Int).Sub(longMax, magnitude)
	limit.Div(limit, magnitude)
	if maximum.Cmp(start) >= 0 && maximum.Cmp(minimum) > 0 && start.Cmp(minimum) >= 0 && cache.Cmp(limit) < 0 {
		return nil
	}
	return &ParseError{Message: "Sequence '" + stmt.Name.Name + "' has out of range value for options", Position: stmt.Loc.Start}
}
