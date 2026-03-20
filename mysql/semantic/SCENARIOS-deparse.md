# MySQL SQL Deparser Scenarios

> Goal: Implement a complete SQL deparser (AST → SQL text) + resolver for MySQL dialect, producing output identical to MySQL 8.0's SHOW CREATE VIEW format
> Verification: Oracle testing — create views on real MySQL 8.0 via testcontainers, compare SHOW CREATE VIEW output
> Reference sources: MySQL 8.0 server source (../mysql-server), mysql/catalog/deparse_rules_test.go empirical data

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial, `[-]` not applicable

Note: Phase 1-3 scenarios test final deparsed output (post-rewrite), verified against MySQL 8.0.
Expression-level tests use views (`CREATE VIEW v AS SELECT expr FROM t`) for oracle comparison.

---

## Phase 1: Deparse Infrastructure & Literals

Foundation: build the `Deparse(node) string` function and handle all literal types.

### 1.1 Infrastructure & Integer/Float/Null Literals

```
[x] Deparse function exists at mysql/deparse/deparse.go with entry point `Deparse(node ast.ExprNode) string`
[x] Integer literal: `1` → `1`
[x] Negative integer literal: `-5` → `-5` (note: may be UnaryExpr wrapping IntLit)
[x] Large integer literal: `9999999999` → `9999999999`
[x] Float literal: `1.5` → `1.5`
[x] Float with exponent: `1.5e10` → `1.5e10`
[x] NULL literal: `NULL` → `NULL` (uppercase)
```

### 1.2 Boolean & String Literals

```
[x] TRUE literal: `TRUE` → `true` (lowercase)
[x] FALSE literal: `FALSE` → `false` (lowercase)
[x] Simple string: `'hello'` → `'hello'`
[x] String with single quote: `'it''s'` → `'it\'s'` (MySQL switches to backslash escaping)
[x] Empty string: `''` → `''`
[x] String with backslash: `'back\\slash'` → `'back\\slash'`
[~] Charset introducer: `_utf8mb4'hello'` → `_utf8mb4'hello'` (deparse logic implemented; parser doesn't support charset introducers)
[~] Charset introducer latin1: `_latin1'world'` → `_latin1'world'` (deparse logic implemented; parser doesn't support charset introducers)
```

### 1.3 Hex, Bit & Date/Time Literals

```
[x] Hex literal 0x form: `0xFF` → `0xff` (lowercase)
[x] Hex literal X'' form: `X'FF'` → `0xff` (normalized to 0x form, lowercase)
[x] Bit literal 0b form: `0b1010` → `0x0a` (converted to hex)
[x] Bit literal b'' form: `b'1010'` → `0x0a` (converted to hex)
[~] DATE literal: `DATE '2024-01-01'` → `DATE'2024-01-01'` (no space, uppercase DATE) — parser doesn't support temporal literals
[~] TIME literal: `TIME '12:00:00'` → `TIME'12:00:00'` (no space, uppercase TIME) — parser doesn't support temporal literals
[~] TIMESTAMP literal: `TIMESTAMP '2024-01-01 12:00:00'` → `TIMESTAMP'2024-01-01 12:00:00'` — parser doesn't support temporal literals
```

## Phase 2: Operators & Expressions

Depends on Phase 1 (literals). Handles all operator types with correct formatting.
Rewrites (!=→<>, MOD→%, +a dropped) are applied here — oracle tests verify final output.

### 2.1 Arithmetic & Unary Operators

```
[x] Addition: `a + b` → `(`a` + `b`)` — spaces around operator, outer parens
[x] Subtraction: `a - b` → `(`a` - `b`)`
[x] Multiplication: `a * b` → `(`a` * `b`)`
[x] Division: `a / b` → `(`a` / `b`)`
[x] Integer division: `a DIV b` → `(`a` DIV `b`)` — DIV stays uppercase
[x] Modulo (MOD): `a MOD b` → `(`a` % `b`)` — MOD normalized to %
[x] Modulo (%): `a % b` → `(`a` % `b`)`
[x] Left-associative chaining: `a + b + c` → `((`a` + `b`) + `c`)`
[x] Unary minus: `-a` → `-(`a`)`
[x] Unary plus: `+a` → `a` (dropped entirely)
```

