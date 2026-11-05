package catalog

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/oceanbase/ast"
	"github.com/bytebase/omni/oceanbase/deparse"
	"github.com/bytebase/omni/oceanbase/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

// TestJSONValueClausesRoundTrip: the clauses change what JSON_VALUE answers, so the analyzed
// query keeps them, and what it deparses to reads back to the same text.
func TestJSONValueClausesRoundTrip(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ sql, want string }{
		{"SELECT JSON_VALUE('{}', '$.a')", "json_value('{}', '$.a' RETURNING char(512) charset utf8mb4)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING SIGNED DEFAULT 7 ON EMPTY)", "RETURNING signed default 7 on empty)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -7 ON EMPTY DEFAULT -1.5 ON ERROR)", "default -7 on empty default -1.5 on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' NULL ON EMPTY ERROR ON ERROR)", "null on empty error on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT 'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT N'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT '001' ON ERROR)", "default '001' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT DATE'2000-01-01' ON ERROR)", "default DATE '2000-01-01' on error)"},
		{"SELECT JSON_VALUE((1 IN (SELECT 1)), '$' NULL ON EMPTY)", "json_value((1 in (select 1 as `1`)), '$'"},
		{"SELECT JSON_VALUE((1 IN (1, 2)), '$')", "json_value((1 in (1, 2)), '$'"},
		{"SELECT JSON_VALUE((EXISTS (SELECT 1)), '$')", "json_value((exists (select 1 as `1`)), '$'"},
		{"SELECT JSON_VALUE((1 = 1), '$')", "json_value((1 = 1), '$'"},
		{"SELECT JSON_VALUE((1 BETWEEN 0 AND 2), '$')", "json_value((1 between 0 and 2), '$'"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT 0x61 ON ERROR)", "default X'61' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT b'1' ON ERROR)", "default X'01' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT TRUE ON ERROR)", "default true on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT FALSE ON ERROR)", "default false on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DECIMAL(5,2))", "RETURNING number(5,2))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DATETIME(3))", "RETURNING datetime(3))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING CHAR(3))", "RETURNING char(3)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING UNSIGNED)", "RETURNING unsigned)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING JSON)", "RETURNING json)"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			_, first, second := roundTrip(t, c, tc.sql)
			if !strings.Contains(first, tc.want) {
				t.Errorf("DeparseQuery(%q) = %q, want it to hold %q", tc.sql, first, tc.want)
			}
			if again := DeparseQuery(second); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
		})
	}
}

// TestJSONValueReturningIsTheResultType: RETURNING says what the call answers.
func TestJSONValueReturningIsTheResultType(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	q, _, _ := roundTrip(t, c, "SELECT JSON_VALUE('{}', '$.a' RETURNING SIGNED) AS x")
	call, ok := q.TargetList[0].Expr.(*FuncCallExprQ)
	if !ok {
		t.Fatalf("target = %T", q.TargetList[0].Expr)
	}
	if call.ResultType == nil || call.ResultType.BaseType != BaseTypeBigInt {
		t.Errorf("ResultType = %+v, want a BIGINT", call.ResultType)
	}
}

