package parser

import (
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
func (p *Parser) parseCreateStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume CREATE

	// OR REPLACE
	orReplace := false
	if p.cur.Type == kwOR {
		p.advance() // consume OR
		if p.cur.Type == kwREPLACE {
			p.advance() // consume REPLACE
			orReplace = true
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
		return p.parseCreateTypeStmt(start, orReplace)
	case kwPROCEDURE:
		return p.parseCreateProcedureStmt(start, orReplace)
	case kwFUNCTION:
		return p.parseCreateFunctionStmt(start, orReplace)
	case kwPACKAGE:
		return p.parseCreatePackageStmt(start, orReplace)
	case kwTRIGGER:
		return p.parseCreateTriggerStmt(start, orReplace)
	case kwAUDIT:
		// CREATE AUDIT POLICY
		p.advance() // consume AUDIT
		if p.isIdentLikeStr("POLICY") {
			p.advance() // consume POLICY
		}
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_AUDIT_POLICY, start)
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
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_JSON_DUALITY_VIEW, start)
	case kwUSER, kwROLE, kwPROFILE,
		kwTABLESPACE, kwDIRECTORY, kwCONTEXT,
		kwCLUSTER, kwJAVA, kwLIBRARY, kwSCHEMA:
		return p.parseCreateAdminObject(start)
	case kwTEMPORARY:
		// CREATE TEMPORARY TABLESPACE (standalone, not GLOBAL TEMPORARY TABLE)
		if !global && !private {
			p.advance() // consume TEMPORARY
			if p.cur.Type == kwTABLESPACE {
				p.advance() // consume TABLESPACE
				return p.parseCreateTablespaceStmt(start, false, false, true, false)
			}
		}
		return nil
	default:
		// Check for "NO FORCE VIEW"
		if p.isIdentLikeStr("NO") || p.cur.Type == kwNOT {
			return p.parseCreateViewStmt(start, orReplace)
		}
		// Check for BIGFILE/SMALLFILE/UNDO TABLESPACE, CONTROLFILE
		if p.isIdentLike() {
			switch p.cur.Str {
			case "BIGFILE":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, true, false, false, false)
				}
			case "SMALLFILE":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, false, true, false, false)
				}
			case "UNDO":
				p.advance()
				if p.cur.Type == kwTABLESPACE {
					p.advance()
					return p.parseCreateTablespaceStmt(start, false, false, false, true)
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
		// Check for INDEXTYPE and OPERATOR (identifier-based keywords)
		if p.isIdentLike() {
			switch p.cur.Str {
			case "INDEXTYPE":
				p.advance() // consume INDEXTYPE
				return p.parseCreateIndextypeStmt(start, orReplace)
			case "OPERATOR":
				p.advance() // consume OPERATOR
				return p.parseCreateOperatorStmt(start, orReplace, false)
			}
		}
		// Check for DIMENSION, FLASHBACK ARCHIVE
		if adminStmt := p.parseCreateAdminObject(start); adminStmt != nil {
			return adminStmt
		}
		return nil
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
func (p *Parser) parseCreateTableStmt(start int, orReplace, global, private, sharded, duplicated bool) *nodes.CreateTableStmt {
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		OrReplace:  orReplace,
		Global:     global,
		Private:    private,
		Sharded:    sharded,
		Duplicated: duplicated,
		Columns:    &nodes.List{},
		Constraints: &nodes.List{},
		Loc:        nodes.Loc{Start: start},
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
	stmt.Name = p.parseObjectName()

	// SHARING = { METADATA | DATA | EXTENDED DATA | NONE }
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
		p.advance() // consume AS
		stmt.AsQuery = p.parseSelectStmt()
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Column definitions and table constraints
	if p.cur.Type == '(' {
		p.advance() // consume '('
		p.parseColumnDefsAndConstraints(stmt)
		if p.cur.Type == ')' {
			p.advance() // consume ')'
		}
	}

	// Immutable table clauses (before table options)
	// BNF (lines 733-741): immutable_table_clauses
	p.parseImmutableTableClauses(stmt)

	// Blockchain table clauses
	// BNF (lines 743-756): blockchain_table_clauses
	p.parseBlockchainTableClauses(stmt)

	// DEFAULT COLLATION collation_name
	if p.cur.Type == kwDEFAULT {
		next := p.peekNext()
		if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "COLLATION" {
			p.advance() // consume DEFAULT
			p.advance() // consume COLLATION
			stmt.Collation = p.parseIdentifier()
		}
	}

	// Table options (physical_properties + table_properties)
	p.parseTableOptions(stmt)

	// MEMOPTIMIZE FOR READ
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
		p.advance() // consume PARENT
		stmt.Parent = p.parseObjectName()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseColumnDefsAndConstraints parses the contents inside the parentheses of a CREATE TABLE.
// It handles both column definitions and table-level constraints.
func (p *Parser) parseColumnDefsAndConstraints(stmt *nodes.CreateTableStmt) {
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Check if this is a table-level constraint
		if p.isTableConstraintStart() {
			tc := p.parseTableConstraint()
			if tc != nil {
				stmt.Constraints.Items = append(stmt.Constraints.Items, tc)
			}
		} else {
			// Column definition
			col := p.parseColumnDef()
			if col != nil {
				stmt.Columns.Items = append(stmt.Columns.Items, col)
			}
		}

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
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
func (p *Parser) parseColumnDef() *nodes.ColumnDef {
	start := p.pos()
	col := &nodes.ColumnDef{
		Constraints: &nodes.List{},
		Loc:         nodes.Loc{Start: start},
	}

	col.Name = p.parseIdentifier()
	if col.Name == "" {
		return nil
	}

	// DOMAIN [ schema. ] domain_name  (instead of datatype)
	if p.isIdentLikeStr("DOMAIN") {
		p.advance() // consume DOMAIN
		col.Domain = p.parseObjectName()
	} else if p.isTypeName() {
		// Data type (optional for some Oracle column types, but typical)
		col.TypeName = p.parseTypeName()
	}

	// Column properties: SORT, VISIBLE, DEFAULT, NOT NULL, NULL, constraints, etc.
	p.parseColumnProperties(col)

	col.Loc.End = p.pos()
	return col
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
func (p *Parser) parseColumnProperties(col *nodes.ColumnDef) {
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
			col.Default = p.parseExpr()

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
						p.parseIdentityOptions(identity)
						if p.cur.Type == ')' {
							p.advance()
						}
					}
					identity.Loc.End = p.pos()
					col.Identity = identity
				} else if p.cur.Type == '(' {
					// virtual column: AS (expr) [VIRTUAL]
					p.advance() // consume '('
					col.Virtual = p.parseExpr()
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
				return
			}

		case kwNULL:
			p.advance() // consume NULL
			col.Null = true

		case kwINVISIBLE:
			p.advance()
			col.Invisible = true

		case kwCONSTRAINT:
			cc := p.parseColumnConstraint()
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwPRIMARY:
			cc := p.parseColumnConstraintInline()
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwUNIQUE:
			cc := p.parseColumnConstraintInline()
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwCHECK:
			cc := p.parseColumnConstraintInline()
			if cc != nil {
				col.Constraints.Items = append(col.Constraints.Items, cc)
			}

		case kwREFERENCES:
			cc := p.parseColumnConstraintInline()
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
				return
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
					// encryption_spec: USING 'algo' IDENTIFIED BY pw SALT|NO SALT
					p.parseEncryptionSpec(col)
				case "DECRYPT":
					p.advance() // consume DECRYPT
					col.Encrypt = "DECRYPT"
				case "COLLATE":
					p.advance() // consume COLLATE
					col.Collation = p.parseIdentifier()
				default:
					return
				}
			} else {
				return
			}
		}
	}
}

