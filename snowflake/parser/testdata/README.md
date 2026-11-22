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

`alter-task`'s `UNSET` list shares its placeholder with the page's `SET` form,
which keeps the class: docs-fixes names the `UNSET` occurrence
`<unset_parameter>`. Reading the sample there, and not the class, is what makes
the page's own `UNSET TAG <tag_name>` the only form that writes `TAG`.
Similarly, `show-semantic-dimensions`, `show-semantic-facts` and
`show-semantic-metrics` sample `<semantic_view_name>` as `a`: Snowflake's
[reserved and limited keywords](https://docs.snowflake.com/en/sql-reference/reserved-keywords)
say `SCHEMA` and `DATABASE` cannot be used as an identifier in `SHOW ... IN
<identifier>`, so the dropped-database near-miss `IN SCHEMA.a` must not read as
`<semantic_view_name>`.

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

## What is left: 8 rows

After the pin bump the docs judged 118 rows that omni rejects. The converter
and parser changes above closed 69 (#1459 left 49), and `execute-dbt-project`
closed 11 more (leaving 38). This pass closes 30 more docs-notation rows
without weakening `snowflake/parser`:

- **`alter-stage` (19 → 0):** Each of the three `externalStageParams` variants
  (Amazon S3, Google Cloud Storage, Microsoft Azure) now requires at least one
  parameter in `docs-fixes.tsv`, keeping every non-empty combination (and
  leaving `[ externalStageParams ]` optional beside `FILE_FORMAT` and
  `COMMENT`) while dropping the invalid empty `ALTER STAGE [ IF EXISTS ] <name>
  SET` derivation
  ([page](https://docs.snowflake.com/en/sql-reference/sql/alter-stage)).
- **`show-semantic-{dimensions,facts,metrics}` (3 → 0):** `<semantic_view_name>`
  is listed in `sampledNames` (`../cmd/docsgrammar`) per Snowflake's
  [reserved and limited keywords](https://docs.snowflake.com/en/sql-reference/reserved-keywords),
  so the `IN SCHEMA.a` near-misses move to `undocumented` while `IN a` keeps
  agreeing.
- **`alter-task-graph` & `create-task-graph` (12 → 4):** `docs-fixes.tsv`
  expands `<property>` and `<default_property>` to the 10 graph-wide properties
  and 10 supported `TASK_DEFAULTS` properties (including `EXECUTE AS USER
  <user_name>`) explicitly listed on
  [create-task-graph](https://docs.snowflake.com/en/sql-reference/sql/create-task-graph)
  and referenced by
  [alter-task-graph](https://docs.snowflake.com/en/sql-reference/sql/alter-task-graph).
  All 12 generic unknown-property rows (`SET a = a`, `UNSET a`,
  `TASK_DEFAULTS(a = a)`, and the `SET CONTACT = a` near-miss) are replaced by
  agreeing rows covering every documented property, leaving only the 4
  repeated-same-property samples below.

All 8 remaining rows are `reject` in omni and `accept` in the docs, and all 8
are the same notation artifact — repeating the identical property twice in an
unordered property or parameter list (0 near-miss disagreements):

| Page | Rows | What remains, in the page's words |
| --- | --- | --- |
| alter-task-graph | 3 | `SET SCHEDULE = 's', SCHEDULE = 's'`, `SET TASK_DEFAULTS(WAREHOUSE = a, WAREHOUSE = a)`, and `UNSET TASK_DEFAULTS(WAREHOUSE, WAREHOUSE)`: the comma-separated repetition of the 10 documented properties allows any unordered combination, so the generator samples its `*` branch by repeating the first property twice ([page](https://docs.snowflake.com/en/sql-reference/sql/alter-task-graph)). |
| create-task-graph | 1 | `TASK_DEFAULTS(WAREHOUSE = a, WAREHOUSE = a)`: the same comma-separated repetition over the 10 supported default properties ([page](https://docs.snowflake.com/en/sql-reference/sql/create-task-graph)). |
| create-interactive-warehouse | 2 | `objectProperties` and `objectParams` are one `NAME = value` per line, which the converter reads as any number of them, so the run may name the same property twice (`WAREHOUSE_SIZE = XSMALL WAREHOUSE_SIZE = XSMALL`, `MAX_CONCURRENCY_LEVEL = 1 MAX_CONCURRENCY_LEVEL = 1`) ([page](https://docs.snowflake.com/en/sql-reference/sql/create-interactive-warehouse)). |
| alter-organization-user, create-organization-user | 2 | The same `objectProperties` run, one row for each page: `EMAIL = 's' EMAIL = 's'` is the converter's sample of a repeat ([alter](https://docs.snowflake.com/en/sql-reference/sql/alter-organization-user), [create](https://docs.snowflake.com/en/sql-reference/sql/create-organization-user)). |

For all 8 rows, the pages' examples never repeat a property, but their
parameters do not state whether a duplicate is accepted or rejected either.
Rewriting those lists in `docs-fixes.tsv` to forbid repeats would fix an
arbitrary property order and stop the docs deriving unordered combinations omni
accepts (`STATEMENT_QUEUED_TIMEOUT_IN_SECONDS` before `MAX_CONCURRENCY_LEVEL`
and `RESOURCE_MONITOR` before `INITIALLY_SUSPENDED` in
`TestParseCreateInteractiveWarehouse`, `EXECUTE AS USER` before `WAREHOUSE` in
`TestParseCreateTaskGraph`, `LOG_LEVEL` before `COMMENT` in
`TestParseAlterTaskGraph`), while relaxing omni's duplicate-property checks
without a live `SNOWFLAKE_COMPAT_DSN` engine oracle would violate `AGENTS.md`.
So these 8 rows stay until an engine oracle settles them.

The agree count moved from 20,483 / 20,521 (`99.81%`) to 20,569 / 20,577
(`99.96%`); the total moved from 73,743 to 74,103 (`undocumented` from 53,222
to 53,526) while the docs pin stays at the same `2026-11-18` digest. No
`expected` cell was edited by hand, and no alternative row that agreed with the
docs stopped agreeing.
