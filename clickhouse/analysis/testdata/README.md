# clickhouse/analysis test data

## span.tsv: query spans judged by ClickHouse

`span.tsv` holds `ResolveQuerySpan` to ClickHouse 26.8.38.9 (the `default`
row of `clickhouse/testdata/pins.tsv`). For each query of `span/queries.sql`,
run in a session's database on a fixture of `span/fixtures`, it has three
rows, one per aspect of the span, each with omni's answer and ClickHouse's:

| aspect | omni | ClickHouse |
| --- | --- | --- |
| `results` | the names of `Results` | the result set's column names |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the table columns under each projection of `EXPLAIN QUERY TREE passes = 1`, the analyzer's resolved tree before any optimization: a column of a subquery, CTE or set operation is followed to its projection, a column of a view through `EXPLAIN QUERY TREE` of the view's stored `as_select` |
| `tables` | the tables and views of `SourceColumns` | `system.query_log`'s `tables` for the query, which holds every view and the tables behind it |

A refused query's rows hold the name of ClickHouse's error code
(`errorCodeToName`): `error:UNKNOWN_IDENTIFIER`, `error:UNKNOWN_TABLE`,
`error:UNKNOWN_DATABASE`, `error:SYNTAX_ERROR`, `error:AMBIGUOUS_IDENTIFIER`.
omni's refusal is a parse error, its `NotFoundError`, or
`ErrAmbiguousColumn`, named the same way.

omni resolves each query with `DefaultDatabase` set to the session's database
and `Catalog` holding what a sync (`clickhouse/driver`) stores for the
fixture: each table's and view's columns, and a view's `create_table_query`
as its definition.

### What ClickHouse leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what ClickHouse does not
say. A row agrees when omni names every entry ClickHouse does state, and
possibly more.

- **lineage.** A column of a table function (`numbers(10)`) is `~`: the tree
  names the function, not what it reads. A table function's own rows read no
  table, which is what omni answers.
- **lineage.** A recursive CTE's working table is a temporary table named
  afresh for each query, so a column read from it is `~`; the CTE's arms
  name the tables it reads.
- **tables.** query_log's `system.one`, which a query without `FROM` reads,
  and a recursive CTE's temporary table are dropped: neither is a table the
  query names. A table function is listed as `_table_function.<name>`,
  which does not say what it reads (`merge` reads tables), so it is `~`.

The conversation with ClickHouse is on `tapes/TestQuerySpanFidelity.tape`
(see `clickhouse/internal/testcontainer/span.go`), so the test starts no
server.

### Running it

```bash
go test ./clickhouse/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./clickhouse/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./clickhouse/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Known differences

- A `USING` column of an inner or left join: ClickHouse computes it from the
  left table's column, and omni names both tables' columns. An equi-join's
  two columns hold the same values, so a masking policy on either needs the
  result column.
- `SELECT *` over `merge()` of tables whose columns differ: ClickHouse orders
  the merged columns its own way, and omni lists the first table's columns,
  then each later table's new ones. Each column's lineage is the same.
