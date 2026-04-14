package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// CREATE DATABASE
// ---------------------------------------------------------------------------

// parseCreateDatabaseStmt parses CREATE [OR REPLACE] [TRANSIENT] DATABASE ...
// The CREATE keyword and optional OR REPLACE / TRANSIENT modifiers have already
// been consumed by parseCreateStmt; start is the Loc of the CREATE token.
func (p *Parser) parseCreateDatabaseStmt(start ast.Loc, orReplace, transient bool) (ast.Node, error) {
	p.advance() // consume DATABASE

	stmt := &ast.CreateDatabaseStmt{
		OrReplace: orReplace,
		Transient: transient,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CLONE source [AT|BEFORE (...)]
	if p.cur.Type == kwCLONE {
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
	}

	// Optional properties
	if err := p.parseDBSchemaPropsInto(&stmt.Props); err != nil {
		return nil, err
	}

	// Optional WITH TAG (...)
	for p.cur.Type == kwWITH || p.cur.Type == kwTAG {
		if p.cur.Type == kwWITH {
			if p.peekNext().Type != kwTAG {
				break
			}
			p.advance() // consume WITH
		}
		tags, err := p.parseTagAssignments()
		if err != nil {
			return nil, err
		}
		stmt.Tags = append(stmt.Tags, tags...)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA
// ---------------------------------------------------------------------------

// parseCreateSchemaStmt parses CREATE [OR REPLACE] [TRANSIENT] SCHEMA ...
func (p *Parser) parseCreateSchemaStmt(start ast.Loc, orReplace, transient bool) (ast.Node, error) {
	p.advance() // consume SCHEMA

	stmt := &ast.CreateSchemaStmt{
		OrReplace: orReplace,
		Transient: transient,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CLONE source [AT|BEFORE (...)]
	if p.cur.Type == kwCLONE {
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
	}

	// Optional WITH MANAGED ACCESS
	if p.cur.Type == kwWITH {
		if p.peekNext().Type == kwMANAGED {
			p.advance() // consume WITH
			p.advance() // consume MANAGED
			if _, err := p.expect(kwACCESS); err != nil {
				return nil, err
			}
			stmt.ManagedAccess = true
		}
	}

	// Optional properties
	if err := p.parseDBSchemaPropsInto(&stmt.Props); err != nil {
		return nil, err
	}

	// Optional WITH TAG (...)
	for p.cur.Type == kwWITH || p.cur.Type == kwTAG {
		if p.cur.Type == kwWITH {
			if p.peekNext().Type != kwTAG {
				break
			}
			p.advance() // consume WITH
		}
		tags, err := p.parseTagAssignments()
		if err != nil {
			return nil, err
		}
		stmt.Tags = append(stmt.Tags, tags...)
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER dispatch
// ---------------------------------------------------------------------------

// parseAlterStmt dispatches ALTER ... to the appropriate sub-parser.
// Handles DATABASE and SCHEMA; falls back to unsupported for other objects.
func (p *Parser) parseAlterStmt() (ast.Node, error) {
	p.advance() // consume ALTER

	switch p.cur.Type {
	case kwDATABASE:
		return p.parseAlterDatabaseStmt()
	case kwSCHEMA:
		return p.parseAlterSchemaStmt()
	default:
		return p.unsupported("ALTER")
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

// parseAlterDatabaseStmt parses ALTER DATABASE ... (all action variants).
// The ALTER keyword has already been consumed.
func (p *Parser) parseAlterDatabaseStmt() (ast.Node, error) {
	altTok := p.advance() // consume DATABASE
	stmt := &ast.AlterDatabaseStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDBRename
		stmt.NewName = newName

	case kwSWAP:
		// SWAP WITH other_name
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		otherName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDBSwap
		stmt.NewName = otherName

	case kwSET:
		p.advance() // consume SET
		// Could be SET TAG (...) or SET <properties>
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBSetTag
			stmt.Tags = tags
		} else {
			props, err := p.parseDBSchemaProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBSet
			stmt.SetProps = props
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			// UNSET TAG ( name, name, ... )
			names, err := p.parseUnsetTagList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET property [, property ...]
			props, err := p.parsePropertyNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBUnset
			stmt.UnsetProps = props
		}

	case kwENABLE:
		p.advance() // consume ENABLE
		switch p.cur.Type {
		case kwREPLICATION:
			// ENABLE REPLICATION TO ACCOUNTS account_list [IGNORE EDITION CHECK]
			p.advance() // consume REPLICATION
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwACCOUNTS); err != nil {
				return nil, err
			}
			// Skip account identifiers until end / IGNORE
			for p.cur.Type != tokEOF && p.cur.Type != kwIGNORE && p.cur.Type != ';' {
				p.advance()
			}
			// Optional IGNORE EDITION CHECK
			if p.cur.Type == kwIGNORE {
				p.advance() // consume IGNORE
				p.advance() // consume EDITION
				p.advance() // consume CHECK
			}
			stmt.Action = ast.AlterDBEnableReplication
		case kwFAILOVER:
			// ENABLE FAILOVER TO ACCOUNTS account_list
			p.advance() // consume FAILOVER
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwACCOUNTS); err != nil {
				return nil, err
			}
			// Skip account identifiers
			for p.cur.Type != tokEOF && p.cur.Type != ';' {
				p.advance()
			}
			stmt.Action = ast.AlterDBEnableFailover
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwDISABLE:
		p.advance() // consume DISABLE
		switch p.cur.Type {
		case kwREPLICATION:
			// DISABLE REPLICATION [TO ACCOUNTS account_list]
			p.advance() // consume REPLICATION
			if p.cur.Type == kwTO {
				p.advance() // consume TO
				if _, err := p.expect(kwACCOUNTS); err != nil {
					return nil, err
				}
				for p.cur.Type != tokEOF && p.cur.Type != ';' {
					p.advance()
				}
			}
			stmt.Action = ast.AlterDBDisableReplication
		case kwFAILOVER:
			// DISABLE FAILOVER [TO ACCOUNTS account_list]
			p.advance() // consume FAILOVER
			if p.cur.Type == kwTO {
				p.advance() // consume TO
				if _, err := p.expect(kwACCOUNTS); err != nil {
					return nil, err
				}
				for p.cur.Type != tokEOF && p.cur.Type != ';' {
					p.advance()
				}
			}
			stmt.Action = ast.AlterDBDisableFailover
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwREFRESH:
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterDBRefresh

	case kwPRIMARY:
		p.advance() // consume PRIMARY
		stmt.Action = ast.AlterDBPrimary

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER SCHEMA
// ---------------------------------------------------------------------------

// parseAlterSchemaStmt parses ALTER SCHEMA ... (all action variants).
// The ALTER keyword has already been consumed.
func (p *Parser) parseAlterSchemaStmt() (ast.Node, error) {
	altTok := p.advance() // consume SCHEMA
	stmt := &ast.AlterSchemaStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaRename
		stmt.NewName = newName

	case kwSWAP:
		// SWAP WITH other_name
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		otherName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaSwap
		stmt.NewName = otherName

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaSetTag
			stmt.Tags = tags
		} else {
			props, err := p.parseDBSchemaProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaSet
			stmt.SetProps = props
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseUnsetTagList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parsePropertyNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaUnset
			stmt.UnsetProps = props
		}

	case kwENABLE:
		// ENABLE MANAGED ACCESS
		p.advance() // consume ENABLE
		if _, err := p.expect(kwMANAGED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaEnableManagedAccess

	case kwDISABLE:
		// DISABLE MANAGED ACCESS
		p.advance() // consume DISABLE
		if _, err := p.expect(kwMANAGED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaDisableManagedAccess

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DROP dispatch
// ---------------------------------------------------------------------------

// parseDropStmt dispatches DROP ... to the appropriate sub-parser.
// (parseDropStmt lives in drop.go; it dispatches DATABASE/SCHEMA here.)

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

// parseDropDatabaseStmt parses DROP DATABASE [IF EXISTS] name [CASCADE|RESTRICT].
// The DROP keyword has already been consumed.
func (p *Parser) parseDropDatabaseStmt() (ast.Node, error) {
	dropTok := p.advance() // consume DATABASE
	stmt := &ast.DropDatabaseStmt{Loc: ast.Loc{Start: dropTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CASCADE | RESTRICT
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		stmt.Cascade = true
	case kwRESTRICT:
		p.advance()
		stmt.Restrict = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DROP SCHEMA
// ---------------------------------------------------------------------------

// parseDropSchemaStmt parses DROP SCHEMA [IF EXISTS] name [CASCADE|RESTRICT].
// The DROP keyword has already been consumed.
func (p *Parser) parseDropSchemaStmt() (ast.Node, error) {
	dropTok := p.advance() // consume SCHEMA
	stmt := &ast.DropSchemaStmt{Loc: ast.Loc{Start: dropTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CASCADE | RESTRICT
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		stmt.Cascade = true
	case kwRESTRICT:
		p.advance()
		stmt.Restrict = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP dispatch
// ---------------------------------------------------------------------------

// (parseUndropStmt lives in drop.go; it dispatches DATABASE/SCHEMA here.)

// ---------------------------------------------------------------------------
// UNDROP DATABASE
// ---------------------------------------------------------------------------

// parseUndropDatabaseStmt parses UNDROP DATABASE name.
// The UNDROP keyword has already been consumed.
func (p *Parser) parseUndropDatabaseStmt() (ast.Node, error) {
	undropTok := p.advance() // consume DATABASE
	stmt := &ast.UndropDatabaseStmt{Loc: ast.Loc{Start: undropTok.Loc.Start}}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP SCHEMA
// ---------------------------------------------------------------------------

// parseUndropSchemaStmt parses UNDROP SCHEMA name.
// The UNDROP keyword has already been consumed.
func (p *Parser) parseUndropSchemaStmt() (ast.Node, error) {
	undropTok := p.advance() // consume SCHEMA
	stmt := &ast.UndropSchemaStmt{Loc: ast.Loc{Start: undropTok.Loc.Start}}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Property parsing helpers
// ---------------------------------------------------------------------------

// parseDBSchemaProps parses zero or more optional DATABASE/SCHEMA properties
// (DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS,
// DEFAULT_DDL_COLLATION, COMMENT). Returns a new *DBSchemaProps; the pointer
// fields inside are nil unless the corresponding property was present.
func (p *Parser) parseDBSchemaProps() (*ast.DBSchemaProps, error) {
	props := &ast.DBSchemaProps{}
	if err := p.parseDBSchemaPropsInto(props); err != nil {
		return nil, err
	}
	return props, nil
}

// parseDBSchemaPropsInto populates *props in-place. Used by CREATE statements
// that hold DBSchemaProps by value.
func (p *Parser) parseDBSchemaPropsInto(props *ast.DBSchemaProps) error {
	for {
		switch p.cur.Type {
		case kwDATA_RETENTION_TIME_IN_DAYS:
			p.advance() // consume DATA_RETENTION_TIME_IN_DAYS
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokInt)
			if err != nil {
				return err
			}
			v := tok.Ival
			props.DataRetention = &v

		case kwMAX_DATA_EXTENSION_TIME_IN_DAYS:
			p.advance() // consume MAX_DATA_EXTENSION_TIME_IN_DAYS
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokInt)
			if err != nil {
				return err
			}
			v := tok.Ival
			props.MaxDataExt = &v

		case kwDEFAULT_DDL_COLLATION:
			p.advance() // consume DEFAULT_DDL_COLLATION / DEFAULT_DDL_COLLATION_
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			props.DefaultDDLCol = &s

		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			props.Comment = &s

		default:
			return nil
		}
	}
}

// parsePropertyNameList parses a comma-separated list of property names
// (identifiers or reserved keywords) for UNSET clauses. Returns the uppercased
// names. Consumes at least one name.
func (p *Parser) parsePropertyNameList() ([]string, error) {
	name, err := p.parsePropertyName()
	if err != nil {
		return nil, err
	}
	names := []string{name}

	for p.cur.Type == ',' {
		p.advance() // consume ','
		name, err = p.parsePropertyName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

// parsePropertyName parses a single property name. Property names can be
// keywords (DATA_RETENTION_TIME_IN_DAYS, COMMENT, etc.) or identifiers.
func (p *Parser) parsePropertyName() (string, error) {
	switch p.cur.Type {
	case kwDATA_RETENTION_TIME_IN_DAYS:
		p.advance()
		return "DATA_RETENTION_TIME_IN_DAYS", nil
	case kwMAX_DATA_EXTENSION_TIME_IN_DAYS:
		p.advance()
		return "MAX_DATA_EXTENSION_TIME_IN_DAYS", nil
	case kwDEFAULT_DDL_COLLATION:
		p.advance()
		return "DEFAULT_DDL_COLLATION", nil
	case kwCOMMENT:
		p.advance()
		return "COMMENT", nil
	case tokIdent:
		tok := p.advance()
		return tok.Str, nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// parseUnsetTagList parses TAG ( name [, name ...] ) for UNSET TAG clauses.
// Consumes the TAG keyword and the parenthesized name list.
func (p *Parser) parseUnsetTagList() ([]*ast.ObjectName, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var names []*ast.ObjectName
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}
