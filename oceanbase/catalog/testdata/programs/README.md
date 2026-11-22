# OceanBase stored-program catalog simulation

`simulate.tsv` holds stored-program DDL statements (`CREATE`/`ALTER`/`DROP` for
`PROCEDURE`, `FUNCTION`, `TRIGGER`, and `EVENT`, including `CREATE DEFINER` and
`ALTER DEFINER`) judged against the pinned OceanBase CE MySQL-mode image
(`../../../testdata/pins.tsv`, `oceanbase/oceanbase-ce:4.3.5.6-106000012026040916`),
run across five fixtures in `../fixtures/`:

- `one_table@missing` — base table `t` (no object named `a`)
- `one_table@table` — base table `a`
- `objects@routine` — stored procedure `a` (`BEGIN END`) and stored function `a`
- `objects@trigger` — base table `t` with trigger `a`
- `objects@event` — recurring scheduled event `a` (`STARTS '2030-01-01 00:00:00'`, future-scheduled so rename paths are exercised)

OceanBase's SQL grammar (`sql_parser_mysql_mode.y`, tracked in
`oceanbase/parser/testdata/compat.tsv`) contains only `DROP FUNCTION`; stored
procedures, stored functions, triggers, and events live in OceanBase's separate
PL grammar (`pl_parser_mysql_mode.y`). `TestProgramSimulation` combines the
non-UDF `compat.tsv` statements (`DROP FUNCTION a`, `DROP FUNCTION IF EXISTS a`)
with a curated matrix of 130 OceanBase stored-program statements (132 statements
total). Loadable UDFs (`CREATE [AGGREGATE] FUNCTION ... SONAME '...'`) are
excluded because they install server-level shared libraries outside the database
schema.

Unlike `../simulate.tsv` (which reads only `SHOW FULL TABLES` so committed
table/view hashes stay stable), `TestProgramSimulation`
(`../../simulate_program_test.go`) reads tables, views, stored functions, stored
procedures, triggers, and events from OceanBase (`SHOW FULL TABLES`,
`information_schema.{ROUTINES,TRIGGERS}` + `SHOW CREATE ...`, and deterministic
reconstruction from `information_schema.EVENTS` without `STARTS` because
OceanBase has no `SHOW CREATE EVENT` and auto-fills `STARTS` with creation time),
replays from `../tapes/TestProgramSimulation.tape` without starting a server,
and diffs every object kind against `oceanbase/catalog`. `oceanbase/review`'s
`TestWalkSimulation` also verifies all rows offline (`programs` subtest) with 0
false positives.

**660 rows (132 statements × 5 fixtures): 660 agree, 0 disagree (100.00%).**

## Commands

From `<repo root>`. Commit rows a catalog fix moved to OceanBase's outcome:

```bash
go test ./oceanbase/catalog/ -run '^TestProgramSimulation$' -update-compat
```

Record all rows against pinned `oceanbase/oceanbase-ce:4.3.5.6-106000012026040916`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./oceanbase/catalog/ -run '^TestProgramSimulation$'
```

Verify recorded answers against pinned `oceanbase/oceanbase-ce:4.3.5.6-106000012026040916` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./oceanbase/catalog/ -run '^TestProgramSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./oceanbase/catalog/ -run '^TestProgramSimulation$' -args -simulate-show='ALTER EVENT'
```
