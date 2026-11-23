# MariaDB query & expression view simulation

`simulate.tsv` implements Spec §10 Decision D6
(`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):
every unique generated query (`SELECT`, `WITH`, `TABLE`, `VALUES`) and `DO`
expression accepted by the pinned MariaDB (`../../../testdata/pins.tsv`,
`mariadb:11.8.8`) in `mariadb/parser/testdata/compat.tsv`, wrapped in
`CREATE VIEW _v AS ...` on the `one_table@column` fixture. MariaDB's stored
`SHOW CREATE VIEW` text is the engine's own reading of each query and
expression. `TestQuerySimulation` (`../../simulate_query_test.go`) replays it
from `../tapes/TestQuerySimulation.tape` with no server.

457 rows: **457 / 457 (100.00%) agree** (every stored view matching
byte-for-byte with 0 schema mismatches and 0 disagreements).

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./mariadb/catalog/ -run '^TestQuerySimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test ./mariadb/catalog/ -run '^TestQuerySimulation$'
```

Verify committed answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test ./mariadb/catalog/ -run '^TestQuerySimulation$'
```
