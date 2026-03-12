package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// DefineStmt: CREATE AGGREGATE / OPERATOR / TYPE / TEXT SEARCH / COLLATION
// ---------------------------------------------------------------------------

// parseDefineStmtAggregate parses CREATE [OR REPLACE] AGGREGATE ...
// CREATE has already been consumed. Current token is AGGREGATE.
func (p *Parser) parseDefineStmtAggregate(replace bool) nodes.Node {
	p.advance() // consume AGGREGATE
	defnames, _ := p.parseFuncName()

	// Save state to try new-style first (aggr_args definition).
	savedCur := p.cur
	savedPrev := p.prev
	savedNext := p.nextBuf
	savedHasNext := p.hasNext

	args := p.parseAggrArgs()
	if args != nil && p.cur.Type == '(' {
		def := p.parseDefinition()
		return &nodes.DefineStmt{
			Kind:       nodes.OBJECT_AGGREGATE,
			Oldstyle:   false,
			Replace:    replace,
			Defnames:   defnames,
			Args:       args,
			Definition: def,
		}
	}

	// Restore and try old-style
	p.cur = savedCur
	p.prev = savedPrev
	p.nextBuf = savedNext
	p.hasNext = savedHasNext

	def := p.parseOldAggrDefinition()
	return &nodes.DefineStmt{
		Kind:       nodes.OBJECT_AGGREGATE,
		Oldstyle:   true,
		Replace:    replace,
		Defnames:   defnames,
		Definition: def,
	}
}

// parseDefineStmtOperator parses CREATE OPERATOR ...
// CREATE has already been consumed. Current token is OPERATOR.
func (p *Parser) parseDefineStmtOperator() nodes.Node {
	p.advance() // consume OPERATOR
	if p.cur.Type == CLASS {
		return p.parseCreateOpClassStmt()
	}
	if p.cur.Type == FAMILY {
		return p.parseCreateOpFamilyStmt()
	}
	opname, _ := p.parseAnyOperator()
	def := p.parseDefinition()
	return &nodes.DefineStmt{
		Kind:       nodes.OBJECT_OPERATOR,
		Defnames:   opname,
		Definition: def,
	}
}

// parseDefineStmtType parses CREATE TYPE ...
// CREATE has already been consumed. Current token is TYPE.
func (p *Parser) parseDefineStmtType() nodes.Node {
	p.advance() // consume TYPE
	typeName, _ := p.parseAnyName()

	switch p.cur.Type {
	case AS:
		p.advance()
		switch p.cur.Type {
		case ENUM_P:
			p.advance()
			p.expect('(')
			vals := p.parseOptEnumValList()
			p.expect(')')
			return &nodes.CreateEnumStmt{TypeName: typeName, Vals: vals}
		case RANGE:
			p.advance()
			params := p.parseDefinition()
			return &nodes.CreateRangeStmt{TypeName: typeName, Params: params}
		default:
			// CompositeTypeStmt
			p.expect('(')
			var coldeflist *nodes.List
			if p.cur.Type != ')' {
				coldeflist = p.parseTableFuncElementList()
			}
			p.expect(')')
			return &nodes.CompositeTypeStmt{
				Typevar:    makeRangeVarFromAnyName(typeName),
				Coldeflist: coldeflist,
			}
		}
	case '(':
		def := p.parseDefinition()
		return &nodes.DefineStmt{Kind: nodes.OBJECT_TYPE, Defnames: typeName, Definition: def}
	default:
		return &nodes.DefineStmt{Kind: nodes.OBJECT_TYPE, Defnames: typeName}
	}
}

// parseDefineStmtTextSearch parses CREATE TEXT SEARCH ...
// CREATE has already been consumed. Current token is TEXT.
func (p *Parser) parseDefineStmtTextSearch() nodes.Node {
	p.advance() // consume TEXT
	p.expect(SEARCH)
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
		return nil
	}
	p.advance()
	defnames, _ := p.parseAnyName()
	def := p.parseDefinition()
	return &nodes.DefineStmt{Kind: kind, Defnames: defnames, Definition: def}
}

