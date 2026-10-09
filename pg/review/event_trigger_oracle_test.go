package review

import (
	"context"
	"testing"
)

// TestEventTriggerEventsAgainstPostgres runs statements on PostgreSQL
// under event triggers that log every event they see, and checks that
// the scan expects each one: a statement that raised ddl_command_end is
// DDL, one that raised sql_drop may drop, and one that raised
// table_rewrite may rewrite. Claiming an event that does not come only
// costs findings; missing one would let a trigger run unseen.
func TestEventTriggerEventsAgainstPostgres(t *testing.T) {
	db := startPostgres(t)
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, setup := range []string{
		`CREATE TABLE ev_log (event text)`,
		`CREATE FUNCTION ev_note() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN INSERT INTO ev_log VALUES (lower(tg_event)); END $$`,
		`CREATE EVENT TRIGGER ev_end ON ddl_command_end EXECUTE FUNCTION ev_note()`,
		`CREATE EVENT TRIGGER ev_drop ON sql_drop EXECUTE FUNCTION ev_note()`,
		`CREATE EVENT TRIGGER ev_rewrite ON table_rewrite EXECUTE FUNCTION ev_note()`,
	} {
		if _, err := conn.ExecContext(ctx, setup); err != nil {
			t.Fatalf("%s: %v", setup, err)
		}
	}
	// Each statement runs after the ones before it, and must succeed.
	for _, sql := range []string{
		"CREATE TABLE ev_p (id int PRIMARY KEY, code text, n int DEFAULT 1)",
		"CREATE UNIQUE INDEX ev_p_code ON ev_p (code)",
		"CREATE TABLE ev_t (id int PRIMARY KEY, p_id int, x int)",
		"CREATE VIEW ev_v AS SELECT 1 AS a",
		"ALTER TABLE ev_t ADD COLUMN a int",
		"ALTER TABLE ev_t ADD COLUMN b int DEFAULT random()::int",
		"ALTER TABLE ev_t ADD COLUMN c serial",
		"ALTER TABLE ev_t ADD CONSTRAINT ev_fk FOREIGN KEY (p_id) REFERENCES ev_p",
		"ALTER TABLE ev_t ADD CONSTRAINT ev_chk CHECK (x > 0)",
		"CREATE UNIQUE INDEX ev_x ON ev_t (x)",
		"ALTER TABLE ev_t ADD CONSTRAINT ev_x_key UNIQUE USING INDEX ev_x",
		"ALTER TABLE ev_t ALTER COLUMN x SET NOT NULL",
		"ALTER TABLE ev_t ALTER COLUMN x DROP NOT NULL",
		"ALTER TABLE ev_t ALTER COLUMN x SET DEFAULT 2",
		"ALTER TABLE ev_t ALTER COLUMN x SET DEFAULT 3, ADD COLUMN d int",
		"ALTER TABLE ev_t ALTER COLUMN x DROP DEFAULT",
		"ALTER TABLE ev_t ALTER COLUMN a TYPE bigint",
		"ALTER TABLE ev_t ALTER COLUMN x SET STATISTICS 100",
		"ALTER TABLE ev_t OWNER TO CURRENT_USER",
		"ALTER TABLE ev_t VALIDATE CONSTRAINT ev_chk",
		"ALTER TABLE ev_t ENABLE ROW LEVEL SECURITY",
		"ALTER TABLE ev_t DISABLE ROW LEVEL SECURITY",
		"ALTER TABLE ev_t RENAME COLUMN a TO aa",
		"ALTER TABLE ev_t RENAME TO ev_t2",
		"CREATE SCHEMA ev_s",
		"ALTER TABLE ev_t2 SET SCHEMA ev_s",
		"ALTER SEQUENCE ev_s.ev_t_c_seq OWNED BY NONE",
		"CREATE OR REPLACE VIEW ev_v AS SELECT 1 AS a, 2 AS b",
		"COMMENT ON TABLE ev_p IS 'x'",
		"GRANT SELECT ON ev_p TO PUBLIC",
		"REVOKE SELECT ON ev_p FROM PUBLIC",
		"CREATE TYPE ev_mood AS ENUM ('a')",
		"ALTER TYPE ev_mood ADD VALUE 'b'",
		"CREATE DOMAIN ev_d AS int",
		"CREATE SEQUENCE ev_q",
		"CREATE FUNCTION ev_f() RETURNS int LANGUAGE sql AS 'SELECT 1'",
		"CREATE OR REPLACE FUNCTION ev_f() RETURNS int LANGUAGE sql AS 'SELECT 2'",
		"CREATE TABLE ev_c AS SELECT 1 AS a",
		"SELECT 1 AS a INTO ev_c2",
		"CREATE INDEX ON ev_p (n)",
		"CREATE SCHEMA ev_s2 CREATE TABLE x (id int UNIQUE) CREATE VIEW y AS SELECT 1 AS a",
		"INSERT INTO ev_p VALUES (1, 'a', 1)",
		"UPDATE ev_p SET n = 2",
		"SELECT * FROM ev_p",
		"EXPLAIN SELECT * FROM ev_p",
		"SET search_path = public",
		"ALTER TABLE ev_s.ev_t2 DROP CONSTRAINT ev_chk",
		"ALTER TABLE ev_s.ev_t2 DROP COLUMN aa",
		"DROP DOMAIN ev_d",
		"DROP VIEW ev_v",
		"DROP TABLE ev_c, ev_c2",
		"DELETE FROM ev_p",
	} {
		if _, err := conn.ExecContext(ctx, sql); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		rows, err := conn.QueryContext(ctx, `DELETE FROM ev_log RETURNING event`)
		if err != nil {
			t.Fatal(err)
		}
		raised := make(map[string]bool)
		for rows.Next() {
			var event string
			if err := rows.Scan(&event); err != nil {
				t.Fatal(err)
			}
			raised[event] = true
		}
		rows.Close()
		stmts, finding := parse(sql, statementRanges(sql))
		if finding != nil || len(stmts) != 1 {
			t.Fatalf("%s: parse: %v", sql, finding)
		}
		n := stmts[0].node
		if raised["ddl_command_end"] && !isDDL(n) {
			t.Errorf("%s raised ddl_command_end, which the scan does not expect", sql)
		}
		if raised["sql_drop"] && !mayDrop(n) {
			t.Errorf("%s raised sql_drop, which the scan does not expect", sql)
		}
		if raised["table_rewrite"] && !mayRewrite(n) {
			t.Errorf("%s raised table_rewrite, which the scan does not expect", sql)
		}
	}
}
