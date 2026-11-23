# MySQL catalog simulation

`simulate.tsv` holds, for each table, index and view DDL statement on a fixture
(`fixtures/`), what `mysql/catalog` does and what the pinned MySQL
(`../../testdata/pins.tsv`, `mysql:8.0.46`) did: the error code, or `ok:` and a
hash of the stored `SHOW CREATE` texts. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` without a
server; `mysql/review`'s `TestWalkSimulation` holds WALK_THROUGH to the same rows.

**3,218 rows: 3,218 agree, 0 disagree (100.00%).**

- **2,877 generated rows:** the 411 statements MySQL accepts in
  `../../parser/testdata/compat.tsv`, on seven fixtures that name their objects
  `utf8mb4` as compat.tsv does (`one_table@{missing,table,column,index}`,
  `objects@{view,partitioned,fk}`).
- **341 family rows** (`../simulate_families_test.go`), labeled `<family>:<case>`,
  on fixtures named as an application names its objects:

| fixture | rows | covers |
| --- | ---: | --- |
| `partitioned@range` | 54 | RANGE and RANGE COLUMNS: ADD, DROP, REORGANIZE, EXCHANGE, TRUNCATE, COALESCE, REMOVE PARTITIONING, repartitioning |
| `partitioned@list` | 24 | LIST and LIST COLUMNS: values, duplicates, reorganizing |
| `partitioned@hash` | 35 | HASH, LINEAR HASH and KEY: ADD/COALESCE PARTITION, column types |
| `features@generated` | 67 | generated columns, functional and invisible indexes, invisible columns, CHECK constraints |
| `features@check` | 11 | dropping a column a CHECK constraint reads, alone or with the other columns it reads |
| `features@charset` | 29 | CONVERT TO, table and column character sets, VARCHAR length and row size |
| `features@fk` | 45 | foreign keys: add, drop, renames, type changes, TRUNCATE, engines |
| `features@spatial` | 27 | SPATIAL keys and SRIDs, FULLTEXT keys |
| `schemas@cross` | 15 | RENAME across schemas, references into another schema |
| `one_table@missing` | 34 | ALTER TABLE, CREATE INDEX and DROP INDEX ALGORITHM and LOCK |

`mysql84/simulate.tsv` holds the same rows on `mysql:8.4.13`, and
`programs/simulate.tsv` the stored-program rows; see their READMEs.
`queries/simulate.tsv` wraps compat.tsv's queries in a view on `one_table@column`.

## Commands

From `<repo root>`. Commit rows a catalog fix moved to MySQL's outcome:

```bash
go test ./mysql/catalog/ -run '^TestSimulation$' -update-compat
```

Record all rows against pinned `mysql:8.0.46` (a new row or fixture needs it):

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestSimulation$'
```

Verify recorded answers against pinned `mysql:8.0.46` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mysql/catalog/ -run '^TestSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./mysql/catalog/ -run '^TestSimulation$' -args -simulate-show='partition:range'
```
