# MySQL Catalog Compatibility Scenarios

> Ultimate goal: every scenario passes oracle testing against real MySQL 8.0.
> Verification surface: `SHOW CREATE TABLE` output exact match + error code/message match.

Status legend: `[ ]` pending, `[x]` passing, `[~]` partial, `[-]` not applicable

---

## Phase 1: SHOW CREATE TABLE Precision

The foundation — our output must be byte-for-byte identical to MySQL 8.0.

### 1.1 Numeric Types

```
[x] INT / INTEGER
[x] INT(11) — display width (deprecated 8.0.17+ but still shown)
[x] INT UNSIGNED
[x] INT UNSIGNED ZEROFILL
[x] TINYINT, SMALLINT, MEDIUMINT, BIGINT
[x] BIGINT UNSIGNED (common for IDs)
[x] FLOAT
[x] FLOAT(7,3)
[x] FLOAT UNSIGNED
[x] DOUBLE / DOUBLE PRECISION
[x] DOUBLE(15,5)
[x] DECIMAL(10,2) / NUMERIC(10,2)
[x] DECIMAL (no precision — MySQL shows DECIMAL(10,0))
[x] BOOLEAN / BOOL (→ tinyint(1))
[x] BIT(1), BIT(8), BIT(64)
[x] SERIAL (→ BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE)
```

### 1.2 String Types

```
[x] CHAR(10)
[x] CHAR (no length → CHAR(1))
[x] VARCHAR(255)
[x] VARCHAR(16383) — max for utf8mb4
[x] TINYTEXT, TEXT, MEDIUMTEXT, LONGTEXT
[x] TEXT(1000) — display shows TEXT, not TEXT(1000)
[x] ENUM('a','b','c')
[x] ENUM with special characters in values (quotes, commas)
[x] SET('x','y','z')
[x] CHAR with CHARACTER SET latin1
[x] VARCHAR with CHARACTER SET utf8mb3 COLLATE utf8mb3_bin
[x] NATIONAL CHAR / NCHAR / NVARCHAR
```

### 1.3 Binary Types

```
[x] BINARY(16)
[x] BINARY (no length → BINARY(1))
[x] VARBINARY(255)
[x] TINYBLOB, BLOB, MEDIUMBLOB, LONGBLOB
[x] BLOB(1000) — display shows BLOB, not BLOB(1000)
```

### 1.4 Date/Time Types

```
[x] DATE
[x] TIME
[x] TIME(3) — fractional seconds
[x] DATETIME
[x] DATETIME(6)
[x] TIMESTAMP
[x] TIMESTAMP(3)
[x] YEAR
[x] YEAR(4) — same as YEAR
```

### 1.5 Spatial Types

```
[x] GEOMETRY
[x] POINT
[x] LINESTRING
[x] POLYGON
[x] MULTIPOINT, MULTILINESTRING, MULTIPOLYGON
[x] GEOMETRYCOLLECTION
[~] POINT NOT NULL SRID 4326 — requires parser support for SRID on column definitions
```

### 1.6 JSON Type

```
[x] JSON
[x] JSON DEFAULT NULL — MySQL 8.0 DOES show DEFAULT NULL for JSON (unlike TEXT/BLOB)
```

### 1.7 Default Values

```
[x] INT DEFAULT 0 → DEFAULT '0' (quoted)
[x] INT DEFAULT NULL
[x] INT NOT NULL (no default shown)
[x] VARCHAR DEFAULT 'hello'
[x] VARCHAR DEFAULT '' (empty string)
[x] FLOAT DEFAULT 3.14 → DEFAULT '3.14'
[x] DECIMAL(10,2) DEFAULT 0.00 → DEFAULT '0.00'
[x] BOOLEAN DEFAULT TRUE → tinyint(1) DEFAULT '1'
[x] BOOLEAN DEFAULT FALSE → tinyint(1) DEFAULT '0'
[x] ENUM DEFAULT 'a'
[x] SET DEFAULT 'x,y'
[x] BIT(8) DEFAULT b'00001111'
[x] BLOB/TEXT — no DEFAULT NULL shown
[x] JSON — MySQL 8.0 shows DEFAULT NULL (unlike TEXT/BLOB)
[x] TIMESTAMP DEFAULT CURRENT_TIMESTAMP
[x] DATETIME DEFAULT CURRENT_TIMESTAMP
[x] TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3)
[x] Expression default (8.0.13+): INT DEFAULT (FLOOR(RAND()*100))
[x] Expression default: JSON DEFAULT (JSON_ARRAY())
[x] Expression default: VARCHAR DEFAULT (UUID())
[x] DATETIME DEFAULT '2024-01-01 00:00:00'
[x] DATE DEFAULT '2024-01-01'
```

