package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseSplitRegionStmt parses
//
//	SPLIT [REGION FOR] [PARTITION] TABLE tbl [PARTITION (p, ...)]
//	    [INDEX idx] split_option
//
//	split_option:
//	    BY (v, ...), (v, ...), ...
//	  | BETWEEN (v, ...) AND (v, ...) REGIONS n
//
// Only the four SplitSyntaxOption spellings upstream has are accepted:
// nothing, REGION FOR, PARTITION, and REGION FOR PARTITION in that order.
// `SPLIT REGION PARTITION TABLE` and `SPLIT PARTITION REGION FOR TABLE` are
// syntax errors on v8.5.5, as are a table list, an index list, and an INDEX
// clause before the PARTITION (p, ...) list.
//
// p.cur is kwSPLIT.
//
// Ref: TiDB v8.5.5 parser.y:3239-3258 (SplitRegionStmt),
// parser.y:3276-3299 (SplitSyntaxOption).
func (p *Parser) parseSplitRegionStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume SPLIT

	stmt := &nodes.SplitRegionStmt{Loc: nodes.Loc{Start: start}, IndexLoc: nodes.NoLoc()}

	// Completion: after SPLIT the grammar allows REGION, PARTITION, or
	// TABLE directly.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwREGION, kwPARTITION, kwTABLE} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// SplitSyntaxOption. REGION is always followed by FOR, and PARTITION
	// here is the bare keyword, never a name list: that list belongs after
	// the table name.
	if _, ok := p.match(kwREGION); ok {
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		stmt.RegionFor = true
	}
	if _, ok := p.match(kwPARTITION); ok {
		stmt.Partition = true
	}

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	// Completion: the table being split.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = table

	// PartitionNameListOpt: PARTITION (p, ...), one or more names. An
	// empty list and a bare PARTITION are both syntax errors.
	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Table.Partitions = parts
		stmt.Table.Loc.End = p.prev.End
	}

	// Completion: INDEX, or the split option directly.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwINDEX, kwBY, kwBETWEEN} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	if _, ok := p.match(kwINDEX); ok {
		// Identifier: a plain identifier or a non-reserved keyword, so
		// `INDEX split` names an index and `INDEX primary` is a syntax
		// error (1064 on v8.5.5), where `INDEX `primary`` is not.
		name, loc, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Index = name
		stmt.IndexLoc = nodes.Loc{Start: loc, End: p.prev.End}
	}

	opt, err := p.parseSplitOption()
	if err != nil {
		return nil, err
	}
	stmt.Opt = opt

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSplitOption parses the one mandatory split specification, in either
// the enumerated BY form or the ranged BETWEEN form. The two are exclusive:
// `BY (...) BETWEEN ...` and `BETWEEN ... REGIONS n BY (...)` are both
// syntax errors on v8.5.5.
//
// Ref: TiDB v8.5.5 parser.y:3260-3274 (SplitOption).
func (p *Parser) parseSplitOption() (*nodes.SplitOption, error) {
	start := p.pos()

	// Completion: the grammar has exactly these two continuations.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwBY)
		p.addTokenCandidate(kwBETWEEN)
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwBY:
		p.advance()
		// ValuesList (parser.y:7681-7688): one or more RowValue, comma
		// separated, no trailing comma.
		var rows [][]nodes.ExprNode
		for {
			row, err := p.parseSplitRowValue()
			if err != nil {
				return nil, err
			}
			rows = append(rows, row)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		return &nodes.SplitOption{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			ValueLists: rows,
		}, nil

	case kwBETWEEN:
		p.advance()
		lower, err := p.parseSplitRowValue()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwAND); err != nil {
			return nil, err
		}
		upper, err := p.parseSplitRowValue()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwREGIONS); err != nil {
			return nil, err
		}
		num, err := p.parseSplitRegionsNum()
		if err != nil {
			return nil, err
		}
		return &nodes.SplitOption{
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
			Ranged:  true,
			Lower:   lower,
			Upper:   upper,
			Regions: num,
		}, nil

	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseSplitRowValue parses one RowValue: a parenthesized list of values,