### 2.2 Comparison Operators

```
[x] Equal: `a = b` → `(`a` = `b`)`
[x] Not equal (<>): `a <> b` → `(`a` <> `b`)`
[x] Not equal (!=): `a != b` → `(`a` <> `b`)` — normalized to <>
[x] Greater: `a > b` → `(`a` > `b`)`
[x] Less: `a < b` → `(`a` < `b`)`
[x] Greater or equal: `a >= b` → `(`a` >= `b`)`
[x] Less or equal: `a <= b` → `(`a` <= `b`)`
[x] Null-safe equal: `a <=> b` → `(`a` <=> `b`)`
[x] SOUNDS LIKE: `a SOUNDS LIKE b` — verify MySQL output format
```

### 2.3 Bitwise Operators

```
[ ] Bitwise OR: `a | b` → `(`a` | `b`)`
[ ] Bitwise AND: `a & b` → `(`a` & `b`)`
[ ] Bitwise XOR: `a ^ b` → `(`a` ^ `b`)`
[ ] Left shift: `a << b` → `(`a` << `b`)`
[ ] Right shift: `a >> b` → `(`a` >> `b`)`
[ ] Bitwise NOT: `~a` → `~(`a`)` — tilde then parens
```

### 2.4 Precedence & Parenthesization

```
[ ] Higher precedence preserved: `a + b * c` → `(`a` + (`b` * `c`))`
[ ] Lower precedence grouping: `(a + b) * c` → `((`a` + `b`) * `c`)`
[ ] Mixed precedence: `a * b + c` → `((`a` * `b`) + `c`)`
[ ] Deeply nested: `a + b + c + a + b + c` → `(((((`a` + `b`) + `c`) + `a`) + `b`) + `c`)`
[ ] Parenthesized expression passthrough: `(a + b)` → `(`a` + `b`)` (outer parens from BinaryExpr)
```

### 2.5 Comparison Predicates

```
[ ] IN list: `a IN (1,2,3)` → `(`a` in (1,2,3))` — lowercase in, no spaces after commas
[ ] NOT IN: `a NOT IN (1,2,3)` → `(`a` not in (1,2,3))`
[ ] BETWEEN: `a BETWEEN 1 AND 10` → `(`a` between 1 and 10)` — lowercase
[ ] NOT BETWEEN: `a NOT BETWEEN 1 AND 10` → `(`a` not between 1 and 10)`
[ ] LIKE: `a LIKE 'foo%'` → `(`a` like 'foo%')`
[ ] LIKE with ESCAPE: `a LIKE 'x' ESCAPE '\\'` → `(`a` like 'x' escape '\\\\')`
[ ] IS NULL: `a IS NULL` → `(`a` is null)`
[ ] IS NOT NULL: `a IS NOT NULL` → `(`a` is not null)`
[ ] IS TRUE: `a IS TRUE` — verify format (may involve boolean wrapping)
[ ] IS FALSE: `a IS FALSE` — verify format
[ ] IS UNKNOWN: `a IS UNKNOWN` — verify MySQL output
[ ] ROW comparison: `ROW(a,b) = ROW(1,2)` — verify format (RowExpr)
```

### 2.6 CASE, CAST & CONVERT

