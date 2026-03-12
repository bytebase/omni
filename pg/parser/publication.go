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
func (p *Parser) parseCreatePublicationStmt() nodes.Node {
	p.advance() // consume PUBLICATION

	name, _ := p.parseName()

	// Check for FOR
	if p.cur.Type == FOR {
		p.advance() // consume FOR

		// Check for ALL TABLES
		if p.cur.Type == ALL {
			p.advance() // consume ALL
			p.expect(TABLES)
			options := p.parseOptDefinition()
			return &nodes.CreatePublicationStmt{
				Pubname:      name,
				Options:      options,
				ForAllTables: true,
			}
		}

		// pub_obj_list
		pubobjects := p.parsePubObjList()
		options := p.parseOptDefinition()
		return &nodes.CreatePublicationStmt{
			Pubname:    name,
			Options:    options,
			Pubobjects: pubobjects,
		}
	}

	// Just name with optional definition
	options := p.parseOptDefinition()
	return &nodes.CreatePublicationStmt{
		Pubname: name,
		Options: options,
	}
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
func (p *Parser) parseAlterPublicationStmt() nodes.Node {
	p.advance() // consume PUBLICATION

	name, _ := p.parseName()

	switch p.cur.Type {
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_PUBLICATION,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_PUBLICATION,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}
	case SET:
		// Could be SET definition (SET '(' ...) or SET pub_obj_list
		next := p.peekNext()
		if next.Type == '(' {
			// SET definition
			p.advance() // consume SET
			p.advance() // consume '('
			list := p.parseDefList()
			p.expect(')')
			return &nodes.AlterPublicationStmt{
				Pubname: name,
				Options: list,
			}
		}
		// SET pub_obj_list
		p.advance() // consume SET
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_SET,
		}
	case ADD_P:
		p.advance() // consume ADD
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_ADD,
		}
	case DROP:
		p.advance() // consume DROP
		pubobjects := p.parsePubObjList()
		return &nodes.AlterPublicationStmt{
			Pubname:    name,
			Pubobjects: pubobjects,
			Action:     nodes.DEFELEM_DROP,
		}
	default:
		return nil
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
		rel := p.parseRelationExpr()
		cols := p.parseOptColumnList()
		where := p.parseOptWhereClausePub()
		pt := &nodes.PublicationTable{
			Relation: rel,
			Columns:  cols,
		}
		if where != nil {
			pt.WhereClause = where
		}
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLE,
			Pubtable:   pt,
			Loc: nodes.Loc{Start: loc, End: -1},
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
				Loc: nodes.Loc{Start: loc, End: -1},
			}
		}
		schemaName, _ := p.parseColId()
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA,
			Name:       schemaName,
			Loc: nodes.Loc{Start: loc, End: -1},
		}
	}

	if p.cur.Type == CURRENT_SCHEMA {
		p.advance()
		return &nodes.PublicationObjSpec{
			Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
			Loc: nodes.Loc{Start: loc, End: -1},
		}
	}

	// CONTINUATION: relation_expr opt_column_list OptWhereClause
	rel := p.parseRelationExpr()
	cols := p.parseOptColumnList()
	where := p.parseOptWhereClausePub()
	pt := &nodes.PublicationTable{
		Relation: rel,
		Columns:  cols,
	}
	if where != nil {
		pt.WhereClause = where
	}
	return &nodes.PublicationObjSpec{
		Pubobjtype: nodes.PUBLICATIONOBJ_CONTINUATION,
		Pubtable:   pt,
		Loc: nodes.Loc{Start: loc, End: -1},
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
	expr := p.parseAExpr(0)
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
func (p *Parser) parseCreateSubscriptionStmt() nodes.Node {
	p.advance() // consume SUBSCRIPTION

	name, _ := p.parseName()
	p.expect(CONNECTION)
	conninfo := p.cur.Str
	p.expect(SCONST)
	p.expect(PUBLICATION)
	pubList, _ := p.parseNameList()
	options := p.parseOptDefinition()

	return &nodes.CreateSubscriptionStmt{
		Subname:     name,
		Conninfo:    conninfo,
		Publication: pubList,
		Options:     options,
	}
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
func (p *Parser) parseAlterSubscriptionStmt() nodes.Node {
	p.advance() // consume SUBSCRIPTION

	name, _ := p.parseName()

	switch p.cur.Type {
	case OWNER:
		p.advance()
		p.expect(TO)
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_SUBSCRIPTION,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}
	case RENAME:
		p.advance()
		p.expect(TO)
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_SUBSCRIPTION,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}

	case CONNECTION:
		p.advance() // consume CONNECTION
		conninfo := p.cur.Str
		p.expect(SCONST)
		return &nodes.AlterSubscriptionStmt{
			Kind:     nodes.ALTER_SUBSCRIPTION_CONNECTION,
			Subname:  name,
			Conninfo: conninfo,
		}

	case REFRESH:
		p.advance() // consume REFRESH
		p.expect(PUBLICATION)
		options := p.parseOptDefinition()
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_REFRESH,
			Subname: name,
			Options: options,
		}

	case ENABLE_P:
		p.advance() // consume ENABLE
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
			Subname: name,
			Options: &nodes.List{Items: []nodes.Node{
				makeDefElem("enabled", &nodes.Boolean{Boolval: true}),
			}},
		}

	case DISABLE_P:
		p.advance() // consume DISABLE
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_ENABLED,
			Subname: name,
			Options: &nodes.List{Items: []nodes.Node{
				makeDefElem("enabled", &nodes.Boolean{Boolval: false}),
			}},
		}

	case SKIP:
		p.advance() // consume SKIP
		p.expect('(')
		list := p.parseDefList()
		p.expect(')')
		return &nodes.AlterSubscriptionStmt{
			Kind:    nodes.ALTER_SUBSCRIPTION_SKIP,
			Subname: name,
			Options: list,
		}

	case SET:
		// Could be SET definition or SET PUBLICATION
		next := p.peekNext()
		if next.Type == '(' {
			// SET definition
			p.advance() // consume SET
			p.advance() // consume '('
			list := p.parseDefList()
			p.expect(')')
			return &nodes.AlterSubscriptionStmt{
				Kind:    nodes.ALTER_SUBSCRIPTION_OPTIONS,
				Subname: name,
				Options: list,
			}
		}
		if next.Type == PUBLICATION {
			p.advance() // consume SET
			p.advance() // consume PUBLICATION
			pubList, _ := p.parseNameList()
			options := p.parseOptDefinition()
			return &nodes.AlterSubscriptionStmt{
				Kind:        nodes.ALTER_SUBSCRIPTION_SET_PUBLICATION,
				Subname:     name,
				Publication: pubList,
				Options:     options,
			}
		}
		return nil

	case ADD_P:
		p.advance() // consume ADD
		p.expect(PUBLICATION)
		pubList, _ := p.parseNameList()
		options := p.parseOptDefinition()
		return &nodes.AlterSubscriptionStmt{
			Kind:        nodes.ALTER_SUBSCRIPTION_ADD_PUBLICATION,
			Subname:     name,
			Publication: pubList,
			Options:     options,
		}

	case DROP:
		p.advance() // consume DROP
		p.expect(PUBLICATION)
		pubList, _ := p.parseNameList()
		options := p.parseOptDefinition()
		return &nodes.AlterSubscriptionStmt{
			Kind:        nodes.ALTER_SUBSCRIPTION_DROP_PUBLICATION,
			Subname:     name,
			Publication: pubList,
			Options:     options,
		}

	default:
		return nil
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
func (p *Parser) parseCreateRuleStmt(replace bool) nodes.Node {
	p.advance() // consume RULE

	name, _ := p.parseName()
	p.expect(AS)
	p.expect(ON)

	// event: SELECT | UPDATE | DELETE | INSERT
	event := p.parseRuleEvent()

	p.expect(TO)

	// qualified_name -> RangeVar
	names, _ := p.parseQualifiedName()
	rel := makeRangeVarFromAnyName(names)

	// where_clause (optional)
	whereClause := p.parseWhereClause()

	p.expect(DO)

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
	}
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
	stmt := p.parseRuleActionStmt()
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
	stmt := p.parseRuleActionStmt()
	if stmt != nil {
		items = append(items, stmt)
	}

	for p.cur.Type == ';' {
		p.advance()
		// RuleActionStmtOrEmpty - could be empty (just ';')
		if p.cur.Type == ')' || p.cur.Type == ';' || p.cur.Type == 0 {
			continue
		}
		stmt = p.parseRuleActionStmt()
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
func (p *Parser) parseRuleActionStmt() nodes.Node {
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
		return nil
	}
}
