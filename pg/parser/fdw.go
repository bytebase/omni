package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateFdwStmt parses CREATE FOREIGN DATA WRAPPER.
// Already consumed: CREATE FOREIGN
func (p *Parser) parseCreateFdwStmt() nodes.Node {
	p.expect(DATA_P)
	p.expect(WRAPPER)
	name, _ := p.parseName()
	funcOptions := p.parseOptFdwOptions()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
		Options:     options,
	}
}

// parseAlterFdwStmt parses ALTER FOREIGN DATA WRAPPER.
// Already consumed: ALTER FOREIGN
func (p *Parser) parseAlterFdwStmt() nodes.Node {
	p.expect(DATA_P)
	p.expect(WRAPPER)
	name, _ := p.parseName()
	funcOptions := p.parseOptFdwOptions()
	if p.cur.Type == OPTIONS {
		options := p.parseAlterGenericOptions()
		return &nodes.AlterFdwStmt{
			Fdwname:     name,
			FuncOptions: funcOptions,
			Options:     options,
		}
	}
	return &nodes.AlterFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
	}
}

func (p *Parser) parseFdwOption() *nodes.DefElem {
	switch p.cur.Type {
	case HANDLER:
		p.advance()
		name := p.parseHandlerName()
		return makeDefElem("handler", name)
	case VALIDATOR:
		p.advance()
		name := p.parseHandlerName()
		return makeDefElem("validator", name)
	case NO:
		p.advance()
		switch p.cur.Type {
		case HANDLER:
			p.advance()
			return makeDefElem("handler", nil)
		case VALIDATOR:
			p.advance()
			return makeDefElem("validator", nil)
		}
	}
	return nil
}

func (p *Parser) parseOptFdwOptions() *nodes.List {
	var items []nodes.Node
	for {
		if p.cur.Type != HANDLER && p.cur.Type != VALIDATOR && p.cur.Type != NO {
			break
		}
		if p.cur.Type == NO {
			next := p.peekNext()
			if next.Type != HANDLER && next.Type != VALIDATOR {
				break
			}
		}
		opt := p.parseFdwOption()
		if opt == nil {
			break
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

func (p *Parser) parseCreateGenericOptions() *nodes.List {
	if p.cur.Type != OPTIONS {
		return nil
	}
	p.advance()
	p.expect('(')
	list := p.parseGenericOptionList()
	p.expect(')')
	return list
}

func (p *Parser) parseGenericOptionList() *nodes.List {
	elem := p.parseGenericOptionElem()
	items := []nodes.Node{elem}
	for p.cur.Type == ',' {
		p.advance()
		elem = p.parseGenericOptionElem()
		items = append(items, elem)
	}
	return &nodes.List{Items: items}
}

// parseHandlerName parses handler_name: name | name attrs
func (p *Parser) parseHandlerName() *nodes.List {
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

func (p *Parser) parseCreateForeignServerStmt() nodes.Node {
	p.advance() // consume SERVER
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}
	name, _ := p.parseName()
	servertype := p.parseOptType()
	version := p.parseOptForeignServerVersion()
	p.expect(FOREIGN)
	p.expect(DATA_P)
	p.expect(WRAPPER)
	fdwname, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateForeignServerStmt{
		Servername:  name,
		Servertype:  servertype,
		Version:     version,
		Fdwname:     fdwname,
		IfNotExists: ifNotExists,
		Options:     options,
	}
}

func (p *Parser) parseAlterForeignServerStmt() nodes.Node {
	p.advance() // consume SERVER
	name, _ := p.parseName()
	hasVersion := false
	version := ""
	if p.cur.Type == VERSION_P {
		version = p.parseForeignServerVersion()
		hasVersion = true
	}
	var options *nodes.List
	if p.cur.Type == OPTIONS {
		options = p.parseAlterGenericOptions()
	}
	return &nodes.AlterForeignServerStmt{
		Servername: name,
		Version:    version,
		Options:    options,
		HasVersion: hasVersion,
	}
}

func (p *Parser) parseOptType() string {
	if p.cur.Type == TYPE_P {
		p.advance()
		val := p.cur.Str
		p.expect(SCONST)
		return val
	}
	return ""
}

func (p *Parser) parseForeignServerVersion() string {
	p.expect(VERSION_P)
	if p.cur.Type == NULL_P {
		p.advance()
		return ""
	}
	val := p.cur.Str
	p.expect(SCONST)
	return val
}

func (p *Parser) parseOptForeignServerVersion() string {
	if p.cur.Type == VERSION_P {
		return p.parseForeignServerVersion()
	}
	return ""
}

func (p *Parser) parseCreateForeignTableStmt() nodes.Node {
	p.expect(TABLE)
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil
	}
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = 'p'

	if p.cur.Type == PARTITION {
		p.advance()
		p.expect(OF)
		inhNames, err := p.parseQualifiedName()
		if err != nil {
			return nil
		}
		inhRv := makeRangeVarFromNames(inhNames)
		var tableElts *nodes.List
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type != ')' {
				tableElts = p.parseOptTableElementList()
			}
			p.expect(')')
		}
		partbound := p.parseForValues()
		p.expect(SERVER)
		servername, _ := p.parseName()
		options := p.parseCreateGenericOptions()
		return &nodes.CreateForeignTableStmt{
			Base: nodes.CreateStmt{
				Relation:     rv,
				InhRelations: &nodes.List{Items: []nodes.Node{inhRv}},
				TableElts:    tableElts,
				Partbound:    partbound,
				IfNotExists:  ifNotExists,
			},
			Servername: servername,
			Options:    options,
		}
	}

	p.expect('(')
	var tableElts *nodes.List
	if p.cur.Type != ')' {
		tableElts = p.parseOptTableElementList()
	}
	p.expect(')')
	inhRelations := p.parseOptInherit()
	p.expect(SERVER)
	servername, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateForeignTableStmt{
		Base: nodes.CreateStmt{
			Relation:     rv,
			TableElts:    tableElts,
			InhRelations: inhRelations,
			IfNotExists:  ifNotExists,
		},
		Servername: servername,
		Options:    options,
	}
}

func (p *Parser) parseImportForeignSchemaStmt() nodes.Node {
	p.expect(FOREIGN)
	p.expect(SCHEMA)
	remoteSchema, _ := p.parseName()
	listType := nodes.FDW_IMPORT_SCHEMA_ALL
	var tableList *nodes.List
	if p.cur.Type == LIMIT {
		p.advance()
		p.expect(TO)
		listType = nodes.FDW_IMPORT_SCHEMA_LIMIT_TO
		p.expect('(')
		tableList = p.parseRelationExprList()
		p.expect(')')
	} else if p.cur.Type == EXCEPT {
		p.advance()
		listType = nodes.FDW_IMPORT_SCHEMA_EXCEPT
		p.expect('(')
		tableList = p.parseRelationExprList()
		p.expect(')')
	}
	p.expect(FROM)
	p.expect(SERVER)
	serverName, _ := p.parseName()
	p.expect(INTO)
	localSchema, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.ImportForeignSchemaStmt{
		ServerName:   serverName,
		RemoteSchema: remoteSchema,
		LocalSchema:  localSchema,
		ListType:     listType,
		TableList:    tableList,
		Options:      options,
	}
}

func (p *Parser) parseCreateUserMappingIfNotExistsStmt() nodes.Node {
	p.advance() // consume USER
	p.expect(MAPPING)
	ifNotExists := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(NOT)
		p.expect(EXISTS)
		ifNotExists = true
	}
	p.expect(FOR)
	user := p.parseAuthIdent()
	if user == nil {
		return nil
	}
	p.expect(SERVER)
	servername, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateUserMappingStmt{
		User:        user,
		Servername:  servername,
		Options:     options,
		IfNotExists: ifNotExists,
	}
}

