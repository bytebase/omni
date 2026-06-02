# Cloud Spanner GoogleSQL truth1 — Lexical Structure

All forms extracted from official Cloud Spanner GoogleSQL lexical reference.
Primary source: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical
version: spanner-current

---

## LEX-001: Identifiers — unquoted

```
identifier:
  { letter | _ } { letter | digit | _ }*

letter: [A-Za-z]
digit:  [0-9]

-- Identifiers are case-insensitive (table and column names)
-- Max length: 128 characters
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#identifiers

```sql
SELECT SingerId, first_name, MyColumn FROM MyTable;
-- SingerId = singerid = SINGERID (case-insensitive)
```

---

## LEX-002: Identifiers — backtick-quoted

```
`identifier_with_special_chars`
`reserved_keyword_as_identifier`
`name with spaces`
`schema.table`    -- dot is part of the identifier name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#identifiers

```sql
SELECT `SELECT` FROM `my table`;
SELECT `my-column` FROM `schema.name`;
CREATE TABLE `order` (id INT64) PRIMARY KEY (id);
```

---

## LEX-003: String literals — single-quoted

```
'string value'
-- Supports escape sequences
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT 'Hello, World!';
SELECT 'It''s a test';   -- escaped single quote
SELECT 'Line1\nLine2';   -- newline escape
```

---

## LEX-004: String literals — double-quoted

```
"string value"
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT "Hello, World!";
SELECT "She said \"hi\"";
```

---

## LEX-005: String literals — triple-quoted

```
'''string value
   spanning multiple lines'''

"""string value
   spanning multiple lines"""
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT '''Line1
Line2
Line3''';

SELECT """Hello
World""";
```

---

## LEX-006: Raw string literals

```
r'raw string — no escape processing'
r"raw string"
R'raw string'
R"raw string"

-- Triple-quoted raw strings:
r'''raw triple'''
r"""raw triple"""
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT r'abc\ndef';     -- literal backslash-n, NOT a newline
SELECT r"a+b*c";
SELECT R'''no\escape''';
```

---

## LEX-007: Bytes literals

```
b'bytes value'
b"bytes value"
B'bytes value'
B"bytes value"

-- Triple-quoted bytes:
b'''bytes triple'''
b"""bytes triple"""
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT b'abc';
SELECT B'\x00\xFF';
SELECT b"""binary\x01data""";
```

---

## LEX-008: Raw bytes literals

```
rb'raw bytes'
rb"raw bytes"
br'raw bytes'
br"raw bytes"
RB'raw bytes'
BR"raw bytes"
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#string_and_bytes_literals

```sql
SELECT rb'no\escape';
SELECT br'\xFF';
```

---

## LEX-009: String/bytes escape sequences

```
\\    backslash
\'    single quote
\"    double quote
\`    backtick
\a    bell
\b    backspace
\f    form feed
\n    newline
\r    carriage return
\t    tab
\v    vertical tab
\0    null character
\ooo  octal (3 digits)
\xhh  hex (2 hex digits)
\uhhhh   Unicode (4 hex digits, strings only)
\Uhhhhhhhh Unicode (8 hex digits, strings only)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#escape_sequences

```sql
SELECT '\n\t\r';
SELECT '\x41';       -- 'A'
SELECT 'A';     -- 'A'
SELECT '\U00000041'; -- 'A'
```

---

## LEX-010: Integer literals

```
decimal_integer: [+-]?[0-9]+
hex_integer:     [+-]?0x[0-9A-Fa-f]+
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#integer_literals

```sql
SELECT 42, -100, +7;
SELECT 0x1A, 0xFF, 0XABC123;
```

---

## LEX-011: Floating-point literals

```
float_literal:
    [+-]? [0-9]+ . [0-9]* [exponent]?
  | [+-]? . [0-9]+ [exponent]?
  | [+-]? [0-9]+ exponent

exponent: [eE] [+-]? [0-9]+
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#floating_point_literals

```sql
SELECT 3.14, .5, 1.;
SELECT 1.23e10, 4.5E-3, 1e2;
SELECT -0.001, +1.5;
```

---

## LEX-012: NUMERIC literals

```
NUMERIC 'decimal_string'
NUMERIC '-12345.6789'
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#numeric_literals

```sql
SELECT NUMERIC '123.456';
SELECT NUMERIC '-9999999999.999999999';
```

---

## LEX-013: Special float values

