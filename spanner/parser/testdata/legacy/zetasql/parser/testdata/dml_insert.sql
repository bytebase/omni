INSERT INTO T
VALUES
  (1);
INSERT OR IGNORE INTO values
VALUES
  (5, 6),
  ('abc')
ASSERT_ROWS_MODIFIED 5;
INSERT INTO values
VALUES
  (a, b);
INSERT INTO mytable
VALUES
  (a, b);
INSERT INTO values(a, `default`)
VALUES
  (1, 3);
INSERT OR IGNORE INTO `replace`
SELECT
  1;
INSERT INTO insert
SELECT
  1;
INSERT INTO `replace`
SELECT
  1;
INSERT INTO T(c1)
WITH
  q1 AS (
    SELECT
      1
  )
SELECT
  *
FROM
  q1
ASSERT_ROWS_MODIFIED @row_count;
INSERT INTO T
VALUES
  ((1, 2), (
      SELECT
        1
    ));
INSERT INTO T
VALUES
  (1, 2),
  ((
      SELECT
        1
    ));
INSERT INTO T
WITH
  q1 AS (
    SELECT
      1
  )
SELECT
  *
FROM
  q1;
INSERT OR REPLACE INTO T
SELECT
  *
FROM
  q1
UNION ALL
SELECT
  *
FROM
  q2
ORDER BY x
LIMIT 5
ASSERT_ROWS_MODIFIED 0;
INSERT INTO T(c1, c2)
VALUES
  (1 + 2, f(), null, DEFAULT),
  (DEFAULT, (
      SELECT
        x
      FROM
        y
    ));
INSERT INTO T
(
SELECT
  5
);
INSERT INTO T
((
SELECT
  1
) UNION ALL(
SELECT
  2
));
INSERT INTO T(c1, c2)
(
SELECT
  5,
  7
) UNION ALL
SELECT
  8,
  9;
INSERT INTO T.(a.b).c
VALUES
  (1);
INSERT INTO T.a[0].b
VALUES
  (1);
@{ a = b }
INSERT INTO t1
VALUES
  (null);
