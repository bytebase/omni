package parser_test

import (
	"fmt"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
	plpgsqlparser "github.com/bytebase/omni/pg/plpgsql/parser"
)

func TestCreateFunctionBodyExtractionForms(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantLang   string
		wantBodies []string
		wantSQLStd bool
	}{
		{
			name:       "as before language",
			sql:        `CREATE FUNCTION f() RETURNS int AS $$SELECT 1$$ LANGUAGE sql`,
			wantLang:   "sql",
			wantBodies: []string{"SELECT 1"},
		},
		{
			name:       "language before tagged dollar body",
			sql:        `CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $body$BEGIN RETURN 1; END$body$`,
			wantLang:   "plpgsql",
			wantBodies: []string{"BEGIN RETURN 1; END"},
		},
		{
			name:       "quoted language and quoted body",
			sql:        `CREATE FUNCTION f() RETURNS int AS 'SELECT 1' LANGUAGE 'sql'`,
			wantLang:   "sql",
			wantBodies: []string{"SELECT 1"},
		},
		{
			name:       "procedure plpgsql body",
			sql:        `CREATE PROCEDURE p() LANGUAGE plpgsql AS $$BEGIN NULL; END$$`,
			wantLang:   "plpgsql",
			wantBodies: []string{"BEGIN NULL; END"},
		},
		{
			name:       "multi string as clause",
			sql:        `CREATE FUNCTION f() RETURNS int AS 'objfile', 'link_symbol' LANGUAGE c`,
			wantLang:   "c",
			wantBodies: []string{"objfile", "link_symbol"},
		},
		{
			name:       "sql standard return body",
			sql:        `CREATE FUNCTION f() RETURNS int LANGUAGE sql RETURN 1`,
			wantLang:   "sql",
			wantSQLStd: true,
		},
		{
			name:       "sql standard atomic body",
			sql:        `CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1; END`,
			wantLang:   "sql",
			wantSQLStd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := parseCreateFunction(t, tt.sql)
			if got := functionLanguage(fn); got != tt.wantLang {
				t.Fatalf("language = %q, want %q", got, tt.wantLang)
			}
			if got := functionASBodies(fn); !sameStrings(got, tt.wantBodies) {
				t.Fatalf("AS bodies = %#v, want %#v", got, tt.wantBodies)
			}
			if got := fn.SqlBody != nil; got != tt.wantSQLStd {
				t.Fatalf("SqlBody present = %v, want %v", got, tt.wantSQLStd)
			}
		})
	}
}

func TestLanguageSQLStringBodyCompatibility(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"select scalar", `SELECT 1`},
		{"select parameter", `SELECT $1 + 1`},
		{"qualified argument reference", `SELECT f.x`},
		{"cte select", `WITH x AS (SELECT 1 AS a) SELECT a FROM x`},
		{"insert returning", `INSERT INTO t(id) VALUES (1) RETURNING id`},
		{"update returning", `UPDATE t SET v = v + 1 WHERE id = 1 RETURNING v`},
		{"delete returning", `DELETE FROM t WHERE id = 1 RETURNING id`},
		{"modifying cte", `WITH moved AS (DELETE FROM t WHERE id > 10 RETURNING *) INSERT INTO archive SELECT * FROM moved`},
		{"multi statement void style", `CREATE TABLE tmp(id int); INSERT INTO tmp VALUES (1); SELECT 1`},
		{"case expression", `SELECT CASE WHEN $1 > 0 THEN $1 ELSE 0 END`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := "CREATE FUNCTION f(x int) RETURNS int LANGUAGE sql AS $fn$" + tt.body + "$fn$"
			fn := parseCreateFunction(t, sql)
			body := singleASBody(t, fn)
			if _, err := pgparser.Parse(body); err != nil {
				t.Fatalf("LANGUAGE sql body parse failed: %v", err)
			}
		})
	}
}

func TestLanguageSQLRejectsPLpgSQLOnlyBody(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"plpgsql block", `BEGIN RETURN 1; END`},
		{"declaration block", `DECLARE x int; BEGIN RETURN x; END`},
		{"if statement", `IF true THEN SELECT 1; END IF`},
		{"loop statement", `LOOP SELECT 1; END LOOP`},
		{"exception clause", `EXCEPTION WHEN others THEN SELECT 1`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := "CREATE FUNCTION f() RETURNS int LANGUAGE sql AS $fn$" + tt.body + "$fn$"
			fn := parseCreateFunction(t, sql)
			body := singleASBody(t, fn)
			if _, err := pgparser.Parse(body); err == nil {
				t.Fatalf("LANGUAGE sql body parsed successfully, want error")
			}
		})
	}
}

