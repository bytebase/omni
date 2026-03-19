package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// DefineStmt: CREATE AGGREGATE / OPERATOR / TYPE / TEXT SEARCH / COLLATION
// ---------------------------------------------------------------------------

// parseDefineStmtAggregate parses CREATE [OR REPLACE] AGGREGATE ...
// CREATE has already been consumed. Current token is AGGREGATE.
func (p *Parser) parseDefineStmtAggregate(replace bool) (nodes.Node, error) {
	p.advance() // consume AGGREGATE
	defnames, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}

	// Save state to try new-style first (aggr_args definition).
	savedCur := p.cur
	savedPrev := p.prev
	savedNext := p.nextBuf
	savedHasNext := p.hasNext
	savedLexerErr := p.lexer.Err
	savedLexerPos := p.lexer.pos
	savedLexerStart := p.lexer.start
	savedLexerState := p.lexer.state

	args, _ := p.parseAggrArgs()
	if args != nil && p.cur.Type == '(' {
		def, err := p.parseDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.DefineStmt{
			Kind:       nodes.OBJECT_AGGREGATE,
			Oldstyle:   false,
			Replace:    replace,
			Defnames:   defnames,
			Args:       args,
			Definition: def,
		}, nil
	}

	// Restore and try old-style
	p.cur = savedCur
	p.prev = savedPrev
	p.nextBuf = savedNext
	p.hasNext = savedHasNext
	p.lexer.Err = savedLexerErr
	p.lexer.pos = savedLexerPos
	p.lexer.start = savedLexerStart
	p.lexer.state = savedLexerState

	def, err := p.parseOldAggrDefinition()
	if err != nil {
		return nil, err
	}
	return &nodes.DefineStmt{
		Kind:       nodes.OBJECT_AGGREGATE,
		Oldstyle:   true,
		Replace:    replace,
		Defnames:   defnames,
		Definition: def,
	}, nil
}

// parseDefineStmtOperator parses CREATE OPERATOR ...
// CREATE has already been consumed. Current token is OPERATOR.
func (p *Parser) parseDefineStmtOperator() (nodes.Node, error) {
	p.advance() // consume OPERATOR
	if p.cur.Type == CLASS {
		return p.parseCreateOpClassStmt()
	}
	if p.cur.Type == FAMILY {
		return p.parseCreateOpFamilyStmt()
	}
	opname, err := p.parseAnyOperator()
	if err != nil {
		return nil, err
	}
	def, err := p.parseDefinition()
	if err != nil {
		return nil, err
	}
	return &nodes.DefineStmt{
		Kind:       nodes.OBJECT_OPERATOR,
		Defnames:   opname,
		Definition: def,
	}, nil
}

