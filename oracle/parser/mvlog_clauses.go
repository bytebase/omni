package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads CREATE MATERIALIZED VIEW LOG
// (oracle/parser/bnf/CREATE-MATERIALIZED-VIEW-LOG.bnf) and ALTER
// MATERIALIZED VIEW LOG (oracle/parser/bnf/ALTER-MATERIALIZED-VIEW-LOG.bnf).
//
// ALTER MATERIALIZED VIEW LOG's page lists alter_table_partitioning, but
// Oracle reads a narrower set of partition operations on a log than ALTER
// TABLE takes, each with fewer options: COALESCE, DROP, TRUNCATE and
// EXCHANGE are refused outright (ORA-12045), and the operations it does
// take refuse table_compression, the inmemory_clause, LOB parameters and
// the index maintenance clauses (ORA-14049, 14020, 22853, 30566). The
// readers below follow what the pinned server takes, which
// oracle/parser/testdata/compat.tsv records; where Oracle reads a clause
// only to refuse what follows it (UPDATE INDEXES, MAPPING TABLE, INMEMORY),
// the reader does the same, so the error lands where Oracle's does.

// parseCreateMviewLogStmt parses a CREATE MATERIALIZED VIEW LOG statement.
// Called after CREATE MATERIALIZED VIEW LOG has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-VIEW-LOG.html
//
//	CREATE MATERIALIZED VIEW LOG
//	    [ IF NOT EXISTS ]
//	    ON [ schema. ] table
//	    [ SHARING = { METADATA | NONE } ]
//	    [ { physical_attributes_clause | TABLESPACE tablespace
//	      | logging_clause | { CACHE | NOCACHE } | parallel_clause
//	      | table_partitioning_clauses }... ]
//	    [ WITH with_item [, with_item ]... ]
//	    [ new_values_clause ]
//	    [ mv_log_purge_clause ]
//	    [ for_refresh_clause ]
//
//	with_item: { OBJECT ID | PRIMARY KEY | ROWID | SEQUENCE | COMMIT SCN }
//	               [ ( column [, column ]... ) ]
//	         | ( column [, column ]... )
//
// Oracle takes the attributes in any order, table_partitioning_clauses
// among them, and a column list after any item of WITH (WITH ROWID (a),
// SEQUENCE).
func (p *Parser) parseCreateMviewLogStmt(start int) (*nodes.CreateMviewLogStmt, error) {
	stmt := &nodes.CreateMviewLogStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF NOT EXISTS: Oracle reports an incomplete one at IF (ORA-11543).
	if p.cur.Type == kwIF {
		toks := p.peekAhead(2)
		if toks[0].Type == kwNOT && toks[1].Type == kwEXISTS {
			p.advance()
			p.advance()
			p.advance()
			stmt.IfNotExists = true
		}
	}

	onTable, err := p.parseMVLogOnTable()
	if err != nil {
		return nil, err
	}
	stmt.OnTable = onTable

	// SHARING = { METADATA | NONE }
	if p.atWord("SHARING") {
		p.advance()
		if err := p.atExpect('='); err != nil {
			return nil, err
		}
		switch {
		case p.atWord("METADATA"), p.atWord("NONE"):
			stmt.Sharing = p.cur.Str
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	partitioned := false
	for {
		if p.isMVLogCreateAttrStart() {
			if err := p.parseMVLogCreateAttr(); err != nil {
				return nil, err
			}
			continue
		}
		// table_partitioning_clauses, once (ORA-14027 at a second).
		if partitioned || p.cur.Type != kwPARTITION {
			break
		}
		if err := p.parseMVLogPartitioning(); err != nil {
			return nil, err
		}
		partitioned = true
	}

	if p.cur.Type == kwWITH {
		p.advance()
		if err := p.parseMVLogWithItems(stmt); err != nil {
			return nil, err
		}
	}

	if p.isMVLogNewValuesStart() {
		if p.atWord("INCLUDING") {
			stmt.Including = true
		} else {
			stmt.Excluding = true
		}
		if err := p.parseMVLogNewValues(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwPURGE {
		if err := p.parseMVLogPurge(stmt, true); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwFOR {
		if err := p.parseMVLogForRefresh(&stmt.ForRefresh, &stmt.StagingLog); err != nil {
			return nil, err
		}
	}

	if p.atWord("INMEMORY") || p.atWord("NOMAPPING") || p.cur.Type == kwIMMEDIATE {
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMVLogPartitioning reads table_partitioning_clauses with CREATE
// TABLE's reader, as for a heap table with a column list, narrowed to what
// Oracle takes on a log (createTablePartitionsOfLog). The AST keeps no
// partitioning of a log.
//
// Oracle takes PARTITION BY REFERENCE ( constraint ) with no partition
// list after it, and fails at the constraint where something else follows
// (ORA-14655). It reads CONSISTENT HASH only to refuse it (ORA-14151 at
// what follows).
func (p *Parser) parseMVLogPartitioning() error {
	ahead := p.peekAhead(3)
	switch {
	case ahead[0].Type == kwBY && ahead[1].Type == kwREFERENCE:
		p.advance() // PARTITION
		p.advance() // BY
		p.advance() // REFERENCE
		if err := p.atExpect('('); err != nil {
			return err
		}
		constraint := p.cur
		if err := p.parseATName(); err != nil {
			return err
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		if !p.isMVLogCreateAttrStart() && !p.isMVLogCreateTailStart() {
			return p.syntaxErrorAtTok(constraint)
		}
		return nil
	case ahead[0].Type == kwBY && isCreateTableWordTok(ahead[1], "CONSISTENT") && isCreateTableWordTok(ahead[2], "HASH"):
		p.advance() // PARTITION
		p.advance() // BY
		p.advance() // CONSISTENT
		p.advance() // HASH
		return p.mvLogRefusal("consistent hash partitioning")
	}
	defer func(saved bool) { p.createTablePartitionsOfLog = saved }(p.createTablePartitionsOfLog)
	p.createTablePartitionsOfLog = true
	table := &nodes.CreateTableStmt{}
	partition, err := p.parsePartitionClause(table, true)
	if err != nil {
		return err
	}
	if partition == nil {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// isMVLogCreateTailStart reports what may follow CREATE MATERIALIZED VIEW
// LOG's attributes: WITH, new_values_clause, mv_log_purge_clause,
// for_refresh_clause or the end of the statement.
func (p *Parser) isMVLogCreateTailStart() bool {
	switch p.cur.Type {
	case kwWITH, kwPURGE, kwFOR, ';', tokEOF:
		return true
	}
	return p.isMVLogNewValuesStart()
}

// parseMVLogOnTable reads ON [ schema. ] table.
func (p *Parser) parseMVLogOnTable() (*nodes.ObjectName, error) {
	if err := p.atExpect(kwON); err != nil {
		return nil, err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if name == nil || name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	return name, nil
}

// isMVLogPhysicalAttrStart reports a physical_attributes_clause item:
// PCTFREE, PCTUSED, INITRANS, MAXTRANS or STORAGE. The ALTER page lists no
// MAXTRANS, but Oracle takes it there as it does in CREATE.
func (p *Parser) isMVLogPhysicalAttrStart() bool {
	return p.atWord("PCTFREE") || p.atWord("PCTUSED") || p.atWord("INITRANS") || p.atWord("MAXTRANS") ||
		p.atWord("STORAGE")
}

// parseMVLogPhysicalAttr reads one physical_attributes_clause item:
// { PCTFREE | PCTUSED | INITRANS | MAXTRANS } integer | storage_clause.
func (p *Parser) parseMVLogPhysicalAttr() error {
	if p.atWord("STORAGE") {
		return p.parseATStorageClause()
	}
	p.advance()
	return p.parseATInteger()
}

// isMVLogLoggingStart reports logging_clause as Oracle takes it on a log:
// LOGGING or NOLOGGING. It refuses FILESYSTEM_LIKE_LOGGING (ORA-12044 in
// CREATE, ORA-12045 in ALTER).
func (p *Parser) isMVLogLoggingStart() bool {
	return p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING
}

// isMVLogParallelStart reports parallel_clause.
func (p *Parser) isMVLogParallelStart() bool {
	return p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL
}

// isMVLogCacheStart reports CACHE or NOCACHE.
func (p *Parser) isMVLogCacheStart() bool {
	return p.cur.Type == kwCACHE || p.cur.Type == kwNOCACHE
}

// isMVLogCreateAttrStart reports one of CREATE MATERIALIZED VIEW LOG's
// attributes: physical_attributes_clause, TABLESPACE, logging_clause,
// CACHE, NOCACHE or parallel_clause.
func (p *Parser) isMVLogCreateAttrStart() bool {
	return p.isMVLogPhysicalAttrStart() || p.cur.Type == kwTABLESPACE || p.isMVLogLoggingStart() ||
		p.isMVLogCacheStart() || p.isMVLogParallelStart()
}

// parseMVLogCreateAttr reads one attribute isMVLogCreateAttrStart reports.
func (p *Parser) parseMVLogCreateAttr() error {
	switch {
	case p.isMVLogPhysicalAttrStart():
		return p.parseMVLogPhysicalAttr()
	case p.cur.Type == kwTABLESPACE:
		p.advance()
		return p.parseATName()
	case p.isMVLogParallelStart():
		return p.parseATParallel()
	}
	p.advance() // LOGGING, NOLOGGING, CACHE, NOCACHE
	return nil
}

// isMVLogWithKeywordStart reports an item of WITH or of ALTER's ADD that is
// a keyword: OBJECT ID, PRIMARY KEY, ROWID or SEQUENCE.
func (p *Parser) isMVLogWithKeywordStart() bool {
	return p.cur.Type == kwPRIMARY || p.cur.Type == kwROWID || p.cur.Type == kwSEQUENCE || p.atWord("OBJECT")
}

// parseMVLogWithKeyword reads OBJECT ID, PRIMARY KEY, ROWID or SEQUENCE and
// returns which.
func (p *Parser) parseMVLogWithKeyword() (string, error) {
	switch {
	case p.cur.Type == kwPRIMARY:
		p.advance()
		if err := p.atExpect(kwKEY); err != nil {
			return "", err
		}
		return "PRIMARY_KEY", nil
	case p.cur.Type == kwROWID:
		p.advance()
		return "ROWID", nil
	case p.cur.Type == kwSEQUENCE:
		p.advance()
		return "SEQUENCE", nil
	}
	p.advance() // OBJECT
	if err := p.atExpectWord("ID"); err != nil {
		return "", err
	}
	return "OBJECT_ID", nil
}

// parseMVLogWithItems reads, after WITH, with_item [, with_item ]...
func (p *Parser) parseMVLogWithItems(stmt *nodes.CreateMviewLogStmt) error {
	for {
		switch {
		case p.cur.Type == kwCOMMIT:
			p.advance()
			if err := p.atExpectWord("SCN"); err != nil {
				return err
			}
			stmt.CommitSCN = true
		case p.isMVLogWithKeywordStart():
			what, err := p.parseMVLogWithKeyword()
			if err != nil {
				return err
			}
			switch what {
			case "PRIMARY_KEY":
				stmt.WithPK = true
			case "ROWID":
				stmt.WithRowID = true
			case "SEQUENCE":
				stmt.WithSeq = true
			case "OBJECT_ID":
				stmt.WithOID = true
			}
		case p.cur.Type != '(':
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == '(' {
			if stmt.Columns == nil {
				stmt.Columns = &nodes.List{}
			}
			if err := p.parseMVLogColumnList(stmt.Columns); err != nil {
				return err
			}
		}
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseMVLogColumnList reads ( column [, column ]... ) into list.
func (p *Parser) parseMVLogColumnList(list *nodes.List) error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		name, err := p.parseATNameString()
		if err != nil {
			return err
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// isMVLogNewValuesStart reports new_values_clause.
func (p *Parser) isMVLogNewValuesStart() bool {
	return p.atWord("INCLUDING") || p.atWord("EXCLUDING")
}

// parseMVLogNewValues reads new_values_clause:
// { INCLUDING | EXCLUDING } NEW VALUES.
func (p *Parser) parseMVLogNewValues() error {
	p.advance() // INCLUDING or EXCLUDING
	if err := p.atExpectWord("NEW"); err != nil {
		return err
	}
	return p.atExpect(kwVALUES)
}

// parseMVLogPurge reads mv_log_purge_clause:
//
//	PURGE { IMMEDIATE [ SYNCHRONOUS | ASYNCHRONOUS ]
//	      | START WITH datetime_expr [ NEXT datetime_expr ]
//	          [ REPEAT INTERVAL interval_expr ] }
//
// REPEAT INTERVAL stands in CREATE only: ALTER refuses it (ORA-03049).
// CREATE's page has IMMEDIATE need SYNCHRONOUS or ASYNCHRONOUS, but Oracle
// takes it alone in both statements, as synchronous.
func (p *Parser) parseMVLogPurge(stmt *nodes.CreateMviewLogStmt, create bool) error {
	p.advance() // PURGE
	if p.cur.Type == kwIMMEDIATE {
		p.advance()
		mode := "IMMEDIATE_SYNC"
		if p.atWord("SYNCHRONOUS") {
			p.advance()
		} else if p.atWord("ASYNCHRONOUS") {
			p.advance()
			mode = "IMMEDIATE_ASYNC"
		}
		if stmt != nil {
			stmt.PurgeMode = mode
		}
		return nil
	}
	if err := p.atExpect(kwSTART); err != nil {
		return err
	}
	if err := p.atExpect(kwWITH); err != nil {
		return err
	}
	startExpr, err := p.parseMVLogExpr()
	if err != nil {
		return err
	}
	var next nodes.ExprNode
	if p.cur.Type == kwNEXT {
		p.advance()
		if next, err = p.parseMVLogExpr(); err != nil {
			return err
		}
	}
	if create && p.atWord("REPEAT") {
		p.advance()
		if err := p.atExpect(kwINTERVAL); err != nil {
			return err
		}
		if err := p.parseMVLogExprOnly(); err != nil {
			return err
		}
	}
	if stmt != nil {
		stmt.PurgeMode = "START_WITH"
		stmt.PurgeStart = startExpr
		stmt.PurgeNext = next
	}
	return nil
}

// parseMVLogExpr reads an expression that must be there.
func (p *Parser) parseMVLogExpr() (nodes.ExprNode, error) {
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return e, nil
}

// parseMVLogExprOnly reads an expression that must be there, where the AST
// keeps none.
func (p *Parser) parseMVLogExprOnly() error {
	e, err := p.parseExpr()
	if err != nil {
		return err
	}
	if e == nil {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseMVLogForRefresh reads for_refresh_clause into its mode and staging
// log:
//
//	FOR { FAST REFRESH | SYNCHRONOUS REFRESH USING staging_log }
func (p *Parser) parseMVLogForRefresh(mode, staging *string) error {
	p.advance() // FOR
	if p.atWord("FAST") {
		p.advance()
		*mode = "FAST"
		return p.atExpect(kwREFRESH)
	}
	if err := p.atExpectWord("SYNCHRONOUS"); err != nil {
		return err
	}
	if err := p.atExpect(kwREFRESH); err != nil {
		return err
	}
	if err := p.atExpect(kwUSING); err != nil {
		return err
	}
	name, err := p.parseATNameString()
	if err != nil {
		return err
	}
	*mode, *staging = "SYNCHRONOUS", name
	return nil
}

// parseAlterMviewLogStmt parses an ALTER MATERIALIZED VIEW LOG statement.
// Called after ALTER MATERIALIZED VIEW LOG has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-VIEW-LOG.html
//
//	ALTER MATERIALIZED VIEW LOG [ IF EXISTS ] [ FORCE ]
//	    ON [ schema. ] table
//	    [ mv_log_options
//	    | MOVE ... | MODIFY ... | RENAME ... | SET ... | SPLIT ... | MERGE ...
//	    | ADD PARTITION ... ]
//
// Oracle takes the statement with nothing after the table.
func (p *Parser) parseAlterMviewLogStmt(start int) (*nodes.AlterMviewLogStmt, error) {
	stmt := &nodes.AlterMviewLogStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	if p.cur.Type == kwFORCE {
		stmt.Force = true
		p.advance()
	}

	onTable, err := p.parseMVLogOnTable()
	if err != nil {
		return nil, err
	}
	stmt.OnTable = onTable

	switch {
	case p.cur.Type == kwMODIFY:
		err = p.parseMVLogModify(stmt)
	case p.cur.Type == kwRENAME:
		err = p.parseMVLogRename(stmt)
	case p.cur.Type == kwSET:
		err = p.parseMVLogSet(stmt)
	case p.atWord("MOVE"):
		err = p.parseMVLogMove(stmt)
	case p.atWord("SPLIT"):
		err = p.parseMVLogSplit(stmt)
	case p.atWord("MERGE"):
		err = p.parseMVLogMerge(stmt)
	case p.cur.Type == kwADD && p.peekNext().Type == kwPARTITION:
		err = p.parseMVLogAddPartitions(stmt)
	case p.atWord("INMEMORY") || p.atWord("NOMAPPING") || p.cur.Type == kwIMMEDIATE:
		p.advance()
		return nil, p.syntaxErrorAtCur()
	default:
		err = p.parseMVLogOptions(stmt)
	}
	if err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// mvLogRefusal fails at the current token, after a clause Oracle reads on
// a materialized view log only to refuse it or what follows it.
func (p *Parser) mvLogRefusal(clause string) error {
	err := p.syntaxErrorAtCur()
	if p.lexer.Err == nil {
		err.Message += ": " + clause + " is not supported on a materialized view log"
	}
	return err
}

// isMVLogIndexMaintenanceStart reports update_index_clauses: UPDATE or
// INVALIDATE.
func (p *Parser) isMVLogIndexMaintenanceStart() bool {
	return p.cur.Type == kwUPDATE || p.atWord("INVALIDATE")
}

// refuseMVLogIndexMaintenance reads the head of update_index_clauses,
// { UPDATE | INVALIDATE } GLOBAL INDEXES or UPDATE INDEXES, and fails at
// what follows: Oracle allows no index maintenance on a log's partitions
// (ORA-30566, at the token after the clause's head).
func (p *Parser) refuseMVLogIndexMaintenance() error {
	update := p.cur.Type == kwUPDATE
	p.advance() // UPDATE or INVALIDATE
	if update && p.atWord("INDEXES") {
		p.advance()
		return p.mvLogRefusal("an index maintenance clause")
	}
	if err := p.atExpectWord("GLOBAL"); err != nil {
		return err
	}
	if err := p.atExpectWord("INDEXES"); err != nil {
		return err
	}
	return p.mvLogRefusal("an index maintenance clause")
}

// isMVLogAttrStart reports an attribute ALTER MATERIALIZED VIEW LOG takes
// before its other options: physical_attributes_clause, logging_clause,
// CACHE, NOCACHE, parallel_clause, allocate_extent_clause,
// deallocate_unused_clause or shrink_clause.
func (p *Parser) isMVLogAttrStart() bool {
	return p.isMVLogPhysicalAttrStart() || p.isMVLogLoggingStart() || p.isMVLogCacheStart() ||
		p.isMVLogParallelStart() || p.isATAllocateExtentStart() || p.isATDeallocateStart() || p.isATShrinkStart()
}

// parseMVLogShrink reads shrink_clause: SHRINK SPACE [ COMPACT ] [ CASCADE ],
// which Oracle takes in either order.
func (p *Parser) parseMVLogShrink() error {
	p.advance() // SHRINK
	if err := p.atExpectWord("SPACE"); err != nil {
		return err
	}
	compact, cascade := false, false
	for {
		switch {
		case !compact && p.atWord("COMPACT"):
			compact = true
		case !cascade && p.cur.Type == kwCASCADE:
			cascade = true
		default:
			return nil
		}
		p.advance()
	}
}

// parseMVLogOptions reads what ALTER MATERIALIZED VIEW LOG takes besides
// its MOVE, MODIFY and partition operations, in this order:
//
//	[ { physical_attributes_clause | logging_clause | { CACHE | NOCACHE }
//	  | parallel_clause | allocate_extent_clause
//	  | deallocate_unused_clause }... [ shrink_clause ] ]
//	[ mv_log_augmentation ] [ mv_log_purge_clause ] [ for_refresh_clause ]
//
// shrink_clause ends the statement (ORA-10630 for anything after it), and
// so does mv_log_augmentation without a new_values_clause (ORA-02000).
func (p *Parser) parseMVLogOptions(stmt *nodes.AlterMviewLogStmt) error {
	for p.isMVLogAttrStart() {
		action := ""
		var err error
		switch {
		case p.isMVLogPhysicalAttrStart():
			action = p.cur.Str
			err = p.parseMVLogPhysicalAttr()
		case p.isMVLogParallelStart():
			action = p.cur.Str
			err = p.parseATParallel()
		case p.isATAllocateExtentStart():
			action = "ALLOCATE_EXTENT"
			err = p.parseATAllocateExtent()
		case p.isATDeallocateStart():
			action = "DEALLOCATE_UNUSED"
			err = p.parseATDeallocateUnused()
		case p.isATShrinkStart():
			if stmt.Action == "" {
				stmt.Action = "SHRINK"
			}
			return p.parseMVLogShrink()
		default:
			action = p.cur.Str // LOGGING, NOLOGGING, CACHE, NOCACHE
			p.advance()
		}
		if err != nil {
			return err
		}
		if stmt.Action == "" {
			stmt.Action = action
		}
	}

	if p.cur.Type == kwADD {
		newValues, err := p.parseMVLogAugmentation(stmt)
		if err != nil || !newValues {
			return err
		}
	}

	if p.cur.Type == kwPURGE {
		if stmt.Action == "" {
			stmt.Action = "PURGE"
		}
		if err := p.parseMVLogPurge(nil, false); err != nil {
			return err
		}
	}

	if p.cur.Type == kwFOR {
		if stmt.Action == "" {
			stmt.Action = "FOR_REFRESH"
		}
		var mode, staging string
		if err := p.parseMVLogForRefresh(&mode, &staging); err != nil {
			return err
		}
	}
	return nil
}

// parseMVLogAugmentation reads mv_log_augmentation and reports whether it
// ended with a new_values_clause:
//
//	ADD [ { OBJECT ID | PRIMARY KEY | ROWID | SEQUENCE }
//	      [, { OBJECT ID | PRIMARY KEY | ROWID | SEQUENCE } ]... ]
//	    [ [,] ( column [, column ]... ) ]
//	    [ new_values_clause ]
//
// The column list ends the list (ORA-03048 for a comma after it), and
// Oracle takes ADD with nothing after it.
func (p *Parser) parseMVLogAugmentation(stmt *nodes.AlterMviewLogStmt) (bool, error) {
	p.advance() // ADD
	action := "ADD"
	for p.isMVLogWithKeywordStart() {
		what, err := p.parseMVLogWithKeyword()
		if err != nil {
			return false, err
		}
		if action == "ADD" {
			action = "ADD_" + what
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		if p.cur.Type == '(' {
			break
		}
		if !p.isMVLogWithKeywordStart() {
			return false, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == '(' {
		stmt.Columns = &nodes.List{}
		if err := p.parseMVLogColumnList(stmt.Columns); err != nil {
			return false, err
		}
	}
	if stmt.Action == "" {
		stmt.Action = action
	}
	if !p.isMVLogNewValuesStart() {
		return false, nil
	}
	return true, p.parseMVLogNewValues()
}

// isMVLogSegmentAttrStart reports an attribute of a moved or new partition
// of a log: physical_attributes_clause, TABLESPACE or logging_clause.
func (p *Parser) isMVLogSegmentAttrStart() bool {
	return p.isMVLogPhysicalAttrStart() || p.cur.Type == kwTABLESPACE || p.isMVLogLoggingStart()
}

// isMVLogSubpartAttrStart reports an attribute of a moved or new
// subpartition: TABLESPACE or logging_clause. Oracle refuses physical
// attributes for a subpartition (ORA-14160).
func (p *Parser) isMVLogSubpartAttrStart() bool {
	return p.cur.Type == kwTABLESPACE || p.isMVLogLoggingStart()
}

// parseMVLogSegmentAttr reads one attribute isMVLogSegmentAttrStart
// reports.
func (p *Parser) parseMVLogSegmentAttr() error {
	switch {
	case p.isMVLogPhysicalAttrStart():
		return p.parseMVLogPhysicalAttr()
	case p.cur.Type == kwTABLESPACE:
		p.advance()
		return p.parseATName()
	}
	p.advance() // LOGGING or NOLOGGING
	return nil
}

// parseMVLogMove reads MOVE: the log's own segment or one of its partitions.
//
//	MOVE [ { physical_attributes_clause | TABLESPACE tablespace
//	       | logging_clause | parallel_clause }... ]
//	MOVE PARTITION partition [ { physical_attributes_clause
//	       | TABLESPACE tablespace | logging_clause | parallel_clause }... ]
//	MOVE SUBPARTITION subpartition [ { TABLESPACE tablespace
//	       | logging_clause | parallel_clause }... ]
//
// MAPPING TABLE and update_index_clauses are read and refused (ORA-14020,
// ORA-30566).
func (p *Parser) parseMVLogMove(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // MOVE
	isAttr := p.isMVLogSegmentAttrStart
	switch p.cur.Type {
	case kwPARTITION:
		stmt.Action = "MOVE_PARTITION"
	case kwSUBPARTITION:
		stmt.Action = "MOVE_SUBPARTITION"
		isAttr = p.isMVLogSubpartAttrStart
	default:
		stmt.Action = "MOVE"
	}
	if stmt.Action != "MOVE" {
		p.advance()
		if err := p.parseATName(); err != nil {
			return err
		}
	}
	for {
		switch {
		case isAttr():
			if err := p.parseMVLogSegmentAttr(); err != nil {
				return err
			}
		case p.isMVLogParallelStart():
			if err := p.parseATParallel(); err != nil {
				return err
			}
		case stmt.Action == "MOVE_PARTITION" && p.atWord("MAPPING"):
			p.advance()
			if err := p.atExpect(kwTABLE); err != nil {
				return err
			}
			return p.mvLogRefusal("MAPPING TABLE")
		case stmt.Action != "MOVE" && p.isMVLogIndexMaintenanceStart():
			return p.refuseMVLogIndexMaintenance()
		default:
			return nil
		}
	}
}

// parseMVLogModify reads MODIFY:
//
//	MODIFY DEFAULT ATTRIBUTES ...
//	MODIFY PARTITION partition ...
//	MODIFY SUBPARTITION subpartition ...
//	MODIFY ( column_definition [, column_definition ]... )
//	MODIFY column_definition
//
// Oracle reads ALTER TABLE's column modification on a log too, and takes a
// column with nothing after it there (MODIFY a, MODIFY (a)).
func (p *Parser) parseMVLogModify(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // MODIFY
	switch {
	case p.cur.Type == kwDEFAULT:
		stmt.Action = "MODIFY_DEFAULT_ATTRIBUTES"
		return p.parseMVLogModifyDefaultAttrs()
	case p.cur.Type == kwPARTITION:
		stmt.Action = "MODIFY_PARTITION"
		return p.parseMVLogModifyPartition()
	case p.cur.Type == kwSUBPARTITION:
		stmt.Action = "MODIFY_SUBPARTITION"
		return p.parseMVLogModifySubpartition()
	}
	stmt.Action = "MODIFY_COLUMN"
	paren := p.cur.Type == '('
	if paren {
		p.advance()
	}
	for {
		if err := p.parseMVLogModifyColumn(); err != nil {
			return err
		}
		if !paren || p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if paren {
		return p.atExpect(')')
	}
	return nil
}

// parseMVLogModifyColumn reads one column of MODIFY: a column_definition,
// or the column's name alone where the column ends there.
func (p *Parser) parseMVLogModifyColumn() error {
	if p.isATNameStart() {
		switch p.peekNext().Type {
		case ';', tokEOF, ',', ')':
			p.advance()
			return nil
		}
	}
	col, err := p.parseColumnDef()
	if err != nil {
		return err
	}
	if col == nil {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseMVLogModifyDefaultAttrs reads, after MODIFY,
//
//	DEFAULT ATTRIBUTES [ FOR PARTITION partition ]
//	    [ { physical_attributes_clause | TABLESPACE tablespace
//	      | logging_clause | parallel_clause }... ]
//
// parallel_clause stands for the table only (ORA-14020 after FOR
// PARTITION), and INMEMORY is read and refused (ORA-12045 after it).
func (p *Parser) parseMVLogModifyDefaultAttrs() error {
	p.advance() // DEFAULT
	if err := p.atExpectWord("ATTRIBUTES"); err != nil {
		return err
	}
	forPartition := false
	if p.cur.Type == kwFOR {
		p.advance()
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		forPartition = true
	}
	for {
		switch {
		case p.isMVLogSegmentAttrStart():
			if err := p.parseMVLogSegmentAttr(); err != nil {
				return err
			}
		case !forPartition && p.isMVLogParallelStart():
			if err := p.parseATParallel(); err != nil {
				return err
			}
		case p.atWord("INMEMORY"):
			p.advance()
			return p.mvLogRefusal("the inmemory_clause")
		default:
			return nil
		}
	}
}

// isMVLogPartitionAttrStart reports an item of partition_attributes as
// Oracle takes it on a log's partition: physical_attributes_clause,
// logging_clause, allocate_extent_clause, deallocate_unused_clause or
// shrink_clause.
func (p *Parser) isMVLogPartitionAttrStart() bool {
	return p.isMVLogPhysicalAttrStart() || p.isMVLogLoggingStart() || p.isATAllocateExtentStart() ||
		p.isATDeallocateStart() || p.isATShrinkStart()
}

// isMVLogSubpartitionAttrStart reports an attribute of a log's
// subpartition: logging_clause, allocate_extent_clause,
// deallocate_unused_clause or shrink_clause. Oracle refuses physical
// attributes there (ORA-14169).
func (p *Parser) isMVLogSubpartitionAttrStart() bool {
	return p.isMVLogLoggingStart() || p.isATAllocateExtentStart() || p.isATDeallocateStart() ||
		p.isATShrinkStart()
}

// parseMVLogPartitionAttrs reads the attributes isAttr reports, in any
// order; shrink_clause comes last (ORA-10630 for anything after it).
func (p *Parser) parseMVLogPartitionAttrs(isAttr func() bool) error {
	for isAttr() {
		var err error
		switch {
		case p.isMVLogPhysicalAttrStart():
			err = p.parseMVLogPhysicalAttr()
		case p.isATAllocateExtentStart():
			err = p.parseATAllocateExtent()
		case p.isATDeallocateStart():
			err = p.parseATDeallocateUnused()
		case p.isATShrinkStart():
			return p.parseMVLogShrink()
		default:
			p.advance() // LOGGING or NOLOGGING
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// parseMVLogModifyPartition reads, after MODIFY,
//
//	PARTITION partition
//	    [ partition_attributes
//	    | { ADD | DROP } VALUES ( value [, value ]... )
//	    | ADD subpartition_desc [, subpartition_desc ]...
//	    | COALESCE [ SUBPARTITION ] [ parallel_clause ]
//	    | [ REBUILD ] UNUSABLE LOCAL INDEXES ]
//
// subpartition_desc: SUBPARTITION [ subpartition ]
// [ range_values_clause | list_values_clause ]. Oracle refuses the
// inmemory_clause (ORA-14049), LOB parameters (ORA-22853) and VARRAY
// storage (ORA-00905) of a log's partition after reading their heads.
func (p *Parser) parseMVLogModifyPartition() error {
	p.advance() // PARTITION
	if err := p.parseATName(); err != nil {
		return err
	}
	switch {
	case p.isMVLogPartitionAttrStart():
		return p.parseMVLogPartitionAttrs(p.isMVLogPartitionAttrStart)
	case p.cur.Type == kwADD:
		p.advance()
		if p.cur.Type == kwVALUES {
			return p.parseMVLogPartitionValues()
		}
		return p.parseMVLogAddSubpartitions()
	case p.cur.Type == kwDROP:
		p.advance()
		return p.parseMVLogPartitionValues()
	case p.atWord("COALESCE"):
		p.advance()
		if p.cur.Type == kwSUBPARTITION {
			p.advance()
		}
		if p.isMVLogIndexMaintenanceStart() {
			return p.refuseMVLogIndexMaintenance()
		}
		if p.isMVLogParallelStart() {
			return p.parseATParallel()
		}
		return nil
	case p.atWord("REBUILD"), p.atWord("UNUSABLE"):
		return p.parseATUnusableLocalIndexes()
	case p.atWord("INMEMORY"):
		p.advance()
		return p.mvLogRefusal("the inmemory_clause")
	case p.atWord("NO") && p.atNextWord("INMEMORY"):
		p.advance()
		p.advance()
		return p.mvLogRefusal("the inmemory_clause")
	}
	return p.parseMVLogLobStorage()
}

// parseMVLogModifySubpartition reads, after MODIFY,
//
//	SUBPARTITION subpartition
//	    [ { logging_clause | allocate_extent_clause
//	      | deallocate_unused_clause | shrink_clause }...
//	    | { ADD | DROP } VALUES ( value [, value ]... )
//	    | [ REBUILD ] UNUSABLE LOCAL INDEXES ]
//
// LOB parameters and VARRAY storage are read and refused, as for a
// partition (ORA-22877, ORA-00905).
func (p *Parser) parseMVLogModifySubpartition() error {
	p.advance() // SUBPARTITION
	if err := p.parseATName(); err != nil {
		return err
	}
	switch {
	case p.isMVLogSubpartitionAttrStart():
		return p.parseMVLogPartitionAttrs(p.isMVLogSubpartitionAttrStart)
	case p.cur.Type == kwADD, p.cur.Type == kwDROP:
		p.advance()
		return p.parseMVLogPartitionValues()
	case p.atWord("REBUILD"), p.atWord("UNUSABLE"):
		return p.parseATUnusableLocalIndexes()
	}
	return p.parseMVLogLobStorage()
}

// parseMVLogLobStorage reads, where a log's partition or subpartition may
// name its LOB or VARRAY storage, LOB ( lob_item ) ( or VARRAY varray, and
// fails at what follows: Oracle takes no LOB parameters there. With
// neither, it reads nothing.
func (p *Parser) parseMVLogLobStorage() error {
	if p.atWord("LOB") {
		p.advance()
		if err := p.atExpect('('); err != nil {
			return err
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		if err := p.atExpect('('); err != nil {
			return err
		}
		return p.mvLogRefusal("LOB storage")
	}
	if p.atWord("VARRAY") {
		p.advance()
		if err := p.parseATName(); err != nil {
			return err
		}
		return p.mvLogRefusal("VARRAY storage")
	}
	return nil
}

// parseMVLogPartitionValues reads VALUES ( value [, value ]... ) after ADD
// or DROP.
func (p *Parser) parseMVLogPartitionValues() error {
	if err := p.atExpect(kwVALUES); err != nil {
		return err
	}
	return p.parseATBoundValues(false)
}

// parseMVLogAddSubpartitions reads, after MODIFY PARTITION partition ADD,
//
//	subpartition_desc [, subpartition_desc ]... [ parallel_clause ]
//
// each SUBPARTITION [ subpartition ] [ range_values_clause |
// list_values_clause ]. parallel_clause follows hash subpartitions only,
// which carry no values (add_hash_subpartition). Where the first
// subpartition's name is followed by what cannot follow it, Oracle fails
// at the name (ADD SUBPARTITION a zzz, ORA-14020 at a); after a later one,
// at what follows. update_index_clauses are read and refused (ORA-30566).
func (p *Parser) parseMVLogAddSubpartitions() error {
	values := false
	named := false
	for first := true; ; first = false {
		if err := p.atExpect(kwSUBPARTITION); err != nil {
			return err
		}
		if p.isATNameStart() && !p.isMVLogParallelStart() {
			name := p.cur
			p.advance()
			named = true
			if first && !p.isMVLogAddSubpartitionFollow() {
				return p.syntaxErrorAtTok(name)
			}
		}
		if p.cur.Type == kwVALUES {
			if err := p.parseATValuesClause(); err != nil {
				return err
			}
			values = true
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.isMVLogIndexMaintenanceStart() {
		if !values && !named {
			return p.syntaxErrorAtCur()
		}
		return p.refuseMVLogIndexMaintenance()
	}
	if !values && p.isMVLogParallelStart() {
		return p.parseATParallel()
	}
	return nil
}

// isMVLogAddSubpartitionFollow reports what may follow a subpartition's
// name in ADD: its values, the next subpartition, update_index_clauses,
// parallel_clause or the end of the statement.
func (p *Parser) isMVLogAddSubpartitionFollow() bool {
	switch p.cur.Type {
	case kwVALUES, ',', ';', tokEOF:
		return true
	}
	return p.isMVLogIndexMaintenanceStart() || p.isMVLogParallelStart()
}

// parseMVLogAddPartitions reads
//
//	ADD add_partition [, add_partition ]...
//
//	add_partition:
//	    PARTITION [ partition ] [ range_values_clause | list_values_clause ]
//	    [ { physical_attributes_clause | TABLESPACE tablespace
//	      | logging_clause | parallel_clause }... ]
//	    [ ( subpartition_desc [, subpartition_desc ]... )
//	    | SUBPARTITIONS integer [ STORE IN ( tablespace [, tablespace ]... ) ] ]
//
// as Oracle takes it on a log. The inmemory_clause and, after the values,
// update_index_clauses are read and refused (ORA-14020 and ORA-30566 at what
// follows); table_compression is refused (ORA-14020).
func (p *Parser) parseMVLogAddPartitions(stmt *nodes.AlterMviewLogStmt) error {
	stmt.Action = "ADD_PARTITION"
	p.advance() // ADD
	for {
		values, err := p.parseMVLogAddPartition()
		if err != nil {
			return err
		}
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type != kwPARTITION {
				return p.syntaxErrorAtCur()
			}
			continue
		}
		if values && p.isMVLogIndexMaintenanceStart() {
			return p.refuseMVLogIndexMaintenance()
		}
		return nil
	}
}

// isMVLogAddPartitionAttrStart reports an attribute of a partition ADD
// adds: physical_attributes_clause, TABLESPACE, logging_clause,
// parallel_clause or the inmemory_clause.
func (p *Parser) isMVLogAddPartitionAttrStart() bool {
	return p.isMVLogSegmentAttrStart() || p.isMVLogParallelStart() || p.atWord("INMEMORY") ||
		p.atWord("NO") && p.atNextWord("INMEMORY")
}

// isMVLogAddPartitionFollow reports what may follow a partition's name in
// ADD: its values, its attributes, its subpartitions, the next partition or
// the end of the statement.
func (p *Parser) isMVLogAddPartitionFollow() bool {
	switch p.cur.Type {
	case kwVALUES, '(', ',', ';', tokEOF:
		return true
	}
	return p.isMVLogAddPartitionAttrStart() || p.atWord("SUBPARTITIONS")
}

// parseMVLogAddPartition reads one add_partition and reports whether it
// had values. Oracle reports a word after PARTITION that can neither name
// the partition nor start its description at PARTITION (ORA-02000 at
// ADD PARTITION COMPRESS), save UPDATE, at itself. A word that starts an
// attribute with more after it is the attribute (ADD PARTITION LOGGING
// zzz fails at zzz); another is the name. Where a name without values is
// followed by attributes Oracle cannot read, it fails at the name
// (ADD PARTITION OVERFLOW TABLESPACE a zzz, ORA-14020 at OVERFLOW).
func (p *Parser) parseMVLogAddPartition() (bool, error) {
	partition := p.cur
	p.advance() // PARTITION
	if p.cur.Type == '(' ||
		p.cur.Type != kwVALUES && p.cur.Type != kwUPDATE && p.cur.Type != tokQIDENT && isOracleSQLReservedKeyword(p.cur) {
		return false, p.syntaxErrorAtTok(partition)
	}
	if next := p.peekNext().Type; p.atWord("INVALIDATE") && next != ';' && next != tokEOF && next != ',' {
		return false, p.refuseMVLogIndexMaintenance()
	}
	var name *Token
	if p.isATNameStart() && p.cur.Type != kwVALUES {
		switch p.peekNext().Type {
		case ';', tokEOF, ',':
			p.advance()
			return false, nil
		}
		if !p.isMVLogAddPartitionAttrStart() {
			tok := p.cur
			name = &tok
			p.advance()
		}
	}
	values := false
	if p.cur.Type == kwVALUES {
		if err := p.parseATValuesClause(); err != nil {
			return false, err
		}
		values = true
	}
	atName := func(err error) error {
		if name != nil && !values {
			return p.syntaxErrorAtTok(*name)
		}
		return err
	}
	for p.isMVLogAddPartitionAttrStart() {
		var err error
		switch {
		case p.isMVLogParallelStart():
			err = p.parseATParallel()
		case p.atWord("INMEMORY"), p.atWord("NO"):
			if p.atWord("NO") {
				p.advance()
			}
			p.advance() // INMEMORY
			err = p.mvLogRefusal("the inmemory_clause")
		default:
			err = p.parseMVLogSegmentAttr()
		}
		if err != nil {
			return false, atName(err)
		}
	}
	switch {
	case p.cur.Type == '(':
		p.advance()
		for {
			if err := p.parseMVLogSubpartitionSpec(true, p.isMVLogSpecListFollow); err != nil {
				return false, err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		return values, p.atExpect(')')
	case p.atWord("SUBPARTITIONS"):
		p.advance()
		if err := p.parseATInteger(); err != nil {
			return false, err
		}
		if p.atWord("STORE") {
			p.advance()
			if err := p.atExpect(kwIN); err != nil {
				return false, err
			}
			if err := p.parseATNameList(); err != nil {
				return false, err
			}
		}
		return values, nil
	}
	switch p.cur.Type {
	case ',', ';', tokEOF:
	default:
		if !values || !p.isMVLogIndexMaintenanceStart() {
			return false, atName(p.syntaxErrorAtCur())
		}
	}
	return values, nil
}

// parseMVLogRename reads RENAME { PARTITION | SUBPARTITION } name TO name.
func (p *Parser) parseMVLogRename(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // RENAME
	switch p.cur.Type {
	case kwPARTITION:
		stmt.Action = "RENAME_PARTITION"
	case kwSUBPARTITION:
		stmt.Action = "RENAME_SUBPARTITION"
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.parseATName(); err != nil {
		return err
	}
	if err := p.atExpect(kwTO); err != nil {
		return err
	}
	return p.parseATName()
}

// parseMVLogSet reads SET:
//
//	SET [ PARTITIONING ] { AUTOMATIC | MANUAL }
//	SET STORE IN ( tablespace [, tablespace ]... )
//	SET INTERVAL ( [ expr ] )
//	SET SUBPARTITION TEMPLATE
//	    { ( [ subpartition_desc [, subpartition_desc ]... ] ) | integer }
//
// subpartition_desc: SUBPARTITION [ subpartition ]
// [ range_values_clause | list_values_clause ]
// [ { TABLESPACE tablespace | logging_clause }... ]. Oracle refuses
// SUBPARTITIONS integer after TEMPLATE, taking the number alone.
func (p *Parser) parseMVLogSet(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // SET
	switch {
	case p.atWord("PARTITIONING"), p.atWord("AUTOMATIC"), p.atWord("MANUAL"):
		stmt.Action = "SET_PARTITIONING"
		if p.atWord("PARTITIONING") {
			p.advance()
		}
		return p.atOneOf("AUTOMATIC", "MANUAL")
	case p.atWord("STORE"):
		stmt.Action = "SET_STORE_IN"
		p.advance()
		if err := p.atExpect(kwIN); err != nil {
			return err
		}
		return p.parseATNameList()
	case p.cur.Type == kwINTERVAL:
		stmt.Action = "SET_INTERVAL"
		p.advance()
		if err := p.atExpect('('); err != nil {
			return err
		}
		if p.cur.Type != ')' {
			if err := p.parseMVLogExprOnly(); err != nil {
				return err
			}
		}
		return p.atExpect(')')
	case p.cur.Type == kwSUBPARTITION:
		stmt.Action = "SET_SUBPARTITION_TEMPLATE"
		p.advance()
		if err := p.atExpectWord("TEMPLATE"); err != nil {
			return err
		}
		if p.cur.Type == tokICONST {
			p.advance()
			return nil
		}
		if err := p.atExpect('('); err != nil {
			return err
		}
		if p.cur.Type == ')' {
			p.advance()
			return nil
		}
		for {
			if err := p.parseMVLogSubpartitionSpec(true, p.isMVLogSpecListFollow); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		return p.atExpect(')')
	}
	return p.syntaxErrorAtCur()
}

// parseMVLogSubpartitionSpec reads SUBPARTITION [ subpartition ], with
// [ range_values_clause | list_values_clause ] where values is set, then
// [ { TABLESPACE tablespace | logging_clause }... ]. follow reports what
// else may come after the subpartition.
func (p *Parser) parseMVLogSubpartitionSpec(values bool, follow func() bool) error {
	if err := p.atExpect(kwSUBPARTITION); err != nil {
		return err
	}
	isRest := p.isMVLogSubpartAttrStart
	if values {
		isRest = func() bool { return p.cur.Type == kwVALUES || p.isMVLogSubpartAttrStart() }
	}
	if err := p.parseMVLogSpecName(isRest, follow); err != nil {
		return err
	}
	if values && p.cur.Type == kwVALUES {
		if err := p.parseATValuesClause(); err != nil {
			return err
		}
	}
	for p.isMVLogSubpartAttrStart() {
		if err := p.parseMVLogSegmentAttr(); err != nil {
			return err
		}
	}
	return nil
}

// parseMVLogPartitionSpec reads PARTITION [ partition ]
// [ { physical_attributes_clause | TABLESPACE tablespace
// | logging_clause }... ]. follow reports what else may come after the
// partition.
func (p *Parser) parseMVLogPartitionSpec(follow func() bool) error {
	if err := p.atExpect(kwPARTITION); err != nil {
		return err
	}
	if err := p.parseMVLogSpecName(p.isMVLogSegmentAttrStart, follow); err != nil {
		return err
	}
	for p.isMVLogSegmentAttrStart() {
		if err := p.parseMVLogSegmentAttr(); err != nil {
			return err
		}
	}
	return nil
}

// parseMVLogSpecName reads the optional name of a partition or
// subpartition spec. A word is the name where what follows it is the
// spec's rest (isRest) or may follow the spec (follow); otherwise Oracle
// fails at the word (MERGE PARTITIONS a, a INTO PARTITION a zzz, ORA-14020
// at the second a). A word that starts the rest, or may follow the spec,
// is not a name.
func (p *Parser) parseMVLogSpecName(isRest, follow func() bool) error {
	if !p.isATNameStart() || isRest() || follow() {
		return nil
	}
	name := p.cur
	p.advance()
	if !isRest() && !follow() {
		return p.syntaxErrorAtTok(name)
	}
	return nil
}

// parseMVLogPartitionOrSubpartitionSpec reads a partition_spec, or a
// subpartition_spec where sub is set.
func (p *Parser) parseMVLogPartitionOrSubpartitionSpec(sub bool, follow func() bool) error {
	if sub {
		return p.parseMVLogSubpartitionSpec(false, follow)
	}
	return p.parseMVLogPartitionSpec(follow)
}

// isMVLogSpecListFollow reports what follows a spec in a parenthesized
// list: a comma or the closing parenthesis.
func (p *Parser) isMVLogSpecListFollow() bool {
	return p.cur.Type == ',' || p.cur.Type == ')'
}

// isMVLogMergeIntoFollow reports what follows MERGE's INTO spec:
// parallel_clause, update_index_clauses or the end of the statement.
func (p *Parser) isMVLogMergeIntoFollow() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF || p.isMVLogParallelStart() || p.isMVLogIndexMaintenanceStart()
}

// parseMVLogSplit reads
//
//	SPLIT { PARTITION | SUBPARTITION } name
//	    { AT ( value [, value ]... ) | VALUES ( value [, value ]... ) }
//	    [ INTO ( spec [, spec ]... ) ] [ parallel_clause ]
//
// with partition_spec or subpartition_spec for spec. Oracle refuses INTO
// without AT or VALUES (ORA-14006 at INTO), and update_index_clauses are
// read and refused (ORA-30566).
func (p *Parser) parseMVLogSplit(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // SPLIT
	sub := false
	switch p.cur.Type {
	case kwPARTITION:
		stmt.Action = "SPLIT_PARTITION"
	case kwSUBPARTITION:
		stmt.Action = "SPLIT_SUBPARTITION"
		sub = true
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.parseATName(); err != nil {
		return err
	}
	if p.atWord("AT") {
		p.advance()
	} else if err := p.atExpect(kwVALUES); err != nil {
		return err
	}
	if err := p.parseATExprList(); err != nil {
		return err
	}
	if p.cur.Type == kwINTO {
		p.advance()
		if p.cur.Type != '(' {
			// Oracle reads the list without its parentheses, and the
			// closing one left over at the end (compat.tsv's
			// ALTER MATERIALIZED VIEW LOG ... AT(1) INTO SUBPARTITION,
			// SUBPARTITION) row).
			for {
				if err := p.parseMVLogPartitionOrSubpartitionSpec(sub, p.isMVLogSpecListFollow); err != nil {
					return err
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type != ')' {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		} else {
			p.advance() // consume (
			for {
				if err := p.parseMVLogPartitionOrSubpartitionSpec(sub, p.isMVLogSpecListFollow); err != nil {
					return err
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
	}
	if p.isMVLogIndexMaintenanceStart() {
		return p.refuseMVLogIndexMaintenance()
	}
	if p.isMVLogParallelStart() {
		return p.parseATParallel()
	}
	return nil
}

// parseMVLogMerge reads
//
//	MERGE { PARTITIONS | SUBPARTITIONS } key { , key | TO key }
//	    [ INTO spec ] [ { attribute | parallel_clause }... ]
//
// key: name | ( value [, value ]... ). Oracle refuses FOR before the
// values (ORA-14006). The attributes are those of the INTO spec,
// partition_spec or subpartition_spec, and come only after it; they and
// parallel_clause stand in any order. update_index_clauses are read and
// refused (ORA-30566).
func (p *Parser) parseMVLogMerge(stmt *nodes.AlterMviewLogStmt) error {
	p.advance() // MERGE
	sub := false
	switch {
	case p.atWord("PARTITIONS"):
		stmt.Action = "MERGE_PARTITIONS"
	case p.atWord("SUBPARTITIONS"):
		stmt.Action = "MERGE_SUBPARTITIONS"
		sub = true
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if err := p.parseMVLogPartitionKey(); err != nil {
		return err
	}
	if p.cur.Type == ',' || p.cur.Type == kwTO {
		p.advance()
	} else {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseMVLogPartitionKey(); err != nil {
		return err
	}
	isAttr := func() bool { return false }
	if p.cur.Type == kwINTO {
		p.advance()
		if err := p.parseMVLogPartitionOrSubpartitionSpec(sub, p.isMVLogMergeIntoFollow); err != nil {
			return err
		}
		isAttr = p.isMVLogSegmentAttrStart
		if sub {
			isAttr = p.isMVLogSubpartAttrStart
		}
	}
	for {
		switch {
		case isAttr():
			if err := p.parseMVLogSegmentAttr(); err != nil {
				return err
			}
		case p.isMVLogParallelStart():
			if err := p.parseATParallel(); err != nil {
				return err
			}
		case p.isMVLogIndexMaintenanceStart():
			return p.refuseMVLogIndexMaintenance()
		default:
			return nil
		}
	}
}

// parseMVLogPartitionKey reads a partition named by MERGE: its name, or its
// key values in parentheses.
func (p *Parser) parseMVLogPartitionKey() error {
	if p.cur.Type == '(' {
		return p.parseATExprList()
	}
	return p.parseATName()
}
