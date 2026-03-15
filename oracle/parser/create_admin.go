package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateMaterializedOrView distinguishes between:
// - CREATE MATERIALIZED VIEW LOG ON ... (mview log)
// - CREATE MATERIALIZED VIEW ... (regular mview)
func (p *Parser) parseCreateMaterializedOrView(start int, orReplace bool) nodes.StmtNode {
	p.advance() // consume MATERIALIZED
	// Check for MATERIALIZED ZONEMAP
	if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
		p.advance() // consume ZONEMAP
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_MATERIALIZED_ZONEMAP, start)
	}
	if p.cur.Type == kwVIEW {
		next := p.peekNext()
		if next.Type == kwLOG {
			p.advance() // consume VIEW
			p.advance() // consume LOG
			return p.parseCreateMviewLogStmt(start)
		}
	}
	// It's a regular MATERIALIZED VIEW — but we already consumed MATERIALIZED.
	// parseCreateViewStmt expects to see kwMATERIALIZED. Since we consumed it,
	// we need to handle VIEW directly.
	stmt := &nodes.CreateViewStmt{
		OrReplace:    orReplace,
		Materialized: true,
		Loc:          nodes.Loc{Start: start},
	}
	if p.cur.Type == kwVIEW {
		p.advance()
	}
	// Delegate the rest to the existing view parsing logic.
	return p.finishCreateViewStmt(stmt)
}

// parseAdminDDLStmt parses generic administrative DDL statements by consuming
// the action keyword (CREATE/ALTER/DROP) and object type keyword, then parsing
// the object name and skipping remaining options until semicolon/EOF.
//
// This handles: TABLESPACE, DIRECTORY, CONTEXT, CLUSTER, DIMENSION,
// FLASHBACK ARCHIVE, JAVA, LIBRARY
func (p *Parser) parseAdminDDLStmt(action string, objType nodes.ObjectType, start int) nodes.StmtNode {
	stmt := &nodes.AdminDDLStmt{
		Action:     action,
		ObjectType: objType,
		Loc:        nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Skip remaining tokens until ;/EOF
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
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
func (p *Parser) parseCreateSchemaStmt(start int) *nodes.CreateSchemaStmt {
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
				stmt.Stmts.Items = append(stmt.Stmts.Items, p.parseCreateTableStmt(nestedStart, false, false, false, false, false))
			case kwVIEW:
				stmt.Stmts.Items = append(stmt.Stmts.Items, p.parseCreateViewStmt(nestedStart, false))
			case kwFORCE:
				stmt.Stmts.Items = append(stmt.Stmts.Items, p.parseCreateViewStmt(nestedStart, false))
			default:
				// Unknown nested CREATE — skip to next CREATE/GRANT/semicolon
				p.skipToSemicolon()
				goto done
			}
		case kwGRANT:
			stmt.Stmts.Items = append(stmt.Stmts.Items, p.parseGrantStmt())
		default:
			// Unexpected token — stop parsing nested statements
			goto done
		}
	}

done:
	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateAdminObject handles CREATE dispatches for admin DDL objects
