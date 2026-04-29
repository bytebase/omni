package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateMaterializedOrView distinguishes between:
// - CREATE MATERIALIZED VIEW LOG ON ... (mview log)
// - CREATE MATERIALIZED VIEW ... (regular mview)
func (p *Parser) parseCreateMaterializedOrView(start int, orReplace bool) (nodes.StmtNode, error) {
	p.advance() // consume MATERIALIZED
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
func (p *Parser) parseAdminDDLStmt(action string, objType nodes.ObjectType, start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     action,
		ObjectType: objType,
		Loc:        nodes.Loc{Start: start},
	}
	var parseErr197 error

	stmt.Name, parseErr197 = p.parseObjectName()
	if parseErr197 !=

		// Skip remaining tokens until ;/EOF
		nil {
		return nil, parseErr197
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
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

				p.skipToSemicolon()
				goto done
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
		p.advance()
		// Check for TABLESPACE SET
		if p.cur.Type == kwSET {
			p.advance() // consume SET
			return p.parseCreateTablespaceSetStmt(start)
		}
		return p.parseCreateTablespaceStmt(start, false, false, false, false, false)
	case kwDIRECTORY:
		p.advance()
		return p.parseCreateDirectoryStmt(start, orReplace)
	case kwCONTEXT:
		p.advance()
		return p.parseCreateContextStmt(start, orReplace)
	case kwCLUSTER:
		p.advance()
		return p.parseCreateClusterStmt(start)
	case kwJAVA:
		p.advance()
		return p.parseCreateJavaStmt(start, orReplace)
	case kwLIBRARY:
		p.advance()
		return p.parseCreateLibraryStmt(start, orReplace)
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
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseCreatePluggableDatabaseStmt(start)
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
				if p.cur.Type == kwPROFILE {
					p.advance() // consume PROFILE
				}
				return p.parseCreateLockdownProfileStmt(start)
			case "OUTLINE":
				p.advance()
				return p.parseCreateOutlineStmt(start, false, false)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if p.cur.Type == kwJOIN {
					p.advance() // consume JOIN
				}
				if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
					p.advance() // consume GROUP
				}
				return p.parseCreateInmemoryJoinGroupStmt(start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
				}
				return p.parseCreateRollbackSegmentStmt(start, false)
			case "EDITION":
				p.advance()
				return p.parseCreateEditionStmt(start)
			case "MLE":
				p.advance() // consume MLE
				if p.isIdentLike() && p.cur.Str == "ENV" {
					p.advance() // consume ENV
					return p.parseCreateMLEEnvStmt(start, orReplace)
				}
				if p.isIdentLike() && p.cur.Str == "MODULE" {
					p.advance() // consume MODULE
					return p.parseCreateMLEModuleStmt(start, orReplace)
				}
				return p.parseCreateMLEEnvStmt(start, orReplace)
			case "PFILE":
				p.advance()
				return p.parseCreatePfileStmt(start)
			case "SPFILE":
				p.advance()
				return p.parseCreateSpfileStmt(start)
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if p.isIdentLike() && p.cur.Str == "GRAPH" {
					p.advance() // consume GRAPH
				}
				return p.parseCreatePropertyGraphStmt(start, false, false)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type == kwINDEX {
					p.advance() // consume INDEX
				}
				return p.parseCreateVectorIndexStmt(start, false)
			case "RESTORE":
				p.advance() // consume RESTORE
				if p.isIdentLike() && p.cur.Str == "POINT" {
					p.advance() // consume POINT
				}
				return p.parseCreateRestorePointStmt(start, false)
			case "CLEAN":
				p.advance() // consume CLEAN
				if p.isIdentLike() && p.cur.Str == "RESTORE" {
					p.advance() // consume RESTORE
				}
				if p.isIdentLike() && p.cur.Str == "POINT" {
					p.advance() // consume POINT
				}
				return p.parseCreateRestorePointStmt(start, true)
			case "LOGICAL":
				p.advance() // consume LOGICAL
				if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
					p.advance() // consume PARTITION
				}
				if p.isIdentLike() && p.cur.Str == "TRACKING" {
					p.advance() // consume TRACKING
				}
				return p.parseCreateLogicalPartitionTrackingStmt(start)
			case "PMEM":
				p.advance() // consume PMEM
				if p.isIdentLike() && p.cur.Str == "FILESTORE" {
					p.advance() // consume FILESTORE
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
		return p.parseAdminDDLStmt("DROP", nodes.OBJECT_ROLE, start)
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
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance()
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_FLASHBACK_ARCHIVE, start)
			case "DISKGROUP":
				p.advance()
				return p.parseDropDiskgroupStmt(start)
			case "PLUGGABLE":
				p.advance() // consume PLUGGABLE
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseDropPluggableDatabaseStmt(start)
			case "ANALYTIC":
				p.advance() // consume ANALYTIC
				if p.cur.Type == kwVIEW {
					p.advance() // consume VIEW
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_ANALYTIC_VIEW, start)
			case "ATTRIBUTE":
				p.advance() // consume ATTRIBUTE
				if p.isIdentLike() && p.cur.Str == "DIMENSION" {
					p.advance() // consume DIMENSION
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
				if p.cur.Type == kwPROFILE {
					p.advance() // consume PROFILE
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_LOCKDOWN_PROFILE, start)
			case "OUTLINE":
				p.advance()
				return p.parseDropSimpleStmt(nodes.OBJECT_OUTLINE, start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if p.cur.Type == kwJOIN {
					p.advance() // consume JOIN
				}
				if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
					p.advance() // consume GROUP
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_INMEMORY_JOIN_GROUP, start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
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
				return p.parseDropSimpleStmt(nodes.OBJECT_MLE_ENV, start)
			case "PROPERTY":
				p.advance() // consume PROPERTY
				if p.isIdentLike() && p.cur.Str == "GRAPH" {
					p.advance() // consume GRAPH
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_PROPERTY_GRAPH, start)
			case "VECTOR":
				p.advance() // consume VECTOR
				if p.cur.Type == kwINDEX {
					p.advance() // consume INDEX
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_VECTOR_INDEX, start)
			case "RESTORE":
				p.advance() // consume RESTORE
				if p.isIdentLike() && p.cur.Str == "POINT" {
					p.advance() // consume POINT
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_RESTORE_POINT, start)
			case "LOGICAL":
				p.advance() // consume LOGICAL
				if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
					p.advance() // consume PARTITION
				}
				if p.isIdentLike() && p.cur.Str == "TRACKING" {
					p.advance() // consume TRACKING
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_LOGICAL_PARTITION_TRACKING, start)
			case "PMEM":
				p.advance() // consume PMEM
				if p.isIdentLike() && p.cur.Str == "FILESTORE" {
					p.advance() // consume FILESTORE
				}
				return p.parseDropSimpleStmt(nodes.OBJECT_PMEM_FILESTORE, start)
			}
		}
		return nil, nil
	}
}

// parseCreateTablespaceStmt parses a CREATE TABLESPACE statement.
//
// BNF: oracle/parser/bnf/CREATE-TABLESPACE.bnf
//
//	CREATE [ BIGFILE | SMALLFILE ]
//	    { permanent_tablespace_clause
//	    | temporary_tablespace_clause
//	    | undo_tablespace_clause
//	    } ;
//
//	permanent_tablespace_clause:
//	    TABLESPACE tablespace [ IF NOT EXISTS ]
//	    [ DATAFILE file_specification [, file_specification ] ... ]
//	    [ permanent_tablespace_attrs ]
//
//	permanent_tablespace_attrs:
//	    [ MINIMUM EXTENT size_clause ]
//	    [ BLOCKSIZE integer [ K ] ]
//	    [ logging_clause ]
//	    [ FORCE LOGGING ]
//	    [ tablespace_encryption_clause ]
//	    [ DEFAULT [ default_tablespace_params ] ]
//	    [ { ONLINE | OFFLINE } ]
//	    [ extent_management_clause ]
//	    [ segment_management_clause ]
//	    [ flashback_mode_clause ]
//	    [ lost_write_protection ]
//
//	temporary_tablespace_clause:
//	    [ LOCAL ] TEMPORARY TABLESPACE tablespace [ IF NOT EXISTS ]
//	    [ TEMPFILE file_specification [, file_specification ] ... ]
//	    [ tablespace_group_clause ]
//	    [ extent_management_clause ]
//	    [ tablespace_encryption_clause ]
//	    [ FOR { ALL | LEAF } ]
//
//	undo_tablespace_clause:
//	    UNDO TABLESPACE tablespace [ IF NOT EXISTS ]
//	    [ DATAFILE file_specification [, file_specification ] ... ]
//	    [ extent_management_clause ]
//	    [ tablespace_retention_clause ]
//	    [ tablespace_encryption_clause ]
//
//	tablespace_encryption_clause:
//	    ENCRYPTION [ tablespace_encryption_spec ] { ENCRYPT | DECRYPT }
//
//	tablespace_encryption_spec:
//	    USING 'encrypt_algorithm'
//
//	default_tablespace_params:
//	    { default_table_compression | default_index_compression
//	    | inmemory_clause | ilm_clause | storage_clause } [ ... ]
//
//	logging_clause:
//	    { LOGGING | NOLOGGING | FILESYSTEM_LIKE_LOGGING }
//
//	extent_management_clause:
//	    EXTENT MANAGEMENT { LOCAL [ { AUTOALLOCATE | UNIFORM [ SIZE size_clause ] } ] | DICTIONARY }
//
//	segment_management_clause:
//	    SEGMENT SPACE MANAGEMENT { AUTO | MANUAL }
//
//	flashback_mode_clause:
//	    FLASHBACK { ON | OFF }
//
//	tablespace_retention_clause:
//	    RETENTION { GUARANTEE | NOGUARANTEE }
//
//	tablespace_group_clause:
//	    TABLESPACE GROUP { tablespace_group_name | '' }
//
//	lost_write_protection:
//	    { ENABLE | DISABLE | SUSPEND | REMOVE } LOST WRITE PROTECTION
//
//	file_specification:
//	    [ 'filename' | 'ASM_filename' ]
//	    [ SIZE size_clause ]
//	    [ REUSE ]
//	    [ autoextend_clause ]
//
//	autoextend_clause:
//	    { AUTOEXTEND OFF | AUTOEXTEND ON [ NEXT size_clause ] [ MAXSIZE { UNLIMITED | size_clause } ] }
//
//	size_clause:
//	    integer [ K | M | G | T | P | E ]
func (p *Parser) parseCreateTablespaceStmt(start int, bigfile, smallfile, local, temporary, undo bool) (*nodes.CreateTablespaceStmt, error) {
	stmt := &nodes.CreateTablespaceStmt{
		Loc:       nodes.Loc{Start: start},
		Bigfile:   bigfile,
		Smallfile: smallfile,
		Local:     local,
		Temporary: temporary,
		Undo:      undo,
	}
	var parseErr200 error

	// Parse tablespace name
	stmt.Name, parseErr200 = p.parseObjectName()
	if parseErr200 !=

		// Optional IF NOT EXISTS
		nil {
		return nil, parseErr200
	}

	if p.cur.Type == kwIF {
		p.advance() // IF
		if p.cur.Type == kwNOT {
			p.advance() // NOT
		}
		if p.cur.Type == kwEXISTS {
			p.advance() // EXISTS
		}
		stmt.IfNotExists = true
	}

	// Parse clauses in any order until ; or EOF
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.isIdentLike() && (p.cur.Str == "DATAFILE" || p.cur.Str == "TEMPFILE"):
			p.advance()
			// Parse one or more file specifications
			for {
				df, parseErr201 := p.parseDatafileClause()
				if parseErr201 != nil {
					return nil, parseErr201
				}
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
			var parseErr202 error
			stmt.Size, parseErr202 = p.parseSizeValue()
			if parseErr202 != nil {
				return nil, parseErr202
			}

		case p.isIdentLike() && p.cur.Str == "AUTOEXTEND":
			var parseErr203 error
			stmt.Autoextend, parseErr203 = p.parseAutoextendClause()
			if parseErr203 != nil {
				return nil, parseErr203
			}

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

		case p.isIdentLike() && p.cur.Str == "FILESYSTEM_LIKE_LOGGING":
			p.advance()
			stmt.Logging = "FILESYSTEM_LIKE_LOGGING"

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

		case p.isIdentLike() && p.cur.Str == "MINIMUM":
			p.advance() // MINIMUM
			if p.isIdentLike() && p.cur.Str == "EXTENT" {
				p.advance() // EXTENT
			}
			var parseErr204 error
			stmt.MinimumExtent, parseErr204 = p.parseSizeValue()
			if parseErr204 != nil {
				return nil, parseErr204
			}

		case p.isIdentLike() && p.cur.Str == "EXTENT":
			var parseErr205 error
			stmt.Extent, parseErr205 = p.parseExtentManagementClause()
			if parseErr205 != nil {
				return nil, parseErr205
			}

		case p.isIdentLike() && p.cur.Str == "SEGMENT":
			var parseErr206 error
			stmt.Segment, parseErr206 = p.parseSegmentManagementClause()
			if parseErr206 != nil {
				return nil, parseErr206
			}

		case p.isIdentLike() && p.cur.Str == "BLOCKSIZE":
			p.advance()
			var parseErr207 error
			stmt.Blocksize, parseErr207 = p.parseSizeValue()
			if parseErr207 != nil {
				return nil, parseErr207
			}

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
			var parseErr208 error
			stmt.Encryption, stmt.EncryptionAlgorithm, parseErr208 = p.parseTablespaceEncryptionClause()
			if parseErr208 != nil {
				return nil, parseErr208
			}

		case p.cur.Type == kwDEFAULT:
			p.advance()
			var parseErr209 error
			stmt.DefaultParams, parseErr209 = p.parseDefaultTablespaceParams()
			if parseErr209 != nil {
				return nil, parseErr209
			}

		case p.isIdentLike() && p.cur.Str == "MAXSIZE":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
				p.advance()
				stmt.MaxSize = "UNLIMITED"
			} else {
				var parseErr210 error
				stmt.MaxSize, parseErr210 = p.parseSizeValue()
				if parseErr210 != nil {
					return nil, parseErr210
				}
			}

		case p.cur.Type == kwSTORAGE:
			p.advance()
			if p.cur.Type == '(' {
				p.skipParens()
			}

		case p.isIdentLike() && p.cur.Str == "FLASHBACK":
			p.advance()
			if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
				p.advance()
				stmt.Flashback = "ON"
			} else if p.isIdentLike() && p.cur.Str == "OFF" {
				p.advance()
				stmt.Flashback = "OFF"
			}

		case p.isIdentLike() && (p.cur.Str == "ENABLE" || p.cur.Str == "DISABLE" || p.cur.Str == "SUSPEND" || p.cur.Str == "REMOVE"):
			var parseErr211 error
			stmt.LostWriteProtection, parseErr211 = p.parseLostWriteProtection()
			if parseErr211 != nil {
				return nil, parseErr211
			}

		case p.cur.Type == kwTABLESPACE:
			// TABLESPACE GROUP clause (in temporary tablespace context)
			p.advance() // TABLESPACE
			if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
				p.advance() // GROUP
				if p.cur.Type == tokSCONST {
					stmt.TablespaceGroup = p.cur.Str
					p.advance()
				} else if p.isIdentLike() || p.cur.Type == tokIDENT {
					stmt.TablespaceGroup = p.cur.Str
					p.advance()
				}
			}

		case p.cur.Type == kwFOR:
			// FOR { ALL | LEAF } in temporary tablespace
			p.advance()
			if p.cur.Type == kwALL || (p.isIdentLike() && p.cur.Str == "ALL") {
				p.advance()
				stmt.ForLeaf = "ALL"
			} else if p.isIdentLike() && p.cur.Str == "LEAF" {
				p.advance()
				stmt.ForLeaf = "LEAF"
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseExtentManagementClause parses EXTENT MANAGEMENT { LOCAL [...] | DICTIONARY }.
func (p *Parser) parseExtentManagementClause() (string, error) {
	p.advance() // EXTENT
	if p.isIdentLike() && p.cur.Str == "MANAGEMENT" {
		p.advance() // MANAGEMENT
	}
	if p.isIdentLike() && p.cur.Str == "DICTIONARY" {
		p.advance()
		return "DICTIONARY", nil
	}
	if p.isIdentLike() && p.cur.Str == "LOCAL" {
		p.advance() // LOCAL
	}
	if p.isIdentLike() && p.cur.Str == "AUTOALLOCATE" {
		p.advance()
		return "AUTOALLOCATE", nil
	} else if p.isIdentLike() && p.cur.Str == "UNIFORM" {
		p.advance()
		if p.cur.Type == kwSIZE {
			p.advance()
			parseValue11, parseErr12 := p.parseSizeValue()
			if parseErr12 != nil {
				return "", parseErr12
			}
			return "UNIFORM SIZE " + parseValue11, nil
		}
		return "UNIFORM", nil
	}
	return "LOCAL", nil
}

// parseSegmentManagementClause parses SEGMENT SPACE MANAGEMENT { AUTO | MANUAL }.
func (p *Parser) parseSegmentManagementClause() (string, error) {
	p.advance() // SEGMENT
	if p.isIdentLike() && p.cur.Str == "SPACE" {
		p.advance() // SPACE
	}
	if p.isIdentLike() && p.cur.Str == "MANAGEMENT" {
		p.advance() // MANAGEMENT
	}
	if p.isIdentLike() && p.cur.Str == "AUTO" {
		p.advance()
		return "AUTO", nil
	} else if p.isIdentLike() && p.cur.Str == "MANUAL" {
		p.advance()
		return "MANUAL", nil
	}
	return "AUTO", nil
}

// parseTablespaceEncryptionClause parses ENCRYPTION [ USING 'algo' ] { ENCRYPT | DECRYPT }.
// Returns (encryption_summary, algorithm).
func (p *Parser) parseTablespaceEncryptionClause() (string, string, error) {
	p.advance() // ENCRYPTION
	algo := ""
	if p.isIdentLike() && p.cur.Str == "USING" {
		p.advance()
		if p.cur.Type == tokSCONST {
			algo = p.cur.Str
			p.advance()
		}
	}
	if p.isIdentLike() && p.cur.Str == "ENCRYPT" {
		p.advance()
		return "ENCRYPT", algo, nil
	} else if p.isIdentLike() && p.cur.Str == "DECRYPT" {
		p.advance()
		return "DECRYPT", algo, nil
	}
	// For ALTER TABLESPACE: ENCRYPTION ONLINE/OFFLINE/FINISH
	return "ENCRYPTION", algo, nil
}

// parseDefaultTablespaceParams parses DEFAULT tablespace params.
// Returns a summary string of the parsed params.
func (p *Parser) parseDefaultTablespaceParams() (string, error) {
	result := ""
	for p.cur.Type != ';' && p.cur.Type != tokEOF && !p.isTablespaceClauseStart() {
		switch {
		case p.cur.Type == kwTABLE:
			p.advance() // TABLE
			if p.cur.Type == kwCOMPRESS {
				p.advance()
				result = p.appendParam(result, "TABLE COMPRESS")
				// FOR OLTP | QUERY | ARCHIVE
				if p.cur.Type == kwFOR {
					p.advance()
					if p.isIdentLike() {
						result = p.appendParam(result, p.cur.Str)
						p.advance()
						// LOW | HIGH
						if p.isIdentLike() && (p.cur.Str == "LOW" || p.cur.Str == "HIGH") {
							result = p.appendParam(result, p.cur.Str)
							p.advance()
						}
					}
				}
			} else if p.cur.Type == kwNOCOMPRESS {
				p.advance()
				result = p.appendParam(result, "TABLE NOCOMPRESS")
			} else if p.isIdentLike() && p.cur.Str == "ROW" {
				p.advance() // ROW
				if p.cur.Type == kwSTORAGE || (p.isIdentLike() && p.cur.Str == "STORE") {
					p.advance() // STORE
				}
				if p.cur.Type == kwCOMPRESS {
					p.advance()
					result = p.appendParam(result, "ROW STORE COMPRESS")
					if p.isIdentLike() && (p.cur.Str == "BASIC" || p.cur.Str == "ADVANCED") {
						result = p.appendParam(result, p.cur.Str)
						p.advance()
					}
				}
			} else if p.isIdentLike() && p.cur.Str == "COLUMN" {
				p.advance() // COLUMN
				if p.isIdentLike() && p.cur.Str == "STORE" {
					p.advance() // STORE
				}
				if p.cur.Type == kwCOMPRESS {
					p.advance()
					result = p.appendParam(result, "COLUMN STORE COMPRESS")
					if p.cur.Type == kwFOR {
						p.advance()
						if p.isIdentLike() {
							result = p.appendParam(result, p.cur.Str)
							p.advance()
							if p.isIdentLike() && (p.cur.Str == "LOW" || p.cur.Str == "HIGH") {
								result = p.appendParam(result, p.cur.Str)
								p.advance()
							}
						}
					}
				}
			}

		case p.cur.Type == kwINDEX:
			p.advance() // INDEX
			if p.cur.Type == kwCOMPRESS {
				p.advance()
				result = p.appendParam(result, "INDEX COMPRESS")
				if p.isIdentLike() && p.cur.Str == "ADVANCED" {
					p.advance()
					result = p.appendParam(result, "ADVANCED")
					if p.isIdentLike() && (p.cur.Str == "LOW" || p.cur.Str == "HIGH") {
						result = p.appendParam(result, p.cur.Str)
						p.advance()
					}
				}
			} else if p.cur.Type == kwNOCOMPRESS {
				p.advance()
				result = p.appendParam(result, "INDEX NOCOMPRESS")
			}

		case p.cur.Type == kwCOMPRESS:
			p.advance()
			result = p.appendParam(result, "COMPRESS")

		case p.cur.Type == kwNOCOMPRESS:
			p.advance()
			result = p.appendParam(result, "NOCOMPRESS")

		case p.isIdentLike() && p.cur.Str == "INMEMORY":
			p.advance()
			result = p.appendParam(result, "INMEMORY")
			// Skip inmemory attributes
			for p.isIdentLike() && (p.cur.Str == "MEMCOMPRESS" || p.cur.Str == "PRIORITY" || p.cur.Str == "DISTRIBUTE" || p.cur.Str == "DUPLICATE") {
				result = p.appendParam(result, p.cur.Str)
				p.advance()
				for p.isIdentLike() && (p.cur.Str == "FOR" || p.cur.Str == "DML" || p.cur.Str == "QUERY" || p.cur.Str == "CAPACITY" || p.cur.Str == "LOW" || p.cur.Str == "HIGH" || p.cur.Str == "NONE" || p.cur.Str == "MEDIUM" || p.cur.Str == "CRITICAL" || p.cur.Str == "AUTO" || p.cur.Str == "BY" || p.cur.Str == "ROWID" || p.cur.Str == "RANGE" || p.cur.Str == "ALL") {
					p.advance()
				}
				if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "SUBPARTITION") {
					p.advance()
				}
			}

		case p.isIdentLike() && p.cur.Str == "NO":
			p.advance()
			if p.isIdentLike() && p.cur.Str == "INMEMORY" {
				p.advance()
				result = p.appendParam(result, "NO INMEMORY")
			} else if p.isIdentLike() && p.cur.Str == "MEMCOMPRESS" {
				p.advance()
				result = p.appendParam(result, "NO MEMCOMPRESS")
			} else if p.isIdentLike() && p.cur.Str == "DUPLICATE" {
				p.advance()
				result = p.appendParam(result, "NO DUPLICATE")
			}

		case p.cur.Type == kwSTORAGE:
			p.advance()
			result = p.appendParam(result, "STORAGE")
			if p.cur.Type == '(' {
				p.skipParens()
			}

		default:
			p.advance()
		}
	}
	if result == "" {
		result = "DEFAULT"
	}
	return result, nil
}

