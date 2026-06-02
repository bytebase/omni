# Trino 481 — Parser-Relevant Syntax Corpus (truth1)
<!-- Source: https://trino.io/docs/current/ — Trino version 481 -->

---

## boolean-type
- **syntax:**
  ```
  BOOLEAN
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(1 AS BOOLEAN)
  true
  false
  ```
- **notes:** Literals are unquoted keywords `true` / `false`.

---

## integer-types
- **syntax:**
  ```
  TINYINT
  SMALLINT
  INTEGER
  INT
  BIGINT
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(x AS TINYINT)
  CAST(x AS SMALLINT)
  CAST(x AS INTEGER)
  CAST(x AS BIGINT)
  ```
- **notes:** `INT` is an alias for `INTEGER`. Range: TINYINT -2^7..2^7-1; SMALLINT -2^15..2^15-1; INTEGER -2^31..2^31-1; BIGINT -2^63..2^63-1.

---

## float-types
- **syntax:**
  ```
  REAL
  DOUBLE
  NUMBER
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(x AS REAL)
  CAST(x AS DOUBLE)
  REAL '10.3'
  DOUBLE '10.3'
  NUMBER 'Infinity'
  NUMBER 'NaN'
  ```
- **notes:** `REAL` is 32-bit IEEE 754; `DOUBLE` is 64-bit IEEE 754. `NUMBER` is a Trino extension supporting ≥50 decimal digits and special values `Infinity`, `-Infinity`, `NaN`. Typed literal syntax `REAL 'value'` and `DOUBLE 'value'` is supported.

---

## decimal-type
- **syntax:**
  ```
  DECIMAL(precision [, scale])
  DECIMAL(precision)
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(x AS DECIMAL(10, 2))
  DECIMAL '10.3'
  1.1
  ```
- **notes:** `precision` = total significant digits (max 38); `scale` = fractional digits (defaults to 0). Typed literal `DECIMAL 'value'` is valid syntax. Plain decimal literals like `1.1` are typed as DECIMAL.

---

## varchar-char-types
- **syntax:**
  ```
  VARCHAR
  VARCHAR(n)
  CHAR
  CHAR(n)
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(x AS VARCHAR(100))
  CAST(x AS CHAR(10))
  ```
- **notes:** `VARCHAR` without length is unbounded. `CHAR` without length defaults to 1. `CHAR` is fixed-length (right-padded with spaces).

---

## varbinary-type
- **syntax:**
  ```
  VARBINARY
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  X'65683F'
  CAST('hello' AS VARBINARY)
  ```
- **notes:** Variable-length binary data.

---

## json-type
- **syntax:**
  ```
  JSON
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST('{"a":1}' AS JSON)
  ```
- **notes:** Represents a JSON value.

---

## variant-type
- **syntax:**
  ```
  VARIANT
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(JSON '{"a": 1, "b": [true, null]}' AS VARIANT)
  ```
- **notes:** Semi-structured type supporting object, array, string, number, boolean, null, date/time values.

---

## date-type
- **syntax:**
  ```
  DATE
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DATE '2001-08-22'
  CAST('2001-08-22' AS DATE)
  ```
- **notes:** Calendar date (no time component).

---

## time-type
- **syntax:**
  ```
  TIME
  TIME(P)
  TIME WITH TIME ZONE
  TIME(P) WITH TIME ZONE
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  TIME '01:02:03.456'
  CAST('01:02:03' AS TIME(6))
  CAST('01:02:03+05:30' AS TIME WITH TIME ZONE)
  ```
- **notes:** `TIME` is an alias for `TIME(3)`. `P` is precision 0–12 (picosecond support). `WITH TIME ZONE` variant stores timezone offset.

---

## timestamp-type
- **syntax:**
  ```
  TIMESTAMP
  TIMESTAMP(P)
  TIMESTAMP WITH TIME ZONE
  TIMESTAMP(P) WITH TIME ZONE
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  TIMESTAMP '2020-06-10 15:55:23.383345'
  TIMESTAMP '2001-08-22 03:04:05.321 UTC'
  CAST('2001-08-22 03:04:05' AS TIMESTAMP(6) WITH TIME ZONE)
  ```
- **notes:** `TIMESTAMP` is alias for `TIMESTAMP(3)`. `TIMESTAMP WITH TIME ZONE` is alias for `TIMESTAMP(3) WITH TIME ZONE`. `P` ranges 0–12.

---

## interval-type
- **syntax:**
  ```
  INTERVAL YEAR TO MONTH
  INTERVAL DAY TO SECOND
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  INTERVAL '3' MONTH
  INTERVAL '2' DAY
  INTERVAL '1' YEAR
  ```
- **notes:** Two distinct interval types: year-month and day-second. The literal syntax uses a quoted value and a single unit keyword (YEAR, MONTH, DAY, HOUR, MINUTE, SECOND).

---

## array-type
- **syntax:**
  ```
  ARRAY(element_type)
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(ARRAY[1,2,3] AS ARRAY(BIGINT))
  ARRAY['foo', 'bar']
  ```
- **notes:** Type spec uses parentheses `ARRAY(T)`; the constructor uses brackets `ARRAY[...]`. See also `array-constructor`.

---

## map-type
- **syntax:**
  ```
  MAP(key_type, value_type)
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(MAP(ARRAY['a','b'], ARRAY[1,2]) AS MAP(VARCHAR, INTEGER))
  ```
- **notes:** Both key and value types are required.

---

## row-type
- **syntax:**
  ```
  ROW(field_name type [, field_name type] ...)
  ROW(type [, type] ...)
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(ROW(1, 2.0) AS ROW(x BIGINT, y DOUBLE))
  ROW(1, 2e0)
  CAST(ROW('hello', 42) AS ROW(name VARCHAR, age INTEGER))
  ```
- **notes:** Fields may be named or unnamed. Named form: `ROW(field_name type, ...)`. Unnamed form: `ROW(type, ...)`. Field access uses `.field_name` notation on named rows.

---

## ipaddress-type
- **syntax:**
  ```
  IPADDRESS
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  IPADDRESS '10.0.0.1'
  IPADDRESS '2001:db8::1'
  CAST('10.0.0.1' AS IPADDRESS)
  ```
- **notes:** Stores IPv4 and IPv6 addresses. Typed literal `IPADDRESS 'string'` is valid syntax.

---

## uuid-type
- **syntax:**
  ```
  UUID
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  UUID '12151fd2-7586-11e9-8f9e-2a86e4085a59'
  CAST('12151fd2-7586-11e9-8f9e-2a86e4085a59' AS UUID)
  ```
- **notes:** RFC 4122 format. Typed literal `UUID 'string'` is valid.

---

## sketch-types
- **syntax:**
  ```
  HyperLogLog
  P4HyperLogLog
  SetDigest
  QDigest(type)
  TDigest
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST(approx_set(x) AS HyperLogLog)
  CAST(approx_set(x) AS P4HyperLogLog)
  make_set_digest(x)
  ```