// parseDefineStmtCollation parses CREATE COLLATION ...
// CREATE has already been consumed. Current token is COLLATION.
func (p *Parser) parseDefineStmtCollation() nodes.Node {
	p.advance() // consume COLLATION
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}
	defnames, _ := p.parseAnyName()
	if p.cur.Type == FROM {
		p.advance()
		fromName, _ := p.parseAnyName()
		return &nodes.DefineStmt{
			Kind:        nodes.OBJECT_COLLATION,
			Defnames:    defnames,
			Definition:  &nodes.List{Items: []nodes.Node{makeDefElem("from", fromName)}},
			IfNotExists: ifNotExists,
		}
	}
	def := p.parseDefinition()
	return &nodes.DefineStmt{
		Kind:        nodes.OBJECT_COLLATION,
		Defnames:    defnames,
		Definition:  def,
		IfNotExists: ifNotExists,
	}
}

// ---------------------------------------------------------------------------
// definition / def_list / def_elem / def_arg
// ---------------------------------------------------------------------------

func (p *Parser) parseDefinition() *nodes.List {
	p.expect('(')
	list := p.parseDefList()
	p.expect(')')
	return list
}

func (p *Parser) parseDefList() *nodes.List {
	elem := p.parseDefElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseDefElem())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseDefElem() *nodes.DefElem {
	label, _ := p.parseColLabel()
	if p.cur.Type == '=' {
		p.advance()
		arg := p.parseDefArg()
		return makeDefElem(label, arg)
	}
	return makeDefElem(label, nil)
}

// parseDefArg parses def_arg:
//   func_type | reserved_keyword | qual_all_Op | NumericOnly | Sconst | NONE
func (p *Parser) parseDefArg() nodes.Node {
	if p.cur.Type == NONE {
		p.advance()
		return &nodes.String{Str: "none"}
	}
	if p.cur.Type == SCONST {
		s := p.cur.Str
		p.advance()
		return &nodes.String{Str: s}
	}
	if p.cur.Type == ICONST || p.cur.Type == FCONST {
		return p.parseNumericOnly()
	}
	if (p.cur.Type == '+' || p.cur.Type == '-') && (p.peekNext().Type == ICONST || p.peekNext().Type == FCONST) {
		return p.parseNumericOnly()
	}
	if p.cur.Type == OPERATOR && p.peekNext().Type == '(' {
		return p.parseQualAllOp()
	}
	if p.cur.Type == Op {
		op := p.cur.Str
		p.advance()
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: op}}}
	}
	switch p.cur.Type {
	case '*', '/', '%', '^', '<', '>', '=', LESS_EQUALS, GREATER_EQUALS, NOT_EQUALS:
		opStr, _ := p.parseMathOp()
		return &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
	}
	if p.isReservedKeyword() {
		s := p.cur.Str
		p.advance()
		return &nodes.String{Str: s}
	}
	tn, _ := p.parseFuncType()
	return tn
}

func (p *Parser) parseQualAllOp() nodes.Node {
	if p.cur.Type == OPERATOR {
		p.advance()
		p.expect('(')
		op, _ := p.parseAnyOperator()
		p.expect(')')
		return op
	}
	opStr, _ := p.parseAllOp()
	return &nodes.List{Items: []nodes.Node{&nodes.String{Str: opStr}}}
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

func (p *Parser) parseOldAggrDefinition() *nodes.List {
	p.expect('(')
	list := p.parseOldAggrList()
	p.expect(')')
	return list
}

func (p *Parser) parseOldAggrList() *nodes.List {
	elem := p.parseOldAggrElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseOldAggrElem())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseOldAggrElem() *nodes.DefElem {
	name := p.cur.Str
	p.advance()
	p.expect('=')
	arg := p.parseDefArg()
	return makeDefElem(name, arg)
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

func (p *Parser) parseTableFuncElementList() *nodes.List {
	elem := p.parseTableFuncElement()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseTableFuncElement())
	}
	return &nodes.List{Items: items}
}

