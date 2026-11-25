# OceanBase query & expression view simulation

`simulate.tsv` implements Spec §10 Decision D6
(`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`):
every unique generated query (`SELECT`, `WITH`, `TABLE`, `VALUES`) and `DO`
expression accepted by the pinned OceanBase (`../../../testdata/pins.tsv`,
`oceanbase/oceanbase-ce:4.3.5.6-106000012026040916`) in
`oceanbase/parser/testdata/compat.tsv`, wrapped in `CREATE VIEW _v AS ...` on
the `one_table@column` fixture (`CREATE TABLE a (b INT PRIMARY KEY, a INT, d INT, KEY idx_d (d))`).
OceanBase's stored `SHOW CREATE VIEW` text is the engine's own reading of each
query and expression. `TestQuerySimulation` (`../../simulate_query_test.go`)
replays it from `../tapes/TestQuerySimulation.tape` with no server.

4,224 rows (141 engine-accepted views, 4,083 engine refusals):
- **Baseline:** 3,625 / 4,022 agree (90.13%), 397 disagree (0 only engine takes, 186 only omni takes, 55 other schema, 156 other code).
- **Final:** 4,224 / 4,224 agree (100.00%: 141 / 141 stored views matching, 4,083 / 4,083 matching refusals) and 0 disagree.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./oceanbase/catalog/ -run '^TestQuerySimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test -timeout 30m ./oceanbase/catalog/ -run '^TestQuerySimulation$'
```

Verify committed answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test -timeout 30m ./oceanbase/catalog/ -run '^TestQuerySimulation$'
```
