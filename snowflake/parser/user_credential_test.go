package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// TestOmitsAlterUserName pins the lookahead that decides whether the text after
// ALTER USER [ IF EXISTS ] is the user's name or already the action.
func TestOmitsAlterUserName(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"SET PASSWORD = 's'", true},
		{"UNSET COMMENT", true},
		{"RESET PASSWORD", true},
		{"ROTATE PAT t", true},
		{"ENROLL MFA", true},
		{"ADD PAT t", true},
		{"u SET PASSWORD = 's'", false},
		{"reset RESET PASSWORD", false}, // a user named RESET
		{"rotate ROTATE PAT t", false},  // a user named ROTATE
		{`"ROTATE" PAT t`, false},       // a quoted name is never an action
		{"modify.x SET COMMENT = 's'", false},
		{"IF ADD PAT t", false}, // a user named IF
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(c.input), input: c.input}
			p.advance()
			if got := p.omitsAlterUserName(); got != c.want {
				t.Errorf("omitsAlterUserName(%q) = %v, want %v", c.input, got, c.want)
			}
		})
	}
}

// TestStartsUserPolicy pins the policy-kind lookahead after ALTER USER SET /
// UNSET.
func TestStartsUserPolicy(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"AUTHENTICATION POLICY p", true},
		{"PASSWORD POLICY p", true},
		{"SESSION POLICY p", true},
		{"INGRESS POLICY p", true},
		{"PASSWORD = 's'", false},
		{"NETWORK_POLICY = 's'", false},
		{`"INGRESS" POLICY p`, false},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			p := &Parser{lexer: NewLexer(c.input), input: c.input}
			p.advance()
			if got := p.startsUserPolicy(true); got != c.want {
				t.Errorf("startsUserPolicy(%q) = %v, want %v", c.input, got, c.want)
			}
		})
	}
}

// TestParseAlterUserOptionalName covers the forms that omit the user name and
// so alter the current user (https://docs.snowflake.com/en/sql-reference/sql/alter-user).
func TestParseAlterUserOptionalName(t *testing.T) {
	cases := []struct {
		sql    string
		name   string // "" when omitted
		ifEx   bool
		action ast.AlterUserAction
	}{
		{"ALTER USER SET PASSWORD = 's'", "", false, ast.AlterUserSet},
		{"ALTER USER SET DEFAULT_SECONDARY_ROLES = ('ALL')", "", false, ast.AlterUserSet},
		{"ALTER USER SET WORKLOAD_IDENTITY = (TYPE = AWS ARN = 's')", "", false, ast.AlterUserSet},
		{"ALTER USER UNSET DEFAULT_ROLE, COMMENT", "", false, ast.AlterUserUnset},
		{"ALTER USER SET TAG a = 's'", "", false, ast.AlterUserSetTag},
		{"ALTER USER RENAME TO u2", "", false, ast.AlterUserRename},
		{"ALTER USER IF EXISTS RESET PASSWORD", "", true, ast.AlterUserResetPassword},
		{"ALTER USER ABORT ALL QUERIES", "", false, ast.AlterUserAbortQueries},
		{"ALTER USER reset RESET PASSWORD", "reset", false, ast.AlterUserResetPassword},
		{"ALTER USER IF ADD PAT t", "IF", false, ast.AlterUserAddCredential},
		{"ALTER USER IF EXISTS u SET COMMENT = 's'", "u", true, ast.AlterUserSet},
		{"ALTER USER REMOVE DELEGATED AUTHORIZATIONS OF ANY ROLE FROM SECURITY INTEGRATION i", "", false, ast.AlterUserRemoveDelegated},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			stmt := mustAlterUser(t, c.sql)
			got := ""
			if stmt.Name != nil {
				got = stmt.Name.String()
			}
			if got != c.name || stmt.IfExists != c.ifEx || stmt.Action != c.action {
				t.Errorf("Name=%q IfExists=%v Action=%v, want %q %v %v", got, stmt.IfExists, stmt.Action, c.name, c.ifEx, c.action)
			}
			if stmt.Loc.End != len(c.sql) {
				t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(c.sql))
			}
		})
	}
}

// TestParseAlterUserPolicy covers SET / UNSET { AUTHENTICATION | PASSWORD |
// SESSION | INGRESS } POLICY.
func TestParseAlterUserPolicy(t *testing.T) {
	s := mustAlterUser(t, "ALTER USER u SET AUTHENTICATION POLICY db.s.p FORCE")
	if s.Action != ast.AlterUserSetPolicy || s.PolicyKind != "AUTHENTICATION POLICY" || s.PolicyName.String() != "db.s.p" || !s.Force {
		t.Errorf("set: %+v", s)
	}
	s = mustAlterUser(t, "ALTER USER u SET INGRESS POLICY p")
	if s.Action != ast.AlterUserSetPolicy || s.PolicyKind != "INGRESS POLICY" || s.Force {
		t.Errorf("set ingress: %+v", s)
	}
	s = mustAlterUser(t, "ALTER USER UNSET SESSION POLICY")
	if s.Action != ast.AlterUserUnsetPolicy || s.PolicyKind != "SESSION POLICY" || s.Name != nil {
		t.Errorf("unset: %+v", s)
	}
	s = mustAlterUser(t, "ALTER USER u UNSET INGRESS POLICY")
	if s.Action != ast.AlterUserUnsetPolicy || s.PolicyKind != "INGRESS POLICY" {
		t.Errorf("unset ingress: %+v", s)
	}
}