```
[ ] Searched CASE: `CASE WHEN a > 0 THEN 'pos' ELSE 'zero' END` → `(case when (`a` > 0) then 'pos' else 'zero' end)`
[ ] Searched CASE multiple WHEN: `CASE WHEN a>0 THEN 'a' WHEN b>0 THEN 'b' ELSE 'c' END`
[ ] Simple CASE: `CASE a WHEN 1 THEN 'one' ELSE 'other' END` → `(case `a` when 1 then 'one' else 'other' end)`
[ ] CASE without ELSE: verify MySQL output
[ ] CAST to CHAR: `CAST(a AS CHAR)` → `cast(`a` as char charset utf8mb4)` — adds charset
[ ] CAST to CHAR(N): `CAST(a AS CHAR(10))` → `cast(`a` as char(10) charset utf8mb4)`
[ ] CAST to BINARY: `CAST(a AS BINARY)` → `cast(`a` as char charset binary)`
[ ] CAST to SIGNED: `CAST(a AS SIGNED)` → `cast(`a` as signed)`
[ ] CAST to UNSIGNED: `CAST(a AS UNSIGNED)` → `cast(`a` as unsigned)`
[ ] CAST to DECIMAL: `CAST(a AS DECIMAL(10,2))` → `cast(`a` as decimal(10,2))`
[ ] CAST to DATE/DATETIME/JSON: lowercase, no charset added
[ ] CONVERT USING: `CONVERT(a USING utf8mb4)` → `convert(`a` using utf8mb4)`
[ ] CONVERT type: `CONVERT(a, CHAR)` → `cast(`a` as char charset utf8mb4)` — rewritten to cast
```

### 2.7 Other Expressions

```
[ ] INTERVAL: `INTERVAL 1 DAY + a` → `(`a` + interval 1 day)` — operand order may swap
[ ] COLLATE: `a COLLATE utf8mb4_bin` → `(`a` collate utf8mb4_bin)`
```

## Phase 3: Functions

Depends on Phase 1-2. Covers all function call formatting.
Function name rewrites (SUBSTRING→substr, etc.) are applied here — oracle tests verify final output.

### 3.1 Regular Functions & Name Rewrites

```
[ ] Simple function: `CONCAT(a, b)` → `concat(`a`,`b`)` — lowercase, no space after comma
[ ] Nested functions: `CONCAT(UPPER(TRIM(a)), LOWER(b))` → `concat(upper(trim(`a`)),lower(`b`))`
[ ] IFNULL: `IFNULL(a, 0)` → `ifnull(`a`,0)`
[ ] COALESCE: `COALESCE(a, b, 0)` → `coalesce(`a`,`b`,0)`
[ ] NULLIF: `NULLIF(a, 0)` → `nullif(`a`,0)`
[ ] IF: `IF(a > 0, 'yes', 'no')` → `if((`a` > 0),'yes','no')` — condition gets parens
[ ] ABS: `ABS(a)` → `abs(`a`)`
[ ] GREATEST/LEAST: → `greatest(`a`,`b`)`, `least(`a`,`b`)`
[ ] SUBSTRING → substr: `SUBSTRING(a, 1, 3)` → `substr(`a`,1,3)`
[ ] CURRENT_TIMESTAMP → now(): `CURRENT_TIMESTAMP` → `now()`
[ ] CURRENT_TIMESTAMP() → now(): `CURRENT_TIMESTAMP()` → `now()`
[ ] CURRENT_DATE → curdate(): `CURRENT_DATE` → `curdate()`
[ ] CURRENT_TIME → curtime(): `CURRENT_TIME` → `curtime()`
[ ] CURRENT_USER → current_user(): `CURRENT_USER` → `current_user()` — adds parens
[ ] NOW() stays: `NOW()` → `now()`
```

### 3.2 Special Function Forms

```
[ ] TRIM simple: `TRIM(a)` → `trim(`a`)`
[ ] TRIM LEADING: `TRIM(LEADING 'x' FROM a)` — verify MySQL format
[ ] TRIM TRAILING: `TRIM(TRAILING 'x' FROM a)` — verify MySQL format
[ ] TRIM BOTH: `TRIM(BOTH 'x' FROM a)` — verify MySQL format
```

### 3.3 Aggregate Functions

```
[ ] COUNT(*): `COUNT(*)` → `count(0)` — * becomes 0
[ ] COUNT(col): `COUNT(a)` → `count(`a`)`
[ ] COUNT(DISTINCT col): `COUNT(DISTINCT a)` → `count(distinct `a`)`
[ ] SUM: `SUM(a)` → `sum(`a`)`
[ ] AVG: `AVG(a)` → `avg(`a`)`
[ ] MAX: `MAX(a)` → `max(`a`)`
[ ] MIN: `MIN(a)` → `min(`a`)`
```