// jsonValueStoredForms are the clauses this package has to print as the server stores them.
// RETURNING DECIMAL(0) is not one: OceanBase 4.3.5.6 answers 4016, an internal error, for it in a
// generated column.
var jsonValueStoredForms = []string{
	"", "RETURNING SIGNED", "RETURNING SIGNED INTEGER", "RETURNING UNSIGNED", "RETURNING UNSIGNED INTEGER",
	"RETURNING CHAR", "RETURNING CHAR(3)", "RETURNING CHAR(3) CHARSET ascii", "RETURNING CHAR CHARSET ascii",
	"RETURNING CHAR(3) BINARY", "RETURNING BINARY", "RETURNING BINARY(3)", "RETURNING NCHAR(3)", "RETURNING CHAR(3) CHARSET utf8mb3", "RETURNING CHAR(3) CHARSET utf8",
	"RETURNING DECIMAL", "RETURNING DECIMAL(5)", "RETURNING DECIMAL(5,2)", "RETURNING NUMBER", "RETURNING NUMBER(5,2)", "RETURNING NUMERIC(5,2)", "RETURNING FLOAT",
	"RETURNING FLOAT(5)", "RETURNING FLOAT(30)", "RETURNING DOUBLE", "RETURNING DATE", "RETURNING TIME",
	"RETURNING TIME(2)", "RETURNING DATETIME", "RETURNING DATETIME(3)", "RETURNING YEAR", "RETURNING JSON",
	"RETURNING CHAR(3) ASCII", "RETURNING CHAR(3) CHARSET ascii ASCII", "RETURNING CHAR(3) TRUNCATE", "RETURNING CHAR(3) TRUNCATE ASCII",
	"NULL ON EMPTY", "ERROR ON EMPTY", "NULL ON ERROR", "ERROR ON ERROR", "NULL ON EMPTY ERROR ON ERROR",
	"DEFAULT 'x' ON EMPTY", "DEFAULT '001' ON EMPTY", "DEFAULT '' ON ERROR", "DEFAULT 'it''s' ON EMPTY",
	`DEFAULT 'a\\b' ON EMPTY`, "DEFAULT _utf8mb4'x' ON EMPTY", "DEFAULT N'x' ON EMPTY", "DEFAULT 7 ON EMPTY",
	"DEFAULT -7 ON EMPTY", "DEFAULT +7 ON EMPTY", "DEFAULT 1.50 ON EMPTY", "DEFAULT -1.5 ON EMPTY", "DEFAULT .5 ON EMPTY",
	"DEFAULT x'61' ON EMPTY", "DEFAULT X'6A6b' ON EMPTY", "DEFAULT 0xABC ON EMPTY", "DEFAULT 0x6a ON EMPTY",
	"DEFAULT b'1' ON EMPTY", "DEFAULT b'101010101' ON EMPTY", "DEFAULT 0b11 ON EMPTY", "DEFAULT DATE'2000-01-01' ON EMPTY",
	"DEFAULT TIME'10:00:00' ON EMPTY", "DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR",
	"RETURNING CHAR(0)", "RETURNING BINARY(0)", "DEFAULT 18446744073709551615 ON EMPTY",
	"DEFAULT -9223372036854775808 ON EMPTY", "DEFAULT 99999999999999999999 ON EMPTY", "DEFAULT 007 ON EMPTY",
	"DEFAULT b'' ON EMPTY", "DEFAULT b'0000000001000001' ON EMPTY", "DEFAULT x'' ON EMPTY", "DEFAULT 0x0041 ON EMPTY",
	"DEFAULT _latin1 X'61' ON EMPTY", "DEFAULT _latin1 b'01100001' ON EMPTY", "DEFAULT _utf8mb4 0x61 ON EMPTY",
	"RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR",
}

