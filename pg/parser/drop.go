package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseDropStmt parses a DROP statement. The DROP keyword has already been consumed by parseStmt.
//
//	DropStmt:
//	    DROP object_type_any_name IF_P EXISTS any_name_list opt_drop_behavior
//	    | DROP object_type_any_name any_name_list opt_drop_behavior
//	    | DROP TYPE_P [IF EXISTS] any_name_list opt_drop_behavior
//	    | DROP DOMAIN_P [IF EXISTS] any_name_list opt_drop_behavior
//	    | DROP SCHEMA|EXTENSION [IF EXISTS] name_list opt_drop_behavior
//	    | DROP FUNCTION|PROCEDURE|ROUTINE|AGGREGATE [IF EXISTS] any_name_list opt_drop_behavior
//	    | DROP TRIGGER|POLICY|RULE [IF EXISTS] name ON any_name opt_drop_behavior
//	    | DROP EVENT TRIGGER [IF EXISTS] name opt_drop_behavior
//	    | DROP LANGUAGE [IF EXISTS] name opt_drop_behavior
//	    | DROP PUBLICATION|SUBSCRIPTION [IF EXISTS] name_list opt_drop_behavior
//	    | DROP INDEX CONCURRENTLY [IF EXISTS] any_name_list opt_drop_behavior
//	    | DROP FOREIGN DATA WRAPPER [IF EXISTS] name_list opt_drop_behavior
//	    | DROP SERVER [IF EXISTS] name_list opt_drop_behavior
//	    | DROP OWNED BY name_list opt_drop_behavior
func (p *Parser) parseDropStmt(stmtLoc int) (nodes.Node, error) {
	switch p.cur.Type {
	case TABLE, SEQUENCE, VIEW, INDEX, COLLATION, CONVERSION_P, STATISTICS:
		return p.parseDropObjectTypeAnyName()
	case MATERIALIZED:
		// MATERIALIZED VIEW
		return p.parseDropObjectTypeAnyName()
	case FOREIGN:
		// Could be FOREIGN TABLE or FOREIGN DATA WRAPPER
		next := p.peekNext()
		if next.Type == TABLE {
			return p.parseDropObjectTypeAnyName()
		}
		// FOREIGN DATA WRAPPER
		return p.parseDropForeignDataWrapper()
	case TEXT_P:
		// TEXT SEARCH PARSER|DICTIONARY|TEMPLATE|CONFIGURATION
		return p.parseDropObjectTypeAnyName()
	case ACCESS:
		// ACCESS METHOD
		return p.parseDropObjectTypeAnyName()
	case TYPE_P:
		return p.parseDropTypeOrDomain(nodes.OBJECT_TYPE)
	case DOMAIN_P:
		return p.parseDropTypeOrDomain(nodes.OBJECT_DOMAIN)
	case SCHEMA:
		return p.parseDropNameList(nodes.OBJECT_SCHEMA)
	case EXTENSION:
		return p.parseDropNameList(nodes.OBJECT_EXTENSION)
	case FUNCTION:
		return p.parseDropFuncStmt(nodes.OBJECT_FUNCTION)
	case PROCEDURE:
		return p.parseDropFuncStmt(nodes.OBJECT_PROCEDURE)
	case ROUTINE:
		return p.parseDropFuncStmt(nodes.OBJECT_ROUTINE)
	case AGGREGATE:
		return p.parseDropAggrStmt()
	case TRIGGER:
		return p.parseDropOnObject(nodes.OBJECT_TRIGGER)
	case POLICY:
		return p.parseDropOnObject(nodes.OBJECT_POLICY)
	case RULE:
		return p.parseDropOnObject(nodes.OBJECT_RULE)
	case EVENT:
		return p.parseDropEventTrigger()
	case LANGUAGE:
		return p.parseDropLanguage()
	case PUBLICATION:
		return p.parseDropNameList(nodes.OBJECT_PUBLICATION)
	case SUBSCRIPTION:
		return p.parseDropSubscription()
	case SERVER:
		return p.parseDropNameList(nodes.OBJECT_FOREIGN_SERVER)
	case OWNED:
		return p.parseDropOwned()
	case DATABASE:
		return p.parseDropdbStmt(stmtLoc)
	case ROLE, GROUP_P:
		return p.parseDropRoleStmt()
	case USER:
		// DROP USER MAPPING or DROP USER (role)
		next := p.peekNext()
		if next.Type == MAPPING {
			return p.parseDropUserMappingStmt()
		}
		return p.parseDropRoleStmt()
	case CAST:
		return p.parseDropCastStmt()
	case TRANSFORM:
		return p.parseDropTransformStmt()
	case OPERATOR:
		return p.parseDropOperatorClassOrFamily()
	case TABLESPACE:
		return p.parseDropTableSpaceStmt()
	default:
		return nil, nil
	}
}

