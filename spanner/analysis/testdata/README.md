# spanner/analysis test data

## span.tsv: query spans judged by the Spanner emulator

`span.tsv` holds `ResolveQuerySpan`, with Bytebase's Spanner options, to the
Cloud Spanner emulator 1.5.54 (the `default` row of
`spanner/testdata/pins.tsv`). For each query of `span/queries.sql`, run on
the fixture of `span/fixtures`, it has three rows, one per aspect of the span,
each with omni's answer and the emulator's:

| aspect | omni | the emulator |
| --- | --- | --- |
| `results` | the names of `Results` | the field names of the query's `ResultSetMetadata.row_type`; an anonymous column is `""` |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | `~` for every column: the emulator states no column's source |
| `tables` | the tables and views of `SourceColumns` | each table or view without which the query fails to resolve (below) |

A query the emulator refuses has `error:<gRPC code>` in the expected cell of
each of its rows; omni's refusal is `error:omni`.

omni resolves each query as Bytebase does for Spanner (`spannerResolveOptions`:
the default schema is `""`, a named schema scopes its tables), against the
metadata a sync stores for the fixture, read from the emulator's
`INFORMATION_SCHEMA` the way `spanner/driver`'s sync reads it: tables, views
with their definitions, and the columns of both in order. Spanner has no
session schema, so every query's session is `default`.

### Where the emulator's answers come from

The emulator states less than a server does. `AnalyzeQuery` (`QueryMode
PLAN`) answers `query plan unavailable`, and it accepts fine-grained access
control DDL but never checks a privilege, so neither names the tables a query
reads. What it does do is resolve names, so the harness
(`spanner/spanner-oracle/span.go`) asks that:

- For a fixture it builds one database with the whole fixture, and, for each
  table and view, one more where that object alone is gone: a table is
  renamed away (`ALTER TABLE … RENAME TO`), which keeps its interleaved
  children, indexes and change streams; a view is dropped. A view that reads
  the object, which the emulator names when it refuses, is dropped first.
- For a query it reads the output names from the whole fixture, then runs the
  query in each of those databases. The query reads an object exactly when it
  fails there (`InvalidArgument`: table not found). So a view's tables are
  read too, through nested views, while a CTE or an alias that shares a
  table's name, an interleaved parent, and an index named in a hint are not.

The fixtures hold no rows, so no query reads one.

The conversation with the harness is on `tapes/TestQuerySpanFidelity.tape`,
so the test starts no server.

### Running it

```bash
go test ./spanner/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./spanner/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./spanner/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Not covered

- Change streams are read through `READ_<stream>(…)`, a streaming read the
  harness cannot end in one request, so the fixture's change stream is only
  there to show that renaming a table it watches leaves it in place.

- `INFORMATION_SCHEMA`: the emulator exits when a query reads it in a
  database where a table was renamed, so the corpus holds no query of a
  system table.
- Pipe syntax (`FROM t |> …`): the emulator's parser refuses it.

### Known differences

- The emulator refuses analytic functions (`… OVER (…)`), `WITH RECURSIVE`,
  and a top-level `SELECT AS STRUCT` or `SELECT AS VALUE` (a struct cannot be
  a column); Cloud Spanner runs the first two. omni returns a span for each;
  the rows keep the emulator's refusal.
- Lineage is unstated throughout, so the suite holds only its arity. A column
  read through a view names the view's column, not the columns the view
  reads.

## query-span/*.yaml

Corpora recorded from Bytebase's old adapters, replayed whole by
`TestQuerySpanCorpusBigQuery` and `TestQuerySpanCorpusSpanner`.
