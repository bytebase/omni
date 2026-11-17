# databricks

Driver-only tree. Databricks has no omni parser, analysis or catalog yet, so
`databricks/driver` is the whole tree: it reaches a workspace through the
Databricks SDK, syncs Unity Catalog into `omni/metadata`, runs statements on a
SQL warehouse, and dumps the server's own DDL. Bytebase's lexical splitter and
`EXPLAIN FORMATTED` rewrite are carried in the driver until a `databricks`
analysis tree holds them, as `clickhouse/driver` and `hive/driver` carry theirs.

`databricks/driver` is a port of Bytebase's `backend/plugin/db/databricks` that
**no test has ever run against a workspace**. It starts no server, has no tapes,
no sync goldens and no `testdata`; the tests that need an account skip with
`DATABRICKS_DSN`. See [`driver/VERIFY.md`](driver/VERIFY.md) for the recording
plan, the pin row, the HTTP tape harness and the parity harness.
