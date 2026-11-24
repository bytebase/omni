package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/mysql/internal/testcontainer"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// TestProgramSimulation holds the catalog to testdata/programs/simulate.tsv:
// every stored-program DDL statement (procedures, functions, triggers, events)
// MySQL accepts in the parser's compat.tsv, run across object fixtures, with
// what MySQL did and what Exec does. Unlike TestSimulation (which reads only
// SHOW FULL TABLES so existing table/view row hashes remain stable), this suite
// reads tables, views, stored procedures, stored functions, triggers, and events
// from MySQL and diffs all of them against the catalog.
//
// MySQL's side is on the test's tape (testdata/tapes/TestProgramSimulation.tape),
// so the test starts no server. OMNI_RECORD=1 runs every row on the pinned
// MySQL again; -update-compat commits the rows a catalog fix brought to
// MySQL's outcome; and -simulate-show=<text> prints both schemas of the rows
// that hold the text.
func TestProgramSimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), programSimulationFixtures, programSimulated)
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     append(rows, eventRows...),
		Omni:     programSimulationOmni,
		Engine:   programSimulationEngine,
	})
}

var programSimulationFixtures = []string{
	"one_table@missing", "one_table@table",
	"objects@routine", "objects@trigger", "objects@event",
}

// programSimulatedStatement matches stored-program DDL whose whole effect is in
// the database a row reads back: stored procedures, stored functions, triggers,
// and scheduled events. Loadable UDFs (CREATE [AGGREGATE] FUNCTION ... SONAME)
// are excluded because they install server-level shared libraries in mysql.func
// outside the database schema (like INSTALL PLUGIN).
var programSimulatedStatement = regexp.MustCompile(`^(?:` +
	`CREATE (?:DEFINER = \S+ )?(?:PROCEDURE|FUNCTION|TRIGGER|EVENT) ` +
	`|ALTER (?:DEFINER = \S+ )?(?:PROCEDURE|FUNCTION|EVENT) ` +
	`|DROP (?:PROCEDURE|FUNCTION|TRIGGER|EVENT) ` +
	`)`)

func programSimulated(statement string) bool {
	if strings.Contains(statement, " SONAME ") {
		return false
	}
	return programSimulatedStatement.MatchString(statement)
}

func programSimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	return programSimulationEngineOn(testcontainer.Default)(t)
}

// programSimulationEngineOn is programSimulationEngine on the server of a
// pins.tsv role.
func programSimulationEngineOn(role string) func(testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	return func(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
		return simulationEngineWith(t, role, programSimulationForms)
	}
}

// programSimulationForms reads tables, views, stored functions, stored
// procedures, triggers, and events from the selected database in deterministic
// order.
func programSimulationForms(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	forms, err := simulationForms(ctx, db)
	if err != nil {
		return nil, err
	}

	// Stored functions and procedures.
	routines, err := queryKindNames(ctx, db,
		"SELECT ROUTINE_TYPE, ROUTINE_NAME FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() ORDER BY ROUTINE_TYPE, ROUTINE_NAME")
	if err != nil {
		return nil, err
	}
	for _, f := range routines {
		text, err := showCreateText(ctx, db, f.Kind, f.Name, 2)
		if err != nil {
			return nil, err
		}
		f.Text = text
		forms = append(forms, f)
	}

	// Triggers.
	triggers, err := queryKindNames(ctx, db,
		"SELECT 'TRIGGER', TRIGGER_NAME FROM information_schema.TRIGGERS WHERE TRIGGER_SCHEMA = DATABASE() ORDER BY TRIGGER_NAME")
	if err != nil {
		return nil, err
	}
	for _, f := range triggers {
		text, err := showCreateText(ctx, db, f.Kind, f.Name, 2)
		if err != nil {
			return nil, err
		}
		f.Text = text
		forms = append(forms, f)
	}

	// Events.
	events, err := queryKindNames(ctx, db,
		"SELECT 'EVENT', EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME")
	if err != nil {
		return nil, err
	}
	for _, f := range events {
		text, err := showCreateText(ctx, db, f.Kind, f.Name, 3)
		if err != nil {
			return nil, err
		}
		f.Text = text
		forms = append(forms, f)
	}

	return forms, nil
}

func queryKindNames(ctx context.Context, db *sql.DB, query string) ([]simulate.Form, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var forms []simulate.Form
	for rows.Next() {
		var kind, name string
		if err := rows.Scan(&kind, &name); err != nil {
			rows.Close()
			return nil, err
		}
		forms = append(forms, simulate.Form{Kind: kind, Name: name})
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return forms, nil
}

func showCreateText(ctx context.Context, db *sql.DB, kind, name string, textCol int) (string, error) {
	show := fmt.Sprintf("SHOW CREATE %s `%s`", kind, strings.ReplaceAll(name, "`", "``"))
	row, err := db.QueryContext(ctx, show)
	if err != nil {
		return "", fmt.Errorf("%s: %w", show, err)
	}
	columns, err := row.Columns()
	if err != nil {
		row.Close()
		return "", err
	}
	values := make([]sql.NullString, len(columns))
	into := make([]any, len(columns))
	for k := range values {
		into[k] = &values[k]
	}
	if !row.Next() || len(values) <= textCol {
		row.Close()
		return "", fmt.Errorf("%s: no text at column %d", show, textCol)
	}
	if err := row.Scan(into...); err != nil {
		row.Close()
		return "", err
	}
	text := values[textCol].String
	if err := row.Close(); err != nil {
		return "", err
	}
	return text, nil
}

func programSimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	return programSimulationOmniAt(MySQL80)(fx, statement, engine)
}

