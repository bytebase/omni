# Oracle Parser Completion Scenarios

This file defines the parser-layer completion target for Oracle. It is intentionally scenario-only: each item describes observable behavior that must pass, not how to implement it.

## Completion Definition

Oracle parser completion means:

- Invalid SQL returns an error instead of a panic, partial AST, or silent success.
- Valid SQL already covered by corpus and parser tests continues to parse successfully.
- Every parsed AST node has a valid byte-span `Loc{Start, End}`, or the explicit unknown sentinel `Loc{-1, -1}` when a location is genuinely unavailable.
- Oracle keyword handling distinguishes reserved, nonreserved, context, function-like, pseudo-column, and quoted identifiers in every identifier context.
- Parser defense matrix status for Oracle can be updated from `--/Partial` to measured states with passing tests.

Current measured gate: `TestOracleParserProgress` reports 445/445 parser `parse*` methods with error returns, 0 silent parser error discards, 62 soft-fail scenarios, 121 strictness scenarios, 344 keyword golden scenarios, 171 classified BNF rows, 249 classified Loc-node rows, 152 direct Loc fixtures, and 20 optional reference-oracle rows. `TestVerifyCorpus` reports 128 statements with 0 parse violations, 0 Loc violations, and 0 crashes. The explicit Oracle Free reference run passes all 20 reference rows and checks 107 word-like reserved/context entries from `V$RESERVED_WORDS`.

## Phase 1: Error Propagation Contract

### 1.1 Parse Function Signatures

- [x] Every required `parse*` function returns an error-capable signature.
- [x] Optional parser probes are named distinctly from required parse functions.
- [x] Required parse functions never signal failure by returning a bare nil node.
- [x] Statement dispatch returns a node and an error for every statement family.
- [x] Expression parsing returns an expression and an error for every required expression.
- [x] DDL parsing returns a statement or command node and an error for every required grammar branch.
- [x] PL/SQL parsing returns an error for malformed block structure.

### 1.2 Error Propagation

- [x] A syntax error raised in a nested expression is returned by `Parse`.
- [x] A syntax error raised in a nested SELECT clause is returned by `Parse`.
- [x] A syntax error raised in a nested DDL option is returned by `Parse`.
- [x] A syntax error raised inside PL/SQL is returned by `Parse`.
- [x] Lexer errors take precedence over EOF syntax errors.
- [x] No production parser call discards an error return with `_`.
- [x] No production parser call replaces a nested error with silent success.

### 1.3 Error Quality

- [x] EOF errors say `syntax error at end of input`.
- [x] Token errors say `syntax error at or near "<token text>"`.
- [x] Error positions are byte offsets inside the input or at EOF.
- [x] Unterminated quoted identifiers report lexer context.
- [x] Unterminated string literals report lexer context.
- [x] Unterminated q-quote literals report lexer context.
- [x] Unterminated block comments report lexer context.

## Phase 2: Loc Sentinel And Span Completeness

### 2.1 Sentinel Semantics

- [x] `Loc{Start: 0, End: 0}` is treated as a real byte span only when it is actually produced by a zero-length construct.
- [x] Unknown locations use `Loc{Start: -1, End: -1}`.
- [x] `Loc.End == 0` is not used as an unknown sentinel.
- [x] Public helper docs describe `-1` as the only unknown sentinel.
- [x] Loc verifier rejects mixed sentinel values such as `{Start: -1, End: 0}`.

### 2.2 Token And Raw Statement Spans

- [x] Every non-EOF token has `End >= Loc`.
- [x] String literal token spans include both quotes in source byte positions.
- [x] Quoted identifier token spans include both quotes in source byte positions.
- [x] Hint token spans include the full hint comment.
- [x] RawStmt `Loc.End` points to the last consumed statement token, not the next token.
- [x] Multi-statement parsing records each RawStmt span independently.

### 2.3 AST Span Coverage

- [x] SELECT target nodes have nonzero positive spans.
- [x] FROM table refs have nonzero positive spans.
- [x] JOIN nodes include the join keyword through the join condition.
- [x] WHERE/HAVING predicates have spans within the SELECT span.
- [x] Function call names have spans.
- [x] Function call args have spans.
- [x] Alias nodes have spans.
- [x] CREATE TABLE column definitions have spans.
- [x] CREATE TABLE constraints have spans.
- [x] ALTER TABLE commands have spans.
- [x] PL/SQL declarations have spans.
- [x] PL/SQL statements have spans.
- [x] Corpus Loc verifier fails the test on any invalid Loc.

## Phase 3: Soft-Fail Matrix

### 3.1 Expression Truncation

- [x] `SELECT 1 +` errors at EOF.
- [x] `SELECT 1 -` errors at EOF when parsed as infix.
- [x] `SELECT 1 *` errors at EOF.
- [x] `SELECT 1 /` errors at EOF.
- [x] `SELECT NOT` errors at EOF.
- [x] `SELECT -` errors at EOF.
- [x] `SELECT PRIOR` errors at EOF.
- [x] `SELECT CONNECT_BY_ROOT` errors at EOF.
- [x] `SELECT CASE WHEN 1 THEN` errors at EOF.
- [x] `SELECT CAST(1 AS` errors at EOF.
- [x] `SELECT DECODE(1,` errors at EOF.
- [x] `SELECT XMLTABLE(` errors at EOF.
- [x] `SELECT JSON_TABLE(` errors at EOF.

