package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
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
func (p *Parser) parseAlterTableStmt(start int) (*nodes.AlterTableStmt, error) {
	p.advance() // consume TABLE
	return p.parseAlterTableBody(&nodes.AlterTableStmt{Loc: nodes.Loc{Start: start}})
}

// parseAlterExternalTableStmt parses ALTER EXTERNAL TABLE, after EXTERNAL. Its
// body is ALTER TABLE's, plus two actions of its own that readAlterTableBody
// takes when the table is external.
//
// Ref: alter_table_stmt's alter_with_opt_hint EXTERNAL TABLE alternatives in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, line 19244).
func (p *Parser) parseAlterExternalTableStmt(start int) (*nodes.AlterTableStmt, error) {
	p.advance() // consume TABLE
	return p.parseAlterTableBody(&nodes.AlterTableStmt{
		Loc:      nodes.Loc{Start: start},
		External: true,
	})
}

// parseAlterTableBody parses what follows ALTER [EXTERNAL] TABLE: the table's
// name, then its actions.
func (p *Parser) parseAlterTableBody(stmt *nodes.AlterTableStmt) (*nodes.AlterTableStmt, error) {
	// Completion: after ALTER TABLE, offer table_ref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("table_ref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Table name
	ref, err := p.parseRelationFactor()
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

	// alter_table_action_list reads `alter_table_actions ',' alter_table_action`
	// with an empty actions arm, so one comma may lead the first action:
	// `ALTER TABLE a, NOPARALLEL` is one table option.
	// The external table's own alternatives take their action directly after the
	// table's name, so no comma leads them.
	if p.cur.Type == ',' && !stmt.External {
		p.advance()
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			// A comma with no action after it is no action list.
			return nil, p.syntaxErrorAtCur()
		}
	}

	// Parse comma-separated alter operations
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}

		// ALTER EXTERNAL TABLE's two own actions, before the generic ones:
		// ADD PARTITION (...) LOCATION 's' and DROP PARTITION LOCATION 's'.
		if stmt.External && (p.cur.Type == kwADD || p.cur.Type == kwDROP) {
			cmd, ok, err := p.parseAlterExternalPartition()
			if err != nil {
				return nil, err
			}
			if ok {
				// The external partition forms are whole statements, not entries
				// of alter_table_actions: no comma follows them.
				if p.cur.Type == ',' {
					return nil, p.syntaxErrorAtCur()
				}
				stmt.Commands = append(stmt.Commands, cmd)
				break
			}
		}

		if p.cur.Type == kwSET {
			cmds, err := p.parseAlterTableSet()
			if err != nil {
				return nil, err
			}
			stmt.Commands = append(stmt.Commands, cmds...)
		} else {
			cmd, err := p.parseAlterTableCmd()
			if err != nil {
				return nil, err
			}
			stmt.Commands = append(stmt.Commands, cmd)
			// {ADD | DROP} COLUMN GROUP is a statement of its own, alter_table_stmt's
			// alter_column_group_option, and no alter_table_action: it stands with no
			// other action before or after it (sql_parser_mysql_mode.y,
			// oceanbase v4.3.5_CE_BP6).
			if cmd.Type == nodes.ATAddColumnGroup || cmd.Type == nodes.ATDropColumnGroup {
				if len(stmt.Commands) > 1 {
					return nil, &ParseError{Message: "syntax error", Position: cmd.Loc.Start}
				}
				if p.cur.Type == ',' {
					return nil, p.syntaxErrorAtCur()
				}
			}
			if cmd.Type == nodes.ATTableOption {
				more, err := p.parseMoreAlterTableOptions()
				if err != nil {
					return nil, err
				}
				stmt.Commands = append(stmt.Commands, more...)
			}
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		// A comma the list never got an action for ends the statement with
		// nothing after it: OceanBase CE 4.3.5.6 answers 1064 for
		// `ALTER TABLE t ENGINE = InnoDB,` and for
		// `ALTER TABLE t ADD COLUMN a INT,`.
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTableSet parses OceanBase's `SET table_option [table_option ...]`, one
// ALTER TABLE action. alter_table_action in sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6) reads `opt_set table_option_list_space_seperated`: SET may lead
// the options, which follow one another without commas, with or without it. 4.3.5.6
// applies `ALTER TABLE t SET COMMENT 'x'` as `ALTER TABLE t COMMENT 'x'`, and reads
// back each option set this way (alterTableSetCases). SET with no table option after
// it is a 1064. Each option becomes its own ATTableOption command, the first one
// spanning the SET.
func (p *Parser) parseAlterTableSet() ([]*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume SET
	opt, ok, err := p.parseTableOption()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	cmds := []*nodes.AlterTableCmd{{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Type:   nodes.ATTableOption,
		Option: opt,
	}}
	more, err := p.parseMoreAlterTableOptions()
	if err != nil {
		return nil, err
	}
	return append(cmds, more...), nil
}

// parseMoreAlterTableOptions parses the table options that follow one in the same
// ALTER TABLE action, without commas between them: `ALTER TABLE t PCTFREE = 5 COMMENT
// 'x'`, which OceanBase reads as table_option_list_space_seperated, with or without a
// leading SET.
func (p *Parser) parseMoreAlterTableOptions() ([]*nodes.AlterTableCmd, error) {
	var cmds []*nodes.AlterTableCmd
	for {
		start := p.pos()
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if !ok {
			return cmds, nil
		}
		cmds = append(cmds, &nodes.AlterTableCmd{
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
			Type:   nodes.ATTableOption,
			Option: opt,
		})
	}
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
		// Known gap: OceanBase refuses a backquoted value, ALGORITHM=`DEFAULT`
		// (oceanbase/review's TestColumnDropAgainstOceanBase records the 1064),
		// and this parser still takes it. Refusing it changes what that test's
		// walk sends, so it waits for the test to be recorded again.
		if _, ok := p.curBareWordOf(alterAlgorithmWords); !ok && p.backtickQuotedAt(p.cur.Loc) {
			for _, w := range alterAlgorithmWords {
				if strings.EqualFold(p.cur.Str, w) {
					cmd.Name = p.advance().Str
					cmd.Loc.End = p.prev.End
					return cmd, nil
				}
			}
		}
		val, err := p.parseWordOf(alterAlgorithmWords)
		if err != nil {
			return nil, err
		}
		cmd.Name = val
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwLOCK:
		p.advance()
		p.match('=')
		cmd.Type = nodes.ATLock
		val, err := p.parseWordOf(alterLockWords)
		if err != nil {
			return nil, err
		}
		cmd.Name = val
		cmd.Loc.End = p.prev.End
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
		if _, ok := p.match(kwSUBPARTITION); ok {
			cmd.Type = nodes.ATTruncateSubpartition
			names, err := p.parseNameList()
			if err != nil {
				return nil, err
			}
			cmd.PartitionNames = names
			cmd.Loc.End = p.prev.End
			return cmd, nil
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
			if _, err := p.expect(kwTABLESPACE); err != nil {
				return nil, err
			}
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// DISCARD TABLESPACE
		if _, err := p.expect(kwTABLESPACE); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATDiscardTablespace
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwIMPORT:
		// IMPORT TABLESPACE | IMPORT PARTITION {partition_names | ALL} TABLESPACE
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd, err := p.parseAlterPartitionNamesOrAll(cmd, nodes.ATImportPartitionTablespace)
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(kwTABLESPACE); err != nil {
				return nil, err
			}
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// IMPORT TABLESPACE
		if _, err := p.expect(kwTABLESPACE); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATImportTablespace
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwFORCE:
		p.advance()
		cmd.Type = nodes.ATForce
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwORDER:
		// ORDER BY col_name [ASC | DESC] [, ...]
		p.advance()
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATOrderBy
		items, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		cmd.OrderByItems = items
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwENABLE:
		// ENABLE KEYS
		p.advance()
		if _, err := p.expect(kwKEYS); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATEnableKeys
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwDISABLE:
		// DISABLE KEYS
		p.advance()
		if _, err := p.expect(kwKEYS); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATDisableKeys
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwWITH:
		// WITH VALIDATION
		p.advance()
		if _, err := p.expect(kwVALIDATION); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATWithValidation
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwWITHOUT:
		// WITHOUT VALIDATION
		p.advance()
		if _, err := p.expect(kwVALIDATION); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATWithoutValidation
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwPARTITION:
		// PARTITION BY ... (repartition)
		if p.peekNext().Type == kwBY {
			part, err := p.parsePartitionClause(false)
			if err != nil {
				return nil, err
			}
			cmd.Type = nodes.ATPartitionBy
			cmd.PartitionBy = part
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		return nil, &ParseError{
			Message:  "expected BY after PARTITION",
			Position: p.cur.Loc,
		}

	case kwREMOVE:
		// REMOVE PARTITIONING, or OceanBase's REMOVE TTL.
		p.advance()
		if p.isBareWord("TTL") {
			p.advance()
			cmd.Type = nodes.ATRemoveTTL
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		if _, err := p.expect(kwPARTITIONING); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATRemovePartitioning
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwSECONDARY_LOAD:
		p.advance()
		cmd.Type = nodes.ATSecondaryLoad
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwSECONDARY_UNLOAD:
		p.advance()
		cmd.Type = nodes.ATSecondaryUnload
		cmd.Loc.End = p.prev.End
		return cmd, nil

	default:
		// alter_table_action's REFRESH refreshes an external table's metadata.
		// REFRESH is no keyword of this lexer's, so it is a bare word here.
		if p.isBareWord("REFRESH") {
			p.advance()
			cmd.Type = nodes.ATRefresh
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// Try table options: ENGINE, CHARSET, etc.
		opt, ok, err := p.parseTableOption()
		if err != nil {
			return nil, err
		}
		if ok {
			cmd.Type = nodes.ATTableOption
			cmd.Option = opt
			cmd.Loc.End = p.prev.End
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

	// alter_index_option's parenthesized arms: ADD (KEY (...)), ADD (UNIQUE
	// (...)), ADD (CONSTRAINT ... ), ADD (SPATIAL (...)), ADD (PRIMARY KEY
	// (...)). They follow ADD with no COLUMN, and only those keywords open one:
	// `ALTER TABLE t ADD COLUMN (INDEX (a))` and `ADD (CHECK (true))` are 1064s.
	if p.cur.Type == '(' && p.parenOpensAddIndex() {
		p.advance()
		inner, err := p.parseAlterAddIndexAction(cmd)
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		inner.Loc.End = p.prev.End
		return inner, nil
	}

	// ADD COLUMN GROUP (...) [DELAYED]. GROUP is reserved, so no column is called
	// it.
	if p.cur.Type == kwCOLUMN && p.peekNext().Type == kwGROUP {
		return p.parseAlterColumnGroup(cmd, nodes.ATAddColumnGroup)
	}

	// ADD PARTITION (partition_definition, ...). MySQL also takes ADD PARTITION
	// PARTITIONS n, and partitions without VALUES for a HASH or KEY table. OceanBase's
	// ADD PARTITION reads opt_partition_range_or_list alone (v4.3.5_CE_BP6), so
	// 4.3.5.6 answers 1064 for both, at PARTITIONS and at the token after the
	// partition's name (partitionOptionCases).
	if p.cur.Type == kwPARTITION {
		p.advance()
		cmd.Type = nodes.ATAddPartition
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		for {
			pd, err := p.parsePartitionDef(true)
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
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}

	// ADD [CONSTRAINT ...] PRIMARY KEY / UNIQUE / FOREIGN KEY / CHECK
	if p.isTableConstraintStart() {
		cmd.Type = nodes.ATAddConstraint
		constr, err := p.parseTableConstraint(true)
		if err != nil {
			return nil, err
		}
		cmd.Constraint = constr
		cmd.Loc.End = p.prev.End
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
		if err := p.parseLobStorageClause(cmd); err != nil {
			return nil, err
		}
		cmd.Loc.End = p.prev.End
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
	if err := p.parseLobStorageClause(cmd); err != nil {
		return nil, err
	}

	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseLobStorageClause reads lob_storage_clause where an ALTER TABLE ADD may
// write it: JSON '(' column_name ')' STORE AS '(' lob_storage_parameters ')',
// whose parameters are lob_chunk_size, `CHUNK n` or `CHUNK 'string'`.
func (p *Parser) parseLobStorageClause(cmd *nodes.AlterTableCmd) error {
	if !p.isBareWord("JSON") {
		return nil
	}
	start := p.pos()
	p.advance() // consume JSON
	if _, err := p.expect('('); err != nil {
		return err
	}
	column, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	if !p.isBareWord("STORE") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if _, err := p.expect(kwAS); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	clause := &nodes.LobStorageClause{Loc: nodes.Loc{Start: start}, Column: column}
	if !p.isBareWord("CHUNK") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != tokICONST && p.cur.Type != tokSCONST {
		return p.syntaxErrorAtCur()
	}
	clause.Chunk = p.cur.Str
	p.advance()
	// lob_storage_parameters is one or more parameters written next to each
	// other; CHUNK is the only one the grammar has.
	for p.isBareWord("CHUNK") {
		p.advance()
		if p.cur.Type != tokICONST && p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		clause.Chunk = p.cur.Str
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	clause.Loc.End = p.prev.End
	cmd.LobStorage = clause
	return nil
}

// parenOpensAddIndex reports whether the '(' at p.cur opens one of
// alter_index_option's parenthesized ADD arms.
func (p *Parser) parenOpensAddIndex() bool {
	if p.cur.Type != '(' {
		return false
	}
	// alter_index_option's parenthesized ADD arms are key_or_index, UNIQUE,
	// CONSTRAINT, SPATIAL and PRIMARY KEY. FULLTEXT is no arm of them, so
	// `ALTER TABLE t ADD (FULLTEXT KEY (a))` is a 1064.
	switch p.peekNext().Type {
	case kwKEY, kwINDEX, kwUNIQUE, kwPRIMARY, kwSPATIAL, kwCONSTRAINT:
		return true
	}
	return false
}

// parseAlterAddIndexAction parses the ADD action inside one of
// alter_index_option's parenthesized arms, with the '(' already consumed.
// p.cur is the constraint or index keyword.
func (p *Parser) parseAlterAddIndexAction(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	if p.isTableConstraintStart() {
		cmd.Type = nodes.ATAddConstraint
		constr, err := p.parseTableConstraint(true)
		if err != nil {
			return nil, err
		}
		cmd.Constraint = constr
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseAlterColumnGroup parses {ADD | DROP} COLUMN GROUP ( column_group_list )
// [DELAYED] at COLUMN.
//
// Ref: alter_column_group_option and alter_column_group_delayed_desc in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, line 19967).
func (p *Parser) parseAlterColumnGroup(cmd *nodes.AlterTableCmd, typ nodes.AlterTableCmdType) (*nodes.AlterTableCmd, error) {
	p.advance() // consume COLUMN
	p.advance() // consume GROUP
	list, err := p.parseColumnGroupList()
	if err != nil {
		return nil, err
	}
	cmd.Type, cmd.Name = typ, list
	if p.isBareWord("DELAYED") {
		start := p.pos()
		p.advance()
		cmd.Option = &nodes.TableOption{Loc: nodes.Loc{Start: start, End: p.prev.End}, Name: "DELAYED"}
	}
	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseAlterDrop parses DROP operations.
//
// Ref: the DROP arms of alter_list_item in MySQL 8.0.46 sql/sql_yacc.yy:
// DROP opt_column ident opt_restrict, DROP FOREIGN KEY ident, DROP PRIMARY
// KEY, DROP key_or_index ident, DROP CHECK ident, DROP CONSTRAINT ident. KEY
// is required after PRIMARY and FOREIGN, and only a column takes RESTRICT or
// CASCADE. No arm takes IF EXISTS, which is MariaDB syntax.
func (p *Parser) parseAlterDrop(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume DROP

	// DROP COLUMN GROUP (...) [DELAYED].
	if p.cur.Type == kwCOLUMN && p.peekNext().Type == kwGROUP {
		return p.parseAlterColumnGroup(cmd, nodes.ATDropColumnGroup)
	}

	// Completion: after DROP, offer what can be dropped.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwCOLUMN, kwINDEX, kwKEY, kwFOREIGN, kwPRIMARY, kwCHECK, kwCONSTRAINT, kwPARTITION} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	switch p.cur.Type {
	case kwTABLEGROUP:
		// DROP TABLEGROUP takes the table out of its table group.
		// Ref: alter_tablegroup_option in sql_parser_mysql_mode.y (oceanbase
		// v4.3.5_CE_BP6, line 20096). TABLEGROUP is reserved, so no column has the name.
		p.advance()
		cmd.Type = nodes.ATDropTablegroup
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwPARTITION, kwSUBPARTITION:
		// {DROP PARTITION | DROP SUBPARTITION} drop_partition_name_list, which is
		// name_list or name_list in parentheses. SUBPARTITION is unreserved, so
		// with no name after it, it is the column DROP drops.
		if next := p.peekNext().Type; p.cur.Type == kwSUBPARTITION &&
			(next == tokEOF || next == ';' || next == ',' || next == kwRESTRICT || next == kwCASCADE) {
			return p.parseAlterDropColumn(cmd)
		}
		cmd.Type = nodes.ATDropPartition
		if p.cur.Type == kwSUBPARTITION {
			cmd.Type = nodes.ATDropSubpartition
		}
		p.advance()
		paren := p.cur.Type == '('
		if paren {
			p.advance()
		}
		names, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		if paren {
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
		cmd.PartitionNames = names
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwPRIMARY:
		// DROP PRIMARY KEY
		p.advance()
		if err := p.expectDropKey(); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATDropConstraint
		cmd.ConstraintKeyword = nodes.ConstraintKeywordPrimaryKey
		cmd.Name = "PRIMARY"
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwINDEX, kwKEY:
		// DROP {INDEX | KEY} index_name. OceanBase 4.3.5 does not take IF EXISTS
		// here, it answers 1064, so the parser must not accept it either: see
		// foreignGrammar and TestOceanBaseRejectsForeignGrammar.
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwFOREIGN:
		// DROP FOREIGN KEY fk_symbol
		p.advance()
		if err := p.expectDropKey(); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATDropConstraint
		cmd.ConstraintKeyword = nodes.ConstraintKeywordForeignKey
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwCHECK:
		// DROP CHECK symbol, or alter_constraint_option's DROP CHECK '(' name_list ')'.
		p.advance()
		cmd.Type = nodes.ATDropConstraint
		cmd.ConstraintKeyword = nodes.ConstraintKeywordCheck
		if p.cur.Type == '(' {
			names, err := p.parseParenNameList()
			if err != nil {
				return nil, err
			}
			cmd.ConstraintNames = names
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwCONSTRAINT:
		// DROP CONSTRAINT symbol, or alter_constraint_option's
		// DROP CONSTRAINT '(' name_list ')'.
		p.advance()
		cmd.Type = nodes.ATDropConstraint
		cmd.ConstraintKeyword = nodes.ConstraintKeywordConstraint
		if p.cur.Type == '(' {
			names, err := p.parseParenNameList()
			if err != nil {
				return nil, err
			}
			cmd.ConstraintNames = names
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	default:
		return p.parseAlterDropColumn(cmd)
	}
}

// parseAlterDropColumn parses DROP [COLUMN] col_name [RESTRICT | CASCADE] after
// DROP.
func (p *Parser) parseAlterDropColumn(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	cmd.Type = nodes.ATDropColumn
	p.match(kwCOLUMN)
	// Completion: after DROP COLUMN, offer columnref.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}
	name, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	cmd.Name = name
	// opt_restrict: the grammar action ignores it.
	p.match(kwRESTRICT, kwCASCADE)
	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseParenNameList parses '(' name_list ')', the form alter_constraint_option
// writes for DROP CONSTRAINT and DROP CHECK.
func (p *Parser) parseParenNameList() ([]string, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	names, err := p.parseNameList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// expectDropKey consumes the KEY that DROP PRIMARY and DROP FOREIGN require.
func (p *Parser) expectDropKey() error {
	// Completion: after DROP PRIMARY or DROP FOREIGN, offer KEY.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwKEY)
		return &ParseError{Message: "collecting"}
	}
	_, err := p.expect(kwKEY)
	return err
}

// parseAlterModify parses MODIFY [COLUMN] col_name column_definition [FIRST | AFTER col_name].
func (p *Parser) parseAlterModify(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume MODIFY

	// MODIFY PARTITION relation_factor ADD SUBPARTITION ( subpartitions ), whose
	// subpartitions are all RANGE or all LIST: each has VALUES (alter_partition_option
	// and opt_subpartition_range_or_list, sql_parser_mysql_mode.y, oceanbase
	// v4.3.5_CE_BP6).
	if p.cur.Type == kwPARTITION {
		p.advance()
		cmd.Type = nodes.ATAddSubpartition
		ref, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		cmd.Name = ref.Name
		if _, err := p.expect(kwADD); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwSUBPARTITION); err != nil {
			return nil, err
		}
		pd := &nodes.PartitionDef{Loc: nodes.Loc{Start: p.pos()}, Name: ref.Name}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		kind := noPartitionValues
		for i := 0; ; i++ {
			if i == 0 && p.cur.Type == kwSUBPARTITION && p.peekNext().Type != tokEOF {
				// The first entry sets the kind, and HASH, which has no VALUES, is
				// not one of the two.
				save := p.saveScan()
				p.advance()
				if _, _, err := p.parseIdent(); err == nil && p.cur.Type != kwVALUES {
					return nil, p.syntaxErrorAtCur()
				}
				p.restoreScan(save)
			}
			spd, got, err := p.parseSubPartitionDef(kind, i == 0)
			if err != nil {
				return nil, err
			}
			kind = got
			pd.SubPartitions = append(pd.SubPartitions, spd)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		pd.Loc.End = p.prev.End
		cmd.PartitionDefs = []*nodes.PartitionDef{pd}
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}
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

	cmd.Loc.End = p.prev.End
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

	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseAlterColumn parses ALTER [COLUMN] col_name {SET DEFAULT | DROP DEFAULT}
// and ALTER INDEX index_name {VISIBLE | INVISIBLE}.
//
// ALTER COLUMN has no nullability: OceanBase CE 4.3.5.6 answers 1064 for `ALTER [COLUMN] c
// SET NOT NULL` and `DROP NOT NULL`, which are PostgreSQL's; MODIFY and CHANGE
// are how a column's nullability changes. ATAlterColumnSetNotNull and
// ATAlterColumnDropNotNull are never built.
//
// ALTER COLUMN has no visibility either: SET VISIBLE and SET INVISIBLE are
// MySQL's, and OceanBase CE 4.3.5.6 answers 1064 for both, though it takes the
// VISIBLE and INVISIBLE column attributes in CREATE TABLE, ADD, MODIFY and
// CHANGE, which are how a column's visibility changes. ATAlterColumnVisible
// and ATAlterColumnInvisible are never built.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/alter-table.html, and the
// `ALTER opt_column ident` arms of alter_list_item in MySQL 8.0.46
// sql/sql_yacc.yy: SET DEFAULT, DROP DEFAULT and SET visibility, of which
// OceanBase takes the first two. OceanBase's own grammar is not in this tree;
// the verdicts above are its.
//
//	ALTER [COLUMN] col_name {
//	    SET DEFAULT {literal | (expr)}
//	  | DROP DEFAULT
//	}
//	ALTER INDEX index_name {VISIBLE | INVISIBLE}
func (p *Parser) parseAlterColumn(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume ALTER

	// {ALTER PARTITION | ALTER SUBPARTITION} name_list STORAGE_CACHE_POLICY = 'policy'
	// (alter_partition_option, sql_parser_mysql_mode.y, oceanbase v4.3.5_CE_BP6).
	// SUBPARTITION is unreserved, so before SET or DROP it is the column ALTER
	// changes the default of.
	if next := p.peekNext().Type; p.cur.Type == kwPARTITION ||
		(p.cur.Type == kwSUBPARTITION && next != kwSET && next != kwDROP) {
		cmd.Type = nodes.ATAlterPartitionStorageCachePolicy
		cmd.Name = strings.ToUpper(p.advance().Str)
		names, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		cmd.PartitionNames = names
		optStart := p.pos()
		if !p.isBareWord("STORAGE_CACHE_POLICY") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		var policy strings.Builder
		for p.cur.Type == tokSCONST {
			policy.WriteString(p.advance().Str)
		}
		cmd.Option = &nodes.TableOption{Name: "STORAGE_CACHE_POLICY", Value: policy.String()}
		cmd.Option.Loc = nodes.Loc{Start: optStart, End: p.prev.End}
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}

	// alter_index_option's ALTER INDEX arms: the index's name is required, and
	// one of a visibility option, STORAGE_CACHE_POLICY (...) and
	// PARALLEL n | NOPARALLEL follows it. `ALTER TABLE a ALTER INDEX a` and
	// `ALTER TABLE a ALTER INDEX NOPARALLEL` are 1064s.
	//
	// Ref: alter_index_option in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
	// 19708-19728).
	if p.cur.Type == kwINDEX {
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		switch {
		case p.cur.Type == kwVISIBLE:
			p.advance()
			cmd.Type = nodes.ATAlterIndexVisible
		case p.cur.Type == kwINVISIBLE:
			p.advance()
			cmd.Type = nodes.ATAlterIndexInvisible
		case p.isBareWord("STORAGE_CACHE_POLICY"):
			opt, ok, err := p.parseStorageCachePolicy()
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			cmd.Type = nodes.ATAlterIndexStorageCachePolicy
			cmd.IndexOptions = []*nodes.IndexOption{opt}
		default:
			par, ok, err := p.parseParallelOption()
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			cmd.Type = nodes.ATAlterIndexParallel
			cmd.Parallel = par
		}
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}

	// ALTER {CHECK | CONSTRAINT} symbol {ENFORCED | NOT ENFORCED}
	//
	// The enforcement is required, not optional: OceanBase CE 4.3.5.6
	// answers 1064 for `ALTER TABLE t ALTER CHECK ident` and for
	// `ALTER TABLE t ALTER CONSTRAINT ident` with nothing after the name.
	// CHECK and CONSTRAINT are the same clause spelled two ways.
	//
	// Ref: `ALTER check_constraint_keyword ident constraint_enforcement` in
	// MySQL 8.0.46 sql/sql_yacc.yy, whose constraint_enforcement has no
	// empty arm.
	if p.cur.Type == kwCHECK || p.cur.Type == kwCONSTRAINT {
		if p.cur.Type == kwCHECK {
			cmd.ConstraintKeyword = nodes.ConstraintKeywordCheck
		} else {
			cmd.ConstraintKeyword = nodes.ConstraintKeywordConstraint
		}
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		cmd.Type = nodes.ATAlterCheckEnforced
		if _, ok := p.match(kwNOT); ok {
			if _, err := p.expect(kwENFORCED); err != nil {
				return nil, err
			}
			cmd.NewName = "NOT ENFORCED"
		} else {
			if _, err := p.expect(kwENFORCED); err != nil {
				return nil, err
			}
			cmd.NewName = "ENFORCED"
		}
		cmd.Loc.End = p.prev.End
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
			// SET DEFAULT takes a signed literal, or any expression inside
			// parentheses, and nothing in between: OceanBase CE 4.3.5.6 answers
			// 1064 for `SET DEFAULT CURRENT_TIMESTAMP`, `SET DEFAULT NOW()` and
			// `SET DEFAULT 1+1`, and runs `SET DEFAULT (CURRENT_TIMESTAMP)`,
			// `SET DEFAULT (1+1)` and `SET DEFAULT -1`.
			//
			// Ref: `ALTER opt_column ident SET DEFAULT signed_literal` and
			// `... SET DEFAULT '(' expr ')'` in MySQL 8.0.46
			// sql/sql_yacc.yy; signed_literal reads
			// `literal | '+' NUM_literal | '-' NUM_literal`.
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
				expr, err := p.parseSignedLiteral()
				if err != nil {
					return nil, err
				}
				cmd.DefaultExpr = expr
			}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	} else if _, ok := p.match(kwDROP); ok {
		if _, ok := p.match(kwDEFAULT); ok {
			cmd.Type = nodes.ATAlterColumnDefault
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	} else {
		// A column name with no action is not a clause: OceanBase CE 4.3.5.6 answers
		// 1064, as it does for a SET or DROP with nothing it takes after it.
		return nil, p.syntaxErrorAtCur()
	}

	cmd.Loc.End = p.prev.End
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
		// TO is required: OceanBase CE 4.3.5.6 answers 1064 for
		// `RENAME COLUMN a b` and for `RENAME INDEX a b`.
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
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
		// TO is required: OceanBase CE 4.3.5.6 answers 1064 for
		// `RENAME COLUMN a b` and for `RENAME INDEX a b`.
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		cmd.NewName = newName

	default:
		// RENAME [TO | AS] new_tbl_name
		//
		// No equals sign: OceanBase CE 4.3.5.6 answers 1064 for
		// `ALTER TABLE t RENAME = t1`, where MySQL 8.0.46 and MariaDB 11.8.8
		// take it, so the separator is not widened here.
		cmd.Type = nodes.ATRenameTable
		p.match(kwTO, kwAS)
		newTable, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		cmd.NewTable = newTable
	}

	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseAlterConvert parses CONVERT TO CHARACTER SET charset_name [COLLATE collation_name].
func (p *Parser) parseAlterConvert(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume CONVERT
	cmd.Type = nodes.ATConvertCharset

	// CONVERT TO CHARACTER SET or CHAR SET is required, and so is the
	// charset after it: OceanBase CE 4.3.5.6 answers 1064 for
	// `CONVERT TO utf8`, for
	// `CONVERT CHARSET utf8` and for a `CONVERT TO` keyword with no charset.
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if err := p.expectCharsetKw(); err != nil {
		return nil, err
	}

	// Completion: after CONVERT TO CHARACTER SET, offer charset candidates.
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("charset")
		return nil, &ParseError{Message: "collecting"}
	}

	// DEFAULT is no charset name here: OceanBase CE 4.3.5.6 answers 1064 for
	// `CONVERT TO CHARACTER SET DEFAULT`, which MySQL 8.0.46, MariaDB 11.8.8
	// and TiDB v8.5.5 all run. consumeOptionValue takes any keyword, so the
	// word has to be refused before it.
	if p.cur.Type == kwDEFAULT {
		return nil, p.syntaxErrorAtCur()
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

	cmd.Loc.End = p.prev.End
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
//
// MySQL's opt_no_write_to_binlog -- the NO_WRITE_TO_BINLOG or LOCAL modifier
// after the PARTITION keyword, which mysql, mariadb and tidb take on ADD,
// REBUILD, OPTIMIZE, REPAIR, ANALYZE, COALESCE and REORGANIZE PARTITION -- is
// deliberately not here, and neither is the trailing partitioning clause those
// three read as `alter_list remove_partitioning` and `alter_list
// opt_partitioning`. OceanBase CE 4.3.5.6 answers 1064 to the whole surface
// they belong to: OPTIMIZE, REPAIR, REBUILD, ANALYZE and COALESCE PARTITION
// are not actions it has at all, REMOVE PARTITIONING is a 1064 standing alone
// or trailing, and PARTITION BY after another action is a 1064 while `ALTER
// TABLE t PARTITION BY HASH(a), ADD COLUMN a INT` -- which the other three
// reject -- parses. Adding the modifier to actions the engine does not have
// would only widen an over-acceptance; tightening that surface to what 4.3.5.6
// carries is its own change.
func (p *Parser) parseAlterPartitionNamesOrAll(cmd *nodes.AlterTableCmd, cmdType nodes.AlterTableCmdType) (*nodes.AlterTableCmd, error) {
	cmd.Type = cmdType
	if p.cur.Type == kwALL {
		p.advance()
		cmd.AllPartitions = true
	} else {
		names, err := p.parseNameList()
		if err != nil {
			return nil, err
		}
		cmd.PartitionNames = names
	}
	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseNameList parses OceanBase's name_list:
//
//	name_list: NAME_OB | unreserved_keyword | name_list [','] NAME_OB
//
// The first name may be an unreserved keyword; each one after it is NAME_OB,
// with a comma before it or without, so `p0 p1` is two names and a comma that
// no NAME_OB follows is left for the next action.
//
// Ref: name_list in sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6).
func (p *Parser) parseNameList() ([]string, error) {
	first, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	names := []string{first}
	for {
		if p.isNameOB(p.cur) {
			names = append(names, p.advance().Str)
			continue
		}
		if p.cur.Type == ',' && p.isNameOB(p.peekNext()) {
			p.advance()
			names = append(names, p.advance().Str)
			continue
		}
		return names, nil
	}
}

// isNameOB reports whether tok is NAME_OB: an identifier that is no keyword of
// OceanBase's, or any backquoted name.
func (p *Parser) isNameOB(tok Token) bool {
	if tok.Type != tokIDENT && tok.Type < 700 {
		return false
	}
	if p.backtickQuotedAt(tok.Loc) {
		return true
	}
	return !oceanBaseKeywords[strings.ToLower(tok.Str)] && !strings.HasPrefix(tok.Str, "@") && !isCharsetIntroducer(tok.Str)
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
	cmd.Loc.End = p.prev.End
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

	// INTO (partition_definitions). OceanBase reads opt_partition_range_or_list here
	// too (v4.3.5_CE_BP6), so a partition without VALUES is a 1064, as in ADD PARTITION.
	if _, err := p.expect(kwINTO); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	for {
		pd, err := p.parsePartitionDef(true)
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
	cmd.Loc.End = p.prev.End
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
	tblRef, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	cmd.ExchangeTable = tblRef

	// Optional [{WITH | WITHOUT} VALIDATION]
	if p.cur.Type == kwWITH {
		p.advance()
		if _, err := p.expect(kwVALIDATION); err != nil {
			return nil, err
		}
		v := true
		cmd.WithValidation = &v
	} else if p.cur.Type == kwWITHOUT {
		p.advance()
		if _, err := p.expect(kwVALIDATION); err != nil {
			return nil, err
		}
		v := false
		cmd.WithValidation = &v
	}

	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseColumnPositioning parses optional FIRST | AFTER col_name.
func (p *Parser) parseColumnPositioning(cmd *nodes.AlterTableCmd) error {
	if _, ok := p.match(kwFIRST); ok {
		cmd.First = true
	} else if _, ok := p.match(kwAFTER); ok {
		// opt_position_column's AFTER column_name, and BEFORE column_name
		// beside it, which ALTER TABLE ... ADD takes.
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		cmd.After = name
	} else if p.isBareWord("BEFORE") {
		// opt_position_column's other arm, which places the column before the one
		// it names where AFTER places it after.
		p.advance()
		name, _, err := p.parseIdentifier()
		if err != nil {
			return err
		}
		cmd.Before = name
	}
	return nil
}

// expectCharsetKw consumes the keyword that introduces a charset name. The
// two-word spellings need both words.
func (p *Parser) expectCharsetKw() error {
	switch p.cur.Type {
	// CHARSET as one word is not a spelling OceanBase CE 4.3.5.6 takes: it answers
	// 1064 for `CONVERT TO CHARSET utf8` and runs
	// `CONVERT TO CHAR SET utf8`.
	case kwCHARACTER, kwCHAR:
		p.advance()
		_, err := p.expect(kwSET)
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseSignedLiteral parses the signed_literal a SET DEFAULT takes outside
// parentheses: one literal, optionally signed. A function call, a keyword such
// as CURRENT_TIMESTAMP and an arithmetic expression are all syntax errors
// here, which is why this is not parseExpr.
//
// Ref: signed_literal in MySQL 8.0.46 sql/sql_yacc.yy, which reads
// `literal | '+' NUM_literal | '-' NUM_literal`.
func (p *Parser) parseSignedLiteral() (nodes.ExprNode, error) {
	start := p.pos()
	if tok, ok := p.match('+', '-'); ok {
		// Only a number may be signed.
		if p.cur.Type != tokICONST && p.cur.Type != tokFCONST {
			return nil, p.syntaxErrorAtCur()
		}
		operand, err := p.parseLiteralOnly()
		if err != nil {
			return nil, err
		}
		op := nodes.UnaryPlus
		opText := "+"
		if tok.Type == '-' {
			op = nodes.UnaryMinus
			opText = "-"
		}
		return &nodes.UnaryExpr{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			Op:         op,
			Operand:    operand,
			OriginalOp: opText,
		}, nil
	}
	return p.parseLiteralOnly()
}

// parseLiteralOnly parses one literal: a number, a string with or without a
// charset introducer, a hex or bit literal, NULL, TRUE or FALSE, and a
// temporal literal.
// Adjacent string literals concatenate, which parsePrimaryExpr already does.
func (p *Parser) parseLiteralOnly() (nodes.ExprNode, error) {
	start := p.pos()
	switch p.cur.Type {
	case tokICONST, tokFCONST, tokSCONST, tokBCONST, tokXCONST:
		expr, err := p.parsePrimaryExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return expr, nil
	case kwNULL:
		p.advance()
		return &nodes.NullLit{Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
	case kwTRUE, kwFALSE:
		value := p.cur.Type == kwTRUE
		p.advance()
		return &nodes.BoolLit{Loc: nodes.Loc{Start: start, End: p.prev.End}, Value: value}, nil
	case kwDATE, kwTIME, kwTIMESTAMP:
		// A temporal literal, which OceanBase CE 4.3.5.6 takes here:
		// `SET DEFAULT DATE '2020-01-01'` runs. The word is only a literal
		// when a string follows it.
		if p.peekNext().Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		return p.parsePrimaryExpr()
	}
	// A charset introducer followed by a string, as in _utf8'x', is a literal
	// too, and the introducer reaches the parser as an identifier.
	if p.cur.Type == tokIDENT && strings.HasPrefix(p.cur.Str, "_") && p.peekNext().Type == tokSCONST {
		return p.parsePrimaryExpr()
	}
	return nil, p.syntaxErrorAtCur()
}

// parseAlterExternalPartition reads an ALTER EXTERNAL TABLE action of its own,
// at ADD or DROP, and reports whether one stood there:
//
//	ADD PARTITION '(' add_external_table_partition_actions ')' LOCATION 's'
//	DROP PARTITION LOCATION 's'
//
// The actions list may be empty or open with an empty entry: `ADD PARTITION()`
// and `ADD PARTITION(, ds = 'x')` both parse. A DROP PARTITION that names a
// partition called LOCATION is the generic action and not this one, so the
// action is read here only when a string follows LOCATION.
//
// Ref: alter_table_stmt's EXTERNAL alternatives and
// add_external_table_partition_action in sql_parser_mysql_mode.y (oceanbase
// v4.3.5_CE_BP6, lines 19244-19313).
func (p *Parser) parseAlterExternalPartition() (*nodes.AlterTableCmd, bool, error) {
	if p.cur.Type == kwDROP {
		if p.peekNext().Type != kwPARTITION {
			return nil, false, nil
		}
		s := p.saveScan()
		p.advance() // DROP
		p.advance() // PARTITION
		literal := p.isBareWord("LOCATION") && p.peekNext().Type == tokSCONST
		p.restoreScan(s)
		if !literal {
			return nil, false, nil
		}
		cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: p.pos()}, Type: nodes.ATDropExternalPartition}
		p.advance() // DROP
		p.advance() // PARTITION
		if !p.isBareWord("LOCATION") {
			return nil, false, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, false, p.syntaxErrorAtCur()
		}
		cmd.Location = p.advance().Str
		cmd.Loc.End = p.prev.End
		return cmd, true, nil
	}

	if !p.atExternalAddPartition() {
		return nil, false, nil
	}
	cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: p.pos()}, Type: nodes.ATAddExternalPartition}
	p.advance() // ADD
	p.advance() // PARTITION
	if _, err := p.expect('('); err != nil {
		return nil, false, err
	}
	// An empty first entry: `ADD PARTITION(, ds = 'x')`.
	if p.cur.Type == ',' {
		p.advance()
		if p.cur.Type == ')' {
			return nil, false, p.syntaxErrorAtCur()
		}
	}
	for p.cur.Type != ')' {
		value, err := p.parseExternalPartitionAction()
		if err != nil {
			return nil, false, err
		}
		cmd.PartitionValues = append(cmd.PartitionValues, value)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == ')' {
			return nil, false, p.syntaxErrorAtCur()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, false, err
	}
	if !p.isBareWord("LOCATION") {
		return nil, false, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != tokSCONST {
		return nil, false, p.syntaxErrorAtCur()
	}
	cmd.Location = p.advance().Str
	cmd.Loc.End = p.prev.End
	return cmd, true, nil
}

// atExternalAddPartition reports whether ADD PARTITION here opens the external
// action, whose partition list LOCATION 's' follows, rather than the generic
// ADD PARTITION (partition definitions).
func (p *Parser) atExternalAddPartition() bool {
	if p.cur.Type != kwADD || p.peekNext().Type != kwPARTITION {
		return false
	}
	s := p.saveScan()
	p.advance() // ADD
	p.advance() // PARTITION
	if p.cur.Type != '(' {
		p.restoreScan(s)
		return false
	}
	p.scanParenBlock()
	external := p.isBareWord("LOCATION") && p.peekNext().Type == tokSCONST
	p.restoreScan(s)
	return external
}

// parseExternalPartitionAction reads add_external_table_partition_action:
// column_name opt_equal_mark expr_const.
func (p *Parser) parseExternalPartitionAction() (*nodes.Assignment, error) {
	start := p.pos()
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	nameEnd := p.prev.End
	p.match('=')
	value, err := p.parseExprConst()
	if err != nil {
		return nil, err
	}
	return &nodes.Assignment{
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
		Column: &nodes.ColumnRef{Loc: nodes.Loc{Start: start, End: nameEnd}, Column: name},
		Value:  value,
	}, nil
}