// parseDropObjectTypeAnyName parses DROP of object_type_any_name objects.
//
//	object_type_any_name:
//	    TABLE | SEQUENCE | VIEW | MATERIALIZED VIEW | INDEX | FOREIGN TABLE
//	    | COLLATION | CONVERSION_P | STATISTICS
//	    | TEXT_P SEARCH PARSER | TEXT_P SEARCH DICTIONARY
//	    | TEXT_P SEARCH TEMPLATE | TEXT_P SEARCH CONFIGURATION
//	    | ACCESS METHOD
func (p *Parser) parseDropObjectTypeAnyName() (nodes.Node, error) {
	objType := p.parseObjectTypeAnyName()

	// Special case: DROP INDEX CONCURRENTLY
	concurrent := false
	if objType == nodes.OBJECT_INDEX && p.cur.Type == CONCURRENTLY {
		p.advance()
		concurrent = true
	}

	missingOk := p.parseOptIfExists()

	objects, err := p.parseAnyNameList()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    objects,
		RemoveType: int(objType),
		Behavior:   behavior,
		Missing_ok: missingOk,
		Concurrent: concurrent,
	}, nil
}

// parseObjectTypeAnyName parses the object_type_any_name rule and returns the ObjectType.
// It consumes the relevant tokens.
//
//	object_type_any_name:
//	    TABLE | SEQUENCE | VIEW | MATERIALIZED VIEW | INDEX | FOREIGN TABLE
//	    | COLLATION | CONVERSION_P | STATISTICS
//	    | TEXT_P SEARCH PARSER | TEXT_P SEARCH DICTIONARY
//	    | TEXT_P SEARCH TEMPLATE | TEXT_P SEARCH CONFIGURATION
//	    | ACCESS METHOD
func (p *Parser) parseObjectTypeAnyName() nodes.ObjectType {
	switch p.cur.Type {
	case TABLE:
		p.advance()
		return nodes.OBJECT_TABLE
	case SEQUENCE:
		p.advance()
		return nodes.OBJECT_SEQUENCE
	case VIEW:
		p.advance()
		return nodes.OBJECT_VIEW
	case MATERIALIZED:
		p.advance() // consume MATERIALIZED
		p.expect(VIEW)
		return nodes.OBJECT_MATVIEW
	case INDEX:
		p.advance()
		return nodes.OBJECT_INDEX
	case FOREIGN:
		p.advance() // consume FOREIGN
		p.expect(TABLE)
		return nodes.OBJECT_FOREIGN_TABLE
	case COLLATION:
		p.advance()
		return nodes.OBJECT_COLLATION
	case CONVERSION_P:
		p.advance()
		return nodes.OBJECT_CONVERSION
	case STATISTICS:
		p.advance()
		return nodes.OBJECT_STATISTIC_EXT
	case TEXT_P:
		p.advance() // consume TEXT
		p.expect(SEARCH)
		switch p.cur.Type {
		case PARSER:
			p.advance()
			return nodes.OBJECT_TSPARSER
		case DICTIONARY:
			p.advance()
			return nodes.OBJECT_TSDICTIONARY
		case TEMPLATE:
			p.advance()
			return nodes.OBJECT_TSTEMPLATE
		case CONFIGURATION:
			p.advance()
			return nodes.OBJECT_TSCONFIGURATION
		}
		return nodes.OBJECT_TSPARSER // fallback
	case ACCESS:
		p.advance() // consume ACCESS
		p.expect(METHOD)
		return nodes.OBJECT_ACCESS_METHOD
	}
	return nodes.OBJECT_TABLE // fallback
}

