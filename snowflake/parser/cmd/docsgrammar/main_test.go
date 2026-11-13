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
		{"( <col_name>, <col_name>, ... )", "'(' ID ',' ID (',' ID)* ')'"},
		{"{ matchedClause | notMatchedClause } [ ... ]", ""}, // fragments: see TestFragments
		{"COMMENT = '<string_literal>'", "'COMMENT' '=' STRING"},
		{"TYPE = 'SNOWPARK-OPTIMIZED'", "'TYPE' '=' '\\'SNOWPARK-OPTIMIZED\\''"},
		{"FROM @[<namespace>.]<stage_name>[/<path>/]", "'FROM' STAGE"},
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
	} {
		if got := form(t, in); got != want {
			t.Errorf("%s\n got %s\nwant %s", in, got, want)
		}
	}
}

func TestElidedTextIsNotConverted(t *testing.T) {
	c := &converter{page: "p", defs: map[string]bool{}}
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
		"col_type": tType, "arg_data_type": tType, "string_literal": tString,
		"object_type": tObjectType, "object_type_plural": tObjectTypes, "integration_name": tName,
	} {
		if got := placeholderClass(name); got != want {
			t.Errorf("<%s> is %s, want %s", name, got, want)
		}
	}
}
