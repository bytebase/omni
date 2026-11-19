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
| `suite:manifest` | 612 | 588 | 2 | 2 | 20 |

The manifest's `compat_comment_*` rows hold the categories `COMMENT ON`
takes and the ones it refuses: `TABLE` (a view is commented through this
form), `MATERIALIZED VIEW`, `COLUMN`, `EDITION`, `INDEXTYPE`, `OPERATOR`,
`AUDIT POLICY`, `MINING MODEL` and `PROPERTY GRAPH` are accepted, and
`VIEW`, `INDEX`, `SEQUENCE`, `FUNCTION`, `PROCEDURE`, `PACKAGE`, `TRIGGER`,
`TYPE`, `SYNONYM`, `DATABASE LINK` and `DIMENSION` are `ORA-32594` at the
category word (byte 11), as `COMMENT.bnf` says.

## What the judge needs

- **`DBMS_SQL.PARSE`, as a user of its own.** The judge runs
  `compat.OracleParse` for each statement as `OMNI_JUDGE`, which may create
  tables, views, sequences, synonyms, materialized views, types, triggers and
  procedures, and nothing else: a statement on the server or on another user
  answers `ORA-01031` after the parse. `PARSE` runs DDL, so after a statement
  that ran the block drops everything the user owns and builds the fixture
  again: the table `t` a trigger's header stands on and the materialized view
  `v` on it, named apart from the `a` every other row writes. It enables table
  locks before it drops a table, since `ALTER TABLE ... DISABLE TABLE LOCK`
  leaves one that cannot be dropped.
- **Lookups before the parse are asked again.** Oracle looks up the object of
  `ALTER TABLE a` (ORA-00942) and `ALTER MATERIALIZED VIEW a` (ORA-12003)
  before it reads the rest, so the judge asks again with `t` or `v` in `a`'s
  place. On a real table `ORA-01735` at byte 0 is an option the table's kind
  does not take (`OVERFLOW` on a heap table): Oracle reads no further, so it is
  an acceptance, as a lookup is.
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

## Rows only Oracle accepts, measured

The 129 rows whose cells are `reject` where Oracle says `accept` were each
measured on the pinned server at judge revision 3: #1438, #1436 and #1445
measured the families they name, #1450 the rest of its own, and a follow-up
after #1450 measured what was left. None is a form omni can read: Oracle
answers a name it looks up, a kind of object that does not take the clause, or
a check of the parsed statement, and never reaches the tokens the row was
generated from. A byte is `DBMS_SQL.LAST_ERROR_POSITION`.

