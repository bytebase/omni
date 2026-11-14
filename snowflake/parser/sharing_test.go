package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func TestParseExternalConsumer(t *testing.T) {
	wantObjectClauses(t, "CREATE OR REPLACE EXTERNAL CONSUMER c REGION = 'r' COMMENT = 'x' EMAIL = 'e'", "")
	wantObjectClauses(t, "CREATE EXTERNAL CONSUMER IF NOT EXISTS db.s.c")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER IF EXISTS c SET EMAIL = 'e' COMMENT = 'x'", "SET")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c UNSET COMMENT, EMAIL", "UNSET")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c RENAME TO db.s.d", "RENAME TO")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c ADD PAT IF NOT EXISTS t DAYS_TO_EXPIRY = 30 COMMENT = 'x'", "ADD PAT IF NOT EXISTS", "")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c ADD PROGRAMMATIC ACCESS TOKEN t", "ADD PROGRAMMATIC ACCESS TOKEN")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c MODIFY PAT IF EXISTS t SET DAYS_TO_EXPIRY = 1 COMMENT = 'x'", "MODIFY PAT IF EXISTS", "SET")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c MODIFY PAT t UNSET COMMENT", "MODIFY PAT", "UNSET")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c MODIFY PAT t", "MODIFY PAT")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER IF EXISTS c REMOVE PAT IF EXISTS t", "REMOVE PAT IF EXISTS")
	wantObjectClauses(t, "ALTER EXTERNAL CONSUMER c ROTATE PROGRAMMATIC ACCESS TOKEN t EXPIRE_ROTATED_TOKEN_AFTER_HOURS = 2",
		"ROTATE PROGRAMMATIC ACCESS TOKEN", "")
	d, ok := mustParseOne(t, "DROP EXTERNAL CONSUMER IF EXISTS c").(*ast.DropStmt)
	if !ok || d.Object != "EXTERNAL CONSUMER" || !d.IfExists {
		t.Errorf("DROP EXTERNAL CONSUMER: got %#v", d)
	}
	for _, in := range []string{
		"SHOW USER PATS FOR EXTERNAL CONSUMER c",
		"SHOW USER PROGRAMMATIC ACCESS TOKENS FOR EXTERNAL CONSUMER db.s.c",
	} {
		s, ok := mustParseOne(t, in).(*ast.ShowStmt)
		if !ok || s.Scope != ast.ShowScopeOther || s.ScopeText != "EXTERNAL CONSUMER" || s.ScopeName == nil {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	// A schema named EXTERNAL is still a scope.
	if s, ok := mustParseOne(t, "SHOW TABLES IN external").(*ast.ShowStmt); !ok || s.Scope != ast.ShowScopeSchema {
		t.Errorf("SHOW TABLES IN external: got %+v", s)
	}
}

func TestParseContactMultiPartyAndRSS(t *testing.T) {
	wantObjectClauses(t, "CREATE OR REPLACE CONTACT c USERS = ('a', 'b') COMMENT = 'x'", "")
	wantObjectClauses(t, "CREATE CONTACT IF NOT EXISTS db.s.c URL = 'u'", "")
	wantObjectClauses(t, "CREATE CONTACT c")
	wantObjectClauses(t, "ALTER CONTACT IF EXISTS c SET EMAIL_DISTRIBUTION_LIST = 'e' COMMENT = 'x'", "SET")
	wantObjectClauses(t, "ALTER CONTACT c RENAME TO d", "RENAME TO")
	wantObjectClauses(t, "CREATE OR REPLACE MULTI PARTY APPROVAL POLICY p AS 'b'", "AS")
	wantObjectClauses(t, "CREATE MULTI PARTY APPROVAL POLICY IF NOT EXISTS db.s.p AS 'b'", "AS")
	wantObjectClauses(t, "ALTER MULTI PARTY APPROVAL POLICY p AS 'b'", "AS")
	wantObjectClauses(t, "ALTER MULTI PARTY APPROVAL POLICY p RENAME TO q", "RENAME TO")
	wantObjectClauses(t, "CREATE RESTRICTED SESSION SCOPE IF NOT EXISTS r COMMENT = 'x' AS 'd'", "", "AS")
	wantObjectClauses(t, "CREATE OR ALTER RSS r AS 'd'", "AS")
	wantObjectClauses(t, "ALTER RSS IF EXISTS r ADD AS 'd'", "ADD AS")
	wantObjectClauses(t, "ALTER RESTRICTED SESSION SCOPE r REMOVE AS 'd'", "REMOVE AS")
	wantObjectClauses(t, "ALTER RSS r AS 'd'", "AS")
	wantObjectClauses(t, "ALTER RSS r SET COMMENT = 'x'", "SET")
	wantObjectClauses(t, "ALTER RSS r UNSET COMMENT", "UNSET")
	for in, kind := range map[string]string{
		"DROP CONTACT c": "CONTACT",
		"DROP MULTI PARTY APPROVAL POLICY IF EXISTS p": "MULTI PARTY APPROVAL POLICY",
		"DROP RSS IF EXISTS r":                         "RSS",
		"DROP RESTRICTED SESSION SCOPE db.s.r":         "RESTRICTED SESSION SCOPE",
	} {
		d, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || d.Object != kind {
			t.Errorf("%q: got %#v", in, d)
		}
	}
}

func TestParseSharingRejects(t *testing.T) {
	for _, in := range []string{
		"CREATE EXTERNAL CONSUMER c ACCOUNT = 'a'",
		"CREATE EXTERNAL CONSUMER c REGION = r",
		"CREATE EXTERNAL CONSUMER c WITH TAG (t = 'v')",
		"ALTER EXTERNAL CONSUMER c",
		"ALTER EXTERNAL CONSUMER c SET REGION = 'r'",
		"ALTER EXTERNAL CONSUMER c UNSET REGION",
		"ALTER EXTERNAL CONSUMER c SET TAG t = 'v'",
		"ALTER EXTERNAL CONSUMER IF EXISTS c ADD PAT t",
		"ALTER EXTERNAL CONSUMER IF EXISTS c ROTATE PAT t",
		"ALTER EXTERNAL CONSUMER c ADD PAT IF EXISTS t",
		"ALTER EXTERNAL CONSUMER c REMOVE PAT IF NOT EXISTS t",
		"ALTER EXTERNAL CONSUMER c ADD PAT db.t",
		"ALTER EXTERNAL CONSUMER c ADD PAT t EXPIRE_ROTATED_TOKEN_AFTER_HOURS = 1",
		"ALTER EXTERNAL CONSUMER c ADD PAT t DAYS_TO_EXPIRY = 'x'",
		"ALTER EXTERNAL CONSUMER c ROTATE PAT t DAYS_TO_EXPIRY = 1",
		"ALTER EXTERNAL CONSUMER c MODIFY PAT t SET EMAIL = 'e'",
		"ALTER EXTERNAL CONSUMER c MODIFY PAT t UNSET DAYS_TO_EXPIRY",
		"ALTER EXTERNAL CONSUMER c MODIFY PAT t UNSET COMMENT SET COMMENT = 'x'",
		"ALTER EXTERNAL CONSUMER c ADD PROGRAMMATIC TOKEN t",
		"SHOW USER PATS FOR EXTERNAL CONSUMER",
		"SHOW TABLES FOR EXTERNAL CONSUMER c",
		"CREATE CONTACT c USERS = ('a') URL = 'u'",
		"CREATE CONTACT c USERS = (a)",
		"CREATE CONTACT c EMAIL = 'e'",
		"CREATE OR REPLACE CONTACT IF NOT EXISTS c",
		"ALTER CONTACT c",
		"ALTER CONTACT c SET URL = 'u' EMAIL_DISTRIBUTION_LIST = 'e'",
		"ALTER CONTACT c UNSET COMMENT",
		"DROP CONTACT IF EXISTS c",
		"CREATE MULTI PARTY APPROVAL POLICY p",
		"CREATE MULTI PARTY APPROVAL POLICY p AS b",
		"ALTER MULTI PARTY APPROVAL POLICY IF EXISTS p AS 'b'",
		"ALTER MULTI PARTY APPROVAL POLICY p SET COMMENT = 'x'",
		"CREATE RSS r",
		"CREATE RSS r AS 'd' COMMENT = 'x'",
		"CREATE OR ALTER RSS IF NOT EXISTS r AS 'd'",
		"CREATE OR REPLACE RSS r AS 'd'",
		"ALTER RSS r SET OWNER = 'x'",
		"ALTER RSS r ADD 'd'",
		"ALTER RSS r",
	} {
		mustRejectObject(t, in)
	}
}