// parseDefineStmtType parses CREATE TYPE ...
// CREATE has already been consumed. Current token is TYPE.
func (p *Parser) parseDefineStmtType() (nodes.Node, error) {
	p.advance() // consume TYPE
	typeName, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case AS:
		p.advance()
		switch p.cur.Type {
		case ENUM_P:
			p.advance()
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			vals := p.parseOptEnumValList()
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.CreateEnumStmt{TypeName: typeName, Vals: vals}, nil
		case RANGE:
			p.advance()
			params, err := p.parseDefinition()
			if err != nil {
				return nil, err
			}
			return &nodes.CreateRangeStmt{TypeName: typeName, Params: params}, nil
		default:
			// CompositeTypeStmt
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			var coldeflist *nodes.List
			if p.cur.Type != ')' {
				coldeflist, err = p.parseTableFuncElementList()
				if err != nil {
					return nil, err
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.CompositeTypeStmt{
				Typevar:    makeRangeVarFromAnyName(typeName),
				Coldeflist: coldeflist,
			}, nil
		}
	case '(':
		def, err := p.parseDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.DefineStmt{Kind: nodes.OBJECT_TYPE, Defnames: typeName, Definition: def}, nil
	default:
		return &nodes.DefineStmt{Kind: nodes.OBJECT_TYPE, Defnames: typeName}, nil
	}
}

// parseDefineStmtTextSearch parses CREATE TEXT SEARCH ...
// CREATE has already been consumed. Current token is TEXT.
func (p *Parser) parseDefineStmtTextSearch() (nodes.Node, error) {
	p.advance() // consume TEXT
	if _, err := p.expect(SEARCH); err != nil {
		return nil, err
	}
	var kind nodes.ObjectType
	switch p.cur.Type {
	case PARSER:
		kind = nodes.OBJECT_TSPARSER
	case DICTIONARY:
		kind = nodes.OBJECT_TSDICTIONARY
	case TEMPLATE:
		kind = nodes.OBJECT_TSTEMPLATE
	case CONFIGURATION:
		kind = nodes.OBJECT_TSCONFIGURATION
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	defnames, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	def, err := p.parseDefinition()
	if err != nil {
		return nil, err
	}
	return &nodes.DefineStmt{Kind: kind, Defnames: defnames, Definition: def}, nil
}

// parseDefineStmtCollation parses CREATE COLLATION ...
// CREATE has already been consumed. Current token is COLLATION.
func (p *Parser) parseDefineStmtCollation() (nodes.Node, error) {
	p.advance() // consume COLLATION
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
	defnames, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == FROM {
		p.advance()
		fromName, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		return &nodes.DefineStmt{
			Kind:        nodes.OBJECT_COLLATION,
			Defnames:    defnames,
			Definition:  &nodes.List{Items: []nodes.Node{makeDefElem("from", fromName)}},
			IfNotExists: ifNotExists,
		}, nil
	}
	def, err := p.parseDefinition()
	if err != nil {
		return nil, err
	}
	return &nodes.DefineStmt{
		Kind:        nodes.OBJECT_COLLATION,
		Defnames:    defnames,
		Definition:  def,
		IfNotExists: ifNotExists,
	}, nil
}

// ---------------------------------------------------------------------------
// definition / def_list / def_elem / def_arg
// ---------------------------------------------------------------------------

func (p *Parser) parseDefinition() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseDefList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

func (p *Parser) parseDefList() (*nodes.List, error) {
	elem, err := p.parseDefElem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err := p.parseDefElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseDefElem() (*nodes.DefElem, error) {
	label, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == '=' {
		p.advance()
		arg, err := p.parseDefArg()
		if err != nil {
			return nil, err
		}
		return makeDefElem(label, arg), nil
	}
	return makeDefElem(label, nil), nil
}

// parseDefArg parses def_arg:
//   func_type | reserved_keyword | qual_all_Op | NumericOnly | Sconst | NONE
func (p *Parser) parseDefArg() (nodes.Node, error) {
	if p.cur.Type == NONE {
		p.advance()
		return &nodes.String{Str: "none"}, nil
	}
	if p.cur.Type == SCONST {
		s := p.cur.Str
		p.advance()
		return &nodes.String{Str: s}, nil
	}
	if p.cur.Type == ICONST || p.cur.Type == FCONST {
		return p.parseNumericOnly(), nil
	}
	if (p.cur.Type == '+' || p.cur.Type == '-') && (p.peekNext().Type == ICONST || p.peekNext().Type == FCONST) {
		return p.parseNumericOnly(), nil
	}
	if p.cur.Type == OPERATOR && p.peekNext().Type == '(' {
		return p.parseQualAllOp()
	}
	if p.cur.Type == Op {
		op := p.cur.Str
		p.advance()
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}}, nil
	}
	switch p.cur.Type {
	case '*', '/', '%', '^', '<', '>', '=', LESS_EQUALS, GREATER_EQUALS, NOT_EQUALS:
		opStr, _ := p.parseMathOp()
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}, nil
	}
	if p.isReservedKeyword() {
		s := p.cur.Str
		p.advance()
		return &nodes.String{Str: s}, nil
	}
	tn, err := p.parseFuncType()
	if err != nil {
		return nil, err
	}
	return tn, nil
}

