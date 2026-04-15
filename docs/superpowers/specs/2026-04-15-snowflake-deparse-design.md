# Snowflake Deparse (T3.2) — Design Spec

**Date:** 2026-04-15
**Node:** T3.2 deparse-core
**Author:** Claude Sonnet 4.6

## Problem

The Snowflake parser (T1–T2, T5) produces an AST but there is no way to convert that AST back to a SQL string. This is needed for:
- T3.3 LIMIT injection rewrite (modify AST → emit SQL)
- Advisor suggestions that reconstruct SQL
- Round-trip testing of the parser

## Goal

Given any P0 statement AST node, produce a valid Snowflake SQL string that re-parses to an equivalent AST.

## Key Decisions

### Keyword casing
Always emit keywords in UPPER CASE. Identifier casing follows the original source: `Ident.Quoted` controls whether double-quotes are added; unquoted identifiers preserve their original case string (which the parser stored as-is).

### Output format
Single-line per statement. No cosmetic newlines or indentation. This simplifies round-trip testing and downstream rewriting. Readability is secondary to correctness.

### Error handling
Return errors for unsupported node types; never panic. This allows callers to handle partial support gracefully.

### Round-trip definition
`Parse(Deparse(Parse(sql))) ≡ Parse(sql)` (modulo Loc fields). The deparsed SQL need not be identical to the original — it may normalize spacing, keyword casing, and syntax variants — but the resulting AST must be structurally equivalent.

## Package structure

```
snowflake/deparse/
    deparse.go       — public API: Deparse(Node), DeparseFile(*File)
    writer.go        — internal writer helper with spacing/quoting utilities
    deparse_expr.go  — all expression node writers
    deparse_select.go — SELECT, SetOperationStmt, CTEs, JoinExpr, TableRef
    deparse_dml.go   — INSERT, UPDATE, DELETE, MERGE
    deparse_ddl.go   — all DDL (CREATE/ALTER/DROP TABLE, VIEW, DB, SCHEMA, etc.)
    deparse_test.go  — round-trip tests
```

## Coverage

All P0 node types from T1.4–T1.7, T2.1–T2.5, T5.1:

- Expressions: Literal, ColumnRef, StarExpr, BinaryExpr, UnaryExpr, ParenExpr, CastExpr (including `::` and TRY_CAST), CaseExpr (simple and searched), FuncCallExpr (including window functions and WITHIN GROUP), IffExpr, CollateExpr, IsExpr, BetweenExpr, InExpr, LikeExpr (and ILIKE/RLIKE/REGEXP/LIKE ANY), AccessExpr (colon, dot, bracket), ArrayLiteralExpr, JsonLiteralExpr, LambdaExpr, SubqueryExpr, ExistsExpr
- Query: SelectStmt (all clauses: WITH CTEs, TOP, DISTINCT/ALL, FROM, WHERE, GROUP BY variants, HAVING, QUALIFY, ORDER BY with NULLS, LIMIT, OFFSET, FETCH FIRST), SetOperationStmt (UNION/INTERSECT/EXCEPT with ALL and BY NAME), TableRef (table/subquery/func/LATERAL), JoinExpr (INNER/LEFT/RIGHT/FULL/CROSS/ASOF, ON/USING/NATURAL/MATCH_CONDITION)
- DML: InsertStmt, InsertMultiStmt (ALL/FIRST with WHEN/ELSE), UpdateStmt, DeleteStmt, MergeStmt
- DDL: CreateTableStmt (all modifiers, CTAS, LIKE, CLONE), ColumnDef, TableConstraint, AlterTableStmt (all 24 action kinds), CreateDatabaseStmt, AlterDatabaseStmt, DropDatabaseStmt, UndropDatabaseStmt, CreateSchemaStmt, AlterSchemaStmt, DropSchemaStmt, UndropSchemaStmt, DropStmt, UndropStmt, CreateViewStmt, AlterViewStmt, CreateMaterializedViewStmt, AlterMaterializedViewStmt

## Special cases

- `StarExpr` in SELECT targets: `t.Star == true` means the `t.Expr` is always a `*StarExpr`; the qualifier is taken from `StarExpr.Qualifier`, not by calling `writeExpr`.
- `InsertMultiStmt` ELSE branches: branches with nil `When` that follow conditional branches in INSERT FIRST are ELSE branches.
- Inline FK constraint: must emit `FOREIGN KEY REFERENCES table (cols)`, not bare `REFERENCES`.
- `writeExprNoLeadSpace`: helper that strips the leading space that `ensureSpace()` would insert, used inside parentheses, comma lists, and after explicit keyword+space sequences.
