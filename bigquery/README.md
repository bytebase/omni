# BigQuery

This tree is driver-only: `bigquery/driver` is the conversation with a BigQuery
project, and the statements are GoogleSQL, which omni's `googlesql` tree parses,
classifies and wraps a result limit in. There is no `bigquery` parser, catalog or
analysis package, and there should not be one: Bytebase's alias table sends
BigQuery to `googlesql` for the language, and only the service client is BigQuery's
own, as `spanner/driver` is Spanner's beside the same grammar.

`bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery`. One run
against the Google Cloud project `bytebase-dev` on 2026-11-17 asked everything that
needs no dataset and agreed with the port apart from three defects it fixed, but
the project holds no dataset, so the sync of a dataset and the writes have never
run against the service. Nothing is recorded either — no tapes, no sync golden and
no pin row — so CI runs none of it. Its package doc says so,
`docs/engine-test-matrix.md` says so, and
[`bigquery/driver/VERIFY.md`](driver/VERIFY.md) holds both the recording plan — the
dataset and the grant that would open it, the corpus, the tests to record in order,
and the parity harness that holds it to the driver it replaces — and what that run
answered.
