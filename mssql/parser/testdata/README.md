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
| One per alternative | 5,510 | 5,276 | 81 | 153 |
| Near-misses | 11,960 | 11,782 | 29 | 149 |
| `suite:lexer` | 67 | 62 | 2 | 3 |

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

17,120 rows agree: 12,212 when the file was recorded, the rest closed family
by family (#1034 to #1070, and then the batch's call by name, `ALTER
AUTHORIZATION` and the refusals around it), each reading its option lists and
clauses by the grammar and the server's checks. What is left, and why:

- **Only omni accepts (305):**
  - The judge has no PolyBase (`IsPolyBaseInstalled` is 0), so it refuses
    every `CREATE EXTERNAL TABLE`, `FILE FORMAT` and `STREAM`, `COPY INTO`
    and Synapse's `CREATE TABLE ... AS SELECT` (73). omni reads them by the
    grammar and the docs rather than refuse a kind of statement outright.
  - Statements omni still reads by their shape: federations (36, Azure SQL
    Database's), security policies (32), sensitivity classifications (29),
    `CREATE RULE` and `CREATE SYNONYM` (7).
  - Table sources and their arguments in `SELECT` (50): `OPENROWSET`,
    `OPENQUERY`, `CONTAINSTABLE` arguments, a trailing comma after a table.
  - Graph `MATCH` patterns (22): omni does not read the arrows.
  - 56 rows of one or two each: `DROP` of an object SQL Server 2022 does not
    know (`DROP AGGREGATE`, `DROP XML SCHEMA`, `DROP PARTITION`), a
    `CREATE TABLE ... AS SELECT`, `ENABLE TRIGGER` without a name, a name of
    two words where the grammar wants one.
- **Only SQL Server accepts (112):**
  - A procedure called by name as a batch's first statement, when another
    statement follows it (1: `zzz SELECT *`). omni reads the call only when
    the batch is the call alone (`../execute.go`, `parseImpliedExecStmt`):
    SQL Server reads a later statement too, but reports Msg 137 for a
    variable that statement uses and the batch does not declare (`my_proc SET
    @a = 1` is refused, `my_proc DECLARE @a INT` is not), and omni cannot
    tell a use from a declaration without resolving a batch's variables. The
    judge cannot declare the variable for the row either: a declaration
    before the call moves the call off the batch's first statement.
  - ODBC escapes (19: `{fn NOW()}`, `{d '2020-01-01'}`), split comparison
    operators (`! =`, `< >`, 7) and `IS DISTINCT FROM` (8).
  - Graph `MATCH` arrows (14), and `ALTER TABLE ... ADD CONNECTION` (8), a
    `COLUMN_SET`, a `CONTAINED` constraint and the rest of what an AST field
    would be needed for.
  - Full-text `LANGUAGE x zzz` (10): a full-text statement's language term is
    looked up while it is parsed (Msg 7678), so the server stops there and
    never reads the words after it — `LANGUAGE 'a' 1 2 3` is accepted too.
    omni reads those words by the grammar and refuses them. The same artifact
    keeps `ALTER DATABASE a COLLATE a zzz` open (Msg 448, then 102).
  - `USE FEDERATION` (9), `DELETE` from `OPENQUERY`, `OPENDATASOURCE` or a
    function (6), `WITH CHANGE_TRACKING_CONTEXT` (5), `FOR PATH` (4), a table
    hint holding a number (3).
