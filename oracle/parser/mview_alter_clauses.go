package parser

import (
	"strconv"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads ALTER MATERIALIZED VIEW by the grammar of its page,
// reusing ALTER TABLE's readers for the clauses the two share.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-VIEW.html
// BNF: oracle/parser/bnf/ALTER-MATERIALIZED-VIEW.bnf

// parseAlterMviewStmt reads an ALTER MATERIALIZED VIEW after its VIEW.
//
//	ALTER MATERIALIZED VIEW [ IF EXISTS ] [ schema. ] materialized_view
//	    [ physical_attributes_clause
//	    | modify_mv_column_clause
//	    | table_compression
//	    | inmemory_table_clause
//	    | modify_LOB_storage_clause [, modify_LOB_storage_clause ]...
//	    | alter_table_partitioning
//	    | parallel_clause
//	    | logging_clause
//	    | allocate_extent_clause
//	    | deallocate_unused_clause
//	    | shrink_clause
//	    | { CACHE | NOCACHE }
//	    | annotations_clause
//	    ]
//	    [ alter_iot_clauses ]
//	    [ USING INDEX physical_attributes_clause ]
//	    [ MODIFY scoped_table_ref_constraint | alter_mv_refresh ]
//	    [ evaluation_edition_clause ]
//	    [ { ENABLE | DISABLE } ON QUERY COMPUTATION ]
//	    [ alter_query_rewrite_clause | COMPILE | CONSIDER FRESH ]
//
// The parts stand in that order, as the 19c page has them, and at least
// one stands. The pinned server refuses the 26 page's LOB_storage_clause
// (ORA-14051 at LOB) and its QUERY REWRITE { ENABLE | DISABLE } (ORA-02000
// at ENABLE): it takes { ENABLE | DISABLE } QUERY REWRITE.
func (p *Parser) parseAlterMviewStmt(start int) (*nodes.AlterMaterializedViewStmt, error) {
	stmt := &nodes.AlterMaterializedViewStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}
	var err error
	stmt.Name, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if err := p.parseAlterMviewProperty(stmt); err != nil {
		return nil, err
	}
	if p.isAlterMviewIOTStart() {
		setAlterMviewAction(stmt, "IOT")
		if err := p.parseAlterMviewIOT(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == kwUSING && p.peekNext().Type == kwINDEX {
		p.advance() // consume USING
		p.advance() // consume INDEX
		setAlterMviewAction(stmt, "USING_INDEX")
		if err := p.parseAlterMviewPhysicalAttrs(stmt); err != nil {
			return nil, err
		}
	}
	switch {
	case p.cur.Type == kwMODIFY && p.atNextWord("SCOPE"):
		p.advance() // consume MODIFY
		setAlterMviewAction(stmt, "MODIFY_SCOPE")
		if err := p.parseAlterMviewScopeFor(stmt); err != nil {
			return nil, err
		}
	case p.cur.Type == kwREFRESH:
		setAlterMviewAction(stmt, "REFRESH")
		if err := p.parseAlterMviewRefresh(stmt); err != nil {
			return nil, err
		}
	}
	if p.atWord("EVALUATE") {
		setAlterMviewAction(stmt, "EVALUATE_USING_EDITION")
		if err := p.parseAlterMviewEvaluation(stmt); err != nil {
			return nil, err
		}
	}
	if (p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE) && p.peekNext().Type == kwON {
		enable := p.cur.Type == kwENABLE
		p.advance() // consume ENABLE / DISABLE
		p.advance() // consume ON
		if err := p.atExpectWord("QUERY"); err != nil {
			return nil, err
		}
		if err := p.atExpectWord("COMPUTATION"); err != nil {
			return nil, err
		}
		if enable {
			setAlterMviewAction(stmt, "ENABLE_ON_QUERY_COMPUTATION")
			stmt.EnableOnQueryComputation = true
		} else {
			setAlterMviewAction(stmt, "DISABLE_ON_QUERY_COMPUTATION")
			stmt.DisableOnQueryComputation = true
		}
	}
	switch {
	case p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
		enable := p.cur.Type == kwENABLE
		p.advance() // consume ENABLE / DISABLE
		if err := p.atExpectWord("QUERY"); err != nil {
			return nil, err
		}
		if err := p.atExpect(kwREWRITE); err != nil {
			return nil, err
		}
		if enable {
			setAlterMviewAction(stmt, "ENABLE_QUERY_REWRITE")
			if err := p.parseMviewUnusableEditions(); err != nil {
				return nil, err
			}
		} else {
			setAlterMviewAction(stmt, "DISABLE_QUERY_REWRITE")
		}
	case p.atWord("QUERY"):
		// The pinned server reads the 26 page's QUERY REWRITE { ENABLE |
		// DISABLE } and refuses whatever follows, the end of the statement
		// included (ORA-02000 at the end of QUERY REWRITE ENABLE).
		p.advance() // consume QUERY
		if err := p.atExpect(kwREWRITE); err != nil {
			return nil, err
		}
		if p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
			p.advance()
		}
		return nil, p.syntaxErrorAtCur()
	case p.atWord("COMPILE"):
		p.advance()
		setAlterMviewAction(stmt, "COMPILE")
	case p.atWord("CONSIDER"):
		p.advance() // consume CONSIDER
		if err := p.atExpectWord("FRESH"); err != nil {
			return nil, err
		}
		setAlterMviewAction(stmt, "CONSIDER_FRESH")
	}
	if stmt.Action == "" {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// setAlterMviewAction records the first clause an ALTER MATERIALIZED VIEW
// holds as its Action.
func setAlterMviewAction(stmt *nodes.AlterMaterializedViewStmt, action string) {
	if stmt.Action == "" {
		stmt.Action = action
	}
}

// parseAlterMviewProperty reads the property clause an ALTER MATERIALIZED
// VIEW may begin with, if one starts here.
func (p *Parser) parseAlterMviewProperty(stmt *nodes.AlterMaterializedViewStmt) error {
	switch {
	case p.cur.Type == kwMODIFY && p.peekNext().Type == '(':
		p.advance() // consume MODIFY
		setAlterMviewAction(stmt, "MODIFY")
		return p.parseAlterMviewModifyColumns()
	case p.cur.Type == kwMODIFY && !p.atNextWord("SCOPE") && !p.isAlterTablePartitioningStart() && !p.atNextWord("LOB"):
		// MODIFY leads a column list, LOB storage, a partition operation or
		// a scope; Oracle refuses anything else at the word after it.
		p.advance() // consume MODIFY
		return p.syntaxErrorAtCur()
	case p.cur.Type == kwMODIFY && p.atNextWord("LOB"):
		setAlterMviewAction(stmt, "MODIFY_LOB")
		for {
			p.advance() // consume MODIFY
			if err := p.parseATLobModify(); err != nil {
				return err
			}
			if p.cur.Type != ',' {
				return nil
			}
			p.advance() // consume ','
			if p.cur.Type != kwMODIFY {
				return p.syntaxErrorAtCur()
			}
		}
	case p.cur.Type == kwADD && p.peekNext().Type != kwPARTITION && !p.atNextWord("OVERFLOW"):
		// ADD leads ADD PARTITION or ADD OVERFLOW; Oracle refuses anything
		// else at the word after it.
		p.advance() // consume ADD
		return p.syntaxErrorAtCur()
	case p.isAlterTablePartitioningStart():
		setAlterMviewAction(stmt, "PARTITION")
		cmd, err := p.parseAlterTablePartitioning()
		_ = cmd // the AST keeps no partition operation of a view
		return err
	case p.isATPhysicalAttrStart():
		return p.parseAlterMviewPhysicalAttrs(stmt)
	case p.isCreateTablePartitionCompressionStart():
		setAlterMviewAction(stmt, p.cur.Str)
		return p.parseCreateTablePartitionCompression()
	case p.isATInmemoryStart():
		if p.atWord("NO") {
			setAlterMviewAction(stmt, "NO_INMEMORY")
		} else {
			setAlterMviewAction(stmt, "INMEMORY")
		}
		cmd, err := p.parseAlterTableInmemory()
		_ = cmd // the AST keeps no in-memory clause
		return err
	case p.cur.Type == kwPARALLEL:
		setAlterMviewAction(stmt, "PARALLEL")
		p.advance()
		if p.cur.Type == tokICONST {
			stmt.ParallelDegree = p.cur.Str
			p.advance()
		}
	case p.cur.Type == kwNOPARALLEL:
		setAlterMviewAction(stmt, "NOPARALLEL")
		p.advance()
	case p.cur.Type == kwLOGGING, p.cur.Type == kwNOLOGGING:
		setAlterMviewAction(stmt, p.cur.Str)
		p.advance()
	case p.isATAllocateExtentStart():
		setAlterMviewAction(stmt, "ALLOCATE_EXTENT")
		return p.parseATAllocateExtent()
	case p.isATDeallocateStart():
		setAlterMviewAction(stmt, "DEALLOCATE_UNUSED")
		return p.parseATDeallocateUnused()
	case p.isATShrinkStart():
		// SHRINK SPACE [ COMPACT ] [ CASCADE ]
		setAlterMviewAction(stmt, "SHRINK")
		p.advance() // consume SHRINK
		if err := p.atExpectWord("SPACE"); err != nil {
			return err
		}
		if p.atWord("COMPACT") {
			p.advance()
			stmt.Compact = true
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
			stmt.Cascade = true
		}
	case p.cur.Type == kwCACHE, p.cur.Type == kwNOCACHE:
		setAlterMviewAction(stmt, p.cur.Str)
		p.advance()
	case p.atWord("ANNOTATIONS"):
		setAlterMviewAction(stmt, "ANNOTATIONS")
		p.advance() // consume ANNOTATIONS
		return p.parseATAnnotations()
	}
	return nil
}

// parseAlterMviewModifyColumns reads modify_mv_column_clause after MODIFY.
//
//	( column [ ENCRYPT encryption_spec | DECRYPT ] [, ...] )
func (p *Parser) parseAlterMviewModifyColumns() error {
	p.advance() // consume '('
	for {
		if err := p.parseATName(); err != nil {
			return err
		}
		switch {
		case p.atWord("ENCRYPT"):
			p.advance()
			if err := p.parseATEncryptionSpec(); err != nil {
				return err
			}
		case p.atWord("DECRYPT"):
			p.advance()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.atExpect(')')
}

// parseAlterMviewPhysicalAttrs reads physical_attributes_clause, one item
// or more: PCTFREE, PCTUSED, INITRANS integer or storage_clause.
func (p *Parser) parseAlterMviewPhysicalAttrs(stmt *nodes.AlterMaterializedViewStmt) error {
	if !p.isATPhysicalAttrStart() {
		return p.syntaxErrorAtCur()
	}
	for p.isATPhysicalAttrStart() {
		word, value := p.cur.Str, p.peekNext()
		setAlterMviewAction(stmt, word)
		if err := p.parseATPhysicalAttr(); err != nil {
			return err
		}
		n, err := strconv.Atoi(value.Str)
		if err != nil || value.Type != tokICONST {
			continue
		}
		switch word {
		case "PCTFREE":
			stmt.PctFree = n
		case "PCTUSED":
			stmt.PctUsed = n
		case "INITRANS":
			stmt.IniTrans = n
		}
	}
	return nil
}

// isAlterMviewIOTStart reports whether an alter_iot_clauses starts here.
func (p *Parser) isAlterMviewIOTStart() bool {
	switch {
	case p.atWord("OVERFLOW"), p.atWord("MAPPING"), p.atWord("NOMAPPING"), p.atWord("INCLUDING"),
		p.atWord("PCTTHRESHOLD"), p.atWord("COALESCE"):
		return true
	case p.cur.Type == kwADD:
		return p.atNextWord("OVERFLOW")
	}
	return false
}

// parseAlterMviewIOT reads alter_iot_clauses.
//
//	{ index_org_table_clause | alter_overflow_clause
//	| alter_mapping_table_clauses | COALESCE }
//	index_org_table_clause:
//	    { mapping_table_clauses | PCTTHRESHOLD integer | prefix_compression
//	    | index_org_overflow_clause }
//	alter_overflow_clause:
//	    { OVERFLOW { allocate_extent_clause | deallocate_unused_clause
//	               | segment_attributes_clause }
//	    | ADD OVERFLOW [ segment_attributes_clause ] }
func (p *Parser) parseAlterMviewIOT() error {
	switch {
	case p.atWord("PCTTHRESHOLD"):
		p.advance()
		return p.parseATInteger()
	case p.cur.Type == kwADD:
		p.advance() // consume ADD
		p.advance() // consume OVERFLOW
		if p.isATSegmentAttrStart() {
			return p.parseATSegmentAttrs(nil)
		}
		return nil
	}
	cmd, err := p.parseAlterTableIOT()
	_ = cmd // the AST keeps no index-organized clause of a view
	return err
}

// parseAlterMviewScopeFor reads scoped_table_ref_constraint after MODIFY.
//
//	SCOPE FOR ( ref_column ) IS [ schema. ] scope_table
func (p *Parser) parseAlterMviewScopeFor(stmt *nodes.AlterMaterializedViewStmt) error {
	p.advance() // consume SCOPE
	if err := p.atExpect(kwFOR); err != nil {
		return err
	}
	if err := p.atExpect('('); err != nil {
		return err
	}
	var err error
	stmt.ScopeColumn, err = p.parseATNameString()
	if err != nil {
		return err
	}
	if err := p.atExpect(')'); err != nil {
		return err
	}
	if err := p.atExpect(kwIS); err != nil {
		return err
	}
	stmt.ScopeTable, err = p.parseRequiredObjectName()
	return err
}

// parseAlterMviewRefresh reads alter_mv_refresh.
//
//	REFRESH
//	    { { FAST | COMPLETE | FORCE }
//	    | ON { DEMAND | COMMIT }
//	    | { START WITH | NEXT } date
//	    | WITH PRIMARY KEY
//	    | USING { DEFAULT [ MASTER | LOCAL ] ROLLBACK SEGMENT
//	            | [ MASTER | LOCAL ] ROLLBACK SEGMENT rollback_segment }
//	    | USING { ENFORCED | TRUSTED } CONSTRAINTS
//	    | { ENABLE | DISABLE } CONCURRENT REFRESH
//	    } ...
//
// The 19c page's form, with the 26 page's concurrent refresh: the parts in
// any order, at least one.
func (p *Parser) parseAlterMviewRefresh(stmt *nodes.AlterMaterializedViewStmt) error {
	p.advance() // consume REFRESH
	var on, start, next, with, constraints, concurrent bool
	for read := 0; ; read++ {
		switch {
		case p.atWord("FAST"), p.atWord("COMPLETE"), p.cur.Type == kwFORCE:
			stmt.RefreshMethod = p.cur.Str
			p.advance()
		case !on && p.cur.Type == kwON:
			p.advance() // consume ON
			switch {
			case p.cur.Type == kwCOMMIT:
				stmt.RefreshMode = "ON_COMMIT"
			case p.atWord("DEMAND"):
				stmt.RefreshMode = "ON_DEMAND"
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
			on = true
		case !start && p.cur.Type == kwSTART:
			p.advance() // consume START
			if err := p.atExpect(kwWITH); err != nil {
				return err
			}
			expr, err := p.parseMviewDateExpr()
			if err != nil {
				return err
			}
			stmt.StartWith = expr
			start = true
		case !next && p.cur.Type == kwNEXT:
			p.advance() // consume NEXT
			expr, err := p.parseMviewDateExpr()
			if err != nil {
				return err
			}
			stmt.Next = expr
			next = true
		case !with && p.cur.Type == kwWITH:
			p.advance() // consume WITH
			if err := p.atExpect(kwPRIMARY); err != nil {
				return err
			}
			if err := p.atExpect(kwKEY); err != nil {
				return err
			}
			stmt.WithPrimaryKey = true
			with = true
		case p.cur.Type == kwUSING:
			p.advance() // consume USING
			if p.atWord("TRUSTED") || p.atWord("ENFORCED") {
				if constraints {
					return p.syntaxErrorAtCur()
				}
				stmt.UsingConstraints = p.cur.Str
				p.advance() // consume TRUSTED / ENFORCED
				if err := p.atExpect(kwCONSTRAINTS); err != nil {
					return err
				}
				constraints = true
				break
			}
			if p.cur.Type != kwDEFAULT {
				if p.atWord("MASTER") || p.atWord("LOCAL") {
					p.advance()
				}
				if err := p.atExpect(kwROLLBACK); err != nil {
					return err
				}
				if err := p.atExpectWord("SEGMENT"); err != nil {
					return err
				}
				var err error
				stmt.UsingRollbackSegment, err = p.parseATNameString()
				if err != nil {
					return err
				}
				break
			}
			if err := p.parseMviewRollbackSegment(); err != nil {
				return err
			}
		case !concurrent && (p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE) && p.atNextWord("CONCURRENT"):
			p.advance() // consume ENABLE / DISABLE
			p.advance() // consume CONCURRENT
			if err := p.atExpect(kwREFRESH); err != nil {
				return err
			}
			concurrent = true
		default:
			if read == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseAlterMviewEvaluation reads ALTER's evaluation_edition_clause.
//
//	EVALUATE USING { CURRENT EDITION | EDITION edition | NULL EDITION }
func (p *Parser) parseAlterMviewEvaluation(stmt *nodes.AlterMaterializedViewStmt) error {
	p.advance() // consume EVALUATE
	if err := p.atExpect(kwUSING); err != nil {
		return err
	}
	if p.atWord("EDITION") {
		p.advance() // consume EDITION
		var err error
		stmt.EditionName, err = p.parseATNameString()
		return err
	}
	return p.parseMviewEditionRef(true)
}
