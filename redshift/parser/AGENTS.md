# redshift/parser

`redshift/parser` is a fork of `pg/parser`. Follow [pg/parser/AGENTS.md](../../pg/parser/AGENTS.md): the FIRST-set discipline and the backtracking rules apply here unchanged, and this package has its own `first_sets.go`, kept in step with the parent's when a predicate changes. Its tests ask no server: the parent's oracles hold pg to PostgreSQL 17, which Redshift does not run, so the FIRST-set oracles were deleted here and the paren-dispatch cases keep only omni's committed routes (`paren_route_*_test.go`).

Redshift-only material stays here because tests read it: `PAREN_AUDIT*.md`, `PARSER_DISPATCH_AUDIT.md`, and `testdata/`. The legacy-parser gap triage lives in `../compat/`.

## Redshift-only FIRST-set predicates

Defined in `first_sets.go` alongside the predicates inherited from `pg/parser/first_sets.go`:

| Predicate | Backing tokens / words | Composes from | Unit test |
|---|---|---|---|
| `isRedshiftDocShowStart` | `USER LOCKOUT`, `COLUMN GRANTS`, `{RLS\|MASKING} POLICIES`, and `{CONSTRAINTS\|FUNCTIONS\|PARAMETERS\|PROCEDURES\|TEMPLATE\|TEMPLATES}` before a clause | unquoted keyword check + `isRedshiftShowClauseFollow` | `TestIsRedshiftDocShowStart` |
| `isRedshiftTemplateLead` | `TEMPLATE` | unquoted `TEMPLATE` token | `TestIsRedshiftTemplateLead` |
| `isRedshiftMultiwordTypenameStart` | the multiword type map: `DOUBLE PRECISION`, `CHARACTER VARYING`, `NATIONAL CHARACTER`, `BIT VARYING`, the `INTERVAL` fields, `TIMESTAMP`/`TIME WITH TIME ZONE` | the lead's continuation tokens | `TestIsRedshiftMultiwordTypenameStart` |
| `isRedshiftDatashareNameStart` | the unquoted word `DATASHARE` | the word and a name after it | `TestIsRedshiftDatashareNameStart` |
| `isRedshiftConnectStart`, `isRedshiftDebugStart` | `CONNECT`, `DEBUG` | unquoted words + `TO`/`ON` | `TestIsRedshiftFederatedPermissionStart` |
| `isRedshiftAdminsOfAccountStart`, `isRedshiftPublicGranteeStart`, `isRedshiftRoleIdentityStart` | `ADMINS OF`, `PUBLIC` (no colon), `ROLE` before a name | unquoted words + continuation | `TestIsRedshiftIdentityGranteeStart` |
| `isRedshiftRoleGranteeStart`, `isRedshiftGroupGranteeStart` | `ROLE`, `GROUP` | the token and a name after it | `TestIsRedshiftNamedGranteeStart` |
| `isRedshiftScopedPermissionStart` | a privilege followed by `FOR` | the next token | `TestIsRedshiftScopedPermissionStart` |
| `isRedshiftIntegrationOptionStart` | `ACCEPTINVCHARS`, `QUERY_ALL_STATES`, `REFRESH_INTERVAL`, `TRUNCATECOLUMNS`, `HISTORY_MODE`, `AUTO_REMEDIATION` | unquoted words | `TestIsRedshiftIntegrationOptionStart` |
| `isRedshiftCreateDatabaseSourceStart` | `FROM`, or `WITH PERMISSIONS` | the token and one lookahead | `TestIsRedshiftCreateDatabaseSourceStart` |
| `isRedshiftAlterDatabaseIntegrationStart` | the unquoted word `INTEGRATION` | the word | `TestIsRedshiftAlterDatabaseIntegrationStart` |
| `isRedshiftDataMovementOptionStart` | a ColLabel that is not `FROM`, `TO`, `WHERE` or a terminator | the category and stop tokens | `TestIsRedshiftDataMovementOptionStart` |
| `isRedshiftCopyColumnStart` | a name immediately before `FROM` or `TO` | the category and one lookahead | `TestIsRedshiftCopyColumnStart` |
| `isRedshiftConnectByStart` | `START WITH`, the unquoted `CONNECT BY` | the word and one lookahead | `TestIsRedshiftConnectByStart` |
| `isRedshiftVacuumTwoWordModeStart` | `SORT ONLY`, `DELETE ONLY` | the word and one lookahead | `TestIsRedshiftVacuumTwoWordModeStart` |
| `isRedshiftPartitionTransformCallStart` | `YEAR`, `MONTH`, `DAY`, `HOUR`, `TRUNCATE` or the word `BUCKET` before `(` | the token and one lookahead | `TestIsRedshiftPartitionTransformCallStart` |
| `isRedshiftRoleSpecStart` | a name, or `CURRENT_USER` / `CURRENT_ROLE` / `SESSION_USER` | the token after a comma | `TestIsRedshiftRoleSpecStart` |
| `isRedshiftNumberStart` | a sign and an `Iconst` | `isNumericOnlyStart`, less `FCONST` | `TestIsRedshiftNumberStart` |
| `isRedshiftIcebergMaterializedViewOptionStart` | `LOCATION`, the word `PARTITIONED`, `TABLE` | the token and the unquoted word | `TestIsRedshiftIcebergMaterializedViewOptionStart` |

The two-word predicates answer through `cursorAfterCur` in completion mode: the
cursor sits after the first word, so the peek sees only the end of the input, and
the continuation is what the user is about to write. `expect` then records it.

The predicates `isRedshiftConnectByStart` and `isRedshiftVacuumTwoWordModeStart`
read a second token because the first one is a name in the other reading:
`SELECT * FROM t start` aliases a table `start`, and `VACUUM SORT` vacuums a
table named `sort`. `isRedshiftConnectByStart` is also what
`isReservedForClause` asks, so a bare `START` or `CONNECT` still names an
alias. `isRedshiftPartitionTransformCallStart` tells `PARTITIONED BY
(BUCKET(16, a))` from a column named `bucket`, and `isRedshiftRoleSpecStart`
ends a `role_list` at the comma before a repeated `WITH` clause. Its set is
`parseNonReservedWord`'s, which is wider than `isColIdToken`'s: a type or
function name keyword such as `FULL` names a role too. `isRedshiftNumberStart`
is the docs' `NUMBER`, which `testdata/tokens.tsv` samples as an integer:
`AUTO REFRESH YES INTERVAL` and `IDENTITY`'s seed and step take no float.

## Redshift-only expression context

`CONNECT BY`'s `PRIOR` is the one prefix operator whose meaning depends on
where it sits: `Parser.inConnectBy` is true only while the clause's condition
is parsed, and `parseAExprAtom` reads `PRIOR` as
`A_Expr{AEXPR_OP, "prior"}` there and as a column name everywhere else.
