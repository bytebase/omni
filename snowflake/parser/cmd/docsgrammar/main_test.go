package main

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/grammar/docs"
)

func form(t *testing.T, notation string) string {
	t.Helper()
	alts, err := docs.Form(dialect, notation)
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
		// Optional clauses with commas between: one or more, a comma between each two.
		{"ALTER T <name> SET [ A = <num> ], [ B = <num> ] [ C = <num> ], [ D ],", "'ALTER' 'T' ID 'SET' ('A' '=' NUMBER | 'B' '=' NUMBER | 'C' '=' NUMBER | 'D') (',' ('A' '=' NUMBER | 'B' '=' NUMBER | 'C' '=' NUMBER | 'D'))*"},
		{"ALTER VIEW <name> UNSET [ SECURE ] [ COMMENT = '<string_literal>' ]", "'ALTER' 'VIEW' ID 'UNSET' ('SECURE' | 'COMMENT') (',' ('SECURE' | 'COMMENT'))*"},
		{"ON_ERROR = { CONTINUE | SKIP_FILE_<num> | 'SKIP_FILE_<num>%' }", "'ON_ERROR' '=' ('CONTINUE' | 'SKIP_FILE_1' | STRING)"},
		{"( <col_name> <col_type> [ , <col_name> <col_type> [ DEFAULT <expr> ] [ ... ] ] )", "'(' ID DATA_TYPE (',' ID DATA_TYPE ('DEFAULT' EXPR)?)* ')'"},
		{"ALTER T <name> MODIFY [ ( ] COLUMN <c> DROP DEFAULT [ ) ]", "'ALTER' 'T' ID 'MODIFY' ('(' 'COLUMN' ID 'DROP' 'DEFAULT' ')' | 'COLUMN' ID 'DROP' 'DEFAULT')"},
		{"ALTER P <name> UNSET [ MIN_LENGTH ] [ MAX_LENGTH ] [ COMMENT ]", "'ALTER' 'P' ID 'UNSET' ('MIN_LENGTH' | 'MAX_LENGTH' | 'COMMENT') (',' ('MIN_LENGTH' | 'MAX_LENGTH' | 'COMMENT'))*"},
		{"IPS = ( [ '<ip>' ] [ , '<ip>' , ... ] )", "'IPS' '=' '(' (STRING (',' STRING)*)? ')'"},
		{"SET RULES <rule_name> = ( A = <a> ), ...", "'SET' 'RULES' ID '=' '(' 'A' '=' ID ')' (',' ID '=' '(' 'A' '=' ID ')')*"},
		{"SET REFRESH_INTERVAL_SECONDS = <value>", "'SET' 'REFRESH_INTERVAL_SECONDS' '=' NUMBER"},
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
		// { <target> SET X | UNSET X }: UNSET keeps the target.
		{"ALTER V <name> { MODIFY COLUMN <col> SET MASKING POLICY <p> | UNSET MASKING POLICY }", "'ALTER' 'V' ID ('MODIFY' 'COLUMN' ID 'SET' 'MASKING' 'POLICY' ID | 'MODIFY' 'COLUMN' ID 'UNSET' 'MASKING' 'POLICY')"},
		{"ALTER V <name> { SET COMMENT = '<c>' | UNSET COMMENT }", "'ALTER' 'V' ID ('SET' 'COMMENT' '=' STRING | 'UNSET' 'COMMENT')"},
		// A | inside SQL's parentheses chooses within them.
		{"[ ENCRYPTION = ( TYPE = 'A' | TYPE = 'B' ) ]", "('ENCRYPTION' '=' '(' ('TYPE' '=' '\\'A\\'' | 'TYPE' '=' '\\'B\\'') ')')?"},
		{"ALTER T [ ( ] <c> DROP DEFAULT [ ) ]", "'ALTER' 'T' ('(' ID 'DROP' 'DEFAULT' ')' | ID 'DROP' 'DEFAULT')"},
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
	alts, err := docs.Statement(dialect, "ALTER TASK [ IF EXISTS ] <name> RESUME | SUSPEND")
	if err != nil {
		t.Fatal(err)
	}
	want := "'ALTER' 'TASK' ('IF' 'EXISTS')? ID 'RESUME' | 'ALTER' 'TASK' ('IF' 'EXISTS')? ID 'SUSPEND'"
	if got := strings.Join(alts, " | "); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	alts, err = docs.Statement(dialect, "ALTER INTEGRATION <name> UNSET ENABLED | COMMENT")
	if err != nil {
		t.Fatal(err)
	}
	want = "'ALTER' 'INTEGRATION' ID 'UNSET' 'ENABLED' | 'ALTER' 'INTEGRATION' ID 'UNSET' 'COMMENT'"
	if got := strings.Join(alts, " | "); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	alts, err = docs.Statement(dialect, "ALTER P <name> SET A = <num> | B = <num>")
	if err != nil {
		t.Fatal(err)
	}
	want = "'ALTER' 'P' ID 'SET' 'A' '=' NUMBER | 'ALTER' 'P' ID 'SET' 'B' '=' NUMBER"
	if got := strings.Join(alts, " | "); got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	if _, err := docs.Statement(dialect, "SHOW TASKS | SUSPEND"); err != docs.ErrScope {
		t.Errorf("a choice with no name to scope it: got %v, want docs.ErrScope", err)
	}
}

