package deparse_test

import (
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
	"github.com/stretchr/testify/require"
)

// VARIABLE, TYPE, INDEX, MCP SERVER, GATEWAY, OPENFLOW CONNECTOR and EVENT
// ROUTING TABLE round-trips.

func TestDeparse_MiscObjects(t *testing.T) {
	for _, sql := range []string{
		"CREATE VARIABLE myschema.tax_rate FLOAT DEFAULT 0.08",
		"CREATE OR REPLACE VARIABLE myschema.api_version VARCHAR DEFAULT 'v2' COMMENT = 'REST API version'",
		"CREATE VARIABLE IF NOT EXISTS myschema.batch_cursor NUMBER DEFAULT 0",
		"CREATE VARIABLE a NUMBER(10, 2) COMMENT = 'x'",
		"ALTER VARIABLE myschema.tax_rate RENAME TO myschema.regional_tax_rate",
		"ALTER VARIABLE IF EXISTS myschema.tax_rate SET DEFAULT 0.10",
		"ALTER VARIABLE myschema.batch_cursor DROP DEFAULT",
		"ALTER VARIABLE a SET COMMENT = 'Regional tax rate'",
		"ALTER VARIABLE a UNSET COMMENT",
		"DROP VARIABLE IF EXISTS myschema.batch_cursor",
		"UNDROP VARIABLE myschema.tax_rate",
		"CREATE TYPE age AS NUMBER(3, 0)",
		"CREATE OR REPLACE TYPE a AS VARCHAR COMMENT = 'x'",
		"CREATE TYPE IF NOT EXISTS a AS ARRAY(INT)",
		"ALTER TYPE age SET COMMENT = 'User-defined type for storing age values'",
		"ALTER TYPE IF EXISTS age UNSET COMMENT",
		"DROP TYPE age",
		"UNDROP TYPE age",
		"CREATE OR REPLACE INDEX vidx ON mytable (val)",
		"CREATE INDEX IF NOT EXISTS i ON db.s.t (a, \"B c\") INCLUDE (d, e)",
		"DROP INDEX mytable.vidx",
		"DROP INDEX IF EXISTS db.s.t.i",
		"CREATE MCP SERVER my_mcp_server FROM SPECIFICATION $$tools: []$$",
		"CREATE MCP SERVER IF NOT EXISTS m FROM SPECIFICATION 'x'",
		"DROP MCP SERVER IF EXISTS my_mcp_server",
		"CREATE GATEWAY split_gateway FROM SPECIFICATION $$spec: {}$$",
		"ALTER GATEWAY IF EXISTS g FROM SPECIFICATION 'x'",
		"DROP GATEWAY g",
		"CREATE OPENFLOW CONNECTOR my_db.my_schema.c IN RUNTIME my_db.my_schema.r FROM DEFINITION OPENFLOW_POSTGRES_CDC DISPLAY_NAME = 'My CDC' COMMENT = 'x'",
		"CREATE OPENFLOW CONNECTOR IF NOT EXISTS c IN RUNTIME r FROM '@my_git_repo/branches/main/' COMMENT = 'Created from Git stage'",
		"ALTER OPENFLOW CONNECTOR c ADD VERSION FROM '@stage/path/'",
		"ALTER OPENFLOW CONNECTOR c ADD LIVE VERSION FROM LAST COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c COMMIT COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c ABORT",
		"ALTER OPENFLOW CONNECTOR IF EXISTS c SET DEFAULT_VERSION = 'v3'",
		"ALTER OPENFLOW CONNECTOR c START",
		"ALTER OPENFLOW CONNECTOR c STOP",
		"ALTER OPENFLOW CONNECTOR c TERMINATE FORCE",
		"ALTER OPENFLOW CONNECTOR c SET DISPLAY_NAME = 'd' COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c RENAME TO d",
		"ALTER OPENFLOW CONNECTOR c PUSH TO 'https://g/b' USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e' COMMENT = 'c'",
		"ALTER OPENFLOW CONNECTOR c PUSH USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e'",
		"ALTER OPENFLOW CONNECTOR c PULL",
		"DROP OPENFLOW CONNECTOR IF EXISTS c",
		"CREATE EVENT ROUTING TABLE org_table WITH RULES default = (REGION_GROUP = 'PUBLIC', REGIONS = ('ALL'), DESTINATION_ACCOUNT = org.account1) aws_us = (REGION_GROUP = PUBLIC, REGIONS = ('AWS_US_EAST_1', 'AWS_US_WEST_2'), DESTINATION_ACCOUNT = org.account1)",
		"ALTER EVENT ROUTING TABLE t FORCE SET RULES a = (REGION_GROUP = PUBLIC, REGIONS = ('ALL'), DESTINATION_ACCOUNT = o.a), \"B\" = (REGION_GROUP = PUBLIC, REGIONS = ('X'), DESTINATION_ACCOUNT = o.b)",
		"ALTER EVENT ROUTING TABLE t SET RULE \"AWS_EU_WEST_RULE\" REGION_GROUP = PUBLIC REGIONS = ('AWS_EU_WEST_1') DESTINATION_ACCOUNT = org.account2",
		"ALTER EVENT ROUTING TABLE t FORCE UNSET RULE r",
		"ALTER EVENT ROUTING TABLE t MODIFY RULE r RENAME TO s",
		"ALTER EVENT ROUTING TABLE t RENAME TO u",
		"DROP EVENT ROUTING TABLE t",
	} {
		assertRoundTrip(t, sql)
	}
}

// TestDeparse_MiscObjectsText pins the text written for forms whose clauses
// the generic writer lays out: the data type, DEFAULT's expression, an
// index's column lists, DROP INDEX's table, and SET RULES' commas.
func TestDeparse_MiscObjectsText(t *testing.T) {
	for in, want := range map[string]string{
		"create variable s.v number(10,2) default 1+2": "CREATE VARIABLE s.v number(10, 2) DEFAULT 1 + 2",
		"create type t as int":                         "CREATE TYPE t AS int",
		"create index i on t(a,b) include(c)":          "CREATE INDEX i ON t (A, B) INCLUDE (C)",
		"drop index s.t.i":                             "DROP INDEX s.t.i",
		"undrop variable v":                            "UNDROP VARIABLE v",
		"alter event routing table t set rules a=(region_group=p, regions=('x'), destination_account=o.a), b=(region_group=p, regions=('y'), destination_account=o.b)": "ALTER EVENT ROUTING TABLE t SET RULES A = (REGION_GROUP = P, REGIONS = ('x'), DESTINATION_ACCOUNT = O.A), B = (REGION_GROUP = P, REGIONS = ('y'), DESTINATION_ACCOUNT = O.B)",
	} {
		f, err := parser.Parse(in)
		require.NoError(t, err, in)
		out, err := deparse.DeparseFile(f)
		require.NoError(t, err, in)
		require.Equal(t, want, out, in)
	}
}
