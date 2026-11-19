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
| One per alternative | 19,775 | 19,638 | 137 | 19,800 | 19,728 | 72 |
| Near-misses | 53,876 | 758 | 49 | 53,948 | 758 | 46 |
| All | 73,651 | 20,396 | 186 | 73,748 | 20,486 | 118 |

`undocumented` rows are the near-misses the docs do not derive; docs say what is
valid, not what Snowflake refuses, so they pass whatever omni says. The new pin
has 53,144 of them; the converter and parser changes below leave 53,289.

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

The object-type sample read `ACCOUNT` first, so `GRANT OWNERSHIP ON ACCOUNT
a`, `GRANT`/`REVOKE CALLER PRIVILEGES ON ACCOUNT a` and `SHOW GRANTS ON
ACCOUNT a` were rows. The converter now reads `<object_type>` on the pages
whose syntax writes it before a name as a type that takes one: `ACCOUNT` is
not among them, because those pages write `ON ACCOUNT`, or `ON ALL
<object_type_plural> IN ACCOUNT`, as a form of their own, and `GRANT
OWNERSHIP`'s `object_type` parameter lists no `ACCOUNT`. `SCHEMA` comes first,
so the statement a row derives is `ON SCHEMA a`. See `objectTypeTakesName` and
`namedObjectTypeChoice` in `../cmd/docsgrammar`.

### The name sample: the page's own form has the name

