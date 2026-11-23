package parser

import (
	"testing"

	bigqueryparser "github.com/bytebase/omni/bigquery/parser"
	googlesqlparser "github.com/bytebase/omni/googlesql/parser"
	pgparser "github.com/bytebase/omni/pg/parser"
)

func TestSpannerAcceptedSyntax(t *testing.T) {
	accepted := []string{
		// Queries & shared GoogleSQL + Spanner-only query forms
		"SELECT 1",
		"SELECT NULL",
		"SELECT NULL FOR UPDATE",
		"SELECT { }",
		"SELECT NEW a { }",
		"SELECT STRUCT { }",
		"SELECT QUALIFY",
		"SELECT NULL QUALIFY",
		"TABLE a() QUALIFY",
		"SELECT * FROM (SELECT 1) AS a",
		"(SELECT 1) UNION ALL (SELECT 2)",
		"SELECT ARRAY { RETURN NULL }",
		"SELECT EXISTS { GRAPH a - }",
		"|> CREATE TABLE a",

		// Spanner CREATE TABLE forms
		"CREATE TABLE T (x INT64) PRIMARY KEY (x)",
		"CREATE TABLE T (x INT64, y STRING(MAX) NOT NULL) PRIMARY KEY (x DESC, y)",
		"CREATE TABLE IF NOT EXISTS T (x INT64) PRIMARY KEY (x)",
		"CREATE TABLE T () PRIMARY KEY ()",
		"CREATE TABLE T (PRIMARY KEY ())",
		"CREATE TABLE T (id INT64 PRIMARY KEY, name STRING(MAX))",
		"CREATE TABLE Albums (s INT64, a INT64) PRIMARY KEY (s, a), INTERLEAVE IN PARENT Singers ON DELETE CASCADE",
		"CREATE TABLE Child (id INT64) PRIMARY KEY (id), INTERLEAVE IN PARENT Parent",
		"CREATE TABLE T (x INT64, c TIMESTAMP OPTIONS (allow_commit_timestamp = true)) PRIMARY KEY (x), ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 30 DAY))",
		"CREATE TABLE T (x INT64 NOT NULL, e ARRAY<FLOAT32>(vector_length=>128)) PRIMARY KEY (x)",

		// Spanner CREATE/ALTER/DROP INDEX forms
		"CREATE INDEX idx ON T (a, b)",
		"CREATE UNIQUE INDEX idx ON T (a)",
		"CREATE NULL_FILTERED INDEX idx ON T (a DESC)",
		"CREATE INDEX idx ON T (a) STORING (b, c)",
		"CREATE INDEX idx ON T (a) STORING (b), INTERLEAVE IN P",
		"CREATE INDEX idx ON T (a) WHERE a IS NOT NULL",
		"CREATE INDEX idx ON T ()",
		"CREATE INDEX idx ON T (COLUMNS)",
		"CREATE INDEX idx ON T ((LOWER(name)))",
		"CREATE INDEX idx ON T ((a + b) DESC, c)",
		"CREATE SEARCH INDEX sidx ON T (tok)",
		"CREATE SEARCH INDEX AlbumsIndex ON Albums(DescriptionTokens) PARTITION BY SingerId",
		"CREATE VECTOR INDEX vidx ON T (emb) OPTIONS (distance_type = 'COSINE')",
		"ALTER INDEX idx ADD STORED COLUMN c",
		"ALTER INDEX idx DROP STORED COLUMN c",
		"ALTER SEARCH INDEX sidx ADD STORED COLUMN c",
		"DROP INDEX idx",
		"DROP SEARCH INDEX sidx",
		"DROP VECTOR INDEX vidx",

		// Spanner VIEW & FUNCTION
		"CREATE VIEW v SQL SECURITY INVOKER AS SELECT 1",
		"CREATE VIEW v AS SELECT 1",
		"CREATE OR REPLACE VIEW v SQL SECURITY INVOKER AS SELECT a FROM T",
		"CREATE FUNCTION f() RETURNS INT64 AS (1)",
		"CREATE FUNCTION f(x INT64) RETURNS INT64 AS (x + 1)",
		"CREATE OR REPLACE FUNCTION f(x INT64) RETURNS INT64 SQL SECURITY INVOKER AS (x)",
		"CREATE FUNCTION f(x INT64) RETURNS INT64 NOT DETERMINISTIC LANGUAGE REMOTE OPTIONS (endpoint = 'https://example.com')",
		"CREATE OR REPLACE FUNCTION f(x INT64) RETURNS INT64 LANGUAGE REMOTE OPTIONS (endpoint = 'https://example.com')",
		"DROP VIEW v",
		"DROP FUNCTION f",

		// Spanner-specific DDL objects
		"CREATE CHANGE STREAM MyStream FOR ALL",
		"CREATE CHANGE STREAM cs FOR ALL OPTIONS (retention_period = '7d')",
		"ALTER CHANGE STREAM MyStream SET FOR ALL",
		"ALTER CHANGE STREAM MyStream SET OPTIONS (retention_period = '2d')",
		"DROP CHANGE STREAM MyStream",
		"CREATE SEQUENCE Seq1",
		"CREATE SEQUENCE Seq1 OPTIONS (sequence_kind = 'bit_reversed_positive')",
		"ALTER SEQUENCE Seq1 SET OPTIONS (start_with_counter = 9000)",
		"DROP SEQUENCE Seq1",
		"CREATE LOCALITY GROUP lg_hot OPTIONS (storage = 'ssd')",
		"ALTER LOCALITY GROUP lg_hot SET OPTIONS (storage = 'hdd')",
		"ALTER LOCALITY GROUP lg_bare",
		"DROP LOCALITY GROUP lg_hot",
		"CREATE PROTO BUNDLE (`a.b.C`, `a.b.D`)",
		"ALTER PROTO BUNDLE INSERT (`a.b.C`) DELETE (`a.b.D`)",
		"CREATE ROLE analyst_role",
		"DROP ROLE analyst_role",
		"CREATE PROPERTY GRAPH g NODE TABLES (T)",
		"DROP PROPERTY GRAPH g",
		"CREATE MODEL m",
		"CREATE OR REPLACE MODEL m REMOTE",
		"DROP MODEL m",
		"CREATE SCHEMA sch",
		"DROP SCHEMA sch",
		"CREATE DATABASE db",
		"ALTER DATABASE db SET OPTIONS (version_retention_period = '7d')",

		// Spanner GRANT / REVOKE, RENAME, ANALYZE, transactions
		"GRANT SELECT ON TABLE T TO ROLE r",
		"GRANT SELECT, INSERT ON TABLE t1, t2 TO ROLE r1, r2",
		"GRANT ROLE r1 TO ROLE r2",
		"REVOKE SELECT ON TABLE T FROM ROLE r",
		"REVOKE ROLE r1 FROM ROLE r2",
		"RENAME TABLE a TO b",
		"ANALYZE",
		"BEGIN",
		"SET a = NULL",
		"DESCRIBE a",
		"LOAD DATA INTO a FROM FILES()",
		"|> EXPORT DATA",
		"|> AGGREGATE GROUP BY NULL",
		"|> INSERT OR REPLACE a",
		"|> INSERT OR UPDATE a",
	}

	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("spanner/parser.Parse(%q) unexpected error: %v", sql, err)
		}
		if _, err := ParseGoogleSQL(sql); err != nil {
			t.Errorf("spanner/parser.ParseGoogleSQL(%q) unexpected error: %v", sql, err)
		}
		if res := ParseBestEffort(sql); len(res.Errors) > 0 {
			t.Errorf("spanner/parser.ParseBestEffort(%q) unexpected errors: %v", sql, res.Errors)
		}
	}
}

