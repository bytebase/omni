-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestParseRefusalFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE TRIGGER tr BEFORE INSERT ON src FOR EACH ROW UPDATE old.t SET old.t.id = 1
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON src FOR EACH ROW SET @a = (SELECT old.t.id FROM old.t LIMIT 1)
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON src FOR EACH ROW SET @a = old.id
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON src FOR EACH ROW SET @a = old.id
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TABLE t (d BLOB DEFAULT (COLUMN_CREATE('a', 1)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d BLOB, e BLOB AS (COLUMN_ADD(d, 'a', 1)) VIRTUAL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d BLOB CHECK (COLUMN_CHECK(d)))
< ok 0 0
> exec ^
: string
| DELETE x.t FROM y.t
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE x.t FROM y.t
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE d.a FROM t AS a
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE d.a FROM t AS a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE a FROM d.t AS a
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| DELETE a FROM d.t AS a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (WEEK(d,3) > 0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (week(d,0) > 0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (WEEK(d) > 0))
< error
| 1901 HY000 Function or expression 'week()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (WEEK(d) > 0))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'week()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (DATE_FORMAT(d,'%a','en_US') <> ''))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (DATE_FORMAT(d,'%a') <> ''))
< error
| 1901 HY000 Function or expression 'date_format()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (DATE_FORMAT(d,'%a') <> ''))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'date_format()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (UNIX_TIMESTAMP(d) > 0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (UNIX_TIMESTAMP() > 0))
< error
| 1901 HY000 Function or expression 'unix_timestamp()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE CHECK (UNIX_TIMESTAMP() > 0))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'unix_timestamp()' cannot be used in the CHECK clause of `d`
> exec ^
: string
| CREATE TABLE t (d DATE, CHECK (WEEK(d) > 0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (d DATE DEFAULT (WEEK(d)))
< ok 0 0
> exec ^
: string
| SELECT SETVAL(s, -9223372036854775807)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT SETVAL(s, -9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SETVAL(s, -9223372036854775808)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT SETVAL(s, -9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SETVAL(s, -09223372036854775808)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT SETVAL(s, -09223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SETVAL(s, -9223372036854775809)
< error
| 1064 42000 %s value is out of range in '%s' near '9223372036854775809)' at line 1
> exec ^
: string
| SELECT SETVAL(s, 9223372036854775807)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT SETVAL(s, 9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN RETURN 1; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN RETURN 1; END
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE PROCEDURE p() BEGIN RETURN 1; END
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE PROCEDURE p() RETURN 1
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE PROCEDURE p() RETURN 1
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE TRIGGER g BEFORE INSERT ON t FOR EACH ROW RETURN 1
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE TRIGGER g BEFORE INSERT ON t FOR EACH ROW RETURN 1
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO RETURN 1
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO RETURN 1
| ; SELECT FROM
< error
| 1313 42000 RETURN is only allowed in a FUNCTION
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN IF 1 THEN RETURN 1; END IF; RETURN 2; END
< ok 0 0
> exec ^
: string
| GRANT `r(1)` TO u
< ok 0 0
> exec ^
: string
| REVOKE `r(1)` FROM u
< ok 0 0
> exec ^
: string
| GRANT 'r(1)' TO u
< ok 0 0
> exec ^
: string
| GRANT r(a) TO u
< error
| 1064 42000 …near '(a) TO u' at line 1
> exec ^
: string
| ALTER EVENT b.e RENAME TO a.f
< ok 0 0
> exec ^
: string
| ALTER EVENT b.e RENAME TO f
< ok 0 0
> exec ^
: string
| CREATE TABLE t1 (a POINT, x DOUBLE DEFAULT x(a), y DOUBLE DEFAULT y(a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t2 (a POINT, b DOUBLE AS (ST_X(a)) VIRTUAL, c DOUBLE AS (st_y(a)) PERSISTENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t3 (a POINT, b TEXT DEFAULT (ST_AsText(a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t4 (a INT, b INT DEFAULT nosuchfn(a))
< error
| 1901 HY000 Function or expression '`nosuchfn`()' cannot be used in the DEFAULT clause of `b`
> exec ^
: string
| CREATE TABLE t4 (a INT, b INT DEFAULT nosuchfn(a))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`nosuchfn`()' cannot be used in the DEFAULT clause of `b`
> exec ^
: string
| EXPLAIN FORMAT=TRADITIONAL SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT=JSON SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT=TREE SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT=TREE SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT='JSON' SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT="json" SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT=`json` SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT=foo SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| EXPLAIN FORMAT=foo SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| EXPLAIN FORMAT='' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: ''
> exec ^
: string
| EXPLAIN FORMAT='' SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: ''
> exec ^
: string
| EXPLAIN ANALYZE FORMAT=TREE SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT=TREE SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FORMAT=JSON SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT=JSON SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FORMAT=TRADITIONAL SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT=TRADITIONAL SELECT 1' at line 1
> exec ^
: string
| EXPLAIN ANALYZE FORMAT=foo SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT=foo SELECT 1' at line 1
> exec ^
: string
| EXPLAIN FORMAT=JSON FOR CONNECTION 1
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT=foo FOR CONNECTION 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| EXPLAIN FORMAT=foo FOR CONNECTION 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| ANALYZE FORMAT=JSON SELECT 1
< ok 0 0
> exec ^
: string
| ANALYZE FORMAT=foo SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| ANALYZE FORMAT=foo SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| DESC FORMAT=foo SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| DESC FORMAT=foo SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| DESCRIBE FORMAT=JSON SELECT 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=INPLACE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=`inplace`
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM='inplace'
< error
| 1064 42000 …near ''inplace'' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM=COPY
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=INSTANT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=NOCOPY
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=`default`
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALGORITHM=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| ALTER TABLE t ALGORITHM=select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| ALTER TABLE t LOCK=NONE
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK=`none`
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK='none'
< error
| 1064 42000 …near ''none'' at line 1
> exec ^
: string
| ALTER TABLE t LOCK=SHARED
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK=EXCLUSIVE
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK=DEFAULT
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK=`default`
< ok 0 0
> exec ^
: string
| ALTER TABLE t LOCK=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) ALGORITHM=NOCOPY
< ok 0 0
> exec ^
: string
| DROP INDEX i ON t ALGORITHM=NOCOPY
< error
| 1064 42000 …near 'ALGORITHM=NOCOPY' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) LOCK=`shared`
< ok 0 0
> exec ^
: string
| DROP INDEX i ON t LOCK=DEFAULT
< error
| 1064 42000 …near 'LOCK=DEFAULT' at line 1
> exec ^
: string
| OPTIMIZE TABLE t
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 CHARACTER SET utf8mb4
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET 'utf8mb4' CHARSET 'latin1'
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET 'utf8mb4' CHARSET 'latin1'
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE latin1_bin
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| CREATE DATABASE d COLLATE latin1_bin CHARACTER SET utf8mb4
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET binary COLLATE binary
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET DEFAULT CHARACTER SET latin1
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE DEFAULT COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE DEFAULT' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE DATABASE d COLLATE DEFAULT COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE DEFAULT' and 'COLLATE latin1_bin'
> exec ^
: string
| ALTER DATABASE d CHARACTER SET utf8mb4 CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER DATABASE d CHARACTER SET utf8mb4 CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER DATABASE d CHARACTER SET latin1 COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| ALTER DATABASE d COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| ALTER DATABASE d COLLATE utf8mb4_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 CHARACTER SET latin1)
< error
| 1064 42000 …near 'CHARACTER SET latin1)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET latin1 COLLATE utf8mb4_bin)
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE latin1_bin COLLATE utf8mb4_bin)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE latin1_bin' and 'COLLATE utf8mb4_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE latin1_bin COLLATE utf8mb4_bin)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE latin1_bin' and 'COLLATE utf8mb4_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin COLLATE ascii_general_ci)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE ascii_general_ci'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin COLLATE ascii_general_ci)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE ascii_general_ci'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COLLATE ascii_bin)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE ascii_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COLLATE ascii_bin)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE ascii_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin PRIMARY KEY COLLATE latin1_bin)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin PRIMARY KEY COLLATE latin1_bin)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a INT COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5)) CHARACTER SET utf8mb4 CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE TABLE t (a CHAR(5)) CHARACTER SET utf8mb4 CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE TABLE t (a CHAR(5)) COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5)) COLLATE utf8mb4_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a CHAR(5)) CHARACTER SET latin1 COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| ALTER TABLE t CHARACTER SET utf8mb4 CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER TABLE t CHARACTER SET utf8mb4 CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET latin1 COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a CHAR(5) COLLATE ascii_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a CHAR(5) COLLATE ascii_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) PARTITIONS 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) PARTITIONS 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 0 (PARTITION p VALUES LESS THAN (1))
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 0 (PARTITION p VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) PARTITIONS 0 (PARTITION p VALUES LESS THAN (1))
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) PARTITIONS 0 (PARTITION p VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH (a) PARTITIONS 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH (a) PARTITIONS 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 0
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a)
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a)
| ; SELECT FROM
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p)
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p)
< error
| 1479 HY000 Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES LESS THAN (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES IN (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES IN (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION p VALUES IN (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION p VALUES IN (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES LESS THAN (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES IN (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES IN (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1,2))
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1,2))
| ; SELECT FROM
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (1,2,3))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (1,2,3))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN MAXVALUE)
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near 'MAXVALUE)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN (1,2))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN (1,2))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1,2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1,2))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2),(3,4,5)))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2),(3,4,5)))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (10,(20,30)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (10,(20,30)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((10,20),30))
< error
| 1064 42000 …near '30))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2 (PARTITION p)
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1 (PARTITION p, PARTITION q)
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s), PARTITION q VALUES LESS THAN (2))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p (SUBPARTITION s))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY HASH (a) SUBPARTITION BY HASH (b) (PARTITION p)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=3 (a)
< error
| 1064 42000 …near '3 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=0 (a)
< error
| 1064 42000 …near '0 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=2 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p DEFAULT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p DEFAULT)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p DEFAULT)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME INTERVAL 1 DAY (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) WITH SYSTEM VERSIONING PARTITION BY SYSTEM_TIME LIMIT 5 (PARTITION x HISTORY, PARTITION y CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y)
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y)
| ; SELECT FROM
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES LESS THAN (1), PARTITION y CURRENT)
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES LESS THAN (1), PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES IN (1), PARTITION y CURRENT)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x VALUES IN (1), PARTITION y CURRENT)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION x HISTORY, PARTITION y CURRENT)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION x HISTORY, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION x CURRENT)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `RANGE`
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION x CURRENT)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `RANGE`
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) (PARTITION p)
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE (a) (PARTITION p)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1,2))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION x HISTORY)
< ok 0 0
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p VALUES IN ((1,2),3))
< error
| 1064 42000 …near '3))' at line 1
> exec ^
: string
| SELECT `a `.t.b FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `a `.t.b FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT `a `.t FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `a `.t FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT 1 FROM `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT 1 FROM ``.t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| SELECT 1 FROM ``.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| SELECT 1 FROM ``
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| SELECT 1 FROM ``
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| SELECT 1 FROM `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT 1 FROM `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT 1 FROM t AS `a `
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t AS `a `
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t AS ``
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t AS ``
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP TABLE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP TABLE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP TABLE ``
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| DROP TABLE ``
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| CREATE TABLE `a `.t (x INT)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE `a `.t (x INT)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE `a ` (x INT)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE `a ` (x INT)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE `` (x INT)
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| CREATE TABLE `` (x INT)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| ALTER TABLE `a `.t ADD COLUMN y INT
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE `a `.t ADD COLUMN y INT
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE `a ` ADD COLUMN y INT
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER TABLE `a ` ADD COLUMN y INT
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE VIEW `a `.v AS SELECT 1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE VIEW `a `.v AS SELECT 1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE VIEW `a ` AS SELECT 1
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE VIEW `a ` AS SELECT 1
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE INDEX i ON `a `.t (x)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE INDEX i ON `a `.t (x)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE INDEX i ON `a ` (x)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE INDEX i ON `a ` (x)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE INDEX `a ` ON t (x)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (`a ` INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (`` INT)
< ok 0 0
> exec ^
: string
| CREATE DATABASE `a `
< ok 0 0
> exec ^
: string
| CREATE DATABASE ``
< ok 0 0
> exec ^
: string
| DROP DATABASE `a `
< ok 0 0
> exec ^
: string
| USE `a `
< ok 0 0
> exec ^
: string
| INSERT INTO `a `.t VALUES (1)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| INSERT INTO `a `.t VALUES (1)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| INSERT INTO `a ` VALUES (1)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| INSERT INTO `a ` VALUES (1)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| UPDATE `a `.t SET a=1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| UPDATE `a `.t SET a=1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| UPDATE `a ` SET a=1
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| UPDATE `a ` SET a=1
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DELETE FROM `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DELETE FROM `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DELETE FROM `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DELETE FROM `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| TRUNCATE TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| TRUNCATE TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| RENAME TABLE `a `.t TO u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| RENAME TABLE `a `.t TO u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| RENAME TABLE t TO `a `.u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| RENAME TABLE t TO `a `.u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT 1 FROM t, t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, u, t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, u, t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t a, t a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, t a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, u a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, u a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, t b
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t a, t b
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM d.t, e.t
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM d.t, e.t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM d.t, d.t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM d.t, d.t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, d.t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t, d.t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t JOIN t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t JOIN t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t JOIN u JOIN t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t JOIN u JOIN t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t a JOIN u a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a JOIN u a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t, (u, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, (u, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM (t, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM (t, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, (SELECT 1) t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t, (SELECT 1) t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a, (SELECT 2) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a, (SELECT 2) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t T, t t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t T, t t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t LEFT JOIN t ON 1
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t LEFT JOIN t ON 1
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT 1 FROM t, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT 1 FROM t, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| UPDATE t, t SET a=1
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| UPDATE t, t SET a=1
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE t FROM t, t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE t FROM t, t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE t, t FROM t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE t, t FROM t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE u FROM t
< error
| 1109 42S02 Unknown table 'u' in MULTI DELETE
> exec ^
: string
| DELETE u FROM t
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'u' in MULTI DELETE
> exec ^
: string
| DELETE t1 FROM t2
< error
| 1109 42S02 Unknown table 't1' in MULTI DELETE
> exec ^
: string
| DELETE t1 FROM t2
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't1' in MULTI DELETE
> exec ^
: string
| DELETE a.t1, b.t1 FROM a.t1, b.t1
< error
| 1146 42S02 Table 'a.t1' doesn't exist
> exec ^
: string
| DELETE a.t1, b.t1 FROM a.t1, b.t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE tablesample FROM t
< error
| 1109 42S02 Unknown table 'tablesample' in MULTI DELETE
> exec ^
: string
| DELETE tablesample FROM t
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'tablesample' in MULTI DELETE
> exec ^
: string
| DELETE t.* FROM t tablesample
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE t.* FROM t tablesample
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| SELECT * FROM t1 tablesample, t2 tablesample
< error
| 1066 42000 Not unique table/alias: 'tablesample'
> exec ^
: string
| SELECT * FROM t1 tablesample, t2 tablesample
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'tablesample'
> exec ^
: string
| INSERT INTO t SELECT 1 FROM u, u
< error
| 1066 42000 Not unique table/alias: 'u'
> exec ^
: string
| INSERT INTO t SELECT 1 FROM u, u
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'u'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t LIMIT 1, 10)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t LIMIT 1, 10)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT a FROM t LIMIT 1))
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT a FROM t LIMIT 1))
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a NOT IN (SELECT a FROM t LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a NOT IN (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT a FROM t LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ALL (SELECT a FROM t LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ALL (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT a FROM t LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = (SELECT a FROM t LIMIT 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT a FROM t LIMIT 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT a FROM t LIMIT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a IN (SELECT a FROM t LIMIT 1) FROM t
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT a IN (SELECT a FROM t LIMIT 1) FROM t
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t ORDER BY a LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT a FROM t ORDER BY a LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT * FROM (SELECT a FROM t LIMIT 1) d)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT * FROM (SELECT a FROM t LIMIT 1) d)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM (SELECT 1) x
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1) AS x
< ok 0 0
> exec ^
: string
| SELECT 1 FROM ((SELECT 1))
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM ((SELECT 1)) x
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1 UNION SELECT 2)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM (VALUES ROW(1))
< error
| 1064 42000 …near 'ROW(1))' at line 1
> exec ^
: string
| SELECT 1 FROM (VALUES ROW(1)) x
< error
| 1064 42000 …near 'ROW(1)) x' at line 1
> exec ^
: string
| SELECT 1 FROM (TABLE t)
< error
| 1064 42000 …near 'TABLE t)' at line 1
> exec ^
: string
| SELECT 1 FROM t JOIN (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM LATERAL (SELECT 1)
< error
| 1064 42000 …near '(SELECT 1)' at line 1
> exec ^
: string
| SELECT 1 FROM LATERAL (SELECT 1) x
< error
| 1064 42000 …near '(SELECT 1) x' at line 1
> exec ^
: string
| UPDATE (SELECT 1) SET a=1
< error
| 1064 42000 …near 'SET a=1' at line 1
> exec ^
: string
| DELETE FROM t WHERE a IN (SELECT 1 FROM (SELECT 1))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a YEAR(2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(4))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN a YEAR(2)
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a YEAR(2)) RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS YEAR(3) RETURN 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(255))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TINYINT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a TINYINT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a BIGINT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a BIGINT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a BIT(65))
< error
| 1439 42000 Display width out of range for 'a' (max = 64)
> exec ^
: string
| CREATE TABLE t (a BIT(65))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 64)
> exec ^
: string
| CREATE TABLE t (a BIT(64))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(256))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a CHAR(256))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a CHAR(255))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(256))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a BINARY(256))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VARCHAR(65535))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(16384))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(4294967296))
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(65536))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(4294967296))
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(4294967295))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(4294967296))
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(66))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a DECIMAL(66))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65,30))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65,31))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65,38))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65,39))
< error
| 1425 42000 Too big scale specified for 'a'. Maximum is 38
> exec ^
: string
| CREATE TABLE t (a DECIMAL(65,39))
| ; SELECT FROM
< error
| 1425 42000 Too big scale specified for 'a'. Maximum is 38
> exec ^
: string
| CREATE TABLE t (a DECIMAL(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a DECIMAL(5,6))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a NUMERIC(66))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a NUMERIC(66))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a FLOAT(54))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(54))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(53))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(256))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(256))
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,31))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,31))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a FLOAT(256,2))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a FLOAT(256,2))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a FLOAT(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a FLOAT(5,6))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a DOUBLE(256,2))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DOUBLE(256,2))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DOUBLE(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a DOUBLE(5,6))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,31))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,31))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a')
> exec ^
: string
| CREATE TABLE t (a TIME(7))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a TIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a TIME(6))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(7))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a DATETIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(7))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(7))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a VARCHAR(5) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VARCHAR(5) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DATE AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DATE AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TEXT AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TEXT AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT AUTO_INCREMENT KEY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT ON UPDATE CURRENT_TIMESTAMP)
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a INT ON UPDATE CURRENT_TIMESTAMP)
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATE ON UPDATE CURRENT_TIMESTAMP)
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATE ON UPDATE CURRENT_TIMESTAMP)
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VARCHAR(5) ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a VARCHAR(5) ON UPDATE NOW())
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW())
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) ON UPDATE NOW(2))
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) ON UPDATE NOW(2))
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) ON UPDATE NOW(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE NOW(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NOW())
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(3) DEFAULT NOW(2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL DEFAULT NULL)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT PRIMARY KEY DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT PRIMARY KEY DEFAULT NULL)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT KEY DEFAULT 1)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT KEY DEFAULT 1)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a TEXT DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON DEFAULT '{}')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a GEOMETRY DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(1) DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) DEFAULT 1)
< error
| 1064 42000 …near 'DEFAULT 1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) AUTO_INCREMENT)
< error
| 1064 42000 …near 'AUTO_INCREMENT)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) ON UPDATE NOW())
< error
| 1064 42000 …near 'ON UPDATE NOW())' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT DEFAULT 1 AS (a))
< error
| 1064 42000 …near 'AS (a))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) NOT NULL)
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) PRIMARY KEY)
< error
| 1064 42000 …near 'PRIMARY KEY)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT AS (a) UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a ENUM())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a ENUM('a','a'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SET('a','a'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SET('a,b'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a ENUM(''))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET nosuch)
< error
| 1115 42000 Unknown character set: 'nosuch'
> exec ^
: string
| CREATE TABLE t (a INT UNSIGNED SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT NULL NOT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1 DEFAULT 2)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT COMMENT 'x' COMMENT 'y')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, a INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY k (a), KEY k (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT PRIMARY KEY PRIMARY KEY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT PRIMARY KEY, PRIMARY KEY (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, PRIMARY KEY (a), PRIMARY KEY (a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT UNIQUE UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 1
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 1
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 10
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5 MAXVALUE 5
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5 MAXVALUE 5
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5 MAXVALUE 4
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5 MAXVALUE 4
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 5 MAXVALUE 6
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s START WITH 0
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 0
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s START WITH 5 MAXVALUE 4
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 5 MAXVALUE 4
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 5 MINVALUE 6
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 5 MINVALUE 6
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH -1
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH -1
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH -1 MINVALUE -5
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 0
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -1 MAXVALUE 5
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -1 MINVALUE -5 MAXVALUE -1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -1 START WITH -1 MINVALUE -1 MAXVALUE -1000
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -1 START WITH -1 MINVALUE -1 MAXVALUE -1000
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -9223372036854775807
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -9223372036854775807
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -9223372036854775808
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT -9223372036854775808
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 9223372036854775807
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 9223372036854775807
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE SEQUENCE s CACHE 0
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s CACHE 1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s CACHE 1000000
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s CACHE 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE SEQUENCE s NOCACHE
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 9223372036854775806
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE -9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE -9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s START WITH 9223372036854775807
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 9223372036854775807
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s START WITH 9223372036854775806
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s NOMAXVALUE NOMINVALUE
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 1 MINVALUE 2
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 1 MINVALUE 2
| ; SELECT FROM
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 5 NOMAXVALUE
< error
| 1225 HY000 Option 'MAXVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 5 NOMAXVALUE
| ; SELECT FROM
< error
| 1225 HY000 Option 'MAXVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s CYCLE NOCYCLE
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s CYCLE NOCYCLE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s AS INT
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT MAXVALUE 128
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT MAXVALUE 127
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT UNSIGNED MAXVALUE 256
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT UNSIGNED MINVALUE -1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT UNSIGNED MAXVALUE 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT UNSIGNED START WITH 18446744073709551614
< ok 0 0
> exec ^
: string
| ALTER SEQUENCE s MAXVALUE 0
< ok 0 0
> exec ^
: string
| ALTER SEQUENCE s MINVALUE 5 MAXVALUE 4
< ok 0 0
> exec ^
: string
| ALTER SEQUENCE s INCREMENT 0
< ok 0 0
> exec ^
: string
| ALTER SEQUENCE s RESTART WITH 0
< ok 0 0
> exec ^
: string
| SELECT DISTINCT ALL 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT DISTINCT ALL 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT ALL DISTINCT 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT ALL DISTINCT 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT DISTINCTROW ALL 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT DISTINCTROW ALL 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT ALL DISTINCTROW 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT ALL DISTINCTROW 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT ALL ALL 1
< ok 0 0
> exec ^
: string
| SELECT DISTINCT DISTINCT 1
< ok 0 0
> exec ^
: string
| SELECT DISTINCT DISTINCTROW 1
< ok 0 0
> exec ^
: string
| SELECT DISTINCT HIGH_PRIORITY ALL 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT DISTINCT HIGH_PRIORITY ALL 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| SELECT SQL_CACHE SQL_NO_CACHE 1
< error
| 1221 HY000 Incorrect usage of SQL_NO_CACHE and SQL_CACHE
> exec ^
: string
| SELECT SQL_CACHE SQL_NO_CACHE 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of SQL_NO_CACHE and SQL_CACHE
> exec ^
: string
| SELECT SQL_NO_CACHE SQL_CACHE 1
< error
| 1221 HY000 Incorrect usage of SQL_NO_CACHE and SQL_CACHE
> exec ^
: string
| SELECT SQL_NO_CACHE SQL_CACHE 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of SQL_NO_CACHE and SQL_CACHE
> exec ^
: string
| SELECT SQL_CACHE SQL_CACHE 1
< error
| 1225 HY000 Option 'SQL_CACHE' used twice in statement
> exec ^
: string
| SELECT SQL_CACHE SQL_CACHE 1
| ; SELECT FROM
< error
| 1225 HY000 Option 'SQL_CACHE' used twice in statement
> exec ^
: string
| SELECT SQL_NO_CACHE SQL_NO_CACHE 1
< error
| 1225 HY000 Option 'SQL_NO_CACHE' used twice in statement
> exec ^
: string
| SELECT SQL_NO_CACHE SQL_NO_CACHE 1
| ; SELECT FROM
< error
| 1225 HY000 Option 'SQL_NO_CACHE' used twice in statement
> exec ^
: string
| SELECT HIGH_PRIORITY HIGH_PRIORITY 1
< ok 0 0
> exec ^
: string
| SELECT STRAIGHT_JOIN STRAIGHT_JOIN 1
< ok 0 0
> exec ^
: string
| SELECT SQL_SMALL_RESULT SQL_BIG_RESULT 1
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 0 DAY
< error
| 1525 HY000 Incorrect DAY value: '0'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 0 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '0'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 65535 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 65536 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 99999999999999999999 DAY
< error
| 1064 42000 …near '99999999999999999999 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 65535
< error
| 1064 42000 …near 'HISTORY 65535' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 65536
< error
| 1064 42000 …near 'HISTORY 65536' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 65535 DAY
< error
| 1064 42000 …near 'REUSE INTERVAL 65535 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 65536 DAY
< error
| 1064 42000 …near 'REUSE INTERVAL 65536 DAY' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 32767
< error
| 1064 42000 …near 'FAILED_LOGIN_ATTEMPTS 32767' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 32768
< error
| 1064 42000 …near 'FAILED_LOGIN_ATTEMPTS 32768' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 32767
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME 32767' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 32768
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME 32768' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 4294967296
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD EXPIRE INTERVAL 65536 DAY
< ok 0 0
> exec ^
: string
| ALTER USER u FAILED_LOGIN_ATTEMPTS 32768
< error
| 1064 42000 …near 'FAILED_LOGIN_ATTEMPTS 32768' at line 1
> exec ^
: string
| ALTER USER u PASSWORD_LOCK_TIME 32768
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME 32768' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='x'
< error
| 1064 42000 …near 'ENCRYPTION='x'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='Y'
< error
| 1064 42000 …near 'ENCRYPTION='Y'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='n'
< error
| 1064 42000 …near 'ENCRYPTION='n'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION=''
< error
| 1064 42000 …near 'ENCRYPTION=''' at line 1
> exec ^
: string
| ALTER DATABASE d ENCRYPTION='x'
< error
| 1064 42000 …near 'ENCRYPTION='x'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTION='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTION='Y'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COMPRESSION='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=FOO
< error
| 1064 42000 …near 'FOO' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=-1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC=2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=0
< error
| 1064 42000 …near '0' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=65536
< error
| 1064 42000 …near '65536' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM=2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM=2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=-1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE=nosuch
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE nosuch
< error
| 1273 HY000 Unknown collation: 'nosuch'
> exec ^
: string
| SELECT ABS()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT ABS()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT ABS(1,2)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT ABS(1,2)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT NOW(7)
< error
| 1426 42000 Too big precision specified for 'current_timestamp'. Maximum is 6
> exec ^
: string
| SELECT NOW(7)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT NOW(6)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(7)
< error
| 1426 42000 Too big precision specified for 'curtime'. Maximum is 6
> exec ^
: string
| SELECT CURTIME(7)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(7)
< error
| 1426 42000 Too big precision specified for 'current_timestamp'. Maximum is 6
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(7)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SYSDATE(7)
< error
| 1426 42000 Too big precision specified for 'sysdate'. Maximum is 6
> exec ^
: string
| SELECT SYSDATE(7)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT UTC_TIME(7)
< error
| 1426 42000 Too big precision specified for 'utc_time'. Maximum is 6
> exec ^
: string
| SELECT UTC_TIME(7)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(1 AS TIME(7))
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 6
> exec ^
: string
| SELECT CAST(1 AS TIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 6
> exec ^
: string
| SELECT 1 LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY 'aa'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ESCAPED BY 'aa'
< ok 0 0
> exec ^
: string
| LOAD DATA INFILE '/x' INTO TABLE t FIELDS ENCLOSED BY 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE '/x' INTO TABLE t FIELDS ENCLOSED BY 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE '/x' INTO TABLE t FIELDS ESCAPED BY 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE '/x' INTO TABLE t FIELDS ESCAPED BY 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 LIKE 'a' ESCAPE 'xy'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| SELECT 1 LIKE 'a' ESCAPE 'xy'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 LIKE 'a' ESCAPE ''
< ok 0 0
> exec ^
: string
| GRANT SELECT ON *.* TO u WITH MAX_QUERIES_PER_HOUR 99999999999999999999
< ok 0 0
> exec ^
: string
| SET @@nosuch = 1
< error
| 1193 HY000 Unknown system variable 'nosuch'
> exec ^
: string
| SET SESSION nosuch = 1
< error
| 1193 HY000 Unknown system variable 'nosuch'
> exec ^
: string
| SET GLOBAL sort_buffer_size = 1, sql_mode = ''
< ok 0 0
> exec ^
: string
| SET NAMES nosuch
< error
| 1115 42000 Unknown character set: 'nosuch'
> exec ^
: string
| SET CHARACTER SET nosuch
< error
| 1115 42000 Unknown character set: 'nosuch'
> exec ^
: string
| SHOW COLLATION LIKE 'a' WHERE 1
< error
| 1064 42000 …near 'WHERE 1' at line 1
> exec ^
: string
| HELP 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| KILL 99999999999999999999
< ok 0 0
> exec ^
: string
| SHOW PROCEDURE CODE pr
< error
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
> exec ^
: string
| SHOW FUNCTION CODE fn
< error
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) UNION (SELECT 2)
< ok 0 0
> exec ^
: string
| CREATE TABLE t SELECT 1 FOR UPDATE
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) STORAGE ENGINE text)
< error
| 1286 42000 Unknown storage engine 'text'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1) ENGINE=nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec ^
: string
| create resource group x ru_per_sec=1 ru_per_sec=2
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 ru_per_sec=2' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1, ru_per_sec=2
< error
| 1064 42000 …near 'resource group x ru_per_sec=1, ru_per_sec=2' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 priority=low priority=high
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 priority=low priority=high' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill action=dryrun)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill acti...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s', exec_elapsed='2s' action=kill)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s', exec_elapsed='2...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='foo' action=kill)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 query_limit=(exec_elapsed='foo' action=kill)' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill watch similar duration 'foo')
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill watc...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1 background=(task_types='a' task_types='b')
< error
| 1064 42000 …near 'resource group x ru_per_sec=1 background=(task_types='a' task_types='b')' at line 1
> exec ^
: string
| create resource group x burstable=18446744073709551615
< error
| 1064 42000 …near 'resource group x burstable=18446744073709551615' at line 1
> exec ^
: string
| create resource group x burstable=9223372036854775808
< error
| 1064 42000 …near 'resource group x burstable=9223372036854775808' at line 1
> exec ^
: string
| create resource group x burstable=9223372036854775807
< error
| 1064 42000 …near 'resource group x burstable=9223372036854775807' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615 action=kill)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615 actio...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808 action=kill)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808 action...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808 action=kill)
< error
| 1064 42000 …near 'resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 92233720368547...' at line 1
> exec ^
: string
| create resource group x ru_per_sec=18446744073709551615
< error
| 1064 42000 …near 'resource group x ru_per_sec=18446744073709551615' at line 1
> exec ^
: string
| query watch remove 9223372036854775808
< error
| 1064 42000 …near 'query watch remove 9223372036854775808' at line 1
> exec ^
: string
| query watch remove 18446744073709551615
< error
| 1064 42000 …near 'query watch remove 18446744073709551615' at line 1
> exec ^
: string
| query watch remove 9223372036854775807
< error
| 1064 42000 …near 'query watch remove 9223372036854775807' at line 1
> exec ^
: string
| query watch add action kill action dryrun sql text exact to 'x'
< error
| 1064 42000 …near 'query watch add action kill action dryrun sql text exact to 'x'' at line 1
> exec ^
: string
| query watch add resource group a resource group b sql text exact to 'x'
< error
| 1064 42000 …near 'query watch add resource group a resource group b sql text exact to 'x'' at line 1
> exec ^
: string
| create placement policy x followers=0
< error
| 1064 42000 …near 'placement policy x followers=0' at line 1
> exec ^
: string
| create placement policy x followers=1
< error
| 1064 42000 …near 'placement policy x followers=1' at line 1
> exec ^
: string
| alter placement policy x followers=0
< error
| 1064 42000 …near 'placement policy x followers=0' at line 1
> exec ^
: string
| create table t (a int) followers=0
< ok 0 0
> exec ^
: string
| create database d followers=0
< error
| 1064 42000 …near 'followers=0' at line 1
> exec ^
: string
| alter table t algorithm=foo
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| alter table t algorithm=foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| alter table t algorithm foo
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| alter table t algorithm foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| alter table t algorithm default
< ok 0 0
> exec ^
: string
| alter table t lock=foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock=foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock default
< ok 0 0
> exec ^
: string
| create index i on t(a) algorithm=foo
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| create index i on t(a) algorithm=foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'foo'
> exec ^
: string
| create index i on t(a) algorithm=`inplace`
< ok 0 0
> exec ^
: string
| create index i on t(a) lock=foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| create index i on t(a) lock=foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| create index i on t(a) lock=`none`
< ok 0 0
> exec ^
: string
| drop index i on t algorithm=foo
< error
| 1064 42000 …near 'algorithm=foo' at line 1
> exec ^
: string
| drop index i on t lock=foo
< error
| 1064 42000 …near 'lock=foo' at line 1
> exec ^
: string
| drop index i on t algorithm=inplace
< error
| 1064 42000 …near 'algorithm=inplace' at line 1
> exec ^
: string
| drop index i on t lock=none
< error
| 1064 42000 …near 'lock=none' at line 1
> exec ^
: string
| flashback cluster to tso 0
< error
| 1064 42000 …near 'flashback cluster to tso 0' at line 1
> exec ^
: string
| flashback cluster to tso 1
< error
| 1064 42000 …near 'flashback cluster to tso 1' at line 1
> exec ^
: string
| flashback table t to tso 0
< error
| 1064 42000 …near 'flashback table t to tso 0' at line 1
> exec ^
: string
| flashback table t to tso 1
< error
| 1064 42000 …near 'flashback table t to tso 1' at line 1
> exec ^
: string
| flashback database d to tso 0
< error
| 1064 42000 …near 'flashback database d to tso 0' at line 1
> exec ^
: string
| flashback schema d to tso 0
< error
| 1064 42000 …near 'flashback schema d to tso 0' at line 1
> exec ^
: string
| flashback database d to tso 1
< error
| 1064 42000 …near 'flashback database d to tso 1' at line 1
> exec ^
: string
| create table t (a int, b timestamp as (a) on update current_timestamp)
< error
| 1064 42000 …near 'on update current_timestamp)' at line 1
> exec ^
: string
| create table t (a int, b int as (a) comment 'c')
< ok 0 0
> exec ^
: string
| create table t (a int, b int auto_increment as (a))
< error
| 1064 42000 …near 'as (a))' at line 1
> exec ^
: string
| alter table t add column b int as (a) default 1
< error
| 1064 42000 …near 'default 1' at line 1
> exec ^
: string
| alter table t add column (b int as (a) default 1)
< error
| 1064 42000 …near 'default 1)' at line 1
> exec ^
: string
| alter table t modify b int as (a) default 1
< error
| 1064 42000 …near 'default 1' at line 1
> exec ^
: string
| alter table t change a b int as (a) default 1
< error
| 1064 42000 …near 'default 1' at line 1
> exec ^
: string
| create global temporary table t (a int)
< error
| 1064 42000 …near 'global temporary table t (a int)' at line 1
> exec ^
: string
| create global temporary table t (a int) on commit delete rows
< error
| 1064 42000 …near 'global temporary table t (a int) on commit delete rows' at line 1
> exec ^
: string
| create global temporary table t (a int) on commit preserve rows
< error
| 1064 42000 …near 'global temporary table t (a int) on commit preserve rows' at line 1
> exec ^
: string
| create temporary table t (a int) on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec ^
: string
| create temporary table t (a int) on commit preserve rows
< error
| 1064 42000 …near 'on commit preserve rows' at line 1
> exec ^
: string
| create table t (a int) on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec ^
: string
| create table t (a int) on commit preserve rows
< error
| 1064 42000 …near 'on commit preserve rows' at line 1
> exec ^
: string
| create global temporary table t like u
< error
| 1064 42000 …near 'global temporary table t like u' at line 1
> exec ^
: string
| create global temporary table t like u on commit delete rows
< error
| 1064 42000 …near 'global temporary table t like u on commit delete rows' at line 1
> exec ^
: string
| create temporary table t like u on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec ^
: string
| create table t like u on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec ^
: string
| create global temporary table t select 1
< error
| 1064 42000 …near 'global temporary table t select 1' at line 1
> exec ^
: string
| create global temporary table t select 1 on commit delete rows
< error
| 1064 42000 …near 'global temporary table t select 1 on commit delete rows' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 1
< ok 0 0
> exec ^
: string
| create table t (a int) partition by key(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by key(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 0 (partition p values less than (1))
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 0 (partition p values less than (1))
| ; SELECT FROM
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (partition x)
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (partition x, partition y)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2) (subpartition b))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 2 (partition x values less than (1) (subpartition a))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near '))' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 1 (partition x values less than (1) (subpartition a))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x (subpartition a))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list(a)
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list(a)
| ; SELECT FROM
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range columns(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range columns(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list columns(a)
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list columns(a)
| ; SELECT FROM
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) partitions 3
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) partitions 3
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by hash(a)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) interval (1) first partition less than (1) last partition less than (5)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) interval (1) first partition less than (1) last partition less than (5)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) (partition x)
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x)
< error
| 1479 HY000 Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values in (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values in (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by key(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by key(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values in (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values in (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1, 2))
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1, 2))
| ; SELECT FROM
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than maxvalue)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (maxvalue))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (1))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (1, 2))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than maxvalue)
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near 'maxvalue)' at line 1
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (maxvalue, maxvalue))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in ((1, 2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (1, 2))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (1, 2))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4)))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), 3))
< error
| 1064 42000 …near '3))' at line 1
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4, 5)))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4, 5)))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), default))
< error
| 1064 42000 …near 'default))' at line 1
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x default)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x default)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (default))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (default, 1))
< error
| 1064 42000 …near ' 1))' at line 1
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (1, default))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x default)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x default)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x default)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x default)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| create table t (a int) partition by system_time interval 1 day (partition x history, partition y current)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x history, partition y current)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x history, partition y current)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| create table t (a int) partition by range(a) (partition x history, partition y current)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `RANGE`
> exec ^
: string
| create table t (a int) partition by range(a) (partition x history, partition y current)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `RANGE`
> exec ^
: string
| alter table t partition by range(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by hash(a)
< ok 0 0
> exec ^
: string
| alter table t partition by hash(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| alter table t partition by hash(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| alter table t partition by hash(a) partitions 2 (partition x)
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| alter table t partition by range(a) (partition x)
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| alter table t partition by range(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax error: RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| alter table t partition by list(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| alter table t partition by list(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| alter table t partition by linear key algorithm=3 (a)
< error
| 1064 42000 …near '3 (a)' at line 1
> exec ^
: string
| alter table t add partition (partition x)
< ok 0 0
> exec ^
: string
| alter table t add partition (partition x values less than (1, 2))
< ok 0 0
> exec ^
: string
| alter table t add partition (partition x values in (1) (subpartition a))
< ok 0 0
> exec ^
: string
| alter table t reorganize partition a into (partition x)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by key algorithm=1 (a)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by linear key algorithm=0 (a)
< error
| 1064 42000 …near '0 (a)' at line 1
> exec ^
: string
| create table t (a int) partition by linear key algorithm=3 (a) partitions 4
< error
| 1064 42000 …near '3 (a) partitions 4' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by key algorithm=3 (a) (partition x values less than (1))
< error
| 1064 42000 …near '3 (a) (partition x values less than (1))' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 day ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 day ago' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 month ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 month ago' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 year ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 year ago' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 quarter ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 quarter ago' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 week ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 week ago' at line 1
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 microsecond ago
< error
| 1064 42000 …near 'database * to 's3://a' snapshot = 1 microsecond ago' at line 1
> exec ^
: string
| select 1 like 'a' escape '||'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape '||'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 like 'a' escape 'x'
< ok 0 0
> exec ^
: string
| select 1 not like 'a' escape 'xy'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 not like 'a' escape 'xy'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select 1 ilike 'a' escape 'xy'
< error
| 1064 42000 …near ''a' escape 'xy'' at line 1
> exec ^
: string
| select 1 not ilike 'a' escape 'xy'
< error
| 1064 42000 …near 'ilike 'a' escape 'xy'' at line 1
> exec ^
: string
| select 1 ilike 'a' escape 'x'
< error
| 1064 42000 …near ''a' escape 'x'' at line 1
> exec ^
: string
| select 1 like 'a' escape '\\'
< ok 0 0
> exec ^
: string
| select 1 like 'a' escape '€'
< ok 0 0
> exec ^
: string
| select cast(1 as float(54))
< error
| 1064 42000 …near '(54))' at line 1
> exec ^
: string
| select cast(1 as float(53))
< error
| 1064 42000 …near '(53))' at line 1
> exec ^
: string
| select cast(1 as float(25))
< error
| 1064 42000 …near '(25))' at line 1
> exec ^
: string
| select cast(1 as float(255))
< error
| 1064 42000 …near '(255))' at line 1
> exec ^
: string
| select cast(1 as float)
< ok 0 0
> exec ^
: string
| select convert(1, float(54))
< error
| 1064 42000 …near '(54))' at line 1
> exec ^
: string
| select convert(1, float(53))
< error
| 1064 42000 …near '(53))' at line 1
> exec ^
: string
| select cast(1 as float(54,2))
< error
| 1064 42000 …near '(54,2))' at line 1
> exec ^
: string
| select cast(1 as float(10,2))
< error
| 1064 42000 …near '(10,2))' at line 1
> exec ^
: string
| select cast(1 as double(10))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select * from ``.t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from ``.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select * from `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select * from `a`.t
< error
| 1146 42S02 Table 'a.t' doesn't exist
> exec ^
: string
| select * from `a`.t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select all 1
< ok 0 0
> exec ^
: string
| select all high_priority distinct 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all high_priority distinct 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| create table t (a int) stats_auto_recalc=0
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc=1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc=default
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='foo'
< error
| 1064 42000 …near '+interval 1 day ttl_enable='foo'' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='ON'
< error
| 1064 42000 …near '+interval 1 day ttl_enable='ON'' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='off'
< error
| 1064 42000 …near '+interval 1 day ttl_enable='off'' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='foo'
< error
| 1064 42000 …near '+interval 1 day ttl_job_interval='foo'' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='1h'
< error
| 1064 42000 …near '+interval 1 day ttl_job_interval='1h'' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='1d'
< error
| 1064 42000 …near '+interval 1 day ttl_job_interval='1d'' at line 1
> exec ^
: string
| create table t (a double(10))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a double(10,2))
< ok 0 0
> exec ^
: string
| create table t (a float(10))
< ok 0 0
> exec ^
: string
| create table t (a real(10))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| select cast(1 as year)
< error
| 1064 42000 …near 'year)' at line 1
> exec ^
: string
| grant r1, select on *.* to u
< error
| 1064 42000 …near ' select on *.* to u' at line 1
> exec ^
: string
| grant select, r1 on *.* to u
< error
| 1064 42000 …near 'r1 on *.* to u' at line 1
> exec ^
: string
| grant r1@h on *.* to u
< error
| 1064 42000 …near '@h on *.* to u' at line 1
> exec ^
: string
| grant select, r1@h on *.* to u
< error
| 1064 42000 …near 'r1@h on *.* to u' at line 1
> exec ^
: string
| grant r1@h, select on *.* to u
< error
| 1064 42000 …near '@h, select on *.* to u' at line 1
> exec ^
: string
| grant select to u
< error
| 1064 42000 …near 'to u' at line 1
> exec ^
: string
| grant r1, select to u
< error
| 1064 42000 …near ' select to u' at line 1
> exec ^
: string
| grant select, r1 to u
< error
| 1064 42000 …near 'r1 to u' at line 1
> exec ^
: string
| grant 'r1'@'h' to u
< error
| 1064 42000 …near '@'h' to u' at line 1
> exec ^
: string
| grant r1 on *.* to u
< error
| 1064 42000 …near 'on *.* to u' at line 1
> exec ^
: string
| revoke r1@h on *.* from u
< error
| 1064 42000 …near '@h on *.* from u' at line 1
> exec ^
: string
| revoke select, r1@h on *.* from u
< error
| 1064 42000 …near 'r1@h on *.* from u' at line 1
> exec ^
: string
| revoke select from u
< error
| 1064 42000 …near 'from u' at line 1
> exec ^
: string
| revoke r1, select from u
< error
| 1064 42000 …near ' select from u' at line 1
> exec ^
: string
| import into t (@a) from select 1
< error
| 1064 42000 …near 'import into t (@a) from select 1' at line 1
> exec ^
: string
| import into t (a) from select 1
< error
| 1064 42000 …near 'import into t (a) from select 1' at line 1
> exec ^
: string
| import into t (a) set a=1 from select 1
< error
| 1064 42000 …near 'import into t (a) set a=1 from select 1' at line 1
> exec ^
: string
| import into t from select 1
< error
| 1064 42000 …near 'import into t from select 1' at line 1
> exec ^
: string
| create sequence s start with -9223372036854775809
< error
| 1064 42000 %s value is out of range in '%s' near '9223372036854775809' at line 1
> exec ^
: string
| create sequence s start with -9223372036854775808
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| create sequence s start with -9223372036854775808
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| create sequence s start with 9223372036854775808
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| create sequence s start with 9223372036854775808
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| create sequence s increment by -9223372036854775809
< error
| 1064 42000 %s value is out of range in '%s' near '9223372036854775809' at line 1
> exec ^
: string
| create sequence s maxvalue 9223372036854775808
< ok 0 0
> exec ^
: string
| create sequence s minvalue -9223372036854775809
< ok 0 0
> exec ^
: string
| select setval(s, -9223372036854775809)
< error
| 1064 42000 %s value is out of range in '%s' near '9223372036854775809)' at line 1
> exec ^
: string
| select setval(s, 9223372036854775808)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| select setval(s, 9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create schema d encryption='yes'
< error
| 1064 42000 …near 'encryption='yes'' at line 1
> exec ^
: string
| create tablespace ts add datafile 'a' encryption='x'
< error
| 1064 42000 …near 'tablespace ts add datafile 'a' encryption='x'' at line 1
> exec ^
: string
| calibrate resource start_time '2020-01-01' start_time '2020-01-01'
< error
| 1064 42000 …near 'calibrate resource start_time '2020-01-01' start_time '2020-01-01'' at line 1
> exec ^
: string
| calibrate resource duration 'foo'
< error
| 1064 42000 …near 'calibrate resource duration 'foo'' at line 1
> exec ^
: string
| calibrate resource duration '1h' duration '1h'
< error
| 1064 42000 …near 'calibrate resource duration '1h' duration '1h'' at line 1
> exec ^
: string
| calibrate resource start_time '2020-01-01' end_time '2020-01-02' duration '1h'
< error
| 1064 42000 …near 'calibrate resource start_time '2020-01-01' end_time '2020-01-02' duration '1h'' at line 1
> exec ^
: string
| calibrate resource start_time '2020-01-01' duration '1h' start_time '2020-01-01'
< error
| 1064 42000 …near 'calibrate resource start_time '2020-01-01' duration '1h' start_time '2020-01-01'' at line 1
> exec ^
: string
| recover table by job 9223372036854775808
< error
| 1064 42000 …near 'recover table by job 9223372036854775808' at line 1
> exec ^
: string
| recover table by job 9223372036854775807
< error
| 1064 42000 …near 'recover table by job 9223372036854775807' at line 1
> exec ^
: string
| recover table t 9223372036854775808
< error
| 1064 42000 …near 'recover table t 9223372036854775808' at line 1
> exec ^
: string
| cancel distribution job 9223372036854775808
< error
| 1064 42000 …near 'cancel distribution job 9223372036854775808' at line 1
> exec ^
: string
| split table t between (1) and (2) regions 9223372036854775808
< error
| 1064 42000 …near 'split table t between (1) and (2) regions 9223372036854775808' at line 1
> exec ^
: string
| split table t between (1) and (2) regions 9223372036854775807
< error
| 1064 42000 …near 'split table t between (1) and (2) regions 9223372036854775807' at line 1
> exec ^
: string
| show br job 9223372036854775808
< error
| 1064 42000 …near 'br job 9223372036854775808' at line 1
> exec ^
: string
| show br job query 9223372036854775808
< error
| 1064 42000 …near 'br job query 9223372036854775808' at line 1
> exec ^
: string
| cancel br job 9223372036854775808
< error
| 1064 42000 …near 'cancel br job 9223372036854775808' at line 1
> exec ^
: string
| cancel import job 9223372036854775808
< error
| 1064 42000 …near 'cancel import job 9223372036854775808' at line 1
> exec ^
: string
| admin show ddl jobs 9223372036854775808
< error
| 1064 42000 …near 'admin show ddl jobs 9223372036854775808' at line 1
> exec ^
: string
| admin alter ddl jobs 9223372036854775808 thread=1
< error
| 1064 42000 …near 'admin alter ddl jobs 9223372036854775808 thread=1' at line 1
> exec ^
: string
| admin cancel ddl jobs 9223372036854775808
< error
| 1064 42000 …near 'admin cancel ddl jobs 9223372036854775808' at line 1
> exec ^
: string
| admin pause ddl jobs 9223372036854775808
< error
| 1064 42000 …near 'admin pause ddl jobs 9223372036854775808' at line 1
> exec ^
: string
| admin resume ddl jobs 9223372036854775808
< error
| 1064 42000 …near 'admin resume ddl jobs 9223372036854775808' at line 1
> exec ^
: string
| admin check index t i (1,2)
< error
| 1064 42000 …near 'admin check index t i (1,2)' at line 1
> exec ^
: string
| admin check index t i (9223372036854775808,2)
< error
| 1064 42000 …near 'admin check index t i (9223372036854775808,2)' at line 1
> exec ^
: string
| show import job 9223372036854775808
< error
| 1064 42000 …near 'import job 9223372036854775808' at line 1
> exec ^
: string
| show distribution job 9223372036854775808
< error
| 1064 42000 …near 'distribution job 9223372036854775808' at line 1
> exec ^
: string
| show profile for query 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| create user u with max_queries_per_hour 9223372036854775808
< ok 0 0
> exec ^
: string
| create user u with max_queries_per_hour 9223372036854775807
< ok 0 0
> exec ^
: string
| create user u with max_user_connections 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| create user u password expire interval 9223372036854775808 day
< error
| 1064 42000 …near '9223372036854775808 day' at line 1
> exec ^
: string
| create user u failed_login_attempts 9223372036854775808
< error
| 1064 42000 …near 'failed_login_attempts 9223372036854775808' at line 1
> exec ^
: string
| create user u password_lock_time 9223372036854775808
< error
| 1064 42000 …near 'password_lock_time 9223372036854775808' at line 1
> exec ^
: string
| select _foo'x'
< error
| 1054 42S22 Unknown column '_foo' in 'SELECT'
> exec ^
: string
| select _foo'x'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create table t (a char(1) unicode)
< ok 0 0
> exec ^
: string
| create table t (a char(1) charset ucs2)
< ok 0 0
> exec ^
: string
| select cast(1 as char unicode)
< ok 0 0
> exec ^
: string
| select cast(1 as char charset foo)
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| load data infile '/x' into table t fields optionally enclosed by 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields optionally enclosed by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields enclosed by X'0D0D'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields enclosed by X'0D0D'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields enclosed by B'110101101101011'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields enclosed by B'110101101101011'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields escaped by X'6B6B'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields escaped by X'6B6B'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields escaped by B'110101101101011'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields escaped by B'110101101101011'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '\\'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '\\'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields enclosed by ''
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields enclosed by ''
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '€'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '€'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields escaped by '€'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields escaped by '€'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'é'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'é'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t columns enclosed by 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t columns enclosed by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t fields terminated by 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t fields terminated by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table t lines terminated by 'aa'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| load data infile '/x' into table t lines terminated by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| import into t from '/x' with fields_enclosed_by='aa'
< error
| 1064 42000 …near 'import into t from '/x' with fields_enclosed_by='aa'' at line 1
> exec ^
: string
| select 1 into outfile '/x' fields optionally enclosed by 'aa'
< ok 0 0
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by 'a'
< ok 0 0
> exec ^
: string
| alter table t lock=select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| create table t (a float8(10))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a float8(10,2))
< ok 0 0
> exec ^
: string
| create table t (a float4(10))
< ok 0 0
> exec ^
: string
| create table t (a double precision(10))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| create table t (a double precision(10,2))
< ok 0 0
> exec ^
: string
| create table t (a double(10,0))
< ok 0 0
> exec ^
: string
| create table t (a real(10,0))
< ok 0 0
> exec ^
: string
| create function f() returns double(10) return 1
< error
| 1064 42000 …near ') return 1' at line 1
> exec ^
: string
| create table t (a int) encryption=Y
< ok 0 0
> exec ^
: string
| create table t (a int) encryption=1
< ok 0 0
> exec ^
: string
| create table t (a int) encryption=''
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval=1
< error
| 1064 42000 …near '+interval 1 day ttl_job_interval=1' at line 1
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval=foo
< error
| 1064 42000 …near '+interval 1 day ttl_job_interval=foo' at line 1
> exec ^
: string
| create table t (a int, b int as (a) serial default value)
< error
| 1064 42000 …near 'serial default value)' at line 1
> exec ^
: string
| create table t (a int, b int serial default value as (a))
< error
| 1064 42000 …near 'as (a))' at line 1
> exec ^
: string
| create table t (a int, b int generated always as (a) default 1)
< error
| 1064 42000 …near 'default 1)' at line 1
> exec ^
: string
| create table t (a int, b int default null generated always as (a))
< error
| 1064 42000 …near 'generated always as (a))' at line 1
> exec ^
: string
| create table t (a int) partition by key algorithm 2 (a)
< error
| 1064 42000 …near '2 (a)' at line 1
> exec ^
: string
| create table t (a int) partition by key algorithm = (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| create table t (a int) partition by key algorithm = 'x' (a)
< error
| 1064 42000 …near ''x' (a)' at line 1
> exec ^
: string
| create table t (a int) partition by key
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) partition by key partitions 2
< error
| 1064 42000 …near 'partitions 2' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions (partition p values less than (1))
< error
| 1064 42000 …near '(partition p values less than (1))' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by (partition p values less than (1))
< error
| 1064 42000 …near '(partition p values less than (1))' at line 1
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by linear (partition p values less than (1))
< error
| 1064 42000 …near '(partition p values less than (1))' at line 1
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D0D'
< ok 0 0
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D'
< ok 0 0
> exec ^
: string
| select 1 into outfile '/x' fields terminated by X'6B6B'
< ok 0 0
> exec ^
: string
| select 1 into outfile '/x' fields escaped by B'110101101101011'
< ok 0 0
> exec ^
: string
| select 1 into outfile '/x' fields enclosed 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| select setval(s, 1+1)
< error
| 1064 42000 …near '+1)' at line 1
> exec ^
: string
| select setval(s)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| select setval(s, 1)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| select setval(s, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select setval(s, -1)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| select setval(s, -1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select setval(s, +1)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| select setval(s, +1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| select setval(1, 1)
< error
| 1064 42000 …near '1, 1)' at line 1
> exec ^
: string
| select setval(s.t, 1)
< error
| 4091 42S02 Unknown SEQUENCE: 't'
> exec ^
: string
| select setval(s.t, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| show profile for query foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| show profile for query 1+1
< error
| 1064 42000 …near '+1' at line 1
> exec ^
: string
| show profile for query 1
< ok 0 0
> exec ^
: string
| show profile for query
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create user u with max_queries_per_hour
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create user u with max_user_connections 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| create user u password expire interval day
< error
| 1064 42000 …near 'day' at line 1
> exec ^
: string
| create user u password expire interval 1
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| create user u failed_login_attempts
< error
| 1064 42000 …near 'failed_login_attempts' at line 1
> exec ^
: string
| create user u password_lock_time
< error
| 1064 42000 …near 'password_lock_time' at line 1
> exec ^
: string
| alter user u with max_queries_per_hour 9223372036854775808
< ok 0 0
> exec ^
: string
| alter user u failed_login_attempts 9223372036854775808
< error
| 1064 42000 …near 'failed_login_attempts 9223372036854775808' at line 1
> exec ^
: string
| select * from ``
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| select * from ``
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| select a from ``.`b`
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select a from ``.`b`
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from `` t
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| select * from `` t
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| select `a `.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `a `.b from t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| load data infile '/x' into table `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| load data infile '/x' into table `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create index i using hnsw on t (a)
< error
| 1064 42000 …near 'hnsw on t (a)' at line 1
> exec ^
: string
| create index i using hnsw on t (a) using btree
< error
| 1064 42000 …near 'hnsw on t (a) using btree' at line 1
> exec ^
: string
| create index i on t (a) using hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| create index i on t (a) type hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| create unique index i on t (a) using hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| create fulltext index i on t (a) using hnsw
< error
| 1064 42000 …near 'using hnsw' at line 1
> exec ^
: string
| create spatial index i on t (a) using hnsw
< error
| 1064 42000 …near 'using hnsw' at line 1
> exec ^
: string
| create index i using btree on t (a) using hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| create table t (a int, key i (a) using hnsw)
< error
| 1064 42000 …near 'hnsw)' at line 1
> exec ^
: string
| alter table t add key i (a) using hnsw
< error
| 1064 42000 …near 'hnsw' at line 1
> exec ^
: string
| alter table t add vector index (a)
< ok 0 0
> exec ^
: string
| alter table t add vector index ((a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| alter table t add vector index ((a), (b))
< error
| 1064 42000 …near '(a), (b))' at line 1
> exec ^
: string
| alter table t add vector index i using hnsw ((a))
< error
| 1064 42000 …near 'using hnsw ((a))' at line 1
> exec ^
: string
| alter table t add vector index i using btree ((a))
< error
| 1064 42000 …near 'using btree ((a))' at line 1
> exec ^
: string
| alter table t add vector index ((a)) using hash
< error
| 1064 42000 …near '(a)) using hash' at line 1
> exec ^
: string
| alter table t add vector index ((a)) using hypo
< error
| 1064 42000 …near '(a)) using hypo' at line 1
> exec ^
: string
| alter table t add vector index i using hnsw ((a)) using btree
< error
| 1064 42000 …near 'using hnsw ((a)) using btree' at line 1
> exec ^
: string
| alter table t add vector index i using btree ((a)) using hnsw
< error
| 1064 42000 …near 'using btree ((a)) using hnsw' at line 1
> exec ^
: string
| alter table t add column c int, add (vector index i (a))
< ok 0 0
> exec ^
: string
| create table t (a int, b int, vector index (b))
< ok 0 0
> exec ^
: string
| create table t (a int, b int, vector index ((b)))
< error
| 1064 42000 …near '(b)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM 1 (a)
< error
| 1064 42000 …near '1 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=1
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=1 a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM=1 ()
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ()
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM='1' (a)
< error
| 1064 42000 …near ''1' (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS -1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1 (PARTITION p)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY foo (b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near 'foo (b) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY LINEAR foo (b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near 'foo (b) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near 'PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY LINEAR KEY (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near 'PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY () (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near ') (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY LINEAR KEY () (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near ') (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY ALGORITHM 1 (b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near '1 (b) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY KEY ALGORITHM=3 (b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near '3 (b) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY LINEAR KEY ALGORITHM=3 (b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near '3 (b) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near '(PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 'x' (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near ''x' (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s, SUBPARTITION u))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) (PARTITION p VALUES LESS THAN (1), PARTITION q VALUES LESS THAN (2) (SUBPARTITION s))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near 'SUBPARTITION s))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s), PARTITION q VALUES LESS THAN (2) (SUBPARTITION u))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s, SUBPARTITION u), PARTITION q VALUES LESS THAN (2) (SUBPARTITION v))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1), PARTITION q VALUES LESS THAN (2) (SUBPARTITION u))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near 'SUBPARTITION u))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1) (SUBPARTITION s), PARTITION q VALUES LESS THAN (2) (SUBPARTITION u))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) (PARTITION p VALUES LESS THAN (1) ())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1),(2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (MAXVALUE))
< error
| 1064 42000 Cannot use MAXVALUE as value in VALUES IN near 'MAXVALUE))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN ((1),(2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (1,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN ((1,2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (MAXVALUE))
< error
| 1064 42000 Cannot use MAXVALUE as value in VALUES IN near 'MAXVALUE))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1)))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near ')))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2),(3,4)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2,3)))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2,3)))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,MAXVALUE)))
< error
| 1064 42000 Cannot use MAXVALUE as value in VALUES IN near 'MAXVALUE)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (DEFAULT))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (DEFAULT))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1,2),3))
< error
| 1064 42000 …near '3))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (3,(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (3,(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1,2)))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (MAXVALUE))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN ())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p, PARTITION q (SUBPARTITION s))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near 'SUBPARTITION s))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p (SUBPARTITION s), PARTITION q)
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p)
< ok 0 0
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p VALUES IN (DEFAULT))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p, PARTITION q (SUBPARTITION s))
< error
| 1064 42000 Wrong number of subpartitions defined, mismatch with previous setting near 'SUBPARTITION s))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE COLUMNS (a) (PARTITION p VALUES LESS THAN (MAXVALUE))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (MAXVALUE,MAXVALUE))
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (MAXVALUE,MAXVALUE))
| ; SELECT FROM
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (1,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (MAXVALUE,MAXVALUE))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (MAXVALUE))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE COLUMNS (a,b) SUBPARTITION BY HASH (a) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES LESS THAN MAXVALUE)
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p VALUES LESS THAN MAXVALUE)
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| DROP VIEW `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP VIEW `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP VIEW `a `.v
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP VIEW `a `.v
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| LOCK TABLES `a ` READ
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| LOCK TABLES `a ` READ
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| LOCK TABLES `a `.t READ
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| LOCK TABLES `a `.t READ
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ANALYZE TABLE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ANALYZE TABLE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ANALYZE TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ANALYZE TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| OPTIMIZE TABLE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| OPTIMIZE TABLE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CHECK TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CHECK TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CHECKSUM TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CHECKSUM TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| REPAIR TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| REPAIR TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| FLUSH TABLES `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| FLUSH TABLES `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| FLUSH TABLES `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| FLUSH TABLES `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DESCRIBE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DESCRIBE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW COLUMNS FROM `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW COLUMNS FROM `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW INDEX FROM `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW INDEX FROM `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON `a `.t FOR EACH ROW SET @a=1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON `a `.t FOR EACH ROW SET @a=1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TRIGGER `a `.tr BEFORE INSERT ON t FOR EACH ROW SET @a=1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TRIGGER `a `.tr BEFORE INSERT ON t FOR EACH ROW SET @a=1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP TRIGGER `a `.tr
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP TRIGGER `a `.tr
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE EVENT `a `.e ON SCHEDULE EVERY 1 DAY DO SET @a=1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE EVENT `a `.e ON SCHEDULE EVERY 1 DAY DO SET @a=1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP EVENT `a `.e
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP EVENT `a `.e
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| HANDLER `a `.t OPEN
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| HANDLER `a `.t OPEN
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| TRUNCATE TABLE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| TRUNCATE TABLE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| TABLE `a `
< error
| 1064 42000 …near 'TABLE `a `' at line 1
> exec ^
: string
| TABLE `a `.t
< error
| 1064 42000 …near 'TABLE `a `.t' at line 1
> exec ^
: string
| SELECT `a `.t.* FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `a `.t.* FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ``.t.b FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT ``.t.b FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM t WHERE `a `.t.b = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t WHERE `a `.t.b = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER VIEW `a `.v AS SELECT 1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER VIEW `a `.v AS SELECT 1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER VIEW `a ` AS SELECT 1
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER VIEW `a ` AS SELECT 1
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t LIKE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t LIKE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t LIKE `a `.u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t LIKE `a `.u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE t RENAME `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER TABLE t RENAME `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER TABLE t RENAME TO `a `.u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE t RENAME TO `a `.u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| REPLACE INTO `a `.t VALUES (1)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| REPLACE INTO `a `.t VALUES (1)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE TABLE `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE TABLE `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE VIEW `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE VIEW `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE FUNCTION `a `.f() RETURNS INT RETURN 1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE FUNCTION `a `.f() RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE PROCEDURE `a `.p() SELECT 1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE PROCEDURE `a `.p() SELECT 1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP FUNCTION `a `.f
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP FUNCTION `a `.f
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP PROCEDURE `a `.p
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP PROCEDURE `a `.p
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CALL `a `.p()
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CALL `a `.p()
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT `a `.f()
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT `a `.f()
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE t EXCHANGE PARTITION p WITH TABLE `a `.u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER TABLE t EXCHANGE PARTITION p WITH TABLE `a `.u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES `a `.u (a))
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES `a `.u (a))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT REFERENCES `a `.u (a))
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t (a INT REFERENCES `a `.u (a))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GRANT SELECT ON `a `.t TO u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| GRANT SELECT ON `a `.t TO u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| GRANT SELECT ON `a `.* TO u
< ok 0 0
> exec ^
: string
| REVOKE SELECT ON `a `.t FROM u
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| REVOKE SELECT ON `a `.t FROM u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP INDEX i ON `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP INDEX i ON `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TEMPORARY TABLE `a ` (x INT)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TEMPORARY TABLE `a ` (x INT)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP TEMPORARY TABLE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP TEMPORARY TABLE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| RENAME TABLE `a ` TO u
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| RENAME TABLE `a ` TO u
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| RENAME TABLE t TO `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| RENAME TABLE t TO `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| RENAME TABLE t TO ``
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| RENAME TABLE t TO ``
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1) AS a
< ok 0 0
> exec ^
: string
| SELECT 1 FROM ((SELECT 1)) a
< ok 0 0
> exec ^
: string
| SELECT 1 FROM ((SELECT 1) a)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (((SELECT 1)))
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM (((SELECT 1))) a
< ok 0 0
> exec ^
: string
| SELECT 1 FROM ((SELECT 1) UNION (SELECT 2))
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM ((SELECT 1) UNION (SELECT 2)) a
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (t)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM (t)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (t, u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM (t, u)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (t JOIN u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM (t JOIN u)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (WITH c AS (SELECT 1) SELECT * FROM c)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM t, (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 FROM t LEFT JOIN (SELECT 1) ON 1
< error
| 1064 42000 …near 'ON 1' at line 1
> exec ^
: string
| SELECT 1 FROM (VALUES ROW(1)) a
< error
| 1064 42000 …near 'ROW(1)) a' at line 1
> exec ^
: string
| SELECT 1 FROM (TABLE t) a
< error
| 1064 42000 …near 'TABLE t) a' at line 1
> exec ^
: string
| SELECT 1 FROM LATERAL (SELECT 1) a
< error
| 1064 42000 …near '(SELECT 1) a' at line 1
> exec ^
: string
| SELECT 1 FROM LATERAL (SELECT 1) AS a
< error
| 1064 42000 …near '(SELECT 1) AS a' at line 1
> exec ^
: string
| SELECT 1 FROM t, LATERAL (VALUES ROW(1)) a
< error
| 1064 42000 …near '(VALUES ROW(1)) a' at line 1
> exec ^
: string
| DELETE a FROM (SELECT 1) a
< error
| 1288 HY000 The target table a of the DELETE is not updatable
> exec ^
: string
| DELETE a FROM (SELECT 1) a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t USING t, (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| DELETE t FROM t, (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| UPDATE t, (SELECT 1) SET t.a=1
< error
| 1064 42000 …near 'SET t.a=1' at line 1
> exec ^
: string
| UPDATE t, (SELECT 1) a SET t.a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t, (SELECT 1) a SET t.a=1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t SELECT 1 FROM (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE VIEW v AS SELECT 1 FROM (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t AS SELECT 1 FROM (SELECT 1)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT (SELECT 1 FROM (SELECT 1))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT 1 WHERE EXISTS (SELECT 1 FROM (SELECT 1))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]','$' COLUMNS (a INT PATH '$'))
< error
| 4177 HY000 Every table function must have an alias.
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]','$' COLUMNS (a INT PATH '$'))
| ; SELECT FROM
< error
| 4177 HY000 Every table function must have an alias.
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]','$' COLUMNS (a INT PATH '$')) j
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a (x)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1) (x)
< error
| 1064 42000 …near '(x)' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE DEFAULT
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE NEVER
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 5
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 5 MONTH
< error
| 1064 42000 …near 'MONTH' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 'x' DAY
< error
| 1064 42000 …near ''x' DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY
< error
| 1064 42000 …near 'HISTORY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY DEFAULT
< error
| 1064 42000 …near 'HISTORY DEFAULT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 0
< error
| 1064 42000 …near 'HISTORY 0' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 'x'
< error
| 1064 42000 …near 'HISTORY 'x'' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE
< error
| 1064 42000 …near 'REUSE' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL
< error
| 1064 42000 …near 'REUSE INTERVAL' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 5
< error
| 1064 42000 …near 'REUSE INTERVAL 5' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL DEFAULT
< error
| 1064 42000 …near 'REUSE INTERVAL DEFAULT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE DEFAULT
< error
| 1064 42000 …near 'REUSE DEFAULT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE 5 DAY
< error
| 1064 42000 …near 'REUSE 5 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 5 DAY
< error
| 1064 42000 …near 'REUSE INTERVAL 5 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 5 MONTH
< error
| 1064 42000 …near 'REUSE INTERVAL 5 MONTH' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 0 DAY
< error
| 1064 42000 …near 'REUSE INTERVAL 0 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE
< error
| 1064 42000 …near 'REQUIRE' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE DEFAULT
< error
| 1064 42000 …near 'REQUIRE DEFAULT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE CURRENT
< error
| 1064 42000 …near 'REQUIRE CURRENT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE CURRENT DEFAULT
< error
| 1064 42000 …near 'REQUIRE CURRENT DEFAULT' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE CURRENT OPTIONAL
< error
| 1064 42000 …near 'REQUIRE CURRENT OPTIONAL' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REQUIRE OPTIONAL
< error
| 1064 42000 …near 'REQUIRE OPTIONAL' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 'x'
< error
| 1064 42000 …near 'FAILED_LOGIN_ATTEMPTS 'x'' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 0
< error
| 1064 42000 …near 'FAILED_LOGIN_ATTEMPTS 0' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 0
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME 0' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME UNBOUNDED
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME UNBOUNDED' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 'x'
< error
| 1064 42000 …near 'PASSWORD_LOCK_TIME 'x'' at line 1
> exec ^
: string
| CREATE USER u ACCOUNT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u ACCOUNT LOCK
< ok 0 0
> exec ^
: string
| CREATE USER u ACCOUNT UNLOCK
< ok 0 0
> exec ^
: string
| CREATE USER u ACCOUNT foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 65536
< error
| 1064 42000 …near 'HISTORY 65536' at line 1
> exec ^
: string
| ALTER USER u PASSWORD REUSE INTERVAL 65536 DAY
< error
| 1064 42000 …near 'REUSE INTERVAL 65536 DAY' at line 1
> exec ^
: string
| ALTER USER u PASSWORD EXPIRE INTERVAL 0 DAY
< error
| 1525 HY000 Incorrect DAY value: '0'
> exec ^
: string
| ALTER USER u PASSWORD EXPIRE INTERVAL 0 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '0'
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=FIXED
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=DYNAMIC
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=COMPRESSED
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=REDUNDANT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=COMPACT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=PAGE
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT='FIXED'
< error
| 1064 42000 …near ''FIXED'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=`fixed`
< error
| 1064 42000 …near '`fixed`' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT FIXED
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT='0'
< error
| 1064 42000 …near ''0'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC='0'
< error
| 1064 42000 …near ''0'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=65535
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES='5'
< error
| 1064 42000 …near ''5'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS='0'
< error
| 1064 42000 …near ''0'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM=0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM=DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) COMPRESSION=foo
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COMPRESSION='foo'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTION=foo
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=FIRST
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) SECONDARY_ENGINE=NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) SECONDARY_ENGINE=foo
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4M
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4M' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=foo
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=foo' at line 1
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551616, 1
< error
| 1064 42000 …near '18446744073709551616, 1' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1, 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 2147483647
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 4294967296
< error
| 1064 42000 …near '4294967296' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< ok 0 0
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551616
< ok 0 0
> exec ^
: string
| KILL 18446744073709551616
< ok 0 0
> exec ^
: string
| EXPLAIN FORMAT='TREE' SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT='TREE' SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT=`tree` SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'tree'
> exec ^
: string
| EXPLAIN FORMAT=`tree` SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'tree'
> exec ^
: string
| DESCRIBE FORMAT=TREE SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| DESCRIBE FORMAT=TREE SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| DESC FORMAT=JSON SELECT 1
< ok 0 0
> exec ^
: string
| DESC ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec ^
: string
| DESCRIBE ANALYZE FORMAT=TREE SELECT 1
< error
| 1064 42000 …near 'ANALYZE FORMAT=TREE SELECT 1' at line 1
> exec ^
: string
| DESCRIBE SELECT 1
< ok 0 0
> exec ^
: string
| DESC SELECT 1
< ok 0 0
> exec ^
: string
| DESCRIBE format
< ok 0 0
> exec ^
: string
| DESCRIBE analyze
< error
| 1064 42000 …near 'analyze' at line 1
> exec ^
: string
| DESCRIBE t
< ok 0 0
> exec ^
: string
| DESC format = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| DESCRIBE FOR CONNECTION 1
< ok 0 0
> exec ^
: string
| DESCRIBE FORMAT=JSON FOR CONNECTION 1
< ok 0 0
> exec ^
: string
| DESCRIBE UPDATE t SET a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DESCRIBE UPDATE t SET a=1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DESCRIBE (SELECT 1)
< ok 0 0
> exec ^
: string
| EXPLAIN t
< ok 0 0
> exec ^
: string
| EXPLAIN format
< ok 0 0
> exec ^
: string
| EXPLAIN analyze
< error
| 1064 42000 …near 'analyze' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(5))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE PRECISION(5))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE PRECISION(5,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a REAL(5))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(5,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(5,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(5,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(04))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR ZEROFILL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT SIGNED UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT UNSIGNED UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT SIGNED SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT ZEROFILL ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT SIGNED ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT ZEROFILL SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT ZEROFILL UNSIGNED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT UNSIGNED ZEROFILL UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL UNSIGNED SIGNED)
< error
| 1064 42000 …near 'SIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT SIGNED UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec ^
: string
| CREATE TABLE t (a, b INT COLLATE ascii_bin AS (1))
< error
| 1064 42000 …near ' b INT COLLATE ascii_bin AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT COLLATE ascii_bin AS (1))
< error
| 1064 42000 …near 'AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin AS ('x'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b CHAR(5) NOT NULL AS ('x'))
< error
| 1064 42000 …near 'AS ('x'))' at line 1
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COMMENT 'x' AS ('x'))
< error
| 1064 42000 …near 'AS ('x'))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT GENERATED ALWAYS AS (1) GENERATED ALWAYS AS (2))
< error
| 1064 42000 …near 'GENERATED ALWAYS AS (2))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) AS (2))
< error
| 1064 42000 …near 'AS (2))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) NOT NULL)
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) STORED NOT NULL)
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) VIRTUAL COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT DEFAULT 1 GENERATED ALWAYS AS (1))
< error
| 1064 42000 …near 'GENERATED ALWAYS AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT PRIMARY KEY AS (1))
< error
| 1064 42000 …near 'AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT UNIQUE AS (1))
< error
| 1064 42000 …near 'AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT NULL AS (1))
< error
| 1064 42000 …near 'AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AUTO_INCREMENT GENERATED ALWAYS AS (1))
< error
| 1064 42000 …near 'GENERATED ALWAYS AS (1))' at line 1
> exec ^
: string
| CREATE TABLE t (b INT COLLATE ascii_bin COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET ascii COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) NOT NULL COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) NOT NULL COLLATE ascii_bin NULL COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) AS ('x') COLLATE ascii_bin)
< error
| 1064 42000 …near 'COLLATE ascii_bin)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin AS ('x') COLLATE ascii_bin)
< error
| 1064 42000 …near 'COLLATE ascii_bin)' at line 1
> exec ^
: string
| CREATE FUNCTION f(a CHAR(5) COLLATE ascii_bin COLLATE latin1_bin) RETURNS INT RETURN 1
< error
| 1064 42000 …near 'COLLATE latin1_bin) RETURNS INT RETURN 1' at line 1
> exec ^
: string
| CREATE FUNCTION f(a CHAR(5) COLLATE ascii_bin) RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8 CHARACTER SET utf8mb3
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8 COLLATE utf8_bin
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE utf8_bin CHARACTER SET utf8mb3
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb3_bin CHARACTER SET utf8
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE binary CHARACTER SET utf8mb4
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| CREATE DATABASE d COLLATE binary CHARACTER SET binary
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE utf8mb4_bin CHARACTER SET binary
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'binary'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET binary COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'binary'
> exec ^
: string
| CREATE DATABASE d COLLATE foo CHARACTER SET utf8mb4
< error
| 1273 HY000 Unknown collation: 'foo'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET foo CHARACTER SET bar
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 CHARACTER SET foo
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET foo CHARACTER SET utf8mb4
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET foo COLLATE utf8mb4_bin
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| CREATE DATABASE d COLLATE latin1_bin CHARACTER SET latin1
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE utf8mb4_bin CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE utf8mb4_bin CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE utf8mb4_bin CHARACTER SET utf8mb4
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE DATABASE d CHARACTER SET latin1 CHARACTER SET DEFAULT
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET DEFAULT CHARACTER SET DEFAULT
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE DEFAULT CHARACTER SET latin1
< ok 0 0
> exec ^
: string
| CREATE DATABASE d CHARACTER SET latin1 COLLATE DEFAULT
< ok 0 0
> exec ^
: string
| CREATE DATABASE d COLLATE latin1_bin COLLATE DEFAULT
< error
| 1302 HY000 Conflicting declarations: 'COLLATE latin1_bin' and 'COLLATE DEFAULT'
> exec ^
: string
| CREATE DATABASE d COLLATE latin1_bin COLLATE DEFAULT
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE latin1_bin' and 'COLLATE DEFAULT'
> exec ^
: string
| CREATE DATABASE d COLLATE latin1_bin CHARACTER SET DEFAULT
< ok 0 0
> exec ^
: string
| CREATE DATABASE d ENCRYPTION=N'Y'
< error
| 1064 42000 …near 'ENCRYPTION=N'Y'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='yes'
< error
| 1064 42000 …near 'ENCRYPTION='yes'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='x' ENCRYPTION='y'
< error
| 1064 42000 …near 'ENCRYPTION='x' ENCRYPTION='y'' at line 1
> exec ^
: string
| CREATE DATABASE d ENCRYPTION='y' ENCRYPTION='n'
< error
| 1064 42000 …near 'ENCRYPTION='y' ENCRYPTION='n'' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY=1
< error
| 1064 42000 …near 'READ ONLY=1' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY=2
< error
| 1064 42000 …near 'READ ONLY=2' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY=DEFAULT
< error
| 1064 42000 …near 'READ ONLY=DEFAULT' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY 0
< error
| 1064 42000 …near 'READ ONLY 0' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET utf8mb4 CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET utf8mb4 CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET utf8mb4 CHARACTER SET utf8mb4
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE latin1_bin CHARACTER SET utf8mb4
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE utf8mb4_bin COLLATE utf8mb4_general_ci
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE utf8mb4_general_ci'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE utf8mb4_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE utf8mb4_bin COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET latin1 COLLATE utf8mb4_bin
< error
| 1253 42000 COLLATION 'utf8mb4_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET DEFAULT CHARACTER SET latin1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE DEFAULT COLLATE utf8mb4_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE DEFAULT' and 'COLLATE utf8mb4_bin'
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE DEFAULT COLLATE utf8mb4_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE DEFAULT' and 'COLLATE utf8mb4_bin'
> exec ^
: string
| ALTER TABLE t COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| ALTER TABLE t COLLATE utf8mb4_bin COLLATE latin1_bin
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE utf8mb4_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| ALTER TABLE t CHARACTER SET utf8mb4, CHARACTER SET latin1
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER TABLE t CHARACTER SET utf8mb4, CHARACTER SET latin1
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'CHARACTER SET utf8mb4' and 'CHARACTER SET latin1'
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4 COLLATE latin1_bin
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_bin COLLATE latin1_bin
< error
| 1064 42000 …near 'COLLATE latin1_bin' at line 1
> exec ^
: string
| GRANT INSERT TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT CREATE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT USAGE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT ALL TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT ALL PRIVILEGES TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT ALTER TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT `select` TO u
< ok 0 0
> exec ^
: string
| GRANT r1, `select` TO u
< error
| 1064 42000 …near ' `select` TO u' at line 1
> exec ^
: string
| GRANT PROCESS TO u
< ok 0 0
> exec ^
: string
| GRANT RELOAD TO u
< ok 0 0
> exec ^
: string
| GRANT FILE TO u
< ok 0 0
> exec ^
: string
| GRANT SHUTDOWN TO u
< ok 0 0
> exec ^
: string
| GRANT SUPER TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION SLAVE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT SHOW VIEW TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT CREATE VIEW TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT CREATE ROLE TO u
< error
| 1064 42000 …near 'ROLE TO u' at line 1
> exec ^
: string
| GRANT BACKUP_ADMIN TO u
< ok 0 0
> exec ^
: string
| GRANT r1@h TO u
< error
| 1064 42000 …near '@h TO u' at line 1
> exec ^
: string
| REVOKE INSERT FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| REVOKE SELECT, r1 FROM u
< error
| 1064 42000 …near 'r1 FROM u' at line 1
> exec ^
: string
| REVOKE PROCESS FROM u
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=4294967295
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=4294967296
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=4294967295
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=4294967296
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=1.0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=1.0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=0
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=0' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4194304
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4194304' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4K
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4K' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4G
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4G' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE='4M'
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE='4M'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=1.5
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE='InnoDB'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE="InnoDB"
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE=`InnoDB`
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENGINE=select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHARSET=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHARSET='utf8mb4'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHARSET=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHARACTER SET=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DEFAULT CHARSET=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHARSET=binary
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHARSET=select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE='utf8mb4_bin'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE=binary
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COLLATE=select
< error
| 1064 42000 …near 'select' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE=DISK
< error
| 1064 42000 …near '=DISK' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE DISK
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE MEMORY
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE=foo
< error
| 1064 42000 …near '=foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE='DISK'
< error
| 1064 42000 …near '='DISK'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE=1
< error
| 1064 42000 …near '=1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE=foo
< error
| 1064 42000 …near '=foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE='foo'
< error
| 1064 42000 …near '='foo'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) TABLESPACE=1
< error
| 1064 42000 …near '=1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD='FIRST'
< error
| 1064 42000 …near ''FIRST'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=LAST
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=NO
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) SECONDARY_ENGINE='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) SECONDARY_ENGINE=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT=foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PAGE_CHECKSUM=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TRANSACTIONAL=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) IETF_QUOTES=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PAGE_COMPRESSED=1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION `f `() RETURNS INT RETURN 1
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| CREATE FUNCTION `f `() RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| CREATE FUNCTION ``() RETURNS INT RETURN 1
< error
| 1458 42000 Incorrect routine name ''
> exec ^
: string
| CREATE FUNCTION ``() RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name ''
> exec ^
: string
| CREATE PROCEDURE `p `() SELECT 1
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| CREATE PROCEDURE `p `() SELECT 1
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| DROP FUNCTION `f `
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| DROP FUNCTION `f `
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| DROP PROCEDURE `p `
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| DROP PROCEDURE `p `
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| CALL `p `()
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| CALL `p `()
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'p '
> exec ^
: string
| ALTER FUNCTION `a `.f COMMENT 'x'
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER FUNCTION `a `.f COMMENT 'x'
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER PROCEDURE `a `.p COMMENT 'x'
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER PROCEDURE `a `.p COMMENT 'x'
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| ALTER FUNCTION `f ` COMMENT 'x'
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| ALTER FUNCTION `f ` COMMENT 'x'
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'f '
> exec ^
: string
| SHOW CREATE FUNCTION `a `.f
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE FUNCTION `a `.f
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE PROCEDURE `a `.p
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE PROCEDURE `a `.p
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE TRIGGER `a `.p
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE TRIGGER `a `.p
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE EVENT `a `.p
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SHOW CREATE EVENT `a `.p
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE EVENT d.e ON SCHEDULE EVERY 1 DAY DO SET @a=1
< ok 0 0
> exec ^
: string
| DROP EVENT d.e
< ok 0 0
> exec ^
: string
| ALTER EVENT d.e COMMENT 'x'
< ok 0 0
> exec ^
: string
| SHOW CREATE EVENT d.e
< ok 0 0
> exec ^
: string
| CREATE TRIGGER d.tr BEFORE INSERT ON t FOR EACH ROW SET @a=1
< ok 0 0
> exec ^
: string
| DROP TRIGGER d.tr
< ok 0 0
> exec ^
: string
| CREATE TRIGGER `tr ` BEFORE INSERT ON t FOR EACH ROW SET @a=1
< error
| 1458 42000 Incorrect routine name 'tr '
> exec ^
: string
| CREATE TRIGGER `tr ` BEFORE INSERT ON t FOR EACH ROW SET @a=1
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'tr '
> exec ^
: string
| CREATE EVENT `e ` ON SCHEDULE EVERY 1 DAY DO SET @a=1
< error
| 1458 42000 Incorrect routine name 'e '
> exec ^
: string
| CREATE EVENT `e ` ON SCHEDULE EVERY 1 DAY DO SET @a=1
| ; SELECT FROM
< error
| 1458 42000 Incorrect routine name 'e '
> exec ^
: string
| CREATE TABLE `t` (`a` INT) COMMENT='x'
< ok 0 0
> exec ^
: string
| GRANT EVENT TO u
< ok 0 0
> exec ^
: string
| GRANT TRIGGER TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT EXECUTE TO u
< ok 0 0
> exec ^
: string
| GRANT REFERENCES TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT INDEX TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT DROP TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT UPDATE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT DELETE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT SHOW TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT LOCK TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT REPLICATION TO u
< ok 0 0
> exec ^
: string
| GRANT EXECUTE, r1 TO u
< error
| 1064 42000 …near 'r1 TO u' at line 1
> exec ^
: string
| GRANT r1, PROCESS TO u
< error
| 1064 42000 …near ' PROCESS TO u' at line 1
> exec ^
: string
| GRANT r1 (a) TO u
< error
| 1064 42000 …near '(a) TO u' at line 1
> exec ^
: string
| GRANT SELECT (a) TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT CREATE ROUTINE TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT ALL, r1 TO u
< error
| 1064 42000 …near ' r1 TO u' at line 1
> exec ^
: string
| GRANT ALL ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT ALL PRIVILEGES ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT PROCESS ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT `process` TO u
< ok 0 0
> exec ^
: string
| GRANT 'process' TO u
< ok 0 0
> exec ^
: string
| GRANT SELECT ON `a ` TO u
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| GRANT SELECT ON `a ` TO u
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| GRANT SELECT ON `` TO u
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| GRANT SELECT ON `` TO u
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name ''
> exec ^
: string
| GRANT SELECT ON ``.t TO u
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| GRANT SELECT ON ``.t TO u
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| GRANT SELECT ON `a `.t.b TO u
< error
| 1064 42000 …near '.b TO u' at line 1
> exec ^
: string
| REVOKE SELECT ON `a `.* FROM u
< ok 0 0
> exec ^
: string
| REVOKE EVENT FROM u
< ok 0 0
> exec ^
: string
| REVOKE ALL FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| REVOKE ALL PRIVILEGES FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| REVOKE ALL PRIVILEGES, GRANT OPTION FROM u
< ok 0 0
> exec ^
: string
| SHOW GRANTS FOR u USING SELECT
< error
| 1064 42000 …near 'USING SELECT' at line 1
> exec ^
: string
| CREATE FUNCTION f(a CHAR(5) CHARSET ascii COLLATE ascii_bin) RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a CHAR(5) CHARSET ascii) RETURNS INT RETURN 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS CHAR(5) COLLATE ascii_bin RETURN 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS CHAR(5) CHARSET ascii COLLATE ascii_bin RETURN 1
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(a CHAR(5) COLLATE ascii_bin) SELECT 1
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p(IN a TEXT COLLATE ascii_bin) SELECT 1
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v CHAR(5) COLLATE ascii_bin; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v CHAR(5) CHARSET ascii COLLATE ascii_bin; END
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (DEFAULT))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (DEFAULT))
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1),(2)))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '),(2)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1),(2,3)))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '),(2,3)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2),(3)))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near ')))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1,(2)))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1,(2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1),2))
< error
| 1064 42000 …near '2))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (1,(2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1)+1))
< error
| 1064 42000 …near '+1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN ((1)+1,2))
< error
| 1064 42000 …near '+1,2))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN ((1)+1))
< error
| 1064 42000 …near '+1))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1),(2)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1,2))
< ok 0 0
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p VALUES IN ((1)))
< error
| 1064 42000 Row expressions in VALUES IN only allowed for multi-field column partitioning near '))' at line 1
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p VALUES LESS THAN (1,2))
< ok 0 0
> exec ^
: string
| ALTER TABLE t PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1,2))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE COLUMNS (a,b) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 Inconsistency in usage of column lists for partitioning near '))' at line 1
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH (a) PARTITIONS 1 (PARTITION p, PARTITION q)
< error
| 1064 42000 Wrong number of partitions defined, mismatch with previous setting near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR KEY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR KEY ALGORITHM=2 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR KEY ALGORITHM 2 (a)
< error
| 1064 42000 …near '2 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR RANGE (a) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 …near 'RANGE (a) (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LINEAR LIST (a) (PARTITION p VALUES IN (1))
< error
| 1064 42000 …near 'LIST (a) (PARTITION p VALUES IN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1 PARTITIONS 2
< error
| 1064 42000 …near 'PARTITIONS 2' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 0
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 0
| ; SELECT FROM
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1))
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) SUBPARTITIONS 2 (PARTITION p VALUES LESS THAN (1))
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| SELECT aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.t.b FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.t.b FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| DROP TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| DROP TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.t
< error
| 1102 42000 Incorrect database name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| DROP TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a INT)
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE TABLE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa (a INT)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE TABLE t (aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa INT)
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| CREATE TABLE t (aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa INT)
| ; SELECT FROM
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| SELECT 1 AS aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
< ok 0 0
> exec ^
: string
| CREATE DATABASE aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
< ok 0 0
> exec ^
: string
| CREATE VIEW aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa AS SELECT 1
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE VIEW aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa AS SELECT 1
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE INDEX aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ON t (a)
< ok 0 0
> exec ^
: string
| RENAME TABLE t TO aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| RENAME TABLE t TO aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
> exec ^
: string
| CREATE FUNCTION aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa() RETURNS INT RETURN 1
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| CREATE FUNCTION aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa() RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1059 42000 Identifier name 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa' is too long
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=65535
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=65536
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE=4294967294
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=2.5
< error
| 1064 42000 …near '2.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=0.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT=1e2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC=1.0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC=1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=0.5
< error
| 1064 42000 …near '0.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=65535.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=65536.5
< error
| 1064 42000 …near '65536.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES=1e2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=2.5
< error
| 1064 42000 …near '2.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS=0.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=4294967295.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH=4294967296.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=1e2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS=.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STORAGE `DISK`
< error
| 1064 42000 …near '`DISK`' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) ROW_FORMAT=NOSUCH
< error
| 1064 42000 …near 'NOSUCH' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) INSERT_METHOD=`NO`
< error
| 1064 42000 …near '`NO`' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) COMPRESSION=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) COMPRESSION=''
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTION=YES
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTED=YES
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) ENCRYPTION_KEY_ID=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DATA DIRECTORY='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) INDEX DIRECTORY='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DATA DIRECTORY=x
< error
| 1064 42000 …near 'x' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) INDEX DIRECTORY=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) UNION=(a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) UNION=()
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) UNION (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) UNION a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) foo=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) foo='x'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) foo=bar
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) foo=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) foo=NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| ALTER TABLE t ENGINE=1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| ALTER TABLE t ROW_FORMAT=FOO
< error
| 1064 42000 …near 'FOO' at line 1
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec ^
: string
| ALTER TABLE t STATS_PERSISTENT=2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| ALTER TABLE t KEY_BLOCK_SIZE=4294967296
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p ENGINE=1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p MAX_ROWS='1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p MIN_ROWS='1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p TABLESPACE=1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p COMMENT=1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p DATA DIRECTORY=1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p NODEGROUP='1')
< error
| 1064 42000 …near ''1')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p NODEGROUP=1)
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY @a
< error
| 1064 42000 …near '@a' at line 1
> exec ^
: string
| SHOW PROFILE FOR 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| SELECT 1 FROM t, LATERAL ((SELECT 1)) a
< error
| 1064 42000 …near '((SELECT 1)) a' at line 1
> exec ^
: string
| SELECT 1 FROM t, LATERAL (WITH c AS (SELECT 1) SELECT * FROM c) a
< error
| 1064 42000 …near '(WITH c AS (SELECT 1) SELECT * FROM c) a' at line 1
> exec ^
: string
| SELECT 1 FROM t, LATERAL (TABLE t) a
< error
| 1064 42000 …near '(TABLE t) a' at line 1
> exec ^
: string
| SELECT 1 FROM t, LATERAL (t) a
< error
| 1064 42000 …near '(t) a' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) SUBPARTITION BY HASH (a) SUBPARTITIONS 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (DEFAULT(a)))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (DEFAULT(a)))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARACTER SET)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET 'ascii')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE 'ascii_bin')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHAR SET ascii)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET select)
< error
| 1064 42000 …near 'select)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE select)
< error
| 1064 42000 …near 'select)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) UNICODE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) UNICODE BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BINARY UNICODE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) ASCII)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) ASCII BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BINARY ASCII)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BYTE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BINARY BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET ascii BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BINARY CHARSET ascii)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) BINARY COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) COLLATE ascii_bin BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) CHARSET ascii COLLATE ascii_bin BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(5) UNICODE COLLATE ucs2_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(5) ASCII COLLATE ascii_bin)
< error
| 1253 42000 COLLATION 'ascii_bin' is not valid for CHARACTER SET 'latin1'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) ASCII CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a TEXT UNICODE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(5) ASCII)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a ENUM('a') UNICODE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a ENUM('a') ASCII BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SET('a') BYTE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT BYTE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR BYTE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE t (a LONG UNICODE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a LONG VARCHAR ASCII)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHAR(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u WITH
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551616
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW PROFILE LIMIT a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW PROFILE LIMIT 1 OFFSET a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW PROFILE LIMIT 1 OFFSET a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW PROFILE LIMIT 1 OFFSET 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 1, 2
< ok 0 0
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW ERRORS LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 1 OFFSET 2
< ok 0 0
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| UPDATE t SET a=1 LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS TERMINATED 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS TERMINATED BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ESCAPED 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ESCAPED BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS OPTIONALLY ENCLOSED 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS OPTIONALLY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES STARTING 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES STARTING BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES TERMINATED 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES TERMINATED BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY X'61'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY 0x61
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY B'01100001'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS TERMINATED BY 0x6161
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ESCAPED BY X'61'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES STARTING BY X'61'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' LINES TERMINATED BY X'61'
< ok 0 0
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY N'a'
< error
| 1064 42000 …near 'N'a'' at line 1
> exec ^
: string
| SELECT 1 INTO OUTFILE '/x' FIELDS ENCLOSED BY _utf8'a'
< error
| 1064 42000 …near '_utf8'a'' at line 1
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED BY 0x61
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED BY 0x61
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED BY X'61'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED BY X'61'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS TERMINATED BY B'01100001'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS TERMINATED BY B'01100001'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t LINES STARTING BY X'61'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t LINES STARTING BY X'61'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t LINES TERMINATED BY X'61'
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t LINES TERMINATED BY X'61'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED BY 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| LOAD DATA INFILE 'x' INTO TABLE t FIELDS ENCLOSED 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=M
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=M' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4MB
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4MB' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=4T
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=4T' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE=`4M`
< error
| 1064 42000 …near 'AUTOEXTEND_SIZE=`4M`' at line 1
> exec ^
: string
| EXPLAIN FORMAT=TREE FOR CONNECTION 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| EXPLAIN FORMAT=TREE FOR CONNECTION 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1,(2,3),4))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN (1,(2,3),4))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((1,2),3,(4,5)))
< error
| 1064 42000 …near '3,(4,5)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY LINEAR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW ERRORS LIMIT 1, 2
< ok 0 0
> exec ^
: string
| SHOW ERRORS LIMIT 1 OFFSET 2
< ok 0 0
> exec ^
: string
| SHOW ERRORS LIMIT a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW ERRORS LIMIT a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 1 OFFSET 2
< ok 0 0
> exec ^
: string
| SHOW RELAYLOG EVENTS LIMIT 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a BLOB COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a BLOB ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| CREATE TABLE t (a TINYBLOB BINARY)
< error
| 1064 42000 …near 'BINARY)' at line 1
> exec ^
: string
| CREATE TABLE t (a TINYBLOB COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a MEDIUMBLOB CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a LONGBLOB UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE t (a LONG VARBINARY ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| CREATE TABLE t (a LONG VARBINARY COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a VARBINARY(5) COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a BINARY(5) COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a NCHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(5) ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| CREATE TABLE t (a NCHAR(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a NCHAR(5) COLLATE utf8mb3_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NVARCHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NVARCHAR(5) COLLATE utf8mb3_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NVARCHAR(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(5) UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHAR(5) COLLATE utf8mb3_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL CHAR(5) UNICODE)
< error
| 1064 42000 …near 'UNICODE)' at line 1
> exec ^
: string
| CREATE TABLE t (a NCHAR VARCHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR VARYING(5) CHARSET ascii)
< error
| 1064 42000 …near 'CHARSET ascii)' at line 1
> exec ^
: string
| CREATE TABLE t (a BIT COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a GEOMETRY COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT COLLATE ascii_bin PRIMARY KEY)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) ASCII)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) UNICODE)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) BYTE)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) BINARY)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(5) CHARSET ascii)
< ok 0 0
> exec ^
: string
| SELECT CONVERT(1, CHAR(5) ASCII)
< ok 0 0
> exec ^
: string
| SELECT CONVERT(1 USING ascii)
< ok 0 0
> exec ^
: string
| SELECT `abs`()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT `abs`()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT `abs`(1,2)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT `abs`(1,2)
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'abs'
> exec ^
: string
| SELECT ABS ()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT ABS ()
| ; SELECT FROM
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec ^
: string
| SELECT d.abs()
< error
| 1305 42000 FUNCTION d.abs does not exist
> exec ^
: string
| SELECT d.abs()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.abs(1,2)
< error
| 1305 42000 FUNCTION d.abs does not exist
> exec ^
: string
| SELECT d.abs(1,2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (b INT AS (1) UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) UNIQUE KEY)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) KEY)
< error
| 1064 42000 …near 'KEY)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) PRIMARY KEY)
< error
| 1064 42000 …near 'PRIMARY KEY)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) INVISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) DEFAULT 1)
< error
| 1064 42000 …near 'DEFAULT 1)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) AUTO_INCREMENT)
< error
| 1064 42000 …near 'AUTO_INCREMENT)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) ON UPDATE NOW())
< error
| 1064 42000 …near 'ON UPDATE NOW())' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) CHECK (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) REFERENCES u (a))
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) COLLATE ascii_bin)
< error
| 1064 42000 …near 'COLLATE ascii_bin)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) VIRTUAL UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) PERSISTENT COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) STORED INVISIBLE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) STORED UNIQUE COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) COMMENT 'x' UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) UNIQUE UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) COMMENT 'x' COMMENT 'y')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT AS (1) WITH SYSTEM VERSIONING)
< error
| 1064 42000 …near 'SYSTEM VERSIONING)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT AS (1) WITHOUT SYSTEM VERSIONING)
< error
| 1064 42000 …near 'WITHOUT SYSTEM VERSIONING)' at line 1
> exec ^
: string
| CREATE TABLE t (b INT GENERATED ALWAYS AS (1) COMMENT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT GENERATED ALWAYS AS (1) NOT NULL)
< error
| 1064 42000 …near 'NOT NULL)' at line 1
> exec ^
: string
| CREATE TABLE t (b CHAR(5) CHARSET ascii AS ('x'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b CHAR(5) BINARY AS ('x'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b INT UNSIGNED AS (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin NOT NULL COLLATE ascii_bin)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin NOT NULL COLLATE ascii_general_ci)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE ascii_general_ci'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin NOT NULL COLLATE ascii_general_ci)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE ascii_general_ci'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin COLLATE latin1_bin)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) COLLATE ascii_bin COLLATE latin1_bin)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) NOT NULL COLLATE ascii_bin NULL COLLATE latin1_bin)
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) NOT NULL COLLATE ascii_bin NULL COLLATE latin1_bin)
| ; SELECT FROM
< error
| 1302 HY000 Conflicting declarations: 'COLLATE ascii_bin' and 'COLLATE latin1_bin'
> exec ^
: string
| CREATE TABLE t (b CHAR(5) CHARSET ascii COLLATE latin1_bin)
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'ascii'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p DEFAULT)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1,2),(3,4,5)))
< error
| 1657 HY000 Cannot have more than one value for this type of LIST partitioning
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN ((1,2),(3,4,5)))
| ; SELECT FROM
< error
| 1657 HY000 Cannot have more than one value for this type of LIST partitioning
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1,2,3))
< ok 0 0
> exec ^
: string
| ALTER TABLE t REORGANIZE PARTITION a INTO (PARTITION p DEFAULT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p DEFAULT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1), PARTITION q DEFAULT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p DEFAULT, PARTITION q DEFAULT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p DEFAULT VALUES IN (1))
< error
| 1064 42000 …near 'IN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x CURRENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y)
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY, PARTITION y)
| ; SELECT FROM
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y CURRENT)
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x, PARTITION y CURRENT)
| ; SELECT FROM
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x DEFAULT)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x DEFAULT)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES DEFAULT in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY SYSTEM_TIME (PARTITION x HISTORY VALUES IN (1))
< error
| 1064 42000 …near 'IN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION x HISTORY)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `KEY (``a```
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY (a) (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `KEY (``a```
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH (a) (PARTITION x HISTORY)
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| ALTER TABLE t PARTITION BY HASH (a) (PARTITION x HISTORY)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partition type `SYSTEM_TIME` for partitioning by `HASH`
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME (PARTITION x)
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| ALTER TABLE t PARTITION BY SYSTEM_TIME (PARTITION x)
| ; SELECT FROM
< error
| 4128 HY000 Wrong partitions for `t`: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| SELECT 1 FROM lateral
< error
| 1146 42S02 Table 'test.lateral' doesn't exist
> exec ^
: string
| SELECT 1 FROM lateral
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM lateral a
< error
| 1146 42S02 Table 'test.lateral' doesn't exist
> exec ^
: string
| SELECT 1 FROM lateral a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT lateral FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT lateral FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE lateral (a INT)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t, LATERAL (SELECT 1) a
< error
| 1064 42000 …near '(SELECT 1) a' at line 1
> exec ^
: string
| SELECT 1 FROM t JOIN LATERAL (SELECT 1) a ON 1
< error
| 1064 42000 …near '(SELECT 1) a ON 1' at line 1
> exec ^
: string
| EXPLAIN EXTENDED SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN PARTITIONS SELECT 1
< ok 0 0
> exec ^
: string
| EXPLAIN ANALYZE SELECT 1
< error
| 1064 42000 …near 'ANALYZE SELECT 1' at line 1
> exec ^
: string
| ANALYZE SELECT 1
< ok 0 0
> exec ^
: string
| ANALYZE FORMAT=TRADITIONAL SELECT 1
< ok 0 0
> exec ^
: string
| ANALYZE FORMAT=TREE SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| ANALYZE FORMAT=TREE SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'TREE'
> exec ^
: string
| ANALYZE FORMAT='JSON' SELECT 1
< ok 0 0
> exec ^
: string
| ANALYZE FORMAT=JSON UPDATE t SET a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE FORMAT=JSON UPDATE t SET a=1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ANALYZE FORMAT=JSON
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ANALYZE UPDATE t SET a=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE UPDATE t SET a=1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ANALYZE TABLE t
< ok 0 0
> exec ^
: string
| ANALYZE format
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ANALYZE t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| DESC EXTENDED SELECT 1
< ok 0 0
> exec ^
: string
| DESC PARTITIONS SELECT 1
< ok 0 0
> exec ^
: string
| DESCRIBE FORMAT=foo SELECT 1
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| DESCRIBE FORMAT=foo SELECT 1
| ; SELECT FROM
< error
| 1791 HY000 Unknown EXPLAIN/ANALYZE format name: 'foo'
> exec ^
: string
| SELECT SQL_BUFFER_RESULT SQL_BUFFER_RESULT 1
< ok 0 0
> exec ^
: string
| SELECT SQL_CALC_FOUND_ROWS SQL_CALC_FOUND_ROWS 1
< ok 0 0
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) UNION (SELECT 2) LOCK IN SHARE MODE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) UNION (SELECT 2) LOCK IN SHARE MODE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t (a INT) SELECT 1 UNION SELECT 2 FOR UPDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) (SELECT 1) FOR UPDATE
< ok 0 0
> exec ^
: string
| CREATE TABLE t AS (SELECT 1) UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE TABLE t AS (SELECT 1) UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a NOT IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a NOT IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ALL (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ALL (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a > SOME (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a > SOME (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = (SELECT 1 LIMIT 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT 1 LIMIT 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE EXISTS (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM (SELECT 1 LIMIT 1) d
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM (SELECT 1 LIMIT 1) d)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM (SELECT 1 LIMIT 1) d)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 UNION SELECT 2 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 UNION SELECT 2 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1))
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1))
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM t WHERE b = (SELECT 1 LIMIT 1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM t WHERE b = (SELECT 1 LIMIT 1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM t LIMIT 1 OFFSET 2)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 1 FROM t LIMIT 1 OFFSET 2)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE (a,b) IN (SELECT 1,2 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE (a,b) IN (SELECT 1,2 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| DELETE FROM t WHERE a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| DELETE FROM t WHERE a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM (t, u) JOIN (v, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM (t, u) JOIN (v, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a, (SELECT 1) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM (SELECT 1) a, (SELECT 1) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM (SELECT 1) t, t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM (SELECT 1) t, t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t JOIN (SELECT 1) t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t JOIN (SELECT 1) t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t a, u `A`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t a, u `A`
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t a, u 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 WHERE EXISTS (SELECT 1 FROM t, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 WHERE EXISTS (SELECT 1 FROM t, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT (SELECT 1 FROM t, t)
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT (SELECT 1 FROM t, t)
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| SELECT 1 FROM t, u WHERE EXISTS (SELECT 1 FROM t, u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t, u WHERE EXISTS (SELECT 1 FROM t, u)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t a, u a SET x=1
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| UPDATE t a, u a SET x=1
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| DELETE t, u FROM t, u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE t, u FROM t, u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE t.* FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE t.* FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE t.*, u.* FROM t, u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE t.*, u.* FROM t, u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t USING t, t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE FROM t USING t, t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE FROM t, t USING t
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE FROM t, t USING t
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 't'
> exec ^
: string
| DELETE FROM t, u USING t, u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t, u USING t, u
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM u USING t
< error
| 1109 42S02 Unknown table 'u' in MULTI DELETE
> exec ^
: string
| DELETE FROM u USING t
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'u' in MULTI DELETE
> exec ^
: string
| DELETE d.t1 FROM d.t1
< error
| 1146 42S02 Table 'd.t1' doesn't exist
> exec ^
: string
| DELETE d.t1 FROM d.t1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE a FROM t a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE a FROM t a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE t FROM t a
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE t FROM t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE d.a FROM t a
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE d.a FROM t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE a FROM t a JOIN u ON 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE a FROM t a JOIN u ON 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE b FROM (t a, u b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE b FROM (t a, u b)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE b FROM t a JOIN (u b, v c) ON 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE b FROM t a JOIN (u b, v c) ON 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE c FROM t a WHERE EXISTS (SELECT 1 FROM u c)
< error
| 1109 42S02 Unknown table 'c' in MULTI DELETE
> exec ^
: string
| DELETE c FROM t a WHERE EXISTS (SELECT 1 FROM u c)
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'c' in MULTI DELETE
> exec ^
: string
| SELECT `a` FROM t `a `
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `a` FROM t `a `
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM `t `
< error
| 1103 42000 Incorrect table name 't '
> exec ^
: string
| SELECT 1 FROM `t `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 't '
> exec ^
: string
| SELECT 1 AS `a `
< ok 0 0
> exec ^
: string
| SELECT 1 AS ``
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) UNION=(`a `)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t (a INT) UNION=(`a `)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t (a INT) UNION=(`a `.t)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t (a INT) UNION=(`a `.t)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE TABLE t (a INT REFERENCES `a ` (a))
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t (a INT REFERENCES `a ` (a))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES `a ` (a))
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES `a ` (a))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES `a `.u (a)
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE SEQUENCE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE SEQUENCE `a `.s
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| CREATE SEQUENCE `a `.s
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| DROP SEQUENCE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| DROP SEQUENCE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER SEQUENCE `a ` RESTART
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| ALTER SEQUENCE `a ` RESTART
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT NEXTVAL(`a `)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT NEXTVAL(`a `)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT NEXT VALUE FOR `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT NEXT VALUE FOR `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SHOW CREATE SEQUENCE `a `
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SHOW CREATE SEQUENCE `a `
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| CREATE PROCEDURE p(a INT(256)) SELECT 1
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE PROCEDURE p(a INT(256)) SELECT 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE PROCEDURE p(a CHAR(256)) SELECT 1
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE PROCEDURE p(a CHAR(256)) SELECT 1
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE PROCEDURE p(a DECIMAL(66)) SELECT 1
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE PROCEDURE p(a DECIMAL(66)) SELECT 1
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE PROCEDURE p(a TIME(7)) SELECT 1
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE PROCEDURE p(a TIME(7)) SELECT 1
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 6
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT(256) RETURN 1
< error
| 1439 42000 Display width out of range for '' (max = 255)
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT(256) RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for '' (max = 255)
> exec ^
: string
| CREATE FUNCTION f() RETURNS DECIMAL(66) RETURN 1
< error
| 1426 42000 Too big precision specified for ''. Maximum is 65
> exec ^
: string
| CREATE FUNCTION f() RETURNS DECIMAL(66) RETURN 1
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for ''. Maximum is 65
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v INT(256); END
< error
| 1439 42000 Display width out of range for 'v' (max = 255)
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v INT(256); END
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'v' (max = 255)
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v DECIMAL(66); END
< error
| 1426 42000 Too big precision specified for 'v'. Maximum is 65
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE v DECIMAL(66); END
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'v'. Maximum is 65
> exec ^
: string
| SELECT CAST(1 AS DECIMAL(66))
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 65
> exec ^
: string
| SELECT CAST(1 AS DECIMAL(66))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 65
> exec ^
: string
| SELECT CAST(1 AS DECIMAL(5,6))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '')
> exec ^
: string
| SELECT CAST(1 AS DECIMAL(5,6))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '')
> exec ^
: string
| SELECT CAST(1 AS CHAR(256))
< ok 0 0
> exec ^
: string
| SELECT CAST(1 AS CHAR(4294967296))
< error
| 1439 42000 Display width out of range for '1' (max = 4294967295)
> exec ^
: string
| SELECT CAST(1 AS CHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for '1' (max = 4294967295)
> exec ^
: string
| SELECT CAST(1 AS DATETIME(7))
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 6
> exec ^
: string
| SELECT CAST(1 AS DATETIME(7))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for '1'. Maximum is 6
> exec ^
: string
| CREATE TABLE t (a BINARY(255))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(4294967295))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(4294967295))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(4294967296))
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(4294967295))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(65,38))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(255,30))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(255,30))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a REAL(255,30))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a REAL(256,2))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a REAL(256,2))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DATETIME(6))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(6))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(256))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BOOL(1))
< error
| 1064 42000 …near '(1))' at line 1
> exec ^
: string
| CREATE TABLE t (a SMALLINT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a SMALLINT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a MEDIUMINT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a MEDIUMINT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT1(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT1(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a MIDDLEINT(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a MIDDLEINT(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INTEGER(256))
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INTEGER(256))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DEC(66))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a DEC(66))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a FIXED(66))
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a FIXED(66))
| ; SELECT FROM
< error
| 1426 42000 Too big precision specified for 'a'. Maximum is 65
> exec ^
: string
| CREATE TABLE t (a NCHAR(256))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a NCHAR(256))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a NVARCHAR(4294967296))
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NVARCHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VECTOR(4294967296))
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VECTOR(4294967296))
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column 'a' (max = 16383); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a CHAR(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(0,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(0,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT AUTO_INCREMENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT AUTO_INCREMENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE AUTO_INCREMENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DECIMAL AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DECIMAL(5,2) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DECIMAL(5,2) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a BIT AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a BIT AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a BOOL AUTO_INCREMENT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a YEAR AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a CHAR(5) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a BLOB AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a BLOB AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a ENUM('a') AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a ENUM('a') AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a SET('a') AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a SET('a') AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TIME AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TIME AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DATETIME AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a DATETIME AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a JSON AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a JSON AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a GEOMETRY AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a GEOMETRY AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a VECTOR(3) AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE NOW())
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) ON UPDATE NOW())
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(3) ON UPDATE CURRENT_TIMESTAMP(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(3) ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a TIME ON UPDATE NOW())
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a YEAR ON UPDATE NOW())
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a YEAR ON UPDATE NOW())
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE 1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE LOCALTIME)
< error
| 1064 42000 …near 'LOCALTIME)' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE LOCALTIMESTAMP(0))
< error
| 1064 42000 …near 'LOCALTIMESTAMP(0))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL NOT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL NOT NULL)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT KEY DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT KEY DEFAULT NULL)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT UNIQUE DEFAULT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL UNIQUE)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL PRIMARY KEY)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL PRIMARY KEY)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL KEY)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NULL KEY)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT NULL DEFAULT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL NULL DEFAULT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT NULL NOT NULL DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT NULL NOT NULL DEFAULT NULL)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT DEFAULT 1)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT DEFAULT 1)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1 AUTO_INCREMENT)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1 AUTO_INCREMENT)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT DEFAULT NULL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT UNIQUE DEFAULT 1)
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT AUTO_INCREMENT UNIQUE DEFAULT 1)
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a DATE DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON DEFAULT 'x')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) DEFAULT NOW(2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(3) DEFAULT NOW())
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT DEFAULT 2)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(1) DEFAULT 'xx')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a ENUM('a') DEFAULT 'b')
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL DEFAULT (NULL))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT NOT NULL DEFAULT (NULL))
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT(256)
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT(256)
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for 'a' (max = 255)
> exec ^
: string
| ALTER TABLE t MODIFY a VARCHAR(5) AUTO_INCREMENT
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t MODIFY a VARCHAR(5) AUTO_INCREMENT
| ; SELECT FROM
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| ALTER TABLE t CHANGE a b INT ON UPDATE NOW()
< error
| 1294 HY000 Invalid ON UPDATE clause for 'b' column
> exec ^
: string
| ALTER TABLE t CHANGE a b INT ON UPDATE NOW()
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'b' column
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT NOT NULL DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT NOT NULL DEFAULT NULL
| ; SELECT FROM
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a INT, j INT) SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| CREATE TABLE t (a INT, j INT) SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t a, JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a, (SELECT 1) a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) a, (SELECT 1) a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t `a`, u `a`
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t `a`, u `a`
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM t `a `, u `a `
< error
| 1066 42000 Not unique table/alias: 'a '
> exec ^
: string
| SELECT 1 FROM t `a `, u `a `
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a '
> exec ^
: string
| SELECT 1 FROM d.t a, d.u a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM d.t a, d.u a
| ; SELECT FROM
< error
| 1066 42000 Not unique table/alias: 'a'
> exec ^
: string
| SELECT 1 FROM d.t a, e.u a
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM d.t a, e.u a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM d.t a, u a
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM d.t a, u a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM dual, dual
< error
| 1064 42000 …near ' dual' at line 1
> exec ^
: string
| DELETE a FROM d.t a
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| DELETE a FROM d.t a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE d.a FROM d.t a
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE d.a FROM d.t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE e.a FROM d.t a
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE e.a FROM d.t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'a' in MULTI DELETE
> exec ^
: string
| DELETE t.* FROM t a
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE t.* FROM t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE a.* FROM t a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE a.* FROM t a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE d.t.* FROM d.t
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| DELETE d.t.* FROM d.t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM d.t USING d.t
< error
| 1146 42S02 Table 'd.t' doesn't exist
> exec ^
: string
| DELETE FROM d.t USING d.t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t USING t a
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE FROM t USING t a
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't' in MULTI DELETE
> exec ^
: string
| DELETE FROM a USING t a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM a USING t a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM a USING (SELECT 1) a
< error
| 1288 HY000 The target table a of the DELETE is not updatable
> exec ^
: string
| DELETE FROM a USING (SELECT 1) a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM a.* USING t a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM a.* USING t a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE t2, a FROM t2 JOIN t3 a ON 1
< error
| 1146 42S02 Table 'test.t2' doesn't exist
> exec ^
: string
| DELETE t2, a FROM t2 JOIN t3 a ON 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE t2, b FROM t2 JOIN t3 a ON 1
< error
| 1109 42S02 Unknown table 'b' in MULTI DELETE
> exec ^
: string
| DELETE t2, b FROM t2 JOIN t3 a ON 1
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 'b' in MULTI DELETE
> exec ^
: string
| DELETE t2 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) t2
< error
| 1109 42S02 Unknown table 't2' in MULTI DELETE
> exec ^
: string
| DELETE t2 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) t2
| ; SELECT FROM
< error
| 1109 42S02 Unknown table 't2' in MULTI DELETE
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1) UNION SELECT 2)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1) UNION SELECT 2)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 2 UNION (SELECT 1 LIMIT 1))
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (SELECT 2 UNION (SELECT 1 LIMIT 1))
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1) UNION (SELECT 2 LIMIT 1))
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 LIMIT 1) UNION (SELECT 2 LIMIT 1))
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 UNION SELECT 2) LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN ((SELECT 1 UNION SELECT 2) LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (WITH c AS (SELECT 1) SELECT * FROM c LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (WITH c AS (SELECT 1) SELECT * FROM c LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (WITH c AS (SELECT 1 LIMIT 1) SELECT * FROM c)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (WITH c AS (SELECT 1 LIMIT 1) SELECT * FROM c)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (VALUES (1) LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a IN (VALUES (1) LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT 1 UNION SELECT 2 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a = ANY (SELECT 1 UNION SELECT 2 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a <> ALL ((SELECT 1 LIMIT 1))
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t WHERE a <> ALL ((SELECT 1 LIMIT 1))
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t HAVING a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t HAVING a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| UPDATE t SET a = 1 WHERE a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| UPDATE t SET a = 1 WHERE a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t JOIN u ON a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t JOIN u ON a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t GROUP BY a IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 FROM t GROUP BY a IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 ORDER BY 1 IN (SELECT 1 LIMIT 1)
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| SELECT 1 ORDER BY 1 IN (SELECT 1 LIMIT 1)
| ; SELECT FROM
< error
| 1235 42000 This version of MariaDB doesn't yet support 'LIMIT & IN/ALL/ANY/SOME subquery'
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IN (SELECT 1 LIMIT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IN (SELECT 1 LIMIT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| SELECT DISTINCT 1
< ok 0 0
> exec ^
: string
| SELECT DISTINCTROW 1
< ok 0 0
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) LOCK IN SHARE MODE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) LOCK IN SHARE MODE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION SELECT 2 FOR UPDATE
< ok 0 0
> exec ^
: string
| SELECT 1 UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| SELECT 1 UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| SELECT 1 UNION SELECT 2 FOR UPDATE
< ok 0 0
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) LIMIT 1 FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION (SELECT 2) LIMIT 1 FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1 FOR UPDATE) UNION (SELECT 2)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM (SELECT 1 UNION SELECT 2 FOR UPDATE) d
< ok 0 0
> exec ^
: string
| SELECT 1 FROM ((SELECT 1) UNION (SELECT 2) FOR UPDATE) d
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| SELECT 1 FROM ((SELECT 1) UNION (SELECT 2) FOR UPDATE) d
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| INSERT INTO t (SELECT 1) UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| INSERT INTO t (SELECT 1) UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE VIEW v AS (SELECT 1) UNION (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE VIEW v AS (SELECT 1) UNION (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) EXCEPT (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) EXCEPT (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) INTERSECT (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) INTERSECT (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION ALL (SELECT 2) FOR UPDATE
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| (SELECT 1) UNION ALL (SELECT 2) FOR UPDATE
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of lock options and SELECT in brackets
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 5 CACHE 1000
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 2 CACHE 5
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 1 MAXVALUE 2 INCREMENT 2
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s MINVALUE 1 MAXVALUE 3 INCREMENT 2
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 10 MAXVALUE 5
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 0 MAXVALUE 2
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s CACHE 1 NOCACHE
< error
| 1225 HY000 Option 'CACHE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s CACHE 1 NOCACHE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CACHE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOCACHE NOCACHE
< error
| 1225 HY000 Option 'CACHE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOCACHE NOCACHE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CACHE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s START 1 START 2
< error
| 1225 HY000 Option 'START' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s START 1 START 2
| ; SELECT FROM
< error
| 1225 HY000 Option 'START' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 1 INCREMENT 2
< error
| 1225 HY000 Option 'INCREMENT' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s INCREMENT 1 INCREMENT 2
| ; SELECT FROM
< error
| 1225 HY000 Option 'INCREMENT' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s AS INT AS BIGINT
< error
| 1225 HY000 Option 'AS' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s AS INT AS BIGINT
| ; SELECT FROM
< error
| 1225 HY000 Option 'AS' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOMINVALUE MINVALUE 1
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOMINVALUE MINVALUE 1
| ; SELECT FROM
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s CYCLE CYCLE
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s CYCLE CYCLE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOCYCLE NOCYCLE
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s NOCYCLE NOCYCLE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| CREATE SEQUENCE s MAXVALUE 9223372036854775807 INCREMENT 1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s CACHE 9223372036854775807
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s CACHE 9223372036854775807
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s CACHE 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT MAXVALUE 1
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT MAXVALUE 1
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT UNSIGNED MAXVALUE 1
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS BIGINT UNSIGNED MAXVALUE 1
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS INT MAXVALUE 1
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS INT MAXVALUE 1
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT START WITH 200
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT START WITH 200
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT MINVALUE 5 MAXVALUE 4
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| CREATE SEQUENCE s AS TINYINT MINVALUE 5 MAXVALUE 4
| ; SELECT FROM
< error
| 4085 HY000 Sequence 'test.s' has out of range value for options
> exec ^
: string
| ALTER SEQUENCE s MINVALUE 1 MINVALUE 2
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| ALTER SEQUENCE s MINVALUE 1 MINVALUE 2
| ; SELECT FROM
< error
| 1225 HY000 Option 'MINVALUE' used twice in statement
> exec ^
: string
| ALTER SEQUENCE s CYCLE NOCYCLE
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| ALTER SEQUENCE s CYCLE NOCYCLE
| ; SELECT FROM
< error
| 1225 HY000 Option 'CYCLE' used twice in statement
> exec ^
: string
| ALTER SEQUENCE s RESTART RESTART
< error
| 1225 HY000 Option 'RESTART' used twice in statement
> exec ^
: string
| ALTER SEQUENCE s RESTART RESTART
| ; SELECT FROM
< error
| 1225 HY000 Option 'RESTART' used twice in statement
> exec ^
: string
| GRANT `select` ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT 'select' ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT tables ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT foo ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT READ ONLY TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| REVOKE READ ONLY FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| GRANT REPLICATION MASTER TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| REVOKE REPLICATION MASTER FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| GRANT READ ONLY ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT REPLICATION MASTER ON *.* TO u
< error
| 1064 42000 …near 'ON *.* TO u' at line 1
> exec ^
: string
| GRANT READ ONLY ADMIN TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT READ_ONLY ADMIN TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT CONNECTION ADMIN TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT BINLOG ADMIN TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT BINLOG MONITOR TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT REPLICA MONITOR TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT SLAVE MONITOR TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT DELETE HISTORY TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| GRANT READ TO u
< error
| 1064 42000 …near 'TO u' at line 1
> exec ^
: string
| REVOKE READ FROM u
< error
| 1064 42000 …near 'FROM u' at line 1
> exec ^
: string
| GRANT BINLOG TO u
< ok 0 0
> exec ^
: string
| GRANT CONNECTION TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION TO u
< ok 0 0
> exec ^
: string
| GRANT SLAVE TO u
< ok 0 0
> exec ^
: string
| GRANT READ_ONLY TO u
< ok 0 0
> exec ^
: string
| GRANT FEDERATED TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICA TO u
< ok 0 0
> exec ^
: string
| GRANT BINLOG ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT BINLOG MONITOR ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT BINLOG REPLAY ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CONNECTION ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT FEDERATED ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT READ_ONLY ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT READ ONLY ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION MASTER ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION SLAVE ADMIN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICA MONITOR ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SLAVE MONITOR ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SET USER ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT DELETE HISTORY ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SHOW CREATE ROUTINE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CREATE TABLESPACE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT EVENT ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT EXECUTE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT FILE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT RELOAD ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SHUTDOWN ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SUPER ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION CLIENT ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT REPLICATION SLAVE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SHOW DATABASES ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT SHOW VIEW ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT LOCK TABLES ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CREATE ROUTINE ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CREATE USER ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT CREATE TEMPORARY TABLES ON *.* TO u
< ok 0 0
> exec ^
: string
| GRANT USAGE ON *.* TO u
< ok 0 0
> exec ^
: string
| REVOKE r1 ON *.* FROM u
< error
| 1064 42000 …near 'ON *.* FROM u' at line 1
> exec ^
: string
| REVOKE foo ON *.* FROM u
< error
| 1064 42000 …near 'ON *.* FROM u' at line 1
> exec ^
: string
| REVOKE tables ON *.* FROM u
< error
| 1064 42000 …near 'ON *.* FROM u' at line 1
> exec ^
: string
| CREATE TABLE t (a INT AS (1) CONSTRAINT c CHECK (1))
< error
| 1064 42000 …near 'CHECK (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2) ON UPDATE NOW(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(2) ON UPDATE CURRENT_TIMESTAMP(1))
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATETIME(2) ON UPDATE CURRENT_TIMESTAMP(1))
| ; SELECT FROM
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE NOW(1+1))
< error
| 1064 42000 …near '+1))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE NOW() + 1)
< error
| 1064 42000 …near '+ 1)' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP ON UPDATE `NOW`())
< error
| 1064 42000 …near '`NOW`())' at line 1
> exec ^
: string
| SELECT SQL_CACHE 1
< ok 0 0
> exec ^
: string
| SELECT sql_cache
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SELECT sql_cache FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec ^
: string
| SELECT `sql_cache` FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT `sql_cache` FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SQL_CACHE a FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SQL_CACHE a FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SQL_NO_CACHE sql_cache FROM t
< error
| 1064 42000 …near 'FROM t' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'a'
< error
| 1064 42000 …near 'TABLESPACE ts ADD DATAFILE 'a'' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p TABLESPACE=x)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p TABLESPACE x)
< ok 0 0
> exec ^
: string
| SELECT 1 AS `Ø`
< ok 0 0
> exec ^
: string
| SELECT setval(d.s, 1)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT setval(d.s, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT setval(s, 1, 0)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT setval(s, 1, 0)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT setval(s, 1, 1, 5)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT setval(s, 1, 1, 5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT setval(s, 1, 2)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT setval(s, 1, 2)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT setval(s, 9223372036854775807)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT setval(s, 9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT setval()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT `setval`(s, 1)
< error
| 1630 42000 FUNCTION test.setval does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `setval`(s, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.setval(s, 1)
< error
| 1630 42000 FUNCTION d.setval does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT d.setval(s, 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nextval(s)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT nextval(s)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nextval(d.s)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT nextval(d.s)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT nextval(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT nextval()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT nextval(s, 1)
< error
| 1064 42000 …near ' 1)' at line 1
> exec ^
: string
| SELECT nextval(d.e.s)
< error
| 1064 42000 …near '.s)' at line 1
> exec ^
: string
| SELECT nextval(`a `.s)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT nextval(`a `.s)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| SELECT `nextval`(s)
< error
| 1630 42000 FUNCTION test.nextval does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `nextval`(s)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.nextval(s)
< error
| 1630 42000 FUNCTION d.nextval does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT d.nextval(s)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT lastval(s)
< error
| 4091 42S02 Unknown SEQUENCE: 's'
> exec ^
: string
| SELECT lastval(s)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT lastval(1)
< error
| 1064 42000 …near '1)' at line 1
> exec ^
: string
| SELECT lastval()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT lastval(`a `)
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT lastval(`a `)
| ; SELECT FROM
< error
| 1103 42000 Incorrect table name 'a '
> exec ^
: string
| SELECT 1 LIMIT a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SELECT 1 LIMIT a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SELECT 1 LIMIT a, 1
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| SELECT 1 LIMIT a, 1
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| UPDATE t SET a=1 LIMIT a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| UPDATE t SET a=1 LIMIT a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DELETE FROM t LIMIT a
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DELETE FROM t LIMIT a
| ; SELECT FROM
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| CREATE PROCEDURE p(a INT) SELECT 1 LIMIT a
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE a INT; SELECT 1 LIMIT a; END
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a INT) RETURNS INT BEGIN DECLARE v INT; SELECT 1 INTO v LIMIT a; RETURN v; END
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (a))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (a))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (a+1))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (a+1))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (a))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (a))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (a+1))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (a+1))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (a))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST COLUMNS (a) (PARTITION p VALUES IN (a))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((a,b)))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY LIST COLUMNS (a,b) (PARTITION p VALUES IN ((a,b)))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (a))
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (a))
| ; SELECT FROM
< error
| 1054 42S22 Unknown column 'a' in 'PARTITION BY'
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1 (PARTITION p VALUES IN (a))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1 (PARTITION p VALUES IN (a))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| ALTER TABLE t ADD COLUMN `` INT
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN `a ` INT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (`` INT, a INT)
< ok 0 0
> exec ^
: string
| DROP TABLESPACE ts
< error
| 1064 42000 …near 'TABLESPACE ts' at line 1
> exec ^
: string
| ALTER TABLESPACE ts ADD DATAFILE 'x'
< error
| 1064 42000 …near 'TABLESPACE ts ADD DATAFILE 'x'' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'x' ENGINE=InnoDB
< error
| 1064 42000 …near 'TABLESPACE ts ADD DATAFILE 'x' ENGINE=InnoDB' at line 1
> exec ^
: string
| CREATE UNDO TABLESPACE ts ADD DATAFILE 'x'
< error
| 1064 42000 …near 'UNDO TABLESPACE ts ADD DATAFILE 'x'' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP g ADD UNDOFILE 'x'
< error
| 1064 42000 …near 'LOGFILE GROUP g ADD UNDOFILE 'x'' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP g ADD UNDOFILE 'x'
< error
| 1064 42000 …near 'LOGFILE GROUP g ADD UNDOFILE 'x'' at line 1
> exec ^
: string
| DROP LOGFILE GROUP g
< error
| 1064 42000 …near 'LOGFILE GROUP g' at line 1
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER w OPTIONS (USER 'x')
< ok 0 0
> exec ^
: string
| ALTER SERVER s OPTIONS (USER 'x')
< ok 0 0
> exec ^
: string
| DROP SERVER s
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IN (SELECT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IN (SELECT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK ((SELECT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK ((SELECT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK ((SELECT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK ((SELECT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (EXISTS (SELECT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (EXISTS (SELECT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CONSTRAINT c CHECK ((SELECT 1)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CONSTRAINT c CHECK ((SELECT 1)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NEXT VALUE FOR s))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NEXT VALUE FOR s))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NEXTVAL(s)))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NEXTVAL(s)))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT CHECK (f(a)))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (f(a)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (d.f(a)))
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (d.f(a)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (ABS(a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (RAND()))
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (RAND()))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NOW()))
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (NOW()))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (@a))
< error
| 1901 HY000 Function or expression '@a' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (@a))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '@a' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (SELECT 1))
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (SELECT 1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT ((SELECT 1)))
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT ((SELECT 1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (f(1)))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (f(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (d.f(1)))
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (d.f(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (NEXT VALUE FOR s))
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (NEXT VALUE FOR s))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (NEXTVAL(s)))
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (NEXTVAL(s)))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (RAND()))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (@a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT @a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS ((SELECT 1)))
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS ((SELECT 1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (f(1)))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (f(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (d.f(1)))
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (d.f(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`d`.`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (NEXT VALUE FOR s))
< error
| 1901 HY000 Function or expression 'nextval()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (NEXT VALUE FOR s))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'nextval()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (RAND()))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS (@a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH ((SELECT 1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (f(a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY k (((SELECT 1))))
< error
| 1064 42000 …near '((SELECT 1))))' at line 1
> exec ^
: string
| ALTER TABLE t ADD CHECK ((SELECT 1))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ADD CHECK ((SELECT 1))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c CHECK (a IN (SELECT 1))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c CHECK (a IN (SELECT 1))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT CHECK ((SELECT 1))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT CHECK ((SELECT 1))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT ((SELECT 1))
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT ((SELECT 1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET NEW.a = 1
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET NEW.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = 1
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET NEW.a = 1
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET NEW.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET NEW.a = 1
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET NEW.a = 1
| ; SELECT FROM
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET NEW.a = 1
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET NEW.a = 1
| ; SELECT FROM
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET OLD.a = 1
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER DELETE ON t FOR EACH ROW SET OLD.a = 1
| ; SELECT FROM
< error
| 1362 HY000 Updating of OLD row is not allowed in trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = OLD.a
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET @x = OLD.a
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW SET @x = NEW.a
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = NEW.a
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = NEW.a
| ; SELECT FROM
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SET @x = OLD.a
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW SET @x = OLD.a + NEW.a
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT NEW.a INTO @x
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT OLD.a INTO @x
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT OLD.a INTO @x
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SELECT NEW.a INTO @x
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW SELECT NEW.a INTO @x
| ; SELECT FROM
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN SET NEW.a = 1; END
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW BEGIN SET NEW.a = 1; END
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW BEGIN SET NEW.a = 1; END
| ; SELECT FROM
< error
| 1362 HY000 Updating of NEW row is not allowed in after trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW BEGIN IF NEW.a THEN SET @x = 1; END IF; END
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE DELETE ON t FOR EACH ROW BEGIN IF NEW.a THEN SET @x = 1; END IF; END
| ; SELECT FROM
< error
| 1363 HY000 There is no NEW row in on DELETE trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE old INT; SET old.a = 1; END
< error
| 1105 HY000 'old' is not a row variable
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE old INT; SET old.a = 1; END
| ; SELECT FROM
< error
| 1105 HY000 'old' is not a row variable
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET `NEW`.a = 1
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW INSERT INTO u VALUES (OLD.a)
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW INSERT INTO u VALUES (OLD.a)
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW UPDATE u SET b = OLD.a
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW UPDATE u SET b = OLD.a
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT 1 FROM old
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT 1 FROM old
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT old.a FROM old
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT old.a FROM old
| ; SELECT FROM
< error
| 1363 HY000 There is no OLD row in on INSERT trigger
> exec ^
: string
| SET NEW.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'NEW'
> exec ^
: string
| SELECT NEW.a
< error
| 1109 42S02 Unknown table 'NEW' in SELECT
> exec ^
: string
| SELECT NEW.a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() SET NEW.a = 1
< error
| 4081 HY000 Unknown structured system variable or ROW routine variable 'NEW'
> exec ^
: string
| CREATE TABLE t (a INT, SERIAL)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT 1 FROM t 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 FROM t AS 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 FROM (SELECT 1) 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 FROM (SELECT 1) AS 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 FROM t, dual
< error
| 1064 42000 …near 'dual' at line 1
> exec ^
: string
| SELECT 1 FROM dual, t
< error
| 1064 42000 …near ' t' at line 1
> exec ^
: string
| SELECT 1 FROM dual JOIN t
< error
| 1064 42000 …near 'JOIN t' at line 1
> exec ^
: string
| SELECT 1 FROM t JOIN dual
< error
| 1064 42000 …near 'dual' at line 1
> exec ^
: string
| SELECT 1 FROM (dual)
< error
| 1064 42000 …near 'dual)' at line 1
> exec ^
: string
| SELECT 1 FROM dual a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| SELECT 1 FROM dual
< ok 0 0
> exec ^
: string
| SELECT 1 FROM DUAL WHERE 1
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1 INTO @a; RETURN 1; END
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1 UNION SELECT 2; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SELECT 1 UNION SELECT 2; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN (SELECT 1); RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN (SELECT 1); RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SHOW TABLES; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN SHOW TABLES; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN EXPLAIN SELECT 1; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN EXPLAIN SELECT 1; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN WITH c AS (SELECT 1) SELECT * FROM c; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN WITH c AS (SELECT 1) SELECT * FROM c; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN TABLE t; RETURN 1; END
< error
| 1064 42000 …near 'TABLE t; RETURN 1; END' at line 1
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN VALUES ROW(1); RETURN 1; END
< error
| 1064 42000 …near 'ROW(1); RETURN 1; END' at line 1
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN IF 1 THEN SELECT 1; END IF; RETURN 1; END
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN IF 1 THEN SELECT 1; END IF; RETURN 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a function
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT RETURN (SELECT 1)
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN DECLARE c CURSOR FOR SELECT 1; RETURN 1; END
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f() RETURNS INT BEGIN INSERT INTO t SELECT 1; RETURN 1; END
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() SELECT 1
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN SELECT 1; END
< ok 0 0
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT 1
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT 1
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SELECT 1 INTO @a
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN SELECT 1; END
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN SELECT 1; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SHOW TABLES
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SHOW TABLES
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXPLAIN SELECT 1
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW EXPLAIN SELECT 1
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN DECLARE c CURSOR FOR SELECT 1; END
< ok 0 0
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN IF 1 THEN (SELECT 1); END IF; END
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW BEGIN IF 1 THEN (SELECT 1); END IF; END
| ; SELECT FROM
< error
| 1415 0A000 Not allowed to return a result set from a trigger
> exec ^
: string
| CREATE TABLE t (a INT CHECK (u.f(a)))
< error
| 1901 HY000 Function or expression '`u`.`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (u.f(a)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`u`.`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (1 + f(a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CASE WHEN f(a) THEN 1 END))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CASE WHEN f(a) THEN 1 END))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (1 + (SELECT 1)))
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (1 + (SELECT 1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'select ...' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT f(1))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT f(1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT ABS(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (@a + 1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (@@sql_mode))
< error
| 1901 HY000 Function or expression '@@sql_mode' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (@@sql_mode))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '@@sql_mode' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (@@sql_mode))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS (@@sql_mode))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IS NULL))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (a IN (1,2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CAST(a AS CHAR) = 'x'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (IF(a, 1, 0)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (LEFT('a', 1)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (REPLACE('a','b','c')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (YEAR('2020-01-01')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (DATE('2020-01-01')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (TIMESTAMP('2020-01-01')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CHAR(65)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (INSERT('a',1,1,'b')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT CHECK (COUNT(*)))
< error
| 1901 HY000 Function or expression 'count()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (COUNT(*)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'count()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (SUM(a)))
< error
| 1901 HY000 Function or expression 'sum()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (SUM(a)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'sum()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (ROW_NUMBER() OVER ()))
< error
| 1901 HY000 Function or expression 'row_number()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (ROW_NUMBER() OVER ()))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'row_number()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CURRENT_USER))
< error
| 1901 HY000 Function or expression 'current_user()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CURRENT_USER))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'current_user()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CURRENT_TIMESTAMP))
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (CURRENT_TIMESTAMP))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'current_timestamp()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (CURRENT_USER))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS (CURRENT_USER))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (RAND()))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (f(a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (d.f(a)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (NOW()))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (@a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (NEXT VALUE FOR s))
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (NEXT VALUE FOR s))
| ; SELECT FROM
< error
| 1970 42000 CHECK does not support subqueries or stored functions
> exec ^
: string
| CREATE TABLE t (a INT, CHECK (COUNT(*)))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CHECK (RAND())
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CHECK (f(a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c CHECK (NOW())
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT CHECK (RAND())
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT CHECK (RAND())
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT DEFAULT (f(1))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT DEFAULT (f(1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT AS (f(1))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| ALTER TABLE t ADD COLUMN a INT AS (f(1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| ALTER TABLE t MODIFY a INT CHECK (RAND())
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| ALTER TABLE t MODIFY a INT CHECK (RAND())
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (f(1))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT (f(1))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT f(1)
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT f(1)
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NEXTVAL(s))
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT NEXTVAL(s))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (SLEEP(1)))
< error
| 1901 HY000 Function or expression 'sleep()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT (SLEEP(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'sleep()' cannot be used in the DEFAULT clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (`abs`(1)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS (`f`(1)))
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (`f`(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`f`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (`rand`()))
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT CHECK (`rand`()))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression 'rand()' cannot be used in the CHECK clause of `a`
> exec ^
: string
| CREATE PROCEDURE p() SELECT f(1)
< ok 0 0
> exec ^
: string
| CREATE FUNCTION g() RETURNS INT RETURN f(1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT AS (d.abs(1)))
< error
| 1901 HY000 Function or expression '`d`.`abs`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT AS (d.abs(1)))
| ; SELECT FROM
< error
| 1901 HY000 Function or expression '`d`.`abs`()' cannot be used in the GENERATED ALWAYS AS clause of `a`
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((a+1)))
< error
| 1064 42000 …near '(a+1)))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY k ((a)))
< error
| 1064 42000 …near '(a)))' at line 1
> exec ^
: string
| CREATE INDEX i ON t ((a+1))
< error
| 1064 42000 …near '(a+1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, UNIQUE KEY k ((a+1)))
< error
| 1064 42000 …near '(a+1)))' at line 1
> exec ^
: string
| ALTER TABLE t ADD KEY k ((a+1))
< error
| 1064 42000 …near '(a+1))' at line 1
> exec ^
: string
| SELECT 1 FROM t AS 'a' WHERE 1
< error
| 1064 42000 …near ''a' WHERE 1' at line 1
> exec ^
: string
| SELECT 1 FROM t a, u AS 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| UPDATE t 'a' SET x=1
< error
| 1064 42000 …near ''a' SET x=1' at line 1
> exec ^
: string
| UPDATE t AS 'a' SET x=1
< error
| 1064 42000 …near ''a' SET x=1' at line 1
> exec ^
: string
| DELETE FROM t 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| DELETE FROM t AS 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| DELETE a FROM t 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| INSERT INTO t AS 'a' VALUES (1)
< error
| 1064 42000 …near 'AS 'a' VALUES (1)' at line 1
> exec ^
: string
| SELECT 1 FROM t "a"
< error
| 1064 42000 …near '"a"' at line 1
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| SELECT 1 FROM JSON_TABLE('[]','$' COLUMNS (x INT PATH '$')) AS 'a'
< error
| 1064 42000 …near ''a'' at line 1
> exec ^
: string
| LOCK TABLES t 'a' READ
< error
| 1064 42000 …near ''a' READ' at line 1
> exec ^
: string
| LOCK TABLES t AS 'a' READ
< error
| 1064 42000 …near ''a' READ' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (f(a)) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ') (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (f(a)) (PARTITION p VALUES IN (1))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ') (PARTITION p VALUES IN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (d.f(a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (`f`(a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (ABS(a))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (RAND())
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a+f(a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (f(b)) (PARTITION p VALUES LESS THAN (1))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ') (PARTITION p VALUES LESS THAN (1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (f(1)))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (f(1)))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 + 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 - 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 * 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 / 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 / 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 DIV 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 % 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 MOD 2))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 | 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 | 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 & 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 & 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 ^ 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 ^ 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 << 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 << 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 >> 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1 >> 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (~1))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (~1))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (-1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (+1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 + 2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 OR 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 OR 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 AND 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 AND 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 XOR 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 XOR 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((NOT 1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((NOT 1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 = 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 = 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 < 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 < 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 IS NULL)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 IS NULL)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 IN (2))))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 IN (2))))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 BETWEEN 2 AND 3)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 BETWEEN 2 AND 3)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 LIKE 2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((1 LIKE 2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((CASE WHEN 1 THEN 2 END)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ((CASE WHEN 1 THEN 2 END)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (NULL))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN ('x'))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (@a))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (@@sql_mode))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (CAST(1 AS UNSIGNED)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p VALUES LESS THAN (CAST(1 AS UNSIGNED)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1 | 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1 | 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY LIST (a) (PARTITION p VALUES IN (1 + 2))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1 | 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES LESS THAN (1 | 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1 | 2))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| ALTER TABLE t ADD PARTITION (PARTITION p VALUES IN (1 | 2))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a | 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a + 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a & 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a OR 1)
< error
| 1064 42000 …near 'OR 1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH ((a OR 1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (@a)
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near ')' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (year('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (year('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (month('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (month('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (day('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (day('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (hour('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (hour('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (minute('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (minute('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (second('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (second('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (date('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (date('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (date('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (date('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (time('10:00:00')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (time('10:00:00')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (time('10:00:00')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (time('10:00:00')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (timestamp('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (timestamp('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (timestamp('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (timestamp('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (if(1,1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (if(1,1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (if(1,1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (if(1,1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (left('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (left('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (left('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (left('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (right('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (right('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (right('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (right('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (mod(5,2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (mod(5,2)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (extract(YEAR FROM '2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (extract(YEAR FROM '2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (char(65)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (char(65)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (char(65)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (char(65)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (truncate(1,1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (truncate(1,1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (truncate(1,1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (truncate(1,1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (insert('a',1,1,'x')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (insert('a',1,1,'x')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (insert('a',1,1,'x')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (insert('a',1,1,'x')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (repeat('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (repeat('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (repeat('a',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (repeat('a',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (replace('a','x','y')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (replace('a','x','y')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (replace('a','x','y')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (replace('a','x','y')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (cast(1 AS CHAR)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (cast(1 AS CHAR)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (cast(1 AS CHAR)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (cast(1 AS CHAR)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (convert(1,CHAR)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (convert(1,CHAR)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (convert(1,CHAR)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (convert(1,CHAR)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (date_add('2020-06-15',INTERVAL 1 DAY)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (date_add('2020-06-15',INTERVAL 1 DAY)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (date_add('2020-06-15',INTERVAL 1 DAY)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (date_add('2020-06-15',INTERVAL 1 DAY)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (timestampdiff(DAY,'2020-06-15','2020-06-16')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (timestampdiff(DAY,'2020-06-15','2020-06-16')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (timestampdiff(DAY,'2020-06-15','2020-06-16')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (timestampdiff(DAY,'2020-06-15','2020-06-16')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (adddate('2020-06-15',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (adddate('2020-06-15',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (adddate('2020-06-15',1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (adddate('2020-06-15',1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (utc_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (utc_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (current_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (current_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (now()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (now()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (rand()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (rand()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (unix_timestamp('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (unix_timestamp('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (unix_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (unix_timestamp()))
< error
| 1064 42000 Constant, random or timezone-dependent expressions in (sub)partitioning function are not allowed near '))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (to_days('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (to_days('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (to_seconds('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (to_seconds('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (coalesce(1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (coalesce(1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (coalesce(1)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (coalesce(1)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (ifnull(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (ifnull(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (ifnull(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (ifnull(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (least(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (least(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (least(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (least(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (greatest(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (greatest(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (greatest(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (greatest(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (nullif(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (nullif(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (nullif(1,2)))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (nullif(1,2)))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (week('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (week('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (week('2020-06-15')))
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (week('2020-06-15')))
| ; SELECT FROM
< error
| 1564 HY000 This partition function is not allowed
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (weekday('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (weekday('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (yearweek('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (yearweek('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (dayofyear('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (dayofyear('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (quarter('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (quarter('2020-06-15')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY RANGE (YEAR(a)) (PARTITION p VALUES LESS THAN (microsecond('10:00:00')))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATE) PARTITION BY LIST (YEAR(a)) (PARTITION p VALUES IN (microsecond('10:00:00')))
< ok 0 0