The docs judge reads a name class (`ID`) as any word, keywords included,
because Snowflake reserves so few of them. So a form of one page with its name
dropped reads as another page's statement: `DROP APPLICATION PACKAGE` is `DROP
APPLICATION` named `PACKAGE`, `USE ROLE` is `USE` named `ROLE`, and `ALTER
ACCOUNT UNSET TAG` is the generic `UNSET <param_name>`. Omni's parser takes the
longer form (it reads `PACKAGE`, `ROLE` and `TAG` as the object word and then
wants the name), and each page writes that word as a form of its own.

`../cmd/docsgrammar` writes the name of the pages that own such a form by the
sample `a`, the text `tokens.tsv` gives `ID`, instead of by the name class
(`sampledNames` there lists each page and placeholder, with the page's URL).
The row's statement is the one the class derived (`DROP APPLICATION a`), and a
name-dropped near-miss is `undocumented` instead of another page's statement.
That is a narrowing: a near-miss that writes a name other than the sample is
`undocumented` too, which passes. No alternative row moved to `undocumented`
here; the alternatives keep their statements and still agree.

Two pages' `UNSET` lists share their placeholder with the page's `SET` form,
which keeps the class: docs-fixes names the `UNSET` occurrence
`<unset_parameter>` (`alter-task`) and `<unset_property>` (`alter-task-graph`).
Reading the sample there, and not the class, is what makes the page's own
`UNSET TAG <tag_name>` the only form that writes `TAG`.

`../../testdata/tokens.tsv` samples every token class the grammar uses.

The parser reads the page's `SHOW GRANTS ... INCLUDE HIERARCHY` now
(`../show.go`, `ast.ShowStmt.IncludeHierarchy`), for `SHOW GRANTS TO` a role, a
database role or a user, and for `SHOW FUTURE GRANTS TO` a role or a database
role, the grantee forms the page writes the clause after:
https://docs.snowflake.com/en/sql-reference/sql/show-grants.

This PR's fixes are in `docs-fixes.tsv` beside the pin bump's: `FOR INTERVAL
'<interval>'` on `grant-role` (the parameter writes a quoted `'<n> <unit>'`),
the two `[ EQUALITY ]` tails of `alter-dynamic-table`'s search optimization
actions (the parameter gives `EQUALITY` a target, and no `EQUALITY` follows
one), and the two `UNSET` names above. Dropping those tails makes DROP's `ON`
list structurally identical to SUSPEND's and RESUME's, and the generator then
attributes the shared alternatives to DROP alone; four `suite:` rows hold the
targeted SUSPEND and RESUME forms the generator stopped writing.

To rebuild the grammar from the live docs (a pin bump), then the rows:

```bash
go run ./snowflake/parser/cmd/docsgrammar -cache /tmp/snowflake-docs
OMNI_RECORD=1 go test ./snowflake/parser/ -run '^TestCompatibility$'
```

After a pin bump or a converter change the row set itself moves; record it, then
commit the verdicts the parser change closed:

```bash
OMNI_RECORD=1 go test ./snowflake/parser/ -run '^TestCompatibility$'
go test ./snowflake/parser/ -run '^TestCompatibility$' -update-compat
```

## What is left: 31 rows

After the pin bump the docs judged 118 rows that omni rejects. The converter
and parser changes above closed 69; #1459 left 49. This pass closes 18 more:

- **`execute-dbt-project` (11).** The page writes two forms, and the second
  brackets `FROM WORKSPACE` and still lists `IMPORTS`, so
  `EXECUTE DBT PROJECT IMPORTS = (...)` with no project at all is the page's
  shape. `dbtExecuteWorkspaceProps` now takes `IMPORTS` beside `PROJECT_ROOT`.
- **`show-semantic-dimensions`, `-facts`, `-metrics` (3).** The pages' `IN`
  choice takes a `<semantic_view_name>`, a QNAME, before the `ACCOUNT`,
  `DATABASE` and `SCHEMA` scopes. Snowflake does not reserve `SCHEMA`, so
  `IN SCHEMA.a` is a semantic view named `a` in database `SCHEMA`, not the
  `SCHEMA` scope with a stray `.a`; `parseShowScope` reads a scope keyword
  followed by `.` as the first part of the name.
- **The property runs (4).** `create-interactive-warehouse`'s
  `objectProperties` and `objectParams`, and `create-` and
  `alter-organization-user`'s `objectProperties`, list their properties one
  per line, and the converter read that as `( ... )+`: any number of them,
  so a sampled run could name the same property twice. The page's examples
  never repeat one and its parameters say nothing about setting one twice, so
  `docs-fixes.tsv` rewrites each list as the alternatives that start at each
  property and take the later ones optionally, in the page's order, as the
  ALTER RESOURCE MONITOR fix already does. No parser change was made to fit
  the samples.

The remaining rows are all `reject` in omni and `accept` in the docs, grouped
by the page whose notation or reading leaves them:

| Page | Rows | What remains, in the page's words |
| --- | --- | --- |
| alter-stage | 19 | `ALTER STAGE [ IF EXISTS ] <name> SET { [ externalStageParams ] [ FILE_FORMAT = ( ... ) ] [ COMMENT = '<string_literal>' ] }`, and the external stage's Amazon S3, Google Cloud Storage and Microsoft Azure `externalStageParams` are all optional too, as the page writes them (the Amazon S3 one after docs-fixes lets an ALTER change the credentials or encryption and keep the URL). Nothing says "at least one", so `ALTER STAGE a SET` is the page's own shape; a page-local fix would invent a required parameter ([page](https://docs.snowflake.com/en/sql-reference/sql/alter-stage)). The `statement` rule holds the same form as `statement/alterStage`. |
| alter-task-graph | 9 | `SET <property> = <value> [ , <property> = <value> ... ]`, `UNSET <property> [ , ... ]`, `SET TASK_DEFAULTS ( <default_property> = <value> [ , ... ] )` and its UNSET. The page says the properties are "the ones listed in CREATE TASK GRAPH" and lists the defaults only in prose, so its notation derives any name, `SET a = a` and a repeated property alike; omni checks the names, which no grammar there can ([page](https://docs.snowflake.com/en/sql-reference/sql/alter-task-graph)). |
| create-task-graph | 3 | `TASK_DEFAULTS ( <default_property> = <value> [ , ... ] )`: the page's list of supported default properties is prose in its TASK_DEFAULTS parameter, so `a = a` and a repeated one are the notation's, not a stated form ([page](https://docs.snowflake.com/en/sql-reference/sql/create-task-graph)). |

The agree count moved from 20,472 to 20,529 and the `undocumented` count from
53,222 to 53,289; the total moved from 73,743 to 73,849, because the four
property-list rules now have one alternative per property instead of the two
a `( ... )+` run had. The pin is the same: `docs-fixes.tsv` changed, the pages
did not. No alternative row that agreed with the docs stopped agreeing.