### 1.8 ON UPDATE

```
[x] TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
[x] DATETIME(3) ON UPDATE CURRENT_TIMESTAMP(3)
[x] TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
[x] DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
```

### 1.9 Generated Columns

```
[x] INT GENERATED ALWAYS AS (col1 + col2) VIRTUAL
[x] INT GENERATED ALWAYS AS (col1 * col2) STORED
[x] VARCHAR(255) AS (CONCAT(first_name, ' ', last_name)) VIRTUAL
[x] Generated column with NOT NULL
[x] Generated column with COMMENT
[x] Generated column with INVISIBLE
[x] JSON GENERATED ALWAYS AS (JSON_EXTRACT(data, '$.name'))
```

### 1.10 Column Attributes Combination

```
[x] INT NOT NULL AUTO_INCREMENT
[x] BIGINT UNSIGNED NOT NULL AUTO_INCREMENT
[x] VARCHAR(100) NOT NULL DEFAULT ''
[x] VARCHAR(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL
[x] INT NOT NULL COMMENT 'user id'
[x] VARCHAR(255) INVISIBLE (8.0.23+)
[x] INT VISIBLE (should not show VISIBLE — it's the default)
[x] Column with all: INT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'count'
```

### 1.11 Primary Key

```
[x] Single column PK: PRIMARY KEY (id)
[x] Multi-column PK: PRIMARY KEY (a, b)
[x] PK on BIGINT UNSIGNED AUTO_INCREMENT
[x] PK column ordering in SHOW CREATE TABLE
[x] PK implicitly makes columns NOT NULL
[x] PK name is never shown (always PRIMARY KEY, not KEY `PRIMARY`)
```

### 1.12 Unique Keys

```
[x] UNIQUE KEY `name` (col)
[x] UNIQUE KEY without name → auto-named from first column
[x] Multi-column UNIQUE KEY
[x] Multiple UNIQUE KEYs on same table
[x] UNIQUE KEY auto-naming collision (col, col_2, col_3)
```

### 1.13 Regular Indexes

```
[x] KEY `idx_name` (col)
[x] KEY without name → auto-named
[x] Multi-column KEY
[x] KEY with prefix length: KEY `idx` (col(10))
[x] KEY with DESC: KEY `idx` (col DESC) (8.0+)
[x] KEY with mixed ASC/DESC: KEY `idx` (a ASC, b DESC)
[x] USING HASH — shown only for non-BTREE
[x] USING BTREE — shown when explicitly specified
```

### 1.14 Fulltext and Spatial Indexes

```
[x] FULLTEXT KEY `ft_idx` (col)
[x] FULLTEXT KEY on multiple TEXT columns
[x] FULLTEXT KEY auto-naming
[x] SPATIAL KEY `sp_idx` (geo_col)
```

### 1.15 Expression Indexes (8.0.13+)

```
[x] KEY `idx` ((UPPER(name)))
[x] KEY `idx` ((col1 + col2))
[x] UNIQUE KEY on expression
[x] Expression index display format in SHOW CREATE TABLE
```

### 1.16 Index Options

```
[x] INDEX with COMMENT 'description'
[x] INDEX INVISIBLE (8.0+)
[x] INDEX VISIBLE — not shown (default)
[x] INDEX with KEY_BLOCK_SIZE=4
```

### 1.17 Foreign Keys

```
[x] Basic FK: FOREIGN KEY (col) REFERENCES parent(id)
[x] FK with name: CONSTRAINT `fk_name` FOREIGN KEY ...
[x] FK ON DELETE CASCADE
[x] FK ON DELETE SET NULL
[x] FK ON DELETE SET DEFAULT
[x] FK ON DELETE RESTRICT — shown when explicitly specified (unlike NO ACTION)
[x] FK ON DELETE NO ACTION — not shown (MySQL default)
[x] FK ON UPDATE CASCADE
[x] FK ON UPDATE SET NULL
[x] FK combined: ON DELETE CASCADE ON UPDATE SET NULL
[x] FK auto-naming: table_ibfk_1, table_ibfk_2
[x] FK auto-generates backing index (KEY `fk_name` (col))
[x] FK self-referencing: REFERENCES same_table(id)
[x] FK multi-column: FOREIGN KEY (a, b) REFERENCES parent(x, y)
[~] FK to different database: explicit schema qualification — requires cross-database test infrastructure
```

