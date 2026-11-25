package parser

import (
	nodes "github.com/bytebase/omni/pg/ast"
)

// parseCreatedbStmt parses a CREATE DATABASE statement.
// The CREATE keyword has already been consumed. Current token is DATABASE.
//
// Ref: gram.y CreatedbStmt
//
//	CREATE DATABASE name opt_with createdb_opt_list
func (p *Parser) parseCreatedbStmt(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// opt_with: WITH | EMPTY
	if p.cur.Type == WITH {
		p.advance()
	}

	options, err := p.parseCreatedbOptList()
	if err != nil {
		return nil, err
	}

	return &nodes.CreatedbStmt{
		Dbname:  name,
		Options: options,
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseCreatedbOptList parses createdb_opt_list.
//
//	createdb_opt_list:
//	    createdb_opt_items | EMPTY
func (p *Parser) parseCreatedbOptList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		item, err := p.parseCreatedbOptItem()
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

// parseCreatedbOptItem parses a single createdb_opt_item.
//
//	createdb_opt_item:
//	    createdb_opt_name opt_equal NumericOnly
//	    | createdb_opt_name opt_equal opt_boolean_or_string
//	    | createdb_opt_name opt_equal DEFAULT
func (p *Parser) parseCreatedbOptItem() (*nodes.DefElem, error) {
	optName, err := p.parseCreatedbOptName()
	if err != nil || optName == "" {
		return nil, err
	}

	// opt_equal: '=' | EMPTY
	if p.cur.Type == '=' {
		p.advance()
	}

	switch p.cur.Type {
	case DEFAULT:
		p.advance()
		return makeDefElem(optName, nil), nil
	case ICONST, FCONST, '+', '-':
		val, err := p.parseNumericOnlyRequired()
		if err != nil {
			return nil, err
		}
		return makeDefElem(optName, val), nil
	case TRUE_P, FALSE_P, ON, SCONST:
		return makeDefElem(optName, &nodes.String{Str: p.parseOptBooleanOrString()}), nil
	}
	// opt_boolean_or_string: ... | NonReservedWord_or_Sconst
	if p.isNonReservedWord() {
		return makeDefElem(optName, &nodes.String{Str: p.parseOptBooleanOrString()}), nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseCreatedbOptName parses createdb_opt_name.
//
//	createdb_opt_name:
//	    IDENT
//	    | CONNECTION LIMIT
//	    | ENCODING
//	    | LOCATION
//	    | OWNER
//	    | TABLESPACE
//	    | TEMPLATE
func (p *Parser) parseCreatedbOptName() (string, error) {
	switch p.cur.Type {
	case CONNECTION:
		p.advance() // consume CONNECTION
		if _, err := p.expect(LIMIT); err != nil {
			return "", err
		}
		return "connection_limit", nil
	case ENCODING:
		p.advance()
		return "encoding", nil
	case LOCATION:
		p.advance()
		return "location", nil
	case OWNER:
		p.advance()
		return "owner", nil
	case TABLESPACE:
		p.advance()
		return "tablespace", nil
	case TEMPLATE:
		p.advance()
		return "template", nil
	case IDENT:
		name := p.cur.Str
		p.advance()
		return name, nil
	default:
		return "", nil
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

// parseAlterDatabaseDispatch parses ALTER DATABASE statements.
// The ALTER keyword has already been consumed. Current token is DATABASE.
//
// This dispatches between:
//   - AlterDatabaseStmt (ALTER DATABASE name [WITH] createdb_opt_list)
//   - AlterDatabaseStmt (ALTER DATABASE name SET TABLESPACE name)
//   - AlterDatabaseSetStmt (ALTER DATABASE name SET/RESET ...)
//   - AlterDatabaseStmt (ALTER DATABASE name CONNECTION LIMIT ...)
func (p *Parser) parseAlterDatabaseDispatch(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case RENAME:
		// ALTER DATABASE name RENAME TO name
		p.advance() // consume RENAME
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		newname, err := p.parseName()
		if err != nil {
			return nil, err
		}
		return &nodes.RenameStmt{
			RenameType: nodes.OBJECT_DATABASE,
			Subname:    name,
			Newname:    newname,
		}, nil

	case OWNER:
		// ALTER DATABASE name OWNER TO RoleSpec; without TO, OWNER is a
		// createdb_opt_name (OWNER [=] value).
		if p.peekNext().Type != TO && !p.completingAtNext() {
			options, err := p.parseCreatedbOptList()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterDatabaseStmt{
				Dbname:  name,
				Options: options,
				Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
			}, nil
		}
		p.advance() // consume OWNER
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec, err := p.parseRoleSpec()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterOwnerStmt{
			ObjectType: nodes.OBJECT_DATABASE,
			Object:     &nodes.String{Str: name},
			Newowner:   roleSpec,
		}, nil

	case SET:
		// Could be SET TABLESPACE or SET variable
		next := p.peekNext()
		if next.Type == TABLESPACE {
			// ALTER DATABASE name SET TABLESPACE name
			p.advance() // consume SET
			p.advance() // consume TABLESPACE
			tbsName, err := p.parseName()
			if err != nil {
				return nil, err
			}
			return &nodes.AlterDatabaseStmt{
				Dbname:  name,
				Options: &nodes.List{Items: []nodes.Node{makeDefElem("tablespace", &nodes.String{Str: tbsName})}},
				Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
			}, nil
		}
		// ALTER DATABASE name SET variable = value (SetResetClause)
		return p.parseAlterDatabaseSetStmt(name, stmtLoc)

	case RESET:
		// ALTER DATABASE name RESET ... (SetResetClause)
		return p.parseAlterDatabaseSetStmt(name, stmtLoc)

	case WITH:
		// ALTER DATABASE name WITH createdb_opt_list
		p.advance() // consume WITH
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	case CONNECTION:
		// ALTER DATABASE name CONNECTION LIMIT connlimit
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	case REFRESH:
		// ALTER DATABASE name REFRESH COLLATION VERSION
		p.advance() // consume REFRESH
		if _, err := p.expect(COLLATION); err != nil {
			return nil, err
		}
		if _, err := p.expect(VERSION_P); err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseRefreshCollStmt{
			Dbname: name,
			Loc:    nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	default:
		// ALTER DATABASE name createdb_opt_list (empty or options)
		options, err := p.parseCreatedbOptList()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil
	}
}

// parseAlterDatabaseSetStmt parses ALTER DATABASE name SET/RESET ... .
// The name has already been parsed. Current token is SET or RESET.
func (p *Parser) parseAlterDatabaseSetStmt(dbname string, stmtLoc int) (nodes.Node, error) {
	setstmt, err := p.parseSetResetClause()
	if err != nil {
		return nil, err
	}
	return &nodes.AlterDatabaseSetStmt{
		Dbname:  dbname,
		Setstmt: setstmt,
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseSetResetClause parses SetResetClause: SET set_rest, which takes
// no LOCAL or SESSION, or RESET reset_rest.
func (p *Parser) parseSetResetClause() (*nodes.VariableSetStmt, error) {
	var result nodes.Node
	var err error
	switch p.cur.Type {
	case SET:
		setLoc := p.pos()
		p.advance()
		if result, err = p.parseSetRest(); err != nil {
			return nil, err
		}
		if vs, ok := result.(*nodes.VariableSetStmt); ok {
			vs.Loc = nodes.Loc{Start: setLoc, End: p.prev.End}
		}
	case RESET:
		p.advance()
		if result, err = p.parseVariableResetStmt(); err != nil {
			return nil, err
		}
	default:
		if p.collectMode() {
			p.addTokenCandidate(SET)
			p.addTokenCandidate(RESET)
			return nil, errCollecting
		}
		return nil, p.syntaxErrorAtCur()
	}
	vs, _ := result.(*nodes.VariableSetStmt)
	return vs, nil
}

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

// parseDropdbStmt parses a DROP DATABASE statement.
// The DROP keyword has already been consumed. Current token is DATABASE.
//
// Ref: gram.y DropdbStmt
//
//	DROP DATABASE name
//	| DROP DATABASE IF_P EXISTS name
//	| DROP DATABASE name opt_with '(' drop_option_list ')'
//	| DROP DATABASE IF_P EXISTS name opt_with '(' drop_option_list ')'
func (p *Parser) parseDropdbStmt(stmtLoc int) (nodes.Node, error) {
	p.advance() // consume DATABASE

	missingOk := false
	if p.cur.Type == IF_P {
		p.advance() // consume IF
		if _, err := p.expect(EXISTS); err != nil {
			return nil, err
		}
		missingOk = true
	}

	name, err := p.parseName()
	if err != nil {
		return nil, err
	}

	// Optional: opt_with '(' drop_option_list ')'
	var options *nodes.List
	if p.cur.Type == WITH || p.cur.Type == '(' {
		if p.cur.Type == WITH {
			p.advance() // consume WITH
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		if options, err = p.parseDropOptionList(); err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	return &nodes.DropdbStmt{
		Dbname:    name,
		MissingOk: missingOk,
		Options:   options,
		Loc:       nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}, nil
}

// parseDropOptionList parses drop_option_list.
//
//	drop_option_list:
//	    drop_option
//	    | drop_option_list ',' drop_option
func (p *Parser) parseDropOptionList() (*nodes.List, error) {
	var items []nodes.Node
	for {
		opt, err := p.parseDropOption()
		if err != nil {
			return nil, err
		}
		items = append(items, opt)
		if p.cur.Type != ',' {
			return &nodes.List{Items: items}, nil
		}
		p.advance()
	}
}

// parseDropOption parses drop_option.
//
//	drop_option:
//	    FORCE
func (p *Parser) parseDropOption() (*nodes.DefElem, error) {
	if _, err := p.expect(FORCE); err != nil {
		return nil, err
	}
	return makeDefElem("force", nil), nil
}
