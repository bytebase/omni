package parser

import (
	"strings"
	"testing"
)

// tvcOracleCase is a statement and whether MariaDB 11.8 parses it.
type tvcOracleCase struct {
	sql    string
	accept bool
}

func tvcAccept(sqls ...string) []tvcOracleCase { return tvcCases(true, sqls) }
func tvcReject(sqls ...string) []tvcOracleCase { return tvcCases(false, sqls) }

func tvcCases(accept bool, sqls []string) []tvcOracleCase {
	cases := make([]tvcOracleCase, len(sqls))
	for i, sql := range sqls {
		cases[i] = tvcOracleCase{sql, accept}
	}
	return cases
}

// tvcOracleCorpus is every position where a TABLE statement, a MySQL VALUES
// ROW(...) row, or a MariaDB table value constructor can stand, with
// MariaDB 11.8's verdict.
//
// MariaDB has no TABLE statement and no ROW(...) row (MySQL 8.0.19+ has
// both). Its table value constructor, VALUES (...), (...), is a query
// primary wherever a SELECT is one (MariaDB 11.8.8 sql/sql_yacc.yy:
// simple_table, line 8963), with ORDER BY, LIMIT, and locking after the query
// expression it ends and INTO only after one that stands as a statement.
// VALUES(col) is a function only in ON DUPLICATE KEY UPDATE and PARTITION BY
// LIST. https://mariadb.com/kb/en/table-value-constructors/
var tvcOracleCorpus = concatTVCCases(
	// MySQL's TABLE statement, in every position omni used to take it.
	tvcReject(
		"TABLE t",
		"TABLE t ORDER BY a LIMIT 1",
		"TABLE t INTO @x",
		"(TABLE t)",
		"SELECT * FROM t UNION TABLE u",
		"SELECT * FROM (TABLE t) x",
		"SELECT * FROM t WHERE a IN (TABLE u)",
		"SELECT (TABLE t)",
		"WITH c AS (TABLE t) SELECT * FROM c",
		"INSERT INTO t TABLE u",
		"REPLACE INTO t TABLE u",
		"CREATE TABLE t2 AS TABLE t",
		"CREATE VIEW v AS TABLE t",
		"EXPLAIN TABLE t",
		"CREATE PROCEDURE tvc_p1() TABLE t",
	),
	// MySQL's ROW(...) rows.
	tvcReject(
		"VALUES ROW(1, 2)",
		"VALUES ROW(1, 2), ROW(3, 4)",
		"VALUES (1), ROW(2)",
		"SELECT * FROM (VALUES ROW(1, 2)) x",
		"SELECT * FROM t UNION VALUES ROW(1, 2)",
		"INSERT INTO t VALUES ROW(1, 2)",
		"INSERT INTO t VALUES (1, 2), ROW(3, 4)",
		"INSERT INTO t VALUE ROW(1, 2)",
		"REPLACE INTO t VALUES ROW(1, 2)",
		"EXPLAIN VALUES ROW(1, 2)",
	),
	// ROW(...) is still a row value in an expression.
	tvcAccept(
		"SELECT ROW(1, 2) = ROW(1, 2)",
		"SELECT * FROM t WHERE (a, b) IN (ROW(1, 2))",
	),
	// The table value constructor as a statement, and its rows.
	tvcAccept(
		"VALUES (1, 2)",
		"VALUES (1, 2), (3, 4)",
		"VALUES ()",
		"VALUES (1, 2), (3)",
		"VALUES (DEFAULT, 1)",
		"VALUES ((1), (2))",
		"VALUES ((SELECT 1), (VALUES (2)))",
		"VALUES (@a := 1)",
		"VALUES (a)",
	),
	tvcReject(
		"VALUES",
		"VALUES 1",
		"VALUES (1,)",
		"VALUES (1) (2)",
		"VALUES (1 AS x)",
		"VALUES (*)",
		"VALUES (t.*)",
		"VALUES (1), (2) AS x",
		"VALUE (1)",
	),
	// The clauses after it, and INTO at statement level.
	tvcAccept(
		"VALUES (1, 2) ORDER BY 1",
		"VALUES (1, 2) ORDER BY 1 DESC LIMIT 1",
		"VALUES (1, 2) LIMIT 1 OFFSET 1",
		"VALUES (1, 2) LIMIT 1, 1",
		"VALUES (1, 2) FOR UPDATE",
		"VALUES (1, 2) LOCK IN SHARE MODE",
		"VALUES (1) ORDER BY 1 LIMIT 1 LOCK IN SHARE MODE",
		"VALUES (1, 2) INTO @x, @y",
		"VALUES (1) INTO OUTFILE '/tmp/tvc_oracle_out'",
		"VALUES (1) INTO @x FOR UPDATE",
		"VALUES (1) FOR UPDATE INTO @x",
		"VALUES (1) ORDER BY 1 LIMIT 1 FOR UPDATE INTO @x",
		"(VALUES (1)) INTO @x",
		"(VALUES (1)) FOR UPDATE",
		"VALUES (1) UNION VALUES (2) INTO @x",
		"VALUES (1) UNION VALUES (2) FOR UPDATE",
		"WITH c AS (SELECT 1) VALUES (1) INTO @x",
		"VALUES (1) ORDER BY (VALUES (1) ORDER BY (SELECT 1))",
		strings.Repeat("VALUES (1) ORDER BY (", 30)+"SELECT 1"+strings.Repeat(")", 30),
		"SELECT ("+strings.Repeat("VALUES (1) ORDER BY (", 30)+"SELECT 1"+strings.Repeat(")", 31),
	),
	tvcReject(
		"VALUES (1) ORDER BY 1 ORDER BY 1",
		"VALUES (1) LIMIT 1 LIMIT 1",
		"VALUES (1) LIMIT 1 ORDER BY 1",
		"VALUES (1) INTO @x INTO @y",
		"VALUES (1) INTO @x ORDER BY 1",
		"VALUES (1) INTO @x UNION VALUES (2)",
		"VALUES (1) ORDER BY 1 WITH ROLLUP",
		"VALUES (1) WHERE 1",
		"VALUES (1) GROUP BY 1",
		"VALUES (1) RETURNING a",
		"SELECT * FROM (VALUES (1) INTO @x) t",
		"SELECT * FROM t WHERE a IN (VALUES (1) INTO @x)",
		"(VALUES (1) INTO @x)",
		"CREATE VIEW v AS VALUES (1) INTO @x",
		"CREATE TABLE t2 VALUES (1) INTO @x",
		"INSERT INTO t VALUES (1, 2) INTO @x",
	),
	// Parenthesized and in set operations.
	tvcAccept(
		"(VALUES (1, 2))",
		"((VALUES (1, 2)))",
		"(VALUES (1, 2)) ORDER BY 1",
		"(VALUES (1, 2) ORDER BY 1) LIMIT 1",
		"(VALUES (1, 2)) UNION (VALUES (3, 4))",
		"VALUES (1, 2) UNION VALUES (3, 4)",
		"VALUES (1, 2) UNION ALL VALUES (3, 4)",
		"VALUES (1, 2) UNION SELECT 3, 4",
		"SELECT 1, 2 UNION VALUES (3, 4)",
		"SELECT 1, 2 UNION VALUES (3, 4) ORDER BY 1",
		"SELECT 1, 2 UNION ALL VALUES (3, 4) LIMIT 1",
		"VALUES (1, 2) INTERSECT VALUES (1, 2)",
		"VALUES (1, 2) EXCEPT VALUES (3, 4)",
		"VALUES (1, 2) EXCEPT ALL VALUES (3, 4)",
		"SELECT 1 UNION (VALUES (2))",
		"VALUES (1) UNION VALUES (2) ORDER BY 1 LIMIT 1",
		"VALUES (1) UNION (VALUES (2) ORDER BY 1 LIMIT 1) ORDER BY 1 LIMIT 1",
		"((VALUES (1)) UNION VALUES (2)) ORDER BY 1",
		"VALUES (1) UNION VALUES (2) INTERSECT VALUES (3)",
		"(VALUES (1) LIMIT 1) UNION (VALUES (2) LIMIT 1)",
		"SELECT 1 UNION VALUES (2) UNION SELECT 3",
	),
	tvcReject(
		"VALUES (1) LIMIT 1 UNION VALUES (2)",
		"VALUES (1) ORDER BY 1 UNION VALUES (2)",
		"VALUES (1) FOR UPDATE UNION VALUES (2)",
		"SELECT 1 UNION VALUES (2) LIMIT 1 UNION SELECT 3",
		"SELECT ("+strings.Repeat("VALUES (1) LIMIT 1 UNION VALUES (", 3)+"1"+strings.Repeat(")", 4),
	),
	// Derived tables, subqueries, and CTEs.
	tvcAccept(
		"SELECT * FROM (VALUES (1, 2)) x",
		"SELECT * FROM (VALUES (1, 2)) AS x",
		"SELECT * FROM (VALUES (1, 2)) AS `x`",
		"SELECT * FROM (VALUES (1, 2), (3, 4)) x (c1, c2)",
		"SELECT * FROM (VALUES (1, 2) ORDER BY 1 LIMIT 1) x",
		"SELECT * FROM (VALUES (1) FOR UPDATE) t",
		"SELECT * FROM ((VALUES (1, 2))) x",
		"SELECT * FROM t JOIN (VALUES (1, 2)) x ON t.a = x.`1`",
		"SELECT * FROM t NATURAL JOIN (VALUES (1, 2)) AS x (a, b)",
		"SELECT * FROM (VALUES (1, 2) UNION VALUES (3, 4)) x",
		"SELECT * FROM (WITH c AS (SELECT 1) VALUES (1)) x",
		"DELETE t FROM t JOIN (VALUES (1)) x ON t.a = x.`1`",
		"UPDATE t JOIN (VALUES (1)) x ON t.a = x.`1` SET t.b = 0",
		"SELECT * FROM t WHERE a IN (VALUES (1), (3))",
		"SELECT * FROM t WHERE a NOT IN (VALUES (1), (3))",
		"SELECT * FROM t WHERE (a, b) IN (VALUES (1, 2))",
		"SELECT * FROM t WHERE a IN ((VALUES (1)))",
		"SELECT * FROM t WHERE a IN (VALUES (1) UNION VALUES (2))",
		"SELECT * FROM t WHERE a IN (VALUES (1) ORDER BY 1 LIMIT 1)",
		"SELECT * FROM t WHERE a IN (SELECT 1 UNION VALUES (2))",
		"SELECT * FROM t WHERE a IN (WITH c AS (SELECT 1) VALUES (1))",
		"SELECT * FROM t WHERE a IN (VALUES ())",
		"SELECT * FROM t WHERE a = (VALUES (1))",
		"SELECT * FROM t WHERE a = ANY (VALUES (1))",
		"SELECT * FROM t WHERE a = ALL (VALUES (1))",
		"SELECT a FROM t WHERE a > SOME (VALUES (1), (2))",
		"SELECT * FROM t WHERE (a, b) = (VALUES (1, 2))",
		"SELECT * FROM t WHERE (VALUES (1)) = a",
		"SELECT * FROM t WHERE EXISTS (VALUES (1))",
		"SELECT NOT EXISTS (VALUES (1))",
		"SELECT (VALUES (1))",
		"SELECT (VALUES (1)) + 1",
		"SELECT (VALUES (1), (2) LIMIT 1)",
		"SELECT (WITH c AS (SELECT 1) VALUES (1))",
		"SELECT 1 IN (VALUES (1))",
		"DO (VALUES (1))",
		"SET @x = (VALUES (1))",
		"UPDATE t SET a = (VALUES (1))",
		"DELETE FROM t WHERE a IN (VALUES (1))",
		"WITH c AS (VALUES (1, 2)) SELECT * FROM c",
		"WITH c (x, y) AS (VALUES (1, 2)) SELECT * FROM c",
		"WITH c AS (VALUES (1, 2)) VALUES (3, 4)",
		"WITH c AS (SELECT 1) VALUES (3, 4)",
		"WITH c AS (SELECT 1) (VALUES (3, 4))",
		"WITH RECURSIVE c (n) AS (VALUES (1) UNION ALL SELECT n + 1 FROM c WHERE n < 3) SELECT * FROM c",
	),
	tvcReject(
		"SELECT * FROM LATERAL (VALUES (1)) x",
		"SELECT * FROM (VALUES (1)) x USE INDEX (i)",
		"SELECT * FROM t WHERE a IN (VALUES 1)",
	),
	// INSERT and REPLACE: rows followed by a set operation or the clauses
	// after a query expression are a query; VALUE starts the first rows.
	tvcAccept(
		"INSERT INTO t VALUES (1, 2)",
		"INSERT INTO t VALUES ()",
		"INSERT INTO t VALUES (), ()",
		"INSERT INTO t VALUE (1, 2)",
		"INSERT INTO t VALUES (1, 2) ORDER BY 1",
		"INSERT INTO t VALUES (1, 2) LIMIT 1",
		"INSERT INTO t VALUES (1, 2) FOR UPDATE",
		"INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)",
		"INSERT INTO t VALUES (1, 2), (3, 4) UNION VALUES (5, 6)",
		"INSERT INTO t VALUES (1, 2) UNION ALL VALUES (3, 4) UNION DISTINCT SELECT 5, 6",
		"INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) ORDER BY 1 LIMIT 1",
		"INSERT INTO t VALUES (1, 2) INTERSECT VALUES (1, 2)",
		"INSERT INTO t VALUES (1, 2) EXCEPT VALUES (3, 4)",
		"INSERT INTO t VALUES (1, 2) UNION (VALUES (3, 4))",
		"INSERT INTO t VALUES (1, 2) UNION (SELECT 3, 4)",
		"INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) ON DUPLICATE KEY UPDATE a = 1",
		"INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4) RETURNING a",
		"INSERT INTO t VALUE (1, 2) UNION VALUES (3, 4)",
		"INSERT INTO t (a, b) VALUES (1, 2) UNION VALUES (3, 4)",
		"INSERT INTO t (a, b) VALUE (1, 2) UNION VALUES (3, 4)",
		"INSERT INTO t (a, b) VALUES (1, 2) UNION SELECT 3, 4",
		"INSERT INTO t PARTITION (p0) VALUES (1, 2) UNION VALUES (3, 4)",
		"INSERT IGNORE INTO t VALUES (1, 2) UNION VALUES (3, 4)",
		"INSERT INTO t SELECT 1, 2 UNION VALUES (3, 4)",
		"INSERT INTO t SELECT * FROM (VALUES (1, 2)) x",
		"INSERT INTO t (VALUES (1, 2))",
		"INSERT INTO t (a, b) (VALUES (1, 2))",
		"INSERT INTO t ((VALUES (1, 2)))",
		"INSERT INTO t ((VALUES (1, 2)) UNION (VALUES (3, 4)))",
		"INSERT INTO t (VALUES (1, 2)) UNION VALUES (3, 4)",
		"INSERT INTO t (VALUES (1, 2)) ON DUPLICATE KEY UPDATE a = 1",
		"INSERT INTO t (VALUES (1, 2)) ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"INSERT INTO t (VALUES (1, 2)) RETURNING a",
		"INSERT INTO t (SELECT 1, 2) ON DUPLICATE KEY UPDATE a = 1",
		"INSERT INTO t (WITH c AS (SELECT 1, 2) VALUES (5, 6))",
		"INSERT INTO t WITH c AS (SELECT 1, 2) VALUES (5, 6)",
		"INSERT INTO t WITH c AS (SELECT 1, 2) VALUES (5, 6) UNION VALUES (7, 8)",
		"INSERT INTO t WITH c AS (VALUES (1, 2)) SELECT * FROM c",
		"REPLACE INTO t VALUES (1, 2) UNION VALUES (3, 4)",
		"REPLACE INTO t VALUE (1, 2) ORDER BY 1",
		"REPLACE INTO t (VALUES (1, 2))",
		"REPLACE INTO t (a, b) (VALUES (1, 2))",
		"REPLACE INTO t SELECT * FROM (VALUES (1, 2)) x",
		"CREATE EVENT tvc_ev ON SCHEDULE EVERY 1 DAY DO INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)",
	),
	tvcReject(
		"INSERT INTO t VALUE (1, 2) UNION VALUE (3, 4)",
		"INSERT INTO t VALUES (1, 2) UNION VALUE (3, 4)",
		"INSERT INTO t SELECT 1, 2 UNION VALUE (3, 4)",
		"INSERT INTO t (VALUE (1, 2))",
		"INSERT INTO t SET a = 1 UNION VALUES (3, 4)",
		"INSERT INTO t VALUES (*)",
		"INSERT INTO t VALUES (t.*)",
		"INSERT INTO t VALUES (1, 2) AS new",
		"INSERT INTO t (VALUES (1, 2)) AS x",
	),
	// VALUES(col) is a function in ON DUPLICATE KEY UPDATE, outside the
	// queries nested there, and of a column only; elsewhere VALUES begins a
	// table value constructor. VALUE(col) is the function everywhere.
	tvcAccept(
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a) + VALUES(b)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(t.a)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(test.t.a)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(`a`)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (VALUES(a))",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (SELECT 1 FROM (VALUES (1)) x)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUE(a)",
		"INSERT INTO t SELECT 1, 2 ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"INSERT INTO t SET a = 1 ON DUPLICATE KEY UPDATE a = VALUES(a)",
		"SELECT VALUE(a) FROM t",
		"UPDATE t SET a = VALUE(a)",
		"CREATE TABLE t3 (a INT) PARTITION BY LIST (VALUES(a)) (PARTITION p0 VALUES IN (1))",
		"CREATE TABLE t3 (a INT) PARTITION BY LIST (a) (PARTITION p0 VALUES IN (1))",
		"CREATE TABLE t3 (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (10))",
		"CREATE TABLE t3 (a INT, CHECK ((VALUES (1)) > 0))",
		"CREATE TABLE t3 (a INT CHECK (a IN (VALUES (1))))",
	),
	tvcReject(
		"SELECT VALUES(a) FROM t",
		"SELECT * FROM t WHERE VALUES(a) = 1",
		"UPDATE t SET a = VALUES(a)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (SELECT VALUES(a))",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = (VALUES (1))",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = a IN (VALUES (1))",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(1)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a + 1)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(@x)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(t.*)",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES()",
		"INSERT INTO t VALUES (1, 2) ON DUPLICATE KEY UPDATE a = VALUES(a) UNION VALUES (1)",
		"CREATE TABLE t3 (a INT) PARTITION BY HASH (VALUES(a))",
		"CREATE TABLE t3 (a INT DEFAULT VALUES(a))",
		"CREATE TABLE t3 (a INT, CHECK (VALUES(a) > 0))",
	),
	// CREATE TABLE ... [AS] query. VALUE reads as VALUES until a '(' or a
	// PARTITION clause follows the table name.
	tvcAccept(
		"CREATE TABLE t2 AS VALUES (1, 2)",
		"CREATE TABLE t2 VALUES (1, 2)",
		"CREATE TABLE t2 (VALUES (1, 2))",
		"CREATE TABLE t2 AS (VALUES (1, 2))",
		"CREATE TABLE t2 ((VALUES (1)))",
		"CREATE TABLE t2 (WITH c AS (SELECT 1) VALUES (1))",
		"CREATE TABLE t2 (VALUES (1)) UNION VALUES (2)",
		"CREATE TABLE t2 AS SELECT * FROM (VALUES (1, 2)) x",
		"CREATE TABLE t2 AS VALUES (1, 2) UNION VALUES (3, 4)",
		"CREATE TABLE t2 VALUES (1) UNION VALUES (2)",
		"CREATE TABLE t2 VALUES (1) ORDER BY 1 LIMIT 1",
		"CREATE TABLE t2 VALUES (1) FOR UPDATE",
		"CREATE TABLE t2 (a INT, b INT) AS VALUES (1, 2)",
		"CREATE TABLE t2 (a INT, b INT) VALUES (1, 2)",
		"CREATE TABLE t2 (a INT) (VALUES (1))",
		"CREATE TABLE t2 (a INT) (SELECT 1)",
		"CREATE TABLE t2 (a INT) IGNORE VALUES (1)",
		"CREATE TABLE t2 (a INT) ENGINE=InnoDB VALUES (1)",
		"CREATE TABLE t2 (a INT) PARTITION BY HASH(a) VALUES (1)",
		"CREATE TABLE t2 AS WITH c AS (SELECT 1) VALUES (1, 2)",
		"CREATE TABLE t2 ENGINE=InnoDB VALUES (1)",
		"CREATE TABLE t2 ENGINE=InnoDB AS VALUES (1)",
		"CREATE TABLE t2 ENGINE=InnoDB VALUE (1)",
		"CREATE TABLE t2 IGNORE VALUES (1)",
		"CREATE TABLE t2 REPLACE VALUES (1)",
		"CREATE TABLE t2 IGNORE AS VALUES (1)",
		"CREATE TABLE t2 VALUE (1)",
		"CREATE TABLE t2 AS VALUE (1)",
		"CREATE TABLE t2 VALUE (1) UNION VALUES (2)",
		"CREATE OR REPLACE TABLE t2 AS VALUES (1)",
		"CREATE TEMPORARY TABLE t2 VALUES (1)",
	),
	tvcReject(
		"CREATE TABLE t2 (a INT) VALUE (1)",
		"CREATE TABLE t2 (a INT) AS VALUE (1)",
		"CREATE TABLE t2 (VALUE (1))",
		"CREATE TABLE t2 VALUE (1) UNION VALUE (2)",
	),
	// CREATE VIEW, EXPLAIN, and routines.
	tvcAccept(
		"CREATE VIEW v AS VALUES (1, 2)",
		"CREATE VIEW v AS (VALUES (1, 2))",
		"CREATE VIEW v AS ((VALUES (1)))",
		"CREATE VIEW v AS VALUES (1, 2) UNION VALUES (3, 4)",
		"CREATE VIEW v AS (VALUES (1)) UNION (VALUES (2)) ORDER BY 1 LIMIT 1",
		"CREATE VIEW v (x, y) AS VALUES (1, 2)",
		"CREATE VIEW v AS VALUES (1, 2) WITH CHECK OPTION",
		"CREATE VIEW v AS VALUES (1) FOR UPDATE",
		"CREATE VIEW v AS WITH c AS (VALUES (1, 2)) SELECT * FROM c",
		"CREATE OR REPLACE VIEW v AS VALUES (1, 2) ORDER BY 1",
		"CREATE ALGORITHM=MERGE VIEW v AS VALUES (1)",
		"ALTER VIEW v AS VALUES (1, 2)",
		"EXPLAIN VALUES (1, 2)",
		"EXPLAIN (VALUES (1, 2))",
		"EXPLAIN VALUES (1, 2) UNION VALUES (3, 4)",
		"EXPLAIN EXTENDED (VALUES (1)) UNION VALUES (2)",
		"EXPLAIN FORMAT=JSON VALUES (1, 2)",
		"EXPLAIN EXTENDED VALUES (1, 2)",
		"EXPLAIN PARTITIONS VALUES (1)",
		"EXPLAIN WITH c AS (SELECT 1) VALUES (1)",
		"EXPLAIN VALUES (1) INTO @x",
		"EXPLAIN VALUES (1) ORDER BY 1 LIMIT 1",
		"EXPLAIN INSERT INTO t VALUES (1, 2) UNION VALUES (3, 4)",
		"CREATE PROCEDURE tvc_p2() BEGIN VALUES (1); END",
		"CREATE PROCEDURE tvc_p3() VALUES (1)",
		"CREATE PROCEDURE tvc_p4() BEGIN DECLARE c CURSOR FOR VALUES (1); END",
		"CREATE PROCEDURE tvc_p5() BEGIN DECLARE c CURSOR FOR (VALUES (1)) UNION VALUES (2); END",
		"CREATE FUNCTION tvc_f1() RETURNS INT RETURN (VALUES (1))",
		"CREATE TRIGGER tvc_tg BEFORE INSERT ON t FOR EACH ROW SET NEW.a = (VALUES (1))",
		"PREPARE tvc_s FROM 'VALUES (1, 2)'",
	),
	tvcReject(
		"CREATE VIEW v AS VALUE (1)",
		"CREATE PROCEDURE tvc_p6() BEGIN DECLARE c CURSOR FOR TABLE t; END",
	),
)

