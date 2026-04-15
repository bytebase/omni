# pg/parser FIRST-set audit

**Generated:** 2026-04-14
**Branch:** pg-first-sets
**Phase:** 3.1 audit

## Method

Three grep variants over `pg/parser/*.go` (excluding `_test.go`), merged and
de-duplicated:

```
grep -nE "case ([A-Z_][A-Z_0-9]*|'[^']+')(,\s*([A-Z_][A-Z_0-9]*|'[^']+'|[a-z_][A-Za-z_0-9]*))+\s*:" pg/parser/*.go
grep -nE "case [^:]+,[^:]+:" pg/parser/*.go
```

66 unique cluster lines were examined. Each was classified by reading its
enclosing function in the source file.

## Classification

- **PROBE**: The cluster is used as a "could the next token start production
  X?" check. All branches share the body or the cluster is wrapped inside a
  predicate helper that already returns `bool`. Migration candidate for
  Phase 3 (replace with an `isXStart()` call).
- **DISPATCH**: The case body implements distinct per-token logic. Removing
  the cluster would break the parser — it IS the implementation. This
  includes "merged alias" branches where a small keyword family shares one
  body inside a larger statement-level switch (e.g. parser.go parseStmt).
- **OTHER**: Operator categorization, lexer state machines, token
  reclassification (`NOT_LA`, `WITH_LA`, …), escape-sequence handlers,
  code-generation lookups. Not FIRST-set checks; left alone.

## Summary

- Total clusters found: **66**
- **PROBE: 11**
  - Covered by existing predicate: **5** (`isSelectStart` already exists)
  - Already wrapped inside a self-contained predicate helper: **4**
    (`isFuncExprCommonSubexprStart`, `isJsonBehaviorStart`,
    `isReservedForClause`, `isJoinKeyword`, `isCreateTableElement`,
    `isSelectStart`) — these are the consolidated state we want, no work
  - Needs NEW predicate: **2** (`isTableConstraintStart` covering three
    call sites; `isTransactionStmtStart` — low-value judgment call)
- **DISPATCH: 37**
- **OTHER: 18**

## Probe migration candidates

These are clusters whose body either (a) delegates to a single parser
function and could be rewritten as `if p.isXStart() { return p.parseX() }`,
or (b) is already wrapped in an `isXStart` helper whose body is the cluster.

### (a) Rewrites — replace inline `case` with predicate call

| File:Line | Function | Tokens | Existing predicate | Notes |
|---|---|---|---|---|
| copy.go:84 | parsePreparableStmt | SELECT, VALUES, TABLE, WITH | `isSelectStart` | Case body: `return p.parseSelectNoParens()`. Straight substitution; also functions as the `default` fallback, so rewrite as leading `if`. |
| publication.go:685 | parseRuleActionStmt | SELECT, VALUES, TABLE, WITH | `isSelectStart` | Case body has an inner `if p.cur.Type == WITH` that dispatches to `parseWithStmt`. Migration: `if p.isSelectStart() { … }` with the WITH branch preserved inside. |
| create_table.go:250 | parseTableElement | CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE | **NEW** `isTableConstraintStart` | Body: `return p.parseTableConstraint(), nil`. Pair with the neighboring `case CONSTRAINT` above it (the predicate must include CONSTRAINT). |
| create_table.go:1043 | parseTypedTableElement | CONSTRAINT, CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE | **NEW** `isTableConstraintStart` | Body: `return p.parseTableConstraint()`. Exact TableConstraint FIRST set. |
| alter_table.go:731 | parseAlterTableAdd | CONSTRAINT, CHECK, UNIQUE, PRIMARY, EXCLUDE, FOREIGN | **NEW** `isTableConstraintStart` | Body: `return AlterTableCmd{AT_AddConstraint, parseTableConstraint}`. Same FIRST set reordered. |
| parser.go:225 | parseStmt | BEGIN_P, START, COMMIT, END_P, ABORT_P, SAVEPOINT, RELEASE | **NEW** `isTransactionStmtStart` (judgment) | Body: `return p.parseTransactionStmt()`. Neighboring `case ROLLBACK` is the seventh member; `case PREPARE` has extra `peekNext()==TRANSACTION` logic so it stays out. Low-value: only one call site, but the predicate would make the FIRST set discoverable. |

### (b) Already-consolidated predicates (no migration work)

These are listed for completeness — they are the state we want for all PROBE
clusters. Phase 3.2 should simply verify none of them have been accidentally
duplicated inline elsewhere.

| File:Line | Helper | Production |
|---|---|---|
| expr.go:1222 | `isFuncExprCommonSubexprStart` | func_expr_common_subexpr |
| expr.go:3147 | `isSelectStart` | SelectStmt FIRST |
| expr.go:3242 | `isReservedForClause` | reserved-clause-keyword set |
| expr.go:3257 | `isJoinKeyword` | join_type / NATURAL / JOIN / ON |
| json.go:153 | `isJsonBehaviorStart` | json_behavior |
| parser.go:953 | `isCreateTableElement` | TableConstraint ∪ {LIKE} |