func TestElidedTextIsNotConverted(t *testing.T) {
	// A lowercase word is a value beside = or |, and an undefined fragment
	// anywhere else.
	if _, err := docs.Form(dialect, "API_PROVIDER = { aws_api_gateway | azure }"); err != nil {
		t.Errorf("lowercase values: %v", err)
	}
	if _, err := docs.Form(dialect, "{ TYPE = app | TYPE = OAUTH2 oauth_parameters | TYPE = NONE }"); err == nil || !strings.Contains(err.Error(), "oauth_parameters") {
		t.Errorf("an undefined lowercase fragment: got %v", err)
	}
	if _, err := docs.Form(dialect, "( TYPE = OAUTH2 {oauth_parameters} )"); err == nil || !strings.Contains(err.Error(), "oauth_parameters") {
		t.Errorf("a lowercase word alone in braces: got %v", err)
	}
	if _, err := docs.Form(dialect, "CREATE S <name> [ RETENTION = ] [ COMMENT = '<s>' ]"); err == nil {
		t.Errorf("an = with no value was converted")
	}
	if _, err := docs.Form(dialect, "ALTER S <name> [ INCREMENT [ BY ] [ = ] <num> ]"); err != nil {
		t.Errorf("an optional = sign: %v", err)
	}
	if _, err := docs.Form(dialect, "CREATE [ OR REPLACE ] <integration_type> INTEGRATION <object_name>"); err == nil {
		t.Errorf("a placeholder for a type's keyword was converted")
	}
	if _, err := docs.Form(dialect, "ALTER [ <integration_type> ] INTEGRATION <object_name> <actions>"); err == nil {
		t.Errorf("a bracketed placeholder for a type's keyword was converted")
	}
	if got, err := docs.Form(dialect, "SHOW X IN APPLICATION {app_name}"); err != nil || got[0] != "'SHOW' 'X' 'IN' 'APPLICATION' ID" {
		t.Errorf("a placeholder in braces: got %v, %v", got, err)
	}
	for _, in := range []string{"SELECT * ILIKE ... REPLACE ...", "[ ... ] SELECT 1", "VALUES ( ... )"} {
		if _, err := docs.Form(dialect, in); err != docs.ErrElided {
			t.Errorf("%s: got %v, want docs.ErrElided", in, err)
		}
	}
	if _, err := docs.Form(dialect, "ALTER X { tableAction }"); err == nil || !strings.Contains(err.Error(), "tableAction") {
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
	pg := docs.ConvertPage(dialect, "alter-warehouse", blocks, nil)
	if len(pg.Failed) > 0 {
		t.Fatalf("not converted: %+v", pg.Failed)
	}
	want := []string{
		"'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'SET' (alterWarehouse_objectProperties alterWarehouse_objectParams? | alterWarehouse_objectParams)",
		"'ALTER' 'WAREHOUSE' ('IF' 'EXISTS')? ID 'RENAME' 'TO' ID",
	}
	if strings.Join(pg.Forms, "\n") != strings.Join(want, "\n") {
		t.Errorf("forms:\n%s\nwant:\n%s", strings.Join(pg.Forms, "\n"), strings.Join(want, "\n"))
	}
	frags := map[string]string{}
	for _, f := range pg.Frags {
		frags[f.Name] = strings.Join(f.Alts, " | ")
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
	pg := docs.ConvertPage(dialect, "create-catalog-integration-delta-sharing", []string{
		"CREATE CATALOG INTEGRATION <name> REST_AUTHENTICATION = ( restAuthenticationParams )",
		"restAuthenticationParams (for Bearer token) ::=\n  TYPE = BEARER\n  BEARER_TOKEN = '<bearer_token>'",
	}, nil)
	if len(pg.Frags) != 1 || strings.Join(pg.Frags[0].Alts, " | ") != "'TYPE' '=' 'BEARER' 'BEARER_TOKEN' '=' STRING" {
		t.Errorf("fragments %+v, failed %+v", pg.Frags, pg.Failed)
	}
}

// TestFormsStartAtTheStatementWord pins the split: COPY INTO's unindented
// clauses continue its form, and ALTER TABLE's indented forms each start
// one.
func TestFormsStartAtTheStatementWord(t *testing.T) {
	pg := docs.ConvertPage(dialect, "copy-into-table", []string{`/* Standard data load */
COPY INTO <t>
     FROM <s>
[ PATTERN = '<regex_pattern>' ]
[ VALIDATION_MODE = RETURN_ERRORS ]`}, nil)
	if len(pg.Forms) != 1 || len(pg.Failed) != 0 {
		t.Errorf("COPY INTO: forms %q, failed %+v", pg.Forms, pg.Failed)
	}
	pg = docs.ConvertPage(dialect, "alter-table", []string{`   ALTER TABLE [ IF EXISTS ] <name> RENAME TO <new_table_name>
   ALTER TABLE [ IF EXISTS ] <name> SWAP WITH <target_table_name>`}, nil)
	if len(pg.Forms) != 2 {
		t.Errorf("ALTER TABLE: forms %q", pg.Forms)
	}
}

// TestAPageMayShowAnotherStatement pins that a DDL verb at the start of a
// line starts a form of its own, on a page about another statement, and
// that a property named like one does not.
func TestAPageMayShowAnotherStatement(t *testing.T) {
	pg := docs.ConvertPage(dialect, "create-table-constraint", []string{`CREATE TABLE <name> ( <col_name> <col_type> )
COMMENT = '<string_literal>'
ALTER TABLE <name> ADD COLUMN <col_name> <col_type>`}, nil)
	want := []string{"'CREATE' 'TABLE' ID '(' ID DATA_TYPE ')' 'COMMENT' '=' STRING", "'ALTER' 'TABLE' ID 'ADD' 'COLUMN' ID DATA_TYPE"}
	if strings.Join(pg.Forms, "\n") != strings.Join(want, "\n") {
		t.Errorf("forms:\n%s", strings.Join(pg.Forms, "\n"))
	}
}

// TestABareValueChoiceIsTheValues pins a line KEY = A | B with no
// brackets: the choice is the value's, inside a statement form too.
func TestABareValueChoiceIsTheValues(t *testing.T) {
	pg := docs.ConvertPage(dialect, "create-security-integration", []string{"CREATE SECURITY INTEGRATION <name>\n  TYPE = OAUTH\n  OAUTH_CLIENT_TYPE = 'CONFIDENTIAL' | 'PUBLIC'\n  OAUTH_REDIRECT_URI = '<uri>'"}, nil)
	want := "'CREATE' 'SECURITY' 'INTEGRATION' ID 'TYPE' '=' 'OAUTH' 'OAUTH_CLIENT_TYPE' '=' ('\\'CONFIDENTIAL\\'' | '\\'PUBLIC\\'') 'OAUTH_REDIRECT_URI' '=' STRING"
	if len(pg.Forms) != 1 || pg.Forms[0] != want {
		t.Errorf("forms %q, failed %+v\nwant %s", pg.Forms, pg.Failed, want)
	}
}

// TestARepeatedFragmentKeepsItsOperator pins intoClause [ ... ]: the
// fragment repeats.
func TestARepeatedFragmentKeepsItsOperator(t *testing.T) {
	pg := docs.ConvertPage(dialect, "insert-multi-table", []string{"INSERT ALL intoClause [ ... ] <subquery>", "intoClause ::=\n  INTO <target_table>"}, nil)
	if len(pg.Forms) != 1 || pg.Forms[0] != "'INSERT' 'ALL' insertMultiTable_intoClause+ QUERY" {
		t.Errorf("forms %q, failed %+v", pg.Forms, pg.Failed)
	}
}

// TestATrailingBarIsTheLists pins a list of names where one is written
// with its value and the list's |: the | is not a choice of values.
func TestATrailingBarIsTheLists(t *testing.T) {
	pg := docs.ConvertPage(dialect, "alter-x", []string{"ALTER X <name> UNSET {\n    A |\n    B = <b> |\n    C }"}, nil)
	if len(pg.Forms) != 1 || pg.Forms[0] != "'ALTER' 'X' ID 'UNSET' ('A' | 'B' '=' ID | 'C')" {
		t.Errorf("forms %q, failed %+v", pg.Forms, pg.Failed)
	}
}

func TestFixesApplyToTheirPage(t *testing.T) {
	fixes := []*docs.Fix{{Page: "get", Find: "'<p>''", Replace: "'<p>'", Reason: "a quote too many"}}
	pg := docs.ConvertPage(dialect, "get", []string{"GET <s> [ PATTERN = '<p>'' ]"}, fixes)
	if !fixes[0].Used || len(pg.Forms) != 1 || len(pg.Fixes) != 1 {
		t.Errorf("forms %q, fixes %q, used %v", pg.Forms, pg.Fixes, fixes[0].Used)
	}
}

func TestPlaceholderClasses(t *testing.T) {
	for name, want := range map[string]string{
		"name": docs.ClassName, "table_name": docs.ClassName, "integer": docs.ClassNumber, "num": docs.ClassNumber, "rows": docs.ClassNumber,
		"expr": docs.ClassExpr, "boolean_expr": docs.ClassExpr, "condition": docs.ClassExpr, "query": docs.ClassQuery,
		"col_type": docs.ClassDataType, "arg_data_type": docs.ClassDataType, "string_literal": docs.ClassString, "copy_statement": tCopy,
		"sql":       docs.ClassQuery,
		"timestamp": docs.ClassString, "procedure_definition": docs.ClassString,
		"threshold": docs.ClassNumber,
		"col1_type": docs.ClassDataType, "specification_text": docs.ClassString,
		"python_version": docs.ClassString, "scala_or_java_runtime_version": docs.ClassString, "version_name": docs.ClassName,
		"object_type": tObjectType, "object_type_plural": tObjectTypes, "integration_name": docs.ClassName,
	} {
		if got := placeholderClass(name); got != want {
			t.Errorf("<%s> is %s, want %s", name, got, want)
		}
	}
}
