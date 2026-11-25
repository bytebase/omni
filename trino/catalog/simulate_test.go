package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	trinoclient "github.com/trinodb/trino-go-client/trino"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestSimulation holds Exec to testdata/simulate.tsv: every schema, table,
// view, materialized view, comment and TRUNCATE statement Trino accepts in
// the parser's compat.tsv, run on each fixture in a session on memory.s,
// with what Trino 482's memory connector did and what Exec does. An outcome
// is Trino's errorName, or ok: and a hash of what the memory catalog holds
// afterwards (simulationForms).
//
// Trino's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned Trino again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// catalogs of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulatedStatement.MatchString)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		// Trino makes up no name a form holds.
		MadeUp: nil,
		Omni:   simulationOmni,
		Engine: simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures. Each
// holds table t in the session's schema s, and a compat.tsv statement writes
// a for every name it holds, so what a statement finds is set by what the
// fixture names a: nothing but t's column, a table, a view, or a schema.
var simulationFixtures = []string{
	"one_table@missing", "one_table@table", "objects@view", "objects@schema",
}

// simulatedStatement matches the statements the file holds: DDL whose whole
// effect is in what the memory catalog's schemas, tables and views hold.
var simulatedStatement = regexp.MustCompile(`^(?:` +
	`CREATE (?:OR REPLACE )?(?:TABLE|VIEW|MATERIALIZED VIEW) |ALTER (?:TABLE|VIEW|MATERIALIZED VIEW) |DROP (?:TABLE|VIEW|MATERIALIZED VIEW) ` +
	`|CREATE SCHEMA |ALTER SCHEMA |DROP SCHEMA ` +
	`|COMMENT ON (?:TABLE|VIEW|COLUMN) |TRUNCATE TABLE ` +
	`)`)

// simulationSchema is the session's schema, which every fixture is built in.
const simulationSchema = "s"

// simulationEngine asks Trino, or the tape. Each row is a conversation that
// reads nothing of omni: the memory catalog emptied, schema s and the
// fixture built in it, the statement, and after one Trino took, what the
// catalog holds.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.OpenDB(t, "memory", simulationSchema)
	ctx := context.Background()
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		if err := emptyMemory(ctx, db); err != nil {
			return simulate.Answer{}, fmt.Errorf("empty the memory catalog: %w", err)
		}
		for _, s := range append([]string{"CREATE SCHEMA " + simulationSchema}, fixtureStatements(fx.SQL)...) {
			if _, err := db.ExecContext(ctx, s); err != nil {
				return simulate.Answer{}, fmt.Errorf("build fixture %s: %s: %w", fx.Name, s, err)
			}
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var refusal *trinoclient.ErrTrino
			if !errors.As(err, &refusal) || refusal.ErrorName == "" {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: refusal.ErrorName}, nil
		}
		forms, err := engineForms(ctx, db)
		if err != nil {
			return simulate.Answer{}, err
		}
		return simulate.Answer{Forms: forms}, nil
	}
}