// parseEncryptionSpec parses the encryption_spec after ENCRYPT keyword.
//
// BNF (CREATE-TABLE.bnf lines 95-98):
//
//	[ USING 'encrypt_algorithm' ]
//	[ IDENTIFIED BY password ]
//	[ SALT | NO SALT ]
func (p *Parser) parseEncryptionSpec(col *nodes.ColumnDef) {
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
				p.advance() // consume BY
				p.parseIdentifier() // consume password
			}
		} else if p.isIdentLikeStr("SALT") {
			p.advance() // consume SALT
		} else if p.isIdentLikeStr("NO") {
			next := p.peekNext()
			if (next.Type == tokIDENT || next.Type >= 2000) && next.Str == "SALT" {
				p.advance() // consume NO
				p.advance() // consume SALT
			} else {
				return
			}
		} else {
			return
		}
	}
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
func (p *Parser) parseIdentityOptions(identity *nodes.IdentityClause) {
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if !p.isIdentLike() {
			return
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
				p.parseExpr() // consume integer
			}
		case "INCREMENT":
			p.advance() // consume INCREMENT
			if p.cur.Type == kwBY {
				p.advance() // consume BY
			}
			p.parseExpr() // consume integer
		case "MAXVALUE":
			p.advance()
			p.parseExpr()
		case "NOMAXVALUE":
			p.advance()
		case "MINVALUE":
			p.advance()
			p.parseExpr()
		case "NOMINVALUE":
			p.advance()
		case "CYCLE":
			p.advance()
		case "NOCYCLE":
			p.advance()
		case "CACHE":
			p.advance() // consume CACHE
			p.parseExpr() // consume integer
		case "NOCACHE":
			p.advance()
		case "ORDER":
			p.advance()
		case "NOORDER":
			p.advance()
		default:
			return
		}
	}
}

