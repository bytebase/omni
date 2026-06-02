# Cloud Spanner GoogleSQL truth1 — Expressions

All forms extracted from official Cloud Spanner GoogleSQL expression reference.
Primary sources:
- https://cloud.google.com/spanner/docs/reference/standard-sql/expressions
- https://cloud.google.com/spanner/docs/reference/standard-sql/functions-and-operators
version: spanner-current

---

## EXPR-001: CASE — simple form

```
CASE expression
  WHEN value THEN result
  [WHEN value THEN result ...]
  [ELSE else_result]
END
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#case_expressions

```sql
SELECT CASE Status
  WHEN 'active' THEN 1
  WHEN 'inactive' THEN 0
  ELSE -1
END AS StatusCode
FROM Accounts;
```

---

## EXPR-002: CASE — searched form

```
CASE
  WHEN bool_condition THEN result
  [WHEN bool_condition THEN result ...]
  [ELSE else_result]
END
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#case_expressions

```sql
SELECT CASE
  WHEN Score > 90 THEN 'A'
  WHEN Score > 80 THEN 'B'
  WHEN Score > 70 THEN 'C'
  ELSE 'F'
END AS Grade
FROM Students;
```

---

## EXPR-003: CAST expression

```
CAST ( expression AS type )

-- Raises an error if the cast fails
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#cast_expressions

```sql
SELECT CAST('123' AS INT64);
SELECT CAST(3.14 AS STRING);
SELECT CAST('2024-01-01' AS DATE);
SELECT CAST(TRUE AS INT64);         -- 1
SELECT CAST(CURRENT_TIMESTAMP() AS DATE);
```

---

## EXPR-004: SAFE_CAST expression

```
SAFE_CAST ( expression AS type )

-- Returns NULL instead of an error if the cast fails
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#safe_cast_expressions

```sql
SELECT SAFE_CAST('not_a_number' AS INT64);   -- NULL
SELECT SAFE_CAST('2024-01-01' AS DATE);
SELECT SAFE_CAST(NULL AS STRING);            -- NULL
```

---

## EXPR-005: EXTRACT expression

