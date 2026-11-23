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

	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// TestSimulation holds the catalog to testdata/simulate.tsv: every table,
// index and view statement MySQL accepts in the parser's compat.tsv, run on
// each fixture, with what MySQL did and what Exec does. An outcome is the
// error number the statement is refused with, or ok: and a hash of the SHOW
// CREATE texts the database holds afterwards.
//
// MySQL's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned MySQL again; -update-compat
// commits the rows a catalog fix brought to MySQL's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
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
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

// simulationFixtures are the starting schemas, in testdata/fixtures. They are
// one table under four namings plus three object-kind fixtures (a view, a
// partitioned table, and a table with a foreign key and check constraint). A
// compat.tsv statement writes utf8mb4 for every name it holds, the table's, a
// column's, an index's, a partition's, and a constraint's alike
// (mysql/testdata/tokens.tsv), so what a statement finds is set by what the
// fixture names utf8mb4.
var simulationFixtures = []string{
	"one_table@missing", "one_table@table", "one_table@column", "one_table@index",
	"objects@view", "objects@partitioned", "objects@fk",
}

// simulationFamilies are the leading words of the statements the file holds:
// DDL on tables, indexes and views, whose whole effect is in the database a
// row reads back. Statements on stored routines, triggers and events are
// simulated in TestProgramSimulation (testdata/programs/simulate.tsv) instead,
// where programSimulationForms reads routines, triggers and events in addition
// to SHOW FULL TABLES, so simulated excludes stored-program DDL even when it
// starts with a shared prefix like "CREATE DEFINER " or "ALTER DEFINER ".
// Statements on databases, tablespaces, accounts and the server are not
// simulated here yet: what they change is outside that database, or is no
// table or view text.
var simulationFamilies = []string{
	"ALTER TABLE ", "CREATE TABLE ", "DROP TABLE ", "RENAME TABLE ", "TRUNCATE ",
	"CREATE INDEX ", "CREATE UNIQUE ", "CREATE FULLTEXT ", "CREATE SPATIAL ", "DROP INDEX ",
	"CREATE VIEW ", "CREATE ALGORITHM ", "CREATE SQL ", "CREATE DEFINER ", "CREATE OR REPLACE ",
	"ALTER VIEW ", "ALTER ALGORITHM ", "ALTER SQL ", "ALTER DEFINER ", "DROP VIEW ",
}

func simulated(statement string) bool {
	if programSimulatedStatement.MatchString(statement) {
		return false
	}
	return slices.ContainsFunc(simulationFamilies, func(words string) bool { return strings.HasPrefix(statement, words) })
}

// simulationDatabase holds every row's fixture, on the server and in omni.
var simulationDatabase = testcontainer.DatabaseName("simulation")

// simulationEngine asks MySQL, or the tape. Each row is three kinds of
// request, none of which reads omni: the database dropped and built again with
// the fixture, the statement, and after one MySQL took, the database's tables
// and views and the SHOW CREATE of each.
func simulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	db := testcontainer.OpenEmptyWith(t, testcontainer.Options{MultiStatements: true})
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

// simulationForms reads the selected database's tables and views as MySQL
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

// parseCode is the code the server answers a statement it does not read with:
// the parser's own lookup refusal where it made one (1273 for a collation no
// server of the pin has), and 1064 for a syntax error.
func parseCode(pe *mysqlparser.ParseError) string {
	if pe.Code != 0 {
		return strconv.Itoa(pe.Code)
	}
	if pe.Message == "Invalid json attribute" {
		return "3980"
	}
	return "1064"
}

// simulationOmni runs the row in the catalog. The catalog stores what MySQL
// stores when a catalog loaded from MySQL's SHOW CREATE texts differs from it
// in neither direction: Diff under the 8.0 normalizer is what this package's
// oracle tests hold to MySQL's stored form, so a text the catalog prints
// otherwise but means the same is no difference.
func simulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	// The server's databases are utf8mb4 unless they say otherwise.
	database := fmt.Sprintf("CREATE DATABASE %[1]s DEFAULT CHARSET=utf8mb4;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQL(database + fx.SQL)
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
	want, err := LoadSQL(stored.String())
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MySQL stores: %v", err)}
	}
	n := NormalizerFor(MySQL80)
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

func TestSimulationExcludesProgramDDL(t *testing.T) {
	for _, stmt := range []string{
		"CREATE DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;",
		"ALTER DEFINER = utf8mb4 VIEW utf8mb4 AS SELECT utf8mb4;",
	} {
		if !simulated(stmt) {
			t.Errorf("simulated(%q) = false, want true", stmt)
		}
		if programSimulated(stmt) {
			t.Errorf("programSimulated(%q) = true, want false", stmt)
		}
	}
	for _, stmt := range []string{
		"CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;",
		"CREATE DEFINER = utf8mb4 FUNCTION utf8mb4() RETURNS INT DETERMINISTIC RETURN 1;",
		"CREATE DEFINER = utf8mb4 TRIGGER utf8mb4 BEFORE INSERT ON utf8mb4 FOR EACH ROW COMMIT;",
		"CREATE DEFINER = utf8mb4 EVENT utf8mb4 ON SCHEDULE AT utf8mb4 DO COMMIT;",
		"ALTER DEFINER = utf8mb4 EVENT utf8mb4 DISABLE;",
	} {
		if simulated(stmt) {
			t.Errorf("simulated(%q) = true, want false", stmt)
		}
		if !programSimulated(stmt) {
			t.Errorf("programSimulated(%q) = false, want true", stmt)
		}
	}

	f, err := compat.ReadFile(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	sawDefinerProc := false
	for _, r := range f.Rows {
		if r.Expected != compat.Accept {
			continue
		}
		if simulated(r.Statement) && programSimulated(r.Statement) {
			t.Errorf("compat.tsv row %s matches both simulated and programSimulated: %s", r.Rule, r.Statement)
		}
		if r.Statement == "CREATE DEFINER = utf8mb4 PROCEDURE utf8mb4() COMMIT;" && programSimulated(r.Statement) {
			sawDefinerProc = true
		}
	}
	if !sawDefinerProc {
		t.Errorf("expected CREATE DEFINER PROCEDURE row from compat.tsv to remain covered by programSimulated")
	}
}
