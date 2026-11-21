package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads CREATE TABLE's table_partitioning_clauses by the grammar
// of the CREATE TABLE page.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html
// BNF: oracle/parser/bnf/CREATE-TABLE.bnf

// createTablePartitionKind is what a partition or subpartition description
// holds as its bounds: VALUES LESS THAN for range, VALUES for list, none
// for hash.
type createTablePartitionKind int

const (
	createTableRangeBound createTablePartitionKind = iota
	createTableListBound
	createTableHashBound
)

// parsePartitionClause parses CREATE TABLE's table_partitioning_clauses.
//
//	range_partitions:
//	    PARTITION BY RANGE ( column [, column ] ... )
//	    [ INTERVAL ( expr ) ]
//	    [ subpartition_by_range | subpartition_by_list | subpartition_by_hash ]
//	    ( PARTITION [ partition ] range_values_clause [ table_partition_description ]
//	      [ subpartitions ] [, ...] )
//	list_partitions:
//	    PARTITION BY LIST ( column [, column ] ... )
//	    [ AUTOMATIC [ STORE IN ( tablespace [, tablespace ] ... ) ] ]
//	    [ subpartition_by_... ]
//	    ( PARTITION [ partition ] list_values_clause [ table_partition_description ]
//	      [ subpartitions ] [, ...] )
//	hash_partitions:
//	    PARTITION BY HASH ( column [, column ] ... )
//	    [ subpartition_by_... ]
//	    { individual_hash_partitions | hash_partitions_by_quantity }
//	reference_partitioning:
//	    PARTITION BY REFERENCE ( constraint )
//	    [ ( PARTITION [ partition ] [ table_partition_description ] [, ...] ) ]
//	system_partitioning:
//	    PARTITION BY SYSTEM
//	    [ PARTITIONS quantity | ( PARTITION [ partition ] [ table_partition_description ] [, ...] ) ]
//
// Oracle 23ai refuses STORE IN after a range table's INTERVAL (ORA-00922 at
// STORE), and a column named twice in a partitioning key, at the second.
// On a table without a column list, it refuses a range or list partition
// without a name, at the token after its values clause (it takes one on a
// table with columns), and reference partitioning, whose constraint is a
// foreign key of the column list.
func (p *Parser) parsePartitionClause(stmt *nodes.CreateTableStmt, hasColumns bool) (*nodes.PartitionClause, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	defer func(saved bool) { p.createTableIndexOrganized = saved }(p.createTableIndexOrganized)
	p.createTableIndexOrganized = stmt.Organization == "INDEX"
	start := p.pos()
	p.advance() // consume PARTITION

	clause := &nodes.PartitionClause{
		Columns:    &nodes.List{},
		Partitions: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}
	if err := p.expectCreateTableToken(kwBY); err != nil {
		return nil, err
	}
	switch {
	case p.isKeywordStr("RANGE"):
		clause.Type = nodes.PARTITION_RANGE
	case p.isKeywordStr("LIST"):
		clause.Type = nodes.PARTITION_LIST
	case p.isKeywordStr("HASH"):
		clause.Type = nodes.PARTITION_HASH
	case p.cur.Type == kwREFERENCE:
		// The constraint is looked up against the table before Oracle asks
		// for a missing column list (ORA-14655), so a table with none is the
		// engine's to check, as enable_disable_clause's is.
		p.createTableEngineCheck = true
		return p.parseCreateTableReferencePartitioning(clause, hasColumns)
	case p.isKeywordStr("SYSTEM"):
		return p.parseCreateTableSystemPartitioning(clause)
	case p.isKeywordStr("CONSISTENT"):
		// consistent_hash_partitions and consistent_hash_with_subpartitions
		// belong to sharded tables: Oracle reads CONSISTENT HASH and refuses
		// the form elsewhere (ORA-14151) at what follows HASH, before it
		// reads the columns.
		p.advance() // CONSISTENT
		if err := p.expectCreateTableToken(kwHASH); err != nil {
			return nil, err
		}
		return nil, p.syntaxErrorAtCur()
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume RANGE / LIST / HASH
	if err := p.parseCreateTablePartitionColumns(clause.Columns); err != nil {
		return nil, err
	}

	switch clause.Type {
	case nodes.PARTITION_RANGE:
		if p.cur.Type == kwINTERVAL {
			p.advance()
			if err := p.expectCreateTableToken('('); err != nil {
				return nil, err
			}
			interval, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if interval == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return nil, err
			}
			clause.Interval = interval
			// An interval takes STORE IN, which a range without one does
			// not (ORA-00922 at STORE), and a log's does not either
			// (ORA-12044 at STORE).
			if p.isKeywordStr("STORE") && !p.createTablePartitionsOfLog {
				if err := p.parseCreateTableStoreIn(); err != nil {
					return nil, err
				}
			}
		}
	case nodes.PARTITION_LIST:
		if p.cur.Type == kwAUTOMATIC {
			p.advance()
			if p.isKeywordStr("STORE") && !p.createTablePartitionsOfLog {
				if err := p.parseCreateTableStoreIn(); err != nil {
					return nil, err
				}
			}
		}
	}

	sub := createTableHashBound
	if p.isKeywordStr("SUBPARTITION") {
		var err error
		clause.Subpartition, sub, err = p.parseCreateTableSubpartitionBy()
		if err != nil {
			return nil, err
		}
	}

	if clause.Type == nodes.PARTITION_HASH && p.isKeywordStr("PARTITIONS") {
		if err := p.parseCreateTableHashPartitionsByQuantity(stmt.Organization == "INDEX"); err != nil {
			return nil, err
		}
		clause.Loc.End = p.prev.End
		return clause, nil
	}

	bound := createTableHashBound
	switch clause.Type {
	case nodes.PARTITION_RANGE:
		bound = createTableRangeBound
	case nodes.PARTITION_LIST:
		bound = createTableListBound
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	for {
		def, err := p.parseCreateTablePartitionDesc(bound, clause.Subpartition != nil, sub)
		if err != nil {
			return nil, err
		}
		clause.Partitions.Items = append(clause.Partitions.Items, def)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseCreateTablePartitionColumns reads ( column [, column ] ... ) into
// columns.
func (p *Parser) parseCreateTablePartitionColumns(columns *nodes.List) error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	seen := map[string]bool{}
	for {
		if !p.isCreateTableNameStart() {
			return p.syntaxErrorAtCur()
		}
		if seen[p.cur.Str] {
			return p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		col, err := p.parseExpr()
		if err != nil {
			return err
		}
		if col == nil {
			return p.syntaxErrorAtCur()
		}
		columns.Items = append(columns.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// parseCreateTableStoreIn reads STORE IN ( tablespace [, tablespace ] ... ).
func (p *Parser) parseCreateTableStoreIn() error {
	p.advance() // consume STORE
	if err := p.expectCreateTableToken(kwIN); err != nil {
		return err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		// Oracle takes a trailing comma before the closing parenthesis
		// (STORE IN (a,)).
		if p.cur.Type == ')' {
			break
		}
	}
	return p.expectCreateTableToken(')')
}

// parseCreateTableHashPartitionsByQuantity reads hash_partitions_by_quantity.
//
//	PARTITIONS hash_partition_quantity
//	[ STORE IN ( tablespace [, tablespace ] ... ) ]
//	[ table_compression | index_compression ]
//	[ OVERFLOW STORE IN ( tablespace [, tablespace ] ... ) ]
//
// OVERFLOW is an index-organized table's: Oracle 23ai refuses it on a heap
// table (ORA-00922 at OVERFLOW).
func (p *Parser) parseCreateTableHashPartitionsByQuantity(indexOrganized bool) error {
	p.advance() // consume PARTITIONS
	if err := p.parseCreateTableInteger(); err != nil {
		return err
	}
	if p.isKeywordStr("STORE") {
		if err := p.parseCreateTableStoreIn(); err != nil {
			return err
		}
	}
	if p.isCreateTablePartitionCompressionStart() && !p.createTablePartitionsOfLog {
		if err := p.parseCreateTablePartitionCompression(); err != nil {
			return err
		}
	}
	if p.isKeywordStr("OVERFLOW") {
		if !indexOrganized {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume OVERFLOW
		if !p.isKeywordStr("STORE") {
			return p.syntaxErrorAtCur()
		}
		return p.parseCreateTableStoreIn()
	}
	return nil
}

// parseCreateTableSubpartitionBy reads subpartition_by_range,
// subpartition_by_list or subpartition_by_hash, and returns the kind of
// bounds its subpartitions hold.
//
//	SUBPARTITION BY { RANGE | LIST } ( column [, column ] ... ) [ subpartition_template ]
//	SUBPARTITION BY HASH ( column [, column ] ... )
//	[ SUBPARTITIONS quantity [ STORE IN ( tablespace [, tablespace ] ... ) ]
//	| subpartition_template
//	]
//	subpartition_template:
//	    SUBPARTITION TEMPLATE ( subpartition_desc [, subpartition_desc ] ... )
//	  | SUBPARTITION TEMPLATE hash_subpartition_quantity
func (p *Parser) parseCreateTableSubpartitionBy() (*nodes.PartitionClause, createTablePartitionKind, error) {
	sub := &nodes.PartitionClause{
		Columns:    &nodes.List{},
		Partitions: &nodes.List{},
		Loc:        nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume SUBPARTITION
	if err := p.expectCreateTableToken(kwBY); err != nil {
		return nil, 0, err
	}
	var bound createTablePartitionKind
	switch {
	case p.isKeywordStr("RANGE"):
		sub.Type, bound = nodes.PARTITION_RANGE, createTableRangeBound
	case p.isKeywordStr("LIST"):
		sub.Type, bound = nodes.PARTITION_LIST, createTableListBound
	case p.isKeywordStr("HASH"):
		sub.Type, bound = nodes.PARTITION_HASH, createTableHashBound
	default:
		return nil, 0, p.syntaxErrorAtCur()
	}
	p.advance() // consume RANGE / LIST / HASH
	if err := p.parseCreateTablePartitionColumns(sub.Columns); err != nil {
		return nil, 0, err
	}
	switch {
	case bound == createTableHashBound && p.isKeywordStr("SUBPARTITIONS"):
		if err := p.parseCreateTableSubpartitionsByQuantity(); err != nil {
			return nil, 0, err
		}
	case p.isKeywordStr("SUBPARTITION"):
		p.advance() // consume SUBPARTITION
		if err := p.expectCreateTableWord("TEMPLATE"); err != nil {
			return nil, 0, err
		}
		if p.cur.Type == tokICONST {
			p.advance()
			if bound != createTableHashBound {
				return nil, 0, p.syntaxErrorAtCur()
			}
			break
		}
		if err := p.parseCreateTableSubpartitionDescs(bound, sub.Partitions); err != nil {
			return nil, 0, err
		}
	}
	sub.Loc.End = p.prev.End
	return sub, bound, nil
}

// parseCreateTableSubpartitionsByQuantity reads hash_subparts_by_quantity.
//
//	SUBPARTITIONS quantity [ STORE IN ( tablespace [, tablespace ] ... ) ]
func (p *Parser) parseCreateTableSubpartitionsByQuantity() error {
	p.advance() // consume SUBPARTITIONS
	if err := p.parseCreateTableInteger(); err != nil {
		return err
	}
	if p.isKeywordStr("STORE") {
		return p.parseCreateTableStoreIn()
	}
	return nil
}

// parseCreateTableSubpartitionDescs reads ( subpartition_desc [, ...] ),
// each a range_subpartition_desc, list_subpartition_desc or
// individual_hash_subparts by bound.
//
//	SUBPARTITION [ subpartition ] [ range_values_clause | list_values_clause ]
//	[ read_only_clause ] [ indexing_clause ] [ partitioning_storage_clause ]
func (p *Parser) parseCreateTableSubpartitionDescs(bound createTablePartitionKind, list *nodes.List) error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		start := p.pos()
		if !p.isKeywordStr("SUBPARTITION") {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume SUBPARTITION
		def := &nodes.PartitionDef{Values: &nodes.List{}, Loc: nodes.Loc{Start: start}}
		if p.isCreateTablePartitionName(bound) {
			def.Name = p.cur.Str
			p.advance()
		}
		if err := p.parseCreateTablePartitionBound(bound, def); err != nil {
			return err
		}
		if err := p.parseCreateTableHashPartitionAttributes(def); err != nil {
			return err
		}
		def.Loc.End = p.prev.End
		list.Items = append(list.Items, def)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// parseCreateTablePartitionDesc reads one partition of a range, list or
// hash partitioned table.
//
//	PARTITION [ partition ] { range_values_clause | list_values_clause }
//	[ table_partition_description ]
//	[ ( subpartition_desc [, ...] ) | SUBPARTITIONS quantity [ STORE IN ( ... ) ] ]
//
//	PARTITION [ partition ] [ read_only_clause ] [ indexing_clause ]
//	[ partitioning_storage_clause ]                      -- hash
func (p *Parser) parseCreateTablePartitionDesc(bound createTablePartitionKind, composite bool, sub createTablePartitionKind) (*nodes.PartitionDef, error) {
	start := p.pos()
	if p.cur.Type != kwPARTITION {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume PARTITION
	def := &nodes.PartitionDef{Values: &nodes.List{}, Loc: nodes.Loc{Start: start}}
	if p.isCreateTablePartitionName(bound) {
		def.Name = p.cur.Str
		p.advance()
	}
	if bound == createTableHashBound {
		if err := p.parseCreateTableHashPartitionAttributes(def); err != nil {
			return nil, err
		}
	} else {
		if err := p.parseCreateTablePartitionBound(bound, def); err != nil {
			return nil, err
		}
		if err := p.parseCreateTablePartitionDescription(def); err != nil {
			return nil, err
		}
	}
	if composite {
		switch {
		case p.cur.Type == '(':
			if err := p.parseCreateTableSubpartitionDescs(sub, &nodes.List{}); err != nil {
				return nil, err
			}
		case sub == createTableHashBound && p.isKeywordStr("SUBPARTITIONS"):
			if err := p.parseCreateTableSubpartitionsByQuantity(); err != nil {
				return nil, err
			}
		}
	}
	def.Loc.End = p.prev.End
	return def, nil
}

// isCreateTablePartitionName reports whether the word after PARTITION or
// SUBPARTITION names it. A range or list partition's bounds start with
// VALUES; a hash partition's description starts with one of its
// attributes, which a word alone before the end of the description names
// instead (PARTITION TABLESPACE, PARTITION READ).
func (p *Parser) isCreateTablePartitionName(bound createTablePartitionKind) bool {
	if !p.isCreateTableNameStart() {
		return false
	}
	if bound != createTableHashBound || p.cur.Type == tokQIDENT {
		return true
	}
	next := p.peekNext()
	if next.Type == ',' || next.Type == ')' {
		return true
	}
	return !p.isCreateTableHashPartitionAttributeStart()
}

// parseCreateTablePartitionBound reads a partition's bounds.
//
//	range_values_clause: VALUES LESS THAN ( { literal | MAXVALUE } [, ...] )
//	list_values_clause: VALUES ( { list_values | DEFAULT } )
//	list_values: { literal | NULL } [, { literal | NULL } ] ...
func (p *Parser) parseCreateTablePartitionBound(bound createTablePartitionKind, def *nodes.PartitionDef) error {
	if bound == createTableHashBound {
		return nil
	}
	if p.cur.Type != kwVALUES {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume VALUES
	if bound == createTableRangeBound {
		if err := p.expectCreateTableWord("LESS"); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("THAN"); err != nil {
			return err
		}
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		val, err := p.parseExpr()
		if err != nil {
			return err
		}
		if val == nil {
			return p.syntaxErrorAtCur()
		}
		def.Values.Items = append(def.Values.Items, val)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// isCreateTablePartitionCompressionStart reports whether a table_compression
// or index_compression starts here.
func (p *Parser) isCreateTablePartitionCompressionStart() bool {
	switch {
	case p.cur.Type == kwCOMPRESS, p.cur.Type == kwNOCOMPRESS:
		return true
	case p.cur.Type == kwROW, p.isKeywordStr("COLUMN"):
		return isCreateTableWordTok(p.peekNext(), "STORE")
	}
	return false
}

// parseCreateTablePartitionCompression reads a table_compression or, for
// an index-organized table, an index_compression.
//
//	table_compression:
//	    { NOCOMPRESS | ROW STORE COMPRESS [ BASIC | ADVANCED ]
//	    | COLUMN STORE COMPRESS FOR { QUERY | ARCHIVE } [ LOW | HIGH ] }
//	index_compression:
//	    { COMPRESS [ integer ] | COMPRESS ADVANCED [ LOW | HIGH ] | NOCOMPRESS }
//
// Oracle 23ai also takes the older COMPRESS [ BASIC | FOR { OLTP | QUERY
// [ LOW | HIGH ] | ARCHIVE [ LOW | HIGH ] | ALL OPERATIONS | DIRECT_LOAD
// OPERATIONS } ], and refuses an index_compression on a heap table
// (ORA-00922 for the table, ORA-14020 for a partition).
func (p *Parser) parseCreateTablePartitionCompression() error {
	switch {
	case p.cur.Type == kwNOCOMPRESS:
		p.advance()
	case p.cur.Type == kwCOMPRESS:
		p.advance()
		switch {
		case p.createTableIndexOrganized && p.cur.Type == tokICONST:
			p.advance()
		case p.createTableIndexOrganized && p.isKeywordStr("ADVANCED"):
			p.advance()
			if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
				p.advance()
			}
		case p.isKeywordStr("BASIC"):
			p.advance()
		case p.cur.Type == kwFOR:
			p.advance() // consume FOR
			switch {
			case p.isKeywordStr("OLTP"):
				p.advance()
			case p.isKeywordStr("QUERY"), p.isKeywordStr("ARCHIVE"):
				p.advance()
				if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
					p.advance()
				}
			case p.cur.Type == kwALL, p.isKeywordStr("DIRECT_LOAD"):
				p.advance()
				if err := p.expectCreateTableWord("OPERATIONS"); err != nil {
					return err
				}
			default:
				return p.syntaxErrorAtCur()
			}
		case p.cur.Type == tokICONST, p.isKeywordStr("ADVANCED"):
			return p.syntaxErrorAtTok(p.prev)
		}
	default:
		return p.parseCreateTableIlmCompression()
	}
	return nil
}

// isCreateTableHashPartitionAttributeStart reports whether an attribute of
// a hash partition or subpartition starts here.
func (p *Parser) isCreateTableHashPartitionAttributeStart() bool {
	switch {
	case p.cur.Type == kwREAD, p.cur.Type == kwTABLESPACE, p.isKeywordStr("INDEXING"),
		p.isKeywordStr("OVERFLOW"), p.isCreateTablePartitionCompressionStart(),
		p.isCreateTableInmemoryStart(), p.isKeywordStr("ILM"), p.isKeywordStr("VARRAY"):
		return true
	}
	return p.isCreateTableLobStorageStart()
}

// parseCreateTableHashPartitionAttributes reads the attributes of a hash
// partition, or of a subpartition.
//
//	[ read_only_clause ] [ indexing_clause ] [ partitioning_storage_clause ]
//
//	partitioning_storage_clause:
//	    [ TABLESPACE tablespace ]
//	    [ OVERFLOW [ TABLESPACE tablespace ] ]
//	    [ table_compression ] [ index_compression ]
//	    [ inmemory_clause ] [ ilm_clause ]
//	    [ LOB_partitioning_storage [, LOB_partitioning_storage ] ... ]
//	    [ VARRAY varray_item STORE AS [ SECUREFILE | BASICFILE ] LOB LOB_segname ]
//
// Oracle 23ai refuses an ilm_clause here (ORA-14176), and OVERFLOW
// followed by anything but TABLESPACE (ORA-14020).
func (p *Parser) parseCreateTableHashPartitionAttributes(def *nodes.PartitionDef) error {
	for {
		if p.createTablePartitionsOfLog && p.cur.Type != kwTABLESPACE && !p.isKeywordStr("OVERFLOW") {
			// A log's hash partition or subpartition takes TABLESPACE and
			// OVERFLOW only.
			if p.cur.Type != ',' && p.cur.Type != ')' {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
		switch {
		case p.cur.Type == kwREAD:
			p.advance() // consume READ
			if !p.isKeywordStr("ONLY") && p.cur.Type != kwWRITE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("INDEXING"):
			p.advance() // consume INDEXING
			if p.cur.Type != kwON && !p.isKeywordStr("OFF") {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			def.Tablespace = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		case p.isKeywordStr("OVERFLOW"):
			overflow := p.cur
			p.advance() // consume OVERFLOW
			if p.cur.Type == kwTABLESPACE {
				p.advance()
				if err := p.parseCreateTableName(); err != nil {
					return err
				}
			} else if p.cur.Type != ',' && p.cur.Type != ')' {
				return p.syntaxErrorAtTok(overflow)
			}
		case p.isCreateTablePartitionCompressionStart():
			if err := p.parseCreateTablePartitionCompression(); err != nil {
				return err
			}
		case p.isCreateTableInmemoryStart():
			if err := p.parseCreateTablePartitionInmemory(); err != nil {
				return err
			}
		case p.isKeywordStr("ILM"):
			return p.syntaxErrorAtCur()
		case p.isCreateTableLobStorageStart():
			if err := p.parseCreateTableLobPartitioningStorage(); err != nil {
				return err
			}
		case p.isKeywordStr("VARRAY"):
			if err := p.parseCreateTablePartitionVarray(); err != nil {
				return err
			}
		default:
			if p.cur.Type != ',' && p.cur.Type != ')' {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseCreateTablePartitionInmemory reads a partition's inmemory_clause.
//
//	{ INMEMORY [ inmemory_attributes ] | NO INMEMORY }
func (p *Parser) parseCreateTablePartitionInmemory() error {
	if p.isKeywordStr("NO") {
		p.advance() // consume NO
		p.advance() // consume INMEMORY
		return nil
	}
	p.advance() // consume INMEMORY
	return p.parseCreateTableInmemoryAttributes(false)
}

// parseCreateTableLobPartitioningStorage reads a LOB_partitioning_storage.
//
//	LOB ( LOB_item ) STORE AS [ SECUREFILE | BASICFILE ]
//	[ LOB_segname [ ( TABLESPACE tablespace | storage_clause ) [ ... ] ]
//	| ( TABLESPACE tablespace | storage_clause ) [ ... ]
//	]
//
// Oracle 23ai requires a segment name, a file kind or the parenthesized
// storage after STORE AS (ORA-00906), and takes only TABLESPACE and a
// storage_clause in the parentheses (ORA-22877).
func (p *Parser) parseCreateTableLobPartitioningStorage() error {
	p.advance() // consume LOB
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	if err := p.parseCreateTableItem(); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("STORE"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	file := false
	if p.isKeywordStr("SECUREFILE") || p.isKeywordStr("BASICFILE") {
		p.advance()
		file = true
	}
	if p.cur.Type != '(' {
		if !p.isCreateTableNameStart() {
			if file {
				return nil
			}
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume LOB_segname
		if p.cur.Type != '(' {
			return nil
		}
	}
	p.advance() // consume '('
	for first := true; first || p.cur.Type != ')'; first = false {
		switch {
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return err
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	p.advance() // consume ')'
	return nil
}

// parseCreateTablePartitionVarray reads the VARRAY of a
// partitioning_storage_clause.
//
//	VARRAY varray_item STORE AS [ SECUREFILE | BASICFILE ] LOB LOB_segname
func (p *Parser) parseCreateTablePartitionVarray() error {
	p.advance() // consume VARRAY
	if err := p.parseCreateTableItem(); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("STORE"); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	file := false
	if p.isKeywordStr("SECUREFILE") || p.isKeywordStr("BASICFILE") {
		p.advance()
		file = true
	}
	if err := p.expectCreateTableWord("LOB"); err != nil {
		return err
	}
	if file && p.isCreateTableRefusedFileSegname() {
		return p.syntaxErrorAtCur()
	}
	return p.parseCreateTableName()
}

// parseCreateTablePartitionDescription reads a range or list partition's
// table_partition_description.
//
//	[ deferred_segment_creation ] [ read_only_clause ] [ indexing_clause ]
//	[ segment_attributes_clause ] [ table_compression | index_compression ]
//	[ inmemory_clause ] [ ilm_clause ] [ OVERFLOW [ segment_attributes_clause ] ]
//	[ { LOB_storage_clause | varray_col_properties | nested_table_col_properties } ... ]
func (p *Parser) parseCreateTablePartitionDescription(def *nodes.PartitionDef) error {
	for {
		if p.createTablePartitionsOfLog && !p.isCreateTableLogPartitionAttributeStart() {
			return nil
		}
		switch {
		case p.isKeywordStr("SEGMENT"):
			p.advance() // consume SEGMENT
			if err := p.expectCreateTableWord("CREATION"); err != nil {
				return err
			}
			if p.cur.Type != kwIMMEDIATE && p.cur.Type != kwDEFERRED {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.cur.Type == kwREAD:
			p.advance() // consume READ
			if !p.isKeywordStr("ONLY") && p.cur.Type != kwWRITE {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("INDEXING"):
			p.advance() // consume INDEXING
			if p.cur.Type != kwON && !p.isKeywordStr("OFF") {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			def.Tablespace = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		case p.cur.Type == kwPCTFREE, p.isKeywordStr("PCTUSED"), p.isKeywordStr("INITRANS"), p.isKeywordStr("MAXTRANS"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return err
			}
		case p.isCreateTableLoggingStart():
			p.advance()
		case p.isCreateTablePartitionCompressionStart():
			if err := p.parseCreateTablePartitionCompression(); err != nil {
				return err
			}
		case p.isCreateTableInmemoryStart():
			if err := p.parseCreateTablePartitionInmemory(); err != nil {
				return err
			}
		case p.isKeywordStr("ILM"):
			if err := p.parseCreateTableIlmClause(); err != nil {
				return err
			}
		case p.isKeywordStr("OVERFLOW"):
			p.advance()
		case p.isCreateTableLobStorageStart():
			if err := p.parseCreateTableLobStorageClause(); err != nil {
				return err
			}
		case p.isKeywordStr("VARRAY"):
			if err := p.parseCreateTableVarrayColProperties(true); err != nil {
				return err
			}
		case p.isKeywordStr("NESTED"):
			if err := p.parseCreateTableNestedTableColProperties(true); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// isCreateTableLogPartitionAttributeStart reports a segment attribute, the
// only description a materialized view log's range or list partition takes:
// TABLESPACE, PCTFREE, PCTUSED, INITRANS, MAXTRANS, STORAGE or
// logging_clause.
func (p *Parser) isCreateTableLogPartitionAttributeStart() bool {
	return p.cur.Type == kwTABLESPACE || p.cur.Type == kwPCTFREE || p.isKeywordStr("PCTUSED") ||
		p.isKeywordStr("INITRANS") || p.isKeywordStr("MAXTRANS") || p.isKeywordStr("STORAGE") ||
		p.isCreateTableLoggingStart()
}

// parseCreateTableReferencePartitioning reads reference_partitioning after
// PARTITION BY.
//
//	REFERENCE ( constraint )
//	[ ( PARTITION [ partition ] [ table_partition_description ] [, ...] ) ]
func (p *Parser) parseCreateTableReferencePartitioning(clause *nodes.PartitionClause, hasColumns bool) (*nodes.PartitionClause, error) {
	p.advance() // consume REFERENCE
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	// The constraint is a foreign key of the table's column list, which
	// Oracle looks up after the parse: PARTITION BY REFERENCE (a) on a
	// table without a column list answers ORA-14655 at the name, as it does
	// on one whose constraint is not there. Anything after the name is
	// refused at the name there (ORA-00922, the byte the row of the bare
	// form reads as the missing constraint's).
	nameTok := p.cur
	if err := p.parseCreateTableName(); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	if !hasColumns {
		if !p.isStatementEnd() {
			return nil, p.syntaxErrorAtTok(nameTok)
		}
		clause.Loc.End = p.prev.End
		return clause, nil
	}
	if p.cur.Type == '(' {
		if err := p.parseCreateTableDescribedPartitions(clause); err != nil {
			return nil, err
		}
	}
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseCreateTableSystemPartitioning reads system_partitioning after
// PARTITION BY.
//
//	SYSTEM
//	[ PARTITIONS quantity
//	| ( PARTITION [ partition ] [ table_partition_description ] [, ...] )
//	]
func (p *Parser) parseCreateTableSystemPartitioning(clause *nodes.PartitionClause) (*nodes.PartitionClause, error) {
	p.advance() // consume SYSTEM
	switch {
	case p.isKeywordStr("PARTITIONS"):
		p.advance()
		if err := p.parseCreateTableInteger(); err != nil {
			return nil, err
		}
	case p.cur.Type == '(':
		if err := p.parseCreateTableDescribedPartitions(clause); err != nil {
			return nil, err
		}
	}
	clause.Loc.End = p.prev.End
	return clause, nil
}

// parseCreateTableDescribedPartitions reads
// ( PARTITION [ partition ] [ table_partition_description ] [, ...] ).
func (p *Parser) parseCreateTableDescribedPartitions(clause *nodes.PartitionClause) error {
	p.advance() // consume '('
	for {
		start := p.pos()
		if p.cur.Type != kwPARTITION {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume PARTITION
		def := &nodes.PartitionDef{Values: &nodes.List{}, Loc: nodes.Loc{Start: start}}
		if p.isCreateTableNameStart() && !p.isCreateTableHashPartitionAttributeStart() {
			def.Name = p.cur.Str
			p.advance()
		} else if (p.cur.Type == kwPARTITION || p.cur.Type == kwREAD || p.atWord("INNER") ||
			p.atWord("ONLY") || p.atWord("NOCOPY")) && (p.peekNext().Type == ')' || p.peekNext().Type == ',') {
			// Oracle reads one of these words as the partition's name when
			// nothing follows it (compat.tsv's PARTITION INNER and kin rows).
			def.Name = p.cur.Str
			p.advance()
		}
		if err := p.parseCreateTablePartitionDescription(def); err != nil {
			return err
		}
		def.Loc.End = p.prev.End
		clause.Partitions.Items = append(clause.Partitions.Items, def)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}
