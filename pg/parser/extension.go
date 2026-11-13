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
	// exhaustive: gram.y:5166 — AlterExtensionStmt enumerates UPDATE / ADD /
	// DROP / SET SCHEMA only. Unknown action after `ALTER EXTENSION name`
	// is a syntax error.
	default:
		return nil, p.syntaxErrorAtCur()
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

// parseAlterExtensionContents parses what ALTER EXTENSION name ADD|DROP adds
// or drops.
//
//	AlterExtensionContentsStmt: ALTER EXTENSION name add_drop
//	    object_type_name name | object_type_any_name any_name
//	    | AGGREGATE aggregate_with_argtypes | CAST '(' Typename AS Typename ')'
//	    | DOMAIN_P Typename | FUNCTION function_with_argtypes
//	    | OPERATOR operator_with_argtypes | OPERATOR CLASS any_name USING name
//	    | OPERATOR FAMILY any_name USING name | PROCEDURE function_with_argtypes
//	    | ROUTINE function_with_argtypes | TRANSFORM FOR Typename LANGUAGE name
//	    | TYPE_P Typename
func (p *Parser) parseAlterExtensionContents(loc int, extname string, action int) (nodes.Node, error) {
	objtype, obj, err := p.parseAlterExtensionObject()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterExtensionContentsStmt{
		Extname: extname, Action: action,
		Objtype: objtype, Object: obj,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseAlterExtensionObject() (nodes.ObjectType, nodes.Node, error) {
	switch p.cur.Type {
	case AGGREGATE:
		p.advance()
		// aggregate_with_argtypes takes aggr_args rather than func_args, so "agg()" is rejected.
		obj, err := p.parseAggregateWithArgtypesLocal()
		return nodes.OBJECT_AGGREGATE, obj, err
	case CAST:
		p.advance()
		if _, err := p.expect('('); err != nil {
			return 0, nil, err
		}
		source, err := p.parseTypename()
		if err != nil {
			return 0, nil, err
		}
		if _, err := p.expect(AS); err != nil {
			return 0, nil, err
		}
		target, err := p.parseTypename()
		if err != nil {
			return 0, nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return 0, nil, err
		}
		return nodes.OBJECT_CAST, &nodes.List{Items: []nodes.Node{source, target}}, nil
	case DOMAIN_P, TYPE_P:
		objtype := nodes.OBJECT_TYPE
		if p.cur.Type == DOMAIN_P {
			objtype = nodes.OBJECT_DOMAIN
		}
		p.advance()
		tn, err := p.parseTypename()
		return objtype, tn, err
	case FUNCTION, PROCEDURE, ROUTINE:
		objtype := map[int]nodes.ObjectType{FUNCTION: nodes.OBJECT_FUNCTION, PROCEDURE: nodes.OBJECT_PROCEDURE, ROUTINE: nodes.OBJECT_ROUTINE}[p.cur.Type]
		p.advance()
		obj, err := p.parseFunctionWithArgtypes()
		return objtype, obj, err
	case OPERATOR:
		p.advance()
		if p.cur.Type == CLASS || p.cur.Type == FAMILY {
			objtype := nodes.OBJECT_OPCLASS
			if p.cur.Type == FAMILY {
				objtype = nodes.OBJECT_OPFAMILY
			}
			p.advance()
			anyName, err := p.parseAnyName()
			if err != nil {
				return 0, nil, err
			}
			if _, err := p.expect(USING); err != nil {
				return 0, nil, err
			}
			usingName, err := p.parseName()
			if err != nil {
				return 0, nil, err
			}
			return objtype, &nodes.List{Items: append([]nodes.Node{&nodes.String{Str: usingName}}, anyName.Items...)}, nil
		}
		opname, err := p.parseAnyOperator()
		if err != nil {
			return 0, nil, err
		}
		args, err := p.parseOperArgtypes()
		if err != nil {
			return 0, nil, err
		}
		return nodes.OBJECT_OPERATOR, &nodes.ObjectWithArgs{Objname: opname, Objargs: args}, nil
	case TRANSFORM:
		p.advance()
		if _, err := p.expect(FOR); err != nil {
			return 0, nil, err
		}
		tn, err := p.parseTypename()
		if err != nil {
			return 0, nil, err
		}
		if _, err := p.expect(LANGUAGE); err != nil {
			return 0, nil, err
		}
		lang, err := p.parseName()
		if err != nil {
			return 0, nil, err
		}
		return nodes.OBJECT_TRANSFORM, &nodes.List{Items: []nodes.Node{tn, &nodes.String{Str: lang}}}, nil
	}
	if objtype, ok, err := p.parseExtObjectTypeName(); ok || err != nil {
		if err != nil {
			return 0, nil, err
		}
		name, err := p.parseName()
		if err != nil {
			return 0, nil, err
		}
		return objtype, &nodes.String{Str: name}, nil
	}
	objtype, err := p.parseObjectTypeAnyName()
	if err != nil {
		return 0, nil, err
	}
	anyName, err := p.parseAnyName()
	if err != nil {
		return 0, nil, err
	}
	return objtype, anyName, nil
}

// parseExtObjectTypeName parses object_type_name, the object types
// ALTER EXTENSION names by a plain name, and reports whether one was there.
//
//	object_type_name: drop_type_name | DATABASE | ROLE | SUBSCRIPTION | TABLESPACE
//	drop_type_name: ACCESS METHOD | EVENT TRIGGER | EXTENSION | FOREIGN DATA_P WRAPPER
//	    | opt_procedural LANGUAGE | PUBLICATION | SCHEMA | SERVER
func (p *Parser) parseExtObjectTypeName() (nodes.ObjectType, bool, error) {
	two := func(tok int, t nodes.ObjectType) (nodes.ObjectType, bool, error) {
		p.advance()
		if _, err := p.expect(tok); err != nil {
			return 0, true, err
		}
		return t, true, nil
	}
	switch p.cur.Type {
	case ACCESS:
		return two(METHOD, nodes.OBJECT_ACCESS_METHOD)
	case EVENT:
		return two(TRIGGER, nodes.OBJECT_EVENT_TRIGGER)
	case PROCEDURAL:
		return two(LANGUAGE, nodes.OBJECT_LANGUAGE)
	case FOREIGN:
		// FOREIGN TABLE is object_type_any_name's.
		if p.peekNext().Type != DATA_P {
			return 0, false, nil
		}
		p.advance()
		return two(WRAPPER, nodes.OBJECT_FDW)
	}
	objtype, ok := map[int]nodes.ObjectType{
		EXTENSION: nodes.OBJECT_EXTENSION, LANGUAGE: nodes.OBJECT_LANGUAGE,
		PUBLICATION: nodes.OBJECT_PUBLICATION, SCHEMA: nodes.OBJECT_SCHEMA,
		SERVER: nodes.OBJECT_FOREIGN_SERVER, DATABASE: nodes.OBJECT_DATABASE,
		ROLE: nodes.OBJECT_ROLE, SUBSCRIPTION: nodes.OBJECT_SUBSCRIPTION,
		TABLESPACE: nodes.OBJECT_TABLESPACE,
	}[p.cur.Type]
	if !ok {
		return 0, false, nil
	}
	p.advance()
	return objtype, true, nil
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
