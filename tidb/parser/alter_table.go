package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
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
		for _, tok := range []int{
			kwADD, kwDROP, kwMODIFY, kwCHANGE, kwRENAME, kwALTER, kwCONVERT,
			kwENGINE, kwDEFAULT, kwORDER, kwALGORITHM, kwLOCK, kwFORCE,
			// TiDB-specific: SET TIFLASH REPLICA n, REMOVE TTL.
			kwSET, kwREMOVE,
		} {
			p.addTokenCandidate(tok)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// The statement-level actions, which are whole arms of ALTER TABLE rather
	// than entries of the comma-separated list, so nothing may precede or
	// follow them.
	switch {
	case p.cur.Type == kwCOMPACT:
		if err := p.parseAlterCompact(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	case p.cur.Type == kwANALYZE && p.peekNext().Type == kwPARTITION:
		if err := p.parseAlterAnalyzePartition(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
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

	// A comma the list never got an action for ends the statement with nothing
	// after it. MySQL 8.0.46, MariaDB 11.8.8 and TiDB v8.5.5 all answer 1064
	// to `ALTER TABLE t ENGINE = InnoDB,` and to `ALTER TABLE t ADD COLUMN a
	// INT,`.
	if pendingComma {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterCompact parses TiDB's COMPACT, which is a whole arm of ALTER TABLE
// rather than an entry of the action list.
//
//	ALTER TABLE tbl COMPACT [TIFLASH REPLICA]
//	ALTER TABLE tbl COMPACT PARTITION partition_names [TIFLASH REPLICA]
//
// Nothing may precede or follow it, because the grammar builds a
// CompactTableStmt from the whole statement rather than an AlterTableSpec:
// v8.5.5 answers 1064 for `COMPACT, CACHE`, for `CACHE, COMPACT`, for
// `COMPACT PARTITION` with no names, for `COMPACT TIFLASH` without REPLICA and
// for `COMPACT REPLICA`. MySQL 8.0.46 and MariaDB 11.8.8 answer 1064 for the
// whole form.
//
// Ref: the four `"ALTER" IgnoreOptional "TABLE" TableName "COMPACT" ...` arms
// of AlterTableStmt in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseAlterCompact(stmt *nodes.AlterTableStmt) error {
	cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: p.pos()}, Type: nodes.ATCompact}
	p.advance() // consume COMPACT
	if _, ok := p.match(kwPARTITION); ok {
		names, err := p.parseIdentList()
		if err != nil {
			return err
		}
		cmd.PartitionNames = names
	}
	if p.cur.Type == kwTIFLASH {
		p.advance()
		if _, err := p.expect(kwREPLICA); err != nil {
			return err
		}
		cmd.CompactTiFlashReplica = true
	}
	cmd.Loc.End = p.prev.End
	stmt.Commands = append(stmt.Commands, cmd)
	return nil
}

// parseAlterAnalyzePartition parses TiDB's ANALYZE PARTITION, which like
// COMPACT is a whole arm of ALTER TABLE rather than an entry of the action
// list.
//
//	ALTER TABLE tbl ANALYZE PARTITION partition_names
//	    [INDEX [index_names]] [WITH analyze_option, ...]
//
// The partition names are a plain list with no ALL, nothing may precede or
// follow the action, and INDEX takes an empty name list: v8.5.5 runs
// `ANALYZE PARTITION p`, `ANALYZE PARTITION p, q`, `ANALYZE PARTITION p INDEX`,
// `ANALYZE PARTITION p INDEX i, j WITH 4 BUCKETS, 3 TOPN` and
// `ANALYZE PARTITION p WITH 0.5 SAMPLERATE`, and answers 1064 for
// `ANALYZE PARTITION ALL`, for `ANALYZE PARTITION p, ADD COLUMN a INT` and for
// `ADD COLUMN a INT, ANALYZE PARTITION p`. MySQL 8.0.46 and MariaDB 11.8.8
// take `ANALYZE PARTITION {names | ALL}` as an ordinary action and answer 1064
// for the INDEX and WITH tails.
//
// Ref: the two `"ALTER" IgnoreOptional "TABLE" TableName "ANALYZE" "PARTITION"
// ...` arms of AlterTableStmt in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseAlterAnalyzePartition(stmt *nodes.AlterTableStmt) error {
	cmd := &nodes.AlterTableCmd{Loc: nodes.Loc{Start: p.pos()}, Type: nodes.ATAnalyzePartition}
	p.advance() // consume ANALYZE
	p.advance() // consume PARTITION
	names, err := p.parseIdentList()
	if err != nil {
		return err
	}
	cmd.PartitionNames = names

	if _, ok := p.match(kwINDEX); ok {
		cmd.AnalyzeIndexFlag = true
		// IndexNameList has an empty arm, so INDEX with no name is the whole
		// clause; PRIMARY is a name here as well as an identifier.
		if p.isIdentToken() || p.cur.Type == kwPRIMARY {
			for {
				var name string
				if p.cur.Type == kwPRIMARY {
					name = p.advance().Str
				} else if name, _, err = p.parseIdent(); err != nil {
					return err
				}
				cmd.AnalyzeIndexes = append(cmd.AnalyzeIndexes, name)
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
	}

	// The WITH list is the one ANALYZE TABLE takes, read by the same
	// function: same units, same integer-or-number rule per unit.
	if p.cur.Type == kwWITH {
		opts, err := p.parseAnalyzeOptions()
		if err != nil {
			return err
		}
		cmd.AnalyzeOptions = opts
	}
	cmd.Loc.End = p.prev.End
	stmt.Commands = append(stmt.Commands, cmd)
	return nil
}

// startsAlterPartitioningTail reports whether the parser sits on one of the
// clauses that close an ALTER TABLE: PARTITION BY ..., REMOVE PARTITIONING,
// `PARTITION name ...`, SPLIT MAXVALUE PARTITION LESS THAN (...) or MERGE
// FIRST PARTITION LESS THAN (...). REMOVE not followed by PARTITIONING, SPLIT
// not followed by MAXVALUE and MERGE not followed by FIRST belong to the
// action list instead.
//
// Ref: AlterTableSpecSingleOpt in TiDB v8.5.5 pkg/parser/parser.y, the one
// optional clause `AlterTableSpecListOpt AlterTableSpecSingleOpt` ends with.
func (p *Parser) startsAlterPartitioningTail() bool {
	switch p.cur.Type {
	case kwPARTITION:
		return true
	case kwREMOVE:
		return p.peekNext().Type == kwPARTITIONING
	case kwSPLIT:
		return p.peekNext().Type == kwMAXVALUE
	case kwMERGE:
		return p.peekNext().Type == kwFIRST
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
	switch {
	case p.cur.Type == kwREMOVE:
		p.advance() // consume REMOVE
		if _, err := p.expect(kwPARTITIONING); err != nil {
			return err
		}
		cmd.Type = nodes.ATRemovePartitioning

	case p.cur.Type == kwSPLIT:
		// SPLIT MAXVALUE PARTITION LESS THAN ( expr )
		p.advance() // consume SPLIT
		p.advance() // consume MAXVALUE
		if err := p.parseIntervalPartitionBound(cmd, nodes.ATReorganizeLastPartition); err != nil {
			return err
		}

	case p.cur.Type == kwMERGE:
		// MERGE FIRST PARTITION LESS THAN ( expr )
		p.advance() // consume MERGE
		p.advance() // consume FIRST
		if err := p.parseIntervalPartitionBound(cmd, nodes.ATReorganizeFirstPartition); err != nil {
			return err
		}

	case p.cur.Type == kwPARTITION && p.peekNext().Type != kwBY:
		if err := p.parseAlterPartitionOptions(cmd); err != nil {
			return err
		}

	default:
		part, err := p.parsePartitionClause()
		if err != nil {
			return err
		}
		cmd.Type = nodes.ATPartitionBy
		cmd.PartitionBy = part
	}
	cmd.Loc.End = p.prev.End
	stmt.Commands = append(stmt.Commands, cmd)
	return nil
}

// parseIntervalPartitionBound parses the `PARTITION LESS THAN ( expr )` that
// the four interval-partition actions name a partition by, from the PARTITION
// keyword on.
//
// The parentheses are required and MAXVALUE is not an expression here: TiDB
// v8.5.5 answers 1064 for `LAST PARTITION LESS THAN 1000` and for
// `LAST PARTITION LESS THAN (MAXVALUE)`.
//
// Ref: the `"LAST" | "FIRST" | "SPLIT" "MAXVALUE" | "MERGE" "FIRST"` arms,
// each spelled `"PARTITION" "LESS" "THAN" '(' BitExpr ')'`, in TiDB v8.5.5
// pkg/parser/parser.y.
func (p *Parser) parseIntervalPartitionBound(cmd *nodes.AlterTableCmd, cmdType nodes.AlterTableCmdType) error {
	cmd.Type = cmdType
	if _, err := p.expect(kwPARTITION); err != nil {
		return err
	}
	if _, err := p.expect(kwLESS); err != nil {
		return err
	}
	if _, err := p.expect(kwTHAN); err != nil {
		return err
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil {
		return p.syntaxErrorAtCur()
	}
	cmd.PartitionBound = expr
	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
}

// parseAlterPartitionOptions parses the `PARTITION name ...` clause that may
// close an ALTER TABLE, from the PARTITION keyword on.
//
//	PARTITION name ATTRIBUTES [=] {'string' | DEFAULT}
//	PARTITION name [partition_definition_option...]
//
// It is a closing clause rather than an entry of the action list, so it takes
// no comma before it and nothing may follow it: TiDB v8.5.5 runs
// `CACHE PARTITION p ATTRIBUTES = 'a'`, `ADD COLUMN a INT PARTITION p ENGINE =
// x` and a bare `PARTITION p` with no options at all, and answers 1064 for
// `CACHE, PARTITION p ATTRIBUTES = 'a'`, for
// `PARTITION p ATTRIBUTES = 'a' , CACHE` and for `PARTITION p ENGINE = x,
// MAX_ROWS 5`. MySQL 8.0.46 and MariaDB 11.8.8 answer 1064 for the whole form.
//
// Ref: the `"PARTITION" Identifier AttributesOpt` and
// `"PARTITION" Identifier PartDefOptionList` arms of AlterTableSpecSingleOpt
// in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseAlterPartitionOptions(cmd *nodes.AlterTableCmd) error {
	p.advance() // consume PARTITION
	name, _, err := p.parseIdent()
	if err != nil {
		return err
	}
	cmd.PartitionNames = []string{name}

	if p.isBareWord("ATTRIBUTES") {
		spec, err := p.parseAttributesSpec()
		if err != nil {
			return err
		}
		cmd.Type = nodes.ATPartitionAttributes
		cmd.Attributes = spec
		return nil
	}

	cmd.Type = nodes.ATPartitionOptions
	for {
		opt, ok, err := p.parsePartDefOption()
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		cmd.PartitionOptions = append(cmd.PartitionOptions, opt)
	}
}

// parseAttributesSpec parses an ATTRIBUTES or STATS_OPTIONS clause from its
// own word on. Both take a string literal or DEFAULT, with an optional = but
// no other value: TiDB v8.5.5 answers 1064 for a bare `ATTRIBUTES`, for
// `ATTRIBUTES = 1` and for `ATTRIBUTES 1`, and runs `ATTRIBUTES 'str'`,
// `ATTRIBUTES = 'str'` and `ATTRIBUTES = DeFaUlT`.
//
// Neither word is in this parser's keyword table, so each is read as a bare
// word; MySQL 8.0.46 has neither option, and MariaDB 11.8.8 takes
// `ATTRIBUTES = 'str'` only as one of the engine-defined attributes its own
// grammar lets any `ident = value` be, which is a different clause.
//
// Ref: AttributesOpt and StatsOptionsOpt in TiDB v8.5.5 pkg/parser/parser.y.
func (p *Parser) parseAttributesSpec() (*nodes.AttributesSpec, error) {
	start := p.pos()
	p.advance() // consume ATTRIBUTES or STATS_OPTIONS
	p.match('=')
	spec := &nodes.AttributesSpec{Loc: nodes.Loc{Start: start}}
	if _, ok := p.match(kwDEFAULT); ok {
		spec.Default = true
		spec.Loc.End = p.prev.End
		return spec, nil
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	spec.Value = p.cur.Str
	p.advance()
	spec.Loc.End = p.prev.End
	return spec, nil
}

// parseNoWriteToBinlog consumes the modifier that follows the PARTITION keyword
// on a partition-maintenance action and reports whether one was there.
//
// Ref: opt_no_write_to_binlog in MySQL 8.0.46 sql/sql_yacc.yy, which reads
// `/* empty */ | NO_WRITE_TO_BINLOG | LOCAL_SYM` and folds both spellings to
// the same value, so one bool is the whole clause. TiDB v8.5.5 takes it on
// ADD, REBUILD, OPTIMIZE, REPAIR, COALESCE and REORGANIZE PARTITION, and
// answers 1064 for it on CHECK, TRUNCATE, DROP and ANALYZE PARTITION -- MySQL
// 8.0.46 and MariaDB 11.8.8 take it on ANALYZE, TiDB does not -- for a second
// modifier after the first, and for one after the partition names.
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
		for _, t := range []int{
			kwADD, kwDROP, kwMODIFY, kwCHANGE, kwRENAME, kwALTER, kwCONVERT,
			kwENGINE, kwDEFAULT, kwORDER, kwALGORITHM, kwLOCK, kwFORCE, kwPARTITION,
			// TiDB-specific: SET TIFLASH REPLICA n, REMOVE TTL.
			kwSET, kwREMOVE,
		} {
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
		val, err := p.parseAlgorithmValue()
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
		val, err := p.parseLockValue()
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
		// FORCE on its own is the standalone action; FORCE AUTO_INCREMENT and
		// FORCE AUTO_RANDOM_BASE are table options, which the default arm
		// below reads.
		if next := p.peekNext().Type; next == kwAUTO_INCREMENT || next == kwAUTO_RANDOM_BASE {
			break
		}
		p.advance()
		cmd.Type = nodes.ATForce
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwCACHE, kwNOCACHE:
		// TiDB's CACHE and NOCACHE, which take no value and chain only with a
		// comma: v8.5.5 runs `CACHE, NOCACHE` and answers 1064 for
		// `CACHE NOCACHE` and for `CACHE = 1`. MySQL 8.0.46 and MariaDB
		// 11.8.8 answer 1064 for both words.
		//
		// Ref: the "CACHE" and "NOCACHE" arms of AlterTableSpec in TiDB
		// v8.5.5 pkg/parser/parser.y.
		if p.cur.Type == kwCACHE {
			cmd.Type = nodes.ATCache
		} else {
			cmd.Type = nodes.ATNoCache
		}
		p.advance()
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwREAD:
		// TiDB's READ ONLY and READ WRITE, which like CACHE chain only with a
		// comma: v8.5.5 runs `READ ONLY, READ WRITE` and answers 1064 for a
		// bare `READ`, for `READ ONLY READ WRITE`, for
		// `ENGINE = InnoDB READ ONLY` and for `READ ONLY ENGINE = InnoDB`.
		// MySQL 8.0.46 and MariaDB 11.8.8 answer 1064 for both forms.
		//
		// Ref: the Writeable arm of AlterTableSpec in TiDB v8.5.5
		// pkg/parser/parser.y.
		p.advance() // consume READ
		cmd.Type = nodes.ATWriteable
		switch {
		case p.cur.Type == kwWRITE:
			cmd.Writeable = true
			p.advance()
		case p.cur.Type == kwONLY:
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwLAST:
		// LAST PARTITION LESS THAN ( expr ) [NO_WRITE_TO_BINLOG | LOCAL],
		// which adds a partition to an interval-partitioned table. It is an
		// entry of the action list, unlike SPLIT MAXVALUE and MERGE FIRST:
		// v8.5.5 runs `ADD COLUMN a INT, LAST PARTITION LESS THAN (1)` and
		// `LAST PARTITION LESS THAN (1), ADD COLUMN a INT`, and answers 1064
		// for the same two statements with SPLIT MAXVALUE PARTITION.
		//
		// Ref: the `"LAST" "PARTITION" "LESS" "THAN" '(' BitExpr ')'
		// NoWriteToBinLogAliasOpt` arm of AlterTableSpec in TiDB v8.5.5
		// pkg/parser/parser.y.
		p.advance() // consume LAST
		if err := p.parseIntervalPartitionBound(cmd, nodes.ATAddLastPartition); err != nil {
			return nil, err
		}
		cmd.NoWriteToBinlog = p.parseNoWriteToBinlog()
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwFIRST:
		// FIRST PARTITION LESS THAN ( expr ) [IF EXISTS], which drops the
		// partitions below a bound.
		//
		// Ref: the `"FIRST" "PARTITION" "LESS" "THAN" '(' BitExpr ')' IfExists`
		// arm of AlterTableSpec in TiDB v8.5.5 pkg/parser/parser.y.
		p.advance() // consume FIRST
		if err := p.parseIntervalPartitionBound(cmd, nodes.ATDropFirstPartition); err != nil {
			return nil, err
		}
		ie, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		cmd.IfExists = ie
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

	case kwREMOVE:
		// TiDB: REMOVE TTL. REMOVE PARTITIONING is the statement's tail, which
		// parseAlterPartitioningTail takes before this dispatch sees it.
		p.advance()
		if _, err := p.expect(kwTTL); err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATRemoveTTL
		cmd.Loc.End = p.prev.End
		return cmd, nil

	case kwSET:
		// TiDB: SET TIFLASH REPLICA n [LOCATION LABELS 'label1', 'label2', ...]
		// Note: no MySQL ALTER TABLE form starts with bare SET at this dispatch level.
		// ALTER COLUMN col SET DEFAULT routes through the kwALTER branch, not here.
		// Ref: parser.y:2193 (AlterTableSpec) + parser.y:2176-2183 (LocationLabelList).
		p.advance()
		if p.cur.Type == kwTIFLASH {
			p.advance()
			if _, ok := p.match(kwREPLICA); !ok {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			count := int(p.cur.Ival)
			p.advance()
			labels, err := p.parseLocationLabelList()
			if err != nil {
				return nil, err
			}
			cmd.Type = nodes.ATSetTiFlashReplica
			cmd.TiFlashReplica = count
			cmd.TiFlashLocationLabels = labels
			cmd.Loc.End = p.prev.End
			return cmd, nil
		}
		return nil, p.syntaxErrorAtCur()

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

	}

	// ATTRIBUTES and STATS_OPTIONS are actions of their own rather than table
	// options, so neither chains with a space after another option: v8.5.5
	// runs `STATS_OPTIONS = 'a', ATTRIBUTES = 'b'` and answers 1064 for
	// `ATTRIBUTES = 'a' ATTRIBUTES = 'b'`, `ATTRIBUTES = 'a' ENGINE = x` and
	// `ENGINE = x ATTRIBUTES = 'a'`. STATS_OPTIONS is an ALTER TABLE action
	// only: v8.5.5 answers 1064 for it on CREATE TABLE.
	switch {
	case p.isBareWord("ATTRIBUTES"):
		spec, err := p.parseAttributesSpec()
		if err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATAttributes
		cmd.Attributes = spec
		cmd.Loc.End = p.prev.End
		return cmd, nil
	case p.isBareWord("STATS_OPTIONS"):
		spec, err := p.parseAttributesSpec()
		if err != nil {
			return nil, err
		}
		cmd.Type = nodes.ATStatsOptions
		cmd.Attributes = spec
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

// parseAlterAdd parses ADD operations.
func (p *Parser) parseAlterAdd(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume ADD

	// Completion: after ADD, offer what can be added.
	p.checkCursor()
	if p.collectMode() {
		for _, t := range []int{kwCOLUMN, kwINDEX, kwKEY, kwUNIQUE, kwPRIMARY, kwFOREIGN, kwCONSTRAINT, kwCHECK, kwPARTITION, kwSPATIAL, kwFULLTEXT, kwVECTOR} {
			p.addTokenCandidate(t)
		}
		return nil, &ParseError{Message: "collecting"}
	}

	// ADD PARTITION [IF NOT EXISTS] [NO_WRITE_TO_BINLOG | LOCAL]
	//     [PARTITIONS n | (partition_definition, ...)]
	if p.cur.Type == kwPARTITION {
		p.advance()
		cmd.Type = nodes.ATAddPartition
		ine, err := p.parseIfNotExists()
		if err != nil {
			return nil, err
		}
		cmd.IfNotExists = ine
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

	// VECTOR after a bare ADD can only open VECTOR INDEX, which was read above. Its
	// precedence over the column reading is what lets the grammar take VECTOR as a
	// type, so TiDB v8.5.5 answers 1064 for `ADD vector INT` and runs `ADD COLUMN
	// vector INT`, ``ADD `vector` INT`` and `ADD comment INT`.
	//
	// Ref: `%precedence vectorType` in TiDB v8.5.5 pkg/parser/parser.y.
	if p.cur.Type == kwVECTOR {
		return nil, p.syntaxErrorAtCur()
	}

	// ADD [COLUMN] [IF NOT EXISTS]
	//
	// IF NOT EXISTS belongs to ADD [COLUMN], ADD PARTITION, ADD {INDEX | KEY}
	// and ADD FOREIGN KEY, and to nothing else: TiDB v8.5.5 answers 1064 for
	// it after UNIQUE, PRIMARY KEY, FULLTEXT and CHECK, and for
	// `ADD COLUMN IF EXISTS`. MySQL 8.0.46 has none of these forms; MariaDB
	// 11.8.8 takes them and more.
	cmd.Type = nodes.ATAddColumn
	p.match(kwCOLUMN)
	ine, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	cmd.IfNotExists = ine

	// Completion: after ADD COLUMN, identifier context (no specific candidates).
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new column name.
		return nil, &ParseError{Message: "collecting"}
	}

	// ADD [COLUMN] (table_element, ...) — the paren form, which reads a table
	// element list and so takes constraints beside columns:
	// `ADD COLUMN (a1 int, CONSTRAINT b1 CHECK (a1>0))` and
	// `ADD (PRIMARY KEY (a))` are one action each. TiDB v8.5.5, MySQL 8.0.46
	// and MariaDB 11.8.8 all run them, and all three answer 1064 for an empty
	// list and for a trailing comma.
	//
	// Ref: `"ADD" ColumnKeywordOpt IfNotExists '(' TableElementList ')'` in
	// TiDB v8.5.5 pkg/parser/parser.y.
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
// Ref: the DROP arms of AlterTableSpec in TiDB v8.5.5 pkg/parser/parser.y:
// "DROP" ColumnKeywordOpt IfExists ColumnName RestrictOrCascadeOpt, "DROP"
// "PRIMARY" "KEY", "DROP" "FOREIGN" "KEY" Symbol. KEY is required after
// PRIMARY and FOREIGN, and only a column takes RESTRICT or CASCADE.
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
		// DROP PARTITION [IF EXISTS] partition_names
		p.advance()
		cmd.Type = nodes.ATDropPartition
		ie, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		cmd.IfExists = ie
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
		// DROP {INDEX | KEY} [IF EXISTS] index_name. IF EXISTS belongs here,
		// to DROP [COLUMN] and to DROP PARTITION, and to nothing else: TiDB
		// v8.5.5 answers 1064 for `DROP CONSTRAINT IF EXISTS c`,
		// `DROP CHECK IF EXISTS c`, `DROP FOREIGN KEY IF EXISTS c` and
		// `DROP PRIMARY KEY IF EXISTS`.
		p.advance()
		cmd.Type = nodes.ATDropIndex
		// Completion: after DROP INDEX, offer index_ref.
		p.checkCursor()
		if p.collectMode() {
			p.addRuleCandidate("index_ref")
			return nil, &ParseError{Message: "collecting"}
		}
		ie, err := p.parseIfExists()
		if err != nil {
			return nil, err
		}
		cmd.IfExists = ie
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
		// DROP [COLUMN] [IF EXISTS] col_name [RESTRICT | CASCADE]
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
		_, _, name, err := p.parseColumnName()
		if err != nil {
			return nil, err
		}
		cmd.Name = name
		// RestrictOrCascadeOpt: the grammar action ignores it.
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
	// MODIFY [COLUMN] [IF EXISTS]: TiDB v8.5.5 and MariaDB 11.8.8 run
	// `MODIFY COLUMN IF EXISTS a varchar(255)` and both answer 1064 for
	// `MODIFY COLUMN IF NOT EXISTS`; MySQL 8.0.46 has neither.
	ie, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	cmd.IfExists = ie

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
	// CHANGE [COLUMN] [IF EXISTS], as MODIFY takes it.
	ie, err := p.parseIfExists()
	if err != nil {
		return nil, err
	}
	cmd.IfExists = ie

	// Completion: after CHANGE [COLUMN], offer columnref (old name).
	p.checkCursor()
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, &ParseError{Message: "collecting"}
	}

	// Old column name
	oldSchema, oldTable, oldName, err := p.parseColumnName()
	if err != nil {
		return nil, err
	}
	cmd.Name = oldName
	cmd.OldColumnSchema = oldSchema
	cmd.OldColumnTable = oldTable

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
	// The enforcement is required, not optional: MySQL 8.0.46, MariaDB 11.8.8
	// and TiDB v8.5.5 all answer 1064 for `ALTER CONSTRAINT ident` and for
	// `ALTER CHECK ident` with nothing after the name. CHECK and CONSTRAINT
	// are the same clause spelled two ways.
	//
	// Ref: `"ALTER" CheckConstraintKeyword Identifier EnforcedOrNot` in TiDB
	// v8.5.5 pkg/parser/parser.y, where CheckConstraintKeyword reads
	// `"CHECK" | "CONSTRAINT"` and EnforcedOrNot reads
	// `"ENFORCED" | NotSym "ENFORCED"`, neither with an empty arm.
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

	schema, table, name, err := p.parseColumnName()
	if err != nil {
		return nil, err
	}
	cmd.Name = name
	qualified := schema != "" || table != ""

	if _, ok := p.match(kwSET); ok {
		if _, ok := p.match(kwDEFAULT); ok {
			cmd.Type = nodes.ATAlterColumnDefault
			// SET DEFAULT takes a signed literal, or any expression inside
			// parentheses, and nothing in between: MySQL 8.0.46 and TiDB
			// v8.5.5 both answer 1064 for `SET DEFAULT CURRENT_TIMESTAMP`,
			// `SET DEFAULT NOW()` and `SET DEFAULT 1+1`, and both run
			// `SET DEFAULT (CURRENT_TIMESTAMP)` and `SET DEFAULT (1+1)`.
			// MariaDB 11.8.8 takes a bare expression here, so this is not
			// ported to that fork.
			//
			// Ref: the `"ALTER" ColumnKeywordOpt ColumnName "SET" "DEFAULT"
			// SignedLiteral` and `... "SET" "DEFAULT" '(' Expression ')'` arms
			// of AlterTableSpec in TiDB v8.5.5 pkg/parser/parser.y;
			// SignedLiteral reads `Literal | '+' NumLiteral | '-' NumLiteral`.
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
		} else if qualified && (p.cur.Type == kwVISIBLE || p.cur.Type == kwINVISIBLE) {
			// TiDB refuses column visibility outright, and omni still reads MySQL's form for
			// a bare name. It takes no qualifier here, so this change adds no false accept.
			return nil, p.syntaxErrorAtCur()
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
		// TO is required: MySQL 8.0.46, MariaDB 11.8.8 and TiDB v8.5.5 all
		// answer 1064 for `RENAME COLUMN a b`.
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
		// TO is required here too: all three answer 1064 for
		// `RENAME INDEX a b`.
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
		// The new name is a table name, so it may carry a database
		// qualifier, and the separator may also be an equals sign: MySQL
		// 8.0.46, MariaDB 11.8.8 and TiDB v8.5.5 all run
		// `RENAME TO d.t1`, `RENAME AS d.t1`, `RENAME = d.t1` and a bare
		// `RENAME d.t1`.
		//
		// Ref: the `"RENAME" "TO" TableName`, `"RENAME" EqOpt TableName` and
		// `"RENAME" "AS" TableName` arms of AlterTableSpec in TiDB v8.5.5
		// pkg/parser/parser.y; MySQL 8.0.46 sql/sql_yacc.yy spells the same
		// three as `RENAME opt_to table_ident`.
		cmd.Type = nodes.ATRenameTable
		if _, ok := p.match(kwTO, kwAS); !ok {
			p.match('=')
		}
		ref, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		cmd.NewTable = ref
		if ref.Schema != "" {
			cmd.NewName = ref.Schema + "." + ref.Name
		} else {
			cmd.NewName = ref.Name
		}
	}

	cmd.Loc.End = p.prev.End
	return cmd, nil
}

// parseAlterConvert parses CONVERT TO CHARACTER SET charset_name [COLLATE collation_name].
func (p *Parser) parseAlterConvert(cmd *nodes.AlterTableCmd) (*nodes.AlterTableCmd, error) {
	p.advance() // consume CONVERT
	cmd.Type = nodes.ATConvertCharset

	// CONVERT TO {CHARACTER SET | CHARSET | CHAR SET} {charset | DEFAULT}
	//     [COLLATE collation]
	//
	// Every word is required and each has a value: MySQL 8.0.46, MariaDB
	// 11.8.8 and TiDB v8.5.5 all run the three spellings of the keyword and
	// `CONVERT TO CHARSET DEFAULT`, and all three answer 1064 for
	// `CONVERT TO utf8`, for `CONVERT CHARSET utf8`, for a `CONVERT TO
	// CHARSET` with no charset after it and for `COLLATE = utf8_bin`, whose
	// equals sign the grammar does not allow.
	//
	// Ref: the `"CONVERT" "TO" CharsetKw CharsetName OptCollate` and
	// `"CONVERT" "TO" CharsetKw "DEFAULT" OptCollate` arms of AlterTableSpec
	// in TiDB v8.5.5 pkg/parser/parser.y; CharsetKw reads
	// `"CHARACTER" "SET" | "CHARSET" | "CHAR" "SET"`.
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
		charset, err := p.consumeOptionValue()
		if err != nil {
			return nil, err
		}
		cmd.Name = charset
	}

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
		_, _, name, err := p.parseColumnName()
		if err != nil {
			return err
		}
		cmd.After = name
	}
	return nil
}

// expectCharsetKw consumes the CharsetKw that introduces a charset name:
// CHARACTER SET, CHARSET, or CHAR SET. The two-word spellings need both words.
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

// parseSignedLiteral parses the SignedLiteral a SET DEFAULT takes outside
// parentheses: one literal, optionally signed. A function call, a keyword such
// as CURRENT_TIMESTAMP and an arithmetic expression are all syntax errors
// here, which is why this is not parseExpr.
//
// Ref: SignedLiteral in TiDB v8.5.5 pkg/parser/parser.y, which reads
// `Literal | '+' NumLiteral | '-' NumLiteral`.
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
// charset introducer, a hex or bit literal, NULL, TRUE or FALSE.
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
		// A temporal literal is not one of these: TiDB v8.5.5 answers 1064 for
		// `SET DEFAULT DATE '2020-01-01'`, which MySQL 8.0.46 and OceanBase CE
		// 4.3.5.6 both run, because its own Literal rule has no temporal arm.
	}
	// A charset introducer followed by a string, as in _utf8'x', is a literal
	// too, and the introducer reaches the parser as an identifier.
	if p.cur.Type == tokIDENT && strings.HasPrefix(p.cur.Str, "_") && p.peekNext().Type == tokSCONST {
		return p.parsePrimaryExpr()
	}
	return nil, p.syntaxErrorAtCur()
}