```
'+inf'   -- positive infinity (as STRING cast to FLOAT)
'-inf'   -- negative infinity
'nan'    -- not a number
-- Case-insensitive when cast: CAST('+inf' AS FLOAT64)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#special_floating_point_values

```sql
SELECT CAST('+inf' AS FLOAT64);
SELECT CAST('-Inf' AS FLOAT32);
SELECT CAST('NaN' AS FLOAT64);
```

---

## LEX-014: Boolean literals

```
TRUE
FALSE
-- Case-insensitive
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#boolean_literals

```sql
SELECT TRUE, FALSE, true, false;
```

---

## LEX-015: NULL literal

```
NULL
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#null_literal

```sql
SELECT NULL;
SELECT COALESCE(NULL, 'default');
```

---

## LEX-016: Date, Timestamp, JSON typed literals

```
DATE 'YYYY-MM-DD'
TIMESTAMP 'YYYY-MM-DD HH:MM:SS[.FFFFFFFFF] [timezone]'
JSON '{ ... }'
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#date_literals

```sql
SELECT DATE '2024-06-01';
SELECT TIMESTAMP '2024-06-01 12:00:00 UTC';
SELECT TIMESTAMP '2024-06-01T12:00:00.000000000Z';
SELECT JSON '{"name": "Alice", "age": 30}';
```

---

## LEX-017: Query parameters

```
@parameter_name

-- Named parameters only; no positional (?) parameters
-- parameter_name follows identifier rules
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#query_parameters

```sql
SELECT * FROM Singers WHERE SingerId = @singer_id;
SELECT * FROM Albums WHERE SingerId = @id AND Title = @title;
INSERT INTO Singers (SingerId, FirstName) VALUES (@id, @name);
```

---

## LEX-018: Comments

```
-- single-line comment
/* multi-line
   comment */
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#comments

```sql
-- This is a single-line comment
SELECT 1; -- inline comment

/* This is a
   multi-line comment */
SELECT /* inline block comment */ 2;
```

---

## LEX-019: Reserved keywords

The following keywords are reserved and cannot be used as unquoted identifiers.
Backtick-quoting is required to use them as identifiers.

```
ALL, AND, ANY, ARRAY, AS, ASC, ASSERT_ROWS_MODIFIED, AT,
BETWEEN, BY,
CASE, CAST, COLLATE, CONTAINS, CREATE, CROSS, CUBE, CURRENT,
DEFAULT, DEFINE, DESC, DISTINCT,
ELSE, END, ENUM, ESCAPE, EXCEPT, EXCLUDE, EXISTS, EXTRACT,
FALSE, FETCH, FOLLOWING, FOR, FROM, FULL,
GRAPH_TABLE, GROUP, GROUPING, GROUPS,
HASH, HAVING,
IF, IGNORE, IN, INNER, INTERSECT, INTERVAL, INTO, IS,
JOIN,
LATERAL, LEFT, LIKE, LIMIT, LOOKUP,
MERGE,
NATURAL, NEW, NO, NOT, NULL, NULLS,
OF, ON, OR, ORDER, OUTER, OVER,
PARTITION, PRECEDING, PROTO,
RANGE, RECURSIVE, RESPECT, RIGHT, ROLLUP, ROWS,
SELECT, SET, SOME, STRUCT,
TABLESAMPLE, THEN, TO, TREAT, TRUE,
UNBOUNDED, UNION, UNNEST, USING,
WHEN, WHERE, WINDOW, WITH, WITHIN
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#reserved_keywords

```sql
-- Using reserved keyword as identifier requires backticks
CREATE TABLE `order` (id INT64 NOT NULL) PRIMARY KEY (id);
SELECT `select` FROM `from`;
```

---

## LEX-020: Case sensitivity

```
-- Keywords: case-insensitive
-- Identifiers (table names, column names): case-insensitive in GoogleSQL
-- String values: case-sensitive
-- Hint keys: case-insensitive
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#case_sensitivity

```sql
-- All equivalent:
SELECT singerid FROM singers;
SELECT SingerId FROM Singers;
SELECT SINGERID FROM SINGERS;
```

---

## LEX-021: Statement terminator

```
-- Semicolons separate multiple DDL/DML statements in a batch
-- Single statements do not require a trailing semicolon in the API
statement1;
statement2;
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#terminating_semicolons

```sql
CREATE TABLE A (id INT64) PRIMARY KEY (id);
CREATE TABLE B (id INT64) PRIMARY KEY (id);
```
