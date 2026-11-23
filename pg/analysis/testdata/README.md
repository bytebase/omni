# pg/analysis test data

## span.tsv: query spans judged by PostgreSQL

`span.tsv` holds `ResolveQuerySpan` to PostgreSQL 18.7 (the `pg18` row of
`pg/testdata/pins.tsv`, the release `pg_proc_generated.go` is read from). For
each query of `span/queries.sql`, run under a search_path on a fixture of
`span/fixtures`, it has three rows, one per aspect of the span, each with
omni's answer and PostgreSQL's:

| aspect | omni | PostgreSQL |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the prepared query's RowDescription |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the table OID and attribute number the RowDescription gives each column, resolved through `pg_class` and `pg_attribute`; `~` where it gives none |
| `tables` | the relations of `SourceColumns` | the relations the rewrite rule of `CREATE TEMPORARY VIEW … AS <query>` depends on (`pg_depend`) |

The session column is the search_path. A query PostgreSQL refuses has
`error:<SQLSTATE>` in the expected cell of each of its rows.

omni resolves each query as Bytebase does: in the fixture's database, under
the session's search_path, against the metadata a sync stores (read with an
empty search_path, as `pg/driver` reads it, so a view's definition is
qualified), with names matched exactly.

How PostgreSQL's answers are read, all in `span_test.go` and
`pg/internal/testcontainer/span.go`:

- The RowDescription names a column's origin only for a plain column
  reference, followed through subqueries, CTEs and joins. For an expression,
  an aggregate, a set operation, a recursive CTE or a FULL JOIN's merged
  column it names none, and the expected lineage is `~`: any sources omni
  gives agree (`enginetest/span.Unstated`).
- A whole-row reference names its relation and attribute 0, read as every
  column of the relation. A system column names its negative attribute
  number, read as its name (`ctid`).
- PostgreSQL's span stops at a view: the RowDescription names the view's
  column, and the view's rule depends on the view, not on what the view
  reads. omni's span stops there too.
- The tables are a temporary view's dependencies, not `EXPLAIN`'s scans:
  the plan replaces a view with what it reads and drops relations the
  planner removes (a constant-false condition, an unused LEFT JOIN), while
  the rule lists every relation the query names, as the span does.

The metadata holds tables and views but no functions, as a sync that
predates a function would. A star over `tags()`, which the fixture defines,
is a function omni cannot expand, so those queries reach `ResolveQuerySpan`'s
parse-tree fallback, which takes the function to return one column.

The conversation with PostgreSQL is on `tapes/TestQuerySpanFidelity.tape`,
with every OID resolved to a name before it is written, so the test starts
no server.

### Running it

```bash
go test ./pg/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./pg/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./pg/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. The tables come from a temporary view over the query, so a query no
view can hold (a data-modifying CTE, an output of a pseudo-type such as an
anonymous `record`) fails the recording; pick another form of it.

### Known differences

- `SELECT j.id, j.name FROM customers JOIN orders USING (id) AS j`:
  PostgreSQL refuses it (`j` names only the merged `id`), and omni's
  analyzer does too, but `ResolveQuerySpan` fails open into its parse-tree
  fallback and returns a best-effort span. Bytebase runs the query and
  PostgreSQL refuses it, so no data flows either way.

## query-span/*.yaml

Corpora recorded from Bytebase's old adapters, replayed by
`TestQuerySpanCorpus`.
