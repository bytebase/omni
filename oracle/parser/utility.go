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
func (p *Parser) parseLockTableStmt() (nodes.StmtNode, error) {
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
		var parseErr1151 error
		item.Table, parseErr1151 = p.parseObjectName()
		if parseErr1151 !=

			// Optional partition_extension_clause:
			//   PARTITION ( partition_name )
			//   PARTITION FOR ( partition_key_value )
			//   SUBPARTITION ( subpartition_name )
			//   SUBPARTITION FOR ( subpartition_key_value )
			nil {
			return nil, parseErr1151
		}

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

		item.Loc.End = p.prev.End
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
		var parseErr1152 error
		stmt.Wait, parseErr1152 = p.parseExpr()
		if parseErr1152 != nil {
			return nil, parseErr1152
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseCallStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume CALL

	stmt := &nodes.CallStmt{
		Args: &nodes.List{},
		Loc:  nodes.Loc{Start: start},
	}
	var parseErr1153 error

	stmt.Name, parseErr1153 = p.parseObjectName()
	if parseErr1153 !=

		// Arguments
		nil {
		return nil, parseErr1153
	}

	if p.cur.Type == '(' {
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg, parseErr1154 := p.parseExpr()
			if parseErr1154 != nil {
				return nil, parseErr1154
			}
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
		var parseErr1155 error
		stmt.Into, parseErr1155 = p.parseExpr()
		if parseErr1155 != nil {
			return nil, parseErr1155
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRenameStmt parses a RENAME statement.
//
// BNF: oracle/parser/bnf/RENAME.bnf
//
//	RENAME old_name TO new_name ;
func (p *Parser) parseRenameStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume RENAME

	stmt := &nodes.RenameStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr1156 error

	stmt.OldName, parseErr1156 = p.parseObjectName()
	if parseErr1156 != nil {
		return nil, parseErr1156
	}

	if p.cur.Type == kwTO {
		p.advance()
	}
	var parseErr1157 error

	stmt.NewName, parseErr1157 = p.parseObjectName()
	if parseErr1157 != nil {
		return nil, parseErr1157
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTruncateStmt parses a TRUNCATE TABLE statement.
//
// BNF: oracle/parser/bnf/TRUNCATE-TABLE.bnf
//
//	TRUNCATE TABLE [ schema. ] table_name
//	    [ { PRESERVE | PURGE } MATERIALIZED VIEW LOG ]
//	    [ { DROP STORAGE | DROP ALL STORAGE | REUSE STORAGE } ]
//	    [ CASCADE ] ;
func (p *Parser) parseTruncateStmt() (nodes.StmtNode, error) {
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
	} else {
		return nil, p.syntaxErrorAtCur()
	}
	var parseErr1158 error

	// Parse table/cluster name
	stmt.Table, parseErr1158 = p.parseObjectName()
	if parseErr1158 !=

		// Parse optional clauses
		nil {
		return nil, parseErr1158
	}
	if stmt.Table == nil || stmt.Table.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

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

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseAnalyzeStmt() (nodes.StmtNode, error) {
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
	var parseErr1159 error

	// Object name (possibly schema-qualified)
	stmt.Table, parseErr1159 = p.parseObjectName()
	if parseErr1159 !=

		// Skip optional partition_extension_clause:
		//   PARTITION ( name ) | SUBPARTITION ( name )
		nil {
		return nil, parseErr1159
	}

	if p.isIdentLike() && (p.cur.Str == "PARTITION" || p.cur.Str == "SUBPARTITION") {
		p.advance() // consume PARTITION/SUBPARTITION
		if p.cur.Type == '(' {
			p.advance()
			parseDiscard1161, // consume (
				parseErr1160 := p.parseExpr()
			_ = parseDiscard1161
			if parseErr1160 != nil {
				return nil,

					// consume )
					parseErr1160
			}
			if p.cur.Type == ')' {
				p.advance()
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
				var parseErr1162 error
				stmt.SampleValue, parseErr1162 = p.parseIntValue()
				if parseErr1162 != nil {
					return nil, parseErr1162
				}
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
				p.advance()
				var // consume INTO
				parseErr1163 error
				stmt.IntoTable, parseErr1163 = p.parseObjectName()
				if parseErr1163 != nil {
					return nil, parseErr1163

					// LIST CHAINED ROWS [INTO table]
				}
			}
		}

	case p.cur.Type == kwLIST:

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
			p.advance()
			var // consume INTO
			parseErr1164 error
			stmt.IntoTable, parseErr1164 = p.parseObjectName()
			if parseErr1164 != nil {
				return nil, parseErr1164
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseExplainPlanStmt parses an EXPLAIN PLAN statement.
//
// BNF: oracle/parser/bnf/EXPLAIN-PLAN.bnf
//
//	EXPLAIN PLAN
//	    [ SET STATEMENT_ID = string ]
//	    [ INTO [ schema. ] table [ @dblink ] ]
//	    FOR statement ;
func (p *Parser) parseExplainPlanStmt() (nodes.StmtNode, error) {
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
		var parseErr1165 error
		stmt.Into, parseErr1165 = p.parseObjectName()
		if parseErr1165 != nil {
			return nil, parseErr1165
		}
	}

	// FOR statement
	if p.cur.Type == kwFOR {
		p.advance()
		var parseErr1166 error
		stmt.Statement, parseErr1166 = p.parseStmt()
		if parseErr1166 != nil {
			return nil, parseErr1166
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseFlashbackTableStmt() (nodes.StmtNode, error) {
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
	first, parseErr1167 := p.parseObjectName()
	if parseErr1167 != nil {
		return nil, parseErr1167
	}
	stmt.Tables = append(stmt.Tables, first)
	stmt.Table = first
	for p.cur.Type == ',' {
		p.advance() // consume ','
		t, parseErr1168 := p.parseObjectName()
		if parseErr1168 != nil {
			return nil, parseErr1168
		}
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
				var parseErr1169 error
				stmt.Rename, parseErr1169 = p.parseIdentifier()
				if parseErr1169 != nil {
					return nil,

						// TO BEFORE SCN expr | TO BEFORE TIMESTAMP expr
						parseErr1169
				}
			}
		} else {

			stmt.Before = true
			switch p.cur.Type {
			case kwSCN:
				p.advance()
				var parseErr1170 error
				stmt.ToSCN, parseErr1170 = p.parseExpr()
				if parseErr1170 != nil {
					return nil, parseErr1170
				}
			case kwTIMESTAMP:
				p.advance()
				var parseErr1171 error
				stmt.ToTimestamp, parseErr1171 = p.parseExpr()
				if parseErr1171 != nil {
					return nil, parseErr1171
				}
			}
		}
	} else {
		switch p.cur.Type {
		case kwSCN:
			p.advance()
			var parseErr1172 error
			stmt.ToSCN, parseErr1172 = p.parseExpr()
			if parseErr1172 != nil {
				return nil, parseErr1172
			}
		case kwTIMESTAMP:
			p.advance()
			var parseErr1173 error
			stmt.ToTimestamp, parseErr1173 = p.parseExpr()
			if parseErr1173 !=

				// TO RESTORE POINT restore_point_name
				nil {
				return nil, parseErr1173
			}
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseFlashbackDatabaseStmt() (nodes.StmtNode, error) {
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
		var parseErr1174 error
		stmt.DatabaseName, parseErr1174 = p.parseObjectName()
		if parseErr1174 !=

			// TO
			nil {
			return nil, parseErr1174
		}
	}

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
		var parseErr1175 error
		stmt.ToSCN, parseErr1175 = p.parseExpr()
		if parseErr1175 != nil {
			return nil, parseErr1175
		}
	case kwTIMESTAMP:
		p.advance()
		var parseErr1176 error
		stmt.ToTimestamp, parseErr1176 = p.parseExpr()
		if parseErr1176 != nil {
			return nil, parseErr1176
		}
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parsePurgeStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume PURGE

	stmt := &nodes.PurgeStmt{
		Loc: nodes.Loc{Start: start},
	}

	switch p.cur.Type {
	case kwTABLE:
		stmt.ObjectType = nodes.OBJECT_TABLE
		p.advance()
		var parseErr1177 error
		stmt.Name, parseErr1177 = p.parseObjectName()
		if parseErr1177 != nil {
			return nil, parseErr1177
		}
	case kwINDEX:
		stmt.ObjectType = nodes.OBJECT_INDEX
		p.advance()
		var parseErr1178 error
		stmt.Name, parseErr1178 = p.parseObjectName()
		if parseErr1178 != nil {
			return nil, parseErr1178
		}
	case kwTABLESPACE:
		stmt.ObjectType = nodes.OBJECT_TABLESPACE
		p.advance()
		var parseErr1179 error
		stmt.Name, parseErr1179 = p.parseObjectName()
		if parseErr1179 !=

			// RECYCLEBIN or DBA_RECYCLEBIN (parsed as identifiers)
			nil {
			return nil, parseErr1179
		}
	default:

		if p.isIdentLike() {
			ident := p.cur.Str
			p.advance()
			stmt.Name = &nodes.ObjectName{
				Name: ident,
				Loc:  nodes.Loc{Start: start, End: p.prev.End},
			}
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseAuditStmt() (nodes.StmtNode, error) {
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
					var parseErr1180 error
					stmt.WithRoles, parseErr1180 = p.parseIdentListForAudit()
					if parseErr1180 != nil {

						// BY user [, user]...
						return nil, parseErr1180
					}
				} else {
					var parseErr1181 error

					stmt.ByUsers, parseErr1181 = p.parseIdentListForAudit()
					if parseErr1181 != nil {
						return nil, parseErr1181
					}
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
					var parseErr1182 error
					stmt.WithRoles, parseErr1182 = p.parseIdentListForAudit()
					if parseErr1182 != nil {
						return nil, parseErr1182
					}
					stmt.WithRoleExcept = true
				} else {
					var parseErr1183 error
					// EXCEPT user [, user]...
					stmt.ExceptUsers, parseErr1183 = p.parseIdentListForAudit()
					if parseErr1183 != nil {
						return nil, parseErr1183
					}
				}
			} else if p.cur.Type == kwWHENEVER {
				var parseErr1184 error
				stmt.When, parseErr1184 = p.parseWheneverClause()
				if parseErr1184 != nil {
					return nil, parseErr1184
				}
			} else {
				break
			}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
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
			var parseErr1185 error
			stmt.ContextAttrs, parseErr1185 = p.parseIdentListForAudit()
			if parseErr1185 != nil {
				return nil, parseErr1185
			}
		}
		if p.cur.Type == kwBY {
			p.advance()
			var parseErr1186 error
			stmt.ByUsers, parseErr1186 = p.parseIdentListForAudit()
			if parseErr1186 != nil {
				return nil, parseErr1186
			}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	var parseErr1187 error

	// Traditional: parse audit actions
	stmt.Actions, parseErr1187 = p.parseAuditActions()
	if parseErr1187 !=

		// ON clause or auditing_by_clause
		nil {
		return nil, parseErr1187
	}

	if p.cur.Type == kwON {
		p.advance()
		parseErr1188 := p.parseTraditionalAuditOnClause(stmt)
		if parseErr1188 != nil {
			return nil, parseErr1188
		}
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
			var parseErr1189 error
			// BY user [, user]...
			stmt.ByUsers2, parseErr1189 = p.parseIdentListForAudit()
			if parseErr1189 !=

				// WHENEVER [NOT] SUCCESSFUL
				nil {
				return nil, parseErr1189
			}
		}
	}

	if p.cur.Type == kwWHENEVER {
		var parseErr1190 error
		stmt.When, parseErr1190 = p.parseWheneverClause()
		if parseErr1190 !=

			// CONTAINER = { CURRENT | ALL }
			nil {
			return nil, parseErr1190
		}
	}

	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		var parseErr1191 error
		stmt.ContainerAll, parseErr1191 = p.parseContainerClause()
		if parseErr1191 != nil {
			return nil, parseErr1191
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
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
func (p *Parser) parseNoauditStmt() (nodes.StmtNode, error) {
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
			var parseErr1192 error
			stmt.ByUsers, parseErr1192 = p.parseIdentListForAudit()
			if parseErr1192 != nil {
				return nil, parseErr1192
			}
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				var parseErr1193 error
				stmt.WithRoles, parseErr1193 = p.parseIdentListForAudit()
				if parseErr1193 != nil {
					return nil, parseErr1193
				}
			}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
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
			var parseErr1194 error
			stmt.ContextAttrs, parseErr1194 = p.parseIdentListForAudit()
			if parseErr1194 != nil {
				return nil, parseErr1194
			}
		}
		if p.cur.Type == kwBY {
			p.advance()
			var parseErr1195 error
			stmt.ByUsers, parseErr1195 = p.parseIdentListForAudit()
			if parseErr1195 != nil {
				return nil, parseErr1195
			}
			if p.cur.Type == kwWITH {
				p.advance()
				if p.cur.Type == kwROLE {
					p.advance()
				}
				var parseErr1196 error
				stmt.WithRoles, parseErr1196 = p.parseIdentListForAudit()
				if parseErr1196 != nil {
					return nil, parseErr1196
				}
			}
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	var parseErr1197 error

	// Traditional: parse audit actions
	stmt.Actions, parseErr1197 = p.parseAuditActions()
	if parseErr1197 !=

		// ON clause
		nil {
		return nil, parseErr1197
	}

	if p.cur.Type == kwON {
		p.advance()
		parseErr1198 := p.parseTraditionalNoauditOnClause(stmt)
		if parseErr1198 != nil {
			return nil, parseErr1198
		}
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
		var parseErr1199 error
		stmt.ByUsers2, parseErr1199 = p.parseIdentListForAudit()
		if parseErr1199 !=

			// WHENEVER [NOT] SUCCESSFUL
			nil {
			return nil, parseErr1199
		}
	}

	if p.cur.Type == kwWHENEVER {
		var parseErr1200 error
		stmt.When, parseErr1200 = p.parseWheneverClause()
		if parseErr1200 !=

			// CONTAINER = { CURRENT | ALL }
			nil {
			return nil, parseErr1200
		}
	}

	if p.isIdentLikeStr("CONTAINER") {
		p.advance()
		var parseErr1201 error
		stmt.ContainerAll, parseErr1201 = p.parseContainerClause()
		if parseErr1201 != nil {
			return nil, parseErr1201
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseTraditionalAuditOnClause parses the ON clause for traditional AUDIT.
// Called after ON has been consumed.
func (p *Parser) parseTraditionalAuditOnClause(stmt *nodes.AuditStmt) error {
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
		var parseErr1202 error
		stmt.Object, parseErr1202 = p.parseObjectName()
		if parseErr1202 != nil {
			return parseErr1202
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
		var parseErr1203 error
		stmt.Object, parseErr1203 = p.parseObjectName()
		if parseErr1203 != nil {
			return parseErr1203
		}
	} else {
		var parseErr1204 error
		stmt.Object, parseErr1204 = p.parseObjectName()
		if parseErr1204 !=

			// parseTraditionalNoauditOnClause parses the ON clause for traditional NOAUDIT.
			// Called after ON has been consumed.
			nil {
			return parseErr1204
		}
	}
	return nil
}

func (p *Parser) parseTraditionalNoauditOnClause(stmt *nodes.NoauditStmt) error {
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
		var parseErr1205 error
		stmt.Object, parseErr1205 = p.parseObjectName()
		if parseErr1205 != nil {
			return parseErr1205
		}
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
		var parseErr1206 error
		stmt.Object, parseErr1206 = p.parseObjectName()
		if parseErr1206 !=

			// parseAuditActions collects audit action identifiers separated by commas.
			nil {
			return parseErr1206
		}
	}
	return nil
}

func (p *Parser) parseAuditActions() ([]string, error) {
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
	return actions, nil
}

// parseWheneverClause parses WHENEVER [NOT] SUCCESSFUL.
func (p *Parser) parseWheneverClause() (string, error) {
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
	return result, nil
}

// parseIdentListForAudit parses a comma-separated list of identifiers,
// stopping at clause keywords (BY, EXCEPT, WHENEVER, WITH, CONTAINER, ;, EOF).
func (p *Parser) parseIdentListForAudit() ([]string, error) {
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
	return list, nil
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
func (p *Parser) parseAssociateStatisticsStmt() (nodes.StmtNode, error) {
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
		name, parseErr1207 := p.parseObjectName()
		if parseErr1207 != nil {
			return nil, parseErr1207
		}
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
		var parseErr1208 error
		stmt.Using, parseErr1208 = p.parseObjectName()
		if parseErr1208 != nil {
			return nil, parseErr1208
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDisassociateStatisticsStmt parses a DISASSOCIATE STATISTICS statement.
//
// BNF: oracle/parser/bnf/DISASSOCIATE-STATISTICS.bnf
//
//	DISASSOCIATE STATISTICS
//	    FROM { COLUMNS | FUNCTIONS | PACKAGES | TYPES | INDEXES | INDEXTYPES }
//	    [ schema. ] object_name [, [ schema. ] object_name ]...
//	    [ FORCE ] ;
func (p *Parser) parseDisassociateStatisticsStmt() (nodes.StmtNode, error) {
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
		name, parseErr1209 := p.parseObjectName()
		if parseErr1209 != nil {
			return nil, parseErr1209
		}
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

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
