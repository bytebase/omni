package catalog

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestProgramSimulation holds the catalog to testdata/programs/simulate.tsv:
// SQL and PL/pgSQL functions and procedures, with real bodies, their ALTER
// and DROP, the dependencies a SQL-standard body records, triggers, event
// triggers, CALL and DO, in programFeatures, run on programs@objects, with
// what PostgreSQL did and what Exec does. testdata/simulate.tsv runs the
// grammar's program statements, whose bodies name nothing; these hold the
// catalog to what PostgreSQL checks when it creates a program
// (check_function_bodies) and to what it records the program depends on.
//
// PostgreSQL's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned PostgreSQL again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestProgramSimulation(t *testing.T) {
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     programRows(t),
		MadeUp:   nil,
		Omni:     simulationOmni,
		Engine:   simulationEngineWithQuery(programFormsQuery),
	})
}

// TestProgramSimulationPG18 holds the catalog to the same rows on the
// pinned PostgreSQL 18, in testdata/pg18/programs/simulate.tsv.
func TestProgramSimulationPG18(t *testing.T) {
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "pg18", "programs", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		PinRole:  testcontainer.PG18,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     programRows(t),
		MadeUp:   nil,
		Omni:     simulationOmniWithVersion(PG18),
		Engine:   simulationEngineWithRole(testcontainer.PG18, programFormsQuery),
	})
}

func programRows(t *testing.T) []simulate.Generated {
	t.Helper()
	rows := make([]simulate.Generated, 0, len(programFeatures)+len(quotedTypeFeatures)+len(quotedAliasFeatures)+len(retypeFeatures)+len(generatedRetypeFeatures))
	seen := map[string]bool{}
	add := func(fixture string, features []struct{ kind, statement string }) {
		for _, f := range features {
			if seen[f.statement] {
				t.Fatalf("the program rows hold %q twice", f.statement)
			}
			seen[f.statement] = true
			rows = append(rows, simulate.Generated{Fixture: fixture, Rule: "program:" + f.kind, Statement: f.statement})
		}
	}
	add("programs@objects", programFeatures)
	add("programs@quoted", quotedTypeFeatures)
	add("programs@quoted", quotedAliasFeatures)
	add("programs@retype", retypeFeatures)
	add("programs@generated", generatedRetypeFeatures)
	return rows
}

// programFormsQuery is simulationFormsQuery with the database's event
// triggers, which live in no schema. Dropping public drops them with the
// functions they run, so a row starts with none but the fixture's.
var programFormsQuery = func() string {
	const end = "\n) o ORDER BY relname COLLATE \"C\", kind"
	if !strings.HasSuffix(simulationFormsQuery, end) {
		panic("simulationFormsQuery no longer ends with " + end)
	}
	return strings.TrimSuffix(simulationFormsQuery, end) + `
	UNION ALL
	SELECT 'EVENT TRIGGER', e.evtname,
		format(E'CREATE EVENT TRIGGER %I ON %I%s EXECUTE FUNCTION %s();\n', e.evtname, e.evtevent,
			coalesce((SELECT ' WHEN TAG IN (' || string_agg(format('%L', t.tag), ', ' ORDER BY t.n) || ')'
				FROM unnest(e.evttags) WITH ORDINALITY AS t(tag, n)), ''),
			e.evtfoid::pg_catalog.regproc)
		|| CASE e.evtenabled
			WHEN 'D' THEN format(E'ALTER EVENT TRIGGER %I DISABLE;\n', e.evtname)
			WHEN 'R' THEN format(E'ALTER EVENT TRIGGER %I ENABLE REPLICA;\n', e.evtname)
			WHEN 'A' THEN format(E'ALTER EVENT TRIGGER %I ENABLE ALWAYS;\n', e.evtname)
			ELSE '' END
		|| coalesce((SELECT format(E'COMMENT ON EVENT TRIGGER %I IS %L;\n', e.evtname, pg_catalog.obj_description(e.oid, 'pg_event_trigger'))
			WHERE pg_catalog.obj_description(e.oid, 'pg_event_trigger') IS NOT NULL), '')
		FROM pg_catalog.pg_event_trigger e` + end
}()

