package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateMaterializedOrView distinguishes between:
// - CREATE MATERIALIZED VIEW LOG ON ... (mview log)
// - CREATE MATERIALIZED VIEW ... (regular mview)
func (p *Parser) parseCreateMaterializedOrView(start int, orReplace bool) (nodes.StmtNode, error) {
	materialized := p.advance() // consume MATERIALIZED
	// Check for MATERIALIZED ZONEMAP
	if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
		p.advance() // consume ZONEMAP
		return p.parseCreateMaterializedZonemapStmt(start)
	}
	if p.cur.Type == kwVIEW {
		next := p.peekNext()
		if next.Type == kwLOG {
			p.advance() // consume VIEW
			p.advance() // consume LOG
			return p.parseCreateMviewLogStmt(start)
		}
	}
	// A materialized view is created, not replaced: Oracle refuses OR
	// REPLACE at MATERIALIZED (ORA-00922).
	if orReplace {
		return nil, p.syntaxErrorAtTok(materialized)
	}
	if err := p.expectCreateTableToken(kwVIEW); err != nil {
		return nil, err
	}
	return p.parseCreateMaterializedViewStmt(&nodes.CreateViewStmt{
		Materialized: true,
		Loc:          nodes.Loc{Start: start},
	})
}

// parseAdminDDLWithOptions parses generic administrative DDL statements by consuming
// the action keyword (CREATE/ALTER/DROP) and object type keyword, then parsing
// the object name and preserving remaining options until semicolon/EOF.
//
// This handles: TABLESPACE, DIRECTORY, CONTEXT, CLUSTER, DIMENSION,
// FLASHBACK ARCHIVE, JAVA, LIBRARY
func (p *Parser) parseAdminDDLWithOptions(action string, objType nodes.ObjectType, start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     action,
		ObjectType: objType,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr197 error

	stmt.Name, parseErr197 = p.parseObjectName()
	if parseErr197 != nil {
		return nil, parseErr197
	}

	if p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		tokens := p.collectDDLTokensUntil(func() bool { return false })
		if len(tokens) > 0 {
			stmt.Options = &nodes.List{Items: []nodes.Node{&nodes.DDLOption{
				Key:   "OPTIONS",
				Value: strings.Join(tokens, " "),
				Loc:   nodes.Loc{Start: optStart, End: p.prev.End},
			}}}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateSchemaStmt parses a CREATE SCHEMA statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-SCHEMA.html
//
//	CREATE SCHEMA AUTHORIZATION schema_name
//	  { create_table_statement
//	  | create_view_statement
//	  | grant_statement
//	  } ...
func (p *Parser) parseCreateSchemaStmt(start int) (*nodes.CreateSchemaStmt, error) {
	stmt := &nodes.CreateSchemaStmt{
		Stmts: &nodes.List{},
		Loc:   nodes.Loc{Start: start},
	}

	// AUTHORIZATION schema_name
	if p.isIdentLikeStr("AUTHORIZATION") {
		p.advance() // consume AUTHORIZATION
	}
	if p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.isIdentLike() {
		stmt.SchemaName = p.cur.Str
		p.advance()
	}

	// Parse nested statements: CREATE TABLE, CREATE VIEW, GRANT
	// In Oracle syntax, nested statements do NOT have their own semicolons.
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		nestedStart := p.pos()
		switch p.cur.Type {
		case kwCREATE:
			p.advance() // consume CREATE
			switch p.cur.Type {
			case kwTABLE:
				parseValue3, parseErr4 := p.parseCreateTableStmt(nestedStart, false, false, false, false, false)
				if parseErr4 != nil {
					return nil, parseErr4
				}
				stmt.Stmts.Items = append(stmt.Stmts.Items, parseValue3)
			case kwVIEW:
				parseValue5, parseErr6 := p.parseCreateViewStmt(nestedStart, false)
				if parseErr6 != nil {
					return nil, parseErr6
				}
				stmt.Stmts.Items = append(stmt.Stmts.Items, parseValue5)
			case kwFORCE:
				parseValue7, parseErr8 := p.parseCreateViewStmt(nestedStart, false)
				if parseErr8 != nil {
					return nil,

						// Unknown nested CREATE — skip to next CREATE/GRANT/semicolon
						parseErr8
				}
				stmt.Stmts.Items = append(stmt.Stmts.Items, parseValue7)
			default:
				return nil, p.syntaxErrorAtCur()
			}
		case kwGRANT:
			parseValue9, parseErr10 := p.parseGrantStmt()
			if parseErr10 != nil {

				// Unexpected token — stop parsing nested statements
				return nil, parseErr10
			}
			stmt.Stmts.Items = append(stmt.Stmts.Items, parseValue9)
		default:

			goto done
		}
	}

done:
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateAdminObject handles CREATE dispatches for admin DDL objects
// (called from parseCreateStmt after consuming CREATE and modifiers).
func (p *Parser) parseCreateAdminObject(start int, orReplace bool) (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		return p.parseCreateUserStmt(start)
	case kwROLE:
		p.advance()
		return p.parseCreateRoleStmt(start)
	case kwPROFILE:
		p.advance()
		return p.parseCreateProfileStmt(start, false)
	case kwTABLESPACE:
		return p.parseCreateTablespaceFamily(start)
	case kwDIRECTORY:
		p.advance()
		return p.parseCreateDirectoryStmt(start, orReplace)
	case kwCLUSTER:
		p.advance()
		return p.parseCreateClusterStmt(start)
	case kwJAVA:
		p.advance()
		return p.parseCreateJavaStmt(start, orReplace, nil)
	case kwSCHEMA:
		p.advance()
		return p.parseCreateSchemaStmt(start)
	default:
		// DIMENSION, FLASHBACK ARCHIVE, MANDATORY PROFILE handled via identifiers
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				if orReplace {
					// CREATE DIMENSION takes no OR REPLACE.
					return nil, p.syntaxErrorAtCur()
				}
				p.advance()
				return p.parseCreateDimensionStmt(start)
			case "FLASHBACK":
				p.advance()
				if !p.isKeywordStr("ARCHIVE") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume ARCHIVE
				return p.parseCreateFlashbackArchiveStmt(start)
			case "MANDATORY":
				p.advance()
				if p.cur.Type == kwPROFILE {
					p.advance()
					return p.parseCreateProfileStmt(start, true)
				}
			case "DISKGROUP":
				p.advance()
				return p.parseCreateDiskgroupStmt(start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type != kwDATABASE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume DATABASE
				return p.parseCreatePluggableDatabaseStmt(start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type != kwVIEW {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume VIEW
				return p.parseCreateAnalyticViewStmt(start, false, false, false)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if !p.isKeywordStr("DIMENSION") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume DIMENSION
				return p.parseCreateAttributeDimensionStmt(start, false, false, false)
			case "HIERARCHY":
				p.advance()
				return p.parseCreateHierarchyStmt(start, false, false, false)
			case "DOMAIN":
				p.advance()
				return p.parseCreateDomainStmt(start, false, false)
			case "INDEXTYPE":
				p.advance()
				return p.parseCreateIndextypeStmt(start, false)
			case "OPERATOR":
				p.advance()
				return p.parseCreateOperatorStmt(start, false, false)
			case "LOCKDOWN":
				p.advance() // consume LOCKDOWN
				if p.cur.Type != kwPROFILE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume PROFILE
				return p.parseCreateLockdownProfileStmt(start)
			case "OUTLINE":
				p.advance()
				return p.parseCreateOutlineStmt(start, false, false)
			case "EDITION":
				p.advance()
				return p.parseCreateEditionStmt(start)
			case "PFILE":
				p.advance()
				return p.parseCreatePfileStmt(start)
			case "SPFILE":
				p.advance()
				return p.parseCreateSpfileStmt(start)
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if err := p.expectWords("GRAPH"); err != nil {
					return nil, err
				}
				return p.parseCreatePropertyGraphStmt(start, false, false)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type != kwINDEX {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume INDEX
				return p.parseCreateVectorIndexStmt(start, false)
			case "LOGICAL":
				p.advance() // consume LOGICAL
				if p.cur.Type != kwPARTITION {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume PARTITION
				if err := p.expectCreateTableWord("TRACKING"); err != nil {
					return nil, err
				}
				return p.parseCreateLogicalPartitionTrackingStmt(start)
			case "PMEM":
				p.advance() // consume PMEM
				if err := p.expectCreateTableWord("FILESTORE"); err != nil {
					return nil, err
				}
				return p.parseCreatePmemFilestoreStmt(start)
			}
		}
		return nil, nil
	}
}

// parseDropAdminObject handles DROP dispatches for admin DDL objects
// (called from parseDropStmt for object types not handled there).
func (p *Parser) parseDropAdminObject(start int) (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_USER,
			Loc:        nodes.Loc{Start: start},
		}
		// IF EXISTS
		if p.cur.Type == kwIF {
			next := p.peekNext()
			if next.Type == kwEXISTS {
				p.advance() // consume IF
				p.advance() // consume EXISTS
			}
		}
		var parseErr198 error
		stmt.Name, parseErr198 = p.parseObjectName()
		if parseErr198 !=
			// CASCADE
			nil {
			return nil, parseErr198
		}

		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	case kwROLE:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_ROLE,
			Loc:        nodes.Loc{Start: start},
		}
		var err error
		stmt.Name, err = p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if stmt.Name == nil || stmt.Name.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	case kwPROFILE:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_PROFILE,
			Loc:        nodes.Loc{Start: start},
		}
		var parseErr199 error
		stmt.Name, parseErr199 = p.parseObjectName()
		if parseErr199 != nil {
			return nil, parseErr199
		}
		if stmt.Name == nil || stmt.Name.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	case kwTABLESPACE:
		p.advance()
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseDropTablespaceStmt(start, true)
		}
		return p.parseDropTablespaceStmt(start, false)
	case kwDIRECTORY:
		p.advance()
		return p.parseDropSimpleStmt(nodes.OBJECT_DIRECTORY, start)
	case kwCONTEXT:
		p.advance()
		return p.parseDropSimpleStmt(nodes.OBJECT_CONTEXT, start)
	case kwCLUSTER:
		p.advance()
		return p.parseDropClusterStmt(start)
	case kwJAVA:
		p.advance()
		return p.parseDropJavaStmt(start)
	case kwLIBRARY:
		p.advance()
		return p.parseDropSimpleStmt(nodes.OBJECT_LIBRARY, start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_DIMENSION, start)
			case "FLASHBACK":
				p.advance()
				if !p.isKeywordStr("ARCHIVE") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume ARCHIVE
				return p.parseDropSimpleStmt(nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			case "DISKGROUP":
				p.advance()
				return p.parseDropDiskgroupStmt(start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if err := p.expectWords("DATABASE"); err != nil {
					return nil, err
				}
				return p.parseDropPluggableDatabaseStmt(start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if err := p.expectWords("VIEW"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_ANALYTIC_VIEW, start)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if err := p.expectWords("DIMENSION"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_ATTRIBUTE_DIMENSION, start)
			case "HIERARCHY":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_HIERARCHY, start)
			case "DOMAIN":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_DOMAIN, start)
			case "INDEXTYPE":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_INDEXTYPE, start)
			case "OPERATOR":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_OPERATOR, start)
			case "LOCKDOWN":
				p.advance() // consume LOCKDOWN
				if p.cur.Type != kwPROFILE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume PROFILE
				return p.parseDropSimpleStmt(nodes.OBJECT_LOCKDOWN_PROFILE, start)
			case "OUTLINE":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_OUTLINE, start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if err := p.expectWords("JOIN", "GROUP"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_INMEMORY_JOIN_GROUP, start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if err := p.expectWords("SEGMENT"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_ROLLBACK_SEGMENT, start)
			case "EDITION":
				p.advance()
				return p.parseDropEditionStmt(start)
			case "MLE":
				p.advance() // consume MLE
				if p.isIdentLike() && p.cur.Str == "ENV" {
					p.advance() // consume ENV
					return p.parseDropSimpleStmt(nodes.OBJECT_MLE_ENV, start)
				}
				if p.isIdentLike() && p.cur.Str == "MODULE" {
					p.advance() // consume MODULE
					return p.parseDropSimpleStmt(nodes.OBJECT_MLE_MODULE, start)
				}
				return nil, p.syntaxErrorAtCur()
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if err := p.expectWords("GRAPH"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_PROPERTY_GRAPH, start)
			case "VECTOR":
				p.advance() // consume VECTOR
				if err := p.expectWords("INDEX"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_VECTOR_INDEX, start)
			case "RESTORE":
				p.advance() // consume RESTORE
				if err := p.expectWords("POINT"); err != nil {
					return nil, err
				}
				return p.parseDropRestorePointStmt(start)
			case "LOGICAL":
				p.advance() // consume LOGICAL
				if err := p.expectWords("PARTITION", "TRACKING"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_LOGICAL_PARTITION_TRACKING, start)
			case "PMEM":
				p.advance() // consume PMEM
				if err := p.expectWords("FILESTORE"); err != nil {
					return nil, err
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_PMEM_FILESTORE, start)
			}
		}
		return nil, nil
	}
}

