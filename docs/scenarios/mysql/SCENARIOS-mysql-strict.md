# MySQL Parser Strictness Alignment Scenarios

> Goal: Align MySQL parser error strictness with PG parser — reject invalid SQL that is currently accepted
> Verification: corpus verifier (mysql/parser/verify_parse_test.go) reports 0 PARSE LENIENT violations; all new invalid SQL cases in corpus are correctly rejected
> Reference sources: PG parser (pg/parser/), MySQL 8.0 Reference Manual grammar, existing MySQL corpus files
> Out of scope: ALTER TABLE strictness (separate effort)

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Infrastructure

### 1.1 Main Loop Token Consumption

- [x] Parser accepts trailing identifier as implicit alias — `SELECT 1 GARBAGE` is valid MySQL (`SELECT 1 AS GARBAGE`)
- [x] Parser rejects trailing keyword after valid statement — `SELECT 1 FROM` returns error
- [x] Parser accepts trailing identifier as implicit alias — `SELECT 1 abc` is valid MySQL (`SELECT 1 AS abc`)
- [x] Parser rejects trailing operator after valid statement — `SELECT 1 +` returns error
- [x] Parser rejects trailing parenthesis after valid statement — `SELECT 1 )` returns error
- [x] Parser accepts trailing identifier in multi-statement — `SELECT 1; SELECT 2 GARBAGE` is valid MySQL (implicit alias)
- [x] Parser still accepts multiple semicolon-separated statements — `SELECT 1; SELECT 2` passes
- [x] Parser still accepts trailing semicolon — `SELECT 1;` passes
- [x] Parser still accepts empty input — `` passes (returns empty list)
- [x] Parser still accepts single semicolon — `;` passes (returns empty list)

---

## Phase 2: Statement-Level Strictness

### 2.1 SELECT — JOIN Keyword Enforcement

- [x] `SELECT * FROM t1 INNER t2` rejected — missing JOIN after INNER
- [x] `SELECT * FROM t1 LEFT t2` rejected — missing JOIN after LEFT
- [x] `SELECT * FROM t1 RIGHT t2` rejected — missing JOIN after RIGHT
- [x] `SELECT * FROM t1 CROSS t2` rejected — missing JOIN after CROSS
- [x] `SELECT * FROM t1 NATURAL t2` rejected — missing JOIN after NATURAL
- [x] `SELECT * FROM t1 LEFT OUTER t2` rejected — missing JOIN after LEFT OUTER
- [x] `SELECT * FROM t1 RIGHT OUTER t2` rejected — missing JOIN after RIGHT OUTER
- [x] `SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id` still accepted
- [x] `SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id` still accepted
- [x] `SELECT * FROM t1 CROSS JOIN t2` still accepted
- [x] `SELECT * FROM t1 NATURAL JOIN t2` still accepted

### 2.2 SELECT — Clause Validation

- [x] `SELECT * WHERE x = 1` rejected — WHERE without FROM
- [x] `SELECT * GROUP BY x` rejected — GROUP BY without FROM
- [x] `SELECT * HAVING count(*) > 1` rejected — HAVING without FROM
- [x] `SELECT * FROM t WHERE 1=1 WHERE 2=2` rejected — duplicate WHERE clause
- [x] `SELECT * FROM t GROUP BY a GROUP BY b` rejected — duplicate GROUP BY
- [x] `SELECT * FROM t ORDER BY a ORDER BY b` rejected — duplicate ORDER BY
- [x] `SELECT * FROM t` still accepted — FROM without WHERE is valid
- [x] `SELECT 1` still accepted — bare SELECT without FROM is valid
- [x] `SELECT * FROM t WHERE x = 1` still accepted — normal SELECT with WHERE

### 2.3 DELETE — FROM Enforcement

- [x] `DELETE WHERE id = 1` rejected — missing FROM and table
- [x] `DELETE LIMIT 10` rejected — missing FROM and table
- [x] `DELETE ORDER BY id LIMIT 1` rejected — missing FROM and table
- [x] `DELETE FROM t ORDER id` rejected — missing BY after ORDER
- [x] `DELETE FROM t WHERE id = 1` still accepted — single-table delete
- [x] `DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id` still accepted — multi-table delete
- [x] `DELETE FROM t1 USING t1 JOIN t2` still accepted — multi-table USING syntax

### 2.4 INSERT / REPLACE — Value Source Enforcement

- [x] `INSERT INTO t (id, name)` rejected — missing VALUES/SELECT/SET
- [x] `INSERT INTO t` rejected — missing value source entirely
- [x] `INSERT INTO t (id) WHERE 1=1` rejected — WHERE is not a valid value source
- [x] `REPLACE INTO t (id, name)` rejected — missing VALUES/SELECT/SET
- [x] `REPLACE INTO t` rejected — missing value source entirely
- [x] `INSERT INTO t VALUES (1) ON DUPLICATE` rejected — incomplete ON DUPLICATE KEY UPDATE chain
- [x] `INSERT INTO t VALUES (1) ON DUPLICATE KEY` rejected — incomplete ON DUPLICATE KEY UPDATE chain
- [x] `INSERT INTO t VALUES (1, 'a')` still accepted
- [x] `INSERT INTO t SET id = 1` still accepted
- [x] `INSERT INTO t SELECT * FROM t2` still accepted
- [x] `INSERT INTO t (id) VALUES (1)` still accepted
- [x] `REPLACE INTO t VALUES (1, 'a')` still accepted
- [x] `INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1` still accepted

### 2.5 UPDATE — SET Enforcement

