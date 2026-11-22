# Databricks parser test data

## `compat.tsv`: Apache Spark's SQL grammar rows

`compat.tsv` holds a statement for each reachable alternative of each rule in
`../../testdata/SqlBaseParser.g4` (paired with `SqlBaseLexer.g4`), followed by
its near-misses (one token dropped, the end cut off, a stray word inserted, or a
keyword swapped). Each row records omni's verdict and the grammar recognizer's
verdict. `TestCompatibility` (`../compat_test.go`) reads it through
`enginetest/compat`, regenerates the rows from the vendor grammar, and asks the
recognizer again on every run with no server, credentials, or network.

| Rows | Total | Agree | Only the grammar accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 2,229 | 2,220 | 9 | 0 |
| Near-misses | 4,209 | 4,177 | 29 | 3 |
| **All** | **6,438** | **6,397** | **38** | **3** |

- **The grammar** is `../../testdata/SqlBaseParser.g4` with its lexer grammar
  `../../testdata/SqlBaseLexer.g4`: Apache Spark's official ANTLR 4 SQL grammar
  (`sql/api/src/main/antlr4/org/apache/spark/sql/catalyst/parser/` in
  [`apache/spark`](https://github.com/apache/spark)), which Databricks SQL is
  built on, pinned at release `v4.0.5` (`../../testdata/pins.tsv`, role
  `default`, matching Databricks SQL 4.2.0 / Spark 4.0). Both files are
  Apache-2.0 (`../../testdata/Spark.LICENSE`) and committed byte for byte; their
  concatenated SHA-256 is the `default` pin's digest.
- **Role `driver` in `../../testdata/pins.tsv` is untouched** and continues to
  pin the Databricks workspace for `databricks/driver` tapes.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the lexer rules
  of `SqlBaseLexer.g4`. One rule override (`primaryExpression` -> `a`) keeps
  generated rows readable because `SqlBaseParser.g4` lists `CURRENT_DATE` before
  identifiers in `primaryExpression`. `.` is sampled as `a` for the `.*?`
  wildcard productions in `REFRESH`, `ADD`/`LIST`, `SET`/`RESET`, and
  `unsupportedHiveNativeCommands`.
- **22 alternatives have no generated row:** Spark's 8 auxiliary non-statement
  entry points (`singleExpression`, `singleTableIdentifier`,
  `singleMultipartIdentifier`, `singleFunctionIdentifier`, `singleDataType`,
  `singleTableSchema`, `singleRoutineParamList`, `singleQuery`) and 4 unused
  helper rules (`executeImmediateArgument`, `executeImmediateArgumentSeq`,
  `functionIdentifier`, `qualifiedNameList`), none of which are reachable from
  the start rule `compoundOrSingleStatement`.

To judge new rows, or every row after a pin bump:

```bash
OMNI_RECORD=1 go test ./databricks/parser/ -run '^TestCompatibility$'
```

To verify committed answers read-only:

```bash
OMNI_RECORD=verify go test ./databricks/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./databricks/parser/ -run '^TestCompatibility$' -update-compat
```

## Who judges what

Databricks SQL is a cloud service with no container CI can start, so the
committed baseline is judged offline by `enginetest/grammar`'s recognizer over
`SqlBaseParser.g4` + `SqlBaseLexer.g4` (`SourceGrammar`). When `DATABRICKS_*`
warehouse credentials are available, this suite can be upgraded to a live
Databricks SQL warehouse judge (`SourceEngine`) without changing the generated
grammar rows.

- **Grammar judge (`SourceGrammar`), 6,438 rows.** A recognizer over
  `SqlBaseParser.g4` + `SqlBaseLexer.g4` accepts any statement the grammar
  derives and rejects every other row.
- **Agreement: 6,397 / 6,438 (99.36%), with 2,220 / 2,229 (99.60%) on
  generated grammar alternatives (the only 9 excluded alternatives are Spark's
  error-capturing rules for unquoted hyphens, invalid pipe-operator forms, and
  mixed/repeated interval continuations), 38 `judge-only` rows, and 3
  `omni-only` near-miss rows.**

## Remaining tail (38 `judge-only`, 3 `omni-only` rows)

All 38 `judge-only` discrepancies are invalid SQL forms that `databricks/parser`
intentionally rejects while the grammar recognizer accepts them due to known
properties of the offline recognizer and upstream error-capturing rules, and
the 3 `omni-only` near-misses are `REFRESH .*?` statements accepted by Spark's
ANTLR grammar that the offline recognizer's `.` wildcard does not match:

1. **9 malformed optimizer hint near-misses** (`SELECT /*+ */ a`,
   `SELECT /*+(a) */ a`, `SELECT /*+ a a) */ a`, `SELECT /*+ a() */ a`,
   `SELECT /*+ a(a */ a`, `SELECT /*+ a(a zzz) */ a`, `SELECT /*+ a(a a) */ a`,
   `SELECT /*+ a(a,) */ a`, `SELECT /*+ a(a, a zzz) */ a`):
   `enginetest/grammar`'s generic SQL lexer skips `/*+ ... */` as an ordinary
   block comment (`enginetest/grammar/recognize.go`), so the recognizer sees
   `SELECT a` and accepts all 9. `databricks/parser` properly lexes `/*+` / `*/`
   and enforces `hintStatement` syntax, rejecting all 9.
2. **9 `.*?` wildcard near-misses** (`SHOW CREATE TABLE a zzz`,
   `DROP INDEX a a`, ``SET `a` `a` ``, ``SET `a` = `a` zzz``, ``SET `a` zzz``,
   ``RESET `a` zzz``, `CREATE INDEX zzz`, `DROP INDEX zzz`,
   `ALTER TABLE a REPLACE COLUMNS zzz`): `SqlBaseParser.g4` contains fallback
   rules ending in `.*?` (`unsupportedHiveNativeCommands .*?`, `SET .*?`,
   `RESET .*?`) whose wildcard matches trailing `zzz` / `a` / `` `a` `` in the
   recognizer, whereas `databricks/parser` parses `SHOW CREATE TABLE`,
   `CREATE INDEX`, `DROP INDEX`, `ALTER TABLE ... REPLACE COLUMNS`, `SET`, and
   `RESET` strictly and rejects these malformed near-misses.
3. **9 error-capturing interval continuation rows** (`errorCapturingMultiUnitsInterval`
   and `errorCapturingUnitToUnitInterval`: 3 alternatives and 6 near-misses for
   mixed or repeated interval continuations like `INTERVAL 1 NANOSECOND 1 SECOND TO SECOND`
   and `INTERVAL 1 SECOND TO SECOND 1 NANOSECOND`): Spark v4.0.5's
   `AstBuilder.parseIntervalLiteral` explicitly rejects these continuations as
   syntax errors; strict `databricks/parser` rejects them directly.
4. **7 unquoted hyphenated identifier (`errorCapturingIdentifierExtra`) rows**
   (3 alternatives `FROM a - a`, `FROM a - a - a` and 4 near-misses including
   `DESC FUNCTION - zzz`): `SqlBaseParser.g4:1479` parses `(MINUS identifier)+`
   solely so Spark's `AstBuilder` can reject unquoted hyphenated names like
   `my-table` with a tailored diagnostic requiring backticks. `databricks/parser`
   rejects these directly with a positioned `ParseError` while accepting
   backtick-quoted identifiers.
5. **4 invalid SQL pipe-operator (`operatorPipeRightSide` /
   `operatorPipeSetAssignmentSeq`) rows** (3 alternatives and 1 near-miss for
   combined `|> PIVOT ... UNPIVOT ...`, `|> UNPIVOT ... PIVOT ...`, and dotted
   `|> SET a.a = a`): `SqlBaseParser.g4:1561-1584` explicitly comments that
   these productions (along with `|> WHERE ... WINDOW ...`) are captured solely
   for diagnostic error messages; strict `databricks/parser` rejects them.
6. **3 `omni-only` `REFRESH .*?` wildcard near-misses**
   (`BEGIN REFRESH REFRESH; END`,
   `BEGIN CASE WHEN a THEN REFRESH WHEN a THEN REFRESH; END CASE; END`,
   `BEGIN CASE a WHEN a THEN REFRESH WHEN a THEN REFRESH; END CASE; END`):
   `SqlBaseParser.g4:325` defines `REFRESH (stringLit | .*?)`, so ANTLR and
   `databricks/parser` consume unquoted resource paths (such as
   `REFRESH /mnt/data`) through the statement boundary (`;`), whereas
   `enginetest/grammar`'s offline recognizer binds `.` only to the identifier
   token class in `tokens.tsv` and therefore rejects keyword tokens (`REFRESH`,
   `WHEN`) inside the `.*?` tail.

## What the judge cannot decide

- **Optimizer hints (`/*+ ... */`).** As noted above, the recognizer's generic
  SQL lexer skips `/*+ ... */` as a block comment; valid and malformed hints are
  tested directly in `../parser_test.go`.
- **Semantic predicates and visitor checks.** `enginetest/grammar`'s ANTLR
  reader drops semantic predicates (`SQL_standard_keyword_behavior`,
  `double_quoted_identifiers`, `legacy_setops_precedence_enabled`,
  `isShiftRightOperator()`, `isValidDecimal()`), and Spark's `AstBuilder`
  enforces additional post-parse checks (e.g. `errorCapturingIdentifier` rejects
  unquoted hyphenated names like `a-b`, `operatorPipeRightSide` rejects `WINDOW`
  after `WHERE` or combined `PIVOT`/`UNPIVOT`).
- **Databricks-only SQL extensions beyond Spark 4.0.** Databricks SQL adds
  clauses such as `QUALIFY` (supported in `databricks/parser` and covered in
  `../parser_test.go`) and Unity Catalog DDL extensions (`CREATE VOLUME`,
  `CREATE CATALOG`, `GRANT`/`REVOKE` privileges) that are not in upstream
  `apache/spark`'s `SqlBaseParser.g4`.
- **Lexing nuances.** The recognizer uses a generic SQL lexer rather than
  `SqlBaseLexer.g4` (e.g. raw strings `R'...'` / `R"..."`, typed numeric
  suffixes `1L`/`1S`/`1Y`/`1F`/`1D`/`1BD`, nested `/* /* */ */` block comments,
  and `--` `\\\n` line continuation are covered in `../parser_test.go`).