// ---------------------------------------------------------------------------
// CreateOpClassStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateOpClassStmt() nodes.Node {
	p.advance() // consume CLASS
	opclassname, _ := p.parseAnyName()
	isDefault := false
	if p.cur.Type == DEFAULT {
		p.advance()
		isDefault = true
	}
	p.expect(FOR)
	p.expect(TYPE_P)
	datatype, _ := p.parseTypename()
	p.expect(USING)
	amname, _ := p.parseName()
	var opfamilyname *nodes.List
	if p.cur.Type == FAMILY {
		p.advance()
		opfamilyname, _ = p.parseAnyName()
	}
	p.expect(AS)
	items := p.parseOpclassItemList()
	return &nodes.CreateOpClassStmt{
		Opclassname:  opclassname,
		IsDefault:    isDefault,
		Datatype:     datatype,
		Amname:       amname,
		Opfamilyname: opfamilyname,
		Items:        items,
	}
}

func (p *Parser) parseOpclassItemList() *nodes.List {
	item := p.parseOpclassItem()
	items := []nodes.Node{item}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseOpclassItem())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseOpclassItem() nodes.Node {
	switch p.cur.Type {
	case OPERATOR:
		p.advance()
		number := int(p.cur.Ival)
		p.advance()
		opname, _ := p.parseAnyOperator()
		var owa *nodes.ObjectWithArgs
		if p.cur.Type == '(' {
			argtypes := p.parseOperArgtypes()
			owa = &nodes.ObjectWithArgs{Objname: opname, Objargs: argtypes}
		} else {
			owa = &nodes.ObjectWithArgs{Objname: opname}
		}
		orderFamily := p.parseOpclassPurpose()
		if p.cur.Type == RECHECK {
			p.advance()
		}
		return &nodes.CreateOpClassItem{
			Itemtype:    nodes.OPCLASS_ITEM_OPERATOR,
			Name:        owa,
			Number:      number,
			OrderFamily: orderFamily,
		}
	case FUNCTION:
		p.advance()
		number := int(p.cur.Ival)
		p.advance()
		// If '(' immediately after Iconst, it's '(' type_list ')' function_with_argtypes
		var classArgs *nodes.List
		if p.cur.Type == '(' {
			p.advance()
			classArgs, _ = p.parseTypeList()
			p.expect(')')
		}
		fwa := p.parseFunctionWithArgtypes()
		return &nodes.CreateOpClassItem{
			Itemtype:  nodes.OPCLASS_ITEM_FUNCTION,
			Name:      fwa,
			Number:    number,
			ClassArgs: classArgs,
		}
	case STORAGE:
		p.advance()
		storedtype, _ := p.parseTypename()
		return &nodes.CreateOpClassItem{
			Itemtype:   nodes.OPCLASS_ITEM_STORAGETYPE,
			Storedtype: storedtype,
		}
	default:
		return nil
	}
}

func (p *Parser) parseOpclassPurpose() *nodes.List {
	if p.cur.Type != FOR {
		return nil
	}
	p.advance()
	if p.cur.Type == SEARCH {
		p.advance()
		return nil
	}
	p.expect(ORDER)
	p.expect(BY)
	name, _ := p.parseAnyName()
	return name
}

// ---------------------------------------------------------------------------
// CreateOpFamilyStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateOpFamilyStmt() nodes.Node {
	p.advance() // consume FAMILY
	opfamilyname, _ := p.parseAnyName()
	p.expect(USING)
	amname, _ := p.parseName()
	return &nodes.CreateOpFamilyStmt{Opfamilyname: opfamilyname, Amname: amname}
}

// ---------------------------------------------------------------------------
// AlterOpFamilyStmt (ADD/DROP items)
// This handles ALTER OPERATOR FAMILY name USING name ADD/DROP ...
// which is distinct from the ALTER ... SET SCHEMA / OWNER / RENAME
// handled by alter_misc.go's parseAlterOperatorClassOrFamily.
// ---------------------------------------------------------------------------

