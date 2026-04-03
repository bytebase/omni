# Parser Defense Matrix

Every SQL parser in omni needs a systematic suite of defensive tests to ensure correct, safe, and predictable behavior across all inputs. This document defines the defense layers, tracks per-engine coverage, and guides prioritization.

## Defense Layers

### L1: Soft-Fail (Error Recovery)

**Goal**: The parser never crashes or silently discards errors on invalid input. It returns meaningful error messages instead.

**Coverage points**:
- All parse functions return `(T, error)` dual values; no silent discard via `_, _, _ :=`
- Expression truncation: binary operators (`SELECT 1 +`), comparison, logical, unary
- Predicate truncation: `BETWEEN`/`LIKE`/`IN`/`REGEXP` missing right operand
- Special expressions: `CASE WHEN` missing `END`, `CAST` missing type, incomplete `CONVERT`
- Statement-level truncation: SELECT/INSERT/UPDATE/DELETE/DDL cut off at each clause
- Error message quality: distinguish "at end of input" vs "at or near 'token'"

**Verification**: Error test files — one test case per truncation pattern, asserting error message content.

---

### L2: Strict Parsing (Rejection of Invalid SQL)

**Goal**: The parser does not accept SQL that the target database would reject. Avoid false positives.

**Coverage points**:
- JOIN keyword enforcement: `INNER`/`LEFT`/`RIGHT`/`CROSS`/`NATURAL` must be followed by `JOIN`
- Clause dependencies: no `WHERE`/`GROUP BY`/`HAVING`/`ORDER BY` without `FROM`
- DML structural completeness: `DELETE` requires `FROM`, `INSERT` requires VALUES/SELECT/SET, `UPDATE` requires `SET`
- DDL chain completeness: `IF NOT EXISTS`/`IF EXISTS` must be fully spelled out
- Parenthesis balance: unclosed parentheses must produce an error
- Trailing garbage detection: unexpected tokens after a valid statement must error
- Implicit alias vs rejection: which cases allow implicit aliases, which must reject

**Verification**: Strictness test files — each pattern has a "should reject" case + a "should accept" control.

---

### L3: Keyword System (Keyword Classification)

**Goal**: Keyword reserved/non-reserved classification matches the target database. No hardcoded eqFold bypasses.

**Coverage points**:
- Keyword registration completeness: all target database keywords registered in the lexer
- Keyword classification correctness: reserved / non-reserved matches target database docs
- No eqFold residue: registered keywords are not matched via `strings.EqualFold` hardcoding
- Context-sensitive keywords: `t.select` (after dot), `CREATE TABLE select` (reserved word as identifier)
- Lexer behavior: dot-context suppression, @@variable tokenization, string escaping

**Verification**: Golden tests (TestKeywordCompleteness / TestKeywordClassification / TestNoEqFoldForRegisteredKeywords) driven to zero differences.

---

### L4: Syntax Coverage

**Goal**: The parser can parse all valid SQL syntax of the target database.

**Coverage points**:
- Type synonyms: `REAL`/`DEC`/`FIXED`/`INT1`-`INT8` aliases
- Window function names: `RANK`/`DENSE_RANK`/`ROW_NUMBER`/`LAG`/`LEAD`
- Interval units: `DAY`/`HOUR`/`MINUTE`/`SECOND` and compound units
- Partition syntax: `PARTITION BY HASH`/`KEY`/`RANGE`/`LIST`
- Engine-specific syntax: PG `RETURNING`, MySQL `ON DUPLICATE KEY`, MSSQL `OUTPUT`/`MERGE`, Oracle `CONNECT BY`
- Statement type completeness: DDL/DML/DCL/TCL/admin full coverage

**Verification**: Corpus syntax files + golden verification (ideally oracle-verified against the real database).

---

### L5: Corpus Testing (Multi-Layer Corpus Verification)

**Goal**: Structured corpus files drive multi-dimensional regression testing.

