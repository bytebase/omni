# BigQuery

BigQuery statements are GoogleSQL, which omni's `googlesql` tree parses,
classifies and wraps a result limit in; Bytebase's alias table sends BigQuery
to `googlesql` for the language, while BigQuery's own packages live here:

- `bigquery/parser` is a thin wrapper delegating `Parse` to `googlesql/parser`
  so BigQuery has its own live-BigQuery-judged syntax compatibility suite in
  [`parser/testdata/compat.tsv`](parser/testdata/compat.tsv) (see
  [`parser/testdata/README.md`](parser/testdata/README.md)), separate from the
  Spanner emulator verdicts in `googlesql/parser/testdata/compat.tsv`.
- `bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery`, and
  all six driver tests replay offline in CI from committed tapes under
  `bigquery/driver/testdata/tapes/` through `bigquery/internal/testcontainer`
  and the driver's unexported `client` seam — requiring no Google Cloud
  credentials or network access. `driver/testdata/sync/omni_record.json` holds
  the synced metadata golden, the pin rows are in
  [`testdata/pins.tsv`](testdata/pins.tsv), and
  [`driver/VERIFY.md`](driver/VERIFY.md) holds the corpus, recording commands,
  and verification record.
