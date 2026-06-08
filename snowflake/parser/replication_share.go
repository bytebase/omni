package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Replication & sharing DDL — CREATE / ALTER FAILOVER GROUP / REPLICATION GROUP
// / ACCOUNT / SHARE (T4.8)
// ---------------------------------------------------------------------------
//
// These account-level objects drive Snowflake's replication and data-sharing
// features. Like the rest of this engine's account-level DDL (STAGE T4.1, FILE
// FORMAT T4.2, integrations T4.7) they carry large, version-growing
// `KEY = <value>` option vocabularies. The distinctive twist — and the reason
// the COPY/STAGE option machinery is NOT reused for them — is that the
// replication/sharing list parameters (OBJECT_TYPES, ALLOWED_DATABASES,
// ALLOWED_SHARES, ALLOWED_INTEGRATION_TYPES, ALLOWED_ACCOUNTS, ...) are
// UNPARENTHESIZED comma lists of multi-word object-type names (`ACCOUNT
// PARAMETERS`, `RESOURCE MONITORS`, `NETWORK POLICIES`) or dotted account names
// (`org.account`). Both the official docs (truth1) and the legacy ANTLR grammar
// (truth2) agree on the no-parentheses form (`OBJECT_TYPES = X, Y` not
// `OBJECT_TYPES = (X, Y)`), so these are captured as ast.GroupOption (a
// `KEY = <value-list | literal>` pair) rather than ast.CopyOption.
//
// Only structural anchors are modeled as dedicated fields: the object kind
// (FAILOVER vs REPLICATION GROUP; MANAGED ACCOUNT), the [WITH] TAG clause, the
// AS REPLICA OF secondary form, and (on ALTER) the action keyword plus its
// operands (ADD/REMOVE/MOVE name lists, the ALLOWED_* target, the MOVE-TO group,
// RENAME's new name, the account-policy forms, ...). The catalog/semantic layer,
// not the parser, validates that an option is real and legal for the chosen
// object (mirroring the merged integration / stage parsers).
//
// Scope notes (per the work order):
//   - CREATE OR ALTER (a preview form) is owned by parser-or-alter; the shared
//     OR-prefix parser recognizes only OR REPLACE, so OR ALTER never reaches
//     here.
//   - $N positional bind references inside option values are owned by
//     expr-dollar-refs and are out of scope.

// ---------------------------------------------------------------------------
// CREATE FAILOVER GROUP / REPLICATION GROUP
// ---------------------------------------------------------------------------

