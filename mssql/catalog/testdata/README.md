# SQL Server catalog test data

`simulate.tsv` holds one row per core DDL statement that the pinned SQL Server
2022 (`../../testdata/pins.tsv`,
`mcr.microsoft.com/mssql/server:2022-CU27-GDR1-ubuntu-22.04`, `16.0.4305.5`)
accepts in `mssql/parser/testdata/compat.tsv`, run once on each fixture in
`fixtures/`:

- Four namings of one plain table: `one_table@missing`, `one_table@table`,
  `one_table@column`, `one_table@index`.
- Eight object-kind fixtures: `objects@view`, `objects@fk`, `objects@schema`,
  `objects@sequence`, `objects@synonym`, `objects@procedure`,
  `objects@function`, `objects@trigger`.

Simulated statement families cover database-local catalog DDL on schemas,
tables, indexes (relational, columnstore, XML, selective XML, spatial), views,
sequences, synonyms, stored procedures, functions, and triggers (excluding
server-, database-, security-, partition-, fulltext-, and external-object
statements).

A row holds omni's outcome and SQL Server's — the primary SQL Server error
number (`mssql.Error.Number`) when refused, or `ok:` and a hash of the
normalized DDL objects stored in the database afterwards. `TestSimulation`
(`../simulate_test.go`) replays it from `tapes/TestSimulation.tape` and starts
no server. The design is
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md` §3.

**3,072 / 3,072 rows (100.00%) agree** (256 statements × 12 fixtures).

## Commands

From `<repo root>`. Commit the rows a catalog fix moved to the engine's answer:

```bash
go test ./mssql/catalog/ -run '^TestSimulation$' -update-compat
```

Run every row on the pinned server again, which a new row or fixture needs:

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/catalog/ -run '^TestSimulation$'
```

Verify that live SQL Server matches the committed tape and `simulate.tsv` with
zero differences:

```bash
OMNI_RECORD=verify go test -v -timeout 20m ./mssql/catalog/ -run '^TestSimulation$'
```

Show both schemas of the rows that hold a statement's text:

```bash
go test -v ./mssql/catalog/ -run '^TestSimulation$' -args -simulate-show='ALTER SEQUENCE'
```

## Query suite

`queries/simulate.tsv` holds one row per query, each wrapped as
`CREATE VIEW _v AS <query>`, so the catalog is held to SQL Server's binding of
the query: the objects and columns it names, the view's column names, and the
rules a view's query keeps (Msg 1033, 4511, 4506, 8120, ...). An accepted row
hashes the stored schema, view definition included. `TestQuerySimulation`
(`../simulate_query_test.go`) replays it from
`tapes/TestQuerySimulation.tape`. It has two parts:

- 334 rows on `one_table@column`: every SELECT, WITH, or parenthesized query
  the parser's `compat.tsv` judges accepted, by its grammar rule.
- 329 rows on `queries@tables`, labelled `feature:<feature>`: T-SQL query
  features the grammar's rows reach only with the one name `a` — select
  lists, joins, predicates, CASE, CAST, CONVERT and PARSE, built-in and
  user-defined functions, grouping, windows, TOP and OFFSET/FETCH, derived
  tables, CTEs, set operations, APPLY, PIVOT and UNPIVOT, FOR XML and FOR
  JSON, table sources, and view options. Some are refused on purpose.

**655 / 663 rows (98.79%) agree** (329 / 334 grammar rows, 326 / 329
feature rows; 249 rows SQL Server accepts). The 8 that do not need something
the catalog does not model:

- `SELECT $PARTITION.a(1)` and its three variants (Msg 208): partition
  functions are outside the catalog, and a catalog loaded from synced
  metadata would not know them either, so `$PARTITION` is not checked.
- `SELECT b FROM nope.dbo.a` (Msg 208): another database's objects are not
  checked, since a catalog holds one database.
- `SELECT name FROM sys.nope` (Msg 208): system views are not modeled, so
  any name in `sys` or `INFORMATION_SCHEMA` binds.
- `SELECT b FROM OPENROWSET('x', 'y', 'z')` (Msg 7222): which OLE DB
  providers a server allows is its configuration, not its catalog.
- `SELECT * SELECT *` (Msg 156): SQL Server ends a view's batch at its
  query; omni parses the second `SELECT` as a statement of its own, as
  Bytebase runs a script statement by statement.

`LoadSQL` does not bind a view's query, since a desired-state script may
define a view before the objects it reads; `Exec` on a loaded catalog does.

Commit the rows a fix moved, or run every row on the pinned server again:

```bash
go test ./mssql/catalog/ -run '^TestQuerySimulation$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/catalog/ -run '^TestQuerySimulation$'
```