// TestParseAlterUserCredential covers the MFA, key pair, programmatic access
// token, and workload identity subcommands.
func TestParseAlterUserCredential(t *testing.T) {
	cases := []struct {
		sql       string
		action    ast.AlterUserAction
		kind      ast.UserCredentialKind
		credName  string
		ifEx      bool
		ifNotEx   bool
		nOptions  int
		nUnset    int
		newName   string
		checkKind bool
	}{
		{sql: "ALTER USER u ADD KEY PAIR k PUBLIC_KEY = 's' ROLE_RESTRICTION = 'r' DAYS_TO_EXPIRY = 1 COMMENT = 'c'", action: ast.AlterUserAddCredential, kind: ast.UserCredentialKeyPair, credName: "k", nOptions: 4, checkKind: true},
		{sql: "ALTER USER ROTATE KEY PAIR k PUBLIC_KEY = 's' EXPIRE_ROTATED_KEY_PAIR_AFTER_HOURS = 1", action: ast.AlterUserRotateCredential, kind: ast.UserCredentialKeyPair, credName: "k", nOptions: 2, checkKind: true},
		{sql: "ALTER USER u REMOVE KEY PAIR k", action: ast.AlterUserRemoveCredential, kind: ast.UserCredentialKeyPair, credName: "k", checkKind: true},
		{sql: "ALTER USER u MODIFY KEY PAIR k RENAME TO k2", action: ast.AlterUserModifyCredentialRename, kind: ast.UserCredentialKeyPair, credName: "k", newName: "k2", checkKind: true},
		{sql: "ALTER USER u MODIFY KEY PAIR k SET DISABLED = TRUE COMMENT = 's'", action: ast.AlterUserModifyCredentialSet, kind: ast.UserCredentialKeyPair, credName: "k", nOptions: 2, checkKind: true},
		{sql: "ALTER USER u MODIFY KEY PAIR k UNSET DISABLED, COMMENT", action: ast.AlterUserModifyCredentialUnset, kind: ast.UserCredentialKeyPair, credName: "k", nUnset: 2, checkKind: true},
		{sql: "ALTER USER u ADD PROGRAMMATIC ACCESS TOKEN IF NOT EXISTS t DAYS_TO_EXPIRY = 1", action: ast.AlterUserAddCredential, kind: ast.UserCredentialPAT, credName: "t", ifNotEx: true, nOptions: 1, checkKind: true},
		{sql: "ALTER USER ADD PAT t", action: ast.AlterUserAddCredential, kind: ast.UserCredentialPAT, credName: "t", checkKind: true},
		{sql: "ALTER USER u ROTATE PAT IF EXISTS t EXPIRE_ROTATED_TOKEN_AFTER_HOURS = 1", action: ast.AlterUserRotateCredential, kind: ast.UserCredentialPAT, credName: "t", ifEx: true, nOptions: 1, checkKind: true},
		{sql: "ALTER USER u MODIFY PAT IF EXISTS t SET MINS_TO_BYPASS_NETWORK_POLICY_REQUIREMENT = 1", action: ast.AlterUserModifyCredentialSet, kind: ast.UserCredentialPAT, credName: "t", ifEx: true, nOptions: 1, checkKind: true},
		{sql: "ALTER USER u REMOVE PROGRAMMATIC ACCESS TOKEN t", action: ast.AlterUserRemoveCredential, kind: ast.UserCredentialPAT, credName: "t", checkKind: true},
		{sql: "ALTER USER u ADD WORKLOAD IDENTITY w TYPE = AWS ARN = 's' OIDC_AUDIENCE_LIST = ('a', 'b')", action: ast.AlterUserAddCredential, kind: ast.UserCredentialWorkloadIdentity, credName: "w", nOptions: 3, checkKind: true},
		{sql: "ALTER USER u MODIFY WORKLOAD IDENTITY w UNSET DISABLED", action: ast.AlterUserModifyCredentialUnset, kind: ast.UserCredentialWorkloadIdentity, credName: "w", nUnset: 1, checkKind: true},
		{sql: "ALTER USER u REMOVE WORKLOAD IDENTITY w", action: ast.AlterUserRemoveCredential, kind: ast.UserCredentialWorkloadIdentity, credName: "w", checkKind: true},
		{sql: "ALTER USER ADD MFA METHOD OTP COUNT = 1", action: ast.AlterUserAddCredential, kind: ast.UserCredentialMFAMethod, credName: "OTP", nOptions: 1, checkKind: true},
		{sql: "ALTER USER u REMOVE MFA METHOD m", action: ast.AlterUserRemoveCredential, kind: ast.UserCredentialMFAMethod, credName: "m", checkKind: true},
		{sql: "ALTER USER u MODIFY MFA METHOD m SET COMMENT = 's'", action: ast.AlterUserModifyCredentialSet, kind: ast.UserCredentialMFAMethod, credName: "m", nOptions: 1, checkKind: true},
		{sql: "ALTER USER ENROLL MFA", action: ast.AlterUserEnrollMFA},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			s := mustAlterUser(t, c.sql)
			if s.Action != c.action {
				t.Fatalf("Action = %v, want %v", s.Action, c.action)
			}
			if c.checkKind {
				if s.Credential != c.kind || s.CredentialName == nil || s.CredentialName.String() != c.credName {
					t.Errorf("Credential=%v CredentialName=%v, want %v %q", s.Credential, s.CredentialName, c.kind, c.credName)
				}
			}
			if s.CredentialIfExists != c.ifEx || s.CredentialIfNotExists != c.ifNotEx {
				t.Errorf("IfExists=%v IfNotExists=%v", s.CredentialIfExists, s.CredentialIfNotExists)
			}
			if len(s.Options) != c.nOptions || len(s.UnsetProps) != c.nUnset {
				t.Errorf("Options=%d UnsetProps=%d, want %d %d", len(s.Options), len(s.UnsetProps), c.nOptions, c.nUnset)
			}
			if c.newName != "" && (s.NewName == nil || s.NewName.String() != c.newName) {
				t.Errorf("NewName = %v, want %q", s.NewName, c.newName)
			}
			if s.Loc.End != len(c.sql) {
				t.Errorf("Loc.End = %d, want %d", s.Loc.End, len(c.sql))
			}
		})
	}
}