func (p *Parser) parseAlterOpFamilyAddDrop(names *nodes.List, amname string) nodes.Node {
	if p.cur.Type == ADD_P {
		p.advance()
		items := p.parseOpclassItemList()
		return &nodes.AlterOpFamilyStmt{
			Opfamilyname: names,
			Amname:       amname,
			IsDrop:       false,
			Items:        items,
		}
	}
	p.expect(DROP)
	items := p.parseOpclassDropList()
	return &nodes.AlterOpFamilyStmt{
		Opfamilyname: names,
		Amname:       amname,
		IsDrop:       true,
		Items:        items,
	}
}

func (p *Parser) parseOpclassDropList() *nodes.List {
	item := p.parseOpclassDrop()
	items := []nodes.Node{item}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseOpclassDrop())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseOpclassDrop() nodes.Node {
	var itemtype int
	switch p.cur.Type {
	case OPERATOR:
		itemtype = nodes.OPCLASS_ITEM_OPERATOR
	case FUNCTION:
		itemtype = nodes.OPCLASS_ITEM_FUNCTION
	default:
		return nil
	}
	p.advance()
	number := int(p.cur.Ival)
	p.advance()
	p.expect('(')
	classArgs, _ := p.parseTypeList()
	p.expect(')')
	return &nodes.CreateOpClassItem{
		Itemtype:  itemtype,
		Number:    number,
		ClassArgs: classArgs,
	}
}

// ---------------------------------------------------------------------------
// CreateStatsStmt
// ---------------------------------------------------------------------------

func (p *Parser) parseCreateStatsStmt() nodes.Node {
	p.advance() // consume STATISTICS
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		defnames, _ := p.parseAnyName()
		statTypes := p.parseOptStatNameList()
		p.expect(ON)
		exprs := p.parseStatsParams()
		p.expect(FROM)
		relations := p.parseFromList()
		return &nodes.CreateStatsStmt{
			Defnames: defnames, StatTypes: statTypes,
			Exprs: exprs, Relations: relations, IfNotExists: true,
		}
	}
	var defnames *nodes.List
	if p.cur.Type != ON && p.cur.Type != '(' {
		defnames, _ = p.parseAnyName()
	}
	statTypes := p.parseOptStatNameList()
	p.expect(ON)
	exprs := p.parseStatsParams()
	p.expect(FROM)
	relations := p.parseFromList()
	return &nodes.CreateStatsStmt{
		Defnames: defnames, StatTypes: statTypes,
		Exprs: exprs, Relations: relations, IfNotExists: false,
	}
}

func (p *Parser) parseOptStatNameList() *nodes.List {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	names, _ := p.parseNameList()
	p.expect(')')
	return names
}

func (p *Parser) parseStatsParams() *nodes.List {
	param := p.parseStatsParam()
	items := []nodes.Node{param}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseStatsParam())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseStatsParam() nodes.Node {
	if p.cur.Type == '(' {
		p.advance()
		expr := p.parseAExpr(0)
		p.expect(')')
		return &nodes.StatsElem{Expr: expr}
	}
	if p.isColId() && p.peekNext().Type != '(' {
		name, _ := p.parseColId()
		return &nodes.StatsElem{Name: name}
	}
	expr := p.parseFuncExprWindowless()
	return &nodes.StatsElem{Expr: expr}
}

// ---------------------------------------------------------------------------
// Helper functions that define.go needs
// ---------------------------------------------------------------------------

// parseAggrArgs parses aggregate arguments.
func (p *Parser) parseAggrArgs() *nodes.List {
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	if p.cur.Type == '*' {
		p.advance()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{nil, &nodes.Integer{Ival: -1}}}
	}
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		args := p.parseAggrArgsList()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: 0}}}
	}
	args := p.parseAggrArgsList()
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		orderedArgs := p.parseAggrArgsList()
		p.expect(')')
		return makeOrderedSetArgs(args, orderedArgs)
	}
	p.expect(')')
	return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: -1}}}
}

