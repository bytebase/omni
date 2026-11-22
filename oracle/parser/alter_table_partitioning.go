package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads alter_table_partitioning, the partition maintenance
// operations of ALTER TABLE (oracle/parser/bnf/ALTER-TABLE.bnf). ALTER
// MATERIALIZED VIEW LOG takes the same operations: a caller checks
// isAlterTablePartitioningStart and calls parseAlterTablePartitioning.

// isAlterTablePartitioningStart reports whether the current token starts
// an alter_table_partitioning operation:
//
//	MODIFY DEFAULT ATTRIBUTES, SET PARTITIONING, SET INTERVAL,
//	SET SUBPARTITION TEMPLATE, MODIFY PARTITION partition,
//	MODIFY SUBPARTITION, MOVE PARTITION, MOVE SUBPARTITION, ADD PARTITION,
//	COALESCE PARTITION, DROP PARTITION, DROP SUBPARTITION, RENAME PARTITION,
//	RENAME SUBPARTITION, TRUNCATE PARTITION, TRUNCATE SUBPARTITION,
//	SPLIT PARTITION, SPLIT SUBPARTITION, MERGE PARTITIONS,
//	MERGE SUBPARTITIONS, EXCHANGE PARTITION, EXCHANGE SUBPARTITION
//
// MODIFY PARTITION BY is modify_to_partitioned, not one of them.
func (p *Parser) isAlterTablePartitioningStart() bool {
	next := p.peekNext()
	switch {
	case p.cur.Type == kwMODIFY:
		switch next.Type {
		case kwPARTITION:
			return !p.atPartitionBy()
		case kwSUBPARTITION, kwDEFAULT:
			return true
		}
	case p.cur.Type == kwSET:
		return next.Type == kwINTERVAL || next.Type == kwSUBPARTITION || p.atNextWord("PARTITIONING") ||
			p.atNextWord("AUTOMATIC") || p.atNextWord("MANUAL") || p.atNextWord("STORE")
	case p.cur.Type == kwADD:
		return next.Type == kwPARTITION
	case p.cur.Type == kwTRUNCATE:
		return !p.alterMview
	case p.cur.Type == kwDROP, p.cur.Type == kwRENAME:
		return next.Type == kwPARTITION || next.Type == kwSUBPARTITION
	case p.atWord("MOVE"):
		return next.Type == kwPARTITION || next.Type == kwSUBPARTITION
	case p.atWord("EXCHANGE"):
		return !p.alterMview || next.Type == kwPARTITION || next.Type == kwSUBPARTITION
	case p.atWord("SPLIT"), p.atWord("MERGE"):
		return true
	case p.atWord("COALESCE"):
		return next.Type == kwPARTITION
	}
	return false
}

// atPartitionBy reports MODIFY PARTITION BY (modify_to_partitioned), with
// the current token on MODIFY.
func (p *Parser) atPartitionBy() bool {
	toks := p.peekAhead(2)
	return toks[0].Type == kwPARTITION && toks[1].Type == kwBY
}

// parseAlterTablePartitioning reads one alter_table_partitioning operation;
// the current token is where isAlterTablePartitioningStart said one starts.
func (p *Parser) parseAlterTablePartitioning() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	cmd := &nodes.AlterTableCmd{Options: &nodes.List{}}
	var err error
	switch {
	case p.cur.Type == kwMODIFY:
		err = p.parseATModifyPartitioning(cmd)
	case p.cur.Type == kwSET:
		err = p.parseATSetPartitioning(cmd)
	case p.cur.Type == kwADD:
		err = p.parseATAddPartition(cmd)
	case p.cur.Type == kwDROP:
		err = p.parseATDropPartition(cmd)
	case p.cur.Type == kwRENAME:
		err = p.parseATRenamePartition(cmd)
	case p.cur.Type == kwTRUNCATE:
		err = p.parseATTruncatePartition(cmd)
	case p.atWord("MOVE"):
		err = p.parseATMovePartition(cmd)
	case p.atWord("SPLIT"):
		err = p.parseATSplitPartition(cmd)
	case p.atWord("MERGE"):
		err = p.parseATMergePartitions(cmd)
	case p.atWord("EXCHANGE"):
		err = p.parseATExchangePartition(cmd)
	case p.atWord("COALESCE"):
		err = p.parseATCoalescePartition(cmd)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	if cmd.Options.Len() == 0 {
		cmd.Options = nil
	}
	cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return cmd, nil
}