// TestOracle_JSONValueStoredForms: a generated column and a view hold the call the way OceanBase
// stores it, so no clause goes missing and SHOW CREATE TABLE agrees with the engine.
func TestOracle_JSONValueStoredForms(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	stored := regexp.MustCompile(`json_value\(.*\)`)

	for _, clauses := range jsonValueStoredForms {
		t.Run(clauses, func(t *testing.T) {
			table := "CREATE TABLE g (a JSON, b VARCHAR(600) AS (JSON_VALUE(a, '$.b' " + clauses + ")))"
			server, err := storeAndRead(ctx, o, "zz_json_value", table, "g")
			if err != nil {
				t.Fatalf("the server refused %q: %v", table, err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil); err != nil {
				t.Fatal(err)
			}
			got, want := stored.FindString(c.ShowCreateTable("d", "g")), stored.FindString(server)
			if got != want {
				t.Errorf("generated column:\nomni:   %s\nserver: %s", got, want)
			}

			viewCall := "JSON_VALUE('{}', '$.b' " + clauses + ")"
			if _, err := o.db.ExecContext(ctx, "CREATE VIEW v AS SELECT "+viewCall+" AS x"); err != nil {
				t.Fatal(err)
			}
			var name, serverView, charset, collation string
			row := o.db.QueryRowContext(ctx, "SHOW CREATE VIEW v")
			if err := row.Scan(&name, &serverView, &charset, &collation); err != nil {
				t.Fatal(err)
			}
			if call := deparsedCall(t, c, viewCall); !strings.Contains(serverView, call) {
				t.Errorf("view:\nomni:   %s\nserver: %s", call, serverView)
			}
		})
	}
}

// deparsedCall is the text the parser's deparser writes for the call.
func deparsedCall(t *testing.T, c *Catalog, call string) string {
	t.Helper()
	stmts, err := parser.Parse("SELECT " + call)
	if err != nil {
		t.Fatal(err)
	}
	return deparse.Deparse(stmts.Items[0].(*nodes.SelectStmt).TargetList[0])
}

// TestJSONValueStoredFormsReadBack: what this package writes for a form reads back to the same text.
func TestJSONValueStoredFormsReadBack(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	for _, clauses := range jsonValueStoredForms {
		t.Run(clauses, func(t *testing.T) {
			_, first, second := roundTrip(t, c, "SELECT JSON_VALUE('{}', '$.b' "+clauses+")")
			if again := DeparseQuery(second); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
		})
	}
}

// TestOracle_NumberIsDecimal: a NUMBER column is a DECIMAL one, on the server and here.
func TestOracle_NumberIsDecimal(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	stored := regexp.MustCompile("`a` [^,\\n]*")
	for _, column := range []string{"NUMBER", "NUMBER(5)", "NUMBER(5,2)", "NUMBER(5,2) UNSIGNED"} {
		t.Run(column, func(t *testing.T) {
			table := "CREATE TABLE g (a " + column + ")"
			server, err := storeAndRead(ctx, o, "zz_number", table, "g")
			if err != nil {
				t.Fatal(err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil); err != nil {
				t.Fatal(err)
			}
			if got, want := stored.FindString(c.ShowCreateTable("d", "g")), stored.FindString(server); got != want {
				t.Errorf("omni: %s\nserver: %s", got, want)
			}
		})
	}
}

// TestOracle_JSONValueFunctionalIndex: a functional index over JSON_VALUE is taken or refused as the
// server takes it, with the server's code.
func TestOracle_JSONValueFunctionalIndex(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for _, returning := range []string{"", "RETURNING SIGNED", "RETURNING CHAR(10)", "RETURNING DECIMAL(5,2)",
		"RETURNING DATETIME(3)", "RETURNING DOUBLE", "RETURNING JSON", "RETURNING CHAR", "RETURNING CHAR(512)",
		"RETURNING CHAR(513)", "RETURNING BINARY(3)", "RETURNING YEAR"} {
		t.Run(returning, func(t *testing.T) {
			table := "CREATE TABLE f (a JSON, KEY k ((JSON_VALUE(a, '$.x' " + returning + "))))"
			serverCode := 0
			if _, err := storeAndRead(ctx, o, "zz_json_value_index", table, "f"); err != nil {
				var merr *gomysql.MySQLError
				if !errors.As(err, &merr) {
					t.Fatal(err)
				}
				serverCode = int(merr.Number)
			}
			c := New()
			results, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil)
			if err != nil {
				t.Fatal(err)
			}
			catalogCode := 0
			var cerr *Error
			if errors.As(results[len(results)-1].Error, &cerr) {
				catalogCode = cerr.Code
			}
			if catalogCode != serverCode {
				t.Errorf("catalog answered %d, the server %d", catalogCode, serverCode)
			}
		})
	}
}

// TestStringLiteralKeepsItsKind: a string stays a string through analysis, introducer and all,
// however its text looks: '001' is no number, and _latin1 X'303031' is the string _latin1'001'.
func TestStringLiteralKeepsItsKind(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	for sql, want := range map[string]string{
		"SELECT '001' AS x":               "select '001' as `x`",
		"SELECT '1.5' AS x":               "select '1.5' as `x`",
		"SELECT 'TRUE' AS x":              "select 'TRUE' as `x`",
		"SELECT _latin1'001' AS x":        "select _latin1'001' as `x`",
		"SELECT _latin1 X'303031' AS x":   "select _latin1'001' as `x`",
		"SELECT _latin1 b'00110001' AS x": "select _latin1'1' as `x`",
		"SELECT _latin1 X'5c' AS x":       "select _latin1'\\\\' as `x`",
		"SELECT _latin1 X'5c6e' AS x":     "select _latin1'\\\\n' as `x`",
		"SELECT 'a\\\\b' AS x":            "select 'a\\\\b' as `x`",
		"SELECT 'it''s' AS x":             "select 'it\\'s' as `x`",
		"SELECT 1 AS x":                   "select 1 as `x`",
		"SELECT 1.5 AS x":                 "select 1.5 as `x`",
	} {
		t.Run(sql, func(t *testing.T) {
			_, first, second := roundTrip(t, c, sql)
			if first != want {
				t.Errorf("DeparseQuery(%q) = %q, want %q", sql, first, want)
			}
			if again := DeparseQuery(second); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
		})
	}
}