- **notes:** These are opaque internal sketch/digest types. `QDigest` is parameterized. `SetDigest` encapsulates HyperLogLog + MinHash.

---

## integer-literal
- **syntax:**
  ```
  decimal_literal        ::= [0-9][0-9_]*
  hex_literal            ::= 0x[0-9a-fA-F][0-9a-fA-F_]*
  octal_literal          ::= 0o[0-7][0-7_]*
  binary_literal         ::= 0b[01][01_]*
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  42
  123_456
  0x0A
  0o40
  0b1001
  ```
- **notes:** Underscores are allowed as digit separators but cannot be leading, trailing, or consecutive.

---

## float-literal
- **syntax:**
  ```
  floating_point_literal ::= digit+ '.' digit* exponent?
                           | '.' digit+ exponent?
                           | digit+ exponent
  exponent               ::= [eE] [+-]? digit+
  typed_float_literal    ::= REAL 'string' | DOUBLE 'string'
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  1.03e1
  10.3e0
  REAL '10.3'
  DOUBLE '10.3'
  NUMBER 'Infinity'
  ```
- **notes:** Typed literal syntax (`REAL 'value'`, `DOUBLE 'value'`) is distinct from function-call CAST.

---

## string-literal
- **syntax:**
  ```
  string_literal          ::= ''' ( non_single_quote | "''" )* '''
  unicode_literal         ::= U&''' ( unicode_char | escape_sequence )* ''' [ UESCAPE 'escape_char' ]
  unicode_escape_sequence ::= '\' hex_digit{4}
                            | '\+' hex_digit{6}
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  'Hello winter !'
  'I am big, it''s the pictures that got small!'
  U&'Hello winter \2603 !'
  U&'Hello winter #2603 !' UESCAPE '#'
  U&'\+01F600'
  ```
- **notes:** Single quotes are escaped by doubling. Unicode literals use `U&'...'` prefix; the escape character defaults to `\` but can be changed with `UESCAPE 'c'`. Six-digit codepoints use `\+` prefix.

---

## binary-literal
- **syntax:**
  ```
  binary_literal ::= X''' hex_pairs* '''
  hex_pairs      ::= [0-9a-fA-F]{2}
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  X'65683F'
  X'FFFF 0FFF  3FFF FFFF'
  ```
- **notes:** Whitespace within the hex string is ignored. Each pair of hex digits forms one byte.

---

## datetime-literals
- **syntax:**
  ```
  date_literal      ::= DATE 'yyyy-MM-dd'
  time_literal      ::= TIME 'HH:mm:ss[.SSS]'
  timestamp_literal ::= TIMESTAMP 'yyyy-MM-dd HH:mm:ss[.SSS] [timezone]'
  interval_literal  ::= INTERVAL 'value' unit
  unit              ::= YEAR | MONTH | DAY | HOUR | MINUTE | SECOND
  ```
- **source_url:** https://trino.io/docs/current/functions/datetime.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DATE '2001-08-22'
  TIME '01:02:03.456'
  TIMESTAMP '2020-06-10 15:55:23.383345'
  TIMESTAMP '2001-08-22 03:04:05.321 UTC'
  INTERVAL '3' MONTH
  INTERVAL '2' DAY
  ```
- **notes:** The `INTERVAL` literal takes a single unit keyword (not `YEAR TO MONTH` compound). Compound interval types appear in type declarations, not literals.

---

## null-literal
- **syntax:**
  ```
  NULL
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT NULL
  COALESCE(x, NULL)
  ```
- **notes:** NULL is the only value of the unknown type.

---

## arithmetic-operators
- **syntax:**
  ```
  expr + expr
  expr - expr
  expr * expr
  expr / expr
  expr % expr
  + expr          (unary plus)
  - expr          (unary minus)
  ```
- **source_url:** https://trino.io/docs/current/functions/math.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  1 + 2
  10 - 3
  4 * 5
  10 / 3
  10 % 3
  -x
  ```
- **notes:** Integer division truncates toward zero. Standard SQL precedence: `*`, `/`, `%` bind tighter than `+`, `-`.

---

## comparison-operators
- **syntax:**
  ```
  expr = expr
  expr <> expr
  expr != expr
  expr < expr
  expr > expr
  expr <= expr
  expr >= expr
  expr IS DISTINCT FROM expr
  expr IS NOT DISTINCT FROM expr
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  a = b
  a <> b
  a != b
  a < b
  a IS DISTINCT FROM b
  NULL IS NOT DISTINCT FROM NULL   -- returns true
  ```
- **notes:** `!=` is non-standard but accepted. `IS DISTINCT FROM` treats NULL as a known value; always returns TRUE or FALSE (never NULL). `IS NOT DISTINCT FROM` is the NULL-safe equality.

---

## logical-operators
- **syntax:**
  ```
  expr AND expr
  expr OR expr
  NOT expr
  ```
- **source_url:** https://trino.io/docs/current/functions/logical.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  a AND b
  a OR b
  NOT a
  a AND (b OR c)
  ```
- **notes:** Precedence (high to low): NOT > AND > OR. NULL propagation: `AND` returns FALSE if either operand is FALSE; `OR` returns TRUE if either operand is TRUE; `NOT NULL` = NULL.

---

## string-concat-operator
- **syntax:**
  ```
  string1 || string2
  array1  || array2
  array1  || element
  element || array1
  ```
- **source_url:** https://trino.io/docs/current/functions/string.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  'foo' || 'bar'
  ARRAY[1,2] || ARRAY[3]
  ARRAY[1,2] || 3
  ```
- **notes:** `||` is overloaded for both string concatenation and array concatenation/element append.

---

## array-subscript
- **syntax:**
  ```
  array_expr [ index ]
  ```
- **source_url:** https://trino.io/docs/current/functions/array.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  my_array[1]
  ARRAY[1, 1.2, 4][2]
  ```
- **notes:** Array indices are 1-based. Returns NULL if index is out of bounds.

---

## row-field-access
- **syntax:**
  ```
  row_expr . field_name
  row_expr . "quoted_field_name"
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  row_col.name
  CAST(ROW(1, 'a') AS ROW(id BIGINT, label VARCHAR)).label
  ```
- **notes:** Field access uses dot notation on ROW values with named fields.

---

## between-predicate
- **syntax:**
  ```
  value BETWEEN min AND max
  value NOT BETWEEN min AND max
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  3 BETWEEN 2 AND 6
  x NOT BETWEEN 1 AND 10
  ```
- **notes:** Equivalent to `value >= min AND value <= max`. Inclusive on both ends.

---

## in-predicate
- **syntax:**
  ```
  value [NOT] IN ( value_list )
  value [NOT] IN ( subquery )
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  name IN ('AMERICA', 'EUROPE')
  id NOT IN (1, 2, 3)
  dept_id IN (SELECT id FROM departments)
  ```
- **notes:** Works with both literal lists and subqueries returning a single column.

---

