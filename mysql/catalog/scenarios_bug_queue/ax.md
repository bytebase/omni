# Section AX — ALTER TABLE sub-command bugs

Bugs discovered while running `TestScenario_AX` against MySQL 8.0 +
the omni catalog. Entries are append-only. When a fix lands, do not
delete the entry — annotate it with the fixing commit.

## AX.3 DROP COLUMN does not reject removal of the last column

- **Discovered**: 2026-04-13
- **Section**: AX
- **Scenario**: `ALTER TABLE t DROP COLUMN a;` on a table whose only
  column is `a`. MySQL rejects with
  `ER_CANT_REMOVE_ALL_FIELDS (1090)` at prepare time.
- **Expected (MySQL 8.0.45)**: ALTER fails with
  `ER_CANT_REMOVE_ALL_FIELDS` before any sub-command is applied.
- **Actual (omni)**: catalog silently executes the DROP and leaves
  the table with zero columns. The `ExecResult.Error` is nil.
- **Severity**: MED
- **Fix hint**: `mysql/catalog/alter_table.go` (or wherever
  `processAlterTable` / `applyDropColumn` lives) — after collecting
  drops + adds, compute the resulting column count and raise
  `ER_CANT_REMOVE_ALL_FIELDS` if it would be zero. Check happens
  BEFORE any sub-command is applied so that
  `ALTER t DROP a, ADD b INT` also fails (MySQL's rejection is
  prepare-time, not final-shape-based).

## AX.9 Column-level `REFERENCES` in CREATE TABLE is not silently ignored

- **Discovered**: 2026-04-13
- **Section**: AX
- **Scenario**:
  ```sql
  CREATE TABLE parent (id INT PRIMARY KEY);
  CREATE TABLE t2 (a INT REFERENCES parent(id));
  ```
- **Expected (MySQL 8.0.45)**: the column-level `REFERENCES` clause
  is parsed but produces NO foreign key (notorious InnoDB pitfall).
  `information_schema.KEY_COLUMN_USAGE` reports zero FKs for `t2`.
- **Actual (omni)**: catalog creates a real foreign key constraint
  from the column-level shorthand, so `t2` ends up with one FK where
  MySQL has none. This will generate spurious FK names and wrong
  SDL diffs on any CREATE TABLE that uses the column-level form.
- **Severity**: HIGH
- **Fix hint**: `mysql/catalog/create_table.go` column-level FK
  handling — drop the column-level `REFERENCES` clause for
  InnoDB/MyISAM without emitting a constraint. Only table-level
  `FOREIGN KEY (...) REFERENCES ...` (or `CONSTRAINT ... FOREIGN
  KEY`) should create a catalog FK. See MySQL ref manual
  create-table-foreign-keys.html: "MySQL accepts REFERENCES clauses
  only when specified as part of a separate FOREIGN KEY
  specification."

## Notes on scenario-doc discrepancies (NOT omni bugs)

The following AX scenarios had incorrect `Expected` text in the
starmap doc; the actual MySQL 8.0 behavior was observed via the
container oracle and both MySQL and omni agree. The test was
written against observed oracle behavior. These are tracked here so
the driver can update SCENARIOS-mysql-implicit-behavior.md during
the doc-reconciliation pass.

- **AX.2**: scenario claims `idx_ab` (composite) is dropped entirely
  when column `a` is dropped. Observed MySQL 8.0 behavior: the
  composite index is NOT dropped — the dropped column is stripped
  from its keypart list and the index remains over the surviving
  column(s). Both oracle and omni report `idx_ab, idx_bc` after
  `DROP COLUMN a`.
- **AX.4**: scenario claims MySQL rejects `DROP COLUMN` referenced
  by a `CHECK` constraint. Observed: MySQL 8.0 in this container
  permits the drop (CHECK expression is removed/invalidated
  silently). The GENERATED-column case still errors as expected.
  The test now asserts omni matches whatever oracle does.
- **AX.15**: scenario's "single-pass" wording implies
  `ADD COLUMN c INT AFTER a, RENAME COLUMN a TO aa` in the same
  ALTER should resolve `AFTER a` against the old schema. Observed
  MySQL 8.0 behavior: ALTER errors with `Unknown column 'a' in 't'`
  because the rename is processed before the AFTER clause is
  resolved. The test was rewritten to use `ADD COLUMN c INT` (no
  AFTER clause) to exercise the drop/add/rename composition that
  actually works.
