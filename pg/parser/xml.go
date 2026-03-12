package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseXmlConcat parses XMLCONCAT(expr, ...).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLCONCAT ( xml [, ...] )
func (p *Parser) parseXmlConcat() nodes.Node {
	p.advance() // consume XMLCONCAT
	p.expect('(')
	args := p.parseExprListFull()
	p.expect(')')
	return &nodes.XmlExpr{
		Op:       nodes.IS_XMLCONCAT,
		Args:     args,
		Location: -1,
	}
}

// parseXmlElement parses XMLELEMENT(NAME name [, xml_attributes] [, expr, ...]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLELEMENT ( NAME name [, XMLATTRIBUTES ( attr_list )] [, content [, ...]] )
func (p *Parser) parseXmlElement() nodes.Node {
	p.advance() // consume XMLELEMENT
	p.expect('(')
	p.expect(NAME_P)
	name, _ := p.parseColLabel()

	result := &nodes.XmlExpr{
		Op:       nodes.IS_XMLELEMENT,
		Name:     name,
		Location: -1,
	}

	if p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == XMLATTRIBUTES {
			result.NamedArgs = p.parseXmlAttributes()
			if p.cur.Type == ',' {
				p.advance()
				result.Args = p.parseExprListFull()
			}
		} else {
			// expr_list starting with the current expression
			result.Args = p.parseExprListFromCurrent()
		}
	}

	p.expect(')')
	return result
}

// parseExprListFromCurrent parses a comma-separated expression list starting
// from the current position (first expression has not been consumed yet).
func (p *Parser) parseExprListFromCurrent() *nodes.List {
	return p.parseExprListFull()
}

// parseXmlExists parses XMLEXISTS(expr PASSING expr).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLEXISTS ( text PASSING [BY {REF|VALUE}] xml [BY {REF|VALUE}] )
func (p *Parser) parseXmlExists() nodes.Node {
	p.advance() // consume XMLEXISTS
	p.expect('(')
	xpath := p.parseCExpr()
	arg := p.parseXmlExistsArgument()
	p.expect(')')
	return &nodes.FuncCall{
		Funcname:   makeFuncName("pg_catalog", "xmlexists"),
		Args:       &nodes.List{Items: []nodes.Node{xpath, arg}},
		FuncFormat: int(nodes.COERCE_SQL_SYNTAX),
		Location:   -1,
	}
}

// parseXmlForest parses XMLFOREST(xml_attribute_list).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLFOREST ( content [AS name] [, ...] )
func (p *Parser) parseXmlForest() nodes.Node {
	p.advance() // consume XMLFOREST
	p.expect('(')
	attrs := p.parseXmlAttributeList()
	p.expect(')')
	return &nodes.XmlExpr{
		Op:        nodes.IS_XMLFOREST,
		NamedArgs: attrs,
		Location:  -1,
	}
}

// parseXmlParse parses XMLPARSE(DOCUMENT|CONTENT expr).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLPARSE ( { DOCUMENT | CONTENT } value )
func (p *Parser) parseXmlParse() nodes.Node {
	p.advance() // consume XMLPARSE
	p.expect('(')
	xmloption := p.parseDocumentOrContent()
	expr := p.parseAExpr(0)
	ws := p.parseXmlWhitespaceOption()
	p.expect(')')
	return &nodes.XmlExpr{
		Op:        nodes.IS_XMLPARSE,
		Args:      &nodes.List{Items: []nodes.Node{expr, makeBoolAConst(ws)}},
		Xmloption: nodes.XmlOptionType(xmloption),
		Location:  -1,
	}
}

// parseXmlPI parses XMLPI(NAME name [, expr]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLPI ( NAME name [, content] )
func (p *Parser) parseXmlPI() nodes.Node {
	p.advance() // consume XMLPI
	p.expect('(')
	p.expect(NAME_P)
	name, _ := p.parseColLabel()

	result := &nodes.XmlExpr{
		Op:       nodes.IS_XMLPI,
		Name:     name,
		Location: -1,
	}

	if p.cur.Type == ',' {
		p.advance()
		expr := p.parseAExpr(0)
		result.Args = &nodes.List{Items: []nodes.Node{expr}}
	}

	p.expect(')')
	return result
}

// parseXmlRoot parses XMLROOT(xml, VERSION expr|NO VALUE [, STANDALONE YES|NO|NO VALUE]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLROOT ( xml, VERSION {value|NO VALUE} [, STANDALONE {YES|NO|NO VALUE}] )
func (p *Parser) parseXmlRoot() nodes.Node {
	p.advance() // consume XMLROOT
	p.expect('(')
	xmlExpr := p.parseAExpr(0)
	p.expect(',')
	version := p.parseXmlRootVersion()
	standalone := p.parseOptXmlRootStandalone()
	p.expect(')')
	return &nodes.XmlExpr{
		Op:       nodes.IS_XMLROOT,
		Args:     &nodes.List{Items: []nodes.Node{xmlExpr, version, standalone}},
		Location: -1,
	}
}

