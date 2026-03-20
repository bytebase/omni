package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreateFdwStmt parses CREATE FOREIGN DATA WRAPPER.
// Already consumed: CREATE FOREIGN
func (p *Parser) parseCreateFdwStmt() (nodes.Node, error) {
	if _, err := p.expect(DATA_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return nil, err
	}
	name, _ := p.parseName()
	funcOptions := p.parseOptFdwOptions()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
		Options:     options,
	}, nil
}

// parseAlterFdwStmt parses ALTER FOREIGN DATA WRAPPER.
// Already consumed: ALTER FOREIGN
func (p *Parser) parseAlterFdwStmt() (nodes.Node, error) {
	if _, err := p.expect(DATA_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return nil, err
	}
	name, _ := p.parseName()

	// Check for OWNER TO and RENAME TO before FDW options
	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_FDW,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_FDW,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	}

	funcOptions := p.parseOptFdwOptions()
	if p.cur.Type == OPTIONS {
		options := p.parseAlterGenericOptions()
		return &nodes.AlterFdwStmt{
			Fdwname:     name,
			FuncOptions: funcOptions,
			Options:     options,
		}, nil
	}
	return &nodes.AlterFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
	}, nil
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

func (p *Parser) parseCreateForeignServerStmt() (nodes.Node, error) {
	p.advance() // consume SERVER
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
	servertype := p.parseOptType()
	version := p.parseOptForeignServerVersion()
	if _, err := p.expect(FOREIGN); err != nil {
		return nil, err
	}
	if _, err := p.expect(DATA_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return nil, err
	}
	fdwname, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateForeignServerStmt{
		Servername:  name,
		Servertype:  servertype,
		Version:     version,
		Fdwname:     fdwname,
		IfNotExists: ifNotExists,
		Options:     options,
	}, nil
}

func (p *Parser) parseAlterForeignServerStmt() (nodes.Node, error) {
	p.advance() // consume SERVER
	name, _ := p.parseName()

	// Check for OWNER TO and RENAME TO before VERSION/OPTIONS
	switch p.cur.Type {
	case OWNER:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_FOREIGN_SERVER,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil
	case RENAME:
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, _ := p.parseName()
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_FOREIGN_SERVER,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	}

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
	}, nil
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

func (p *Parser) parseCreateForeignTableStmt() (nodes.Node, error) {
	if _, err := p.expect(TABLE); err != nil {
		return nil, err
	}
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
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names)
	rv.Relpersistence = 'p'

	if p.cur.Type == PARTITION {
		p.advance()
		if _, err := p.expect(OF); err != nil {
			return nil, err
		}
		inhNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		inhRv := makeRangeVarFromNames(inhNames)
		var tableElts *nodes.List
		if p.cur.Type == '(' {
			p.advance()
			if p.cur.Type != ')' {
				var err error
				tableElts, err = p.parseOptTableElementList()
				if err != nil {
					return nil, err
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		partbound := p.parseForValues()
		if _, err := p.expect(SERVER); err != nil {
			return nil, err
		}
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
		}, nil
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var tableElts *nodes.List
	if p.cur.Type != ')' {
		var err error
		tableElts, err = p.parseOptTableElementList()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	inhRelations := p.parseOptInherit()
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
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
	}, nil
}

func (p *Parser) parseImportForeignSchemaStmt() (nodes.Node, error) {
	if _, err := p.expect(FOREIGN); err != nil {
		return nil, err
	}
	if _, err := p.expect(SCHEMA); err != nil {
		return nil, err
	}
	remoteSchema, _ := p.parseName()
	listType := nodes.FDW_IMPORT_SCHEMA_ALL
	var tableList *nodes.List
	if p.cur.Type == LIMIT {
		p.advance()
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		listType = nodes.FDW_IMPORT_SCHEMA_LIMIT_TO
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		tableList = p.parseRelationExprList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == EXCEPT {
		p.advance()
		listType = nodes.FDW_IMPORT_SCHEMA_EXCEPT
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		tableList = p.parseRelationExprList()
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(FROM); err != nil {
		return nil, err
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	serverName, _ := p.parseName()
	if _, err := p.expect(INTO); err != nil {
		return nil, err
	}
	localSchema, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.ImportForeignSchemaStmt{
		ServerName:   serverName,
		RemoteSchema: remoteSchema,
		LocalSchema:  localSchema,
		ListType:     listType,
		TableList:    tableList,
		Options:      options,
	}, nil
}

func (p *Parser) parseCreateUserMappingIfNotExistsStmt() (nodes.Node, error) {
	p.advance() // consume USER
	if _, err := p.expect(MAPPING); err != nil {
		return nil, err
	}
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
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	user := p.parseAuthIdent()
	if user == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, _ := p.parseName()
	options := p.parseCreateGenericOptions()
	return &nodes.CreateUserMappingStmt{
		User:        user,
		Servername:  servername,
		Options:     options,
		IfNotExists: ifNotExists,
	}, nil
}

func (p *Parser) parseAlterUserMappingStmt() (nodes.Node, error) {
	p.advance() // consume USER
	if _, err := p.expect(MAPPING); err != nil {
		return nil, err
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	user := p.parseAuthIdent()
	if user == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, _ := p.parseName()
	options := p.parseAlterGenericOptions()
	return &nodes.AlterUserMappingStmt{
		User:       user,
		Servername: servername,
		Options:    options,
	}, nil
}

func (p *Parser) parseDropUserMappingStmt() (nodes.Node, error) {
	p.advance() // consume USER
	if _, err := p.expect(MAPPING); err != nil {
		return nil, err
	}
	missingOk := false
	if p.cur.Type == IF_P {
		p.advance()
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}
	if _, err := p.expect(FOR); err != nil {
		return nil, err
	}
	user := p.parseAuthIdent()
	if user == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, _ := p.parseName()
	return &nodes.DropUserMappingStmt{
		User:       user,
		Servername: servername,
		MissingOk:  missingOk,
	}, nil
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

func (p *Parser) parseCreatePLangStmt(replace bool) (nodes.Node, error) {
	trusted := false
	if p.cur.Type == TRUSTED {
		p.advance()
		trusted = true
	}
	if p.cur.Type == PROCEDURAL {
		p.advance()
	}
	if _, err := p.expect(LANGUAGE); err != nil {
		return nil, err
	}
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
	return stmt, nil
}