// programFeatures are the statements the suite runs on programs@objects,
// labelled program:<kind>. Many are refused by PostgreSQL on purpose.
var programFeatures = []struct{ kind, statement string }{
	// SQL functions with a string body: check_function_bodies parses and
	// analyzes every statement of the body at CREATE.
	{"function-sql", "CREATE FUNCTION f(c int) RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM orders WHERE customer_id = c'"},
	{"function-sql", "CREATE FUNCTION f(c int) RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM orders WHERE customer_id = $1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM missing_t'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS text LANGUAGE sql AS 'SELECT nope FROM customers'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT name FROM customers'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELEC 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS 'CREATE TABLE made_in_f (x int); INSERT INTO made_in_f VALUES (1)'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS 'CREATE TABLE made_in_f (x int)'"},
	{"function-sql", "CREATE FUNCTION f(x anyelement) RETURNS anyelement LANGUAGE sql AS 'SELECT nope FROM missing_t'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS SETOF orders LANGUAGE sql AS 'SELECT * FROM orders'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS TABLE (id int, name text) LANGUAGE sql AS 'SELECT id, name FROM customers'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS TABLE (id int, name text) LANGUAGE sql AS 'SELECT id FROM customers'"},
	{"function-sql", "CREATE FUNCTION f(INOUT a int, OUT b text) LANGUAGE sql AS 'SELECT a, ''x'''"},
	{"function-sql", "CREATE FUNCTION f(OUT a int, OUT b int) LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f(VARIADIC xs int[]) RETURNS int LANGUAGE sql AS 'SELECT array_length(xs, 1)'"},
	{"function-sql", "CREATE FUNCTION f(a int, b int DEFAULT 1) RETURNS int LANGUAGE sql AS 'SELECT a + b'"},
	{"function-sql", "CREATE FUNCTION f(a int DEFAULT 1, b int) RETURNS int LANGUAGE sql AS 'SELECT a + b'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql STABLE STRICT PARALLEL SAFE COST 10 SECURITY DEFINER SET search_path = public AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS trigger LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS event_trigger LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT 1; SELECT 2'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'INSERT INTO logs (msg) VALUES (''x'')'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'INSERT INTO logs (id) VALUES (1) RETURNING id'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS 'UPDATE orders SET nope = 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS 'DELETE FROM missing_t'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS 'SELECT; SELECT nope FROM missing_t'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT f_double(1)'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT nope_fn(1)'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT f_double(1, 2)'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT $1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS 'SELECT id FROM customers WHERE region = nope'"},
	{"function-sql", "CREATE FUNCTION f(c int) RETURNS SETOF orders LANGUAGE sql AS 'SELECT * FROM orders WHERE customer_id = c ORDER BY nope'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS SETOF text LANGUAGE sql AS 'SELECT name FROM v_customers'"},
	{"function-sql", "CREATE FUNCTION f(m mood) RETURNS boolean LANGUAGE sql AS 'SELECT m = ''nope'''"},
	{"function-sql", "CREATE FUNCTION f() RETURNS record LANGUAGE sql AS 'SELECT 1, 2'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS ''"},
	{"function-sql", "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS ''"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS $$ SELECT 1 $$"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE internal AS 'nope_symbol'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int LANGUAGE nope AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS nope LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f(x nope) RETURNS int LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f() RETURNS int AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION nope_schema.f() RETURNS int LANGUAGE sql AS 'SELECT 1'"},
	{"function-sql", "CREATE FUNCTION f_double(x int) RETURNS int LANGUAGE sql AS 'SELECT x'"},
	{"function-sql", "CREATE FUNCTION f_double(x bigint) RETURNS bigint LANGUAGE sql AS 'SELECT x * 2'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION f_double(x int) RETURNS int LANGUAGE sql AS 'SELECT x * 3'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION f_double(x int) RETURNS bigint LANGUAGE sql AS 'SELECT x::bigint'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION f_double(y int) RETURNS int LANGUAGE sql AS 'SELECT y'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION f_double(x int) RETURNS int LANGUAGE sql AS 'SELECT nope FROM missing_t'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION p_log(m text) RETURNS void LANGUAGE sql AS 'SELECT'"},
	{"function-sql", "CREATE OR REPLACE FUNCTION f_log_count() RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM audit'"},

	// SQL-standard bodies (RETURN, BEGIN ATOMIC) are parsed and analyzed at
	// CREATE, and record what they use.
	{"function-atomic", "CREATE FUNCTION f(i int) RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT qty FROM items WHERE id = i; END"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT qty FROM missing_t; END"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT nope FROM items; END"},
	{"function-atomic", "CREATE FUNCTION f(x int) RETURNS int RETURN x * 2"},
	{"function-atomic", "CREATE FUNCTION f(x int) RETURNS int RETURN nope_fn(x)"},
	{"function-atomic", "CREATE FUNCTION f(x int) RETURNS int RETURN nope"},
	{"function-atomic", "CREATE FUNCTION f(x int) RETURNS int LANGUAGE plpgsql RETURN x"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC END"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS void LANGUAGE sql BEGIN ATOMIC END"},
	{"function-atomic", "CREATE FUNCTION f(x anyelement) RETURNS anyelement BEGIN ATOMIC SELECT x; END"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS bigint LANGUAGE sql BEGIN ATOMIC INSERT INTO logs (msg) VALUES ('a'); SELECT count(*) FROM logs; END"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS text RETURN (SELECT name FROM v_customers LIMIT 1)"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS text RETURN (SELECT name FROM v_customers WHERE nope)"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS int RETURN f_quad(1)"},
	{"function-atomic", "CREATE FUNCTION f(m mood) RETURNS boolean RETURN m = 'happy'"},
	{"function-atomic", "CREATE FUNCTION f(m mood) RETURNS boolean RETURN m = 'nope'"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS int RETURN 'a'"},
	{"function-atomic", "CREATE FUNCTION f() RETURNS SETOF int LANGUAGE sql BEGIN ATOMIC SELECT id FROM orders; END"},
	{"function-atomic", "CREATE OR REPLACE FUNCTION f_item_qty(i int) RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT qty * 2 FROM items WHERE id = i; END"},

	// Procedures.
	{"procedure", "CREATE PROCEDURE p(m text) LANGUAGE sql AS 'INSERT INTO logs (msg) VALUES (m)'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql AS 'INSERT INTO missing_t VALUES (1)'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql AS 'SELECT nope FROM logs'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE plpgsql AS $$ BEGIN INSERT INTO missing_t VALUES (1); END $$"},
	{"procedure", "CREATE PROCEDURE p(m text) BEGIN ATOMIC INSERT INTO logs (msg) VALUES (m); END"},
	{"procedure", "CREATE PROCEDURE p() BEGIN ATOMIC INSERT INTO missing_t VALUES (1); END"},
	{"procedure", "CREATE PROCEDURE p(OUT n int) LANGUAGE plpgsql AS $$ BEGIN n := 1; END $$"},
	{"procedure", "CREATE PROCEDURE p() RETURNS int LANGUAGE sql AS 'SELECT 1'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE plpgsql AS $$ BEGIN RETURN 1; END $$"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE plpgsql AS $$ BEGIN COMMIT; END $$"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql STRICT AS 'SELECT'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql SECURITY DEFINER SET search_path = public AS 'SELECT'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql AS 'CALL p_log(''x'')'"},
	{"procedure", "CREATE PROCEDURE p() LANGUAGE sql AS 'CALL nope()'"},
	{"procedure", "CREATE PROCEDURE p_log(m text) LANGUAGE sql AS 'SELECT'"},
	{"procedure", "CREATE OR REPLACE PROCEDURE p_log(m text) LANGUAGE plpgsql AS $$ BEGIN NULL; END $$"},
	{"procedure", "CREATE OR REPLACE PROCEDURE f_double(x int) LANGUAGE sql AS 'SELECT'"},

	// PL/pgSQL: the validator compiles the body, which checks its syntax,
	// variables and declared types, and leaves the tables it names to run
	// time.
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS bigint LANGUAGE plpgsql AS $$ BEGIN RETURN (SELECT count(*) FROM missing_t); END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x int; BEGIN SELECT id INTO x FROM missing_t; RETURN x; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1; $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN IF true THEN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN x := 1; RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN SELECT 1 INTO nope; RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 + ; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN UPDATE orders SET WHERE; RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x nope_type; BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE r orders%ROWTYPE; BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE r missing_t%ROWTYPE; BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x orders.amount%TYPE; BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x orders.nope%TYPE; BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE c CONSTANT int := 1; BEGIN c := 2; RETURN c; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS void LANGUAGE plpgsql AS $$ BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN NEXT 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS SETOF int LANGUAGE plpgsql AS $$ BEGIN RETURN NEXT 1; RETURN; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS SETOF int LANGUAGE plpgsql AS $$ BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f(OUT a int) LANGUAGE plpgsql AS $$ BEGIN RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN PERFORM nope_fn(); RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN EXECUTE 'SELECT nope FROM missing_t'; RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE r record; BEGIN FOR r IN SELECT id FROM customers LOOP NULL; END LOOP; RETURN 1; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'boom'; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RAISE nope_condition; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1; EXCEPTION WHEN unique_violation THEN RETURN 0; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1; EXCEPTION WHEN nope_condition THEN RETURN 0; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS 'BEGIN RETURN 1; END'"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RETURN NEW; END $$"},
	{"function-plpgsql", "CREATE FUNCTION f() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN RAISE NOTICE 'e'; END $$"},
	{"function-plpgsql", "CREATE OR REPLACE FUNCTION f_total(c int) RETURNS numeric LANGUAGE plpgsql AS $$ BEGIN RETURN 0; END $$"},
	{"function-plpgsql", "CREATE OR REPLACE FUNCTION f_total(c int) RETURNS numeric LANGUAGE plpgsql AS $$ BEGIN RETURN 0 END $$"},

	// ALTER FUNCTION, PROCEDURE and ROUTINE.
	{"alter-routine", "ALTER FUNCTION f_double(int) RENAME TO f_twice"},
	{"alter-routine", "ALTER FUNCTION f_double(int) RENAME TO f_quad"},
	{"alter-routine", "ALTER FUNCTION nope(int) RENAME TO f_twice"},
	{"alter-routine", "ALTER FUNCTION f_double STABLE"},
	{"alter-routine", "ALTER FUNCTION f_double(int) STRICT"},
	{"alter-routine", "ALTER FUNCTION f_double(int) SET search_path = public"},
	{"alter-routine", "ALTER FUNCTION f_double(int) RESET ALL"},
	{"alter-routine", "ALTER FUNCTION f_double(int) SET SCHEMA nope_schema"},
	{"alter-routine", "ALTER FUNCTION f_double(int) OWNER TO nope_role"},
	{"alter-routine", "ALTER FUNCTION f_double(bigint) STABLE"},
	{"alter-routine", "ALTER FUNCTION f_over STABLE"},
	{"alter-routine", "ALTER FUNCTION f_over(text) STABLE"},
	{"alter-routine", "ALTER FUNCTION p_log(text) STABLE"},
	{"alter-routine", "ALTER FUNCTION f_item_qty(int) RENAME TO f_qty"},
	{"alter-routine", "ALTER PROCEDURE p_log(text) RENAME TO p_write"},
	{"alter-routine", "ALTER PROCEDURE f_double(int) RENAME TO p_write"},
	{"alter-routine", "ALTER PROCEDURE p_log(text) SECURITY DEFINER"},
	{"alter-routine", "ALTER PROCEDURE p_log(text) STRICT"},
	{"alter-routine", "ALTER ROUTINE f_double(int) COST 5"},
	{"alter-routine", "ALTER ROUTINE p_log RENAME TO p_write"},
	{"alter-routine", "COMMENT ON FUNCTION f_double(int) IS 'doubles'"},
	{"alter-routine", "COMMENT ON PROCEDURE p_log(text) IS 'logs'"},
	{"alter-routine", "COMMENT ON FUNCTION nope() IS 'x'"},

	// DROP FUNCTION, PROCEDURE and ROUTINE.
	{"drop-routine", "DROP FUNCTION f_double(int)"},
	{"drop-routine", "DROP FUNCTION f_double(int) CASCADE"},
	{"drop-routine", "DROP FUNCTION f_total(int)"},
	{"drop-routine", "DROP FUNCTION f_total"},
	{"drop-routine", "DROP FUNCTION f_total(text)"},
	{"drop-routine", "DROP FUNCTION f_over"},
	{"drop-routine", "DROP FUNCTION f_over(text)"},
	{"drop-routine", "DROP FUNCTION nope()"},
	{"drop-routine", "DROP FUNCTION IF EXISTS nope()"},
	{"drop-routine", "DROP FUNCTION f_total(int), nope()"},
	{"drop-routine", "DROP FUNCTION trg_status()"},
	{"drop-routine", "DROP FUNCTION trg_status() CASCADE"},
	{"drop-routine", "DROP FUNCTION evt_note()"},
	{"drop-routine", "DROP FUNCTION evt_note() CASCADE"},
	{"drop-routine", "DROP FUNCTION p_log(text)"},
	{"drop-routine", "DROP PROCEDURE p_log(text)"},
	{"drop-routine", "DROP PROCEDURE f_total(int)"},
	{"drop-routine", "DROP PROCEDURE IF EXISTS nope(int)"},
	{"drop-routine", "DROP ROUTINE p_log, f_total"},
	{"drop-routine", "DROP FUNCTION f_log_count()"},
	{"drop-routine", "DROP FUNCTION f_item_qty(int)"},
	{"drop-routine", "DROP FUNCTION f_quad(int)"},
	{"drop-routine", "DROP FUNCTION f_mood(mood)"},

	// What a body depends on: a SQL-standard body records the tables,
	// columns, views, types and routines it uses; a string body records
	// none of them.
	{"dependency", "DROP TABLE items"},
	{"dependency", "DROP TABLE items CASCADE"},
	{"dependency", "DROP TABLE logs"},
	{"dependency", "DROP TABLE audit"},
	{"dependency", "DROP TABLE orders"},
	{"dependency", "DROP TABLE customers"},
	{"dependency", "ALTER TABLE items DROP COLUMN qty"},
	{"dependency", "ALTER TABLE items DROP COLUMN note"},
	{"dependency", "ALTER TABLE items DROP COLUMN qty CASCADE"},
	{"dependency", "ALTER TABLE items RENAME COLUMN qty TO quantity"},
	{"dependency", "ALTER TABLE items RENAME TO goods"},
	{"dependency", "ALTER TABLE items ALTER COLUMN qty TYPE bigint"},
	{"dependency", "ALTER TABLE logs DROP COLUMN msg"},
	{"dependency", "ALTER TABLE audit DROP COLUMN msg"},
	{"dependency", "DROP VIEW v_customers"},
	{"dependency", "DROP VIEW v_customers CASCADE"},
	{"dependency", "DROP TYPE mood"},
	{"dependency", "DROP TYPE mood CASCADE"},

	// Triggers.
	{"trigger", "CREATE TRIGGER t BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER INSERT OR UPDATE OR DELETE ON orders FOR EACH STATEMENT EXECUTE PROCEDURE trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE UPDATE OF status ON orders FOR EACH ROW WHEN (OLD.status IS DISTINCT FROM NEW.status) EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW WHEN (OLD.status IS NULL) EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE DELETE ON orders FOR EACH ROW WHEN (NEW.id > 0) EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW WHEN (NEW.nope > 0) EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER INSERT ON orders FOR EACH STATEMENT WHEN (NEW.id > 0) EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE UPDATE OF nope ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON missing_t FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION nope()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION evt_note()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION f_log_count()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION trg_status('a', 'b')"},
	{"trigger", "CREATE TRIGGER orders_status BEFORE INSERT ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE OR REPLACE TRIGGER orders_status AFTER INSERT ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t INSTEAD OF INSERT ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t INSTEAD OF INSERT ON v_customers FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON v_customers FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t INSTEAD OF INSERT ON v_customers FOR EACH STATEMENT EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER INSERT ON orders REFERENCING NEW TABLE AS nt FOR EACH STATEMENT EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t BEFORE INSERT ON orders REFERENCING NEW TABLE AS nt FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER DELETE ON orders REFERENCING NEW TABLE AS nt FOR EACH STATEMENT EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE CONSTRAINT TRIGGER t AFTER INSERT ON orders DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER TRUNCATE ON orders FOR EACH STATEMENT EXECUTE FUNCTION trg_status()"},
	{"trigger", "CREATE TRIGGER t AFTER TRUNCATE ON orders FOR EACH ROW EXECUTE FUNCTION trg_status()"},
	{"trigger", "ALTER TRIGGER orders_status ON orders RENAME TO orders_status2"},
	{"trigger", "ALTER TRIGGER nope ON orders RENAME TO t"},
	{"trigger", "ALTER TRIGGER orders_status ON customers RENAME TO t"},
	{"trigger", "DROP TRIGGER orders_status ON orders"},
	{"trigger", "DROP TRIGGER nope ON orders"},
	{"trigger", "DROP TRIGGER IF EXISTS nope ON orders"},
	{"trigger", "DROP TRIGGER IF EXISTS nope ON missing_t"},
	{"trigger", "ALTER TABLE orders DISABLE TRIGGER orders_status"},
	{"trigger", "ALTER TABLE orders ENABLE ALWAYS TRIGGER orders_status"},
	{"trigger", "ALTER TABLE orders DISABLE TRIGGER nope"},
	{"trigger", "COMMENT ON TRIGGER orders_status ON orders IS 'status'"},

	// Event triggers.
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_start EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON sql_drop EXECUTE PROCEDURE evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON table_rewrite EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_end WHEN TAG IN ('CREATE TABLE', 'ALTER TABLE') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON sql_drop WHEN TAG IN ('DROP TABLE') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_end WHEN TAG IN ('NOPE') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_end WHEN TAG IN ('CREATE EVENT TRIGGER') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON table_rewrite WHEN TAG IN ('CREATE TABLE') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_end WHEN nope IN ('CREATE TABLE') EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON nope_event EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_start EXECUTE FUNCTION trg_status()"},
	{"event-trigger", "CREATE EVENT TRIGGER e ON ddl_command_start EXECUTE FUNCTION nope()"},
	{"event-trigger", "CREATE EVENT TRIGGER evt_existing ON ddl_command_start EXECUTE FUNCTION evt_note()"},
	{"event-trigger", "ALTER EVENT TRIGGER evt_existing DISABLE"},
	{"event-trigger", "ALTER EVENT TRIGGER evt_existing ENABLE REPLICA"},
	{"event-trigger", "ALTER EVENT TRIGGER evt_existing ENABLE ALWAYS"},
	{"event-trigger", "ALTER EVENT TRIGGER evt_existing RENAME TO evt_renamed"},
	{"event-trigger", "ALTER EVENT TRIGGER nope DISABLE"},
	{"event-trigger", "ALTER EVENT TRIGGER evt_existing OWNER TO nope_role"},
	{"event-trigger", "DROP EVENT TRIGGER evt_existing"},
	{"event-trigger", "DROP EVENT TRIGGER nope"},
	{"event-trigger", "DROP EVENT TRIGGER IF EXISTS nope"},
	{"event-trigger", "COMMENT ON EVENT TRIGGER evt_existing IS 'notes'"},

	// CALL resolves the procedure, then runs it.
	{"call", "CALL p_log('x')"},
	{"call", "CALL p_log(m => 'x')"},
	{"call", "CALL p_audit('x')"},
	{"call", "CALL p_bump(1)"},
	{"call", "CALL p_bump()"},
	{"call", "CALL nope()"},
	{"call", "CALL p_log()"},
	{"call", "CALL p_log('x', 'y')"},
	{"call", "CALL p_log(nope => 'x')"},
	{"call", "CALL f_double(1)"},
	{"call", "CALL p_broken()"},

	// DO compiles its block, then runs it.
	{"do", "DO $$ BEGIN NULL; END $$"},
	{"do", "DO $$ BEGIN RAISE NOTICE 'x'; END $$"},
	{"do", "DO $$ BEGIN NULL $$"},
	{"do", "DO $$ BEGIN x := 1; END $$"},
	{"do", "DO LANGUAGE sql $$ SELECT 1 $$"},
	{"do", "DO LANGUAGE nope $$ x $$"},
	{"do", "DO $$ BEGIN CALL p_log('x'); END $$"},
	{"do", "DO $$ DECLARE r record; BEGIN FOR r IN SELECT id FROM customers LOOP NULL; END LOOP; END $$"},
	{"do", "DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = 'customers') THEN CREATE TABLE customers (id int); END IF; END $$"},
	{"do", "DO $$ BEGIN RAISE EXCEPTION 'boom'; END $$"},
	{"do", "DO $$ BEGIN INSERT INTO missing_t VALUES (1); END $$"},
	{"do", "DO $$ BEGIN CREATE TABLE made_by_do (x int); END $$"},
}

