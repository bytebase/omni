# GoogleSQL Oracle

Empirically established 2026-06-02. The oracle converts grammar correctness from agent-judgment into a mechanical accept/reject comparison, consumed by the differential harness and every grammar node.

## Summary

| Field | Value |
|---|---|
| engine | `googlesql` (serves bytebase `BIGQUERY` + `SPANNER`) |
| tier | **2 — Docker** (Cloud Spanner emulator) |
| dialect coverage | **Spanner GoogleSQL only** — a *subset* of the full GoogleSQL union (see caveat) |
| access | gRPC `localhost:9010` (data + admin). **REST :9020 is unusable — see below.** |
| client | Go `cloud.google.com/go/spanner` v1.91.0 + `admin/database/apiv1`, env `SPANNER_EMULATOR_HOST=localhost:9010` |
| image | `gcr.io/cloud-spanner-emulator/emulator:latest` |
| digest | `sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2` (image built 2026-05-30, pulled 2026-06-02) |
| syntax_reject_signal | gRPC `InvalidArgument` **AND** message prefix `Syntax error:` (queries/DML) or `Error parsing Spanner DDL statement:` (DDL) |

## Standing it up

```bash
docker run -d --name spanner-emulator -p 9010:9010 -p 9020:9020 \
  gcr.io/cloud-spanner-emulator/emulator
```

Bootstrap (once per emulator process; state is in-memory and lost on container restart). Done via the Go admin client or REST :9020 *for setup only*:
- project `test-project`, instance config `emulator-config`
- instance `test-instance`, database `testdb`

The differential harness owns this bootstrap and a scratch table for query/DML probes.

## ⚠️ REST gateway (:9020) cannot report errors

The emulator's grpc-gateway REST endpoint collapses **every** error — syntax and semantic alike — to:

```
HTTP 500  {"code": 13, "message": "failed to marshal error message"}
```

Valid SQL returns `HTTP 200` with rows; everything else is an indistinguishable 500. **REST is therefore only usable for the binary "did it 200?" and cannot classify rejects.** The oracle/harness MUST use the **gRPC** API (:9010) via the Go client, which surfaces the full `google.rpc.Status` (code + message). REST is fine for instance/database/table *bootstrap* (no error classification needed there).

## Error classification (verified empirically)

The gRPC **code is NOT a sufficient discriminator** — both syntax and semantic failures return `InvalidArgument (3)` for queries/DML. Classification is by **message prefix**:

| Statement kind | API | Grammar REJECT signal | Grammar-ACCEPTED (semantic / unsupported) examples |
|---|---|---|---|
| Query / DML | `ExecuteSql` (gRPC) | `InvalidArgument` + msg starts `Syntax error:` | `InvalidArgument` `Table not found:` · `Unrecognized name:` · `Column X is not present` · `QUALIFY is not supported` |
| DDL | `UpdateDatabaseDdl` (gRPC) | `InvalidArgument` + msg starts `Error parsing Spanner DDL statement:` | `FailedPrecondition` (`Duplicate name in schema`, `does not reference parent key column`) · `Internal` (`GOOGLESQL_RET_CHECK failure` on unknown type — emulator quirk) |

