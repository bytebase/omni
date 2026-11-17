# BigQuery

This tree is driver-only: `bigquery/driver` is the conversation with a BigQuery
project, and the statements are GoogleSQL, which omni's `googlesql` tree parses,
classifies and wraps a result limit in. There is no `bigquery` parser, catalog or
analysis package, and there should not be one: Bytebase's alias table sends
BigQuery to `googlesql` for the language, and only the service client is BigQuery's
own, as `spanner/driver` is Spanner's beside the same grammar.

`bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery`, and every
live test passed against a scratch dataset it created in the Google Cloud project
`bytebase-dev` on 2026-11-17: the sync of a dataset, the reads, a command's affected
rows, `Execute`'s script and the not-found mapping. Six defects that run found are
fixed, and `driver/testdata/sync/omni_record.json` is what the sync built.

What is still missing is the half a project cannot give: the seam that would let a
tape hold those answers, so no tape exists and CI runs none of it — the live tests
skip without `BIGQUERY_CREDENTIALS`. Its package doc says so,
`docs/engine-test-matrix.md` says so, the pin row is in
[`testdata/pins.tsv`](testdata/pins.tsv), and
[`driver/VERIFY.md`](driver/VERIFY.md) holds the recording plan, what that run
answered, and what step 4 would still add.