// parseColumnConstraint parses a named column constraint: CONSTRAINT name <constraint_body>.
func (p *Parser) parseColumnConstraint() *nodes.ColumnConstraint {
	start := p.pos()
	p.advance() // consume CONSTRAINT

	name := p.parseIdentifier()

	cc := p.parseColumnConstraintInline()
	if cc == nil {
		return nil
	}
	cc.Name = name
	cc.Loc.Start = start
	return cc
}

// parseColumnConstraintInline parses an inline (unnamed) column constraint body.
func (p *Parser) parseColumnConstraintInline() *nodes.ColumnConstraint {
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
			p.advance() // consume '('
			cc.Expr = p.parseExpr()
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
		}
		cc.Type = nodes.CONSTRAINT_CHECK

	case kwREFERENCES:
		p.advance() // consume REFERENCES
		cc.Type = nodes.CONSTRAINT_FOREIGN
		cc.RefTable = p.parseObjectName()
		if p.cur.Type == '(' {
			p.advance() // consume '('
			cc.RefColumns = p.parseIdentifierList()
			if p.cur.Type == ')' {
				p.advance() // consume ')'
			}
		}
		// ON DELETE
		if p.cur.Type == kwON {
			next := p.peekNext()
			if next.Type == kwDELETE {
				p.advance() // consume ON
				p.advance() // consume DELETE
				cc.OnDelete = p.parseDeleteAction()
			}
		}

	default:
		return nil
	}

	// DEFERRABLE / NOT DEFERRABLE
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

	cc.Loc.End = p.pos()
	return cc
}

