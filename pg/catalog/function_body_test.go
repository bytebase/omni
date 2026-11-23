package catalog

import "testing"

// TestSQLFunctionBodyAccepts holds the checks of a SQL function's body to
// what PostgreSQL's regression suite accepts.
func TestSQLFunctionBodyAccepts(t *testing.T) {
	for _, sql := range []string{
		// The validator runs once the function exists, so a body may call it.
		`CREATE FUNCTION fact(n int) RETURNS int LANGUAGE sql AS 'SELECT CASE WHEN n <= 1 THEN 1 ELSE n * fact(n - 1) END'`,
		// The catalog holds no system catalog, which a body may read.
		`CREATE FUNCTION error1(text) RETURNS text LANGUAGE sql AS $$ SELECT relname::text FROM pg_class c WHERE c.oid = $1::regclass $$`,
		// A polymorphic operator's result is not resolved, so it is not compared.
		`CREATE FUNCTION tfnp(int[], int) RETURNS int[] AS 'select $1 || $2' LANGUAGE SQL`,
		// pg_dump turns the checks off to create a function before its table.
		`SET check_function_bodies = false; CREATE FUNCTION later() RETURNS bigint LANGUAGE sql AS 'SELECT count(*) FROM made_later'`,
		`SET check_function_bodies = false; CREATE FUNCTION later() RETURNS int LANGUAGE plpgsql AS $$ BEGIN RETURN 1 END $$`,
		// CALL passes NULL for an OUT parameter.
		`CREATE PROCEDURE p(OUT n int) LANGUAGE plpgsql AS $$ BEGIN n := 1; END $$; CALL p(NULL)`,
		// A quoted mixed-case parameter is assignable.
		`CREATE FUNCTION f("X" int) RETURNS int LANGUAGE plpgsql AS $$ BEGIN "X" := 2; RETURN "X"; END $$`,
		// SELECT INTO may target a record field.
		`CREATE TABLE t(x int); CREATE FUNCTION f() RETURNS int LANGUAGE plpgsql AS $$ DECLARE r t%ROWTYPE; BEGIN SELECT 1 INTO r.x; RETURN r.x; END $$`,
		// A replacement's body sees the defaults it adds.
		`CREATE FUNCTION f(n int) RETURNS int LANGUAGE sql AS 'SELECT n'; CREATE OR REPLACE FUNCTION f(n int DEFAULT 0) RETURNS int LANGUAGE sql AS 'SELECT CASE WHEN n > 0 THEN f() ELSE 0 END'`,
	} {
		results, err := New().Exec(sql, nil)
		for _, r := range results {
			if err == nil {
				err = r.Error
			}
		}
		if err != nil {
			t.Errorf("%s: %v", sql, err)
		}
	}
}
