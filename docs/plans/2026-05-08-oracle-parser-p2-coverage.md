# Oracle Parser P2 Coverage Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Move Oracle parser coverage from classified-but-coarse BNF accounting to P2 parser completeness: every Oracle BNF row is mapped to a structured AST target, every obvious skip/stub site is owned by a BNF family, and P2-target parser-visible clauses stop being silently skipped.

**Architecture:** Keep P2 strictly parser-layer. P2 requires structured AST for parser-visible statement kind, object name, option clauses, subclauses, and byte spans; it does not require catalog execution, dependency resolution, privilege validation, or storage semantics. Add a new P2 manifest and gate before changing parser behavior, then fix skip/stub hotspots in batches ordered by user-facing AST value and shared file risk.

**Tech Stack:** Go recursive-descent parser in `oracle/parser`, AST nodes in `oracle/ast`, BNF source files in `oracle/parser/bnf`, current BNF manifest in `oracle/parser/testdata/coverage/bnf_coverage.tsv`, optional Oracle reference harness in `oracle/parser/reference_oracle_test.go`.

## P2 Definition

P2 is the phase-one target for every Oracle BNF:

- Valid syntax parses without panic and without silent statement loss.
- The AST records statement kind, action, object type, object name, and source `Loc`.
- Parser-visible clauses are represented in typed AST fields or a deliberate structured option list.
- Required operands and delimiters still raise parser errors.
- Complex clauses may be generic option nodes at first, but they must preserve keyword, value, nested parenthesized payload, and byte span.
- No BNF is exempt merely because later catalog or semantic work is needed.

For completion, the P2 target set is every row in `oracle/parser/bnf/*.bnf` except rows proven to be `p2_not_parser_visible` after the parser-visible AST is complete. Target rows may use `p2_partial`, `p2_stub`, or `p2_deferred` while work is in progress, but those statuses are not completion states.

P2 explicitly excludes:

- Object existence checks.
- Privilege and role validation.
- Tablespace, storage, flashback, Data Guard, and PDB execution semantics.
- Dependency graph, migration planning, or deparse fidelity beyond preserving parsed structure.

## Status Vocabulary

Replace parser-layer use of `catalog` with P2-specific statuses:

| Status | Meaning | Default Gate |
|---|---|---|
| `p2_done` | Parser-visible BNF surface is structurally represented in AST and covered by positive, negative, and Loc tests. | Allowed |
| `p2_partial` | Statement parses and has some AST structure, but one or more parser-visible clauses are skipped or collapsed too coarsely. | Allowed with owned next action |
| `p2_stub` | Statement returns a coarse placeholder AST or skips most of the body. | Allowed only while explicitly tracked |
| `p2_deferred` | No useful P2 parser implementation yet. | Allowed only with phase owner |
| `p2_not_parser_visible` | Remaining behavior is truly semantic/catalog-only after P2 AST is complete. | Allowed |
| `p2_unknown` | Not reviewed under P2 rules. | Never allowed |

Only `p2_done` and justified `p2_not_parser_visible` are valid completion statuses. `p2_partial`, `p2_stub`, and `p2_deferred` are execution-time debt states.

The existing `covered/partial/deferred/catalog` manifest remains historical input. The new P2 manifest is the source of truth for this effort.

## Current Baseline

Current `bnf_coverage.tsv` status counts:

- `covered`: 11
- `partial`: 86
- `deferred`: 47
- `catalog`: 27

Current skip/stub hotspot counts from parser source scanning:

- `oracle/parser/alter_table.go`: 57
- `oracle/parser/create_admin.go`: 29
- `oracle/parser/alter_misc.go`: 25
- `oracle/parser/create_table.go`: 18
- `oracle/parser/create_index.go`: 16

Important current examples:

- `parseAdminDDLStmt` parses action/object/name, then skips remaining body.
- `parseAlterStmt` can skip unknown ALTER targets and return `nil, nil`.
- `parseAlterGeneric` is documented as a placeholder.
- `ALTER TABLE` has many `skipAlterTableClauseDetails` branches for partition/storage/ILM/supplemental-log details.
- `CREATE TYPE BODY` uses `skipToEndBlock` as a placeholder for full PL/SQL body parsing.

## Canonical Proof Commands

Run before and after every batch:

```bash
go test -run 'TestOracleP2|TestOracleCoverage|TestOracleParserProgress|TestVerifyCorpus' -count=1 -v ./oracle/parser
go test -count=1 ./oracle/parser
git diff --check
```