Decision rule the harness uses (**fail-closed** — an oracle must never let an
infra failure masquerade as ACCEPT, or it silently corrupts every grammar
node's differential):

```
classify(err):
  if err == nil:                                       ACCEPT (ok)
  if err is not a gRPC status (dial/context error):    ERROR  (infra)   # cannot decide
  code, msg = grpcStatus(err)
  if code==InvalidArgument and msg startswith "Syntax error:":               REJECT (syntax)
  if code==InvalidArgument and msg startswith "Error parsing Spanner DDL statement:": REJECT (syntax)
  if code in {InvalidArgument, FailedPrecondition, NotFound, AlreadyExists}: ACCEPT (semantic)  # parsed, then a semantic/feature result
  if code==Internal and msg contains "GOOGLESQL_RET_CHECK":                  ACCEPT (semantic)  # emulator quirk on unknown type (parsed)
  else:  # Unavailable / DeadlineExceeded / Canceled / Aborted / ResourceExhausted / Unknown / generic Internal
                                                                            ERROR  (infra)   # emulator down/slow/crashed => NOT a grammar verdict
```

Consumers MUST treat a verdict of `error` as "the oracle could not decide" — fail the test / drop the fixture, **never** fold it into accept or reject.

Verified probes (2026-06-02, against emulator @ above digest):

```
SELECT 1                                  -> ACCEPT
SELECT id, name FROM T                     -> ACCEPT
SELEC 1                                    -> REJECT  Syntax error: Unexpected identifier "SELEC"
SELECT FROM                                -> REJECT  Syntax error: SELECT list must not be empty
@@@ garbage !!                             -> REJECT  Syntax error: Unexpected "@@"
SELECT * FROM no_such_table                -> ACCEPT  (semantic: Table not found)
SELECT bogus FROM T                        -> ACCEPT  (semantic: Unrecognized name)
INSERT INTO T (id,name) VALUES (1,'a')     -> ACCEPT
INSERT INTO T (id,name) VALUE (1,'a')      -> REJECT  Syntax error: Unexpected keyword VALUE
CREATE TABLE T2 (x INT64) PRIMARY KEY (x)  -> ACCEPT
CREATE TABL T3 ...                         -> REJECT  Error parsing Spanner DDL statement
CREATE TABLE T4 (x INT64) PRIMARY KEY x    -> REJECT  Error parsing Spanner DDL statement (missing parens)
CREATE TABLE T (id INT64) PRIMARY KEY (id) -> ACCEPT  (semantic: Duplicate name / FailedPrecondition)
```

## Non-mutating validation (harness note)

To validate queries/DML without side effects, use **`QueryMode = PLAN`** (Go: `tx.AnalyzeQuery` for SELECT; the low-level `ExecuteSqlRequest{QueryMode: PLAN}` for DML inside a read-write txn). PLAN runs parse + semantic analysis and returns the query plan or the same `Syntax error:`/semantic error, but does not execute or mutate. DDL is validated by submitting to `UpdateDatabaseDdl` against a scratch database (it does mutate schema — the harness uses throwaway table names or a fresh database per batch).

## ⚠️ Dialect caveat — Spanner is a SUBSET of GoogleSQL

The bytebase `googlesql` parser must accept the **union** of BigQuery + Spanner GoogleSQL, but this oracle only speaks **Spanner's** dialect. Consequences:

1. **BigQuery-only constructs are NOT authoritative against this oracle.** A Spanner reject of a BigQuery-only form is expected and must NOT be trusted as a grammar verdict. Two failure shapes observed:
   - *Parsed-but-unsupported* — e.g. `QUALIFY` → `QUALIFY is not supported` (classified ACCEPT by the rule above; harmless).
   - *Hard syntax reject* — Spanner's grammar can't tokenize it at all → false `Syntax error:`. Expected for BigQuery scripting (`DECLARE`, `EXECUTE IMMEDIATE`, `BEGIN…EXCEPTION…END`), `EXPORT DATA` / `LOAD DATA`, `CREATE MODEL`, capacity/reservation DDL, `PIVOT`/`UNPIVOT`, BigQuery-style partitioning/clustering options, etc.
2. **Triangulation for BigQuery-only forms** (per the approved plan): treat as canonical-accept when ≥2 of these agree, ignoring the Spanner emulator:
   - official BigQuery docs corpus → `truth1/bigquery/`
   - the legacy `GoogleSQLParser.g4` (it parses the form) → `antlr_rules.md`
   - the third-party `bq-parser` ANTLR grammar (`/Users/h3n4l/OpenSource/bq-parser`)
3. **Divergence ledger.** Any case where the Spanner emulator and the docs/grammar disagree on a form the parser must accept is recorded as a `divergence` row `{ node, case, legacy_behavior, new_behavior, claim, evidence[], confidence, status }`. Each grammar node tagged BigQuery-only carries: "Spanner-emulator verdict is non-authoritative; triangulated."

### Practical routing for the harness

| Form class (from truth1 / antlr_rules tags) | Oracle |
|---|---|
| shared GoogleSQL core (SELECT, joins, set-ops, basic DDL/DML, expressions, types) | Spanner emulator (authoritative) |
| Spanner-only (INTERLEAVE, CHANGE STREAM, sequences, statement hints `@{}`, role DDL) | Spanner emulator (authoritative) |
| BigQuery-only (scripting, EXPORT/LOAD DATA, MODEL, PIVOT/UNPIVOT, BQ partition/cluster, reservations) | triangulation (docs + .g4 + bq-parser); emulator verdict ignored |

## Tiers considered

Per `oracle-setup.md`, probe strongest-first. Chosen: tier 2 (Docker) — matches the local-Docker workflow, zero credentials, fully scriptable over gRPC. Rejected for now: BigQuery dry-run (tier ~1, truest for the BigQuery dialect, but needs GCP credentials — revisit if BigQuery-only coverage gaps prove costly); ZetaSQL reference build (tier 3, canonical for the union, but heavy Bazel/C++ build). The triangulation backup above covers the Spanner-subset blind spot without those.