### 1.18 Check Constraints (8.0.16+)

```
[x] CHECK (col > 0)
[x] CHECK with name: CONSTRAINT `chk_name` CHECK (col > 0)
[x] CHECK NOT ENFORCED → /*!80016 NOT ENFORCED */
[x] CHECK auto-naming: table_chk_1, table_chk_2
[x] CHECK expression display format (parenthesization)
[x] CHECK referencing multiple columns
```

### 1.19 Table Options

```
[x] ENGINE=InnoDB (default)
[x] ENGINE=MyISAM
[x] ENGINE=MEMORY
[x] DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
[x] DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci
[x] DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci (non-default collation)
[x] COMMENT='table description'
[x] COMMENT with special chars (quotes, backslash)
[x] ROW_FORMAT=DYNAMIC
[x] ROW_FORMAT=COMPRESSED
[x] AUTO_INCREMENT=1000 (shown only if > 1)
[x] KEY_BLOCK_SIZE=8
[x] Multiple options ordering in output
```

### 1.20 Charset/Collation Inheritance

```
[x] Table charset inherited from database
[x] Column charset inherited from table
[x] Column charset overrides table charset
[x] Column collation overrides table collation
[x] Column charset/collation display rules (when shown vs omitted)
[x] Binary charset on column
```

---

## Phase 2: DDL Behavior Precision

Beyond output format — the catalog must behave identically to MySQL.

### 2.1 CREATE TABLE Variants

```
[x] CREATE TABLE IF NOT EXISTS — no error when table exists
[x] CREATE TEMPORARY TABLE
[x] CREATE TABLE t2 LIKE t1 — copies structure
[~] CREATE TABLE t2 AS SELECT * FROM t1 — CTAS — requires query execution engine (not catalog-level)
[x] CREATE TABLE with same name as existing view — error
[x] CREATE TABLE with reserved word as name (backtick-quoted)
```

### 2.2 ALTER TABLE — Column Operations

```
[x] ADD COLUMN at end
[x] ADD COLUMN FIRST
[x] ADD COLUMN AFTER specific_col
[x] ADD multiple columns in one ALTER
[x] DROP COLUMN
[x] DROP COLUMN that's part of an index → index removed or shrunk
[x] DROP COLUMN that's the only column in an index → index removed
[x] DROP COLUMN that's part of PK → error
[x] DROP COLUMN that's referenced by FK → error
[~] DROP COLUMN IF EXISTS — MySQL 8.0 does not support IF EXISTS on DROP COLUMN (syntax error); supported only in 8.0.32+
[x] MODIFY COLUMN — change type, keep name
[x] MODIFY COLUMN — widen VARCHAR
[x] MODIFY COLUMN — narrow VARCHAR (data truncation warning in real MySQL)
[x] MODIFY COLUMN — change INT to BIGINT
[x] MODIFY COLUMN — add NOT NULL
[x] MODIFY COLUMN — remove NOT NULL
[x] MODIFY COLUMN FIRST / AFTER
[x] CHANGE COLUMN — rename + change type
[x] CHANGE COLUMN — same name, different type
[x] CHANGE COLUMN — update index column references
[x] RENAME COLUMN old TO new
[x] RENAME COLUMN — update index column references
[x] ALTER COLUMN col SET DEFAULT value
[x] ALTER COLUMN col DROP DEFAULT
[x] ALTER COLUMN col SET VISIBLE
[x] ALTER COLUMN col SET INVISIBLE
```

### 2.3 ALTER TABLE — Index Operations

```
[x] ADD INDEX idx_name (col)
[x] ADD UNIQUE INDEX
[x] ADD FULLTEXT INDEX
[x] ADD PRIMARY KEY
[x] ADD PRIMARY KEY when one already exists → error 1068
[x] DROP INDEX idx_name
[~] DROP INDEX IF EXISTS — MySQL 8.0 does not support IF EXISTS on DROP INDEX (syntax error)
[x] DROP PRIMARY KEY
[x] RENAME INDEX old TO new
[x] ALTER INDEX idx VISIBLE
[x] ALTER INDEX idx INVISIBLE
```

### 2.4 ALTER TABLE — Constraint Operations

```
[x] ADD CONSTRAINT fk FOREIGN KEY ...
[x] ADD CHECK (expr)
[x] DROP FOREIGN KEY fk_name
[x] DROP CHECK chk_name
[x] DROP CONSTRAINT name (generic)
[x] ALTER CHECK name ENFORCED
[x] ALTER CHECK name NOT ENFORCED
```

