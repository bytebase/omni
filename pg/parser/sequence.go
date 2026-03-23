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

	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)
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
	if p.isSeqOptStart() {
		return p.parseSeqOptList()
	}
	return nil, nil
}

// isSeqOptStart returns true if the current token can start a SeqOptElem.
func (p *Parser) isSeqOptStart() bool {
	switch p.cur.Type {
	case AS, CACHE, CYCLE, NO, INCREMENT, MAXVALUE, MINVALUE,
		OWNED, SEQUENCE, START, RESTART, LOGGED, UNLOGGED:
		return true
	}
	return false
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

	vals := p.parseOptEnumValList()

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
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
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
//	    | RESTART [opt_with NumericOnly]
//	    | LOGGED | UNLOGGED
func (p *Parser) parseSeqOptElem() (*nodes.DefElem, error) {
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
		val := p.parseNumericOnly()
		return makeDefElem("cache", val), nil
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
		return nil, nil
	case INCREMENT:
		p.advance()
		// opt_by: BY | EMPTY
		if p.cur.Type == BY {
			p.advance()
		}
		val := p.parseNumericOnly()
		return makeDefElem("increment", val), nil
	case MAXVALUE:
		p.advance()
		val := p.parseNumericOnly()
		return makeDefElem("maxvalue", val), nil
	case MINVALUE:
		p.advance()
		val := p.parseNumericOnly()
		return makeDefElem("minvalue", val), nil
	case OWNED:
		p.advance()
		p.expect(BY)
		name, _ := p.parseAnyName()
		return makeDefElem("owned_by", name), nil
	case SEQUENCE:
		p.advance()
		p.expect(NAME_P)
		name, _ := p.parseAnyName()
		return makeDefElem("sequence_name", name), nil
	case START:
		p.advance()
		// opt_with: WITH | EMPTY
		if p.cur.Type == WITH {
			p.advance()
		}
		val := p.parseNumericOnly()
		return makeDefElem("start", val), nil
	case RESTART:
		p.advance()
		// opt_with NumericOnly or empty
		if p.cur.Type == WITH {
			p.advance()
			val := p.parseNumericOnly()
			return makeDefElem("restart", val), nil
		}
		if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
			val := p.parseNumericOnly()
			return makeDefElem("restart", val), nil
		}
		return makeDefElem("restart", nil), nil
	case LOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: true}), nil
	case UNLOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: false}), nil
	}
	return nil, nil
}

// parseNumericOnly parses a numeric constant (integer or float, optionally signed).
//
//	NumericOnly:
//	    FCONST | '+' FCONST | '-' FCONST
//	    | SignedIconst
func (p *Parser) parseNumericOnly() nodes.Node {
	negative := false
	if p.cur.Type == '+' {
		p.advance()
	} else if p.cur.Type == '-' {
		p.advance()
		negative = true
	}
	if p.cur.Type == ICONST {
		val := p.cur.Ival
		p.advance()
		if negative {
			val = -val
		}
		return &nodes.Integer{Ival: val}
	}
	if p.cur.Type == FCONST {
		str := p.cur.Str
		p.advance()
		if negative {
			str = "-" + str
		}
		return &nodes.Float{Fval: str}
	}
	return &nodes.Integer{Ival: 0}
}

// parseOptEnumValList parses an optional enum value list.
//
//	opt_enum_val_list:
//	    enum_val_list
//	    | /* EMPTY */
func (p *Parser) parseOptEnumValList() *nodes.List {
	if p.cur.Type == SCONST {
		return p.parseEnumValList()
	}
	return nil
}

// parseEnumValList parses a comma-separated list of string constants for enum values.
//
//	enum_val_list:
//	    Sconst
//	    | enum_val_list ',' Sconst
func (p *Parser) parseEnumValList() *nodes.List {
	tok := p.advance() // consume first Sconst
	items := []nodes.Node{&nodes.String{Str: tok.Str}}

	for p.cur.Type == ',' {
		p.advance()
		tok = p.advance() // consume Sconst
		items = append(items, &nodes.String{Str: tok.Str})
	}

	return &nodes.List{Items: items}
}