// parseDropTypeOrDomain parses DROP TYPE or DROP DOMAIN statements.
//
//	DROP TYPE_P [IF EXISTS] any_name_list opt_drop_behavior
//	DROP DOMAIN_P [IF EXISTS] any_name_list opt_drop_behavior
func (p *Parser) parseDropTypeOrDomain(objType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume TYPE or DOMAIN

	missingOk := p.parseOptIfExists()

	objects, err := p.parseAnyNameList()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    objects,
		RemoveType: int(objType),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropNameList parses DROP statements that use name_list (simple names, not qualified).
// Used for SCHEMA, EXTENSION, PUBLICATION, SUBSCRIPTION, SERVER.
//
//	DROP SCHEMA|EXTENSION|PUBLICATION|SUBSCRIPTION|SERVER [IF EXISTS] name_list opt_drop_behavior
func (p *Parser) parseDropNameList(objType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume the keyword

	missingOk := p.parseOptIfExists()

	nameList, err := p.parseNameList()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    makeNameListAsAnyNameList(nameList),
		RemoveType: int(objType),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropAnyNameList parses DROP statements that use any_name_list.
// Used for FUNCTION, PROCEDURE, ROUTINE, AGGREGATE.
//
//	DROP FUNCTION|PROCEDURE|ROUTINE|AGGREGATE [IF EXISTS] any_name_list opt_drop_behavior
func (p *Parser) parseDropAnyNameList(objType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume the keyword

	missingOk := p.parseOptIfExists()

	objects, err := p.parseAnyNameList()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    objects,
		RemoveType: int(objType),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropOnObject parses DROP TRIGGER|POLICY|RULE name ON any_name statements.
//
//	DROP TRIGGER [IF EXISTS] name ON any_name opt_drop_behavior
//	DROP POLICY [IF EXISTS] name ON any_name opt_drop_behavior
//	DROP RULE [IF EXISTS] name ON any_name opt_drop_behavior
func (p *Parser) parseDropOnObject(objType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume TRIGGER/POLICY/RULE

	missingOk := p.parseOptIfExists()

	// parse name (the trigger/policy/rule name — not stored in Objects per the yacc grammar)
	_, _ = p.parseName()

	// ON
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}

	// any_name (the table name)
	anyName, err := p.parseAnyName()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{anyName}},
		RemoveType: int(objType),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropEventTrigger parses DROP EVENT TRIGGER statements.
//
//	DROP EVENT TRIGGER [IF EXISTS] name opt_drop_behavior
func (p *Parser) parseDropEventTrigger() (nodes.Node, error) {
	p.advance() // consume EVENT
	if _, err := p.expect(TRIGGER); err != nil {
		return nil, err
	}

	missingOk := p.parseOptIfExists()

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	nameList := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
	return &nodes.DropStmt{
		Objects:    makeNameListAsAnyNameList(nameList),
		RemoveType: int(nodes.OBJECT_EVENT_TRIGGER),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropLanguage parses DROP LANGUAGE statements.
//
//	DROP LANGUAGE [IF EXISTS] name opt_drop_behavior
func (p *Parser) parseDropLanguage() (nodes.Node, error) {
	p.advance() // consume LANGUAGE

	missingOk := p.parseOptIfExists()

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	nameList := &nodes.List{Items: []nodes.Node{&nodes.String{Str: name}}}
	return &nodes.DropStmt{
		Objects:    makeNameListAsAnyNameList(nameList),
		RemoveType: int(nodes.OBJECT_LANGUAGE),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropForeignDataWrapper parses DROP FOREIGN DATA WRAPPER statements.
//
//	DROP FOREIGN DATA WRAPPER [IF EXISTS] name_list opt_drop_behavior
func (p *Parser) parseDropForeignDataWrapper() (nodes.Node, error) {
	p.advance() // consume FOREIGN
	if _, err := p.expect(DATA_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return nil, err
	}

	missingOk := p.parseOptIfExists()

	nameList, err := p.parseNameList()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		Objects:    makeNameListAsAnyNameList(nameList),
		RemoveType: int(nodes.OBJECT_FDW),
		Behavior:   behavior,
		Missing_ok: missingOk,
	}, nil
}

// parseDropOwned parses DROP OWNED BY statements.
//
//	DropOwnedStmt:
//	    DROP OWNED BY role_list opt_drop_behavior
func (p *Parser) parseDropOwned() (nodes.Node, error) {
	p.advance() // consume OWNED
	if _, err := p.expect(BY); err != nil {
		return nil, err
	}

	roles := p.parseRoleList()

	behavior := p.parseOptDropBehavior()

	return &nodes.DropOwnedStmt{
		Roles:    roles,
		Behavior: nodes.DropBehavior(behavior),
	}, nil
}

// parseDropSubscription parses DROP SUBSCRIPTION statements.
//
//	DropSubscriptionStmt:
//	    DROP SUBSCRIPTION [IF EXISTS] name opt_drop_behavior
func (p *Parser) parseDropSubscription() (nodes.Node, error) {
	p.advance() // consume SUBSCRIPTION

	missingOk := p.parseOptIfExists()

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	behavior := p.parseOptDropBehavior()

	return &nodes.DropSubscriptionStmt{
		Subname:   name,
		MissingOk: missingOk,
		Behavior:  nodes.DropBehavior(behavior),
	}, nil
}

// parseOptIfExists parses the optional IF EXISTS clause.
//
//	opt_if_exists:
//	    IF_P EXISTS
//	    | /* EMPTY */
func (p *Parser) parseOptIfExists() bool {
	if p.cur.Type == IF_P {
		next := p.peekNext()
		if next.Type == EXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			return true
		}
	}
	return false
}

// makeNameListAsAnyNameList wraps each name in a name_list as a single-element
// any_name (a *List wrapping a *String), matching the yacc helper
// makeNameListAsAnyNameList.
func makeNameListAsAnyNameList(nameList *nodes.List) *nodes.List {
	if nameList == nil {
		return nil
	}
	result := &nodes.List{}
	for _, item := range nameList.Items {
		// Each item is a *String; wrap it in a *List (as a single-element any_name)
		result.Items = append(result.Items, &nodes.List{Items: []nodes.Node{item}})
	}
	return result
}

// parseTruncateStmt parses a TRUNCATE statement. The TRUNCATE keyword has already been consumed by parseStmt.
//
//	TruncateStmt:
//	    TRUNCATE opt_table relation_expr_list opt_restart_seqs opt_drop_behavior
func (p *Parser) parseTruncateStmt() (nodes.Node, error) {
	// opt_table
	if p.cur.Type == TABLE {
		p.advance()
	}

	// relation_expr_list
	relations := p.parseRelationExprList()

	// opt_restart_seqs
	restartSeqs := p.parseOptRestartSeqs()

	// opt_drop_behavior
	behavior := p.parseOptDropBehavior()

	return &nodes.TruncateStmt{
		Relations:   relations,
		RestartSeqs: restartSeqs,
		Behavior:    nodes.DropBehavior(behavior),
	}, nil
}

// parseRelationExprList parses a comma-separated list of relation_expr.
//
//	relation_expr_list:
//	    relation_expr
//	    | relation_expr_list ',' relation_expr
func (p *Parser) parseRelationExprList() *nodes.List {
	rel, _ := p.parseRelationExpr()
	if rel == nil {
		return nil
	}
	result := &nodes.List{Items: []nodes.Node{rel}}

	for p.cur.Type == ',' {
		p.advance()
		rel, _ = p.parseRelationExpr()
		if rel == nil {
			break
		}
		result.Items = append(result.Items, rel)
	}

	return result
}

// parseOptRestartSeqs parses the optional RESTART IDENTITY / CONTINUE IDENTITY clause.
//
//	opt_restart_seqs:
//	    CONTINUE_P IDENTITY_P  -> false
//	    | RESTART IDENTITY_P   -> true
//	    | /* EMPTY */          -> false
func (p *Parser) parseOptRestartSeqs() bool {
	if p.cur.Type == RESTART {
		p.advance()
		p.expect(IDENTITY_P)
		return true
	}
	if p.cur.Type == CONTINUE_P {
		p.advance()
		p.expect(IDENTITY_P)
		return false
	}
	return false
}

// parseDropOperatorClassOrFamily parses DROP OPERATOR CLASS/FAMILY.
// Current token is OPERATOR (DROP already consumed).
func (p *Parser) parseDropOperatorClassOrFamily() (nodes.Node, error) {
	p.advance() // consume OPERATOR

	var removeType nodes.ObjectType
	switch p.cur.Type {
	case CLASS:
		p.advance()
		removeType = nodes.OBJECT_OPCLASS
	case FAMILY:
		p.advance()
		removeType = nodes.OBJECT_OPFAMILY
	default:
		// DROP OPERATOR (not CLASS/FAMILY) - this is RemoveOperStmt
		return p.parseDropOperator()
	}

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance() // IF
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

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
	behavior := p.parseOptDropBehavior()

	// Prepend access method name to the object list (matches yacc grammar)
	items := []nodes.Node{&nodes.String{Str: amName}}
	if names != nil {
		items = append(items, names.Items...)
	}

	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: items}}},
		RemoveType: int(removeType),
		Behavior:   int(behavior),
		Missing_ok: missingOk,
	}, nil
}