// quotedTypeFeatures run on programs@quoted, whose enum "Mood" and table "Foo"
// have quoted mixed-case names. A quoted type name keeps its case wherever it
// names a column, argument, result or base type, and an unquoted one folds.
var quotedTypeFeatures = []struct{ kind, statement string }{
	{"quoted-type", `CREATE TABLE tq (m "Mood")`},
	{"quoted-type", `CREATE TABLE tq (m public."Mood")`},
	{"quoted-type", `CREATE TABLE tq (m "Mood"[])`},
	{"quoted-type", `CREATE TABLE tq (f "Foo")`},
	{"quoted-type", `CREATE TABLE tq (m mood)`},
	{"quoted-type", `CREATE TABLE tq (m "MOOD")`},
	{"quoted-type", `ALTER TABLE "Foo" ADD COLUMN m "Mood"`},
	{"quoted-type", `CREATE DOMAIN dq AS "Mood"`},
	{"quoted-type", `CREATE TYPE cq AS (m "Mood")`},
	{"quoted-type", `CREATE VIEW vq AS SELECT 'sad'::"Mood" AS m`},
	{"quoted-type", `CREATE FUNCTION fq(m "Mood") RETURNS text LANGUAGE sql RETURN m::text`},
	{"quoted-type", `CREATE FUNCTION fq() RETURNS "Mood" LANGUAGE sql AS $$SELECT 'sad'::"Mood"$$`},
	{"quoted-type", `CREATE FUNCTION fq(m "Mood") RETURNS int LANGUAGE plpgsql AS $$ DECLARE x "Mood"; BEGIN x := m; RETURN 1; END $$`},
}

