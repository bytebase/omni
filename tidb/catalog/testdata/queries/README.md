# TiDB query & expression view simulation

`simulate.tsv` implements Spec §10 Decision D6
(`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):
every unique generated query (`SELECT`, `WITH`, `TABLE`, `VALUES`) and `DO`
expression accepted by the pinned TiDB (`../../../testdata/pins.tsv`,
`pingcap/tidb:v8.5.5`) in `tidb/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on the `one_table@column` fixture. TiDB's stored
`SHOW CREATE VIEW` text is the engine's own reading of each query and
expression. `TestQuerySimulation` (`../../simulate_query_test.go`) replays it
from `../tapes/TestQuerySimulation.tape` with no server.

382 rows: 382 agree (187 stored views matching byte-for-byte, 195 matching
refusals) and 0 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./tidb/catalog/ -run '^TestQuerySimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test ./tidb/catalog/ -run '^TestQuerySimulation$'
```

Verify committed answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test ./tidb/catalog/ -run '^TestQuerySimulation$'
```
