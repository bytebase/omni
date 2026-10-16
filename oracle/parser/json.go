package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseJSONAggregateOrSerialize parses the arguments of JSON_ARRAYAGG,
// JSON_OBJECTAGG, or JSON_SERIALIZE into fc, the current token after the
// '('. Their clauses follow an argument where a plain call takes none, in
// an order Oracle 23ai fixes; neither aggregate takes DISTINCT (ORA-30482)
// or an OVER or KEEP clause after it.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_ARRAYAGG.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_OBJECTAGG.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/JSON_SERIALIZE.html
//
//	JSON_ARRAYAGG ( [ ALL ] expr [ FORMAT JSON ] [ order_by_clause ]
//	    [ JSON_on_null_clause ] [ JSON_returning_clause ] [ TRUNCATE ] [ STRICT ] )
//	JSON_OBJECTAGG ( [ KEY ] key_expr { VALUE | IS } expr [ FORMAT JSON ]
//	    [ JSON_on_null_clause ] [ JSON_returning_clause ] [ TRUNCATE ] [ STRICT ]
//	    [ { WITH | WITHOUT } UNIQUE KEYS ] )
//	JSON_SERIALIZE ( expr [ FORMAT JSON ] [ JSON_returning_clause ]
//	    [ { PRETTY | ASCII | TRUNCATE }... ] [ ORDERED ] [ EXTENDED ]
//	    [ { ERROR | NULL | EMPTY [ ARRAY | OBJECT ] } ON ERROR ] )
func (p *Parser) parseJSONAggregateOrSerialize(fc *nodes.FuncCallExpr) (nodes.ExprNode, error) {
	name := fc.FuncName.Name
	if name == "JSON_ARRAYAGG" && p.cur.Type == kwALL {
		p.advance()
	}
	if name == "JSON_OBJECTAGG" && p.isKeywordStr("KEY") && !p.atJSONKeyValueSeparator(p.peekNext()) {
		p.advance() // consume KEY
	}

	// The key of JSON_OBJECTAGG stops before VALUE or IS, which would
	// otherwise read as a condition.
	var arg nodes.ExprNode
	var err error
	if name == "JSON_OBJECTAGG" {
		arg, err = p.parseExprPrec(precComp)
	} else {
		arg, err = p.parseExpr()
	}
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, p.syntaxErrorAtCur()
	}
	// The key of JSON_OBJECTAGG is a scalar, not * or t.* (ORA-00936,
	// ORA-01747).
	if name == "JSON_OBJECTAGG" && p.isWildcard(arg) {
		return nil, p.syntaxErrorAtNode(arg)
	}
	fc.Args.Items = append(fc.Args.Items, arg)

	if name == "JSON_OBJECTAGG" {
		// The key takes VALUE or IS, not ':' (ORA-02000).
		if !p.isKeywordStr("VALUE") && p.cur.Type != kwIS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, p.syntaxErrorAtCur()
		}
		fc.Args.Items = append(fc.Args.Items, value)
	}

	if err := p.parseOptionalFormatJSON(); err != nil {
		return nil, err
	}

	switch name {
	case "JSON_ARRAYAGG":
		if p.cur.Type == kwORDER {
			p.advance()
			if p.cur.Type != kwBY {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			fc.OrderBy, err = p.parseOrderByList()
			if err != nil {
				return nil, err
			}
			// An item follows each comma (ORA-00936).
			if fc.OrderBy.Len() == 0 || p.prev.Type == ',' {
				return nil, p.syntaxErrorAtCur()
			}
		}
		if err := p.parseJSONConstructorClauses(false, false); err != nil {
			return nil, err
		}
	case "JSON_OBJECTAGG":
		if err := p.parseJSONConstructorClauses(true, false); err != nil {
			return nil, err
		}
	case "JSON_SERIALIZE":
		if err := p.parseJSONSerializeClauses(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	fc.Loc.End = p.prev.End
	return fc, nil
}

// atJSONKeyValueSeparator reports whether tok separates a JSON key from its
// value: VALUE, IS, or ':'. A KEY before one is a column named KEY.
func (p *Parser) atJSONKeyValueSeparator(tok Token) bool {
	return tok.Type == ':' || tok.Type == kwIS || isOneOfKeywords(tok, []string{"VALUE"})
}

// parseJSONConstructorClauses parses the clauses that end the arguments of
// a JSON constructor aggregate, in the order Oracle 23ai takes them
// (ORA-02000 for any other).
//
//	[ { NULL | ABSENT } ON NULL ] [ RETURNING datatype ] [ TRUNCATE ] [ STRICT ]
//	[ { WITH | WITHOUT } TYPENAME ]        -- with typename only
//	[ { WITH | WITHOUT } UNIQUE KEYS ]     -- with uniqueKeys only
//	[ WITH TYPENAME ]                      -- with typename only
//
// JSON_OBJECT takes one TYPENAME clause, before or after the UNIQUE KEYS
// one: Oracle 23ai takes WITH TYPENAME WITH UNIQUE KEYS and WITH UNIQUE
// KEYS WITH TYPENAME, and rejects WITH TYPENAME WITHOUT UNIQUE KEYS, WITH
// UNIQUE KEYS WITHOUT TYPENAME (ORA-02000), and a second TYPENAME clause
// (ORA-48503). Which datatypes RETURNING takes, and whether one takes
// TRUNCATE, is checked when the statement is analyzed (ORA-40449,
// ORA-40779).
func (p *Parser) parseJSONConstructorClauses(uniqueKeys, typename bool) error {
	if (p.cur.Type == kwNULL || p.isKeywordStr("ABSENT")) && p.peekNext().Type == kwON {
		p.advance() // consume NULL or ABSENT
		p.advance() // consume ON
		if p.cur.Type != kwNULL {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if err := p.parseJSONReturningClause(); err != nil {
		return err
	}
	if p.isKeywordStr("TRUNCATE") {
		p.advance()
	}
	if p.isKeywordStr("STRICT") {
		p.advance()
	}
	hasTypename := false
	if typename && (p.cur.Type == kwWITH || p.isKeywordStr("WITHOUT")) && isOneOfKeywords(p.peekNext(), []string{"TYPENAME"}) {
		p.advance() // consume WITH or WITHOUT
		p.advance() // consume TYPENAME
		hasTypename = true
	}
	if uniqueKeys && (p.cur.Type == kwWITH || p.isKeywordStr("WITHOUT") && !hasTypename) {
		p.advance()
		if !p.isKeywordStr("UNIQUE") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if !p.isKeywordStr("KEYS") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if typename && !hasTypename && p.cur.Type == kwWITH && isOneOfKeywords(p.peekNext(), []string{"TYPENAME"}) {
			p.advance() // consume WITH
			p.advance() // consume TYPENAME
		}
	}
	return nil
}

// parseJSONReturningClause parses an optional RETURNING datatype. A CLOB or
// BLOB may say how the LOB is returned, VALUE or REFERENCE (Oracle 23ai runs
// JSON_OBJECT('a' VALUE 1 RETURNING CLOB VALUE); after VARCHAR2 it is
// ORA-02000).
//
//	RETURNING datatype [ { VALUE | REFERENCE } ]
func (p *Parser) parseJSONReturningClause() error {
	if p.cur.Type != kwRETURNING {
		return nil
	}
	p.advance()
	start := p.cur
	tn, err := p.parseTypeName()
	if err != nil {
		return err
	}
	if tn == nil || tn.Names.Len() == 0 {
		return p.syntaxErrorAtCur()
	}
	if tn.Names.Len() == 1 && start.Type != tokQIDENT && (start.Type == kwCLOB || start.Type == kwBLOB) &&
		(p.isKeywordStr("VALUE") || p.isKeywordStr("REFERENCE")) {
		p.advance() // consume VALUE or REFERENCE
	}
	return nil
}

// parseJSONSerializeClauses parses the clauses of JSON_SERIALIZE after its
// argument: RETURNING, then PRETTY, ASCII, and TRUNCATE in any order and
// once each (ORA-00939 for a repeat), then ORDERED, then EXTENDED, then one
// ON ERROR clause (ORA-00907 for another order, ORA-40450 for DEFAULT or a
// second ON ERROR).
func (p *Parser) parseJSONSerializeClauses() error {
	if err := p.parseJSONReturningClause(); err != nil {
		return err
	}
	seen := map[string]bool{}
	for p.isKeywordStr("PRETTY") || p.isKeywordStr("ASCII") || p.isKeywordStr("TRUNCATE") {
		if seen[p.cur.Str] {
			return p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		p.advance()
	}
	if p.isKeywordStr("ORDERED") {
		p.advance()
	}
	if p.isKeywordStr("EXTENDED") {
		p.advance()
	}
	switch {
	case p.isKeywordStr("ERROR"), p.cur.Type == kwNULL:
		p.advance()
	case p.isKeywordStr("EMPTY"):
		p.advance()
		if p.isKeywordStr("ARRAY") || p.isKeywordStr("OBJECT") {
			p.advance()
		}
	default:
		return nil
	}
	if p.cur.Type != kwON {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isKeywordStr("ERROR") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// isWildcard reports whether expr is * or t.*. A quoted "*" or t."*" names
// a column (Oracle 23ai takes JSON_OBJECTAGG("*" VALUE 1)), so the
// reference must end in the * itself.
func (p *Parser) isWildcard(expr nodes.ExprNode) bool {
	switch e := expr.(type) {
	case *nodes.Star:
		return true
	case *nodes.ColumnRef:
		end := e.Loc.End
		return e.Column == "*" && end > 0 && end <= len(p.source) && p.source[end-1] == '*'
	}
	return false
}