func (p *Parser) parseQualAllOp() (nodes.Node, error) {
	if p.cur.Type == OPERATOR {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		op, err := p.parseAnyOperator()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return op, nil
	}
	opStr, _ := p.parseAllOp()
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}, nil
}

func (p *Parser) isReservedKeyword() bool {
	if kw := LookupKeyword(p.cur.Str); kw != nil && kw.Token == p.cur.Type {
		return kw.Category == ReservedKeyword
	}
	return false
}

// ---------------------------------------------------------------------------
// old_aggr_definition / old_aggr_list / old_aggr_elem
// ---------------------------------------------------------------------------

func (p *Parser) parseOldAggrDefinition() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseOldAggrList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

func (p *Parser) parseOldAggrList() (*nodes.List, error) {
	elem, err := p.parseOldAggrElem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err := p.parseOldAggrElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseOldAggrElem() (*nodes.DefElem, error) {
	name := p.cur.Str
	p.advance()
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	arg, err := p.parseDefArg()
	if err != nil {
		return nil, err
	}
	return makeDefElem(name, arg), nil
}

// ---------------------------------------------------------------------------
// aggr_args_list helper and makeOrderedSetArgs
// ---------------------------------------------------------------------------

func (p *Parser) parseAggrArgsList() *nodes.List {
	arg := p.parseFuncArg()
	items := []nodes.Node{arg}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseFuncArg())
	}
	return &nodes.List{Items: items}
}

func makeOrderedSetArgs(directArgs *nodes.List, orderedArgs *nodes.List) *nodes.List {
	combined := make([]nodes.Node, 0)
	if directArgs != nil {
		combined = append(combined, directArgs.Items...)
	}
	directCount := 0
	if directArgs != nil {
		directCount = len(directArgs.Items)
	}
	if orderedArgs != nil {
		combined = append(combined, orderedArgs.Items...)
	}
	return &nodes.List{Items: []nodes.Node{
		&nodes.List{Items: combined},
		&nodes.Integer{Ival: int64(directCount)},
	}}
}

// ---------------------------------------------------------------------------
// TableFuncElementList (for CompositeTypeStmt)
// parseTableFuncElement is in alter_misc.go
// ---------------------------------------------------------------------------

