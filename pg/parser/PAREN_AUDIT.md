# PAREN_AUDIT — pg/parser `(` / `)` dispatch sites

**Status:** Phase 1 audit complete (2026-04-21); Phase 4 §5.3 lock-in proofs + §4.2 citation sweep complete.
**Scope:** every `(` / `)` dispatch site in `pg/parser/*.go`
**Reference:** `docs/plans/2026-04-21-pg-paren-dispatch.md` §4.1 scope, §3 technique catalogue, §5.3 alignment bar

**Note on site identifiers:** `file:line` in this document is the **original audit coordinate** captured at Phase 1 scope-lock (2026-04-21). Line numbers drift as code evolves; we retain the original coordinate as a stable row ID so audit history remains cross-referenceable. For the current live location of any site, grep the function name.
**Machine-readable mirror:** `PAREN_AUDIT.json` — **source of truth** for the full `proof_notes` text (caller-context paragraphs and test citations). The markdown cluster tables below carry a terse one-line technique/status summary; for the §5.3 two-proof bar text (caller-context + ≥5 pinned tests for ambiguity-present sites) or test-file citations (non-ambiguous sites), consult the matching site row in `PAREN_AUDIT.json`.

## Summary

- **94 total dispatch sites** (69 `(` opens, 25 `)` closes)
- **5 ambiguity-present sites:** all aligned (Phase 0 + Phase 1 §§1.1–1.4)
- **89 non-ambiguous sites:** aligned via grammar structure (unconditional `expect('(')` / `expect(')')` or T1 peek on optional list)
- **0 unclear** (was 1: `expr.go:1610` parseArraySubscript ARRAY subquery contract — locked in by §1.4 with 6 empirical tests + caller-context proof; see paren_array_expr_test.go)

## By cluster

| Cluster | File(s) | Sites | Ambiguity-present | Misaligned |
|---------|---------|-------|-------------------|------------|
| C1 | select.go | 14 | 3 | 0 |
| C2 | expr.go | 24 | 2 | 0 |
| C3 | type.go | 9 | 0 | 0 |
| C4 | create_table.go, create_index.go, define.go | 14 | 0 | 0 |
| C5 | (11 utility files — see §5.1 unbounded-tail note in PLAN) | 33 | 0 | 0 |

## Attention list

All previously attention-listed sites are resolved (Phase 1 BATCH 1 §§1.1–1.4). The table below is kept for historical reference.

| # | Site | Function | Priority | Ambiguity | Resolution |
|---|------|----------|----------|-----------|------------|
| 1 | `expr.go:2053` | `parseInExpr` | **high** | `IN (expr_list)` vs `IN (select)` | §1.1: T3 snapshot-scan + T8 FIRST-set; completion-mode guard on lookahead |
| 2 | `select.go:1347` | `parseLateralTableRef` | **high** | LATERAL variants | §1.2: T6/T7 hybrid — dedicated dispatch arms + post-parse ColumnRef reject |
| 3 | `expr.go:1609` | `parseArrayCExpr` | med | `ARRAY[...]` vs `ARRAY(select)` | §1.3: T1 peek on `[` vs `(` at entry |
| 4 | `expr.go:1610` | `parseArrayCExpr` | med | `ARRAY(select)` content contract | §1.4: T7 post-parse typed-nil guard rejects ARRAY() / ARRAY(1) / ARRAY(ROWS FROM …) |

## Cluster tables

