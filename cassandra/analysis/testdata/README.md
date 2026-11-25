# cassandra/analysis test data

## span.tsv: query spans judged by Cassandra

Bytebase masks a CQL query's results by `ResolveQuerySpan`'s span. `span.tsv`
holds that span to Cassandra 5.0.9 (the `driver` row of
`cassandra/testdata/pins.tsv`). For each query of `span/queries.sql`, run in a
session keyspace on the fixture `span/fixtures/base.sql`, it has three rows,
one per aspect of the span, each with omni's answer and Cassandra's:

| aspect | omni | Cassandra |
| --- | --- | --- |
| `results` | the names of `Results` | the column names of the result metadata |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | the keyspace, table and name the result metadata gives each column, where the name is a column of that relation; `~` otherwise |
| `tables` | the relations of `SourceColumns` | the keyspaces and tables the result metadata names |

The session column is the session keyspace. 97 queries: plain columns,
aliases, `*` (with static columns, counters and tuples), `DISTINCT`, `JSON`,
`WRITETIME`, `MAXWRITETIME` and `TTL`, aggregates and `GROUP BY`, native
functions, casts, arithmetic, type hints, the dynamic data masking functions,
user-defined type fields and collection elements and slices, materialized
views, quoted and mixed-case names, and tables of other keyspaces.

omni resolves each query as Bytebase does: in the session keyspace, against
the metadata `cassandra/driver`'s sync stores for a keyspace (tables from
`system_schema.tables`, materialized views from `system_schema.views`, each
with its columns from `system_schema.columns` in `SELECT *` order).

How Cassandra's answers are read, in `span_test.go` and
`cassandra/internal/testcontainer/span.go`:

- gocql prepares the query, and Cassandra's result metadata gives every output
  column a keyspace, a table and a name. The fixture's tables are empty, so no
  row follows it.
- The keyspace and table are those of the relation the query reads, for every
  column, an expression's too. They are the `tables` row.
- The name is the column's own for a plain column reference, the alias where
  there is one, and the selector's text for anything else
  (`writetime(name)`, `system.tojson(home)`, `count`, `[json]`). A name that is
  a column of the relation is that column's lineage; any other is `~`
  (`enginetest/span.Unstated`): the metadata does not say which columns an
  expression or an aliased column reads, so any sources omni gives agree.
- A materialized view is named as itself. The metadata says nothing of the
  base table its rows come from, so a materialized view's column and the view
  among the tables are followed by `~`.

The conversation with Cassandra is on
`tapes/TestQuerySpanFidelity.tape`, so the test starts no server. The tape
server starts with `materialized_views_enabled: true`, which Cassandra 5
ships off.

### Agreement

| aspect | rows | agree | differ | of which `~` |
| --- | --- | --- | --- | --- |
| results | 97 | 97 | 0 | 0 |
| lineage | 97 | 96 | 1 | 61 |
| tables | 97 | 97 | 0 | 7 |
| all | 291 | 290 | 1 | 68 |

The first recording agreed on 231 of 291 rows. The fixes:

- CQL folds an unquoted name to lower case. omni kept `SELECT NAME FROM Users`
  as written, so it named a column and a table no policy is set on.
- `SELECT JSON` returns one column, `[json]`, which holds every selected
  column. omni answered one result per selected column, so the results did not
  line up with the one column a masking policy is applied to. It now answers
  one `[json]` result that reads every selected column.
- An expression's result now reads the columns of its arguments. A function
  call or a cast had no sources, or, with an alias, the bogus column `""`, so
  `SELECT blobAsText(payload)` escaped a policy on `payload`. A field of a
  user-defined type, `home.city`, named the column `city`; it now reads
  `home`. A collection element or slice reads its collection.
- A column without an alias is named as Cassandra names it: the selector's
  text, with a native function in the `system` keyspace and `COUNT(*)` as
  `count`.
- `cassandra/driver`'s sync ordered a static column among the regular ones by
  name, but `SELECT *` returns the static columns before every regular one, so
  a star over `events` put `region` in the wrong place and a policy by
  position on the wrong column. The sync, and so omni's star, now order them
  as Cassandra does.
- The sync stores each materialized view, with its columns and the base-table
  column each copies (`DependencyColumns`). omni expands a star over a view,
  which it refused before, and every column of a view also reads the
  base-table column, so a policy on `users.email` applies to
  `users_by_email.email`. Cassandra states only the view (`~`).

### Known difference

- `SELECT name AS email FROM users`: the metadata names the alias, `email`,
  which is also a column of `users`, so it reads as the column `email`. omni
  traces the column to `name`, where its values come from.