// parseTableConstraint parses a table-level constraint.
//
//	[ CONSTRAINT name ] { PRIMARY KEY (cols) | UNIQUE (cols) | CHECK (expr) | FOREIGN KEY (cols) REFERENCES ... }
func (p *Parser) parseTableConstraint() *nodes.TableConstraint {
	start := p.pos()
	tc := &nodes.TableConstraint{
		Loc: nodes.Loc{Start: start},
	}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		tc.Name = p.parseIdentifier()
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
			tc.Columns = p.parseIdentifierListAsStrings()
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		tc.Type = nodes.CONSTRAINT_UNIQUE
		if p.cur.Type == '(' {
			p.advance()
			tc.Columns = p.parseIdentifierListAsStrings()
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case kwCHECK:
		p.advance() // consume CHECK
		tc.Type = nodes.CONSTRAINT_CHECK
		if p.cur.Type == '(' {
			p.advance()
			tc.Expr = p.parseExpr()
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
			tc.Columns = p.parseIdentifierListAsStrings()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		if p.cur.Type == kwREFERENCES {
			p.advance() // consume REFERENCES
			tc.RefTable = p.parseObjectName()
			if p.cur.Type == '(' {
				p.advance()
				tc.RefColumns = p.parseIdentifierListAsStrings()
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
				p.advance() // consume DELETE
				tc.OnDelete = p.parseDeleteAction()
			}
		}

	default:
		return nil
	}

	// DEFERRABLE / NOT DEFERRABLE
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

	tc.Loc.End = p.pos()
	return tc
}

// parseIdentifierList parses a comma-separated list of identifiers,
// returning a *List of *String nodes.
func (p *Parser) parseIdentifierList() *nodes.List {
	list := &nodes.List{}
	for {
		name := p.parseIdentifier()
		if name == "" {
			break
		}
		list.Items = append(list.Items, &nodes.String{Str: name})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list
}

// parseIdentifierListAsStrings parses a comma-separated list of identifiers,
// returning a *List of *String nodes. Same as parseIdentifierList.
func (p *Parser) parseIdentifierListAsStrings() *nodes.List {
	return p.parseIdentifierList()
}

// parseDeleteAction parses the ON DELETE action (CASCADE, SET NULL, etc.).
func (p *Parser) parseDeleteAction() string {
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		return "CASCADE"
	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwNULL {
			p.advance() // consume NULL
			return "SET NULL"
		}
		return "SET"
	default:
		if p.isIdentLikeStr("RESTRICT") || p.cur.Type == kwRESTRICT {
			p.advance()
			return "RESTRICT"
		}
		if p.isIdentLikeStr("NO") {
			p.advance()
			if p.isIdentLikeStr("ACTION") {
				p.advance()
				return "NO ACTION"
			}
		}
		return ""
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
func (p *Parser) parseTableOptions(stmt *nodes.CreateTableStmt) {
	for {
		switch p.cur.Type {
		case kwTABLESPACE:
			p.advance() // consume TABLESPACE
			stmt.Tablespace = p.parseIdentifier()

		case kwON:
			// ON COMMIT { PRESERVE | DELETE } ROWS
			next := p.peekNext()
			if next.Type == kwCOMMIT {
				p.advance() // consume ON
				p.advance() // consume COMMIT
				stmt.OnCommit = p.parseOnCommitAction()
			} else {
				return
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
			stmt.Partition = p.parsePartitionClause()

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
					stmt.FlashbackArchive = p.parseIdentifier()
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
				return
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
				return
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
				// ENABLE [VALIDATE|NOVALIDATE] constraint_clause - skip it
				p.parseEnableDisableClause()
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
				p.parseEnableDisableClause()
			}

		case kwAS:
			// AS subquery (for CTAS after options)
			p.advance() // consume AS
			stmt.AsQuery = p.parseSelectStmt()
			return

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
					// ORGANIZATION { HEAP | INDEX | EXTERNAL }
					p.advance() // consume ORGANIZATION
					if p.isIdentLikeStr("HEAP") {
						stmt.Organization = "HEAP"
						p.advance()
					} else if p.cur.Type == kwINDEX {
						stmt.Organization = "INDEX"
						p.advance()
					} else if p.isIdentLikeStr("EXTERNAL") {
						stmt.Organization = "EXTERNAL"
						p.advance()
						// skip external table clause details
						if p.cur.Type == '(' {
							p.skipParenthesized()
						}
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
						return
					}
				case "INMEMORY":
					// INMEMORY [attributes] - skip details
					p.advance()
				case "LOB":
					// LOB storage clause - skip parenthesized content
					p.advance() // consume LOB
					if p.cur.Type == '(' {
						p.skipParenthesized()
					}
					// STORE AS ...
					if p.isIdentLikeStr("STORE") {
						p.advance()
						if p.cur.Type == kwAS {
							p.advance()
						}
						// skip rest of LOB storage
						if p.cur.Type == '(' {
							p.skipParenthesized()
						} else if p.isIdentLike() {
							p.advance() // LOB segment name
							if p.cur.Type == '(' {
								p.skipParenthesized()
							}
						}
					}
				case "VARRAY":
					// VARRAY col STORE AS ...
					p.advance() // consume VARRAY
					p.parseIdentifier() // consume varray item
					if p.isIdentLikeStr("STORE") {
						p.advance()
						if p.cur.Type == kwAS {
							p.advance()
						}
						// skip rest
						for p.isIdentLike() && p.cur.Type != tokEOF {
							p.advance()
						}
						if p.cur.Type == '(' {
							p.skipParenthesized()
						}
					}
				case "NESTED":
					// NESTED TABLE ... STORE AS ...
					p.advance() // consume NESTED
					if p.cur.Type == kwTABLE {
						p.advance() // consume TABLE
					}
					p.parseIdentifier() // nested item
					// skip to STORE AS
					for p.cur.Type != ';' && p.cur.Type != tokEOF {
						if p.isIdentLikeStr("STORE") {
							p.advance()
							if p.cur.Type == kwAS {
								p.advance()
							}
							p.parseIdentifier() // storage table name
							if p.cur.Type == '(' {
								p.skipParenthesized()
							}
							break
						}
						p.advance()
					}
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
						return
					}
				case "STORAGE":
					// STORAGE ( ... ) - skip
					p.advance()
					if p.cur.Type == '(' {
						p.skipParenthesized()
					}
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
					p.parseIdentifier()
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
					// ILM clause - skip details
					p.advance()
					for p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwPARTITION && p.cur.Type != kwENABLE && p.cur.Type != kwDISABLE {
						p.advance()
					}
				case "CLUSTERING":
					// attribute clustering - skip
					p.advance()
					for p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
						if p.cur.Type == '(' {
							p.skipParenthesized()
						} else {
							p.advance()
						}
						// stop at known table option keywords
						if p.cur.Type == kwCACHE || p.cur.Type == kwPARALLEL || p.cur.Type == kwNOPARALLEL || p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
							break
						}
					}
				case "SUPPLEMENTAL":
					// SUPPLEMENTAL LOG ... - skip
					p.advance() // consume SUPPLEMENTAL
					if p.isIdentLikeStr("LOG") {
						p.advance() // consume LOG
					}
					// skip until a known boundary
					for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != ';' && p.cur.Type != tokEOF {
						if p.cur.Type == '(' {
							p.skipParenthesized()
						} else {
							p.advance()
						}
					}
				default:
					return
				}
			} else {
				return
			}
		}
	}
}

