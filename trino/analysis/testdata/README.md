# trino/analysis test data

## span.tsv: query spans judged by Trino

`span.tsv` holds `ResolveQuerySpan` to Trino 482 (the `default` row of
`trino/testdata/pins.tsv`). For each query of `span/queries.sql`, run in a
session on a fixture of `span/fixtures`, it has three rows, one per aspect of
the span, each with omni's answer and Trino's:

| aspect | omni | Trino |
| --- | --- | --- |
| `results` | the names of `Results` | the output column names the statement protocol returns |
| `lineage` | each result's `SourceColumns`, results separated by `\|` | each output column's `sourceColumns` in the `output` of `/v1/query/{id}` for `CREATE TABLE … AS SELECT * FROM (<query>) AS q(c1, …)`, Trino's own column lineage |
| `tables` | the tables of `SourceColumns` | `referencedTables` in `/v1/query/{id}` for the query, from Trino's analyzer |

A query Trino refuses has `error:<errorName>` in the expected cell of each of its rows.

omni resolves each query as Bytebase does: the session's catalog is
`DefaultDatabase` and its schema `DefaultSchema`, against the metadata a sync
stores for the fixture (read the way `trino/driver` syncs it), with
`IgnoreCaseSensitive` false, as Bytebase passes it for Trino.

Two readings of Trino's answers, both in `span_test.go`:

- Trino's lineage stops at a view's column. The fixture's answer records the
  lineage of each view's definition, run in the view's schema, and the test
  follows a view's column through it, so a column read through a view names
  the view's column and every column on the way to the tables, as omni does.
- Trino's lineage spells a view's column as its definition wrote it
  (`VipName`); its metadata stores every name in lower case, and that is the
  name the test compares.

The conversation with Trino is on `tapes/TestQuerySpanFidelity.tape` (see
`trino/internal/testcontainer/span.go`), so the test starts no server.

### Running it

```bash
go test ./trino/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./trino/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./trino/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or a fixture
```

A row that agreed and no longer does fails as a regression; `-update-compat`
never commits one. `-span-show <text>` prints both whole answers for the
queries that hold the text. `enginetest/cmd/compatcheck` holds a change to
the file to the rules of `enginetest/span.Guard`.

### Adding a query

Add a line to `span/queries.sql` under the `-- session:` it runs in and record.
The lineage comes from a `CREATE TABLE AS` over the query, so a query whose
output Trino cannot store (an unnamed `NULL`, a type the memory connector
lacks) fails the recording; pick another form of it.

### Known differences

- `JOIN … USING (id)`: Trino's lineage gives the coalesced `id` no source
  column. omni names the column of both sides, which is what a masking policy
  on either needs.
- `SELECT * FROM hr."Vip_Names"`: Trino names the column `VipName`, as the
  view's definition wrote it, while its metadata, and so Bytebase's sync,
  stores `vipname`.

## query-span/*.yaml

Corpora recorded from Bytebase's old adapters. `TestQuerySpanMetadataCorpus`
checks only each statement's type.
