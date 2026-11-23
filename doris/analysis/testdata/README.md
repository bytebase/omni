# doris/analysis test data

## span.tsv: query spans judged by Doris

`span.tsv` holds `ResolveQuerySpan` to Doris 3.1.4 (the `default` row of
`doris/testdata/pins.tsv`). For each query of `span/queries.sql`, run with a
database selected on a fixture of `span/fixtures`, it has three rows, one per
aspect of the span, each with omni's answer and Doris's:

| aspect | omni | Doris |
| --- | --- | --- |
| `results` | the names of `Results` | each column definition's name in the query's result set |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | each column definition's `org_table` and `org_name` |
| `tables` | the tables and views of `SourceColumns` | every table and view the analyzed plan (`EXPLAIN ANALYZED PLAN`) reads: a scan's `qualified=` name, and the three-part qualifier of the alias over a view or a system table |

Doris answers every refusal with error 1105 and a message, so a refused
query's rows name the cause the message states: `error:syntax`,
`error:unknown-column`, `error:unknown-table`, `error:unknown-database` or
`error:ambiguous-column`. omni's refusal is a parse error or its
`NotFoundError`, named the same way.

omni resolves each query with `DefaultDatabase` set to the session's database
and `Catalog` holding the metadata a sync stores for the fixture (each
table's and view's columns, and a view's `VIEW_DEFINITION`), case-sensitively,
as Doris matches table names.

### What Doris states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what Doris does not say.
A row agrees when omni names every entry Doris does state, and possibly more.

- **lineage.** Doris fills `org_table` and `org_name` for a select item that
  is one column reference, through derived tables and CTEs, to the table it
  reads. Through a view it names the view's column, so the columns of the
  view's definition are `~`. An expression, a set operation's column and a
  table function's column name nothing, so they are `~` alone.
- **tables.** The analyzed plan binds each view to its definition and keeps
  the view, and the optimizer has dropped nothing yet (an impossible `WHERE`
  still scans its table), so it states every table and view and needs no
  `~`.

The conversation with Doris is on `tapes/TestQuerySpanFidelity.tape` (see
`doris/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./doris/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./doris/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./doris/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Known differences

- An unqualified column two joined tables share, `USING` columns included:
  Doris refuses it as ambiguous, and omni names the column of both tables.
  Doris returns no rows, so no masked value can pass, and both sides are what
  a masking policy on either needs.
- `SELECT * EXCEPT (amount, status) FROM orders`: Doris accepts it
  (`selectColumnClause` in `DorisParser.g4`), and omni's parser refuses it.
- `QUALIFY` and `WITH RECURSIVE`: 3.1.4's grammar has neither, and Doris
  refuses both; omni's parser accepts them. omni refuses the recursive CTE
  anyway, as a table it cannot find.
