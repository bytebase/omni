package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// DROP statement dispatch
// ---------------------------------------------------------------------------

// parseDropStmt parses DROP <object_type> [IF EXISTS] name [CASCADE|RESTRICT].
// The DROP keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseDropStmt() (ast.Node, error) {
	dropTok := p.advance() // consume DROP
	start := dropTok.Loc

	switch p.cur.Type {
	case kwDATABASE:
		// DATABASE handled by T2.1; parseDropDatabaseStmt consumes DATABASE.
		return p.parseDropDatabaseStmt()

	case kwSCHEMA:
		// SCHEMA handled by T2.1; parseDropSchemaStmt consumes SCHEMA.
		return p.parseDropSchemaStmt()

	case kwTABLE:
		p.advance() // consume TABLE
		return p.parseDropObject(ast.DropTable, start, true, true)

	case kwVIEW:
		p.advance() // consume VIEW
		return p.parseDropObject(ast.DropView, start, true, false)

	case kwMATERIALIZED:
		// MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropMaterializedView, start, true, false)

	case kwDYNAMIC:
		// DYNAMIC TABLE — no IF EXISTS in legacy grammar
		p.advance() // consume DYNAMIC
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropDynamicTable, start, false, false)

	case kwEXTERNAL:
		// EXTERNAL TABLE
		p.advance() // consume EXTERNAL
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropExternalTable, start, true, true)

	case kwSTREAM:
		p.advance() // consume STREAM
		return p.parseDropObject(ast.DropStream, start, true, false)

	case kwTASK:
		p.advance() // consume TASK
		return p.parseDropObject(ast.DropTask, start, true, false)

	case kwSEQUENCE:
		p.advance() // consume SEQUENCE
		return p.parseDropObject(ast.DropSequence, start, true, true)

	case kwSTAGE:
		p.advance() // consume STAGE
		return p.parseDropObject(ast.DropStage, start, true, false)

	case kwFILE_FORMAT:
		p.advance() // consume FILE_FORMAT (single keyword token)
		return p.parseDropObject(ast.DropFileFormat, start, true, false)

	case kwFILE:
		// FILE FORMAT — two-word form using separate FILE and FORMAT tokens
		p.advance() // consume FILE
		if _, err := p.expect(kwFORMAT); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropFileFormat, start, true, false)

	case kwFUNCTION:
		p.advance() // consume FUNCTION
		return p.parseDropObject(ast.DropFunction, start, true, false)

	case kwPROCEDURE:
		p.advance() // consume PROCEDURE
		return p.parseDropObject(ast.DropProcedure, start, true, false)

	case kwPIPE:
		p.advance() // consume PIPE
		return p.parseDropObject(ast.DropPipe, start, true, false)

	case kwTAG:
		p.advance() // consume TAG
		return p.parseDropObject(ast.DropTag, start, true, false)

	case kwROLE:
		p.advance() // consume ROLE
		return p.parseDropObject(ast.DropRole, start, true, false)

	case kwWAREHOUSE:
		p.advance() // consume WAREHOUSE
		return p.parseDropObject(ast.DropWarehouse, start, true, false)

	case kwNETWORK:
		// DROP NETWORK RULE (gap-network-rule) and DROP NETWORK POLICY. NETWORK is
		// reserved but RULE is not, so NETWORK RULE lexes as kwNETWORK followed by a
		// "RULE" identifier; POLICY is a reserved keyword, so NETWORK POLICY lexes as
		// kwNETWORK followed by kwPOLICY.
		if p.peekIsWord("RULE") {
			p.advance() // consume NETWORK
			p.advance() // consume RULE
			return p.parseDropObject(ast.DropNetworkRule, start, true, false)
		}
		if p.peekNext().Type == kwPOLICY {
			p.advance() // consume NETWORK
			p.advance() // consume POLICY
			return p.parseDropObject(ast.DropNetworkPolicy, start, true, false)
		}
		return p.unsupportedDrop()

	case kwALERT:
		p.advance() // consume ALERT
		return p.parseDropObject(ast.DropAlert, start, true, false)

	case kwCONNECTION:
		p.advance() // consume CONNECTION
		return p.parseDropObject(ast.DropConnection, start, true, false)

	case kwFAILOVER:
		// FAILOVER GROUP
		p.advance() // consume FAILOVER
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropFailoverGroup, start, true, false)

	case kwINTEGRATION:
		p.advance() // consume INTEGRATION
		return p.parseDropObject(ast.DropIntegration, start, true, false)

	case kwMANAGED:
		// MANAGED ACCOUNT
		p.advance() // consume MANAGED
		if _, err := p.expect(kwACCOUNT); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropManagedAccount, start, true, false)

	case kwMASKING:
		// MASKING POLICY
		p.advance() // consume MASKING
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropMaskingPolicy, start, true, false)

	case kwREPLICATION:
		// REPLICATION GROUP
		p.advance() // consume REPLICATION
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropReplicationGroup, start, true, false)

	case kwRESOURCE:
		// RESOURCE MONITOR
		p.advance() // consume RESOURCE
		if _, err := p.expect(kwMONITOR); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropResourceMonitor, start, true, false)

	case kwROW:
		// ROW ACCESS POLICY
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropRowAccessPolicy, start, true, false)

	case kwSESSION:
		// SESSION POLICY
		p.advance() // consume SESSION
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropSessionPolicy, start, true, false)

	case kwSHARE:
		p.advance() // consume SHARE
		return p.parseDropObject(ast.DropShare, start, true, false)

	case kwUSER:
		p.advance() // consume USER
		return p.parseDropObject(ast.DropUser, start, true, false)

	default:
		// Emit a targeted error for recognized-but-unimplemented DROP forms,
		// or a generic error for completely unknown object types.
		return p.unsupportedDrop()
	}
}