// which may be empty. `()` is a syntax-valid row that fails when the
// statement runs (1136 for BY, 1105 for BETWEEN on v8.5.5), so the parser
// returns an empty, non-nil row for it rather than an error.
//
// A value is an expression or the DEFAULT keyword, which is what makes
// `BY (1+1)`, `BY (now())`, `BY (-1)`, and `BY (?)` valid. The ROW(...)
// constructor an INSERT row accepts is not part of this rule: `BY row('a')`
// is a syntax error on v8.5.5.
//
// Ref: TiDB v8.5.5 parser.y:7690-7717 (RowValue, ValuesOpt, Values,
// ExprOrDefault).
func (p *Parser) parseSplitRowValue() ([]nodes.ExprNode, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	row := []nodes.ExprNode{}
	if p.cur.Type != ')' {
		vals, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		row = vals
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return row, nil
}

// parseSplitRegionsNum parses the REGIONS count. The grammar is Int64Num, so
// the value is a bare unsigned integer literal that fits int64: a sign, a
// decimal, an exponent, a hex literal, a string, a placeholder, and an
// out-of-range value are all syntax errors on v8.5.5. REGIONS 0 parses and
// fails when the statement runs.
//
// Ref: TiDB v8.5.5 parser.y:6189-6201 (Int64Num, NUM).
func (p *Parser) parseSplitRegionsNum() (int64, error) {
	if p.cur.Type != tokICONST {
		return 0, p.syntaxErrorAtCur()
	}
	tok := p.cur
	v, err := strconv.ParseInt(tok.Str, 10, 64)
	if err != nil {
		return 0, &ParseError{
			Message:  "REGIONS count out of range [0, 9223372036854775807]",
			Position: tok.Loc,
		}
	}
	p.advance()
	return v, nil
}

// parseDistributeTableStmt parses
//
//	DISTRIBUTE TABLE tbl [PARTITION (p, ...)]
//	    RULE [=] 'rule' ENGINE [=] 'engine' [TIMEOUT [=] 'timeout']
//
// RULE and ENGINE are both mandatory and fixed in this order, each value a
// string literal; TIMEOUT is optional and last. On v8.5.5 the reversed
// order, a missing RULE or ENGINE, a bare identifier or integer value, a
// comma between the options, a repeated option, an unknown option, and an
// INDEX clause are all syntax errors.
//
// p.cur is kwDISTRIBUTE.
//
// Ref: TiDB v8.5.5 parser.y:3202-3221 (DistributeTableStmt).
func (p *Parser) parseDistributeTableStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume DISTRIBUTE

	// Completion: TABLE is the only continuation.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTABLE)
		return nil, &ParseError{Message: "collecting"}
	}

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	// Completion: the table being distributed.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	table, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt := &nodes.DistributeTableStmt{Loc: nodes.Loc{Start: start}, Table: table}

	if p.cur.Type == kwPARTITION {
		p.advance()
		parts, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Table.Partitions = parts
		stmt.Table.Loc.End = p.prev.End
	}

	// Completion: RULE is the only continuation.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwRULE)
		return nil, &ParseError{Message: "collecting"}
	}

	if _, err := p.expect(kwRULE); err != nil {
		return nil, err
	}
	rule, err := p.parseDistributeOptionValue()
	if err != nil {
		return nil, err
	}
	stmt.Rule = rule

	if _, err := p.expect(kwENGINE); err != nil {
		return nil, err
	}
	engine, err := p.parseDistributeOptionValue()
	if err != nil {
		return nil, err
	}
	stmt.Engine = engine

	// Completion: TIMEOUT is the one optional continuation.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwTIMEOUT)
		return nil, &ParseError{Message: "collecting"}
	}

	if _, ok := p.match(kwTIMEOUT); ok {
		timeout, err := p.parseDistributeOptionValue()
		if err != nil {
			return nil, err
		}
		stmt.Timeout = timeout
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDistributeOptionValue reads `[=] 'value'`: EqOpt followed by a string
// literal. An identifier, an integer, and a placeholder are syntax errors on
// v8.5.5.
func (p *Parser) parseDistributeOptionValue() (string, error) {
	p.match('=')
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return tok.Str, nil
}