// parseXmlSerialize parses XMLSERIALIZE(DOCUMENT|CONTENT expr AS typename [INDENT]).
//
// Ref: https://www.postgresql.org/docs/17/functions-xml.html
//
//	XMLSERIALIZE ( { DOCUMENT | CONTENT } value AS type [INDENT] )
func (p *Parser) parseXmlSerialize() nodes.Node {
	p.advance() // consume XMLSERIALIZE
	p.expect('(')
	xmloption := p.parseDocumentOrContent()
	expr := p.parseAExpr(0)
	p.expect(AS)
	typeName, _ := p.parseSimpleTypename()
	indent := p.parseXmlIndentOption()
	p.expect(')')
	return &nodes.XmlSerialize{
		Xmloption: nodes.XmlOptionType(xmloption),
		Expr:      expr,
		TypeName:  typeName,
		Indent:    indent != 0,
		Location:  -1,
	}
}

// parseXmlRootVersion parses VERSION a_expr | VERSION NO VALUE.
func (p *Parser) parseXmlRootVersion() nodes.Node {
	p.expect(VERSION_P)
	if p.cur.Type == NO {
		p.advance()
		p.expect(VALUE_P)
		return &nodes.A_Const{Isnull: true}
	}
	return p.parseAExpr(0)
}

// parseOptXmlRootStandalone parses optional STANDALONE YES|NO|NO VALUE.
func (p *Parser) parseOptXmlRootStandalone() nodes.Node {
	if p.cur.Type != ',' {
		return makeIntConst(int64(nodes.XML_STANDALONE_OMITTED))
	}
	p.advance() // consume ','
	p.expect(STANDALONE_P)
	switch p.cur.Type {
	case YES_P:
		p.advance()
		return makeIntConst(int64(nodes.XML_STANDALONE_YES))
	case NO:
		p.advance()
		if p.cur.Type == VALUE_P {
			p.advance()
			return makeIntConst(int64(nodes.XML_STANDALONE_NO_VALUE))
		}
		return makeIntConst(int64(nodes.XML_STANDALONE_NO))
	default:
		return makeIntConst(int64(nodes.XML_STANDALONE_OMITTED))
	}
}

// parseXmlAttributes parses XMLATTRIBUTES '(' xml_attribute_list ')'.
func (p *Parser) parseXmlAttributes() *nodes.List {
	p.advance() // consume XMLATTRIBUTES
	p.expect('(')
	list := p.parseXmlAttributeList()
	p.expect(')')
	return list
}

// parseXmlAttributeList parses a comma-separated list of xml_attribute_el.
func (p *Parser) parseXmlAttributeList() *nodes.List {
	result := &nodes.List{}
	result.Items = append(result.Items, p.parseXmlAttributeEl())
	for p.cur.Type == ',' {
		p.advance()
		result.Items = append(result.Items, p.parseXmlAttributeEl())
	}
	return result
}

// parseXmlAttributeEl parses a_expr [AS ColLabel].
func (p *Parser) parseXmlAttributeEl() nodes.Node {
	expr := p.parseAExpr(0)
	rt := &nodes.ResTarget{
		Val:      expr,
		Location: -1,
	}
	if p.cur.Type == AS {
		p.advance()
		name, _ := p.parseColLabel()
		rt.Name = name
	}
	return rt
}