// parseDropOperator parses DROP OPERATOR (RemoveOperStmt).
// Current token is the first token after OPERATOR (DROP OPERATOR already consumed).
//
// Ref: https://www.postgresql.org/docs/17/sql-dropoperator.html
//
//	DROP OPERATOR [ IF EXISTS ] operator_with_argtypes_list [ CASCADE | RESTRICT ]
//	DROP OPERATOR IF EXISTS operator_with_argtypes_list [ CASCADE | RESTRICT ]
func (p *Parser) parseDropOperator() (nodes.Node, error) {
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	objects := p.parseOperatorWithArgtypesList()
	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_OPERATOR),
		Objects:    objects,
		Behavior:   int(behavior),
		Missing_ok: missingOk,
	}, nil
}

// parseOperatorWithArgtypesList parses a comma-separated list of operator_with_argtypes.
func (p *Parser) parseOperatorWithArgtypesList() *nodes.List {
	owa := p.parseOperatorWithArgtypes()
	items := []nodes.Node{owa}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseOperatorWithArgtypes())
	}
	return &nodes.List{Items: items}
}

// parseOperatorWithArgtypes parses operator_with_argtypes: any_operator oper_argtypes.
func (p *Parser) parseOperatorWithArgtypes() *nodes.ObjectWithArgs {
	opName, _ := p.parseAnyOperator()
	owa := &nodes.ObjectWithArgs{Objname: opName}
	owa.Objargs, _ = p.parseOperArgtypes()
	return owa
}