### 2.5 ALTER TABLE — Table-level

```
[x] RENAME TO new_name
[x] ENGINE=MyISAM (change engine)
[x] CONVERT TO CHARACTER SET utf8mb4
[x] DEFAULT CHARACTER SET latin1
[x] COMMENT='new comment'
[x] AUTO_INCREMENT=1000
[x] ROW_FORMAT=COMPRESSED
```

### 2.6 DROP TABLE

```
[x] DROP TABLE t1
[x] DROP TABLE IF EXISTS t1
[x] DROP TABLE t1, t2, t3 — multi-table
[x] DROP TABLE nonexistent → error 1051
[x] DROP TEMPORARY TABLE
[x] DROP TABLE that has FK references → error (without foreign_key_checks=0)
```

### 2.7 TRUNCATE TABLE

```
[x] TRUNCATE TABLE t1
[x] TRUNCATE resets AUTO_INCREMENT
[x] TRUNCATE nonexistent table → error
```

### 2.8 CREATE/DROP INDEX (standalone)

```
[x] CREATE INDEX idx ON t (col)
[x] CREATE UNIQUE INDEX
[x] CREATE FULLTEXT INDEX
[x] CREATE SPATIAL INDEX
[~] CREATE INDEX IF NOT EXISTS — MySQL 8.0 does not support IF NOT EXISTS on CREATE INDEX (syntax error); omni accepts it as extension
[x] CREATE INDEX — duplicate name → error 1061
[x] DROP INDEX idx ON t
[x] DROP INDEX nonexistent → error 1091
```

### 2.9 RENAME TABLE

```
[x] RENAME TABLE t1 TO t2
[x] RENAME TABLE t1 TO db2.t1 (cross-database)
[x] RENAME TABLE t1 TO t2, t3 TO t4 (multi-pair)
[x] RENAME TABLE nonexistent → error
[x] RENAME TABLE to existing name → error
```

### 2.10 CREATE/DROP VIEW

```
[x] CREATE VIEW v AS SELECT ...
[x] CREATE OR REPLACE VIEW
[x] CREATE VIEW with column list
[x] CREATE VIEW with ALGORITHM, DEFINER, SQL_SECURITY
[x] CREATE VIEW with CHECK OPTION
[x] DROP VIEW v
[x] DROP VIEW IF EXISTS
[x] DROP VIEW v1, v2 (multi-view)
```

### 2.11 CREATE/DROP/ALTER DATABASE

```
[ ] CREATE DATABASE db
[ ] CREATE DATABASE IF NOT EXISTS
[ ] CREATE DATABASE with CHARACTER SET
[ ] CREATE DATABASE with COLLATE
[ ] DROP DATABASE db
[ ] DROP DATABASE IF EXISTS
[ ] ALTER DATABASE db CHARACTER SET utf8mb4
[ ] ALTER DATABASE db COLLATE utf8mb4_unicode_ci
[ ] Operations on nonexistent database → proper errors
```

---

## Phase 3: Error Precision

Every error must match MySQL's errno, SQLSTATE, and message format.

### 3.1 Database Errors

```
[ ] 1007 (HY000) Can't create database 'db'; database exists
[ ] 1049 (42000) Unknown database 'db'
[ ] 1046 (3D000) No database selected
```

### 3.2 Table Errors

```
[ ] 1050 (42S01) Table 'tbl' already exists
[ ] 1051 (42S02) Unknown table 'db.tbl'
[ ] 1146 (42S02) Table 'db.tbl' doesn't exist
```

### 3.3 Column Errors

```
[ ] 1054 (42S22) Unknown column 'col' in 'table definition'
[ ] 1060 (42S21) Duplicate column name 'col'
[ ] 1068 (42000) Multiple primary key defined
```

### 3.4 Index/Key Errors

```
[ ] 1061 (42000) Duplicate key name 'idx'
[ ] 1091 (42000) Can't DROP 'idx'; check that column/key exists
```

### 3.5 FK Errors

```
[ ] 1824 (HY000) Failed to open the referenced table 'tbl'
[ ] 1822 (HY000) A foreign key constraint fails
[ ] FK column type mismatch detection
```

### 3.6 Error Context

```
[ ] Error message identifier quoting matches MySQL
[ ] Error position (line number) for multi-statement SQL
[ ] IF EXISTS suppresses errors correctly
[ ] IF NOT EXISTS suppresses errors correctly
```

---

## Phase 4: Advanced DDL

### 4.1 Partitioning

