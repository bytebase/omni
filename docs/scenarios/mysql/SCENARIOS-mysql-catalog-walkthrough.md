# MySQL Catalog Walk-Through Test Coverage

> Goal: Comprehensive test coverage for MySQL catalog Exec() to support bytebase walk-through migration — verify DDL execution produces correct catalog state, emits correct error codes, handles ContinueOnError, and works with realistic multi-step migration scripts.
> Verification: `go test ./mysql/catalog/ -run TestWalkThrough -v` — all scenarios are Go test cases
> Reference sources: PG catalog exec_test.go + round tests, MySQL 8.0 error codes, bytebase walk_through_omni.go error mapping, mysql/catalog/errors.go, mysql/catalog/exec.go processUtility switch

Status: [ ] pending, [x] passing, [~] partial

---

## Phase 1: Exec Mechanics

Foundation tests for the Exec() function itself — error handling modes, result metadata, and edge cases. These must pass before state verification tests make sense.

### 1.1 Exec Result Basics

- [ ] Empty SQL returns nil results, nil error
- [ ] Whitespace-only SQL returns nil results, nil error
- [ ] Comment-only SQL returns nil results, nil error
- [ ] Single DDL statement returns 1 result with no error
- [ ] Multiple DDL statements return correct number of results
- [ ] Result.Index matches statement position (0-based)
- [ ] DML statements (SELECT, INSERT, UPDATE, DELETE) have Skipped=true
- [ ] DML statements do not modify catalog state
- [ ] Unknown/unsupported statements (FLUSH, ANALYZE) return nil error (silently ignored)

### 1.2 Error Propagation and ContinueOnError

- [ ] Default mode: execution stops at first error
- [ ] Default mode: statements after error are not executed
- [ ] Default mode: catalog state reflects only pre-error statements
- [ ] ContinueOnError: all statements are attempted regardless of failures
- [ ] ContinueOnError: successful statements after a failure modify catalog
- [ ] ContinueOnError: Result.Error is set per-statement (nil for success, error for failure)
- [ ] ContinueOnError: multiple errors collected across results
- [ ] Parse error returns top-level error (not per-result), nil results
- [ ] DELIMITER-containing SQL: statements split correctly, results match

### 1.3 Result Metadata (Line, SQL)

- [ ] Single-line multi-statement: each Result.Line is 1
- [ ] Multi-line statements: Result.Line matches first line of each statement
- [ ] DELIMITER mode: Result.Line points to correct line in original SQL
- [ ] Statements after blank lines: Line numbers account for blank lines
- [ ] Result.Line for DML (skipped) statements is still correct

---

## Phase 2: Error Codes

Every error code in mysql/catalog/errors.go must be produced by the correct operation. bytebase maps these to user-facing advice. Tests assert error type (`*catalog.Error`) and `.Code` value.

### 2.1 Database Errors

- [ ] CREATE DATABASE duplicate → ErrDupDatabase (1007)
- [ ] CREATE DATABASE IF NOT EXISTS on existing → no error (not duplicate)
- [ ] DROP DATABASE unknown → ErrUnknownDatabase (1049)
- [ ] DROP DATABASE IF EXISTS on unknown → no error
- [ ] USE unknown database → ErrUnknownDatabase (1049)
- [ ] CREATE TABLE without USE → ErrNoDatabaseSelected (1046)
- [ ] ALTER TABLE without USE → ErrNoDatabaseSelected (1046)
- [ ] ALTER DATABASE unknown → ErrUnknownDatabase (1049)
- [ ] TRUNCATE unknown table → ErrNoSuchTable (1146)

### 2.2 Table and Column Errors

- [ ] CREATE TABLE duplicate name → ErrDupTable (1050)
- [ ] CREATE TABLE IF NOT EXISTS on existing → no error
- [ ] CREATE TABLE with duplicate column names → ErrDupColumn (1060)
- [ ] CREATE TABLE with multiple PRIMARY KEYs → ErrMultiplePriKey (1068)
- [ ] DROP TABLE unknown → ErrNoSuchTable (1146) or ErrUnknownTable (1051)
- [ ] DROP TABLE IF EXISTS on unknown → no error
- [ ] ALTER TABLE on unknown table → ErrNoSuchTable (1146)
- [ ] ALTER TABLE ADD COLUMN duplicate name → ErrDupColumn (1060)
- [ ] ALTER TABLE DROP COLUMN unknown → ErrNoSuchColumn (1054)
- [ ] ALTER TABLE MODIFY COLUMN unknown → ErrNoSuchColumn (1054)
- [ ] ALTER TABLE CHANGE COLUMN unknown source → ErrNoSuchColumn (1054)
- [ ] ALTER TABLE CHANGE COLUMN to duplicate target → ErrDupColumn (1060)
- [ ] ALTER TABLE ADD PRIMARY KEY when PK exists → ErrMultiplePriKey (1068)
- [ ] ALTER TABLE RENAME COLUMN unknown → ErrNoSuchColumn (1054)
- [ ] RENAME TABLE unknown source → ErrNoSuchTable (1146)
- [ ] RENAME TABLE to existing target → ErrDupTable (1050)
- [ ] DROP TABLE t1, t2 where t2 unknown → error for t2, t1 still dropped

