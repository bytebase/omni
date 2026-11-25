package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Account-level integration-object DDL — CREATE / ALTER (T4.7)
// ---------------------------------------------------------------------------
//
// CREATE / ALTER for the account/schema-level objects:
//
//	{ STORAGE | API | NOTIFICATION | SECURITY | CATALOG | EXTERNAL ACCESS } INTEGRATION
//	RESOURCE MONITOR
//	SECRET
//	CONNECTION
//	EXTERNAL VOLUME
//	GIT REPOSITORY
//
// Every one carries a large, version-growing vocabulary of `KEY = <value>`
// configuration parameters. Rather than mirror the legacy ANTLR grammar's
// finite, already-stale per-type enumerations (its create_*_integration rules
// pin a fixed option order and lack newer params; it has no CREATE EXTERNAL
// VOLUME rule at all), every parameter that follows the object name is parsed as
// an open-ended `KEY = <value>` pair (ast.CopyOption), reusing the merged COPY
// (T5.2) machinery (parseCopyOption / startsCopyOption). Only the structural
// anchors are modeled: the object Kind, RESOURCE MONITOR's WITH keyword + its
// TRIGGERS clause, the GIT REPOSITORY [WITH] TAG clause, CONNECTION's AS REPLICA
// OF clause, and (on ALTER) the action keyword and EXTERNAL VOLUME's
// ADD/REMOVE/UPDATE STORAGE_LOCATION actions. The catalog/semantic layer, not the
// parser, validates that an option is real and legal for the chosen Kind. This
// mirrors the merged STAGE (T4.1) / FILE FORMAT (T4.2) / pipeline (T4.3)
// open-ended approach.
//
// Two structural notes that distinguish these from STAGE/FILE FORMAT:
//   - EXTERNAL VOLUME's STORAGE_LOCATIONS = ( ( NAME=... ... ), ( ... ) ) is a
//     comma-separated *list of parenthesized groups*, a shape the shared
//     parseCopyOptionParen does not handle (it expects either a key/value group
//     or a flat literal list). parseStorageLocationsOption captures it as a
//     CopyOption whose Group holds one unnamed sub-CopyOption per inner group.
//   - VOLUME is not a reserved keyword, so EXTERNAL VOLUME lexes as the kwEXTERNAL
//     keyword followed by a plain "VOLUME" identifier; the CREATE/ALTER dispatch
//     in create_table.go / database_schema.go branches EXTERNAL TABLE vs EXTERNAL
//     VOLUME on that identifier.

// ---------------------------------------------------------------------------
// CREATE
// ---------------------------------------------------------------------------