// appendParam appends a parameter to a space-separated string.
func (p *Parser) appendParam(base, param string) string {
	if base == "" {
		return param
	}
	return base + " " + param
}

// parseLostWriteProtection parses { ENABLE | DISABLE | SUSPEND | REMOVE } LOST WRITE PROTECTION.
func (p *Parser) parseLostWriteProtection() (string, error) {
	action := p.cur.Str
	p.advance() // ENABLE/DISABLE/SUSPEND/REMOVE
	if p.isIdentLike() && p.cur.Str == "LOST" {
		p.advance() // LOST
	}
	if p.isIdentLike() && p.cur.Str == "WRITE" {
		p.advance() // WRITE
	}
	if p.isIdentLike() && p.cur.Str == "PROTECTION" {
		p.advance() // PROTECTION
	}
	return action, nil
}

// parseAlterTablespaceStmt parses an ALTER TABLESPACE statement.
//
// BNF: oracle/parser/bnf/ALTER-TABLESPACE.bnf
//
//	ALTER TABLESPACE [ IF EXISTS ] tablespace
//	    alter_tablespace_attrs
//
//	alter_tablespace_attrs:
//	    { default_tablespace_params
//	    | MINIMUM EXTENT size_clause
//	    | RESIZE size_clause
//	    | COALESCE
//	    | SHRINK SPACE [ KEEP size_clause ]
//	    | RENAME TO new_tablespace_name
//	    | BEGIN BACKUP
//	    | END BACKUP
//	    | datafile_tempfile_clauses
//	    | tablespace_logging_clauses
//	    | tablespace_group_clause
//	    | tablespace_state_clauses
//	    | autoextend_clause
//	    | flashback_mode_clause
//	    | tablespace_retention_clause
//	    | alter_tablespace_encryption
//	    | lost_write_protection
//	    }
//
//	datafile_tempfile_clauses:
//	    { ADD { DATAFILE | TEMPFILE } [ file_specification [, file_specification ]... ]
//	    | DROP { DATAFILE | TEMPFILE } { 'filename' | file_number }
//	    | SHRINK TEMPFILE { 'filename' | file_number } [ KEEP size_clause ]
//	    | RENAME DATAFILE 'filename' [, 'filename' ]... TO 'filename' [, 'filename' ]...
//	    | { DATAFILE | TEMPFILE } { ONLINE | OFFLINE }
//	    }
//
//	tablespace_logging_clauses:
//	    { logging_clause | [ NO ] FORCE LOGGING }
//
//	tablespace_state_clauses:
//	    { ONLINE | OFFLINE [ NORMAL | TEMPORARY | IMMEDIATE ]
//	    | READ ONLY | READ WRITE | PERMANENT | TEMPORARY }
//
//	alter_tablespace_encryption:
//	    ENCRYPTION
//	        { ONLINE [ tablespace_encryption_spec ] [ ts_file_name_convert ]
//	        | OFFLINE { ENCRYPT [ tablespace_encryption_spec ] | DECRYPT }
//	        | FINISH [ ENCRYPT | DECRYPT ] [ ts_file_name_convert ]
//	        }
//
//	lost_write_protection:
//	    { ENABLE | DISABLE | REMOVE | SUSPEND } LOST WRITE PROTECTION
func (p *Parser) parseAlterTablespaceStmt(start int, isSet bool) (*nodes.AlterTablespaceStmt, error) {
	stmt := &nodes.AlterTablespaceStmt{
		IsSet: isSet,
		Loc:   nodes.Loc{Start: start},
	}

	// Optional IF EXISTS (not for SET)
	if !isSet && p.cur.Type == kwIF {
		p.advance() // IF
		if p.cur.Type == kwEXISTS {
			p.advance() // EXISTS
		}
		stmt.IfExists = true
	}
	var parseErr212 error

	// Parse tablespace name
	stmt.Name, parseErr212 = p.parseObjectName()
	if parseErr212 !=

		// Parse alter clauses
		nil {
		return nil, parseErr212
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwDEFAULT:
			p.advance()
			var parseErr213 error
			stmt.DefaultParams, parseErr213 = p.parseDefaultTablespaceParams()
			if parseErr213 != nil {
				return nil, parseErr213
			}

		case p.isIdentLike() && p.cur.Str == "MINIMUM":
			p.advance() // MINIMUM
			if p.isIdentLike() && p.cur.Str == "EXTENT" {
				p.advance() // EXTENT
			}
			var parseErr214 error
			stmt.MinimumExtent, parseErr214 = p.parseSizeValue()
			if parseErr214 != nil {
				return nil, parseErr214
			}

		case p.isIdentLike() && p.cur.Str == "RESIZE":
			p.advance()
			var parseErr215 error
			stmt.Resize, parseErr215 = p.parseSizeValue()
			if parseErr215 != nil {
				return nil, parseErr215
			}

		case p.isIdentLike() && p.cur.Str == "COALESCE":
			p.advance()
			stmt.Coalesce = true

		case p.isIdentLike() && p.cur.Str == "SHRINK":
			p.advance() // SHRINK
			if p.isIdentLike() && p.cur.Str == "TEMPFILE" {
				p.advance() // TEMPFILE
				// { 'filename' | file_number }
				if p.cur.Type == tokSCONST {
					stmt.ShrinkTempfile = p.cur.Str
					p.advance()
				} else if p.cur.Type == tokICONST {
					stmt.ShrinkTempfile = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() && p.cur.Str == "KEEP" {
					p.advance()
					var parseErr216 error
					stmt.ShrinkTempfileKeep, parseErr216 = p.parseSizeValue()
					if parseErr216 != nil {
						return nil,

							// SHRINK SPACE
							parseErr216
					}
				}
			} else {

				if p.isIdentLike() && p.cur.Str == "SPACE" {
					p.advance()
				}
				stmt.ShrinkSpace = true
				if p.isIdentLike() && p.cur.Str == "KEEP" {
					p.advance()
					var parseErr217 error
					stmt.ShrinkKeep, parseErr217 = p.parseSizeValue()
					if parseErr217 != nil {
						return nil, parseErr217
					}
				}
			}

		case p.isIdentLike() && p.cur.Str == "RENAME":
			p.advance() // RENAME
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance() // DATAFILE
				stmt.RenameDatafile = true
				// Parse old filenames
				for {
					if p.cur.Type == tokSCONST {
						stmt.RenameFrom = append(stmt.RenameFrom, p.cur.Str)
						p.advance()
					}
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					break
				}
				// TO
				if p.cur.Type == kwTO {
					p.advance()
				}
				// Parse new filenames
				for {
					if p.cur.Type == tokSCONST {
						stmt.RenameTo2 = append(stmt.RenameTo2, p.cur.Str)
						p.advance()
					}
					if p.cur.Type == ',' {
						p.advance()
						continue
					}
					break
				}
			} else {
				// RENAME TO new_name
				if p.cur.Type == kwTO {
					p.advance()
				}
				var parseErr218 error
				stmt.RenameTo, parseErr218 = p.parseIdentifier()
				if parseErr218 != nil {
					return nil, parseErr218
				}
			}

		case p.isIdentLike() && p.cur.Str == "BEGIN":
			p.advance() // BEGIN
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			stmt.BeginBackup = true

		case p.isIdentLike() && p.cur.Str == "END":
			p.advance() // END
			if p.isIdentLike() && p.cur.Str == "BACKUP" {
				p.advance()
			}
			stmt.EndBackup = true

		case p.cur.Type == kwADD:
			p.advance() // ADD
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance()
				stmt.AddDatafile = true
			} else if p.isIdentLike() && p.cur.Str == "TEMPFILE" {
				p.advance()
				stmt.AddTempfile = true
			}
			// Parse file specifications
			for {
				df, parseErr219 := p.parseDatafileClause()
				if parseErr219 != nil {
					return nil, parseErr219
				}
				if df != nil {
					stmt.Datafiles = append(stmt.Datafiles, df)
				}
				if p.cur.Type == ',' {
					p.advance()
					continue
				}
				break
			}

		case p.cur.Type == kwDROP:
			p.advance() // DROP
			if p.isIdentLike() && p.cur.Str == "DATAFILE" {
				p.advance()
				stmt.DropDatafile = true
			} else if p.isIdentLike() && p.cur.Str == "TEMPFILE" {
				p.advance()
				stmt.DropTempfile = true
			}
			// { 'filename' | file_number }
			if p.cur.Type == tokSCONST {
				stmt.DropFileRef = p.cur.Str
				p.advance()
			} else if p.cur.Type == tokICONST {
				stmt.DropFileRef = p.cur.Str
				p.advance()
			}

		case p.isIdentLike() && p.cur.Str == "DATAFILE":
			p.advance() // DATAFILE
			if p.cur.Type == kwONLINE {
				p.advance()
				stmt.DatafileOnline = true
			} else if p.cur.Type == kwOFFLINE {
				p.advance()
				stmt.DatafileOffline = true
			}

		case p.isIdentLike() && p.cur.Str == "TEMPFILE":
			p.advance() // TEMPFILE
			if p.cur.Type == kwONLINE {
				p.advance()
				stmt.TempfileOnline = true
			} else if p.cur.Type == kwOFFLINE {
				p.advance()
				stmt.TempfileOffline = true
			}

		case p.cur.Type == kwLOGGING:
			p.advance()
			stmt.Logging = "LOGGING"

		case p.cur.Type == kwNOLOGGING:
			p.advance()
			stmt.Logging = "NOLOGGING"

		case p.isIdentLike() && p.cur.Str == "FILESYSTEM_LIKE_LOGGING":
			p.advance()
			stmt.Logging = "FILESYSTEM_LIKE_LOGGING"

		case p.isIdentLike() && p.cur.Str == "FORCE":
			p.advance() // FORCE
			if p.cur.Type == kwLOGGING {
				p.advance()
			}
			stmt.ForceLogging = "FORCE LOGGING"

		case p.isIdentLike() && p.cur.Str == "NO":
			next := p.peekNext()
			if next.Type == kwFORCE || (next.Str == "FORCE" && (next.Type == tokIDENT || next.Type >= 2000)) {
				p.advance() // NO
				p.advance() // FORCE
				if p.cur.Type == kwLOGGING {
					p.advance()
				}
				stmt.ForceLogging = "NO FORCE LOGGING"
			} else {
				p.advance()
			}

		case p.cur.Type == kwONLINE:
			p.advance()
			stmt.Online = true

		case p.cur.Type == kwOFFLINE:
			p.advance()
			stmt.Offline = true
			// Optional NORMAL | TEMPORARY | IMMEDIATE
			if p.isIdentLike() && p.cur.Str == "NORMAL" {
				p.advance()
				stmt.OfflineMode = "NORMAL"
			} else if p.cur.Type == kwTEMPORARY {
				p.advance()
				stmt.OfflineMode = "TEMPORARY"
			} else if p.isIdentLike() && p.cur.Str == "IMMEDIATE" {
				p.advance()
				stmt.OfflineMode = "IMMEDIATE"
			}

		case p.cur.Type == kwREAD:
			p.advance() // READ
			if p.isIdentLike() && p.cur.Str == "ONLY" {
				p.advance()
				stmt.ReadOnly = true
			} else if p.isIdentLike() && p.cur.Str == "WRITE" {
				p.advance()
				stmt.ReadWrite = true
			}

		case p.isIdentLike() && p.cur.Str == "PERMANENT":
			p.advance()
			stmt.Permanent = true

		case p.cur.Type == kwTEMPORARY:
			p.advance()
			stmt.TempState = true

		case p.isIdentLike() && p.cur.Str == "AUTOEXTEND":
			var parseErr220 error
			stmt.Autoextend, parseErr220 = p.parseAutoextendClause()
			if parseErr220 != nil {
				return nil, parseErr220
			}

		case p.isIdentLike() && p.cur.Str == "FLASHBACK":
			p.advance()
			if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
				p.advance()
				stmt.Flashback = "ON"
			} else if p.isIdentLike() && p.cur.Str == "OFF" {
				p.advance()
				stmt.Flashback = "OFF"
			}

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
			p.advance() // ENCRYPTION
			// ALTER has ONLINE/OFFLINE/FINISH sub-clauses
			enc := "ENCRYPTION"
			if p.isIdentLike() && p.cur.Str == "ONLINE" {
				p.advance()
				enc = "ONLINE"
			} else if p.cur.Type == kwOFFLINE {
				p.advance()
				enc = "OFFLINE"
			} else if p.isIdentLike() && p.cur.Str == "FINISH" {
				p.advance()
				enc = "FINISH"
			}
			// Skip remaining encryption sub-clause tokens
			for p.cur.Type != ';' && p.cur.Type != tokEOF && !p.isAlterTablespaceClauseStart() {
				p.advance()
			}
			stmt.Encryption = enc

		case p.isIdentLike() && (p.cur.Str == "ENABLE" || p.cur.Str == "DISABLE" || p.cur.Str == "SUSPEND" || p.cur.Str == "REMOVE"):
			var parseErr221 error
			stmt.LostWriteProtection, parseErr221 = p.parseLostWriteProtection()
			if parseErr221 != nil {
				return nil, parseErr221
			}

		case p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP"):
			p.advance() // GROUP
			if p.cur.Type == tokSCONST {
				stmt.TablespaceGroup = p.cur.Str
				p.advance()
			} else if p.isIdentLike() || p.cur.Type == tokIDENT {
				stmt.TablespaceGroup = p.cur.Str
				p.advance()
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// isAlterTablespaceClauseStart returns true if the current token starts a known ALTER TABLESPACE clause.
func (p *Parser) isAlterTablespaceClauseStart() bool {
	if p.isIdentLike() {
		switch p.cur.Str {
		case "DATAFILE", "TEMPFILE", "AUTOEXTEND", "EXTENT", "SEGMENT",
			"BLOCKSIZE", "RETENTION", "ENCRYPTION", "MAXSIZE", "FORCE",
			"MINIMUM", "RESIZE", "COALESCE", "SHRINK", "RENAME",
			"BEGIN", "END", "FLASHBACK", "ENABLE", "DISABLE",
			"SUSPEND", "REMOVE", "PERMANENT", "NO", "GROUP",
			"FILESYSTEM_LIKE_LOGGING":
			return true
		}
	}
	switch p.cur.Type {
	case kwLOGGING, kwNOLOGGING, kwONLINE, kwOFFLINE, kwDEFAULT, kwREAD, kwADD, kwDROP, kwTEMPORARY, kwGROUP:
		return true
	}
	return false
}

// isTablespaceClauseStart returns true if the current token starts a known tablespace clause.
func (p *Parser) isTablespaceClauseStart() bool {
	if p.isIdentLike() {
		switch p.cur.Str {
		case "DATAFILE", "TEMPFILE", "AUTOEXTEND", "REUSE", "EXTENT", "SEGMENT",
			"BLOCKSIZE", "RETENTION", "ENCRYPTION", "MAXSIZE", "FORCE",
			"MINIMUM", "FLASHBACK", "ENABLE", "DISABLE", "SUSPEND", "REMOVE",
			"FILESYSTEM_LIKE_LOGGING":
			return true
		}
	}
	switch p.cur.Type {
	case kwSIZE, kwLOGGING, kwNOLOGGING, kwONLINE, kwOFFLINE, kwDEFAULT, kwSTORAGE, kwTABLESPACE, kwFOR:
		return true
	}
	return false
}

// parseCreateTablespaceSetStmt parses a CREATE TABLESPACE SET statement.
//
// BNF: oracle/parser/bnf/CREATE-TABLESPACE-SET.bnf
//
//	CREATE TABLESPACE SET tablespace_set
//	    [ IN SHARDSPACE shardspace_name ]
//	    [ USING TEMPLATE
//	      ( [ DATAFILE [ file_specification ] ]
//	        [ permanent_tablespace_attrs ]
//	      )
//	    ] ;
//
//	permanent_tablespace_attrs:
//	    [ logging_clause ]
//	    [ tablespace_encryption_clause ]
//	    [ DEFAULT default_tablespace_params ]
//	    [ extent_management_clause ]
//	    [ segment_management_clause ]
//	    [ flashback_mode_clause ]
func (p *Parser) parseCreateTablespaceSetStmt(start int) (*nodes.CreateTablespaceSetStmt, error) {
	stmt := &nodes.CreateTablespaceSetStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr222 error

	stmt.Name, parseErr222 = p.parseObjectName()
	if parseErr222 !=

		// Optional IN SHARDSPACE
		nil {
		return nil, parseErr222
	}

	if p.cur.Type == kwIN {
		p.advance() // IN
		if p.isIdentLike() && p.cur.Str == "SHARDSPACE" {
			p.advance() // SHARDSPACE
		}
		var parseErr223 error
		stmt.Shardspace, parseErr223 = p.parseIdentifier()
		if parseErr223 !=

			// Optional USING TEMPLATE ( ... )
			nil {
			return nil, parseErr223
		}
	}

	if p.isIdentLike() && p.cur.Str == "USING" {
		p.advance() // USING
		if p.isIdentLike() && p.cur.Str == "TEMPLATE" {
			p.advance() // TEMPLATE
		}
		if p.cur.Type == '(' {
			p.advance() // consume (
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				switch {
				case p.isIdentLike() && p.cur.Str == "DATAFILE":
					p.advance()
					for {
						df, parseErr224 := p.parseDatafileClause()
						if parseErr224 != nil {
							return nil, parseErr224
						}
						if df != nil {
							stmt.Datafiles = append(stmt.Datafiles, df)
						}
						if p.cur.Type == ',' {
							p.advance()
							continue
						}
						break
					}

				case p.cur.Type == kwLOGGING:
					p.advance()
					stmt.Logging = "LOGGING"

				case p.cur.Type == kwNOLOGGING:
					p.advance()
					stmt.Logging = "NOLOGGING"

				case p.isIdentLike() && p.cur.Str == "FILESYSTEM_LIKE_LOGGING":
					p.advance()
					stmt.Logging = "FILESYSTEM_LIKE_LOGGING"

				case p.isIdentLike() && p.cur.Str == "ENCRYPTION":
					enc, encryptionAlgorithm, parseErr225 := p.parseTablespaceEncryptionClause()
					_ = encryptionAlgorithm
					if parseErr225 != nil {
						return nil, parseErr225
					}
					stmt.Encryption = enc

				case p.cur.Type == kwDEFAULT:
					p.advance()
					var parseErr226 error
					stmt.DefaultParams, parseErr226 = p.parseDefaultTablespaceParams()
					if parseErr226 != nil {
						return nil, parseErr226
					}

				case p.isIdentLike() && p.cur.Str == "EXTENT":
					var parseErr227 error
					stmt.Extent, parseErr227 = p.parseExtentManagementClause()
					if parseErr227 != nil {
						return nil, parseErr227
					}

				case p.isIdentLike() && p.cur.Str == "SEGMENT":
					var parseErr228 error
					stmt.Segment, parseErr228 = p.parseSegmentManagementClause()
					if parseErr228 != nil {
						return nil, parseErr228
					}

				case p.isIdentLike() && p.cur.Str == "FLASHBACK":
					p.advance()
					if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
						p.advance()
						stmt.Flashback = "ON"
					} else if p.isIdentLike() && p.cur.Str == "OFF" {
						p.advance()
						stmt.Flashback = "OFF"
					}

				default:
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDropTablespaceStmt parses a DROP TABLESPACE or DROP TABLESPACE SET statement.
//
// BNF: oracle/parser/bnf/DROP-TABLESPACE.bnf
//
//	DROP TABLESPACE [ IF EXISTS ] tablespace
//	    [ { DROP | KEEP } QUOTA ]
//	    [ INCLUDING CONTENTS [ { AND DATAFILES | KEEP DATAFILES } ]
//	      [ CASCADE CONSTRAINTS ] ] ;
//
// BNF: oracle/parser/bnf/DROP-TABLESPACE-SET.bnf
//
//	DROP TABLESPACE SET tablespace_set
//	    [ INCLUDING CONTENTS [ { AND DATAFILES | KEEP DATAFILES } ]
//	      [ CASCADE CONSTRAINTS ] ] ;
func (p *Parser) parseDropTablespaceStmt(start int, isSet bool) (*nodes.DropTablespaceStmt, error) {
	stmt := &nodes.DropTablespaceStmt{
		IsSet: isSet,
		Loc:   nodes.Loc{Start: start},
	}

	// Optional IF EXISTS (only for non-SET)
	if !isSet && p.cur.Type == kwIF {
		p.advance() // IF
		if p.cur.Type == kwEXISTS {
			p.advance() // EXISTS
		}
		stmt.IfExists = true
	}
	var parseErr229 error

	stmt.Name, parseErr229 = p.parseObjectName()
	if parseErr229 !=

		// Parse optional clauses
		nil {
		return nil, parseErr229
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwDROP:
			p.advance() // DROP
			if p.isIdentLike() && p.cur.Str == "QUOTA" {
				p.advance()
				stmt.DropQuota = true
			}

		case p.isIdentLike() && p.cur.Str == "KEEP":
			p.advance() // KEEP
			if p.isIdentLike() && p.cur.Str == "QUOTA" {
				p.advance()
				stmt.KeepQuota = true
			} else if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
				stmt.KeepDatafiles = true
			}

		case p.isIdentLike() && p.cur.Str == "INCLUDING":
			p.advance() // INCLUDING
			if p.isIdentLike() && p.cur.Str == "CONTENTS" {
				p.advance()
				stmt.IncludingContents = true
			}

		case p.cur.Type == kwAND:
			p.advance() // AND
			if p.isIdentLike() && p.cur.Str == "DATAFILES" {
				p.advance()
				stmt.AndDatafiles = true
			}

		case p.cur.Type == kwCASCADE:
			p.advance() // CASCADE
			if p.cur.Type == kwCONSTRAINT || (p.isIdentLike() && p.cur.Str == "CONSTRAINTS") {
				p.advance()
				stmt.CascadeConstraints = true
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
//
// parseCreateClusterStmt parses a CREATE CLUSTER statement.
//
// BNF: oracle/parser/bnf/CREATE-CLUSTER.bnf
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
//	    [ PCTFREE integer ]
//	    [ PCTUSED integer ]
//	    [ INITRANS integer ]
//	    [ storage_clause ]
//
//	parallel_clause:
//	    { NOPARALLEL | PARALLEL [ integer ] }
//
//	cluster_range_partitions:
//	    PARTITION BY RANGE ( column [, column ]... )
//	    ( PARTITION [ partition ]
//	        VALUES LESS THAN ( { value | MAXVALUE } [, { value | MAXVALUE } ]... )
//	        [ table_partition_description ]
//	      [, PARTITION [ partition ]
//	          VALUES LESS THAN ( { value | MAXVALUE } [, { value | MAXVALUE } ]... )
//	          [ table_partition_description ] ]...
//	    )
func (p *Parser) parseCreateClusterStmt(start int) (*nodes.CreateClusterStmt, error) {
	stmt := &nodes.CreateClusterStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // IF
			p.advance() // NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // EXISTS
			}
		}
	}
	var parseErr234 error

	stmt.Name, parseErr234 = p.parseObjectName()
	if parseErr234 !=

		// Parse column list in parentheses: ( col datatype [SORT] [, ...] )
		nil {
		return nil, parseErr234
	}

	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col := &nodes.ClusterColumn{
				Loc: nodes.Loc{Start: p.pos()},
			}
			var parseErr235 error
			col.Name, parseErr235 = p.parseIdentifier()
			if parseErr235 != nil {
				return nil, parseErr235
			}

			// Optional COLLATE
			var parseErr236 error
			col.DataType, parseErr236 = p.parseTypeName()
			if parseErr236 != nil {
				return nil, parseErr236
			}

			if p.isIdentLike() && p.cur.Str == "COLLATE" {
				p.advance()
				parseDiscard238, parseErr237 := p.parseIdentifier()
				_ = // collation name
					parseDiscard238
				if parseErr237 !=

					// Optional SORT
					nil {
					return nil, parseErr237
				}
			}

			if p.cur.Type == kwORDER || (p.isIdentLike() && p.cur.Str == "SORT") {
				p.advance()
				col.Sort = true
			}
			col.Loc.End = p.prev.End
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
				v, parseErr239 := p.parseIntValue()
				if parseErr239 != nil {
					return nil, parseErr239
				}
				stmt.PctFree = &v
			}

		case p.isIdentLike() && p.cur.Str == "PCTUSED":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr240 := p.parseIntValue()
				if parseErr240 != nil {
					return nil, parseErr240
				}
				stmt.PctUsed = &v
			}

		case p.isIdentLike() && p.cur.Str == "INITRANS":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr241 := p.parseIntValue()
				if parseErr241 != nil {
					return nil, parseErr241
				}
				stmt.InitTrans = &v
			}

		case p.cur.Type == kwSIZE:
			p.advance()
			var parseErr242 error
			stmt.Size, parseErr242 = p.parseSizeValue()
			if parseErr242 != nil {
				return nil, parseErr242
			}

		case p.cur.Type == kwTABLESPACE:
			p.advance()
			var parseErr243 error
			stmt.Tablespace, parseErr243 = p.parseIdentifier()
			if parseErr243 != nil {
				return nil, parseErr243
			}

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
				var parseErr244 error
				stmt.HashExpr, parseErr244 = p.parseExpr()
				if parseErr244 != nil {
					return nil, parseErr244
				}
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
			parseDiscard246, parseErr245 := p.parseIdentifier()
			_ = // METADATA, DATA, NONE, etc.
				parseDiscard246
			if parseErr245 != nil {
				return nil, parseErr245
			}

		case p.isIdentLike() && p.cur.Str == "PARTITION":
			// cluster_range_partitions - skip until end
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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

// parseCreateDimensionStmt parses a CREATE DIMENSION statement.
//
// BNF: oracle/parser/bnf/CREATE-DIMENSION.bnf
//
//	CREATE DIMENSION [ schema. ] dimension
//	    level_clause [ level_clause ]...
//	    { hierarchy_clause
//	    | dimension_join_clause
//	    | attribute_clause
//	    | extended_attribute_clause
//	    }... ;
//
//	level_clause:
//	    LEVEL level IS ( [ schema. ] table. column [, [ schema. ] table. column ]... )
//	    [ SKIP WHEN NULL ]
//
//	hierarchy_clause:
//	    HIERARCHY hierarchy (
//	        child_level CHILD OF parent_level
//	        [ CHILD OF parent_level ]...
//	        [ dimension_join_clause ]...
//	    )
//
//	dimension_join_clause:
//	    JOIN KEY ( [ [ schema. ] table. ] child_key_column
//	        [, [ [ schema. ] table. ] child_key_column ]... )
//	    REFERENCES parent_level
//
//	attribute_clause:
//	    ATTRIBUTE level DETERMINES ( [ schema. ] table. column
//	        [, [ schema. ] table. column ]... )
//
//	extended_attribute_clause:
//	    ATTRIBUTE attribute LEVEL level DETERMINES ( [ schema. ] table. column
//	        [, [ schema. ] table. column ]... )
func (p *Parser) parseCreateDimensionStmt(start int) (*nodes.CreateDimensionStmt, error) {
	stmt := &nodes.CreateDimensionStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr247 error

	stmt.Name, parseErr247 = p.parseObjectName()
	if parseErr247 !=

		// Parse clauses: LEVEL, HIERARCHY, ATTRIBUTE
		nil {
		return nil, parseErr247
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwLEVEL:
			parseValue13, parseErr14 := p.parseDimensionLevel()
			if parseErr14 != nil {
				return nil, parseErr14
			}
			stmt.Levels = append(stmt.Levels, parseValue13)

		case p.isIdentLike() && p.cur.Str == "HIERARCHY":
			parseValue15, parseErr16 := p.parseDimensionHierarchy()
			if parseErr16 != nil {
				return nil, parseErr16
			}
			stmt.Hierarchies = append(stmt.Hierarchies, parseValue15)

		case p.isIdentLike() && p.cur.Str == "ATTRIBUTE":
			parseValue17, parseErr18 := p.parseDimensionAttribute()
			if parseErr18 != nil {
				return nil, parseErr18
			}
			stmt.Attributes = append(stmt.Attributes, parseValue17)

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDimensionLevel parses a LEVEL clause.
//
//	LEVEL level IS ( level_table.level_column [, ...] ) [ SKIP WHEN NULL ]
//	-- or without parens:
//	LEVEL level IS level_table.level_column [ SKIP WHEN NULL ]
func (p *Parser) parseDimensionLevel() (*nodes.DimensionLevel, error) {
	lvl := &nodes.DimensionLevel{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance()
	var // consume LEVEL
	parseErr248 error

	lvl.Name, parseErr248 = p.parseIdentifier()
	if parseErr248 !=

		// IS
		nil {
		return nil, parseErr248
	}

	if p.isIdentLike() && p.cur.Str == "IS" {
		p.advance()
	}

	// Column list: ( col [, ...] ) or just col
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			parseValue19, parseErr20 := p.parseObjectName()
			if parseErr20 != nil {
				return nil, parseErr20
			}
			lvl.Columns = append(lvl.Columns, parseValue19)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	} else {
		parseValue21,
			// Single column reference: table.column
			parseErr22 := p.parseObjectName()
		if parseErr22 !=

			// Optional SKIP WHEN NULL
			nil {
			return nil, parseErr22
		}
		lvl.Columns = append(lvl.Columns, parseValue21)
	}

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

	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// parseDimensionHierarchy parses a HIERARCHY clause.
//
//	HIERARCHY hierarchy_name (
//	    child_level CHILD OF parent_level [ CHILD OF ... ]
//	    [ JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level ] ...
//	)
func (p *Parser) parseDimensionHierarchy() (*nodes.DimensionHierarchy, error) {
	hier := &nodes.DimensionHierarchy{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance()
	var // consume HIERARCHY
	parseErr249 error

	hier.Name, parseErr249 = p.parseIdentifier()
	if parseErr249 !=

		// Parse ( ... )
		nil {
		return nil, parseErr249
	}

	if p.cur.Type == '(' {
		p.advance()

		// Parse level chain: child CHILD OF parent CHILD OF ...
		// First level
		if p.isIdentLike() {
			parseValue23, parseErr24 := p.parseIdentifier()
			if parseErr24 != nil {
				return nil, parseErr24
			}
			hier.Levels = append(hier.Levels, parseValue23)
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
					parseValue25, parseErr26 := p.parseIdentifier()
					if parseErr26 != nil {
						return nil, parseErr26
					}
					hier.Levels = append(hier.Levels, parseValue25)
				}
			} else if p.isIdentLike() && p.cur.Str == "JOIN" {
				// JOIN KEY clause
				jk, parseErr250 := p.parseDimensionJoinKey()
				if parseErr250 != nil {
					return nil, parseErr250
				}
				hier.JoinKeys = append(hier.JoinKeys, jk)
			} else {
				break
			}
		}

		if p.cur.Type == ')' {
			p.advance()
		}
	}

	hier.Loc.End = p.prev.End
	return hier, nil
}

// parseDimensionJoinKey parses a JOIN KEY clause.
//
//	JOIN KEY ( child_key_column [, ...] ) REFERENCES parent_level
func (p *Parser) parseDimensionJoinKey() (*nodes.DimensionJoinKey, error) {
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
			parseValue27, parseErr28 := p.parseObjectName()
			if parseErr28 != nil {
				return nil, parseErr28
			}
			jk.ChildKeys = append(jk.ChildKeys, parseValue27)
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
		var parseErr251 error
		jk.ParentLevel, parseErr251 = p.parseIdentifier()
		if parseErr251 != nil {
			return nil, parseErr251
		}
	}

	jk.Loc.End = p.prev.End
	return jk, nil
}

// parseDimensionAttribute parses an ATTRIBUTE clause.
//
//	ATTRIBUTE level DETERMINES ( dependent_column [, ...] )
//	-- or extended form:
//	ATTRIBUTE attr_name LEVEL level DETERMINES ( dependent_column [, ...] )
func (p *Parser) parseDimensionAttribute() (*nodes.DimensionAttribute, error) {
	attr := &nodes.DimensionAttribute{
		Loc: nodes.Loc{Start: p.pos()},
	}
	p.advance()
	var // consume ATTRIBUTE
	parseErr252 error

	attr.AttrName, parseErr252 = p.parseIdentifier()
	if parseErr252 !=

		// Check for extended form: LEVEL level DETERMINES ...
		nil {
		return nil, parseErr252
	}

	if p.cur.Type == kwLEVEL {
		p.advance()
		var parseErr253 error
		attr.LevelName, parseErr253 = p.parseIdentifier()
		if parseErr253 !=

			// DETERMINES
			nil {
			return nil, parseErr253
		}
	}

	if p.isIdentLike() && p.cur.Str == "DETERMINES" {
		p.advance()
	}

	// ( dependent_column [, ...] ) or single column
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			parseValue29, parseErr30 := p.parseObjectName()
			if parseErr30 != nil {
				return nil, parseErr30
			}
			attr.Columns = append(attr.Columns, parseValue29)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	} else {
		parseValue31, parseErr32 := p.parseObjectName()
		if parseErr32 != nil {
			return nil, parseErr32
		}
		attr.Columns = append(attr.Columns, parseValue31)
	}

	attr.Loc.End = p.prev.End
	return attr, nil
}

// ---------------------------------------------------------------------------
// ALTER CLUSTER
// ---------------------------------------------------------------------------

// parseAlterClusterStmt parses an ALTER CLUSTER statement.
//
// BNF: oracle/parser/bnf/ALTER-CLUSTER.bnf
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
//	    [ PCTUSED integer ]
//	    [ PCTFREE integer ]
//	    [ INITRANS integer ]
//	    [ STORAGE storage_clause ]
//
//	allocate_extent_clause:
//	    ALLOCATE EXTENT
//	    [ ( { SIZE size_clause
//	        | DATAFILE 'filename'
//	        | INSTANCE integer
//	        }...
//	      )
//	    ]
//
//	deallocate_unused_clause:
//	    DEALLOCATE UNUSED [ KEEP size_clause ]
//
//	parallel_clause:
//	    { PARALLEL [ integer ] | NOPARALLEL }
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
	var parseErr254 error

	stmt.Name, parseErr254 = p.parseObjectName()
	if parseErr254 !=

		// Parse the action
		nil {
		return nil, parseErr254
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwSIZE:
			p.advance()
			stmt.Action = "SIZE"
			var parseErr255 error
			stmt.Size, parseErr255 = p.parseSizeValue()
			if parseErr255 != nil {
				return nil, parseErr255
			}

		case p.isIdentLike() && p.cur.Str == "PCTUSED":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr256 := p.parseIntValue()
				if parseErr256 != nil {
					return nil, parseErr256
				}
				stmt.PctUsed = &v
			}
			if stmt.Action == "" {
				stmt.Action = "PHYSICAL_ATTRS"
			}

		case p.cur.Type == kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr257 := p.parseIntValue()
				if parseErr257 != nil {
					return nil, parseErr257
				}
				stmt.PctFree = &v
			}
			if stmt.Action == "" {
				stmt.Action = "PHYSICAL_ATTRS"
			}

		case p.isIdentLike() && p.cur.Str == "INITRANS":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr258 := p.parseIntValue()
				if parseErr258 != nil {
					return nil, parseErr258
				}
				stmt.InitTrans = &v
			}
			if stmt.Action == "" {
				stmt.Action = "PHYSICAL_ATTRS"
			}

		case p.cur.Type == kwSTORAGE:
			p.advance()
			if p.cur.Type == '(' {
				p.skipParens()
			}
			if stmt.Action == "" {
				stmt.Action = "PHYSICAL_ATTRS"
			}

		case p.isIdentLike() && p.cur.Str == "MODIFY":
			p.advance() // MODIFY
			stmt.Action = "MODIFY_PARTITION"
			if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
				p.advance() // PARTITION
			}
			var parseErr259 error
			stmt.ModifyPartition, parseErr259 = p.parseIdentifier()
			if parseErr259 !=
				// allocate_extent_clause follows
				nil {
				return nil, parseErr259
			}

			if p.isIdentLike() && p.cur.Str == "ALLOCATE" {
				p.advance() // ALLOCATE
				if p.isIdentLike() && p.cur.Str == "EXTENT" {
					p.advance() // EXTENT
				}
				if p.cur.Type == '(' {
					p.skipParens()
				}
			}

		case p.isIdentLike() && p.cur.Str == "ALLOCATE":
			p.advance() // ALLOCATE
			stmt.Action = "ALLOCATE_EXTENT"
			if p.isIdentLike() && p.cur.Str == "EXTENT" {
				p.advance() // EXTENT
			}
			if p.cur.Type == '(' {
				p.skipParens()
			}

		case p.isIdentLike() && p.cur.Str == "DEALLOCATE":
			p.advance() // DEALLOCATE
			stmt.Action = "DEALLOCATE_UNUSED"
			if p.isIdentLike() && p.cur.Str == "UNUSED" {
				p.advance() // UNUSED
			}
			if p.isIdentLike() && p.cur.Str == "KEEP" {
				p.advance()
				parseDiscard261, // KEEP
					parseErr260 := p.parseSizeValue()
				_ = parseDiscard261
				if parseErr260 != nil {
					return nil, parseErr260
				}
			}

		case p.isIdentLike() && p.cur.Str == "PARALLEL":
			p.advance()
			stmt.Action = "PARALLEL"
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			} else {
				stmt.Parallel = "PARALLEL"
			}

		case p.isIdentLike() && p.cur.Str == "NOPARALLEL":
			p.advance()
			stmt.Action = "PARALLEL"
			stmt.Parallel = "NOPARALLEL"

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER DIMENSION
// ---------------------------------------------------------------------------

