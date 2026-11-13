package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateSeqStmt parses a CREATE SEQUENCE statement.
// The CREATE keyword has already been consumed by the caller.
// The current token is SEQUENCE (or the position after OptTemp has been parsed).
//
//	CreateSeqStmt:
//	    CREATE OptTemp SEQUENCE qualified_name OptSeqOptList
//	    | CREATE OptTemp SEQUENCE IF NOT EXISTS qualified_name OptSeqOptList
func (p *Parser) parseCreateSeqStmt(startLoc int, relpersistence byte) (nodes.Node, error) {
	if _, err := p.expect(SEQUENCE); err != nil {
		return nil, err
	}

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}

	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	rv.Relpersistence = relpersistence

	options, err := p.parseOptSeqOptList()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateSeqStmt{
		Sequence:    rv,
		Options:     options,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: startLoc, End: p.prev.End},
	}, nil
}

// parseOptSeqOptList parses an optional sequence option list.
//
//	OptSeqOptList:
//	    SeqOptList
//	    | /* EMPTY */
func (p *Parser) parseOptSeqOptList() (*nodes.List, error) {
	if p.collectMode() || p.isSeqOptElemLead() {
		return p.parseSeqOptList()
	}
	return nil, nil
}

// parseCreateDomainStmt parses a CREATE DOMAIN statement.
// The CREATE keyword has already been consumed. The current token is DOMAIN.
//
//	CreateDomainStmt:
//	    CREATE DOMAIN any_name opt_as Typename opt_column_constraints
func (p *Parser) parseCreateDomainStmt(startLoc int) (nodes.Node, error) {
	p.advance() // consume DOMAIN

	domainname, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	// opt_as: AS | /* EMPTY */
	if p.cur.Type == AS {
		p.advance()
	}

	typname, err := p.parseTypename()
	if err != nil {
		return nil, err
	}

	constraints, err := p.parseOptColumnConstraints()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateDomainStmt{
		Domainname:  domainname,
		Typname:     typname,
		Constraints: constraints,
		Loc:         nodes.Loc{Start: startLoc, End: p.prev.End},
	}, nil
}

// parseCreateTypeEnumStmt parses a CREATE TYPE ... AS ENUM statement.
// The CREATE keyword has already been consumed. The current token is TYPE.
//
//	CreateEnumStmt:
//	    CREATE TYPE any_name AS ENUM '(' opt_enum_val_list ')'
func (p *Parser) parseCreateTypeEnumStmt() (nodes.Node, error) {
	p.advance() // consume TYPE

	typeName, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	if _, err := p.expect(ENUM_P); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	vals, err := p.parseOptEnumValList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.CreateEnumStmt{
		TypeName: typeName,
		Vals:     vals,
	}, nil
}

// parseSeqOptList parses a sequence option list.
//
//	SeqOptList:
//	    SeqOptElem
//	    | SeqOptList SeqOptElem
func (p *Parser) parseSeqOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt, err := p.parseSeqOptElem()
		if err != nil {
			return nil, err
		}
		items = append(items, opt)
		if !p.collectMode() && !p.isSeqOptElemLead() {
			return &nodes.List{Items: items}, nil
		}
	}
}