func (p *Parser) parseTableFuncElementList() (*nodes.List, error) {
	elem, err := p.parseTableFuncElement()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err := p.parseTableFuncElement()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

// ---------------------------------------------------------------------------
// CreateOpClassStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateOpClassStmt() (nodes.Node, error) {
	p.advance() // consume CLASS
	opclassname, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	isDefault := false
	if p.cur.Type == DEFAULT {
		p.advance()
		isDefault = true
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	if _, err := p.expect(TYPE_P); err != nil {
		return nil, err
	}
	datatype, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(USING); err != nil {
		return nil, err
	}
	amname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	var opfamilyname *nodes.List
	if p.cur.Type == FAMILY {
		p.advance()
		opfamilyname, err = p.parseAnyName()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	items, err := p.parseOpclassItemList()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateOpClassStmt{
		Opclassname:  opclassname,
		IsDefault:    isDefault,
		Datatype:     datatype,
		Amname:       amname,
		Opfamilyname: opfamilyname,
		Items:        items,
	}, nil
}

func (p *Parser) parseOpclassItemList() (*nodes.List, error) {
	item, err := p.parseOpclassItem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{item}
	for p.cur.Type == ',' {
		p.advance()
		item, err := p.parseOpclassItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseOpclassItem() (nodes.Node, error) {
	switch p.cur.Type {
	case OPERATOR:
		p.advance()
		number := int(p.cur.Ival)
		p.advance()
		opname, err := p.parseAnyOperator()
		if err != nil {
			return nil, err
		}
		var owa *nodes.ObjectWithArgs
		if p.cur.Type == '(' {
			argtypes, err := p.parseOperArgtypes()
			if err != nil {
				return nil, err
			}
			owa = &nodes.ObjectWithArgs{Objname: opname, Objargs: argtypes}
		} else {
			owa = &nodes.ObjectWithArgs{Objname: opname}
		}
		orderFamily, err := p.parseOpclassPurpose()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == RECHECK {
			p.advance()
		}
		return &nodes.CreateOpClassItem{
			Itemtype:    nodes.OPCLASS_ITEM_OPERATOR,
			Name:        owa,
			Number:      number,
			OrderFamily: orderFamily,
		}, nil
	case FUNCTION:
		p.advance()
		number := int(p.cur.Ival)
		p.advance()
		// If '(' immediately after Iconst, it's '(' type_list ')' function_with_argtypes
		var classArgs *nodes.List
		if p.cur.Type == '(' {
			p.advance()
			var err error
			classArgs, err = p.parseTypeList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		fwa, err := p.parseFunctionWithArgtypes()
		if err != nil {
			return nil, err
		}
		return &nodes.CreateOpClassItem{
			Itemtype:  nodes.OPCLASS_ITEM_FUNCTION,
			Name:      fwa,
			Number:    number,
			ClassArgs: classArgs,
		}, nil
	case STORAGE:
		p.advance()
		storedtype, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		return &nodes.CreateOpClassItem{
			Itemtype:   nodes.OPCLASS_ITEM_STORAGETYPE,
			Storedtype: storedtype,
		}, nil
	default:
		return nil, nil
	}
}

func (p *Parser) parseOpclassPurpose() (*nodes.List, error) {
	if p.cur.Type != FOR {
		return nil, nil
	}
	p.advance()
	if p.cur.Type == SEARCH {
		p.advance()
		return nil, nil
	}
	if _, err := p.expect(ORDER); err != nil {
		return nil, err
	}
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}
	name, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	return name, nil
}

// ---------------------------------------------------------------------------
// CreateOpFamilyStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateOpFamilyStmt() (nodes.Node, error) {
	p.advance() // consume FAMILY
	opfamilyname, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(USING); err != nil {
		return nil, err
	}
	amname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateOpFamilyStmt{Opfamilyname: opfamilyname, Amname: amname}, nil
}

// ---------------------------------------------------------------------------
// AlterOpFamilyStmt (ADD/DROP items)
// This handles ALTER OPERATOR FAMILY name USING name ADD/DROP ...
// which is distinct from the ALTER ... SET SCHEMA / OWNER / RENAME
// handled by alter_misc.go's parseAlterOperatorClassOrFamily.
// ---------------------------------------------------------------------------

func (p *Parser) parseAlterOpFamilyAddDrop(names *nodes.List, amname string) (nodes.Node, error) {
	if p.cur.Type == ADD_P {
		p.advance()
		items, err := p.parseOpclassItemList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterOpFamilyStmt{
			Opfamilyname: names,
			Amname:       amname,
			IsDrop:       false,
			Items:        items,
		}, nil
	}
	if _, err := p.expect(DROP); err != nil {
		return nil, err
	}
	items, err := p.parseOpclassDropList()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterOpFamilyStmt{
		Opfamilyname: names,
		Amname:       amname,
		IsDrop:       true,
		Items:        items,
	}, nil
}

func (p *Parser) parseOpclassDropList() (*nodes.List, error) {
	item, err := p.parseOpclassDrop()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{item}
	for p.cur.Type == ',' {
		p.advance()
		item, err := p.parseOpclassDrop()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseOpclassDrop() (nodes.Node, error) {
	var itemtype int
	switch p.cur.Type {
	case OPERATOR:
		itemtype = nodes.OPCLASS_ITEM_OPERATOR
	case FUNCTION:
		itemtype = nodes.OPCLASS_ITEM_FUNCTION
	default:
		return nil, nil
	}
	p.advance()
	number := int(p.cur.Ival)
	p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	classArgs, err := p.parseTypeList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateOpClassItem{
		Itemtype:  itemtype,
		Number:    number,
		ClassArgs: classArgs,
	}, nil
}

// ---------------------------------------------------------------------------
// CreateStatsStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateStatsStmt() (nodes.Node, error) {
	p.advance() // consume STATISTICS
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		defnames, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		statTypes, err := p.parseOptStatNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		exprs, err := p.parseStatsParams()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(FROM); err != nil {
			return nil, err
		}
		relations, err := p.parseFromList()
		if err != nil {
			return nil, err
		}
		return &nodes.CreateStatsStmt{
			Defnames: defnames, StatTypes: statTypes,
			Exprs: exprs, Relations: relations, IfNotExists: true,
		}, nil
	}
	var defnames *nodes.List
	if p.cur.Type != ON && p.cur.Type != '(' {
		var err error
		defnames, err = p.parseAnyName()
		if err != nil {
			return nil, err
		}
	}
	statTypes, err := p.parseOptStatNameList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	exprs, err := p.parseStatsParams()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	relations, err := p.parseFromList()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateStatsStmt{
		Defnames: defnames, StatTypes: statTypes,
		Exprs: exprs, Relations: relations, IfNotExists: false,
	}, nil
}

func (p *Parser) parseOptStatNameList() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

func (p *Parser) parseStatsParams() (*nodes.List, error) {
	param, err := p.parseStatsParam()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{param}
	for p.cur.Type == ',' {
		p.advance()
		param, err := p.parseStatsParam()
		if err != nil {
			return nil, err
		}
		items = append(items, param)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseStatsParam() (nodes.Node, error) {
	if p.cur.Type == '(' {
		p.advance()
		expr, err := p.parseAExpr(0)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.StatsElem{Expr: expr}, nil
	}
	if p.isColId() && p.peekNext().Type != '(' {
		name, err := p.parseColId()
		if err != nil {
			return nil, err
		}
		return &nodes.StatsElem{Name: name}, nil
	}
	expr, err := p.parseFuncExprWindowless()
	if err != nil {
		return nil, err
	}
	return &nodes.StatsElem{Expr: expr}, nil
}

// ---------------------------------------------------------------------------
// Helper functions that define.go needs
// ---------------------------------------------------------------------------

// parseAggrArgs parses aggregate arguments.
func (p *Parser) parseAggrArgs() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()
	if p.cur.Type == '*' {
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: []nodes.Node{nil, &nodes.Integer{Ival: -1}}}, nil
	}
	if p.cur.Type == ORDER {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		args := p.parseAggrArgsList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: 0}}}, nil
	}
	args := p.parseAggrArgsList()
	if p.cur.Type == ORDER {
		p.advance()
		if _, err := p.expect(BY); err != nil {
			return nil, err
		}
		orderedArgs := p.parseAggrArgsList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return makeOrderedSetArgs(args, orderedArgs), nil
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: -1}}}, nil
}

