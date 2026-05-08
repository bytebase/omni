package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAlterTableStmt parses an ALTER TABLE statement after ALTER has been consumed.
// The caller has already consumed ALTER; this function consumes TABLE and the rest.
//
// BNF: oracle/parser/bnf/ALTER-TABLE.bnf
//
//	ALTER TABLE [ schema. ] table
//	    { alter_table_properties
//	    | column_clauses
//	    | constraint_clauses
//	    | alter_table_partitioning
//	    | alter_table_partitionset
//	    | alter_external_table
//	    | move_table_clause
//	    | modify_to_partitioned
//	    | modify_opaque_type
//	    | immutable_table_clauses
//	    | blockchain_table_clauses
//	    | duplicated_table_refresh
//	    | enable_disable_clause
//	    }
//	    [ enable_disable_clause ]...
//	    [ { ENABLE | DISABLE } TABLE LOCK ]
//	    [ { ENABLE | DISABLE } ALL TRIGGERS ]
func (p *Parser) parseAlterTableStmt(start int) (*nodes.AlterTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.AlterTableStmt{
		Actions: &nodes.List{},
		Loc:     nodes.Loc{Start: start},
	}
	var parseErr132 error

	// Table name
	stmt.Name, parseErr132 = p.parseObjectName()
	if parseErr132 !=

		// Parse one or more ALTER TABLE actions
		nil {
		return nil, parseErr132
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		cmd, parseErr133 := p.parseAlterTableAction()
		if parseErr133 != nil {
			return nil, parseErr133
		}
		if cmd == nil {
			break
		}
		stmt.Actions.Items = append(stmt.Actions.Items, cmd)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterTableAction parses a single ALTER TABLE action.
//
//	alter_table_properties | column_clauses | constraint_clauses |
//	alter_table_partitioning | move_table_clause | modify_to_partitioned |
//	modify_opaque_type | immutable_table_clauses | blockchain_table_clauses |
//	duplicated_table_refresh | enable_disable_clause
func (p *Parser) parseAlterTableAction() (*nodes.AlterTableCmd, error) {
	switch p.cur.Type {
	case kwADD:
		return p.parseAlterTableAdd()
	case kwMODIFY:
		return p.parseAlterTableModify()
	case kwDROP:
		return p.parseAlterTableDrop()
	case kwRENAME:
		return p.parseAlterTableRename()
	case kwTRUNCATE:
		return p.parseAlterTableTruncatePartition()
	case kwENABLE, kwDISABLE:
		return p.parseAlterTableEnableDisable()
	case kwSET:
		return p.parseAlterTableSet()
	case kwLOGGING, kwNOLOGGING:
		return p.parseAlterTableProperty()
	case kwCOMPRESS, kwNOCOMPRESS:
		return p.parseAlterTableProperty()
	case kwPARALLEL, kwNOPARALLEL:
		return p.parseAlterTableProperty()
	case kwCACHE:
		return p.parseAlterTableProperty()
	case kwREFRESH:
		return p.parseAlterTableRefresh()
	case kwFLASHBACK:
		return p.parseAlterTableFlashbackArchive()
	case kwREAD:
		return p.parseAlterTableReadOnly()
	case kwROW:
		return p.parseAlterTableRowStore()
	case kwRESULT_CACHE:
		return p.parseAlterTableProperty()
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "MOVE":
				return p.parseAlterTableMove()
			case "SPLIT":
				return p.parseAlterTableSplit()
			case "MERGE":
				return p.parseAlterTableMerge()
			case "EXCHANGE":
				return p.parseAlterTableExchange()
			case "COALESCE":
				return p.parseAlterTableCoalesce()
			case "SHRINK":
				return p.parseAlterTableShrinkSpace()
			case "INMEMORY":
				return p.parseAlterTableInmemory()
			case "ILM":
				return p.parseAlterTableILM()
			case "SUPPLEMENTAL":
				return p.parseAlterTableSupplementalLog()
			case "ALLOCATE":
				return p.parseAlterTableAllocateExtent()
			case "DEALLOCATE":
				return p.parseAlterTableDeallocateUnused()
			case "NOCOMPRESS":
				return p.parseAlterTableProperty()
			case "NOCACHE":
				return p.parseAlterTableProperty()
			case "PCTFREE", "PCTUSED", "INITRANS":
				return p.parseAlterTablePhysicalAttr()
			case "STORAGE":
				return p.parseAlterTableStorageClause()
			case "MEMOPTIMIZE":
				return p.parseAlterTableProperty()
			case "COLUMN":
				// COLUMN in MODIFY COLUMN ... SUBSTITUTABLE
				return nil, nil
			case "OVERFLOW":
				return p.parseAlterTableProperty()
			case "MAPPING":
				return p.parseAlterTableProperty()
			case "NOMAPPING":
				return p.parseAlterTableProperty()
			case "NO":
				return p.parseAlterTableNoPrefix()
			case "RECORDS_PER_BLOCK":
				return p.parseAlterTableProperty()
			case "UPGRADE":
				return p.parseAlterTableProperty()
			case "INDEXING":
				return p.parseAlterTableProperty()
			case "SEGMENT":
				return p.parseAlterTableProperty()
			case "ANNOTATIONS":
				return p.parseAlterTableProperty()
			case "PCTTHRESHOLD":
				return p.parseAlterTablePhysicalAttr()
			}
		}
		return nil, nil
	}
}