// programSimulationOmniAt is programSimulationOmni with the catalog and the
// normalizer at version v.
func programSimulationOmniAt(v Version) func(simulate.Fixture, string, simulate.Answer) simulate.Omni {
	return func(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
		return programSimulationOmniVersion(v, fx, statement, engine)
	}
}

func programSimulationOmniVersion(v Version, fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	database := fmt.Sprintf("CREATE DATABASE %[1]s DEFAULT CHARSET=utf8mb4;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQLWithVersion(database+fx.SQL, v)
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

	want, err := LoadSQLWithVersion(database+"SET FOREIGN_KEY_CHECKS = 0;\n", v)
	if err != nil {
		return simulate.Omni{Why: fmt.Sprintf("the catalog does not initialize want: %v", err)}
	}
	for _, kind := range []string{"TABLE", "FUNCTION", "PROCEDURE", "VIEW", "TRIGGER", "EVENT"} {
		for _, f := range engine.Forms {
			if f.Kind != kind {
				continue
			}
			res, execErr := want.Exec(f.Text, nil)
			if execErr != nil {
				return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MySQL stores (%s %s): %v", f.Kind, f.Name, execErr)}
			}
			for _, r := range res {
				if r.Error != nil {
					return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what MySQL stores (%s %s): %v", f.Kind, f.Name, r.Error)}
				}
			}
		}
	}
	if _, err := want.Exec("SET FOREIGN_KEY_CHECKS = 1", nil); err != nil {
		return simulate.Omni{Why: fmt.Sprintf("restore FK checks: %v", err)}
	}

	n := NormalizerFor(v)
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

func TestProgramSimulationBugFixes(t *testing.T) {
	t.Run("ALTER EVENT RENAME TO same name returns 1551", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			setup string
			stmt  string
			want  int
		}{
			{
				name:  "missing event unqualified",
				setup: "CREATE DATABASE a; USE a;",
				stmt:  "ALTER EVENT e RENAME TO e",
				want:  ErrEventSameName,
			},
			{
				name:  "existing event unqualified",
				setup: "CREATE DATABASE a; USE a; CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO COMMIT;",
				stmt:  "ALTER EVENT e RENAME TO e",
				want:  ErrEventSameName,
			},
			{
				name:  "existing event case-insensitive",
				setup: "CREATE DATABASE a; USE a; CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO COMMIT;",
				stmt:  "ALTER EVENT a.e RENAME TO E",
				want:  ErrEventSameName,
			},
			{
				name:  "missing event qualified same schema",
				setup: "CREATE DATABASE a; CREATE DATABASE b; USE a;",
				stmt:  "ALTER EVENT b.e RENAME TO B.E",
				want:  ErrEventSameName,
			},
			{
				name:  "rename to nonexistent target db still 1049",
				setup: "CREATE DATABASE a; USE a; CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO COMMIT;",
				stmt:  "ALTER EVENT e RENAME TO nosuch.e",
				want:  ErrUnknownDatabase,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				c := New()
				mustExec(t, c, tc.setup)
				if got := execErrorCode(t, c, tc.stmt); got != tc.want {
					t.Errorf("%s: got %d, want %d", tc.stmt, got, tc.want)
				}
			})
		}

		// Cross-database rename of b.e to unqualified e (moves to current db a) succeeds.
		c := New()
		mustExec(t, c, "CREATE DATABASE a; CREATE DATABASE b; USE a; CREATE EVENT b.e ON SCHEDULE EVERY 1 DAY DO COMMIT;")
		mustExec(t, c, "ALTER EVENT b.e RENAME TO e")
		if c.GetDatabase("a").Events["e"] == nil || c.GetDatabase("b").Events["e"] != nil {
			t.Fatalf("expected event moved from b.e to a.e")
		}
	})

	t.Run("qualified DROP/ALTER routine in nonexistent schema returns 1305", func(t *testing.T) {
		for _, tc := range []struct {
			stmt string
			want int
		}{
			{"DROP FUNCTION nosuch.f", ErrNoSuchFunction},
			{"DROP PROCEDURE nosuch.p", ErrNoSuchProcedure},
			{"ALTER FUNCTION nosuch.f COMMENT 'x'", ErrNoSuchFunction},
			{"ALTER PROCEDURE nosuch.p COMMENT 'x'", ErrNoSuchProcedure},
		} {
			c := New()
			mustExec(t, c, "CREATE DATABASE a; USE a;")
			if got := execErrorCode(t, c, tc.stmt); got != tc.want {
				t.Errorf("%s: got %d, want %d", tc.stmt, got, tc.want)
			}
		}

		c := New()
		mustExec(t, c, "CREATE DATABASE a; USE a;")
		mustExec(t, c, "DROP FUNCTION IF EXISTS nosuch.f; DROP PROCEDURE IF EXISTS nosuch.p;")
	})

	t.Run("LoadSQL parses SHOW CREATE PROCEDURE with DEFINER and IF body", func(t *testing.T) {
		ddl := "CREATE DATABASE a; USE a;\n" +
			"CREATE DEFINER=`root`@`%` PROCEDURE `p`()\n" +
			"    DETERMINISTIC\n" +
			"    COMMENT 'c'\n" +
			"IF 1 THEN COMMIT; END IF"
		c, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		p := c.GetDatabase("a").Procedures["p"]
		if p == nil {
			t.Fatalf("procedure p not loaded")
		}
		if p.Body != "IF 1 THEN COMMIT; END IF" {
			t.Errorf("got body %q, want %q", p.Body, "IF 1 THEN COMMIT; END IF")
		}
	})
}