// parseTableFuncElement parses: ColId Typename opt_collate_clause
func (p *Parser) parseTableFuncElement() (*nodes.ColumnDef, error) {
	colname, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	typename, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	coldef := &nodes.ColumnDef{
		Colname:  colname,
		TypeName: typename,
		IsLocal:  true,
		Loc: nodes.NoLoc(),
	}
	if p.cur.Type == COLLATE {
		p.advance()
		collname, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		coldef.CollClause = &nodes.CollateClause{
			Collname: collname,
			Loc: nodes.NoLoc(),
		}
	}
	return coldef, nil
}

// parseOperArgtypes parses oper_argtypes: '(' Typename ',' Typename ')' etc.
func (p *Parser) parseOperArgtypes() (*nodes.List, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if p.cur.Type == NONE {
		p.advance()
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		right, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: []nodes.Node{nil, right}}, nil
	}
	left, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == ')' {
		p.advance()
		return &nodes.List{Items: []nodes.Node{left}}, nil
	}
	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	if p.cur.Type == NONE {
		p.advance()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.List{Items: []nodes.Node{left, nil}}, nil
	}
	right, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{left, right}}, nil
}

// parseFunctionWithArgtypes parses function_with_argtypes.
func (p *Parser) parseFunctionWithArgtypes() (*nodes.ObjectWithArgs, error) {
	funcname, err := p.parseFuncName()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			p.advance()
			return &nodes.ObjectWithArgs{Objname: funcname, Objargs: nil}, nil
		}
		args := p.parseFuncArgsList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		argtypes := extractArgTypes(args)
		return &nodes.ObjectWithArgs{Objname: funcname, Objargs: argtypes}, nil
	}
	return &nodes.ObjectWithArgs{Objname: funcname, ArgsUnspecified: true}, nil
}