### 2.3 Index and Constraint Errors

- [ ] CREATE INDEX duplicate name → ErrDupKeyName (1061) or ErrDupIndex (1831)
- [ ] CREATE INDEX IF NOT EXISTS on existing → no error
- [ ] DROP INDEX unknown → ErrCantDropKey (1091)
- [ ] ALTER TABLE DROP INDEX unknown → ErrCantDropKey (1091)
- [ ] ALTER TABLE ADD UNIQUE INDEX with duplicate name → ErrDupKeyName (1061)

### 2.4 Foreign Key Errors

- [ ] CREATE TABLE FK references unknown table → ErrFKNoRefTable (1824) (when fk_checks=1)
- [ ] CREATE TABLE FK references unknown column → error (when fk_checks=1)
- [ ] DROP TABLE referenced by FK → ErrFKCannotDropParent (3730) (when fk_checks=1)
- [ ] DROP TABLE referenced by FK with fk_checks=0 → no error
- [ ] ALTER TABLE ADD FK references unknown table → ErrFKNoRefTable (1824)
- [ ] ALTER TABLE DROP COLUMN used in FK → appropriate error
- [ ] ALTER TABLE ADD FK where referenced table lacks matching index → ErrFKMissingIndex (1822)
- [ ] ALTER TABLE ADD FK where column types are incompatible → ErrFKIncompatibleColumns (3780)
- [ ] SET foreign_key_checks=0 disables FK validation
- [ ] SET foreign_key_checks=1 re-enables FK validation

### 2.5 View Errors

- [ ] CREATE VIEW on existing name (not OR REPLACE) → error
- [ ] CREATE OR REPLACE VIEW on existing → no error
- [ ] DROP VIEW unknown → error
- [ ] DROP VIEW IF EXISTS on unknown → no error

### 2.6 Routine, Trigger, and Event Errors

- [ ] CREATE PROCEDURE duplicate → ErrDupProcedure (1304)
- [ ] CREATE FUNCTION duplicate → ErrDupFunction (1304)
- [ ] DROP PROCEDURE unknown → ErrNoSuchProcedure (1305)
- [ ] DROP FUNCTION unknown → ErrNoSuchFunction (1305)
- [ ] DROP PROCEDURE IF EXISTS unknown → no error
- [ ] CREATE TRIGGER duplicate → ErrDupTrigger (1359)
- [ ] DROP TRIGGER unknown → ErrNoSuchTrigger (1360)
- [ ] DROP TRIGGER IF EXISTS unknown → no error
- [ ] CREATE EVENT duplicate → ErrDupEvent (1537)
- [ ] DROP EVENT unknown → ErrNoSuchEvent (1539)
- [ ] DROP EVENT IF EXISTS unknown → no error

---

## Phase 3: State Verification

After DDL execution, the catalog state must accurately reflect all changes. Each scenario executes DDL then inspects the catalog objects.

### 3.1 CREATE TABLE State

- [ ] Table exists in database after CREATE TABLE
- [ ] Column count matches definition
- [ ] Column names match in order
- [ ] Column positions are sequential (1-based)
- [ ] Column types preserved (INT, VARCHAR(100), DECIMAL(10,2), DATETIME, TEXT, BLOB, JSON)
- [ ] NOT NULL constraint reflected in Column.Nullable
- [ ] DEFAULT value preserved in Column.Default
- [ ] AUTO_INCREMENT column flagged correctly
- [ ] Column COMMENT preserved
- [ ] Table ENGINE preserved
- [ ] Table CHARSET and COLLATION preserved
- [ ] Table COMMENT preserved
- [ ] Table AUTO_INCREMENT starting value preserved
- [ ] UNSIGNED modifier in ColumnType
- [ ] Generated column (VIRTUAL) recorded with expression
- [ ] Generated column (STORED) recorded with expression and Stored=true
- [ ] Column visibility (INVISIBLE) preserved
- [ ] CREATE TABLE ... LIKE copies columns, indexes, constraints from source

### 3.2 Database, Rename, and Truncate State

- [ ] ALTER DATABASE changes charset: database.Charset updated
- [ ] ALTER DATABASE changes collation: database.Collation updated
- [ ] RENAME TABLE: old name gone, new name present, same columns/indexes
- [ ] RENAME TABLE cross-database: table moves between databases
- [ ] TRUNCATE TABLE: table still exists, AUTO_INCREMENT reset to 0

### 3.3 Index and Constraint State

- [ ] PRIMARY KEY creates index with Primary=true
- [ ] PRIMARY KEY columns marked NOT NULL automatically
- [ ] UNIQUE KEY creates index with Unique=true and matching constraint
- [ ] Regular INDEX created with correct columns
- [ ] Multi-column index preserves column order
- [ ] FULLTEXT index flagged correctly
- [ ] SPATIAL index flagged correctly
- [ ] Index COMMENT preserved
- [ ] Index visibility (INVISIBLE) preserved
- [ ] FOREIGN KEY constraint records RefTable, RefColumns, OnDelete, OnUpdate
- [ ] FOREIGN KEY creates implicit backing index
- [ ] CHECK constraint records expression and enforcement state
- [ ] Named constraints preserve user-specified names
- [ ] Unnamed constraints get auto-generated names

