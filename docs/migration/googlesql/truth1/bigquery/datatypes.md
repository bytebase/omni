# BigQuery GoogleSQL truth1 — Data Types

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types

---

## DT-001: Scalar data types

```
-- Integer
INT64   (aliases: INT, SMALLINT, INTEGER, BIGINT, TINYINT, BYTEINT)
Range: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807

-- Decimal / Fixed-point
NUMERIC   (alias: DECIMAL)
  Precision: 38, Scale: 9
  Range: -9.9999999999999999999999999999999999999E+28 to 9.9999999999999999999999999999999999999E+28

BIGNUMERIC  (alias: BIGDECIMAL)
  Precision: ~76.8 digits, Scale: 38

-- Floating point
FLOAT64
  Double precision IEEE 754 approximate

-- Boolean
BOOL   (alias: BOOLEAN)
  Values: TRUE, FALSE

-- String
STRING
  Variable-length Unicode (UTF-8)

-- Bytes
BYTES
  Variable-length binary

-- Date/Time
DATE       -- Gregorian calendar date, no time zone
           -- Range: 0001-01-01 to 9999-12-31
DATETIME   -- Date and time, no time zone
           -- Range: 0001-01-01 00:00:00 to 9999-12-31 23:59:59.999999
TIME       -- Time of day, no time zone
           -- Range: 00:00:00 to 23:59:59.999999
TIMESTAMP  -- Absolute point in time (microsecond precision), UTC
           -- Range: 0001-01-01 00:00:00 to 9999-12-31 23:59:59.999999 UTC

-- Interval
INTERVAL   -- Duration of time (without reference to a specific point in time)

-- JSON
JSON

-- Geography
GEOGRAPHY  -- Collection of points, linestrings, polygons on Earth's surface
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#data_type_list
version: bigquery-current

```sql
SELECT INT64 '42', NUMERIC '3.14', BIGNUMERIC '1.23456789', FLOAT64 '1.5';
SELECT BOOL 'TRUE', STRING 'hello', BYTES b'\x00\x01', DATE '2024-01-01';
SELECT TIMESTAMP '2024-01-01 12:00:00 UTC', INTERVAL 5 DAY, JSON '{"k": 1}';
```

---

## DT-002: Parameterized data types

```
-- Parameterized string
STRING(L)              -- max L Unicode characters

-- Parameterized bytes
BYTES(L)               -- max L bytes

-- Parameterized numeric
NUMERIC(P[, S])        -- precision P (max 29+S), scale S (0-9)
DECIMAL(P[, S])        -- alias for NUMERIC(P[,S])
BIGNUMERIC(P[, S])     -- precision P (max 38+S), scale S (0-38)
BIGDECIMAL(P[, S])     -- alias for BIGNUMERIC(P[,S])
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#parameterized_data_types
version: bigquery-current

```sql
CREATE TABLE t (
  name STRING(100),
  data BYTES(256),
  price NUMERIC(10, 2),
  exact_val BIGNUMERIC(38, 10)
);
```

---

## DT-003: ARRAY type

```
ARRAY<T>

-- T can be any non-ARRAY type (including STRUCT)
-- Arrays of arrays are not allowed directly; use ARRAY<STRUCT<ARRAY<T>>>
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#array_type
version: bigquery-current

```sql
-- Declare: ARRAY<INT64>, ARRAY<STRING>, ARRAY<STRUCT<INT64, INT64>>
SELECT [1, 2, 3] AS nums;
SELECT ARRAY<STRING>['a', 'b', 'c'] AS letters;
SELECT ARRAY<INT64>[] AS empty_arr;
```

---

## DT-004: STRUCT type

```
STRUCT<[field_name] field_type [, ...]>

-- Anonymous fields:  STRUCT<INT64, STRING>
-- Named fields:      STRUCT<x INT64, y STRING>
-- Nested:            STRUCT<x STRUCT<y INT64, z INT64>>
-- With array:        STRUCT<inner_array ARRAY<INT64>>
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#struct_type
version: bigquery-current

```sql
SELECT STRUCT<x INT64, y STRING>(1, 'hello') AS s;
SELECT STRUCT(1 AS a, 'abc' AS b) AS s;
SELECT (1, 'hello') AS s;  -- tuple syntax (2+ fields)
SELECT STRUCT<INT64, STRING>(42, 'world') AS s;
```

---

## DT-005: RANGE type

```
RANGE<T>  where T is DATE, DATETIME, or TIMESTAMP

-- Range literal:
RANGE<T> '[lower_bound, upper_bound)'
-- lower_bound and upper_bound can be UNBOUNDED or NULL
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#range_type
version: bigquery-current

```sql
SELECT RANGE<DATE> '[2024-01-01, 2024-12-31)' AS r;
SELECT RANGE<TIMESTAMP> '[2024-01-01 00:00:00 UTC, UNBOUNDED)' AS r;
INSERT mytable (emp_id, duration)
VALUES (10, RANGE<DATE> '[2010-01-10, 2010-03-10)');
```

---

## DT-006: Type aliases

```
-- Integer aliases (all equivalent to INT64):
INT, SMALLINT, INTEGER, BIGINT, TINYINT, BYTEINT

-- Decimal aliases:
DECIMAL   = NUMERIC
BIGDECIMAL = BIGNUMERIC

-- Boolean alias:
BOOLEAN = BOOL
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-types#data_type_list
version: bigquery-current

```sql
SELECT CAST(1 AS INT), CAST(3.14 AS DECIMAL), CAST(TRUE AS BOOLEAN);
```