`TestVerifyCorpus` is a smoke/regression proof only. It is not the P2 malformed-syntax gate, because accepted invalid SQL is currently reported as diagnostic output rather than a hard failure. Each P2 batch must add focused negative tests for required operands, delimiters, and clause boundaries.

Optional Oracle reference proof:

```bash
ORACLE_PARSER_REF_DSN="$ORACLE_PARSER_REF_DSN" go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
ORACLE_PARSER_REF_CONTAINER=1 go test -tags oracle_ref -run 'TestOracleReference|TestOracleVReservedWordsKeywordAudit' -count=1 -v ./oracle/parser
```

## Task 1: Add P2 BNF Manifest

**Files:**
- Create: `oracle/parser/testdata/coverage/p2_coverage.tsv`
- Create: `oracle/parser/p2_coverage_test.go`
- Modify: `oracle/parser/SCENARIOS-oracle-parser-completion.md`

**Step 1: Write the failing gate**

Add `TestOracleP2BNFManifestCompleteness`.

It must:

- Require exactly one row per `oracle/parser/bnf/*.bnf`.
- Validate `p2_status` against the P2 status vocabulary and fail on `p2_unknown`.
- Fail if any previous `catalog` row is carried over unchanged as a P2 status.
- Require non-empty `p2_surface`, `ast_target`, `parser_entrypoint`, `current_gap`, `next_action`, `positive_test`, `negative_test`, and `loc_test`.
- Reject placeholder evidence values such as `todo`, `unknown`, `n/a`, and `none` in required proof columns.
- Verify `bnf_file` exists under `oracle/parser/bnf`.
- For `p2_done` rows, verify referenced test files exist and at least one referenced test name is present in the file. Where feasible, verify `parser_entrypoint` and `ast_target` symbols exist in `oracle/parser` and `oracle/ast`.
- Require every row with `p2_partial`, `p2_stub`, or `p2_deferred` to have an owner phase.
- Require every `p2_not_parser_visible` row to name the parser-visible AST evidence that made the remaining behavior semantic-only.

Expected initial result before manifest exists: FAIL because `p2_coverage.tsv` is missing.

**Step 2: Create `p2_coverage.tsv`**

Columns:

```text
bnf_file	family	p2_status	p2_surface	parser_entrypoint	ast_target	current_gap	skip_hotspot	owner_phase	next_action	positive_test	negative_test	loc_test	reference_evidence
```

Initial classification rule:

- Current `covered` rows become `p2_done` only after verifying clause structure, not automatically.
- Current `partial` rows usually become `p2_partial`.
- Current `deferred` rows become `p2_deferred`.
- Current `catalog` rows become `p2_stub` or `p2_partial` unless the parser already preserves the visible clause structure.

**Step 3: Run the gate**

Run:

```bash
go test -run 'TestOracleP2BNFManifestCompleteness' -count=1 -v ./oracle/parser
```

Expected: PASS with 171 rows, zero `p2_unknown`, no historical `catalog` status, and valid proof columns for any `p2_done` row.

## Task 2: Add Skip/Stub Inventory

**Files:**
- Create: `oracle/parser/testdata/coverage/p2_skip_inventory.tsv`
- Modify: `oracle/parser/p2_coverage_test.go`

**Step 1: Write the failing drift gate**

Add `TestOracleP2SkipInventoryDrift`.

It must scan non-test `oracle/parser/*.go` for these patterns:

- `skipToSemicolon`
- `skipAlterTableClauseDetails`
- `skipParenthesizedBlock`
- `skipToNextClause`
- `parseAdminDDLStmt`
- `parseAlterGeneric`
- comments containing `placeholder`
- comments containing `skip remaining`, `skip details`, or `skip unrecognized`
- required parse paths returning `nil, nil` after consuming tokens

The test fails when a new skip/stub site is not listed in the inventory.

**Step 2: Create inventory**

Columns:

```text
site	file	function	pattern	bnf_family	bnf_file	p2_status	reason	next_action
```

Use stable `file:function:pattern` identities instead of raw line numbers where possible, because line numbers will churn during implementation.

**Step 3: Link inventory to manifest**

The test must fail if `p2_skip_inventory.tsv` references a BNF file absent from `p2_coverage.tsv`, or if a row marked `p2_done` still has an owned skip/stub site.

## Task 3: Remove Silent Statement Loss

**Files:**
- Modify: `oracle/parser/alter_misc.go`
- Modify: `oracle/parser/create_admin.go`
- Test: `oracle/parser/soft_fail_strict_test.go`

**Scope:**

Do not structure every ALTER/admin clause yet. First eliminate parser paths that consume a real statement prefix, skip to semicolon, and return `nil, nil`.

