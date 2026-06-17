-- MySQL DDL corpus — keyword chain enforcement and structural delimiters
-- Covers: §2.6 IF NOT EXISTS / IF EXISTS chains, §2.7 column constraint chains, §2.8 structural delimiters
-- Note: §2.8 cases require no parser changes — expect(')') already enforces closing parens

-- ============================================================
-- Section 2.6: IF NOT EXISTS / IF EXISTS Keyword Chains
-- ============================================================

-- @name: CREATE TABLE IF missing NOT EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
CREATE TABLE IF t (id INT)

-- @name: CREATE TABLE IF NOT missing EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
CREATE TABLE IF NOT t (id INT)

-- @name: CREATE INDEX IF missing NOT EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
CREATE INDEX IF idx ON t (id)

-- @name: DROP TABLE IF missing EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP TABLE IF t

-- @name: DROP INDEX IF missing EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP INDEX IF idx ON t

-- @name: DROP VIEW IF missing EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP VIEW IF v

-- @name: DROP DATABASE IF missing EXISTS
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP DATABASE IF db1

-- @name: CREATE TABLE IF NOT EXISTS accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.6
CREATE TABLE IF NOT EXISTS t (id INT)

-- @name: CREATE INDEX IF NOT EXISTS accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.6
CREATE INDEX IF NOT EXISTS idx ON t (id)

-- @name: DROP TABLE IF EXISTS accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP TABLE IF EXISTS t

-- @name: DROP VIEW IF EXISTS accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP VIEW IF EXISTS v

-- @name: DROP DATABASE IF EXISTS accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.6
DROP DATABASE IF EXISTS db1

-- ============================================================
-- Section 2.7: Column Constraint Keyword Chains
-- ============================================================

-- @name: NOT missing NULL
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT NOT)

-- @name: PRIMARY missing KEY
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT PRIMARY)

-- @name: NOT NULL accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT NOT NULL)

-- @name: PRIMARY KEY accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT PRIMARY KEY)

-- @name: UNIQUE accepted without KEY
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT UNIQUE)

-- @name: UNIQUE KEY accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.7
CREATE TABLE t (id INT UNIQUE KEY)

-- ============================================================
-- Section 2.8: Structural Delimiter Enforcement
-- ============================================================

-- @name: LIKE missing closing paren
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.8
CREATE TABLE t (LIKE old_t

-- @name: single column missing closing paren
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.8
CREATE TABLE t (id INT

-- @name: multiple columns missing closing paren
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md §2.8
CREATE TABLE t (id INT, name VARCHAR(50)

-- @name: single column with closing paren accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.8
CREATE TABLE t (id INT)

-- @name: LIKE with closing paren accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md §2.8
CREATE TABLE t (LIKE old_t)
