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

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	if p.cur.Type == WITH {
		p.advance()
	}

	opts, err := p.parseCreateExtensionOptList()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateExtensionStmt{
		Extname:     name,
		IfNotExists: ifNotExists,
		Options:     opts,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

func (p *Parser) parseCreateExtensionOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseCreateExtensionOptItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseCreateExtensionOptItem parses create_extension_opt_item, or returns
// nil when none starts here.
//
//	create_extension_opt_item: SCHEMA name | VERSION_P NonReservedWord_or_Sconst | FROM NonReservedWord_or_Sconst | CASCADE
func (p *Parser) parseCreateExtensionOptItem() (*nodes.DefElem, error) {
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return makeDefElem("schema", &nodes.String{Str: name}), nil
	case VERSION_P:
		p.advance()
		ver, err := p.parseNonReservedWordOrSconstStrict()
		if err != nil {
			return nil, err
		}
		return makeDefElem("new_version", &nodes.String{Str: ver}), nil
	case CASCADE:
		p.advance()
		return makeDefElem("cascade", &nodes.Boolean{Boolval: true}), nil
	case FROM:
		// FROM NonReservedWord_or_Sconst: the grammar reads it and its
		// action refuses it at FROM.
		from := p.cur.Loc
		p.advance()
		if _, err := p.parseNonReservedWordOrSconstStrict(); err != nil {
			return nil, err
		}
		if p.collectMode() {
			return nil, errCollecting
		}
		return nil, &ParseError{Code: "0A000", Message: "CREATE EXTENSION ... FROM is no longer supported", Position: from}
	// optional-probe: create_extension_opt_list ends at any other token
	default:
		return nil, nil
	}
}

