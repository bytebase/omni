# Trino catalog test data

`simulate.tsv` holds one row per schema, table, view, materialized view,
comment and `TRUNCATE` statement that Trino 482 (the `default` row of
`../../testdata/pins.tsv`) accepts in `trino/parser/testdata/compat.tsv`, run
once on each fixture in `fixtures/` (320 statements × 4 fixtures = 1,280 rows).
Each fixture holds table `t` in schema `memory.s`, the session's schema, and
names `a` differently: nothing (`one_table@missing`), a table
(`one_table@table`), a view (`objects@view`), or a schema holding a table
(`objects@schema`).

A row holds what `Exec` does and what Trino did on its memory connector:
Trino's `errorName` when it refuses the statement, or `ok:` and a hash of what
the memory catalog holds afterwards. That is each schema, each table as `SHOW
CREATE TABLE` prints it, and each view as the head of `SHOW CREATE VIEW` (name,
comment, security) with its columns and their types. A view's query is left
out: Trino stores it rewritten by its SqlFormatter, which omni has no copy of.
`TestSimulation` (`../simulate_test.go`) replays Trino's side from
`tapes/TestSimulation.tape` and starts no server.

**1,242 / 1,280 rows agree.** The 38 that do not are `ALTER TABLE` column
changes on a view (`objects@view`): the memory connector fails on them with
`GENERIC_INTERNAL_ERROR`, a NullPointerException in `MemoryMetadata`, where
omni refuses with `TABLE_NOT_FOUND`, as a connector that does not hand a view
to `ALTER TABLE` does.

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to Trino's answer:

```bash
go test ./trino/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned Trino again, after adding a statement or fixture:

```bash
OMNI_RECORD=1 go test ./trino/catalog/ -run '^TestSimulation$'
```
