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


## Rows only omni accepts, measured

The other direction: the 750 rows whose cells are `accept` where the pinned
Oracle says `reject`. Each was measured the same way, at judge revision 3,
with the code `DBMS_SQL.PARSE` answers and `LAST_ERROR_POSITION`'s byte. The
question a row asks is whether Oracle's refusal is the grammar's: the page
`bnf/` quotes under it is the grammar omni must read, so a statement the page
derives stays accepted even where this server refuses it, and is recorded
here instead of closed. Fifteen rows are text the page refuses too, and are
closed, and so are the thirty-two `PARTITIONSET` rows: the page's clause is
only reachable in a composite-sharded database, which `OracleDocs.g4` records
where the converter could drop the rule (its `fix:` notes on
`alter_table_partitioning_clauses`' neighbours) and this server is not. The
other 703 are a name Oracle looks up, a check of the parsed statement, or the
server's own setup, and stay rows.

### Closed

`TestOverAcceptsRefused` holds each of these to the byte Oracle reports, and
holds the form beside it, which Oracle takes, to the parser.

| Family | Rows | Oracle's answer | The rule it breaks |
| --- | --- | --- | --- |
| `ANALYZE a VALIDATE STRUCTURE` | 1 | `ORA-01490` at `VALIDATE` (10) | ANALYZE names `TABLE`, `INDEX` or `CLUSTER` first; the three near-misses (`TABLES`, `INDEXES`, `CLUSTERING` for it) now answer at the word they leave over, 15, 16 and 19, as Oracle does |
| `CREATE SCHEMA a ...` | 1 | `ORA-02420` at 14 | `CREATE SCHEMA AUTHORIZATION schema`, the word required |
| `DROP USER` | 1 | `ORA-01935` at 9 | a user is required |
| `DROP AUDIT POLICY` | 1 | `ORA-46353` at 17 | a policy is required |
| `ALTER TRIGGER a RENAME TO` | 1 | `ORA-04070` at 25 | `RENAME TO new_name` |
| `PURGE a`, `PURGE REDO`, `PURGE DBTIMEZONE` | 3 | `ORA-38302` at 6 | `PURGE { TABLE \| INDEX \| TABLESPACE \| RECYCLEBIN \| DBA_RECYCLEBIN }`; five more rows (`PURGE TABLES a` among them) moved to Oracle's byte without changing cell |
| `ALTER TABLE a SET PARTITIONING`, `ALTER MATERIALIZED VIEW a SET PARTITIONING` | 2 | `ORA-14853` at 0 | `SET PARTITIONING { AUTOMATIC \| MANUAL }` |
| `CREATE FUNCTION IF NOT a ...` | 1 | `ORA-11543` at 16 | `[ IF NOT EXISTS ]`, all three words or none |
| `CREATE MATERIALIZED VIEW a INMEMORY(a) ...` | 2 | `ORA-00922` at 56, 77 | a materialized view's `inmemory_column_clause` is `( column [, ...] ) { INMEMORY \| NO INMEMORY }`, the keyword after the columns; a table's is the other way round and stays |
| `SELECT 1 "a""b" FROM dual` | 1 | `ORA-25716` at 9 | a quoted identifier holds no `"`, doubled or not (`"a b"` still parses) |
| `SELECT 1 FROM dual WHERE 1 ~= 2` | 1 | `ORA-00911` at 27 | SQL has no `~=`; PL/SQL does, and `BEGIN IF 1 ~= 2 THEN NULL; END IF; END;` still parses |
| `ALTER TABLE a MODIFY PARTITIONSET ...` | 32 | `ORA-02577` at `PARTITIONSET` (21) | a partitionset belongs to a composite-sharded table, so the word is refused before the clause is read; no form of it parses on this server, so these have no control beside them |

### Measured

A byte is `DBMS_SQL.LAST_ERROR_POSITION` unless the row says otherwise. Each
group's code was measured on the pinned server, and the control beside it,
which is the same text with the missing feature added, is what tells a check
from the grammar; the controls are named where a reader would ask.

| Cause | Rows | Oracle's answer |
| --- | --- | --- |
| `CREATE CLUSTER ... PARTITION BY RANGE(...)`, and every clause inside its partitions | 146 | `ORA-14179` at the `PARTITION` clause (45): a cluster is partitioned only when it is a hash cluster, so `CREATE CLUSTER ... HASHKEYS 10 PARTITION BY RANGE(...)` parses and answers `ORA-01031` when it runs (the clause is the page's, the check is not) |
| `CREATE PLUGGABLE DATABASE a USING 's' ...` | 82 | `ORA-65125` at 34, the file name |
| `SELECT`/`DELETE ... FROM EXTERNAL(... TYPE a ...)` (42), a materialized view's `INMEMORY (columns) INMEMORY` or `... NO INMEMORY` (8) | 50 | `ORA-00922` at 39 and 48, and at the query (65, 68, 80): `TYPE ORACLE_LOADER` parses, and the column clause's documented form is refused here: this server takes no inmemory column list for a materialized view, the page takes one |
| `CREATE TABLE a ENABLE ... USING INDEX ...` | 47 | `ORA-02432` at 34, `USING`, and the message names the check: "cannot enable primary key - primary key not defined for table". The byte is the token after the `enable_disable_clause`, not `ENABLE`; the control `CREATE TABLE a (b NUMBER, PRIMARY KEY(b)) ENABLE PRIMARY KEY USING INDEX` parses, and `... USING INDEX bb` answers `ORA-03049` at the name, so the clause is read and the table's own element list fails the check |
| `ALTER TYPE a ...` | 48 | `ORA-04043` at 0: the type `a` is not there. With a type, `ADD ATTRIBUTE`, `ADD ATTRIBUTE (...)`, `ADD MAP MEMBER FUNCTION ... RETURN NUMBER`, `ADD MEMBER PROCEDURE` and `INSTANTIABLE` parse |
| a column, property, attribute or key named twice | 41 | `ORA-00957` at the second name |
| `SHARING = { METADATA \| DATA \| NONE }` | 21 | `ORA-65021` at the clause: not an application root |
| property graph: an element table, property or key named twice | 21 | `ORA-42404` at 41-58 |
| `ALTER INMEMORY JOIN GROUP a(...)`, `CREATE INMEMORY JOIN GROUP a(...)` | 14 | `ORA-00903` at the table name (29-43) |
| `CREATE PFILE ...`, `CREATE SPFILE ...` | 14 | `ORA-65168` at 13, 24, 25: a name is required in this container, and `CREATE PFILE = 'x.ora' FROM SPFILE` parses (`ORA-01031`) |
| `COLLATE a`, `DEFAULT COLLATION a` | 12 | `ORA-43929` at 0: `MAX_STRING_SIZE=STANDARD` |
| `ALTER DATABASE RECOVER MANAGED STANDBY DATABASE { FINISH \| TO LOGICAL STANDBY \| USING INSTANCES }` | 11 | `ORA-00274` at 0, naming the option: the database is no standby. In the root container `FINISH` and `USING INSTANCES` parse (`ORA-01031`) and `TO LOGICAL STANDBY` answers `ORA-00274` still |
| a virtual column of a datatype that is no type | 11 | `ORA-54003` at the datatype (15-21), as `ORA-00902` is for a column |
| `DROP TABLESPACE SET a [ INCLUDING CONTENTS ... ]` | 10 | `ORA-02557` at 20: shard DDL is disabled |
| `SPLIT ... NESTED TABLE a INTO ...` | 10 | `ORA-22912` at the column (51, 63): not a nested table |
| `ALTER PLUGGABLE DATABASE ... RECOVER ...` | 8 | `ORA-00905` at 0, in the root container too: the page derives the clause and this release's parser does not read it |
| `CREATE ROLE`/`ALTER ROLE`/`CREATE USER`/`ALTER USER ... IDENTIFIED GLOBALLY` | 8 | `ORA-00439` at 0: Enterprise User Security is not enabled |
| `ENCRYPT USING 's'`, and a tablespace's encryption spec | 8 | `ORA-28339` at the algorithm |
| `CREATE AUDIT POLICY a` with no audit option | 8 | `ORA-46373` at 0 |
| `ALTER PLUGGABLE DATABASE ... SNAPSHOT ...` | 7 | `ORA-12754` at 0: the feature's capability is missing |
| `SIZE 1 P` in `DEALLOCATE UNUSED KEEP 1 P` and `ALLOCATE EXTENT(SIZE ...)` | 6 | `ORA-03277` at the value |
| `ALTER INDEXTYPE a ADD a(NUMBER, OBJECT)`, `CREATE OPERATOR a BINDING(NUMBER, OBJECT) ...` | 6 | `ORA-00902` at the datatype (26-38) |
| `AUDIT ...` | 5 | `ORA-46401` at 6: traditional auditing is desupported |
| `SHARD`, `NOSHARD` on a sequence | 5 | `ORA-02511` at 17, 18: shard DDL is disabled |
| `XMLTYPE ... STORE AS BINARY XML`, `ALLOW NONSCHEMA` | 5 | `ORA-19001` at 47, 57 |
| `ALTER AUDIT POLICY a CONDITION 's' EVALUATE PER ...` | 4 | `ORA-46368` at 0: no simple rule condition |
| `INITRANS 1` | 4 | `ORA-02207` at the value: a value check, and `INITRANS 2` parses |
| `ALTER INDEXTYPE a.a COMPILE`, `DROP INDEXTYPE a.a`, `DROP OPERATOR a.a` | 4 | `ORA-00987` at 0: no user `a` |
| a resource named twice in `ALTER PROFILE`/`ALTER RESOURCE COST` | 4 | `ORA-02376` at the end |
| `ALTER TABLESPACE SET ... DATAFILE` | 4 | `ORA-02580` at the end |
| `SIZE 1 G` on `CREATE CLUSTER` | 4 | `ORA-02229` at 32: a value check, `SIZE 1 K` parses |
| `ALTER DIMENSION a.a COMPILE`, `CREATE MATERIALIZED VIEW a.a ...`, `CREATE MATERIALIZED ZONEMAP a.a ...` | 3 | `ORA-01435` at 0: no user `a` |
| `ALTER MATERIALIZED VIEW a REFRESH FAST USING ROLLBACK SEGMENT a` | 3 | `ORA-12061` at 45 |
| `XMLTYPE ... XMLSCHEMA 's'`, `STORE AS OBJECT RELATIONAL` | 3 | `ORA-19002` at the end: the URL or the schema is missing |
| `... ILM ADD POLICY COMPRESS ROW ...` on a tablespace | 3 | `ORA-38336` at 0: the policy is invalid for a tablespace |
| `CREATE FUNCTION a ... SHARING = { METADATA \| NONE }` | 3 | `ORA-24344` with no byte, a unit created with compilation errors: as `ORA-65021` is for a sequence, SHARING is the application root's |

Forty-three rows answer a code no other row has, or one they share with a
single other row: a date literal out of range (`ORA-01841` at 74), `UPDATE
INDEXES(a(PARTITION NOCOMPRESS))`
(`ORA-14006` at 40), a segment name (`ORA-64223` at 0), `ALTER TABLE a ENABLE
PRIMARY KEY USING INDEX a` (`ORA-03049` at 45: `USING INDEX bb` parses on the
same table and answers `ORA-02432`, so the clause takes a name and the word
`a` does not), a flashback archive with another operation (`ORA-55643` at the
end), a LOB option on a hash partition (`ORA-22877`), `MAXVALUE` twice
(`ORA-02278`), `ASSOCIATE STATISTICS ... USING NULL` (`ORA-29936`), a nested
table partitioned under a plain one (`ORA-30758`), a system-partitioned
materialized view (`ORA-32362`), a subpartition's physical attribute
(`ORA-14160`), `SCALE 1` on a zonemap (`ORA-31951`), `CREATE PFILE = FROM
SPFILE` (`ORA-02236` at 15), a mandatory profile's limit (`ORA-02387` at the
end), `CREATE SPFILE ...` in a pluggable database (`ORA-65136`), `CREATE TABLE
a PARENT a` (`ORA-02523`), `CREATE TABLE a ENABLE CONSTRAINT a`
(`ORA-02250`: the table defines no such constraint), a tablespace set's file
name (`ORA-02537`), an annotation twice (`ORA-11562`), a property graph's key
and EXCEPT list named twice (`ORA-42416`, `ORA-42426`), a dimension hierarchy
with no child (`ORA-30336`), `CREATE DATABASE a` (`ORA-01501` at 0: the judge
ran it and the database could not be created), a retention outside its range
(`ORA-05741`), Java not installed (`ORA-29538`), a partition bound that is a
keyword rather than a literal (`ORA-14019`, `ORA-14017`, `ORA-14308`), an
implementation type that is not there (`ORA-29829`), a user or role that is
not there (`ORA-01917` at 0) and a graph property named twice (`ORA-42407`).

Two rows are a bare `reject`, with no byte: `ALTER MATERIALIZED VIEW LOG ON a
SPLIT PARTITION a AT(1) PARALLEL` and its `SUBPARTITION` twin fail in Oracle's
own code (`ORA-00600` at `PTI, prspo_set_up`), which the judge reads as a
rejection with no position. Nothing says where it stopped, so neither row is
closed.

## Rows both refuse, at another byte

A row here is a statement omni and the pinned Oracle both refuse, at a
different byte: 5,298 rows, 1,021 of them where omni stops first and 3,875
where Oracle does. Omni's byte is the token its grammar cannot take. Oracle's
is where its own parser stopped, which may be a name it looked up, a check of
the parsed statement, or its own code failing, so a row is read from the
outside in: Oracle first means omni took a token Oracle does not, and omni
first means it stopped at something Oracle reads past, or reported a rewrite
point instead of the offending token.

Four hundred and two rows carry no byte at all (`reject` alone), a stored unit
compiled with errors or a statement the engine never answered, so there is
nothing to match. Four hundred and forty carry `reject@0`: Oracle's
answer at the start of the statement, which is "no such form" or a name it
looks up before it reads the rest.

#1473 moved 567 rows to Oracle's byte. No `suite:` row is among them: those
statements are the corpora's, whose cells move only with the corpus.

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| a statement whose text ends inside a group: the delimiter scan's end-of-input byte is reported only where the statement parsers read through to that group, so a statement they stop before reports its own first error | 395 | the token the parsers stop at, which Oracle refuses too |
| `ALTER TABLE a MODIFY PARTITIONSET ...` | 163 | `PARTITIONSET` (21): `ORA-02577`, before the clause is read |
| `CREATE MATERIALIZED VIEW a PARTITION RANGE(a) ...` | 9 | the word where `BY` belongs (37): a materialized view's properties enter the clause on `PARTITION` |

#1475 moved 280 more, and closed 2 rows omni accepted where Oracle
refuses (`... ORGANIZATION EXTERNAL(TYPE a) AS QUERY`). Both are readings
the grammar had not made: Oracle reads two words, or the table's column
list, before it refuses. One row Oracle accepts, `CREATE TABLE a
ORGANIZATION EXTERNAL(TYPE a` (its type is looked up), still rejects, now at
Oracle's byte 13 rather than 43, so it stays in "Rows only Oracle accepts".

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `... PARTITION BY CONSISTENT HASH (...) ...` | 217 | the token after `HASH`: Oracle reads `CONSISTENT HASH` and answers `ORA-14151` there outside a sharded database, before the columns. omni refused at `CONSISTENT` without reading them. 209 rows name `HASH`; 8 are near-misses that drop it |
| `CREATE TABLE a ORGANIZATION EXTERNAL (...)` | 63 | the table's name (13): Oracle reads the access driver type, then answers `ORA-00906` because the table has no column list, before the rest of the clause. A clause that does not start with `TYPE` is refused inside it, at its own token |

This change moved 48 more. It is a clause the page derives and omni did not
read, so Oracle's byte is where the clause it read ends and the token that
cannot follow it stands. No verdict moved: every one of the 48 was a
rejection on both sides, at another byte.

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `CREATE TABLE a(SUPPLEMENTAL LOG ...)` | 48 | the ')' of a column list that holds no column, or the token inside or after the clause that is not the one it takes: a supplemental_logging_props clause is one of relational_properties (`CREATE-TABLE.bnf` lines 58-65 and 183), so Oracle reads it and then wants a column, as it refuses a list of constraints alone at its ')' (`ORA-03079`). `CREATE TABLE a OF a(SUPPLEMENTAL LOG GROUP a(a))` is the control: an object table's `object_properties` took the clause all along and Oracle accepts it |

This change closed the derived analytic view's family, and the grammar that
generated it was wrong. `SELECT.bnf` read `sub_av_clause` as `ANALYTIC VIEW
USING analytic_view_name`, so the generator wrote `FROM ANALYTIC
VIEW(ANALYTIC VIEW USING a)`. The page draws the clause as `USING [ schema. ]
base_av_name [ hierarchies_clause ] [ filter_clauses ] [ add_meas_clause ]`
and puts `ANALYTIC VIEW` on the inline view and on the WITH element instead
(its diagrams, and its own transitory and inline analytic view examples on the
SELECT page), and the pinned Oracle reads it that way: `ANALYTIC
VIEW(ANALYTIC VIEW USING a)` is `ORA-02000`, "missing USING keyword", at the
token inside the parens (28), while `ANALYTIC VIEW(USING a)` is read to its
name lookup (`ORA-18307`) and `WITH a ANALYTIC VIEW AS(USING a)` to
`ORA-18385`, a check of the parsed statement. So the family was a reading of
the quote, not a lookup, and the quote was corrected: `sub_av_clause`,
`hierarchies_clause` (whose list is entirely optional), `filter_clauses` (one
FILTER FACT around a list of `hier_ids TO condition`), `add_meas_clause`
(whose item is `cube_meas`, a measure name and a base or calculated clause)
and `subav_factoring_clause`. `inline_analytic_view` keeps its parentheses,
which both the server and the page's example write. `analytic_inline.go` reads
the clause.

All 103 rows the family was measured with are gone, at 28 on both sides. The
corrected quote removed 221 rows and added 345: 231 of them agree (204 in an
inline analytic view, 19 in a WITH element, 8 beside them), 54 are the "only
omni accepts" rows below, and 60 are near-misses inside the clause, where the
two sides stop at different tokens. 84 rows that had a row before and after
closed, and 2 moved to another byte behind the delimiter scan.

| Cause | Rows | What moved |
| --- | --- | --- |
| `... FROM ANALYTIC VIEW(ANALYTIC VIEW USING a)` | 204 | agree: the clause starts at `USING`, is read in a SELECT and a DELETE, with an alias, an empty HIERARCHIES list, a filter and a calculated measure |
| `WITH a ANALYTIC VIEW AS(USING a) ...` | 19 | agree: the WITH element's `ANALYTIC VIEW` stands before `AS`, not inside the clause |
| `... ADD MEASURES(a)`, `... ADD MEASURES(a SUM)`, `... FILTER FACT(a, a TO 1 = 1)` | 54 | the other way: the page derives them (a `base_meas_clause` holds both of its parts optionally, and `hier_ids` is a list) and this server answers at the ')' (50, 51) or at the second hierarchy name (49), so they stay rows |

The parameter family moved. `parameters.go` carries the 566 names
`V$PARAMETER` lists on the pinned server and the 7,137 hidden names `X$KSPPI`
adds (7,700 distinct), with the single-word values `V$PARAMETER_VALID_VALUES`
lists for them, and `TestAlterParameterReading` holds the family: a name the
lists do not hold is refused at the name
(`ORA-02248` for ALTER SESSION, `ORA-02065` for ALTER SYSTEM, `ORA-00922` at 0
for ALTER SYSTEM RESET), while a name it holds is read on, so a value where the
`=` belongs is refused at the value. 136 rows moved, 51 of them from omni
accepting to both refusing; the 10 rows below are what the list cannot answer,
all at Oracle's byte 0.

The remaining rows are spread thin: 3,033 name no family below. These are the
largest, and a row may be one of several, so the counts are not disjoint.

| Family | Rows | Why the byte stays |
| --- | --- | --- |
| `DELETE ... MATCH_RECOGNIZE (...)` | 252 | the byte is inside the clause: a DML `FROM` clause refuses the word, because Oracle 23ai fails in its own code on the form there (`Parser.engineCrash`), so its answer is no token of the grammar |
| `... PARTITION BY CONSISTENT HASH (...) ...` behind a `)` | 80 | the delimiter scan answers for the `)` before the statement parsers run, so the reading above never reaches Oracle's token: 60 are the `MODIFY PARTITIONSET` family below, and 20 are near-misses that put `a)` or a `SUBPARTITION BY RANGE a)` where `(a)` belongs |
| rows whose cell is `reject@0` | 458 | Oracle's answer is at the start of the statement: no such form, or a name it looks up before reading the rest |
| rows whose cell is `reject` | 402 | a stored unit compiled with errors, or a statement the engine never answered: no position to match |
| `CREATE PLUGGABLE DATABASE a USING 's' ...` | 78 | `ORA-65125` at the file name (34), the engine's check on the name it is given |
| `suite:reference` | 76 | these statements are the corpus's, so their cells move only with the corpus |
| `... MODIFY PARTITIONSET ...` behind a `)` in the middle of the statement | 60 | `PARTITIONSET` (21): a near-miss puts a `)` where a `(` belongs, and the delimiter scan answers for it before the statement parsers run |
| `CREATE TABLE a ORGANIZATION EXTERNAL ...` behind a `)` | 2 | the clause is read and refused at the table's name, but the statement carries a stray `)` the scan answers first: `... EXTERNAL TYPE a)` and `... EXTERNAL(TYPE a LOCATION 's'))` |
| `ALTER DIMENSION ...` | 39 | the dimension's name (16): Oracle reads a complete `ADD`/`DROP` clause and then resolves the dimension, which is not there, before it reads what follows the clause, so its byte is the name and not a token (`ALTER DIMENSION a DROP LEVEL a IS(a.a)`, `... ADD LEVEL a IS(a.a) zzz`). A clause broken inside itself answers at its own token and already agrees. Of the 72 `ALTER DIMENSION` rows here, 30 more carry `reject@0` (a clause word Oracle has no form for) and 3 sit behind the scan's `)` and end-of-input answers |
| `ALTER SESSION SET DEFAULT_COLLATION ...`, `ROW VISIBILITY ...`, and `ALTER SESSION SET` with no setting | 10 | `ORA-43929` or `ORA-00922` at 0: this server runs `MAX_STRING_SIZE=STANDARD`, so `DEFAULT_COLLATION` is refused at the statement's start; `ROW VISIBILITY` is no setting, and an empty `SET` is refused whole. The parameter names themselves are read now (see below) |
| `... PARTITION BY RANGE(a)` on a `CREATE CLUSTER` | 323 | the clause's `(` (45): `ORA-14179`, because a cluster is partitioned only when it is a hash cluster. This is the check "Rows only omni accepts" measures (146 rows there), not a reading, and the near-misses that carry it stay |
| `CREATE TABLE ... ENABLE ... USING INDEX ...` | 329 | the token after the `enable_disable_clause` (34), which is `USING` here: `ORA-02432`, "cannot enable primary key - primary key not defined for table". Not a token the grammar cannot take: the check is of the table's element list, and where that token is the end of the statement the recording reads the check as an accept (`CREATE TABLE a ENABLE PRIMARY KEY` answers `ORA-02432` at 33 and agrees today). With a primary key defined the clause is read — `... (b NUMBER, PRIMARY KEY(b)) ENABLE PRIMARY KEY USING INDEX` parses — and a bare index name after `USING INDEX` is `ORA-03049` at the name. omni leaves the check to the engine, as `parser.go`'s `createTableEngineCheck` says and names ORA-02432 for, so the rows stay |
| `CREATE TABLE a (...)` with `PARTITION`, `VARRAY`, `NESTED`, `LOB` or `DOMAIN` first | 48 | the token after the word: Oracle reads `PARTITION` as a column name and the next word as its datatype, as "Rows only Oracle accepts" records (`ORA-00902` at 25 for `CREATE TABLE a (PARTITION ...)`), but for `VARRAY`, `NESTED TABLE` and `TABLE` it reads past the second word and stops at 32. Matching one server's reading of a word that is no datatype is not a reading of the page |
| `CREATE TABLE a XMLTYPE ...` | 73 | 15 for some texts and the token after the XMLType column for others: Oracle refuses the clause where it stands for a table with no column list, but reads `XMLTYPE a` before refusing when more follows it, so no single reading of the page orders them |

The scan's stray `)` answer is still not deferred. Asked with these readings
in place, returning the parsers' answer for a `)` the scan refuses moves 732
rows to Oracle's byte, but 133 to another wrong byte and 15 where omni agreed
that would stop agreeing: `CREATE VECTOR INDEX a ON a a)`, a LOB storage
option, `ALTER TABLE a ADD PERIOD FOR a)`, identity options and the
`INSERT`/`MERGE` `VALUES 1)` rows, where omni's grammar stops at the token
before the `)` and Oracle reads it. Those 15 are the ones #1473 measured, so
the deferral waits on their readings.

What is left of those is near-misses whose own token is Oracle's answer
(`ALTER MATERIALIZED VIEW a TRUNCATE a`, `ALTER SESSION SET a = 1`, a LOB
storage option, an external table's directory): each needs a reading omni does
not have, which is an accept gap rather than a position.

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
