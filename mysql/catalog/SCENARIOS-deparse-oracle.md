# MySQL Deparser Oracle Verification Scenarios

> Goal: Comprehensive end-to-end oracle testing of the MySQL deparser against real MySQL 8.0 SHOW CREATE VIEW output
> Verification: Create views on both MySQL 8.0 (testcontainers) and omni catalog, compare SELECT body from SHOW CREATE VIEW
> Reference: Real MySQL 8.0 output is authoritative
> Notes: (1) Oracle mismatches may require deparser fixes — tracked inline. (2) INTERSECT/EXCEPT require MySQL 8.0.31+. (3) Recursive CTEs and correlated subqueries may hit parser limitations — mark [~] if so.

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial

---

## Phase 1: Operators & Literals

### 1.1 Arithmetic & Comparison Operators

```
[x] `SELECT a + b FROM t` — addition with spacing and parens
[x] `SELECT a - b FROM t` — subtraction
[x] `SELECT a * b FROM t` — multiplication
[x] `SELECT a / b FROM t` — division
[x] `SELECT a DIV b FROM t` — integer division (DIV uppercase)
[x] `SELECT a MOD b FROM t` — MOD → %
[x] `SELECT a = b FROM t` — equals
[x] `SELECT a != b FROM t` — != → <>
[x] `SELECT a <> b FROM t` — not equals
[x] `SELECT a > b, a < b, a >= b, a <= b FROM t` — comparison ops
[x] `SELECT a <=> b FROM t` — null-safe equals
```

### 1.2 Logical, Bitwise & IS Operators

```
[x] `SELECT a AND b, a OR b FROM t` — logical with boolean wrapping
[x] `SELECT a XOR b FROM t` — XOR with boolean wrapping
[x] `SELECT NOT a FROM t` — NOT on non-boolean
[x] `SELECT a | b, a & b, a ^ b FROM t` — bitwise ops
[x] `SELECT a << b, a >> b FROM t` — shifts
[x] `SELECT ~a FROM t` — bitwise NOT
[x] `SELECT a IS NULL, a IS NOT NULL FROM t` — NULL tests
[x] `SELECT a IS TRUE, a IS FALSE FROM t` — boolean tests on INT
[x] `SELECT a IN (1,2,3), a NOT IN (1,2,3) FROM t` — membership
[x] `SELECT a BETWEEN 1 AND 10, a NOT BETWEEN 1 AND 10 FROM t` — range
[x] `SELECT a LIKE 'foo%', a NOT LIKE 'bar%' FROM t` — pattern matching
[x] `SELECT a LIKE 'x' ESCAPE '\\' FROM t` — LIKE with ESCAPE
```

### 1.3 Literals & Spacing Rules

```
[x] `SELECT 1, 1.5, 'hello', NULL, TRUE, FALSE FROM t` — all basic literal types
[x] `SELECT 0xFF, X'FF', 0b1010, b'1010' FROM t` — hex/bit literals
[x] `SELECT _utf8mb4'hello', _latin1'world' FROM t` — charset introducers
[x] `SELECT '' FROM t` — empty string
[x] `SELECT 'it''s' FROM t` — escaped quotes
[x] `SELECT 'back\\slash' FROM t` — escaped backslash
[x] `SELECT DATE '2024-01-01', TIME '12:00:00', TIMESTAMP '2024-01-01 12:00:00' FROM t` — temporal literals
[x] `SELECT CONCAT(a, b, c) FROM t` — no spaces after commas in function args
[x] `SELECT a IN (1, 2, 3) FROM t` — no spaces after commas in IN list
```

## Phase 2: Functions & Rewrites

### 2.1 Function Name Rewrites

```
[x] `SELECT SUBSTRING('abc', 1, 2) FROM t` — SUBSTRING → substr
[x] `SELECT CURRENT_TIMESTAMP FROM t` — CURRENT_TIMESTAMP → now()
[x] `SELECT CURRENT_TIMESTAMP() FROM t` — CURRENT_TIMESTAMP() → now()
[x] `SELECT CURRENT_DATE FROM t` — CURRENT_DATE → curdate()
[x] `SELECT CURRENT_TIME FROM t` — CURRENT_TIME → curtime()
[x] `SELECT CURRENT_USER FROM t` — CURRENT_USER → current_user()
[x] `SELECT NOW() FROM t` — now() stays lowercase
[x] `SELECT COUNT(*) FROM t` — COUNT(*) → count(0)
[x] `SELECT COUNT(DISTINCT a) FROM t` — count(distinct ...)
```

### 2.2 Regular Functions & Aggregates

