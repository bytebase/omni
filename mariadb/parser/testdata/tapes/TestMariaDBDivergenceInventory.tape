-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestMariaDBDivergenceInventory default#1
> exec
| CREATE TABLE IF NOT EXISTS t (id INT, val INT, score INT, status INT, name VARCHAR(100), dept VARCHAR(50), salary DECIMAL(10,2), d DATE, created_at DATETIME)
< ok 0 0
> exec
| -- MySQL SELECT corpus — quality gate scenarios
| -- Each entry: -- @name / -- @valid / SQL
< ok 0 0
> exec
| -- ============================================================
| -- Section 2.1: JOIN keyword enforcement
| -- ============================================================
< ok 0 0
> exec
| SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 INNER JOIN t2 ON t1.id = t2.id
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 LEFT JOIN t2 ON t1.id = t2.id
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t1 CROSS JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t1 CROSS JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 CROSS JOIN t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t1 NATURAL JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t1 NATURAL JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SELECT * FROM t1 NATURAL JOIN t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t1 INNER t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 LEFT t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 RIGHT t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 CROSS t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 NATURAL t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 LEFT OUTER t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| SELECT * FROM t1 RIGHT OUTER t2
< error
| 1064 42000 …near 't2' at line 1
> exec
| -- ============================================================
| -- Section 2.2: Clause validation
| -- ============================================================
< ok 0 0
> exec
| SELECT * WHERE x = 1
< error
| 1096 HY000 No tables used
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * WHERE x = 1
< error
| 1096 HY000 No tables used
> exec ^
: string
| SELECT * WHERE x = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * GROUP BY x
< error
| 1096 HY000 No tables used
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * GROUP BY x
< error
| 1096 HY000 No tables used
> exec ^
: string
| SELECT * GROUP BY x
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * HAVING count(*) > 1
< error
| 1096 HY000 No tables used
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * HAVING count(*) > 1
< error
| 1096 HY000 No tables used
> exec ^
: string
| SELECT * HAVING count(*) > 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t WHERE 1=1 WHERE 2=2
< error
| 1064 42000 …near 'WHERE 2=2' at line 1
> exec
| SELECT * FROM t GROUP BY a GROUP BY b
< error
| 1064 42000 …near 'GROUP BY b' at line 1
> exec
| SELECT * FROM t ORDER BY a ORDER BY b
< error
| 1064 42000 …near 'ORDER BY b' at line 1
> exec
| SELECT * FROM t
< ok 0 0
> exec
| SELECT 1
< ok 0 0
> exec
| SELECT * FROM t WHERE x = 1
< error
| 1054 42S22 Unknown column 'x' in 'WHERE'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE x = 1
< error
| 1054 42S22 Unknown column 'x' in 'WHERE'
> exec ^
: string
| SELECT * FROM t WHERE x = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| -- ============================================================
| -- Section 3.1: Parenthesis balance
| -- ============================================================
< ok 0 0
> exec
| SELECT (1 + 2 FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec
| SELECT ((1 + 2) FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec
| SELECT * FROM t WHERE (a = 1 AND b = 2
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM (SELECT 1
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT COUNT(
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT (1 + 2)
< ok 0 0
> exec
| SELECT ((1 + 2))
< ok 0 0
> exec
| SELECT * FROM (SELECT 1) AS sub
< ok 0 0
> exec
| -- ============================================================
| -- Section 3.2: Binary operator right-operand
| -- ============================================================
< ok 0 0
> exec
| SELECT 1 +
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 *
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE a =
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE a AND
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE a OR
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE NOT
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT 1 + 2
< ok 0 0
> exec
| SELECT * FROM t WHERE a = 1 AND b = 2
< error
| 1054 42S22 Unknown column 'a' in 'WHERE'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE a = 1 AND b = 2
< error
| 1054 42S22 Unknown column 'a' in 'WHERE'
> exec ^
: string
| SELECT * FROM t WHERE a = 1 AND b = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| -- ============================================================
| -- Section 3.3: IN/BETWEEN/LIKE expression completeness
| -- ============================================================
< ok 0 0
> exec
| SELECT * FROM t WHERE id IN
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE id IN (
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE id BETWEEN 1
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE id BETWEEN
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE name LIKE
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE name REGEXP
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT * FROM t WHERE id IN (1, 2, 3)
< ok 0 0
> exec
| SELECT * FROM t WHERE id BETWEEN 1 AND 10
< ok 0 0
> exec
| SELECT * FROM t WHERE id NOT IN (1, 2)
< ok 0 0
> exec
| SELECT * FROM t WHERE id NOT BETWEEN 1 AND 10
< ok 0 0
> exec
| SELECT * FROM t WHERE name LIKE '%foo%'
< ok 0 0
> exec
| -- ============================================================
| -- Section 3.4: CASE expression completeness
| -- ============================================================
< ok 0 0
> exec
| SELECT CASE WHEN 1 FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec
| SELECT CASE WHEN 1 THEN 'a'
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT CASE END
< error
| 1064 42000 …near '' at line 1
> exec
| SELECT CASE WHEN 1 THEN 'a' END
< ok 0 0
> exec
| SELECT CASE x WHEN 1 THEN 'a' ELSE 'b' END
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT CASE x WHEN 1 THEN 'a' ELSE 'b' END
< error
| 1054 42S22 Unknown column 'x' in 'SELECT'
> exec ^
: string
| SELECT CASE x WHEN 1 THEN 'a' ELSE 'b' END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DELETE WHERE id = 1
< error
| 1064 42000 …near 'WHERE id = 1' at line 1
> exec
| DELETE LIMIT 10
< error
| 1064 42000 …near 'LIMIT 10' at line 1
> exec
| DELETE ORDER BY id LIMIT 1
< error
| 1064 42000 …near 'ORDER BY id LIMIT 1' at line 1
> exec
| DELETE FROM t ORDER id
< error
| 1064 42000 …near 'id' at line 1
> exec
| DELETE FROM t WHERE id = 1
< ok 0 0
> exec
| DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DELETE FROM t1 USING t1 JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| DELETE FROM t1 USING t1 JOIN t2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| DELETE FROM t1 USING t1 JOIN t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| UPDATE t WHERE id = 1
< error
| 1064 42000 …near 'WHERE id = 1' at line 1
> exec
| UPDATE t LIMIT 10
< error
| 1064 42000 …near 'LIMIT 10' at line 1
> exec
| UPDATE t ORDER BY id
< error
| 1064 42000 …near 'ORDER BY id' at line 1
> exec
| UPDATE t SET x = 1 ORDER id
< error
| 1064 42000 …near 'id' at line 1
> exec
| UPDATE t SET x = 1 WHERE id = 1
< error
| 1054 42S22 Unknown column 'x' in 'SET'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| UPDATE t SET x = 1 WHERE id = 1
< error
| 1054 42S22 Unknown column 'x' in 'SET'
> exec ^
: string
| UPDATE t SET x = 1 WHERE id = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| UPDATE t SET x = 1
< error
| 1054 42S22 Unknown column 'x' in 'SET'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| UPDATE t SET x = 1
< error
| 1054 42S22 Unknown column 'x' in 'SET'
> exec ^
: string
| UPDATE t SET x = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| INSERT INTO t (id, name)
< error
| 1064 42000 …near '' at line 1
> exec
| INSERT INTO t
< error
| 1064 42000 …near '' at line 1
> exec
| INSERT INTO t (id) WHERE 1=1
< error
| 1064 42000 …near 'WHERE 1=1' at line 1
> exec
| REPLACE INTO t (id, name)
< error
| 1064 42000 …near '' at line 1
> exec
| REPLACE INTO t
< error
| 1064 42000 …near '' at line 1
> exec
| INSERT INTO t VALUES (1) ON DUPLICATE
< error
| 1064 42000 …near '' at line 1
> exec
| INSERT INTO t VALUES (1) ON DUPLICATE KEY
< error
| 1064 42000 …near '' at line 1
> exec
| INSERT INTO t VALUES (1, 'a')
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO t VALUES (1, 'a')
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec ^
: string
| INSERT INTO t VALUES (1, 'a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| INSERT INTO t SET id = 1
< ok 1 0
> exec
| INSERT INTO t SELECT * FROM t2
< error
| 1146 42S02 Table 'test.t2' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO t SELECT * FROM t2
< error
| 1146 42S02 Table 'test.t2' doesn't exist
> exec ^
: string
| INSERT INTO t SELECT * FROM t2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| INSERT INTO t (id) VALUES (1)
< ok 1 0
> exec
| REPLACE INTO t VALUES (1, 'a')
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| REPLACE INTO t VALUES (1, 'a')
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec ^
: string
| REPLACE INTO t VALUES (1, 'a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1
< error
| 1136 21S01 Column count doesn't match value count at row 1
> exec ^
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| -- MySQL DDL corpus — keyword chain enforcement and structural delimiters
| -- Covers: §2.6 IF NOT EXISTS / IF EXISTS chains, §2.7 column constraint chains, §2.8 structural delimiters
| -- Note: §2.8 cases require no parser changes — expect(')') already enforces closing parens
< ok 0 0
> exec
| -- ============================================================
| -- Section 2.6: IF NOT EXISTS / IF EXISTS Keyword Chains
| -- ============================================================
< ok 0 0
> exec
| CREATE TABLE IF t (id INT)
< error
| 1064 42000 …near 't (id INT)' at line 1
> exec
| CREATE TABLE IF NOT t (id INT)
< error
| 1064 42000 …near 't (id INT)' at line 1
> exec
| CREATE INDEX IF idx ON t (id)
< error
| 1064 42000 …near 'idx ON t (id)' at line 1
> exec
| DROP TABLE IF t
< error
| 1064 42000 …near 't' at line 1
> exec
| DROP INDEX IF idx ON t
< error
| 1064 42000 …near 'idx ON t' at line 1
> exec
| DROP VIEW IF v
< error
| 1064 42000 …near 'v' at line 1
> exec
| DROP DATABASE IF db1
< error
| 1064 42000 …near 'db1' at line 1
> exec
| CREATE TABLE IF NOT EXISTS t (id INT)
< ok 0 0
> exec
| CREATE INDEX IF NOT EXISTS idx ON t (id)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db1
< ok 0 0
> exec
| -- ============================================================
| -- Section 2.7: Column Constraint Keyword Chains
| -- ============================================================
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (id INT PRIMARY)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (id INT NOT NULL)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT UNIQUE)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (id INT UNIQUE)
< ok 0 0
> exec
| CREATE TABLE t (id INT UNIQUE KEY)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (id INT UNIQUE KEY)
< ok 0 0
> exec
| -- ============================================================
| -- Section 2.8: Structural Delimiter Enforcement
| -- ============================================================
< ok 0 0
> exec
| CREATE TABLE t (LIKE old_t
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (id INT
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (id INT, name VARCHAR(50)
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (id INT)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (id INT)
< ok 0 0
> exec
| CREATE TABLE t (LIKE old_t)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TABLE t (LIKE old_t)
< error
| 1146 42S02 Table 'test.old_t' doesn't exist
> exec ^
: string
| CREATE TABLE t (LIKE old_t)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @a = 1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET CHARACTER SET utf8mb4
< ok 0 0
> exec
| SET GLOBAL max_connections = 100
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES'
< ok 0 0
> exec
| SHOW DATABASES
< ok 0 0
> exec
| SHOW TABLES
< ok 0 0
> exec
| SHOW COLUMNS FROM employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SHOW COLUMNS FROM employees
< ok 0 0
> exec
| SHOW INDEX FROM employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SHOW INDEX FROM employees
< ok 0 0
> exec
| SHOW CREATE TABLE employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SHOW CREATE TABLE employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec ^
: string
| SHOW CREATE TABLE employees
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SHOW TABLE STATUS
< ok 0 0
> exec
| SHOW VARIABLES LIKE 'max%'
< ok 0 0
> exec
| SHOW PROCESSLIST
< ok 0 0
> exec
| SHOW WARNINGS
< ok 0 0
> exec
| SHOW ERRORS
< ok 0 0
> exec
| USE mydb
< error
| 1049 42000 Unknown database 'mydb'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| USE mydb
< ok 0 0
> exec
| BEGIN
< ok 0 0
> exec
| START TRANSACTION
< ok 0 0
> exec
| START TRANSACTION READ ONLY
< ok 0 0
> exec
| COMMIT
< ok 0 0
> exec
| ROLLBACK
< ok 0 0
> exec
| SAVEPOINT sp1
< ok 0 0
> exec
| ROLLBACK TO SAVEPOINT sp1
< error
| 1305 42000 SAVEPOINT sp1 does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ROLLBACK TO SAVEPOINT sp1
< ok 0 0
> exec
| RELEASE SAVEPOINT sp1
< error
| 1305 42000 SAVEPOINT sp1 does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| RELEASE SAVEPOINT sp1
< ok 0 0
> exec
| EXPLAIN SELECT * FROM employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| EXPLAIN SELECT * FROM employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec ^
: string
| EXPLAIN SELECT * FROM employees
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| DESCRIBE employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| DESCRIBE employees
< ok 0 0
> exec
| PREPARE stmt FROM 'SELECT ?'
< ok 0 0
> exec
| EXECUTE stmt USING @a
< ok 0 0
> exec
| DEALLOCATE PREPARE stmt
< ok 0 0
> exec
| FLUSH TABLES
< ok 0 0
> exec
| FLUSH PRIVILEGES
< ok 0 0
> exec
| LOCK TABLES t1 WRITE, t2 READ
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| LOCK TABLES t1 WRITE, t2 READ
< ok 0 0
> exec
| UNLOCK TABLES
< ok 0 0
> exec
| ANALYZE TABLE employees
< ok 0 0
> exec
| OPTIMIZE TABLE employees
< ok 0 0
> exec
| CHECK TABLE employees
< ok 0 0
> exec
| REPAIR TABLE employees
< ok 0 0
> exec
| GRANT SELECT ON mydb.* TO 'appuser'@'localhost'
< ok 0 0
> exec
| GRANT SELECT, INSERT, UPDATE ON mydb.employees TO 'appuser'@'localhost'
< error
| 1146 42S02 Table 'mydb.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| GRANT SELECT, INSERT, UPDATE ON mydb.employees TO 'appuser'@'localhost'
< ok 0 0
> exec
| GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost'
< ok 0 0
> exec
| GRANT SELECT ON mydb.* TO 'lead'@'localhost' WITH GRANT OPTION
< ok 0 0
> exec
| GRANT app_role TO 'user1'@'localhost'
< error
| 1959 OP000 Invalid role specification `app_role`
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| GRANT app_role TO 'user1'@'localhost'
< ok 0 0
> exec
| REVOKE SELECT ON mydb.* FROM 'appuser'@'localhost'
< ok 0 0
> exec
| REVOKE ALL PRIVILEGES ON *.* FROM 'appuser'@'localhost'
< ok 0 0
> exec
| CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'password123'
< ok 0 0
> exec
| CREATE USER IF NOT EXISTS 'testuser'@'localhost' IDENTIFIED BY 'password123'
< ok 0 0
> exec
| DROP USER 'testuser'@'localhost'
< ok 0 0
> exec
| DROP USER IF EXISTS 'testuser'@'localhost'
< ok 0 0
> exec
| ALTER USER 'testuser'@'localhost' IDENTIFIED BY 'newpassword'
< error
| 1396 HY000 Operation ALTER USER failed for 'testuser'@'localhost'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER USER 'testuser'@'localhost' IDENTIFIED BY 'newpassword'
< ok 0 0
> exec
| CREATE ROLE app_role
< ok 0 0
> exec
| DROP ROLE app_role
< ok 0 0
> exec
| SET PASSWORD FOR 'testuser'@'localhost' = 'newpass'
< error
| 1133 28000 Can't find any matching row in the user table
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SET PASSWORD FOR 'testuser'@'localhost' = 'newpass'
< ok 0 0
> exec
| CREATE TRIGGER trg_audit BEFORE INSERT ON employees FOR EACH ROW SET NEW.created_at = NOW()
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TRIGGER trg_audit BEFORE INSERT ON employees FOR EACH ROW SET NEW.created_at = NOW()
< ok 0 0
> exec
| CREATE TRIGGER trg_update BEFORE UPDATE ON employees FOR EACH ROW SET NEW.updated_at = NOW()
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TRIGGER trg_update BEFORE UPDATE ON employees FOR EACH ROW SET NEW.updated_at = NOW()
< ok 0 0
> exec
| DROP TRIGGER trg_audit
< error
| 1360 HY000 Trigger does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| DROP TRIGGER trg_audit
< ok 0 0
> exec
| DROP TRIGGER IF EXISTS trg_audit
< ok 0 0
> exec
| DROP PROCEDURE IF EXISTS my_proc
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS my_func
< ok 0 0
> exec
| CALL my_proc()
< error
| 1305 42000 PROCEDURE test.my_proc does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CALL my_proc()
< ok 0 0
> exec
| CALL my_proc(1, 'test')
< error
| 1305 42000 PROCEDURE test.my_proc does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CALL my_proc(1, 'test')
< ok 0 0
> exec
| DO 1
< ok 0 0
> exec
| DO SLEEP(1)
< ok 0 0
> exec
| LOAD DATA INFILE '/tmp/data.csv' INTO TABLE employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| LOAD DATA INFILE '/tmp/data.csv' INTO TABLE employees
< error
| 1146 42S02 Table 'test.employees' doesn't exist
> exec ^
: string
| LOAD DATA INFILE '/tmp/data.csv' INTO TABLE employees
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE TABLE t1 (a INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t3 (a INT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t4 (a INT(11))
< ok 0 0
> exec
| CREATE TABLE t5 (a INT UNSIGNED ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t6 (a TINYINT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t7 (a TINYINT UNSIGNED ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t8 (a SMALLINT(5) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t9 (a MEDIUMINT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t10 (a BIGINT(20) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t11 (a INT1 UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t12 (a INT8 SIGNED)
< ok 0 0
> exec
| CREATE TABLE t13 (a DECIMAL)
< ok 0 0
> exec
| CREATE TABLE t14 (a DECIMAL(10))
< ok 0 0
> exec
| CREATE TABLE t15 (a DECIMAL(10,2))
< ok 0 0
> exec
| CREATE TABLE t16 (a DECIMAL(10,2) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t17 (a NUMERIC(10,2))
< ok 0 0
> exec
| CREATE TABLE t18 (a DEC(10,2))
< ok 0 0
> exec
| CREATE TABLE t19 (a FIXED(10,2))
< ok 0 0
> exec
| CREATE TABLE t20 (a FLOAT)
< ok 0 0
> exec
| CREATE TABLE t21 (a FLOAT(10,2))
< ok 0 0
> exec
| CREATE TABLE t22 (a FLOAT(24))
< ok 0 0
> exec
| CREATE TABLE t23 (a FLOAT(25))
< ok 0 0
> exec
| CREATE TABLE t24 (a DOUBLE)
< ok 0 0
> exec
| CREATE TABLE t25 (a DOUBLE PRECISION)
< ok 0 0
> exec
| CREATE TABLE t26 (a REAL)
< ok 0 0
> exec
| CREATE TABLE t27 (a FLOAT4)
< ok 0 0
> exec
| CREATE TABLE t28 (a FLOAT8)
< ok 0 0
> exec
| CREATE TABLE t29 (a CHAR(10))
< ok 0 0
> exec
| CREATE TABLE t30 (a CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)
< ok 0 0
> exec
| CREATE TABLE t31 (a VARCHAR(255))
< ok 0 0
> exec
| CREATE TABLE t32 (a TEXT)
< ok 0 0
> exec
| CREATE TABLE t33 (a TEXT(1000))
< ok 0 0
> exec
| CREATE TABLE t34 (a TINYTEXT)
< ok 0 0
> exec
| CREATE TABLE t35 (a MEDIUMTEXT CHARACTER SET latin1)
< ok 0 0
> exec
| CREATE TABLE t36 (a LONGTEXT)
< ok 0 0
> exec
| CREATE TABLE t37 (a LONG)
< ok 0 0
> exec
| CREATE TABLE t38 (a LONG VARCHAR)
< ok 0 0
> exec
| CREATE TABLE t39 (a BINARY(16))
< ok 0 0
> exec
| CREATE TABLE t40 (a VARBINARY(255))
< ok 0 0
> exec
| CREATE TABLE t41 (a BLOB)
< ok 0 0
> exec
| CREATE TABLE t42 (a BLOB(1000))
< ok 0 0
> exec
| CREATE TABLE t43 (a TINYBLOB)
< ok 0 0
> exec
| CREATE TABLE t44 (a MEDIUMBLOB)
< ok 0 0
> exec
| CREATE TABLE t45 (a LONGBLOB)
< ok 0 0
> exec
| CREATE TABLE t46 (a LONG VARBINARY)
< ok 0 0
> exec
| CREATE TABLE t47 (a NATIONAL CHAR(10))
< ok 0 0
> exec
| CREATE TABLE t48 (a NCHAR(10))
< ok 0 0
> exec
| CREATE TABLE t49 (a NVARCHAR(100))
< ok 0 0
> exec
| CREATE TABLE t50 (a DATE)
< ok 0 0
> exec
| CREATE TABLE t51 (a TIME)
< ok 0 0
> exec
| CREATE TABLE t52 (a TIME(3))
< ok 0 0
> exec
| CREATE TABLE t53 (a DATETIME)
< ok 0 0
> exec
| CREATE TABLE t54 (a DATETIME(6))
< ok 0 0
> exec
| CREATE TABLE t55 (a TIMESTAMP)
< ok 0 0
> exec
| CREATE TABLE t56 (a TIMESTAMP(3))
< ok 0 0
> exec
| CREATE TABLE t57 (a YEAR)
< ok 0 0
> exec
| CREATE TABLE t58 (a BIT(8))
< ok 0 0
> exec
| CREATE TABLE t59 (a BOOL)
< ok 0 0
> exec
| CREATE TABLE t60 (a JSON)
< ok 0 0
> exec
| CREATE TABLE t61 (a SERIAL)
< ok 0 0
> exec
| CREATE TABLE t62 (a ENUM('a','b','c'))
< ok 0 0
> exec
| CREATE TABLE t63 (a SET('x','y','z'))
< ok 0 0
> exec
| CREATE TABLE t64 (a GEOMETRY)
< ok 0 0
> exec
| CREATE TABLE t65 (a POINT)
< ok 0 0
> exec
| CREATE TABLE t66 (a LINESTRING)
< ok 0 0
> exec
| CREATE TABLE t67 (a POLYGON)
< ok 0 0
> exec
| CREATE TABLE t68 (a GEOMETRYCOLLECTION)
< ok 0 0
> exec
| SELECT RANK() OVER w FROM t WINDOW w AS (ORDER BY id)
< ok 0 0
> exec
| SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t
< error
| 1054 42S22 Unknown column 'dept' in 'GROUP BY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t
< error
| 1054 42S22 Unknown column 'dept' in 'GROUP BY'
> exec ^
: string
| SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT LAG(val) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT LAG(val) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec ^
: string
| SELECT LAG(val) OVER (ORDER BY id) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT LAG(val, 2, 0) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ' 0) OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec ^
: string
| SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec ^
: string
| SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t
< error
| 1054 42S22 Unknown column 'score' in 'ORDER BY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t
< error
| 1054 42S22 Unknown column 'score' in 'ORDER BY'
> exec ^
: string
| SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT NTILE(4) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t
< error
| 1054 42S22 Unknown column 'val' in 'SELECT'
> exec ^
: string
| SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT COUNT(*) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT MEDIAN(val) OVER (PARTITION BY dept) FROM t
< error
| 1054 42S22 Unknown column 'dept' in 'GROUP BY'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT MEDIAN(val) OVER (PARTITION BY dept) FROM t
< error
| 1054 42S22 Unknown column 'dept' in 'GROUP BY'
> exec ^
: string
| SELECT MEDIAN(val) OVER (PARTITION BY dept) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT NTILE() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NTILE(2, 3) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ' 3) OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT LEAD() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT FIRST_VALUE() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT ROW_NUMBER(id) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'id) OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT ABS(val) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT `sum`(val) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT LAST_VALUE(1, 2) FROM t
< ok 0 0
> exec
| SELECT LAST_VALUE(1, 2) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NOW() + INTERVAL 1 DAY
< ok 0 0
> exec
| SELECT NOW() - INTERVAL 2 HOUR
< ok 0 0
> exec
| SELECT DATE_ADD('2024-01-01', INTERVAL 1 MONTH)
< ok 0 0
> exec
| SELECT DATE_SUB(NOW(), INTERVAL 30 MINUTE)
< ok 0 0
> exec
| SELECT TIMESTAMPADD(MINUTE, 30, NOW())
< ok 0 0
> exec
| SELECT EXTRACT(HOUR FROM NOW())
< ok 0 0
> exec
| SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY
< error
| 1054 42S22 Unknown column 'created_at' in 'WHERE'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY
< error
| 1054 42S22 Unknown column 'created_at' in 'WHERE'
> exec ^
: string
| SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()
< error
| 1054 42S22 Unknown column 'created_at' in 'WHERE'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()
< error
| 1054 42S22 Unknown column 'created_at' in 'WHERE'
> exec ^
: string
| SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE TABLE t_rej1 (a VARCHAR UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| CREATE TABLE t_rej2 (a TEXT ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec
| CREATE TABLE t_rej3 (a JSON UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| CREATE TABLE select (a INT)
< error
| 1064 42000 …near 'select (a INT)' at line 1
> exec
| CREATE TABLE t_rej5 (select INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| ALTER TABLE t PARTITION BY
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t PARTITION BY RANGE(id)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| ALTER TABLE t PARTITION BY RANGE(id)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE(id)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t
< error
| 1064 42000 …near 'INVALID_UNIT) FROM t' at line 1
> exec
| SELECT 1 MINUS SELECT 2
< error
| 1064 42000 …near 'SELECT 2' at line 1
> exec
| SELECT 1 MINUS SELECT 2 MINUS SELECT 3
< error
| 1064 42000 …near 'SELECT 2 MINUS SELECT 3' at line 1
> exec
| SELECT name -> '$.b' FROM t
< error
| 1064 42000 …near '> '$.b' FROM t' at line 1
> exec
| SELECT name ->> '$.b' FROM t
< error
| 1064 42000 …near '>> '$.b' FROM t' at line 1
> exec
| SELECT * FROM t WHERE id = 1 FOR SHARE
< error
| 1064 42000 …near 'SHARE' at line 1
> exec
| SELECT CAST(name AS JSON ARRAY) FROM t
< error
| 1064 42000 …near 'JSON ARRAY) FROM t' at line 1
> exec
| /*!50000 SELECT 'a*/' */
< ok 0 0
> exec
| SELECT WEIGHT_STRING(name AS CHAR(4)) FROM t
< error
| 1054 42S22 Unknown column 'name' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT WEIGHT_STRING(name AS CHAR(4)) FROM t
< error
| 1054 42S22 Unknown column 'name' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(name AS CHAR(4)) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT WEIGHT_STRING(name LEVEL 1) FROM t
< error
| 1054 42S22 Unknown column 'name' in 'SELECT'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT WEIGHT_STRING(name LEVEL 1) FROM t
< error
| 1054 42S22 Unknown column 'name' in 'SELECT'
> exec ^
: string
| SELECT WEIGHT_STRING(name LEVEL 1) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| CREATE TABLE t_adjacent_default (c VARCHAR(8) DEFAULT 'a' 'b')
< ok 0 0
> exec
| SELECT * FROM t WHERE name = 'a' 'b'
< error
| 1054 42S22 Unknown column 'name' in 'WHERE'
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT * FROM t WHERE name = 'a' 'b'
< error
| 1054 42S22 Unknown column 'name' in 'WHERE'
> exec ^
: string
| SELECT * FROM t WHERE name = 'a' 'b'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT * FROM (VALUES ('a' 'b')) x
< ok 0 0