// quotedAliasFeatures run on programs@quoted. The grammar maps an unquoted
// SQL-standard type name (integer, double precision) to its catalog name
// (int4, float8); a quoted one names a type of that exact name, which
// pg_catalog lacks.
var quotedAliasFeatures = []struct{ kind, statement string }{
	{"quoted-alias", `CREATE TABLE tq (a "integer")`},
	{"quoted-alias", `CREATE TABLE tq (a "int")`},
	{"quoted-alias", `CREATE TABLE tq (a "bigint")`},
	{"quoted-alias", `CREATE TABLE tq (a "boolean")`},
	{"quoted-alias", `CREATE TABLE tq (a "decimal"(5, 2))`},
	{"quoted-alias", `CREATE TABLE tq (a "double precision")`},
	{"quoted-alias", `CREATE TABLE tq (a "character varying"(5))`},
	{"quoted-alias", `CREATE TABLE tq (a "timestamp with time zone")`},
	{"quoted-alias", `CREATE TABLE tq (a pg_catalog."integer")`},
	{"quoted-alias", `CREATE TABLE tq (a "integer"[])`},
	{"quoted-alias", `CREATE TABLE tq (a "int4", b "varchar"(5), c "bool", d "char", e "float8")`},
	{"quoted-alias", `CREATE TABLE tq (a integer, b double precision, c character varying(5), d timestamp with time zone)`},
	{"quoted-alias", `ALTER TABLE "Foo" ADD COLUMN b "smallint"`},
	{"quoted-alias", `ALTER TABLE "Foo" ALTER COLUMN a TYPE "bigint"`},
	{"quoted-alias", `CREATE DOMAIN dq AS "integer"`},
	{"quoted-alias", `CREATE TYPE cq AS (a "real")`},
	{"quoted-alias", `CREATE TYPE rq AS RANGE (subtype = "integer")`},
	{"quoted-alias", `CREATE VIEW vq AS SELECT 1::"integer" AS a`},
	{"quoted-alias", `CREATE FUNCTION fq(a "integer") RETURNS int LANGUAGE sql RETURN a`},
	{"quoted-alias", `CREATE FUNCTION fq() RETURNS "integer" LANGUAGE sql RETURN 1`},
	{"quoted-alias", `CREATE FUNCTION fq() RETURNS int LANGUAGE plpgsql AS $$ DECLARE x "integer"; BEGIN x := 1; RETURN x; END $$`},
}