```
[x] `SELECT CONCAT(a, b), UPPER(a), LOWER(a) FROM t` — string functions lowercase
[x] `SELECT IFNULL(a, 0), COALESCE(a, b, 0), NULLIF(a, 0) FROM t` — NULL handling
[x] `SELECT IF(a > 0, 'yes', 'no') FROM t` — IF function
[x] `SELECT ABS(a), GREATEST(a, b), LEAST(a, b) FROM t` — numeric functions
[x] `SELECT SUM(a), AVG(a), MAX(a), MIN(a) FROM t` — aggregates
[x] `SELECT CONCAT(UPPER(TRIM(a)), LOWER(b)) FROM t` — nested functions
[x] `SELECT COUNT(*), SUM(a), AVG(b), MAX(c) FROM t GROUP BY a` — multiple aggregates
```

### 2.3 Special Functions

```
[x] `SELECT TRIM(a) FROM t` — simple TRIM
[x] `SELECT TRIM(LEADING 'x' FROM a) FROM t` — directional TRIM
[x] `SELECT TRIM(TRAILING 'x' FROM a) FROM t`
[x] `SELECT TRIM(BOTH 'x' FROM a) FROM t`
[x] `SELECT GROUP_CONCAT(a ORDER BY a SEPARATOR ',') FROM t` — GROUP_CONCAT basic
[x] `SELECT GROUP_CONCAT(DISTINCT a ORDER BY a DESC SEPARATOR ';') FROM t` — GROUP_CONCAT full
[x] `SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' ELSE 'other' END FROM t` — simple CASE (distinct from searched CASE)
```

### 2.4 CAST, CONVERT & Operator-to-Function Rewrites

```
[x] `SELECT CAST(a AS CHAR) FROM t` — CAST adds charset
[x] `SELECT CAST(a AS CHAR(10)) FROM t` — CAST with length
[x] `SELECT CAST(a AS BINARY) FROM t` — BINARY → char charset binary
[x] `SELECT CAST(a AS SIGNED), CAST(a AS UNSIGNED) FROM t` — numeric CAST
[x] `SELECT CAST(a AS DECIMAL(10,2)) FROM t` — DECIMAL CAST
[x] `SELECT CAST(a AS DATE), CAST(a AS DATETIME), CAST(a AS JSON) FROM t` — other CAST types
[x] `SELECT CONVERT(a, CHAR) FROM t` — CONVERT → cast
[x] `SELECT CONVERT(a USING utf8mb4) FROM t` — CONVERT USING
[x] `SELECT a REGEXP 'pattern' FROM t` — REGEXP → regexp_like()
[x] `SELECT a NOT REGEXP 'pattern' FROM t` — NOT REGEXP → not(regexp_like())
[x] `SELECT a->'$.key' FROM t` — -> → json_extract()
[x] `SELECT a->>'$.key' FROM t` — ->> → json_unquote(json_extract())
```

## Phase 3: Boolean Context & Precedence

### 3.1 Boolean Context Wrapping

```
[x] `SELECT a AND b FROM t` — column refs get (0 <> ...) wrapping
[x] `SELECT (a + 1) AND b FROM t` — arithmetic expression in AND
[x] `SELECT (a > 0) AND (b + 1) FROM t` — comparison + expression mixed
[x] `SELECT (a > 0) AND (b > 0) FROM t` — comparisons NOT wrapped
[x] `SELECT ABS(a) AND b FROM t` — function result in AND
[x] `SELECT CASE WHEN a > 0 THEN 1 ELSE 0 END AND b FROM t` — CASE in AND
[x] `SELECT IF(a > 0, 1, 0) AND b FROM t` — IF in AND
[x] `SELECT (SELECT MAX(a) FROM t) AND b FROM t` — subquery in AND
[x] `SELECT 'hello' AND 1 FROM t` — string literal in AND
[x] `SELECT IFNULL(a, 0) AND b FROM t` — IFNULL in AND
[x] `SELECT COALESCE(a, b) AND 1 FROM t` — COALESCE in AND
[x] `SELECT NULLIF(a, 0) AND b FROM t` — NULLIF in AND
[x] `SELECT GREATEST(a, b) AND 1 FROM t` — GREATEST in AND
[x] `SELECT LEAST(a, b) AND 1 FROM t` — LEAST in AND
```

### 3.2 NOT Folding & No-Double-Wrapping

```
[x] `SELECT NOT (a > 0) FROM t` — NOT folds into <=
[x] `SELECT NOT (a + 1) FROM t` — NOT non-boolean → (0 = ...)
[x] `SELECT !a FROM t` — ! same as NOT
[x] `SELECT (a = b) AND (a > 0) FROM t` — comparisons NOT double-wrapped
[x] `SELECT (a IN (1,2,3)) AND (b BETWEEN 1 AND 10) FROM t` — predicates NOT wrapped
[x] `SELECT (a IS NULL) AND (b LIKE 'x%') FROM t` — IS/LIKE NOT wrapped
[x] `SELECT EXISTS(SELECT 1 FROM t WHERE a > 0) AND (b > 0) FROM t` — EXISTS NOT wrapped
```

### 3.3 Complex Precedence

```
[x] `SELECT a + b + c FROM t` — left-associative arithmetic
[x] `SELECT a * b + c FROM t` — multiplication before addition
[x] `SELECT a + b * c FROM t` — preserves higher precedence
[x] `SELECT (a + b) * c FROM t` — explicit grouping
[x] `SELECT a OR b AND c FROM t` — AND before OR with boolean wrapping
[x] `SELECT (a OR b) AND c FROM t` — explicit grouping with boolean wrapping
[x] `SELECT a > 0 AND b < 10 OR c = 5 FROM t` — mixed comparison/logical
```

## Phase 4: JOINs & FROM

### 4.1 All JOIN Types

```
[x] `SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a` — INNER JOIN
[x] `SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a` — LEFT JOIN
[x] `SELECT t1.a, t2.b FROM t1 RIGHT JOIN t2 ON t1.a = t2.a` — RIGHT → LEFT swap
[x] `SELECT t1.a, t2.b FROM t1 CROSS JOIN t2` — CROSS → join (no ON)
[x] `SELECT * FROM t1 NATURAL JOIN t2` — NATURAL expanded to ON conditions
[x] `SELECT t1.a, t2.b FROM t1 STRAIGHT_JOIN t2 ON t1.a = t2.a` — straight_join
[x] `SELECT t1.a, t2.b FROM t1 JOIN t2 USING (a)` — USING expanded to ON
[x] `SELECT t1.a, t2.b FROM t1, t2 WHERE t1.a = t2.a` — comma → explicit join
```

### 4.2 Multi-Table & Derived Tables

```
[x] `SELECT t1.a, t2.b, t3.c FROM t1 JOIN t2 ON t1.a = t2.a JOIN t3 ON t2.b = t3.b` — 3-table JOIN
[x] `SELECT t1.a FROM t1 LEFT JOIN t2 ON t1.a = t2.a LEFT JOIN t3 ON t1.a = t3.a` — chained LEFT JOINs
[x] `SELECT d.x FROM (SELECT a AS x FROM t) d` — derived table
[x] `SELECT d.x FROM (SELECT a AS x FROM t WHERE a > 0) AS d WHERE d.x < 10` — derived with WHERE
[x] `SELECT x.a FROM t AS x` — table alias (no AS in output)
[x] `SELECT x.a FROM t x` — table alias without AS (same output)
```

## Phase 5: SELECT Clauses

### 5.1 WHERE, GROUP BY, HAVING, ORDER BY Combined

```
[x] `SELECT a FROM t WHERE a > 1 GROUP BY a HAVING COUNT(*) > 1 ORDER BY a LIMIT 10` — all clauses combined
[x] `SELECT a, COUNT(*) cnt FROM t GROUP BY a HAVING COUNT(*) > 1 ORDER BY cnt DESC` — alias in ORDER BY
[x] `SELECT DISTINCT a FROM t ORDER BY a DESC` — DISTINCT + ORDER BY
[x] `SELECT a FROM t ORDER BY a, b DESC` — multi-column ORDER BY
[x] `SELECT a FROM t LIMIT 10 OFFSET 5` — LIMIT with OFFSET
[x] `SELECT a + b FROM t GROUP BY a + b` — expression in GROUP BY
```

### 5.2 Set Operations

```
[x] `SELECT a FROM t UNION SELECT b FROM t` — UNION
[x] `SELECT a FROM t UNION ALL SELECT b FROM t` — UNION ALL
[x] `SELECT a FROM t UNION SELECT b FROM t UNION SELECT c FROM t` — multiple UNION
[x] `SELECT a FROM t INTERSECT SELECT b FROM t` — INTERSECT
[x] `SELECT a FROM t EXCEPT SELECT b FROM t` — EXCEPT
[x] `SELECT a FROM t UNION SELECT b FROM t ORDER BY 1 LIMIT 5` — UNION + ORDER BY + LIMIT
```

### 5.3 Column & Alias Patterns

```
[x] `SELECT a AS col1, b col2 FROM t` — explicit alias AS vs space (both output AS)
[x] `SELECT a + b AS sum_col FROM t` — expression with explicit alias
[x] `SELECT 1 FROM t` — literal auto-alias
[x] `SELECT * FROM t` — star expansion (all columns listed)
[~] `SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a` — same-name columns from different tables (MySQL 8.0 rejects: duplicate column name in view)
```

## Phase 6: Subqueries & CTEs

### 6.1 Subquery Patterns

```
[ ] `SELECT (SELECT MAX(a) FROM t) FROM t` — scalar subquery in SELECT
[ ] `SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)` — IN subquery
[ ] `SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t WHERE a > 0)` — EXISTS subquery
[ ] `SELECT a, (SELECT COUNT(*) FROM t t2 WHERE t2.a = t1.a) FROM t t1` — correlated subquery
[ ] `SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE b IN (SELECT b FROM t WHERE c > 0))` — nested subqueries (2 levels)
```

### 6.2 CTE Patterns

```
[ ] `WITH cte AS (SELECT a FROM t) SELECT * FROM cte` — simple CTE
[ ] `WITH cte(x) AS (SELECT a FROM t) SELECT x FROM cte` — CTE with column list
[ ] `WITH RECURSIVE cte AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM cte WHERE n < 10) SELECT * FROM cte` — recursive CTE
[ ] `WITH c1 AS (SELECT a FROM t), c2 AS (SELECT b FROM t) SELECT c1.a, c2.b FROM c1, c2` — multiple CTEs
[ ] `WITH cte AS (SELECT a FROM t) SELECT * FROM cte UNION SELECT * FROM cte` — CTE used in UNION
```

## Phase 7: Window Functions

### 7.1 Window Function Patterns

```
[ ] `SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t` — basic window
[ ] `SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b) FROM t` — PARTITION BY + ORDER BY
[ ] `SELECT a, SUM(b) OVER (PARTITION BY a ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t` — ROWS frame
[ ] `SELECT a, AVG(b) OVER (ORDER BY a RANGE BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t` — RANGE frame
[ ] `SELECT a, RANK() OVER w, DENSE_RANK() OVER w FROM t WINDOW w AS (ORDER BY a)` — named window
[ ] `SELECT a, ROW_NUMBER() OVER (ORDER BY a), SUM(b) OVER (ORDER BY a) FROM t` — multiple window functions
[ ] `SELECT a, LAG(a, 1) OVER (ORDER BY a), LEAD(a, 1) OVER (ORDER BY a) FROM t` — LAG/LEAD
```

## Phase 8: Edge Cases & Advanced

### 8.1 View-of-View & Complex Structures

```
[ ] View referencing another view: `CREATE VIEW v1 AS SELECT a FROM t; CREATE VIEW v2 AS SELECT * FROM v1`
[ ] View with 10+ columns: `SELECT a, b, c, a+1, b+1, c+1, a*b, b*c, a*c, a+b+c FROM t`
[ ] View with reserved word aliases: `SELECT a AS `select`, b AS `from`, c AS `where` FROM t`
[ ] View with CASE without ELSE: `SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' END FROM t`
[ ] View with BETWEEN using column bounds: `SELECT a BETWEEN b AND c FROM t`
```

### 8.2 Expression Edge Cases & Stress Tests

```
[ ] `SELECT CAST(a + b AS CHAR) FROM t` — CAST with expression (not just column)
[ ] `SELECT a COLLATE utf8mb4_unicode_ci FROM t` — COLLATE expression
[ ] `SELECT INTERVAL 1 DAY + a FROM t` — INTERVAL arithmetic
[ ] `SELECT a SOUNDS LIKE b FROM t` — SOUNDS LIKE operator
[ ] `SELECT -a, +a FROM t` — unary operators (minus preserved, plus dropped)
[ ] Long expression auto-alias: `SELECT CASE WHEN a > 0 THEN CONCAT(a, b, c) WHEN a < 0 THEN CONCAT(c, b, a) ELSE NULL END FROM t` — Name_exp_N pattern
[ ] Reserved word as column name: `SELECT a FROM t` where column is named `select` (table has reserved word column)
[ ] `SELECT a + b, NOT (a > 0), CAST(a AS CHAR), COUNT(*), a REGEXP 'x' FROM t GROUP BY a, b HAVING COUNT(*) > 1` — multiple rewrites in one view
[ ] `SELECT a AND b, a OR b, NOT a, a XOR b, !a FROM t` — all logical operators with wrapping
[ ] `SELECT IFNULL(a, 0) AND COALESCE(b, 0) FROM t` — function results in boolean context
[ ] `SELECT (a > 0) AND (b + 1) OR (c IS NULL) FROM t` — mixed boolean/non-boolean precedence
```
