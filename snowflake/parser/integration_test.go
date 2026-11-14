package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

func mustCreateIntegration(t *testing.T, input string) *ast.CreateIntegrationStmt {
	t.Helper()
	s, ok := mustParseOne(t, input).(*ast.CreateIntegrationStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CreateIntegrationStmt", input, mustParseOne(t, input))
	}
	if s.Loc.Start != 0 || s.Loc.End != len(input) {
		t.Errorf("parse %q: Loc = %+v, want {0 %d}", input, s.Loc, len(input))
	}
	return s
}

// mustAlterIntegration parses an ALTER of an integration-style object, whose
// Loc starts at the object keyword after ALTER.
func mustAlterIntegration(t *testing.T, input string) *ast.AlterIntegrationStmt {
	t.Helper()
	s, ok := mustParseOne(t, input).(*ast.AlterIntegrationStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.AlterIntegrationStmt", input, mustParseOne(t, input))
	}
	if s.Loc.Start != len("ALTER ") || s.Loc.End != len(input) {
		t.Errorf("parse %q: Loc = %+v, want {%d %d}", input, s.Loc, len("ALTER "), len(input))
	}
	return s
}

func TestParseCreateCatalogIntegration(t *testing.T) {
	in := "CREATE OR REPLACE CATALOG INTEGRATION c CATALOG_SOURCE = ICEBERG_REST TABLE_FORMAT = ICEBERG " +
		"CATALOG_NAMESPACE = 'ns' REST_CONFIG = ( CATALOG_URI = 'https://x' CATALOG_API_TYPE = PUBLIC ) " +
		"REST_AUTHENTICATION = ( TYPE = OAUTH OAUTH_CLIENT_ID = 'id' OAUTH_CLIENT_SECRET = 's' " +
		"OAUTH_ALLOWED_SCOPES = ( 'PRINCIPAL_ROLE:ALL', 'x' ) ) ENABLED = TRUE REFRESH_INTERVAL_SECONDS = 30 COMMENT = 'c'"
	s := mustCreateIntegration(t, in)
	if s.Kind != ast.CatalogIntegration || !s.OrReplace || s.Name.String() != "c" || len(s.Options) != 8 {
		t.Fatalf("got %+v", s)
	}
	rest := findOption(s.Options, "REST_CONFIG")
	if rest == nil || len(rest.Group) != 2 || rest.Group[1].Name != "CATALOG_API_TYPE" || rest.Group[1].Words != "PUBLIC" {
		t.Fatalf("REST_CONFIG = %+v", rest)
	}
	auth := findOption(s.Options, "REST_AUTHENTICATION")
	if scopes := findOption(auth.Group, "OAUTH_ALLOWED_SCOPES"); scopes == nil || len(scopes.List) != 2 {
		t.Fatalf("OAUTH_ALLOWED_SCOPES = %+v", scopes)
	}

	for _, in := range []string{
		"CREATE CATALOG INTEGRATION IF NOT EXISTS c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = DELTA ENABLED = FALSE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = GLUE TABLE_FORMAT = ICEBERG REST_CONFIG = ( CATALOG_API_TYPE = AWS_GLUE CATALOG_NAME = '1' ) " +
			"REST_AUTHENTICATION = ( TYPE = SIGV4 SIGV4_IAM_ROLE = 'arn' ) ENABLED = TRUE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = POLARIS TABLE_FORMAT = ICEBERG REST_CONFIG = ( CATALOG_URI = 'u' CATALOG_NAME = 'n' ) " +
			"REST_AUTHENTICATION = ( TYPE = OAUTH OAUTH_GRANT_TYPE = TOKEN_EXCHANGE OAUTH_ALLOWED_SCOPES = ( 's' ) ) " +
			"DEFAULT_STORAGE_CONFIG = ( USE_PRIVATELINK_ENDPOINT = TRUE ) ENABLED = TRUE",
	} {
		if s := mustCreateIntegration(t, in); s.Kind != ast.CatalogIntegration {
			t.Errorf("%q: Kind = %v", in, s.Kind)
		}
	}

	for _, in := range []string{
		// CATALOG_SOURCE, TABLE_FORMAT and ENABLED are required.
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG",
		"CREATE CATALOG INTEGRATION c",
		// Every property is KEY = value.
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE zzz",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = POLARIS TABLE_FORMAT = ICEBERG REST_CONFIG = ( CATALOG_URI 'u' ) ENABLED = TRUE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = POLARIS TABLE_FORMAT = ICEBERG REST_CONFIG = ( ) ENABLED = TRUE",
		// The page documents no CREATE OR ALTER.
		"CREATE OR ALTER CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlterCatalogIntegration(t *testing.T) {
	s := mustAlterIntegration(t, "ALTER CATALOG INTEGRATION IF EXISTS c SET REST_AUTHENTICATION = ( OAUTH_CLIENT_SECRET = 's' ) "+
		"DEFAULT_STORAGE_CONFIG = ( USE_PRIVATELINK_ENDPOINT = FALSE ) REFRESH_INTERVAL_SECONDS = 60 COMMENT = 'c'")
	if s.Kind != ast.CatalogIntegration || !s.IfExists || s.Action != ast.AlterIntegrationSet || len(s.Options) != 4 {
		t.Fatalf("got %+v", s)
	}
	for _, in := range []string{
		"ALTER CATALOG INTEGRATION c SET",
		"ALTER CATALOG INTEGRATION c UNSET COMMENT",
		"ALTER CATALOG INTEGRATION c SET TAG t = 's'",
		"ALTER CATALOG INTEGRATION c REFRESH",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseExternalAccessIntegration(t *testing.T) {
	s := mustCreateIntegration(t, "CREATE OR REPLACE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( r1, db.s.r2 ) "+
		"ALLOWED_API_AUTHENTICATION_INTEGRATIONS = none ALLOWED_AUTHENTICATION_SECRETS = ( s1 ) ENABLED = TRUE COMMENT = 'c'")
	if s.Kind != ast.ExternalAccessIntegration || len(s.Options) != 5 {
		t.Fatalf("got %+v", s)
	}
	rules := findOption(s.Options, "ALLOWED_NETWORK_RULES")
	if rules == nil || len(rules.Group) != 2 || !rules.Group[1].Bare || rules.Group[1].Name != "DB.S.R2" {
		t.Fatalf("ALLOWED_NETWORK_RULES = %+v", rules)
	}

	for in, want := range map[string]ast.AlterIntegrationAction{
		"ALTER EXTERNAL ACCESS INTEGRATION e SET ALLOWED_NETWORK_RULES = ( r ) ENABLED = FALSE": ast.AlterIntegrationSet,
		"ALTER EXTERNAL ACCESS INTEGRATION e SET COMMENT = 'c' TAG t = 's', u = 's'":            ast.AlterIntegrationSet,
		"ALTER EXTERNAL ACCESS INTEGRATION e SET TAG t = 's'":                                   ast.AlterIntegrationSetTag,
		"ALTER EXTERNAL ACCESS INTEGRATION IF EXISTS e UNSET ALLOWED_NETWORK_RULES, COMMENT":    ast.AlterIntegrationUnset,
		"ALTER EXTERNAL ACCESS INTEGRATION e UNSET COMMENT, TAG t":                              ast.AlterIntegrationUnset,
		"ALTER EXTERNAL ACCESS INTEGRATION e UNSET TAG t, TAG u":                                ast.AlterIntegrationUnsetTag,
	} {
		s := mustAlterIntegration(t, in)
		if s.Kind != ast.ExternalAccessIntegration || s.Action != want {
			t.Errorf("%q: Kind = %v, Action = %v, want %v", in, s.Kind, s.Action, want)
		}
	}
	s2 := mustAlterIntegration(t, "ALTER EXTERNAL ACCESS INTEGRATION e SET COMMENT = 'c' TAG t = 's', u = 's'")
	if len(s2.Options) != 1 || len(s2.Tags) != 2 {
		t.Errorf("SET = %+v", s2)
	}

	for _, in := range []string{
		"CREATE EXTERNAL ACCESS INTEGRATION e ENABLED = TRUE",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( r )",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( r s ) ENABLED = TRUE",
		"ALTER EXTERNAL ACCESS INTEGRATION e UNSET ALLOWED_NETWORK_RULES zzz",
		"ALTER EXTERNAL ACCESS INTEGRATION e UNSET",
		"ALTER EXTERNAL ACCESS INTEGRATION e REFRESH",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseAlterSecurityIntegrationRefresh(t *testing.T) {
	for in, want := range map[string][]string{
		"ALTER INTEGRATION i REFRESH":                                          nil,
		"ALTER SECURITY INTEGRATION IF EXISTS i REFRESH OAUTH_CLIENT_SECRET_2": {"OAUTH_CLIENT_SECRET_2"},
		"ALTER INTEGRATION i REFRESH SAML2_SNOWFLAKE_PRIVATE_KEY METADATA_URL": {"SAML2_SNOWFLAKE_PRIVATE_KEY", "METADATA_URL"},
		"ALTER INTEGRATION i REFRESH METADATA_URL":                             {"METADATA_URL"},
	} {
		s := mustAlterIntegration(t, in)
		if s.Action != ast.AlterIntegrationRefresh || len(s.Refresh) != len(want) {
			t.Errorf("%q: Action = %v, Refresh = %q", in, s.Action, s.Refresh)
			continue
		}
		for i := range want {
			if s.Refresh[i] != want[i] {
				t.Errorf("%q: Refresh = %q, want %q", in, s.Refresh, want)
			}
		}
	}
	for _, in := range []string{
		"ALTER INTEGRATION i REFRESH OAUTH_CLIENT_SECRET METADATA_URL",
		"ALTER INTEGRATION i REFRESH METADATA_URL SAML2_SNOWFLAKE_PRIVATE_KEY",
		"ALTER NOTIFICATION INTEGRATION i REFRESH",
		"ALTER API INTEGRATION i REFRESH",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseNotificationIntegrationProperties(t *testing.T) {
	s := mustCreateIntegration(t, "CREATE NOTIFICATION INTEGRATION n ENABLED = TRUE TYPE = MULTI_QUEUE DIRECTION = INBOUND "+
		"QUEUES = ( ( NAME = 'a' NOTIFICATION_PROVIDER = AWS_SNS AWS_SNS_TOPIC_ARN = 't' ), "+
		"( NAME = 'b' NOTIFICATION_PROVIDER = GCP_PUBSUB GCP_PUBSUB_SUBSCRIPTION_NAME = 's' ) ) ACTIVE = 'a'")
	queues := findOption(s.Options, "QUEUES")
	if s.Kind != ast.NotificationIntegration || queues == nil || len(queues.Group) != 2 || len(queues.Group[1].Group) != 3 {
		t.Fatalf("QUEUES = %+v", queues)
	}
	s = mustCreateIntegration(t, "CREATE NOTIFICATION INTEGRATION n TYPE = WEBHOOK ENABLED = TRUE WEBHOOK_URL = 'u' "+
		"WEBHOOK_HEADERS = ( 'Content-Type' = 'application/json', 'X' = 'y' )")
	if h := findOption(s.Options, "WEBHOOK_HEADERS"); h == nil || len(h.Group) != 2 || h.Group[0].Name != "Content-Type" {
		t.Fatalf("WEBHOOK_HEADERS = %+v", h)
	}
	for _, in := range []string{
		"CREATE NOTIFICATION INTEGRATION n ENABLED TRUE TYPE = QUEUE",
		"CREATE NOTIFICATION INTEGRATION n TYPE = WEBHOOK ENABLED = TRUE WEBHOOK_URL = 'u' zzz",
		"CREATE NOTIFICATION INTEGRATION n ENABLED = TRUE QUEUES = ( ( NAME = 'a' ) ( NAME = 'b' ) )",
		"CREATE NOTIFICATION INTEGRATION n WEBHOOK_HEADERS = ( 'a' = 'b' 'c' = 'd' )",
	} {
		mustRejectObject(t, in)
	}
}

func TestParseDropIntegrationKinds(t *testing.T) {
	for _, in := range []string{
		"DROP INTEGRATION i",
		"DROP API INTEGRATION i",
		"DROP CATALOG INTEGRATION IF EXISTS i",
		"DROP EXTERNAL ACCESS INTEGRATION i",
		"DROP NOTIFICATION INTEGRATION i",
		"DROP SECURITY INTEGRATION i",
		"DROP STORAGE INTEGRATION IF EXISTS i",
	} {
		s, ok := mustParseOne(t, in).(*ast.DropStmt)
		if !ok || s.Kind != ast.DropIntegration || s.Loc.Start != 0 || s.Loc.End != len(in) {
			t.Errorf("%q: got %+v", in, s)
		}
	}
	for _, in := range []string{
		"DROP EXTERNAL ACCESS i",
		"DROP CATALOG i",
	} {
		mustRejectObject(t, in)
	}
}

func TestStartsCatalogOrAccessIntegration(t *testing.T) {
	for in, want := range map[string]bool{
		"CATALOG INTEGRATION c":         true,
		"EXTERNAL ACCESS INTEGRATION e": true,
		"EXTERNAL ACCESS e":             false,
		"EXTERNAL TABLE t":              false,
		"CATALOG c":                     false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsCatalogOrAccessIntegration(); got != want {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}

func TestParseIntegrationPropertiesByPage(t *testing.T) {
	for _, in := range []string{
		"CREATE NOTIFICATION INTEGRATION n TYPE = EMAIL ENABLED = TRUE ALLOWED_RECIPIENTS = ( 'a@b.c' ) DEFAULT_SUBJECT = 's'",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( r ) ALLOWED_AUTHENTICATION_SECRETS = all ENABLED = TRUE",
		"ALTER CATALOG INTEGRATION c SET REST_AUTHENTICATION = ( BEARER_TOKEN = 't' )",
	} {
		mustParseOne(t, in)
	}
	for _, in := range []string{
		// Unknown property names, at the top and inside a group.
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE NOT_A_PROPERTY = 1",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = POLARIS TABLE_FORMAT = ICEBERG REST_CONFIG = ( NOT_A_KEY = 's' ) ENABLED = TRUE",
		"ALTER CATALOG INTEGRATION c SET CATALOG_SOURCE = GLUE",
		"ALTER CATALOG INTEGRATION c SET REST_AUTHENTICATION = ( TYPE = OAUTH )",
		"ALTER EXTERNAL ACCESS INTEGRATION e SET a = 's'",
		"ALTER EXTERNAL ACCESS INTEGRATION e UNSET ENABLED",
		"CREATE NOTIFICATION INTEGRATION n TYPE = WEBHOOK ENABLED = TRUE WEBHOOK_URL = 'u' NOT_A_PROPERTY = 's'",
		"CREATE NOTIFICATION INTEGRATION n ENABLED = TRUE QUEUES = ( ( NAME = 'a' NOT_A_KEY = 's' ) )",
		// Quoted keys.
		`CREATE CATALOG INTEGRATION c "CATALOG_SOURCE" = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE`,
		// Wrong value shapes.
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = 'GLUE' TABLE_FORMAT = ICEBERG ENABLED = TRUE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = PARQUET ENABLED = TRUE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = 'x'",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE REFRESH_INTERVAL_SECONDS = 'x'",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = POLARIS TABLE_FORMAT = ICEBERG REST_CONFIG = 's' ENABLED = TRUE",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( 'r' ) ENABLED = TRUE",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = r ENABLED = TRUE",
		"CREATE EXTERNAL ACCESS INTEGRATION e ALLOWED_NETWORK_RULES = ( r ) ALLOWED_AUTHENTICATION_SECRETS = some ENABLED = TRUE",
		"CREATE NOTIFICATION INTEGRATION n TYPE = WEBHOOK ENABLED = TRUE WEBHOOK_URL = 'u' WEBHOOK_HEADERS = ( a = 'b' )",
		"CREATE NOTIFICATION INTEGRATION n TYPE = EMAIL ENABLED = TRUE ALLOWED_RECIPIENTS = ( a )",
		// A key twice.
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = OBJECT_STORE TABLE_FORMAT = ICEBERG ENABLED = TRUE ENABLED = FALSE",
		"CREATE CATALOG INTEGRATION c CATALOG_SOURCE = ICEBERG_REST TABLE_FORMAT = ICEBERG REST_CONFIG = ( CATALOG_URI = 'u' CATALOG_URI = 'u' ) ENABLED = TRUE",
	} {
		mustRejectObject(t, in)
	}
}
