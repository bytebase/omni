# MySQL Catalog Walkthrough Test Scenarios

> Goal: Achieve production-ready in-memory MySQL 8.0 catalog simulation by systematically verifying all catalog operations against real MySQL 8.0
> Verification: Each scenario compares omni catalog SHOW CREATE TABLE / catalog state against a MySQL 8.0 container via testcontainers
> Reference sources: MySQL 8.0 Reference Manual, real MySQL 8.0 container behavior, existing container_scenarios_test.go patterns

Status: [ ] pending, [x] passing, [~] partial (needs upstream change)

---

## Phase 1: Multi-Command ALTER TABLE

### 1.1 Column Repositioning Interactions

- [x] ADD COLUMN x AFTER a, ADD COLUMN y AFTER x — second command references column added by first
- [x] ADD COLUMN x FIRST, ADD COLUMN y FIRST — both FIRST, y should end up before x
- [x] ADD COLUMN x AFTER a, DROP COLUMN a — add after a column that is then dropped in same statement
- [x] MODIFY COLUMN a INT AFTER c, MODIFY COLUMN b INT AFTER a — chain of AFTER references
- [x] MODIFY COLUMN a INT FIRST, ADD COLUMN x INT AFTER a — FIRST then AFTER the moved column
- [x] DROP COLUMN a, ADD COLUMN a INT — drop then re-add same column name
- [x] CHANGE COLUMN a b INT, ADD COLUMN c INT AFTER b — rename then reference new name
- [x] ADD COLUMN x INT, ADD COLUMN y INT, ADD COLUMN z INT — three appends, verify final order

### 1.2 Column + Index Interactions

- [x] ADD COLUMN x INT, ADD INDEX idx_x (x) — add column then index on it in same ALTER
- [x] ADD COLUMN x INT, ADD UNIQUE INDEX ux (x) — add column then unique index
- [x] DROP COLUMN x, DROP INDEX idx_x — drop column and its index simultaneously
- [x] MODIFY COLUMN x VARCHAR(200), DROP INDEX idx_x, ADD INDEX idx_x (x) — rebuild index after type change
- [x] CHANGE COLUMN x y INT, ADD INDEX idx_y (y) — rename column then index with new name
- [x] ADD COLUMN x INT, ADD PRIMARY KEY (id, x) — add column then include in new PK
- [x] DROP INDEX idx_x, ADD INDEX idx_x (x, y) — drop and recreate index with extra column

### 1.3 Column + FK Interactions

- [x] ADD COLUMN parent_id INT, ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id) — add column then FK
- [x] ADD COLUMN parent_id INT, ADD INDEX idx (parent_id), ADD CONSTRAINT fk FOREIGN KEY (parent_id) REFERENCES parent(id) — explicit index before FK
- [x] DROP FOREIGN KEY fk, DROP COLUMN parent_id — drop FK then its column
- [x] DROP FOREIGN KEY fk, DROP INDEX fk — drop FK then its backing index
- [x] ADD FOREIGN KEY fk1 (...), ADD INDEX idx (...) on same column in same ALTER — FK uses explicit index

### 1.4 Error Semantics in Multi-Command ALTER

- [x] ADD COLUMN x INT, ADD COLUMN x INT — duplicate column in same ALTER, error on second command
- [x] ADD COLUMN x INT, DROP COLUMN nonexistent — first succeeds, second errors; verify x was NOT added (MySQL rolls back entire ALTER)
- [x] MODIFY COLUMN nonexistent INT, ADD COLUMN y INT — first errors, second never runs
- [x] DROP INDEX nonexistent, ADD COLUMN y INT — error on first, verify y not added
- [x] ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x — add, index, then drop same column

---

## Phase 2: Foreign Key Lifecycle

### 2.1 FK Backing Index Management