## like-predicate
- **syntax:**
  ```
  string [NOT] LIKE pattern [ESCAPE escape_char]
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  continent LIKE 'E%'
  name LIKE 'J_ne'
  path LIKE '%\%%' ESCAPE '\'
  name NOT LIKE '%test%'
  ```
- **notes:** `%` matches zero or more characters; `_` matches exactly one character. `ESCAPE` defines the escape character for literal `%` and `_`.

---

## quantified-comparison
- **syntax:**
  ```
  expression operator { ALL | ANY | SOME } ( subquery )
  operator ::= = | <> | != | < | > | <= | >=
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  x = ANY (SELECT col FROM t)
  x <> ALL (SELECT col FROM t)
  x < SOME (SELECT col FROM t)
  ```
- **notes:** `ANY` and `SOME` are synonyms. `= ANY (subquery)` is equivalent to `IN (subquery)`. `<> ALL (subquery)` is equivalent to `NOT IN (subquery)`.

---

## is-null-predicate
- **syntax:**
  ```
  expr IS NULL
  expr IS NOT NULL
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  x IS NULL
  x IS NOT NULL
  ```
- **notes:** Always returns TRUE or FALSE (never NULL).

---

## exists-predicate
- **syntax:**
  ```
  EXISTS ( subquery )
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  WHERE EXISTS (SELECT * FROM orders WHERE customer_id = c.id)
  ```
- **notes:** Returns TRUE if subquery produces at least one row.

---

## case-expression
- **syntax:**
  ```sql
  -- Simple (value) form:
  CASE expression
      WHEN value THEN result
      [ WHEN value THEN result ] ...
      [ ELSE result ]
  END

  -- Searched form:
  CASE
      WHEN condition THEN result
      [ WHEN condition THEN result ] ...
      [ ELSE result ]
  END
  ```
- **source_url:** https://trino.io/docs/current/functions/conditional.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CASE status WHEN 1 THEN 'active' WHEN 0 THEN 'inactive' ELSE 'unknown' END
  CASE WHEN score >= 90 THEN 'A' WHEN score >= 80 THEN 'B' ELSE 'C' END
  ```
- **notes:** Simple form compares expression to each WHEN value for equality. Searched form evaluates boolean conditions. ELSE is optional; without ELSE returns NULL when no condition matches.

---

## if-expression
- **syntax:**
  ```
  if(condition, true_value)
  if(condition, true_value, false_value)
  ```
- **source_url:** https://trino.io/docs/current/functions/conditional.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  if(x > 0, x, -x)
  if(name IS NULL, 'unknown', name)
  ```
- **notes:** Two-argument form returns NULL when condition is false. Three-argument form returns false_value when condition is false.

---

## coalesce-expression
- **syntax:**
  ```
  COALESCE(value1, value2 [, ...])
  ```
- **source_url:** https://trino.io/docs/current/functions/conditional.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  COALESCE(a, b, 0)
  COALESCE(name, 'unknown')
  ```
- **notes:** Returns the first non-NULL argument. Evaluates arguments left to right; stops at first non-NULL.

---

## nullif-expression
- **syntax:**
  ```
  NULLIF(value1, value2)
  ```
- **source_url:** https://trino.io/docs/current/functions/conditional.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  NULLIF(x, 0)
  NULLIF(status, 'unknown')
  ```
- **notes:** Returns NULL if value1 = value2; otherwise returns value1.

---

## try-expression
- **syntax:**
  ```
  TRY(expression)
  ```
- **source_url:** https://trino.io/docs/current/functions/conditional.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  TRY(1 / 0)
  TRY(CAST(x AS INTEGER))
  TRY(json_parse(bad_json))
  ```
- **notes:** Evaluates expression and returns NULL on certain errors: division by zero, invalid casts, numeric overflow, JSON parsing errors. Does NOT suppress all errors.

---

## cast-expression
- **syntax:**
  ```
  CAST(value AS type)
  TRY_CAST(value AS type)
  ```
- **source_url:** https://trino.io/docs/current/functions/conversion.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CAST('123' AS BIGINT)
  CAST(x AS DECIMAL(10, 2))
  CAST(ROW(1, 2.0) AS ROW(x BIGINT, y DOUBLE))
  TRY_CAST('abc' AS INTEGER)
  ```
- **notes:** `TRY_CAST` returns NULL instead of raising an error on failed conversions. `type` is any Trino data type expression.

---

## greatest-least
- **syntax:**
  ```
  GREATEST(value1, value2 [, ...])
  LEAST(value1, value2 [, ...])
  ```
- **source_url:** https://trino.io/docs/current/functions/comparison.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  GREATEST(1, 2, 3)
  LEAST(a, b, c)
  ```
- **notes:** Returns NULL if ANY argument is NULL. Return type is the common supertype of all arguments.

---

## lambda-expression
- **syntax:**
  ```
  identifier -> expression
  ( identifier [, identifier] ... ) -> expression
  ```
- **source_url:** https://trino.io/docs/current/functions/lambda.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  transform(ARRAY[1,2,3], x -> x * 2)
  filter(ARRAY[1,2,3], x -> x > 1)
  reduce(ARRAY[1,2,3], 0, (acc, x) -> acc + x, acc -> acc)
  zip_with(ARRAY[1,2], ARRAY[3,4], (x, y) -> x + y)
  x -> a * x + b
  ```
- **notes:** Single-argument lambdas do not require parentheses. Multi-argument lambdas require parentheses. Lambdas can capture outer-scope columns. Subqueries and aggregations are NOT allowed inside lambda bodies.

---

## array-constructor
- **syntax:**
  ```
  ARRAY [ element [, element] ... ]
  ARRAY []
  ```
- **source_url:** https://trino.io/docs/current/functions/array.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  ARRAY[1, 2, 4]
  ARRAY['foo', 'bar', 'bazz']
  ARRAY[]
  ```
- **notes:** The constructor uses square brackets `ARRAY[...]`. The type spec uses parentheses `ARRAY(T)`. Distinct syntactic forms.

---

## row-constructor
- **syntax:**
  ```
  ROW( value [, value] ... )
  ( value [, value] ... )
  ```
- **source_url:** https://trino.io/docs/current/language/types.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  ROW(1, 2e0)
  CAST(ROW(1, 2.0) AS ROW(x BIGINT, y DOUBLE))
  ```
- **notes:** `ROW(...)` keyword is explicit; the short form `(v1, v2)` may also produce a row in expression contexts. Commonly combined with CAST to name fields.

---

## map-constructor
- **syntax:**
  ```
  MAP(ARRAY[key ...], ARRAY[value ...])
  MAP()
  ```
