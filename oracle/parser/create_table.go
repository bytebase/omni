package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateStmt dispatches CREATE statements based on the object type that follows.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLE.html
//
//	CREATE [ OR REPLACE ] [ GLOBAL TEMPORARY | PRIVATE TEMPORARY ] TABLE ...
//	CREATE [ OR REPLACE ] [ UNIQUE | BITMAP ] INDEX ...
//	CREATE [ OR REPLACE ] [ FORCE | NO FORCE ] VIEW ...
//	...
func (p *Parser) parseCreateStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume CREATE

	// OR REPLACE. Oracle 23ai takes IF NOT EXISTS after the object's kind,
	// never before it (ORA-00901 at the IF of CREATE IF NOT EXISTS).
	orReplace := false
	if p.cur.Type == kwOR {
		p.advance() // consume OR
		if p.cur.Type != kwREPLACE {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume REPLACE
		orReplace = true
	} else if p.cur.Type == kwIF {
		return nil, p.syntaxErrorAtCur()
	}

	// EDITIONABLE | NONEDITIONABLE (only for PL/SQL objects: PROCEDURE, FUNCTION, PACKAGE, TRIGGER, TYPE, LIBRARY)
	// Note: Views handle EDITIONABLE internally in parseCreateViewStmt.
	editionable := false
	nonEditionable := false
	editionTok := p.cur // the EDITIONABLE or NONEDITIONABLE, when consumed
	if p.isIdentLikeStr("EDITIONABLE") {
		// Peek to see if the next token is a PL/SQL object keyword
		next := p.peekNext()
		if next.Type == kwPROCEDURE || next.Type == kwFUNCTION || next.Type == kwPACKAGE ||
			next.Type == kwTRIGGER || next.Type == kwTYPE || next.Type == kwLIBRARY {
			editionable = true
			p.advance()
		}
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		next := p.peekNext()
		if next.Type == kwPROCEDURE || next.Type == kwFUNCTION || next.Type == kwPACKAGE ||
			next.Type == kwTRIGGER || next.Type == kwTYPE || next.Type == kwLIBRARY {
			nonEditionable = true
			p.advance()
		}
	}

	// GLOBAL TEMPORARY | PRIVATE TEMPORARY: GLOBAL and PRIVATE lead only a
	// temporary table (ORA-00901 for CREATE GLOBAL TABLE).
	global := false
	private := false
	if p.cur.Type == kwGLOBAL {
		globalTok := p.advance() // consume GLOBAL
		if p.cur.Type != kwTEMPORARY {
			return nil, p.syntaxErrorAtTok(globalTok)
		}
		p.advance() // consume TEMPORARY
		global = true
	}

	// PRIVATE also leads CREATE PRIVATE OUTLINE.
	if p.cur.Type == kwPRIVATE {
		privateTok := p.advance() // consume PRIVATE
		switch p.cur.Type {
		case kwTEMPORARY:
			p.advance() // consume TEMPORARY
			private = true
		default:
			if !p.isKeywordStr("OUTLINE") {
				return nil, p.syntaxErrorAtTok(privateTok)
			}
		}
	}

	// IMMUTABLE TABLE | BLOCKCHAIN TABLE
	if (p.isKeywordStr("IMMUTABLE") || p.isKeywordStr("BLOCKCHAIN")) && p.peekNext().Type == kwTABLE &&
		!orReplace && !global && !private {
		stmt := &nodes.CreateTableStmt{
			Immutable:  p.cur.Str == "IMMUTABLE",
			Blockchain: p.cur.Str == "BLOCKCHAIN",
			Loc:        nodes.Loc{Start: start},
		}
		p.advance() // consume IMMUTABLE / BLOCKCHAIN
		return p.parseCreateTableOfKind(stmt)
	}

	// SHARDED / DUPLICATED
	sharded := false
	duplicated := false
	if p.isIdentLikeStr("SHARDED") {
		p.advance() // consume SHARDED
		sharded = true
	} else if p.isIdentLikeStr("DUPLICATED") {
		p.advance() // consume DUPLICATED
		duplicated = true
	}

	// SHARED [ PUBLIC ] DATABASE LINK
	if !orReplace && !editionable && !nonEditionable && !global && !private && !sharded && !duplicated && p.isKeywordStr("SHARED") {
		p.advance() // consume SHARED
		sharedPublic := false
		if p.cur.Type == kwPUBLIC {
			sharedPublic = true
			p.advance()
		}
		return p.parseCreateDatabaseLinkStmt(start, sharedPublic, true)
	}

	// PUBLIC (for CREATE PUBLIC SYNONYM / DATABASE LINK)
	public := false
	var publicTok *Token
	if p.cur.Type == kwPUBLIC {
		public = true
		tok := p.cur
		publicTok = &tok
		p.advance()
	}

	// [ AND { RESOLVE | COMPILE } ] [ NOFORCE ] — only for CREATE [OR REPLACE]
	// JAVA, before the JAVA keyword. Record them for the JAVA statement.
	var javaMods []*nodes.DDLOption
	if p.cur.Type == kwAND {
		optStart := p.pos()
		p.advance() // consume AND
		if !p.isKeywordStr("RESOLVE") && !p.isKeywordStr("COMPILE") {
			return nil, p.syntaxErrorAtCur()
		}
		value := p.cur.Str
		p.advance()
		javaMods = append(javaMods, &nodes.DDLOption{Key: "AND", Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	}
	if p.isKeywordStr("NOFORCE") && p.peekNext().Type == kwJAVA {
		javaMods = append(javaMods, &nodes.DDLOption{Key: "NOFORCE", Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}})
		p.advance() // consume NOFORCE
	}
	if len(javaMods) > 0 {
		if p.cur.Type != kwJAVA {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume JAVA
		return p.parseCreateJavaStmt(start, orReplace, javaMods)
	}

	if p.cur.Type == kwLIBRARY {
		p.advance() // consume LIBRARY
		var edition *Token
		if editionable || nonEditionable {
			edition = &editionTok
		}
		return p.parseCreateLibraryStmt(start, orReplace, edition)
	}

	// CONTEXT, RESTORE POINT, ROLLBACK SEGMENT, INMEMORY JOIN GROUP
	if stmt, ok, err := p.parseCreateAdminTail(start, orReplace, publicTok); ok {
		return stmt, err
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, global, private, sharded, duplicated)
	case kwUNIQUE, kwBITMAP, kwINDEX:
		return p.parseCreateIndexStmt(start)
	case kwMATERIALIZED:
		return p.parseCreateMaterializedOrView(start, orReplace)
	case kwFORCE, kwVIEW:
		if p.cur.Type == kwFORCE && isAnalyticObjectWord(p.peekNext()) {
			return p.parseCreateForcedAnalyticObject(start, orReplace)
		}
		return p.parseCreateViewStmt(start, orReplace)
	case kwSEQUENCE:
		return p.parseCreateSequenceStmt(start)
	case kwSYNONYM:
		return p.parseCreateSynonymStmt(start, orReplace, public)
	case kwDATABASE:
		// Distinguish CREATE DATABASE LINK from CREATE DATABASE
		next := p.peekNext()
		if next.Type == kwLINK {
			return p.parseCreateDatabaseLinkStmt(start, public, false)
		}
		return p.parseCreateDatabaseStmt(start)
	case kwTYPE:
		return p.parseCreateTypeStmt(start, orReplace, false, editionable, nonEditionable)
	case kwPROCEDURE:
		return p.parseCreateProcedureStmt(start, orReplace, false, editionable, nonEditionable)
	case kwFUNCTION:
		return p.parseCreateFunctionStmt(start, orReplace, false, editionable, nonEditionable)
	case kwPACKAGE:
		return p.parseCreatePackageStmt(start, orReplace, false, editionable, nonEditionable)
	case kwTRIGGER:
		return p.parseCreateTriggerStmt(start, orReplace, false, editionable, nonEditionable)
	case kwAUDIT:
		// CREATE AUDIT POLICY
		p.advance() // consume AUDIT
		if !p.isKeywordStr("POLICY") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume POLICY
		return p.parseCreateAuditPolicyStmt(start)
	case kwJSON:
		// CREATE JSON RELATIONAL DUALITY VIEW
		p.advance() // consume JSON
		if err := p.parseJsonDualityViewKind(false); err != nil {
			return nil, err
		}
		return p.parseCreateJsonDualityViewStmt(start, orReplace, false)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwDIRECTORY,
		kwCLUSTER, kwJAVA, kwSCHEMA:
		return p.parseCreateAdminObject(start, orReplace)
	case kwTEMPORARY:
		// CREATE TEMPORARY TABLESPACE (standalone, not GLOBAL TEMPORARY TABLE)
		if !global && !private {
			if p.peekNext().Type == kwTABLESPACE {
				return p.parseCreateTablespaceFamily(start)
			}
		}
		return nil, p.syntaxErrorAtCur()
	default:
		if p.isKeywordStr("NOFORCE") {
			if isAnalyticObjectWord(p.peekNext()) {
				return p.parseCreateForcedAnalyticObject(start, orReplace)
			}
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for "NO FORCE VIEW"
		if p.isIdentLikeStr("NO") || p.cur.Type == kwNOT {
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for EDITIONABLE / NONEDITIONABLE [ PUBLIC ] SYNONYM
		if !public && (p.isKeywordStr("EDITIONABLE") || p.isKeywordStr("NONEDITIONABLE")) {
			next := p.peekNext()
			if next.Type == kwSYNONYM ||
				(next.Type == kwPUBLIC && len(p.peekAhead(2)) == 2 && p.peekAhead(2)[1].Type == kwSYNONYM) {
				return p.parseCreateSynonymStmt(start, orReplace, public)
			}
		}
		// Check for EDITIONING / EDITIONABLE / NONEDITIONABLE VIEW
		if p.isIdentLike() && (p.cur.Str == "EDITIONING" || p.cur.Str == "EDITIONABLE" || p.cur.Str == "NONEDITIONABLE") {
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for BIGFILE/SMALLFILE/UNDO TABLESPACE, CONTROLFILE
		if p.isIdentLike() {
			switch p.cur.Str {
			case "BIGFILE", "SMALLFILE", "UNDO", "LOCAL":
				return p.parseCreateTablespaceFamily(start)
			case "CONTROLFILE":
				p.advance() // consume CONTROLFILE
				return p.parseCreateControlfileStmt(start)
			case "SNAPSHOT":
				if p.cur.Type == tokQIDENT {
					return nil, p.syntaxErrorAtCur()
				}
				// Oracle treats CREATE [OR REPLACE] SNAPSHOT as a legacy synonym
				// for CREATE [OR REPLACE] MATERIALIZED VIEW (and SNAPSHOT LOG).
				p.advance() // consume SNAPSHOT
				if p.cur.Type == kwLOG {
					p.advance() // consume LOG
					return p.parseCreateMviewLogStmt(start)
				}
				return p.parseCreateMaterializedViewStmt(&nodes.CreateViewStmt{
					OrReplace:    orReplace,
					Materialized: true,
					Loc:          nodes.Loc{Start: start},
				})
			}
		}
		// Check for MULTIVALUE INDEX
		if p.isIdentLikeStr("MULTIVALUE") {
			return p.parseCreateIndexStmt(start)
		}
		// Check for INDEXTYPE, OPERATOR, ANALYTIC VIEW, ATTRIBUTE DIMENSION, HIERARCHY, DOMAIN
		if p.isIdentLike() {
			switch p.cur.Str {
			case "INDEXTYPE":
				p.advance() // consume INDEXTYPE
				return p.parseCreateIndextypeStmt(start, orReplace)
			case "OPERATOR":
				p.advance() // consume OPERATOR
				return p.parseCreateOperatorStmt(start, orReplace, false)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type != kwVIEW {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume VIEW
				return p.parseCreateAnalyticViewStmt(start, orReplace, false, false)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if !p.isKeywordStr("DIMENSION") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume DIMENSION
				return p.parseCreateAttributeDimensionStmt(start, orReplace, false, false)
			case "HIERARCHY":
				p.advance() // consume HIERARCHY
				return p.parseCreateHierarchyStmt(start, orReplace, false, false)
			case "DOMAIN":
				p.advance() // consume DOMAIN
				return p.parseCreateDomainStmt(start, orReplace, false)
			case "USECASE":
				p.advance() // consume USECASE
				if p.isIdentLike() && p.cur.Str == "DOMAIN" {
					p.advance() // consume DOMAIN
					return p.parseCreateDomainStmt(start, orReplace, true)
				}
			case "OUTLINE":
				p.advance() // consume OUTLINE
				return p.parseCreateOutlineStmt(start, orReplace, public)
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if err := p.expectWords("GRAPH"); err != nil {
					return nil, err
				}
				return p.parseCreatePropertyGraphStmt(start, orReplace, false)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type != kwINDEX {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume INDEX
				return p.parseCreateVectorIndexStmt(start, false)
			case "MLE":
				p.advance() // consume MLE
				switch {
				case p.isKeywordStr("ENV"):
					p.advance() // consume ENV
					return p.parseCreateMLEEnvStmt(start, orReplace)
				case p.isKeywordStr("MODULE"):
					p.advance() // consume MODULE
					return p.parseCreateMLEModuleStmt(start, orReplace)
				}
				return nil, p.syntaxErrorAtCur()
			}
		}
		// Check for DIMENSION, FLASHBACK ARCHIVE
		adminStmt, err := p.parseCreateAdminObject(start, orReplace)
		if err != nil {
			return nil, err
		}
		if adminStmt != nil {
			return adminStmt, nil
		}
		return nil, nil
	}
}

// parseCreateTableStmt parses a CREATE TABLE statement after the TABLE keyword.
// The caller has already consumed CREATE [OR REPLACE] [GLOBAL TEMPORARY | PRIVATE TEMPORARY | SHARDED | DUPLICATED].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLE.html
//
// BNF (CREATE-TABLE.bnf lines 1-14):
//
//	CREATE [ { GLOBAL TEMPORARY | PRIVATE TEMPORARY | SHARDED | DUPLICATED } ]
//	    TABLE [ IF NOT EXISTS ] [ schema. ] table
//	    [ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
//	    { relational_table | object_table | XMLType_table | JSON_collection_table }
//	    [ MEMOPTIMIZE FOR READ ]
//	    [ MEMOPTIMIZE FOR WRITE ]
//	    [ PARENT [ schema. ] table ] ;
func (p *Parser) parseCreateTableStmt(start int, orReplace, global, private, sharded, duplicated bool) (*nodes.CreateTableStmt, error) {
	return p.parseCreateTableOfKind(&nodes.CreateTableStmt{
		OrReplace:  orReplace,
		Global:     global,
		Private:    private,
		Sharded:    sharded,
		Duplicated: duplicated,
		Loc:        nodes.Loc{Start: start},
	})
}

// parseCreateTableOfKind parses a CREATE TABLE statement after the TABLE
// keyword into stmt, which holds the kind of table the words before TABLE
// named and the statement's start. CREATE IMMUTABLE TABLE and CREATE
// BLOCKCHAIN TABLE take the immutable_table_clauses and
// blockchain_table_clauses.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html
func (p *Parser) parseCreateTableOfKind(stmt *nodes.CreateTableStmt) (*nodes.CreateTableStmt, error) {
	p.advance() // consume TABLE
	immutable, blockchain := stmt.Immutable, stmt.Blockchain
	stmt.Columns = &nodes.List{}
	stmt.Constraints = &nodes.List{}

	// IF NOT EXISTS (Oracle 23c): IF EXISTS and IF NOT without EXISTS fail at IF.
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		if p.peekNext().Type == kwNOT {
			ifTok := p.advance() // consume IF
			p.advance()          // consume NOT
			if p.cur.Type != kwEXISTS {
				return nil, p.syntaxErrorAtTok(ifTok)
			}
			p.advance() // consume EXISTS
			stmt.IfNotExists = true
		}
	}

	// Table name
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	nameTok := p.cur
	var parseErr487 error
	stmt.Name, parseErr487 = p.parseReservedCheckedObjectName()
	if parseErr487 !=

		// SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
		nil {
		return nil, parseErr487
	}
	if stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	// A relational table is created from its column list or from its AS
	// subquery: Oracle refuses one with neither at the word after its name
	// (ORA-00906), once it has read the rest of the statement.
	afterName := p.cur
	defer func(saved bool) { p.createTableEngineCheck = saved }(p.createTableEngineCheck)
	p.createTableEngineCheck = false
	defer func(saved map[string]bool) { p.createTableConstraintNames = saved }(p.createTableConstraintNames)
	p.createTableConstraintNames = map[string]bool{}

	if p.isIdentLikeStr("SHARING") {
		if err := p.parseCreateTableSharing(stmt); err != nil {
			return nil, err
		}
	}

	// Check for CTAS: AS subquery
	if p.cur.Type == kwAS {
		p.advance()
		var // consume AS
		parseErr488 error
		stmt.AsQuery, parseErr488 = p.parseSelectStmt()
		if parseErr488 != nil {
			return nil, parseErr488
		}
		if !selectStmtHasTargets(stmt.AsQuery) {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Column definitions and table constraints
	hasColumns := false
	if p.cur.Type == kwOF {
		if err := p.parseCreateTableObjectTable(stmt); err != nil {
			return nil, err
		}
		hasColumns = true
	} else if p.cur.Type == '(' {
		p.advance()
		parseErr489 := // consume '('
			p.parseColumnDefsAndConstraints(stmt)
		if parseErr489 != nil {
			return nil, parseErr489
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume ')'
		hasColumns = true
	}

	// immutable_table_clauses and blockchain_table_clauses belong to CREATE
	// IMMUTABLE TABLE and CREATE BLOCKCHAIN TABLE: Oracle 23ai refuses NO
	// DROP, NO DELETE and HASHING on any other table (ORA-00922).
	if immutable {
		if err := p.parseImmutableTableClauses(stmt); err != nil {
			return nil, err
		}
	}
	if blockchain {
		if err := p.parseBlockchainTableClauses(stmt); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if err := p.expectCreateTableWord("COLLATION"); err != nil {
			return nil, err
		}
		p.createTableEngineCheck = true
		var parseErr492 error
		stmt.Collation, parseErr492 = p.parseIdentifier()
		if parseErr492 != nil {
			return nil, parseErr492
		}
	}
	parseErr493 := p.parseTableOptions(stmt, hasColumns)
	if parseErr493 != nil {
		return nil, parseErr493
	}

	// [ MEMOPTIMIZE FOR READ ] [ MEMOPTIMIZE FOR WRITE ]
	for p.isKeywordStr("MEMOPTIMIZE") {
		p.advance() // consume MEMOPTIMIZE
		if err := p.expectCreateTableToken(kwFOR); err != nil {
			return nil, err
		}
		switch {
		case p.cur.Type == kwREAD && !stmt.MemoptimizeRead:
			stmt.MemoptimizeRead = true
		case p.cur.Type == kwWRITE && !stmt.MemoptimizeWrite:
			stmt.MemoptimizeWrite = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	// PARENT [ schema. ] table
	if p.isIdentLikeStr("PARENT") {
		p.advance()
		p.createTableEngineCheck = true
		parentTok := p.cur
		var // consume PARENT
		parseErr494 error
		stmt.Parent, parseErr494 = p.parseObjectName()
		if parseErr494 != nil {
			return nil, parseErr494
		}
		if stmt.Parent != nil && stmt.Name != nil && stmt.Parent.Name != "" &&
			stmt.Parent.Schema == stmt.Name.Schema && stmt.Parent.Name == stmt.Name.Name {
			return nil, p.syntaxErrorAtTok(parentTok)
		}
	}

	if !hasColumns && stmt.AsQuery == nil && p.isStatementEnd() && !p.createTableEngineCheck {
		if stmt.Sharded {
			return nil, p.syntaxErrorAtTok(nameTok)
		}
		return nil, p.syntaxErrorAtTok(afterName)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

func selectStmtHasTargets(stmt nodes.StmtNode) bool {
	selectStmt, ok := stmt.(*nodes.SelectStmt)
	if !ok {
		return false
	}
	// A query around a parenthesized subquery has its targets there.
	if selectStmt.Larg != nil {
		return true
	}
	return selectStmt.TargetList != nil && selectStmt.TargetList.Len() > 0
}

// parseColumnDefsAndConstraints parses the contents inside the parentheses of a CREATE TABLE.
// It handles both column definitions and table-level constraints.
//
// A supplemental_logging_props clause is one of relational_properties too
// (CREATE-TABLE.bnf lines 58-65 and 183): Oracle reads SUPPLEMENTAL LOG
// GROUP a(a) inside the list and refuses the closing ')' of a list that
// holds no column, as it refuses a list of constraints alone. It is kept as
// a SUPPLEMENTAL LOG DDLOption, whose group ends the walk in
// oracle/review (the group is a constraint), as the object_properties path
// already keeps it.
func (p *Parser) parseColumnDefsAndConstraints(stmt *nodes.CreateTableStmt) error {
	defer func(saved bool) { p.createTableColumns = saved }(p.createTableColumns)
	p.createTableColumns = true
	periods, supplemental, domains := 0, 0, 0
	seenCols := map[string]bool{}
	for progress := -1; p.cur.Type != ')' && p.cur.Type != tokEOF; {
		// An item that consumes nothing would come round again.
		if p.pos() == progress {
			return p.syntaxErrorAtCur()
		}
		progress = p.pos()
		before := stmt.Columns.Len() + stmt.Constraints.Len() + periods + supplemental + domains
		if p.atWord("SUPPLEMENTAL") && (p.atNextWord("LOG") || p.peekNext().Type == kwGROUP) {
			start, valueStart := p.pos(), p.peekNext().End
			if err := p.parseCreateTableSupplementalLogging(); err != nil {
				return err
			}
			appendCreateTableOption(stmt, "SUPPLEMENTAL LOG", p.createTableOptionValue(valueStart), nodes.Loc{Start: start, End: p.prev.End})
			supplemental++
		} else if p.isCreateTableDomainClauseStart() {
			start := p.pos()
			if err := p.parseCreateTableDomainClause(stmt, start); err != nil {
				return err
			}
			domains++
			// Check if this is a table-level constraint
		} else if p.isCreateTablePeriodStart() {
			if err := p.parseCreateTablePeriod(stmt); err != nil {
				return err
			}
			periods++
		} else if p.isOutOfLineRefConstraintStart() {
			tc, err := p.parseOutOfLineRefConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
		} else if p.isTableConstraintStart() {
			tc, parseErr495 := p.parseTableConstraint()
			if parseErr495 != nil {
				return parseErr495
			}
			if tc != nil {
				stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
			}
		} else {
			// Column definition
			colTok := p.cur
			if p.isIdentLike() && seenCols[colTok.Str] {
				ahead := p.peekNext()
				if ahead.Type == ',' || ahead.Type == ')' || (ahead.Type == tokIDENT && !isOracleBuiltInScalarTypeName(ahead) && !isColumnPropertyKeyword(ahead) && p.peekAhead(2)[1].Type == ')') {
					return p.syntaxErrorAtTok(colTok)
				}
			}
			col, parseErr496 := p.parseColumnDef()
			if parseErr496 != nil {
				return parseErr496
			}
			if col != nil {
				if col.TypeName == nil && col.Domain == nil && seenCols[col.Name] {
					return p.syntaxErrorAtTok(colTok)
				}
				seenCols[col.Name] = true
				stmt.Columns.Items = append(stmt.Columns.Items, col)
			}
		}
		if stmt.Columns.Len()+stmt.Constraints.Len()+periods+supplemental+domains == before {
			return p.syntaxErrorAtCur()
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		if p.cur.Type == ')' || p.cur.Type == tokEOF {
			return p.syntaxErrorAtCur()
		}
	}
	// relational_properties holds at least one item (ORA-00931 for
	// CREATE TABLE a()), and a table at least one column, which a
	// period_definition adds: a list of constraints alone is ORA-03079 at
	// its closing parenthesis.
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	if stmt.Columns.Len()+stmt.Constraints.Len()+periods == 0 {
		return p.syntaxErrorAtCur()
	}
	if stmt.Columns.Len()+periods == 0 {
		return p.syntaxErrorAtCur()
	}
	return p.checkOutOfLineRefConstraints(stmt, true)
}

// isCreateTableDomainClauseStart reports a table-level multi-column
// domain_clause inside CREATE TABLE's relational_properties:
//
//	DOMAIN [ schema. ] domain_name ( column [, column ] ... )
//
// A column named DOMAIN with a built-in datatype (e.g. DOMAIN VARCHAR2(10))
// or without a parenthesized column list (e.g. DOMAIN my_type) falls through
// to parseColumnDef.
func (p *Parser) isCreateTableDomainClauseStart() bool {
	if !p.atWord("DOMAIN") {
		return false
	}
	ahead := p.peekAhead(5)
	if len(ahead) < 2 {
		return false
	}
	if isOracleTypeToken(ahead[0].Type) {
		return false
	}
	if !p.isCreateTableNameTok(ahead[0]) {
		return false
	}
	lparen := 1
	if ahead[1].Type == '.' {
		if len(ahead) < 4 || !p.isCreateTableNameTok(ahead[2]) {
			return false
		}
		lparen = 3
	}
	if ahead[lparen].Type != '(' {
		return false
	}
	// A numeric argument like DOMAIN my_type(10) is a type modifier, not a
	// multi-column domain column list.
	if lparen+1 < len(ahead) && ahead[lparen+1].Type == tokICONST {
		return false
	}
	return true
}

// isCreateTableNameTok reports whether tok can be an object/domain name token.
func (p *Parser) isCreateTableNameTok(tok Token) bool {
	if tok.Type == tokIDENT || tok.Type == tokQIDENT {
		return true
	}
	return tok.Type >= 2000 && !isOracleSQLReservedKeyword(tok)
}

// parseCreateTableDomainClause reads CREATE TABLE's table-level domain_clause:
//
//	DOMAIN [ schema. ] domain_name ( column [, column ] ... )
func (p *Parser) parseCreateTableDomainClause(stmt *nodes.CreateTableStmt, start int) error {
	p.advance() // consume DOMAIN
	valueStart := p.pos()
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	if name == nil || name.Name == "" {
		return p.syntaxErrorAtCur()
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
		p.advance()
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	appendCreateTableOption(stmt, "DOMAIN", p.createTableOptionValue(valueStart), nodes.Loc{Start: start, End: p.prev.End})
	return nil
}

// isTableConstraintStart checks if the current position starts a table-level constraint.
func (p *Parser) isTableConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwCONSTRAINTS:
		// A column named CONSTRAINT takes a reserved datatype word alone:
		// Oracle 23ai reads CONSTRAINT NUMBER as a column and CONSTRAINT
		// NCHAR(1), CLOB, or REF u as a constraint (ORA-03050, ORA-00931).
		switch p.peekNext().Type {
		case kwNUMBER, kwINTEGER, kwSMALLINT, kwDECIMAL, kwFLOAT, kwCHAR, kwVARCHAR2, kwVARCHAR,
			kwDATE, kwRAW, kwLONG, kwROWID:
			return false
		}
		return true
	case kwPRIMARY, kwFOREIGN:
		// PRIMARY and FOREIGN are not reserved: without KEY they start a
		// column definition (engine-checked by TestOracleNonReservedKeywordsAsColumns).
		return p.peekNext().Type == kwKEY
	case kwUNIQUE:
		// UNIQUE could be a column constraint when it follows a type,
		// but at the top level of the column list it's a table constraint
		// only if followed by '('
		next := p.peekNext()
		return next.Type == '('
	case kwCHECK:
		return true
	}
	return false
}

// parseColumnDef parses a single column definition.
//
// BNF (CREATE-TABLE.bnf lines 69-80):
//
//	column [ datatype | DOMAIN [ schema. ] domain_name ]
//	[ SORT ]
//	[ VISIBLE | INVISIBLE ]
//	[ DEFAULT [ ON NULL [ FOR INSERT ONLY ] ] expr | identity_clause ]
//	[ ENCRYPT encryption_spec ]
//	[ COLLATE column_collation_name ]
//	[ inline_constraint ... ]
//	[ inline_ref_constraint ]
func (p *Parser) parseColumnDef() (*nodes.ColumnDef, error) {
	start := p.pos()
	col := &nodes.ColumnDef{
		Constraints: &nodes.List{},
		Loc:         nodes.Loc{Start: start},
	}

	if err := p.syntaxErrorIfReservedColumnIdentifier(); err != nil {
		return nil, err
	}
	colStartTok := p.cur
	var parseErr497 error
	col.Name, parseErr497 = p.parseIdentifier()
	if parseErr497 != nil {
		return nil, parseErr497
	}
	if col.Name == "" {
		return nil, nil
	}

	userDefinedType := false
	// DOMAIN [ schema. ] domain_name  (instead of datatype)
	if p.isIdentLikeStr("DOMAIN") {
		p.advance()
		var // consume DOMAIN
		parseErr498 error
		col.Domain, parseErr498 = p.parseObjectName()
		if parseErr498 != nil {
			return nil, parseErr498
		}
		if col.Domain == nil || col.Domain.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
	} else if (p.cur.Type == kwIDENTITY || p.cur.Type == kwKEY || p.cur.Type == kwNULLS ||
		((p.cur.Type == kwSAMPLE || p.cur.Type == kwUSING) && p.alterMview) ||
		(p.cur.Type == kwREAD && p.createTableColumns)) &&
		(p.alterTableModifyColumn ||
			(p.cur.Type == kwREAD && p.createTableColumns) ||
			((p.cur.Type == kwKEY || p.cur.Type == kwNULLS) && p.createTableColumns &&
				(p.peekNext().Type == ')' || p.peekNext().Type == ','))) {
		// Oracle leaves the datatype of a MODIFY column to the resolve
		// phase too: MODIFY (a KEY), MODIFY (a IDENTITY) and MODIFY (a
		// NULLS) are accepted (compat.tsv's rows).
		col.TypeName = &nodes.TypeName{
			Names:    &nodes.List{Items: []nodes.Node{&nodes.String{Str: p.cur.Str}}},
			TypeMods: &nodes.List{},
			Loc:      nodes.Loc{Start: p.pos(), End: p.cur.End},
		}
		p.advance()
	} else if p.isTypeName() && !p.isCreateTableColumnClauseStart() && !isColumnPropertyKeyword(p.cur) {
		// COLLATE, ENCRYPT, ANNOTATIONS and SCOPE IS start clauses of a
		// column without a datatype, not a datatype.
		typeTok := p.cur
		var parseErr499 error

		col.TypeName, parseErr499 = p.parseTypeName()
		if parseErr499 != nil {
			return nil, parseErr499
		}
		if !isOracleTypeToken(typeTok.Type) && typeTok.Type != kwREF && !isOracleBuiltInScalarTypeName(typeTok) {
			userDefinedType = true
		}
		// A user-defined type name (not a built-in parameterized type) cannot
		// take a parenthesized modifier; Oracle rejects at the type name
		// (ORA-03060: MODIFY (a LOB(a) ...), MODIFY (a AS(1), a(1)), ADD LOG DATA(ALL)).
		if (p.alterTableModifyColumn || p.createTableColumns || p.alterTableAddColumn || p.alterTableAddList) &&
			p.cur.Type == '(' && !isOracleTypeToken(typeTok.Type) &&
			col.TypeName != nil && !col.TypeName.IsPercType && !col.TypeName.IsPercRowtype {
			return nil, p.syntaxErrorAtTok(typeTok)
		}
	} else if (p.alterTableAddColumn ||
		((p.createTableColumns || p.alterTableModifyColumn || p.alterTableAddList) && !isColumnPropertyKeyword(p.cur))) &&
		p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur) && !p.isCreateTableColumnClauseStart() {
		// ALTER TABLE ... ADD's and CREATE TABLE's datatype is the resolve
		// phase's: a word Oracle does not reserve (and that does not start a
		// column property) names a user-defined type, as DEFERRED does in ADD
		// INITIALLY DEFERRED and VARRAY/NESTED/PARTITION do in CREATE TABLE
		// column definitions. Oracle reserves IMMEDIATE (ORA-00904 at it), so
		// ADD INITIALLY IMMEDIATE is no column. A size after such a word is
		// ORA-03060 at it (ADD PRIORITY KEY(a)), which is a syntax error: only
		// a type the server has takes one.
		if p.peekNext().Type == '(' {
			return nil, p.syntaxErrorAtCur()
		}
		if !isOracleBuiltInScalarTypeName(p.cur) {
			userDefinedType = true
		}
		col.TypeName = &nodes.TypeName{
			Names:    &nodes.List{Items: []nodes.Node{&nodes.String{Str: p.cur.Str}}},
			TypeMods: &nodes.List{},
			Loc:      nodes.Loc{Start: p.pos(), End: p.cur.End},
		}
		p.advance()
	}
	// PARTITION TABLE ... in CREATE TABLE column list: Oracle reads TABLE as
	// datatype and rejects at the following token.
	if col.TypeName == nil && p.createTableColumns && col.Name == "PARTITION" && p.cur.Type == kwTABLE {
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}
	// Bare ALWAYS AS (...) / ALWAYS AS IDENTITY without GENERATED:
	// Oracle reads ALWAYS as an invalid virtual-column datatype (ORA-54003 at
	// column name) before AS(...), and past IDENTITY before rejecting ALWAYS AS IDENTITY.
	if col.TypeName == nil && p.cur.Type == kwALWAYS && p.peekNext().Type == kwAS {
		ahead := p.peekAhead(2)
		if len(ahead) == 2 && ahead[1].Type == '(' {
			return nil, p.syntaxErrorAtTok(colStartTok)
		}
		if len(ahead) == 2 && ahead[1].Type == kwIDENTITY {
			p.advance() // ALWAYS
			p.advance() // AS
			p.advance() // IDENTITY
			return nil, p.syntaxErrorAtCur()
		}
	}
	parseErr500 := p.parseColumnProperties(col)
	if parseErr500 != nil {
		return nil, parseErr500
	}
	if p.alterTableAddColumn && col.TypeName != nil && p.isKeywordStr("DEFERRED") {
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}
	// A virtual column cannot take a non-built-in/user-defined datatype (ORA-54003 at column name).
	if col.Virtual != nil && userDefinedType {
		return nil, p.syntaxErrorAtTok(colStartTok)
	}
	// CREATE TABLE's column_definition takes its datatype as optional
	// (CREATE TABLE t (a, b) AS SELECT ...); a bare column elsewhere is
	// an error, but in ALTER TABLE ... ADD it is Oracle's ORA-02263, which
	// the statement survives, in a list as out of one (ADD (a)).
	if col.Domain == nil && col.TypeName == nil && !columnDefHasProperties(col) && !p.createTableColumns &&
		!p.alterTableAddColumn && !p.alterTableAddList {
		return nil, p.syntaxErrorAtCur()
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// isTypeName returns true if the current token can begin a type name.
func (p *Parser) isTypeName() bool {
	// REF object_type is a column datatype (Oracle 23ai).
	if isOracleTypeToken(p.cur.Type) || p.cur.Type == kwREF {
		return true
	}
	// User-defined type: identifier that is NOT a constraint keyword
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	return false
}

func isOracleTypeToken(tokenType int) bool {
	switch tokenType {
	case kwNUMBER, kwINTEGER, kwSMALLINT, kwDECIMAL, kwFLOAT,
		kwCHAR, kwVARCHAR2, kwVARCHAR, kwNCHAR, kwNVARCHAR2,
		kwCLOB, kwBLOB, kwNCLOB, kwJSON,
		kwDATE, kwTIMESTAMP, kwINTERVAL,
		kwRAW, kwLONG, kwROWID:
		return true
	default:
		return false
	}
}

func isOracleBuiltInScalarTypeName(tok Token) bool {
	if tok.Type == tokQIDENT {
		return false
	}
	switch tok.Str {
	case "INT", "DEC", "NUMERIC", "REAL", "DOUBLE",
		"CHARACTER", "NATIONAL", "BOOLEAN", "BOOL",
		"BINARY_FLOAT", "BINARY_DOUBLE", "UROWID",
		"BFILE", "XMLTYPE", "VECTOR", "ANYDATA", "SDO_GEOMETRY":
		return true
	}
	return false
}

// isColumnPropertyKeyword reports whether tok is a keyword that starts a column
// property/constraint clause rather than a user-defined datatype name in a
// CREATE TABLE column definition.
func isColumnPropertyKeyword(tok Token) bool {
	switch tok.Type {
	case kwDEFAULT, kwAS, kwGENERATED, kwNOT, kwNULL, kwINVISIBLE, kwWITH,
		kwCONSTRAINT, kwCONSTRAINTS, kwPRIMARY, kwUNIQUE, kwCHECK, kwREFERENCES,
		kwDROP, kwALWAYS, kwIDENTITY:
		return true
	}
	if tok.Type != tokQIDENT {
		switch tok.Str {
		case "SORT", "VISIBLE", "ENCRYPT", "DECRYPT", "COLLATE", "ANNOTATIONS", "SCOPE":
			return true
		}
	}
	return false
}

func columnDefHasProperties(col *nodes.ColumnDef) bool {
	return (col.Annotations != nil && col.Annotations.Len() > 0) || col.Sort || col.Visible || col.Invisible || col.Default != nil ||
		col.Identity != nil || col.Virtual != nil || col.NotNull || col.Null ||
		col.DropIdentity || col.Encrypt != "" || col.Collation != "" ||
		(col.Constraints != nil && col.Constraints.Len() > 0)
}

// parseColumnProperties parses the properties after the column type:
// SORT, VISIBLE, INVISIBLE, DEFAULT, NOT NULL, NULL, GENERATED (identity/virtual),
// ENCRYPT, COLLATE, constraints, etc.
//
// BNF (CREATE-TABLE.bnf lines 69-80, 100-109):
//
//	[ SORT ]
//	[ VISIBLE | INVISIBLE ]
//	[ DEFAULT [ ON NULL [ FOR INSERT ONLY ] ] expr | identity_clause ]
//	[ ENCRYPT encryption_spec ]
//	[ COLLATE column_collation_name ]
//	[ inline_constraint ... ]
func (p *Parser) parseColumnProperties(col *nodes.ColumnDef) error {
	for {
		switch p.cur.Type {
		case kwDEFAULT:
			p.advance() // consume DEFAULT
			// DEFAULT ON NULL [ FOR INSERT ONLY ] expr
			if p.cur.Type == kwON {
				p.advance() // consume ON
				if err := p.expectCreateTableToken(kwNULL); err != nil {
					return err
				}
				col.DefaultOnNull = true
				// FOR INSERT ONLY
				if p.cur.Type == kwFOR {
					p.advance() // consume FOR
					if err := p.expectCreateTableToken(kwINSERT); err != nil {
						return err
					}
					if err := p.expectCreateTableWord("ONLY"); err != nil {
						return err
					}
					col.DefaultOnNullInsertOnly = true
				}
			}
			if p.cur.Type == kwINSERT {
				return p.syntaxErrorAtCur()
			}
			var parseErr501 error
			col.Default, parseErr501 = p.parseExpr()
			if parseErr501 != nil {
				return parseErr501
			}
			if col.Default == nil {

				return p.syntaxErrorAtCur()
			}

		case kwAS:
			// GENERATED ALWAYS is optional for Oracle virtual columns.
			p.advance() // consume AS
			if err := p.parseVirtualColumnExpression(col); err != nil {
				return err
			}

		case kwGENERATED:
			// identity_clause: GENERATED { ALWAYS | BY DEFAULT [ ON NULL ] } AS IDENTITY [ ( options ) ]
			// virtual_column:  GENERATED ALWAYS AS ( expr ) [ VIRTUAL ]
			p.advance() // consume GENERATED
			always := false
			byDefault := false
			if p.isIdentLikeStr("ALWAYS") {
				p.advance() // consume ALWAYS
				always = true
			} else if p.cur.Type == kwBY {
				p.advance() // consume BY
				if err := p.expectCreateTableToken(kwDEFAULT); err != nil {
					return err
				}
				byDefault = true
				// ON NULL
				if p.cur.Type == kwON {
					p.advance() // consume ON
					if err := p.expectCreateTableToken(kwNULL); err != nil {
						return err
					}
				}
			}
			// GENERATED AS IDENTITY is GENERATED ALWAYS; a virtual column
			// says ALWAYS (ORA-02000 at the parenthesis of GENERATED AS
			// (expr)).
			if p.cur.Type != kwAS {
				return p.syntaxErrorAtCur()
			}
			if !always && !byDefault && p.peekNext().Type == '(' {
				return p.syntaxErrorAtCur()
			}
			if p.cur.Type == kwAS {
				p.advance() // consume AS
				if p.cur.Type == kwIDENTITY {
					// identity_clause
					p.advance() // consume IDENTITY
					identity := &nodes.IdentityClause{
						Always:  always,
						Options: &nodes.List{},
						Loc:     nodes.Loc{Start: col.Loc.Start},
					}
					// ( identity_options ) or bare identity_options (Oracle accepts
					// sequence options after AS IDENTITY without outer parens too).
					if p.cur.Type == '(' {
						p.advance()
						parseErr502 := p.parseIdentityOptions(identity)
						if parseErr502 != nil {
							return parseErr502
						}
						if err := p.expectCreateTableToken(')'); err != nil {
							return err
						}
					} else if p.isIdentityOptionStart() {
						parseErr502 := p.parseIdentityOptions(identity)
						if parseErr502 != nil {
							return parseErr502
						}
					}
					identity.Loc.End = p.prev.End
					col.Identity = identity
				} else if p.cur.Type == '(' {
					// virtual column: AS (expr) [VIRTUAL]
					if err := p.parseVirtualColumnExpression(col); err != nil {
						return err
					}
				} else {
					return p.syntaxErrorAtCur()
				}
			}
			_ = byDefault // stored implicitly: !always means BY DEFAULT

		case kwNOT:
			// NOT starts NOT NULL here: Oracle 23ai answers a NOT followed
			// by anything else with ORA-02000 at the word after it.
			not := p.cur
			p.advance() // consume NOT
			if p.cur.Type != kwNULL {
				return p.syntaxErrorAtCur()
			}
			// A second NULL or NOT NULL is ORA-02258 at it.
			if col.Null || col.NotNull {
				return p.syntaxErrorAtTok(not)
			}
			p.advance() // consume NULL
			col.NotNull = true
			var cs constraintState
			if err := p.parseConstraintState(&cs); err != nil {
				return err
			}

		case kwNULL:
			// A NULL constraint takes no constraint_state: Oracle 23ai
			// answers NULL DEFERRABLE, NULL RELY, NULL USING INDEX and
			// NULL ENABLE with ORA-03099 at the state.
			if col.Null || col.NotNull {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume NULL
			col.Null = true
			if p.isCreateTableConstraintStateStart() {
				return p.syntaxErrorAtCur()
			}

		case kwINVISIBLE:
			p.advance()
			col.Invisible = true

		case kwWITH:
			if err := p.parseInlineRefConstraint(col); err != nil {
				return err
			}

		case kwCONSTRAINT, kwCONSTRAINTS:
			cc, parseErr504 := p.parseColumnConstraint()
			if parseErr504 != nil {
				return parseErr504
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwPRIMARY:
			cc, parseErr505 := p.parseColumnConstraintInline()
			if parseErr505 != nil {
				return parseErr505
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwUNIQUE:
			cc, parseErr506 := p.parseColumnConstraintInline()
			if parseErr506 != nil {
				return parseErr506
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwCHECK:
			cc, parseErr507 := p.parseColumnConstraintInline()
			if parseErr507 != nil {
				return parseErr507
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwREFERENCES:
			cc, parseErr508 := p.parseColumnConstraintInline()
			if parseErr508 != nil {
				return parseErr508
			}
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwDROP:
			// DROP IDENTITY (for ALTER TABLE MODIFY column)
			next := p.peekNext()
			if next.Type == kwIDENTITY {
				p.advance() // consume DROP
				p.advance() // consume IDENTITY
				col.DropIdentity = true
			} else {
				return nil
			}

		default:
			// Handle SORT, VISIBLE, ENCRYPT, DECRYPT, COLLATE as identifier-like keywords
			if p.isIdentLike() {
				switch p.cur.Str {
				case "SCOPE":
					// The keyword, not a quoted "SCOPE" (ORA-03062 on Oracle 23ai).
					if !p.isKeywordStr("SCOPE") {
						return p.syntaxErrorAtCur()
					}
					if err := p.parseInlineRefConstraint(col); err != nil {
						return err
					}
				case "ANNOTATIONS":
					anns, err := p.parseCreateTableAnnotations()
					if err != nil {
						return err
					}
					col.Annotations = anns
				case "SORT":
					p.advance()
					col.Sort = true
				case "VISIBLE":
					p.advance()
					col.Visible = true
				case "ENCRYPT":
					p.advance() // consume ENCRYPT
					col.Encrypt = "ENCRYPT"
					parseErr509 :=
						// encryption_spec: USING 'algo' IDENTIFIED BY pw SALT|NO SALT
						p.parseEncryptionSpec(col)
					if parseErr509 != nil {
						return parseErr509
					}
				case "DECRYPT":
					p.advance() // consume DECRYPT
					col.Encrypt = "DECRYPT"
				case "COLLATE":
					// A column collation needs the column's datatype:
					// Oracle 23ai answers a COLLATE without one with
					// ORA-03099 at COLLATE.
					if col.TypeName == nil && col.Domain == nil {
						return p.syntaxErrorAtCur()
					}
					p.advance()
					var // consume COLLATE
					parseErr510 error
					col.Collation, parseErr510 = p.parseIdentifier()
					if parseErr510 != nil {
						return parseErr510
					}
					if col.Collation == "" {
						return p.syntaxErrorAtCur()
					}
				default:
					return nil
				}
			} else {
				return nil
			}
		}
	}
}

// parseInlineRefConstraint parses the SCOPE IS or WITH ROWID constraint of
// a REF column. Oracle 23ai rejects one on another column (ORA-22893) and
// a second SCOPE (ORA-22888).
//
//	SCOPE IS [ schema . ] scope_table | WITH ROWID
func (p *Parser) parseInlineRefConstraint(col *nodes.ColumnDef) error {
	// Oracle reads SCOPE IS and WITH ROWID on a column with no datatype, and
	// on any column of ALTER TABLE ... ADD (compat.tsv's inline_ref_constraint
	// rows); CREATE TABLE's typed non-REF column still refuses them, and so
	// does a MODIFY column, whose grammar has no inline_ref_constraint.
	if col.TypeName == nil && p.alterTableModifyColumn {
		return p.syntaxErrorAtCur()
	}
	if col.TypeName != nil && !col.TypeName.IsRef && !p.alterTableAddColumn && !p.alterTableAddList {
		return p.syntaxErrorAtCur()
	}
	cc := &nodes.ColumnConstraint{Loc: nodes.Loc{Start: p.pos()}}
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if err := p.expectCreateTableToken(kwROWID); err != nil {
			return err
		}
		cc.Type = nodes.CONSTRAINT_WITH_ROWID
	} else {
		for _, item := range col.Constraints.Items {
			if other, ok := item.(*nodes.ColumnConstraint); ok && other.Type == nodes.CONSTRAINT_SCOPE {
				return p.syntaxErrorAtCur()
			}
		}
		p.advance() // consume SCOPE
		if p.cur.Type != kwIS {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		table, err := p.parseRefScopeTable()
		if err != nil {
			return err
		}
		cc.Type = nodes.CONSTRAINT_SCOPE
		cc.RefTable = table
	}
	cc.Loc.End = p.prev.End
	col.Constraints.Items = append(col.Constraints.Items, cc)
	return nil
}

// parseRefScopeTable parses the scope table of a REF column, which is
// local (ORA-25124 on Oracle 23ai for one over a database link).
func (p *Parser) parseRefScopeTable() (*nodes.ObjectName, error) {
	table, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	// A schema names its table (ORA-00903 for SCOPE IS s.), and an @ its
	// database link (ORA-01729 for SCOPE IS o@).
	if table.Name == "" || p.prev.Type == '.' || p.prev.Type == '@' {
		return nil, p.syntaxErrorAtCur()
	}
	if table.DBLink != "" {
		return nil, p.syntaxErrorAtNode(table)
	}
	return table, nil
}

// checkOutOfLineRefConstraints checks the REF constraints of stmt that
// span its relational properties, as Oracle 23ai does. A constraint names
// a column of the table, which may come before or after it, or an
// attribute of an object column (o.r); a quoted name ("O.R") is a column.
// An attribute's column is of an object type (ORA-02337 for a scalar or
// REF one). A SCOPE or WITH ROWID constraint names a REF (ORA-22893), and
// a REF column or attribute has one SCOPE (ORA-22888) and not both a SCOPE
// and a referential constraint (ORA-22896). A REF is in no key (ORA-02329,
// ORA-02267 for a composite foreign key), has one referential constraint
// (ORA-30771), and that one names the object table alone (ORA-30726 for a
// column list). An attribute counts as a REF where a SCOPE or WITH ROWID
// constraint names it; otherwise its type is the engine's to resolve. When
// stmt holds every column of the table, complete, as a CREATE TABLE does,
// a table constraint names one of them, or an attribute of one (ORA-00904
// on Oracle 23ai for SCOPE FOR (missing), FOREIGN KEY (missing), and
// UNIQUE (missing.x)); ALTER TABLE ADD may name a column the table already
// has.
func (p *Parser) checkOutOfLineRefConstraints(stmt *nodes.CreateTableStmt, complete bool) error {
	type target struct {
		name      string
		attribute bool
	}
	type constraint struct {
		node       nodes.Node
		typ        nodes.ConstraintType
		keys       []target
		refColumns *nodes.List
		inline     bool
	}
	columns := map[string]*nodes.ColumnDef{}
	refs := map[target]bool{}
	var constraints []constraint
	for _, item := range stmt.Columns.Items {
		col, ok := item.(*nodes.ColumnDef)
		if !ok {
			continue
		}
		columns[col.Name] = col
		key := target{name: col.Name}
		refs[key] = col.TypeName != nil && col.TypeName.IsRef
		for _, cItem := range col.Constraints.Items {
			if cc, ok := cItem.(*nodes.ColumnConstraint); ok {
				constraints = append(constraints, constraint{cc, cc.Type, []target{key}, cc.RefColumns, true})
			}
		}
	}
	for _, item := range stmt.Constraints.Items {
		tc, ok := item.(*nodes.TableConstraint)
		if !ok || tc.Columns == nil {
			continue
		}
		c := constraint{node: tc, typ: tc.Type, refColumns: tc.RefColumns}
		for _, cItem := range tc.Columns.Items {
			column, ok := cItem.(*nodes.String)
			if !ok {
				continue
			}
			parts := p.attributePaths[column]
			root := column.Str
			if parts != nil {
				root = parts[0]
			}
			if complete && columns[root] == nil {
				return p.syntaxErrorAtNode(tc)
			}
			// The attribute of a column of a type that has none.
			if parts != nil && columns[root] != nil && p.scalarColumn(columns[root]) {
				return p.syntaxErrorAtNode(tc)
			}
			// An attribute is keyed by its names, which a quoted one may
			// join with dots alike: o."R.X" is not "O.R".x.
			key := target{name: column.Str}
			if parts != nil {
				key = target{name: strings.Join(parts, "\x00"), attribute: true}
			}
			c.keys = append(c.keys, key)
		}
		constraints = append(constraints, c)
	}
	for _, c := range constraints {
		if (c.typ == nodes.CONSTRAINT_SCOPE || c.typ == nodes.CONSTRAINT_WITH_ROWID) && len(c.keys) == 1 {
			refs[c.keys[0]] = true
		}
	}
	scoped := map[target]bool{}
	references := map[target]int{}
	for _, c := range constraints {
		if len(c.keys) != 1 {
			// No REF in a composite key.
			if c.typ == nodes.CONSTRAINT_PRIMARY || c.typ == nodes.CONSTRAINT_UNIQUE || c.typ == nodes.CONSTRAINT_FOREIGN {
				for _, key := range c.keys {
					if refs[key] {
						return p.syntaxErrorAtNode(c.node)
					}
				}
			}
			continue
		}
		key := c.keys[0]
		col := columns[key.name]
		if key.attribute {
			col = nil
		}
		switch c.typ {
		case nodes.CONSTRAINT_PRIMARY, nodes.CONSTRAINT_UNIQUE:
			if refs[key] {
				return p.syntaxErrorAtNode(c.node)
			}
		case nodes.CONSTRAINT_SCOPE, nodes.CONSTRAINT_WITH_ROWID:
			// CREATE TABLE's typed non-REF column refuses these, and so does
			// an out-of-line one; a column with no datatype, and an inline
			// one of ALTER TABLE ... ADD, takes them (compat.tsv's
			// inline_ref_constraint rows).
			if col != nil && col.TypeName != nil && !col.TypeName.IsRef && (complete || !c.inline) {
				return p.syntaxErrorAtNode(c.node)
			}
			if c.typ == nodes.CONSTRAINT_SCOPE {
				if scoped[key] || references[key] > 0 {
					return p.syntaxErrorAtNode(c.node)
				}
				scoped[key] = true
			}
		case nodes.CONSTRAINT_FOREIGN:
			if scoped[key] || refs[key] && (references[key] > 0 || c.refColumns != nil && c.refColumns.Len() > 0) {
				return p.syntaxErrorAtNode(c.node)
			}
			references[key]++
		}
	}
	return nil
}

// scalarColumn reports whether col is of a datatype the text shows to
// have no attributes: a REF, or a predefined SQL datatype.
func (p *Parser) scalarColumn(col *nodes.ColumnDef) bool {
	if col.TypeName == nil {
		return false
	}
	lead := p.predefinedTypeSpelling(col.TypeName)
	return col.TypeName.IsRef || plsqlNonCharacterTypes[lead] || lead == "CHAR" || lead == "CHARACTER" ||
		lead == "VARCHAR2" || lead == "VARCHAR" || lead == "CLOB"
}

// isOutOfLineRefConstraintStart reports whether a table's relational
// property here is SCOPE FOR ( column ) or REF ( column ).
func (p *Parser) isOutOfLineRefConstraintStart() bool {
	return p.isKeywordStr("SCOPE") && p.peekNext().Type == kwFOR || p.cur.Type == kwREF && p.peekNext().Type == '('
}

// parseOutOfLineRefConstraint parses the out-of-line constraint of a REF
// column.
//
//	SCOPE FOR ( column ) IS [ schema . ] scope_table | REF ( column ) WITH ROWID
func (p *Parser) parseOutOfLineRefConstraint() (*nodes.TableConstraint, error) {
	tc := &nodes.TableConstraint{Loc: nodes.Loc{Start: p.pos()}}
	scope := p.cur.Type != kwREF
	p.advance() // consume SCOPE or REF
	if scope {
		p.advance() // consume FOR
	}
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	column, err := p.parseColumnPath()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	tc.Columns = &nodes.List{Items: []nodes.Node{column}}
	if scope {
		if p.cur.Type != kwIS {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if tc.RefTable, err = p.parseRefScopeTable(); err != nil {
			return nil, err
		}
		tc.Type = nodes.CONSTRAINT_SCOPE
	} else {
		if err := p.expectCreateTableToken(kwWITH); err != nil {
			return nil, err
		}
		if err := p.expectCreateTableToken(kwROWID); err != nil {
			return nil, err
		}
		tc.Type = nodes.CONSTRAINT_WITH_ROWID
	}
	tc.Loc.End = p.prev.End
	return tc, nil
}

func (p *Parser) parseVirtualColumnExpression(col *nodes.ColumnDef) error {
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if expr == nil {
		return p.syntaxErrorAtCur()
	}
	col.Virtual = expr
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	if p.cur.Type == kwVIRTUAL {
		p.advance() // consume VIRTUAL
	}
	return p.parseCreateTableEditionClauses()
}

// isValidOracleEncryptionAlgorithm reports whether alg is a valid Oracle 23ai
// TDE encryption algorithm name (case-insensitive, ORA-28339 for any other).
func isValidOracleEncryptionAlgorithm(alg string) bool {
	switch strings.ToUpper(alg) {
	case "3DES168", "AES128", "AES192", "AES256",
		"ARIA128", "ARIA192", "ARIA256", "GOST256", "SEED128", "SM4":
		return true
	}
	return false
}

// isValidOracleIntegrityAlgorithm reports whether alg is a valid Oracle 23ai
// TDE integrity algorithm name (case-insensitive, ORA-28339 for any other).
func isValidOracleIntegrityAlgorithm(alg string) bool {
	switch strings.ToUpper(alg) {
	case "SHA-1", "NOMAC":
		return true
	}
	return false
}

// parseEncryptionSpec parses the encryption_spec after ENCRYPT keyword.
//
// BNF (CREATE-TABLE.bnf lines 95-98):
//
//	[ USING 'encrypt_algorithm' ]
//	[ IDENTIFIED BY password ]
//	[ SALT | NO SALT ]
func (p *Parser) parseEncryptionSpec(col *nodes.ColumnDef) error {
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		if p.cur.Type != tokSCONST || !isValidOracleEncryptionAlgorithm(p.cur.Str) {
			return p.syntaxErrorAtCur()
		}
		col.Encrypt = "ENCRYPT USING " + p.cur.Str
		p.advance()
	}
	if p.isKeywordStr("IDENTIFIED") {
		p.advance() // consume IDENTIFIED
		if err := p.expectCreateTableToken(kwBY); err != nil {
			return err
		}
		if p.cur.Type == ')' {
			p.advance()
		} else if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	switch {
	case p.isKeywordStr("SALT"):
		p.advance()
	case p.isKeywordStr("NO"):
		p.advance() // consume NO
		if err := p.expectCreateTableWord("SALT"); err != nil {
			return err
		}
	}
	return nil
}

// isIdentityOptionStart reports whether the current token starts an
// identity_options keyword (for bare identity options without outer parens).
func (p *Parser) isIdentityOptionStart() bool {
	if !p.isIdentLike() || p.cur.Type == tokQIDENT {
		return false
	}
	switch p.cur.Str {
	case "START", "INCREMENT", "MAXVALUE", "MINVALUE", "CACHE",
		"NOMAXVALUE", "NOMINVALUE", "CYCLE", "NOCYCLE", "NOCACHE", "ORDER", "NOORDER":
		return true
	}
	return false
}

// parseIdentityOptions parses identity column options inside parentheses.
//
// BNF (CREATE-TABLE.bnf lines 86-93):
//
//	[ START WITH { integer | LIMIT VALUE } ]
//	[ INCREMENT BY integer ]
//	[ { MAXVALUE integer | NOMAXVALUE } ]
//	[ { MINVALUE integer | NOMINVALUE } ]
//	[ { CYCLE | NOCYCLE } ]
//	[ { CACHE integer | NOCACHE } ]
//	[ { ORDER | NOORDER } ]
func (p *Parser) parseIdentityOptions(identity *nodes.IdentityClause) error {
	seen := map[string]bool{}
	for progress := -1; p.cur.Type != ')' && p.cur.Type != tokEOF; {
		// An item that consumes nothing would come round again.
		if p.pos() == progress {
			return p.syntaxErrorAtCur()
		}
		progress = p.pos()
		if !p.isIdentLike() {
			return nil
		}
		switch p.cur.Str {
		case "START":
			p.advance() // consume START
			if err := p.expectCreateTableToken(kwWITH); err != nil {
				return err
			}
			// START WITH LIMIT VALUE is ALTER TABLE MODIFY's: CREATE TABLE
			// answers it with ORA-03089 at LIMIT.
			if p.isIdentLikeStr("LIMIT") && !p.createTableColumns {
				p.advance() // consume LIMIT
				if err := p.expectCreateTableWord("VALUE"); err != nil {
					return err
				}
			} else if err := p.parseIdentitySignedInteger(); err != nil {
				return err
			}
		case "INCREMENT":
			p.advance() // consume INCREMENT
			if err := p.expectCreateTableToken(kwBY); err != nil {
				return err
			}
			if err := p.parseIdentitySignedInteger(); err != nil {
				return err
			}
		case "MAXVALUE", "MINVALUE", "CACHE":
			p.advance()
			if err := p.parseIdentitySignedInteger(); err != nil {
				return err
			}
		case "NOMAXVALUE", "NOMINVALUE", "CYCLE", "NOCYCLE", "NOCACHE", "ORDER", "NOORDER":
			if seen[p.cur.Str] {
				return p.syntaxErrorAtCur()
			}
			seen[p.cur.Str] = true
			p.advance()
		default:
			return nil
		}
	}
	return nil
}

// parseIdentitySignedInteger reads the integer of an identity option,
// which may carry a sign. Oracle 23ai answers a missing one with ORA-03089
// at the word after the option.
func (p *Parser) parseIdentitySignedInteger() error {
	if p.cur.Type == '-' || p.cur.Type == '+' {
		p.advance()
	}
	return p.parseCreateTableInteger()
}

// parseColumnConstraint parses a named column constraint: CONSTRAINT name
// <constraint_body>. A [ NOT ] NULL body is left to the caller; no other
// clause follows a name (ORA-02253 on Oracle 23ai for CONSTRAINT c DEFAULT
// 1, ORA-22890 for CONSTRAINT c SCOPE IS t and CONSTRAINT c WITH ROWID).
func (p *Parser) parseColumnConstraint() (*nodes.ColumnConstraint, error) {
	start := p.pos()
	constTok := p.advance() // consume CONSTRAINT or CONSTRAINTS

	if p.cur.Type == kwNULL {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != tokQIDENT && isOracleSQLReservedKeyword(p.cur) {
		return nil, p.syntaxErrorAtTok(constTok)
	}
	name, parseErr523 := p.parseIdentifier()
	if parseErr523 != nil {
		return nil, parseErr523
	}
	if name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	cc, parseErr524 := p.parseColumnConstraintInline()
	if parseErr524 != nil {
		return nil, parseErr524
	}
	if cc == nil {
		if p.cur.Type != kwNULL && (p.cur.Type != kwNOT || p.peekNext().Type != kwNULL) {
			return nil, p.syntaxErrorAtTok(constTok)
		}
		if p.createTableConstraintNames != nil {
			p.createTableConstraintNames[name] = true
		}
		return nil, nil
	}
	if p.createTableConstraintNames != nil {
		p.createTableConstraintNames[name] = true
	}
	cc.Name = name
	cc.Loc.Start = start
	return cc, nil
}

// parseColumnConstraintInline parses an inline (unnamed) column constraint body.
func (p *Parser) parseColumnConstraintInline() (*nodes.ColumnConstraint, error) {
	start := p.pos()
	cc := &nodes.ColumnConstraint{
		Loc: nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if err := p.expectCreateTableToken(kwKEY); err != nil {
			return nil, err
		}
		cc.Type = nodes.CONSTRAINT_PRIMARY

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		if p.alterTableAddColumn && p.cur.Type == '(' {
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
		cc.Type = nodes.CONSTRAINT_UNIQUE

	case kwCHECK:
		p.advance() // consume CHECK
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr525 error
		cc.Expr, parseErr525 = p.parseExpr()
		if parseErr525 != nil {
			return nil, parseErr525
		}
		if cc.Expr == nil || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume ')'
		cc.Type = nodes.CONSTRAINT_CHECK

	case kwREFERENCES:
		p.advance() // consume REFERENCES
		cc.Type = nodes.CONSTRAINT_FOREIGN
		var parseErr526 error
		cc.RefTable, parseErr526 = p.parseReferencedTable()
		if parseErr526 != nil {
			return nil, parseErr526
		}
		if p.cur.Type == '(' {
			p.advance()
			var // consume '('
			parseErr527 error
			cc.RefColumns, parseErr527 = p.parseUniqueIdentifierList()
			if parseErr527 != nil {
				return nil, parseErr527
			}
			if cc.RefColumns.Len() == 0 || p.prev.Type == ',' {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return nil, err
			}
		}
		// ON DELETE
		if p.cur.Type == kwON {
			p.advance() // consume ON
			if err := p.expectCreateTableToken(kwDELETE); err != nil {
				return nil, err
			}
			var parseErr528 error
			cc.OnDelete, parseErr528 = p.parseDeleteAction()
			if parseErr528 != nil {
				return nil, parseErr528
			}
		}

	default:
		return nil, nil
	}

	var cs constraintState
	if err := p.parseConstraintState(&cs); err != nil {
		return nil, err
	}
	cc.Deferrable = cs.Deferrable
	cc.Initially = cs.Initially
	cc.Tablespace = cs.Tablespace
	cc.UsingIndex = cs.UsingIndex
	cc.UsingIndexLocal = cs.UsingIndexLocal

	cc.Loc.End = p.prev.End
	return cc, nil
}

// parseTableConstraint parses a table-level constraint including its
// trailing constraint_state clause.
//
//	[ CONSTRAINT name ] { PRIMARY KEY (cols) | UNIQUE (cols) | CHECK (expr) | FOREIGN KEY (cols) REFERENCES ... } constraint_state
func (p *Parser) parseTableConstraint() (*nodes.TableConstraint, error) {
	tc, err := p.parseTableConstraintBody(true)
	if err != nil {
		return nil, err
	}

	var cs constraintState
	if err := p.parseConstraintState(&cs); err != nil {
		return nil, err
	}
	tc.Deferrable = cs.Deferrable
	tc.Initially = cs.Initially
	tc.Tablespace = cs.Tablespace
	tc.UsingIndex = cs.UsingIndex
	tc.UsingIndexLocal = cs.UsingIndexLocal

	tc.Loc.End = p.prev.End
	return tc, nil
}

// parseTableConstraintBody parses a table-level constraint without its
// constraint_state clause. Used directly by view constraint parsing, where
// the allowed state differs — see parseViewConstraintState.
func (p *Parser) parseTableConstraintBody(allowOnDelete bool) (*nodes.TableConstraint, error) {
	start := p.pos()
	tc := &nodes.TableConstraint{
		Loc: nodes.Loc{Start: start},
	}

	// Optional CONSTRAINT name; Oracle takes CONSTRAINTS here too.
	hasConstraintName := false
	if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwCONSTRAINTS {
		p.advance()
		var // consume CONSTRAINT
		parseErr529 error
		tc.Name, parseErr529 = p.parseObjectNameIdentifier()
		if parseErr529 != nil {
			return nil, parseErr529
		}
		if tc.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		hasConstraintName = true
		if tc.Name == "FOREIGN" && p.cur.Type == kwKEY && p.peekNext().Type == '(' {
			p.advance() // KEY
			return nil, p.syntaxErrorAtCur()
		}
		if allowOnDelete && p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur) &&
			p.cur.Type != kwPRIMARY && p.cur.Type != kwUNIQUE && p.cur.Type != kwCHECK && p.cur.Type != kwFOREIGN {
			next := p.peekNext()
			if next.Type == kwUNIQUE {
				p.advance() // unexpected word
				return nil, p.syntaxErrorAtCur()
			}
			if next.Type == kwFOREIGN {
				p.advance() // unexpected word
				p.advance() // FOREIGN
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	_ = hasConstraintName

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		tc.Type = nodes.CONSTRAINT_PRIMARY
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr530 error
		tc.Columns, parseErr530 = p.parseUniqueIdentifierList()
		if parseErr530 != nil {
			return nil, parseErr530
		}
		if tc.Columns.Len() == 0 || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		tc.Type = nodes.CONSTRAINT_UNIQUE
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr531 error
		tc.Columns, parseErr531 = p.parseUniqueIdentifierList()
		if parseErr531 != nil {
			return nil, parseErr531
		}
		if tc.Columns.Len() == 0 || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

	case kwCHECK:
		checkTok := p.advance() // consume CHECK
		tc.Type = nodes.CONSTRAINT_CHECK
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtTok(checkTok)
		}
		p.advance()
		var parseErr532 error
		tc.Expr, parseErr532 = p.parseExpr()
		if parseErr532 != nil {
			return nil, parseErr532
		}
		if tc.Expr == nil || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		tc.Type = nodes.CONSTRAINT_FOREIGN
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		// A column, or an attribute of an object column (o.ref_attr), as in
		// the out-of-line REF constraint.
		var parseErr533 error
		tc.Columns, parseErr533 = p.parseColumnPathList()
		if parseErr533 != nil {
			return nil, parseErr533
		}
		if tc.Columns.Len() == 0 || p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwREFERENCES {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		var parseErr534 error
		tc.RefTable, parseErr534 = p.parseReferencedTable()
		if parseErr534 != nil {
			return nil, parseErr534
		}
		if p.cur.Type == '(' {
			p.advance()
			var parseErr535 error
			tc.RefColumns, parseErr535 = p.parseUniqueIdentifierList()
			if parseErr535 != nil {
				return nil, parseErr535
			}
			if tc.RefColumns.Len() == 0 || p.cur.Type != ')' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		}
		// ON DELETE
		if allowOnDelete && p.cur.Type == kwON {
			p.advance() // consume ON
			if err := p.expectCreateTableToken(kwDELETE); err != nil {
				return nil, err
			}
			var parseErr536 error
			tc.OnDelete, parseErr536 = p.parseDeleteAction()
			if parseErr536 != nil {
				return nil, parseErr536
			}
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	tc.Loc.End = p.prev.End
	return tc, nil
}

// parseIdentifierList parses a comma-separated list of identifiers,
// returning a *List of *String nodes.
func (p *Parser) parseIdentifierList() (*nodes.List, error) {
	return p.parseIdentifierListUnique(false)
}

func (p *Parser) parseUniqueIdentifierList() (*nodes.List, error) {
	return p.parseIdentifierListUnique(true)
}

func (p *Parser) parseIdentifierListUnique(unique bool) (*nodes.List, error) {
	list := &nodes.List{}
	var seen map[string]bool
	if unique {
		seen = map[string]bool{}
	}
	for {
		tok := p.cur
		name, parseErr537 := p.parseIdentifier()
		if parseErr537 != nil {
			return nil, parseErr537
		}
		if name == "" {
			// A name follows each ',' (ORA-00936 at a ')' after one).
			if len(list.Items) > 0 {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		if unique {
			if seen[name] {
				return nil, p.syntaxErrorAtTok(tok)
			}
			seen[name] = true
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
}

// parseColumnPathList parses a comma-separated list of columns, each
// possibly a dotted attribute path (o.r), as *String nodes.
func (p *Parser) parseColumnPathList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		column, err := p.parseColumnPath()
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, column)
		if p.cur.Type != ',' {
			return list, nil
		}
		p.advance() // consume ','
	}
}

// parseColumnPath parses a column, or an attribute of an object column
// (o.r), whose names it records in p.attributePaths, as a *String of the
// names joined by dots.
func (p *Parser) parseColumnPath() (*nodes.String, error) {
	var parts []string
	for {
		part, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		if part == "" {
			return nil, p.syntaxErrorAtCur()
		}
		parts = append(parts, part)
		if p.cur.Type != '.' {
			break
		}
		p.advance()
	}
	column := &nodes.String{Str: strings.Join(parts, ".")}
	if len(parts) > 1 {
		if p.attributePaths == nil {
			p.attributePaths = map[*nodes.String][]string{}
		}
		p.attributePaths[column] = parts
	}
	return column, nil
}

// parseReferencedTable parses the table a referential constraint
// references. A referential constraint is DDL on that table, which Oracle
// 23ai never runs over a database link (ORA-02021).
func (p *Parser) parseReferencedTable() (*nodes.ObjectName, error) {
	table, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	// A schema names its table (ORA-00903 for REFERENCES s.).
	if table == nil || table.Name == "" || p.prev.Type == '.' {
		return nil, p.syntaxErrorAtCur()
	}
	// The table is local, with no database link and no bare @ (ORA-02021
	// for REFERENCES o@ on Oracle 23ai, reported at '@').
	if table.DBLink != "" || p.prev.Type == '@' {
		inQuote := false
		for i := table.Loc.Start; i < table.Loc.End; i++ {
			switch p.source[i] {
			case '"':
				inQuote = !inQuote
			case '@':
				if !inQuote {
					return nil, &ParseError{
						Message:  `syntax error at or near "@"`,
						Position: i,
					}
				}
			}
		}
		return nil, p.syntaxErrorAtNode(table)
	}
	return table, nil
}

// parseIdentifierListAsStrings parses a comma-separated list of identifiers,
// returning a *List of *String nodes. Same as parseIdentifierList.
func (p *Parser) parseIdentifierListAsStrings() (*nodes.List, error) {
	return p.parseIdentifierList()
}

// parseDeleteAction parses the ON DELETE action (CASCADE, SET NULL, etc.).
func (p *Parser) parseDeleteAction() (string, error) {
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		return "CASCADE", nil
	case kwSET:
		p.advance() // consume SET
		if err := p.expectCreateTableToken(kwNULL); err != nil {
			return "", err
		}
		return "SET NULL", nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseTableOptions parses table-level options after the column definitions.
// These include physical_properties and table_properties from the BNF.
//
// BNF (CREATE-TABLE.bnf lines 200-351):
//
//	physical_properties: deferred_segment_creation | segment_attributes_clause |
//	    table_compression | heap_org_table_clause | index_org_table_clause | external_table_clause
//	table_properties: column_properties | read_only_clause | indexing_clause |
//	    table_partitioning_clauses | CACHE/NOCACHE | RESULT_CACHE |
//	    parallel_clause | ROWDEPENDENCIES | enable_disable_clause |
//	    row_movement_clause | flashback_archive_clause | AS subquery
func (p *Parser) parseTableOptions(stmt *nodes.CreateTableStmt, hasColumns bool) error {
	inmemoryTable := false
	for {
		switch p.cur.Type {
		case kwTABLESPACE:
			p.advance() // consume TABLESPACE
			var parseErr538 error
			stmt.Tablespace, parseErr538 = p.parseIdentifier()
			if parseErr538 != nil {
				return parseErr538
			}
			if stmt.Tablespace == "" {
				return p.syntaxErrorAtCur()
			}

		case kwON:

			next := p.peekNext()
			if next.Type == kwCOMMIT {
				// ON COMMIT is a temporary table's: Oracle 23ai refuses it
				// on any other table (ORA-00922 at ON).
				if !stmt.Global && !stmt.Private {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume ON
				p.advance()
				var // consume COMMIT
				parseErr539 error
				stmt.OnCommit, parseErr539 = p.parseOnCommitAction()
				if parseErr539 != nil {
					return parseErr539
				}
			} else {
				return nil
			}

		case kwPARALLEL:
			p.advance()
			stmt.Parallel = "PARALLEL"
			// optional integer degree
			if p.cur.Type == tokICONST {
				p.advance()
			}

		case kwNOPARALLEL:
			p.advance()
			stmt.Parallel = "NOPARALLEL"

		case kwCOMPRESS:
			p.createTableIndexOrganized = stmt.Organization == "INDEX"
			err := p.parseCreateTablePartitionCompression()
			p.createTableIndexOrganized = false
			if err != nil {
				return err
			}
			stmt.Compress = "COMPRESS"

		case kwNOCOMPRESS:
			p.advance()
			stmt.Compress = "NOCOMPRESS"

		case kwPARTITION:
			var parseErr540 error
			stmt.Partition, parseErr540 = p.parsePartitionClause(stmt, hasColumns)
			if parseErr540 != nil {
				return parseErr540
			}

		case kwLOGGING:
			p.advance()
			stmt.Logging = "LOGGING"

		case kwNOLOGGING:
			p.advance()
			stmt.Logging = "NOLOGGING"

		case kwCACHE:
			p.advance()
			stmt.Cache = "CACHE"

		case kwRESULT_CACHE:
			// RESULT_CACHE ( MODE { DEFAULT | FORCE } )
			p.advance() // consume RESULT_CACHE
			if err := p.expectCreateTableToken('('); err != nil {
				return err
			}
			if err := p.expectCreateTableWord("MODE"); err != nil {
				return err
			}
			switch {
			case p.cur.Type == kwDEFAULT:
				stmt.ResultCache = "DEFAULT"
			case p.isKeywordStr("FORCE"):
				stmt.ResultCache = "FORCE"
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.expectCreateTableToken(')'); err != nil {
				return err
			}

		case kwFLASHBACK:
			// FLASHBACK ARCHIVE [ flashback_archive ]
			p.advance() // consume FLASHBACK
			if err := p.expectCreateTableWord("ARCHIVE"); err != nil {
				return err
			}
			{
				if p.isIdentLike() && !p.isStatementEnd() && !p.isCreateTableOptionBoundary() {
					var parseErr541 error
					stmt.FlashbackArchive, parseErr541 = p.parseIdentifier()
					if parseErr541 != nil {
						return parseErr541
					}
				} else {
					stmt.FlashbackArchive = "FLASHBACK ARCHIVE"
				}
			}

		case kwROW:
			// ROW STORE [ COMPRESS [ BASIC | ADVANCED ] ]
			p.advance() // consume ROW
			if err := p.expectCreateTableWord("STORE"); err != nil {
				return err
			}
			if p.cur.Type == kwCOMPRESS {
				p.advance() // consume COMPRESS
				stmt.Compress = "ROW STORE COMPRESS"
				if p.isIdentLikeStr("BASIC") || p.isIdentLikeStr("ADVANCED") {
					stmt.Compress += " " + p.cur.Str
					p.advance()
				}
			}

		case kwREAD:
			// READ ONLY | READ WRITE
			p.advance() // consume READ
			if p.isIdentLikeStr("ONLY") {
				p.advance() // consume ONLY
				stmt.ReadOnly = "READ ONLY"
			} else if p.cur.Type == kwWRITE {
				p.advance() // consume WRITE
				stmt.ReadOnly = "READ WRITE"
			} else {
				return p.syntaxErrorAtCur()
			}

		case kwJSON:
			optStart := p.pos()
			p.advance() // consume JSON
			if err := p.parseCreateTableColumnNames(); err != nil {
				return err
			}
			if p.isKeywordStr("STORE") {
				p.advance()
			}
			if err := p.expectCreateTableToken(kwAS); err != nil {
				return err
			}
			if p.cur.Type == '(' {
				if err := p.parseCreateTableLobStorageParameters(); err != nil {
					return err
				}
			} else if p.isCreateTableNameStart() && !p.isCreateTableOptionBoundary() {
				p.advance() // consume LOB_segname
				if p.cur.Type == '(' {
					if err := p.parseCreateTableLobStorageParameters(); err != nil {
						return err
					}
				}
			}
			appendCreateTableOption(stmt, "JSON", p.createTableOptionValue(optStart+len("JSON")), nodes.Loc{Start: optStart, End: p.prev.End})

		case kwENABLE:
			// ENABLE ROW MOVEMENT, ENABLE LOGICAL REPLICATION, or ENABLE constraint
			next := p.peekNext()
			if next.Type == kwROW {
				p.advance() // consume ENABLE
				p.advance() // consume ROW
				if err := p.expectCreateTableWord("MOVEMENT"); err != nil {
					return err
				}
				stmt.RowMovement = "ENABLE"
			} else if isCreateTableWordTok(next, "LOGICAL") {
				optStart := p.pos()
				p.advance() // consume ENABLE
				p.advance() // consume LOGICAL
				if err := p.expectCreateTableWord("REPLICATION"); err != nil {
					return err
				}
				appendCreateTableOption(stmt, "ENABLE LOGICAL REPLICATION", "", nodes.Loc{Start: optStart, End: p.prev.End})
			} else {
				p.createTableEngineCheck = true
				parseErr542 :=
					// ENABLE [VALIDATE|NOVALIDATE] constraint_clause - skip it
					p.parseEnableDisableClause(stmt)
				if parseErr542 != nil {
					return parseErr542
				}
			}

		case kwDISABLE:
			next := p.peekNext()
			if next.Type == kwROW {
				p.advance() // consume DISABLE
				p.advance() // consume ROW
				if err := p.expectCreateTableWord("MOVEMENT"); err != nil {
					return err
				}
				stmt.RowMovement = "DISABLE"
			} else if isCreateTableWordTok(next, "LOGICAL") {
				optStart := p.pos()
				p.advance() // consume DISABLE
				p.advance() // consume LOGICAL
				if err := p.expectCreateTableWord("REPLICATION"); err != nil {
					return err
				}
				appendCreateTableOption(stmt, "DISABLE LOGICAL REPLICATION", "", nodes.Loc{Start: optStart, End: p.prev.End})
			} else {
				p.createTableEngineCheck = true
				parseErr543 := p.parseEnableDisableClause(stmt)
				if parseErr543 != nil {
					return parseErr543

					// AS subquery (for CTAS after options)
				}
			}

		case kwAS:

			p.advance()
			var // consume AS
			parseErr544 error
			stmt.AsQuery, parseErr544 = p.parseSelectStmt()
			if parseErr544 != nil {
				return parseErr544
			}
			if !selectStmtHasTargets(stmt.AsQuery) {
				return p.syntaxErrorAtCur()
			}
			return nil

		default:
			if p.isIdentLike() {
				switch p.cur.Str {
				case "NOCACHE":
					p.advance()
					stmt.Cache = "NOCACHE"
				case "MONITORING", "NOMONITORING":
					p.advance()
				case "SEGMENT":
					// SEGMENT CREATION { IMMEDIATE | DEFERRED }
					p.advance() // consume SEGMENT
					if err := p.expectCreateTableWord("CREATION"); err != nil {
						return err
					}
					switch p.cur.Type {
					case kwIMMEDIATE:
						stmt.SegmentCreation = "IMMEDIATE"
					case kwDEFERRED:
						stmt.SegmentCreation = "DEFERRED"
					default:
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "ORGANIZATION":
					if stmt.Organization != "" {
						return p.syntaxErrorAtCur()
					}
					// ORGANIZATION { HEAP | INDEX | EXTERNAL }
					p.advance() // consume ORGANIZATION
					if p.isIdentLikeStr("HEAP") {
						stmt.Organization = "HEAP"
						p.advance()
					} else if p.cur.Type == kwINDEX {
						indexTok := p.cur
						stmt.Organization = "INDEX"
						p.advance()
						// An index-organized table must define a PRIMARY KEY
						// constraint (ORA-25175 at INDEX).
						if !createTableHasPrimaryKey(stmt) {
							return p.syntaxErrorAtTok(indexTok)
						}
					} else if p.isIdentLikeStr("EXTERNAL") {
						optStart := p.pos()
						stmt.Organization = "EXTERNAL"
						p.advance()
						if err := p.parseCreateTableExternalData(hasColumns, stmt.Name.Loc); err != nil {
							return err
						}
						appendCreateTableOption(stmt, "ORGANIZATION EXTERNAL", p.createTableOptionValue(optStart+len("EXTERNAL")), nodes.Loc{Start: optStart, End: p.prev.End})
					}
				case "INDEXING":
					// INDEXING { ON | OFF }
					p.advance() // consume INDEXING
					switch {
					case p.cur.Type == kwON:
						stmt.Indexing = "ON"
					case p.isKeywordStr("OFF"):
						stmt.Indexing = "OFF"
					default:
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "ROWDEPENDENCIES":
					p.advance()
					stmt.RowDependencies = "ROWDEPENDENCIES"
				case "NOROWDEPENDENCIES":
					p.advance()
					stmt.RowDependencies = "NOROWDEPENDENCIES"
				case "PCTFREE":
					p.advance() // consume PCTFREE
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "PCTUSED":
					p.advance() // consume PCTUSED
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "INITRANS":
					p.advance() // consume INITRANS
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "MAXTRANS":
					p.advance() // consume MAXTRANS
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "PCTTHRESHOLD":
					// PCTTHRESHOLD integer — index-organized tables only;
					// Oracle rejects it on heap tables (ORA-00922).
					if stmt.Organization != "INDEX" {
						return p.syntaxErrorAtCur()
					}
					p.advance() // consume PCTTHRESHOLD
					if p.cur.Type != tokICONST {
						return p.syntaxErrorAtCur()
					}
					p.advance()
				case "COLUMN":
					// COLUMN STORE COMPRESS [ FOR { QUERY | ARCHIVE } [ LOW | HIGH ] ]
					// or object_type_col_properties: COLUMN column substitutable_column_clause
					if isCreateTableWordTok(p.peekNext(), "STORE") {
						p.advance() // consume COLUMN
						p.advance() // consume STORE
						if err := p.expectCreateTableToken(kwCOMPRESS); err != nil {
							return err
						}
						stmt.Compress = "COLUMN STORE COMPRESS"
						if p.cur.Type == kwFOR {
							p.advance() // consume FOR
							if !p.isKeywordStr("QUERY") && !p.isKeywordStr("ARCHIVE") {
								return p.syntaxErrorAtCur()
							}
							stmt.Compress += " FOR " + p.cur.Str
							p.advance()
							if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
								stmt.Compress += " " + p.cur.Str
								p.advance()
							}
						}
					} else {
						optStart := p.pos()
						if err := p.parseCreateTableObjectTypeColProperties(); err != nil {
							return err
						}
						appendCreateTableOption(stmt, "COLUMN", p.createTableOptionValue(optStart+len("COLUMN")), nodes.Loc{Start: optStart, End: p.prev.End})
					}
				case "INMEMORY":
					if err := p.parseCreateTableInmemory(stmt, hasColumns, &inmemoryTable, false); err != nil {
						return err
					}
				case "LOB":
					optStart := p.pos()
					if err := p.parseCreateTableLobStorageClause(); err != nil {
						return err
					}
					appendCreateTableOption(stmt, "LOB", p.createTableOptionValue(optStart+len("LOB")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "XMLTYPE":
					optStart := p.pos()
					if err := p.parseATXMLTypeColumnProperties(); err != nil {
						return err
					}
					appendCreateTableOption(stmt, "XMLTYPE", p.createTableOptionValue(optStart+len("XMLTYPE")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "VARRAY":
					optStart := p.pos()
					kwTok := p.cur
					if err := p.parseCreateTableVarrayColProperties(hasColumns); err != nil {
						if !hasColumns && p.cur.Type == tokEOF && p.prev.Type == kwAS {
							return p.syntaxErrorAtTok(kwTok)
						}
						return err
					}
					appendCreateTableOption(stmt, "VARRAY", p.createTableOptionValue(optStart+len("VARRAY")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "NESTED":
					optStart := p.pos()
					kwTok := p.cur
					if err := p.parseCreateTableNestedTableColProperties(hasColumns); err != nil {
						if !hasColumns && p.cur.Type == tokEOF && p.prev.Type == kwAS {
							return p.syntaxErrorAtTok(kwTok)
						}
						return err
					}
					appendCreateTableOption(stmt, "NESTED", p.createTableOptionValue(optStart+len("NESTED")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "NO":
					// NO FLASHBACK ARCHIVE / NO INMEMORY
					next := p.peekNext()
					if next.Type == kwFLASHBACK {
						p.advance() // consume NO
						p.advance() // consume FLASHBACK
						if err := p.expectCreateTableWord("ARCHIVE"); err != nil {
							return err
						}
						stmt.FlashbackArchive = "NO FLASHBACK ARCHIVE"
					} else if isCreateTableWordTok(next, "INMEMORY") {
						if err := p.parseCreateTableInmemory(stmt, hasColumns, &inmemoryTable, false); err != nil {
							return err
						}
					} else {
						return nil
					}
				case "STORAGE":
					optStart := p.pos()
					if err := p.parseCreateTableStorageClause(); err != nil {
						return err
					}
					appendCreateTableOption(stmt, "STORAGE", p.createTableOptionValue(optStart+len("STORAGE")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "MAPPING", "NOMAPPING":
					// { MAPPING TABLE | NOMAPPING } belongs to
					// index_org_table_clause: on any other table Oracle
					// 23ai answers ORA-00922 at the word after it.
					mapping := p.cur.Str == "MAPPING"
					p.advance()
					if mapping {
						if err := p.expectCreateTableToken(kwTABLE); err != nil {
							return err
						}
					}
					if stmt.Organization != "INDEX" {
						return p.syntaxErrorAtCur()
					}
				case "INCLUDING":
					// INCLUDING column_name OVERFLOW
					p.advance()
					parseDiscard552, parseErr551 := p.parseIdentifier()
					_ = parseDiscard552
					if parseErr551 != nil {
						return parseErr551
					}
				case "OVERFLOW":
					p.advance()
				case "REJECT":
					// REJECT LIMIT { integer | UNLIMITED }
					p.advance() // consume REJECT
					if p.isIdentLikeStr("LIMIT") {
						p.advance() // consume LIMIT
						if p.cur.Type == tokICONST {
							p.advance()
						} else if p.isIdentLikeStr("UNLIMITED") {
							p.advance()
						}
					}
				case "ILM":
					optStart := p.pos()
					if p.peekNext().Type != kwADD && p.peekNext().Type != kwDELETE &&
						p.peekNext().Type != kwENABLE && p.peekNext().Type != kwDISABLE {
						p.createTableEngineCheck = true
					}
					if err := p.parseCreateTableIlmClause(); err != nil {
						return err
					}
					appendCreateTableOption(stmt, "ILM", p.createTableOptionValue(optStart+len("ILM")), nodes.Loc{Start: optStart, End: p.prev.End})
				case "CLUSTERING":
					optStart := p.pos()
					p.advance()
					valueStart := p.pos()
					if err := p.parseCreateTableClustering(); err != nil {
						return err
					}
					appendCreateTableOption(stmt, "CLUSTERING", p.createTableOptionValue(valueStart), nodes.Loc{Start: optStart, End: p.prev.End})
				case "IMMEDIATE", "EXTERNAL":
					// Oracle advances past bare IMMEDIATE or EXTERNAL in table
					// options before refusing what follows (ORA-00922).
					p.advance()
					return p.syntaxErrorAtCur()
				case "SUPPLEMENTAL":
					optStart := p.pos()
					p.advance() // consume SUPPLEMENTAL
					if p.isIdentLikeStr("LOG") {
						p.advance() // consume LOG
					}
					value := p.collectCreateTableOptionValue()
					appendCreateTableOption(stmt, "SUPPLEMENTAL LOG", value, nodes.Loc{Start: optStart, End: p.prev.End})
				default:
					return nil
				}
			} else {
				return nil
			}
		}
	}
}

func appendCreateTableOption(stmt *nodes.CreateTableStmt, key, value string, loc nodes.Loc) {
	if loc.End <= loc.Start {
		return
	}
	if stmt.Options == nil {
		stmt.Options = &nodes.List{}
	}
	stmt.Options.Items = append(stmt.Options.Items, &nodes.DDLOption{Key: key, Value: value, Loc: loc})
}

func (p *Parser) collectCreateTableOptionValue() string {
	tokens := p.collectDDLTokensUntil(p.isCreateTableOptionBoundary)
	return strings.Join(tokens, " ")
}

func (p *Parser) isCreateTableOptionBoundary() bool {
	switch p.cur.Type {
	case ';', tokEOF, kwAS, kwCACHE, kwDISABLE, kwENABLE, kwFLASHBACK, kwJSON, kwLOGGING,
		kwNOLOGGING, kwON, kwPARALLEL, kwNOPARALLEL, kwPARTITION, kwREAD,
		kwRESULT_CACHE, kwTABLESPACE:
		return true
	}
	if p.isIdentLike() {
		switch p.cur.Str {
		case "CLUSTERING", "COLUMN", "FILESYSTEM_LIKE_LOGGING", "ILM", "INCLUDING",
			"INDEXING", "INMEMORY", "INITRANS", "LOB", "MAPPING", "MEMOPTIMIZE",
			"MONITORING", "NESTED", "NO", "NOCACHE", "NOMAPPING", "NOMONITORING",
			"NOROWDEPENDENCIES", "ORGANIZATION", "OVERFLOW", "PCTFREE", "PCTUSED",
			"REJECT", "ROWDEPENDENCIES", "SEGMENT", "STORAGE", "SUPPLEMENTAL",
			"VARRAY", "XMLTYPE":
			return true
		}
	}
	return false
}

func (p *Parser) collectDDLTokensUntil(isBoundary func() bool) []string {
	var tokens []string
	depth := 0
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if depth == 0 && len(tokens) > 0 && isBoundary() {
			break
		}
		tokens = append(tokens, p.ddlOptionTokenText(p.cur))
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

func (p *Parser) ddlOptionTokenText(tok Token) string {
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

// parseImmutableTableClauses parses the immutable_table_clauses of CREATE
// IMMUTABLE TABLE.
//
// BNF (CREATE-TABLE.bnf):
//
//	immutable_table_clauses:
//	    [ immutable_table_no_drop_clause ] [ immutable_table_no_delete_clause ]
//	immutable_table_no_drop_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	immutable_table_no_delete_clause: NO DELETE [ LOCKED | UNTIL integer DAYS AFTER INSERT ]
func (p *Parser) parseImmutableTableClauses(stmt *nodes.CreateTableStmt) error {
	if err := p.parseCreateTableNoDrop(stmt); err != nil {
		return err
	}
	return p.parseCreateTableNoDelete(stmt, false)
}

// parseCreateTableNoDrop reads NO DROP [ UNTIL integer DAYS IDLE ] when it
// is next.
func (p *Parser) parseCreateTableNoDrop(stmt *nodes.CreateTableStmt) error {
	if !p.isKeywordStr("NO") || p.peekNext().Type != kwDROP {
		return nil
	}
	p.advance() // consume NO
	p.advance() // consume DROP
	stmt.ImmutableNoDrop = "NO DROP"
	if p.isKeywordStr("UNTIL") {
		p.advance() // consume UNTIL
		days := p.cur.Str
		if err := p.parseCreateTableInteger(); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("DAYS"); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("IDLE"); err != nil {
			return err
		}
		stmt.ImmutableNoDrop = "NO DROP UNTIL " + days + " DAYS IDLE"
	}
	return nil
}

// parseCreateTableNoDelete reads NO DELETE [ LOCKED | UNTIL integer DAYS
// AFTER INSERT ] when it is next. A blockchain table's
// blockchain_row_retention_clause requires LOCKED or UNTIL.
func (p *Parser) parseCreateTableNoDelete(stmt *nodes.CreateTableStmt, retention bool) error {
	if !p.isKeywordStr("NO") || p.peekNext().Type != kwDELETE {
		return nil
	}
	p.advance() // consume NO
	p.advance() // consume DELETE
	stmt.ImmutableNoDel = "NO DELETE"
	switch {
	case p.isKeywordStr("LOCKED"):
		p.advance()
		stmt.ImmutableNoDel = "NO DELETE LOCKED"
	case p.isKeywordStr("UNTIL"):
		p.advance() // consume UNTIL
		days := p.cur.Str
		if err := p.parseCreateTableInteger(); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("DAYS"); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwAFTER); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwINSERT); err != nil {
			return err
		}
		stmt.ImmutableNoDel = "NO DELETE UNTIL " + days + " DAYS AFTER INSERT"
	case retention:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseBlockchainTableClauses parses the blockchain_table_clauses of CREATE
// BLOCKCHAIN TABLE.
//
// BNF (CREATE-TABLE.bnf):
//
//	blockchain_table_clauses:
//	    [ blockchain_drop_table_clause ]
//	    [ blockchain_row_retention_clause ]
//	    [ blockchain_hash_and_data_format_clause ]
//	blockchain_drop_table_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	blockchain_row_retention_clause: NO DELETE { LOCKED | UNTIL integer DAYS AFTER INSERT }
//	blockchain_hash_and_data_format_clause: HASHING USING 'hash_algorithm' VERSION 'version_string'
func (p *Parser) parseBlockchainTableClauses(stmt *nodes.CreateTableStmt) error {
	if err := p.parseCreateTableNoDrop(stmt); err != nil {
		return err
	}
	if err := p.parseCreateTableNoDelete(stmt, true); err != nil {
		return err
	}
	if !p.isKeywordStr("HASHING") {
		return nil
	}
	p.advance() // consume HASHING
	if err := p.expectCreateTableToken(kwUSING); err != nil {
		return err
	}
	stmt.BlockchainHash = p.cur.Str
	if err := p.expectCreateTableToken(tokSCONST); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("VERSION"); err != nil {
		return err
	}
	stmt.BlockchainVer = p.cur.Str
	return p.expectCreateTableToken(tokSCONST)
}

// parseEnableDisableClause parses ENABLE/DISABLE [VALIDATE|NOVALIDATE] constraint clause.
//
// BNF (CREATE-TABLE.bnf lines 463-472):
//
//	{ ENABLE | DISABLE } [ VALIDATE | NOVALIDATE ]
//	{ UNIQUE ( column [,...] ) | PRIMARY KEY | CONSTRAINT constraint_name }
//	[ using_index_clause ] [ exceptions_clause ] [ CASCADE ] [ { KEEP | DROP } INDEX ]
func (p *Parser) parseEnableDisableClause(stmt *nodes.CreateTableStmt) error {
	p.advance() // consume ENABLE or DISABLE

	// [ VALIDATE | NOVALIDATE ]
	if p.isIdentLikeStr("VALIDATE") || p.isIdentLikeStr("NOVALIDATE") {
		p.advance()
	}

	// { UNIQUE (cols) | PRIMARY KEY | CONSTRAINT name }
	switch p.cur.Type {
	case kwUNIQUE:
		p.advance()
		if err := p.parseCreateTableColumnNames(); err != nil {
			return err
		}
	case kwPRIMARY:
		p.advance()
		if err := p.expectCreateTableToken(kwKEY); err != nil {
			return err
		}
		// In CREATE TABLE, ENABLE/DISABLE PRIMARY KEY followed by any trailing
		// clause requires a PRIMARY KEY constraint defined on the table
		// (ORA-02432 at the token immediately following PRIMARY KEY).
		if !createTableHasPrimaryKey(stmt) && !p.isStatementEnd() {
			return p.syntaxErrorAtCur()
		}
	case kwCONSTRAINT:
		p.advance()
		constraintName := p.cur.Str
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		// In CREATE TABLE, ENABLE/DISABLE CONSTRAINT name requires that
		// constraint to be defined in the statement (ORA-02250).
		if !p.createTableHasConstraint(stmt, constraintName) {
			return p.syntaxErrorAtCur()
		}
	default:
		return p.syntaxErrorAtCur()
	}

	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
		// skip index properties
		if p.cur.Type == '(' {
			p.skipParenthesized()
		} else if p.isIdentLike() {
			parseDiscard556, parseErr555 := p.parseIdentifier()
			_ = parseDiscard556

			// [ EXCEPTIONS INTO table ]
			if parseErr555 != nil {
				return parseErr555
			}
		}
	}

	if p.isIdentLikeStr("EXCEPTIONS") {
		p.advance()
		if p.cur.Type == kwINTO {
			p.advance()
			parseDiscard558, parseErr557 := p.parseObjectName()
			_ = parseDiscard558

			// [ CASCADE ]
			if parseErr557 != nil {
				return parseErr557
			}
		}
	}

	if p.cur.Type == kwCASCADE {
		p.advance()
	}

	// [ { KEEP | DROP } INDEX ]
	if p.isIdentLikeStr("KEEP") || p.cur.Type == kwDROP {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
	}
	return nil
}

func createTableHasPrimaryKey(stmt *nodes.CreateTableStmt) bool {
	if stmt == nil {
		return false
	}
	if stmt.Constraints != nil {
		for _, item := range stmt.Constraints.Items {
			if tc, ok := item.(*nodes.TableConstraint); ok && tc.Type == nodes.CONSTRAINT_PRIMARY {
				return true
			}
		}
	}
	if stmt.Columns != nil {
		for _, item := range stmt.Columns.Items {
			col, ok := item.(*nodes.ColumnDef)
			if !ok || col.Constraints == nil {
				continue
			}
			for _, citem := range col.Constraints.Items {
				if cc, ok := citem.(*nodes.ColumnConstraint); ok && cc.Type == nodes.CONSTRAINT_PRIMARY {
					return true
				}
			}
		}
	}
	return false
}

func (p *Parser) createTableHasConstraint(stmt *nodes.CreateTableStmt, name string) bool {
	if stmt == nil || name == "" {
		return false
	}
	if p.createTableConstraintNames[name] {
		return true
	}
	if stmt.Constraints != nil {
		for _, item := range stmt.Constraints.Items {
			if tc, ok := item.(*nodes.TableConstraint); ok && tc.Name == name {
				return true
			}
		}
	}
	if stmt.Columns != nil {
		for _, item := range stmt.Columns.Items {
			col, ok := item.(*nodes.ColumnDef)
			if !ok || col.Constraints == nil {
				continue
			}
			for _, citem := range col.Constraints.Items {
				if cc, ok := citem.(*nodes.ColumnConstraint); ok && cc.Name == name {
					return true
				}
			}
		}
	}
	return false
}

// isStatementEnd returns true if the current token is at a statement boundary.
func (p *Parser) isStatementEnd() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF
}

// parseOnCommitAction parses the action after ON COMMIT.
//
//	PRESERVE ROWS | DELETE ROWS
func (p *Parser) parseOnCommitAction() (string, error) {
	var action string
	switch {
	case p.isKeywordStr("PRESERVE"):
		action = "PRESERVE ROWS"
	case p.cur.Type == kwDELETE:
		action = "DELETE ROWS"
	default:
		return "", p.syntaxErrorAtCur()
	}
	p.advance() // consume PRESERVE / DELETE
	if err := p.expectCreateTableToken(kwROWS); err != nil {
		return "", err
	}
	return action, nil
}

// statementEndsAfterNext reports whether the statement's text ends right
// after the current token.
func (p *Parser) statementEndsAfterNext() bool {
	next := p.peekNext()
	return next.Type == ';' || next.Type == tokEOF
}
