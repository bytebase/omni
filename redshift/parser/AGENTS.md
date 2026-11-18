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