**Test cases:**

- `ALTER PUBLIC SELECT x` returns a syntax error, not success with no AST.
- `ALTER SHARED PUBLIC SELECT x` returns a syntax error.
- Unknown top-level `ALTER <known prefix but invalid target>` returns a syntax error.
- Unknown nested `CREATE SCHEMA` child returns a syntax error with position.

**Implementation rule:**

If the statement prefix is recognized enough to enter a parser branch, malformed or unsupported parser-visible syntax must return an error or a structured `AdminDDLStmt`, never silent nil success.

## Task 4: Convert Admin DDL From Body Skip To Structured Option List

**Files:**
- Modify: `oracle/ast/parsenodes.go`
- Modify: `oracle/ast/outfuncs.go`
- Modify: `oracle/parser/create_admin.go`
- Test: `oracle/parser/create_admin_test.go`
- Test: `oracle/parser/database_test.go`

**Scope:**

Target current `catalog` rows first, under P2 rules:

- `ADMINISTER-KEY-MANAGEMENT.bnf`
- `CREATE/ALTER/DROP DATABASE LINK`
- `CREATE/ALTER/DROP TABLESPACE`
- `CREATE/ALTER/DROP TABLESPACE SET`
- `CREATE/ALTER/DROP DATABASE`
- `CREATE/ALTER/DROP PLUGGABLE DATABASE`
- `CREATE/DROP RESTORE POINT`
- `CREATE PFILE`
- `CREATE SPFILE`
- `CREATE CONTROLFILE`
- `CREATE/ALTER/DROP DISKGROUP`
- `FLASHBACK DATABASE`

**Design:**

Add a generic parser-layer option representation, for example:

```go
type DDLClause struct {
	Name string
	Args *List
	Raw string
	Loc Loc
}
```

Use this only where typed AST fields would be premature. It must preserve clause boundaries and nested parenthesized content. It is not acceptable for a generic clause node to become a raw statement tail.

Minimum `DDLClause` semantics:

- `Name` is the leading clause keyword or keyword phrase, normalized consistently with existing AST output.
- `Loc` covers the clause keyword through the last token owned by that clause, not the whole statement.
- `Args` or a structured payload records parser-visible operands and nested option nodes when the grammar exposes them.
- `Raw` is allowed only as a debug/source-preservation field; tests must not rely on `Raw` alone as proof of P2 structure.
- Clause parsing must stop at sibling clause boundaries and must keep balanced parentheses, quoted strings, and PL/SQL blocks from splitting incorrectly.
- Required clause operands must still error when missing; a generic clause cannot hide malformed required syntax.

**P2 acceptance:**

- The statement AST has action, object type, object name when present, and a list of clauses.
- Parser-visible clauses are not dropped.
- Tests assert clause count, order, keyword spans, operand spans, and sibling boundaries.
- Tests include one case where a nested parenthesized payload contains tokens that look like sibling clauses.
- Tests include one malformed case proving a missing required operand does not become a raw generic payload.
- Malformed required syntax errors.
- Loc covers the statement and every clause node.

## Task 5: ALTER TABLE P2 Structure

**Files:**
- Modify: `oracle/ast/parsenodes.go`
- Modify: `oracle/ast/outfuncs.go`
- Modify: `oracle/parser/alter_table.go`
- Test: `oracle/parser/alter_table_test.go`
- Test: `oracle/parser/loc_node_coverage_test.go`

**Scope:**

Replace the most visible `skipAlterTableClauseDetails` branches with structured commands:

- `ADD SUPPLEMENTAL LOG`
- `ADD PARTITION` / `ADD SUBPARTITION`
- `MODIFY PARTITION`
- `DROP PARTITION`
- `TRUNCATE PARTITION`
- `MOVE PARTITION`
- `SPLIT PARTITION`
- `MERGE PARTITIONS`
- `EXCHANGE PARTITION`
- `SET SUBPARTITION TEMPLATE`
- generic `SET` table property
- storage, logging, compression, cache, parallel, row movement, flashback archive, ILM/inmemory property clauses

**P2 acceptance:**

Each branch may initially use generic `DDLClause` payloads for deep storage attributes, but command type, target partition/subpartition name, and source span must be explicit.

## Task 6: CREATE TABLE And CREATE INDEX P2 Structure

**Files:**
- Modify: `oracle/ast/parsenodes.go`
- Modify: `oracle/ast/outfuncs.go`
- Modify: `oracle/parser/create_table.go`
- Modify: `oracle/parser/create_index.go`
- Test: `oracle/parser/create_table_test.go`
- Test: `oracle/parser/create_index_test.go`

