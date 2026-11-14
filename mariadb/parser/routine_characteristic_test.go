package parser

import "testing"

// TestRoutineCharacteristicWords holds a routine's characteristics, a loadable function's
// return type and EXECUTE ... USING to the words the grammar spells (sp_chistic, sp_suid,
// udf_type, execute_var_ident). code is what mariadb:11.8.8 answers to PREPARE of the statement, 0
// when the compat judge reads that as parsed (compat.PrepareVerdict).
func TestRoutineCharacteristicWords(t *testing.T) {
	cases := []struct {
		sql  string
		code int
	}{
		{sql: "CREATE FUNCTION f RETURNS STRING SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS REAL SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS DECIMAL SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS DEC SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS INT SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS INTEGER SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS INT4 SONAME 'x.so'"},
		{sql: "CREATE FUNCTION f RETURNS NUMERIC SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS FIXED SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS DOUBLE SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS CHAR SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS `INT` SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS 'INT' SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS BIGINT SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS FLOAT8 SONAME 'x.so'", code: 1064},
		{sql: "CREATE FUNCTION f RETURNS FLOAT SONAME 'x.so'", code: 1064},
		{sql: "ALTER FUNCTION f SQL SECURITY INVOKER"},
		{sql: "ALTER FUNCTION f LANGUAGE SQL"},
		{sql: "ALTER FUNCTION f NO SQL"},
		{sql: "ALTER FUNCTION f CONTAINS SQL"},
		{sql: "ALTER FUNCTION f READS SQL DATA"},
		{sql: "ALTER FUNCTION f MODIFIES SQL DATA"},
		{sql: "ALTER FUNCTION f COMMENT 'c' NO SQL"},
		{sql: "ALTER FUNCTION f DETERMINISTIC", code: 1064},
		{sql: "CREATE PROCEDURE p() CONTAINS SELECT 1", code: 1064},
		{sql: "CREATE PROCEDURE p() READS SQL SELECT 1", code: 1064},
		{sql: "CREATE PROCEDURE p() SQL SELECT 1", code: 1064},
		{sql: "CREATE PROCEDURE p() MODIFIES SQL DATA SELECT 1"},
		{sql: "EXECUTE s USING @a, @b"},
		{sql: "EXECUTE s USING @'a'"},
		{sql: "EXECUTE s USING @@sql_mode"},
		{sql: "EXECUTE s USING 1"},
		{sql: "EXECUTE s USING a", code: 1054},
		{sql: "EXECUTE s USING @a + 1"},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if c.code == 0 && err != nil {
				t.Fatalf("Parse refused what MariaDB takes: %v", err)
			}
			if c.code != 0 && err == nil {
				t.Fatalf("Parse took what MariaDB refuses with %d", c.code)
			}
		})
	}
}
