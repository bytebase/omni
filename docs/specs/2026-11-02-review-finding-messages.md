# Review finding messages: standard and style

Status: accepted 2026-11-03. Applies to `review.Finding.Message` in every engine's `Review` (`pg`, `mysql`, `mariadb`, `tidb`, `oceanbase`, `mssql`, `oracle`, `snowflake`), and states what Bytebase shows around it.

Evidence: a survey of 170 changes on a live Bytebase dev build (omni `a4ffab6b`) across PostgreSQL 17, MySQL 8.0, MariaDB 11.8, TiDB 8.5, SQL Server 2022, and Oracle 23. It includes screenshots of every finding and a before/after demo of Markdown messages rendered by Bytebase's real comment renderer. The findings that drive this spec are summarized in [Appendix A](#appendix-a-what-the-survey-found).

## 1. Who reads a finding, and what they need

A finding becomes an issue comment, anchored to the statement it is about. Four readers meet it:

| Reader | When | What they do next |
| --- | --- | --- |
| Author of the change | Right after submitting, and after each edit | Fix the SQL, or explain why it stays |
| Reviewer | While approving | Accept the risk knowingly, or send it back |
| Whoever deploys or bypasses | When the rollout is held | Decide whether holding is right |
| Whoever audits later | After an incident | Learn what was accepted and why |

P0 findings speak mainly to the author: what will fail, and how to fix it. P1 findings speak mainly to the reviewer: exactly what the change does that needs someone's acceptance.

Every reader comes with the same five questions:

1. **What is the problem?** Which object, and which part of the statement.
2. **What happens?** The statement fails, data is lost, the backup is not taken, the deploy is held.
3. **What do I do?** Change the SQL, turn a setting off, or have a reviewer accept it.
4. **Can I trust it?** Is it certain, and what is it based on?
5. **Where, and on which database?** Bytebase answers this with the statement anchor and the "Applies to" targets.

A good finding lets a reader answer the first three questions without opening the rule's documentation or asking anyone.

## 2. The standard

Each criterion is ordered by importance and comes with a check.

| # | Criterion | Check |
| --- | --- | --- |
| 1 | **True.** Every claim is literally true on the engine. | The engine oracle confirms it. "Table `v` does not exist" for a view fails this check. |
| 2 | **Understandable.** Uses only words the reader has seen: in their SQL, in the engine's documentation, or in the Bytebase UI. | Ask of each noun where the reader met it. `gh-ost` internals, "the task's database", "review unit", and PostgreSQL's "relation" fail this check. |
| 3 | **States the consequence** for this change, not the rule's rationale. | A reader who does not know the rule can tell why it was flagged. |
| 4 | **Actionable.** The next step is clear: an edit, a setting, or a reviewer's acceptance. | The reader can name their next action. |
| 5 | **Matches what Bytebase does with it.** P0 wording means the deploy is held; P1 wording means a person must accept it. | No P1 text says "disallowed"; no P0 text reads like a suggestion. |
| 6 | **One problem, one finding.** | Fixing one root cause closes exactly one finding. |
| 7 | **Stands alone.** Readable in the activity timeline, a notification, or an API response, without the SQL beside it. | Read the findings as a bare list; each is still clear. |
| 8 | **Consistent.** The same problem reads the same across rules and engines. | Run one scenario on every engine and diff the messages. |
| 9 | **Concise.** The first line carries the point. | The first eight or so words name the problem. |
| 10 | **Neutral.** No blame, no alarm. A risk that needs acceptance is not called forbidden. | Read it aloud to the author. |

When criteria conflict, the lower number wins. A message that is shorter but less true fails.

## 3. Who carries what

A finding is shown with Bytebase's own framing. The message must not repeat it. The framing is localized; the message is English only, by decision (§8).

| Shown by Bytebase | Carried by the message |
| --- | --- |
| Rule title and its general description (localized) | The specific problem in this change |
| Priority, and what it means for the deploy | The objects involved, named as the SQL names them |
| The anchored statement and line | The consequence specific to this case |
| The targets the finding applies to | The next step specific to this case |
| Author, time, thread state | |

The message never names:

- a target database: the existing contract, so findings from several targets merge;
- the rule, or its priority;
- the general reason the rule exists, which belongs to the rule description.

## 4. Format

### 4.1 Structure

```
<summary line>

<body: optional, at most three short paragraphs, or one list>
```

- **Summary line.**
  - Required.
  - One sentence. Sentence case. Ends with a period.
  - States the problem and names the object.
  - Aim for 100 characters or fewer, not counting the contents of code spans.
  - It may contain code spans and nothing else from Markdown. Bytebase shows this line on its own in places, such as a collapsed row and a finding title.
