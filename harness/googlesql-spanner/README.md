# googlesql-spanner harness

The differential **oracle** side for the omni `googlesql` parser migration. It
runs SQL through a live **Cloud Spanner emulator** (gRPC) and emits a compact
accept/reject verdict per statement, so grammar-node tests can diff omni's
parser accept/reject against a real GoogleSQL implementation.

This is the GoogleSQL analogue of `harness/mssql-scriptdom`: an external
reference engine omni diffs against. Unlike ScriptDOM (which exposes an AST),
the emulator only answers "did it parse?", so this harness emits a **verdict +
error classification**, not an AST shape.

It is a **separate Go module** (own `go.mod`) on purpose — it pulls the heavy
`cloud.google.com/go/spanner` gRPC client, which must stay out of omni's main
module. The parser packages never import it; only test code shells out to it.

Full rationale, the empirically-verified error signal, and the dialect caveat
live in [`docs/migration/googlesql/oracle.md`](../../docs/migration/googlesql/oracle.md).

## Prerequisite: the emulator

```bash
# Pinned by digest: classification keys mostly on the gRPC code, but the
# query/DML "Syntax error:" prefix is message-based, so bumping the image
# requires a green `go test` (live TestOracleLive) to re-validate the format.
docker run -d --name spanner-emulator -p 9010:9010 -p 9020:9020 \
  gcr.io/cloud-spanner-emulator/emulator@sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2
```

State is in-memory; the harness creates a fresh scratch instance+database on
each run. Always talk to the **gRPC** port (`:9010`) — the REST gateway
(`:9020`) cannot report error classes (see oracle.md).

## Build & run

```bash
# batch line mode (default): one base64-encoded SQL per line on stdin,
# one JSON verdict per line out (output line N <-> input line N).
printf '%s\n' "$(printf 'SELECT 1' | base64)" \
  | SPANNER_EMULATOR_HOST=localhost:9010 go run .
# => {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}

# single mode: all of stdin is one statement.
echo 'CREATE TABL x' | SPANNER_EMULATOR_HOST=localhost:9010 GOOGLESQL_HARNESS_LINE=0 go run .
# => {"verdict":"reject","kind":"ddl","reason":"syntax","code":"InvalidArgument","message":"Error parsing Spanner DDL statement: ..."}
```

## Verdict JSON

```
{"verdict":"accept|reject|error","kind":"query|dml|ddl","reason":"ok|syntax|semantic","code":"<grpc code>","message":"<truncated>"}
```

- `verdict` — `accept` / `reject` (the grammar verdict) or **`error`** (the
  oracle could not decide).
- `kind` — how the statement was routed (by leading keyword): query → executed
  (first row only); dml → run in an aborted read-write txn (no mutation);
  ddl → submitted to `UpdateDatabaseDdl` on the scratch db.
- `reason` — `ok` (clean parse), `syntax` (grammar reject), `semantic` (parsed
  fine, failed later: unknown table/column, "X is not supported", etc.), or for
  `error`: `infra` (emulator down/slow/crashed, timeout, cancellation), `empty`
  (blank or comment/hint-only input), `blank line`, `bad base64`.

> **`error` is fail-closed.** Consumers MUST treat `error` as "no verdict" —
> fail the test / drop the fixture — and **never** fold it into accept or
> reject. An infra blip that silently became `accept` would mask real parser
> bugs in every node that trusts this oracle.

## Classification (see oracle.md for the verified evidence)

The gRPC **code is not a discriminator** (syntax and semantic are both
`InvalidArgument`). The message **prefix** decides accept-vs-reject; codes that
aren't a parser/semantic result become `error`:

| outcome | verdict |
|---|---|
| `InvalidArgument` + msg `Syntax error:` (query/DML) | **reject** |
| `InvalidArgument` + msg `Error parsing Spanner DDL statement:` (DDL) | **reject** |
| `InvalidArgument` / `FailedPrecondition` / `NotFound` / `AlreadyExists` (other msgs: `Table not found`, `Unrecognized name`, `… is not supported`, `Duplicate name`) | **accept** (grammar parsed) |
| `Internal` + msg contains `GOOGLESQL_RET_CHECK` (emulator quirk on unknown type) | **accept** (parsed) |
| `Unavailable` / `DeadlineExceeded` / `Canceled` / `Aborted` / `ResourceExhausted` / `Unknown` / generic `Internal` / non-gRPC error | **error** (infra — fail closed) |

## ⚠️ Spanner is a SUBSET of GoogleSQL

The emulator speaks Spanner's dialect; the omni parser must accept the BigQuery
**+** Spanner union. A **reject of a BigQuery-only form is NOT authoritative**
(e.g. `DECLARE`, `EXECUTE IMMEDIATE`, `EXPORT DATA`, `CREATE MODEL`,
`PIVOT`/`UNPIVOT`, GQL). This harness reports only the Spanner verdict — the
caller decides authoritativeness using the form's dialect tag (from
`truth1/` / `antlr_rules.md`) and triangulates BigQuery-only forms against the
docs corpus + the legacy `.g4`.

## How grammar nodes consume it

A future `googlesql/parser` diff test (build-tagged, like
`mssql/parser/scriptdom_harness_test.go`) will:

1. base64-encode each corpus statement (one statement, no trailing `;` — the
   DDL path rejects a terminator), pipe the batch to this harness;
2. parse the same statement with omni's googlesql parser;
3. on `verdict == "error"`, **fail loudly** (the oracle didn't decide — fix the
   emulator / fixture, don't proceed) — never treat it as accept or reject;
4. otherwise assert `omniAccepts == (verdict == "accept")` — **except** for
   statements tagged BigQuery-only, where a harness `reject` is ignored and the
   docs/`.g4` triangulation is authoritative;
5. record any genuine disagreement in the divergence ledger.

## Tests

```bash
# pure unit tests (classify + kind routing) — no emulator needed:
go test ./...
# include the live integration test (auto-skips if the emulator is down):
SPANNER_EMULATOR_HOST=localhost:9010 go test ./...
```