### 3.2 Predicate Truncation

- [x] `SELECT 1 BETWEEN` errors at EOF.
- [x] `SELECT 1 BETWEEN 0 AND` errors at EOF.
- [x] `SELECT 1 IN (` errors at EOF.
- [x] `SELECT 1 IN (1,` errors at EOF.
- [x] `SELECT 'a' LIKE` errors at EOF.
- [x] `SELECT 'a' LIKE 'b' ESCAPE` errors at EOF.
- [x] `SELECT 1 IS` errors at EOF.
- [x] `SELECT 1 IS NOT` errors at EOF.

### 3.3 Statement Truncation

- [x] `SELECT` errors at EOF.
- [x] `SELECT * FROM` errors at EOF.
- [x] `SELECT * FROM t WHERE` errors at EOF.
- [x] `SELECT * FROM t JOIN` errors at EOF.
- [x] `SELECT * FROM t JOIN t2 ON` errors at EOF.
- [x] `SELECT 1 GROUP BY` errors at EOF.
- [x] `SELECT 1 ORDER BY` errors at EOF.
- [x] `SELECT 1 UNION` errors at EOF.
- [x] `INSERT INTO t (` errors at EOF.
- [x] `INSERT INTO t VALUES (` errors at EOF.
- [x] `UPDATE t SET` errors at EOF.
- [x] `DELETE FROM` errors at EOF.
- [x] `MERGE INTO t USING` errors at EOF.
- [x] `CREATE TABLE t (` errors at EOF.
- [x] `CREATE TABLE t (a NUMBER DEFAULT` errors at EOF.
- [x] `ALTER TABLE t ADD` errors at EOF.
- [x] `CREATE PROCEDURE p IS BEGIN` errors at EOF.
- [x] `DECLARE x NUMBER; BEGIN` errors at EOF.

## Phase 4: Strictness Matrix

### 4.1 Duplicate Clauses

- [x] Duplicate SELECT `WHERE` is rejected.
- [x] Duplicate SELECT `GROUP BY` is rejected.
- [x] Duplicate SELECT `HAVING` is rejected.
- [x] Duplicate SELECT `ORDER BY` is rejected.
- [x] Duplicate SELECT `FETCH` is rejected.
- [x] Duplicate CREATE TABLE table options with exclusive semantics are rejected where Oracle rejects them.
- [x] Oracle-supported ALTER TABLE multi-action forms remain accepted while unknown ALTER TABLE actions are rejected.

### 4.2 Unknown Options

- [x] Unknown CREATE TABLE option is rejected.
- [x] Unknown ALTER TABLE option is rejected.
- [x] Unknown CREATE INDEX option is rejected.
- [x] Unknown DROP option is rejected.
- [x] Unknown GRANT option is rejected.
- [x] Unknown ALTER SESSION option is rejected.
- [x] Unknown CREATE USER option is rejected.

### 4.3 Illegal Keyword Position

- [x] `SELECT FROM t` is rejected.
- [x] `SELECT * FROM WHERE` is rejected.
- [x] `CREATE TABLE SELECT (a NUMBER)` is rejected.
- [x] `CREATE TABLE t (SELECT NUMBER)` is rejected.
- [x] `SELECT 1 FROM SELECT` is rejected.
- [x] `ALTER TABLE t ADD SELECT NUMBER` is rejected.
- [x] `CREATE INDEX SELECT ON t(a)` is rejected.

### 4.4 Parentheses And Statement Boundaries

- [x] Unclosed parenthesized expression is rejected.
- [x] Unclosed subquery is rejected.
- [x] Unclosed column definition list is rejected.
- [x] Unclosed PL/SQL block is rejected.
- [x] Trailing garbage after a valid statement is rejected.
- [x] Two statements without a semicolon separator are rejected.

### 4.5 Strictness V2 Coverage

- [x] Table-driven strictness manifest covers at least 100 negative parser scenarios.
- [x] SELECT row limiting requires valid `OFFSET` and `FETCH` operands.
- [x] JOIN `USING` requires parentheses and at least one identifier.
- [x] CASE and EXISTS reject missing operands and malformed subqueries.
- [x] MERGE rejects missing `ON`, missing `WHEN`, missing action, and malformed INSERT/UPDATE branches.
- [x] CREATE SEQUENCE, SYNONYM, INDEX, COMMENT, TRUNCATE, REVOKE, and ALTER SESSION reject missing mandatory operands.
- [x] PL/SQL blocks reject missing `BEGIN`, missing `END`, and malformed loop endings while still accepting declarations before `BEGIN`.

## Phase 5: Oracle Keyword Matrix

### 5.1 Classification Completeness

