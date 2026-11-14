-- pin: default gvenzl/oracle-free:23.26.3-slim-faststart sha256:f5ff19033860d662c821cb04eb10483fa94f14f78eae252d054291ea07028093 23.26.3.0.0

== TestOracleNonReservedKeywordsAsColumns default#1
> query
|
| SELECT keyword FROM v$reserved_words
|  WHERE keyword IS NOT NULL
|    AND (reserved = 'Y' OR res_type = 'Y' OR res_attr = 'Y' OR res_semi = 'Y')
|  ORDER BY keyword
< result
< column "KEYWORD" "NCHAR"
< row
: string
| !
< row
: string
| &
< row
: string
| (
< row
: string
| )
< row
: string
| *
< row
: string
| +
< row
: string
| ,
< row
: string
| -
< row
: string
| .
< row
: string
| /
< row
: string
| :
< row
: string
| <
< row
: string
| <#>
< row
: string
| <->
< row
: string
| <=>
< row
: string
| =
< row
: string
| >
< row
: string
| @
< row
: string
| ACCESS
< row
: string
| ADD
< row
: string
| ALL
< row
: string
| ALTER
< row
: string
| AND
< row
: string
| ANY
< row
: string
| AS
< row
: string
| ASC
< row
: string
| AUDIT
< row
: string
| BETWEEN
< row
: string
| BY
< row
: string
| CHAR
< row
: string
| CHECK
< row
: string
| CLUSTER
< row
: string
| COLUMN
< row
: string
| COLUMN_VALUE
< row
: string
| COMMENT
< row
: string
| COMPRESS
< row
: string
| CONNECT
< row
: string
| CREATE
< row
: string
| CURRENT
< row
: string
| DATE
< row
: string
| DECIMAL
< row
: string
| DEFAULT
< row
: string
| DELETE
< row
: string
| DESC
< row
: string
| DISTINCT
< row
: string
| DROP
< row
: string
| ELSE
< row
: string
| EXCEPT
< row
: string
| EXCLUSIVE
< row
: string
| EXISTS
< row
: string
| FILE
< row
: string
| FLOAT
< row
: string
| FOR
< row
: string
| FROM
< row
: string
| GRANT
< row
: string
| GROUP
< row
: string
| HAVING
< row
: string
| IDENTIFIED
< row
: string
| IMMEDIATE
< row
: string
| IN
< row
: string
| INCREMENT
< row
: string
| INDEX
< row
: string
| INITIAL
< row
: string
| INSERT
< row
: string
| INTEGER
< row
: string
| INTERSECT
< row
: string
| INTO
< row
: string
| IS
< row
: string
| LEVEL
< row
: string
| LIKE
< row
: string
| LOCK
< row
: string
| LONG
< row
: string
| MAXEXTENTS
< row
: string
| MINUS
< row
: string
| MLSLABEL
< row
: string
| MODE
< row
: string
| MODIFY
< row
: string
| NESTED_TABLE_ID
< row
: string
| NOAUDIT
< row
: string
| NOCOMPRESS
< row
: string
| NOT
< row
: string
| NOWAIT
< row
: string
| NULL
< row
: string
| NUMBER
< row
: string
| OF
< row
: string
| OFFLINE
< row
: string
| ON
< row
: string
| ONLINE
< row
: string
| OPTION
< row
: string
| OR
< row
: string
| ORDER
< row
: string
| PCTFREE
< row
: string
| PRIOR
< row
: string
| PUBLIC
< row
: string
| RAW
< row
: string
| RENAME
< row
: string
| RESOURCE
< row
: string
| REVOKE
< row
: string
| ROW
< row
: string
| ROWID
< row
: string
| ROWNUM
< row
: string
| ROWS
< row
: string
| SELECT
< row
: string
| SESSION
< row
: string
| SET
< row
: string
| SHARE
< row
: string
| SIZE
< row
: string
| SMALLINT
< row
: string
| START
< row
: string
| SUCCESSFUL
< row
: string
| SYNONYM
< row
: string
| SYSDATE
< row
: string
| TABLE
< row
: string
| THEN
< row
: string
| TO
< row
: string
| TRIGGER
< row
: string
| UID
< row
: string
| UNION
< row
: string
| UNIQUE
< row
: string
| UPDATE
< row
: string
| USER
< row
: string
| VALIDATE
< row
: string
| VALUES
< row
: string
| VARCHAR
< row
: string
| VARCHAR2
< row
: string
| VIEW
< row
: string
| WHENEVER
< row
: string
| WHERE
< row
: string
| WITH
< row
: string
| [
< row
: string
| ]
< row
: string
| ^
< row
: string
| |

== TestOracleNonReservedKeywordsAsColumns default#2
> exec
| CREATE USER omni_kw IDENTIFIED BY omni_test QUOTA UNLIMITED ON users
< ok 0 0
> exec
| GRANT CREATE SESSION, CREATE TABLE, CREATE VIEW, CREATE SEQUENCE, CREATE SYNONYM, CREATE TYPE, CREATE PROCEDURE, CREATE TRIGGER TO omni_kw
< ok 0 0