### C1: table_ref / select_with_parens / joined_table

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| select.go:77 | parseSelectWithParens | select_with_parens | no | T5 | yes | high | Unconditional expect('('); no disambiguation |
| select.go:166 | parseSelectClausePrimary | select_with_parens | yes | T5 | yes | high | Left-factored to parseSelectNoParens |
| select.go:770 | parseCommonTableExpr | name_list | no | T1 | yes | med | Optional paren for column names |
| select.go:1162 | parseParenTableRef | select_with_parens, joined_table | yes | T3 | yes | high | Phase 0: parenBeginsSubquery scanner |
| select.go:1225 | parseJoinedTable | joined_table | no | T6 | yes | high | Phase 0: T6/T7 dedicated nonterminal |
| select.go:1347 | parseLateralTableRef | select_with_parens, xmltable, json_table, func_table | yes | T6+T7 | yes | high | §1.2 landed: four disjoint leads ('(' / XMLTABLE / JSON_TABLE / ROWS / default) routed via dedicated arms; post-parse ColumnRef rejection closes `LATERAL <relation>` gap. Caller-context: only reachable from parseTableRefPrimary's LATERAL arm. Tests: paren_lateral_ref_test.go TestParenLateralRefDispatch (12 scenarios: 10 accept, 2 reject) + AST-shape tests for subquery/func_table/rows-from/xmltable/json_table/alias-coldeflist/ordinality. |
| select.go:1504 | parseRelationExpr | func_table | no | T1 | yes | med | After identifier; '(' always func_table |
| select.go:1528 | parseRelationExpr | func_table | no | T1 | yes | med | After schema.attr; '(' always func_table |
| select.go:1552 | parseRelationExpr | func_table | no | T1 | yes | med | After catalog.schema.attr; '(' always func_table |
| select.go:1674 | parseOptSearchClause | columnList | no | none | yes | low | Simple column list in SEARCH clause |
| select.go:1767 | parseOptCycleClause | CycleClause | no | none | yes | low | Unconditional expect |
| select.go:1782 | parseOptCycleClause | CycleClause | no | none | yes | low | Unconditional expect |
| select.go:1821 | parseOptCycleClause | CycleClause | no | none | yes | low | Unconditional expect |
| select.go:1834 | parseOptCycleClause | CycleClause | no | none | yes | low | Unconditional expect |

