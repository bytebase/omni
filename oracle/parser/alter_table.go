package parser

import (
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
//
// Oracle combines property, column and constraint clauses in one statement
// (ADD (c NUMBER) LOGGING, PCTFREE 5 NOPARALLEL), but refuses to combine
// some operations with any other (alterTableCmdAlone), and after an
// enable_disable_clause takes only more of them (ORA-03049 for ENABLE
// CONSTRAINT c LOGGING).
func (p *Parser) parseAlterTableStmt(start int) (*nodes.AlterTableStmt, error) {
	p.advance() // consume TABLE

	stmt := &nodes.AlterTableStmt{
		Actions: &nodes.List{},
		Loc:     nodes.Loc{Start: start},
	}
	var err error
	stmt.Name, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}

	// ALTER TABLE t alone is refused (ORA-02210).
	if p.isStatementEnd() {
		return nil, p.syntaxErrorAtCur()
	}
	enableDisable := false
	for !p.isStatementEnd() {
		tok := p.cur
		if enableDisable && p.cur.Type != kwENABLE && p.cur.Type != kwDISABLE {
			return nil, p.syntaxErrorAtCur()
		}
		cmd, err := p.parseAlterTableAction()
		if err != nil {
			return nil, err
		}
		if cmd == nil {
			return nil, p.syntaxErrorAtCur()
		}
		alone, last := alterTableCmdAlone(cmd)
		if alone && stmt.Actions.Len() > 0 {
			return nil, p.syntaxErrorAtTok(tok)
		}
		stmt.Actions.Items = append(stmt.Actions.Items, cmd)
		// The constraints of an ADD ( ... ) list follow its columns.
		for _, extra := range p.addedConstraints {
			stmt.Actions.Items = append(stmt.Actions.Items, extra)
		}
		p.addedConstraints = nil
		if (alone || last) && !p.isStatementEnd() {
			return nil, p.syntaxErrorAtCur()
		}
		switch cmd.Action {
		case nodes.AT_ENABLE_DISABLE, nodes.AT_ENABLE_DISABLE_TABLE_LOCK, nodes.AT_ENABLE_DISABLE_TRIGGERS:
			enableDisable = true
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// alterTableCmdAlone reports an operation Oracle runs only alone in its
// ALTER TABLE statement: a partition operation (ORA-14048), MOVE
// (ORA-14133), MODIFY DEFAULT ATTRIBUTES (ORA-14121), a RENAME (ORA-14047,
// ORA-23290) or a dropped or unused column (ORA-12987); and, as last, one
// nothing may follow: MODIFY CONSTRAINT and SHRINK SPACE (ORA-03049,
// ORA-10630 for what follows them).
func alterTableCmdAlone(cmd *nodes.AlterTableCmd) (alone, last bool) {
	switch cmd.Action {
	case nodes.AT_MODIFY_CONSTRAINT, nodes.AT_SHRINK_SPACE:
		return false, true
	case nodes.AT_ADD_PARTITION, nodes.AT_DROP_PARTITION, nodes.AT_DROP_SUBPARTITION,
		nodes.AT_TRUNCATE_PARTITION, nodes.AT_MODIFY_PARTITION, nodes.AT_MODIFY_SUBPARTITION,
		nodes.AT_MODIFY_DEFAULT_ATTRS, nodes.AT_MOVE, nodes.AT_SPLIT_PARTITION,
		nodes.AT_SPLIT_SUBPARTITION, nodes.AT_MERGE_PARTITIONS, nodes.AT_MERGE_SUBPARTITIONS,
		nodes.AT_EXCHANGE_PARTITION, nodes.AT_COALESCE_PARTITION, nodes.AT_RENAME_PARTITION,
		nodes.AT_SET_INTERVAL, nodes.AT_SET_PARTITIONING, nodes.AT_SET_SUBPARTITION_TEMPLATE,
		nodes.AT_RENAME, nodes.AT_RENAME_COLUMN, nodes.AT_RENAME_CONSTRAINT,
		nodes.AT_DROP_COLUMN, nodes.AT_SET_UNUSED, nodes.AT_DROP_UNUSED_COLUMNS:
		return true, false
	case nodes.AT_ALTER_PROPERTY:
		return cmd.Subtype == "RENAME LOB", false
	}
	return false, false
}

// parseAlterTableAction parses a single ALTER TABLE action, or returns nil
// when the current token starts none.
func (p *Parser) parseAlterTableAction() (*nodes.AlterTableCmd, error) {
	if p.isAlterTablePartitioningStart() {
		return p.parseAlterTablePartitioning()
	}
	switch p.cur.Type {
	case kwADD:
		return p.parseAlterTableAdd()
	case kwMODIFY:
		return p.parseAlterTableModify()
	case kwDROP:
		return p.parseAlterTableDrop()
	case kwRENAME:
		return p.parseAlterTableRename()
	case kwENABLE, kwDISABLE:
		return p.parseAlterTableEnableDisable()
	case kwSET:
		return p.parseAlterTableSet()
	case kwFLASHBACK:
		return p.parseAlterTableFlashbackArchive()
	case kwREAD:
		return p.parseAlterTableReadOnly()
	case kwDEFAULT:
		return p.parseAlterTableExternal()
	case kwREFRESH:
		return p.parseAlterTableRefresh()
	}
	switch {
	case p.cur.Type == kwLOGGING, p.cur.Type == kwNOLOGGING,
		p.cur.Type == kwPARALLEL, p.cur.Type == kwNOPARALLEL,
		p.cur.Type == kwCACHE, p.atWord("NOCACHE"),
		p.cur.Type == kwRESULT_CACHE, p.atWord("MEMOPTIMIZE"),
		p.atWord("UPGRADE"), p.atWord("ANNOTATIONS"),
		p.atWord("MINIMIZE"), p.atWord("NOMINIMIZE"):
		return p.parseAlterTableProperty()
	case p.isATTableCompressionStart():
		return p.parseAlterTableCompression()
	case p.atWord("INMEMORY"):
		return p.parseAlterTableInmemory()
	case p.atWord("ILM"):
		return p.parseAlterTableClause("ILM", p.parseATIlmClause)
	case p.atWord("ALLOCATE"):
		return p.parseAlterTableClause("ALLOCATE EXTENT", p.parseATAllocateExtent)
	case p.atWord("DEALLOCATE"):
		return p.parseAlterTableClause("DEALLOCATE UNUSED", p.parseATDeallocateUnused)
	case p.atWord("SHRINK"):
		return p.parseAlterTableShrinkSpace()
	case p.isATPhysicalAttrStart(), p.atWord("PCTTHRESHOLD"):
		return p.parseAlterTablePhysicalAttr()
	case p.atWord("OVERFLOW"), p.atWord("MAPPING"), p.atWord("NOMAPPING"),
		p.atWord("INCLUDING"), p.atWord("COALESCE"):
		return p.parseAlterTableIOT()
	case p.atWord("NO"):
		return p.parseAlterTableNoPrefix()
	case p.atWord("MOVE"):
		return p.parseAlterTableMove()
	case p.atWord("ACCESS"), p.atWord("LOCATION"), p.atWord("REJECT"), p.atWord("PROJECT"):
		return p.parseAlterTableExternal()
	}
	return nil, nil
}

// parseAlterTableClause parses a property clause read by parse and records
// it as one option.
func (p *Parser) parseAlterTableClause(subtype string, parse func() error) (*nodes.AlterTableCmd, error) {
	start := p.pos()
	if err := parse(); err != nil {
		return nil, err
	}
	cmd := &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: subtype,
		Options: &nodes.List{},
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}
	p.atOption(cmd.Options, start)
	return cmd, nil
}

// parseAlterTableAdd parses ADD column, constraint, period, overflow or
// supplemental logging; ADD PARTITION is alter_table_partitioning.
//
//	add_column_clause:
//	    ADD ( { column_definition | virtual_column_definition } [, ...] )
//	        [ column_properties ] [ out_of_line_part_storage ]
//
//	constraint_clauses:
//	    ADD { out_of_line_constraint | out_of_line_ref_constraint | constraint_state }
//
//	add_period_clause:
//	    ADD ( PERIOD FOR valid_time_column [ ( start, end ) ] )
//
//	alter_overflow_clause:
//	    ADD OVERFLOW [ segment_attributes_clause ]
//
//	supplemental_table_logging:
//	    ADD SUPPLEMENTAL LOG { supplemental_log_grp_clause | supplemental_id_key_clause } [, ...]
func (p *Parser) parseAlterTableAdd() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume ADD
	if p.isStatementEnd() {
		return nil, p.syntaxErrorAtCur()
	}

	if p.atWord("SUPPLEMENTAL") {
		p.advance() // consume SUPPLEMENTAL
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: "ADD SUPPLEMENTAL LOG", Options: &nodes.List{}}
		if err := p.parseATSupplementalLogging(true, cmd.Options); err != nil {
			return nil, err
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return cmd, nil
	}

	if p.atWord("OVERFLOW") {
		p.advance() // consume OVERFLOW
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: "ADD OVERFLOW", Options: &nodes.List{}}
		if p.isATSegmentAttrStart() {
			if err := p.parseATSegmentAttrs(cmd.Options); err != nil {
				return nil, err
			}
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil
	}

	// SCOPE FOR ( column ) IS table, or REF ( column ) WITH ROWID.
	if p.isOutOfLineRefConstraintStart() {
		tc, err := p.parseOutOfLineRefConstraint()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_CONSTRAINT,
			Constraint: tc,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// Check if this starts a constraint: CONSTRAINT, PRIMARY, UNIQUE, FOREIGN, CHECK
	if p.isTableConstraintStart() {
		tc, err := p.parseTableConstraint()
		if err != nil {
			return nil, err
		}
		// [ EXCEPTIONS INTO table ] — legal on ADD CONSTRAINT (unlike
		// CREATE TABLE constraints), verified against Oracle 23ai.
		if err := p.parseExceptionsIntoClause(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_ADD_CONSTRAINT,
			Constraint: tc,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// ADD ( column_def [, ...] ) or ADD ( PERIOD FOR ... )
	if p.cur.Type == '(' {
		p.advance() // consume '('

		// ADD ( PERIOD FOR valid_time_column [ ( start, end ) ] )
		if p.atWord("PERIOD") && p.peekNext().Type == kwFOR {
			p.advance() // consume PERIOD
			if err := p.atExpect(kwFOR); err != nil {
				return nil, err
			}
			name, err := p.parseATNameString()
			if err != nil {
				return nil, err
			}
			if p.cur.Type == '(' {
				p.advance()
				if err := p.parseATName(); err != nil {
					return nil, err
				}
				if err := p.atExpect(','); err != nil {
					return nil, err
				}
				if err := p.parseATName(); err != nil {
					return nil, err
				}
				if err := p.atExpect(')'); err != nil {
					return nil, err
				}
			}
			if err := p.atExpect(')'); err != nil {
				return nil, err
			}
			return &nodes.AlterTableCmd{
				Action:     nodes.AT_ADD_PERIOD,
				ColumnName: name,
				Loc:        nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// Column definitions and out-of-line constraints, in any order
		// (Oracle 23ai adds (m NUMBER, SCOPE FOR (d) IS t)); the
		// constraints become ADD CONSTRAINT actions after the columns'.
		cols := &nodes.List{}
		var col *nodes.ColumnDef
		var constraints []*nodes.AlterTableCmd
		for {
			itemStart := p.pos()
			var tc *nodes.TableConstraint
			var err error
			switch {
			case p.isOutOfLineRefConstraintStart():
				tc, err = p.parseOutOfLineRefConstraint()
			case p.isTableConstraintStart():
				tc, err = p.parseTableConstraint()
			default:
				col, err = p.parseColumnDef()
				if err == nil && col == nil {
					err = p.syntaxErrorAtCur()
				}
				if err == nil {
					cols.Items = append(cols.Items, col)
				}
			}
			if err != nil {
				return nil, err
			}
			if tc != nil {
				constraints = append(constraints, &nodes.AlterTableCmd{
					Action:     nodes.AT_ADD_CONSTRAINT,
					Constraint: tc,
					Loc:        nodes.Loc{Start: itemStart, End: p.prev.End},
				})
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume ')'
		// The list's REF constraints follow the CREATE TABLE rules over
		// its own columns (ORA-22893 for ADD (n NUMBER, SCOPE FOR (n) IS t)).
		listed := &nodes.CreateTableStmt{Columns: cols, Constraints: &nodes.List{}}
		for _, c := range constraints {
			listed.Constraints.Items = append(listed.Constraints.Items, c.Constraint)
		}
		if err := p.checkOutOfLineRefConstraints(listed, false); err != nil {
			return nil, err
		}
		options := &nodes.List{}
		if err := p.parseATColumnProperties(options); err != nil {
			return nil, err
		}
		if cols.Len() == 0 {
			// Constraints alone: the first is this action.
			p.addedConstraints = constraints[1:]
			constraints[0].Loc.Start = start
			return constraints[0], nil
		}
		p.addedConstraints = constraints

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
	col, err := p.parseColumnDef()
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, p.syntaxErrorAtCur()
	}
	// Its REF constraints follow the rules of a list's (ORA-02329 for ADD c
	// REF t PRIMARY KEY, ORA-30726 for REFERENCES o (n), ORA-22896 for
	// SCOPE and REFERENCES).
	single := &nodes.CreateTableStmt{Columns: &nodes.List{Items: []nodes.Node{col}}, Constraints: &nodes.List{}}
	if err := p.checkOutOfLineRefConstraints(single, false); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Action:    nodes.AT_ADD_COLUMN,
		ColumnDef: col,
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableModify parses MODIFY actions other than the partition
// operations of alter_table_partitioning.
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
//	MODIFY CLUSTERING [ clustering_when ] [ zonemap_clause ]
//	MODIFY NESTED TABLE collection_item RETURN AS { LOCATOR | VALUE }
//	MODIFY LOB ( lob_item ) ( modify_LOB_parameters )
//	MODIFY VARRAY varray_item ( modify_LOB_parameters )
//	MODIFY OPAQUE TYPE column STORE ( type [, type ]... ) UNPACKED
//	MODIFY table_partitioning_clauses [ filter_condition ] [ ONLINE ] [ update_index_clauses ]
//	MODIFY PARTITIONSET partitionset ...
func (p *Parser) parseAlterTableModify() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume MODIFY

	switch {
	case p.cur.Type == kwCONSTRAINT, p.cur.Type == kwCONSTRAINTS, p.cur.Type == kwPRIMARY, p.cur.Type == kwUNIQUE:
		return p.parseAlterTableModifyConstraint(start)
	case p.cur.Type == kwPARTITION, p.atWord("PARTITIONSET") && p.peekNext().Type == kwBY:
		// modify_to_partitioned
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_MODIFY_TO_PARTITIONED, Options: &nodes.List{}}
		if err := p.parseATTablePartitioning(cmd.Options); err != nil {
			return nil, err
		}
		if err := p.parseATClauses(atcFilter|atcOnline|atcIndexMaintenance, cmd.Options); err != nil {
			return nil, err
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil
	case p.atWord("PARTITIONSET"):
		// alter_table_partitionset: MODIFY PARTITIONSET partitionset
		// { ADD | MODIFY } PARTITIONSET partition_spec. Partitionsets
		// belong to sharded tables (ORA-02577 elsewhere).
		p.advance()
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		if p.cur.Type != kwADD && p.cur.Type != kwMODIFY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.atExpectWord("PARTITIONSET"); err != nil {
			return nil, err
		}
		if err := p.parseATName(); err != nil {
			return nil, err
		}
		if err := p.parseATClauses(atcTablePartitionDescription, nil); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_PARTITIONSET,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	case p.cur.Type == kwNESTED:
		// MODIFY NESTED TABLE collection_item RETURN AS { LOCATOR | VALUE }
		p.advance()
		if err := p.atExpect(kwTABLE); err != nil {
			return nil, err
		}
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		if err := p.atExpectWord("RETURN"); err != nil {
			return nil, err
		}
		if p.cur.Type == kwAS {
			p.advance()
		}
		if !p.atWord("LOCATOR") && !p.atWord("VALUE") {
			return nil, p.syntaxErrorAtCur()
		}
		subtype := p.cur.Str
		p.advance()
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_MODIFY_NESTED_TABLE,
			ColumnName: name,
			Subtype:    subtype,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	case p.atWord("LOB") && p.peekNext().Type == '(', p.atWord("VARRAY"):
		action := nodes.AT_MODIFY_LOB
		if p.atWord("VARRAY") {
			action = nodes.AT_MODIFY_VARRAY
		}
		if err := p.parseATLobModify(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action: action,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	case p.atWord("OPAQUE") && p.atNextWord("TYPE"):
		return p.parseAlterTableModifyOpaqueType(start)
	case p.atWord("CLUSTERING"):
		p.advance() // consume CLUSTERING
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: "MODIFY CLUSTERING", Options: &nodes.List{}}
		if err := p.parseATClusteringWhen(cmd.Options); err != nil {
			return nil, err
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil
	case p.cur.Type == kwCOLUMN:
		// modify_col_substitutable:
		//   COLUMN column [ NOT ] SUBSTITUTABLE AT ALL LEVELS [ FORCE ]
		p.advance()
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		if p.cur.Type == kwNOT {
			p.advance()
		}
		if err := p.atExpectWord("SUBSTITUTABLE"); err != nil {
			return nil, err
		}
		if err := p.atExpectWord("AT"); err != nil {
			return nil, err
		}
		if err := p.atExpect(kwALL); err != nil {
			return nil, err
		}
		// Oracle refuses the page's [ FORCE ] here (ORA-03049).
		if err := p.atExpectWord("LEVELS"); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_ALTER_PROPERTY,
			Subtype:    "MODIFY COLUMN SUBSTITUTABLE",
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// MODIFY ( column_def [, ...] )
	if p.cur.Type == '(' {
		p.advance() // consume '('
		cols := &nodes.List{}
		var col *nodes.ColumnDef
		for {
			var err error
			colTok := p.cur
			col, err = p.parseATModifyColumn(true)
			if err != nil {
				return nil, err
			}
			if col == nil {
				return nil, p.syntaxErrorAtCur()
			}
			// A column stands once in the list (ORA-00957).
			for _, item := range cols.Items {
				if item.(*nodes.ColumnDef).Name == col.Name {
					return nil, p.syntaxErrorAtTok(colTok)
				}
			}
			cols.Items = append(cols.Items, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
		if err := p.atExpect(')'); err != nil {
			return nil, err
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
	col, err := p.parseATModifyColumn(false)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.AlterTableCmd{
		Action:    nodes.AT_MODIFY_COLUMN,
		ColumnDef: col,
		Loc:       nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseATModifyColumn reads modify_col_properties. Every part after the
// column is optional there: Oracle runs ALTER TABLE t MODIFY c, MODIFY (c)
// and MODIFY c LOGGING, which CREATE TABLE's column_definition does not
// take. In a list (inList), a column followed by ',' or ')' is read alone;
// otherwise one followed by the end of the statement or by another action.
func (p *Parser) parseATModifyColumn(inList bool) (*nodes.ColumnDef, error) {
	if p.isATNameStart() {
		next := p.peekNext()
		var alone bool
		if inList {
			alone = next.Type == ',' || next.Type == ')'
		} else {
			alone = next.Type == ';' || next.Type == tokEOF || p.isATActionAfterColumn(next)
		}
		if alone {
			start := p.pos()
			name, err := p.parseATNameString()
			if err != nil {
				return nil, err
			}
			return &nodes.ColumnDef{
				Name:        name,
				Constraints: &nodes.List{},
				Loc:         nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}
	}
	return p.parseColumnDef()
}

// isATActionAfterColumn reports whether tok, after a column of MODIFY,
// starts another ALTER TABLE action rather than a column property (DROP
// IDENTITY is one).
func (p *Parser) isATActionAfterColumn(tok Token) bool {
	switch tok.Type {
	case kwDROP:
		return !p.isIdentLikeStrAt(p.peekAhead(2)[1], "IDENTITY")
	case kwADD, kwMODIFY, kwENABLE, kwDISABLE, kwLOGGING, kwNOLOGGING,
		kwPARALLEL, kwNOPARALLEL, kwCOMPRESS, kwNOCOMPRESS, kwCACHE, kwREAD, kwPCTFREE:
		return true
	}
	return false
}

// parseAlterTableModifyConstraint parses, after MODIFY,
//
//	{ CONSTRAINT constraint_name | PRIMARY KEY | UNIQUE ( column [, column ]... ) }
//	    constraint_state [ CASCADE ]
//
// Oracle needs a constraint_state that says something here (ORA-03049 for
// MODIFY CONSTRAINT c CASCADE, ORA-00905 for MODIFY PRIMARY KEY alone), and
// takes CONSTRAINTS for CONSTRAINT, as it does after ADD, DROP and ENABLE.
func (p *Parser) parseAlterTableModifyConstraint(start int) (*nodes.AlterTableCmd, error) {
	cmd := &nodes.AlterTableCmd{Action: nodes.AT_MODIFY_CONSTRAINT}
	tc := &nodes.TableConstraint{}
	switch p.cur.Type {
	case kwCONSTRAINT, kwCONSTRAINTS:
		p.advance()
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		tc.Name = name
	case kwPRIMARY:
		p.advance()
		if err := p.atExpect(kwKEY); err != nil {
			return nil, err
		}
		cmd.Subtype = "PRIMARY KEY"
		tc.Type = nodes.CONSTRAINT_PRIMARY
	default:
		p.advance() // UNIQUE
		if err := p.parseATNameList(); err != nil {
			return nil, err
		}
		cmd.Subtype = "UNIQUE"
		tc.Type = nodes.CONSTRAINT_UNIQUE
	}
	stateStart := p.pos()
	var cs constraintState
	if err := p.parseConstraintState(&cs); err != nil {
		return nil, err
	}
	if err := p.parseExceptionsIntoClause(); err != nil {
		return nil, err
	}
	if p.pos() == stateStart {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwCASCADE {
		p.advance()
	}
	tc.Deferrable = cs.Deferrable
	tc.Initially = cs.Initially
	tc.Tablespace = cs.Tablespace
	tc.UsingIndexLocal = cs.UsingIndexLocal
	tc.Loc = nodes.Loc{Start: start, End: p.prev.End}
	cmd.Constraint = tc
	cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return cmd, nil
}

// parseAlterTableModifyOpaqueType parses, after MODIFY,
//
//	OPAQUE TYPE column STORE ( type [, type ]... ) UNPACKED
//
// The page gives STORE AS [ LOB | OBJECT RELATIONAL ] ( column_properties ),
// which Oracle refuses whatever the parentheses hold (ORA-19001, ORA-22853,
// ORA-00906); the clause above is the one of its MODIFY OPAQUE TYPE.
func (p *Parser) parseAlterTableModifyOpaqueType(start int) (*nodes.AlterTableCmd, error) {
	p.advance() // OPAQUE
	p.advance() // TYPE
	name, err := p.parseATNameString()
	if err != nil {
		return nil, err
	}
	cmd := &nodes.AlterTableCmd{Action: nodes.AT_MODIFY_OPAQUE_TYPE, ColumnName: name, Options: &nodes.List{}}
	if err := p.atExpectWord("STORE"); err != nil {
		return nil, err
	}
	optStart := p.pos()
	if p.cur.Type == kwAS {
		// Oracle refuses what follows the page's STORE AS.
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}
	if err := p.atExpect('('); err != nil {
		return nil, err
	}
	for {
		if err := p.parseATQualifiedName(); err != nil {
			return nil, err
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.atExpect(')'); err != nil {
		return nil, err
	}
	if err := p.atExpectWord("UNPACKED"); err != nil {
		return nil, err
	}
	p.atOption(cmd.Options, optStart)
	cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return withAlterTableOptions(cmd, cmd.Options), nil
}

// parseAlterTableDrop parses DROP actions other than DROP PARTITION and
// DROP SUBPARTITION.
//
//	drop_column_clause:
//	    DROP { COLUMN column | ( column [, column ]... ) }
//	        [ { CASCADE CONSTRAINTS | INVALIDATE }... ] [ CHECKPOINT integer ]
//	    DROP { UNUSED COLUMNS | COLUMNS CONTINUE } [ CHECKPOINT integer ]
//
//	drop_constraint_clause:
//	    DROP { PRIMARY KEY | UNIQUE ( column [, column ]... ) }
//	        [ CASCADE ] [ { KEEP | DROP } INDEX ] [ ONLINE ]
//	    DROP CONSTRAINT constraint_name [ CASCADE ] [ ONLINE ]
//
//	drop_period_clause:
//	    DROP ( PERIOD FOR valid_time_column )
//
//	supplemental_table_logging:
//	    DROP SUPPLEMENTAL LOG { supplemental_id_key_clause | GROUP log_group } [, ...]
//
//	DROP CLUSTERING
func (p *Parser) parseAlterTableDrop() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume DROP

	switch {
	case p.cur.Type == kwCOLUMN:
		p.advance() // consume COLUMN
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		if err := p.parseATDropColumnTrailing(true); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_COLUMN,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == kwCONSTRAINT, p.cur.Type == kwCONSTRAINTS:
		p.advance() // consume CONSTRAINT
		name := ""
		var err error
		if name, err = p.parseATNameString(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		if p.cur.Type == kwONLINE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action: nodes.AT_DROP_CONSTRAINT,
			Constraint: &nodes.TableConstraint{
				Name: name,
				Loc:  nodes.Loc{Start: start, End: p.prev.End},
			},
			Loc: nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == kwPRIMARY, p.cur.Type == kwUNIQUE:
		subtype := "UNIQUE"
		if p.cur.Type == kwPRIMARY {
			subtype = "PRIMARY KEY"
			p.advance()
			if err := p.atExpect(kwKEY); err != nil {
				return nil, err
			}
		} else {
			p.advance()
			if err := p.parseATNameList(); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		if err := p.parseATKeepDropIndex(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwONLINE {
			p.advance()
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_DROP_CONSTRAINT,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == '(':
		p.advance() // consume '('

		// DROP ( PERIOD FOR valid_time_column )
		if p.atWord("PERIOD") {
			p.advance() // consume PERIOD
			if err := p.atExpect(kwFOR); err != nil {
				return nil, err
			}
			name, err := p.parseATNameString()
			if err != nil {
				return nil, err
			}
			if err := p.atExpect(')'); err != nil {
				return nil, err
			}
			return &nodes.AlterTableCmd{
				Action:     nodes.AT_DROP_PERIOD,
				ColumnName: name,
				Loc:        nodes.Loc{Start: start, End: p.prev.End},
			}, nil
		}

		// DROP ( column [, column ]... )
		name, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		for p.cur.Type == ',' {
			p.advance()
			if err := p.parseATName(); err != nil {
				return nil, err
			}
		}
		if err := p.atExpect(')'); err != nil {
			return nil, err
		}
		if err := p.parseATDropColumnTrailing(true); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     nodes.AT_DROP_COLUMN,
			ColumnName: name,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.atWord("UNUSED"), p.atWord("COLUMNS"):
		subtype := "UNUSED COLUMNS"
		if p.atWord("UNUSED") {
			p.advance()
			if err := p.atExpect(kwCOLUMNS); err != nil {
				return nil, err
			}
		} else {
			subtype = "COLUMNS CONTINUE"
			p.advance()
			if err := p.atExpectWord("CONTINUE"); err != nil {
				return nil, err
			}
		}
		if err := p.parseATCheckpoint(); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_DROP_UNUSED_COLUMNS,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.atWord("SUPPLEMENTAL"):
		p.advance() // consume SUPPLEMENTAL
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: "DROP SUPPLEMENTAL LOG", Options: &nodes.List{}}
		if err := p.parseATSupplementalLogging(false, cmd.Options); err != nil {
			return nil, err
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return cmd, nil

	case p.atWord("CLUSTERING"):
		p.advance()
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "DROP CLUSTERING",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseATDropColumnTrailing reads [ { CASCADE CONSTRAINTS | INVALIDATE }... ]
// and, after DROP, [ CHECKPOINT integer ] or, after SET UNUSED, [ ONLINE ].
func (p *Parser) parseATDropColumnTrailing(drop bool) error {
	for {
		if p.cur.Type == kwCASCADE {
			p.advance()
			if err := p.atExpect(kwCONSTRAINTS); err != nil {
				return err
			}
			continue
		}
		if p.atWord("INVALIDATE") {
			p.advance()
			continue
		}
		break
	}
	if drop {
		return p.parseATCheckpoint()
	}
	if p.cur.Type == kwONLINE {
		p.advance()
	}
	return nil
}

// parseATCheckpoint reads [ CHECKPOINT integer ].
func (p *Parser) parseATCheckpoint() error {
	if !p.atWord("CHECKPOINT") {
		return nil
	}
	p.advance()
	// Oracle takes CHECKPOINT without its integer as well.
	if p.cur.Type != tokICONST {
		return nil
	}
	return p.parseATInteger()
}

// parseATKeepDropIndex reads [ { KEEP | DROP } INDEX ].
func (p *Parser) parseATKeepDropIndex() error {
	if !p.atWord("KEEP") && p.cur.Type != kwDROP {
		return nil
	}
	p.advance()
	return p.atExpect(kwINDEX)
}

// parseAlterTableRename parses RENAME actions other than RENAME PARTITION
// and RENAME SUBPARTITION.
//
//	RENAME COLUMN old_name TO new_name
//	RENAME CONSTRAINT old_name TO new_name
//	RENAME TO new_table_name
//	RENAME LOB ( lob_item ) TO new_lob_name
func (p *Parser) parseAlterTableRename() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume RENAME

	renameTo := func() (string, error) {
		if err := p.atExpect(kwTO); err != nil {
			return "", err
		}
		return p.parseATNameString()
	}
	switch {
	case p.cur.Type == kwCOLUMN, p.cur.Type == kwCONSTRAINT:
		action := nodes.AT_RENAME_COLUMN
		if p.cur.Type == kwCONSTRAINT {
			action = nodes.AT_RENAME_CONSTRAINT
		}
		p.advance()
		oldName, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		newName, err := renameTo()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:     action,
			ColumnName: oldName,
			NewName:    newName,
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == kwTO:
		p.advance() // consume TO
		newName, err := p.parseATNameString()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_RENAME,
			NewName: newName,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.atWord("LOB"):
		p.advance() // consume LOB
		if err := p.atExpect('('); err != nil {
			return nil, err
		}
		if err := p.parseATItemPath(); err != nil {
			return nil, err
		}
		if err := p.atExpect(')'); err != nil {
			return nil, err
		}
		newName, err := renameTo()
		if err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: "RENAME LOB",
			NewName: newName,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
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
//	{ ENABLE | DISABLE } LOGICAL REPLICATION [ ALL KEYS | ALLOW NOVALIDATE KEYS ]
//
// CASCADE and { KEEP | DROP } INDEX follow DISABLE only (ORA-03049 after
// ENABLE).
func (p *Parser) parseAlterTableEnableDisable() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	disable := p.cur.Type == kwDISABLE
	subtype := p.cur.Str // "ENABLE" or "DISABLE"
	p.advance()          // consume ENABLE/DISABLE

	switch {
	case p.cur.Type == kwTABLE:
		p.advance() // consume TABLE
		if err := p.atExpect(kwLOCK); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ENABLE_DISABLE_TABLE_LOCK,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == kwALL:
		p.advance() // consume ALL
		if err := p.atExpectWord("TRIGGERS"); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ENABLE_DISABLE_TRIGGERS,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.cur.Type == kwROW:
		p.advance() // consume ROW
		if err := p.atExpectWord("MOVEMENT"); err != nil {
			return nil, err
		}
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: subtype + " ROW MOVEMENT",
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}, nil

	case p.atWord("LOGICAL"):
		p.advance() // consume LOGICAL
		if err := p.atExpectWord("REPLICATION"); err != nil {
			return nil, err
		}
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: subtype + " LOGICAL REPLICATION", Options: &nodes.List{}}
		optStart := p.pos()
		switch {
		case p.cur.Type == kwALL:
			p.advance()
			if err := p.atExpectWord("KEYS"); err != nil {
				return nil, err
			}
		case p.atWord("ALLOW"):
			p.advance()
			if err := p.atExpectWord("NOVALIDATE"); err != nil {
				return nil, err
			}
			if err := p.atExpectWord("KEYS"); err != nil {
				return nil, err
			}
		}
		p.atOption(cmd.Options, optStart)
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil
	}

	// ENABLE/DISABLE [VALIDATE|NOVALIDATE] { UNIQUE (cols) | PRIMARY KEY | CONSTRAINT name }
	if p.cur.Type == kwVALIDATE || p.atWord("NOVALIDATE") {
		p.advance()
	}

	name := ""
	constraintType := ""
	switch p.cur.Type {
	case kwUNIQUE:
		constraintType = "UNIQUE"
		p.advance()
		if err := p.parseATNameList(); err != nil {
			return nil, err
		}
	case kwPRIMARY:
		constraintType = "PRIMARY KEY"
		p.advance()
		if err := p.atExpect(kwKEY); err != nil {
			return nil, err
		}
	case kwCONSTRAINT, kwCONSTRAINTS:
		constraintType = "CONSTRAINT"
		p.advance()
		var err error
		if name, err = p.parseATNameString(); err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// [ USING INDEX ... ]
	if p.cur.Type == kwUSING && p.peekNext().Type == kwINDEX {
		var cs constraintState
		if err := p.parseUsingIndexClause(&cs); err != nil {
			return nil, err
		}
	}

	// [ EXCEPTIONS INTO table ]
	if err := p.parseExceptionsIntoClause(); err != nil {
		return nil, err
	}

	if disable {
		// [ CASCADE ] [ { KEEP | DROP } INDEX ]
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		if err := p.parseATKeepDropIndex(); err != nil {
			return nil, err
		}
	}

	return &nodes.AlterTableCmd{
		Action:     nodes.AT_ENABLE_DISABLE,
		Subtype:    subtype + " " + constraintType,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableSet parses SET UNUSED; SET PARTITIONING, SET INTERVAL and
// SET SUBPARTITION TEMPLATE are alter_table_partitioning.
//
//	SET UNUSED { COLUMN column | ( column [, column ]... ) }
//	    [ { CASCADE CONSTRAINTS | INVALIDATE }... ] [ ONLINE ]
func (p *Parser) parseAlterTableSet() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume SET
	if err := p.atExpectWord("UNUSED"); err != nil {
		return nil, err
	}
	var name string
	var err error
	if p.cur.Type == kwCOLUMN {
		p.advance()
		if name, err = p.parseATNameString(); err != nil {
			return nil, err
		}
	} else {
		if err := p.atExpect('('); err != nil {
			return nil, err
		}
		if name, err = p.parseATNameString(); err != nil {
			return nil, err
		}
		for p.cur.Type == ',' {
			p.advance()
			if err := p.parseATName(); err != nil {
				return nil, err
			}
		}
		if err := p.atExpect(')'); err != nil {
			return nil, err
		}
	}
	if err := p.parseATDropColumnTrailing(false); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Action:     nodes.AT_SET_UNUSED,
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableMove parses move_table_clause; MOVE PARTITION and MOVE
// SUBPARTITION are alter_table_partitioning.
//
//	MOVE [ filter_condition ] [ ONLINE ]
//	    [ segment_attributes_clause ] [ table_compression ]
//	    [ index_org_table_clause ]
//	    [ { LOB_storage_clause | varray_col_properties }... ]
//	    [ parallel_clause ] [ allow_disallow_clustering ]
//	    [ update_index_clauses ]
//
// Oracle reads these in any order, as it reads a partition's, and refuses
// INVALIDATE GLOBAL INDEXES here (ORA-14133).
func (p *Parser) parseAlterTableMove() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume MOVE
	cmd := &nodes.AlterTableCmd{Action: nodes.AT_MOVE, Options: &nodes.List{}}
	const allowed = atcFilter | atcOnline | atcSegment | atcCompression | atcLobStorage |
		atcParallel | atcClustering | atcUpdateIndexes
	for {
		if err := p.parseATClauses(allowed, cmd.Options); err != nil {
			return nil, err
		}
		// index_org_table_clause: mapping_table_clauses | PCTTHRESHOLD
		// integer | prefix_compression | index_org_overflow_clause.
		optStart := p.pos()
		switch {
		case p.atWord("MAPPING"):
			p.advance()
			if err := p.atExpect(kwTABLE); err != nil {
				return nil, err
			}
		case p.atWord("NOMAPPING"):
			p.advance()
		case p.atWord("PCTTHRESHOLD"):
			p.advance()
			if err := p.parseATInteger(); err != nil {
				return nil, err
			}
		case p.atWord("INCLUDING"), p.atWord("OVERFLOW"):
			if err := p.parseATIndexOrgOverflow(); err != nil {
				return nil, err
			}
		default:
			cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
			return withAlterTableOptions(cmd, cmd.Options), nil
		}
		p.atOption(cmd.Options, optStart)
	}
}

// parseATIndexOrgOverflow reads index_org_overflow_clause:
// [ INCLUDING column_name ] OVERFLOW [ segment_attributes_clause ].
//
// Oracle takes INCLUDING column_name without OVERFLOW as well.
func (p *Parser) parseATIndexOrgOverflow() error {
	if p.atWord("INCLUDING") {
		p.advance()
		if err := p.parseATName(); err != nil {
			return err
		}
		if !p.atWord("OVERFLOW") {
			return nil
		}
	}
	if err := p.atExpectWord("OVERFLOW"); err != nil {
		return err
	}
	if p.isATSegmentAttrStart() {
		return p.parseATSegmentAttrs(nil)
	}
	return nil
}

// parseAlterTableShrinkSpace parses SHRINK SPACE [ COMPACT ] [ CASCADE ].
func (p *Parser) parseAlterTableShrinkSpace() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	if err := p.parseATShrink(); err != nil {
		return nil, err
	}
	cascade := ""
	if p.prev.Type == kwCASCADE {
		cascade = "CASCADE"
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_SHRINK_SPACE,
		Subtype: cascade,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableInmemory parses inmemory_table_clause from INMEMORY:
//
//	{ INMEMORY [ inmemory_attributes ] | NO INMEMORY } [ inmemory_column_clause ]
//	inmemory_column_clause:
//	    { INMEMORY [ inmemory_memcompress ] | NO INMEMORY } ( column [, column ]... ) ...
//
// The table part may be left out: INMEMORY (c) is a column clause alone.
// NO INMEMORY reaches here through parseAlterTableNoPrefix.
func (p *Parser) parseAlterTableInmemory() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	cmd := &nodes.AlterTableCmd{Action: nodes.AT_ALTER_PROPERTY, Subtype: "INMEMORY", Options: &nodes.List{}}
	if p.atWord("NO") {
		cmd.Subtype = "NO INMEMORY"
		p.advance() // NO
		p.advance() // INMEMORY
		if p.cur.Type == '(' {
			if err := p.parseATNameList(); err != nil {
				return nil, err
			}
		}
	} else {
		p.advance() // INMEMORY
		if err := p.parseATInmemoryMemcompress(); err != nil {
			return nil, err
		}
		if p.cur.Type == '(' {
			// A column clause with no table part.
			if err := p.parseATNameList(); err != nil {
				return nil, err
			}
		} else if err := p.parseATInmemoryRest(); err != nil {
			return nil, err
		}
	}
	p.atOption(cmd.Options, start)
	for p.isATInmemoryStart() {
		optStart := p.pos()
		if err := p.parseATInmemoryColumnItem(); err != nil {
			return nil, err
		}
		p.atOption(cmd.Options, optStart)
	}
	cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return cmd, nil
}

// parseAlterTableFlashbackArchive parses flashback_archive_clause:
// FLASHBACK ARCHIVE [ flashback_archive ]; NO FLASHBACK ARCHIVE reaches
// here through parseAlterTableNoPrefix.
func (p *Parser) parseAlterTableFlashbackArchive() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume FLASHBACK
	if err := p.atExpectWord("ARCHIVE"); err != nil {
		return nil, err
	}
	name := ""
	if p.isATNameStart() && !p.isATAlterTableActionWord() {
		var err error
		if name, err = p.parseATNameString(); err != nil {
			return nil, err
		}
	}
	return &nodes.AlterTableCmd{
		Action:     nodes.AT_ALTER_PROPERTY,
		Subtype:    "FLASHBACK ARCHIVE",
		ColumnName: name,
		Loc:        nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// isATAlterTableActionWord reports a word that starts an ALTER TABLE
// action Oracle combines with others, which an optional name before it
// does not take.
func (p *Parser) isATAlterTableActionWord() bool {
	switch p.cur.Type {
	case kwADD, kwMODIFY, kwDROP, kwENABLE, kwDISABLE, kwLOGGING, kwNOLOGGING,
		kwPARALLEL, kwNOPARALLEL, kwCOMPRESS, kwNOCOMPRESS, kwCACHE, kwREAD, kwPCTFREE:
		return true
	}
	return false
}

// parseAlterTableRefresh parses the refresh of a duplicated table. The
// page gives duplicated_table_refresh as { REFRESH | NO REFRESH }; Oracle
// refuses both (ORA-00905 after REFRESH, ORA-03049 at NO) and takes
// REFRESH ON DEMAND.
func (p *Parser) parseAlterTableRefresh() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume REFRESH
	if err := p.atExpect(kwON); err != nil {
		return nil, err
	}
	if err := p.atExpectWord("DEMAND"); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_DUPLICATED_REFRESH,
		Subtype: "REFRESH ON DEMAND",
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableReadOnly parses read_only_clause: { READ ONLY | READ WRITE }.
func (p *Parser) parseAlterTableReadOnly() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	p.advance() // consume READ
	value := "READ ONLY"
	if p.cur.Type == kwWRITE {
		value = "READ WRITE"
		p.advance()
	} else if err := p.atExpectWord("ONLY"); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: value,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableCompression parses table_compression, or
// prefix_compression (COMPRESS integer) of an index-organized table.
func (p *Parser) parseAlterTableCompression() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	subtype := p.cur.Str
	switch p.cur.Type {
	case kwROW:
		subtype = "ROW STORE COMPRESS"
	case kwCOLUMN:
		subtype = "COLUMN STORE COMPRESS"
	}
	if err := p.parseATTableCompression(true); err != nil {
		return nil, err
	}
	cmd := &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: subtype,
		Options: &nodes.List{},
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}
	if subtype == "COLUMN STORE COMPRESS" {
		p.atOption(cmd.Options, start)
	}
	return withAlterTableOptions(cmd, cmd.Options), nil
}

// parseAlterTableProperty parses the single-clause properties:
//
//	logging_clause:       { LOGGING | NOLOGGING }
//	parallel_clause:      { NOPARALLEL | PARALLEL [ integer ] }
//	{ CACHE | NOCACHE }
//	RESULT_CACHE ( MODE { DEFAULT | FORCE } )
//	MEMOPTIMIZE FOR { READ | WRITE }
//	upgrade_table_clause: UPGRADE [ [ NOT ] INCLUDING DATA [ column_properties ] ]
//	annotations_clause:   ANNOTATIONS ( ... )
//	records_per_block_clause: { MINIMIZE | NOMINIMIZE } RECORDS_PER_BLOCK
//
// Oracle refuses the page's alter_XMLSchema_clause here (ORA-01735 at ALLOW
// ANYSCHEMA) and its RECORDS_PER_BLOCK integer, which is not a clause of
// ALTER TABLE: records_per_block_clause is.
func (p *Parser) parseAlterTableProperty() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	options := &nodes.List{}
	switch {
	case p.cur.Type == kwPARALLEL, p.cur.Type == kwNOPARALLEL:
		if err := p.parseATParallel(); err != nil {
			return nil, err
		}
	case p.cur.Type == kwRESULT_CACHE:
		p.advance()
		if err := p.atExpect('('); err != nil {
			return nil, err
		}
		if err := p.atExpect(kwMODE); err != nil {
			return nil, err
		}
		if err := p.atOneOf("DEFAULT", "FORCE"); err != nil {
			return nil, err
		}
		if err := p.atExpect(')'); err != nil {
			return nil, err
		}
	case p.atWord("MEMOPTIMIZE"):
		p.advance()
		if err := p.atExpect(kwFOR); err != nil {
			return nil, err
		}
		if p.cur.Type != kwREAD && p.cur.Type != kwWRITE {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	case p.atWord("UPGRADE"):
		// Oracle takes UPGRADE alone as well.
		p.advance()
		if p.cur.Type == kwNOT || p.atWord("INCLUDING") {
			if p.cur.Type == kwNOT {
				p.advance()
			}
			if err := p.atExpectWord("INCLUDING"); err != nil {
				return nil, err
			}
			if err := p.atExpectWord("DATA"); err != nil {
				return nil, err
			}
			if err := p.parseATColumnProperties(options); err != nil {
				return nil, err
			}
		}
	case p.atWord("ANNOTATIONS"):
		p.advance()
		if err := p.parseATAnnotations(); err != nil {
			return nil, err
		}
	case p.atWord("MINIMIZE"), p.atWord("NOMINIMIZE"):
		p.advance()
		if err := p.atExpectWord("RECORDS_PER_BLOCK"); err != nil {
			return nil, err
		}
		propName += " RECORDS_PER_BLOCK"
	default:
		// LOGGING, NOLOGGING, CACHE, NOCACHE
		p.advance()
	}
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseATAnnotations reads ( { ADD | DROP | REPLACE } annotation [, ...] )
// after ANNOTATIONS, where annotation is name [ 'value' ]; ADD may be left
// out.
func (p *Parser) parseATAnnotations() error {
	if err := p.atExpect('('); err != nil {
		return err
	}
	for {
		switch {
		case p.cur.Type == kwADD, p.cur.Type == kwDROP:
			p.advance()
			if p.cur.Type == kwIF {
				p.advance()
				if p.cur.Type == kwNOT {
					p.advance()
				}
				if err := p.atExpect(kwEXISTS); err != nil {
					return err
				}
			}
		case p.atWord("REPLACE"):
			p.advance()
		}
		if err := p.parseATName(); err != nil {
			return err
		}
		if p.cur.Type == tokSCONST {
			p.advance()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return p.atExpect(')')
}

// parseAlterTablePhysicalAttr parses PCTFREE, PCTUSED, INITRANS,
// storage_clause or PCTTHRESHOLD integer.
func (p *Parser) parseAlterTablePhysicalAttr() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	if p.atWord("PCTTHRESHOLD") {
		p.advance()
		if err := p.parseATInteger(); err != nil {
			return nil, err
		}
	} else if p.atWord("STORAGE") {
		// ALTER TABLE cannot change INITIAL (ORA-02203).
		if p.peekAhead(2)[1].Type == kwINITIAL {
			p.advance()
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.parseATStorageClause(); err != nil {
			return nil, err
		}
	} else if err := p.parseATPhysicalAttr(); err != nil {
		return nil, err
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableIOT parses alter_iot_clauses:
//
//	{ index_org_table_clause | alter_overflow_clause
//	| alter_mapping_table_clauses | COALESCE }
//	alter_overflow_clause:
//	    OVERFLOW { allocate_extent_clause | deallocate_unused_clause
//	             | { physical_attributes_clause | logging_clause }... }
//	alter_mapping_table_clauses:
//	    MAPPING TABLE { allocate_extent_clause | deallocate_unused_clause }
//	mapping_table_clauses: { MAPPING TABLE | NOMAPPING }
//	index_org_overflow_clause: [ INCLUDING column_name ] OVERFLOW [ segment_attributes_clause ]
//
// Oracle takes no TABLESPACE after OVERFLOW (ORA-01735), and each physical
// attribute once (ORA-14102 for OVERFLOW LOGGING LOGGING).
func (p *Parser) parseAlterTableIOT() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	options := &nodes.List{}
	switch {
	case p.atWord("COALESCE"):
		p.advance()
	case p.atWord("NOMAPPING"):
		// Oracle reads NOMAPPING here and refuses whatever follows it, the
		// end of the statement included (ORA-01735); MOVE takes it.
		p.advance()
		return nil, p.syntaxErrorAtCur()
	case p.atWord("MAPPING"):
		p.advance()
		if err := p.atExpect(kwTABLE); err != nil {
			return nil, err
		}
		optStart := p.pos()
		ok, err := p.parseATClause(atcAllocate|atcDeallocate, nil)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		p.atOption(options, optStart)
	case p.atWord("INCLUDING"):
		if err := p.parseATIndexOrgOverflow(); err != nil {
			return nil, err
		}
	default:
		p.advance() // OVERFLOW
		optStart := p.pos()
		switch {
		case p.isStatementEnd():
			// Oracle takes OVERFLOW alone as well.
		case p.isATAllocateExtentStart():
			if err := p.parseATAllocateExtent(); err != nil {
				return nil, err
			}
		case p.isATDeallocateStart():
			if err := p.parseATDeallocateUnused(); err != nil {
				return nil, err
			}
		default:
			seen := map[string]bool{}
			for n := 0; ; n++ {
				key := p.cur.Str
				if p.cur.Type == kwNOLOGGING {
					key = "LOGGING"
				}
				isAttr := p.cur.Type == kwLOGGING || p.cur.Type == kwNOLOGGING || p.isATPhysicalAttrStart()
				if !isAttr || seen[key] {
					if n == 0 || seen[key] {
						return nil, p.syntaxErrorAtCur()
					}
					break
				}
				seen[key] = true
				if p.isATLoggingStart() {
					p.advance()
				} else if err := p.parseATPhysicalAttr(); err != nil {
					return nil, err
				}
			}
		}
		p.atOption(options, optStart)
	}
	return withAlterTableOptions(&nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, options), nil
}

// parseAlterTableExternal parses the clauses of alter_external_table that
// no other table takes:
//
//	external_table_data_props:
//	    [ DEFAULT DIRECTORY directory ]
//	    [ ACCESS PARAMETERS { ( opaque_format_spec ) | USING CLOB subquery } ]
//	    [ LOCATION ( [ directory: ] 'location_specifier' [, ...] ) ]
//	REJECT LIMIT { integer | UNLIMITED }
//	PROJECT COLUMN { ALL | REFERENCED }
func (p *Parser) parseAlterTableExternal() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	propName := p.cur.Str
	switch {
	case p.cur.Type == kwDEFAULT, p.atWord("LOCATION"):
		if err := p.parseATExternalDataProp(); err != nil {
			return nil, err
		}
	case p.atWord("ACCESS"):
		p.advance()
		if err := p.atExpectWord("PARAMETERS"); err != nil {
			return nil, err
		}
		if p.cur.Type == '(' {
			p.skipParenthesized()
		} else {
			if err := p.atExpect(kwUSING); err != nil {
				return nil, err
			}
			if err := p.atExpect(kwCLOB); err != nil {
				return nil, err
			}
			if p.cur.Type != kwSELECT && p.cur.Type != kwWITH && p.cur.Type != '(' {
				return nil, p.syntaxErrorAtCur()
			}
			q, err := p.parseSelectStmt()
			if err != nil {
				return nil, err
			}
			if q == nil {
				return nil, p.syntaxErrorAtCur()
			}
		}
	case p.atWord("REJECT"):
		p.advance()
		if err := p.atExpectWord("LIMIT"); err != nil {
			return nil, err
		}
		if !p.atOptOneOf("UNLIMITED") {
			if err := p.parseATInteger(); err != nil {
				return nil, err
			}
		}
	default: // PROJECT
		p.advance()
		if err := p.atExpect(kwCOLUMN); err != nil {
			return nil, err
		}
		if p.cur.Type == kwALL {
			p.advance()
		} else if err := p.atExpectWord("REFERENCED"); err != nil {
			return nil, err
		}
	}
	return &nodes.AlterTableCmd{
		Action:  nodes.AT_ALTER_PROPERTY,
		Subtype: propName,
		Loc:     nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseAlterTableNoPrefix parses the clauses that start with NO:
//
//	NO INMEMORY [ inmemory_column_clause ]
//	NO FLASHBACK ARCHIVE
//	NO MEMOPTIMIZE FOR { READ | WRITE }
//	NO ROW ARCHIVAL
//	NO DROP [ UNTIL integer DAYS IDLE ]
//	NO DELETE [ LOCKED ] [ UNTIL integer { DAY | DAYS } AFTER INSERT ]
//
// Oracle takes only DAYS in NO DROP UNTIL (ORA-00905 for DAY).
func (p *Parser) parseAlterTableNoPrefix() (*nodes.AlterTableCmd, error) {
	start := p.pos()
	if p.atNextWord("INMEMORY") {
		return p.parseAlterTableInmemory()
	}
	p.advance() // consume NO
	prop := func(subtype string) *nodes.AlterTableCmd {
		return &nodes.AlterTableCmd{
			Action:  nodes.AT_ALTER_PROPERTY,
			Subtype: subtype,
			Loc:     nodes.Loc{Start: start, End: p.prev.End},
		}
	}

	switch {
	case p.cur.Type == kwFLASHBACK:
		p.advance() // consume FLASHBACK
		if err := p.atExpectWord("ARCHIVE"); err != nil {
			return nil, err
		}
		return prop("NO FLASHBACK ARCHIVE"), nil

	case p.atWord("MEMOPTIMIZE"):
		p.advance() // consume MEMOPTIMIZE
		if err := p.atExpect(kwFOR); err != nil {
			return nil, err
		}
		if p.cur.Type != kwREAD && p.cur.Type != kwWRITE {
			return nil, p.syntaxErrorAtCur()
		}
		dir := p.cur.Str
		p.advance()
		return prop("NO MEMOPTIMIZE FOR " + dir), nil

	case p.cur.Type == kwROW:
		p.advance()
		if err := p.atExpectWord("ARCHIVAL"); err != nil {
			return nil, err
		}
		return prop("NO ROW ARCHIVAL"), nil

	case p.cur.Type == kwDROP:
		p.advance() // consume DROP
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_IMMUTABLE_TABLE, Subtype: "NO DROP", Options: &nodes.List{}}
		if p.atWord("UNTIL") {
			optStart := p.pos()
			p.advance()
			if err := p.parseATInteger(); err != nil {
				return nil, err
			}
			if err := p.atExpectWord("DAYS"); err != nil {
				return nil, err
			}
			if err := p.atExpectWord("IDLE"); err != nil {
				return nil, err
			}
			p.atOption(cmd.Options, optStart)
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil

	case p.cur.Type == kwDELETE:
		p.advance() // consume DELETE
		cmd := &nodes.AlterTableCmd{Action: nodes.AT_IMMUTABLE_TABLE, Subtype: "NO DELETE", Options: &nodes.List{}}
		optStart := p.pos()
		if p.atWord("LOCKED") {
			p.advance()
			p.atOption(cmd.Options, optStart)
		}
		if p.atWord("UNTIL") {
			optStart := p.pos()
			p.advance()
			if err := p.parseATInteger(); err != nil {
				return nil, err
			}
			if err := p.atOneOf("DAY", "DAYS"); err != nil {
				return nil, err
			}
			if err := p.atExpectWord("AFTER"); err != nil {
				return nil, err
			}
			if err := p.atExpect(kwINSERT); err != nil {
				return nil, err
			}
			p.atOption(cmd.Options, optStart)
		}
		cmd.Loc = nodes.Loc{Start: start, End: p.prev.End}
		return withAlterTableOptions(cmd, cmd.Options), nil
	}
	return nil, p.syntaxErrorAtCur()
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
	} else {
		cmd.Options = nil
	}
	return cmd
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
