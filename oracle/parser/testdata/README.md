# Oracle parser test data

## `compat.tsv`: Oracle's verdicts

`compat.tsv` holds a statement for each alternative of each rule in
`../../testdata/OracleDocs.g4`, then its near-misses, with omni's verdict and
the pinned Oracle's (`../../testdata/pins.tsv`). `TestCompatibility`
(`../compat_test.go`) reads it through `enginetest/compat` and starts no
server. Recorded 2026-11-14 in 35 minutes, most of them on the statements
Oracle never answers (below); a second recording on a fresh server wrote the
same file. A rejection carries the byte Oracle reports it at. The counts here
are the file as it stands at judge revision 3, after #1473, #1475, #1480, #1483,
#1486, #1530, Wave 2 (unclosed-`(` deferral and follow-up grammar readings),
Wave 3 (partition/nested-table/column-datatype/DML/DDL grammar readings),
Wave 4 (inline external table driver types and remaining grammar readings), and
Wave 5 (intra-statement duplicate checks, lexer/ORDER BY offsets, and remaining
grammar readings):
45,852 generated rows and 1,010 hand rows, 46,862 in all.

| Rows | Total | Same cell | Both reject, another byte | Only Oracle accepts | Only omni accepts |
| --- | --- | --- | --- | --- | --- |
| One per alternative | 14,534 | 13,680 | 349 | 4 | 501 |
| Near-misses | 31,318 | 29,491 | 1,688 | 100 | 39 |
| Hand rows | 1,010 | 911 | 80 | 17 | 2 |

At the recording, most of what only Oracle accepts was materialized views
(2,018 rows of `CREATE` and `ALTER MATERIALIZED VIEW`), `ALTER TABLE` and
`CREATE TABLE` (913), `GRANT` and `REVOKE` (382) and `SELECT` (189). Most of
what only omni accepted was `ALTER TABLE` (1,881), materialized view logs
(1,072), databases, pluggable databases and tablespaces (3,566) and `CREATE
TABLE` (648): forms omni read loosely. No parser gap is fixed here.

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
| `suite:lexer` | 77 | 76 | 1 | 0 | 0 |
| `suite:reference` | 321 | 241 | 64 | 16 | 0 |
| `suite:manifest` | 612 | 594 | 15 | 1 | 2 |

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

The 128 rows whose cells are `reject` where Oracle says `accept` were each
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
| `CREATE TABLE a (PARTITION ...)` | 4 | `ORA-00902` at 25: `PARTITION` is a column, the next word its datatype |
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

