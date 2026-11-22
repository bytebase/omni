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

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
	obparser "github.com/bytebase/omni/oceanbase/parser"
)

// TestProgramSimulation holds the catalog to testdata/programs/simulate.tsv:
// stored-program DDL statements (procedures, functions, triggers, events)
// across object fixtures, with what pinned OceanBase 4.3.5.6 did and what
// Exec does. Unlike TestSimulation (which reads only SHOW FULL TABLES so
// existing table/view row hashes remain stable), this suite reads tables,
// views, stored procedures, stored functions, triggers, and events from
// OceanBase and diffs all of them against the catalog.
//
// OceanBase's SQL grammar (sql_parser_mysql_mode.y, tracked in compat.tsv)
// contains only DROP FUNCTION; stored procedures, stored functions, triggers,
// and events live in OceanBase's separate PL grammar (pl_parser_mysql_mode.y).
// This suite combines the compat.tsv statements with a curated matrix covering
// procedures, functions, triggers, and events.
//
// OceanBase's side is on the test's tape
// (testdata/tapes/TestProgramSimulation.tape), so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned OceanBase again; -update-compat
// commits the rows a catalog fix brought to OceanBase's outcome; and
// -simulate-show=<text> prints both schemas of the rows that hold the text.
func TestProgramSimulation(t *testing.T) {
	rows, err := programSimulationRows(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   testcontainer.Codec{}.MadeUp(),
		Omni:     programSimulationOmni,
		Engine:   programSimulationEngine,
	})
}

var programSimulationFixtures = []string{
	"one_table@missing", "one_table@table",
	"objects@routine", "objects@trigger", "objects@event",
}

// programSimulatedStatement matches stored-program DDL in compat.tsv whose
// whole effect is in the database a row reads back. Loadable UDFs (SONAME)
// live outside the database schema and are excluded.
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

type programSpec struct {
	rule string
	stmt string
}

