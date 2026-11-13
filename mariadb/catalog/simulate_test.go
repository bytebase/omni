package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mariadb/internal/testcontainer"
	"github.com/bytebase/omni/mariadb/parser"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: every table,
// index and view statement MariaDB accepts in the parser's compat.tsv, run
// on each fixture, with what MariaDB did and what Exec does. An outcome is
// the error number the statement is refused with, or ok: and a hash of the
// SHOW CREATE texts the database holds afterwards.
//
// MariaDB's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned MariaDB again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		// MariaDB makes up no name that differs between runs: an unnamed
		// foreign key is <table>_ibfk_<n>, and an unnamed CHECK
		// CONSTRAINT_<n>.
		MadeUp: nil,
		Omni:   simulationOmni,
		Engine: simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures: one
// table under four namings. A compat.tsv statement writes a for every name it
// holds, the table's, a column's and an index's alike, so what a statement
// finds is set by what the fixture names a: nothing, the table, a column of
// it too, and an index of it too.
var simulationFixtures = []string{"one_table@missing", "one_table@table", "one_table@column", "one_table@index"}

// simulatedStatement matches the statements the file holds: DDL on tables,
// indexes and views, whose whole effect is in the database a row reads back.
// A temporary table is in no SHOW FULL TABLES, so CREATE TEMPORARY is left
// out; so are databases, sequences, routines, triggers, events, accounts and
// servers, which a later fixture can take.
var simulatedStatement = regexp.MustCompile(`^(?:` +
	`ALTER (?:ONLINE |IGNORE )*TABLE |CREATE (?:OR REPLACE )?TABLE |DROP TABLE |RENAME TABLE |TRUNCATE ` +
	`|CREATE (?:OR REPLACE )?(?:UNIQUE |FULLTEXT |SPATIAL |VECTOR )?INDEX |DROP INDEX ` +
	`|(?:CREATE (?:OR REPLACE )?|ALTER )(?:ALGORITHM = \w+ )?(?:DEFINER = \S+ )?(?:SQL SECURITY \w+ )?VIEW |DROP VIEW ` +
	`)`)

func simulated(statement string) bool { return simulatedStatement.MatchString(statement) }

// simulationDatabase holds every row's fixture, on the server and in omni.
const simulationDatabase = "omni_simulation"

// simulationEngine asks MariaDB, or the tape. Each row is three kinds of
// request, none of which reads omni: the database dropped and built again
// with the fixture, the statement, and after one MariaDB took, the
// database's tables and views and the SHOW CREATE of each.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		_, err := db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+simulationDatabase)
		return err
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true})
	ctx := context.Background()
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		build := fmt.Sprintf("DROP DATABASE IF EXISTS %[1]s; CREATE DATABASE %[1]s; USE %[1]s; %[2]s", simulationDatabase, fx.SQL)
		if _, err := db.ExecContext(ctx, build); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var refusal *gomysql.MySQLError
			if !errors.As(err, &refusal) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: strconv.Itoa(int(refusal.Number))}, nil
		}
		forms, err := simulationForms(ctx, db)
		return simulate.Answer{Forms: forms}, err
	}
}

// simulationForms reads the fixture database's tables and views as MariaDB
// prints them, in name order. It names the database, since a statement may
// have selected another.
func simulationForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	rows, err := db.QueryContext(ctx, "SHOW FULL TABLES FROM "+simulationDatabase)
	if err != nil {
		return nil, err
	}
	var forms []simulate.Form
	for rows.Next() {
		var name, kind string
		if err := rows.Scan(&name, &kind); err != nil {
			rows.Close()
			return nil, err
		}
		f := simulate.Form{Kind: "TABLE", Name: name}
		if kind == "VIEW" {
			f.Kind = "VIEW"
		}
		forms = append(forms, f)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	slices.SortFunc(forms, func(a, b simulate.Form) int { return strings.Compare(a.Name, b.Name) })
	for i, f := range forms {
		show := fmt.Sprintf("SHOW CREATE %s %s.`%s`", f.Kind, simulationDatabase, strings.ReplaceAll(f.Name, "`", "``"))
		row, err := db.QueryContext(ctx, show)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", show, err)
		}
		// SHOW CREATE TABLE answers two columns and SHOW CREATE VIEW four;
		// the text is the second of either.
		columns, err := row.Columns()
		if err != nil {
			row.Close()
			return nil, err
		}
		values := make([]sql.NullString, len(columns))
		into := make([]any, len(columns))
		for k := range values {
			into[k] = &values[k]
		}
		if !row.Next() || len(values) < 2 {
			row.Close()
			return nil, fmt.Errorf("%s: no text", show)
		}
		if err := row.Scan(into...); err != nil {
			row.Close()
			return nil, err
		}
		forms[i].Text = values[1].String
		if err := row.Close(); err != nil {
			return nil, err
		}
	}
	return forms, nil
}

// simulationOmni runs the row in the catalog. The catalog stores what MariaDB
// stores when it prints the same tables and views as a catalog loaded from
// MariaDB's SHOW CREATE texts: both sides are printed by the catalog, so what
// it does not keep, such as a view's definer, does not count, and what it
// keeps and stores otherwise does.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	database := fmt.Sprintf("CREATE DATABASE %[1]s;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQL(database + fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *parser.ParseError
	switch {
	case errors.As(err, &syntax):
		return simulate.Omni{Code: "1064", Why: err.Error()}
	case err != nil:
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		var refusal *Error
		switch {
		case r.Error == nil:
		case errors.As(r.Error, &refusal) && refusal.Code > 0:
			return simulate.Omni{Code: strconv.Itoa(refusal.Code), Why: r.Error.Error()}
		default:
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}

	var stored strings.Builder
	stored.WriteString(database)
	for _, kind := range []string{"TABLE", "VIEW"} { // a view reads the tables
		for _, f := range engine.Forms {
			if f.Kind == kind {
				stored.WriteString(f.Text + ";\n")
			}
		}
	}
	want, err := LoadSQL(stored.String())
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MariaDB stores: %v", err)}
	}
	got, wanted := simulationPrint(c), simulationPrint(want)
	if got == wanted {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "the catalog stores\n" + got + "\nwhere MariaDB's stored form prints\n" + wanted}
}

// simulationPrint prints the fixture database's tables and views as the
// catalog shows them, in name order.
func simulationPrint(c *Catalog) string {
	db := c.GetDatabase(simulationDatabase)
	if db == nil {
		return "-- no database\n"
	}
	var b strings.Builder
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, t := range db.Tables {
			if !yield(t.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- TABLE %s\n%s\n", name, c.ShowCreateTable(simulationDatabase, name))
	}
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, v := range db.Views {
			if !yield(v.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- VIEW %s\n%s\n", name, c.ShowCreateView(simulationDatabase, name))
	}
	return b.String()
}
