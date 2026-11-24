# mssql/analysis test data

## span.tsv: query spans judged by SQL Server

`span.tsv` holds `ResolveQuerySpan` to SQL Server 2022 (the `default` row of
`mssql/testdata/pins.tsv`). For each query of `span/queries.sql`, run in a
session on a fixture of `span/fixtures`, it has three rows, one per aspect of
the span, each with omni's answer and SQL Server's:

| aspect | omni | SQL Server |
| --- | --- | --- |
| `results` | the names of `Results` | `name` of each column `sys.dm_exec_describe_first_result_set(@sql, NULL, 1)` returns, but the hidden ones browse mode adds |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | `source_database`, `source_schema`, `source_table` and `source_column` of the same rows: browse information on, which names the table column a value is, read through any view |
| `tables` | the tables of `SourceColumns` | the `Object` elements of the query's `SET SHOWPLAN_XML` plan, and the tables and views `sys.dm_sql_referenced_entities` names for a procedure made over the query, with the tables and views each view names, in turn |

The plan reads a view's tables and never names the view, which is why the
procedure's references are read too.

A query SQL Server refuses has `error:<number>` in the expected cell of each of
its rows. omni's refusal is named by the number SQL Server gives the same
refusal (207, 208, 209, 4104), and `error:omni` otherwise.

Where SQL Server names no source, `~` (`span.Unstated`) stands in the lineage:
for an expression, and where browse mode names a view's column because it
cannot see through the view, as for a view over a `UNION`. omni agrees with
such a column whatever sources it names.

omni resolves each query as Bytebase does: the session's database is
`DefaultDatabase` and its schema `DefaultSchema`, against the metadata a sync
stores for the fixture (read the way `mssql/driver` syncs it), with
`IgnoreCaseSensitive` true, as Bytebase passes it for SQL Server. The
sessions are `dbo`, the default schema of `sa`, which the recording runs as.

The conversation with SQL Server is on `tapes/TestQuerySpanFidelity.tape`, so
the test starts no server.

### Running it

```bash
go test ./mssql/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./mssql/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./mssql/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and record.
Fixture tables have no keys, so the optimizer drops no join from the plan.
A query has to fit in a procedure (`CREATE PROCEDURE … AS <query>`), which
is how its references are read.

### Known differences

- `INTERSECT` and `EXCEPT`: SQL Server names the left arm's column; omni
  names both arms'. The right arm decides which rows come back, and a
  value `INTERSECT` returns is in both, so a masking policy on either side
  applies.
- PIVOT: omni refuses to trace it (`TypeNotSupportedError`), and Bytebase
  refuses the query rather than return its columns unmasked.
- A view: omni names the view and not the tables beneath it in its own
  database. SQL Server checks a reader's permission on the view alone when
  its owner owns the tables, and a Bytebase grant on the view names it;
  naming the tables too would refuse that reader. What a view reads in
  another database omni names, since ownership chains stop at a database by
  default. Lineage, which masking reads, names the tables' columns either
  way.

## query-span/*.yaml

Bytebase's query-span corpora; `TestQuerySpanCorpus` reproduces them.
