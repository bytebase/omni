# PostgreSQL stored-program catalog simulation

`simulate.tsv` holds the statements in `programFeatures`
(`../../simulate_program_test.go`) on `../fixtures/programs@objects.sql`:
SQL and PL/pgSQL functions and procedures with real bodies, their `ALTER` and
`DROP`, what a body depends on, triggers, event triggers, `CALL` and `DO`. The
statements in `quotedTypeFeatures` run on `../fixtures/programs@quoted.sql` and
name its quoted mixed-case enum and table as column, argument and result types;
those in `quotedAliasFeatures` quote SQL-standard type names there (`"integer"`),
which PostgreSQL does not map. The statements in `retypeFeatures` run
`ALTER COLUMN TYPE`, with and without `USING`, on `../fixtures/programs@retype.sql`.
`TestProgramSimulation` replays it from `../tapes/TestProgramSimulation.tape`
(postgres:17.12), and `TestProgramSimulationPG18` replays the same rows in
`../pg18/programs/simulate.tsv` from postgres:18.7. Both pins answer every row
alike. The readout adds the database's event triggers to `../simulate.tsv`'s.
`pg/review`'s `TestWalkSimulation` checks WALK_THROUGH on these rows
(`programs` subtest).

**323 rows: all agree, on both pins.**

The catalog checks what PostgreSQL checks at `CREATE` with
`check_function_bodies` on: a SQL body is parsed and analyzed, and its last
statement must return the declared result; a PL/pgSQL body is compiled (syntax,
assigned variables, declared types, condition names, `RETURN`), with the tables
it names left to run time. A quoted literal takes the type it is assigned or
compared to, so an enum refuses a label it lacks (22P02), and a SQL-standard
body (`RETURN`, `BEGIN ATOMIC`) prints it with that type, as PostgreSQL stores
it. Such a body records what it uses, so dropping or retyping it fails with
2BP01 or 0A000. `LANGUAGE internal` names one of the server's built-in
functions (42883 otherwise).

`DO` and `CALL` run a PL/pgSQL body that is a straight list of statements (no
variables, control flow or handlers): its DDL executes, its DML is analyzed, and
`RAISE EXCEPTION` fails it. `CALL` runs only a procedure without parameters,
`SET` options or `SECURITY DEFINER`. A body the catalog does not run, or one
that calls such a body, is compiled only.

## Commands

From `<repo root>`. Commit rows a catalog fix moved:

```bash
go test ./pg/catalog/ -run '^TestProgramSimulation(PG18)?$' -update-compat
```

Record all rows on the pinned servers:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./pg/catalog/ -run '^TestProgramSimulation(PG18)?$'
```

Show both schemas for matching rows:

```bash
go test -v ./pg/catalog/ -run '^TestProgramSimulation$' -args -simulate-show='BEGIN ATOMIC'
```
