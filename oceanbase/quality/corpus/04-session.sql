-- @name: set user variable
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET @a = 1

-- @name: set names
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET NAMES utf8mb4

-- @name: set character set
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET CHARACTER SET utf8mb4

-- @name: set global variable
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET GLOBAL max_connections = 100

-- @name: set session variable
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SET SESSION sql_mode = 'STRICT_TRANS_TABLES'

-- @name: show databases
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW DATABASES

-- @name: show tables
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW TABLES

-- @name: show columns
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW COLUMNS FROM employees

-- @name: show index
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW INDEX FROM employees

-- @name: show create table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW CREATE TABLE employees

-- @name: show table status
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW TABLE STATUS

-- @name: show variables
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW VARIABLES LIKE 'max%'

-- @name: show processlist
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW PROCESSLIST

-- @name: show warnings
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW WARNINGS

-- @name: show errors
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SHOW ERRORS

-- @name: use database
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
USE mydb

-- @name: begin
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
BEGIN

-- @name: start transaction
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
START TRANSACTION

-- @name: start transaction read only
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
START TRANSACTION READ ONLY

-- @name: commit
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
COMMIT

-- @name: rollback
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
ROLLBACK

-- @name: savepoint
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
SAVEPOINT sp1

-- @name: rollback to savepoint
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
ROLLBACK TO SAVEPOINT sp1

-- @name: release savepoint
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
RELEASE SAVEPOINT sp1

-- @name: explain select
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
EXPLAIN SELECT * FROM employees

-- @name: describe table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DESCRIBE employees

-- @name: prepare statement
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
PREPARE stmt FROM 'SELECT ?'

-- @name: execute statement
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
EXECUTE stmt USING @a

-- @name: deallocate prepare
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
DEALLOCATE PREPARE stmt

-- @name: flush tables
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
FLUSH TABLES

-- @name: flush privileges
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
FLUSH PRIVILEGES

-- @name: lock tables
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
LOCK TABLES t1 WRITE, t2 READ

-- @name: unlock tables
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
UNLOCK TABLES

-- @name: analyze table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
ANALYZE TABLE employees

-- @name: optimize table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
OPTIMIZE TABLE employees

-- @name: check table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
CHECK TABLE employees

-- @name: repair table
-- @valid: true
-- @source: MySQL 8.0 Reference Manual
REPAIR TABLE employees
