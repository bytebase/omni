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
has 53,144 of them; the converter and parser changes below leave 53,222.

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

## What is left: 49 rows

After the pin bump the docs judged 118 rows that omni rejects. The converter
and parser changes above close 69 of them. Every remaining row is `reject` in
omni and `accept` in the docs, grouped here by the page whose notation or
reading leaves it. No alternative row that agreed with the docs stopped
agreeing; the only rows whose expected moved from `accept` to `undocumented`
are near-misses that write a name other than the sample (81), as the name
sample says above.

| Page | Rows | What remains |
| --- | --- | --- |
| alter-stage | 19 | `ALTER STAGE <name> SET` with no property, which the `statement` rule holds as `statement/alterStage` too; the external form's `SET { [ externalStageParams ] [ FILE_FORMAT ] [ COMMENT ] }` and its all-optional GCS/Azure `externalStageParams` have no notation for "at least one", so a page-local fix would invent a required parameter. |
| execute-dbt-project | 11 | `EXECUTE DBT PROJECT IMPORTS = (...)` with no project: the page's `FROM WORKSPACE` is optional, and omni does not read these targetless forms. The system-function calls are the five the parameters name, not a bare identifier. |
| alter-task-graph | 9 | `<property> = <value>`, `<default_property>` and a repeated property run. The page says the properties are the ones in CREATE TASK GRAPH rather than listing them, and no grammar can stop `SET a = a` or a property named twice. `UNSET a` is the same: `a` is not one of the page's properties. |
| create-task-graph | 3 | `TASK_DEFAULTS ( <default_property> = <value> )`; the page's list of supported default properties is prose, and the run may repeat one. |
| create-interactive-warehouse | 2 | A property run `( ... )+` may name the same property twice. |
| show-semantic-dimensions, -facts, -metrics | 3 | `IN SCHEMA.a` with the database part dropped; the page also lets `IN` take a `<semantic_view_name>`, a qualified name, and `SCHEMA.a` is one, so the docs say this is valid. |
| alter-organization-user | 1 | `SET EMAIL = 's' EMAIL = 's'`: the property run repeats a property. |
| create-organization-user | 1 | A property run may name the same property twice. |

The agree count moved from 20,486 to 20,472 and the `undocumented` count from
53,144 to 53,222; the total moved from 73,748 to 73,743. The object-type choice
lost its `ACCOUNT` alternative on four rules, the search-optimization actions
lost their trailing `EQUALITY` (with four `suite:` rows keeping the SUSPEND and
RESUME targets), and `use-schema` gained the page's qualified form
(`USE <db_name>.<schema_name>`) beside the `USE SCHEMA` one, as its parameters
say.
