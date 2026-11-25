# elasticsearch/analysis test data

## span.tsv: query spans judged by Elasticsearch

`span.tsv` holds the analysis Bytebase reads for masking
(`elasticsearch.GetQuerySpan`, whose `ElasticsearchAnalysis` this package
builds) to Elasticsearch 8.19.25 (the `driver` row of
`elasticsearch/testdata/pins.tsv`). For each request of `span/queries.http`,
run on a fixture of `span/fixtures`, it has three rows:

| aspect | omni | Elasticsearch |
| --- | --- | --- |
| `results` | the leaf fields of each index in `tables` that the request's `_source` filtering (`SourceFields`, `SourceExcludes`, `SourceDisabled`) keeps | the dotted path of every value in the `_source` of the documents the request returns |
| `lineage` | for each field, `<index>.<field>` for each index in `tables` whose mapping defines it | for each field, the `_index` of each returned document that holds it, plus `~` |
| `tables` | `ExpandIndices(Indices, synced indices)` | the targets the request names (its path, a multi-get's docs, a multi-search's headers; every index when it names none) through `GET /_resolve/index/<target>?expand_wildcards=open`: each index, each alias's indices, each data stream |

A request Elasticsearch refuses has `error:<type>` in each of its expected
cells. A data stream is named for itself: its backing index
(`.ds-<stream>-<date>-<generation>`) is the stream in both columns.

omni resolves each request as Bytebase can: against what
`elasticsearch/driver`'s sync stores for the fixture, the indices the cat API
lists less the hidden ones, each with the fields its mapping defines (which
an `ObjectSchema` describes). The test reads those with the sync's own two
requests and each index's `_mapping`.

The fixture is `span/fixtures/shop.sql`, a list of console requests separated
by blank lines (the harness reads every fixture as `<name>.sql`). Every
document holds every field its index maps, so a field missing from a result
is one the filtering removed. One run cannot show that a field has no other
source, so every engine lineage column holds `~`.

The conversation with Elasticsearch is on `tapes/TestQuerySpanFidelity.tape`
(through `elasticsearch/internal/testcontainer`), so the test starts no
server.

### Running it

```bash
go test ./elasticsearch/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./elasticsearch/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./elasticsearch/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or the fixture
```

### Known differences

- **Aliases** (`people`, `everyone`, `adults`, `peo*`): Elasticsearch reads
  the indices behind the alias; the sync stores no aliases, so omni keeps the
  name as written and knows none of its fields. This is the one difference
  that is not safer for masking: a policy on `users` does not reach
  `GET /people/_search`. Closing it needs the sync to store aliases, which
  the metadata proto has no place for, or Bytebase to mask each document by
  the `_index` it names.
- **Data streams and hidden indices** (`events-app`, `hidden_audit`, and the
  stream that `*`, `_all` and a request with no index also read): the sync
  stores neither, so no masking policy can name them and omni knows none of
  their fields.
- **Multi-search** with a `_source` per search: omni keeps one filter for the
  whole request, so it names every field of every index, more than the
  engine returns.
- **A missing index**: Elasticsearch refuses the request
  (`index_not_found_exception`); omni names it.