### 3.4 GROUP_CONCAT

Prerequisite: AST may need a dedicated Separator field on FuncCallExpr or a GroupConcatExpr node.

```
[ ] Basic: `GROUP_CONCAT(a)` → `group_concat(`a` separator ',')`
[ ] With ORDER BY: `GROUP_CONCAT(a ORDER BY a)` → `group_concat(`a` order by `a` ASC separator ',')`
[ ] With SEPARATOR: `GROUP_CONCAT(a SEPARATOR ';')` → `group_concat(`a` separator ';')`
[ ] With DISTINCT + ORDER BY DESC + SEPARATOR: full combination
```

### 3.5 Window Functions

```
[ ] ROW_NUMBER: `ROW_NUMBER() OVER (ORDER BY a)` → `row_number() OVER (ORDER BY `a` )` — OVER uppercase, trailing space
[ ] SUM OVER PARTITION BY: `SUM(b) OVER (PARTITION BY a ORDER BY b)` — PARTITION BY uppercase
[ ] Frame clause: `ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW` — uppercase frame keywords
[ ] Named window: `WINDOW w AS (ORDER BY a)` — verify format
```

### 3.6 Operator-to-Function Rewrites

```
[ ] REGEXP → regexp_like(): `a REGEXP 'pattern'` → `regexp_like(`a`,'pattern')`
[ ] NOT REGEXP: `a NOT REGEXP 'p'` → `(not(regexp_like(`a`,'p')))`
[ ] -> → json_extract(): `a->'$.key'` → `json_extract(`a`,'$.key')`
[ ] ->> → json_unquote(json_extract()): `a->>'$.key'` → `json_unquote(json_extract(`a`,'$.key'))`
```

## Phase 4: AST Rewrites (No Schema Required)

Depends on Phase 2-3. These modify the AST before deparsing.
Phase 2-3 tests already verify rewrites in expression context.
This phase focuses on the standalone rewrite logic: NOT folding and boolean context wrapping.

### 4.1 NOT Folding

```
[ ] NOT (a > 0) → (a <= 0) — comparison operator inversion
[ ] NOT (a < 0) → (a >= 0)
[ ] NOT (a >= 0) → (a < 0)
[ ] NOT (a <= 0) → (a > 0)
[ ] NOT (a = 0) → (a <> 0)
[ ] NOT (a <> 0) → (a = 0)
[ ] NOT (non-boolean) → (0 = expr) — e.g., NOT(a+1) → (0 = (a+1))
[ ] NOT LIKE → not((expr like pattern)) — NOT wrapping, not folding
[ ] ! operator on column: `!a` → `(0 = `a`)` — C-style NOT same as NOT
[ ] NOT (a > 0) in view: oracle verification of folded output
```

### 4.2 Boolean Context Wrapping (isBooleanExpr)

```
[ ] isBooleanExpr() identifies comparison ops (=,<>,<,>,<=,>=,<=>) as boolean
[ ] isBooleanExpr() identifies IN/BETWEEN/LIKE/IS NULL/IS NOT NULL as boolean
[ ] isBooleanExpr() identifies AND/OR/NOT/XOR results as boolean
[ ] isBooleanExpr() identifies EXISTS as boolean
[ ] isBooleanExpr() identifies TRUE/FALSE literals as boolean
[ ] Column ref in AND: `a AND b` → `(0 <> `a`) and (0 <> `b`)`
[ ] Arithmetic in AND: `(a+1) AND b` → `(0 <> (`a` + 1)) and (0 <> `b`)`
[ ] Function in AND: `ABS(a) AND b` → `(0 <> abs(`a`)) and (0 <> `b`)`
[ ] CASE in AND: wrapped with `(0 <> (case...end))`
[ ] IF in AND: wrapped with `(0 <> if(...))`
[ ] Subquery in AND: wrapped with `(0 <> (select...))`
[ ] Literal in AND: `'hello' AND 1` → `(0 <> 'hello') and (0 <> 1)`
[ ] Comparison NOT wrapped: `(a > 0) AND (b > 0)` → no (0 <> ...) wrapping
[ ] Mixed: `(a > 0) AND (b + 1)` → `(`a` > 0) and (0 <> (`b` + 1))`
[ ] IS TRUE on non-boolean: `a IS TRUE` → `(0 <> `a`) is true`
[ ] IS FALSE on non-boolean: `a IS FALSE` → `(0 <> `a`) is false`
[ ] XOR: `a XOR b` → `(0 <> `a`) xor (0 <> `b`)`
```