// curatedProgramStatements covers stored procedures, stored functions,
// triggers, and scheduled events in OceanBase MySQL mode across creation,
// alteration, dropping, characteristics, parameters, compound bodies, trigger
// timings/events/ordering/validation, event schedules/status/completion/rename,
// and schema qualification.
var curatedProgramStatements = []programSpec{
	// ---- Stored Procedures ----
	{"procedure/create_simple", "CREATE PROCEDURE a() COMMIT"},
	{"procedure/create_if_not_exists", "CREATE PROCEDURE IF NOT EXISTS a() COMMIT"},
	{"procedure/create_param_plain", "CREATE PROCEDURE a(x INT) COMMIT"},
	{"procedure/create_param_in", "CREATE PROCEDURE a(IN x INT) COMMIT"},
	{"procedure/create_param_out", "CREATE PROCEDURE a(OUT x INT) COMMIT"},
	{"procedure/create_param_inout", "CREATE PROCEDURE a(INOUT x INT) COMMIT"},
	{"procedure/create_param_types", "CREATE PROCEDURE a(IN s VARCHAR(20), OUT n DECIMAL(10,2), INOUT e ENUM('x','y')) COMMIT"},
	{"procedure/create_deterministic", "CREATE PROCEDURE a() DETERMINISTIC COMMIT"},
	{"procedure/create_not_deterministic", "CREATE PROCEDURE a() NOT DETERMINISTIC COMMIT"},
	{"procedure/create_no_sql_begin_end", "CREATE PROCEDURE a() NO SQL BEGIN END"},
	{"procedure/create_contains_sql_begin_end", "CREATE PROCEDURE a() CONTAINS SQL BEGIN END"},
	{"procedure/create_reads_sql_data_begin_end", "CREATE PROCEDURE a() READS SQL DATA BEGIN END"},
	{"procedure/create_modifies_sql_data_begin_end", "CREATE PROCEDURE a() MODIFIES SQL DATA BEGIN END"},
	{"procedure/create_no_sql", "CREATE PROCEDURE a() NO SQL COMMIT"},
	{"procedure/create_contains_sql", "CREATE PROCEDURE a() CONTAINS SQL COMMIT"},
	{"procedure/create_reads_sql_data", "CREATE PROCEDURE a() READS SQL DATA COMMIT"},
	{"procedure/create_modifies_sql_data", "CREATE PROCEDURE a() MODIFIES SQL DATA COMMIT"},
	{"procedure/create_security_definer", "CREATE PROCEDURE a() SQL SECURITY DEFINER COMMIT"},
	{"procedure/create_security_invoker", "CREATE PROCEDURE a() SQL SECURITY INVOKER COMMIT"},
	{"procedure/create_comment", "CREATE PROCEDURE a() COMMENT 's' COMMIT"},
	{"procedure/create_language_sql", "CREATE PROCEDURE a() LANGUAGE SQL COMMIT"},
	{"procedure/create_definer_root", "CREATE DEFINER = `root`@`%` PROCEDURE a() COMMIT"},
	{"procedure/create_definer_current_user", "CREATE DEFINER = CURRENT_USER PROCEDURE a() COMMIT"},
	{"procedure/create_begin_end", "CREATE PROCEDURE a() BEGIN END"},
	{"procedure/create_declare_var", "CREATE PROCEDURE a() BEGIN DECLARE x INT DEFAULT 1; SET x = x + 1; END"},
	{"procedure/create_declare_condition_handler", "CREATE PROCEDURE a() BEGIN DECLARE c CONDITION FOR 1062; DECLARE EXIT HANDLER FOR c COMMIT; END"},
	{"procedure/create_declare_continue_handler", "CREATE PROCEDURE a() BEGIN DECLARE CONTINUE HANDLER FOR SQLWARNING, NOT FOUND, SQLEXCEPTION COMMIT; END"},
	{"procedure/create_if", "CREATE PROCEDURE a() IF 1 THEN COMMIT; ELSEIF 0 THEN COMMIT; ELSE COMMIT; END IF"},
	{"procedure/create_case_searched", "CREATE PROCEDURE a() CASE WHEN 1 THEN COMMIT; ELSE COMMIT; END CASE"},
	{"procedure/create_case_simple", "CREATE PROCEDURE a() CASE 1 WHEN 1 THEN COMMIT; ELSE COMMIT; END CASE"},
	{"procedure/create_labeled_block", "CREATE PROCEDURE a() lbl: BEGIN END lbl"},
	{"procedure/create_loop", "CREATE PROCEDURE a() lbl: LOOP LEAVE lbl; END LOOP"},
	{"procedure/create_while", "CREATE PROCEDURE a() WHILE 0 DO COMMIT; END WHILE"},
	{"procedure/create_repeat", "CREATE PROCEDURE a() REPEAT COMMIT; UNTIL 1 END REPEAT"},
	{"procedure/create_qualified_unknown_db", "CREATE PROCEDURE a.a() COMMIT"},
	{"procedure/alter_empty", "ALTER PROCEDURE a"},
	{"procedure/alter_comment", "ALTER PROCEDURE a COMMENT 's'"},
	{"procedure/alter_language_sql", "ALTER PROCEDURE a LANGUAGE SQL"},
	{"procedure/alter_no_sql", "ALTER PROCEDURE a NO SQL"},
	{"procedure/alter_contains_sql", "ALTER PROCEDURE a CONTAINS SQL"},
	{"procedure/alter_reads_sql_data", "ALTER PROCEDURE a READS SQL DATA"},
	{"procedure/alter_modifies_sql_data", "ALTER PROCEDURE a MODIFIES SQL DATA"},
	{"procedure/alter_security_definer", "ALTER PROCEDURE a SQL SECURITY DEFINER"},
	{"procedure/alter_security_invoker", "ALTER PROCEDURE a SQL SECURITY INVOKER"},
	{"procedure/alter_qualified_unknown_db", "ALTER PROCEDURE a.a COMMENT 's'"},
	{"procedure/drop", "DROP PROCEDURE a"},
	{"procedure/drop_if_exists", "DROP PROCEDURE IF EXISTS a"},
	{"procedure/drop_qualified_unknown_db", "DROP PROCEDURE a.a"},
	{"procedure/drop_if_exists_qualified_unknown_db", "DROP PROCEDURE IF EXISTS a.a"},

	// ---- Stored Functions ----
	{"function/create_simple", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC RETURN 1"},
	{"function/create_if_not_exists", "CREATE FUNCTION IF NOT EXISTS a() RETURNS INT DETERMINISTIC RETURN 1"},
	{"function/create_param_int", "CREATE FUNCTION a(x INT) RETURNS INT DETERMINISTIC RETURN x"},
	{"function/create_param_varchar", "CREATE FUNCTION a(s VARCHAR(20)) RETURNS VARCHAR(40) DETERMINISTIC RETURN s"},
	{"function/create_return_decimal", "CREATE FUNCTION a(x DECIMAL(10,2)) RETURNS DECIMAL(12,4) DETERMINISTIC RETURN x"},
	{"function/create_return_enum", "CREATE FUNCTION a() RETURNS ENUM('x','y') DETERMINISTIC RETURN 'x'"},
	{"function/create_not_deterministic", "CREATE FUNCTION a() RETURNS INT NOT DETERMINISTIC NO SQL RETURN 1"},
	{"function/create_no_sql", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC NO SQL RETURN 1"},
	{"function/create_contains_sql", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC CONTAINS SQL RETURN 1"},
	{"function/create_reads_sql_data", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC READS SQL DATA RETURN 1"},
	{"function/create_modifies_sql_data", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC MODIFIES SQL DATA RETURN 1"},
	{"function/create_security_definer", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC SQL SECURITY DEFINER RETURN 1"},
	{"function/create_security_invoker", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC SQL SECURITY INVOKER RETURN 1"},
	{"function/create_comment", "CREATE FUNCTION a() RETURNS INT DETERMINISTIC COMMENT 's' RETURN 1"},
	{"function/create_definer_root", "CREATE DEFINER = `root`@`%` FUNCTION a() RETURNS INT DETERMINISTIC RETURN 1"},
	{"function/create_begin_end", "CREATE FUNCTION a(x INT) RETURNS INT DETERMINISTIC BEGIN DECLARE r INT DEFAULT 1; SET r = x + r; RETURN r; END"},
	{"function/create_qualified_unknown_db", "CREATE FUNCTION a.a() RETURNS INT DETERMINISTIC RETURN 1"},
	{"function/alter_empty", "ALTER FUNCTION a"},
	{"function/alter_comment", "ALTER FUNCTION a COMMENT 's'"},
	{"function/alter_language_sql", "ALTER FUNCTION a LANGUAGE SQL"},
	{"function/alter_no_sql", "ALTER FUNCTION a NO SQL"},
	{"function/alter_contains_sql", "ALTER FUNCTION a CONTAINS SQL"},
	{"function/alter_reads_sql_data", "ALTER FUNCTION a READS SQL DATA"},
	{"function/alter_modifies_sql_data", "ALTER FUNCTION a MODIFIES SQL DATA"},
	{"function/alter_security_definer", "ALTER FUNCTION a SQL SECURITY DEFINER"},
	{"function/alter_security_invoker", "ALTER FUNCTION a SQL SECURITY INVOKER"},
	{"function/alter_qualified_unknown_db", "ALTER FUNCTION a.a COMMENT 's'"},
	{"function/drop_qualified_unknown_db", "DROP FUNCTION a.a"},
	{"function/drop_if_exists_qualified_unknown_db", "DROP FUNCTION IF EXISTS a.a"},

	// ---- Triggers ----
	{"trigger/create_before_insert_on_a", "CREATE TRIGGER a BEFORE INSERT ON a FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_after_insert_on_a", "CREATE TRIGGER a AFTER INSERT ON a FOR EACH ROW SET @x = 1"},
	{"trigger/create_before_update_on_a", "CREATE TRIGGER a BEFORE UPDATE ON a FOR EACH ROW SET NEW.c = OLD.c + 1"},
	{"trigger/create_after_update_on_a", "CREATE TRIGGER a AFTER UPDATE ON a FOR EACH ROW SET @x = NEW.c"},
	{"trigger/create_before_delete_on_a", "CREATE TRIGGER a BEFORE DELETE ON a FOR EACH ROW SET @x = OLD.c"},
	{"trigger/create_after_delete_on_a", "CREATE TRIGGER a AFTER DELETE ON a FOR EACH ROW SET @x = OLD.c"},
	{"trigger/create_before_insert_on_t", "CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_after_insert_on_t_begin", "CREATE TRIGGER a AFTER INSERT ON t FOR EACH ROW BEGIN SET @x = 1; END"},
	{"trigger/create_if_not_exists_on_t", "CREATE TRIGGER IF NOT EXISTS a BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_follows", "CREATE TRIGGER a2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS a SET NEW.c = 2"},
	{"trigger/create_precedes", "CREATE TRIGGER a2 BEFORE INSERT ON t FOR EACH ROW PRECEDES a SET NEW.c = 2"},
	{"trigger/create_definer_root", "CREATE DEFINER = `root`@`%` TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_old_in_insert", "CREATE TRIGGER a BEFORE INSERT ON a FOR EACH ROW SET @x = OLD.c"},
	{"trigger/create_new_in_delete", "CREATE TRIGGER a BEFORE DELETE ON a FOR EACH ROW SET @x = NEW.c"},
	{"trigger/create_set_new_in_after", "CREATE TRIGGER a AFTER INSERT ON a FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_qualified_unknown_db", "CREATE TRIGGER a.a BEFORE INSERT ON a FOR EACH ROW SET NEW.c = 1"},
	{"trigger/create_qualified_table_unknown_db", "CREATE TRIGGER a BEFORE INSERT ON a.a FOR EACH ROW SET NEW.c = 1"},
	{"trigger/drop", "DROP TRIGGER a"},
	{"trigger/drop_if_exists", "DROP TRIGGER IF EXISTS a"},
	{"trigger/drop_qualified_unknown_db", "DROP TRIGGER a.a"},
	{"trigger/drop_if_exists_qualified_unknown_db", "DROP TRIGGER IF EXISTS a.a"},

	// ---- Scheduled Events ----
	{"event/create_every_starts", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO COMMIT"},
	{"event/create_if_not_exists", "CREATE EVENT IF NOT EXISTS a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO COMMIT"},
	{"event/create_every_default", "CREATE EVENT a ON SCHEDULE EVERY 1 HOUR DO COMMIT"},
	{"event/create_every_starts_ends", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ENDS '2031-01-01 00:00:00' DO COMMIT"},
	{"event/create_at_preserve", "CREATE EVENT a ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION PRESERVE DO COMMIT"},
	{"event/create_at_not_preserve", "CREATE EVENT a ON SCHEDULE AT '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DO COMMIT"},
	{"event/create_enable", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION PRESERVE ENABLE DO COMMIT"},
	{"event/create_disable", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE DISABLE DO COMMIT"},
	{"event/create_comment", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' COMMENT 's' DO COMMIT"},
	{"event/create_definer_root", "CREATE DEFINER = `root`@`%` EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO COMMIT"},
	{"event/create_begin_end", "CREATE EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO BEGIN COMMIT; END"},
	{"event/create_at_colref", "CREATE EVENT a ON SCHEDULE AT a DO COMMIT"},
	{"event/alter_schedule_every", "ALTER EVENT a ON SCHEDULE EVERY 2 DAY STARTS '2030-01-01 00:00:00'"},
	{"event/alter_schedule_ends", "ALTER EVENT a ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ENDS '2031-01-01 00:00:00'"},
	{"event/alter_schedule_at", "ALTER EVENT a ON SCHEDULE AT '2030-06-01 00:00:00'"},
	{"event/alter_on_completion_preserve", "ALTER EVENT a ON COMPLETION PRESERVE"},
	{"event/alter_on_completion_not_preserve", "ALTER EVENT a ON COMPLETION NOT PRESERVE"},
	{"event/alter_enable", "ALTER EVENT a ENABLE"},
	{"event/alter_disable", "ALTER EVENT a DISABLE"},
	{"event/alter_comment", "ALTER EVENT a COMMENT 's'"},
	{"event/alter_do_commit", "ALTER EVENT a DO COMMIT"},
	{"event/alter_do_begin_end", "ALTER EVENT a DO BEGIN COMMIT; END"},
	{"event/alter_do_set", "ALTER EVENT a DO SET @x = 1"},
	{"event/alter_rename_same", "ALTER EVENT a RENAME TO a"},
	{"event/alter_rename_other", "ALTER EVENT a RENAME TO b"},
	{"event/alter_schedule_at_colref", "ALTER EVENT a ON SCHEDULE AT a"},
	{"event/alter_schedule_starts_colref", "ALTER EVENT a ON SCHEDULE EVERY 1 DAY STARTS a"},
	{"event/alter_schedule_ends_colref", "ALTER EVENT a ON SCHEDULE EVERY 1 DAY ENDS a"},
	{"event/alter_definer_root", "ALTER DEFINER = `root`@`%` EVENT a COMMENT 's'"},
	{"event/drop", "DROP EVENT a"},
	{"event/drop_if_exists", "DROP EVENT IF EXISTS a"},
}

func programSimulationRows(compatPath string) ([]simulate.Generated, error) {
	fromCompat, err := simulate.FromCompat(compatPath, []string{"_"}, programSimulated)
	if err != nil {
		return nil, err
	}
	var statements []programSpec
	seen := make(map[string]bool)
	for _, r := range fromCompat {
		if seen[r.Statement] {
			continue
		}
		seen[r.Statement] = true
		statements = append(statements, programSpec{rule: r.Rule, stmt: r.Statement})
	}
	for _, s := range curatedProgramStatements {
		if seen[s.stmt] {
			continue
		}
		seen[s.stmt] = true
		statements = append(statements, s)
	}
	rows := make([]simulate.Generated, 0, len(programSimulationFixtures)*len(statements))
	for _, fx := range programSimulationFixtures {
		for _, s := range statements {
			rows = append(rows, simulate.Generated{
				Fixture:   fx,
				Rule:      s.rule,
				Statement: s.stmt,
			})
		}
	}
	return rows, nil
}

func programSimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	testcontainer.Reset(t, func(ctx context.Context, db *sql.DB) error {
		rows, err := db.QueryContext(ctx,
			"SELECT EVENT_NAME FROM information_schema.EVENTS WHERE EVENT_SCHEMA = ?", simulationDatabase)
		if err != nil {
			return err
		}
		var names []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				rows.Close()
				return err
			}
			names = append(names, name)
		}
		if err := rows.Close(); err != nil {
			return err
		}
		if err := rows.Err(); err != nil {
			return err
		}
		if len(names) > 0 {
			if _, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+simulationDatabase); err != nil {
				return err
			}
			if _, err := db.ExecContext(ctx, "USE "+simulationDatabase); err != nil {
				return err
			}
			for _, name := range names {
				if _, err := db.ExecContext(ctx, fmt.Sprintf("DROP EVENT IF EXISTS `%s`", strings.ReplaceAll(name, "`", "``"))); err != nil {
					return err
				}
			}
		}
		_, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS "+simulationDatabase)
		return err
	})
	db := testcontainer.OpenWith(t, testcontainer.Options{MultiStatements: true})
	ctx := context.Background()
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		build := fmt.Sprintf("DROP DATABASE IF EXISTS %[1]s; CREATE DATABASE %[1]s; USE %[1]s; %[2]s", simulationDatabase, fx.SQL)
		if _, err := db.ExecContext(ctx, build); err != nil {
			return simulate.Answer{}, fmt.Errorf("build fixture %s: %w", fx.Name, err)
		}
		var ans simulate.Answer
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var refusal *gomysql.MySQLError
			if !errors.As(err, &refusal) {
				return simulate.Answer{}, err
			}
			ans = simulate.Answer{Code: strconv.Itoa(int(refusal.Number))}
		} else {
			forms, err := programSimulationForms(ctx, db)
			if err != nil {
				return simulate.Answer{}, err
			}
			ans = simulate.Answer{Forms: forms}
		}
		// OceanBase 4.3.5.6 does not drop scheduled events when their database
		// is dropped, so remove any events this row's fixture or statement may
		// have created before the next row drops the database.
		if fx.Name == "objects@event" || strings.Contains(statement, "EVENT") {
			if _, err := db.ExecContext(ctx, "DROP EVENT IF EXISTS a; DROP EVENT IF EXISTS b"); err != nil {
				return simulate.Answer{}, fmt.Errorf("cleanup events after %s: %w", fx.Name, err)
			}
		}
		return ans, nil
	}
}

// programSimulationForms reads tables, views, stored functions, stored
// procedures, triggers, and events from the selected database in deterministic
// order. OceanBase has no SHOW CREATE EVENT (1064), so events are rebuilt
// deterministically from information_schema.EVENTS without STARTS (which
// OceanBase auto-fills with the creation timestamp when unspecified).
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

	// Events (rebuilt from information_schema.EVENTS because OceanBase has no SHOW CREATE EVENT).
	events, err := queryStoredEvents(ctx, db)
	if err != nil {
		return nil, err
	}
	forms = append(forms, events...)

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

func queryStoredEvents(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	rows, err := db.QueryContext(ctx, `SELECT EVENT_NAME, EVENT_TYPE, IFNULL(DATE_FORMAT(EXECUTE_AT, '%Y-%m-%d %H:%i:%s'), ''),
		IFNULL(INTERVAL_VALUE, ''), IFNULL(INTERVAL_FIELD, ''), IFNULL(DATE_FORMAT(ENDS, '%Y-%m-%d %H:%i:%s'), ''),
		STATUS, ON_COMPLETION, IFNULL(EVENT_COMMENT, ''), EVENT_DEFINITION
	FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() ORDER BY EVENT_NAME`)
	if err != nil {
		return nil, err
	}
	var forms []simulate.Form
	for rows.Next() {
		var name, kind, at, every, field, ends, status, completion, comment, body string
		if err := rows.Scan(&name, &kind, &at, &every, &field, &ends, &status, &completion, &comment, &body); err != nil {
			rows.Close()
			return nil, err
		}
		var b strings.Builder
		fmt.Fprintf(&b, "CREATE EVENT `%s` ON SCHEDULE ", strings.ReplaceAll(name, "`", "``"))
		if kind == "ONE TIME" {
			fmt.Fprintf(&b, "AT '%s'", at)
		} else {
			fmt.Fprintf(&b, "EVERY %s %s", every, field)
			if ends != "" {
				fmt.Fprintf(&b, " ENDS '%s'", ends)
			}
		}
		fmt.Fprintf(&b, " ON COMPLETION %s", completion)
		switch status {
		case "ENABLED":
			b.WriteString(" ENABLE")
		case "DISABLED":
			b.WriteString(" DISABLE")
		case "SLAVESIDE_DISABLED":
			b.WriteString(" DISABLE ON SLAVE")
		default:
			rows.Close()
			return nil, fmt.Errorf("event %s has STATUS %q", name, status)
		}
		if comment != "" {
			fmt.Fprintf(&b, " COMMENT %s", quoteStringLiteral(comment))
		}
		fmt.Fprintf(&b, " DO %s", strings.TrimSuffix(strings.TrimSpace(body), ";"))
		forms = append(forms, simulate.Form{Kind: "EVENT", Name: name, Text: b.String()})
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return forms, nil
}

func programSimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	database := fmt.Sprintf("CREATE DATABASE %[1]s DEFAULT CHARSET=utf8mb4;\nUSE %[1]s;\n", simulationDatabase)
	c, err := LoadSQL(database + fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	results, err := c.Exec(statement, nil)
	var syntax *obparser.ParseError
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

	want, err := LoadSQL(database + "SET foreign_key_checks = 0;\n")
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
				return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what OceanBase stores (%s %s): %v", f.Kind, f.Name, execErr)}
			}
			for _, r := range res {
				if r.Error != nil {
					return simulate.Omni{Why: fmt.Sprintf("the catalog does not load what OceanBase stores (%s %s): %v", f.Kind, f.Name, r.Error)}
				}
			}
		}
	}
	if _, err := want.Exec("SET foreign_key_checks = 1", nil); err != nil {
		return simulate.Omni{Why: fmt.Sprintf("restore FK checks: %v", err)}
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
	return simulate.Omni{Why: "to store what OceanBase stores, the catalog would run: " + strings.Join(strings.Fields(plan), " ")}
}

func TestProgramSimulationBugFixes(t *testing.T) {
	t.Run("duplicate trigger returns OceanBase 5734", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE t (b INT PRIMARY KEY, c INT); CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		if got := omniCode(c, "CREATE TRIGGER a BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1"); got != ErrDupTrigger {
			t.Errorf("got %d, want %d", got, ErrDupTrigger)
		}
	})

	t.Run("FOLLOWS/PRECEDES missing referenced trigger returns 3011", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE TABLE t (b INT PRIMARY KEY, c INT);")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		for _, stmt := range []string{
			"CREATE TRIGGER a2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS a SET NEW.c = 2",
			"CREATE TRIGGER a2 BEFORE INSERT ON t FOR EACH ROW PRECEDES a SET NEW.c = 2",
		} {
			if got := omniCode(c, stmt); got != ErrReferencedTriggerNotFound {
				t.Errorf("%s: got %d, want %d", stmt, got, ErrReferencedTriggerNotFound)
			}
		}
		if len(c.GetDatabase("d").Triggers) != 0 {
			t.Fatalf("failed trigger creation mutated catalog")
		}
	})

	t.Run("trigger qualified table unknown db returns 1146 and qualified DROP TRIGGER unknown db returns 1360", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		if got := omniCode(c, "CREATE TRIGGER a BEFORE INSERT ON nosuch.a FOR EACH ROW SET NEW.c = 1"); got != ErrNoSuchTable {
			t.Errorf("CREATE TRIGGER on unknown db table: got %d, want %d", got, ErrNoSuchTable)
		}
		if got := omniCode(c, "CREATE TRIGGER nosuch.a BEFORE INSERT ON a FOR EACH ROW SET NEW.c = 1"); got != ErrUnknownDatabase {
			t.Errorf("CREATE TRIGGER in unknown db: got %d, want %d", got, ErrUnknownDatabase)
		}
		if got := omniCode(c, "DROP TRIGGER nosuch.a"); got != ErrNoSuchTrigger {
			t.Errorf("DROP TRIGGER in unknown db: got %d, want %d", got, ErrNoSuchTrigger)
		}
		if got := omniCode(c, "DROP TRIGGER IF EXISTS nosuch.a"); got != 0 {
			t.Errorf("DROP TRIGGER IF EXISTS in unknown db: got %d, want 0", got)
		}
	})

	t.Run("DROP PROCEDURE/FUNCTION IF EXISTS unknown db returns 1049", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		for _, stmt := range []string{
			"DROP PROCEDURE IF EXISTS nosuch.p",
			"DROP FUNCTION IF EXISTS nosuch.f",
		} {
			if got := omniCode(c, stmt); got != ErrUnknownDatabase {
				t.Errorf("%s: got %d, want %d", stmt, got, ErrUnknownDatabase)
			}
		}
	})

	t.Run("procedure containing COMMIT defaults DATA ACCESS to MODIFIES SQL DATA", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE PROCEDURE p_commit() NO SQL COMMIT; CREATE PROCEDURE p_empty() NO SQL BEGIN END;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		if got := c.GetDatabase("d").Procedures["p_commit"].Characteristics["DATA ACCESS"]; got != "MODIFIES SQL DATA" {
			t.Errorf("p_commit DATA ACCESS = %q, want MODIFIES SQL DATA", got)
		}
		if got := c.GetDatabase("d").Procedures["p_empty"].Characteristics["DATA ACCESS"]; got != "NO SQL" {
			t.Errorf("p_empty DATA ACCESS = %q, want NO SQL", got)
		}
	})

	t.Run("event schedule colref returns 1525 and RENAME TO same name returns 1551", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO COMMIT;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		for _, tc := range []struct {
			stmt string
			want int
		}{
			{"CREATE EVENT e ON SCHEDULE AT a DO COMMIT", ErrWrongValue},
			{"CREATE EVENT missing ON SCHEDULE AT a DO COMMIT", ErrWrongValue},
			{"ALTER EVENT missing ON SCHEDULE AT a", ErrWrongValue},
			{"ALTER EVENT missing ON SCHEDULE EVERY 1 DAY STARTS a", ErrWrongValue},
			{"ALTER EVENT missing ON SCHEDULE EVERY 1 DAY ENDS a", ErrWrongValue},
			{"ALTER EVENT e RENAME TO e", ErrEventSameName},
			{"ALTER EVENT e RENAME TO E", ErrEventSameName},
			{"ALTER EVENT missing RENAME TO missing", ErrEventSameName},
		} {
			if got := omniCode(c, tc.stmt); got != tc.want {
				t.Errorf("%s: got %d, want %d", tc.stmt, got, tc.want)
			}
		}
	})

	t.Run("LoadSQL parses SHOW CREATE PROCEDURE with DEFINER and IF body", func(t *testing.T) {
		ddl := "CREATE DATABASE a; USE a;\n" +
			"CREATE DEFINER = `root`@`%` PROCEDURE `p`()\n" +
			"    MODIFIES SQL DATA\n" +
			"    COMMENT 'c'\n" +
			"IF 1 THEN COMMIT; ELSEIF 0 THEN COMMIT; ELSE COMMIT; END IF"
		c, err := LoadSQL(ddl)
		if err != nil {
			t.Fatalf("LoadSQL failed: %v", err)
		}
		p := c.GetDatabase("a").Procedures["p"]
		if p == nil {
			t.Fatalf("procedure p not loaded")
		}
		want := "IF 1 THEN COMMIT; ELSEIF 0 THEN COMMIT; ELSE COMMIT; END IF"
		if p.Body != want {
			t.Errorf("got body %q, want %q", p.Body, want)
		}
	})

	t.Run("event single-statement BEGIN END body is unwrapped", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO BEGIN COMMIT; END;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		if got := c.GetDatabase("d").Events["e"].Body; got != "COMMIT" {
			t.Errorf("CREATE EVENT body = %q, want COMMIT", got)
		}
		if got := omniCode(c, "ALTER EVENT e DO BEGIN SET @x = 1; END"); got != 0 {
			t.Fatalf("ALTER EVENT failed with %d", got)
		}
		if got := c.GetDatabase("d").Events["e"].Body; got != "SET @x = 1" {
			t.Errorf("ALTER EVENT body = %q, want %q", got, "SET @x = 1")
		}
	})
}
