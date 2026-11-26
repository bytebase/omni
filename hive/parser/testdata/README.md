# Hive parser test data

## `compat.tsv`: Apache Hive 4.0.1's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/Hive.g4`, then its near-misses: the statement with one token
dropped, the end cut off, a stray word inserted, or a keyword swapped. Each
row records omni's verdict (`hive/parser.Parse`) and the pinned Apache Hive
4.0.1 server's (`../../testdata/pins.tsv`, `apache/hive:4.0.1`).
`TestCompatibility` (`../compat_test.go`) reads the file offline through
`enginetest/compat` and starts no Docker container. Recorded against
`apache/hive:4.0.1` (`OMNI_RECORD=1`) and verified byte-identical on a second
fresh container (`OMNI_RECORD=verify`).

### Counts at `apache/hive:4.0.1` (`rel/release-4.0.1`)

| Rows | Total | Agree | Only Hive accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| Generated alternatives | 2,556 | 2,556 | 0 | 0 |
| Generated near-misses | 5,477 | 5,477 | 0 | 0 |
| **All** | **8,033** | **8,033** | **0** | **0** |

Apache Hive 4.0.1 accepts 2,535 of the 2,556 generated alternatives and rejects
21: 12 `MERGE` rows with zero `WHEN` clauses (where Hive 4.0.1's `HiveParser`
throws ANTLR 3's `RewriteEmptyStreamException` on `whenClauses` during parsing),
4 `SELECT a(*, ...)` rows (where `function`'s `(STAR) => STAR` syntactic
predicate matches bare `*` before `selectExpression`), and 5 `CREATE TABLE ...
LIKE FILE` rows followed by `PARTITIONED BY`, `STORED BY` or `TBLPROPERTIES`
(where HiveParser.g's k=3 decision reads the keyword as the file format and
then wants a URI string). It rejects 5,025 of the 5,477 near-misses and accepts
452 whose single-token mutation still produces valid HiveQL. `hive/parser`
agrees with live `apache/hive:4.0.1` on all 8,033 rows (100.0%).

## Grammar source

- **The grammar** is Apache Hive 4.0.1's official ANTLR 3 HiveQL grammar
  (`parser/src/java/org/apache/hadoop/hive/ql/parse/` at tag
  `rel/release-4.0.1`, commit `3af4517eb8cfd9407ad34ed78a0b48b57dfaa264`,
  `grammar` role in `../../testdata/pins.tsv`), vendored byte-for-byte in
  `../../testdata/vendor/*.g` under the Apache License 2.0 (`Hive.LICENSE`,
  `Hive.NOTICE`). `TestGrammarVendorDigest` holds all 10 files to their
  pinned SHA-256.
- **`Hive.g4` is converted from them** by `../cmd/antlr3to4`, since they are
  ANTLR 3. Each rule is preceded by the file and line it came from, with AST
  rewrite rules (`-> ...`), actions, arguments, labels, tree operators, and
  predicates dropped, and empty parent lexer stubs (`StringLiteral`,
  `QuotedIdentifier` in `HiveLexerParent.g`) superseded by `HiveLexer.g`.
  `TestHiveG4IsTheConversion` fails if `Hive.g4` drifts from the converter's
  output.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to `HiveLexer.g`
  and `HiveLexerParent.g`. Ten terminals referenced in the parser grammar have
  no lexer rule in `HiveLexerParent.g` (`KW_ANTI`, fragment `QuotedIdentifier`,
  and 8 stale `nonReserved` entries) and are marked `(unreachable)`.
  A rule that chooses among a few forms and is named at several sites, such
  as `tableName` or `principalName`, gets a row for each form at each site:
  `tokens.tsv` names it with the form `(every alternative)`.
- **25 alternatives have no generated row:** the 10 `(unreachable)` terminals
  above plus unused standalone/legacy rules in Apache Hive 4.0.1's `.g` files
  (`alterDatabaseSuffixSetManagedLocation`, `alterStatementSuffixStatsPart`,
  `disable`, `triggerExpressionStandalone`, `triggerOrExpression`,
  `triggerAndExpression`, `triggerActionExpressionStandalone`).

## How the live Hive 4.0.1 judge works

`compatJudge` (`../compat_test.go`) starts the pinned `apache/hive:4.0.1`
container via `hive/internal/testcontainer.Judge` (verifying both the running
container's image digest and `SELECT VERSION()`) only when `OMNI_RECORD=1` or
`OMNI_RECORD=verify` is set.

Every probed statement `<stmt>` is wrapped as:

```sql
EXPLAIN AST CREATE SCHEDULED QUERY omni_probe EVERY 1 DAY AS <stmt>
```

In `HiveParser.g`, `createScheduledQueryStatement -> definedAsSpec` directly
invokes the top-level start rule `statement : explainStatement EOF |
execStatement EOF`. Wrapping every statement this way:

1. Never executes `<stmt>` (not even DDL, DML, `USE`, `SET ROLE`,
   `KILL QUERY`, or `EXPLAIN ANALYZE`), leaving server state untouched.
2. Parses `<stmt>` through Hive 4.0.1's exact `HiveLexer` and
   `HiveParser.statement()` rule to `EOF` with no prefix ambiguity (avoiding
   bare `EXPLAIN <stmt>` collisions with `EXPLAIN ...`, `ANALYZE TABLE ...`,
   and near-misses beginning with `explainOption` keywords such as `EXTENDED`,
   `FORMATTED`, or `AST`).

`hiveVerdict` classifies HiveServer2's response fail-closed:

- **`err == nil`** -> **`accept`**.
- **Non-`*testcontainer.Error`** (transport, timeout, closed connection) ->
  **fail the recording**.
- **`FAILED: ParseException`**, **`FAILED: RewriteEmptyStreamException`**, or
  **`FAILED: RewriteCardinalityException`** (thrown directly by ANTLR 3
  inside `ParseDriver.parse`) -> **`reject`**.
- **Known post-parse Hive exceptions** (`SemanticException` and subclasses such
  as `CalciteSubquerySemanticException`, `IllegalStateException`,
  `Execution Error`, `HiveAccessControlException`, `InvalidOperationException`,
  `UnsupportedOperationException`, `IllegalArgumentException`,
  `NumberFormatException`, `NullPointerException`, `RuntimeException`,
  `HiveException`, etc.) -> **`accept`** (`ParseDriver.parse` succeeded before
  compilation/analysis failed).
- **Code 0 `"Error running query"`** (unhandled post-parse `Throwable` in
  HiveServer2, observed on `CREATE [MANAGED] TABLE a LIKE FILE`): verified with
  a tail probe (`stmt + "\n)"`). Only a probe `FAILED: ParseException` proves
  `stmt` parsed cleanly to `EOF` -> **`accept`**; any other probe outcome fails
  the recording.
- **Any other error** -> **fail the recording** (never guesses).

`TestHiveVerdict` (`../compat_test.go`) pins every branch of this classifier
offline.

## Commands

From `<repo root>`. Replay offline (no Docker):

```bash
go test ./hive/parser/ -run '^TestCompatibility$'
```

Regenerate `Hive.g4` after a vendor grammar update, then record and verify
against the pinned `apache/hive:4.0.1` image:

```bash
go run ./hive/parser/cmd/antlr3to4
OMNI_RECORD=1 go test -timeout 30m ./hive/parser/ -run '^TestCompatibility$'
OMNI_RECORD=verify go test -timeout 30m ./hive/parser/ -run '^TestCompatibility$'
```

Commit rows that a parser change moved:

```bash
go test ./hive/parser/ -run '^TestCompatibility$' -update-compat
```
