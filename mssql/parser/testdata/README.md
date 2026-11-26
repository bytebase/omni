# SQL Server parser test data

## `compat.tsv`: SQL Server's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/TSql160.g4`, then its near-misses (a token dropped, the end cut
off, a stray word inserted, a keyword swapped), then the token forms of
`suite:lexer`, with omni's verdict and the pinned SQL Server's
(`../../testdata/pins.tsv`, role `parser`, judge rev 3). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no server.
Re-recorded 2026-11-22 under judge rev 3 in under a minute; `OMNI_RECORD=verify`
on a fresh server confirmed every row.

| Rows | Total | Agree | Only SQL Server accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 6,044 | 6,044 | 0 | 0 |
| Near-misses | 12,268 | 12,268 | 0 | 0 |
| `suite:lexer` | 66 | 66 | 0 | 0 |
| `suite:optional-clauses` | 100 | 100 | 0 | 0 |
| `issue:BYT-10538` | 356 | 356 | 0 | 0 |
| `suite:search-condition` | 168 | 168 | 0 | 0 |

6,044 alternatives include the rows `tokens.tsv` adds with an optional clause
written out at each site a rule names it (`(every site)`), such as
`ALTER TABLE a DROP COLUMN IF EXISTS a`, and with each form of a keyword
choice at each site that names it (`(every alternative)`), such as
`WHERE a ! = 1` in an index filter. `suite:optional-clauses` holds
`IF EXISTS`, `IF NOT EXISTS`, `CREATE OR ALTER`, `CREATE OR REPLACE`,
`TEMPORARY` and `CASCADE` where T-SQL has none, which no grammar row reaches.
`suite:search-condition` holds the comparisons, `IS NULL`, `IN`, `LIKE` and
`BETWEEN`, which SQL Server reads only in a search condition: in a select
list, a SET, a value and an operand, where it refuses them, and in each place
that takes a search condition, where it accepts them.

The rev 3 judge accepts 3,208 of the 5,592 alternatives. Five more alternatives get
no row: they name a token neither keyword table gives a text
(`Credential`, `Connection`, `Object`, `Disk`, `Precision`), so ScriptDOM's
lexer never makes it.

## The grammar

- **ScriptDOM's, not SQL Server's.** SQL Server's grammar is not published.
  Microsoft's ScriptDOM (MIT) is its T-SQL parser for tools, and `TSql160.g` is
  the one for compatibility level 160, SQL Server 2022. `TSql160.g`,
  `TSqlTokenTypes.g` (the keywords every version shares) and
  `CodeGenerationSupporter.cs` (the words its actions match) are beside
  `TSql160.g4` byte for byte, at the commit `pins.tsv` names, with
  `ScriptDom.LICENSE`.
- **`TSql160.g4` is converted from them** by `../cmd/scriptdomg4`, since they
  are ANTLR 2 with C# actions. Each rule is under the line it came from, and
  what the conversion lost is noted beside it. `TestTSql160G4IsTheConversion`
  fails while `TSql160.g4` is anything else. After a pin bump:

  ```bash
  go run ./mssql/parser/cmd/scriptdomg4
  OMNI_RECORD=1 go test ./mssql/parser/ -run '^TestCompatibility$'
  ```
- **Most words are identifiers to ScriptDOM.** It has 182 keywords and
  reads every other word as an `Identifier` whose text an action or a
  predicate checks. The converter writes such an identifier as its words where
  the check names them: `Match(t, CodeGenerationSupporter.X)`, a chain of
  `TryMatch` and `Match`, and a predicate on the next word,
  `{NextTokenMatches(CodeGenerationSupporter.X)}? rule`, which becomes a copy
  of `rule` that starts with `X` (`rule_X`). An option a helper reads by its
  text (`ParseOption(tOption)`) stays an identifier, sampled `a`, so its rows
  are refused: most of the 2,384 alternatives the rev 3 judge refuses are such
  options, and Azure SQL forms outside the rev 3 fallback families.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the lexer rules
  of `TSql160.g`. One row picks a rule's text: a `statement` inside another is
  `SELECT 1`, since the grammar's first, `BREAK`, is refused outside a `WHILE`.
  `suite:lexer` holds the forms one sample cannot: `]]` and `""` in quoted
  names, `N'...'`, hex, money and number forms, nested comments, `@@` names,
  `$` pseudo-columns, and ODBC escapes. Identifiers longer than 128 characters
  are intentionally accepted by `mssql/parser` so `mssql/review` (`WT403` /
  `NameTooLong`) reports a named finding rather than a syntax error; that
  contract is pinned in `../lexer_test.go`
  (`TestLexerLongIdentifierAcceptedForReviewLayer`).

After a parser fix, commit the rows that now agree:

