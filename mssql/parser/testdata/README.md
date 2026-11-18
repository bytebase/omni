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
| Near-misses | 11,960 | 11,936 | 20 | 4 |
| `suite:lexer` | 67 | 65 | 0 | 2 |

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

17,406 rows agree: 12,212 when the file was recorded, the rest closed family
by family (#1034 to #1070, and then the batch's call by name, `ALTER
AUTHORIZATION` and the refusals around it, `SELECT`'s table sources, ODBC
escapes, graph `MATCH`, the security policies, sensitivity classifications and
federation shapes, and then `ALTER TABLE`'s remaining clauses, DML targets,
`WITH CHANGE_TRACKING_CONTEXT`, the `CREATE`/`DROP`/`ALTER` statement shapes,
query-expression parentheses, `GROUP BY`'s `DISTRIBUTED_AGG`, `TSEQUAL`, the
search-condition-only boolean operators, and the shapes a query expression's
parentheses, a name's dot and a label's colon carry). 131 rows disagree. What
is left, and why:

- **Only omni accepts (111):**
  - The judge has no PolyBase (`IsPolyBaseInstalled` is 0), so it refuses
    every `CREATE EXTERNAL TABLE`, `FILE FORMAT` and `STREAM`, `COPY INTO`
    and Synapse's `CREATE TABLE ... AS SELECT` (80). omni reads them by the
    grammar and the docs rather than refuse a kind of statement outright.
  - Federations (12, Azure SQL Database's): every `CREATE`, `ALTER` and `DROP
    FEDERATION` met Msg 343 "Unknown object type 'FEDERATION'" on the pinned
    server, so it has no answer for the statements themselves and only the
    near-misses above are closed. omni reads them by `TSql160.g` for the
    service that has them.
  - Full-text `LANGUAGE x zzz` (10): a full-text statement's language term is
    looked up while it is parsed (Msg 7678), so the server stops there and
    never reads the words after it — `LANGUAGE 'a' 1 2 3` is accepted too.
    omni reads those words by the grammar and refuses them. The same artifact
    keeps `ALTER DATABASE a COLLATE a zzz` open (Msg 448, then 102).
  - `SELECT PARSE(1 AS a)` and `TRY_PARSE` (6): these resolve their `AS` type
    while parsing and refuse one outside the set they convert, so
    `PARSE(1 AS nvarchar(10))` is Msg 10761 as `PARSE(1 AS a)` is. omni has no
    table of the types they take.
  - A CTE the statement does not use (4): `WITH a AS (SELECT 1) SELECT 1` is
    Msg 422, "Common table expression defined but not used", which needs the
    query's name resolution rather than its text.
  - `GOTO a` (3): the label is looked up as the statement is read, and a batch
    of one statement cannot declare the label it jumps to, so the judge reads
    Msg 133 as a rejection. `BREAK` outside a `WHILE` (135) is the same kind of
    check and is already recorded below.
  - 2 single rows: a name of 129 characters, which SQL Server refuses as Msg
    103 while `mssql/review`'s walk reports it as WT403 (the name-length check
    is the review layer's, and moving it into the lexer is its own change), and
    `SELECT 1 --(* a *)`, where the server's lexer refuses the four characters
    `--(*` inside a line comment (Msg 102, near `--(*`) although `--( *`,
    `-- (*` and `--a(*` are comments.
- **Only SQL Server accepts (20):**
  - `USE FEDERATION` (9): `USE` takes the next word as a database name and
    looks it up while it parses, so `USE FEDERATION ROOT WITH zzz RESET` is
    Msg 911 for the database `FEDERATION` and nothing after it is read — an
    acceptance, like every other lookup. Closing these would make omni accept
    the garbage.
  - The pinned server's full-text statements answer Msg 7678 and 7609 for the
    ten `LANGUAGE x zzz` rows above, whose words are read after the lookup
    stops it: omni reads them by the grammar and refuses them.
  - 1 single row: `ALTER DATABASE a COLLATE a zzz`, whose collation lookup
    (Msg 448) stops the server before the words after it.
