package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateSchemaStmt parses a CREATE SCHEMA statement.
//
// Ref: gram.y CreateSchemaStmt
//
//	CreateSchemaStmt:
//	    CREATE SCHEMA opt_single_name AUTHORIZATION RoleSpec OptSchemaEltList
//	    | CREATE SCHEMA ColId OptSchemaEltList
//	    | CREATE SCHEMA IF_P NOT EXISTS opt_single_name AUTHORIZATION RoleSpec OptSchemaEltList
//	    | CREATE SCHEMA IF_P NOT EXISTS ColId OptSchemaEltList
//
// The CREATE keyword has already been consumed. Current token is SCHEMA.
func (p *Parser) parseCreateSchemaStmt() (nodes.Node, error) {
	p.advance() // consume SCHEMA

	stmt := &nodes.CreateSchemaStmt{}

	// Check for IF NOT EXISTS
	if p.cur.Type == IF_P {
		p.advance() // consume IF
		if _, err := p.expect(NOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	if p.cur.Type == AUTHORIZATION {
		p.advance() // consume AUTHORIZATION
		stmt.Authrole = p.parseRoleSpec()
		elts, err := p.parseOptSchemaEltList()
		if err != nil {
			return nil, err
		}
		stmt.SchemaElts = elts
		return stmt, nil
	}

	name, err := p.parseColId()
	if err != nil {
		return nil, err
	}
	stmt.Schemaname = name

	if p.cur.Type == AUTHORIZATION {
		p.advance() // consume AUTHORIZATION
		stmt.Authrole = p.parseRoleSpec()
	}

	elts, err := p.parseOptSchemaEltList()
	if err != nil {
		return nil, err
	}
	stmt.SchemaElts = elts
	return stmt, nil
}

// parseOptSchemaEltList parses an optional list of schema elements.
func (p *Parser) parseOptSchemaEltList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		elt, err := p.parseSchemaStmt()
		if err != nil {
			return nil, err
		}
		if elt == nil {
			break
		}
		items = append(items, elt)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseSchemaStmt parses a schema_stmt (a statement allowed inside CREATE SCHEMA).
func (p *Parser) parseSchemaStmt() (nodes.Node, error) {
	if p.cur.Type != CREATE {
		return nil, nil
	}
	next := p.peekNext()
	switch next.Type {
	case TABLE:
		return p.parseCreateOrCTAS()
	case INDEX, UNIQUE:
		loc := p.pos()
		p.advance() // consume CREATE
		stmt, err := p.parseIndexStmt()
		if stmt != nil { stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End} }
		return stmt, err
	case SEQUENCE:
		p.advance() // consume CREATE
		return p.parseCreateSeqStmt(byte(nodes.RELPERSISTENCE_PERMANENT))
	case VIEW:
		loc := p.pos()
		p.advance() // consume CREATE
		stmt, err := p.parseViewStmt(false)
		if stmt != nil { stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End} }
		return stmt, err
	case OR:
		loc := p.pos()
		p.advance() // consume CREATE
		p.advance() // consume OR
		if _, err := p.expect(REPLACE); err != nil {
			return nil, err
		}
		stmt, err := p.parseViewStmt(true)
		if stmt != nil { stmt.Loc = nodes.Loc{Start: loc, End: p.prev.End} }
		return stmt, err
	default:
		return nil, nil
	}
}

// parseCreateTableSpaceStmt parses a CREATE TABLESPACE statement.
// The CREATE keyword has already been consumed. Current token is TABLESPACE.
//
// Ref: https://www.postgresql.org/docs/17/sql-createtablespace.html
//
//	CREATE TABLESPACE tablespace_name
//	    [ OWNER { new_owner | CURRENT_ROLE | CURRENT_USER | SESSION_USER } ]
//	    LOCATION 'directory'
//	    [ WITH ( tablespace_option = value [, ... ] ) ]
func (p *Parser) parseCreateTableSpaceStmt() (nodes.Node, error) {
	p.advance() // consume TABLESPACE
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	var owner *nodes.RoleSpec
	if p.cur.Type == OWNER {
		p.advance()
		owner = p.parseRoleSpec()
	}

	if _, err := p.expect(LOCATION); err != nil {
		return nil, err
	}
	loc := p.cur.Str
	if _, err := p.expect(SCONST); err != nil {
		return nil, err
	}

	var options *nodes.List
	if p.cur.Type == WITH {
		p.advance()
		options = p.parseReloptions()
	}

	return &nodes.CreateTableSpaceStmt{
		Tablespacename: name,
		Owner:          owner,
		Location:       loc,
		Options:        options,
	}, nil
}