// parseAlterDimensionStmt parses an ALTER DIMENSION statement.
//
// BNF: oracle/parser/bnf/ALTER-DIMENSION.bnf
//
//	ALTER DIMENSION [ schema . ] dimension
//	    { ADD level_clause
//	    | ADD hierarchy_clause
//	    | ADD attribute_clause
//	    | ADD extended_attribute_clause
//	    | DROP level_clause
//	    | DROP hierarchy_clause
//	    | DROP attribute_clause
//	    | DROP extended_attribute_clause
//	    | COMPILE
//	    }
//
//	level_clause:
//	    LEVEL level IS ( table . column [, table . column ]... )
//	    [ SKIP WHEN expression ]
//
//	hierarchy_clause:
//	    HIERARCHY hierarchy_name ( child_level CHILD OF parent_level
//	      [ JOIN KEY child_key_column REFERENCES parent_level ]
//	      [, child_level CHILD OF parent_level
//	        [ JOIN KEY child_key_column REFERENCES parent_level ] ]...
//	    )
//
//	attribute_clause:
//	    ATTRIBUTE level DETERMINES ( dependent_column [, dependent_column ]... )
//
//	extended_attribute_clause:
//	    ATTRIBUTE attribute_name OF level_name
//	      DETERMINES ( dependent_column [, dependent_column ]... )
func (p *Parser) parseAlterDimensionStmt(start int) (*nodes.AlterDimensionStmt, error) {
	stmt := &nodes.AlterDimensionStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr262 error

	stmt.Name, parseErr262 = p.parseObjectName()
	if parseErr262 !=

		// Parse actions
		nil {
		return nil, parseErr262
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.isIdentLike() && p.cur.Str == "COMPILE":
			p.advance()
			stmt.Compile = true

		case p.cur.Type == kwADD:
			p.advance() // ADD
			switch {
			case p.cur.Type == kwLEVEL:
				parseValue33, parseErr34 := p.parseDimensionLevel()
				if parseErr34 != nil {
					return nil, parseErr34
				}
				stmt.AddLevels = append(stmt.AddLevels, parseValue33)
			case p.isIdentLike() && p.cur.Str == "HIERARCHY":
				parseValue35, parseErr36 := p.parseDimensionHierarchy()
				if parseErr36 != nil {
					return nil, parseErr36
				}
				stmt.AddHierarchies = append(stmt.AddHierarchies, parseValue35)
			case p.isIdentLike() && p.cur.Str == "ATTRIBUTE":
				parseValue37, parseErr38 := p.parseDimensionAttribute()
				if parseErr38 != nil {
					return nil, parseErr38
				}
				stmt.AddAttributes = append(stmt.AddAttributes, parseValue37)
			default:
				p.advance()
			}

		case p.cur.Type == kwDROP:
			p.advance() // DROP
			switch {
			case p.cur.Type == kwLEVEL:
				p.advance()
				parseValue39, // LEVEL
					parseErr40 := p.parseIdentifier()
				if parseErr40 != nil {
					return nil, parseErr40
				}
				stmt.DropLevels = append(stmt.DropLevels, parseValue39)
			case p.isIdentLike() && p.cur.Str == "HIERARCHY":
				p.advance()
				parseValue41, // HIERARCHY
					parseErr42 := p.parseIdentifier()
				if parseErr42 != nil {
					return nil, parseErr42
				}
				stmt.DropHierarchies = append(stmt.DropHierarchies, parseValue41)
			case p.isIdentLike() && p.cur.Str == "ATTRIBUTE":
				p.advance()
				parseValue43, // ATTRIBUTE
					parseErr44 := p.parseIdentifier()
				if parseErr44 != nil {
					return nil, parseErr44
				}
				stmt.DropAttributes = append(stmt.DropAttributes, parseValue43)
			default:
				p.advance()
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE MATERIALIZED ZONEMAP
// ---------------------------------------------------------------------------

// parseCreateMaterializedZonemapStmt parses a CREATE MATERIALIZED ZONEMAP statement.
//
// BNF: oracle/parser/bnf/CREATE-MATERIALIZED-ZONEMAP.bnf
//
//	CREATE MATERIALIZED ZONEMAP [ IF NOT EXISTS ]
//	    [ schema. ] zonemap_name
//	    [ zonemap_attributes ]
//	    [ zonemap_refresh_clause ]
//	    [ { ENABLE | DISABLE } PRUNING ]
//	    { create_zonemap_on_table | create_zonemap_as_subquery }
//
//	create_zonemap_on_table:
//	    ON [ schema. ] table ( column [, column ]... )
//
//	create_zonemap_as_subquery:
//	    [ ( column_alias [, column_alias ]... ) ]
//	    AS query_block
//
//	zonemap_attributes:
//	    [ TABLESPACE tablespace_name ]
//	    [ SCALE integer ]
//	    [ PCTFREE integer ]
//	    [ PCTUSED integer ]
//	    [ { CACHE | NOCACHE } ]
//
//	zonemap_refresh_clause:
//	    REFRESH
//	    [ { FAST | COMPLETE | FORCE } ]
//	    [ { ON DEMAND
//	      | ON COMMIT
//	      | ON LOAD
//	      | ON DATA MOVEMENT
//	      | ON LOAD DATA MOVEMENT
//	      }
//	    ]
func (p *Parser) parseCreateMaterializedZonemapStmt(start int) (*nodes.CreateMaterializedZonemapStmt, error) {
	stmt := &nodes.CreateMaterializedZonemapStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // IF
			p.advance() // NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // EXISTS
			}
			stmt.IfNotExists = true
		}
	}
	var parseErr263 error

	stmt.Name, parseErr263 = p.parseObjectName()
	if parseErr263 !=

		// Parse options until ON, AS, (, or ; / EOF
		nil {
		return nil, parseErr263
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		// zonemap_attributes
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			var parseErr264 error
			stmt.Tablespace, parseErr264 = p.parseIdentifier()
			if parseErr264 != nil {
				return nil, parseErr264
			}

		case p.isIdentLike() && p.cur.Str == "SCALE":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr265 := p.parseIntValue()
				if parseErr265 != nil {
					return nil, parseErr265
				}
				stmt.Scale = &v
			}

		case p.cur.Type == kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr266 := p.parseIntValue()
				if parseErr266 != nil {
					return nil, parseErr266
				}
				stmt.PctFree = &v
			}

		case p.isIdentLike() && p.cur.Str == "PCTUSED":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr267 := p.parseIntValue()
				if parseErr267 != nil {
					return nil, parseErr267
				}
				stmt.PctUsed = &v
			}

		case p.cur.Type == kwCACHE:
			p.advance()
			stmt.Cache = true

		case p.cur.Type == kwNOCACHE:
			p.advance()
			stmt.NoCache = true

		// zonemap_refresh_clause
		case p.isIdentLike() && p.cur.Str == "REFRESH":
			p.advance()
			parseErr268 := // REFRESH
				p.parseZonemapRefresh(stmt, nil)
			if parseErr268 !=

				// ENABLE/DISABLE PRUNING
				nil {
				return nil, parseErr268
			}

		case p.cur.Type == kwENABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "PRUNING" {
				p.advance()
				stmt.EnablePruning = true
			}

		case p.cur.Type == kwDISABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "PRUNING" {
				p.advance()
				stmt.DisablePruning = true
			}

		// create_zonemap_on_table
		case p.cur.Type == kwON:
			p.advance()
			var // ON
			parseErr269 error
			stmt.OnTable, parseErr269 = p.parseObjectName()
			if parseErr269 != nil {
				return nil, parseErr269
			}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					parseValue45, parseErr46 := p.parseIdentifier()
					if parseErr46 != nil {
						return nil, parseErr46
					}
					stmt.OnColumns = append(stmt.OnColumns, parseValue45)
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}

		// create_zonemap_as_subquery: ( column_alias, ... ) AS query_block
		case p.cur.Type == '(' && stmt.OnTable == nil && stmt.AsQuery == nil:
			// column aliases before AS
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseValue47, parseErr48 := p.parseIdentifier()
				if parseErr48 != nil {
					return nil, parseErr48
				}
				stmt.ColumnAliases = append(stmt.ColumnAliases, parseValue47)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}

		case p.cur.Type == kwAS:
			p.advance()
			var // AS
			parseErr270 error
			stmt.AsQuery, parseErr270 = p.parseSelectStmt()
			if parseErr270 != nil {
				return nil, parseErr270
			}

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseZonemapRefresh parses the refresh clause for a zonemap.
// It sets refresh fields on either a CreateMaterializedZonemapStmt or AlterMaterializedZonemapStmt.
func (p *Parser) parseZonemapRefresh(create *nodes.CreateMaterializedZonemapStmt, alter *nodes.AlterMaterializedZonemapStmt) error {
	// Optional method: FAST | COMPLETE | FORCE
	method := ""
	if p.isIdentLike() {
		switch p.cur.Str {
		case "FAST", "COMPLETE", "FORCE":
			method = p.cur.Str
			p.advance()
		}
	}

	// Optional ON ...
	refreshOn := ""
	if p.cur.Type == kwON {
		p.advance() // ON
		if p.isIdentLike() {
			switch p.cur.Str {
			case "DEMAND":
				refreshOn = "ON DEMAND"
				p.advance()
			case "COMMIT":
				refreshOn = "ON COMMIT"
				p.advance()
			case "LOAD":
				p.advance()
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance() // DATA
					if p.isIdentLike() && p.cur.Str == "MOVEMENT" {
						p.advance() // MOVEMENT
					}
					refreshOn = "ON LOAD DATA MOVEMENT"
				} else {
					refreshOn = "ON LOAD"
				}
			case "DATA":
				p.advance() // DATA
				if p.isIdentLike() && p.cur.Str == "MOVEMENT" {
					p.advance() // MOVEMENT
				}
				refreshOn = "ON DATA MOVEMENT"
			case "STATEMENT":
				refreshOn = "ON STATEMENT"
				p.advance()
			}
		}
	}

	if create != nil {
		create.RefreshMethod = method
		create.RefreshOn = refreshOn
	}
	if alter != nil {
		alter.RefreshMethod = method
		alter.RefreshOn = refreshOn
	}
	return nil
}