| Cause | Rows | Oracle's answer |
| --- | --- | --- |
| `suite:reference` | 16 | refused, or `INVALID` (`PLS-00589`); only `PLS-001xx` counts as syntax |
| `ALTER MATERIALIZED VIEW ... SUBPARTITIONS 1 ...` | 12 | `ORA-14156` at `SUBPARTITIONS` (52) |
| `GRANT`/`REVOKE` to and from program units | 32 | `ORA-28700`, `ORA-28701` at 0: no role list is read |
| `GRANT a TO a, a zzz`, `... BY a`, `... IDENTIFIED BY a,` | 4 | `ORA-01700` at 14, `ORA-03001` at 13, `ORA-00989` at 28 |
| `ALTER TABLE ... UPGRADE [NOT] INCLUDING DATA` | 5 | `ORA-01735` at 0 on `t`; no column property is read |
| `MODIFY [SUB]PARTITION a LOB(a(RETENTION)` | 3 | `ORA-03059` at `LOB` (25, 37, 40) |
| `ALTER TABLE a MODIFY (PARTITION PARTITION)` | 1 | `ORA-00904` at `PARTITION` (21) on `t` |
| `ALTER TABLE a ADD (a NUMBER, a` | 1 | `ORA-00957` at 28 (a column named twice) |
| `CREATE MATERIALIZED VIEW LOG ON a ...` | 4 | `ORA-00942` at 0 (3); `ORA-14017` at the end of the truncated bound (1) |
| `CREATE MATERIALIZED VIEW a ...` | 5 | `ORA-02538` for `TABLESPACE SET` (2); `ORA-14030` or `ORA-14017` (2); `ORA-00904`/`ORA-22912` at the nested column (1) |
| `CREATE TABLE a (PARTITION ...)` | 5 | `ORA-00902` at 25: `PARTITION` is a column, the next word its datatype |
| `CREATE TABLE a PARTITION BY RANGE (a) (PARTITION VALUES LESS THAN (1,` | 1 | `ORA-14017` at 66 |
| `CREATE TABLE ... LOB (a) STORE AS ...` on a hash partition | 2 | `ORA-22877` at 74, 76 |
| `CREATE TABLE a COLUMN a IS OF (a)` | 1 | `ORA-03001` at 32 |
| `CREATE TABLE a INMEMORY INMEMORY NO` | 1 | `ORA-64350` at 33 (one in-memory clause) |
| `CREATE TABLE a ORGANIZATION EXTERNAL (TYPE a` | 1 | `ORA-29829` at 13 (no such implementation type) |
| `CREATE TABLESPACE SET a USING TEMPLATE (` | 4 | `ORA-02199` at the end (3); `ORA-25147` at 75 (1) |
| `CREATE SCHEMA AUTHORIZATION a ...` | 3 | `ORA-02421` at 0 (no such user) |
| `CREATE TYPE BODY`/`FUNCTION`/`PACKAGE`/`PROCEDURE` `a.name` | 6 | `ORA-01435` at 0 (no such user) |
| `ADMINISTER KEY MANAGEMENT ...` | 5 | `ORA-46603`, `ORA-46601` at the command word (2); `ORA-01031` at 0 (3) |
| `ALTER DISKGROUP a REBALANCE ...` | 2 | `ORA-15260` at 0 (permission on the disk group) |
| `CREATE LIBRARY a IS 's' SHARING =` | 1 | `ORA-24344` at 0 (compiled with errors) |
| `CREATE AUDIT POLICY a ACTIONS COMPONENT =` | 1 | `ORA-46373` at 0 (no audit option) |
| `CREATE MLE MODULE a LANGUAGE JAVASCRIPT AS` | 1 | `ORA-04125` at 0 (empty module) |
| `DROP CONTEXT` | 1 | `ORA-01804` at 12, the end (time zone files) |
| `LOCK TABLE a PARTITION (1)` / `SUBPARTITION (1)` | 2 | `ORA-00942` at 11; on a partitioned table `ORA-14108`, `ORA-14173` at the number |
| `SELECT`/`DELETE` ... `EXTERNAL ((a NUMBER, a` | 2 | `ORA-00957` at 34 |
| `CREATE CLUSTER a (a NUMBER, a` | 1 | `ORA-00957` at 27 |
| `CREATE VECTOR INDEX a ON a (a) INCLUDE (a` | 1 | `ORA-00957` at 38 |
| `CREATE PROPERTY GRAPH a ... EDGE TABLES (a` | 1 | `ORA-42404` at 53 (a graph table named twice) |
| `MERGE INTO a` / `MERGE INTO USING s ...` | 2 | `ORA-38107` at 0 (no `USING`); `ORA-00942` at 11 |
| `SELECT * FOR UPDATE WAIT` | 1 | `ORA-30020` at 24, the end |
| `ASSOCIATE STATISTICS ... USING a` | 1 | `ORA-29820` at 0 (no statistics type) |

Two families are left as rows for another reason. `TABLESPACE SET` in a
materialized view or its log answers `ORA-02538` at the clause (byte 0 in the
log's), which is this server's setup: it has no sharding and cannot create a
tablespace set (`CREATE TABLESPACE SET ... ORA-02537`). The tablespace-set
template rows do not set the clause at all, so their `ORA-02199` is the
missing `DATAFILE` the truncated text leaves.


## Tapes

Every other test that asked Oracle replays a tape through
`oracle/internal/testcontainer.Open(t, user)`, here (`tapes/`, 4 tapes) and in
`oracle/review` (6 tapes, the walk-through oracle and its sweep among them).
A recording starts each top-level test on a server without the users,
public synonyms and SYSTEM objects another test made. Oracle's generated
names (`SYS_C…`, `SYS_C_SNAP$_…`, LOB segments, identity sequences) get
stand-ins. Two recordings on fresh servers wrote the same tapes.

`TestOraclePLSQLDictionary` reads the PL/SQL Oracle ships, which a tape may
not hold, so it runs only in a recording (`tape.Live`).
