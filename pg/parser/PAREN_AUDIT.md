# PAREN_AUDIT — pg/parser `(` / `)` dispatch sites

**Status:** Phase 1 audit complete (2026-04-21)
**Scope:** every `(` / `)` dispatch site in `pg/parser/*.go`
**Reference:** `/Users/rebeliceyang/Github/omni-byt-9315/docs/plans/2026-04-21-pg-paren-dispatch.md` §4.1 scope, §3 technique catalogue, §5.3 alignment bar
**Machine-readable mirror:** `PAREN_AUDIT.json`

## Summary

- **94 total dispatch sites** (69 `(` opens, 25 `)` closes)
- **5 ambiguity-present sites:** 2 already aligned (Phase 0), 3 not aligned → Phase 2 fix targets
- **88 non-ambiguous sites:** aligned via grammar structure (unconditional `expect('(')` / `expect(')')` or T1 peek on optional list)
- **1 unclear:** `expr.go:1610` parseArraySubscript (ARRAY subquery contract) — needs oracle cross-check, not a known bug

## By cluster

| Cluster | File(s) | Sites | Ambiguity-present | Misaligned |
|---------|---------|-------|-------------------|------------|
| C1 | select.go | 14 | 3 | 1 (select.go:1347 parseLateralTableRef) |
| C2 | expr.go | 24 | 2 | 2 (expr.go:2053 parseInExpr, expr.go:1609 parseArraySubscript) |
| C3 | type.go | 9 | 0 | 0 |
| C4 | create_table.go, create_index.go, define.go | 14 | 0 | 0 |
| C5 | (11 utility files — see §5.1 unbounded-tail note in PLAN) | 33 | 0 | 0 |

## Attention list

| # | Site | Function | Priority | Ambiguity | Recommended technique |
|---|------|----------|----------|-----------|----------------------|
| 1 | `expr.go:2053` | `parseInExpr` | **high** | `IN (expr_list)` vs `IN (select)` | T3 scan or T8 FIRST-set predicate for subquery-start |
| 2 | `select.go:1347` | `parseLateralTableRef` | **high** | `LATERAL (select)` vs `LATERAL XMLTABLE(...)` vs `LATERAL JSON_TABLE(...)` | T6 separate dispatch arms keyed by XMLTABLE / JSON_TABLE keywords |
| 3 | `expr.go:1609` | `parseArraySubscript` | med | `ARRAY[...]` (constructor) vs `ARRAY(select)` (sublink) | T1 peek on `[` vs `(` at callsite entry |
| 4 | `expr.go:1610` | `parseArraySubscript` | med | `ARRAY(select)` contract — does omni route all non-`[` to subquery correctly? | Oracle cross-check before classifying |

## Cluster tables

### C1: table_ref / select_with_parens / joined_table

| Site | Function | Nonterminal | Ambiguity | Technique | Aligned | Priority | Notes |
|------|----------|-------------|-----------|-----------|---------|----------|-------|
| select.go:77 | parseSelectWithParens | select_with_parens | no | T5 | yes | high | Unconditional expect('('); no disambiguation |
| select.go:166 | parseSelectClausePrimary | select_with_parens | yes | T5 | yes | high | Left-factored to parseSelectNoParens |
| select.go:770 | parseCommonTableExpr | name_list | no | T1 | yes | med | Optional paren for column names |
| select.go:1162 | parseParenTableRef | select_with_parens, joined_table | yes | T3 | yes | high | Phase 0: parenBeginsSubquery scanner |
| select.go:1225 | parseJoinedTable | joined_table | no | T6 | yes | high | Phase 0: T6/T7 dedicated nonterminal |
| select.go:1347 | parseLateralTableRef | select_with_parens, xmltable, json_table | yes | T1 | no | high | Peek only; needs T3/T6 for xmltable/json_table |
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
| expr.go:1609 | parseArraySubscript | ARRAY_SUBLINK | yes | T1 | no | med | ARRAY [ ... ] vs ARRAY ( SELECT ) ambiguity; needs T3 |
| expr.go:1610 | parseArraySubscript | ARRAY_SUBLINK | no | T1 | unclear | med | ARRAY(...) always subquery? Cross-check PG 17 |
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
| expr.go:2053 | parseInExpr | in_expr | yes | T1 | no | high | IN (list) vs IN (SELECT) ambiguity; needs T3/T8 |
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

## TOP 10 PRIORITY MISALIGNED SITES (Phase 2 Scenario Ordering)

1. **PRIORITY HIGH** — expr.go:2053 | parseInExpr | IN (expr_list) vs IN (SELECT) ambiguity | Affects all IN expressions with subqueries | **Test:** `WHERE x IN (1, 2)` vs `WHERE x IN (SELECT y FROM t)` vs `WHERE x IN (SELECT 1 UNION SELECT 2)` | **PG 17 behavior:** FIRST set includes SELECT/VALUES/WITH/TABLE; omni 1-token peek can't distinguish | **Technique:** T3 (scan for SELECT/VALUES/WITH inside parens) or T8 (FIRST-set predicate for subquery-start)

2. **PRIORITY HIGH** — select.go:1347 | parseLateralTableRef | LATERAL (SELECT) vs LATERAL XMLTABLE vs LATERAL JSON_TABLE | Peek only checks '('; can't route to correct handler | **Test:** `FROM LATERAL (SELECT 1)` vs `FROM LATERAL XMLTABLE(...)` vs `FROM LATERAL JSON_TABLE(...)` | **PG 17 behavior:** Each is a distinct production with disjoint FIRST sets after paren | **Technique:** T3 scan or T6 dedicated nonterminals for xmltable/json_table paths

3. **PRIORITY HIGH** — expr.go:1609 | parseArraySubscript | ARRAY [ ... ] vs ARRAY ( SELECT ) syntax | Current code treats all as ARRAY(SELECT); [ operator is unhandled | **Test:** `ARRAY[1, 2, 3]` vs `ARRAY(SELECT 1)` vs bare `ARRAY` | **PG 17 behavior:** ARRAY followed by [ or ( disambiguates; omni's code path confusion | **Technique:** T1 peek on '(' vs '[' at callsite entry

4. **PRIORITY MED** — expr.go:1610 | parseArraySubscript | ARRAY(SELECT) inside expression context | Assumes all ARRAY(...) are subqueries; need to check against PG | **Test:** `SELECT ARRAY(SELECT 1)` vs `SELECT ARRAY(ROWS FROM f(...))` | **PG 17 behavior:** ARRAY ( ... ) expects a subquery (SELECT/TABLE/VALUES) | **Technique:** T8 FIRST-set check to confirm parseSelectStmtForExpr is correct delegation

5. **PRIORITY MED** — select.go:770 | parseCommonTableExpr | Paren-wrapped column list optional | Currently T1 peek; confirm PG alignment | **Test:** `WITH cte AS (...) SELECT ...` vs `WITH cte(a, b) AS (...) SELECT ...` | **PG 17 behavior:** name opt_name_list AS — opt_name_list is optional | **Technique:** Already T1; verify test coverage exists

6. **PRIORITY MED** — type.go:313 through type.go:687 | parseOptFloat / parseBit / parseVarcharType / etc. | All type-modifier parens use T1 peek; cluster alignment check | **Test:** `FLOAT(24)` vs `FLOAT` vs `VARCHAR(255)` vs `VARCHAR` vs `BIT(10)` vs `BIT` vs `CHAR(20)` vs `CHAR` | **PG 17 behavior:** All are optional parens after keyword | **Technique:** T1 sufficient; batch test coverage validation