- [x] CREATE TABLE with named FK, no explicit index — implicit index uses constraint name
- [x] CREATE TABLE with unnamed FK — implicit index uses first column name
- [x] CREATE TABLE with explicit index on FK columns, then named FK — no duplicate index created
- [x] CREATE TABLE with FK on column already in UNIQUE KEY — no duplicate index created
- [x] CREATE TABLE with FK on column already in PRIMARY KEY — no duplicate index created
- [x] CREATE TABLE with multi-column FK, partial index exists — implicit index still created (partial doesn't cover)
- [x] ALTER TABLE ADD FK when column already has index — no duplicate index
- [x] ALTER TABLE ADD FK when column has no index — implicit index created
- [x] ALTER TABLE DROP FOREIGN KEY — FK removed but backing index remains (MySQL behavior)
- [x] ALTER TABLE DROP FOREIGN KEY, DROP INDEX fk_name — explicit index cleanup after FK drop

### 2.2 FK Validation Matrix

- [x] DROP TABLE parent when child FK exists, foreign_key_checks=1 — error 3730
- [x] DROP TABLE parent when child FK exists, foreign_key_checks=0 — succeeds, child FK becomes orphan
- [x] DROP TABLE child then parent, foreign_key_checks=1 — succeeds (child dropped first)
- [x] DROP COLUMN used in FK on same table, foreign_key_checks=1 — error 1828
- [x] CREATE TABLE with FK referencing nonexistent table, foreign_key_checks=0 — succeeds
- [x] CREATE TABLE with FK referencing nonexistent column, foreign_key_checks=0 — succeeds
- [x] ALTER TABLE ADD FK with type mismatch (INT vs VARCHAR), foreign_key_checks=1 — error
- [x] ALTER TABLE ADD FK where referenced table has no index on referenced columns — error 1822
- [x] SET foreign_key_checks=0 then CREATE circular FKs then SET foreign_key_checks=1 — both tables valid
- [x] Self-referencing FK (table references itself) — column references own table PK
- [x] FK column count mismatch (single-column FK referencing composite PK) — error
- [x] Cross-database FK: FOREIGN KEY (col) REFERENCES other_db.parent(id) — stored and rendered correctly

### 2.3 FK Actions Rendering

- [x] FK with ON DELETE CASCADE ON UPDATE SET NULL — both actions rendered in SHOW CREATE
- [x] FK with no action specified — defaults rendered correctly (RESTRICT/NO ACTION omitted in SHOW CREATE)
- [x] Multi-column FK with actions — actions on composite FK rendered correctly

---

## Phase 3: Partition Operations

### 3.1 RANGE Partitioning

- [x] CREATE TABLE PARTITION BY RANGE (expr) with 3 partitions + MAXVALUE — SHOW CREATE matches MySQL
- [x] CREATE TABLE PARTITION BY RANGE COLUMNS (col) — single column, SHOW CREATE matches
- [x] CREATE TABLE PARTITION BY RANGE COLUMNS (col1, col2) — multi-column range
- [x] ALTER TABLE ADD PARTITION to RANGE table — new partition appended
- [x] ALTER TABLE DROP PARTITION from RANGE table — partition removed, others unchanged
- [x] ALTER TABLE REORGANIZE PARTITION p3 INTO (p3a, p3b) — split one range into two
- [x] ALTER TABLE REORGANIZE PARTITION p1, p2 INTO (p_merged) — merge two ranges into one
- [x] RANGE partition with date expression (YEAR(col)) — expression rendering

### 3.2 LIST Partitioning

- [x] CREATE TABLE PARTITION BY LIST (expr) with VALUES IN — SHOW CREATE matches MySQL
- [x] CREATE TABLE PARTITION BY LIST COLUMNS (col) — single column list
- [~] CREATE TABLE PARTITION BY LIST COLUMNS (col1, col2) — multi-column list (parser lacks tuple syntax)
- [x] ALTER TABLE ADD PARTITION with new VALUES IN — partition added
- [x] ALTER TABLE DROP PARTITION from LIST table — partition removed
- [x] ALTER TABLE REORGANIZE PARTITION in LIST table — values redistributed

### 3.3 HASH and KEY Partitioning

- [x] CREATE TABLE PARTITION BY HASH (expr) PARTITIONS 4 — implicit partition defs
- [x] CREATE TABLE PARTITION BY LINEAR HASH (expr) PARTITIONS 4 — LINEAR keyword rendered
- [x] CREATE TABLE PARTITION BY KEY (col) PARTITIONS 4 — KEY partition
- [~] CREATE TABLE PARTITION BY KEY () PARTITIONS 4 — KEY with empty column list (parser gap)
- [x] CREATE TABLE PARTITION BY KEY (col) ALGORITHM=2 PARTITIONS 4 — ALGORITHM rendered
- [x] ALTER TABLE COALESCE PARTITION 2 on HASH table with 4 partitions — reduces to 2
- [x] ALTER TABLE COALESCE PARTITION on KEY table — same behavior as HASH
- [x] ALTER TABLE ADD PARTITION on HASH table — error in MySQL (HASH does not support ADD PARTITION)
- [x] ALTER TABLE REMOVE PARTITIONING on partitioned table — table becomes unpartitioned

### 3.4 Subpartitions and Partition Options

- [x] CREATE TABLE RANGE partitions with SUBPARTITION BY HASH — subpartition rendering
- [x] CREATE TABLE RANGE partitions with SUBPARTITION BY KEY — KEY subpartitions
- [x] Explicit subpartition definitions with names — SUBPARTITION sp1, sp2
- [x] Partition with ENGINE option — ENGINE=InnoDB per partition
- [x] Partition with COMMENT option — COMMENT='desc' per partition
- [x] SUBPARTITIONS N count without explicit defs — count rendering
- [x] ALTER TABLE TRUNCATE PARTITION p1 — no structural change, SHOW CREATE unchanged
- [x] ALTER TABLE EXCHANGE PARTITION p1 WITH TABLE t2 — validation only, structure unchanged

---

## Phase 4: Character Set and Collation

### 4.1 Inheritance Chain

- [x] CREATE DATABASE with CHARSET utf8mb4 → CREATE TABLE (no charset) → column inherits utf8mb4
- [x] CREATE DATABASE with CHARSET latin1 → CREATE TABLE CHARSET utf8mb4 → column inherits utf8mb4 (table overrides DB)
- [x] Table with CHARSET utf8mb4 → ADD COLUMN VARCHAR (no charset) → column inherits table charset
- [x] Table with CHARSET utf8mb4 → ADD COLUMN VARCHAR CHARSET latin1 → column overrides
- [x] CREATE DATABASE CHARSET utf8mb4 → table inherits → SHOW CREATE TABLE shows DEFAULT CHARSET=utf8mb4
- [x] CREATE TABLE with CHARSET only (no COLLATE) — default collation derived from charset
- [x] CREATE TABLE with COLLATE only (no CHARSET) — charset derived from collation

### 4.2 ALTER CHARACTER SET and CONVERT

- [x] ALTER TABLE DEFAULT CHARACTER SET latin1 — table default changes but existing column charsets unchanged
- [x] CONVERT TO CHARACTER SET utf8mb4 — table + all VARCHAR/TEXT/ENUM columns updated
- [x] CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci — non-default collation
- [x] CONVERT TO CHARACTER SET on table with INT columns — INT columns unchanged
- [x] CONVERT TO CHARACTER SET on table with mixed column types — only string types updated
- [x] CONVERT when column already has explicit charset — column charset overwritten
- [x] CONVERT then SHOW CREATE TABLE — column charsets not shown when matching table default

### 4.3 SHOW CREATE TABLE Charset Rendering

- [x] Column charset same as table — CHARACTER SET not shown in column def
- [x] Column charset differs from table — CHARACTER SET shown in column def
- [x] Column collation is non-default for its charset but same as table — COLLATE shown
- [x] Column collation differs from table — both CHARACTER SET and COLLATE shown
- [x] Table with utf8mb4 default collation — COLLATE always shown for utf8mb4 (MySQL 8.0 behavior)
- [x] Table with latin1 and default collation — COLLATE not shown (latin1_swedish_ci is default)
- [x] BINARY charset on column — rendered as CHARACTER SET binary

---

## Phase 5: Generated Columns

### 5.1 Generated Column CRUD

- [x] CREATE TABLE with VIRTUAL generated column (arithmetic) — SHOW CREATE matches
- [x] CREATE TABLE with STORED generated column (arithmetic) — SHOW CREATE matches
- [x] CREATE TABLE with VIRTUAL generated column (string function: CONCAT) — expression rendering
- [x] CREATE TABLE with STORED generated column + NOT NULL — NOT NULL after STORED keyword
- [x] CREATE TABLE with generated column + COMMENT — COMMENT rendered after STORED/VIRTUAL
- [x] CREATE TABLE with generated column + INVISIBLE — INVISIBLE rendered
- [x] CREATE TABLE with generated column using JSON_EXTRACT — JSON expression rendering
- [x] ALTER TABLE ADD COLUMN with GENERATED ALWAYS AS — column added with generation info
- [x] MODIFY COLUMN to change generated expression — expression updated
- [x] ALTER TABLE MODIFY generated column to regular column — generation info removed
- [x] ALTER TABLE MODIFY VIRTUAL to STORED — MySQL 8.0 error (cannot change storage type in-place)

### 5.2 Generated Column Dependencies

- [x] DROP COLUMN referenced by VIRTUAL generated column — MySQL 8.0 error (verify omni behavior)
- [x] DROP COLUMN referenced by STORED generated column — MySQL 8.0 error (verify omni behavior)
- [x] MODIFY base column type when generated column uses it — expression preserved, no validation error
- [x] Generated column referencing another generated column — verify creation and SHOW CREATE
- [x] Index on generated column — index created, rendered correctly
- [x] UNIQUE index on generated column — unique constraint on generated column

---

## Phase 6: CREATE TABLE LIKE

### 6.1 Basic LIKE Completeness

- [x] LIKE copies all column definitions (name, type, nullability, default, comment) — verify each attribute
- [x] LIKE copies PRIMARY KEY — constraint and index present
- [x] LIKE copies UNIQUE KEYs — all unique indexes copied
- [x] LIKE copies regular indexes — all non-unique indexes copied
- [x] LIKE copies FULLTEXT indexes — fulltext index copied
- [x] LIKE copies CHECK constraints — check constraint copied with expression
- [x] LIKE does NOT copy FOREIGN KEY constraints — MySQL 8.0 behavior, FKs excluded
- [x] LIKE copies AUTO_INCREMENT column attribute — but counter resets to 0
- [x] LIKE copies ENGINE, CHARSET, COLLATION, COMMENT — table options preserved

### 6.2 LIKE Edge Cases

- [x] LIKE copies generated columns — expression and VIRTUAL/STORED preserved
- [x] LIKE copies INVISIBLE columns — invisibility preserved
- [x] LIKE does NOT copy partitioning — target table is unpartitioned (MySQL behavior)
- [x] LIKE from table with prefix index — prefix length preserved
- [x] LIKE into TEMPORARY TABLE — target is temporary
- [x] LIKE cross-database — source in different database
- [x] LIKE then ALTER TABLE ADD COLUMN — verify table is independently modifiable

---

## Phase 7: Clone and Catalog Isolation

### 7.1 Catalog State Isolation

- [ ] Execute DDL on catalog A, verify catalog B (separate New()) is unaffected
- [ ] Create table, get reference, execute DROP — reference should not be reusable on new table with same name
- [ ] Two Exec() calls building up state incrementally — state accumulates correctly
- [ ] Exec with ContinueOnError — successful statements applied, failed ones not
- [ ] Exec empty string — no state change, no error
- [ ] Exec with only DML — no state change, statements skipped
- [ ] DROP DATABASE cascade — all tables, views, routines, triggers, events removed
- [ ] Operations across two databases — CREATE TABLE in db1 and db2 independently

### 7.2 Table State Consistency

- [ ] After ADD COLUMN, all column positions are sequential (1-based, no gaps)
- [ ] After DROP COLUMN, remaining column positions are resequenced
- [ ] After MODIFY COLUMN FIRST, positions reflect new order
- [ ] After RENAME COLUMN, index column references updated
- [ ] After DROP INDEX, remaining indexes unaffected
- [ ] colByName index is consistent after every ALTER TABLE operation

---

## Phase 8: Index Edge Cases

### 8.1 Prefix and Expression Index Rendering

- [x] KEY idx (col(10)) — prefix length rendered in SHOW CREATE
- [x] KEY idx (col1(10), col2(20)) — multi-column prefix index
- [x] KEY idx (col(10), col2) — mixed prefix and full column
- [x] KEY idx ((UPPER(col))) — expression index with function
- [x] KEY idx ((col1 + col2)) — expression index with arithmetic
- [x] UNIQUE KEY idx ((UPPER(col))) — unique expression index
- [x] KEY idx (col1, (UPPER(col2))) — mixed regular and expression columns
- [x] KEY idx (col DESC) — descending index column

### 8.2 Index Rendering in SHOW CREATE TABLE

- [ ] Table with PK + UNIQUE + KEY + FULLTEXT + SPATIAL — verify rendering order matches MySQL 8.0
- [ ] Table with regular index + expression index — verify relative ordering matches MySQL 8.0
- [ ] Index with COMMENT — COMMENT rendered in index definition
- [ ] Index with INVISIBLE — rendered with version comment /*!80000 INVISIBLE */
- [ ] Index with KEY_BLOCK_SIZE — rendered in index definition
- [ ] Index with USING BTREE — rendered when non-default for index type
- [ ] Index with USING HASH — rendered for MEMORY engine indexes

---

## Phase 9: SET Variables and Miscellaneous

### 9.1 SET Variable Effects

- [x] SET foreign_key_checks = 0 then CREATE TABLE with invalid FK — succeeds
- [x] SET foreign_key_checks = 1 then CREATE TABLE with invalid FK — fails
- [x] SET foreign_key_checks = OFF — accepts OFF as 0
- [x] SET NAMES utf8mb4 — silently accepted, no catalog state change
- [x] SET CHARACTER SET latin1 — silently accepted
- [x] SET sql_mode = 'STRICT_TRANS_TABLES' — silently accepted
- [x] SET unknown_variable = 'value' — silently accepted (MySQL session variable)

### 9.2 SHOW CREATE TABLE Fidelity

- [ ] Table with no explicit options — ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 rendered
- [ ] Table with ROW_FORMAT=DYNAMIC — ROW_FORMAT rendered
- [ ] Table with KEY_BLOCK_SIZE=8 — KEY_BLOCK_SIZE rendered
- [ ] Table with COMMENT='description' — COMMENT rendered with proper quoting
- [ ] Table with AUTO_INCREMENT=1000 — AUTO_INCREMENT rendered
- [ ] TEMPORARY TABLE — SHOW CREATE TABLE works (no special rendering difference)
- [ ] Table with all column types (INT, BIGINT, DECIMAL, VARCHAR, TEXT, BLOB, JSON, ENUM, SET, DATE, DATETIME, TIMESTAMP) — all types rendered correctly
- [ ] Column with ON UPDATE CURRENT_TIMESTAMP — rendered correctly
- [ ] Column DEFAULT expression (CURRENT_TIMESTAMP, literal, NULL) — all default types rendered