func TestSQLStandardFunctionBodyRejectsMalformedBody(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "return without expression",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE sql RETURN`,
		},
		{
			name: "atomic return without expression",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN; END`,
		},
		{
			name: "atomic invalid sql statement",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC SELECT FROM; END`,
		},
		{
			name: "atomic body missing atomic keyword",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN RETURN 1; END`,
		},
		{
			name: "atomic body missing end keyword",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE sql BEGIN ATOMIC RETURN 1;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := pgparser.Parse(tt.sql); err == nil {
				t.Fatalf("Parse succeeded, want error")
			}
		})
	}
}

func TestCreateFunctionRejectsMalformedOptions(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "as without string body",
			sql:  `CREATE FUNCTION f() RETURNS int AS LANGUAGE sql`,
		},
		{
			name: "as second string missing",
			sql:  `CREATE FUNCTION f() RETURNS int AS 'objfile', LANGUAGE c`,
		},
		{
			name: "language without name",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE`,
		},
		{
			name: "language name cannot be AS keyword",
			sql:  `CREATE FUNCTION f() RETURNS int LANGUAGE AS 'SELECT 1'`,
		},
		{
			name: "transform missing for type",
			sql:  `CREATE FUNCTION f() RETURNS int TRANSFORM int LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "called missing on",
			sql:  `CREATE FUNCTION f() RETURNS int CALLED NULL INPUT LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "security missing mode",
			sql:  `CREATE FUNCTION f() RETURNS int SECURITY LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "not without leakproof",
			sql:  `CREATE FUNCTION f() RETURNS int NOT STABLE LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set time missing zone",
			sql:  `CREATE FUNCTION f() RETURNS int SET TIME 'UTC' LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set catalog requires string",
			sql:  `CREATE FUNCTION f() RETURNS int SET CATALOG current_catalog LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set session requires authorization",
			sql:  `CREATE FUNCTION f() RETURNS int SET SESSION ROLE LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set xml requires option",
			sql:  `CREATE FUNCTION f() RETURNS int SET XML DOCUMENT LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set transaction requires snapshot",
			sql:  `CREATE FUNCTION f() RETURNS int SET TRANSACTION 'abc' LANGUAGE sql AS 'SELECT 1'`,
		},
		{
			name: "set from requires current",
			sql:  `CREATE FUNCTION f() RETURNS int SET search_path FROM 'public' LANGUAGE sql AS 'SELECT 1'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := pgparser.Parse(tt.sql); err == nil {
				t.Fatalf("Parse succeeded, want error")
			}
		})
	}
}

func TestLanguagePLpgSQLStringBodyCompatibility(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"empty block", `BEGIN END`},
		{"declarations", `DECLARE x character varying(255) := 'a'; y timestamp with time zone; BEGIN NULL; END`},
		{"assignment", `BEGIN x := y + 1; END`},
		{"perform", `BEGIN PERFORM do_work(1, 2); END`},
		{"select into", `BEGIN SELECT a, b INTO x, y FROM t WHERE id = 1; END`},
		{"returning into", `BEGIN INSERT INTO t(a) VALUES (1) RETURNING id INTO new_id; END`},
		{"utility sql", `BEGIN REFRESH MATERIALIZED VIEW CONCURRENTLY mv; GRANT USAGE ON SCHEMA api TO web_anon; END`},
		{"if case", `BEGIN IF x > 0 THEN y := 1; ELSIF x = 0 THEN y := 0; ELSE y := -1; END IF; CASE y WHEN 1 THEN NULL; ELSE NULL; END CASE; END`},
		{"loops", `BEGIN WHILE x < 10 LOOP x := x + 1; END LOOP; FOR i IN 1..10 BY 2 LOOP NULL; END LOOP; END`},
		{"query for target list", `BEGIN FOR a, b IN SELECT x, y FROM t LOOP NULL; END LOOP; END`},
		{"dynamic return query", `BEGIN RETURN QUERY EXECUTE 'SELECT * FROM ' || quote_ident(tbl) USING id; END`},
		{"exception block", `BEGIN INSERT INTO t VALUES (1); EXCEPTION WHEN unique_violation THEN NULL; WHEN others THEN RAISE; END`},
		{"cursor", `DECLARE c CURSOR FOR SELECT id FROM t; v int; BEGIN OPEN c; FETCH c INTO v; CLOSE c; END`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $fn$" + tt.body + "$fn$"
			fn := parseCreateFunction(t, sql)
			body := singleASBody(t, fn)
			if _, err := plpgsqlparser.Parse(body); err != nil {
				t.Fatalf("LANGUAGE plpgsql body parse failed: %v", err)
			}
		})
	}
}