// parseSeqOptElem parses a single sequence option.
//
//	SeqOptElem:
//	    AS SimpleTypename
//	    | CACHE NumericOnly
//	    | CYCLE | NO CYCLE
//	    | INCREMENT opt_by NumericOnly
//	    | MAXVALUE NumericOnly | NO MAXVALUE
//	    | MINVALUE NumericOnly | NO MINVALUE
//	    | OWNED BY any_name
//	    | SEQUENCE NAME_P any_name
//	    | START opt_with NumericOnly
//	    | RESTART | RESTART opt_with NumericOnly
//	    | LOGGED | UNLOGGED
func (p *Parser) parseSeqOptElem() (*nodes.DefElem, error) {
	if p.collectMode() {
		for _, t := range seqOptElemLeadTokens {
			p.addTokenCandidate(t)
		}
		return nil, errCollecting
	}
	// numeric reads the NumericOnly of the option named name.
	numeric := func(name string) (*nodes.DefElem, error) {
		val, err := p.parseNumericOnlyRequired()
		if err != nil {
			return nil, err
		}
		return makeDefElem(name, val), nil
	}
	switch p.cur.Type {
	case AS:
		p.advance()
		tn, err := p.parseSimpleTypename()
		if err != nil {
			return nil, err
		}
		return makeDefElem("as", tn), nil
	case CACHE:
		p.advance()
		return numeric("cache")
	case CYCLE:
		p.advance()
		return makeDefElem("cycle", &nodes.Boolean{Boolval: true}), nil
	case NO:
		p.advance()
		switch p.cur.Type {
		case CYCLE:
			p.advance()
			return makeDefElem("cycle", &nodes.Boolean{Boolval: false}), nil
		case MAXVALUE:
			p.advance()
			return makeDefElem("maxvalue", nil), nil
		case MINVALUE:
			p.advance()
			return makeDefElem("minvalue", nil), nil
		}
		return nil, p.syntaxErrorAtCur()
	case INCREMENT:
		p.advance()
		// opt_by: BY | EMPTY
		if p.cur.Type == BY {
			p.advance()
		}
		return numeric("increment")
	case MAXVALUE:
		p.advance()
		return numeric("maxvalue")
	case MINVALUE:
		p.advance()
		return numeric("minvalue")
	case OWNED:
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		return makeDefElem("owned_by", name), nil
	case SEQUENCE:
		p.advance()
		if _, err := p.expect(NAME_P); err != nil {
			return nil, err
		}
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		return makeDefElem("sequence_name", name), nil
	case START:
		p.advance()
		// opt_with: WITH | EMPTY
		if p.cur.Type == WITH {
			p.advance()
		}
		return numeric("start")
	case RESTART:
		p.advance()
		// After WITH the number is required; bare RESTART takes none.
		if p.cur.Type == WITH {
			p.advance()
			return numeric("restart")
		}
		if p.isNumericOnlyStart() {
			return numeric("restart")
		}
		return makeDefElem("restart", nil), nil
	case LOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: true}), nil
	case UNLOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: false}), nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseNumericOnly parses a numeric constant (integer or float, optionally signed).
//
//	NumericOnly:
//	    FCONST | '+' FCONST | '-' FCONST
//	    | SignedIconst
func (p *Parser) parseNumericOnly() nodes.Node {
	n, err := p.parseNumericOnlyRequired()
	if err != nil {
		// Callers that predate parseNumericOnlyRequired read a missing number as 0.
		return &nodes.Integer{Ival: 0}
	}
	return n
}

// parseNumericOnlyRequired parses a NumericOnly and fails where none is.
func (p *Parser) parseNumericOnlyRequired() (nodes.Node, error) {
	negative := false
	if p.cur.Type == '+' || p.cur.Type == '-' {
		negative = p.cur.Type == '-'
		p.advance()
	}
	switch p.cur.Type {
	case ICONST:
		val := p.cur.Ival
		p.advance()
		if negative {
			val = -val
		}
		return &nodes.Integer{Ival: val}, nil
	case FCONST:
		str := p.cur.Str
		p.advance()
		if negative {
			str = "-" + str
		}
		return &nodes.Float{Fval: str}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseOptEnumValList parses an optional enum value list.
//
//	opt_enum_val_list:
//	    enum_val_list
//	    | /* EMPTY */
func (p *Parser) parseOptEnumValList() (*nodes.List, error) {
	if p.cur.Type == SCONST {
		return p.parseEnumValList()
	}
	return nil, nil
}

// parseEnumValList parses a comma-separated list of string constants for enum values.
//
//	enum_val_list:
//	    Sconst
//	    | enum_val_list ',' Sconst
func (p *Parser) parseEnumValList() (*nodes.List, error) {
	tok := p.advance() // consume first Sconst
	items := []nodes.Node{&nodes.String{Str: tok.Str}}

	for p.cur.Type == ',' {
		p.advance()
		tok, err := p.expect(SCONST)
		if err != nil {
			return nil, err
		}
		items = append(items, &nodes.String{Str: tok.Str})
	}

	return &nodes.List{Items: items}, nil
}
