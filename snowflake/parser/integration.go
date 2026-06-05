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
//	{ STORAGE | API | NOTIFICATION | SECURITY } INTEGRATION
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
// RESOURCE / SECRET / CONNECTION / GIT). EXTERNAL VOLUME enters through
// parseCreateExternalVolumeStmt instead (its dispatch already consumed EXTERNAL +
// VOLUME).
func (p *Parser) parseCreateIntegrationStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	kind, err := p.consumeIntegrationObjectKeyword()
	if err != nil {
		return nil, err
	}

	stmt := &ast.CreateIntegrationStmt{
		Kind:      kind,
		OrReplace: orReplace,
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
		if p.cur.Type == kwTRIGGERS {
			triggers, err := p.parseResourceMonitorTriggers()
			if err != nil {
				return nil, err
			}
			stmt.Triggers = triggers
		}

	case ast.Connection:
		// CREATE CONNECTION <name> [ AS REPLICA OF a.b.c ] [ COMMENT = '...' ]
		if p.cur.Type == kwAS {
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

	default:
		// STORAGE/API/NOTIFICATION/SECURITY INTEGRATION, SECRET, GIT REPOSITORY:
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
func (p *Parser) parseCreateExternalVolumeStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	stmt := &ast.CreateIntegrationStmt{
		Kind:      ast.ExternalVolume,
		OrReplace: orReplace,
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

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwSET:
		if err := p.parseAlterIntegrationSet(stmt); err != nil {
			return nil, err
		}

	case kwUNSET:
		if err := p.parseAlterIntegrationUnset(stmt); err != nil {
			return nil, err
		}

	case kwFETCH:
		// GIT REPOSITORY FETCH.
		p.advance() // consume FETCH
		stmt.Action = ast.AlterIntegrationFetch

	case kwENABLE, kwDISABLE:
		// CONNECTION ENABLE/DISABLE FAILOVER [ TO ACCOUNTS a.b [, ...] ].
		if err := p.parseAlterConnectionFailover(stmt); err != nil {
			return nil, err
		}

	case kwPRIMARY:
		// CONNECTION PRIMARY.
		p.advance() // consume PRIMARY
		stmt.Action = ast.AlterIntegrationPrimary

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
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

	switch p.cur.Type {
	case kwADD:
		// ADD STORAGE_LOCATION = ( NAME=... <provider params> )
		p.advance() // consume ADD
		opt, err := p.parseStorageLocationSingleOption()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterIntegrationAddLocation
		stmt.Options = []*ast.CopyOption{opt}

	case kwREMOVE:
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

	case kwUPDATE:
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

	case kwSET:
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

	if p.cur.Type == kwTAG {
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
	if p.cur.Type == kwTRIGGERS {
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

	if p.cur.Type == kwTAG {
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
	enable := p.cur.Type == kwENABLE
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
		if p.cur.Type == kwIGNORE {
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
	if p.cur.Type == kwTO {
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
		if p.cur.Type == ',' {
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
	for p.cur.Type == kwON {
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

	switch p.cur.Type {
	case kwSUSPEND:
		p.advance()
		trig.Action = "SUSPEND"
	case kwSUSPEND_IMMEDIATE:
		p.advance()
		trig.Action = "SUSPEND_IMMEDIATE"
	case kwNOTIFY:
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
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		inner, err := p.parseStorageLocationGroup()
		if err != nil {
			return nil, err
		}
		groups = append(groups, inner)
		if p.cur.Type == ',' {
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
	if p.cur.Type != '(' {
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
	switch p.cur.Type {
	case kwSTORAGE:
		p.advance() // consume STORAGE
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.StorageIntegration, nil
	case kwAPI:
		p.advance() // consume API
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.APIIntegration, nil
	case kwNOTIFICATION:
		p.advance() // consume NOTIFICATION
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.NotificationIntegration, nil
	case kwSECURITY:
		p.advance() // consume SECURITY
		if _, err := p.expect(kwINTEGRATION); err != nil {
			return 0, err
		}
		return ast.SecurityIntegration, nil
	case kwINTEGRATION:
		// Bare ALTER INTEGRATION <name> ... — subtype qualifier omitted.
		p.advance() // consume INTEGRATION
		return ast.StorageIntegration, nil
	case kwRESOURCE:
		p.advance() // consume RESOURCE
		if _, err := p.expect(kwMONITOR); err != nil {
			return 0, err
		}
		return ast.ResourceMonitor, nil
	case kwSECRET:
		p.advance() // consume SECRET
		return ast.Secret, nil
	case kwCONNECTION:
		p.advance() // consume CONNECTION
		return ast.Connection, nil
	case kwGIT:
		p.advance() // consume GIT
		if _, err := p.expect(kwREPOSITORY); err != nil {
			return 0, err
		}
		return ast.GitRepository, nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
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

// startsIntegrationOption reports whether the current token can begin an option.
// It is the COPY option predicate minus the structural keywords that anchor a
// trailing clause: TRIGGERS (RESOURCE MONITOR), and WITH / TAG (GIT REPOSITORY's
// [WITH] TAG). Those must not be swallowed as an option name.
func (p *Parser) startsIntegrationOption() bool {
	switch p.cur.Type {
	case kwTRIGGERS, kwWITH, kwTAG:
		return false
	}
	return p.startsCopyOption()
}

// startsIntegrationTags reports whether the current position begins a GIT
// REPOSITORY [WITH] TAG clause (TAG directly, or WITH immediately followed by
// TAG). Reuses the STAGE predicate semantics.
func (p *Parser) startsIntegrationTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
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
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwNOT {
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
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			p.advance() // consume IF
			p.advance() // consume EXISTS
			*flag = true
		}
	}
}