// (called from parseCreateStmt after consuming CREATE and modifiers).
func (p *Parser) parseCreateAdminObject(start int) nodes.StmtNode {
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
		p.advance()
		// Check for TABLESPACE SET
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_TABLESPACE_SET, start)
		}
		return p.parseCreateTablespaceStmt(start, false, false, false, false)
	case kwDIRECTORY:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_DIRECTORY, start)
	case kwCONTEXT:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_CONTEXT, start)
	case kwCLUSTER:
		p.advance()
		return p.parseCreateClusterStmt(start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_LIBRARY, start)
	case kwSCHEMA:
		p.advance()
		return p.parseCreateSchemaStmt(start)
	default:
		// DIMENSION, FLASHBACK ARCHIVE, MANDATORY PROFILE handled via identifiers
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseCreateDimensionStmt(start)
			case "FLASHBACK":
				p.advance()
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance()
				}
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			case "MANDATORY":
				p.advance()
				if p.cur.Type == kwPROFILE {
					p.advance()
					return p.parseCreateProfileStmt(start, true)
				}
			case "DISKGROUP":
				p.advance()
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_DISKGROUP, start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_PLUGGABLE_DATABASE, start)
		case "ANALYTIC":
			p.advance() // consume ANALYTIC
			if p.cur.Type == kwVIEW {
				p.advance() // consume VIEW
			}
			return p.parseCreateAnalyticViewStmt(start, false, false, false)
		case "ATTRIBUTE":
			p.advance() // consume ATTRIBUTE
			if p.isIdentLike() && p.cur.Str == "DIMENSION" {
				p.advance() // consume DIMENSION
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_ATTRIBUTE_DIMENSION, start)
		case "HIERARCHY":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_HIERARCHY, start)
		case "DOMAIN":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_DOMAIN, start)
		case "INDEXTYPE":
			p.advance()
			return p.parseCreateIndextypeStmt(start, false)
		case "OPERATOR":
			p.advance()
			return p.parseCreateOperatorStmt(start, false, false)
		case "LOCKDOWN":
			p.advance() // consume LOCKDOWN
			if p.cur.Type == kwPROFILE {
				p.advance() // consume PROFILE
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_LOCKDOWN_PROFILE, start)
		case "OUTLINE":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_OUTLINE, start)
		case "INMEMORY":
			p.advance() // consume INMEMORY
			if p.cur.Type == kwJOIN {
				p.advance() // consume JOIN
			}
			if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
				p.advance() // consume GROUP
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_INMEMORY_JOIN_GROUP, start)
		case "ROLLBACK":
			p.advance() // consume ROLLBACK
			if p.isIdentLike() && p.cur.Str == "SEGMENT" {
				p.advance() // consume SEGMENT
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_ROLLBACK_SEGMENT, start)
		case "EDITION":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_EDITION, start)
		case "MLE":
			p.advance() // consume MLE
			if p.isIdentLike() && p.cur.Str == "ENV" {
				p.advance() // consume ENV
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_MLE_ENV, start)
			}
			if p.isIdentLike() && p.cur.Str == "MODULE" {
				p.advance() // consume MODULE
				return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_MLE_MODULE, start)
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_MLE_ENV, start)
		case "PFILE":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_PFILE, start)
		case "SPFILE":
			p.advance()
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_SPFILE, start)
		case "PROPERTY":
			p.advance() // consume PROPERTY
			if p.isIdentLike() && p.cur.Str == "GRAPH" {
				p.advance() // consume GRAPH
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_PROPERTY_GRAPH, start)
		case "VECTOR":
			p.advance() // consume VECTOR
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_VECTOR_INDEX, start)
		case "RESTORE":
			p.advance() // consume RESTORE
			if p.isIdentLike() && p.cur.Str == "POINT" {
				p.advance() // consume POINT
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_RESTORE_POINT, start)
		case "LOGICAL":
			p.advance() // consume LOGICAL
			if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
				p.advance() // consume PARTITION
			}
			if p.isIdentLike() && p.cur.Str == "TRACKING" {
				p.advance() // consume TRACKING
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_LOGICAL_PARTITION_TRACKING, start)
		case "PMEM":
			p.advance() // consume PMEM
			if p.isIdentLike() && p.cur.Str == "FILESTORE" {
				p.advance() // consume FILESTORE
			}
			return p.parseAdminDDLStmt("CREATE", nodes.OBJECT_PMEM_FILESTORE, start)
		}
	}
	return nil
	}
}

