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

	// OR REPLACE | IF NOT EXISTS
	orReplace := false
	ifNotExists := false
	if p.cur.Type == kwOR {
		p.advance() // consume OR
		if p.cur.Type == kwREPLACE {
			p.advance() // consume REPLACE
			orReplace = true
		}
	} else if p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
				ifNotExists = true
			}
		}
	}

	// EDITIONABLE | NONEDITIONABLE (only for PL/SQL objects: PROCEDURE, FUNCTION, PACKAGE, TRIGGER, TYPE)
	// Note: Views handle EDITIONABLE internally in parseCreateViewStmt.
	editionable := false
	nonEditionable := false
	if p.isIdentLikeStr("EDITIONABLE") {
		// Peek to see if the next token is a PL/SQL object keyword
		next := p.peekNext()
		if next.Type == kwPROCEDURE || next.Type == kwFUNCTION || next.Type == kwPACKAGE ||
			next.Type == kwTRIGGER || next.Type == kwTYPE {
			editionable = true
			p.advance()
		}
	} else if p.isIdentLikeStr("NONEDITIONABLE") {
		next := p.peekNext()
		if next.Type == kwPROCEDURE || next.Type == kwFUNCTION || next.Type == kwPACKAGE ||
			next.Type == kwTRIGGER || next.Type == kwTYPE {
			nonEditionable = true
			p.advance()
		}
	}

	// GLOBAL TEMPORARY
	global := false
	private := false
	if p.cur.Type == kwGLOBAL {
		p.advance() // consume GLOBAL
		if p.cur.Type == kwTEMPORARY {
			p.advance() // consume TEMPORARY
			global = true
		}
	}

	// PRIVATE TEMPORARY
	if p.cur.Type == kwPRIVATE {
		p.advance() // consume PRIVATE
		if p.cur.Type == kwTEMPORARY {
			p.advance() // consume TEMPORARY
			private = true
		}
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

	// PUBLIC (for CREATE PUBLIC SYNONYM / DATABASE LINK)
	public := false
	if p.cur.Type == kwPUBLIC {
		public = true
		p.advance()
	}

	// [ AND { RESOLVE | COMPILE } ] [ NOFORCE ] — only for CREATE [OR REPLACE] JAVA
	// These modifiers appear before the JAVA keyword. Consume them and dispatch.
	if p.cur.Type == kwAND {
		p.advance() // consume AND
		if p.isIdentLike() && (p.cur.Str == "RESOLVE" || p.cur.Str == "COMPILE") {
			p.advance()
		}
	}
	if p.isIdentLike() && p.cur.Str == "NOFORCE" {
		next := p.peekNext()
		if next.Type == kwJAVA {
			p.advance() // consume NOFORCE
		}
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, global, private, sharded, duplicated)
	case kwUNIQUE, kwBITMAP, kwINDEX:
		return p.parseCreateIndexStmt(start)
	case kwMATERIALIZED:
		return p.parseCreateMaterializedOrView(start, orReplace)
	case kwFORCE, kwVIEW:
		return p.parseCreateViewStmt(start, orReplace)
	case kwSEQUENCE:
		return p.parseCreateSequenceStmt(start)
	case kwSYNONYM:
		return p.parseCreateSynonymStmt(start, orReplace, public)
	case kwDATABASE:
		// Distinguish CREATE DATABASE LINK from CREATE DATABASE
		next := p.peekNext()
		if next.Type == kwLINK {
			return p.parseCreateDatabaseLinkStmt(start, public)
		}
		return p.parseCreateDatabaseStmt(start)
	case kwTYPE:
		return p.parseCreateTypeStmt(start, orReplace, ifNotExists, editionable, nonEditionable)
	case kwPROCEDURE:
		return p.parseCreateProcedureStmt(start, orReplace, ifNotExists, editionable, nonEditionable)
	case kwFUNCTION:
		return p.parseCreateFunctionStmt(start, orReplace, ifNotExists, editionable, nonEditionable)
	case kwPACKAGE:
		return p.parseCreatePackageStmt(start, orReplace, ifNotExists, editionable, nonEditionable)
	case kwTRIGGER:
		return p.parseCreateTriggerStmt(start, orReplace, ifNotExists, editionable, nonEditionable)
	case kwAUDIT:
		// CREATE AUDIT POLICY
		p.advance() // consume AUDIT
		if p.isIdentLikeStr("POLICY") {
			p.advance() // consume POLICY
		}
		return p.parseCreateAuditPolicyStmt(start)
	case kwJSON:
		// CREATE JSON RELATIONAL DUALITY VIEW
		p.advance() // consume JSON
		if p.isIdentLike() && p.cur.Str == "RELATIONAL" {
			p.advance() // consume RELATIONAL
		}
		if p.isIdentLike() && p.cur.Str == "DUALITY" {
			p.advance() // consume DUALITY
		}
		if p.cur.Type == kwVIEW {
			p.advance() // consume VIEW
		}
		return p.parseCreateJsonDualityViewStmt(start, orReplace, false)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwDIRECTORY, kwCONTEXT,
		kwCLUSTER, kwJAVA, kwLIBRARY, kwSCHEMA:
		return p.parseCreateAdminObject(start, orReplace)
	case kwTEMPORARY:
		// CREATE TEMPORARY TABLESPACE (standalone, not GLOBAL TEMPORARY TABLE)
		if !global && !private {
			p.advance() // consume TEMPORARY
			if p.cur.Type == kwTABLESPACE {
				p.advance() // consume TABLESPACE
				return p.parseCreateTablespaceStmt(start, false, false, false, true, false)
			}
		}
		return nil, nil
	default:
		// Check for "NO FORCE VIEW"
		if p.isIdentLikeStr("NO") || p.cur.Type == kwNOT {
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for EDITIONING / EDITIONABLE / NONEDITIONABLE VIEW
		if p.isIdentLike() && (p.cur.Str == "EDITIONING" || p.cur.Str == "EDITIONABLE" || p.cur.Str == "NONEDITIONABLE") {
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for BIGFILE/SMALLFILE/UNDO TABLESPACE, CONTROLFILE
		if p.isIdentLike() {
			switch p.cur.Str {
			case "BIGFILE":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, true, false, false, false, false)
				}
			case "SMALLFILE":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, false, true, false, false, false)
				}
			case "UNDO":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, false, false, false, false, true)
				}
			case "LOCAL":
				// CREATE LOCAL TEMPORARY TABLESPACE
				p.advance() // consume LOCAL
				if p.cur.Type == kwTEMPORARY {
					p.advance() // consume TEMPORARY
					if p.cur.Type == kwTABLESPACE {
						p.advance() // consume TABLESPACE
						return p.parseCreateTablespaceStmt(start, false, false, true, true, false)
					}
				}
			case "CONTROLFILE":
				p.advance() // consume CONTROLFILE
				return p.parseCreateControlfileStmt(start)
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
				if p.cur.Type == kwVIEW {
					p.advance() // consume VIEW
				}
				return p.parseCreateAnalyticViewStmt(start, orReplace, false, false)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if p.isIdentLike() && p.cur.Str == "DIMENSION" {
					p.advance() // consume DIMENSION
				}
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
				if p.isIdentLike() && p.cur.Str == "GRAPH" {
					p.advance() // consume GRAPH
				}
				return p.parseCreatePropertyGraphStmt(start, orReplace, ifNotExists)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type == kwINDEX {
					p.advance() // consume INDEX
				}
				return p.parseCreateVectorIndexStmt(start, ifNotExists)
			case "ROLLBACK":
				// CREATE [PUBLIC] ROLLBACK SEGMENT
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
				}
				return p.parseCreateRollbackSegmentStmt(start, public)
			case "PURE":
				// CREATE [PURE] MLE ENV
				p.advance() // consume PURE
				if p.isIdentLike() && p.cur.Str == "MLE" {
					p.advance() // consume MLE
					if p.isIdentLike() && p.cur.Str == "ENV" {
						p.advance() // consume ENV
						return p.parseCreateMLEEnvStmt(start, orReplace)
					}
				}
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
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		OrReplace:   orReplace,
		Global:      global,
		Private:     private,
		Sharded:     sharded,
		Duplicated:  duplicated,
		Columns:     &nodes.List{},
		Constraints: &nodes.List{},
		Loc:         nodes.Loc{Start: start},
	}

	// IF NOT EXISTS (Oracle 23c)
	if p.cur.Type == kwIF {
		p.advance() // consume IF
		if p.cur.Type == kwNOT {
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
				stmt.IfNotExists = true
			}
		}
	}

	// Table name
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	var parseErr487 error
	stmt.Name, parseErr487 = p.parseReservedCheckedObjectName()
	if parseErr487 !=

		// SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
		nil {
		return nil, parseErr487
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume '='
		}
		switch {
		case p.isIdentLikeStr("METADATA"):
			stmt.Sharing = "METADATA"
			p.advance()
		case p.isIdentLikeStr("DATA"):
			stmt.Sharing = "DATA"
			p.advance()
		case p.isIdentLikeStr("EXTENDED"):
			p.advance() // consume EXTENDED
			if p.isIdentLikeStr("DATA") {
				p.advance() // consume DATA
			}
			stmt.Sharing = "EXTENDED DATA"
		case p.isIdentLikeStr("NONE"):
			stmt.Sharing = "NONE"
			p.advance()
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
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// Column definitions and table constraints
	if p.cur.Type == '(' {
		p.advance()
		parseErr489 := // consume '('
			p.parseColumnDefsAndConstraints(stmt)
		if parseErr489 != nil {
			return nil, parseErr489
		}
		if p.cur.Type == ')' {
			p.advance() // consume ')'
		}
	}
	parseErr490 :=

		// Immutable table clauses (before table options)
		// BNF (lines 733-741): immutable_table_clauses
		p.parseImmutableTableClauses(stmt)
	if parseErr490 !=

		// Blockchain table clauses
		// BNF (lines 743-756): blockchain_table_clauses
		nil {
		return nil, parseErr490
	}
	parseErr491 := p.parseBlockchainTableClauses(stmt)
	if parseErr491 !=

		// DEFAULT COLLATION collation_name
		nil {
		return nil, parseErr491
	}

	if p.cur.Type == kwDEFAULT {
		next := p.peekNext()
		if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "COLLATION" {
			p.advance() // consume DEFAULT
			p.advance()
			var // consume COLLATION
			parseErr492 error
			stmt.Collation, parseErr492 = p.parseIdentifier()
			if parseErr492 !=

				// Table options (physical_properties + table_properties)
				nil {
				return nil, parseErr492
			}
		}
	}
	parseErr493 := p.parseTableOptions(stmt)
	if parseErr493 !=

		// MEMOPTIMIZE FOR READ
		nil {
		return nil, parseErr493
	}

	if p.isIdentLikeStr("MEMOPTIMIZE") {
		p.advance() // consume MEMOPTIMIZE
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			if p.cur.Type == kwREAD {
				p.advance() // consume READ
				stmt.MemoptimizeRead = true
			} else if p.cur.Type == kwWRITE {
				p.advance() // consume WRITE
				stmt.MemoptimizeWrite = true
			}
		}
	}
	// MEMOPTIMIZE FOR WRITE (can appear after READ)
	if p.isIdentLikeStr("MEMOPTIMIZE") {
		p.advance() // consume MEMOPTIMIZE
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			if p.cur.Type == kwWRITE {
				p.advance() // consume WRITE
				stmt.MemoptimizeWrite = true
			} else if p.cur.Type == kwREAD {
				p.advance() // consume READ
				stmt.MemoptimizeRead = true
			}
		}
	}

	// PARENT [ schema. ] table
	if p.isIdentLikeStr("PARENT") {
		p.advance()
		var // consume PARENT
		parseErr494 error
		stmt.Parent, parseErr494 = p.parseObjectName()
		if parseErr494 != nil {
			return nil, parseErr494
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseColumnDefsAndConstraints parses the contents inside the parentheses of a CREATE TABLE.
// It handles both column definitions and table-level constraints.
func (p *Parser) parseColumnDefsAndConstraints(stmt *nodes.CreateTableStmt) error {
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Check if this is a table-level constraint
		if p.isTableConstraintStart() {
			tc, parseErr495 := p.parseTableConstraint()
			if parseErr495 != nil {
				return parseErr495
			}
			if tc != nil {
				stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
			}
		} else {
			// Column definition
			col, parseErr496 := p.parseColumnDef()
			if parseErr496 != nil {
				return parseErr496
			}
			if col != nil {
				stmt.Columns.Items = append(stmt.Columns.Items, col)
			}
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return nil
}

// isTableConstraintStart checks if the current position starts a table-level constraint.
func (p *Parser) isTableConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT:
		return true
	case kwPRIMARY:
		return true
	case kwFOREIGN:
		return true
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

	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return nil, err
	}
	var parseErr497 error
	col.Name, parseErr497 = p.parseIdentifier()
	if parseErr497 != nil {
		return nil, parseErr497
	}
	if col.Name == "" {
		return nil, nil
	}

	// DOMAIN [ schema. ] domain_name  (instead of datatype)
	if p.isIdentLikeStr("DOMAIN") {
		p.advance()
		var // consume DOMAIN
		parseErr498 error
		col.Domain, parseErr498 = p.parseObjectName()
		if parseErr498 != nil {
			return nil, parseErr498

			// Data type (optional for some Oracle column types, but typical)
		}
	} else if p.isTypeName() {
		var parseErr499 error

		col.TypeName, parseErr499 = p.parseTypeName()
		if parseErr499 !=

			// Column properties: SORT, VISIBLE, DEFAULT, NOT NULL, NULL, constraints, etc.
			nil {
			return nil, parseErr499
		}
	}
	parseErr500 := p.parseColumnProperties(col)
	if parseErr500 != nil {
		return nil, parseErr500
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// isTypeName returns true if the current token can begin a type name.
func (p *Parser) isTypeName() bool {
	switch p.cur.Type {
	case kwNUMBER, kwINTEGER, kwSMALLINT, kwDECIMAL, kwFLOAT,
		kwCHAR, kwVARCHAR2, kwVARCHAR, kwNCHAR, kwNVARCHAR2,
		kwCLOB, kwBLOB, kwNCLOB,
		kwDATE, kwTIMESTAMP, kwINTERVAL,
		kwRAW, kwLONG, kwROWID:
		return true
	}
	// User-defined type: identifier that is NOT a constraint keyword
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		return true
	}
	return false
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
				next := p.peekNext()
				if next.Type == kwNULL {
					p.advance() // consume ON
					p.advance() // consume NULL
					col.DefaultOnNull = true
					// FOR INSERT ONLY
					if p.cur.Type == kwFOR {
						p.advance() // consume FOR
						if p.cur.Type == kwINSERT {
							p.advance() // consume INSERT
							if p.isIdentLikeStr("ONLY") {
								p.advance() // consume ONLY
							}
							col.DefaultOnNullInsertOnly = true
						}
					}
				}
			}
			var parseErr501 error
			col.Default, parseErr501 = p.parseExpr()
			if parseErr501 != nil {
				return parseErr501
			}
			if col.Default == nil {

				return p.syntaxErrorAtCur()
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
				if p.cur.Type == kwDEFAULT {
					p.advance() // consume DEFAULT
					byDefault = true
					// ON NULL
					if p.cur.Type == kwON {
						next := p.peekNext()
						if next.Type == kwNULL {
							p.advance() // consume ON
							p.advance() // consume NULL
						}
					}
				}
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
					// ( identity_options )
					if p.cur.Type == '(' {
						p.advance()
						parseErr502 := p.parseIdentityOptions(identity)
						if parseErr502 != nil {
							return parseErr502
						}
						if p.cur.Type == ')' {
							p.advance()
						}
					}
					identity.Loc.End = p.prev.End
					col.Identity = identity
				} else if p.cur.Type == '(' {
					// virtual column: AS (expr) [VIRTUAL]
					p.advance()
					var // consume '('
					parseErr503 error
					col.Virtual, parseErr503 = p.parseExpr()
					if parseErr503 != nil {
						return parseErr503
					}
					if p.cur.Type == ')' {
						p.advance()
					}
					if p.cur.Type == kwVIRTUAL {
						p.advance() // consume VIRTUAL
					}
				}
			}
			_ = byDefault // stored implicitly: !always means BY DEFAULT

		case kwNOT:
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // consume NOT
				p.advance() // consume NULL
				col.NotNull = true
			} else {
				return nil
			}

		case kwNULL:
			p.advance() // consume NULL
			col.Null = true

		case kwINVISIBLE:
			p.advance()
			col.Invisible = true

		case kwCONSTRAINT:
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
					p.advance()
					var // consume COLLATE
					parseErr510 error
					col.Collation, parseErr510 = p.parseIdentifier()
					if parseErr510 != nil {
						return parseErr510
					}
				default:
					return nil
				}
			} else {
				return nil
			}
		}
	}
	return nil
}