```bash
go test ./mssql/parser/ -run '^TestCompatibility$' -update-compat
```

## What the judge needs

- **`SET PARSEONLY ON` + scoped ScriptDOM fallback** (`parseOnlyVerdictRev3` in
  `../compat_test.go`, extending `parseOnlyVerdictRev2` and
  `compat.ParseOnlyVerdict`). The judge sends each row as a batch of its own on
  one session under `SET PARSEONLY ON`, which parses and neither compiles nor
  runs it, and puts the session back after each: `PARSEONLY`,
  `QUOTED_IDENTIFIER` and `FIPS_FLAGGER` take effect while a batch is parsed,
  and so may `USE`. When SQL Server 2022 rejects a statement belonging to the
  documented PolyBase, Synapse, or Azure SQL Federation families
  (`CREATE`/`DROP EXTERNAL TABLE`, `CREATE EXTERNAL FILE FORMAT`,
  `CREATE EXTERNAL STREAM`, `COPY INTO`, Synapse CTAS
  `CREATE TABLE ... WITH (...) AS SELECT`, and
  `CREATE`/`ALTER`/`DROP`/`USE FEDERATION`), judge rev 3 falls back to
  Microsoft's pinned `mssql/scriptdom` (`TSql160Parser`, role `scriptdom`). If
  ScriptDOM rejects a Synapse CTAS solely for omitting `DISTRIBUTION`, judge
  rev 3 retries with `, DISTRIBUTION = ROUND_ROBIN` appended to the non-empty
  `WITH(...)` list so ScriptDOM parses and rejects broken tails (such as
  `... AS SELECT * zzz`). SQL Server 2022's native `EXTERNAL` statements
  (`DATA SOURCE`, `LIBRARY`, `LANGUAGE`, `RESOURCE POOL`) and all statements
  SQL Server accepts never invoke ScriptDOM.
- **Every error in the batch is checked.** SQL Server's parser checks more than
  a grammar does: a `BREAK` outside a `WHILE` (135), a column with no type
  (173), an option it does not know (155, 195), a statement that must start its
  batch (111). Judge rev 2 inspects all errors returned in the batch rather than
  stopping at the first, so a lookup paired with a syntax error (such as
  `ALTER DATABASE a COLLATE a zzz`, Msg 448 + Msg 102) is rejected.
- **A lookup of what the server holds is an acceptance only when no syntax
  error accompanies it:** a database (911), an object (208, 15151), or a
  collation (448). The judge's session runs in database `a` holding table
  `a (a INT)`, and scratch databases `FEDERATION` and `ROOT` exist so
  `USE FEDERATION ...` and `USE ROOT ...` resolve past the database name
  instead of aborting on Msg 911 before reading the rest of the statement.
- **Full-text language lookups probe the tail (7609, 7678).** Without
  `mssql-server-fts` installed, SQL Server aborts inside `LANGUAGE <term>`
  before reading the rest of the statement. On Msg 7609 or 7678, judge rev 2
  re-asks with the validated `LANGUAGE <term>` clause removed (in
  `FULLTEXT INDEX`) or normalized to `DROP ALL` (in `ALTER FULLTEXT STOPLIST`),
  so SQL Server parses and checks the statement's tail and required semicolon.
- **Batch context and placeholder values are supplied on retry:**
  - An undeclared variable (137, 1087, 16948 for a cursor) is declared before
    the statement and asked again; statements that must start their batch keep
    their rejection (Msg 111).
  - An undeclared `GOTO` label (133) is declared after the statement (`<label>:`)
    and asked again.
  - An unused CTE (422) is referenced by appending `FROM <cte>` to a trailing
    table-less `SELECT *` and asked again, so valid CTEs accept while invalid
    combinations such as `WITH CHANGE_TRACKING_CONTEXT(0x00), a AS(SELECT *) SELECT *`
    still reject (Msg 22108).
  - A placeholder `PARSE` / `TRY_PARSE` target type (10761) is retried with
    `INT` so SQL Server parses the rest of the call and trailing tokens.
  - An `ALTER EXTERNAL LANGUAGE ... ENVIRONMENT_VARIABLES` placeholder string
    (39133) is retried with `'{"a":"b"}'` so SQL Server parses past the JSON
    check and rejects broken tails such as `... = 'a' zzz)`.
- **A batch ends in a newline.** go-mssqldb sends a text that looks like a
  procedure name (`GO`, `a`, `BREAK`) as a call of that procedure, not a
  batch.
- **`GO` is the client's.** sqlcmd splits a script at a `GO` line; the server
  never sees one. A row is one line, so `GO` reaches the server as a name.

## The work list

