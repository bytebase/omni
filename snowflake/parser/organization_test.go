package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseOrganizationObjects(t *testing.T) {
	for in, want := range map[string][]string{
		"CREATE ORGANIZATION USER IF NOT EXISTS u":                                                                                                 nil,
		"CREATE ORGANIZATION USER u EMAIL = 'e' LOGIN_NAME = 'l' TYPE = PERSON COMMENT = 'c'":                                                      {""},
		"ALTER ORGANIZATION USER IF EXISTS u SET DISPLAY_NAME = 'd' TYPE = SERVICE":                                                                {"SET"},
		"ALTER ORGANIZATION USER u UNSET EMAIL, TYPE":                                                                                              {"UNSET"},
		"CREATE ORGANIZATION USER GROUP g IS_GRANTABLE = TRUE":                                                                                     {""},
		"ALTER ORGANIZATION USER GROUP IF EXISTS g ADD ORGANIZATION USERS a, b":                                                                    {"ADD ORGANIZATION USERS"},
		"ALTER ORGANIZATION USER GROUP g REMOVE ORGANIZATION USERS a":                                                                              {"REMOVE ORGANIZATION USERS"},
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = ALL":                                                                                     {"SET VISIBILITY = ALL"},
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = ACCOUNTS a, b":                                                                           {"SET VISIBILITY = ACCOUNTS"},
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = REGION GROUPS 'x', 'y'":                                                                  {"SET VISIBILITY = REGION GROUPS"},
		"ALTER ORGANIZATION USER GROUP g SET IS_GRANTABLE = FALSE":                                                                                 {"SET"},
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' ADMIN_RSA_PUBLIC_KEY = 'k' EMAIL = 'e' EDITION = BUSINESS_CRITICAL REGION = aws_us_west_2": {""},
		"ALTER ORGANIZATION ACCOUNT SET RESOURCE_MONITOR = m":                                                                                      {"SET"},
		"ALTER ORGANIZATION ACCOUNT SET TIMEZONE = 'UTC' DATA_RETENTION_TIME_IN_DAYS = 1":                                                          nil,
		"ALTER ORGANIZATION ACCOUNT SET SESSION POLICY db.s.p":                                                                                     {"SET SESSION POLICY"},
		"ALTER ORGANIZATION ACCOUNT UNSET PASSWORD POLICY":                                                                                         {"UNSET PASSWORD POLICY"},
		"ALTER ORGANIZATION ACCOUNT UNSET PASSWORD, TIMEZONE":                                                                                      {"UNSET"},
		"ALTER ORGANIZATION ACCOUNT SET TAG t = 'v', db.s.u = 'w'":                                                                                 {"SET TAG"},
		"ALTER ORGANIZATION ACCOUNT UNSET TAG t, db.s.u":                                                                                           {"UNSET TAG"},
		"CREATE ORGANIZATION PROFILE p":                                                                                                            nil,
		"CREATE ORGANIZATION PROFILE IF NOT EXISTS p AS 'm' VERSION v PUBLISH = FALSE":                                                             {"AS", "VERSION", ""},
		"CREATE ORGANIZATION PROFILE p FROM @db.s.st/dir PUBLISH = TRUE":                                                                           {"FROM", ""},
		"ALTER ORGANIZATION PROFILE IF EXISTS p AS 'm'":                                                                                            {"AS"},
		"ALTER ORGANIZATION PROFILE p RENAME TO q":                                                                                                 {"RENAME TO"},
		"ALTER ORGANIZATION PROFILE p PUBLISH":                                                                                                     {"PUBLISH"},
		"ALTER ORGANIZATION PROFILE p ADD VERSION FROM @st":                                                                                        {"ADD VERSION", "FROM"},
		"ALTER ORGANIZATION PROFILE p ADD VERSION IF NOT EXISTS v FROM @st":                                                                        {"ADD VERSION IF NOT EXISTS", "FROM"},
		"ALTER ORGANIZATION PROFILE p ADD LIVE VERSION v FROM LAST":                                                                                {"ADD LIVE VERSION", "FROM LAST"},
		"ALTER ORGANIZATION PROFILE p COMMIT":                                                                                                      {"COMMIT"},
		"ALTER ORGANIZATION PROFILE p ABORT":                                                                                                       {"ABORT"},
		"ALTER ORGANIZATION SET DEFAULT_TENANT_TYPE = EXTERNAL":                                                                                    {"SET"},
		"ALTER ORGANIZATION SET DOMAIN_NAMES = ('a.com', 'b.com')":                                                                                 {"SET"},
		"ALTER ORGANIZATION UNSET DOMAIN_NAMES":                                                                                                    {"UNSET"},
		"ALTER ORGANIZATION SET EVENT ROUTING TABLE db.s.t FOR ALL APPLICATION LISTINGS":                                                           {"SET EVENT ROUTING TABLE", "FOR ALL APPLICATION LISTINGS"},
		"ALTER ORGANIZATION UNSET EVENT ROUTING TABLE FOR ALL APPLICATION LISTINGS":                                                                {"UNSET EVENT ROUTING TABLE FOR ALL APPLICATION LISTINGS"},
	} {
		s := mustObjectStmt(t, in)
		got := clauseKeywords(s)
		if want == nil {
			continue
		}
		if len(got) != len(want) {
			t.Errorf("%q: clauses = %q, want %q", in, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%q: clauses = %q, want %q", in, got, want)
				break
			}
		}
	}
	s := mustObjectStmt(t, "ALTER ORGANIZATION ACCOUNT SET PASSWORD POLICY p")
	if s.Name != nil || s.Kind != "ORGANIZATION ACCOUNT" || s.Clauses[0].Names[0].String() != "p" {
		t.Errorf("got %+v", s)
	}
	s = mustObjectStmt(t, "ALTER ORGANIZATION USER GROUP g SET VISIBILITY = REGION GROUPS 'x', 'y'")
	if c := s.Clauses[0]; len(c.Options) != 2 || c.Options[1].Lit.Value != "y" {
		t.Errorf("REGION GROUPS = %+v", c.Options)
	}
	s = mustObjectStmt(t, "CREATE ORGANIZATION PROFILE p FROM @db.s.st/dir")
	if o := s.Clauses[0].Options[0]; o.Words != "@db.s.st/dir" {
		t.Errorf("FROM = %q", o.Words)
	}
}

