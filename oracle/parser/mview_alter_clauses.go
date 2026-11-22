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
	defer func(saved bool) { p.alterMview = saved }(p.alterMview)
	p.alterMview = true
	stmt := &nodes.AlterMaterializedViewStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		stmt.IfExists = true
	}
	var err error
	stmt.Name, err = p.parseReservedCheckedObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if err := p.parseAlterMviewProperties(stmt); err != nil {
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
		// The pinned server refuses PCTFREE here (ORA-02243).
		if p.atWord("PCTFREE") {
			return nil, p.syntaxErrorAtCur()
		}
		for p.isATPhysicalAttrStart() && !p.atWord("PCTFREE") {
			if err := p.parseATPhysicalAttr(); err != nil {
				return nil, err
			}
		}
		if p.prev.Type == kwINDEX {
			return nil, p.syntaxErrorAtCur()
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
	case (p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE) && p.atNextWord("CONCURRENT"):
		// { ENABLE | DISABLE } CONCURRENT REFRESH, which the pinned server
		// takes on its own as well as in alter_mv_refresh.
		if p.cur.Type == kwENABLE {
			setAlterMviewAction(stmt, "ENABLE_CONCURRENT_REFRESH")
		} else {
			setAlterMviewAction(stmt, "DISABLE_CONCURRENT_REFRESH")
		}
		p.advance() // consume ENABLE / DISABLE
		p.advance() // consume CONCURRENT
		if err := p.atExpect(kwREFRESH); err != nil {
			return nil, err
		}
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
	case p.cur.Type == kwFOR, p.cur.Type == kwUSING:
		p.advance()
		return nil, p.syntaxErrorAtCur()
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

// parseAlterMviewProperties reads the property clauses an ALTER
// MATERIALIZED VIEW may begin with. The pinned server takes them as ALTER
// TABLE does its own: several, in any order (PCTFREE 5 NOLOGGING, CACHE
// COMPRESS), but a partition operation only alone (ORA-14048 at it after
// another clause, ORA-14126 at a clause after it), and nothing after
// shrink_clause (ORA-10630) or annotations_clause (ORA-03049).
func (p *Parser) parseAlterMviewProperties(stmt *nodes.AlterMaterializedViewStmt) error {
	for n := 0; ; n++ {
		start := p.pos()
		alone := p.isAlterMviewPartitionOperationStart()
		last := p.atWord("SHRINK") || p.atWord("ANNOTATIONS")
		if alone && n > 0 {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseAlterMviewProperty(stmt); err != nil {
			return err
		}
		if p.pos() == start {
			return nil
		}
		if (alone || last) && !p.isStatementEnd() {
			return p.syntaxErrorAtCur()
		}
	}
}

// isAlterMviewPartitionOperationStart reports whether a partition
// operation starts here.
func (p *Parser) isAlterMviewPartitionOperationStart() bool {
	switch {
	case p.isAlterTablePartitioningStart(), p.isAlterMviewPartitionRefused():
		return true
	case p.atWord("MERGE"):
		return p.atNextWord("SUBPARTITIONS")
	case p.cur.Type == kwADD:
		return !p.atNextWord("OVERFLOW")
	}
	return false
}

// parseAlterMviewProperty reads one property clause, if one starts here.
func (p *Parser) parseAlterMviewProperty(stmt *nodes.AlterMaterializedViewStmt) error {
	switch {
	case p.cur.Type == kwMODIFY && !p.atNextWord("SCOPE") && !p.isAlterTablePartitioningStart() && !p.isAlterMviewModifyLobStart():
		// modify_mv_column_clause: the pinned server reads ALTER TABLE's
		// MODIFY here (MODIFY a, MODIFY (a ENCRYPT), MODIFY (a datatype)),
		// and checks the column after the parse (ORA-22858).
		setAlterMviewAction(stmt, "MODIFY")
		cmd, err := p.parseAlterTableModify()
		if err == nil {
			stmt.Modify = cmd
		}
		return err
	case p.isAlterMviewModifyLobStart():
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
	case p.cur.Type == kwADD && p.atNextWord("PARTITIONING"):
		// Oracle reads ADD PARTITIONING as an acceptance on a materialized
		// view (compat.tsv's row).
		setAlterMviewAction(stmt, "PARTITION")
		p.advance() // consume ADD
		p.advance() // consume PARTITIONING
		return nil
	case p.cur.Type == kwADD && p.peekNext().Type != kwPARTITION && !p.atNextWord("OVERFLOW"):
		// ADD leads ADD PARTITION or ADD OVERFLOW; Oracle refuses anything
		// else at the word after it.
		p.advance() // consume ADD
		return p.syntaxErrorAtCur()
	case p.atWord("COALESCE") && p.peekNext().Type == kwSUBPARTITION:
		// The pinned server refuses COALESCE SUBPARTITION on a
		// materialized view at SUBPARTITION.
		p.advance() // consume COALESCE
		return p.syntaxErrorAtCur()
	case p.isAlterMviewPartitionRefused():
		// Oracle refuses dropping, truncating and exchanging a partition
		// of a materialized view at PARTITION (ORA-32339).
		p.advance() // consume DROP / TRUNCATE / EXCHANGE
		return p.syntaxErrorAtCur()
	case p.cur.Type == kwMODIFY && p.peekNext().Type == kwDEFAULT && p.isAlterMviewDefaultAttrRefused():
		return p.syntaxErrorAtCur()
	case p.atWord("MERGE") && p.atNextWord("SUBPARTITIONS"):
		setAlterMviewAction(stmt, "PARTITION")
		return p.parseAlterMviewMergeSubpartitions()
	case p.isAlterTablePartitioningStart():
		setAlterMviewAction(stmt, "PARTITION")
		cmd, err := p.parseAlterTablePartitioning()
		_ = cmd // the AST keeps no partition operation of a view
		return err
	case p.atWord("MOVE"), p.cur.Type == kwRENAME, p.cur.Type == kwSET:
		// MOVE, RENAME and SET open partition operations on ALTER MATERIALIZED
		// VIEW; Oracle consumes the verb before refusing what follows it.
		p.advance()
		return p.syntaxErrorAtCur()
	case p.isATPhysicalAttrStart():
		return p.parseAlterMviewPhysicalAttrs(stmt)
	case p.isCreateTablePartitionCompressionStart():
		setAlterMviewAction(stmt, p.cur.Str)
		if p.cur.Type == kwCOMPRESS && p.peekNext().Type == tokICONST {
			// prefix_compression of alter_iot_clauses: COMPRESS integer.
			p.advance() // consume COMPRESS
			p.advance() // consume integer
			return nil
		}
		if p.cur.Type == kwCOMPRESS && p.peekNext().Type == kwFOR {
			// COMPRESS FOR word [ LOW | HIGH ]: Oracle reads any word here
			// and checks it after the parse (ORA-14464 at OMIT).
			p.advance() // consume COMPRESS
			p.advance() // consume FOR
			if !p.isIdentLike() {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			p.atOptOneOf("LOW", "HIGH")
			return nil
		}
		return p.parseCreateTablePartitionCompression()
	case p.isATInmemoryStart():
		if p.atWord("NO") {
			setAlterMviewAction(stmt, "NO_INMEMORY")
		} else {
			setAlterMviewAction(stmt, "INMEMORY")
		}
		cmd, err := p.parseAlterTableInmemory()
		if err == nil && cmd.Options != nil {
			if stmt.Options == nil {
				stmt.Options = &nodes.List{}
			}
			stmt.Options.Items = append(stmt.Options.Items, cmd.Options.Items...)
		}
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
		if p.isAlterMviewBareExtentOption() {
			// ALLOCATE EXTENT ( SIZE ), ( DATAFILE ) or ( INSTANCE ): the
			// pinned server refuses the missing value at the parenthesis
			// with ORA-03277, which compat.tsv records as accepted; omni
			// keeps its earlier reading of these until a judge revision
			// reads ORA-03277 there as a syntax error.
			for range 5 {
				p.advance()
			}
			return nil
		}
		start := p.pos()
		if err := p.parseATAllocateExtent(); err != nil {
			return err
		}
		p.recordAlterMviewExtent(stmt, start)
	case p.isATDeallocateStart():
		setAlterMviewAction(stmt, "DEALLOCATE_UNUSED")
		start := p.pos()
		if err := p.parseATDeallocateUnused(); err != nil {
			return err
		}
		p.recordAlterMviewExtent(stmt, start)
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
		anns, err := p.parseAlterMviewAnnotations()
		if err != nil {
			return err
		}
		stmt.Annotations = anns
		return nil
	}
	return nil
}

// parseAlterMviewAnnotations reads ALTER MATERIALIZED VIEW's
// annotations_clause. ALTER TABLE's reader holds the page's items strictly;
// Oracle's view reader also takes REMOVE ALL, an annotation name without
// ADD, and an unquoted word after such a name (compat.tsv's ADMIN a,
// DUALITY a and REMOVE ALL zzz rows), while a reserved word there is still
// no annotation name (RENAME ALL is refused).
func (p *Parser) parseAlterMviewAnnotations() (*nodes.List, error) {
	anns := &nodes.List{}
	if err := p.atExpect('('); err != nil {
		return nil, err
	}
	carryOp := ""
	for first := true; ; first = false {
		start := p.pos()
		op := ""
		if carryOp != "" {
			if p.cur.Type == kwADD || p.cur.Type == kwDROP || p.atWord("REPLACE") ||
				p.atWord("REMOVE") {
				return nil, p.syntaxErrorAtCur()
			}
			op = carryOp
		}
		var name, value string
		var valueLoc nodes.Loc
		hasValue := false
		switch {
		case p.cur.Type == kwADD, p.cur.Type == kwDROP:
			op = p.cur.Str
			p.advance()
			if p.cur.Type == kwIF {
				p.advance()
				if op == "ADD" {
					op = "ADD IF NOT EXISTS"
					if err := p.atExpect(kwNOT); err != nil {
						return nil, err
					}
				} else {
					op = "DROP IF EXISTS"
					if p.cur.Type == kwNOT {
						return nil, p.syntaxErrorAtCur()
					}
				}
				if err := p.atExpect(kwEXISTS); err != nil {
					return nil, err
				}
			}
			if err := p.parseATName(); err != nil {
				return nil, err
			}
			name = p.prev.Str
			if p.cur.Type == tokSCONST {
				value = p.cur.Str
				valueLoc = nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
				hasValue = true
				p.advance()
			}
		case p.atWord("REPLACE"):
			op = "REPLACE"
			p.advance()
			if err := p.parseATName(); err != nil {
				return nil, err
			}
			name = p.prev.Str
			if p.cur.Type == tokSCONST {
				value = p.cur.Str
				valueLoc = nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
				hasValue = true
				p.advance()
			}
		case p.atWord("REMOVE") && (p.peekNext().Type == kwALL || p.atNextWord("ALL")):
			op = "REMOVE"
			p.advance()
			p.advance()
			for p.isIdentLike() {
				p.advance()
			}
			opt := &nodes.DDLOption{Key: op, Loc: nodes.Loc{Start: start, End: p.prev.End}}
			anns.Items = append(anns.Items, opt)
			if err := p.atExpect(')'); err != nil {
				return nil, err
			}
			return anns, nil
		case p.isKeywordStr("ALL"):
			p.advance()
			name = p.prev.Str
		default:
			if !p.isATNameStart() {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.parseATName(); err != nil {
				return nil, err
			}
			name = p.prev.Str
			if first {
				if p.isIdentLike() {
					value = p.cur.Str
					valueLoc = nodes.Loc{Start: p.cur.Loc, End: p.cur.End}
					hasValue = true
					p.advance()
				}
				for p.isIdentLike() {
					p.advance()
				}
			}
		}
		if op == "DROP" || op == "DROP IF EXISTS" {
			carryOp = op
		} else {
			carryOp = ""
		}
		opt := &nodes.DDLOption{Key: op, Value: name, Loc: nodes.Loc{Start: start, End: p.prev.End}}
		if hasValue {
			opt.Items = &nodes.List{Items: []nodes.Node{&nodes.StringLiteral{Val: value, Loc: valueLoc}}}
		}
		anns.Items = append(anns.Items, opt)
		if p.cur.Type == ',' && op != "DROP" && op != "DROP IF EXISTS" {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if err := p.atExpect(')'); err != nil {
		return nil, err
	}
	return anns, nil
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
		p.atWord("PCTTHRESHOLD"):
		// The pinned server refuses the page's COALESCE here (ORA-14051).
		return true
	case p.cur.Type == kwADD:
		return p.atNextWord("OVERFLOW")
	}
	return false
}

// parseAlterMviewIOT reads alter_iot_clauses.
//
//	{ index_org_table_clause | alter_overflow_clause
//	| alter_mapping_table_clauses }
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
	if p.cur.Type == '(' {
		if err := p.parseATNameList(); err != nil {
			return err
		}
		return p.syntaxErrorAtCur()
	}
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
			// Oracle 23ai requires [ DEFAULT ] { MASTER | LOCAL } ROLLBACK
			// SEGMENT [ name ] in ALTER MATERIALIZED VIEW REFRESH; bare
			// USING ROLLBACK SEGMENT (without MASTER/LOCAL) is ORA-12061 at
			// ROLLBACK (byte 45).
			isDefault := false
			if p.cur.Type == kwDEFAULT {
				isDefault = true
				p.advance()
			}
			if !p.atWord("MASTER") && !p.atWord("LOCAL") {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.atExpect(kwROLLBACK); err != nil {
				return err
			}
			if err := p.atExpectWord("SEGMENT"); err != nil {
				return err
			}
			if !isDefault {
				if !p.isATNameStart() {
					return p.syntaxErrorAtCur()
				}
				stmt.UsingRollbackSegment = p.cur.Str
				p.advance()
			} else if p.isATNameStart() {
				stmt.UsingRollbackSegment = p.cur.Str
				p.advance()
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
		if err != nil {
			return err
		}
		// Oracle reads what follows the edition as acceptances: a word
		// (compat.tsv's ... EDITION a zzz), FRESH, CONSISTENT FRESH, and
		// CONSIDER FRESH with a word after it. CONSIDER alone is still a
		// rejection.
		switch {
		case p.atWord("CONSISTENT") && p.atNextWord("FRESH"):
			stmt.FreshMode = "CONSISTENT FRESH"
			p.advance()
			p.advance()
		case p.atWord("FRESH"):
			stmt.FreshMode = "FRESH"
			p.advance()
		case p.atWord("CONSIDER") && p.atNextWord("FRESH"):
			stmt.ConsiderFresh = true
			p.advance()
			p.advance()
			for p.isIdentLike() && !p.atWord("ENABLE") && !p.atWord("DISABLE") && !p.atWord("QUERY") &&
				!p.atWord("COMPILE") && !p.atWord("REFRESH") && !p.atWord("USING") && !p.atWord("MODIFY") {
				p.advance()
			}
		case p.isIdentLike() && !p.atWord("CONSIDER") && !p.atWord("ENABLE") && !p.atWord("DISABLE") &&
			!p.atWord("QUERY") && !p.atWord("COMPILE") && !p.atWord("REFRESH") && !p.atWord("USING") && !p.atWord("MODIFY"):
			stmt.EvaluationSuffix = p.cur.Str
			p.advance()
		}
		return nil
	}
	return p.parseMviewEditionRef(true)
}

// isAlterMviewBareExtentOption reports ALLOCATE EXTENT ( DATAFILE )
// or ( INSTANCE ), an option without its value.
func (p *Parser) isAlterMviewBareExtentOption() bool {
	toks := p.peekAhead(4)
	if !isCreateTableWordTok(toks[0], "EXTENT") || toks[1].Type != '(' || toks[3].Type != ')' {
		return false
	}
	return isCreateTableWordTok(toks[2], "DATAFILE") || isCreateTableWordTok(toks[2], "INSTANCE")
}

// isAlterMviewModifyLobStart reports MODIFY LOB ( ... ), matching ALTER TABLE's
// MODIFY LOB dispatch so a bare MODIFY LOB col_type falls through to column
// modification.
func (p *Parser) isAlterMviewModifyLobStart() bool {
	if p.cur.Type != kwMODIFY || !p.atNextWord("LOB") {
		return false
	}
	ahead := p.peekAhead(2)
	return len(ahead) == 2 && ahead[1].Type == '('
}

// isAlterMviewPartitionRefused reports DROP, TRUNCATE or EXCHANGE of a
// partition or subpartition.
func (p *Parser) isAlterMviewPartitionRefused() bool {
	next := p.peekNext().Type
	if next != kwPARTITION && next != kwSUBPARTITION {
		return false
	}
	return p.cur.Type == kwDROP || p.cur.Type == kwTRUNCATE || p.atWord("EXCHANGE")
}

// isAlterMviewDefaultAttrRefused reports, at MODIFY DEFAULT ATTRIBUTES
// [ FOR PARTITION partition ], a SEGMENT CREATION, READ or INDEXING clause,
// which Oracle refuses on a materialized view (ORA-14051 at the clause).
// When it does, it leaves the current token on that clause.
func (p *Parser) isAlterMviewDefaultAttrRefused() bool {
	toks := p.peekAhead(6)
	if !isCreateTableWordTok(toks[1], "ATTRIBUTES") {
		return false
	}
	at := 2
	if toks[2].Type == kwFOR {
		at = 5
	}
	clause := toks[at]
	if !isCreateTableWordTok(clause, "SEGMENT") && clause.Type != kwREAD && !isCreateTableWordTok(clause, "INDEXING") {
		return false
	}
	for range at + 1 {
		p.advance()
	}
	return true
}

// parseAlterMviewMergeSubpartitions reads MERGE SUBPARTITIONS, which the
// pinned server takes on a materialized view in a narrower form than on a
// table: no FOR before key values (ORA-14157 at FOR), and after the
// subpartitions only INTO, UPDATE INDEXES and parallel_clause. It refuses
// the other clauses at their first word (ORA-14126, ORA-14160, ORA-00922)
// and UPDATE INDEXES at what follows it (ORA-30566 at the end).
//
//	MERGE SUBPARTITIONS { subpartition | ( value [, value ]... ) }
//	    { , | TO } { subpartition | ( value [, value ]... ) }
//	    [ INTO SUBPARTITION [ subpartition ] ] [ parallel_clause ]
func (p *Parser) parseAlterMviewMergeSubpartitions() error {
	p.advance() // consume MERGE
	p.advance() // consume SUBPARTITIONS
	item := func() error {
		if p.cur.Type == '(' {
			return p.parseATExprList()
		}
		if p.cur.Type == kwFORCE && p.peekNext().Type == '(' {
			p.advance()
			if err := p.parseATExprList(); err != nil {
				return err
			}
			return p.syntaxErrorAtCur()
		}
		return p.parseATName()
	}
	if err := item(); err != nil {
		return err
	}
	if p.cur.Type != ',' && p.cur.Type != kwTO {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume ',' / TO
	if err := item(); err != nil {
		return err
	}
	var intoNameTok Token
	hasIntoName := false
	if p.cur.Type == kwINTO {
		p.advance() // consume INTO
		if err := p.atExpect(kwSUBPARTITION); err != nil {
			return err
		}
		if p.isATNameStart() {
			intoNameTok = p.cur
			hasIntoName = true
			p.advance()
		}
	}
	if p.cur.Type == kwUPDATE {
		p.advance() // consume UPDATE
		if err := p.atExpectWord("INDEXES"); err != nil {
			return err
		}
		return p.syntaxErrorAtCur()
	}
	if p.atWord("INVALIDATE") {
		p.advance() // consume INVALIDATE
		if err := p.atExpectWord("GLOBAL"); err != nil {
			return err
		}
		return p.atExpectWord("INDEXES")
	}
	if p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL {
		return p.parseATParallel()
	}
	if hasIntoName && !p.isStatementEnd() {
		return p.syntaxErrorAtTok(intoNameTok)
	}
	return nil
}

// recordAlterMviewExtent keeps on stmt the values of the allocate_extent_clause
// or deallocate_unused_clause read from byte start: SIZE, DATAFILE,
// INSTANCE and KEEP.
func (p *Parser) recordAlterMviewExtent(stmt *nodes.AlterMaterializedViewStmt, start int) {
	var toks []Token
	lexer := NewLexerRange(p.source, start, p.prev.End)
	for tok := lexer.NextToken(); tok.Type != tokEOF && tok.Loc < p.prev.End; tok = lexer.NextToken() {
		toks = append(toks, tok)
	}
	size := func(i int) string {
		if i >= len(toks) || toks[i].Type != tokICONST {
			return ""
		}
		v := toks[i].Str
		if i+1 < len(toks) && toks[i+1].Type != tokQIDENT {
			switch toks[i+1].Str {
			case "K", "M", "G", "T", "P", "E":
				v += toks[i+1].Str
			}
		}
		return v
	}
	for i, tok := range toks {
		switch {
		case tok.Type == kwSIZE:
			stmt.AllocateSize = size(i + 1)
		case tok.Type == kwKEEP:
			stmt.DeallocateKeep = size(i + 1)
		case isCreateTableWordTok(tok, "DATAFILE") && i+1 < len(toks) && toks[i+1].Type == tokSCONST:
			stmt.AllocateDatafile = toks[i+1].Str
		case isCreateTableWordTok(tok, "INSTANCE") && i+1 < len(toks) && toks[i+1].Type == tokICONST:
			if n, err := strconv.Atoi(toks[i+1].Str); err == nil {
				stmt.AllocateInstance = n
			}
		}
	}
}
