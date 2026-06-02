# Trino Oracle

| Field | Value |
|---|---|
| **engine** | `trino` (`storepb.Engine_TRINO`, value 24; no separate Presto engine) |
| **tier** | **2 — Docker** (live local Trino server). Far stronger than antlr-fallback; flips correctness from judgment to mechanical comparison. |
| **version** | **Trino 481** (`trinodb/trino:latest`, image id `5b5e0a97f599`, pulled 2026-06-02). Matches the `trino.io/docs/current` corpus used for truth1. |
| **access** | `docker run -d --name trino-oracle -p 18080:8080 trinodb/trino:latest` |
| **client** | Go REST client at `trino/internal/trinooracle` (`Oracle.CheckSyntax`). Connects to `$TRINO_ORACLE_URL` (default `http://localhost:18080`). `StartContainer` auto-starts a disposable server via testcontainers for CI. |
| **syntax_reject_signal** | REST statement error object with **`errorName == "SYNTAX_ERROR"`** (`errorCode` 1, `errorType` `USER_ERROR`). This — and only this — means Trino's parser rejected the input. |
| **semantic_error_examples** | Grammar **accepted**, failed later: `TABLE_NOT_FOUND` (46), `CATALOG_NOT_FOUND` (44), `COLUMN_NOT_FOUND` (47), `NOT_SUPPORTED` (13, e.g. memory-connector `MERGE`), `TYPE_MISMATCH`. Any non-`SYNTAX_ERROR` outcome (including success) = accepted. |
| **catalogs** | `memory` (supports CREATE/INSERT → DDL/DML probing), `tpch` + `tpcds` (real tables/columns → query-span & lineage probing), `system`, `jmx`. Harness defaults to `memory.default`. |

## How classification works

Trino runs parse → analyze → plan → execute. A `SYNTAX_ERROR` is raised during
parse, before analysis. The oracle submits the statement over `/v1/statement`,
follows `nextUri`, and reads `stats.state`:

- `error.errorName == "SYNTAX_ERROR"` → **rejected**.
- any other `error.errorName`, or reaching `RUNNING`/`FINISHED` → **accepted**.

Once the query reaches `RUNNING`, parsing+planning have demonstrably succeeded,
so the client cancels (`DELETE nextUri`) to avoid streaming result data. This
keeps even `SELECT … FROM tpch.sf1.*` checks sub-second.

## Notes

- **Robust to side effects.** Because classification keys on `SYNTAX_ERROR`, the
  oracle may *execute* a DDL/DML statement against the disposable `memory`
  catalog (e.g. actually create a table). That never changes the syntax verdict.
  Grammar-node authors who want a pure parse with no execution can wrap an
  EXPLAIN-able statement in `EXPLAIN …` (not universal — many utility/DDL forms
  are not EXPLAIN-able), but it is not required for syntax classification.
- **Boot time.** The Trino JVM takes ~40–60s to become ready (`/v1/info` →
  `starting:false`). Prefer one long-lived container + `TRINO_ORACLE_URL` for
  iterative work; `StartContainer` (≈1 min) is for self-contained CI.
- **Gating.** `TestOracleClassification` (15 cases) skips cleanly in `-short`
  mode or when no oracle is reachable, so `go test ./...` stays green without a
  running Trino.
- **Verified.** `go vet ./trino/...` clean; self-test green against Trino 481
  (valid / Trino-specific / semantic-accept / syntax-reject all classified
  correctly).
