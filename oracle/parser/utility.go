package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseLockTableStmt parses a LOCK TABLE statement.
//
// BNF: oracle/parser/bnf/LOCK-TABLE.bnf
//
//	LOCK TABLE [ schema. ] { table | view } [ @dblink ]
//	    [ partition_extension_clause ]
//	    [, [ schema. ] { table | view } [ @dblink ]
//	       [ partition_extension_clause ] ]...
//	    IN lockmode MODE
//	    [ NOWAIT | WAIT integer ] ;
//
//	partition_extension_clause::=
//	    PARTITION ( partition )
//	  | PARTITION FOR ( partition_key_value )
//	  | SUBPARTITION ( subpartition )
//	  | SUBPARTITION FOR ( subpartition_key_value )
//
//	lockmode::=
//	    ROW SHARE
//	  | ROW EXCLUSIVE
//	  | SHARE UPDATE
//	  | SHARE
//	  | SHARE ROW EXCLUSIVE
//	  | EXCLUSIVE
func (p *Parser) parseLockTableStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume LOCK

	if p.cur.Type == kwTABLE {
		p.advance()
	}

	stmt := &nodes.LockTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Parse one or more table entries separated by commas.
	// Each entry is: [schema.]table[@dblink] [partition_extension_clause]
	for {
		itemStart := p.pos()
		item := &nodes.LockTableItem{
			Loc: nodes.Loc{Start: itemStart},
		}
		item.Table = p.parseObjectName()

		// Optional partition_extension_clause:
		//   PARTITION ( partition_name )
		//   PARTITION FOR ( partition_key_value )
		//   SUBPARTITION ( subpartition_name )
		//   SUBPARTITION FOR ( subpartition_key_value )
		if p.cur.Type == kwPARTITION || p.cur.Type == kwSUBPARTITION {
			if p.cur.Type == kwPARTITION {
				item.PartitionType = "PARTITION"
			} else {
				item.PartitionType = "SUBPARTITION"
			}
			p.advance()
			if p.cur.Type == kwFOR {
				item.PartitionFor = true
				p.advance()
			}
			if p.cur.Type == '(' {
				p.advance() // consume '('
				// Collect partition name/value as a string
				name := ""
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					if name != "" {
						name += " "
					}
					name += p.cur.Str
					p.advance()
				}
				item.PartitionName = name
				if p.cur.Type == ')' {
					p.advance() // consume ')'
				}
			}
		}

		item.Loc.End = p.pos()
		stmt.Tables = append(stmt.Tables, item)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}

	// IN
	if p.cur.Type == kwIN {
		p.advance()
	}

	// Lock mode: collect words until MODE
	mode := ""
	for p.cur.Type != kwMODE && p.cur.Type != tokEOF && p.cur.Type != ';' {
		if mode != "" {
			mode += " "
		}
		if p.cur.Type == kwSHARE {
			mode += "SHARE"
		} else if p.cur.Type == kwROW {
			mode += "ROW"
		} else if p.cur.Type == kwEXCLUSIVE {
			mode += "EXCLUSIVE"
		} else if p.isIdentLike() {
			mode += p.cur.Str
		}
		p.advance()
	}
	stmt.LockMode = mode

	// MODE
	if p.cur.Type == kwMODE {
		p.advance()
	}

	// NOWAIT or WAIT n
	if p.cur.Type == kwNOWAIT {
		stmt.Nowait = true
		p.advance()
	} else if p.cur.Type == kwWAIT {
		p.advance()
		stmt.Wait = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCallStmt parses a CALL statement.
//
// BNF: oracle/parser/bnf/CALL.bnf
//
//	CALL
//	    { routine_clause | object_access_expression }
//	    [ INTO :host_variable [ [ INDICATOR ] :indicator_variable ] ] ;
//
//	routine_clause:
//	    [ schema. ] [ { type_name | package_name } . ] routine_name [ @dblink_name ]
//	    ( [ argument [, argument ]... ] )
//
//	object_access_expression:
//	    ( expr ) . method_name ( [ argument [, argument ]... ] )
func (p *Parser) parseCallStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume CALL

	stmt := &nodes.CallStmt{
		Args: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}

	stmt.Name = p.parseObjectName()

	// Arguments
	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg := p.parseExpr()
			if arg != nil {
				stmt.Args.Items = append(stmt.Args.Items, arg)
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

	// INTO :bind_variable
	if p.cur.Type == kwINTO {
		p.advance()
		stmt.Into = p.parseExpr()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRenameStmt parses a RENAME statement.
//
// BNF: oracle/parser/bnf/RENAME.bnf
//
//	RENAME old_name TO new_name ;
func (p *Parser) parseRenameStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume RENAME

	stmt := &nodes.RenameStmt{
		Loc: nodes.Loc{Start: start},
	}

	stmt.OldName = p.parseObjectName()

	if p.cur.Type == kwTO {
		p.advance()
	}

	stmt.NewName = p.parseObjectName()

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// BNF: oracle/parser/bnf/TRUNCATE-TABLE.bnf
//
//	TRUNCATE TABLE [ schema. ] table_name
//	    [ { PRESERVE | PURGE } MATERIALIZED VIEW LOG ]
//	    [ { DROP STORAGE | DROP ALL STORAGE | REUSE STORAGE } ]
//	    [ CASCADE ] ;
func (p *Parser) parseTruncateStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume TRUNCATE

	stmt := &nodes.TruncateStmt{
		Loc: nodes.Loc{Start: start},
	}

	// TRUNCATE TABLE or TRUNCATE CLUSTER
	if p.cur.Type == kwTABLE {
		p.advance()
	} else if p.cur.Type == kwCLUSTER {
		stmt.Cluster = true
		p.advance()
	}

	// Parse table/cluster name
	stmt.Table = p.parseObjectName()

	// Parse optional clauses
	for {
		if p.cur.Type == kwPURGE {
			// PURGE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				stmt.PurgeMVLog = true
				p.advance() // consume MATERIALIZED
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else if p.cur.Type == kwCASCADE {
			stmt.Cascade = true
			p.advance()
		} else if p.cur.Type == kwDROP || p.isIdentLike() && p.cur.Str == "REUSE" {
			// DROP STORAGE or REUSE STORAGE
			p.advance()
			if p.cur.Type == kwSTORAGE {
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "PRESERVE" {
			// PRESERVE MATERIALIZED VIEW LOG
			p.advance()
			if p.cur.Type == kwMATERIALIZED {
				p.advance()
				if p.cur.Type == kwVIEW {
					p.advance()
				}
				if p.cur.Type == kwLOG {
					p.advance()
				}
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAnalyzeStmt parses an ANALYZE statement.
//
// BNF: oracle/parser/bnf/ANALYZE.bnf
//
//	ANALYZE
//	    { TABLE [ schema. ] table [ partition_extension_clause ]
//	    | INDEX [ schema. ] index [ partition_extension_clause ]
//	    | CLUSTER [ schema. ] cluster
//	    }
//	    { validation_clauses
//	    | DELETE [ SYSTEM ] STATISTICS
//	    } ;
//
//	partition_extension_clause:
//	    { PARTITION ( partition_name )
//	    | SUBPARTITION ( subpartition_name )
//	    }
//
//	validation_clauses:
//	    { VALIDATE REF UPDATE [ SET DANGLING TO NULL ]
//	    | VALIDATE STRUCTURE [ CASCADE [ FAST ] ] [ ONLINE | OFFLINE ] [ into_clause ]
//	    | LIST CHAINED ROWS [ into_clause ]
//	    }
//
//	into_clause:
//	    INTO [ schema. ] table
func (p *Parser) parseAnalyzeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ANALYZE

	stmt := &nodes.AnalyzeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Object type: TABLE, INDEX, or CLUSTER
	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
	case kwCLUSTER:
		stmt.ObjectType = nodes.OBJECT_CLUSTER
		p.advance()
	default:
		stmt.ObjectType = nodes.OBJECT_TABLE
	}

	// Object name (possibly schema-qualified)
	stmt.Table = p.parseObjectName()

	// Skip optional partition_extension_clause:
	//   PARTITION ( name ) | SUBPARTITION ( name )
	if p.isIdentLike() && (p.cur.Str == "PARTITION" || p.cur.Str == "SUBPARTITION") {
		p.advance() // consume PARTITION/SUBPARTITION
		if p.cur.Type == '(' {
			p.advance() // consume (
			p.parseExpr()
			if p.cur.Type == ')' {
				p.advance() // consume )
			}
		}
	}

	// Action clause
	switch {
	case p.isIdentLike() && p.cur.Str == "COMPUTE":
		// COMPUTE STATISTICS
		p.advance() // consume COMPUTE
		stmt.Action = "COMPUTE STATISTICS"
		if p.cur.Type == kwSTATISTICS {
			p.advance() // consume STATISTICS
		}
		// Optional FOR clause - consume until end of statement
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			// Consume the rest of the FOR clause (TABLE | ALL [INDEXED] COLUMNS [SIZE n] | COLUMNS ...)
			for p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
		}

	case p.isIdentLike() && p.cur.Str == "ESTIMATE":
		// ESTIMATE STATISTICS [FOR ...] [SAMPLE n {ROWS|PERCENT}]
		p.advance() // consume ESTIMATE
		stmt.Action = "ESTIMATE STATISTICS"
		if p.cur.Type == kwSTATISTICS {
			p.advance() // consume STATISTICS
		}
		// Optional FOR clause
		if p.cur.Type == kwFOR {
			p.advance() // consume FOR
			for p.cur.Type != kwSAMPLE && p.cur.Type != ';' && p.cur.Type != tokEOF {
				p.advance()
			}
		}
		// Optional SAMPLE n {ROWS|PERCENT}
		if p.cur.Type == kwSAMPLE {
			p.advance() // consume SAMPLE
			if p.cur.Type == tokICONST {
				stmt.SampleValue = p.parseIntValue()
			}
			if p.cur.Type == kwROWS {
				stmt.SampleUnit = "ROWS"
				p.advance()
			} else if p.cur.Type == kwPERCENT {
				stmt.SampleUnit = "PERCENT"
				p.advance()
			}
		}

	case p.cur.Type == kwDELETE:
		// DELETE [SYSTEM] STATISTICS
		p.advance() // consume DELETE
		stmt.Action = "DELETE STATISTICS"
		if p.cur.Type == kwSYSTEM {
			p.advance() // consume SYSTEM
			stmt.DeleteSystem = true
			stmt.Action = "DELETE SYSTEM STATISTICS"
		}
		if p.cur.Type == kwSTATISTICS {
			p.advance() // consume STATISTICS
		}

	case p.cur.Type == kwVALIDATE:
		// VALIDATE REF UPDATE [SET DANGLING TO NULL]
		// VALIDATE STRUCTURE [CASCADE [FAST]] [ONLINE|OFFLINE] [INTO table]
		p.advance() // consume VALIDATE
		if p.cur.Type == kwREF {
			// VALIDATE REF UPDATE [SET DANGLING TO NULL]
			p.advance() // consume REF
			stmt.Action = "VALIDATE REF UPDATE"
			if p.cur.Type == kwUPDATE {
				p.advance() // consume UPDATE
			}
			if p.cur.Type == kwSET {
				p.advance() // consume SET
				// DANGLING is an identifier
				if p.isIdentLike() && p.cur.Str == "DANGLING" {
					p.advance() // consume DANGLING
					if p.cur.Type == kwTO {
						p.advance() // consume TO
					}
					if p.cur.Type == kwNULL {
						p.advance() // consume NULL
						stmt.SetDanglingNull = true
					}
				}
			}
		} else if p.isIdentLike() && p.cur.Str == "STRUCTURE" {
			// VALIDATE STRUCTURE [CASCADE [FAST]] [ONLINE|OFFLINE] [INTO table]
			p.advance() // consume STRUCTURE
			stmt.Action = "VALIDATE STRUCTURE"
			// Optional CASCADE [FAST]
			if p.cur.Type == kwCASCADE {
				p.advance() // consume CASCADE
				if p.isIdentLike() && p.cur.Str == "FAST" {
					p.advance() // consume FAST
					stmt.CascadeFast = true
				}
			}
			// Optional ONLINE | OFFLINE
			if p.cur.Type == kwONLINE {
				p.advance() // consume ONLINE
				stmt.Online = true
			} else if p.cur.Type == kwOFFLINE {
				p.advance() // consume OFFLINE
				stmt.Offline = true
			}
			// Optional INTO table
			if p.cur.Type == kwINTO {
				p.advance() // consume INTO
				stmt.IntoTable = p.parseObjectName()
			}
		}

	case p.cur.Type == kwLIST:
		// LIST CHAINED ROWS [INTO table]
		p.advance() // consume LIST
		stmt.Action = "LIST CHAINED ROWS"
		// CHAINED is an identifier
		if p.isIdentLike() && p.cur.Str == "CHAINED" {
			p.advance() // consume CHAINED
		}
		if p.cur.Type == kwROWS {
			p.advance() // consume ROWS
		}
		// Optional INTO table
		if p.cur.Type == kwINTO {
			p.advance() // consume INTO
			stmt.IntoTable = p.parseObjectName()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseExplainPlanStmt parses an EXPLAIN PLAN statement.
//
// BNF: oracle/parser/bnf/EXPLAIN-PLAN.bnf
//
//	EXPLAIN PLAN
//	    [ SET STATEMENT_ID = string ]
//	    [ INTO [ schema. ] table [ @dblink ] ]
//	    FOR statement ;
func (p *Parser) parseExplainPlanStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume EXPLAIN

	// Expect PLAN
	if p.cur.Type == kwPLAN {
		p.advance()
	}

	stmt := &nodes.ExplainPlanStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional SET STATEMENT_ID = 'id'
	if p.cur.Type == kwSET {
		p.advance() // consume SET
		// STATEMENT_ID is an identifier
		if p.isIdentLike() && p.cur.Str == "STATEMENT_ID" {
			p.advance() // consume STATEMENT_ID
			if p.cur.Type == '=' {
				p.advance() // consume =
			}
			if p.cur.Type == tokSCONST {
				stmt.StatementID = p.cur.Str
				p.advance()
			}
		}
	}

	// Optional INTO [schema.]table
	if p.cur.Type == kwINTO {
		p.advance()
		stmt.Into = p.parseObjectName()
	}

	// FOR statement
	if p.cur.Type == kwFOR {
		p.advance()
		stmt.Statement = p.parseStmt()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFlashbackTableStmt parses a FLASHBACK TABLE statement.
//
// BNF: oracle/parser/bnf/FLASHBACK-TABLE.bnf
//
//	FLASHBACK TABLE [ schema. ] table [, [ schema. ] table ]...
//	    { TO SCN expr
//	    | TO TIMESTAMP expr
//	    | TO RESTORE POINT restore_point_name
//	    | TO BEFORE DROP [ RENAME TO table ]
//	    }
//	    [ { ENABLE | DISABLE } TRIGGERS ] ;
func (p *Parser) parseFlashbackTableStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume FLASHBACK

	// Expect TABLE
	if p.cur.Type == kwTABLE {
		p.advance()
	}

	stmt := &nodes.FlashbackTableStmt{
		Loc: nodes.Loc{Start: start},
	}

	// One or more table names (comma-separated), ending before TO
	first := p.parseObjectName()
	stmt.Tables = append(stmt.Tables, first)
	stmt.Table = first
	for p.cur.Type == ',' {
		p.advance() // consume ','
		t := p.parseObjectName()
		stmt.Tables = append(stmt.Tables, t)
	}

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// BEFORE variant or RESTORE POINT or SCN expr | TIMESTAMP expr
	if p.cur.Type == kwBEFORE {
		p.advance() // consume BEFORE
		if p.cur.Type == kwDROP {
			p.advance() // consume DROP
			stmt.ToBeforeDrop = true
			// Optional RENAME TO name
			if p.cur.Type == kwRENAME {
				p.advance() // consume RENAME
				if p.cur.Type == kwTO {
					p.advance()
				}
				stmt.Rename = p.parseIdentifier()
			}
		} else {
			// TO BEFORE SCN expr | TO BEFORE TIMESTAMP expr
			stmt.Before = true
			switch p.cur.Type {
			case kwSCN:
				p.advance()
				stmt.ToSCN = p.parseExpr()
			case kwTIMESTAMP:
				p.advance()
				stmt.ToTimestamp = p.parseExpr()
			}
		}
	} else {
		switch p.cur.Type {
		case kwSCN:
			p.advance()
			stmt.ToSCN = p.parseExpr()
		case kwTIMESTAMP:
			p.advance()
			stmt.ToTimestamp = p.parseExpr()
		default:
			// TO RESTORE POINT restore_point_name
			if p.isIdentLike() && p.cur.Str == "RESTORE" {
				p.advance()
				if p.isIdentLike() && p.cur.Str == "POINT" {
					p.advance()
				}
				if p.isIdentLike() {
					stmt.ToRestorePoint = p.cur.Str
					p.advance()
				}
			}
		}
	}

	// Optional { ENABLE | DISABLE } TRIGGERS
	if p.cur.Type == kwENABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "TRIGGERS" {
			p.advance()
		}
		t := true
		stmt.EnableTriggers = &t
	} else if p.cur.Type == kwDISABLE {
		p.advance()
		if p.isIdentLike() && p.cur.Str == "TRIGGERS" {
			p.advance()
		}
		f := false
		stmt.EnableTriggers = &f
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseFlashbackDatabaseStmt parses a FLASHBACK DATABASE statement.
//
// BNF: oracle/parser/bnf/FLASHBACK-DATABASE.bnf
//
//	FLASHBACK [ STANDBY | PLUGGABLE ] DATABASE [ database ]
//	    { TO SCN scn_number
//	    | TO BEFORE SCN scn_number
//	    | TO TIMESTAMP timestamp_expression
//	    | TO BEFORE TIMESTAMP timestamp_expression
//	    | TO RESTORE POINT restore_point_name
//	    | TO BEFORE RESETLOGS
//	    } ;
func (p *Parser) parseFlashbackDatabaseStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume FLASHBACK

	stmt := &nodes.FlashbackDatabaseStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional STANDBY | PLUGGABLE
	if p.isIdentLike() && p.cur.Str == "STANDBY" {
		stmt.Modifier = "STANDBY"
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "PLUGGABLE" {
		stmt.Modifier = "PLUGGABLE"
		p.advance()
	}

	// DATABASE
	if p.cur.Type == kwDATABASE {
		p.advance()
	}

	// Optional database name (not TO keyword)
	if p.isIdentLike() && p.cur.Str != "TO" {
		stmt.DatabaseName = p.parseObjectName()
	}

	// TO
	if p.cur.Type == kwTO {
		p.advance()
	}

	// Optional BEFORE
	if p.isIdentLike() && p.cur.Str == "BEFORE" {
		stmt.Before = true
		p.advance()
	}

	switch p.cur.Type {
	case kwSCN:
		p.advance()
		stmt.ToSCN = p.parseExpr()
	case kwTIMESTAMP:
		p.advance()
		stmt.ToTimestamp = p.parseExpr()
	default:
		if p.isIdentLike() && p.cur.Str == "RESTORE" {
			p.advance()
			if p.isIdentLike() && p.cur.Str == "POINT" {
				p.advance()
			}
			if p.isIdentLike() {
				stmt.ToRestorePoint = p.cur.Str
				p.advance()
			}
		} else if p.isIdentLike() && p.cur.Str == "RESETLOGS" {
			stmt.ToResetlogs = true
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parsePurgeStmt parses a PURGE statement.
//
// BNF: oracle/parser/bnf/PURGE.bnf
//
//	PURGE { TABLE [ schema. ] table
//	      | INDEX [ schema. ] index
//	      | TABLESPACE tablespace [ USER user ]
//	      | TABLESPACE SET tablespace_set [ USER user ]
//	      | RECYCLEBIN
//	      | DBA_RECYCLEBIN
//	      } ;
func (p *Parser) parsePurgeStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume PURGE

	stmt := &nodes.PurgeStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
		stmt.Name = p.parseObjectName()
	case kwTABLESPACE:
		stmt.ObjectType = nodes.OBJECT_TABLESPACE
		p.advance()
		stmt.Name = p.parseObjectName()
	default:
		// RECYCLEBIN or DBA_RECYCLEBIN (parsed as identifiers)
		if p.isIdentLike() {
			ident := p.cur.Str
			p.advance()
			stmt.Name = &nodes.ObjectName{
				Name: ident,
				Loc:  nodes.Loc{Start: start, End: p.pos()},
			}
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseAuditStmt parses an AUDIT statement (Traditional + Unified Auditing).
//
// BNF: oracle/parser/bnf/AUDIT-Traditional-Auditing.bnf
//
//	AUDIT { audit_operation_clause [ auditing_by_clause | auditing_on_clause ]
//	      | audit_schema_object_clause
//	      }
//	    [ BY { SESSION | ACCESS } ]
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ]
//
//	audit_operation_clause:
//	    { sql_statement_shortcut
//	    | system_privilege
//	    | ALL
//	    | ALL STATEMENTS
//	    | ALL PRIVILEGES
//	    }
//	    [, { sql_statement_shortcut | system_privilege } ]...
//
//	auditing_by_clause:
//	    BY { user [, user ]...
//	       | SESSION CURRENT
//	       }
//
//	audit_schema_object_clause:
//	    sql_operation [, sql_operation ]...
//	    auditing_on_clause
//
//	auditing_on_clause:
//	    ON [ schema. ] object
//	  | ON DEFAULT
//	  | ON DIRECTORY directory_name
//	  | ON MINING MODEL model_name
//	  | ON SQL TRANSLATION PROFILE profile_name
//	  | NETWORK
//	  | DIRECT_PATH LOAD
//
// BNF: oracle/parser/bnf/AUDIT-Unified-Auditing.bnf
//
//	AUDIT
//	    { POLICY policy_name
//	        [ { BY | EXCEPT } user [, user ]... ]
//	        [ { BY | EXCEPT } USERS WITH ROLE role [, role ]... ]
//	        [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    | CONTEXT NAMESPACE namespace
//	        ATTRIBUTES attribute [, attribute ]...
//	        [ BY user [, user ]... ]
//	    } ;
func (p *Parser) parseAuditStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume AUDIT

	stmt := &nodes.AuditStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Unified: AUDIT POLICY ...
	if p.isIdentLikeStr("POLICY") {
		p.advance()
		if p.isIdentLike() {
			stmt.Policy = p.cur.Str
			p.advance()
		}
		// [ { BY | EXCEPT } user [, user]... ]
		// [ { BY | EXCEPT } USERS WITH ROLE role [, role]... ]
		for p.cur.Type != ';' && p.cur.Type != tokEOF {
			if p.cur.Type == kwBY {
				p.advance()
				if p.isIdentLikeStr("USERS") {
					// BY USERS WITH ROLE role [, role]...
					p.advance()
					if p.cur.Type == kwWITH {
						p.advance()
					}
					if p.cur.Type == kwROLE {
						p.advance()
					}
					stmt.WithRoles = p.parseIdentListForAudit()
				} else {
					// BY user [, user]...
					stmt.ByUsers = p.parseIdentListForAudit()
				}
			} else if p.cur.Type == kwEXCEPT {
				p.advance()
				if p.isIdentLikeStr("USERS") {
					// EXCEPT USERS WITH ROLE role [, role]...
					p.advance()
					if p.cur.Type == kwWITH {
						p.advance()
					}
					if p.cur.Type == kwROLE {
						p.advance()
					}
					stmt.WithRoles = p.parseIdentListForAudit()
					stmt.WithRoleExcept = true
				} else {
					// EXCEPT user [, user]...
					stmt.ExceptUsers = p.parseIdentListForAudit()
				}
			} else if p.cur.Type == kwWHENEVER {
				stmt.When = p.parseWheneverClause()
			} else {
				break
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Unified: AUDIT CONTEXT NAMESPACE ...
	if p.cur.Type == kwCONTEXT {
		p.advance()
		if p.isIdentLikeStr("NAMESPACE") {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.ContextNS = p.cur.Str
			p.advance()
		}
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance()
			stmt.ContextAttrs = p.parseIdentListForAudit()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Traditional: parse audit actions
	stmt.Actions = p.parseAuditActions()

	// ON clause or auditing_by_clause
	if p.cur.Type == kwON {
		p.advance()
		p.parseTraditionalAuditOnClause(stmt)
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	}

	// BY { SESSION | ACCESS } or BY user [, user]...
	if p.cur.Type == kwBY {
		p.advance()
		tok := p.cur.Str
		if tok == "SESSION" || tok == "ACCESS" {
			stmt.By = tok
			p.advance()
			// Check for CURRENT after SESSION
			if tok == "SESSION" && p.isIdentLikeStr("CURRENT") {
				stmt.By = "SESSION CURRENT"
				p.advance()
			}
		} else {
			// BY user [, user]...
			stmt.ByUsers2 = p.parseIdentListForAudit()
		}
	}

	// WHENEVER [NOT] SUCCESSFUL
	if p.cur.Type == kwWHENEVER {
		stmt.When = p.parseWheneverClause()
	}

	// CONTAINER = { CURRENT | ALL }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseNoauditStmt parses a NOAUDIT statement (Traditional + Unified Auditing).
//
// BNF: oracle/parser/bnf/NOAUDIT-Traditional-Auditing.bnf
//
//	NOAUDIT
//	    { audit_operation_clause [ auditing_by_clause ]
//	    | audit_schema_object_clause
//	    }
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ] ;
//
//	audit_operation_clause::=
//	    { statement_option [, statement_option ]...
//	    | ALL
//	    | ALL STATEMENTS
//	    | system_privilege [, system_privilege ]...
//	    | ALL PRIVILEGES
//	    }
//
//	auditing_by_clause::=
//	    BY user [, user ]...
//
//	audit_schema_object_clause::=
//	    { sql_operation [, sql_operation ]...
//	    | ALL
//	    }
//	    ON auditing_on_clause
//
//	auditing_on_clause::=
//	    [ schema. ] object
//	  | DIRECTORY directory_name
//	  | SQL TRANSLATION PROFILE [ schema. ] profile
//	  | DEFAULT
//	  | NETWORK
//	  | DIRECT_PATH LOAD
//
// BNF: oracle/parser/bnf/NOAUDIT-Unified-Auditing.bnf
//
//	NOAUDIT POLICY policy
//	    [ BY user [, user ]...
//	      [ WITH ROLE role [, role ]... ]
//	    ] ;
//
//	NOAUDIT CONTEXT NAMESPACE namespace
//	    ATTRIBUTES attribute [, attribute ]...
//	    [ BY user [, user ]...
//	      [ WITH ROLE role [, role ]... ]
//	    ] ;
func (p *Parser) parseNoauditStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume NOAUDIT

	stmt := &nodes.NoauditStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Unified: NOAUDIT POLICY ...
	if p.isIdentLikeStr("POLICY") {
		p.advance()
		if p.isIdentLike() {
			stmt.Policy = p.cur.Str
			p.advance()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				stmt.WithRoles = p.parseIdentListForAudit()
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Unified: NOAUDIT CONTEXT NAMESPACE ...
	if p.cur.Type == kwCONTEXT {
		p.advance()
		if p.isIdentLikeStr("NAMESPACE") {
			p.advance()
		}
		if p.isIdentLike() {
			stmt.ContextNS = p.cur.Str
			p.advance()
		}
		if p.isIdentLikeStr("ATTRIBUTES") {
			p.advance()
			stmt.ContextAttrs = p.parseIdentListForAudit()
		}
		if p.cur.Type == kwBY {
			p.advance()
			stmt.ByUsers = p.parseIdentListForAudit()
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				stmt.WithRoles = p.parseIdentListForAudit()
			}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Traditional: parse audit actions
	stmt.Actions = p.parseAuditActions()

	// ON clause
	if p.cur.Type == kwON {
		p.advance()
		p.parseTraditionalNoauditOnClause(stmt)
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	}

	// BY user [, user]...
	if p.cur.Type == kwBY {
		p.advance()
		stmt.ByUsers2 = p.parseIdentListForAudit()
	}

	// WHENEVER [NOT] SUCCESSFUL
	if p.cur.Type == kwWHENEVER {
		stmt.When = p.parseWheneverClause()
	}

	// CONTAINER = { CURRENT | ALL }
	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		stmt.ContainerAll = p.parseContainerClause()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseTraditionalAuditOnClause parses the ON clause for traditional AUDIT.
// Called after ON has been consumed.
func (p *Parser) parseTraditionalAuditOnClause(stmt *nodes.AuditStmt) {
	if p.cur.Type == kwDEFAULT {
		stmt.OnDefault = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECTORY") {
		p.advance()
		if p.isIdentLike() {
			stmt.OnDirectory = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("MINING") {
		p.advance()
		if p.cur.Type == kwMODEL {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else if p.isIdentLikeStr("SQL") {
		// SQL TRANSLATION PROFILE
		p.advance()
		if p.isIdentLikeStr("TRANSLATION") {
			p.advance()
		}
		if p.cur.Type == kwPROFILE {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else {
		stmt.Object = p.parseObjectName()
	}
}

// parseTraditionalNoauditOnClause parses the ON clause for traditional NOAUDIT.
// Called after ON has been consumed.
func (p *Parser) parseTraditionalNoauditOnClause(stmt *nodes.NoauditStmt) {
	if p.cur.Type == kwDEFAULT {
		stmt.OnDefault = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECTORY") {
		p.advance()
		if p.isIdentLike() {
			stmt.OnDirectory = p.cur.Str
			p.advance()
		}
	} else if p.isIdentLikeStr("SQL") {
		// SQL TRANSLATION PROFILE
		p.advance()
		if p.isIdentLikeStr("TRANSLATION") {
			p.advance()
		}
		if p.cur.Type == kwPROFILE {
			p.advance()
		}
		stmt.Object = p.parseObjectName()
	} else if p.isIdentLikeStr("NETWORK") {
		stmt.OnNetwork = true
		p.advance()
	} else if p.isIdentLikeStr("DIRECT_PATH") {
		stmt.OnDirectPath = true
		p.advance()
		if p.isIdentLikeStr("LOAD") {
			p.advance()
		}
	} else {
		stmt.Object = p.parseObjectName()
	}
}

// parseAuditActions collects audit action identifiers separated by commas.
func (p *Parser) parseAuditActions() []string {
	var actions []string
	for {
		// Collect multi-word action (e.g., "CREATE TABLE", "ALTER SESSION")
		action := ""
		for p.isIdentLike() || p.cur.Type == kwSELECT || p.cur.Type == kwINSERT ||
			p.cur.Type == kwUPDATE || p.cur.Type == kwDELETE || p.cur.Type == kwCREATE ||
			p.cur.Type == kwALTER || p.cur.Type == kwDROP || p.cur.Type == kwGRANT ||
			p.cur.Type == kwEXECUTE || p.cur.Type == kwINDEX || p.cur.Type == kwALL {
			// Stop before special clause identifiers
			if p.isIdentLikeStr("NETWORK") || p.isIdentLikeStr("DIRECT_PATH") || p.isIdentLikeStr("CONTAINER") {
				break
			}
			if action != "" {
				action += " "
			}
			action += p.cur.Str
			p.advance()
			// Stop if we hit keywords that start a clause
			if p.cur.Type == kwON || p.cur.Type == kwBY || p.cur.Type == kwWHENEVER {
				break
			}
		}
		if action != "" {
			actions = append(actions, action)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return actions
}

// parseWheneverClause parses WHENEVER [NOT] SUCCESSFUL.
func (p *Parser) parseWheneverClause() string {
	result := "WHENEVER"
	p.advance() // consume WHENEVER
	if p.cur.Type == kwNOT {
		result += " NOT"
		p.advance()
	}
	if p.cur.Type == kwSUCCESSFUL {
		result += " SUCCESSFUL"
		p.advance()
	}
	return result
}

// parseIdentListForAudit parses a comma-separated list of identifiers,
// stopping at clause keywords (BY, EXCEPT, WHENEVER, WITH, CONTAINER, ;, EOF).
func (p *Parser) parseIdentListForAudit() []string {
	var list []string
	for {
		if !p.isIdentLike() {
			break
		}
		list = append(list, p.cur.Str)
		p.advance()
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	return list
}

// parseAssociateStatisticsStmt parses an ASSOCIATE STATISTICS statement.
//
// BNF: oracle/parser/bnf/ASSOCIATE-STATISTICS.bnf
//
//	ASSOCIATE STATISTICS
//	    WITH { column_association | function_association }
//	    using_statistics_type
//	    [ default_cost_clause ]
//	    [ default_selectivity_clause ]
//	    [ storage_table_clause ] ;
//
//	column_association:
//	    COLUMNS [ schema. ] table . column [, [ schema. ] table . column ]...
//
//	function_association:
//	    { FUNCTIONS [ schema. ] function [, [ schema. ] function ]...
//	    | PACKAGES [ schema. ] package [, [ schema. ] package ]...
//	    | TYPES [ schema. ] type [, [ schema. ] type ]...
//	    | DOMAIN INDEXES [ schema. ] index [, [ schema. ] index ]...
//	    | INDEXTYPES [ schema. ] indextype [, [ schema. ] indextype ]...
//	    }
//
//	using_statistics_type:
//	    USING { [ schema. ] statistics_type | NULL }
//
//	default_cost_clause:
//	    DEFAULT COST ( cpu_cost , io_cost , network_cost )
//
//	default_selectivity_clause:
//	    DEFAULT SELECTIVITY default_selectivity
//
//	storage_table_clause:
//	    WITH { SYSTEM MANAGED STORAGE TABLES
//	         | USER MANAGED STORAGE TABLES
//	         }
func (p *Parser) parseAssociateStatisticsStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume ASSOCIATE
	if p.cur.Type == kwSTATISTICS {
		p.advance()
	}

	stmt := &nodes.AssociateStatisticsStmt{
		Loc: nodes.Loc{Start: start},
	}

	// WITH
	if p.cur.Type == kwWITH {
		p.advance()
	}

	// Object type: COLUMNS, FUNCTIONS, PACKAGES, TYPES, INDEXES
	if p.isIdentLike() || p.cur.Type == kwINDEX {
		if p.cur.Type == kwINDEX {
			stmt.ObjectType = "INDEXES"
		} else {
			stmt.ObjectType = p.cur.Str
		}
		p.advance()
	}

	// Object names
	for {
		name := p.parseObjectName()
		if name != nil {
			stmt.Objects = append(stmt.Objects, name)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// USING statistics_type
	if p.isIdentLike() && p.cur.Str == "USING" {
		p.advance()
		stmt.Using = p.parseObjectName()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseDisassociateStatisticsStmt parses a DISASSOCIATE STATISTICS statement.
//
// BNF: oracle/parser/bnf/DISASSOCIATE-STATISTICS.bnf
//
//	DISASSOCIATE STATISTICS
//	    FROM { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES | INDEXTYPES }
//	    [ schema. ] object_name [, [ schema. ] object_name ]...
//	    [ FORCE ] ;
func (p *Parser) parseDisassociateStatisticsStmt() nodes.StmtNode {
	start := p.pos()
	p.advance() // consume DISASSOCIATE
	if p.cur.Type == kwSTATISTICS {
		p.advance()
	}

	stmt := &nodes.DisassociateStatisticsStmt{
		Loc: nodes.Loc{Start: start},
	}

	// FROM
	if p.cur.Type == kwFROM {
		p.advance()
	}

	// Object type
	if p.isIdentLike() || p.cur.Type == kwINDEX {
		if p.cur.Type == kwINDEX {
			stmt.ObjectType = "INDEXES"
		} else {
			stmt.ObjectType = p.cur.Str
		}
		p.advance()
	}

	// Object names
	for {
		name := p.parseObjectName()
		if name != nil {
			stmt.Objects = append(stmt.Objects, name)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	// FORCE
	if p.cur.Type == kwFORCE {
		stmt.Force = true
		p.advance()
	}

	stmt.Loc.End = p.pos()
	return stmt
}