// ---------------------------------------------------------------------------
// ALTER MATERIALIZED ZONEMAP
// ---------------------------------------------------------------------------

// parseAlterMaterializedZonemapStmt parses an ALTER MATERIALIZED ZONEMAP statement.
//
// BNF: oracle/parser/bnf/ALTER-MATERIALIZED-ZONEMAP.bnf
//
//	ALTER MATERIALIZED ZONEMAP [ IF EXISTS ] [ schema. ] zonemap_name
//	    { alter_zonemap_attributes
//	    | zonemap_refresh_clause
//	    | { ENABLE | DISABLE } PRUNING
//	    | COMPILE
//	    | REBUILD
//	    | UNUSABLE
//	    } ;
//
//	alter_zonemap_attributes:
//	    [ PCTFREE integer ]
//	    [ PCTUSED integer ]
//	    [ { CACHE | NOCACHE } ]
//
//	zonemap_refresh_clause:
//	    REFRESH
//	    [ { FAST | COMPLETE | FORCE } ]
//	    [ { ON COMMIT | ON DEMAND | ON LOAD | ON DATA MOVEMENT | ON STATEMENT } ]
func (p *Parser) parseAlterMaterializedZonemapStmt(start int) (*nodes.AlterMaterializedZonemapStmt, error) {
	stmt := &nodes.AlterMaterializedZonemapStmt{
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
	var parseErr271 error

	stmt.Name, parseErr271 = p.parseObjectName()
	if parseErr271 !=

		// Parse action
		nil {
		return nil, parseErr271
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch {
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr272 := p.parseIntValue()
				if parseErr272 != nil {
					return nil, parseErr272
				}
				stmt.PctFree = &v
			}
			if stmt.Action == "" {
				stmt.Action = "ATTRS"
			}

		case p.isIdentLike() && p.cur.Str == "PCTUSED":
			p.advance()
			if p.cur.Type == tokICONST {
				v, parseErr273 := p.parseIntValue()
				if parseErr273 != nil {
					return nil, parseErr273
				}
				stmt.PctUsed = &v
			}
			if stmt.Action == "" {
				stmt.Action = "ATTRS"
			}

		case p.cur.Type == kwCACHE:
			p.advance()
			stmt.Cache = true
			if stmt.Action == "" {
				stmt.Action = "ATTRS"
			}

		case p.cur.Type == kwNOCACHE:
			p.advance()
			stmt.NoCache = true
			if stmt.Action == "" {
				stmt.Action = "ATTRS"
			}

		case p.isIdentLike() && p.cur.Str == "REFRESH":
			p.advance()
			stmt.Action = "REFRESH"
			parseErr274 := p.parseZonemapRefresh(nil, stmt)
			if parseErr274 != nil {
				return nil, parseErr274
			}

		case p.cur.Type == kwENABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "PRUNING" {
				p.advance()
			}
			stmt.Action = "ENABLE_PRUNING"

		case p.cur.Type == kwDISABLE:
			p.advance()
			if p.isIdentLike() && p.cur.Str == "PRUNING" {
				p.advance()
			}
			stmt.Action = "DISABLE_PRUNING"

		case p.isIdentLike() && p.cur.Str == "COMPILE":
			p.advance()
			stmt.Action = "COMPILE"

		case p.isIdentLike() && p.cur.Str == "REBUILD":
			p.advance()
			stmt.Action = "REBUILD"

		case p.isIdentLike() && p.cur.Str == "UNUSABLE":
			p.advance()
			stmt.Action = "UNUSABLE"

		default:
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// parseCreateInmemoryJoinGroupStmt parses a CREATE INMEMORY JOIN GROUP statement.
//
// BNF: oracle/parser/bnf/CREATE-INMEMORY-JOIN-GROUP.bnf
//
//	CREATE INMEMORY JOIN GROUP [ IF NOT EXISTS ] [ schema. ] join_group
//	    ( [ schema. ] table ( column )
//	      [, [ schema. ] table ( column ) ]... ) ;
func (p *Parser) parseCreateInmemoryJoinGroupStmt(start int) (*nodes.CreateInmemoryJoinGroupStmt, error) {
	stmt := &nodes.CreateInmemoryJoinGroupStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // IF
			p.advance() // NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // EXISTS
			}
			stmt.IfNotExists = true
		}
	}
	var parseErr275 error

	stmt.Name, parseErr275 = p.parseObjectName()
	if parseErr275 !=

		// Parse member list: ( table(col) [, table(col) ]... )
		nil {
		return nil, parseErr275
	}

	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			member, parseErr276 := p.parseJoinGroupMember()
			if parseErr276 != nil {
				return nil, parseErr276
			}
			stmt.Members = append(stmt.Members, member)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseJoinGroupMember parses a table(column) member of a join group.
func (p *Parser) parseJoinGroupMember() (*nodes.JoinGroupMember, error) {
	m := &nodes.JoinGroupMember{
		Loc: nodes.Loc{Start: p.pos()},
	}
	var parseErr277 error
	m.Table, parseErr277 = p.parseObjectName()
	if parseErr277 != nil {
		return nil, parseErr277
	}
	if p.cur.Type == '(' {
		p.advance()
		var parseErr278 error
		m.Column, parseErr278 = p.parseIdentifier()
		if parseErr278 != nil {
			return nil, parseErr278
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}
	m.Loc.End = p.prev.End
	return m, nil
}

// ---------------------------------------------------------------------------
// ALTER INMEMORY JOIN GROUP
// ---------------------------------------------------------------------------

// parseAlterInmemoryJoinGroupStmt parses an ALTER INMEMORY JOIN GROUP statement.
//
// BNF: oracle/parser/bnf/ALTER-INMEMORY-JOIN-GROUP.bnf
//
//	ALTER INMEMORY JOIN GROUP [ IF EXISTS ] [ schema. ] join_group
//	    { ADD | REMOVE } ( [ schema. ] table ( column ) ) ;
func (p *Parser) parseAlterInmemoryJoinGroupStmt(start int) (*nodes.AlterInmemoryJoinGroupStmt, error) {
	stmt := &nodes.AlterInmemoryJoinGroupStmt{
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
	var parseErr279 error

	stmt.Name, parseErr279 = p.parseObjectName()
	if parseErr279 !=

		// ADD or REMOVE
		nil {
		return nil, parseErr279
	}

	if p.cur.Type == kwADD {
		p.advance()
		stmt.Action = "ADD"
	} else if p.isIdentLike() && p.cur.Str == "REMOVE" {
		p.advance()
		stmt.Action = "REMOVE"
	}

	// ( table(column) )
	if p.cur.Type == '(' {
		p.advance()
		var parseErr280 error
		stmt.Member, parseErr280 = p.parseJoinGroupMember()
		if parseErr280 != nil {
			return nil, parseErr280
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

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
		if p.cur.Type == kwTABLE || (p.isIdentLike() && p.cur.Str == "TABLES") {
			p.advance() // TABLES
		}
		// Optional CASCADE CONSTRAINTS
		if p.cur.Type == kwCASCADE {
			p.advance() // CASCADE
			if p.isIdentLike() && p.cur.Str == "CONSTRAINTS" {
				p.advance() // CONSTRAINTS
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
				if p.cur.Type == kwDATABASE {
					p.advance() // consume DATABASE
				}
				return p.parseAlterPluggableDatabaseStmt(start)
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
				if p.cur.Type == kwPROFILE {
					p.advance() // consume PROFILE
				}
				return p.parseAlterLockdownProfileStmt(start)
			case "OUTLINE":
				p.advance()
				return p.parseAlterOutlineStmt(start)
			case "INMEMORY":
				p.advance() // consume INMEMORY
				if p.cur.Type == kwJOIN {
					p.advance() // consume JOIN
				}
				if p.cur.Type == kwGROUP || (p.isIdentLike() && p.cur.Str == "GROUP") {
					p.advance() // consume GROUP
				}
				return p.parseAlterInmemoryJoinGroupStmt(start)
			case "FLASHBACK":
				p.advance() // consume FLASHBACK
				if p.isIdentLike() && p.cur.Str == "ARCHIVE" {
					p.advance() // consume ARCHIVE
				}
				return p.parseAlterFlashbackArchiveStmt(start)
			case "RESOURCE":
				p.advance() // consume RESOURCE
				if p.isIdentLike() && p.cur.Str == "COST" {
					p.advance() // consume COST
				}
				return p.parseAlterResourceCostStmt(start)
			case "ROLLBACK":
				p.advance() // consume ROLLBACK
				if p.isIdentLike() && p.cur.Str == "SEGMENT" {
					p.advance() // consume SEGMENT
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
// CREATE ATTRIBUTE DIMENSION
// ---------------------------------------------------------------------------

// parseCreateAttributeDimensionStmt parses a CREATE ATTRIBUTE DIMENSION statement.
//
// BNF: oracle/parser/bnf/CREATE-ATTRIBUTE-DIMENSION.bnf
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] ATTRIBUTE DIMENSION
//	    [ IF NOT EXISTS ] [ schema. ] attr_dimension
//	    [ SHARING = { METADATA | NONE } ]
//	    [ classification_clause ]
//	    [ DIMENSION TYPE { STANDARD | TIME } ]
//	    attr_dim_using_clause
//	    attributes_clause
//	    attr_dim_level_clause [ attr_dim_level_clause ]...
//	    [ all_clause ] ;
//
//	classification_clause:
//	    { CAPTION 'caption'
//	    | DESCRIPTION 'description'
//	    | CLASSIFICATION classification_name [ LANGUAGE language ] VALUE 'value'
//	    } [ classification_clause ]
//
//	attr_dim_using_clause:
//	    USING source_clause [, source_clause ]...
//
//	source_clause:
//	    [ REMOTE ] [ schema. ] { table | view } [ [ AS ] alias ]
//	    | join_path_clause
//
//	join_path_clause:
//	    JOIN PATH join_path_name ON join_condition
//
//	join_condition:
//	    join_condition_elem [ AND join_condition_elem ]...
//
//	join_condition_elem:
//	    [ table_alias. ] column = [ table_alias. ] column
//
//	attributes_clause:
//	    ATTRIBUTES ( attr_dim_attribute_clause [, attr_dim_attribute_clause ]... )
//
//	attr_dim_attribute_clause:
//	    column [ AS alias ]
//	    [ classification_clause ]
//
//	attr_dim_level_clause:
//	    LEVEL level_name
//	    [ LEVEL TYPE { STANDARD | YEARS | HALF_YEARS | QUARTERS | MONTHS | WEEKS | DAYS | HOURS | MINUTES | SECONDS } ]
//	    [ classification_clause ]
//	    key_clause
//	    [ alternate_key_clause ]
//	    [ MEMBER NAME expression ]
//	    [ MEMBER CAPTION expression ]
//	    [ MEMBER DESCRIPTION expression ]
//	    [ dim_order_clause ]
//	    [ DETERMINES ( attribute [, attribute ]... ) ]
//
//	key_clause:
//	    KEY { attribute [ NOT NULL | SKIP WHEN NULL ]
//	        | ( attribute [, attribute ]... ) }
//
//	alternate_key_clause:
//	    ALTERNATE KEY { attribute
//	                  | ( attribute [, attribute ]... ) }
//
//	dim_order_clause:
//	    ORDER BY { attribute [ ASC | DESC ]
//	             | ( attribute [ ASC | DESC ] [, attribute [ ASC | DESC ] ]... ) }
//
//	all_clause:
//	    ALL [ MEMBER NAME expression ]
//	        [ MEMBER CAPTION expression ]
//	        [ MEMBER DESCRIPTION expression ]
func (p *Parser) parseCreateAttributeDimensionStmt(start int, orReplace, force, noForce bool) (*nodes.CreateAttributeDimensionStmt, error) {
	stmt := &nodes.CreateAttributeDimensionStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance()
		p.advance()
		if p.cur.Type == kwEXISTS {
			p.advance()
		}
		stmt.IfNotExists = true
	}
	var parseErr282 error

	// name
	stmt.Name, parseErr282 = p.parseObjectName()
	if parseErr282 !=

		// SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr282
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}
	var parseErr283 error

	// classification_clause(s) at top level
	stmt.Classifications, parseErr283 = p.parseClassificationClauses()
	if parseErr283 !=

		// DIMENSION TYPE { STANDARD | TIME }
		nil {
		return nil, parseErr283
	}

	if p.isIdentLike() && p.cur.Str == "DIMENSION" {
		next := p.peekNext()
		if next.Type == kwTYPE {
			p.advance() // consume DIMENSION
			p.advance() // consume TYPE
			if p.isIdentLike() {
				stmt.DimensionType = p.cur.Str
				p.advance()
			}
		}
	}

	// USING source_clause [, source_clause]...
	if p.cur.Type == kwUSING {
		p.advance()
		stmt.Sources = &nodes.List{}
		for {
			src, parseErr284 := p.parseAttrDimSourceClause()
			if parseErr284 != nil {
				return nil, parseErr284
			}
			if src != nil {
				stmt.Sources.Items = append(stmt.Sources.Items, src)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume comma
		}
	}

	// ATTRIBUTES ( ... )
	if p.isIdentLikeStr("ATTRIBUTES") {
		p.advance()
		stmt.Attributes = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				attr, parseErr285 := p.parseAttrDimAttributeClause()
				if parseErr285 != nil {
					return nil, parseErr285
				}
				if attr != nil {
					stmt.Attributes.Items = append(stmt.Attributes.Items, attr)
				}
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// LEVEL clauses (one or more)
	stmt.Levels = &nodes.List{}
	for p.cur.Type == kwLEVEL {
		lvl, parseErr286 := p.parseAttrDimLevelClause()
		if parseErr286 != nil {
			return nil, parseErr286
		}
		if lvl != nil {
			stmt.Levels.Items = append(stmt.Levels.Items, lvl)
		}
	}

	// ALL clause
	if p.cur.Type == kwALL {
		p.advance()
		allc := &nodes.AttrDimAllClause{
			Loc: nodes.Loc{Start: p.pos()},
		}
		var parseErr287 error
		allc.MemberName, allc.MemberCaption, allc.MemberDesc, parseErr287 = p.parseMemberExprs()
		if parseErr287 != nil {
			return nil, parseErr287
		}
		allc.Loc.End = p.prev.End
		stmt.AllClause = allc
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseClassificationClauses parses zero or more classification_clause items.
func (p *Parser) parseClassificationClauses() (*nodes.List, error) {
	var items []nodes.Node
	for {
		if p.isIdentLikeStr("CAPTION") {
			p.advance()
			if p.cur.Type == tokSCONST {
				items = append(items, &nodes.DDLOption{Key: "CAPTION", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLikeStr("DESCRIPTION") {
			p.advance()
			if p.cur.Type == tokSCONST {
				items = append(items, &nodes.DDLOption{Key: "DESCRIPTION", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLikeStr("CLASSIFICATION") {
			p.advance()
			opt := &nodes.DDLOption{Key: "CLASSIFICATION"}
			if p.isIdentLike() || p.cur.Type == tokSCONST {
				opt.Value = p.cur.Str
				p.advance()
			}
			// [ LANGUAGE language ]
			if p.isIdentLikeStr("LANGUAGE") {
				p.advance()
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					opt.Value += " LANGUAGE " + p.cur.Str
					p.advance()
				}
			}
			// VALUE 'value'
			if p.isIdentLikeStr("VALUE") {
				p.advance()
				if p.cur.Type == tokSCONST {
					opt.Value += " VALUE " + p.cur.Str
					p.advance()
				}
			}
			items = append(items, opt)
		} else {
			break
		}
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseAttrDimSourceClause parses a source_clause in an attribute dimension USING clause.
func (p *Parser) parseAttrDimSourceClause() (*nodes.AttrDimSourceClause, error) {
	src := &nodes.AttrDimSourceClause{Loc: nodes.Loc{Start: p.pos()}}

	// JOIN PATH join_path_name ON join_condition
	if p.cur.Type == kwJOIN {
		p.advance() // consume JOIN
		if p.isIdentLikeStr("PATH") {
			p.advance() // consume PATH
		}
		src.IsJoinPath = true
		var parseErr288 error
		src.JoinPathName, parseErr288 = p.parseIdentifier()
		if parseErr288 != nil {
			return nil, parseErr288
		}
		if p.cur.Type == kwON {
			p.advance()
		}
		src.JoinCondition = &nodes.List{}
		for {
			elem, parseErr289 := p.parseAttrDimJoinCondElem()
			if parseErr289 != nil {
				return nil, parseErr289
			}
			if elem != nil {
				src.JoinCondition.Items = append(src.JoinCondition.Items, elem)
			}
			if p.cur.Type != kwAND {
				break
			}
			p.advance() // consume AND
		}
		src.Loc.End = p.prev.End
		return src, nil
	}

	// [ REMOTE ]
	if p.isIdentLikeStr("REMOTE") {
		src.Remote = true
		p.advance()
	}
	var parseErr290 error

	// [ schema. ] table/view
	src.Name, parseErr290 = p.parseObjectName()
	if parseErr290 !=

		// [ [ AS ] alias ]
		nil {
		return nil, parseErr290
	}

	if p.cur.Type == kwAS {
		p.advance()
		var parseErr291 error
		src.Alias, parseErr291 = p.parseIdentifier()
		if parseErr291 != nil {
			return nil, parseErr291
		}
	} else if p.isIdentLike() && p.cur.Str != "ATTRIBUTES" && p.cur.Str != "JOIN" &&
		p.cur.Str != "DIMENSION" && p.cur.Str != "SHARING" && p.cur.Str != "CAPTION" &&
		p.cur.Str != "DESCRIPTION" && p.cur.Str != "CLASSIFICATION" &&
		p.cur.Type != ',' && p.cur.Type != ';' && p.cur.Type != tokEOF {
		var parseErr292 error
		src.Alias, parseErr292 = p.parseIdentifier()
		if parseErr292 != nil {
			return nil, parseErr292
		}
	}

	src.Loc.End = p.prev.End
	return src, nil
}

// parseAttrDimJoinCondElem parses a join condition element: [table.]col = [table.]col
func (p *Parser) parseAttrDimJoinCondElem() (*nodes.AttrDimJoinCondElem, error) {
	elem := &nodes.AttrDimJoinCondElem{Loc: nodes.Loc{Start: p.pos()}}
	// Left side
	name1, parseErr293 := p.parseIdentifier()
	if parseErr293 != nil {
		return nil, parseErr293
	}
	if p.cur.Type == '.' {
		p.advance()
		elem.LeftTable = name1
		var parseErr294 error
		elem.LeftCol, parseErr294 = p.parseIdentifier()
		if parseErr294 != nil {
			return nil, parseErr294

			// =
		}
	} else {
		elem.LeftCol = name1
	}

	if p.cur.Type == '=' {
		p.advance()
	}
	// Right side
	name2, parseErr295 := p.parseIdentifier()
	if parseErr295 != nil {
		return nil, parseErr295
	}
	if p.cur.Type == '.' {
		p.advance()
		elem.RightTable = name2
		var parseErr296 error
		elem.RightCol, parseErr296 = p.parseIdentifier()
		if parseErr296 != nil {
			return nil, parseErr296
		}
	} else {
		elem.RightCol = name2
	}
	elem.Loc.End = p.prev.End
	return elem, nil
}

// parseAttrDimAttributeClause parses an attribute in the ATTRIBUTES clause.
func (p *Parser) parseAttrDimAttributeClause() (*nodes.AttrDimAttribute, error) {
	attr := &nodes.AttrDimAttribute{Loc: nodes.Loc{Start: p.pos()}}
	var parseErr297 error
	attr.Column, parseErr297 = p.parseIdentifier()
	if parseErr297 != nil {
		return nil, parseErr297
	}
	if p.cur.Type == kwAS {
		p.advance()
		var parseErr298 error
		attr.Alias, parseErr298 = p.parseIdentifier()
		if parseErr298 != nil {
			return nil, parseErr298
		}
	}
	var parseErr299 error
	attr.Classifications, parseErr299 = p.parseClassificationClauses()
	if parseErr299 != nil {
		return nil, parseErr299
	}
	attr.Loc.End = p.prev.End
	return attr, nil
}

// parseAttrDimLevelClause parses a LEVEL clause in CREATE ATTRIBUTE DIMENSION.
func (p *Parser) parseAttrDimLevelClause() (*nodes.AttrDimLevel, error) {
	lvl := &nodes.AttrDimLevel{Loc: nodes.Loc{Start: p.pos()}}
	p.advance()
	var // consume LEVEL
	parseErr300 error

	lvl.Name, parseErr300 = p.parseIdentifier()
	if parseErr300 !=

		// LEVEL TYPE { STANDARD | YEARS | ... }
		nil {
		return nil, parseErr300
	}

	if p.cur.Type == kwLEVEL {
		next := p.peekNext()
		if next.Type == kwTYPE {
			p.advance() // consume LEVEL
			p.advance() // consume TYPE
			if p.isIdentLike() {
				lvl.LevelType = p.cur.Str
				p.advance()
			}
		}
	}
	var parseErr301 error

	// classification_clause(s)
	lvl.Classifications, parseErr301 = p.parseClassificationClauses()
	if parseErr301 !=

		// KEY clause
		nil {
		return nil, parseErr301
	}

	if p.cur.Type == kwKEY {
		p.advance()
		lvl.KeyAttrs = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseValue49, parseErr50 := p.parseIdentifier()
				if parseErr50 != nil {
					return nil, parseErr50
				}
				lvl.KeyAttrs.Items = append(lvl.KeyAttrs.Items, &nodes.String{Str: parseValue49})
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			parseValue51, parseErr52 := p.parseIdentifier()
			if parseErr52 !=
				// NOT NULL | SKIP WHEN NULL
				nil {
				return nil, parseErr52
			}
			lvl.KeyAttrs.Items = append(lvl.KeyAttrs.Items, &nodes.String{Str: parseValue51})

			if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
				lvl.KeyNotNull = true
				p.advance()
				p.advance()
			} else if p.cur.Type == kwSKIP {
				p.advance() // consume SKIP
				if p.cur.Type == kwWHEN {
					p.advance() // consume WHEN
				}
				if p.cur.Type == kwNULL {
					p.advance() // consume NULL
				}
				lvl.KeySkipWhenNull = true
			}
		}
	}

	// ALTERNATE KEY
	if p.isIdentLikeStr("ALTERNATE") {
		p.advance() // consume ALTERNATE
		if p.cur.Type == kwKEY {
			p.advance() // consume KEY
		}
		lvl.AltKeyAttrs = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseValue53, parseErr54 := p.parseIdentifier()
				if parseErr54 != nil {
					return nil, parseErr54
				}
				lvl.AltKeyAttrs.Items = append(lvl.AltKeyAttrs.Items, &nodes.String{Str: parseValue53})
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			parseValue55, parseErr56 := p.parseIdentifier()
			if parseErr56 != nil {
				return nil, parseErr56

				// MEMBER NAME / CAPTION / DESCRIPTION
			}
			lvl.AltKeyAttrs.Items = append(lvl.AltKeyAttrs.Items, &nodes.String{Str: parseValue55})
		}
	}
	var parseErr302 error

	lvl.MemberName, lvl.MemberCaption, lvl.MemberDesc, parseErr302 = p.parseMemberExprs()
	if parseErr302 !=

		// ORDER BY
		nil {
		return nil, parseErr302
	}

	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		lvl.OrderByAttrs = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				item := &nodes.AttrDimOrderByItem{Loc: nodes.Loc{Start: p.pos()}}
				var parseErr303 error
				item.Attribute, parseErr303 = p.parseIdentifier()
				if parseErr303 != nil {
					return nil, parseErr303
				}
				if p.cur.Type == kwASC {
					item.Direction = "ASC"
					p.advance()
				} else if p.cur.Type == kwDESC {
					item.Direction = "DESC"
					p.advance()
				}
				item.Loc.End = p.prev.End
				lvl.OrderByAttrs.Items = append(lvl.OrderByAttrs.Items, item)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			item := &nodes.AttrDimOrderByItem{Loc: nodes.Loc{Start: p.pos()}}
			var parseErr304 error
			item.Attribute, parseErr304 = p.parseIdentifier()
			if parseErr304 != nil {
				return nil, parseErr304
			}
			if p.cur.Type == kwASC {
				item.Direction = "ASC"
				p.advance()
			} else if p.cur.Type == kwDESC {
				item.Direction = "DESC"
				p.advance()
			}
			item.Loc.End = p.prev.End
			lvl.OrderByAttrs.Items = append(lvl.OrderByAttrs.Items, item)
		}
	}

	// DETERMINES ( attribute [, attribute ]... )
	if p.isIdentLikeStr("DETERMINES") {
		p.advance()
		lvl.Determines = &nodes.List{}
		if p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseValue57, parseErr58 := p.parseIdentifier()
				if parseErr58 != nil {
					return nil, parseErr58
				}
				lvl.Determines.Items = append(lvl.Determines.Items, &nodes.String{Str: parseValue57})
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// parseMemberExprs parses MEMBER NAME/CAPTION/DESCRIPTION expression triples.
func (p *Parser) parseMemberExprs() (name, caption, desc nodes.ExprNode, parseErr error) {
	for p.isIdentLikeStr("MEMBER") {
		p.advance() // consume MEMBER
		switch {
		case p.cur.Type == kwNAME || p.isIdentLikeStr("NAME"):
			p.advance()
			var parseErr305 error
			name, parseErr305 = p.parseExpr()
			if parseErr305 != nil {
				return nil, nil, nil, parseErr305
			}
		case p.isIdentLikeStr("CAPTION"):
			p.advance()
			var parseErr306 error
			caption, parseErr306 = p.parseExpr()
			if parseErr306 != nil {
				return nil, nil, nil, parseErr306
			}
		case p.isIdentLikeStr("DESCRIPTION"):
			p.advance()
			var parseErr307 error
			desc, parseErr307 = p.parseExpr()
			if parseErr307 != nil {
				return nil, nil, nil, parseErr307

				// ---------------------------------------------------------------------------
				// ALTER ATTRIBUTE DIMENSION
				// ---------------------------------------------------------------------------
			}
		default:
			return
		}
	}
	return
}

// parseAlterAttributeDimensionStmt parses an ALTER ATTRIBUTE DIMENSION statement.
//
// BNF: oracle/parser/bnf/ALTER-ATTRIBUTE-DIMENSION.bnf
//
//	ALTER ATTRIBUTE DIMENSION [ IF EXISTS ] [ schema . ] attr_dim_name
//	    {
//	        RENAME TO new_attr_dim_name
//	      | COMPILE
//	    }
func (p *Parser) parseAlterAttributeDimensionStmt(start int) (*nodes.AlterAttributeDimensionStmt, error) {
	stmt := &nodes.AlterAttributeDimensionStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var parseErr308 error

	stmt.Name, parseErr308 = p.parseObjectName()
	if parseErr308 != nil {
		return nil, parseErr308
	}

	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwRENAME || p.isIdentLikeStr("RENAME"):
		stmt.Action = "RENAME"
		p.advance() // consume RENAME
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		var parseErr309 error
		stmt.NewName, parseErr309 = p.parseObjectName()
		if parseErr309 != nil {
			return nil, parseErr309
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE HIERARCHY
// ---------------------------------------------------------------------------

// parseCreateHierarchyStmt parses a CREATE HIERARCHY statement.
//
// BNF: oracle/parser/bnf/CREATE-HIERARCHY.bnf
//
//	CREATE [ OR REPLACE ] [ { FORCE | NOFORCE } ] HIERARCHY
//	    [ IF NOT EXISTS ] [ schema. ] hierarchy
//	    [ SHARING = { METADATA | NONE } ]
//	    [ classification_clause ]...
//	    hier_using_clause
//	    ( level_hier_clause )
//	    [ hier_attrs_clause ] ;
//
//	classification_clause:
//	    { CAPTION 'caption'
//	    | DESCRIPTION 'description'
//	    | CLASSIFICATION classification_name VALUE 'classification_value'
//	        [ LANGUAGE language ] }
//
//	hier_using_clause:
//	    USING [ schema. ] attr_dimension
//
//	level_hier_clause:
//	    level_name [ classification_clause ]...
//	        [ CHILD OF level_hier_clause ]
//
//	hier_attrs_clause:
//	    HIERARCHICAL ATTRIBUTES ( hier_attr_clause [, hier_attr_clause ]... )
//
//	hier_attr_clause:
//	    hier_attr_name [ classification_clause ]...
//
//	hier_attr_name:
//	    { HIER_ORDER | DEPTH | IS_LEAF | IS_ROOT
//	    | MEMBER_NAME | MEMBER_UNIQUE_NAME | MEMBER_CAPTION | MEMBER_DESCRIPTION
//	    | PARENT_LEVEL_NAME | PARENT_UNIQUE_NAME }
func (p *Parser) parseCreateHierarchyStmt(start int, orReplace, force, noForce bool) (*nodes.CreateHierarchyStmt, error) {
	stmt := &nodes.CreateHierarchyStmt{
		OrReplace: orReplace,
		Force:     force,
		NoForce:   noForce,
		Loc:       nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance()
		p.advance()
		if p.cur.Type == kwEXISTS {
			p.advance()
		}
		stmt.IfNotExists = true
	}
	var parseErr310 error

	stmt.Name, parseErr310 = p.parseObjectName()
	if parseErr310 !=

		// SHARING = { METADATA | NONE }
		nil {
		return nil, parseErr310
	}

	if p.isIdentLikeStr("SHARING") {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}
	var parseErr311 error

	// classification_clause(s)
	stmt.Classifications, parseErr311 = p.parseClassificationClauses()
	if parseErr311 !=

		// USING [ schema. ] attr_dimension
		nil {
		return nil, parseErr311
	}

	if p.cur.Type == kwUSING {
		p.advance()
		var parseErr312 error
		stmt.UsingAttrDim, parseErr312 = p.parseObjectName()
		if parseErr312 !=

			// ( level_hier_clause )
			nil {
			return nil, parseErr312
		}
	}

	if p.cur.Type == '(' {
		p.advance()
		var parseErr313 error
		stmt.LevelHier, parseErr313 = p.parseHierLevelClause()
		if parseErr313 != nil {
			return nil, parseErr313
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// HIERARCHICAL ATTRIBUTES ( ... )
	if p.isIdentLikeStr("HIERARCHICAL") {
		p.advance() // consume HIERARCHICAL
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance() // consume ATTRIBUTES
		}
		if p.cur.Type == '(' {
			p.advance()
			stmt.HierAttrs = &nodes.List{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				ha := &nodes.HierAttr{Loc: nodes.Loc{Start: p.pos()}}
				var parseErr314 error
				ha.Name, parseErr314 = p.parseIdentifier()
				if parseErr314 != nil {
					return nil, parseErr314
				}
				var parseErr315 error
				ha.Classifications, parseErr315 = p.parseClassificationClauses()
				if parseErr315 != nil {
					return nil, parseErr315
				}
				ha.Loc.End = p.prev.End
				stmt.HierAttrs.Items = append(stmt.HierAttrs.Items, ha)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseHierLevelClause parses a level_hier_clause recursively.
func (p *Parser) parseHierLevelClause() (*nodes.HierLevelClause, error) {
	lvl := &nodes.HierLevelClause{Loc: nodes.Loc{Start: p.pos()}}
	var parseErr316 error
	lvl.Name, parseErr316 = p.parseIdentifier()
	if parseErr316 != nil {
		return nil, parseErr316
	}
	var parseErr317 error
	lvl.Classifications, parseErr317 = p.parseClassificationClauses()
	if parseErr317 !=

		// CHILD OF level_hier_clause
		nil {
		return nil, parseErr317
	}

	if p.isIdentLikeStr("CHILD") {
		p.advance() // consume CHILD
		if p.cur.Type == kwOF {
			p.advance() // consume OF
		}
		var parseErr318 error
		lvl.ChildOf, parseErr318 = p.parseHierLevelClause()
		if parseErr318 != nil {
			return nil, parseErr318
		}
	}

	lvl.Loc.End = p.prev.End
	return lvl, nil
}

// ---------------------------------------------------------------------------
// ALTER HIERARCHY
// ---------------------------------------------------------------------------

// parseAlterHierarchyStmt parses an ALTER HIERARCHY statement.
//
// BNF: oracle/parser/bnf/ALTER-HIERARCHY.bnf
//
//	ALTER HIERARCHY [ IF EXISTS ] [ schema. ] hierarchy_name
//	    { RENAME TO new_hier_name
//	    | COMPILE
//	    } ;
func (p *Parser) parseAlterHierarchyStmt(start int) (*nodes.AlterHierarchyStmt, error) {
	stmt := &nodes.AlterHierarchyStmt{
		Loc: nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var parseErr319 error

	stmt.Name, parseErr319 = p.parseObjectName()
	if parseErr319 != nil {
		return nil, parseErr319
	}

	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwRENAME || p.isIdentLikeStr("RENAME"):
		stmt.Action = "RENAME"
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
		}
		var parseErr320 error
		stmt.NewName, parseErr320 = p.parseObjectName()
		if parseErr320 != nil {
			return nil, parseErr320
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE DOMAIN
// ---------------------------------------------------------------------------

// parseCreateDomainStmt parses a CREATE DOMAIN statement.
//
// BNF: oracle/parser/bnf/create-domain.bnf
//
//	CREATE [ OR REPLACE ] [ USECASE ] DOMAIN [ IF NOT EXISTS ] [ schema. ] domain_name
//	    AS { create_single_column_domain
//	       | create_multi_column_domain
//	       | create_flexible_domain } ;
//
//	create_single_column_domain:
//	    datatype [ STRICT ]
//	    [ DEFAULT [ ON NULL ] default_expression ]
//	    [ { constraint_clause }... ]
//	    [ COLLATE collation_name ]
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//	    [ annotations_clause ]
//
//	create_single_column_domain:  -- ENUM variant
//	    ENUM ( enum_list ) [ STRICT ]
//	    [ DEFAULT [ ON NULL ] default_expression ]
//	    [ { constraint_clause }... ]
//	    [ COLLATE collation_name ]
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//	    [ annotations_clause ]
//
//	enum_list:
//	    enum_item_list [, enum_item_list ]...
//
//	enum_item_list:
//	    name [ = enum_alias_list ] [ = value ]
//
//	enum_alias_list:
//	    alias [ = alias ]...
//
//	column_properties_clause:
//	    [ DEFAULT [ ON NULL ] default_expression ]
//	    [ { constraint_clause }... ]
//	    [ COLLATE collation_name ]
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//
//	create_multi_column_domain:
//	    ( column_name AS datatype [ annotations_clause ]
//	      [, column_name AS datatype [ annotations_clause ] ]... )
//	    [ { constraint_clause }... ]
//	    [ COLLATE collation_name ]
//	    [ DISPLAY display_expression ]
//	    [ ORDER order_expression ]
//	    [ annotations_clause ]
//
//	create_flexible_domain:
//	    FLEXIBLE DOMAIN [ schema. ] domain_name
//	        ( column_name [, column_name ]... )
//	    CHOOSE DOMAIN USING ( domain_discriminant_column datatype
//	        [, domain_discriminant_column datatype ]... )
//	    FROM { DECODE ( expr , comparison_expr , result_expr
//	                    [, comparison_expr , result_expr ]... )
//	         | CASE { WHEN condition THEN result_expr }... END }
//
//	result_expr:
//	    [ schema. ] domain_name ( column_name [, column_name ]... )
//
//	default_clause:
//	    DEFAULT [ ON NULL ] default_expression
//
//	constraint_clause:
//	    [ CONSTRAINT constraint_name ]
//	    { NOT NULL | NULL | CHECK ( condition ) }
//	    [ constraint_state ]
//
//	annotations_clause:
//	    ANNOTATIONS ( annotation [, annotation ]... )
func (p *Parser) parseCreateDomainStmt(start int, orReplace, usecase bool) (*nodes.CreateDomainStmt, error) {
	stmt := &nodes.CreateDomainStmt{
		OrReplace: orReplace,
		Usecase:   usecase,
		Loc:       nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		p.advance()
		p.advance()
		if p.cur.Type == kwEXISTS {
			p.advance()
		}
		stmt.IfNotExists = true
	}
	var parseErr321 error

	stmt.Name, parseErr321 = p.parseObjectName()
	if parseErr321 !=

		// AS
		nil {
		return nil, parseErr321
	}

	if p.cur.Type == kwAS {
		p.advance()
	}

	// Determine variant: flexible (FLEXIBLE), enum (ENUM), multi-column '(', or single-column (datatype)
	switch {
	case p.isIdentLikeStr("FLEXIBLE"):
		stmt.DomainType = "FLEXIBLE"
		p.advance() // consume FLEXIBLE
		// DOMAIN [schema.]domain_name (column_name, ...)
		if p.isIdentLike() && p.cur.Str == "DOMAIN" {
			p.advance()
		}
		var parseErr322 error
		stmt.FlexDomainName, parseErr322 = p.parseObjectName()
		if parseErr322 != nil {
			return nil, parseErr322
		}
		if p.cur.Type == '(' {
			p.advance()
			stmt.FlexColumns = &nodes.List{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				parseValue59, parseErr60 := p.parseIdentifier()
				if parseErr60 != nil {
					return nil, parseErr60
				}
				stmt.FlexColumns.Items = append(stmt.FlexColumns.Items, &nodes.String{Str: parseValue59})
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// CHOOSE DOMAIN USING ( ... )
		if p.isIdentLikeStr("CHOOSE") {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DOMAIN" {
				p.advance()
			}
			if p.cur.Type == kwUSING {
				p.advance()
			}
			if p.cur.Type == '(' {
				p.advance()
				stmt.ChooseUsing = &nodes.List{}
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					col := &nodes.DomainColumn{Loc: nodes.Loc{Start: p.pos()}}
					var parseErr323 error
					col.Name, parseErr323 = p.parseIdentifier()
					if parseErr323 != nil {
						return nil, parseErr323
					}
					var parseErr324 error
					col.DataType, parseErr324 = p.parseTypeName()
					if parseErr324 != nil {
						return nil, parseErr324
					}
					col.Loc.End = p.prev.End
					stmt.ChooseUsing.Items = append(stmt.ChooseUsing.Items, col)
					if p.cur.Type == ',' {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
		// FROM { DECODE(...) | CASE ... END }
		if p.cur.Type == kwFROM {
			p.advance()
			var parseErr325 error
			stmt.ChooseExpr, parseErr325 = p.parseExpr()
			if parseErr325 != nil {
				return nil, parseErr325
			}
		}

	case p.isIdentLikeStr("ENUM"):
		stmt.DomainType = "ENUM"
		p.advance() // consume ENUM
		// ( enum_list )
		if p.cur.Type == '(' {
			p.advance()
			stmt.EnumItems = &nodes.List{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				item := &nodes.DomainEnumItem{Loc: nodes.Loc{Start: p.pos()}}
				var parseErr326 error
				item.Name, parseErr326 = p.parseIdentifier()
				if parseErr326 !=
					// [ = alias [= alias]... ] [ = value ]
					nil {
					return nil, parseErr326
				}

				for p.cur.Type == '=' {
					p.advance()
					if p.cur.Type == tokICONST || p.cur.Type == tokFCONST || p.cur.Type == tokSCONST {
						var parseErr327 error
						item.Value, parseErr327 = p.parseExpr()
						if parseErr327 != nil {
							return nil, parseErr327
						}
					} else if p.isIdentLike() || p.cur.Type == tokIDENT {
						parseValue61, parseErr62 := p.parseIdentifier()
						if parseErr62 != nil {
							return nil, parseErr62
						}
						item.Aliases = append(item.Aliases, parseValue61)
					} else {
						var parseErr328 error
						item.Value, parseErr328 = p.parseExpr()
						if parseErr328 != nil {
							return nil, parseErr328
						}
					}
				}
				item.Loc.End = p.prev.End
				stmt.EnumItems.Items = append(stmt.EnumItems.Items, item)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// STRICT
		if p.isIdentLikeStr("STRICT") {
			stmt.Strict = true
			p.advance()
		}
		parseErr329 := p.parseDomainProperties(stmt)
		if parseErr329 != nil {
			return nil, parseErr329

			// Multi-column domain: ( column_name AS datatype [, ...] )
		}

	case p.cur.Type == '(':

		stmt.DomainType = "MULTI"
		p.advance() // consume (
		stmt.Columns = &nodes.List{}
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			col := &nodes.DomainColumn{Loc: nodes.Loc{Start: p.pos()}}
			var parseErr330 error
			col.Name, parseErr330 = p.parseIdentifier()
			if parseErr330 != nil {
				return nil, parseErr330
			}
			if p.cur.Type == kwAS {
				p.advance()
			}
			var parseErr331 error
			col.DataType, parseErr331 = p.parseTypeName()
			if parseErr331 !=
				// annotations_clause per column
				nil {
				return nil, parseErr331
			}

			if p.isIdentLikeStr("ANNOTATIONS") {
				var parseErr332 error
				col.Annotations, parseErr332 = p.parseDomainAnnotations()
				if parseErr332 != nil {
					return nil, parseErr332
				}
			}
			col.Loc.End = p.prev.End
			stmt.Columns.Items = append(stmt.Columns.Items, col)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
		parseErr333 := p.parseDomainProperties(stmt)
		if parseErr333 !=

			// Single-column domain: datatype [STRICT] [properties...]
			nil {
			return nil, parseErr333
		}

	default:

		stmt.DomainType = "SINGLE"
		var parseErr334 error
		stmt.DataType, parseErr334 = p.parseTypeName()
		if parseErr334 != nil {
			return nil, parseErr334
		}
		if p.isIdentLikeStr("STRICT") {
			stmt.Strict = true
			p.advance()
		}
		parseErr335 := p.parseDomainProperties(stmt)
		if parseErr335 != nil {
			return nil, parseErr335
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDomainProperties parses the common properties after a domain type definition.
func (p *Parser) parseDomainProperties(stmt *nodes.CreateDomainStmt) error {
	// DEFAULT [ON NULL] expr
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if p.cur.Type == kwON && p.peekNext().Type == kwNULL {
			stmt.DefaultOnNull = true
			p.advance()
			p.advance()
		}
		var parseErr336 error
		stmt.Default, parseErr336 = p.parseExpr()
		if parseErr336 !=

			// constraint_clause(s)
			nil {
			return parseErr336
		}
	}
	var parseErr337 error

	stmt.Constraints, parseErr337 = p.parseDomainConstraints()
	if parseErr337 !=

		// COLLATE collation_name
		nil {
		return parseErr337
	}

	if p.isIdentLikeStr("COLLATE") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokSCONST {
			stmt.Collation = p.cur.Str
			p.advance()
		}
	}

	// DISPLAY display_expression
	if p.isIdentLikeStr("DISPLAY") {
		p.advance()
		var parseErr338 error
		stmt.Display, parseErr338 = p.parseExpr()
		if parseErr338 !=

			// ORDER order_expression
			nil {
			return parseErr338
		}
	}

	if p.cur.Type == kwORDER {
		p.advance()
		var parseErr339 error
		stmt.Order, parseErr339 = p.parseExpr()
		if parseErr339 !=

			// ANNOTATIONS ( ... )
			nil {
			return parseErr339
		}
	}

	if p.isIdentLikeStr("ANNOTATIONS") {
		var parseErr340 error
		stmt.Annotations, parseErr340 = p.parseDomainAnnotations()
		if parseErr340 !=

			// parseDomainConstraints parses zero or more constraint_clause items.
			nil {
			return parseErr340
		}
	}
	return nil
}

func (p *Parser) parseDomainConstraints() (*nodes.List, error) {
	var items []nodes.Node
	for {
		var c *nodes.DomainConstraint
		if p.cur.Type == kwCONSTRAINT {
			c = &nodes.DomainConstraint{Loc: nodes.Loc{Start: p.pos()}}
			p.advance()
			var parseErr341 error
			c.Name, parseErr341 = p.parseIdentifier()
			if parseErr341 != nil {
				return nil, parseErr341
			}
		}
		if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
			if c == nil {
				c = &nodes.DomainConstraint{Loc: nodes.Loc{Start: p.pos()}}
			}
			c.Type = "NOT_NULL"
			p.advance()
			p.advance()
		} else if p.cur.Type == kwNULL {
			if c == nil {
				c = &nodes.DomainConstraint{Loc: nodes.Loc{Start: p.pos()}}
			}
			c.Type = "NULL"
			p.advance()
		} else if p.cur.Type == kwCHECK {
			if c == nil {
				c = &nodes.DomainConstraint{Loc: nodes.Loc{Start: p.pos()}}
			}
			c.Type = "CHECK"
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				var parseErr342 error
				c.CheckExpr, parseErr342 = p.parseExpr()
				if parseErr342 != nil {
					return nil, parseErr342
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		} else if c != nil {
			// CONSTRAINT name without a recognized type — shouldn't happen, keep going
			c.Type = "UNKNOWN"
		} else {
			break
		}
		c.Loc.End = p.prev.End
		items = append(items, c)
	}
	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseDomainAnnotations parses ANNOTATIONS ( ... ) for domains.
func (p *Parser) parseDomainAnnotations() (*nodes.List, error) {
	p.advance() // consume ANNOTATIONS
	result := &nodes.List{}
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			key, parseErr343 := p.parseIdentifier()
			if parseErr343 != nil {
				return nil, parseErr343
			}
			var val string
			if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			}
			result.Items = append(result.Items, &nodes.DDLOption{Key: key, Value: val})
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}
	return result, nil
}

// ---------------------------------------------------------------------------
// ALTER DOMAIN
// ---------------------------------------------------------------------------

// parseAlterDomainStmt parses an ALTER DOMAIN statement.
//
// BNF: oracle/parser/bnf/alter-domain.bnf
//
//	ALTER [ USECASE ] DOMAIN [ IF EXISTS ] [ schema. ] domain_name
//	    { ADD DISPLAY display_expression
//	    | MODIFY DISPLAY display_expression
//	    | DROP DISPLAY
//	    | ADD ORDER order_expression
//	    | MODIFY ORDER order_expression
//	    | DROP ORDER
//	    | annotations_clause
//	    } ;
func (p *Parser) parseAlterDomainStmt(start int, usecase bool) (*nodes.AlterDomainStmt, error) {
	stmt := &nodes.AlterDomainStmt{
		Usecase: usecase,
		Loc:     nodes.Loc{Start: start},
	}

	// IF EXISTS
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var parseErr344 error

	stmt.Name, parseErr344 = p.parseObjectName()
	if parseErr344 !=

		// Action
		nil {
		return nil, parseErr344
	}

	switch {
	case p.isIdentLikeStr("ADD"):
		p.advance()
		if p.isIdentLikeStr("DISPLAY") {
			stmt.Action = "ADD_DISPLAY"
			p.advance()
			var parseErr345 error
			stmt.Display, parseErr345 = p.parseExpr()
			if parseErr345 != nil {
				return nil, parseErr345
			}
		} else if p.cur.Type == kwORDER {
			stmt.Action = "ADD_ORDER"
			p.advance()
			var parseErr346 error
			stmt.Order, parseErr346 = p.parseExpr()
			if parseErr346 != nil {
				return nil, parseErr346
			}
		}
	case p.isIdentLikeStr("MODIFY"):
		p.advance()
		if p.isIdentLikeStr("DISPLAY") {
			stmt.Action = "MODIFY_DISPLAY"
			p.advance()
			var parseErr347 error
			stmt.Display, parseErr347 = p.parseExpr()
			if parseErr347 != nil {
				return nil, parseErr347
			}
		} else if p.cur.Type == kwORDER {
			stmt.Action = "MODIFY_ORDER"
			p.advance()
			var parseErr348 error
			stmt.Order, parseErr348 = p.parseExpr()
			if parseErr348 != nil {
				return nil, parseErr348
			}
		}
	case p.isIdentLikeStr("DROP"):
		p.advance()
		if p.isIdentLikeStr("DISPLAY") {
			stmt.Action = "DROP_DISPLAY"
			p.advance()
		} else if p.cur.Type == kwORDER {
			stmt.Action = "DROP_ORDER"
			p.advance()
		}
	case p.cur.Type == kwRENAME:
		stmt.Action = "RENAME"
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
		}
		parseDiscard350, parseErr349 := p.parseObjectName()
		_ = parseDiscard350
		if parseErr349 != nil {
			return nil, parseErr349
		}
	case p.isIdentLikeStr("ANNOTATIONS"):
		stmt.Action = "ANNOTATIONS"
		var parseErr351 error
		stmt.Annotations, parseErr351 = p.parseDomainAnnotations()
		if parseErr351 != nil {
			return nil, parseErr351
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE PROPERTY GRAPH
// ---------------------------------------------------------------------------

// parseCreatePropertyGraphStmt parses a CREATE PROPERTY GRAPH statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] PROPERTY GRAPH has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE ] PROPERTY GRAPH [ IF NOT EXISTS ]
//	    [ schema. ] graph_name
//	    vertex_tables_clause
//	    [ edge_tables_clause ]
//	    [ graph_options ] ;
func (p *Parser) parseCreatePropertyGraphStmt(start int, orReplace, ifNotExists bool) (*nodes.CreatePropertyGraphStmt, error) {
	stmt := &nodes.CreatePropertyGraphStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}

	// IF NOT EXISTS (may also be parsed here if not consumed by caller)
	if !ifNotExists && p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
				stmt.IfNotExists = true
			}
		}
	}
	var parseErr352 error

	// Graph name
	stmt.Name, parseErr352 = p.parseObjectName()
	if parseErr352 !=

		// VERTEX TABLES ( ... )
		nil {
		return nil, parseErr352
	}

	if p.isIdentLikeStr("VERTEX") {
		p.advance() // consume VERTEX
		if p.isIdentLikeStr("TABLES") {
			p.advance() // consume TABLES
		}
		if p.cur.Type == '(' {
			p.advance() // consume (
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				vtd, parseErr353 := p.parseGraphTableDef()
				if parseErr353 != nil {
					return nil, parseErr353
				}
				stmt.VertexTables = append(stmt.VertexTables, vtd)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// EDGE TABLES ( ... )
	if p.isIdentLikeStr("EDGE") {
		p.advance() // consume EDGE
		if p.isIdentLikeStr("TABLES") {
			p.advance() // consume TABLES
		}
		if p.cur.Type == '(' {
			p.advance() // consume (
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				etd, parseErr354 := p.parseGraphEdgeDef()
				if parseErr354 != nil {
					return nil, parseErr354
				}
				stmt.EdgeTables = append(stmt.EdgeTables, etd)
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// OPTIONS ( ... )
	if p.isIdentLikeStr("OPTIONS") {
		p.advance() // consume OPTIONS
		if p.cur.Type == '(' {
			p.advance() // consume (
			opts := &nodes.GraphOptions{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLikeStr("ENFORCED") {
					opts.Mode = "ENFORCED"
					p.advance()
					if p.isIdentLikeStr("MODE") {
						p.advance()
					}
				} else if p.isIdentLikeStr("TRUSTED") {
					opts.Mode = "TRUSTED"
					p.advance()
					if p.isIdentLikeStr("MODE") {
						p.advance()
					}
				} else if p.isIdentLikeStr("ALLOW") {
					p.advance()
					opts.MixedPropTypes = "ALLOW"
					// skip MIXED PROPERTY TYPES
					for p.isIdentLike() && p.cur.Type != ')' {
						p.advance()
					}
				} else if p.isIdentLikeStr("DISALLOW") {
					p.advance()
					opts.MixedPropTypes = "DISALLOW"
					// skip MIXED PROPERTY TYPES
					for p.isIdentLike() && p.cur.Type != ')' {
						p.advance()
					}
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
			stmt.Options = opts
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseGraphTableDef parses a vertex table definition.
func (p *Parser) parseGraphTableDef() (*nodes.GraphTableDef, error) {
	def := &nodes.GraphTableDef{
		Loc: nodes.Loc{Start: p.pos()},
	}
	var parseErr355 error

	// [schema.]table_name
	def.Name, parseErr355 = p.parseObjectName()
	if parseErr355 !=

		// AS graph_element_name
		nil {
		return nil, parseErr355
	}

	if p.cur.Type == kwAS {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			def.Alias = p.cur.Str
			p.advance()
		}
	}

	// KEY ( col, ... )
	if p.cur.Type == kwKEY {
		p.advance()
		var parseErr356 error
		def.KeyColumns, parseErr356 = p.parseParenIdentList()
		if parseErr356 !=

			// Label and properties
			nil {
			return nil, parseErr356
		}
	}
	parseErr357 := p.parseGraphLabelAndProperties(def)
	if parseErr357 != nil {
		return nil, parseErr357
	}

	def.Loc.End = p.prev.End
	return def, nil
}

// parseGraphEdgeDef parses an edge table definition.
func (p *Parser) parseGraphEdgeDef() (*nodes.GraphEdgeDef, error) {
	def := &nodes.GraphEdgeDef{
		Loc: nodes.Loc{Start: p.pos()},
	}
	var parseErr358 error

	// [schema.]table_name
	def.Name, parseErr358 = p.parseObjectName()
	if parseErr358 !=

		// AS graph_element_name
		nil {
		return nil, parseErr358
	}

	if p.cur.Type == kwAS {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			def.Alias = p.cur.Str
			p.advance()
		}
	}

	// KEY ( col, ... )
	if p.cur.Type == kwKEY {
		p.advance()
		var parseErr359 error
		def.KeyColumns, parseErr359 = p.parseParenIdentList()
		if parseErr359 !=

			// SOURCE [ KEY ( col, ... ) REFERENCES ] vertex_table_ref
			nil {
			return nil, parseErr359
		}
	}

	if p.isIdentLikeStr("SOURCE") {
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
			var parseErr360 error
			def.SourceKeyColumns, parseErr360 = p.parseParenIdentList()
			if parseErr360 != nil {
				return nil, parseErr360
			}
			if p.cur.Type == kwREFERENCES {
				p.advance()
			}
		}
		// vertex_table_reference: name ( col, ... )
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			def.SourceRef = p.cur.Str
			p.advance()
			if p.cur.Type == '(' {
				var parseErr361 error
				def.SourceRefColumns, parseErr361 = p.parseParenIdentList()
				if parseErr361 != nil {

					// DESTINATION [ KEY ( col, ... ) REFERENCES ] vertex_table_ref
					return nil, parseErr361
				}
			}
		}
	}

	if p.isIdentLikeStr("DESTINATION") {
		p.advance()
		if p.cur.Type == kwKEY {
			p.advance()
			var parseErr362 error
			def.DestKeyColumns, parseErr362 = p.parseParenIdentList()
			if parseErr362 != nil {
				return nil, parseErr362
			}
			if p.cur.Type == kwREFERENCES {
				p.advance()
			}
		}
		// vertex_table_reference: name ( col, ... )
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			def.DestRef = p.cur.Str
			p.advance()
			if p.cur.Type == '(' {
				var parseErr363 error
				def.DestRefColumns, parseErr363 = p.parseParenIdentList()
				if parseErr363 != nil {

					// Label and properties (reuse vertex table logic for label/properties part)
					return nil, parseErr363
				}
			}
		}
	}

	vtd := &nodes.GraphTableDef{}
	parseErr364 := p.parseGraphLabelAndProperties(vtd)
	if parseErr364 != nil {
		return nil, parseErr364
	}
	def.Labels = vtd.Labels
	def.Properties = vtd.Properties
	def.PropColumns = vtd.PropColumns
	def.PropAliases = vtd.PropAliases
	def.DefaultLabel = vtd.DefaultLabel

	def.Loc.End = p.prev.End
	return def, nil
}

// parseGraphLabelAndProperties parses the label and properties clauses for graph tables.
func (p *Parser) parseGraphLabelAndProperties(def *nodes.GraphTableDef) error {
	// LABEL label_name [ LABEL label_name ]...
	for p.isIdentLikeStr("LABEL") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			def.Labels = append(def.Labels, p.cur.Str)
			p.advance()
		}
	}

	// PROPERTIES clause
	if p.isIdentLikeStr("PROPERTIES") {
		p.advance()
		// ARE ALL COLUMNS [EXCEPT ...]
		if p.isIdentLikeStr("ARE") {
			p.advance()
		}
		if p.cur.Type == kwALL {
			p.advance()
			// ALL COLUMNS EXCEPT (...)
			if p.isIdentLikeStr("COLUMNS") {
				p.advance()
			}
			if p.cur.Type == kwEXCEPT || p.isIdentLikeStr("EXCEPT") {
				p.advance()
				def.Properties = "ALL_EXCEPT"
				if p.cur.Type == '(' {
					var parseErr365 error
					def.PropColumns, parseErr365 = p.parseParenIdentList()
					if parseErr365 != nil {
						return parseErr365
					}
				}
			} else {
				def.Properties = "ALL"
			}
		} else if p.cur.Type == '(' {
			def.Properties = "LIST"
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					def.PropColumns = append(def.PropColumns, p.cur.Str)
					p.advance()
				}
				if p.cur.Type == kwAS {
					p.advance()
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
						def.PropAliases = append(def.PropAliases, p.cur.Str)
						p.advance()
					}
				}
				if p.cur.Type == ',' {
					p.advance()
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	} else if p.isIdentLikeStr("NO") {
		p.advance()
		if p.isIdentLikeStr("PROPERTIES") {
			p.advance()
			def.Properties = "NO"
		}
	}

	// DEFAULT LABEL
	if p.cur.Type == kwDEFAULT {
		p.advance()
		if p.isIdentLikeStr("LABEL") {
			p.advance()
			def.DefaultLabel = true
		}
	}
	return nil
}

// parseParenIdentList parses ( ident, ident, ... ) and returns the list of identifiers.
func (p *Parser) parseParenIdentList() ([]string, error) {
	var result []string
	if p.cur.Type != '(' {
		return result, nil
	}
	p.advance() // consume (
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
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
// BNF:
//
//	CREATE VECTOR INDEX [ IF NOT EXISTS ] [ schema. ] index_name
//	    ON [ schema. ] table_name ( column_name )
//	    [ INCLUDE ( column_name [, column_name ]... ) ]
//	    [ vector_index_organization_clause ]
//	    [ DISTANCE metric_name ]
//	    [ WITH TARGET ACCURACY integer [ PARAMETERS ( ... ) ] ]
//	    [ vector_index_hnsw_replication_clause ]
//	    [ ONLINE ]
//	    [ PARALLEL integer ] ;
func (p *Parser) parseCreateVectorIndexStmt(start int, ifNotExists bool) (*nodes.CreateVectorIndexStmt, error) {
	stmt := &nodes.CreateVectorIndexStmt{
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}

	// IF NOT EXISTS
	if !ifNotExists && p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
				stmt.IfNotExists = true
			}
		}
	}
	var parseErr366 error

	// Index name
	stmt.Name, parseErr366 = p.parseObjectName()
	if parseErr366 !=

		// ON [schema.]table_name (column_name)
		nil {
		return nil, parseErr366
	}

	if p.cur.Type == kwON {
		p.advance()
		var parseErr367 error
		stmt.TableName, parseErr367 = p.parseObjectName()
		if parseErr367 != nil {
			return nil, parseErr367
		}
		if p.cur.Type == '(' {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
				stmt.Column = p.cur.Str
				p.advance()
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	}

	// INCLUDE ( col, ... )
	if p.cur.Type == kwINCLUDE {
		p.advance()
		var parseErr368 error
		stmt.IncludeColumns, parseErr368 = p.parseParenIdentList()
		if parseErr368 !=

			// ORGANIZATION { INMEMORY NEIGHBOR GRAPH | NEIGHBOR PARTITIONS }
			nil {
			return nil, parseErr368
		}
	}

	if p.isIdentLikeStr("ORGANIZATION") {
		p.advance()
		if p.isIdentLikeStr("INMEMORY") {
			p.advance()
			stmt.Organization = "INMEMORY_NEIGHBOR_GRAPH"
			// consume NEIGHBOR GRAPH
			if p.isIdentLikeStr("NEIGHBOR") {
				p.advance()
			}
			if p.isIdentLikeStr("GRAPH") {
				p.advance()
			}
		} else if p.isIdentLikeStr("NEIGHBOR") {
			p.advance()
			stmt.Organization = "NEIGHBOR_PARTITIONS"
			if p.isIdentLikeStr("PARTITIONS") {
				p.advance()
			}
		}
	}

	// DISTANCE metric_name
	if p.isIdentLikeStr("DISTANCE") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Distance = p.cur.Str
			p.advance()
		}
	}

	// WITH TARGET ACCURACY integer [ PARAMETERS (...) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.isIdentLikeStr("TARGET") {
			p.advance()
		}
		if p.isIdentLikeStr("ACCURACY") {
			p.advance()
		}
		if p.cur.Type == tokICONST {
			var parseErr369 error
			stmt.TargetAccuracy, parseErr369 = p.parseIntValue()
			if parseErr369 !=

				// PARAMETERS ( ... )
				nil {
				return nil, parseErr369
			}
		}

		if p.isIdentLikeStr("PARAMETERS") {
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				// TYPE HNSW | IVF
				if p.cur.Type == kwTYPE || p.isIdentLikeStr("TYPE") {
					p.advance()
					if p.isIdentLikeStr("HNSW") {
						stmt.ParameterType = "HNSW"
						p.advance()
					} else if p.isIdentLikeStr("IVF") {
						stmt.ParameterType = "IVF"
						p.advance()
					}
				}
				// Parse remaining HNSW/IVF parameters
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.cur.Type == ',' {
						p.advance()
					}
					if p.isIdentLikeStr("NEIGHBORS") || p.isIdentLikeStr("M") {
						p.advance()
						if p.cur.Type == tokICONST {
							var parseErr370 error
							stmt.Neighbors, parseErr370 = p.parseIntValue()
							if parseErr370 != nil {
								return nil, parseErr370
							}
						}
					} else if p.isIdentLikeStr("EFCONSTRUCTION") {
						p.advance()
						if p.cur.Type == tokICONST {
							var parseErr371 error
							stmt.EfConstruction, parseErr371 = p.parseIntValue()
							if parseErr371 != nil {
								return nil, parseErr371
							}
						}
					} else if p.isIdentLikeStr("RESCORE_FACTOR") {
						p.advance()
						if p.cur.Type == tokICONST {
							var parseErr372 error
							stmt.RescoreFactor, parseErr372 = p.parseIntValue()
							if parseErr372 != nil {
								return nil, parseErr372
							}
						}
					} else if p.isIdentLikeStr("QUANTIZATION") {
						p.advance()
						if p.isIdentLike() || p.cur.Type == tokIDENT {
							stmt.Quantization = p.cur.Str
							p.advance()
						}
					} else if p.isIdentLikeStr("NEIGHBOR") {
						p.advance()
						if p.isIdentLikeStr("PARTITIONS") {
							p.advance()
						}
						if p.cur.Type == tokICONST {
							var parseErr373 error
							stmt.NeighborParts, parseErr373 = p.parseIntValue()
							if parseErr373 != nil {
								return nil, parseErr373
							}
						}
					} else if p.isIdentLikeStr("SAMPLES_PER_PARTITION") {
						p.advance()
						if p.cur.Type == tokICONST {
							var parseErr374 error
							stmt.SamplesPerPart, parseErr374 = p.parseIntValue()
							if parseErr374 != nil {
								return nil, parseErr374
							}
						}
					} else if p.isIdentLikeStr("MIN_VECTORS_PER_PARTITION") {
						p.advance()
						if p.cur.Type == tokICONST {
							var parseErr375 error
							stmt.MinVecsPerPart, parseErr375 = p.parseIntValue()
							if parseErr375 != nil {
								return nil, parseErr375
							}
						}
					} else {
						p.advance()
					}
				}
				if p.cur.Type == ')' {
					p.advance()
				}
			}
		}
	}

	// Replication clause: DUPLICATE ALL | DISTRIBUTE [AUTO | BY ...]
	if p.isIdentLikeStr("DUPLICATE") {
		p.advance()
		if p.cur.Type == kwALL {
			p.advance()
			stmt.Replication = "DUPLICATE_ALL"
		}
	} else if p.isIdentLikeStr("DISTRIBUTE") {
		p.advance()
		stmt.Replication = "DISTRIBUTE"
		if p.isIdentLikeStr("AUTO") {
			p.advance()
			stmt.Replication = "DISTRIBUTE_AUTO"
		} else if p.cur.Type == kwBY {
			p.advance()
			if p.isIdentLikeStr("ROWID") {
				p.advance()
				if p.isIdentLikeStr("RANGE") {
					p.advance()
				}
				stmt.Replication = "DISTRIBUTE_BY_ROWID_RANGE"
			} else if p.cur.Type == kwPARTITION {
				p.advance()
				stmt.Replication = "DISTRIBUTE_BY_PARTITION"
			} else if p.isIdentLikeStr("SUBPARTITION") {
				p.advance()
				stmt.Replication = "DISTRIBUTE_BY_SUBPARTITION"
			}
		}
	}

	// ONLINE
	if p.cur.Type == kwONLINE {
		p.advance()
		stmt.Online = true
	}

	// PARALLEL integer
	if p.cur.Type == kwPARALLEL {
		p.advance()
		if p.cur.Type == tokICONST {
			var parseErr376 error
			stmt.Parallel, parseErr376 = p.parseIntValue()
			if parseErr376 != nil {
				return nil, parseErr376
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE / ALTER LOCKDOWN PROFILE
// ---------------------------------------------------------------------------

// parseCreateLockdownProfileStmt parses a CREATE LOCKDOWN PROFILE statement.
// Called after CREATE LOCKDOWN PROFILE has been consumed.
//
// BNF:
//
//	CREATE LOCKDOWN PROFILE profile_name
//	    [ USING base_profile_name | INCLUDING base_profile_name ]
func (p *Parser) parseCreateLockdownProfileStmt(start int) (*nodes.CreateLockdownProfileStmt, error) {
	stmt := &nodes.CreateLockdownProfileStmt{
		Loc: nodes.Loc{Start: start},
	}

	// profile_name
	if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// USING base_profile_name
	if p.isIdentLikeStr("USING") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Using = p.cur.Str
			p.advance()
		}
	}

	// INCLUDING base_profile_name
	if p.isIdentLikeStr("INCLUDING") {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Including = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterLockdownProfileStmt parses an ALTER LOCKDOWN PROFILE statement.
// Called after ALTER LOCKDOWN PROFILE has been consumed.
//
// BNF:
//
//	ALTER LOCKDOWN PROFILE profile_name
//	    { lockdown_features | lockdown_options | lockdown_statements }
//	    [ USERS = { ALL | COMMON | LOCAL } ] ;
func (p *Parser) parseAlterLockdownProfileStmt(start int) (*nodes.AlterLockdownProfileStmt, error) {
	stmt := &nodes.AlterLockdownProfileStmt{
		Loc: nodes.Loc{Start: start},
	}

	// profile_name
	if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// { DISABLE | ENABLE }
	if p.cur.Type == kwDISABLE {
		stmt.Action = "DISABLE"
		p.advance()
	} else if p.cur.Type == kwENABLE {
		stmt.Action = "ENABLE"
		p.advance()
	}

	// { FEATURE = (...) | OPTION = (...) | STATEMENT = (...) }
	if p.isIdentLikeStr("FEATURE") {
		stmt.RuleType = "FEATURE"
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var parseErr377 error
		stmt.RuleItems, stmt.AllItems, stmt.ExceptItems, parseErr377 = p.parseLockdownItemList()
		if parseErr377 != nil {
			return nil, parseErr377
		}
	} else if p.cur.Type == kwOPTION || p.isIdentLikeStr("OPTION") {
		stmt.RuleType = "OPTION"
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var parseErr378 error
		stmt.RuleItems, stmt.AllItems, stmt.ExceptItems, parseErr378 = p.parseLockdownItemList()
		if parseErr378 != nil {
			return nil, parseErr378
		}
	} else if p.isIdentLikeStr("STATEMENT") {
		stmt.RuleType = "STATEMENT"
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		var parseErr379 error
		stmt.RuleItems, stmt.AllItems, stmt.ExceptItems, parseErr379 = p.parseLockdownItemList()
		if parseErr379 !=

			// CLAUSE = (...)
			nil {
			return nil, parseErr379
		}

		if p.isIdentLikeStr("CLAUSE") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			var parseErr380 error
			stmt.Clauses, stmt.ClauseAll, stmt.ClauseExceptItems, parseErr380 = p.parseLockdownItemList()
			if parseErr380 !=

				// OPTION = (...)
				nil {
				return nil, parseErr380
			}

			if p.cur.Type == kwOPTION || p.isIdentLikeStr("OPTION") {
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
				}
				var clauseAllItems bool
				var clauseExceptItems []string
				var parseErr381 error
				stmt.ClauseOptions, clauseAllItems, clauseExceptItems, parseErr381 = p.parseLockdownItemList()
				_ = clauseAllItems
				_ = clauseExceptItems
				if parseErr381 != nil {
					return nil, parseErr381
				}

				// VALUE / MINVALUE / MAXVALUE
				if p.isIdentLikeStr("VALUE") {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					var valueAllItems bool
					var valueExceptItems []string
					var parseErr382 error
					stmt.ValueItems, valueAllItems, valueExceptItems, parseErr382 = p.parseLockdownItemList()
					_ = valueAllItems
					_ = valueExceptItems
					if parseErr382 != nil {
						return nil, parseErr382
					}
				}
				if p.cur.Type == kwMINVALUE || p.isIdentLikeStr("MINVALUE") {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
						stmt.MinValue = p.cur.Str
						p.advance()
					}
				}
				if p.cur.Type == kwMAXVALUE || p.isIdentLikeStr("MAXVALUE") {
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
						stmt.MaxValue = p.cur.Str
						p.advance()
					}
				}
			}
		}
	}

	// USERS = { ALL | COMMON | LOCAL }
	if p.isIdentLikeStr("USERS") {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == kwALL {
			stmt.Users = "ALL"
			p.advance()
		} else if p.isIdentLikeStr("COMMON") {
			stmt.Users = "COMMON"
			p.advance()
		} else if p.isIdentLikeStr("LOCAL") {
			stmt.Users = "LOCAL"
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseLockdownItemList parses ( item, item, ... ) or ( ALL ) or ( ALL EXCEPT = ( ... ) ).
func (p *Parser) parseLockdownItemList() (items []string, allItems bool, exceptItems []string, parseErr error) {
	if p.cur.Type != '(' {
		return
	}
	p.advance() // consume (
	if p.cur.Type == kwALL {
		p.advance()
		// ALL EXCEPT = ( ... )
		if p.cur.Type == kwEXCEPT || p.isIdentLikeStr("EXCEPT") {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == '(' {
				p.advance()
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.cur.Type == tokSCONST {
						exceptItems = append(exceptItems, p.cur.Str)
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
			}
		}
		allItems = true
	} else {
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT || p.cur.Type == tokSCONST {
				items = append(items, p.cur.Str)
				p.advance()
			} else {
				p.advance()
			}
			if p.cur.Type == ',' {
				p.advance()
			}
		}
	}
	if p.cur.Type == ')' {
		p.advance()
	}
	return
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

	// Optional outline name — the name comes before FOR/FROM/ON.
	// We need to check if the current token is an identifier that is NOT a keyword
	// like FOR, FROM, ON which start the next clause.
	if p.cur.Type != kwFOR && p.cur.Type != kwFROM && p.cur.Type != kwON &&
		p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Name = p.cur.Str
			p.advance()
		}
	}

	// FOR CATEGORY category
	if p.cur.Type == kwFOR {
		p.advance()
		if p.isIdentLikeStr("CATEGORY") {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Category = p.cur.Str
			p.advance()
		}
	}

	// FROM [ PRIVATE ] source_outline
	if p.cur.Type == kwFROM {
		p.advance()
		if p.cur.Type == kwPRIVATE {
			stmt.FromPrivate = true
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.FromSource = p.cur.Str
			p.advance()
		}
	}

	// ON statement — skip to semicolon
	if p.cur.Type == kwON {
		p.advance()
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			p.advance()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterOutlineStmt parses an ALTER OUTLINE statement.
// Called after ALTER OUTLINE has been consumed.
//
// BNF:
//
//	ALTER OUTLINE [ PUBLIC | PRIVATE ] outline
//	    { REBUILD | RENAME TO new_outline_name | CHANGE CATEGORY TO new_category_name | { ENABLE | DISABLE } } ;
func (p *Parser) parseAlterOutlineStmt(start int) (*nodes.AlterOutlineStmt, error) {
	stmt := &nodes.AlterOutlineStmt{
		Loc: nodes.Loc{Start: start},
	}

	// PUBLIC | PRIVATE
	if p.cur.Type == kwPUBLIC {
		stmt.Public = true
		p.advance()
	} else if p.cur.Type == kwPRIVATE {
		stmt.Private = true
		p.advance()
	}

	// outline name
	if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// Action
	if p.isIdentLikeStr("REBUILD") {
		stmt.Action = "REBUILD"
		p.advance()
	} else if p.cur.Type == kwRENAME {
		stmt.Action = "RENAME"
		p.advance()
		if p.cur.Type == kwTO {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.NewName = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("CHANGE") {
		stmt.Action = "CHANGE_CATEGORY"
		p.advance()
		if p.isIdentLikeStr("CATEGORY") {
			p.advance()
		}
		if p.cur.Type == kwTO {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			stmt.Category = p.cur.Str
			p.advance()
		}
	} else if p.cur.Type == kwENABLE {
		stmt.Action = "ENABLE"
		p.advance()
	} else if p.cur.Type == kwDISABLE {
		stmt.Action = "DISABLE"
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Batch 105 — small objects bundle
// ---------------------------------------------------------------------------

// parseCreateJavaStmt parses a CREATE JAVA statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] JAVA has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] [ AND { RESOLVE | COMPILE } ] [ NOFORCE ]
//	    JAVA { SOURCE | CLASS | RESOURCE }
//	    [ NAMED [ schema. ] primary_name ]
//	    [ SHARING = { METADATA | NONE } ]
//	    [ invoker_rights_clause ]
//	    [ resolver_clause ]
//	    { USING { BFILE ( directory_object_name , server_file_name )
//	            | { CLOB | BLOB | BFILE } subquery
//	            | key_for_BLOB }
//	    | AS source_char }
func (p *Parser) parseCreateJavaStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_JAVA,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ AND { RESOLVE | COMPILE } ]
	if p.isIdentLike() && p.cur.Str == "AND" {
		p.advance()
		if p.isIdentLike() && (p.cur.Str == "RESOLVE" || p.cur.Str == "COMPILE") {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AND", Value: p.cur.Str})
			p.advance()
		}
	}

	// [ NOFORCE ]
	if p.isIdentLike() && p.cur.Str == "NOFORCE" {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "NOFORCE"})
		p.advance()
	}

	// { SOURCE | CLASS | RESOURCE }
	if p.isIdentLike() && (p.cur.Str == "SOURCE" || p.cur.Str == "CLASS" || p.cur.Str == "RESOURCE") {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "JAVA_TYPE", Value: p.cur.Str})
		p.advance()
	}

	// [ NAMED [ schema. ] primary_name ]
	if p.isIdentLike() && p.cur.Str == "NAMED" {
		p.advance()
		var parseErr383 error
		stmt.Name, parseErr383 = p.parseObjectName()
		if parseErr383 !=

			// Parse remaining clauses until ; or EOF
			nil {
			return nil, parseErr383
		}
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "SHARING" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHARING", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "AUTHID" {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AUTHID", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "RESOLVER" {
			p.advance()
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "RESOLVER"})
		} else if p.isIdentLike() && p.cur.Str == "USING" {
			p.advance()
			val := ""
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				val = p.cur.Str
				p.advance()
			}
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "USING", Value: val})
		} else if p.cur.Type == kwAS {
			p.advance()
			// source_char — typically a string constant
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS", Value: p.cur.Str})
				p.advance()
			}
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterJavaStmt parses an ALTER JAVA statement.
// Called after ALTER JAVA has been consumed.
//
// BNF:
//
//	ALTER JAVA [ IF EXISTS ] { SOURCE | CLASS } [ schema. ] object_name
//	    [ RESOLVER ( ( match_string schema_name ) [, ...] ) ]
//	    [ invoker_rights_clause ]
//	    { RESOLVE | COMPILE } ;
func (p *Parser) parseAlterJavaStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_JAVA,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ IF EXISTS ]
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance()
			p.advance()
			stmt.IfExists = true
		}
	}

	// { SOURCE | CLASS }
	if p.isIdentLike() && (p.cur.Str == "SOURCE" || p.cur.Str == "CLASS") {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "JAVA_TYPE", Value: p.cur.Str})
		p.advance()
	}
	var parseErr384 error

	stmt.Name, parseErr384 = p.parseObjectName()
	if parseErr384 !=

		// Parse remaining clauses
		nil {
		return nil, parseErr384
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "RESOLVER" {
			p.advance()
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "RESOLVER"})
		} else if p.isIdentLike() && p.cur.Str == "AUTHID" {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AUTHID", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && (p.cur.Str == "RESOLVE" || p.cur.Str == "COMPILE") {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: p.cur.Str})
			p.advance()
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateLibraryStmt parses a CREATE LIBRARY statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] [EDITIONABLE|NONEDITIONABLE] LIBRARY has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] [ EDITIONABLE | NONEDITIONABLE ]
//	    LIBRARY [ schema. ] library_name
//	    { IS | AS } library_path
//	    [ AGENT agent_dblink ]
//	    [ CREDENTIAL credential_name ]
//	    [ SHARING = { METADATA | NONE } ]
func (p *Parser) parseCreateLibraryStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_LIBRARY,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr385 error

	stmt.Name, parseErr385 = p.parseObjectName()
	if parseErr385 !=

		// { IS | AS } library_path
		nil {
		return nil, parseErr385
	}

	if p.cur.Type == kwIS || p.cur.Type == kwAS {
		p.advance()
	}
	if p.cur.Type == tokSCONST {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PATH", Value: p.cur.Str})
		p.advance()
	}

	// Parse optional clauses
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "AGENT" {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AGENT", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "CREDENTIAL" {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "CREDENTIAL", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "SHARING" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHARING", Value: p.cur.Str})
				p.advance()
			}
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterLibraryStmt parses an ALTER LIBRARY statement.
// Called after ALTER LIBRARY has been consumed.
//
// BNF:
//
//	ALTER LIBRARY [ IF EXISTS ] [ schema. ] library_name
//	    { library_compile_clause }
//	    [ EDITIONABLE | NONEDITIONABLE ] ;
//
//	library_compile_clause:
//	    COMPILE [ DEBUG ] [ compiler_parameters_clause ] [ REUSE SETTINGS ]
func (p *Parser) parseAlterLibraryStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_LIBRARY,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ IF EXISTS ]
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance()
			p.advance()
			stmt.IfExists = true
		}
	}
	var parseErr386 error

	stmt.Name, parseErr386 = p.parseObjectName()
	if parseErr386 !=

		// Parse remaining clauses
		nil {
		return nil, parseErr386
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "COMPILE" {
			p.advance()
			opt := &nodes.DDLOption{Key: "COMPILE"}
			if p.isIdentLike() && p.cur.Str == "DEBUG" {
				opt.Value = "DEBUG"
				p.advance()
			}
			opts.Items = append(opts.Items, opt)
			// compiler_parameters_clause: name = value [, ...]
			for p.isIdentLike() && p.cur.Str != "REUSE" &&
				p.cur.Type != ';' && p.cur.Type != tokEOF {
				if p.cur.Str == "EDITIONABLE" || p.cur.Str == "NONEDITIONABLE" {
					break
				}
				paramName := p.cur.Str
				p.advance()
				if p.cur.Type == '=' {
					p.advance()
					paramVal := ""
					if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokICONST || p.cur.Type == tokSCONST {
						paramVal = p.cur.Str
						p.advance()
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{Key: paramName, Value: paramVal})
				}
			}
			// [ REUSE SETTINGS ]
			if p.isIdentLike() && p.cur.Str == "REUSE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "SETTINGS" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "REUSE SETTINGS"})
			}
		} else if p.isIdentLike() && (p.cur.Str == "EDITIONABLE" || p.cur.Str == "NONEDITIONABLE") {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: p.cur.Str})
			p.advance()
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
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
	var parseErr387 error

	stmt.Name, parseErr387 = p.parseObjectName()
	if parseErr387 !=

		// [ SHARING = { METADATA | NONE } ]
		nil {
		return nil, parseErr387
	}

	if p.isIdentLike() && p.cur.Str == "SHARING" {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHARING", Value: p.cur.Str})
			p.advance()
		}
	}

	// AS 'path_name'
	if p.cur.Type == kwAS {
		p.advance()
		if p.cur.Type == tokSCONST {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS", Value: p.cur.Str})
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateContextStmt parses a CREATE CONTEXT statement.
// Called after CREATE [OR REPLACE] CONTEXT has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE ] CONTEXT namespace
//	    USING [ schema. ] package
//	    [ SHARING = { METADATA | NONE } ]
//	    [ INITIALIZED { EXTERNALLY | GLOBALLY } ]
//	    [ ACCESSED GLOBALLY ] ;
func (p *Parser) parseCreateContextStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_CONTEXT,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr388 error

	stmt.Name, parseErr388 = p.parseObjectName()
	if parseErr388 !=

		// USING [ schema. ] package
		nil {
		return nil, parseErr388
	}

	if p.isIdentLike() && p.cur.Str == "USING" {
		p.advance()
		usingName, parseErr389 := p.parseObjectName()
		if parseErr389 != nil {
			return nil, parseErr389
		}
		if usingName != nil {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "USING", Value: usingName.Name})
		}
	}

	// Parse remaining clauses
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "SHARING" {
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHARING", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "INITIALIZED" {
			p.advance()
			if p.isIdentLike() && (p.cur.Str == "EXTERNALLY" || p.cur.Str == "GLOBALLY") {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "INITIALIZED", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "ACCESSED" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "GLOBALLY" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "ACCESSED GLOBALLY"})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateMLEEnvStmt parses a CREATE MLE ENV statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] [PURE] MLE ENV has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] [ PURE ] MLE ENV
//	    [ schema. ] environment_name
//	    [ CLONE [ schema. ] existing_environment ]
//	    [ imports_clause ]
//	    [ language_options_clause ]
//
//	imports_clause: IMPORTS ( import_item [, import_item ]... )
//	import_item: import_name MODULE [ schema. ] module_name
//	language_options_clause: LANGUAGE OPTIONS language_options_string
func (p *Parser) parseCreateMLEEnvStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_MLE_ENV,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr390 error

	stmt.Name, parseErr390 = p.parseObjectName()
	if parseErr390 !=

		// [ CLONE [ schema. ] existing_environment ]
		nil {
		return nil, parseErr390
	}

	if p.isIdentLike() && p.cur.Str == "CLONE" {
		p.advance()
		cloneName, parseErr391 := p.parseObjectName()
		if parseErr391 != nil {
			return nil, parseErr391
		}
		if cloneName != nil {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "CLONE", Value: cloneName.Name})
		}
	}

	// Parse remaining clauses
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "IMPORTS" {
			p.advance()
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "IMPORTS"})
		} else if p.isIdentLike() && p.cur.Str == "LANGUAGE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "OPTIONS" {
				p.advance()
			}
			val := ""
			if p.cur.Type == tokSCONST {
				val = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "LANGUAGE OPTIONS", Value: val})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateMLEModuleStmt parses a CREATE MLE MODULE statement.