// parseImmutableTableClauses parses immutable table clauses.
//
// BNF (CREATE-TABLE.bnf lines 733-741):
//
//	immutable_table_no_drop_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	immutable_table_no_delete_clause: NO DELETE [ LOCKED | UNTIL integer DAYS AFTER INSERT ]
func (p *Parser) parseImmutableTableClauses(stmt *nodes.CreateTableStmt) {
	if !p.isIdentLikeStr("NO") {
		return
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
}

// parseBlockchainTableClauses parses blockchain table clauses.
//
// BNF (CREATE-TABLE.bnf lines 743-756):
//
//	blockchain_drop_table_clause: NO DROP [ UNTIL integer DAYS IDLE ]
//	blockchain_row_retention_clause: NO DELETE { LOCKED | UNTIL integer DAYS AFTER INSERT }
//	blockchain_hash_and_data_format_clause: HASHING USING 'hash_algorithm' VERSION 'version_string'
func (p *Parser) parseBlockchainTableClauses(stmt *nodes.CreateTableStmt) {
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
}

// parseEnableDisableClause parses ENABLE/DISABLE [VALIDATE|NOVALIDATE] constraint clause.
//
// BNF (CREATE-TABLE.bnf lines 463-472):
//
//	{ ENABLE | DISABLE } [ VALIDATE | NOVALIDATE ]
//	{ UNIQUE ( column [,...] ) | PRIMARY KEY | CONSTRAINT constraint_name }
//	[ using_index_clause ] [ exceptions_clause ] [ CASCADE ] [ { KEEP | DROP } INDEX ]
func (p *Parser) parseEnableDisableClause() {
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
		p.parseIdentifier()
	default:
		return
	}

	// [ USING INDEX ... ]
	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type == kwINDEX {
			p.advance()
		}
		// skip index properties
		if p.cur.Type == '(' {
			p.skipParenthesized()
		} else if p.isIdentLike() {
			p.parseIdentifier()
		}
	}

	// [ EXCEPTIONS INTO table ]
	if p.isIdentLikeStr("EXCEPTIONS") {
		p.advance()
		if p.cur.Type == kwINTO {
			p.advance()
			p.parseObjectName()
		}
	}

	// [ CASCADE ]
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
}

// isStatementEnd returns true if the current token is at a statement boundary.
func (p *Parser) isStatementEnd() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF
}

// parsePartitionClause parses a PARTITION BY clause.
//
//	PARTITION BY { RANGE | LIST | HASH } (columns)
//	    ( partition_def [,...] )
func (p *Parser) parsePartitionClause() *nodes.PartitionClause {
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
			col := p.parseExpr()
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

	// Optional SUBPARTITION BY
	if p.isIdentLike() && p.cur.Str == "SUBPARTITION" {
		p.advance()
		if p.cur.Type == kwBY {
			p.advance()
		}
		clause.Subpartition = p.parsePartitionClause()
	}

	// Partition definitions: ( partition p1 ... [,...] )
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			pDef := p.parsePartitionDef()
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

	clause.Loc.End = p.pos()
	return clause
}

// parsePartitionDef parses a single partition definition.
//
//	PARTITION name VALUES LESS THAN (expr) [TABLESPACE ts]
//	PARTITION name VALUES (expr [,...]) [TABLESPACE ts]
func (p *Parser) parsePartitionDef() *nodes.PartitionDef {
	start := p.pos()

	if p.cur.Type != kwPARTITION {
		// Skip unknown tokens until , or )
		for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
			p.advance()
		}
		return nil
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
				val := p.parseExpr()
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
		def.Tablespace = p.parseIdentifier()
	}

	// Skip any remaining options (LOGGING, etc.) until , or )
	for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
		p.advance()
	}

	def.Loc.End = p.pos()
	return def
}

// parseOnCommitAction parses the action after ON COMMIT.
//
//	PRESERVE ROWS | DELETE ROWS
func (p *Parser) parseOnCommitAction() string {
	switch {
	case p.isIdentLikeStr("PRESERVE"):
		p.advance() // consume PRESERVE
		if p.cur.Type == kwROWS {
			p.advance() // consume ROWS
		}
		return "PRESERVE ROWS"
	case p.cur.Type == kwDELETE:
		p.advance() // consume DELETE
		if p.cur.Type == kwROWS {
			p.advance() // consume ROWS
		}
		return "DELETE ROWS"
	default:
		return ""
	}
}