func TestSpannerRejectsBigQueryOnlySyntax(t *testing.T) {
	rejected := []string{
		// BigQuery procedural scripting
		"DECLARE x INT64",
		"DECLARE a DEFAULT NULL",
		"IF NULL THEN END IF",
		"CASE WHEN NULL THEN END CASE",
		"WHILE NULL DO END WHILE",
		"LOOP END LOOP",
		"REPEAT UNTIL NULL END REPEAT",
		"FOR a IN (SELECT NULL) DO END FOR",
		"BREAK",
		"LEAVE",
		"CONTINUE",
		"ITERATE",
		"RETURN",
		"RAISE",
		"RAISE USING MESSAGE = 'err'",
		"lbl: LOOP END LOOP",
		"lbl: WHILE NULL DO END WHILE",
		"BEGIN END",
		"a : BEGIN END",
		"BEGIN EXCEPTION WHEN ERROR THEN END",

		// Spanner-rejected pipe / DESCRIBE / LOAD DATA / bare INSERT REPLACE|UPDATE forms
		"DESCRIBE a |> AGGREGATE",
		"LOAD DATA INTO TEMP TABLE a FROM FILES()",
		"LOAD DATA INTO TEMPORARY TABLE a FROM FILES()",
		"LOAD DATA INTO a() FROM FILES()",
		"LOAD DATA INTO a PARTITIONS(NULL) FROM FILES()",
		"LOAD DATA INTO a OVERWRITE PARTITIONS(NULL) FROM FILES()",
		"|> EXPORT",
		"|> EXPORT AS",
		"|> MATCH_RECOGNIZE(ORDER BY NULL MEASURES NULL AS a PATTERN() DEFINE a AS NULL)",
		"|> AGGREGATE GROUP AND ORDER BY NULL",
		"|> INSERT REPLACE a",
		"|> INSERT UPDATE a",

		// BigQuery-only query clauses
		"SELECT NULL QUALIFY NULL",
		"SELECT NULL WHERE NULL QUALIFY NULL",
		"SELECT NULL GROUP BY ALL QUALIFY NULL",
		"SELECT NULL HAVING NULL QUALIFY NULL",
		"SELECT NULL IS DISTINCT FROM NULL",
		"SELECT NULL IS NOT DISTINCT FROM NULL",
		"(SELECT NULL) AS a",
		"FROM a",

		// BigQuery-only CREATE TABLE forms
		"CREATE TABLE a",
		"CREATE TABLE a PRIMARY KEY (a)",
		"CREATE TABLE a (a INT64)",
		"CREATE TABLE a (CHECK (NULL))",
		"CREATE TABLE a AS SELECT 1",
		"CREATE TABLE a LIKE b",
		"CREATE TABLE a CLONE b",
		"CREATE TABLE a COPY b",
		"CREATE TABLE a PARTITION BY d",
		"CREATE TABLE a CLUSTER BY c",
		"CREATE TABLE a OPTIONS ()",
		"CREATE TABLE a (x INT64) PRIMARY KEY (x) PARTITION BY x",
		"CREATE TABLE a (x INT64) PRIMARY KEY (x) OPTIONS (description = 'x')",

		// BigQuery-only CREATE/ALTER/DROP INDEX forms
		"CREATE INDEX a ON a (NULL)",
		"CREATE UNIQUE INDEX a ON a (NULL)",
		"CREATE NULL_FILTERED INDEX a ON a (NULL)",
		"CREATE SEARCH INDEX a ON a (NULL)",
		"CREATE INDEX a ON a (ALL COLUMNS)",
		"CREATE SEARCH INDEX idx ON t (ALL COLUMNS)",
		"CREATE INDEX a ON a (x) STORING (NULL)",
		"CREATE INDEX a ON a (x) PARTITION BY x",
		"CREATE INDEX a ON a (x) OPTIONS ()",
		"CREATE SEARCH INDEX sidx ON T ((LOWER(name)))",
		"CREATE VECTOR INDEX vidx ON T (emb) PARTITION BY p",
		"DROP SEARCH INDEX a ON t",
		"DROP VECTOR INDEX a ON t",
		"ALTER VECTOR INDEX idx ON t REBUILD",

		// BigQuery-only VIEW forms
		"CREATE MATERIALIZED VIEW a AS SELECT 1",
		"CREATE APPROX VIEW a AS SELECT 1",
		"CREATE RECURSIVE VIEW a AS SELECT 1",
		"CREATE VIEW a (c) AS SELECT 1",
		"CREATE VIEW a OPTIONS (description = 'v') AS SELECT 1",

		// BigQuery-only FUNCTION / PROCEDURE forms
		"CREATE FUNCTION a()",
		"CREATE FUNCTION a() RETURNS INT64",
		"CREATE FUNCTION a() AS (1)",
		"CREATE FUNCTION a() RETURNS INT64 AS 'return 1'",
		"CREATE FUNCTION a() RETURNS TABLE<x INT64> AS (1)",
		"CREATE FUNCTION a(INT64) RETURNS INT64 AS (1)",
		"CREATE FUNCTION a(x INT64 AS y) RETURNS INT64 AS (x)",
		"CREATE FUNCTION a(x INT64 NOT AGGREGATE) RETURNS INT64 AS (x)",
		"CREATE FUNCTION a() RETURNS INT64 DETERMINISTIC AS (1)",
		"CREATE FUNCTION a() RETURNS INT64 NOT DETERMINISTIC AS (1)",
		"CREATE FUNCTION a() RETURNS INT64 LANGUAGE js AS '1'",
		"CREATE FUNCTION a() RETURNS INT64 REMOTE AS (1)",
		"CREATE FUNCTION a() RETURNS INT64 LANGUAGE REMOTE",
		"CREATE FUNCTION a() RETURNS INT64 LANGUAGE REMOTE OPTIONS (endpoint = 'https://example.com') AS (1)",
		"CREATE AGGREGATE FUNCTION a(x INT64) RETURNS INT64 AS (SUM(x))",
		"CREATE TEMP FUNCTION a() RETURNS INT64 AS (1)",
		"CREATE TABLE FUNCTION a() AS SELECT 1",
		"CREATE PROCEDURE a() BEGIN END",

		// BigQuery-only ROW ACCESS POLICY / SNAPSHOT / generic entity / DROP forms
		"CREATE ROW POLICY ON a USING (NULL)",
		"CREATE ROW ACCESS POLICY p ON a GRANT TO ('u') FILTER USING (TRUE)",
		"CREATE SNAPSHOT TABLE a CLONE b",
		"CREATE PROJECT a",
		"CREATE CAPACITY a OPTIONS (slot_count = 100)",
		"DROP ALL ROW POLICIES ON a",
		"DROP ROW ACCESS POLICY p ON a",
		"DROP MATERIALIZED VIEW a",
		"DROP APPROX VIEW a",
		"DROP DATABASE a",
		"DROP EXTERNAL SCHEMA a",
		"DROP PROCEDURE a",
		"DROP TABLE FUNCTION a",
		"DROP SNAPSHOT TABLE a",
		"DROP PROJECT a",
		"ALTER MATERIALIZED VIEW a SET OPTIONS (enable_refresh = true)",

		// BigQuery-only GRANT / REVOKE / ANALYZE / TRUNCATE / RENAME / empty DDL OPTIONS()
		"GRANT ALL ON a TO 'a'",
		"GRANT SELECT ON a TO 'a'",
		"GRANT SELECT ON TABLE a TO 'a'",
		"REVOKE ALL ON a FROM 'a'",
		"REVOKE SELECT ON TABLE a FROM 'a'",
		"ANALYZE a",
		"ANALYZE OPTIONS ()",
		"TRUNCATE TABLE a",
		"RENAME a b TO c",
		"CREATE SCHEMA a OPTIONS ()",
		"CREATE SEQUENCE a OPTIONS ()",
		"CREATE LOCALITY GROUP a OPTIONS ()",
		"ALTER DATABASE db SET OPTIONS ()",
		"ALTER TABLE t ALTER COLUMN c SET OPTIONS ()",

		// Spanner's CREATE VIEW takes OR REPLACE but no IF NOT EXISTS.
		"CREATE VIEW IF NOT EXISTS a SQL SECURITY INVOKER AS SELECT 1",
	}

	for _, sql := range rejected {
		if _, err := Parse(sql); err == nil {
			t.Errorf("spanner/parser.Parse(%q) expected error, got nil", sql)
		}
	}
}

