# BigQuery

This tree is driver-only: `bigquery/driver` is the conversation with a BigQuery
project, and the statements are GoogleSQL, which omni's `googlesql` tree parses,
classifies and wraps a result limit in. There is no `bigquery` parser, catalog or
analysis package, and there should not be one: Bytebase's alias table sends
BigQuery to `googlesql` for the language, and only the service client is BigQuery's
own, as `spanner/driver` is Spanner's beside the same grammar.

`bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery`, and all
six driver tests replay offline in CI from committed tapes under
`bigquery/driver/testdata/tapes/` through `bigquery/internal/testcontainer` and
the driver's unexported `client` seam — requiring no Google Cloud credentials or
network access. `driver/testdata/sync/omni_record.json` holds the synced metadata
golden, the pin row is in [`testdata/pins.tsv`](testdata/pins.tsv), and
[`driver/VERIFY.md`](driver/VERIFY.md) holds the corpus, recording commands, and
verification record.
