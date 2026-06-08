package parser

import "testing"

func TestLocBroadScan(t *testing.T) {
	tests := []string{
		// ALTER TABLE
		"ALTER TABLE t ADD COLUMN c int",
		"ALTER TABLE t DROP COLUMN c",
		"ALTER TABLE t ALTER COLUMN c SET NOT NULL",
		"ALTER TABLE t RENAME TO t2",
		// CREATE FUNCTION
		"CREATE FUNCTION f() RETURNS int AS 'SELECT 1' LANGUAGE SQL",
		// CREATE TRIGGER
		"CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW EXECUTE FUNCTION f()",
		// CREATE TYPE
		"CREATE TYPE t AS (a int, b text)",
		// DROP
		"DROP TABLE t",
		"DROP TABLE IF EXISTS t CASCADE",
		"DROP INDEX idx",
		"DROP VIEW v",
		"DROP FUNCTION f()",
		// GRANT/REVOKE
		"GRANT SELECT ON t TO PUBLIC",
		"REVOKE ALL ON t FROM PUBLIC",
		// SCHEMA
		"CREATE SCHEMA s",
		"DROP SCHEMA s CASCADE",
		// DATABASE
		"CREATE DATABASE d",
		"DROP DATABASE d",
		// UTILITY
		"EXPLAIN SELECT 1",
		"VACUUM t",
		"ANALYZE t",
		// TRANSACTION
		"BEGIN",
		"COMMIT",
		"ROLLBACK",
		// SET
		"SET search_path TO public",
		"SHOW search_path",
		"RESET ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}
