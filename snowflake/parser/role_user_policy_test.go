package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// T4.6 — CREATE / ALTER ROLE / USER / { MASKING | ROW ACCESS | SESSION |
// PASSWORD | NETWORK | AUTHENTICATION } POLICY.
//
// Oracle: triangulation of the legacy SnowflakeParser.g4 rules (create_role /
// create_user / create_masking_policy / create_row_access_policy /
// create_{session,password,network}_policy + alter_*) against the official docs
// (docs win on conflict) and the testdata/official corpus. No live Snowflake
// instance is available, so divergences where the docs are followed over the
// (stale) ANTLR grammar are flagged in the migration ledger.
// ---------------------------------------------------------------------------

func mustCreateRole(t *testing.T, input string) *ast.CreateRoleStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.CreateRoleStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.CreateRoleStmt", input)
	}
	return stmt
}

func mustCreateUser(t *testing.T, input string) *ast.CreateUserStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.CreateUserStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.CreateUserStmt", input)
	}
	return stmt
}

func mustCreatePolicy(t *testing.T, input string) *ast.CreatePolicyStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.CreatePolicyStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.CreatePolicyStmt", input)
	}
	return stmt
}

func mustAlterRole(t *testing.T, input string) *ast.AlterRoleStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.AlterRoleStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.AlterRoleStmt", input)
	}
	return stmt
}

func mustAlterUser(t *testing.T, input string) *ast.AlterUserStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.AlterUserStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.AlterUserStmt", input)
	}
	return stmt
}

func mustAlterPolicy(t *testing.T, input string) *ast.AlterPolicyStmt {
	t.Helper()
	stmt, ok := mustParseOne(t, input).(*ast.AlterPolicyStmt)
	if !ok {
		t.Fatalf("parse %q: want *ast.AlterPolicyStmt", input)
	}
	return stmt
}

// assertSecurityLoc asserts an AST node's Loc has non-negative, ordered bounds.
// Every new T4.6 node must carry a valid Loc (Start/End >= 0, Start <= End).
func assertSecurityLoc(t *testing.T, loc ast.Loc, what string) {
	t.Helper()
	if loc.Start < 0 || loc.End < 0 {
		t.Errorf("%s: negative Loc %+v", what, loc)
	}
	if loc.Start > loc.End {
		t.Errorf("%s: Loc.Start %d > Loc.End %d", what, loc.Start, loc.End)
	}
}

// ---------------------------------------------------------------------------
// CREATE ROLE
//
// Docs grammar (truth1, authoritative):
//
//	CREATE [ OR REPLACE ] ROLE [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
//	CREATE [ OR REPLACE ] DATABASE ROLE [ IF NOT EXISTS ] [<db>.]<name>
//	  [ COMMENT = '<string_literal>' ]
// ---------------------------------------------------------------------------