// parseDropAdminObject handles DROP dispatches for admin DDL objects
// (called from parseDropStmt for object types not handled there).
func (p *Parser) parseDropAdminObject(start int) nodes.StmtNode {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_USER,
			Loc:        nodes.Loc{Start: start},
		}
		stmt.Name = p.parseObjectName()
		// CASCADE
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.pos()
		return stmt
	case kwROLE:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ROLE, start)
	case kwPROFILE:
		p.advance()
		stmt := &nodes.AdminDDLStmt{
			Action:     "DROP",
			ObjectType: nodes.OBJECT_PROFILE,
			Loc:        nodes.Loc{Start: start},
		}
		stmt.Name = p.parseObjectName()
		if p.cur.Type == kwCASCADE {
			p.advance()
		}
		stmt.Loc.End = p.pos()
		return stmt
	case kwTABLESPACE:
		p.advance()
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseAdminDDLStmt("DROP", nodes.OBJECT_TABLESPACE_SET, start)
		}
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_TABLESPACE, start)
	case kwDIRECTORY:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DIRECTORY, start)
	case kwCONTEXT:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_CONTEXT, start)
	case kwCLUSTER:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_CLUSTER, start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_LIBRARY, start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DIMENSION, start)
			case "FLASHBACK":
				p.advance()
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance()
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			case "DISKGROUP":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DISKGROUP, start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_PLUGGABLE_DATABASE, start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type == kwVIEW {
					p.advance() // consume VIEW
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ANALYTIC_VIEW, start)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if p.isIdentLike() && p.cur.Str == "DIMENSION" {
					p.advance() // consume DIMENSION
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ATTRIBUTE_DIMENSION, start)
			case "HIERARCHY":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_HIERARCHY, start)
			case "DOMAIN":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_DOMAIN, start)
			case "INDEXTYPE":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_INDEXTYPE, start)
			case "OPERATOR":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_OPERATOR, start)
			case "LOCKDOWN":
				p.advance() // consume LOCKDOWN
				if p.cur.Type == kwPROFILE {
					p.advance() // consume PROFILE
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_LOCKDOWN_PROFILE, start)
			case "OUTLINE":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_OUTLINE, start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if p.cur.Type == kwJOIN {
					p.advance() // consume JOIN
				}
				if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
					p.advance() // consume GROUP
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_INMEMORY_JOIN_GROUP, start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ROLLBACK_SEGMENT, start)
			case "EDITION":
				p.advance()
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_EDITION, start)
			case "MLE":
				p.advance() // consume MLE
				if p.isIdentLike() && p.cur.Str == "ENV" {
					p.advance() // consume ENV
					return p.parseAdminDDLStmt("DROP", nodes.OBJECT_MLE_ENV, start)
				}
				if p.isIdentLike() && p.cur.Str == "MODULE" {
					p.advance() // consume MODULE
					return p.parseAdminDDLStmt("DROP", nodes.OBJECT_MLE_MODULE, start)
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_MLE_ENV, start)
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if p.isIdentLike() && p.cur.Str == "GRAPH" {
					p.advance() // consume GRAPH
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_PROPERTY_GRAPH, start)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type == kwINDEX {
					p.advance() // consume INDEX
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_VECTOR_INDEX, start)
			case "RESTORE":
				p.advance() // consume RESTORE
				if p.isIdentLike() && p.cur.Str == "POINT" {
					p.advance() // consume POINT
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_RESTORE_POINT, start)
			case "LOGICAL":
				p.advance() // consume LOGICAL
				if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
					p.advance() // consume PARTITION
				}
				if p.isIdentLike() && p.cur.Str == "TRACKING" {
					p.advance() // consume TRACKING
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_LOGICAL_PARTITION_TRACKING, start)
			case "PMEM":
				p.advance() // consume PMEM
				if p.isIdentLike() && p.cur.Str == "FILESTORE" {
					p.advance() // consume FILESTORE
				}
				return p.parseAdminDDLStmt("DROP", nodes.OBJECT_PMEM_FILESTORE, start)
			}
		}
		return nil
	}
}

