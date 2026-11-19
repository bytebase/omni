# SQL Server parser test data

## `compat.tsv`: SQL Server's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/TSql160.g4`, then its near-misses (a token dropped, the end cut
off, a stray word inserted, a keyword swapped), then the token forms of
`suite:lexer`, with omni's verdict and the pinned SQL Server's
(`../../testdata/pins.tsv`). `TestCompatibility` (`../compat_test.go`) reads it
through `enginetest/compat` and starts no server. Recorded 2026-11-14 in under
a minute; a second recording on a fresh server wrote the same file.

| Rows | Total | Agree | Only SQL Server accepts | Only omni accepts |
| --- | --- | --- | --- | --- |
| One per alternative | 5,510 | 5,405 | 0 | 105 |
| Near-misses | 11,960 | 11,938 | 20 | 4 |
| `suite:lexer` | 67 | 66 | 0 | 2 |

SQL Server accepts 3,034 of the 5,510 alternatives. Five more alternatives get
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
  are refused: most of the 2,476 alternatives SQL Server refuses are such
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

- **`SET PARSEONLY ON`** (`compat.ParseOnlyVerdict`). The judge sends each row
  as a batch of its own on one session under `SET PARSEONLY ON`, which parses
  and neither compiles nor runs it, and puts the session back after each:
  `PARSEONLY`, `QUOTED_IDENTIFIER` and `FIPS_FLAGGER` take effect while a batch
  is parsed, and so may `USE`.
- **Every error is the parser's, and a rejection.** SQL Server's parser checks
  more than a grammar does: a `BREAK` outside a `WHILE` (135), a column with no
  type (173), an option it does not know (155, 195), a statement that must
  start its batch (111). The first error decides; the parser stops reading a
  statement there.
- **A lookup of what the server holds is an acceptance:** a database (911), an
  object (208, 15151), a collation (448), a language (7678), full-text search,
  which the image does not install (7609). The judge's session is in a
  database `a` holding a table `a (a INT)`, so the rows' names mostly resolve
  and the parser reads on.
- **A variable is the batch's to declare.** A row is one statement, so an
  undeclared variable (137, 1087, 16948 for a cursor) is declared before it
  and the row asked again. A statement that must start its batch, such as
  `CREATE PROCEDURE`, keeps its rejection, since nothing outside it could
  declare the variable.
- **A batch ends in a newline.** go-mssqldb sends a text that looks like a
  procedure name (`GO`, `a`, `BREAK`) as a call of that procedure, not a
  batch.
- **`GO` is the client's.** sqlcmd splits a script at a `GO` line; the server
  never sees one. A row is one line, so `GO` reaches the server as a name.

## The work list