- **Body.** Optional. A blank line separates it from the summary. The body holds:
  - the consequence specific to this case, when the summary does not already say it;
  - the next step, when one exists: what to change, or which setting to turn off;
  - for several causes of one problem, a bullet list with one cause per item, at most five items;
  - for SYNTAX, a fenced code block showing the statement line with a caret under the error.

### 4.2 Markdown subset

Bytebase renders review comments as untrusted Markdown (markdown-it, no raw HTML, images become links, bare hosts are linkified). A message may use only:

| Allowed | For |
| --- | --- |
| Code spans | Every piece of SQL-derived text (§4.3) |
| Fenced code blocks | SYNTAX snippets only |
| Bullet lists | Several causes of one problem |
| **Bold** | The name of a Bytebase setting the reader should change: **Prior backup**, **Online migration** |

No headings, links, tables, images, italics, block quotes, or HTML.

### 4.3 SQL-derived text goes in code spans

Every piece of text taken from the SQL, the synced schema, or the parser goes in a code span. That covers identifiers, qualified names, expressions, quoted keywords, data types, and parser tokens. This is a correctness rule, not a style rule. Plain text from the SQL is rendered as Markdown today and comes out wrong:

| Message today | Rendered by Bytebase |
| --- | --- |
| `drops table _tmp_` | "drops table *tmp*": underscores lost |
| `truncates retail.shop` | a link to `http://retail.shop` |
| `drops user a@example.com` | a mailto link |
| `"a*b = NULL OR c*d = NULL" is never ...` | emphasis between the asterisks |
| a SYNTAX text that quotes the rest of the script | headings and lists from `#` and `-` lines |

Use one helper, `review.Code(s string) string`, and never build code spans by hand. The helper:

- picks a backtick fence longer than any backtick run in `s`;
- pads with a space when `s` starts or ends with a backtick;
- writes line breaks and tabs as `\n` and `\t`;
- truncates to 120 characters with `…`.

Write SQL keywords that are not quoted from the statement in upper case inside a code span when they name syntax: `` `WHERE` ``, `` `RENAME COLUMN` ``, `` `PRIMARY KEY` ``.

### 4.4 Naming objects

- **In a form the engine reads back as the same object.** A name appears:
  - qualified exactly as the SQL qualified it;
  - quoted only when the engine needs the quotes to read it back.

  Each engine renders an unquoted name its own way: PostgreSQL folds it to lower case and Oracle to upper case, as their catalogs store it, while the others keep the SQL's case. Every one of these forms reads back as the same object. Writing every name exactly as typed would mean re-reading the source for each one, since most names carry no source position.
- **As the engine stores it.** An object the SQL does not name, such as a foreign key found through the synced schema, appears as the engine stores it.
- **Kind words.** Use one term per kind across engines:
  - `table`, `view`, `materialized view`, `column`, `index`, `sequence`, `schema`, `database`, `function`, `procedure`, `trigger`;
  - `primary key`, `foreign key`, `unique constraint`, `check constraint`.

  MySQL-family docs say "unique key", but "unique constraint" is the term every other engine uses and MySQL readers understand it. When the kind is unknown, omit the kind word rather than guess: write `` `order_items` doesn't exist ``.
- **The owning object.** Name a column or constraint together with its table: ``column `status` of `orders` ``.

### 4.5 Words

| Use | Not |
| --- | --- |
| **Online migration** (the setting), "gh-ost" only when the limitation is gh-ost's | the gh-ost directive, `approve-renamed-columns`, flag names, unless the user wrote them |
| **Prior backup** (the setting); the backup database or schema by name, in a code span | "the backup", without saying which; `bbdataarchive` without saying it is the backup location |
| "this change", "this statement" | task, spec, sheet, unit, walk, snapshot, ledger |
| "the schema Bytebase last synced" | "the metadata", "the snapshot" |
| table, view, ... (§4.4) | relation, object (when the kind is known) |
| "fails", "is rejected" (P0 predictions) | "error", "invalid", with no subject |
| "drops", "renames", "truncates" (P1 facts) | "disallowed", "forbidden", "not allowed" |

Use American English. Contractions such as "doesn't" and "can't" are fine. Write numbers as digits, and sizes with binary units (`10 MiB`, not `10485760-byte`). Use correct plurals.

### 4.6 By priority

**P0** (SYNTAX, WALK_THROUGH, ONLINE_MIGRATION, PRIOR_BACKUP, REQUIRE_IS_NULL):

- The summary says what fails, or what the reader asked for that will not happen.
- The body says how to fix it, or which setting to turn off.
- Never say "must": a P0 finding can be resolved without a change, and the deploy then proceeds (decision 1, §8). State the fact and the fix.
- A WALK_THROUGH claim is conditional on the earlier statements running. When an earlier statement of the change caused the failure, name it: "line 2 drops it".
- No error code in the message text. Finding codes are a separate field (decision 5, §8, pending).