func TestParseCreateRole(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE ROLE myrole")
		if stmt.Name.String() != "myrole" {
			t.Errorf("Name = %q, want myrole", stmt.Name.String())
		}
		if stmt.OrReplace || stmt.Database || stmt.IfNotExists {
			t.Errorf("unexpected modifier: %+v", stmt)
		}
		if stmt.Comment != nil || stmt.Tags != nil {
			t.Errorf("unexpected comment/tags: %+v", stmt)
		}
		assertSecurityLoc(t, stmt.Loc, "CreateRoleStmt")
	})

	t.Run("or replace", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE OR REPLACE ROLE r")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
	})

	t.Run("if not exists", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE ROLE IF NOT EXISTS r")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
	})

	t.Run("comment", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE ROLE r COMMENT = 'desc'")
		if stmt.Comment == nil || *stmt.Comment != "desc" {
			t.Errorf("Comment = %v, want desc", stmt.Comment)
		}
	})

	t.Run("with tag", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE ROLE r WITH TAG (cost_center = 'finance')")
		if len(stmt.Tags) != 1 || stmt.Tags[0].Name.String() != "cost_center" || stmt.Tags[0].Value != "finance" {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
	})

	t.Run("bare tag (no WITH)", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE ROLE r TAG (t = 'v')")
		if len(stmt.Tags) != 1 {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
	})

	t.Run("comment then tag", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE OR REPLACE ROLE IF NOT EXISTS r COMMENT = 'c' WITH TAG (a = 'b')")
		if stmt.Comment == nil || *stmt.Comment != "c" || len(stmt.Tags) != 1 {
			t.Errorf("Comment=%v Tags=%+v", stmt.Comment, stmt.Tags)
		}
	})

	t.Run("tag then comment (docs/antlr order divergence)", func(t *testing.T) {
		// The legacy ANTLR create_role lists `with_tags? comment_clause?` (TAG before
		// COMMENT); the docs list COMMENT first. Both orders are accepted so neither
		// is regressed and a trailing COMMENT after TAG is not dropped.
		stmt := mustCreateRole(t, "CREATE ROLE r TAG (a = 'b') COMMENT = 'c'")
		if stmt.Comment == nil || *stmt.Comment != "c" || len(stmt.Tags) != 1 {
			t.Errorf("Comment=%v Tags=%+v", stmt.Comment, stmt.Tags)
		}
	})

	t.Run("quoted name", func(t *testing.T) {
		stmt := mustCreateRole(t, `CREATE ROLE "My Role"`)
		if stmt.Name.String() != `"My Role"` {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})

	t.Run("database role", func(t *testing.T) {
		stmt := mustCreateRole(t, "CREATE DATABASE ROLE dr")
		if !stmt.Database {
			t.Error("Database not set")
		}
		if stmt.Name.String() != "dr" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})

	t.Run("database role qualified name", func(t *testing.T) {
		// DATABASE ROLE names may be db-qualified (db_name.role_name) per the docs.
		stmt := mustCreateRole(t, "CREATE OR REPLACE DATABASE ROLE IF NOT EXISTS mydb.dr COMMENT = 'c'")
		if !stmt.Database || stmt.Name.String() != "mydb.dr" {
			t.Errorf("Database=%v Name=%q", stmt.Database, stmt.Name.String())
		}
		if stmt.Comment == nil || *stmt.Comment != "c" {
			t.Errorf("Comment = %v", stmt.Comment)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE USER
//
// Docs grammar (truth1):
//
//	CREATE [ OR REPLACE ] USER [ IF NOT EXISTS ] <name>
//	  [ objectProperties ] [ objectParams ] [ sessionParams ]
//	  [ [ WITH ] TAG ( ... ) ]
//
// Every property/parameter is an open-ended KEY = value pair (CopyOption).
// ---------------------------------------------------------------------------

func TestParseCreateUser(t *testing.T) {
	t.Run("minimal", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE USER u")
		if stmt.Name.String() != "u" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
		if len(stmt.Options) != 0 {
			t.Errorf("unexpected options: %+v", stmt.Options)
		}
		assertSecurityLoc(t, stmt.Loc, "CreateUserStmt")
	})

	t.Run("or replace if not exists", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE OR REPLACE USER IF NOT EXISTS u")
		if !stmt.OrReplace || !stmt.IfNotExists {
			t.Errorf("OrReplace=%v IfNotExists=%v", stmt.OrReplace, stmt.IfNotExists)
		}
	})

	t.Run("properties", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE USER u PASSWORD = 'abc' LOGIN_NAME = 'ln' DISPLAY_NAME = 'dn' DEFAULT_ROLE = r DEFAULT_WAREHOUSE = wh")
		if findOption(stmt.Options, "PASSWORD").Lit.Value != "abc" {
			t.Errorf("PASSWORD = %+v", findOption(stmt.Options, "PASSWORD"))
		}
		if findOption(stmt.Options, "DEFAULT_ROLE").Words != "R" {
			t.Errorf("DEFAULT_ROLE = %+v", findOption(stmt.Options, "DEFAULT_ROLE"))
		}
		for _, n := range []string{"LOGIN_NAME", "DISPLAY_NAME", "DEFAULT_WAREHOUSE"} {
			if findOption(stmt.Options, n) == nil {
				t.Errorf("missing option %s", n)
			}
		}
	})

	t.Run("rsa key + bool param", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE USER u RSA_PUBLIC_KEY = 'MIIB...' MUST_CHANGE_PASSWORD = TRUE")
		if findOption(stmt.Options, "RSA_PUBLIC_KEY") == nil {
			t.Error("missing RSA_PUBLIC_KEY")
		}
		mcp := findOption(stmt.Options, "MUST_CHANGE_PASSWORD")
		if mcp == nil || mcp.Words != "TRUE" {
			t.Errorf("MUST_CHANGE_PASSWORD = %+v", mcp)
		}
	})

	t.Run("default_secondary_roles parenthesized list", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE USER u DEFAULT_SECONDARY_ROLES = ('ALL')")
		dsr := findOption(stmt.Options, "DEFAULT_SECONDARY_ROLES")
		if dsr == nil || len(dsr.List) != 1 || dsr.List[0].Value != "ALL" {
			t.Errorf("DEFAULT_SECONDARY_ROLES = %+v", dsr)
		}
	})

	t.Run("default_secondary_roles empty list", func(t *testing.T) {
		// DEFAULT_SECONDARY_ROLES = () activates no roles (a documented form).
		stmt := mustCreateUser(t, "CREATE USER u DEFAULT_SECONDARY_ROLES = ()")
		dsr := findOption(stmt.Options, "DEFAULT_SECONDARY_ROLES")
		if dsr == nil || dsr.Group == nil {
			t.Errorf("DEFAULT_SECONDARY_ROLES = %+v (want non-nil empty group)", dsr)
		}
	})

	t.Run("with tag", func(t *testing.T) {
		stmt := mustCreateUser(t, "CREATE USER u LOGIN_NAME = 'ln' WITH TAG (dept = 'eng')")
		if len(stmt.Tags) != 1 || stmt.Tags[0].Value != "eng" {
			t.Errorf("Tags = %+v", stmt.Tags)
		}
		if findOption(stmt.Options, "LOGIN_NAME") == nil {
			t.Error("LOGIN_NAME dropped")
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE MASKING POLICY / ROW ACCESS POLICY
//
// Docs grammar (truth1):
//
//	CREATE [ OR REPLACE ] MASKING POLICY [ IF NOT EXISTS ] <name>
//	  AS ( <arg> <type> [ , ... ] ) RETURNS <type> -> <body>
//	  [ COMMENT = '...' ] [ EXEMPT_OTHER_POLICIES = { TRUE | FALSE } ]
//
//	CREATE [ OR REPLACE ] ROW ACCESS POLICY [ IF NOT EXISTS ] <name>
//	  AS ( <arg> <type> [ , ... ] ) RETURNS BOOLEAN -> <body> [ COMMENT = '...' ]
// ---------------------------------------------------------------------------

func TestParseCreateMaskingPolicy(t *testing.T) {
	t.Run("single arg", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE MASKING POLICY mp AS (val string) RETURNS string -> val")
		if stmt.Kind != ast.PolicyMasking {
			t.Errorf("Kind = %v, want Masking", stmt.Kind)
		}
		if stmt.Name.String() != "mp" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
		if len(stmt.Args) != 1 || stmt.Args[0].Name.String() != "val" {
			t.Fatalf("Args = %+v", stmt.Args)
		}
		if strings.ToUpper(stmt.Args[0].DataType.Name) != "STRING" {
			t.Errorf("arg type = %q", stmt.Args[0].DataType.Name)
		}
		if stmt.Returns == nil || strings.ToUpper(stmt.Returns.Name) != "STRING" {
			t.Errorf("Returns = %+v", stmt.Returns)
		}
		if stmt.Body == nil {
			t.Error("Body is nil")
		}
		assertSecurityLoc(t, stmt.Loc, "CreatePolicyStmt(masking)")
		assertSecurityLoc(t, stmt.Args[0].Loc, "PolicyArg")
	})

	t.Run("multiple args", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE MASKING POLICY mp AS (email varchar, visibility string) RETURNS varchar -> case when visibility = 'Public' then email else '***' end")
		if len(stmt.Args) != 2 {
			t.Fatalf("Args = %+v", stmt.Args)
		}
		if stmt.Args[1].Name.String() != "visibility" {
			t.Errorf("arg[1] = %q", stmt.Args[1].Name.String())
		}
		if _, ok := stmt.Body.(*ast.CaseExpr); !ok {
			t.Errorf("Body = %T, want *ast.CaseExpr", stmt.Body)
		}
	})

	t.Run("case body", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE OR REPLACE MASKING POLICY email_mask AS (val string) RETURNS string -> CASE WHEN current_role() IN ('ANALYST') THEN VAL ELSE '********' END")
		if !stmt.OrReplace {
			t.Error("OrReplace not set")
		}
		if _, ok := stmt.Body.(*ast.CaseExpr); !ok {
			t.Errorf("Body = %T", stmt.Body)
		}
		if len(stmt.Options) != 0 {
			t.Errorf("unexpected options: %+v", stmt.Options)
		}
	})

	t.Run("trailing comment and exempt_other_policies", func(t *testing.T) {
		// Trailing options follow the body. The expression parser must stop at
		// 'comment' rather than swallowing it.
		stmt := mustCreatePolicy(t, "CREATE MASKING POLICY mp AS (val string) RETURNS string -> case when 1=1 then val else '*' end COMMENT = 'c' EXEMPT_OTHER_POLICIES = true")
		if findOption(stmt.Options, "COMMENT") == nil || findOption(stmt.Options, "COMMENT").Lit.Value != "c" {
			t.Errorf("COMMENT option = %+v", findOption(stmt.Options, "COMMENT"))
		}
		eop := findOption(stmt.Options, "EXEMPT_OTHER_POLICIES")
		if eop == nil || eop.Words != "TRUE" {
			t.Errorf("EXEMPT_OTHER_POLICIES = %+v", eop)
		}
	})

	t.Run("if not exists + qualified name", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE MASKING POLICY IF NOT EXISTS governance.policies.email_mask AS (val string) RETURNS string -> val")
		if !stmt.IfNotExists {
			t.Error("IfNotExists not set")
		}
		if stmt.Name.String() != "governance.policies.email_mask" {
			t.Errorf("Name = %q", stmt.Name.String())
		}
	})
}

