package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseXmlConcat parses XMLCONCAT(expr, ...).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLCONCAT ( xml [, ...] )
func (p *Parser) parseXmlConcat() (nodes.Node, error) {
	p.advance() // consume XMLCONCAT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	args, err := p.parseExprListFull()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.XmlExpr{
		Op:   nodes.IS_XMLCONCAT,
		Args: args,
		Loc:  nodes.NoLoc(),
	}, nil
}

// parseXmlElement parses XMLELEMENT(NAME name [, xml_attributes] [, expr, ...]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLELEMENT ( NAME name [, XMLATTRIBUTES ( attr_list )] [, content [, ...]] )
func (p *Parser) parseXmlElement() (nodes.Node, error) {
	p.advance() // consume XMLELEMENT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(NAME_P); err != nil {
		return nil, err
	}
	name, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}

	result := &nodes.XmlExpr{
		Op:   nodes.IS_XMLELEMENT,
		Name: name,
		Loc:  nodes.NoLoc(),
	}

	if p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == XMLATTRIBUTES {
			namedArgs, err := p.parseXmlAttributes()
			if err != nil {
				return nil, err
			}
			result.NamedArgs = namedArgs
			if p.cur.Type == ',' {
				p.advance()
				result.Args, _ = p.parseExprListFull()
			}
		} else {
			// expr_list starting with the current expression
			result.Args = p.parseExprListFromCurrent()
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return result, nil
}

// parseExprListFromCurrent parses a comma-separated expression list starting
// from the current position (first expression has not been consumed yet).
func (p *Parser) parseExprListFromCurrent() *nodes.List {
	result, _ := p.parseExprListFull()
	return result
}

// parseXmlExists parses XMLEXISTS(expr PASSING expr).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLEXISTS ( text PASSING [BY {REF|VALUE}] xml [BY {REF|VALUE}] )
func (p *Parser) parseXmlExists() (nodes.Node, error) {
	p.advance() // consume XMLEXISTS
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	xpath, err := p.parseCExpr()
	arg, err := p.parseXmlExistsArgument()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "xmlexists"),
		Args:       &nodes.List{Items: []nodes.Node{xpath, arg}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Loc:        nodes.NoLoc(),
	}, nil
}

// parseXmlForest parses XMLFOREST(xml_attribute_list).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLFOREST ( content [AS name] [, ...] )
func (p *Parser) parseXmlForest() (nodes.Node, error) {
	p.advance() // consume XMLFOREST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	attrs, err := p.parseXmlAttributeList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.XmlExpr{
		Op:        nodes.IS_XMLFOREST,
		NamedArgs: attrs,
		Loc:       nodes.NoLoc(),
	}, nil
}

// parseXmlParse parses XMLPARSE(DOCUMENT|CONTENT expr).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLPARSE ( { DOCUMENT | CONTENT } value )
func (p *Parser) parseXmlParse() (nodes.Node, error) {
	p.advance() // consume XMLPARSE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	xmloption, err := p.parseDocumentOrContent()
	if err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	ws, err := p.parseXmlWhitespaceOption()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.XmlExpr{
		Op:        nodes.IS_XMLPARSE,
		Args:      &nodes.List{Items: []nodes.Node{expr, makeBoolAConst(ws)}},
		Xmloption: nodes.XmlOptionType(xmloption),
		Loc:       nodes.NoLoc(),
	}, nil
}

// parseXmlPI parses XMLPI(NAME name [, expr]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLPI ( NAME name [, content] )
func (p *Parser) parseXmlPI() (nodes.Node, error) {
	p.advance() // consume XMLPI
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(NAME_P); err != nil {
		return nil, err
	}
	name, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}

	result := &nodes.XmlExpr{
		Op:   nodes.IS_XMLPI,
		Name: name,
		Loc:  nodes.NoLoc(),
	}

	if p.cur.Type == ',' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		result.Args = &nodes.List{Items: []nodes.Node{expr}}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return result, nil
}

