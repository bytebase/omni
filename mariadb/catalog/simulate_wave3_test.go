package catalog

import (
	"strings"
	"testing"
)

func TestWave3CatalogFixes(t *testing.T) {
	t.Run("definer formatting on routine trigger event", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE DATABASE d;
USE d;
CREATE TABLE t (b INT PRIMARY KEY, c INT);
CREATE DEFINER=` + "`root`@`%`" + ` PROCEDURE p() COMMIT;
CREATE DEFINER=` + "`root`@`%`" + ` FUNCTION f() RETURNS INT RETURN 1;
CREATE DEFINER=` + "`root`@`%`" + ` TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.c = 1;
CREATE DEFINER=` + "`root`@`%`" + ` EVENT ev ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' DO BEGIN END;
`)
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		for kind, got := range map[string]string{
			"procedure": c.ShowCreateProcedure("d", "p"),
			"function":  c.ShowCreateFunction("d", "f"),
			"trigger":   c.ShowCreateTrigger("d", "tr"),
			"event":     c.ShowCreateEvent("d", "ev"),
		} {
			if !strings.Contains(got, "DEFINER=`root`@`%`") {
				t.Errorf("%s ShowCreate did not backtick-quote DEFINER: %s", kind, got)
			}
		}
	})

	t.Run("routine default characteristics omitted", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE DATABASE d;
USE d;
CREATE PROCEDURE p() COMMIT;
ALTER PROCEDURE p CONTAINS SQL SQL SECURITY DEFINER;
`)
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		got := c.ShowCreateProcedure("d", "p")
		if strings.Contains(got, "CONTAINS SQL") || strings.Contains(got, "SQL SECURITY DEFINER") {
			t.Errorf("default characteristics should be omitted from SHOW CREATE PROCEDURE:\n%s", got)
		}
		if _, err := c.Exec("ALTER PROCEDURE p READS SQL DATA SQL SECURITY INVOKER COMMENT 'hi';", nil); err != nil {
			t.Fatalf("alter: %v", err)
		}
		got = c.ShowCreateProcedure("d", "p")
		for _, want := range []string{"READS SQL DATA", "SQL SECURITY INVOKER", "COMMENT 'hi'"} {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q in SHOW CREATE PROCEDURE:\n%s", want, got)
			}
		}
	})

	t.Run("routine drop and alter unknown schema answer 1305", func(t *testing.T) {
		c, err := LoadSQL("CREATE DATABASE d; USE d;")
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		for _, stmt := range []string{
			"DROP FUNCTION nosuchdb.f",
			"DROP PROCEDURE nosuchdb.p",
			"ALTER FUNCTION nosuchdb.f",
			"ALTER PROCEDURE nosuchdb.p",
		} {
			res, err := c.Exec(stmt, nil)
			if err != nil {
				t.Fatalf("%s: parse err %v", stmt, err)
			}
			if got := errorCode(res[0].Error); got != 1305 {
				t.Errorf("%s: got code %d (%v), want 1305", stmt, got, res[0].Error)
			}
		}
	})

	t.Run("sequence view errors, AS narrowing overflow, increment range", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE DATABASE d;
USE d;
CREATE TABLE t (b INT PRIMARY KEY);
CREATE VIEW v AS SELECT b FROM t;
CREATE SEQUENCE s;
`)
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		cases := []struct {
			stmt string
			want int
		}{
			{"ALTER SEQUENCE v NOMINVALUE", 4089},
			{"DROP SEQUENCE v", 1965},
			{"ALTER SEQUENCE s AS INT", 1264},
			{"ALTER SEQUENCE s AS TINYINT", 1264},
			{"ALTER SEQUENCE s AS INT UNSIGNED", 1264},
			{"ALTER SEQUENCE s AS BIGINT", 0},
			{"ALTER SEQUENCE s INCREMENT - 9223372036854775808", 4085},
			{"ALTER SEQUENCE s INCREMENT - 2147483648", 0},
		}
		for _, tc := range cases {
			res, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatalf("%s: parse err %v", tc.stmt, err)
			}
			if got := errorCode(res[0].Error); got != tc.want {
				t.Errorf("%s: got code %d (%v), want %d", tc.stmt, got, res[0].Error, tc.want)
			}
		}
	})

	t.Run("event schedule column ref 1054 and rename same name 1551", func(t *testing.T) {
		c, err := LoadSQL(`
CREATE DATABASE d;
USE d;
CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2026-11-13 00:00:00' DO BEGIN END;
`)
		if err != nil {
			t.Fatalf("LoadSQL: %v", err)
		}
		cases := []struct {
			stmt string
			want int
		}{
			{"CREATE EVENT e ON SCHEDULE AT a DO COMMIT", 1054},
			{"CREATE EVENT missing ON SCHEDULE AT a DO COMMIT", 1054},
			{"ALTER EVENT e ON SCHEDULE EVERY a DAY", 1054},
			{"ALTER EVENT missing ON SCHEDULE AT a", 1054},
			{"ALTER EVENT e ON SCHEDULE EVERY 1 DAY STARTS a", 1054},
			{"ALTER EVENT e ON SCHEDULE EVERY 1 DAY ENDS a", 1054},
			{"ALTER EVENT e RENAME TO e", 1551},
			{"ALTER EVENT missing RENAME TO missing", 1551},
			{"ALTER EVENT e RENAME TO e2", 0},
		}
		for _, tc := range cases {
			res, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatalf("%s: parse err %v", tc.stmt, err)
			}
			if got := errorCode(res[0].Error); got != tc.want {
				t.Errorf("%s: got code %d (%v), want %d", tc.stmt, got, res[0].Error, tc.want)
			}
		}
		if c.GetDatabase("d").Events["e2"] == nil || c.GetDatabase("d").Events["e"] != nil {
			t.Errorf("expected event renamed from e to e2")
		}
	})

	t.Run("table DDL targeting a sequence", func(t *testing.T) {
		newCat := func(t *testing.T) *Catalog {
			t.Helper()
			c, err := LoadSQL("CREATE DATABASE d; USE d; CREATE SEQUENCE a;")
			if err != nil {
				t.Fatalf("LoadSQL: %v", err)
			}
			return c
		}
		cases := []struct {
			stmt string
			want int
		}{
			{"TRUNCATE a", 1031},
			{"CREATE INDEX idx ON a(a)", 4086},
			{"DROP INDEX idx ON a", 1091},
			{"DROP INDEX IF EXISTS idx ON a", 0},
			{"ALTER TABLE a ADD x INT", 4086},
			{"ALTER TABLE a ADD x INT AFTER a", 1054},
			{"ALTER TABLE a ADD KEY(a)", 4086},
			{"ALTER TABLE a ADD CHECK(a)", 4086},
			{"ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a", 1072},
			{"ALTER TABLE a DROP COLUMN a", 1091},
			{"ALTER TABLE a DROP IF EXISTS a", 0},
			{"ALTER TABLE a DROP KEY a", 1091},
			{"ALTER TABLE a MODIFY a INT", 1054},
			{"ALTER TABLE a ALTER KEY a IGNORED", 1176},
			{"ALTER TABLE a OPTIMIZE PARTITION ALL", 1505},
			{"ALTER TABLE a REMOVE PARTITIONING", 1505},
			{"ALTER TABLE a PARTITION BY SYSTEM_TIME", 4124},
			{"ALTER TABLE a WITH SYSTEM VERSIONING", 4086},
			{"ALTER TABLE a DROP SYSTEM VERSIONING", 4124},
			{"ALTER TABLE a ADD x INT AS ROW START", 4124},
			{"ALTER TABLE a ORDER BY a", 4086},
			{"ALTER TABLE a a = 's'", 1911},
			{"ALTER TABLE a IMPORT TABLESPACE", 1813},
			{"ALTER TABLE a DISCARD TABLESPACE", 0},
			{"ALTER TABLE a COMMENT 'hello'", 0},
			{"ALTER TABLE a ENGINE = Aria", 0},
			{"ALTER TABLE a ROW_FORMAT FIXED", 0},
			{"ALTER TABLE a RENAME TO b", 0},
		}
		for _, tc := range cases {
			c := newCat(t)
			res, err := c.Exec(tc.stmt, nil)
			if err != nil {
				t.Fatalf("%s: parse err %v", tc.stmt, err)
			}
			if got := errorCode(res[0].Error); got != tc.want {
				t.Errorf("%s: got code %d (%v), want %d", tc.stmt, got, res[0].Error, tc.want)
			}
		}

		// Verify RENAME, COMMENT, ENGINE, and DISCARD TABLESPACE mutate the sequence.
		c := newCat(t)
		if _, err := c.Exec("ALTER TABLE a COMMENT 'hi', ENGINE = Aria, RENAME TO b;", nil); err != nil {
			t.Fatalf("alter sequence table: %v", err)
		}
		if c.GetDatabase("d").Sequences["a"] != nil || c.GetDatabase("d").Sequences["b"] == nil {
			t.Fatalf("expected sequence a renamed to b")
		}
		got := c.ShowCreateSequence("d", "b")
		if !strings.Contains(got, "ENGINE=Aria") || !strings.Contains(got, "COMMENT='hi'") {
			t.Errorf("unexpected ShowCreateSequence after ALTER TABLE: %s", got)
		}
	})
}
