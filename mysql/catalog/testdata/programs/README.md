# MySQL stored-program catalog simulation

`simulate.tsv` holds every stored-program DDL statement (`CREATE`/`ALTER`/`DROP`
for `PROCEDURE`, `FUNCTION`, `TRIGGER`, and `EVENT`, including `CREATE DEFINER`)
that the pinned MySQL (`../../../testdata/pins.tsv`, `mysql:8.0.46`) accepts in
`mysql/parser/testdata/compat.tsv`, run across five fixtures in `../fixtures/`:

- `one_table@missing` — base table `t` (no object named `utf8mb4`)
- `one_table@table` — base table `utf8mb4`
- `objects@routine` — stored procedure `utf8mb4` and stored function `utf8mb4`
- `objects@trigger` — base table `t` with trigger `utf8mb4`
- `objects@event` — recurring scheduled event `utf8mb4`

Loadable UDFs (`CREATE [AGGREGATE] FUNCTION ... SONAME '...'`) are excluded
because they load server-level shared libraries into `mysql.func` outside the
database schema (like `INSTALL PLUGIN`).

Unlike `../simulate.tsv` (which reads only `SHOW FULL TABLES` so committed
table/view hashes stay stable), `TestProgramSimulation`
(`../../simulate_program_test.go`) reads tables, views, stored functions, stored
procedures, triggers, and events from MySQL (`SHOW FULL TABLES` +
`information_schema.{ROUTINES,TRIGGERS,EVENTS}` + `SHOW CREATE ...`), replays
from `../tapes/TestProgramSimulation.tape` without starting a server, and diffs
every object kind against `mysql/catalog`. `mysql/review`'s `TestWalkSimulation`
also verifies all rows offline (`programs` subtest).

**440 rows (83 statements × 5 fixtures, and 25 `event:*` rows from
`../../simulate_families_test.go` on `objects@event` and `one_table@missing`):
440 agree, 0 disagree (100.00%).**

## Commands

From `<repo root>`. Commit rows a catalog fix moved to MySQL's outcome:

```bash
go test ./mysql/catalog/ -run '^TestProgramSimulation$' -update-compat
```

Record all rows against pinned `mysql:8.0.46`:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mysql/catalog/ -run '^TestProgramSimulation$'
```

Verify recorded answers against pinned `mysql:8.0.46` without writing:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mysql/catalog/ -run '^TestProgramSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./mysql/catalog/ -run '^TestProgramSimulation$' -args -simulate-show='ALTER EVENT'
```