// parseCreateReplicationGroupStmt parses the body of a CREATE statement for a
// FAILOVER GROUP or a REPLICATION GROUP. The CREATE keyword and optional OR
// REPLACE modifier have already been consumed by parseCreateStmt; start is the
// Loc of the CREATE token, and cur is the FAILOVER or REPLICATION keyword (which
// is followed by the GROUP keyword).
//
//	CREATE [ OR REPLACE ] { FAILOVER | REPLICATION } GROUP [ IF NOT EXISTS ] <name>
//	  OBJECT_TYPES = <object_type> [ , ... ]  <list/scalar options...>
//	  ALLOWED_ACCOUNTS = <org>.<account> [ , ... ]  [ IGNORE EDITION CHECK ]  ...
//	  [ [ WITH ] TAG (...) ]
//	CREATE [ OR REPLACE ] { FAILOVER | REPLICATION } GROUP [ IF NOT EXISTS ] <name>
//	  AS REPLICA OF <org>.<source_account>.<name>
func (p *Parser) parseCreateReplicationGroupStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	failover := p.cur.Type == kwFAILOVER
	p.advance() // consume FAILOVER / REPLICATION
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	stmt := &ast.CreateReplicationGroupStmt{
		Failover:  failover,
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

	// Secondary form: AS REPLICA OF <org>.<account>.<name>.
	if p.cur.Type == kwAS {
		p.advance() // consume AS
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
		stmt.Replica = src
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Primary form: an open-ended run of list/scalar options, the bare IGNORE
	// EDITION CHECK flag, and the trailing [WITH] TAG clause, in any order (so a
	// REPLICATION_SCHEDULE that trails the TAG clause is not dropped — mirrors the
	// STAGE / integration loops).
	for {
		if p.startsGroupTags() {
			tags, err := p.parseStageWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		if p.cur.Type == kwIGNORE {
			if err := p.parseIgnoreEditionCheck(); err != nil {
				return nil, err
			}
			stmt.IgnoreEditionCheck = true
			continue
		}
		if p.startsGroupOption() {
			opt, err := p.parseGroupOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER FAILOVER GROUP / REPLICATION GROUP
// ---------------------------------------------------------------------------

// parseAlterReplicationGroupStmt parses the body of an ALTER statement for a
// FAILOVER GROUP or REPLICATION GROUP. The ALTER keyword has already been
// consumed; cur is the FAILOVER or REPLICATION keyword (followed by GROUP).
func (p *Parser) parseAlterReplicationGroupStmt() (ast.Node, error) {
	startLoc := p.cur.Loc // object keyword anchors Loc.Start (ALTER convention)
	failover := p.cur.Type == kwFAILOVER
	p.advance() // consume FAILOVER / REPLICATION
	if _, err := p.expect(kwGROUP); err != nil {
		return nil, err
	}

	stmt := &ast.AlterReplicationGroupStmt{Failover: failover, Loc: ast.Loc{Start: startLoc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwRENAME:
		// RENAME TO <new_name>
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterGroupRename
		stmt.NewName = newName

	case kwSET:
		if err := p.parseAlterGroupSet(stmt); err != nil {
			return nil, err
		}

	case kwUNSET:
		if err := p.parseAlterGroupUnset(stmt); err != nil {
			return nil, err
		}

	case kwADD:
		// ADD <name-list> TO ALLOWED_{DATABASES|SHARES|ACCOUNTS} [ IGNORE EDITION CHECK ]
		p.advance() // consume ADD
		names, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		target, err := p.parseAllowedListKeyword()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterGroupAdd
		stmt.Names = names
		stmt.ListTarget = target
		if p.cur.Type == kwIGNORE {
			if err := p.parseIgnoreEditionCheck(); err != nil {
				return nil, err
			}
			stmt.IgnoreEditionCheck = true
		}

	case kwREMOVE:
		// REMOVE <name-list> FROM ALLOWED_{DATABASES|SHARES|ACCOUNTS}
		p.advance() // consume REMOVE
		names, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, err
		}
		target, err := p.parseAllowedListKeyword()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterGroupRemove
		stmt.Names = names
		stmt.ListTarget = target

	case kwMOVE:
		// MOVE { DATABASES | SHARES } <name-list> TO { FAILOVER | REPLICATION } GROUP <name>
		p.advance() // consume MOVE
		switch p.cur.Type {
		case kwDATABASES:
			p.advance()
			stmt.MoveKind = "DATABASES"
		case kwSHARES:
			p.advance()
			stmt.MoveKind = "SHARES"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		names, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		// TO { FAILOVER | REPLICATION } GROUP <name>.
		switch p.cur.Type {
		case kwFAILOVER, kwREPLICATION:
			p.advance()
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(kwGROUP); err != nil {
			return nil, err
		}
		moveTo, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterGroupMove
		stmt.Names = names
		stmt.MoveTo = moveTo

	case kwREFRESH:
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterGroupRefresh

	case kwPRIMARY:
		p.advance() // consume PRIMARY
		stmt.Action = ast.AlterGroupPrimary

	case kwSUSPEND:
		p.advance() // consume SUSPEND
		stmt.Action = ast.AlterGroupSuspend
		if p.cur.Type == kwIMMEDIATE {
			p.advance() // consume IMMEDIATE
			stmt.Immediate = true
		}

	case kwRESUME:
		p.advance() // consume RESUME
		stmt.Action = ast.AlterGroupResume

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterGroupSet parses the SET branch of ALTER FAILOVER/REPLICATION GROUP.
// cur is the SET keyword.
//
//	SET TAG <tag> = '<value>' [ , ... ]   -> AlterGroupSetTag
//	SET <options>                         -> AlterGroupSet
func (p *Parser) parseAlterGroupSet(stmt *ast.AlterReplicationGroupStmt) error {
	p.advance() // consume SET

	if p.cur.Type == kwTAG {
		tags, err := p.parseTagSetList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterGroupSetTag
		stmt.Tags = tags
		return nil
	}

	opts, err := p.parseGroupOptions()
	if err != nil {
		return err
	}
	if len(opts) == 0 {
		// SET with nothing settable is a syntax error.
		return p.syntaxErrorAtCur()
	}
	stmt.Action = ast.AlterGroupSet
	stmt.Options = opts
	return nil
}

// parseAlterGroupUnset parses the UNSET branch of ALTER FAILOVER/REPLICATION
// GROUP. cur is the UNSET keyword.
//
//	UNSET TAG <tag> [ , ... ]    -> AlterGroupUnsetTag
//	UNSET <property> [ , ... ]   -> AlterGroupUnset (e.g. COMMENT / REPLICATION_SCHEDULE)
func (p *Parser) parseAlterGroupUnset(stmt *ast.AlterReplicationGroupStmt) error {
	p.advance() // consume UNSET

	if p.cur.Type == kwTAG {
		names, err := p.parseUnsetTagNameList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterGroupUnsetTag
		stmt.UnsetTags = names
		return nil
	}

	props, err := p.parseStageUnsetProps()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterGroupUnset
	stmt.UnsetProps = props
	return nil
}

// ---------------------------------------------------------------------------
// CREATE ACCOUNT / MANAGED ACCOUNT
// ---------------------------------------------------------------------------

// parseCreateAccountStmt parses CREATE ACCOUNT / CREATE [OR REPLACE] MANAGED
// ACCOUNT. The CREATE keyword and optional OR REPLACE modifier have already been
// consumed; start is the Loc of the CREATE token, and cur is either the ACCOUNT
// keyword or the MANAGED keyword (followed by ACCOUNT).
//
//	CREATE ACCOUNT <name> ADMIN_NAME = ... <space-separated params...>
//	CREATE [ OR REPLACE ] MANAGED ACCOUNT <name> ADMIN_NAME = ..., <comma-separated params...>
func (p *Parser) parseCreateAccountStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	managed := p.cur.Type == kwMANAGED
	if managed {
		p.advance() // consume MANAGED
	}
	if _, err := p.expect(kwACCOUNT); err != nil {
		return nil, err
	}

	stmt := &ast.CreateAccountStmt{
		Managed:   managed,
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// The parameter list is an open-ended run of `KEY = value` pairs. CREATE
	// ACCOUNT separates them with whitespace; CREATE MANAGED ACCOUNT separates
	// them with commas. parseGroupOptions tolerates both (it consumes an optional
	// comma between options).
	opts, err := p.parseGroupOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER ACCOUNT
// ---------------------------------------------------------------------------

// parseAlterAccountStmt parses ALTER ACCOUNT. The ALTER keyword has already been
// consumed; cur is the ACCOUNT keyword. The current-account forms omit the name
// (Name stays nil); the cross-account forms (org admins) carry a name. The
// disambiguation is structural: if the token after ACCOUNT is SET or UNSET, it
// is a current-account form; otherwise it is a name followed by an action.
//
//	ALTER ACCOUNT SET <param> = <value> [ , ... ]   |   SET RESOURCE_MONITOR = <name>
//	ALTER ACCOUNT UNSET <param> [ , ... ]
//	ALTER ACCOUNT SET TAG ...   |   UNSET TAG ...
//	ALTER ACCOUNT SET { <kind> } POLICY <name> [ <scope> ] [ FORCE ]
//	ALTER ACCOUNT UNSET { <kind> } POLICY [ <scope> ]
//	ALTER ACCOUNT <name> SET <param> = <value> [ , ... ]
//	ALTER ACCOUNT <name> RENAME TO <new_name> [ SAVE_OLD_URL = { TRUE | FALSE } ]
//	ALTER ACCOUNT <name> DROP OLD [ ORGANIZATION ] URL
func (p *Parser) parseAlterAccountStmt() (ast.Node, error) {
	startLoc := p.cur.Loc // ACCOUNT keyword anchors Loc.Start (ALTER convention)
	p.advance()           // consume ACCOUNT

	stmt := &ast.AlterAccountStmt{Loc: ast.Loc{Start: startLoc.Start}}

	// Cross-account form: a name precedes the action when the next token is not
	// SET / UNSET.
	if p.cur.Type != kwSET && p.cur.Type != kwUNSET {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	switch p.cur.Type {
	case kwSET:
		if err := p.parseAlterAccountSet(stmt); err != nil {
			return nil, err
		}

	case kwUNSET:
		if err := p.parseAlterAccountUnset(stmt); err != nil {
			return nil, err
		}

	case kwRENAME:
		// <name> RENAME TO <new_name> [ SAVE_OLD_URL = { TRUE | FALSE } ].
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAccountRename
		stmt.NewName = newName
		if p.cur.Type == kwSAVE_OLD_URL {
			p.advance() // consume SAVE_OLD_URL
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			b, err := p.parseTrueFalse()
			if err != nil {
				return nil, err
			}
			stmt.SaveOldURL = &b
		}

	case kwDROP:
		// <name> DROP OLD [ ORGANIZATION ] URL.
		p.advance() // consume DROP
		if _, err := p.expect(kwOLD); err != nil {
			return nil, err
		}
		if p.cur.Type == kwORGANIZATION {
			p.advance() // consume ORGANIZATION
			stmt.Organization = true
		}
		if _, err := p.expect(kwURL); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterAccountDropURL

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterAccountSet parses the SET branch of ALTER ACCOUNT. cur is the SET
// keyword.
//
//	SET TAG <tag> = '<value>' [ , ... ]                      -> AlterAccountSetTag
//	SET { <kind> } POLICY <name> [ <scope> ] [ FORCE ]       -> AlterAccountSetPolicy
//	SET <param> = <value> [ , ... ] / SET RESOURCE_MONITOR=. -> AlterAccountSet
func (p *Parser) parseAlterAccountSet(stmt *ast.AlterAccountStmt) error {
	p.advance() // consume SET

	if p.cur.Type == kwTAG {
		tags, err := p.parseTagSetList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterAccountSetTag
		stmt.Tags = tags
		return nil
	}

	if p.startsAccountPolicy() {
		return p.parseAlterAccountSetPolicy(stmt)
	}

	opts, err := p.parseGroupOptions()
	if err != nil {
		return err
	}
	if len(opts) == 0 {
		return p.syntaxErrorAtCur()
	}
	stmt.Action = ast.AlterAccountSet
	stmt.Options = opts
	return nil
}

// parseAlterAccountUnset parses the UNSET branch of ALTER ACCOUNT. cur is the
// UNSET keyword.
//
//	UNSET TAG <tag> [ , ... ]              -> AlterAccountUnsetTag
//	UNSET { <kind> } POLICY [ <scope> ]    -> AlterAccountUnsetPolicy
//	UNSET <param> [ , ... ]                -> AlterAccountUnset
func (p *Parser) parseAlterAccountUnset(stmt *ast.AlterAccountStmt) error {
	p.advance() // consume UNSET

	if p.cur.Type == kwTAG {
		names, err := p.parseUnsetTagNameList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterAccountUnsetTag
		stmt.UnsetTags = names
		return nil
	}

	if p.startsAccountPolicy() {
		kind, err := p.parseAccountPolicyKind()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterAccountUnsetPolicy
		stmt.PolicyKind = kind
		// Optional trailing scope (FOR ALL ... ).
		if p.cur.Type == kwFOR {
			scope, err := p.parseAccountPolicyScope()
			if err != nil {
				return err
			}
			stmt.PolicyScope = scope
		}
		return nil
	}

	props, err := p.parseStageUnsetProps()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterAccountUnset
	stmt.UnsetProps = props
	return nil
}

// parseAlterAccountSetPolicy parses the SET { <kind> } POLICY <name> [ <scope> ]
// [ FORCE ] form. cur is the first policy-kind word (the caller confirmed
// startsAccountPolicy).
func (p *Parser) parseAlterAccountSetPolicy(stmt *ast.AlterAccountStmt) error {
	kind, err := p.parseAccountPolicyKind()
	if err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterAccountSetPolicy
	stmt.PolicyKind = kind
	stmt.PolicyName = name

	// Optional FOR ALL ... scope.
	if p.cur.Type == kwFOR {
		scope, err := p.parseAccountPolicyScope()
		if err != nil {
			return err
		}
		stmt.PolicyScope = scope
	}
	// Optional trailing FORCE.
	if p.cur.Type == kwFORCE {
		p.advance() // consume FORCE
		stmt.Force = true
	}
	return nil
}

// startsAccountPolicy reports whether the current position begins an ALTER
// ACCOUNT policy form, i.e. a policy-kind word (AUTHENTICATION / SESSION /
// PACKAGES / PASSWORD / FEATURE) immediately followed by the POLICY keyword.
// AUTHENTICATION and FEATURE are not reserved keywords, so they arrive as plain
// identifiers; SESSION / PACKAGES / PASSWORD are keywords. The next-token POLICY
// check is what distinguishes the policy form from an ordinary
// `SET <param> = value` whose param happens to share a prefix.
func (p *Parser) startsAccountPolicy() bool {
	switch {
	case p.cur.Type == kwSESSION, p.cur.Type == kwPACKAGES, p.cur.Type == kwPASSWORD:
		return p.peekNext().Type == kwPOLICY
	case p.curIsWord("AUTHENTICATION"), p.curIsWord("FEATURE"):
		return p.peekNext().Type == kwPOLICY
	}
	return false
}

// parseAccountPolicyKind consumes the policy-kind word(s) and the POLICY keyword
// and returns the uppercased "<KIND> POLICY" text (e.g. "PACKAGES POLICY"). The
// caller has confirmed startsAccountPolicy.
func (p *Parser) parseAccountPolicyKind() (string, error) {
	kindTok := p.advance() // consume the policy-kind word
	if _, err := p.expect(kwPOLICY); err != nil {
		return "", err
	}
	return strings.ToUpper(kindTok.Str) + " POLICY", nil
}

// parseAccountPolicyScope parses a trailing policy scope clause and returns it
// verbatim (uppercased, space-joined). cur is the FOR keyword. The two documented
// shapes are `FOR ALL { PERSON | SERVICE } USERS` and `FOR ALL APPLICATIONS`;
// rather than enumerate them, the run of words after FOR is captured whole (it
// runs to the FORCE keyword or the statement boundary).
func (p *Parser) parseAccountPolicyScope() (string, error) {
	if _, err := p.expect(kwFOR); err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("FOR")
	for p.isOptionWord(p.cur.Type) && p.cur.Type != kwFORCE {
		b.WriteByte(' ')
		b.WriteString(strings.ToUpper(p.advance().Str))
	}
	return b.String(), nil
}

// ---------------------------------------------------------------------------
// CREATE SHARE
// ---------------------------------------------------------------------------

// parseCreateShareStmt parses
//
//	CREATE [ OR REPLACE ] SHARE [ IF NOT EXISTS ] <name> [ COMMENT = '<string>' ]
//
// The CREATE keyword and optional OR REPLACE modifier have already been consumed;
// start is the Loc of the CREATE token, and cur is the SHARE keyword.
func (p *Parser) parseCreateShareStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume SHARE

	stmt := &ast.CreateShareStmt{
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

	// COMMENT is the only documented parameter, but capture any trailing
	// `KEY = value` option open-ended (mirrors the engine's parse-permissively
	// stance).
	opts, err := p.parseGroupOptions()
	if err != nil {
		return nil, err
	}
	stmt.Options = opts

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER SHARE
// ---------------------------------------------------------------------------

// parseAlterShareStmt parses ALTER SHARE. The ALTER keyword has already been
// consumed; cur is the SHARE keyword.
//
//	ALTER SHARE [ IF EXISTS ] <name> { ADD | REMOVE } ACCOUNTS = <a> [ , ... ] [ SHARE_RESTRICTIONS = ... ]
//	ALTER SHARE [ IF EXISTS ] <name> SET [ ACCOUNTS = <a> [ , ... ] ] [ COMMENT = '...' ]
//	ALTER SHARE [ IF EXISTS ] <name> SET TAG <tag> = '<value>' [ , ... ]
//	ALTER SHARE <name> UNSET TAG <tag> [ , ... ]
//	ALTER SHARE [ IF EXISTS ] <name> UNSET COMMENT
func (p *Parser) parseAlterShareStmt() (ast.Node, error) {
	altTok := p.advance() // consume SHARE
	stmt := &ast.AlterShareStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	p.parseOptionalIfExists(&stmt.IfExists)

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwADD, kwREMOVE:
		// { ADD | REMOVE } ACCOUNTS = <a> [ , ... ] [ SHARE_RESTRICTIONS = ... ]
		add := p.cur.Type == kwADD
		p.advance() // consume ADD / REMOVE
		if _, err := p.expect(kwACCOUNTS); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		accounts, err := p.parseObjectNameList()
		if err != nil {
			return nil, err
		}
		if add {
			stmt.Action = ast.AlterShareAdd
		} else {
			stmt.Action = ast.AlterShareRemove
		}
		stmt.Accounts = accounts
		if p.cur.Type == kwSHARE_RESTRICTIONS {
			p.advance() // consume SHARE_RESTRICTIONS
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			b, err := p.parseTrueFalse()
			if err != nil {
				return nil, err
			}
			stmt.ShareRestrictions = &b
		}

	case kwSET:
		if err := p.parseAlterShareSet(stmt); err != nil {
			return nil, err
		}

	case kwUNSET:
		if err := p.parseAlterShareUnset(stmt); err != nil {
			return nil, err
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterShareSet parses the SET branch of ALTER SHARE. cur is the SET keyword.
//
//	SET TAG <tag> = '<value>' [ , ... ]                  -> AlterShareSetTag
//	SET [ ACCOUNTS = <a> [ , ... ] ] [ COMMENT = '...' ] -> AlterShareSet
func (p *Parser) parseAlterShareSet(stmt *ast.AlterShareStmt) error {
	p.advance() // consume SET

	if p.cur.Type == kwTAG {
		tags, err := p.parseTagSetList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterShareSetTag
		stmt.Tags = tags
		return nil
	}

	stmt.Action = ast.AlterShareSet
	settable := false

	// Optional ACCOUNTS = <a> [ , ... ].
	if p.cur.Type == kwACCOUNTS {
		p.advance() // consume ACCOUNTS
		if _, err := p.expect('='); err != nil {
			return err
		}
		accounts, err := p.parseObjectNameList()
		if err != nil {
			return err
		}
		stmt.Accounts = accounts
		settable = true
	}

	// Optional COMMENT / other open-ended params.
	opts, err := p.parseGroupOptions()
	if err != nil {
		return err
	}
	if len(opts) > 0 {
		stmt.Options = opts
		settable = true
	}

	if !settable {
		// SET with nothing settable is a syntax error.
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterShareUnset parses the UNSET branch of ALTER SHARE. cur is the UNSET
// keyword.
//
//	UNSET TAG <tag> [ , ... ]   -> AlterShareUnsetTag
//	UNSET COMMENT [ , ... ]     -> AlterShareUnset
func (p *Parser) parseAlterShareUnset(stmt *ast.AlterShareStmt) error {
	p.advance() // consume UNSET

	if p.cur.Type == kwTAG {
		names, err := p.parseUnsetTagNameList()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterShareUnsetTag
		stmt.UnsetTags = names
		return nil
	}

	props, err := p.parseStageUnsetProps()
	if err != nil {
		return err
	}
	stmt.Action = ast.AlterShareUnset
	stmt.UnsetProps = props
	return nil
}

// ---------------------------------------------------------------------------
// Shared helpers — GroupOption (KEY = <value-list | literal>)
// ---------------------------------------------------------------------------

// parseGroupOptions parses a run of zero or more GroupOptions. Each is a
// `KEY = <value>` pair where the value is an unparenthesized comma list of
// word/dotted items or a string/number literal. The run continues while the
// current token begins another option, optionally consuming a comma separator
// between options (CREATE MANAGED ACCOUNT separates its options with commas; the
// group/share/account list-options are space-separated). It terminates at a
// statement boundary, at the [WITH] TAG clause, at the IGNORE EDITION CHECK
// flag, or at any token that does not begin an option.
func (p *Parser) parseGroupOptions() ([]*ast.GroupOption, error) {
	var opts []*ast.GroupOption
	for p.startsGroupOption() {
		opt, err := p.parseGroupOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
		if p.cur.Type == ',' {
			p.advance() // consume optional ',' separator (CREATE MANAGED ACCOUNT)
		}
	}
	return opts, nil
}

// startsGroupOption reports whether the current token can begin a GroupOption.
// An option name is an identifier or a keyword. The structural keywords that
// anchor a trailing clause (WITH / TAG for [WITH] TAG, IGNORE for IGNORE EDITION
// CHECK) are excluded so they are not swallowed as an option name.
func (p *Parser) startsGroupOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG, kwIGNORE:
		return false
	}
	if p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent {
		return true
	}
	if p.cur.Type == tokEOF || p.cur.Type == ';' {
		return false
	}
	return p.cur.Type >= 700
}

// startsGroupTags reports whether the current position begins a [WITH] TAG
// clause: the TAG keyword directly, or WITH immediately followed by TAG.
func (p *Parser) startsGroupTags() bool {
	if p.cur.Type == kwTAG {
		return true
	}
	return p.cur.Type == kwWITH && p.peekNext().Type == kwTAG
}

// parseGroupOption parses one `KEY = <value>` option. The value is either a
// string/number literal (REPLICATION_SCHEDULE = '10 MINUTE', ADMIN_PASSWORD =
// 'secret') or a run of word/dotted items. Whether the value is a comma list is
// driven by the option NAME, not by punctuation: only the closed set of
// replication/sharing list options (OBJECT_TYPES, ALLOWED_DATABASES,
// ALLOWED_EXTERNAL_VOLUMES, ALLOWED_SHARES, ALLOWED_INTEGRATION_TYPES,
// ALLOWED_ACCOUNTS) take an UNPARENTHESIZED comma list of elements (each itself a
// multi-word object-type name or a dotted org.account). Every other option is
// scalar (a single word/dotted value: ERROR_INTEGRATION = my_int,
// OPTIMIZED_REFRESH = TRUE, ADMIN_NAME = admin, EDITION = STANDARD, TYPE =
// READER, ...), so a trailing comma is an OPTION separator (CREATE MANAGED
// ACCOUNT comma-separates its params) — it is left for parseGroupOptions, NOT
// swallowed into this option's value. This name-driven rule is what
// disambiguates `OBJECT_TYPES = ROLES, DATABASES` (one option, two list
// elements) from `ADMIN_NAME = admin, ADMIN_PASSWORD = '...'` (two options) with
// only one token of lookahead. cur is the option name.
func (p *Parser) parseGroupOption() (*ast.GroupOption, error) {
	nameTok := p.advance() // consume the option name
	opt := &ast.GroupOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start, End: nameTok.Loc.End},
	}

	if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// A string/number literal value.
	switch p.cur.Type {
	case tokString:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	case tokInt:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	case tokFloat, tokReal:
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}
		opt.Loc.End = tok.Loc.End
		return opt, nil
	}

	// A run of word/dotted items: a comma list for the known list options, a
	// single element otherwise.
	values, end, err := p.parseGroupValueList(isGroupListOption(opt.Name))
	if err != nil {
		return nil, err
	}
	opt.Values = values
	opt.Loc.End = end
	return opt, nil
}

// isGroupListOption reports whether name is one of the replication/sharing list
// options whose value is an unparenthesized comma list of elements (so a comma
// within its value is a list separator, not an option separator). Every other
// option is scalar. The set is closed per the docs + legacy grammar.
func isGroupListOption(name string) bool {
	switch name {
	case "OBJECT_TYPES",
		"ALLOWED_DATABASES",
		"ALLOWED_EXTERNAL_VOLUMES",
		"ALLOWED_SHARES",
		"ALLOWED_INTEGRATION_TYPES",
		"ALLOWED_ACCOUNTS":
		return true
	}
	return false
}

// parseGroupValueList parses a value made of one or more word/dotted/'$' items,
// each captured verbatim and uppercased. Each element may be multi-word
// (`ACCOUNT PARAMETERS`, `RESOURCE MONITORS`, `SECURITY INTEGRATIONS`) or dotted
// (`org.account`). When list is true it parses a comma-separated list (the comma
// is the element separator); when list is false it parses exactly one element
// (any trailing comma is an option separator, left for the caller). Consumes at
// least one element. Returns the joined elements and the end offset.
//
// The element reader stops at a comma, at a token that begins the next
// `KEY = value` option (a word immediately followed by '='), or at a non-word
// token / clause anchor / statement boundary — which keeps a space-separated
// follow-on option (e.g. the ALLOWED_ACCOUNTS after an unparenthesized
// OBJECT_TYPES list) from being absorbed into the current value.
func (p *Parser) parseGroupValueList(list bool) ([]string, int, error) {
	if !p.isOptionWord(p.cur.Type) {
		return nil, 0, p.syntaxErrorAtCur()
	}
	var values []string
	end := p.cur.Loc.End
	for {
		elem, elemEnd, err := p.parseGroupValueElement()
		if err != nil {
			return nil, 0, err
		}
		values = append(values, elem)
		end = elemEnd
		if list && p.cur.Type == ',' {
			p.advance() // consume ',' list separator
			continue
		}
		break
	}
	return values, end, nil
}

// parseGroupValueElement parses one comma-list element: a run of one or more
// word tokens (joined with single spaces) plus any source-adjacent dotted / '$'
// continuations (joined verbatim with no space). The run stops at:
//   - a comma (the list separator), or
//   - a word that begins the NEXT option, detected as a word immediately
//     followed by '=' (e.g. the ALLOWED_ACCOUNTS in `OBJECT_TYPES = DATABASES
//     ALLOWED_ACCOUNTS = ...`), or
//   - a non-word token / the WITH / TAG / IGNORE clause anchors / EOF.
func (p *Parser) parseGroupValueElement() (string, int, error) {
	if !p.isOptionWord(p.cur.Type) {
		return "", 0, p.syntaxErrorAtCur()
	}
	var b strings.Builder
	first := p.advance()
	b.WriteString(strings.ToUpper(first.Str))
	end := first.Loc.End

	for {
		// Adjacency-joined dotted / '$' continuations (no intervening space):
		// dotted account names (org.account) and $-segments.
		if (p.cur.Type == '.' || p.cur.Type == '$') && p.cur.Loc.Start == end {
			sep := p.advance()
			b.WriteString(p.srcSlice(sep.Loc.Start, sep.Loc.End))
			end = sep.Loc.End
			if p.isOptionWord(p.cur.Type) && p.cur.Loc.Start == end {
				part := p.advance()
				b.WriteString(strings.ToUpper(part.Str))
				end = part.Loc.End
			}
			continue
		}
		// A further space-separated word continues this element ONLY if it is part
		// of a multi-word object-type name, i.e. it is not the start of the next
		// `KEY = value` option (a word followed by '=') and not a clause anchor.
		if p.continuesGroupElement() {
			part := p.advance()
			b.WriteByte(' ')
			b.WriteString(strings.ToUpper(part.Str))
			end = part.Loc.End
			continue
		}
		break
	}
	return b.String(), end, nil
}

// continuesGroupElement reports whether the current token extends the current
// multi-word list element (e.g. PARAMETERS in `ACCOUNT PARAMETERS`, MONITORS in
// `RESOURCE MONITORS`, INTEGRATIONS in `SECURITY INTEGRATIONS`) rather than
// beginning the next option or a trailing clause. A word that is immediately
// followed by '=' starts the next option; the WITH / TAG / IGNORE anchors and
// any non-word token end the element.
func (p *Parser) continuesGroupElement() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG, kwIGNORE:
		return false
	}
	if !p.isOptionWord(p.cur.Type) {
		return false
	}
	// A word followed by '=' is the next option name, not a continuation.
	return p.peekNext().Type != '='
}

// ---------------------------------------------------------------------------
// Shared helpers — name lists, ALLOWED_* targets, IGNORE EDITION CHECK, booleans
// ---------------------------------------------------------------------------

// parseObjectNameList parses a comma-separated list of object names (each a
// 1/2/3-part dotted ObjectName — used for database/share name lists and the
// dotted org.account lists). Consumes at least one.
func (p *Parser) parseObjectNameList() ([]*ast.ObjectName, error) {
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

// parseAllowedListKeyword consumes an ALLOWED_DATABASES / ALLOWED_SHARES /
// ALLOWED_ACCOUNTS keyword (the target of an ADD ... TO / REMOVE ... FROM
// action) and returns its uppercased text.
func (p *Parser) parseAllowedListKeyword() (string, error) {
	switch p.cur.Type {
	case kwALLOWED_DATABASES:
		p.advance()
		return "ALLOWED_DATABASES", nil
	case kwALLOWED_SHARES:
		p.advance()
		return "ALLOWED_SHARES", nil
	case kwALLOWED_ACCOUNTS:
		p.advance()
		return "ALLOWED_ACCOUNTS", nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// parseIgnoreEditionCheck consumes the bare IGNORE EDITION CHECK flag. cur is the
// IGNORE keyword.
func (p *Parser) parseIgnoreEditionCheck() error {
	p.advance() // consume IGNORE
	if _, err := p.expect(kwEDITION); err != nil {
		return err
	}
	if _, err := p.expect(kwCHECK); err != nil {
		return err
	}
	return nil
}

// parseTrueFalse parses a { TRUE | FALSE } boolean value, returning the bool. The
// lexer keywordizes TRUE / FALSE; they are also accepted as bare identifiers
// defensively (some are not reserved).
func (p *Parser) parseTrueFalse() (bool, error) {
	switch {
	case p.curIsWord("TRUE"):
		p.advance()
		return true, nil
	case p.curIsWord("FALSE"):
		p.advance()
		return false, nil
	default:
		return false, p.syntaxErrorAtCur()
	}
}
