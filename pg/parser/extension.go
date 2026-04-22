package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateExtensionStmt parses a CREATE EXTENSION statement.
func (p *Parser) parseCreateExtensionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume EXTENSION

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

	name, _ := p.parseName()

	if p.cur.Type == WITH {
		p.advance()
	}

	opts := p.parseCreateExtensionOptList()

	return &nodes.CreateExtensionStmt{
		Extname:     name,
		IfNotExists: ifNotExists,
		Options:     opts,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseCreateExtensionOptList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseCreateExtensionOptItem()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseCreateExtensionOptItem() *nodes.DefElem {
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		name, _ := p.parseName()
		return makeDefElem("schema", &nodes.String{Str: name})
	case VERSION_P:
		p.advance()
		ver := p.parseNonReservedWordOrSconst()
		return makeDefElem("new_version", &nodes.String{Str: ver})
	case CASCADE:
		p.advance()
		return makeDefElem("cascade", &nodes.Boolean{Boolval: true})
	default:
		return nil
	}
}

// parseAlterExtensionStmt parses ALTER EXTENSION statements.
func (p *Parser) parseAlterExtensionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance()       // consume EXTENSION

	name, _ := p.parseName()

	switch p.cur.Type {
	case UPDATE:
		p.advance()
		opts := p.parseAlterExtensionOptList()
		return &nodes.AlterExtensionStmt{
			Extname: name,
			Options: opts,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case ADD_P:
		p.advance()
		return p.parseAlterExtensionContents(loc, name, 1)
	case DROP:
		p.advance()
		return p.parseAlterExtensionContents(loc, name, -1)
	case SET:
		p.advance()
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, err
		}
		newschema, _ := p.parseName()
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_EXTENSION,
			Object:     &nodes.String{Str: name},
			Newschema:  newschema,
		}, nil
	// known-gap: not a KB-2 blocker; tracked in PARSER_DISPATCH_AUDIT.md §2 for future fix
	default:
		return nil, nil
	}
}

func (p *Parser) parseAlterExtensionOptList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseAlterExtensionOptItem()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseAlterExtensionOptItem() *nodes.DefElem {
	if p.cur.Type == TO {
		p.advance()
		ver := p.parseNonReservedWordOrSconst()
		return makeDefElem("new_version", &nodes.String{Str: ver})
	}
	return nil
}