**P1** (REQUIRE_WHERE, DISALLOW_*, REQUIRE_PRIMARY_KEY):

- The summary states exactly what the change does to which object. That is what the reviewer accepts.
- The body is optional, and only for scope this case adds that the summary cannot say, such as what `CASCADE` also drops. It must be true in every case the template covers: "so it changes every row" is false for a joined `UPDATE` or one with `LIMIT`, so no template says it.
- The general risk of the operation is the rule description's job, not the message's.

### 4.7 One problem, one finding

- **Several reasons, one finding.** When one rule finds several reasons why the same request cannot be met, it reports one finding that lists them. This applies to ONLINE_MIGRATION and PRIOR_BACKUP. Anchor it on the statement when there is one, else on the whole change.
- **Overlapping rules: precedence.** When two rules describe the same subcommand, only the more specific rule reports it, if that rule is enabled. Otherwise the other rule reports it.
  - `DROP PRIMARY KEY`, or dropping the primary key constraint: REQUIRE_PRIMARY_KEY reports it when the change adds no primary key back, else DISALLOW_DROP_CONSTRAINT.
  - MySQL-family `DROP INDEX` of a unique key: DISALLOW_DROP_CONSTRAINT reports it, not DISALLOW_DROP_OBJECT. This is the one the proto's rule description already promises.
- **Repeated statements.** A repeated statement, such as `DROP TABLE t1` twice, gets one P1 finding per statement. The WALK_THROUGH finding on the second one says why it fails.

## 5. Examples

Before: the message omni reports today. After: what this spec asks for. Rendered examples are in the survey's Markdown demo page.

### SYNTAX

Before: `syntax error at or near ")"`

After:

````
Syntax error at line 2, column 24, near `)`.

```
CREATE TABLE t (id int,);
                       ^
```

The other rules check the change once it parses.
````

Line and column come from `review.Position`. The anchor is zero-width, and Bytebase highlights only the whole line, so the message carries the column.

When the parser knows that the construct is valid SQL elsewhere but not on this engine, say so instead of naming a token:

- Before: `expected identifier` (MariaDB, `DROP CHECK`)
- After: ``MariaDB doesn't support `DROP CHECK`. Use `DROP CONSTRAINT orders_amount_check`.``

### WALK_THROUGH

| Before | After |
| --- | --- |
| `relation order_items does not exist` (pg) | ``` `order_items` doesn't exist, so this `INSERT` fails. ``` |
| `Table big_orders does not exist` (mysql, `TRUNCATE` of a view) | ``` `big_orders` is a view, not a table, so `TRUNCATE TABLE` fails on it. ``` |
| `Table t1 does not exist` (second `DROP TABLE t1`) | ``` Table `t1` no longer exists here: line 2 drops it. ``` |
| `Column STATUS being added already exists in table ORDERS` (oracle) | ``` Column `status` already exists in `orders`. ``` |
| `an object named dbo.customers already exists` (mssql) | ``` `dbo.customers` already exists. ``` |

### ONLINE_MIGRATION

Before, three findings:

- `online migration needs database bbdataarchive, which does not exist`
- `online migration does not support orders, which has a foreign key`
- `gh-ost does not read RENAME COLUMN as a rename, and the migrated table loses the data of column status`

After, one finding:

```
Online migration can't run this change.

- The backup database `bbdataarchive` doesn't exist on this instance.
- `orders` has a foreign key, and gh-ost can't migrate a table with foreign keys.
- gh-ost doesn't read `RENAME COLUMN` as a rename, so column `status` would lose its data.

Turn off **Online migration**, or fix these and rerun the review.
```

On an engine that has no online migration:

- Before:
  - `the change requests online migration, which PostgreSQL does not support`
  - TiDB: `the change requests online migration, which runs gh-ost, and TiDB has no binary log for gh-ost to read`
- After: ``PostgreSQL doesn't support **Online migration**.``, with the body `Turn it off for this change.`
- For TiDB, add one body line giving the reason: `gh-ost needs a binary log, which TiDB doesn't have.`

### PRIOR_BACKUP

**Backup location missing**

Before: `prior backup needs schema bbdataarchive, which does not exist`

After:

```
Prior backup can't run: this database has no `bbdataarchive` schema to store the backup in.

Create the schema `bbdataarchive`, or turn off **Prior backup** for this change.
```

**UPDATE and DELETE on one table**

Before: `prior backup cannot back up both UPDATE and DELETE on orders`

After:

