# BigQuery GoogleSQL truth1 — Special-Syntax Expressions

All forms extracted from official BigQuery GoogleSQL documentation.

---

## EXPR-001: CASE expression (searched)

```
CASE
  WHEN condition THEN result
  [WHEN condition THEN result ...]
  [ELSE else_result]
END
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/conditional_expressions#case_expr
version: bigquery-current

```sql
SELECT CASE
  WHEN x > 0 THEN 'positive'
  WHEN x < 0 THEN 'negative'
  ELSE 'zero'
END AS sign
FROM t;
```

---

## EXPR-002: CASE expression (simple)

```
CASE expr
  WHEN value THEN result
  [WHEN value THEN result ...]
  [ELSE else_result]
END
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/conditional_expressions#case_expr
version: bigquery-current

```sql
SELECT CASE status
  WHEN 'A' THEN 'Active'
  WHEN 'I' THEN 'Inactive'
  ELSE 'Unknown'
END AS status_label
FROM accounts;
```

---

## EXPR-003: CAST expression

```
CAST(expr AS type)

-- Safe cast (returns NULL on failure instead of error):
SAFE_CAST(expr AS type)

-- Format clause (for date/time conversions):
CAST(expr AS type FORMAT format_string [AT TIME ZONE tz])
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/conversion_functions#cast
version: bigquery-current

```sql
SELECT CAST(3.14 AS INT64);
SELECT CAST('2024-01-01' AS DATE);
SELECT SAFE_CAST('not a number' AS INT64);  -- returns NULL
SELECT CAST(date_col AS STRING FORMAT 'YYYY/MM/DD');
SELECT CAST('01/01/2024' AS DATE FORMAT 'MM/DD/YYYY');
```

---

## EXPR-004: EXTRACT expression

```
EXTRACT(part FROM datetime_expression [AT TIME ZONE tz])

-- part: MICROSECOND, MILLISECOND, SECOND, MINUTE, HOUR, DAYOFWEEK, DAY,
--       DAYOFYEAR, WEEK, WEEK(<WEEKDAY>), ISOWEEK, MONTH, QUARTER, YEAR,
--       ISOYEAR, DATE, DATETIME, TIME
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/date_functions#extract
version: bigquery-current

```sql
SELECT EXTRACT(YEAR FROM DATE '2014-09-27');  -- 2014
SELECT EXTRACT(MONTH FROM DATETIME '2014-09-27 12:30:00');  -- 9
SELECT EXTRACT(HOUR FROM TIMESTAMP '2014-09-27 12:30:00 UTC' AT TIME ZONE 'America/Los_Angeles');
SELECT EXTRACT(DAY FROM date_col), EXTRACT(WEEK FROM date_col);
```

---

## EXPR-005: Aggregate function call (with modifiers)

```
function_name([DISTINCT] [argument_list])
  [IGNORE NULLS | RESPECT NULLS]
  [ORDER BY expression [ASC | DESC] [, ...]]
  [LIMIT n]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/aggregate_functions
version: bigquery-current

```sql
SELECT COUNT(DISTINCT user_id) FROM events;
SELECT ARRAY_AGG(x IGNORE NULLS) FROM t;
SELECT ARRAY_AGG(x ORDER BY x DESC LIMIT 10) FROM t;
SELECT STRING_AGG(name, ', ' ORDER BY name) FROM users;
SELECT PERCENTILE_CONT(score, 0.5) OVER () FROM scores;
```

---

## EXPR-006: Window function (OVER clause)

```
function_name ( [ argument_list ] ) OVER over_clause

over_clause:
  { named_window | ( [ window_specification ] ) }

window_specification:
  [ named_window ]
  [ PARTITION BY partition_expression [, ...] ]
  [ ORDER BY expression [ { ASC | DESC } ] [, ...] ]
  [ window_frame_clause ]

window_frame_clause:
  { rows_range } { frame_start | frame_between }

rows_range:
  { ROWS | RANGE }

frame_between:
  BETWEEN { UNBOUNDED PRECEDING | numeric_expression PRECEDING | CURRENT ROW | numeric_expression FOLLOWING }
    AND
    { numeric_expression PRECEDING | CURRENT ROW | numeric_expression FOLLOWING | UNBOUNDED FOLLOWING }

frame_start:
  UNBOUNDED PRECEDING | numeric_expression PRECEDING | CURRENT ROW
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/window-function-calls
version: bigquery-current

```sql
SELECT SUM(amount) OVER () AS grand_total FROM sales;
SELECT SUM(amount) OVER (PARTITION BY dept) AS dept_total FROM sales;
SELECT SUM(amount) OVER (ORDER BY date ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total FROM sales;
SELECT RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS rank FROM employees;
SELECT AVG(price) OVER (ORDER BY date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) AS moving_avg FROM prices;
SELECT LAST_VALUE(name) OVER (PARTITION BY cat ORDER BY score ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING) FROM t;
```

---

## EXPR-007: Array element access (OFFSET / ORDINAL / SAFE_OFFSET / SAFE_ORDINAL)

