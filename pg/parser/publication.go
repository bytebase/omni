package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
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
			options := p.parseOptDefinition()
			return &nodes.CreatePublicationStmt{
				Pubname:      name,
				Options:      options,
				ForAllTables: true,
				Loc:          nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}

		// pub_obj_list
		pubobjects := p.parsePubObjList()
		options := p.parseOptDefinition()
		return &nodes.CreatePublicationStmt{
			Pubname:    name,
			Options:    options,
			Pubobjects: pubobjects,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	}

	// Just name with optional definition
	options := p.parseOptDefinition()
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
		roleSpec := p.parseRoleSpec()
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
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_SET,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case ADD_P:
		p.advance() // consume ADD
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_ADD,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	case DROP:
		p.advance() // consume DROP
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_DROP,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}, nil
	default:
		return nil, nil
	}
}

// =============================================================================
// pub_obj_list / PublicationObjSpec
// =============================================================================

// parsePubObjList parses pub_obj_list.
//
//	pub_obj_list:
//	    PublicationObjSpec
//	    | pub_obj_list ',' PublicationObjSpec
func (p *Parser) parsePubObjList() *nodes.List {
	obj := p.parsePublicationObjSpec()
	items := []nodes.Node{obj}
	for p.cur.Type == ',' {
		p.advance()
		obj = p.parsePublicationObjSpec()
		items = append(items, obj)
	}
	return &nodes.List{Items: items}
}

// parsePublicationObjSpec parses PublicationObjSpec.
//
//	PublicationObjSpec:
//	    TABLE relation_expr opt_column_list OptWhereClause
//	    | TABLES IN_P SCHEMA ColId
//	    | TABLES IN_P SCHEMA CURRENT_SCHEMA
//	    | relation_expr opt_column_list OptWhereClause     (CONTINUATION)
//	    | CURRENT_SCHEMA                                   (CONTINUATION)
func (p *Parser) parsePublicationObjSpec() *nodes.PublicationObjSpec {
	loc := p.pos()

	if p.cur.Type == TABLE {
		p.advance() // consume TABLE
		ptLoc := p.pos()
		rel, _ := p.parseRelationExpr()
		cols, _ := p.parseOptColumnList()
		where := p.parseOptWhereClausePub()
		pt := &nodes.PublicationTable{
			Relation: rel,
			Columns:  cols,
			Loc:      nodes.Loc{Start: ptLoc, End: p.prev.End},
		}
		if where != nil {
			pt.WhereClause = where
			pt.Loc.End = p.prev.End
		}
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLE,
			Pubtable:   pt,
			Loc:        nodes.Loc{Start: loc, End: p.prev.End},
		}
	}

	if p.cur.Type == TABLES {
		p.advance() // consume TABLES
		p.expect(IN_P)
		p.expect(SCHEMA)
		if p.cur.Type == CURRENT_SCHEMA {
			p.advance()
			return &nodes.PublicationObjSpec{
				Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_CUR_SCHEMA,
				Loc: nodes.Loc{Start: loc, End: p.prev.End},
			}
		}
		schemaName, _ := p.parseColId()
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA,
			Name:       schemaName,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}
	}

	if p.cur.Type == CURRENT_SCHEMA {
		p.advance()
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
			Loc: nodes.Loc{Start: loc, End: p.prev.End},
		}
	}

	// CONTINUATION: relation_expr opt_column_list OptWhereClause
	ptLoc := p.pos()
	rel, _ := p.parseRelationExpr()
	cols, _ := p.parseOptColumnList()
	where := p.parseOptWhereClausePub()
	pt := &nodes.PublicationTable{
		Relation: rel,
		Columns:  cols,
		Loc:      nodes.Loc{Start: ptLoc, End: p.prev.End},
	}
	if where != nil {
		pt.WhereClause = where
		pt.Loc.End = p.prev.End
	}
	return &nodes.PublicationObjSpec{
		Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
		Pubtable:   pt,
		Loc: nodes.Loc{Start: loc, End: p.prev.End},
	}
}