// parseCommentStmt parses a COMMENT ON statement.
// COMMENT has already been consumed. Current token is ON.
func (p *Parser) parseCommentStmt() (nodes.Node, error) {
	if _, err := p.expect(ON); err != nil {
		return nil, err
	}

	stmt := &nodes.CommentStmt{}

	if p.collectMode() {
		p.cachedCollect("parseCommentStmt", func() {
			for _, t := range []int{COLUMN, TYPE_P, TABLE, SCHEMA, INDEX, SEQUENCE, FUNCTION, PROCEDURE, TRIGGER, VIEW, MATERIALIZED, EXTENSION, POLICY, DOMAIN_P} {
				p.addTokenCandidate(t)
			}
		})
		return nil, nil
	}

	switch p.cur.Type {
	case COLUMN:
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("qualified_name")
			return nil, nil
		}
		stmt.Objtype = nodes.OBJECT_COLUMN
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case TYPE_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_TYPE
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case DOMAIN_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_DOMAIN
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case CONSTRAINT:
		p.advance()
		constraintName, err := p.parseName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(ON); err != nil {
			return nil, err
		}
		if p.cur.Type == DOMAIN_P {
			p.advance()
			stmt.Objtype = nodes.OBJECT_DOMCONSTRAINT
			anyName, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			stmt.Object = &nodes.List{Items: []nodes.Node{makeTypeNameFromNameList(anyName), &nodes.String{Str: constraintName}}}
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}
		stmt.Objtype = nodes.OBJECT_TABCONSTRAINT
		anyName, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = appendNodeToList(anyName, &nodes.String{Str: constraintName})
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case LARGE_P:
		p.advance()
		if _, err := p.expect(OBJECT_P); err != nil {
			return nil, err
		}
		stmt.Objtype = nodes.OBJECT_LARGEOBJECT
		stmt.Object = p.parseNumericOnly()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case CAST:
		p.advance()
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
		dstType, err := p.parseTypename()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.Objtype = nodes.OBJECT_CAST
		stmt.Object = &nodes.List{Items: []nodes.Node{srcType, dstType}}
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case EVENT:
		p.advance()
		if _, err := p.expect(TRIGGER); err != nil {
			return nil, err
		}
		stmt.Objtype = nodes.OBJECT_EVENT_TRIGGER
		name, err := p.parseName()
		if err != nil {
			return nil, err
		}
		stmt.Object = &nodes.String{Str: name}
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case TRANSFORM:
		p.advance()
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
		stmt.Objtype = nodes.OBJECT_TRANSFORM
		stmt.Object = &nodes.List{Items: []nodes.Node{typeName, &nodes.String{Str: lang}}}
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case OPERATOR:
		p.advance()
		if p.cur.Type == CLASS {
			p.advance()
			stmt.Objtype = nodes.OBJECT_OPCLASS
			name, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(USING); err != nil {
				return nil, err
			}
			accessMethod, err := p.parseName()
			if err != nil {
				return nil, err
			}
			stmt.Object = prependNodeToList(&nodes.String{Str: accessMethod}, name)
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}
		if p.cur.Type == FAMILY {
			p.advance()
			stmt.Objtype = nodes.OBJECT_OPFAMILY
			name, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(USING); err != nil {
				return nil, err
			}
			accessMethod, err := p.parseName()
			if err != nil {
				return nil, err
			}
			stmt.Object = prependNodeToList(&nodes.String{Str: accessMethod}, name)
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}
		stmt.Objtype = nodes.OBJECT_OPERATOR
		stmt.Object = p.parseOperatorWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case AGGREGATE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_AGGREGATE
		stmt.Object = p.parseAggregateWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case FUNCTION:
		p.advance()
		stmt.Objtype = nodes.OBJECT_FUNCTION
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case PROCEDURE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_PROCEDURE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case ROUTINE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_ROUTINE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case FOREIGN:
		p.advance()
		if p.cur.Type == DATA_P {
			p.advance()
			if _, err := p.expect(WRAPPER); err != nil {
				return nil, err
			}
			stmt.Objtype = nodes.OBJECT_FDW
			name, err := p.parseName()
			if err != nil {
				return nil, err
			}
			stmt.Object = &nodes.String{Str: name}
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}
		if _, err := p.expect(TABLE); err != nil {
			return nil, err
		}
		stmt.Objtype = nodes.OBJECT_FOREIGN_TABLE
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	case MATERIALIZED:
		p.advance()
		if _, err := p.expect(VIEW); err != nil {
			return nil, err
		}
		stmt.Objtype = nodes.OBJECT_MATVIEW
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil

	default:
		// Try object_type_name_on_any_name: POLICY | RULE | TRIGGER
		if objType, ok := p.tryParseObjectTypeNameOnAnyName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			objName, err := p.parseName()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(ON); err != nil {
				return nil, err
			}
			anyName, err := p.parseAnyName()
			if err != nil {
				return nil, err
			}
			stmt.Object = appendNodeToList(anyName, &nodes.String{Str: objName})
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}

		// Try object_type_name: SCHEMA | DATABASE | ROLE | TABLESPACE | ...
		if objType, ok := p.tryParseObjectTypeName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			name, err := p.parseName()
			if err != nil {
				return nil, err
			}
			stmt.Object = &nodes.String{Str: name}
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Comment = p.parseCommentText()
			return stmt, nil
		}

		// Remaining object_type_any_name
		stmt.Objtype = p.parseObjectTypeAnyName()
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Comment = p.parseCommentText()
		return stmt, nil
	}
}