// TestParseAlterUserCredentialNegatives keeps the subcommands strict.
func TestParseAlterUserCredentialNegatives(t *testing.T) {
	bad := []string{
		"ALTER USER",                                           // no action
		"ALTER USER SET",                                       // SET nothing
		"ALTER USER u SET AUTHENTICATION POLICY",               // policy name missing
		"ALTER USER u ADD KEY PAIR k",                          // PUBLIC_KEY required
		"ALTER USER u ADD KEY PAIR k COMMENT = 's'",            // PUBLIC_KEY must come first
		"ALTER USER u ADD WORKLOAD IDENTITY w",                 // TYPE required
		"ALTER USER u ROTATE KEY PAIR k",                       // PUBLIC_KEY required
		"ALTER USER u ROTATE WORKLOAD IDENTITY w",              // no ROTATE for workload identities
		"ALTER USER u ROTATE MFA METHOD m",                     // no ROTATE for MFA methods
		"ALTER USER u ADD MFA METHOD m",                        // only OTP is added
		"ALTER USER u MODIFY MFA METHOD m RENAME TO m2",        // MFA methods take only SET COMMENT
		"ALTER USER u MODIFY MFA METHOD m SET DISABLED = TRUE", // MFA methods take only SET COMMENT
		"ALTER USER u MODIFY PAT t SET",                        // SET nothing
		"ALTER USER u MODIFY PAT t",                            // no MODIFY action
		"ALTER USER u REMOVE KEY PAIR",                         // name missing
		"ALTER USER u ADD KEY k PUBLIC_KEY = 's'",              // KEY without PAIR
		"ALTER USER u ADD PROGRAMMATIC TOKEN t",                // ACCESS missing
		"ALTER USER u ENROLL",                                  // MFA missing
		"ALTER USER SET INGRESS POLICY p",                      // INGRESS POLICY names its user
		"ALTER USER UNSET INGRESS POLICY",                      // INGRESS POLICY names its user
		"ALTER USER u SET INGRESS POLICY p FORCE zzz",          // trailing word
		"ALTER USER u ADD MFA METHOD OTP CREATE = 1",           // only COUNT
		"ALTER USER u ADD PAT t zzz",                           // a bare word is no property
		"ALTER USER u SET MUST_CHANGE_PASSWORD TRUE",           // '=' required
		"ALTER USER u ADD WORKLOAD IDENTITY w TYPE AWS",        // '=' required
		"ALTER USER u REMOVE DELEGATED AUTHORIZATIONS OF ROLE r FROM SECURITY INTEGRATION i", // OF ANY ROLE only
	}
	for _, sql := range bad {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("expected parse error for %q", sql)
			}
		})
	}
}