// parseATPartitionName reads the name after PARTITION or SUBPARTITION.
func (p *Parser) parseATPartitionName() (string, error) {
	return p.parseATNameString()
}

// parseATPartitionOrKey reads partition_or_key_value (or its subpartition
// form): partition | FOR ( partition_key_value [, partition_key_value ]... ).
// Oracle takes the key values without FOR as well. It sets *name to the
// partition's name, if one is given and name is not nil.
func (p *Parser) parseATPartitionOrKey(name *string) error {
	if p.cur.Type == kwFOR {
		p.advance()
		return p.parseATExprList()
	}
	if p.cur.Type == kwFORCE && p.peekNext().Type == '(' {
		p.advance()
		if err := p.parseATExprList(); err != nil {
			return err
		}
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == '(' {
		return p.parseATExprList()
	}
	if name == nil {
		return p.parseATName()
	}
	var err error
	*name, err = p.parseATNameString()
	return err
}

// atcPartitionAttributes is partition_attributes:
//
//	[ physical_attributes_clause ] [ logging_clause ] [ allocate_extent_clause ]
//	[ deallocate_unused_clause ] [ SHRINK SPACE [ CASCADE ] ]
//	[ { table_compression | prefix_compression } ] [ inmemory_clause ]
//	[ LOB ( lob_item ) ( modify_LOB_parameters ) ]...
//	[ OVERFLOW { physical_attributes_clause | logging_clause
//	  | allocate_extent_clause | deallocate_unused_clause | SHRINK SPACE [ CASCADE ] } ]
const atcPartitionAttributes = atcPhysical | atcAllocate | atcDeallocate | atcShrink |
	atcCompression | atcInmemory | atcLobModify | atcLobStorage | atcClustering

// parseATModifyPartitioning reads the MODIFY operations of
// alter_table_partitioning: modify_table_default_attrs,
// modify_table_partition and modify_table_subpartition.
func (p *Parser) parseATModifyPartitioning(cmd *nodes.AlterTableCmd) error {
	p.advance() // MODIFY
	switch p.cur.Type {
	case kwDEFAULT:
		return p.parseATModifyDefaultAttrs(cmd)
	case kwSUBPARTITION:
		return p.parseATModifySubpartition(cmd)
	}
	return p.parseATModifyPartition(cmd)
}

// parseATModifyDefaultAttrs reads modify_table_default_attrs after MODIFY:
//
//	DEFAULT ATTRIBUTES [ FOR PARTITION partition ]
//	    [ deferred_segment_creation ] [ read_only_clause ] [ indexing_clause ]
//	    [ segment_attributes_clause ] [ table_compression ] [ inmemory_clause ]
//	    [ prefix_compression ] [ OVERFLOW segment_attributes_clause ]
//	    [ { LOB ( lob_item ) | VARRAY varray } ( LOB_parameters ) ]...
func (p *Parser) parseATModifyDefaultAttrs(cmd *nodes.AlterTableCmd) error {
	cmd.Action = nodes.AT_MODIFY_DEFAULT_ATTRS
	p.advance() // DEFAULT
	if err := p.atExpectWord("ATTRIBUTES"); err != nil {
		return err
	}
	if p.cur.Type == kwFOR {
		start := p.pos()
		p.advance()
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		name, err := p.parseATPartitionName()
		if err != nil {
			return err
		}
		cmd.ColumnName = name
		p.atOption(cmd.Options, start)
	}
	const allowed = atcSegmentCreation | atcReadOnly | atcIndexing | atcSegment |
		atcCompression | atcInmemory | atcOverflow
	// Oracle needs an attribute to modify (ORA-14121 for MODIFY DEFAULT
	// ATTRIBUTES alone).
	attrStart := p.pos()
	for {
		if err := p.parseATClauses(allowed, cmd.Options); err != nil {
			return err
		}
		if !(p.atWord("LOB") && p.peekNext().Type == '(') && !p.atWord("VARRAY") {
			if p.pos() == attrStart {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
		start := p.pos()
		if p.atWord("LOB") {
			p.advance()
			if err := p.parseATLobItemList(); err != nil {
				return err
			}
		} else {
			p.advance()
			if err := p.parseATName(); err != nil {
				return err
			}
		}
		if err := p.parseATLobStorageParameters(); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
	}
}

// parseATModifyPartition reads modify_table_partition after MODIFY:
//
//	PARTITION partition
//	    { partition_attributes
//	    | { add_range_subpartition | add_hash_subpartition | add_list_subpartition }
//	    | COALESCE SUBPARTITION [ update_index_clauses ] [ parallel_clause ]
//	        [ allow_disallow_clustering ]
//	    | [ REBUILD ] UNUSABLE LOCAL INDEXES
//	    | read_only_clause
//	    | indexing_clause }
func (p *Parser) parseATModifyPartition(cmd *nodes.AlterTableCmd) error {
	cmd.Action = nodes.AT_MODIFY_PARTITION
	p.advance() // PARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	start := p.pos()
	switch {
	case p.cur.Type == kwADD:
		// add_{range,hash,list}_subpartition:
		//   ADD subpartition_desc [, subpartition_desc ]...
		//       [ dependent_tables_clause ] [ update_index_clauses ] [ parallel_clause ]
		p.advance()
		if p.cur.Type == kwVALUES {
			if p.atNextWord("LESS") {
				// Oracle consumes VALUES before refusing LESS on ADD VALUES LESS THAN.
				p.advance()
				return p.syntaxErrorAtCur()
			}
			// Oracle reads ADD VALUES on a partition (compat.tsv's row).
			if err := p.parseATValuesClause(); err != nil {
				return err
			}
			p.atOption(cmd.Options, start)
			return p.parseATClauses(atcDependentTables|atcIndexMaintenance|atcParallel, cmd.Options)
		}
		var lastSubUpdateTok Token
		subMethod := atMethodAny
		for {
			subStartPos := p.pos()
			beforeValues := p.cur
			if err := p.parseATSubpartitionDesc(subMethod); err != nil {
				return err
			}
			_ = beforeValues
			if subMethod == atMethodAny && p.atLastSubpartitionMethod != atMethodAny {
				subMethod = p.atLastSubpartitionMethod
			}
			if p.prev.Type == kwSUBPARTITION && p.cur.Type == kwUPDATE {
				lastSubUpdateTok = p.cur
			} else if p.pos() > subStartPos {
				lastSubUpdateTok = Token{}
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		p.atOption(cmd.Options, start)
		if err := p.parseATClauses(atcDependentTables|atcIndexMaintenance|atcParallel, cmd.Options); err != nil {
			if lastSubUpdateTok.Type == kwUPDATE {
				return p.syntaxErrorAtTok(lastSubUpdateTok)
			}
			return err
		}
		if lastSubUpdateTok.Type == kwUPDATE && p.cur.Type != ';' && p.cur.Type != tokEOF {
			return p.syntaxErrorAtTok(lastSubUpdateTok)
		}
		return nil
	case p.atWord("COALESCE"):
		// Oracle takes COALESCE alone as well, as it does at the top of
		// ALTER TABLE for an index-organized table.
		p.advance()
		if p.isStatementEnd() {
			return nil
		}
		if err := p.atExpect(kwSUBPARTITION); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
		return p.parseATClauses(atcIndexMaintenance|atcParallel|atcClustering, cmd.Options)
	case p.atWord("REBUILD"), p.atWord("UNUSABLE"):
		if err := p.parseATUnusableLocalIndexes(); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
		return nil
	case p.cur.Type == kwREAD, p.atWord("INDEXING"):
		if err := p.parseATClauses(atcReadOnly|atcIndexing, cmd.Options); err != nil {
			return err
		}
		return nil
	}
	if err := p.parseATClauses(atcPartitionAttributes, cmd.Options); err != nil {
		return err
	}
	if p.atWord("OVERFLOW") {
		start := p.pos()
		p.advance()
		// Oracle takes OVERFLOW alone as well.
		ok, err := p.parseATClause(atcPhysical|atcAllocate|atcDeallocate|atcShrink|atcClustering, nil)
		if err != nil {
			return err
		}
		if !ok && !p.isStatementEnd() {
			return p.syntaxErrorAtCur()
		}
		p.atOption(cmd.Options, start)
	}
	return nil
}

// parseATUnusableLocalIndexes reads [ REBUILD ] UNUSABLE LOCAL INDEXES.
func (p *Parser) parseATUnusableLocalIndexes() error {
	if p.atWord("REBUILD") {
		p.advance()
	}
	if err := p.atExpectWord("UNUSABLE"); err != nil {
		return err
	}
	if err := p.atExpectWord("LOCAL"); err != nil {
		return err
	}
	return p.atExpectWord("INDEXES")
}

// parseATModifySubpartition reads modify_table_subpartition after MODIFY:
//
//	SUBPARTITION subpartition
//	    { allocate_extent_clause | deallocate_unused_clause
//	    | SHRINK SPACE [ CASCADE ]
//	    | { LOB ( lob_item ) | VARRAY varray } ( modify_LOB_parameters )
//	    | [ REBUILD ] UNUSABLE LOCAL INDEXES
//	    | { ADD | DROP } VALUES ( value [, value ]... )
//	    | read_only_clause | indexing_clause }
func (p *Parser) parseATModifySubpartition(cmd *nodes.AlterTableCmd) error {
	defer func(saved bool) { p.alterSubpartition = saved }(p.alterSubpartition)
	p.alterSubpartition = true
	cmd.Action = nodes.AT_MODIFY_SUBPARTITION
	p.advance() // SUBPARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	if p.isStatementEnd() {
		// Oracle takes MODIFY SUBPARTITION s alone, as it takes MODIFY
		// PARTITION p with no partition_attributes.
		return nil
	}
	start := p.pos()
	switch {
	case p.atWord("REBUILD"), p.atWord("UNUSABLE"):
		err = p.parseATUnusableLocalIndexes()
	case p.cur.Type == kwADD, p.cur.Type == kwDROP:
		p.advance()
		if err = p.atExpect(kwVALUES); err == nil {
			err = p.parseATExprList()
		}
	default:
		var ok bool
		ok, err = p.parseATClause(atcAllocate|atcDeallocate|atcShrink|atcLobModify|atcLobStorage|atcReadOnly|atcIndexing|atcClustering, nil)
		if err == nil && !ok {
			err = p.syntaxErrorAtCur()
		}
	}
	if err != nil {
		return err
	}
	p.atOption(cmd.Options, start)
	return nil
}

// parseATSetPartitioning reads alter_automatic_partitioning,
// alter_interval_partitioning and set_subpartition_template:
//
//	SET PARTITIONING { AUTOMATIC | MANUAL }
//	SET INTERVAL ( [ expr ] )
//	SET SUBPARTITION TEMPLATE { ( subpartition_desc [, ...] ) | integer }
//
// Oracle takes an integer, not SUBPARTITIONS integer, as the template's
// hash_subparts_by_quantity, and an empty list.
func (p *Parser) parseATSetPartitioning(cmd *nodes.AlterTableCmd) error {
	p.advance() // SET
	start := p.pos()
	switch {
	case p.atWord("AUTOMATIC"), p.atWord("MANUAL"):
		// Oracle takes SET { AUTOMATIC | MANUAL } for SET PARTITIONING.
		cmd.Action = nodes.AT_SET_PARTITIONING
		cmd.Subtype = p.cur.Str
		p.advance()
		return nil
	case p.atWord("STORE"):
		// SET STORE IN ( tablespace [, tablespace ]... ) of interval partitioning.
		cmd.Action = nodes.AT_SET_INTERVAL
		p.advance()
		if err := p.atExpect(kwIN); err != nil {
			return err
		}
		if err := p.parseATNameList(); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
		return nil
	case p.atWord("PARTITIONING"):
		cmd.Action = nodes.AT_SET_PARTITIONING
		p.advance()
		// Oracle reads no SET PARTITIONING alone: ALTER TABLE a SET
		// PARTITIONING answers ORA-14853, an invalid option for the clause,
		// at byte 0, and takes SET PARTITIONING AUTOMATIC or MANUAL.
		if !p.atWord("AUTOMATIC") && !p.atWord("MANUAL") {
			return p.syntaxErrorAtCur()
		}
		cmd.Subtype = p.cur.Str
		p.advance()
		return nil
	case p.cur.Type == kwINTERVAL:
		cmd.Action = nodes.AT_SET_INTERVAL
		p.advance()
		if err := p.atExpect('('); err != nil {
			return err
		}
		if p.cur.Type != ')' {
			e, err := p.parseExpr()
			if err != nil {
				return err
			}
			if e == nil {
				return p.syntaxErrorAtCur()
			}
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
		return nil
	}
	cmd.Action = nodes.AT_SET_SUBPARTITION_TEMPLATE
	p.advance() // SUBPARTITION
	if err := p.atExpectWord("TEMPLATE"); err != nil {
		return err
	}
	start = p.pos()
	defer func(saved bool) { p.atSubpartitionTemplate = saved }(p.atSubpartitionTemplate)
	p.atSubpartitionTemplate = true
	switch {
	case p.cur.Type == tokICONST:
		p.advance()
	case p.cur.Type == '(' && p.peekNext().Type == ')':
		p.advance()
		p.advance()
	default:
		if err := p.parseATSubpartitionDescList(atMethodAny); err != nil {
			return err
		}
	}
	p.atOption(cmd.Options, start)
	return nil
}

// parseATAddPartition reads add_table_partition:
//
//	ADD PARTITION [ partition ] [ range_values_clause | list_values_clause ]
//	    [ table_partition_description | partitioning_storage_clause ]
//	    [ external_part_subpart_data_props ]
//	    [ ( subpartition_desc [, ...] ) | hash_subparts_by_quantity ]
//	    [ update_index_clauses ] [ parallel_clause ] [ read_only_clause ]
//	    [ indexing_clause ]
//
// add_range_partition_clause, add_list_partition_clause,
// add_system_partition_clause and add_hash_partition_clause differ in
// what the table's partitioning method allows, which Oracle decides after
// the parse.
func (p *Parser) parseATAddPartition(cmd *nodes.AlterTableCmd) error {
	cmd.Action = nodes.AT_ADD_PARTITION
	cmd.Subtype = "PARTITION"
	p.advance() // ADD
	const allowed = atcTablePartitionDescription | atcPartitioningStorage | atcExternal |
		atcIndexMaintenance | atcParallel
	for {
		p.advance() // PARTITION
		// Oracle refuses COMPRESS and NOCOMPRESS straight after PARTITION
		// (ORA-00902); after a name or VALUES they parse.
		if p.cur.Type == kwCOMPRESS || p.cur.Type == kwNOCOMPRESS {
			return p.syntaxErrorAtCur()
		}
		var nameTok Token
		hasName := false
		if p.isATPartitionNameHere(allowed) {
			nameTok = p.cur
			hasName = true
			name, err := p.parseATPartitionName()
			if err != nil {
				return err
			}
			if cmd.ColumnName == "" {
				cmd.ColumnName = name
			}
		}
		afterNamePos := p.pos()
		firstTokAfterPart := p.cur
		if p.cur.Type == kwVALUES {
			start := p.pos()
			if err := p.parseATValuesClause(); err != nil {
				return err
			}
			p.atOption(cmd.Options, start)
		}
		if err := p.parseATClauses(allowed, cmd.Options); err != nil {
			if !hasName && firstTokAfterPart.Type == kwUPDATE {
				return p.syntaxErrorAtTok(firstTokAfterPart)
			}
			return err
		}
		start := p.pos()
		switch {
		case p.cur.Type == '(':
			if !hasName && p.pos() == afterNamePos && p.peekNext().Type != kwSUBPARTITION {
				return p.syntaxErrorAtCur()
			}
			if err := p.parseATSubpartitionDescList(atMethodAny); err != nil {
				return err
			}
		case p.atWord("SUBPARTITIONS"):
			if err := p.parseATHashSubpartsByQuantity(); err != nil {
				return err
			}
		}
		p.atOption(cmd.Options, start)
		if err := p.parseATClauses(atcIndexMaintenance|atcParallel, cmd.Options); err != nil {
			return err
		}
		if p.cur.Type != ',' && p.cur.Type != ';' && p.cur.Type != tokEOF {
			if hasName && p.pos() == afterNamePos {
				return p.syntaxErrorAtTok(nameTok)
			}
			if !hasName && firstTokAfterPart.Type == kwUPDATE {
				return p.syntaxErrorAtTok(firstTokAfterPart)
			}
		}
		if p.cur.Type != ',' || p.peekNext().Type != kwPARTITION {
			return nil
		}
		p.advance() // ,
	}
}

// parseATDropPartition reads drop_table_partition and
// drop_table_subpartition:
//
//	DROP { PARTITION | SUBPARTITION } partition_or_key_value
//	    [, partition_or_key_value ]...
//	    [ update_index_clauses [ parallel_clause ] ]
//
// Oracle takes parallel_clause only after update_index_clauses
// (ORA-14048 for DROP PARTITION p NOPARALLEL).
func (p *Parser) parseATDropPartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // DROP
	cmd.Action = nodes.AT_DROP_PARTITION
	if p.cur.Type == kwSUBPARTITION {
		cmd.Action = nodes.AT_DROP_SUBPARTITION
	}
	p.advance() // PARTITION | SUBPARTITION
	if err := p.parseATPartitionOrKeyList(cmd); err != nil {
		return err
	}
	return p.parseATUpdateIndexesParallel(cmd)
}

// parseATPartitionOrKeyList reads partition_or_key_value [, ...] and keeps
// the first name.
func (p *Parser) parseATPartitionOrKeyList(cmd *nodes.AlterTableCmd) error {
	for {
		var name string
		if err := p.parseATPartitionOrKey(&name); err != nil {
			return err
		}
		if cmd.ColumnName == "" {
			cmd.ColumnName = name
		}
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseATUpdateIndexesParallel reads [ update_index_clauses [ parallel_clause ] ].
func (p *Parser) parseATUpdateIndexesParallel(cmd *nodes.AlterTableCmd) error {
	if !p.isATUpdateIndexStart() {
		return nil
	}
	start := p.pos()
	if err := p.parseATUpdateIndexClauses(); err != nil {
		return err
	}
	p.atOption(cmd.Options, start)
	if p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL {
		start := p.pos()
		if err := p.parseATParallel(); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
	}
	return nil
}

// parseATRenamePartition reads rename_partition_subpart:
// RENAME { PARTITION | SUBPARTITION } name TO new_name.
func (p *Parser) parseATRenamePartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // RENAME
	cmd.Action = nodes.AT_RENAME_PARTITION
	cmd.Subtype = p.cur.Str
	p.advance() // PARTITION | SUBPARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	if err := p.atExpect(kwTO); err != nil {
		return err
	}
	cmd.NewName, err = p.parseATNameString()
	return err
}

// parseATTruncatePartition reads truncate_partition_subpart:
//
//	TRUNCATE { PARTITION | SUBPARTITION } partition_or_key_value [, ...]
//	    [ { DROP [ ALL ] | REUSE } STORAGE ] [ update_index_clauses [ parallel_clause ] ]
//	    [ CASCADE ]
func (p *Parser) parseATTruncatePartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // TRUNCATE
	cmd.Action = nodes.AT_TRUNCATE_PARTITION
	if p.cur.Type != kwPARTITION && p.cur.Type != kwSUBPARTITION && !(!p.alterMview && p.atWord("SUBPARTITIONS")) {
		return p.syntaxErrorAtCur()
	}
	cmd.Subtype = p.cur.Str
	if !p.alterMview && p.atWord("SUBPARTITIONS") {
		cmd.Subtype = "SUBPARTITION"
	}
	p.advance() // PARTITION | SUBPARTITION
	if err := p.parseATPartitionOrKeyList(cmd); err != nil {
		return err
	}
	if p.cur.Type == kwDROP || p.atWord("REUSE") {
		start := p.pos()
		if p.cur.Type == kwDROP {
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
			}
		} else {
			p.advance()
		}
		if err := p.atExpectWord("STORAGE"); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
	}
	if err := p.parseATUpdateIndexesParallel(cmd); err != nil {
		return err
	}
	if p.cur.Type == kwCASCADE {
		start := p.pos()
		p.advance()
		p.atOption(cmd.Options, start)
	}
	return nil
}

// parseATMovePartition reads move_table_partition and
// move_table_subpartition:
//
//	MOVE PARTITION partition [ MAPPING TABLE ] [ table_partition_description ]
//	    [ filter_condition ] [ update_index_clauses ] [ parallel_clause ]
//	    [ allow_disallow_clustering ] [ ONLINE ]
//	MOVE SUBPARTITION subpartition [ indexing_clause ]
//	    [ partitioning_storage_clause ] [ update_index_clauses ]
//	    [ filter_condition ] [ parallel_clause ] [ allow_disallow_clustering ]
//	    [ ONLINE ]
func (p *Parser) parseATMovePartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // MOVE
	cmd.Action = nodes.AT_MOVE
	cmd.Subtype = p.cur.Str
	sub := p.cur.Type == kwSUBPARTITION
	p.advance() // PARTITION | SUBPARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	allowed := atcMappingTable | atcTablePartitionDescription
	if sub {
		allowed = atcIndexing | atcPartitioningStorage
	}
	allowed |= atcFilter | atcIndexMaintenance | atcParallel | atcClustering | atcOnline
	return p.parseATClauses(allowed, cmd.Options)
}

// atcSplitMerge is what follows the partitions of SPLIT and MERGE:
//
//	[ filter_condition ] [ dependent_tables_clause ] [ update_index_clauses ]
//	[ parallel_clause ] [ allow_disallow_clustering ] [ ONLINE ]
//
// Oracle refuses filter_condition after MERGE (ORA-14126); SPLIT takes it.
const atcSplitMerge = atcDependentTables | atcIndexMaintenance | atcParallel |
	atcClustering | atcOnline

// parseATSplitPartition reads split_table_partition and
// split_table_subpartition:
//
//	SPLIT { PARTITION | SUBPARTITION } name
//	    { AT ( value [, value ]... ) [ INTO ( spec , spec ) ]
//	    | VALUES ( value [, value ]... ) [ INTO ( spec , spec ) ]
//	    | INTO ( spec [, spec ]... , spec ) }
//	    [ split_nested_table_part ] [ filter_condition ] [ dependent_tables_clause ]
//	    [ update_index_clauses ] [ parallel_clause ] [ allow_disallow_clustering ]
//	    [ ONLINE ]
func (p *Parser) parseATSplitPartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // SPLIT
	cmd.Action = nodes.AT_SPLIT_PARTITION
	cmd.Subtype = "PARTITION"
	if p.cur.Type != kwPARTITION && p.cur.Type != kwSUBPARTITION {
		return p.syntaxErrorAtCur()
	}
	sub := p.cur.Type == kwSUBPARTITION
	if sub {
		cmd.Action = nodes.AT_SPLIT_SUBPARTITION
		cmd.Subtype = "SUBPARTITION"
	}
	p.advance() // PARTITION | SUBPARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	start := p.pos()
	into := true
	intoOnly := false
	switch {
	case p.atWord("AT"), p.cur.Type == kwVALUES:
		p.advance()
		if err := p.parseATBoundValues(false); err != nil {
			return err
		}
		into = p.cur.Type == kwINTO
	case p.cur.Type != kwINTO:
		return p.syntaxErrorAtCur()
	default:
		intoOnly = true
	}
	if into {
		p.advance() // INTO
		if p.cur.Type != '(' {
			if !sub {
				return p.syntaxErrorAtCur()
			}
			// Oracle reads the subpartition list without its parentheses, and
			// the closing one left over at the end, after AT(values)
			// (compat.tsv's SPLIT SUBPARTITION ... AT(1) INTO SUBPARTITION,
			// SUBPARTITION) rows).
			for {
				var had bool
				values := &had
				if !intoOnly {
					values = nil
				}
				if err = p.parseATSubpartitionSpec(values); err != nil {
					return err
				}
				if p.cur.Type != ',' {
					break
				}
				if intoOnly && !had {
					return p.syntaxErrorAtCur()
				}
				p.advance()
			}
			if p.cur.Type != ')' {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		} else {
			p.advance() // consume (
			// INTO alone lists range or list partition descriptions, each
			// with its values but the last (ORA-00926 for one without).
			for {
				var had bool
				values := &had
				if !intoOnly {
					values = nil
				}
				if sub {
					err = p.parseATSubpartitionSpec(values)
				} else {
					err = p.parseATPartitionSpec(values)
				}
				if err != nil {
					return err
				}
				if p.cur.Type != ',' {
					break
				}
				if intoOnly && !had {
					return p.syntaxErrorAtCur()
				}
				p.advance()
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
	}
	p.atOption(cmd.Options, start)
	if !sub && p.cur.Type == kwNESTED {
		start := p.pos()
		if err := p.parseATSplitNestedTablePart(); err != nil {
			return err
		}
		p.atOption(cmd.Options, start)
	}
	return p.parseATClauses(atcSplitMerge|atcSplitFilter, cmd.Options)
}

// parseATSplitNestedTablePart reads split_nested_table_part:
//
//	NESTED TABLE column INTO ( nested_table_partition_spec ,
//	    nested_table_partition_spec [ split_nested_table_part ] )
//	nested_table_partition_spec: PARTITION partition [ segment_attributes_clause ]
func (p *Parser) parseATSplitNestedTablePart() error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	p.advance() // NESTED
	if err := p.atExpect(kwTABLE); err != nil {
		return err
	}
	colTok := p.cur
	if err := p.parseATName(); err != nil {
		return err
	}
	if err := p.atExpect(kwINTO); err != nil {
		if p.cur.Type == '(' {
			return p.syntaxErrorAtTok(colTok)
		}
		return err
	}
	if err := p.atExpect('('); err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		if i == 1 {
			if err := p.atExpect(','); err != nil {
				return err
			}
		}
		if err := p.atExpect(kwPARTITION); err != nil {
			return err
		}
		if p.isATNameStart() && !p.isATSegmentAttrStart() {
			if err := p.parseATName(); err != nil {
				return err
			}
		}
		if p.isATSegmentAttrStart() {
			if err := p.parseATSegmentAttrs(nil); err != nil {
				return err
			}
		}
	}
	if p.cur.Type == kwNESTED {
		if err := p.parseATSplitNestedTablePart(); err != nil {
			return err
		}
	}
	return p.atExpect(')')
}

// parseATMergePartitions reads merge_table_partitions and
// merge_table_subpartitions:
//
//	MERGE { PARTITIONS | SUBPARTITIONS }
//	    { partition_or_key_value , partition_or_key_value [, ...]
//	    | partition_or_key_value TO partition_or_key_value }
//	    [ INTO { partition_spec | subpartition_spec } ]
//	    [ filter_condition ] [ dependent_tables_clause ] [ update_index_clauses ]
//	    [ parallel_clause ] [ allow_disallow_clustering ] [ ONLINE ]
func (p *Parser) parseATMergePartitions(cmd *nodes.AlterTableCmd) error {
	p.advance() // MERGE
	cmd.Action = nodes.AT_MERGE_PARTITIONS
	cmd.Subtype = "PARTITIONS"
	if !p.atWord("PARTITIONS") && !p.atWord("SUBPARTITIONS") {
		return p.syntaxErrorAtCur()
	}
	sub := p.atWord("SUBPARTITIONS")
	if sub {
		cmd.Action = nodes.AT_MERGE_SUBPARTITIONS
		cmd.Subtype = "SUBPARTITIONS"
	}
	p.advance() // PARTITIONS | SUBPARTITIONS
	start := p.pos()
	if err := p.parseATPartitionOrKey(&cmd.ColumnName); err != nil {
		return err
	}
	var err error
	switch p.cur.Type {
	case kwTO:
		p.advance()
		if err := p.parseATPartitionOrKey(nil); err != nil {
			return err
		}
	case ',':
		for p.cur.Type == ',' {
			p.advance()
			if err := p.parseATPartitionOrKey(nil); err != nil {
				return err
			}
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwINTO {
		p.advance()
		if sub {
			err = p.parseATSubpartitionSpec(nil)
		} else {
			err = p.parseATPartitionSpec(nil)
		}
		if err != nil {
			return err
		}
	}
	p.atOption(cmd.Options, start)
	return p.parseATClauses(atcSplitMerge, cmd.Options)
}

// parseATExchangePartition reads exchange_partition_subpart:
//
//	EXCHANGE { PARTITION partition | SUBPARTITION subpartition }
//	    WITH TABLE [ schema. ] table
//	    [ { INCLUDING | EXCLUDING } INDEXES ] [ { WITH | WITHOUT } VALIDATION ]
//	    [ exceptions_clause ] [ update_index_clauses [ parallel_clause ] ]
//	    [ CASCADE ]
//
// Oracle refuses the page's ONLINE here and a parallel_clause not after
// update_index_clauses (ORA-14094, ORA-14020).
func (p *Parser) parseATExchangePartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // EXCHANGE
	cmd.Action = nodes.AT_EXCHANGE_PARTITION
	if p.cur.Type != kwPARTITION && p.cur.Type != kwSUBPARTITION {
		return p.syntaxErrorAtCur()
	}
	cmd.Subtype = p.cur.Str
	p.advance() // PARTITION | SUBPARTITION
	name, err := p.parseATPartitionName()
	if err != nil {
		return err
	}
	cmd.ColumnName = name
	start := p.pos()
	if err := p.atExpect(kwWITH); err != nil {
		return err
	}
	if err := p.atExpect(kwTABLE); err != nil {
		return err
	}
	if err := p.parseATQualifiedName(); err != nil {
		return err
	}
	p.atOption(cmd.Options, start)
	if err := p.parseATClauses(atcExchange, cmd.Options); err != nil {
		return err
	}
	if err := p.parseATUpdateIndexesParallel(cmd); err != nil {
		return err
	}
	return p.parseATClauses(atcExchange|atcCascade, cmd.Options)
}

// parseATCoalescePartition reads coalesce_table_partition:
//
//	COALESCE PARTITION [ update_index_clauses ] [ parallel_clause ]
//	    [ allow_disallow_clustering ]
//
// coalesce_table_subpartition stands only in modify_table_partition:
// Oracle refuses COALESCE SUBPARTITION here (ORA-14004).
func (p *Parser) parseATCoalescePartition(cmd *nodes.AlterTableCmd) error {
	p.advance() // COALESCE
	p.advance() // PARTITION
	cmd.Action = nodes.AT_COALESCE_PARTITION
	return p.parseATClauses(atcIndexMaintenance|atcParallel|atcClustering, cmd.Options)
}