19,002 / 19,002 rows agree (**100.00%**). `issue:BYT-10538` holds the forms a
T-SQL codebase failed on and the forms beside them, every one closed.
`suite:search-condition` found 58 comparisons omni read outside a search
condition (`SELECT 1 = 1`, `UPDATE a SET a = 1 = 1`, `WHERE a = (1 = 1)`),
closed. The `(every alternative)` rows found 5
gaps, closed: a split comparison operator in an index filter, HOLDLOCK and no
condition in a table-valued function's list, a variable as OPENXML's column
pattern, and a pseudo-column as a text column's last part. The `(every site)` and
`suite:optional-clauses` rows found 21 gaps, closed: `CREATE OR ALTER` on
anything but a function, procedure, trigger or view (13), `IF EXISTS` on
`DROP PARTITION FUNCTION`/`SCHEME` and `DROP XML SCHEMA COLLECTION` (3), a DML
`FROM` with no table source (2), `NOT FOR` without `REPLICATION` (1), and a
constraint's `ON` with a string filegroup or none (2). Before them,
17,536 / 17,536 rows agreed (0 `only_mssql`, 0 `only_omni`, up from
17,506 / 17,536 = 99.83% after judge rev 3): 12,212 when the file was first
recorded, the rest closed family by family (#1034 to #1070, and then the
batch's call by name, `ALTER AUTHORIZATION` and the refusals around it,
`SELECT`'s table sources, ODBC escapes, graph `MATCH`, the security policies,
sensitivity classifications and federation shapes, and then `ALTER TABLE`'s
remaining clauses, DML targets, `WITH CHANGE_TRACKING_CONTEXT`, the
`CREATE`/`DROP`/`ALTER` statement shapes, query-expression parentheses,
`GROUP BY`'s `DISTRIBUTED_AGG`, `TSEQUAL`, the search-condition-only boolean
operators, the shapes a query expression's parentheses, a name's dot and a
label's colon carry, and then the CTAS alternatives and a line comment's
opening), followed by judge rev 2 closing all 20 `only_mssql`
lookup-hides-syntax rows and 18 context-only `only_omni` rows (`GOTO`, unused
CTE, `PARSE`/`TRY_PARSE`, and `ALTER EXTERNAL LANGUAGE`), closing the 27 parser
bugs judge rev 2 exposed (18 loose `USE FEDERATION` near-misses in `utility.go`
and 9 `ALTER FULLTEXT STOPLIST ... LANGUAGE ...` forms lacking `;` in
`fulltext.go`), judge rev 3 falling back to the pinned `mssql/scriptdom`
(`TSql160Parser`) for the PolyBase, Synapse, and Azure SQL Federation families
(closing 93 valid `only_omni` rows while moving the 1 intentionally
review-layer 129-char identifier test into `../lexer_test.go`), and finally
closing all 30 parser gaps judge rev 3 exposed in `external.go` (27 CETAS
`CREATE EXTERNAL TABLE ... WITH (...) AS SELECT` rows, 1 `CREATE EXTERNAL
STREAM` option-duplication row, 1 `COPY INTO ... WITH(IDENTITY_INSERT = a)`
row, and 1 duplicate `FORMAT_OPTIONS(USE_TYPE_DEFAULT = ...)` row).

0 rows disagree: all 6,044 alternatives, all 12,268 near-misses, all 66
`suite:lexer`, all 100 `suite:optional-clauses`, all 356 `issue:BYT-10538`
and all 168 `suite:search-condition` rows agree.

### Closed over-accepts (3)

`TestOverAcceptsRefused` and `TestLineCommentOpening` hold each of these to the
byte measured, and hold the form beside it to the parser too: the CTAS forms are
the grammar's, which ScriptDOM parses, and `CREATE TABLE a(a INT)` is the plain
form SQL Server takes.

| Rows | Statement | SQL Server's answer | The rule it breaks |
| --- | --- | --- | --- |
| 1 | `CREATE TABLE a(HEAP) AS SELECT *` | Msg 173 at the column `HEAP` (15), which `(HEAP)` is read as | `ctasCreateTableStatement` is `( columnNameList )? withTableOptions As selectStatement` (TSql160.g:523-535) and neither part is optional, so a column list with no `WITH` is no CTAS, and no alternative beside it reaches `AS` with a column list. ScriptDOM refuses at `AS` (21), and omni now refuses at 21 |
| 1 | `CREATE TABLE a WITH() AS SELECT *` | Msg 102 near `)` (20) | `withTableOptions` is one option at least (TSql160.g:26710-26723): SQL Server, ScriptDOM and omni all refuse at 20 |
| 1 | `SELECT 1 --(* a *)` | Msg 102 near `--(*` (9) | a line comment opening with `(*` is no comment: ScriptDOM's lexer refuses the same four characters, and a comment that opens otherwise may hold them later (`SELECT 1 --x --(* a *)` parses). omni now refuses at 9 |
