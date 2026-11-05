package catalog

import (
	"strings"
	"testing"
)

// TestJSONValueClausesRoundTrip: the clauses change what JSON_VALUE answers, so the analyzed
// query keeps them, and what it deparses to reads back to the same text.
func TestJSONValueClausesRoundTrip(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d", nil); err != nil {
		t.Fatal(err)
	}
	cases := []struct{ sql, want string }{
		{"SELECT JSON_VALUE('{}', '$.a')", "json_value('{}', '$.a' returning char(512))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING SIGNED DEFAULT 7 ON EMPTY)", "returning signed default 7 on empty)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -7 ON EMPTY DEFAULT -1.5 ON ERROR)", "default -7 on empty default -1.5 on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' NULL ON EMPTY ERROR ON ERROR)", "null on empty error on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT 'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT N'x' ON ERROR)", "default 'x' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT '001' ON ERROR)", "default '001' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT DATE'2000-01-01' ON ERROR)", "default DATE'2000-01-01' on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT 0x61 ON ERROR)", "default 0x61 on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' DEFAULT b'1' ON ERROR)", "default 0x01 on error)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DECIMAL(5,2))", "returning decimal(5, 2))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING DATETIME(3))", "returning datetime(3))"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING CHAR(3))", "returning char(3)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING UNSIGNED)", "returning unsigned)"},
		{"SELECT JSON_VALUE('{}', '$.a' RETURNING JSON)", "returning json)"},
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

// jsonValueStoredForms are the clauses this package has to print as the server stores them, each
// with what is behind it in sql_yacc.yy: a RETURNING cast type, then a response's literal.
var jsonValueStoredForms = []string{
	"", "RETURNING SIGNED", "RETURNING SIGNED INTEGER", "RETURNING UNSIGNED", "RETURNING UNSIGNED INTEGER",
	"RETURNING CHAR", "RETURNING CHAR(3)", "RETURNING CHAR(3) CHARSET latin1", "RETURNING CHAR CHARSET latin1",
	"RETURNING CHAR CHARSET utf8mb3", "RETURNING CHAR CHARSET utf8", "RETURNING CHAR CHARSET gbk", "RETURNING CHAR CHARSET utf32",
	"RETURNING CHAR CHARSET binary", "RETURNING CHAR(3) CHARSET binary", "RETURNING CHAR BINARY", "RETURNING NCHAR",
	"RETURNING CHAR(3) BINARY", "RETURNING BINARY", "RETURNING BINARY(3)", "RETURNING NCHAR(3)",
	"RETURNING DECIMAL", "RETURNING DECIMAL(5)", "RETURNING DECIMAL(5,2)", "RETURNING FLOAT", "RETURNING FLOAT(5)",
	"RETURNING FLOAT(30)", "RETURNING DOUBLE", "RETURNING REAL", "RETURNING DATE", "RETURNING TIME", "RETURNING TIME(2)",
	"RETURNING DATETIME", "RETURNING DATETIME(3)", "RETURNING YEAR", "RETURNING JSON",
	"NULL ON EMPTY", "ERROR ON EMPTY", "NULL ON ERROR", "ERROR ON ERROR", "NULL ON EMPTY ERROR ON ERROR",
	"DEFAULT 'x' ON EMPTY", "DEFAULT '001' ON EMPTY", "DEFAULT '' ON ERROR", "DEFAULT 'it''s' ON EMPTY",
	`DEFAULT 'a\\b' ON EMPTY`, "DEFAULT _latin1'x' ON EMPTY", "DEFAULT N'x' ON EMPTY", "DEFAULT 7 ON EMPTY",
	"DEFAULT -7 ON EMPTY", "DEFAULT +7 ON EMPTY", "DEFAULT 1.50 ON EMPTY", "DEFAULT -1.5 ON EMPTY", "DEFAULT .5 ON EMPTY",
	"DEFAULT 1E2 ON EMPTY", "DEFAULT TRUE ON EMPTY", "DEFAULT FALSE ON ERROR", "DEFAULT x'61' ON EMPTY",
	"DEFAULT X'6A6b' ON EMPTY", "DEFAULT 0xABC ON EMPTY", "DEFAULT 0x6a ON EMPTY", "DEFAULT b'1' ON EMPTY",
	"DEFAULT b'101010101' ON EMPTY", "DEFAULT 0b11 ON EMPTY", "DEFAULT DATE'2000-01-01' ON EMPTY",
	"DEFAULT TIME'10:00:00' ON EMPTY", "DEFAULT TIMESTAMP'2000-01-01 00:00:00' ON ERROR",
	"RETURNING SIGNED DEFAULT 7 ON EMPTY DEFAULT -1 ON ERROR",
}