- **source_url:** https://trino.io/docs/current/functions/map.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  MAP(ARRAY[1, 3], ARRAY[2, 4])
  MAP(ARRAY['foo', 'bar'], ARRAY[1, 2])
  MAP()
  ```
- **notes:** No SQL-standard map literal syntax. Maps are built via the `MAP(keys_array, values_array)` constructor function. `MAP()` returns an empty map.

---

## at-time-zone
- **syntax:**
  ```
  timestamp_expr AT TIME ZONE timezone_str
  ```
- **source_url:** https://trino.io/docs/current/functions/datetime.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  TIMESTAMP '2012-10-31 01:00 UTC' AT TIME ZONE 'America/Los_Angeles'
  current_timestamp AT TIME ZONE 'Europe/Paris'
  ```
- **notes:** Converts a TIMESTAMP WITH TIME ZONE to another timezone. `timezone_str` is a varchar timezone name or offset.

---

## current-datetime-functions
- **syntax:**
  ```
  CURRENT_DATE
  CURRENT_TIME
  CURRENT_TIMESTAMP
  LOCALTIME
  LOCALTIMESTAMP
  current_timestamp(p)
  localtimestamp(p)
  localtime(p)
  ```
- **source_url:** https://trino.io/docs/current/functions/datetime.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT CURRENT_DATE
  SELECT CURRENT_TIMESTAMP
  SELECT LOCALTIME
  SELECT LOCALTIMESTAMP
  SELECT current_timestamp(6)
  ```
- **notes:** `CURRENT_DATE`, `CURRENT_TIME`, `CURRENT_TIMESTAMP`, `LOCALTIME`, `LOCALTIMESTAMP` require NO parentheses (SQL standard). The precision-parameterized forms `current_timestamp(p)` and `localtimestamp(p)` use function-call syntax.

---

## extract-function
- **syntax:**
  ```
  EXTRACT( field FROM expr )
  field ::= YEAR | QUARTER | MONTH | WEEK | DAY | DAY_OF_MONTH
          | DAY_OF_WEEK | DOW | DAY_OF_YEAR | DOY
          | YEAR_OF_WEEK | YOW
          | HOUR | MINUTE | SECOND
          | TIMEZONE_HOUR | TIMEZONE_MINUTE
  ```
- **source_url:** https://trino.io/docs/current/functions/datetime.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  EXTRACT(YEAR FROM TIMESTAMP '2022-10-20 05:10:00')
  EXTRACT(DOW FROM DATE '2022-10-20')
  EXTRACT(HOUR FROM CURRENT_TIMESTAMP)
  ```
- **notes:** Returns BIGINT. `DOW` is an alias for `DAY_OF_WEEK` (ISO: 1=Monday). `YOW` and `YEAR_OF_WEEK` refer to ISO week-year. `TIMEZONE_HOUR` and `TIMEZONE_MINUTE` apply to `TIMESTAMP WITH TIME ZONE` and `TIME WITH TIME ZONE` values.

---

## substring-function
- **syntax:**
  ```
  substring(string, start)
  substring(string, start, length)
  substr(string, start)
  substr(string, start, length)
  ```
- **source_url:** https://trino.io/docs/current/functions/string.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  substring('hello world', 7)
  substring('hello world', 1, 5)
  substr('abc', 2, 2)
  ```
- **notes:** Positions are 1-based. Negative `start` is relative to end of string. `substr` is an alias for `substring`. The SQL-standard `SUBSTRING(x FROM start FOR length)` form is NOT documented for Trino; use the function-call form.

---

## trim-function
- **syntax:**
  ```
  trim(string)
  trim([ { LEADING | TRAILING | BOTH } ] [ chars ] FROM source)
  ltrim(string)
  rtrim(string)
  ```
- **source_url:** https://trino.io/docs/current/functions/string.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  trim('  hello  ')
  trim(LEADING FROM '  abcd')
  trim(BOTH '$' FROM '$var$')
  trim(TRAILING '!' FROM '!foo!')
  ltrim('  hello')
  rtrim('hello  ')
  ```
- **notes:** The special-form `TRIM([spec] [chars] FROM source)` supports `LEADING`, `TRAILING`, or `BOTH` (default). The `chars` argument specifies the character set to remove (defaults to whitespace). This is SQL-standard syntax, not a function call.

---

## position-function
- **syntax:**
  ```
  POSITION(substring IN string)
  ```
- **source_url:** https://trino.io/docs/current/functions/string.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  POSITION('lo' IN 'hello world')
  POSITION('x' IN 'abcdef')
  ```
- **notes:** SQL-standard special-form syntax (not `position(a, b)`). Returns 1-based position of first occurrence; returns 0 if not found.

---

## normalize-function
- **syntax:**
  ```
  NORMALIZE(string [, form])
  form ::= NFD | NFC | NFKD | NFKC
  ```
- **source_url:** https://trino.io/docs/current/functions/string.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  NORMALIZE('café', NFC)
  NORMALIZE(name, NFKD)
  NORMALIZE('hello')
  ```
- **notes:** Unicode normalization special-form. `form` is an unquoted keyword (not a string). Defaults to NFC if form is omitted.

---

## aggregate-call-syntax
- **syntax:**
  ```
  aggregate_function(expr [, expr] ...)
  aggregate_function(DISTINCT expr)
  aggregate_function(expr [ORDER BY sort_item [, ...]])
  aggregate_function(expr) FILTER (WHERE condition)
  aggregate_function(expr [ORDER BY ...]) FILTER (WHERE condition)
  ```
- **source_url:** https://trino.io/docs/current/functions/aggregate.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  count(*)
  count(DISTINCT user_id)
  array_agg(x ORDER BY y DESC)
  array_agg(x ORDER BY x, y, z)
  count(*) FILTER (WHERE petal_length_cm > 4)
  sum(amount) FILTER (WHERE status = 'paid')
  ```
- **notes:** `ORDER BY` inside aggregate controls ordering for order-sensitive aggregates (`array_agg`, `string_agg`, etc.). `FILTER (WHERE ...)` removes rows before aggregation. `DISTINCT` is supported inside aggregates.

---

## listagg-syntax
- **syntax:**
  ```
  LISTAGG( expression [, separator]
           [ON OVERFLOW { ERROR
                        | TRUNCATE [filler_string] { WITH | WITHOUT } COUNT }]
  )
  WITHIN GROUP (ORDER BY sort_item [, ...])
  [FILTER (WHERE condition)]
  ```
- **source_url:** https://trino.io/docs/current/functions/aggregate.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  LISTAGG(value, ',') WITHIN GROUP (ORDER BY value)
  LISTAGG(value, ',' ON OVERFLOW ERROR) WITHIN GROUP (ORDER BY id)
  LISTAGG(value, ',' ON OVERFLOW TRUNCATE '.....' WITH COUNT)
      WITHIN GROUP (ORDER BY value)
      FILTER (WHERE value IS NOT NULL)
  ```
- **notes:** `WITHIN GROUP (ORDER BY ...)` is required. `ON OVERFLOW ERROR` is the default (throws when output exceeds 1,048,576 bytes). `TRUNCATE` appends `filler_string` (default `'...'`) and optionally appends the overflow count. This is an SQL-standard ordered-set aggregate.

---

