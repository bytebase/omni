# PostgreSQL 18 catalog simulation

`simulate.tsv` holds `pg/catalog`'s DDL simulation suite recorded against the
pinned PostgreSQL 18 server (`../../../testdata/pins.tsv`, role `pg18`:
`postgres:18.7-alpine`, tag `REL_18_7`) and replayed offline from
`../tapes/TestSimulationPG18.tape`.

## Coverage

5,424 rows total (5,391 generated from `pg/parser/testdata/compat.tsv` across
9 core fixtures in `../fixtures/` + 33 curated `pg18/` feature rows):

- `one_table@missing` (611 rows) — object creation, missing targets, PG 18 built-ins (`uuidv7()`, `casefold()`, `array_sort()`), SQL/JSON views, `UNIQUE NULLS NOT DISTINCT`, unlogged sequences, SQL-body functions, named `NOT NULL` constraints
- `one_table@table` (599 rows) — table DDL and PG 18 `ATTACH`/`DETACH PARTITION` `42809` checks on non-partitioned tables
- `one_table@column` (614 rows) — column, constraint, generated expression, identity, `LIKE INCLUDING ALL`, inheritance, and PG 18 cataloged `NOT NULL` constraints (`DROP`/`RENAME`/`COMMENT ON CONSTRAINT`)
- `objects@view` (599 rows) — view DDL
- `objects@matview` (599 rows) — materialized view DDL
- `objects@sequence` (599 rows) — sequence DDL
- `objects@partitioned` (602 rows) — partitioned table DDL, PG 18 `SET LOGGED`/`SET UNLOGGED` `42809` rejection, partition detach PK naming, identity and `NOT NULL` propagation
- `objects@fk` (600 rows) — foreign key and check constraints
- `objects@domain` (601 rows) — domain DDL and PG 18 `ALTER DOMAIN ... ADD [CONSTRAINT ...] NOT NULL NO INHERIT` `0A000` rejection

**5,424 rows: 5,424 agree, 0 disagree (100.00%).**

## PG 17 vs PG 18 catalog differences verified

1. **`ALTER DOMAIN ... ADD [CONSTRAINT name] NOT NULL NO INHERIT`**: rejected with `0A000` before domain lookup in PG 18 (PG 17 accepted it on domains).
2. **`ALTER TABLE <non-partitioned> ATTACH/DETACH PARTITION [FINALIZE]`**: rejected with `42809` (`ERRCODE_WRONG_OBJECT_TYPE`) in PG 18 (PG 17 returned `42P17` / `42P01`).
3. **`ALTER TABLE <partitioned> SET LOGGED / SET UNLOGGED`**: rejected with `42809` in PG 18 (PG 17 silently accepted as a no-op).
4. **Cataloged table `NOT NULL` constraints (`pg_constraint.contype = 'n'`)**: PG 18 tracks `<rel>_<col>_not_null` (or explicit `CONSTRAINT name NOT NULL`) in `pg_constraint`, supporting `RENAME CONSTRAINT`, `COMMENT ON CONSTRAINT`, and `DROP CONSTRAINT` (refusing PK columns with `42P16` and identity columns with `42P17`).
5. **Partition cloned PK constraint naming & SQL/JSON expressions**: cloned partition primary key constraints use `<part>_pkey` (matching both PG 17 and PG 18 on `DETACH PARTITION`), and `JSON_VALUE` / `JSON_QUERY` preserve `RETURNING` type, wrapper, and quotes clauses.

## Commands

From `<repo root>`. Commit rows a catalog fix moved to PostgreSQL 18's outcome:

```bash
go test ./pg/catalog/ -run '^TestSimulationPG18$' -update-compat
```

Record all rows against pinned `postgres:18.7-alpine`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./pg/catalog/ -run '^TestSimulationPG18$'
```

Verify recorded answers against pinned `postgres:18.7-alpine` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./pg/catalog/ -run '^TestSimulationPG18$'
```

Show both schemas for matching rows:

```bash
go test -v ./pg/catalog/ -run '^TestSimulationPG18$' -args -simulate-show='pg18/'
```