// TestJSONValueStoredForms: a generated column and a view hold the call the way the server stores
// it, clauses and all, so SHOW CREATE TABLE agrees with the engine and no clause goes missing.
func TestJSONValueStoredForms(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	for _, clauses := range jsonValueStoredForms {
		t.Run(clauses, func(t *testing.T) {
			if err := ctr.execSQL("DROP TABLE IF EXISTS g; DROP VIEW IF EXISTS v"); err != nil {
				t.Fatal(err)
			}
			call := "JSON_VALUE(a, '$.b' " + clauses + ")"
			table := "CREATE TABLE g (a JSON, b VARCHAR(600) AS (" + call + "))"
			if err := ctr.execSQL(table); err != nil {
				t.Fatalf("the server refused %q: %v", table, err)
			}
			server, err := ctr.showCreateTable("g")
			if err != nil {
				t.Fatal(err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil); err != nil {
				t.Fatal(err)
			}
			if got := c.ShowCreateTable("d", "g"); got != server {
				t.Errorf("generated column:\nomni:   %s\nserver: %s", got, server)
			}

			viewCall := "JSON_VALUE('{}', '$.b' " + clauses + ")"
			if err := ctr.execSQL("CREATE VIEW v AS SELECT " + viewCall + " AS x"); err != nil {
				t.Fatal(err)
			}
			serverView, err := ctr.showCreateView("v")
			if err != nil {
				t.Fatal(err)
			}
			_, first, second := roundTrip(t, c, "SELECT "+viewCall+" AS x")
			if !strings.Contains(serverView, strings.TrimSuffix(strings.TrimPrefix(first, "select "), " as `x`")) {
				t.Errorf("view:\nomni:   %s\nserver: %s", first, serverView)
			}
			if again := DeparseQuery(second); again != first {
				t.Errorf("deparsed %q, which reads back as %q", first, again)
			}
		})
	}
}

// TestJSONValueStoredInConstraintsAndIndexes: CHECK, a functional index and DEFAULT store the call
// the same way, clauses and all. This package writes strings in those three places without the
// introducer the server adds, for every call and not only this one, so the comparison sets that
// apart.
func TestJSONValueStoredInConstraintsAndIndexes(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	for name, table := range map[string]string{
		"check":      "CREATE TABLE g (a JSON, CONSTRAINT c CHECK (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY) > 0))",
		"index":      "CREATE TABLE g (a JSON, KEY k ((JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY))))",
		"default":    "CREATE TABLE g (a JSON, b INT DEFAULT (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY)))",
		"stored col": "CREATE TABLE g (a JSON, b INT AS (JSON_VALUE(a, '$.b' RETURNING SIGNED DEFAULT 7 ON EMPTY)) STORED)",
	} {
		t.Run(name, func(t *testing.T) {
			if err := ctr.execSQL("DROP TABLE IF EXISTS g"); err != nil {
				t.Fatal(err)
			}
			if err := ctr.execSQL(table); err != nil {
				t.Fatalf("the server refused %q: %v", table, err)
			}
			server, err := ctr.showCreateTable("g")
			if err != nil {
				t.Fatal(err)
			}
			c := New()
			if _, err := c.Exec("CREATE DATABASE d; USE d; "+table, nil); err != nil {
				t.Fatal(err)
			}
			got := strings.ReplaceAll(c.ShowCreateTable("d", "g"), "_utf8mb4'", "'")
			if want := strings.ReplaceAll(server, "_utf8mb4'", "'"); got != want {
				t.Errorf("\nomni:   %s\nserver: %s", got, server)
			}
		})
	}
}
