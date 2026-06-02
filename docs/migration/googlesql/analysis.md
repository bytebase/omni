# GoogleSQL Migration Analysis (v2 ANALYZE)

Synthesis + index for the omni `googlesql` engine migration. The detailed
artifacts are the deliverables; this file ties them together and hands off to
`omni-engine-v2-planning`.

- **Engine:** omni `googlesql` (hand-written recursive-descent parser, full legacy parity target).
- **Scope:** **both** bytebase consumers — `BIGQUERY` (enum 19) and `SPANNER` — which both import `github.com/bytebase/parser/googlesql`. One omni parser must serve both.
- **Migration target:** full coverage of the legacy ANTLR `bytebase/parser/googlesql` grammar; bytebase consumption *prioritizes* order, never *scopes down* legacy features.
- **Analyzed:** 2026-06-02. Oracle tier: **docker** (Spanner emulator). Provenance: the legacy grammar is a **hand-port of ZetaSQL** (Google's open-source GoogleSQL reference).

## Artifact index

| Artifact | What it is | Size / counts |
|---|---|---|
| `truth1/bigquery/` | documented BigQuery GoogleSQL forms (9 files + INDEX) | **165 forms** |
| `truth1/spanner/` | documented Spanner GoogleSQL forms (7 files + INDEX) | **185 forms** |
| `antlr_rules.md` | legacy ANTLR grammar catalog (truth2, a *hint*) | 1069 lines; **625 rules**, 53 SQL + 12 scripting statements |
| `contract.md` | bytebase consumption surface (must-work; marks P0) | 321 lines; **6 P0 handlers** |
| `oracle.md` | oracle tier + access + verified syntax-reject signal | Spanner emulator, gRPC :9010 |
| `harness/googlesql-spanner/` | differential oracle (code; reused by grammar nodes) | own Go module; builds + tests green |
| legacy examples | `/Users/h3n4l/OpenSource/parser/googlesql/examples/` (Corpus A) | **72 `.sql`** (mostly ZetaSQL testdata) |

## Grammar (truth2 highlights)

ZetaSQL hand-port; **625 parser rules**, **307 keyword tokens → 211 non-reserved / 97 reserved**, 43 operator/punct tokens. **53 active top-level SQL statements** + **12 procedural-script forms** (IF/CASE/WHILE/LOOP/REPEAT/FOR-IN/DECLARE/BREAK/CONTINUE/RETURN/RAISE/BEGIN-END). Notable: a **40-rule GQL / property-graph** sub-language; rich expression precedence chain, joins, set-ops, CTEs, PIVOT/UNPIVOT, window, hints, OPTIONS, and a generic-entity extensibility mechanism (used for CHANGE STREAM, named schemas).

**10 TODO/FIXME gaps** in the legacy grammar (keep-until-ruled): `DEFINE MACRO` only error-rejected; **pipe syntax `|>` is lexed but never parsed** (`query_without_pipe_operators`); `with_group_rows` accepts no body; lambda arg-kind check skipped; `@<int>` hint missing ABORT_CHECK; `SCRIPT_LABEL` token unhandled; error messages emit literal `<KEYWORD>`. Many alternatives exist only to emit ZetaSQL-style errors.

`bq-parser` (third-party) is a ~25-rule WIP stub — **not** a viable reference; useful only as a hint for a few array-access keywords (`OFFSET/ORDINAL/SAFE_OFFSET/SAFE_ORDINAL`, `TREAT`, `WITHIN`, `LATERAL`, `FETCH`).

## Contract (P0 surface — small and clean)

The legacy `googlesql` import is confined to `backend/plugin/parser/{bigquery,spanner}/` (only external touchpoint: a blank `init()` import in `server/ultimate.go`). Registered handlers — **the entire must-work surface**:

| Handler | bigquery | spanner | Notes |
|---|---|---|---|
| `RegisterGetQuerySpan` (GetQuerySpan) | ✅ | ✅ | drives field lineage + masking; embeds the `queryType` listener |
| `RegisterSplitterFunc` (SplitSQL) | ✅ | ✅ | **divergent**: BigQuery = lexer/token split; Spanner = full parse-tree split (BEGIN/END/CASE) |
| `RegisterDiagnoseFunc` (Diagnose) | ✅ | ✅ | syntax diagnostics for the editor |

**Explicitly out of scope** (no registration / gated `false` in `backend/common/engine.go`): auto-complete, SQL review/lint, statement report, **statement ranges**, prior backup/restore, resource-change extraction, **query validation** (lives in `standard/query.go`, regex-based, shared — not googlesql). So the omni googlesql parser needs **no** `completion/`, `deparse/`, `quality/`, or statement-ranges packages.

**bigquery vs spanner divergences:** (1) split impl (above); (2) metadata resolution — bigquery per-dataset (`project.dataset.table`, no schema layer), spanner per-schema under one DB; (3) system schemas — bigquery `INFORMATION_SCHEMA` only, spanner adds `SPANNER_SYS`; default join-type fallthrough differs. Otherwise `query_type.go` + the query-span extractor core are byte-identical → share omni code, parameterize the deltas.

## Oracle (see oracle.md for full detail)

Cloud Spanner emulator (Docker), gRPC **:9010** via the Go client. **Verified** syntax-reject signal: gRPC `InvalidArgument` is *not* a discriminator — match the message **prefix** `Syntax error:` (query/DML) or `Error parsing Spanner DDL statement:` (DDL). REST :9020 is unusable (collapses all errors to HTTP 500). The harness measures **parseability**, not feature-support — `"X is not supported"` / `"Statement not supported"` mean the grammar *accepted* the input.

**⚠ Spanner is a SUBSET of the BigQuery+Spanner union.** BigQuery-only forms are **non-authoritative** against this oracle and must be triangulated against `truth1/bigquery/` + the legacy `.g4` (+ `bq-parser` hints). Routing table in oracle.md.

## Dialect matrix (for P0/P1 tagging)

- **Shared core** (P0 backbone): SELECT/query stack, joins, set-ops, CTEs, full expression + type + literal + identifier grammars, DML (INSERT/UPDATE/DELETE/MERGE), basic DDL, transactions, CALL, EXECUTE IMMEDIATE, hints, OPTIONS, statement-level `@{}` hints (both dialects).
- **Spanner-only:** `INTERLEAVE IN [PARENT] … ON DELETE`, `NULL_FILTERED` indexes, generated `AS (expr) STORED`, sequences, Spanner `ALTER COLUMN`/`SET ON DELETE`, CHANGE STREAM + named schemas (via generic-entity), `INSERT OR UPDATE|IGNORE`, query hints `@{FORCE_INDEX=…}`.
- **BigQuery-only:** procedural scripting, `EXPORT DATA`/`EXPORT MODEL`/`LOAD DATA`/`CLONE DATA`, ML models, REMOTE/CONNECTION functions, differential privacy (`WITH DIFFERENTIAL_PRIVACY`), snapshots, SEARCH/VECTOR indexes, materialized/approx/recursive views, dashed/slashed table paths, `#` comments, `AT SYSTEM TIME`, property-graph + GQL.

## P0 vs P1 coverage target

**P0 — required before bytebase cutover** (drives GetQuerySpan / Split / Diagnose for what BigQuery & Spanner users actually write):
lexer · parser entry + error collection · statement splitter (**both** lexer-split and parse-tree-split variants) · identifiers/paths/normalization · data types · expressions (incl. window, subqueries, array/struct, CASE/CAST) · SELECT core · joins · set-ops · CTEs · DML (INSERT incl. Spanner upserts / UPDATE / DELETE / MERGE) · core DDL (CREATE/ALTER/DROP TABLE/VIEW/INDEX/SCHEMA/DATABASE incl. Spanner INTERLEAVE/CHANGE STREAM/SEQUENCE and BigQuery PARTITION/CLUSTER/OPTIONS) · query-span extractor (×2 dialect resolution) · query-type classifier · diagnose.

**P1 — legacy parity** (parse-only; not consumed today): procedural scripting · EXPORT/LOAD/CLONE DATA · ML/MODEL · differential privacy · snapshots · SEARCH/VECTOR index · materialized/approx views · generic-entity DDL · GQL/property-graph · pipe `|>` (currently a grammar gap) · the remaining 53-statement long tail (GRANT/REVOKE, ANALYZE, transactions, CALL, etc. — promote to P0 any that appear in real corpora).

## Preliminary architecture

```
googlesql/
  ast/         # AST node types (ZetaSQL-shaped; recursive-descent friendly)
  parser/      # lexer + recursive-descent parser + 2 splitter variants
    testdata/
      legacy/    # lifted from parser/googlesql/examples (72 files)
      official/  # from truth1 (per-form fixtures, snowflake-style layout)
  analysis/    # query span / field lineage + query-type, dialect-parameterized
```

No `completion/`, `deparse/`, `quality/`, `semantic/` packages (not consumed). Reference engines: `snowflake/` (closest GoogleSQL-style surface: QUALIFY, PIVOT/UNPIVOT, EXECUTE IMMEDIATE, rich DDL) and `pg/` (most mature ast/parser split). Differential testing reuses `harness/googlesql-spanner` (Spanner-authoritative for shared + Spanner-only forms; triangulated for BigQuery-only).

## Divergences / risks → divergence ledger (for PLAN to seed)

Oracle-verified during ANALYZE (the Spanner truth1 crawl hit Claude API overload (HTTP 529); some Spanner forms were written from training knowledge — so these were re-checked against the live emulator):

| Case | truth1 claim | Oracle (Spanner emulator) | Disposition |
|---|---|---|---|
| `QUALIFY` (Spanner) | "supported" | parses, **feature-rejected** `QUALIFY is not supported` | corpus feature-claim wrong; **parser-union must still parse it** (BigQuery). Confidence: high |
| `MERGE` (Spanner) | not supported | parses, `Statement not supported: MergeStatement` | corpus correct; grammar parses MERGE. high |
| `WITH RECURSIVE` (Spanner) | not supported | parses, `RECURSIVE is not supported in the WITH clause` | corpus correct; grammar parses it. high |
| `INSERT OR UPDATE` / `OR IGNORE` | supported | **accepted** (`Table not found`) | confirmed Spanner upsert grammar. high |
| inline column `PRIMARY KEY` | "trailing-only" | **accepted** by emulator | divergence — verify vs `.g4` + production Spanner at node time. medium |
| `CREATE TABLE` with no PK | (PK required) | **reject** `Error parsing Spanner DDL` | confirmed PK required. high |
| INTERLEAVE / CHANGE STREAM / SEQUENCE / `ALTER ADD COLUMN` / `VIEW SQL SECURITY` | supported | **accepted** | confirmed. high |

Other risks: (a) **Spanner-subset oracle** — BigQuery-only forms blind-spot, mitigated by triangulation; (b) **pipe `|>`** lexed-not-parsed in legacy grammar — decide keep/drop in PLAN; (c) corpus feature-accuracy on un-reverified Spanner forms — node-level differential will catch them.

## Hand-off to PLAN

Ready for `omni-engine-v2-planning` (build DAG + seed SQLite store). Seed: engine=`googlesql`, repo=`bytebase/omni`, **oracle_tier=`docker`**, scope=both. P0 set above is the cutover gate (GetQuerySpan/SplitSQL/Diagnose ×{bigquery,spanner}); P1 is the long tail incl. GQL and scripting.