// parseTableFuncElement parses: ColId Typename opt_collate_clause
func (p *Parser) parseTableFuncElement() *nodes.ColumnDef {
	colname, _ := p.parseColId()
	typename, _ := p.parseTypename()
	coldef := &nodes.ColumnDef{
		Colname:  colname,
		TypeName: typename,
		IsLocal:  true,
		Location: -1,
	}
	if p.cur.Type == COLLATE {
		p.advance()
		collname, _ := p.parseAnyName()
		coldef.CollClause = &nodes.CollateClause{
			Collname: collname,
			Location: -1,
		}
	}
	return coldef
}

// parseOperArgtypes parses oper_argtypes: '(' Typename ',' Typename ')' etc.
func (p *Parser) parseOperArgtypes() *nodes.List {
	p.expect('(')
	if p.cur.Type == NONE {
		p.advance()
		p.expect(',')
		right, _ := p.parseTypename()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{nil, right}}
	}
	left, _ := p.parseTypename()
	if p.cur.Type == ')' {
		p.advance()
		return &nodes.List{Items: []nodes.Node{left}}
	}
	p.expect(',')
	if p.cur.Type == NONE {
		p.advance()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{left, nil}}
	}
	right, _ := p.parseTypename()
	p.expect(')')
	return &nodes.List{Items: []nodes.Node{left, right}}
}

