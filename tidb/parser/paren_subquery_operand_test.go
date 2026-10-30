package parser

import (
	"errors"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"
)

// parenSubqueryOperandAccepted are parentheses around an expression whose
// first operand is a parenthesized subquery, and the derived-table join
// groups TiDB writes in view definitions. TiDB v8.5 accepts them all
// (TestParenSubqueryOperandOracle).
var parenSubqueryOperandAccepted = []string{
	"select (((select 1) + 0))",
	"select (((select 1) = 0))",
	"select ((((select 1) = 0)))",
	"select ((((select 1) = 0)) or (((select 2)) = 2))",
	"select (((select count(0) from n1) = 0))",
	"select 1 in (((select 1) = 1))",
	"select 1 in ((((select 1)) = 1))",
	"select if((((select count(0) from n1) = 0)), 'NO', 'YES')",
	"select ((select 1) + 0)",
	"select (((select 1)) + 0)",
	// a parenthesized query stays a subquery
	"select (((select 1)))",
	"select (((select 1) union (select 2)))",
	"select 1 in (((select 1)))",
	"select 1 in (((select 1) union (select 2)))",
	// join groups as TiDB writes them in view definitions
	"SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `test`.`j1`) AS `d` JOIN `test`.`j2` ON `d`.`k1`=`j2`.`k1`) JOIN `test`.`j3` ON `j3`.`k1`=`d`.`k1`",
	"SELECT `d`.`k1` AS `k1` FROM `test`.`j2` JOIN ((SELECT `k1` AS `k1` FROM `test`.`j1`) AS `d` JOIN `test`.`j3` USING (`k1`)) ON `j2`.`k1`=`d`.`k1`",
	"SELECT `d`.`k1` AS `k1` FROM ((SELECT `k1` AS `k1` FROM `test`.`j1`) UNION (SELECT `k1` AS `k1` FROM `test`.`j2`)) AS `d` JOIN `test`.`j3` ON `d`.`k1`=`j3`.`k1`",
}

// parenSubqueryOperandRejected are rejected by TiDB v8.5 with 1064.
var parenSubqueryOperandRejected = []string{
	"select (((select 1) + 0)",
	"select (((select 1) +))",
	"select (((select 1) d))",
	"select (((select 1) union (select 2) + 0))",
}

func TestParenSubqueryOperand(t *testing.T) {
	for _, sql := range parenSubqueryOperandAccepted {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
	for _, sql := range parenSubqueryOperandRejected {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestParenSubqueryOperandOracle asks TiDB v8.5 about every statement of
// TestParenSubqueryOperand.
func TestParenSubqueryOperandOracle(t *testing.T) {
	tc := startTiDB(t)
	// One connection in a database of its own, so USE holds and the tables
	// and views stay out of the other tests' way.
	conn, err := tc.db.Conn(tc.ctx)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer conn.Close()
	exec := func(sql string) error {
		_, err := conn.ExecContext(tc.ctx, sql)
		return err
	}
	defer exec("DROP DATABASE IF EXISTS omni_paren_operand")
	for _, ddl := range []string{
		"DROP DATABASE IF EXISTS omni_paren_operand",
		"CREATE DATABASE omni_paren_operand",
		"USE omni_paren_operand",
		"CREATE TABLE j1 (k1 INT, a INT)",
		"CREATE TABLE j2 (k1 INT, b INT)",
		"CREATE TABLE j3 (k1 INT, c INT)",
		"CREATE TABLE n1 (a INT)",
	} {
		if err := exec(ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	check := func(sql string, wantAccept bool) {
		t.Helper()
		err := exec(sql)
		var me *gomysql.MySQLError
		accepted := err == nil || (errors.As(err, &me) && me.Number != 1064)
		if err != nil && !errors.As(err, &me) {
			t.Fatalf("TiDB errored on %q: %v", sql, err)
		}
		if accepted != wantAccept {
			t.Errorf("TiDB on %q: accepted %v (%v), want %v", sql, accepted, err, wantAccept)
		}
		if _, err := Parse(sql); (err == nil) != wantAccept {
			t.Errorf("Parse(%q) = %v, want accepted %v", sql, err, wantAccept)
		}
	}
	for _, sql := range parenSubqueryOperandAccepted {
		check(sql, true)
	}
	for _, sql := range parenSubqueryOperandRejected {
		check(sql, false)
	}

	// TiDB writes a view definition back in its own spelling; each must parse.
	for _, ddl := range []string{
		"CREATE OR REPLACE VIEW v1 AS SELECT d.k1 FROM (SELECT k1 FROM j1) d JOIN j2 ON d.k1 = j2.k1 JOIN j3 ON j3.k1 = d.k1",
		"CREATE OR REPLACE VIEW v2 AS SELECT d.k1 FROM j2 JOIN ((SELECT k1 FROM j1) d JOIN j3 USING (k1)) ON j2.k1 = d.k1",
		"CREATE OR REPLACE VIEW v3 AS SELECT d.k1 FROM ((SELECT k1 FROM j1) UNION (SELECT k1 FROM j2)) d JOIN j3 ON d.k1 = j3.k1",
		"CREATE OR REPLACE VIEW v4 AS SELECT (((SELECT a FROM n1 LIMIT 1) + 0)) AS x, if((((SELECT count(0) FROM n1) = 0)), 'NO', 'YES') AS y",
		"CREATE OR REPLACE VIEW v5 AS SELECT a FROM n1 WHERE a IN ((((SELECT a FROM n1 LIMIT 1) + 0)), 2)",
	} {
		if err := exec(ddl); err != nil {
			t.Fatalf("%s: %v", ddl, err)
		}
		name := strings.Fields(ddl)[4]
		var view, createView, charset, collation string
		if err := conn.QueryRowContext(tc.ctx, "SHOW CREATE VIEW "+name).Scan(&view, &createView, &charset, &collation); err != nil {
			t.Fatalf("SHOW CREATE VIEW %s: %v", name, err)
		}
		if _, err := Parse(createView); err != nil {
			t.Errorf("Parse(%q): %v", createView, err)
		}
	}
}

// TestParenSubqueryOperandStaysLinear: the verdict of every block the
// lookahead passes is recorded, so asking at every level of a scalar paren
// layer around a subquery operand costs nothing.
func TestParenSubqueryOperandStaysLinear(t *testing.T) {
	raiseNestingLimit(t)
	const depth = 20000
	for _, tc := range []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"subquery operand in parentheses", "SELECT " + nest(depth, "(", "(SELECT 1) = 0", ")"), false},
		{"subquery operand in an IN list", "SELECT 1 IN " + nest(depth, "(", "(SELECT 1) = 0", ")"), false},
		{"unclosed subquery operand in parentheses", "SELECT " + strings.Repeat("(", depth) + "(SELECT 1) = 0", true},
	} {
		start := time.Now()
		if _, err := Parse(tc.sql); (err != nil) != tc.wantErr {
			t.Fatalf("%s: Parse error = %v, want error %v", tc.name, briefErr(err), tc.wantErr)
		}
		if elapsed := time.Since(start); elapsed > 2*time.Second {
			t.Errorf("%s: Parse through %d levels took %v", tc.name, depth, elapsed)
		}
	}
}