// Called after CREATE [OR REPLACE] [IF NOT EXISTS] MLE MODULE has been consumed.
//
// BNF:
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] MLE MODULE
//	    [ schema. ] module_name
//	    LANGUAGE JAVASCRIPT
//	    [ VERSION version_string ]
//	    { USING { CLOB ( subquery ) | BLOB ( subquery ) | BFILE ( subquery )
//	            | BFILE ( directory_object_name , server_file_name ) }
//	    | AS source_code }
func (p *Parser) parseCreateMLEModuleStmt(start int, orReplace bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_MLE_MODULE,
		OrReplace:  orReplace,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr392 error

	stmt.Name, parseErr392 = p.parseObjectName()
	if parseErr392 !=

		// LANGUAGE JAVASCRIPT
		nil {
		return nil, parseErr392
	}

	if p.isIdentLike() && p.cur.Str == "LANGUAGE" {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "JAVASCRIPT" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "LANGUAGE", Value: "JAVASCRIPT"})
			p.advance()
		}
	}

	// [ VERSION version_string ]
	if p.isIdentLike() && p.cur.Str == "VERSION" {
		p.advance()
		if p.cur.Type == tokSCONST {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "VERSION", Value: p.cur.Str})
			p.advance()
		}
	}

	// USING or AS
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "USING" {
			p.advance()
			val := ""
			if p.isIdentLike() || p.cur.Type == tokIDENT {
				val = p.cur.Str
				p.advance()
			}
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "USING", Value: val})
		} else if p.cur.Type == kwAS {
			p.advance()
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS", Value: p.cur.Str})
				p.advance()
			}
		} else {
			p.advance()
		}
	}

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

	// FROM
	if p.cur.Type == kwFROM {
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

	// [ = 'spfile_name' ] or [ 'spfile_name' ]
	if p.cur.Type == '=' {
		p.advance()
	}
	if p.cur.Type == tokSCONST {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SPFILE", Value: p.cur.Str})
		p.advance()
	}

	// FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

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
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateFlashbackArchiveStmt parses a CREATE FLASHBACK ARCHIVE statement.
// Called after CREATE FLASHBACK ARCHIVE has been consumed.
//
// BNF:
//
//	CREATE FLASHBACK ARCHIVE [ DEFAULT ] flashback_archive
//	    TABLESPACE tablespace_name
//	    [ flashback_archive_quota ]
//	    [ { NO OPTIMIZE DATA | OPTIMIZE DATA } ]
//	    flashback_archive_retention ;
//
//	flashback_archive_quota: QUOTA integer { K | M | G | T | P | E }
//	flashback_archive_retention: RETENTION integer { DAY | DAYS | MONTH | MONTHS | YEAR | YEARS }
func (p *Parser) parseCreateFlashbackArchiveStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_FLASHBACK_ARCHIVE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// [ DEFAULT ]
	if p.isIdentLike() && p.cur.Str == "DEFAULT" {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "DEFAULT"})
		p.advance()
	}
	var parseErr393 error

	stmt.Name, parseErr393 = p.parseObjectName()
	if parseErr393 !=

		// Parse remaining clauses
		nil {
		return nil, parseErr393
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwTABLESPACE {
			p.advance()
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "TABLESPACE", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "QUOTA" {
			p.advance()
			size, parseErr394 := p.parseSizeClause()
			if parseErr394 != nil {
				return nil, parseErr394
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "QUOTA", Value: size})
		} else if p.isIdentLike() && p.cur.Str == "NO" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "OPTIMIZE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "NO OPTIMIZE DATA"})
			}
		} else if p.isIdentLike() && p.cur.Str == "OPTIMIZE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DATA" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "OPTIMIZE DATA"})
		} else if p.isIdentLike() && p.cur.Str == "RETENTION" {
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.cur.Str
				p.advance()
			}
			if p.isIdentLike() {
				val += " " + p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "RETENTION", Value: val})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterFlashbackArchiveStmt parses an ALTER FLASHBACK ARCHIVE statement.
