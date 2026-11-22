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

4,021 rows (140 engine-accepted views, 3,881 engine refusals):
- **Baseline:** 3,625 / 4,022 agree (90.13%), 397 disagree (0 only engine takes, 186 only omni takes, 55 other schema, 156 other code).
- **Final:** 3,998 / 4,021 agree (99.43%: 140 / 140 stored views matching, 3,858 / 3,881 matching refusals) and 23 disagree (0 only engine takes, 13 only omni takes, 0 other schema, 10 other code).

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

## What the 23 rows still disagree on

| Rows | Category | Omni vs OceanBase | Why parked |
|---:|---|---|---|
| 13 | Function existence (`STATISTICS()`, `STATISTICS('s')`, `a()`, `a.a()`, `a.STATISTICS()`, `a.CURRENT_USER()`, `COALESCE_AGGR()`, `REPLACE_CONST()`, `DECRYPT()`) | `ok` vs `1305` | Unqualified builtin function existence depends on the engine's full builtin registry, and schema-qualified calls may target unmodeled external stored functions in SDL. |
| 4 | `CALC_PARTITION_ID(a, a)` / `CALC_PARTITION_ID(a, a, a)` | `1054` vs `1210` | OceanBase-specific builtin argument validation precedes column resolution. |
| 4 | `/*+ LOAD_BATCH_SIZE(...) */` hints | `1054` vs `5063` | Optimizer hint name validation (`Unknown hint`) is outside catalog view modeling. |
| 2 | `/*+ DIRECT(TRUE, 1, 's') */` hints | `1054` vs `1210` | Optimizer hint argument validation (`Incorrect arguments to load method in direct hint`) is outside catalog view modeling. |
