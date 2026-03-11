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

	// PUBLIC (for CREATE PUBLIC SYNONYM / DATABASE LINK)
	public := false
	if p.cur.Type == kwPUBLIC {
		public = true
		p.advance()
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, global, private)
	case kwUNIQUE, kwBITMAP, kwINDEX:
		return p.parseCreateIndexStmt(start)
	case kwMATERIALIZED, kwFORCE, kwVIEW:
		return p.parseCreateViewStmt(start, orReplace)
	case kwSEQUENCE:
		return p.parseCreateSequenceStmt(start)
	case kwSYNONYM:
		return p.parseCreateSynonymStmt(start, orReplace, public)
	case kwDATABASE:
		return p.parseCreateDatabaseLinkStmt(start, public)
	case kwTYPE:
		return p.parseCreateTypeStmt(start, orReplace)
	default:
		// Check for "NO FORCE VIEW"
		if p.isIdentLikeStr("NO") || p.cur.Type == kwNOT {
			return p.parseCreateViewStmt(start, orReplace)
		}
		return nil
	}
}

// parseCreateTableStmt parses a CREATE TABLE statement after the TABLE keyword.
// The caller has already consumed CREATE [OR REPLACE] [GLOBAL TEMPORARY] [PRIVATE TEMPORARY].
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLE.html
//
//	TABLE [ IF NOT EXISTS ] [ schema . ] table_name
//	    ( column_def | table_constraint [, ...] )
//	    [ ON COMMIT { PRESERVE | DELETE } ROWS ]
//	    [ TABLESPACE tablespace_name ]
//	    [ ... ]
//	  | TABLE [ schema . ] table_name AS subquery
func (p *Parser) parseCreateTableStmt(start int, orReplace, global, private bool) *nodes.CreateTableStmt {
	p.advance() // consume TABLE

	stmt := &nodes.CreateTableStmt{
		OrReplace:   orReplace,
		Global:      global,
		Private:     private,
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
	stmt.Name = p.parseObjectName()

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

	// Table options
	p.parseTableOptions(stmt)

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
//	column_name datatype [ DEFAULT expr ] [ column_constraint ... ]
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

	// Data type (optional for some Oracle column types, but typical)
	if p.isTypeName() {
		col.TypeName = p.parseTypeName()
	}

	// Column properties: DEFAULT, NOT NULL, NULL, constraints, etc.
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
// DEFAULT, NOT NULL, NULL, INVISIBLE, constraints, etc.
func (p *Parser) parseColumnProperties(col *nodes.ColumnDef) {
	for {
		switch p.cur.Type {
		case kwDEFAULT:
			p.advance() // consume DEFAULT
			col.Default = p.parseExpr()

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
// These include TABLESPACE, ON COMMIT, PARALLEL, COMPRESS, etc.
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

		case kwNOPARALLEL:
			p.advance()
			stmt.Parallel = "NOPARALLEL"

		case kwCOMPRESS:
			p.advance()
			stmt.Compress = "COMPRESS"

		case kwNOCOMPRESS:
			p.advance()
			stmt.Compress = "NOCOMPRESS"

		default:
			return
		}
	}
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