// parseCreateTablespaceStmt parses a CREATE TABLESPACE statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-TABLESPACE.html
//
//	CREATE [ BIGFILE | SMALLFILE ]
//	  [ { TEMPORARY | UNDO } ] TABLESPACE tablespace_name
//	  { DATAFILE | TEMPFILE } file_specification [, ...]
//	  [ SIZE size_clause ]
//	  [ REUSE ]
//	  [ AUTOEXTEND { ON [ NEXT size_clause ] [ MAXSIZE { size_clause | UNLIMITED } ] | OFF } ]
//	  [ BLOCKSIZE integer [ K ] ]
//	  [ LOGGING | NOLOGGING | FORCE LOGGING ]
//	  [ ONLINE | OFFLINE ]
//	  [ EXTENT MANAGEMENT LOCAL [ AUTOALLOCATE | UNIFORM [ SIZE size_clause ] ] ]
//	  [ SEGMENT SPACE MANAGEMENT { AUTO | MANUAL } ]
//	  [ RETENTION { GUARANTEE | NOGUARANTEE } ]
//	  [ DEFAULT { COMPRESS | NOCOMPRESS | table_compression } ]
//	  [ ENCRYPTION tablespace_encryption_clause ]
//	  [ MAXSIZE { size_clause | UNLIMITED } ]
func (p *Parser) parseCreateTablespaceStmt(start int, bigfile, smallfile, temporary, undo bool) *nodes.CreateTablespaceStmt {
	stmt := &nodes.CreateTablespaceStmt{
		Loc:       nodes.Loc{Start: start},
		Bigfile:   bigfile,
		Smallfile: smallfile,
		Temporary: temporary,
		Undo:      undo,
	}

	// Parse tablespace name
	stmt.Name = p.parseObjectName()

	// Parse clauses in any order until ; or EOF
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.isIdentLike() && (p.cur.Str == "DATAFILE" || p.cur.Str == "TEMPFILE"):
			p.advance()
			// Parse one or more file specifications
			for {
				df := p.parseDatafileClause()
				if df != nil {
					stmt.Datafiles = append(stmt.Datafiles, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}

		case p.cur.Type == kwSIZE:
			p.advance()
			stmt.Size = p.parseSizeValue()

		case p.isIdentLike() && p.cur.Str == "AUTOEXTEND":
			stmt.Autoextend = p.parseAutoextendClause()

		case p.isIdentLike() && p.cur.Str == "REUSE":
			p.advance()
			// Attach REUSE to last datafile if exists
			if len(stmt.Datafiles) > 0 {
				stmt.Datafiles[len(stmt.Datafiles)-1].Reuse = true
			}

		case p.cur.Type == kwLOGGING:
			p.advance()
			stmt.Logging = "LOGGING"

		case p.cur.Type == kwNOLOGGING:
			p.advance()
			stmt.Logging = "NOLOGGING"

		case p.isIdentLike() && p.cur.Str == "FORCE":
			p.advance()
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			stmt.Logging = "FORCE LOGGING"

		case p.cur.Type == kwONLINE:
			p.advance()
			stmt.Online = true

		case p.cur.Type == kwOFFLINE:
			p.advance()
			stmt.Offline = true

		case p.isIdentLike() && p.cur.Str == "EXTENT":
			p.advance() // EXTENT
			if p.isIdentLike() && p.cur.Str == "MANAGEMENT" {
				p.advance() // MANAGEMENT
			}
			if p.isIdentLike() && p.cur.Str == "LOCAL" {
				p.advance() // LOCAL
			}
			if p.isIdentLike() && p.cur.Str == "AUTOALLOCATE" {
				p.advance()
				stmt.Extent = "AUTOALLOCATE"
			} else if p.isIdentLike() && p.cur.Str == "UNIFORM" {
				p.advance()
				if p.cur.Type == kwSIZE {
					p.advance()
					stmt.Extent = "UNIFORM SIZE " + p.parseSizeValue()
				} else {
					stmt.Extent = "UNIFORM"
				}
			} else {
				stmt.Extent = "LOCAL"
			}

		case p.isIdentLike() && p.cur.Str == "SEGMENT":
			p.advance() // SEGMENT
			if p.isIdentLike() && p.cur.Str == "SPACE" {
				p.advance() // SPACE
			}
			if p.isIdentLike() && p.cur.Str == "MANAGEMENT" {
				p.advance() // MANAGEMENT
			}
			if p.isIdentLike() && p.cur.Str == "AUTO" {
				p.advance()
				stmt.Segment = "AUTO"
			} else if p.isIdentLike() && p.cur.Str == "MANUAL" {
				p.advance()
				stmt.Segment = "MANUAL"
			}

		case p.isIdentLike() && p.cur.Str == "BLOCKSIZE":
			p.advance()
			stmt.Blocksize = p.parseSizeValue()

		case p.isIdentLike() && p.cur.Str == "RETENTION":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "GUARANTEE" {
				p.advance()
				stmt.Retention = "GUARANTEE"
			} else if p.isIdentLike() && p.cur.Str == "NOGUARANTEE" {
				p.advance()
				stmt.Retention = "NOGUARANTEE"
			}

		case p.isIdentLike() && p.cur.Str == "ENCRYPTION":
			p.advance()
			stmt.Encryption = "ENCRYPTION"
			// Skip encryption sub-clause details
			for p.cur.Type != ';' && p.cur.Type != tokEOF && !p.isTablespaceClauseStart() {
				p.advance()
			}

		case p.cur.Type == kwDEFAULT:
			p.advance()
			if p.cur.Type == kwCOMPRESS {
				p.advance()
				stmt.Compress = "COMPRESS"
			} else if p.cur.Type == kwNOCOMPRESS {
				p.advance()
				stmt.Compress = "NOCOMPRESS"
			} else {
				// Skip other DEFAULT sub-clauses
				for p.cur.Type != ';' && p.cur.Type != tokEOF && !p.isTablespaceClauseStart() {
					p.advance()
				}
			}

		case p.isIdentLike() && p.cur.Str == "MAXSIZE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				stmt.MaxSize = "UNLIMITED"
			} else {
				stmt.MaxSize = p.parseSizeValue()
			}

		case p.cur.Type == kwSTORAGE:
			// Skip STORAGE (...)
			p.advance()
			if p.cur.Type == '(' {
				p.skipParens()
			}

		default:
			// Skip unknown tokens
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// isTablespaceClauseStart returns true if the current token starts a known tablespace clause.
func (p *Parser) isTablespaceClauseStart() bool {
	if p.isIdentLike() {
		switch p.cur.Str {
		case "DATAFILE", "TEMPFILE", "AUTOEXTEND", "REUSE", "EXTENT", "SEGMENT",
			"BLOCKSIZE", "RETENTION", "ENCRYPTION", "MAXSIZE", "FORCE":
			return true
		}
	}
	switch p.cur.Type {
	case kwSIZE, kwLOGGING, kwNOLOGGING, kwONLINE, kwOFFLINE, kwDEFAULT, kwSTORAGE:
		return true
	}
	return false
}

// parseDatafileClause parses a single file specification (path string and optional SIZE).
func (p *Parser) parseDatafileClause() *nodes.DatafileClause {
	if p.cur.Type != tokSCONST {
		return nil
	}
	df := &nodes.DatafileClause{
		Loc: nodes.Loc{Start: p.pos()},
	}
	df.Filename = p.cur.Str
	p.advance()

	// Optional SIZE
	if p.cur.Type == kwSIZE {
		p.advance()
		df.Size = p.parseSizeValue()
	}

	// Optional REUSE
	if p.isIdentLike() && p.cur.Str == "REUSE" {
		p.advance()
		df.Reuse = true
	}

	// Optional AUTOEXTEND
	if p.isIdentLike() && p.cur.Str == "AUTOEXTEND" {
		df.Autoextend = p.parseAutoextendClause()
	}

	df.Loc.End = p.pos()
	return df
}

// parseAutoextendClause parses AUTOEXTEND ON/OFF with optional NEXT and MAXSIZE.
func (p *Parser) parseAutoextendClause() *nodes.AutoextendClause {
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
			ac.Next = p.parseSizeValue()
		}
		// Optional MAXSIZE
		if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				ac.MaxSize = "UNLIMITED"
			} else {
				ac.MaxSize = p.parseSizeValue()
			}
		}
	} else if p.isIdentLike() && p.cur.Str == "OFF" {
		p.advance()
		ac.On = false
	}

	ac.Loc.End = p.pos()
	return ac
}