func TestDialectsRemainIndependent(t *testing.T) {
	// BigQuery-only forms must still parse in googlesql/parser (DialectDefault)
	// and bigquery/parser (DialectBigQuery) while rejecting in spanner/parser.
	bqOnly := []string{
		"DECLARE x INT64",
		"SELECT NULL QUALIFY NULL",
		"SELECT NULL IS DISTINCT FROM NULL",
		"CREATE TABLE a",
		"CREATE PROCEDURE a() BEGIN END",
		"GRANT ALL ON a TO 'a'",
		"TRUNCATE TABLE a",
	}
	for _, sql := range bqOnly {
		if _, err := Parse(sql); err == nil {
			t.Errorf("spanner/parser.Parse(%q): want error, got nil", sql)
		}
		if _, err := googlesqlparser.Parse(sql); err != nil {
			t.Errorf("googlesql/parser.Parse(%q): unexpected error: %v", sql, err)
		}
		if _, err := bigqueryparser.Parse(sql); err != nil {
			t.Errorf("bigquery/parser.Parse(%q): unexpected error: %v", sql, err)
		}
	}

	// Spanner-only forms must parse in spanner/parser and googlesql/parser,
	// while rejecting in bigquery/parser.
	spannerOnly := []string{
		"SELECT { }",
		"SELECT NULL FOR UPDATE",
		"CREATE LOCALITY GROUP lg",
		"CREATE NULL_FILTERED INDEX idx ON t (a)",
		"CREATE TABLE t (a INT64) PRIMARY KEY (a), INTERLEAVE IN PARENT p",
	}
	for _, sql := range spannerOnly {
		if _, err := Parse(sql); err != nil {
			t.Errorf("spanner/parser.Parse(%q): unexpected error: %v", sql, err)
		}
		if _, err := googlesqlparser.Parse(sql); err != nil {
			t.Errorf("googlesql/parser.Parse(%q): unexpected error: %v", sql, err)
		}
		if _, err := bigqueryparser.Parse(sql); err == nil {
			t.Errorf("bigquery/parser.Parse(%q): want error, got nil", sql)
		}
	}
}

