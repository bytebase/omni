package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// alterDropCodeCase is a statement and the error TiDB raises while parsing it:
// 0 when it parses, else the server error number.
type alterDropCodeCase struct {
	sql  string
	code uint16
}

var alterDropFormsTiDBOracleCases = []alterDropCodeCase{
	// RestrictOrCascadeOpt follows DROP [COLUMN] [IF EXISTS] ColumnName and
	// nothing else.
	{"ALTER TABLE t DROP COLUMN m RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN m CASCADE", 0},
	{"ALTER TABLE t DROP m RESTRICT", 0},
	{"ALTER TABLE t DROP k CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN m RESTRICT, DROP COLUMN k CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN m CASCADE, ADD COLUMN z INT", 0},
	{"ALTER TABLE t DROP COLUMN m /* c */ RESTRICT", 0},
	{"ALTER TABLE t DROP COLUMN `restrict`", 0},
	{"ALTER TABLE t DROP `cascade`", 0},
	{"ALTER TABLE t DROP COLUMN IF EXISTS m RESTRICT", 0},
	{"ALTER TABLE t DROP IF EXISTS m CASCADE", 0},
	{"ALTER TABLE t DROP COLUMN IF EXISTS m", 0},
	{"ALTER TABLE t DROP IF EXISTS no_such", 0},
	{"ALTER TABLE t DROP COLUMN m RESTRICT CASCADE", 1064},
	{"ALTER TABLE t DROP COLUMN m RESTRICT RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN m RESTRICT FIRST", 1064},
	{"ALTER TABLE t DROP COLUMN m, RESTRICT", 1064},
	{"ALTER TABLE t DROP COLUMN restrict", 1064},
	{"ALTER TABLE t DROP cascade", 1064},
	{"ALTER TABLE t DROP INDEX i RESTRICT", 1064},
	{"ALTER TABLE t DROP KEY i CASCADE", 1064},
	{"ALTER TABLE t DROP KEY IF EXISTS i CASCADE", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY f RESTRICT", 1064},
	{"ALTER TABLE t DROP CONSTRAINT ck_n CASCADE", 1064},
	{"ALTER TABLE t DROP CHECK ck_n RESTRICT", 1064},
	{"ALTER TABLE t DROP PRIMARY KEY CASCADE", 1064},
	{"ALTER TABLE pt DROP PARTITION p0 RESTRICT", 1064},
	// KEY is required after PRIMARY and FOREIGN, and DROP FOREIGN KEY takes
	// no IF EXISTS.
	{"ALTER TABLE t DROP PRIMARY KEY", 0},
	{"ALTER TABLE t DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY KEY, DROP FOREIGN KEY f", 0},
	{"ALTER TABLE t DROP PRIMARY", 1064},
	{"ALTER TABLE t DROP PRIMARY, ADD PRIMARY KEY (id)", 1064},
	{"ALTER TABLE t DROP PRIMARY INDEX", 1064},
	{"ALTER TABLE t DROP FOREIGN f", 1064},
	{"ALTER TABLE t DROP FOREIGN INDEX f", 1064},
	{"ALTER TABLE t DROP FOREIGN IF EXISTS f", 1064},
	{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f", 1064},
	{"ALTER TABLE t DROP CONSTRAINT IF EXISTS ck_n", 1064},
	{"ALTER TABLE t DROP CHECK IF EXISTS ck_n", 1064},
}

// TestAlterTableDropFormsTiDBOracle asserts omni agrees with a live TiDB
// v8.5.5 on the ALTER TABLE ... DROP forms in alterDropFormsTiDBOracleCases
// (container test). TiDB's verdict comes from PREPARE, which parses without
// running the statement. The expected code pins the engine, so an image
// upgrade that changes a verdict fails here instead of passing silently.
func TestAlterTableDropFormsTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	// USE and PREPARE are session state: a dedicated connection keeps every
	// statement on the database it selects.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()
	for _, ddl := range []string{
		"DROP DATABASE IF EXISTS omni_alter_drop",
		"CREATE DATABASE omni_alter_drop",
		"USE omni_alter_drop",
		"CREATE TABLE p (id INT PRIMARY KEY)",
		"CREATE TABLE t (id INT PRIMARY KEY, n INT, m INT, k INT, pid INT, CONSTRAINT ck_n CHECK (n > 0), INDEX i (n), CONSTRAINT f FOREIGN KEY (pid) REFERENCES p (id))",
		"CREATE TABLE pt (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)",
	} {
		if _, err := conn.ExecContext(tc.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_alter_drop")

	for _, c := range alterDropFormsTiDBOracleCases {
		var code uint16
		_, err := conn.ExecContext(tc.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(c.sql, "'", "''")+"'")
		if err != nil {
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Errorf("PREPARE %q: %v", c.sql, err)
				continue
			}
			code = myErr.Number
		}
		if code != c.code {
			t.Errorf("TiDB code %d, want %d for %q: %v", code, c.code, c.sql, err)
			continue
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, TiDB code %d for %q: %v", omniAccepts, code, c.sql, perr)
		}
	}
}