// parseCommentText parses comment_text: Sconst | NULL_P
func (p *Parser) parseCommentText() string {
	if p.cur.Type == NULL_P {
		p.advance()
		return ""
	}
	if p.cur.Type == SCONST {
		tok := p.advance()
		return tok.Str
	}
	return ""
}

// tryParseObjectTypeName tries to parse object_type_name.
func (p *Parser) tryParseObjectTypeName() (int64, bool) {
	switch p.cur.Type {
	case SCHEMA:
		p.advance()
		return int64(nodes.OBJECT_SCHEMA), true
	case DATABASE:
		p.advance()
		return int64(nodes.OBJECT_DATABASE), true
	case ROLE:
		p.advance()
		return int64(nodes.OBJECT_ROLE), true
	case TABLESPACE:
		p.advance()
		return int64(nodes.OBJECT_TABLESPACE), true
	case SUBSCRIPTION:
		p.advance()
		return int64(nodes.OBJECT_SUBSCRIPTION), true
	case PUBLICATION:
		p.advance()
		return int64(nodes.OBJECT_PUBLICATION), true
	case SERVER:
		p.advance()
		return int64(nodes.OBJECT_FOREIGN_SERVER), true
	}
	return 0, false
}

// tryParseObjectTypeNameOnAnyName tries to parse object_type_name_on_any_name.
func (p *Parser) tryParseObjectTypeNameOnAnyName() (int64, bool) {
	switch p.cur.Type {
	case POLICY:
		p.advance()
		return int64(nodes.OBJECT_POLICY), true
	case RULE:
		p.advance()
		return int64(nodes.OBJECT_RULE), true
	case TRIGGER:
		p.advance()
		return int64(nodes.OBJECT_TRIGGER), true
	}
	return 0, false
}

// prependNodeToList prepends a node to a List.
func prependNodeToList(n nodes.Node, list *nodes.List) *nodes.List {
	if list == nil {
		return &nodes.List{Items: []nodes.Node{n}}
	}
	items := make([]nodes.Node, 0, len(list.Items)+1)
	items = append(items, n)
	items = append(items, list.Items...)
	return &nodes.List{Items: items}
}

// appendNodeToList appends a node to a List.
func appendNodeToList(list *nodes.List, n nodes.Node) *nodes.List {
	if list == nil {
		return &nodes.List{Items: []nodes.Node{n}}
	}
	list.Items = append(list.Items, n)
	return list
}

