package main

import (
	"strings"
	"testing"
)

func form(t *testing.T, notation string) string {
	t.Helper()
	c := &converter{page: "p", defs: map[string]bool{}}
	alts, err := c.form(notation)
	if err != nil {
		t.Fatalf("%s: %v", notation, err)
	}
	return strings.Join(alts, " | ")
}

// TestNotationReadsAsANTLR pins how the docs' notation is read: [ ] is
// optional, { | } a choice, <placeholders> token classes, and each way the
// docs write ... a repetition.
func TestNotationReadsAsANTLR(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"ALTER WAREHOUSE [ IF EXISTS ] [ <name> ] { SUSPEND | RESUME [ IF SUSPENDED ] }",
			"'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID? ('SUSPEND' | 'RESUME' ('IF' 'SUSPENDED')?)"},
		{"UNSET TAG <tag_name> [ , <tag_name> ... ]", "'UNSET' 'TAG' ID (',' ID)*"},
		{"( <col_name> [ , ... ] )", "'(' ID (',' ID)* ')'"},
		{"( <expr> [ , <expr> , ... ] )", "'(' EXPR (',' EXPR)* ')'"},
		// A list element the docs spell twice is the element, however long.
		{"( <col_name>, <col_name>, ... )", "'(' ID (',' ID)* ')'"},
		{"RENAME ( <col_name> AS <col_alias>, <col_name> AS <col_alias>, ... )", "'RENAME' '(' ID 'AS' ID (',' ID 'AS' ID)* ')'"},
		{"PUT file://<path_to_file>/<filename> @<stage>", "'PUT' FILE STAGE"},
		{"{ matchedClause | notMatchedClause } [ ... ]", ""}, // fragments: see TestFragments
		{"COMMENT = '<string_literal>'", "'COMMENT' '=' STRING"},
		{"TYPE = 'SNOWPARK-OPTIMIZED'", "'TYPE' '=' '\\'SNOWPARK-OPTIMIZED\\''"},
		{"FROM @[<namespace>.]<stage_name>[/<path>/]", "'FROM' STAGE"},
		{"START = <initial_value>", "'START' '=' NUMBER"},
		{"AS SELECT <query>", "'AS' QUERY"},
		{"F ( [ <arg_name> <arg_data_type> [ , ... ] ] )", "'F' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')'"},
		{"UPDATE <t> SET <col_name> = <value> [ , ... ]", "'UPDATE' ID 'SET' ID '=' ID (',' ID '=' ID)*"},
		{"GRANT READ [ , WRITE ] [ , ... ]", "'GRANT' 'READ' (',' 'WRITE')? (',' 'READ' (',' 'WRITE')?)*"},
		{"VALUES ( <v> ) [ , ... ]", "'VALUES' '(' ID ')' (',' '(' ID ')')*"},
		// The element is the target with its optional alias.
		{"SELECT { <col_name> | <expr> } [ [ AS ] <col_alias> ] [ , ... ]", "'SELECT' (ID | EXPR) ('AS'? ID)? (',' (ID | EXPR) ('AS'? ID)?)*"},
		// [ X ] [ , ... ]: an optional list, a comma always before an X.
		{"F ( [ <arg_name> <arg_data_type> ] [ , ... ] )", "'F' '(' (ID DATA_TYPE (',' ID DATA_TYPE)*)? ')'"},
		// [ X , ... ] is an optional list, commas kept.
		{"DROP FUNCTION <name> ( [ <arg_data_type> , ... ] )", "'DROP' 'FUNCTION' ID '(' (DATA_TYPE (',' DATA_TYPE)*)? ')'"},
		// [ KEY = A | B ]: each choice is a value of KEY.
		{"SET [ SECRETS = ( <s> ) | all | none ]", "'SET' ('SECRETS' '=' '(' ID ')' | 'SECRETS' '=' 'all' | 'SECRETS' '=' 'none')"},
		{"DESC[RIBE] RESULT '<query_id>'", "('DESC' | 'DESCRIBE') 'RESULT' STRING"},
		{"STORAGE_PROVIDER = '{ S3 | S3CHINA }'", "'STORAGE_PROVIDER' '=' ('\\'S3\\'' | '\\'S3CHINA\\'')"},
		// A | inside SQL's parentheses chooses within them.
		{"[ ENCRYPTION = ( TYPE = 'A' | TYPE = 'B' ) ]", "('ENCRYPTION' '=' '(' ('TYPE' '=' '\\'A\\'' | 'TYPE' '=' '\\'B\\'') ')')?"},
		{"ALTER T [ ( ] <c> DROP DEFAULT [ ) ]", "'ALTER' 'T' ('(')? ID 'DROP' 'DEFAULT' (')')?"},
		{"SELECT $<col_position> FROM <t>", "'SELECT' POSITION 'FROM' ID"},
		{"API_PROVIDER = { aws_api_gateway | azure_api_management }", "'API_PROVIDER' '=' ('aws_api_gateway' | 'azure_api_management')"},
		{"SHOW TABLES IN { ACCOUNT | DATABASE | }", "'SHOW' 'TABLES' 'IN' ('ACCOUNT' | 'DATABASE')"},
		{"CREATE TABLE <t> AS <query>;", "'CREATE' 'TABLE' ID 'AS' QUERY"},
		{"ALTER TABLE <t> ALTER COLUMN <c> TYPE <col_type> -- a comment", "'ALTER' 'TABLE' ID 'ALTER' 'COLUMN' ID 'TYPE' DATA_TYPE"},
		{"AS $$ <yaml> $$", "'AS' STRING"},
	} {
		if c.want == "" {
			continue
		}
		if got := form(t, c.in); got != c.want {
			t.Errorf("%s\n got %s\nwant %s", c.in, got, c.want)
		}
	}
}