// fixtureStatements splits a fixture into its statements, one to a line
// ending in ";".
func fixtureStatements(sql string) []string {
	var out []string
	for _, s := range strings.Split(sql, ";\n") {
		if s = strings.TrimSuffix(strings.TrimSpace(s), ";"); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// emptyMemory drops every schema of the memory catalog but default, with
// what it holds, and every table and view of default.
func emptyMemory(ctx context.Context, db *sql.DB) error {
	schemas, err := queryStrings(ctx, db, "SELECT schema_name FROM memory.information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'default') ORDER BY 1")
	if err != nil {
		return err
	}
	var drops []string
	for _, s := range schemas {
		drops = append(drops, "DROP SCHEMA memory."+quoteName(s[0])+" CASCADE")
	}
	rels, err := queryStrings(ctx, db, "SELECT table_name, table_type FROM memory.information_schema.tables WHERE table_schema = 'default' ORDER BY 1")
	if err != nil {
		return err
	}
	for _, r := range rels {
		kind := "TABLE"
		if r[1] == "VIEW" {
			kind = "VIEW"
		}
		drops = append(drops, "DROP "+kind+" memory.\"default\"."+quoteName(r[0]))
	}
	for _, d := range drops {
		if _, err := db.ExecContext(ctx, d); err != nil {
			return fmt.Errorf("%s: %w", d, err)
		}
	}
	return nil
}

// engineForms reads what the memory catalog holds, as simulationForms
// writes omni's: each schema, then each table as SHOW CREATE TABLE prints
// it, and each view as the head of SHOW CREATE VIEW (its name, comment and
// security) and its columns with their types. A view's query is left out:
// Trino stores it as its SqlFormatter writes it again, which omni has no
// copy of.
func engineForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	schemas, err := queryStrings(ctx, db, "SELECT schema_name FROM memory.information_schema.schemata WHERE schema_name <> 'information_schema' ORDER BY 1")
	if err != nil {
		return nil, err
	}
	var forms []simulate.Form
	for _, s := range schemas {
		forms = append(forms, simulate.Form{Kind: "SCHEMA", Name: s[0], Text: "CREATE SCHEMA memory." + quoteName(s[0])})
	}
	rels, err := queryStrings(ctx, db, "SELECT table_schema, table_name, table_type FROM memory.information_schema.tables WHERE table_schema <> 'information_schema' ORDER BY 1, 2")
	if err != nil {
		return nil, err
	}
	cols, err := queryStrings(ctx, db, "SELECT c.table_schema, c.table_name, c.column_name, c.data_type FROM memory.information_schema.columns c "+
		"JOIN memory.information_schema.tables t ON t.table_schema = c.table_schema AND t.table_name = c.table_name "+
		"WHERE t.table_type = 'VIEW' ORDER BY 1, 2, c.ordinal_position")
	if err != nil {
		return nil, err
	}
	viewCols := map[[2]string][]string{}
	for _, c := range cols {
		k := [2]string{c[0], c[1]}
		viewCols[k] = append(viewCols[k], "   "+quoteName(c[2])+" "+c[3])
	}
	for _, r := range rels {
		name := "memory." + quoteName(r[0]) + "." + quoteName(r[1])
		if r[2] != "VIEW" {
			text, err := queryStrings(ctx, db, "SHOW CREATE TABLE "+name)
			if err != nil {
				return nil, err
			}
			forms = append(forms, simulate.Form{Kind: "TABLE", Name: r[0] + "." + r[1], Text: text[0][0]})
			continue
		}
		text, err := queryStrings(ctx, db, "SHOW CREATE VIEW "+name)
		if err != nil {
			return nil, err
		}
		head, _, _ := strings.Cut(text[0][0], " AS\n")
		forms = append(forms, simulate.Form{Kind: "VIEW", Name: r[0] + "." + r[1], Text: head + "\n" + strings.Join(viewCols[[2]string{r[0], r[1]}], "\n")})
	}
	return forms, nil
}

// queryStrings runs a query and reads every column of every row as text.
func queryStrings(ctx context.Context, db *sql.DB, query string) ([][]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", query, err)
	}
	defer rows.Close()
	names, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out [][]string
	for rows.Next() {
		vals := make([]sql.NullString, len(names))
		ptrs := make([]any, len(names))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make([]string, len(vals))
		for i, v := range vals {
			row[i] = v.String
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

// simulationOmni runs the statement in a catalog Exec built the fixture in,
// on the session memory.s, and holds what it stores to Trino's forms.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c := New()
	db := c.EnsureCatalog("memory")
	db.EnsureSchema("default")
	db.EnsureSchema(simulationSchema)
	c.SetCurrentCatalog("memory")
	c.SetCurrentSchema(simulationSchema)
	if err := execAll(c, fx.SQL); err != nil {
		return simulate.Omni{Code: simulate.NoCode, Why: "the fixture fails: " + err.Error()}
	}
	if err := execAll(c, statement); err != nil {
		var e *Error
		if errors.As(err, &e) {
			return simulate.Omni{Code: e.Name, Why: e.Message}
		}
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	got := simulate.Answer{Forms: simulationForms(c)}.Text()
	if got == engine.Text() {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "omni stores:\n" + got}
}

// execAll runs sql in c and returns the first error.
func execAll(c *Catalog, sql string) error {
	results, err := c.Exec(sql, nil)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}

// simulationForms writes what c's memory catalog holds as engineForms reads
// Trino's.
func simulationForms(c *Catalog) []simulate.Form {
	db := c.GetCatalog("memory")
	var forms []simulate.Form
	for _, s := range db.Schemas() {
		forms = append(forms, simulate.Form{Kind: "SCHEMA", Name: s, Text: "CREATE SCHEMA memory." + quoteName(s)})
	}
	for _, sn := range db.Schemas() {
		s := db.GetSchema(sn)
		for _, name := range s.Relations() {
			if t := s.GetTable(name); t != nil {
				forms = append(forms, simulate.Form{Kind: "TABLE", Name: sn + "." + name, Text: t.ShowCreate("memory")})
				continue
			}
			v := s.GetView(name)
			text := v.ShowCreateHead("memory")
			for _, col := range v.Columns() {
				text += "\n   " + quoteName(col.Name) + " " + col.Type
			}
			forms = append(forms, simulate.Form{Kind: "VIEW", Name: sn + "." + name, Text: text})
		}
	}
	return forms
}