// parseEncryptionSpec parses the encryption_spec after ENCRYPT keyword.
//
// BNF (CREATE-TABLE.bnf lines 95-98):
//
//	[ USING 'encrypt_algorithm' ]
//	[ IDENTIFIED BY password ]
//	[ SALT | NO SALT ]
func (p *Parser) parseEncryptionSpec(col *nodes.ColumnDef) error {
	for {
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
			if p.cur.Type == tokSCONST {
				col.Encrypt = "ENCRYPT USING " + p.cur.Str
				p.advance()
			}
		} else if p.isIdentLikeStr("IDENTIFIED") {
			p.advance() // consume IDENTIFIED
			if p.cur.Type == kwBY {
				p.advance()
				parseDiscard512, // consume BY
					parseErr511 := p.parseIdentifier()
				_ = // consume password
					parseDiscard512
				if parseErr511 != nil {
					return parseErr511
				}
			}
		} else if p.isIdentLikeStr("SALT") {
			p.advance() // consume SALT
		} else if p.isIdentLikeStr("NO") {
			next := p.peekNext()
			if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "SALT" {
				p.advance() // consume NO
				p.advance() // consume SALT
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
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
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if !p.isIdentLike() {
			return nil
		}
		switch p.cur.Str {
		case "START":
			p.advance() // consume START
			if p.cur.Type == kwWITH {
				p.advance() // consume WITH
			}
			if p.isIdentLikeStr("LIMIT") {
				p.advance() // consume LIMIT
				if p.isIdentLike() {
					p.advance() // consume VALUE
				}
			} else {
				parseDiscard514, parseErr513 := p.parseExpr()
				_ = // consume integer
					parseDiscard514
				if parseErr513 != nil {
					return parseErr513
				}
			}
		case "INCREMENT":
			p.advance() // consume INCREMENT
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			parseDiscard516, parseErr515 := p.parseExpr()
			_ = // consume integer
				parseDiscard516
			if parseErr515 != nil {
				return parseErr515
			}
		case "MAXVALUE":
			p.advance()
			parseDiscard518, parseErr517 := p.parseExpr()
			_ = parseDiscard518
			if parseErr517 != nil {
				return parseErr517
			}
		case "NOMAXVALUE":
			p.advance()
		case "MINVALUE":
			p.advance()
			parseDiscard520, parseErr519 := p.parseExpr()
			_ = parseDiscard520
			if parseErr519 != nil {
				return parseErr519
			}
		case "NOMINVALUE":
			p.advance()
		case "CYCLE":
			p.advance()
		case "NOCYCLE":
			p.advance()
		case "CACHE":
			p.advance()
			parseDiscard522, // consume CACHE
				parseErr521 := p.parseExpr()
			_ = // consume integer
				parseDiscard522
			if parseErr521 != nil {
				return parseErr521
			}
		case "NOCACHE":
			p.advance()
		case "ORDER":
			p.advance()
		case "NOORDER":
			p.advance()
		default:
			return nil
		}
	}
	return nil
}

// parseColumnConstraint parses a named column constraint: CONSTRAINT name <constraint_body>.
func (p *Parser) parseColumnConstraint() (*nodes.ColumnConstraint, error) {
	start := p.pos()
	p.advance() // consume CONSTRAINT

	name, parseErr523 := p.parseIdentifier()
	if parseErr523 != nil {
		return nil, parseErr523
	}

	cc, parseErr524 := p.parseColumnConstraintInline()
	if parseErr524 != nil {
		return nil, parseErr524
	}
	if cc == nil {
		return nil, nil
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
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		cc.Type = nodes.CONSTRAINT_PRIMARY

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		cc.Type = nodes.CONSTRAINT_UNIQUE

	case kwCHECK:
		p.advance() // consume CHECK
		if p.cur.Type == '(' {
			p.advance()
			var // consume '('
			parseErr525 error
			cc.Expr, parseErr525 = p.parseExpr()
			if parseErr525 != nil {
				return nil, parseErr525
			}
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
		}
		cc.Type = nodes.CONSTRAINT_CHECK

	case kwREFERENCES:
		p.advance() // consume REFERENCES
		cc.Type = nodes.CONSTRAINT_FOREIGN
		var parseErr526 error
		cc.RefTable, parseErr526 = p.parseObjectName()
		if parseErr526 != nil {
			return nil, parseErr526
		}
		if p.cur.Type == '(' {
			p.advance()
			var // consume '('
			parseErr527 error
			cc.RefColumns, parseErr527 = p.parseIdentifierList()
			if parseErr527 != nil {
				return nil, parseErr527
			}
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
		}
		// ON DELETE
		if p.cur.Type == kwON {
			next := p.peekNext()
			if next.Type == kwDELETE {
				p.advance() // consume ON
				p.advance()
				var // consume DELETE
				parseErr528 error
				cc.OnDelete, parseErr528 = p.parseDeleteAction()
				if parseErr528 != nil {
					return nil, parseErr528

					// DEFERRABLE / NOT DEFERRABLE
				}
			}
		}

	default:
		return nil, nil
	}

	if p.cur.Type == kwDEFERRABLE {
		cc.Deferrable = true
		p.advance()
	} else if p.cur.Type == kwNOT {
		next := p.peekNext()
		if next.Type == kwDEFERRABLE {
			p.advance() // consume NOT
			p.advance() // consume DEFERRABLE
			cc.Deferrable = false
		}
	}

	// INITIALLY DEFERRED / INITIALLY IMMEDIATE
	if p.cur.Type == kwINITIALLY {
		p.advance() // consume INITIALLY
		if p.cur.Type == kwDEFERRED {
			cc.Initially = "DEFERRED"
			p.advance()
		} else if p.cur.Type == kwIMMEDIATE {
			cc.Initially = "IMMEDIATE"
			p.advance()
		}
	}

	cc.Loc.End = p.prev.End
	return cc, nil
}

// parseTableConstraint parses a table-level constraint.
//
//	[ CONSTRAINT name ] { PRIMARY KEY (cols) | UNIQUE (cols) | CHECK (expr) | FOREIGN KEY (cols) REFERENCES ... }
func (p *Parser) parseTableConstraint() (*nodes.TableConstraint, error) {
	start := p.pos()
	tc := &nodes.TableConstraint{
		Loc: nodes.Loc{Start: start},
	}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		var // consume CONSTRAINT
		parseErr529 error
		tc.Name, parseErr529 = p.parseIdentifier()
		if parseErr529 != nil {
			return nil, parseErr529
		}
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		tc.Type = nodes.CONSTRAINT_PRIMARY
		if p.cur.Type == '(' {
			p.advance()
			var parseErr530 error
			tc.Columns, parseErr530 = p.parseIdentifierListAsStrings()
			if parseErr530 != nil {
				return nil, parseErr530
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		tc.Type = nodes.CONSTRAINT_UNIQUE
		if p.cur.Type == '(' {
			p.advance()
			var parseErr531 error
			tc.Columns, parseErr531 = p.parseIdentifierListAsStrings()
			if parseErr531 != nil {
				return nil, parseErr531
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case kwCHECK:
		p.advance() // consume CHECK
		tc.Type = nodes.CONSTRAINT_CHECK
		if p.cur.Type == '(' {
			p.advance()
			var parseErr532 error
			tc.Expr, parseErr532 = p.parseExpr()
			if parseErr532 != nil {
				return nil, parseErr532
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		tc.Type = nodes.CONSTRAINT_FOREIGN
		if p.cur.Type == '(' {
			p.advance()
			var parseErr533 error
			tc.Columns, parseErr533 = p.parseIdentifierListAsStrings()
			if parseErr533 != nil {
				return nil, parseErr533
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		if p.cur.Type == kwREFERENCES {
			p.advance()
			var // consume REFERENCES
			parseErr534 error
			tc.RefTable, parseErr534 = p.parseObjectName()
			if parseErr534 != nil {
				return nil, parseErr534
			}
			if p.cur.Type == '(' {
				p.advance()
				var parseErr535 error
				tc.RefColumns, parseErr535 = p.parseIdentifierListAsStrings()
				if parseErr535 != nil {
					return nil, parseErr535
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
		// ON DELETE
		if p.cur.Type == kwON {
			next := p.peekNext()
			if next.Type == kwDELETE {
				p.advance() // consume ON
				p.advance()
				var // consume DELETE
				parseErr536 error
				tc.OnDelete, parseErr536 = p.parseDeleteAction()
				if parseErr536 != nil {
					return nil, parseErr536

					// DEFERRABLE / NOT DEFERRABLE
				}
			}
		}

	default:
		return nil, nil
	}

	if p.cur.Type == kwDEFERRABLE {
		tc.Deferrable = true
		p.advance()
	} else if p.cur.Type == kwNOT {
		next := p.peekNext()
		if next.Type == kwDEFERRABLE {
			p.advance() // consume NOT
			p.advance() // consume DEFERRABLE
			tc.Deferrable = false
		}
	}

	// INITIALLY DEFERRED / INITIALLY IMMEDIATE
	if p.cur.Type == kwINITIALLY {
		p.advance() // consume INITIALLY
		if p.cur.Type == kwDEFERRED {
			tc.Initially = "DEFERRED"
			p.advance()
		} else if p.cur.Type == kwIMMEDIATE {
			tc.Initially = "IMMEDIATE"
			p.advance()
		}
	}

	tc.Loc.End = p.prev.End
	return tc, nil
}

// parseIdentifierList parses a comma-separated list of identifiers,
// returning a *List of *String nodes.
func (p *Parser) parseIdentifierList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		name, parseErr537 := p.parseIdentifier()
		if parseErr537 != nil {
			return nil, parseErr537
		}
		if name == "" {
			break
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list, nil
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
		if p.cur.Type == kwNULL {
			p.advance() // consume NULL
			return "SET NULL", nil
		}
		return "SET", nil
	default:
		if p.isIdentLikeStr("RESTRICT") || p.cur.Type == kwRESTRICT {
			p.advance()
			return "RESTRICT", nil
		}
		if p.isIdentLikeStr("NO") {
			p.advance()
			if p.isIdentLikeStr("ACTION") {
				p.advance()
				return "NO ACTION", nil
			}
		}
		return "", nil
	}
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
func (p *Parser) parseTableOptions(stmt *nodes.CreateTableStmt) error {
	for {
		switch p.cur.Type {
		case kwTABLESPACE:
			p.advance()
			var // consume TABLESPACE
			parseErr538 error
			stmt.Tablespace, parseErr538 = p.parseIdentifier()
			if parseErr538 != nil {

				// ON COMMIT { PRESERVE | DELETE } ROWS
				return parseErr538
			}

		case kwON:

			next := p.peekNext()
			if next.Type == kwCOMMIT {
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
			p.advance()
			stmt.Compress = "COMPRESS"

		case kwNOCOMPRESS:
			p.advance()
			stmt.Compress = "NOCOMPRESS"

		case kwPARTITION:
			var parseErr540 error
			stmt.Partition, parseErr540 = p.parsePartitionClause()
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
			if p.cur.Type == '(' {
				p.advance() // consume '('
				if p.isIdentLikeStr("MODE") {
					p.advance() // consume MODE
				}
				if p.cur.Type == kwDEFAULT {
					stmt.ResultCache = "DEFAULT"
					p.advance()
				} else if p.isIdentLikeStr("FORCE") {
					stmt.ResultCache = "FORCE"
					p.advance()
				}
				if p.cur.Type == ')' {
					p.advance() // consume ')'
				}
			}

		case kwFLASHBACK:
			// FLASHBACK ARCHIVE [ flashback_archive ]
			p.advance() // consume FLASHBACK
			if p.isIdentLikeStr("ARCHIVE") {
				p.advance() // consume ARCHIVE
				if p.isIdentLike() && !p.isStatementEnd() {
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
			// ROW STORE COMPRESS [ BASIC | ADVANCED ]
			next := p.peekNext()
			if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "STORE" {
				p.advance() // consume ROW
				p.advance() // consume STORE
				if p.cur.Type == kwCOMPRESS {
					p.advance() // consume COMPRESS
					stmt.Compress = "ROW STORE COMPRESS"
					if p.isIdentLikeStr("BASIC") || p.isIdentLikeStr("ADVANCED") {
						stmt.Compress += " " + p.cur.Str
						p.advance()
					}
				}
			} else {
				return nil
			}

		case kwREAD:
			// READ ONLY | READ WRITE
			next := p.peekNext()
			if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "ONLY" {
				p.advance() // consume READ
				p.advance() // consume ONLY
				stmt.ReadOnly = "READ ONLY"
			} else if next.Type == kwWRITE {
				p.advance() // consume READ
				p.advance() // consume WRITE
				stmt.ReadOnly = "READ WRITE"
			} else {
				return nil
			}

		case kwENABLE:
			// ENABLE ROW MOVEMENT or ENABLE VALIDATE/NOVALIDATE constraint
			next := p.peekNext()
			if next.Type == kwROW {
				p.advance() // consume ENABLE
				p.advance() // consume ROW
				if p.isIdentLikeStr("MOVEMENT") {
					p.advance() // consume MOVEMENT
				}
				stmt.RowMovement = "ENABLE"
			} else {
				parseErr542 :=
					// ENABLE [VALIDATE|NOVALIDATE] constraint_clause - skip it
					p.parseEnableDisableClause()
				if parseErr542 != nil {
					return parseErr542
				}
			}

		case kwDISABLE:
			next := p.peekNext()
			if next.Type == kwROW {
				p.advance() // consume DISABLE
				p.advance() // consume ROW
				if p.isIdentLikeStr("MOVEMENT") {
					p.advance() // consume MOVEMENT
				}
				stmt.RowMovement = "DISABLE"
			} else {
				parseErr543 := p.parseEnableDisableClause()
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
					if p.isIdentLikeStr("CREATION") {
						p.advance() // consume CREATION
						if p.cur.Type == kwIMMEDIATE {
							stmt.SegmentCreation = "IMMEDIATE"
							p.advance()
						} else if p.cur.Type == kwDEFERRED {
							stmt.SegmentCreation = "DEFERRED"
							p.advance()
						}
					}
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
						stmt.Organization = "INDEX"
						p.advance()
					} else if p.isIdentLikeStr("EXTERNAL") {
						optStart := p.pos()
						stmt.Organization = "EXTERNAL"
						p.advance()
						value := p.collectCreateTableOptionValue()
						appendCreateTableOption(stmt, "ORGANIZATION EXTERNAL", value, nodes.Loc{Start: optStart, End: p.prev.End})
					}
				case "INDEXING":
					// INDEXING { ON | OFF }
					p.advance() // consume INDEXING
					if p.cur.Type == kwON {
						stmt.Indexing = "ON"
						p.advance()
					} else if p.isIdentLikeStr("OFF") {
						stmt.Indexing = "OFF"
						p.advance()
					}
				case "ROWDEPENDENCIES":
					p.advance()
					stmt.RowDependencies = "ROWDEPENDENCIES"
				case "NOROWDEPENDENCIES":
					p.advance()
					stmt.RowDependencies = "NOROWDEPENDENCIES"
				case "PCTFREE":
					p.advance() // consume PCTFREE
					if p.cur.Type == tokICONST {
						p.advance()
					}
				case "PCTUSED":
					p.advance() // consume PCTUSED
					if p.cur.Type == tokICONST {
						p.advance()
					}
				case "INITRANS":
					p.advance() // consume INITRANS
					if p.cur.Type == tokICONST {
						p.advance()
					}
				case "COLUMN":
					// COLUMN STORE COMPRESS FOR { QUERY | ARCHIVE } [ LOW | HIGH ]
					next := p.peekNext()
					if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "STORE" {
						p.advance() // consume COLUMN
						p.advance() // consume STORE
						if p.cur.Type == kwCOMPRESS {
							p.advance() // consume COMPRESS
							stmt.Compress = "COLUMN STORE COMPRESS"
							if p.cur.Type == kwFOR {
								p.advance() // consume FOR
								if p.isIdentLikeStr("QUERY") || p.isIdentLikeStr("ARCHIVE") {
									stmt.Compress += " FOR " + p.cur.Str
									p.advance()
									if p.isIdentLikeStr("LOW") || p.isIdentLikeStr("HIGH") {
										stmt.Compress += " " + p.cur.Str
										p.advance()
									}
								}
							}
						}
					} else {
						return nil
					}
				case "INMEMORY":
					optStart := p.pos()
					p.advance()
					value := p.collectCreateTableOptionValue()
					appendCreateTableOption(stmt, "INMEMORY", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "LOB":
					optStart := p.pos()
					p.advance() // consume LOB
					if p.cur.Type != '(' {
						return p.syntaxErrorAtCur()
					}
					value := p.collectCreateTableStorageOptionValue()
					appendCreateTableOption(stmt, "LOB", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "VARRAY":
					optStart := p.pos()
					p.advance()
					value := p.collectCreateTableStorageOptionValue()
					appendCreateTableOption(stmt, "VARRAY", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "NESTED":
					optStart := p.pos()
					p.advance() // consume NESTED
					value := p.collectCreateTableStorageOptionValue()
					appendCreateTableOption(stmt, "NESTED", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "NO":
					// NO FLASHBACK ARCHIVE / NO INMEMORY
					next := p.peekNext()
					if next.Type == kwFLASHBACK {
						p.advance() // consume NO
						p.advance() // consume FLASHBACK
						if p.isIdentLikeStr("ARCHIVE") {
							p.advance() // consume ARCHIVE
						}
						stmt.FlashbackArchive = "NO FLASHBACK ARCHIVE"
					} else if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "INMEMORY" {
						p.advance() // consume NO
						p.advance() // consume INMEMORY
					} else {
						return nil
					}
				case "STORAGE":
					optStart := p.pos()
					p.advance()
					value := p.collectCreateTableOptionValue()
					appendCreateTableOption(stmt, "STORAGE", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "FILESYSTEM_LIKE_LOGGING":
					p.advance()
					stmt.Logging = "FILESYSTEM_LIKE_LOGGING"
				case "MAPPING":
					// MAPPING TABLE
					p.advance()
					if p.cur.Type == kwTABLE {
						p.advance()
					}
				case "NOMAPPING":
					p.advance()
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
					p.advance()
					value := p.collectCreateTableOptionValue()
					appendCreateTableOption(stmt, "ILM", value, nodes.Loc{Start: optStart, End: p.prev.End})
				case "CLUSTERING":
					optStart := p.pos()
					p.advance()
					value := p.collectCreateTableOptionValue()
					appendCreateTableOption(stmt, "CLUSTERING", value, nodes.Loc{Start: optStart, End: p.prev.End})
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
	return nil
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

func (p *Parser) collectCreateTableStorageOptionValue() string {
	tokens := p.collectDDLTokensUntil(p.isCreateTableStorageOptionBoundary)
	return strings.Join(tokens, " ")
}

func (p *Parser) isCreateTableStorageOptionBoundary() bool {
	if p.cur.Type == kwAS {
		return false
	}
	return p.isCreateTableOptionBoundary()
}

func (p *Parser) isCreateTableOptionBoundary() bool {
	switch p.cur.Type {
	case ';', tokEOF, kwAS, kwCACHE, kwDISABLE, kwENABLE, kwFLASHBACK, kwLOGGING,
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
			"VARRAY":
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

// parseImmutableTableClauses parses immutable table clauses.
//
// BNF (CREATE-TABLE.bnf lines 733-741):
//
//	immutable_table_no_drop_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	immutable_table_no_delete_clause: NO DELETE [ LOCKED | UNTIL integer DAYS AFTER INSERT ]
func (p *Parser) parseImmutableTableClauses(stmt *nodes.CreateTableStmt) error {
	if !p.isIdentLikeStr("NO") {
		return nil
	}
	next := p.peekNext()
	// NO DROP
	if next.Type == kwDROP {
		p.advance() // consume NO
		p.advance() // consume DROP
		noDrop := "NO DROP"
		if p.isIdentLikeStr("UNTIL") {
			p.advance() // consume UNTIL
			days := ""
			if p.cur.Type == tokICONST {
				days = p.cur.Str
				p.advance()
			}
			if p.isIdentLikeStr("DAYS") {
				p.advance() // consume DAYS
			}
			if p.isIdentLikeStr("IDLE") {
				p.advance() // consume IDLE
			}
			noDrop = "NO DROP UNTIL " + days + " DAYS IDLE"
		}
		stmt.ImmutableNoDrop = noDrop
	}

	// NO DELETE
	if p.isIdentLikeStr("NO") {
		next2 := p.peekNext()
		if next2.Type == kwDELETE {
			p.advance() // consume NO
			p.advance() // consume DELETE
			noDel := "NO DELETE"
			if p.isIdentLikeStr("LOCKED") {
				p.advance()
				noDel = "NO DELETE LOCKED"
			} else if p.isIdentLikeStr("UNTIL") {
				p.advance() // consume UNTIL
				days := ""
				if p.cur.Type == tokICONST {
					days = p.cur.Str
					p.advance()
				}
				if p.isIdentLikeStr("DAYS") {
					p.advance() // consume DAYS
				}
				if p.cur.Type == kwAFTER {
					p.advance() // consume AFTER
				}
				if p.cur.Type == kwINSERT {
					p.advance() // consume INSERT
				}
				noDel = "NO DELETE UNTIL " + days + " DAYS AFTER INSERT"
			}
			stmt.ImmutableNoDel = noDel
		}
	}
	return nil
}

// parseBlockchainTableClauses parses blockchain table clauses.
//
// BNF (CREATE-TABLE.bnf lines 743-756):
//
//	blockchain_drop_table_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	blockchain_row_retention_clause: NO DELETE { LOCKED | UNTIL integer DAYS AFTER INSERT }
//	blockchain_hash_and_data_format_clause: HASHING USING 'hash_algorithm' VERSION 'version_string'
func (p *Parser) parseBlockchainTableClauses(stmt *nodes.CreateTableStmt) error {
	// HASHING USING 'hash_algorithm'
	if p.isIdentLikeStr("HASHING") {
		p.advance() // consume HASHING
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
		}
		if p.cur.Type == tokSCONST {
			stmt.BlockchainHash = p.cur.Str
			p.advance()
		}
	}
	// VERSION 'version_string'
	if p.isIdentLikeStr("VERSION") {
		p.advance() // consume VERSION
		if p.cur.Type == tokSCONST {
			stmt.BlockchainVer = p.cur.Str
			p.advance()
		}
	}
	return nil
}

// parseEnableDisableClause parses ENABLE/DISABLE [VALIDATE|NOVALIDATE] constraint clause.
//
// BNF (CREATE-TABLE.bnf lines 463-472):
//
//	{ ENABLE | DISABLE } [ VALIDATE | NOVALIDATE ]
//	{ UNIQUE ( column [,...] ) | PRIMARY KEY | CONSTRAINT constraint_name }
//	[ using_index_clause ] [ exceptions_clause ] [ CASCADE ] [ { KEEP | DROP } INDEX ]
func (p *Parser) parseEnableDisableClause() error {
	p.advance() // consume ENABLE or DISABLE

	// [ VALIDATE | NOVALIDATE ]
	if p.isIdentLikeStr("VALIDATE") || p.isIdentLikeStr("NOVALIDATE") {
		p.advance()
	}

	// { UNIQUE (cols) | PRIMARY KEY | CONSTRAINT name }
	switch p.cur.Type {
	case kwUNIQUE:
		p.advance()
		if p.cur.Type == '(' {
			p.skipParenthesized()
		}
	case kwPRIMARY:
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
		}
	case kwCONSTRAINT:
		p.advance()
		parseDiscard554, parseErr553 := p.parseIdentifier()
		_ = parseDiscard554
		if parseErr553 !=

			// [ USING INDEX ... ]
			nil {
			return parseErr553
		}
	default:
		return nil
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

// isStatementEnd returns true if the current token is at a statement boundary.
func (p *Parser) isStatementEnd() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF
}

// parsePartitionClause parses a PARTITION BY clause.
//
//	PARTITION BY { RANGE | LIST | HASH } (columns)
//	    ( partition_def [,...] )
func (p *Parser) parsePartitionClause() (*nodes.PartitionClause, error) {
	start := p.pos()
	p.advance() // consume PARTITION

	clause := &nodes.PartitionClause{
		Columns:    &nodes.List{},
		Partitions: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// BY
	if p.cur.Type == kwBY {
		p.advance()
	}

	// RANGE / LIST / HASH
	switch {
	case p.isIdentLike() && p.cur.Str == "RANGE":
		clause.Type = nodes.PARTITION_RANGE
		p.advance()
	case p.isIdentLike() && p.cur.Str == "LIST":
		clause.Type = nodes.PARTITION_LIST
		p.advance()
	case p.isIdentLike() && p.cur.Str == "HASH":
		clause.Type = nodes.PARTITION_HASH
		p.advance()
	}

	// (columns)
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col, parseErr559 := p.parseExpr()
			if parseErr559 != nil {
				return nil, parseErr559
			}
			if col != nil {
				clause.Columns.Items = append(clause.Columns.Items, col)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Optional INTERVAL (expr) for range-interval partitioned tables.
	if p.cur.Type == kwINTERVAL {
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		interval, parseErr560 := p.parseExpr()
		if parseErr560 != nil {
			return nil, parseErr560
		}
		if interval == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		clause.Interval = interval
	}

	// Optional SUBPARTITION BY
	if p.isIdentLike() && p.cur.Str == "SUBPARTITION" {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		var parseErr560 error
		clause.Subpartition, parseErr560 = p.parsePartitionClause()
		if parseErr560 !=

			// Partition definitions: ( partition p1 ... [,...] )
			nil {
			return nil, parseErr560
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			pDef, parseErr561 := p.parsePartitionDef()
			if parseErr561 != nil {
				return nil, parseErr561
			}
			if pDef != nil {
				clause.Partitions.Items = append(clause.Partitions.Items, pDef)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	if clause.Interval != nil && (clause.Partitions == nil || clause.Partitions.Len() == 0) {
		return nil, p.syntaxErrorAtCur()
	}

	clause.Loc.End = p.prev.End
	return clause, nil
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION name VALUES LESS THAN (expr) [TABLESPACE ts]
//	PARTITION name VALUES (expr [,...]) [TABLESPACE ts]
func (p *Parser) parsePartitionDef() (*nodes.PartitionDef, error) {
	start := p.pos()

	if p.cur.Type != kwPARTITION {
		// Skip unknown tokens until , or )
		for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
			p.advance()
		}
		return nil, nil
	}
	p.advance() // consume PARTITION

	def := &nodes.PartitionDef{
		Values: &nodes.List{},
		Loc:    nodes.Loc{Start: start},
	}

	// Partition name
	if p.isIdentLike() {
		def.Name = p.cur.Str
		p.advance()
	}

	// VALUES LESS THAN (expr) or VALUES (expr,...)
	if p.isIdentLike() && p.cur.Str == "VALUES" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "LESS" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "THAN" {
				p.advance()
			}
		}
		// (expr [,...])
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				val, parseErr562 := p.parseExpr()
				if parseErr562 != nil {
					return nil, parseErr562
				}
				if val != nil {
					def.Values.Items = append(def.Values.Items, val)
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// TABLESPACE
	if p.cur.Type == kwTABLESPACE {
		p.advance()
		var parseErr563 error
		def.Tablespace, parseErr563 = p.parseIdentifier()
		if parseErr563 !=

			// Skip any remaining options (LOGGING, etc.) until , or )
			nil {
			return nil, parseErr563
		}
	}

	for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
		p.advance()
	}

	def.Loc.End = p.prev.End
	return def, nil
}

// parseOnCommitAction parses the action after ON COMMIT.
//
//	PRESERVE ROWS | DELETE ROWS
func (p *Parser) parseOnCommitAction() (string, error) {
	switch {
	case p.isIdentLikeStr("PRESERVE"):
		p.advance() // consume PRESERVE
		if p.cur.Type == kwROWS {
			p.advance() // consume ROWS
		}
		return "PRESERVE ROWS", nil
	case p.cur.Type == kwDELETE:
		p.advance() // consume DELETE
		if p.cur.Type == kwROWS {
			p.advance() // consume ROWS
		}
		return "DELETE ROWS", nil
	default:
		return "", nil
	}
}
