# BigQuery

This tree is driver-only: `bigquery/driver` is the conversation with a BigQuery
project, and the statements are GoogleSQL, which omni's `googlesql` tree parses,
classifies and wraps a result limit in. There is no `bigquery` parser, catalog or
analysis package, and there should not be one: Bytebase's alias table sends
BigQuery to `googlesql` for the language, and only the service client is BigQuery's
own, as `spanner/driver` is Spanner's beside the same grammar.

`bigquery/driver` is a port of Bytebase's `backend/plugin/db/bigquery` and has never
asked a real project: BigQuery has no image omni can start and no emulator that
answers its API, so there are no tapes, no sync goldens and no pin row. Its package
doc says so, `docs/engine-test-matrix.md` says so, and
[`bigquery/driver/VERIFY.md`](driver/VERIFY.md) holds the recording plan — the
project and account to pin, the corpus, the tests to record in order, and the parity
harness that holds it to the driver it replaces.
