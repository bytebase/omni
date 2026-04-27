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

	// Completion: after ALTER TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	stmt := &nodes.AlterTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Table name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Table = ref

	// Completion: after ALTER TABLE table-name, offer sub-command keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, tok := range []int{kwADD, kwDROP, kwMODIFY, kwCHANGE, kwRENAME, kwALTER, kwCONVERT, kwENGINE, kwDEFAULT, kwORDER, kwALGORITHM, kwLOCK, kwFORCE} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

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
	// Completion: offer ALTER TABLE operation keywords.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwADD, kwDROP, kwMODIFY, kwCHANGE, kwRENAME, kwALTER, kwCONVERT, kwENGINE, kwDEFAULT, kwORDER, kwALGORITHM, kwLOCK, kwFORCE, kwPARTITION} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

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
		// Completion: after ALGORITHM=, offer algorithm keywords.
		p.checkCursor()
		if p.collectMode() {
			p.addTokenCandidate(kwDEFAULT)
			// INPLACE, COPY, INSTANT are identifiers, not keywords, but we add them via a rule
			// that the completion module will handle. For now, add DEFAULT as keyword.
			return nil, &ParseError{Message: "collecting"}
		}
		cmd.Type = nodes.ATAlgorithm
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, err
		}
		cmd.Name = val
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwLOCK:
		p.advance()
		p.match('=')
		cmd.Type = nodes.ATLock
		val, err := p.consumeOptionValue()
		if err != nil {
			return nil, err
		}
		cmd.Name = val
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwCOALESCE:
		return p.parseAlterCoalescePartition(cmd)

	case kwREORGANIZE:
		return p.parseAlterReorganizePartition(cmd)

	case kwEXCHANGE:
		return p.parseAlterExchangePartition(cmd)

	case kwTRUNCATE:
		// TRUNCATE PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATTruncatePartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after TRUNCATE",
			Position: p.cur.Loc,
		}

	case kwANALYZE:
		// ANALYZE PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATAnalyzePartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after ANALYZE",
			Position: p.cur.Loc,
		}

	case kwCHECK:
		// CHECK PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATCheckPartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after CHECK",
			Position: p.cur.Loc,
		}

	case kwOPTIMIZE:
		// OPTIMIZE PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATOptimizePartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after OPTIMIZE",
			Position: p.cur.Loc,
		}

	case kwREBUILD:
		// REBUILD PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATRebuildPartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after REBUILD",
			Position: p.cur.Loc,
		}

	case kwREPAIR:
		// REPAIR PARTITION ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATRepairPartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after REPAIR",
			Position: p.cur.Loc,
		}

	case kwDISCARD:
		// DISCARD TABLESPACE | DISCARD PARTITION {partition_names | ALL} TABLESPACE
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd, err := p.parseAlterPartitionNamesOrAll(cmd, nodes.ATDiscardPartitionTablespace)
			if err != nil {
				return nil, err
			}
			p.match(kwTABLESPACE)
			cmd.Loc.End = p.pos()
			return cmd, nil
		}
		// DISCARD TABLESPACE
		p.match(kwTABLESPACE)
		cmd.Type = nodes.ATDiscardTablespace
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwIMPORT:
		// IMPORT TABLESPACE | IMPORT PARTITION {partition_names | ALL} TABLESPACE
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd, err := p.parseAlterPartitionNamesOrAll(cmd, nodes.ATImportPartitionTablespace)
			if err != nil {
				return nil, err
			}
			p.match(kwTABLESPACE)
			cmd.Loc.End = p.pos()
			return cmd, nil
		}
		// IMPORT TABLESPACE
		p.match(kwTABLESPACE)
		cmd.Type = nodes.ATImportTablespace
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwFORCE:
		p.advance()
		cmd.Type = nodes.ATForce
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwORDER:
		// ORDER BY col_name [ASC | DESC] [, ...]
		p.advance()
		p.match(kwBY)
		cmd.Type = nodes.ATOrderBy
		items, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		cmd.OrderByItems = items
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwENABLE:
		// ENABLE KEYS
		p.advance()
		p.match(kwKEYS)
		cmd.Type = nodes.ATEnableKeys
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwDISABLE:
		// DISABLE KEYS
		p.advance()
		p.match(kwKEYS)
		cmd.Type = nodes.ATDisableKeys
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwWITH:
		// WITH VALIDATION
		p.advance()
		p.match(kwVALIDATION)
		cmd.Type = nodes.ATWithValidation
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwWITHOUT:
		// WITHOUT VALIDATION
		p.advance()
		p.match(kwVALIDATION)
		cmd.Type = nodes.ATWithoutValidation
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwPARTITION:
		// PARTITION BY ... (repartition)
		if p.peekNext().Type == kwBY {
			part, err := p.parsePartitionClause()
			if err != nil {
				return nil, err
			}
			cmd.Type = nodes.ATPartitionBy
			cmd.PartitionBy = part
			cmd.Loc.End = p.pos()
			return cmd, nil
		}
		return nil, &ParseError{
			Message:  "expected BY after PARTITION",
			Position: p.cur.Loc,
		}

	case kwREMOVE:
		// REMOVE PARTITIONING
		p.advance()
		p.match(kwPARTITIONING)
		cmd.Type = nodes.ATRemovePartitioning
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwSECONDARY_LOAD:
		p.advance()
		cmd.Type = nodes.ATSecondaryLoad
		cmd.Loc.End = p.pos()
		return cmd, nil

	case kwSECONDARY_UNLOAD:
		p.advance()
		cmd.Type = nodes.ATSecondaryUnload
		cmd.Loc.End = p.pos()
		return cmd, nil

	default:
		// Try table options: ENGINE, CHARSET, etc.
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
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

	// Completion: after ADD, offer what can be added.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwCOLUMN, kwINDEX, kwKEY, kwUNIQUE, kwPRIMARY, kwFOREIGN, kwCONSTRAINT, kwCHECK, kwPARTITION, kwSPATIAL, kwFULLTEXT} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// ADD PARTITION (partition_definition, ...)
	if p.cur.Type == kwPARTITION {
		p.advance()
		cmd.Type = nodes.ATAddPartition
		if p.cur.Type == kwPARTITIONS {
			p.advance()
			if p.cur.Type != tokICONST {
				return nil, &ParseError{Message: "expected partition count", Position: p.cur.Loc}
			}
			cmd.Number = int(p.cur.Ival)
			p.advance()
			cmd.Loc.End = p.pos()
			return cmd, nil
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for {
			if p.cur.Type == ')' {
				break
			}
			pd, err := p.parsePartitionDef()
			if err != nil {
				return nil, err
			}
			cmd.PartitionDefs = append(cmd.PartitionDefs, pd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

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

	// Completion: after ADD COLUMN, identifier context (no specific candidates).
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new column name.
		return nil, &ParseError{Message: "collecting"}
	}

	// ADD (col1 INT, col2 INT, ...) — multi-column form
	if p.cur.Type == '(' {
		p.advance()
		for {
			col, err := p.parseColumnDef()
			if err != nil {
				return nil, err
			}
			cmd.Columns = append(cmd.Columns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col

	// FIRST | AFTER col_name
	if err := p.parseColumnPositioning(cmd); err != nil {
		return nil, err
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterDrop parses DROP operations.
func (p *Parser) parseAlterDrop(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume DROP

	// Completion: after DROP, offer what can be dropped.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwCOLUMN, kwINDEX, kwKEY, kwFOREIGN, kwPRIMARY, kwCHECK, kwCONSTRAINT, kwPARTITION} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwPARTITION:
		// DROP PARTITION partition_names
		p.advance()
		cmd.Type = nodes.ATDropPartition
		names, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		cmd.PartitionNames = names
		cmd.Loc.End = p.pos()
		return cmd, nil

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
		// Completion: after DROP INDEX, offer index_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("index_ref")
			return nil, &ParseError{Message: "collecting"}
		}
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
		// Completion: after DROP FOREIGN KEY, offer constraint_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("constraint_ref")
			return nil, &ParseError{Message: "collecting"}
		}
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

	case kwCONSTRAINT:
		// DROP CONSTRAINT symbol
		p.advance()
		cmd.Type = nodes.ATDropConstraint
		// Completion: after DROP CONSTRAINT, offer constraint_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("constraint_ref")
			return nil, &ParseError{Message: "collecting"}
		}
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
		// Completion: after DROP COLUMN, offer columnref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}
		if p.cur.Type == kwIF {
			p.advance()
			if _, err := p.expect(kwEXISTS_KW); err != nil {
				return nil, err
			}
			cmd.IfExists = true
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

	// Completion: after MODIFY [COLUMN], offer columnref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col
	cmd.Name = col.Name

	if err := p.parseColumnPositioning(cmd); err != nil {
		return nil, err
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterChange parses CHANGE [COLUMN] old_col_name new_col_name column_definition [FIRST | AFTER col_name].
func (p *Parser) parseAlterChange(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume CHANGE
	cmd.Type = nodes.ATChangeColumn
	p.match(kwCOLUMN)

	// Completion: after CHANGE [COLUMN], offer columnref (old name).
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

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

	if err := p.parseColumnPositioning(cmd); err != nil {
		return nil, err
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterColumn parses ALTER [COLUMN] col_name {SET DEFAULT | DROP DEFAULT | SET NOT NULL | DROP NOT NULL | SET VISIBLE | SET INVISIBLE}
// and ALTER INDEX index_name {VISIBLE | INVISIBLE}.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
//
//	ALTER [COLUMN] col_name {
//	    SET DEFAULT {literal | (expr)}
//	  | SET {VISIBLE | INVISIBLE}
//	  | DROP DEFAULT
//	  | SET NOT NULL
//	  | DROP NOT NULL
//	}
//	ALTER INDEX index_name {VISIBLE | INVISIBLE}
func (p *Parser) parseAlterColumn(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume ALTER

	// ALTER INDEX index_name {VISIBLE | INVISIBLE}
	if p.cur.Type == kwINDEX {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		if _, ok := p.match(kwVISIBLE); ok {
			cmd.Type = nodes.ATAlterIndexVisible
		} else if _, ok := p.match(kwINVISIBLE); ok {
			cmd.Type = nodes.ATAlterIndexInvisible
		}
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

	// ALTER CHECK symbol [NOT] ENFORCED
	if p.cur.Type == kwCHECK {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Type = nodes.ATAlterCheckEnforced
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == kwENFORCED {
				p.advance()
			}
			cmd.NewName = "NOT ENFORCED"
		} else if p.cur.Type == kwENFORCED {
			p.advance()
			cmd.NewName = "ENFORCED"
		}
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

	// ALTER CONSTRAINT symbol [NOT] ENFORCED
	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Type = nodes.ATAlterCheckEnforced
		if p.cur.Type == kwNOT {
			p.advance()
			if p.cur.Type == kwENFORCED {
				p.advance()
			}
			cmd.NewName = "NOT ENFORCED"
		} else if p.cur.Type == kwENFORCED {
			p.advance()
			cmd.NewName = "ENFORCED"
		}
		cmd.Loc.End = p.pos()
		return cmd, nil
	}

	// ALTER [COLUMN] col_name ...
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
			if p.cur.Type == '(' {
				p.advance()
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if expr == nil {
					return nil, p.syntaxErrorAtCur()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				cmd.DefaultExpr = expr
			} else {
				expr, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				if expr == nil {
					return nil, p.syntaxErrorAtCur()
				}
				cmd.DefaultExpr = expr
			}
		} else if p.cur.Type == kwNOT {
			p.advance()
			p.match(kwNULL)
			cmd.Type = nodes.ATAlterColumnSetNotNull
		} else if _, ok := p.match(kwVISIBLE); ok {
			cmd.Type = nodes.ATAlterColumnVisible
		} else if _, ok := p.match(kwINVISIBLE); ok {
			cmd.Type = nodes.ATAlterColumnInvisible
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

	// Completion: after RENAME, no specific rule needed — identifier context for table rename.
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new name.
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwCOLUMN:
		// RENAME COLUMN old_col_name TO new_col_name
		p.advance()
		cmd.Type = nodes.ATRenameColumn
		// Completion: after RENAME COLUMN, offer columnref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, &ParseError{Message: "collecting"}
		}
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
		// Completion: after RENAME INDEX, offer index_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("index_ref")
			return nil, &ParseError{Message: "collecting"}
		}
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

	// Completion: after CONVERT TO CHARACTER SET, offer charset candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("charset")
		return nil, &ParseError{Message: "collecting"}
	}

	charset, err := p.consumeOptionValue()
	if err != nil {
		return nil, err
	}
	cmd.Name = charset

	if _, ok := p.match(kwCOLLATE); ok {
		collation, err := p.consumeOptionValue()
		if err != nil {
			return nil, err
		}
		cmd.NewName = collation
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseIdentList parses a comma-separated list of identifiers.
func (p *Parser) parseIdentList() ([]string, error) {
	var names []string
	for {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return names, nil
}

// parseAlterPartitionNamesOrAll parses {partition_names | ALL} for partition operations.
func (p *Parser) parseAlterPartitionNamesOrAll(cmd *nodes.AlterTableCmd, cmdType nodes.AlterTableCmdType) (*nodes.AlterTableCmd, error) {
	cmd.Type = cmdType
	if p.cur.Type == kwALL {
		p.advance()
		cmd.AllPartitions = true
	} else {
		names, err := p.parseIdentList()
		if err != nil {
			return nil, err
		}
		cmd.PartitionNames = names
	}
	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterCoalescePartition parses COALESCE PARTITION number.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
//
//	COALESCE PARTITION number
func (p *Parser) parseAlterCoalescePartition(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume COALESCE
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	cmd.Type = nodes.ATCoalescePartition
	if p.cur.Type == tokICONST {
		cmd.Number = int(p.cur.Ival)
		p.advance()
	}
	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterReorganizePartition parses REORGANIZE PARTITION partition_names INTO (partition_definitions).
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
//
//	REORGANIZE PARTITION partition_names INTO (partition_definitions)
func (p *Parser) parseAlterReorganizePartition(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume REORGANIZE
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	cmd.Type = nodes.ATReorganizePartition

	// Parse partition names
	names, err := p.parseIdentList()
	if err != nil {
		return nil, err
	}
	cmd.PartitionNames = names

	// INTO (partition_definitions)
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		if p.cur.Type == ')' {
			break
		}
		pd, err := p.parsePartitionDef()
		if err != nil {
			return nil, err
		}
		cmd.PartitionDefs = append(cmd.PartitionDefs, pd)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseAlterExchangePartition parses EXCHANGE PARTITION partition_name WITH TABLE tbl_name [{WITH | WITHOUT} VALIDATION].
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html
//
//	EXCHANGE PARTITION partition_name WITH TABLE tbl_name [{WITH | WITHOUT} VALIDATION]
func (p *Parser) parseAlterExchangePartition(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume EXCHANGE
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	cmd.Type = nodes.ATExchangePartition

	// Partition name
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cmd.Name = name

	// WITH TABLE tbl_name
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	tblRef, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	cmd.ExchangeTable = tblRef

	// Optional [{WITH | WITHOUT} VALIDATION]
	if p.cur.Type == kwWITH {
		p.advance()
		p.match(kwVALIDATION)
		v := true
		cmd.WithValidation = &v
	} else if p.cur.Type == kwWITHOUT {
		p.advance()
		p.match(kwVALIDATION)
		v := false
		cmd.WithValidation = &v
	}

	cmd.Loc.End = p.pos()
	return cmd, nil
}

// parseColumnPositioning parses optional FIRST | AFTER col_name.
func (p *Parser) parseColumnPositioning(cmd *nodes.AlterTableCmd) error {
	if _, ok := p.match(kwFIRST); ok {
		cmd.First = true
	} else if _, ok := p.match(kwAFTER); ok {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		cmd.After = name
	}
	return nil
}