### 3.4 ALTER TABLE State — Column Operations

- [ ] ADD COLUMN: new column appears at end by default
- [ ] ADD COLUMN FIRST: new column at position 1, others shift
- [ ] ADD COLUMN AFTER x: new column after specified column
- [ ] DROP COLUMN: column removed, remaining positions resequenced
- [ ] MODIFY COLUMN: type change reflected, position unchanged (unless FIRST/AFTER)
- [ ] MODIFY COLUMN: nullability change reflected
- [ ] CHANGE COLUMN: old name removed, new name appears
- [ ] CHANGE COLUMN: type and attributes updated
- [ ] RENAME COLUMN: name changed, position unchanged
- [ ] ALTER COLUMN SET DEFAULT: default value updated
- [ ] ALTER COLUMN DROP DEFAULT: default removed
- [ ] ALTER COLUMN SET VISIBLE / SET INVISIBLE: visibility toggled

### 3.5 ALTER TABLE State — Index and Constraint Operations

- [ ] ADD INDEX: index appears in table.Indexes
- [ ] ADD UNIQUE INDEX: creates both index and constraint
- [ ] ADD PRIMARY KEY: index created with Primary=true
- [ ] ADD FOREIGN KEY: constraint created with correct refs
- [ ] ADD CHECK: constraint created with expression
- [ ] DROP INDEX: index removed from table.Indexes
- [ ] DROP PRIMARY KEY: PK index removed (if supported)
- [ ] RENAME INDEX: name changed in table.Indexes
- [ ] ALTER INDEX SET VISIBLE / INVISIBLE: visibility toggled
- [ ] ALTER CHECK SET ENFORCED / NOT ENFORCED: enforcement toggled
- [ ] CONVERT TO CHARACTER SET: all string columns updated

### 3.6 View State

- [ ] CREATE VIEW: view exists in database.Views
- [ ] CREATE VIEW: Definition stores deparsed SQL
- [ ] CREATE VIEW: Algorithm, Definer, SqlSecurity preserved
- [ ] CREATE VIEW: Column list derived from SELECT
- [ ] CREATE OR REPLACE VIEW: updates existing view
- [ ] ALTER VIEW: updates definition and attributes
- [ ] DROP VIEW: view removed from database.Views
- [ ] View referencing table columns resolves correctly

### 3.7 Routine, Trigger, and Event State

- [ ] CREATE PROCEDURE: exists in database.Procedures with correct params
- [ ] CREATE FUNCTION: exists in database.Functions with return type
- [ ] ALTER PROCEDURE: characteristics updated
- [ ] DROP PROCEDURE: removed from database.Procedures
- [ ] CREATE TRIGGER: exists in database.Triggers with correct timing/event/table
- [ ] DROP TRIGGER: removed from database.Triggers
- [ ] CREATE EVENT: exists in database.Events with schedule and status
- [ ] ALTER EVENT: schedule/status/body updated
- [ ] ALTER EVENT RENAME: event renamed
- [ ] DROP EVENT: removed from database.Events

---

## Phase 4: Multi-Step Migrations

Realistic DDL sequences that simulate bytebase walk-through. Each scenario is a complete migration script that tests cross-object interactions and cumulative state.

### 4.1 Schema Setup Migrations

- [ ] Create database + multiple tables + indexes in one Exec: all objects present
- [ ] Create table then add FK referencing it: FK resolves correctly
- [ ] Create table then create view on it: view resolves columns
- [ ] Create database + SET CHARSET + tables: charset inheritance correct
- [ ] mysqldump-style output: SET vars + DELIMITER + procedures + triggers + tables

### 4.2 Schema Modification Migrations

- [ ] Add column then create index on it: index references new column
- [ ] Rename column then verify index still references correct column name
- [ ] Add column then add FK using it: FK constraint correct
- [ ] Drop index then re-create with different columns: old gone, new present
- [ ] ALTER TABLE multiple sub-commands in sequence: cumulative effect correct
- [ ] RENAME TABLE then CREATE VIEW on new name: view resolves
- [ ] Change column type then verify dependent generated column still recorded
- [ ] CONVERT TO CHARACTER SET then verify all string columns updated

### 4.3 Error Detection in Migrations

- [ ] ContinueOnError migration: first CREATE succeeds, duplicate CREATE fails, third ALTER on first table succeeds
- [ ] ContinueOnError: correct error code for each failing statement
- [ ] Migration with FK cycle: A references B, B references A (with fk_checks=0 then verify after fk_checks=1)
- [ ] DROP TABLE cascade: attempt drop parent with FK child → correct error
- [ ] Migration referencing non-existent table: ALTER on missing table → correct error at correct line
- [ ] Multiple errors in one migration: all error codes correct, all lines correct
- [ ] Mixed DML and DDL: DML skipped, DDL executed, final state correct
