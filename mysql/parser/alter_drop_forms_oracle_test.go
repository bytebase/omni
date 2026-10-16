package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// alterDropCodeCase is a statement and the error MySQL raises while parsing
// it: 0 when it parses, else the server error number.
type alterDropCodeCase struct {
	sql  string
	code uint16
}

var alterDropFormsOracleCases = []alterDropCodeCase{
	// opt_restrict follows DROP [COLUMN] ident and nothing else.
	{"ALTER TABLE t DROP COLUMN m RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN m CASCADE", 0},
	{"ALTER TABLE t DROP m RESTRICT", 0},
	{"ALTER TABLE t DROP k CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN m RESTRICT, DROP COLUMN k CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN m CASCADE, ADD COLUMN z INT", 0},
	{"ALTER TABLE t DROP COLUMN m /* c */ RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN `restrict`", 0},
	{"ALTER TABLE t DROP `cascade`", 0},
	{"ALTER TABLE t DROP COLUMN m RESTRICT CASCADE", 1064},
	{"ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN m RESTRICT FIRST", 1064},
	{"ALTER TABLE t DROP COLUMN m, RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN restrict", 1064},
	{"ALTER TABLE t DROP cascade", 1064},
	{"ALTER TABLE t DROP INDEX i RESTRICT", 1064},
	{"ALTER TABLE t DROP KEY i CASCADE", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY f RESTRICT", 1064},
	{"ALTER TABLE t DROP CONSTRAINT ck_n CASCADE", 1064},
	{"ALTER TABLE t DROP CHECK ck_n RESTRICT", 1064},
	{"ALTER TABLE t DROP PRIMARY KEY CASCADE", 1064},
	{"ALTER TABLE pt DROP PARTITION p0 RESTRICT", 1064},
	// KEY is required after PRIMARY and FOREIGN.
	{"ALTER TABLE t DROP PRIMARY KEY", 0},
	{"ALTER TABLE t DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY", 1064},
	{"ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)", 1064},
	{"ALTER TABLE t DROP PRIMARY INDEX", 1064},
	{"ALTER TABLE t DROP FOREIGN f", 1064},
	{"ALTER TABLE t DROP FOREIGN INDEX f", 1064},
	{"ALTER TABLE t DROP FOREIGN IF EXISTS f", 1064},
	// IF EXISTS on a table element is MariaDB syntax.
	{"ALTER TABLE t DROP COLUMN IF EXISTS m", 1064},
	{"ALTER TABLE t DROP IF EXISTS m", 1064},
	{"ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT", 1064},
	{"ALTER TABLE t DROP INDEX IF EXISTS i", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f", 1064},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n", 1064},
	{"ALTER TABLE t DROP CHECK IF EXISTS ck_n", 1064},
	{"ALTER TABLE pt DROP PARTITION IF EXISTS p0", 1064},
}

// TestAlterTableDropFormsOracle asserts omni agrees with a live mysql:8.0 on
// the ALTER TABLE ... DROP forms in alterDropFormsOracleCases (container
// test). MySQL's verdict comes from PREPARE, which parses without running the
// statement. The expected code pins the engine, so an image upgrade that
// changes a verdict fails here instead of passing silently.
func TestAlterTableDropFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, ddl := range []string{
		"CREATE TABLE p (id INT PRIMARY KEY)",
		"CREATE TABLE t (id INT PRIMARY KEY, n INT, m INT, k INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))",
		"CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	for _, tc := range alterDropFormsOracleCases {
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
			t.Errorf("MySQL code %d, want %d for %q: %v", code, tc.code, tc.sql, err)
			continue
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, MySQL code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
