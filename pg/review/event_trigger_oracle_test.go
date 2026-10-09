package review

import (
	"context"
	"slices"
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
		`CREATE TABLE ev_log (event text, tag text)`,
		`CREATE FUNCTION ev_note() RETURNS event_trigger LANGUAGE plpgsql AS $$ BEGIN INSERT INTO ev_log VALUES (lower(tg_event), tg_tag); END $$`,
		`CREATE FOREIGN DATA WRAPPER ev_fdw`,
		`CREATE SERVER ev_srv FOREIGN DATA WRAPPER ev_fdw`,
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
		"ALTER TABLE ev_t ADD COLUMN e int DEFAULT 5",
		"CREATE DOMAIN ev_pos AS int CHECK (VALUE > 0)",
		"ALTER TABLE ev_t ADD COLUMN f ev_pos",
		"ALTER TABLE ev_t ADD COLUMN g int GENERATED ALWAYS AS (x * 2) STORED",
		"ALTER TABLE ev_t ADD COLUMN h int GENERATED ALWAYS AS IDENTITY",
		"ALTER TABLE ev_t ADD COLUMN i text NOT NULL, ADD COLUMN j varchar(10) COLLATE \"C\", ADD COLUMN k int[], ADD COLUMN l pg_catalog.int8",
		"ALTER TABLE ev_t SET UNLOGGED",
		"ALTER TABLE ev_t SET LOGGED",
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
		"CREATE PROCEDURE ev_pr() LANGUAGE sql AS 'SELECT 1'",
		"DROP PROCEDURE ev_pr()",
		"CREATE MATERIALIZED VIEW ev_mv AS SELECT 1 AS a",
		"ALTER VIEW ev_v RENAME TO ev_v2",
		"ALTER VIEW ev_v2 RENAME TO ev_v",
		"ALTER INDEX ev_p_code RENAME TO ev_p_code2",
		"CREATE FOREIGN TABLE ev_ft (a int) SERVER ev_srv",
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
		rows, err := conn.QueryContext(ctx, `DELETE FROM ev_log RETURNING event, tag`)
		if err != nil {
			t.Fatal(err)
		}
		raised := make(map[string]bool)
		var tags []string
		for rows.Next() {
			var event, tag string
			if err := rows.Scan(&event, &tag); err != nil {
				t.Fatal(err)
			}
			raised[event] = true
			if event == "ddl_command_end" {
				tags = append(tags, tag)
			}
		}
		rows.Close()
		stmts, finding := parse(sql, statementRanges(sql))
		if finding != nil || len(stmts) != 1 {
			t.Fatalf("%s: parse: %v", sql, finding)
		}
		n := stmts[0].node
		// A tag the scan names is the server's.
		if tag := commandTag(n); tag != "" && len(tags) > 0 && !slices.Contains(tags, tag) {
			t.Errorf("%s: the scan names tag %q, the server %q", sql, tag, tags)
		}
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

// TestEventTriggerClassification pins the statements the scan expects
// not to raise sql_drop or table_rewrite, which
// TestEventTriggerEventsAgainstPostgres runs on the server.
func TestEventTriggerClassification(t *testing.T) {
	for _, c := range []struct {
		sql             string
		drops, rewrites bool
	}{
		{"ALTER TABLE ev_t ADD COLUMN a int", false, false},
		{`ALTER TABLE ev_t ADD COLUMN i text NOT NULL, ADD COLUMN j varchar(10) COLLATE "C", ADD COLUMN k int[], ADD COLUMN l pg_catalog.int8`, false, false},
		{"ALTER TABLE ev_t ADD CONSTRAINT ev_fk FOREIGN KEY (p_id) REFERENCES ev_p", false, false},
		{"ALTER TABLE ev_t ALTER COLUMN x SET DEFAULT 3, ADD COLUMN d int", false, false},
		{"ALTER TABLE ev_t ALTER COLUMN x DROP DEFAULT", true, false},
		{"ALTER TABLE ev_s.ev_t2 DROP COLUMN aa", true, false},
		{"ALTER TABLE ev_t ADD COLUMN e int DEFAULT 5", false, true},
		{"ALTER TABLE ev_t ADD COLUMN f ev_pos", false, true},
		{"ALTER TABLE ev_t ALTER COLUMN a TYPE bigint", true, true},
		{"ALTER TABLE ev_t SET LOGGED", true, true},
		{"ALTER TABLE ev_t RENAME TO ev_t2", false, false},
	} {
		stmts, finding := parse(c.sql, statementRanges(c.sql))
		if finding != nil || len(stmts) != 1 {
			t.Fatalf("%s: parse: %v", c.sql, finding)
		}
		if got := mayDrop(stmts[0].node); got != c.drops {
			t.Errorf("mayDrop(%s) = %v, want %v", c.sql, got, c.drops)
		}
		if got := mayRewrite(stmts[0].node); got != c.rewrites {
			t.Errorf("mayRewrite(%s) = %v, want %v", c.sql, got, c.rewrites)
		}
	}
}
