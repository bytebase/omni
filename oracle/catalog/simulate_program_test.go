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

	"github.com/sijms/go-ora/v2/network"

	"github.com/bytebase/omni/enginetest/simulate"
	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/internal/testcontainer"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// TestProgramSimulation holds the catalog to testdata/programs/simulate.tsv:
// procedures, functions, packages and package bodies, triggers, and object
// and collection types with bodies, in programFeatures, run on
// programs@objects, with what Oracle did and what Exec does. testdata's
// simulate.tsv leaves PL/SQL out; these hold the catalog to how Oracle
// compiles a unit when it creates it, and to what a change to an object a
// unit references does to the unit.
//
// A row's outcome is the ORA- error Oracle refuses the statement with, or
// ok: and a hash of what the schema holds afterwards: dump_forms's tables,
// views, sequences and synonyms, and each PL/SQL unit's kind, name and
// status, a trigger's ENABLED or DISABLED, and the first error USER_ERRORS
// lists for a unit that did not compile. Oracle creates such a unit, INVALID,
// and the driver reports no error (SQL*Plus prints a warning).
//
// Oracle's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned Oracle again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestProgramSimulation(t *testing.T) {
	rows := make([]simulate.Generated, 0, len(programFeatures))
	seen := map[string]bool{}
	for _, f := range programFeatures {
		if seen[f.statement] {
			t.Fatalf("programFeatures holds %q twice", f.statement)
		}
		seen[f.statement] = true
		rows = append(rows, simulate.Generated{Fixture: programFixture, Rule: "program:" + f.kind, Statement: f.statement})
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     rows,
		MadeUp:   (testcontainer.Codec{}).MadeUp(),
		Omni:     programSimulationOmni,
		Engine:   programSimulationEngine,
	})
}

const programFixture = "programs@objects"

// programGrants add the PL/SQL privileges to the scratch account's.
const programGrants = simGrants + ", CREATE PROCEDURE, CREATE TRIGGER, CREATE TYPE"

// programSimulationEngine builds the fixture one statement at a time, as
// reset_schema's split on semicolons cannot hold a PL/SQL unit, runs the
// row, and reads what the schema holds.
func programSimulationEngine(t testing.TB) func(simulate.Fixture, string) (simulate.Answer, error) {
	ctx := context.Background()
	db := simulationAccount(t, programGrants)
	return func(fx simulate.Fixture, statement string) (simulate.Answer, error) {
		if _, err := db.ExecContext(ctx, "BEGIN system.omni_sim.reset_schema(NULL); END;"); err != nil {
			return simulate.Answer{}, fmt.Errorf("reset for fixture %s: %w", fx.Name, err)
		}
		for _, seg := range oracleparser.Split(fx.SQL) {
			if seg.Empty() {
				continue
			}
			if _, err := db.ExecContext(ctx, strings.TrimSpace(seg.Text)); err != nil {
				return simulate.Answer{}, fmt.Errorf("build fixture %s (%s): %w", fx.Name, strings.TrimSpace(seg.Text), err)
			}
		}
		if _, err := db.ExecContext(ctx, statement); err != nil {
			var oe *network.OracleError
			if !errors.As(err, &oe) {
				return simulate.Answer{}, err
			}
			return simulate.Answer{Code: strconv.Itoa(oe.ErrCode)}, nil
		}
		var raw sql.NullString
		if err := db.QueryRowContext(ctx, "SELECT system.omni_sim.dump_forms FROM dual").Scan(&raw); err != nil {
			return simulate.Answer{}, fmt.Errorf("read stored schema: %w", err)
		}
		forms := parseSimulationForms(raw.String)
		units, err := programUnits(ctx, db)
		if err != nil {
			return simulate.Answer{}, err
		}
		return simulate.Answer{Forms: append(forms, units...)}, nil
	}
}

// programUnits reads each PL/SQL unit's kind, name and status, a trigger's
// ENABLED or DISABLED, and the first error USER_ERRORS lists for the unit,
// as one form each. Oracle creates a unit it cannot compile, INVALID, and
// the driver reports no error for it.
func programUnits(ctx context.Context, db *sql.DB) ([]simulate.Form, error) {
	first := map[[2]string]string{}
	errRows, err := db.QueryContext(ctx, `SELECT type, name, text FROM user_errors WHERE attribute = 'ERROR' ORDER BY type, name, sequence`)
	if err != nil {
		return nil, fmt.Errorf("read errors: %w", err)
	}
	defer errRows.Close()
	for errRows.Next() {
		var kind, name, text string
		if err := errRows.Scan(&kind, &name, &text); err != nil {
			return nil, err
		}
		key := [2]string{kind, name}
		if code := plsqlErrorCode(text); code != "" && first[key] == "" {
			first[key] = code
		}
	}
	if err := errRows.Err(); err != nil {
		return nil, err
	}
	rows, err := db.QueryContext(ctx, `SELECT o.object_type, o.object_name, o.status, tr.status
FROM user_objects o LEFT JOIN user_triggers tr ON o.object_type = 'TRIGGER' AND tr.trigger_name = o.object_name
WHERE o.object_type IN ('PROCEDURE', 'FUNCTION', 'PACKAGE', 'PACKAGE BODY', 'TRIGGER', 'TYPE', 'TYPE BODY')
  AND o.object_name NOT LIKE 'SYS\_PLSQL\_%' ESCAPE '\'
ORDER BY o.object_type, o.object_name`)
	if err != nil {
		return nil, fmt.Errorf("read units: %w", err)
	}
	defer rows.Close()
	var forms []simulate.Form
	for rows.Next() {
		var kind, name, status string
		var enabled sql.NullString
		if err := rows.Scan(&kind, &name, &status, &enabled); err != nil {
			return nil, err
		}
		forms = append(forms, programUnitForm(unitState{kind: kind, name: name, status: status, enabled: enabled.String, err: first[[2]string{kind, name}]}))
	}
	return forms, rows.Err()
}

