# BigQuery GoogleSQL truth1 — Lexical Structure

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical

---

## LEX-001: Unquoted identifiers

```
-- Must begin with a letter (a-z, A-Z) or underscore (_)
-- Subsequent characters: letters, numbers, underscores
-- Case-insensitive for keywords; case-sensitive for identifiers by default
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#identifiers
version: bigquery-current

```sql
-- Valid unquoted identifiers:
abc5
_5abc
dataField

-- Invalid (starts with number, not quoted):
-- 5abc  -> error
```

---

## LEX-002: Quoted (backtick) identifiers

```
`identifier`

-- Can contain any characters including spaces, symbols, reserved keywords
-- Cannot be empty
-- Same escape sequences as string literals
-- Required when identifier is a reserved keyword
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#identifiers
version: bigquery-current

```sql
SELECT `GROUP`, `column name with spaces`, `5abc` FROM `my-table`;
SELECT `FROM`.id FROM `project.dataset.table` AS `FROM`;
```

---

## LEX-003: Path expressions (table / column paths)

```
path:
  [path_expression][. ...]

path_expression:
  [first_part]/subsequent_part[ { / | : | - } subsequent_part ][...]

first_part:
  { unquoted_identifier | quoted_identifier }

subsequent_part:
  { unquoted_identifier | quoted_identifier | number }
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#path_expressions
version: bigquery-current

```sql
foo.bar
foo.bar/25
foo/bar:25
foo/bar/25-31
```

---

## LEX-004: Table name syntax (dashed project IDs)

```
-- Table name: project.dataset.table  (fully qualified)
-- Only the first part (project ID) can contain dashes (unquoted)
-- Dataset names cannot contain dashes
-- Table names can only contain dashes when in FROM/TABLE clause (unquoted first part)

my-project.mydataset.mytable     -- valid
mydataset.mytable                -- valid  
my-table                         -- valid (unquoted, in FROM/TABLE)
`myproject.mydataset.my-table`   -- valid (backtick-quoted)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#table_names
version: bigquery-current

```sql
SELECT * FROM my-project.mydataset.mytable;
SELECT * FROM `287mytable`;
SELECT * FROM my-table;
```

---

## LEX-005: String and bytes literals

```
-- Quoted string (single or double quotes):
'abc'
"abc"
'it\'s'
"it's"

-- Triple-quoted string (allows embedded newlines and quotes):
'''multi
line'''
"""also multi
line"""

-- Raw string (backslashes not treated as escape):
r"abc+"
r'f\(abc,(.*),def\)'
R'''abc+'''

-- Bytes literal:
b"abc"
B'''abc'''
b"\x00\xFF"

-- Raw bytes:
br'abc+'
rb"abc+"
RB'''abc'''

-- String concatenation (adjacent literals of same type):
'foo' 'bar'  --> 'foobar'
r'\n' '\n'   --> '\\n\n'
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#string_and_bytes_literals
version: bigquery-current

```sql
SELECT 'hello', "world", '''triple
quoted''', r"raw\string", b"bytes", br'raw bytes';
SELECT r'\n' '\n' "b" """c"d"e""";  -- string concatenation
```

---

## LEX-006: Escape sequences for string/bytes literals

```
\a    Bell
\b    Backspace
\f    Formfeed
\n    Newline
\r    Carriage Return
\t    Tab
\v    Vertical Tab
\\    Backslash
\?    Question Mark
\"    Double Quote
\'    Single Quote
\`    Backtick
\ooo  Octal escape (exactly 3 digits, 0-7)
\xhh or \Xhh   Hex escape (exactly 2 hex digits)
\uhhhh          Unicode escape (lowercase 'u', exactly 4 hex digits; strings only)
\Uhhhhhhhh      Unicode escape (uppercase 'U', exactly 8 hex digits; strings only)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#escape_sequences
version: bigquery-current

```sql
SELECT '\n', '\t', '\x41', 'A', '\U00000041';  -- newline, tab, 'A', 'A', 'A'
```

---

## LEX-007: Integer literals

```
-- Decimal:  [+-]DIGITS
-- Hexadecimal:  [+-]0x[HEX_DIGITS]

