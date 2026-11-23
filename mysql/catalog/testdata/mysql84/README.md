# MySQL 8.4 LTS catalog simulation

`simulate.tsv` holds the full MySQL table, index, and view DDL simulation corpus
(the same 411 statements from `mysql/parser/testdata/compat.tsv` × 7 fixtures as
`../simulate.tsv`, totaling 2,877 generated rows, and the same 341 family rows,
see `../README.md`) plus 7 hand-added MySQL 8.4 LTS-specific rows (`mysql84:*`),
recorded live against the pinned MySQL 8.4 LTS server
(`../../../testdata/pins.tsv`, role `mysql84`: `mysql:8.4.13`).

`TestSimulationMySQL84` (`../../simulate_mysql84_test.go`) replays offline from
`../tapes/TestSimulationMySQL84.tape` without starting a server and compares
`mysql/catalog` configured at `MySQL84` (`LoadSQLWithVersion(..., MySQL84)` +
`NormalizerFor(MySQL84)`) against MySQL 8.4.13's stored `SHOW CREATE` forms
and error codes.

**3,225 rows (2,877 generated + 341 family + 7 `mysql84:*` hand-added): 3,225 agree, 0 disagree (100.00%).**

## MySQL 8.4 LTS behaviors covered

- **`restrict_fk_on_non_standard_key=ON` (default in 8.4 LTS):** InnoDB requires
  parent-side foreign key references to match an exact-order, full-column
  `PRIMARY KEY` or `UNIQUE KEY` with no prefix lengths. Missing/non-unique,
  unindexed, partial-composite, or reordered parent keys fail with error `6125`
  (`ER_FK_NO_UNIQUE_INDEX_PARENT`), even when `foreign_key_checks = 0` against
  an existing parent table, and dropping a parent unique key needed by a foreign
  key fails with error `1553` even when a non-unique fallback index exists.
- **Removal of `AUTO_INCREMENT` on `FLOAT` / `DOUBLE`:** `CREATE TABLE` and
  `ALTER TABLE ... MODIFY` fail with error `1063` (`ER_WRONG_FIELD_SPEC`) on
  8.4 LTS (while remaining accepted on `MySQL80`).
- **Deterministic `KEY_BLOCK_SIZE` recording:** MySQL 8.4 changed the Linux
  default of `innodb_flush_method` from `fsync` to `O_DIRECT`, which causes
  `KEY_BLOCK_SIZE=1` / `2` to fail with `EINVAL` (reported as `3675` disk full)
  only on 4K-sector block devices while succeeding on 512B-sector devices. The
  `mysql84` testcontainer pins `--innodb-flush-method=fsync` so recordings are
  reproducible across host storage hardware.

## Commands

From `<repo root>`. Commit rows a catalog fix moved to MySQL 8.4's outcome:

```bash
go test ./mysql/catalog/ -run '^TestSimulationMySQL84$' -update-compat
```

Record all rows against pinned `mysql:8.4.13`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestSimulationMySQL84$'
```

Verify recorded answers against pinned `mysql:8.4.13` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mysql/catalog/ -run '^TestSimulationMySQL84$'
```

Show both schemas for matching rows:

```bash
go test -v ./mysql/catalog/ -run '^TestSimulationMySQL84$' -args -simulate-show='KEY_BLOCK_SIZE'
```