func TestParseCreateRowAccessPolicy(t *testing.T) {
	t.Run("case body returns boolean", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE OR REPLACE ROW ACCESS POLICY rap_it AS (empl_id varchar) RETURNS BOOLEAN -> case when 'it_admin' = current_role() then true else false end")
		if stmt.Kind != ast.PolicyRowAccess {
			t.Errorf("Kind = %v, want RowAccess", stmt.Kind)
		}
		if stmt.Returns == nil || strings.ToUpper(stmt.Returns.Name) != "BOOLEAN" {
			t.Errorf("Returns = %+v", stmt.Returns)
		}
		if _, ok := stmt.Body.(*ast.CaseExpr); !ok {
			t.Errorf("Body = %T", stmt.Body)
		}
	})

	t.Run("boolean expression body with subquery", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE OR REPLACE ROW ACCESS POLICY rap AS (sales_region varchar) RETURNS BOOLEAN -> 'exec' = current_role() or exists (select 1 from t where region = sales_region)")
		if stmt.Body == nil {
			t.Fatal("Body is nil")
		}
		if _, ok := stmt.Body.(*ast.BinaryExpr); !ok {
			t.Errorf("Body = %T, want *ast.BinaryExpr (OR)", stmt.Body)
		}
	})

	t.Run("trivial true body, two args", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE OR REPLACE ROW ACCESS POLICY rap_test2 AS (n number, v varchar) RETURNS BOOLEAN -> true")
		if len(stmt.Args) != 2 {
			t.Errorf("Args = %+v", stmt.Args)
		}
		if lit, ok := stmt.Body.(*ast.Literal); !ok || lit.Kind != ast.LitBool {
			t.Errorf("Body = %T", stmt.Body)
		}
	})
}