// parseXmlRoot parses XMLROOT(xml, VERSION expr|NO VALUE [, STANDALONE YES|NO|NO VALUE]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLROOT ( xml, VERSION {value|NO VALUE} [, STANDALONE {YES|NO|NO VALUE}] )
func (p *Parser) parseXmlRoot() (nodes.Node, error) {
	p.advance() // consume XMLROOT
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	xmlExpr, err := p.parseAExpr(0)
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	version, err := p.parseXmlRootVersion()
	if err != nil {
		return nil, err
	}
	standalone, err := p.parseOptXmlRootStandalone()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.XmlExpr{
		Op:   nodes.IS_XMLROOT,
		Args: &nodes.List{Items: []nodes.Node{xmlExpr, version, standalone}},
		Loc:  nodes.NoLoc(),
	}, nil
}

// parseXmlSerialize parses XMLSERIALIZE(DOCUMENT|CONTENT expr AS typename [INDENT]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLSERIALIZE ( { DOCUMENT | CONTENT } value AS type [INDENT] )
func (p *Parser) parseXmlSerialize() (nodes.Node, error) {
	p.advance() // consume XMLSERIALIZE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	xmloption, err := p.parseDocumentOrContent()
	if err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	typeName, err := p.parseSimpleTypename()
	if err != nil {
		return nil, err
	}
	indent := p.parseXmlIndentOption()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.XmlSerialize{
		Xmloption: nodes.XmlOptionType(xmloption),
		Expr:      expr,
		TypeName:  typeName,
		Indent:    indent != 0,
		Loc:       nodes.NoLoc(),
	}, nil
}

// parseXmlRootVersion parses VERSION a_expr | VERSION NO VALUE.
func (p *Parser) parseXmlRootVersion() (nodes.Node, error) {
	if _, err := p.expect(VERSION_P); err != nil {
		return nil, err
	}
	if p.cur.Type == NO {
		p.advance()
		if _, err := p.expect(VALUE_P); err != nil {
			return nil, err
		}
		return &nodes.A_Const{Isnull: true}, nil
	}
	return p.parseAExpr(0)
}

// parseOptXmlRootStandalone parses optional STANDALONE YES|NO|NO VALUE.
func (p *Parser) parseOptXmlRootStandalone() (nodes.Node, error) {
	if p.cur.Type != ',' {
		return makeIntConst(int64(nodes.XML_STANDALONE_OMITTED)), nil
	}
	p.advance() // consume ','
	if _, err := p.expect(STANDALONE_P); err != nil {
		return nil, err
	}
	switch p.cur.Type {
	case YES_P:
		p.advance()
		return makeIntConst(int64(nodes.XML_STANDALONE_YES)), nil
	case NO:
		p.advance()
		if p.cur.Type == VALUE_P {
			p.advance()
			return makeIntConst(int64(nodes.XML_STANDALONE_NO_VALUE)), nil
		}
		return makeIntConst(int64(nodes.XML_STANDALONE_NO)), nil
	default:
		return makeIntConst(int64(nodes.XML_STANDALONE_OMITTED)), nil
	}
}

// parseXmlAttributes parses XMLATTRIBUTES '(' xml_attribute_list ')'.
func (p *Parser) parseXmlAttributes() (*nodes.List, error) {
	p.advance() // consume XMLATTRIBUTES
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseXmlAttributeList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseXmlAttributeList parses a comma-separated list of xml_attribute_el.
func (p *Parser) parseXmlAttributeList() (*nodes.List, error) {
	result := &nodes.List{}
	el, err := p.parseXmlAttributeEl()
	if err != nil {
		return nil, err
	}
	result.Items = append(result.Items, el)
	for p.cur.Type == ',' {
		p.advance()
		el, err := p.parseXmlAttributeEl()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, el)
	}
	return result, nil
}

// parseXmlAttributeEl parses a_expr [AS ColLabel].
func (p *Parser) parseXmlAttributeEl() (nodes.Node, error) {
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	rt := &nodes.ResTarget{
		Val: expr,
		Loc: nodes.NoLoc(),
	}
	if p.cur.Type == AS {
		p.advance()
		name, err := p.parseColLabel()
		if err != nil {
			return nil, err
		}
		rt.Name = name
	}
	return rt, nil
}

