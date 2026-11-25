package catalog

import (
	"strings"
	"testing"
)

// TestRunPLpgSQL pins what runPLpgSQL runs and what it leaves alone. The
// programs simulate suites hold the cases PostgreSQL was asked; these are
// the edges around them.
func TestRunPLpgSQL(t *testing.T) {
	const setup = `
CREATE TYPE mood AS ENUM ('happy');
CREATE SCHEMA app;
CREATE SCHEMA s2;
CREATE FUNCTION app.mt(mood, text) RETURNS boolean LANGUAGE sql AS 'SELECT true';
CREATE OPERATOR app.= (LEFTARG = mood, RIGHTARG = text, FUNCTION = app.mt);
CREATE TABLE t (id int);
CREATE TABLE flags (flag boolean);
CREATE PROCEDURE q(int) LANGUAGE plpgsql AS $$ BEGIN NULL; END $$;
CREATE PROCEDURE p(x int) LANGUAGE plpgsql AS $$ BEGIN CALL q(x); CREATE TABLE made_p (a int); END $$;
CREATE PROCEDURE qs() LANGUAGE plpgsql AS $$ DECLARE u int; BEGIN RAISE EXCEPTION 'stop'; END $$;
CREATE PROCEDURE ps() LANGUAGE plpgsql AS $$ BEGIN CALL qs(); CREATE TABLE made_ps (a int); END $$;
CREATE PROCEDURE pset() LANGUAGE plpgsql SET search_path = s2 AS $$ BEGIN CREATE TABLE made_set (a int); END $$;
CREATE PROCEDURE pdef() LANGUAGE plpgsql SECURITY DEFINER AS $$ BEGIN CREATE TABLE made_def (a int); END $$;
CREATE PROCEDURE pfail() LANGUAGE plpgsql AS $$ BEGIN CREATE TABLE made_fail (a int); INSERT INTO missing_t VALUES (1); END $$;
`
	for _, tc := range []struct {
		stmt    string
		code    string // "" when the statement succeeds
		made    string // a table the statement leaves in public
		notMade string // a table it leaves out
	}{
		// A procedure with parameters is not run: the catalog binds no
		// PL/pgSQL variable.
		{stmt: `CALL p(1)`, notMade: "made_p"},
		// Nor is one with SET options or SECURITY DEFINER, which the
		// catalog does not apply.
		{stmt: `CALL pset()`, notMade: "made_set"},
		{stmt: `CALL pdef()`, notMade: "made_def"},
		// A callee the catalog does not run stops the caller's run.
		{stmt: `CALL ps()`, notMade: "made_ps"},
		// A failed body leaves nothing behind.
		{stmt: `CALL pfail()`, code: "42P01", notMade: "made_fail"},
		{stmt: `DO $$ BEGIN CREATE INDEX CONCURRENTLY ix ON t (id); END $$`, code: "25001"},
		{stmt: `DO $$ BEGIN DO 'BEGIN CREATE TABLE inner_do (a int); END'; END $$`, made: "inner_do"},
		{stmt: `DO $$ BEGIN CREATE TABLE made_raise (a int); RAISE EXCEPTION 'x'; END $$`, code: "P0001", notMade: "made_raise"},
		// A statement that returns rows, or names FOUND, is not run.
		{stmt: `DO $$ BEGIN SHOW search_path; CREATE TABLE made_show (a int); END $$`, notMade: "made_show"},
		{stmt: `DO $$ BEGIN INSERT INTO flags VALUES (found); CREATE TABLE made_found (a int); END $$`, notMade: "made_found"},
		{stmt: `DO $$ BEGIN CREATE TABLE made_copy (a int); COPY made_copy TO STDOUT; END $$`, notMade: "made_copy"},
		// A user-defined operator takes the literal as its own type.
		{stmt: `CREATE FUNCTION fo(m mood) RETURNS boolean RETURN m OPERATOR(app.=) 'arbitrary'`},
		{stmt: `CREATE FUNCTION fe(m mood) RETURNS boolean RETURN m = 'arbitrary'`, code: "22P02"},
	} {
		c, err := LoadSQL(setup)
		if err != nil {
			t.Fatal(err)
		}
		res, err := c.Exec(tc.stmt, nil)
		if err != nil {
			t.Fatal(err)
		}
		code := ""
		if e, ok := res[0].Error.(*Error); ok {
			code = e.Code
		} else if res[0].Error != nil {
			t.Fatalf("%s: %v", tc.stmt, res[0].Error)
		}
		if code != tc.code {
			t.Errorf("%s: code %q, want %q (%v)", tc.stmt, code, tc.code, res[0].Error)
		}
		public := c.schemaByName["public"]
		if tc.made != "" && public.Relations[tc.made] == nil {
			t.Errorf("%s: no table %s", tc.stmt, tc.made)
		}
		if tc.notMade != "" && public.Relations[tc.notMade] != nil {
			t.Errorf("%s: table %s was made", tc.stmt, tc.notMade)
		}
	}
}

// TestLiteralTypeName pins that a literal's type prints as a cast to it
// does: quoted, and qualified outside public.
func TestLiteralTypeName(t *testing.T) {
	c, err := LoadSQL(`CREATE TYPE "Mood" AS ENUM ('sad'); CREATE SCHEMA app; CREATE TYPE app.e AS ENUM ('x');`)
	if err != nil {
		t.Fatal(err)
	}
	d := &ruleDeparser{cat: c}
	for name, want := range map[string]string{"public.Mood": `"Mood"`, "app.e": "app.e"} {
		schema, typ, _ := strings.Cut(name, ".")
		oid := c.typeByName[typeKey{ns: c.schemaByName[schema].OID, name: typ}]
		if oid == nil {
			t.Fatalf("no type %s", name)
		}
		if got := d.literalTypeName(oid.OID); got != want {
			t.Errorf("%s: %s, want %s", name, got, want)
		}
	}
	if got := d.literalTypeName(TEXTOID); got != "text" {
		t.Errorf("text: %s", got)
	}
}
