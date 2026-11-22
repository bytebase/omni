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

// TestSequenceSimulation holds the catalog to testdata/sequences/simulate.tsv:
// every generated sequence DDL statement TiDB accepts in the parser's
// compat.tsv, run on each fixture, with what TiDB did and what Exec does.
// Unlike TestSimulation (which reads only tables and views so existing row
// hashes stay stable), this suite reads tables, views, and sequences from TiDB
// (SHOW FULL TABLES + SHOW CREATE TABLE/VIEW/SEQUENCE) and compares every
// object's SHOW CREATE text against the catalog.
//
// TiDB's side is on testdata/tapes/TestSequenceSimulation.tape, so the test
// starts no server. OMNI_RECORD=1 runs every row on the pinned TiDB again;
// -update-compat commits the rows a catalog fix brought to TiDB's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
func TestSequenceSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), sequenceSimulationFixtures, sequenceSimulated)
	if err != nil {
		t.Fatal(err)
	}
	// Only the grammar's rows: compat.tsv's hand rows (suite:, issue:) are
	// corpora that grow on their own, and this file should not move with them.
	rows = slices.DeleteFunc(rows, func(g simulate.Generated) bool { return handLabel(g.Rule) })
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "sequences", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		Omni:     sequenceSimulationOmni,
		Engine:   sequenceSimulationEngine,
	})
}

// sequenceSimulationFixtures are the starting schemas in testdata/fixtures.
// A compat.tsv statement writes ident for every object name, so these five
// fixtures cover an absent name (one_table@missing), a base table named ident
// (one_table@table), a view named ident (objects@view), a default ascending
// sequence named ident (objects@sequence), and a descending cyclic sequence
// with custom bounds/cache/comment named ident (objects@sequence_options).
var sequenceSimulationFixtures = []string{
	"one_table@missing", "one_table@table", "objects@view",
	"objects@sequence", "objects@sequence_options",
}

// sequenceSimulationFamilies are the leading words of TiDB sequence DDL
// statements whose whole effect is in the database a row reads back.
var sequenceSimulationFamilies = []string{
	"CREATE SEQUENCE ", "ALTER SEQUENCE ", "DROP SEQUENCE ",
}

func sequenceSimulated(statement string) bool {
	return slices.ContainsFunc(sequenceSimulationFamilies, func(words string) bool {
		return strings.HasPrefix(statement, words)
	})
}

// sequenceSimulationEngine asks TiDB, or the tape.
func sequenceSimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	tt, ok := t.(*testing.T)
	if !ok {
		t.Fatalf("TestSequenceSimulation runs as a test")
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
		forms, err := sequenceSimulationForms(ctx, db)
		return simulate.Answer{Forms: forms}, err
	}
}

// sequenceSimulationForms reads the selected database's tables, views, and
// sequences as TiDB prints them, in deterministic name order.
func sequenceSimulationForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
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
	slices.SortFunc(forms, func(a, b simulate.Form) int {
		if c := strings.Compare(a.Name, b.Name); c != 0 {
			return c
		}
		return strings.Compare(a.Kind, b.Kind)
	})
	for i, f := range forms {
		show := fmt.Sprintf("SHOW CREATE %s `%s`", f.Kind, strings.ReplaceAll(f.Name, "`", "``"))
		row, err := db.QueryContext(ctx, show)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", show, err)
		}
		// SHOW CREATE TABLE and SHOW CREATE SEQUENCE answer two columns and
		// SHOW CREATE VIEW four; the text is the second of each.
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

// sequenceSimulationOmni runs the row in the catalog and compares every stored
// table, view, and sequence against TiDB's SHOW CREATE texts.
func sequenceSimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
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

	got := sequenceSimulationStored(c)
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

// sequenceSimulationStored is what SHOW FULL TABLES and SHOW CREATE read from
// the catalog's simulation database: its tables, views, and sequences in name
// order.
func sequenceSimulationStored(c *Catalog) []simulate.Form {
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
	for _, seq := range db.Sequences {
		forms = append(forms, simulate.Form{Kind: "SEQUENCE", Name: seq.Name, Text: c.ShowCreateSequence(db.Name, seq.Name)})
	}
	slices.SortFunc(forms, func(a, b simulate.Form) int {
		if c := strings.Compare(a.Name, b.Name); c != 0 {
			return c
		}
		return strings.Compare(a.Kind, b.Kind)
	})
	return forms
}