== TestOracleNonReservedKeywordsAsColumns default#3
> exec
| CREATE TABLE omni_kw_0 (a NUMBER, AFTER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT AFTER FROM omni_kw_0
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_0 WHERE AFTER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_0 WHERE AFTER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_0 ORDER BY AFTER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_0 SET AFTER = NULL WHERE AFTER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_0 (AFTER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_0_i ON omni_kw_0 (AFTER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_0 WHERE AFTER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE AFTER WHEN 1 THEN 2 END FROM omni_kw_0
< ok 1 0
> exec ^
: string
| SELECT CASE k.AFTER WHEN 1 THEN 2 END FROM omni_kw_0 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + AFTER WHEN 2 THEN 3 END FROM omni_kw_0
< ok 1 0
> exec
| DROP TABLE omni_kw_0 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_1 (a NUMBER, ALWAYS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ALWAYS FROM omni_kw_1
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_1 WHERE ALWAYS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_1 WHERE ALWAYS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_1 ORDER BY ALWAYS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_1 SET ALWAYS = NULL WHERE ALWAYS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_1 (ALWAYS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_1_i ON omni_kw_1 (ALWAYS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_1 WHERE ALWAYS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ALWAYS WHEN 1 THEN 2 END FROM omni_kw_1
< ok 1 0
> exec ^
: string
| SELECT CASE k.ALWAYS WHEN 1 THEN 2 END FROM omni_kw_1 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ALWAYS WHEN 2 THEN 3 END FROM omni_kw_1
< ok 1 0
> exec
| DROP TABLE omni_kw_1 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_2 (a NUMBER, ANALYZE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ANALYZE FROM omni_kw_2
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_2 WHERE ANALYZE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_2 WHERE ANALYZE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_2 ORDER BY ANALYZE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_2 SET ANALYZE = NULL WHERE ANALYZE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_2 (ANALYZE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_2_i ON omni_kw_2 (ANALYZE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_2 WHERE ANALYZE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ANALYZE WHEN 1 THEN 2 END FROM omni_kw_2
< ok 1 0
> exec ^
: string
| SELECT CASE k.ANALYZE WHEN 1 THEN 2 END FROM omni_kw_2 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ANALYZE WHEN 2 THEN 3 END FROM omni_kw_2
< ok 1 0
> exec
| DROP TABLE omni_kw_2 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_3 (a NUMBER, ASSOCIATE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ASSOCIATE FROM omni_kw_3
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_3 WHERE ASSOCIATE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_3 WHERE ASSOCIATE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_3 ORDER BY ASSOCIATE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_3 SET ASSOCIATE = NULL WHERE ASSOCIATE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_3 (ASSOCIATE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_3_i ON omni_kw_3 (ASSOCIATE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_3 WHERE ASSOCIATE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ASSOCIATE WHEN 1 THEN 2 END FROM omni_kw_3
< ok 1 0
> exec ^
: string
| SELECT CASE k.ASSOCIATE WHEN 1 THEN 2 END FROM omni_kw_3 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ASSOCIATE WHEN 2 THEN 3 END FROM omni_kw_3
< ok 1 0
> exec
| DROP TABLE omni_kw_3 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_4 (a NUMBER, AUTOMATIC NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT AUTOMATIC FROM omni_kw_4
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_4 WHERE AUTOMATIC IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_4 WHERE AUTOMATIC NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_4 ORDER BY AUTOMATIC
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_4 SET AUTOMATIC = NULL WHERE AUTOMATIC IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_4 (AUTOMATIC) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_4_i ON omni_kw_4 (AUTOMATIC)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_4 WHERE AUTOMATIC(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE AUTOMATIC WHEN 1 THEN 2 END FROM omni_kw_4
< ok 1 0
> exec ^
: string
| SELECT CASE k.AUTOMATIC WHEN 1 THEN 2 END FROM omni_kw_4 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + AUTOMATIC WHEN 2 THEN 3 END FROM omni_kw_4
< ok 1 0
> exec
| DROP TABLE omni_kw_4 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_5 (a NUMBER, BEFORE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BEFORE FROM omni_kw_5
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_5 WHERE BEFORE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_5 WHERE BEFORE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_5 ORDER BY BEFORE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_5 SET BEFORE = NULL WHERE BEFORE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_5 (BEFORE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_5_i ON omni_kw_5 (BEFORE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_5 WHERE BEFORE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BEFORE WHEN 1 THEN 2 END FROM omni_kw_5
< ok 1 0
> exec ^
: string
| SELECT CASE k.BEFORE WHEN 1 THEN 2 END FROM omni_kw_5 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BEFORE WHEN 2 THEN 3 END FROM omni_kw_5
< ok 1 0
> exec
| DROP TABLE omni_kw_5 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_6 (a NUMBER, BEGIN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BEGIN FROM omni_kw_6
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_6 WHERE BEGIN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_6 WHERE BEGIN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_6 ORDER BY BEGIN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_6 SET BEGIN = NULL WHERE BEGIN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_6 (BEGIN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_6_i ON omni_kw_6 (BEGIN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_6 WHERE BEGIN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BEGIN WHEN 1 THEN 2 END FROM omni_kw_6
< ok 1 0
> exec ^
: string
| SELECT CASE k.BEGIN WHEN 1 THEN 2 END FROM omni_kw_6 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BEGIN WHEN 2 THEN 3 END FROM omni_kw_6
< ok 1 0
> exec
| DROP TABLE omni_kw_6 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_7 (a NUMBER, BITMAP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BITMAP FROM omni_kw_7
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_7 WHERE BITMAP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_7 WHERE BITMAP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_7 ORDER BY BITMAP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_7 SET BITMAP = NULL WHERE BITMAP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_7 (BITMAP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_7_i ON omni_kw_7 (BITMAP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_7 WHERE BITMAP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BITMAP WHEN 1 THEN 2 END FROM omni_kw_7
< ok 1 0
> exec ^
: string
| SELECT CASE k.BITMAP WHEN 1 THEN 2 END FROM omni_kw_7 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BITMAP WHEN 2 THEN 3 END FROM omni_kw_7
< ok 1 0
> exec
| DROP TABLE omni_kw_7 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_8 (a NUMBER, BLOB NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BLOB FROM omni_kw_8
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_8 WHERE BLOB IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_8 WHERE BLOB NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_8 ORDER BY BLOB
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_8 SET BLOB = NULL WHERE BLOB IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_8 (BLOB) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_8_i ON omni_kw_8 (BLOB)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_8 WHERE BLOB(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BLOB WHEN 1 THEN 2 END FROM omni_kw_8
< ok 1 0
> exec ^
: string
| SELECT CASE k.BLOB WHEN 1 THEN 2 END FROM omni_kw_8 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BLOB WHEN 2 THEN 3 END FROM omni_kw_8
< ok 1 0
> exec
| DROP TABLE omni_kw_8 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_9 (a NUMBER, BLOCK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BLOCK FROM omni_kw_9
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_9 WHERE BLOCK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_9 WHERE BLOCK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_9 ORDER BY BLOCK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_9 SET BLOCK = NULL WHERE BLOCK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_9 (BLOCK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_9_i ON omni_kw_9 (BLOCK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_9 WHERE BLOCK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BLOCK WHEN 1 THEN 2 END FROM omni_kw_9
< ok 1 0
> exec ^
: string
| SELECT CASE k.BLOCK WHEN 1 THEN 2 END FROM omni_kw_9 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BLOCK WHEN 2 THEN 3 END FROM omni_kw_9
< ok 1 0
> exec
| DROP TABLE omni_kw_9 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_10 (a NUMBER, BODY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BODY FROM omni_kw_10
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_10 WHERE BODY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_10 WHERE BODY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_10 ORDER BY BODY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_10 SET BODY = NULL WHERE BODY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_10 (BODY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_10_i ON omni_kw_10 (BODY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_10 WHERE BODY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BODY WHEN 1 THEN 2 END FROM omni_kw_10
< ok 1 0
> exec ^
: string
| SELECT CASE k.BODY WHEN 1 THEN 2 END FROM omni_kw_10 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BODY WHEN 2 THEN 3 END FROM omni_kw_10
< ok 1 0
> exec
| DROP TABLE omni_kw_10 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_11 (a NUMBER, BULK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT BULK FROM omni_kw_11
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_11 WHERE BULK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_11 WHERE BULK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_11 ORDER BY BULK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_11 SET BULK = NULL WHERE BULK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_11 (BULK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_11_i ON omni_kw_11 (BULK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_11 WHERE BULK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE BULK WHEN 1 THEN 2 END FROM omni_kw_11
< ok 1 0
> exec ^
: string
| SELECT CASE k.BULK WHEN 1 THEN 2 END FROM omni_kw_11 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + BULK WHEN 2 THEN 3 END FROM omni_kw_11
< ok 1 0
> exec
| DROP TABLE omni_kw_11 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_12 (a NUMBER, CACHE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CACHE FROM omni_kw_12
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_12 WHERE CACHE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_12 WHERE CACHE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_12 ORDER BY CACHE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_12 SET CACHE = NULL WHERE CACHE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_12 (CACHE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_12_i ON omni_kw_12 (CACHE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_12 WHERE CACHE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CACHE WHEN 1 THEN 2 END FROM omni_kw_12
< ok 1 0
> exec ^
: string
| SELECT CASE k.CACHE WHEN 1 THEN 2 END FROM omni_kw_12 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CACHE WHEN 2 THEN 3 END FROM omni_kw_12
< ok 1 0
> exec
| DROP TABLE omni_kw_12 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_13 (a NUMBER, CALL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CALL FROM omni_kw_13
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_13 WHERE CALL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_13 WHERE CALL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_13 ORDER BY CALL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_13 SET CALL = NULL WHERE CALL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_13 (CALL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_13_i ON omni_kw_13 (CALL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_13 WHERE CALL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CALL WHEN 1 THEN 2 END FROM omni_kw_13
< ok 1 0
> exec ^
: string
| SELECT CASE k.CALL WHEN 1 THEN 2 END FROM omni_kw_13 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CALL WHEN 2 THEN 3 END FROM omni_kw_13
< ok 1 0
> exec
| DROP TABLE omni_kw_13 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_14 (a NUMBER, CASCADE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CASCADE FROM omni_kw_14
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_14 WHERE CASCADE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_14 WHERE CASCADE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_14 ORDER BY CASCADE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_14 SET CASCADE = NULL WHERE CASCADE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_14 (CASCADE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_14_i ON omni_kw_14 (CASCADE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_14 WHERE CASCADE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CASCADE WHEN 1 THEN 2 END FROM omni_kw_14
< ok 1 0
> exec ^
: string
| SELECT CASE k.CASCADE WHEN 1 THEN 2 END FROM omni_kw_14 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CASCADE WHEN 2 THEN 3 END FROM omni_kw_14
< ok 1 0
> exec
| DROP TABLE omni_kw_14 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_15 (a NUMBER, CASE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CASE FROM omni_kw_15
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_15 WHERE CASE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_15 WHERE CASE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_15 ORDER BY CASE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_15 SET CASE = NULL WHERE CASE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_15 (CASE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_15_i ON omni_kw_15 (CASE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_15 WHERE CASE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CASE WHEN 1 THEN 2 END FROM omni_kw_15
< error
| 923 ORA-00923: FROM keyword not found where expected
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE k.CASE WHEN 1 THEN 2 END FROM omni_kw_15 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CASE WHEN 2 THEN 3 END FROM omni_kw_15
< error
| 923 ORA-00923: FROM keyword not found where expected
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec
| DROP TABLE omni_kw_15 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_16 (a NUMBER, CAST NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CAST FROM omni_kw_16
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_16 WHERE CAST IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_16 WHERE CAST NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_16 ORDER BY CAST
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_16 SET CAST = NULL WHERE CAST IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_16 (CAST) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_16_i ON omni_kw_16 (CAST)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_16 WHERE CAST(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CAST WHEN 1 THEN 2 END FROM omni_kw_16
< ok 1 0
> exec ^
: string
| SELECT CASE k.CAST WHEN 1 THEN 2 END FROM omni_kw_16 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CAST WHEN 2 THEN 3 END FROM omni_kw_16
< ok 1 0
> exec
| DROP TABLE omni_kw_16 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_17 (a NUMBER, CLOB NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CLOB FROM omni_kw_17
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_17 WHERE CLOB IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_17 WHERE CLOB NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_17 ORDER BY CLOB
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_17 SET CLOB = NULL WHERE CLOB IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_17 (CLOB) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_17_i ON omni_kw_17 (CLOB)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_17 WHERE CLOB(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CLOB WHEN 1 THEN 2 END FROM omni_kw_17
< ok 1 0
> exec ^
: string
| SELECT CASE k.CLOB WHEN 1 THEN 2 END FROM omni_kw_17 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CLOB WHEN 2 THEN 3 END FROM omni_kw_17
< ok 1 0
> exec
| DROP TABLE omni_kw_17 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_18 (a NUMBER, CLOSE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CLOSE FROM omni_kw_18
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_18 WHERE CLOSE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_18 WHERE CLOSE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_18 ORDER BY CLOSE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_18 SET CLOSE = NULL WHERE CLOSE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_18 (CLOSE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_18_i ON omni_kw_18 (CLOSE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_18 WHERE CLOSE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CLOSE WHEN 1 THEN 2 END FROM omni_kw_18
< ok 1 0
> exec ^
: string
| SELECT CASE k.CLOSE WHEN 1 THEN 2 END FROM omni_kw_18 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CLOSE WHEN 2 THEN 3 END FROM omni_kw_18
< ok 1 0
> exec
| DROP TABLE omni_kw_18 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_19 (a NUMBER, COLLECT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT COLLECT FROM omni_kw_19
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_19 WHERE COLLECT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_19 WHERE COLLECT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_19 ORDER BY COLLECT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_19 SET COLLECT = NULL WHERE COLLECT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_19 (COLLECT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_19_i ON omni_kw_19 (COLLECT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_19 WHERE COLLECT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE COLLECT WHEN 1 THEN 2 END FROM omni_kw_19
< ok 1 0
> exec ^
: string
| SELECT CASE k.COLLECT WHEN 1 THEN 2 END FROM omni_kw_19 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + COLLECT WHEN 2 THEN 3 END FROM omni_kw_19
< ok 1 0
> exec
| DROP TABLE omni_kw_19 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_20 (a NUMBER, COLUMNS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT COLUMNS FROM omni_kw_20
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_20 WHERE COLUMNS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_20 WHERE COLUMNS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_20 ORDER BY COLUMNS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_20 SET COLUMNS = NULL WHERE COLUMNS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_20 (COLUMNS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_20_i ON omni_kw_20 (COLUMNS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_20 WHERE COLUMNS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE COLUMNS WHEN 1 THEN 2 END FROM omni_kw_20
< ok 1 0
> exec ^
: string
| SELECT CASE k.COLUMNS WHEN 1 THEN 2 END FROM omni_kw_20 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + COLUMNS WHEN 2 THEN 3 END FROM omni_kw_20
< ok 1 0
> exec
| DROP TABLE omni_kw_20 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_21 (a NUMBER, COMMIT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT COMMIT FROM omni_kw_21
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_21 WHERE COMMIT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_21 WHERE COMMIT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_21 ORDER BY COMMIT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_21 SET COMMIT = NULL WHERE COMMIT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_21 (COMMIT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_21_i ON omni_kw_21 (COMMIT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_21 WHERE COMMIT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE COMMIT WHEN 1 THEN 2 END FROM omni_kw_21
< ok 1 0
> exec ^
: string
| SELECT CASE k.COMMIT WHEN 1 THEN 2 END FROM omni_kw_21 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + COMMIT WHEN 2 THEN 3 END FROM omni_kw_21
< ok 1 0
> exec
| DROP TABLE omni_kw_21 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_22 (a NUMBER, COMPOUND NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT COMPOUND FROM omni_kw_22
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_22 WHERE COMPOUND IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_22 WHERE COMPOUND NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_22 ORDER BY COMPOUND
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_22 SET COMPOUND = NULL WHERE COMPOUND IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_22 (COMPOUND) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_22_i ON omni_kw_22 (COMPOUND)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_22 WHERE COMPOUND(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE COMPOUND WHEN 1 THEN 2 END FROM omni_kw_22
< ok 1 0
> exec ^
: string
| SELECT CASE k.COMPOUND WHEN 1 THEN 2 END FROM omni_kw_22 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + COMPOUND WHEN 2 THEN 3 END FROM omni_kw_22
< ok 1 0
> exec
| DROP TABLE omni_kw_22 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_23 (a NUMBER, CONNECT_BY_ROOT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONNECT_BY_ROOT FROM omni_kw_23
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT a FROM omni_kw_23 WHERE CONNECT_BY_ROOT IS NULL
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT a FROM omni_kw_23 WHERE CONNECT_BY_ROOT NOT IN (1, 2)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT a FROM omni_kw_23 ORDER BY CONNECT_BY_ROOT
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| UPDATE omni_kw_23 SET CONNECT_BY_ROOT = NULL WHERE CONNECT_BY_ROOT IS NULL
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| INSERT INTO omni_kw_23 (CONNECT_BY_ROOT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_23_i ON omni_kw_23 (CONNECT_BY_ROOT)
< error
| 976 ORA-00976: Specified pseudocolumn or operator not allowed here.
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_23 WHERE CONNECT_BY_ROOT(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE CONNECT_BY_ROOT WHEN 1 THEN 2 END FROM omni_kw_23
< error
| 923 ORA-00923: FROM keyword not found where expected
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE k.CONNECT_BY_ROOT WHEN 1 THEN 2 END FROM omni_kw_23 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONNECT_BY_ROOT WHEN 2 THEN 3 END FROM omni_kw_23
< error
| 923 ORA-00923: FROM keyword not found where expected
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec
| DROP TABLE omni_kw_23 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_24 (a NUMBER, CONSTANT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONSTANT FROM omni_kw_24
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_24 WHERE CONSTANT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_24 WHERE CONSTANT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_24 ORDER BY CONSTANT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_24 SET CONSTANT = NULL WHERE CONSTANT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_24 (CONSTANT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_24_i ON omni_kw_24 (CONSTANT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_24 WHERE CONSTANT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONSTANT WHEN 1 THEN 2 END FROM omni_kw_24
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONSTANT WHEN 1 THEN 2 END FROM omni_kw_24 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONSTANT WHEN 2 THEN 3 END FROM omni_kw_24
< ok 1 0
> exec
| DROP TABLE omni_kw_24 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_25 (a NUMBER, CONSTRAINT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONSTRAINT FROM omni_kw_25
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_25 WHERE CONSTRAINT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_25 WHERE CONSTRAINT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_25 ORDER BY CONSTRAINT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_25 SET CONSTRAINT = NULL WHERE CONSTRAINT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_25 (CONSTRAINT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_25_i ON omni_kw_25 (CONSTRAINT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_25 WHERE CONSTRAINT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONSTRAINT WHEN 1 THEN 2 END FROM omni_kw_25
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONSTRAINT WHEN 1 THEN 2 END FROM omni_kw_25 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONSTRAINT WHEN 2 THEN 3 END FROM omni_kw_25
< ok 1 0
> exec
| DROP TABLE omni_kw_25 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_26 (a NUMBER, CONSTRAINTS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONSTRAINTS FROM omni_kw_26
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_26 WHERE CONSTRAINTS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_26 WHERE CONSTRAINTS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_26 ORDER BY CONSTRAINTS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_26 SET CONSTRAINTS = NULL WHERE CONSTRAINTS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_26 (CONSTRAINTS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_26_i ON omni_kw_26 (CONSTRAINTS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_26 WHERE CONSTRAINTS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONSTRAINTS WHEN 1 THEN 2 END FROM omni_kw_26
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONSTRAINTS WHEN 1 THEN 2 END FROM omni_kw_26 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONSTRAINTS WHEN 2 THEN 3 END FROM omni_kw_26
< ok 1 0
> exec
| DROP TABLE omni_kw_26 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_27 (a NUMBER, CONTENT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONTENT FROM omni_kw_27
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_27 WHERE CONTENT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_27 WHERE CONTENT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_27 ORDER BY CONTENT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_27 SET CONTENT = NULL WHERE CONTENT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_27 (CONTENT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_27_i ON omni_kw_27 (CONTENT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_27 WHERE CONTENT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONTENT WHEN 1 THEN 2 END FROM omni_kw_27
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONTENT WHEN 1 THEN 2 END FROM omni_kw_27 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONTENT WHEN 2 THEN 3 END FROM omni_kw_27
< ok 1 0
> exec
| DROP TABLE omni_kw_27 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_28 (a NUMBER, CONTEXT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONTEXT FROM omni_kw_28
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_28 WHERE CONTEXT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_28 WHERE CONTEXT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_28 ORDER BY CONTEXT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_28 SET CONTEXT = NULL WHERE CONTEXT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_28 (CONTEXT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_28_i ON omni_kw_28 (CONTEXT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_28 WHERE CONTEXT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONTEXT WHEN 1 THEN 2 END FROM omni_kw_28
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONTEXT WHEN 1 THEN 2 END FROM omni_kw_28 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONTEXT WHEN 2 THEN 3 END FROM omni_kw_28
< ok 1 0
> exec
| DROP TABLE omni_kw_28 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_29 (a NUMBER, CONTINUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CONTINUE FROM omni_kw_29
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_29 WHERE CONTINUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_29 WHERE CONTINUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_29 ORDER BY CONTINUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_29 SET CONTINUE = NULL WHERE CONTINUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_29 (CONTINUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_29_i ON omni_kw_29 (CONTINUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_29 WHERE CONTINUE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CONTINUE WHEN 1 THEN 2 END FROM omni_kw_29
< ok 1 0
> exec ^
: string
| SELECT CASE k.CONTINUE WHEN 1 THEN 2 END FROM omni_kw_29 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CONTINUE WHEN 2 THEN 3 END FROM omni_kw_29
< ok 1 0
> exec
| DROP TABLE omni_kw_29 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_30 (a NUMBER, CROSS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CROSS FROM omni_kw_30
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_30 WHERE CROSS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_30 WHERE CROSS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_30 ORDER BY CROSS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_30 SET CROSS = NULL WHERE CROSS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_30 (CROSS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_30_i ON omni_kw_30 (CROSS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_30 WHERE CROSS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CROSS WHEN 1 THEN 2 END FROM omni_kw_30
< ok 1 0
> exec ^
: string
| SELECT CASE k.CROSS WHEN 1 THEN 2 END FROM omni_kw_30 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CROSS WHEN 2 THEN 3 END FROM omni_kw_30
< ok 1 0
> exec
| DROP TABLE omni_kw_30 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_31 (a NUMBER, CUBE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CUBE FROM omni_kw_31
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_31 WHERE CUBE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_31 WHERE CUBE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_31 ORDER BY CUBE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_31 SET CUBE = NULL WHERE CUBE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_31 (CUBE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_31_i ON omni_kw_31 (CUBE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_31 WHERE CUBE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CUBE WHEN 1 THEN 2 END FROM omni_kw_31
< ok 1 0
> exec ^
: string
| SELECT CASE k.CUBE WHEN 1 THEN 2 END FROM omni_kw_31 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CUBE WHEN 2 THEN 3 END FROM omni_kw_31
< ok 1 0
> exec
| DROP TABLE omni_kw_31 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_32 (a NUMBER, CURSOR NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CURSOR FROM omni_kw_32
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_32 WHERE CURSOR IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_32 WHERE CURSOR NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_32 ORDER BY CURSOR
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_32 SET CURSOR = NULL WHERE CURSOR IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_32 (CURSOR) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_32_i ON omni_kw_32 (CURSOR)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_32 WHERE CURSOR(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CURSOR WHEN 1 THEN 2 END FROM omni_kw_32
< ok 1 0
> exec ^
: string
| SELECT CASE k.CURSOR WHEN 1 THEN 2 END FROM omni_kw_32 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CURSOR WHEN 2 THEN 3 END FROM omni_kw_32
< ok 1 0
> exec
| DROP TABLE omni_kw_32 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_33 (a NUMBER, CYCLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT CYCLE FROM omni_kw_33
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_33 WHERE CYCLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_33 WHERE CYCLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_33 ORDER BY CYCLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_33 SET CYCLE = NULL WHERE CYCLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_33 (CYCLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_33_i ON omni_kw_33 (CYCLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_33 WHERE CYCLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE CYCLE WHEN 1 THEN 2 END FROM omni_kw_33
< ok 1 0
> exec ^
: string
| SELECT CASE k.CYCLE WHEN 1 THEN 2 END FROM omni_kw_33 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + CYCLE WHEN 2 THEN 3 END FROM omni_kw_33
< ok 1 0
> exec
| DROP TABLE omni_kw_33 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_34 (a NUMBER, DATABASE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DATABASE FROM omni_kw_34
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_34 WHERE DATABASE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_34 WHERE DATABASE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_34 ORDER BY DATABASE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_34 SET DATABASE = NULL WHERE DATABASE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_34 (DATABASE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_34_i ON omni_kw_34 (DATABASE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_34 WHERE DATABASE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DATABASE WHEN 1 THEN 2 END FROM omni_kw_34
< ok 1 0
> exec ^
: string
| SELECT CASE k.DATABASE WHEN 1 THEN 2 END FROM omni_kw_34 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DATABASE WHEN 2 THEN 3 END FROM omni_kw_34
< ok 1 0
> exec
| DROP TABLE omni_kw_34 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_35 (a NUMBER, DECLARE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DECLARE FROM omni_kw_35
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_35 WHERE DECLARE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_35 WHERE DECLARE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_35 ORDER BY DECLARE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_35 SET DECLARE = NULL WHERE DECLARE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_35 (DECLARE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_35_i ON omni_kw_35 (DECLARE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_35 WHERE DECLARE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DECLARE WHEN 1 THEN 2 END FROM omni_kw_35
< ok 1 0
> exec ^
: string
| SELECT CASE k.DECLARE WHEN 1 THEN 2 END FROM omni_kw_35 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DECLARE WHEN 2 THEN 3 END FROM omni_kw_35
< ok 1 0
> exec
| DROP TABLE omni_kw_35 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_36 (a NUMBER, DECODE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DECODE FROM omni_kw_36
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_36 WHERE DECODE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_36 WHERE DECODE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_36 ORDER BY DECODE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_36 SET DECODE = NULL WHERE DECODE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_36 (DECODE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_36_i ON omni_kw_36 (DECODE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_36 WHERE DECODE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DECODE WHEN 1 THEN 2 END FROM omni_kw_36
< ok 1 0
> exec ^
: string
| SELECT CASE k.DECODE WHEN 1 THEN 2 END FROM omni_kw_36 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DECODE WHEN 2 THEN 3 END FROM omni_kw_36
< ok 1 0
> exec
| DROP TABLE omni_kw_36 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_37 (a NUMBER, DECREMENT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DECREMENT FROM omni_kw_37
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_37 WHERE DECREMENT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_37 WHERE DECREMENT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_37 ORDER BY DECREMENT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_37 SET DECREMENT = NULL WHERE DECREMENT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_37 (DECREMENT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_37_i ON omni_kw_37 (DECREMENT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_37 WHERE DECREMENT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DECREMENT WHEN 1 THEN 2 END FROM omni_kw_37
< ok 1 0
> exec ^
: string
| SELECT CASE k.DECREMENT WHEN 1 THEN 2 END FROM omni_kw_37 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DECREMENT WHEN 2 THEN 3 END FROM omni_kw_37
< ok 1 0
> exec
| DROP TABLE omni_kw_37 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_38 (a NUMBER, DEFERRABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DEFERRABLE FROM omni_kw_38
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_38 WHERE DEFERRABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_38 WHERE DEFERRABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_38 ORDER BY DEFERRABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_38 SET DEFERRABLE = NULL WHERE DEFERRABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_38 (DEFERRABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_38_i ON omni_kw_38 (DEFERRABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_38 WHERE DEFERRABLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DEFERRABLE WHEN 1 THEN 2 END FROM omni_kw_38
< ok 1 0
> exec ^
: string
| SELECT CASE k.DEFERRABLE WHEN 1 THEN 2 END FROM omni_kw_38 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DEFERRABLE WHEN 2 THEN 3 END FROM omni_kw_38
< ok 1 0
> exec
| DROP TABLE omni_kw_38 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_39 (a NUMBER, DEFERRED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DEFERRED FROM omni_kw_39
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_39 WHERE DEFERRED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_39 WHERE DEFERRED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_39 ORDER BY DEFERRED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_39 SET DEFERRED = NULL WHERE DEFERRED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_39 (DEFERRED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_39_i ON omni_kw_39 (DEFERRED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_39 WHERE DEFERRED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DEFERRED WHEN 1 THEN 2 END FROM omni_kw_39
< ok 1 0
> exec ^
: string
| SELECT CASE k.DEFERRED WHEN 1 THEN 2 END FROM omni_kw_39 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DEFERRED WHEN 2 THEN 3 END FROM omni_kw_39
< ok 1 0
> exec
| DROP TABLE omni_kw_39 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_40 (a NUMBER, DENSE_RANK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DENSE_RANK FROM omni_kw_40
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_40 WHERE DENSE_RANK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_40 WHERE DENSE_RANK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_40 ORDER BY DENSE_RANK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_40 SET DENSE_RANK = NULL WHERE DENSE_RANK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_40 (DENSE_RANK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_40_i ON omni_kw_40 (DENSE_RANK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_40 WHERE DENSE_RANK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DENSE_RANK WHEN 1 THEN 2 END FROM omni_kw_40
< ok 1 0
> exec ^
: string
| SELECT CASE k.DENSE_RANK WHEN 1 THEN 2 END FROM omni_kw_40 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DENSE_RANK WHEN 2 THEN 3 END FROM omni_kw_40
< ok 1 0
> exec
| DROP TABLE omni_kw_40 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_41 (a NUMBER, DETERMINISTIC NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DETERMINISTIC FROM omni_kw_41
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_41 WHERE DETERMINISTIC IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_41 WHERE DETERMINISTIC NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_41 ORDER BY DETERMINISTIC
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_41 SET DETERMINISTIC = NULL WHERE DETERMINISTIC IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_41 (DETERMINISTIC) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_41_i ON omni_kw_41 (DETERMINISTIC)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_41 WHERE DETERMINISTIC(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DETERMINISTIC WHEN 1 THEN 2 END FROM omni_kw_41
< ok 1 0
> exec ^
: string
| SELECT CASE k.DETERMINISTIC WHEN 1 THEN 2 END FROM omni_kw_41 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DETERMINISTIC WHEN 2 THEN 3 END FROM omni_kw_41
< ok 1 0
> exec
| DROP TABLE omni_kw_41 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_42 (a NUMBER, DIMENSION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DIMENSION FROM omni_kw_42
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_42 WHERE DIMENSION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_42 WHERE DIMENSION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_42 ORDER BY DIMENSION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_42 SET DIMENSION = NULL WHERE DIMENSION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_42 (DIMENSION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_42_i ON omni_kw_42 (DIMENSION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_42 WHERE DIMENSION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DIMENSION WHEN 1 THEN 2 END FROM omni_kw_42
< ok 1 0
> exec ^
: string
| SELECT CASE k.DIMENSION WHEN 1 THEN 2 END FROM omni_kw_42 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DIMENSION WHEN 2 THEN 3 END FROM omni_kw_42
< ok 1 0
> exec
| DROP TABLE omni_kw_42 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_43 (a NUMBER, DIRECTORY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DIRECTORY FROM omni_kw_43
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_43 WHERE DIRECTORY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_43 WHERE DIRECTORY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_43 ORDER BY DIRECTORY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_43 SET DIRECTORY = NULL WHERE DIRECTORY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_43 (DIRECTORY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_43_i ON omni_kw_43 (DIRECTORY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_43 WHERE DIRECTORY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DIRECTORY WHEN 1 THEN 2 END FROM omni_kw_43
< ok 1 0
> exec ^
: string
| SELECT CASE k.DIRECTORY WHEN 1 THEN 2 END FROM omni_kw_43 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DIRECTORY WHEN 2 THEN 3 END FROM omni_kw_43
< ok 1 0
> exec
| DROP TABLE omni_kw_43 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_44 (a NUMBER, DISABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DISABLE FROM omni_kw_44
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_44 WHERE DISABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_44 WHERE DISABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_44 ORDER BY DISABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_44 SET DISABLE = NULL WHERE DISABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_44 (DISABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_44_i ON omni_kw_44 (DISABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_44 WHERE DISABLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DISABLE WHEN 1 THEN 2 END FROM omni_kw_44
< ok 1 0
> exec ^
: string
| SELECT CASE k.DISABLE WHEN 1 THEN 2 END FROM omni_kw_44 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DISABLE WHEN 2 THEN 3 END FROM omni_kw_44
< ok 1 0
> exec
| DROP TABLE omni_kw_44 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_45 (a NUMBER, DISASSOCIATE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT DISASSOCIATE FROM omni_kw_45
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_45 WHERE DISASSOCIATE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_45 WHERE DISASSOCIATE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_45 ORDER BY DISASSOCIATE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_45 SET DISASSOCIATE = NULL WHERE DISASSOCIATE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_45 (DISASSOCIATE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_45_i ON omni_kw_45 (DISASSOCIATE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_45 WHERE DISASSOCIATE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE DISASSOCIATE WHEN 1 THEN 2 END FROM omni_kw_45
< ok 1 0
> exec ^
: string
| SELECT CASE k.DISASSOCIATE WHEN 1 THEN 2 END FROM omni_kw_45 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + DISASSOCIATE WHEN 2 THEN 3 END FROM omni_kw_45
< ok 1 0
> exec
| DROP TABLE omni_kw_45 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_46 (a NUMBER, EACH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT EACH FROM omni_kw_46
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_46 WHERE EACH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_46 WHERE EACH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_46 ORDER BY EACH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_46 SET EACH = NULL WHERE EACH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_46 (EACH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_46_i ON omni_kw_46 (EACH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_46 WHERE EACH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE EACH WHEN 1 THEN 2 END FROM omni_kw_46
< ok 1 0
> exec ^
: string
| SELECT CASE k.EACH WHEN 1 THEN 2 END FROM omni_kw_46 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + EACH WHEN 2 THEN 3 END FROM omni_kw_46
< ok 1 0
> exec
| DROP TABLE omni_kw_46 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_47 (a NUMBER, ELSIF NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ELSIF FROM omni_kw_47
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_47 WHERE ELSIF IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_47 WHERE ELSIF NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_47 ORDER BY ELSIF
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_47 SET ELSIF = NULL WHERE ELSIF IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_47 (ELSIF) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_47_i ON omni_kw_47 (ELSIF)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_47 WHERE ELSIF(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ELSIF WHEN 1 THEN 2 END FROM omni_kw_47
< ok 1 0
> exec ^
: string
| SELECT CASE k.ELSIF WHEN 1 THEN 2 END FROM omni_kw_47 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ELSIF WHEN 2 THEN 3 END FROM omni_kw_47
< ok 1 0
> exec
| DROP TABLE omni_kw_47 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_48 (a NUMBER, ENABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ENABLE FROM omni_kw_48
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_48 WHERE ENABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_48 WHERE ENABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_48 ORDER BY ENABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_48 SET ENABLE = NULL WHERE ENABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_48 (ENABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_48_i ON omni_kw_48 (ENABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_48 WHERE ENABLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ENABLE WHEN 1 THEN 2 END FROM omni_kw_48
< ok 1 0
> exec ^
: string
| SELECT CASE k.ENABLE WHEN 1 THEN 2 END FROM omni_kw_48 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ENABLE WHEN 2 THEN 3 END FROM omni_kw_48
< ok 1 0
> exec
| DROP TABLE omni_kw_48 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_49 (a NUMBER, END NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT END FROM omni_kw_49
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_49 WHERE END IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_49 WHERE END NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_49 ORDER BY END
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_49 SET END = NULL WHERE END IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_49 (END) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_49_i ON omni_kw_49 (END)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_49 WHERE END(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE END WHEN 1 THEN 2 END FROM omni_kw_49
< ok 1 0
> exec ^
: string
| SELECT CASE k.END WHEN 1 THEN 2 END FROM omni_kw_49 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + END WHEN 2 THEN 3 END FROM omni_kw_49
< ok 1 0
> exec
| DROP TABLE omni_kw_49 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_50 (a NUMBER, ERRORS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ERRORS FROM omni_kw_50
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_50 WHERE ERRORS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_50 WHERE ERRORS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_50 ORDER BY ERRORS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_50 SET ERRORS = NULL WHERE ERRORS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_50 (ERRORS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_50_i ON omni_kw_50 (ERRORS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_50 WHERE ERRORS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ERRORS WHEN 1 THEN 2 END FROM omni_kw_50
< ok 1 0
> exec ^
: string
| SELECT CASE k.ERRORS WHEN 1 THEN 2 END FROM omni_kw_50 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ERRORS WHEN 2 THEN 3 END FROM omni_kw_50
< ok 1 0
> exec
| DROP TABLE omni_kw_50 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_51 (a NUMBER, ESCAPE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ESCAPE FROM omni_kw_51
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_51 WHERE ESCAPE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_51 WHERE ESCAPE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_51 ORDER BY ESCAPE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_51 SET ESCAPE = NULL WHERE ESCAPE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_51 (ESCAPE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_51_i ON omni_kw_51 (ESCAPE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_51 WHERE ESCAPE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ESCAPE WHEN 1 THEN 2 END FROM omni_kw_51
< ok 1 0
> exec ^
: string
| SELECT CASE k.ESCAPE WHEN 1 THEN 2 END FROM omni_kw_51 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ESCAPE WHEN 2 THEN 3 END FROM omni_kw_51
< ok 1 0
> exec
| DROP TABLE omni_kw_51 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_52 (a NUMBER, EXCEPTION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT EXCEPTION FROM omni_kw_52
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_52 WHERE EXCEPTION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_52 WHERE EXCEPTION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_52 ORDER BY EXCEPTION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_52 SET EXCEPTION = NULL WHERE EXCEPTION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_52 (EXCEPTION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_52_i ON omni_kw_52 (EXCEPTION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_52 WHERE EXCEPTION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE EXCEPTION WHEN 1 THEN 2 END FROM omni_kw_52
< ok 1 0
> exec ^
: string
| SELECT CASE k.EXCEPTION WHEN 1 THEN 2 END FROM omni_kw_52 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + EXCEPTION WHEN 2 THEN 3 END FROM omni_kw_52
< ok 1 0
> exec
| DROP TABLE omni_kw_52 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_53 (a NUMBER, EXECUTE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT EXECUTE FROM omni_kw_53
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_53 WHERE EXECUTE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_53 WHERE EXECUTE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_53 ORDER BY EXECUTE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_53 SET EXECUTE = NULL WHERE EXECUTE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_53 (EXECUTE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_53_i ON omni_kw_53 (EXECUTE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_53 WHERE EXECUTE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE EXECUTE WHEN 1 THEN 2 END FROM omni_kw_53
< ok 1 0
> exec ^
: string
| SELECT CASE k.EXECUTE WHEN 1 THEN 2 END FROM omni_kw_53 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + EXECUTE WHEN 2 THEN 3 END FROM omni_kw_53
< ok 1 0
> exec
| DROP TABLE omni_kw_53 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_54 (a NUMBER, EXIT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT EXIT FROM omni_kw_54
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_54 WHERE EXIT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_54 WHERE EXIT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_54 ORDER BY EXIT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_54 SET EXIT = NULL WHERE EXIT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_54 (EXIT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_54_i ON omni_kw_54 (EXIT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_54 WHERE EXIT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE EXIT WHEN 1 THEN 2 END FROM omni_kw_54
< ok 1 0
> exec ^
: string
| SELECT CASE k.EXIT WHEN 1 THEN 2 END FROM omni_kw_54 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + EXIT WHEN 2 THEN 3 END FROM omni_kw_54
< ok 1 0
> exec
| DROP TABLE omni_kw_54 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_55 (a NUMBER, EXPLAIN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT EXPLAIN FROM omni_kw_55
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_55 WHERE EXPLAIN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_55 WHERE EXPLAIN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_55 ORDER BY EXPLAIN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_55 SET EXPLAIN = NULL WHERE EXPLAIN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_55 (EXPLAIN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_55_i ON omni_kw_55 (EXPLAIN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_55 WHERE EXPLAIN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE EXPLAIN WHEN 1 THEN 2 END FROM omni_kw_55
< ok 1 0
> exec ^
: string
| SELECT CASE k.EXPLAIN WHEN 1 THEN 2 END FROM omni_kw_55 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + EXPLAIN WHEN 2 THEN 3 END FROM omni_kw_55
< ok 1 0
> exec
| DROP TABLE omni_kw_55 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_56 (a NUMBER, FETCH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FETCH FROM omni_kw_56
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_56 WHERE FETCH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_56 WHERE FETCH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_56 ORDER BY FETCH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_56 SET FETCH = NULL WHERE FETCH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_56 (FETCH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_56_i ON omni_kw_56 (FETCH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_56 WHERE FETCH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FETCH WHEN 1 THEN 2 END FROM omni_kw_56
< ok 1 0
> exec ^
: string
| SELECT CASE k.FETCH WHEN 1 THEN 2 END FROM omni_kw_56 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FETCH WHEN 2 THEN 3 END FROM omni_kw_56
< ok 1 0
> exec
| DROP TABLE omni_kw_56 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_57 (a NUMBER, FIRST NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FIRST FROM omni_kw_57
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_57 WHERE FIRST IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_57 WHERE FIRST NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_57 ORDER BY FIRST
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_57 SET FIRST = NULL WHERE FIRST IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_57 (FIRST) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_57_i ON omni_kw_57 (FIRST)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_57 WHERE FIRST(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FIRST WHEN 1 THEN 2 END FROM omni_kw_57
< ok 1 0
> exec ^
: string
| SELECT CASE k.FIRST WHEN 1 THEN 2 END FROM omni_kw_57 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FIRST WHEN 2 THEN 3 END FROM omni_kw_57
< ok 1 0
> exec
| DROP TABLE omni_kw_57 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_58 (a NUMBER, FLASHBACK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FLASHBACK FROM omni_kw_58
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_58 WHERE FLASHBACK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_58 WHERE FLASHBACK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_58 ORDER BY FLASHBACK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_58 SET FLASHBACK = NULL WHERE FLASHBACK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_58 (FLASHBACK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_58_i ON omni_kw_58 (FLASHBACK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_58 WHERE FLASHBACK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FLASHBACK WHEN 1 THEN 2 END FROM omni_kw_58
< ok 1 0
> exec ^
: string
| SELECT CASE k.FLASHBACK WHEN 1 THEN 2 END FROM omni_kw_58 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FLASHBACK WHEN 2 THEN 3 END FROM omni_kw_58
< ok 1 0
> exec
| DROP TABLE omni_kw_58 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_59 (a NUMBER, FOLLOWING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FOLLOWING FROM omni_kw_59
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_59 WHERE FOLLOWING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_59 WHERE FOLLOWING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_59 ORDER BY FOLLOWING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_59 SET FOLLOWING = NULL WHERE FOLLOWING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_59 (FOLLOWING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_59_i ON omni_kw_59 (FOLLOWING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_59 WHERE FOLLOWING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FOLLOWING WHEN 1 THEN 2 END FROM omni_kw_59
< ok 1 0
> exec ^
: string
| SELECT CASE k.FOLLOWING WHEN 1 THEN 2 END FROM omni_kw_59 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FOLLOWING WHEN 2 THEN 3 END FROM omni_kw_59
< ok 1 0
> exec
| DROP TABLE omni_kw_59 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_60 (a NUMBER, FORALL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FORALL FROM omni_kw_60
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_60 WHERE FORALL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_60 WHERE FORALL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_60 ORDER BY FORALL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_60 SET FORALL = NULL WHERE FORALL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_60 (FORALL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_60_i ON omni_kw_60 (FORALL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_60 WHERE FORALL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FORALL WHEN 1 THEN 2 END FROM omni_kw_60
< ok 1 0
> exec ^
: string
| SELECT CASE k.FORALL WHEN 1 THEN 2 END FROM omni_kw_60 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FORALL WHEN 2 THEN 3 END FROM omni_kw_60
< ok 1 0
> exec
| DROP TABLE omni_kw_60 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_61 (a NUMBER, FORCE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FORCE FROM omni_kw_61
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_61 WHERE FORCE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_61 WHERE FORCE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_61 ORDER BY FORCE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_61 SET FORCE = NULL WHERE FORCE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_61 (FORCE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_61_i ON omni_kw_61 (FORCE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_61 WHERE FORCE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FORCE WHEN 1 THEN 2 END FROM omni_kw_61
< ok 1 0
> exec ^
: string
| SELECT CASE k.FORCE WHEN 1 THEN 2 END FROM omni_kw_61 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FORCE WHEN 2 THEN 3 END FROM omni_kw_61
< ok 1 0
> exec
| DROP TABLE omni_kw_61 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_62 (a NUMBER, FOREIGN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FOREIGN FROM omni_kw_62
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_62 WHERE FOREIGN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_62 WHERE FOREIGN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_62 ORDER BY FOREIGN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_62 SET FOREIGN = NULL WHERE FOREIGN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_62 (FOREIGN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_62_i ON omni_kw_62 (FOREIGN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_62 WHERE FOREIGN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FOREIGN WHEN 1 THEN 2 END FROM omni_kw_62
< ok 1 0
> exec ^
: string
| SELECT CASE k.FOREIGN WHEN 1 THEN 2 END FROM omni_kw_62 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FOREIGN WHEN 2 THEN 3 END FROM omni_kw_62
< ok 1 0
> exec
| DROP TABLE omni_kw_62 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_63 (a NUMBER, FORMAT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FORMAT FROM omni_kw_63
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_63 WHERE FORMAT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_63 WHERE FORMAT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_63 ORDER BY FORMAT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_63 SET FORMAT = NULL WHERE FORMAT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_63 (FORMAT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_63_i ON omni_kw_63 (FORMAT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_63 WHERE FORMAT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FORMAT WHEN 1 THEN 2 END FROM omni_kw_63
< ok 1 0
> exec ^
: string
| SELECT CASE k.FORMAT WHEN 1 THEN 2 END FROM omni_kw_63 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FORMAT WHEN 2 THEN 3 END FROM omni_kw_63
< ok 1 0
> exec
| DROP TABLE omni_kw_63 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_64 (a NUMBER, FULL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FULL FROM omni_kw_64
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_64 WHERE FULL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_64 WHERE FULL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_64 ORDER BY FULL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_64 SET FULL = NULL WHERE FULL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_64 (FULL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_64_i ON omni_kw_64 (FULL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_64 WHERE FULL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FULL WHEN 1 THEN 2 END FROM omni_kw_64
< ok 1 0
> exec ^
: string
| SELECT CASE k.FULL WHEN 1 THEN 2 END FROM omni_kw_64 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FULL WHEN 2 THEN 3 END FROM omni_kw_64
< ok 1 0
> exec
| DROP TABLE omni_kw_64 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_65 (a NUMBER, FUNCTION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT FUNCTION FROM omni_kw_65
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_65 WHERE FUNCTION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_65 WHERE FUNCTION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_65 ORDER BY FUNCTION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_65 SET FUNCTION = NULL WHERE FUNCTION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_65 (FUNCTION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_65_i ON omni_kw_65 (FUNCTION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_65 WHERE FUNCTION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE FUNCTION WHEN 1 THEN 2 END FROM omni_kw_65
< ok 1 0
> exec ^
: string
| SELECT CASE k.FUNCTION WHEN 1 THEN 2 END FROM omni_kw_65 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + FUNCTION WHEN 2 THEN 3 END FROM omni_kw_65
< ok 1 0
> exec
| DROP TABLE omni_kw_65 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_66 (a NUMBER, GENERATED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT GENERATED FROM omni_kw_66
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_66 WHERE GENERATED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_66 WHERE GENERATED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_66 ORDER BY GENERATED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_66 SET GENERATED = NULL WHERE GENERATED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_66 (GENERATED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_66_i ON omni_kw_66 (GENERATED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_66 WHERE GENERATED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE GENERATED WHEN 1 THEN 2 END FROM omni_kw_66
< ok 1 0
> exec ^
: string
| SELECT CASE k.GENERATED WHEN 1 THEN 2 END FROM omni_kw_66 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + GENERATED WHEN 2 THEN 3 END FROM omni_kw_66
< ok 1 0
> exec
| DROP TABLE omni_kw_66 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_67 (a NUMBER, GLOBAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT GLOBAL FROM omni_kw_67
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_67 WHERE GLOBAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_67 WHERE GLOBAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_67 ORDER BY GLOBAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_67 SET GLOBAL = NULL WHERE GLOBAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_67 (GLOBAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_67_i ON omni_kw_67 (GLOBAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_67 WHERE GLOBAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE GLOBAL WHEN 1 THEN 2 END FROM omni_kw_67
< ok 1 0
> exec ^
: string
| SELECT CASE k.GLOBAL WHEN 1 THEN 2 END FROM omni_kw_67 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + GLOBAL WHEN 2 THEN 3 END FROM omni_kw_67
< ok 1 0
> exec
| DROP TABLE omni_kw_67 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_68 (a NUMBER, GOTO NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT GOTO FROM omni_kw_68
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_68 WHERE GOTO IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_68 WHERE GOTO NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_68 ORDER BY GOTO
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_68 SET GOTO = NULL WHERE GOTO IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_68 (GOTO) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_68_i ON omni_kw_68 (GOTO)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_68 WHERE GOTO(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE GOTO WHEN 1 THEN 2 END FROM omni_kw_68
< ok 1 0
> exec ^
: string
| SELECT CASE k.GOTO WHEN 1 THEN 2 END FROM omni_kw_68 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + GOTO WHEN 2 THEN 3 END FROM omni_kw_68
< ok 1 0
> exec
| DROP TABLE omni_kw_68 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_69 (a NUMBER, GROUPING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT GROUPING FROM omni_kw_69
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_69 WHERE GROUPING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_69 WHERE GROUPING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_69 ORDER BY GROUPING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_69 SET GROUPING = NULL WHERE GROUPING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_69 (GROUPING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_69_i ON omni_kw_69 (GROUPING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_69 WHERE GROUPING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE GROUPING WHEN 1 THEN 2 END FROM omni_kw_69
< ok 1 0
> exec ^
: string
| SELECT CASE k.GROUPING WHEN 1 THEN 2 END FROM omni_kw_69 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + GROUPING WHEN 2 THEN 3 END FROM omni_kw_69
< ok 1 0
> exec
| DROP TABLE omni_kw_69 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_70 (a NUMBER, GROUPS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT GROUPS FROM omni_kw_70
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_70 WHERE GROUPS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_70 WHERE GROUPS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_70 ORDER BY GROUPS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_70 SET GROUPS = NULL WHERE GROUPS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_70 (GROUPS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_70_i ON omni_kw_70 (GROUPS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_70 WHERE GROUPS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE GROUPS WHEN 1 THEN 2 END FROM omni_kw_70
< ok 1 0
> exec ^
: string
| SELECT CASE k.GROUPS WHEN 1 THEN 2 END FROM omni_kw_70 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + GROUPS WHEN 2 THEN 3 END FROM omni_kw_70
< ok 1 0
> exec
| DROP TABLE omni_kw_70 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_71 (a NUMBER, HASH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT HASH FROM omni_kw_71
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_71 WHERE HASH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_71 WHERE HASH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_71 ORDER BY HASH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_71 SET HASH = NULL WHERE HASH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_71 (HASH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_71_i ON omni_kw_71 (HASH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_71 WHERE HASH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE HASH WHEN 1 THEN 2 END FROM omni_kw_71
< ok 1 0
> exec ^
: string
| SELECT CASE k.HASH WHEN 1 THEN 2 END FROM omni_kw_71 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + HASH WHEN 2 THEN 3 END FROM omni_kw_71
< ok 1 0
> exec
| DROP TABLE omni_kw_71 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_72 (a NUMBER, IDENTITY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT IDENTITY FROM omni_kw_72
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_72 WHERE IDENTITY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_72 WHERE IDENTITY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_72 ORDER BY IDENTITY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_72 SET IDENTITY = NULL WHERE IDENTITY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_72 (IDENTITY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_72_i ON omni_kw_72 (IDENTITY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_72 WHERE IDENTITY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE IDENTITY WHEN 1 THEN 2 END FROM omni_kw_72
< ok 1 0
> exec ^
: string
| SELECT CASE k.IDENTITY WHEN 1 THEN 2 END FROM omni_kw_72 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + IDENTITY WHEN 2 THEN 3 END FROM omni_kw_72
< ok 1 0
> exec
| DROP TABLE omni_kw_72 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_73 (a NUMBER, IF NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT IF FROM omni_kw_73
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_73 WHERE IF IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_73 WHERE IF NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_73 ORDER BY IF
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_73 SET IF = NULL WHERE IF IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_73 (IF) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_73_i ON omni_kw_73 (IF)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_73 WHERE IF(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE IF WHEN 1 THEN 2 END FROM omni_kw_73
< ok 1 0
> exec ^
: string
| SELECT CASE k.IF WHEN 1 THEN 2 END FROM omni_kw_73 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + IF WHEN 2 THEN 3 END FROM omni_kw_73
< ok 1 0
> exec
| DROP TABLE omni_kw_73 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_74 (a NUMBER, INCLUDE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INCLUDE FROM omni_kw_74
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_74 WHERE INCLUDE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_74 WHERE INCLUDE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_74 ORDER BY INCLUDE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_74 SET INCLUDE = NULL WHERE INCLUDE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_74 (INCLUDE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_74_i ON omni_kw_74 (INCLUDE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_74 WHERE INCLUDE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INCLUDE WHEN 1 THEN 2 END FROM omni_kw_74
< ok 1 0
> exec ^
: string
| SELECT CASE k.INCLUDE WHEN 1 THEN 2 END FROM omni_kw_74 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INCLUDE WHEN 2 THEN 3 END FROM omni_kw_74
< ok 1 0
> exec
| DROP TABLE omni_kw_74 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_75 (a NUMBER, INITIALLY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INITIALLY FROM omni_kw_75
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_75 WHERE INITIALLY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_75 WHERE INITIALLY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_75 ORDER BY INITIALLY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_75 SET INITIALLY = NULL WHERE INITIALLY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_75 (INITIALLY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_75_i ON omni_kw_75 (INITIALLY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_75 WHERE INITIALLY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INITIALLY WHEN 1 THEN 2 END FROM omni_kw_75
< ok 1 0
> exec ^
: string
| SELECT CASE k.INITIALLY WHEN 1 THEN 2 END FROM omni_kw_75 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INITIALLY WHEN 2 THEN 3 END FROM omni_kw_75
< ok 1 0
> exec
| DROP TABLE omni_kw_75 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_76 (a NUMBER, INNER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INNER FROM omni_kw_76
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_76 WHERE INNER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_76 WHERE INNER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_76 ORDER BY INNER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_76 SET INNER = NULL WHERE INNER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_76 (INNER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_76_i ON omni_kw_76 (INNER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_76 WHERE INNER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INNER WHEN 1 THEN 2 END FROM omni_kw_76
< ok 1 0
> exec ^
: string
| SELECT CASE k.INNER WHEN 1 THEN 2 END FROM omni_kw_76 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INNER WHEN 2 THEN 3 END FROM omni_kw_76
< ok 1 0
> exec
| DROP TABLE omni_kw_76 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_77 (a NUMBER, INSTEAD NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INSTEAD FROM omni_kw_77
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_77 WHERE INSTEAD IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_77 WHERE INSTEAD NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_77 ORDER BY INSTEAD
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_77 SET INSTEAD = NULL WHERE INSTEAD IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_77 (INSTEAD) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_77_i ON omni_kw_77 (INSTEAD)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_77 WHERE INSTEAD(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INSTEAD WHEN 1 THEN 2 END FROM omni_kw_77
< ok 1 0
> exec ^
: string
| SELECT CASE k.INSTEAD WHEN 1 THEN 2 END FROM omni_kw_77 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INSTEAD WHEN 2 THEN 3 END FROM omni_kw_77
< ok 1 0
> exec
| DROP TABLE omni_kw_77 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_78 (a NUMBER, INTERVAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INTERVAL FROM omni_kw_78
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_78 WHERE INTERVAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_78 WHERE INTERVAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_78 ORDER BY INTERVAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_78 SET INTERVAL = NULL WHERE INTERVAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_78 (INTERVAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_78_i ON omni_kw_78 (INTERVAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_78 WHERE INTERVAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INTERVAL WHEN 1 THEN 2 END FROM omni_kw_78
< ok 1 0
> exec ^
: string
| SELECT CASE k.INTERVAL WHEN 1 THEN 2 END FROM omni_kw_78 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INTERVAL WHEN 2 THEN 3 END FROM omni_kw_78
< ok 1 0
> exec
| DROP TABLE omni_kw_78 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_79 (a NUMBER, INVISIBLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT INVISIBLE FROM omni_kw_79
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_79 WHERE INVISIBLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_79 WHERE INVISIBLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_79 ORDER BY INVISIBLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_79 SET INVISIBLE = NULL WHERE INVISIBLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_79 (INVISIBLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_79_i ON omni_kw_79 (INVISIBLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_79 WHERE INVISIBLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE INVISIBLE WHEN 1 THEN 2 END FROM omni_kw_79
< ok 1 0
> exec ^
: string
| SELECT CASE k.INVISIBLE WHEN 1 THEN 2 END FROM omni_kw_79 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + INVISIBLE WHEN 2 THEN 3 END FROM omni_kw_79
< ok 1 0
> exec
| DROP TABLE omni_kw_79 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_80 (a NUMBER, ISOLATION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ISOLATION FROM omni_kw_80
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_80 WHERE ISOLATION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_80 WHERE ISOLATION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_80 ORDER BY ISOLATION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_80 SET ISOLATION = NULL WHERE ISOLATION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_80 (ISOLATION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_80_i ON omni_kw_80 (ISOLATION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_80 WHERE ISOLATION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ISOLATION WHEN 1 THEN 2 END FROM omni_kw_80
< ok 1 0
> exec ^
: string
| SELECT CASE k.ISOLATION WHEN 1 THEN 2 END FROM omni_kw_80 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ISOLATION WHEN 2 THEN 3 END FROM omni_kw_80
< ok 1 0
> exec
| DROP TABLE omni_kw_80 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_81 (a NUMBER, ITERATE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ITERATE FROM omni_kw_81
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_81 WHERE ITERATE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_81 WHERE ITERATE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_81 ORDER BY ITERATE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_81 SET ITERATE = NULL WHERE ITERATE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_81 (ITERATE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_81_i ON omni_kw_81 (ITERATE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_81 WHERE ITERATE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ITERATE WHEN 1 THEN 2 END FROM omni_kw_81
< ok 1 0
> exec ^
: string
| SELECT CASE k.ITERATE WHEN 1 THEN 2 END FROM omni_kw_81 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ITERATE WHEN 2 THEN 3 END FROM omni_kw_81
< ok 1 0
> exec
| DROP TABLE omni_kw_81 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_82 (a NUMBER, JAVA NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JAVA FROM omni_kw_82
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_82 WHERE JAVA IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_82 WHERE JAVA NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_82 ORDER BY JAVA
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_82 SET JAVA = NULL WHERE JAVA IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_82 (JAVA) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_82_i ON omni_kw_82 (JAVA)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_82 WHERE JAVA(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE JAVA WHEN 1 THEN 2 END FROM omni_kw_82
< ok 1 0
> exec ^
: string
| SELECT CASE k.JAVA WHEN 1 THEN 2 END FROM omni_kw_82 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JAVA WHEN 2 THEN 3 END FROM omni_kw_82
< ok 1 0
> exec
| DROP TABLE omni_kw_82 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_83 (a NUMBER, JOIN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JOIN FROM omni_kw_83
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_83 WHERE JOIN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_83 WHERE JOIN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_83 ORDER BY JOIN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_83 SET JOIN = NULL WHERE JOIN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_83 (JOIN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_83_i ON omni_kw_83 (JOIN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_83 WHERE JOIN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE JOIN WHEN 1 THEN 2 END FROM omni_kw_83
< ok 1 0
> exec ^
: string
| SELECT CASE k.JOIN WHEN 1 THEN 2 END FROM omni_kw_83 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JOIN WHEN 2 THEN 3 END FROM omni_kw_83
< ok 1 0
> exec
| DROP TABLE omni_kw_83 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_84 (a NUMBER, JSON NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON FROM omni_kw_84
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_84 WHERE JSON IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_84 WHERE JSON NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_84 ORDER BY JSON
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_84 SET JSON = NULL WHERE JSON IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_84 (JSON) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_84_i ON omni_kw_84 (JSON)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_84 WHERE JSON(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON WHEN 1 THEN 2 END FROM omni_kw_84
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON WHEN 1 THEN 2 END FROM omni_kw_84 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON WHEN 2 THEN 3 END FROM omni_kw_84
< ok 1 0
> exec
| DROP TABLE omni_kw_84 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_85 (a NUMBER, JSON_ARRAY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_ARRAY FROM omni_kw_85
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_85 WHERE JSON_ARRAY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_85 WHERE JSON_ARRAY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_85 ORDER BY JSON_ARRAY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_85 SET JSON_ARRAY = NULL WHERE JSON_ARRAY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_85 (JSON_ARRAY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_85_i ON omni_kw_85 (JSON_ARRAY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_85 WHERE JSON_ARRAY(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_ARRAY WHEN 1 THEN 2 END FROM omni_kw_85
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_ARRAY WHEN 1 THEN 2 END FROM omni_kw_85 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_ARRAY WHEN 2 THEN 3 END FROM omni_kw_85
< ok 1 0
> exec
| DROP TABLE omni_kw_85 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_86 (a NUMBER, JSON_EXISTS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_EXISTS FROM omni_kw_86
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_86 WHERE JSON_EXISTS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_86 WHERE JSON_EXISTS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_86 ORDER BY JSON_EXISTS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_86 SET JSON_EXISTS = NULL WHERE JSON_EXISTS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_86 (JSON_EXISTS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_86_i ON omni_kw_86 (JSON_EXISTS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_86 WHERE JSON_EXISTS(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_EXISTS WHEN 1 THEN 2 END FROM omni_kw_86
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_EXISTS WHEN 1 THEN 2 END FROM omni_kw_86 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_EXISTS WHEN 2 THEN 3 END FROM omni_kw_86
< ok 1 0
> exec
| DROP TABLE omni_kw_86 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_87 (a NUMBER, JSON_MERGEPATCH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_MERGEPATCH FROM omni_kw_87
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_87 WHERE JSON_MERGEPATCH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_87 WHERE JSON_MERGEPATCH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_87 ORDER BY JSON_MERGEPATCH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_87 SET JSON_MERGEPATCH = NULL WHERE JSON_MERGEPATCH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_87 (JSON_MERGEPATCH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_87_i ON omni_kw_87 (JSON_MERGEPATCH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_87 WHERE JSON_MERGEPATCH(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_MERGEPATCH WHEN 1 THEN 2 END FROM omni_kw_87
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_MERGEPATCH WHEN 1 THEN 2 END FROM omni_kw_87 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_MERGEPATCH WHEN 2 THEN 3 END FROM omni_kw_87
< ok 1 0
> exec
| DROP TABLE omni_kw_87 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_88 (a NUMBER, JSON_OBJECT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_OBJECT FROM omni_kw_88
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_88 WHERE JSON_OBJECT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_88 WHERE JSON_OBJECT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_88 ORDER BY JSON_OBJECT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_88 SET JSON_OBJECT = NULL WHERE JSON_OBJECT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_88 (JSON_OBJECT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_88_i ON omni_kw_88 (JSON_OBJECT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_88 WHERE JSON_OBJECT(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_OBJECT WHEN 1 THEN 2 END FROM omni_kw_88
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_OBJECT WHEN 1 THEN 2 END FROM omni_kw_88 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_OBJECT WHEN 2 THEN 3 END FROM omni_kw_88
< ok 1 0
> exec
| DROP TABLE omni_kw_88 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_89 (a NUMBER, JSON_QUERY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_QUERY FROM omni_kw_89
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_89 WHERE JSON_QUERY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_89 WHERE JSON_QUERY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_89 ORDER BY JSON_QUERY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_89 SET JSON_QUERY = NULL WHERE JSON_QUERY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_89 (JSON_QUERY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_89_i ON omni_kw_89 (JSON_QUERY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_89 WHERE JSON_QUERY(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_QUERY WHEN 1 THEN 2 END FROM omni_kw_89
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_QUERY WHEN 1 THEN 2 END FROM omni_kw_89 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_QUERY WHEN 2 THEN 3 END FROM omni_kw_89
< ok 1 0
> exec
| DROP TABLE omni_kw_89 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_90 (a NUMBER, JSON_TABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_TABLE FROM omni_kw_90
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_90 WHERE JSON_TABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_90 WHERE JSON_TABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_90 ORDER BY JSON_TABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_90 SET JSON_TABLE = NULL WHERE JSON_TABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_90 (JSON_TABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_90_i ON omni_kw_90 (JSON_TABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_90 WHERE JSON_TABLE(+) = LENGTH(d.dummy)
< error
| 61785 ORA-61785: Wrong usage of JSON_TABLE
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_TABLE WHEN 1 THEN 2 END FROM omni_kw_90
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_TABLE WHEN 1 THEN 2 END FROM omni_kw_90 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_TABLE WHEN 2 THEN 3 END FROM omni_kw_90
< ok 1 0
> exec
| DROP TABLE omni_kw_90 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_91 (a NUMBER, JSON_VALUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT JSON_VALUE FROM omni_kw_91
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_91 WHERE JSON_VALUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_91 WHERE JSON_VALUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_91 ORDER BY JSON_VALUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_91 SET JSON_VALUE = NULL WHERE JSON_VALUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_91 (JSON_VALUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_91_i ON omni_kw_91 (JSON_VALUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_91 WHERE JSON_VALUE(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE JSON_VALUE WHEN 1 THEN 2 END FROM omni_kw_91
< ok 1 0
> exec ^
: string
| SELECT CASE k.JSON_VALUE WHEN 1 THEN 2 END FROM omni_kw_91 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + JSON_VALUE WHEN 2 THEN 3 END FROM omni_kw_91
< ok 1 0
> exec
| DROP TABLE omni_kw_91 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_92 (a NUMBER, KEEP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT KEEP FROM omni_kw_92
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_92 WHERE KEEP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_92 WHERE KEEP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_92 ORDER BY KEEP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_92 SET KEEP = NULL WHERE KEEP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_92 (KEEP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_92_i ON omni_kw_92 (KEEP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_92 WHERE KEEP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE KEEP WHEN 1 THEN 2 END FROM omni_kw_92
< ok 1 0
> exec ^
: string
| SELECT CASE k.KEEP WHEN 1 THEN 2 END FROM omni_kw_92 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + KEEP WHEN 2 THEN 3 END FROM omni_kw_92
< ok 1 0
> exec
| DROP TABLE omni_kw_92 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_93 (a NUMBER, KEY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT KEY FROM omni_kw_93
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_93 WHERE KEY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_93 WHERE KEY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_93 ORDER BY KEY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_93 SET KEY = NULL WHERE KEY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_93 (KEY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_93_i ON omni_kw_93 (KEY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_93 WHERE KEY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE KEY WHEN 1 THEN 2 END FROM omni_kw_93
< ok 1 0
> exec ^
: string
| SELECT CASE k.KEY WHEN 1 THEN 2 END FROM omni_kw_93 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + KEY WHEN 2 THEN 3 END FROM omni_kw_93
< ok 1 0
> exec
| DROP TABLE omni_kw_93 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_94 (a NUMBER, LAST NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LAST FROM omni_kw_94
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_94 WHERE LAST IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_94 WHERE LAST NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_94 ORDER BY LAST
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_94 SET LAST = NULL WHERE LAST IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_94 (LAST) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_94_i ON omni_kw_94 (LAST)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_94 WHERE LAST(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LAST WHEN 1 THEN 2 END FROM omni_kw_94
< ok 1 0
> exec ^
: string
| SELECT CASE k.LAST WHEN 1 THEN 2 END FROM omni_kw_94 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LAST WHEN 2 THEN 3 END FROM omni_kw_94
< ok 1 0
> exec
| DROP TABLE omni_kw_94 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_95 (a NUMBER, LATERAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LATERAL FROM omni_kw_95
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_95 WHERE LATERAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_95 WHERE LATERAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_95 ORDER BY LATERAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_95 SET LATERAL = NULL WHERE LATERAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_95 (LATERAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_95_i ON omni_kw_95 (LATERAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_95 WHERE LATERAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LATERAL WHEN 1 THEN 2 END FROM omni_kw_95
< ok 1 0
> exec ^
: string
| SELECT CASE k.LATERAL WHEN 1 THEN 2 END FROM omni_kw_95 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LATERAL WHEN 2 THEN 3 END FROM omni_kw_95
< ok 1 0
> exec
| DROP TABLE omni_kw_95 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_96 (a NUMBER, LEFT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LEFT FROM omni_kw_96
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_96 WHERE LEFT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_96 WHERE LEFT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_96 ORDER BY LEFT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_96 SET LEFT = NULL WHERE LEFT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_96 (LEFT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_96_i ON omni_kw_96 (LEFT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_96 WHERE LEFT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LEFT WHEN 1 THEN 2 END FROM omni_kw_96
< ok 1 0
> exec ^
: string
| SELECT CASE k.LEFT WHEN 1 THEN 2 END FROM omni_kw_96 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LEFT WHEN 2 THEN 3 END FROM omni_kw_96
< ok 1 0
> exec
| DROP TABLE omni_kw_96 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_97 (a NUMBER, LIBRARY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIBRARY FROM omni_kw_97
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_97 WHERE LIBRARY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_97 WHERE LIBRARY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_97 ORDER BY LIBRARY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_97 SET LIBRARY = NULL WHERE LIBRARY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_97 (LIBRARY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_97_i ON omni_kw_97 (LIBRARY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_97 WHERE LIBRARY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIBRARY WHEN 1 THEN 2 END FROM omni_kw_97
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIBRARY WHEN 1 THEN 2 END FROM omni_kw_97 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIBRARY WHEN 2 THEN 3 END FROM omni_kw_97
< ok 1 0
> exec
| DROP TABLE omni_kw_97 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_98 (a NUMBER, LIKE2 NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIKE2 FROM omni_kw_98
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_98 WHERE LIKE2 IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_98 WHERE LIKE2 NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_98 ORDER BY LIKE2
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_98 SET LIKE2 = NULL WHERE LIKE2 IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_98 (LIKE2) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_98_i ON omni_kw_98 (LIKE2)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_98 WHERE LIKE2(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIKE2 WHEN 1 THEN 2 END FROM omni_kw_98
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIKE2 WHEN 1 THEN 2 END FROM omni_kw_98 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIKE2 WHEN 2 THEN 3 END FROM omni_kw_98
< ok 1 0
> exec
| DROP TABLE omni_kw_98 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_99 (a NUMBER, LIKE4 NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIKE4 FROM omni_kw_99
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_99 WHERE LIKE4 IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_99 WHERE LIKE4 NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_99 ORDER BY LIKE4
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_99 SET LIKE4 = NULL WHERE LIKE4 IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_99 (LIKE4) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_99_i ON omni_kw_99 (LIKE4)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_99 WHERE LIKE4(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIKE4 WHEN 1 THEN 2 END FROM omni_kw_99
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIKE4 WHEN 1 THEN 2 END FROM omni_kw_99 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIKE4 WHEN 2 THEN 3 END FROM omni_kw_99
< ok 1 0
> exec
| DROP TABLE omni_kw_99 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_100 (a NUMBER, LIKEC NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIKEC FROM omni_kw_100
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_100 WHERE LIKEC IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_100 WHERE LIKEC NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_100 ORDER BY LIKEC
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_100 SET LIKEC = NULL WHERE LIKEC IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_100 (LIKEC) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_100_i ON omni_kw_100 (LIKEC)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_100 WHERE LIKEC(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIKEC WHEN 1 THEN 2 END FROM omni_kw_100
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIKEC WHEN 1 THEN 2 END FROM omni_kw_100 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIKEC WHEN 2 THEN 3 END FROM omni_kw_100
< ok 1 0
> exec
| DROP TABLE omni_kw_100 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_101 (a NUMBER, LIMIT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIMIT FROM omni_kw_101
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_101 WHERE LIMIT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_101 WHERE LIMIT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_101 ORDER BY LIMIT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_101 SET LIMIT = NULL WHERE LIMIT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_101 (LIMIT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_101_i ON omni_kw_101 (LIMIT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_101 WHERE LIMIT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIMIT WHEN 1 THEN 2 END FROM omni_kw_101
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIMIT WHEN 1 THEN 2 END FROM omni_kw_101 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIMIT WHEN 2 THEN 3 END FROM omni_kw_101
< ok 1 0
> exec
| DROP TABLE omni_kw_101 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_102 (a NUMBER, LINK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LINK FROM omni_kw_102
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_102 WHERE LINK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_102 WHERE LINK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_102 ORDER BY LINK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_102 SET LINK = NULL WHERE LINK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_102 (LINK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_102_i ON omni_kw_102 (LINK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_102 WHERE LINK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LINK WHEN 1 THEN 2 END FROM omni_kw_102
< ok 1 0
> exec ^
: string
| SELECT CASE k.LINK WHEN 1 THEN 2 END FROM omni_kw_102 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LINK WHEN 2 THEN 3 END FROM omni_kw_102
< ok 1 0
> exec
| DROP TABLE omni_kw_102 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_103 (a NUMBER, LIST NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LIST FROM omni_kw_103
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_103 WHERE LIST IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_103 WHERE LIST NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_103 ORDER BY LIST
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_103 SET LIST = NULL WHERE LIST IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_103 (LIST) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_103_i ON omni_kw_103 (LIST)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_103 WHERE LIST(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LIST WHEN 1 THEN 2 END FROM omni_kw_103
< ok 1 0
> exec ^
: string
| SELECT CASE k.LIST WHEN 1 THEN 2 END FROM omni_kw_103 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LIST WHEN 2 THEN 3 END FROM omni_kw_103
< ok 1 0
> exec
| DROP TABLE omni_kw_103 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_104 (a NUMBER, LOCAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LOCAL FROM omni_kw_104
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_104 WHERE LOCAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_104 WHERE LOCAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_104 ORDER BY LOCAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_104 SET LOCAL = NULL WHERE LOCAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_104 (LOCAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_104_i ON omni_kw_104 (LOCAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_104 WHERE LOCAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LOCAL WHEN 1 THEN 2 END FROM omni_kw_104
< ok 1 0
> exec ^
: string
| SELECT CASE k.LOCAL WHEN 1 THEN 2 END FROM omni_kw_104 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LOCAL WHEN 2 THEN 3 END FROM omni_kw_104
< ok 1 0
> exec
| DROP TABLE omni_kw_104 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_105 (a NUMBER, LOCKED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LOCKED FROM omni_kw_105
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_105 WHERE LOCKED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_105 WHERE LOCKED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_105 ORDER BY LOCKED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_105 SET LOCKED = NULL WHERE LOCKED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_105 (LOCKED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_105_i ON omni_kw_105 (LOCKED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_105 WHERE LOCKED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LOCKED WHEN 1 THEN 2 END FROM omni_kw_105
< ok 1 0
> exec ^
: string
| SELECT CASE k.LOCKED WHEN 1 THEN 2 END FROM omni_kw_105 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LOCKED WHEN 2 THEN 3 END FROM omni_kw_105
< ok 1 0
> exec
| DROP TABLE omni_kw_105 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_106 (a NUMBER, LOG NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LOG FROM omni_kw_106
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_106 WHERE LOG IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_106 WHERE LOG NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_106 ORDER BY LOG
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_106 SET LOG = NULL WHERE LOG IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_106 (LOG) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_106_i ON omni_kw_106 (LOG)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_106 WHERE LOG(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LOG WHEN 1 THEN 2 END FROM omni_kw_106
< ok 1 0
> exec ^
: string
| SELECT CASE k.LOG WHEN 1 THEN 2 END FROM omni_kw_106 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LOG WHEN 2 THEN 3 END FROM omni_kw_106
< ok 1 0
> exec
| DROP TABLE omni_kw_106 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_107 (a NUMBER, LOGGING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LOGGING FROM omni_kw_107
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_107 WHERE LOGGING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_107 WHERE LOGGING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_107 ORDER BY LOGGING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_107 SET LOGGING = NULL WHERE LOGGING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_107 (LOGGING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_107_i ON omni_kw_107 (LOGGING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_107 WHERE LOGGING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LOGGING WHEN 1 THEN 2 END FROM omni_kw_107
< ok 1 0
> exec ^
: string
| SELECT CASE k.LOGGING WHEN 1 THEN 2 END FROM omni_kw_107 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LOGGING WHEN 2 THEN 3 END FROM omni_kw_107
< ok 1 0
> exec
| DROP TABLE omni_kw_107 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_108 (a NUMBER, LOOP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT LOOP FROM omni_kw_108
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_108 WHERE LOOP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_108 WHERE LOOP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_108 ORDER BY LOOP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_108 SET LOOP = NULL WHERE LOOP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_108 (LOOP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_108_i ON omni_kw_108 (LOOP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_108 WHERE LOOP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE LOOP WHEN 1 THEN 2 END FROM omni_kw_108
< ok 1 0
> exec ^
: string
| SELECT CASE k.LOOP WHEN 1 THEN 2 END FROM omni_kw_108 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + LOOP WHEN 2 THEN 3 END FROM omni_kw_108
< ok 1 0
> exec
| DROP TABLE omni_kw_108 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_109 (a NUMBER, MAIN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MAIN FROM omni_kw_109
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_109 WHERE MAIN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_109 WHERE MAIN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_109 ORDER BY MAIN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_109 SET MAIN = NULL WHERE MAIN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_109 (MAIN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_109_i ON omni_kw_109 (MAIN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_109 WHERE MAIN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MAIN WHEN 1 THEN 2 END FROM omni_kw_109
< ok 1 0
> exec ^
: string
| SELECT CASE k.MAIN WHEN 1 THEN 2 END FROM omni_kw_109 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MAIN WHEN 2 THEN 3 END FROM omni_kw_109
< ok 1 0
> exec
| DROP TABLE omni_kw_109 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_110 (a NUMBER, MATCHED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MATCHED FROM omni_kw_110
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_110 WHERE MATCHED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_110 WHERE MATCHED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_110 ORDER BY MATCHED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_110 SET MATCHED = NULL WHERE MATCHED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_110 (MATCHED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_110_i ON omni_kw_110 (MATCHED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_110 WHERE MATCHED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MATCHED WHEN 1 THEN 2 END FROM omni_kw_110
< ok 1 0
> exec ^
: string
| SELECT CASE k.MATCHED WHEN 1 THEN 2 END FROM omni_kw_110 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MATCHED WHEN 2 THEN 3 END FROM omni_kw_110
< ok 1 0
> exec
| DROP TABLE omni_kw_110 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_111 (a NUMBER, MATERIALIZED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MATERIALIZED FROM omni_kw_111
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_111 WHERE MATERIALIZED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_111 WHERE MATERIALIZED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_111 ORDER BY MATERIALIZED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_111 SET MATERIALIZED = NULL WHERE MATERIALIZED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_111 (MATERIALIZED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_111_i ON omni_kw_111 (MATERIALIZED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_111 WHERE MATERIALIZED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MATERIALIZED WHEN 1 THEN 2 END FROM omni_kw_111
< ok 1 0
> exec ^
: string
| SELECT CASE k.MATERIALIZED WHEN 1 THEN 2 END FROM omni_kw_111 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MATERIALIZED WHEN 2 THEN 3 END FROM omni_kw_111
< ok 1 0
> exec
| DROP TABLE omni_kw_111 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_112 (a NUMBER, MAXVALUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MAXVALUE FROM omni_kw_112
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_112 WHERE MAXVALUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_112 WHERE MAXVALUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_112 ORDER BY MAXVALUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_112 SET MAXVALUE = NULL WHERE MAXVALUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_112 (MAXVALUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_112_i ON omni_kw_112 (MAXVALUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_112 WHERE MAXVALUE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MAXVALUE WHEN 1 THEN 2 END FROM omni_kw_112
< ok 1 0
> exec ^
: string
| SELECT CASE k.MAXVALUE WHEN 1 THEN 2 END FROM omni_kw_112 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MAXVALUE WHEN 2 THEN 3 END FROM omni_kw_112
< ok 1 0
> exec
| DROP TABLE omni_kw_112 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_113 (a NUMBER, MEASURES NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MEASURES FROM omni_kw_113
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_113 WHERE MEASURES IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_113 WHERE MEASURES NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_113 ORDER BY MEASURES
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_113 SET MEASURES = NULL WHERE MEASURES IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_113 (MEASURES) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_113_i ON omni_kw_113 (MEASURES)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_113 WHERE MEASURES(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MEASURES WHEN 1 THEN 2 END FROM omni_kw_113
< ok 1 0
> exec ^
: string
| SELECT CASE k.MEASURES WHEN 1 THEN 2 END FROM omni_kw_113 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MEASURES WHEN 2 THEN 3 END FROM omni_kw_113
< ok 1 0
> exec
| DROP TABLE omni_kw_113 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_114 (a NUMBER, MERGE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MERGE FROM omni_kw_114
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_114 WHERE MERGE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_114 WHERE MERGE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_114 ORDER BY MERGE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_114 SET MERGE = NULL WHERE MERGE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_114 (MERGE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_114_i ON omni_kw_114 (MERGE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_114 WHERE MERGE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MERGE WHEN 1 THEN 2 END FROM omni_kw_114
< ok 1 0
> exec ^
: string
| SELECT CASE k.MERGE WHEN 1 THEN 2 END FROM omni_kw_114 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MERGE WHEN 2 THEN 3 END FROM omni_kw_114
< ok 1 0
> exec
| DROP TABLE omni_kw_114 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_115 (a NUMBER, MINVALUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MINVALUE FROM omni_kw_115
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_115 WHERE MINVALUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_115 WHERE MINVALUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_115 ORDER BY MINVALUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_115 SET MINVALUE = NULL WHERE MINVALUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_115 (MINVALUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_115_i ON omni_kw_115 (MINVALUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_115 WHERE MINVALUE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MINVALUE WHEN 1 THEN 2 END FROM omni_kw_115
< ok 1 0
> exec ^
: string
| SELECT CASE k.MINVALUE WHEN 1 THEN 2 END FROM omni_kw_115 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MINVALUE WHEN 2 THEN 3 END FROM omni_kw_115
< ok 1 0
> exec
| DROP TABLE omni_kw_115 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_116 (a NUMBER, MODEL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MODEL FROM omni_kw_116
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_116 WHERE MODEL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_116 WHERE MODEL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_116 ORDER BY MODEL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_116 SET MODEL = NULL WHERE MODEL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_116 (MODEL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_116_i ON omni_kw_116 (MODEL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_116 WHERE MODEL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MODEL WHEN 1 THEN 2 END FROM omni_kw_116
< ok 1 0
> exec ^
: string
| SELECT CASE k.MODEL WHEN 1 THEN 2 END FROM omni_kw_116 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MODEL WHEN 2 THEN 3 END FROM omni_kw_116
< ok 1 0
> exec
| DROP TABLE omni_kw_116 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_117 (a NUMBER, MULTISET NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT MULTISET FROM omni_kw_117
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_117 WHERE MULTISET IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_117 WHERE MULTISET NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_117 ORDER BY MULTISET
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_117 SET MULTISET = NULL WHERE MULTISET IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_117 (MULTISET) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_117_i ON omni_kw_117 (MULTISET)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_117 WHERE MULTISET(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE MULTISET WHEN 1 THEN 2 END FROM omni_kw_117
< ok 1 0
> exec ^
: string
| SELECT CASE k.MULTISET WHEN 1 THEN 2 END FROM omni_kw_117 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + MULTISET WHEN 2 THEN 3 END FROM omni_kw_117
< ok 1 0
> exec
| DROP TABLE omni_kw_117 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_118 (a NUMBER, NAME NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NAME FROM omni_kw_118
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_118 WHERE NAME IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_118 WHERE NAME NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_118 ORDER BY NAME
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_118 SET NAME = NULL WHERE NAME IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_118 (NAME) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_118_i ON omni_kw_118 (NAME)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_118 WHERE NAME(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NAME WHEN 1 THEN 2 END FROM omni_kw_118
< ok 1 0
> exec ^
: string
| SELECT CASE k.NAME WHEN 1 THEN 2 END FROM omni_kw_118 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NAME WHEN 2 THEN 3 END FROM omni_kw_118
< ok 1 0
> exec
| DROP TABLE omni_kw_118 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_119 (a NUMBER, NATURAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NATURAL FROM omni_kw_119
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_119 WHERE NATURAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_119 WHERE NATURAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_119 ORDER BY NATURAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_119 SET NATURAL = NULL WHERE NATURAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_119 (NATURAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_119_i ON omni_kw_119 (NATURAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_119 WHERE NATURAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NATURAL WHEN 1 THEN 2 END FROM omni_kw_119
< ok 1 0
> exec ^
: string
| SELECT CASE k.NATURAL WHEN 1 THEN 2 END FROM omni_kw_119 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NATURAL WHEN 2 THEN 3 END FROM omni_kw_119
< ok 1 0
> exec
| DROP TABLE omni_kw_119 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_120 (a NUMBER, NAV NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NAV FROM omni_kw_120
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_120 WHERE NAV IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_120 WHERE NAV NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_120 ORDER BY NAV
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_120 SET NAV = NULL WHERE NAV IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_120 (NAV) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_120_i ON omni_kw_120 (NAV)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_120 WHERE NAV(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NAV WHEN 1 THEN 2 END FROM omni_kw_120
< ok 1 0
> exec ^
: string
| SELECT CASE k.NAV WHEN 1 THEN 2 END FROM omni_kw_120 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NAV WHEN 2 THEN 3 END FROM omni_kw_120
< ok 1 0
> exec
| DROP TABLE omni_kw_120 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_121 (a NUMBER, NCHAR NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NCHAR FROM omni_kw_121
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_121 WHERE NCHAR IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_121 WHERE NCHAR NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_121 ORDER BY NCHAR
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_121 SET NCHAR = NULL WHERE NCHAR IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_121 (NCHAR) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_121_i ON omni_kw_121 (NCHAR)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_121 WHERE NCHAR(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NCHAR WHEN 1 THEN 2 END FROM omni_kw_121
< ok 1 0
> exec ^
: string
| SELECT CASE k.NCHAR WHEN 1 THEN 2 END FROM omni_kw_121 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NCHAR WHEN 2 THEN 3 END FROM omni_kw_121
< ok 1 0
> exec
| DROP TABLE omni_kw_121 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_122 (a NUMBER, NCLOB NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NCLOB FROM omni_kw_122
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_122 WHERE NCLOB IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_122 WHERE NCLOB NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_122 ORDER BY NCLOB
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_122 SET NCLOB = NULL WHERE NCLOB IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_122 (NCLOB) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_122_i ON omni_kw_122 (NCLOB)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_122 WHERE NCLOB(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NCLOB WHEN 1 THEN 2 END FROM omni_kw_122
< ok 1 0
> exec ^
: string
| SELECT CASE k.NCLOB WHEN 1 THEN 2 END FROM omni_kw_122 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NCLOB WHEN 2 THEN 3 END FROM omni_kw_122
< ok 1 0
> exec
| DROP TABLE omni_kw_122 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_123 (a NUMBER, NESTED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NESTED FROM omni_kw_123
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_123 WHERE NESTED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_123 WHERE NESTED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_123 ORDER BY NESTED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_123 SET NESTED = NULL WHERE NESTED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_123 (NESTED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_123_i ON omni_kw_123 (NESTED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_123 WHERE NESTED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NESTED WHEN 1 THEN 2 END FROM omni_kw_123
< ok 1 0
> exec ^
: string
| SELECT CASE k.NESTED WHEN 1 THEN 2 END FROM omni_kw_123 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NESTED WHEN 2 THEN 3 END FROM omni_kw_123
< ok 1 0
> exec
| DROP TABLE omni_kw_123 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_124 (a NUMBER, NEXT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NEXT FROM omni_kw_124
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_124 WHERE NEXT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_124 WHERE NEXT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_124 ORDER BY NEXT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_124 SET NEXT = NULL WHERE NEXT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_124 (NEXT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_124_i ON omni_kw_124 (NEXT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_124 WHERE NEXT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NEXT WHEN 1 THEN 2 END FROM omni_kw_124
< ok 1 0
> exec ^
: string
| SELECT CASE k.NEXT WHEN 1 THEN 2 END FROM omni_kw_124 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NEXT WHEN 2 THEN 3 END FROM omni_kw_124
< ok 1 0
> exec
| DROP TABLE omni_kw_124 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_125 (a NUMBER, NOCACHE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOCACHE FROM omni_kw_125
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_125 WHERE NOCACHE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_125 WHERE NOCACHE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_125 ORDER BY NOCACHE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_125 SET NOCACHE = NULL WHERE NOCACHE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_125 (NOCACHE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_125_i ON omni_kw_125 (NOCACHE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_125 WHERE NOCACHE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOCACHE WHEN 1 THEN 2 END FROM omni_kw_125
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOCACHE WHEN 1 THEN 2 END FROM omni_kw_125 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOCACHE WHEN 2 THEN 3 END FROM omni_kw_125
< ok 1 0
> exec
| DROP TABLE omni_kw_125 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_126 (a NUMBER, NOCOPY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOCOPY FROM omni_kw_126
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_126 WHERE NOCOPY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_126 WHERE NOCOPY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_126 ORDER BY NOCOPY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_126 SET NOCOPY = NULL WHERE NOCOPY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_126 (NOCOPY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_126_i ON omni_kw_126 (NOCOPY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_126 WHERE NOCOPY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOCOPY WHEN 1 THEN 2 END FROM omni_kw_126
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOCOPY WHEN 1 THEN 2 END FROM omni_kw_126 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOCOPY WHEN 2 THEN 3 END FROM omni_kw_126
< ok 1 0
> exec
| DROP TABLE omni_kw_126 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_127 (a NUMBER, NOCYCLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOCYCLE FROM omni_kw_127
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_127 WHERE NOCYCLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_127 WHERE NOCYCLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_127 ORDER BY NOCYCLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_127 SET NOCYCLE = NULL WHERE NOCYCLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_127 (NOCYCLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_127_i ON omni_kw_127 (NOCYCLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_127 WHERE NOCYCLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOCYCLE WHEN 1 THEN 2 END FROM omni_kw_127
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOCYCLE WHEN 1 THEN 2 END FROM omni_kw_127 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOCYCLE WHEN 2 THEN 3 END FROM omni_kw_127
< ok 1 0
> exec
| DROP TABLE omni_kw_127 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_128 (a NUMBER, NOLOGGING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOLOGGING FROM omni_kw_128
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_128 WHERE NOLOGGING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_128 WHERE NOLOGGING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_128 ORDER BY NOLOGGING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_128 SET NOLOGGING = NULL WHERE NOLOGGING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_128 (NOLOGGING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_128_i ON omni_kw_128 (NOLOGGING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_128 WHERE NOLOGGING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOLOGGING WHEN 1 THEN 2 END FROM omni_kw_128
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOLOGGING WHEN 1 THEN 2 END FROM omni_kw_128 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOLOGGING WHEN 2 THEN 3 END FROM omni_kw_128
< ok 1 0
> exec
| DROP TABLE omni_kw_128 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_129 (a NUMBER, NOMAXVALUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOMAXVALUE FROM omni_kw_129
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_129 WHERE NOMAXVALUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_129 WHERE NOMAXVALUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_129 ORDER BY NOMAXVALUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_129 SET NOMAXVALUE = NULL WHERE NOMAXVALUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_129 (NOMAXVALUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_129_i ON omni_kw_129 (NOMAXVALUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_129 WHERE NOMAXVALUE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOMAXVALUE WHEN 1 THEN 2 END FROM omni_kw_129
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOMAXVALUE WHEN 1 THEN 2 END FROM omni_kw_129 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOMAXVALUE WHEN 2 THEN 3 END FROM omni_kw_129
< ok 1 0
> exec
| DROP TABLE omni_kw_129 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_130 (a NUMBER, NOMINVALUE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOMINVALUE FROM omni_kw_130
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_130 WHERE NOMINVALUE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_130 WHERE NOMINVALUE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_130 ORDER BY NOMINVALUE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_130 SET NOMINVALUE = NULL WHERE NOMINVALUE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_130 (NOMINVALUE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_130_i ON omni_kw_130 (NOMINVALUE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_130 WHERE NOMINVALUE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOMINVALUE WHEN 1 THEN 2 END FROM omni_kw_130
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOMINVALUE WHEN 1 THEN 2 END FROM omni_kw_130 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOMINVALUE WHEN 2 THEN 3 END FROM omni_kw_130
< ok 1 0
> exec
| DROP TABLE omni_kw_130 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_131 (a NUMBER, NOORDER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOORDER FROM omni_kw_131
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_131 WHERE NOORDER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_131 WHERE NOORDER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_131 ORDER BY NOORDER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_131 SET NOORDER = NULL WHERE NOORDER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_131 (NOORDER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_131_i ON omni_kw_131 (NOORDER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_131 WHERE NOORDER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOORDER WHEN 1 THEN 2 END FROM omni_kw_131
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOORDER WHEN 1 THEN 2 END FROM omni_kw_131 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOORDER WHEN 2 THEN 3 END FROM omni_kw_131
< ok 1 0
> exec
| DROP TABLE omni_kw_131 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_132 (a NUMBER, NOPARALLEL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NOPARALLEL FROM omni_kw_132
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_132 WHERE NOPARALLEL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_132 WHERE NOPARALLEL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_132 ORDER BY NOPARALLEL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_132 SET NOPARALLEL = NULL WHERE NOPARALLEL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_132 (NOPARALLEL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_132_i ON omni_kw_132 (NOPARALLEL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_132 WHERE NOPARALLEL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NOPARALLEL WHEN 1 THEN 2 END FROM omni_kw_132
< ok 1 0
> exec ^
: string
| SELECT CASE k.NOPARALLEL WHEN 1 THEN 2 END FROM omni_kw_132 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NOPARALLEL WHEN 2 THEN 3 END FROM omni_kw_132
< ok 1 0
> exec
| DROP TABLE omni_kw_132 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_133 (a NUMBER, NULLS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NULLS FROM omni_kw_133
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_133 WHERE NULLS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_133 WHERE NULLS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_133 ORDER BY NULLS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_133 SET NULLS = NULL WHERE NULLS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_133 (NULLS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_133_i ON omni_kw_133 (NULLS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_133 WHERE NULLS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NULLS WHEN 1 THEN 2 END FROM omni_kw_133
< ok 1 0
> exec ^
: string
| SELECT CASE k.NULLS WHEN 1 THEN 2 END FROM omni_kw_133 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NULLS WHEN 2 THEN 3 END FROM omni_kw_133
< ok 1 0
> exec
| DROP TABLE omni_kw_133 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_134 (a NUMBER, NVARCHAR2 NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT NVARCHAR2 FROM omni_kw_134
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_134 WHERE NVARCHAR2 IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_134 WHERE NVARCHAR2 NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_134 ORDER BY NVARCHAR2
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_134 SET NVARCHAR2 = NULL WHERE NVARCHAR2 IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_134 (NVARCHAR2) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_134_i ON omni_kw_134 (NVARCHAR2)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_134 WHERE NVARCHAR2(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE NVARCHAR2 WHEN 1 THEN 2 END FROM omni_kw_134
< ok 1 0
> exec ^
: string
| SELECT CASE k.NVARCHAR2 WHEN 1 THEN 2 END FROM omni_kw_134 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + NVARCHAR2 WHEN 2 THEN 3 END FROM omni_kw_134
< ok 1 0
> exec
| DROP TABLE omni_kw_134 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_135 (a NUMBER, OFFSET NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT OFFSET FROM omni_kw_135
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_135 WHERE OFFSET IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_135 WHERE OFFSET NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_135 ORDER BY OFFSET
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_135 SET OFFSET = NULL WHERE OFFSET IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_135 (OFFSET) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_135_i ON omni_kw_135 (OFFSET)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_135 WHERE OFFSET(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE OFFSET WHEN 1 THEN 2 END FROM omni_kw_135
< ok 1 0
> exec ^
: string
| SELECT CASE k.OFFSET WHEN 1 THEN 2 END FROM omni_kw_135 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + OFFSET WHEN 2 THEN 3 END FROM omni_kw_135
< ok 1 0
> exec
| DROP TABLE omni_kw_135 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_136 (a NUMBER, ONLY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ONLY FROM omni_kw_136
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_136 WHERE ONLY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_136 WHERE ONLY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_136 ORDER BY ONLY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_136 SET ONLY = NULL WHERE ONLY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_136 (ONLY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_136_i ON omni_kw_136 (ONLY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_136 WHERE ONLY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ONLY WHEN 1 THEN 2 END FROM omni_kw_136
< ok 1 0
> exec ^
: string
| SELECT CASE k.ONLY WHEN 1 THEN 2 END FROM omni_kw_136 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ONLY WHEN 2 THEN 3 END FROM omni_kw_136
< ok 1 0
> exec
| DROP TABLE omni_kw_136 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_137 (a NUMBER, OPEN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT OPEN FROM omni_kw_137
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_137 WHERE OPEN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_137 WHERE OPEN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_137 ORDER BY OPEN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_137 SET OPEN = NULL WHERE OPEN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_137 (OPEN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_137_i ON omni_kw_137 (OPEN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_137 WHERE OPEN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE OPEN WHEN 1 THEN 2 END FROM omni_kw_137
< ok 1 0
> exec ^
: string
| SELECT CASE k.OPEN WHEN 1 THEN 2 END FROM omni_kw_137 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + OPEN WHEN 2 THEN 3 END FROM omni_kw_137
< ok 1 0
> exec
| DROP TABLE omni_kw_137 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_138 (a NUMBER, ORDINALITY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ORDINALITY FROM omni_kw_138
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_138 WHERE ORDINALITY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_138 WHERE ORDINALITY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_138 ORDER BY ORDINALITY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_138 SET ORDINALITY = NULL WHERE ORDINALITY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_138 (ORDINALITY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_138_i ON omni_kw_138 (ORDINALITY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_138 WHERE ORDINALITY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ORDINALITY WHEN 1 THEN 2 END FROM omni_kw_138
< ok 1 0
> exec ^
: string
| SELECT CASE k.ORDINALITY WHEN 1 THEN 2 END FROM omni_kw_138 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ORDINALITY WHEN 2 THEN 3 END FROM omni_kw_138
< ok 1 0
> exec
| DROP TABLE omni_kw_138 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_139 (a NUMBER, OUT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT OUT FROM omni_kw_139
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_139 WHERE OUT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_139 WHERE OUT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_139 ORDER BY OUT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_139 SET OUT = NULL WHERE OUT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_139 (OUT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_139_i ON omni_kw_139 (OUT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_139 WHERE OUT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE OUT WHEN 1 THEN 2 END FROM omni_kw_139
< ok 1 0
> exec ^
: string
| SELECT CASE k.OUT WHEN 1 THEN 2 END FROM omni_kw_139 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + OUT WHEN 2 THEN 3 END FROM omni_kw_139
< ok 1 0
> exec
| DROP TABLE omni_kw_139 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_140 (a NUMBER, OUTER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT OUTER FROM omni_kw_140
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_140 WHERE OUTER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_140 WHERE OUTER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_140 ORDER BY OUTER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_140 SET OUTER = NULL WHERE OUTER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_140 (OUTER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_140_i ON omni_kw_140 (OUTER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_140 WHERE OUTER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE OUTER WHEN 1 THEN 2 END FROM omni_kw_140
< ok 1 0
> exec ^
: string
| SELECT CASE k.OUTER WHEN 1 THEN 2 END FROM omni_kw_140 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + OUTER WHEN 2 THEN 3 END FROM omni_kw_140
< ok 1 0
> exec
| DROP TABLE omni_kw_140 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_141 (a NUMBER, OVER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT OVER FROM omni_kw_141
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_141 WHERE OVER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_141 WHERE OVER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_141 ORDER BY OVER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_141 SET OVER = NULL WHERE OVER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_141 (OVER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_141_i ON omni_kw_141 (OVER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_141 WHERE OVER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE OVER WHEN 1 THEN 2 END FROM omni_kw_141
< ok 1 0
> exec ^
: string
| SELECT CASE k.OVER WHEN 1 THEN 2 END FROM omni_kw_141 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + OVER WHEN 2 THEN 3 END FROM omni_kw_141
< ok 1 0
> exec
| DROP TABLE omni_kw_141 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_142 (a NUMBER, PACKAGE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PACKAGE FROM omni_kw_142
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_142 WHERE PACKAGE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_142 WHERE PACKAGE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_142 ORDER BY PACKAGE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_142 SET PACKAGE = NULL WHERE PACKAGE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_142 (PACKAGE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_142_i ON omni_kw_142 (PACKAGE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_142 WHERE PACKAGE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PACKAGE WHEN 1 THEN 2 END FROM omni_kw_142
< ok 1 0
> exec ^
: string
| SELECT CASE k.PACKAGE WHEN 1 THEN 2 END FROM omni_kw_142 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PACKAGE WHEN 2 THEN 3 END FROM omni_kw_142
< ok 1 0
> exec
| DROP TABLE omni_kw_142 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_143 (a NUMBER, PARALLEL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PARALLEL FROM omni_kw_143
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_143 WHERE PARALLEL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_143 WHERE PARALLEL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_143 ORDER BY PARALLEL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_143 SET PARALLEL = NULL WHERE PARALLEL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_143 (PARALLEL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_143_i ON omni_kw_143 (PARALLEL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_143 WHERE PARALLEL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PARALLEL WHEN 1 THEN 2 END FROM omni_kw_143
< ok 1 0
> exec ^
: string
| SELECT CASE k.PARALLEL WHEN 1 THEN 2 END FROM omni_kw_143 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PARALLEL WHEN 2 THEN 3 END FROM omni_kw_143
< ok 1 0
> exec
| DROP TABLE omni_kw_143 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_144 (a NUMBER, PARALLEL_ENABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PARALLEL_ENABLE FROM omni_kw_144
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_144 WHERE PARALLEL_ENABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_144 WHERE PARALLEL_ENABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_144 ORDER BY PARALLEL_ENABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_144 SET PARALLEL_ENABLE = NULL WHERE PARALLEL_ENABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_144 (PARALLEL_ENABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_144_i ON omni_kw_144 (PARALLEL_ENABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_144 WHERE PARALLEL_ENABLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PARALLEL_ENABLE WHEN 1 THEN 2 END FROM omni_kw_144
< ok 1 0
> exec ^
: string
| SELECT CASE k.PARALLEL_ENABLE WHEN 1 THEN 2 END FROM omni_kw_144 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PARALLEL_ENABLE WHEN 2 THEN 3 END FROM omni_kw_144
< ok 1 0
> exec
| DROP TABLE omni_kw_144 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_145 (a NUMBER, PARTITION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PARTITION FROM omni_kw_145
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_145 WHERE PARTITION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_145 WHERE PARTITION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_145 ORDER BY PARTITION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_145 SET PARTITION = NULL WHERE PARTITION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_145 (PARTITION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_145_i ON omni_kw_145 (PARTITION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_145 WHERE PARTITION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PARTITION WHEN 1 THEN 2 END FROM omni_kw_145
< ok 1 0
> exec ^
: string
| SELECT CASE k.PARTITION WHEN 1 THEN 2 END FROM omni_kw_145 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PARTITION WHEN 2 THEN 3 END FROM omni_kw_145
< ok 1 0
> exec
| DROP TABLE omni_kw_145 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_146 (a NUMBER, PASSING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PASSING FROM omni_kw_146
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_146 WHERE PASSING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_146 WHERE PASSING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_146 ORDER BY PASSING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_146 SET PASSING = NULL WHERE PASSING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_146 (PASSING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_146_i ON omni_kw_146 (PASSING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_146 WHERE PASSING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PASSING WHEN 1 THEN 2 END FROM omni_kw_146
< ok 1 0
> exec ^
: string
| SELECT CASE k.PASSING WHEN 1 THEN 2 END FROM omni_kw_146 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PASSING WHEN 2 THEN 3 END FROM omni_kw_146
< ok 1 0
> exec
| DROP TABLE omni_kw_146 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_147 (a NUMBER, PATH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PATH FROM omni_kw_147
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_147 WHERE PATH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_147 WHERE PATH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_147 ORDER BY PATH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_147 SET PATH = NULL WHERE PATH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_147 (PATH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_147_i ON omni_kw_147 (PATH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_147 WHERE PATH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PATH WHEN 1 THEN 2 END FROM omni_kw_147
< ok 1 0
> exec ^
: string
| SELECT CASE k.PATH WHEN 1 THEN 2 END FROM omni_kw_147 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PATH WHEN 2 THEN 3 END FROM omni_kw_147
< ok 1 0
> exec
| DROP TABLE omni_kw_147 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_148 (a NUMBER, PERCENT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PERCENT FROM omni_kw_148
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_148 WHERE PERCENT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_148 WHERE PERCENT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_148 ORDER BY PERCENT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_148 SET PERCENT = NULL WHERE PERCENT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_148 (PERCENT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_148_i ON omni_kw_148 (PERCENT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_148 WHERE PERCENT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PERCENT WHEN 1 THEN 2 END FROM omni_kw_148
< ok 1 0
> exec ^
: string
| SELECT CASE k.PERCENT WHEN 1 THEN 2 END FROM omni_kw_148 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PERCENT WHEN 2 THEN 3 END FROM omni_kw_148
< ok 1 0
> exec
| DROP TABLE omni_kw_148 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_149 (a NUMBER, PIPE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PIPE FROM omni_kw_149
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_149 WHERE PIPE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_149 WHERE PIPE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_149 ORDER BY PIPE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_149 SET PIPE = NULL WHERE PIPE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_149 (PIPE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_149_i ON omni_kw_149 (PIPE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_149 WHERE PIPE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PIPE WHEN 1 THEN 2 END FROM omni_kw_149
< ok 1 0
> exec ^
: string
| SELECT CASE k.PIPE WHEN 1 THEN 2 END FROM omni_kw_149 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PIPE WHEN 2 THEN 3 END FROM omni_kw_149
< ok 1 0
> exec
| DROP TABLE omni_kw_149 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_150 (a NUMBER, PIPELINED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PIPELINED FROM omni_kw_150
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_150 WHERE PIPELINED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_150 WHERE PIPELINED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_150 ORDER BY PIPELINED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_150 SET PIPELINED = NULL WHERE PIPELINED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_150 (PIPELINED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_150_i ON omni_kw_150 (PIPELINED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_150 WHERE PIPELINED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PIPELINED WHEN 1 THEN 2 END FROM omni_kw_150
< ok 1 0
> exec ^
: string
| SELECT CASE k.PIPELINED WHEN 1 THEN 2 END FROM omni_kw_150 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PIPELINED WHEN 2 THEN 3 END FROM omni_kw_150
< ok 1 0
> exec
| DROP TABLE omni_kw_150 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_151 (a NUMBER, PIVOT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PIVOT FROM omni_kw_151
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_151 WHERE PIVOT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_151 WHERE PIVOT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_151 ORDER BY PIVOT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_151 SET PIVOT = NULL WHERE PIVOT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_151 (PIVOT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_151_i ON omni_kw_151 (PIVOT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_151 WHERE PIVOT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PIVOT WHEN 1 THEN 2 END FROM omni_kw_151
< ok 1 0
> exec ^
: string
| SELECT CASE k.PIVOT WHEN 1 THEN 2 END FROM omni_kw_151 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PIVOT WHEN 2 THEN 3 END FROM omni_kw_151
< ok 1 0
> exec
| DROP TABLE omni_kw_151 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_152 (a NUMBER, PLAN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PLAN FROM omni_kw_152
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_152 WHERE PLAN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_152 WHERE PLAN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_152 ORDER BY PLAN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_152 SET PLAN = NULL WHERE PLAN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_152 (PLAN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_152_i ON omni_kw_152 (PLAN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_152 WHERE PLAN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PLAN WHEN 1 THEN 2 END FROM omni_kw_152
< ok 1 0
> exec ^
: string
| SELECT CASE k.PLAN WHEN 1 THEN 2 END FROM omni_kw_152 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PLAN WHEN 2 THEN 3 END FROM omni_kw_152
< ok 1 0
> exec
| DROP TABLE omni_kw_152 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_153 (a NUMBER, PRAGMA NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PRAGMA FROM omni_kw_153
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_153 WHERE PRAGMA IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_153 WHERE PRAGMA NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_153 ORDER BY PRAGMA
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_153 SET PRAGMA = NULL WHERE PRAGMA IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_153 (PRAGMA) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_153_i ON omni_kw_153 (PRAGMA)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_153 WHERE PRAGMA(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PRAGMA WHEN 1 THEN 2 END FROM omni_kw_153
< ok 1 0
> exec ^
: string
| SELECT CASE k.PRAGMA WHEN 1 THEN 2 END FROM omni_kw_153 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PRAGMA WHEN 2 THEN 3 END FROM omni_kw_153
< ok 1 0
> exec
| DROP TABLE omni_kw_153 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_154 (a NUMBER, PRECEDING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PRECEDING FROM omni_kw_154
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_154 WHERE PRECEDING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_154 WHERE PRECEDING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_154 ORDER BY PRECEDING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_154 SET PRECEDING = NULL WHERE PRECEDING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_154 (PRECEDING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_154_i ON omni_kw_154 (PRECEDING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_154 WHERE PRECEDING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PRECEDING WHEN 1 THEN 2 END FROM omni_kw_154
< ok 1 0
> exec ^
: string
| SELECT CASE k.PRECEDING WHEN 1 THEN 2 END FROM omni_kw_154 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PRECEDING WHEN 2 THEN 3 END FROM omni_kw_154
< ok 1 0
> exec
| DROP TABLE omni_kw_154 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_155 (a NUMBER, PRIMARY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PRIMARY FROM omni_kw_155
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_155 WHERE PRIMARY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_155 WHERE PRIMARY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_155 ORDER BY PRIMARY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_155 SET PRIMARY = NULL WHERE PRIMARY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_155 (PRIMARY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_155_i ON omni_kw_155 (PRIMARY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_155 WHERE PRIMARY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PRIMARY WHEN 1 THEN 2 END FROM omni_kw_155
< ok 1 0
> exec ^
: string
| SELECT CASE k.PRIMARY WHEN 1 THEN 2 END FROM omni_kw_155 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PRIMARY WHEN 2 THEN 3 END FROM omni_kw_155
< ok 1 0
> exec
| DROP TABLE omni_kw_155 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_156 (a NUMBER, PRIVATE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PRIVATE FROM omni_kw_156
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_156 WHERE PRIVATE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_156 WHERE PRIVATE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_156 ORDER BY PRIVATE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_156 SET PRIVATE = NULL WHERE PRIVATE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_156 (PRIVATE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_156_i ON omni_kw_156 (PRIVATE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_156 WHERE PRIVATE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PRIVATE WHEN 1 THEN 2 END FROM omni_kw_156
< ok 1 0
> exec ^
: string
| SELECT CASE k.PRIVATE WHEN 1 THEN 2 END FROM omni_kw_156 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PRIVATE WHEN 2 THEN 3 END FROM omni_kw_156
< ok 1 0
> exec
| DROP TABLE omni_kw_156 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_157 (a NUMBER, PRIVILEGES NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PRIVILEGES FROM omni_kw_157
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_157 WHERE PRIVILEGES IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_157 WHERE PRIVILEGES NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_157 ORDER BY PRIVILEGES
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_157 SET PRIVILEGES = NULL WHERE PRIVILEGES IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_157 (PRIVILEGES) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_157_i ON omni_kw_157 (PRIVILEGES)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_157 WHERE PRIVILEGES(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PRIVILEGES WHEN 1 THEN 2 END FROM omni_kw_157
< ok 1 0
> exec ^
: string
| SELECT CASE k.PRIVILEGES WHEN 1 THEN 2 END FROM omni_kw_157 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PRIVILEGES WHEN 2 THEN 3 END FROM omni_kw_157
< ok 1 0
> exec
| DROP TABLE omni_kw_157 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_158 (a NUMBER, PROCEDURE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PROCEDURE FROM omni_kw_158
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_158 WHERE PROCEDURE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_158 WHERE PROCEDURE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_158 ORDER BY PROCEDURE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_158 SET PROCEDURE = NULL WHERE PROCEDURE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_158 (PROCEDURE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_158_i ON omni_kw_158 (PROCEDURE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_158 WHERE PROCEDURE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PROCEDURE WHEN 1 THEN 2 END FROM omni_kw_158
< ok 1 0
> exec ^
: string
| SELECT CASE k.PROCEDURE WHEN 1 THEN 2 END FROM omni_kw_158 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PROCEDURE WHEN 2 THEN 3 END FROM omni_kw_158
< ok 1 0
> exec
| DROP TABLE omni_kw_158 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_159 (a NUMBER, PROFILE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PROFILE FROM omni_kw_159
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_159 WHERE PROFILE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_159 WHERE PROFILE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_159 ORDER BY PROFILE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_159 SET PROFILE = NULL WHERE PROFILE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_159 (PROFILE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_159_i ON omni_kw_159 (PROFILE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_159 WHERE PROFILE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PROFILE WHEN 1 THEN 2 END FROM omni_kw_159
< ok 1 0
> exec ^
: string
| SELECT CASE k.PROFILE WHEN 1 THEN 2 END FROM omni_kw_159 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PROFILE WHEN 2 THEN 3 END FROM omni_kw_159
< ok 1 0
> exec
| DROP TABLE omni_kw_159 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_160 (a NUMBER, PURGE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT PURGE FROM omni_kw_160
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_160 WHERE PURGE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_160 WHERE PURGE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_160 ORDER BY PURGE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_160 SET PURGE = NULL WHERE PURGE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_160 (PURGE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_160_i ON omni_kw_160 (PURGE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_160 WHERE PURGE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE PURGE WHEN 1 THEN 2 END FROM omni_kw_160
< ok 1 0
> exec ^
: string
| SELECT CASE k.PURGE WHEN 1 THEN 2 END FROM omni_kw_160 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + PURGE WHEN 2 THEN 3 END FROM omni_kw_160
< ok 1 0
> exec
| DROP TABLE omni_kw_160 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_161 (a NUMBER, RAISE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RAISE FROM omni_kw_161
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_161 WHERE RAISE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_161 WHERE RAISE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_161 ORDER BY RAISE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_161 SET RAISE = NULL WHERE RAISE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_161 (RAISE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_161_i ON omni_kw_161 (RAISE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_161 WHERE RAISE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RAISE WHEN 1 THEN 2 END FROM omni_kw_161
< ok 1 0
> exec ^
: string
| SELECT CASE k.RAISE WHEN 1 THEN 2 END FROM omni_kw_161 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RAISE WHEN 2 THEN 3 END FROM omni_kw_161
< ok 1 0
> exec
| DROP TABLE omni_kw_161 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_162 (a NUMBER, RANGE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RANGE FROM omni_kw_162
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_162 WHERE RANGE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_162 WHERE RANGE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_162 ORDER BY RANGE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_162 SET RANGE = NULL WHERE RANGE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_162 (RANGE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_162_i ON omni_kw_162 (RANGE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_162 WHERE RANGE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RANGE WHEN 1 THEN 2 END FROM omni_kw_162
< ok 1 0
> exec ^
: string
| SELECT CASE k.RANGE WHEN 1 THEN 2 END FROM omni_kw_162 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RANGE WHEN 2 THEN 3 END FROM omni_kw_162
< ok 1 0
> exec
| DROP TABLE omni_kw_162 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_163 (a NUMBER, READ NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT READ FROM omni_kw_163
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_163 WHERE READ IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_163 WHERE READ NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_163 ORDER BY READ
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_163 SET READ = NULL WHERE READ IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_163 (READ) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_163_i ON omni_kw_163 (READ)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_163 WHERE READ(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE READ WHEN 1 THEN 2 END FROM omni_kw_163
< ok 1 0
> exec ^
: string
| SELECT CASE k.READ WHEN 1 THEN 2 END FROM omni_kw_163 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + READ WHEN 2 THEN 3 END FROM omni_kw_163
< ok 1 0
> exec
| DROP TABLE omni_kw_163 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_164 (a NUMBER, RECORD NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RECORD FROM omni_kw_164
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_164 WHERE RECORD IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_164 WHERE RECORD NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_164 ORDER BY RECORD
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_164 SET RECORD = NULL WHERE RECORD IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_164 (RECORD) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_164_i ON omni_kw_164 (RECORD)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_164 WHERE RECORD(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RECORD WHEN 1 THEN 2 END FROM omni_kw_164
< ok 1 0
> exec ^
: string
| SELECT CASE k.RECORD WHEN 1 THEN 2 END FROM omni_kw_164 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RECORD WHEN 2 THEN 3 END FROM omni_kw_164
< ok 1 0
> exec
| DROP TABLE omni_kw_164 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_165 (a NUMBER, RECURSIVE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RECURSIVE FROM omni_kw_165
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_165 WHERE RECURSIVE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_165 WHERE RECURSIVE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_165 ORDER BY RECURSIVE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_165 SET RECURSIVE = NULL WHERE RECURSIVE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_165 (RECURSIVE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_165_i ON omni_kw_165 (RECURSIVE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_165 WHERE RECURSIVE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RECURSIVE WHEN 1 THEN 2 END FROM omni_kw_165
< ok 1 0
> exec ^
: string
| SELECT CASE k.RECURSIVE WHEN 1 THEN 2 END FROM omni_kw_165 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RECURSIVE WHEN 2 THEN 3 END FROM omni_kw_165
< ok 1 0
> exec
| DROP TABLE omni_kw_165 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_166 (a NUMBER, REF NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REF FROM omni_kw_166
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_166 WHERE REF IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_166 WHERE REF NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_166 ORDER BY REF
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_166 SET REF = NULL WHERE REF IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_166 (REF) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_166_i ON omni_kw_166 (REF)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_166 WHERE REF(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REF WHEN 1 THEN 2 END FROM omni_kw_166
< ok 1 0
> exec ^
: string
| SELECT CASE k.REF WHEN 1 THEN 2 END FROM omni_kw_166 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REF WHEN 2 THEN 3 END FROM omni_kw_166
< ok 1 0
> exec
| DROP TABLE omni_kw_166 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_167 (a NUMBER, REFERENCE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REFERENCE FROM omni_kw_167
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_167 WHERE REFERENCE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_167 WHERE REFERENCE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_167 ORDER BY REFERENCE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_167 SET REFERENCE = NULL WHERE REFERENCE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_167 (REFERENCE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_167_i ON omni_kw_167 (REFERENCE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_167 WHERE REFERENCE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REFERENCE WHEN 1 THEN 2 END FROM omni_kw_167
< ok 1 0
> exec ^
: string
| SELECT CASE k.REFERENCE WHEN 1 THEN 2 END FROM omni_kw_167 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REFERENCE WHEN 2 THEN 3 END FROM omni_kw_167
< ok 1 0
> exec
| DROP TABLE omni_kw_167 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_168 (a NUMBER, REFERENCES NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REFERENCES FROM omni_kw_168
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_168 WHERE REFERENCES IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_168 WHERE REFERENCES NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_168 ORDER BY REFERENCES
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_168 SET REFERENCES = NULL WHERE REFERENCES IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_168 (REFERENCES) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_168_i ON omni_kw_168 (REFERENCES)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_168 WHERE REFERENCES(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REFERENCES WHEN 1 THEN 2 END FROM omni_kw_168
< ok 1 0
> exec ^
: string
| SELECT CASE k.REFERENCES WHEN 1 THEN 2 END FROM omni_kw_168 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REFERENCES WHEN 2 THEN 3 END FROM omni_kw_168
< ok 1 0
> exec
| DROP TABLE omni_kw_168 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_169 (a NUMBER, REFRESH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REFRESH FROM omni_kw_169
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_169 WHERE REFRESH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_169 WHERE REFRESH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_169 ORDER BY REFRESH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_169 SET REFRESH = NULL WHERE REFRESH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_169 (REFRESH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_169_i ON omni_kw_169 (REFRESH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_169 WHERE REFRESH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REFRESH WHEN 1 THEN 2 END FROM omni_kw_169
< ok 1 0
> exec ^
: string
| SELECT CASE k.REFRESH WHEN 1 THEN 2 END FROM omni_kw_169 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REFRESH WHEN 2 THEN 3 END FROM omni_kw_169
< ok 1 0
> exec
| DROP TABLE omni_kw_169 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_170 (a NUMBER, REJECT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REJECT FROM omni_kw_170
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_170 WHERE REJECT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_170 WHERE REJECT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_170 ORDER BY REJECT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_170 SET REJECT = NULL WHERE REJECT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_170 (REJECT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_170_i ON omni_kw_170 (REJECT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_170 WHERE REJECT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REJECT WHEN 1 THEN 2 END FROM omni_kw_170
< ok 1 0
> exec ^
: string
| SELECT CASE k.REJECT WHEN 1 THEN 2 END FROM omni_kw_170 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REJECT WHEN 2 THEN 3 END FROM omni_kw_170
< ok 1 0
> exec
| DROP TABLE omni_kw_170 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_171 (a NUMBER, RELY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RELY FROM omni_kw_171
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_171 WHERE RELY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_171 WHERE RELY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_171 ORDER BY RELY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_171 SET RELY = NULL WHERE RELY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_171 (RELY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_171_i ON omni_kw_171 (RELY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_171 WHERE RELY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RELY WHEN 1 THEN 2 END FROM omni_kw_171
< ok 1 0
> exec ^
: string
| SELECT CASE k.RELY WHEN 1 THEN 2 END FROM omni_kw_171 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RELY WHEN 2 THEN 3 END FROM omni_kw_171
< ok 1 0
> exec
| DROP TABLE omni_kw_171 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_172 (a NUMBER, REPLACE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REPLACE FROM omni_kw_172
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_172 WHERE REPLACE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_172 WHERE REPLACE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_172 ORDER BY REPLACE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_172 SET REPLACE = NULL WHERE REPLACE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_172 (REPLACE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_172_i ON omni_kw_172 (REPLACE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_172 WHERE REPLACE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REPLACE WHEN 1 THEN 2 END FROM omni_kw_172
< ok 1 0
> exec ^
: string
| SELECT CASE k.REPLACE WHEN 1 THEN 2 END FROM omni_kw_172 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REPLACE WHEN 2 THEN 3 END FROM omni_kw_172
< ok 1 0
> exec
| DROP TABLE omni_kw_172 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_173 (a NUMBER, RESTRICT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RESTRICT FROM omni_kw_173
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_173 WHERE RESTRICT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_173 WHERE RESTRICT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_173 ORDER BY RESTRICT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_173 SET RESTRICT = NULL WHERE RESTRICT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_173 (RESTRICT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_173_i ON omni_kw_173 (RESTRICT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_173 WHERE RESTRICT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RESTRICT WHEN 1 THEN 2 END FROM omni_kw_173
< ok 1 0
> exec ^
: string
| SELECT CASE k.RESTRICT WHEN 1 THEN 2 END FROM omni_kw_173 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RESTRICT WHEN 2 THEN 3 END FROM omni_kw_173
< ok 1 0
> exec
| DROP TABLE omni_kw_173 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_174 (a NUMBER, RESULT_CACHE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RESULT_CACHE FROM omni_kw_174
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_174 WHERE RESULT_CACHE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_174 WHERE RESULT_CACHE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_174 ORDER BY RESULT_CACHE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_174 SET RESULT_CACHE = NULL WHERE RESULT_CACHE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_174 (RESULT_CACHE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_174_i ON omni_kw_174 (RESULT_CACHE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_174 WHERE RESULT_CACHE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RESULT_CACHE WHEN 1 THEN 2 END FROM omni_kw_174
< ok 1 0
> exec ^
: string
| SELECT CASE k.RESULT_CACHE WHEN 1 THEN 2 END FROM omni_kw_174 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RESULT_CACHE WHEN 2 THEN 3 END FROM omni_kw_174
< ok 1 0
> exec
| DROP TABLE omni_kw_174 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_175 (a NUMBER, RETURN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RETURN FROM omni_kw_175
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_175 WHERE RETURN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_175 WHERE RETURN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_175 ORDER BY RETURN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_175 SET RETURN = NULL WHERE RETURN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_175 (RETURN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_175_i ON omni_kw_175 (RETURN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_175 WHERE RETURN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RETURN WHEN 1 THEN 2 END FROM omni_kw_175
< ok 1 0
> exec ^
: string
| SELECT CASE k.RETURN WHEN 1 THEN 2 END FROM omni_kw_175 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RETURN WHEN 2 THEN 3 END FROM omni_kw_175
< ok 1 0
> exec
| DROP TABLE omni_kw_175 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_176 (a NUMBER, RETURNING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RETURNING FROM omni_kw_176
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_176 WHERE RETURNING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_176 WHERE RETURNING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_176 ORDER BY RETURNING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_176 SET RETURNING = NULL WHERE RETURNING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_176 (RETURNING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_176_i ON omni_kw_176 (RETURNING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_176 WHERE RETURNING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RETURNING WHEN 1 THEN 2 END FROM omni_kw_176
< ok 1 0
> exec ^
: string
| SELECT CASE k.RETURNING WHEN 1 THEN 2 END FROM omni_kw_176 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RETURNING WHEN 2 THEN 3 END FROM omni_kw_176
< ok 1 0
> exec
| DROP TABLE omni_kw_176 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_177 (a NUMBER, REVERSE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REVERSE FROM omni_kw_177
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_177 WHERE REVERSE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_177 WHERE REVERSE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_177 ORDER BY REVERSE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_177 SET REVERSE = NULL WHERE REVERSE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_177 (REVERSE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_177_i ON omni_kw_177 (REVERSE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_177 WHERE REVERSE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REVERSE WHEN 1 THEN 2 END FROM omni_kw_177
< ok 1 0
> exec ^
: string
| SELECT CASE k.REVERSE WHEN 1 THEN 2 END FROM omni_kw_177 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REVERSE WHEN 2 THEN 3 END FROM omni_kw_177
< ok 1 0
> exec
| DROP TABLE omni_kw_177 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_178 (a NUMBER, REWRITE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT REWRITE FROM omni_kw_178
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_178 WHERE REWRITE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_178 WHERE REWRITE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_178 ORDER BY REWRITE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_178 SET REWRITE = NULL WHERE REWRITE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_178 (REWRITE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_178_i ON omni_kw_178 (REWRITE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_178 WHERE REWRITE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE REWRITE WHEN 1 THEN 2 END FROM omni_kw_178
< ok 1 0
> exec ^
: string
| SELECT CASE k.REWRITE WHEN 1 THEN 2 END FROM omni_kw_178 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + REWRITE WHEN 2 THEN 3 END FROM omni_kw_178
< ok 1 0
> exec
| DROP TABLE omni_kw_178 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_179 (a NUMBER, RIGHT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RIGHT FROM omni_kw_179
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_179 WHERE RIGHT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_179 WHERE RIGHT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_179 ORDER BY RIGHT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_179 SET RIGHT = NULL WHERE RIGHT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_179 (RIGHT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_179_i ON omni_kw_179 (RIGHT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_179 WHERE RIGHT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RIGHT WHEN 1 THEN 2 END FROM omni_kw_179
< ok 1 0
> exec ^
: string
| SELECT CASE k.RIGHT WHEN 1 THEN 2 END FROM omni_kw_179 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RIGHT WHEN 2 THEN 3 END FROM omni_kw_179
< ok 1 0
> exec
| DROP TABLE omni_kw_179 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_180 (a NUMBER, ROLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ROLE FROM omni_kw_180
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_180 WHERE ROLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_180 WHERE ROLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_180 ORDER BY ROLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_180 SET ROLE = NULL WHERE ROLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_180 (ROLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_180_i ON omni_kw_180 (ROLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_180 WHERE ROLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ROLE WHEN 1 THEN 2 END FROM omni_kw_180
< ok 1 0
> exec ^
: string
| SELECT CASE k.ROLE WHEN 1 THEN 2 END FROM omni_kw_180 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ROLE WHEN 2 THEN 3 END FROM omni_kw_180
< ok 1 0
> exec
| DROP TABLE omni_kw_180 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_181 (a NUMBER, ROLLBACK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ROLLBACK FROM omni_kw_181
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_181 WHERE ROLLBACK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_181 WHERE ROLLBACK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_181 ORDER BY ROLLBACK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_181 SET ROLLBACK = NULL WHERE ROLLBACK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_181 (ROLLBACK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_181_i ON omni_kw_181 (ROLLBACK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_181 WHERE ROLLBACK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ROLLBACK WHEN 1 THEN 2 END FROM omni_kw_181
< ok 1 0
> exec ^
: string
| SELECT CASE k.ROLLBACK WHEN 1 THEN 2 END FROM omni_kw_181 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ROLLBACK WHEN 2 THEN 3 END FROM omni_kw_181
< ok 1 0
> exec
| DROP TABLE omni_kw_181 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_182 (a NUMBER, ROLLUP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ROLLUP FROM omni_kw_182
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_182 WHERE ROLLUP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_182 WHERE ROLLUP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_182 ORDER BY ROLLUP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_182 SET ROLLUP = NULL WHERE ROLLUP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_182 (ROLLUP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_182_i ON omni_kw_182 (ROLLUP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_182 WHERE ROLLUP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ROLLUP WHEN 1 THEN 2 END FROM omni_kw_182
< ok 1 0
> exec ^
: string
| SELECT CASE k.ROLLUP WHEN 1 THEN 2 END FROM omni_kw_182 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ROLLUP WHEN 2 THEN 3 END FROM omni_kw_182
< ok 1 0
> exec
| DROP TABLE omni_kw_182 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_183 (a NUMBER, ROWTYPE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ROWTYPE FROM omni_kw_183
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_183 WHERE ROWTYPE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_183 WHERE ROWTYPE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_183 ORDER BY ROWTYPE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_183 SET ROWTYPE = NULL WHERE ROWTYPE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_183 (ROWTYPE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_183_i ON omni_kw_183 (ROWTYPE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_183 WHERE ROWTYPE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ROWTYPE WHEN 1 THEN 2 END FROM omni_kw_183
< ok 1 0
> exec ^
: string
| SELECT CASE k.ROWTYPE WHEN 1 THEN 2 END FROM omni_kw_183 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ROWTYPE WHEN 2 THEN 3 END FROM omni_kw_183
< ok 1 0
> exec
| DROP TABLE omni_kw_183 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_184 (a NUMBER, RULES NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT RULES FROM omni_kw_184
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_184 WHERE RULES IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_184 WHERE RULES NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_184 ORDER BY RULES
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_184 SET RULES = NULL WHERE RULES IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_184 (RULES) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_184_i ON omni_kw_184 (RULES)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_184 WHERE RULES(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE RULES WHEN 1 THEN 2 END FROM omni_kw_184
< ok 1 0
> exec ^
: string
| SELECT CASE k.RULES WHEN 1 THEN 2 END FROM omni_kw_184 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + RULES WHEN 2 THEN 3 END FROM omni_kw_184
< ok 1 0
> exec
| DROP TABLE omni_kw_184 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_185 (a NUMBER, SAMPLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SAMPLE FROM omni_kw_185
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_185 WHERE SAMPLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_185 WHERE SAMPLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_185 ORDER BY SAMPLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_185 SET SAMPLE = NULL WHERE SAMPLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_185 (SAMPLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_185_i ON omni_kw_185 (SAMPLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_185 WHERE SAMPLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SAMPLE WHEN 1 THEN 2 END FROM omni_kw_185
< ok 1 0
> exec ^
: string
| SELECT CASE k.SAMPLE WHEN 1 THEN 2 END FROM omni_kw_185 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SAMPLE WHEN 2 THEN 3 END FROM omni_kw_185
< ok 1 0
> exec
| DROP TABLE omni_kw_185 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_186 (a NUMBER, SAVEPOINT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SAVEPOINT FROM omni_kw_186
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_186 WHERE SAVEPOINT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_186 WHERE SAVEPOINT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_186 ORDER BY SAVEPOINT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_186 SET SAVEPOINT = NULL WHERE SAVEPOINT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_186 (SAVEPOINT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_186_i ON omni_kw_186 (SAVEPOINT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_186 WHERE SAVEPOINT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SAVEPOINT WHEN 1 THEN 2 END FROM omni_kw_186
< ok 1 0
> exec ^
: string
| SELECT CASE k.SAVEPOINT WHEN 1 THEN 2 END FROM omni_kw_186 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SAVEPOINT WHEN 2 THEN 3 END FROM omni_kw_186
< ok 1 0
> exec
| DROP TABLE omni_kw_186 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_187 (a NUMBER, SCHEMA NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SCHEMA FROM omni_kw_187
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_187 WHERE SCHEMA IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_187 WHERE SCHEMA NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_187 ORDER BY SCHEMA
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_187 SET SCHEMA = NULL WHERE SCHEMA IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_187 (SCHEMA) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_187_i ON omni_kw_187 (SCHEMA)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_187 WHERE SCHEMA(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SCHEMA WHEN 1 THEN 2 END FROM omni_kw_187
< ok 1 0
> exec ^
: string
| SELECT CASE k.SCHEMA WHEN 1 THEN 2 END FROM omni_kw_187 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SCHEMA WHEN 2 THEN 3 END FROM omni_kw_187
< ok 1 0
> exec
| DROP TABLE omni_kw_187 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_188 (a NUMBER, SCN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SCN FROM omni_kw_188
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_188 WHERE SCN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_188 WHERE SCN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_188 ORDER BY SCN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_188 SET SCN = NULL WHERE SCN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_188 (SCN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_188_i ON omni_kw_188 (SCN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_188 WHERE SCN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SCN WHEN 1 THEN 2 END FROM omni_kw_188
< ok 1 0
> exec ^
: string
| SELECT CASE k.SCN WHEN 1 THEN 2 END FROM omni_kw_188 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SCN WHEN 2 THEN 3 END FROM omni_kw_188
< ok 1 0
> exec
| DROP TABLE omni_kw_188 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_189 (a NUMBER, SEED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SEED FROM omni_kw_189
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_189 WHERE SEED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_189 WHERE SEED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_189 ORDER BY SEED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_189 SET SEED = NULL WHERE SEED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_189 (SEED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_189_i ON omni_kw_189 (SEED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_189 WHERE SEED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SEED WHEN 1 THEN 2 END FROM omni_kw_189
< ok 1 0
> exec ^
: string
| SELECT CASE k.SEED WHEN 1 THEN 2 END FROM omni_kw_189 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SEED WHEN 2 THEN 3 END FROM omni_kw_189
< ok 1 0
> exec
| DROP TABLE omni_kw_189 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_190 (a NUMBER, SEQUENCE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SEQUENCE FROM omni_kw_190
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_190 WHERE SEQUENCE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_190 WHERE SEQUENCE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_190 ORDER BY SEQUENCE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_190 SET SEQUENCE = NULL WHERE SEQUENCE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_190 (SEQUENCE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_190_i ON omni_kw_190 (SEQUENCE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_190 WHERE SEQUENCE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SEQUENCE WHEN 1 THEN 2 END FROM omni_kw_190
< ok 1 0
> exec ^
: string
| SELECT CASE k.SEQUENCE WHEN 1 THEN 2 END FROM omni_kw_190 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SEQUENCE WHEN 2 THEN 3 END FROM omni_kw_190
< ok 1 0
> exec
| DROP TABLE omni_kw_190 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_191 (a NUMBER, SEQUENTIAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SEQUENTIAL FROM omni_kw_191
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_191 WHERE SEQUENTIAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_191 WHERE SEQUENTIAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_191 ORDER BY SEQUENTIAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_191 SET SEQUENTIAL = NULL WHERE SEQUENTIAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_191 (SEQUENTIAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_191_i ON omni_kw_191 (SEQUENTIAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_191 WHERE SEQUENTIAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SEQUENTIAL WHEN 1 THEN 2 END FROM omni_kw_191
< ok 1 0
> exec ^
: string
| SELECT CASE k.SEQUENTIAL WHEN 1 THEN 2 END FROM omni_kw_191 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SEQUENTIAL WHEN 2 THEN 3 END FROM omni_kw_191
< ok 1 0
> exec
| DROP TABLE omni_kw_191 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_192 (a NUMBER, SETS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SETS FROM omni_kw_192
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_192 WHERE SETS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_192 WHERE SETS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_192 ORDER BY SETS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_192 SET SETS = NULL WHERE SETS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_192 (SETS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_192_i ON omni_kw_192 (SETS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_192 WHERE SETS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SETS WHEN 1 THEN 2 END FROM omni_kw_192
< ok 1 0
> exec ^
: string
| SELECT CASE k.SETS WHEN 1 THEN 2 END FROM omni_kw_192 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SETS WHEN 2 THEN 3 END FROM omni_kw_192
< ok 1 0
> exec
| DROP TABLE omni_kw_192 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_193 (a NUMBER, SKIP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SKIP FROM omni_kw_193
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_193 WHERE SKIP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_193 WHERE SKIP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_193 ORDER BY SKIP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_193 SET SKIP = NULL WHERE SKIP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_193 (SKIP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_193_i ON omni_kw_193 (SKIP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_193 WHERE SKIP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SKIP WHEN 1 THEN 2 END FROM omni_kw_193
< ok 1 0
> exec ^
: string
| SELECT CASE k.SKIP WHEN 1 THEN 2 END FROM omni_kw_193 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SKIP WHEN 2 THEN 3 END FROM omni_kw_193
< ok 1 0
> exec
| DROP TABLE omni_kw_193 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_194 (a NUMBER, STATISTICS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT STATISTICS FROM omni_kw_194
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_194 WHERE STATISTICS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_194 WHERE STATISTICS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_194 ORDER BY STATISTICS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_194 SET STATISTICS = NULL WHERE STATISTICS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_194 (STATISTICS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_194_i ON omni_kw_194 (STATISTICS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_194 WHERE STATISTICS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE STATISTICS WHEN 1 THEN 2 END FROM omni_kw_194
< ok 1 0
> exec ^
: string
| SELECT CASE k.STATISTICS WHEN 1 THEN 2 END FROM omni_kw_194 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + STATISTICS WHEN 2 THEN 3 END FROM omni_kw_194
< ok 1 0
> exec
| DROP TABLE omni_kw_194 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_195 (a NUMBER, STORAGE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT STORAGE FROM omni_kw_195
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_195 WHERE STORAGE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_195 WHERE STORAGE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_195 ORDER BY STORAGE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_195 SET STORAGE = NULL WHERE STORAGE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_195 (STORAGE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_195_i ON omni_kw_195 (STORAGE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_195 WHERE STORAGE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE STORAGE WHEN 1 THEN 2 END FROM omni_kw_195
< ok 1 0
> exec ^
: string
| SELECT CASE k.STORAGE WHEN 1 THEN 2 END FROM omni_kw_195 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + STORAGE WHEN 2 THEN 3 END FROM omni_kw_195
< ok 1 0
> exec
| DROP TABLE omni_kw_195 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_196 (a NUMBER, SUBPARTITION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SUBPARTITION FROM omni_kw_196
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_196 WHERE SUBPARTITION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_196 WHERE SUBPARTITION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_196 ORDER BY SUBPARTITION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_196 SET SUBPARTITION = NULL WHERE SUBPARTITION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_196 (SUBPARTITION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_196_i ON omni_kw_196 (SUBPARTITION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_196 WHERE SUBPARTITION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SUBPARTITION WHEN 1 THEN 2 END FROM omni_kw_196
< ok 1 0
> exec ^
: string
| SELECT CASE k.SUBPARTITION WHEN 1 THEN 2 END FROM omni_kw_196 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SUBPARTITION WHEN 2 THEN 3 END FROM omni_kw_196
< ok 1 0
> exec
| DROP TABLE omni_kw_196 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_197 (a NUMBER, SUBTYPE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SUBTYPE FROM omni_kw_197
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_197 WHERE SUBTYPE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_197 WHERE SUBTYPE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_197 ORDER BY SUBTYPE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_197 SET SUBTYPE = NULL WHERE SUBTYPE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_197 (SUBTYPE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_197_i ON omni_kw_197 (SUBTYPE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_197 WHERE SUBTYPE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SUBTYPE WHEN 1 THEN 2 END FROM omni_kw_197
< ok 1 0
> exec ^
: string
| SELECT CASE k.SUBTYPE WHEN 1 THEN 2 END FROM omni_kw_197 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SUBTYPE WHEN 2 THEN 3 END FROM omni_kw_197
< ok 1 0
> exec
| DROP TABLE omni_kw_197 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_198 (a NUMBER, SYSTEM NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SYSTEM FROM omni_kw_198
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_198 WHERE SYSTEM IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_198 WHERE SYSTEM NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_198 ORDER BY SYSTEM
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_198 SET SYSTEM = NULL WHERE SYSTEM IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_198 (SYSTEM) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_198_i ON omni_kw_198 (SYSTEM)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_198 WHERE SYSTEM(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SYSTEM WHEN 1 THEN 2 END FROM omni_kw_198
< ok 1 0
> exec ^
: string
| SELECT CASE k.SYSTEM WHEN 1 THEN 2 END FROM omni_kw_198 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SYSTEM WHEN 2 THEN 3 END FROM omni_kw_198
< ok 1 0
> exec
| DROP TABLE omni_kw_198 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_199 (a NUMBER, SYS_CONNECT_BY_PATH NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT SYS_CONNECT_BY_PATH FROM omni_kw_199
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_199 WHERE SYS_CONNECT_BY_PATH IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_199 WHERE SYS_CONNECT_BY_PATH NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_199 ORDER BY SYS_CONNECT_BY_PATH
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_199 SET SYS_CONNECT_BY_PATH = NULL WHERE SYS_CONNECT_BY_PATH IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_199 (SYS_CONNECT_BY_PATH) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_199_i ON omni_kw_199 (SYS_CONNECT_BY_PATH)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_199 WHERE SYS_CONNECT_BY_PATH(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE SYS_CONNECT_BY_PATH WHEN 1 THEN 2 END FROM omni_kw_199
< ok 1 0
> exec ^
: string
| SELECT CASE k.SYS_CONNECT_BY_PATH WHEN 1 THEN 2 END FROM omni_kw_199 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + SYS_CONNECT_BY_PATH WHEN 2 THEN 3 END FROM omni_kw_199
< ok 1 0
> exec
| DROP TABLE omni_kw_199 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_200 (a NUMBER, TABLESPACE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TABLESPACE FROM omni_kw_200
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_200 WHERE TABLESPACE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_200 WHERE TABLESPACE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_200 ORDER BY TABLESPACE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_200 SET TABLESPACE = NULL WHERE TABLESPACE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_200 (TABLESPACE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_200_i ON omni_kw_200 (TABLESPACE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_200 WHERE TABLESPACE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TABLESPACE WHEN 1 THEN 2 END FROM omni_kw_200
< ok 1 0
> exec ^
: string
| SELECT CASE k.TABLESPACE WHEN 1 THEN 2 END FROM omni_kw_200 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TABLESPACE WHEN 2 THEN 3 END FROM omni_kw_200
< ok 1 0
> exec
| DROP TABLE omni_kw_200 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_201 (a NUMBER, TEMPORARY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TEMPORARY FROM omni_kw_201
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_201 WHERE TEMPORARY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_201 WHERE TEMPORARY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_201 ORDER BY TEMPORARY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_201 SET TEMPORARY = NULL WHERE TEMPORARY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_201 (TEMPORARY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_201_i ON omni_kw_201 (TEMPORARY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_201 WHERE TEMPORARY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TEMPORARY WHEN 1 THEN 2 END FROM omni_kw_201
< ok 1 0
> exec ^
: string
| SELECT CASE k.TEMPORARY WHEN 1 THEN 2 END FROM omni_kw_201 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TEMPORARY WHEN 2 THEN 3 END FROM omni_kw_201
< ok 1 0
> exec
| DROP TABLE omni_kw_201 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_202 (a NUMBER, TIES NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TIES FROM omni_kw_202
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_202 WHERE TIES IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_202 WHERE TIES NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_202 ORDER BY TIES
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_202 SET TIES = NULL WHERE TIES IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_202 (TIES) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_202_i ON omni_kw_202 (TIES)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_202 WHERE TIES(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TIES WHEN 1 THEN 2 END FROM omni_kw_202
< ok 1 0
> exec ^
: string
| SELECT CASE k.TIES WHEN 1 THEN 2 END FROM omni_kw_202 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TIES WHEN 2 THEN 3 END FROM omni_kw_202
< ok 1 0
> exec
| DROP TABLE omni_kw_202 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_203 (a NUMBER, TIMESTAMP NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TIMESTAMP FROM omni_kw_203
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_203 WHERE TIMESTAMP IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_203 WHERE TIMESTAMP NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_203 ORDER BY TIMESTAMP
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_203 SET TIMESTAMP = NULL WHERE TIMESTAMP IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_203 (TIMESTAMP) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_203_i ON omni_kw_203 (TIMESTAMP)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_203 WHERE TIMESTAMP(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TIMESTAMP WHEN 1 THEN 2 END FROM omni_kw_203
< ok 1 0
> exec ^
: string
| SELECT CASE k.TIMESTAMP WHEN 1 THEN 2 END FROM omni_kw_203 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TIMESTAMP WHEN 2 THEN 3 END FROM omni_kw_203
< ok 1 0
> exec
| DROP TABLE omni_kw_203 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_204 (a NUMBER, TRANSACTION NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TRANSACTION FROM omni_kw_204
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_204 WHERE TRANSACTION IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_204 WHERE TRANSACTION NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_204 ORDER BY TRANSACTION
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_204 SET TRANSACTION = NULL WHERE TRANSACTION IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_204 (TRANSACTION) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_204_i ON omni_kw_204 (TRANSACTION)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_204 WHERE TRANSACTION(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TRANSACTION WHEN 1 THEN 2 END FROM omni_kw_204
< ok 1 0
> exec ^
: string
| SELECT CASE k.TRANSACTION WHEN 1 THEN 2 END FROM omni_kw_204 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TRANSACTION WHEN 2 THEN 3 END FROM omni_kw_204
< ok 1 0
> exec
| DROP TABLE omni_kw_204 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_205 (a NUMBER, TREAT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TREAT FROM omni_kw_205
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_205 WHERE TREAT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_205 WHERE TREAT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_205 ORDER BY TREAT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_205 SET TREAT = NULL WHERE TREAT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_205 (TREAT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_205_i ON omni_kw_205 (TREAT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_205 WHERE TREAT(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE TREAT WHEN 1 THEN 2 END FROM omni_kw_205
< ok 1 0
> exec ^
: string
| SELECT CASE k.TREAT WHEN 1 THEN 2 END FROM omni_kw_205 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TREAT WHEN 2 THEN 3 END FROM omni_kw_205
< ok 1 0
> exec
| DROP TABLE omni_kw_205 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_206 (a NUMBER, TRUNCATE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TRUNCATE FROM omni_kw_206
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_206 WHERE TRUNCATE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_206 WHERE TRUNCATE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_206 ORDER BY TRUNCATE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_206 SET TRUNCATE = NULL WHERE TRUNCATE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_206 (TRUNCATE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_206_i ON omni_kw_206 (TRUNCATE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_206 WHERE TRUNCATE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TRUNCATE WHEN 1 THEN 2 END FROM omni_kw_206
< ok 1 0
> exec ^
: string
| SELECT CASE k.TRUNCATE WHEN 1 THEN 2 END FROM omni_kw_206 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TRUNCATE WHEN 2 THEN 3 END FROM omni_kw_206
< ok 1 0
> exec
| DROP TABLE omni_kw_206 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_207 (a NUMBER, TYPE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT TYPE FROM omni_kw_207
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_207 WHERE TYPE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_207 WHERE TYPE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_207 ORDER BY TYPE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_207 SET TYPE = NULL WHERE TYPE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_207 (TYPE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_207_i ON omni_kw_207 (TYPE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_207 WHERE TYPE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE TYPE WHEN 1 THEN 2 END FROM omni_kw_207
< ok 1 0
> exec ^
: string
| SELECT CASE k.TYPE WHEN 1 THEN 2 END FROM omni_kw_207 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + TYPE WHEN 2 THEN 3 END FROM omni_kw_207
< ok 1 0
> exec
| DROP TABLE omni_kw_207 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_208 (a NUMBER, UNBOUNDED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UNBOUNDED FROM omni_kw_208
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_208 WHERE UNBOUNDED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_208 WHERE UNBOUNDED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_208 ORDER BY UNBOUNDED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_208 SET UNBOUNDED = NULL WHERE UNBOUNDED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_208 (UNBOUNDED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_208_i ON omni_kw_208 (UNBOUNDED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_208 WHERE UNBOUNDED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UNBOUNDED WHEN 1 THEN 2 END FROM omni_kw_208
< ok 1 0
> exec ^
: string
| SELECT CASE k.UNBOUNDED WHEN 1 THEN 2 END FROM omni_kw_208 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UNBOUNDED WHEN 2 THEN 3 END FROM omni_kw_208
< ok 1 0
> exec
| DROP TABLE omni_kw_208 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_209 (a NUMBER, UNDER NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UNDER FROM omni_kw_209
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_209 WHERE UNDER IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_209 WHERE UNDER NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_209 ORDER BY UNDER
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_209 SET UNDER = NULL WHERE UNDER IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_209 (UNDER) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_209_i ON omni_kw_209 (UNDER)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_209 WHERE UNDER(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UNDER WHEN 1 THEN 2 END FROM omni_kw_209
< ok 1 0
> exec ^
: string
| SELECT CASE k.UNDER WHEN 1 THEN 2 END FROM omni_kw_209 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UNDER WHEN 2 THEN 3 END FROM omni_kw_209
< ok 1 0
> exec
| DROP TABLE omni_kw_209 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_210 (a NUMBER, UNPIVOT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UNPIVOT FROM omni_kw_210
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_210 WHERE UNPIVOT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_210 WHERE UNPIVOT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_210 ORDER BY UNPIVOT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_210 SET UNPIVOT = NULL WHERE UNPIVOT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_210 (UNPIVOT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_210_i ON omni_kw_210 (UNPIVOT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_210 WHERE UNPIVOT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UNPIVOT WHEN 1 THEN 2 END FROM omni_kw_210
< ok 1 0
> exec ^
: string
| SELECT CASE k.UNPIVOT WHEN 1 THEN 2 END FROM omni_kw_210 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UNPIVOT WHEN 2 THEN 3 END FROM omni_kw_210
< ok 1 0
> exec
| DROP TABLE omni_kw_210 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_211 (a NUMBER, UNTIL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UNTIL FROM omni_kw_211
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_211 WHERE UNTIL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_211 WHERE UNTIL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_211 ORDER BY UNTIL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_211 SET UNTIL = NULL WHERE UNTIL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_211 (UNTIL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_211_i ON omni_kw_211 (UNTIL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_211 WHERE UNTIL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UNTIL WHEN 1 THEN 2 END FROM omni_kw_211
< ok 1 0
> exec ^
: string
| SELECT CASE k.UNTIL WHEN 1 THEN 2 END FROM omni_kw_211 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UNTIL WHEN 2 THEN 3 END FROM omni_kw_211
< ok 1 0
> exec
| DROP TABLE omni_kw_211 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_212 (a NUMBER, UPDATED NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UPDATED FROM omni_kw_212
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_212 WHERE UPDATED IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_212 WHERE UPDATED NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_212 ORDER BY UPDATED
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_212 SET UPDATED = NULL WHERE UPDATED IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_212 (UPDATED) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_212_i ON omni_kw_212 (UPDATED)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_212 WHERE UPDATED(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UPDATED WHEN 1 THEN 2 END FROM omni_kw_212
< ok 1 0
> exec ^
: string
| SELECT CASE k.UPDATED WHEN 1 THEN 2 END FROM omni_kw_212 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UPDATED WHEN 2 THEN 3 END FROM omni_kw_212
< ok 1 0
> exec
| DROP TABLE omni_kw_212 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_213 (a NUMBER, UPSERT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT UPSERT FROM omni_kw_213
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_213 WHERE UPSERT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_213 WHERE UPSERT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_213 ORDER BY UPSERT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_213 SET UPSERT = NULL WHERE UPSERT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_213 (UPSERT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_213_i ON omni_kw_213 (UPSERT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_213 WHERE UPSERT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE UPSERT WHEN 1 THEN 2 END FROM omni_kw_213
< ok 1 0
> exec ^
: string
| SELECT CASE k.UPSERT WHEN 1 THEN 2 END FROM omni_kw_213 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + UPSERT WHEN 2 THEN 3 END FROM omni_kw_213
< ok 1 0
> exec
| DROP TABLE omni_kw_213 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_214 (a NUMBER, USING NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT USING FROM omni_kw_214
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_214 WHERE USING IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_214 WHERE USING NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_214 ORDER BY USING
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_214 SET USING = NULL WHERE USING IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_214 (USING) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_214_i ON omni_kw_214 (USING)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_214 WHERE USING(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE USING WHEN 1 THEN 2 END FROM omni_kw_214
< ok 1 0
> exec ^
: string
| SELECT CASE k.USING WHEN 1 THEN 2 END FROM omni_kw_214 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + USING WHEN 2 THEN 3 END FROM omni_kw_214
< ok 1 0
> exec
| DROP TABLE omni_kw_214 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_215 (a NUMBER, VARRAY NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT VARRAY FROM omni_kw_215
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_215 WHERE VARRAY IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_215 WHERE VARRAY NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_215 ORDER BY VARRAY
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_215 SET VARRAY = NULL WHERE VARRAY IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_215 (VARRAY) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_215_i ON omni_kw_215 (VARRAY)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_215 WHERE VARRAY(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE VARRAY WHEN 1 THEN 2 END FROM omni_kw_215
< ok 1 0
> exec ^
: string
| SELECT CASE k.VARRAY WHEN 1 THEN 2 END FROM omni_kw_215 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + VARRAY WHEN 2 THEN 3 END FROM omni_kw_215
< ok 1 0
> exec
| DROP TABLE omni_kw_215 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_216 (a NUMBER, VERSIONS NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT VERSIONS FROM omni_kw_216
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_216 WHERE VERSIONS IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_216 WHERE VERSIONS NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_216 ORDER BY VERSIONS
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_216 SET VERSIONS = NULL WHERE VERSIONS IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_216 (VERSIONS) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_216_i ON omni_kw_216 (VERSIONS)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_216 WHERE VERSIONS(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE VERSIONS WHEN 1 THEN 2 END FROM omni_kw_216
< ok 1 0
> exec ^
: string
| SELECT CASE k.VERSIONS WHEN 1 THEN 2 END FROM omni_kw_216 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + VERSIONS WHEN 2 THEN 3 END FROM omni_kw_216
< ok 1 0
> exec
| DROP TABLE omni_kw_216 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_217 (a NUMBER, VIRTUAL NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT VIRTUAL FROM omni_kw_217
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_217 WHERE VIRTUAL IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_217 WHERE VIRTUAL NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_217 ORDER BY VIRTUAL
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_217 SET VIRTUAL = NULL WHERE VIRTUAL IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_217 (VIRTUAL) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_217_i ON omni_kw_217 (VIRTUAL)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_217 WHERE VIRTUAL(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE VIRTUAL WHEN 1 THEN 2 END FROM omni_kw_217
< ok 1 0
> exec ^
: string
| SELECT CASE k.VIRTUAL WHEN 1 THEN 2 END FROM omni_kw_217 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + VIRTUAL WHEN 2 THEN 3 END FROM omni_kw_217
< ok 1 0
> exec
| DROP TABLE omni_kw_217 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_218 (a NUMBER, WAIT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WAIT FROM omni_kw_218
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_218 WHERE WAIT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_218 WHERE WAIT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_218 ORDER BY WAIT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_218 SET WAIT = NULL WHERE WAIT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_218 (WAIT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_218_i ON omni_kw_218 (WAIT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_218 WHERE WAIT(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WAIT WHEN 1 THEN 2 END FROM omni_kw_218
< ok 1 0
> exec ^
: string
| SELECT CASE k.WAIT WHEN 1 THEN 2 END FROM omni_kw_218 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WAIT WHEN 2 THEN 3 END FROM omni_kw_218
< ok 1 0
> exec
| DROP TABLE omni_kw_218 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_219 (a NUMBER, WHEN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WHEN FROM omni_kw_219
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_219 WHERE WHEN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_219 WHERE WHEN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_219 ORDER BY WHEN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_219 SET WHEN = NULL WHERE WHEN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_219 (WHEN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_219_i ON omni_kw_219 (WHEN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_219 WHERE WHEN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WHEN WHEN 1 THEN 2 END FROM omni_kw_219
< error
| 2000 ORA-02000: missing THEN keyword
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE k.WHEN WHEN 1 THEN 2 END FROM omni_kw_219 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WHEN WHEN 2 THEN 3 END FROM omni_kw_219
< ok 1 0
> exec
| DROP TABLE omni_kw_219 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_220 (a NUMBER, WHILE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WHILE FROM omni_kw_220
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_220 WHERE WHILE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_220 WHERE WHILE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_220 ORDER BY WHILE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_220 SET WHILE = NULL WHERE WHILE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_220 (WHILE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_220_i ON omni_kw_220 (WHILE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_220 WHERE WHILE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WHILE WHEN 1 THEN 2 END FROM omni_kw_220
< ok 1 0
> exec ^
: string
| SELECT CASE k.WHILE WHEN 1 THEN 2 END FROM omni_kw_220 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WHILE WHEN 2 THEN 3 END FROM omni_kw_220
< ok 1 0
> exec
| DROP TABLE omni_kw_220 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_221 (a NUMBER, WITHIN NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WITHIN FROM omni_kw_221
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_221 WHERE WITHIN IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_221 WHERE WITHIN NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_221 ORDER BY WITHIN
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_221 SET WITHIN = NULL WHERE WITHIN IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_221 (WITHIN) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_221_i ON omni_kw_221 (WITHIN)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_221 WHERE WITHIN(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WITHIN WHEN 1 THEN 2 END FROM omni_kw_221
< ok 1 0
> exec ^
: string
| SELECT CASE k.WITHIN WHEN 1 THEN 2 END FROM omni_kw_221 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WITHIN WHEN 2 THEN 3 END FROM omni_kw_221
< ok 1 0
> exec
| DROP TABLE omni_kw_221 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_222 (a NUMBER, WORK NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WORK FROM omni_kw_222
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_222 WHERE WORK IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_222 WHERE WORK NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_222 ORDER BY WORK
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_222 SET WORK = NULL WHERE WORK IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_222 (WORK) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_222_i ON omni_kw_222 (WORK)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_222 WHERE WORK(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WORK WHEN 1 THEN 2 END FROM omni_kw_222
< ok 1 0
> exec ^
: string
| SELECT CASE k.WORK WHEN 1 THEN 2 END FROM omni_kw_222 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WORK WHEN 2 THEN 3 END FROM omni_kw_222
< ok 1 0
> exec
| DROP TABLE omni_kw_222 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_223 (a NUMBER, WRITE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT WRITE FROM omni_kw_223
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_223 WHERE WRITE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_223 WHERE WRITE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_223 ORDER BY WRITE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_223 SET WRITE = NULL WHERE WRITE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_223 (WRITE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_223_i ON omni_kw_223 (WRITE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_223 WHERE WRITE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE WRITE WHEN 1 THEN 2 END FROM omni_kw_223
< ok 1 0
> exec ^
: string
| SELECT CASE k.WRITE WHEN 1 THEN 2 END FROM omni_kw_223 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + WRITE WHEN 2 THEN 3 END FROM omni_kw_223
< ok 1 0
> exec
| DROP TABLE omni_kw_223 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_224 (a NUMBER, XMLAGG NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLAGG FROM omni_kw_224
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_224 WHERE XMLAGG IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_224 WHERE XMLAGG NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_224 ORDER BY XMLAGG
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_224 SET XMLAGG = NULL WHERE XMLAGG IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_224 (XMLAGG) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_224_i ON omni_kw_224 (XMLAGG)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_224 WHERE XMLAGG(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE XMLAGG WHEN 1 THEN 2 END FROM omni_kw_224
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLAGG WHEN 1 THEN 2 END FROM omni_kw_224 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLAGG WHEN 2 THEN 3 END FROM omni_kw_224
< ok 1 0
> exec
| DROP TABLE omni_kw_224 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_225 (a NUMBER, XMLELEMENT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLELEMENT FROM omni_kw_225
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_225 WHERE XMLELEMENT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_225 WHERE XMLELEMENT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_225 ORDER BY XMLELEMENT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_225 SET XMLELEMENT = NULL WHERE XMLELEMENT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_225 (XMLELEMENT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_225_i ON omni_kw_225 (XMLELEMENT)
< error
| 906 ORA-00906: missing left parenthesis
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_225 WHERE XMLELEMENT(+) = LENGTH(d.dummy)
< error
| 931 ORA-00931: missing identifier
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE XMLELEMENT WHEN 1 THEN 2 END FROM omni_kw_225
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLELEMENT WHEN 1 THEN 2 END FROM omni_kw_225 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLELEMENT WHEN 2 THEN 3 END FROM omni_kw_225
< ok 1 0
> exec
| DROP TABLE omni_kw_225 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_226 (a NUMBER, XMLFOREST NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLFOREST FROM omni_kw_226
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_226 WHERE XMLFOREST IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_226 WHERE XMLFOREST NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_226 ORDER BY XMLFOREST
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_226 SET XMLFOREST = NULL WHERE XMLFOREST IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_226 (XMLFOREST) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_226_i ON omni_kw_226 (XMLFOREST)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_226 WHERE XMLFOREST(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE XMLFOREST WHEN 1 THEN 2 END FROM omni_kw_226
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLFOREST WHEN 1 THEN 2 END FROM omni_kw_226 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLFOREST WHEN 2 THEN 3 END FROM omni_kw_226
< ok 1 0
> exec
| DROP TABLE omni_kw_226 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_227 (a NUMBER, XMLPARSE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLPARSE FROM omni_kw_227
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_227 WHERE XMLPARSE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_227 WHERE XMLPARSE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_227 ORDER BY XMLPARSE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_227 SET XMLPARSE = NULL WHERE XMLPARSE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_227 (XMLPARSE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_227_i ON omni_kw_227 (XMLPARSE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_227 WHERE XMLPARSE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE XMLPARSE WHEN 1 THEN 2 END FROM omni_kw_227
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLPARSE WHEN 1 THEN 2 END FROM omni_kw_227 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLPARSE WHEN 2 THEN 3 END FROM omni_kw_227
< ok 1 0
> exec
| DROP TABLE omni_kw_227 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_228 (a NUMBER, XMLROOT NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLROOT FROM omni_kw_228
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_228 WHERE XMLROOT IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_228 WHERE XMLROOT NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_228 ORDER BY XMLROOT
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_228 SET XMLROOT = NULL WHERE XMLROOT IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_228 (XMLROOT) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_228_i ON omni_kw_228 (XMLROOT)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_228 WHERE XMLROOT(+) = LENGTH(d.dummy)
< error
| 936 ORA-00936: missing expression
| ORA-06512: at line 14
| ORA-06512: at "SYS.DBMS_SQL", line 1251
| ORA-06512: at line 7
|
> exec ^
: string
| SELECT CASE XMLROOT WHEN 1 THEN 2 END FROM omni_kw_228
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLROOT WHEN 1 THEN 2 END FROM omni_kw_228 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLROOT WHEN 2 THEN 3 END FROM omni_kw_228
< ok 1 0
> exec
| DROP TABLE omni_kw_228 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_229 (a NUMBER, XMLSERIALIZE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLSERIALIZE FROM omni_kw_229
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_229 WHERE XMLSERIALIZE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_229 WHERE XMLSERIALIZE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_229 ORDER BY XMLSERIALIZE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_229 SET XMLSERIALIZE = NULL WHERE XMLSERIALIZE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_229 (XMLSERIALIZE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_229_i ON omni_kw_229 (XMLSERIALIZE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_229 WHERE XMLSERIALIZE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE XMLSERIALIZE WHEN 1 THEN 2 END FROM omni_kw_229
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLSERIALIZE WHEN 1 THEN 2 END FROM omni_kw_229 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLSERIALIZE WHEN 2 THEN 3 END FROM omni_kw_229
< ok 1 0
> exec
| DROP TABLE omni_kw_229 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_230 (a NUMBER, XMLTABLE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT XMLTABLE FROM omni_kw_230
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_230 WHERE XMLTABLE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_230 WHERE XMLTABLE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_230 ORDER BY XMLTABLE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_230 SET XMLTABLE = NULL WHERE XMLTABLE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_230 (XMLTABLE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_230_i ON omni_kw_230 (XMLTABLE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_230 WHERE XMLTABLE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE XMLTABLE WHEN 1 THEN 2 END FROM omni_kw_230
< ok 1 0
> exec ^
: string
| SELECT CASE k.XMLTABLE WHEN 1 THEN 2 END FROM omni_kw_230 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + XMLTABLE WHEN 2 THEN 3 END FROM omni_kw_230
< ok 1 0
> exec
| DROP TABLE omni_kw_230 PURGE
< ok 0 0
> exec
| CREATE TABLE omni_kw_231 (a NUMBER, ZONE NUMBER)
< ok 0 0
> exec
|
| DECLARE
|   c INTEGER;
| BEGIN
|   c := DBMS_SQL.OPEN_CURSOR;
|   BEGIN
|     DBMS_SQL.PARSE(c, :1, DBMS_SQL.NATIVE);
|     DBMS_SQL.CLOSE_CURSOR(c);
|   EXCEPTION
|     WHEN OTHERS THEN
|       IF DBMS_SQL.IS_OPEN(c) THEN
|         DBMS_SQL.CLOSE_CURSOR(c);
|       END IF;
|       RAISE;
|   END;
| END;
: string
| SELECT ZONE FROM omni_kw_231
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_231 WHERE ZONE IS NULL
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_231 WHERE ZONE NOT IN (1, 2)
< ok 1 0
> exec ^
: string
| SELECT a FROM omni_kw_231 ORDER BY ZONE
< ok 1 0
> exec ^
: string
| UPDATE omni_kw_231 SET ZONE = NULL WHERE ZONE IS NULL
< ok 1 0
> exec ^
: string
| INSERT INTO omni_kw_231 (ZONE) VALUES (1)
< ok 1 0
> exec ^
: string
| CREATE INDEX omni_kw_231_i ON omni_kw_231 (ZONE)
< ok 1 0
> exec ^
: string
| SELECT d.dummy FROM dual d, omni_kw_231 WHERE ZONE(+) = LENGTH(d.dummy)
< ok 1 0
> exec ^
: string
| SELECT CASE ZONE WHEN 1 THEN 2 END FROM omni_kw_231
< ok 1 0
> exec ^
: string
| SELECT CASE k.ZONE WHEN 1 THEN 2 END FROM omni_kw_231 k
< ok 1 0
> exec ^
: string
| SELECT CASE 1 + ZONE WHEN 2 THEN 3 END FROM omni_kw_231
< ok 1 0
> exec
| DROP TABLE omni_kw_231 PURGE
< ok 0 0
