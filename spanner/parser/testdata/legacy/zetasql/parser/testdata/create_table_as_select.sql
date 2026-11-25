CREATE TABLE t1 AS
SELECT
  *
FROM
  t2;
CREATE TABLE t1 LIKE
t2 AS
SELECT
  *
FROM
  t2;
CREATE TABLE t1
(
  x int64
) LIKE
t2 AS
SELECT
  *
FROM
  t2;
CREATE PRIVATE TABLE pkg1.pkg2.`t 2` AS
SELECT
  1 AS a;
CREATE TEMP TABLE pkg1.pkg2.`t 2` AS
SELECT
  1
FROM
  t2
UNION ALL
SELECT
  2
FROM
  t3
  CROSS JOIN
  t4;
CREATE TABLE t3 AS
WITH
  tt AS (
    SELECT
      *
    FROM
      KeyValue
  )
SELECT
  *
FROM
  tt;
CREATE TABLE t4 OPTIONS(x = y) AS
WITH
  t1 AS (
    SELECT
      1
  ),
  t2 AS (
    SELECT
      2
  )
SELECT
  3;
CREATE TABLE tt OPTIONS(x = 1) AS
SELECT
  1;
CREATE TEMP TABLE tt OPTIONS(x = 5, y = 'abc', z = @param, zz = ident, zzz = @@sysvar) AS
SELECT
  2;
CREATE TABLE tt OPTIONS() AS
SELECT
  2;
CREATE TABLE tt OPTIONS(x = 5.5, y = a, z = b.c) AS
SELECT
  2;
CREATE TABLE tt OPTIONS(y = 'b.c', z = `b.c`) AS
SELECT
  2;
CREATE TABLE options AS
SELECT
  1 AS x;
CREATE TABLE options OPTIONS(x = y) AS
SELECT
  1 AS x;
CREATE TABLE options OPTIONS(x = y) AS
SELECT
  1 AS x;
CREATE OR REPLACE TABLE xyz AS
(
SELECT
  1
);
CREATE TABLE IF NOT EXISTS xyz AS
(
SELECT
  1
);
CREATE OR REPLACE TEMP TABLE IF NOT EXISTS a.b.c OPTIONS(d = e) AS
SELECT
  1;
CREATE TABLE t AS
WITH
  q AS (
    SELECT
      1
  ),
  q2 AS (
    SELECT
      *
    FROM
      q
  )
SELECT
  *
FROM
  q2;
CREATE TABLE t1
(
  a int64,
  b string
) AS
SELECT
  1 AS a,
  'hi' AS b;
CREATE TEMP TABLE t
(
  param1 int64,
  param2 int64
) AS
SELECT
  1,
  2,
  3;
CREATE TABLE `function` AS
SELECT
  1,
  2,
  3;
CREATE TABLE t1 PARTITION BY key, value AS
SELECT
  key,
  value
FROM
  KeyValue;
CREATE TABLE t1
(
  a int64,
  b string
) PARTITION BY b AS
SELECT
  1 AS a,
  'hi' AS b;
CREATE TABLE t1 CLUSTER BY key, value AS
SELECT
  key,
  value
FROM
  KeyValue;
CREATE TABLE t1
(
  a int64,
  b string
) CLUSTER BY b AS
SELECT
  1 AS a,
  'hi' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) PARTITION BY b CLUSTER BY a OPTIONS(key = 'value') AS
SELECT
  1 AS a,
  'hi' AS b;
CREATE TABLE t1 AS
SELECT
  1 AS a,
  2 AS b;
CREATE TABLE t1
(
  a int64,
  b string
) AS
SELECT
  1 AS a,
  'foo' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) WITH CONNECTION DEFAULT AS
SELECT
  1 AS a,
  'foo' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) WITH CONNECTION connection1 AS
SELECT
  1 AS a,
  'foo' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) WITH CONNECTION connection1 OPTIONS(foo = x, bar = b) AS
SELECT
  1 AS a,
  'foo' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) WITH CONNECTION us.connection1 OPTIONS(foo = x, bar = b) AS
SELECT
  1 AS a,
  'foo' AS b;
CREATE TABLE t1
(
  a int64,
  b string
) WITH CONNECTION myproject.us.connection1 OPTIONS(foo = x, bar = b) AS
SELECT
  1 AS a,
  'foo' AS b;
