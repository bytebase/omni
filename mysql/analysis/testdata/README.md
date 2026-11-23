# mysql/analysis test data

## span.tsv: query spans judged by MySQL

`span.tsv` holds `ResolveQuerySpan` to MySQL 8.4.13 (the `mysql84` row of
`mysql/testdata/pins.tsv`). For each query of `span/queries.sql`, run with a
database selected on a fixture of `span/fixtures`, it has three rows, one per
aspect of the span, each with omni's answer and MySQL's:

| aspect | omni | MySQL |
| --- | --- | --- |
| `results` | the names of `Results` | each column definition's name in the result set of the prepared statement (`COM_STMT_PREPARE`) |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the columns each select item names in the query MySQL resolved: the fully qualified text of `EXPLAIN FORMAT=JSON` (its `query`), after MySQL merged views and derived tables |
| `tables` | the tables of `SourceColumns` | every table the plan of `EXPLAIN FORMAT=JSON` accesses (`schema_name`, `table_name`) |

A query MySQL refuses has `error:<number>` in the expected cell of each of its
rows. omni's refusal is the number MySQL gives for the same cause: 1064 for a
parse error, 1054 for a missing column, 1146 for a missing table.

omni resolves each query as Bytebase does: the session's database is
`DefaultDatabase`, against the metadata a sync stores for the fixture (each
table's and view's columns, and a view's `VIEW_DEFINITION`), with
`IgnoreCaseSensitive` false, as on a server with `lower_case_table_names` 0.

### What MySQL states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what MySQL does not say.
A row agrees when omni names every entry MySQL does state, and possibly more.

- **lineage.** A select item that is one column reference states its column
  exactly, through materialized derived tables, CTEs and views (a view reads
  through its stored definition). An item that names no column states it has
  none (`-`). Any other expression states the columns it names, plus `~`:
  MySQL says nothing about which of them a masking policy should follow, and
  nothing about the columns of a subquery inside it.
- **tables.** The plan names the tables behind a view and never the view, and
  leaves out a table it read while optimizing (a constant table, an
  aggregate answered from an index, an impossible `WHERE`). A query that
  reads a view (an optimizer trace `view` step) or whose plan hides a table
  states its tables plus `~`.

### Why not `org_table` and `org_name`

The column definitions of a result set carry `org_table` and `org_name`, but
through a derived table MySQL fills them with the derived table's column name
on its underlying table: `SELECT * FROM (SELECT amount AS id FROM orders) d`
reports `orders.id` for a column that reads `orders.amount`. The resolved
query has no such error, so the lineage comes from it. The tape keeps the
column definitions whole.

### Why 8.4 and not 8.0

8.4's `explain_json_format_version = 2` names the database and table of each
access; 8.0's JSON names the alias. Names and lineage come from the same
resolver in both.

The conversation with MySQL is on `tapes/TestQuerySpanFidelity.tape` (see
`mysql/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./mysql/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./mysql/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./mysql/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text. `enginetest/cmd/compatcheck` holds a change to
the file to the rules of `enginetest/span.Guard`.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The reader of the resolved query (`span_resolved_test.go`) knows the
form MySQL prints; a recording fails on a query whose text it cannot read.

### Known differences

- `JOIN … USING (id)` and `NATURAL JOIN`: MySQL resolves the merged `id` to
  the first table's column. omni names the column of both sides, which is
  what a masking policy on either needs, since in a matched row both hold
  the value.
- `SELECT id FROM orders JOIN customers ON …`, where both tables have `id`:
  MySQL refuses the ambiguous name (1052), and omni resolves it to one table.
  MySQL returns no rows, so no masked value can pass.

## query-span/*.yaml

Bytebase's MySQL query-span fixture, recorded from the resolver masking ran on
before this package existed (`query_span_corpus_test.go`). Five of its CTE
cases listed the CTE's name among the tables read; span.tsv shows MySQL reads
only the tables behind it, and those entries were removed.
