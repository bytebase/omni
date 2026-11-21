# Redis parser test data

## `compat.tsv`: Redis 7.4.12 command call-shape rows

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/RedisCommands.g4`, then its near-misses (one token dropped, the
end cut off, a stray word inserted, or a command keyword swapped), plus
hand-added `suite:lexical` edge-case rows. Each row records omni's verdict
(`redis/parser.ParseCommand`) and the pinned Redis server's verdict
(`redis:7.4.12-alpine`, `default` role in `../../testdata/pins.tsv`, leaving the
`driver` role untouched for `redis/driver`'s tapes).

`TestCompatibility` (`../compat_test.go`) reads `compat.tsv` offline through
`enginetest/compat`, regenerates the grammar rows from `RedisCommands.g4` and
`../../testdata/tokens.tsv`, and starts no Docker container during normal test
runs.

| Rows | Total | Agree | Only Redis accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 390 | 390 | 0 | 0 |
| Near-misses | 978 | 978 | 0 | 0 |
| `suite:lexical` | 13 | 13 | 0 | 0 |
| All | 1,381 | 1,381 | 0 | 0 |

## Command table and grammar source

Redis has no SQL statement grammar; per §2.2 and §8 of
`docs/specs/2026-11-11-engine-compatibility-without-servers-design.md`, its rows
cover command names, container subcommands, exact/minimum arity call shapes, and
command-line token forms:

- **Command table snapshot (`../../testdata/commands.tsv`):** captured directly
  from `COMMAND` on the pinned `redis:7.4.12-alpine` server
  (`sha256:7f7d30b7fdd09467bac58198a1f166dd69789d348d4fe859cccb20e0cf287fbe`,
  `7.4.12`). Covers all **250** top-level commands (**235** standalone commands
  plus **15** container commands with **129** subcommands, including `COMMAND`
  which is both callable standalone with arity `-1` and a 7-subcommand
  container — **365** callable command/subcommand call shapes in total).
- **Grammar (`../../testdata/RedisCommands.g4`) and Go registry
  (`../commands_gen.go`):** generated deterministically from `commands.tsv` by
  `../cmd/commandgrammar`. Exact-arity commands (`arity > 0`) emit fixed
  argument counts so `!insert` near-misses exercise extra-argument rejection;
  variadic commands (`arity < 0`) emit minimum required arguments followed by
  `variadicArgs` so `!insert` exercises variadic extra-argument acceptance.
- **Token forms (`../../testdata/tokens.tsv`):** unquoted words (`WORD`),
  single-quoted verbatim strings (`SINGLE_QUOTED`), double-quoted strings with
  escapes (`DOUBLE_QUOTED`), `#`-prefixed argument text (`HASH_ARG`, since Redis
  has no comment syntax), integers (`INT_LITERAL`), and floats (`FLOAT_LITERAL`).

To regenerate `commands.tsv`, `RedisCommands.g4`, and `commands_gen.go` from the
pinned server after a pin bump, then re-record `compat.tsv`:

```bash
go run ./redis/parser/cmd/commandgrammar -from-live
OMNI_RECORD=1 go test ./redis/parser/ -run '^TestCompatibility$'
```

To verify every committed verdict read-only against the pinned live server:

```bash
OMNI_RECORD=verify go test ./redis/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit rows that now agree:

```bash
go test ./redis/parser/ -run '^TestCompatibility$' -update-compat
```

## How the live Redis 7.4.12 judge works

`compatJudge` (`../compat_test.go`) starts the pinned `redis:7.4.12-alpine`
container through `redis/internal/testcontainer.Judge` (verifying both the
container's image digest and reported `redis_version:7.4.12`) only when
`OMNI_RECORD=1` or `OMNI_RECORD=verify` is set.

Rather than executing destructive, mutating, or blocking commands (`SHUTDOWN`,
`FLUSHALL`, `BLPOP`, `SUBSCRIBE`, `MULTI`, `REPLICAOF`, `FAILOVER`, `MIGRATE`,
`WAIT`, `DEBUG`, etc.) against the server, the judge:

1. Splits the command line using an independent implementation of Redis's own
   inline argument splitter (`sdssplitargs`), rejecting empty lines, multi-line
   strings, unbalanced quotes, and pipe-containing first tokens such as
   `acl|list` (which Redis's command dispatcher rejects with `ERR unknown
   command`).
2. Queries the live pinned server with `COMMAND INFO <cmd>`.
3. Uses the server's live reply to validate command existence, container
   subcommand lookup (including `COMMAND`'s standalone arity `-1` vs subcommand
   dispatch), and exact (`arity > 0`) or minimum (`arity < 0`) argument count.

On startup, `verifyJudgeAgainstExecution` sends 35 safe read-only / non-mutating
command lines directly over raw TCP to Redis 7.4.12's inline protocol and
asserts 100% agreement with `judgeCommandLine` across quoting, escapes, fixed
arity, variadic arity, container commands, `COMMAND` standalone vs subcommand
dispatch, unknown commands/subcommands, and `acl|list` rejection.

## What the judge leaves out

- **Per-command option and type validation:** `COMMAND INFO` checks lexical
  command-line syntax, command/subcommand names, and Redis's `arity` contract
  (`ERR unknown command`, `ERR unknown subcommand`, `ERR wrong number of
  arguments`). Command-specific keyword options, numeric range/type parsing, and
  key-type checks happen inside each command handler at runtime and are outside
  call-shape syntax scope.
- **Shell-style unquoted escapes (`redis/driver` compatibility):** `redis/parser`
  matches `redis/driver` and `redis/analysis`'s shell-like tokenization for
  unquoted backslash escapes (`two\ words`), trailing-backslash errors (`a\`),
  and adjacent quote concatenation (`b"c d"e`), covered by unit tests in
  `../parser_test.go`.
- **Redis Stack module commands:** Core `redis:7.4.12-alpine` ships no loaded
  modules, so module-provided commands (`JSON.*`, `FT.*`, `TS.*`, `BF.*`) are
  not part of the pinned core command table.