17,409 rows agree: 12,212 when the file was recorded, the rest closed family
by family (#1034 to #1070, and then the batch's call by name, `ALTER
AUTHORIZATION` and the refusals around it, `SELECT`'s table sources, ODBC
escapes, graph `MATCH`, the security policies, sensitivity classifications and
federation shapes, and then `ALTER TABLE`'s remaining clauses, DML targets,
`WITH CHANGE_TRACKING_CONTEXT`, the `CREATE`/`DROP`/`ALTER` statement shapes,
query-expression parentheses, `GROUP BY`'s `DISTRIBUTED_AGG`, `TSEQUAL`, the
search-condition-only boolean operators, the shapes a query expression's
parentheses, a name's dot and a label's colon carry, and then the CTAS
alternatives and a line comment's opening).

128 rows disagree: 108 omni accepts and the pinned server refuses, and 20 the
other way. Each of the 128 was measured on the pinned image and on the same
release with the features below installed, and each was put to ScriptDOM — the
pin's own reader of the grammar the rows come from — which tells a check of a
parsed statement from a refusal of its text. Three are the grammar's and are
closed; the other 125 are a feature the image lacks, a statement of another
product, or a check made after the statement is read, and stay rows. The two
sections below hold the measurement so that no one measures them again.

A byte in them is the offset of the token the server's message names: SQL
Server reports a line, not a byte, and each row is one line. A row that says
"no byte" has a message that names what it refuses and reports no position.

## Rows only omni accepts, measured (108)

`TSql160.g4` is the grammar omni reads, so a statement the rule under the row
derives stays accepted where this server refuses it: the row is recorded here
instead of closed. Three rows are the grammar's, and are closed; the rest are
the image's features, another product's statements, or checks of the parsed
statement.

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

### Measured (105)

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
| Federations (Azure SQL Database's) | 12 | Msg 343 "Unknown object type 'FEDERATION'" at the object type word (7 for `CREATE`, 6 for `ALTER`, 5 for `DROP`), the answer `CREATE FOO a(a INT)` gets too: the word is looked up in the set of object types this product has, not read by a grammar |
| `GOTO a` | 3 | Msg 133, no byte: the label is looked up as the statement is read, and a batch of one statement declares none. `GOTO a` after a line `a:` parses |
| A CTE the statement does not use | 4 | Msg 422, no byte: `WITH a AS(SELECT *) SELECT * FROM a` parses |
| The type `PARSE` and `TRY_PARSE` convert to | 10 | Msg 10761 naming the type, no byte: the two functions take the datetime and numeric types, so `PARSE(1 AS int)` and `PARSE('1' AS datetime)` parse while `PARSE(1 AS nvarchar(10))` is refused as `PARSE(1 AS a)` is |
| `ALTER EXTERNAL LANGUAGE`'s environment variables | 1 | Msg 39133 naming the string, no byte: `ALTER EXTERNAL LANGUAGE a SET(CONTENT = 0x00)` and `... SET(CONTENT = 0x00, FILE_NAME = 'a')` parse |
| A name of 129 characters | 1 | Msg 103 at 7, the name; 128 characters parse, and ScriptDOM's lexer refuses the 129 as well. omni reads a name of any length and `mssql/review` reports it as WT403 (`NameTooLong`): the limit is the review layer's, and moving it into the lexer would report the text as a syntax error instead, so it stays its own change |

## Rows only SQL Server accepts, measured (20)

The server accepts each of these by looking a name up while it parses, which the
judge reads as an acceptance: in all but the last row the read stops there, and
the words after the name are never read. omni refuses those words, and so does
the grammar, so the rows stay and the server's acceptance is the row's
`expected` cell rather than omni's answer.

| Cause | Rows | The server's answer, and the control |
| --- | --- | --- |
| A full-text language term is looked up while it is parsed | 10 | Msg 7678 for the six rows whose term is a string or a name this server has no alias for, and Msg 7609 for the four whose term is an integer or a binary, since full-text search is not installed: both stop the read at the term (41 in `ALTER FULLTEXT INDEX`, 43 or 44 in `ALTER FULLTEXT STOPLIST`), so the `zzz` after it is never read. With `mssql-server-fts` installed, `ALTER FULLTEXT INDEX ON a ADD(a LANGUAGE 'English' zzz)` is Msg 102 at `zzz` (51) and `... LANGUAGE 'English'` parses, and `ALTER FULLTEXT STOPLIST a ADD 'a' LANGUAGE 'English' zzz;` is Msg 102 at `zzz` (53) and parses without it: the trailing word is the grammar's to refuse, and omni refuses it at 45, 43, 43, 46, 47, 45, 45, 48, 48 and 48. ScriptDOM refuses all ten as well — the four `ALTER FULLTEXT INDEX` rows at `zzz`, and the six stoplist rows at the stoplist name (24), where its reader wants the statement to end (a stoplist statement's rule ends in `requiredSemicolon`, TSql160.g:1720, which these rows do not have either) |
| `USE` takes the next word as a database name and looks it up | 9 | Msg 911 naming `FEDERATION` (4) or `ROOT` (4), and nothing after it is read. The control is `USE a zzz`: with a database that exists the server reads on and answers Msg 102 at `zzz` (6). ScriptDOM refuses the nine at 14, 31, 25, 45, 24, 44, 41, 45 and 41 |
| A collation is looked up while it is parsed | 1 | `ALTER DATABASE a COLLATE a zzz` answers Msg 448 for the collation *and* Msg 102 at `zzz` (27) in the same batch, and the judge reads the first error, a lookup, as an acceptance although the text is refused. `ALTER DATABASE a COLLATE Latin1_General_100_CI_AS zzz` answers Msg 102 at `zzz` (50), and the statement without the word parses. Closing this row means the judge reading the last error of a batch instead of the first, which is a change of the judge's reading and a pin bump of its own |
