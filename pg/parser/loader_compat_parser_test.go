package parser

import "testing"

func TestLoaderCompatParserRejectsInvalidTails(t *testing.T) {
	tests := []string{
		`DROP FUNCTION ();`,
		`ALTER FUNCTION f();`,
		`COMMENT ON FUNCTION f(integer IS 'comment';`,
		`GRANT EXECUTE ON FUNCTION f(integer TO role_name;`,
		`CREATE TABLE t (x integer, FOREIGN KEY (x) REFERENCES p MATCH);`,
		`ALTER INDEX i ATTACH PARTITION;`,
		`WITH c AS SELECT 1 SELECT * FROM c;`,
		`SELECT * FROM t, LATERAL ();`,
		`SELECT * FROM t, LATERAL relation_name;`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", sql)
			}
		})
	}
}

func TestLoaderCompatParserAcceptsKeywordIdentifiers(t *testing.T) {
	tests := []string{
		`CREATE TABLE "select" (id integer);`,
		`CREATE TABLE t ("select" integer);`,
		`CREATE FUNCTION "select"(a integer) RETURNS integer LANGUAGE sql AS 'SELECT $1';`,
		`CREATE TYPE "select" AS ENUM ('a');`,
		`CREATE SCHEMA "select";`,
		`COMMENT ON TABLE "select" IS 'comment';`,
		`GRANT SELECT ON "select" TO PUBLIC;`,
		`DROP TABLE "select";`,
		`ALTER TABLE "select" RENAME TO select_new;`,
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q): %v", sql, err)
			}
		})
	}
}