```
Prior backup can't back up `orders`: this change both updates and deletes its rows.

Move the `UPDATE` and the `DELETE` into separate changes, or turn off **Prior backup**.
```

**Nothing to back up (Oracle)**

Before: `prior backup copies the rows of top-level UPDATE and DELETE statements on a table, and the change has none, which fails the backup`

After:

```
Prior backup fails on this change: it has no `UPDATE` or `DELETE` whose rows it can copy.

Turn off **Prior backup**.
```

### REQUIRE_IS_NULL

Before: `"status = NULL" is never true; use IS NULL`

After:

```
`status = NULL` is never true.

Use `IS NULL` instead.
```

"so the `WHERE` clause matches no rows" would be false under `OR`, so the message does not say it.

### REQUIRE_WHERE

| Before | After |
| --- | --- |
| `UPDATE orders has no WHERE clause` | ``` `UPDATE` of `orders` has no `WHERE` clause. ``` |
| `DELETE FROM logs has no WHERE clause` | ``` `DELETE` from `logs` has no `WHERE` clause. ``` |
| `DELETE FROM t, p has no WHERE clause` (MySQL multi-table) | ``` `DELETE` from `t` and `p` has no `WHERE` clause. ``` |

### DISALLOW_DROP_OBJECT, DISALLOW_TRUNCATE

| Before | After |
| --- | --- |
| `drops table logs` | ``` Drops table `logs`. ``` |
| `drops table t, p` | ``` Drops tables `t` and `p`. ``` |
| `drops column status of orders` | ``` Drops column `status` of `orders`. ``` |
| `drops index orders_status_idx on dbo.orders` | ``` Drops index `orders_status_idx` of `dbo.orders`. ``` |
| `truncates table LOGS` (oracle) | ``` Truncates `LOGS`. ``` |

### DISALLOW_DROP_CONSTRAINT, DISALLOW_RENAME

| Before | After |
| --- | --- |
| `drops foreign key orders_customer_fk of orders` | ``` Drops foreign key `orders_customer_fk` of `orders`. ``` |
| `drops unique key customers_email_key of customers` + `drops index customers_email_key of customers` | one finding: ``` Drops unique constraint `customers_email_key` of `customers`. ``` |
| `renames column status of orders to state` | ``` Renames column `status` of `orders` to `state`. ``` |

### REQUIRE_PRIMARY_KEY

**Creating a table without a primary key**

- Before: `creates table events without a primary key`
- After: ``Creates table `events` without a primary key.``

**Dropping the primary key**

Before, two findings:

- `drops primary key customers_pkey of customers`
- `removes the primary key of customers, and the change adds none back`

After, one finding:

```
Drops the primary key `customers_pkey` of `customers`, and the change doesn't add one back.

`CASCADE` also drops foreign key `orders_customer_fk` of `orders`, which references it.
```

The `CASCADE` line needs new logic. Report it only when the synced schema shows the dependent object.

## 6. omni implementation

- **Contract.** Replace the `Finding.Message` doc in `review/review.go` ("one line naming the object as the SQL wrote it") with a pointer to the rules here. The new doc says:
  - the message is Markdown in the subset of §4.2;
  - the first line is the summary;
  - SQL-derived text goes through `review.Code`;
  - the message never names a target.
- **Helpers in `review`.** Add `Code`, which escapes text into a code span (§4.3). Optionally add a small builder for "summary, body, list". Engines must not build messages from raw `fmt.Sprintf` of names.
- **One message table per rule, shared shape across engines.** Each engine keeps its own detection, but a rule's message templates follow the same wording in every engine. Change a rule's messages in all engines in one PR, forks included.
- **Tests.**
  - A hostile-name test per engine. It runs each rule's cases with identifiers such as `_t_`, `a.shop`, `` `x` ``, names containing a line break, and expressions containing `*`. It asserts that every SQL-derived piece of the message sits inside a code span: render it with a CommonMark library in the test only, and compare the code-span text with the input.
  - A cross-engine scenario test. The same scenario on every engine yields the same message modulo names, with documented exceptions.
  - Existing expectations change with the wording. Update them in the same PR as the template.
- **Order of work.**
  1. Contract, helper, and the escaping change (correctness).
  2. Merging and precedence (§4.7).
  3. Rewording rule by rule, P0 rules first.

## 7. What Bytebase is expected to show

These are Bytebase changes. omni does not depend on them, but the message format assumes them.

1. **Rule framing.** Under each finding, show the localized rule title, and what its priority means for the deploy:
   - P0: holds the deploy until the SQL is fixed or the finding is resolved;
   - P1: needs a reviewer's acceptance.

   Offer the rule's description, also localized, on demand.
