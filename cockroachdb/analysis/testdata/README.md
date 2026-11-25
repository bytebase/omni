# cockroachdb/analysis test data

## span.tsv: query spans judged by CockroachDB

Bytebase resolves a CockroachDB query's span with `pg/analysis`'s
`ResolveQuerySpan`, as it does PostgreSQL's. `span.tsv` holds that span to
CockroachDB v25.2.3 (the `default` row of `cockroachdb/testdata/pins.tsv`).
For each query of `span/queries.sql`, run under a search_path on a fixture of
`span/fixtures`, it has three rows, one per aspect of the span, each with
omni's answer and CockroachDB's:

| aspect | omni | CockroachDB |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the prepared query's RowDescription |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the table OID and attribute number the RowDescription gives each column, resolved through `pg_class` and `pg_attribute`; `~` where it gives none |
| `tables` | the relations of `SourceColumns` | the relations a view over the query depends on (`crdb_internal.backward_dependencies`) |

The session column is the search_path. 157 queries, most of them
`pg/analysis`'s span corpus where CockroachDB accepts it (it has no
partitioned tables, system columns, `TABLESAMPLE`, `JSON_TABLE`,
`USING (...) AS` or grouping sets), and a section on hidden columns.

omni resolves each query as Bytebase does: in the fixture's database, under
the session's search_path, against the metadata `cockroachdb/driver`'s sync
stores (`information_schema.columns`, with a hidden column invisible, and
`pg_views.definition`, which CockroachDB stores with every name qualified by
database and schema), with names matched exactly, and with `ThroughViews`
set.

How CockroachDB's answers are read, all in `span_test.go` and
`cockroachdb/internal/testcontainer/span.go`:

- The RowDescription names a column's origin only for a plain column
  reference, followed through subqueries, joins, `INTERSECT` and `EXCEPT`.
  For an expression, an aggregate, a `UNION`, any CTE or a `FULL JOIN`'s
  merged column it names none, and the expected lineage is `~`: any sources
  omni gives agree (`enginetest/span.Unstated`).
- Unlike PostgreSQL, CockroachDB's RowDescription follows a view to the
  table column behind it and never names the view's column. omni, with
  `ThroughViews`, names both, so a masking policy on either applies. A
  column whose origin is also the origin of a column of a view the query
  reads therefore expects that origin and `~`.
- The tables are the view's dependencies, not `EXPLAIN`'s scans: the plan
  replaces a view with what it reads and drops relations the optimizer
  removes, while the dependencies list every relation the query names, in
  any clause, as the span does. A view is a dependency; what it reads is
  not.

The conversation with CockroachDB is on `tapes/TestQuerySpanFidelity.tape`,
with every descriptor ID resolved to a name before it is written, so the
test starts no server.

### Agreement

| aspect | rows | agree | differ | of which `~` |
| --- | --- | --- | --- | --- |
| results | 157 | 156 | 1 | 0 |
| lineage | 157 | 155 | 2 | 77 |
| tables | 157 | 156 | 1 | 0 |
| all | 471 | 467 | 4 | 77 |

The first recording agreed on 417 of 471 rows. Three fixes closed the rest:

- A table with no primary key has a hidden `rowid`, which `*` leaves out.
  The sync stored it as an ordinary column, so omni's `*` added it and every
  later column of the result moved one place, out of line with the column
  a masking policy applies to. The sync now marks it `IsInvisible`, and
  `pg/catalog` installs such a column hidden: out of `*` and a join's
  columns, reached by its name.
- CockroachDB stores `CREATE VIEW v (a, b) AS SELECT x, y ...` with the
  definition alone. `pg/catalog` now names a loaded view's columns by the
  synced columns, so `SELECT a FROM v` resolves.
- `ResolveOptions.ThroughViews` follows a view's column into its
  definition.

### Known differences

- `SELECT email FROM customers INTERSECT SELECT status FROM orders`:
  CockroachDB names the left column, omni both. A row of the result is in
  both inputs, so a policy on either applies.
- `SELECT * FROM items JOIN accounts USING (rowid)`: CockroachDB refuses a
  hidden column in `USING`, and omni's analyzer does too, but
  `ResolveQuerySpan` fails open into its parse-tree fallback and returns a
  best-effort span. CockroachDB refuses the query, so no data flows.

### Running it

```bash
go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./cockroachdb/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The tables come from a view over the query, so a query no view can
hold (an output of an anonymous record type, a call of a user-defined
set-returning function) fails the recording; pick another form of it.
