# Oracle Parser Completion Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Finish the Oracle parser-layer quality gap against MySQL/PG/MSSQL by completing error propagation, Loc semantics, soft-fail/strictness, keyword classification, and defense-matrix verification.

**Architecture:** Move from ad hoc parser fixes to contract-driven parser infrastructure. First establish scenario and guard tests that define completion, then migrate parser internals to dual-return error propagation, make Loc validation a hard gate, replace broad keyword-as-identifier behavior with contextual identifier APIs, and finally expand strict/soft-fail coverage until Oracle has measured defense-matrix status.

**Tech Stack:** Go recursive-descent parser in `oracle/parser`, AST nodes in `oracle/ast`, corpus verifier in `oracle/parser/verify_parse_test.go`, strategy docs in `docs/` and `oracle/quality/`.

## Baseline

Initial baseline when this plan was written:

- `rg -n '^func \(p \*Parser\) parse' oracle/parser --glob '*.go' | wc -l` reports `442`.
- `rg -n '^func \(p \*Parser\) parse.* error' oracle/parser --glob '*.go' | wc -l` reports `1`.
- Current foundation slice has `parseStmt() (nodes.StmtNode, error)`, parser-level syntax helpers, lexer error propagation, partial soft-fail checks, `p.prev.End` Loc cleanup, and initial keyword classification.
- `go test -count=1 ./oracle/parser` passes.
- `go test ./...` is not a required parser-layer gate because existing MySQL/TiDB catalog container tests fail independently.

Current measured status after execution:

- `TestOracleParserProgress` reports 445/445 parser `parse*` methods return `error`, with 0 silent parser error discards.
- `TestSoftFail*` covers 62 truncation and missing-operand scenarios.
- `TestStrict*` plus `TestStrictV2CoverageMatrix` cover 121 duplicate-clause, unknown-option, illegal-keyword-position, statement-boundary, reserved-identifier, parenthesis-balance, DML/DDL utility, and PL/SQL strictness scenarios.
- `TestOracleKeyword*` covers an exhaustive 344-row manifest for the declared Oracle lexer keyword set, plus a pinned Oracle 26ai SQL reserved-word audit.
- `TestOracleCoverage` enforces 171 classified BNF rows, 0 unknown BNF rows, 249 classified Loc-node rows, 152 direct Loc fixtures, 0 unknown Loc-node rows, and 20 optional reference-oracle rows.
- `TestVerifyCorpus` reports 128 statements, 125 parser accepts, 3 expected parser rejects, 0 parse violations, 0 Loc violations, and 0 crashes.

## Global Rules

- Do not change the public `Parse(sql string) (*nodes.List, error)` API.
- Do not introduce a parser generator or external grammar runtime.
- Do not hide parser errors in transitional fields once a function is migrated.
- Do not use `End == 0` as unknown; use `Loc{Start: -1, End: -1}` only.
- Do not make `parseIdentifier()` globally stricter until context-specific identifier APIs exist.
- Do not run `gofmt` over the whole repository; use targeted `gofmt -w` on touched Go files.
- Every task must add or tighten a failing test first, then implement, then run the listed proof command.

## Canonical Proof Commands

Use these commands throughout:

```bash
go test -run 'TestFoundation|TestVerifyCorpus|TestOracleKeyword|TestSoftFail|TestStrict|TestOracleCoverage' -count=1 -v ./oracle/parser
go test -count=1 ./oracle/parser
git diff --check
```

For global awareness only, not as a blocker for Oracle parser completion:

```bash
go test ./...
```

If `go test ./...` fails in non-Oracle packages, record the package and first failure, then continue only if `./oracle/parser` is green.

## Phase 0: Coverage Contract

### Task 0.1: Land Completion Scenarios

**Files:**
- Create: `oracle/parser/SCENARIOS-oracle-parser-completion.md`

**Step 1: Write the scenario file**

Create a scenario-only document with phases for error propagation, Loc, soft-fail, strictness, keywords, corpus, and docs status.

**Step 2: Review the scenario scope**

Run:

```bash
sed -n '1,260p' oracle/parser/SCENARIOS-oracle-parser-completion.md
```

Expected: the file contains checkboxes and no implementation details such as file paths in scenario items.

**Step 3: Commit**

```bash
git add oracle/parser/SCENARIOS-oracle-parser-completion.md
git commit -m "docs: define oracle parser completion scenarios"
```

### Task 0.2: Add Contract Guard Tests

**Files:**
- Create: `oracle/parser/parser_contract_test.go`
- Modify: `oracle/parser/verify_parse_test.go`

