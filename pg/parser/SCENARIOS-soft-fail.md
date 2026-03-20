# PG Parser Soft-Fail Error Alignment Scenarios

> Goal: Fix all locations where a keyword is consumed via advance() but the subsequent sub-parser soft-fails (nil, nil) without reporting a syntax error
> Verification: `go build ./pg/parser/ && go test ./pg/parser/ ./pg/parsertest/ ./pg/completion/ ./pg/pgregress/ -count=1`; error_test.go table-driven tests per scenario
> Reference sources: PostgreSQL 17 — truncated SQL should produce `syntax error at end of input` or `syntax error at or near "X"`

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: expr.go — a_expr Binary Operators

### 1.1 Logical & Comparison Operators

For each: after consuming the operator, `parseAExpr()` can return (nil, nil). The nil `right` operand is used directly in AST node construction.

- [x] `SELECT 1 OR` — OR consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 AND` — AND consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 <` — less-than consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 >` — greater-than consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 =` — equals consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 <=` — LESS_EQUALS consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 >=` — GREATER_EQUALS consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 <>` — NOT_EQUALS consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 ||` — Op consumed, no right operand → syntax error at end of input

### 1.2 Arithmetic Operators

- [x] `SELECT 1 +` — plus consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 -` — minus consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 *` — multiply consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 /` — divide consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 %` — modulo consumed, no right operand → syntax error at end of input
- [x] `SELECT 1 ^` — exponent consumed, no right operand → syntax error at end of input

### 1.3 IS DISTINCT FROM

- [x] `SELECT 1 IS DISTINCT FROM` — DISTINCT FROM consumed, no right expr → syntax error at end of input
- [x] `SELECT 1 IS NOT DISTINCT FROM` — NOT DISTINCT FROM consumed, no right expr → syntax error at end of input

## Phase 2: expr.go — Pattern Matching & Special Expressions

### 2.1 BETWEEN / LIKE / ILIKE / SIMILAR TO

- [x] `SELECT 1 BETWEEN` — BETWEEN consumed, no lower bound → syntax error at end of input
- [x] `SELECT 1 BETWEEN 0 AND` — lower bound parsed, AND consumed, no upper bound → syntax error at end of input
- [x] `SELECT 1 NOT BETWEEN` — NOT BETWEEN consumed, no lower bound → syntax error at end of input
- [x] `SELECT 'a' LIKE` — LIKE consumed, no pattern → syntax error at end of input
- [x] `SELECT 'a' LIKE 'b' ESCAPE` — ESCAPE consumed, no escape char → syntax error at end of input
- [x] `SELECT 'a' NOT LIKE` — NOT LIKE consumed, no pattern → syntax error at end of input
- [x] `SELECT 'a' ILIKE` — ILIKE consumed, no pattern → syntax error at end of input
- [x] `SELECT 'a' ILIKE 'b' ESCAPE` — ESCAPE consumed, no escape char → syntax error at end of input
- [x] `SELECT 'a' SIMILAR TO` — SIMILAR TO consumed, no pattern → syntax error at end of input
- [x] `SELECT 'a' SIMILAR TO 'b' ESCAPE` — ESCAPE consumed, no escape char → syntax error at end of input

### 2.2 COLLATE, TYPECAST & String Functions

- [x] `SELECT 'a' COLLATE` — COLLATE consumed, no collation name → syntax error at end of input
- [x] `SELECT 1::` — TYPECAST (::) consumed, no type name → syntax error at end of input
- [x] `SELECT now() AT TIME ZONE` — AT TIME ZONE consumed, no timezone expr → syntax error at end of input
- [x] `SELECT OVERLAY('abc' PLACING` — PLACING consumed, no replacement expr → syntax error at end of input
- [x] `SELECT OVERLAY('abc' PLACING 'x' FROM 1 FOR` — FOR consumed, no length expr → syntax error at end of input
- [x] `SELECT POSITION('a' IN` — IN consumed, no string expr → syntax error at end of input
- [x] `SELECT SUBSTRING('abc' FROM` — FROM consumed, no start expr → syntax error at end of input
- [x] `SELECT SUBSTRING('abc' FROM 1 FOR` — FOR consumed, no length expr → syntax error at end of input
- [x] `SELECT SUBSTRING('abc' SIMILAR` — SIMILAR consumed, no pattern → syntax error at end of input
- [x] `SELECT TRIM(LEADING FROM` — FROM consumed, no source string → syntax error at end of input

