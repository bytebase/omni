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

	// An object read by its documented shape (object.go) whose page documents
	// DROP <kind> [ IF EXISTS ] <name>.
	if f := p.matchObjectForm(formDrops); f != nil && f.drop {
		for range f.words {
			p.advance()
		}
		node, err := p.parseDropObject(ast.DropObject, start, !f.dropNoIfExists, false)
		if err == nil && f.nameParts > 0 {
			_, err = p.checkFormNameParts(f, node.(*ast.DropStmt).Name)
		}
		if err != nil {
			return nil, err
		}
		ds := node.(*ast.DropStmt)
		ds.Object = f.kind
		if f.dropCascade && p.curIs(kwCASCADE) {
			p.advance() // consume CASCADE
			ds.Cascade = true
			ds.Loc.End = p.prev.Loc.End
		}
		if f.dropForce && p.curIs(kwFORCE) {
			p.advance() // consume FORCE
			ds.Force = true
			ds.Loc.End = p.prev.Loc.End
		}
		return ds, nil
	}
	// DROP INDEX [ IF EXISTS ] <table_name>.<index_name> (index.go).
	if p.curIs(kwINDEX) {
		p.advance() // consume INDEX
		return p.parseDropIndexStmt(start)
	}
	// DROP [ <integration kind> ] INTEGRATION (integration.go).
	if p.acceptDropIntegrationKind() {
		return p.parseDropIntegrationStmt(start)
	}
	// Policies read by their docs pages (policy.go).
	if d := p.startsDocPolicy(); d != nil {
		return p.parseDropDocPolicy(start, d)
	}
	// Objects created and altered on other paths (drop_docs.go).
	if d := p.startsDocDrop(); d != nil {
		return p.parseDocDrop(start, d)
	}

	switch {
	case p.curIs(kwDATABASE):
		// DROP DATABASE ROLE [ IF EXISTS ] <name>
		// (https://docs.snowflake.com/en/sql-reference/sql/drop-database-role),
		// whose name a database may qualify.
		if p.startsDropDatabaseRole() {
			p.advance() // consume DATABASE
			p.advance() // consume ROLE
			node, err := p.parseDropObject(ast.DropObject, start, true, false)
			if err != nil {
				return nil, err
			}
			ds := node.(*ast.DropStmt)
			ds.Object = "DATABASE ROLE"
			if !ds.Name.Database.IsEmpty() {
				return nil, &ParseError{
					Position: ds.Name.Loc.Start, End: ds.Name.Loc.End,
					Message: "syntax error: a database role's name has at most 2 part(s)",
				}
			}
			return ds, nil
		}
		// DATABASE handled by T2.1; parseDropDatabaseStmt consumes DATABASE.
		return p.parseDropDatabaseStmt()

	case p.curIs(kwSCHEMA):
		// SCHEMA handled by T2.1; parseDropSchemaStmt consumes SCHEMA.
		return p.parseDropSchemaStmt()

	case p.curIs(kwTABLE):
		p.advance() // consume TABLE
		return p.parseDropTableStmt(start, false)

	case p.curIs(kwVIEW):
		p.advance() // consume VIEW
		return p.parseDropObject(ast.DropView, start, true, false)

	case p.curIs(kwMATERIALIZED):
		// MATERIALIZED VIEW
		p.advance() // consume MATERIALIZED
		if _, err := p.expect(kwVIEW); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropMaterializedView, start, true, false)

	case p.curIs(kwDYNAMIC):
		// DROP DYNAMIC TABLE [ IF EXISTS ] <name>
		// (https://docs.snowflake.com/en/sql-reference/sql/drop-dynamic-table)
		p.advance() // consume DYNAMIC
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropDynamicTable, start, true, false)

	case p.curIs(kwEXTERNAL):
		// EXTERNAL TABLE
		p.advance() // consume EXTERNAL
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropExternalTable, start, true, true)

	case p.curIs(kwSTREAM):
		p.advance() // consume STREAM
		return p.parseDropObject(ast.DropStream, start, true, false)

	case p.curIs(kwTASK):
		// DROP TASK, or DROP TASK GRAPH (task_graph.go).
		task := func() (ast.Node, error) {
			p.advance() // consume TASK
			return p.parseDropObject(ast.DropTask, start, true, false)
		}
		if p.startsTaskGraph() {
			return p.parseTaskOrGraph(func() (ast.Node, error) {
				p.advance() // consume TASK
				p.advance() // consume GRAPH
				node, err := p.parseDropObject(ast.DropObject, start, true, false)
				if err != nil {
					return nil, err
				}
				node.(*ast.DropStmt).Object = taskGraphForm.kind
				return node, nil
			}, task)
		}
		return task()

	case p.curIs(kwSEQUENCE):
		p.advance() // consume SEQUENCE
		return p.parseDropObject(ast.DropSequence, start, true, true)

	case p.curIs(kwSTAGE):
		p.advance() // consume STAGE
		return p.parseDropObject(ast.DropStage, start, true, false)

	case p.curIs(kwFILE_FORMAT):
		p.advance() // consume FILE_FORMAT (single keyword token)
		return p.parseDropObject(ast.DropFileFormat, start, true, false)

	case p.curIs(kwFILE):
		// FILE FORMAT — two-word form using separate FILE and FORMAT tokens
		p.advance() // consume FILE
		if _, err := p.expect(kwFORMAT); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropFileFormat, start, true, false)

	case p.curIs(kwFUNCTION):
		p.advance() // consume FUNCTION
		return p.parseDropObject(ast.DropFunction, start, true, false)

	case p.curIs(kwPROCEDURE):
		p.advance() // consume PROCEDURE
		return p.parseDropObject(ast.DropProcedure, start, true, false)

	case p.curIs(kwPIPE):
		p.advance() // consume PIPE
		return p.parseDropObject(ast.DropPipe, start, true, false)

	case p.curIs(kwTAG):
		p.advance() // consume TAG
		return p.parseDropObject(ast.DropTag, start, true, false)

	case p.curIs(kwROLE):
		p.advance() // consume ROLE
		return p.parseDropObject(ast.DropRole, start, true, false)

	case p.curIs(kwWAREHOUSE):
		p.advance() // consume WAREHOUSE
		return p.parseDropObject(ast.DropWarehouse, start, true, false)

	case p.curIs(kwNETWORK):
		// DROP NETWORK RULE (gap-network-rule) and DROP NETWORK POLICY. NETWORK is
		// reserved but RULE is not, so NETWORK RULE lexes as kwNETWORK followed by a
		// "RULE" identifier; POLICY is a reserved keyword, so NETWORK POLICY lexes as
		// kwNETWORK followed by kwPOLICY.
		if p.peekIsWord("RULE") {
			p.advance() // consume NETWORK
			p.advance() // consume RULE
			return p.parseDropObject(ast.DropNetworkRule, start, true, false)
		}
		if p.nextIs(kwPOLICY) {
			p.advance() // consume NETWORK
			p.advance() // consume POLICY
			return p.parseDropObject(ast.DropNetworkPolicy, start, true, false)
		}
		return p.unsupportedDrop()

	case p.curIs(kwALERT):
		p.advance() // consume ALERT
		return p.parseDropObject(ast.DropAlert, start, true, false)

	case p.curIs(kwCONNECTION):
		p.advance() // consume CONNECTION
		return p.parseDropObject(ast.DropConnection, start, true, false)

	case p.curIs(kwFAILOVER):
		// FAILOVER GROUP
		p.advance() // consume FAILOVER
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropFailoverGroup, start, true, false)

	case p.curIs(kwMANAGED):
		// MANAGED ACCOUNT
		p.advance() // consume MANAGED
		if _, err := p.expect(kwACCOUNT); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropManagedAccount, start, true, false)

	case p.curIs(kwMASKING):
		// MASKING POLICY
		p.advance() // consume MASKING
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropMaskingPolicy, start, true, false)

	case p.curIs(kwREPLICATION):
		// REPLICATION GROUP
		p.advance() // consume REPLICATION
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropReplicationGroup, start, true, false)

	case p.curIs(kwRESOURCE):
		// RESOURCE MONITOR
		p.advance() // consume RESOURCE
		if _, err := p.expect(kwMONITOR); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropResourceMonitor, start, true, false)

	case p.curIs(kwROW):
		// ROW ACCESS POLICY
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropRowAccessPolicy, start, true, false)

	case p.curIs(kwSESSION):
		// SESSION POLICY
		p.advance() // consume SESSION
		if _, err := p.expect(kwPOLICY); err != nil {
			return nil, err
		}
		return p.parseDropObject(ast.DropSessionPolicy, start, true, false)

	case p.curIs(kwSHARE):
		p.advance() // consume SHARE
		return p.parseDropObject(ast.DropShare, start, true, false)

	case p.curIs(kwUSER):
		p.advance() // consume USER
		return p.parseDropObject(ast.DropUser, start, true, false)

	default:
		// DROP ICEBERG TABLE. ICEBERG is not a reserved keyword, so it is
		// matched by text and only before TABLE.
		if p.curIsWord("ICEBERG") && p.nextIs(kwTABLE) {
			p.advance() // consume ICEBERG
			p.advance() // consume TABLE
			return p.parseDropTableStmt(start, true)
		}
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
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: "DROP " + objText + " statement parsing is not yet supported",
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
	if ifExistsOK && p.curIs(kwIF) {
		if p.nextIs(kwEXISTS) {
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

	// A data metric function's signature, ( TABLE ( ... ) [ , TABLE ( ... ) ] )
	// (data_metric_function.go).
	if kind == ast.DropFunction && p.startsDataMetricTableArgs() {
		args, err := p.parseDataMetricTableArgs()
		if err != nil {
			return nil, err
		}
		stmt.TableArgs = args
	}

	// Optional argument-type signature for FUNCTION / PROCEDURE overload
	// disambiguation: DROP FUNCTION f(NUMBER, VARCHAR) / DROP PROCEDURE p().
	if (kind == ast.DropFunction || kind == ast.DropProcedure) && stmt.TableArgs == nil && p.curIs('(') {
		p.advance() // consume '('
		stmt.HasArgs = true
		for !p.curIs(')') && !p.curIs(tokEOF) {
			argType, err := p.parseDataType()
			if err != nil {
				return nil, err
			}
			stmt.ArgTypes = append(stmt.ArgTypes, argType)
			if p.curIs(',') {
				p.advance() // consume ','
				continue
			}
			break
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// Optional CASCADE | RESTRICT
	if cascadeOK {
		switch {
		case p.curIs(kwCASCADE):
			p.advance()
			stmt.Cascade = true
		case p.curIs(kwRESTRICT):
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

	// An object read by its documented shape (object.go) whose page documents
	// UNDROP <kind> <name>.
	if f := p.matchObjectForm(formUndrops); f != nil && f.undrop {
		for range f.words {
			p.advance()
		}
		node, err := p.parseUndropObject(ast.UndropObject, start)
		if err != nil {
			return nil, err
		}
		us := node.(*ast.UndropStmt)
		if _, err := p.checkFormNameParts(f, us.Name); err != nil {
			return nil, err
		}
		us.Object = f.kind
		if f.undropRename {
			if err := p.parseUndropRename(us); err != nil {
				return nil, err
			}
		}
		return us, nil
	}

	switch {
	case p.curIs(kwDATABASE):
		// DATABASE handled by T2.1; parseUndropDatabaseStmt consumes DATABASE.
		return p.parseUndropDatabaseStmt()

	case p.curIs(kwSCHEMA):
		// SCHEMA handled by T2.1; parseUndropSchemaStmt consumes SCHEMA.
		return p.parseUndropSchemaStmt()

	case p.curIs(kwTABLE):
		p.advance() // consume TABLE
		return p.parseUndropObject(ast.UndropTable, start)

	case p.curIs(kwDYNAMIC):
		// UNDROP DYNAMIC TABLE
		p.advance() // consume DYNAMIC
		if _, err := p.expect(kwTABLE); err != nil {
			return nil, err
		}
		return p.parseUndropObject(ast.UndropDynamicTable, start)

	case p.curIs(kwTAG):
		p.advance() // consume TAG
		return p.parseUndropObject(ast.UndropTag, start)

	default:
		if o := p.startsUndropObject(); o != nil {
			return p.parseUndropOther(o, start)
		}
		// UNDROP ICEBERG TABLE
		// (https://docs.snowflake.com/en/sql-reference/sql/undrop-iceberg-table).
		if p.curIsWord("ICEBERG") && p.nextIs(kwTABLE) {
			p.advance() // consume ICEBERG
			p.advance() // consume TABLE
			node, err := p.parseUndropObject(ast.UndropTable, start)
			if err != nil {
				return nil, err
			}
			node.(*ast.UndropStmt).Iceberg = true
			return node, nil
		}
		objText := p.cur.Str
		if objText == "" {
			objText = TokenName(p.cur.Type)
		}
		err := &ParseError{
			Position: p.cur.Loc.Start, End: p.cur.Loc.End,
			Message: "UNDROP " + objText + " statement parsing is not yet supported",
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

// parseDropTableStmt parses the tail of DROP [ ICEBERG ] TABLE, after TABLE:
//
//	[ IF EXISTS ] <name> [ CASCADE | RESTRICT ] [ PURGE | NOPURGE ]
//
// PURGE / NOPURGE come from DROP ICEBERG TABLE, whose ICEBERG keyword is
// optional (https://docs.snowflake.com/en/sql-reference/sql/drop-iceberg-table).
func (p *Parser) parseDropTableStmt(start ast.Loc, iceberg bool) (ast.Node, error) {
	node, err := p.parseDropObject(ast.DropTable, start, true, true)
	if err != nil {
		return nil, err
	}
	stmt := node.(*ast.DropStmt)
	stmt.Iceberg = iceberg
	switch {
	case p.curIsWord("PURGE"):
		p.advance() // consume PURGE
		stmt.Purge = true
	case p.curIsWord("NOPURGE"):
		p.advance() // consume NOPURGE
		stmt.NoPurge = true
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsDropDatabaseRole reports whether cur, after DROP, begins DATABASE ROLE
// <name>: DATABASE, then ROLE, then a token other than the end of the
// statement, CASCADE or RESTRICT, which would make ROLE the name of the
// database DROP DATABASE drops. It consumes nothing.
func (p *Parser) startsDropDatabaseRole() bool {
	if !p.curIs(kwDATABASE) || !p.nextIs(kwROLE) {
		return false
	}
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance() // DATABASE
	p.advance() // ROLE
	switch {
	case p.curIs(tokEOF), p.curIs(';'), p.curIs(tokFlow), p.curIs(kwCASCADE), p.curIs(kwRESTRICT):
		return false
	}
	return true
}