// programUnitForm renders a unit as one line: its kind, name and status,
// then ENABLED or DISABLED for a trigger and the first error for a unit
// that did not compile.
func programUnitForm(u unitState) simulate.Form {
	text := fmt.Sprintf("%s %s %s", u.kind, quoteIdent(u.name), u.status)
	for _, w := range []string{u.enabled, u.err} {
		if w != "" {
			text += " " + w
		}
	}
	return simulate.Form{Kind: u.kind, Name: u.name, Text: text + "\n"}
}

var plsqlErrorText = regexp.MustCompile(`^(?:PL/SQL: )?((?:PLS|ORA)-\d{5}):`)

// plsqlErrorCode reads the code an error line of USER_ERRORS starts with:
// PLS-nnnnn, or the ORA-nnnnn of a SQL statement as "PL/SQL: ORA-nnnnn".
// The "Statement ignored" line that follows each has none.
func plsqlErrorCode(text string) string {
	m := plsqlErrorText.FindStringSubmatch(strings.TrimSpace(text))
	if m == nil || m[1] == "PLS-00000" {
		return ""
	}
	return m[1]
}

// programSimulationOmni runs the statement on the fixture in the catalog and
// renders the schema as the engine's side does.
func programSimulationOmni(fx simulate.Fixture, statement string, engine simulate.Answer) simulate.Omni {
	c, err := LoadSQL(fx.SQL)
	if err != nil {
		panic(fmt.Sprintf("the catalog does not load fixture %s: %v", fx.Name, err))
	}
	c.meta.Name = simUser
	if len(c.meta.Schemas) > 0 {
		c.meta.Schemas[0].Name = simUser
	}
	c.SetServerConfig(StandardOracle23ServerConfig())
	results, err := c.Exec(statement, nil)
	if err != nil {
		return simulate.Omni{Code: simulate.NoCode, Why: err.Error()}
	}
	for _, r := range results {
		if r.Error == nil {
			continue
		}
		var refusal *Error
		if !errors.As(r.Error, &refusal) {
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
		code := oracleErrorCode(refusal.Code)
		if code == "" {
			return simulate.Omni{Code: simulate.NoCode, Why: r.Error.Error()}
		}
		return simulate.Omni{Code: code, Why: r.Error.Error()}
	}
	if engine.Refused() {
		return simulate.Omni{}
	}
	var tables, units []simulate.Form
	for _, f := range engine.Forms {
		if programUnitKinds[f.Kind] {
			units = append(units, f)
		} else {
			tables = append(tables, f)
		}
	}
	var why []string
	if got, want := renderUnits(catalogUnitForms(c)), renderUnits(units); got != want {
		why = append(why, "the catalog's units:\n"+got+"Oracle's units:\n"+want)
	}
	if o := compareStoredSchema(withoutPrograms(c.Metadata()), tables); !o.Same {
		why = append(why, o.Why)
	}
	if len(why) == 0 {
		return simulate.Omni{Same: true}
	}
	return simulate.Omni{Why: strings.Join(why, "\n")}
}

var programUnitKinds = map[string]bool{
	"PROCEDURE": true, "FUNCTION": true, "PACKAGE": true, "PACKAGE BODY": true,
	"TRIGGER": true, "TYPE": true, "TYPE BODY": true,
}

func renderUnits(forms []simulate.Form) string {
	var lines []string
	for _, f := range forms {
		lines = append(lines, f.Text)
	}
	slices.Sort(lines)
	return strings.Join(lines, "")
}

// programFeatures are the statements the suite runs, labelled
// program:<kind>. Many are refused by Oracle on purpose.
var programFeatures = []struct{ kind, statement string }{
	// Procedures: static SQL.
	{"procedure", "CREATE PROCEDURE p IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p (p_id IN NUMBER) IS v_name VARCHAR2(100); BEGIN SELECT name INTO v_name FROM customers WHERE id = p_id; END;"},
	{"procedure", "CREATE PROCEDURE p (p_id IN NUMBER) IS v_name VARCHAR2(100); BEGIN SELECT nope INTO v_name FROM customers WHERE id = p_id; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT id INTO v FROM missing_t; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT c.nope INTO v FROM customers c; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM customers c JOIN orders o ON o.customer_id = c.id WHERE c.region = 'EU'; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN SELECT id FROM customers; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO orders (id, customer_id, amount) VALUES (1, 1, 2.5); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO orders (nope) VALUES (1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO orders (id) VALUES (1, 2); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO orders (id, amount) VALUES (1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO missing_t VALUES (1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO orders (id, customer_id) SELECT id, id FROM customers; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO cust (id, name) VALUES (1, 'x'); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO v_customers (id, name) VALUES (1, 'x'); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN UPDATE orders SET status = 'x' WHERE id = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN UPDATE orders SET nope = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN UPDATE orders SET status = 'x' WHERE nope = 1; END;"},
	{"procedure", "CREATE PROCEDURE p (p_id IN NUMBER) IS BEGIN DELETE FROM orders WHERE customer_id = p_id; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN DELETE FROM orders WHERE nope = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN DELETE FROM missing_t; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN MERGE INTO orders o USING customers c ON (o.customer_id = c.id) WHEN MATCHED THEN UPDATE SET o.status = c.region; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN MERGE INTO orders o USING customers c ON (o.customer_id = c.nope) WHEN MATCHED THEN UPDATE SET o.status = c.region; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT order_seq.NEXTVAL INTO v FROM dual; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := order_seq.NEXTVAL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := nope_seq.NEXTVAL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT f_double(id) INTO v FROM customers WHERE ROWNUM = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT nope_f(id) INTO v FROM customers WHERE ROWNUM = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM v_customers WHERE name = 'x'; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM user_tables; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN COMMIT; ROLLBACK; SAVEPOINT s; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN CREATE TABLE x (a NUMBER); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN EXECUTE IMMEDIATE 'SELECT nope FROM nowhere'; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN EXECUTE IMMEDIATE 'CREATE TABLE x (a NUMBER)'; END;"},
	// Procedures: cursors and loops.
	{"procedure", "CREATE PROCEDURE p IS BEGIN FOR r IN (SELECT id, name FROM customers) LOOP NULL; END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS v VARCHAR2(100); BEGIN FOR r IN (SELECT id FROM customers) LOOP v := r.name; END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN FOR r IN (SELECT nope FROM customers) LOOP NULL; END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS CURSOR c IS SELECT id FROM customers; v NUMBER; BEGIN OPEN c; FETCH c INTO v; CLOSE c; END;"},
	{"procedure", "CREATE PROCEDURE p IS CURSOR c IS SELECT nope FROM customers; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN OPEN nope_c; END;"},
	{"procedure", "CREATE PROCEDURE p IS rc SYS_REFCURSOR; BEGIN OPEN rc FOR SELECT id FROM orders; CLOSE rc; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER := 0; BEGIN FOR i IN 1 .. 10 LOOP v := v + i; END LOOP; WHILE v > 0 LOOP v := v - 1; END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS TYPE t_ids IS TABLE OF NUMBER; ids t_ids; BEGIN SELECT id BULK COLLECT INTO ids FROM customers; FORALL i IN 1 .. ids.COUNT DELETE FROM orders WHERE customer_id = ids(i); END;"},
	// Procedures: names and calls.
	{"procedure", "CREATE PROCEDURE p IS BEGIN nope_proc; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p_existing(1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p_existing; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p_existing(1, 2); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p_existing(p_id => 1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p_existing(nope => 1); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := f_double(2); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := f_double(2, 3); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN f_double(2); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := nope_f(2); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := nope; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN nope := 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := 'x'; END;"},
	{"procedure", "CREATE PROCEDURE p IS v DATE; BEGIN v := SYSDATE + 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := LENGTH(USER) + SQLCODE; DBMS_OUTPUT.PUT_LINE(SQLERRM); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN pk.log_it('x'); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN pk.nope; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN pk.g_count := pk.g_count + 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := pk.total(1); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN pk.log_it; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN nope_pkg.p; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN DBMS_OUTPUT.PUT_LINE('x'); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN RAISE_APPLICATION_ERROR(-20001, 'x'); END;"},
	{"procedure", "CREATE PROCEDURE p IS PROCEDURE inner IS BEGIN NULL; END; BEGIN inner; END;"},
	{"procedure", "CREATE PROCEDURE p IS FUNCTION twice (x NUMBER) RETURN NUMBER IS BEGIN RETURN x * 2; END; v NUMBER; BEGIN v := twice(1, 2); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN p; END;"},
	// Procedures: declarations.
	{"procedure", "CREATE PROCEDURE p IS v customers.name%TYPE; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v customers.nope%TYPE; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v missing_t.c%TYPE; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS r customers%ROWTYPE; BEGIN SELECT * INTO r FROM customers WHERE id = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS r customers%ROWTYPE; BEGIN r.nope := 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS r missing_t%ROWTYPE; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v nope_type; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v t_point := t_point(1, 2); n NUMBER; BEGIN n := v.norm(); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; v NUMBER; BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; v NUMBER; BEGIN v := 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS c CONSTANT NUMBER := 1; BEGIN c := 2; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x IN NUMBER) IS BEGIN p_x := 1; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x OUT NUMBER) IS BEGIN p_x := 1; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x IN OUT NOCOPY VARCHAR2, p_y IN NUMBER DEFAULT 0) IS BEGIN p_x := p_x || p_y; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x IN nope_type) IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x IN customers.id%TYPE) IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p (p_x IN NUMBER, p_x IN NUMBER) IS BEGIN NULL; END;"},
	// Procedures: control flow and exceptions.
	{"procedure", "CREATE PROCEDURE p IS BEGIN RETURN; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN RETURN 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN NULL; EXCEPTION WHEN NO_DATA_FOUND THEN NULL; WHEN OTHERS THEN RAISE; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN NULL; EXCEPTION WHEN nope_exc THEN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS e EXCEPTION; BEGIN RAISE e; EXCEPTION WHEN e THEN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN RAISE nope_exc; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN GOTO nowhere; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN <<l>> NULL; GOTO l; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER := 1; BEGIN IF v > 0 THEN NULL; ELSIF v < 0 THEN NULL; ELSE NULL; END IF; CASE v WHEN 1 THEN NULL; ELSE NULL; END CASE; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN NULL END;"},
	{"procedure", "CREATE PROCEDURE p IS PRAGMA AUTONOMOUS_TRANSACTION; BEGIN INSERT INTO audit_log (id) VALUES (1); COMMIT; END;"},
	{"procedure", "CREATE PROCEDURE p AUTHID CURRENT_USER IS BEGIN NULL; END;"},
	// Procedures: everyday forms that name nothing missing.
	{"procedure", "CREATE PROCEDURE p IS CURSOR c IS SELECT id FROM customers; n NUMBER; BEGIN UPDATE orders SET status = 'x'; n := SQL%ROWCOUNT; IF SQL%NOTFOUND THEN NULL; END IF; OPEN c; IF c%ISOPEN THEN CLOSE c; END IF; END;"},
	{"procedure", "CREATE PROCEDURE p IS TYPE t IS TABLE OF NUMBER INDEX BY PLS_INTEGER; a t; BEGIN a(1) := 2; IF a.COUNT > 0 THEN a.DELETE; END IF; DBMS_OUTPUT.PUT_LINE(TO_CHAR(a.FIRST)); END;"},
	{"procedure", "CREATE PROCEDURE p IS v VARCHAR2(10); BEGIN v := q'[x]'; v := CASE WHEN v IS NULL THEN 'a' ELSE v END; EXECUTE IMMEDIATE 'SELECT 1 FROM dual' INTO v; v := SUBSTR(v, 1, 2) || USER; END;"},
	{"procedure", "CREATE PROCEDURE p (x OUT SYS_REFCURSOR) IS BEGIN OPEN x FOR SELECT * FROM orders o WHERE o.amount > 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS e EXCEPTION; PRAGMA EXCEPTION_INIT(e, -20001); BEGIN NULL; EXCEPTION WHEN e THEN DBMS_OUTPUT.PUT_LINE(SQLERRM); END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN FOR i IN REVERSE 1 .. 3 LOOP CONTINUE WHEN i = 2; EXIT WHEN i = 1; END LOOP; <<outer>> FOR r IN (SELECT id FROM orders) LOOP NULL; END LOOP outer; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM orders o WHERE EXISTS (SELECT 1 FROM customers c WHERE c.id = o.customer_id AND c.name = USER); v := v + p_existing_count; END;"},
	{"procedure", "CREATE PROCEDURE p IS CURSOR c (p_c NUMBER) IS SELECT id, amount FROM orders WHERE customer_id = p_c; r c%ROWTYPE; TYPE tt IS TABLE OF c%ROWTYPE; rs tt; BEGIN OPEN c(1); FETCH c BULK COLLECT INTO rs LIMIT 10; CLOSE c; FOR x IN c(2) LOOP r := x; END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; d DATE := DATE '2024-01-01'; BEGIN SELECT 1 INTO v FROM dual; INSERT INTO orders (id, amount) VALUES (order_seq.NEXTVAL, 1) RETURNING id INTO v; UPDATE orders SET amount = amount + 1 WHERE id = v RETURNING amount INTO v; d := TO_DATE('2024', 'YYYY') + NVL(v, 0); EXCEPTION WHEN OTHERS THEN ROLLBACK; RAISE; END;"},
	{"procedure", "CREATE PROCEDURE p IS PROCEDURE b; PROCEDURE a IS BEGIN b; END; PROCEDURE b IS BEGIN NULL; END; BEGIN a; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN DECLARE x NUMBER := 1; BEGIN x := x + 1; END; MERGE INTO orders o USING (SELECT id, name FROM customers) s ON (o.customer_id = s.id) WHEN MATCHED THEN UPDATE SET o.status = s.name WHEN NOT MATCHED THEN INSERT (id, customer_id) VALUES (s.id, s.id); END;"},
	{"procedure", "CREATE PROCEDURE p IS v VARCHAR2(100); BEGIN v := SYS_CONTEXT('USERENV', 'SESSION_USER'); v := pk.g_count; pk.g_count := 1; IF pk.total(1) > 0 THEN NULL; END IF; v := f_double(p_x => 2); END;"},
	{"procedure", "CREATE PROCEDURE p IS TYPE ids IS TABLE OF NUMBER; l ids := ids(1, 2, 3); BEGIN FORALL i IN 1 .. l.COUNT SAVE EXCEPTIONS DELETE FROM orders WHERE id = l(i); FOR i IN 1 .. l.COUNT LOOP DBMS_OUTPUT.PUT_LINE(l(i)); END LOOP; END;"},
	{"procedure", "CREATE PROCEDURE p IS p2 t_point := t_point(1, 2); n NUMBER; BEGIN n := p2.x + p2.norm; n := p2.norm(); END;"},
	{"procedure", "CREATE PROCEDURE p IS v customers.name%TYPE; c SYS_REFCURSOR; BEGIN OPEN c FOR 'SELECT name FROM customers'; LOOP FETCH c INTO v; EXIT WHEN c%NOTFOUND; END LOOP; CLOSE c; EXECUTE IMMEDIATE 'BEGIN NULL; END;'; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT o.amount INTO v FROM orders o JOIN customers c ON c.id = o.customer_id WHERE c.region = 'EU' AND ROWNUM = 1 ORDER BY o.id; SELECT MAX(amount) INTO v FROM orders GROUP BY customer_id HAVING COUNT(*) > 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN WITH x AS (SELECT id FROM orders) SELECT COUNT(*) INTO v FROM x; SELECT id INTO v FROM (SELECT id FROM orders WHERE amount > v) WHERE ROWNUM = 1; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM user_tables WHERE table_name = 'X'; SELECT COUNT(*) INTO v FROM all_objects; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN INSERT INTO audit_log SELECT id, name, SYSDATE FROM customers; DELETE FROM audit_log a WHERE a.id IN (SELECT id FROM customers); END;"},
	{"procedure", "CREATE PROCEDURE p (o OUT NUMBER, io IN OUT VARCHAR2) IS BEGIN o := 1; io := io || 'x'; p_existing(o); END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := f_double(f_double(1)); v := MOD(v, 2) + ABS(-1) + GREATEST(1, 2) + LENGTH('x') + INSTR('a', 'b'); END;"},
	{"procedure", "CREATE PROCEDURE p IS j JSON; x XMLTYPE; b BLOB; c CLOB := EMPTY_CLOB(); BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := 1; GOTO done; v := 2; <<done>> NULL; END;"},
	{"procedure", "CREATE PROCEDURE p IS v NUMBER; BEGIN v := 1; <<blk>> DECLARE w NUMBER; BEGIN w := blk.w; END blk; END;"},
	{"procedure", "CREATE PROCEDURE p IS BEGIN COMMIT WORK; ROLLBACK TO SAVEPOINT s; LOCK TABLE orders IN EXCLUSIVE MODE; END;"},
	{"procedure", "CREATE PROCEDURE p IS x NUMBER; BEGIN x := g_missing.y; END;"},
	// Procedures: names, replacement, drops.
	{"procedure", "CREATE PROCEDURE p_existing IS BEGIN NULL; END;"},
	{"procedure", "CREATE OR REPLACE PROCEDURE p_existing (p_id IN NUMBER) IS BEGIN NULL; END;"},
	{"procedure", "CREATE OR REPLACE PROCEDURE p_existing (p_id IN NUMBER, p_extra IN NUMBER) IS BEGIN NULL; END;"},
	{"procedure", "CREATE OR REPLACE PROCEDURE p_existing IS BEGIN NULL; END;"},
	{"procedure", "CREATE OR REPLACE PROCEDURE customers IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE cust IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE f_double IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE pk IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE omni_sim.p IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE nope.p IS BEGIN NULL; END;"},
	{"procedure", "CREATE PROCEDURE IF NOT EXISTS p_existing IS BEGIN NULL; END;"},
	{"procedure", "DROP PROCEDURE p_existing"},
	{"procedure", "DROP PROCEDURE nope"},
	{"procedure", "DROP PROCEDURE f_double"},
	{"procedure", "DROP PROCEDURE IF EXISTS nope"},
	{"procedure", "ALTER PROCEDURE p_existing COMPILE"},
	{"procedure", "ALTER PROCEDURE nope COMPILE"},
	{"procedure", "ALTER PROCEDURE p_calls COMPILE"},

	// Functions.
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END;"},
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN NULL; END;"},
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN; END;"},
	{"function", "CREATE FUNCTION f (p_c IN NUMBER) RETURN NUMBER IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM orders WHERE customer_id = p_c; RETURN v; END;"},
	{"function", "CREATE FUNCTION f (p_c IN NUMBER) RETURN NUMBER IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM orders WHERE nope = p_c; RETURN v; END;"},
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN f_double(1) + pk.total(1); END;"},
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN nope_f(1); END;"},
	{"function", "CREATE FUNCTION f RETURN customers%ROWTYPE IS r customers%ROWTYPE; BEGIN SELECT * INTO r FROM customers WHERE ROWNUM = 1; RETURN r; END;"},
	{"function", "CREATE FUNCTION f RETURN missing_t%ROWTYPE IS BEGIN RETURN NULL; END;"},
	{"function", "CREATE FUNCTION f RETURN nope_type IS BEGIN RETURN NULL; END;"},
	{"function", "CREATE FUNCTION f RETURN t_point IS BEGIN RETURN t_point(1, 2); END;"},
	{"function", "CREATE FUNCTION f RETURN VARCHAR2 DETERMINISTIC IS BEGIN RETURN 'x'; END;"},
	{"function", "CREATE FUNCTION f (n IN NUMBER) RETURN NUMBER RESULT_CACHE IS BEGIN RETURN n; END;"},
	{"function", "CREATE FUNCTION f (n IN NUMBER) RETURN NUMBER IS BEGIN IF n <= 1 THEN RETURN 1; END IF; RETURN n * f(n - 1); END;"},
	{"function", "CREATE FUNCTION f RETURN NUMBER IS BEGIN RETURN 'x'; END;"},
	{"function", "CREATE FUNCTION f RETURN DATE IS BEGIN RETURN 1; END;"},
	{"function", "CREATE FUNCTION f RETURN BOOLEAN IS BEGIN RETURN TRUE; END;"},
	{"function", "CREATE FUNCTION f_double (p_x IN NUMBER) RETURN NUMBER IS BEGIN RETURN p_x; END;"},
	{"function", "CREATE OR REPLACE FUNCTION f_double (p_x IN NUMBER) RETURN NUMBER IS BEGIN RETURN p_x + p_x; END;"},
	{"function", "CREATE OR REPLACE FUNCTION f_double RETURN NUMBER IS BEGIN RETURN 2; END;"},
	{"function", "CREATE OR REPLACE FUNCTION p_existing RETURN NUMBER IS BEGIN RETURN 1; END;"},
	{"function", "CREATE FUNCTION orders RETURN NUMBER IS BEGIN RETURN 1; END;"},
	{"function", "DROP FUNCTION f_double"},
	{"function", "DROP FUNCTION nope"},
	{"function", "DROP FUNCTION p_existing"},

	// Package specs.
	{"package", "CREATE PACKAGE p2 AS PROCEDURE a; END;"},
	{"package", "CREATE PACKAGE p2 AS PROCEDURE a (x IN nope_type); END;"},
	{"package", "CREATE PACKAGE p2 AS v customers.nope%TYPE; END;"},
	{"package", "CREATE PACKAGE p2 AS v customers.name%TYPE; c CONSTANT NUMBER := 1; TYPE t_rec IS RECORD (a NUMBER, b VARCHAR2(10)); TYPE t_tab IS TABLE OF t_rec; END;"},
	{"package", "CREATE PACKAGE p2 AS CURSOR c RETURN customers%ROWTYPE; END;"},
	{"package", "CREATE PACKAGE p2 AS FUNCTION f RETURN NUMBER; e EXCEPTION; END p2;"},
	{"package", "CREATE PACKAGE p2 AS PROCEDURE a; END nope;"},
	{"package", "CREATE PACKAGE pk AS PROCEDURE a; END;"},
	{"package", "CREATE OR REPLACE PACKAGE pk AS g_count NUMBER := 0; PROCEDURE log_it (p_what IN VARCHAR2); FUNCTION total (p_customer IN NUMBER) RETURN NUMBER; END pk;"},
	{"package", "CREATE OR REPLACE PACKAGE pk AS PROCEDURE log_it (p_what IN VARCHAR2); END pk;"},
	{"package", "CREATE OR REPLACE PACKAGE pk AS g_count NUMBER := 0; PROCEDURE log_it (p_what IN VARCHAR2); FUNCTION total (p_customer IN NUMBER) RETURN NUMBER; PROCEDURE extra; END pk;"},
	{"package", "CREATE OR REPLACE PACKAGE orders AS PROCEDURE a; END;"},
	{"package", "DROP PACKAGE pk"},
	{"package", "DROP PACKAGE nope"},
	{"package", "ALTER PACKAGE pk COMPILE BODY"},

	{"package", "CREATE PACKAGE p2 AS TYPE rec IS RECORD (id customers.id%TYPE, n customers.name%TYPE); TYPE recs IS TABLE OF rec INDEX BY PLS_INTEGER; c_max CONSTANT PLS_INTEGER := 10; FUNCTION get (p_id customers.id%TYPE) RETURN rec; PROCEDURE put (p IN rec); PROCEDURE put (p IN recs); END;"},

	// Package bodies.
	{"package-body", "CREATE PACKAGE BODY nope AS PROCEDURE a IS BEGIN NULL; END; END;"},
	{"package-body", "CREATE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN NUMBER) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_other IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN VARCHAR2 IS BEGIN RETURN 'x'; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE helper IS BEGIN NULL; END; PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN helper; g_count := g_count + 1; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN g_count; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN INSERT INTO audit_log (nope) VALUES (p_what); END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS v NUMBER; BEGIN SELECT SUM(amount) INTO v FROM missing_t; RETURN v; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN NULL; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; BEGIN g_count := 10; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN nope; END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"package-body", "CREATE OR REPLACE PACKAGE BODY pk AS g_private NUMBER; PROCEDURE log_it (p_what IN VARCHAR2) IS BEGIN g_private := 1; g_count := g_count + g_private; INSERT INTO audit_log (id, what, at) VALUES (g_count, p_what, SYSTIMESTAMP); END; FUNCTION total (p_customer IN NUMBER) RETURN NUMBER IS BEGIN RETURN pk.g_count + g_private; END; END;"},
	{"package-body", "DROP PACKAGE BODY pk"},
	{"package-body", "DROP PACKAGE BODY nope"},

	// Triggers.
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN :NEW.amount := NVL(:NEW.amount, 0); END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN :NEW.nope := 0; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN :NEW.amount := :OLD.amount; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON missing_t FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg AFTER UPDATE OF status ON orders FOR EACH ROW BEGIN INSERT INTO audit_log (id, what, at) VALUES (:NEW.id, :NEW.status, SYSDATE); END;"},
	{"trigger", "CREATE TRIGGER trg AFTER UPDATE OF nope ON orders FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg AFTER INSERT ON orders BEGIN INSERT INTO audit_log (id) VALUES (:NEW.id); END;"},
	{"trigger", "CREATE TRIGGER trg AFTER INSERT ON orders BEGIN pk.log_it('insert'); END;"},
	{"trigger", "CREATE TRIGGER trg AFTER INSERT ON orders FOR EACH ROW BEGIN :NEW.amount := 0; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE DELETE ON orders FOR EACH ROW BEGIN :NEW.amount := 0; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW WHEN (NEW.amount > 100) BEGIN :NEW.status := 'big'; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders WHEN (NEW.amount > 100) BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders REFERENCING NEW AS n FOR EACH ROW BEGIN :n.status := 'x'; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders REFERENCING NEW AS n FOR EACH ROW BEGIN :NEW.status := 'x'; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT OR UPDATE OR DELETE ON orders FOR EACH ROW BEGIN IF INSERTING THEN :NEW.status := 'i'; ELSIF UPDATING('status') THEN NULL; END IF; END;"},
	{"trigger", "CREATE TRIGGER trg INSTEAD OF INSERT ON v_customers FOR EACH ROW BEGIN INSERT INTO customers (id, name) VALUES (:NEW.id, :NEW.name); END;"},
	{"trigger", "CREATE TRIGGER trg INSTEAD OF INSERT ON customers FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON v_customers FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN INSERT INTO audit_log (nope) VALUES (1); END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN nope_proc; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW BEGIN :NEW.id := order_seq.NEXTVAL; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON cust FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT ON orders FOR EACH ROW DISABLE BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER trg FOR INSERT ON orders COMPOUND TRIGGER BEFORE EACH ROW IS BEGIN :NEW.status := 'c'; END BEFORE EACH ROW; END trg;"},
	{"trigger", "CREATE TRIGGER trg BEFORE UPDATE ON orders FOR EACH ROW BEGIN IF UPDATING('STATUS') AND :OLD.status <> :NEW.status THEN :NEW.amount := :OLD.amount; END IF; END;"},
	{"trigger", "CREATE TRIGGER trg BEFORE INSERT OR UPDATE ON orders FOR EACH ROW WHEN (NEW.amount IS NULL) DECLARE v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM customers WHERE id = :NEW.customer_id; IF v = 0 THEN RAISE_APPLICATION_ERROR(-20000, 'no customer ' || :NEW.customer_id); END IF; END;"},
	{"trigger", "CREATE TRIGGER trg INSTEAD OF UPDATE ON v_customers FOR EACH ROW BEGIN UPDATE customers SET name = :NEW.name WHERE id = :OLD.id; END;"},
	{"trigger", "CREATE TRIGGER trg AFTER INSERT OR DELETE ON orders BEGIN pk.log_it('stmt'); END;"},
	{"trigger", "CREATE TRIGGER trg_existing BEFORE INSERT ON orders FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE OR REPLACE TRIGGER trg_existing BEFORE INSERT ON orders FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE OR REPLACE TRIGGER trg_existing BEFORE INSERT ON customers FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "CREATE TRIGGER p_existing BEFORE INSERT ON orders FOR EACH ROW BEGIN NULL; END;"},
	{"trigger", "DROP TRIGGER trg_existing"},
	{"trigger", "DROP TRIGGER nope"},
	{"trigger", "ALTER TRIGGER trg_existing DISABLE"},
	{"trigger", "ALTER TRIGGER trg_existing COMPILE"},
	{"trigger", "ALTER TRIGGER nope ENABLE"},

	// Types.
	{"type", "CREATE TYPE t AS OBJECT (a NUMBER, b VARCHAR2(10))"},
	{"type", "CREATE TYPE t AS OBJECT (a nope_type)"},
	{"type", "CREATE TYPE t AS OBJECT (a t_point)"},
	{"type", "CREATE TYPE t AS OBJECT (a VARCHAR2)"},
	{"type", "CREATE TYPE t AS OBJECT (a NUMBER, MEMBER FUNCTION f RETURN NUMBER)"},
	{"type", "CREATE TYPE t AS OBJECT (a NUMBER) NOT FINAL"},
	{"type", "CREATE TYPE t AS TABLE OF NUMBER"},
	{"type", "CREATE TYPE t AS TABLE OF t_point"},
	{"type", "CREATE TYPE t AS TABLE OF nope_type"},
	{"type", "CREATE TYPE t AS VARRAY(10) OF VARCHAR2(10)"},
	{"type", "CREATE TYPE t UNDER t_point (z NUMBER)"},
	{"type", "CREATE TYPE t_point AS OBJECT (x NUMBER)"},
	{"type", "CREATE OR REPLACE TYPE t_point AS OBJECT (x NUMBER, y NUMBER, MEMBER FUNCTION norm RETURN NUMBER)"},
	{"type", "CREATE OR REPLACE TYPE t_point AS OBJECT (x NUMBER)"},
	{"type", "CREATE TYPE customers AS OBJECT (a NUMBER)"},
	{"type", "DROP TYPE t_point"},
	{"type", "DROP TYPE t_point FORCE"},
	{"type", "DROP TYPE nope"},

	// Type bodies.
	{"type-body", "CREATE TYPE BODY nope AS MEMBER FUNCTION f RETURN NUMBER IS BEGIN RETURN 1; END; END;"},
	{"type-body", "CREATE OR REPLACE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS BEGIN RETURN x; END; END;"},
	{"type-body", "CREATE OR REPLACE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS BEGIN RETURN nope; END; END;"},
	{"type-body", "CREATE OR REPLACE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS BEGIN RETURN SELF.x; END; MEMBER FUNCTION extra RETURN NUMBER IS BEGIN RETURN 1; END; END;"},
	{"type-body", "CREATE OR REPLACE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS v NUMBER; BEGIN SELECT COUNT(*) INTO v FROM missing_t; RETURN v; END; END;"},
	{"type-body", "CREATE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS BEGIN RETURN 0; END; END;"},
	{"type-body", "CREATE OR REPLACE TYPE BODY t_point AS MEMBER FUNCTION norm RETURN NUMBER IS BEGIN RETURN SQRT(SELF.x * SELF.x + y * y); END; END;"},
	{"type-body", "DROP TYPE BODY t_point"},
	{"type-body", "DROP TYPE BODY nope"},

	// What a change to a referenced object does to the units that read it.
	{"reference", "DROP TABLE audit_log"},
	{"reference", "DROP TABLE customers CASCADE CONSTRAINTS"},
	{"reference", "ALTER TABLE audit_log ADD (extra NUMBER)"},
	{"reference", "ALTER TABLE audit_log DROP COLUMN at"},
	{"reference", "ALTER TABLE customers DROP COLUMN region"},
	{"reference", "ALTER TABLE customers MODIFY (name VARCHAR2(200))"},
	{"reference", "RENAME audit_log TO audit_log2"},
	{"reference", "DROP SEQUENCE order_seq"},
	{"reference", "DROP SYNONYM cust"},
	{"reference", "DROP VIEW v_customers"},
	{"reference", "CREATE TABLE pts (id NUMBER, p t_point)"},
	{"reference", "CREATE TABLE pts (id NUMBER, p nope_type)"},
}

// catalogUnitForms renders the catalog's PL/SQL units as programUnits reads
// Oracle's.
func catalogUnitForms(c *Catalog) []simulate.Form {
	var forms []simulate.Form
	for _, u := range c.unitStates() {
		forms = append(forms, programUnitForm(u))
	}
	return forms
}

// withoutPrograms is meta with its PL/SQL units left out, as dump_forms
// reads a schema.
func withoutPrograms(meta *metadatapb.DatabaseSchemaMetadata) *metadatapb.DatabaseSchemaMetadata {
	out := cloneMetadata(meta)
	for _, s := range out.GetSchemas() {
		s.Procedures, s.Functions, s.Packages = nil, nil, nil
		for _, t := range s.GetTables() {
			t.Triggers = nil
		}
		for _, v := range s.GetViews() {
			v.Triggers = nil
		}
		for _, v := range s.GetMaterializedViews() {
			v.Triggers = nil
		}
	}
	return out
}