The other direction: the 696 rows whose cells are `accept` where the pinned
Oracle says `reject`. Each was measured the same way, at judge revision 3,
with the code `DBMS_SQL.PARSE` answers and `LAST_ERROR_POSITION`'s byte. The
question a row asks is whether Oracle's refusal is the grammar's: the page
`bnf/` quotes under it is the grammar omni must read, so a statement the page
derives stays accepted even where this server refuses it, and is recorded
here instead of closed. Fifteen rows were text the page refuses too, and were
closed, and so were the thirty-two `PARTITIONSET` rows: the page's clause is
only reachable in a composite-sharded database, which `OracleDocs.g4` records
where the converter could drop the rule (its `fix:` notes on
`alter_table_partitioning_clauses`' neighbours) and this server is not.
#1475, #1480 and #1483 closed rows of this direction too: an external table's
`AS QUERY` (2), the `ALTER SYSTEM SET`, `ALTER SESSION SET` and `ALTER SYSTEM
RESET` rows whose parameter names `oracle/parser` now carries (51), and the
derived analytic view the `SELECT.bnf` quote drew wrongly. All 696 here are a
name Oracle looks up, a check of the parsed statement, or the server's own
setup, and stay rows.

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
from the grammar; the controls are named where a reader would ask. These are
the causes measured, largest first: they are the largest groups, not a
partition of the 696, since one row can answer more than one cause and the
long tail below holds the rest.

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

Forty-two rows answer a code no other row has, or one they share with a
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
different byte: 2,969 rows, 475 of them where omni stops first and 2,092
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
corrected quote removed 221 rows and added 345: 321 of them agree (294 in an
inline analytic view, 19 in a WITH element, 8 beside them) and 24 are
near-misses inside the clause, where the two sides stop at different tokens.

| Cause | Rows | What moved |
| --- | --- | --- |
| `... FROM ANALYTIC VIEW(ANALYTIC VIEW USING a)` | 294 | agree: the clause starts at `USING`, is read in a SELECT and a DELETE, with an alias, an empty HIERARCHIES list, a single-hierarchy filter and a calculated or `AGGREGATE BY` measure |
| `WITH a ANALYTIC VIEW AS(USING a) ...` | 19 | agree: the WITH element's `ANALYTIC VIEW` stands before `AS`, not inside the clause |

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

#1486 moved 241 more to Oracle's byte (221 where omni stopped first and
20 where Oracle did), with no verdict change and no `suite:` row moved.

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `... FROM a CONTAINERS(...)` and `... SHARDS(...)` in `SELECT` and `DELETE` | 44 | the `(` (26 and 22): `CONTAINERS` and `SHARDS` open a clause only in place of a table; after one they name its alias |
| `CREATE TABLE a COLUMN ...` | 32 | inside or after `object_type_col_properties` (`ELEMENT` is refused at 24 with `ORA-00905`, unlike on a collection column; `CREATE TABLE a COLUMN a IS OF(a)` stays in "Rows only Oracle accepts" at 15), plus `COLUMN COMPRESS` (22) and `COLUMN STORE COMPRESS FOR` (40) |
| `MODIFY [SUB]PARTITION a VARRAY ...` in `ALTER TABLE` and `ALTER MATERIALIZED VIEW` | 30 | where `STORE` belongs (41, 44, 53, 56) or at the missing item: a partition's `VARRAY` reads `varray_col_properties`, not `(modify_LOB_parameters)` |
| `CREATE TABLE` and `ALTER TABLE` constraint, `DEFAULT COLLATION` and `OF ... NOT` near-misses | 29 | the unconsumed token before an unclosed column list's `)` (16), after `REFERENCES ... ON` where `DELETE` belongs (5), after `REF(a) WITH` or inline `WITH` where `ROWID` belongs (6), `DEFAULT a` at 23 (1), and `OF a NOT zzz` at 24 (1) |
| `ALLOW` and `DISALLOW` in `ALTER TABLE` and `ALTER MATERIALIZED VIEW` | 20 | the token where `CLUSTERING` belongs: Oracle reads `allow_disallow_clustering` in `ALTER TABLE` properties, its `OVERFLOW`, and `MODIFY [SUB]PARTITION` |
| `INSERT INTO a SELECT 1 FROM dual RETURNING ...` | 20 | the token after `RETURNING` (43): `RETURNING` is not reserved in a query and is read as `dual`'s alias |
| `CREATE VIEW` and `ALTER VIEW` out-of-line `FOREIGN KEY ... ON ...` | 20 | `ON` (45): a view's declarative foreign key takes no `ON DELETE` (`ORA-00922`) |
| `... INMEMORY NO ...` in `CREATE`/`ALTER TABLE`, `ILM` and `CREATE MATERIALIZED VIEW` | 13 | the token after `NO`, where `MEMCOMPRESS` or `DUPLICATE` belongs |
| `CREATE [SHARED] [PUBLIC] DATABASE LINK` and `AUTHENTICATED BY` | 11 | `USING` on a `SHARED` link without `AUTHENTICATED BY` (30), the non-`DATABASE` token after `SHARED` (14), and the token after `AUTHENTICATED` on a non-`SHARED` link (37, `ORA-00905`) |
| `CREATE SYNONYM a SHARING = ...` | 8 | the token after `=` when not `NONE` (`ORA-65021`, `ORA-65014`), or where `=` belongs (25) |
| `DELETE`/`UPDATE` `(SELECT 1 FROM dual WITH ...)` | 7 | the token after `WITH` (31, 38), where `READ` or `CHECK` belongs |
| `CREATE NOFORCE ...` and `CREATE VIEW ... CONSTRAINTS` | 7 | the token after `NOFORCE` when no view or analytic object follows (5), and past `CONSTRAINTS` on a view constraint (2) |

