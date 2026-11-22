package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseClusterInteger reads an option's integer value into *v when v is
// not nil.
func (p *Parser) parseClusterInteger(v **int) error {
	if p.cur.Type != tokICONST {
		return p.syntaxErrorAtCur()
	}
	n, err := p.parseIntValue()
	if err != nil {
		return err
	}
	if v != nil {
		*v = &n
	}
	return nil
}

// parseClusterSize reads a size_clause: integer [ K | M | G | T | P | E ].
func (p *Parser) parseClusterSize() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.parseSizeValue()
}

// parseCreateClusterStmt parses a CREATE CLUSTER statement after CLUSTER.
//
// BNF: oracle/parser/bnf/CREATE-CLUSTER.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-CLUSTER.html
//
//	CREATE CLUSTER [ IF NOT EXISTS ] [ schema. ] cluster
//	    [ SHARING = { METADATA | NONE } ]
//	    ( column datatype [ COLLATE column_collation_name ] [ SORT ]
//	      [, column datatype [ COLLATE column_collation_name ] [ SORT ] ]... )
//	    [ physical_attributes_clause ]
//	    [ SIZE size_clause ]
//	    [ TABLESPACE tablespace ]
//	    [ { INDEX
//	      | HASHKEYS integer [ HASH IS expr ] }
//	    ]
//	    [ SINGLE TABLE ]
//	    [ parallel_clause ]
//	    [ { NOROWDEPENDENCIES | ROWDEPENDENCIES } ]
//	    [ { CACHE | NOCACHE } ]
//	    [ cluster_range_partitions ] ;
//
//	physical_attributes_clause:
//	    [ PCTFREE integer ] [ PCTUSED integer ] [ INITRANS integer ] [ storage_clause ]
//
//	parallel_clause:
//	    { NOPARALLEL | PARALLEL [ integer ] }
//
//	cluster_range_partitions:
//	    PARTITION BY RANGE ( column [, column ]... )
//	    ( PARTITION [ partition ]
//	        VALUES LESS THAN ( { value | MAXVALUE } [, { value | MAXVALUE } ]... )
//	        [ table_partition_description ]
//	      [, ... ]
//	    )
//
// Oracle 23ai takes the options after the column list in any order, and
// MAXTRANS among the physical attributes. It refuses range partitioning of
// a cluster that is not a hash cluster (ORA-14179), which the page states
// but its syntax does not, so omni does not check it.
func (p *Parser) parseCreateClusterStmt(start int) (*nodes.CreateClusterStmt, error) {
	stmt := &nodes.CreateClusterStmt{
		Loc: nodes.Loc{Start: start},
	}

	// [ IF NOT EXISTS ]: Oracle answers an IF not followed by NOT EXISTS
	// with ORA-11543 at the IF.
	if p.cur.Type == kwIF {
		if p.peekNext().Type != kwNOT {
			return nil, p.syntaxErrorAtCur()
		}
		ifTok := p.cur
		p.advance() // IF
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance() // EXISTS
	}

	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	var err error
	stmt.Name, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	// [ SHARING = { METADATA | NONE } ]
	if p.isKeywordStr("SHARING") {
		p.advance() // SHARING
		if err := p.expectCreateTableToken('='); err != nil {
			return nil, err
		}
		if !p.isKeywordStr("METADATA") && !p.isKeywordStr("NONE") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	if err := p.parseCreateClusterColumns(stmt); err != nil {
		return nil, err
	}

	for !p.isStatementEnd() {
		switch {
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if err := p.parseClusterInteger(&stmt.PctFree); err != nil {
				return nil, err
			}
		case p.isKeywordStr("PCTUSED"):
			p.advance()
			if err := p.parseClusterInteger(&stmt.PctUsed); err != nil {
				return nil, err
			}
		case p.isKeywordStr("INITRANS"):
			p.advance()
			if err := p.parseClusterInteger(&stmt.InitTrans); err != nil {
				return nil, err
			}
		case p.isKeywordStr("MAXTRANS"):
			p.advance()
			if err := p.parseClusterInteger(nil); err != nil {
				return nil, err
			}
		case p.cur.Type == kwSTORAGE:
			if err := p.parseCreateTableStorageClause(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwSIZE:
			p.advance()
			if stmt.Size, err = p.parseClusterSize(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			stmt.Tablespace = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwINDEX:
			p.advance()
			stmt.IsIndex = true
		case p.isKeywordStr("HASHKEYS"):
			p.advance()
			stmt.IsHash = true
			stmt.HashKeys = p.cur.Str
			if err := p.parseCreateTableInteger(); err != nil {
				return nil, err
			}
			if p.cur.Type == kwHASH {
				p.advance() // HASH
				if err := p.expectCreateTableToken(kwIS); err != nil {
					return nil, err
				}
				if stmt.HashExpr, err = p.parseExpr(); err != nil {
					return nil, err
				}
				if stmt.HashExpr == nil {
					return nil, p.syntaxErrorAtCur()
				}
			}
		case p.isKeywordStr("SINGLE"):
			p.advance()
			stmt.SingleTable = true
			if err := p.expectCreateTableToken(kwTABLE); err != nil {
				return nil, err
			}
		case p.isKeywordStr("PARALLEL"):
			p.advance()
			stmt.Parallel = "PARALLEL"
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			}
		case p.isKeywordStr("NOPARALLEL"):
			p.advance()
			stmt.Parallel = "NOPARALLEL"
		case p.isKeywordStr("ROWDEPENDENCIES"):
			p.advance()
			stmt.RowDep = true
		case p.isKeywordStr("NOROWDEPENDENCIES"):
			p.advance()
			stmt.NoRowDep = true
		case p.cur.Type == kwCACHE:
			p.advance()
			stmt.Cache = true
		case p.cur.Type == kwNOCACHE:
			p.advance()
			stmt.NoCache = true
		case p.cur.Type == kwPARTITION:
			if err := p.parseClusterRangePartitions(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("INMEMORY"):
			p.advance()
			return nil, p.syntaxErrorAtCur()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateClusterColumns reads the cluster key.
//
//	( column datatype [ COLLATE column_collation_name ] [ SORT ] [, ...] )
func (p *Parser) parseCreateClusterColumns(stmt *nodes.CreateClusterStmt) error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		col := &nodes.ClusterColumn{Loc: nodes.Loc{Start: p.pos()}}
		if !p.isCreateTableNameStart() {
			return p.syntaxErrorAtCur()
		}
		var err error
		if col.Name, err = p.parseIdentifier(); err != nil {
			return err
		}
		if col.DataType, err = p.parseTypeName(); err != nil {
			return err
		}
		if col.DataType == nil || col.DataType.Names.Len() == 0 {
			return p.syntaxErrorAtCur()
		}
		if p.isKeywordStr("COLLATE") {
			p.advance()
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		}
		if p.isKeywordStr("SORT") {
			p.advance()
			col.Sort = true
		}
		col.Loc.End = p.prev.End
		stmt.Columns = append(stmt.Columns, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ','
	}
	return p.expectCreateTableToken(')')
}

// parseClusterRangePartitions reads cluster_range_partitions, whose
// partitions are those of a range partitioned table.
//
//	PARTITION BY RANGE ( column [, column ]... )
//	( PARTITION [ partition ] range_values_clause [ table_partition_description ] [, ...] )
func (p *Parser) parseClusterRangePartitions() error {
	p.advance() // PARTITION
	if err := p.expectCreateTableToken(kwBY); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("RANGE"); err != nil {
		return err
	}
	if err := p.parseCreateTablePartitionColumns(&nodes.List{}); err != nil {
		return err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		def, err := p.parseCreateTablePartitionDesc(createTableRangeBound, false, createTableHashBound)
		if err != nil {
			return err
		}
		_ = def
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ','
	}
	return p.expectCreateTableToken(')')
}

// parseAlterClusterStmt parses an ALTER CLUSTER statement after CLUSTER.
//
// BNF: oracle/parser/bnf/ALTER-CLUSTER.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-CLUSTER.html
//
//	ALTER CLUSTER [ IF EXISTS ] [ schema . ] cluster
//	    { physical_attributes_clause
//	    | SIZE integer
//	    | MODIFY PARTITION partition allocate_extent_clause
//	    | allocate_extent_clause
//	    | deallocate_unused_clause
//	    | parallel_clause
//	    }
//
//	physical_attributes_clause:
//	    [ PCTUSED integer ] [ PCTFREE integer ] [ INITRANS integer ] [ STORAGE storage_clause ]
//
//	allocate_extent_clause:
//	    ALLOCATE EXTENT
//	    [ ( { SIZE size_clause | DATAFILE 'filename' | INSTANCE integer }... ) ]
//
//	deallocate_unused_clause:
//	    DEALLOCATE UNUSED [ KEEP size_clause ]
//
//	parallel_clause:
//	    { PARALLEL [ integer ] | NOPARALLEL }
//
// Oracle 23ai takes several of these clauses in one statement (ALTER
// CLUSTER c PCTFREE 10 SIZE 1K), MAXTRANS among the physical attributes,
// and CACHE or NOCACHE.
func (p *Parser) parseAlterClusterStmt(start int) (*nodes.AlterClusterStmt, error) {
	stmt := &nodes.AlterClusterStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // IF
			p.advance() // EXISTS
			stmt.IfExists = true
		}
	}
	var err error
	stmt.Name, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	physical := func() {
		if stmt.Action == "" {
			stmt.Action = "PHYSICAL_ATTRS"
		}
	}
	for !p.isStatementEnd() {
		switch {
		case p.cur.Type == kwSIZE:
			p.advance()
			stmt.Action = "SIZE"
			if stmt.Size, err = p.parseClusterSize(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("PCTUSED"):
			p.advance()
			if err := p.parseClusterInteger(&stmt.PctUsed); err != nil {
				return nil, err
			}
			physical()
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if err := p.parseClusterInteger(&stmt.PctFree); err != nil {
				return nil, err
			}
			physical()
		case p.isKeywordStr("INITRANS"):
			p.advance()
			if err := p.parseClusterInteger(&stmt.InitTrans); err != nil {
				return nil, err
			}
			physical()
		case p.isKeywordStr("MAXTRANS"):
			p.advance()
			if err := p.parseClusterInteger(nil); err != nil {
				return nil, err
			}
			physical()
		case p.cur.Type == kwSTORAGE:
			if err := p.parseCreateTableStorageClause(); err != nil {
				return nil, err
			}
			physical()
		case p.cur.Type == kwCACHE, p.cur.Type == kwNOCACHE:
			p.advance()
			physical()
		case p.isKeywordStr("INMEMORY"):
			p.advance()
			return nil, p.syntaxErrorAtCur()
		case p.cur.Type == kwMODIFY:
			p.advance() // MODIFY
			stmt.Action = "MODIFY_PARTITION"
			if err := p.expectCreateTableToken(kwPARTITION); err != nil {
				return nil, err
			}
			stmt.ModifyPartition = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
			if !p.isKeywordStr("ALLOCATE") {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.parseClusterAllocateExtent(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("ALLOCATE"):
			stmt.Action = "ALLOCATE_EXTENT"
			if err := p.parseClusterAllocateExtent(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("DEALLOCATE"):
			p.advance() // DEALLOCATE
			stmt.Action = "DEALLOCATE_UNUSED"
			if err := p.expectCreateTableWord("UNUSED"); err != nil {
				return nil, err
			}
			if p.isKeywordStr("KEEP") {
				p.advance()
				keep, err := p.parseClusterSize()
				if err != nil {
					return nil, err
				}
				_ = keep
			}
		case p.isKeywordStr("PARALLEL"):
			p.advance()
			stmt.Action = "PARALLEL"
			stmt.Parallel = "PARALLEL"
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			}
		case p.isKeywordStr("NOPARALLEL"):
			p.advance()
			stmt.Action = "PARALLEL"
			stmt.Parallel = "NOPARALLEL"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}
	if stmt.Action == "" {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseClusterAllocateExtent reads an allocate_extent_clause, each of
// whose parameters Oracle takes once (ORA-03278 for a second SIZE).
//
//	ALLOCATE EXTENT
//	[ ( { SIZE size_clause | DATAFILE 'filename' | INSTANCE integer }... ) ]
func (p *Parser) parseClusterAllocateExtent() error {
	p.advance() // ALLOCATE
	if err := p.expectCreateTableWord("EXTENT"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance() // '('
	var size, datafile, instance bool
	for first := true; first || p.cur.Type != ')'; first = false {
		switch {
		case !size && p.cur.Type == kwSIZE:
			p.advance()
			extent, err := p.parseClusterSize()
			if err != nil {
				return err
			}
			_ = extent
			size = true
		case !datafile && p.isKeywordStr("DATAFILE"):
			p.advance()
			if err := p.expectCreateTableToken(tokSCONST); err != nil {
				return err
			}
			datafile = true
		case !instance && p.isKeywordStr("INSTANCE"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
			instance = true
		default:
			return p.syntaxErrorAtCur()
		}
	}
	p.advance() // ')'
	return nil
}