// TestSetSetsAtLeastOne pins the reading of SET or UNSET followed only by
// optional parts: at least one of them.
func TestSetSetsAtLeastOne(t *testing.T) {
	for in, want := range map[string]string{
		"ALTER X <name> SET [ A = <num> ] [ B = <num> ]":     "'ALTER' 'X' ID 'SET' ('A' '=' NUMBER ('B' '=' NUMBER)? | 'B' '=' NUMBER)",
		"ALTER X <name> SET { [ A = <num> ] [ B = <num> ] }": "'ALTER' 'X' ID 'SET' ('A' '=' NUMBER ('B' '=' NUMBER)? | 'B' '=' NUMBER)",
		"ALTER X <name> SET A = <num> [ B = <num> ]":         "'ALTER' 'X' ID 'SET' 'A' '=' NUMBER ('B' '=' NUMBER)?",
		"ALTER X SET { [ A = <num> ] | [ B = <num> ] }":      "'ALTER' 'X' 'SET' ('A' '=' NUMBER | 'B' '=' NUMBER)",
	} {
		if got := form(t, in); got != want {
			t.Errorf("%s\n got %s\nwant %s", in, got, want)
		}
	}
}

// TestBareChoicesChooseWhatFollowsTheName pins a | outside brackets in a
// statement form: a choice of what follows the object's name, not of
// whole statements.
func TestBareChoicesChooseWhatFollowsTheName(t *testing.T) {
	c := &converter{page: "p", defs: map[string]bool{}}
	alts, err := c.statement("ALTER TASK [ IF EXISTS ] <name> RESUME | SUSPEND")
	if err != nil {
		t.Fatal(err)
	}
	want := "'ALTER' 'TASK' ('IF' 'EXISTS')? ID 'RESUME' | 'ALTER' 'TASK' ('IF' 'EXISTS')? ID 'SUSPEND'"
	if got := strings.Join(alts, " | "); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	alts, err = c.statement("ALTER INTEGRATION <name> UNSET ENABLED | COMMENT")
	if err != nil {
		t.Fatal(err)
	}
	want = "'ALTER' 'INTEGRATION' ID 'UNSET' 'ENABLED' | 'ALTER' 'INTEGRATION' ID 'UNSET' 'COMMENT'"
	if got := strings.Join(alts, " | "); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	if _, err := c.statement("SHOW TASKS | SUSPEND"); err != errScope {
		t.Errorf("a choice with no name to scope it: got %v, want errScope", err)
	}
}