func concatTVCCases(groups ...[]tvcOracleCase) []tvcOracleCase {
	var all []tvcOracleCase
	for _, g := range groups {
		all = append(all, g...)
	}
	return all
}

// TestTableValueConstructorOracle runs every case on MariaDB 11.8 and through
// Parse, and requires both to give the pinned verdict. A statement Parse
// accepts must be one statement: the parser does not require ';' between
// statements, so a text it read as two would pass for an accept.
func TestTableValueConstructorOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, ddl := range []string{
		"CREATE TABLE t (a INT, b INT)",
		"CREATE TABLE u (a INT, b INT)",
		"INSERT INTO t VALUES (1, 2), (3, 4)",
	} {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}
	var cases []tvcOracleCase
	cases = append(cases, tvcOracleCorpus...)
	for _, tc := range sourceClauseCases {
		cases = append(cases, tvcOracleCase{tc.sql, true})
	}
	t.Logf("%d statements", len(cases))
	seen := map[string]bool{}
	for _, tc := range cases {
		if seen[tc.sql] {
			t.Errorf("duplicate case %q", tc.sql)
		}
		seen[tc.sql] = true

		v, code, msg := o.classify(tc.sql)
		if v == vErrored {
			t.Errorf("ERRORED (infra, not a parse verdict) for %q: %s", tc.sql, msg)
			continue
		}
		if mdbAccepts := v == vAccepted; mdbAccepts != tc.accept {
			t.Errorf("MariaDB accept=%v (code %d: %s), pinned %v: %q", mdbAccepts, code, msg, tc.accept, tc.sql)
		}

		list, err := Parse(tc.sql)
		switch {
		case tc.accept && err != nil:
			t.Errorf("Parse(%q): %v", tc.sql, err)
		case tc.accept && len(list.Items) != 1:
			t.Errorf("Parse(%q): %d statements, want 1", tc.sql, len(list.Items))
		case !tc.accept && err == nil:
			t.Errorf("Parse(%q) accepted what MariaDB rejects", tc.sql)
		}
	}
}
