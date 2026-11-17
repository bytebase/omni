package driver

import (
	"strings"
	"testing"
)

// definition is the shape the dump's transforms expect GET_DDL to state. It is
// written by hand from what the transforms read, not recorded from a database:
// no GET_DDL answer has been recorded, and VERIFY.md says how one is taken
// before this path is trusted.
const definition = `create or replace database DEMO_DB;

create or replace schema DEMO_DB.PUBLIC;

create or replace table DEMO_DB.PUBLIC.PARENT (
	ID NUMBER(38,0) NOT NULL,
	constraint PARENT_PK primary key (ID)
);
create or replace table DEMO_DB.PUBLIC.CHILD (
	ID NUMBER(38,0) NOT NULL,
	PARENT_ID NUMBER(38,0),
	constraint CHILD_FK foreign key (PARENT_ID) references DEMO_DB.PUBLIC.PARENT (ID)
);
create or replace view DEMO_DB.PUBLIC.RECENT as select ID from DEMO_DB.PUBLIC.CHILD;
`

// TestDatabaseDefinition pins the four passages a dump takes out of GET_DDL's
// answer: the database it opens with, the default PUBLIC schema, the database
// qualifier, and CREATE OR REPLACE.
func TestDatabaseDefinition(t *testing.T) {
	got := databaseDefinition("DEMO_DB", definition)
	for _, unwanted := range []string{
		"create or replace database",
		"create or replace schema",
		"create or replace",
	} {
		if strings.Contains(got, unwanted) {
			t.Errorf("the definition still holds %q:\n%s", unwanted, got)
		}
	}
	for _, wanted := range []string{
		"create table PUBLIC.PARENT (",
		"create table PUBLIC.CHILD (",
		"create view PUBLIC.RECENT as select ID from",
		"constraint PARENT_PK primary key (ID)",
	} {
		if !strings.Contains(got, wanted) {
			t.Errorf("the definition does not hold %q:\n%s", wanted, got)
		}
	}
	// Only a CREATE statement's own head is unqualified, as in Bytebase's
	// driver: a qualifier inside a statement's body stays. A recorded run
	// settles whether the server restores such a script into the same database.
	if !strings.Contains(got, "references DEMO_DB.PUBLIC.PARENT") {
		t.Errorf("the qualifier inside the foreign key was rewritten:\n%s", got)
	}
}

// TestDropCreateDatabase pins what a dump does with the statement GET_DDL opens
// with: it takes it and the blank line after it out, and leaves a definition
// that opens with something else whole.
func TestDropCreateDatabase(t *testing.T) {
	got := dropCreateDatabase("DEMO_DB", definition)
	if strings.Contains(got, "create or replace database") {
		t.Errorf("the CREATE DATABASE statement is still there:\n%s", got)
	}
	if !strings.HasPrefix(got, "create or replace schema") {
		t.Errorf("the definition does not start at the schema:\n%s", got)
	}

	const withoutDatabase = "create or replace table DEMO_DB.PUBLIC.T (ID NUMBER);\n"
	if got := dropCreateDatabase("DEMO_DB", withoutDatabase); got != withoutDatabase {
		t.Errorf("dropCreateDatabase = %q, want the definition unchanged", got)
	}
	// The header of a transient database states its modifier, and it is still
	// the database's own statement.
	for _, line := range []string{
		"create or replace database DEMO_DB;",
		"create transient database DEMO_DB;",
		"create or replace transient database DEMO_DB;\n",
	} {
		if got := dropCreateDatabase("DEMO_DB", line); got != "" {
			t.Errorf("dropCreateDatabase(%q) = %q, want nothing", line, got)
		}
	}
}

// TestDropPublicSchema pins that the default schema's CREATE is removed whether
// or not the server leaves a blank line after it.
func TestDropPublicSchema(t *testing.T) {
	for _, definition := range []string{
		"create or replace schema DEMO_DB.PUBLIC;\n\ncreate or replace table DEMO_DB.PUBLIC.T (ID NUMBER);",
		"create or replace schema DEMO_DB.PUBLIC;\ncreate or replace table DEMO_DB.PUBLIC.T (ID NUMBER);",
	} {
		got := dropPublicSchema("DEMO_DB", definition)
		if strings.Contains(got, "PUBLIC;") {
			t.Errorf("the PUBLIC schema's CREATE is still there:\n%s", got)
		}
		if !strings.Contains(got, "create or replace table") {
			t.Errorf("the table was removed too:\n%s", got)
		}
	}
	// The database is named as the quoted identifier its DDL requires.
	const quoted = `create or replace schema "demo_db".PUBLIC;` + "\n" + `create or replace table "demo_db".PUBLIC.T (ID NUMBER);`
	if got := dropPublicSchema("demo_db", quoted); strings.Contains(got, "PUBLIC;") {
		t.Errorf("the quoted PUBLIC schema's CREATE is still there:\n%s", got)
	}
	// Another database's PUBLIC schema is not this database's.
	const other = "create or replace schema OTHER_DB.PUBLIC;\n"
	if got := dropPublicSchema("DEMO_DB", other); got != other {
		t.Errorf("dropPublicSchema = %q, want the definition unchanged", got)
	}
}

// TestKeepCreate pins the two rewrites a dump's CREATE lines take, and that
// neither reaches past the statement's head.
func TestKeepCreate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "a non-CREATE line is left alone",
			input: "create or replace table DEMO_DB.PUBLIC.T (ID NUMBER);\n-- create or replace table DEMO_DB.PUBLIC.U (ID NUMBER);\n",
			want:  "create table PUBLIC.T (ID NUMBER);\n-- create or replace table DEMO_DB.PUBLIC.U (ID NUMBER);\n",
		},
		{
			name:  "a literal in a view's body keeps its text",
			input: "create or replace view DEMO_DB.PUBLIC.V as select 'create or replace table DEMO_DB.U' as s from DEMO_DB.PUBLIC.T;",
			want:  "create view PUBLIC.V as select 'create or replace table DEMO_DB.U' as s from DEMO_DB.PUBLIC.T;",
		},
		{
			name:  "an unqualified object name does not rewrite a literal that names the database",
			input: "create table PUBLIC.T (TEXT varchar default 'DEMO_DB.U');",
			want:  "create table PUBLIC.T (TEXT varchar default 'DEMO_DB.U');",
		},
		{
			name:  "a quoted object name loses the qualifier as written",
			input: `create or replace table "DEMO_DB".PUBLIC.T (ID NUMBER);`,
			want:  "create table PUBLIC.T (ID NUMBER);",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := keepCreate("DEMO_DB", tc.input); got != tc.want {
				t.Errorf("keepCreate = %q, want %q", got, tc.want)
			}
		})
	}
}
