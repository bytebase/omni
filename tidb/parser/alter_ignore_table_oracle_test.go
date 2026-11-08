package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// alterIgnoreCase is a statement and the error TiDB raises while parsing it:
// 0 when it parses, else the server error number.
type alterIgnoreCase struct {
	sql  string
	code uint16
}

// IgnoreOptional is "IGNORE or empty" and appears only in the AlterTableStmt
// family, so the IGNORE keyword is legal between ALTER and TABLE and nowhere
// else. Its value is discarded (TiDB does not model it either).
//
// Ref: TiDB v8.5.5 pkg/parser/parser.y IgnoreOptional and
// AlterTableStmt: "ALTER" IgnoreOptional "TABLE" ...
var alterIgnoreTiDBOracleCases = []alterIgnoreCase{
	// Legacy copy-ALTER syntax TiDB still accepts.
	{"ALTER IGNORE TABLE t", 0},
	{"ALTER IGNORE TABLE t ADD COLUMN z INT", 0},
	{"ALTER IGNORE TABLE t RENAME TO t2", 0},
	{"ALTER IGNORE TABLE t ANALYZE PARTITION p0", 0},
	// IGNORE is valid only immediately before TABLE.
	{"ALTER IGNORE DATABASE d", 1064},
	{"ALTER IGNORE SEQUENCE s", 1064},
	{"ALTER TABLE IGNORE t", 1064},
	{"ALTER IGNORE", 1064},
}

// TestAlterIgnoreTableTiDBOracle asserts omni agrees with a live TiDB v8.5.5
// on the optional IGNORE keyword between ALTER and TABLE (container test).
// TiDB's verdict comes from PREPARE, which parses without running the
// statement. The expected code pins the engine, so an image upgrade that
// changes a verdict fails here instead of passing silently.
func TestAlterIgnoreTableTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()
	for _, ddl := range []string{
		"DROP DATABASE IF EXISTS omni_alter_ignore",
		"CREATE DATABASE omni_alter_ignore",
		"USE omni_alter_ignore",
		"CREATE TABLE t (id INT PRIMARY KEY) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN MAXVALUE)",
	} {
		if _, err := conn.ExecContext(tc.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_alter_ignore")

	for _, c := range alterIgnoreTiDBOracleCases {
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