// parseFunctionWithArgtypesForComment parses function_with_argtypes for COMMENT ON.
func (p *Parser) parseFunctionWithArgtypesForComment() nodes.Node {
	funcName, _ := p.parseFuncName()
	owa := &nodes.ObjectWithArgs{Objname: funcName}
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == ')' {
			p.advance()
		} else {
			owa.Objargs = p.parseCommentFuncArgsList()
			p.expect(')')
		}
	} else {
		owa.ArgsUnspecified = true
	}
	return owa
}

// parseCommentFuncArgsList parses a comma-separated list of function argument types.
func (p *Parser) parseCommentFuncArgsList() *nodes.List {
	var items []nodes.Node
	for {
		typeName, _ := p.parseTypename()
		if typeName != nil {
			items = append(items, typeName)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseAggregateWithArgtypesForComment parses aggregate_with_argtypes.
func (p *Parser) parseAggregateWithArgtypesForComment() nodes.Node {
	funcName, _ := p.parseFuncName()
	owa := &nodes.ObjectWithArgs{Objname: funcName}
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == '*' {
			p.advance()
			p.expect(')')
			owa.Objargs = &nodes.List{Items: []nodes.Node{&nodes.String{Str: "*"}}}
		} else if p.cur.Type == ')' {
			p.advance()
		} else {
			owa.Objargs = p.parseCommentFuncArgsList()
			p.expect(')')
		}
	} else {
		owa.ArgsUnspecified = true
	}
	return owa
}

// parseOperatorWithArgtypesForComment parses operator_with_argtypes.
func (p *Parser) parseOperatorWithArgtypesForComment() nodes.Node {
	opName, _ := p.parseAnyOperator()
	owa := &nodes.ObjectWithArgs{Objname: opName}
	if p.cur.Type == '(' {
		p.advance()
		var left, right nodes.Node
		if p.cur.Type == NONE {
			p.advance()
			p.expect(',')
			right, _ = p.parseTypename()
		} else {
			left, _ = p.parseTypename()
			if p.cur.Type == ',' {
				p.advance()
				if p.cur.Type == NONE {
					p.advance()
				} else {
					right, _ = p.parseTypename()
				}
			}
		}
		owa.Objargs = &nodes.List{Items: []nodes.Node{left, right}}
		p.expect(')')
	}
	return owa
}

// parseSecLabelStmt parses a SECURITY LABEL statement.
// SECURITY has already been consumed. Current token is LABEL.
func (p *Parser) parseSecLabelStmt() (nodes.Node, error) {
	if _, err := p.expect(LABEL); err != nil {
		return nil, err
	}

	provider := ""
	if p.cur.Type == FOR {
		p.advance()
		provider = p.parseNonReservedWordOrSconst()
	}

	if _, err := p.expect(ON); err != nil {
		return nil, err
	}

	stmt := &nodes.SecLabelStmt{Provider: provider}

	switch p.cur.Type {
	case COLUMN:
		p.advance()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			p.addRuleCandidate("qualified_name")
			return nil, nil
		}
		stmt.Objtype = nodes.OBJECT_COLUMN
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case TYPE_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_TYPE
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case DOMAIN_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_DOMAIN
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case AGGREGATE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_AGGREGATE
		stmt.Object = p.parseAggregateWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case FUNCTION:
		p.advance()
		stmt.Objtype = nodes.OBJECT_FUNCTION
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case PROCEDURE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_PROCEDURE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	case ROUTINE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_ROUTINE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil

	default:
		if objType, ok := p.tryParseObjectTypeName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			name, err := p.parseName()
			if err != nil {
				return nil, err
			}
			stmt.Object = &nodes.String{Str: name}
			if _, err := p.expect(IS); err != nil {
				return nil, err
			}
			stmt.Label = p.parseSecurityLabel()
			return stmt, nil
		}

		stmt.Objtype = p.parseObjectTypeAnyName()
		name, err := p.parseAnyName()
		if err != nil {
			return nil, err
		}
		stmt.Object = name
		if _, err := p.expect(IS); err != nil {
			return nil, err
		}
		stmt.Label = p.parseSecurityLabel()
		return stmt, nil
	}
}

// parseSecurityLabel parses security_label: Sconst | NULL_P
func (p *Parser) parseSecurityLabel() string {
	if p.cur.Type == NULL_P {
		p.advance()
		return ""
	}
	if p.cur.Type == SCONST {
		tok := p.advance()
		return tok.Str
	}
	return ""
}
