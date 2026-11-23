# MariaDB catalog test data

`simulate.tsv` holds one row per DDL statement of the simulated families that
the pinned MariaDB (`../../testdata/pins.tsv`, `mariadb:11.8.8`) takes in
`mariadb/parser/testdata/compat.tsv`, run once on each fixture in `fixtures/`:

- Four namings of one plain table: `one_table@missing`, `one_table@table`,
  `one_table@column`, `one_table@index`.
- Seven object-kind fixtures: `objects@view`, `objects@partitioned`,
  `objects@fk`, `objects@sequence`, `objects@routine` (stored procedure +
  stored function), `objects@trigger`, `objects@event`.

After them come 20 hand-written `check:drop` rows
(`../simulate_families_test.go`) on `features@check`: dropping columns that a
column CHECK, a one-column and a two-column table CHECK, and a CHECK reading
no column read, and renaming them. Then 39 `exists:` rows on `objects@view`,
`objects@partitioned` and `objects@fk` write `IF [NOT] EXISTS` over names the
fixture holds, which the generated rows, failing earlier, do not reach, and 7
`fk:generated_key` rows on `objects@fk` add a second foreign key over the
column whose key the fixture's foreign key generated.

Then 362 rows carry mysql/catalog's DDL families to MariaDB, written in
MariaDB's syntax, on eight fixtures of their own and two of the ones above:

| fixture | rows | families |
| --- | ---: | --- |
| `partitioned@range` | 54 | partition maintenance, repartitioning and bounds on RANGE and RANGE COLUMNS |
| `partitioned@list` | 24 | the same on LIST and LIST COLUMNS |
| `partitioned@hash` | 35 | the same on HASH, KEY and LINEAR HASH |
| `features@generated` | 63 | generated columns, IGNORED keys, invisible columns, CHECK constraints |
| `features@charset` | 30 | CONVERT TO, table and column character sets, VARCHAR and row lengths |
| `features@fk` | 45 | foreign keys added, dropped, renamed and redefined across three tables |
| `features@spatial` | 27 | SPATIAL and FULLTEXT keys, REF_SYSTEM_ID |
| `schemas@cross` | 15 | RENAME TABLE and views across two schemas |
| `one_table@missing` | 52 | ALTER TABLE and CREATE INDEX under ALGORITHM and LOCK, CREATE EVENT |
| `objects@event` | 17 | ALTER, CREATE (OR REPLACE) and DROP EVENT |

`schemas@cross` keeps a second database, `omni_simulation_other`, which the
next row's build drops only after the simulation database: no row may leave a
table there that references one in the simulation database.

Simulated statement families cover database-local DDL on tables, indexes,
views, sequences, stored routines (`PROCEDURE`/`FUNCTION`, excluding
server-level loadable UDFs with `SONAME`), triggers, and events (excluding
temporary objects, databases, packages, accounts, and servers).

A row holds omni's outcome and MariaDB's — an error number, or `ok:` and a hash
of the `SHOW CREATE` texts the database holds afterwards. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` and starts
no server. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.

**5,906 / 5,906 rows (100%) agree** (498 statements × 11 fixtures, and
428 family rows).

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./mariadb/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned server again, which a new row or fixture needs:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mariadb/catalog/ -run '^TestSimulation$'
```

Verify that live MariaDB matches the committed tape and `simulate.tsv` with
zero differences:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mariadb/catalog/ -run '^TestSimulation$'
```

Show both schemas of the rows that hold a statement's text:

```bash
go test ./mariadb/catalog/ -run '^TestSimulation$' -args -simulate-show='ALTER SEQUENCE'
```
