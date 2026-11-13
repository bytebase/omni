package parser

import (
	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseJSONQueryFunc reads what follows `JSON_QUERY(`, which has its own arm in the grammar as
// JSON_VALUE has, and not the generic call's expression list:
//
//	JSON_QUERY ( simple_expr , string_literal [RETURNING cast_type] [TRUNCATE]
//	    [{ALLOW | DISALLOW} SCALARS] [PRETTY] [ASCII] [wrapper] [ASIS]
//	    [on_error] [on_empty] [on_mismatch] [MULTIVALUE] )
//	wrapper: WITHOUT [ARRAY] WRAPPER | WITH [CONDITIONAL | UNCONDITIONAL] [ARRAY] WRAPPER
//	on_error: {ERROR | NULL | EMPTY | EMPTY ARRAY | EMPTY OBJECT} ON ERROR
//	on_empty: {ERROR | NULL | EMPTY | EMPTY ARRAY | EMPTY OBJECT} ON EMPTY
//	on_mismatch: {ERROR | NULL | DOT} ON MISMATCH
//
// The clauses come in that order and once each, but that ON ERROR and ON EMPTY swap. A
// clause's words are keywords, so a quoted one is no clause, and quoted or qualified the
// name is an ordinary call. So OceanBase 4.3.5.6 answers 1064 for `JSON_QUERY(a)`,
// `JSON_QUERY(1, 2)`, `JSON_QUERY(a, '$' PRETTY ALLOW SCALARS)`, `JSON_QUERY(a, '$' NULL ON
// MISMATCH NULL ON EMPTY)` and `JSON_QUERY(a, '$' EMPTY ON MISMATCH)`.
//
// It is how OceanBase prints a multi-valued index: a key written (CAST(j->'$.ids' AS
// UNSIGNED ARRAY)) reads back as json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY
// WRAPPER asis error on error null on empty null on mismatch), and the readback, written
// as a key, makes the same index.
//
// Ref: json_query_expr in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, line 24917).
// TestJSONQueryFormsOceanBaseOracle holds the forms to the engine.
func (p *Parser) parseJSONQueryFunc(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	if err := p.parseJSONDocAndPath(fc); err != nil {
		return nil, err
	}

	opts := &nodes.JSONQueryOptions{Loc: nodes.Loc{Start: p.pos()}}
	if _, ok := p.match(kwRETURNING); ok {
		dt, err := p.parseCastDataType()
		if err != nil {
			return nil, err
		}
		opts.Returning = dt
	}
	if _, ok := p.match(kwTRUNCATE); ok {
		opts.Truncate = true
	}
	for _, scalars := range []string{"ALLOW", "DISALLOW"} {
		if p.isBareWord(scalars) {
			p.advance()
			if !p.isBareWord("SCALARS") {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			opts.Scalars = scalars
			break
		}
	}
	if p.isBareWord("PRETTY") {
		p.advance()
		opts.Pretty = true
	}
	if _, ok := p.match(kwASCII); ok {
		opts.ASCII = true
	}
	if p.cur.Type == kwWITH || p.cur.Type == kwWITHOUT {
		wrapper, err := p.parseJSONQueryWrapper()
		if err != nil {
			return nil, err
		}
		opts.Wrapper = wrapper
	}
	if p.isBareWord("ASIS") {
		p.advance()
		opts.Asis = true
	}
	if err := p.parseJSONQueryResponses(opts); err != nil {
		return nil, err
	}
	if p.isBareWord("MULTIVALUE") {
		p.advance()
		opts.Multivalue = true
	}

	fc.JSONQueryArm = true
	if p.prev.End > opts.Loc.Start {
		opts.Loc.End = p.prev.End
		fc.JSONQuery = opts
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	fc.Loc.End = p.prev.End
	return fc, nil
}

// parseJSONQueryWrapper reads JSON_QUERY's wrapper clause from its WITH or WITHOUT, and
// returns it in capitals. WITHOUT takes ARRAY alone before WRAPPER; WITH also takes
// CONDITIONAL or UNCONDITIONAL, before ARRAY.
func (p *Parser) parseJSONQueryWrapper() (string, error) {
	with := p.cur.Type == kwWITH
	wrapper := "WITHOUT"
	if with {
		wrapper = "WITH"
	}
	p.advance()
	if with {
		for _, condition := range []string{"CONDITIONAL", "UNCONDITIONAL"} {
			if p.isBareWord(condition) {
				p.advance()
				wrapper += " " + condition
				break
			}
		}
	}
	if _, ok := p.match(kwARRAY); ok {
		wrapper += " ARRAY"
	}
	if _, err := p.expect(kwWRAPPER); err != nil {
		return "", err
	}
	return wrapper + " WRAPPER", nil
}

// parseJSONQueryResponses reads JSON_QUERY's ON ERROR, ON EMPTY and ON MISMATCH. The first
// two come in either order and ON MISMATCH after them, each at most once.
func (p *Parser) parseJSONQueryResponses(opts *nodes.JSONQueryOptions) error {
	for {
		response, ok := p.jsonQueryResponse()
		if !ok {
			return nil
		}
		if _, err := p.expect(kwON); err != nil {
			return err
		}
		empties := response != "DOT"
		simple := response == "ERROR" || response == "NULL"
		switch {
		case p.cur.Type == kwERROR_KW && empties && opts.OnError == "":
			opts.OnError = response
		case p.cur.Type == kwEMPTY && empties && opts.OnEmpty == "":
			opts.OnEmpty = response
		case p.isBareWord("MISMATCH") && (simple || response == "DOT"):
			p.advance()
			opts.OnMismatch = response
			return nil
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
}

// jsonQueryResponse reads what JSON_QUERY answers with, before its ON: ERROR, NULL, EMPTY,
// EMPTY ARRAY, EMPTY OBJECT, or the DOT only ON MISMATCH takes.
func (p *Parser) jsonQueryResponse() (string, bool) {
	switch {
	case p.cur.Type == kwERROR_KW:
		p.advance()
		return "ERROR", true
	case p.cur.Type == kwNULL:
		p.advance()
		return "NULL", true
	case p.isBareWord("DOT"):
		p.advance()
		return "DOT", true
	case p.cur.Type == kwEMPTY:
		p.advance()
		if _, ok := p.match(kwARRAY); ok {
			return "EMPTY ARRAY", true
		}
		if p.isBareWord("OBJECT") {
			p.advance()
			return "EMPTY OBJECT", true
		}
		return "EMPTY", true
	}
	return "", false
}