The stray-`)` deferral and follow-up grammar wave moved 1,005 more rows to
Oracle's byte (241 where omni stopped first and 764 where Oracle did), with no
verdict change and no `suite:` row moved:

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| delimiter scan stray-`)` deferral (`ParseRange` reports the statement parser's error when it stops at or before a stray `)`, after fixing the 13 prerequisite readings: `CREATE VECTOR INDEX ... ON table [alias]`, `ALTER MATERIALIZED VIEW MODIFY LOB`, bare `ALTER TABLE ADD PERIOD FOR`, `LOB ... STORE AS <unreserved>` segment name, bare `AS IDENTITY` options, `CREATE TABLE XMLTYPE`, and `INSERT`/`MERGE` `VALUES expr)`) | 761 | the earlier token the statement parsers stop at, closing all remaining `MODIFY PARTITIONSET`, `CONSISTENT HASH`, and `ORGANIZATION EXTERNAL` rows that sat behind a `)` |
| `CREATE TABLE a XMLTYPE ...` | 69 | inside or after `XMLType_column_properties`, or at `XMLTYPE` (15) when complete with no column list |
| `SUBPARTITION` missing `BY` / `TEMPLATE` and non-hash `SUBPARTITION TEMPLATE` quantity | 63 | the token after `SUBPARTITION` or after the quantity (`PARTITIONS` / `zzz`) |
| `ALTER TABLE` and `ALTER MATERIALIZED VIEW` partition verbs (`SPLIT`, `MERGE`, `EXCHANGE`, `MOVE`, `RENAME`, `SET`), `DEPENDENT TABLES`, `SPLIT ... INCLUDING ROWS`, and `MODIFY [SUB]PARTITION LOB` | 28 | the missing keyword after the verb or clause word |
| `CREATE`/`ALTER TABLE` and `CREATE MATERIALIZED VIEW` `JSON`/`LOB`/`ROW`/`READ`/`TABLESPACE`/`ENABLE LOGICAL REPLICATION`/`USING INDEX`/`EXCEPTIONS INTO`/`DEFAULT ON`/`AS (expr)`/`ENCRYPT NO` | 22 | the missing keyword or unconsumed token inside the clause |
| `ALTER MATERIALIZED VIEW` `INMEMORY DUPLICATE FOR` and `CREATE MATERIALIZED VIEW` `INMEMORY(a) INMEMORY, ...` | 16 | the token after `FOR` (49) or after `,` (48, 49) |
| `ASSOCIATE STATISTICS ... USING ... WITH ...` | 15 | the token after `WITH` (51) |
| `ILM ADD POLICY` (`ROW`, `COLUMN`, `MODIFY`, `NO`, `COMPRESS ROW`) | 10 | the missing `STORE`, `INMEMORY`, or token after `ROW` |
| `CREATE PLUGGABLE DATABASE` (`USER`, proxy `HOST`/`PORT`) | 7 | the token after `USER`, `HOST`, or `PORT` |
| `SET` statement and `ALTER DISKGROUP a REPAIR ...` | 9 | the token after `SET` (3, 4) or `REPAIR` (25) |
| `SELECT`/`DELETE`/`UPDATE` `APPLY` alias, `FROM TABLE 1)`, and inline `ANALYTIC VIEW(USING a ADD(...))` | 5 | the token after `APPLY`, `TABLE`, or `ADD` |

Wave 2 moved 1,048 more rows to Oracle's exact cell (85 net fewer `omni_first`,
962 fewer `oracle_first`, and 1 `CREATE TABLE a(PARTITION READ)` moved from
only-Oracle-accepts to agreeing `accept`, bringing exact byte match from 89.67%
to 91.91%):

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| delimiter scan unclosed-`(` deferral (`ParseRange` reports the statement parser's error whenever it stops at or before `scanErr`, after fixing the prerequisite readings: `MOVE KEYS ... WITH IDENTIFIER IN (...)`, `ALTER MATERIALIZED VIEW MODIFY LOB(a, ...)`, repeated `LOB_partition_storage` clauses, `LOCAL(PARTITION PARAMETERS('s')` EOF handling, `CREATE TABLE` table-level `domain_clause`, `SUPPLEMENTAL LOG GROUP ... NO LOG`, hash partition `PARTITION ILM` at EOF, and `CREATE TABLE (col READ)` datatype reading) | 636 | the earlier token the statement parsers stop at inside an unterminated parenthesis group (closing `USING TEMPLATE(DEFAULT DEFAULT ...)`, inline `EXTERNAL` column clauses, missing `)` before `SUBPARTITION BY`, `CREATE MATERIALIZED VIEW a(PARTITION ...)`, `SELECT * WINDOW a AS(ROWS ...)`, `ALTER LOCKDOWN PROFILE`, and `INSERT INTO t (a VALUES (1)` in `suite:manifest`) |
| `CREATE TABLE` and `CREATE MATERIALIZED VIEW` `VARRAY ... STORE AS LOB seg (PARTITION ...)` | 135 | `(` (38, 50): `LOB_partition_storage` is refused at `(` on `VARRAY ... STORE AS LOB seg` here, while non-`PARTITION` parameters keep failing inside `(` |
| `ALTER MATERIALIZED VIEW a ADD PARTITION ...` | 91 | the token after `ADD PARTITION` (40) when a bare partition name has an unexpected trailing word or an unprefixed `UPDATE [GLOBAL] INDEXES` fails to complete |
| `INMEMORY DISTRIBUTE BY ...` in `CREATE`/`ALTER TABLE` and `CREATE`/`ALTER MATERIALIZED VIEW` | 32 | `BY` (34, 35, 46, 47) when not followed by `ROWID RANGE`, `PARTITION`, or `SUBPARTITION` |
| `CREATE MATERIALIZED VIEW a PARTITION BY REFERENCE(a)` | 24 | the constraint name `a` (50) once `)` closes (`ORA-14655`) |
| `CREATE TABLE` and `CREATE MATERIALIZED VIEW` `attribute_clustering_clause` | 17 | `BY` / `ON` / `ORDER` (25–53) via full BNF parsing of `CLUSTERING [ table [ alias ] ] [ JOIN ... ON (...) ] BY [ LINEAR \| INTERLEAVED ] ORDER (...)` |
| `INMEMORY`, `NOMAPPING`, `IMMEDIATE`, and `INMEMORY MEMOPTIMIZE FOR` readings on `CLUSTER`, `INDEX`, `MATERIALIZED VIEW`, and `MATERIALIZED VIEW LOG` | 17 | the token after `INMEMORY`, `NOMAPPING`, `IMMEDIATE`, or `MEMOPTIMIZE FOR` |
| `ADMINISTER KEY MANAGEMENT MOVE KEYS ... WITH IDENTIFIER IN(...)` | 14 | end of input when a parenthesized `WITH IDENTIFIER IN (...)` is not a subquery followed by `WITH BACKUP` (`ORA-00921`) |
| `CREATE TABLE a(DOMAIN ...)` table-level `domain_clause` | 12 | `)` (26, 28, 29) when no column is defined, or the token inside/after `DOMAIN [ schema. ] name ( col [, ...] )` |
| `ALTER DISKGROUP a` (`MOVE FILE`, `MOVE FILEGROUP`, `REPLACE DISK`, `MODIFY QUOTAGROUP`) | 9 | the token after the keyword read as an unquoted object/file name (`FILEGROUP`, `a`, `QUOTAGROUP`, `'s'`, `QUOTA`) |
| `ALTER VIEW a ADD` (`CHECK`, `(`, reserved `CONSTRAINT` name, `PRIMARY`/`FOREIGN` missing `KEY`) | 9 | `CHECK` (17), inside `(` (17), reserved constraint name (28), or `(` missing `KEY` (24) |
| `ALTER TABLE a MODIFY(a <type>(...))` (`LOB(a)`, `LOBS(a)`, `CHECKPOINT(...)`) | 9 | the user-defined type name (`ORA-03060`: parenthesized modifiers belong to built-in parameterized types) |
| `ALTER MATERIALIZED VIEW` `TRUNCATE` / `MODIFY LOB(a, ...)` and `ALTER TABLE UPGRADE` repeated `LOB_partition_storage` | 14 | `TRUNCATE` (26), after `,` in `MODIFY LOB(a, ...)` (39, 40), or inside/at `)` of the second `VARRAY` storage clause |
| `ALTER TABLE a NO ...`, `CONTAINER_MAP UPDATE`, `SUPPLEMENTAL LOG DATA` / `ANNOTATIONS` duplicates, `MODIFY PARTITION ADD VALUES LESS THAN`, hash `TABLESPACE SET`, `ENCRYPT NO zzz SALT`, `CREATE TABLE (PARTITION READ ...)`, and `SELECT ... VERSIONS FOR` | 29 | `NO` (14), `VALUES` / EOF in `CONTAINER_MAP UPDATE`, the token after duplicate `ALL`/`a`, `LESS` (44, 56), `SET` (47, 72), `zzz` (40), `ONLY` (30), or `a` (29) |

Wave 3 moved 514 more rows to Oracle's exact cell (57 fewer `omni_first`, 435 net
fewer `oracle_first`, 5 fewer only-Oracle-accepts, and 17 fewer only-omni-accepts,
bringing exact byte match from 91.91% to 93.00% and verdict agreement from 98.24%
to 98.29%):

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `CREATE TABLE` / `CREATE MATERIALIZED VIEW` `NESTED TABLE ... STORE AS storage ( ...` | 101 | the nested item name when the storage table's parenthesized group is unclosed (`skipParenthesizedBalanced`) |
| `CREATE`/`ALTER TABLE` & `CREATE VIEW` column datatypes, constraints, and properties (non-reserved keyword UDT datatypes, user-defined type `(`, `GENERATED AS(`, `DEFAULT ON NULL INSERT`, duplicate identity options, `SUPPLEMENTAL LOG GROUP(`, `RENAME LOB`, `DROP(PERIOD`, `ENABLE LOGICAL REPLICATION`, `MODIFY CLUSTERING`, `INMEMORY NO ...`, `ILM ENABLE/DISABLE`, `CONSTRAINT` missing body, `GLOBAL`/`PRIVATE` table modifiers) | 96 | the column datatype/clause keyword or missing `old_lob_name` / `KEYS` / `POLICY` / constraint keyword |
| Query & DML clauses (`SELECT`, `WITH`, `INSERT`, `UPDATE`, `DELETE`, `MERGE`: required `INTO`, reserved target/CTE/SET names, `RETURNING INTO`, `AS SCN`/`TIMESTAMP`/`PERIOD`, `START WITH` without `CONNECT BY`, `SAMPLE BLOCK`, unclosed `PARTITION(`/`SUBPARTITION(`/`TABLE()(`/`(SELECT`) | 70 | the missing keyword, reserved word, or opening `(` / clause word |
| Partition & subpartition clauses (`ALTER TABLE`/`MVIEW`/`MVLOG`/`INDEX` `ADD`/`MODIFY`/`MERGE` partition/subpartition bare name & `COMPRESS` trailing-token blame, `VALUES` identifier rejection, `OVERFLOW TABLESPACE`, `MVIEW` `SUBPARTITION TEMPLATE ... INDEXING`) | 69 | the partition/subpartition name, `COMPRESS`, `OVERFLOW`, `INDEXING`, `UPDATE`, or `)` after an invalid bound identifier |
| Reserved object, schema, user, and dblink name checks (`DROP`/`ALTER ... EXISTS`, `CREATE DIRECTORY`/`PACKAGE`/`TYPE BODY`/`SCHEMA`/`JSON DUALITY VIEW`/`USER`/`DATABASE LINK`) | 50 | the reserved keyword in the object/schema/user name position |
| Utility statements (`LOCK TABLE`, `CALL`, `FLASHBACK TABLE`, `RENAME`, `ANALYZE`, `EXPLAIN`, `COMMENT ON COLUMN`) | 31 | the reserved name, routine start, `PARTITION`/`SUBPARTITION`, `BEFORE`, or missing `(` / `:` |
| `CREATE TABLE` / `CREATE MATERIALIZED VIEW` / `CREATE MATERIALIZED VIEW LOG` `PARTITION BY HASH` attributes | 21 | the bare hash partition name or `[NO]COMPRESS` when followed by an unexpected non-EOF token, and `LOB` storage start |
| `GRANT ALL` and non-privilege column-list checks | 15 | `ALL` at EOF / `PROCEDURE`, or the non-privilege name before `(` |
| `OPAQUE`, `XMLTYPE`, `XMLSCHEMA`, and `MVIEW REFRESH START WITH`/`NEXT` clauses | 10 | the missing `TYPE`/`ALL`, reserved `AS` in `XMLTYPE`/`REFRESH`, or `STORE` |
| `CREATE ... IF EXISTS` and `IF NOT <non-EXISTS>` | 10 | `IF` across table, procedure, function, package, trigger, type, index, view, materialized view, zonemap, and admin DDL |
| `ALTER TABLE`/`MVIEW`/`MVLOG` `SPLIT [SUB]PARTITION INTO` and `SPLIT PARTITION ... NESTED TABLE` | 8 | `PARTITION` when `(` is missing after `INTO`, `,` on `SPLIT SUBPARTITION ... INTO`, or the nested table column when `INTO` is missing |
| Other DDL (`PDB` `AS SEED`/`CONTAINER_MAP`/`CREATE SNAPSHOT`/`SET NOLOGGING`, `CREATE MLE MODULE`, `INMEMORY JOIN GROUP`, `KEY MANAGEMENT`, `DOMAIN`, `INDEX LOCAL(PARTITION PARAMETERS('s'`, `OPERATOR BINDING`, `RESTORE POINT AS OF`, `AUDIT POLICY BY WITH`, `PROPERTY GRAPH`) | 33 | the clause keyword, object/partition name, or `'s'` literal at EOF |

