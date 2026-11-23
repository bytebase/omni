# tidb/analysis test data

## span.tsv: query spans judged by TiDB

`span.tsv` holds `ResolveQuerySpan` to TiDB v8.5.5 (the `default` row of
`tidb/testdata/pins.tsv`). For each query of `span/queries.sql`, run with a
database selected on a fixture of `span/fixtures`, it has three rows, one per
aspect of the span, each with omni's answer and TiDB's:

| aspect | omni | TiDB |
| --- | --- | --- |
| `results` | the names of `Results` | each column definition's name in the result set of the prepared statement (`COM_STMT_PREPARE`) |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the expression of each output column in the top projection of `EXPLAIN FORMAT = 'brief'`, or, where the plan has none, the column definition's `org_table` and `org_name` |
| `tables` | the tables of `SourceColumns` | the `TABLE_NAMES` `information_schema.STATEMENTS_SUMMARY` records once the query ran |

The corpus is mysql/analysis's. TiDB runs with `lower_case_table_names` 2, so
the fixture has `MixedCase` and not also `mixedcase`, and
`SELECT * FROM mixedcase` reads `MixedCase`. `LATERAL`, derived column lists
and `VALUES ROW(...)` stay as queries TiDB refuses.

A query TiDB refuses has `error:<number>` in the expected cell of each of its
rows. omni's refusal is the number TiDB gives for the same cause: 1064 for a
parse error, 1054 for a missing column, 1146 for a missing table.

omni resolves each query as Bytebase does: the session's database is
`DefaultDatabase`, against the metadata a sync stores for the fixture (each
table's and view's columns, and a view's `VIEW_DEFINITION`).

### What TiDB states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what TiDB does not say.
A row agrees when omni names every entry TiDB does state, and possibly more.

- **lineage.** The plan names each column it reads as `db.table.column`, by
  the table's own name, through aliases, views, derived tables and CTEs, and
  each column it computes as `Column#N`, which an operator defines with
  `expr->Column#N`. An output column that is one column, through such
  definitions, states it exactly. Any other expression states the columns it
  names, plus `~`: TiDB says nothing about which of them a masking policy
  should follow, and a constant says nothing either, since TiDB evaluates an
  uncorrelated subquery while planning and the plan holds its value. A plan
  whose top operator is no projection (a plain read, a join, an aggregate, a
  union) lists no output columns; there a plain column reference states the
  `org_table`.`org_name` of its column definition, in the definition's
  database, when the fixture has that column and the definition's table is no
  view. Through a view TiDB takes the name from the view's select list and the
  database from the view (`BuildDataSourceFromView`,
  `pkg/planner/core/logical_plan_builder.go`), so such a column is `~`, as is
  any other column without an origin.
- **tables.** The statement summary names the tables and views a query
  names, without the CTEs it defines, and never the tables behind a view, so
  a query that reads a view states its tables plus `~`.

The conversation with TiDB is on `tapes/TestQuerySpanFidelity.tape` (see
`tidb/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./tidb/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./tidb/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./tidb/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text. `enginetest/cmd/compatcheck` holds a change to
the file to the rules of `enginetest/span.Guard`.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The plan reader (`span_plan_test.go`) knows the operators TiDB
prints; a recording fails on a plan whose top projection does not match the
result set.

### Known differences

- `JOIN … USING (id)` and `NATURAL JOIN`: TiDB resolves the merged `id` to
  the first table's column. omni names the column of both sides, which is
  what a masking policy on either needs, since in a matched row both hold
  the value.
- `INTERSECT` and `EXCEPT`: TiDB's output column is the first branch's.
  omni names the columns of every branch, as for `UNION`: an `INTERSECT`
  value is in both, and a masking policy on either side applies.
- `SELECT id FROM orders JOIN customers ON …`, where both tables have `id`:
  TiDB refuses the ambiguous name (1052), and omni resolves it to one table.
  TiDB returns no rows, so no masked value can pass.
- `VALUES ROW(1, 2 + 3), ROW(4, 5)`: TiDB's grammar takes it (`RowValue`,
  `parser.y`), and its planner refuses it with 1051. omni parses it, as the
  grammar does, and resolves its columns, which read nothing.
