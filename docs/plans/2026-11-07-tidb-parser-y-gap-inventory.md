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

| | first | after #724 | after #725 | after #730 | after #731 | now |
|---|---:|---:|---:|---:|---:|---:|
| both accept | 2,670 | 2,766 | 2,782 | 2,790 | 2,799 | 2,805 |
| gap: TiDB accepts, omni rejects | 152 | 56 | 40 | 32 | 23 | 17 |
| over-acceptance: TiDB rejects, omni accepts | 0 | 0 | 0 | 0 | 0 | 0 |
| both reject | 121 | 121 | 121 | 121 | 121 | 121 |
| no verdict | 0 | 0 | 0 | 0 | 0 | 0 |

omni parses **2,805 of the 2,822 alternatives TiDB accepts — 99.4%**, with zero
over-acceptances. The 121 both-reject rows are alternatives the server itself
refuses on a semantic check (`CREATE USER CURRENT_USER IDENTIFIED BY PASSWORD`
with a non-hash, a repeated option, and so on); they are not work.

The **after #724** column is after the `CURRENT_USER`-in-username-position
cluster closed (`fix(tidb,mysql,mariadb,oceanbase): accept CURRENT_USER as a
username`): 96 account-management rows flipped from gap to both-accept, zero new
over-acceptances. The **after #725** column is after the set-operation / LIMIT
cluster closed (`fix(tidb): accept ? as a LIMIT operand and ASC/DESC on GROUP BY
items`): 16 rows flipped from gap to both-accept, zero new over-acceptances. The
**after #730** column is after the stored-procedure cluster closed
(`fix(tidb,mysql,mariadb,oceanbase): accept a bare IF ... END IF procedure
body`): 8 rows flipped from gap to both-accept, zero new over-acceptances. The
**now** column is after the statistics / hypothetical-index / TiFlash cluster
closed (`fix(tidb): accept extended statistics, HYPO index and TiFlash forms`):
9 rows flipped from gap to both-accept, zero new over-acceptances. The **now**
column is after the SHOW / import-job cluster closed (`fix(tidb,mysql,mariadb,
oceanbase): accept SHOW IMPORT JOB(S), CANCEL IMPORT JOB, SHOW AFFINITY and SHOW
CHAR SET`): 6 rows flipped from gap to both-accept, zero new over-acceptances.

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

The stored-procedure cluster (8 rows: `procedurceElseIfs`,
`ProcedureProcStmt1s`, `ProcedureIf`, `ProcedureIfstmt`, `ProcedureProcStmt`)
**closed** in the splitter, not the parser. A routine body is a single
`sp_proc_stmt`, which may be a bare `IF ... THEN ... END IF` with its own inner
`;` and no surrounding `BEGIN ... END`. `Split` broke that body at the inner `;`
because `IF` is a compound opener only at a statement-start position (to keep
`SELECT IF(...)` and `DROP TABLE IF EXISTS` out), and the position right after a
`CREATE PROCEDURE` parameter list was not marked as one. The splitter now tracks
that parameter list and marks the body that follows as statement-start, so the
bare `IF` opens a compound and the inner `;` no longer splits. `WHILE`/`CASE`/
`REPEAT`/`LOOP` bodies already opened unconditionally, so only `IF` was affected.
The same splitter fix ported to `mysql`, `mariadb`, and `oceanbase`, each
measured against its own engine (a bare `IF`/`WHILE`/`CASE` body is valid in all
three); zero new over-acceptances.