// Called after ALTER FLASHBACK ARCHIVE has been consumed.
//
// BNF:
//
//	ALTER FLASHBACK ARCHIVE flashback_archive_name
//	    { SET DEFAULT
//	    | { ADD | MODIFY } TABLESPACE tablespace_name [ flashback_archive_quota ]
//	    | REMOVE TABLESPACE tablespace_name
//	    | MODIFY RETENTION flashback_archive_retention
//	    | PURGE { ALL | BEFORE SCN scn_value | BEFORE TIMESTAMP timestamp_value }
//	    | [ NO ] OPTIMIZE DATA
//	    } ;
func (p *Parser) parseAlterFlashbackArchiveStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_FLASHBACK_ARCHIVE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr395 error

	stmt.Name, parseErr395 = p.parseObjectName()
	if parseErr395 !=

		// Parse action clause
		nil {
		return nil, parseErr395
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.cur.Type == kwSET {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DEFAULT" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SET DEFAULT"})
		} else if p.isIdentLike() && (p.cur.Str == "ADD" || p.cur.Str == "MODIFY") {
			action := p.cur.Str
			p.advance()
			if p.cur.Type == kwTABLESPACE {
				p.advance()
				tsName := ""
				if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
					tsName = p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: action + " TABLESPACE", Value: tsName})
				// optional quota
				if p.isIdentLike() && p.cur.Str == "QUOTA" {
					p.advance()
					size, parseErr396 := p.parseSizeClause()
					if parseErr396 != nil {
						return nil, parseErr396
					}
					opts.Items = append(opts.Items, &nodes.DDLOption{Key: "QUOTA", Value: size})
				}
			} else if p.isIdentLike() && p.cur.Str == "RETENTION" {
				p.advance()
				val := ""
				if p.cur.Type == tokICONST {
					val = p.cur.Str
					p.advance()
				}
				if p.isIdentLike() {
					val += " " + p.cur.Str
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "MODIFY RETENTION", Value: val})
			}
		} else if p.isIdentLike() && p.cur.Str == "REMOVE" {
			p.advance()
			if p.cur.Type == kwTABLESPACE {
				p.advance()
			}
			tsName := ""
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
				tsName = p.cur.Str
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "REMOVE TABLESPACE", Value: tsName})
		} else if p.isIdentLike() && p.cur.Str == "PURGE" {
			p.advance()
			val := ""
			if p.isIdentLike() && p.cur.Str == "ALL" {
				val = "ALL"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "BEFORE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "SCN" {
					p.advance()
					val = "BEFORE SCN"
					if p.cur.Type == tokICONST {
						val += " " + p.cur.Str
						p.advance()
					}
				} else if p.isIdentLike() && p.cur.Str == "TIMESTAMP" {
					p.advance()
					val = "BEFORE TIMESTAMP"
					// Skip the timestamp expression
					for p.cur.Type != ';' && p.cur.Type != tokEOF {
						p.advance()
					}
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PURGE", Value: val})
		} else if p.isIdentLike() && p.cur.Str == "NO" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "OPTIMIZE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "DATA" {
					p.advance()
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "NO OPTIMIZE DATA"})
			}
		} else if p.isIdentLike() && p.cur.Str == "OPTIMIZE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "DATA" {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "OPTIMIZE DATA"})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateRollbackSegmentStmt parses a CREATE ROLLBACK SEGMENT statement.