// parseDocumentOrContent parses DOCUMENT | CONTENT.
func (p *Parser) parseDocumentOrContent() (int64, error) {
	if p.cur.Type == DOCUMENT_P {
		p.advance()
		return int64(nodes.XMLOPTION_DOCUMENT), nil
	}
	if _, err := p.expect(CONTENT_P); err != nil {
		return 0, err
	}
	return int64(nodes.XMLOPTION_CONTENT), nil
}

// parseXmlIndentOption parses optional INDENT | NO INDENT.
func (p *Parser) parseXmlIndentOption() int64 {
	if p.cur.Type == INDENT {
		p.advance()
		return 1
	}
	if p.cur.Type == NO {
		next := p.peekNext()
		if next.Type == INDENT {
			p.advance() // consume NO
			p.advance() // consume INDENT
			return 0
		}
	}
	return 0
}

// parseXmlWhitespaceOption parses optional PRESERVE WHITESPACE | STRIP WHITESPACE.
func (p *Parser) parseXmlWhitespaceOption() (int64, error) {
	if p.cur.Type == PRESERVE {
		p.advance()
		if _, err := p.expect(WHITESPACE_P); err != nil {
			return 0, err
		}
		return 1, nil
	}
	if p.cur.Type == STRIP_P {
		p.advance()
		if _, err := p.expect(WHITESPACE_P); err != nil {
			return 0, err
		}
		return 0, nil
	}
	return 0, nil
}

// parseXmlExistsArgument parses PASSING [BY {REF|VALUE}] c_expr [BY {REF|VALUE}].
func (p *Parser) parseXmlExistsArgument() (nodes.Node, error) {
	if _, err := p.expect(PASSING); err != nil {
		return nil, err
	}
	// Check for optional leading BY REF/VALUE
	if p.cur.Type == BY {
		p.parseXmlPassingMech()
	}
	expr, err := p.parseCExpr()
	if err != nil {
		return nil, err
	}
	// Check for optional trailing BY REF/VALUE
	if p.cur.Type == BY {
		p.parseXmlPassingMech()
	}
	return expr, nil
}

// parseXmlPassingMech consumes BY REF or BY VALUE (ignored semantically).
func (p *Parser) parseXmlPassingMech() {
	p.advance() // consume BY
	if p.cur.Type == REF_P || p.cur.Type == VALUE_P {
		p.advance()
	}
}

// parseXmlTable parses XMLTABLE(...).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLTABLE ( [XMLNAMESPACES(...),] row_expr PASSING doc_expr COLUMNS col_list )
func (p *Parser) parseXmlTable() (nodes.Node, error) {
	p.advance() // consume XMLTABLE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var namespaces *nodes.List

	// Check for XMLNAMESPACES
	if p.cur.Type == XMLNAMESPACES {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var err error
		namespaces, err = p.parseXmlNamespaceList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
	}

	rowExpr, err := p.parseCExpr()
	docExpr, err := p.parseXmlExistsArgument()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(COLUMNS); err != nil {
		return nil, err
	}
	columns, err := p.parseXmlTableColumnList()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	return &nodes.RangeTableFunc{
		Rowexpr:    rowExpr,
		Docexpr:    docExpr,
		Columns:    columns,
		Namespaces: namespaces,
		Loc:        nodes.NoLoc(),
	}, nil
}

// parseXmlTableColumnList parses a comma-separated list of xmltable_column_el.
func (p *Parser) parseXmlTableColumnList() (*nodes.List, error) {
	result := &nodes.List{}
	el, err := p.parseXmlTableColumnEl()
	if err != nil {
		return nil, err
	}
	result.Items = append(result.Items, el)
	for p.cur.Type == ',' {
		p.advance()
		el, err := p.parseXmlTableColumnEl()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, el)
	}
	return result, nil
}

