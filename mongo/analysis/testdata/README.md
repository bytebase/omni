# mongo/analysis test data

## span.tsv: collections read, judged by MongoDB

`span.tsv` holds `Analyze`'s `Reads`, with views followed by `ExpandViews`,
to MongoDB 8.0.38 (the `default` row of `mongo/testdata/pins.tsv`). For each
statement of `span/queries.js`, run in the database its `// session:` line
names on the fixture `span/fixtures/base.sql` builds, it has three rows:

| aspect | omni | MongoDB |
| --- | --- | --- |
| `results` | `~` | `~` |
| `lineage` | `~` | `~` |
| `tables` | `Reads`, each view expanded against the views the server stored, as `database.collection` | the collections the statement's `explain("queryPlanner")` plan reads |

A statement MongoDB refuses has `error:<codeName>` in each expected cell.

How MongoDB's answers are read, in `mongo/internal/testcontainer/span.go`:

- The fixture is mongosh statements (in a `.sql` file, the name
  `enginetest/span` reads) run through gomongo, as Bytebase runs
  them, each in the database the last `// database:` line names; a view is
  made with the `create` command. The answer is each database's views with
  their listCollections options, which is the definition `mongo/driver`'s
  sync stores and omni expands.
- A plan names each collection it scans (`queryPlanner.namespace`, one per
  pipeline a `$unionWith` runs) and each a stage joins: `$lookup` and
  `$graphLookup`'s `from`, `$unionWith`'s `coll`, and a lowered lookup's
  `foreignCollection`. Only the plan's structure is read (its plan stages
  and pipeline stages), never a filter, a `$match` or the `command` the plan
  echoes, where a document is the query's data. The target of `$out` or
  `$merge` is written, not read, and neither side names it.
- A view the statement targets is already resolved in the plan. A view a
  stage names is not, so it is planned in turn (`aggregate([])` over it) and
  the answer names what it reads. MongoDB's answer is therefore collections
  only.
- gomongo has no `db.getSiblingDB` and does not plan `findOne`, so the judge
  runs `db.getSiblingDB("x").c.find()` as `db.c.find()` in database `x`, and
  plans a `findOne` as the `find` MongoDB runs for it.

Neither side states output fields. A queryPlanner plan does not name the
fields a statement returns, and `Analyze` does not either: Bytebase masks
each returned document against its collection's schema, with each
`$lookup`'s `as` field taking the joined collection's (`Joins`). So
`results` and `lineage` hold `~` (`enginetest/span.Unstated`) on both sides,
and only `tables` is judged.

### Kept differences

Every remaining difference has omni naming more than the plan, which is
the safer answer: a masking policy or access check on any collection named
still applies.

- A view is named with what it reads. The plan names the collections under
  a view and never the view, but a policy may be on the view.
- A `$lookup` on a collection that does not exist: the plan of an empty
  collection is `EOF` and drops the stages after it, while omni names the
  joined collection too.
- `$lookup` with `from: {db, coll}`: MongoDB refuses it outside its own
  internal collections (`FailedToParse`), and omni names the other
  database's collection.

### Running it

```bash
go test ./mongo/analysis/ -run '^TestQuerySpanFidelity$'                  # replay
go test ./mongo/analysis/ -run '^TestQuerySpanFidelity$' -update-compat   # commit a closed gap
OMNI_RECORD=1 go test ./mongo/analysis/ -run '^TestQuerySpanFidelity$'    # after changing the corpus or the fixture
```
