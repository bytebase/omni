# SQL Server parser test data

## `compat.tsv`: SQL Server's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/TSql160.g4`, then its near-misses (a token dropped, the end cut
off, a stray word inserted, a keyword swapped), then the token forms of
`suite:lexer`, with omni's verdict and the pinned SQL Server's
(`../../testdata/pins.tsv`, role `parser`, judge rev 2). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no server.
Re-recorded 2026-11-22 under judge rev 2 in under a minute; `OMNI_RECORD=verify`
on a fresh server confirmed every row.

| Rows | Total | Agree | Only SQL Server accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 5,510 | 5,415 | 0 | 95 |
| Near-misses | 11,960 | 11,960 | 0 | 0 |
| `suite:lexer` | 67 | 66 | 0 | 1 |

SQL Server accepts 3,035 of the 5,510 alternatives. Five more alternatives get
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
  are refused: most of the 2,475 alternatives SQL Server refuses are such
  options, and Azure SQL and Synapse forms SQL Server 2022 does not take.
- **Token forms** are in `../../testdata/tokens.tsv`, cited to the lexer rules
  of `TSql160.g`. One row picks a rule's text: a `statement` inside another is
  `SELECT 1`, since the grammar's first, `BREAK`, is refused outside a `WHILE`.
  `suite:lexer` holds the forms one sample cannot: `]]` and `""` in quoted
  names, `N'...'`, hex, money and number forms, nested comments, `@@` names,
  `$` pseudo-columns, ODBC escapes, a name of 129 characters.

After a parser fix, commit the rows that now agree:

```bash
go test ./mssql/parser/ -run '^TestCompatibility$' -update-compat
```

## What the judge needs

- **`SET PARSEONLY ON`** (`parseOnlyVerdictRev2` in `../compat_test.go`,
  extending `compat.ParseOnlyVerdict`). The judge sends each row as a batch of
  its own on one session under `SET PARSEONLY ON`, which parses and neither
  compiles nor runs it, and puts the session back after each: `PARSEONLY`,
  `QUOTED_IDENTIFIER` and `FIPS_FLAGGER` take effect while a batch is parsed,
  and so may `USE`.
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

17,441 rows agree (99.45%, up from 17,414 / 99.30% after judge rev 2): 12,212
when the file was first recorded, the rest closed family by family (#1034 to
#1070, and then the batch's call by name, `ALTER AUTHORIZATION` and the refusals
around it, `SELECT`'s table sources, ODBC escapes, graph `MATCH`, the security
policies, sensitivity classifications and federation shapes, and then
`ALTER TABLE`'s remaining clauses, DML targets, `WITH CHANGE_TRACKING_CONTEXT`,
the `CREATE`/`DROP`/`ALTER` statement shapes, query-expression parentheses,
`GROUP BY`'s `DISTRIBUTED_AGG`, `TSEQUAL`, the search-condition-only boolean
operators, the shapes a query expression's parentheses, a name's dot and a
label's colon carry, and then the CTAS alternatives and a line comment's
opening), followed by judge rev 2 closing all 20 `only_mssql` lookup-hides-syntax
rows and 18 context-only `only_omni` rows (`GOTO`, unused CTE, `PARSE`/`TRY_PARSE`,
and `ALTER EXTERNAL LANGUAGE`), and then closing the 27 parser bugs judge rev 2
exposed (18 loose `USE FEDERATION` near-misses in `utility.go` and 9
`ALTER FULLTEXT STOPLIST ... LANGUAGE ...` forms lacking `;` in `fulltext.go`).

96 rows disagree: all 96 are accepted only by omni, and 0 are accepted only by
SQL Server (`only_mssql` is 0, and all 11,960 near-misses agree).

A byte in them is the offset of the token the server's message names: SQL
Server reports a line, not a byte, and each row is one line. A row that says
"no byte" has a message that names what it refuses and reports no position.

## Rows only omni accepts, measured (96)

`TSql160.g4` is the grammar omni reads, so a statement the rule under the row
derives stays accepted where this server refuses it: the row is recorded here
instead of closed. Three rows of the 111 before #1468 were the grammar's, and
are closed; judge rev 2 closed 18 context-only rows (`GOTO`, unused CTE,
`PARSE`/`TRY_PARSE`, `ALTER EXTERNAL LANGUAGE`) and exposed 33 rows previously
hidden behind Msg 911 (`USE FEDERATION`, 24) and Msg 7609/7678
(`ALTER FULLTEXT STOPLIST` without `;`, 9), of which the 27 malformed forms
(18 loose `USE FEDERATION` near-misses and all 9 `ALTER FULLTEXT STOPLIST` rows
missing `;`) are now closed.

### Closed (3)

`TestOverAcceptsRefused` and `TestLineCommentOpening` hold each of these to the
byte measured, and hold the form beside it to the parser too: the CTAS forms are
the grammar's, which ScriptDOM parses, and `CREATE TABLE a(a INT)` is the plain
form SQL Server takes.

