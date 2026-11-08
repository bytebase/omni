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

| | first | after #725 | after #731 | after #733 | after #736 | after #737 | after #738 | now |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| both accept | 2,670 | 2,782 | 2,799 | 2,805 | 2,810 | 2,814 | 2,818 | 2,821 |
| gap: TiDB accepts, omni rejects | 152 | 40 | 23 | 17 | 12 | 8 | 4 | 1 |
| over-acceptance: TiDB rejects, omni accepts | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
| both reject | 121 | 121 | 121 | 121 | 121 | 121 | 121 | 121 |
| no verdict | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

omni parses **2,821 of the 2,822 alternatives TiDB accepts — 100.0%**, with zero
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
**after #731** column is after the statistics / hypothetical-index / TiFlash
cluster closed (`fix(tidb): accept extended statistics, HYPO index and TiFlash
forms`): 9 rows flipped from gap to both-accept, zero new over-acceptances. The
**after #733** column is after the SHOW / import-job cluster closed
(`fix(tidb,mysql,mariadb, oceanbase): accept SHOW IMPORT JOB(S), CANCEL IMPORT
JOB, SHOW AFFINITY and SHOW CHAR SET`): 6 rows flipped from gap to both-accept,
zero new over-acceptances. The **after #736** column is after the collate /
convert / weight_string cluster closed (`fix(tidb,mysql,mariadb,oceanbase):
accept COLLATE on column option and operand, CONVERT(... USING ...),
WEIGHT_STRING(... AS ...)`): 5 rows flipped from gap to both-accept, zero new
over-acceptances. The reserved word BINARY became a collation name
(`ColumnOption`/`SimpleExpr` COLLATE) and a charset name (`CONVERT ... USING`),
and `WEIGHT_STRING` gained its AS CHAR(n) / CHARACTER(n) / BINARY(n) clause. The
dialect split is sharp and was measured per engine: TiDB takes `COLLATE BINARY`
as both a column option and an operand; MySQL takes it only as a column option;
MariaDB and OceanBase answer 1064 to both. All four take `CONVERT(... USING
BINARY)` and the WEIGHT_STRING AS clause, but only TiDB takes its zero length
(`AS CHAR(0)`).

The **after #737** column is after the LOAD DATA / OUTFILE field-option cluster closed
(`fix(tidb,mysql,mariadb,oceanbase): accept CHARACTER SET on LOAD DATA and
DEFINED NULL BY / repeated OPTIONALLY ENCLOSED field options`): 4 rows flipped
from gap to both-accept, zero new over-acceptances. LOAD DATA's `CHARACTER SET`
name became a `CharsetName`, so the reserved word BINARY and a quoted name parse
there (`CharsetOpt`). SELECT ... INTO OUTFILE's FIELDS clause, read until now as
a fixed one-of-each set, became the same non-empty, any-order, repeatable
`FieldItemList` LOAD DATA takes, with TiDB's own `DEFINED NULL BY` among the
items (`FieldItem`, `FieldItemList`). The OUTFILE changes are TiDB's alone:
MySQL, MariaDB and OceanBase already read the repeatable list and all three
answer 1064 to DEFINED NULL BY. The `CHARACTER SET` fix was ported to the three
forks, each measured to take `CHARACTER SET BINARY` on LOAD DATA; OceanBase
diverges on the short `CHARSET` spelling, which it answers 1064 to where MySQL
and MariaDB take it, so that arm stays a plain identifier there.

The **after #738** column is after the `?`-parameter-marker-in-expression cluster
closed (`fix(tidb,mysql,mariadb,oceanbase): accept ? parameter marker as
operand, LEAD/LAG offset and window-frame bound`): 4 rows flipped from gap to
both-accept, zero new over-acceptances. The existing leaf `ParamMarkerExpr`
node — wired into the LIMIT operand by #725 — became a `SimpleExpr` alternative
(`parsePrimaryExpr` learned `?`, parser.y:7886), so a `?` is now a valid operand
anywhere an expression is: a bare `DO ?`, the LEAD/LAG offset (`OptLeadLagInfo`,
parser.y:9827) and a window-frame bound (`WindowFrameStart`/`WindowFrameBound`,
parser.y:9700, 9729) all go through that leaf. The analyzer boundary still
refuses a standalone marker, so `?` cannot reach catalog Exec. The fix was ported
to the three forks, each measured against its live engine: MySQL and OceanBase
accept `?` in all three positions; MariaDB accepts it as an operand and LEAD/LAG
offset but answers 1064 to a `?` frame bound (its bound is a numeric literal or
INTERVAL), so `mariadb`'s `parseFrameBound` keeps the marker out of that one
position alone.

The **now** column is after the CREATE SEQUENCE table-option cluster closed
(`fix(tidb): accept STORAGE ENGINE / NODEGROUP / AFFINITY table options on
CREATE SEQUENCE`): 3 rows flipped from gap to both-accept, zero new
over-acceptances. `CreateSequenceStmt`'s option list is the general
`CreateTableOptionListOpt` (parser.y:12804), so every `TableOption` TiDB admits
on CREATE TABLE it admits on a sequence; omni's `parseTableOption` was missing
three of them. `STORAGE ENGINE [=] name` (the `PartDefOption` arm that names the
same engine ENGINE does), `NODEGROUP [=] n` and `AFFINITY [=] name` (an
unreserved keyword omni lexes as an identifier, so matched as a bare word) each
gained a case, with STORAGE ENGINE and AFFINITY taking a `StringName` and
NODEGROUP a `LengthNum`, each refusal held to the engine. All three are
TiDB-only: `mysql` and `oceanbase` have no CREATE SEQUENCE, and `mariadb`'s takes
no trailing table options, so nothing ported.

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
`OptLeadLagInfo`) were a separate param-marker-as-expression cluster, since
**closed** when `parsePrimaryExpr` made the same `ParamMarkerExpr` leaf a
`SimpleExpr` alternative (see the **after #738** column above).

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

The 1 that remains is unrelated to usernames, LIMIT, procedures, statistics,
SHOW, collation/convert/weight_string, LOAD DATA / OUTFILE field options, the
`?` parameter marker, or CREATE SEQUENCE options:

- **ALTER IGNORE TABLE (1 row).** `IgnoreOptional`.

A fix lands as a small PR of the kind §3.2's work landed as: close the rule,
rerun `TestTiDBParserYCorpusOmniVerdicts -update`, and the diff is the rows it
bought. Check `mysql`, `mariadb`, and `oceanbase` for the same gap.

## Every rule with a gap

One example statement per rule; `n` is how many of the rule's alternatives gap.


| Rule | n | Example (shortest form) |
|---|---:|---|
| `IgnoreOptional` | 1 | ALTER IGNORE TABLE ident |
