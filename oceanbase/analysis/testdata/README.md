# oceanbase/analysis test data

## span.tsv: query spans judged by OceanBase

`span.tsv` holds `ResolveQuerySpan` to OceanBase CE 4.3.5.6 in MySQL mode (the
`default` row of `oceanbase/testdata/pins.tsv`). For each query of
`span/queries.sql`, run with a database selected on a fixture of
`span/fixtures`, it has three rows, one per aspect of the span, each with
omni's answer and OceanBase's:

| aspect | omni | OceanBase |
| --- | --- | --- |
| `results` | the names of `Results` | each column definition's name in the result set of the prepared statement (`COM_STMT_PREPARE`) |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the column a definition's `org_table` and `org_name` name, read through a view's stored definition |
| `tables` | the tables of `SourceColumns` | every table and view OceanBase checks the SELECT privilege on when the query runs |

The corpus is mysql/analysis's. OceanBase runs with `lower_case_table_names`
1, so the fixture has `MixedCase` and not also `mixedcase`, and stores names
in lower case.

A query OceanBase refuses has `error:<number>` in the expected cell of each of
its rows. omni's refusal is the number OceanBase gives for the same cause:
1064 for a parse error, 1054 for a missing column, 1146 for a missing table.

omni resolves each query as Bytebase does: the session's database is
`DefaultDatabase`, against the metadata a sync stores for the fixture (each
table's and view's columns, and a view's `VIEW_DEFINITION`), with
`IgnoreCaseSensitive` true, as on the pinned server.

### What OceanBase states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what OceanBase does not
say. A row agrees when omni names every entry OceanBase does state, and
possibly more.

- **lineage.** A column definition's `org_table` and `org_name` name the
  relation the query reads and its column: a table, or a view, whose stored
  definition (fully qualified, as OceanBase resolved it) is read through to
  its tables, through views over views. A column of a derived table or CTE
  names that relation, which has no database, and an expression names
  nothing, so both are `~`. OceanBase's plan names the tables it reads by
  their aliases, so it is no source.
- **tables.** An account that holds no SELECT privilege runs the query, and
  is granted each table or view its refusal (1142) names, until the query
  runs: the tables and views it names, not the CTEs it defines. OceanBase
  checks a view as the view, never the tables behind it, so a query that
  reads a view states its tables plus `~`.

The conversation with OceanBase is on `tapes/TestQuerySpanFidelity.tape` (see
`oceanbase/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./oceanbase/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./oceanbase/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./oceanbase/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text. `enginetest/cmd/compatcheck` holds a change to
the file to the rules of `enginetest/span.Guard`.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The reader of a view's definition (`span_resolved_test.go`) knows the
form OceanBase stores.

### Known differences

- `JOIN … USING (id)` and `NATURAL JOIN`: OceanBase resolves the merged `id`
  to the first table's column. omni names the column of both sides, which is
  what a masking policy on either needs, since in a matched row both hold
  the value.
- `SELECT id FROM orders JOIN customers ON …`, where both tables have `id`:
  OceanBase refuses the ambiguous name (1052), and omni resolves it to one
  table. OceanBase returns no rows, so no masked value can pass.

## query-span/*.yaml

Bytebase's OceanBase query-span fixture, recorded from the resolver masking
ran on before this package existed (`query_span_corpus_test.go`). Five of its
CTE cases listed the CTE's name among the tables read; span.tsv shows
OceanBase reads only the tables behind it, and those entries were removed.
`case_insensitive.yaml` reads table `T` of a catalog that stores `t`, and now
lists it as `t`, the name OceanBase checks.
