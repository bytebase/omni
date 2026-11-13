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
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	funcOptions, err := p.parseOptFdwOptions()
	if err != nil {
		return nil, err
	}
	options, err := p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
		Options:     options,
		Loc:         nodes.Loc{Start: p.stmtStart, End: p.prev.End},
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
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Check for OWNER TO and RENAME TO before FDW options
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
			ObjectType: nodes.OBJECT_FDW,
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
			RenameType: nodes.OBJECT_FDW,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	}

	funcOptions, err := p.parseOptFdwOptions()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == OPTIONS {
		options, err := p.parseAlterGenericOptions()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterFdwStmt{
			Fdwname:     name,
			FuncOptions: funcOptions,
			Options:     options,
			Loc:         nodes.Loc{Start: p.stmtStart, End: p.prev.End},
		}, nil
	}
	// AlterFdwStmt: ... name opt_fdw_options alter_generic_options | ... name fdw_options
	if funcOptions == nil {
		if p.collectMode() {
			for _, t := range []int{HANDLER, NO, VALIDATOR, OPTIONS, OWNER, RENAME} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.AlterFdwStmt{
		Fdwname:     name,
		FuncOptions: funcOptions,
		Loc:         nodes.Loc{Start: p.stmtStart, End: p.prev.End},
	}, nil
}

// parseFdwOption parses fdw_option. The caller has seen HANDLER, VALIDATOR
// or NO HANDLER|VALIDATOR.
//
//	fdw_option: HANDLER handler_name | NO HANDLER | VALIDATOR handler_name | NO VALIDATOR
func (p *Parser) parseFdwOption() (*nodes.DefElem, error) {
	switch p.cur.Type {
	case HANDLER, VALIDATOR:
		defname := "handler"
		if p.cur.Type == VALIDATOR {
			defname = "validator"
		}
		p.advance()
		name, err := p.parseHandlerName()
		if err != nil {
			return nil, err
		}
		return makeDefElem(defname, name), nil
	case NO:
		p.advance()
		switch p.cur.Type {
		case HANDLER:
			p.advance()
			return makeDefElem("handler", nil), nil
		case VALIDATOR:
			p.advance()
			return makeDefElem("validator", nil), nil
		}
	}
	return nil, p.syntaxErrorAtCur()
}

// parseOptFdwOptions parses opt_fdw_options: a run of fdw_option, or nothing.
func (p *Parser) parseOptFdwOptions() (*nodes.List, error) {
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
		opt, err := p.parseFdwOption()
		if err != nil {
			return nil, err
		}
		items = append(items, opt)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseCreateGenericOptions parses create_generic_options.
//
//	create_generic_options: OPTIONS '(' generic_option_list ')' | /* EMPTY */
func (p *Parser) parseCreateGenericOptions() (*nodes.List, error) {
	if p.cur.Type != OPTIONS {
		return nil, nil
	}
	p.advance()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	list, err := p.parseGenericOptionList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return list, nil
}

// parseGenericOptionList parses generic_option_list.
//
//	generic_option_list: generic_option_elem | generic_option_list ',' generic_option_elem
func (p *Parser) parseGenericOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		elem, err := p.parseGenericOptionElem()
		if err != nil {
			return nil, err
		}
		items = append(items, elem)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseHandlerName parses handler_name: name | name attrs
func (p *Parser) parseHandlerName() (*nodes.List, error) {
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	items := []nodes.Node{&nodes.String{Str: name}}
	if p.cur.Type == '.' {
		attrs, err := p.parseAttrs()
		if err != nil {
			return nil, err
		}
		items = append(items, attrs.Items...)
	}
	return &nodes.List{Items: items}, nil
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
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	servertype, err := p.parseOptType()
	if err != nil {
		return nil, err
	}
	version, err := p.parseOptForeignServerVersion()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(FOREIGN); err != nil {
		return nil, err
	}
	if _, err := p.expect(DATA_P); err != nil {
		return nil, err
	}
	if _, err := p.expect(WRAPPER); err != nil {
		return nil, err
	}
	fdwname, err := p.parseName()
	if err != nil {
		return nil, err
	}
	options, err := p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateForeignServerStmt{
		Servername:  name,
		Servertype:  servertype,
		Version:     version,
		Fdwname:     fdwname,
		IfNotExists: ifNotExists,
		Options:     options,
		Loc:         nodes.Loc{Start: p.stmtStart, End: p.prev.End},
	}, nil
}

func (p *Parser) parseAlterForeignServerStmt() (nodes.Node, error) {
	p.advance() // consume SERVER
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Check for OWNER TO and RENAME TO before VERSION/OPTIONS
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
			ObjectType: nodes.OBJECT_FOREIGN_SERVER,
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
			RenameType: nodes.OBJECT_FOREIGN_SERVER,
			Object:     &nodes.String{Str: name},
			Newname:    newname,
		}, nil
	}

	hasVersion := false
	version := ""
	if p.cur.Type == VERSION_P {
		if version, err = p.parseForeignServerVersion(); err != nil {
			return nil, err
		}
		hasVersion = true
	}
	var options *nodes.List
	if p.cur.Type == OPTIONS {
		if options, err = p.parseAlterGenericOptions(); err != nil {
			return nil, err
		}
	} else if !hasVersion {
		// AlterForeignServerStmt needs foreign_server_version,
		// alter_generic_options, or both.
		if p.collectMode() {
			for _, t := range []int{VERSION_P, OPTIONS, OWNER, RENAME} {
				p.addTokenCandidate(t)
			}
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.AlterForeignServerStmt{
		Servername: name,
		Version:    version,
		Options:    options,
		HasVersion: hasVersion,
		Loc:        nodes.Loc{Start: p.stmtStart, End: p.prev.End},
	}, nil
}

// parseOptType parses opt_type: TYPE_P Sconst | /*EMPTY*/.
func (p *Parser) parseOptType() (string, error) {
	if p.cur.Type != TYPE_P {
		return "", nil
	}
	p.advance()
	tok, err := p.expect(SCONST)
	if err != nil {
		return "", err
	}
	return tok.Str, nil
}

// parseForeignServerVersion parses foreign_server_version:
// VERSION_P Sconst | VERSION_P NULL_P.
func (p *Parser) parseForeignServerVersion() (string, error) {
	if _, err := p.expect(VERSION_P); err != nil {
		return "", err
	}
	if p.cur.Type == NULL_P {
		p.advance()
		return "", nil
	}
	tok, err := p.expect(SCONST)
	if err != nil {
		return "", err
	}
	return tok.Str, nil
}

func (p *Parser) parseOptForeignServerVersion() (string, error) {
	if p.cur.Type == VERSION_P {
		return p.parseForeignServerVersion()
	}
	return "", nil
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
	nameStart := p.pos()
	names, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	rv := makeRangeVarFromNames(names, nodes.Loc{Start: nameStart, End: p.prev.End})
	rv.Relpersistence = 'p'

	if p.cur.Type == PARTITION {
		p.advance()
		if _, err := p.expect(OF); err != nil {
			return nil, err
		}
		inhStart := p.pos()
		inhNames, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		inhRv := makeRangeVarFromNames(inhNames, nodes.Loc{Start: inhStart, End: p.prev.End})
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
		partbound, err := p.parseForValues()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(SERVER); err != nil {
			return nil, err
		}
		servername, err := p.parseName()
		if err != nil {
			return nil, err
		}
		options, err := p.parseCreateGenericOptions()
		if err != nil {
			return nil, err
		}
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
			Loc:        nodes.Loc{Start: p.stmtStart, End: p.prev.End},
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
	servername, err := p.parseName()
	if err != nil {
		return nil, err
	}
	options, err := p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateForeignTableStmt{
		Base: nodes.CreateStmt{
			Relation:     rv,
			TableElts:    tableElts,
			InhRelations: inhRelations,
			IfNotExists:  ifNotExists,
		},
		Servername: servername,
		Options:    options,
		Loc:        nodes.Loc{Start: p.stmtStart, End: p.prev.End},
	}, nil
}

func (p *Parser) parseImportForeignSchemaStmt() (nodes.Node, error) {
	if _, err := p.expect(FOREIGN); err != nil {
		return nil, err
	}
	if _, err := p.expect(SCHEMA); err != nil {
		return nil, err
	}
	remoteSchema, err := p.parseName()
	if err != nil {
		return nil, err
	}
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
		if tableList, err = p.parseRelationExprList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	} else if p.cur.Type == EXCEPT {
		p.advance()
		listType = nodes.FDW_IMPORT_SCHEMA_EXCEPT
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if tableList, err = p.parseRelationExprList(); err != nil {
			return nil, err
		}
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
	serverName, err := p.parseName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(INTO); err != nil {
		return nil, err
	}
	localSchema, err := p.parseName()
	if err != nil {
		return nil, err
	}
	options, err := p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.ImportForeignSchemaStmt{
		ServerName:   serverName,
		RemoteSchema: remoteSchema,
		LocalSchema:  localSchema,
		ListType:     listType,
		TableList:    tableList,
		Options:      options,
		Loc:          nodes.Loc{Start: p.stmtStart, End: p.prev.End},
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
	user, err := p.parseAuthIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, err := p.parseName()
	if err != nil {
		return nil, err
	}
	options, err := p.parseCreateGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.CreateUserMappingStmt{
		User:        user,
		Servername:  servername,
		Options:     options,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: p.stmtStart, End: p.prev.End},
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
	user, err := p.parseAuthIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, err := p.parseName()
	if err != nil {
		return nil, err
	}
	options, err := p.parseAlterGenericOptions()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterUserMappingStmt{
		User:       user,
		Servername: servername,
		Options:    options,
		Loc:        nodes.Loc{Start: p.stmtStart, End: p.prev.End},
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
	user, err := p.parseAuthIdent()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(SERVER); err != nil {
		return nil, err
	}
	servername, err := p.parseName()
	if err != nil {
		return nil, err
	}
	return &nodes.DropUserMappingStmt{
		User:       user,
		Servername: servername,
		MissingOk:  missingOk,
		Loc:        nodes.Loc{Start: p.stmtStart, End: p.prev.End},
	}, nil
}

func (p *Parser) parseAuthIdent() (*nodes.RoleSpec, error) {
	if p.cur.Type == USER {
		p.advance()
		return &nodes.RoleSpec{
			Roletype: int(nodes.ROLESPEC_CURRENT_USER),
		}, nil
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
	name, err := p.parseName()
	if err != nil {
		return nil, err
	}
	stmt := &nodes.CreatePLangStmt{
		Replace:   replace,
		Plname:    name,
		Pltrusted: trusted,
	}
	if p.cur.Type == HANDLER {
		p.advance()
		if stmt.Plhandler, err = p.parseHandlerName(); err != nil {
			return nil, err
		}
		if p.cur.Type == INLINE_P {
			p.advance()
			if stmt.Plinline, err = p.parseHandlerName(); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == VALIDATOR {
			p.advance()
			if stmt.Plvalidator, err = p.parseHandlerName(); err != nil {
				return nil, err
			}
		} else if p.cur.Type == NO {
			next := p.peekNext()
			if next.Type == VALIDATOR {
				p.advance()
				p.advance()
			}
		}
	}
	stmt.Loc = nodes.Loc{Start: p.stmtStart, End: p.prev.End}
	return stmt, nil
}