// parseDocumentOrContent parses DOCUMENT | CONTENT.
func (p *Parser) parseDocumentOrContent() int64 {
	if p.cur.Type == DOCUMENT_P {
		p.advance()
		return int64(nodes.XMLOPTION_DOCUMENT)
	}
	p.expect(CONTENT_P)
	return int64(nodes.XMLOPTION_CONTENT)
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
func (p *Parser) parseXmlWhitespaceOption() int64 {
	if p.cur.Type == PRESERVE {
		p.advance()
		p.expect(WHITESPACE_P)
		return 1
	}
	if p.cur.Type == STRIP_P {
		p.advance()
		p.expect(WHITESPACE_P)
		return 0
	}
	return 0
}

// parseXmlExistsArgument parses PASSING [BY {REF|VALUE}] c_expr [BY {REF|VALUE}].
func (p *Parser) parseXmlExistsArgument() nodes.Node {
	p.expect(PASSING)
	// Check for optional leading BY REF/VALUE
	if p.cur.Type == BY {
		p.parseXmlPassingMech()
	}
	expr := p.parseCExpr()
	// Check for optional trailing BY REF/VALUE
	if p.cur.Type == BY {
		p.parseXmlPassingMech()
	}
	return expr
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
func (p *Parser) parseXmlTable() nodes.Node {
	p.advance() // consume XMLTABLE
	p.expect('(')

	var namespaces *nodes.List

	// Check for XMLNAMESPACES
	if p.cur.Type == XMLNAMESPACES {
		p.advance()
		p.expect('(')
		namespaces = p.parseXmlNamespaceList()
		p.expect(')')
		p.expect(',')
	}

	rowExpr := p.parseCExpr()
	docExpr := p.parseXmlExistsArgument()
	p.expect(COLUMNS)
	columns := p.parseXmlTableColumnList()

	p.expect(')')

	return &nodes.RangeTableFunc{
		Rowexpr:    rowExpr,
		Docexpr:    docExpr,
		Columns:    columns,
		Namespaces: namespaces,
		Location:   -1,
	}
}

// parseXmlTableColumnList parses a comma-separated list of xmltable_column_el.
func (p *Parser) parseXmlTableColumnList() *nodes.List {
	result := &nodes.List{}
	result.Items = append(result.Items, p.parseXmlTableColumnEl())
	for p.cur.Type == ',' {
		p.advance()
		result.Items = append(result.Items, p.parseXmlTableColumnEl())
	}
	return result
}

// parseXmlTableColumnEl parses a single XMLTABLE column definition.
//
//	ColId Typename [options...] | ColId FOR ORDINALITY
func (p *Parser) parseXmlTableColumnEl() nodes.Node {
	colname, _ := p.parseColId()

	// ColId FOR ORDINALITY
	if p.cur.Type == FOR {
		p.advance()
		p.expect(ORDINALITY)
		return &nodes.RangeTableFuncCol{
			Colname:       colname,
			ForOrdinality: true,
			Location:      -1,
		}
	}

	// ColId Typename [option_list]
	typeName, _ := p.parseTypename()
	fc := &nodes.RangeTableFuncCol{
		Colname:  colname,
		TypeName: typeName,
		Location: -1,
	}

	// Parse optional column options (PATH, DEFAULT, NOT NULL, NULL)
	for {
		opt := p.parseXmlTableColumnOptionEl()
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

	return fc
}

// parseXmlTableColumnOptionEl tries to parse a single xmltable column option.
// Returns nil if the current token doesn't start an option.
func (p *Parser) parseXmlTableColumnOptionEl() nodes.Node {
	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		expr := p.parseBExpr(0)
		return &nodes.DefElem{Defname: "default", Arg: expr, Location: -1}
	case NOT:
		p.advance()
		p.expect(NULL_P)
		return &nodes.DefElem{Defname: "__pg__is_not_null", Arg: &nodes.Boolean{Boolval: true}, Location: -1}
	case NULL_P:
		p.advance()
		return &nodes.DefElem{Defname: "__pg__is_not_null", Arg: &nodes.Boolean{Boolval: false}, Location: -1}
	case PATH:
		p.advance()
		expr := p.parseBExpr(0)
		return &nodes.DefElem{Defname: "path", Arg: expr, Location: -1}
	case IDENT:
		// Generic IDENT option (rare, but grammar supports it)
		name := p.cur.Str
		p.advance()
		expr := p.parseBExpr(0)
		return &nodes.DefElem{Defname: name, Arg: expr, Location: -1}
	default:
		return nil
	}
}

// parseXmlNamespaceList parses a comma-separated list of xml_namespace_el.
func (p *Parser) parseXmlNamespaceList() *nodes.List {
	result := &nodes.List{}
	result.Items = append(result.Items, p.parseXmlNamespaceEl())
	for p.cur.Type == ',' {
		p.advance()
		result.Items = append(result.Items, p.parseXmlNamespaceEl())
	}
	return result
}

// parseXmlNamespaceEl parses b_expr AS ColLabel | DEFAULT b_expr.
func (p *Parser) parseXmlNamespaceEl() nodes.Node {
	if p.cur.Type == DEFAULT {
		p.advance()
		expr := p.parseBExpr(0)
		return &nodes.ResTarget{
			Val:      expr,
			Location: -1,
		}
	}
	expr := p.parseBExpr(0)
	p.expect(AS)
	name, _ := p.parseColLabel()
	return &nodes.ResTarget{
		Name:     name,
		Val:      expr,
		Location: -1,
	}
}

// makeBoolAConst creates an A_Const representing a boolean value as "t" or "f" string.
func makeBoolAConst(val int64) nodes.Node {
	if val != 0 {
		return &nodes.A_Const{Val: &nodes.String{Str: "t"}}
	}
	return &nodes.A_Const{Val: &nodes.String{Str: "f"}}
}