// parseDatafileClause parses a single file specification (path string and optional SIZE).
func (p *Parser) parseDatafileClause() (*nodes.DatafileClause, error) {
	if p.cur.Type != tokSCONST {
		return nil, nil
	}
	df := &nodes.DatafileClause{
		Loc: nodes.Loc{Start: p.pos()},
	}
	df.Filename = p.cur.Str
	p.advance()

	// Optional SIZE
	if p.cur.Type == kwSIZE {
		p.advance()
		var parseErr230 error
		df.Size, parseErr230 = p.parseSizeValue()
		if parseErr230 !=

			// Optional REUSE
			nil {
			return nil, parseErr230
		}
	}

	if p.isIdentLike() && p.cur.Str == "REUSE" {
		p.advance()
		df.Reuse = true
	}

	// Optional AUTOEXTEND
	if p.isIdentLike() && p.cur.Str == "AUTOEXTEND" {
		var parseErr231 error
		df.Autoextend, parseErr231 = p.parseAutoextendClause()
		if parseErr231 != nil {
			return nil, parseErr231
		}
	}

	df.Loc.End = p.prev.End
	return df, nil
}

// parseAutoextendClause parses AUTOEXTEND ON/OFF with optional NEXT and MAXSIZE.
func (p *Parser) parseAutoextendClause() (*nodes.AutoextendClause, error) {
	ac := &nodes.AutoextendClause{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume AUTOEXTEND

	if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
		p.advance()
		ac.On = true
		// Optional NEXT size
		if p.cur.Type == kwNEXT {
			p.advance()
			var parseErr232 error
			ac.Next, parseErr232 = p.parseSizeValue()
			if parseErr232 !=

				// Optional MAXSIZE
				nil {
				return nil, parseErr232
			}
		}

		if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				ac.MaxSize = "UNLIMITED"
			} else {
				var parseErr233 error
				ac.MaxSize, parseErr233 = p.parseSizeValue()
				if parseErr233 != nil {
					return nil, parseErr233
				}
			}
		}
	} else if p.isIdentLike() && p.cur.Str == "OFF" {
		p.advance()
		ac.On = false
	}

	ac.Loc.End = p.prev.End
	return ac, nil
}