## Phase 5: SELECT Statement Deparsing

Depends on Phase 1-4. Handles full SELECT clause formatting.

### 5.1 Target List & Aliases

```
[ ] Single column: `SELECT a FROM t` → `select `a` from `t``
[ ] Multiple columns: comma-separated, no space after comma
[ ] Column alias with AS: `SELECT a AS col1` → `select `a` AS `col1``
[ ] Column alias without AS: same output (always adds AS)
[ ] Expression alias: `SELECT a + b AS sum_col` → `(`a` + `b`) AS `sum_col``
[ ] Auto-alias literal: `SELECT 1` → `1 AS `1``
[ ] Auto-alias expression: `SELECT a + b` → `(`a` + `b`) AS `a + b``
[ ] Auto-alias empty/complex: long expressions → `Name_exp_1` pattern
```

### 5.2 FROM Clause

```
[ ] Single table: `FROM t` → `from `t``
[ ] Table alias with AS: `FROM t AS t1` → `from `t` `t1`` — no AS keyword for table alias
[ ] Table alias without AS: `FROM t t1` → `from `t` `t1`` — same output
[ ] Multiple tables (implicit cross join): `FROM t1, t2` → `from (`t1` join `t2`)`
[ ] Derived table: `FROM (SELECT a FROM t) d` → `from (select ...) `d`` — no AS
[ ] Derived table with AS: same output (AS stripped)
```

### 5.3 JOIN Clause

```
[ ] INNER JOIN: `t1 JOIN t2 ON t1.a = t2.a` → `(`t1` join `t2` on((`t1`.`a` = `t2`.`a`)))`
[ ] LEFT JOIN: `(`t1` left join `t2` on((...)))`
[ ] RIGHT JOIN → LEFT JOIN: `t1 RIGHT JOIN t2 ON ...` → `(`t2` left join `t1` on((...)))` — tables swapped
[ ] CROSS JOIN: `t1 CROSS JOIN t2` → `(`t1` join `t2`)` — CROSS removed
[ ] STRAIGHT_JOIN: `(`t1` straight_join `t2` on((...)))` — lowercase
[ ] NATURAL JOIN: expanded to `join ... on((col equality conditions))`
[ ] NATURAL LEFT JOIN: expanded similarly with `left join`
[ ] NATURAL RIGHT JOIN: expanded as `left join` with table swap
[ ] USING: `JOIN t2 USING (a)` → `join `t2` on((`t1`.`a` = `t2`.`a`))` — expanded to ON
```

### 5.4 WHERE, GROUP BY, HAVING, ORDER BY, LIMIT

```
[ ] WHERE simple: `WHERE a > 1` → `where (`a` > 1)`
[ ] WHERE compound: `WHERE a > 1 AND b < 10` → `where ((`a` > 1) and (`b` < 10))`
[ ] GROUP BY single: `GROUP BY a` → `group by `a``
[ ] GROUP BY multiple: `GROUP BY a, b` → `group by `a`,`b`` (no space after comma)
[ ] GROUP BY WITH ROLLUP: `GROUP BY a WITH ROLLUP` → `group by `a` with rollup`
[ ] HAVING: `HAVING COUNT(*) > 1` → `having (count(0) > 1)`
[ ] ORDER BY ASC (default): `ORDER BY a` → `order by `a``
[ ] ORDER BY DESC: `ORDER BY a DESC` → `order by `a` desc`
[ ] ORDER BY multiple: `ORDER BY a, b DESC` → `order by `a`,`b` desc`
[ ] LIMIT: `LIMIT 10` → `limit 10`
[ ] LIMIT with OFFSET: `LIMIT 10 OFFSET 5` → `limit 5,10` (MySQL comma syntax)
[ ] DISTINCT: `SELECT DISTINCT a` → `select distinct `a``
```