| Rows | Statement | SQL Server's answer | The rule it breaks |
| --- | --- | --- | --- |
| 1 | `CREATE TABLE a(HEAP) AS SELECT *` | Msg 173 at the column `HEAP` (15), which `(HEAP)` is read as | `ctasCreateTableStatement` is `( columnNameList )? withTableOptions As selectStatement` (TSql160.g:523-535) and neither part is optional, so a column list with no `WITH` is no CTAS, and no alternative beside it reaches `AS` with a column list. ScriptDOM refuses at `AS` (21), and omni now refuses at 21 |
| 1 | `CREATE TABLE a WITH() AS SELECT *` | Msg 102 near `)` (20) | `withTableOptions` is one option at least (TSql160.g:26710-26723): SQL Server, ScriptDOM and omni all refuse at 20 |
| 1 | `SELECT 1 --(* a *)` | Msg 102 near `--(*` (9) | a line comment opening with `(*` is no comment: ScriptDOM's lexer refuses the same four characters, and a comment that opens otherwise may hold them later (`SELECT 1 --x --(* a *)` parses). omni now refuses at 9 |

### Measured (96)

Each group's code was measured on both servers, and the control beside it — the
form the server takes, or the same text with the feature it lacks — is what
tells a check from the grammar.

| Cause | Rows | The server's answer, and the control |
| --- | --- | --- |
| The image has no PolyBase (`IsPolyBaseInstalled` 0) | 28 | Msg 102 at `EXTERNAL` (7) for `CREATE EXTERNAL TABLE` and `DROP EXTERNAL TABLE` (5), and at `FORMAT_TYPE` (35) for `CREATE EXTERNAL FILE FORMAT`, whose head and `WITH (` are read: the option word is not. With `mssql-server-polybase` 16.0.4305.5 installed and `polybase enabled` 1 on this release, `CREATE EXTERNAL TABLE a(a a) WITH(DATA_SOURCE = a)`, `DROP EXTERNAL TABLE a` and `CREATE EXTERNAL FILE FORMAT a WITH(FORMAT_TYPE = DELIMITEDTEXT)` all parse, and `CREATE EXTERNAL DATA SOURCE a WITH(LOCATION = 'a')` parses on the image as recorded |
| An external option's value | 11 | Msg 46508 "Incorrect syntax on external DDL option 'SERDE_METHOD'" or 'DATA_COMPRESSION', no byte, with PolyBase installed: the option is read and its value checked. `SERDE_METHOD = 'org.apache.hadoop.hive.serde2.columnar.ColumnarSerDe'` and `DATA_COMPRESSION = 'org.apache.hadoop.io.compress.DefaultCodec'` parse |
| The pin's file formats are five, not six | 1 | `FORMAT_TYPE = JSON` is Msg 102 at `JSON` (49) with PolyBase installed too; `DELIMITEDTEXT`, `RCFILE`, `ORC`, `PARQUET` and `DELTA` parse |
| Synapse's `COPY INTO` | 23 | Msg 156 near `INTO` (5) on both servers, and Msg 319 near `WITH` (21) for the row whose `WITH (IDENTITY_INSERT = a)` is read as a CTE clause: this product reads `COPY` as no statement at all. ScriptDOM parses all of them, and refuses the `WITH` row for its option's value |
| Synapse's `CREATE EXTERNAL STREAM` | 9 | Msg 102 at `EXTERNAL` (7), with PolyBase installed too: `EXTERNAL` is a word the shape takes only for the kinds this product has |
| Synapse's CTAS options | 5 | Msg 102 at `HEAP` (20) for `CREATE TABLE a WITH(HEAP) AS SELECT *`, and Msg 173 at the column `a` (15) then Msg 319 near `WITH` (18) for `CREATE TABLE a(a) WITH(HEAP) AS SELECT *`. `CREATE TABLE a(a INT) WITH(DATA_COMPRESSION = PAGE)` parses and `CREATE TABLE a WITH(DATA_COMPRESSION = PAGE) AS SELECT *` answers Msg 156 near `AS` (45): the option list is read and the `AS SELECT` is not. ScriptDOM takes both texts and refuses each for the distribution option a CTAS requires |
| Federations (`CREATE`/`ALTER`/`DROP`, Azure SQL Database's) | 12 | Msg 343 "Unknown object type 'FEDERATION'" at the object type word (7 for `CREATE`, 6 for `ALTER`, 5 for `DROP`), the answer `CREATE FOO a(a INT)` gets too: the word is looked up in the set of object types this product has, not read by a grammar |
| Federations (`USE FEDERATION`, Azure SQL Database's) | 6 | With scratch database `FEDERATION` present, SQL Server 2022 reads past the database name and refuses all 6 valid Azure SQL `USE FEDERATION` alternatives with Msg 102/319 (`parseUseFederationStmt` in `utility.go` requires every clause per `TSql160.g:13647`, closing all 18 loose near-misses) |
| A name of 129 characters | 1 | Msg 103 at 7, the name; 128 characters parse, and ScriptDOM's lexer refuses the 129 as well. omni reads a name of any length and `mssql/review` reports it as WT403 (`NameTooLong`): the limit is the review layer's, and moving it into the lexer would report the text as a syntax error instead, so it stays its own change |

## Rows only SQL Server accepts (0)

All 20 rows previously accepted only by SQL Server under judge rev 1 were broken
near-misses hidden behind early lookups (`ALTER DATABASE a COLLATE a zzz` via
Msg 448 + Msg 102, 9 `USE FEDERATION`/`USE ROOT` rows via Msg 911, and 10
`ALTER FULLTEXT INDEX`/`STOPLIST ... LANGUAGE ... zzz` rows via Msg 7609/7678).
Judge rev 2 resolves all three cases and SQL Server now rejects all 20.