// unsupportedDrop emits a targeted "DROP <obj> ... not yet supported" error for
// a recognized-but-unimplemented (or unknown) DROP object type at cur, and
// skips to the next statement for error recovery.
func (p *Parser) unsupportedDrop() (ast.Node, error) {
	objText := p.cur.Str
	if objText == "" {
		objText = TokenName(p.cur.Type)
	}
	err := &ParseError{
		Loc: p.cur.Loc,
		Msg: "DROP " + objText + " statement parsing is not yet supported",
	}
	p.skipToNextStatement()
	return nil, err
}

// parseDropObject is the shared helper that parses the [IF EXISTS] name
// [CASCADE|RESTRICT] tail common to most DROP forms.
//
// kind is the already-determined DropObjectKind.
// start is the Loc of the DROP token.
// ifExistsOK controls whether IF EXISTS is accepted.
// cascadeOK controls whether CASCADE/RESTRICT are accepted.
func (p *Parser) parseDropObject(kind ast.DropObjectKind, start ast.Loc, ifExistsOK, cascadeOK bool) (ast.Node, error) {
	stmt := &ast.DropStmt{
		Kind: kind,
		Loc:  ast.Loc{Start: start.Start},
	}

	// Optional IF EXISTS
	if ifExistsOK && p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			stmt.IfExists = true
		}
	}

	// Object name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional CASCADE | RESTRICT
	if cascadeOK {
		switch p.cur.Type {
		case kwCASCADE:
			p.advance()
			stmt.Cascade = true
		case kwRESTRICT:
			p.advance()
			stmt.Restrict = true
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// UNDROP statement dispatch
// ---------------------------------------------------------------------------

// parseUndropStmt parses UNDROP <object_type> name.
// The UNDROP keyword has NOT yet been consumed when this function is called.
func (p *Parser) parseUndropStmt() (ast.Node, error) {
	undropTok := p.advance() // consume UNDROP
	start := undropTok.Loc

	switch p.cur.Type {
	case kwDATABASE:
		// DATABASE handled by T2.1; parseUndropDatabaseStmt consumes DATABASE.
		return p.parseUndropDatabaseStmt()

	case kwSCHEMA:
		// SCHEMA handled by T2.1; parseUndropSchemaStmt consumes SCHEMA.
		return p.parseUndropSchemaStmt()

	case kwTABLE:
		p.advance() // consume TABLE
		return p.parseUndropObject(ast.UndropTable, start)

	case kwDYNAMIC:
		// UNDROP DYNAMIC TABLE
		p.advance() // consume DYNAMIC
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseUndropObject(ast.UndropDynamicTable, start)

	case kwTAG:
		p.advance() // consume TAG
		return p.parseUndropObject(ast.UndropTag, start)

	default:
		objText := p.cur.Str
		if objText == "" {
			objText = TokenName(p.cur.Type)
		}
		err := &ParseError{
			Loc: p.cur.Loc,
			Msg: "UNDROP " + objText + " statement parsing is not yet supported",
		}
		p.skipToNextStatement()
		return nil, err
	}
}

// parseUndropObject parses the name following the UNDROP <type> keywords.
func (p *Parser) parseUndropObject(kind ast.UndropObjectKind, start ast.Loc) (ast.Node, error) {
	stmt := &ast.UndropStmt{
		Kind: kind,
		Loc:  ast.Loc{Start: start.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