Wave 4 moved 150 more rows to Oracle's exact cell (32 fewer `omni_first`, 66 net
fewer `oracle_first`, 2 fewer only-Oracle-accepts, and 50 fewer only-omni-accepts
including 2 rows narrowed from `accept` to both-reject, bringing exact byte match
from 93.00% to 93.32% and verdict agreement from 98.29% to 98.40%):

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `SELECT`/`DELETE ... FROM EXTERNAL((...) TYPE <driver> ...)` | 95 | the invalid driver name (`ORA-00922`: inline external tables accept `ORACLE_LOADER`, `ORACLE_DATAPUMP`, `ORACLE_HDFS`, `ORACLE_HIVE`, `ORACLE_BIGDATA`); 2 more duplicate-column rows moved from `accept` to both-reject |
| `ALTER TABLE MODIFY PARTITION BY` & `CREATE INDEX GLOBAL PARTITION BY` expression column lists | 9 | `VALUES` inside an unclosed `(a(PARTITION VALUES ...)` (matching `CREATE TABLE`) |
| `CREATE MATERIALIZED VIEW ... XMLTYPE` (`XMLSCHEMA` name vs string literal, `ALLOW`/`DISALLOW NONSCHEMA`/`ANYSCHEMA`) | 9 | `'s'`, `a`, `ALLOW`, or `DISALLOW` (matching `ALTER TABLE`) |
| `ALTER TABLE UPGRADE INCLUDING DATA VARRAY ... STORE AS LOB seg (...)` partition storage | 6 | the unexpected token inside `(` (`LOGGING`, `PARTITIONING`, `a`) or inside `PARTITION a ...` |
| `ALTER TABLE MODIFY (a AS(1), a ...)` duplicate column name check | 5 | the repeated column name `a` (`ORA-00957`) before reading its trailing properties |
| `CREATE TABLE` / `CREATE MATERIALIZED VIEW` `NESTED TABLE ... LOCAL STORE AS` | 4 | `LOCAL` (matching `ALTER TABLE`, which allows `GLOBAL` only) |
| `CREATE TABLE` / `CREATE MATERIALIZED VIEW` / `CREATE MATERIALIZED VIEW LOG` empty `SUBPARTITION TEMPLATE ()` | 3 | accepts empty `SUBPARTITION TEMPLATE ()` on `MVIEW`/`MVLOG` (2 `only_oracle` closed) and reaches missing-columns blame on `CREATE TABLE` |
| `CREATE`/`ALTER INDEXTYPE` `WITH ARRAY DML (a, a), ...` and `WITHOUT SYSTEM MANAGED STORAGE TABLES` | 3 | the missing `(` after `,` or `SYSTEM` after `WITHOUT` |
| Query & DML table-alias and expression readings (`PERIOD`, `WINDOW` in `DELETE`, `LEFT` before `PARTITION BY`, `LATERAL` without `(`, reserved `ADD` in SQL expressions) | 8 | `FOR`, `a`, `1`, `SELECT`, or `ADD` |
| Other DDL (`ALTER TABLE` top-level `NESTED TABLE`/`VARRAY` column properties, `ADD LOG DATA(ALL)`, `CREATE TABLE` duplicate `ANNOTATIONS` closing `)`, `CREATE SNAPSHOT`, `ALTER MVIEW MERGE SUBPARTITIONS INVALIDATE`, `ALTER MVLOG SPLIT PARTITION INTO`, `MOVE KEYS WITH IDENTIFIER IN a)`) | 8 | `RETURN`, `(`, `DATA`, closing `)`, `a`, `SUBPARTITION`, `PARTITION`, or `a` |

