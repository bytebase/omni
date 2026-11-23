package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/mariadb/parser"
)

// TestViewTemporalPrecisionWrapOracle holds omni's temporal opt_time_precision
// uint16 truncation validation and SHOW CREATE VIEW/TABLE deparsing to the
// pinned mariadb:11.8.8 server (replayed from tape; record with OMNI_RECORD=1).
func TestViewTemporalPrecisionWrapOracle(t *testing.T) {
	conn := mariadbDatabase(t)
	ctx := context.Background()

	if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS a"); err != nil {
		t.Fatalf("drop a: %v", err)
	}
	if _, err := conn.ExecContext(ctx, "CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT)"); err != nil {
		t.Fatalf("create a: %v", err)
	}

	viewCases := []struct {
		name string
		sql  string
	}{
		{"curtime_0", "CREATE VIEW _v AS SELECT CURRENT_TIME(0) FROM a"},
		{"curtime_6", "CREATE VIEW _v AS SELECT CURRENT_TIME(6) FROM a"},
		{"curtime_7", "CREATE VIEW _v AS SELECT CURRENT_TIME(7) FROM a"},
		{"curtime_65535", "CREATE VIEW _v AS SELECT CURRENT_TIME(65535) FROM a"},
		{"curtime_65536", "CREATE VIEW _v AS SELECT CURRENT_TIME(65536) FROM a"},
		{"curtime_maxint32", "CREATE VIEW _v AS SELECT CURRENT_TIME(2147483647) FROM a"},
		{"curtime_long_wrap0", "CREATE VIEW _v AS SELECT CURRENT_TIME(2147483648) FROM a"},
		{"curtime_long_wrap1", "CREATE VIEW _v AS SELECT CURRENT_TIME(2147483649) FROM a"},
		{"curtime_long_wrap6", "CREATE VIEW _v AS SELECT CURRENT_TIME(2147483654) FROM a"},
		{"curtime_long_wrap7", "CREATE VIEW _v AS SELECT CURRENT_TIME(2147483655) FROM a"},
		{"curtime_ulonglong_wrap0", "CREATE VIEW _v AS SELECT CURRENT_TIME(9223372036854775808) FROM a"},
		{"curtime_maxuint64", "CREATE VIEW _v AS SELECT CURRENT_TIME(18446744073709551615) FROM a"},
		{"curtime_hex_6", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x06) FROM a"},
		{"curtime_hex_65", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x41) FROM a"},
		{"curtime_hex_wrap0", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x10000) FROM a"},
		{"curtime_hex_wrap1", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x10001) FROM a"},
		{"curtime_hex_overflow", "CREATE VIEW _v AS SELECT CURRENT_TIME(0x8000000000000000) FROM a"},
		{"now_long_wrap0", "CREATE VIEW _v AS SELECT NOW(2147483648) FROM a"},
		{"now_long_wrap1", "CREATE VIEW _v AS SELECT NOW(2147483649) FROM a"},
		{"sysdate_long_wrap0", "CREATE VIEW _v AS SELECT SYSDATE(2147483648) FROM a"},
	}

	for _, tc := range viewCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := conn.ExecContext(ctx, "DROP VIEW IF EXISTS _v"); err != nil {
				t.Fatalf("drop _v: %v", err)
			}
			var engCode int
			_, engErr := conn.ExecContext(ctx, tc.sql)
			if engErr != nil {
				var me *gomysql.MySQLError
				if !errors.As(engErr, &me) {
					t.Fatalf("engine %q: non-MySQL error: %v", tc.sql, engErr)
				}
				engCode = int(me.Number)
			}

			c, err := LoadSQL("CREATE DATABASE test; USE test; CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT);")
			if err != nil {
				t.Fatalf("omni setup: %v", err)
			}
			omniCode := execCode(t, c, tc.sql)
			if omniCode != engCode {
				t.Fatalf("%s: omni code=%d, mariadb code=%d (engine err: %v)", tc.sql, omniCode, engCode, engErr)
			}
			if engCode != 0 {
				return
			}

			var viewName, engCreate, csClient, collConn string
			if err := conn.QueryRowContext(ctx, "SHOW CREATE VIEW _v").Scan(&viewName, &engCreate, &csClient, &collConn); err != nil {
				t.Fatalf("SHOW CREATE VIEW _v: %v", err)
			}
			omniCreate := c.ShowCreateView("test", "_v")
			if omniCreate != engCreate {
				t.Errorf("%s:\n  omni:    %s\n  mariadb: %s", tc.sql, omniCreate, engCreate)
			}
		})
	}

	// Column DEFAULT CURRENT_TIMESTAMP(2147483648) wraps to precision 0 on MariaDB
	// and stores as `DEFAULT current_timestamp()`.
	t.Run("default_current_timestamp_wrap0", func(t *testing.T) {
		const ddl = "CREATE TABLE t_wrap_ct (c DATETIME DEFAULT CURRENT_TIMESTAMP(2147483648))"
		if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS t_wrap_ct"); err != nil {
			t.Fatalf("drop t_wrap_ct: %v", err)
		}
		if _, err := conn.ExecContext(ctx, ddl); err != nil {
			t.Fatalf("engine %q: %v", ddl, err)
		}
		var tblName, engCreate string
		if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE t_wrap_ct").Scan(&tblName, &engCreate); err != nil {
			t.Fatalf("SHOW CREATE TABLE t_wrap_ct: %v", err)
		}

		c, err := LoadSQL("CREATE DATABASE test; USE test;")
		if err != nil {
			t.Fatalf("omni setup: %v", err)
		}
		if code := execCode(t, c, ddl); code != 0 {
			t.Fatalf("omni %q failed with code %d", ddl, code)
		}
		omniCreate := c.ShowCreateTable("test", "t_wrap_ct")
		if !strings.Contains(engCreate, "DEFAULT current_timestamp()") || !strings.Contains(omniCreate, "DEFAULT current_timestamp()") {
			t.Errorf("want DEFAULT current_timestamp() on both sides:\n  omni:    %s\n  mariadb: %s", omniCreate, engCreate)
		}
	})

	// Column TIME(65536) is rejected on both MariaDB and omni (no uint16 wrap on column types).
	t.Run("col_time_65536_rejected", func(t *testing.T) {
		const ddl = "CREATE TABLE t_bad_time (c TIME(65536))"
		if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS t_bad_time"); err != nil {
			t.Fatalf("drop t_bad_time: %v", err)
		}
		_, engErr := conn.ExecContext(ctx, ddl)
		var me *gomysql.MySQLError
		if !errors.As(engErr, &me) || me.Number != 1426 {
			t.Fatalf("engine %q: want 1426, got %v", ddl, engErr)
		}
		c, err := LoadSQL("CREATE DATABASE test; USE test;")
		if err != nil {
			t.Fatalf("omni setup: %v", err)
		}
		if code := execCode(t, c, ddl); code == 0 {
			t.Fatalf("omni %q: expected rejection, got success", ddl)
		}
	})
}

func execCode(t *testing.T, c *Catalog, stmt string) int {
	t.Helper()
	results, parseErr := c.Exec(stmt, nil)
	if parseErr != nil {
		var pe *parser.ParseError
		if errors.As(parseErr, &pe) && pe.Code != 0 {
			return pe.Code
		}
		return 1064
	}
	for _, r := range results {
		if r.Error != nil {
			var ce *Error
			if errors.As(r.Error, &ce) {
				return ce.Code
			}
			t.Fatalf("Exec(%q): unexpected non-catalog error: %v", stmt, r.Error)
		}
	}
	return 0
}
