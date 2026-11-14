package parser

import (
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// VARIABLE, TYPE, INDEX, MCP SERVER, GATEWAY, OPENFLOW CONNECTOR and EVENT
// ROUTING TABLE, read as ObjectStmts by their docs pages.

func TestParseVariable(t *testing.T) {
	s := mustObjectStmt(t, "CREATE VARIABLE myschema.api_version VARCHAR DEFAULT 'v2' COMMENT = 'REST API version'")
	if s.Kind != "VARIABLE" || s.Name.String() != "myschema.api_version" {
		t.Fatalf("got %+v", s)
	}
	if got := clauseKeywords(s); len(got) != 3 || got[0] != "" || got[1] != "DEFAULT" || got[2] != "" {
		t.Fatalf("clauses = %q", got)
	}
	if s.Clauses[0].Type == nil || s.Clauses[0].Type.Name != "VARCHAR" || s.Clauses[1].Options[0].Expr == nil {
		t.Errorf("type/default = %+v %+v", s.Clauses[0], s.Clauses[1].Options[0])
	}
	for _, ok := range []string{
		"CREATE VARIABLE myschema.tax_rate FLOAT DEFAULT 0.08",
		"CREATE OR REPLACE VARIABLE a INT",
		"CREATE VARIABLE IF NOT EXISTS myschema.batch_cursor NUMBER DEFAULT 0",
		"CREATE VARIABLE a INT DEFAULT 1 + 2",
		"ALTER VARIABLE myschema.tax_rate RENAME TO myschema.regional_tax_rate",
		"ALTER VARIABLE IF EXISTS myschema.tax_rate SET DEFAULT 0.10",
		"ALTER VARIABLE myschema.batch_cursor DROP DEFAULT",
		"ALTER VARIABLE a SET COMMENT = 'Regional tax rate'",
		"ALTER VARIABLE a UNSET COMMENT",
	} {
		mustObjectStmt(t, ok)
	}
	for _, bad := range []string{
		"CREATE VARIABLE a",                                    // a data type is required
		"CREATE VARIABLE a INT COMMENT = 1",                    // a comment is a string
		"CREATE VARIABLE a INT COMMENT = 'x' DEFAULT 1",        // DEFAULT comes first
		"CREATE VARIABLE a INT DEFAULT",                        // DEFAULT takes a value
		"CREATE VARIABLE a INT TAG (t = 's')",                  // no tags documented
		"CREATE VARIABLE a INT SECURE = TRUE",                  // not a documented property
		"CREATE OR REPLACE VARIABLE IF NOT EXISTS a INT",       // IF NOT EXISTS only without OR REPLACE
		"CREATE OR ALTER VARIABLE a INT",                       // no OR ALTER
		"ALTER VARIABLE a SET TAG t = 's'",                     // no tags documented
		"ALTER VARIABLE a SET DEFAULT_VALUE = 1",               // not a documented property
		"ALTER VARIABLE a UNSET DEFAULT",                       // DROP DEFAULT clears it
		"ALTER VARIABLE a DROP",                                // DROP DEFAULT is whole
		"ALTER VARIABLE a",                                     // an action is required
		"ALTER VARIABLE a SET COMMENT = 'x' COMMENT = 'y'",     // set twice
		"ALTER VARIABLE a RENAME TO b SET COMMENT = 'x'",       // one action
		"UNDROP VARIABLE",                                      // a name is required
		"UNDROP VARIABLE IF EXISTS a",                          // the page writes no IF EXISTS
		"DROP VARIABLE a CASCADE",                              // no CASCADE documented
		"ALTER VARIABLE a SET DEFAULT 1 COMMENT = 'x'",         // one action
		"CREATE VARIABLE a INT DEFAULT 1 DEFAULT 2",            // DEFAULT once
		"CREATE VARIABLE a INT COMMENT = 'x' COMMENT = 'y'",    // set twice
		"CREATE VARIABLE a INT WITH TAG (t = 's')",             // no tags documented
		"CREATE VARIABLE a INT \"COMMENT\" = 'x'",              // a property's name is never quoted
		"CREATE \"VARIABLE\" a INT",                            // a quoted word is a name
		"ALTER VARIABLE a UNSET COMMENT, COMMENT_TOO",          // not a documented property
		"CREATE VARIABLE a INT DEFAULT 1 COMMENT = 'x' extra",  // trailing tokens
		"DROP VARIABLE IF EXISTS",                              // a name is required
		"CREATE VARIABLE db.s.a.b INT",                         // at most three parts
		"ALTER VARIABLE a SET",                                 // nothing to set
		"ALTER VARIABLE a UNSET",                               // nothing to unset
		"ALTER VARIABLE a RENAME TO",                           // a new name is required
		"ALTER VARIABLE a SET DEFAULT",                         // a value is required
		"CREATE VARIABLE a INT DEFAULT 1 COMMENT",              // a value is required
		"CREATE VARIABLE a INT COMMENT 'x'",                    // '=' is required
		"CREATE VARIABLE a INT DEFAULT 1 TAG (t = 's')",        // no tags documented
		"ALTER VARIABLE a UNSET TAG t",                         // no tags documented
		"CREATE VARIABLE a INT COMMENT = 'x' WITH TAG (t='s')", // no tags documented
	} {
		mustRejectObject(t, bad)
	}
	u, ok := mustParseOne(t, "UNDROP VARIABLE myschema.tax_rate").(*ast.UndropStmt)
	if !ok || u.Kind != ast.UndropObject || u.ObjectType() != "VARIABLE" || u.Name.String() != "myschema.tax_rate" {
		t.Errorf("UNDROP VARIABLE: got %+v", u)
	}
	d, ok := mustParseOne(t, "DROP VARIABLE IF EXISTS myschema.batch_cursor").(*ast.DropStmt)
	if !ok || d.ObjectType() != "VARIABLE" || !d.IfExists {
		t.Errorf("DROP VARIABLE: got %+v", d)
	}
}

func TestParseType(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE TYPE age AS NUMBER(3, 0) COMMENT = 'years'")
	if s.Kind != "TYPE" || !s.OrReplace || len(s.Clauses) != 2 || s.Clauses[0].Keywords != "AS" || s.Clauses[0].Type.Name != "NUMBER" {
		t.Fatalf("got %+v", s)
	}
	for _, ok := range []string{
		"CREATE TYPE IF NOT EXISTS db.s.age AS INT",
		"ALTER TYPE age SET COMMENT = 'User-defined type for storing age values'",
		"ALTER TYPE IF EXISTS age UNSET COMMENT",
	} {
		mustObjectStmt(t, ok)
	}
	for _, bad := range []string{
		"CREATE TYPE a INT",                    // AS is required
		"CREATE TYPE a AS",                     // a type is required
		"CREATE TYPE a AS INT COMMENT = 1",     // a comment is a string
		"CREATE TYPE a AS INT TAG (t = 's')",   // no tags documented
		"ALTER TYPE a RENAME TO b",             // not documented
		"ALTER TYPE a SET TAG t = 's'",         // no tags documented
		"ALTER TYPE a SET DATA_TYPE = INT",     // not a documented property
		"DROP TYPE IF EXISTS a",                // the page writes no IF EXISTS
		"UNDROP TYPE IF EXISTS a",              // the page writes no IF EXISTS
		"CREATE OR ALTER TYPE a AS INT",        // no OR ALTER
		"ALTER TYPE a SET COMMENT = 'x' x",     // trailing tokens
		"CREATE TYPE a AS INT COMMENT = 'x' x", // trailing tokens
	} {
		mustRejectObject(t, bad)
	}
	u, ok := mustParseOne(t, "UNDROP TYPE age").(*ast.UndropStmt)
	if !ok || u.ObjectType() != "TYPE" {
		t.Errorf("UNDROP TYPE: got %+v", u)
	}
	if d, ok := mustParseOne(t, "DROP TYPE age").(*ast.DropStmt); !ok || d.ObjectType() != "TYPE" {
		t.Errorf("DROP TYPE: got %+v", d)
	}
}

func TestParseIndex(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE INDEX vidx ON db.s.mytable (val, \"Val2\") INCLUDE (c)")
	if got := clauseKeywords(s); len(got) != 2 || got[0] != "ON" || got[1] != "INCLUDE" {
		t.Fatalf("clauses = %q", got)
	}
	on := s.Clauses[0]
	if on.Names[0].String() != "db.s.mytable" || len(on.Options) != 1 || len(on.Options[0].Group) != 2 ||
		on.Options[0].Group[0].Name != "VAL" || on.Options[0].Group[1].Name != `"Val2"` {
		t.Errorf("ON = %+v", on)
	}
	mustObjectStmt(t, "CREATE INDEX IF NOT EXISTS i ON t (a)")
	for _, bad := range []string{
		"CREATE INDEX i ON t",                              // columns are required
		"CREATE INDEX i ON t ()",                           // at least one column
		"CREATE INDEX i ON t (a,)",                         // trailing comma
		"CREATE INDEX i ON t (a.b)",                        // a column is one part
		"CREATE INDEX db.i ON t (a)",                       // the index name is one part
		"CREATE INDEX i ON t (a) INCLUDE ()",               // at least one column
		"CREATE INDEX i (a)",                               // ON is required
		"CREATE OR REPLACE INDEX IF NOT EXISTS i ON t (a)", // IF NOT EXISTS only without OR REPLACE
		"ALTER INDEX i RENAME TO j",                        // no ALTER INDEX page
		"CREATE INDEX i ON t (a) COMMENT = 'x'",            // not documented
		"DROP INDEX i",                                     // <table_name>.<index_name>
		"DROP INDEX a.b.c.d.e",                             // at most a three-part table
		"DROP INDEX t.i CASCADE",                           // not documented
		"DROP INDEX IF EXISTS",                             // a name is required
	} {
		mustRejectObject(t, bad)
	}
	d, ok := mustParseOne(t, "DROP INDEX IF EXISTS db.s.t.i").(*ast.DropStmt)
	if !ok || d.ObjectType() != "INDEX" || !d.IfExists || d.Table.String() != "db.s.t" || d.Name.String() != "i" {
		t.Fatalf("DROP INDEX: got %+v", d)
	}
	if d.Loc.Start != 0 || d.Loc.End != len("DROP INDEX IF EXISTS db.s.t.i") {
		t.Errorf("DROP INDEX Loc = %+v", d.Loc)
	}
	d = mustParseOne(t, "DROP INDEX t0.c_idx").(*ast.DropStmt)
	if d.Table.String() != "t0" || d.Name.String() != "c_idx" {
		t.Errorf("DROP INDEX t0.c_idx: got %+v", d)
	}
}

func TestParseMcpServerAndGateway(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OR REPLACE MCP SERVER my_mcp_server FROM SPECIFICATION $$tools: []$$")
	if s.Kind != "MCP SERVER" || len(s.Clauses) != 1 || s.Clauses[0].Keywords != "FROM SPECIFICATION" || s.Clauses[0].Options[0].Lit.Value != "tools: []" {
		t.Fatalf("got %+v", s.Clauses[0])
	}
	for _, ok := range []string{
		"CREATE MCP SERVER IF NOT EXISTS db.s.m FROM SPECIFICATION 'x'",
		"CREATE GATEWAY split_gateway FROM SPECIFICATION $$spec: {}$$",
		"CREATE OR REPLACE GATEWAY g FROM SPECIFICATION 'x'",
		"ALTER GATEWAY IF EXISTS g FROM SPECIFICATION 'x'",
	} {
		mustObjectStmt(t, ok)
	}
	for _, bad := range []string{
		"CREATE MCP SERVER m",                                   // FROM SPECIFICATION is required
		"CREATE MCP SERVER m FROM SPECIFICATION x",              // the specification is a string
		"CREATE MCP SERVER m FROM 'x'",                          // SPECIFICATION is required
		"ALTER MCP SERVER m FROM SPECIFICATION 'x'",             // no ALTER MCP SERVER page
		"CREATE GATEWAY g FROM SPECIFICATION 'x' COMMENT = 'c'", // not documented
		"ALTER GATEWAY g SET COMMENT = 'c'",                     // not documented
		"CREATE OR ALTER GATEWAY g FROM SPECIFICATION 'x'",      // no OR ALTER
		"DROP MCP SERVER m CASCADE",                             // not documented
	} {
		mustRejectObject(t, bad)
	}
	for _, in := range []string{"DROP MCP SERVER IF EXISTS my_mcp_server", "DROP GATEWAY g"} {
		if _, ok := mustParseOne(t, in).(*ast.DropStmt); !ok {
			t.Errorf("%q: not a DropStmt", in)
		}
	}
}

func TestParseOpenflowConnector(t *testing.T) {
	s := mustObjectStmt(t, "CREATE OPENFLOW CONNECTOR my_db.my_schema.c IN RUNTIME my_db.my_schema.r FROM DEFINITION OPENFLOW_POSTGRES_CDC DISPLAY_NAME = 'My CDC' COMMENT = 'x'")
	if got := clauseKeywords(s); len(got) != 3 || got[0] != "IN RUNTIME" || got[1] != "FROM DEFINITION" || got[2] != "" {
		t.Fatalf("clauses = %q", got)
	}
	for _, ok := range []string{
		"CREATE OPENFLOW CONNECTOR IF NOT EXISTS c IN RUNTIME r FROM '@my_git_repo/branches/main/' COMMENT = 'Created from Git stage'",
		"ALTER OPENFLOW CONNECTOR c ADD VERSION FROM '@stage/path/'",
		"ALTER OPENFLOW CONNECTOR c ADD LIVE VERSION FROM LAST COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c ADD LIVE VERSION FROM v1",
		"ALTER OPENFLOW CONNECTOR c COMMIT",
		"ALTER OPENFLOW CONNECTOR c COMMIT COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c ABORT",
		"ALTER OPENFLOW CONNECTOR IF EXISTS c SET DEFAULT_VERSION = FIRST",
		"ALTER OPENFLOW CONNECTOR c SET DEFAULT_VERSION = 'v3'",
		"ALTER OPENFLOW CONNECTOR IF EXISTS my_db.my_schema.c START",
		"ALTER OPENFLOW CONNECTOR c STOP",
		"ALTER OPENFLOW CONNECTOR c TERMINATE FORCE",
		"ALTER OPENFLOW CONNECTOR c SET DISPLAY_NAME = 'd' COMMENT = 'x'",
		"ALTER OPENFLOW CONNECTOR c RENAME TO d",
		"ALTER OPENFLOW CONNECTOR c PUSH TO 'https://g/b' USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e' COMMENT = 'c'",
		"ALTER OPENFLOW CONNECTOR c PUSH USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e'",
		"ALTER OPENFLOW CONNECTOR c PULL",
	} {
		mustObjectStmt(t, ok)
	}
	for _, bad := range []string{
		"CREATE OPENFLOW CONNECTOR c FROM DEFINITION d",                          // IN RUNTIME is required
		"CREATE OPENFLOW CONNECTOR c IN RUNTIME r",                               // FROM is required
		"CREATE OPENFLOW CONNECTOR c IN RUNTIME r FROM 's' DISPLAY_NAME = 'x'",   // the stage form takes COMMENT alone
		"CREATE OPENFLOW CONNECTOR c IN RUNTIME r FROM DEFINITION db.d",          // the definition is one part
		"CREATE OPENFLOW CONNECTOR c IN RUNTIME r FROM DEFINITION d COMMENT = 1", // a comment is a string
		"CREATE OPENFLOW CONNECTOR c IN RUNTIME r FROM DEFINITION d TAG (t = 's')",
		"CREATE OR REPLACE OPENFLOW CONNECTOR c IN RUNTIME r FROM 's'",             // no OR REPLACE
		"ALTER OPENFLOW CONNECTOR IF EXISTS c ABORT",                               // the form writes no IF EXISTS
		"ALTER OPENFLOW CONNECTOR IF EXISTS c PULL",                                // the form writes no IF EXISTS
		"ALTER OPENFLOW CONNECTOR c SET DEFAULT_VERSION = FIRST COMMENT = 'x'",     // DEFAULT_VERSION goes alone
		"ALTER OPENFLOW CONNECTOR c SET DEFAULT_VERSION = NEWEST",                  // FIRST, LAST or a string
		"ALTER OPENFLOW CONNECTOR c SET TAG t = 's'",                               // no tags documented
		"ALTER OPENFLOW CONNECTOR c UNSET COMMENT",                                 // no UNSET documented
		"ALTER OPENFLOW CONNECTOR c PUSH USERNAME = 'u' PASSWORD = 'p' NAME = 'n'", // EMAIL is required
		"ALTER OPENFLOW CONNECTOR c PUSH USERNAME = 'u' PASSWORD = 'p' NAME = 'n' EMAIL = 'e' BRANCH = 'b'",
		"ALTER OPENFLOW CONNECTOR c COMMIT DISPLAY_NAME = 'x'", // COMMIT takes a COMMENT alone
		"ALTER OPENFLOW CONNECTOR c TERMINATE NOW",             // FORCE alone
		"ALTER OPENFLOW CONNECTOR c ADD VERSION FROM LAST",     // a stage reference is a string
	} {
		mustRejectObject(t, bad)
	}
}

func TestParseEventRoutingTable(t *testing.T) {
	s := mustObjectStmt(t, "CREATE EVENT ROUTING TABLE org_table WITH RULES "+
		"default = (REGION_GROUP='PUBLIC', REGIONS=('ALL'), DESTINATION_ACCOUNT = org.account1) "+
		"aws_us = (REGION_GROUP=PUBLIC, REGIONS=('AWS_US_EAST_1', 'AWS_US_WEST_2'), DESTINATION_ACCOUNT = org.account1)")
	if len(s.Clauses) != 1 || s.Clauses[0].Keywords != "WITH RULES" || len(s.Clauses[0].Options) != 2 || s.Clauses[0].Options[0].Name != "DEFAULT" {
		t.Fatalf("got %+v", s.Clauses)
	}
	for _, ok := range []string{
		"ALTER EVENT ROUTING TABLE org_table FORCE SET RULES a = (REGION_GROUP=PUBLIC, REGIONS=('ALL'), DESTINATION_ACCOUNT = org.account1), b = (REGION_GROUP=PUBLIC, REGIONS=('X'), DESTINATION_ACCOUNT = org.a2)",
		"ALTER EVENT ROUTING TABLE org_table FORCE SET RULE \"AWS_EU_WEST_RULE\" REGION_GROUP=PUBLIC REGIONS=('AWS_EU_WEST_1') DESTINATION_ACCOUNT = org.account2",
		"ALTER EVENT ROUTING TABLE org_table FORCE UNSET RULE aws_eu_west_rule",
		"ALTER EVENT ROUTING TABLE org_table UNSET RULE r",
		"ALTER EVENT ROUTING TABLE org_table MODIFY RULE aws_us_east_rule RENAME TO rule_1",
		"ALTER EVENT ROUTING TABLE org_table RENAME TO org_table_2",
	} {
		mustObjectStmt(t, ok)
	}
	for _, bad := range []string{
		"CREATE EVENT ROUTING TABLE t",            // rules are required
		"CREATE EVENT ROUTING TABLE t WITH RULES", // at least one rule
		"CREATE OR REPLACE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a)",
		"CREATE EVENT ROUTING TABLE IF NOT EXISTS t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a)",
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, DESTINATION_ACCOUNT = o.a)",                       // REGIONS is required
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGIONS=('x'), REGION_GROUP=P, DESTINATION_ACCOUNT = o.a)",        // in order
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = a)",          // <organization>.<account_name>
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a.b)",      // two parts
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=(x), DESTINATION_ACCOUNT = o.a)",          // regions are strings
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a, X = 1)", // three settings
		"CREATE EVENT ROUTING TABLE t WITH RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a),",       // no comma between rules
		"ALTER EVENT ROUTING TABLE t SET RULES a = (REGION_GROUP=P, REGIONS=('x'), DESTINATION_ACCOUNT = o.a),",         // trailing comma
		"ALTER EVENT ROUTING TABLE t SET RULE r REGION_GROUP=P REGIONS=('x')",                                           // DESTINATION_ACCOUNT is required
		"ALTER EVENT ROUTING TABLE t SET RULE r REGION_GROUP=P DESTINATION_ACCOUNT = o.a REGIONS=('x')",                 // in order
		"ALTER EVENT ROUTING TABLE t FORCE RENAME TO u",                                                                 // FORCE goes with the rule actions
		"ALTER EVENT ROUTING TABLE t FORCE MODIFY RULE r RENAME TO s",                                                   // FORCE goes with the rule actions
		"ALTER EVENT ROUTING TABLE t MODIFY RULE r",                                                                     // RENAME TO is required
		"ALTER EVENT ROUTING TABLE IF EXISTS t RENAME TO u",                                                             // the page writes no IF EXISTS
		"ALTER EVENT ROUTING TABLE t UNSET RULE",                                                                        // a rule is required
		"DROP EVENT ROUTING TABLE IF EXISTS t",                                                                          // the page writes no IF EXISTS
	} {
		mustRejectObject(t, bad)
	}
	if d, ok := mustParseOne(t, "DROP EVENT ROUTING TABLE org_table").(*ast.DropStmt); !ok || d.ObjectType() != "EVENT ROUTING TABLE" {
		t.Errorf("DROP: got %+v", d)
	}
}

