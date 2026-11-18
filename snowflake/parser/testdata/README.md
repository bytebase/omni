# Snowflake parser test data

## `compat.tsv`: the docs' rows

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/SnowflakeDocs.g4`, then its near-misses: the statement with one
token dropped, the end cut off, a stray word inserted, or a keyword swapped.
Each row has omni's verdict and the judge's. `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat`, starts no server,
regenerates the rows from the grammar, and asks the recognizer again on every
run.

The pin moved from the pages read 2026-11-13
(`sha256:334be655...`) to the pages read 2026-11-18
(`sha256:9c6ceb51...`). The live pages also gained the `PLAN` privilege on
`GRANT`/`REVOKE`, new `SHOW GRANTS` targets, and `INCLUDE HIERARCHY`; the
generated grammar and rows moved with them.

| Rows | Old pin total | Old agree | Old docs only | New total | New agree | New docs only |
| --- | --- | --- | --- | --- | --- | --- |
| One per alternative | 19,775 | 19,638 | 137 | 19,792 | 19,724 | 68 |
| Near-misses | 53,876 | 758 | 49 | 53,914 | 758 | 46 |
| All | 73,651 | 20,396 | 186 | 73,706 | 20,482 | 114 |

`undocumented` rows are the near-misses the docs do not derive; docs say what is
valid, not what Snowflake refuses, so they pass whatever omni says. The new pin
has 53,110 of them.

## The grammar and its fixes

`SnowflakeDocs.g4` is the Syntax sections of the 779 pages
https://docs.snowflake.com/en/sql-reference/sql-all lists, converted by
`../cmd/docsgrammar` with `enginetest/grammar/docs`. 765 pages have syntax and
758 have a rule. `../../testdata/pins.tsv` pins the index URL, the SHA-256 of
the blocks read, and the date.

`../../testdata/docs-fixes.tsv` holds the hand corrections. This pin bump drops
the two rows whose pages moved on 2026-11-18: `create-dbt-project` and
`execute-dbt-project` now write `DBT_VERSION = '<version_number>'` themselves.
The other fixes added by this bump are page by page, each with its reason in
the file: the required `INCREMENT` and `SET` action of `ALTER SEQUENCE`, the
internal stage's empty `SET`, the resource monitor's `SET`/`TRIGGERS` choice,
the enumerated DCM `LOG_LEVEL`, the service role's
`<service_name>!<service_role_name>` name, the workspace form of `EXECUTE DBT
PROJECT`, the Iceberg REST column parentheses and `[ ... ]` tails, the `ALTER
LISTING` version name and stage/manifest values, the missing `|` in `ALTER
SCHEMA`'s `UNSET` list, the DCM `USING` name/value pairs and bind variable, and
`CREATE EXTERNAL TABLE`'s template `LOCATION` and non-empty column list. The
DCM bind variable is `$<session_variable>` or the client `?`; the scripting
`:name` form has no generated row because this notation writes `:` and the name
as two tokens, which the row generator separates.

The object-type sample still reads `ACCOUNT` as an object type that takes a
name. Removing it is a converter change, which a pin-bump PR may not carry
(`compatcheck` refuses one), so the `GRANT OWNERSHIP`, `GRANT`/`REVOKE CALLER`
and `SHOW GRANTS ON ACCOUNT <name>` rows it makes up remain for a follow-up.

`../../testdata/tokens.tsv` samples every token class the grammar uses.

To rebuild the grammar from the live docs (a pin bump), then the rows:

```bash
go run ./snowflake/parser/cmd/docsgrammar -cache /tmp/snowflake-docs
OMNI_RECORD=1 go test ./snowflake/parser/ -run '^TestCompatibility$'
```

After a parser fix, commit the rows that now agree:

```bash
go test ./snowflake/parser/ -run '^TestCompatibility$' -update-compat
```

## What is left: 114 rows

Every remaining row is `reject` in omni and `accept` in the docs. They are
grouped here by the page whose notation or reading leaves them. None is a
statement omni used to accept: no row that agreed with the docs stopped
agreeing. Two class-2 rows moved from the docs' `accept` to `undocumented`,
because the grammar no longer derives them; they still pass: the `ALTER
SEQUENCE IF` and `ALTER RESOURCE MONITOR IF` near-misses.

| Page | Rows | What remains |
| --- | --- | --- |
| alter-stage | 18 | `ALTER STAGE <name> SET` with no property; the external form's `SET { [ externalStageParams ] [ FILE_FORMAT ] [ COMMENT ] }` and its all-optional GCS/Azure `externalStageParams` have no notation for "at least one", so a page-local fix would invent a required parameter. |
| the combined `statement` rule | 19 | Dropping a name leaves another page's statement: `DROP APPLICATION PACKAGE` is `DROP APPLICATION` named `PACKAGE`, `USE ROLE` is `USE [ DATABASE ] ROLE`, `ALTER DATABASE ROLE UNSET COMMENT` is `ALTER DATABASE` named `ROLE`, and `GRANT ... TO APPLICATION` loses its grantee. The `GRANT OWNERSHIP`, `GRANT`/`REVOKE CALLER` and `SHOW GRANTS ON ACCOUNT <name>` statements made by the `ACCOUNT` object-type sample are here too. These are judge readings, not a notation defect on one page. |
| alter-task-graph | 11 | `<property> = <value>`, `<default_property>` and a repeated property run. The page says the properties are the ones in CREATE TASK GRAPH rather than listing them, and no grammar can stop `SET a = a` or a property named twice. |
| grant-ownership | 10 | `ON ACCOUNT a`: the object-type sample includes `ACCOUNT`, so the name after it is a form the docs do not write; the page's `ON ACCOUNT` takes no name. Removing `ACCOUNT` from the sample is a converter follow-up. |
| execute-dbt-project | 7 | `EXECUTE DBT PROJECT IMPORTS = (...)` with no project, or an identifier where the page takes a source location or system function; the optional `FROM WORKSPACE` alternative is kept, as the page and the parser have it. |
| grant-caller | 6 | `ON ACCOUNT a` from the object-type sample, and the inherited form with `IN APPLICATION PACKAGE` and no name. |
| revoke-caller | 5 | `ON ACCOUNT a` from the object-type sample, and the inherited form with `IN APPLICATION PACKAGE` and no name. |
| show-grants | 5 | Four `ON ACCOUNT a` rows from the object-type sample; `SHOW GRANTS TO APPLICATION ROLE` without its name; and the new `SHOW GRANTS ... INCLUDE HIERARCHY`, which the live page added and omni's parser does not read yet (a parser follow-up, not a docs fix). |
| alter-database-role | 4 | The role's name is dropped by a mutation and the statement reads as `ALTER DATABASE` named `ROLE`. |
| alter-application-package | 3 | The package's name is dropped and the statement reads as `ALTER APPLICATION` named `PACKAGE`. |
| create-task-graph | 3 | `TASK_DEFAULTS ( <default_property> = <value> )`; the page's list of supported default properties is prose, and the run may repeat one. |
| alter-application-role | 2 | The role's name is dropped and the statement reads as `ALTER APPLICATION` named `ROLE`. |
| create-interactive-warehouse | 2 | A property run `( ... )+` may name the same property twice. |
| alter-account | 1 | `UNSET TAG` without a tag: the page's generic `UNSET <param_name>` reads `TAG` as the parameter name. |
| alter-application-service | 1 | The service's name is dropped and the statement reads as `ALTER APPLICATION` named `SERVICE` and `UPGRADE`. |
| alter-dynamic-table | 1 | `[ EQUALITY ]` after `DROP SEARCH OPTIMIZATION ON <col>`; the parameters list no trailing `EQUALITY`. |
| alter-organization-account | 1 | `UNSET TAG` without a tag, as `ALTER ACCOUNT`'s generic `UNSET <param_name>`. |
| alter-organization-user | 1 | `SET EMAIL = 's' EMAIL = 's'`: the property run repeats a property. |
| alter-pipe | 1 | `UNSET TAG` without a tag, as the generic `UNSET { <property_name> | COMMENT }`. |
| alter-snapshot-policy | 1 | The policy's name is dropped and the statement reads as another `ALTER` form. |
| alter-task | 1 | `UNSET TAG` without a tag, as the generic `UNSET <property_name>`. |
| alter-user | 1 | `UNSET TAG` without a tag, as the generic `UNSET { <object_property_name> | ... }`. |
| alter-warehouse | 1 | `UNSET TAG` without a tag, as the generic `UNSET { <property_name> | <param_name> }`. |
| create-clone | 1 | The database role's name is dropped and the statement reads as `CREATE <object>` named `DATABASE ROLE`. |
| create-database-role | 1 | The role's name is dropped and `CREATE OR ALTER DATABASE ROLE` reads as `CREATE OR ALTER DATABASE` named `ROLE`. |
| create-organization-user | 1 | A property run may name the same property twice. |
| grant-application-role | 1 | The grantee's name is dropped after `APPLICATION ROLE`. |
| grant-role | 1 | `FOR INTERVAL 1`: `<interval>` is read as a number, while the parameters require a quoted `'<n> <unit>'`. |
| revoke-application-role | 1 | The grantee's name is dropped after `APPLICATION ROLE`. |
| show-semantic-dimensions | 1 | `IN SCHEMA.a` with the database part dropped; the page writes `SCHEMA <db_name>.<schema_name>`. |
| show-semantic-facts | 1 | `IN SCHEMA.a` with the database part dropped. |
| show-semantic-metrics | 1 | `IN SCHEMA.a` with the database part dropped. |