func (p *Parser) parseAlterExtensionContents(loc int, extname string, action int) (nodes.Node, error) {
	switch p.cur.Type {
	case AGGREGATE:
		p.advance()
		obj := p.parseExtFuncWithArgtypes()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_AGGREGATE, Object: obj,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case FUNCTION:
		p.advance()
		obj := p.parseExtFuncWithArgtypes()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_FUNCTION, Object: obj,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case PROCEDURE:
		p.advance()
		obj := p.parseExtFuncWithArgtypes()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_PROCEDURE, Object: obj,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case ROUTINE:
		p.advance()
		obj := p.parseExtFuncWithArgtypes()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_ROUTINE, Object: obj,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case OPERATOR:
		p.advance()
		if p.cur.Type == CLASS {
			p.advance()
			anyName, _ := p.parseAnyName()
			p.expect(USING)
			usingName, _ := p.parseName()
			obj := &nodes.List{Items: append([]nodes.Node{&nodes.String{Str: usingName}}, anyName.Items...)}
			return &nodes.AlterExtensionContentsStmt{
				Extname: extname, Action: action,
				Objtype: nodes.OBJECT_OPCLASS, Object: obj,
				Loc: nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if p.cur.Type == FAMILY {
			p.advance()
			anyName, _ := p.parseAnyName()
			p.expect(USING)
			usingName, _ := p.parseName()
			obj := &nodes.List{Items: append([]nodes.Node{&nodes.String{Str: usingName}}, anyName.Items...)}
			return &nodes.AlterExtensionContentsStmt{
				Extname: extname, Action: action,
				Objtype: nodes.OBJECT_OPFAMILY, Object: obj,
				Loc: nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		obj := p.parseExtOperWithArgtypes()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_OPERATOR, Object: obj,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case DOMAIN_P:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_DOMAIN, Object: tn,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case TYPE_P:
		p.advance()
		tn, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: nodes.OBJECT_TYPE, Object: tn,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		objType, ok := p.tryParseExtObjTypeName()
		if ok {
			objName, _ := p.parseName()
			return &nodes.AlterExtensionContentsStmt{
				Extname: extname, Action: action,
				Objtype: objType, Object: &nodes.String{Str: objName},
				Loc: nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		objType = p.parseObjectTypeAnyName()
		anyName, _ := p.parseAnyName()
		return &nodes.AlterExtensionContentsStmt{
			Extname: extname, Action: action,
			Objtype: objType, Object: anyName,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}
}

func (p *Parser) tryParseExtObjTypeName() (nodes.ObjectType, bool) {
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		return nodes.OBJECT_SCHEMA, true
	case DATABASE:
		p.advance()
		return nodes.OBJECT_DATABASE, true
	case ROLE:
		p.advance()
		return nodes.OBJECT_ROLE, true
	case TABLESPACE:
		p.advance()
		return nodes.OBJECT_TABLESPACE, true
	case SUBSCRIPTION:
		p.advance()
		return nodes.OBJECT_SUBSCRIPTION, true
	case PUBLICATION:
		p.advance()
		return nodes.OBJECT_PUBLICATION, true
	case SERVER:
		p.advance()
		return nodes.OBJECT_FOREIGN_SERVER, true
	default:
		return 0, false
	}
}

// parseExtFuncWithArgtypes parses function_with_argtypes for extension statements.
func (p *Parser) parseExtFuncWithArgtypes() *nodes.ObjectWithArgs {
	funcname, _ := p.parseFuncName()
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			p.advance()
			return &nodes.ObjectWithArgs{Objname: funcname, Objargs: nil}
		}
		args := p.parseExtFuncArgsList()
		p.expect(')')
		argtypes := extExtractArgTypes(args)
		return &nodes.ObjectWithArgs{Objname: funcname, Objargs: argtypes}
	}
	return &nodes.ObjectWithArgs{Objname: funcname, ArgsUnspecified: true}
}

func (p *Parser) parseExtFuncArgsList() *nodes.List {
	arg := p.parseFuncArg()
	items := []nodes.Node{arg}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseFuncArg())
	}
	return &nodes.List{Items: items}
}

func extExtractArgTypes(args *nodes.List) *nodes.List {
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

// parseExtOperWithArgtypes parses operator_with_argtypes for extension statements.
func (p *Parser) parseExtOperWithArgtypes() *nodes.ObjectWithArgs {
	opname, _ := p.parseAnyOperator()
	operArgs := p.parseExtOperArgtypes()
	return &nodes.ObjectWithArgs{Objname: opname, Objargs: operArgs}
}

func (p *Parser) parseExtOperArgtypes() *nodes.List {
	p.expect('(')
	if p.cur.Type == NONE {
		p.advance()
		p.expect(',')
		typ, _ := p.parseTypename()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{nil, typ}}
	}
	typ1, _ := p.parseTypename()
	if p.cur.Type == ')' {
		p.advance()
		return &nodes.List{Items: []nodes.Node{typ1}}
	}
	p.expect(',')
	if p.cur.Type == NONE {
		p.advance()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{typ1, nil}}
	}
	typ2, _ := p.parseTypename()
	p.expect(')')
	return &nodes.List{Items: []nodes.Node{typ1, typ2}}
}

// parseExtAggrArgs parses aggr_args for extension statements.
func (p *Parser) parseExtAggrArgs() *nodes.List {
	p.expect('(')
	if p.cur.Type == '*' {
		p.advance()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{nil, &nodes.Integer{Ival: -1}}}
	}
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		args := p.parseExtFuncArgsList()
		p.expect(')')
		return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: 0}}}
	}
	args := p.parseExtFuncArgsList()
	if p.cur.Type == ORDER {
		p.advance()
		p.expect(BY)
		orderedArgs := p.parseExtFuncArgsList()
		p.expect(')')
		directCount := len(args.Items)
		merged := &nodes.List{Items: append(args.Items, orderedArgs.Items...)}
		return &nodes.List{Items: []nodes.Node{merged, &nodes.Integer{Ival: int64(directCount)}}}
	}
	p.expect(')')
	return &nodes.List{Items: []nodes.Node{args, &nodes.Integer{Ival: -1}}}
}

func extExtractAggrArgTypes(args *nodes.List) *nodes.List {
	if args == nil || len(args.Items) < 1 {
		return nil
	}
	argsList, ok := args.Items[0].(*nodes.List)
	if !ok || argsList == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range argsList.Items {
		if fp, ok := item.(*nodes.FunctionParameter); ok {
			result.Items = append(result.Items, fp.ArgType)
		}
	}
	return result
}