123
-123
0xABC
-0xDEF
+456
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#integer_literals
version: bigquery-current

```sql
SELECT 123, -456, 0xABC, 0xFF;
```

---

## LEX-008: NUMERIC and BIGNUMERIC literals

```
NUMERIC 'numeric_value_string'
DECIMAL 'numeric_value_string'

BIGNUMERIC 'numeric_value_string'
BIGDECIMAL 'numeric_value_string'
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#numeric_literals
version: bigquery-current

```sql
SELECT NUMERIC '0', NUMERIC '123456', NUMERIC '-3.14', NUMERIC '1.23456e05';
SELECT BIGNUMERIC '-9.876e-3';
```

---

## LEX-009: Floating-point literals

```
[+-]DIGITS.[DIGITS][e[+-]DIGITS]
[+-][DIGITS].DIGITS[e[+-]DIGITS]
DIGITSe[+-]DIGITS
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#floating_point_literals
version: bigquery-current

```sql
SELECT 123.456e-67, .1E4, 58., 4e2, -1.5, +0.5;
-- NaN and infinity via explicit cast only:
SELECT CAST('NaN' AS FLOAT64), CAST('inf' AS FLOAT64), CAST('-inf' AS FLOAT64);
```

---

## LEX-010: Array literals

```
[expr1, expr2, ...]                  -- implicit array
ARRAY[expr1, expr2, ...]             -- explicit ARRAY keyword
ARRAY<T>[expr1, expr2, ...]          -- typed array literal
ARRAY<T>[]                           -- empty typed array
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#array_literals
version: bigquery-current

```sql
SELECT [1, 2, 3];
SELECT ['x', 'y', 'xy'];
SELECT ARRAY[1, 2, 3];
SELECT ARRAY<STRING>['x', 'y', 'xy'];
SELECT ARRAY<INT64>[];
```

---

## LEX-011: Struct literals

```
-- Tuple syntax (2+ fields required):
(expr1, expr2 [, ...])

-- Typeless struct syntax:
STRUCT(expr1 [AS field_name] [, ...])

-- Typed struct syntax:
STRUCT<[field_name] field_type [, ...]>(expr1 [, ...])
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#struct_literals
version: bigquery-current

```sql
SELECT (1, 2, 3);
SELECT (1, 'abc');
SELECT STRUCT(1 AS foo, 'abc' AS bar);
SELECT STRUCT<INT64, STRING>(1, 'abc');
SELECT STRUCT<x INT64, y STRING>(1, 'hello');
```

---

## LEX-012: Date, Time, Datetime, Timestamp literals

```
DATE 'YYYY-MM-DD'
TIME 'HH:MM:SS[.F]'
DATETIME 'YYYY-MM-DD HH:MM:SS[.F]'
TIMESTAMP 'YYYY-MM-DD HH:MM:SS[.F] [TZ]'

-- String literals implicitly coerce in typed contexts
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#date_literals
version: bigquery-current

```sql
SELECT DATE '2014-09-27';
SELECT TIME '12:30:00.45';
SELECT DATETIME '2014-09-27 12:30:00.45';
SELECT TIMESTAMP '2018-10-01 12:00:00+08';
SELECT TIMESTAMP '2018-10-01 12:00:00 UTC';
-- Implicit coercion:
SELECT * FROM foo WHERE date_col = "2014-09-27";
```

---

## LEX-013: Interval literals

```
-- Single datetime part:
INTERVAL int64_expression datetime_part
  -- datetime_part: YEAR, MONTH, DAY, HOUR, MINUTE, SECOND, MILLISECOND, MICROSECOND, NANOSECOND,
  --                QUARTER, WEEK

-- Datetime part range:
INTERVAL datetime_parts_string starting_datetime_part TO ending_datetime_part
  -- e.g.:  INTERVAL '10:20:30.52' HOUR TO SECOND
  --        INTERVAL '1-2' YEAR TO MONTH
  --        INTERVAL '1 5:30' DAY TO MINUTE
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#interval_literals
version: bigquery-current