func (p *Parser) parseAlterUserMappingStmt() nodes.Node {
	p.advance() // consume USER
	p.expect(MAPPING)
	p.expect(FOR)
	user := p.parseAuthIdent()
	if user == nil {
		return nil
	}
	p.expect(SERVER)
	servername, _ := p.parseName()
	options := p.parseAlterGenericOptions()
	return &nodes.AlterUserMappingStmt{
		User:       user,
		Servername: servername,
		Options:    options,
	}
}

func (p *Parser) parseDropUserMappingStmt() nodes.Node {
	p.advance() // consume USER
	p.expect(MAPPING)
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		p.expect(EXISTS)
		missingOk = true
	}
	p.expect(FOR)
	user := p.parseAuthIdent()
	if user == nil {
		return nil
	}
	p.expect(SERVER)
	servername, _ := p.parseName()
	return &nodes.DropUserMappingStmt{
		User:      user,
		Servername: servername,
		MissingOk: missingOk,
	}
}

func (p *Parser) parseAuthIdent() *nodes.RoleSpec {
	if p.cur.Type == USER {
		p.advance()
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CURRENT_USER),
		}
	}
	return p.parseRoleSpec()
}

func (p *Parser) parseCreatePLangStmt(replace bool) nodes.Node {
	trusted := false
	if p.cur.Type == TRUSTED {
		p.advance()
		trusted = true
	}
	if p.cur.Type == PROCEDURAL {
		p.advance()
	}
	p.expect(LANGUAGE)
	name, _ := p.parseName()
	stmt := &nodes.CreatePLangStmt{
		Replace:   replace,
		Plname:    name,
		Pltrusted: trusted,
	}
	if p.cur.Type == HANDLER {
		p.advance()
		stmt.Plhandler = p.parseHandlerName()
		if p.cur.Type == INLINE_P {
			p.advance()
			stmt.Plinline = p.parseHandlerName()
		}
		if p.cur.Type == VALIDATOR {
			p.advance()
			stmt.Plvalidator = p.parseHandlerName()
		} else if p.cur.Type == NO {
			next := p.peekNext()
			if next.Type == VALIDATOR {
				p.advance()
				p.advance()
			}
		}
	}
	return stmt
}