func TestParseDropOrganizationObjects(t *testing.T) {
	for in, kind := range map[string]string{
		"DROP ORGANIZATION USER IF EXISTS u":       "ORGANIZATION USER",
		"DROP ORGANIZATION USER GROUP IF EXISTS g": "ORGANIZATION USER GROUP",
		"DROP ORGANIZATION PROFILE p":              "ORGANIZATION PROFILE",
	} {
		d, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || d.Kind != ast.DropObject || d.Object != kind {
			t.Errorf("%q: got %#v", in, d)
		}
	}
}

func TestParseShowOrganizationScopes(t *testing.T) {
	for in, scope := range map[string]string{
		"SHOW ORGANIZATION USERS IN ORGANIZATION USER GROUP g":                  "ORGANIZATION USER GROUP",
		"SHOW VERSIONS IN ORGANIZATION PROFILE p":                               "ORGANIZATION PROFILE",
		"SHOW EVENT ROUTING TABLE ON ORGANIZATION FOR ALL APPLICATION LISTINGS": "ALL APPLICATION LISTINGS",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeText != scope {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"SHOW ORGANIZATION USERS IN ORGANIZATION USER GROUP db.g",
		"SHOW ORGANIZATION USERS IN ORGANIZATION USER GROUP",
		"SHOW VERSIONS IN ORGANIZATION PROFILE",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseOrganizationRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE OR REPLACE ORGANIZATION USER u",
		"CREATE ORGANIZATION USER db.u",
		"CREATE ORGANIZATION USER u NOT_A_PROPERTY = 's'",
		"CREATE ORGANIZATION USER u TYPE = ROBOT",
		"CREATE ORGANIZATION USER u EMAIL = 1",
		"CREATE ORGANIZATION USER u WITH TAG (t = 'v')",
		"ALTER ORGANIZATION USER u",
		"ALTER ORGANIZATION USER u SET LOGIN_NAME = 'l'", // CREATE's alone
		"ALTER ORGANIZATION USER u UNSET LOGIN_NAME",
		"ALTER ORGANIZATION USER IF EXISTS u UNSET EMAIL",
		"ALTER ORGANIZATION USER u SET TAG t = 'v'",
		"ALTER ORGANIZATION USER u RENAME TO v",
		"CREATE ORGANIZATION USER GROUP g IS_GRANTABLE = 'yes'",
		"ALTER ORGANIZATION USER GROUP g ADD ORGANIZATION USERS",
		"ALTER ORGANIZATION USER GROUP g ADD ORGANIZATION USERS db.a",
		"ALTER ORGANIZATION USER GROUP g ADD USERS a",
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = NONE",
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = REGION GROUPS x",
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = REGION GROUPS 'x',",
		"ALTER ORGANIZATION USER GROUP g SET VISIBILITY = ACCOUNTS",
		"ALTER ORGANIZATION USER GROUP g SET COMMENT = 'c'",
		"ALTER ORGANIZATION USER GROUP g UNSET IS_GRANTABLE",
		"DROP ORGANIZATION USER GROUP g CASCADE",
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' EMAIL = 'e' EDITION = ENTERPRISE", // no password or key
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' ADMIN_PASSWORD = 'p' ADMIN_RSA_PUBLIC_KEY = 'k' EMAIL = 'e' EDITION = ENTERPRISE",
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' ADMIN_PASSWORD = 'p' EMAIL = 'e'", // EDITION is required
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' ADMIN_PASSWORD = 'p' EMAIL = 'e' EDITION = STANDARD",
		"CREATE ORGANIZATION ACCOUNT IF NOT EXISTS a ADMIN_NAME = 'n' ADMIN_PASSWORD = 'p' EMAIL = 'e' EDITION = ENTERPRISE",
		"CREATE ORGANIZATION ACCOUNT a ADMIN_NAME = 'n' ADMIN_PASSWORD = 'p' EMAIL = 'e' EDITION = ENTERPRISE REGION = a.b",
		"ALTER ORGANIZATION ACCOUNT a SET TIMEZONE = 'UTC'", // names no account
		"ALTER ORGANIZATION ACCOUNT IF EXISTS SET TIMEZONE = 'UTC'",
		"ALTER ORGANIZATION ACCOUNT SET",
		"ALTER ORGANIZATION ACCOUNT SET PASSWORD POLICY",
		"ALTER ORGANIZATION ACCOUNT UNSET SESSION POLICY p",
		"ALTER ORGANIZATION ACCOUNT SET TIMEZONE = 'UTC' TAG t = 'v'",
		"DROP ORGANIZATION ACCOUNT a",
		"CREATE ORGANIZATION PROFILE p AS 'm' COMMENT = 'c'",
		"CREATE ORGANIZATION PROFILE p AS m",
		"CREATE ORGANIZATION PROFILE p FROM 'm'",
		"CREATE ORGANIZATION PROFILE p PUBLISH = TRUE", // PUBLISH follows AS or FROM
		"CREATE ORGANIZATION PROFILE p AS 'm' VERSION db.v",
		"ALTER ORGANIZATION PROFILE p",
		"ALTER ORGANIZATION PROFILE IF EXISTS p COMMIT",
		"ALTER ORGANIZATION PROFILE IF EXISTS p ADD VERSION FROM @st",
		"ALTER ORGANIZATION PROFILE p ADD VERSION v",
		"ALTER ORGANIZATION PROFILE p ADD VERSION v GARBAGE @st",
		"ALTER ORGANIZATION PROFILE p ADD VERSION v FROM @",
		"CREATE ORGANIZATION PROFILE p FROM @ PUBLISH = TRUE",
		"ALTER ORGANIZATION PROFILE p ADD VERSION IF NOT EXISTS FROM @st",
		"ALTER ORGANIZATION PROFILE p ADD LIVE VERSION FROM @st",
		"ALTER ORGANIZATION PROFILE p RENAME TO db.q",
		"DROP ORGANIZATION PROFILE IF EXISTS p",
		"ALTER ORGANIZATION SET DEFAULT_TENANT_TYPE = OTHER",
		"ALTER ORGANIZATION SET DEFAULT_TENANT_TYPE = INTERNAL DOMAIN_NAMES = ('a')",
		"ALTER ORGANIZATION SET DOMAIN_NAMES = ()",
		"ALTER ORGANIZATION SET COMMENT = 'c'",
		"ALTER ORGANIZATION UNSET DEFAULT_TENANT_TYPE",
		"ALTER ORGANIZATION o SET DEFAULT_TENANT_TYPE = INTERNAL",
		"ALTER ORGANIZATION SET EVENT ROUTING TABLE t",
		"ALTER ORGANIZATION UNSET EVENT ROUTING TABLE t FOR ALL APPLICATION LISTINGS",
		"CREATE ORGANIZATION o",
	} {
		mustRejectObject(t, in)
	}
}