### 5.5 Set Operations

```
[ ] UNION: `SELECT a FROM t UNION SELECT b FROM t` → `select ... union select ...`
[ ] UNION ALL: → `... union all ...`
[ ] Multiple UNION: three SELECTs chained flat (not nested parens)
[ ] INTERSECT: `... intersect ...` (MySQL 8.0.31+)
[ ] EXCEPT: `... except ...` (MySQL 8.0.31+)
```

### 5.6 Subqueries

```
[ ] Scalar subquery: `(SELECT MAX(a) FROM t)` → `(select max(`a`) from `t`)`
[ ] IN subquery: `a IN (SELECT a FROM t)` → `... in (select ...)`
[ ] EXISTS: `EXISTS (SELECT 1 FROM t)` → `exists(select 1 from `t`)` — no space after exists
```

### 5.7 CTE (WITH clause)

```
[ ] Simple CTE: `WITH cte AS (...) SELECT ...` → `with `cte` as (...) select ...`
[ ] CTE with column list: `WITH cte(x) AS (...)` → `with `cte` (`x`) as (...)`
[ ] RECURSIVE CTE: `WITH RECURSIVE cte AS (...)` → `with recursive `cte` as (...)`
[ ] Multiple CTEs: `WITH c1 AS (...), c2 AS (...) SELECT ...`
```

### 5.8 FOR UPDATE / FOR SHARE

```
[ ] FOR UPDATE: → `for update` (lowercase)
[ ] FOR SHARE: → `for share`
[ ] LOCK IN SHARE MODE: → `lock in share mode`
[ ] FOR UPDATE OF table: → `for update of `t``
[ ] FOR UPDATE NOWAIT: → `for update nowait`
[ ] FOR UPDATE SKIP LOCKED: → `for update skip locked`
```

## Phase 6: Schema-Aware Resolver

Depends on Phase 4-5. Requires Catalog access.

### 6.1 Column Name Qualification

```
[ ] Resolver interface: accepts *Catalog + *SelectStmt, returns resolved *SelectStmt
[ ] Single table: `SELECT a FROM t` → `select `t`.`a` AS `a` from `t``
[ ] Multiple columns: all columns get table prefix
[ ] Qualified column preserved: `SELECT t.a FROM t` → `select `t`.`a` AS `a``
[ ] Ambiguous column (two tables with same column name): error or correct resolution
[ ] Table alias: `SELECT a FROM t AS x` → `select `x`.`a` AS `a` from `t` `x``
[ ] Column in WHERE: `WHERE a > 0` → `where (`t`.`a` > 0)`
[ ] Column in ORDER BY: qualified
[ ] Column in GROUP BY: qualified
[ ] Column in HAVING: qualified
[ ] Column in ON condition: both sides qualified
[ ] Qualified star: `SELECT t1.*, t2.a` — t1.* expanded, t2.a qualified
```

### 6.2 SELECT * Expansion