func TestStartsEventRuleName(t *testing.T) {
	for in, want := range map[string]bool{
		"r":       true,
		`"Rule"`:  true,
		"default": true,
		"'r'":     false,
		"(":       false,
		"1":       false,
	} {
		p := &Parser{lexer: NewLexer(in), input: in}
		p.advance()
		if got := p.startsEventRuleName(); got != want || p.cur.Loc.Start != 0 {
			t.Errorf("%q: got %v, want %v", in, got, want)
		}
	}
}

func TestParseShowConnectorAndRoutingScopes(t *testing.T) {
	s, ok := mustParseOne(t, "SHOW VERSIONS IN OPENFLOW CONNECTOR my_db.my_schema.c").(*ast.ShowStmt)
	if !ok || s.ScopeText != "OPENFLOW CONNECTOR" || s.ScopeName.String() != "my_db.my_schema.c" {
		t.Errorf("SHOW VERSIONS: got %+v", s)
	}
	s, ok = mustParseOne(t, "SHOW RULES IN EVENT ROUTING TABLE (org_table)").(*ast.ShowStmt)
	if !ok || s.ScopeText != "EVENT ROUTING TABLE" || s.ScopeName.String() != "org_table" {
		t.Errorf("SHOW RULES: got %+v", s)
	}
	for _, bad := range []string{
		"SHOW VERSIONS IN OPENFLOW CONNECTOR",
		"SHOW VERSIONS IN OPENFLOW CONNECTOR c LIMIT 1",
		"SHOW VERSIONS LIKE 'v%' IN OPENFLOW CONNECTOR c",
		"SHOW RULES IN EVENT ROUTING TABLE org_table",
		"SHOW RULES IN EVENT ROUTING TABLE (org_table",
		"SHOW RULES IN EVENT ROUTING TABLE (t) STARTS WITH 'a'",
	} {
		mustRejectObject(t, bad)
	}
}
