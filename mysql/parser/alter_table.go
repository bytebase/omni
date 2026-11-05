package parser

import (
	"strings"

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
func (p *Parser) parseAlterTableStmt(start int) (*nodes.AlterTableStmt, error) {
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

	// Parse comma-separated alter operations, which one partitioning clause
	// may close.
	pendingComma := false
	for {
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			break
		}
		// A partitioning clause belongs to the tail below, not to this list.
		if p.startsAlterPartitioningTail() {
			break
		}

		cmd, err := p.parseAlterTableCmd()
		if err != nil {
			return nil, err
		}
		stmt.Commands = append(stmt.Commands, cmd)
		if cmd.Type == nodes.ATTableOption {
			more, err := p.parseMoreAlterTableOptions()
			if err != nil {
				return nil, err
			}
			stmt.Commands = append(stmt.Commands, more...)
		}
		pendingComma = false

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		pendingComma = true
	}

	if err := p.parseAlterPartitioningTail(stmt, pendingComma); err != nil {
		return nil, err
	}

	// A comma the list never got an action for ends the statement with
	// nothing after it: MySQL 8.0.46 answers 1064 for
	// `ALTER TABLE t ENGINE = InnoDB,` and for
	// `ALTER TABLE t ADD COLUMN a INT,`.
	if pendingComma {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// startsAlterPartitioningTail reports whether the parser sits on the
// partitioning clause that closes an ALTER TABLE: PARTITION BY ... or REMOVE
// PARTITIONING. PARTITION not followed by BY, and REMOVE not followed by
// PARTITIONING, belong to the action list instead.
func (p *Parser) startsAlterPartitioningTail() bool {
	switch p.cur.Type {
	case kwPARTITION:
		return p.peekNext().Type == kwBY
	case kwREMOVE:
		return p.peekNext().Type == kwPARTITIONING
	default:
		return false
	}
}

// parseAlterPartitioningTail parses the one partitioning clause an ALTER TABLE
// may end with, and is a no-op when there is none.
//
// Ref: alter_commands in MySQL 8.0.46 sql/sql_yacc.yy, which spells the
// statement's tail as `alter_list opt_partitioning`, `alter_list
// remove_partitioning`, `partitioning` or `remove_partitioning`: PARTITION BY
// ... or REMOVE PARTITIONING follows the comma-separated action list without a
// comma of its own, at most once, with nothing after it. MySQL 8.0.46, MariaDB
// 11.8.8 and TiDB v8.5.5 all run
// `ALTER TABLE t ADD COLUMN a INT REMOVE PARTITIONING` and
// `ALTER TABLE t COMMENT 'c' PARTITION BY HASH(a)`, and all three answer 1064
// to the same statements with a comma before the clause, to two clauses in a
// row, and to an action after one.
func (p *Parser) parseAlterPartitioningTail(stmt *nodes.AlterTableStmt, pendingComma bool) error {
	if !p.startsAlterPartitioningTail() {
		return nil
	}
	if pendingComma {
		return &ParseError{
			Message:  "a partitioning clause takes no comma before it",
			Position: p.cur.Loc,
		}
	}
	cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == kwREMOVE {
		p.advance() // consume REMOVE
		if _, err := p.expect(kwPARTITIONING); err != nil {
			return err
		}
		cmd.Type = nodes.ATRemovePartitioning
		cmd.Loc.End = p.prev.End
		stmt.Commands = append(stmt.Commands, cmd)
		return nil
	}
	part, err := p.parsePartitionClause()
	if err != nil {
		return err
	}
	cmd.Type = nodes.ATPartitionBy
	cmd.PartitionBy = part
	cmd.Loc.End = p.prev.End
	stmt.Commands = append(stmt.Commands, cmd)
	return nil
}

// parseNoWriteToBinlog consumes the modifier that follows the PARTITION keyword
// on a partition-maintenance action and reports whether one was there.
//
// Ref: opt_no_write_to_binlog in MySQL 8.0.46 sql/sql_yacc.yy, which reads
// `/* empty */ | NO_WRITE_TO_BINLOG | LOCAL_SYM` and folds both spellings to
// the same value, so one bool is the whole clause. MySQL 8.0.46 and MariaDB
// 11.8.8 take it on ADD, REBUILD, OPTIMIZE, REPAIR, ANALYZE, COALESCE and
// REORGANIZE PARTITION, and answer 1064 for it on CHECK, TRUNCATE and DROP
// PARTITION, for a second modifier after the first, and for one after the
// partition names -- the last of which MariaDB allows on OPTIMIZE alone.
func (p *Parser) parseNoWriteToBinlog() bool {
	if _, ok := p.match(kwNO_WRITE_TO_BINLOG); ok {
		return true
	}
	if _, ok := p.match(kwLOCAL); ok {
		return true
	}
	return false
}

// parseMoreAlterTableOptions parses the table options that follow one in the
// same ALTER TABLE action, without commas between them, each its own
// ATTableOption command: `ALTER TABLE t ENGINE = InnoDB ROW_FORMAT = DYNAMIC`
// and `ALTER TABLE t DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin`, which MySQL
// 8.0, MariaDB 11.8.8 and TiDB v8.5.5 all run, as CREATE TABLE takes them. An
// action other than a table option still needs a comma before it:
// `ALTER TABLE t COMMENT 'x' DISABLE KEYS` is a 1064 on all three.
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
		p.matchEqual()
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwLOCK:
		p.advance()
		p.matchEqual()
		cmd.Type = nodes.ATLock
		val, err := p.consumeOptionValue()
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
		return nil, &ParseError{
			Message:  "expected PARTITION after TRUNCATE",
			Position: p.cur.Loc,
		}

	case kwANALYZE:
		// ANALYZE PARTITION [NO_WRITE_TO_BINLOG | LOCAL] ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
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
		// OPTIMIZE PARTITION [NO_WRITE_TO_BINLOG | LOCAL] ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATOptimizePartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after OPTIMIZE",
			Position: p.cur.Loc,
		}

	case kwREBUILD:
		// REBUILD PARTITION [NO_WRITE_TO_BINLOG | LOCAL] ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
			return p.parseAlterPartitionNamesOrAll(cmd, nodes.ATRebuildPartition)
		}
		return nil, &ParseError{
			Message:  "expected PARTITION after REBUILD",
			Position: p.cur.Loc,
		}

	case kwREPAIR:
		// REPAIR PARTITION [NO_WRITE_TO_BINLOG | LOCAL] ...
		p.advance()
		if _, ok := p.match(kwPARTITION); ok {
			cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
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
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// DISCARD TABLESPACE
		p.match(kwTABLESPACE)
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
			p.match(kwTABLESPACE)
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// IMPORT TABLESPACE
		p.match(kwTABLESPACE)
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
		p.match(kwBY)
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
		p.match(kwKEYS)
		cmd.Type = nodes.ATEnableKeys
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwDISABLE:
		// DISABLE KEYS
		p.advance()
		p.match(kwKEYS)
		cmd.Type = nodes.ATDisableKeys
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwWITH:
		// WITH VALIDATION
		p.advance()
		p.match(kwVALIDATION)
		cmd.Type = nodes.ATWithValidation
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwWITHOUT:
		// WITHOUT VALIDATION
		p.advance()
		p.match(kwVALIDATION)
		cmd.Type = nodes.ATWithoutValidation
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwPARTITION:
		// PARTITION BY ... is the statement's tail, which
		// parseAlterPartitioningTail takes before this dispatch sees it.
		p.advance() // consume PARTITION
		return nil, &ParseError{
			Message:  "expected BY after PARTITION",
			Position: p.cur.Loc,
		}

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

	// ADD PARTITION [NO_WRITE_TO_BINLOG | LOCAL] [PARTITIONS n | (partition_definition, ...)]
	if p.cur.Type == kwPARTITION {
		p.advance()
		cmd.Type = nodes.ATAddPartition
		cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
		if p.cur.Type == kwPARTITIONS {
			p.advance()
			if p.cur.Type != tokICONST {
				return nil, &ParseError{Message: "expected partition count", Position: p.cur.Loc}
			}
			cmd.Number = int(p.cur.Ival)
			p.advance()
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		// add_part_extra is optional: MySQL 8.0.46, MariaDB 11.8.8 and TiDB
		// v8.5.5 all run `ALTER TABLE t ADD PARTITION` and
		// `ADD PARTITION LOCAL` with nothing after them.
		if p.cur.Type != '(' {
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		p.advance() // consume '('
		for {
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
		cmd.Loc.End = p.prev.End
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

	// ADD [COLUMN] (table_element, ...) — the paren form, which reads a table
	// element list and so takes constraints beside columns:
	// `ADD COLUMN (a1 int, CONSTRAINT b1 CHECK (a1>0))` and
	// `ADD (PRIMARY KEY (a))` are one action each. MySQL 8.0.46 runs them,
	// and answers 1064 for an empty list and for a trailing comma.
	//
	// Ref: `ADD opt_column '(' table_element_list ')'` in MySQL 8.0.46
	// sql/sql_yacc.yy.
	if p.cur.Type == '(' {
		p.advance()
		for {
			if p.isTableConstraintStart() {
				constr, err := p.parseTableConstraint()
				if err != nil {
					return nil, err
				}
				cmd.Constraints = append(cmd.Constraints, constr)
			} else {
				col, err := p.parseColumnDef()
				if err != nil {
					return nil, err
				}
				cmd.Columns = append(cmd.Columns, col)
			}
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

	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	cmd.Column = col

	// FIRST | AFTER col_name
	if err := p.parseColumnPositioning(cmd); err != nil {
		return nil, err
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwPRIMARY:
		// DROP PRIMARY KEY
		p.advance()
		if err := p.expectDropKey(); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATDropConstraint
		cmd.Name = "PRIMARY"
		cmd.Loc.End = p.prev.End
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwFOREIGN:
		// DROP FOREIGN KEY fk_symbol
		p.advance()
		if err := p.expectDropKey(); err != nil {
			return nil, err
		}
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
		cmd.Loc.End = p.prev.End
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
		cmd.Loc.End = p.prev.End
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
		cmd.Loc.End = p.prev.End
		return cmd, nil

	default:
		// DROP [COLUMN] col_name [RESTRICT | CASCADE]
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
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}

	// ALTER {CHECK | CONSTRAINT} symbol {ENFORCED | NOT ENFORCED}
	//
	// The enforcement is required, not optional: MySQL 8.0.46 answers 1064
	// for `ALTER TABLE t ALTER CHECK ident` and for
	// `ALTER TABLE t ALTER CONSTRAINT ident` with nothing after the name.
	// CHECK and CONSTRAINT are the same clause spelled two ways.
	//
	// Ref: `ALTER check_constraint_keyword ident constraint_enforcement` in
	// MySQL 8.0.46 sql/sql_yacc.yy, whose constraint_enforcement has no
	// empty arm.
	if p.cur.Type == kwCHECK || p.cur.Type == kwCONSTRAINT {
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
			// parentheses, and nothing in between: MySQL 8.0.46 answers 1064
			// for `SET DEFAULT CURRENT_TIMESTAMP`, `SET DEFAULT NOW()` and
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
		// TO is required: MySQL 8.0.46 answers 1064 for
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
		// TO is required: MySQL 8.0.46 answers 1064 for
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
		// RENAME [TO | AS | =] new_tbl_name
		//
		// The separator may also be an equals sign: MySQL 8.0.46 runs
		// `ALTER TABLE t RENAME = t1` and `ALTER TABLE db.t RENAME = db1.t1`.
		//
		// Ref: `RENAME opt_to table_ident` in MySQL 8.0.46 sql/sql_yacc.yy,
		// where opt_to reads `TO | = | AS | <empty>`.
		cmd.Type = nodes.ATRenameTable
		if _, ok := p.match(kwTO, kwAS); !ok {
			p.match('=')
		}
		newTable, err := p.parseTableRef()
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

	// CONVERT TO CHARACTER SET, CHARSET or CHAR SET is required, and so is
	// the charset after it: MySQL 8.0.46 answers 1064 for `CONVERT TO utf8`, for
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

	if _, ok := p.match(kwDEFAULT); ok {
		cmd.Name = "DEFAULT"
	} else {
		charset, err := p.parseCharsetName()
		if err != nil {
			return nil, err
		}
		cmd.Name = charset
	}

	if _, ok := p.match(kwCOLLATE); ok {
		collation, err := p.parseCollationName()
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
	cmd.Loc.End = p.prev.End
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
	cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
	// real_ulong_num is not optional: MySQL 8.0.46, MariaDB 11.8.8 and TiDB
	// v8.5.5 all answer 1064 to `ALTER TABLE t COALESCE PARTITION` and to
	// `COALESCE PARTITION LOCAL`.
	if p.cur.Type != tokICONST {
		return nil, p.syntaxErrorAtCur()
	}
	cmd.Number = int(p.cur.Ival)
	p.advance()
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
	cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()

	// REORGANIZE PARTITION with no names is its own arm and takes no INTO:
	// MySQL 8.0.46, MariaDB 11.8.8 and TiDB v8.5.5 all run
	// `ALTER TABLE t REORGANIZE PARTITION LOCAL` and all answer 1064 to
	// `REORGANIZE PARTITION LOCAL INTO (...)`.
	if p.cur.Type == tokEOF || p.cur.Type == ';' || p.cur.Type == ',' {
		cmd.Loc.End = p.prev.End
		return cmd, nil
	}

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

	cmd.Loc.End = p.prev.End
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

// expectCharsetKw consumes the keyword that introduces a charset name. The
// two-word spellings need both words.
func (p *Parser) expectCharsetKw() error {
	switch p.cur.Type {
	case kwCHARSET:
		p.advance()
		return nil
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
		// A temporal literal, which MySQL 8.0.46 takes here:
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