// parseCreateAmStmt parses CREATE ACCESS METHOD statement.
func (p *Parser) parseCreateAmStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume ACCESS
	if _, err := p.expect(METHOD); err != nil {
		return nil, err
	}
	name, _ := p.parseName()
	if _, err := p.expect(TYPE_P); err != nil {
		return nil, err
	}
	amtype := p.parseAmType()
	if _, err := p.expect(HANDLER); err != nil {
		return nil, err
	}
	handlerName := p.parseExtHandlerName()
	return &nodes.CreateAmStmt{
		Amname: name, HandlerName: handlerName, Amtype: amtype,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseAmType() byte {
	switch p.cur.Type {
	case INDEX:
		p.advance()
		return nodes.AMTYPE_INDEX
	case TABLE:
		p.advance()
		return nodes.AMTYPE_TABLE
	default:
		p.advance()
		return nodes.AMTYPE_INDEX
	}
}

func (p *Parser) parseExtHandlerName() *nodes.List {
	name, _ := p.parseName()
	items := []nodes.Node{&nodes.String{Str: name}}
	if p.cur.Type == '.' {
		attrs, err := p.parseAttrs()
		if err == nil && attrs != nil {
			items = append(items, attrs.Items...)
		}
	}
	return &nodes.List{Items: items}
}

// parseCreateCastStmt parses CREATE CAST statement.
func (p *Parser) parseCreateCastStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	srcType, _ := p.parseTypename()
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tgtType, _ := p.parseTypename()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	stmt := &nodes.CreateCastStmt{Sourcetype: srcType, Targettype: tgtType}

	if p.cur.Type == WITHOUT {
		p.advance()
		if _, err := p.expect(FUNCTION); err != nil {
			return nil, err
		}
		stmt.Context = p.parseCastContext()
		stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		return stmt, nil
	}
	if _, err := p.expect(WITH); err != nil {
		return nil, err
	}
	if p.cur.Type == INOUT {
		p.advance()
		stmt.Context = p.parseCastContext()
		stmt.Inout = true
		stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		return stmt, nil
	}
	if _, err := p.expect(FUNCTION); err != nil {
		return nil, err
	}
	stmt.Func = p.parseExtFuncWithArgtypes()
	stmt.Context = p.parseCastContext()
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

func (p *Parser) parseCastContext() nodes.CoercionContext {
	if p.cur.Type == AS {
		next := p.peekNext()
		if next.Type == IMPLICIT_P {
			p.advance()
			p.advance()
			return nodes.COERCION_IMPLICIT
		}
		if next.Type == ASSIGNMENT {
			p.advance()
			p.advance()
			return nodes.COERCION_ASSIGNMENT
		}
	}
	return nodes.COERCION_EXPLICIT
}

// parseDropCastStmt parses DROP CAST statement.
func (p *Parser) parseDropCastStmt() (nodes.Node, error) {
	p.advance() // consume CAST
	missingOk := p.parseOptIfExists()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	srcType, _ := p.parseTypename()
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tgtType, _ := p.parseTypename()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	return &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_CAST),
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{srcType, tgtType}}}},
		Behavior:   behavior, Missing_ok: missingOk,
	}, nil
}

// parseCreateTransformStmt parses CREATE [OR REPLACE] TRANSFORM statement.
func (p *Parser) parseCreateTransformStmt(replace bool, loc int) (nodes.Node, error) {
	p.advance() // consume TRANSFORM
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	typeName, _ := p.parseTypename()
	if _, err := p.expect(LANGUAGE); err != nil {
		return nil, err
	}
	lang, _ := p.parseName()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	fromsql, tosql := p.parseTransformElementList()
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateTransformStmt{
		Replace: replace, TypeName: typeName, Lang: lang,
		Fromsql: fromsql, Tosql: tosql,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseTransformElementList() (fromsql *nodes.ObjectWithArgs, tosql *nodes.ObjectWithArgs) {
	if p.cur.Type == FROM {
		p.advance()
		p.expect(SQL_P)
		p.expect(WITH)
		p.expect(FUNCTION)
		fromsql = p.parseExtFuncWithArgtypes()
		if p.cur.Type == ',' {
			p.advance()
			p.expect(TO)
			p.expect(SQL_P)
			p.expect(WITH)
			p.expect(FUNCTION)
			tosql = p.parseExtFuncWithArgtypes()
		}
		return
	}
	p.expect(TO)
	p.expect(SQL_P)
	p.expect(WITH)
	p.expect(FUNCTION)
	tosql = p.parseExtFuncWithArgtypes()
	if p.cur.Type == ',' {
		p.advance()
		p.expect(FROM)
		p.expect(SQL_P)
		p.expect(WITH)
		p.expect(FUNCTION)
		fromsql = p.parseExtFuncWithArgtypes()
	}
	return
}

// parseDropTransformStmt parses DROP TRANSFORM statement.
func (p *Parser) parseDropTransformStmt() (nodes.Node, error) {
	p.advance() // consume TRANSFORM
	missingOk := p.parseOptIfExists()
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	typeName, _ := p.parseTypename()
	if _, err := p.expect(LANGUAGE); err != nil {
		return nil, err
	}
	lang, _ := p.parseName()
	behavior := p.parseOptDropBehavior()
	return &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_TRANSFORM),
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{typeName, &nodes.String{Str: lang}}}}},
		Behavior:   behavior, Missing_ok: missingOk,
	}, nil
}