// parseSizeValue parses a size value like "100M", "10G", "512", "8K".
// It combines the number and optional unit suffix into a single string.
func (p *Parser) parseSizeValue() string {
	if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
		val := p.cur.Str
		p.advance()
		// Check for size suffix (K, M, G, T) as an identifier
		if p.isIdentLike() {
			switch p.cur.Str {
			case "K", "M", "G", "T":
				val += p.cur.Str
				p.advance()
			}
		}
		return val
	}
	return ""
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

// parseCreateClusterStmt parses a CREATE CLUSTER statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-CLUSTER.html
//
//	CREATE CLUSTER [ schema. ] cluster
//	    ( { column datatype [ COLLATE collation ] [ SORT ] } [, ...] )
//	    [ physical_attributes_clause ]
//	    [ SIZE size_clause ]
//	    [ TABLESPACE tablespace ]
//	    [ { INDEX
//	      | [ SINGLE TABLE ] HASHKEYS integer [ HASH IS expr ] } ]
//	    [ parallel_clause ]
//	    [ NOROWDEPENDENCIES | ROWDEPENDENCIES ]
//	    [ CACHE | NOCACHE ]
//	    [ cluster_range_partitions ]
func (p *Parser) parseCreateClusterStmt(start int) *nodes.CreateClusterStmt {
	stmt := &nodes.CreateClusterStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Parse column list in parentheses: ( col datatype [SORT] [, ...] )
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col := &nodes.ClusterColumn{
				Loc: nodes.Loc{Start: p.pos()},
			}
			col.Name = p.parseIdentifier()
			col.DataType = p.parseTypeName()
			// Optional COLLATE
			if p.isIdentLike() && p.cur.Str == "COLLATE" {
				p.advance()
				p.parseIdentifier() // collation name
			}
			// Optional SORT
			if p.cur.Type == kwORDER || (p.isIdentLike() && p.cur.Str == "SORT") {
				p.advance()
				col.Sort = true
			}
			col.Loc.End = p.pos()
			stmt.Columns = append(stmt.Columns, col)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Parse options until ; or EOF
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				v := p.parseIntValue()
				stmt.PctFree = &v
			}

		case p.isIdentLike() && p.cur.Str == "PCTUSED":
			p.advance()
			if p.cur.Type == tokICONST {
				v := p.parseIntValue()
				stmt.PctUsed = &v
			}

		case p.isIdentLike() && p.cur.Str == "INITRANS":
			p.advance()
			if p.cur.Type == tokICONST {
				v := p.parseIntValue()
				stmt.InitTrans = &v
			}

		case p.cur.Type == kwSIZE:
			p.advance()
			stmt.Size = p.parseSizeValue()

		case p.cur.Type == kwTABLESPACE:
			p.advance()
			stmt.Tablespace = p.parseIdentifier()

		case p.cur.Type == kwINDEX:
			p.advance()
			stmt.IsIndex = true

		case p.isIdentLike() && p.cur.Str == "SINGLE":
			p.advance()
			if p.cur.Type == kwTABLE {
				p.advance()
			}
			stmt.SingleTable = true

		case p.isIdentLike() && p.cur.Str == "HASHKEYS":
			p.advance()
			stmt.IsHash = true
			if p.cur.Type == tokICONST {
				stmt.HashKeys = p.cur.Str
				p.advance()
			}

		case p.cur.Type == kwHASH:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "IS" {
				p.advance()
				stmt.HashExpr = p.parseExpr()
			}

		case p.cur.Type == kwCACHE:
			p.advance()
			stmt.Cache = true

		case p.cur.Type == kwNOCACHE:
			p.advance()
			stmt.NoCache = true

		case p.isIdentLike() && p.cur.Str == "PARALLEL":
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			} else {
				stmt.Parallel = "PARALLEL"
			}

		case p.isIdentLike() && p.cur.Str == "NOPARALLEL":
			p.advance()
			stmt.Parallel = "NOPARALLEL"

		case p.isIdentLike() && p.cur.Str == "ROWDEPENDENCIES":
			p.advance()
			stmt.RowDep = true

		case p.isIdentLike() && p.cur.Str == "NOROWDEPENDENCIES":
			p.advance()
			stmt.NoRowDep = true

		case p.cur.Type == kwSTORAGE:
			p.advance()
			if p.cur.Type == '(' {
				// Collect storage tokens
				start := p.pos()
				p.skipParens()
				_ = start
			}

		case p.isIdentLike() && p.cur.Str == "SHARING":
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			p.parseIdentifier() // METADATA, DATA, NONE, etc.

		case p.isIdentLike() && p.cur.Str == "PARTITION":
			// cluster_range_partitions - skip until end
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseIntValue parses an integer constant and returns its value.
func (p *Parser) parseIntValue() int {
	if p.cur.Type == tokICONST {
		val := 0
		for _, c := range p.cur.Str {
			val = val*10 + int(c-'0')
		}
		p.advance()
		return val
	}
	return 0
}

