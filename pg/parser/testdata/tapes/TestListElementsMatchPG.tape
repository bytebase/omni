-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestListElementsMatchPG default#1
> exec
| SELECT 1/0;
| SELECT a, b, FROM t
< error
| 42601 syntax error at or near "FROM"
| position: 26
> exec
| SELECT 1/0;
| SELECT a, FROM t
< error
| 42601 syntax error at or near "FROM"
| position: 23
> exec
| SELECT 1/0;
| SELECT a, b,
< error
| 42601 syntax error at end of input
| position: 25
> exec
| SELECT 1/0;
| SELECT a, b, WHERE true
< error
| 42601 syntax error at or near "WHERE"
| position: 26
> exec
| SELECT 1/0;
| INSERT INTO t VALUES (1) RETURNING a,
< error
| 42601 syntax error at end of input
| position: 50
> exec
| SELECT 1/0;
| UPDATE t SET a = 1 RETURNING a,
< error
| 42601 syntax error at end of input
| position: 44
> exec
| SELECT 1/0;
| DELETE FROM t RETURNING a,
< error
| 42601 syntax error at end of input
| position: 39
> exec
| SELECT 1/0;
| SELECT a, b FROM t
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT FROM t
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| INSERT INTO t VALUES (1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 40
> exec
| SELECT 1/0;
| VALUES (1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 26
> exec
| SELECT 1/0;
| VALUES ()
< error
| 42601 syntax error at or near ")"
| position: 21
> exec
| SELECT 1/0;
| INSERT INTO t VALUES ()
< error
| 42601 syntax error at or near ")"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1 WHERE 1 IN (1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 39
> exec
| SELECT 1/0;
| SELECT 1 WHERE (1, 2) IN ((1, 2), )
< error
| 42601 syntax error at or near ")"
| position: 47
> exec
| SELECT 1/0;
| SELECT ARRAY[1, 2,]
< error
| 42601 syntax error at or near "]"
| position: 31
> exec
| SELECT 1/0;
| SELECT ARRAY[[1], ]
< error
| 42601 syntax error at or near "]"
| position: 31
> exec
| SELECT 1/0;
| SELECT ROW(1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 29
> exec
| SELECT 1/0;
| SELECT (1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 26
> exec
| SELECT 1/0;
| SELECT ()
< error
| 42601 syntax error at or near ")"
| position: 21
> exec
| SELECT 1/0;
| SELECT COALESCE(1, )
< error
| 42601 syntax error at or near ")"
| position: 32
> exec
| SELECT 1/0;
| SELECT GREATEST(1, )
< error
| 42601 syntax error at or near ")"
| position: 32
> exec
| SELECT 1/0;
| SELECT xmlconcat(1, )
< error
| 42601 syntax error at or near ")"
| position: 33
> exec
| SELECT 1/0;
| SELECT xmlelement(NAME x, )
< error
| 42601 syntax error at or near ")"
| position: 39
> exec
| SELECT 1/0;
| SELECT DISTINCT ON (a, ) a FROM t
< error
| 42601 syntax error at or near ")"
| position: 36
> exec
| SELECT 1/0;
| SELECT DISTINCT ON () a FROM t
< error
| 42601 syntax error at or near ")"
| position: 33
> exec
| SELECT 1/0;
| SELECT 1 FROM t WINDOW w AS (PARTITION BY a, )
< error
| 42601 syntax error at or near ")"
| position: 58
> exec
| SELECT 1/0;
| SELECT 1 FROM t GROUP BY ROLLUP (a, )
< error
| 42601 syntax error at or near ")"
| position: 49
> exec
| SELECT 1/0;
| SELECT * FROM t TABLESAMPLE system (1, )
< error
| 42601 syntax error at or near ")"
| position: 52
> exec
| SELECT 1/0;
| EXECUTE p (1, )
< error
| 42601 syntax error at or near ")"
| position: 27
> exec
| SELECT 1/0;
| MERGE INTO t USING u ON true WHEN NOT MATCHED THEN INSERT VALUES (1, )
< error
| 42601 syntax error at or near ")"
| position: 82
> exec
| SELECT 1/0;
| CREATE TABLE t (a numeric(10, ))
< error
| 42601 syntax error at or near ")"
| position: 43
> exec
| SELECT 1/0;
| CREATE TABLE t (a numeric())
< error
| 42601 syntax error at or near ")"
| position: 39
> exec
| SELECT 1/0;
| SELECT 1::numeric(10, )
< error
| 42601 syntax error at or near ")"
| position: 35
> exec
| SELECT 1/0;
| SELECT 1::bit(1, )
< error
| 42601 syntax error at or near ")"
| position: 30
> exec
| SELECT 1/0;
| INSERT INTO t VALUES (1, 2), (3, 4)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT ARRAY[]
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT ROW()
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FROM t GROUP BY ()
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FROM t GROUP BY GROUPING SETS ((), (a))
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| EXECUTE p (1, 2)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| MERGE INTO t USING u ON true WHEN NOT MATCHED THEN INSERT VALUES (1, 2)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT f(1, 2,)
< error
| 42601 syntax error at or near ")"
| position: 27
> exec
| SELECT 1/0;
| SELECT f(a => 1, )
< error
| 42601 syntax error at or near ")"
| position: 30
> exec
| SELECT 1/0;
| SELECT f(a => )
< error
| 42601 syntax error at or near ")"
| position: 27
> exec
| SELECT 1/0;
| SELECT f(VARIADIC )
< error
| 42601 syntax error at or near ")"
| position: 31
> exec
| SELECT 1/0;
| SELECT count(DISTINCT )
< error
| 42601 syntax error at or near ")"
| position: 35
> exec
| SELECT 1/0;
| SELECT f()
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT f(1, VARIADIC ARRAY[2])
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT f(a => 1, b => 2)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT 1 FROM t ORDER BY a,
< error
| 42601 syntax error at end of input
| position: 40
> exec
| SELECT 1/0;
| SELECT 1 ORDER BY
< error
| 42601 syntax error at end of input
| position: 30
> exec
| SELECT 1/0;
| SELECT 1 FROM t WINDOW w AS (ORDER BY a, )
< error
| 42601 syntax error at or near ")"
| position: 54
> exec
| SELECT 1/0;
| SELECT 1 FROM t ORDER BY a, b DESC
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SELECT * FROM t,
< error
| 42601 syntax error at end of input
| position: 29
> exec
| SELECT 1/0;
| SELECT * FROM t, WHERE true
< error
| 42601 syntax error at or near "WHERE"
| position: 30
> exec
| SELECT 1/0;
| UPDATE t SET a = 1 FROM u,
< error
| 42601 syntax error at end of input
| position: 39
> exec
| SELECT 1/0;
| DELETE FROM t USING u,
< error
| 42601 syntax error at end of input
| position: 35
> exec
| SELECT 1/0;
| SELECT * FROM t, u
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE VIEW v (a, ) AS SELECT 1
< error
| 42601 syntax error at or near ")"
| position: 31
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, UNIQUE (a, ))
< error
| 42601 syntax error at or near ")"
| position: 47
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, PRIMARY KEY (a) INCLUDE (a, ))
< error
| 42601 syntax error at or near ")"
| position: 64
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, FOREIGN KEY (a, ) REFERENCES u)
< error
| 42601 syntax error at or near ")"
| position: 52
> exec
| SELECT 1/0;
| CREATE TABLE t (a int REFERENCES u (b, ))
< error
| 42601 syntax error at or near ")"
| position: 52
> exec
| SELECT 1/0;
| CREATE TABLE t (a int, EXCLUDE USING gist (a WITH =, ))
< error
| 42601 syntax error at or near ")"
| position: 66
> exec
| SELECT 1/0;
| CREATE INDEX i ON t (a, )
< error
| 42601 syntax error at or near ")"
| position: 37
> exec
| SELECT 1/0;
| CREATE INDEX i ON t (a) INCLUDE (b, )
< error
| 42601 syntax error at or near ")"
| position: 49
> exec
| SELECT 1/0;
| INSERT INTO t VALUES (1) ON CONFLICT (a, ) DO NOTHING
< error
| 42601 syntax error at or near ")"
| position: 54
> exec
| SELECT 1/0;
| COPY t FROM stdin WITH CSV FORCE NOT NULL a,
< error
| 42601 syntax error at end of input
| position: 57
> exec
| SELECT 1/0;
| CREATE TRIGGER tr AFTER UPDATE OF a, ON t FOR EACH ROW EXECUTE FUNCTION f()
< error
| 42601 syntax error at or near "ON"
| position: 50
> exec
| SELECT 1/0;
| WITH RECURSIVE a(x) AS (SELECT 1) SEARCH DEPTH FIRST BY x, SET o SELECT 1
< error
| 42601 syntax error at or near "o"
| position: 76
> exec
| SELECT 1/0;
| WITH RECURSIVE a(x) AS (SELECT 1) CYCLE x, SET c USING p SELECT 1
< error
| 42601 syntax error at or near "c"
| position: 60
> exec
| SELECT 1/0;
| CREATE VIEW v (a, b) AS SELECT 1, 2
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE INDEX i ON t (a, lower(b), (a + 1)) INCLUDE (c, d)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| COPY t FROM stdin CSV FORCE QUOTE a, b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE TRIGGER tr AFTER UPDATE OF a, b ON t FOR EACH ROW EXECUTE FUNCTION f()
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| TRUNCATE a,
< error
| 42601 syntax error at end of input
| position: 24
> exec
| SELECT 1/0;
| LOCK TABLE a,
< error
| 42601 syntax error at end of input
| position: 26
> exec
| SELECT 1/0;
| IMPORT FOREIGN SCHEMA s LIMIT TO (a, ) FROM SERVER x INTO y
< error
| 42601 syntax error at or near ")"
| position: 50
> exec
| SELECT 1/0;
| TRUNCATE a, b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| GRANT SELECT, ON t TO a
< error
| 42601 syntax error at or near "ON"
| position: 27
> exec
| SELECT 1/0;
| GRANT SELECT ON t, TO a
< error
| 42601 syntax error at or near "TO"
| position: 32
> exec
| SELECT 1/0;
| GRANT SELECT ON TABLE a, TO r
< error
| 42601 syntax error at or near "TO"
| position: 38
> exec
| SELECT 1/0;
| GRANT USAGE ON SCHEMA a, TO r
< error
| 42601 syntax error at or near "TO"
| position: 38
> exec
| SELECT 1/0;
| GRANT USAGE ON TYPE a, TO r
< error
| 42601 syntax error at or near "TO"
| position: 36
> exec
| SELECT 1/0;
| GRANT CONNECT ON DATABASE a, TO r
< error
| 42601 syntax error at or near "TO"
| position: 42
> exec
| SELECT 1/0;
| GRANT SELECT ON ALL TABLES IN SCHEMA a, TO r
< error
| 42601 syntax error at or near "TO"
| position: 53
> exec
| SELECT 1/0;
| GRANT EXECUTE ON FUNCTION f(int, ) TO a
< error
| 42601 syntax error at or near ")"
| position: 46
> exec
| SELECT 1/0;
| GRANT SELECT, INSERT ON t TO r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| GRANT SELECT (a, b), UPDATE (a) ON t TO r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| GRANT USAGE ON SCHEMA a, b TO r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| GRANT EXECUTE ON FUNCTION f(int), g() TO r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| SET search_path = a,
< error
| 42601 syntax error at end of input
| position: 33
> exec
| SELECT 1/0;
| SET search_path TO a, b,
< error
| 42601 syntax error at end of input
| position: 37
> exec
| SELECT 1/0;
| BEGIN ISOLATION LEVEL SERIALIZABLE,
< error
| 42601 syntax error at end of input
| position: 48
> exec
| SELECT 1/0;
| BEGIN ISOLATION LEVEL SERIALIZABLE READ ONLY,
< error
| 42601 syntax error at end of input
| position: 58
> exec
| SELECT 1/0;
| SET search_path = a, b
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| BEGIN ISOLATION LEVEL SERIALIZABLE, READ ONLY
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| BEGIN ISOLATION LEVEL SERIALIZABLE READ ONLY
< error
| 22012 division by zero
