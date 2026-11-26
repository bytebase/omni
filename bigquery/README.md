# BigQuery

BigQuery's own tree, forked from `spanner` (then `googlesql`) the way `redshift` forked from
`pg`: a mechanical copy of `ast`, `parser`, `analysis`, `diagnostics` and
`v1review`, then BigQuery's dialect deltas. `spanner` holds Spanner's grammar.

- `bigquery/parser` is judged by live BigQuery in
  [`parser/testdata/compat.tsv`](parser/testdata/compat.tsv) (see
  [`parser/testdata/README.md`](parser/testdata/README.md)). Its rows come from
  `spanner/testdata/GoogleSQL.g4`, the GoogleSQL grammar both engines share.
- `bigquery/analysis`'s query spans are judged by BigQuery's dry run in
  [`analysis/testdata/span.tsv`](analysis/testdata/span.tsv) (see
  [`analysis/testdata/README.md`](analysis/testdata/README.md)), recorded on a
  dataset the recording builds and deletes, and replayed from a tape.
- `bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery`, and
  all six driver tests replay offline in CI from committed tapes under
  `bigquery/driver/testdata/tapes/` through `bigquery/internal/testcontainer`
  and the driver's unexported `client` seam — requiring no Google Cloud
  credentials or network access. `driver/testdata/sync/omni_record.json` holds
  the synced metadata golden, the pin rows are in
  [`testdata/pins.tsv`](testdata/pins.tsv), and
  [`driver/VERIFY.md`](driver/VERIFY.md) holds the corpus, recording commands,
  and verification record.