2. **Summary line outside the renderer.** Wherever only the summary line is shown, such as the collapsed stacked-thread row (`ThreadStack.tsx`), render it as inline Markdown, or strip the Markdown. Today it shows raw backticks.
3. **No duplicate reports.** Under SQL Review V2, hide V1's Statement report when it repeats a SYNTAX finding, and the gh-ost dry run when it repeats ONLINE_MIGRATION.
4. **Failures in plain words.** A failed review run currently shows `1 of 1 review units failed: projects/.../SHOP: metadata not synced`. Show the reason in plain words and say what to do, for example "Bytebase hasn't synced `SHOP` yet. Sync the database, then rerun the review."
5. **Bypass and header copy.** The bypass sheet and the header popover should count findings by priority, and not call them discussions or threads.

## 8. Decisions and open questions

Decided 2026-11-02:

1. **A P0 finding can be resolved without a change.** It stays possible, and the deploy then proceeds. P0 wording states what fails and how to fix it; it never says "must".
2. **No localization.** Messages stay English free text. There are no template IDs or arguments in the contract.
3. **Precedence falls back.** When the rule §4.7 prefers is switched off, the other rule reports the subcommand.
4. **P1 findings stay on a statement that will fail.** A WALK_THROUGH finding does not suppress P1 findings on the same statement. Deciding when the P1 is moot is too hard to get right.

5. **Finding codes.** Bytebase has its own finding codes in place of engine error codes. The scheme is in §9.
6. **SYNTAX hints.** A failure names its line, column, and token everywhere. A hint is added where a known mistake is confirmed (§10). PostgreSQL is the prototype.

## 9. Finding codes

Each finding carries a code that names the specific condition it reports. The code is a field, never part of the message text.