// parseFuncArgsList parses: func_arg (',' func_arg)*
func (p *Parser) parseFuncArgsList() *nodes.List {
	arg := p.parseFuncArg()
	items := []nodes.Node{arg}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseFuncArg())
	}
	return &nodes.List{Items: items}
}

// extractArgTypes extracts TypeName nodes from a list of FunctionParameter.
func extractArgTypes(args *nodes.List) *nodes.List {
	if args == nil {
		return nil
	}
	var types []nodes.Node
	for _, item := range args.Items {
		if fp, ok := item.(*nodes.FunctionParameter); ok {
			if fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
				continue
			}
			types = append(types, fp.ArgType)
		}
	}
	if len(types) == 0 {
		return nil
	}
	return &nodes.List{Items: types}
}

// parseAlterDefaultPrivilegesStmt parses ALTER DEFAULT PRIVILEGES ...
// ALTER has already been consumed. Current token is DEFAULT.
func (p *Parser) parseAlterDefaultPrivilegesStmt() (nodes.Node, error) {
	p.advance() // consume DEFAULT
	if _, err := p.expect(PRIVILEGES); err != nil {
		return nil, err
	}
	var optItems []nodes.Node
	for {
		if p.cur.Type == IN_P {
			p.advance()
			if _, err := p.expect(SCHEMA); err != nil {
				return nil, err
			}
			names, err := p.parseNameList()
			if err != nil {
				return nil, err
			}
			optItems = append(optItems, &nodes.DefElem{Defname: "schemas", Arg: names, Loc: nodes.NoLoc()})
		} else if p.cur.Type == FOR {
			p.advance()
			if p.cur.Type == ROLE || p.cur.Type == USER {
				p.advance()
			}
			roles := p.parseRoleList()
			optItems = append(optItems, &nodes.DefElem{Defname: "roles", Arg: roles, Loc: nodes.NoLoc()})
		} else {
			break
		}
	}
	var options *nodes.List
	if len(optItems) > 0 {
		options = &nodes.List{Items: optItems}
	}
	var action *nodes.GrantStmt
	if p.cur.Type == GRANT {
		p.advance()
		var err error
		action, err = p.parseDefACLGrantAction()
		if err != nil {
			return nil, err
		}
	} else if p.cur.Type == REVOKE {
		p.advance()
		var err error
		action, err = p.parseDefACLRevokeAction()
		if err != nil {
			return nil, err
		}
	}
	return &nodes.AlterDefaultPrivilegesStmt{Options: options, Action: action}, nil
}

func (p *Parser) parseDefACLGrantAction() (*nodes.GrantStmt, error) {
	privs := p.parsePrivileges()
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	objtype := p.parseDefACLPrivilegeTarget()
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	grantees := p.parseGranteeList()
	grantOption := false
	if p.cur.Type == WITH {
		p.advance()
		if _, err := p.expect(GRANT); err != nil {
			return nil, err
		}
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		grantOption = true
	}
	grantor := p.parseOptGrantedBy()
	return &nodes.GrantStmt{
		IsGrant: true, Privileges: privs, Targtype: nodes.ACL_TARGET_DEFAULTS,
		Objtype: nodes.ObjectType(objtype), Grantees: grantees,
		GrantOption: grantOption, Grantor: grantor,
	}, nil
}

