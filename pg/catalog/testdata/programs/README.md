# PostgreSQL stored-program catalog simulation

`simulate.tsv` holds the statements in `programFeatures`
(`../../simulate_program_test.go`) on `../fixtures/programs@objects.sql`:
SQL and PL/pgSQL functions and procedures with real bodies, their `ALTER` and
`DROP`, what a body depends on, triggers, event triggers, `CALL` and `DO`.
`TestProgramSimulation` replays it from `../tapes/TestProgramSimulation.tape`
(postgres:17.12), and `TestProgramSimulationPG18` replays the same rows in
`../pg18/programs/simulate.tsv` from postgres:18.7. Both pins answer every row
alike. The readout adds the database's event triggers to `../simulate.tsv`'s.
`pg/review`'s `TestWalkSimulation` checks WALK_THROUGH on these rows
(`programs` subtest).

**266 rows: 257 agree, 9 differ (96.6%), on both pins.**

The catalog checks what PostgreSQL checks at `CREATE` with
`check_function_bodies` on: a SQL body is parsed and analyzed, and its last
statement must return the declared result; a PL/pgSQL body is compiled (syntax,
assigned variables, declared types, condition names, `RETURN`), with the tables
it names left to run time. A SQL-standard body (`RETURN`, `BEGIN ATOMIC`)
records what it uses, so dropping or retyping it fails with 2BP01 or 0A000.

## Gaps

- Run time: `CALL p_broken()` (42P01) and `DO` blocks that raise (P0001), insert
  into a missing table (42P01) or create a table. The catalog resolves the
  procedure and compiles the block, and runs neither.
- An enum literal is not checked against the enum's labels (`m = 'nope'`,
  22P02, twice).
- `LANGUAGE internal AS 'nope_symbol'` (42883): the catalog has no list of the
  server's internal function symbols.
- A SQL-standard body prints a literal without the type PostgreSQL resolves it
  to (`'happy'::mood`, `'a'::text`) and a result column without its `AS` name,
  so two accepted rows store another body text.

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