**Coverage points**:
- Parse correctness: SQL annotated `@valid: true/false`, verify parser accept/reject matches annotation
- Crash protection: no input should cause a panic
- Loc field validity: AST nodes have `Start >= 0`, `End > Start`
- Error position reasonableness: error positions on failed parses fall within input range
- Round-trip consistency: parse -> deparse -> parse produces the same result (where a deparser exists)

**Verification**: `quality/corpus/*.sql` golden files + multi-layer verification harness.

---

### L6: Completion Instrumentation

**Goal**: The parser provides correct completion candidates on partial input.

**Coverage points**:
- Statement-level dispatch: identify the statement type being written
- Clause-level context: `SELECT ... FROM ... WHERE |` should complete column names
- Expression context: function arguments, subqueries, CASE WHEN internals
- Identifier resolution: qualified name (`schema.table.column`) completion
- Full DML/DDL/DCL coverage: not just SELECT

**Verification**: Completion test files — cursor position -> expected candidates.

---

### L7: Position Tracking (Loc)

**Goal**: Every AST node carries accurate source position information.

**Coverage points**:
- Token.End accuracy: correct end position for every token
- Loc.Start / Loc.End consistency: all AST nodes have positions set
- Sentinel value handling: unset Loc uses a consistent zero value
- Line:column conversion: byte offset -> line:column public API
- Nesting correctness: child expression Loc falls within parent node Loc range

**Verification**: Loc accuracy tests + corpus-level Loc validity assertions.

---

## Per-Engine Coverage Matrix

```
            L1         L2         L3         L4         L5         L6         L7
            Soft-Fail  Strict     Keyword    Coverage   Corpus     Complete   Loc
            ---------  ---------  ---------  ---------  ---------  ---------  ---------
PG          Done       --         --         Partial    --         Done       Partial
MySQL       Partial*   Partial*   Done       Done       Done       Done       Done
MSSQL       Done       Partial    Done       --         --         Done*      Done
Oracle      --         --         --         --         Partial    --         Partial
CosmosDB    --         --         N/A        --         --         --         --
MongoDB     Partial    --         N/A        --         --         --         --
```

Legend: `Done` = complete, `Partial` = in progress, `--` = not started, `N/A` = not applicable

`*` marks verified status (see notes below):
- MySQL L1: Core path complete (138 test runs), but ~18 silent error discards remain in production code
- MySQL L2: 111 strictness scenarios complete, but 3 checks removed (WHERE/GROUP BY/HAVING requires FROM — MySQL allows these without FROM)
- MSSQL L2: 14 oracle mismatches remaining (12 option validation + 2 multi-statement). Option validation is the L2 Strict core work.
- MSSQL L6: Core instrumentation complete (9 phases, 3659 tests), but 4 secondary CREATE statements uninstrumented + catalog resolution stubbed

### PG
- **L1 Soft-Fail**: Dual-return migration and soft-fail fixes complete
- **L2 Strict**: Not systematically tested
- **L3 Keyword**: Not audited
- **L4 Coverage**: pgregress regression suite provides broad coverage, but not systematically quantified
- **L5 Corpus**: No structured corpus files
- **L6 Completion**: Full completion engine + tests
- **L7 Loc**: SCENARIOS-pg-loc exists, partially complete