```
[ ] Simple *: `SELECT * FROM t` → all columns listed with aliases
[ ] * with alias per column: each gets `col` AS `col``
[ ] * ordering: columns in table definition order (Column.Position)
[ ] * with WHERE: columns listed, WHERE still applies
```

### 6.3 Auto-Alias Generation

```
[ ] Column ref: `SELECT a FROM t` → `` `a` AS `a` `` — alias = column name
[ ] Qualified column: `SELECT t.a FROM t` → alias = unqualified name `a`
[ ] Expression: auto-alias from original expression text
[ ] Literal: `SELECT 1` → `` 1 AS `1` ``
[ ] String literal: `SELECT 'hello'` → `` 'hello' AS `hello` ``
[ ] NULL: `SELECT NULL` → `` NULL AS `NULL` ``
[ ] Complex expression: long expressions get `Name_exp_N` auto-alias
[ ] Explicit alias preserved: `SELECT a AS x` → AS `x` not overwritten
```

### 6.4 JOIN Normalization (Schema-Aware)

```
[ ] NATURAL JOIN: find common columns → expand to ON((...) and (...))
[ ] USING single column: `USING (a)` → `on((`t1`.`a` = `t2`.`a`))`
[ ] USING multiple columns: `USING (a, b)` → `on((...) and (...))`
[ ] RIGHT JOIN → LEFT JOIN: swap table order in AST
[ ] CROSS JOIN → JOIN: remove CROSS keyword
[ ] Implicit cross join → explicit: `FROM t1, t2` → `FROM (`t1` join `t2`)`
```

### 6.5 CAST Charset from Catalog

```
[ ] CAST to CHAR uses database default charset (not hardcoded utf8mb4)
[ ] Database with latin1 charset: CAST adds `charset latin1`
```

## Phase 7: Integration — SHOW CREATE VIEW

Depends on all previous phases. End-to-end oracle testing against MySQL 8.0.

### 7.1 View Creation Pipeline

```
[ ] createView() calls resolver + deparser instead of storing raw SelectText
[ ] View.Definition contains deparsed SQL, not raw input
[ ] Preamble: CREATE ALGORITHM=UNDEFINED DEFINER=... SQL SECURITY DEFINER VIEW ... AS
```

### 7.2 Simple Views (oracle match)

```
[ ] SELECT constant: `CREATE VIEW v AS SELECT 1`
[ ] SELECT column: `CREATE VIEW v AS SELECT a FROM t`
[ ] SELECT with alias: `CREATE VIEW v AS SELECT a AS col1 FROM t`
[ ] SELECT multiple columns: `CREATE VIEW v AS SELECT a, b FROM t`
[ ] SELECT with WHERE: `CREATE VIEW v AS SELECT a FROM t WHERE a > 0`
[ ] SELECT with ORDER BY and LIMIT: `SELECT a FROM t ORDER BY a LIMIT 10`
```

### 7.3 Expression Views (oracle match)

```
[ ] Arithmetic expression view: `SELECT a + b FROM t`
[ ] Function call view: `SELECT CONCAT(a, b) FROM t`
[ ] CASE expression view: `SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t`
[ ] CAST expression view: `SELECT CAST(a AS CHAR) FROM t`
[ ] Aggregate view: `SELECT COUNT(*), SUM(a) FROM t GROUP BY a HAVING SUM(a) > 10`
```

### 7.4 Join Views (oracle match)

```
[ ] INNER JOIN view: `SELECT t1.a, t2.b FROM t1 JOIN t2 ON t1.a = t2.a`
[ ] LEFT JOIN view: `SELECT t1.a, t2.b FROM t1 LEFT JOIN t2 ON t1.a = t2.a`
[ ] Multiple table view: `SELECT t1.a FROM t1, t2 WHERE t1.a = t2.a`
[ ] Subquery in FROM view: `SELECT d.x FROM (SELECT a AS x FROM t) d`
```

### 7.5 Advanced Views (oracle match)

```
[ ] UNION view: `SELECT a FROM t UNION SELECT b FROM t`
[ ] CTE view: `WITH cte AS (SELECT a FROM t) SELECT * FROM cte`
[ ] Window function view: `SELECT a, ROW_NUMBER() OVER (ORDER BY a) FROM t`
[ ] Nested subquery view: `SELECT * FROM t WHERE a IN (SELECT a FROM t WHERE a > 0)`
[ ] Boolean expression view: `SELECT a AND b, a OR b FROM t` (with (0<>) wrapping)
[ ] Combined rewrite view: all rewrite rules in one complex view
```

### 7.6 Regression & Compatibility

```
[ ] All existing mysql/catalog tests pass (go test ./mysql/catalog/ -short)
[ ] All existing mysql/parser tests pass (go test ./mysql/parser/ -short)
[ ] SHOW CREATE VIEW for simple views matches MySQL 8.0 exactly
[ ] SHOW CREATE VIEW for complex views matches MySQL 8.0 exactly
[ ] View with explicit column aliases matches
```
