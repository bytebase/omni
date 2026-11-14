package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// ALTER USER — MFA actions and per-user credentials
// ---------------------------------------------------------------------------
//
// The subcommands that work on one of a user's MFA methods, key pairs,
// programmatic access tokens, or workload identities. Each is documented on
// its own page:
//
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user (mfaActions)
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-key-pair
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-key-pair
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-key-pair
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-key-pair
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-programmatic-access-token
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-programmatic-access-token
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-programmatic-access-token
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-rotate-programmatic-access-token
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-add-workload-identity
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-modify-workload-identity
//	https://docs.snowflake.com/en/sql-reference/sql/alter-user-remove-workload-identity
//
// The trailing properties (PUBLIC_KEY, ROLE_RESTRICTION, DAYS_TO_EXPIRY, TYPE,
// ARN, OIDC_AUDIENCE_LIST, ...) are open-ended KEY = value options, the house
// pattern for property bags; the one each page requires first (PUBLIC_KEY for
// a key pair, TYPE for a workload identity) is checked.

// parseAlterUserCredential parses, after ALTER USER [ IF EXISTS ] [ <name> ],
//
//	ENROLL MFA
//	ADD MFA METHOD OTP [ COUNT = <n> ]
//	REMOVE MFA METHOD <method>
//	MODIFY MFA METHOD <method> SET COMMENT = '<string>'
//	ADD KEY PAIR <k> PUBLIC_KEY = '<key>' [ <options> ]
//	ROTATE KEY PAIR <k> PUBLIC_KEY = '<key>' [ <options> ]
//	ADD { PROGRAMMATIC ACCESS TOKEN | PAT } [ IF NOT EXISTS ] <t> [ <options> ]
//	ROTATE { PROGRAMMATIC ACCESS TOKEN | PAT } [ IF EXISTS ] <t> [ <options> ]
//	ADD WORKLOAD IDENTITY <w> TYPE = { AWS | AZURE | GCP | OIDC } [ <options> ]
//	REMOVE <credential> [ IF EXISTS ] <name>
//	MODIFY <credential> [ IF EXISTS ] <name>
//	  { RENAME TO <new_name> | SET <options> | UNSET <property> [ , ... ] }
//
// where IF [ NOT ] EXISTS belongs to the token forms only. cur is ADD, REMOVE,
// MODIFY, ROTATE, or ENROLL.
func (p *Parser) parseAlterUserCredential(stmt *ast.AlterUserStmt) error {
	if p.curIsUnquotedWord("ENROLL") {
		p.advance() // consume ENROLL
		if !p.curIsUnquotedWord("MFA") {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume MFA
		stmt.Action = ast.AlterUserEnrollMFA
		return nil
	}

	verb := p.advance() // consume ADD / REMOVE / MODIFY / ROTATE
	rotate := verb.Type == tokIdent
	kind, err := p.parseUserCredentialKind(rotate)
	if err != nil {
		return err
	}
	stmt.Credential = kind

	if kind == ast.UserCredentialPAT {
		if verb.Type == kwADD {
			if err := p.parseIfNotExists(&stmt.CredentialIfNotExists); err != nil {
				return err
			}
		} else if err := p.parseIfExists(&stmt.CredentialIfExists); err != nil {
			return err
		}
	}

	// ADD MFA METHOD takes the method type, OTP, where the others take a name.
	if verb.Type == kwADD && kind == ast.UserCredentialMFAMethod && !p.curIsUnquotedWord("OTP") {
		return p.syntaxErrorAtCur()
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.CredentialName = name

	switch {
	case verb.Type == kwREMOVE:
		stmt.Action = ast.AlterUserRemoveCredential
		return nil
	case verb.Type == kwMODIFY:
		return p.parseModifyUserCredential(stmt)
	case rotate:
		stmt.Action = ast.AlterUserRotateCredential
	default:
		stmt.Action = ast.AlterUserAddCredential
	}

	// ADD MFA METHOD OTP takes nothing but an optional COUNT = <n>.
	if kind == ast.UserCredentialMFAMethod {
		if p.cur.Type == kwCOUNT && p.peekNext().Type == '=' {
			opt, err := p.parseCopyOption()
			if err != nil {
				return err
			}
			stmt.Options = []*ast.CopyOption{opt}
		}
		return nil
	}

	// ADD / ROTATE: the property the page requires first, then the rest.
	lead := ""
	switch kind {
	case ast.UserCredentialKeyPair:
		lead = "PUBLIC_KEY"
	case ast.UserCredentialWorkloadIdentity:
		lead = "TYPE"
	}
	if lead != "" && (!p.curIsWord(lead) || p.peekNext().Type != '=') {
		return p.syntaxErrorAtCur()
	}
	opts, err := p.parseUserOptions()
	if err != nil {
		return err
	}
	stmt.Options = opts
	return nil
}

// parseUserOptions parses a run of `KEY = value` options, the property lists
// of ALTER USER SET and of the credential subcommands. Unlike parseCopyOptions
// it takes no bare word as an option: every documented user property has a
// value, so the run stops at a word not followed by '=' and that word is left
// to fail as a trailing token.
func (p *Parser) parseUserOptions() ([]*ast.CopyOption, error) {
	var opts []*ast.CopyOption
	for p.startsCopyOption() && p.peekNext().Type == '=' {
		opt, err := p.parseCopyOption()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}
	return opts, nil
}

// parseUserCredentialKind consumes the credential words after ADD / REMOVE /
// MODIFY / ROTATE: KEY PAIR, PROGRAMMATIC ACCESS TOKEN, PAT, WORKLOAD
// IDENTITY, or MFA METHOD. ROTATE takes only KEY PAIR and the token.
func (p *Parser) parseUserCredentialKind(rotate bool) (ast.UserCredentialKind, error) {
	switch {
	case p.cur.Type == kwKEY:
		p.advance() // consume KEY
		if !p.curIsUnquotedWord("PAIR") {
			return 0, p.syntaxErrorAtCur()
		}
		p.advance() // consume PAIR
		return ast.UserCredentialKeyPair, nil
	case p.curIsUnquotedWord("PAT"):
		p.advance() // consume PAT
		return ast.UserCredentialPAT, nil
	case p.curIsUnquotedWord("PROGRAMMATIC"):
		p.advance() // consume PROGRAMMATIC
		if _, err := p.expect(kwACCESS); err != nil {
			return 0, err
		}
		if !p.curIsUnquotedWord("TOKEN") {
			return 0, p.syntaxErrorAtCur()
		}
		p.advance() // consume TOKEN
		return ast.UserCredentialPAT, nil
	case !rotate && p.curIsUnquotedWord("WORKLOAD"):
		p.advance() // consume WORKLOAD
		if _, err := p.expect(kwIDENTITY); err != nil {
			return 0, err
		}
		return ast.UserCredentialWorkloadIdentity, nil
	case !rotate && p.curIsUnquotedWord("MFA"):
		p.advance() // consume MFA
		if !p.curIsUnquotedWord("METHOD") {
			return 0, p.syntaxErrorAtCur()
		}
		p.advance() // consume METHOD
		return ast.UserCredentialMFAMethod, nil
	}
	return 0, p.syntaxErrorAtCur()
}

// parseModifyUserCredential parses the tail of MODIFY <credential> <name>:
// RENAME TO <new_name>, SET <options>, or UNSET <property> [ , ... ]. An MFA
// method takes only SET COMMENT = '<string>'.
func (p *Parser) parseModifyUserCredential(stmt *ast.AlterUserStmt) error {
	mfa := stmt.Credential == ast.UserCredentialMFAMethod
	switch {
	case p.cur.Type == kwRENAME && !mfa:
		newName, err := p.parseRenameTo()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterUserModifyCredentialRename
		stmt.NewName = newName
	case p.cur.Type == kwSET:
		p.advance() // consume SET
		if mfa {
			if !p.curIsWord("COMMENT") || p.peekNext().Type != '=' {
				return p.syntaxErrorAtCur()
			}
			opt, err := p.parseCopyOption()
			if err != nil {
				return err
			}
			stmt.Options = []*ast.CopyOption{opt}
		} else {
			opts, err := p.parseUserOptions()
			if err != nil {
				return err
			}
			if len(opts) == 0 {
				return p.syntaxErrorAtCur()
			}
			stmt.Options = opts
		}
		stmt.Action = ast.AlterUserModifyCredentialSet
	case p.cur.Type == kwUNSET && !mfa:
		p.advance() // consume UNSET
		props, err := p.parseUnsetPropertyNames()
		if err != nil {
			return err
		}
		stmt.Action = ast.AlterUserModifyCredentialUnset
		stmt.UnsetProps = props
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}