Wave 5 moved 349 more rows to Oracle's exact cell (97 fewer `omni_first`, 165
fewer `oracle_first`, and 87 fewer only-omni-accepts, bringing exact byte match
from 93.32% to 94.07% and verdict agreement from 98.40% to 98.59%):

| Cause | Rows | Where Oracle stops |
| --- | --- | --- |
| `CREATE PROPERTY GRAPH` duplicate element tables/aliases, `KEY`, `EXCEPT`, and `PROPERTIES` | 91 | the repeated table/alias (`ORA-42404`), key column (`ORA-42416`), excluded column (`ORA-42426`), or property (`ORA-42407`) |
| `ALTER TABLE` duplicate column lists (`UNIQUE`, `PRIMARY KEY`, `REFERENCES`, `INMEMORY`, `LOB`, `PARTITION BY`, `SUPPLEMENTAL LOG GROUP`), virtual column non-built-in datatypes, `FORCE(1)` partition key, `INCLUDING ROWS`, `IMMEDIATE`, `ADD DEFERRABLE`, `UPGRADE INCLUDING DATA(`, and reserved table name | 68 | the repeated name, virtual column name (`ORA-54003`), `(` after `FORCE`, missing `WHERE`, `ADD`, or trailing token |
| `CREATE TABLE` duplicate column/key lists, bare column duplicates, virtual column datatypes & `ALWAYS AS IDENTITY`, range bound arity, `CONSTRAINT` body, `SUPPLEMENTAL GROUP`, `ENCRYPT IDENTIFIED BY)`, `IMMEDIATE`/`EXTERNAL`, and `VARRAY`/`NESTED TABLE ... STORE AS` | 37 | the repeated column, virtual column name, extra bound expression, missing keyword, or `)` |
| `CREATE` / `ALTER MATERIALIZED VIEW` & `LOG` duplicate column lists, `ALLOCATE EXTENT(SIZE)`, `SCOPE FOR`, `MERGE SUBPARTITIONS FORCE(1)` / trailing token, and `ENCRYPT IDENTIFIED BY)` | 31 | the first repeated alias (for `CREATE MATERIALIZED VIEW`) or second column name, `)`, `FORCE`, or trailing token |
| Utility statements (`LOCK TABLE`, `FLASHBACK TABLE`, `TRUNCATE`) | 25 | the token after an unrecognized lock mode, `PARTITION` on multi-partition `LOCK TABLE`, missing `TO`, or missing `TABLE`/`CLUSTER`/`MATERIALIZED` |
| Query & DML (`EXTERNAL` duplicate columns, `INSERT ... SET` duplicate target columns, `FETCH FIRST 1 zzz ROW`, `UPDATE ... FOR(1)`) | 19 | the repeated column (`ORA-00957`), `ROW`, or `(` after `FOR` |
| Other DDL (`ALTER DISKGROUP` volume/quotagroup, `ALTER DATABASE`, `ALTER PLUGGABLE DATABASE`, `CREATE USER`, `ADMINISTER KEY MANAGEMENT`, `CREATE AUDIT POLICY`, `CREATE INMEMORY JOIN GROUP`, `CREATE [OR REPLACE] JAVA SOURCE IF [NOT] EXISTS`) | 18 | the clause keyword, literal, `IF`, or trailing token |
| `CREATE` / `ALTER VIEW` duplicate aliases/annotations/constraints, `CREATE MATERIALIZED ZONEMAP`, `CREATE VECTOR INDEX INCLUDE`, `CREATE CLUSTER`, `ALTER INDEX`, `CREATE DIMENSION` duplicate `LEVEL`, `CREATE`/`ALTER PROFILE` & `ALTER RESOURCE COST` duplicate limits | 39 | the repeated alias/column/level/limit name, `UNUSABLE`, or constraint keyword |
| `CREATE`/`ALTER OPERATOR` and `CREATE`/`ALTER INDEXTYPE` `OBJECT` datatype refusal & `CREATE PFILE`/`SPFILE` syntax | 10 | `OBJECT` (`ORA-00902`), `FROM`, or EOF after `AS` |
| `suite:lexer` error offsets (unterminated string/comment/q-quote, `""`, reserved `DATE` in expressions) & parenthesized trailing `ORDER BY` (`suite:reference`) | 11 | the opening lexeme byte, `DATE`, or `ORDER` (`ORA-00907`) |