// parseXmlTableColumnEl parses a single XMLTABLE column definition.
//
//	ColId Typename [options...] | ColId FOR ORDINALITY
func (p *Parser) parseXmlTableColumnEl() (nodes.Node, error) {
	colname, err := p.parseColId()
	if err != nil {
		return nil, err
	}

	// ColId FOR ORDINALITY
	if p.cur.Type == FOR {
		p.advance()
		if _, err := p.expect(ORDINALITY); err != nil {
			return nil, err
		}
		return &nodes.RangeTableFuncCol{
			Colname:       colname,
			ForOrdinality: true,
			Loc:           nodes.NoLoc(),
		}, nil
	}

	// ColId Typename [option_list]
	typeName, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	fc := &nodes.RangeTableFuncCol{
		Colname:  colname,
		TypeName: typeName,
		Loc:      nodes.NoLoc(),
	}

	// Parse optional column options (PATH, DEFAULT, NOT NULL, NULL)
	for {
		opt, err := p.parseXmlTableColumnOptionEl()
		if err != nil {
			return nil, err
		}
		if opt == nil {
			break
		}
		defel := opt.(*nodes.DefElem)
		switch defel.Defname {
		case "default":
			fc.Coldefexpr = defel.Arg
		case "path":
			fc.Colexpr = defel.Arg
		case "__pg__is_not_null":
			if defel.Arg != nil {
				if b, ok := defel.Arg.(*nodes.Boolean); ok && b.Boolval {
					fc.IsNotNull = true
				}
			}
		}
	}

	return fc, nil
}

// parseXmlTableColumnOptionEl tries to parse a single xmltable column option.
// Returns nil, nil if the current token doesn't start an option.
func (p *Parser) parseXmlTableColumnOptionEl() (nodes.Node, error) {
	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		expr, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "default", Arg: expr, Loc: nodes.NoLoc()}, nil
	case NOT:
		p.advance()
		if _, err := p.expect(NULL_P); err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "__pg__is_not_null", Arg: &nodes.Boolean{Boolval: true}, Loc: nodes.NoLoc()}, nil
	case NULL_P:
		p.advance()
		return &nodes.DefElem{Defname: "__pg__is_not_null", Arg: &nodes.Boolean{Boolval: false}, Loc: nodes.NoLoc()}, nil
	case PATH:
		p.advance()
		expr, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: "path", Arg: expr, Loc: nodes.NoLoc()}, nil
	case IDENT:
		// Generic IDENT option (rare, but grammar supports it)
		name := p.cur.Str
		p.advance()
		expr, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: name, Arg: expr, Loc: nodes.NoLoc()}, nil
	default:
		return nil, nil
	}
}

// parseXmlNamespaceList parses a comma-separated list of xml_namespace_el.
func (p *Parser) parseXmlNamespaceList() (*nodes.List, error) {
	result := &nodes.List{}
	el, err := p.parseXmlNamespaceEl()
	if err != nil {
		return nil, err
	}
	result.Items = append(result.Items, el)
	for p.cur.Type == ',' {
		p.advance()
		el, err := p.parseXmlNamespaceEl()
		if err != nil {
			return nil, err
		}
		result.Items = append(result.Items, el)
	}
	return result, nil
}

// parseXmlNamespaceEl parses b_expr AS ColLabel | DEFAULT b_expr.
func (p *Parser) parseXmlNamespaceEl() (nodes.Node, error) {
	if p.cur.Type == DEFAULT {
		p.advance()
		expr, err := p.parseBExpr(0)
		if err != nil {
			return nil, err
		}
		return &nodes.ResTarget{
			Val: expr,
			Loc: nodes.NoLoc(),
		}, nil
	}
	expr, err := p.parseBExpr(0)
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	name, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	return &nodes.ResTarget{
		Name: name,
		Val:  expr,
		Loc:  nodes.NoLoc(),
	}, nil
}

// makeBoolAConst creates an A_Const representing a boolean value as "t" or "f" string.
func makeBoolAConst(val int64) nodes.Node {
	if val != 0 {
		return &nodes.A_Const{Val: &nodes.String{Str: "t"}}
	}
	return &nodes.A_Const{Val: &nodes.String{Str: "f"}}
}
