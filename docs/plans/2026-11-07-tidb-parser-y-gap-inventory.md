# TiDB parser.y gap inventory

Date: 2026-11-07. The first measurement of the parser.y completeness corpus
(`tidb/parser/testdata/parser_y_corpus.tsv`, see the design's §3.3). This is the
grammar work list it calls for: it replaces the pingcap-corpus inventory
(`2026-11-03-tidb-grammar-gap-inventory.md`) as the measure of what is left.

## The measurement

parser.y at v8.5.5 has 690 rules and 2,946 alternatives. One (`SimpleExpr#15`,
the `||` concat, which needs `PIPES_AS_CONCAT`) is reported uncovered because it
cannot be exercised in the default SQL mode the oracle measures, leaving 2,945.
The shortest statement exercising each was measured against omni and against
TiDB v8.5.5 via `PREPARE`:

| | count |
|---|---:|
| both accept | 2,672 |
| gap: TiDB accepts, omni rejects | 152 |
| over-acceptance: TiDB rejects, omni accepts | 0 |
| both reject | 121 |
| no verdict | 0 |

omni parses **2,672 of the 2,824 alternatives TiDB accepts — 94.6%**, with zero
over-acceptances. The 121 both-reject rows are alternatives the server itself
refuses on a semantic check (`CREATE USER CURRENT_USER IDENTIFIED BY PASSWORD`
with a non-hash, a repeated option, and so on); they are not work.

## The 152 gaps, by theme

The dominant root cause is one form, not 72 unrelated ones: the generator writes
`CURRENT_USER` wherever the grammar wants a username (its shortest derivation),
and omni does not accept `CURRENT_USER` in username position. Most of the
account-management cluster below turns on that single form, so one fix closes
many rows.

- **Account management (~50 rows).** `CreateUserStmt`, `AlterUserStmt`,
  `DropUserStmt`, `RenameUserStmt`, `GrantRoleStmt`, `RevokeRoleStmt`,
  `GrantProxyStmt`, `SetDefaultRoleStmt`, and their option lists —
  `AuthOption`, `AuthString`, `AuthPlugin`, `HashString`,
  `PasswordOrLockOption(s)(List)`, `RequireClause(Opt)`, `RequireList(Element)`,
  `ConnectionOption(s)(List)`, `CommentOrAttributeOption`,
  `ResourceGroupNameOption`, `UserSpec(List)`, `Username(List)`,
  `UserToUser(List)`, `RolenameWithoutIdent`, `RoleOrPrivElem(List)`,
  `ExtendedPriv`, `IfExists`, `IfNotExists`.
- **Set operations with ORDER BY / LIMIT (~16 rows).** `SetOprStmt`,
  `SetOprStmtWithLimitOrderBy`, `SelectStmtLimit(Opt)`, `LimitOption`,
  `LimitClause`, `FetchFirstOpt`, `Order`, `ByItem`.
- **Stored procedures (~12 rows).** `procedurceElseIfs`, `ProcedureProcStmt(1s)`,
  `ProcedureIf(stmt)`, `ProcedureProcStmt`.
- **Statistics, hypo indexes, TiFlash (~9 rows).** `StatsType`, `AlterTableSpec`
  (`SET HYPO TIFLASH REPLICA`), `DropIndexStmt` (`DROP HYPO INDEX`),
  `IndexNameList`.
- **LOAD DATA field options (~5 rows).** `FieldItem(List)`, `CharsetOpt`.
- **Window frames (~3 rows).** `WindowFrameStart`, `WindowFrameBound`,
  `OptLeadLagInfo`.
- **Expressions (~4 rows).** `SimpleExpr`/`ColumnOption` `COLLATE` on an
  operand, `FunctionCallNonKeyword` (`WEIGHT_STRING ... AS`).
- **SHOW / IMPORT (~5 rows).** `ShowTargetFilterable`, `ShowStmt`,
  `CancelImportStmt`, `CharsetKw`, `CharsetOpt`, `IgnoreOptional`,
  `TableOption`/`PartDefOption` on `CREATE SEQUENCE`,
  `SetStmt` (`SET PASSWORD FOR`).

A fix lands as a small PR of the kind §3.2's work landed as: close the rule,
rerun `TestTiDBParserYCorpusOmniVerdicts -update`, and the diff is the rows it
bought. Check `mysql`, `mariadb`, and `oceanbase` for the same gap.

## Every rule with a gap

One example statement per rule; `n` is how many of the rule's alternatives gap.

| Rule | n | Example (shortest form) |
|---|---:|---|
| `PasswordOrLockOption` | 14 | CREATE USER CURRENT_USER ACCOUNT LOCK |
| `Statement` | 9 | CANCEL IMPORT JOB 1 |
| `AuthOption` | 6 | CREATE USER CURRENT_USER |
| `RequireListElement` | 5 | CREATE USER CURRENT_USER REQUIRE SAN 's' |
| `SetOprStmtWithLimitOrderBy` | 4 | (SELECT *) LIMIT ? |
| `RoleOrPrivElem` | 4 | GRANT 's' TO CURRENT_USER |
| `RequireClause` | 4 | CREATE USER CURRENT_USER REQUIRE SSL |
| `ConnectionOption` | 4 | CREATE USER CURRENT_USER WITH MAX_QUERIES_PER_HOUR 1 |
| `procedurceElseIfs` | 3 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `StatsType` | 3 | ALTER TABLE ident ADD STATS_EXTENDED ident DEPENDENCY(ident) |
| `SimpleExpr` | 3 | DO ? |
| `SelectStmtLimit` | 3 | SELECT * LIMIT ? |
| `RequireList` | 3 | CREATE USER CURRENT_USER REQUIRE ISSUER 's' |
| `CommentOrAttributeOption` | 3 | CREATE USER CURRENT_USER |
| `AlterTableSpec` | 3 | ALTER TABLE ident DROP STATS_EXTENDED ident |
| `UsernameList` | 2 | DROP USER CURRENT_USER |
| `UserToUserList` | 2 | RENAME USER CURRENT_USER TO CURRENT_USER |
| `UserSpecList` | 2 | CREATE USER CURRENT_USER |
| `ShowTargetFilterable` | 2 | SHOW AFFINITY |
| `SetOprStmt` | 2 | (SELECT *) LIMIT ? |
| `RolenameWithoutIdent` | 2 | GRANT 's' TO CURRENT_USER |
| `RoleOrPrivElemList` | 2 | GRANT 's' TO CURRENT_USER |
| `ResourceGroupNameOption` | 2 | CREATE USER CURRENT_USER |
| `RequireClauseOpt` | 2 | CREATE USER CURRENT_USER |
| `ProcedureProcStmt1s` | 2 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `PasswordOrLockOptions` | 2 | CREATE USER CURRENT_USER |
| `PasswordOrLockOptionList` | 2 | CREATE USER CURRENT_USER ACCOUNT UNLOCK |
| `PartDefOption` | 2 | CREATE SEQUENCE ident NODEGROUP 1 |
| `Order` | 2 | SELECT * GROUP BY ident ASC |
| `IndexNameList` | 2 | ANALYZE TABLE ident INDEX PRIMARY |
| `IfNotExists` | 2 | CREATE USER CURRENT_USER |
| `IfExists` | 2 | ALTER USER CURRENT_USER |
| `HashString` | 2 | CREATE USER CURRENT_USER IDENTIFIED BY PASSWORD 's' |
| `FunctionCallNonKeyword` | 2 | DO WEIGHT_STRING(ident AS BINARY(1)) |
| `FieldItem` | 2 | SELECT * INTO OUTFILE 's' FIELDS DEFINED NULL BY 's' |
| `ExtendedPriv` | 2 | GRANT ident TO CURRENT_USER |
| `DropUserStmt` | 2 | DROP USER CURRENT_USER |
| `ConnectionOptions` | 2 | CREATE USER CURRENT_USER |
| `ConnectionOptionList` | 2 | CREATE USER CURRENT_USER WITH MAX_QUERIES_PER_HOUR 1 |
| `WindowFrameStart` | 1 | DO ROW_NUMBER() OVER(ROWS ? PRECEDING) |
| `WindowFrameBound` | 1 | DO ROW_NUMBER() OVER(ROWS BETWEEN ? FOLLOWING AND UNBOUNDED PRECEDING) |
| `Username` | 1 | DROP USER CURRENT_USER |
| `UserToUser` | 1 | RENAME USER CURRENT_USER TO CURRENT_USER |
| `UserSpec` | 1 | CREATE USER CURRENT_USER |
| `TableOption` | 1 | CREATE SEQUENCE ident AFFINITY 's' |
| `ShowStmt` | 1 | SHOW IMPORT JOB 1 |
| `SetStmt` | 1 | SET PASSWORD FOR CURRENT_USER = 's' |
| `SetDefaultRoleStmt` | 1 | SET DEFAULT ROLE NONE TO CURRENT_USER |
| `SelectStmtLimitOpt` | 1 | SELECT * LIMIT ? |
| `RevokeRoleStmt` | 1 | REVOKE 's' FROM CURRENT_USER |
| `RenameUserStmt` | 1 | RENAME USER CURRENT_USER TO CURRENT_USER |
| `ProcedureProcStmt` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ProcedureIfstmt` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `ProcedureIf` | 1 | CREATE PROCEDURE ident() IF ident THEN COMMIT; END IF |
| `OptLeadLagInfo` | 1 | DO LEAD(ident, ?) OVER ident |
| `LimitOption` | 1 | SELECT * LIMIT ? |
| `LimitClause` | 1 | DELETE FROM ident LIMIT ? |
| `IgnoreOptional` | 1 | ALTER IGNORE TABLE ident |
| `GrantRoleStmt` | 1 | GRANT 's' TO CURRENT_USER |
| `GrantProxyStmt` | 1 | GRANT PROXY ON CURRENT_USER TO CURRENT_USER |
| `FieldItemList` | 1 | SELECT * INTO OUTFILE 's' FIELDS OPTIONALLY ENCLOSED BY 's' OPTIONALLY ENCLOSED BY 's' |
| `FetchFirstOpt` | 1 | SELECT * FETCH FIRST ? ROW ONLY |
| `DropIndexStmt` | 1 | DROP HYPO INDEX ident ON ident |
| `CreateUserStmt` | 1 | CREATE USER CURRENT_USER |
| `ColumnOption` | 1 | ALTER TABLE ident ADD ident JSON COLLATE BINARY |
| `CharsetOpt` | 1 | LOAD DATA INFILE 's' INTO TABLE ident CHARACTER SET BINARY |
| `CharsetKw` | 1 | SHOW CHAR SET |
| `CancelImportStmt` | 1 | CANCEL IMPORT JOB 1 |
| `ByItem` | 1 | SELECT * GROUP BY ident ASC |
| `AuthString` | 1 | CREATE USER CURRENT_USER IDENTIFIED BY 's' |
| `AuthPlugin` | 1 | CREATE USER CURRENT_USER IDENTIFIED WITH 's' |
| `AlterUserStmt` | 1 | ALTER USER CURRENT_USER |