## Phase 3: expr.go — b_expr Operators

### 3.1 b_expr Binary Operators

These are the same operators in `parseBExprInfix()` context (restricted expression without AND/OR/NOT).

- [x] `SELECT CAST(1 + AS int)` — plus consumed in b_expr, no right operand → syntax error at or near "AS"
- [x] `SELECT CAST(1 - AS int)` — minus consumed in b_expr, no right operand → syntax error at or near "AS"
- [x] `SELECT CAST(1 * AS int)` — multiply consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 < AS int)` — less-than consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 > AS int)` — greater-than consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 = AS int)` — equals consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 <= AS int)` — LESS_EQUALS consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 >= AS int)` — GREATER_EQUALS consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 <> AS int)` — NOT_EQUALS consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 || AS int)` — Op consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1 IS DISTINCT FROM AS int)` — IS DISTINCT FROM consumed in b_expr → syntax error at or near "AS"
- [x] `SELECT CAST(1:: AS int)` — TYPECAST in b_expr, no type → syntax error at or near "AS"

## Phase 4: select.go — JOIN & GROUP BY

### 4.1 JOIN Clauses

After consuming JOIN keyword(s), `parseTableRefPrimary()` can return (nil, nil). The nil `right` is used directly in JoinExpr construction.

- [ ] `SELECT * FROM t CROSS JOIN` — CROSS JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t JOIN` — JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t INNER JOIN` — INNER JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t LEFT JOIN` — LEFT JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t RIGHT JOIN` — RIGHT JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t FULL JOIN` — FULL JOIN consumed, no right table → syntax error at end of input
- [ ] `SELECT * FROM t NATURAL JOIN` — NATURAL JOIN consumed, no right table → syntax error at end of input

### 4.2 GROUP BY, WHERE & Aggregation

- [ ] `SELECT * FROM t WHERE` — WHERE consumed, no condition → syntax error at end of input
- [ ] `SELECT 1 FROM t GROUP BY 1,` — comma consumed, no next group expr → syntax error at end of input
- [ ] `SELECT 1 FROM t GROUP BY CUBE(` — CUBE and ( consumed, no expr list → syntax error at end of input
- [ ] `SELECT 1 FROM t GROUP BY ROLLUP(` — ROLLUP and ( consumed, no expr list → syntax error at end of input
- [ ] `SELECT 1 FROM t GROUP BY GROUPING SETS(` — GROUPING SETS and ( consumed, no content → syntax error at end of input

## Phase 5: DDL — Constraints & ALTER

### 5.1 CREATE TABLE Constraints

- [ ] `CREATE TABLE t (a int DEFAULT` — DEFAULT consumed, no default expr → syntax error at end of input
- [ ] `CREATE TABLE t (a int COLLATE` — COLLATE consumed, no collation name → syntax error at end of input
- [ ] `CREATE TABLE t (a int REFERENCES` — REFERENCES consumed, no table name → syntax error at end of input
- [ ] `CREATE TABLE t (LIKE` — LIKE consumed, no source table name → syntax error at end of input
- [ ] `CREATE TABLE t (a int CHECK (` — CHECK ( consumed, no check expr → syntax error at end of input
- [ ] `CREATE TABLE t (a int COMPRESSION` — COMPRESSION consumed, no method → syntax error at end of input
- [ ] `CREATE TABLE t (a int STORAGE` — STORAGE consumed, no mode → syntax error at end of input

### 5.2 ALTER TABLE, GRANT, CREATE FUNCTION

- [ ] `ALTER TABLE t OF` — OF consumed, no type name → syntax error at end of input
- [ ] `ALTER TABLE t INHERIT` — INHERIT consumed, no parent table → syntax error at end of input
- [ ] `ALTER TABLE t ALTER COLUMN a SET DATA TYPE` — TYPE consumed, no type name → syntax error at end of input
- [ ] `ALTER TABLE t RENAME COLUMN` — COLUMN consumed, no column name → syntax error at end of input
- [ ] `GRANT SELECT(` — ( consumed in privilege column list, no column → syntax error at end of input
- [ ] `CREATE FUNCTION f() RETURNS int LANGUAGE sql RETURN` — RETURN consumed, no expr → syntax error at end of input
- [ ] `CREATE FUNCTION f(a DEFAULT` — DEFAULT consumed, no default expr → syntax error at end of input
- [ ] `CREATE FUNCTION f(a =` — = consumed, no default expr → syntax error at end of input