**Step 1: Write failing signature/count tests**

Add tests that:

- Count parser methods whose names match `func (p *Parser) parse`.
- Count parser methods whose signature does not include an `error` return.
- Allow only a small explicit transitional allowlist.
- Fail on production parser error discards such as `_, _ := p.parse`.

Suggested command inside the test may use Go source parsing rather than `rg`, so it works cross-platform.

**Step 2: Run the guard and confirm RED**

Run:

```bash
go test -run 'TestOracleParserContract' -count=1 ./oracle/parser
```

Expected: FAIL with current bare parser methods listed.

**Step 3: Add Loc hard-gate test mode**

Change corpus verification so any Loc violation fails the test, not just logs.

**Step 4: Run and confirm current state**

Run:

```bash
go test -run 'TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected after current foundation slice: PASS with `LOC VIOLATIONS: 0`.

**Step 5: Commit**

```bash
git add oracle/parser/parser_contract_test.go oracle/parser/verify_parse_test.go
git commit -m "test: add oracle parser completion guards"
```

## Phase 1: Dual-Return Error Propagation

### Task 1.1: Define Parser Error Conventions

**Files:**
- Modify: `oracle/parser/parser.go`
- Modify: `oracle/parser/README.md` if present, otherwise document in `oracle/parser/SKILL.md`

**Step 1: Add conventions to parser docs**

Document these rules:

- Required parse functions return `(T, error)`.
- Optional probe helpers use `tryParse*` or `parseOptional*` naming and may return `(nil, nil)`.
- Required parse functions return syntax errors, never placeholder nodes, for missing required constituents.
- Nested errors propagate directly to `Parse()`.

**Step 2: Run docs/contract tests**

Run:

```bash
go test -run 'TestOracleParserContract' -count=1 ./oracle/parser
```

Expected: still FAIL until migrations happen, but error text should reflect the intended allowlist.

### Task 1.2: Migrate Identifier And Type Helpers

**Files:**
- Modify: `oracle/parser/name.go`
- Modify: `oracle/parser/type.go`
- Modify: parser call sites in files that consume identifier/type helpers

**Step 1: Write failing tests**

Add tests for:

- Required object name missing after `CREATE TABLE`.
- Required column name missing inside `CREATE TABLE`.
- Required datatype missing after `CAST(... AS`.
- Missing identifier after schema dot.

**Step 2: Change helper signatures**

Target signatures:

```go
func (p *Parser) parseIdentifier() (string, error)
func (p *Parser) parseAlias() (*nodes.Alias, error)
func (p *Parser) parseObjectName() (*nodes.ObjectName, error)
func (p *Parser) parseColumnRef() (*nodes.ColumnRef, error)
func (p *Parser) parseTypeName() (*nodes.TypeName, error)
```

Optional or permissive callers should be renamed or wrapped explicitly.

**Step 3: Update call sites**

Propagate returned errors immediately. Do not use `_` for errors.

**Step 4: Run focused proof**

```bash
go test -run 'TestSoftFail|TestStrict|TestOracleKeywordClassificationGolden' -count=1 ./oracle/parser
```

Expected: PASS.

**Step 5: Commit**

```bash
git add oracle/parser/name.go oracle/parser/type.go oracle/parser/*_test.go
git commit -m "refactor: propagate oracle identifier parse errors"
```

### Task 1.3: Migrate Expression Parser

**Files:**
- Modify: `oracle/parser/expr.go`
- Modify: expression call sites in `oracle/parser/select.go`, `insert.go`, `update_delete.go`, `merge.go`, `create_table.go`, `plsql_block.go`
- Test: `oracle/parser/soft_fail_strict_test.go`

**Step 1: Expand expression truncation tests**

Add RED tests for `CASE`, `CAST`, `DECODE`, `XMLTABLE`, `JSON_TABLE`, parenthesized expression, and subquery truncation.

**Step 2: Change expression signatures**

Target signatures:

```go
func (p *Parser) parseExpr() (nodes.ExprNode, error)
func (p *Parser) parseExprPrec(minPrec int) (nodes.ExprNode, error)
func (p *Parser) parsePrefix() (nodes.ExprNode, error)
func (p *Parser) parsePrimary() (nodes.ExprNode, error)
```

All postfix/infix helpers should return `(nodes.ExprNode, error)`.

**Step 3: Remove expression use of `Parser.err`**

Expression errors must return directly.

**Step 4: Run focused proof**

```bash
go test -run 'TestSoftFailTruncatedExpressions|TestSoftFailTruncatedPredicates|TestParseOraclePlusJoin' -count=1 ./oracle/parser
```

Expected: PASS.

**Step 5: Commit**

```bash
git add oracle/parser/expr.go oracle/parser/*_test.go
git commit -m "refactor: propagate oracle expression parse errors"
```

### Task 1.4: Migrate SELECT Parser

**Files:**
- Modify: `oracle/parser/select.go`
- Test: `oracle/parser/soft_fail_strict_test.go`

**Step 1: Add SELECT matrix tests**

Cover truncated SELECT list, FROM, JOIN, WHERE, GROUP BY, HAVING, ORDER BY, FETCH, set operators, subquery refs, lateral refs, XMLTABLE/JSON_TABLE table refs.

**Step 2: Change SELECT signatures**

Target signatures include:

```go
func (p *Parser) parseSelectStmt() (*nodes.SelectStmt, error)
func (p *Parser) parseSelectList() (*nodes.List, error)
func (p *Parser) parseFromClause() (*nodes.List, error)
func (p *Parser) parseTableRef() (nodes.TableExpr, error)
func (p *Parser) parseJoinContinuation(left nodes.TableExpr) (nodes.TableExpr, error)
```

**Step 3: Preserve optional grammar explicitly**

Use optional helpers for truly optional clauses. Required clauses after consumed keywords must error when missing.

**Step 4: Run focused proof**

```bash
go test -run 'TestSoftFailTruncatedClauses|TestStrictDuplicateClauses|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS with `PARSE VIOLATIONS: 0`, `LOC VIOLATIONS: 0`, `CRASHES: 0`.

**Step 5: Commit**

```bash
git add oracle/parser/select.go oracle/parser/*_test.go
git commit -m "refactor: propagate oracle select parse errors"
```

### Task 1.5: Migrate DML Parsers

**Files:**
- Modify: `oracle/parser/insert.go`
- Modify: `oracle/parser/update_delete.go`
- Modify: `oracle/parser/merge.go`
- Test: `oracle/parser/soft_fail_strict_test.go`

**Step 1: Add DML truncation tests**

Cover INSERT, UPDATE, DELETE, MERGE required constituents and subqueries.

**Step 2: Change DML signatures**

Target each statement parser to return `(*nodes.XStmt, error)` or `(nodes.StmtNode, error)` where dispatch requires an interface.

**Step 3: Run focused proof**

```bash
go test -run 'TestSoftFailDML|TestParseInsert|TestParseUpdate|TestParseDelete|TestParseMerge' -count=1 ./oracle/parser
```

Expected: PASS.

**Step 4: Commit**

```bash
git add oracle/parser/insert.go oracle/parser/update_delete.go oracle/parser/merge.go oracle/parser/*_test.go
git commit -m "refactor: propagate oracle dml parse errors"
```

### Task 1.6: Migrate DDL Core Parsers

**Files:**
- Modify: `oracle/parser/create_table.go`
- Modify: `oracle/parser/create_index.go`
- Modify: `oracle/parser/create_view.go`
- Modify: `oracle/parser/drop.go`
- Modify: `oracle/parser/alter_table.go`
- Test: `oracle/parser/soft_fail_strict_test.go`

**Step 1: Add DDL truncation and unknown-option tests**

Cover CREATE TABLE, CREATE INDEX, CREATE VIEW, DROP, and ALTER TABLE required grammar.

**Step 2: Change DDL signatures**

Every DDL parser returns `(nodes.StmtNode, error)` or concrete node plus error.

**Step 3: Make unknown option rejection explicit**

Known option parsers should validate option names rather than skipping arbitrary identifier-like tokens.

**Step 4: Run focused proof**

```bash
go test -run 'TestStrictUnknownOptions|TestSoftFailTruncatedClauses|TestParseCreate|TestParseAlterTable|TestParseDrop' -count=1 ./oracle/parser
```

Expected: PASS.

**Step 5: Commit**

```bash
git add oracle/parser/create_table.go oracle/parser/create_index.go oracle/parser/create_view.go oracle/parser/drop.go oracle/parser/alter_table.go oracle/parser/*_test.go
git commit -m "refactor: propagate oracle ddl parse errors"
```

### Task 1.7: Migrate Admin, Security, Session, Database, PL/SQL

**Files:**
- Modify: `oracle/parser/create_admin.go`
- Modify: `oracle/parser/security.go`
- Modify: `oracle/parser/session.go`
- Modify: `oracle/parser/database.go`
- Modify: `oracle/parser/plsql_block.go`
- Modify: remaining `oracle/parser/*.go`

**Step 1: Add truncation tests for each family**

Cover CREATE/ALTER USER, GRANT/REVOKE, ALTER SESSION, CREATE/ALTER DATABASE, anonymous blocks, procedure/function/package/trigger parsing.

**Step 2: Change remaining signatures**

Use the contract test to drive the allowlist to zero or to only explicit optional probe helpers.

**Step 3: Remove transitional parser error accumulator**

Delete `Parser.err` after all migrated functions return errors directly.

**Step 4: Run full parser proof**

```bash
go test -count=1 ./oracle/parser
```

Expected: PASS.

**Step 5: Commit**

```bash
git add oracle/parser oracle/ast
git commit -m "refactor: complete oracle parser error propagation"
```

## Phase 2: Loc Hard Gate

### Task 2.1: Centralize Loc Sentinel Helpers

**Files:**
- Modify: `oracle/ast/node.go`
- Modify: `oracle/ast/loc.go`
- Test: `oracle/parser/foundation_test.go` or new `oracle/parser/loc_contract_test.go`

**Step 1: Add sentinel tests**

Assert:

- `NoLoc()` returns `Loc{-1, -1}`.
- `Loc{0, 0}` is not treated as unknown by docs or helper code.
- Mixed sentinels are invalid.

**Step 2: Add or normalize helpers**

Use a single helper for unknown checks, for example:

```go
func (l Loc) IsUnknown() bool {
	return l.Start == -1 && l.End == -1
}
```

Do not introduce `End == 0` checks as unknown checks.

**Step 3: Run Loc proof**

```bash
go test -run 'TestFoundation|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS.

### Task 2.2: Make Corpus Loc Violations Fatal

**Files:**
- Modify: `oracle/parser/verify_parse_test.go`
- Test: `oracle/parser/verify_parse_test.go`

**Step 1: Add a negative test or fixture**

Inject or simulate one invalid Loc result and verify the Loc verifier fails.

**Step 2: Change verifier behavior**

`LOC VIOLATIONS > 0` must call `t.Errorf` or `t.Fatalf`, not only `t.Logf`.

**Step 3: Run proof**

```bash
go test -run 'TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS with `LOC VIOLATIONS: 0`.

### Task 2.3: Expand Synthetic Loc Coverage

**Files:**
- Create: `oracle/parser/loc_contract_test.go`

**Step 1: Add table-driven Loc coverage**

Cover SELECT, DML, CREATE TABLE, ALTER TABLE, PL/SQL, hints, aliases, quoted identifiers, Unicode identifiers, function calls, old-style `(+)`, and nested expressions.

**Step 2: Run proof**

```bash
go test -run 'TestOracleLocContract|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS.

## Phase 3: Keyword System

### Task 3.1: Create Keyword Category API

**Files:**
- Create: `oracle/parser/keywords.go`
- Modify: `oracle/parser/name.go`
- Test: `oracle/parser/keyword_class_test.go`

**Step 1: Expand golden tests**

Add bidirectional golden coverage for:

- reserved keywords
- nonreserved keywords
- context keywords
- type keywords
- function-like keywords
- pseudo-columns
- clause starters

**Step 2: Implement category API**

Preferred shape:

```go
type oracleKeywordCategory int

const (
	oracleKeywordIdentifier oracleKeywordCategory = iota
	oracleKeywordReserved
	oracleKeywordNonReserved
	oracleKeywordContext
	oracleKeywordType
	oracleKeywordFunction
	oracleKeywordPseudoColumn
	oracleKeywordClauseStarter
)
```

**Step 3: Run keyword proof**

```bash
go test -run 'TestOracleKeyword' -count=1 ./oracle/parser
```

Expected: PASS.

### Task 3.2: Replace Global Keyword-As-Identifier Behavior

**Files:**
- Modify: `oracle/parser/name.go`
- Modify: identifier call sites throughout `oracle/parser`

**Step 1: Add context tests**

Test object names, column names, aliases, dotted names, DB links, labels, roles, option values, and function names.

**Step 2: Introduce context-specific APIs**

Target APIs:

```go
func (p *Parser) parseObjectIdentifier() (string, error)
func (p *Parser) parseColumnIdentifier() (string, error)
func (p *Parser) parseAliasIdentifier() (string, error)
func (p *Parser) parseRoleIdentifier() (string, error)
func (p *Parser) parseOptionKeyword(allowed map[int]string) (string, error)
func (p *Parser) parseKeywordOrIdentifier() (string, error)
```

**Step 3: Convert call sites incrementally**

Start with `CREATE TABLE`, `FROM`, aliases, then DDL option parsers. Keep compatibility tests green.

**Step 4: Run proof**

```bash
go test -run 'TestOracleKeyword|TestStrictReservedKeywordIdentifiers|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS.

### Task 3.3: Eliminate String Keyword Matching Where Practical

**Files:**
- Modify: `oracle/parser/*.go`
- Test: `oracle/parser/keyword_class_test.go`

**Step 1: Add a source guard test**

Fail on new `strings.EqualFold` or broad `p.cur.Str ==` keyword checks for registered keywords unless allowlisted.

**Step 2: Convert high-risk matches**

Prioritize option and clause dispatch code that currently accepts arbitrary identifier-like tokens.

**Step 3: Run proof**

```bash
go test -run 'TestOracleKeyword|TestStrict' -count=1 ./oracle/parser
```

Expected: PASS.

## Phase 4: Soft-Fail And Strictness Matrix

### Task 4.1: Expand Soft-Fail Matrix

**Files:**
- Modify: `oracle/parser/soft_fail_strict_test.go`

**Step 1: Add table-driven scenario groups**

Add groups for:

- expression truncation
- predicate truncation
- SELECT truncation
- DML truncation
- DDL truncation
- PL/SQL truncation

**Step 2: Run RED**

```bash
go test -run 'TestSoftFail' -count=1 ./oracle/parser
```

Expected: FAIL for unimplemented families.

**Step 3: Implement family fixes**

Fix one family at a time. Do not add broad token skipping to pass tests.

**Step 4: Run proof**

```bash
go test -run 'TestSoftFail' -count=1 ./oracle/parser
```

Expected: PASS.

### Task 4.2: Expand Strictness Matrix

**Files:**
- Modify: `oracle/parser/soft_fail_strict_test.go`
- Modify: parser files by failing family

**Step 1: Add strict scenario groups**

Add groups for:

- duplicate clauses
- unknown options
- illegal keyword positions
- parenthesis balance
- statement separators
- clause dependency rules

**Step 2: Run RED**

```bash
go test -run 'TestStrict' -count=1 ./oracle/parser
```

Expected: FAIL for unimplemented families.

**Step 3: Implement with precise validation**

Known option sets should reject unknown tokens. Required grammar positions should error. Optional grammar positions should remain optional.

**Step 4: Run proof**

```bash
go test -run 'TestStrict|TestVerifyCorpus' -count=1 -v ./oracle/parser
```

Expected: PASS.

## Phase 5: Coverage And Documentation Status

### Task 5.1: Add Progress Reporting

**Files:**
- Create: `oracle/parser/parser_progress_test.go` or `oracle/parser/testdata/README.md`

**Step 1: Add a test or scriptable test log**

Report:

- total parse functions
- parse functions with error returns
- silent error discards
- Loc violations
- keyword golden mismatch counts
- strict/soft-fail scenario counts

**Step 2: Run proof**

```bash
go test -run 'TestOracleParserProgress' -count=1 -v ./oracle/parser
```

Expected: PASS with measured counts at target thresholds.

### Task 5.2: Update Defense Documentation

**Files:**
- Modify: `docs/PARSER-DEFENSE-MATRIX.md`
- Modify: `docs/engine-capability-guide.md`
- Modify: `oracle/quality/strategy.md`
- Modify: `oracle/parser/SCENARIOS-oracle-parser-completion.md`

**Step 1: Update measured status**

Change Oracle rows only when backed by fresh test output.

**Step 2: Remove resolved blind spots**

Do not remove a blind spot until a test fails without the implementation and passes with it.

**Step 3: Run docs-adjacent proof**

```bash
go test -run 'TestOracleParserProgress|TestVerifyCorpus|TestOracleKeyword|TestSoftFail|TestStrict' -count=1 -v ./oracle/parser
git diff --check
```

Expected: PASS.

## Final Acceptance

Oracle parser completion is ready to report when all are true:

- `go test -count=1 ./oracle/parser` passes.
- `go test -run 'TestVerifyCorpus' -count=1 -v ./oracle/parser` reports zero parse violations, zero Loc violations, and zero crashes.
- Parser contract test reports no bare required `parse*` functions outside explicit optional allowlist.
- Source guard reports no silent parser error discards.
- Keyword golden tests report zero missing/misclassified categories for the declared Oracle keyword set.
- Strict/soft-fail matrix tests pass.
- `docs/PARSER-DEFENSE-MATRIX.md` has measured Oracle status for L1, L2, L3, L5, and L7.
- `git diff --check` passes.
