# SQL Server catalog test data

`simulate.tsv` holds one row per core DDL statement that the pinned SQL Server
2022 (`../../testdata/pins.tsv`,
`mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04`, `16.0.4305.5`)
accepts in `mssql/parser/testdata/compat.tsv`, run once on each fixture in
`fixtures/`:

- Four namings of one plain table: `one_table@missing`, `one_table@table`,
  `one_table@column`, `one_table@index`.
- Eight object-kind fixtures: `objects@view`, `objects@fk`, `objects@schema`,
  `objects@sequence`, `objects@synonym`, `objects@procedure`,
  `objects@function`, `objects@trigger`.

Simulated statement families cover database-local catalog DDL on schemas,
tables, indexes (relational, columnstore, XML, selective XML, spatial), views,
sequences, synonyms, stored procedures, functions, and triggers (excluding
server-, database-, security-, partition-, fulltext-, and external-object
statements).

A row holds omni's outcome and SQL Server's — the primary SQL Server error
number (`mssql.Error.Number`) when refused, or `ok:` and a hash of the
normalized DDL objects stored in the database afterwards. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` and starts
no server. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.

**3,072 / 3,072 rows (100.00%) agree** (256 statements × 12 fixtures).

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./mssql/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned server again, which a new row or fixture needs:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/catalog/ -run '^TestSimulation$'
```

Verify that live SQL Server matches the committed tape and `simulate.tsv` with
zero differences:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mssql/catalog/ -run '^TestSimulation$'
```

Show both schemas of the rows that hold a statement's text:

```bash
go test -v ./mssql/catalog/ -run '^TestSimulation$' -args -simulate-show='ALTER SEQUENCE'
```