func (p *Parser) parseDefACLRevokeAction() (*nodes.GrantStmt, error) {
	grantOption := false
	if p.cur.Type == GRANT {
		p.advance()
		if _, err := p.expect(OPTION); err != nil {
			return nil, err
		}
		if _, err := p.expect(FOR); err != nil {
			return nil, err
		}
		grantOption = true
	}
	privs := p.parsePrivileges()
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}
	objtype := p.parseDefACLPrivilegeTarget()
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	grantees := p.parseGranteeList()
	grantor := p.parseOptGrantedBy()
	behavior := p.parseOptDropBehavior()
	return &nodes.GrantStmt{
		IsGrant: false, GrantOption: grantOption, Privileges: privs,
		Targtype: nodes.ACL_TARGET_DEFAULTS, Objtype: nodes.ObjectType(objtype),
		Grantees: grantees, Grantor: grantor, Behavior: nodes.DropBehavior(behavior),
	}, nil
}

func (p *Parser) parseDefACLPrivilegeTarget() int {
	switch p.cur.Type {
	case TABLES:
		p.advance()
		return int(nodes.OBJECT_TABLE)
	case FUNCTIONS:
		p.advance()
		return int(nodes.OBJECT_FUNCTION)
	case ROUTINES:
		p.advance()
		return int(nodes.OBJECT_FUNCTION)
	case SEQUENCES:
		p.advance()
		return int(nodes.OBJECT_SEQUENCE)
	case TYPES_P:
		p.advance()
		return int(nodes.OBJECT_TYPE)
	case SCHEMAS:
		p.advance()
		return int(nodes.OBJECT_SCHEMA)
	default:
		return int(nodes.OBJECT_TABLE)
	}
}

// parseAlterStatisticsStmt parses ALTER STATISTICS ...
// ALTER has already been consumed. Current token is STATISTICS.
//
// Ref: https://www.postgresql.org/docs/17/sql-alterstatistics.html
//
//	ALTER STATISTICS name RENAME TO new_name
//	ALTER STATISTICS name OWNER TO { new_owner | CURRENT_ROLE | CURRENT_USER | SESSION_USER }
//	ALTER STATISTICS name SET SCHEMA new_schema
//	ALTER STATISTICS name SET STATISTICS target
//	ALTER STATISTICS IF EXISTS name SET STATISTICS target
func (p *Parser) parseAlterStatisticsStmt() (nodes.Node, error) {
	p.advance() // consume STATISTICS

	// IF EXISTS only applies to SET STATISTICS variant
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	defnames, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_STATISTIC_EXT,
			Object:     defnames,
			Newname:    newname,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_STATISTIC_EXT,
			Object:     defnames,
			Newowner:   roleSpec,
		}, nil
	}

	if _, err := p.expect(SET); err != nil {
		return nil, err
	}

	// SET SCHEMA vs SET STATISTICS
	if p.cur.Type == SCHEMA {
		p.advance()
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_STATISTIC_EXT,
			Object:     defnames,
			Newschema:  newschema,
		}, nil
	}

	if _, err := p.expect(STATISTICS); err != nil {
		return nil, err
	}
	stmt := &nodes.AlterStatsStmt{Defnames: defnames, MissingOk: missingOk}
	if p.cur.Type == DEFAULT {
		p.advance()
	} else {
		val := p.parseSignedIconst()
		stmt.Stxstattarget = int(val)
	}
	return stmt, nil
}