The statistics / hypothetical-index / TiFlash cluster (9 rows: `StatsType`,
`AlterTableSpec` ×3, `DropIndexStmt`, `IndexNameList` ×2) **closed** in the
parser, all TiDB-only extensions the forks do not share. `ALTER TABLE` learned
three specs: `SET HYPO TIFLASH REPLICA n` (a `Hypo` flag on the existing
`ATSetTiFlashReplica`), `ADD STATS_EXTENDED name type (cols)` and `DROP
STATS_EXTENDED name` (two new `ATAddStatistics`/`ATDropStatistics` commands over
`CARDINALITY`/`DEPENDENCY`/`CORRELATION`, the same closed `StatsType` list CREATE
STATISTICS already read). `DROP HYPO INDEX name ON tbl` became its own
`DropIndexStmt` arm (an `IsHypo` flag; the arm takes no lock or algorithm).
`IndexNameList` learned `PRIMARY` as a name, in and beyond the first slot, so
`ANALYZE TABLE t INDEX PRIMARY` and `INDEX, PRIMARY` parse. `mysql`, `mariadb`
and `oceanbase` reject every one of these forms, measured; nothing ported.

The SHOW / import-job cluster (6 rows: `ShowTargetFilterable` ×2, `ShowStmt`,
`CancelImportStmt`, `Statement`, `CharsetKw`) **closed** in the parser. SHOW
learned three targets: `IMPORT JOB n` / `IMPORT JOBS` (the same shape as
`DISTRIBUTION JOB{,S}`, reusing `ShowStmt.JobID`), `AFFINITY` (an unreserved
keyword omni lexes as an identifier, so matched as a bare word), and `CHAR SET`
— `CharsetKw`'s third spelling, a synonym for `CHARACTER SET`. `CANCEL IMPORT
JOB n` became a new `CancelImportJobStmt` beside `CancelDistributionJobStmt`.
The import-job and `AFFINITY` forms are TiDB-only (`mysql`, `mariadb`,
`oceanbase` reject them, measured), but `SHOW CHAR SET` is shared: all three
engines accept it, and the one-line `CHAR SET` arm ported to each, measured,
zero new over-acceptances.

The 17 that remain are unrelated to usernames, LIMIT, procedures, statistics, or
SHOW:

- **LOAD DATA field options (~4 rows).** `FieldItem(List)`, `CharsetOpt`.
- **Window frames (~3 rows).** `WindowFrameStart`, `WindowFrameBound`,
  `OptLeadLagInfo`.
- **Expressions (~6 rows).** `SimpleExpr`/`ColumnOption` `COLLATE` on an
  operand, `FunctionCallNonKeyword` (`WEIGHT_STRING ... AS`).
- **CREATE SEQUENCE options and ALTER IGNORE TABLE (~4 rows).**
  `TableOption`/`PartDefOption` on `CREATE SEQUENCE`, `IgnoreOptional`.

A fix lands as a small PR of the kind §3.2's work landed as: close the rule,
rerun `TestTiDBParserYCorpusOmniVerdicts -update`, and the diff is the rows it
bought. Check `mysql`, `mariadb`, and `oceanbase` for the same gap.

## Every rule with a gap

One example statement per rule; `n` is how many of the rule's alternatives gap.


| Rule | n | Example (shortest form) |
|---|---:|---|
| `SimpleExpr` | 3 | DO ? |
| `FieldItem` | 2 | SELECT * INTO OUTFILE 's' FIELDS DEFINED NULL BY 's' |
| `FunctionCallNonKeyword` | 2 | DO WEIGHT_STRING(ident AS BINARY(1)) |
| `PartDefOption` | 2 | CREATE SEQUENCE ident NODEGROUP 1 |
| `CharsetOpt` | 1 | LOAD DATA INFILE 's' INTO TABLE ident CHARACTER SET BINARY |
| `ColumnOption` | 1 | ALTER TABLE ident ADD ident JSON COLLATE BINARY |
| `FieldItemList` | 1 | SELECT * INTO OUTFILE 's' FIELDS OPTIONALLY ENCLOSED BY 's' OPTIONALLY ENCLOSED BY 's' |
| `IgnoreOptional` | 1 | ALTER IGNORE TABLE ident |
| `OptLeadLagInfo` | 1 | DO LEAD(ident, ?) OVER ident |
| `TableOption` | 1 | CREATE SEQUENCE ident AFFINITY 's' |
| `WindowFrameBound` | 1 | DO ROW_NUMBER() OVER(ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) |
| `WindowFrameStart` | 1 | DO ROW_NUMBER() OVER(ROWS ? PRECEDING) |