## New predicate candidates

| File:Line | Function | Production | Suggested predicate name | FIRST set |
|---|---|---|---|---|
| create_table.go:250, 1043; alter_table.go:731 | parseTableElement / parseTypedTableElement / parseAlterTableAdd | TableConstraint | `isTableConstraintStart` | CONSTRAINT, CHECK, UNIQUE, PRIMARY, FOREIGN, EXCLUDE |
| parser.go:225 (judgment) | parseStmt | TransactionStmt | `isTransactionStmtStart` | BEGIN_P, START, COMMIT, END_P, ABORT_P, SAVEPOINT, RELEASE, ROLLBACK |

Notes:

- `isTableConstraintStart` is the highest-value new predicate: three
  independent call sites, and the set is already duplicated in
  `isCreateTableElement` (which is TableConstraint + LIKE). Phase 3.2 should
  either have `isCreateTableElement` delegate to `isTableConstraintStart`,
  or define `isTableConstraintStart` as a slice and reuse it.
- `isTransactionStmtStart` is single-use, inside the top-level `parseStmt`
  switch. Adding it is stylistic; it makes the FIRST set greppable and
  oracle-testable, but doesn't remove duplication. Gate on user preference.
- No oracle-level new FIRST sets are required — all PROBE candidates here
  are either already covered or map to a single named PG production.

## Dispatch switches (no action)

Each row below is a cluster whose body contains distinct per-branch logic or
whose enclosing switch IS the canonical implementation of its production.
Removing the cluster would break parsing; these stay as-is.

| File:Line | Function | Reason |
|---|---|---|
| alter_misc.go:84 | parseAlterFunctionStmt | NO/DEPENDS alias inside an ALTER FUNCTION sub-switch |
| create_table.go:108 | parseOptTemp | TEMPORARY/TEMP alias |
| cursor.go:305 | parseFetchDirection | FROM/IN_P alias — from_in production |
| cursor.go:318 | parseFetchDirection | ICONST/'+'/'-' — SignedIconst FIRST (trivial 3-token inline; not worth a predicate) |
| drop.go:49 | parseDropStmtInner | object_type_any_name dispatch, each branch is a distinct drop path |
| drop.go:104 | parseDropStmtInner | ROLE/GROUP_P alias |
| expr.go:202 | parseAExprInfix | `<`/`>`/`=` comparison operator, body parses subquery-op or binary op |
| expr.go:349 | parseAExprInfix | `+`/`-` infix arithmetic |
| expr.go:365 | parseAExprInfix | `*`/`/`/'%' infix arithmetic |
| expr.go:507 | parseIsPostfix | NFC/NFD/NFKC/NFKD — `IS [NOT] NORMALIZED` form, distinct form arg per token |
| expr.go:1121 | parseBExprInfix | `<`/`>`/`=` in b_expr |
| expr.go:1193 | parseBExprInfix | `+`/`-`/`*`/`/`/'%'/'^' merged infix arithmetic |
| grant.go:769 | parseCreateOptRoleElem | ROLE/USER alias |
| maintenance.go:304 | parseReindexStmt | INDEX/TABLE — reindex_target_type, kind switched by token |
| maintenance.go:347 | parseReindexStmt | SYSTEM_P/DATABASE — reindex_target_multitable |
| maintenance.go:404 | parseUtilityOptionName | ANALYZE/ANALYSE alias |
| maintenance.go:435 | parseUtilityOptionArg | ICONST/FCONST numeric literal branch |
| maintenance.go:437 | parseUtilityOptionArg | `+`/`-` signed numeric branch |
| parser.go:108 | parseStmt | SELECT/VALUES/TABLE — top-level stmt dispatch |
| parser.go:184 | parseStmt (ALTER sub-switch) | FUNCTION/PROCEDURE/ROUTINE — alias for parseAlterFunctionStmt |
| parser.go:285 | parseStmt | FETCH/MOVE alias for parseFetchStmt |
| parser.go:291 | parseStmt | ANALYZE/ANALYSE alias |
| parser.go:379 | parseCreateDispatch | FUNCTION/PROCEDURE alias |
| parser.go:381 | parseCreateDispatch | TRIGGER/CONSTRAINT — CREATE OR REPLACE TRIGGER / CONSTRAINT TRIGGER |
| parser.go:385 | parseCreateDispatch | TRUSTED/PROCEDURAL/LANGUAGE — CREATE [TRUSTED \| PROCEDURAL] LANGUAGE |
| parser.go:424 | parseCreateDispatch | TEMPORARY/TEMP alias |
| parser.go:955 | isCreateTableElement | Already a consolidated predicate (included here because the grep matched the cluster inside its body) |
| schema.go:95 | parseSchemaStmt | INDEX/UNIQUE — CREATE [UNIQUE] INDEX |
| select.go:974 | parseOptTempTableName | TEMPORARY/TEMP alias |
| select.go:978 | parseOptTempTableName | LOCAL/GLOBAL — two prefixes with same following-body |
| set.go:263 | parseVariableSetStmt | TO/'=' alias |
| type.go:113 | parseSimpleTypename | INT_P/INTEGER alias — canonical SimpleTypename implementation |
| utility.go:72 | parseExplainableStmt | SELECT/VALUES/TABLE/'(' — includes `'('` for parenthesized select, distinct from isSelectStart |
| grant.go:400 | parseGrantFuncArg | IN_P/OUT_P/INOUT/VARIADIC — arg_class optional prefix; 4-token inline, no reuse, left alone |
| parser.go:3149-ish via expr.go:3149 | `isSelectStart` body | The cluster IS the predicate body; counted under "already consolidated" above |