// ---------------------------------------------------------------------------
// CREATE { SESSION | PASSWORD | NETWORK | AUTHENTICATION } POLICY (option bags)
// ---------------------------------------------------------------------------

func TestParseCreateOptionBagPolicies(t *testing.T) {
	t.Run("session policy", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE SESSION POLICY sp SESSION_IDLE_TIMEOUT_MINS = 30 SESSION_UI_IDLE_TIMEOUT_MINS = 10 COMMENT = 'c'")
		if stmt.Kind != ast.PolicySession {
			t.Errorf("Kind = %v", stmt.Kind)
		}
		if stmt.Args != nil || stmt.Returns != nil || stmt.Body != nil {
			t.Errorf("option-bag policy should have nil Args/Returns/Body: %+v", stmt)
		}
		if findOption(stmt.Options, "SESSION_IDLE_TIMEOUT_MINS").Lit.Ival != 30 {
			t.Errorf("SESSION_IDLE_TIMEOUT_MINS = %+v", findOption(stmt.Options, "SESSION_IDLE_TIMEOUT_MINS"))
		}
		if findOption(stmt.Options, "COMMENT").Lit.Value != "c" {
			t.Errorf("COMMENT = %+v", findOption(stmt.Options, "COMMENT"))
		}
		assertSecurityLoc(t, stmt.Loc, "CreatePolicyStmt(session)")
	})

	t.Run("password policy", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE OR REPLACE PASSWORD POLICY IF NOT EXISTS pp PASSWORD_MIN_LENGTH = 10 PASSWORD_MAX_RETRIES = 3")
		if stmt.Kind != ast.PolicyPassword || !stmt.OrReplace || !stmt.IfNotExists {
			t.Errorf("flags: %+v", stmt)
		}
		if findOption(stmt.Options, "PASSWORD_MIN_LENGTH").Lit.Ival != 10 {
			t.Errorf("PASSWORD_MIN_LENGTH = %+v", findOption(stmt.Options, "PASSWORD_MIN_LENGTH"))
		}
	})

	t.Run("network policy (docs ALLOWED_NETWORK_RULE_LIST)", func(t *testing.T) {
		// The docs' current network policy options (ALLOWED_NETWORK_RULE_LIST /
		// BLOCKED_NETWORK_RULE_LIST) supersede the legacy ANTLR grammar's
		// ALLOWED_IP_LIST / BLOCKED_IP_LIST; the open-ended bag accepts both.
		stmt := mustCreatePolicy(t, "CREATE NETWORK POLICY np ALLOWED_NETWORK_RULE_LIST = ('a') BLOCKED_NETWORK_RULE_LIST = ('b')")
		if stmt.Kind != ast.PolicyNetwork {
			t.Errorf("Kind = %v", stmt.Kind)
		}
		anrl := findOption(stmt.Options, "ALLOWED_NETWORK_RULE_LIST")
		if anrl == nil || len(anrl.List) != 1 || anrl.List[0].Value != "a" {
			t.Errorf("ALLOWED_NETWORK_RULE_LIST = %+v", anrl)
		}
	})

	t.Run("network policy (legacy ALLOWED_IP_LIST)", func(t *testing.T) {
		stmt := mustCreatePolicy(t, "CREATE NETWORK POLICY np ALLOWED_IP_LIST = ('192.168.1.0/24') BLOCKED_IP_LIST = ('192.168.1.99')")
		if findOption(stmt.Options, "ALLOWED_IP_LIST") == nil {
			t.Error("ALLOWED_IP_LIST dropped")
		}
	})

	t.Run("authentication policy (no legacy ANTLR rule; docs only)", func(t *testing.T) {
		// AUTHENTICATION POLICY has no rule in the legacy grammar at all — it is a
		// docs-only object. AUTHENTICATION is a non-reserved identifier, so the
		// dispatcher keys on the AUTHENTICATION+POLICY token pair.
		stmt := mustCreatePolicy(t, "CREATE AUTHENTICATION POLICY ap AUTHENTICATION_METHODS = ('SAML', 'PASSWORD') CLIENT_TYPES = ('SNOWFLAKE_UI', 'DRIVERS') MFA_ENROLLMENT = 'REQUIRED' COMMENT = 'c'")
		if stmt.Kind != ast.PolicyAuthentication {
			t.Errorf("Kind = %v", stmt.Kind)
		}
		am := findOption(stmt.Options, "AUTHENTICATION_METHODS")
		if am == nil || len(am.List) != 2 {
			t.Errorf("AUTHENTICATION_METHODS = %+v", am)
		}
		mfa := findOption(stmt.Options, "MFA_ENROLLMENT")
		if mfa == nil || mfa.Lit == nil || mfa.Lit.Value != "REQUIRED" {
			t.Errorf("MFA_ENROLLMENT = %+v", mfa)
		}
	})

	t.Run("authentication policy nested group option", func(t *testing.T) {
		// MFA_POLICY = ( <sub-properties> ) is a nested key/value group.
		stmt := mustCreatePolicy(t, "CREATE AUTHENTICATION POLICY ap CLIENT_POLICY = (SNOWFLAKE_UI = (MINIMUM_VERSION = '1.0'))")
		cp := findOption(stmt.Options, "CLIENT_POLICY")
		if cp == nil || len(cp.Group) == 0 {
			t.Errorf("CLIENT_POLICY = %+v (want nested group)", cp)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER ROLE
//
// Docs/legacy alter_role:
//
//	ALTER [ DATABASE ] ROLE [ IF EXISTS ] <name>
//	  { RENAME TO <n> | SET COMMENT = '...' | UNSET COMMENT
//	  | SET TAG ... | UNSET TAG ... }
// ---------------------------------------------------------------------------

func TestParseAlterRole(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER ROLE r RENAME TO r2")
		if stmt.Action != ast.AlterRoleRename || stmt.NewName.String() != "r2" {
			t.Errorf("Action=%v NewName=%v", stmt.Action, stmt.NewName)
		}
		assertSecurityLoc(t, stmt.Loc, "AlterRoleStmt")
	})

	t.Run("if exists set comment", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER ROLE IF EXISTS r SET COMMENT = 'c'")
		if !stmt.IfExists || stmt.Action != ast.AlterRoleSetComment {
			t.Errorf("IfExists=%v Action=%v", stmt.IfExists, stmt.Action)
		}
		if stmt.Comment == nil || *stmt.Comment != "c" {
			t.Errorf("Comment = %v", stmt.Comment)
		}
	})

	t.Run("unset comment", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER ROLE r UNSET COMMENT")
		if stmt.Action != ast.AlterRoleUnset {
			t.Errorf("Action = %v", stmt.Action)
		}
	})

	t.Run("set tag", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER ROLE r SET TAG a = 'b', c = 'd'")
		if stmt.Action != ast.AlterRoleSetTag || len(stmt.Tags) != 2 {
			t.Errorf("Action=%v Tags=%+v", stmt.Action, stmt.Tags)
		}
	})

	t.Run("unset tag", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER ROLE r UNSET TAG a, b")
		if stmt.Action != ast.AlterRoleUnsetTag || len(stmt.UnsetTags) != 2 {
			t.Errorf("Action=%v UnsetTags=%+v", stmt.Action, stmt.UnsetTags)
		}
	})

	t.Run("database role rename", func(t *testing.T) {
		stmt := mustAlterRole(t, "ALTER DATABASE ROLE mydb.dr RENAME TO mydb.dr2")
		if !stmt.Database || stmt.Action != ast.AlterRoleRename {
			t.Errorf("Database=%v Action=%v", stmt.Database, stmt.Action)
		}
		if stmt.Name.String() != "mydb.dr" || stmt.NewName.String() != "mydb.dr2" {
			t.Errorf("Name=%q NewName=%q", stmt.Name.String(), stmt.NewName.String())
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER USER
// ---------------------------------------------------------------------------

func TestParseAlterUser(t *testing.T) {
	t.Run("rename", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u RENAME TO u2")
		if stmt.Action != ast.AlterUserRename || stmt.NewName.String() != "u2" {
			t.Errorf("Action=%v NewName=%v", stmt.Action, stmt.NewName)
		}
		assertSecurityLoc(t, stmt.Loc, "AlterUserStmt")
	})

	t.Run("reset password", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u RESET PASSWORD")
		if stmt.Action != ast.AlterUserResetPassword {
			t.Errorf("Action = %v", stmt.Action)
		}
	})

	t.Run("abort all queries", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER IF EXISTS u ABORT ALL QUERIES")
		if !stmt.IfExists || stmt.Action != ast.AlterUserAbortQueries {
			t.Errorf("IfExists=%v Action=%v", stmt.IfExists, stmt.Action)
		}
	})

	t.Run("set options", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u SET DEFAULT_ROLE = r DISABLED = TRUE DEFAULT_WAREHOUSE = wh")
		if stmt.Action != ast.AlterUserSet {
			t.Errorf("Action = %v", stmt.Action)
		}
		if findOption(stmt.Options, "DEFAULT_ROLE") == nil || findOption(stmt.Options, "DISABLED") == nil {
			t.Errorf("options = %+v", stmt.Options)
		}
	})

	t.Run("unset properties (keyword names)", func(t *testing.T) {
		// DEFAULT_ROLE / DISPLAY_NAME lex as keywords; the open-ended UNSET name
		// list must accept any name word, not just the DB/SCHEMA vocabulary.
		stmt := mustAlterUser(t, "ALTER USER u UNSET DEFAULT_ROLE, DISPLAY_NAME, COMMENT")
		if stmt.Action != ast.AlterUserUnset {
			t.Errorf("Action = %v", stmt.Action)
		}
		want := []string{"DEFAULT_ROLE", "DISPLAY_NAME", "COMMENT"}
		if len(stmt.UnsetProps) != 3 {
			t.Fatalf("UnsetProps = %+v", stmt.UnsetProps)
		}
		for i, w := range want {
			if stmt.UnsetProps[i] != w {
				t.Errorf("UnsetProps[%d] = %q, want %q", i, stmt.UnsetProps[i], w)
			}
		}
	})

	t.Run("set tag / unset tag", func(t *testing.T) {
		s1 := mustAlterUser(t, "ALTER USER u SET TAG dept = 'eng'")
		if s1.Action != ast.AlterUserSetTag || len(s1.Tags) != 1 {
			t.Errorf("set tag: %+v", s1)
		}
		s2 := mustAlterUser(t, "ALTER USER u UNSET TAG dept")
		if s2.Action != ast.AlterUserUnsetTag || len(s2.UnsetTags) != 1 {
			t.Errorf("unset tag: %+v", s2)
		}
	})

	t.Run("add delegated authorization", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u ADD DELEGATED AUTHORIZATION OF ROLE r TO SECURITY INTEGRATION i")
		if stmt.Action != ast.AlterUserAddDelegated {
			t.Errorf("Action = %v", stmt.Action)
		}
		if !strings.Contains(stmt.Raw, "DELEGATED") {
			t.Errorf("Raw = %q", stmt.Raw)
		}
	})

	t.Run("remove delegated authorization of role", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u REMOVE DELEGATED AUTHORIZATION OF ROLE r FROM SECURITY INTEGRATION i")
		if stmt.Action != ast.AlterUserRemoveDelegated {
			t.Errorf("Action = %v", stmt.Action)
		}
	})

	t.Run("remove delegated authorizations", func(t *testing.T) {
		stmt := mustAlterUser(t, "ALTER USER u REMOVE DELEGATED AUTHORIZATIONS FROM SECURITY INTEGRATION i")
		if stmt.Action != ast.AlterUserRemoveDelegated {
			t.Errorf("Action = %v", stmt.Action)
		}
	})
}

