# TiDB sequence catalog simulation

`simulate.tsv` holds every generated sequence DDL statement (`CREATE SEQUENCE`,
`ALTER SEQUENCE`, and `DROP SEQUENCE`) that the pinned TiDB
(`../../../testdata/pins.tsv`, `pingcap/tidb:v8.5.5`) accepts in
`tidb/parser/testdata/compat.tsv`, run across five fixtures in `../fixtures/`:

- `one_table@missing` — base table `t` (no object named `ident`)
- `one_table@table` — base table `ident`
- `objects@view` — base table `t` and view `ident`
- `objects@sequence` — default ascending sequence `ident`
- `objects@sequence_options` — descending cyclic sequence `ident` with custom
  `START`, `MINVALUE`, `MAXVALUE`, `INCREMENT`, `CACHE`, and `COMMENT`

Unlike `../simulate.tsv` (which reads only tables and views so committed row
hashes stay stable), `TestSequenceSimulation`
(`../../simulate_sequence_test.go`) reads tables, views, and sequences from TiDB
(`SHOW FULL TABLES` + `SHOW CREATE TABLE`/`VIEW`/`SEQUENCE`), replays from
`../tapes/TestSequenceSimulation.tape` without starting a server, and compares
every object's `SHOW CREATE` text directly against `tidb/catalog`.
`tidb/review`'s `TestWalkSimulation` also checks every row offline
(`sequences` subtest).

**135 rows (27 statements × 5 fixtures): 135 agree, 0 disagree (100.00%).**

## Commands

From `<repo root>`. Commit rows a catalog fix moved to TiDB's outcome:

```bash
go test ./tidb/catalog/ -run '^TestSequenceSimulation$' -update-compat
```

Record all rows against pinned `pingcap/tidb:v8.5.5`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./tidb/catalog/ -run '^TestSequenceSimulation$'
```

Verify recorded answers against pinned `pingcap/tidb:v8.5.5` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./tidb/catalog/ -run '^TestSequenceSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./tidb/catalog/ -run '^TestSequenceSimulation$' -args -simulate-show='ALTER SEQUENCE'
```
