# starrocks/analysis test data

## span.tsv: query spans judged by StarRocks

`span.tsv` holds `ResolveQuerySpan` to StarRocks 3.4.10 (the `default` row of
`starrocks/testdata/pins.tsv`). For each query of `span/queries.sql`, run with
a database selected on a fixture of `span/fixtures`, it has three rows, one
per aspect of the span, each with omni's answer and StarRocks's:

| aspect | omni | StarRocks |
| --- | --- | --- |
| `results` | the names of `Results` | each column definition's name in the query's result set |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | each output slot of the logical plan (`EXPLAIN LOGICAL`), traced to the scanned columns it reads (`span_plan_test.go`) |
| `tables` | the tables and views of `SourceColumns` | the tables and views the query dump (`POST /api/query_dump`) says the analyzer bound: its `table_meta` and `view_meta` |

StarRocks leaves `org_table` and `org_name` empty in a column definition, so
the plan is where lineage comes from. A scan names its table without its
database; the query dump's tables supply it.

A refused query's rows name the cause StarRocks's message states:
`error:syntax`, `error:unknown-column`, `error:unknown-table`,
`error:unknown-database` or `error:ambiguous-column`. omni's refusal is a
parse error or its `NotFoundError`, named the same way.

omni resolves each query with `DefaultDatabase` set to the session's database
and `Catalog` holding the metadata a sync stores for the fixture (each
table's and view's columns, a view's `VIEW_DEFINITION`, and a materialized
view's `MATERIALIZED_VIEW_DEFINITION`, its whole `CREATE` statement),
case-sensitively, as StarRocks matches table names.

The fixture holds a few rows: StarRocks plans a scan of an empty table as
`EMPTY`, which names no column.

### What StarRocks states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what StarRocks does not
say. A row agrees when omni names every entry StarRocks does state, and
possibly more.

- **lineage.** An output slot that is a scanned column, or a rename of one,
  states it exactly. An expression, an aggregate or a set operation's
  column states the scanned columns its definition reads, plus `~`. A slot
  the optimizer folded away (an impossible `WHERE`, a materialized view not
  yet refreshed) is `~` alone. The plan inlines views, so it states no
  view's column: every column of a query that reads a view gets `~`.
- **tables.** The query dump lists every table and view the analyzer bound,
  before the optimizer dropped any, so it needs no `~`. For a materialized
  view it lists the view's base tables too.

The conversation with StarRocks is on `tapes/TestQuerySpanFidelity.tape` (see
`starrocks/internal/testcontainer/span.go`), so the test starts no server.
The recording turns off materialized-view rewrite, which would swap the
tables a plan scans.

### Running it

```bash
go test ./starrocks/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./starrocks/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./starrocks/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Known differences

- An unqualified column two joined tables share, `USING` columns included:
  StarRocks refuses it as ambiguous, and omni names the column of both
  tables. StarRocks returns no rows, so no masked value can pass.
- `SELECT * EXCEPT (amount, status) FROM orders`: StarRocks 3.4.10 refuses
  it, and omni's parser accepts it.
- `rank() OVER (PARTITION BY status ORDER BY amount)`: StarRocks names the
  column with the sort direction it adds (`ORDER BY amount ASC `); omni keeps
  the text as written. Names do not decide masking, which follows position.
