package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// CREATE DATABASE
// ---------------------------------------------------------------------------

// parseCreateDatabaseStmt parses
//
//	CREATE [ OR REPLACE | OR ALTER ] [ TRANSIENT ] DATABASE [ IF NOT EXISTS ] <name>
//	  [ CLONE <source> [ { AT | BEFORE } ( ... ) ]
//	      [ IGNORE TABLES WITH INSUFFICIENT DATA RETENTION ] [ IGNORE HYBRID TABLES ]
//	  | FROM SHARE <provider>.<share>
//	  | FROM LISTING '<listing>'
//	  | FROM BACKUP SET <set> IDENTIFIER '<id>'
//	  | AS REPLICA OF <account>.<database> ]
//	  [ <property> = <value> ... ] [ [ WITH ] TAG ( ... ) ] [ WITH CONTACT ( ... ) ]
//
// per https://docs.snowflake.com/en/sql-reference/sql/create-database and
// https://docs.snowflake.com/en/sql-reference/sql/create-database-catalog-linked
// (whose LINKED_CATALOG = ( ... ) is one more KEY = value property). The CREATE
// keyword and optional OR REPLACE / TRANSIENT modifiers have already been
// consumed by parseCreateStmt; start is the Loc of the CREATE token.
func (p *Parser) parseCreateDatabaseStmt(start ast.Loc, orReplace, orAlter, transient bool) (ast.Node, error) {
	p.advance() // consume DATABASE

	stmt := &ast.CreateDatabaseStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Transient: transient,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional source: CLONE, FROM { SHARE | LISTING | BACKUP SET }, or
	// AS REPLICA OF.
	switch {
	case p.cur.Type == kwCLONE:
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		stmt.IgnoreInsufficientRetention, stmt.IgnoreHybridTables, err = p.parseCloneIgnoreClauses()
		if err != nil {
			return nil, err
		}

	case p.cur.Type == kwFROM && p.peekNext().Type == kwSHARE:
		p.advance() // consume FROM
		p.advance() // consume SHARE
		share, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.FromShare = share

	case p.cur.Type == kwFROM && p.peekNext().Type == kwLISTING:
		p.advance() // consume FROM
		p.advance() // consume LISTING
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		listing := tok.Str
		stmt.FromListing = &listing

	case p.cur.Type == kwFROM && tokIsWord(p.peekNext(), "BACKUP"):
		set, id, err := p.parseFromBackupSet()
		if err != nil {
			return nil, err
		}
		stmt.FromBackupSet = set
		stmt.BackupID = id

	case p.cur.Type == kwAS && p.peekNext().Type == kwREPLICA:
		p.advance() // consume AS
		p.advance() // consume REPLICA
		if _, err := p.expect(kwOF); err != nil {
			return nil, err
		}
		primary, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.ReplicaOf = primary
	}

	// Properties, [WITH] TAG (...) and WITH CONTACT (...), in any order.
	tags, err := p.parseCreateDBSchemaTail(&stmt.Props)
	if err != nil {
		return nil, err
	}
	stmt.Tags = tags

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseFromBackupSet parses FROM BACKUP SET <set> IDENTIFIER '<id>' (CREATE
// DATABASE / CREATE SCHEMA, restoring from a backup:
// https://docs.snowflake.com/en/sql-reference/sql/create-database). On entry
// cur is FROM and the next token is BACKUP.
func (p *Parser) parseFromBackupSet() (*ast.ObjectName, string, error) {
	p.advance() // consume FROM
	p.advance() // consume BACKUP
	if _, err := p.expect(kwSET); err != nil {
		return nil, "", err
	}
	set, err := p.parseObjectName()
	if err != nil {
		return nil, "", err
	}
	if err := p.expectWord("IDENTIFIER"); err != nil {
		return nil, "", err
	}
	tok, err := p.expect(tokString)
	if err != nil {
		return nil, "", err
	}
	return set, tok.Str, nil
}

// parseCloneIgnoreClauses parses the optional clauses that follow a database or
// schema CLONE source:
//
//	[ IGNORE TABLES WITH INSUFFICIENT DATA RETENTION ] [ IGNORE HYBRID TABLES ]
//
// per https://docs.snowflake.com/en/sql-reference/sql/create-database and
// https://docs.snowflake.com/en/sql-reference/sql/create-schema.
func (p *Parser) parseCloneIgnoreClauses() (insufficient, hybrid bool, err error) {
	if p.cur.Type == kwIGNORE && p.peekNext().Type == kwTABLES {
		p.advance() // consume IGNORE
		p.advance() // consume TABLES
		for _, w := range []string{"WITH", "INSUFFICIENT", "DATA", "RETENTION"} {
			if err := p.expectWord(w); err != nil {
				return false, false, err
			}
		}
		insufficient = true
	}
	if p.cur.Type == kwIGNORE && tokIsWord(p.peekNext(), "HYBRID") {
		p.advance() // consume IGNORE
		p.advance() // consume HYBRID
		if _, err := p.expect(kwTABLES); err != nil {
			return false, false, err
		}
		hybrid = true
	}
	return insufficient, hybrid, nil
}

// parseCreateDBSchemaTail parses the trailing clauses of CREATE DATABASE /
// CREATE SCHEMA: properties, [ WITH ] TAG ( ... ) and WITH CONTACT ( ... ). The
// docs list properties on both sides of the TAG and CONTACT clauses (e.g.
// OBJECT_VISIBILITY after WITH CONTACT), so the three may come in any order.
// Each iteration consumes at least one token or ends the loop.
func (p *Parser) parseCreateDBSchemaTail(props *ast.DBSchemaProps) ([]*ast.TagAssignment, error) {
	var tags []*ast.TagAssignment
	for {
		if err := p.parseDBSchemaPropsInto(props, false); err != nil {
			return nil, err
		}
		switch {
		case p.cur.Type == kwTAG || (p.cur.Type == kwWITH && p.peekNext().Type == kwTAG):
			if p.cur.Type == kwWITH {
				p.advance() // consume WITH
			}
			t, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			tags = append(tags, t...)
		case p.cur.Type == kwWITH && tokIsWord(p.peekNext(), "CONTACT"):
			p.advance() // consume WITH
			p.advance() // consume CONTACT
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			contacts, err := p.parseContactAssignments()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			props.Contacts = append(props.Contacts, contacts...)
		default:
			return tags, nil
		}
	}
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA
// ---------------------------------------------------------------------------

// parseCreateSchemaStmt parses
//
//	CREATE [ OR REPLACE | OR ALTER ] [ TRANSIENT ] SCHEMA [ IF NOT EXISTS ] <name>
//	  [ CLONE <source> [ { AT | BEFORE } ( ... ) ]
//	      [ IGNORE TABLES WITH INSUFFICIENT DATA RETENTION ] [ IGNORE HYBRID TABLES ]
//	  | FROM BACKUP SET <set> IDENTIFIER '<id>' ]
//	  [ WITH MANAGED ACCESS ]
//	  [ <property> = <value> ... ] [ [ WITH ] TAG ( ... ) ] [ WITH CONTACT ( ... ) ]
//
// per https://docs.snowflake.com/en/sql-reference/sql/create-schema, and
//
//	CREATE OR ALTER VERSIONED SCHEMA <name> [ WITH MANAGED ACCESS ] [ <property> = <value> ... ]
//
// per https://docs.snowflake.com/en/sql-reference/sql/create-versioned-schema
// (versioned is set when parseCreateStmt consumed VERSIONED; it never takes a
// CLONE or backup source).
func (p *Parser) parseCreateSchemaStmt(start ast.Loc, orReplace, orAlter, transient, versioned bool) (ast.Node, error) {
	p.advance() // consume SCHEMA

	stmt := &ast.CreateSchemaStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Transient: transient,
		Versioned: versioned,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional source: CLONE or FROM BACKUP SET (not for a versioned schema).
	switch {
	case versioned:
	case p.cur.Type == kwCLONE:
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		stmt.IgnoreInsufficientRetention, stmt.IgnoreHybridTables, err = p.parseCloneIgnoreClauses()
		if err != nil {
			return nil, err
		}
	case p.cur.Type == kwFROM && tokIsWord(p.peekNext(), "BACKUP"):
		set, id, err := p.parseFromBackupSet()
		if err != nil {
			return nil, err
		}
		stmt.FromBackupSet = set
		stmt.BackupID = id
	}

	// Optional WITH MANAGED ACCESS
	if p.cur.Type == kwWITH {
		if p.peekNext().Type == kwMANAGED {
			p.advance() // consume WITH
			p.advance() // consume MANAGED
			if _, err := p.expect(kwACCESS); err != nil {
				return nil, err
			}
			stmt.ManagedAccess = true
		}
	}

	// Properties, [WITH] TAG (...) and WITH CONTACT (...), in any order. A
	// versioned schema takes properties only.
	if versioned {
		if err := p.parseDBSchemaPropsInto(&stmt.Props, false); err != nil {
			return nil, err
		}
	} else {
		tags, err := p.parseCreateDBSchemaTail(&stmt.Props)
		if err != nil {
			return nil, err
		}
		stmt.Tags = tags
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER dispatch
// ---------------------------------------------------------------------------

// parseAlterStmt dispatches ALTER ... to the appropriate sub-parser.
// Handles DATABASE, SCHEMA, VIEW, and MATERIALIZED VIEW; falls back to
// unsupported for other objects.
func (p *Parser) parseAlterStmt() (ast.Node, error) {
	alterTok := p.advance() // consume ALTER

	// An object read by its documented shape (object.go) is dispatched by the
	// registry before the switch.
	if f := p.matchObjectForm(); f != nil {
		return p.parseObjectStmt("ALTER", f, alterTok.Loc, false, false)
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseAlterTableStmt()
	case kwDATABASE:
		// ALTER DATABASE ROLE ... (T4.6) vs ALTER DATABASE ... (T2.1). DATABASE ROLE
		// is disambiguated by the ROLE keyword following DATABASE.
		if p.peekNext().Type == kwROLE {
			p.advance() // consume DATABASE
			return p.parseAlterRoleStmt(true)
		}
		return p.parseAlterDatabaseStmt()
	case kwSCHEMA:
		return p.parseAlterSchemaStmt()
	case kwVIEW:
		return p.parseAlterViewStmt()
	case kwMATERIALIZED:
		p.advance() // consume MATERIALIZED
		return p.parseAlterMaterializedViewStmt()
	case kwSTAGE:
		// ALTER STAGE ... (T4.1).
		return p.parseAlterStageStmt()
	case kwFILE_FORMAT, kwFILE:
		// ALTER FILE FORMAT ... (T4.2). FILE FORMAT lexes as one FILE_FORMAT token
		// or two FILE+FORMAT tokens; both dispatch here.
		return p.parseAlterFileFormatStmt()
	case kwPIPE:
		// ALTER PIPE ... (T4.3).
		return p.parseAlterPipeStmt()
	case kwSTREAM:
		// ALTER STREAM ... (T4.3).
		return p.parseAlterStreamStmt()
	case kwTASK:
		// ALTER TASK ... (T4.3), or TASK GRAPH (task_graph.go).
		if p.startsTaskGraph() {
			return p.parseTaskOrGraph(func() (ast.Node, error) {
				return p.parseObjectStmt("ALTER", taskGraphForm, alterTok.Loc, false, false)
			}, p.parseAlterTaskStmt)
		}
		return p.parseAlterTaskStmt()
	case kwALERT:
		// ALTER ALERT ... (T4.3).
		return p.parseAlterAlertStmt()
	case kwFUNCTION:
		// ALTER FUNCTION ... (T4.5).
		return p.parseAlterFunctionStmt()
	case kwPROCEDURE:
		// ALTER PROCEDURE ... (T4.5).
		return p.parseAlterProcedureStmt()
	case kwDYNAMIC:
		// ALTER DYNAMIC TABLE ... (T4.4). The sub-parser consumes DYNAMIC + TABLE.
		return p.parseAlterDynamicTableStmt()
	case kwEXTERNAL:
		// ALTER EXTERNAL { TABLE (T4.4) | VOLUME (T4.7) }. The TABLE sub-parser
		// consumes EXTERNAL + TABLE. VOLUME is not a reserved keyword, so EXTERNAL
		// VOLUME lexes as kwEXTERNAL followed by a "VOLUME" identifier.
		if p.peekNext().Type == kwTABLE {
			return p.parseAlterExternalTableStmt()
		}
		extTok := p.advance() // consume EXTERNAL (anchors Loc.Start)
		if p.curIsWord("VOLUME") {
			p.advance() // consume VOLUME
			return p.parseAlterExternalVolumeStmt(extTok.Loc)
		}
		return p.unsupported("ALTER")
	case kwSEQUENCE:
		// ALTER SEQUENCE ... (T4.4). (ALTER EVENT TABLE goes through ALTER TABLE.)
		return p.parseAlterSequenceStmt()
	case kwWAREHOUSE:
		// ALTER WAREHOUSE ... (gap-warehouse).
		return p.parseAlterWarehouseStmt()
	case kwSTORAGE, kwAPI, kwNOTIFICATION, kwSECURITY, kwINTEGRATION, kwRESOURCE, kwSECRET, kwCONNECTION, kwGIT:
		// ALTER account-level integration objects (T4.7): { [STORAGE] | API |
		// NOTIFICATION | SECURITY } INTEGRATION, RESOURCE MONITOR, SECRET,
		// CONNECTION, GIT REPOSITORY. A bare ALTER INTEGRATION (qualifier omitted)
		// dispatches here too. parseAlterIntegrationStmt consumes the keyword(s).
		return p.parseAlterIntegrationStmt()
	case kwTAG:
		// ALTER TAG ... (T4.9).
		return p.parseAlterTagStmt()
	case kwSEMANTIC:
		// ALTER SEMANTIC VIEW ... (T4.9). The sub-parser consumes SEMANTIC + VIEW.
		return p.parseAlterSemanticViewStmt()
	case kwFAILOVER, kwREPLICATION:
		// ALTER { FAILOVER | REPLICATION } GROUP ... (T4.8). The sub-parser consumes
		// the kind keyword + GROUP.
		return p.parseAlterReplicationGroupStmt()
	case kwACCOUNT:
		// ALTER ACCOUNT ... (T4.8). The sub-parser consumes ACCOUNT.
		return p.parseAlterAccountStmt()
	case kwSHARE:
		// ALTER SHARE ... (T4.8).
		return p.parseAlterShareStmt()
	case kwROLE:
		// ALTER ROLE ... (T4.6). DATABASE ROLE is dispatched by the kwDATABASE case.
		return p.parseAlterRoleStmt(false)
	case kwUSER:
		// ALTER USER ... (T4.6).
		return p.parseAlterUserStmt()
	case kwMASKING, kwSESSION, kwPASSWORD, kwNETWORK, kwROW:
		// ALTER NETWORK RULE ... (gap-network-rule) vs the policy objects below.
		// NETWORK is reserved but RULE is not, so NETWORK RULE lexes as kwNETWORK
		// followed by a "RULE" identifier; NETWORK POLICY lexes as kwNETWORK followed
		// by kwPOLICY. The RULE form is dispatched here before the shared policy path.
		if p.cur.Type == kwNETWORK && p.peekIsWord("RULE") {
			return p.parseAlterNetworkRuleStmt()
		}
		// ALTER SESSION { SET | UNSET } <param> ... (session parameters,
		// gap-alter-family) vs ALTER SESSION POLICY <name> ... (a policy object).
		// SESSION followed by SET/UNSET is session parameters; SESSION followed by
		// POLICY stays on the policy path below.
		if p.cur.Type == kwSESSION {
			if next := p.peekNext().Type; next == kwSET || next == kwUNSET {
				return p.parseAlterSessionStmt()
			}
		}
		// ALTER { MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK } POLICY ...
		// (T4.6). Only routed to the policy parser when the kind keyword is actually
		// followed by POLICY (and, for ROW, ACCESS POLICY); otherwise a same-prefix
		// non-policy statement (e.g. ALTER SESSION SET ...) is left unsupported here.
		if p.startsPolicyKeyword() {
			return p.parseAlterPolicyDispatch()
		}
		return p.unsupported("ALTER")
	default:
		// ALTER ICEBERG TABLE: ICEBERG is not a keyword.
		if p.startsAlterIcebergTable() {
			return p.parseAlterIcebergTableStmt()
		}
		// AUTHENTICATION is not a reserved keyword, so ALTER AUTHENTICATION POLICY
		// lexes as an "AUTHENTICATION" identifier followed by POLICY (T4.6).
		if p.startsPolicyKeyword() {
			return p.parseAlterPolicyDispatch()
		}
		return p.unsupported("ALTER")
	}
}

// parseAlterPolicyDispatch consumes the policy-kind keyword(s) + the POLICY
// keyword and delegates to parseAlterPolicyStmt (T4.6). On entry cur is the
// first policy-kind word; the caller has confirmed startsPolicyKeyword.
func (p *Parser) parseAlterPolicyDispatch() (ast.Node, error) {
	startLoc := p.cur.Loc // policy-kind keyword anchors Loc.Start (ALTER convention)
	kind, err := p.consumePolicyKeywords()
	if err != nil {
		return nil, err
	}
	return p.parseAlterPolicyStmt(startLoc, kind)
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

// parseAlterDatabaseStmt parses ALTER DATABASE ... (all action variants) per
// https://docs.snowflake.com/en/sql-reference/sql/alter-database and
// https://docs.snowflake.com/en/sql-reference/sql/alter-database-catalog-linked.
// The ALTER keyword has already been consumed.
func (p *Parser) parseAlterDatabaseStmt() (ast.Node, error) {
	altTok := p.advance() // consume DATABASE
	stmt := &ast.AlterDatabaseStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDBRename
		stmt.NewName = newName

	case kwSWAP:
		// SWAP WITH other_name
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		otherName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDBSwap
		stmt.NewName = otherName

	case kwSET:
		p.advance() // consume SET
		switch {
		case p.startsFeaturePolicy():
			// SET FEATURE POLICY <name> [ FORCE ]
			p.advance() // consume FEATURE
			p.advance() // consume POLICY
			policy, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBSetFeaturePolicy
			stmt.FeaturePolicy = policy
			if p.cur.Type == kwFORCE {
				p.advance() // consume FORCE
				stmt.Force = true
			}
		case p.cur.Type == kwTAG:
			tags, err := p.parseDBSchemaSetTags()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBSetTag
			stmt.Tags = tags
		default:
			props, err := p.parseDBSchemaSetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBSet
			stmt.SetProps = props
		}

	case kwUNSET:
		p.advance() // consume UNSET
		switch {
		case p.startsFeaturePolicy():
			// UNSET FEATURE POLICY
			p.advance() // consume FEATURE
			p.advance() // consume POLICY
			stmt.Action = ast.AlterDBUnsetFeaturePolicy
		case p.cur.Type == kwTAG:
			names, err := p.parseDBSchemaUnsetTags()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBUnsetTag
			stmt.UnsetTags = names
		default:
			props, err := p.parseDBSchemaUnsetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDBUnset
			stmt.UnsetProps = props
		}

	case kwSUSPEND, kwRESUME:
		// { SUSPEND | RESUME } DISCOVERY (catalog-linked database)
		action := ast.AlterDBSuspendDiscovery
		if p.cur.Type == kwRESUME {
			action = ast.AlterDBResumeDiscovery
		}
		p.advance() // consume SUSPEND / RESUME
		if err := p.expectWord("DISCOVERY"); err != nil {
			return nil, err
		}
		stmt.Action = action

	case kwUPDATE:
		// UPDATE LINKED_CATALOG ... (catalog-linked database)
		p.advance() // consume UPDATE
		upd, err := p.parseLinkedCatalogUpdate()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDBUpdateLinkedCatalog
		stmt.LinkedCatalog = upd

	case kwENABLE:
		p.advance() // consume ENABLE
		switch p.cur.Type {
		case kwREPLICATION:
			// ENABLE REPLICATION TO ACCOUNTS account_list [IGNORE EDITION CHECK]
			p.advance() // consume REPLICATION
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwACCOUNTS); err != nil {
				return nil, err
			}
			// Skip account identifiers until end / IGNORE
			for p.cur.Type != tokEOF && p.cur.Type != kwIGNORE && p.cur.Type != ';' {
				p.advance()
			}
			// Optional IGNORE EDITION CHECK
			if p.cur.Type == kwIGNORE {
				p.advance() // consume IGNORE
				p.advance() // consume EDITION
				p.advance() // consume CHECK
			}
			stmt.Action = ast.AlterDBEnableReplication
		case kwFAILOVER:
			// ENABLE FAILOVER TO ACCOUNTS account_list
			p.advance() // consume FAILOVER
			if _, err := p.expect(kwTO); err != nil {
				return nil, err
			}
			if _, err := p.expect(kwACCOUNTS); err != nil {
				return nil, err
			}
			// Skip account identifiers
			for p.cur.Type != tokEOF && p.cur.Type != ';' {
				p.advance()
			}
			stmt.Action = ast.AlterDBEnableFailover
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwDISABLE:
		p.advance() // consume DISABLE
		switch p.cur.Type {
		case kwREPLICATION:
			// DISABLE REPLICATION [TO ACCOUNTS account_list]
			p.advance() // consume REPLICATION
			if p.cur.Type == kwTO {
				p.advance() // consume TO
				if _, err := p.expect(kwACCOUNTS); err != nil {
					return nil, err
				}
				for p.cur.Type != tokEOF && p.cur.Type != ';' {
					p.advance()
				}
			}
			stmt.Action = ast.AlterDBDisableReplication
		case kwFAILOVER:
			// DISABLE FAILOVER [TO ACCOUNTS account_list]
			p.advance() // consume FAILOVER
			if p.cur.Type == kwTO {
				p.advance() // consume TO
				if _, err := p.expect(kwACCOUNTS); err != nil {
					return nil, err
				}
				for p.cur.Type != tokEOF && p.cur.Type != ';' {
					p.advance()
				}
			}
			stmt.Action = ast.AlterDBDisableFailover
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwREFRESH:
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterDBRefresh

	case kwPRIMARY:
		p.advance() // consume PRIMARY
		stmt.Action = ast.AlterDBPrimary

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsFeaturePolicy reports whether cur begins FEATURE POLICY (ALTER DATABASE
// { SET | UNSET } FEATURE POLICY). FEATURE is not a reserved keyword.
func (p *Parser) startsFeaturePolicy() bool {
	return tokIsWord(p.cur, "FEATURE") && p.peekNext().Type == kwPOLICY
}

// parseLinkedCatalogUpdate parses the body of ALTER DATABASE ... UPDATE
// LINKED_CATALOG (cur is LINKED_CATALOG):
//
//	ADD ( '<namespace>' [ , ... ] ) TO { ALLOWED_NAMESPACES | BLOCKED_NAMESPACES }
//	REMOVE ( '<namespace>' [ , ... ] ) FROM { ALLOWED_NAMESPACES | BLOCKED_NAMESPACES }
//	UNSET { ALLOWED_NAMESPACES | BLOCKED_NAMESPACES }
//	SET <key> = <value> [ ... ]
//
// per https://docs.snowflake.com/en/sql-reference/sql/alter-database-catalog-linked.
// The namespace-list name and the SET keys are open-ended name words.
func (p *Parser) parseLinkedCatalogUpdate() (*ast.LinkedCatalogUpdate, error) {
	if err := p.expectWord("LINKED_CATALOG"); err != nil {
		return nil, err
	}
	upd := &ast.LinkedCatalogUpdate{}
	switch p.cur.Type {
	case kwADD, kwREMOVE:
		upd.Op = "ADD"
		joiner := kwTO
		if p.cur.Type == kwREMOVE {
			upd.Op = "REMOVE"
			joiner = kwFROM
		}
		p.advance() // consume ADD / REMOVE
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		names, err := p.parseTagStringList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		upd.Namespaces = names
		if _, err := p.expect(joiner); err != nil {
			return nil, err
		}
		target, err := p.parseDBPropName()
		if err != nil {
			return nil, err
		}
		upd.Target = target
	case kwUNSET:
		p.advance() // consume UNSET
		upd.Op = "UNSET"
		target, err := p.parseDBPropName()
		if err != nil {
			return nil, err
		}
		upd.Target = target
	case kwSET:
		p.advance() // consume SET
		upd.Op = "SET"
		for i := 0; p.startsDBSchemaOption(); i++ {
			if i >= dbSchemaOptionCap {
				return nil, p.syntaxErrorAtCur()
			}
			opt, err := p.parseDBSchemaOption()
			if err != nil {
				return nil, err
			}
			upd.Options = append(upd.Options, opt)
		}
		if len(upd.Options) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return upd, nil
}

// ---------------------------------------------------------------------------
// ALTER SCHEMA
// ---------------------------------------------------------------------------

// parseAlterSchemaStmt parses ALTER SCHEMA ... (all action variants) per
// https://docs.snowflake.com/en/sql-reference/sql/alter-schema.
// The ALTER keyword has already been consumed.
func (p *Parser) parseAlterSchemaStmt() (ast.Node, error) {
	altTok := p.advance() // consume SCHEMA
	stmt := &ast.AlterSchemaStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Action branch
	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO new_name
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaRename
		stmt.NewName = newName

	case kwSWAP:
		// SWAP WITH other_name
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		otherName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaSwap
		stmt.NewName = otherName

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseDBSchemaSetTags()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaSetTag
			stmt.Tags = tags
		} else {
			props, err := p.parseDBSchemaSetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaSet
			stmt.SetProps = props
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseDBSchemaUnsetTags()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseDBSchemaUnsetProps()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterSchemaUnset
			stmt.UnsetProps = props
		}

	case kwENABLE:
		// ENABLE MANAGED ACCESS
		p.advance() // consume ENABLE
		if _, err := p.expect(kwMANAGED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaEnableManagedAccess

	case kwDISABLE:
		// DISABLE MANAGED ACCESS
		p.advance() // consume DISABLE
		if _, err := p.expect(kwMANAGED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSchemaDisableManagedAccess

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER SESSION (session parameters)
// ---------------------------------------------------------------------------

// sessionParamCap bounds the open-ended `<param> = <value>` run of an
// ALTER SESSION SET so a malformed input can never spin.
const sessionParamCap = 1024

// parseAlterSessionStmt parses ALTER SESSION { SET | UNSET } <session_parameter>
// ... (session parameters). The ALTER keyword has already been consumed; cur is
// SESSION. ALTER SESSION POLICY ... is routed to the policy parser by the
// dispatch in parseAlterStmt, so SESSION here is always followed by SET/UNSET.
func (p *Parser) parseAlterSessionStmt() (ast.Node, error) {
	sessTok := p.advance() // consume SESSION (anchors Loc.Start at the object keyword)
	stmt := &ast.AlterSessionStmt{Loc: ast.Loc{Start: sessTok.Loc.Start}}

	switch p.cur.Type {
	case kwSET:
		p.advance() // consume SET
		opts, err := p.parseSessionParamOptions()
		if err != nil {
			return nil, err
		}
		if len(opts) == 0 {
			// SET with nothing settable is a syntax error.
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterSessionSet
		stmt.Options = opts

	case kwUNSET:
		p.advance() // consume UNSET
		keys, err := p.parseSessionParamUnsetKeys()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterSessionUnset
		stmt.UnsetKeys = keys

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSessionParamUnsetKeys reads the comma-separated parameter-name list of an
// ALTER SESSION UNSET. Session parameter names are open-ended: most lex as
// reserved keywords (LOCK_TIMEOUT, AUTOCOMMIT, ...), so any option word is
// accepted. The names are uppercased; the loop is bounded.
func (p *Parser) parseSessionParamUnsetKeys() ([]string, error) {
	var keys []string
	for i := 0; ; i++ {
		if i >= sessionParamCap {
			return nil, p.syntaxErrorAtCur()
		}
		if !p.isOptionWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		keys = append(keys, strings.ToUpper(p.advance().Str))
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	return keys, nil
}

// parseSessionParamOptions reads the open-ended `<param> = <value>` run of an
// ALTER SESSION SET (e.g. AUTOCOMMIT = TRUE, LOCK_TIMEOUT = 3600,
// DEFAULT_NULL_ORDERING = 'LAST'). Each entry is a CopyOption so any name/value
// shape (word, number, or string literal) is captured. The loop is bounded and
// guards against a non-advancing token.
func (p *Parser) parseSessionParamOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for i := 0; p.startsCopyOption(); i++ {
		if i >= sessionParamCap {
			return nil, p.syntaxErrorAtCur()
		}
		before := p.cur.Loc.Start
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if p.cur.Loc.Start == before && p.cur.Type != tokEOF {
			// No forward progress on a non-EOF token: abort rather than spin.
			return nil, p.syntaxErrorAtCur()
		}
	}
	return opts, nil
}

// ---------------------------------------------------------------------------
// DROP dispatch
// ---------------------------------------------------------------------------

// parseDropStmt dispatches DROP ... to the appropriate sub-parser.
// (parseDropStmt lives in drop.go; it dispatches DATABASE/SCHEMA here.)

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

// parseDropDatabaseStmt parses DROP DATABASE [IF EXISTS] name [CASCADE|RESTRICT].
// The DROP keyword has already been consumed.
func (p *Parser) parseDropDatabaseStmt() (ast.Node, error) {
	dropTok := p.advance() // consume DATABASE
	stmt := &ast.DropDatabaseStmt{Loc: ast.Loc{Start: dropTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Database name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CASCADE | RESTRICT
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		stmt.Cascade = true
	case kwRESTRICT:
		p.advance()
		stmt.Restrict = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// DROP SCHEMA
// ---------------------------------------------------------------------------

// parseDropSchemaStmt parses DROP SCHEMA [IF EXISTS] name [CASCADE|RESTRICT].
// The DROP keyword has already been consumed.
func (p *Parser) parseDropSchemaStmt() (ast.Node, error) {
	dropTok := p.advance() // consume SCHEMA
	stmt := &ast.DropSchemaStmt{Loc: ast.Loc{Start: dropTok.Loc.Start}}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Schema name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CASCADE | RESTRICT
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		stmt.Cascade = true
	case kwRESTRICT:
		p.advance()
		stmt.Restrict = true
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP dispatch
// ---------------------------------------------------------------------------

// (parseUndropStmt lives in drop.go; it dispatches DATABASE/SCHEMA here.)

// ---------------------------------------------------------------------------
// UNDROP DATABASE
// ---------------------------------------------------------------------------

// parseUndropDatabaseStmt parses UNDROP DATABASE name.
// The UNDROP keyword has already been consumed.
func (p *Parser) parseUndropDatabaseStmt() (ast.Node, error) {
	undropTok := p.advance() // consume DATABASE
	stmt := &ast.UndropDatabaseStmt{Loc: ast.Loc{Start: undropTok.Loc.Start}}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP SCHEMA
// ---------------------------------------------------------------------------

// parseUndropSchemaStmt parses UNDROP SCHEMA name.
// The UNDROP keyword has already been consumed.
func (p *Parser) parseUndropSchemaStmt() (ast.Node, error) {
	undropTok := p.advance() // consume SCHEMA
	stmt := &ast.UndropSchemaStmt{Loc: ast.Loc{Start: undropTok.Loc.Start}}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Property parsing helpers
// ---------------------------------------------------------------------------

// dbSchemaOptionCap bounds an open-ended DATABASE/SCHEMA property or UNSET run.
// Real statements carry a handful of properties; a far larger run signals a
// runaway parse and aborts loudly rather than spinning.
const dbSchemaOptionCap = 1024

// parseDBSchemaSetProps parses the property list of ALTER { DATABASE | SCHEMA }
// ... SET, which also takes CONTACT <purpose> = <contact> [ , ... ] inline
// (https://docs.snowflake.com/en/sql-reference/sql/alter-database). SET with
// nothing to set is a syntax error.
func (p *Parser) parseDBSchemaSetProps() (*ast.DBSchemaProps, error) {
	props := &ast.DBSchemaProps{}
	if err := p.parseDBSchemaPropsInto(props, true); err != nil {
		return nil, err
	}
	if props.DataRetention == nil && props.MaxDataExt == nil && props.DefaultDDLCol == nil &&
		props.Comment == nil && len(props.Options) == 0 && len(props.Contacts) == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	return props, nil
}

// parseDBSchemaPropsInto populates *props in place with a run of DATABASE/SCHEMA
// properties. DATA_RETENTION_TIME_IN_DAYS, MAX_DATA_EXTENSION_TIME_IN_DAYS,
// DEFAULT_DDL_COLLATION and COMMENT keep their typed fields; every other
// property is an open-ended `KEY = value` pair (ast.CopyOption) in
// props.Options, since the documented set (EXTERNAL_VOLUME, CATALOG,
// LOG_LEVEL, OBJECT_VISIBILITY, LINKED_CATALOG = ( ... ), ...) keeps growing.
// With inlineContact (ALTER ... SET), CONTACT <purpose> = <contact> [ , ... ]
// is read into props.Contacts. The run ends at the first token that begins no
// property.
func (p *Parser) parseDBSchemaPropsInto(props *ast.DBSchemaProps, inlineContact bool) error {
	for i := 0; ; i++ {
		if i >= dbSchemaOptionCap {
			return p.syntaxErrorAtCur()
		}
		switch p.cur.Type {
		case kwDATA_RETENTION_TIME_IN_DAYS:
			p.advance() // consume DATA_RETENTION_TIME_IN_DAYS
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokInt)
			if err != nil {
				return err
			}
			v := tok.Ival
			props.DataRetention = &v
			continue

		case kwMAX_DATA_EXTENSION_TIME_IN_DAYS:
			p.advance() // consume MAX_DATA_EXTENSION_TIME_IN_DAYS
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokInt)
			if err != nil {
				return err
			}
			v := tok.Ival
			props.MaxDataExt = &v
			continue

		case kwDEFAULT_DDL_COLLATION:
			p.advance() // consume DEFAULT_DDL_COLLATION / DEFAULT_DDL_COLLATION_
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			props.DefaultDDLCol = &s
			continue

		case kwCOMMENT:
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			props.Comment = &s
			continue
		}

		switch {
		case inlineContact && p.startsInlineContact():
			p.advance() // consume CONTACT
			contacts, err := p.parseContactAssignments()
			if err != nil {
				return err
			}
			props.Contacts = append(props.Contacts, contacts...)
		case p.startsDBSchemaOption():
			opt, err := p.parseDBSchemaOption()
			if err != nil {
				return err
			}
			props.Options = append(props.Options, opt)
		default:
			return nil
		}
	}
}

// startsDBSchemaOption reports whether cur begins an open-ended `KEY = value`
// DATABASE/SCHEMA property: a name word followed by '='. WITH and TAG never
// start one; they anchor the [ WITH ] TAG and WITH CONTACT clauses.
func (p *Parser) startsDBSchemaOption() bool {
	if p.cur.Type == kwWITH || p.cur.Type == kwTAG || !isNameWordTok(p.cur) {
		return false
	}
	return p.peekNext().Type == '='
}

// parseDBSchemaOption parses one `KEY = value` DATABASE/SCHEMA property; the
// caller has confirmed startsDBSchemaOption. The value is one of
//
//	'<string>' | <number> | <name> [ . <name> ... ]
//	( '<string>' [ , ... ] )                 e.g. ALLOWED_NAMESPACES = ( 'a', 'b' )
//	( <key> = <value> [ , ... ] )            e.g. LINKED_CATALOG = ( CATALOG = 'c', ... )
//
// Unlike the COPY option groups (parseCopyOptionParen), a group here separates
// its entries with commas and every entry has a value, as in
// https://docs.snowflake.com/en/sql-reference/sql/create-database-catalog-linked.
func (p *Parser) parseDBSchemaOption() (*ast.CopyOption, error) {
	nameTok := p.advance() // consume the property name
	opt := &ast.CopyOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start},
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if err := p.parseDBSchemaOptionValue(opt); err != nil {
		return nil, err
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseDBSchemaOptionValue parses the value of a parseDBSchemaOption property
// into opt.
func (p *Parser) parseDBSchemaOptionValue(opt *ast.CopyOption) error {
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
		return nil
	case tokInt:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}
		return nil
	case tokFloat, tokReal:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}
		return nil
	case '(':
		return p.parseDBSchemaOptionParen(opt)
	}
	words, err := p.parseDBSchemaWordValue()
	if err != nil {
		return err
	}
	opt.Words = words
	return nil
}

// parseDBSchemaOptionParen parses a parenthesized property value: a non-empty
// string list or a non-empty comma-separated `KEY = value` group. cur is '('.
func (p *Parser) parseDBSchemaOptionParen(opt *ast.CopyOption) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	p.advance() // consume '('

	if p.cur.Type == tokString {
		for i := 0; ; i++ {
			if i >= dbSchemaOptionCap {
				return p.syntaxErrorAtCur()
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			opt.List = append(opt.List, &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc})
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	} else {
		for i := 0; ; i++ {
			if i >= dbSchemaOptionCap || !isNameWordTok(p.cur) {
				return p.syntaxErrorAtCur()
			}
			sub, err := p.parseDBSchemaOption()
			if err != nil {
				return err
			}
			opt.Group = append(opt.Group, sub)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // consume ','
		}
	}
	_, err := p.expect(')')
	return err
}

// parseDBSchemaWordValue parses a name-valued property (TRUE, COMPATIBLE, an
// external volume or a qualified event table): one or more dot-separated name
// parts. Unquoted parts are uppercased; a quoted part keeps its source text,
// quotes included, so its case survives.
func (p *Parser) parseDBSchemaWordValue() (string, error) {
	var b strings.Builder
	for i := 0; ; i++ {
		if i >= dbSchemaOptionCap {
			return "", p.syntaxErrorAtCur()
		}
		switch {
		case p.cur.Type == tokQuotedIdent:
			tok := p.advance()
			b.WriteString(p.srcSlice(tok.Loc.Start, tok.Loc.End))
		case isNameWordTok(p.cur):
			b.WriteString(strings.ToUpper(p.advance().Str))
		default:
			return "", p.syntaxErrorAtCur()
		}
		if p.cur.Type != '.' {
			return b.String(), nil
		}
		p.advance() // consume '.'
		b.WriteByte('.')
	}
}

// startsInlineContact reports whether cur begins the ALTER ... SET form
// CONTACT <purpose> = <contact>: the word CONTACT followed by a name word.
func (p *Parser) startsInlineContact() bool {
	return tokIsWord(p.cur, "CONTACT") && isNameWordTok(p.peekNext())
}

// parseContactAssignments parses <purpose> = <contact> [ , ... ], the body of
// WITH CONTACT ( ... ) on CREATE and of CONTACT ... on ALTER ... SET
// (https://docs.snowflake.com/en/sql-reference/sql/create-database). Each pair
// is a CopyOption whose Name is the uppercased purpose and whose Words is the
// contact name.
func (p *Parser) parseContactAssignments() ([]*ast.CopyOption, error) {
	var contacts []*ast.CopyOption
	for i := 0; ; i++ {
		if i >= dbSchemaOptionCap {
			return nil, p.syntaxErrorAtCur()
		}
		if !isNameWordTok(p.cur) {
			return nil, p.syntaxErrorAtCur()
		}
		purpose := p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		contact, err := p.parseDBSchemaWordValue()
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, &ast.CopyOption{
			Name:  strings.ToUpper(purpose.Str),
			Words: contact,
			Loc:   ast.Loc{Start: purpose.Loc.Start, End: p.prev.Loc.End},
		})
		if p.cur.Type != ',' {
			return contacts, nil
		}
		p.advance() // consume ','
	}
}

// parseDBSchemaSetTags parses SET TAG <tag> = '<value>' [ , ... ] of ALTER
// { DATABASE | SCHEMA } (the documented, unparenthesized form), and the
// parenthesized TAG ( ... ) form omni has long accepted. cur is TAG.
func (p *Parser) parseDBSchemaSetTags() ([]*ast.TagAssignment, error) {
	if p.peekNext().Type == '(' {
		return p.parseTagAssignments()
	}
	return p.parseTagSetList()
}

// parseDBSchemaUnsetTags parses UNSET TAG <tag> [ , ... ] of ALTER { DATABASE |
// SCHEMA } (the documented, unparenthesized form), and the parenthesized
// TAG ( ... ) form omni has long accepted. cur is TAG.
func (p *Parser) parseDBSchemaUnsetTags() ([]*ast.ObjectName, error) {
	if p.peekNext().Type == '(' {
		return p.parseUnsetTagList()
	}
	return p.parseUnsetTagNameList()
}

// parseDBSchemaUnsetProps parses the comma-separated property list of ALTER
// { DATABASE | SCHEMA } ... UNSET. Property names are open-ended name words,
// uppercased. CONTACT takes an optional purpose (UNSET CONTACT <purpose>),
// recorded as "CONTACT <PURPOSE>". Consumes at least one name.
func (p *Parser) parseDBSchemaUnsetProps() ([]string, error) {
	var names []string
	for i := 0; ; i++ {
		if i >= dbSchemaOptionCap {
			return nil, p.syntaxErrorAtCur()
		}
		name, err := p.parseDBPropName()
		if err != nil {
			return nil, err
		}
		if name == "CONTACT" && isNameWordTok(p.cur) {
			name += " " + strings.ToUpper(p.advance().Str)
		}
		names = append(names, name)
		if p.cur.Type != ',' {
			return names, nil
		}
		p.advance() // consume ','
	}
}

// parseDBPropName parses one property name word and returns it uppercased.
func (p *Parser) parseDBPropName() (string, error) {
	if !isNameWordTok(p.cur) {
		return "", p.syntaxErrorAtCur()
	}
	return strings.ToUpper(p.advance().Str), nil
}

// isNameWordTok reports whether tok is an unquoted name word: an identifier or
// a keyword.
func isNameWordTok(tok Token) bool {
	return tok.Type == tokIdent || tok.Type >= 700
}

// tokIsWord reports whether tok is the unquoted name word upper (an identifier
// or keyword spelled so, case-insensitively). Unlike curIsWord it never
// matches a string literal or quoted identifier with that text.
func tokIsWord(tok Token, upper string) bool {
	return isNameWordTok(tok) && strings.EqualFold(tok.Str, upper)
}

// expectWord consumes cur when it is the unquoted name word upper, and fails
// with a syntax error otherwise.
func (p *Parser) expectWord(upper string) error {
	if !tokIsWord(p.cur, upper) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseUnsetTagList parses TAG ( name [, name ...] ) for UNSET TAG clauses.
// Consumes the TAG keyword and the parenthesized name list.
func (p *Parser) parseUnsetTagList() ([]*ast.ObjectName, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var names []*ast.ObjectName
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return names, nil
}