// Called after CREATE [PUBLIC] ROLLBACK SEGMENT has been consumed.
//
// BNF:
//
//	CREATE [ PUBLIC ] ROLLBACK SEGMENT rollback_segment
//	    [ TABLESPACE tablespace ]
//	    [ storage_clause ] ;
func (p *Parser) parseCreateRollbackSegmentStmt(start int, public bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_ROLLBACK_SEGMENT,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	if public {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PUBLIC"})
	}
	var parseErr397 error

	stmt.Name, parseErr397 = p.parseObjectName()
	if parseErr397 !=

		// [ TABLESPACE tablespace ]
		nil {
		return nil, parseErr397
	}

	if p.cur.Type == kwTABLESPACE {
		p.advance()
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "TABLESPACE", Value: p.cur.Str})
			p.advance()
		}
	}

	// [ storage_clause ] and remaining
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "STORAGE" {
			p.advance()
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "STORAGE"})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterRollbackSegmentStmt parses an ALTER ROLLBACK SEGMENT statement.
// Called after ALTER ROLLBACK SEGMENT has been consumed.
//
// BNF:
//
//	ALTER ROLLBACK SEGMENT rollback_segment
//	    { ONLINE | OFFLINE | storage_clause | SHRINK [ TO size_clause ] }
func (p *Parser) parseAlterRollbackSegmentStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_ROLLBACK_SEGMENT,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr398 error

	stmt.Name, parseErr398 = p.parseObjectName()
	if parseErr398 !=

		// Parse action
		nil {
		return nil, parseErr398
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "ONLINE" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "ONLINE"})
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "OFFLINE" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "OFFLINE"})
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "SHRINK" {
			p.advance()
			val := ""
			if p.cur.Type == kwTO {
				p.advance()
				var parseErr399 error
				val, parseErr399 = p.parseSizeClause()
				if parseErr399 != nil {
					return nil, parseErr399
				}
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SHRINK", Value: val})
		} else if p.isIdentLike() && p.cur.Str == "STORAGE" {
			p.advance()
			if p.cur.Type == '(' {
				p.skipParenthesized()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "STORAGE"})
		} else {
			p.advance()
		}
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

	// [ IF NOT EXISTS ] — may have been consumed by the caller already,
	// but check here for cases where EDITION is dispatched through parseCreateAdminObject.
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // IF
			p.advance() // NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // EXISTS
			}
		}
	}
	var parseErr400 error

	stmt.Name, parseErr400 = p.parseObjectName()
	if parseErr400 !=

		// [ AS CHILD OF parent_edition ]
		nil {
		return nil, parseErr400
	}

	if p.cur.Type == kwAS {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "CHILD" {
			p.advance()
		}
		if p.isIdentLike() && p.cur.Str == "OF" {
			p.advance()
		}
		if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS CHILD OF", Value: p.cur.Str})
			p.advance()
		}
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

