# TiDB catalog test data

`simulate.tsv` holds one row per DDL statement of the simulated families that
the pinned TiDB (`../../testdata/pins.tsv`, `pingcap/tidb:v8.5.5`) takes in
`tidb/parser/testdata/compat.tsv`, run once on each fixture in `fixtures/`:
four namings of one plain table (`one_table@missing`, `@table`, `@column`,
`@index`) plus three object-kind fixtures (`objects@view`,
`objects@partitioned`, `objects@fk`). A row holds omni's outcome and the
engine's — an error number, or `ok:` and a hash of the `SHOW CREATE` texts the
database holds afterwards. `TestSimulation` (`../simulate_test.go`) replays it
from `tapes/TestSimulation.tape` and starts no server. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.
After them come 20 hand-written `check:drop` rows
(`../simulate_families_test.go`) on `features@check`, which turns
`tidb_enable_check_constraint` on: dropping and renaming columns that a column
CHECK, a one-column and a two-column table CHECK, and a CHECK reading no column
read. 3,170 rows: 3,170 agree and 0 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./tidb/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned server again, which a new row or fixture needs:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestSimulation$'
```

Show both schemas of the rows that hold a statement's text:

```bash
go test ./tidb/catalog/ -run '^TestSimulation$' -args -simulate-show='PARTITION BY'
```

Record `TestSimulationCausesOracle`'s tape after changing what it sends. Each
cause of the file is one case there, answered by the server and by the catalog:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestSimulationCausesOracle$'
```
