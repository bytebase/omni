# TiDB parser.y gap inventory

Date: 2026-11-07. The first measurement of the parser.y completeness corpus
(`tidb/parser/testdata/parser_y_corpus.tsv`, see the design's §3.3). This is the
grammar work list it calls for: it replaces the pingcap-corpus inventory
(`2026-11-03-tidb-grammar-gap-inventory.md`) as the measure of what is left.

## The measurement

parser.y at v8.5.5 has 690 rules and 2,946 alternatives. Three are reported
uncovered because they cannot be exercised in the default SQL mode the oracle
measures — the `||` concat (`SimpleExpr#15`, needs `PIPES_AS_CONCAT`) and the
high-precedence `NOT` (`NotSym#2`, `SimpleExpr#16`, need `HIGH_NOT_PRECEDENCE`)
— leaving 2,943. The shortest statement exercising each was measured against
omni and against TiDB v8.5.5 via `PREPARE`:

| | first | after #724 | now |
|---|---:|---:|---:|
| both accept | 2,670 | 2,766 | 2,782 |
| gap: TiDB accepts, omni rejects | 152 | 56 | 40 |
| over-acceptance: TiDB rejects, omni accepts | 0 | 0 | 0 |
| both reject | 121 | 121 | 121 |
| no verdict | 0 | 0 | 0 |

omni parses **2,782 of the 2,822 alternatives TiDB accepts — 98.6%**, with zero
over-acceptances. The 121 both-reject rows are alternatives the server itself
refuses on a semantic check (`CREATE USER CURRENT_USER IDENTIFIED BY PASSWORD`
with a non-hash, a repeated option, and so on); they are not work.

The **after #724** column is after the `CURRENT_USER`-in-username-position
cluster closed (`fix(tidb,mysql,mariadb,oceanbase): accept CURRENT_USER as a
username`): 96 account-management rows flipped from gap to both-accept, zero new
over-acceptances. The **now** column is after the set-operation / LIMIT cluster
closed (`fix(tidb): accept ? as a LIMIT operand and ASC/DESC on GROUP BY items`):
16 rows flipped from gap to both-accept, zero new over-acceptances.

## The gaps, by theme

The first measurement's dominant root cause was one form, not 72 unrelated ones:
the generator writes `CURRENT_USER` wherever the grammar wants a username (its
shortest derivation), and omni did not accept `CURRENT_USER` in username
position. That cluster — 96 account-management rows across `CreateUserStmt`,
`AlterUserStmt`, `DropUserStmt`, `RenameUserStmt`, `GrantRoleStmt`,
`RevokeRoleStmt`, `GrantProxyStmt`, `SetDefaultRoleStmt`, `SET PASSWORD FOR`,
and every option list that embeds a `Username`/`UserSpec`/`UserToUser` —
**closed** when `parseUserAndHost` learned the `"CURRENT_USER" OptionalBraces`
alternative (parser.y:11346). CURRENT_USER is a `Username` and not a
`Rolename`, so it stays refused in every role position (CREATE/DROP/SET ROLE,
the DEFAULT ROLE clause, SET DEFAULT ROLE's role list, SHOW GRANTS ... USING).
The same form was ported to the `mysql`, `mariadb`, and `oceanbase` forks, each
guarded against over-acceptance at the positions its engine, measured, refuses
CURRENT_USER in (mysql: `SET DEFAULT ROLE` target; oceanbase: `RENAME USER` and
`SET PASSWORD FOR`) and kept at the few it reads it in (mysql: `SHOW GRANTS ...
USING`; oceanbase: `DROP ROLE`).

The set-operation / LIMIT cluster (~16 rows: `SetOprStmt`,
`SetOprStmtWithLimitOrderBy`, `SelectStmtLimit(Opt)`, `LimitOption`,
`LimitClause`, `FetchFirstOpt`, `Order`, `ByItem`) **closed** when
`parseLimitValue` learned the `?` parameter marker (`LimitOption: LengthNum |
paramMarker`, parser.y) — a new leaf `ParamMarkerExpr` node, wired only into the
LIMIT/FETCH operand, not yet general expression position — and GROUP BY items
learned the optional ASC/DESC direction (`ByItem: Expression Order`), which moved
`SelectStmt.GroupBy` to `[]*OrderByItem` so a direction can ride on each item.
The general `?`-in-expression rows (`SimpleExpr`, `WindowFrameStart/Bound`,
`OptLeadLagInfo`) stay gaps: they are a separate param-marker-as-expression
cluster.

The 40 that remain are unrelated to usernames or LIMIT:

- **Stored procedures (~8 rows).** `procedurceElseIfs`, `ProcedureProcStmt(1s)`,
  `ProcedureIf(stmt)`, `ProcedureProcStmt`.
- **Statistics, hypo indexes, TiFlash (~9 rows).** `StatsType`, `AlterTableSpec`
  (`SET HYPO TIFLASH REPLICA`), `DropIndexStmt` (`DROP HYPO INDEX`),
  `IndexNameList`.
- **LOAD DATA field options (~4 rows).** `FieldItem(List)`, `CharsetOpt`.
- **Window frames (~3 rows).** `WindowFrameStart`, `WindowFrameBound`,
  `OptLeadLagInfo`.
- **Expressions (~6 rows).** `SimpleExpr`/`ColumnOption` `COLLATE` on an
  operand, `FunctionCallNonKeyword` (`WEIGHT_STRING ... AS`).
- **SHOW / IMPORT / SEQUENCE (~10 rows).** `ShowTargetFilterable`, `ShowStmt`,
  `CancelImportStmt`, `Statement` (`CANCEL IMPORT JOB`), `CharsetKw`,
  `IgnoreOptional`, `TableOption`/`PartDefOption` on `CREATE SEQUENCE`.

A fix lands as a small PR of the kind §3.2's work landed as: close the rule,
rerun `TestTiDBParserYCorpusOmniVerdicts -update`, and the diff is the rows it
bought. Check `mysql`, `mariadb`, and `oceanbase` for the same gap.

## Every rule with a gap

One example statement per rule; `n` is how many of the rule's alternatives gap.


| Rule | n | Example (shortest form) |
|---|---:|---|
| `AlterTableSpec` | 3 | ALTER TABLE ident DROP STATS_EXTENDED ident |
| `SimpleExpr` | 3 | DO ? |
| `StatsType` | 3 | ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident) |
| `procedurceElseIfs` | 3 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `FieldItem` | 2 | SELECT * INTO OUTFILE 's' FIELDS DEFINED NULL BY 's' |
| `FunctionCallNonKeyword` | 2 | DO WEIGHT_STRING(ident AS BINARY(1)) |
| `IndexNameList` | 2 | ANALYZE TABLE ident INDEX PRIMARY |
| `PartDefOption` | 2 | CREATE SEQUENCE ident NODEGROUP 1 |
| `ProcedureProcStmt1s` | 2 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ShowTargetFilterable` | 2 | SHOW AFFINITY |
| `CancelImportStmt` | 1 | CANCEL IMPORT JOB 1 |
| `CharsetKw` | 1 | SHOW CHAR SET |
| `CharsetOpt` | 1 | LOAD DATA INFILE 's' INTO TABLE ident CHARACTER SET BINARY |
| `ColumnOption` | 1 | ALTER TABLE ident ADD ident JSON COLLATE BINARY |
| `DropIndexStmt` | 1 | DROP HYPO INDEX ident ON ident |
| `FieldItemList` | 1 | SELECT * INTO OUTFILE 's' FIELDS OPTIONALLY ENCLOSED BY 's' OPTIONALLY ENCLOSED BY 's' |
| `IgnoreOptional` | 1 | ALTER IGNORE TABLE ident |
| `OptLeadLagInfo` | 1 | DO LEAD(ident, ?) OVER ident |
| `ProcedureIf` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ProcedureIfstmt` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ProcedureProcStmt` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ShowStmt` | 1 | SHOW IMPORT JOB 1 |
| `Statement` | 1 | CANCEL IMPORT JOB 1 |
| `TableOption` | 1 | CREATE SEQUENCE ident AFFINITY 's' |
| `WindowFrameBound` | 1 | DO ROW_NUMBER() OVER(ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) |
| `WindowFrameStart` | 1 | DO ROW_NUMBER() OVER(ROWS ? PRECEDING) |
