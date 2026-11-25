# cosmosdb/analysis test data

## span.tsv: query spans judged by the Cosmos DB emulator

`span.tsv` holds `Analyze` to the Cosmos DB emulator 2.14.28 (the `emulator`
row of `cosmosdb/testdata/pins.tsv`), which is what Bytebase's query span for
Cosmos DB, and so its document masking, is built on. For each query of
`span/queries.sql` it has three rows, one per aspect of the span:

| aspect | omni | the emulator |
| --- | --- | --- |
| `results` | each projection's `Name` | the members of the items the query returns, in the order the gateway sends them; one unnamed item (`""`) when the gateway's query plan states `hasSelectValue` |
| `lineage` | each projection's `SourcePaths`, projections separated by `\|` | the stored properties each returned value holds (below), and `~` |
| `tables` | nothing (`-`) | `~` |

A query the emulator refuses has `error:<code>` in its expected cells, the
query compiler's code (`SC1001` for a syntax error). The corpus holds none:
a refused query returns nothing to mask.

Paths are written as Bytebase's masker reads them: the container's own name
(a path's first selector) is dropped, so `c.address.city` is
`$.address.city` whatever the FROM clause calls the container, and array
indexes are dropped, since a masking policy names a property and covers every
element of an array.

### Where the emulator's lineage comes from

The emulator states no lineage: its query plan names ORDER BY and GROUP BY
expressions, but nothing about the SELECT list. The values it returns do. The
fixture, `span/fixtures/document.sql`, is one JSON document (the harness reads
a fixture from a `.sql` file) whose every leaf is a value no other leaf holds,
such as `"«$.address.city»"`. A returned value names a stored property when it
is a stored subtree (`c.address` returns the whole of `address`), or, for a
string, when it holds a stored string in any case (`CONCAT` holds both of its
arguments, `UPPER` keeps the letters). Otherwise the value is computed
(`c.n + 1`, `IS_DEFINED`, `COUNT`) and names nothing.

One run on one document shows a value's sources, never that it has no other:
`c.name ?? c.nick` returns `c.name` here and `c.nick` on another document. So
every lineage column of the emulator's answer also holds `~`, and agrees with
any omni answer that names every source the emulator showed. A missing source
is what leaks a masked value, and that is what a row catches; a source omni
names beyond the emulator's, like `c.nick` above, only masks more.

The emulator adds system properties to a stored document (`_rid`, `_etag`,
`_ts`, ...), whose values name the store that answered. The tape holds
`"(system)"` in their place, and a value shared by several properties names
none of them.

Cosmos DB reads the one container a request names: the FROM clause only
names an alias for it (`FROM c`, `FROM root r` and `FROM mycontainer x` all
read the same container). So neither side names the tables read.

omni's `SELECT *` is expanded over the stored document at its `StarPaths`,
as Bytebase's masker looks each returned property up in the stored schema.

The conversation with the emulator is on `tapes/TestQuerySpanFidelity.tape`
(see `cosmosdb/internal/testcontainer/span.go`), so the test starts nothing.

### Running it

```bash
go test ./cosmosdb/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./cosmosdb/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./cosmosdb/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or the fixture
```

A recording starts the emulator, which binds port 8081, so one runs per
machine; it uses several GiB, so check `free -g` first. `-span-show <text>`
prints both whole answers for the queries that hold the text.

### Adding a query

Add a line to `span/queries.sql` and record. A query must return a row,
since an empty result states nothing, and every returned property must be
defined: Cosmos DB leaves an undefined property out of the item.

### Known differences

None: every row agrees. omni names more sources than the emulator showed
wherever a value may come from either of several properties (`??`, `IIF`,
`? :`, and the condition those read), and for values computed from a
property (`-c.n`, `ARRAY_LENGTH(c.tags)`, `IS_DEFINED(c.nick)`), which is the
safe answer for masking.
