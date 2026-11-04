package parser

import (
	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The reader for what goes inside a function call's parentheses: the
// quantifier prefix an aggregate's arguments carry, and the keyword-argument
// forms POSITION and the INTERVAL function take.
//
// Which quantifier an aggregate accepts is per-builtin, not uniform. The
// grammar's `in_sum_expr: opt_all expr` carries the optional ALL, and only the
// aggregates whose production reads a DISTINCT before it take one: so
// `sum(distinct all c1)` parses, `std(distinct c1)` is a syntax error where
// `std(all c1)` parses, and GROUP_CONCAT, which takes an expr_list rather than
// an in_sum_expr, rejects ALL altogether. DISTINCTROW is the lexer's synonym
// for DISTINCT. Every row of the table below is measured against a live
// mariadb:11.8.8 by TestMariaDBFunctionArgumentOracle.

// quantSet is the set of quantifier prefixes one aggregate's production
// accepts. DISTINCT and DISTINCTROW are always accepted together: the lexer
// maps DISTINCTROW to the DISTINCT token.
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

// distinctInSumExpr is the quantifier set of a production that reads DISTINCT
// and then an in_sum_expr, whose own opt_all makes DISTINCT ALL legal too.
const distinctInSumExpr = quantDistinct | quantAll | quantDistinctAll

var aggArgSpecs = map[string]aggArgSpec{
	"AVG": {quant: distinctInSumExpr, args: 1},
	"MAX": {quant: distinctInSumExpr, args: 1},
	"MIN": {quant: distinctInSumExpr, args: 1},
	"SUM": {quant: distinctInSumExpr, args: 1},
	// variance_or_standard_deviation reads an in_sum_expr with no DISTINCT
	// before it, so these seven take ALL and nothing else: mariadb:11.8.8 answers
	// 1064 for std(distinct c1).
	"STD":         {quant: quantAll, args: 1},
	"STDDEV":      {quant: quantAll, args: 1},
	"STDDEV_POP":  {quant: quantAll, args: 1},
	"STDDEV_SAMP": {quant: quantAll, args: 1},
	"VARIANCE":    {quant: quantAll, args: 1},
	"VAR_POP":     {quant: quantAll, args: 1},
	"VAR_SAMP":    {quant: quantAll, args: 1},
	// GROUP_CONCAT's arguments are an expr_list, which carries no opt_all:
	// mariadb:11.8.8 answers 1064 for group_concat(all c1).
	"GROUP_CONCAT": {quant: quantDistinct, args: argList},
	"BIT_AND":      {quant: quantAll, args: 1},
	"BIT_OR":       {quant: quantAll, args: 1},
	"BIT_XOR":      {quant: quantAll, args: 1},
	// This engine's JSON aggregates differ from MySQL's: JSON_ARRAYAGG takes a
	// DISTINCT and an argument list where MySQL takes one argument with an
	// optional ALL, and JSON_OBJECTAGG takes no quantifier at all
	// (mariadb:11.8.8 answers 1064 for json_arrayagg(all c1) and for
	// json_objectagg(all c1, c2), and parses json_arrayagg(distinct c1, c2)).
	"JSON_ARRAYAGG":  {quant: quantDistinct, args: argList},
	"JSON_OBJECTAGG": {args: 2},
}

// isStarArg reports whether expr is `*` or a qualified `t.*`. Only COUNT
// takes a star, and only in its own production: mariadb:11.8.8 answers 1064 for
// sum(*), count(t.*), abs(*) and position(* in 'a').
func isStarArg(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.StarExpr:
		return true
	case *nodes.ColumnRef:
		return e.Star
	}
	return false
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
		if arg == nil || isStarArg(arg) {
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
// three productions are each their own shape:
//
//	COUNT '(' opt_all '*' ')'
//	COUNT '(' in_sum_expr ')'
//	COUNT '(' DISTINCT expr_list ')'
//
// So DISTINCT takes a list where the other forms take one argument, and the
// star takes a bare ALL but not a DISTINCT: mariadb:11.8.8 parses `count(all *)`
// and answers 1064 for `count(distinct *)` and `count(distinct all c1)`.
func (p *Parser) parseCountFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	q, err := p.parseAggQuantifier(quantDistinct | quantAll)
	if err != nil {
		return nil, err
	}
	fc.Quantifier = q

	if (q == nodes.AggQuantNone || q == nodes.AggQuantAll) && p.cur.Type == '*' {
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
	return p.finishFuncCall(fc, true)
}

// parsePositionFunc parses POSITION(substr IN str), the '(' already consumed.
// The keyword is the only separator the grammar has — `POSITION_SYM '('
// bit_expr IN_SYM expr ')'` — so POSITION(a, b) is a syntax error, and the
// first operand stops short of a comparison: mariadb:11.8.8 answers 1064 for
// `position(1 = 1 in 'x')` and parses `position('a' + 1 in 'b')`.
func (p *Parser) parsePositionFunc(start int) (nodes.ExprNode, error) {
	substr, err := p.parseExprPrec(precBitOr)
	if err != nil {
		return nil, err
	}
	if substr == nil || isStarArg(substr) {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	str, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if str == nil || isStarArg(str) {
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

// parenGroupHasComma reports whether the paren group the current '(' opens
// holds a comma at its own level. It is the INTERVAL disambiguator: the
// grammar has both `INTERVAL_SYM '(' expr ',' expr_list ')'`, the function,
// and `INTERVAL_SYM expr interval`, the literal, whose value may be
// parenthesized. One comma tells them apart — mariadb:11.8.8 parses
// `select a + interval (1) day from t` as the literal and answers 1064 for
// `select interval(1) day`, the function with one argument.
//
// The scan is a loop, so a deeply nested group costs tokens and not stack,
// and it rewinds what it consumed.
func (p *Parser) parenGroupHasComma() bool {
	if p.cur.Type != '(' {
		return false
	}
	saved := p.saveScan()
	defer p.restoreScan(saved)
	depth := 0
	for {
		switch p.cur.Type {
		case tokEOF:
			return false
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return false
			}
		case ',':
			if depth == 1 {
				return true
			}
		}
		p.advance()
	}
}