func TestParsePostgres(t *testing.T) {
	accepted := []string{
		"SELECT 1",
		"SELECT 'a'::text, ARRAY[1, 2, 3]",
		"CREATE TABLE singers (id bigint PRIMARY KEY, name varchar(1024) NOT NULL)",
		"CREATE TABLE albums (id bigint, album_id bigint, PRIMARY KEY (id, album_id)) INTERLEAVE IN PARENT singers ON DELETE CASCADE",
		"CREATE TABLE events (id bigint PRIMARY KEY, ts timestamptz) TTL INTERVAL '30 days' ON ts",
		"CREATE CHANGE STREAM cs FOR ALL WITH (retention_period = '7d')",
		"CREATE VIEW singer_names SQL SECURITY INVOKER AS SELECT id, name FROM singers",
		"CREATE SEQUENCE seq BIT_REVERSED_POSITIVE SKIP RANGE 1000 5000 START COUNTER WITH 100",
		"CREATE INDEX idx ON albums (album_id) INCLUDE (id) INTERLEAVE IN singers",
		"CREATE SEARCH INDEX sidx ON singers (name_tokens) PARTITION BY id ORDER BY id DESC",
		"CREATE LOCALITY GROUP IF NOT EXISTS lg",
		"ALTER DATABASE db SET spanner.default_sequence_kind = 'bit_reversed_positive'",
		"INSERT INTO singers (id, name) VALUES (1, 'Ada') RETURNING id",
		"ALTER TABLE singers ADD COLUMN bio text",
		"DROP TABLE IF EXISTS singers",
	}
	for _, sql := range accepted {
		list, err := ParsePostgres(sql)
		if err != nil {
			t.Errorf("ParsePostgres(%q) unexpected error: %v", sql, err)
		}
		if list == nil || len(list.Items) != 1 {
			t.Errorf("ParsePostgres(%q) got %v, want 1 statement", sql, list)
		}
	}

	rejected := []string{
		"",
		";",
		"SELECT FROM",
		"CREATE TABLE",
		"INSERT INTO",
		"SELECT (1 + )",
		"SELECT XMLELEMENT(NAME a)",
		"SELECT JSON_QUERY(a, a PASSING a AS a)",
		"SELECT a AT LOCAL",
		"CREATE VIEW v SQL SECURITY UNKNOWN AS SELECT 1",
		"CREATE SEQUENCE s SKIP RANGE 1000",
		"ALTER TABLE a SET zzz",
		"ALTER TABLE a ALTER COLUMN c SET",
		"ALTER DATABASE a TARGET",
		"ALTER DATABASE a JSON 1.5",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE NO",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE SET",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE DEFAULT",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE SETOF DEFAULT",
		"CREATE DOMAIN a a REFERENCES a ON DELETE ACTION",
		"CREATE TABLE t (a int PRIMARY KEY, FOREIGN KEY (a) REFERENCES p ON UPDATE ACTION)",
		"CREATE TABLE t (a int PRIMARY KEY, b int ON UPDATE CASCADE)",
		"CREATE TABLE t (a int PRIMARY KEY, b int ON UPDATE a AND b)",
	}
	for _, sql := range rejected {
		if _, err := ParsePostgres(sql); err == nil {
			t.Errorf("ParsePostgres(%q) expected error, got nil", sql)
		}
	}

	// Spanner-PG 1.5.54 grammar forms accepted by ParsePostgres (DialectSpannerPG)
	// and strictly rejected by standard PostgreSQL 17 (pg/parser.Parse).
	spannerPGOnly := []string{
		"ALTER TABLE a SET",
		"ALTER TABLE a SET, ADD COLUMN b int",
		"ALTER DATABASE a TARGET 1.5",
		"ALTER DATABASE a SOURCE = 1.5",
		"CREATE DATABASE a TARGET 1.5",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE ACTION",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE RETURN",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE CASCADED",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE NULL",
		"CREATE DOMAIN a a REFERENCES a ON UPDATE SETOF NULL",
		"CREATE TABLE t (a bigint PRIMARY KEY, ts timestamptz DEFAULT spanner.pending_commit_timestamp() ON UPDATE spanner.pending_commit_timestamp())",
	}
	for _, sql := range spannerPGOnly {
		if _, err := ParsePostgres(sql); err != nil {
			t.Errorf("ParsePostgres(%q) unexpected error: %v", sql, err)
		}
		if _, err := pgparser.Parse(sql); err == nil {
			t.Errorf("pg/parser.Parse(%q) unexpectedly succeeded in DialectDefault; want error", sql)
		}
	}
}
