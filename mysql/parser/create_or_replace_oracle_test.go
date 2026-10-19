package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestCreateOrReplaceOracle pins MySQL 8.0's verdict on createOrReplaceCases
// and checks omni against it (container test). MySQL's verdict comes from
// PREPARE, which parses without running the statement. The expected code pins
// the engine, so an image upgrade that changes a verdict fails here instead of
// passing silently. For a 1064, MySQL's message must quote the same rest of
// the statement as the case's near, which places omni's error at the same
// token as MySQL's.
func TestCreateOrReplaceOracle(t *testing.T) {
	o := startParserOracle(t)
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE t (id INT PRIMARY KEY, n INT)"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	for _, tc := range createOrReplaceCases {
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
		if code == 1064 {
			// ER_PARSE_ERROR quotes at most 80 bytes of the rest ("near '%-.80s'").
			near := tc.near
			quote := "near '" + near + "'"
			if len(near) > 80 {
				quote = "near '" + near[:80]
			}
			if !strings.Contains(err.Error(), quote) {
				t.Errorf("MySQL error for %q does not quote %q: %v", tc.sql, quote, err)
			}
		}
		_, perr := Parse(tc.sql)
		if omniAccepts := perr == nil; omniAccepts != parsesInMySQL(code) {
			t.Errorf("DIVERGENCE omni accepts=%v, MySQL code %d for %q: %v", omniAccepts, code, tc.sql, perr)
		}
	}
}
