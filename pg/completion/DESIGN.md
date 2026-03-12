# PG Completion Engine: Parse Table Based

## Goal

Replace the ANTLR C3-based completion for PostgreSQL with a completion engine built on pgparser's goyacc parse tables. Grammar-level completeness (100% keyword coverage). No regression on existing 21 test cases.

## Architecture

```
Input: SQL text + caret position
  │
  ├─ Preprocessor
  │   ├─ splitStatements(): split by semicolon, keep caret statement
  │   └─ findLatestStatement(): fallback for incomplete SQL
  │
  ├─ pgparser Lexer: tokenize the statement
  │
  ├─ Parse Table Engine (replaces C3)
  │   ├─ simulateParse(): run tokens through LALR state machine
  │   ├─ collectValidTokens(): expand reduce chains, collect shiftable tokens
  │   └─ inferContexts(): determine grammar context (columnref/relation_expr/...)
  │   Output: CompletionHint{ValidTokens, Contexts}
  │
  ├─ Semantic Expander
  │   ├─ ValidTokens → keyword candidates
  │   ├─ InColumnRef → collect table refs, query metadata for columns
  │   ├─ InRelationExpr → query metadata for tables/views/sequences
  │   ├─ InFuncName → built-in function list
  │   ├─ InTypeName → type list (for DDL)
  │   └─ InQualifiedName → schema-qualified object resolution
  │
  └─ Output: []Candidate
```

## Components

### 1. Parse Table Access (`tables.go`)

pgparser's generated parse tables are unexported (lowercase `pg` prefix). Add exported accessors in `parser/` package:

```go
// parser/tables_export.go
func PactTable() []int32  { ... }
func ActTable() []int16   { ... }
// etc.
```

The completion package imports `github.com/pgplex/pgparser/parser` to access these.

### 2. Parse Simulator (`simulator.go`)

Replicate the goyacc parse loop (`parser/parser.go` lines 17888-18150) without semantic actions:
- Shift: push new state
- Reduce: pop RHS, lookup goto, push new state
- Error: skip token, continue (completion-mode recovery)

No AST construction. Pure state stack tracking.

### 3. Token Collector (`collector.go`)

From a state stack, compute all valid next tokens:
1. Current state: find shiftable tokens via `Pact[state] + tok` / `Chk` / `Act`
2. Default reduce (`Def[state]`): simulate reduce, recurse on new state
3. Exception table (`Exca`): check for additional reduces
4. Depth limit: 20 (reduce chains are typically 3-5 deep)

Reference: `pgErrorMessage()` in `parser/parser.go` already implements expected-token extraction for error messages.

### 4. Context Inferrer (`context.go`)

Determine which grammar rule the caret is inside. Two methods:

**Primary — "shift IDENT, trace reduce":**
1. Clone state stack
2. Simulate shifting a generic IDENT token
3. Follow reduce chain: `pgR1[prod]` = LHS nonterminal
4. If nonterminal is a preferred rule → record as context

**Fallback — state-to-rule mapping:**
Pre-built from `y.output`: `map[int][]string` mapping state → active grammar rules.

Preferred rules (matching C3 config):
- `columnref` — column references
- `relation_expr` — table/view references
- `qualified_name` — schema.object references
- `func_name` — function names
- `typename` — type names (DDL extension)

### 5. Preprocessor (`preprocess.go`)

Isolate the caret statement using pgparser's lexer:
1. Tokenize full input, track semicolons
2. Find statement containing caret, adjust position
3. Fallback: find latest statement-starting keyword at column 0

### 6. Semantic Expander (`expander.go`)

Convert CompletionHint into candidates using metadata. Migrated from bytebase's `pg/completion.go`:
- Table reference collection (FROM/JOIN) — rewritten to use pgparser AST
- CTE extraction — rewritten using pgparser AST
- Column/table/function/type expansion — reuse logic with metadata interface
- Alias resolution, schema qualification, quoting

### 7. State-to-Rule Mapping (`cmd/gen_rulemap/`)

Codegen tool: parse `y.output` (582 MB) → generate `rulemap_generated.go`.

For each state, extract which nonterminals have items with dot at relevant positions. Output a static `map[int][]string`.

## Error Recovery

**Before parse simulation (preprocessor):**
- Multiple statements: split by semicolon, keep caret statement
- Invalid preceding statement: find latest keyword at line start

**During parse simulation:**
- Token mismatch: skip token, keep state stack
- Consecutive errors (3+): pop stack to find accepting state

**Not needed:**
- Recovery after caret — we never parse past the caret

## DDL/DML Extension

Parse tables already contain all DDL/DML grammar. Extension requires:
1. Add `typename` to preferred rules
2. Add semantic expanders for: data types, constraints, indexes
3. Add test cases for CREATE/ALTER/UPDATE/INSERT/DELETE

## File Layout

```
pgparser/
├── completion/
│   ├── DESIGN.md              # this file
│   ├── tables.go              # parse table accessor wrappers
│   ├── simulator.go           # LALR state machine simulator
│   ├── collector.go           # valid token collection
│   ├── context.go             # grammar context inference
│   ├── preprocess.go          # statement isolation
│   ├── expander.go            # semantic expansion
│   ├── completion.go          # entry point
│   ├── completion_test.go     # tests
│   ├── testdata/
│   │   └── completion.yaml    # test cases
│   └── cmd/
│       ├── gen_rulemap/       # y.output → rulemap_generated.go
│       └── diff_c3/           # C3 vs new engine comparison
├── parser/
│   ├── tables_export.go       # NEW: exported parse table accessors
│   └── ...                    # existing parser code
└── ...
```

## Milestones

| # | Milestone | Verification |
|---|-----------|-------------|
| M1 | Parse table export | `go build ./parser/` passes |
| M2 | Parse simulator | Unit: known SQL → expected state stack |
| M3 | Token collector | Diff: 21 cases, tokens ⊇ C3 tokens |
| M4 | Context inferrer | Unit: known positions → expected contexts |
| M5 | Preprocessor | Unit: multi-statement → correct isolation |
| M6 | Semantic expander | Integration: 21 cases, full candidate match |
| M7 | Full pipeline | All 21 existing tests pass exactly |
| M8 | DDL/DML extension | New tests for CREATE/ALTER/UPDATE/INSERT |
| M9 | bytebase integration | bytebase CI completion tests pass |