// parseOptWhereClausePub parses the OptWhereClause for publication objects.
//
//	OptWhereClause: WHERE '(' a_expr ')' | /* EMPTY */
func (p *Parser) parseOptWhereClausePub() nodes.Node {
	if p.cur.Type != WHERE {
		return nil
	}
	p.advance() // consume WHERE
	p.expect('(')
	expr, _ := p.parseAExpr(0)
	p.expect(')')
	return expr
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
	options := p.parseOptDefinition()

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
		roleSpec := p.parseRoleSpec()
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
		options := p.parseOptDefinition()
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
			options := p.parseOptDefinition()
			return &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_SET_PUBLICATION,
				Subname:     name,
				Publication: pubList,
				Options:     options,
				Loc:         nodes.Loc{Start: loc, End: p.prev.End},
			}, nil
		}
		return nil, nil

	case ADD_P:
		p.advance() // consume ADD
		if _, err := p.expect(PUBLICATION); err != nil {
			return nil, err
		}
		pubList, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		options := p.parseOptDefinition()
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
		options := p.parseOptDefinition()
		return &nodes.AlterSubscriptionStmt{
			Kind:        nodes.ALTER_SUBSCRIPTION_DROP_PUBLICATION,
			Subname:     name,
			Publication: pubList,
			Loc:         nodes.Loc{Start: loc, End: p.prev.End},
			Options:     options,
		}, nil

	default:
		return nil, nil
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
	event := p.parseRuleEvent()

	if _, err := p.expect(TO); err != nil {
		return nil, err
	}

	// qualified_name -> RangeVar
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rel := makeRangeVarFromAnyName(names)

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
	actions := p.parseRuleActionList()

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
func (p *Parser) parseRuleEvent() int {
	switch p.cur.Type {
	case SELECT:
		p.advance()
		return int(nodes.CMD_SELECT)
	case UPDATE:
		p.advance()
		return int(nodes.CMD_UPDATE)
	case DELETE_P:
		p.advance()
		return int(nodes.CMD_DELETE)
	case INSERT:
		p.advance()
		return int(nodes.CMD_INSERT)
	default:
		p.advance() // consume whatever is there
		return int(nodes.CMD_SELECT)
	}
}

// parseRuleActionList parses RuleActionList.
//
//	RuleActionList:
//	    NOTHING
//	    | RuleActionStmt
//	    | '(' RuleActionMulti ')'
func (p *Parser) parseRuleActionList() *nodes.List {
	if p.cur.Type == NOTHING {
		p.advance()
		return nil
	}

	if p.cur.Type == '(' {
		p.advance() // consume '('
		actions := p.parseRuleActionMulti()
		p.expect(')')
		return actions
	}

	// Single RuleActionStmt
	stmt, _ := p.parseRuleActionStmt()
	if stmt == nil {
		return nil
	}
	return &nodes.List{Items: []nodes.Node{stmt}}
}

// parseRuleActionMulti parses RuleActionMulti.
//
//	RuleActionMulti:
//	    RuleActionMulti ';' RuleActionStmtOrEmpty
//	    | RuleActionStmtOrEmpty
func (p *Parser) parseRuleActionMulti() *nodes.List {
	var items []nodes.Node

	// First item
	stmt, _ := p.parseRuleActionStmt()
	if stmt != nil {
		items = append(items, stmt)
	}

	for p.cur.Type == ';' {
		p.advance()
		// RuleActionStmtOrEmpty - could be empty (just ';')
		if p.cur.Type == ')' || p.cur.Type == ';' || p.cur.Type == 0 {
			continue
		}
		stmt, _ = p.parseRuleActionStmt()
		if stmt != nil {
			items = append(items, stmt)
		}
	}

	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseRuleActionStmt parses a single rule action statement.
//
//	RuleActionStmt:
//	    SelectStmt | InsertStmt | UpdateStmt | DeleteStmt | NotifyStmt
func (p *Parser) parseRuleActionStmt() (nodes.Node, error) {
	switch p.cur.Type {
	case SELECT, VALUES, TABLE, WITH:
		if p.cur.Type == WITH {
			return p.parseWithStmt()
		}
		return p.parseSelectNoParens()
	case INSERT:
		return p.parseInsertStmt(nil)
	case UPDATE:
		return p.parseUpdateStmt(nil)
	case DELETE_P:
		return p.parseDeleteStmt(nil)
	default:
		return nil, nil
	}
}