```sql
SELECT INTERVAL 5 DAY;
SELECT INTERVAL -5 DAY;
SELECT INTERVAL 1 SECOND;
SELECT INTERVAL '10:20:30.52' HOUR TO SECOND;
SELECT INTERVAL '1-2' YEAR TO MONTH;
SELECT INTERVAL '1 5:30' DAY TO MINUTE;
SELECT INTERVAL '-23-2 10 -12:30' YEAR TO MINUTE;
```

---

## LEX-014: Range literals

```
RANGE<T> '[lower_bound, upper_bound)'
  -- T: DATE, DATETIME, or TIMESTAMP
  -- lower_bound or upper_bound may be UNBOUNDED or NULL
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#range_literals
version: bigquery-current

```sql
SELECT RANGE<DATE> '[2014-01-01, 2015-01-01)';
SELECT RANGE<TIMESTAMP> '[2024-01-01 00:00:00 UTC, UNBOUNDED)';
```

---

## LEX-015: JSON literals

```
JSON 'json_formatted_data'
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#json_literals
version: bigquery-current

```sql
SELECT JSON '{"id": 10, "type": "fruit", "name": "apple", "on_menu": true}';
```

---

## LEX-016: Named query parameters

```
@parameter_name

-- parameter_name is an identifier (can start with letter or reserved keyword, quoted or unquoted)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#named_query_parameters
version: bigquery-current

```sql
SELECT * FROM Roster WHERE LastName = @myparam;
SELECT @x + @y AS result;
```

---

## LEX-017: Positional query parameters

```
?
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#positional_query_parameters
version: bigquery-current

```sql
SELECT * FROM Roster WHERE FirstName = ? AND LastName = ?;
```

---

## LEX-018: Comments

```
-- Single-line (hash):
# comment text
-- Single-line (double-dash):
-- comment text
-- Block / inline:
/* comment text */
/* multi
   line
   comment */
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#comments
version: bigquery-current

```sql
# this is a comment
SELECT book FROM library; -- inline comment
/* multi-line
   comment */
SELECT book FROM library /* inline block */ WHERE book = "Ulysses";
```

---

## LEX-019: Reserved keywords

```
ALL AND ANY ARRAY AS ASC ASSERT_ROWS_MODIFIED AT BETWEEN BY CASE CAST COLLATE
CONTAINS CREATE CROSS CUBE CURRENT DEFAULT DEFINE DESC DISTINCT ELSE END ENUM
ESCAPE EXCEPT EXCLUDE EXISTS EXTRACT FALSE FETCH FOLLOWING FOR FROM FULL
GRAPH_TABLE GROUP GROUPING GROUPS HASH HAVING IF IGNORE IN INNER INTERSECT
INTERVAL INTO IS JOIN LATERAL LEFT LIKE LIMIT LOOKUP MERGE NATURAL NEW NO NOT
NULL NULLS OF ON OR ORDER OUTER OVER PARTITION PRECEDING PROTO QUALIFY RANGE
RECURSIVE RESPECT RIGHT ROLLUP ROWS SELECT SET SOME STRUCT TABLESAMPLE THEN TO
TREAT TRUE UNBOUNDED UNION UNNEST USING WHEN WHERE WINDOW WITH WITHIN
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#reserved_keywords
version: bigquery-current

```sql
-- Reserved keywords must be quoted when used as identifiers:
SELECT `GROUP`, `FROM`, `SELECT` FROM `RANGE`;
```

---

## LEX-020: Trailing commas and terminating semicolons

```
-- Trailing comma after last column in SELECT list is allowed:
SELECT name, age, FROM table1

-- Semicolons terminate statements; required when multiple statements in one request:
SELECT 1;
SELECT 2;
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#terminating_semicolons
version: bigquery-current

```sql
SELECT name, release_date, FROM Books;
```
