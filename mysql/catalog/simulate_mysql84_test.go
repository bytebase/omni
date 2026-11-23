package catalog

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// TestSimulationMySQL84 holds the catalog to testdata/mysql84/simulate.tsv:
// the same table, index and view DDL statements and fixtures as TestSimulation,
// plus MySQL 8.4 LTS-specific hand-added cases, recorded against pinned
// MySQL 8.4 LTS (mysql84 pin in mysql/testdata/pins.tsv).
//
// MySQL's side is on the test's tape (testdata/tapes/TestSimulationMySQL84.tape),
// so the test starts no server. OMNI_RECORD=1 runs every row on the pinned
// MySQL 8.4 LTS server; -update-compat commits the rows a catalog fix brought
// to MySQL's outcome; and -simulate-show=<text> prints both schemas of the rows
// that hold the text.
func TestSimulationMySQL84(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), simulationFixtures, simulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "mysql84", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.MySQL84,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		Omni:     simulationOmniMySQL84,
		Engine:   simulationEngineMySQL84,
	})
}

func simulationEngineMySQL84(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.OpenEmptyWith(t, testcontainer.Options{Role: testcontainer.MySQL84, MultiStatements: true})
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

func simulationOmniMySQL84(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	database := fmt.Sprintf("CREATE DATABASE %[1]s DEFAULT CHARSET=utf8mb4;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQLWithVersion(database+fx.SQL, MySQL84)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *mysqlparser.ParseError
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

	var stored strings.Builder
	stored.WriteString(database)
	stored.WriteString("SET FOREIGN_KEY_CHECKS = 0;\n")
	for _, kind := range []string{"TABLE", "VIEW"} { // a view reads the tables
		for _, f := range engine.Forms {
			if f.Kind == kind {
				stored.WriteString(f.Text + ";\n")
			}
		}
	}
	stored.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n")
	want, err := LoadSQLWithVersion(stored.String(), MySQL84)
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MySQL stores: %v", err)}
	}
	n := NormalizerFor(MySQL84)
	there := DiffWithNormalizer(c, want, n)
	back := DiffWithNormalizer(want, c, n)
	if there.IsEmpty() && back.IsEmpty() {
		return simulate.Omni{Same: true}
	}
	plan := GenerateMigrationWithNormalizer(c, want, there, n).SQL()
	if there.IsEmpty() {
		plan = "nothing, though from there it would run: " + GenerateMigrationWithNormalizer(want, c, back, n).SQL()
	}
	return simulate.Omni{Why: "to store what MySQL stores, the catalog would run: " + strings.Join(strings.Fields(plan), " ")}
}