// parseAlterTableAdd parses ADD column, constraint, partition, or period.
//
//	add_column_clause:
//	    ADD ( { column_definition | virtual_column_definition } [, ...] )
//	        [ column_properties ] [ out_of_line_part_storage ]
//
//	constraint_clauses:
//	    ADD { out_of_line_constraint | out_of_line_ref_constraint | constraint_state }
//
//	add_table_partition:
//	    ADD { PARTITION ... }
//
//	add_period_clause:
//	    ADD ( PERIOD FOR valid_time_column [ ( start, end ) ] )
//
//	add_range_subpartition / add_hash_subpartition / add_list_subpartition
//	    (within MODIFY PARTITION context, but can also follow ADD in partition modifications)
func (p *Parser) parseAlterTableAdd() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume ADD
	if p.isStatementEnd() {

		return nil, p.syntaxErrorAtCur(

		// ADD SUPPLEMENTAL LOG ...
		)
	}

	if p.isIdentLikeStr("SUPPLEMENTAL") {
		p.advance() // consume SUPPLEMENTAL
		if p.isIdentLikeStr("LOG") {
			p.advance() // consume LOG
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "ADD SUPPLEMENTAL LOG",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// ADD OVERFLOW
	if p.isIdentLikeStr("OVERFLOW") {
		p.advance() // consume OVERFLOW
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "ADD OVERFLOW",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// Check if this starts a constraint: CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, CHECK
	if p.isTableConstraintStart() {
		tc, parseErr134 := p.parseTableConstraint()
		if parseErr134 != nil {
			return nil, parseErr134
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_CONSTRAINT,
			Constraint: tc,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// ADD PARTITION
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		name := ""
		if p.isIdentLike() && !p.isAlterTablePartitionKeyword() {
			var parseErr135 error
			name, parseErr135 = p.parseIdentifier()
			if parseErr135 !=

				// Skip partition details (VALUES, table_partition_description, subpartitions, update_index_clauses)
				nil {
				return nil, parseErr135
			}
		}

		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_PARTITION,
			Subtype:    "PARTITION",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// ADD SUBPARTITION (within MODIFY PARTITION context — rare at top level)
	if p.cur.Type == kwSUBPARTITION {
		p.advance() // consume SUBPARTITION
		name := ""
		if p.isIdentLike() {
			var parseErr136 error
			name, parseErr136 = p.parseIdentifier()
			if parseErr136 != nil {
				return nil, parseErr136
			}
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_PARTITION,
			Subtype:    "SUBPARTITION",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// ADD ( column_def [, ...] ) or ADD ( PERIOD FOR ... )
	if p.cur.Type == '(' {
		p.advance() // consume '('

		// ADD ( PERIOD FOR valid_time_column [ ( start, end ) ] )
		if p.isIdentLikeStr("PERIOD") {
			p.advance() // consume PERIOD
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
			}
			name, parseErr137 := p.parseIdentifier()
			if parseErr137 !=
				// optional ( start_time_column , end_time_column )
				nil {
				return nil, parseErr137
			}

			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
			return &nodes.AlterTableCmd{
				Action:     nodes.AT_ADD_PERIOD,
				ColumnName: name,
				Loc:        nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// Parse column definitions (one or more)
		cols := &nodes.List{}
		col, parseErr138 := p.parseColumnDef()
		if parseErr138 != nil {
			return nil, parseErr138
		}
		cols.Items = append(cols.Items, col)
		for p.cur.Type == ',' {
			p.advance()
			var // consume ','
			parseErr139 error
			col, parseErr139 = p.parseColumnDef()
			if parseErr139 != nil {
				return nil, parseErr139
			}
			cols.Items = append(cols.Items, col)
		}
		if p.cur.Type == ')' {
			p.advance() // consume ')'
		}
		options := p.collectAlterTableClauseDetails()

		if len(cols.Items) == 1 {
			return withAlterTableOptions(&nodes.AlterTableCmd{
				Action:    nodes.AT_ADD_COLUMN,
				ColumnDef: col,
				Loc:       nodes.Loc{Start: start, End: p.prev.End},
			}, options), nil
		}
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_COLUMN,
			ColumnDefs: cols,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// ADD column_def (without parentheses)
	col, parseErr140 := p.parseColumnDef()
	if parseErr140 != nil {
		return nil, parseErr140
	}
	return &nodes.AlterTableCmd{
		Action:    nodes.AT_ADD_COLUMN,
		ColumnDef: col,
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableModify parses MODIFY actions.
//
//	modify_column_clauses:
//	    MODIFY { ( modify_col_properties [, ...] )
//	           | ( modify_virtcol_properties [, ...] )
//	           | modify_col_visibility
//	           | modify_col_substitutable }
//
//	constraint_clauses:
//	    MODIFY { CONSTRAINT name | PRIMARY KEY | UNIQUE (cols) } constraint_state [ CASCADE ]
//
//	MODIFY PARTITION partition { partition_attributes | ... }
//	MODIFY SUBPARTITION subpartition { ... }
//	MODIFY DEFAULT ATTRIBUTES [ FOR PARTITION partition ] ...
//	MODIFY CLUSTERING attribute_clustering_clause
//	MODIFY NESTED TABLE collection_item RETURN AS { LOCATOR | VALUE }
//	MODIFY LOB ( lob_item ) ( modify_LOB_parameters )
//	MODIFY VARRAY varray_item ( modify_LOB_parameters )
//	MODIFY OPAQUE TYPE column_name ...
//	MODIFY PARTITIONSET partitionset ...
func (p *Parser) parseAlterTableModify() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume MODIFY

	// MODIFY CONSTRAINT name constraint_state [CASCADE]
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, parseErr141 := p.parseIdentifier()
		if parseErr141 !=
			// Skip constraint_state and CASCADE
			nil {
			return nil, parseErr141
		}

		p.skipConstraintState()
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		tc := &nodes.TableConstraint{
			Name: name,
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_CONSTRAINT,
			Constraint: tc,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY PRIMARY KEY constraint_state [CASCADE]
	if p.cur.Type == kwPRIMARY {
		p.advance() // consume PRIMARY
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		p.skipConstraintState()
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_MODIFY_CONSTRAINT,
			Subtype: "PRIMARY KEY",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY UNIQUE (column [, column]...) constraint_state [CASCADE]
	if p.cur.Type == kwUNIQUE {
		p.advance() // consume UNIQUE
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		p.skipConstraintState()
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_MODIFY_CONSTRAINT,
			Subtype: "UNIQUE",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY PARTITION
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if p.isAlterTableDetailClauseStart() || p.isAlterTableActionStart() || p.isStatementEnd() {
			return nil, p.syntaxErrorAtCur()
		}
		name, parseErr142 := p.parseIdentifier()
		if parseErr142 != nil {
			return nil, parseErr142
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_PARTITION,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY SUBPARTITION
	if p.cur.Type == kwSUBPARTITION {
		p.advance() // consume SUBPARTITION
		if p.isAlterTableDetailClauseStart() || p.isAlterTableActionStart() || p.isStatementEnd() {
			return nil, p.syntaxErrorAtCur()
		}
		name, parseErr143 := p.parseIdentifier()
		if parseErr143 != nil {
			return nil, parseErr143
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_SUBPARTITION,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY DEFAULT ATTRIBUTES
	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance() // consume ATTRIBUTES
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action: nodes.AT_MODIFY_DEFAULT_ATTRS,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY NESTED TABLE collection_item RETURN AS { LOCATOR | VALUE }
	if p.cur.Type == kwNESTED {
		p.advance() // consume NESTED
		if p.cur.Type == kwTABLE {
			p.advance() // consume TABLE
		}
		name, parseErr144 := p.parseIdentifier()
		if parseErr144 !=
			// RETURN AS { LOCATOR | VALUE }
			nil {
			return nil, parseErr144
		}

		if p.isIdentLikeStr("RETURN") {
			p.advance()
		}
		if p.cur.Type == kwAS {
			p.advance()
		}
		subtype := ""
		if p.isIdentLikeStr("LOCATOR") || p.isIdentLikeStr("VALUE") {
			subtype = p.cur.Str
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_NESTED_TABLE,
			ColumnName: name,
			Subtype:    subtype,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY LOB (lob_item) (modify_LOB_parameters)
	if p.isIdentLikeStr("LOB") {
		p.advance() // consume LOB
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		return &nodes.AlterTableCmd{
			Action: nodes.AT_MODIFY_LOB,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY VARRAY varray_item (modify_LOB_parameters)
	if p.isIdentLikeStr("VARRAY") {
		p.advance()
		parseDiscard146, // consume VARRAY
			parseErr145 := p.parseIdentifier()
		_ = parseDiscard146
		if parseErr145 != nil {
			return nil, parseErr145
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		return &nodes.AlterTableCmd{
			Action: nodes.AT_MODIFY_VARRAY,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY OPAQUE TYPE column_name ...
	if p.isIdentLikeStr("OPAQUE") {
		p.advance() // consume OPAQUE
		if p.cur.Type == kwTYPE {
			p.advance() // consume TYPE
		}
		name, parseErr147 := p.parseIdentifier()
		if parseErr147 != nil {
			return nil, parseErr147
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_OPAQUE_TYPE,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY CLUSTERING attribute_clustering_clause
	if p.isIdentLikeStr("CLUSTERING") {
		p.advance() // consume CLUSTERING
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "MODIFY CLUSTERING",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY PARTITIONSET
	if p.isIdentLikeStr("PARTITIONSET") {
		p.advance() // consume PARTITIONSET
		name, parseErr148 := p.parseIdentifier()
		if parseErr148 != nil {
			return nil, parseErr148
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_PARTITIONSET,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MODIFY ( column_def [, ...] )
	if p.cur.Type == '(' {
		p.advance() // consume '('
		cols := &nodes.List{}
		col, parseErr149 := p.parseColumnDef()
		if parseErr149 != nil {
			return nil, parseErr149
		}
		cols.Items = append(cols.Items, col)
		for p.cur.Type == ',' {
			p.advance()
			var // consume ','
			parseErr150 error
			col, parseErr150 = p.parseColumnDef()
			if parseErr150 != nil {
				return nil, parseErr150
			}
			cols.Items = append(cols.Items, col)
		}
		if p.cur.Type == ')' {
			p.advance() // consume ')'
		}
		if len(cols.Items) == 1 {
			return &nodes.AlterTableCmd{
				Action:    nodes.AT_MODIFY_COLUMN,
				ColumnDef: col,
				Loc:       nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_COLUMN,
			ColumnDefs: cols,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY column_def (without parentheses) — includes modify_col_visibility
	col, parseErr151 := p.parseColumnDef()
	if parseErr151 != nil {
		return nil, parseErr151
	}
	return &nodes.AlterTableCmd{
		Action:    nodes.AT_MODIFY_COLUMN,
		ColumnDef: col,
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableDrop parses DROP actions.
//
//	drop_column_clause:
//	    SET UNUSED { COLUMN column | ( column [, column ]... ) } ...
//	    DROP { COLUMN column | ( column [, column ]... ) } ...
//	    DROP { UNUSED COLUMNS | COLUMNS CONTINUE } ...
//
//	drop_constraint_clause:
//	    DROP { PRIMARY KEY | UNIQUE ( column [, column ]... ) } ...
//	    DROP CONSTRAINT constraint_name ...
//
//	drop_table_partition / drop_table_subpartition:
//	    DROP PARTITION { partition | FOR (...) } [, ...] ...
//	    DROP SUBPARTITION { subpartition | FOR (...) } [, ...] ...
//
//	drop_period_clause:
//	    DROP ( PERIOD FOR valid_time_column )
func (p *Parser) parseAlterTableDrop() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume DROP

	switch p.cur.Type {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		name, parseErr152 := p.parseIdentifier()
		if parseErr152 !=
			// Optional CASCADE CONSTRAINTS | INVALIDATE
			nil {
			return nil, parseErr152
		}

		p.skipDropColumnTrailing()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_COLUMN,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwCONSTRAINT:
		p.advance() // consume CONSTRAINT
		name, parseErr153 := p.parseIdentifier()
		if parseErr153 !=
			// Optional CASCADE, ONLINE
			nil {
			return nil, parseErr153
		}

		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		if p.cur.Type == kwONLINE {
			p.advance()
		}
		tc := &nodes.TableConstraint{
			Name: name,
			Loc:  nodes.Loc{Start: start, End: p.prev.End},
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_CONSTRAINT,
			Constraint: tc,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwPRIMARY:
		// DROP PRIMARY KEY [CASCADE] [{KEEP|DROP} INDEX] [ONLINE]
		p.advance() // consume PRIMARY
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		p.skipKeepDropIndex()
		if p.cur.Type == kwONLINE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_DROP_CONSTRAINT,
			Subtype: "PRIMARY KEY",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwUNIQUE:
		// DROP UNIQUE (column [, column]...) [CASCADE] [{KEEP|DROP} INDEX] [ONLINE]
		p.advance() // consume UNIQUE
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		p.skipKeepDropIndex()
		if p.cur.Type == kwONLINE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_DROP_CONSTRAINT,
			Subtype: "UNIQUE",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwPARTITION:
		p.advance() // consume PARTITION
		name, parseErr154 := p.parsePartitionNameOrFor()
		if parseErr154 !=
			// Additional partitions: , { partition | FOR (...) }
			nil {
			return nil, parseErr154
		}

		for p.cur.Type == ',' {
			p.advance()
			parseDiscard156, parseErr155 := p.parsePartitionNameOrFor()
			_ = parseDiscard156

			// [ update_index_clauses ] [ parallel_clause ]
			if parseErr155 != nil {
				return nil, parseErr155
			}
		}

		p.skipUpdateIndexAndParallel()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_PARTITION,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwSUBPARTITION:
		p.advance() // consume SUBPARTITION
		name, parseErr157 := p.parsePartitionNameOrFor()
		if parseErr157 != nil {
			return nil, parseErr157
		}
		for p.cur.Type == ',' {
			p.advance()
			parseDiscard159, parseErr158 := p.parsePartitionNameOrFor()
			_ = parseDiscard159
			if parseErr158 != nil {
				return nil, parseErr158
			}
		}
		p.skipUpdateIndexAndParallel()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_SUBPARTITION,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	default:
		// DROP (col_name, ...) or DROP ( PERIOD FOR ... )
		if p.cur.Type == '(' {
			p.advance() // consume '('

			// Check for PERIOD FOR — drop_period_clause
			if p.isIdentLikeStr("PERIOD") {
				p.advance() // consume PERIOD
				if p.cur.Type == kwFOR {
					p.advance() // consume FOR
				}
				name, parseErr160 := p.parseIdentifier()
				if parseErr160 != nil {
					return nil, parseErr160
				}
				if p.cur.Type == ')' {
					p.advance() // consume ')'
				}
				return &nodes.AlterTableCmd{
					Action:     nodes.AT_DROP_PERIOD,
					ColumnName: name,
					Loc:        nodes.Loc{Start: start, End: p.prev.End},
				}, nil
			}

			// DROP (col_name, ...) — drop multiple columns
			name, parseErr161 := p.parseIdentifier()
			if parseErr161 != nil {
				return nil, parseErr161
			}
			for p.cur.Type == ',' {
				p.advance()
				parseDiscard163, parseErr162 := p.parseIdentifier()
				_ = parseDiscard163
				if parseErr162 != nil {
					return nil, parseErr162
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
			p.skipDropColumnTrailing()
			return &nodes.AlterTableCmd{
				Action:     nodes.AT_DROP_COLUMN,
				ColumnName: name,
				Loc:        nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// DROP UNUSED COLUMNS or DROP COLUMNS CONTINUE
		if p.isIdentLikeStr("UNUSED") {
			p.advance() // consume UNUSED
			if p.isIdentLikeStr("COLUMNS") {
				p.advance() // consume COLUMNS
			}
			// [ CHECKPOINT integer ]
			if p.isIdentLikeStr("CHECKPOINT") {
				p.advance()
				if p.cur.Type == tokICONST {
					p.advance()
				}
			}
			return &nodes.AlterTableCmd{
				Action:  nodes.AT_DROP_UNUSED_COLUMNS,
				Subtype: "UNUSED COLUMNS",
				Loc:     nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}
		if p.isIdentLikeStr("COLUMNS") {
			p.advance() // consume COLUMNS
			if p.isIdentLikeStr("CONTINUE") {
				p.advance() // consume CONTINUE
			}
			if p.isIdentLikeStr("CHECKPOINT") {
				p.advance()
				if p.cur.Type == tokICONST {
					p.advance()
				}
			}
			return &nodes.AlterTableCmd{
				Action:  nodes.AT_DROP_UNUSED_COLUMNS,
				Subtype: "COLUMNS CONTINUE",
				Loc:     nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// DROP SUPPLEMENTAL LOG ...
		if p.isIdentLikeStr("SUPPLEMENTAL") {
			p.advance() // consume SUPPLEMENTAL
			if p.isIdentLikeStr("LOG") {
				p.advance() // consume LOG
			}
			options := p.collectAlterTableClauseDetails()
			return withAlterTableOptions(&nodes.AlterTableCmd{
				Action:  nodes.AT_ALTER_PROPERTY,
				Subtype: "DROP SUPPLEMENTAL LOG",
				Loc:     nodes.Loc{Start: start, End: p.prev.End},
			}, options), nil
		}

		// Bare DROP col_name (without COLUMN keyword)
		name, parseErr164 := p.parseIdentifier()
		if parseErr164 != nil {
			return nil, parseErr164
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_COLUMN,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
}

// parseAlterTableRename parses RENAME actions.
//
//	RENAME COLUMN old_name TO new_name
//	RENAME CONSTRAINT old_name TO new_name
//	RENAME TO new_table_name
//	RENAME { PARTITION partition | SUBPARTITION subpartition } TO new_name
//	RENAME LOB ( lob_item ) TO new_lob_name
func (p *Parser) parseAlterTableRename() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume RENAME

	switch p.cur.Type {
	case kwCOLUMN:
		p.advance() // consume COLUMN
		oldName, parseErr165 := p.parseIdentifier()
		if parseErr165 != nil {
			return nil, parseErr165
		}
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		newName, parseErr166 := p.parseIdentifier()
		if parseErr166 != nil {
			return nil, parseErr166
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_RENAME_COLUMN,
			ColumnName: oldName,
			NewName:    newName,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwCONSTRAINT:
		p.advance() // consume CONSTRAINT
		oldName, parseErr167 := p.parseIdentifier()
		if parseErr167 != nil {
			return nil, parseErr167
		}
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		newName, parseErr168 := p.parseIdentifier()
		if parseErr168 != nil {
			return nil, parseErr168
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_RENAME_CONSTRAINT,
			ColumnName: oldName,
			NewName:    newName,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwPARTITION:
		p.advance() // consume PARTITION
		oldName, parseErr169 := p.parseIdentifier()
		if parseErr169 != nil {
			return nil, parseErr169
		}
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		newName, parseErr170 := p.parseIdentifier()
		if parseErr170 != nil {
			return nil, parseErr170
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_RENAME_PARTITION,
			Subtype:    "PARTITION",
			ColumnName: oldName,
			NewName:    newName,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwSUBPARTITION:
		p.advance() // consume SUBPARTITION
		oldName, parseErr171 := p.parseIdentifier()
		if parseErr171 != nil {
			return nil, parseErr171
		}
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		newName, parseErr172 := p.parseIdentifier()
		if parseErr172 != nil {
			return nil, parseErr172
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_RENAME_PARTITION,
			Subtype:    "SUBPARTITION",
			ColumnName: oldName,
			NewName:    newName,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case kwTO:
		p.advance() // consume TO
		newName, parseErr173 := p.parseIdentifier()
		if parseErr173 != nil {
			return nil, parseErr173
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_RENAME,
			NewName: newName,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	default:
		// RENAME LOB (lob_item) TO new_lob_name
		if p.isIdentLikeStr("LOB") {
			p.advance() // consume LOB
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			if p.cur.Type == kwTO {
				p.advance() // consume TO
			}
			newName, parseErr174 := p.parseIdentifier()
			if parseErr174 != nil {
				return nil, parseErr174
			}
			return &nodes.AlterTableCmd{
				Action:  nodes.AT_ALTER_PROPERTY,
				Subtype: "RENAME LOB",
				NewName: newName,
				Loc:     nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// RENAME TO new_name (without explicit TO first)
		if p.isIdentLike() {
			newName, parseErr175 := p.parseIdentifier()
			if parseErr175 != nil {
				return nil, parseErr175
			}
			return &nodes.AlterTableCmd{
				Action:  nodes.AT_RENAME,
				NewName: newName,
				Loc:     nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}
		return &nodes.AlterTableCmd{
			Action: nodes.AT_RENAME,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
}

// parseAlterTableTruncatePartition parses TRUNCATE PARTITION/SUBPARTITION.
//
//	truncate_partition_subpart:
//	    TRUNCATE { PARTITION | SUBPARTITION }
//	        { name | FOR ( key_value [, key_value ]... ) }
//	        [, { name | FOR ( key_value [, key_value ]... ) } ]...
//	        [ { DROP [ ALL ] | REUSE } STORAGE ]
//	        [ update_index_clauses ] [ parallel_clause ]
//	        [ CASCADE ]
func (p *Parser) parseAlterTableTruncatePartition() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume TRUNCATE

	subtype := ""
	action := nodes.AT_TRUNCATE_PARTITION
	if p.cur.Type == kwPARTITION {
		subtype = "PARTITION"
		p.advance()
	} else if p.cur.Type == kwSUBPARTITION {
		subtype = "SUBPARTITION"
		p.advance()
	} else {
		return nil, nil
	}

	name, parseErr176 := p.parsePartitionNameOrFor()
	if parseErr176 != nil {
		return nil, parseErr176
	}
	for p.cur.Type == ',' {
		p.advance()
		parseDiscard178, parseErr177 := p.parsePartitionNameOrFor()
		_ = parseDiscard178

		// [ { DROP [ALL] | REUSE } STORAGE ]
		if parseErr177 != nil {
			return nil, parseErr177
		}
	}

	if p.cur.Type == kwDROP || p.isIdentLikeStr("REUSE") {
		p.advance()
		if p.cur.Type == kwALL {
			p.advance()
		}
		if p.isIdentLikeStr("STORAGE") {
			p.advance()
		}
	}
	p.skipUpdateIndexAndParallel()
	if p.cur.Type == kwCASCADE {
		p.advance()
	}

	return &nodes.AlterTableCmd{
		Action:     action,
		Subtype:    subtype,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableEnableDisable parses ENABLE/DISABLE actions.
//
//	enable_disable_clause:
//	    { ENABLE | DISABLE }
//	        [ VALIDATE | NOVALIDATE ]
//	        { UNIQUE ( column [, column ]... )
//	        | PRIMARY KEY
//	        | CONSTRAINT constraint_name
//	        }
//	        [ using_index_clause ]
//	        [ exceptions_clause ]
//	        [ CASCADE ]
//	        [ { KEEP | DROP } INDEX ]
//
//	{ ENABLE | DISABLE } TABLE LOCK
//	{ ENABLE | DISABLE } ALL TRIGGERS
//	{ ENABLE | DISABLE } ROW MOVEMENT
//	{ ENABLE | DISABLE } LOGICAL REPLICATION
func (p *Parser) parseAlterTableEnableDisable() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	subtype := p.cur.Str // "ENABLE" or "DISABLE"
	p.advance()          // consume ENABLE/DISABLE

	// ENABLE/DISABLE TABLE LOCK
	if p.cur.Type == kwTABLE {
		p.advance() // consume TABLE
		if p.cur.Type == kwLOCK {
			p.advance() // consume LOCK
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ENABLE_DISABLE_TABLE_LOCK,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// ENABLE/DISABLE ALL TRIGGERS
	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		if p.cur.Type == kwTRIGGER || p.isIdentLikeStr("TRIGGERS") {
			p.advance() // consume TRIGGERS
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ENABLE_DISABLE_TRIGGERS,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// ENABLE/DISABLE ROW MOVEMENT
	if p.cur.Type == kwROW {
		p.advance() // consume ROW
		if p.isIdentLikeStr("MOVEMENT") {
			p.advance() // consume MOVEMENT
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: subtype + " ROW MOVEMENT",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// ENABLE/DISABLE LOGICAL REPLICATION
	if p.isIdentLikeStr("LOGICAL") {
		p.advance() // consume LOGICAL
		if p.isIdentLikeStr("REPLICATION") {
			p.advance() // consume REPLICATION
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: subtype + " LOGICAL REPLICATION",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// ENABLE/DISABLE [VALIDATE|NOVALIDATE] { UNIQUE (cols) | PRIMARY KEY | CONSTRAINT name }
	if p.isIdentLikeStr("VALIDATE") || p.isIdentLikeStr("NOVALIDATE") {
		p.advance()
	}

	name := ""
	constraintType := ""
	switch p.cur.Type {
	case kwUNIQUE:
		constraintType = "UNIQUE"
		p.advance()
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
	case kwPRIMARY:
		constraintType = "PRIMARY KEY"
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
		}
	case kwCONSTRAINT:
		constraintType = "CONSTRAINT"
		p.advance()
		var parseErr179 error
		name, parseErr179 = p.parseIdentifier()
		if parseErr179 != nil {
			return nil, parseErr179
		}
	default:
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ENABLE_DISABLE,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// [ USING INDEX ... ]
	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		} else if p.isIdentLike() && !p.isAlterTableActionStart() {
			// index attributes or index name
			p.collectAlterTableClauseDetails()
		}
	}

	// [ EXCEPTIONS INTO table ]
	if p.isIdentLikeStr("EXCEPTIONS") {
		p.advance()
		if p.cur.Type == kwINTO {
			p.advance()
			parseDiscard181, parseErr180 := p.parseObjectName()
			_ = parseDiscard181

			// [ CASCADE ]
			if parseErr180 != nil {
				return nil, parseErr180
			}
		}
	}

	if p.cur.Type == kwCASCADE {
		p.advance()
	}

	// [ { KEEP | DROP } INDEX ]
	p.skipKeepDropIndex()

	return &nodes.AlterTableCmd{
		Action:     nodes.AT_ENABLE_DISABLE,
		Subtype:    subtype + " " + constraintType,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableSet parses SET actions.
//
//	SET UNUSED { COLUMN column | ( column [, column ]... ) } ...
//	SET INTERVAL ( [ expr ] )
//	SET PARTITIONING { AUTOMATIC | MANUAL }
//	SET SUBPARTITION TEMPLATE ...
func (p *Parser) parseAlterTableSet() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume SET

	if p.isIdentLikeStr("UNUSED") {
		p.advance() // consume UNUSED
		name := ""
		if p.cur.Type == kwCOLUMN {
			p.advance()
			var // consume COLUMN
			parseErr182 error
			name, parseErr182 = p.parseIdentifier()
			if parseErr182 != nil {
				return nil, parseErr182
			}
		} else if p.cur.Type == '(' {
			p.advance()
			var // consume '('
			parseErr183 error
			name, parseErr183 = p.parseIdentifier()
			if parseErr183 != nil {
				return nil, parseErr183
			}
			for p.cur.Type == ',' {
				p.advance()
				parseDiscard185, parseErr184 := p.parseIdentifier()
				_ = parseDiscard185
				if parseErr184 != nil {
					return nil, parseErr184
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// [ CASCADE CONSTRAINTS | INVALIDATE ]
		p.skipDropColumnTrailing()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_SET_UNUSED,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.isIdentLikeStr("INTERVAL") {
		p.advance() // consume INTERVAL
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		return &nodes.AlterTableCmd{
			Action: nodes.AT_SET_INTERVAL,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.isIdentLikeStr("PARTITIONING") {
		p.advance() // consume PARTITIONING
		value := ""
		if p.isIdentLikeStr("AUTOMATIC") || p.isIdentLikeStr("MANUAL") {
			value = p.cur.Str
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_SET_PARTITIONING,
			Subtype: value,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == kwSUBPARTITION {
		p.advance() // consume SUBPARTITION
		if p.isIdentLikeStr("TEMPLATE") {
			p.advance() // consume TEMPLATE
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action: nodes.AT_SET_SUBPARTITION_TEMPLATE,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	options := p.collectAlterTableClauseDetails()
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "SET",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableMove parses MOVE [ONLINE] clause.
//
//	move_table_clause:
//	    MOVE [ ONLINE ]
//	        [ segment_attributes_clause ]
//	        [ table_compression ]
//	        [ index_org_table_clause ]
//	        [ { LOB (lob_item) STORE AS ... }... ]
//	        [ { VARRAY varray STORE AS ... }... ]
//	        [ parallel_clause ]
//	        [ allow_disallow_clustering ]
//	        [ update_index_clauses ]
//
//	move_table_partition:
//	    MOVE PARTITION partition ...
//
//	move_table_subpartition:
//	    MOVE SUBPARTITION subpartition ...
func (p *Parser) parseAlterTableMove() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume MOVE

	// MOVE PARTITION
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if p.isAlterTableDetailClauseStart() || p.isAlterTableActionStart() || p.isStatementEnd() {
			return nil, p.syntaxErrorAtCur()
		}
		name, parseErr186 := p.parseIdentifier()
		if parseErr186 != nil {
			return nil, parseErr186
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MOVE,
			Subtype:    "PARTITION",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MOVE SUBPARTITION
	if p.cur.Type == kwSUBPARTITION {
		p.advance() // consume SUBPARTITION
		if p.isAlterTableDetailClauseStart() || p.isAlterTableActionStart() || p.isStatementEnd() {
			return nil, p.syntaxErrorAtCur()
		}
		name, parseErr187 := p.parseIdentifier()
		if parseErr187 != nil {
			return nil, parseErr187
		}
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:     nodes.AT_MOVE,
			Subtype:    "SUBPARTITION",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// MOVE [ONLINE] [segment_attributes] [compression] ...
	options := p.collectAlterTableClauseDetails()
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action: nodes.AT_MOVE,
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableSplit parses SPLIT PARTITION/SUBPARTITION.
//
//	split_table_partition:
//	    SPLIT PARTITION partition
//	        { AT ( value [, value ]... ) [ INTO ( partition_spec, partition_spec ) ]
//	        | VALUES ( value [, value ]... ) [ INTO ( partition_spec, partition_spec ) ]
//	        | INTO ( partition_spec [, ...], partition_spec )
//	        }
//	        ...
//
//	split_table_subpartition:
//	    SPLIT SUBPARTITION subpartition ...
func (p *Parser) parseAlterTableSplit() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume SPLIT

	action := nodes.AT_SPLIT_PARTITION
	subtype := "PARTITION"
	if p.cur.Type == kwSUBPARTITION {
		action = nodes.AT_SPLIT_SUBPARTITION
		subtype = "SUBPARTITION"
	}
	p.advance() // consume PARTITION/SUBPARTITION

	name, parseErr188 := p.parseIdentifier()
	if parseErr188 != nil {
		return nil, parseErr188
	}
	options := p.collectAlterTableClauseDetails()

	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:     action,
		Subtype:    subtype,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableMerge parses MERGE PARTITIONS/SUBPARTITIONS.
//
//	merge_table_partitions:
//	    MERGE PARTITIONS
//	        { partition_or_key_value , partition_or_key_value [ INTO partition_spec ]
//	        | partition_or_key_value TO partition_or_key_value [ INTO partition_spec ]
//	        } ...
//
//	merge_table_subpartitions:
//	    MERGE SUBPARTITIONS ...
func (p *Parser) parseAlterTableMerge() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume MERGE

	action := nodes.AT_MERGE_PARTITIONS
	subtype := "PARTITIONS"
	if p.cur.Type == kwSUBPARTITION || p.isIdentLikeStr("SUBPARTITIONS") {
		action = nodes.AT_MERGE_SUBPARTITIONS
		subtype = "SUBPARTITIONS"
	}
	// consume PARTITIONS or SUBPARTITIONS
	p.advance()

	options := p.collectAlterTableClauseDetails()

	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  action,
		Subtype: subtype,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableExchange parses EXCHANGE PARTITION/SUBPARTITION.
//
//	exchange_partition_subpart:
//	    EXCHANGE { PARTITION partition | SUBPARTITION subpartition }
//	        WITH TABLE [ schema. ] table
//	        [ { INCLUDING | EXCLUDING } INDEXES ]
//	        [ { WITH | WITHOUT } VALIDATION ]
//	        [ exceptions_clause ]
//	        [ update_index_clauses ]
//	        [ parallel_clause ]
//	        [ CASCADE ]
//	        [ ONLINE ]
func (p *Parser) parseAlterTableExchange() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume EXCHANGE

	subtype := "PARTITION"
	if p.cur.Type == kwSUBPARTITION {
		subtype = "SUBPARTITION"
	}
	p.advance() // consume PARTITION/SUBPARTITION

	name, parseErr189 := p.parseIdentifier()
	if parseErr189 !=

		// WITH TABLE [schema.]table
		nil {
		return nil, parseErr189
	}

	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwTABLE {
			p.advance() // consume TABLE
		}
		parseDiscard191, parseErr190 := p.parseObjectName()
		_ = parseDiscard191

		if parseErr190 != nil {
			return nil, parseErr190
		}
	}

	options := p.collectAlterTableClauseDetails()

	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:     nodes.AT_EXCHANGE_PARTITION,
		Subtype:    subtype,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableCoalesce parses COALESCE PARTITION/SUBPARTITION.
//
//	coalesce_table_partition:
//	    COALESCE PARTITION [ update_index_clauses ] [ parallel_clause ]
//	        [ allow_disallow_clustering ]
//
//	coalesce_table_subpartition:
//	    COALESCE SUBPARTITION [ subpartition ]
//	        [ update_index_clauses ] [ parallel_clause ]
//	        [ allow_disallow_clustering ]
func (p *Parser) parseAlterTableCoalesce() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume COALESCE

	if p.cur.Type == kwSUBPARTITION {
		p.advance() // consume SUBPARTITION
		name := ""
		if p.isIdentLike() && !p.isAlterTableActionStart() {
			var parseErr192 error
			name, parseErr192 = p.parseIdentifier()
			if parseErr192 != nil {
				return nil, parseErr192
			}
		}
		p.skipUpdateIndexAndParallel()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_COALESCE_SUBPARTITION,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
	}
	p.skipUpdateIndexAndParallel()
	return &nodes.AlterTableCmd{
		Action: nodes.AT_COALESCE_PARTITION,
		Loc:    nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableShrinkSpace parses SHRINK SPACE [CASCADE].
//
//	SHRINK SPACE [ CASCADE ]
func (p *Parser) parseAlterTableShrinkSpace() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume SHRINK
	if p.isIdentLikeStr("SPACE") {
		p.advance() // consume SPACE
	}
	cascade := ""
	if p.cur.Type == kwCASCADE {
		cascade = "CASCADE"
		p.advance()
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_SHRINK_SPACE,
		Subtype: cascade,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableInmemory parses INMEMORY / NO INMEMORY clauses.
//
//	inmemory_table_clause:
//	    { INMEMORY [ inmemory_memcompress ] [ inmemory_priority ]
//	        [ inmemory_distribute ] [ inmemory_duplicate ]
//	    | NO INMEMORY
//	    }
//	    [ inmemory_column_clause ]
func (p *Parser) parseAlterTableInmemory() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume INMEMORY
	options := p.collectAlterTableClauseDetails()
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "INMEMORY",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableILM parses ILM clause.
//
//	ilm_clause:
//	    ILM { ADD POLICY ilm_policy_clause
//	        | DELETE POLICY policy_name
//	        | DELETE_ALL
//	        | ENABLE_ALL
//	        | DISABLE_ALL
//	        }
func (p *Parser) parseAlterTableILM() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume ILM
	options := p.collectAlterTableClauseDetails()
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "ILM",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableSupplementalLog parses SUPPLEMENTAL LOG clause.
//
//	supplemental_table_logging:
//	    { ADD SUPPLEMENTAL LOG
//	        { DATA ( { ALL | PRIMARY KEY | UNIQUE | FOREIGN KEY } [, ...] ) COLUMNS
//	        | GROUP log_group ( column [ NO LOG ] [, ...] ) [ ALWAYS ]
//	        }
//	    | DROP SUPPLEMENTAL LOG ...
//	    }
func (p *Parser) parseAlterTableSupplementalLog() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	// ADD SUPPLEMENTAL LOG is handled in parseAlterTableAdd
	// This handles standalone SUPPLEMENTAL LOG (should not normally be called)
	p.advance() // consume SUPPLEMENTAL
	options := p.collectAlterTableClauseDetails()
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "SUPPLEMENTAL LOG",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableAllocateExtent parses ALLOCATE EXTENT clause.
//
//	allocate_extent_clause:
//	    ALLOCATE EXTENT [ ( { SIZE size_clause | DATAFILE 'filename' | INSTANCE integer }... ) ]
func (p *Parser) parseAlterTableAllocateExtent() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume ALLOCATE
	if p.isIdentLikeStr("EXTENT") {
		p.advance() // consume EXTENT
	}
	if p.cur.Type == '(' {
		p.skipParenthesized()
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "ALLOCATE EXTENT",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableDeallocateUnused parses DEALLOCATE UNUSED clause.
//
//	deallocate_unused_clause:
//	    DEALLOCATE UNUSED [ KEEP size_clause ]
func (p *Parser) parseAlterTableDeallocateUnused() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume DEALLOCATE
	if p.isIdentLikeStr("UNUSED") {
		p.advance() // consume UNUSED
	}
	if p.isIdentLikeStr("KEEP") {
		p.advance() // consume KEEP
		p.skipSizeClause()
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "DEALLOCATE UNUSED",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableFlashbackArchive parses FLASHBACK ARCHIVE / NO FLASHBACK ARCHIVE.
//
//	flashback_archive_clause:
//	    FLASHBACK ARCHIVE [ flashback_archive ]
//	    | NO FLASHBACK ARCHIVE
func (p *Parser) parseAlterTableFlashbackArchive() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume FLASHBACK
	if p.isIdentLikeStr("ARCHIVE") {
		p.advance() // consume ARCHIVE
		name := ""
		if p.isIdentLike() && !p.isStatementEnd() && !p.isAlterTableActionStart() {
			var parseErr193 error
			name, parseErr193 = p.parseIdentifier()
			if parseErr193 != nil {
				return nil, parseErr193
			}
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_ALTER_PROPERTY,
			Subtype:    "FLASHBACK ARCHIVE",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "FLASHBACK",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableReadOnly parses READ ONLY / READ WRITE.
//
//	read_only_clause:
//	    { READ ONLY | READ WRITE }
func (p *Parser) parseAlterTableReadOnly() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume READ
	value := "READ"
	if p.isIdentLikeStr("ONLY") {
		value = "READ ONLY"
		p.advance()
	} else if p.cur.Type == kwWRITE {
		value = "READ WRITE"
		p.advance()
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: value,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableRowStore parses ROW STORE COMPRESS clause.
//
//	ROW STORE COMPRESS [ BASIC | ADVANCED ]
func (p *Parser) parseAlterTableRowStore() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	// Check ROW STORE COMPRESS vs ROW MOVEMENT
	next := p.peekNext()
	if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "STORE" {
		p.advance() // consume ROW
		p.advance() // consume STORE
		if p.cur.Type == kwCOMPRESS {
			p.advance() // consume COMPRESS
			if p.isIdentLikeStr("BASIC") || p.isIdentLikeStr("ADVANCED") {
				p.advance()
			}
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "ROW STORE COMPRESS",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
	// Not ROW STORE — shouldn't get here normally
	return nil, nil
}

// parseAlterTableRefresh parses REFRESH / NO REFRESH.
//
//	duplicated_table_refresh:
//	    { REFRESH | NO REFRESH }
func (p *Parser) parseAlterTableRefresh() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume REFRESH
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_DUPLICATED_REFRESH,
		Subtype: "REFRESH",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableProperty parses generic property changes.
//
//	logging_clause:       { LOGGING | NOLOGGING }
//	table_compression:    { COMPRESS [BASIC] | NOCOMPRESS | ... }
//	parallel_clause:      { NOPARALLEL | PARALLEL [integer] }
//	RESULT_CACHE, CACHE, NOCACHE, INDEXING, SEGMENT, ANNOTATIONS, etc.
func (p *Parser) parseAlterTableProperty() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	p.advance()

	// PARALLEL may have an integer degree
	if propName == "PARALLEL" && p.cur.Type == tokICONST {
		p.advance()
	}

	// COMPRESS may have BASIC
	if propName == "COMPRESS" && (p.isIdentLikeStr("BASIC") || p.isIdentLikeStr("ADVANCED")) {
		p.advance()
	}

	// COLUMN STORE COMPRESS
	var options *nodes.List
	if propName == "COLUMN" && p.isIdentLikeStr("STORE") {
		p.advance() // consume STORE
		if p.cur.Type == kwCOMPRESS {
			p.advance() // consume COMPRESS
			options = p.collectAlterTableClauseDetails()
		}
		propName = "COLUMN STORE COMPRESS"
	}

	// RESULT_CACHE ( MODE { DEFAULT | FORCE } )
	if propName == "RESULT_CACHE" && p.cur.Type == '(' {
		p.skipParenthesized()
	}

	// INDEXING { ON | OFF }
	if propName == "INDEXING" {
		if p.cur.Type == kwON || p.isIdentLikeStr("OFF") {
			p.advance()
		}
	}

	// SEGMENT CREATION { IMMEDIATE | DEFERRED }
	if propName == "SEGMENT" {
		if p.isIdentLikeStr("CREATION") {
			p.advance()
			if p.cur.Type == kwIMMEDIATE || p.cur.Type == kwDEFERRED {
				p.advance()
			}
		}
	}

	// ANNOTATIONS (...)
	if propName == "ANNOTATIONS" && p.cur.Type == '(' {
		p.skipParenthesized()
	}

	// UPGRADE [NOT] INCLUDING DATA [column_properties]
	if propName == "UPGRADE" {
		if p.cur.Type == kwNOT {
			p.advance()
		}
		if p.isIdentLikeStr("INCLUDING") {
			p.advance()
		}
		if p.isIdentLikeStr("DATA") {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
	}

	// OVERFLOW / MAPPING TABLE / NOMAPPING
	if propName == "OVERFLOW" || propName == "MAPPING" || propName == "NOMAPPING" {
		options = p.collectAlterTableClauseDetails()
	}

	// MEMOPTIMIZE FOR READ/WRITE
	if propName == "MEMOPTIMIZE" {
		if p.cur.Type == kwFOR {
			p.advance()
		}
		if p.cur.Type == kwREAD || p.cur.Type == kwWRITE {
			p.advance()
		}
	}

	// RECORDS_PER_BLOCK integer
	if propName == "RECORDS_PER_BLOCK" && p.cur.Type == tokICONST {
		p.advance()
	}

	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTablePhysicalAttr parses PCTFREE/PCTUSED/INITRANS/PCTTHRESHOLD integer.
func (p *Parser) parseAlterTablePhysicalAttr() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	p.advance() // consume PCTFREE/PCTUSED/INITRANS/PCTTHRESHOLD
	if p.cur.Type == tokICONST {
		p.advance() // consume integer
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableStorageClause parses STORAGE (...).
func (p *Parser) parseAlterTableStorageClause() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume STORAGE
	if p.cur.Type == '(' {
		p.skipParenthesized()
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "STORAGE",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableNoPrefix parses NO-prefixed properties.
//
//	NO INMEMORY
//	NO FLASHBACK ARCHIVE
//	NO MEMOPTIMIZE FOR READ
//	NO MEMOPTIMIZE FOR WRITE
//	NO REFRESH
//	NO DROP [UNTIL integer {DAY|DAYS} IDLE]  (immutable_table_clauses)
//	NO DELETE [LOCKED] [UNTIL ...]           (immutable_table_clauses / blockchain_table_clauses)
func (p *Parser) parseAlterTableNoPrefix() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume NO

	if p.isIdentLikeStr("INMEMORY") {
		p.advance()
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "NO INMEMORY",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == kwFLASHBACK {
		p.advance() // consume FLASHBACK
		if p.isIdentLikeStr("ARCHIVE") {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "NO FLASHBACK ARCHIVE",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.isIdentLikeStr("MEMOPTIMIZE") {
		p.advance() // consume MEMOPTIMIZE
		if p.cur.Type == kwFOR {
			p.advance()
		}
		dir := ""
		if p.cur.Type == kwREAD {
			dir = "READ"
			p.advance()
		} else if p.cur.Type == kwWRITE {
			dir = "WRITE"
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "NO MEMOPTIMIZE FOR " + dir,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.cur.Type == kwREFRESH {
		p.advance() // consume REFRESH
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_DUPLICATED_REFRESH,
			Subtype: "NO REFRESH",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	if p.isIdentLikeStr("DUPLICATE") {
		p.advance()
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "NO DUPLICATE",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// NO DROP [UNTIL integer {DAY|DAYS} IDLE]  — immutable/blockchain
	if p.cur.Type == kwDROP {
		p.advance() // consume DROP
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_IMMUTABLE_TABLE,
			Subtype: "NO DROP",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// NO DELETE [LOCKED] [UNTIL ...]  — immutable/blockchain
	if p.cur.Type == kwDELETE {
		p.advance() // consume DELETE
		options := p.collectAlterTableClauseDetails()
		return withAlterTableOptions(&nodes.AlterTableCmd{
			Action:  nodes.AT_IMMUTABLE_TABLE,
			Subtype: "NO DELETE",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, options), nil
	}

	// Generic NO prefix
	if p.isIdentLike() {
		name := p.cur.Str
		p.advance()
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "NO " + name,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: "NO",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// skipParenthesized skips a parenthesized group if present.
func (p *Parser) skipParenthesized() {
	if p.cur.Type != '(' {
		return
	}
	depth := 1
	p.advance() // consume '('
	for depth > 0 && p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
		}
		p.advance()
	}
}

func withAlterTableOptions(cmd *nodes.AlterTableCmd, options *nodes.List) *nodes.AlterTableCmd {
	if options != nil && len(options.Items) > 0 {
		cmd.Options = options
	}
	return cmd
}

// collectAlterTableClauseDetails preserves complex ALTER TABLE sub-clauses as
// ordered parser-visible options until the next ALTER TABLE action or statement
// end. Parenthesized payloads stay inside their owning option.
func (p *Parser) collectAlterTableClauseDetails() *nodes.List {
	options := &nodes.List{}
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isAlterTableActionStart() {
			return options
		}
		if p.cur.Type == '(' {
			start := p.pos()
			tokens := p.consumeAlterTableOptionTokens()
			if len(tokens) > 0 {
				options.Items = append(options.Items, &nodes.DDLOption{
					Key:   tokens[0],
					Value: strings.Join(tokens[1:], " "),
					Loc:   nodes.Loc{Start: start, End: p.prev.End},
				})
			}
			continue
		}
		start := p.pos()
		tokens := p.consumeAlterTableOptionTokens()
		if len(tokens) == 0 {
			continue
		}
		options.Items = append(options.Items, &nodes.DDLOption{
			Key:   tokens[0],
			Value: strings.Join(tokens[1:], " "),
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
		})
	}
	return options
}

func (p *Parser) consumeAlterTableOptionTokens() []string {
	var tokens []string
	depth := 0
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if depth == 0 && len(tokens) > 0 && (p.isAlterTableActionStart() || p.isAlterTableDetailClauseStart()) {
			break
		}
		if depth == 0 && len(tokens) == 0 && p.isAlterTableActionStart() {
			break
		}

		tokens = append(tokens, p.alterTableTokenText(p.cur))
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
		p.advance()
	}
	return tokens
}

func (p *Parser) alterTableTokenText(tok Token) string {
	if tok.Str != "" {
		return tok.Str
	}
	if tok.Loc >= 0 && tok.End <= len(p.source) && tok.Loc < tok.End {
		return p.source[tok.Loc:tok.End]
	}
	if tok.Type > 0 && tok.Type < 128 {
		return string(rune(tok.Type))
	}
	return ""
}

// isAlterTableActionStart returns true if the current token starts a new ALTER TABLE action.
func (p *Parser) isAlterTableActionStart() bool {
	switch p.cur.Type {
	case kwADD, kwMODIFY, kwDROP, kwRENAME, kwTRUNCATE,
		kwENABLE, kwDISABLE, kwSET,
		kwLOGGING, kwNOLOGGING,
		kwCOMPRESS, kwNOCOMPRESS,
		kwPARALLEL, kwNOPARALLEL,
		kwCACHE, kwREFRESH, kwFLASHBACK, kwREAD, kwROW,
		kwRESULT_CACHE:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "MOVE", "SPLIT", "MERGE", "EXCHANGE", "COALESCE",
			"SHRINK", "INMEMORY", "ILM", "SUPPLEMENTAL",
			"ALLOCATE", "DEALLOCATE", "NOCACHE", "NOCOMPRESS",
			"PCTFREE", "PCTUSED", "INITRANS", "STORAGE",
			"MEMOPTIMIZE", "NO", "OVERFLOW", "MAPPING", "NOMAPPING",
			"RECORDS_PER_BLOCK", "UPGRADE", "INDEXING",
			"SEGMENT", "ANNOTATIONS", "PCTTHRESHOLD",
			"HASHING", "VERSION":
			return true
		}
	}
	return false
}

func (p *Parser) isAlterTableDetailClauseStart() bool {
	switch p.cur.Type {
	case kwAS, kwCASCADE, kwDEFAULT, kwDROP, kwINTO, kwLOGGING, kwNOLOGGING,
		kwONLINE, kwPARALLEL, kwTABLESPACE, kwUPDATE, kwUSING, kwVALUES, kwWITH:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "ALLOW", "ALWAYS", "COMPRESS", "COLUMNS", "DATA", "DISALLOW",
			"EXCEPTIONS", "EXCLUDING", "GROUP", "INCLUDING", "INVALIDATE", "LOB",
			"NO", "NOCOMPRESS", "OVERFLOW", "STORE", "STORAGE", "SUBPARTITION",
			"TABLE", "TABLESPACE", "UPDATE", "VARRAY", "WITHOUT":
			return true
		}
	}
	return false
}

// isAlterTablePartitionKeyword checks if current token is a partition-related keyword
// that shouldn't be consumed as a partition name.
func (p *Parser) isAlterTablePartitionKeyword() bool {
	if p.cur.Type == kwVALUES || p.cur.Type == kwFOR {
		return true
	}
	if p.isIdentLikeStr("LESS") || p.isIdentLikeStr("SUBPARTITIONS") {
		return true
	}
	return false
}

// parsePartitionNameOrFor parses a partition name or FOR (key_value, ...).
func (p *Parser) parsePartitionNameOrFor() (string, error) {
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
		return "", nil
	}
	if p.isIdentLike() {
		return p.parseIdentifier()
	}
	return "", nil
}

// skipConstraintState skips constraint state tokens:
// [ [NOT] DEFERRABLE [INITIALLY {DEFERRED|IMMEDIATE}] ]
// [ RELY | NORELY ]
// [ USING INDEX ... ]
// [ { ENABLE | DISABLE } ]
// [ { VALIDATE | NOVALIDATE } ]
// [ exceptions_clause ]
func (p *Parser) skipConstraintState() {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch p.cur.Type {
		case kwNOT:
			p.advance()
			if p.cur.Type == kwDEFERRED || p.isIdentLikeStr("DEFERRABLE") {
				p.advance()
			}
		case kwDEFERRED:
			p.advance()
		case kwIMMEDIATE:
			p.advance()
		case kwENABLE, kwDISABLE:
			p.advance()
		case kwVALIDATE:
			p.advance()
		case kwUSING:
			p.advance()
			if p.cur.Type == kwINDEX {
				p.advance()
			}
			if p.cur.Type == '(' {
				p.skipParenthesized()
			} else if p.isIdentLike() && !p.isAlterTableActionStart() {
				p.parseIdentifier()
			}
		default:
			if p.isIdentLikeStr("DEFERRABLE") || p.isIdentLikeStr("INITIALLY") {
				p.advance()
				continue
			}
			if p.isIdentLikeStr("NOVALIDATE") || p.isIdentLikeStr("NORELY") {
				p.advance()
				continue
			}
			if p.cur.Type == kwRELY {
				p.advance()
				continue
			}
			if p.isIdentLikeStr("EXCEPTIONS") {
				p.advance()
				if p.cur.Type == kwINTO {
					p.advance()
					p.parseObjectName()
				}
				continue
			}
			return
		}
	}
}

// skipDropColumnTrailing skips CASCADE CONSTRAINTS, INVALIDATE, CHECKPOINT, ONLINE after DROP COLUMN.
func (p *Parser) skipDropColumnTrailing() {
	for {
		if p.cur.Type == kwCASCADE {
			p.advance()
			if p.cur.Type == kwCONSTRAINTS {
				p.advance()
			}
			continue
		}
		if p.isIdentLikeStr("INVALIDATE") {
			p.advance()
			continue
		}
		if p.isIdentLikeStr("CHECKPOINT") {
			p.advance()
			if p.cur.Type == tokICONST {
				p.advance()
			}
			continue
		}
		if p.cur.Type == kwONLINE {
			p.advance()
			continue
		}
		break
	}
}

// skipKeepDropIndex skips { KEEP | DROP } INDEX if present.
func (p *Parser) skipKeepDropIndex() {
	if p.isIdentLikeStr("KEEP") || p.cur.Type == kwDROP {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
	}
}

// skipUpdateIndexAndParallel skips update_index_clauses and parallel_clause.
func (p *Parser) skipUpdateIndexAndParallel() {
	// { UPDATE | INVALIDATE } { GLOBAL INDEXES | INDEXES [...] }
	if p.cur.Type == kwUPDATE || p.isIdentLikeStr("INVALIDATE") {
		p.advance()
		if p.isIdentLikeStr("GLOBAL") {
			p.advance()
		}
		if p.isIdentLikeStr("INDEXES") {
			p.advance()
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
	}

	// parallel_clause
	if p.cur.Type == kwPARALLEL {
		p.advance()
		if p.cur.Type == tokICONST {
			p.advance()
		}
	} else if p.cur.Type == kwNOPARALLEL {
		p.advance()
	}

	// allow_disallow_clustering
	if p.isIdentLikeStr("ALLOW") || p.isIdentLikeStr("DISALLOW") {
		p.advance()
		if p.isIdentLikeStr("CLUSTERING") {
			p.advance()
		}
	}

	// ONLINE
	if p.cur.Type == kwONLINE {
		p.advance()
	}
}

// skipSizeClause skips a size clause like "10M", "1G", etc.
func (p *Parser) skipSizeClause() {
	if p.cur.Type == tokICONST {
		p.advance()
		// optional unit suffix (K, M, G, T)
		if p.isIdentLike() {
			p.advance()
		}
	} else if p.isIdentLike() {
		p.advance()
	}
}