- [x] Every lexer keyword token has an explicit Oracle category and manifest row.
- [x] Reserved keyword golden list matches the documented local Oracle reserved-identifier policy.
- [x] Nonreserved keyword golden list matches documented local policy.
- [x] Context keyword list has exhaustive manifest coverage.
- [x] Function-like keyword list has exhaustive manifest coverage.
- [x] Pseudo-column list has exhaustive manifest coverage.

### 5.2 Identifier Contexts

- [x] Unquoted reserved table names are rejected.
- [x] Quoted reserved table names are accepted.
- [x] Unquoted nonreserved table names are accepted where Oracle accepts them.
- [x] Unquoted reserved column names are rejected.
- [x] Quoted reserved column names are accepted.
- [x] Alias contexts accept legal Oracle aliases.
- [x] Alias contexts do not consume clause-start keywords.
- [x] Dotted names allow reserved words only when quoted.
- [x] DB link names follow the same identifier rules as object names.

### 5.3 Keyword-As-Function And Pseudo-Columns

- [x] Built-in function keywords parse as functions in call position.
- [x] Function-like keywords parse as identifiers in non-call identifier contexts where allowed.
- [x] Pseudo-columns parse as pseudo-column expressions in expression position.
- [x] Quoted pseudo-column names parse as identifiers, not pseudo-columns.
- [x] Keywords after `.` are parsed according to qualified-name rules, not global clause rules.
- [x] Oracle 26ai SQL reserved words are pinned in a separate official-source manifest.
- [x] Every pinned Oracle 26ai SQL reserved word appears in the local lexer keyword manifest.
- [x] Pinned Oracle SQL reserved words are not lexed as plain identifiers.

## Phase 6: Coverage And Defense Matrix

### 6.1 Corpus

- [x] Oracle corpus positive cases parse successfully.
- [x] Oracle corpus negative cases reject with errors.
- [x] Oracle corpus parser mismatch count is zero.
- [x] Oracle corpus Loc violation count is zero and fails the test when nonzero.
- [x] Oracle corpus crash count is zero.
- [x] Error position validation is enabled for failed corpus cases.

### 6.2 BNF Coverage Accounting

- [x] Every `oracle/parser/bnf/*.bnf` file appears in `testdata/coverage/bnf_coverage.tsv`.
- [x] Every BNF row has a classified status; `unknown` is rejected by test.
- [x] High-value parser families have no `missing` or `unknown` rows.
- [x] Coverage statuses distinguish parser coverage, partial parser debt, catalog-only semantics, unsupported grammar, and deferred work.
- [x] Every non-covered BNF row has explicit debt class, approval, and next action metadata.
- [x] Every Oracle BNF row appears in the P2 manifest with parser-layer status, AST target, parser entrypoint, next action, and positive/negative/Loc evidence fields.
- [x] Historical `catalog` status is not reused as a P2 parser-layer status.
- [x] P2 skip/stub hotspots are tracked in an inventory keyed by stable `file:function:pattern` identities.
- [x] P2 skip/stub inventory rows must link back to P2 BNF manifest rows, and `p2_done` rows cannot retain owned skip/stub sites.
- [x] `TestOracleCoverage` enforces soft-fail, strictness, keyword, BNF, Loc-node, and reference-oracle minimum gates.

### 6.3 Loc Node And Reference Oracle Accounting

- [x] Every Loc-bearing AST node in `oracle/ast/loc.go` appears in `testdata/coverage/loc_node_coverage.tsv`.
- [x] Loc-node manifest rows have classified status; `unknown` is rejected by test.
- [x] Covered Loc-node fixtures must parse, contain the expected node type, and pass `CheckLocations`.
- [x] Direct Loc fixture coverage has a floor of 80 rows and currently covers 152 rows.
- [x] Deferred Loc-node rows have explicit debt class, approval, and next action metadata.
- [x] Optional reference-oracle manifest has at least 20 rows across parser families.
- [x] Reference-oracle execution is build-tag gated and skipped unless `ORACLE_PARSER_REF_DSN` or `ORACLE_PARSER_REF_CONTAINER=1` is set.
- [x] Reference-oracle execution can also use an explicit Oracle Free testcontainer with `ORACLE_PARSER_REF_CONTAINER=1`.
- [x] Strict reference mode fails when no real Oracle backend is provided.
- [x] Oracle Free testcontainer reference proof passes all 20 manifest rows.
- [x] Oracle Free `V$RESERVED_WORDS` proof checks 107 word-like reserved/context entries against the local manifest.

### 6.4 Documentation And Status

- [x] `docs/PARSER-DEFENSE-MATRIX.md` reports measured Oracle status for L1, L2, L3, L4, L5, and L7.
- [x] `docs/engine-capability-guide.md` Oracle status reflects parser-layer progress.
- [x] Oracle quality strategy no longer lists resolved foundation blind spots.
- [x] A developer can run one command to print dual-return, Loc, keyword, strictness, and corpus status.
- [x] Completion implementation remains explicitly out of scope for parser-coverage completion and has a separate scope plan.