### MySQL
- **L1 Soft-Fail**: Core path complete — 32 silent errors fixed, 138 test runs across 21 test functions. **Gap: ~18 silent error discards (`_, _ :=` / `_, _ =`) remain in production parser code.** These are in non-critical paths (mostly `parseIdentifier` result ignoring) but violate the "no silent discard" principle.
- **L2 Strict**: 111 strictness scenarios verified. **3 checks removed** (WHERE/GROUP BY/HAVING "requires FROM") — these were too strict; MySQL 8.0 allows `SELECT 1 WHERE 1=1` without FROM. Current strictness aligns with MySQL 8.0 actual behavior.
- **L3 Keyword**: **Done.** Full alignment with MySQL 8.0 `lex.h` + `sql_yacc.yy`. 4 golden tests all at zero: TestKeywordCompleteness (0 missing from 796 non-operator/non-hint keywords), TestKeywordClassification (0 misclassified), TestNoEqFoldForRegisteredKeywords (0 violations), TestKeywordFunctions (0 failures across 78 keyword-named functions). 6-category classification system (reserved, unambiguous, ambiguous 1-4). 5 context-dependent identifier functions (parseIdent, parseLabelIdent, parseRoleIdent, parseLvalueIdent, parseKeywordOrIdent). Lexer fixes: dot-context suppression (no keyword lookup after `.`), @@variable multi-token emission, string escaping (`\b`, `\Z`, `\_`, `\%`).
- **L4 Coverage**: 91 oracle-verified corpus cases (against real MySQL 8.0 via testcontainers), 78 keyword-function tests, type synonym coverage (REAL/DEC/FIXED/LONG/INT1-8/FLOAT4/8), ALTER TABLE PARTITION BY, EXTRACT function, INTERVAL unit validation. ParseError includes RelatedText for error context.
- **L5 Corpus**: 5 golden files (`01-select.sql` through `05-dcl.sql`), multi-layer verification harness (parse correctness, crash protection, Loc validity, error position).
- **L6 Completion**: Full completion engine + tests.
- **L7 Loc**: Loc fields fully implemented across all AST node types (235 Loc instances in parsenodes.go). TestLocAudit covers 100+ SQL statement patterns with reflection-based walker. Line:column conversion available.

### MSSQL
- **L1 Soft-Fail**: **Done.** 234 silent error discards fixed (TestNoSilentErrorDiscard at 0). All parse functions propagate errors.
- **L2 Strict**: **In progress.** Oracle tests (SQL Server 2022 via testcontainers) identified 14 remaining mismatches:
  - 12 "omni too permissive" — `isAnyKeywordIdent()` in 459 option-parsing positions accepts any keyword as option name (e.g., `WITH (SELECT = 1)` accepted). Root cause: omni lacks option-name validation. SqlScriptDOM uses `OptionsHelper.ParseOption()` pattern — accept `Identifier` token then validate against known option list, reject unknown.
  - 2 "omni too permissive" — Core keywords `SET`/`INTO` trigger multi-statement parsing instead of error (`SELECT 1 set FROM t` → two statements). SqlScriptDOM/PARSEONLY rejects.
  - Context keyword disambiguation (4 cases) and bare alias (4 cases) already fixed.
  - **Next step**: Implement option-name validation per SqlScriptDOM's OptionsHelper pattern. This is the L2 Strict core work for MSSQL.
- **L3 Keyword**: **Done.** Full alignment with SqlScriptDOM TSql170. 6 golden enforcement tests all pass:
  - TestKeywordCompleteness: 0 unregistered keywords (509 total: 180 Core + 329 Context)
  - TestNoStringKeywordMatch: 0 string-based keyword matches (eliminated 746 eqFold + matchesKeywordCI)
  - TestKeywordClassification: Core set = SqlScriptDOM 180, bidirectional golden list
  - TestCoreKeywordNotIdentifier: Core keywords rejected as unquoted identifiers
  - TestContextKeywordAsIdentifier: Context keywords accepted as identifiers + bare aliases
  - TestKeywordCasePreservation: keyword tokens preserve original case
  - `matchesKeywordCI` function deleted. All keyword matching via token type checks.
- **L4 Coverage**: compare_test.go has 193 test functions, but coverage scope not systematically quantified. Oracle test infrastructure established (SQL Server 2022 testcontainers).
- **L5 Corpus**: No structured corpus files
- **L6 Completion**: 9 phases complete, 3659 tests passing. Core DML/DDL/control-flow/security instrumented. **Gap: 4 secondary CREATE statements (SEQUENCE, SYNONYM, TYPE, STATISTICS) and non-TABLE ALTER statements lack instrumentation. Catalog resolution is stubbed** (`TODO: resolve against mssql/catalog once it exists`).
- **L7 Loc**: All 6 phases complete, 179 scenarios verified, 350 migration sites. Public API with line:column conversion. **No gaps found.**

