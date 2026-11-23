# MariaDB stored-program catalog simulation

`simulate.tsv` holds stored procedures, stored functions, triggers and events
written as an application writes them, judged against the pinned MariaDB
(`../../../testdata/pins.tsv`, `mariadb:11.8.8`), run across five fixtures in
`../fixtures/`:

- `one_table@missing` — base table `t` (no object named `a`)
- `one_table@table` — base table `a`
- `objects@routine` — stored procedure `a` and stored function `a`
- `objects@trigger` — base table `t` with trigger `a` (`BEFORE INSERT`)
- `objects@event` — recurring event `a`

The compat.tsv statements on these objects are in `../simulate.tsv` already.
These rows carry `mysql/catalog`'s programs corpus to MariaDB with real names,
values and bodies, and add MariaDB's own forms: `CREATE OR REPLACE` and its
refusal beside `IF NOT EXISTS` (1221), `IF NOT EXISTS` on routines and
triggers, parameter modes and `DEFAULT` values on a function's parameters,
`ROW` and `TYPE OF` parameters and return types, `FOR` loops over a range or a
cursor, stored `AGGREGATE FUNCTION`s with `FETCH GROUP NEXT ROW`, the errors
MariaDB's parser raises in a stored program's body, `FOLLOWS`/`PRECEDES`, and
what a table's `DROP` and `RENAME` do to its triggers. Packages are left out:
the catalog models none.

`TestProgramSimulation` (`../../simulate_program_test.go`) replays MariaDB's
side from `../tapes/TestProgramSimulation.tape` with no server, and
`mariadb/review`'s `TestWalkSimulation` checks every row offline (`programs`
subtest).

**1,280 rows (256 statements × 5 fixtures): 1,275 agree, 5 disagree (99.61%).**

The 5 are `CREATE PROCEDURE a() SET zzz = 1` on each fixture. MariaDB looks a
name a body's `SET` assigns up while it parses: no variable of the program, so
a system variable, and there is none named `zzz` (1193). The catalog keeps no
list of MariaDB's system variables (`mariadb/validate`'s is MySQL 8.0's), so it
takes the procedure.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to MariaDB's outcome:

```bash
go test ./mariadb/catalog/ -run '^TestProgramSimulation$' -update-compat
```

Record every row on the pinned server again:

```bash
OMNI_RECORD=1 go test ./mariadb/catalog/ -run '^TestProgramSimulation$'
```

Verify the recorded answers against the live pinned server without writing:

```bash
OMNI_RECORD=verify go test ./mariadb/catalog/ -run '^TestProgramSimulation$'
```

Show both schemas for matching rows:

```bash
go test -v ./mariadb/catalog/ -run '^TestProgramSimulation$' -args -simulate-show='ALTER EVENT'
```