- **Two levels.** The rule, one of the 11 switches, stays the unit a project turns on and off. Under it, a code names one condition: `WT101` is "object doesn't exist".
- **One code per way of handling it.** Conditions a reader handles the same way share a code; the message names the object kind and the cause. For a P1 rule the reader weighs a risk, so a rule has several codes only where the risk differs: DISALLOW_DROP_OBJECT separates dropping data, dropping a definition the change could create again, and dropping access. WALK_THROUGH keeps its hundreds digit as the condition (1 doesn't exist, 2 name taken, 3 another kind, 4 other), so a finer code can be added later without renumbering. Appendix C lists the 40 codes.
- **Format.** A code is a two-letter rule prefix and three digits:
  - `SY` SYNTAX, `WT` WALK_THROUGH, `OM` ONLINE_MIGRATION, `PB` PRIOR_BACKUP;
  - `NU` REQUIRE_IS_NULL, `WH` REQUIRE_WHERE;
  - `DO` DISALLOW_DROP_OBJECT, `TR` DISALLOW_TRUNCATE, `DC` DISALLOW_DROP_CONSTRAINT, `RN` DISALLOW_RENAME;
  - `PK` REQUIRE_PRIMARY_KEY.

  Each code also has a readable name, such as `WT101 object-not-found`.
- **Engine-neutral.** A code means the same condition on every engine. MySQL 1146, PostgreSQL 42P01, ORA-00942, and SQL Server 208 are all `WT101`. A condition only one engine has gets its own code only when its handling differs, and its documentation names the engine.
- **Stable.** A published code is never renumbered or reused. A code no longer reported stays documented as such.
- **Merged findings.** A finding that lists several reasons (§4.7) carries one code per reason. Reasons sharing a code give it once.
- **Where it lives.**
  - omni: `Finding.Codes`, at least one per finding, each of the finding's rule; the table in `review/codes.go` (`Codes()`, `FindingCode.Info`), which is the source of truth for Appendix C. An engine reports a code, and the rule follows from it.
  - Bytebase stores the code in the comment's review metadata and shows it as a chip that links to the code's page.
  - The page has the explanation, examples per engine, the fix, and the engine error codes it corresponds to.

## 10. Syntax hints

Deferred: messages for the other rules come first. A PostgreSQL prototype of what follows exists.

- **Everywhere.**
  - The summary names the line, column, and token: "Syntax error at line 2, column 24, near `)`."
  - A code block shows the line with a caret under the error.
- **Hints.** Each engine keeps a table of known mistakes, such as MySQL or SQL Server syntax, trailing commas, and reserved words used as names. A hint is offered only when:
  - its pattern matches at the error, read from the statement's tokens;
  - the statement, repaired the way the hint says, parses.
- **Precision.** The parser's confirmation is only as good as the parser. Each hint therefore has engine tests:
  - the engine refuses the statement the hint fires on;
  - the engine accepts the repaired statement;
  - the hint stays silent on look-alike failures.

  A pattern also refuses a repair the parser is known to over-accept. For example, `TOP` gives no hint when the query already has `LIMIT`, because omni's PostgreSQL parser takes a second `LIMIT` that PostgreSQL refuses.
- **Codes.** Hints get `SY2xx` codes: `SY201` reserved word, `SY202` trailing comma, `SY203` `TOP`, `SY204` backticks, `SY205` `AUTO_INCREMENT`.

## Appendix A: what the survey found

- **Plain text is rendered as Markdown.** Underscores in names become emphasis, `retail.shop` becomes an external link, and SQL keywords lose their markers (§4.3).
- **Two voices.** MySQL-family and Oracle WALK_THROUGH messages are capitalized engine-style sentences; every other message is a lower-case fragment with no subject or period.
- **Only REQUIRE_IS_NULL says what to do.** P1 messages restate the statement.
- **Duplicates.**
  - Dropping a primary key gives two findings in every engine.
  - A MySQL-family unique-key drop gives two.
  - One gh-ost request can give three ONLINE_MIGRATION findings plus a DISALLOW_RENAME on the same statement.
- **Internal terms.** `bbdataarchive` with no explanation, "the task's database", `relation`, and gh-ost details.
- **Untrue or misleading.** `Table big_orders does not exist` for a view (MySQL family).
- **Name rendering varies by engine.** Oracle upper-cases names the SQL wrote in lower case.
- **SYNTAX is the raw parser text.** It has no column, and some messages name nothing (`expected identifier`).

## 11. Rollout

Each wave is one PR across every engine, forks included.

1. REQUIRE_IS_NULL, REQUIRE_WHERE, DISALLOW_DROP_OBJECT, DISALLOW_TRUNCATE, DISALLOW_DROP_CONSTRAINT, DISALLOW_RENAME, and REQUIRE_PRIMARY_KEY. This wave rewords them to Appendix B and adds the `review` helpers.
2. Precedence and merging between overlapping rules (§4.7).
3. ONLINE_MIGRATION and PRIOR_BACKUP, with several reasons merged into one finding.
4. WALK_THROUGH.
5. Finding codes (§9) and SYNTAX (§10).

## Appendix B: phrases for wave 1

These phrases are exact; engines differ only in the names.

**Helpers in `review`.**

- `Code(s)` writes a name or SQL text as a code span.
- `Names(names...)` writes each name as a code span and joins them: `` `a` ``, `` `a` and `b` ``, `` `a`, `b`, and `c` ``.
- `Join(items...)` joins phrases that are already written, the same way.
- `Plural(kind, n)` gives the kind's plural when n is not 1, changing only its last word: `tables`, `indexes`, `policies`, `operator classes`, `statistics`.

**Shape of every summary.** A capitalized sentence ending with a period. Every name goes through `Code`, `Names`, or `Join`. A wave-1 message has no body.

| Rule | Shape | Examples |
| --- | --- | --- |
| REQUIRE_IS_NULL | `` `<comparison>` is never true.`` + body ``Use `IS NULL` instead.`` or ``Use `IS NOT NULL` instead.`` | ``` `status = NULL` is never true. ``` / ``` Use `IS NULL` instead. ``` |
| REQUIRE_WHERE | ``` `UPDATE` of <Names> has no `WHERE` clause. ``` / ``` `DELETE` from <Names> has no `WHERE` clause. ``` | ``` `UPDATE` of `orders` has no `WHERE` clause. ```, ``` `DELETE` from `t` and `p` has no `WHERE` clause. ```, ``` `UPDATE` of an inline view has no `WHERE` clause. ``` |
| DISALLOW_DROP_OBJECT | ``Drops <Plural(kind)> <Names>.`` | ``Drops table `logs`.``, ``Drops tables `t` and `p`.``, ``Drops function `f(integer, text)`.`` |
| | with an owner: ``Drops <kind> <Names> of <Code(owner)>.`` | ``Drops column `status` of `orders`.``, ``Drops partitions `p0` and `p1` of `t`.``, ``Drops period `p` of `t`.`` |
| | several owners: ``Drops <Plural(kind)> <Join(`name` of `owner`, ...)>.`` | ``Drops indexes `ix` of `dbo.t` and `ix` of `dbo.p`.`` |
| | no name: ``Drops the <kind>.`` | ``Drops the database.``, ``Drops the master key.`` |
| | fixed phrases | ``Drops the objects owned by `r`.``, ``Drops the user mapping for `r` on server `s`.``, ``Drops the cast from `integer` to `text`.``, ``Drops trigger `tr` on `t`.``, ``Drops operator class `c` using `btree`.``, ``Drops system versioning of `t`.``, ``Drops the unused columns of `t`.``, ``Drops the columns an interrupted drop left in `t`.``, ``Drops the sensitivity classification of `c1` and `c2`.``, ``Drops the signature of `m`.`` |
| DISALLOW_TRUNCATE | ``Truncates <Names>.`` | ``Truncates `logs`.``, ``Truncates `t` and `p`.``, ``Truncates materialized view `mv`.`` |
| | partitions | ``Truncates partitions `p0` and `p1` of `t`.``, ``Truncates all partitions of `t`.``, ``Truncates partitions of `dbo.p`.``, ``Truncates the error table of `t`.`` |
| DISALLOW_DROP_CONSTRAINT | ``Drops <kind> <Code(name)> of <Code(table)>.`` with kind `primary key`, `foreign key`, `unique constraint`, `check constraint`, or `constraint` when unknown | ``Drops foreign key `fk_p` of `t`.``, ``Drops unique constraint `uk_code` of `t`.`` |
| | no name | ``Drops the primary key of `t`.``, ``Drops a unique constraint of `t`.``, ``Drops the unique constraint on (`a`, `b`) of `t`.`` |
| DISALLOW_RENAME | ``Renames <kind> <Code(old)> to <Code(new)>.`` / ``Renames column <Code(old)> of <Code(table)> to <Code(new)>.`` | ``Renames table `logs` to `audit_logs`.``, ``Renames column `status` of `orders` to `state`.``, ``Swaps the names of tables `t` and `p`.`` |
| REQUIRE_PRIMARY_KEY | ``Creates table <Code(t)> without a primary key.`` / ``Removes the primary key of <Code(t)>, and the change doesn't add one back.`` | |

## Appendix C: finding codes

The table in `review/codes.go` is the source of truth; this copy says what each code covers.

| Code | Name | Title | Covers | Handling |
| --- | --- | --- | --- | --- |
| SY101 | `syntax-error` | Statement doesn't parse | Any statement the engine's parser refuses. `SY2xx` is reserved for syntax hints (§10). | Fix the syntax |
| WT101 | `object-not-found` | Object doesn't exist | A table, view, sequence, synonym, column, index, constraint, primary key, schema, database, routine, trigger, type, or event the statement needs doesn't exist, in the synced schema or after an earlier statement dropped or renamed it. | Correct the name, or create the object earlier in the change |
| WT201 | `name-taken` | Name is already taken | The name the statement creates or renames to is taken in its namespace, including a second primary key and a row type's name. | Choose another name, drop the old object first, or use `IF NOT EXISTS` |
| WT301 | `wrong-object-kind` | Object is of another kind | The name is a view where a table is needed, a table where an index is needed, a routine or type of another kind, and so on. | Use the statement for the object's kind |
| WT401 | `duplicate-name` | Statement uses one name twice | One statement adds, declares, or drops the same name twice. | Remove the repeat |
| WT402 | `no-current-schema` | No database or schema selected | An unqualified name has no current database or schema to resolve in. | Qualify the name, or select a database |
| WT403 | `name-too-long` | Name is too long | An identifier is longer than the engine allows. | Shorten it |
| WT404 | `reserved-name` | Name is reserved | An index named `PRIMARY`, a schema named `pg_…`, a column named like a system column. | Choose another name |
| WT405 | `schema-not-empty` | Schema still holds objects | `DROP SCHEMA` without `CASCADE` of a schema that holds objects. | Drop its objects first, or add `CASCADE` |
| WT406 | `system-column` | System column can't be changed | The statement drops, alters, or renames a system column (PostgreSQL). | Leave the column alone |
| WT407 | `cross-database-reference` | Name refers to another database | A name qualified with another database, which PostgreSQL refuses. | Name an object of this database |
| WT408 | `unsupported-form` | Engine rejects this form of statement | Clauses or objects the engine refuses together, such as `DROP INDEX CONCURRENTLY` with `CASCADE`. | Split or rewrite the statement |
| WT409 | `not-first-in-batch` | Statement must start its batch | SQL Server `CREATE VIEW`, `PROCEDURE`, `FUNCTION`, `TRIGGER`, or `SCHEMA` not first in its batch. | Put `GO` before it |
| WT410 | `trigger-schema-mismatch` | Trigger and table are in different schemas | A SQL Server DML trigger in another schema than its table. | Create it in the table's schema |
| WT411 | `index-backs-constraint` | Index belongs to a constraint | SQL Server `DROP INDEX` of an index that enforces a primary key or unique constraint. | Drop the constraint instead |
| OM101 | `online-migration-unsupported` | Engine has no online migration | Online migration asked for on PostgreSQL, TiDB, SQL Server, Oracle, or Snowflake. | Turn off **Online migration** |
| OM102 | `directive-invalid` | gh-ost directive is invalid | The directive is not a JSON object of strings, sets a flag Bytebase doesn't pass, a value of the wrong type or range, or flags that conflict. | Fix the directive |
| OM103 | `unsupported-change` | Change doesn't fit online migration | Several statements, a statement other than `ALTER TABLE`, an `ALTER` that renames the table, names another database's table, or targets a view. | Make it one `ALTER TABLE` of a table in this database |
| OM104 | `ghost-column-rename` | gh-ost can't rename the column | `CHANGE` that renames, a column name with whitespace, `RENAME COLUMN` (which loses the column's data). | Rename the column without online migration |
| OM105 | `archive-database-missing` | Backup database doesn't exist | The instance has no `bbdataarchive` database for gh-ost's tables. | Create it, or turn off **Online migration** |
| OM106 | `table-not-eligible` | Table doesn't meet gh-ost's requirements | No usable primary or unique key, a foreign key, a referencing foreign key, or a trigger. | Fix the table, or migrate without online migration |
| PB101 | `prior-backup-unsupported` | Engine has no prior backup | Prior backup enabled on Snowflake or OceanBase. | Turn off **Prior backup** |
| PB102 | `change-too-large` | Change exceeds the backup size limit | The change's SQL is larger than prior backup handles. | Split the change, or turn off **Prior backup** |
| PB103 | `backup-location-missing` | Backup database or schema doesn't exist | No `bbdataarchive` database, or schema in PostgreSQL. | Create it, or turn off **Prior backup** |
| PB104 | `nothing-to-back-up` | Change has nothing to back up | No top-level `UPDATE` or `DELETE` of a table (Oracle fails the backup). | Turn off **Prior backup** |
| PB105 | `statements-share-backup` | Statements can't share one backup | One table both updated and deleted, statements needing different `WITH` clauses or aliases, or more statements than the engine's backup takes together. | Split the change |
| PB106 | `backup-reads-wrong-table` | Backup may read another table | The target isn't found in its `FROM` list, a `WITH` query shadows it, or the backup's search path finds another table. | Qualify the name, or match the alias |
| PB107 | `table-not-backupable` | Table can't be backed up | A table in another database or behind a database link, a view, a table variable or `WITH` query target, a name with a quote, a `LONG` column. | Turn off **Prior backup** for this change |
| PB108 | `statement-not-backupable` | Statement form can't be backed up | A data-modifying `WITH`, `WHERE CURRENT OF`, `WITH XMLNAMESPACES`, `TOP`, `RETURNING`, TiDB `BATCH`, `REPLACE` or upsert, `EXPLAIN ANALYZE`. | Rewrite the statement, or turn off **Prior backup** |
| PB109 | `backup-reads-new-object` | Backup reads what the change creates | The backup runs before the change and reads a relation, column, function, or role the change creates. | Create the object in an earlier change |
| NU101 | `null-comparison` | Comparison with NULL is never true | `= NULL`, `<> NULL`, `!=` and the like, and an Oracle comparison with `''`, which Oracle reads as NULL. | Use `IS NULL` or `IS NOT NULL` |
| WH101 | `missing-where` | UPDATE or DELETE has no WHERE clause | An `UPDATE` or `DELETE` of any target form without `WHERE`. | Add a `WHERE`, or accept a whole-table change |
| DO101 | `drops-data` | Drops data | Drops a table, column, partition, sequence, schema, or database, MariaDB system versioning, or PostgreSQL `DROP OWNED`. | Accept the data loss, or keep the object |
| DO102 | `drops-definition` | Drops a definition | Drops a view, index, routine, trigger, type, an external or dynamic table whose rows live outside it, or another object the change could create again. | Check what depends on it |
| DO103 | `drops-access` | Drops a user, role, or security object | Drops a user, role, login, policy, network rule, key, certificate, credential, audit, user mapping, signature, share, managed account, or security integration. | Check who loses access |
| TR101 | `truncates-data` | Truncates data | `TRUNCATE` of tables, partitions, a cluster, a materialized view, or an error table. | Accept the data loss |
| DC101 | `drops-constraint` | Drops a constraint | Drops a primary key, foreign key, unique, check, or other constraint. | Check the data stays valid without it |
| RN101 | `renames-object` | Renames a table, view, or column | Renames a table, view, or column, or swaps two tables' names. | Update what reads the old name |
| PK101 | `table-without-primary-key` | Creates a table without a primary key | A `CREATE TABLE` whose table ends the change without a primary key, reported on the `CREATE`. | Declare a primary key |
| PK102 | `primary-key-removed` | Removes a table's primary key | A statement drops a table's primary key, by name, keyword, or column, and the change doesn't add one back; reported on the drop. | Keep the key, or add one back |