// ---------------------------------------------------------------------------
// ALTER { MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK | AUTHENTICATION }
// POLICY
// ---------------------------------------------------------------------------

func TestParseAlterPolicy(t *testing.T) {
	t.Run("masking set body", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER MASKING POLICY mp SET BODY -> case when 1=1 then a else '*' end")
		if stmt.Kind != ast.PolicyMasking || stmt.Action != ast.AlterPolicySetBody {
			t.Errorf("Kind=%v Action=%v", stmt.Kind, stmt.Action)
		}
		if _, ok := stmt.Body.(*ast.CaseExpr); !ok {
			t.Errorf("Body = %T", stmt.Body)
		}
		assertSecurityLoc(t, stmt.Loc, "AlterPolicyStmt")
	})

	t.Run("masking rename", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER MASKING POLICY IF EXISTS mp RENAME TO mp2")
		if !stmt.IfExists || stmt.Action != ast.AlterPolicyRename || stmt.NewName.String() != "mp2" {
			t.Errorf("%+v", stmt)
		}
	})

	t.Run("masking set comment", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER MASKING POLICY mp SET COMMENT = 'c'")
		if stmt.Action != ast.AlterPolicySet {
			t.Errorf("Action = %v", stmt.Action)
		}
		if findOption(stmt.Options, "COMMENT").Lit.Value != "c" {
			t.Errorf("COMMENT = %+v", findOption(stmt.Options, "COMMENT"))
		}
	})

	t.Run("row access set body", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER ROW ACCESS POLICY rap SET BODY -> true")
		if stmt.Kind != ast.PolicyRowAccess || stmt.Action != ast.AlterPolicySetBody {
			t.Errorf("Kind=%v Action=%v", stmt.Kind, stmt.Action)
		}
	})

	t.Run("session set / unset", func(t *testing.T) {
		s1 := mustAlterPolicy(t, "ALTER SESSION POLICY sp SET SESSION_IDLE_TIMEOUT_MINS = 5")
		if s1.Kind != ast.PolicySession || s1.Action != ast.AlterPolicySet {
			t.Errorf("set: %+v", s1)
		}
		s2 := mustAlterPolicy(t, "ALTER SESSION POLICY sp UNSET SESSION_IDLE_TIMEOUT_MINS, COMMENT")
		if s2.Action != ast.AlterPolicyUnset || len(s2.UnsetProps) != 2 {
			t.Errorf("unset: %+v", s2)
		}
	})

	t.Run("password unset", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER PASSWORD POLICY pp UNSET PASSWORD_MIN_LENGTH")
		if stmt.Kind != ast.PolicyPassword || stmt.Action != ast.AlterPolicyUnset {
			t.Errorf("%+v", stmt)
		}
	})

	t.Run("network set", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER NETWORK POLICY np SET ALLOWED_IP_LIST = ('1.2.3.4')")
		if stmt.Kind != ast.PolicyNetwork || stmt.Action != ast.AlterPolicySet {
			t.Errorf("%+v", stmt)
		}
	})

	t.Run("policy set tag / unset tag", func(t *testing.T) {
		s1 := mustAlterPolicy(t, "ALTER MASKING POLICY mp SET TAG t = 'v'")
		if s1.Action != ast.AlterPolicySetTag || len(s1.Tags) != 1 {
			t.Errorf("set tag: %+v", s1)
		}
		s2 := mustAlterPolicy(t, "ALTER MASKING POLICY mp UNSET TAG t")
		if s2.Action != ast.AlterPolicyUnsetTag || len(s2.UnsetTags) != 1 {
			t.Errorf("unset tag: %+v", s2)
		}
	})

	t.Run("authentication set", func(t *testing.T) {
		stmt := mustAlterPolicy(t, "ALTER AUTHENTICATION POLICY ap SET COMMENT = 'c'")
		if stmt.Kind != ast.PolicyAuthentication || stmt.Action != ast.AlterPolicySet {
			t.Errorf("%+v", stmt)
		}
	})
}

