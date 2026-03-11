package parser

import (
	nodes "github.com/bytebase/omni/mysql/ast"
)

// parseAlterTableStmt parses an ALTER TABLE statement.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
//
//	ALTER TABLE tbl_name
//	    [alter_option [, alter_option] ...]
//
//	alter_option:
//	    ADD [COLUMN] col_name column_definition [FIRST | AFTER col_name]
//	    | ADD [COLUMN] (col_name column_definition, ...)
//	    | ADD {INDEX | KEY} [index_name] (key_part,...) [index_option] ...
//	    | ADD {FULLTEXT | SPATIAL} [INDEX | KEY] [index_name] (key_part,...) [index_option] ...
//	    | ADD [CONSTRAINT [symbol]] PRIMARY KEY (key_part,...) [index_option] ...
//	    | ADD [CONSTRAINT [symbol]] UNIQUE [INDEX | KEY] [index_name] (key_part,...) [index_option] ...
//	    | ADD [CONSTRAINT [symbol]] FOREIGN KEY [index_name] (col,...) reference_definition
//	    | ADD [CONSTRAINT [symbol]] CHECK (expr)
//	    | DROP [COLUMN] col_name
//	    | DROP {INDEX | KEY} index_name
//	    | DROP PRIMARY KEY
//	    | DROP FOREIGN KEY fk_symbol
//	    | DROP CHECK symbol
//	    | MODIFY [COLUMN] col_name column_definition [FIRST | AFTER col_name]
//	    | CHANGE [COLUMN] old_col_name new_col_name column_definition [FIRST | AFTER col_name]
//	    | ALTER [COLUMN] col_name {SET DEFAULT {literal | (expr)} | DROP DEFAULT}
//	    | RENAME [TO | AS] new_tbl_name
//	    | RENAME COLUMN old_col_name TO new_col_name
//	    | RENAME {INDEX | KEY} old_index_name TO new_index_name
//	    | CONVERT TO CHARACTER SET charset_name [COLLATE collation_name]
//	    | [DEFAULT] CHARACTER SET [=] charset_name [COLLATE [=] collation_name]
//	    | ALGORITHM [=] {DEFAULT | INSTANT | INPLACE | COPY}
//	    | LOCK [=] {DEFAULT | NONE | SHARED | EXCLUSIVE}
func (p *Parser) parseAlterTableStmt() (*nodes.AlterTableStmt, error) {
	start := p.pos()
	p.advance() // consume TABLE

	stmt := &nodes.AlterTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Table name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// Parse comma-separated alter operations
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}

		cmd, err := p.parseAlterTableCmd()
		if err != nil {
			return nil, err
		}
		stmt.Commands = append(stmt.Commands, cmd)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}

// parseAlterTableCmd parses a single ALTER TABLE operation.
func (p *Parser) parseAlterTableCmd() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: start}}

	switch p.cur.Type {
	case kwADD:
		return p.parseAlterAdd(cmd)

	case kwDROP:
		return p.parseAlterDrop(cmd)

	case kwMODIFY:
		return p.parseAlterModify(cmd)

	case kwCHANGE:
		return p.parseAlterChange(cmd)

	case kwALTER:
		return p.parseAlterColumn(cmd)

	case kwRENAME:
		return p.parseAlterRename(cmd)

	case kwCONVERT:
		return p.parseAlterConvert(cmd)

	case kwALGORITHM:
		p.advance()
		p.match('=')
		cmd.Type = nodes.ATAlgorithm
		cmd.Name = p.consumeOptionValue()
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwLOCK:
		p.advance()
		p.match('=')
		cmd.Type = nodes.ATLock
		cmd.Name = p.consumeOptionValue()
		cmd.Loc.End = p.pos()
		return cmd, nil

	default:
		// Try table options: ENGINE, CHARSET, etc.
		opt, ok := p.parseTableOption()
		if ok {
			cmd.Type = nodes.ATTableOption
			cmd.Option = opt
			cmd.Loc.End = p.pos()
			return cmd, nil
		}
		return nil, &ParseError{
			Message:  "expected ALTER TABLE operation",
			Position: p.cur.Loc,
		}
	}
}

