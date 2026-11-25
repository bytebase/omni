# Cloud Spanner GoogleSQL truth1 — Data Types

All forms extracted from official Cloud Spanner GoogleSQL data types reference.
Primary source: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types
version: spanner-current

---

## TYPE-001: BOOL

```
BOOL
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#boolean_type

```sql
col1 BOOL
col1 BOOL NOT NULL
-- literals: TRUE, FALSE
SELECT TRUE, FALSE, NULL;
```

---

## TYPE-002: INT64

```
INT64
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#integer_type

```sql
col1 INT64
col1 INT64 NOT NULL
-- Range: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
SELECT 42, -100, 0x1A2B;
```

---

## TYPE-003: FLOAT32

```
FLOAT32
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#floating_point_type

```sql
col1 FLOAT32
-- 32-bit IEEE 754 floating point
-- Special values: '+inf', '-inf', 'nan' (case insensitive)
SELECT CAST(1.5 AS FLOAT32);
```

---

## TYPE-004: FLOAT64

```
FLOAT64
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#floating_point_type

```sql
col1 FLOAT64
-- 64-bit IEEE 754 floating point
SELECT 3.14, 1.23e10, CAST('nan' AS FLOAT64);
```

---

## TYPE-005: NUMERIC

```
NUMERIC
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#numeric_type

```sql
col1 NUMERIC
-- Precision: 38 digits, scale: 9 decimal digits
-- Range: -9.9999999999999999999999999999999999999E+28 to 9.9999999999999999999999999999999999999E+28
SELECT NUMERIC '123.456789';
SELECT CAST('123.456' AS NUMERIC);
```

---

## TYPE-006: STRING(length) and STRING(MAX)

```
STRING ( { length | MAX } )

length: positive integer literal
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#string_type

```sql
col1 STRING(100)
col2 STRING(MAX)
-- Variable-length unicode text, max 2,621,440 bytes encoded in UTF-8
SELECT 'hello world';
SELECT N'unicode string';
```

---

## TYPE-007: BYTES(length) and BYTES(MAX)

```
BYTES ( { length | MAX } )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#bytes_type

```sql
col1 BYTES(256)
col2 BYTES(MAX)
-- Variable-length binary data, max 10,485,760 bytes
SELECT B'binary data', b'abc';
```

---

## TYPE-008: JSON

```
JSON
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#json_type

```sql
col1 JSON
-- Stores valid JSON in a normalized, canonical format
SELECT JSON '{"key": "value", "n": 42}';
SELECT JSON_VALUE(col1, '$.key') FROM t;
```

---

## TYPE-009: DATE

```
DATE
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#date_type

```sql
col1 DATE
-- Valid range: 0001-01-01 to 9999-12-31
SELECT DATE '2024-06-01';
SELECT CURRENT_DATE();
```

---

## TYPE-010: TIMESTAMP

```
TIMESTAMP
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#timestamp_type

```sql
col1 TIMESTAMP
col2 TIMESTAMP OPTIONS (allow_commit_timestamp = true)
-- Stores UTC timestamp, precision up to nanoseconds
-- Special value: PENDING_COMMIT_TIMESTAMP()
SELECT TIMESTAMP '2024-06-01 12:00:00 UTC';
SELECT CURRENT_TIMESTAMP();
SELECT PENDING_COMMIT_TIMESTAMP();
```

---

## TYPE-011: ARRAY

```
ARRAY < element_type >

-- Vector array with fixed length:
ARRAY < FLOAT32 > ( vector_length => integer_literal )
ARRAY < FLOAT64 > ( vector_length => integer_literal )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#array_type

```sql
col1 ARRAY<INT64>
col2 ARRAY<STRING(MAX)>
col3 ARRAY<FLOAT32>(vector_length=>128)
-- Arrays cannot be NULL at top level (the array itself can be NULL, but cannot have NULL elements unless explicitly cast)
SELECT [1, 2, 3];
SELECT ARRAY<INT64>[1, 2, 3];
SELECT ARRAY(SELECT SingerId FROM Singers);
```

---

## TYPE-012: STRUCT

```
STRUCT < field_name type [, ...] >
STRUCT < type [, ...] >        -- anonymous fields
STRUCT < >                     -- empty struct
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#struct_type

```sql
col1 STRUCT<FirstName STRING(100), LastName STRING(100), Age INT64>
-- Struct literals
SELECT STRUCT(1, 2, 3);
SELECT STRUCT<INT64, STRING(100)>(1, 'hello');
SELECT STRUCT(1 AS x, 'abc' AS y);
```

---

## TYPE-013: TOKENLIST

```
TOKENLIST
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#tokenlist_type

```sql
-- Used for full-text and vector search indexes
-- Typically created as a generated stored column using TOKENIZE_* functions
col1 TOKENLIST AS (TOKENIZE_FULLTEXT(text_col)) STORED HIDDEN
col2 TOKENLIST AS (TOKENIZE_NGRAMS(text_col, ngram_size_min=>2, ngram_size_max=>3)) STORED HIDDEN
```

---

## TYPE-014: PROTO (protocol buffer types)

```
proto_type_name
-- Fully qualified proto type name, backtick-quoted if contains dots

-- Must first be loaded via CREATE PROTO BUNDLE
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#proto_type

```sql
-- After CREATE PROTO BUNDLE (`my.package.MyMessage`)
col1 `my.package.MyMessage`
SELECT col1.field_name FROM t;
```

---

## TYPE-015: INTERVAL

```
INTERVAL

-- Interval literals:
INTERVAL integer { YEAR | MONTH | DAY | HOUR | MINUTE | SECOND }
INTERVAL string { YEAR | MONTH | DAY | HOUR | MINUTE | SECOND | ... TO ... }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#interval_type

```sql
-- INTERVAL is used in ROW DELETION POLICY and date arithmetic
INTERVAL 30 DAY
INTERVAL 1 YEAR
INTERVAL '10:20:30' HOUR TO SECOND
SELECT DATE_ADD(CURRENT_DATE(), INTERVAL 7 DAY);
SELECT TIMESTAMP_ADD(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR);
```

---

## TYPE-016: ENUM (proto enum types)

```
-- Proto-based enums, loaded via CREATE PROTO BUNDLE
enum_type_name

-- Enum literals use integer or string representation
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#enum_type

```sql
-- After CREATE PROTO BUNDLE with a proto containing an enum
col1 `my.package.MyEnum`
-- Accessed as integers or via proto field path
```

---

## TYPE-017: Type casting and CAST

```
CAST ( expression AS target_type )
SAFE_CAST ( expression AS target_type )
-- SAFE_CAST returns NULL instead of error on cast failure
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#casting

```sql
SELECT CAST('123' AS INT64);
SELECT CAST(3.14 AS FLOAT64);
SELECT SAFE_CAST('not_a_number' AS INT64);  -- returns NULL
SELECT CAST(TRUE AS INT64);                  -- returns 1
```

---

## TYPE-018: NULL type

```
NULL
-- NULL is compatible with any type; untyped NULL has a special NULL type
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-types#null_type

```sql
SELECT NULL;
SELECT CAST(NULL AS INT64);
```