// parseSizeValue parses a size value like "100M", "10G", "512", "8K".
// It combines the number and optional unit suffix into a single string.
func (p *Parser) parseSizeValue() (string, error) {
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		// Check for size suffix (K, M, G, T, P, E) as an identifier
		if p.isIdentLike() {
			switch p.cur.Str {
			case "K", "M", "G", "T", "P", "E":
				val += p.cur.Str
				p.advance()
			}
		}
		return val, nil
	}
	return "", nil
}

// skipParens skips balanced parentheses.
func (p *Parser) skipParens() {
	depth := 0
	if p.cur.Type == '(' {
		depth = 1
		p.advance()
	}
	for depth > 0 && p.cur.Type != tokEOF {
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			depth--
		}
		p.advance()
	}
}

// parseIntValue parses an integer constant and returns its value.
func (p *Parser) parseIntValue() (int, error) {
	if p.cur.Type == tokICONST {
		val := 0
		for _, c := range p.cur.Str {
			val = val*10 + int(c-'0')
		}
		p.advance()
		return val, nil
	}
	return 0, nil
}

// ---------------------------------------------------------------------------
// ALTER CLUSTER
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// CREATE INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// ALTER INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// DROP CLUSTER (with INCLUDING TABLES / CASCADE CONSTRAINTS)
// ---------------------------------------------------------------------------