// parseAlterExtensionStmt parses ALTER EXTENSION statements.
func (p *Parser) parseAlterExtensionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance()       // consume EXTENSION

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case UPDATE:
		p.advance()
		opts, err := p.parseAlterExtensionOptList()
		if err != nil {
			return nil, err
		}
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
		newschema, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterObjectSchemaStmt{
			ObjectType: nodes.OBJECT_EXTENSION,
			Object:     &nodes.String{Str: name},
			Newschema:  newschema,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	// exhaustive: gram.y:5166 — AlterExtensionStmt enumerates UPDATE / ADD /
	// DROP / SET SCHEMA only. Unknown action after `ALTER EXTENSION name`
	// is a syntax error.
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

func (p *Parser) parseAlterExtensionOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseAlterExtensionOptItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseAlterExtensionOptItem parses alter_extension_opt_item: TO
// NonReservedWord_or_Sconst, or returns nil when TO does not follow.
func (p *Parser) parseAlterExtensionOptItem() (*nodes.DefElem, error) {
	if p.cur.Type != TO {
		return nil, nil
	}
	p.advance()
	ver, err := p.parseNonReservedWordOrSconstStrict()
	if err != nil {
		return nil, err
	}
	return makeDefElem("new_version", &nodes.String{Str: ver}), nil
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

// parseCreateAmStmt parses CREATE ACCESS METHOD statement.
func (p *Parser) parseCreateAmStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume ACCESS
	if _, err := p.expect(METHOD); err != nil {
		return nil, err
	}
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(TYPE_P); err != nil {
		return nil, err
	}
	amtype, err := p.parseAmType()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(HANDLER); err != nil {
		return nil, err
	}
	handlerName, err := p.parseHandlerName()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateAmStmt{
		Amname: name, HandlerName: handlerName, Amtype: amtype,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseAmType parses am_type: INDEX | TABLE.
func (p *Parser) parseAmType() (byte, error) {
	if p.collectMode() {
		p.addTokenCandidate(INDEX)
		p.addTokenCandidate(TABLE)
		return 0, errCollecting
	}
	switch p.cur.Type {
	case INDEX:
		p.advance()
		return nodes.AMTYPE_INDEX, nil
	case TABLE:
		p.advance()
		return nodes.AMTYPE_TABLE, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// parseCreateCastStmt parses CREATE CAST statement.
func (p *Parser) parseCreateCastStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume CAST
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	srcType, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tgtType, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	stmt := &nodes.CreateCastStmt{Sourcetype: srcType, Targettype: tgtType}

	if p.cur.Type == WITHOUT {
		p.advance()
		if _, err := p.expect(FUNCTION); err != nil {
			return nil, err
		}
		if stmt.Context, err = p.parseCastContext(); err != nil {
			return nil, err
		}
		stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		return stmt, nil
	}
	if _, err := p.expect(WITH); err != nil {
		return nil, err
	}
	if p.cur.Type == INOUT {
		p.advance()
		if stmt.Context, err = p.parseCastContext(); err != nil {
			return nil, err
		}
		stmt.Inout = true
		stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
		return stmt, nil
	}
	if _, err := p.expect(FUNCTION); err != nil {
		return nil, err
	}
	fn, err := p.parseFunctionWithArgtypes()
	if err != nil {
		return nil, err
	}
	stmt.Func = fn
	if stmt.Context, err = p.parseCastContext(); err != nil {
		return nil, err
	}
	stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End}
	return stmt, nil
}

// parseCastContext parses cast_context, which ends CREATE CAST: AS starts
// nothing else there.
//
//	cast_context: AS IMPLICIT_P | AS ASSIGNMENT | /* EMPTY */
func (p *Parser) parseCastContext() (nodes.CoercionContext, error) {
	if p.cur.Type != AS {
		return nodes.COERCION_EXPLICIT, nil
	}
	p.advance() // AS
	switch p.cur.Type {
	case IMPLICIT_P:
		p.advance()
		return nodes.COERCION_IMPLICIT, nil
	case ASSIGNMENT:
		p.advance()
		return nodes.COERCION_ASSIGNMENT, nil
	}
	if p.collectMode() {
		p.addTokenCandidate(IMPLICIT_P)
		p.addTokenCandidate(ASSIGNMENT)
		return nodes.COERCION_EXPLICIT, errCollecting
	}
	return nodes.COERCION_EXPLICIT, p.syntaxErrorAtCur()
}

// parseDropCastStmt parses DROP CAST statement.
func (p *Parser) parseDropCastStmt() (nodes.Node, error) {
	p.advance() // consume CAST
	missingOk, err := p.parseIfExistsBeforeKeyword()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	srcType, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	tgtType, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
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
	typeName, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(LANGUAGE); err != nil {
		return nil, err
	}
	lang, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	fromsql, tosql, err := p.parseTransformElementList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.CreateTransformStmt{
		Replace: replace, TypeName: typeName, Lang: lang,
		Fromsql: fromsql, Tosql: tosql,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// parseTransformElementList parses transform_element_list: a FROM SQL
// element, a TO SQL element, or both in either order.
//
//	transform_element_list: FROM SQL_P WITH FUNCTION function_with_argtypes ',' TO SQL_P WITH FUNCTION function_with_argtypes
//	    | TO SQL_P WITH FUNCTION function_with_argtypes ',' FROM SQL_P WITH FUNCTION function_with_argtypes
//	    | FROM SQL_P WITH FUNCTION function_with_argtypes
//	    | TO SQL_P WITH FUNCTION function_with_argtypes
func (p *Parser) parseTransformElementList() (fromsql, tosql *nodes.ObjectWithArgs, err error) {
	if p.collectMode() {
		p.addTokenCandidate(FROM)
		p.addTokenCandidate(TO)
		return nil, nil, errCollecting
	}
	first := p.cur.Type
	if first != FROM && first != TO {
		return nil, nil, p.syntaxErrorAtCur()
	}
	fn, err := p.parseTransformElement(first)
	if err != nil {
		return nil, nil, err
	}
	if first == FROM {
		fromsql = fn
	} else {
		tosql = fn
	}
	if p.cur.Type != ',' {
		return fromsql, tosql, nil
	}
	p.advance()
	second := TO
	if first == TO {
		second = FROM
	}
	fn, err = p.parseTransformElement(second)
	if err != nil {
		return nil, nil, err
	}
	if second == FROM {
		fromsql = fn
	} else {
		tosql = fn
	}
	return fromsql, tosql, nil
}

// parseTransformElement parses one element of transform_element_list:
// lead SQL_P WITH FUNCTION function_with_argtypes, lead being FROM or TO.
func (p *Parser) parseTransformElement(lead int) (*nodes.ObjectWithArgs, error) {
	for _, tok := range []int{lead, SQL_P, WITH, FUNCTION} {
		if _, err := p.expect(tok); err != nil {
			return nil, err
		}
	}
	return p.parseFunctionWithArgtypes()
}

// parseDropTransformStmt parses DROP TRANSFORM statement.
func (p *Parser) parseDropTransformStmt() (nodes.Node, error) {
	p.advance() // consume TRANSFORM
	missingOk, err := p.parseIfExistsBeforeKeyword()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	typeName, err := p.parseTypename()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(LANGUAGE); err != nil {
		return nil, err
	}
	lang, err := p.parseName()
	if err != nil {
		return nil, err
	}
	behavior := p.parseOptDropBehavior()
	return &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_TRANSFORM),
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: []nodes.Node{typeName, &nodes.String{Str: lang}}}}},
		Behavior:   behavior, Missing_ok: missingOk,
	}, nil
}