// ---------------------------------------------------------------------------
// Negative tests — the parser must reject these (it must not be over-permissive)
// ---------------------------------------------------------------------------

func TestParseSecurityNegatives(t *testing.T) {
	bad := []string{
		"CREATE ROLE",          // missing name
		"CREATE DATABASE ROLE", // missing name
		"CREATE USER",          // missing name
		"CREATE MASKING POLICY mp RETURNS string -> val",              // missing AS (args)
		"CREATE MASKING POLICY mp AS (val string) -> val",             // missing RETURNS
		"CREATE MASKING POLICY mp AS (val string) RETURNS string",     // missing -> body
		"CREATE MASKING POLICY mp AS (val string) RETURNS string -> ", // missing body expr
		"CREATE MASKING POLICY mp AS () RETURNS string -> val",        // empty arg list
		"CREATE MASKING POLICY mp AS (val) RETURNS string -> val",     // arg missing type
		"CREATE ROW ACCESS POLICY rap AS (a int) RETURNS BOOLEAN val", // missing arrow
		"CREATE MASKING mp AS (a int) RETURNS int -> a",               // MASKING not followed by POLICY
		"CREATE ROW POLICY rap AS (a int) RETURNS BOOLEAN -> true",    // ROW not followed by ACCESS
		"CREATE SESSION sp COMMENT = 'c'",                             // SESSION not followed by POLICY
		"CREATE AUTHENTICATION ap COMMENT = 'c'",                      // AUTHENTICATION not followed by POLICY
		"ALTER ROLE r SET",                                            // SET nothing
		"ALTER ROLE r SET DEFAULT_ROLE = x",                           // ROLE has no settable options other than COMMENT/TAG
		"ALTER ROLE r FOO",                                            // unknown action
		"ALTER USER u SET",                                            // SET nothing
		"ALTER USER u UNSET",                                          // UNSET nothing
		"ALTER USER u ADD DELEGATED AUTHORIZATION OF ROLE r",          // truncated (missing TO SECURITY INTEGRATION)
		"ALTER USER u REMOVE DELEGATED FROM SECURITY INTEGRATION i",   // REMOVE DELEGATED missing AUTHORIZATION(S)
		"ALTER MASKING POLICY mp SET BODY val",                        // missing arrow
		"ALTER MASKING POLICY mp SET",                                 // SET nothing
		"ALTER MASKING POLICY mp FROBNICATE",                          // unknown action
	}
	for _, sql := range bad {
		t.Run(sql, func(t *testing.T) {
			res := ParseBestEffort(sql)
			if len(res.Errors) == 0 {
				t.Errorf("expected parse error for %q, got AST %T", sql, firstStmt(res))
			}
		})
	}
}