**CREATE TABLE scope:**

- external table clauses
- LOB storage clauses
- nested table storage clauses
- XMLType storage clauses
- segment attributes
- table compression and inmemory clauses
- partition and subpartition descriptions
- ILM and attribute clustering clauses

**CREATE INDEX scope:**

- local/global partitioned index clauses
- domain index parameters
- XML/search/spatial/vector index parameters
- physical attributes and parallel/logging clauses
- deferred invalidation / online / usable-unusable states

**P2 acceptance:**

The AST must preserve every parser-visible option as typed fields or generic option nodes with Loc. It does not need to validate tablespace existence, index method availability, or storage semantics.

## Task 7: PL/SQL Body P2 Structure

**Files:**
- Modify: `oracle/parser/create_type.go`
- Modify: `oracle/parser/create_proc.go`
- Modify: `oracle/parser/create_trigger.go`
- Modify: `oracle/parser/plsql_block.go`
- Modify: `oracle/ast/parsenodes.go`
- Modify: `oracle/ast/outfuncs.go`
- Test: `oracle/parser/create_type_test.go`
- Test: `oracle/parser/plsql_block_test.go`

**Scope:**

- Replace `skipToEndBlock` placeholder for `CREATE TYPE BODY`.
- Preserve member function/procedure declarations and bodies.
- Preserve package body initialization block and subprogram definitions.
- Preserve trigger body timing sections, including compound trigger sections.

**P2 acceptance:**

PL/SQL body internals do not need semantic resolution, but parser-visible declarations, statements, and section boundaries must be represented.

## Task 8: Deferred Modern Oracle Feature Families

**Files:**
- Modify parser files according to each family.
- Modify `oracle/parser/testdata/coverage/p2_coverage.tsv`.
- Add focused tests per family.

**Families:**

- analytic view
- attribute dimension
- hierarchy
- domain
- JSON relational duality view
- property graph
- vector index
- MLE module/environment
- Java and library
- lockdown profile
- materialized zonemap
- inmemory join group
- indextype/operator/outline

**Execution shape:**

Handle these after Tasks 3-7 because they are less likely to be common user AST dependencies and may require new AST node families.

## Task 9: Expand Reference Oracle Rows For P2

**Files:**
- Modify: `oracle/parser/testdata/coverage/reference_oracle.tsv`
- Modify: `oracle/parser/reference_oracle_test.go`

**Scope:**

Add accept/reject rows for every P2 phase above. The reference test remains optional, but the manifest should become broad enough to catch parser-visible mismatches. A `p2_done` row for a risky or grammar-rich surface must either link reference evidence or give an explicit `reference_evidence` reason such as `no_reference_runtime`, `oracle_version_gap`, or `local_structure_tests_only`.

**P2 acceptance:**

Every P2 batch adds at least:

- one valid reference row
- one malformed required-syntax reference row
- one keyword/identifier edge row when applicable

## Implementation Order

1. Land Task 1 and Task 2 gates with no parser behavior changes.
2. Land Task 3 to remove silent statement loss.
3. Land Task 4 for previous `catalog` rows, because those are currently most likely to be wrongly discounted.
4. Land Task 5 for `ALTER TABLE`, the largest skip hotspot.
5. Land Task 6 for `CREATE TABLE` and `CREATE INDEX`.
6. Land Task 7 for PL/SQL body structure.
7. Land Task 8 feature families in independent batches.
8. Land Task 9 reference-oracle breadth as each implementation batch completes.

## Completion Criteria

This P2 effort is complete when:

- `p2_coverage.tsv` has exactly one row for every `oracle/parser/bnf/*.bnf` file.
- All target rows are `p2_done`; non-target rows are only justified `p2_not_parser_visible`.
- The manifest has zero `p2_unknown`, zero `p2_stub`, zero `p2_deferred`, and zero `p2_partial` rows.
- No BNF row uses old `catalog` as a parser-layer status.
- Every tracked skip/stub site is mapped to a BNF row and owner phase.
- No `p2_done` row has an active skip/stub inventory row.
- Every `p2_done` row links positive, negative, and Loc evidence that the gate validates.
- Generic clause AST tests prove clause boundaries, nested payload preservation, and malformed required syntax handling.
- Focused P2 negative tests pass for malformed required syntax; `TestVerifyCorpus` remains a smoke test and still reports zero parse violations, zero Loc violations, and zero crashes.
- `go test -count=1 ./oracle/parser` passes.