### C2: expression parens (a_expr, row, subquery, function calls)

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| expr.go:1609 | parseArrayCExpr | A_ArrayExpr, ARRAY_SUBLINK | yes | T1 | yes | med | §1.3 landed: T1 peek on '[' vs '(' after ARRAY routes to parseArrayExpr (array constructor) vs parseSelectStmtForExpr (sublink). Caller-context: parseArrayCExpr is the sole consumer of the ARRAY keyword lead (dispatched from c_expr at expr.go:1491). Tests: paren_array_expr_test.go TestParenArrayExprDispatch (10 scenarios: 6 constructor '[' forms, 4 sublink '(' forms), TestParenArrayExprConstructorAST (A_ArrayExpr shape), TestParenArrayExprNestedConstructorAST, TestParenArrayExprSublinkAST (ARRAY_SUBLINK shape). Ref: gram.y:15440-15459 + 16583-16595. |
| expr.go:1610 | parseArrayCExpr | ARRAY_SUBLINK | no | T1+T7 | yes | med | §1.4 landed: ARRAY(...) content contract via T7 post-parse typed-nil guard. Caller-context: parseArrayCExpr is invoked only after the ARRAY keyword (c_expr dispatch at expr.go:1491), and the '(' branch here is entered only after T1 peek confirms '(' (not '[', not other). The enclosing parseArrayCExpr calls parseSelectStmtForExpr (= parseSelectNoParens) which returns a typed-nil *SelectStmt for anything not in the select_no_parens FIRST set; we type-assert and reject. Pinned scenarios (paren_array_expr_test.go TestParenArrayExprSubqueryContract at line 151): accept — SELECT ARRAY(TABLE foo) (line 162-165), SELECT ARRAY(WITH cte AS (SELECT 1) SELECT * FROM cte) (line 169-173), SELECT ARRAY(VALUES (1), (2)) (line 175-179); reject — SELECT ARRAY(ROWS FROM (f(1))) (line 186-190), SELECT ARRAY(1) (line 194-198), SELECT ARRAY() (line 202-206). Additional AST-shape proof: TestParenArrayExprSublinkAST (line 224) asserts SubLink{SubLinkType=ARRAY_SUBLINK}. Ref: gram.y:15440-15451. |
| expr.go:1741 | parseExplicitRow | explicit_row | no | none | yes | low | ROW keyword disambiguates |
| expr.go:1747 | parseExplicitRow | explicit_row | no | none | yes | low | Empty row check |
| expr.go:1773 | parseCastExpr | cast_expr | no | none | yes | low | After CAST; unconditional |
| expr.go:1787 | parseCastExpr | cast_expr | no | none | yes | low | Expect AS then closing |
| expr.go:1796 | parseCoalesceExpr | coalesce_expr | no | none | yes | low | After COALESCE; unconditional |
| expr.go:1801 | parseCoalesceExpr | coalesce_expr | no | none | yes | low | Expect closing paren |
| expr.go:1812 | parseMinMaxExpr | greatest_expr / least_expr | no | none | yes | low | After GREATEST/LEAST; unconditional |
| expr.go:1818 | parseMinMaxExpr | greatest_expr / least_expr | no | none | yes | low | Expect closing paren |
| expr.go:1828 | parseNullIfExpr | nullif_expr | no | none | yes | low | After NULLIF; unconditional |
| expr.go:1834 | parseNullIfExpr | nullif_expr | no | none | yes | low | Expect comma then expr |
| expr.go:1840 | parseNullIfExpr | nullif_expr | no | none | yes | low | Expect closing paren |
| expr.go:1855 | parseExtractExpr | extract_expr | no | none | yes | low | After EXTRACT; unconditional |
| expr.go:1862 | parseExtractExpr | extract_expr | no | none | yes | low | Check for empty extract parens |
| expr.go:1878 | parseExtractExpr | extract_expr | no | none | yes | low | Expect closing paren |
| expr.go:2053 | parseInExpr | in_expr | yes | T3+T8 | yes | high | §1.1 landed: fast-path T8 (isSelectStart) for SELECT/VALUES/WITH/TABLE; T3 snapshot/restore lookaheadInIsSubquery scan for the `IN ((...))` ambiguous case (top-level ',' → expr_list, top-level UNION/INTERSECT/EXCEPT → subquery, single-element → FIRST-set probe past leading '('). Completion-mode guard added at lookaheadInIsSubquery entry (snapshotTokenStream does not cover collecting/candidates). Caller-context: only entered after IN token consumed and '(' confirmed; empty `IN ()` rejected. Tests: paren_in_expr_test.go TestParenInExprDispatch (20 scenarios spanning expr_list, subquery, nested parens, VALUES/WITH/TABLE leads, empty list reject, row-constructor LHS, NOT IN variants, partition_prune regression) + AST-shape tests TestParenInExprSubqueryAST (ANY_SUBLINK), TestParenInExprExprListAST (A_Expr AEXPR_IN with "="), TestParenInExprNotInSubqueryAST (BoolExpr NOT wrap), TestParenInExprNotInExprListAST (operator "<>"). Ref: gram.y:14973-14998. |
| expr.go:2112 | parseSubqueryExprList | IN_SUBLINK | no | none | yes | med | After IN; unconditional |
| expr.go:2404 | parseSVFWithOptionalPrecision | sql_value_function | no | T1 | yes | med | Optional precision; T1 peek sufficient |
| expr.go:2452 | parseColumnRefOrFuncCall | func_application | no | T1 | yes | med | After column ID; '(' always func_call |
| expr.go:2494 | parseColumnRefOrFuncCall | func_application | no | T1 | yes | med | After schema.attr; '(' always func_call |
| expr.go:2577 | parseFuncApplication | func_application | no | none | yes | low | After '('; check for empty/*/DISTINCT/ALL |
| expr.go:2797 | parseArraySubscript | array_subscript | no | none | yes | low | Inside brackets; array subscript/slice |

### C3: type parens (SimpleTypename, numeric, character, timestamp, interval)

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| type.go:245 | parseOptTypeModifiers | opt_type_modifiers | no | T1 | yes | low | Optional typmod list |
| type.go:313 | parseOptFloat | opt_float | no | T1 | yes | low | After FLOAT; optional precision |
| type.go:340 | parseBit | bit_type | no | T1 | yes | low | After BIT; optional length |
| type.go:395 | parseVarcharType | varchar_type | no | T1 | yes | low | After VARCHAR; optional length |
| type.go:440 | finishCharType | character_type | no | T1 | yes | low | After CHARACTER/CHAR; optional length |
| type.go:463 | parseTimestampType | timestamp_type | no | T1 | yes | low | After TIMESTAMP; optional precision |
| type.go:496 | parseTimeType | time_type | no | T1 | yes | low | After TIME; optional precision |
| type.go:553 | parseIntervalType | interval_type | no | T1 | yes | low | After INTERVAL; optional precision |
| type.go:687 | parseIntervalSecond | interval_second | no | T1 | yes | low | After SECOND_P; optional fractional precision |

### C4: DDL element lists (table elements, index elements, constraints)

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| create_table.go:136 | parseTableElement | constraint_elem | no | T1 | yes | low | After constraint keyword; paren-wrapped |
| create_table.go:176 | parseTableElement | constraint_elem | no | T1 | yes | low | Constraint element list |
| create_table.go:211 | parseTableElement | constraint_elem | no | none | yes | low | Unconditional closing paren |
| create_table.go:1475 | parsePartitionSpec | partition_spec | no | T1 | yes | low | Optional BY or parens |
| create_index.go:169 | parseIndexElement | index_elem | no | T1 | yes | low | Optional paren-wrapped expression |
| create_index.go:185 | parseIndexElement | index_elem | no | none | yes | low | Unconditional closing paren |
| create_index.go:258 | parseIndexWithClause | with_clause | no | T1 | yes | low | WITH (...) parameter list |
| define.go:32 | parseDefineArgList | arg_list | no | T1 | yes | low | Optional paren-wrapped argument list |
| define.go:551 | parseDefineArgListItem | arg_item | no | T1 | yes | low | Optional paren for type/expression |
| define.go:580 | parseDefineArgListItem | arg_item | no | T1 | yes | low | Closing paren |
| define.go:856 | parseDefineArgNameAndType | define_stmt | no | T1 | yes | low | Optional paren-wrapped value |
| define.go:980 | parseDefineValue | define_value | no | none | yes | low | Unconditional closing paren |
| define.go:1010 | parseDefineValue | define_value | no | T1 | yes | low | Optional paren-wrapped value |
| define.go:1012 | parseDefineValue | define_value | no | none | yes | low | Unconditional closing paren |

### C5: Utility statements (short tail)

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| database.go:360 | parseWithClauseOptions | option_list | no | T1 | yes | low | WITH (...) options in CREATE DATABASE |
| copy.go:18 | parseCopyFromTo | copy_from_to | no | T1 | yes | low | Optional FROM/TO paren syntax |
| copy.go:142 | parseCopyOptions | copy_option_list | no | T1 | yes | low | Optional paren for copy option |
| grant.go:367 | parseGrantRoleList | role_list | no | T1 | yes | low | Optional paren for role list |
| grant.go:369 | parseGrantRoleList | role_list | no | none | yes | low | Unconditional closing paren |
| extension.go:274 | parseExtensionOptionList | option_list | no | T1 | yes | low | Optional WITH (...) options |
| extension.go:276 | parseExtensionOptionList | option_list | no | none | yes | low | Unconditional closing paren |
| extension.go:334 | parseExtensionOptions | option_item | no | none | yes | low | Unconditional closing paren |
| fdw.go:312 | parseForeignServerForOptions | foreign_server_for | no | T1 | yes | low | Optional FOR SERVER/USER parens |
| json.go:556 | parseJsonWrapper | json_wrapper | no | none | yes | low | Unconditional closing paren |
| json.go:697 | parseJsonReturning | json_returning | no | none | yes | low | Unconditional closing paren |
| trigger.go:207 | parseTriggerActionTime | trigger_action_time | no | none | yes | low | Check for closing paren |
| prepare.go:24 | parsePrepareStmt | prepare_stmt | no | T1 | yes | low | Optional paren for parameter types |
| maintenance.go:17 | parseAnalyzeStmt | analyze_stmt | no | T1 | yes | low | Optional paren for table list |
| maintenance.go:78 | parseVacuumRelList | vacuum_rel_list | no | T1 | yes | low | Optional paren for relation list |
| maintenance.go:173 | parseVacuumOptions | vacuum_option_list | no | T1 | yes | low | Optional paren for option list |
| maintenance.go:295 | parseReindexIndexOptions | reindex_option_list | no | T1 | yes | low | Optional paren for reindex options |
| schema.go:638 | parseGrantOnClause | privilege_list | no | T1 | yes | low | Optional paren for privilege list |
| schema.go:640 | parseGrantOnClause | privilege_list | no | none | yes | low | Unconditional closing paren |
| schema.go:675 | parseGrantStmt | grantee_list | no | T1 | yes | low | Optional paren for grantee list |
| schema.go:681 | parseGrantStmt | grantee_list | no | none | yes | low | Unconditional closing paren or branch |
| schema.go:697 | parseGrantStmt | grantee_list | no | T1 | yes | low | Paren for optional GROUP/ROLE syntax |
| publication.go:638 | parsePublicationForClause | publication_for_list | no | T1 | yes | low | Optional paren for FOR list |
| publication.go:670 | parsePublicationForClause | publication_for_list | no | none | yes | low | Unconditional closing paren or terminator |
| insert.go:207 | parseInsertSelectClause | insert_select_clause | no | T1 | yes | low | Optional paren for column list |
| insert.go:399 | parseInsertColumnsClause | column_list | no | T1 | yes | low | Paren-wrapped column list |
| create_function.go:142 | parseParamsList | param_list | no | none | yes | low | Unconditional closing paren |
| utility.go:11 | parsePreparableStmt | select_with_parens | no | T1 | yes | low | Optional paren-wrapped select |
| parser.go:659 | parseCreateDispatch | create_stmt | no | T1 | yes | low | After CREATE [OR REPLACE]; optional parens |
| parser.go:730 | parseCreateDispatch | temp_table | no | T1 | yes | low | Optional temp table clause parens |
| parser.go:778 | parseCreateDispatch | create_table_as | no | T1 | yes | low | CREATE TABLE AS with options parens |
| parser.go:847 | parseCreateDispatch | with_clause | no | none | yes | low | Unconditional closing paren |
| parser.go:897 | parseCreateDispatch | column_list | no | none | yes | low | Check for closing paren or comma |

---

## Historical attention list (pre-Phase-1)

Originally surfaced as the top misaligned sites. All four have since been closed by Phase 1 of the starmap — the rows in the main cluster tables above reflect current state (`aligned: yes`). Kept here for audit-trail continuity; do not use for future prioritization.

1. ~~**HIGH** expr.go:2053 parseInExpr — IN (list) vs IN (SELECT)~~ **CLOSED in §1.1** (commit ad700fa). Technique: T3+T8 hybrid (snapshot scan + FIRST-set). Proof: `TestParenInExprDispatch` in `paren_in_expr_test.go`.
2. ~~**HIGH** select.go:1347 parseLateralTableRef — LATERAL dispatch~~ **CLOSED in §1.2** (commit 284f39e). Technique: T6/T7 hybrid (dedicated nonterminal + post-parse ColumnRef reject). Proof: `TestParenLateralRef*` in `paren_lateral_ref_test.go`.
3. ~~**HIGH** expr.go:1609 parseArrayCExpr — ARRAY[ vs ARRAY(~~ **CLOSED in §1.3** (commit be8af80). Technique: T1 peek. Proof: `TestParenArrayExprDispatch` in `paren_array_expr_test.go`.
4. ~~**MED** expr.go:1610 parseArrayCExpr — ARRAY( content contract~~ **CLOSED in §1.4** (commit c1158b7 + codex-fix b97477f). Technique: T1 + T7 (typed-nil guard). Proof: `TestParenArrayExprSubqueryContract` in `paren_array_expr_test.go`.

Remaining low-priority items from the original list (already informally aligned via T1 peek; kept for §5.3 lock-in in Phase 4):

5. select.go:770 parseCommonTableExpr — optional paren-wrapped column list (T1; Phase 4 §4.2 citation: pgregress `with.sql`, `with_recursive.sql`)
6. type.go:313-687 parseOptFloat / parseBit / parseVarcharType / … — type-modifier paren family (T1; Phase 4 §4.2 citations: pgregress `float4.sql`, `bit.sql`, `varchar.sql`, `timestamp.sql`, `interval.sql`, etc.)

---

## Phase 4 §5.3 lock-in summary (2026-04-21)

- **§4.1 two-proof bar** applied to the 4 Phase 0 `select_with_parens` / `joined_table` sites: `select.go:77 parseSelectWithParens`, `select.go:166 parseSelectClausePrimary`, `select.go:1162 parseParenTableRef`, `select.go:1225 parseJoinedTable`. Each row in `PAREN_AUDIT.json` now carries (a) an explicit caller-context paragraph identifying every caller, (b) ≥5 pinned empirical tests from `paren_multi_join_test.go` (plus the oracle + pgregress anti-regression fences). The two Phase 1 ambiguity-present sites (`expr.go:1609-1610 parseArrayCExpr`, `expr.go:2053 parseInExpr`, `select.go:1347 parseLateralTableRef`) already carry the §5.3 bar from Phase 1 commits ad700fa / 284f39e / be8af80 / c1158b7.
- **§4.2 citation sweep** applied to all 85+ non-ambiguous C3/C4/C5 rows: every `proof_notes` field now cites at least one concrete test file (dedicated `pg/parser/*_test.go`) or pgregress SQL corpus entry (`pg/pgregress/testdata/sql/*.sql`) that exercises the grammar point. Grammar structure remains the primary proof; citations are for anti-regression traceability only.
- **Gaps surfaced:** none blocking. All 94 rows have at least one citation. A handful of rows cite only the broad pgregress corpus (no dedicated unit test) — these are low-priority utility sites (e.g. `trigger.go:207`, `extension.go:274-334`) where adding a bespoke paren test would be redundant with pgregress; deferred as minor TODOs rather than §4.2 blockers.