// parseDropFuncStmt parses DROP FUNCTION/PROCEDURE/ROUTINE (RemoveFuncStmt).
// Current token is FUNCTION/PROCEDURE/ROUTINE.
//
//	DROP FUNCTION [IF EXISTS] function_with_argtypes_list [CASCADE|RESTRICT]
//	DROP PROCEDURE [IF EXISTS] function_with_argtypes_list [CASCADE|RESTRICT]
//	DROP ROUTINE [IF EXISTS] function_with_argtypes_list [CASCADE|RESTRICT]
func (p *Parser) parseDropFuncStmt(objType nodes.ObjectType) (nodes.Node, error) {
	p.advance() // consume FUNCTION/PROCEDURE/ROUTINE

	missingOk := p.parseOptIfExists()

	objects := p.parseFunctionWithArgtypesList()
	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		RemoveType: int(objType),
		Objects:    objects,
		Behavior:   int(behavior),
		Missing_ok: missingOk,
	}, nil
}

// parseFunctionWithArgtypesList parses a comma-separated list of function_with_argtypes.
func (p *Parser) parseFunctionWithArgtypesList() *nodes.List {
	fwa, _ := p.parseFunctionWithArgtypes()
	items := []nodes.Node{fwa}
	for p.cur.Type == ',' {
		p.advance()
		f, _ := p.parseFunctionWithArgtypes()
		items = append(items, f)
	}
	return &nodes.List{Items: items}
}

// parseDropAggrStmt parses DROP AGGREGATE (RemoveAggrStmt).
// Current token is AGGREGATE.
//
//	DROP AGGREGATE [IF EXISTS] aggregate_with_argtypes_list [CASCADE|RESTRICT]
func (p *Parser) parseDropAggrStmt() (nodes.Node, error) {
	p.advance() // consume AGGREGATE

	missingOk := p.parseOptIfExists()

	objects := p.parseAggregateWithArgtypesList()
	behavior := p.parseOptDropBehavior()

	return &nodes.DropStmt{
		RemoveType: int(nodes.OBJECT_AGGREGATE),
		Objects:    objects,
		Behavior:   int(behavior),
		Missing_ok: missingOk,
	}, nil
}

// parseAggregateWithArgtypesList parses a comma-separated list of aggregate_with_argtypes.
func (p *Parser) parseAggregateWithArgtypesList() *nodes.List {
	agg, _ := p.parseAggregateWithArgtypesLocal()
	items := []nodes.Node{agg}
	for p.cur.Type == ',' {
		p.advance()
		a, _ := p.parseAggregateWithArgtypesLocal()
		items = append(items, a)
	}
	return &nodes.List{Items: items}
}

// parseDropTableSpaceStmt parses a DROP TABLESPACE statement.
// The DROP keyword has already been consumed. Current token is TABLESPACE.
//
// Ref: https://www.postgresql.org/docs/17/sql-droptablespace.html
//
//	DROP TABLESPACE [ IF EXISTS ] name
func (p *Parser) parseDropTableSpaceStmt() (nodes.Node, error) {
	p.advance() // consume TABLESPACE

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.DropTableSpaceStmt{
		Tablespacename: name,
		MissingOk:      missingOk,
	}, nil
}