- [x] `UPDATE t WHERE id = 1` rejected — missing SET clause
- [x] `UPDATE t LIMIT 10` rejected — missing SET clause
- [x] `UPDATE t ORDER BY id` rejected — missing SET clause
- [x] `UPDATE t SET x = 1 ORDER id` rejected — missing BY after ORDER
- [x] `UPDATE t SET x = 1 WHERE id = 1` still accepted
- [x] `UPDATE t SET x = 1` still accepted — UPDATE without WHERE is valid

### 2.6 DDL — IF NOT EXISTS / IF EXISTS Keyword Chains

- [x] `CREATE TABLE IF t (id INT)` rejected — missing NOT EXISTS after IF
- [x] `CREATE TABLE IF NOT t (id INT)` rejected — missing EXISTS after NOT
- [x] `CREATE INDEX IF idx ON t (id)` rejected — missing NOT EXISTS after IF
- [x] `DROP TABLE IF t` rejected — missing EXISTS after IF
- [x] `DROP INDEX IF idx ON t` rejected — missing EXISTS after IF
- [x] `DROP VIEW IF v` rejected — missing EXISTS after IF
- [x] `DROP DATABASE IF db1` rejected — missing EXISTS after IF
- [x] `CREATE TABLE IF NOT EXISTS t (id INT)` still accepted
- [x] `CREATE INDEX IF NOT EXISTS idx ON t (id)` still accepted
- [x] `DROP TABLE IF EXISTS t` still accepted
- [x] `DROP VIEW IF EXISTS v` still accepted
- [x] `DROP DATABASE IF EXISTS db1` still accepted

### 2.7 DDL — Column Constraint Keyword Chains

- [x] `CREATE TABLE t (id INT NOT)` rejected — missing NULL after NOT
- [x] `CREATE TABLE t (id INT PRIMARY)` rejected — missing KEY after PRIMARY
- [x] `CREATE TABLE t (id INT NOT NULL)` still accepted
- [x] `CREATE TABLE t (id INT PRIMARY KEY)` still accepted
- [x] `CREATE TABLE t (id INT UNIQUE)` still accepted — KEY is optional after UNIQUE
- [x] `CREATE TABLE t (id INT UNIQUE KEY)` still accepted

### 2.8 DDL — Structural Delimiter Enforcement

- [x] `CREATE TABLE t (LIKE old_t` rejected — missing closing parenthesis after LIKE
- [x] `CREATE TABLE t (id INT` rejected — missing closing parenthesis
- [x] `CREATE TABLE t (id INT, name VARCHAR(50)` rejected — missing closing paren with multiple columns
- [x] `CREATE TABLE t (id INT)` still accepted
- [x] `CREATE TABLE t (LIKE old_t)` still accepted

---

## Phase 3: Expression Strictness

### 3.1 Parenthesis Balance

- [x] `SELECT (1 + 2 FROM t` rejected — unclosed parenthesis in expression
- [x] `SELECT ((1 + 2) FROM t` rejected — unclosed nested parenthesis
- [x] `SELECT * FROM t WHERE (a = 1 AND b = 2` rejected — unclosed parenthesis in WHERE
- [x] `SELECT * FROM (SELECT 1` rejected — unclosed subquery parenthesis
- [x] `SELECT COUNT(` rejected — unclosed function call parenthesis
- [x] `SELECT (1 + 2)` still accepted — balanced parentheses
- [x] `SELECT ((1 + 2))` still accepted — nested balanced parentheses
- [x] `SELECT * FROM (SELECT 1) AS sub` still accepted — subquery with closing paren

### 3.2 Binary Operator Right-Operand

- [x] `SELECT 1 +` rejected — missing right operand for +
- [x] `SELECT 1 *` rejected — missing right operand for *
- [x] `SELECT * FROM t WHERE a =` rejected — missing right operand for =
- [x] `SELECT * FROM t WHERE a AND` rejected — missing right operand for AND
- [x] `SELECT * FROM t WHERE a OR` rejected — missing right operand for OR
- [x] `SELECT * FROM t WHERE NOT` rejected — missing operand for NOT
- [x] `SELECT 1 + 2` still accepted — complete binary expression
- [x] `SELECT * FROM t WHERE a = 1 AND b = 2` still accepted — complete compound condition

### 3.3 IN / BETWEEN / LIKE Expression Completeness

- [x] `SELECT * FROM t WHERE id IN` rejected — missing parenthesized list after IN
- [x] `SELECT * FROM t WHERE id IN (` rejected — unclosed IN list
- [x] `SELECT * FROM t WHERE id BETWEEN 1` rejected — missing AND in BETWEEN
- [x] `SELECT * FROM t WHERE id BETWEEN` rejected — missing range in BETWEEN
- [x] `SELECT * FROM t WHERE name LIKE` rejected — missing pattern after LIKE
- [x] `SELECT * FROM t WHERE name REGEXP` rejected — missing pattern after REGEXP
- [x] `SELECT * FROM t WHERE id IN (1, 2, 3)` still accepted
- [x] `SELECT * FROM t WHERE id BETWEEN 1 AND 10` still accepted
- [x] `SELECT * FROM t WHERE id NOT IN (1, 2)` still accepted
- [x] `SELECT * FROM t WHERE id NOT BETWEEN 1 AND 10` still accepted
- [x] `SELECT * FROM t WHERE name LIKE '%foo%'` still accepted

### 3.4 CASE Expression Completeness

- [x] `SELECT CASE WHEN 1 FROM t` rejected — missing THEN and END
- [x] `SELECT CASE WHEN 1 THEN 'a'` rejected — missing END
- [x] `SELECT CASE END` rejected — missing WHEN clause
- [x] `SELECT CASE WHEN 1 THEN 'a' END` still accepted
- [x] `SELECT CASE x WHEN 1 THEN 'a' ELSE 'b' END` still accepted