// retypeFeatures run on programs@retype: ALTER COLUMN TYPE with and without
// USING. PostgreSQL analyzes the USING expression against the table and
// casts its result to the new type by assignment; without USING it casts the
// column itself.
var retypeFeatures = []struct{ kind, statement string }{
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING t::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN v TYPE mood USING v::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING notes.t::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING 'happy'`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN n TYPE mood USING CASE WHEN n > 0 THEN 'happy' ELSE 'sad' END::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING t`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING n::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING missing::mood`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE mood USING t::nope`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN m TYPE text`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN m TYPE int USING m::text::int`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING t::int`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING length(t)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE varchar(3) USING t`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING (SELECT 1)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING count(*)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING row_number() OVER ()`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING generate_series(1, 2)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING $1`},
	{"alter-using", `ALTER TABLE notes ALTER t TYPE mood USING t::mood, ALTER v TYPE mood USING v::mood`},
	{"alter-using", `ALTER TABLE notes ALTER t TYPE mood USING t::mood, ALTER v TYPE int USING length(t)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING 'not-an-integer'`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING '12'`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN n TYPE json USING JSON_ARRAYAGG(n)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN n TYPE json USING JSON_OBJECTAGG('k' VALUE n)`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN n TYPE json USING JSON_ARRAYAGG(n) OVER ()`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE text USING tableoid::text`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE varchar(40) USING ctid::text`},
	{"alter-using", `ALTER TABLE notes ALTER COLUMN t TYPE int USING xmin::text::int`},
}

// generatedRetypeFeatures run ALTER COLUMN TYPE on programs@generated, whose
// table has a generated column and columns with defaults. A generated column
// takes no USING, and a default must cast to the new type by assignment
// whatever USING does to the rows.
var generatedRetypeFeatures = []struct{ kind, statement string }{
	{"alter-using", `ALTER TABLE gen ALTER COLUMN b TYPE int USING b::int`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN b TYPE varchar(5)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN b TYPE int`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d TYPE int USING length(d)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d TYPE varchar(5) USING upper(d)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d TYPE varchar(5)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN e TYPE text USING e::text`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN e TYPE bigint`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d DROP DEFAULT, ALTER COLUMN d TYPE int USING length(d)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d TYPE int USING length(d), ALTER COLUMN d DROP DEFAULT`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d SET DEFAULT 1, ALTER COLUMN d TYPE int USING length(d)`},
	{"alter-using", `ALTER TABLE gen ALTER COLUMN d DROP DEFAULT, ALTER COLUMN d SET DEFAULT 2`},
}
