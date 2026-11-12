package parser

import (
	nodes "github.com/bytebase/omni/tidb/ast"
)

// The reader for what goes inside a function call's parentheses: the
// quantifier prefix an aggregate's arguments carry, and the keyword-argument
// forms POSITION and the INTERVAL function take.
//
// Which quantifier an aggregate accepts is per-builtin, not uniform, because
// the grammar spells out one production per aggregate (pingcap/tidb
// pkg/parser/parser.y, SumExpr): AVG takes BuggyDefaultFalseDistinctOpt, so
// `avg(distinct all c1)` parses; BIT_AND takes a bare "ALL", so `bit_and(all
// c1)` parses and `bit_and(distinct c1)` is a syntax error; COUNT has four
// productions of its own. Every row of the table below is measured against a
// live TiDB v8.5.5 by TestTiDBFunctionArgumentOracle.

// quantSet is the set of quantifier prefixes one aggregate's production
// accepts. DISTINCT and DISTINCTROW are always accepted together: the grammar
// reads both through DistinctKwd.
type quantSet uint8

const (
	quantDistinct    quantSet = 1 << iota // DISTINCT, DISTINCTROW
	quantAll                              // ALL
	quantDistinctAll                      // DISTINCT ALL, DISTINCTROW ALL
)

// argList is the argument count of an aggregate that takes a list of one or
// more arguments rather than a fixed number.
const argList = -1

// aggArgSpec is one aggregate's argument grammar: the quantifier prefixes it
// accepts and how many arguments follow. COUNT is not here — its four
// productions differ in both, so parseCountFunc reads it.
type aggArgSpec struct {
	quant quantSet
	args  int
	// allPerArg marks JSON_OBJECTAGG, whose four productions put an optional
	// "ALL" before each of its two arguments. ALL is the default and the
	// engine keeps no record of which argument carried it, so the call records
	// one AggQuantAll wherever it appeared.
	allPerArg bool
}

// buggyDistinct is the quantifier set of BuggyDefaultFalseDistinctOpt: the
// name is the grammar's own (parser.y).
const buggyDistinct = quantDistinct | quantAll | quantDistinctAll

var aggArgSpecs = map[string]aggArgSpec{
	"AVG": {quant: buggyDistinct, args: 1},
	"MAX": {quant: buggyDistinct, args: 1},
	"MIN": {quant: buggyDistinct, args: 1},
	"SUM": {quant: buggyDistinct, args: 1},
	// STD and STDDEV are spellings of STDDEV_POP, VARIANCE of VAR_POP, so all
	// seven take the same production.
	"STD":          {quant: buggyDistinct, args: 1},
	"STDDEV":       {quant: buggyDistinct, args: 1},
	"STDDEV_POP":   {quant: buggyDistinct, args: 1},
	"STDDEV_SAMP":  {quant: buggyDistinct, args: 1},
	"VARIANCE":     {quant: buggyDistinct, args: 1},
	"VAR_POP":      {quant: buggyDistinct, args: 1},
	"VAR_SAMP":     {quant: buggyDistinct, args: 1},
	"GROUP_CONCAT": {quant: buggyDistinct, args: argList},
	"BIT_AND":      {quant: quantAll, args: 1},
	"BIT_OR":       {quant: quantAll, args: 1},
	"BIT_XOR":      {quant: quantAll, args: 1},
	// The JSON aggregates take ALL and not DISTINCT, and JSON_OBJECTAGG takes
	// it before either argument.
	"JSON_ARRAYAGG":  {quant: quantAll, args: 1},
	"JSON_OBJECTAGG": {quant: quantAll, args: 2, allPerArg: true},
}

