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
// index, view, sequence, stored routine, trigger and event statement MariaDB
// accepts in the parser's compat.tsv, run on each fixture, with what MariaDB
// did and what Exec does. An outcome is the error number the statement is
// refused with, or ok: and a hash of the SHOW CREATE texts the database holds
// afterwards.
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

// simulationFixtures are the starting schemas, in testdata/fixtures. A
// compat.tsv statement writes a for every name it holds, so what a statement
// finds is set by what the fixture names a: nothing, a table, a column, an
// index, a view, a partition of a RANGE-partitioned table, a foreign key on
// a constrained table, a sequence, a stored procedure and function, a trigger,
// or an event.
var simulationFixtures = []string{
	"one_table@missing", "one_table@table", "one_table@column", "one_table@index",
	"objects@view", "objects@partitioned", "objects@fk",
	"objects@sequence", "objects@routine", "objects@trigger", "objects@event",
}

// simulatedStatement matches the statements the file holds: DDL on tables,
// indexes, views, sequences, stored routines, triggers and events, whose whole
// effect is in the database a row reads back. Temporary objects (CREATE
// TEMPORARY TABLE/SEQUENCE) are in no SHOW FULL TABLES, loadable UDFs (SONAME)
// live outside any database, and databases, packages, accounts and servers are
// left out.
var simulatedStatement = regexp.MustCompile(`^(?:` +
	`ALTER (?:ONLINE |IGNORE )*TABLE |CREATE (?:OR REPLACE )?TABLE |DROP TABLE |RENAME TABLE |TRUNCATE ` +
	`|CREATE (?:OR REPLACE )?(?:UNIQUE |FULLTEXT |SPATIAL |VECTOR )?INDEX |DROP INDEX ` +
	`|(?:CREATE (?:OR REPLACE )?|ALTER )(?:ALGORITHM = \w+ )?(?:DEFINER = \S+ )?(?:SQL SECURITY \w+ )?VIEW |DROP VIEW ` +
	`|(?:CREATE (?:OR REPLACE )?|ALTER |DROP )SEQUENCE ` +
	`|(?:CREATE (?:OR REPLACE )?(?:DEFINER = \S+ )?|ALTER |DROP )(?:PROCEDURE|FUNCTION) ` +
	`|(?:CREATE (?:OR REPLACE )?(?:DEFINER = \S+ )?|DROP )TRIGGER ` +
	`|(?:CREATE (?:OR REPLACE )?(?:DEFINER = \S+ )?|ALTER (?:DEFINER = \S+ )?|DROP )EVENT ` +
	`)`)

func simulated(statement string) bool {
	return simulatedStatement.MatchString(statement) && !strings.Contains(statement, " SONAME ")
}

// simulationDatabase holds every row's fixture, on the server and in omni.
const simulationDatabase = "omni_simulation"

// simulationEngine asks MariaDB, or the tape. Each row is three kinds of
// request, none of which reads omni: the database dropped and built again
// with the fixture, the statement, and after one MariaDB took, the
// database's objects and the SHOW CREATE of each.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	return simulationEngineWith(t, true)
}

func simulationEngineTablesOnly(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	return simulationEngineWith(t, false)
}

func simulationEngineWith(t testing.TB, includeRoutinesEventsTriggers bool) func(simulate.Fixture, string) (simulate.Answer, error) {
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
		forms, err := simulationForms(ctx, db, includeRoutinesEventsTriggers)
		return simulate.Answer{Forms: forms}, err
	}
}

