package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Access-control DDL — CREATE / ALTER ROLE / USER / { MASKING | ROW ACCESS |
// SESSION | PASSWORD | NETWORK | AUTHENTICATION } POLICY (T4.6)
// ---------------------------------------------------------------------------
//
// Roles, users, and the six policy objects reuse the open-ended `KEY = <value>`
// option-bag machinery already merged for STAGE (T4.1) / FILE FORMAT (T4.2) /
// COPY (T5.2) / pipeline (T4.3) (parseCopyOption / startsCopyOption): every
// option that CREATE/ALTER USER and the SESSION / PASSWORD / NETWORK /
// AUTHENTICATION policies carry is a version-growing vocabulary the catalog /
// semantic layer — not the parser — validates, so each is captured as a
// CopyOption rather than enumerated against the already-stale legacy ANTLR rules
// (whose create_network_policy lists only ALLOWED_IP_LIST / BLOCKED_IP_LIST, not
// the docs' ALLOWED_NETWORK_RULE_LIST, and which has no AUTHENTICATION POLICY
// rule at all). The two structurally distinct objects — MASKING POLICY and ROW
// ACCESS POLICY — carry a typed argument list, a RETURNS type, and an expression
// body after `->`; those anchor the grammar (Args / Returns / Body), and the
// body reuses the expression parser so it participates in walks / query-span /
// analysis. Trailing `KEY = value` options after a MASKING / ROW ACCESS body
// (COMMENT / EXEMPT_OTHER_POLICIES) are captured open-ended in Options.
//
// Oracle (no live Snowflake instance): triangulation of the legacy
// SnowflakeParser.g4 rules (create_role / create_user / create_masking_policy /
// create_row_access_policy / create_{session,password,network}_policy +
// alter_*) against the official docs (docs win on conflict) and the
// testdata/official/{create-role,create-user,create-masking-policy,
// create-row-access-policy,create-network-policy} corpus.

// ---------------------------------------------------------------------------
// CREATE ROLE
// ---------------------------------------------------------------------------

// parseCreateRoleStmt parses
//
//	CREATE [ OR REPLACE ] [ DATABASE ] ROLE [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//
// The CREATE keyword and the optional OR REPLACE modifier have already been
// consumed by parseCreateStmt; start is the Loc of the CREATE token. database is
// true when the DATABASE keyword preceded ROLE (CREATE DATABASE ROLE), whose
// <name> may be db-qualified. On entry cur is the ROLE keyword.
//
// The docs list COMMENT before [WITH] TAG; the legacy ANTLR create_role lists
// `with_tags? comment_clause?` (TAG before COMMENT). Accepting either order (a)
// follows the docs, (b) regresses neither, and (c) avoids silently dropping a
// COMMENT that trails the TAG clause — matching the merged STAGE (T4.1) handling.
func (p *Parser) parseCreateRoleStmt(start ast.Loc, orReplace, orAlter, database bool) (ast.Node, error) {
	p.advance() // consume ROLE

	stmt := &ast.CreateRoleStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Database:  database,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CREATE [ OR REPLACE ] DATABASE ROLE [ IF NOT EXISTS ] <name> CLONE <source_role>
	// (https://docs.snowflake.com/en/sql-reference/sql/create-clone).
	if p.curIs(kwCLONE) {
		if !database || orAlter || (orReplace && stmt.IfNotExists) || !name.Database.IsEmpty() {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume CLONE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if !src.Database.IsEmpty() {
			return nil, &ParseError{
				Position: src.Loc.Start, End: src.Loc.End,
				Message: "syntax error: a database role's name has at most 2 parts",
			}
		}
		stmt.Clone = src
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Trailing COMMENT / [WITH] TAG, in either order.
	for {
		if p.curIs(kwCOMMENT) {
			if err := p.parseEqComment(&stmt.Comment); err != nil {
				return nil, err
			}
			continue
		}
		if p.startsWithTags() {
			tags, err := p.parseWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE USER
// ---------------------------------------------------------------------------

// parseCreateUserStmt parses
//
//	CREATE [ OR REPLACE ] USER [ IF NOT EXISTS ] <name>
//	  [ objectProperties ] [ objectParams ] [ sessionParams ]
//	  [ [ WITH ] TAG ( <tag> = '<value>' [ , ... ] ) ]
//
// Every property/parameter is an open-ended `KEY = value` pair (CopyOption).
// The CREATE keyword and OR REPLACE have already been consumed; cur is USER.
func (p *Parser) parseCreateUserStmt(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	p.advance() // consume USER

	stmt := &ast.CreateUserStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Open-ended option bag interleaved with the trailing [WITH] TAG clause.
	// WITH/TAG are excluded as option starters (startsUserOption) so the tag
	// clause is not swallowed as an option named WITH or TAG.
	for {
		if p.startsWithTags() {
			tags, err := p.parseWithTags()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		}
		if p.startsUserOption() {
			opt, err := p.parseCopyOption()
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
// CREATE { MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK | AUTHENTICATION }
// POLICY
// ---------------------------------------------------------------------------

// parseCreatePolicyStmt parses the CREATE form of one of the six policy objects.
// The CREATE keyword and OR REPLACE have already been consumed by
// parseCreateStmt; the policy-kind keyword(s) (MASKING / ROW ACCESS / SESSION /
// PASSWORD / NETWORK / AUTHENTICATION) and the trailing POLICY keyword have ALSO
// already been consumed by the dispatcher. start is the Loc of the CREATE token
// and kind identifies which policy.
//
// MASKING / ROW ACCESS:
//
//	... [ IF NOT EXISTS ] <name> AS ( <arg> <type> [ , ... ] )
//	      RETURNS <type> -> <body> [ <trailing options> ]
//
// SESSION / PASSWORD / NETWORK / AUTHENTICATION:
//
//	... [ IF NOT EXISTS ] <name> [ <option> ... ]
func (p *Parser) parseCreatePolicyStmt(start ast.Loc, orReplace, orAlter bool, kind ast.PolicyKind) (ast.Node, error) {
	stmt := &ast.CreatePolicyStmt{
		Kind:      kind,
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExists(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch kind {
	case ast.PolicyMasking, ast.PolicyRowAccess:
		// AS ( <arg> <type> [ , ... ] ) RETURNS <type> -> <body>
		if _, err := p.expect(kwAS); err != nil {
			return nil, err
		}
		args, err := p.parsePolicyArgs()
		if err != nil {
			return nil, err
		}
		stmt.Args = args

		if _, err := p.expect(kwRETURNS); err != nil {
			return nil, err
		}
		ret, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		stmt.Returns = ret

		if _, err := p.expect(tokArrow); err != nil {
			return nil, err
		}
		body, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.Body = body

		// Trailing open-ended options (COMMENT / EXEMPT_OTHER_POLICIES). The
		// expression parser stops at the first non-expression token, which is the
		// first trailing option name (or the statement boundary).
		opts, err := p.parseCopyOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts

	default:
		// SESSION / PASSWORD / NETWORK / AUTHENTICATION: open-ended option bag.
		opts, err := p.parseCopyOptions()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parsePolicyArgs parses the parenthesized typed signature of a MASKING / ROW
// ACCESS policy: ( <arg_name> <arg_type> [ , <arg_name> <arg_type> ... ] ).
// At least one argument is required (per the docs and the legacy grammar).
func (p *Parser) parsePolicyArgs() ([]*ast.PolicyArg, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var args []*ast.PolicyArg
	for {
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		args = append(args, &ast.PolicyArg{
			Name:     name,
			DataType: dt,
			Loc:      ast.Loc{Start: name.Loc.Start, End: dt.Loc.End},
		})

		if p.curIs(',') {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return args, nil
}

// ---------------------------------------------------------------------------
// ALTER ROLE
// ---------------------------------------------------------------------------

// parseAlterRoleStmt parses
//
//	ALTER [ DATABASE ] ROLE [ IF EXISTS ] <name>
//	  { RENAME TO <new_name>
//	  | SET COMMENT = '<string>'
//	  | UNSET COMMENT
//	  | SET TAG <tag> = '<value>' [ , ... ]
//	  | UNSET TAG <tag> [ , ... ] }
//
// The ALTER keyword has already been consumed by parseAlterStmt; database is
// true when DATABASE preceded ROLE. On entry cur is the ROLE keyword.
func (p *Parser) parseAlterRoleStmt(database bool) (ast.Node, error) {
	altTok := p.advance() // consume ROLE
	stmt := &ast.AlterRoleStmt{Database: database, Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExists(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.curIs(kwRENAME):
		newName, err := p.parseRenameTo()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterRoleRename
		stmt.NewName = newName

	case p.curIs(kwSET):
		p.advance() // consume SET
		if p.curIs(kwTAG) {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterRoleSetTag
			stmt.Tags = tags
		} else {
			// SET COMMENT = '...'
			if !p.curIs(kwCOMMENT) {
				return nil, p.syntaxErrorAtCur()
			}
			if err := p.parseEqComment(&stmt.Comment); err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterRoleSetComment
		}

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterRoleUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET COMMENT (the only unsettable role property).
			if !p.curIs(kwCOMMENT) {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // consume COMMENT
			stmt.Action = ast.AlterRoleUnset
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER USER
// ---------------------------------------------------------------------------

// parseAlterUserStmt parses
//
//	ALTER USER [ IF EXISTS ] [ <name> ]
//	  { RENAME TO <new_name>
//	  | RESET PASSWORD
//	  | ABORT ALL QUERIES
//	  | SET <options>
//	  | UNSET <property> [ , ... ]
//	  | SET TAG <tag> = '<value>' [ , ... ]
//	  | UNSET TAG <tag> [ , ... ]
//	  | SET { AUTHENTICATION | PASSWORD | SESSION | INGRESS } POLICY <p> [ FORCE ]
//	  | UNSET { AUTHENTICATION | PASSWORD | SESSION | INGRESS } POLICY
//	  | ADD DELEGATED AUTHORIZATION OF ROLE <r> TO SECURITY INTEGRATION <i>
//	  | REMOVE DELEGATED { AUTHORIZATION OF ROLE <r> | AUTHORIZATIONS [ OF ANY ROLE ] }
//	      FROM SECURITY INTEGRATION <i>
//	  | <mfaActions>
//	  | { ADD | MODIFY | REMOVE | ROTATE } <credential> ... }
//
// The <name> is optional: "If the identifier is omitted, the statement modifies
// the active (i.e. logged in) user"
// (https://docs.snowflake.com/en/sql-reference/sql/alter-user). The credential
// and MFA subcommands are parsed in user_credential.go; INGRESS POLICY is
// https://docs.snowflake.com/en/sql-reference/sql/alter-user-set-ingress-policy.
//
// The ALTER keyword has already been consumed; cur is the USER keyword.
func (p *Parser) parseAlterUserStmt() (ast.Node, error) {
	altTok := p.advance() // consume USER
	stmt := &ast.AlterUserStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExists(&stmt.IfExists); err != nil {
		return nil, err
	}

	if !p.omitsAlterUserName() {
		name, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
	}

	switch {
	case p.curIs(kwRENAME):
		newName, err := p.parseRenameTo()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterUserRename
		stmt.NewName = newName

	case p.curIs(kwRESET):
		// RESET PASSWORD
		p.advance() // consume RESET
		if _, err := p.expect(kwPASSWORD); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterUserResetPassword

	case p.curIs(kwABORT):
		// ABORT ALL QUERIES
		p.advance() // consume ABORT
		if _, err := p.expect(kwALL); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwQUERIES); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterUserAbortQueries

	case p.curIs(kwADD) && p.nextIs(kwDELEGATED):
		// ADD DELEGATED AUTHORIZATION OF ROLE <r> TO SECURITY INTEGRATION <i>
		rawStart := p.cur.Loc.Start
		if err := p.parseAddDelegatedAuthorization(); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterUserAddDelegated
		stmt.Raw = p.srcSlice(rawStart, p.prev.Loc.End)

	case p.curIs(kwREMOVE) && p.nextIs(kwDELEGATED):
		// REMOVE DELEGATED { AUTHORIZATION OF ROLE <r> | AUTHORIZATIONS
		//   [ OF ANY ROLE ] } FROM SECURITY INTEGRATION <i>
		rawStart := p.cur.Loc.Start
		if err := p.parseRemoveDelegatedAuthorization(); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterUserRemoveDelegated
		stmt.Raw = p.srcSlice(rawStart, p.prev.Loc.End)

	case p.curIs(kwADD), p.curIs(kwREMOVE), p.curIs(kwMODIFY),
		p.curIsUnquotedWord("ROTATE"), p.curIsUnquotedWord("ENROLL"):
		if err := p.parseAlterUserCredential(stmt); err != nil {
			return nil, err
		}

	case p.curIs(kwSET):
		p.advance() // consume SET
		switch {
		case p.curIs(kwTAG):
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterUserSetTag
			stmt.Tags = tags
		case p.startsUserPolicy(stmt.Name != nil):
			stmt.PolicyKind = p.parseUserPolicyKind()
			name, err := p.parseObjectName()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterUserSetPolicy
			stmt.PolicyName = name
			if p.curIs(kwFORCE) {
				p.advance() // consume FORCE
				stmt.Force = true
			}
		default:
			// SET <options> — open-ended KEY = value params.
			opts, err := p.parseUserOptions()
			if err != nil {
				return nil, err
			}
			if len(opts) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterUserSet
			stmt.Options = opts
		}

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		switch {
		case p.curIs(kwTAG):
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterUserUnsetTag
			stmt.UnsetTags = names
		case p.startsUserPolicy(stmt.Name != nil):
			stmt.PolicyKind = p.parseUserPolicyKind()
			stmt.Action = ast.AlterUserUnsetPolicy
		default:
			// UNSET <property> [ , ... ] — open-ended name list.
			props, err := p.parseUnsetPropertyNames()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterUserUnset
			stmt.UnsetProps = props
		}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// omitsAlterUserName reports whether ALTER USER leaves out the user name, so
// that cur already begins the action. That holds when cur is an action word
// and the token after it is neither another action word nor a '.': in
// `ALTER USER SET PASSWORD = 's'` SET begins the action, while in
// `ALTER USER reset RESET PASSWORD` the first RESET is the user's name.
func (p *Parser) omitsAlterUserName() bool {
	if !isAlterUserActionWord(p.cur) {
		return false
	}
	next := p.peekNext()
	return !p.nextIs('.') && !isAlterUserActionWord(next)
}

// isAlterUserActionWord reports whether tok is a word that can begin an ALTER
// USER action: RENAME, RESET, ABORT, ADD, REMOVE, SET, UNSET, MODIFY, or the
// unreserved ROTATE and ENROLL (which lex as identifiers). A double-quoted
// "ROTATE" is a name, never an action.
func isAlterUserActionWord(tok Token) bool {
	switch tok.Type {
	case kwRENAME, kwRESET, kwABORT, kwADD, kwREMOVE, kwSET, kwUNSET, kwMODIFY:
		return true
	case tokIdent:
		u := strings.ToUpper(tok.Str)
		return u == "ROTATE" || u == "ENROLL"
	}
	return false
}

// curIsUnquotedWord reports whether cur is the unquoted identifier upper (a
// word the lexer does not keywordize), so a double-quoted name never matches.
func (p *Parser) curIsUnquotedWord(upper string) bool {
	if p.curIs(tokIdent) && strings.ToUpper(p.cur.Str) == upper {
		return true
	}
	p.wantWords(upper)
	return false
}

// startsUserPolicy reports whether cur begins the policy-kind words of ALTER
// USER SET / UNSET { AUTHENTICATION | PASSWORD | SESSION | INGRESS } POLICY: a
// kind word directly followed by POLICY. The POLICY check keeps
// `SET PASSWORD = '...'` an ordinary property. INGRESS counts only when the
// statement names its user (named): its page writes `ALTER USER <name>`
// (https://docs.snowflake.com/en/sql-reference/sql/alter-user-set-ingress-policy).
func (p *Parser) startsUserPolicy(named bool) bool {
	if !p.nextIs(kwPOLICY) {
		return false
	}
	switch {
	case p.curIs(kwPASSWORD), p.curIs(kwSESSION):
		return true
	case p.curIsUnquotedWord("AUTHENTICATION"):
		return true
	case p.curIsUnquotedWord("INGRESS"):
		return named
	}
	return false
}

// parseUserPolicyKind consumes the kind word and POLICY (the caller confirmed
// startsUserPolicy) and returns the uppercased "<KIND> POLICY".
func (p *Parser) parseUserPolicyKind() string {
	kindTok := p.advance() // consume the policy-kind word
	p.advance()            // consume POLICY
	return strings.ToUpper(kindTok.Str) + " POLICY"
}

// parseAddDelegatedAuthorization parses the body following ADD:
//
//	ADD DELEGATED AUTHORIZATION OF ROLE <r> TO SECURITY INTEGRATION <i>
//
// On entry cur is the ADD keyword. The identifiers are consumed but not modeled
// (the caller captures the verbatim tail via Raw).
func (p *Parser) parseAddDelegatedAuthorization() error {
	p.advance() // consume ADD
	if _, err := p.expect(kwDELEGATED); err != nil {
		return err
	}
	if _, err := p.expect(kwAUTHORIZATION); err != nil {
		return err
	}
	if _, err := p.expect(kwOF); err != nil {
		return err
	}
	if _, err := p.expect(kwROLE); err != nil {
		return err
	}
	if _, err := p.parseObjectName(); err != nil {
		return err
	}
	if _, err := p.expect(kwTO); err != nil {
		return err
	}
	if _, err := p.expect(kwSECURITY); err != nil {
		return err
	}
	if _, err := p.expect(kwINTEGRATION); err != nil {
		return err
	}
	if _, err := p.parseObjectName(); err != nil {
		return err
	}
	return nil
}

// parseRemoveDelegatedAuthorization parses the body following REMOVE:
//
//	REMOVE DELEGATED { AUTHORIZATION OF ROLE <r> | AUTHORIZATIONS [ OF ANY ROLE ] }
//	  FROM SECURITY INTEGRATION <i>
//
// On entry cur is the REMOVE keyword.
func (p *Parser) parseRemoveDelegatedAuthorization() error {
	p.advance() // consume REMOVE
	if _, err := p.expect(kwDELEGATED); err != nil {
		return err
	}
	switch {
	case p.curIs(kwAUTHORIZATION):
		p.advance() // consume AUTHORIZATION
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		if _, err := p.expect(kwROLE); err != nil {
			return err
		}
		if _, err := p.parseObjectName(); err != nil {
			return err
		}
	case p.curIs(kwAUTHORIZATIONS):
		p.advance() // consume AUTHORIZATIONS
		// AUTHORIZATIONS OF ANY ROLE is the same removal written out
		// (https://docs.snowflake.com/en/sql-reference/sql/alter-user).
		if p.curIs(kwOF) {
			p.advance() // consume OF
			if _, err := p.expect(kwANY); err != nil {
				return err
			}
			if _, err := p.expect(kwROLE); err != nil {
				return err
			}
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if _, err := p.expect(kwSECURITY); err != nil {
		return err
	}
	if _, err := p.expect(kwINTEGRATION); err != nil {
		return err
	}
	if _, err := p.parseObjectName(); err != nil {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------
// ALTER { MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK | AUTHENTICATION }
// POLICY
// ---------------------------------------------------------------------------

// parseAlterPolicyStmt parses the ALTER form of one of the six policy objects.
// The ALTER keyword and the policy-kind keyword(s) + the POLICY keyword have
// already been consumed by the dispatcher; startLoc anchors Loc.Start and kind
// identifies which policy.
//
//	... [ IF EXISTS ] <name>
//	  { RENAME TO <new_name>
//	  | SET BODY -> <expr>                     -- MASKING / ROW ACCESS only
//	  | SET { COMMENT = '...' | <options> }
//	  | UNSET { COMMENT | <property> } [ , ... ]
//	  | SET TAG <tag> = '<value>' [ , ... ]
//	  | UNSET TAG <tag> [ , ... ]
//	  | { ADD | REMOVE } { ALLOWED_NETWORK_RULE_LIST | BLOCKED_NETWORK_RULE_LIST }
//	      = ( '<network_rule>' [ , ... ] ) }      -- NETWORK only, without IF EXISTS
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-network-policy
func (p *Parser) parseAlterPolicyStmt(startLoc ast.Loc, kind ast.PolicyKind) (ast.Node, error) {
	stmt := &ast.AlterPolicyStmt{Kind: kind, Loc: ast.Loc{Start: startLoc.Start}}

	if err := p.parseIfExists(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch {
	case p.curIs(kwRENAME):
		newName, err := p.parseRenameTo()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterPolicyRename
		stmt.NewName = newName

	case p.curIs(kwSET):
		p.advance() // consume SET
		switch {
		case p.curIs(kwBODY):
			// SET BODY -> <expr> (MASKING / ROW ACCESS).
			p.advance() // consume BODY
			if _, err := p.expect(tokArrow); err != nil {
				return nil, err
			}
			body, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPolicySetBody
			stmt.Body = body
		case p.curIs(kwTAG):
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPolicySetTag
			stmt.Tags = tags
		default:
			// SET { COMMENT = '...' | <option> ... } — open-ended option bag.
			opts, err := p.parseCopyOptions()
			if err != nil {
				return nil, err
			}
			if len(opts) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Action = ast.AlterPolicySet
			stmt.Options = opts
		}

	case p.curIs(kwUNSET):
		p.advance() // consume UNSET
		if p.curIs(kwTAG) {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPolicyUnsetTag
			stmt.UnsetTags = names
		} else {
			// UNSET { COMMENT | <param_name> } [ , ... ] — open-ended name list.
			props, err := p.parseUnsetPropertyNames()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterPolicyUnset
			stmt.UnsetProps = props
		}

	case p.curIs(kwADD), p.curIs(kwREMOVE):
		if kind != ast.PolicyNetwork || stmt.IfExists {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterPolicyAdd
		if p.advance().Type == kwREMOVE {
			stmt.Action = ast.AlterPolicyRemove
		}
		if !(p.atWords("ALLOWED_NETWORK_RULE_LIST") || p.atWords("BLOCKED_NETWORK_RULE_LIST")) || !p.nextIs('=') {
			return nil, p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return nil, err
		}
		if !propStrings.matches(opt) {
			return nil, &ParseError{
				Position: opt.Loc.Start, End: opt.Loc.End,
				Message: "syntax error: " + opt.Name + " takes ( '<network_rule>' [ , ... ] )",
			}
		}
		stmt.Options = []*ast.CopyOption{opt}

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared helpers
// ---------------------------------------------------------------------------

// startsPolicyKeyword reports whether the current token begins a policy-kind
// prefix ({ MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK |
// AUTHENTICATION } POLICY) such that consumePolicyKeywords would succeed. It is
// used by the ALTER dispatcher, whose switch must first decide whether ALTER is
// followed by a policy object. AUTHENTICATION is matched as a non-reserved
// identifier; the others are keywords. The trailing POLICY (and, for ROW, the
// intervening ACCESS) is confirmed so non-policy statements (e.g. ALTER NETWORK
// RULE, ALTER SESSION) are not misrouted.
func (p *Parser) startsPolicyKeyword() bool {
	switch {
	case p.curIs(kwMASKING), p.curIs(kwSESSION), p.curIs(kwPASSWORD), p.curIs(kwNETWORK):
		return p.nextIs(kwPOLICY)
	case p.curIs(kwROW):
		// ROW ACCESS POLICY — only the ROW + ACCESS prefix is checkable with one
		// token of lookahead; consumePolicyKeywords validates the POLICY keyword.
		return p.nextIs(kwACCESS)
	}
	return p.curIsWord("AUTHENTICATION") && p.nextIs(kwPOLICY)
}

// consumePolicyKeywords consumes the policy-kind keyword(s) and the trailing
// POLICY keyword, returning the identified PolicyKind. On entry cur is the first
// policy-kind word. Returns a syntax error if the token sequence is not a valid
// policy prefix.
func (p *Parser) consumePolicyKeywords() (ast.PolicyKind, error) {
	switch {
	case p.curIs(kwMASKING):
		p.advance() // consume MASKING
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicyMasking, nil
	case p.curIs(kwROW):
		p.advance() // consume ROW
		if _, err := p.expect(kwACCESS); err != nil {
			return 0, err
		}
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicyRowAccess, nil
	case p.curIs(kwSESSION):
		p.advance() // consume SESSION
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicySession, nil
	case p.curIs(kwPASSWORD):
		p.advance() // consume PASSWORD
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicyPassword, nil
	case p.curIs(kwNETWORK):
		p.advance() // consume NETWORK
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicyNetwork, nil
	}
	// AUTHENTICATION (non-reserved identifier) POLICY.
	if p.curIsWord("AUTHENTICATION") {
		p.advance() // consume AUTHENTICATION
		if _, err := p.expect(kwPOLICY); err != nil {
			return 0, err
		}
		return ast.PolicyAuthentication, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// startsUserOption reports whether the current token can begin a CREATE USER
// option. It is the COPY option predicate minus WITH and TAG, which anchor the
// trailing [WITH] TAG clause and must not be swallowed as an option name.
func (p *Parser) startsUserOption() bool {
	switch {
	case p.curIs(kwWITH), p.curIs(kwTAG):
		return false
	}
	return p.startsCopyOption()
}

// startsWithTags reports whether the current position begins a [WITH] TAG
// clause: either the TAG keyword directly, or WITH immediately followed by TAG.
func (p *Parser) startsWithTags() bool {
	if p.curIs(kwTAG) {
		return true
	}
	return p.curIs(kwWITH) && p.nextIs(kwTAG)
}

// parseWithTags parses a [WITH] TAG ( <tag> = '<value>' [ , ... ] ) clause and
// returns the tag assignments. The caller must have confirmed startsWithTags.
func (p *Parser) parseWithTags() ([]*ast.TagAssignment, error) {
	if p.curIs(kwWITH) {
		p.advance() // consume WITH (peekNext == TAG guaranteed by startsWithTags)
	}
	return p.parseTagAssignments()
}

// parseUnsetPropertyNames parses a comma-separated UNSET property-name list for
// USER / POLICY objects, where each name is a single open-ended name word
// (identifier or keyword) — the settable property/parameter vocabulary is large
// and version-growing (DEFAULT_ROLE / DISPLAY_NAME / DEFAULT_WAREHOUSE /
// SESSION_IDLE_TIMEOUT_MINS / PASSWORD_MIN_LENGTH / COMMENT / ...), much of it
// lexed as keywords, so it is NOT enumerated. Returns the uppercased names;
// consumes at least one.
func (p *Parser) parseUnsetPropertyNames() ([]string, error) {
	var names []string
	for {
		if !p.isOptionWord(p.cur.Type) {
			return nil, p.syntaxErrorAtCur()
		}
		names = append(names, strings.ToUpper(p.advance().Str))
		if p.curIs(',') {
			p.advance() // consume ','
			continue
		}
		break
	}
	return names, nil
}

// parseIfNotExists consumes an optional IF NOT EXISTS clause, setting *flag.
func (p *Parser) parseIfNotExists(flag *bool) error {
	if p.curIs(kwIF) && p.nextIs(kwNOT) {
		p.advance() // consume IF
		p.advance() // consume NOT
		if _, err := p.expect(kwEXISTS); err != nil {
			return err
		}
		*flag = true
	}
	return nil
}

// parseIfExists consumes an optional IF EXISTS clause, setting *flag.
func (p *Parser) parseIfExists(flag *bool) error {
	if p.curIs(kwIF) && p.nextIs(kwEXISTS) {
		p.advance() // consume IF
		p.advance() // consume EXISTS
		*flag = true
	}
	return nil
}

// parseRenameTo consumes RENAME TO <new_name> and returns the new name. On entry
// cur is the RENAME keyword.
func (p *Parser) parseRenameTo() (*ast.ObjectName, error) {
	p.advance() // consume RENAME
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	return p.parseObjectName()
}

// parseEqComment parses a `COMMENT [=] '<string>'` clause, storing the string
// into *dst. On entry cur is the COMMENT keyword. The '=' is required by the
// docs (COMMENT = '...') but tolerated-optional defensively, matching other
// merged DDL.
func (p *Parser) parseEqComment(dst **string) error {
	p.advance() // consume COMMENT
	if p.curIs('=') {
		p.advance() // consume '='
	}
	tok, err := p.expect(tokString)
	if err != nil {
		return err
	}
	s := tok.Str
	*dst = &s
	return nil
}
