# databricks

Databricks SQL support in `omni`:

- `databricks/parser`: hand-written recursive-descent parser implementing the
  Apache Spark `v4.0.5` SQL grammar (`SqlBaseParser.g4` + `SqlBaseLexer.g4`,
  plus the `QUALIFY` clause), backed by the vendor-grammar compatibility suite
  in [`parser/testdata/README.md`](parser/testdata/README.md). Statements
  currently return statement-level `ParsedStmt` spans (`Kind`, `Target`, `Loc`)
  rather than full clause/expression AST trees. Databricks-only SQL extensions
  outside Spark 4.0 — including `USE CATALOG`, `CREATE CATALOG`,
  `CREATE SCHEMA ... MANAGED LOCATION`, `ALTER TABLE ... OWNER TO` /
  `ADD CONSTRAINT` / `SET ROW FILTER`, `PRIMARY KEY` constraints,
  `OPTIMIZE ... ZORDER BY`, `VACUUM`, `COPY INTO`, `CREATE MATERIALIZED VIEW`,
  `VARIANT` colon paths (`a:b.c`), and `SHOW GRANTS ON TABLE` — are not yet
  implemented and are tracked for follow-up work against the pinned Databricks
  workspace.
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