// simulationForms reads the fixture database's tables, views, sequences, and
// (when includeExtra is true) routines, triggers and events as MariaDB prints
// them, in deterministic order. It names the database, since a statement may
// have selected another.
func simulationForms(ctx context.Context, db *sql.DB, includeExtra bool) ([]simulate.Form, error) {
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
		switch kind {
		case "VIEW":
			f.Kind = "VIEW"
		case "SEQUENCE":
			f.Kind = "SEQUENCE"
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

	if includeExtra {
		// Read routines, triggers and events in one query and append them in fixed
		// kind order (PROCEDURE, FUNCTION, TRIGGER, EVENT) and name order so
		// existing table/view fixtures keep the exact form order they were recorded with.
		const extraQuery = "" +
			"SELECT 'PROCEDURE', ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = '" + simulationDatabase + "' AND ROUTINE_TYPE = 'PROCEDURE' " +
			"UNION ALL " +
			"SELECT 'FUNCTION', ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = '" + simulationDatabase + "' AND ROUTINE_TYPE = 'FUNCTION' " +
			"UNION ALL " +
			"SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = '" + simulationDatabase + "' " +
			"UNION ALL " +
			"SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = '" + simulationDatabase + "'"
		extraRows, err := db.QueryContext(ctx, extraQuery)
		if err != nil {
			return nil, err
		}
		var extra []simulate.Form
		for extraRows.Next() {
			var kind, name string
			if err := extraRows.Scan(&kind, &name); err != nil {
				extraRows.Close()
				return nil, err
			}
			extra = append(extra, simulate.Form{Kind: kind, Name: name})
		}
		if err := extraRows.Close(); err != nil {
			return nil, err
		}
		if err := extraRows.Err(); err != nil {
			return nil, err
		}
		kindRank := map[string]int{"PROCEDURE": 1, "FUNCTION": 2, "TRIGGER": 3, "EVENT": 4}
		slices.SortFunc(extra, func(a, b simulate.Form) int {
			if c := kindRank[a.Kind] - kindRank[b.Kind]; c != 0 {
				return c
			}
			return strings.Compare(a.Name, b.Name)
		})
		forms = append(forms, extra...)
	}

	for i, f := range forms {
		show := fmt.Sprintf("SHOW CREATE %s %s.`%s`", f.Kind, simulationDatabase, strings.ReplaceAll(f.Name, "`", "``"))
		row, err := db.QueryContext(ctx, show)
		if err != nil {
			// After ALTER TABLE <seq> DISCARD TABLESPACE on an InnoDB sequence,
			// both SHOW CREATE SEQUENCE and SHOW CREATE TABLE fail with 1030
			// because opening a sequence reads its single data row from the
			// discarded .ibd tablespace. Record the error text deterministically.
			var me *gomysql.MySQLError
			if f.Kind == "SEQUENCE" && errors.As(err, &me) && me.Number == 1030 {
				forms[i].Text = fmt.Sprintf("-- ERROR %d: %s", me.Number, me.Message)
				continue
			}
			return nil, fmt.Errorf("%s: %w", show, err)
		}
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
		textCol := 1
		switch f.Kind {
		case "PROCEDURE", "FUNCTION", "TRIGGER":
			textCol = 2
		case "EVENT":
			textCol = 3
		}
		if !row.Next() || len(values) <= textCol {
			row.Close()
			return nil, fmt.Errorf("%s: no text", show)
		}
		if err := row.Scan(into...); err != nil {
			row.Close()
			return nil, err
		}
		forms[i].Text = values[textCol].String
		if err := row.Close(); err != nil {
			return nil, err
		}
	}
	return forms, nil
}

// simulationOmni runs the row in the catalog. The catalog stores what MariaDB
// stores when it prints the same objects as a catalog loaded from MariaDB's
// SHOW CREATE texts: both sides are printed by the catalog, so what it does
// not keep, such as a view's definer, does not count, and what it keeps and
// stores otherwise does.
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
		return simulate.Omni{Code: parseCode(syntax), Why: err.Error()}
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

	// engine.Forms are in name order, so a child table (e.g. a) can precede
	// its FK parent (e.g. parent). Disable FK checks while loading SHOW CREATE
	// forms from MariaDB, and load each object statement individually so
	// compound routine/trigger/event bodies need no delimiter rewriting.
	want, err := LoadSQL(database + "SET FOREIGN_KEY_CHECKS = 0;\n")
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not initialize want database: %v", err)}
	}
	for _, kind := range []string{"TABLE", "SEQUENCE", "VIEW", "PROCEDURE", "FUNCTION", "TRIGGER", "EVENT"} {
		for _, f := range engine.Forms {
			if f.Kind != kind {
				continue
			}
			// After ALTER TABLE <seq> DISCARD TABLESPACE, SHOW CREATE SEQUENCE
			// fails with 1030 on MariaDB; mirror that state on want directly.
			if f.Kind == "SEQUENCE" && strings.HasPrefix(f.Text, "-- ERROR 1030:") {
				wantDB := want.GetDatabase(simulationDatabase)
				wantDB.Sequences[toLower(f.Name)] = &Sequence{
					Name:                f.Name,
					Database:            wantDB,
					TablespaceDiscarded: true,
				}
				continue
			}
			res, err := want.Exec(f.Text, nil)
			if err != nil {
				return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MariaDB stores (%s %s): %v", f.Kind, f.Name, err)}
			}
			for _, r := range res {
				if r.Error != nil {
					return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MariaDB stores (%s %s): %v", f.Kind, f.Name, r.Error)}
				}
			}
		}
	}
	got, wanted := simulationPrint(c), simulationPrint(want)
	if got == wanted {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: "the catalog stores\n" + got + "\nwhere MariaDB's stored form prints\n" + wanted}
}

// parseCode is the outcome of a statement the parser refuses: the code the parser
// records for a name MariaDB looks up while it parses, as 1253 is for a collation of
// another character set and 1128 for a FULLTEXT parser the server does not ship, and
// 1064 for a syntax error.
func parseCode(pe *parser.ParseError) string {
	if pe.Code != 0 {
		return strconv.Itoa(pe.Code)
	}
	return "1064"
}

// simulationPrint prints the fixture database's objects as the catalog shows
// them, in kind and name order.
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
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, s := range db.Sequences {
			if !yield(s.Name) {
				return
			}
		}
	}) {
		if seq := db.Sequences[toLower(name)]; seq != nil && seq.TablespaceDiscarded {
			fmt.Fprintf(&b, "-- SEQUENCE %s\n-- ERROR 1030: Got error 194 \"Tablespace is missing for a table\" from storage engine InnoDB\n", name)
			continue
		}
		fmt.Fprintf(&b, "-- SEQUENCE %s\n%s\n", name, c.ShowCreateSequence(simulationDatabase, name))
	}
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, p := range db.Procedures {
			if !yield(p.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- PROCEDURE %s\n%s\n", name, c.ShowCreateProcedure(simulationDatabase, name))
	}
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, fn := range db.Functions {
			if !yield(fn.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- FUNCTION %s\n%s\n", name, c.ShowCreateFunction(simulationDatabase, name))
	}
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, tr := range db.Triggers {
			if !yield(tr.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- TRIGGER %s\n%s\n", name, c.ShowCreateTrigger(simulationDatabase, name))
	}
	for _, name := range slices.Sorted(func(yield func(string) bool) {
		for _, ev := range db.Events {
			if !yield(ev.Name) {
				return
			}
		}
	}) {
		fmt.Fprintf(&b, "-- EVENT %s\n%s\n", name, c.ShowCreateEvent(simulationDatabase, name))
	}
	return b.String()
}