// parseCreateDimensionStmt parses a CREATE DIMENSION statement.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-DIMENSION.html
//
//	CREATE DIMENSION [ schema. ] dimension
//	    { LEVEL level IS ( level_table.level_column [, ...] ) [ SKIP WHEN NULL ] } ...
//	    { HIERARCHY hierarchy (
//	        child_level CHILD OF parent_level [ CHILD OF ... ]
//	        [ JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level ] ...
//	      )
//	    | ATTRIBUTE level DETERMINES ( dependent_column [, ...] )
//	    | ATTRIBUTE attr_name LEVEL level DETERMINES ( dependent_column [, ...] )
//	    } ...
func (p *Parser) parseCreateDimensionStmt(start int) *nodes.CreateDimensionStmt {
	stmt := &nodes.CreateDimensionStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Parse clauses: LEVEL, HIERARCHY, ATTRIBUTE
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwLEVEL:
			stmt.Levels = append(stmt.Levels, p.parseDimensionLevel())

		case p.isIdentLike() && p.cur.Str == "HIERARCHY":
			stmt.Hierarchies = append(stmt.Hierarchies, p.parseDimensionHierarchy())

		case p.isIdentLike() && p.cur.Str == "ATTRIBUTE":
			stmt.Attributes = append(stmt.Attributes, p.parseDimensionAttribute())

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDimensionLevel parses a LEVEL clause.
//
//	LEVEL level IS ( level_table.level_column [, ...] ) [ SKIP WHEN NULL ]
//	-- or without parens:
//	LEVEL level IS level_table.level_column [ SKIP WHEN NULL ]
func (p *Parser) parseDimensionLevel() *nodes.DimensionLevel {
	lvl := &nodes.DimensionLevel{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume LEVEL

	lvl.Name = p.parseIdentifier()

	// IS
	if p.isIdentLike() && p.cur.Str == "IS" {
		p.advance()
	}

	// Column list: ( col [, ...] ) or just col
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			lvl.Columns = append(lvl.Columns, p.parseObjectName())
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	} else {
		// Single column reference: table.column
		lvl.Columns = append(lvl.Columns, p.parseObjectName())
	}

	// Optional SKIP WHEN NULL
	if p.cur.Type == kwSKIP {
		p.advance() // SKIP
		if p.cur.Type == kwWHEN {
			p.advance() // WHEN
		}
		if p.cur.Type == kwNULL {
			p.advance() // NULL
		}
		lvl.SkipWhenNull = true
	}

	lvl.Loc.End = p.pos()
	return lvl
}

// parseDimensionHierarchy parses a HIERARCHY clause.
//
//	HIERARCHY hierarchy_name (
//	    child_level CHILD OF parent_level [ CHILD OF ... ]
//	    [ JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level ] ...
//	)
func (p *Parser) parseDimensionHierarchy() *nodes.DimensionHierarchy {
	hier := &nodes.DimensionHierarchy{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume HIERARCHY

	hier.Name = p.parseIdentifier()

	// Parse ( ... )
	if p.cur.Type == '(' {
		p.advance()

		// Parse level chain: child CHILD OF parent CHILD OF ...
		// First level
		if p.isIdentLike() {
			hier.Levels = append(hier.Levels, p.parseIdentifier())
		}

		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			// Check for CHILD OF
			if p.isIdentLike() && p.cur.Str == "CHILD" {
				p.advance() // CHILD
				if p.cur.Type == kwOF {
					p.advance() // OF
				}
				// Next level
				if p.isIdentLike() {
					hier.Levels = append(hier.Levels, p.parseIdentifier())
				}
			} else if p.isIdentLike() && p.cur.Str == "JOIN" {
				// JOIN KEY clause
				jk := p.parseDimensionJoinKey()
				hier.JoinKeys = append(hier.JoinKeys, jk)
			} else {
				break
			}
		}

		if p.cur.Type == ')' {
			p.advance()
		}
	}

	hier.Loc.End = p.pos()
	return hier
}

// parseDimensionJoinKey parses a JOIN KEY clause.
//
//	JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level
func (p *Parser) parseDimensionJoinKey() *nodes.DimensionJoinKey {
	jk := &nodes.DimensionJoinKey{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume JOIN

	if p.isIdentLike() && p.cur.Str == "KEY" {
		p.advance() // KEY
	}

	// ( child_key_column [, ...] )
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			jk.ChildKeys = append(jk.ChildKeys, p.parseObjectName())
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// REFERENCES parent_level
	if p.cur.Type == kwREFERENCES {
		p.advance()
		jk.ParentLevel = p.parseIdentifier()
	}

	jk.Loc.End = p.pos()
	return jk
}

// parseDimensionAttribute parses an ATTRIBUTE clause.
//
//	ATTRIBUTE level DETERMINES ( dependent_column [, ...] )
//	-- or extended form:
//	ATTRIBUTE attr_name LEVEL level DETERMINES ( dependent_column [, ...] )
func (p *Parser) parseDimensionAttribute() *nodes.DimensionAttribute {
	attr := &nodes.DimensionAttribute{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance() // consume ATTRIBUTE

	attr.AttrName = p.parseIdentifier()

	// Check for extended form: LEVEL level DETERMINES ...
	if p.cur.Type == kwLEVEL {
		p.advance()
		attr.LevelName = p.parseIdentifier()
	}

	// DETERMINES
	if p.isIdentLike() && p.cur.Str == "DETERMINES" {
		p.advance()
	}

	// ( dependent_column [, ...] ) or single column
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			attr.Columns = append(attr.Columns, p.parseObjectName())
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	} else {
		attr.Columns = append(attr.Columns, p.parseObjectName())
	}

	attr.Loc.End = p.pos()
	return attr
}

// parseAlterAdminObject handles ALTER dispatches for admin DDL objects.
func (p *Parser) parseAlterAdminObject(start int) nodes.StmtNode {
	switch p.cur.Type {
	case kwUSER:
		p.advance()
		return p.parseAlterUserStmt(start)
	case kwROLE:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_ROLE, start)
	case kwPROFILE:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_PROFILE, start)
	case kwTABLESPACE:
		p.advance()
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_TABLESPACE_SET, start)
		}
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_TABLESPACE, start)
	case kwCLUSTER:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_CLUSTER, start)
	case kwJAVA:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_JAVA, start)
	case kwLIBRARY:
		p.advance()
		return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_LIBRARY, start)
	default:
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DIMENSION":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_DIMENSION, start)
			case "DISKGROUP":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_DISKGROUP, start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_PLUGGABLE_DATABASE, start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type == kwVIEW {
					p.advance() // consume VIEW
				}
				return p.parseAlterAnalyticViewStmt(start)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if p.isIdentLike() && p.cur.Str == "DIMENSION" {
					p.advance() // consume DIMENSION
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_ATTRIBUTE_DIMENSION, start)
			case "HIERARCHY":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_HIERARCHY, start)
			case "DOMAIN":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_DOMAIN, start)
			case "INDEXTYPE":
				p.advance()
				return p.parseAlterIndextypeStmt(start)
			case "OPERATOR":
				p.advance()
				return p.parseAlterOperatorStmt(start)
			case "LOCKDOWN":
				p.advance() // consume LOCKDOWN
				if p.cur.Type == kwPROFILE {
					p.advance() // consume PROFILE
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_LOCKDOWN_PROFILE, start)
			case "OUTLINE":
				p.advance()
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_OUTLINE, start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if p.cur.Type == kwJOIN {
					p.advance() // consume JOIN
				}
				if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
					p.advance() // consume GROUP
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_INMEMORY_JOIN_GROUP, start)
			case "FLASHBACK":
				p.advance() // consume FLASHBACK
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance() // consume ARCHIVE
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			case "RESOURCE":
				p.advance() // consume RESOURCE
				if p.isIdentLike() && p.cur.Str == "COST" {
					p.advance() // consume COST
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_RESOURCE_COST, start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
				}
				return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_ROLLBACK_SEGMENT, start)
			case "MATERIALIZED":
				p.advance() // consume MATERIALIZED (already a keyword, handled above for MVIEW)
				if p.isIdentLike() && p.cur.Str == "ZONEMAP" {
					p.advance() // consume ZONEMAP
					return p.parseAdminDDLStmt("ALTER", nodes.OBJECT_MATERIALIZED_ZONEMAP, start)
				}
				return nil
			}
		}
		return nil
	}
}
