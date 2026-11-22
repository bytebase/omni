# databricks

Databricks SQL support in `omni`:

- `databricks/parser`: hand-written recursive-descent parser for Databricks /
  Spark SQL, backed by the vendor-grammar compatibility suite in
  [`parser/testdata/README.md`](parser/testdata/README.md) (Apache Spark
  `v4.0.5` `SqlBaseParser.g4` + `SqlBaseLexer.g4`).
- `databricks/analysis`: lexical SQL statement splitter, row-returning
  statement classifier, and `EXPLAIN FORMATTED` rewrite.
- `databricks/driver`: Databricks SDK driver syncing Unity Catalog into
  `omni/metadata`, running statements on a SQL warehouse, and dumping server
  DDL.

The driver was verified against a scratch workspace on 2026-11-17. Its live
tests passed, `driver/testdata/tapes` holds the recordings, `driver/testdata/sync`
holds the synced catalog, and `databricks/testdata/pins.tsv` pins both the
upstream Spark SQL grammar (`default`) and the Databricks workspace (`driver`).
The live driver tests still skip without `DATABRICKS_DSN`.
See [`driver/VERIFY.md`](driver/VERIFY.md) for the recording and parity record.
