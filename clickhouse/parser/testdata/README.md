# ClickHouse parser test data

## `compat.tsv`: ClickHouse 26.8.38.9 verdicts

`compat.tsv` holds one statement for each alternative of ClickHouse's official
ANTLR 4 grammar (`../../testdata/ClickHouseParser.g4` and
`../../testdata/ClickHouseLexer.g4`, Apache-2.0, vendored byte-for-byte from
`ClickHouse/ClickHouse` commit `71391c9b2bbc4d2c104996451586a396138a4785` at
tag `v26.8.38.9-lts`), followed by its near-misses, with omni's verdict and the
pinned `clickhouse/clickhouse-server:26.8.38.9` server's (`../../testdata/pins.tsv`).
`TestCompatibility` (`../compat_test.go`) reads `compat.tsv` through
`enginetest/compat` offline without Docker; `TestGrammarVendorDigest` verifies
the SHA-256 of the vendored grammar files against the `grammar` row of
`../../testdata/pins.tsv`.

| Rows | Total | Agree | Only ClickHouse accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 1,207 | 1,207 (100.0%) | 0 | 0 |
| Near-misses | 2,341 | 2,341 (100.0%) | 0 | 0 |
| **All** | **3,548** | **3,548 (100.0%)** | **0** | **0** |

To record or verify against the pinned ClickHouse container:

```bash
OMNI_RECORD=1 go test -timeout 20m ./clickhouse/parser/ -run '^TestCompatibility$'
OMNI_RECORD=verify go test -timeout 20m ./clickhouse/parser/ -run '^TestCompatibility$'
```

## How the live judge works (`EXPLAIN AST`)

ClickHouse parses a single query per native-protocol request and rejects
multi-statement strings with error code `62` (`SYNTAX_ERROR`). `compatJudge`
uses `EXPLAIN AST` on the live `clickhouse/clickhouse-server:26.8.38.9`
container so every DDL, DML, query, and admin statement is parsed by
ClickHouse's real C++ parser without executing it or mutating server state.

Because ClickHouse's `ParserExplainQuery` accepts optional bare `key = value`
settings between `EXPLAIN AST` and the inner query (such as `graph = 1`),
`clickhouseVerdict` classifies each row fail-closed in up to three steps:

1. **Bare `EXPLAIN AST <stmt>`:** catches syntax errors immediately (`Code 62 -> reject`),
   including leading-comma fragments such as `, a = 1 SELECT 1` that would
   otherwise continue a settings list.
2. **Fixed-setting `EXPLAIN AST graph = 0 <stmt>`:** catches fragments whose
   prefix would be absorbed as an `EXPLAIN AST` setting by Step 1 (such as
   `graph = 1 SELECT 1` or `EXPLAIN = 1e1 SELECT *` -> `Code 62 -> reject`).
   When both Step 1 and Step 2 succeed (`err == nil`), the verdict is `accept`.
3. **EOF probe `EXPLAIN AST graph = 0 <stmt>\n; SELECT 1`:** only when both
   Step 1 and Step 2 return the same allowlisted ClickHouse semantic error code
   (`36` `BAD_ARGUMENTS`, `56` `UNKNOWN_STORAGE`, `115` `UNKNOWN_SETTING`,
   `376` `CANNOT_PARSE_UUID`, `476` `NOT_IMPLEMENTED`, or `558`
   `EXPECTED_ALL_OR_DISTINCT`), the probe separates:
   - **Post-parse checks** (e.g. unknown `SETTINGS` name `115`, unknown table
     engine `56`, bare `UNION` with empty `union_default_mode` `558`): the
     probe fails at `; SELECT 1` with `Code 62 (SYNTAX_ERROR)` -> `accept`.
   - **Parse-time rejections** (e.g. invalid `UUID` literal `376`, `WITH TIES`
     without `ORDER BY` `476`, window frame start `UNBOUNDED FOLLOWING` `36`):
     the probe aborts inside the parser with the same code before reaching EOF
     -> `reject`.
   - Any transport error, timeout, non-allowlisted server error, or inconsistent
     probe fails the recording rather than guessing.

## Remaining disagreements

None — all 3,548 rows (1,207 grammar alternatives and 2,341 near-misses) agree
with ClickHouse 26.8.38.9 (`0` ClickHouse-only, `0` omni-only).