### Oracle
- **L1 Soft-Fail**: Not started (SCENARIOS-oracle-foundation focuses on Loc and ParseError infrastructure)
- **L2 Strict**: Not started
- **L3 Keyword**: Not audited
- **L4 Coverage**: Not systematically quantified
- **L5 Corpus**: quality/corpus/ directory exists, partial coverage
- **L6 Completion**: Not started
- **L7 Loc**: SCENARIOS-oracle-foundation in progress

### CosmosDB / MongoDB
- Smaller engines — lower defense priority than the four SQL engines
- MongoDB has basic error tests
- CosmosDB has virtually no defensive tests

---

## Layer Dependency Graph

The defense layers are not independent — they form a clear dependency chain:

```
Infrastructure          Functional              Verification
--------------          ----------              ------------
L7 Loc --------+
               +--> L1 Soft-Fail --> L2 Strict --> L4 Coverage
L3 Keyword ----+                                        |
               +--> L6 Completion                       v
                                                   L5 Corpus
```

- **L7 Loc** is foundational: error messages need positions (L1), completion needs cursor-to-AST mapping (L6), corpus verification checks Loc validity (L5)
- **L3 Keyword** is foundational: correct lexer tokenization is prerequisite for parsing, strictness, and completion
- **L1 Soft-Fail** depends on L7: error messages must report `line:column`
- **L2 Strict** depends on L1: the parser must be able to report errors before we discuss which errors to report
- **L6 Completion** depends on L7 + L3: needs accurate positions and correct token stream
- **L5 Corpus** is the terminal verification layer: multi-dimensional regression after all other layers are in place

Recommended progression: **infrastructure first (L7, L3), then functional layers (L1, L2, L6), then verification (L4, L5)**.

---

## Priority Recommendations

Ordered by dependency chain. Within each tier, sorted by impact.

### Tier 1: Infrastructure (Prerequisites for All Other Layers)

| Engine | Layer | Rationale |
|--------|-------|-----------|
| MSSQL | L1 Soft-Fail (cleanup) | 66 silent error discards remain; must finish before L2 Strict can begin |
| MySQL | L1 Soft-Fail (cleanup) | 18 silent error discards remain in production parser code |
| Oracle | L7 Loc | SCENARIOS-oracle-foundation in progress but incomplete; blocks all downstream layers |
| PG | L7 Loc | SCENARIOS-pg-loc partially complete; needs push to full coverage |
| MSSQL | L3 Keyword | Needs audit for eqFold hardcoding issues; affects lexer correctness |

### Tier 2: Functional (Directly Affects User Experience)

| Engine | Layer | Rationale |
|--------|-------|-----------|
| Oracle | L1 Soft-Fail | Missing basic error handling; users may hit panics or meaningless errors |
| MSSQL | L2 Strict | MySQL already found extensive "too lenient" issues; MSSQL very likely has the same (blocked by L1 cleanup) |
| PG | L2 Strict | Largest engine with no systematic strictness testing |
| Oracle | L6 Completion | Only major engine without a completion engine |

### Tier 3: Verification (Global Regression Safety Net)

| Engine | Layer | Rationale |
|--------|-------|-----------|
| PG | L5 Corpus | Largest engine but no structured corpus; highest regression risk |
| MSSQL | L5 Corpus | Multi-layer corpus verification covering dimensions beyond compare_test.go |
| MSSQL | L4 Coverage | 193 test functions exist but coverage scope not quantified |
| Oracle | L2 Strict | Complex syntax; likely has significant strictness gaps |

---

## Design Principles

1. **Follow the dependency chain**: Infrastructure first (Loc, Keyword), then functional layers (Soft-Fail, Strict, Completion), then verification (Corpus)
2. **Golden-standard driven**: Use real databases as oracle verification wherever possible
3. **Zero-difference target**: Golden tests (keyword/corpus/strict) target zero discrepancies
4. **Incremental progress**: Each engine progresses via phased SCENARIOS files; each phase independently verifiable
5. **Depth over breadth**: One engine at L1-L5 full coverage is more valuable than five engines each at L1 only
