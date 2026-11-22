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

4,022 rows (141 engine-accepted views, 3,881 engine refusals):
- **Baseline:** 3,625 / 4,022 agree (90.13%), 397 disagree (0 only engine takes, 186 only omni takes, 55 other schema, 156 other code).
- **Final:** 3,973 / 4,022 agree (98.78%: 136 / 141 stored views matching, 3,837 / 3,881 matching refusals) and 49 disagree (0 only engine takes, 26 only omni takes, 5 other schema, 18 other code).

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

## What the 49 rows still disagree on

| Rows | Category | Omni vs OceanBase | Why parked |
|---:|---|---|---|
| 13 | Function existence (`STATISTICS()`, `a()`, `a.a()`, `a.CURRENT_USER()`, `COALESCE_AGGR()`, `REPLACE_CONST()`, `DECRYPT()`) | `ok` vs `1305` | Unqualified builtin function existence depends on the engine's full builtin registry, and schema-qualified calls may target unmodeled external stored functions in SDL. |
| 13 | Builtin zero-argument calls (`DUMP()`, `CHARSET()`, `COLLATION()`, `COALESCE()`, `REPEAT()`, `REVERSE()`, `RIGHT()`, `REPLACE()`, `TRUNCATE()`, `FORMAT()`, `NORMAL()`, `PASSWORD()`, `ENCRYPT()`) | `ok` vs `1582` | Native builtin arity validation requires maintaining the full engine builtin signature table across releases. |
| 4 | `STATISTICS(a zzz)`, `STATISTICS(a 's')` | `1054` vs `1583` | Native-function aliased-argument rejection (`Incorrect parameters in the call to native function 'statistics'`) precedes column resolution. |
| 4 | `NTH_VALUE(a, a) OVER a` | `5929` vs `1210` | Builtin argument-type validation for `NTH_VALUE`'s second argument precedes window name lookup (`5929`). |
| 4 | `CALC_PARTITION_ID(a, a)` / `CALC_PARTITION_ID(a, a, a)` | `1054` vs `1210` | OceanBase-specific builtin argument validation precedes column resolution. |
| 4 | `/*+ LOAD_BATCH_SIZE(...) */` hints | `1054` vs `5063` | Optimizer hint name validation (`Unknown hint`) is outside catalog view modeling. |
| 2 | `/*+ DIRECT(TRUE, 1, 's') */` hints | `1054` vs `1210` | Optimizer hint argument validation (`Incorrect arguments to load method in direct hint`) is outside catalog view modeling. |
| 2 | `SELECT TRUE`, `SELECT TRUE zzz` | `ok` vs `ok` | OceanBase stores `TRUE` as `(1 = 1)` in `SHOW CREATE VIEW`, while existing catalog view-name tests (`TestViewNamesAndJoinColumns`) intentionally pin `true`. |
| 2 | `SELECT a FROM a FULL`, `SELECT a FROM a NATURAL JOIN a FULL` | `ok` vs `ok` | OceanBase folds the unquoted keyword table alias `FULL` to lowercase `full` in `SHOW CREATE VIEW`. |
| 1 | `SELECT _utf16le 's'` | `ok` vs `ok` | OceanBase reinterprets raw bytes under UTF-16LE (`'猀'`) when storing the literal in `SHOW CREATE VIEW`. |
