# Oracle parser test data

## `compat.tsv`: Oracle's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/OracleDocs.g4`, then its near-misses, with omni's verdict and
the pinned Oracle's (`../../testdata/pins.tsv`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no
server. Recorded 2026-11-14 in 35 minutes, most of them on the statements
Oracle never answers (below); a second recording on a fresh server wrote the
same file. A rejection carries the byte Oracle reports it at.

| Rows | Total | Same cell | Both reject, another byte | Only Oracle accepts | Only omni accepts |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 14,809 | 7,549 | 1,723 | 2,488 | 3,049 |
| Near-misses | 31,920 | 9,725 | 10,334 | 2,136 | 9,725 |

Most of what only Oracle accepts is materialized views (2,018 rows of
`CREATE` and `ALTER MATERIALIZED VIEW`), `ALTER TABLE` and `CREATE TABLE`
(913), `GRANT` and `REVOKE` (382) and `SELECT` (189). Most of what only omni
accepts is `ALTER TABLE` (1,881), materialized view logs (1,072), databases,
pluggable databases and tablespaces (3,566) and `CREATE TABLE` (648): forms
omni reads loosely. No parser gap is fixed here.

## The grammar

`OracleDocs.g4` is the syntax of the 171 statement pages of the SQL Language
Reference that `../bnf/` quotes, converted by `../cmd/docsgrammar` with
`enginetest/grammar/docs`, each rule under its page's URL. The `docs` row of
`pins.tsv` is the SHA-256 of what the converter read and the date the pages
were quoted. Its readings of the quoted BNF:

- **A name in lower case is a fragment** where a page defines it (`name ::=`
  or `name:`), its own page's definition first, then the first page in file
  order. A fragment another page defines is copied under the page that uses
  it, except SELECT's `subquery` and `query_block`, which are `QUERY`.
- **A name no page defines is a value,** sampled by class in
  `../../testdata/tokens.tsv`: `ID`, `NUMBER`, `STRING`, `EXPR`, `CONDITION`,
  `DATA_TYPE`, `QUERY`, `HINT`. The PL/SQL a unit's header names
  (`plsql_function_source`, ...), CREATE SCHEMA's statements and a type's
  method specifications get a sample each.
- **A clause no quoted page defines** (`storage_clause`, `LOB_parameters`,
  ...) is left out where it is optional: its alternative, or its `[ ]`. A
  definition that needs one is left out the same way. Each page lists what
  it left out above its rules.
- `bnf-fixes.tsv` corrects five places where the quote does not read as the
  notation; `docs-fixes.tsv` is the converter's own and is empty.

## Suites

The corpora that held `Parse` to a live Oracle are `suite:` rows, judged by
the same recording. `TestCompatibilitySuitesAreTheCorpora` keeps them equal
to `lexerCorpus` and to the statements of `coverage/reference_oracle.tsv` and
`coverage/compat_oracle.tsv`, which other tests still read.

| Suite | Rows | Same cell | Both reject, another byte | Only Oracle accepts | Only omni accepts |
| --- | --- | --- | --- | --- | --- |
| `suite:lexer` | 77 | 57 | 7 | 9 | 4 |
| `suite:reference` | 321 | 223 | 77 | 21 | 0 |
| `suite:manifest` | 592 | 553 | 21 | 16 | 2 |

## What the judge needs

- **`DBMS_SQL.PARSE`, as a user of its own.** The judge runs
  `compat.OracleParse` for each statement as `OMNI_JUDGE`, which may create
  tables, views, sequences, synonyms, materialized views, types, triggers and
  procedures, and nothing else: a statement on the server or on another user
  answers `ORA-01031` after the parse. `PARSE` runs DDL, so after a statement
  that ran the block drops everything the user owns and builds the fixture
  again: the table `t` a trigger's header stands on, named apart from the `a`
  every other row writes.
- **A probe tells syntax from the rest** (`compat.OracleVerdict`): the
  statement again with `\n ~` after it. Oracle reads all of a statement before
  it resolves a name or checks a privilege, so a parser error comes back at the
  same byte and anything later is answered at the tail. A rejection carries
  `DBMS_SQL.LAST_ERROR_POSITION`, a byte offset.
- **PL/SQL is the compiler's.** A unit created with compilation errors
  (`ORA-24344`) or a block that did not compile (`ORA-06550`) is a rejection
  when an error is a syntax error: `PLS-001xx`, or the SQL parser's on a
  statement inside. A stored unit's rejection carries no byte.
- **Oracle fails on some statements in its own code:** it raises `ORA-00600`,
  leaves the session failing with `ORA-00000`, kills its process (`ORA-07445`),
  or never answers. Most of the forms of `ALTER MATERIALIZED VIEW LOG ... MODIFY
  PARTITION ... ADD SUBPARTITION`, broken ones too, parse on and on, and the
  session's process keeps its locks after the session is killed. The judge
  gives a text 15 seconds, then kills its server process (`docker exec`) and
  asks again on a new session. Failing twice is a rejection with no byte.
- **`ORACLE_PARSER_REF_DSN`** judges against another server in a recording,
  instead of starting the pinned image. It must run the pinned release, and
  its user must be a scratch account: the block drops everything that user
  owns after each statement that ran. The judge cannot kill a process there,
  so a statement Oracle never answers fails the recording.