// parseCreateIntegrationStmt parses the body of a CREATE statement for one of the
// T4.7 objects. The CREATE keyword and the optional OR REPLACE modifier have
// already been consumed by parseCreateStmt; start is the Loc of the CREATE token,
// and cur is the object-type keyword (STORAGE / API / NOTIFICATION / SECURITY /
// CATALOG / EXTERNAL ACCESS / RESOURCE / SECRET / CONNECTION / GIT). EXTERNAL
// VOLUME enters through parseCreateExternalVolumeStmt instead (its dispatch
// already consumed EXTERNAL + VOLUME).
func (p *Parser) parseCreateIntegrationStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	kind, err := p.consumeIntegrationObjectKeyword()
	if err != nil {
		return nil, err
	}

	// The CATALOG and EXTERNAL ACCESS INTEGRATION pages document CREATE
	// [ OR REPLACE ] only.
	if orAlter && (kind == ast.CatalogIntegration || kind == ast.ExternalAccessIntegration) {
		return p.unsupported("CREATE")
	}

	stmt := &ast.CreateIntegrationStmt{
		Kind:      kind,
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseOptionalIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseIntegrationObjectName(kind)
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch kind {
	case ast.ResourceMonitor:
		// CREATE RESOURCE MONITOR <name> WITH <options...> [ TRIGGERS <trigger>... ]
		// The WITH keyword is mandatory per the docs; the option list after it is
		// all-optional. TRIGGERS is a distinct trailing clause.
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		stmt.With = true
		opts, err := p.parseIntegrationOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
		if p.curIs(kwTRIGGERS) {
			triggers, err := p.parseResourceMonitorTriggers()
			if err != nil {
				return nil, err
			}
			stmt.Triggers = triggers
		}

	case ast.Connection:
		// CREATE CONNECTION <name> [ AS REPLICA OF a.b.c ] [ COMMENT = '...' ]
		if p.curIs(kwAS) {
			replica, err := p.parseConnectionReplica()
			if err != nil {
				return nil, err
			}
			stmt.Replica = replica
		}
		opts, err := p.parseIntegrationOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts

	case ast.CatalogIntegration, ast.ExternalAccessIntegration, ast.NotificationIntegration:
		// `KEY = value` properties by their pages' names and shapes, the
		// required ones checked.
		props := notificationIntegrationProps
		switch kind {
		case ast.CatalogIntegration:
			props = catalogIntegrationProps
		case ast.ExternalAccessIntegration:
			props = accessIntegrationProps
		}
		opts, err := p.parseIntegrationProperties(props)
		if err != nil {
			return nil, err
		}
		if err := p.requireIntegrationOptions(kind, opts); err != nil {
			return nil, err
		}
		if err := checkIntegrationForm(kind, opts, stmt.Loc.Start, p.prev.Loc.End); err != nil {
			return nil, err
		}
		stmt.Options = opts

	default:
		// STORAGE/API/SECURITY INTEGRATION, SECRET, GIT REPOSITORY:
		// an open-ended option bag, optionally followed (GIT REPOSITORY) by a
		// [WITH] TAG clause. The loop accepts options and tag clauses in any order
		// so a COMMENT that trails the TAG clause is not dropped (mirrors STAGE).
		for {
			if p.startsIntegrationTags() {
				tags, err := p.parseStageWithTags()
				if err != nil {
					return nil, err
				}
				stmt.Tags = append(stmt.Tags, tags...)
				continue
			}
			if p.startsIntegrationOption() {
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, opt)
				continue
			}
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseCreateExternalVolumeStmt parses
//
//	CREATE [ OR REPLACE ] EXTERNAL VOLUME [ IF NOT EXISTS ] <name>
//	  STORAGE_LOCATIONS = ( ( NAME=... <provider params> ) [, ( ... ) ] )
//	  [ ALLOW_WRITES = { TRUE | FALSE } ]
//	  [ COMMENT = '<string_literal>' ]
//
// The CREATE / OR REPLACE / EXTERNAL / VOLUME tokens have already been consumed by
// the dispatch in create_table.go; start is the Loc of the CREATE token. Every
// parameter is open-ended; only STORAGE_LOCATIONS needs the list-of-groups reader.
func (p *Parser) parseCreateExternalVolumeStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	stmt := &ast.CreateIntegrationStmt{
		Kind:      ast.ExternalVolume,
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseOptionalIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	opts, err := p.parseIntegrationOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER
// ---------------------------------------------------------------------------

// parseAlterIntegrationStmt parses the body of an ALTER statement for one of the
// T4.7 integration-style objects (STORAGE/API/NOTIFICATION/SECURITY INTEGRATION,
// RESOURCE MONITOR, SECRET, CONNECTION, GIT REPOSITORY). The ALTER keyword has
// already been consumed; cur is the object-type keyword. EXTERNAL VOLUME enters
// through parseAlterExternalVolumeStmt instead.
//
// The shared action grammar is:
//
//	[ IF EXISTS ] <name> SET <options>
//	[ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	<name> SET TAG <tag> = '<value>' [ , ... ]
//	<name> UNSET TAG <tag> [ , ... ]
//
// with per-kind extras: RESOURCE MONITOR's SET carries NOTIFY_USERS + TRIGGERS;
// GIT REPOSITORY adds FETCH; CONNECTION adds ENABLE/DISABLE FAILOVER and PRIMARY.
func (p *Parser) parseAlterIntegrationStmt() (ast.Node, error) {
	startLoc := p.cur.Loc // object keyword anchors Loc.Start (ALTER convention)
	kind, err := p.consumeIntegrationObjectKeyword()
	if err != nil {
		return nil, err
	}

	stmt := &ast.AlterIntegrationStmt{Kind: kind, Loc: ast.Loc{Start: startLoc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseIntegrationObjectName(kind)
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.curIs(kwSET):
		if err := p.parseAlterIntegrationSet(stmt); err != nil {
			return nil, err
		}

	case p.curIs(kwUNSET):
		if err := p.parseAlterIntegrationUnset(stmt); err != nil {
			return nil, err
		}

	case p.curIs(kwTRIGGERS):
		// ALTER RESOURCE MONITOR [ IF EXISTS ] <name> [ SET { ... } ]
		//   [ TRIGGERS triggerDefinition [ triggerDefinition ... ] ]
		// (https://docs.snowflake.com/en/sql-reference/sql/alter-resource-monitor):
		// TRIGGERS without SET is held as a SET that sets nothing else.
		if kind != ast.ResourceMonitor {
			return nil, p.syntaxErrorAtCur()
		}
		triggers, err := p.parseResourceMonitorTriggers()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterIntegrationSet
		stmt.Triggers = triggers

	case p.curIs(kwFETCH):
		// GIT REPOSITORY FETCH.
		p.advance() // consume FETCH
		stmt.Action = ast.AlterIntegrationFetch

	case p.curIs(kwENABLE), p.curIs(kwDISABLE):
		// CONNECTION ENABLE/DISABLE FAILOVER [ TO ACCOUNTS a.b [, ...] ].
		if err := p.parseAlterConnectionFailover(stmt); err != nil {
			return nil, err
		}

	case p.curIs(kwPRIMARY):
		// CONNECTION PRIMARY.
		p.advance() // consume PRIMARY
		stmt.Action = ast.AlterIntegrationPrimary

	case p.curIs(kwREFRESH):
		// [ SECURITY ] INTEGRATION REFRESH.
		p.parseAlterIntegrationRefresh(stmt)

	default:
		return nil, p.syntaxErrorAtCur()
	}

	if !integrationAlterActionDocumented(kind, stmt.Action) {
		return nil, &ParseError{
			Position: startLoc.Start, End: p.prev.Loc.End,
			Message: "syntax error: ALTER " + kind.String() + " does not take this action",
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// integrationAlterActionDocumented reports whether ALTER of kind documents
// action. Only the kinds whose pages name a short list are checked here:
//
//	ALTER CATALOG INTEGRATION [ IF EXISTS ] <name> SET ...
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-catalog-integration
//
//	ALTER EXTERNAL ACCESS INTEGRATION [ IF EXISTS ] <name> SET ...
//	ALTER EXTERNAL ACCESS INTEGRATION [ IF EXISTS ] <name> UNSET ...
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration
// (whose SET and UNSET carry TAG <tag> = '<value>' and TAG <tag>), and REFRESH,
// which only a [ SECURITY ] INTEGRATION takes.
func integrationAlterActionDocumented(kind ast.IntegrationObjectKind, action ast.AlterIntegrationAction) bool {
	switch kind {
	case ast.CatalogIntegration:
		return action == ast.AlterIntegrationSet
	case ast.ExternalAccessIntegration:
		switch action {
		case ast.AlterIntegrationSet, ast.AlterIntegrationUnset, ast.AlterIntegrationSetTag, ast.AlterIntegrationUnsetTag:
			return true
		}
		return false
	case ast.StorageIntegration, ast.SecurityIntegration:
		// A bare ALTER INTEGRATION reads as StorageIntegration (see
		// consumeIntegrationObjectKeyword), so it may be a security integration.
		return true
	}
	return action != ast.AlterIntegrationRefresh
}

// parseAlterIntegrationRefresh reads a security integration's REFRESH action.
// cur is REFRESH.
//
//	ALTER [ SECURITY ] INTEGRATION [ IF EXISTS ] <name> REFRESH { OAUTH_CLIENT_SECRET | OAUTH_CLIENT_SECRET_2 }
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-oauth-snowflake
//
//	ALTER [ SECURITY ] INTEGRATION <name> REFRESH [ SAML2_SNOWFLAKE_PRIVATE_KEY ] [ METADATA_URL ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-security-integration-saml2
func (p *Parser) parseAlterIntegrationRefresh(stmt *ast.AlterIntegrationStmt) {
	p.advance() // consume REFRESH
	stmt.Action = ast.AlterIntegrationRefresh
	switch {
	case p.curIsWord("OAUTH_CLIENT_SECRET"), p.curIsWord("OAUTH_CLIENT_SECRET_2"):
		stmt.Refresh = append(stmt.Refresh, strings.ToUpper(p.advance().Str))
		return
	case p.curIsWord("SAML2_SNOWFLAKE_PRIVATE_KEY"):
		stmt.Refresh = append(stmt.Refresh, strings.ToUpper(p.advance().Str))
	}
	if p.curIsWord("METADATA_URL") {
		stmt.Refresh = append(stmt.Refresh, strings.ToUpper(p.advance().Str))
	}
}

// parseAlterExternalVolumeStmt parses
//
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> ADD STORAGE_LOCATION = (...)
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> REMOVE STORAGE_LOCATION '<name>'
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> UPDATE STORAGE_LOCATION = '<name>' CREDENTIALS = (...)
//	ALTER EXTERNAL VOLUME [ IF EXISTS ] <name> SET { ALLOW_WRITES | COMMENT } = ...
//
// EXTERNAL + VOLUME have already been consumed by the dispatch in
// database_schema.go; startLoc is the EXTERNAL keyword Loc.
func (p *Parser) parseAlterExternalVolumeStmt(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.AlterIntegrationStmt{Kind: ast.ExternalVolume, Loc: ast.Loc{Start: startLoc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.curIs(kwADD):
		// ADD STORAGE_LOCATION = ( NAME=... <provider params> )
		p.advance() // consume ADD
		opt, err := p.parseStorageLocationSingleOption()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterIntegrationAddLocation
		stmt.Options = []*ast.CopyOption{opt}

	case p.curIs(kwREMOVE):
		// REMOVE STORAGE_LOCATION '<name>'
		p.advance() // consume REMOVE
		if err := p.expectStorageLocationKeyword(); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterIntegrationRemoveLocation
		stmt.Location = tok.Str

	case p.curIs(kwUPDATE):
		// UPDATE STORAGE_LOCATION = '<name>' <params...>
		p.advance() // consume UPDATE
		if err := p.expectStorageLocationKeyword(); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		stmt.Location = tok.Str
		opts, err := p.parseIntegrationOptions()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterIntegrationUpdateLocation
		stmt.Options = opts

	case p.curIs(kwSET):
		// SET <options> — open-ended (ALLOW_WRITES / COMMENT).
		p.advance() // consume SET
		opts, err := p.parseIntegrationOptions()
		if err != nil {
			return nil, err
		}
		if len(opts) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterIntegrationSet
		stmt.Options = opts

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER action helpers
// ---------------------------------------------------------------------------

// parseAlterIntegrationSet parses the SET branch shared by the integration-style
// objects. cur is the SET keyword. The branches are:
//
//	SET TAG <tag> = '<value>' [ , ... ]            -> AlterIntegrationSetTag
//	SET <options> [ NOTIFY_USERS=(...) ] [ TRIGGERS ... ]  -> AlterIntegrationSet
//
// NOTIFY_USERS is captured as an ordinary open-ended option; only TRIGGERS needs
// the dedicated reader (it is the one non-`KEY = value` clause). RESOURCE MONITOR
// is the only kind that carries TRIGGERS, but accepting it for any kind here is
// harmless — a non-RESOURCE-MONITOR statement simply never contains it.
func (p *Parser) parseAlterIntegrationSet(stmt *ast.AlterIntegrationStmt) error {
	p.advance() // consume SET

	if stmt.Kind == ast.CatalogIntegration || (stmt.Kind == ast.ExternalAccessIntegration && !p.curIs(kwTAG)) {
		return p.parseAlterIntegrationSetProperties(stmt)
	}

	if p.curIs(kwTAG) {
		tags, err := p.parseTagSetList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterIntegrationSetTag
		stmt.Tags = tags
		return nil
	}

	opts, err := p.parseIntegrationOptions()
	if err != nil {
		return err
	}
	if p.curIs(kwTRIGGERS) {
		triggers, err := p.parseResourceMonitorTriggers()
		if err != nil {
			return err
		}
		stmt.Triggers = triggers
	}
	if len(opts) == 0 && len(stmt.Triggers) == 0 {
		// SET with nothing settable is a syntax error.
		return p.syntaxErrorAtCur()
	}
	stmt.Action = ast.AlterIntegrationSet
	stmt.Options = opts
	return nil
}

// parseAlterIntegrationUnset parses the UNSET branch shared by the
// integration-style objects. cur is the UNSET keyword.
//
//	UNSET TAG <tag> [ , ... ]            -> AlterIntegrationUnsetTag
//	UNSET <property> [ , ... ]           -> AlterIntegrationUnset (e.g. ENABLED, COMMENT)
func (p *Parser) parseAlterIntegrationUnset(stmt *ast.AlterIntegrationStmt) error {
	p.advance() // consume UNSET

	if stmt.Kind == ast.ExternalAccessIntegration {
		return p.parseAlterAccessIntegrationUnset(stmt)
	}

	if p.curIs(kwTAG) {
		names, err := p.parseUnsetTagNameList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterIntegrationUnsetTag
		stmt.UnsetTags = names
		return nil
	}

	// UNSET <property> [, ...]. Each property is an open-ended run of name words
	// (reusing the STAGE UNSET-property reader, which captures multi-word
	// properties whole and uppercases them). A comma separates properties.
	props, err := p.parseStageUnsetProps()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterIntegrationUnset
	stmt.UnsetProps = props
	return nil
}

// parseAlterConnectionFailover parses CONNECTION's ENABLE / DISABLE FAILOVER
// action. cur is the ENABLE or DISABLE keyword.
//
//	ENABLE FAILOVER TO ACCOUNTS <org>.<account> [ , ... ] [ IGNORE EDITION CHECK ]
//	DISABLE FAILOVER [ TO ACCOUNTS <org>.<account> [ , ... ] ]
func (p *Parser) parseAlterConnectionFailover(stmt *ast.AlterIntegrationStmt) error {
	enable := p.curIs(kwENABLE)
	p.advance() // consume ENABLE / DISABLE
	if _, err := p.expect(kwFAILOVER); err != nil {
		return err
	}

	if enable {
		stmt.Action = ast.AlterIntegrationEnableFailover
		// ENABLE FAILOVER requires TO ACCOUNTS.
		if _, err := p.expect(kwTO); err != nil {
			return err
		}
		if _, err := p.expect(kwACCOUNTS); err != nil {
			return err
		}
		accounts, err := p.parseAccountNameList()
		if err != nil {
			return err
		}
		stmt.Accounts = accounts
		// Optional IGNORE EDITION CHECK.
		if p.curIs(kwIGNORE) {
			p.advance() // consume IGNORE
			if _, err := p.expect(kwEDITION); err != nil {
				return err
			}
			if _, err := p.expect(kwCHECK); err != nil {
				return err
			}
		}
		return nil
	}

	stmt.Action = ast.AlterIntegrationDisableFailover
	// DISABLE FAILOVER [ TO ACCOUNTS ... ].
	if p.curIs(kwTO) {
		p.advance() // consume TO
		if _, err := p.expect(kwACCOUNTS); err != nil {
			return err
		}
		accounts, err := p.parseAccountNameList()
		if err != nil {
			return err
		}
		stmt.Accounts = accounts
	}
	return nil
}

// parseAccountNameList parses a comma-separated list of <org>.<account> names
// (each a dotted ObjectName). Consumes at least one.
func (p *Parser) parseAccountNameList() ([]*ast.ObjectName, error) {
	var names []*ast.ObjectName
	for {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
		if p.curIs(',') {
			p.advance() // consume ','
			continue
		}
		break
	}
	return names, nil
}

// ---------------------------------------------------------------------------
// CONNECTION AS REPLICA OF
// ---------------------------------------------------------------------------

// parseConnectionReplica parses an AS REPLICA OF <org>.<account>.<connection>
// clause. cur is the AS keyword.
func (p *Parser) parseConnectionReplica() (*ast.ConnectionReplica, error) {
	asTok := p.advance() // consume AS
	if _, err := p.expect(kwREPLICA); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwOF); err != nil {
		return nil, err
	}
	src, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	return &ast.ConnectionReplica{
		Source: src,
		Loc:    ast.Loc{Start: asTok.Loc.Start, End: p.prev.Loc.End},
	}, nil
}

// ---------------------------------------------------------------------------
// RESOURCE MONITOR triggers
// ---------------------------------------------------------------------------

// parseResourceMonitorTriggers parses a TRIGGERS clause: the TRIGGERS keyword
// followed by one or more space-separated trigger definitions. cur is the
// TRIGGERS keyword. Per the docs the trigger definitions are NOT comma-separated.
//
//	TRIGGERS ( ON <num> PERCENT DO { SUSPEND | SUSPEND_IMMEDIATE | NOTIFY } )+
func (p *Parser) parseResourceMonitorTriggers() ([]*ast.ResourceMonitorTrigger, error) {
	if _, err := p.expect(kwTRIGGERS); err != nil {
		return nil, err
	}

	var triggers []*ast.ResourceMonitorTrigger
	for p.curIs(kwON) {
		trig, err := p.parseResourceMonitorTrigger()
		if err != nil {
			return nil, err
		}
		triggers = append(triggers, trig)
	}
	if len(triggers) == 0 {
		// TRIGGERS with no trigger definition is a syntax error.
		return nil, p.syntaxErrorAtCur()
	}
	return triggers, nil
}

// parseResourceMonitorTrigger parses one trigger definition:
//
//	ON <num> PERCENT DO { SUSPEND | SUSPEND_IMMEDIATE | NOTIFY }
//
// cur is the ON keyword.
func (p *Parser) parseResourceMonitorTrigger() (*ast.ResourceMonitorTrigger, error) {
	onTok := p.advance() // consume ON
	trig := &ast.ResourceMonitorTrigger{Loc: ast.Loc{Start: onTok.Loc.Start}}

	thresholdTok, err := p.expect(tokInt)
	if err != nil {
		return nil, err
	}
	trig.Threshold = thresholdTok.Ival

	if _, err := p.expect(kwPERCENT); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}

	switch {
	case p.curIs(kwSUSPEND):
		p.advance()
		trig.Action = "SUSPEND"
	case p.curIs(kwSUSPEND_IMMEDIATE):
		p.advance()
		trig.Action = "SUSPEND_IMMEDIATE"
	case p.curIs(kwNOTIFY):
		p.advance()
		trig.Action = "NOTIFY"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	trig.Loc.End = p.prev.Loc.End
	return trig, nil
}

// ---------------------------------------------------------------------------
// EXTERNAL VOLUME STORAGE_LOCATIONS list-of-groups
// ---------------------------------------------------------------------------

// parseStorageLocationsOption parses the EXTERNAL VOLUME
//
//	STORAGE_LOCATIONS = ( ( NAME=... <params> ) [, ( ... ) ] )
//
// option, whose value is a comma-separated list of parenthesized key/value
// groups. cur is the STORAGE_LOCATIONS option name. The result is a CopyOption
// named STORAGE_LOCATIONS whose Group holds one unnamed sub-CopyOption per inner
// group (each carrying that inner group's own Group of NAME/STORAGE_PROVIDER/...
// entries). This shape is not handled by the shared parseCopyOptionParen (which
// expects a key/value group or a flat literal list, never a list of groups), so
// it is read here.
func (p *Parser) parseStorageLocationsOption() (*ast.CopyOption, error) {
	nameTok := p.advance() // consume STORAGE_LOCATIONS
	opt := &ast.CopyOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start, End: nameTok.Loc.End},
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var groups []*ast.CopyOption
	for !p.curIs(')') && !p.curIs(tokEOF) {
		inner, err := p.parseStorageLocationGroup()
		if err != nil {
			return nil, err
		}
		groups = append(groups, inner)
		if p.curIs(',') {
			p.advance() // consume optional ',' separator between groups
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	opt.Group = groups
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseStorageLocationGroup parses one parenthesized inner group of a
// STORAGE_LOCATIONS list, e.g. ( NAME = 'x' STORAGE_PROVIDER = 'S3' ENCRYPTION =
// (...) ). cur is the inner '('. The group is returned as an unnamed CopyOption
// (Name == "") whose Group holds the inner key/value entries, parsed by the
// shared parseCopyOptionParen.
func (p *Parser) parseStorageLocationGroup() (*ast.CopyOption, error) {
	if !p.curIs('(') {
		return nil, p.syntaxErrorAtCur()
	}
	start := p.cur.Loc.Start
	group := &ast.CopyOption{Loc: ast.Loc{Start: start}}
	if err := p.parseCopyOptionParen(group); err != nil {
		return nil, err
	}
	group.Loc.End = p.prev.Loc.End
	return group, nil
}

// parseStorageLocationSingleOption parses an ALTER EXTERNAL VOLUME ADD
// STORAGE_LOCATION = ( NAME=... <params> ) option (a single parenthesized group,
// not the bracketed list form). cur is the STORAGE_LOCATION option name.
func (p *Parser) parseStorageLocationSingleOption() (*ast.CopyOption, error) {
	if !p.curIsStorageLocationWord() {
		return nil, p.syntaxErrorAtCur()
	}
	// STORAGE_LOCATION = ( ... ) is exactly the COPY option shape, so the shared
	// parseCopyOption reader handles it directly.
	return p.parseCopyOption()
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// consumeIntegrationObjectKeyword consumes the object-type keyword(s) at the
// current position and returns the corresponding IntegrationObjectKind. On entry
// cur is the leading keyword. It handles the optional INTEGRATION word for the
// {STORAGE|API|NOTIFICATION|SECURITY} INTEGRATION forms, the two-word RESOURCE
// MONITOR and GIT REPOSITORY forms, and the single-word SECRET / CONNECTION.
//
// Per the docs and the legacy grammar, ALTER accepts a bare INTEGRATION keyword
// with the leading qualifier word omitted (ALTER [STORAGE] INTEGRATION ...). On
// CREATE the qualifier is required and INTEGRATION follows it; this helper maps
// both spellings since a bare ALTER INTEGRATION cannot recover the original
// subtype anyway (the semantic layer resolves it from the catalog).
func (p *Parser) consumeIntegrationObjectKeyword() (ast.IntegrationObjectKind, error) {
	switch {
	case p.curIs(kwSTORAGE):
		p.advance() // consume STORAGE
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.StorageIntegration, nil
	case p.curIs(kwAPI):
		p.advance() // consume API
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.APIIntegration, nil
	case p.curIs(kwNOTIFICATION):
		p.advance() // consume NOTIFICATION
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.NotificationIntegration, nil
	case p.curIs(kwSECURITY):
		p.advance() // consume SECURITY
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.SecurityIntegration, nil
	case p.curIs(kwINTEGRATION):
		// Bare ALTER INTEGRATION <name> ... — subtype qualifier omitted.
		p.advance() // consume INTEGRATION
		return ast.StorageIntegration, nil
	case p.curIs(kwRESOURCE):
		p.advance() // consume RESOURCE
		if _, err := p.expect(kwMONITOR); err != nil {
			return 0, err
		}
		return ast.ResourceMonitor, nil
	case p.curIs(kwSECRET):
		p.advance() // consume SECRET
		return ast.Secret, nil
	case p.curIs(kwCONNECTION):
		p.advance() // consume CONNECTION
		return ast.Connection, nil
	case p.curIs(kwGIT):
		p.advance() // consume GIT
		if _, err := p.expect(kwREPOSITORY); err != nil {
			return 0, err
		}
		return ast.GitRepository, nil
	}
	switch {
	case p.acceptWords("CATALOG", "INTEGRATION"):
		return ast.CatalogIntegration, nil
	case p.acceptWords("EXTERNAL", "ACCESS", "INTEGRATION"):
		return ast.ExternalAccessIntegration, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// startsCatalogOrAccessIntegration reports whether cur begins CATALOG
// INTEGRATION or EXTERNAL ACCESS INTEGRATION, the integrations the CREATE and
// ALTER dispatch switches do not reach by their first keyword (CATALOG is not
// a keyword; EXTERNAL also starts TABLE, VOLUME and FUNCTION). It consumes
// nothing.
func (p *Parser) startsCatalogOrAccessIntegration() bool {
	return p.atWords("CATALOG", "INTEGRATION") || p.atWords("EXTERNAL", "ACCESS", "INTEGRATION")
}

// acceptDropIntegrationKind consumes the object words of
//
//	DROP [ { API | CATALOG | EXTERNAL ACCESS | NOTIFICATION | SECURITY | STORAGE } ] INTEGRATION
//
// when they start at cur, and reports whether it did. The qualifier names no
// different object, so every spelling reads as ast.DropIntegration.
//
// https://docs.snowflake.com/en/sql-reference/sql/drop-integration
// https://docs.snowflake.com/en/sql-reference/sql/drop-catalog-integration
func (p *Parser) acceptDropIntegrationKind() bool {
	for _, q := range [][]string{{}, {"API"}, {"CATALOG"}, {"EXTERNAL", "ACCESS"}, {"NOTIFICATION"}, {"SECURITY"}, {"STORAGE"}} {
		if p.acceptWords(append(q, "INTEGRATION")...) {
			return true
		}
	}
	return false
}

// parseDropIntegrationStmt reads [ IF EXISTS ] <name> once
// acceptDropIntegrationKind has consumed the object words. An integration is
// an account-level object, so its name takes no qualifier.
func (p *Parser) parseDropIntegrationStmt(start ast.Loc) (ast.Node, error) {
	node, err := p.parseDropObject(ast.DropIntegration, start, true, false)
	if err != nil {
		return nil, err
	}
	if name := node.(*ast.DropStmt).Name; !name.Schema.IsEmpty() || !name.Database.IsEmpty() {
		return nil, &ParseError{
			Position: name.Loc.Start, End: name.Loc.End,
			Message: "syntax error: " + name.String() + " takes no qualifier here",
		}
	}
	return node, nil
}

// parseIntegrationObjectName reads the name of an object of kind: one part
// for an integration, which is an account-level object (the pages write
// <name>, never a qualified one), and a qualified name for the other kinds.
func (p *Parser) parseIntegrationObjectName(kind ast.IntegrationObjectKind) (*ast.ObjectName, error) {
	switch kind {
	case ast.StorageIntegration, ast.APIIntegration, ast.NotificationIntegration, ast.SecurityIntegration,
		ast.CatalogIntegration, ast.ExternalAccessIntegration:
		return p.parseObjectSimpleName()
	}
	return p.parseObjectName()
}

// parseIntegrationOptions parses a run of zero or more open-ended `KEY = value`
// options. STORAGE_LOCATIONS is special-cased to the list-of-groups reader; every
// other option uses the shared parseCopyOption. The run terminates at a statement
// boundary, at a structural keyword that anchors a trailing clause (TRIGGERS /
// WITH / TAG), or at any token that does not begin an option.
func (p *Parser) parseIntegrationOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsIntegrationOption() {
		var (
			opt *ast.CopyOption
			err error
		)
		if p.curIsWord("STORAGE_LOCATIONS") {
			opt, err = p.parseStorageLocationsOption()
		} else {
			opt, err = p.parseCopyOption()
		}
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// requireIntegrationOptions fails unless opts sets each property the CREATE
// page of kind marks required. Only the kinds read with such a list are
// checked:
//
//	CREATE CATALOG INTEGRATION <name> CATALOG_SOURCE = ... TABLE_FORMAT = ... ENABLED = ...
//
// on every catalog-source page under
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration, and
//
//	CREATE EXTERNAL ACCESS INTEGRATION <name> ALLOWED_NETWORK_RULES = ( ... ) ENABLED = ...
//
// https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration
func (p *Parser) requireIntegrationOptions(kind ast.IntegrationObjectKind, opts []*ast.CopyOption) error {
	var required []string
	switch kind {
	case ast.CatalogIntegration:
		required = []string{"CATALOG_SOURCE", "TABLE_FORMAT", "ENABLED"}
	case ast.ExternalAccessIntegration:
		required = []string{"ALLOWED_NETWORK_RULES", "ENABLED"}
	}
	for _, n := range required {
		found := false
		for _, o := range opts {
			if o.Name == n {
				found = true
				break
			}
		}
		if !found {
			return &ParseError{
				Position: p.cur.Loc.Start, End: p.cur.Loc.End,
				Message: "syntax error: " + n + " is required",
			}
		}
	}
	return nil
}

// startsIntegrationOption reports whether the current token can begin an option.
// It is the COPY option predicate minus the structural keywords that anchor a
// trailing clause: TRIGGERS (RESOURCE MONITOR), and WITH / TAG (GIT REPOSITORY's
// [WITH] TAG). Those must not be swallowed as an option name.
func (p *Parser) startsIntegrationOption() bool {
	switch {
	case p.curIs(kwTRIGGERS), p.curIs(kwWITH), p.curIs(kwTAG):
		return false
	}
	return p.startsCopyOption()
}

// startsIntegrationTags reports whether the current position begins a GIT
// REPOSITORY [WITH] TAG clause (TAG directly, or WITH immediately followed by
// TAG). Reuses the STAGE predicate semantics.
func (p *Parser) startsIntegrationTags() bool {
	if p.curIs(kwTAG) {
		return true
	}
	return p.curIs(kwWITH) && p.nextIs(kwTAG)
}

// curIsStorageLocationWord reports whether cur is the STORAGE_LOCATION option name
// (used by ALTER EXTERNAL VOLUME ADD / UPDATE).
func (p *Parser) curIsStorageLocationWord() bool {
	return p.curIsWord("STORAGE_LOCATION")
}

// expectStorageLocationKeyword consumes the STORAGE_LOCATION option name, erroring
// if the current token is not it.
func (p *Parser) expectStorageLocationKeyword() error {
	if !p.curIsStorageLocationWord() {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume STORAGE_LOCATION
	return nil
}

// parseOptionalIfNotExists consumes an optional IF NOT EXISTS clause, setting
// *flag when present. Shared by the T4.7 CREATE parsers.
func (p *Parser) parseOptionalIfNotExists(flag *bool) error {
	if p.curIs(kwIF) {
		if p.nextIs(kwNOT) {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return err
			}
			*flag = true
		}
	}
	return nil
}

// parseOptionalIfExists consumes an optional IF EXISTS clause, setting *flag when
// present. Shared by the T4.7 ALTER parsers.
func (p *Parser) parseOptionalIfExists(flag *bool) {
	if p.curIs(kwIF) {
		if p.nextIs(kwEXISTS) {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			*flag = true
		}
	}
}

// ---------------------------------------------------------------------------
// CATALOG, EXTERNAL ACCESS and NOTIFICATION INTEGRATION properties
// ---------------------------------------------------------------------------
//
// These integrations' CREATE (and the first two's ALTER ... SET) are read by
// their pages' shape rather than through the open-ended COPY option reader the
// other kinds use: every property is `KEY = value`, its key one the pages
// document, and its value of the shape they give it. A group's entries are
// read the same way against the group's own keys.

// integrationProp is a documented property's value: a leaf of one of shapes,
// a group of keys, a list of such groups, or string pairs.
type integrationProp struct {
	shapes []objectValue    // a leaf of one of these shapes (object.go)
	group  integrationProps // ( <key> = <value> [ <key> = <value> ... ] )
	groups bool             // with group: ( ( ... ) [ , ( ... ) ... ] )
	pairs  bool             // ( '<key>' = '<value>' [ , ... ] )
	// oneParts holds a list of names to one-part names, as the pages write
	// a list of account-level objects (grammar ID).
	oneParts bool
}

// integrationProps maps a documented property's name to its value.
type integrationProps map[string]integrationProp

func intLeaf(shapes ...objectValue) integrationProp  { return integrationProp{shapes: shapes} }
func intGroup(keys integrationProps) integrationProp { return integrationProp{group: keys} }

var (
	intString = intLeaf(propString)
	intBool   = intLeaf(propBool)
)

// catalogRestConfig is REST_CONFIG's keys, across the catalog-source pages.
var catalogRestConfig = integrationProps{
	"CATALOG_URI":  intString,
	"PREFIX":       intString,
	"CATALOG_NAME": intString,
	"CATALOG_API_TYPE": intLeaf(propOneOf("PUBLIC", "PRIVATE", "AWS_API_GATEWAY", "AWS_PRIVATE_API_GATEWAY",
		"AWS_GLUE", "AWS_PRIVATE_GLUE", "AWS_S3TABLES", "AWS_PRIVATE_S3TABLES")),
	"ACCESS_DELEGATION_MODE": intLeaf(propOneOf("VENDED_CREDENTIALS", "EXTERNAL_VOLUME_CREDENTIALS")),
	"POSTGRES_INSTANCE":      intString,
}

// catalogRestAuthentication is REST_AUTHENTICATION's keys, across the
// catalog-source pages.
var catalogRestAuthentication = integrationProps{
	"TYPE":                 intLeaf(propOneOf("OAUTH", "BEARER", "SIGV4", "OIDC")),
	"OAUTH_API_TYPE":       intLeaf(propOneOf("PUBLIC", "PRIVATE")),
	"OAUTH_TOKEN_URI":      intString,
	"OAUTH_CLIENT_ID":      intString,
	"OAUTH_CLIENT_SECRET":  intString,
	"OAUTH_ALLOWED_SCOPES": intLeaf(propStrings),
	"OAUTH_GRANT_TYPE":     intLeaf(propOneOf("TOKEN_EXCHANGE")),
	"OAUTH_AUDIENCE":       intString,
	"BEARER_TOKEN":         intString,
	"SIGV4_IAM_ROLE":       intString,
	"SIGV4_SIGNING_REGION": intString,
	"SIGV4_EXTERNAL_ID":    intString,
	"OIDC_AUDIENCE":        intString,
}

// catalogDefaultStorageConfig is DEFAULT_STORAGE_CONFIG's keys.
var catalogDefaultStorageConfig = integrationProps{"USE_PRIVATELINK_ENDPOINT": intBool}

// catalogIntegrationProps is CREATE CATALOG INTEGRATION's properties, across
// the pages of its catalog sources:
//
//	CREATE [ OR REPLACE ] CATALOG INTEGRATION [ IF NOT EXISTS ] <name>
//	  CATALOG_SOURCE = { GLUE | OBJECT_STORE | POLARIS | ICEBERG_REST | SNOWFLAKE_POSTGRES | DELTA_SHARING }
//	  TABLE_FORMAT = { ICEBERG | DELTA | NONE | HIVE }
//	  [ CATALOG_NAMESPACE = '<namespace>' ]
//	  [ REST_CONFIG = ( ... ) ] [ REST_AUTHENTICATION = ( ... ) ] [ DEFAULT_STORAGE_CONFIG = ( ... ) ]
//	  ENABLED = { TRUE | FALSE }
//	  [ REFRESH_INTERVAL_SECONDS = <value> ] [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-open-catalog
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-horizon
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-snowflake-postgres
// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-delta-sharing
var catalogIntegrationProps = integrationProps{
	"CATALOG_SOURCE": intLeaf(propOneOf("GLUE", "OBJECT_STORE", "POLARIS", "ICEBERG_REST",
		"SNOWFLAKE_POSTGRES", "DELTA_SHARING")),
	"TABLE_FORMAT":             intLeaf(propOneOf("ICEBERG", "DELTA", "NONE", "HIVE")),
	"CATALOG_NAMESPACE":        intString,
	"REST_CONFIG":              intGroup(catalogRestConfig),
	"REST_AUTHENTICATION":      intGroup(catalogRestAuthentication),
	"DEFAULT_STORAGE_CONFIG":   intGroup(catalogDefaultStorageConfig),
	"ENABLED":                  intBool,
	"REFRESH_INTERVAL_SECONDS": intLeaf(propNumber),
	"COMMENT":                  intString,
}

// catalogIntegrationSetProps is ALTER CATALOG INTEGRATION SET's properties:
//
//	SET [ REST_AUTHENTICATION = ( { OAUTH_CLIENT_SECRET | BEARER_TOKEN } = '<value>' ) ]
//	    [ DEFAULT_STORAGE_CONFIG = ( USE_PRIVATELINK_ENDPOINT = { TRUE | FALSE } ) ]
//	    [ REFRESH_INTERVAL_SECONDS = <value> ] [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-catalog-integration
var catalogIntegrationSetProps = integrationProps{
	"REST_AUTHENTICATION":      intGroup(integrationProps{"OAUTH_CLIENT_SECRET": intString, "BEARER_TOKEN": intString}),
	"DEFAULT_STORAGE_CONFIG":   intGroup(catalogDefaultStorageConfig),
	"REFRESH_INTERVAL_SECONDS": intLeaf(propNumber),
	"COMMENT":                  intString,
}

// accessIntegrationProps is EXTERNAL ACCESS INTEGRATION's properties, which
// its CREATE and ALTER ... SET share:
//
//	ALLOWED_NETWORK_RULES = ( <rule_name> [ , ... ] )
//	[ ALLOWED_API_AUTHENTICATION_INTEGRATIONS = ( <integration_name> [ , ... ] ) | none ]
//	[ ALLOWED_AUTHENTICATION_SECRETS = ( <secret_name> [ , ... ] ) | all | none ]
//	ENABLED = { TRUE | FALSE }
//	[ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-external-access-integration
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration
var accessIntegrationProps = integrationProps{
	"ALLOWED_NETWORK_RULES":                   intLeaf(propNames),
	"ALLOWED_API_AUTHENTICATION_INTEGRATIONS": {shapes: []objectValue{propNames, propOneOf("NONE")}, oneParts: true},
	"ALLOWED_AUTHENTICATION_SECRETS":          intLeaf(propNames, propOneOf("ALL", "NONE")),
	"ENABLED":                                 intBool,
	"COMMENT":                                 intString,
}

// accessIntegrationUnset is the properties ALTER EXTERNAL ACCESS INTEGRATION
// UNSET names.
var accessIntegrationUnset = []string{
	"ALLOWED_NETWORK_RULES", "ALLOWED_API_AUTHENTICATION_INTEGRATIONS", "ALLOWED_AUTHENTICATION_SECRETS", "COMMENT",
}

// notificationProvider is NOTIFICATION_PROVIDER's values.
var notificationProvider = intLeaf(propOneOf("AWS_SNS", "AZURE_EVENT_GRID", "AZURE_STORAGE_QUEUE", "GCP_PUBSUB"))

// notificationIntegrationProps is CREATE NOTIFICATION INTEGRATION's properties,
// across the pages of its types: queues (inbound Azure and GCP, outbound AWS,
// Azure and GCP, multi-queue), webhooks and email.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-azure
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-gcp
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-aws
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-azure
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-gcp
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-multi-queue
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-webhooks
// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-email
var notificationIntegrationProps = integrationProps{
	"ENABLED":                         intBool,
	"TYPE":                            intLeaf(propOneOf("QUEUE", "MULTI_QUEUE", "WEBHOOK", "EMAIL")),
	"DIRECTION":                       intLeaf(propOneOf("INBOUND", "OUTBOUND")),
	"NOTIFICATION_PROVIDER":           notificationProvider,
	"AWS_SNS_TOPIC_ARN":               intString,
	"AWS_SNS_ROLE_ARN":                intString,
	"AZURE_EVENT_GRID_TOPIC_ENDPOINT": intString,
	"AZURE_STORAGE_QUEUE_PRIMARY_URI": intString,
	"AZURE_TENANT_ID":                 intString,
	"USE_PRIVATELINK_ENDPOINT":        intBool,
	"GCP_PUBSUB_SUBSCRIPTION_NAME":    intString,
	"GCP_PUBSUB_TOPIC_NAME":           intString,
	"QUEUES": {groups: true, group: integrationProps{
		"NAME":                            intString,
		"NOTIFICATION_PROVIDER":           notificationProvider,
		"AWS_SNS_TOPIC_ARN":               intString,
		"GCP_PUBSUB_SUBSCRIPTION_NAME":    intString,
		"AZURE_STORAGE_QUEUE_PRIMARY_URI": intString,
		"AZURE_TENANT_ID":                 intString,
	}},
	"ACTIVE":                intString,
	"WEBHOOK_URL":           intString,
	"WEBHOOK_SECRET":        intLeaf(propQName),
	"WEBHOOK_BODY_TEMPLATE": intString,
	"WEBHOOK_HEADERS":       {pairs: true},
	"ALLOWED_RECIPIENTS":    intLeaf(propStrings),
	"DEFAULT_RECIPIENTS":    intLeaf(propStrings),
	"DEFAULT_SUBJECT":       intString,
	"COMMENT":               intString,
}

// parseIntegrationProperties reads a run of `KEY = value` properties, each a
// key of props with a value of its shape, and no key twice.
func (p *Parser) parseIntegrationProperties(props integrationProps) ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	seen := map[string]bool{}
	for p.startsObjectProperty() {
		if seen[strings.ToUpper(p.cur.Str)] {
			return nil, p.syntaxErrorAtCur()
		}
		opt, err := p.parseIntegrationProperty(props)
		if err != nil {
			return nil, err
		}
		seen[opt.Name] = true
		opts = append(opts, opt)
	}
	return opts, nil
}

// parseIntegrationProperty reads one `KEY = value` property whose key props
// names. cur is the key, which startsObjectProperty has seen followed by '='.
func (p *Parser) parseIntegrationProperty(props integrationProps) (*ast.CopyOption, error) {
	spec, ok := props[strings.ToUpper(p.cur.Str)]
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	key := p.advance()
	opt := &ast.CopyOption{Name: strings.ToUpper(key.Str), Loc: ast.Loc{Start: key.Loc.Start}}
	p.advance() // consume '='
	var err error
	switch {
	case spec.groups:
		err = p.parseIntegrationGroupList(opt, spec.group)
	case spec.group != nil:
		err = p.parseIntegrationGroup(opt, spec.group)
	case spec.pairs:
		err = p.parseIntegrationStringPairs(opt)
	default:
		err = p.parseIntegrationLeaf(opt, spec.shapes)
	}
	if err != nil {
		return nil, err
	}
	if spec.oneParts {
		for _, e := range opt.Group {
			if nameParts(e.Name) != 1 {
				return nil, &ParseError{
					Position: e.Loc.Start, End: e.Loc.End,
					Message: "syntax error: " + e.Name + " takes no qualifier here",
				}
			}
		}
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parseIntegrationGroup reads ( <key> = <value> [ <key> = <value> ... ] ),
// entries separated by spaces as the pages write them, each a key of keys.
// cur is '('.
func (p *Parser) parseIntegrationGroup(opt *ast.CopyOption, keys integrationProps) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	if _, err := p.expect('('); err != nil {
		return err
	}
	if !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	entries, err := p.parseIntegrationProperties(keys)
	if err != nil {
		return err
	}
	opt.Group = entries
	_, err = p.expect(')')
	return err
}

// parseIntegrationGroupList reads ( ( ... ) [ , ( ... ) ... ] ), a list of
// groups of keys, as a multi-queue integration's QUEUES. cur is '('.
func (p *Parser) parseIntegrationGroupList(opt *ast.CopyOption, keys integrationProps) error {
	if err := p.enterNesting(); err != nil {
		return err
	}
	defer p.leaveNesting()
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		entry := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
		if err := p.parseIntegrationGroup(entry, keys); err != nil {
			return err
		}
		entry.Loc.End = p.prev.Loc.End
		opt.Group = append(opt.Group, entry)
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(')')
	return err
}

// parseIntegrationStringPairs reads ( '<key>' = '<value>' [ , ... ] ), as a
// webhook integration's WEBHOOK_HEADERS, each entry named by its quoted key.
// cur is '('.
func (p *Parser) parseIntegrationStringPairs(opt *ast.CopyOption) error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		key, err := p.expect(tokString)
		if err != nil {
			return err
		}
		if _, err := p.expect('='); err != nil {
			return err
		}
		val, err := p.parseCopyLiteral()
		if err != nil {
			return err
		}
		if val.Kind != ast.LitString {
			return &ParseError{Position: val.Loc.Start, End: val.Loc.End, Message: "syntax error: expected a string"}
		}
		// The key is kept single-quoted, however it was written ('...' or
		// $$...$$), so that it reads apart from a word key and writes back
		// as the same string.
		opt.Group = append(opt.Group, &ast.CopyOption{Name: quoteIntegrationKey(key.Str), Lit: val, Loc: ast.Loc{Start: key.Loc.Start, End: val.Loc.End}})
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(')')
	return err
}

// parseIntegrationLeaf reads a value that is no group into opt, and fails
// unless it has one of shapes:
//
//	'<string>' | <number> | <word>   -> Lit, or Words
//	( '<string>' [ , ... ] )          -> List
//	( <name> [ , ... ] )              -> Group of bare names
func (p *Parser) parseIntegrationLeaf(opt *ast.CopyOption, shapes []objectValue) error {
	start := p.cur.Loc.Start
	if !p.curIs('(') {
		if err := p.parseObjectValue(opt); err != nil {
			return err
		}
	} else {
		p.advance() // consume '('
		switch {
		case p.curIsLiteral():
			for {
				lit, err := p.parseCopyLiteral()
				if err != nil {
					return err
				}
				opt.List = append(opt.List, lit)
				if !p.curIs(',') {
					break
				}
				p.advance() // consume ','
			}
		case p.isOptionWord(p.cur.Type):
			for {
				nameStart := p.cur.Loc.Start
				name, err := p.parseObjectWord()
				if err != nil {
					return err
				}
				opt.Group = append(opt.Group, &ast.CopyOption{Name: name, Bare: true, Loc: ast.Loc{Start: nameStart, End: p.prev.Loc.End}})
				if !p.curIs(',') {
					break
				}
				p.advance() // consume ','
			}
		default:
			return p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return err
		}
	}
	for _, v := range shapes {
		if v.matches(opt) {
			return nil
		}
	}
	return &ParseError{
		Position: start, End: p.prev.Loc.End,
		Message: "syntax error: unexpected value for " + opt.Name,
	}
}

// parseAlterIntegrationSetProperties reads the SET of a CATALOG or EXTERNAL
// ACCESS INTEGRATION once SET is consumed:
//
//	ALTER CATALOG INTEGRATION [ IF EXISTS ] <name> SET
//	  [ REST_AUTHENTICATION = ( ... ) ] [ DEFAULT_STORAGE_CONFIG = ( ... ) ]
//	  [ REFRESH_INTERVAL_SECONDS = <value> ] [ COMMENT = '<string>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-catalog-integration
//
//	ALTER EXTERNAL ACCESS INTEGRATION [ IF EXISTS ] <name> SET
//	  [ ALLOWED_NETWORK_RULES = ( ... ) ] [ ... ] [ COMMENT = '<string>' ]
//	  [ TAG <tag_name> = '<tag_value>' [ , <tag_name> = '<tag_value>' ... ] ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration
func (p *Parser) parseAlterIntegrationSetProperties(stmt *ast.AlterIntegrationStmt) error {
	props := catalogIntegrationSetProps
	if stmt.Kind == ast.ExternalAccessIntegration {
		props = accessIntegrationProps
	}
	opts, err := p.parseIntegrationProperties(props)
	if err != nil {
		return err
	}
	if len(opts) == 0 {
		return p.syntaxErrorAtCur()
	}
	stmt.Action = ast.AlterIntegrationSet
	stmt.Options = opts
	if stmt.Kind == ast.ExternalAccessIntegration && p.curIs(kwTAG) {
		tags, err := p.parseTagSetList()
		if err != nil {
			return err
		}
		stmt.Tags = tags
	}
	return nil
}

// parseAlterAccessIntegrationUnset reads the UNSET list of an EXTERNAL ACCESS
// INTEGRATION once UNSET is consumed. Each item is one property the page
// names (accessIntegrationUnset) or TAG <tag_name>:
//
//	ALTER EXTERNAL ACCESS INTEGRATION [ IF EXISTS ] <name> UNSET {
//	  ALLOWED_NETWORK_RULES | ALLOWED_API_AUTHENTICATION_INTEGRATIONS |
//	  ALLOWED_AUTHENTICATION_SECRETS | COMMENT | TAG <tag_name> } [ , ... ]
//
// A list of tags alone reads as AlterIntegrationUnsetTag, as the other kinds'
// UNSET TAG does.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-access-integration
func (p *Parser) parseAlterAccessIntegrationUnset(stmt *ast.AlterIntegrationStmt) error {
	for {
		switch {
		case p.curIs(kwTAG):
			p.advance() // consume TAG
			name, err := p.parseObjectName()
			if err != nil {
				return err
			}
			stmt.UnsetTags = append(stmt.UnsetTags, name)
		case p.atAccessIntegrationUnsetKey():
			stmt.UnsetProps = append(stmt.UnsetProps, strings.ToUpper(p.advance().Str))
		default:
			return p.syntaxErrorAtCur()
		}
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	stmt.Action = ast.AlterIntegrationUnset
	if len(stmt.UnsetProps) == 0 {
		stmt.Action = ast.AlterIntegrationUnsetTag
	}
	return nil
}

// atAccessIntegrationUnsetKey reports whether cur is one of
// accessIntegrationUnset, as an unquoted identifier or a keyword. It consumes
// nothing.
func (p *Parser) atAccessIntegrationUnsetKey() bool {
	for _, w := range accessIntegrationUnset {
		if p.curWord(w) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Forms: which properties go together
// ---------------------------------------------------------------------------
//
// The property maps above read a CREATE's properties by their names and
// shapes across all of a kind's pages. Each page then documents one form:
// which of those properties it takes, which it requires, and the values a
// property takes on that page. A CATALOG INTEGRATION is held to a form of its
// CATALOG_SOURCE, a NOTIFICATION INTEGRATION to a form of its TYPE.

// integrationField is one property of a form: a leaf of shape value, or a
// group whose entries match one of variants (each entry of a list of groups,
// when list is set). parsed takes the value as the property map read it.
type integrationField struct {
	value    objectValue
	variants []integrationForm
	list     bool
	parsed   bool
	required bool
}

// integrationForm maps a form's property names to their fields.
type integrationForm map[string]integrationField

func reqField(v objectValue) integrationField { return integrationField{value: v, required: true} }
func optField(v objectValue) integrationField { return integrationField{value: v} }
func reqGroup(variants ...integrationForm) integrationField {
	return integrationField{variants: variants, required: true}
}
func optGroup(variants ...integrationForm) integrationField {
	return integrationField{variants: variants}
}

// matchesIntegrationForm reports whether opts set only form's properties,
// each with a value the form allows, and every property it requires.
func matchesIntegrationForm(opts []*ast.CopyOption, form integrationForm) bool {
	seen := map[string]bool{}
	for _, o := range opts {
		f, ok := form[o.Name]
		if !ok {
			return false
		}
		seen[o.Name] = true
		switch {
		case f.parsed:
		case f.list:
			for _, e := range o.Group {
				if !matchesAnyIntegrationForm(e.Group, f.variants) {
					return false
				}
			}
		case f.variants != nil:
			if !matchesAnyIntegrationForm(o.Group, f.variants) {
				return false
			}
		default:
			if !f.value.matches(o) {
				return false
			}
		}
	}
	for k, f := range form {
		if f.required && !seen[k] {
			return false
		}
	}
	return true
}

func matchesAnyIntegrationForm(opts []*ast.CopyOption, forms []integrationForm) bool {
	for _, form := range forms {
		if matchesIntegrationForm(opts, form) {
			return true
		}
	}
	return false
}

// withIntegrationFields returns form with fields added.
func withIntegrationFields(form integrationForm, fields integrationForm) integrationForm {
	out := integrationForm{}
	for k, v := range form {
		out[k] = v
	}
	for k, v := range fields {
		out[k] = v
	}
	return out
}

var (
	accessDelegationMode = propOneOf("VENDED_CREDENTIALS", "EXTERNAL_VOLUME_CREDENTIALS")
	vendedCredentials    = propOneOf("VENDED_CREDENTIALS")

	// catalogCommon is what every catalog source's page takes.
	catalogCommon = integrationForm{
		"CATALOG_SOURCE": reqField(propName),
		"ENABLED":        reqField(propBool),
		"COMMENT":        optField(propString),
	}

	// openCatalogAuth is Snowflake Open Catalog's REST_AUTHENTICATION;
	// withTokenURI says whether it takes OAUTH_TOKEN_URI.
	openCatalogAuth = func(withTokenURI bool) integrationForm {
		f := integrationForm{
			"TYPE":                 reqField(propOneOf("OAUTH")),
			"OAUTH_API_TYPE":       optField(propOneOf("PUBLIC", "PRIVATE")),
			"OAUTH_CLIENT_ID":      reqField(propString),
			"OAUTH_CLIENT_SECRET":  reqField(propString),
			"OAUTH_ALLOWED_SCOPES": reqField(propStrings),
		}
		if withTokenURI {
			f["OAUTH_TOKEN_URI"] = optField(propString)
		}
		return f
	}
	bearerAuth = integrationForm{"TYPE": reqField(propOneOf("BEARER")), "BEARER_TOKEN": reqField(propString)}

	// catalogSourceForms maps a CATALOG_SOURCE to the forms its pages give.
	catalogSourceForms = map[string][]integrationForm{
		// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-glue
		"GLUE": {withIntegrationFields(catalogCommon, integrationForm{
			"TABLE_FORMAT":      reqField(propOneOf("ICEBERG")),
			"CATALOG_NAMESPACE": optField(propString),
			"REST_CONFIG": reqGroup(integrationForm{
				"CATALOG_API_TYPE":       reqField(propOneOf("AWS_GLUE", "AWS_PRIVATE_GLUE")),
				"CATALOG_NAME":           reqField(propString),
				"ACCESS_DELEGATION_MODE": optField(accessDelegationMode),
			}),
			"REST_AUTHENTICATION": reqGroup(integrationForm{
				"TYPE":                 reqField(propOneOf("SIGV4")),
				"SIGV4_IAM_ROLE":       reqField(propString),
				"SIGV4_SIGNING_REGION": optField(propString),
			}),
			"REFRESH_INTERVAL_SECONDS": optField(propNumber),
		})},
		// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-object-storage
		"OBJECT_STORE": {withIntegrationFields(catalogCommon, integrationForm{
			"TABLE_FORMAT":             reqField(propOneOf("ICEBERG", "DELTA", "NONE", "HIVE")),
			"REFRESH_INTERVAL_SECONDS": optField(propNumber),
		})},
		"POLARIS": {
			// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-open-catalog
			withIntegrationFields(catalogCommon, integrationForm{
				"TABLE_FORMAT":      reqField(propOneOf("ICEBERG")),
				"CATALOG_NAMESPACE": optField(propString),
				"REST_CONFIG": reqGroup(integrationForm{
					"CATALOG_URI":            reqField(propString),
					"CATALOG_API_TYPE":       optField(propOneOf("PUBLIC")),
					"CATALOG_NAME":           reqField(propString),
					"ACCESS_DELEGATION_MODE": optField(accessDelegationMode),
				}),
				"REST_AUTHENTICATION":      reqGroup(openCatalogAuth(true)),
				"REFRESH_INTERVAL_SECONDS": optField(propNumber),
			}),
			withIntegrationFields(catalogCommon, integrationForm{
				"TABLE_FORMAT":      reqField(propOneOf("ICEBERG")),
				"CATALOG_NAMESPACE": optField(propString),
				"REST_CONFIG": reqGroup(integrationForm{
					"CATALOG_URI":            reqField(propString),
					"CATALOG_API_TYPE":       reqField(propOneOf("PRIVATE")),
					"CATALOG_NAME":           reqField(propString),
					"ACCESS_DELEGATION_MODE": optField(accessDelegationMode),
				}),
				"REST_AUTHENTICATION":      reqGroup(openCatalogAuth(false)),
				"REFRESH_INTERVAL_SECONDS": optField(propNumber),
			}),
			// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-horizon
			withIntegrationFields(catalogCommon, integrationForm{
				"TABLE_FORMAT": reqField(propOneOf("ICEBERG")),
				"REST_CONFIG": reqGroup(integrationForm{
					"CATALOG_URI":            reqField(propString),
					"CATALOG_API_TYPE":       optField(propOneOf("PUBLIC", "PRIVATE")),
					"CATALOG_NAME":           reqField(propString),
					"ACCESS_DELEGATION_MODE": optField(accessDelegationMode),
				}),
				"REST_AUTHENTICATION": reqGroup(integrationForm{
					"TYPE":                 reqField(propOneOf("OAUTH")),
					"OAUTH_GRANT_TYPE":     reqField(propOneOf("TOKEN_EXCHANGE")),
					"OAUTH_ALLOWED_SCOPES": reqField(propStrings),
					"OAUTH_AUDIENCE":       optField(propString),
				}),
				"DEFAULT_STORAGE_CONFIG":   optGroup(catalogDefaultStorageForm),
				"REFRESH_INTERVAL_SECONDS": optField(propNumber),
			}),
		},
		// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-rest
		"ICEBERG_REST": {withIntegrationFields(catalogCommon, integrationForm{
			"TABLE_FORMAT":      reqField(propOneOf("ICEBERG")),
			"CATALOG_NAMESPACE": optField(propString),
			"REST_CONFIG": reqGroup(integrationForm{
				"CATALOG_URI":  reqField(propString),
				"PREFIX":       optField(propString),
				"CATALOG_NAME": optField(propString),
				"CATALOG_API_TYPE": optField(propOneOf("PUBLIC", "PRIVATE", "AWS_API_GATEWAY", "AWS_PRIVATE_API_GATEWAY",
					"AWS_GLUE", "AWS_PRIVATE_GLUE", "AWS_S3TABLES", "AWS_PRIVATE_S3TABLES")),
				"ACCESS_DELEGATION_MODE": optField(accessDelegationMode),
			}),
			"REST_AUTHENTICATION": reqGroup(
				openCatalogAuth(true),
				bearerAuth,
				integrationForm{
					"TYPE":                 reqField(propOneOf("SIGV4")),
					"SIGV4_IAM_ROLE":       reqField(propString),
					"SIGV4_SIGNING_REGION": optField(propString),
					"SIGV4_EXTERNAL_ID":    optField(propString),
				},
			),
			"REFRESH_INTERVAL_SECONDS": optField(propNumber),
		})},
		// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-snowflake-postgres
		"SNOWFLAKE_POSTGRES": {withIntegrationFields(catalogCommon, integrationForm{
			"TABLE_FORMAT":      reqField(propOneOf("ICEBERG")),
			"CATALOG_NAMESPACE": optField(propString),
			"REST_CONFIG": reqGroup(integrationForm{
				"POSTGRES_INSTANCE":      reqField(propString),
				"ACCESS_DELEGATION_MODE": reqField(vendedCredentials),
				"CATALOG_NAME":           optField(propString),
			}),
		})},
		// https://docs.snowflake.com/en/sql-reference/sql/create-catalog-integration-delta-sharing
		"DELTA_SHARING": {withIntegrationFields(catalogCommon, integrationForm{
			"TABLE_FORMAT": reqField(propOneOf("DELTA")),
			"REST_CONFIG": reqGroup(integrationForm{
				"CATALOG_URI":            reqField(propString),
				"CATALOG_NAME":           reqField(propString),
				"ACCESS_DELEGATION_MODE": reqField(vendedCredentials),
			}),
			"REST_AUTHENTICATION": reqGroup(
				bearerAuth,
				integrationForm{"TYPE": reqField(propOneOf("OIDC")), "OIDC_AUDIENCE": reqField(propString)},
				integrationForm{
					"TYPE":                reqField(propOneOf("OAUTH")),
					"OAUTH_API_TYPE":      optField(propOneOf("PUBLIC", "PRIVATE")),
					"OAUTH_CLIENT_ID":     reqField(propString),
					"OAUTH_CLIENT_SECRET": reqField(propString),
					"OAUTH_TOKEN_URI":     reqField(propString),
				},
			),
		})},
	}

	catalogDefaultStorageForm = integrationForm{"USE_PRIVATELINK_ENDPOINT": reqField(propBool)}

	// notificationCommon is what every notification type's page takes.
	notificationCommon = integrationForm{
		"ENABLED": reqField(propBool),
		"TYPE":    reqField(propName),
		"COMMENT": optField(propString),
	}

	// notificationTypeForms maps a notification integration's TYPE to the
	// forms its pages give.
	notificationTypeForms = map[string][]integrationForm{
		"QUEUE": {
			// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-azure
			withIntegrationFields(notificationCommon, integrationForm{
				"NOTIFICATION_PROVIDER":           reqField(propOneOf("AZURE_STORAGE_QUEUE")),
				"AZURE_STORAGE_QUEUE_PRIMARY_URI": reqField(propString),
				"AZURE_TENANT_ID":                 reqField(propString),
				"USE_PRIVATELINK_ENDPOINT":        optField(propBool),
			}),
			// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-inbound-gcp
			withIntegrationFields(notificationCommon, integrationForm{
				"NOTIFICATION_PROVIDER":        reqField(propOneOf("GCP_PUBSUB")),
				"GCP_PUBSUB_SUBSCRIPTION_NAME": reqField(propString),
			}),
			// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-aws
			withIntegrationFields(notificationCommon, integrationForm{
				"DIRECTION":             reqField(propOneOf("OUTBOUND")),
				"NOTIFICATION_PROVIDER": reqField(propOneOf("AWS_SNS")),
				"AWS_SNS_TOPIC_ARN":     reqField(propString),
				"AWS_SNS_ROLE_ARN":      reqField(propString),
			}),
			// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-azure
			withIntegrationFields(notificationCommon, integrationForm{
				"DIRECTION":                       reqField(propOneOf("OUTBOUND")),
				"NOTIFICATION_PROVIDER":           reqField(propOneOf("AZURE_EVENT_GRID")),
				"AZURE_EVENT_GRID_TOPIC_ENDPOINT": reqField(propString),
				"AZURE_TENANT_ID":                 reqField(propString),
			}),
			// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-queue-outbound-gcp
			withIntegrationFields(notificationCommon, integrationForm{
				"DIRECTION":             reqField(propOneOf("OUTBOUND")),
				"NOTIFICATION_PROVIDER": reqField(propOneOf("GCP_PUBSUB")),
				"GCP_PUBSUB_TOPIC_NAME": reqField(propString),
			}),
		},
		// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-multi-queue
		"MULTI_QUEUE": {withIntegrationFields(notificationCommon, integrationForm{
			"DIRECTION": reqField(propOneOf("INBOUND")),
			"QUEUES": {list: true, required: true, variants: []integrationForm{
				{
					"NAME":                  reqField(propString),
					"NOTIFICATION_PROVIDER": reqField(propOneOf("AWS_SNS")),
					"AWS_SNS_TOPIC_ARN":     reqField(propString),
				},
				{
					"NAME":                         reqField(propString),
					"NOTIFICATION_PROVIDER":        reqField(propOneOf("GCP_PUBSUB")),
					"GCP_PUBSUB_SUBSCRIPTION_NAME": reqField(propString),
				},
				{
					"NAME":                            reqField(propString),
					"NOTIFICATION_PROVIDER":           reqField(propOneOf("AZURE_STORAGE_QUEUE")),
					"AZURE_STORAGE_QUEUE_PRIMARY_URI": reqField(propString),
					"AZURE_TENANT_ID":                 reqField(propString),
				},
			}},
			"ACTIVE": reqField(propString),
		})},
		// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-webhooks
		"WEBHOOK": {withIntegrationFields(notificationCommon, integrationForm{
			"WEBHOOK_URL":           reqField(propString),
			"WEBHOOK_SECRET":        optField(propQName),
			"WEBHOOK_BODY_TEMPLATE": optField(propString),
			"WEBHOOK_HEADERS":       {parsed: true},
		})},
		// https://docs.snowflake.com/en/sql-reference/sql/create-notification-integration-email
		"EMAIL": {withIntegrationFields(notificationCommon, integrationForm{
			"ALLOWED_RECIPIENTS": optField(propStrings),
			"DEFAULT_RECIPIENTS": optField(propStrings),
			"DEFAULT_SUBJECT":    optField(propString),
		})},
	}
)

// checkIntegrationForm fails unless a CATALOG or NOTIFICATION INTEGRATION's
// properties match a form its CATALOG_SOURCE or TYPE documents. start is the
// statement's start, where a mismatch is reported.
func checkIntegrationForm(kind ast.IntegrationObjectKind, opts []*ast.CopyOption, start, end int) error {
	var forms []integrationForm
	switch kind {
	case ast.CatalogIntegration:
		forms = catalogSourceForms[integrationOptionWords(opts, "CATALOG_SOURCE")]
	case ast.NotificationIntegration:
		forms = notificationTypeForms[integrationOptionWords(opts, "TYPE")]
	default:
		return nil
	}
	if matchesAnyIntegrationForm(opts, forms) {
		return nil
	}
	return &ParseError{
		Position: start, End: end,
		Message: "syntax error: the properties of " + kind.String() + " match no form its documentation gives",
	}
}

// integrationOptionWords returns the word value of the option named name.
func integrationOptionWords(opts []*ast.CopyOption, name string) string {
	for _, o := range opts {
		if o.Name == name {
			return o.Words
		}
	}
	return ""
}

// quoteIntegrationKey spells a string key as a single-quoted literal that
// reads back as s.
func quoteIntegrationKey(s string) string {
	return "'" + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", "''") + "'"
}