func TestElidedTextIsNotConverted(t *testing.T) {
	c := &converter{page: "p", defs: map[string]bool{}}
	// A lowercase word is a value beside = or |, and an undefined fragment
	// anywhere else.
	if _, err := c.form("API_PROVIDER = { aws_api_gateway | azure }"); err != nil {
		t.Errorf("lowercase values: %v", err)
	}
	if _, err := c.form("{ TYPE = app | TYPE = OAUTH2 oauth_parameters | TYPE = NONE }"); err == nil || !strings.Contains(err.Error(), "oauth_parameters") {
		t.Errorf("an undefined lowercase fragment: got %v", err)
	}
	if _, err := c.form("( TYPE = OAUTH2 {oauth_parameters} )"); err == nil || !strings.Contains(err.Error(), "oauth_parameters") {
		t.Errorf("a lowercase word alone in braces: got %v", err)
	}
	if _, err := c.form("CREATE S <name> [ RETENTION = ] [ COMMENT = '<s>' ]"); err == nil {
		t.Errorf("an = with no value was converted")
	}
	if _, err := c.form("ALTER S <name> [ INCREMENT [ BY ] [ = ] <num> ]"); err != nil {
		t.Errorf("an optional = sign: %v", err)
	}
	if got, err := c.form("SHOW X IN APPLICATION {app_name}"); err != nil || got[0] != "'SHOW' 'X' 'IN' 'APPLICATION' ID" {
		t.Errorf("a placeholder in braces: got %v, %v", got, err)
	}
	for _, in := range []string{"SELECT * ILIKE ... REPLACE ...", "[ ... ] SELECT 1", "VALUES ( ... )"} {
		if _, err := c.form(in); err != errElided {
			t.Errorf("%s: got %v, want errElided", in, err)
		}
	}
	if _, err := c.form("ALTER X { tableAction }"); err == nil || !strings.Contains(err.Error(), "tableAction") {
		t.Errorf("an undefined fragment: got %v", err)
	}
}

// TestPages pins how a page's blocks split into statement forms and
// definitions, and how definitions read.
func TestPages(t *testing.T) {
	blocks := []string{
		`ALTER WAREHOUSE [ IF EXISTS ] <name> SET [ objectProperties ]
                                         [ objectParams ]

ALTER WAREHOUSE [ IF EXISTS ] <name> RENAME TO <new_name>`,
		`objectProperties ::=
  WAREHOUSE_SIZE = XSMALL | SMALL
  AUTO_RESUME = { TRUE | FALSE }`,
		`objectParams ::=
-- For one
  { A | B } [ , ... ]
-- For another
  C`,
	}
	pg := convertPage("alter-warehouse", blocks, nil)
	if len(pg.failed) > 0 {
		t.Fatalf("not converted: %+v", pg.failed)
	}
	want := []string{
		"'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'SET' (alterWarehouse_objectProperties alterWarehouse_objectParams? | alterWarehouse_objectParams)",
		"'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID",
	}
	if strings.Join(pg.forms, "\n") != strings.Join(want, "\n") {
		t.Errorf("forms:\n%s\nwant:\n%s", strings.Join(pg.forms, "\n"), strings.Join(want, "\n"))
	}
	frags := map[string]string{}
	for _, f := range pg.frags {
		frags[f.name] = strings.Join(f.alts, " | ")
	}
	// One property per line: any number of them, each NAME = its choice.
	if got, want := frags["alterWarehouse_objectProperties"], "('WAREHOUSE_SIZE' '=' ('XSMALL' | 'SMALL') | 'AUTO_RESUME' '=' ('TRUE' | 'FALSE'))+"; got != want {
		t.Errorf("objectProperties:\n got %s\nwant %s", got, want)
	}
	// Variants between -- comments: one of them.
	if got, want := frags["alterWarehouse_objectParams"], "('A' | 'B') (',' ('A' | 'B'))* | 'C'"; got != want {
		t.Errorf("objectParams:\n got %s\nwant %s", got, want)
	}
}

