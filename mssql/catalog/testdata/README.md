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

**3,576 / 3,576 rows (100.00%) agree** (298 statements × 12 fixtures).

## Collations

`../collations.txt` is every collation the pinned server lists in
`sys.fn_helpcollations()` (5,508). `TestCollations`
(`../collations_test.go`) holds the file to the server and to what SQL Server
2022 accepts as a column's collation, for every listed name in either case and
for 39,047 other forms: version-less and versioned, with and without `_SC` and
`_UTF8`, with flags swapped. SQL Server 2022 accepts no form the list does not
hold, and `database_default` and `catalog_default` only for a column or an
expression. The catalog refuses a `COLLATE` (Msg 448) only when no version can
know the name (`UnknownCollation`, which mssql/review's WALK_THROUGH shares):
it takes every listed name, the 32 `_160` collations SQL Server 2025 adds,
which 2022 refuses, and `Compatibility_` names such as
`Compatibility_52_409_30003`, which every version takes and none lists.
`TestUnknownCollation` pins those. `TestCollations` replays
`tapes/TestCollations.tape`; rewrite the file from the server with:

```bash
OMNI_RECORD=1 go test ./mssql/catalog/ -run '^TestCollations$'
```

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

- 384 rows on `one_table@column`: every SELECT, WITH, or parenthesized query
  the parser's `compat.tsv` judges accepted, by its grammar rule.
- 364 rows on `queries@tables`, labelled `feature:<feature>`: T-SQL query
  features the grammar's rows reach only with the one name `a` — select
  lists, joins, predicates, CASE, CAST, CONVERT and PARSE, built-in and
  user-defined functions, grouping, windows, TOP and OFFSET/FETCH, derived
  tables, CTEs, set operations, APPLY, PIVOT and UNPIVOT, FOR XML and FOR
  JSON, table sources, view options, and `COLLATE` (in a query, and in the
  table, type, partition function and database DDL that names a collation).
  Some are refused on purpose.

**745 / 750 rows (99.33%) agree** (382 / 384 grammar rows, 363 / 366
feature rows; 265 rows SQL Server accepts). The 5 that do not:

- `SELECT a.$PARTITION.a(1)` (Msg 208): Bytebase's metadata has no
  partition functions, so a catalog loaded from it could not tell a missing
  one from a synced one, and `$PARTITION` is not checked.
- `WITH CHANGE_TRACKING_CONTEXT(@a) SELECT *` (Msg 137): SQL Server's parser
  declares variables before it checks where `CHANGE_TRACKING_CONTEXT` stands,
  and a view declares none; omni's parser refuses the clause in a view's
  query (Msg 22108) without knowing which variables a batch declares.
- `SELECT b FROM nope.dbo.a` (Msg 208): a catalog holds one database, so
  another database's objects are not checked.
- `SELECT name FROM sys.nope` (Msg 208): the system views are the server's,
  not the schema's, and differ by version and edition (Azure SQL Database
  has its own, such as `sys.database_firewall_rules`), so any name in `sys`
  or `INFORMATION_SCHEMA` binds.
- `SELECT b FROM OPENROWSET('x', 'y', 'z')` (Msg 7222): which OLE DB
  providers a server allows is its configuration, not its catalog.

A view ends its batch, and Bytebase sends a script batch by batch, so `Exec`
refuses a statement after a view before `GO` (Msg 156, or 102 when it does
not start with a keyword), as `SELECT * SELECT *` measures.

`LoadSQL` does neither: a desired-state script need not separate its views
with `GO`, and may define a view before the objects it reads. `Exec` on a
loaded catalog binds a view's query.

Commit the rows a fix moved, or run every row on the pinned server again:

```bash
go test ./mssql/catalog/ -run '^TestQuerySimulation$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/catalog/ -run '^TestQuerySimulation$'
```

## Program suite

`programs/simulate.tsv` holds one row per stored program with a real body,
run on `programs@objects`, so the catalog is held to how SQL Server binds a
body when it creates the program. `simulate.tsv` already runs the grammar's
program statements, whose bodies name nothing. `TestProgramSimulation`
(`../simulate_program_test.go`) replays it from
`tapes/TestProgramSimulation.tape`. Rows are labelled `program:<kind>`:

| Kind | Rows | SQL Server accepts |
| --- | ---: | ---: |
| `procedure` | 73 | 42 |
| `scalar-function` | 26 | 15 |
| `inline-function` | 18 | 9 |
| `multi-statement-function` | 11 | 4 |
| `trigger` (DML) | 27 | 15 |
| `ddl-trigger` | 7 | 4 |

**162 / 162 rows (100.00%) agree**; the first recording agreed on 92.

What a row holds the catalog to: variables are declared in the order the
body is read (Msg 137, 1087); a statement whose tables all exist is bound
(Msg 207, 209, 8120, 110, 120, 8144, ...), and one that names a missing or
temporary table is deferred, as is a call of a user-defined function that
does not exist; a function writes only its table variables and returns no
rows (Msg 443, 444) and ends with RETURN (Msg 455, 178); an inline function
is bound as a view is (Msg 4514, 4506, 1033); a DML trigger reads
`inserted` and `deleted` with its table's columns and lives in its table's
schema (Msg 2103), and ALTER TRIGGER keeps its table (Msg 2110). A
collation the body gives a column is checked as the body is parsed, even in a
deferred statement, and a COLLATE operator as its statement is bound (Msg
448).

SQL Server stores a database DDL trigger, but the catalog models none, so a
DDL trigger row is held to its refusal only: the stored schema read back has
no DDL triggers.

Each row first drops the schema-bound functions and views the row before it
created (`programResetSQL`), which would keep `_sim_reset` from dropping
their tables.

```bash
go test ./mssql/catalog/ -run '^TestProgramSimulation$' -update-compat
```

```bash
OMNI_RECORD=1 go test -v -timeout 20m ./mssql/catalog/ -run '^TestProgramSimulation$'
```
