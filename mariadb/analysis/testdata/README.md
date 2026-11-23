# mariadb/analysis test data

## span.tsv: query spans judged by MariaDB

`span.tsv` holds `ResolveQuerySpan` to MariaDB 11.8.8 (the `default` row of
`mariadb/testdata/pins.tsv`). For each query of `span/queries.sql`, run with a
database selected on a fixture of `span/fixtures`, it has three rows, one per
aspect of the span, each with omni's answer and MariaDB's:

| aspect | omni | MariaDB |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the query's result set |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the columns each select item names in the query MariaDB resolved: the fully qualified text of the note (1003) `EXPLAIN EXTENDED` leaves, after MariaDB merged views, derived tables and CTEs |
| `tables` | the tables of `SourceColumns` | every table a FROM clause of that text names |

The corpus is mysql/analysis's, with `VALUES ROW(...)` spelled `VALUES (...)`.
`TABLE` and `LATERAL`, which MariaDB lacks, stay as refused queries.

A query MariaDB refuses has `error:<number>` in the expected cell of each of
its rows. omni's refusal is the number MariaDB gives for the same cause: 1064
for a parse error, 1054 for a missing column, 1146 for a missing table.

omni resolves each query as Bytebase does: the session's database is
`DefaultDatabase`, against the metadata a sync stores for the fixture (each
table's and view's columns, and a view's `VIEW_DEFINITION`), with
`IgnoreCaseSensitive` false, as on a server with `lower_case_table_names` 0.

### What MariaDB states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what MariaDB does not
say. A row agrees when omni names every entry MariaDB does state, and possibly
more.

- **lineage.** A select item that is one column reference states its column
  exactly, through materialized derived tables, CTEs and views (a view reads
  through its stored definition). An item that names no column states it has
  none (`-`). Any other expression states the columns it names, plus `~`:
  MariaDB says nothing about which of them a masking policy should follow,
  nothing about the columns of a subquery inside it, and the key of an
  `<expr_cache>` is not read as a source. The text is printed after
  optimization, which replaces the columns of a table read as a constant
  (`const` or `system` in the plan) by their values, so in such a query an
  item that names no column is `~` too.
- **tables.** The text names a merged view's tables and never the view, and
  leaves out a table the optimizer eliminated from a join. A query whose
  optimizer trace merges or materializes a view, eliminates a table, or whose
  plan reads no table (an impossible `WHERE`, an aggregate answered from an
  index) states its tables plus `~`.

The conversation with MariaDB is on `tapes/TestQuerySpanFidelity.tape` (see
`mariadb/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./mariadb/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./mariadb/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./mariadb/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text. `enginetest/cmd/compatcheck` holds a change to
the file to the rules of `enginetest/span.Guard`.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The reader of the resolved query (`span_resolved_test.go`) knows the
form MariaDB prints; a recording fails on a query whose text it cannot read.

### Known differences

- `JOIN … USING (id)` and `NATURAL JOIN`: MariaDB resolves the merged `id` to
  the first table's column. omni names the column of both sides, which is
  what a masking policy on either needs, since in a matched row both hold
  the value.
- `SELECT id FROM orders JOIN customers ON …`, where both tables have `id`:
  MariaDB refuses the ambiguous name (1052), and omni resolves it to one
  table. MariaDB returns no rows, so no masked value can pass.
- In a `WITH RECURSIVE` clause MariaDB lets a CTE read a later one
  (`WITH RECURSIVE c AS (SELECT id FROM orders), orders AS (SELECT 1 AS id)`
  reads no table). omni's tables follow it; its lineage still resolves such a
  name to the table, which names more sources than MariaDB reads.

## query-span/*.yaml

Bytebase's MariaDB query-span fixture, recorded from the resolver masking ran
on before this package existed (`query_span_corpus_test.go`). Five of its CTE
cases listed the CTE's name among the tables read; span.tsv shows MariaDB
reads only the tables behind it, and those entries were removed.
