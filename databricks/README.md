# databricks

Driver-only tree. Databricks has no omni parser, analysis or catalog yet, so
`databricks/driver` is the whole tree: it reaches a workspace through the
Databricks SDK, syncs Unity Catalog into `omni/metadata`, runs statements on a
SQL warehouse, and dumps the server's own DDL. Bytebase's lexical splitter and
`EXPLAIN FORMATTED` rewrite are carried in the driver until a `databricks`
analysis tree holds them, as `clickhouse/driver` and `hive/driver` carry theirs.

The package was verified against a scratch workspace on 2026-11-17. Its live
tests passed, `driver/testdata/tapes` holds the recordings, `driver/testdata/sync`
holds the synced catalog, and `databricks/testdata/pins.tsv` pins the workspace,
SQL warehouse and catalog. The live tests still skip without `DATABRICKS_DSN`.
See [`driver/VERIFY.md`](driver/VERIFY.md) for the recording and parity record.
