package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestTiDBTemporalFormsMatchEngine holds every row of temporalFormCases to a live
// TiDB v8.5.5 (container test), so the recorded verdicts are the engine's and not
// MySQL's inherited. They differ here: NOW takes an expression list, a function
// word with a space before its '(' is mostly a syntax error rather than an ordinary
// call, and a literal up to 18446744073709551615 parses as a precision.
//
// PREPARE reads these through tidbTemporalVerdict rather than prepareVerdict, because
// TiDB answers 1064 for two different things here.
func TestTiDBTemporalFormsMatchEngine(t *testing.T) {
	tc := startTiDB(t)
	// USE and the user variable are session state, so every statement takes one
	// connection.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS omni_temporal",
		"USE omni_temporal",
		"CREATE TABLE IF NOT EXISTS t (a DATETIME, b INT, c INT, `now` INT, `curdate` INT, `curtime` INT, " +
			"`date_add` INT, `date_sub` INT, `sysdate` INT, `adddate` INT, `subdate` INT)",
	} {
		if _, err := conn.ExecContext(tc.ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	defer conn.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_temporal")

	for _, c := range temporalFormCases {
		if _, err := conn.ExecContext(tc.ctx, "SET @omni_temporal = ?", c.sql); err != nil {
			t.Fatalf("binding %q: %v", c.sql, err)
		}
		_, err := conn.ExecContext(tc.ctx, "PREPARE omni_temporal_stmt FROM @omni_temporal")
		if err == nil {
			_, _ = conn.ExecContext(tc.ctx, "DEALLOCATE PREPARE omni_temporal_stmt")
		}
		accepts, err := tidbTemporalVerdict(err)
		if err != nil {
			t.Fatalf("no verdict for %q: %v", c.sql, err)
		}
		if accepts != c.accept {
			t.Errorf("TiDB accepts %q = %v, this file records %v", c.sql, accepts, c.accept)
		}
	}
}

// tidbTemporalVerdict reads the error PREPARE returned. The parser's 1064 carries a
// position: `line 1 column 16 near "(1)"`. A builtin that refuses a value when the
// call is type-checked answers 1064 too, for now(-1) and `curtime`(-1), and that
// message has no position. The second is no syntax error: a quoted name is not even
// the keyword, and the grammar lists NOW with an ExpressionListOpt (parser.y:8123,
// parser.y:8160), and neither production's action refuses anything. Prepared with
// "\n; SELECT FROM" after it, each of these answers the tail's 1064 at line 2, so
// the parser read past the call. It counts as parsed. Anything that is not a server
// error is no verdict.
func tidbTemporalVerdict(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	var myErr *gomysql.MySQLError
	if !errors.As(err, &myErr) {
		return false, err
	}
	return myErr.Number != 1064 || !strings.Contains(myErr.Message, " near "), nil
}