func TestLanguagePLpgSQLRejectsMalformedBody(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"not null without default", `DECLARE x int NOT NULL; BEGIN END`},
		{"missing end if", `BEGIN IF true THEN NULL; END`},
		{"for missing loop", `BEGIN FOR i IN 1..10 NULL; END LOOP; END`},
		{"execute missing query", `BEGIN EXECUTE; END`},
		{"select into missing target", `BEGIN SELECT 1 INTO; END`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := "CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $fn$" + tt.body + "$fn$"
			fn := parseCreateFunction(t, sql)
			body := singleASBody(t, fn)
			if _, err := plpgsqlparser.Parse(body); err == nil {
				t.Fatalf("LANGUAGE plpgsql body parsed successfully, want error")
			}
		})
	}
}

func parseCreateFunction(t *testing.T, sql string) *nodes.CreateFunctionStmt {
	t.Helper()
	fn, err := pgParseCreateFunction(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	return fn
}

func pgParseWithBodyCheck(sql string) (*nodes.CreateFunctionStmt, error) {
	fn, err := pgParseCreateFunction(sql)
	if err != nil {
		return nil, err
	}
	bodies := functionASBodies(fn)
	if len(bodies) != 1 {
		return fn, nil
	}
	switch functionLanguage(fn) {
	case "sql":
		if _, err := pgparser.Parse(bodies[0]); err != nil {
			return nil, err
		}
	case "plpgsql":
		if _, err := plpgsqlparser.Parse(bodies[0]); err != nil {
			return nil, err
		}
	}
	return fn, nil
}

func pgParseCreateFunction(sql string) (*nodes.CreateFunctionStmt, error) {
	stmts, err := pgparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if stmts == nil || len(stmts.Items) != 1 {
		return nil, fmt.Errorf("expected one statement, got %#v", stmts)
	}
	raw, ok := stmts.Items[0].(*nodes.RawStmt)
	if !ok {
		return nil, fmt.Errorf("expected RawStmt, got %T", stmts.Items[0])
	}
	fn, ok := raw.Stmt.(*nodes.CreateFunctionStmt)
	if !ok {
		return nil, fmt.Errorf("expected CreateFunctionStmt, got %T", raw.Stmt)
	}
	return fn, nil
}

func functionLanguage(fn *nodes.CreateFunctionStmt) string {
	if fn.Options == nil {
		return ""
	}
	for _, item := range fn.Options.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok || !strings.EqualFold(d.Defname, "language") {
			continue
		}
		if s, ok := d.Arg.(*nodes.String); ok {
			return strings.ToLower(s.Str)
		}
	}
	return ""
}

func functionASBodies(fn *nodes.CreateFunctionStmt) []string {
	if fn.Options == nil {
		return nil
	}
	for _, item := range fn.Options.Items {
		d, ok := item.(*nodes.DefElem)
		if !ok || !strings.EqualFold(d.Defname, "as") {
			continue
		}
		l, ok := d.Arg.(*nodes.List)
		if !ok {
			return nil
		}
		var out []string
		for _, body := range l.Items {
			s, ok := body.(*nodes.String)
			if !ok {
				continue
			}
			out = append(out, s.Str)
		}
		return out
	}
	return nil
}

func singleASBody(t *testing.T, fn *nodes.CreateFunctionStmt) string {
	t.Helper()
	bodies := functionASBodies(fn)
	if len(bodies) != 1 {
		t.Fatalf("expected exactly one AS body, got %#v", bodies)
	}
	return bodies[0]
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
