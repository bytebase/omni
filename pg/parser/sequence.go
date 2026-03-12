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
func (p *Parser) parseCreateSeqStmt(relpersistence byte) nodes.Node {
	p.expect(SEQUENCE) // consume SEQUENCE

	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}

	names, _ := p.parseQualifiedName()
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = relpersistence

	options := p.parseOptSeqOptList()

	return &nodes.CreateSeqStmt{
		Sequence:    rv,
		Options:     options,
		IfNotExists: ifNotExists,
	}
}

// parseOptSeqOptList parses an optional sequence option list.
//
//	OptSeqOptList:
//	    SeqOptList
//	    | /* EMPTY */
func (p *Parser) parseOptSeqOptList() *nodes.List {
	if p.isSeqOptStart() {
		return p.parseSeqOptList()
	}
	return nil
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
func (p *Parser) parseCreateDomainStmt() nodes.Node {
	p.advance() // consume DOMAIN

	domainname, _ := p.parseAnyName()

	// opt_as: AS | /* EMPTY */
	if p.cur.Type == AS {
		p.advance()
	}

	typname, err := p.parseTypename()
	if err != nil {
		return nil
	}

	constraints := p.parseOptColumnConstraints()

	return &nodes.CreateDomainStmt{
		Domainname:  domainname,
		Typname:     typname,
		Constraints: constraints,
	}
}

// parseCreateTypeEnumStmt parses a CREATE TYPE ... AS ENUM statement.
// The CREATE keyword has already been consumed. The current token is TYPE.
//
//	CreateEnumStmt:
//	    CREATE TYPE any_name AS ENUM '(' opt_enum_val_list ')'
func (p *Parser) parseCreateTypeEnumStmt() nodes.Node {
	p.advance() // consume TYPE

	typeName, _ := p.parseAnyName()

	p.expect(AS)
	p.expect(ENUM_P)
	p.expect('(')

	vals := p.parseOptEnumValList()

	p.expect(')')

	return &nodes.CreateEnumStmt{
		TypeName: typeName,
		Vals:     vals,
	}
}

// parseSeqOptList parses a sequence option list.
//
//	SeqOptList:
//	    SeqOptElem
//	    | SeqOptList SeqOptElem
func (p *Parser) parseSeqOptList() *nodes.List {
	var items []nodes.Node
	for {
		opt := p.parseSeqOptElem()
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
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
func (p *Parser) parseSeqOptElem() *nodes.DefElem {
	switch p.cur.Type {
	case AS:
		p.advance()
		tn, _ := p.parseSimpleTypename()
		return makeDefElem("as", tn)
	case CACHE:
		p.advance()
		val := p.parseNumericOnly()
		return makeDefElem("cache", val)
	case CYCLE:
		p.advance()
		return makeDefElem("cycle", &nodes.Boolean{Boolval: true})
	case NO:
		p.advance()
		switch p.cur.Type {
		case CYCLE:
			p.advance()
			return makeDefElem("cycle", &nodes.Boolean{Boolval: false})
		case MAXVALUE:
			p.advance()
			return makeDefElem("maxvalue", nil)
		case MINVALUE:
			p.advance()
			return makeDefElem("minvalue", nil)
		}
		return nil
	case INCREMENT:
		p.advance()
		// opt_by: BY | EMPTY
		if p.cur.Type == BY {
			p.advance()
		}
		val := p.parseNumericOnly()
		return makeDefElem("increment", val)
	case MAXVALUE:
		p.advance()
		val := p.parseNumericOnly()
		return makeDefElem("maxvalue", val)
	case MINVALUE:
		p.advance()
		val := p.parseNumericOnly()
		return makeDefElem("minvalue", val)
	case OWNED:
		p.advance()
		p.expect(BY)
		name, _ := p.parseAnyName()
		return makeDefElem("owned_by", name)
	case SEQUENCE:
		p.advance()
		p.expect(NAME_P)
		name, _ := p.parseAnyName()
		return makeDefElem("sequence_name", name)
	case START:
		p.advance()
		// opt_with: WITH | EMPTY
		if p.cur.Type == WITH {
			p.advance()
		}
		val := p.parseNumericOnly()
		return makeDefElem("start", val)
	case RESTART:
		p.advance()
		// opt_with NumericOnly or empty
		if p.cur.Type == WITH {
			p.advance()
			val := p.parseNumericOnly()
			return makeDefElem("restart", val)
		}
		if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
			val := p.parseNumericOnly()
			return makeDefElem("restart", val)
		}
		return makeDefElem("restart", nil)
	case LOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: true})
	case UNLOGGED:
		p.advance()
		return makeDefElem("logged", &nodes.Boolean{Boolval: false})
	}
	return nil
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