```
EXTRACT ( part FROM datetime_expression [AT TIME ZONE timezone] )

part:
    YEAR | MONTH | DAY
  | HOUR | MINUTE | SECOND | MILLISECOND | MICROSECOND | NANOSECOND
  | DATE | TIME
  | DAYOFWEEK | DAYOFYEAR
  | WEEK | WEEK(<weekday>)
  | ISOWEEK | ISOYEAR
  | QUARTER
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#extract_expressions

```sql
SELECT EXTRACT(YEAR FROM BirthDate) AS birth_year FROM Singers;
SELECT EXTRACT(HOUR FROM CURRENT_TIMESTAMP() AT TIME ZONE 'America/New_York');
SELECT EXTRACT(DATE FROM TIMESTAMP '2024-06-01 12:30:00 UTC');
SELECT EXTRACT(DAYOFWEEK FROM CURRENT_DATE());
```

---

## EXPR-006: Function call — basic

```
function_name ( [argument [, ...]] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#function_calls

```sql
SELECT CONCAT(FirstName, ' ', LastName) FROM Singers;
SELECT UPPER(LastName), LOWER(FirstName) FROM Singers;
SELECT CURRENT_TIMESTAMP(), CURRENT_DATE();
```

---

## EXPR-007: Function call — named arguments

```
function_name ( argument_name => value [, ...] )
-- Mix of positional and named:
function_name ( positional_arg, named_arg => value )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#function_calls

```sql
SELECT TOKENIZE_NGRAMS(text_col, ngram_size_min => 2, ngram_size_max => 4);
SELECT ARRAY<FLOAT32>(vector_length => 128);
```

---

## EXPR-008: Aggregate function call with DISTINCT

```
aggregate_function ( DISTINCT expression )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls

```sql
SELECT COUNT(DISTINCT SingerId) FROM Albums;
SELECT SUM(DISTINCT Score) FROM Scores;
```

---

## EXPR-009: Aggregate function with IGNORE NULLS / RESPECT NULLS

```
aggregate_function ( expression [IGNORE NULLS | RESPECT NULLS] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls

```sql
SELECT ARRAY_AGG(SingerId IGNORE NULLS) FROM Singers;
SELECT STRING_AGG(Name RESPECT NULLS ORDER BY Name) FROM Tags;
SELECT LAST_VALUE(Score IGNORE NULLS) OVER (PARTITION BY Player ORDER BY Ts) FROM Scores;
```

---

## EXPR-010: Aggregate function with ORDER BY

```
aggregate_function ( expression [ORDER BY expression [ASC|DESC]] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls

```sql
SELECT ARRAY_AGG(SingerId ORDER BY LastName) FROM Singers;
SELECT STRING_AGG(Title, ', ' ORDER BY AlbumId DESC) FROM Albums;
```

---

## EXPR-011: Aggregate function with LIMIT

```
aggregate_function ( expression [ORDER BY ...] LIMIT n )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#aggregate_function_calls

```sql
SELECT ARRAY_AGG(SingerId ORDER BY LastName LIMIT 5) FROM Singers;
```

---

## EXPR-012: Window function — OVER clause

```
window_function ( [args] ) OVER window_specification

window_specification:
    window_name
  | ( [window_name]
      [PARTITION BY expression [, ...]]
      [ORDER BY expression [{ ASC | DESC }] [, ...]]
      [window_frame_clause]
    )

window_frame_clause:
  { ROWS | RANGE } { frame_start | BETWEEN frame_start AND frame_end }

frame_start / frame_end:
    UNBOUNDED PRECEDING
  | integer_expr PRECEDING
  | CURRENT ROW
  | integer_expr FOLLOWING
  | UNBOUNDED FOLLOWING
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#window_function_calls

```sql
SELECT SingerId, AlbumId,
  ROW_NUMBER() OVER (PARTITION BY SingerId ORDER BY AlbumId) AS rn,
  SUM(MarketingBudget) OVER (PARTITION BY SingerId
                             ORDER BY AlbumId
                             ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_budget
FROM Albums;

-- Named window
SELECT AVG(Score) OVER w AS avg_score
FROM Scores
WINDOW w AS (PARTITION BY PlayerId ORDER BY Ts ROWS BETWEEN 5 PRECEDING AND CURRENT ROW);
```

---

## EXPR-013: Array subscript — OFFSET() accessor

```
array_expression [ OFFSET ( index ) ]
-- 0-based index; error if out of bounds
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator

```sql
SELECT [1, 2, 3][OFFSET(0)];       -- 1
SELECT [1, 2, 3][OFFSET(2)];       -- 3
SELECT arr[OFFSET(0)] FROM t;
```

---

## EXPR-014: Array subscript — ORDINAL() accessor

```
array_expression [ ORDINAL ( index ) ]
-- 1-based index; error if out of bounds
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator

```sql
SELECT [1, 2, 3][ORDINAL(1)];      -- 1
SELECT [1, 2, 3][ORDINAL(3)];      -- 3
```

---

## EXPR-015: Array subscript — SAFE_OFFSET() accessor

```
array_expression [ SAFE_OFFSET ( index ) ]
-- 0-based; returns NULL if out of bounds
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator

```sql
SELECT [1, 2, 3][SAFE_OFFSET(10)];  -- NULL (no error)
```

---

## EXPR-016: Array subscript — SAFE_ORDINAL() accessor

```
array_expression [ SAFE_ORDINAL ( index ) ]
-- 1-based; returns NULL if out of bounds
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subscript_operator

```sql
SELECT [1, 2, 3][SAFE_ORDINAL(10)];  -- NULL
```

---

## EXPR-017: Array constructor

```
ARRAY [ element [, ...] ]
ARRAY < type > [ element [, ...] ]
[ element [, ...] ]            -- shorthand array literal
ARRAY ( subquery )             -- array from subquery
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_constructors

```sql
SELECT [1, 2, 3];
SELECT ARRAY<INT64>[1, 2, 3];
SELECT ARRAY['a', 'b', 'c'];
SELECT ARRAY(SELECT SingerId FROM Singers ORDER BY SingerId);
```

---

## EXPR-018: STRUCT constructor

```
STRUCT ( value [AS field_name] [, ...] )
STRUCT < field_type [, ...] > ( value [, ...] )
( value, value [, ...] )    -- tuple/struct shorthand in some contexts
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#struct_constructors

```sql
SELECT STRUCT(1, 'hello', TRUE);
SELECT STRUCT(1 AS id, 'Alice' AS name);
SELECT STRUCT<INT64, STRING(100)>(1, 'hello');
SELECT (1, 2, 3);   -- anonymous struct
```

---

## EXPR-019: Field access on STRUCT

```
struct_expression . field_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#field_access_operator

```sql
SELECT s.name FROM StructTable AS s;
SELECT STRUCT(1 AS id, 'Alice' AS name).name;    -- 'Alice'
SELECT (SELECT AS STRUCT SingerId, FirstName FROM Singers LIMIT 1).FirstName;
```

---

## EXPR-020: INTERVAL literals

```
INTERVAL integer { YEAR | MONTH | DAY | HOUR | MINUTE | SECOND }
INTERVAL 'value' { part TO part }
-- Used with date/timestamp arithmetic
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#interval_expressions

```sql
SELECT TIMESTAMP_ADD(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR);
SELECT DATE_ADD(CURRENT_DATE(), INTERVAL 30 DAY);
SELECT TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 7 DAY);
SELECT DATE_DIFF(DATE '2024-12-31', DATE '2024-01-01', DAY);
```

---

## EXPR-021: Conditional expressions (IF, IFNULL, NULLIF, COALESCE, IIF)

```
IF ( condition, true_result, false_result )
IFNULL ( expression, null_replacement )
NULLIF ( expression, expression_to_match )
COALESCE ( expression [, ...] )
IIF ( condition, true_result, false_result )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#conditional_expressions

```sql
SELECT IF(Score > 50, 'pass', 'fail') FROM Students;
SELECT IFNULL(NickName, FirstName) FROM Singers;
SELECT NULLIF(Status, 'unknown') FROM Orders;
SELECT COALESCE(NickName, FirstName, 'Anonymous') FROM Singers;
```

---

## EXPR-022: IN expressions

```
expression [NOT] IN ( value [, ...] )
expression [NOT] IN ( subquery )
expression [NOT] IN UNNEST ( array_expression )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#in_operators

```sql
SELECT * FROM Singers WHERE SingerId IN (1, 2, 3);
SELECT * FROM Singers WHERE SingerId NOT IN (SELECT SingerId FROM Banned);
SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids);
```

---

## EXPR-023: EXISTS expression

```
EXISTS ( subquery )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#exists_expression

```sql
SELECT * FROM Singers
WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerId = Singers.SingerId);
```

---

## EXPR-024: LIKE and NOT LIKE expressions

```
expression [NOT] LIKE pattern
-- % matches any sequence; _ matches single character
-- No regex; only SQL LIKE patterns
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#like_operator

```sql
SELECT * FROM Singers WHERE LastName LIKE 'Sm%';
SELECT * FROM Singers WHERE FirstName NOT LIKE '_a%';
```

---

## EXPR-025: IS NULL / IS NOT NULL

```
expression IS [NOT] NULL
expression IS [NOT] TRUE
expression IS [NOT] FALSE
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#is_expressions

```sql
SELECT * FROM Singers WHERE BirthDate IS NULL;
SELECT * FROM Orders WHERE Processed IS NOT NULL;
SELECT * FROM Flags WHERE Active IS TRUE;
```

---

## EXPR-026: BETWEEN expression

```
expression [NOT] BETWEEN lower AND upper
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#between_operator

```sql
SELECT * FROM Singers WHERE SingerId BETWEEN 1 AND 100;
SELECT * FROM Albums WHERE MarketingBudget NOT BETWEEN 0 AND 10000;
```

---

## EXPR-027: Arithmetic operators

```
expression + expression
expression - expression
expression * expression
expression / expression
- expression    (unary minus)
+ expression    (unary plus)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#arithmetic_operators

```sql
SELECT 1 + 2, 10 - 3, 4 * 5, 8 / 2;
SELECT -Score, +Amount FROM t;
```

---

## EXPR-028: Comparison operators

```
expression = expression
expression != expression  (also <>)
expression < expression
expression <= expression
expression > expression
expression >= expression
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#comparison_operators

```sql
SELECT * FROM t WHERE a = b;
SELECT * FROM t WHERE x != y;
SELECT * FROM t WHERE Score >= 90;
```

---

## EXPR-029: Logical operators

```
bool_expression AND bool_expression
bool_expression OR bool_expression
NOT bool_expression
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#logical_operators

```sql
SELECT * FROM t WHERE a AND b;
SELECT * FROM t WHERE x OR y;
SELECT * FROM t WHERE NOT active;
```

---

## EXPR-030: Bitwise operators

```
expression & expression    -- bitwise AND
expression | expression    -- bitwise OR
expression ^ expression    -- bitwise XOR
~expression                -- bitwise NOT
expression << expression   -- left shift
expression >> expression   -- right shift
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#bitwise_operators

```sql
SELECT 5 & 3, 5 | 3, 5 ^ 3, ~5;
SELECT flags << 2, flags >> 1 FROM t;
```

---

## EXPR-031: String concatenation operator

```
string1 || string2
-- Concatenates two strings (or arrays)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#concatenation_operator

```sql
SELECT 'Hello' || ', ' || 'World';
SELECT FirstName || ' ' || LastName AS FullName FROM Singers;
SELECT [1,2] || [3,4];   -- array concatenation: [1,2,3,4]
```

---

## EXPR-032: Scalar subquery

```
( SELECT expression FROM ... [WHERE ...] )
-- Must return exactly 0 or 1 row; 0 rows = NULL; >1 rows = error
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#scalar_subqueries

```sql
SELECT (SELECT MAX(MarketingBudget) FROM Albums WHERE SingerId = s.SingerId)
FROM Singers s;
```

---

## EXPR-033: ARRAY subquery

```
ARRAY ( SELECT expression FROM ... )
-- Constructs an array from a subquery; empty subquery = empty array
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#array_subqueries

```sql
SELECT ARRAY(SELECT AlbumId FROM Albums WHERE SingerId = s.SingerId) AS albums
FROM Singers s;
```

---

## EXPR-034: ALL, SOME, ANY quantified subquery

```
expression comparison_op { ALL | SOME | ANY } ( subquery )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#quantified_subqueries

```sql
SELECT * FROM t WHERE x > ALL (SELECT score FROM scores);
SELECT * FROM t WHERE x = ANY (SELECT val FROM vals);
```

---

## EXPR-035: AT TIME ZONE

```
timestamp_expression AT TIME ZONE timezone_string
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#at_time_zone

```sql
SELECT EXTRACT(HOUR FROM CURRENT_TIMESTAMP() AT TIME ZONE 'America/New_York');
SELECT TIMESTAMP '2024-06-01 12:00:00 UTC' AT TIME ZONE 'Europe/London';
```

---

## EXPR-036: PENDING_COMMIT_TIMESTAMP() — special function

```
PENDING_COMMIT_TIMESTAMP()
-- Used as DEFAULT for TIMESTAMP columns with allow_commit_timestamp=true
-- Returns a placeholder that becomes the actual commit timestamp at write time
```

source_url: https://cloud.google.com/spanner/docs/commit-timestamp

```sql
INSERT INTO Events (EventId, LastMod) VALUES (1, PENDING_COMMIT_TIMESTAMP());
UPDATE Events SET LastMod = PENDING_COMMIT_TIMESTAMP() WHERE EventId = 1;
-- Also usable as DEFAULT in column definition:
-- col TIMESTAMP DEFAULT (PENDING_COMMIT_TIMESTAMP()) OPTIONS (allow_commit_timestamp = true)
```

---

## EXPR-037: NEXT VALUE FOR / GET_NEXT_SEQUENCE_VALUE

```
NEXT VALUE FOR sequence_name
GET_NEXT_SEQUENCE_VALUE(SEQUENCE sequence_name)
```

source_url: https://cloud.google.com/spanner/docs/sequences

```sql
INSERT INTO Orders (OrderId, Description) VALUES (NEXT VALUE FOR OrderSeq, 'test');
SELECT GET_NEXT_SEQUENCE_VALUE(SEQUENCE MySeq);
-- As DEFAULT: col INT64 DEFAULT (NEXT VALUE FOR MySeq)
```

---

## EXPR-038: IN UNNEST (array parameter binding)

```
expression [NOT] IN UNNEST ( array_expression )
-- Commonly used to check membership in an array parameter
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#in_operators

```sql
SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids);
SELECT * FROM Albums WHERE Title NOT IN UNNEST(@excluded_titles);
```

---

## EXPR-039: NEW proto constructor

```
NEW proto_type_name ( field_name: value [, ...] )
-- Constructs a proto message value
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#proto_constructors

```sql
SELECT NEW `my.package.MyMessage`(field1: 'value', field2: 42);
```

---

## EXPR-040: TREAT expression (proto type conversion)

```
TREAT ( expression AS proto_type )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/expressions#treat_expression

```sql
SELECT TREAT(col AS `my.package.SubType`) FROM t;
```