```
array_expr[OFFSET(index)]          -- 0-based, error on out-of-range
array_expr[ORDINAL(index)]         -- 1-based, error on out-of-range
array_expr[SAFE_OFFSET(index)]     -- 0-based, NULL on out-of-range
array_expr[SAFE_ORDINAL(index)]    -- 1-based, NULL on out-of-range
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/arrays#accessing_array_elements
version: bigquery-current

```sql
SELECT arr[OFFSET(0)] AS first FROM t;
SELECT arr[ORDINAL(1)] AS first FROM t;
SELECT arr[SAFE_OFFSET(5)] AS maybe_sixth FROM t;
SELECT arr[SAFE_ORDINAL(10)] AS maybe_tenth FROM t;
-- In FROM clause:
SELECT * FROM t, t.array_col AS item WITH OFFSET AS pos;
```

---

## EXPR-008: Struct field access

```
struct_expr.field_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#struct_type
version: bigquery-current

```sql
SELECT s.name, s.address.city FROM users AS u, u.info AS s;
SELECT (SELECT x.a FROM (SELECT STRUCT(1 AS a, 2 AS b) AS x)).a;
```

---

## EXPR-009: AT TIME ZONE

```
timestamp_expr AT TIME ZONE timezone_string

-- Used in EXTRACT, CAST, and timestamp comparisons
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/timestamp_functions#at_time_zone
version: bigquery-current

```sql
SELECT EXTRACT(HOUR FROM TIMESTAMP '2014-09-27 12:30:00 UTC' AT TIME ZONE 'America/Los_Angeles');
SELECT CAST(ts AS STRING FORMAT 'YYYY-MM-DD HH24:MI:SS' AT TIME ZONE 'US/Pacific') FROM t;
```

---

## EXPR-010: Scalar subquery

```
( subquery )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/expression_subqueries#scalar_subqueries
version: bigquery-current

```sql
SELECT username,
  (SELECT mascot FROM Mascots WHERE Players.team = Mascots.team) AS player_mascot
FROM Players;

SELECT username, level,
  (SELECT AVG(level) FROM Players) AS avg_level
FROM Players;
```

---

## EXPR-011: Array subquery

```
ARRAY( subquery )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/expression_subqueries#array_subqueries
version: bigquery-current

```sql
SELECT ARRAY(SELECT username FROM NPCs WHERE team = 'red') AS red;
SELECT ARRAY(SELECT AS STRUCT id, name FROM users) AS user_structs;
```

---

## EXPR-012: IN subquery / NOT IN subquery

```
value [ NOT ] IN ( subquery )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/expression_subqueries#in_subqueries
version: bigquery-current

```sql
SELECT 'corba' IN (SELECT username FROM Players) AS result;
SELECT id FROM t WHERE status NOT IN (SELECT status FROM blocked);
```

---

## EXPR-013: EXISTS subquery

```
EXISTS( subquery )
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/expression_subqueries#exists_subqueries
version: bigquery-current

```sql
SELECT EXISTS(SELECT 1 FROM Players WHERE team = 'yellow') AS result;
SELECT * FROM t WHERE EXISTS(SELECT 1 FROM related WHERE related.id = t.id);
```

---

## EXPR-014: Table subquery (in FROM)

```
FROM ( subquery ) [ [ AS ] alias ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/expression_subqueries#table_subqueries
version: bigquery-current

```sql
SELECT * FROM (SELECT id, name FROM users WHERE active = true) AS active_users;
SELECT t.*, sq.avg_val FROM t JOIN (SELECT grp, AVG(val) AS avg_val FROM data GROUP BY grp) AS sq ON t.grp = sq.grp;
```

---

## EXPR-015: INTERVAL literal expression

```
INTERVAL int64_expression datetime_part
INTERVAL datetime_parts_string starting_datetime_part TO ending_datetime_part
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#interval_literals
version: bigquery-current

```sql
SELECT TIMESTAMP_ADD(ts, INTERVAL 1 DAY) FROM t;
SELECT TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 7 DAY);
SELECT DATE_ADD(d, INTERVAL '1-2' YEAR TO MONTH);
SELECT DATE_DIFF(d1, d2, DAY);
```

---

## EXPR-016: Operators and precedence

```
-- Unary:
- expr        (negation)
NOT expr      (logical NOT)
~ expr        (bitwise NOT)

-- Arithmetic:
expr * expr
expr / expr
expr + expr
expr - expr

-- Bitwise:
expr << expr  (shift left)
expr >> expr  (shift right)
expr & expr   (bitwise AND)
expr ^ expr   (bitwise XOR)
expr | expr   (bitwise OR)

-- Comparison:
expr = expr
expr != expr  (or expr <> expr)
expr < expr
expr > expr
expr <= expr
expr >= expr
expr LIKE expr
expr NOT LIKE expr
expr IN (...)
expr NOT IN (...)
expr BETWEEN expr AND expr
expr NOT BETWEEN expr AND expr
IS NULL
IS NOT NULL
IS TRUE / IS NOT TRUE
IS FALSE / IS NOT FALSE

-- Logical:
AND
OR

-- Conditional:
expr IS DISTINCT FROM expr
expr IS NOT DISTINCT FROM expr

-- Concatenation:
expr || expr   (string/array concatenation)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/operators
version: bigquery-current

```sql
SELECT -x, NOT b, ~n, a + b, a - b, a * b, a / b;
SELECT a << 2, a >> 1, a & b, a | b, a ^ b;
SELECT x = 1, x != 2, x < 3, x >= 0;
SELECT name LIKE 'A%', id IN (1, 2, 3), age BETWEEN 18 AND 65;
SELECT x IS NULL, y IS NOT NULL, z IS TRUE;
SELECT 'foo' || 'bar', [1,2] || [3,4];
SELECT x IS DISTINCT FROM y;
```