// parseCreateRestorePointStmt parses a CREATE RESTORE POINT statement.
// Called after CREATE [CLEAN] RESTORE POINT has been consumed.
//
// BNF:
//
//	CREATE [ CLEAN ] RESTORE POINT restore_point
//	    [ FOR PLUGGABLE DATABASE pdb_name ]
//	    [ AS OF { TIMESTAMP | SCN } expr ]
//	    [ PRESERVE ]
//	    [ GUARANTEE FLASHBACK DATABASE ] ;
func (p *Parser) parseCreateRestorePointStmt(start int, clean bool) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_RESTORE_POINT,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	if clean {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "CLEAN"})
	}
	var parseErr402 error

	stmt.Name, parseErr402 = p.parseObjectName()
	if parseErr402 !=

		// Parse remaining clauses
		nil {
		return nil, parseErr402
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "FOR" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "PLUGGABLE" {
				p.advance()
			}
			if p.cur.Type == kwDATABASE {
				p.advance()
			}
			if p.isIdentLike() || p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "FOR PLUGGABLE DATABASE", Value: p.cur.Str})
				p.advance()
			}
		} else if p.cur.Type == kwAS {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "OF" {
				p.advance()
			}
			if p.isIdentLike() && (p.cur.Str == "TIMESTAMP" || p.cur.Str == "SCN") {
				asOfType := p.cur.Str
				p.advance()
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AS OF", Value: asOfType})
				// Skip the expression
				for p.cur.Type != ';' && p.cur.Type != tokEOF &&
					!(p.isIdentLike() && (p.cur.Str == "PRESERVE" || p.cur.Str == "GUARANTEE")) {
					p.advance()
				}
			}
		} else if p.isIdentLike() && p.cur.Str == "PRESERVE" {
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PRESERVE"})
			p.advance()
		} else if p.isIdentLike() && p.cur.Str == "GUARANTEE" {
			p.advance()
			if p.cur.Type == kwFLASHBACK {
				p.advance()
			}
			if p.cur.Type == kwDATABASE {
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "GUARANTEE FLASHBACK DATABASE"})
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateLogicalPartitionTrackingStmt parses a CREATE LOGICAL PARTITION TRACKING statement.
// Called after CREATE LOGICAL PARTITION TRACKING has been consumed.
//
// BNF:
//
//	CREATE LOGICAL PARTITION TRACKING ON [ schema. ] table_name
//	    PARTITION BY { RANGE | INTERVAL } ( column_name )
//	    ( partition_definition [, partition_definition ]... )
//
//	partition_definition: PARTITION partition_name VALUES LESS THAN ( value )
func (p *Parser) parseCreateLogicalPartitionTrackingStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_LOGICAL_PARTITION_TRACKING,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}

	// ON [ schema. ] table_name
	if p.cur.Type == kwON {
		p.advance()
	}
	var parseErr403 error
	stmt.Name, parseErr403 = p.parseObjectName()
	if parseErr403 !=

		// PARTITION BY { RANGE | INTERVAL } ( column_name )
		nil {
		return nil, parseErr403
	}

	if p.cur.Type == kwPARTITION || (p.isIdentLike() && p.cur.Str == "PARTITION") {
		p.advance()
	}
	if p.cur.Type == kwBY {
		p.advance()
	}
	if p.cur.Type == kwRANGE || (p.isIdentLike() && (p.cur.Str == "RANGE" || p.cur.Str == "INTERVAL")) {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PARTITION BY", Value: p.cur.Str})
		p.advance()
	}
	if p.cur.Type == '(' {
		p.skipParenthesized()
	}

	// ( partition_definition [, ...] )
	if p.cur.Type == '(' {
		p.skipParenthesized()
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PARTITIONS"})
	}

	// Skip remaining
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreatePmemFilestoreStmt parses a CREATE PMEM FILESTORE statement.
// Called after CREATE PMEM FILESTORE has been consumed.
//
// BNF:
//
//	CREATE PMEM FILESTORE filestore_name
//	    MOUNTPOINT 'file_path'
//	    BACKINGFILE 'backing_file_path' SIZE size_value BLOCKSIZE blocksize_value
//	    [ AUTOEXTEND { ON | OFF } [ NEXT size_value ] [ MAXSIZE { UNLIMITED | size_value } ] ] ;
func (p *Parser) parseCreatePmemFilestoreStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PMEM_FILESTORE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	var parseErr404 error

	stmt.Name, parseErr404 = p.parseObjectName()
	if parseErr404 !=

		// Parse clauses
		nil {
		return nil, parseErr404
	}

	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		if p.isIdentLike() && p.cur.Str == "MOUNTPOINT" {
			p.advance()
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "MOUNTPOINT", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "BACKINGFILE" {
			p.advance()
			if p.cur.Type == tokSCONST {
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "BACKINGFILE", Value: p.cur.Str})
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "SIZE" {
			p.advance()
			size, parseErr405 := p.parseSizeClause()
			if parseErr405 != nil {
				return nil, parseErr405
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "SIZE", Value: size})
		} else if p.isIdentLike() && p.cur.Str == "BLOCKSIZE" {
			p.advance()
			size, parseErr406 := p.parseSizeClause()
			if parseErr406 != nil {
				return nil, parseErr406
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "BLOCKSIZE", Value: size})
		} else if p.isIdentLike() && p.cur.Str == "AUTOEXTEND" {
			p.advance()
			val := ""
			if p.cur.Type == kwON || (p.isIdentLike() && p.cur.Str == "ON") {
				val = "ON"
				p.advance()
			} else if p.isIdentLike() && p.cur.Str == "OFF" {
				val = "OFF"
				p.advance()
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "AUTOEXTEND", Value: val})
			// [ NEXT size_value ]
			if p.isIdentLike() && p.cur.Str == "NEXT" {
				p.advance()
				nextSize, parseErr407 := p.parseSizeClause()
				if parseErr407 != nil {
					return nil, parseErr407
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "NEXT", Value: nextSize})
			}
			// [ MAXSIZE { UNLIMITED | size_value } ]
			if p.isIdentLike() && p.cur.Str == "MAXSIZE" {
				p.advance()
				maxVal := ""
				if p.isIdentLike() && p.cur.Str == "UNLIMITED" {
					maxVal = "UNLIMITED"
					p.advance()
				} else {
					var parseErr408 error
					maxVal, parseErr408 = p.parseSizeClause()
					if parseErr408 != nil {
						return nil, parseErr408
					}
				}
				opts.Items = append(opts.Items, &nodes.DDLOption{Key: "MAXSIZE", Value: maxVal})
			}
		} else {
			p.advance()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
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
	if p.isIdentLike() && (p.cur.Str == "SOURCE" || p.cur.Str == "CLASS" || p.cur.Str == "RESOURCE") {
		p.advance()
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
