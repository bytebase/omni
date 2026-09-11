# SCENARIOS: MSSQL Keyword System Alignment with SqlScriptDOM

## Goal

Align omni's MSSQL parser keyword system with SqlScriptDOM TSql170. All 6 enforcement tests PASS:
- `TestKeywordCompleteness`: PASS — 0 unregistered keywords
- `TestNoStringKeywordMatch`: PASS — 0 string-based keyword matches
- `TestKeywordClassification`: PASS — 180 Core + 372 Context = 552 keywords classified
- `TestCoreKeywordNotIdentifier`: PASS — Core keywords rejected as unquoted identifiers
- `TestContextKeywordAsIdentifier`: PASS — Context keywords accepted as identifiers + bare aliases
- `TestKeywordCasePreservation`: PASS — Keyword tokens preserve original case

Known remaining mismatch (strictness, not keyword classification):
- `SELECT 1 into FROM dbo.t` — SQL Server rejects, omni accepts (`INTO` as bare alias)

## Reference

- Enforcement tests: `mssql/parser/keyword_classification_test.go`
- Build check: `go build ./mssql/...`
- Full regression: `go test ./mssql/... -count=1`
- SqlScriptDOM reference: `../SqlScriptDOM/SqlScriptDom/Parser/TSql/`

---

## Phase 1: Infrastructure (sequential — shared files)

### Section 1.1: Lexer case preservation (1 site)

- [x] keyword tokens preserve original case in Str field (currently lowercased)
- [x] `SELECT MyPartition FROM t` → AST stores "MyPartition" not "mypartition"
- [x] `CREATE TABLE MyWindow (a INT)` → table name preserves case
- [x] non-keyword identifiers still preserve case (no regression)
- [x] reverseKeywordMap still works for completion (maps token type → uppercase string)

Verification: `go test ./mssql/parser/ -run TestKeywordCasePreservation -count=1`

### Section 1.2: Classification infrastructure

- [x] `KeywordCategory` type defined (Core, Context)
- [x] `Keyword` struct with Name, Token, Category fields
- [x] classification table mapping each keyword to its category
- [x] 180 Core keywords match SqlScriptDOM golden list exactly
- [x] 372 Context keywords added to keywordMap and classified (was 279 in original estimate)
- [x] all existing keywordMap entries classified (no unclassified)
- [x] `lookupKeywordCategory(token int) KeywordCategory` helper available
- [x] `isContextKeyword(token int) bool` helper available
- [x] completion: reverseKeywordMap includes all new context keywords
- [x] completion: `Collect()` on empty input returns context keywords as candidates

Verification: `go test ./mssql/parser/ -run TestKeywordClassification -count=1` + `go test ./mssql/parser/ -run TestKeywordCompleteness -count=1`

### Section 1.3: parseIdentifier classification-aware

- [x] `isIdentLike()` rejects Core keywords, accepts Context keywords and tokIDENT
- [x] `parseIdentifier()` rejects Core keywords, accepts Context keywords and tokIDENT
- [x] bracket-quoted Core keywords accepted: `[select]` parses as identifier
- [x] `isAnyKeywordIdent()` accepts ALL keywords (for alias positions)
- [x] Core keyword as table name fails: `CREATE TABLE select (a INT)` → error
- [x] Context keyword as table name works: `CREATE TABLE window (a INT)` → success

Verification: `go test ./mssql/parser/ -run TestCoreKeywordNotIdentifier -count=1`

### Section 1.4: Fix tokIDENT-gated paths (85 sites across 18 files)

These `.Type == tokIDENT` checks were widened to also accept Context keyword tokens via `isIdentLike()` and `isContextKeyword()`.

- [x] parser.go: statement dispatch lookahead uses isIdentLike/isContextKeyword
- [x] create_table.go: PERSISTED, SPARSE, HIDDEN, MASKED, ENCRYPTED, GENERATED, ALWAYS, NODE, EDGE, PERIOD, SYSTEM_TIME, FILESTREAM_ON, TEXTIMAGE_ON
- [x] select.go: bare alias detection, FETCH NEXT/FIRST/ONLY/ROW/TIES, OVER clause
- [x] merge.go: USING, MATCHED, TARGET, SOURCE
- [x] execute.go: OUT parameter, exec string parsing
- [x] alter_table.go: PERIOD, SYSTEM_TIME
- [x] backup_restore.go: option parsing
- [x] create_proc.go: OUT parameter
- [x] name.go: base identifier checks
- [x] update_delete.go: output clause context
- [x] control_flow.go: TIME in WAITFOR
- [x] create_database.go: database options
- [x] create_sequence.go: sequence options
- [x] create_trigger.go: trigger timing
- [x] drop.go: CASCADE
- [x] expr.go: expression parsing
- [x] fulltext.go: fulltext options
- [x] type.go: type name parsing

Verification: `go test ./mssql/parser/ -run TestContextKeywordAsIdentifier -count=1` + `go test ./mssql/... -count=1`

---

## Phase 2: String match migration (per-file, parallelizable after Phase 1)

All `strings.EqualFold` and `matchesKeywordCI` calls have been replaced with token type checks. `matchesKeywordCI` was never implemented / already removed. Only 1 legitimate `strings.EqualFold` remains (AST field check in create_table.go:554 for `col.DataType.Name`).

### Section 2.1a: parser.go — CREATE dispatch (~70 sites)

- [x] all CREATE statement recognition uses keyword tokens

### Section 2.1b: parser.go — ALTER/DROP dispatch (~70 sites)

- [x] all ALTER/DROP statement recognition uses keyword tokens

### Section 2.1c: parser.go — remaining dispatch + eqFold (~72 sites)

- [x] parser.go contributes 0 violations

### Section 2.2: service_broker.go (60 sites)

- [x] all keyword matching uses token types

### Section 2.3: fulltext.go (45 sites)

- [x] all keyword matching uses token types

### Section 2.4: select.go (38 sites)

- [x] all keyword matching uses token types

### Section 2.5: security_principals.go (36 sites)

- [x] all keyword matching uses token types

### Section 2.6: create_table.go (34 sites) + alter_table.go (18 sites)

- [x] all keyword matching uses token types (1 legitimate EqualFold on AST field retained)

### Section 2.7: server.go (30 sites) + availability.go (23 sites)

- [x] all keyword matching uses token types

### Section 2.8: alter_objects.go (23 sites) + external.go (22 sites)

- [x] all keyword matching uses token types

### Section 2.9: utility.go (20 sites) + endpoint.go (17 sites)

- [x] all keyword matching uses token types

### Section 2.10: event.go (15 sites) + security_misc.go (14 sites) + security_audit.go (10 sites)

- [x] all keyword matching uses token types

### Section 2.11: backup_restore.go (14 sites) + create_database.go (13 sites)

- [x] all keyword matching uses token types

### Section 2.12: declare_set.go (12 sites) + partition.go (10 sites) + security_keys.go (10 sites)

- [x] all keyword matching uses token types

### Section 2.13: remaining 18 files (70 sites total)

- [x] all keyword matching uses token types across all files

### Section 2.14: Delete matchesKeywordCI function

- [x] `matchesKeywordCI` does not exist in the codebase (never implemented or already removed)

---

## Proof

### Global proof
All 6 enforcement tests pass. Full regression green (excluding pre-existing `TestKeywordOracleCoreAsIdentifier/alias_into_bare` strictness mismatch).
