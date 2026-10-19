package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// parseCodeCase is a statement and the error MariaDB raises while parsing it:
// 0 when it parses, else the server error number.
type parseCodeCase struct {
	sql  string
	code uint16
}

var dropClauseOracleCases = []parseCodeCase{
	// IF EXISTS on every named table element (opt_if_exists_table_element).
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS no_such", 0},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n", 0},
	{"ALTER TABLE t DROP INDEX IF EXISTS no_such", 0},
	{"ALTER TABLE t DROP KEY IF EXISTS i", 0},
	{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f", 0},
	{"ALTER TABLE t DROP COLUMN IF EXISTS no_such", 0},
	{"ALTER TABLE t DROP IF EXISTS no_such", 0},
	{"ALTER TABLE t DROP INDEX IF EXISTS `if`", 0},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS c1, DROP INDEX IF EXISTS i1, DROP FOREIGN KEY IF EXISTS f1", 0},
	{"ALTER TABLE pt DROP PARTITION IF EXISTS no_such", 0},
	{"ALTER TABLE pt DROP PARTITION IF EXISTS p0, p1", 0},
	{"ALTER TABLE t DROP CONSTRAINT ck_n", 0},
	{"ALTER TABLE t DROP CONSTRAINT PRIMARY KEY", 0},
	// MySQL 8.0's DROP CHECK does not exist in MariaDB.
	{"ALTER TABLE t DROP CHECK ck_n", 1064},
	{"ALTER TABLE t DROP CHECK IF EXISTS ck_n", 1064},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS", 1064},
	{"ALTER TABLE t DROP INDEX IF no_such", 1064},
	{"ALTER TABLE t DROP INDEX EXISTS no_such", 1064},
	{"ALTER TABLE t DROP CONSTRAINT PRIMARY", 1064},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS PRIMARY KEY", 1064},
	{"ALTER TABLE t DROP PRIMARY KEY IF EXISTS", 1064},
	{"ALTER TABLE pt DROP PARTITION IF EXISTS", 1064},
}

var orReplaceTableOracleCases = []parseCodeCase{
	{"CREATE OR REPLACE TABLE r (id INT)", 0},
	{"CREATE OR REPLACE TEMPORARY TABLE r (id INT)", 0},
	{"CREATE OR REPLACE TABLE r LIKE t", 0},
	{"CREATE OR REPLACE TABLE r AS SELECT 1 AS a", 0},
	{"CREATE OR REPLACE TABLE r (id INT) SELECT 1 AS id", 0},
	{"CREATE OR REPLACE TABLE IF NOT EXISTS r (id INT)", 1221},
	{"CREATE OR REPLACE TEMPORARY TABLE IF NOT EXISTS r (id INT)", 1221},
	{"CREATE OR REPLACE TABLE IF NOT EXISTS r LIKE t", 1221},
	{"CREATE OR TABLE r (id INT)", 1064},
	{"CREATE OR REPLACE TABLE IF EXISTS r (id INT)", 1064},
}

// TestMariaDBDropClauseOracle asserts omni agrees with a live mariadb:11.8.8 on
// the ALTER TABLE ... DROP forms in dropClauseOracleCases (container test).
func TestMariaDBDropClauseOracle(t *testing.T) {
	assertParseCodes(t, dropClauseOracleCases)
}

// TestMariaDBOrReplaceTableOracle asserts omni agrees with a live
// mariadb:11.8.8 on the CREATE OR REPLACE TABLE forms in
// orReplaceTableOracleCases (container test).
func TestMariaDBOrReplaceTableOracle(t *testing.T) {
	assertParseCodes(t, orReplaceTableOracleCases)
}

// assertParseCodes checks each case against MariaDB and omni. MariaDB's verdict
// comes from PREPARE, which parses without running the statement. Unlike
// classify, it does not count only 1064 as a rejection: ER_WRONG_USAGE (1221)
// for OR REPLACE with IF NOT EXISTS is raised by a grammar action, so the
// statement never parses. The expected code pins the engine, so an image
// upgrade that changes a verdict fails here instead of passing silently.
func assertParseCodes(t *testing.T, cases []parseCodeCase) {
	t.Helper()
	o := startMariaDB(t)
	for _, ddl := range []string{
		"CREATE TABLE p (id INT PRIMARY KEY)",
		"CREATE TABLE t (id INT PRIMARY KEY, n INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))",
		"CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20), PARTITION p2 VALUES LESS THAN MAXVALUE)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	for _, tc := range cases {
		var code uint16
		_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(tc.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", tc.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != tc.code {
			t.Errorf("MariaDB code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