The remaining 2,780 disagreeing rows are non-grammar causes (99 `suite:*` corpus
rows, 402 bare `reject` without position, 440 `reject@0` lookups/sub-parsers,
121 measured `only_oracle` lookups/checks, 542 measured `only_omni`
engine/setup/value/duplicate checks including `CREATE TABLE`/`ALTER TABLE ADD`/`DROP`
and `CREATE INDEX` duplicates intentionally left to semantic review, and ~1,176
`another_byte` rows from the documented engine checks below). These are the
largest, and a row may be one of several, so the counts are not disjoint.

| Family | Rows | Why the byte stays |
| --- | --- | --- |
| `DELETE ... MATCH_RECOGNIZE (...)` | 235 | the byte is inside the clause: a DML `FROM` clause refuses the word, because Oracle 23ai fails in its own code on the form there (`Parser.engineCrash`), so its answer is no token of the grammar |
| rows whose cell is `reject@0` | 440 | Oracle's answer is at the start of the statement: no such form, or a name it looks up before reading the rest |
| rows whose cell is `reject` | 402 | a stored unit compiled with errors, or a statement the engine never answered: no position to match |
| `CREATE PLUGGABLE DATABASE a USING 's' ...` | 111 | `ORA-65125` at the file name (34), the engine's check on the name it is given |
| `suite:reference` | 76 | these statements are the corpus's, so their cells move only with the corpus |
| `CREATE FLASHBACK ARCHIVE ...` and `ALTER FLASHBACK ARCHIVE ...` | 62 | the token after the offending one, or the end of the statement: Oracle's archive parser advances past the token it refuses before raising `ORA-55603` or `ORA-55613` (`ALTER FLASHBACK ARCHIVE a SETS DEFAULT` at `DEFAULT`, 31; `... QUOTA 1 K RETENTION 1 DAY` at `RETENTION`, 50), as `FLASHBACK TABLE` (`ORA-02000`) and `LOCK TABLE ... IN ... MODE` (`ORA-01737`) do on theirs |
| `CREATE TABLE a ORGANIZATION INDEX ...` | 52 | `INDEX` (28): `ORA-25175`, "no PRIMARY KEY constraint found", checked on the table's element list as soon as `ORGANIZATION INDEX` is read and before the rest of `index_org_table_clause` (`CREATE TABLE a (id NUMBER PRIMARY KEY) ORGANIZATION INDEX` parses, and `CREATE TABLE a (a NUMBER) ORGANIZATION INDEX` answers `ORA-25175` at `INDEX` too) |
| `CREATE VIEW a OF XMLTYPE ...` | 34 | `'s'` (35) on the 31 `XMLSCHEMA 's'` rows (`ORA-19002`, the schema URL looked up before `WITH OBJECT IDENTIFIER` or `AS subquery` is read, so even the page's valid forms carry `reject@35`), and 25 or 43 on the 3 near-misses without `'s'` (`CREATE VIEW a OF XMLTYPE WITH OBJECT IDENTIFIER DEFAULT AS SELECT 1 FROM dual`, with no `XMLSCHEMA`, parses to `ORA-22980`) |
| `ALTER DIMENSION ...` | 70 | the dimension's name (16): Oracle reads a complete `ADD`/`DROP` clause and then resolves the dimension, which is not there, before it reads what follows the clause, so its byte is the name and not a token (`ALTER DIMENSION a DROP LEVEL a IS(a.a)`, `... ADD LEVEL a IS(a.a) zzz`). A clause broken inside itself answers at its own token and already agrees. Of the 70 `ALTER DIMENSION` rows here, 30 carry `reject@0` (a clause word Oracle has no form for) |
| `ALTER SESSION SET DEFAULT_COLLATION ...`, `ROW VISIBILITY ...`, and `ALTER SESSION SET` with no setting | 6 | `ORA-43929` or `ORA-00922` at 0: this server runs `MAX_STRING_SIZE=STANDARD`, so `DEFAULT_COLLATION` is refused at the statement's start; `ROW VISIBILITY` is no setting, and an empty `SET` is refused whole. The parameter names themselves are read now (see above) |
| `... PARTITION BY RANGE(a)` on a `CREATE CLUSTER` | 323 | the clause's `(` (45): `ORA-14179`, because a cluster is partitioned only when it is a hash cluster. This is the check "Rows only omni accepts" measures (146 rows there), not a reading, and the near-misses that carry it stay |
| `CREATE TABLE ... ENABLE ... USING INDEX ...` | 329 | the token after the `enable_disable_clause` (34), which is `USING` here: `ORA-02432`, "cannot enable primary key - primary key not defined for table". Not a token the grammar cannot take: the check is of the table's element list, and where that token is the end of the statement the recording reads the check as an accept (`CREATE TABLE a ENABLE PRIMARY KEY` answers `ORA-02432` at 33 and agrees today). With a primary key defined the clause is read — `... (b NUMBER, PRIMARY KEY(b)) ENABLE PRIMARY KEY USING INDEX` parses — and a bare index name after `USING INDEX` is `ORA-03049` at the name. omni leaves the check to the engine, as `parser.go`'s `createTableEngineCheck` says and names ORA-02432 for, so the rows stay |
| `CREATE TABLE a (...)` with `PARTITION`, `VARRAY`, `NESTED`, `LOB` or `DOMAIN` first | 59 | the token after the word: Oracle reads `PARTITION` as a column name and the next word as its datatype, as "Rows only Oracle accepts" records (`ORA-00902` at 25 for `CREATE TABLE a (PARTITION ...)`), but for `VARRAY`, `NESTED TABLE` and `TABLE` it reads past the second word and stops at 32. Matching one server's reading of a word that is no datatype is not a reading of the page |

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
