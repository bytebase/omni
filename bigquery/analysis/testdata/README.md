# bigquery/analysis test data

## span.tsv: query spans judged by BigQuery

`span.tsv` holds `ResolveQuerySpan`, with Bytebase's BigQuery options, to
BigQuery itself in the Google Cloud project and location of the `span` row of
`bigquery/testdata/pins.tsv`. For each query of `span/queries.sql`, run on the
fixture of `span/fixtures`, it has three rows, one per aspect of the span, each
with omni's answer and BigQuery's:

| aspect | omni | BigQuery |
| --- | --- | --- |
| `results` | the names of `Results` | the field names of the dry run's `schema` |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | `~` for every column: BigQuery states no column's source |
| `tables` | the tables and views of `SourceColumns`, as `dataset.table` | the dry run's `referencedTables` and `referencedLogicalViews`, as `dataset.table` |

A query BigQuery refuses has `error:<reason>` (`invalidQuery`, `notFound`) in
the expected cell of each of its rows; omni's refusal is `error:omni`.

omni resolves each query as Bytebase does for BigQuery
(`bigqueryResolveOptions`, with the session's dataset as the default one),
against the metadata `bigquery/driver`'s sync stores for the fixture: one
unnamed schema holding every table of the dataset with its top-level columns in
order, and every view with its columns and its definition as
`INFORMATION_SCHEMA.VIEWS` states it. The sync stores no routines.

### Where BigQuery's answers come from

A recording (`bigquery/internal/testcontainer/span.go`) creates the dataset
`omni_span` in the pinned project, labelled `omni=span`, runs the fixture's DDL
in it, and reads back each table and view the dataset lists, with each view's
definition from `INFORMATION_SCHEMA.VIEWS`. Each query is then
a dry run with `omni_span` as the default dataset: BigQuery resolves it and
states its output schema and the tables and views it references, reads no
data, and bills nothing. The fixture's tables hold no rows. When the test ends
the recording deletes the dataset with everything in it; it refuses to delete
a dataset of that name without the label.

The Go client drops `referencedLogicalViews`, so the dry run is read from the
REST answer. A table read at a point in time (`FOR SYSTEM_TIME AS OF`) is
named with the snapshot it read, `table@<ms>`; the snapshot is dropped, since
it moves every run.

BigQuery reports no column lineage a dry run can read: the job's statistics
name tables, not columns, and Dataplex lineage records only jobs that ran.
Every column's lineage is therefore `~`, and the suite holds its arity alone.

The conversation is on `tapes/TestQuerySpanFidelity.tape`, so a replay needs
no credentials and reaches no server.

### Running it

```bash
go test ./bigquery/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./bigquery/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 BIGQUERY_CREDENTIALS=default BIGQUERY_PROJECT=bytebase-dev go test ./bigquery/analysis/ -run '^TestQuerySpanFidelity$'   # after changing the corpus or the fixture
```

`BIGQUERY_CREDENTIALS` is a key file, key JSON, or `default` for the default
credential chain. A recording takes about a minute and a half. A row that
agreed and no longer does fails as a regression; `-update-compat` never
commits one. `-span-show <text>` prints both whole answers for the queries
that hold the text.

### Known differences

- Table functions: BigQuery reports the tables a table function reads. The
  sync stores no routines, so omni lists none of them, and fails closed on a
  star over a table function and on its range variable used as a value
  (`SELECT x FROM f() AS x`), whose columns it cannot list.
- Pipe syntax: omni fails closed on `CALL`, `MATCH_RECOGNIZE` and a
  recursive union, and on any operator that changes the projection of an
  input it cannot list column by column (a wildcard table, a table without
  metadata). The other operators agree.
- Stars omni does not expand: a struct column's fields (`Ship.*`, and
  `UNNEST` of an array of structs under `*`), and a wildcard table's columns
  (`` `events_*` ``, whose schema BigQuery takes from the newest matching
  table).
- `FULL` and `LEFT UNION ALL BY NAME` do not parse.
- omni resolves names BigQuery refuses: a table name in another case (table
  names are case-sensitive), a table or column that does not exist, a column
  qualified by its dataset, an ambiguous column, a range variable inside a
  subquery that has an alias, `= ANY (subquery)`, a CTE column list, an
  unqualified table function, and a top-level `SELECT AS VALUE` that is not a
  struct.
- Both refuse, so the answers differ only in the error: `FOR SYSTEM TIME`
  without the underscore, and an unqualified table function under `*` or
  `CALL`.

## query-span/*.yaml

The corpus recorded from Bytebase's old BigQuery adapter, replayed whole by
`TestQuerySpanCorpusBigQuery`. Where BigQuery answers otherwise it holds
BigQuery's answer (result names, and no CTE listed as a table).
