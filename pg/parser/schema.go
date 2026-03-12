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
func (p *Parser) parseCreateSchemaStmt() nodes.Node {
	p.advance() // consume SCHEMA

	stmt := &nodes.CreateSchemaStmt{}

	// Check for IF NOT EXISTS
	if p.cur.Type == IF_P {
		p.advance() // consume IF
		p.expect(NOT)
		p.expect(EXISTS)
		stmt.IfNotExists = true
	}

	if p.cur.Type == AUTHORIZATION {
		p.advance() // consume AUTHORIZATION
		stmt.Authrole = p.parseRoleSpec()
		stmt.SchemaElts = p.parseOptSchemaEltList()
		return stmt
	}

	name, _ := p.parseColId()
	stmt.Schemaname = name

	if p.cur.Type == AUTHORIZATION {
		p.advance() // consume AUTHORIZATION
		stmt.Authrole = p.parseRoleSpec()
	}

	stmt.SchemaElts = p.parseOptSchemaEltList()
	return stmt
}

// parseOptSchemaEltList parses an optional list of schema elements.
func (p *Parser) parseOptSchemaEltList() *nodes.List {
	var items []nodes.Node
	for {
		elt := p.parseSchemaStmt()
		if elt == nil {
			break
		}
		items = append(items, elt)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseSchemaStmt parses a schema_stmt (a statement allowed inside CREATE SCHEMA).
func (p *Parser) parseSchemaStmt() nodes.Node {
	if p.cur.Type != CREATE {
		return nil
	}
	next := p.peekNext()
	switch next.Type {
	case TABLE:
		return p.parseCreateOrCTAS()
	case INDEX, UNIQUE:
		p.advance() // consume CREATE
		return p.parseIndexStmt()
	case SEQUENCE:
		p.advance() // consume CREATE
		return p.parseCreateSeqStmt(byte(nodes.RELPERSISTENCE_PERMANENT))
	case VIEW:
		p.advance() // consume CREATE
		return p.parseViewStmt(false)
	case OR:
		p.advance() // consume CREATE
		p.advance() // consume OR
		p.expect(REPLACE)
		return p.parseViewStmt(true)
	default:
		return nil
	}
}

// parseCommentStmt parses a COMMENT ON statement.
// COMMENT has already been consumed. Current token is ON.
func (p *Parser) parseCommentStmt() nodes.Node {
	p.expect(ON)

	stmt := &nodes.CommentStmt{}

	switch p.cur.Type {
	case COLUMN:
		p.advance()
		stmt.Objtype = nodes.OBJECT_COLUMN
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case TYPE_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_TYPE
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case DOMAIN_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_DOMAIN
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case CONSTRAINT:
		p.advance()
		constraintName, _ := p.parseName()
		p.expect(ON)
		if p.cur.Type == DOMAIN_P {
			p.advance()
			stmt.Objtype = nodes.OBJECT_DOMCONSTRAINT
			anyName, _ := p.parseAnyName()
			stmt.Object = &nodes.List{Items: []nodes.Node{makeTypeNameFromNameList(anyName), &nodes.String{Str: constraintName}}}
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}
		stmt.Objtype = nodes.OBJECT_TABCONSTRAINT
		anyName, _ := p.parseAnyName()
		stmt.Object = appendNodeToList(anyName, &nodes.String{Str: constraintName})
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case LARGE_P:
		p.advance()
		p.expect(OBJECT_P)
		stmt.Objtype = nodes.OBJECT_LARGEOBJECT
		stmt.Object = p.parseNumericOnly()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case CAST:
		p.advance()
		p.expect('(')
		srcType, _ := p.parseTypename()
		p.expect(AS)
		dstType, _ := p.parseTypename()
		p.expect(')')
		stmt.Objtype = nodes.OBJECT_CAST
		stmt.Object = &nodes.List{Items: []nodes.Node{srcType, dstType}}
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case EVENT:
		p.advance()
		p.expect(TRIGGER)
		stmt.Objtype = nodes.OBJECT_EVENT_TRIGGER
		name, _ := p.parseName()
		stmt.Object = &nodes.String{Str: name}
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case TRANSFORM:
		p.advance()
		p.expect(FOR)
		typeName, _ := p.parseTypename()
		p.expect(LANGUAGE)
		lang, _ := p.parseName()
		stmt.Objtype = nodes.OBJECT_TRANSFORM
		stmt.Object = &nodes.List{Items: []nodes.Node{typeName, &nodes.String{Str: lang}}}
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case OPERATOR:
		p.advance()
		if p.cur.Type == CLASS {
			p.advance()
			stmt.Objtype = nodes.OBJECT_OPCLASS
			name, _ := p.parseAnyName()
			p.expect(USING)
			accessMethod, _ := p.parseName()
			stmt.Object = prependNodeToList(&nodes.String{Str: accessMethod}, name)
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}
		if p.cur.Type == FAMILY {
			p.advance()
			stmt.Objtype = nodes.OBJECT_OPFAMILY
			name, _ := p.parseAnyName()
			p.expect(USING)
			accessMethod, _ := p.parseName()
			stmt.Object = prependNodeToList(&nodes.String{Str: accessMethod}, name)
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}
		stmt.Objtype = nodes.OBJECT_OPERATOR
		stmt.Object = p.parseOperatorWithArgtypesForComment()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case AGGREGATE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_AGGREGATE
		stmt.Object = p.parseAggregateWithArgtypesForComment()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case FUNCTION:
		p.advance()
		stmt.Objtype = nodes.OBJECT_FUNCTION
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case PROCEDURE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_PROCEDURE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case ROUTINE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_ROUTINE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case FOREIGN:
		p.advance()
		if p.cur.Type == DATA_P {
			p.advance()
			p.expect(WRAPPER)
			stmt.Objtype = nodes.OBJECT_FDW
			name, _ := p.parseName()
			stmt.Object = &nodes.String{Str: name}
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}
		p.expect(TABLE)
		stmt.Objtype = nodes.OBJECT_FOREIGN_TABLE
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	case MATERIALIZED:
		p.advance()
		p.expect(VIEW)
		stmt.Objtype = nodes.OBJECT_MATVIEW
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt

	default:
		// Try object_type_name_on_any_name: POLICY | RULE | TRIGGER
		if objType, ok := p.tryParseObjectTypeNameOnAnyName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			objName, _ := p.parseName()
			p.expect(ON)
			anyName, _ := p.parseAnyName()
			stmt.Object = appendNodeToList(anyName, &nodes.String{Str: objName})
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}

		// Try object_type_name: SCHEMA | DATABASE | ROLE | TABLESPACE | ...
		if objType, ok := p.tryParseObjectTypeName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			name, _ := p.parseName()
			stmt.Object = &nodes.String{Str: name}
			p.expect(IS)
			stmt.Comment = p.parseCommentText()
			return stmt
		}

		// Remaining object_type_any_name
		stmt.Objtype = p.parseObjectTypeAnyName()
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Comment = p.parseCommentText()
		return stmt
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
func (p *Parser) parseSecLabelStmt() nodes.Node {
	p.expect(LABEL)

	provider := ""
	if p.cur.Type == FOR {
		p.advance()
		provider = p.parseNonReservedWordOrSconst()
	}

	p.expect(ON)

	stmt := &nodes.SecLabelStmt{Provider: provider}

	switch p.cur.Type {
	case COLUMN:
		p.advance()
		stmt.Objtype = nodes.OBJECT_COLUMN
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case TYPE_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_TYPE
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case DOMAIN_P:
		p.advance()
		stmt.Objtype = nodes.OBJECT_DOMAIN
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case AGGREGATE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_AGGREGATE
		stmt.Object = p.parseAggregateWithArgtypesForComment()
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case FUNCTION:
		p.advance()
		stmt.Objtype = nodes.OBJECT_FUNCTION
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case PROCEDURE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_PROCEDURE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	case ROUTINE:
		p.advance()
		stmt.Objtype = nodes.OBJECT_ROUTINE
		stmt.Object = p.parseFunctionWithArgtypesForComment()
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt

	default:
		if objType, ok := p.tryParseObjectTypeName(); ok {
			stmt.Objtype = nodes.ObjectType(objType)
			name, _ := p.parseName()
			stmt.Object = &nodes.String{Str: name}
			p.expect(IS)
			stmt.Label = p.parseSecurityLabel()
			return stmt
		}

		stmt.Objtype = p.parseObjectTypeAnyName()
		name, _ := p.parseAnyName()
		stmt.Object = name
		p.expect(IS)
		stmt.Label = p.parseSecurityLabel()
		return stmt
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