func firstStmt(r *ParseResult) ast.Node {
	if len(r.File.Stmts) > 0 {
		return r.File.Stmts[0]
	}
	return nil
}

// ---------------------------------------------------------------------------
// DATABASE ROLE vs DATABASE disambiguation. `CREATE DATABASE ROLE x` is a
// database role; a database literally named ROLE must be quoted
// (CREATE DATABASE "ROLE"). The dispatcher keys on the ROLE keyword after
// DATABASE, matching Snowflake's grammar precedence.
// ---------------------------------------------------------------------------

func TestSecurityDatabaseRoleDisambiguation(t *testing.T) {
	cases := []struct {
		sql      string
		wantRole bool // true => Create/AlterRoleStmt; false => Create/AlterDatabaseStmt
	}{
		{`CREATE DATABASE "ROLE"`, false},
		{"CREATE DATABASE mydb", false},
		{"CREATE DATABASE ROLE dr", true},
		{"CREATE OR REPLACE DATABASE ROLE IF NOT EXISTS db.dr", true},
		{"ALTER DATABASE ROLE db.dr RENAME TO db.dr2", true},
		{`ALTER DATABASE "ROLE" RENAME TO d2`, false},
		{"ALTER DATABASE mydb RENAME TO d2", false},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			node := mustParseOne(t, c.sql)
			_, isCreateRole := node.(*ast.CreateRoleStmt)
			_, isAlterRole := node.(*ast.AlterRoleStmt)
			gotRole := isCreateRole || isAlterRole
			if gotRole != c.wantRole {
				t.Errorf("%q: role-node=%v want=%v (got %T)", c.sql, gotRole, c.wantRole, node)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Dispatch isolation — same-prefix statements owned by other nodes must NOT be
// mis-parsed as a T4.6 node. (They may be unsupported or syntax errors, but they
// must never become a Create/AlterPolicyStmt.)
// ---------------------------------------------------------------------------

func TestSecurityDispatchIsolation(t *testing.T) {
	// ALTER SESSION SET ... is the session statement (a different node); it must
	// not be routed to ALTER SESSION POLICY.
	res := ParseBestEffort("ALTER SESSION SET TIMEZONE = 'UTC'")
	for _, s := range res.File.Stmts {
		if _, ok := s.(*ast.AlterPolicyStmt); ok {
			t.Errorf("ALTER SESSION SET misrouted to AlterPolicyStmt")
		}
	}
	// CREATE NETWORK RULE is a different object; it must not become a NETWORK
	// POLICY. (It is currently unsupported; the assertion is that it is not a
	// CreatePolicyStmt.)
	res = ParseBestEffort("CREATE NETWORK RULE nr TYPE = IPV4 VALUE_LIST = ('1.2.3.4')")
	for _, s := range res.File.Stmts {
		if _, ok := s.(*ast.CreatePolicyStmt); ok {
			t.Errorf("CREATE NETWORK RULE misrouted to CreatePolicyStmt")
		}
	}
}

// ---------------------------------------------------------------------------
// Walk coverage — the policy body expression must be reachable via Walk (so
// query-span / analysis can descend into it). A round-trip Walk must visit a
// node inside the body.
// ---------------------------------------------------------------------------

func TestSecurityWalkBody(t *testing.T) {
	stmt := mustCreatePolicy(t, "CREATE MASKING POLICY mp AS (val string) RETURNS string -> current_role()")
	found := false
	ast.Inspect(stmt, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncCallExpr); ok {
			found = true
		}
		return true
	})
	if !found {
		t.Error("policy body FuncCallExpr not reachable via Walk/Inspect")
	}
}

// ---------------------------------------------------------------------------
// Official-docs corpus — every CREATE statement in the security corpus dirs
// must parse to the expected AST type. Context statements (USE ...) and the
// DESC NETWORK POLICY statement are owned by other nodes and skipped.
// ---------------------------------------------------------------------------

func TestSecurityOfficialCorpus(t *testing.T) {
	dirs := []string{
		"testdata/official/create-role",
		"testdata/official/create-user",
		"testdata/official/create-masking-policy",
		"testdata/official/create-row-access-policy",
		"testdata/official/create-network-policy",
	}
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertSecurityStatementsParse(t, string(data))
			})
		}
	}
}