// parseAlterOperatorStmt parses ALTER OPERATOR ...
// ALTER has already been consumed. Current token is OPERATOR.
func (p *Parser) parseAlterOperatorStmt() (nodes.Node, error) {
	p.advance() // consume OPERATOR
	if p.cur.Type == CLASS || p.cur.Type == FAMILY {
		// ALTER OPERATOR CLASS/FAMILY ... (SET SCHEMA, OWNER, RENAME, or ADD/DROP)
		isClass := p.cur.Type == CLASS
		p.advance()
		names, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(USING); err != nil {
			return nil, err
		}
		amName, err := p.parseName()
		if err != nil {
			return nil, err
		}
		// Check for ADD/DROP (AlterOpFamilyStmt) vs SET/OWNER/RENAME
		if !isClass && (p.cur.Type == ADD_P || p.cur.Type == DROP) {
			return p.parseAlterOpFamilyAddDrop(names, amName)
		}
		// Handle SET SCHEMA, OWNER TO, RENAME TO
		var objtype nodes.ObjectType
		if isClass {
			objtype = nodes.OBJECT_OPCLASS
		} else {
			objtype = nodes.OBJECT_OPFAMILY
		}
		object := &nodes.List{Items: append([]nodes.Node{&nodes.String{Str: amName}}, names.Items...)}
		if p.cur.Type == SET {
			next := p.peekNext()
			if next.Type == SCHEMA {
				p.advance()
				p.advance()
				newschema, err := p.parseName()
				if err != nil {
					return nil, err
				}
				return &nodes.AlterObjectSchemaStmt{ObjectType: objtype, Object: object, Newschema: newschema}, nil
			}
		}
		if p.cur.Type == OWNER {
			p.advance()
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			roleSpec := p.parseRoleSpec()
			return &nodes.AlterOwnerStmt{ObjectType: objtype, Object: object, Newowner: roleSpec}, nil
		}
		if p.cur.Type == RENAME {
			p.advance()
			if _, err := p.expect(TO); err != nil {
				return nil, err
			}
			newname, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.RenameStmt{RenameType: objtype, Object: object, Newname: newname}, nil
		}
		return nil, nil
	}
	// Plain ALTER OPERATOR op_with_argtypes ...
	opname, err := p.parseAnyOperator()
	if err != nil {
		return nil, err
	}
	operArgs, err := p.parseOperArgtypes()
	if err != nil {
		return nil, err
	}
	owa := &nodes.ObjectWithArgs{Objname: opname, Objargs: operArgs}

	switch p.cur.Type {
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance()
			p.advance()
			newschema, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_OPERATOR,
				Object:     owa,
				Newschema:  newschema,
			}, nil
		}
		// ALTER OPERATOR op SET (operator_def_list)
		p.advance() // consume SET
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		opts, err := p.parseOperatorDefList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return &nodes.AlterOperatorStmt{
			Opername: owa,
			Options:  opts,
		}, nil
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_OPERATOR,
			Object:     owa,
			Newowner:   roleSpec,
		}, nil
	default:
		return nil, nil
	}
}

// parseOperatorDefList parses: operator_def_elem (',' operator_def_elem)*
func (p *Parser) parseOperatorDefList() (*nodes.List, error) {
	elem, err := p.parseOperatorDefElem()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem, err := p.parseOperatorDefElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
	}
	return &nodes.List{Items: items}, nil
}

func (p *Parser) parseOperatorDefElem() (*nodes.DefElem, error) {
	label, err := p.parseColLabel()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type == NONE {
			p.advance()
			return &nodes.DefElem{Defname: label, Loc: nodes.NoLoc()}, nil
		}
		arg, err := p.parseDefArg()
		if err != nil {
			return nil, err
		}
		return &nodes.DefElem{Defname: label, Arg: arg, Loc: nodes.NoLoc()}, nil
	}
	return &nodes.DefElem{Defname: label, Loc: nodes.NoLoc()}, nil
}

// parseCreateConversionStmt parses CREATE [DEFAULT] CONVERSION.
// CREATE has already been consumed. Current token is DEFAULT or CONVERSION_P.
//
// Ref: https://www.postgresql.org/docs/17/sql-createconversion.html
//
//	CREATE [ DEFAULT ] CONVERSION name
//	    FOR source_encoding TO dest_encoding FROM function_name
func (p *Parser) parseCreateConversionStmt(isDef bool) (nodes.Node, error) {
	p.advance() // consume CONVERSION_P

	convName, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	srcEnc := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	if _, err := p.expect(TO); err != nil {
		return nil, err
	}
	dstEnc := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	funcName, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateConversionStmt{
		ConversionName:  convName,
		ForEncodingName: srcEnc,
		ToEncodingName:  dstEnc,
		FuncName:        funcName,
		Def:             isDef,
	}, nil
}
