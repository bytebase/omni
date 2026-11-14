package parser

import (
	nodes "github.com/bytebase/omni/redshift/ast"
)

// =============================================================================
// CREATE PUBLICATION
// =============================================================================

// parseCreatePublicationStmt parses a CREATE PUBLICATION statement.
// The CREATE keyword has already been consumed; current token is PUBLICATION.
//
//	CreatePublicationStmt:
//	    CREATE PUBLICATION name opt_definition
//	    | CREATE PUBLICATION name FOR ALL TABLES opt_definition
//	    | CREATE PUBLICATION name FOR pub_obj_list opt_definition
func (p *Parser) parseCreatePublicationStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume PUBLICATION

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Check for FOR
	if p.cur.Type == FOR {
		p.advance() // consume FOR

		// Check for ALL TABLES
		if p.cur.Type == ALL {
			p.advance() // consume ALL
			if _, err := p.expect(TABLES); err != nil {
				return nil, err
			}
			options, err := p.parseOptDefinition()
			if err != nil {
				return nil, err
			}
			return &nodes.CreatePublicationStmt{
				Pubname:      name,
				Options:      options,
				ForAllTables: true,
				Loc:          nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}

		// pub_obj_list
		pubobjects, shapes, err := p.parsePubObjList()
		if err != nil {
			return nil, err
		}
		options, err := p.parseOptDefinition()
		if err != nil {
			return nil, err
		}
		if err := p.checkPubObjList(pubobjects, shapes); err != nil {
			return nil, err
		}
		return &nodes.CreatePublicationStmt{
			Pubname:    name,
			Options:    options,
			Pubobjects: pubobjects,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}

	// Just name with optional definition
	options, err := p.parseOptDefinition()
	if err != nil {
		return nil, err
	}
	return &nodes.CreatePublicationStmt{
		Pubname: name,
		Options: options,
		Loc:     nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// =============================================================================
// ALTER PUBLICATION
// =============================================================================

// parseAlterPublicationStmt parses an ALTER PUBLICATION statement.
// The ALTER keyword has already been consumed; current token is PUBLICATION.
//
//	AlterPublicationStmt:
//	    ALTER PUBLICATION name SET definition
//	    | ALTER PUBLICATION name ADD_P pub_obj_list
//	    | ALTER PUBLICATION name SET pub_obj_list
//	    | ALTER PUBLICATION name DROP pub_obj_list
func (p *Parser) parseAlterPublicationStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance()       // consume PUBLICATION

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec, err := p.parseRoleSpec()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_PUBLICATION,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil
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
			RenameType: nodes.OBJECT_PUBLICATION,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	case SET:
		// Could be SET definition (SET '(' ...) or SET pub_obj_list
		next := p.peekNext()
		if next.Type == '(' {
			// SET definition
			p.advance() // consume SET
			p.advance() // consume '('
			list, err := p.parseDefList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.AlterPublicationStmt{
				Pubname: name,
				Options: list,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		// SET pub_obj_list
		p.advance() // consume SET
		pubobjects, shapes, err := p.parsePubObjList()
		if err != nil {
			return nil, err
		}
		if err := p.checkPubObjList(pubobjects, shapes); err != nil {
			return nil, err
		}
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_SET,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case ADD_P:
		p.advance() // consume ADD
		pubobjects, shapes, err := p.parsePubObjList()
		if err != nil {
			return nil, err
		}
		if err := p.checkPubObjList(pubobjects, shapes); err != nil {
			return nil, err
		}
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_ADD,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case DROP:
		p.advance() // consume DROP
		pubobjects, shapes, err := p.parsePubObjList()
		if err != nil {
			return nil, err
		}
		if err := p.checkPubObjList(pubobjects, shapes); err != nil {
			return nil, err
		}
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_DROP,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		if p.collectMode() {
			for _, t := range []int{OWNER, RENAME, SET, ADD_P, DROP} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// =============================================================================
// pub_obj_list / PublicationObjSpec
// =============================================================================

// pubObjShape records how a PublicationObjSpec was written, which
// checkPubObjList needs and the AST does not keep.
type pubObjShape struct {
	loc      int  // gram.y's pubobj->location
	name     bool // a bare ColId: gram.y keeps it as a name
	pubtable bool // gram.y built a PublicationTable for it
	cols     bool
	where    bool
}

// parsePubObjList parses pub_obj_list. The caller checks the list with
// checkPubObjList where gram.y calls preprocess_pubobj_list.
//
//	pub_obj_list:
//	    PublicationObjSpec
//	    | pub_obj_list ',' PublicationObjSpec
func (p *Parser) parsePubObjList() (*nodes.List, []pubObjShape, error) {
	var items []nodes.Node
	var shapes []pubObjShape
	for {
		obj, shape, err := p.parsePublicationObjSpec()
		if err != nil {
			return nil, nil, err
		}
		items = append(items, obj)
		shapes = append(shapes, shape)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, shapes, nil
		}
		p.advance()
	}
}

// checkPubObjList reports the errors gram.y's preprocess_pubobj_list raises
// as syntax errors: a list that starts with a bare name, a bare
// CURRENT_SCHEMA after TABLE, and a table-only clause on a schema.
func (p *Parser) checkPubObjList(list *nodes.List, shapes []pubObjShape) error {
	if p.collectMode() || list == nil {
		return nil
	}
	prev := nodes.PUBLICATIONOBJ_CONTINUATION
	for i, item := range list.Items {
		obj := item.(*nodes.PublicationObjSpec)
		shape := shapes[i]
		objtype := obj.Pubobjtype
		if objtype == nodes.PUBLICATIONOBJ_CONTINUATION {
			if i == 0 {
				return &ParseError{Message: "invalid publication object list", Position: shape.loc}
			}
			objtype = prev
		}
		switch objtype {
		case nodes.PUBLICATIONOBJ_TABLE:
			if !shape.name && !shape.pubtable {
				return &ParseError{Message: "invalid table name", Position: shape.loc}
			}
		case nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA, nodes.PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA:
			switch {
			case shape.where:
				return &ParseError{Message: "WHERE clause not allowed for schema", Position: shape.loc}
			case shape.cols:
				return &ParseError{Message: "column specification not allowed for schema", Position: shape.loc}
			case shape.pubtable:
				return &ParseError{Message: "invalid schema name", Position: shape.loc}
			}
			if shape.name {
				objtype = nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA
			} else {
				objtype = nodes.PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA
			}
		}
		prev = objtype
	}
	return nil
}

// parsePublicationObjSpec parses PublicationObjSpec.
//
//	PublicationObjSpec:
//	    TABLE relation_expr opt_column_list OptWhereClause
//	    | TABLES IN_P SCHEMA ColId
//	    | TABLES IN_P SCHEMA CURRENT_SCHEMA
//	    | ColId opt_column_list OptWhereClause                 (CONTINUATION)
//	    | ColId indirection opt_column_list OptWhereClause     (CONTINUATION)
//	    | extended_relation_expr opt_column_list OptWhereClause (CONTINUATION)
//	    | CURRENT_SCHEMA                                       (CONTINUATION)
//
// TABLES is an unreserved keyword, so without IN after it it is the ColId of
// a continuation.
func (p *Parser) parsePublicationObjSpec() (*nodes.PublicationObjSpec, pubObjShape, error) {
	loc := p.pos()

	if p.cur.Type == TABLE {
		p.advance() // consume TABLE
		pt, _, err := p.parsePublicationTable()
		if err != nil {
			return nil, pubObjShape{}, err
		}
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLE,
			Pubtable:   pt,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, pubObjShape{loc: loc, pubtable: true}, nil
	}

	if p.cur.Type == TABLES && (p.peekNext().Type == IN_P || p.completingAtNext()) {
		p.advance() // consume TABLES
		if _, err := p.expect(IN_P); err != nil {
			return nil, pubObjShape{}, err
		}
		if _, err := p.expect(SCHEMA); err != nil {
			return nil, pubObjShape{}, err
		}
		schemaLoc := p.pos()
		if p.cur.Type == CURRENT_SCHEMA && !p.collectMode() {
			p.advance()
			return &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA,
				Loc:        nodes.Loc{Start: loc, End: p.prev.End},
			}, pubObjShape{loc: schemaLoc}, nil
		}
		schemaName, err := p.parseColId()
		if err != nil {
			return nil, pubObjShape{}, err
		}
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA,
			Name:       schemaName,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, pubObjShape{loc: schemaLoc, name: true}, nil
	}

	if p.cur.Type == CURRENT_SCHEMA && !p.collectMode() {
		p.advance()
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, pubObjShape{loc: loc}, nil
	}

	// CONTINUATION: a table or schema name, with an optional column list
	// and row filter.
	pt, extended, err := p.parsePublicationTable()
	if err != nil {
		return nil, pubObjShape{}, err
	}
	rel := pt.Relation
	shape := pubObjShape{loc: loc, cols: pt.Columns != nil, where: pt.WhereClause != nil}
	switch {
	case extended:
		// gram.y's extended_relation_expr alternative never sets
		// pubobj->location, so PostgreSQL reports it at offset 0 of the
		// query; omni points at the statement instead.
		shape.loc = p.stmtStart
		shape.pubtable = true
	case rel.Schemaname != "" || rel.Catalogname != "":
		shape.pubtable = true
	default:
		shape.pubtable = shape.cols || shape.where
		shape.name = !shape.pubtable
	}
	return &nodes.PublicationObjSpec{
		Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
		Pubtable:   pt,
		Loc:        nodes.Loc{Start: loc, End: p.prev.End},
	}, shape, nil
}

// parsePublicationTable parses relation_expr opt_column_list OptWhereClause,
// and reports whether the relation_expr was an extended_relation_expr
// (name *, ONLY name, ONLY (name)).
func (p *Parser) parsePublicationTable() (*nodes.PublicationTable, bool, error) {
	ptLoc := p.pos()
	extended := p.cur.Type == ONLY
	rel, err := p.parseRelationExpr()
	if err != nil {
		return nil, false, err
	}
	extended = extended || p.prev.Type == '*'
	cols, err := p.parseOptColumnList()
	if err != nil {
		return nil, false, err
	}
	where, err := p.parseOptWhereClausePub()
	if err != nil {
		return nil, false, err
	}
	return &nodes.PublicationTable{
		Relation:    rel,
		Columns:     cols,
		WhereClause: where,
		Loc:         nodes.Loc{Start: ptLoc, End: p.prev.End},
	}, extended, nil
}

// parseOptWhereClausePub parses the OptWhereClause for publication objects.
//
//	OptWhereClause: WHERE '(' a_expr ')' | /* EMPTY */
func (p *Parser) parseOptWhereClausePub() (nodes.Node, error) {
	if p.cur.Type != WHERE {
		return nil, nil
	}
	p.advance() // consume WHERE
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseAExpr(0)
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return expr, nil
}

// =============================================================================
// CREATE SUBSCRIPTION
// =============================================================================

// parseCreateSubscriptionStmt parses a CREATE SUBSCRIPTION statement.
// The CREATE keyword has already been consumed; current token is SUBSCRIPTION.
//
//	CreateSubscriptionStmt:
//	    CREATE SUBSCRIPTION name CONNECTION Sconst PUBLICATION name_list opt_definition
func (p *Parser) parseCreateSubscriptionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of CREATE
	p.advance()       // consume SUBSCRIPTION

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(CONNECTION); err != nil {
		return nil, err
	}
	conninfo := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}
	if _, err := p.expect(PUBLICATION); err != nil {
		return nil, err
	}
	pubList, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	options, err := p.parseOptDefinition()
	if err != nil {
		return nil, err
	}

	return &nodes.CreateSubscriptionStmt{
		Subname:     name,
		Conninfo:    conninfo,
		Publication: pubList,
		Options:     options,
		Loc:         nodes.Loc{Start: loc, End: p.prev.End},
	}, nil
}

// =============================================================================
// ALTER SUBSCRIPTION
// =============================================================================

// parseAlterSubscriptionStmt parses an ALTER SUBSCRIPTION statement.
// The ALTER keyword has already been consumed; current token is SUBSCRIPTION.
//
//	AlterSubscriptionStmt:
//	    ALTER SUBSCRIPTION name SET definition
//	    | ALTER SUBSCRIPTION name CONNECTION Sconst
//	    | ALTER SUBSCRIPTION name REFRESH PUBLICATION opt_definition
//	    | ALTER SUBSCRIPTION name ADD_P PUBLICATION name_list opt_definition
//	    | ALTER SUBSCRIPTION name DROP PUBLICATION name_list opt_definition
//	    | ALTER SUBSCRIPTION name SET PUBLICATION name_list opt_definition
//	    | ALTER SUBSCRIPTION name ENABLE_P
//	    | ALTER SUBSCRIPTION name DISABLE_P
//	    | ALTER SUBSCRIPTION name SKIP definition
func (p *Parser) parseAlterSubscriptionStmt() (nodes.Node, error) {
	loc := p.prev.Loc // start of ALTER
	p.advance()       // consume SUBSCRIPTION

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec, err := p.parseRoleSpec()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_SUBSCRIPTION,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil
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
			RenameType: nodes.OBJECT_SUBSCRIPTION,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil

	case CONNECTION:
		p.advance() // consume CONNECTION
		conninfo := p.cur.Str
		if _, err := p.expect(SCONST); err != nil {
			return nil, err
		}
		return &nodes.AlterSubscriptionStmt{
			Kind:     nodes.ALTER_SUBSCRIPTION_CONNECTION,
			Subname:  name,
			Conninfo: conninfo,
			Loc:      nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case REFRESH:
		p.advance() // consume REFRESH
		if _, err := p.expect(PUBLICATION); err != nil {
			return nil, err
		}
		options, err := p.parseOptDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_REFRESH,
			Subname: name,
			Options: options,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case ENABLE_P:
		p.advance() // consume ENABLE
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
			Subname: name,
			Options: &nodes.List{Items: []nodes.Node{
				makeDefElem("enabled", &nodes.Boolean{Boolval: true}),
			}},
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case DISABLE_P:
		p.advance() // consume DISABLE
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
			Subname: name,
			Options: &nodes.List{Items: []nodes.Node{
				makeDefElem("enabled", &nodes.Boolean{Boolval: false}),
			}},
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case SKIP:
		p.advance() // consume SKIP
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
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_SKIP,
			Subname: name,
			Options: list,
			Loc:     nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case SET:
		// Could be SET definition or SET PUBLICATION
		next := p.peekNext()
		if next.Type == '(' {
			// SET definition
			p.advance() // consume SET
			p.advance() // consume '('
			list, err := p.parseDefList()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			return &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_OPTIONS,
				Subname: name,
				Options: list,
				Loc:     nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		if next.Type == PUBLICATION {
			p.advance() // consume SET
			p.advance() // consume PUBLICATION
			pubList, err := p.parseNameList()
			if err != nil {
				return nil, err
			}
			options, err := p.parseOptDefinition()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_SET_PUBLICATION,
				Subname:     name,
				Publication: pubList,
				Options:     options,
				Loc:         nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		p.advance() // consume SET
		if p.collectMode() {
			p.addTokenCandidate('(')
			p.addTokenCandidate(PUBLICATION)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()

	case ADD_P:
		p.advance() // consume ADD
		if _, err := p.expect(PUBLICATION); err != nil {
			return nil, err
		}
		pubList, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		options, err := p.parseOptDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterSubscriptionStmt{
			Kind:        nodes.ALTER_SUBSCRIPTION_ADD_PUBLICATION,
			Subname:     name,
			Publication: pubList,
			Options:     options,
			Loc:         nodes.Loc{Start: loc, End: p.prev.End},
		}, nil

	case DROP:
		p.advance() // consume DROP
		if _, err := p.expect(PUBLICATION); err != nil {
			return nil, err
		}
		pubList, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		options, err := p.parseOptDefinition()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterSubscriptionStmt{
			Kind:        nodes.ALTER_SUBSCRIPTION_DROP_PUBLICATION,
			Subname:     name,
			Publication: pubList,
			Loc:         nodes.Loc{Start: loc, End: p.prev.End},
			Options:     options,
		}, nil

	default:
		if p.collectMode() {
			for _, t := range []int{OWNER, RENAME, CONNECTION, REFRESH, ENABLE_P, DISABLE_P, SKIP, SET, ADD_P, DROP} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
}

// =============================================================================
// CREATE RULE
// =============================================================================

// parseCreateRuleStmt parses a CREATE RULE statement.
// The CREATE keyword has been consumed; OR REPLACE may have been consumed.
// Current token is RULE.
//
//	RuleStmt:
//	    CREATE opt_or_replace RULE name AS
//	    ON event TO qualified_name where_clause
//	    DO opt_instead RuleActionList
func (p *Parser) parseCreateRuleStmt(replace bool, stmtLoc int) (nodes.Node, error) {
	p.advance() // consume RULE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(AS); err != nil {
		return nil, err
	}
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}

	// event: SELECT | UPDATE | DELETE | INSERT
	event, err := p.parseRuleEvent()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(TO); err != nil {
		return nil, err
	}

	// qualified_name -> RangeVar
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rel := makeRangeVarFromAnyName(names, nodes.Loc{Start: nameStart, End: p.prev.End})

	// where_clause (optional)
	whereClause, err := p.parseWhereClause()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(DO); err != nil {
		return nil, err
	}

	// opt_instead: INSTEAD | ALSO | /* EMPTY */
	instead := false
	if p.cur.Type == INSTEAD {
		p.advance()
		instead = true
	} else if p.cur.Type == ALSO {
		p.advance()
		instead = false
	}

	// RuleActionList: NOTHING | RuleActionStmt | '(' RuleActionMulti ')'
	actions, err := p.parseRuleActionList()
	if err != nil {
		return nil, err
	}

	return &nodes.RuleStmt{
		Replace:     replace,
		Relation:    rel,
		Rulename:    name,
		WhereClause: whereClause,
		Event:       nodes.CmdType(event),
		Instead:     instead,
		Actions:     actions,
		Loc:         nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseRuleEvent parses the event keyword for a rule.
//
//	event: SELECT | UPDATE | DELETE_P | INSERT
func (p *Parser) parseRuleEvent() (int, error) {
	switch p.cur.Type {
	case SELECT:
		p.advance()
		return int(nodes.CMD_SELECT), nil
	case UPDATE:
		p.advance()
		return int(nodes.CMD_UPDATE), nil
	case DELETE_P:
		p.advance()
		return int(nodes.CMD_DELETE), nil
	case INSERT:
		p.advance()
		return int(nodes.CMD_INSERT), nil
	}
	if p.collectMode() {
		for _, t := range []int{SELECT, UPDATE, DELETE_P, INSERT} {
			p.addTokenCandidate(t)
		}
		return 0, errCollecting
	}
	return 0, p.syntaxErrorAtCur()
}

// parseRuleActionList parses RuleActionList.
//
//	RuleActionList:
//	    NOTHING
//	    | RuleActionStmt
//	    | '(' RuleActionMulti ')'
func (p *Parser) parseRuleActionList() (*nodes.List, error) {
	if p.cur.Type == NOTHING {
		p.advance()
		return nil, nil
	}

	if p.cur.Type == '(' {
		p.advance() // consume '('
		actions, err := p.parseRuleActionMulti()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		return actions, nil
	}

	// Single RuleActionStmt
	stmt, err := p.parseRuleActionStmt()
	if err != nil {
		return nil, err
	}
	return &nodes.List{Items: []nodes.Node{stmt}}, nil
}

// parseRuleActionMulti parses RuleActionMulti.
//
//	RuleActionMulti:
//	    RuleActionMulti ';' RuleActionStmtOrEmpty
//	    | RuleActionStmtOrEmpty
func (p *Parser) parseRuleActionMulti() (*nodes.List, error) {
	var items []nodes.Node
	for {
		// RuleActionStmtOrEmpty: RuleActionStmt | /*EMPTY*/
		if p.cur.Type != ')' && p.cur.Type != ';' {
			stmt, err := p.parseRuleActionStmt()
			if err != nil {
				return nil, err
			}
			items = append(items, stmt)
		}
		if p.cur.Type != ';' {
			break
		}
		p.advance()
	}

	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseRuleActionStmt parses a single rule action statement.
//
//	RuleActionStmt:
//	    SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | NotifyStmt
func (p *Parser) parseRuleActionStmt() (nodes.Node, error) {
	if p.isSelectStart() {
		// WITH is in the SELECT FIRST set but dispatches to parseWithStmt
		// rather than parseSelectNoParens, so handle it as a sub-branch
		// inside the isSelectStart() gate.
		var stmt nodes.Node
		var err error
		if p.cur.Type == WITH {
			stmt, err = p.parseWithStmt()
		} else {
			stmt, err = p.parseSelectNoParens()
		}
		if err == nil && stmt == nil {
			err = p.syntaxErrorAtCur()
		}
		return stmt, err
	}
	switch p.cur.Type {
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	case NOTIFY:
		p.advance() // consume NOTIFY
		return p.parseNotifyStmt()
	}
	return nil, p.syntaxErrorAtCur()
}
