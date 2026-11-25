# cockroachdb/analysis test data

## span.tsv: query spans judged by CockroachDB

`ResolveQuerySpan` resolves a CockroachDB query's span on CockroachDB's own
AST, against the metadata Bytebase synced. `span.tsv` holds that span to
CockroachDB v25.2.3 (the `default` row of `cockroachdb/testdata/pins.tsv`).
For each query of `span/queries.sql`, run under a search_path on a fixture of
`span/fixtures`, it has three rows, one per aspect of the span, each with
omni's answer and CockroachDB's:

| aspect | omni | CockroachDB |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the prepared query's RowDescription |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the table OID and attribute number the RowDescription gives each column, resolved through `pg_class` and `pg_attribute`; `~` where it gives none |
| `tables` | the relations of `SourceColumns` | the relations a view over the query depends on (`crdb_internal.backward_dependencies`); `~` for a statement no view can hold |

The session column is the search_path. 184 queries: most of them
`pg/analysis`'s span corpus where CockroachDB accepts it (it has no
partitioned tables, system columns, `TABLESAMPLE`, `JSON_TABLE`,
`USING (...) AS` or grouping sets), a section on hidden columns, and a
section on CockroachDB's own syntax: `AS OF SYSTEM TIME`, index hints, join
hints, `ORDER BY PRIMARY KEY`, `IF` and `IFNULL`, statement sources, `SHOW`,
`UPSERT` and `EXPLAIN`.

omni resolves each query as Bytebase does: in the fixture's database, under
the session's search_path, against the metadata `cockroachdb/driver`'s sync
stores (`information_schema.columns`, with a hidden column invisible, and
`pg_views.definition`, which CockroachDB stores with every name qualified by
database and schema), with names matched exactly. A refusal names the
SQLSTATE CockroachDB would give (42703 for a column no relation has).

How CockroachDB's answers are read, all in `span_test.go` and
`cockroachdb/internal/testcontainer/span.go`:

- The RowDescription names a column's origin only for a plain column
  reference, followed through subqueries, joins, `INTERSECT` and `EXCEPT`.
  For an expression, an aggregate, a `UNION`, any CTE, a statement source or
  a `FULL JOIN`'s merged column it names none, and the expected lineage is
  `~`: any sources omni gives agree (`enginetest/span.Unstated`).
- CockroachDB's RowDescription follows a view to the table column behind it
  and never names the view's column. omni names both, so a masking policy on
  either applies. A column whose origin is also the origin of a column of a
  view the query reads therefore expects that origin and `~`.
- The tables are the view's dependencies, not `EXPLAIN`'s scans: the plan
  replaces a view with what it reads and drops relations the optimizer
  removes, while the dependencies list every relation the query names, in
  any clause, as the span does. A view is a dependency; what it reads is
  not. CockroachDB refuses a view over `AS OF SYSTEM TIME`, `UPSERT`, `SHOW`
  or `EXPLAIN` with a syntax error, so for those the tables are `~`.

The conversation with CockroachDB is on `tapes/TestQuerySpanFidelity.tape`,
with every descriptor ID resolved to a name before it is written, so the
test starts no server.

### Agreement

| aspect | rows | agree | differ | of which `~` |
| --- | --- | --- | --- | --- |
| results | 184 | 178 | 6 | 0 |
| lineage | 184 | 181 | 3 | 88 |
| tables | 184 | 184 | 0 | 8 |
| all | 552 | 543 | 9 | 96 |

On the first 157 queries, the span Bytebase resolved through `pg/analysis`
agreed on 467 of 471 rows; the native resolver agrees on 470, all but the
`INTERSECT` row below.

### Known differences

- A `SHOW`, an `EXPLAIN`, or a query of system relations only (such as a
  `[SHOW ...]` statement source) gives no `Results`: as in `pg/analysis`,
  only a SELECT of user tables has them. CockroachDB names their columns, so
  the results row differs, and the lineage row too where there are several
  columns.
- `SELECT email FROM customers INTERSECT SELECT status FROM orders`:
  CockroachDB names the left column, omni both. A row of the result is in
  both inputs, so a policy on either applies.

### Running it

```bash
go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The tables come from a view over the query. A statement CockroachDB
refuses in a view with a syntax error records them as `~`; one it refuses for
another reason (an output of an anonymous record type, a call of a
user-defined set-returning function) fails the recording, so pick another
form of it.