## Other clusters (no action)

| File:Line | Function | Category | Reason |
|---|---|---|---|
| complete.go:21 | TokenName | Token-name lookup | Non-keyword token sentinel, returns `""` |
| define.go:320 | parseAnyOperator | Operator categorization | Math-op set routed to parseMathOp |
| expr.go:72 | aExprInfixPrec | Operator precedence table | Precedence lookup, not a FIRST set |
| expr.go:74 | aExprInfixPrec | Operator precedence | `<`/`>`/`=` precComparison |
| expr.go:76 | aExprInfixPrec | Operator precedence | LESS_EQUALS etc. precComparison |
| expr.go:78 | aExprInfixPrec | Operator precedence | BETWEEN/IN/LIKE/ILIKE/SIMILAR precIn |
| expr.go:84 | aExprInfixPrec | Operator precedence | Op/OPERATOR precOp |
| expr.go:86 | aExprInfixPrec | Operator precedence | +/- precAdd |
| expr.go:88 | aExprInfixPrec | Operator precedence | `*`/`/`/'%' precMul |
| expr.go:1058 | bExprInfixPrec | Operator precedence | `<`/`>`/`=` |
| expr.go:1060 | bExprInfixPrec | Operator precedence | LESS_EQUALS etc. |
| expr.go:1062 | bExprInfixPrec | Operator precedence | Op/OPERATOR |
| expr.go:1064 | bExprInfixPrec | Operator precedence | +/- |
| expr.go:1066 | bExprInfixPrec | Operator precedence | `*`/`/`/'%' |
| lexer.go:161 | Lexer.NextToken | Lexer state machine | stateXB/stateXH dispatch |
| lexer.go:163 | Lexer.NextToken | Lexer state machine | stateXQ/stateXE/stateXUS dispatch |
| lexer.go:169 | Lexer.NextToken | Lexer state machine | stateXD/stateXUI dispatch |
| lexer.go:498 | lexEscape | Character literal | Octal digit escape |
| lexer.go:613 | lexQuoteContinue | Lexer token-typing | Returns the right literal token type based on entry state |
| name.go:602 | parseMathOp | Operator categorization | Math-op set → string |
| parser.go:1115 | Parser.advance | Lookahead reclassification | NOT → NOT_LA when followed by BETWEEN/IN/LIKE/ILIKE/SIMILAR |
| parser.go:1122 | Parser.advance | Lookahead reclassification | WITH → WITH_LA when followed by TIME/ORDINALITY |
| parser.go:1134 | Parser.advance | Lookahead reclassification | NULLS → NULLS_LA when followed by FIRST/LAST |

## Notes for Phase 3.2

Of the **11** PROBE entries:

- **5** are straightforward substitutions requiring NO new predicate:
  - copy.go:84 → `isSelectStart()`
  - publication.go:685 → `isSelectStart()` (with inner WITH branch preserved)
  - create_table.go:250, create_table.go:1043, alter_table.go:731 → one NEW
    predicate `isTableConstraintStart` shared across all three (plus a
    refactor of `isCreateTableElement` to delegate).
- **1** is a low-value judgment call (parser.go:225 / `isTransactionStmtStart`).
  Single call site, no duplication. Add only if the team wants the FIRST set
  made explicit and oracle-tested.
- **5** are already-consolidated helpers (`isFuncExprCommonSubexprStart`,
  `isJsonBehaviorStart`, `isReservedForClause`, `isJoinKeyword`,
  `isCreateTableElement`, `isSelectStart`). No work — but Phase 3.2 should
  verify the helpers are the only definition and that no inline duplicate
  of the cluster has crept in elsewhere (greps already clean as of this
  audit).

Phase 3.2 scope is therefore **small**: one new predicate
(`isTableConstraintStart`), three call-site rewrites using it, two
call-site rewrites using `isSelectStart`, optionally
`isTransactionStmtStart` plus its single call site. All other clusters
in the grep output are DISPATCH or OTHER and must stay as-is.
