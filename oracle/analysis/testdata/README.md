# oracle/analysis test data

## span.tsv: query spans judged by Oracle

`span.tsv` holds `ResolveQuerySpan` to Oracle Database 23ai Free 23.26.3 (the
`default` row of `oracle/testdata/pins.tsv`). For each query of
`span/queries.sql`, run with the session's current schema set, on the fixture
of `span/fixtures`, it has three rows, one per aspect of the span, each with
omni's answer and Oracle's:

| aspect | omni | Oracle |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the query's cursor, as go-ora's describe information gives them |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the domain Oracle gives each column of a view over the query: the test gives every table column a domain of its own, and Oracle carries a column's domain to an output that reads it |
| `tables` | the tables of `SourceColumns` | the objects of the query's `EXPLAIN PLAN`, and the objects a view over the query depends on (`ALL_DEPENDENCIES`), read through synonyms and, in turn, through the views they name |

A query Oracle refuses has `error:<ORA number>` in the expected cell of each of
its rows. omni's refusal is named 904 for a missing column, 942 for a missing
table or schema, and `omni` otherwise.

omni resolves each query as Bytebase does: the session's schema is
`DefaultDatabase`, against the metadata a sync stores for the fixture (read
the way `oracle/driver`'s sync reads it: user-generated columns in `COLUMN_ID`
order with invisible ones last, a view's `TEXT`, a materialized view's
`QUERY`, the private synonyms, and the public synonyms that name a fixture
schema), with names matched exactly. The database link resolves to the schema
it names, or to its user's.

### What Oracle states, and what it leaves unstated

`~` (`span.Unstated`) in an expected cell stands for what Oracle does not say.
A row agrees when omni names every entry Oracle does state, and possibly more.

- **lineage.** Oracle 23ai keeps a column's domain on a view column that is
  the column, read through subqueries, WITH elements, joins, views and scalar
  subqueries, and through `NVL`, `COALESCE`, `CASE` and `DECODE`, which keep
  one argument's domain even when another column can be returned
  (`NVL(email, name)` carries `email`'s). That names one source; it does not
  say there is no other (that other argument, a `CASE` condition, the other
  side of a `NATURAL` join), so every lineage cell is `<column> ~`, or `~`
  alone where Oracle carries no domain: an expression, an aggregate, a set operation over two
  columns, a pseudo-column (`ROWNUM`, `ROWID`, `LEVEL`), a materialized view,
  and a table through a database link.
- **tables.** The plan names the tables a merged view reads and drops none
  (the fixture tables have no keys); the dependencies name what the plan
  merges away, such as a synonym's target and the view itself. A plan
  operation through a database link, and a dependency on a synonym through
  one, are another database's objects: `~`. The plan's own objects, which no
  dictionary view lists (the temporary table of a materialized WITH element,
  `VW_FOJ_0` of a FULL JOIN, `VW_DCL_…` of a decorrelated lateral view), are
  left out.

The conversation with Oracle is on `tapes/TestQuerySpanFidelity.tape`, so the
test starts no server. The test runs as SYSTEM; the fixture grants each schema
`SELECT ANY TABLE`, which the view over a query needs.

### Running it

```bash
go test ./oracle/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./oracle/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./oracle/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or the fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and
record. A query has to fit in a view (`CREATE VIEW … (C1, …) AS <query>`),
which is how its lineage and dependencies are read. The fixture's schemas
avoid the names Bytebase counts as Oracle's system schemas (HR, SCOTT, …):
omni gives a query on one no lineage.

### Known differences

- DUAL: omni names no table for it. It holds no data, and naming `SYS.DUAL`
  would make a query that reads a user table and DUAL a mix of user and
  system tables, which Bytebase refuses. Oracle's plan of an `UNPIVOT` reads
  DUAL too.
- A query Oracle refuses for an unknown or ambiguous column (ORA-00904,
  ORA-00918), and `SELECT *` over a `NATURAL` join of two tables of the same
  name (ORA-00918): omni returns a span. Oracle runs no such query, so no data
  flows either way; refusing every name omni cannot resolve would refuse
  queries that call a function without parentheses.

## query-span/*.yaml

Bytebase's query-span corpora; `TestQuerySpanCorpus` reproduces them. Two
view cases named only the view among the tables read; span.tsv shows Oracle
reads the tables beneath it too, and those were added.