// TestVariantParametersAreRequired pins that a definition written for one
// variant, name (for X) ::=, needs every line it lists, where an unnamed
// list of properties is any number of them.
func TestVariantParametersAreRequired(t *testing.T) {
	pg := convertPage("create-catalog-integration-delta-sharing", []string{
		"CREATE CATALOG INTEGRATION <name> REST_AUTHENTICATION = ( restAuthenticationParams )",
		"restAuthenticationParams (for Bearer token) ::=\n  TYPE = BEARER\n  BEARER_TOKEN = '<bearer_token>'",
	}, nil)
	if len(pg.frags) != 1 || strings.Join(pg.frags[0].alts, " | ") != "'TYPE' '=' 'BEARER' 'BEARER_TOKEN' '=' STRING" {
		t.Errorf("fragments %+v, failed %+v", pg.frags, pg.failed)
	}
}

// TestFormsStartAtTheStatementWord pins the split: COPY INTO's unindented
// clauses continue its form, and ALTER TABLE's indented forms each start
// one.
func TestFormsStartAtTheStatementWord(t *testing.T) {
	pg := convertPage("copy-into-table", []string{`/* Standard data load */
COPY INTO <t>
     FROM <s>
[ PATTERN = '<regex_pattern>' ]
[ VALIDATION_MODE = RETURN_ERRORS ]`}, nil)
	if len(pg.forms) != 1 || len(pg.failed) != 0 {
		t.Errorf("COPY INTO: forms %q, failed %+v", pg.forms, pg.failed)
	}
	pg = convertPage("alter-table", []string{`   ALTER TABLE [ IF EXISTS ] <name> RENAME TO <new_table_name>
   ALTER TABLE [ IF EXISTS ] <name> SWAP WITH <target_table_name>`}, nil)
	if len(pg.forms) != 2 {
		t.Errorf("ALTER TABLE: forms %q", pg.forms)
	}
}

// TestAPageMayShowAnotherStatement pins that a DDL verb at the start of a
// line starts a form of its own, on a page about another statement, and
// that a property named like one does not.
func TestAPageMayShowAnotherStatement(t *testing.T) {
	pg := convertPage("create-table-constraint", []string{`CREATE TABLE <name> ( <col_name> <col_type> )
COMMENT = '<string_literal>'
ALTER TABLE <name> ADD COLUMN <col_name> <col_type>`}, nil)
	want := []string{"'CREATE' 'TABLE' ID '(' ID DATA_TYPE ')' 'COMMENT' '=' STRING", "'ALTER' 'TABLE' ID 'ADD' 'COLUMN' ID DATA_TYPE"}
	if strings.Join(pg.forms, "\n") != strings.Join(want, "\n") {
		t.Errorf("forms:\n%s", strings.Join(pg.forms, "\n"))
	}
}

func TestFixesApplyToTheirPage(t *testing.T) {
	fixes := []*fix{{page: "get", find: "'<p>''", replace: "'<p>'", reason: "a quote too many"}}
	pg := convertPage("get", []string{"GET <s> [ PATTERN = '<p>'' ]"}, fixes)
	if !fixes[0].used || len(pg.forms) != 1 || len(pg.fixes) != 1 {
		t.Errorf("forms %q, fixes %q, used %v", pg.forms, pg.fixes, fixes[0].used)
	}
}

func TestPlaceholderClasses(t *testing.T) {
	for name, want := range map[string]string{
		"name": tName, "table_name": tName, "integer": tNumber, "num": tNumber, "rows": tNumber,
		"expr": tExpr, "boolean_expr": tExpr, "condition": tExpr, "query": tQuery,
		"col_type": tType, "arg_data_type": tType, "string_literal": tString, "copy_statement": tCopy,
		"sql":         tQuery,
		"object_type": tObjectType, "object_type_plural": tObjectTypes, "integration_name": tName,
	} {
		if got := placeholderClass(name); got != want {
			t.Errorf("<%s> is %s, want %s", name, got, want)
		}
	}
}