// parseAlterAdd parses ADD operations.
func (p *Parser) parseAlterAdd(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume ADD

	// ADD [CONSTRAINT ...] PRIMARY KEY / UNIQUE / FOREIGN KEY / CHECK
	if p.isTableConstraintStart() {
		cmd.Type = nodes.ATAddConstraint
		constr, err := p.parseTableConstraint()
		if err != nil {
			return nil, err
		}
		cmd.Constraint = constr
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

	// ADD [COLUMN]
	cmd.Type = nodes.ATAddColumn
	p.match(kwCOLUMN)

	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col

	// FIRST | AFTER col_name
	p.parseColumnPositioning(cmd)

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterDrop parses DROP operations.
func (p *Parser) parseAlterDrop(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume DROP

	switch p.cur.Type {
	case kwPRIMARY:
		// DROP PRIMARY KEY
		p.advance()
		p.match(kwKEY)
		cmd.Type = nodes.ATDropConstraint
		cmd.Name = "PRIMARY"
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwINDEX, kwKEY:
		// DROP {INDEX | KEY} index_name
		p.advance()
		cmd.Type = nodes.ATDropIndex
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwFOREIGN:
		// DROP FOREIGN KEY fk_symbol
		p.advance()
		p.match(kwKEY)
		cmd.Type = nodes.ATDropConstraint
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwCHECK:
		// DROP CHECK symbol
		p.advance()
		cmd.Type = nodes.ATDropConstraint
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Loc.End = p.pos()
		return cmd, nil

	default:
		// DROP [COLUMN] [IF EXISTS] col_name
		cmd.Type = nodes.ATDropColumn
		p.match(kwCOLUMN)
		if p.cur.Type == kwIF {
			p.advance()
			p.match(kwNOT) // IF EXISTS (not IF NOT)
			// Actually it's IF EXISTS for DROP
			// But let's check: MySQL uses DROP [COLUMN] [IF EXISTS] col_name
			// Here we consumed IF, let's check for EXISTS
			if p.cur.Type == kwEXISTS_KW {
				p.advance()
				cmd.IfExists = true
			}
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Loc.End = p.pos()
		return cmd, nil
	}
}

// parseAlterModify parses MODIFY [COLUMN] col_name column_definition [FIRST | AFTER col_name].
func (p *Parser) parseAlterModify(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume MODIFY
	cmd.Type = nodes.ATModifyColumn
	p.match(kwCOLUMN)

	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col
	cmd.Name = col.Name

	p.parseColumnPositioning(cmd)

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterChange parses CHANGE [COLUMN] old_col_name new_col_name column_definition [FIRST | AFTER col_name].
func (p *Parser) parseAlterChange(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume CHANGE
	cmd.Type = nodes.ATChangeColumn
	p.match(kwCOLUMN)

	// Old column name
	oldName, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cmd.Name = oldName

	// New column definition (includes new name)
	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col
	cmd.NewName = col.Name

	p.parseColumnPositioning(cmd)

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterColumn parses ALTER [COLUMN] col_name {SET DEFAULT ... | DROP DEFAULT | SET NOT NULL | DROP NOT NULL}.
func (p *Parser) parseAlterColumn(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume ALTER
	p.match(kwCOLUMN)

	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cmd.Name = name

	if _, ok := p.match(kwSET); ok {
		if _, ok := p.match(kwDEFAULT); ok {
			cmd.Type = nodes.ATAlterColumnDefault
			// Default value: literal or (expr)
			// We don't store the default value in cmd currently
		} else if p.cur.Type == kwNOT {
			p.advance()
			p.match(kwNULL)
			cmd.Type = nodes.ATAlterColumnSetNotNull
		}
	} else if _, ok := p.match(kwDROP); ok {
		if _, ok := p.match(kwDEFAULT); ok {
			cmd.Type = nodes.ATAlterColumnDefault
		} else if _, ok := p.match(kwNOT); ok {
			p.match(kwNULL)
			cmd.Type = nodes.ATAlterColumnDropNotNull
		}
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterRename parses RENAME operations.
func (p *Parser) parseAlterRename(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume RENAME

	switch p.cur.Type {
	case kwCOLUMN:
		// RENAME COLUMN old_col_name TO new_col_name
		p.advance()
		cmd.Type = nodes.ATRenameColumn
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = oldName
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.NewName = newName

	case kwINDEX, kwKEY:
		// RENAME {INDEX | KEY} old_index_name TO new_index_name
		p.advance()
		cmd.Type = nodes.ATRenameIndex
		oldName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = oldName
		p.match(kwTO)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.NewName = newName

	default:
		// RENAME [TO | AS] new_tbl_name
		cmd.Type = nodes.ATRenameTable
		p.match(kwTO, kwAS)
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.NewName = newName
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterConvert parses CONVERT TO CHARACTER SET charset_name [COLLATE collation_name].
func (p *Parser) parseAlterConvert(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume CONVERT
	cmd.Type = nodes.ATConvertCharset

	p.match(kwTO)
	p.match(kwCHARACTER)
	p.match(kwSET)

	charset := p.consumeOptionValue()
	cmd.Name = charset

	if _, ok := p.match(kwCOLLATE); ok {
		cmd.NewName = p.consumeOptionValue()
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseColumnPositioning parses optional FIRST | AFTER col_name.
func (p *Parser) parseColumnPositioning(cmd *nodes.AlterTableCmd) {
	if _, ok := p.match(kwFIRST); ok {
		cmd.First = true
	} else if _, ok := p.match(kwAFTER); ok {
		name, _, _ := p.parseIdentifier()
		cmd.After = name
	}
}