```
[ ] PARTITION BY RANGE (col)
[ ] PARTITION BY RANGE COLUMNS (col)
[ ] PARTITION BY LIST (col)
[ ] PARTITION BY HASH (col)
[ ] PARTITION BY KEY (col)
[ ] LINEAR HASH
[ ] SUBPARTITION
[ ] SHOW CREATE TABLE output for partitioned tables
[ ] ALTER TABLE ADD PARTITION
[ ] ALTER TABLE DROP PARTITION
[ ] ALTER TABLE REORGANIZE PARTITION
[ ] ALTER TABLE TRUNCATE PARTITION
[ ] ALTER TABLE COALESCE PARTITION
[ ] ALTER TABLE EXCHANGE PARTITION
```

### 4.2 Stored Routines

```
[ ] CREATE FUNCTION — store metadata (name, params, return type, body)
[ ] CREATE PROCEDURE — store metadata
[ ] DROP FUNCTION / PROCEDURE
[ ] ALTER ROUTINE (characteristics only)
[ ] SHOW CREATE FUNCTION output
[ ] SHOW CREATE PROCEDURE output
```

### 4.3 Triggers

```
[ ] CREATE TRIGGER — store metadata (name, timing, event, table, body)
[ ] DROP TRIGGER
[ ] SHOW CREATE TRIGGER output
[ ] Multiple triggers per table/event (MySQL 8.0 supports ordering)
```

### 4.4 Events

```
[ ] CREATE EVENT — store metadata
[ ] ALTER EVENT
[ ] DROP EVENT
[ ] SHOW CREATE EVENT output
```

### 4.5 Views — Deep

```
[ ] ALTER VIEW
[ ] View dependency tracking (base table dropped → view invalid)
[ ] SHOW CREATE VIEW output format
[ ] View with column aliases
```

---

## Phase 5: Session & System Behavior

### 5.1 USE Statement

```
[ ] USE db — sets current database
[ ] USE nonexistent → error 1049
```

### 5.2 SET Variables (affecting DDL behavior)

```
[ ] SET foreign_key_checks = 0 — skip FK validation
[ ] SET foreign_key_checks = 1 — enforce FK validation
[ ] SET NAMES utf8mb4
[ ] SET CHARACTER SET utf8mb4
[ ] SET sql_mode (affects parsing and behavior)
```

### 5.3 User/Role Management

```
[ ] CREATE USER 'user'@'host'
[ ] DROP USER
[ ] ALTER USER
[ ] CREATE ROLE
[ ] DROP ROLE
[ ] GRANT privileges
[ ] REVOKE privileges
```

---

## Phase 6: Query APIs

### 6.1 SHOW CREATE TABLE

```
[x] Basic table output
[x] Column formatting
[x] Index formatting
[x] Constraint formatting
[x] Table options
[ ] All data types (Phase 1.1-1.6)
[ ] All default value forms (Phase 1.7)
[ ] All index types (Phase 1.13-1.16)
[ ] All constraint forms (Phase 1.17-1.18)
[ ] Partitioned table output (Phase 4.1)
```

### 6.2 SHOW CREATE VIEW/FUNCTION/PROCEDURE/TRIGGER/EVENT

```
[ ] SHOW CREATE VIEW
[ ] SHOW CREATE FUNCTION
[ ] SHOW CREATE PROCEDURE
[ ] SHOW CREATE TRIGGER
[ ] SHOW CREATE EVENT
```

### 6.3 INFORMATION_SCHEMA Consistency

```
[ ] INFORMATION_SCHEMA.COLUMNS matches catalog state
[ ] INFORMATION_SCHEMA.STATISTICS matches catalog indexes
[ ] INFORMATION_SCHEMA.TABLE_CONSTRAINTS matches
[ ] INFORMATION_SCHEMA.KEY_COLUMN_USAGE matches
[ ] INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS matches
[ ] INFORMATION_SCHEMA.CHECK_CONSTRAINTS matches
```

---

## Progress Tracking

| Phase | Total | Done | % |
|-------|-------|------|---|
| 1. SHOW CREATE TABLE Precision | ~130 | 59 | 45% |
| 2. DDL Behavior Precision | ~75 | ~30 | 40% |
| 3. Error Precision | ~18 | ~8 | 44% |
| 4. Advanced DDL | ~25 | 0 | 0% |
| 5. Session & System | ~12 | 0 | 0% |
| 6. Query APIs | ~15 | 5 | 33% |
| **Total** | **~275** | **~102** | **~37%** |