// parseDropClusterStmt parses a DROP CLUSTER statement.
//
// BNF: oracle/parser/bnf/DROP-CLUSTER.bnf
//
//	DROP CLUSTER [ IF EXISTS ] [ schema. ] cluster
//	    [ INCLUDING TABLES [ CASCADE CONSTRAINTS ] ]
func (p *Parser) parseDropClusterStmt(start int) (*nodes.DropStmt, error) {
	stmt := &nodes.DropStmt{
		ObjectType: nodes.OBJECT_CLUSTER,
		Names:      &nodes.List{},
		Loc:        nodes.Loc{Start: start},
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

	name, parseErr281 := p.parseObjectName()
	if parseErr281 != nil {
		return nil, parseErr281
	}
	stmt.Names.Items = append(stmt.Names.Items, name)

	// Optional INCLUDING TABLES
	if p.isIdentLike() && p.cur.Str == "INCLUDING" {
		p.advance() // INCLUDING
		if err := p.expectWords("TABLES"); err != nil {
			return nil, err
		}
		// Optional CASCADE CONSTRAINTS
		if p.cur.Type == kwCASCADE {
			p.advance() // CASCADE
			if err := p.expectWords("CONSTRAINTS"); err != nil {
				return nil, err
			}
			stmt.Cascade = true
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterAdminObject handles ALTER dispatches for admin DDL objects.
func (p *Parser) parseAlterAdminObject(start int) (nodes.StmtNode, error) {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		return p.parseAlterUserStmt(start)
	case kwROLE:
		p.advance()
		return p.parseAlterRoleStmt(start)
	case kwPROFILE:
		p.advance()
		return p.parseAlterProfileStmt(start)
	case kwTABLESPACE:
		p.advance()
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseAlterTablespaceStmt(start, true)
		}
		return p.parseAlterTablespaceStmt(start, false)
	case kwCLUSTER:
		p.advance()
		return p.parseAlterClusterStmt(start)
	case kwJAVA:
		p.advance()
		return p.parseAlterJavaStmt(start)
	case kwLIBRARY:
		p.advance()
		return p.parseAlterLibraryStmt(start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAlterDimensionStmt(start)
			case "DISKGROUP":
				p.advance()
				return p.parseAlterDiskgroupStmt(start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type != kwDATABASE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume DATABASE
				return p.parseAlterPluggableDatabaseStmt(start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type != kwVIEW {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume VIEW
				return p.parseAlterAnalyticViewStmt(start)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if !p.isKeywordStr("DIMENSION") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume DIMENSION
				return p.parseAlterAttributeDimensionStmt(start)
			case "HIERARCHY":
				p.advance()
				return p.parseAlterHierarchyStmt(start)
			case "DOMAIN":
				p.advance()
				return p.parseAlterDomainStmt(start, false)
			case "INDEXTYPE":
				p.advance()
				return p.parseAlterIndextypeStmt(start)
			case "OPERATOR":
				p.advance()
				return p.parseAlterOperatorStmt(start)
			case "LOCKDOWN":
				p.advance() // consume LOCKDOWN
				if p.cur.Type != kwPROFILE {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume PROFILE
				return p.parseAlterLockdownProfileStmt(start)
			case "OUTLINE":
				p.advance()
				return p.parseAlterOutlineStmt(start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if err := p.expectCreateTableWord("JOIN"); err != nil {
					return nil, err
				}
				if err := p.expectCreateTableWord("GROUP"); err != nil {
					return nil, err
				}
				return p.parseAlterInmemoryJoinGroupStmt(start)
			case "FLASHBACK":
				p.advance() // consume FLASHBACK
				if !p.isKeywordStr("ARCHIVE") {
					return nil, p.syntaxErrorAtCur()
				}
				p.advance() // consume ARCHIVE
				return p.parseAlterFlashbackArchiveStmt(start)
			case "RESOURCE":
				p.advance() // consume RESOURCE
				if err := p.expectWords("COST"); err != nil {
					return nil, err
				}
				return p.parseAlterResourceCostStmt(start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if err := p.expectCreateTableWord("SEGMENT"); err != nil {
					return nil, err
				}
				return p.parseAlterRollbackSegmentStmt(start)
			case "MATERIALIZED":
				p.advance() // consume MATERIALIZED (already a keyword, handled above for MVIEW)
				if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
					p.advance() // consume ZONEMAP
					return p.parseAlterMaterializedZonemapStmt(start)
				}
				return nil, nil
			}
		}
		return nil, nil
	}
}

// ---------------------------------------------------------------------------
// CREATE DOMAIN
// ---------------------------------------------------------------------------

// parseParenIdentList parses ( ident, ident, ... ) and returns the list of identifiers.
func (p *Parser) parseParenIdentList() ([]string, error) {
	var result []string
	if p.cur.Type != '(' {
		return result, nil
	}
	p.advance() // consume (
	for progress := -1; p.cur.Type != ')' && p.cur.Type != tokEOF; {
		// An item that consumes nothing would come round again.
		if p.pos() == progress {
			return nil, p.syntaxErrorAtCur()
		}
		progress = p.pos()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			result = append(result, p.cur.Str)
			p.advance()
		} else {
			p.advance()
		}
		if p.cur.Type == ',' {
			p.advance()
		}
	}
	if p.cur.Type == ')' {
		p.advance()
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// CREATE VECTOR INDEX
// ---------------------------------------------------------------------------

// parseCreateVectorIndexStmt parses a CREATE VECTOR INDEX statement.
// Called after CREATE VECTOR INDEX has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-vector-index.html
//
//	CREATE VECTOR INDEX [ IF NOT EXISTS ] [ schema. ] index_name
//	    ON [ schema. ] table_name ( column_name )
//	    [ INCLUDE ( column_name [, column_name ]... ) ]
//	    [ vector_index_organization_clause ]
//	    [ DISTANCE metric_name ]
//	    [ WITH TARGET ACCURACY integer [ PARAMETERS ( vector_index_parameters_clause ) ] ]
//	    [ vector_index_hnsw_replication_clause ]
//	    [ ONLINE ]
//	    [ PARALLEL integer ] ;
//
// Where Oracle answers otherwise than the page: it takes DISTANCE and WITH
// TARGET ACCURACY in either order, but neither after PARAMETERS
// (ORA-02158); ONLINE,
// PARALLEL and LOCAL in any order; each clause and each parameter once
// (ORA-51918). Of the parameters it takes NEIGHBORS, M and EFCONSTRUCTION
// for HNSW and NEIGHBOR PARTITIONS, SAMPLES_PER_PARTITION and
// MIN_VECTORS_PER_PARTITION for IVF; RESCORE_FACTOR and QUANTIZATION are
// ORA-00922, a parameter of the other type ORA-51916.
//
// Oracle requires the organization clause the page shows as optional:
// without it the statement is ORA-51914 (ORGANIZATION clause required).
func (p *Parser) parseCreateVectorIndexStmt(start int, ifNotExists bool) (*nodes.CreateVectorIndexStmt, error) {
	stmt := &nodes.CreateVectorIndexStmt{
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}

	// IF NOT EXISTS: IF NOT without EXISTS is ORA-11543 at IF.
	if !ifNotExists && p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		ifTok := p.advance()
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance()
		stmt.IfNotExists = true
	}

	var err error
	if stmt.Name, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}
	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if stmt.TableName, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if stmt.Column, err = p.ixName(); err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// INCLUDE ( column_name [, column_name ]... )
	if p.cur.Type == kwINCLUDE {
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		for {
			col, err := p.ixName()
			if err != nil {
				return nil, err
			}
			stmt.IncludeColumns = append(stmt.IncludeColumns, col)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != ')' {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	// ORGANIZATION { INMEMORY [ NEIGHBOR ] GRAPH | [ NEIGHBOR ] PARTITIONS }:
	// Oracle takes INMEMORY GRAPH without NEIGHBOR, and PARTITIONS without it
	// too (measured on the pinned server, where the page shows NEIGHBOR
	// PARTITIONS alone).
	if err := p.ixWord("ORGANIZATION"); err != nil {
		return nil, err
	}
	switch {
	case p.isKeywordStr("INMEMORY"):
		p.advance()
		if p.isKeywordStr("NEIGHBOR") {
			p.advance()
		}
		if err := p.ixWord("GRAPH"); err != nil {
			return nil, err
		}
		stmt.Organization = "INMEMORY_NEIGHBOR_GRAPH"
	case p.isKeywordStr("NEIGHBOR"):
		p.advance()
		if err := p.ixWord("PARTITIONS"); err != nil {
			return nil, err
		}
		stmt.Organization = "NEIGHBOR_PARTITIONS"
	case p.isKeywordStr("PARTITIONS"):
		p.advance()
		stmt.Organization = "NEIGHBOR_PARTITIONS"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	seen := map[string]bool{}
	for {
		if p.isKeywordStr("DISTANCE") {
			if err := p.ixOnce(seen, "DISTANCE"); err != nil {
				return nil, err
			}
			p.advance()
			if p.cur.Type == tokQIDENT || !p.isIxName() {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Distance = p.advance().Str
			continue
		}
		if p.cur.Type == kwWITH {
			if err := p.ixOnce(seen, "ACCURACY"); err != nil {
				return nil, err
			}
			p.advance()
			if err := p.ixWord("TARGET"); err != nil {
				return nil, err
			}
			if err := p.ixWord("ACCURACY"); err != nil {
				return nil, err
			}
			if p.cur.Type == tokICONST {
				v, err := p.parseIntValue()
				if err != nil {
					return nil, err
				}
				stmt.TargetAccuracy = v
			} else {
				return nil, p.syntaxErrorAtCur()
			}
			if p.isKeywordStr("PARAMETERS") {
				if err := p.parseVectorIndexParameters(stmt); err != nil {
					return nil, err
				}
				break
			}
			continue
		}
		break
	}

	for !p.ixEnd() {
		switch {
		case p.isKeywordStr("DUPLICATE"):
			if err := p.ixOnce(seen, "REPLICATION"); err != nil {
				return nil, err
			}
			p.advance()
			if p.cur.Type != kwALL {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			stmt.Replication = "DUPLICATE_ALL"
		case p.isKeywordStr("DISTRIBUTE"):
			if err := p.ixOnce(seen, "REPLICATION"); err != nil {
				return nil, err
			}
			p.advance()
			stmt.Replication = "DISTRIBUTE"
			if p.isKeywordStr("AUTO") {
				p.advance()
				stmt.Replication = "DISTRIBUTE_AUTO"
			} else if p.cur.Type == kwBY {
				p.advance()
				switch {
				case p.isKeywordStr("ROWID"):
					p.advance()
					if err := p.ixWord("RANGE"); err != nil {
						return nil, err
					}
					stmt.Replication = "DISTRIBUTE_BY_ROWID_RANGE"
				case p.cur.Type == kwPARTITION:
					p.advance()
					stmt.Replication = "DISTRIBUTE_BY_PARTITION"
				case p.cur.Type == kwSUBPARTITION:
					p.advance()
					stmt.Replication = "DISTRIBUTE_BY_SUBPARTITION"
				default:
					return nil, p.syntaxErrorAtCur()
				}
			}
		case p.cur.Type == kwONLINE:
			if err := p.ixOnce(seen, "ONLINE"); err != nil {
				return nil, err
			}
			p.advance()
			stmt.Online = true
		case p.cur.Type == kwPARALLEL:
			if err := p.ixOnce(seen, "PARALLEL"); err != nil {
				return nil, err
			}
			p.advance()
			if p.cur.Type == tokICONST {
				v, err := p.parseIntValue()
				if err != nil {
					return nil, err
				}
				stmt.Parallel = v
			} else {
				return nil, p.syntaxErrorAtCur()
			}
		case p.cur.Type == kwLOCAL:
			if err := p.ixOnce(seen, "LOCAL"); err != nil {
				return nil, err
			}
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseVectorIndexParameters reads PARAMETERS ( vector_index_parameters_clause ):
//
//	TYPE HNSW [, { NEIGHBORS | M } integer ] [, EFCONSTRUCTION integer ]
//	TYPE IVF [, NEIGHBOR PARTITIONS integer ] [, SAMPLES_PER_PARTITION integer ]
//	         [, MIN_VECTORS_PER_PARTITION integer ]
//
// See parseCreateVectorIndexStmt for the parameters Oracle refuses.
func (p *Parser) parseVectorIndexParameters(stmt *nodes.CreateVectorIndexStmt) error {
	p.advance() // PARAMETERS
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwTYPE {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.isKeywordStr("HNSW"), p.isKeywordStr("IVF"):
		stmt.ParameterType = p.advance().Str
	default:
		return p.syntaxErrorAtCur()
	}
	hnsw := stmt.ParameterType == "HNSW"
	seen := map[string]bool{}
	for p.cur.Type == ',' {
		p.advance()
		var key string
		var target *int
		switch {
		case hnsw && (p.isKeywordStr("NEIGHBORS") || p.isKeywordStr("M")):
			key, target = "NEIGHBORS", &stmt.Neighbors
			p.advance()
		case hnsw && p.isKeywordStr("EFCONSTRUCTION"):
			key, target = "EFCONSTRUCTION", &stmt.EfConstruction
			p.advance()
		case !hnsw && p.isKeywordStr("NEIGHBOR"):
			p.advance()
			if err := p.ixWord("PARTITIONS"); err != nil {
				return err
			}
			key, target = "NEIGHBOR PARTITIONS", &stmt.NeighborParts
		case !hnsw && p.isKeywordStr("SAMPLES_PER_PARTITION"):
			key, target = "SAMPLES_PER_PARTITION", &stmt.SamplesPerPart
			p.advance()
		case !hnsw && p.isKeywordStr("MIN_VECTORS_PER_PARTITION"):
			key, target = "MIN_VECTORS_PER_PARTITION", &stmt.MinVecsPerPart
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
		if seen[key] {
			return p.syntaxErrorAtTok(p.prev)
		}
		seen[key] = true
		if p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		v, err := p.parseIntValue()
		if err != nil {
			return err
		}
		*target = v
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// ---------------------------------------------------------------------------
// CREATE / ALTER OUTLINE
// ---------------------------------------------------------------------------

// parseCreateOutlineStmt parses a CREATE OUTLINE statement.
// Called after CREATE [OR REPLACE] [PUBLIC|PRIVATE] OUTLINE has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE ] [ PUBLIC | PRIVATE ] OUTLINE [ outline ]
//	    [ FOR CATEGORY category ]
//	    { FROM [ PRIVATE ] source_outline | ON statement } ;
func (p *Parser) parseCreateOutlineStmt(start int, orReplace, public bool) (*nodes.CreateOutlineStmt, error) {
	stmt := &nodes.CreateOutlineStmt{
		OrReplace: orReplace,
		Public:    public,
		Loc:       nodes.Loc{Start: start},
	}

	// Optional outline name: one name (ORA-00969 at a dot after it).
	if p.cur.Type != kwFOR && p.cur.Type != kwFROM && p.cur.Type != kwON && p.isIdentLike() {
		var err error
		if stmt.Name, err = p.parseIdentifier(); err != nil {
			return nil, err
		}
	}

	switch p.cur.Type {
	case kwFROM:
		// FROM [ PRIVATE ] source_outline
		p.advance()
		if p.cur.Type == kwPRIVATE {
			stmt.FromPrivate = true
			p.advance()
		}
		var err error
		if stmt.FromSource, err = p.parseRequiredIdentifier(); err != nil {
			return nil, err
		}
	case kwFOR, kwON:
		// [ FOR CATEGORY category ] ON statement
		if p.cur.Type == kwFOR {
			p.advance()
			if err := p.expectWords("CATEGORY"); err != nil {
				return nil, err
			}
			var err error
			if stmt.Category, err = p.parseRequiredIdentifier(); err != nil {
				return nil, err
			}
			if p.cur.Type != kwON {
				return nil, p.syntaxErrorAtCur()
			}
		}
		// ON statement — skip to semicolon
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRequiredIdentifier parses one name, failing at the current token when
// none is there or when it is a reserved word (ORA-00931 at PUBLIC).
func (p *Parser) parseRequiredIdentifier() (string, error) {
	if p.cur.Type != tokQIDENT && isOracleSQLReservedKeyword(p.cur) {
		return "", p.syntaxErrorAtCur()
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return "", err
	}
	if name == "" {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// parseAlterOutlineStmt parses an ALTER OUTLINE statement.
// Called after ALTER OUTLINE has been consumed.
//
// BNF:
//
//	ALTER OUTLINE outline
//	    { REBUILD | RENAME TO new_outline_name | CHANGE CATEGORY TO new_category_name | { ENABLE | DISABLE } }... ;
//
// Oracle 23ai takes no PUBLIC or PRIVATE here (ORA-00931 at PUBLIC; it
// reads PRIVATE as the name), and more than one action (REBUILD ENABLE).
// Action records the first.
func (p *Parser) parseAlterOutlineStmt(start int) (*nodes.AlterOutlineStmt, error) {
	stmt := &nodes.AlterOutlineStmt{
		Loc: nodes.Loc{Start: start},
	}

	var err error
	if stmt.Name, err = p.parseRequiredIdentifier(); err != nil {
		return nil, err
	}

	for {
		action := ""
		switch {
		case p.isKeywordStr("REBUILD"):
			action = "REBUILD"
			p.advance()
		case p.cur.Type == kwRENAME:
			action = "RENAME"
			p.advance()
			if p.cur.Type != kwTO {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if stmt.NewName, err = p.parseRequiredIdentifier(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("CHANGE"):
			action = "CHANGE_CATEGORY"
			p.advance()
			if err := p.expectWords("CATEGORY"); err != nil {
				return nil, err
			}
			if p.cur.Type != kwTO {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if stmt.Category, err = p.parseRequiredIdentifier(); err != nil {
				return nil, err
			}
		case p.cur.Type == kwENABLE:
			action = "ENABLE"
			p.advance()
		case p.cur.Type == kwDISABLE:
			action = "DISABLE"
			p.advance()
		}
		if action == "" {
			break
		}
		if stmt.Action == "" {
			stmt.Action = action
		}
	}
	if stmt.Action == "" {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Batch 105 — small objects bundle
// ---------------------------------------------------------------------------

// parseSourceText records the source text after the current AS of a CREATE
// JAVA or CREATE MLE MODULE statement and moves the parser to the end of the
// parsed range, without lexing the Java or JavaScript text. The text runs from
// its first non-blank byte through the end of the range, trailing blanks
// included, as Split keeps it: Split already leaves out the line break and
// "/" that frame it.
func (p *Parser) parseSourceText(opts *nodes.List) error {
	asTok := p.cur
	srcStart := skipSpace(p.source, asTok.End, p.lexer.end)
	srcEnd := p.lexer.end
	if srcStart >= trimRightSpace(p.source, srcEnd) {
		p.advance() // consume AS; the source is missing
		return p.syntaxErrorAtCur()
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{
		Key:   "AS",
		Value: p.source[srcStart:srcEnd],
		Loc:   nodes.Loc{Start: asTok.Loc, End: srcEnd},
	})
	p.prev = Token{Type: tokIDENT, Loc: srcStart, End: srcEnd}
	p.hasNext = false
	p.lexer.pos = p.lexer.end
	p.cur = Token{Type: tokEOF, Loc: p.lexer.end, End: p.lexer.end}
	return nil
}

// skipSpace returns the first offset in source[pos:end] that is not
// whitespace, or end.
func skipSpace(source string, pos, end int) int {
	for pos < end {
		switch source[pos] {
		case ' ', '\t', '\n', '\r', '\f':
			pos++
		default:
			return pos
		}
	}
	return pos
}

// parseCreateDirectoryStmt parses a CREATE DIRECTORY statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] DIRECTORY has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] DIRECTORY directory
//	    [ SHARING = { METADATA | NONE } ]
//	    AS 'path_name' ;
func (p *Parser) parseCreateDirectoryStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_DIRECTORY,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	ifNotExists, err := p.parseIfNotExistsAfterKind(orReplace)
	if err != nil {
		return nil, err
	}
	if ifNotExists {
		stmt.IfNotExists = true
	}
	var parseErr387 error

	stmt.Name, parseErr387 = p.parseObjectName()
	if parseErr387 != nil {
		return nil, parseErr387
	}

	// [ SHARING = { METADATA | NONE } ]
	if p.isKeywordStr("SHARING") {
		p.advance()
		value, err := p.parseSharingValue()
		if err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHARING", Value: value})
	}

	// AS 'path_name'
	if p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS", Value: p.cur.Str})
	p.advance()

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreatePfileStmt parses a CREATE PFILE statement.
// Called after CREATE PFILE has been consumed.
//
// BNF:
//
//	CREATE PFILE [ = 'pfile_name' ]
//	    FROM { SPFILE [ = 'spfile_name' ] | MEMORY } ;
func (p *Parser) parseCreatePfileStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PFILE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ = 'pfile_name' ] or [ 'pfile_name' ]
	if p.cur.Type == '=' {
		p.advance()
	}
	if p.cur.Type == tokSCONST {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PFILE", Value: p.cur.Str})
		p.advance()
	}

	seenFrom := false
	if p.cur.Type == kwFROM {
		seenFrom = true
		p.advance()
	}

	// { SPFILE [ = 'spfile_name' ] | MEMORY }
	if p.isIdentLike() && p.cur.Str == "SPFILE" {
		p.advance()
		val := ""
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM", Value: "SPFILE"})
		if val != "" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SPFILE", Value: val})
		}
	} else if p.isIdentLike() && p.cur.Str == "MEMORY" {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM", Value: "MEMORY"})
		p.advance()
	} else if seenFrom {
		return nil, p.syntaxErrorAtCur()
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateSpfileStmt parses a CREATE SPFILE statement.
// Called after CREATE SPFILE has been consumed.
//
// BNF:
//
//	CREATE SPFILE [ = 'spfile_name' ]
//	    FROM { PFILE [ = 'pfile_name' ] [ AS COPY ] | MEMORY } ;
func (p *Parser) parseCreateSpfileStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_SPFILE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ = 'spfile_name' ]
	if p.cur.Type == '=' {
		p.advance()
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SPFILE", Value: p.cur.Str})
		p.advance()
	}

	// FROM
	if p.cur.Type != kwFROM {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// { PFILE [ = 'pfile_name' ] [ AS COPY ] | MEMORY }
	if p.isIdentLike() && p.cur.Str == "PFILE" {
		p.advance()
		val := ""
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokSCONST {
			val = p.cur.Str
			p.advance()
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM", Value: "PFILE"})
		if val != "" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PFILE", Value: val})
		}
		// [ AS COPY ]
		if p.cur.Type == kwAS {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "COPY" {
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS COPY"})
			}
		}
	} else if p.isIdentLike() && p.cur.Str == "MEMORY" {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FROM", Value: "MEMORY"})
		p.advance()
	} else {
		return nil, p.syntaxErrorAtCur()
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateEditionStmt parses a CREATE EDITION statement.
// Called after CREATE [IF NOT EXISTS] EDITION has been consumed.
//
// BNF:
//
//	CREATE EDITION [ IF NOT EXISTS ] edition
//	    [ AS CHILD OF parent_edition ] ;
func (p *Parser) parseCreateEditionStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_EDITION,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ IF NOT EXISTS ]
	var err error
	if stmt.IfNotExists, err = p.parseIfNotExistsAfterKind(false); err != nil {
		return nil, err
	}

	if stmt.Name, err = p.parseRequiredObjectName(); err != nil {
		return nil, err
	}

	// [ AS CHILD OF parent_edition ]
	if p.cur.Type == kwAS {
		p.advance()
		if err := p.expectWords("CHILD"); err != nil {
			return nil, err
		}
		if p.cur.Type != kwOF {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		parent, err := p.parseRequiredIdentifier()
		if err != nil {
			return nil, err
		}
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS CHILD OF", Value: parent})
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropEditionStmt parses a DROP EDITION statement.
// Called after DROP EDITION has been consumed.
//
// BNF:
//
//	DROP EDITION [ IF EXISTS ] edition [ CASCADE ] ;
func (p *Parser) parseDropEditionStmt(start int) (*nodes.DropStmt, error) {
	stmt := &nodes.DropStmt{
		ObjectType: nodes.OBJECT_EDITION,
		Names:      &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// [ IF EXISTS ]
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance()
			p.advance()
			stmt.IfExists = true
		}
	}

	name, parseErr401 := p.parseObjectName()
	if parseErr401 != nil {
		return nil, parseErr401
	}
	if name != nil {
		stmt.Names.Items = append(stmt.Names.Items, name)
	}

	// [ CASCADE ]
	if p.cur.Type == kwCASCADE {
		stmt.Cascade = true
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropJavaStmt parses a DROP JAVA statement.
// Called after DROP JAVA has been consumed.
//
// BNF:
//
//	DROP JAVA { SOURCE | CLASS | RESOURCE } [ IF EXISTS ] [ schema. ] object_name ;
func (p *Parser) parseDropJavaStmt(start int) (*nodes.DropStmt, error) {
	stmt := &nodes.DropStmt{
		ObjectType: nodes.OBJECT_JAVA,
		Names:      &nodes.List{},
		Loc:        nodes.Loc{Start: start},
	}

	// { SOURCE | CLASS | RESOURCE }
	if !p.isKeywordStr("SOURCE") && !p.isKeywordStr("CLASS") && !p.isKeywordStr("RESOURCE") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// [ IF EXISTS ]
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance()
			p.advance()
			stmt.IfExists = true
		}
	}

	name, parseErr409 := p.parseObjectName()
	if parseErr409 != nil {
		return nil, parseErr409
	}
	if name != nil {
		stmt.Names.Items = append(stmt.Names.Items, name)
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