// parseAggQuantifier reads the [DISTINCT|DISTINCTROW] [ALL] prefix of an
// aggregate's argument list, accepting only the forms in allowed. It consumes
// nothing and returns AggQuantNone when no quantifier is written.
func (p *Parser) parseAggQuantifier(allowed quantSet) (nodes.AggQuantifier, error) {
	switch p.cur.Type {
	case kwDISTINCT, kwDISTINCTROW:
		if allowed&(quantDistinct|quantDistinctAll) == 0 {
			return nodes.AggQuantNone, p.syntaxErrorAtCur()
		}
		distinctRow := p.cur.Type == kwDISTINCTROW
		p.advance()
		if p.cur.Type == kwALL {
			if allowed&quantDistinctAll == 0 {
				return nodes.AggQuantNone, p.syntaxErrorAtCur()
			}
			p.advance()
			if distinctRow {
				return nodes.AggQuantDistinctRowAll, nil
			}
			return nodes.AggQuantDistinctAll, nil
		}
		if allowed&quantDistinct == 0 {
			return nodes.AggQuantNone, p.syntaxErrorAtCur()
		}
		if distinctRow {
			return nodes.AggQuantDistinctRow, nil
		}
		return nodes.AggQuantDistinct, nil
	case kwALL:
		if allowed&quantAll == 0 {
			return nodes.AggQuantNone, p.syntaxErrorAtCur()
		}
		p.advance()
		return nodes.AggQuantAll, nil
	}
	return nodes.AggQuantNone, nil
}

// parseAggArgs reads an aggregate's argument list after its quantifier,
// through the closing ')', and holds it to the count spec gives. The engine
// fixes that count in the grammar, so a wrong one is a syntax error and not
// the 1582 an ordinary builtin answers.
func (p *Parser) parseAggArgs(fc *nodes.FuncCallExpr, spec aggArgSpec) error {
	for {
		if spec.allPerArg && p.cur.Type == kwALL {
			p.advance()
			fc.Quantifier = nodes.AggQuantAll
		}
		arg, err := p.parseExpr()
		if err != nil {
			return err
		}
		if arg == nil || isBareInterval(arg) {
			return p.syntaxErrorAtCur()
		}
		fc.Args = append(fc.Args, arg)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if spec.args != argList && len(fc.Args) != spec.args {
		return p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
}

// parseCountFunc parses COUNT's argument list, the '(' already consumed. Its
// four productions are each their own shape (parser.y, SumExpr):
//
//	COUNT '(' DistinctKwd ExpressionList ')'
//	COUNT '(' "ALL" Expression ')' OptWindowingClause
//	COUNT '(' Expression ')' OptWindowingClause
//	COUNT '(' '*' ')' OptWindowingClause
//
// So DISTINCT takes a list where the other forms take one argument, '*' takes
// no quantifier at all, and a de-duplicating COUNT takes no window — TiDB
// v8.5.5 answers 1064 for `count(distinct c1) over ()`, where it parses
// `sum(distinct c1) over ()`.
func (p *Parser) parseCountFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	q, err := p.parseAggQuantifier(quantDistinct | quantAll)
	if err != nil {
		return nil, err
	}
	fc.Quantifier = q

	if q == nodes.AggQuantNone && p.cur.Type == '*' {
		p.advance()
		fc.Star = true
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return p.finishFuncCall(fc, true)
	}

	spec := aggArgSpec{args: 1}
	if q.RemovesDuplicates() {
		spec.args = argList
	}
	if err := p.parseAggArgs(fc, spec); err != nil {
		return nil, err
	}
	if q.RemovesDuplicates() {
		// The de-duplicating production ends at the ')': no window follows.
		fc.Loc.End = p.prev.End
		return fc, nil
	}
	return p.finishFuncCall(fc, true)
}

// parsePositionFunc parses POSITION(substr IN str), the '(' already consumed.
// The keyword is the only separator the grammar has — `builtinPosition '('
// BitExpr "IN" Expression ')'` (parser.y) — so POSITION(a, b) is a syntax
// error, and the first operand stops short of a comparison: TiDB v8.5.5
// answers 1064 for `position(1 = 1 in 'x')` and parses `position('a' + 1 in
// 'b')`.
func (p *Parser) parsePositionFunc(start int) (nodes.ExprNode, error) {
	substr, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}
	if substr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	str, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if str == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.PositionExpr{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Substr: substr,
		Str:    str,
	}, nil
}
