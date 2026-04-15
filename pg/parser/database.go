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

	options := p.parseCreatedbOptList()

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
func (p *Parser) parseCreatedbOptList() *nodes.List {
	var items []nodes.Node
	for {
		item := p.parseCreatedbOptItem()
		if item == nil {
			break
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil
	}
	return &nodes.List{Items: items}
}

// parseCreatedbOptItem parses a single createdb_opt_item.
//
//	createdb_opt_item:
//	    createdb_opt_name opt_equal NumericOnly
//	    | createdb_opt_name opt_equal opt_boolean_or_string
//	    | createdb_opt_name opt_equal DEFAULT
func (p *Parser) parseCreatedbOptItem() *nodes.DefElem {
	optName := p.parseCreatedbOptName()
	if optName == "" {
		return nil
	}

	// opt_equal: '=' | EMPTY
	if p.cur.Type == '=' {
		p.advance()
	}

	// Value: NumericOnly | opt_boolean_or_string | DEFAULT
	if p.cur.Type == DEFAULT {
		p.advance()
		return makeDefElem(optName, nil)
	}

	// Try NumericOnly first (ICONST, FCONST, +/- followed by number)
	if p.cur.Type == ICONST || p.cur.Type == FCONST || p.cur.Type == '+' || p.cur.Type == '-' {
		val := p.parseNumericOnly()
		return makeDefElem(optName, val)
	}

	// opt_boolean_or_string
	val := p.parseOptBooleanOrString()
	return makeDefElem(optName, &nodes.String{Str: val})
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
func (p *Parser) parseCreatedbOptName() string {
	switch p.cur.Type {
	case CONNECTION:
		p.advance() // consume CONNECTION
		p.expect(LIMIT)
		return "connection_limit"
	case ENCODING:
		p.advance()
		return "encoding"
	case LOCATION:
		p.advance()
		return "location"
	case OWNER:
		p.advance()
		return "owner"
	case TABLESPACE:
		p.advance()
		return "tablespace"
	case TEMPLATE:
		p.advance()
		return "template"
	case IDENT:
		name := p.cur.Str
		p.advance()
		return name
	default:
		return ""
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
		// ALTER DATABASE name OWNER TO RoleSpec
		p.advance() // consume OWNER
		if _, err := p.expect(TO); err != nil {
			return nil, err
		}
		roleSpec := p.parseRoleSpec()
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
		return p.parseAlterDatabaseSetStmt(name, stmtLoc), nil

	case RESET:
		// ALTER DATABASE name RESET ... (SetResetClause)
		return p.parseAlterDatabaseSetStmt(name, stmtLoc), nil

	case WITH:
		// ALTER DATABASE name WITH createdb_opt_list
		p.advance() // consume WITH
		options := p.parseCreatedbOptList()
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil

	case CONNECTION:
		// ALTER DATABASE name CONNECTION LIMIT connlimit
		options := p.parseCreatedbOptList()
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
		options := p.parseCreatedbOptList()
		return &nodes.AlterDatabaseStmt{
			Dbname:  name,
			Options: options,
			Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
		}, nil
	}
}

// parseAlterDatabaseSetStmt parses ALTER DATABASE name SET/RESET ... .
// The name has already been parsed. Current token is SET or RESET.
func (p *Parser) parseAlterDatabaseSetStmt(dbname string, stmtLoc int) nodes.Node {
	setstmt := p.parseSetResetClause()
	return &nodes.AlterDatabaseSetStmt{
		Dbname:  dbname,
		Setstmt: setstmt,
		Loc:     nodes.Loc{Start: stmtLoc, End: p.prev.End},
	}
}

// parseSetResetClause parses SetResetClause.
//
//	SetResetClause:
//	    SET set_rest
//	    | VariableResetStmt
func (p *Parser) parseSetResetClause() *nodes.VariableSetStmt {
	if p.cur.Type == SET {
		p.advance() // consume SET
		result, _ := p.parseSetRest()
		if vs, ok := result.(*nodes.VariableSetStmt); ok {
			return vs
		}
		return nil
	}
	// RESET
	p.advance() // consume RESET
	result := p.parseResetRest()
	if vs, ok := result.(*nodes.VariableSetStmt); ok {
		return vs
	}
	return nil
}

// parseResetRest parses the rest of a RESET statement (after RESET keyword).
//
//	reset_rest:
//	    generic_reset
//	    | TIME ZONE
//	    | TRANSACTION ISOLATION LEVEL
//	    | SESSION AUTHORIZATION
func (p *Parser) parseResetRest() nodes.Node {
	switch p.cur.Type {
	case TIME:
		p.advance()
		p.expect(ZONE)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "timezone",
		}
	case TRANSACTION:
		p.advance()
		p.expect(ISOLATION)
		p.expect(LEVEL)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "transaction_isolation",
		}
	case SESSION:
		p.advance()
		p.expect(AUTHORIZATION)
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: "session_authorization",
		}
	case ALL:
		p.advance()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET_ALL,
		}
	default:
		name := p.parseVarName()
		return &nodes.VariableSetStmt{
			Kind: nodes.VAR_RESET,
			Name: name,
		}
	}
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
	if p.cur.Type == WITH {
		p.advance() // consume WITH
	}
	if p.cur.Type == '(' {
		p.advance() // consume '('
		options = p.parseDropOptionList()
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
func (p *Parser) parseDropOptionList() *nodes.List {
	first := p.parseDropOption()
	items := []nodes.Node{first}
	for p.cur.Type == ',' {
		p.advance()
		items = append(items, p.parseDropOption())
	}
	return &nodes.List{Items: items}
}

// parseDropOption parses drop_option.
//
//	drop_option:
//	    FORCE
func (p *Parser) parseDropOption() *nodes.DefElem {
	if p.cur.Type == FORCE {
		p.advance()
		return makeDefElem("force", nil)
	}
	// Consume whatever token is there (shouldn't happen for valid SQL)
	p.advance()
	return makeDefElem("force", nil)
}
