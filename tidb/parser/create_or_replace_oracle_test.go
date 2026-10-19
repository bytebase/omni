package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestCreateOrReplaceTiDBOracle pins TiDB v8.5.5's verdict on
// createOrReplaceCases and checks omni against it (container test). TiDB's
// verdict comes from PREPARE, which parses without running the statement. The
// expected code pins the engine, so an image upgrade that changes a verdict
// fails here instead of passing silently. For a 1064, TiDB's message must
// quote the same rest of the statement as the case's near, which places
// omni's error at the same token as TiDB's.
func TestCreateOrReplaceTiDBOracle(t *testing.T) {
	tc := startTiDB(t)
	// USE and PREPARE are session state: keep every statement on one
	// connection.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("acquire connection: %v", err)
	}
	defer conn.Close()

	conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_or_replace")
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_or_replace")
	for _, stmt := range []string{
		"CREATE DATABASE omni_or_replace",
		"USE omni_or_replace",
		"CREATE TABLE t (id INT PRIMARY KEY, n INT)",
	} {
		if _, err := conn.ExecContext(tc.ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}

	for _, c := range createOrReplaceCases {
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
		if code == 1064 {
			if quote := `near "` + c.near + `"`; !strings.Contains(err.Error(), quote) {
				t.Errorf("TiDB error for %q does not quote %q: %v", c.sql, quote, err)
			}
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != (code == 0) {
			t.Errorf("DIVERGENCE omni accepts=%v, TiDB code %d for %q: %v", omniAccepts, code, c.sql, perr)
		}
	}
}
