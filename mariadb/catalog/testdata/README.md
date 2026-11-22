# MariaDB catalog test data

`simulate.tsv` holds one row per DDL statement of the simulated families that
the pinned MariaDB (`../../testdata/pins.tsv`, `mariadb:11.8.8`) takes in
`mariadb/parser/testdata/compat.tsv`, run once on each fixture in `fixtures/`:

- Four namings of one plain table: `one_table@missing`, `one_table@table`,
  `one_table@column`, `one_table@index`.
- Seven object-kind fixtures: `objects@view`, `objects@partitioned`,
  `objects@fk`, `objects@sequence`, `objects@routine` (stored procedure +
  stored function), `objects@trigger`, `objects@event`.

Simulated statement families cover database-local DDL on tables, indexes,
views, sequences, stored routines (`PROCEDURE`/`FUNCTION`, excluding
server-level loadable UDFs with `SONAME`), triggers, and events (excluding
temporary objects, databases, packages, accounts, and servers).

A row holds omni's outcome and MariaDB's — an error number, or `ok:` and a hash
of the `SHOW CREATE` texts the database holds afterwards. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` and starts
no server. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.

**5,038 / 5,049 rows (99.78%) agree** (459 statements × 11 fixtures; 11 rows
disagree on error code for `ALTER TABLE a DROP a.a.a;`, where MariaDB answers
`1102` and omni answers `1064`).

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