// assertSecurityStatementsParse parses sql and asserts that every CREATE ROLE /
// USER / POLICY and ALTER ... statement parses with no errors and to a T4.6 AST
// type. Statements owned by other nodes (USE, DESC/DESCRIBE, GRANT, and the
// CREATE OR ALTER preview form) are skipped with a note.
func assertSecurityStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)

		// CREATE OR ALTER preview form is owned by parser-or-alter; skip (and
		// surface if it unexpectedly starts parsing here).
		if strings.Contains(upper, "CREATE OR ALTER") {
			if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) == 0 {
				t.Logf("note: CREATE OR ALTER now parses, revisit skip: %q", text)
			}
			continue
		}

		var want func(ast.Node) bool
		switch {
		case strings.HasPrefix(upper, "CREATE ROLE"),
			strings.HasPrefix(upper, "CREATE OR REPLACE ROLE"),
			strings.HasPrefix(upper, "CREATE DATABASE ROLE"),
			strings.HasPrefix(upper, "CREATE OR REPLACE DATABASE ROLE"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.CreateRoleStmt); return ok }
		case strings.HasPrefix(upper, "CREATE USER"),
			strings.HasPrefix(upper, "CREATE OR REPLACE USER"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.CreateUserStmt); return ok }
		case strings.Contains(upper, "POLICY") && strings.HasPrefix(upper, "CREATE"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.CreatePolicyStmt); return ok }
		case strings.HasPrefix(upper, "ALTER") && strings.Contains(upper, "POLICY"):
			want = func(n ast.Node) bool { _, ok := n.(*ast.AlterPolicyStmt); return ok }
		default:
			// USE / DESC / DESCRIBE / GRANT context statements owned by other nodes.
			continue
		}

		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
			continue
		}
		if !want(node) {
			t.Errorf("statement %q parsed to unexpected type %T", text, node)
		}
	}
}