## window-function-syntax
- **syntax:**
  ```
  function_name([args]) OVER (window_specification)
  function_name([args]) OVER window_name

  window_specification ::=
      [ existing_window_name ]
      [ PARTITION BY partition_expr [, ...] ]
      [ ORDER BY sort_expr [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...] ]
      [ frame_clause ]

  frame_clause ::=
      { ROWS | RANGE | GROUPS }
      { frame_start | BETWEEN frame_start AND frame_end }
      [ { EXCLUDE CURRENT ROW
        | EXCLUDE GROUP
        | EXCLUDE TIES
        | EXCLUDE NO OTHERS } ]

  frame_start / frame_end ::=
      UNBOUNDED PRECEDING
    | offset PRECEDING
    | CURRENT ROW
    | offset FOLLOWING
    | UNBOUNDED FOLLOWING
  ```
- **source_url:** https://trino.io/docs/current/functions/window.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  rank() OVER (PARTITION BY dept ORDER BY salary DESC)
  sum(amount) OVER (PARTITION BY user_id ORDER BY ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
  first_value(x) OVER w
  lead(x, 2, 0) OVER (ORDER BY ts)
  nth_value(x, 3) IGNORE NULLS OVER (ORDER BY id)
  ```
- **notes:** Default frame is `RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW` when ORDER BY is present. `ROWS` uses physical row offsets; `RANGE` uses logical value ranges; `GROUPS` uses peer groups. Ranking functions (`rank()`, `dense_rank()`, `row_number()`, `ntile()`, `percent_rank()`, `cume_dist()`) do NOT accept a frame clause.

---

## null-treatment-in-window
- **syntax:**
  ```
  function_name(expr) { IGNORE NULLS | RESPECT NULLS } OVER (window_specification)
  ```
- **source_url:** https://trino.io/docs/current/functions/window.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  first_value(x) IGNORE NULLS OVER (ORDER BY id)
  last_value(x) RESPECT NULLS OVER (PARTITION BY dept ORDER BY ts)
  nth_value(x, 2) IGNORE NULLS OVER (ORDER BY id)
  lead(x) IGNORE NULLS OVER (ORDER BY id)
  lag(x) IGNORE NULLS OVER (ORDER BY id)
  ```
- **notes:** Applies to value window functions: `first_value`, `last_value`, `nth_value`, `lead`, `lag`. `RESPECT NULLS` is the default. The clause appears between the function call and the OVER keyword.

---

## window-definition-clause
- **syntax:**
  ```
  WINDOW window_name AS (window_specification) [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT rank() OVER w, sum(x) OVER w
  FROM t
  WINDOW w AS (PARTITION BY dept ORDER BY salary DESC)
  ```
- **notes:** Named window definitions appear in the `WINDOW` clause of SELECT. Can be referenced by name in the `OVER window_name` form in both SELECT and ORDER BY clauses.

---

## window-ranking-functions
- **syntax:**
  ```
  cume_dist()
  dense_rank()
  ntile(n)
  percent_rank()
  rank()
  row_number()
  ```
- **source_url:** https://trino.io/docs/current/functions/window.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  rank() OVER (PARTITION BY dept ORDER BY salary DESC)
  dense_rank() OVER (ORDER BY score DESC)
  ntile(4) OVER (ORDER BY amount DESC)
  row_number() OVER (ORDER BY id)
  ```
- **notes:** These functions require OVER clause. Frame clause is NOT allowed for ranking functions.

---

## window-value-functions
- **syntax:**
  ```
  first_value(x) [{ IGNORE NULLS | RESPECT NULLS }] OVER (...)
  last_value(x)  [{ IGNORE NULLS | RESPECT NULLS }] OVER (...)
  nth_value(x, offset) [{ IGNORE NULLS | RESPECT NULLS }] OVER (...)
  lead(x [, offset [, default_value]]) [{ IGNORE NULLS | RESPECT NULLS }] OVER (...)
  lag(x [, offset [, default_value]])  [{ IGNORE NULLS | RESPECT NULLS }] OVER (...)
  ```
- **source_url:** https://trino.io/docs/current/functions/window.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  first_value(amount) OVER (PARTITION BY user ORDER BY ts)
  last_value(amount) IGNORE NULLS OVER (ORDER BY ts)
  nth_value(amount, 3) OVER (ORDER BY ts)
  lead(amount, 1, 0) OVER (ORDER BY ts)
  lag(amount) OVER (PARTITION BY user ORDER BY ts)
  ```
- **notes:** `offset` defaults to 1 for lead/lag; `default_value` defaults to NULL.

---

## json-exists
- **syntax:**
  ```
  JSON_EXISTS(
      json_input [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]],
      json_path
      [PASSING json_argument AS parameter_name [, ...]]
      [{ TRUE | FALSE | UNKNOWN | ERROR } ON ERROR]
  )
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  JSON_EXISTS(description, 'lax $.children[*]?(@ > 10)')
  JSON_EXISTS(doc, 'strict $.name' PASSING 'Alice' AS name_param)
  JSON_EXISTS(col, '$.x' FALSE ON ERROR)
  ```
- **notes:** Returns BOOLEAN. Default ON ERROR is FALSE. Handles both input conversion errors and JSON path evaluation errors.

---

## json-query
- **syntax:**
  ```
  JSON_QUERY(
      json_input [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]],
      json_path
      [PASSING json_argument AS parameter_name [, ...]]
      [RETURNING type [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]]]
      [WITHOUT [ARRAY] WRAPPER
       | WITH [{ CONDITIONAL | UNCONDITIONAL }] [ARRAY] WRAPPER]
      [{ KEEP | OMIT } QUOTES [ON SCALAR STRING]]
      [{ ERROR | NULL | EMPTY ARRAY | EMPTY OBJECT } ON EMPTY]
      [{ ERROR | NULL | EMPTY ARRAY | EMPTY OBJECT } ON ERROR]
  )
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  JSON_QUERY(description, 'lax $.children[last]' WITH ARRAY WRAPPER)
  JSON_QUERY(doc, 'strict $.tags' RETURNING varchar)
  JSON_QUERY(col, 'lax $.x' OMIT QUOTES ON SCALAR STRING NULL ON EMPTY)
  ```
- **notes:** Default RETURNING is varchar. `WITH ARRAY WRAPPER` wraps result in a JSON array. `KEEP/OMIT QUOTES` applies to scalar string extraction. Default ON EMPTY and ON ERROR both return NULL.

---

## json-value
- **syntax:**
  ```
  JSON_VALUE(
      json_input [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]],
      json_path
      [PASSING json_argument AS parameter_name [, ...]]
      [RETURNING type]
      [{ ERROR | NULL | DEFAULT expression } ON EMPTY]
      [{ ERROR | NULL | DEFAULT expression } ON ERROR]
  )
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  JSON_VALUE(description, 'lax $.children[0]' RETURNING tinyint)
  JSON_VALUE(doc, '$.name' DEFAULT 'unknown' ON EMPTY)
  JSON_VALUE(col, 'strict $.id' ERROR ON ERROR)
  ```
- **notes:** Extracts a scalar SQL value. Default RETURNING is varchar. `DEFAULT expression` is only evaluated when the ON EMPTY or ON ERROR branch is taken. Subqueries not supported in DEFAULT.

---

## json-object
- **syntax:**
  ```
  JSON_OBJECT(
      [ key_value [, ...]
        [{ NULL ON NULL | ABSENT ON NULL }] ],
      [{ WITH UNIQUE [KEYS] | WITHOUT UNIQUE [KEYS] }]
      [RETURNING type [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]]]
  )

  key_value ::=
      'key' : value
    | KEY 'key' VALUE value
    | 'key' VALUE value
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  JSON_OBJECT('x' : true, 'y' : 1.2, 'z' : 'text')
  JSON_OBJECT(KEY 'name' VALUE 'Alice', KEY 'age' VALUE 30)
  JSON_OBJECT('x' : null, 'y' : 1 ABSENT ON NULL)
  JSON_OBJECT('a' : 1 WITH UNIQUE KEYS)
  ```
- **notes:** Default null handling for values is `NULL ON NULL` (includes JSON nulls). `WITH UNIQUE KEYS` enforces uniqueness and may throw on duplicates. Keys must be character strings; never NULL.

---

## json-array
- **syntax:**
  ```
  JSON_ARRAY(
      [ array_element [, ...]
        [{ NULL ON NULL | ABSENT ON NULL }] ],
      [RETURNING type [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]]]
  )
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  JSON_ARRAY(true, 12e-1, 'text')
  JSON_ARRAY(true, null, 1)
  JSON_ARRAY(1, 2, 3 NULL ON NULL)
  JSON_ARRAY(x, y RETURNING varbinary FORMAT JSON ENCODING UTF8)
  ```
- **notes:** Default null handling is `ABSENT ON NULL` (omits nulls). `NULL ON NULL` includes JSON nulls. Default RETURNING is varchar.

---

## json-table
- **syntax:**
  ```
  JSON_TABLE(
      json_input,
      json_path [AS path_name]
      [PASSING value AS parameter_name [, ...]]
      COLUMNS ( column_definition [, ...] )
      [PLAN (json_table_specific_plan)
       | PLAN DEFAULT (json_table_default_plan)]
      [{ ERROR | EMPTY } ON ERROR]
  )

  column_definition ::=
      column_name FOR ORDINALITY
    | column_name type
        [FORMAT JSON [ENCODING { UTF8 | UTF16 | UTF32 }]]
        [PATH json_path]
        [{ WITHOUT | WITH { CONDITIONAL | UNCONDITIONAL } } [ARRAY] WRAPPER]
        [{ KEEP | OMIT } QUOTES [ON SCALAR STRING]]
        [{ ERROR | NULL | EMPTY { [ARRAY] | OBJECT } | DEFAULT expression } ON EMPTY]
        [{ ERROR | NULL | DEFAULT expression } ON ERROR]
    | NESTED [PATH] json_path [AS path_name]
        COLUMNS ( column_definition [, ...] )
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM json_table(
      '[{"id":1,"name":"Africa"},{"id":2,"name":"Americas"}]',
      'strict $' COLUMNS (
          NESTED PATH 'strict $[*]' COLUMNS (
              id integer PATH 'strict $.id',
              name varchar PATH 'strict $.name'
          )
      )
  )

  SELECT * FROM json_table(
      doc, 'lax $.items[*]'
      COLUMNS (
          ord FOR ORDINALITY,
          item_id bigint PATH 'lax $.id' NULL ON EMPTY,
          item_name varchar PATH 'lax $.name' DEFAULT 'unknown' ON EMPTY
      )
      EMPTY ON ERROR
  )
  ```
- **notes:** `FOR ORDINALITY` adds a 1-based row counter. `NESTED [PATH]` descends into nested JSON arrays/objects. `PLAN` clause controls join strategy between NESTED paths (OUTER, INNER, CROSS, UNION). `ON ERROR` at table level: `ERROR` throws; `EMPTY` returns no rows.

---

## json-path-syntax
- **syntax:**
  ```
  path_mode ::= strict | lax

  json_path ::= path_mode path_expression

  path_expression ::=
      $                                 -- context variable
    | $ parameter_name                  -- named variable from PASSING
    | @                                 -- current item (in filter)
    | last                              -- last array index
    | path_expression . member_name     -- member accessor
    | path_expression . "member_name"   -- quoted member accessor
    | path_expression . *               -- wildcard member
    | path_expression .. member_name    -- descendant member
    | path_expression [ subscript ]     -- array subscript (0-based)
    | path_expression [ n to m ]        -- array range
    | path_expression [ * ]             -- wildcard array
    | path_expression ? ( filter_expr ) -- filter predicate
    | path_expression . method()        -- item method
    | ( path_expression )               -- grouping
    | unary_op path_expression
    | path_expression binary_op path_expression

  method ::= double() | ceiling() | floor() | abs()
           | keyvalue() | type() | size()

  unary_op ::= + | -
  binary_op ::= + | - | * | / | %

  filter_expr ::=
      filter_expr && filter_expr
    | filter_expr || filter_expr
    | ! filter_expr
    | exists ( path_expression )
    | path_expression starts with "text"
    | ( filter_expr ) is unknown
    | path_expression == path_expression
    | path_expression <> path_expression
    | path_expression != path_expression
    | path_expression < path_expression
    | path_expression > path_expression
    | path_expression <= path_expression
    | path_expression >= path_expression
  ```
- **source_url:** https://trino.io/docs/current/functions/json.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  'lax $.children[*]?(@ > 10)'
  'strict $[*].name'
  'lax $.items[0 to 2]'
  'lax $.*'
  'lax $.x.abs()'
  'lax $[*]?(@ starts with "A")'
  ```
- **notes:** `lax` mode auto-wraps/unwraps arrays and suppresses structural errors. `strict` mode requires exact schema match. Array subscripts are 0-based in JSON path (unlike SQL 1-based). `last` refers to the last array index. Descendant `..key` traverses all levels.

---

## select-synopsis
- **syntax:**
  ```
  [ WITH SESSION [ name = expression [, ...] ] ]
  [ WITH [ FUNCTION udf ] [, ...] ]
  [ WITH [ RECURSIVE ] with_query [, ...] ]
  SELECT [ ALL | DISTINCT ] select_expression [, ...]
  [ FROM from_item [, ...] ]
  [ WHERE condition ]
  [ GROUP BY [ ALL | DISTINCT ] grouping_element [, ...] ]
  [ HAVING condition ]
  [ WINDOW window_definition_list ]
  [ { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select ]
  [ ORDER BY expression [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...] ]
  [ OFFSET count [ ROW | ROWS ] ]
  [ LIMIT { count | ALL } ]
  [ FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } { ONLY | WITH TIES } ]

  select_expression ::=
      expression [ [ AS ] column_alias ]
    | row_expression . * [ AS ( column_alias [, ...] ) ]
    | relation . *
    | *

  from_item ::=
      table_name [ [ AS ] alias [ ( column_alias [, ...] ) ] ]
    | from_item join_type from_item [ ON condition | USING ( column [, ...] ) ]
    | from_item CROSS JOIN from_item
    | LATERAL ( subquery ) [ AS alias ]
    | UNNEST ( array_or_map_expr ) [ WITH ORDINALITY ] [ AS alias ( col [, ...] ) ]
    | JSON_TABLE ( ... )
    | TABLE ( table_function_invocation ) [ AS alias ]
    | MATCH_RECOGNIZE ( ... )

  join_type ::=
      [ INNER ] JOIN
    | LEFT [ OUTER ] JOIN
    | RIGHT [ OUTER ] JOIN
    | FULL [ OUTER ] JOIN
    | CROSS JOIN

  grouping_element ::=
      expression
    | GROUPING SETS ( ( expression [, ...] ) [, ...] )
    | CUBE ( expression [, ...] )
    | ROLLUP ( expression [, ...] )

  with_query ::= name [ ( col [, ...] ) ] AS ( query )
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT name, count(*) FROM orders GROUP BY name HAVING count(*) > 1
  SELECT DISTINCT dept FROM employees ORDER BY dept NULLS LAST
  SELECT * FROM t LIMIT 10 OFFSET 5
  FETCH FIRST 5 ROWS WITH TIES
  ```
- **notes:** `WITH RECURSIVE` supports recursive CTEs with UNION shape. `GROUP BY ALL` infers non-aggregated columns. `GROUP BY DISTINCT` removes duplicate grouping sets. `OFFSET` uses `ROW` or `ROWS` keyword (both accepted). `FETCH FIRST` is an alternative to LIMIT.

---

## cte-syntax
- **syntax:**
  ```
  WITH [ RECURSIVE ] cte_name [ ( col_name [, ...] ) ] AS ( query )
                   [, cte_name ...] ...
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  WITH orders_cte AS (SELECT * FROM orders WHERE status = 'open')
  SELECT * FROM orders_cte

  WITH RECURSIVE t(n) AS (
      VALUES (1)
      UNION ALL
      SELECT n + 1 FROM t WHERE n < 4
  )
  SELECT * FROM t
  ```
- **notes:** `WITH RECURSIVE` requires UNION shape (base case UNION ALL recursive step). Default max recursion depth is 10 (configurable via `max_recursion_depth` session property). Multiple CTEs are comma-separated; later ones may reference earlier.

---

## unnest-syntax
- **syntax:**
  ```
  UNNEST( array_or_map_expr [, ...] ) [ WITH ORDINALITY ] [ AS alias ( col [, ...] ) ]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM UNNEST(ARRAY[1,2,3]) AS t(x)
  SELECT * FROM UNNEST(ARRAY['a','b'], ARRAY[1,2]) WITH ORDINALITY AS t(s, n, ord)
  SELECT * FROM t CROSS JOIN UNNEST(t.tags) AS u(tag)
  ```
- **notes:** Arrays expand to single columns; maps expand to key and value columns. `WITH ORDINALITY` adds a 1-based row number column. Multiple arrays are expanded side by side with NULL padding for shorter arrays.

---

## tablesample-syntax
- **syntax:**
  ```
  table_name TABLESAMPLE { BERNOULLI | SYSTEM } ( percentage )
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM orders TABLESAMPLE BERNOULLI (10)
  SELECT * FROM orders TABLESAMPLE SYSTEM (1)
  ```
- **notes:** `BERNOULLI` does probabilistic row-level sampling (scans all blocks). `SYSTEM` does segment-level sampling (connector-dependent, faster but less precise). `percentage` is 0–100.

---

## lateral-syntax
- **syntax:**
  ```
  LATERAL ( subquery ) [ AS alias ]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM orders o, LATERAL (SELECT max(amount) FROM items WHERE order_id = o.id) AS m(max_amt)
  ```
- **notes:** `LATERAL` allows the subquery to reference columns from preceding FROM items (correlated subquery in FROM). Used to produce one result per row of the outer query.

---

## subquery-expressions
- **syntax:**
  ```
  EXISTS ( subquery )
  expr IN ( subquery )
  expr = ( subquery )          -- scalar subquery
  expr operator ANY/ALL/SOME ( subquery )
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  WHERE EXISTS (SELECT 1 FROM t WHERE t.id = outer.id)
  WHERE dept_id IN (SELECT id FROM departments)
  WHERE amount = (SELECT max(amount) FROM orders)
  ```
- **notes:** Scalar subqueries return zero or one row; returning more than one row is an error. Returns NULL if zero rows returned.

---

## datetime-arithmetic-operators
- **syntax:**
  ```
  date      + interval  → date
  date      - interval  → date
  timestamp + interval  → timestamp
  timestamp - interval  → timestamp
  interval  + interval  → interval
  interval  - interval  → interval
  ```
- **source_url:** https://trino.io/docs/current/functions/datetime.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DATE '2012-08-08' + INTERVAL '2' DAY
  TIMESTAMP '2020-01-01 00:00:00' - INTERVAL '1' HOUR
  ```
- **notes:** Date/time arithmetic uses the INTERVAL literal form. No date - date subtraction syntax; use `date_diff()` for that.

---

## match-recognize-syntax
- **syntax:**
  ```
  MATCH_RECOGNIZE (
      [ PARTITION BY column [, ...] ]
      [ ORDER BY column [, ...] ]
      [ MEASURES measure_definition [, ...] ]
      [ rows_per_match ]
      [ AFTER MATCH skip_to ]
      PATTERN ( row_pattern )
      [ SUBSET subset_definition [, ...] ]
      DEFINE variable_definition [, ...]
  )

  rows_per_match ::=
      ONE ROW PER MATCH
    | ALL ROWS PER MATCH [ SHOW EMPTY MATCHES
                          | OMIT EMPTY MATCHES
                          | WITH UNMATCHED ROWS ]

  skip_to ::=
      SKIP PAST LAST ROW
    | SKIP TO NEXT ROW
    | SKIP TO [ FIRST | LAST ] pattern_variable

  row_pattern ::= pattern_element [ | pattern_element ] ...
  pattern_element ::= pattern_term [ pattern_quantifier ]
  pattern_quantifier ::= * | + | ? | { n } | { n, m } | { n, } | { , m }

  -- anchors: ^ (start of partition), $ (end of partition)
  -- pattern grouping: ( row_pattern )
  ```
- **source_url:** https://trino.io/docs/current/sql/match-recognize.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM orders
  MATCH_RECOGNIZE (
      PARTITION BY customer_id
      ORDER BY order_date
      MEASURES A.totalprice AS starting_price, LAST(B.totalprice) AS bottom_price
      ONE ROW PER MATCH
      AFTER MATCH SKIP PAST LAST ROW
      PATTERN (A B+ C*)
      DEFINE B AS B.totalprice < PREV(B.totalprice),
             C AS C.totalprice > PREV(C.totalprice)
  )
  ```
- **notes:** Used in the FROM clause for row-pattern recognition. `MEASURES` specifies output columns. `SUBSET` declares union pattern variables. Default is `ONE ROW PER MATCH` and `SKIP PAST LAST ROW`.

---

## bitwise-functions
- **syntax:**
  ```
  bit_count(x, bits)
  bitwise_and(x, y)
  bitwise_or(x, y)
  bitwise_xor(x, y)
  bitwise_not(x)
  bitwise_left_shift(value, shift)
  bitwise_right_shift(value, shift)
  bitwise_right_shift_arithmetic(value, shift)
  ```
- **source_url:** https://trino.io/docs/current/functions/bitwise.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  bitwise_and(19, 25)
  bitwise_or(19, 25)
  bitwise_not(19)
  bitwise_left_shift(1, 2)
  bitwise_right_shift(8, 3)
  ```
- **notes:** Trino does NOT document symbolic infix operators `&`, `|`, `^`, `~` for bitwise operations. All bitwise operations are performed via named functions. All take and return BIGINT.

---

## typeof-function
- **syntax:**
  ```
  typeof(expr)
  ```
- **source_url:** https://trino.io/docs/current/functions/conversion.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  typeof(123)           -- 'integer'
  typeof('cat')         -- 'varchar(3)'
  typeof(ARRAY[1,2,3])  -- 'array(integer)'
  ```
- **notes:** Returns the type name as varchar. Useful for dynamic type inspection.

---

## format-function
- **syntax:**
  ```
  format(format_string, args ...)
  ```
- **source_url:** https://trino.io/docs/current/functions/conversion.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  format('%s%%', 123)        -- '123%'
  format('%.5f', pi())       -- '3.14159'
  format('%03d', 8)          -- '008'
  ```
- **notes:** Uses Java `Formatter` syntax for format strings. Not SQL-standard.

---

## geospatial-constructors
- **syntax:**
  ```
  ST_GeometryFromText(wkt_string)
  ST_Point(longitude, latitude)
  ST_LineString(array_of_points)
  ST_Polygon(wkt_string)
  ST_GeomFromBinary(wkb_varbinary)
  ```
- **source_url:** https://trino.io/docs/current/functions/geospatial.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  ST_GeometryFromText('POINT (0 0)')
  ST_GeometryFromText('LINESTRING (0 0, 1 1, 1 2)')
  ST_GeometryFromText('POLYGON ((0 0, 4 0, 4 4, 0 4, 0 0))')
  ST_Point(-73.9857, 40.7484)
  ```
- **notes:** Trino geospatial uses WKT (Well-Known Text) strings for geometry literals. There is no native SQL geometry literal syntax (no `GEOMETRY '...'` typed literal). The `GEOMETRY` and `SphericalGeography` types exist but are accessed only through constructor functions.

---

## order-by-syntax
- **syntax:**
  ```
  ORDER BY expression [ ASC | DESC ] [ NULLS { FIRST | LAST } ] [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  ORDER BY salary DESC NULLS LAST
  ORDER BY name ASC NULLS FIRST
  ORDER BY 1, 2 DESC
  ```
- **notes:** Default is ASC. Default null ordering is `NULLS LAST` regardless of direction. Column positions (1-based integers) are supported in ORDER BY.

---

## fetch-offset-limit
- **syntax:**
  ```
  OFFSET count [ ROW | ROWS ]
  LIMIT { count | ALL }
  FETCH { FIRST | NEXT } [ count ] { ROW | ROWS } { ONLY | WITH TIES }
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT * FROM t LIMIT 10
  SELECT * FROM t LIMIT ALL
  SELECT * FROM t OFFSET 5 ROWS
  SELECT * FROM t FETCH FIRST 10 ROWS ONLY
  SELECT * FROM t ORDER BY score DESC FETCH FIRST 3 ROWS WITH TIES
  ```
- **notes:** `LIMIT` and `FETCH FIRST` are alternatives; cannot use both. `WITH TIES` requires ORDER BY and returns additional rows that tie with the last selected row. `ROW`/`ROWS` are interchangeable keywords.

---

## grouping-sets-syntax
- **syntax:**
  ```
  GROUP BY GROUPING SETS ( ( expr [, ...] ) [, ...] )
  GROUP BY CUBE ( expr [, ...] )
  GROUP BY ROLLUP ( expr [, ...] )
  GROUP BY [ ALL | DISTINCT ] expr [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  GROUP BY GROUPING SETS ((dept, year), (dept), ())
  GROUP BY CUBE (dept, year, quarter)
  GROUP BY ROLLUP (continent, country, city)
  GROUP BY ALL dept, year
  GROUP BY DISTINCT GROUPING SETS ((a,b), (a), ())
  ```
- **notes:** `CUBE(a,b,c)` generates all 2^n grouping sets. `ROLLUP(a,b,c)` generates n+1 hierarchical subtotals. `GROUP BY ALL` (Trino extension) infers grouping from non-aggregated SELECT columns. `GROUP BY DISTINCT` removes duplicate grouping combinations.

---

## set-operations
- **syntax:**
  ```
  query { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] query
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SELECT a FROM t1 UNION ALL SELECT a FROM t2
  SELECT a FROM t1 INTERSECT SELECT a FROM t2
  SELECT a FROM t1 EXCEPT DISTINCT SELECT a FROM t2
  ```
- **notes:** Default is `DISTINCT` (deduplication). `ALL` preserves duplicates. Standard SQL precedence: INTERSECT binds tighter than UNION/EXCEPT.

---

## values-clause
- **syntax:**
  ```
  VALUES row_value [, row_value ...]
  row_value ::= ( expr [, expr] ... ) | expr
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  VALUES (1, 'a'), (2, 'b'), (3, 'c')
  VALUES 1, 2, 3
  WITH t(n) AS (VALUES (1) UNION ALL SELECT n+1 FROM t WHERE n < 5)
  SELECT * FROM (VALUES (1,'foo'), (2,'bar')) AS t(id, name)
  ```
- **notes:** `VALUES` can be used as a standalone statement, in CTEs, or as a derived table in FROM.

---

## with-session-clause
- **syntax:**
  ```
  WITH SESSION name = expression [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  WITH SESSION max_recursion_depth = 20
  SELECT ...
  ```
- **notes:** Trino extension that sets session properties for the duration of the query. Appears before the main CTE/SELECT clause.

---

## with-function-clause
- **syntax:**
  ```
  WITH [ FUNCTION function_name (...) RETURNS type ... ] [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/select.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  WITH FUNCTION double_it(x INTEGER) RETURNS INTEGER
      RETURN x * 2
  SELECT double_it(val) FROM t
  ```
- **notes:** Trino extension allowing inline UDF definitions scoped to the query. The UDF is declared inline before the main SELECT.