// parseFunctionWithArgtypes parses function_with_argtypes.
func (p *Parser) parseFunctionWithArgtypes() *nodes.ObjectWithArgs {
	funcname, _ := p.parseFuncName()
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			p.advance()
			return &nodes.ObjectWithArgs{Objname: funcname, Objargs: nil}
		}
		args := p.parseFuncArgsList()
		p.expect(')')
		argtypes := extractArgTypes(args)
		return &nodes.ObjectWithArgs{Objname: funcname, Objargs: argtypes}
	}
	return &nodes.ObjectWithArgs{Objname: funcname, ArgsUnspecified: true}
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
func (p *Parser) parseAlterDefaultPrivilegesStmt() nodes.Node {
	p.advance() // consume DEFAULT
	p.expect(PRIVILEGES)
	var optItems []nodes.Node
	for {
		if p.cur.Type == IN_P {
			p.advance()
			p.expect(SCHEMA)
			names, _ := p.parseNameList()
			optItems = append(optItems, &nodes.DefElem{Defname: "schemas", Arg: names, Location: -1})
		} else if p.cur.Type == FOR {
			p.advance()
			if p.cur.Type == ROLE || p.cur.Type == USER {
				p.advance()
			}
			roles := p.parseRoleList()
			optItems = append(optItems, &nodes.DefElem{Defname: "roles", Arg: roles, Location: -1})
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
		action = p.parseDefACLGrantAction()
	} else if p.cur.Type == REVOKE {
		p.advance()
		action = p.parseDefACLRevokeAction()
	}
	return &nodes.AlterDefaultPrivilegesStmt{Options: options, Action: action}
}

func (p *Parser) parseDefACLGrantAction() *nodes.GrantStmt {
	privs := p.parsePrivileges()
	p.expect(ON)
	objtype := p.parseDefACLPrivilegeTarget()
	p.expect(TO)
	grantees := p.parseGranteeList()
	grantOption := false
	if p.cur.Type == WITH {
		p.advance()
		p.expect(GRANT)
		p.expect(OPTION)
		grantOption = true
	}
	grantor := p.parseOptGrantedBy()
	return &nodes.GrantStmt{
		IsGrant: true, Privileges: privs, Targtype: nodes.ACL_TARGET_DEFAULTS,
		Objtype: nodes.ObjectType(objtype), Grantees: grantees,
		GrantOption: grantOption, Grantor: grantor,
	}
}

func (p *Parser) parseDefACLRevokeAction() *nodes.GrantStmt {
	grantOption := false
	if p.cur.Type == GRANT {
		p.advance()
		p.expect(OPTION)
		p.expect(FOR)
		grantOption = true
	}
	privs := p.parsePrivileges()
	p.expect(ON)
	objtype := p.parseDefACLPrivilegeTarget()
	p.expect(FROM)
	grantees := p.parseGranteeList()
	grantor := p.parseOptGrantedBy()
	behavior := p.parseOptDropBehavior()
	return &nodes.GrantStmt{
		IsGrant: false, GrantOption: grantOption, Privileges: privs,
		Targtype: nodes.ACL_TARGET_DEFAULTS, Objtype: nodes.ObjectType(objtype),
		Grantees: grantees, Grantor: grantor, Behavior: nodes.DropBehavior(behavior),
	}
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
func (p *Parser) parseAlterStatisticsStmt() nodes.Node {
	p.advance() // consume STATISTICS
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}
	defnames, _ := p.parseAnyName()
	p.expect(SET)
	p.expect(STATISTICS)
	stmt := &nodes.AlterStatsStmt{Defnames: defnames, MissingOk: missingOk}
	if p.cur.Type == DEFAULT {
		p.advance()
	} else {
		val := p.parseSignedIconst()
		stmt.Stxstattarget = int(val)
	}
	return stmt
}

// parseAlterOperatorStmt parses ALTER OPERATOR ...
// ALTER has already been consumed. Current token is OPERATOR.
func (p *Parser) parseAlterOperatorStmt() nodes.Node {
	p.advance() // consume OPERATOR
	if p.cur.Type == CLASS || p.cur.Type == FAMILY {
		// ALTER OPERATOR CLASS/FAMILY ... (SET SCHEMA, OWNER, RENAME, or ADD/DROP)
		isClass := p.cur.Type == CLASS
		p.advance()
		names, _ := p.parseAnyName()
		p.expect(USING)
		amName, _ := p.parseName()
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
				newschema, _ := p.parseName()
				return &nodes.AlterObjectSchemaStmt{ObjectType: objtype, Object: object, Newschema: newschema}
			}
		}
		if p.cur.Type == OWNER {
			p.advance()
			p.expect(TO)
			roleSpec := p.parseRoleSpec()
			return &nodes.AlterOwnerStmt{ObjectType: objtype, Object: object, Newowner: roleSpec}
		}
		if p.cur.Type == RENAME {
			p.advance()
			p.expect(TO)
			newname, _ := p.parseName()
			return &nodes.RenameStmt{RenameType: objtype, Object: object, Newname: newname}
		}
		return nil
	}
	// Plain ALTER OPERATOR op_with_argtypes ...
	opname, _ := p.parseAnyOperator()
	operArgs := p.parseOperArgtypes()
	owa := &nodes.ObjectWithArgs{Objname: opname, Objargs: operArgs}

	switch p.cur.Type {
	case SET:
		next := p.peekNext()
		if next.Type == SCHEMA {
			p.advance()
			p.advance()
			newschema, _ := p.parseName()
			return &nodes.AlterObjectSchemaStmt{
				ObjectType: nodes.OBJECT_OPERATOR,
				Object:     owa,
				Newschema:  newschema,
			}
		}
		// ALTER OPERATOR op SET (operator_def_list)
		p.advance() // consume SET
		p.expect('(')
		opts := p.parseOperatorDefList()
		p.expect(')')
		return &nodes.AlterOperatorStmt{
			Opername: owa,
			Options:  opts,
		}
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_OPERATOR,
			Object:     owa,
			Newowner:   roleSpec,
		}
	default:
		return nil
	}
}

// parseOperatorDefList parses: operator_def_elem (',' operator_def_elem)*
func (p *Parser) parseOperatorDefList() *nodes.List {
	elem := p.parseOperatorDefElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseOperatorDefElem())
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseOperatorDefElem() *nodes.DefElem {
	label, _ := p.parseColLabel()
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type == NONE {
			p.advance()
			return &nodes.DefElem{Defname: label, Location: -1}
		}
		arg := p.parseDefArg()
		return &nodes.DefElem{Defname: label, Arg: arg, Location: -1}
	}
	return &nodes.DefElem{Defname: label, Location: -1}
}
