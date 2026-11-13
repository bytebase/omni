package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/tidb/internal/testcontainer"
	tidbparser "github.com/bytebase/omni/tidb/parser"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: every table,
// index and view statement TiDB accepts in the parser's compat.tsv, run on
// each fixture, with what TiDB did and what Exec does. An outcome is the
// error number the statement is refused with, or ok: and a hash of the SHOW
// CREATE texts the database holds afterwards.
//
// TiDB's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned TiDB again; -update-compat
// commits the rows a catalog fix brought to TiDB's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
func TestSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	// Only the grammar's rows: compat.tsv's hand rows (suite:, issue:) are
	// corpora that grow on their own, and this file should not move with them.
	// FromCompat keeps a statement's first rule, and hand rows follow the
	// generated ones, so a generated statement is never lost here.
	rows = slices.DeleteFunc(rows, func(g simulate.Generated) bool { return handLabel(g.Rule) })
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

// handLabel reports whether a compat.tsv rule is a label a person wrote
// (issue:BYT-1, suite:cited) rather than a grammar rule: a word and a colon
// before any slash.
func handLabel(rule string) bool {
	i := strings.IndexAny(rule, ":/")
	return i >= 0 && rule[i] == ':'
}

// simulationFixtures are the starting schemas, in testdata/fixtures. They
// are one table under four namings. A compat.tsv statement writes ident for
// every name it holds, the table's, a column's and an index's alike, so what
// a statement finds is set by what the fixture names ident: nothing, the
// table, a column of it too, and an index of it too.
var simulationFixtures = []string{"one_table@missing", "one_table@table", "one_table@column", "one_table@index"}

// simulationFamilies are the leading words of the statements the file holds:
// DDL on tables, indexes and views, whose whole effect is in the database a
// row reads back. Statements on databases, sequences, placement policies,
// resource groups, bindings, accounts and the server are not simulated here
// yet: what they change is outside that database.
var simulationFamilies = []string{
	"ALTER TABLE ", "ALTER IGNORE ", "CREATE TABLE ", "CREATE TEMPORARY ", "DROP TABLE ", "RENAME TABLE ", "TRUNCATE ",
	"CREATE INDEX ", "CREATE UNIQUE ", "CREATE FULLTEXT ", "CREATE SPATIAL ", "DROP INDEX ",
	"CREATE VIEW ", "CREATE OR ", "CREATE ALGORITHM ", "CREATE DEFINER ", "CREATE SQL ", "DROP VIEW ",
}

func simulated(statement string) bool {
	return slices.ContainsFunc(simulationFamilies, func(words string) bool { return strings.HasPrefix(statement, words) })
}

// simulationDatabase holds every row's fixture, on the server and in omni.
const simulationDatabase = "simulation"

// simulationEngine asks TiDB, or the tape. Each row is three kinds of
// request, none of which reads omni: the database dropped and built again
// with the fixture, the statement, and after one TiDB took, the database's
// tables and views and the SHOW CREATE of each.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	tt, ok := t.(*testing.T)
	if !ok {
		t.Fatalf("TestSimulation runs as a test")
	}
	ctr, done := startContainer(tt)
	t.Cleanup(done)
	db, ctx := ctr.db, ctr.ctx
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

// simulationForms reads the selected database's tables and views as TiDB
// prints them, in name order.
func simulationForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	rows, err := db.QueryContext(ctx, "SHOW FULL TABLES")
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
		show := fmt.Sprintf("SHOW CREATE %s `%s`", f.Kind, strings.ReplaceAll(f.Name, "`", "``"))
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

// simulationOmni runs the row in the catalog. The catalog stores what TiDB
// stores when it prints every table and view as TiDB's SHOW CREATE does, the
// test this package's differential tests hold it to: ShowCreateTable and
// ShowCreateView write TiDB's text.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	database := fmt.Sprintf("CREATE DATABASE %[1]s;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQL(database + fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *tidbparser.ParseError
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

	got := simulationStored(c)
	for i := 0; i < max(len(got), len(engine.Forms)); i++ {
		var mine, theirs simulate.Form
		if i < len(got) {
			mine = got[i]
		}
		if i < len(engine.Forms) {
			theirs = engine.Forms[i]
		}
		if mine != theirs {
			return simulate.Omni{Why: fmt.Sprintf("TiDB stores %s %s as %q, the catalog %s %s as %q",
				theirs.Kind, theirs.Name, oneLine(theirs.Text), mine.Kind, mine.Name, oneLine(mine.Text))}
		}
	}
	return simulate.Omni{Same: true}
}

// simulationStored is what SHOW FULL TABLES and SHOW CREATE would read from
// the catalog's simulation database: its tables and views in name order. A
// local temporary table is not listed, as TiDB lists none.
func simulationStored(c *Catalog) []simulate.Form {
	db := c.GetDatabase(simulationDatabase)
	if db == nil {
		return nil
	}
	var forms []simulate.Form
	for _, tbl := range db.Tables {
		if tbl.Temporary {
			continue
		}
		forms = append(forms, simulate.Form{Kind: "TABLE", Name: tbl.Name, Text: c.ShowCreateTable(db.Name, tbl.Name)})
	}
	for _, v := range db.Views {
		forms = append(forms, simulate.Form{Kind: "VIEW", Name: v.Name, Text: c.ShowCreateView(db.Name, v.Name)})
	}
	slices.SortFunc(forms, func(a, b simulate.Form) int { return strings.Compare(a.Name, b.Name) })
	return forms
}

func oneLine(text string) string { return strings.Join(strings.Fields(text), " ") }
