# Trino Oracle

| Field | Value |
|---|---|
| **engine** | `trino` (`storepb.Engine_TRINO`, value 24; no separate Presto engine) |
| **tier** | **2 — Docker**, recorded. The server's answers are committed: `trino/parser/testdata/compat.tsv` for syntax, and a tape per test (`testdata/tapes/<Test>.tape`) for the rest. Tests replay them and start no server. |
| **version** | **Trino 482**, pinned by digest in `trino/testdata/pins.tsv`. The tests were first written against 481 (`trinodb/trino:latest`, 2026-06-02). |
| **access** | `OMNI_RECORD=1 go test ./trino/...` starts the pinned image through testcontainers and records. |
| **client** | Go REST client at `trino/internal/testcontainer` (`Oracle.CheckSyntax`). `Open(t)` replays the test's tape, or records it; `Judge(t)` judges `compat.tsv` in a recording. |
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
  `starting:false`), once per package in a recording. A replay needs none.
- **State.** A statement that changes something (an `updateType`) runs to the
  end, so the next finds its change made, and a recording empties the `memory`
  catalog before each top-level test, so a test recorded alone reads what it
  reads in a run of its package.
