-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestNumberTokenFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| CREATE TABLE t (a POINT SRID 4294967295)
< error
| 1064 42000 …line 1 column 23 near "POINT SRID 4294967295)" 
> exec ^
: string
| CREATE TABLE t (a POINT SRID 4294967296)
< error
| 1064 42000 …line 1 column 23 near "POINT SRID 4294967296)" 
> exec ^
: string
| CREATE TABLE t (a POINT SRID 0x10)
< error
| 1064 42000 …line 1 column 23 near "POINT SRID 0x10)" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x7fffffff
< error
| 1064 42000 …line 1 column 6 near "CHANGE MASTER TO MASTER_DELAY = 0x7fffffff" 
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x80000000
< error
| 1064 42000 …line 1 column 6 near "CHANGE MASTER TO MASTER_DELAY = 0x80000000" 
> exec ^
: string
| RESET MASTER TO 0x10
< error
| 1064 42000 …line 1 column 5 near "RESET MASTER TO 0x10" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 0x10
< error
| 1064 42000 …line 1 column 60 near "0x10" 
> exec ^
: string
| CREATE TABLE t (a CHAR(0x10))
< error
| 1064 42000 …line 1 column 27 near "0x10))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(1.))
< error
| 1064 42000 …line 1 column 25 near "1.))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(.5))
< error
| 1064 42000 …line 1 column 25 near ".5))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999.5))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999.5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999.5))
< error
| 1064 42000 …line 1 column 48 near "99999999999999999999.5))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967295))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967295))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775806))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775806))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551615.5))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551615.5))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999.5))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999.5))" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999K
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999K" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 4K
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 4K" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 4K
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 4K" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)
< error
| 1064 42000 …line 1 column 73 near "1.5)" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 0x10
< error
| 1064 42000 …line 1 column 45 near "0x10" 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 0x10
< error
| 1064 42000 …line 1 column 44 near "0x10" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0.5 BEGIN END; END
< error
| 1064 42000 …line 1 column 59 near "0.5 BEGIN END; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 LIMIT 0x10
< error
| 1064 42000 …line 1 column 19 near "0x10" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647M
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647M" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648M
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648M" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE 4M
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE 4M" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE 16384
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE 16384" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE 4M MAX_SIZE 1G EXTENT_SIZE 1M
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE 4M MAX_SIZE 1G EXTENT_SIZE 1M" 
> exec ^
: string
| ALTER TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
< error
| 1064 42000 …line 1 column 16 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M" 
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE 4M UNDO_BUFFER_SIZE 4M
< error
| 1064 42000 …line 1 column 14 near "LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE 4M UNDO_BUFFER_SIZE 4M" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 0x10)
< error
| 1064 42000 …line 1 column 74 near "0x10)" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END
< error
| 1064 42000 …line 1 column 46 near "CONDITION FOR 0; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0.9; END
< error
| 1064 42000 …line 1 column 46 near "CONDITION FOR 0.9; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.9; END
< error
| 1064 42000 …line 1 column 46 near "CONDITION FOR 1.9; END" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …line 1 column 29 near "1.5,0))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1.5))
< error
| 1064 42000 …line 1 column 32 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a TIME(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(1.5) CHARACTER SET utf8mb4)
< error
| 1064 42000 …line 1 column 26 near "1.5) CHARACTER SET utf8mb4)" 
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1.5)) RETURNS VARCHAR(2.5) RETURN a
< error
| 1064 42000 …line 1 column 15 near "FUNCTION f(a CHAR(1.5)) RETURNS VARCHAR(2.5) RETURN a" 
> exec ^
: string
| SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a YEAR(.5))
< error
| 1064 42000 …line 1 column 25 near ".5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(.5))
< error
| 1064 42000 …line 1 column 28 near ".5))" 
> exec ^
: string
| CREATE TABLE t (a TIME(.5))
< error
| 1064 42000 …line 1 column 25 near ".5))" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(.5))
< error
| 1064 42000 …line 1 column 28 near ".5))" 
> exec ^
: string
| CREATE TABLE t (a INT(.5))
< error
| 1064 42000 …line 1 column 24 near ".5))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(.5))
< error
| 1064 42000 …line 1 column 24 near ".5))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(.5))
< error
| 1064 42000 …line 1 column 27 near ".5))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(1.))
< error
| 1064 42000 …line 1 column 24 near "1.))" 
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1.5))
< error
| 1064 42000 …line 1 column 31 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …line 1 column 29 near "1.5,0))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …line 1 column 27 near "1.5,0))" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR .5 BEGIN END; END
< error
| 1064 42000 …line 1 column 58 near ".5 BEGIN END; END" 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 0 DAY
< ok 0 0
> exec ^
: string
| EXPLAIN FOR CONNECTION .5
< error
| 1064 42000 …line 1 column 25 near ".5" 
> exec ^
: string
| RESET MASTER TO .5
< error
| 1064 42000 …line 1 column 5 near "RESET MASTER TO .5" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 254 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 254 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'" 
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10
< error
| 1064 42000 …line 1 column 17 near "TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10" 
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967295
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967296
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65535
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65536
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
< error
| 1064 42000 …line 1 column 59 near "1.5 BEGIN END; END" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1.5)))
< error
| 1064 42000 …line 1 column 37 near "1.5)))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483647))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483648))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 46 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 46 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a BINARY(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483647))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483648))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551616))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a BINARY(99999999999999999999))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a BINARY(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(1.5))
< error
| 1064 42000 …line 1 column 28 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a BINARY(1e1))
< error
| 1064 42000 …line 1 column 28 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483647))
< error
| 1074 42000 Column length too big for column 'a' (max = 65535); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483648))
< error
| 1074 42000 Column length too big for column 'a' (max = 65535); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551616))
< error
| 1064 42000 …line 1 column 48 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(99999999999999999999))
< error
| 1064 42000 …line 1 column 48 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1.5))
< error
| 1064 42000 …line 1 column 31 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1e1))
< error
| 1064 42000 …line 1 column 31 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a BLOB(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a BLOB(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a BLOB(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a BLOB(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a TEXT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a TEXT(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a TEXT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a TEXT(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a BIT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(2147483647))
< error
| 1439 42000 Display width out of range for column 'a' (max = 64)
> exec ^
: string
| CREATE TABLE t (a BIT(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BIT(2147483648))
< error
| 1439 42000 Display width out of range for column 'a' (max = 64)
> exec ^
: string
| CREATE TABLE t (a BIT(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775808))
< error
| 3013 HY000 Invalid size for column 'a'.
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551616))
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a BIT(99999999999999999999))
< error
| 1064 42000 …line 1 column 42 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a BIT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a BIT(1e1))
< error
| 1064 42000 …line 1 column 25 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a INT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551616))
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
< error
| 1064 42000 …line 1 column 42 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a INT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a INT(1e1))
< error
| 1064 42000 …line 1 column 25 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a BIGINT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551616))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a BIGINT(99999999999999999999))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a BIGINT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(1.5))
< error
| 1064 42000 …line 1 column 28 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a BIGINT(1e1))
< error
| 1064 42000 …line 1 column 28 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483647))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483648))
< error
| 1074 42000 Column length too big for column 'a' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a NCHAR(1e1))
< error
| 1064 42000 …line 1 column 27 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 55 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 55 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1.5))
< error
| 1064 42000 …line 1 column 38 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1e1))
< error
| 1064 42000 …line 1 column 38 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551616))
< error
| 1064 42000 …line 1 column 51 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(99999999999999999999))
< error
| 1064 42000 …line 1 column 51 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1.5))
< error
| 1064 42000 …line 1 column 34 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1e1))
< error
| 1064 42000 …line 1 column 34 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616))
< error
| 1064 42000 …line 1 column 46 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999))
< error
| 1064 42000 …line 1 column 46 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647,0))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648,0))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …line 1 column 46 near "18446744073709551616,0))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …line 1 column 46 near "99999999999999999999,0))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …line 1 column 29 near "1.5,0))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1,0))
< error
| 1064 42000 …line 1 column 29 near "1e1,0))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483647))
< error
| 1425 42000 Too big scale 2147483647 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483648))
< error
| 1425 42000 Too big scale 2147483648 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775807))
< error
| 1425 42000 Too big scale 9223372036854775807 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,18446744073709551616))
< error
| 1064 42000 …line 1 column 49 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,99999999999999999999))
< error
| 1064 42000 …line 1 column 49 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1.5))
< error
| 1064 42000 …line 1 column 32 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1e1))
< error
| 1064 42000 …line 1 column 32 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775808))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551616))
< error
| 1064 42000 …line 1 column 46 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(99999999999999999999))
< error
| 1064 42000 …line 1 column 46 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a FIXED(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775808))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551616))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a FIXED(99999999999999999999))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a FIXED(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FIXED(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a FIXED(1e1))
< error
| 1064 42000 …line 1 column 27 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648))
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1e1))
< error
| 1064 42000 …line 1 column 27 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616,0))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616,0))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999,0))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999,0))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …line 1 column 27 near "1.5,0))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(1e1,0))
< error
| 1064 42000 …line 1 column 27 near "1e1,0))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483647))
< error
| 1425 42000 Too big scale 2147483647 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483648))
< error
| 1425 42000 Too big scale 2147483648 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775807))
< error
| 1425 42000 Too big scale 9223372036854775807 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,18446744073709551616))
< error
| 1064 42000 …line 1 column 47 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,99999999999999999999))
< error
| 1064 42000 …line 1 column 47 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1.5))
< error
| 1064 42000 …line 1 column 30 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1e1))
< error
| 1064 42000 …line 1 column 30 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483647,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483647,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483648,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483648,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775807,0))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775807,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775808,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775808,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551615,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551615,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551616,0))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616,0))" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(99999999999999999999,0))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999,0))" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1.5,0))
< error
| 1064 42000 …line 1 column 28 near "1.5,0))" 
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1e1,0))
< error
| 1064 42000 …line 1 column 28 near "1e1,0))" 
> exec ^
: string
| CREATE TABLE t (a TIME(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775807))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a TIME(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a TIME(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a TIME(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775807))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(18446744073709551616))
< error
| 1064 42000 …line 1 column 47 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME(99999999999999999999))
< error
| 1064 42000 …line 1 column 47 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(1.5))
< error
| 1064 42000 …line 1 column 30 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME(1e1))
< error
| 1064 42000 …line 1 column 30 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775807))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(18446744073709551616))
< error
| 1064 42000 …line 1 column 48 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(99999999999999999999))
< error
| 1064 42000 …line 1 column 48 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1.5))
< error
| 1064 42000 …line 1 column 31 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1e1))
< error
| 1064 42000 …line 1 column 31 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 42 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS CHAR(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS CHAR(1e1))
< error
| 1064 42000 …line 1 column 25 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(18446744073709551616))
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(99999999999999999999))
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS BINARY(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS BINARY(1e1))
< error
| 1064 42000 …line 1 column 27 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807))
< error
| 1426 42000 Too-big precision 9223372036854775807 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '`a`').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1064 42000 …line 1 column 28 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1e1))
< error
| 1064 42000 …line 1 column 28 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647,0))
< error
| 1426 42000 Too-big precision 2147483647 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648,0))
< error
| 1426 42000 Too-big precision 2147483648 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807,0))
< error
| 1426 42000 Too-big precision 9223372036854775807 specified for '`a`'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '`a`').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616,0))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999,0))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001,0))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5,0))
< error
| 1064 42000 …line 1 column 28 near "1.5,0))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1e1,0))
< error
| 1064 42000 …line 1 column 28 near "1e1,0))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483647))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '`a`').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483648))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '`a`').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775807))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '`a`').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551616))
< error
| 1064 42000 …line 1 column 48 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,99999999999999999999))
< error
| 1064 42000 …line 1 column 48 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1.5))
< error
| 1064 42000 …line 1 column 31 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1e1))
< error
| 1064 42000 …line 1 column 31 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS TIME(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551616))
< error
| 1064 42000 …line 1 column 42 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS TIME(99999999999999999999))
< error
| 1064 42000 …line 1 column 42 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS TIME(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS TIME(1.5))
< error
| 1064 42000 …line 1 column 25 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS TIME(1e1))
< error
| 1064 42000 …line 1 column 25 near "1e1))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483648))
< error
| 1426 42000 Too-big precision 2147483648 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775807))
< error
| 1426 42000 Too-big precision 9223372036854775807 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551616))
< error
| 1064 42000 …line 1 column 46 near "18446744073709551616))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(99999999999999999999))
< error
| 1064 42000 …line 1 column 46 near "99999999999999999999))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(1.5))
< error
| 1064 42000 …line 1 column 29 near "1.5))" 
> exec ^
: string
| SELECT CAST(a AS DATETIME(1e1))
< error
| 1064 42000 …line 1 column 29 near "1e1))" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| SELECT CONVERT(a, CHAR(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551616))
< error
| 1064 42000 …line 1 column 47 near "18446744073709551616))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(99999999999999999999))
< error
| 1064 42000 …line 1 column 47 near "99999999999999999999))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1.5))
< error
| 1064 42000 …line 1 column 30 near "1.5))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1e1))
< error
| 1064 42000 …line 1 column 30 near "1e1))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483647))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483648))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551616))
< error
| 1064 42000 …line 1 column 49 near "18446744073709551616))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(99999999999999999999))
< error
| 1064 42000 …line 1 column 49 near "99999999999999999999))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1.5))
< error
| 1064 42000 …line 1 column 32 near "1.5))" 
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1e1))
< error
| 1064 42000 …line 1 column 32 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775807))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(18446744073709551616))
< error
| 1064 42000 …line 1 column 61 near "18446744073709551616))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(99999999999999999999))
< error
| 1064 42000 …line 1 column 61 near "99999999999999999999))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1.5))
< error
| 1064 42000 …line 1 column 44 near "1.5))" 
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1e1))
< error
| 1064 42000 …line 1 column 44 near "1e1))" 
> exec ^
: string
| CREATE INDEX i ON t (a(1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(2147483647))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(2147483648))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775807))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775808))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551615))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551616))
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616))" 
> exec ^
: string
| CREATE INDEX i ON t (a(99999999999999999999))
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999))" 
> exec ^
: string
| CREATE INDEX i ON t (a(000000000000000000000001))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a(1.5))
< error
| 1064 42000 …line 1 column 26 near "1.5))" 
> exec ^
: string
| CREATE INDEX i ON t (a(1e1))
< error
| 1064 42000 …line 1 column 26 near "1e1))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(2147483647)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(2147483648)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775807)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775808)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(18446744073709551615)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(18446744073709551616)))
< error
| 1064 42000 …line 1 column 54 near "18446744073709551616)))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(99999999999999999999)))
< error
| 1064 42000 …line 1 column 54 near "99999999999999999999)))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(000000000000000000000001)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1.5)))
< error
| 1064 42000 …line 1 column 37 near "1.5)))" 
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1e1)))
< error
| 1064 42000 …line 1 column 37 near "1e1)))" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 …line 1 column 61 near "18446744073709551616" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 …line 1 column 61 near "99999999999999999999" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1.5
< error
| 1064 42000 …line 1 column 44 near "1.5" 
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1e1
< error
| 1064 42000 …line 1 column 44 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 18446744073709551616)
< error
| 1064 42000 …line 1 column 66 near "18446744073709551616)" 
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 99999999999999999999)
< error
| 1064 42000 …line 1 column 66 near "99999999999999999999)" 
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1.5)
< error
| 1064 42000 …line 1 column 49 near "1.5)" 
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1e1)
< error
| 1064 42000 …line 1 column 49 near "1e1)" 
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 18446744073709551616
< error
| 1064 42000 …line 1 column 54 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 99999999999999999999
< error
| 1064 42000 …line 1 column 54 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1.5
< error
| 1064 42000 …line 1 column 37 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1e1
< error
| 1064 42000 …line 1 column 37 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 18446744073709551616
< error
| 1064 42000 …line 1 column 61 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 99999999999999999999
< error
| 1064 42000 …line 1 column 61 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1.5
< error
| 1064 42000 …line 1 column 44 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1e1
< error
| 1064 42000 …line 1 column 44 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 18446744073709551616
< error
| 1064 42000 …line 1 column 54 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 99999999999999999999
< error
| 1064 42000 …line 1 column 54 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1.5
< error
| 1064 42000 …line 1 column 37 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1e1
< error
| 1064 42000 …line 1 column 37 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 18446744073709551616
< error
| 1064 42000 …line 1 column 54 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 99999999999999999999
< error
| 1064 42000 …line 1 column 54 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1.5
< error
| 1064 42000 …line 1 column 37 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1e1
< error
| 1064 42000 …line 1 column 37 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 18446744073709551616
< error
| 1064 42000 …line 1 column 55 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 99999999999999999999
< error
| 1064 42000 …line 1 column 55 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1.5
< error
| 1064 42000 …line 1 column 38 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1e1
< error
| 1064 42000 …line 1 column 38 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 18446744073709551616
< error
| 1064 42000 …line 1 column 62 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 99999999999999999999
< error
| 1064 42000 …line 1 column 62 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1.5
< error
| 1064 42000 …line 1 column 45 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1e1
< error
| 1064 42000 …line 1 column 45 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551616
< error
| 1064 42000 …line 1 column 64 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 99999999999999999999
< error
| 1064 42000 …line 1 column 64 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1.5
< error
| 1064 42000 …line 1 column 47 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1e1
< error
| 1064 42000 …line 1 column 47 near "1e1" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551616
< error
| 1064 42000 …line 1 column 51 near "18446744073709551616" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 99999999999999999999
< error
| 1064 42000 …line 1 column 51 near "99999999999999999999" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1.5
< error
| 1064 42000 …line 1 column 34 near "1.5" 
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1e1
< error
| 1064 42000 …line 1 column 34 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 18446744073709551616
< error
| 1064 42000 …line 1 column 76 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 99999999999999999999
< error
| 1064 42000 …line 1 column 76 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1.5
< error
| 1064 42000 …line 1 column 59 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1e1
< error
| 1064 42000 …line 1 column 59 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2147483647 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2147483648 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 9223372036854775807 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 9223372036854775808 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 18446744073709551615 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 18446744073709551616 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …line 1 column 112 near "18446744073709551616 (PARTITION p0 VALUES LESS THAN (10))" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …line 1 column 112 near "99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 000000000000000000000001 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1.5 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …line 1 column 95 near "1.5 (PARTITION p0 VALUES LESS THAN (10))" 
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1e1 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …line 1 column 95 near "1e1 (PARTITION p0 VALUES LESS THAN (10))" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 18446744073709551616)
< error
| 1064 42000 …line 1 column 90 near "18446744073709551616)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 99999999999999999999)
< error
| 1064 42000 …line 1 column 90 near "99999999999999999999)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)
< error
| 1064 42000 …line 1 column 73 near "1.5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1e1)
< error
| 1064 42000 …line 1 column 73 near "1e1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 18446744073709551616)
< error
| 1064 42000 …line 1 column 90 near "18446744073709551616)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 99999999999999999999)
< error
| 1064 42000 …line 1 column 90 near "99999999999999999999)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 1.5)
< error
| 1064 42000 …line 1 column 73 near "1.5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 1e1)
< error
| 1064 42000 …line 1 column 73 near "1e1)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 18446744073709551616)
< error
| 1064 42000 …line 1 column 91 near "18446744073709551616)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 99999999999999999999)
< error
| 1064 42000 …line 1 column 91 near "99999999999999999999)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 1.5)
< error
| 1064 42000 …line 1 column 74 near "1.5)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 1e1)
< error
| 1064 42000 …line 1 column 74 near "1e1)" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 18446744073709551616
< error
| 1064 42000 …line 1 column 59 near "18446744073709551616" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 99999999999999999999
< error
| 1064 42000 …line 1 column 59 near "99999999999999999999" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1.5
< error
| 1064 42000 …line 1 column 42 near "1.5" 
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1e1
< error
| 1064 42000 …line 1 column 42 near "1e1" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 18446744073709551616
< error
| 1064 42000 …line 1 column 53 near "18446744073709551616" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 99999999999999999999
< error
| 1064 42000 …line 1 column 53 near "99999999999999999999" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1.5
< error
| 1064 42000 …line 1 column 36 near "1.5" 
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1e1
< error
| 1064 42000 …line 1 column 36 near "1e1" 
> exec ^
: string
| SELECT 1 LIMIT 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 2147483647
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 2147483648
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 9223372036854775807
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 9223372036854775808
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551616
< error
| 1064 42000 …line 1 column 35 near "18446744073709551616" 
> exec ^
: string
| SELECT 1 LIMIT 99999999999999999999
< error
| 1064 42000 …line 1 column 35 near "99999999999999999999" 
> exec ^
: string
| SELECT 1 LIMIT 000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1.5
< error
| 1064 42000 …line 1 column 18 near "1.5" 
> exec ^
: string
| SELECT 1 LIMIT 1e1
< error
| 1064 42000 …line 1 column 18 near "1e1" 
> exec ^
: string
| SELECT 1 LIMIT 1, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 2147483647, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 2147483648, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 9223372036854775807, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 9223372036854775808, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551615, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 18446744073709551616, 1
< error
| 1064 42000 …line 1 column 35 near "18446744073709551616, 1" 
> exec ^
: string
| SELECT 1 LIMIT 99999999999999999999, 1
< error
| 1064 42000 …line 1 column 35 near "99999999999999999999, 1" 
> exec ^
: string
| SELECT 1 LIMIT 000000000000000000000001, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1.5, 1
< error
| 1064 42000 …line 1 column 18 near "1.5, 1" 
> exec ^
: string
| SELECT 1 LIMIT 1e1, 1
< error
| 1064 42000 …line 1 column 18 near "1e1, 1" 
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 2147483647
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 2147483648
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 9223372036854775807
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 9223372036854775808
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 18446744073709551616
< error
| 1064 42000 …line 1 column 44 near "18446744073709551616" 
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 99999999999999999999
< error
| 1064 42000 …line 1 column 44 near "99999999999999999999" 
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 1.5
< error
| 1064 42000 …line 1 column 27 near "1.5" 
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 1e1
< error
| 1064 42000 …line 1 column 27 near "1e1" 
> exec ^
: string
| DELETE FROM t LIMIT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551616
< error
| 1064 42000 …line 1 column 40 near "18446744073709551616" 
> exec ^
: string
| DELETE FROM t LIMIT 99999999999999999999
< error
| 1064 42000 …line 1 column 40 near "99999999999999999999" 
> exec ^
: string
| DELETE FROM t LIMIT 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t LIMIT 1.5
< error
| 1064 42000 …line 1 column 23 near "1.5" 
> exec ^
: string
| DELETE FROM t LIMIT 1e1
< error
| 1064 42000 …line 1 column 23 near "1e1" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551616
< error
| 1064 42000 …line 1 column 45 near "18446744073709551616" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 99999999999999999999
< error
| 1064 42000 …line 1 column 45 near "99999999999999999999" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1.5
< error
| 1064 42000 …line 1 column 28 near "1.5" 
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1e1
< error
| 1064 42000 …line 1 column 28 near "1e1" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551616 LINES
< error
| 1064 42000 …line 1 column 61 near "18446744073709551616 LINES" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 99999999999999999999 LINES
< error
| 1064 42000 …line 1 column 61 near "99999999999999999999 LINES" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 000000000000000000000001 LINES
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 000000000000000000000001 LINES
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1.5 LINES
< error
| 1064 42000 …line 1 column 44 near "1.5 LINES" 
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1e1 LINES
< error
| 1064 42000 …line 1 column 44 near "1e1 LINES" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483647
< error
| 1094 HY000 Unknown thread id: 2147483647
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483648
< error
| 1094 HY000 Unknown thread id: 2147483648
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775807
< error
| 1094 HY000 Unknown thread id: 9223372036854775807
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775808
< error
| 1094 HY000 Unknown thread id: 9223372036854775808
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< error
| 1094 HY000 Unknown thread id: 18446744073709551615
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551616
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 99999999999999999999
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 000000000000000000000001
< error
| 1094 HY000 Unknown thread id: 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1.5
< error
| 1064 42000 …line 1 column 26 near "1.5" 
> exec ^
: string
| EXPLAIN FOR CONNECTION 1e1
< error
| 1064 42000 …line 1 column 26 near "1e1" 
> exec ^
: string
| KILL 1
< ok 0 0
> exec ^
: string
| KILL 2147483647
< ok 0 0
> exec ^
: string
| KILL 2147483648
< ok 0 0
> exec ^
: string
| KILL 9223372036854775807
< ok 0 0
> exec ^
: string
| KILL 9223372036854775808
< ok 0 0
> exec ^
: string
| KILL 18446744073709551615
< ok 0 0
> exec ^
: string
| KILL 18446744073709551616
< error
| 1064 42000 …line 1 column 25 near "18446744073709551616" 
> exec ^
: string
| KILL 99999999999999999999
< error
| 1064 42000 …line 1 column 25 near "99999999999999999999" 
> exec ^
: string
| KILL 000000000000000000000001
< ok 0 0
> exec ^
: string
| KILL 1.5
< error
| 1064 42000 …line 1 column 8 near "1.5" 
> exec ^
: string
| KILL 1e1
< error
| 1064 42000 …line 1 column 8 near "1e1" 
> exec ^
: string
| KILL QUERY 1
< ok 0 0
> exec ^
: string
| KILL QUERY 2147483647
< ok 0 0
> exec ^
: string
| KILL QUERY 2147483648
< ok 0 0
> exec ^
: string
| KILL QUERY 9223372036854775807
< ok 0 0
> exec ^
: string
| KILL QUERY 9223372036854775808
< ok 0 0
> exec ^
: string
| KILL QUERY 18446744073709551615
< ok 0 0
> exec ^
: string
| KILL QUERY 18446744073709551616
< error
| 1064 42000 …line 1 column 31 near "18446744073709551616" 
> exec ^
: string
| KILL QUERY 99999999999999999999
< error
| 1064 42000 …line 1 column 31 near "99999999999999999999" 
> exec ^
: string
| KILL QUERY 000000000000000000000001
< ok 0 0
> exec ^
: string
| KILL QUERY 1.5
< error
| 1064 42000 …line 1 column 14 near "1.5" 
> exec ^
: string
| KILL QUERY 1e1
< error
| 1064 42000 …line 1 column 14 near "1e1" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483647 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483647 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483648 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483648 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775807 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775807 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775808 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775808 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551615 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551615 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551616 BEGIN END; END
< error
| 1064 42000 …line 1 column 76 near "18446744073709551616 BEGIN END; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 99999999999999999999 BEGIN END; END
< error
| 1064 42000 …line 1 column 76 near "99999999999999999999 BEGIN END; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 000000000000000000000001 BEGIN END; END
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 000000000000000000000001 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
< error
| 1064 42000 …line 1 column 59 near "1.5 BEGIN END; END" 
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1e1 BEGIN END; END
< error
| 1064 42000 …line 1 column 59 near "1e1 BEGIN END; END" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483648 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483648 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775807 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775807 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775808 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775808 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551615 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551615 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551616 BUCKETS
< error
| 1064 42000 …line 1 column 71 near "BUCKETS" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 99999999999999999999 BUCKETS
< error
| 1064 42000 …line 1 column 71 near "BUCKETS" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 000000000000000000000001 BUCKETS
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 000000000000000000000001 BUCKETS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1.5 BUCKETS
< error
| 1064 42000 …line 1 column 54 near "BUCKETS" 
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1e1 BUCKETS
< error
| 1064 42000 …line 1 column 54 near "BUCKETS" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 18446744073709551616
< error
| 1064 42000 …line 1 column 52 near "18446744073709551616" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 99999999999999999999
< error
| 1064 42000 …line 1 column 52 near "99999999999999999999" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1.5
< error
| 1064 42000 …line 1 column 35 near "1.5" 
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1e1
< error
| 1064 42000 …line 1 column 35 near "1e1" 
> exec ^
: string
| SELECT 1
< ok 0 0
> exec ^
: string
| SELECT 2147483647
< ok 0 0
> exec ^
: string
| SELECT 2147483648
< ok 0 0
> exec ^
: string
| SELECT 9223372036854775807
< ok 0 0
> exec ^
: string
| SELECT 9223372036854775808
< ok 0 0
> exec ^
: string
| SELECT 18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT 18446744073709551616
< ok 0 0
> exec ^
: string
| SELECT 99999999999999999999
< ok 0 0
> exec ^
: string
| SELECT 000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT 1.5
< ok 0 0
> exec ^
: string
| SELECT 1e1
< ok 0 0
> exec ^
: string
| SELECT -1
< ok 0 0
> exec ^
: string
| SELECT -2147483647
< ok 0 0
> exec ^
: string
| SELECT -2147483648
< ok 0 0
> exec ^
: string
| SELECT -9223372036854775807
< ok 0 0
> exec ^
: string
| SELECT -9223372036854775808
< ok 0 0
> exec ^
: string
| SELECT -18446744073709551615
< ok 0 0
> exec ^
: string
| SELECT -18446744073709551616
< ok 0 0
> exec ^
: string
| SELECT -99999999999999999999
< ok 0 0
> exec ^
: string
| SELECT -000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT -1.5
< ok 0 0
> exec ^
: string
| SELECT -1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 18446744073709551616)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 99999999999999999999)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1.5)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1e1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -9223372036854775807)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -9223372036854775808)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -18446744073709551615)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -18446744073709551616)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -99999999999999999999)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1.5)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1e1)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551616
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 99999999999999999999
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1.5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1e1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1e1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551616
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -99999999999999999999
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1.5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1e1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1e1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (2147483648))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (9223372036854775807))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (9223372036854775808))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (18446744073709551615))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (18446744073709551616))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (99999999999999999999))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1e1))
< ok 0 0
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483647 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483647 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483648 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483648 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775807 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775807 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775808 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775808 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551615 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551615 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551616 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551616 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 99999999999999999999 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 99999999999999999999 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 000000000000000000000001 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 000000000000000000000001 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1.5 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1.5 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1e1 PRECEDING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1e1 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483647 PRECEDING AND 2147483647 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483647 PRECEDING AND 2147483647 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483648 PRECEDING AND 2147483648 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483648 PRECEDING AND 2147483648 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775807 PRECEDING AND 9223372036854775807 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775807 PRECEDING AND 9223372036854775807 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775808 PRECEDING AND 9223372036854775808 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775808 PRECEDING AND 9223372036854775808 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551615 PRECEDING AND 18446744073709551615 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551615 PRECEDING AND 18446744073709551615 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551616 PRECEDING AND 18446744073709551616 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551616 PRECEDING AND 18446744073709551616 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 99999999999999999999 PRECEDING AND 99999999999999999999 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 99999999999999999999 PRECEDING AND 99999999999999999999 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 000000000000000000000001 PRECEDING AND 000000000000000000000001 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 000000000000000000000001 PRECEDING AND 000000000000000000000001 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1.5 PRECEDING AND 1.5 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1.5 PRECEDING AND 1.5 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1e1 PRECEDING AND 1e1 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1e1 PRECEDING AND 1e1 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 1 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 1 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 2147483647 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 2147483647 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 2147483648 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 2147483648 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 9223372036854775807 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 9223372036854775807 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 9223372036854775808 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 9223372036854775808 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 18446744073709551615 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 18446744073709551615 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 18446744073709551616 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 18446744073709551616 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 99999999999999999999 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 99999999999999999999 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 000000000000000000000001 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 000000000000000000000001 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 1.5 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 1.5 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT INTERVAL 1e1 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 1e1 DAY + a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483647
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483648
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775808
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551615
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551616
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 99999999999999999999
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 000000000000000000000001
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 1.5
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t ORDER BY 1e1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1e1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT CHAR(1 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(2147483647 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(2147483648 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(9223372036854775807 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(9223372036854775808 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(18446744073709551615 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(18446744073709551616 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(99999999999999999999 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(000000000000000000000001 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(1.5 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT CHAR(1e1 USING utf8mb4)
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(1) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(2147483647) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(2147483648) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(9223372036854775807) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(9223372036854775808) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(18446744073709551615) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(18446744073709551616) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(99999999999999999999) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(000000000000000000000001) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(1.5) */ 1
< ok 0 0
> exec ^
: string
| SELECT /*+ MAX_EXECUTION_TIME(1e1) */ 1
< ok 0 0
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 2147483647
< error
| 1292 22007 Incorrect time value: '2147483647'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 2147483647
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 2147483648
< error
| 1292 22007 Incorrect time value: '2147483648'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 2147483648
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 9223372036854775807
< error
| 1292 22007 Incorrect time value: '9223372036854775807'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 9223372036854775808
< error
| 1105 HY000 cannot convert datum from unsigned bigint to type timestamp
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 18446744073709551615
< error
| 1105 HY000 cannot convert datum from unsigned bigint to type timestamp
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 18446744073709551616
< error
| 1292 22007 Incorrect datetime value: '18446744073709551616'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 99999999999999999999
< error
| 1292 22007 Incorrect datetime value: '99999999999999999999'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 000000000000000000000001
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1.5
< error
| 1292 22007 Incorrect datetime value: '1.5'
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1.5
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1e1
< error
| 1105 HY000 cannot convert datum from double to type timestamp
> exec ^
: string
| SELECT * FROM t AS OF TIMESTAMP 1e1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SET @a = 1
< ok 0 0
> exec ^
: string
| SET @a = 2147483647
< ok 0 0
> exec ^
: string
| SET @a = 2147483648
< ok 0 0
> exec ^
: string
| SET @a = 9223372036854775807
< ok 0 0
> exec ^
: string
| SET @a = 9223372036854775808
< ok 0 0
> exec ^
: string
| SET @a = 18446744073709551615
< ok 0 0
> exec ^
: string
| SET @a = 18446744073709551616
< ok 0 0
> exec ^
: string
| SET @a = 99999999999999999999
< ok 0 0
> exec ^
: string
| SET @a = 000000000000000000000001
< ok 0 0
> exec ^
: string
| SET @a = 1.5
< ok 0 0
> exec ^
: string
| SET @a = 1e1
< ok 0 0
> exec ^
: string
| SELECT a FROM t WHERE a IN (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483647)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483647)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483648)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483648)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775807)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775808)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551615)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551616)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (99999999999999999999)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (000000000000000000000001)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (1.5)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1.5)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT a FROM t WHERE a IN (1e1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1e1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (2147483647)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (2147483647)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (2147483648)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (2147483648)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775807)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775808)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551615)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551616)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (99999999999999999999)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (000000000000000000000001)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (1.5)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1.5)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT INTO t VALUES (1e1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1e1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(1)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(2147483647)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(2147483648)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(2147483648)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(9223372036854775807)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(9223372036854775808)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(18446744073709551615)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(18446744073709551616)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(99999999999999999999)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(000000000000000000000001)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(1.5)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(1.5)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| CALL p(1e1)
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| CALL p(1e1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW PROFILE FOR QUERY 1
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 2147483647
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 2147483648
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 9223372036854775807
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 9223372036854775808
< error
| 1064 42000 …line 1 column 42 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| SHOW PROFILE FOR QUERY 18446744073709551615
< error
| 1064 42000 …line 1 column 43 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| SHOW PROFILE FOR QUERY 18446744073709551616
< error
| 1064 42000 …line 1 column 43 near "18446744073709551616" 
> exec ^
: string
| SHOW PROFILE FOR QUERY 99999999999999999999
< error
| 1064 42000 …line 1 column 43 near "99999999999999999999" 
> exec ^
: string
| SHOW PROFILE FOR QUERY 000000000000000000000001
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 1.5
< error
| 1064 42000 …line 1 column 26 near "1.5" 
> exec ^
: string
| SHOW PROFILE FOR QUERY 1e1
< error
| 1064 42000 …line 1 column 26 near "1e1" 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483647 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483648 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775807 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775808 DAY
< error
| 1064 42000 …line 1 column 62 near "DAY"9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 18446744073709551615 DAY
< error
| 1064 42000 …line 1 column 63 near "DAY"18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 18446744073709551616 DAY
< error
| 1064 42000 …line 1 column 59 near "18446744073709551616 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 99999999999999999999 DAY
< error
| 1064 42000 …line 1 column 59 near "99999999999999999999 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 000000000000000000000001 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1.5 DAY
< error
| 1064 42000 …line 1 column 42 near "1.5 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1e1 DAY
< error
| 1064 42000 …line 1 column 42 near "1e1 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 1
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 18446744073709551616
< error
| 1064 42000 …line 1 column 51 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 99999999999999999999
< error
| 1064 42000 …line 1 column 51 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 1.5
< error
| 1064 42000 …line 1 column 34 near "1.5" 
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 1e1
< error
| 1064 42000 …line 1 column 34 near "1e1" 
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 1 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 2147483647 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 2147483648 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 9223372036854775807 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 9223372036854775808 DAY
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 9223372036854775808 DAY
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 18446744073709551615 DAY
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 18446744073709551615 DAY
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 18446744073709551616 DAY
< error
| 1064 42000 …line 1 column 58 near "18446744073709551616 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 99999999999999999999 DAY
< error
| 1064 42000 …line 1 column 58 near "99999999999999999999 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 000000000000000000000001 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 1.5 DAY
< error
| 1064 42000 …line 1 column 41 near "1.5 DAY" 
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 1e1 DAY
< error
| 1064 42000 …line 1 column 41 near "1e1 DAY" 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775808
< error
| 1064 42000 …line 1 column 55 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 18446744073709551615
< error
| 1064 42000 …line 1 column 56 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 18446744073709551616
< error
| 1064 42000 …line 1 column 56 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 99999999999999999999
< error
| 1064 42000 …line 1 column 56 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1.5
< error
| 1064 42000 …line 1 column 39 near "1.5" 
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1e1
< error
| 1064 42000 …line 1 column 39 near "1e1" 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775808
< error
| 1064 42000 …line 1 column 52 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 18446744073709551615
< error
| 1064 42000 …line 1 column 53 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 18446744073709551616
< error
| 1064 42000 …line 1 column 53 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 99999999999999999999
< error
| 1064 42000 …line 1 column 53 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1.5
< error
| 1064 42000 …line 1 column 36 near "1.5" 
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1e1
< error
| 1064 42000 …line 1 column 36 near "1e1" 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 9223372036854775808
< error
| 1064 42000 …line 1 column 59 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551615
< error
| 1064 42000 …line 1 column 60 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 9223372036854775808
< error
| 1064 42000 …line 1 column 59 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 18446744073709551615
< error
| 1064 42000 …line 1 column 60 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 9223372036854775808
< error
| 1064 42000 …line 1 column 63 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 18446744073709551615
< error
| 1064 42000 …line 1 column 64 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 18446744073709551616
< error
| 1064 42000 …line 1 column 64 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 99999999999999999999
< error
| 1064 42000 …line 1 column 64 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1.5
< error
| 1064 42000 …line 1 column 47 near "1.5" 
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1e1
< error
| 1064 42000 …line 1 column 47 near "1e1" 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 2147483648
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 9223372036854775808
< error
| 1064 42000 …line 1 column 59 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 18446744073709551615
< error
| 1064 42000 …line 1 column 60 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 18446744073709551616
< error
| 1064 42000 …line 1 column 60 near "18446744073709551616" 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 99999999999999999999
< error
| 1064 42000 …line 1 column 60 near "99999999999999999999" 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1e1
< error
| 1064 42000 …line 1 column 43 near "1e1" 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 1
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 2147483647
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 2147483648
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 9223372036854775808
< error
| 1064 42000 …line 1 column 58 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551615
< error
| 1064 42000 …line 1 column 59 near ""18446744073709551615 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551616
< error
| 1064 42000 …line 1 column 59 near "18446744073709551616" 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 99999999999999999999
< error
| 1064 42000 …line 1 column 59 near "99999999999999999999" 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 1.5
< error
| 1064 42000 …line 1 column 42 near "1.5" 
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 1e1
< error
| 1064 42000 …line 1 column 42 near "1e1" 
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 1
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 2147483647
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 2147483648
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 18446744073709551616
< error
| 1064 42000 …line 1 column 50 near "18446744073709551616" 
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 99999999999999999999
< error
| 1064 42000 …line 1 column 50 near "99999999999999999999" 
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 1.5
< error
| 1064 42000 …line 1 column 33 near "1.5" 
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 1e1
< error
| 1064 42000 …line 1 column 33 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483647
< error
| 1064 42000 …line 1 column 53 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483648
< error
| 1064 42000 …line 1 column 53 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775807
< error
| 1064 42000 …line 1 column 62 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775808
< error
| 1064 42000 …line 1 column 62 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 18446744073709551615
< error
| 1064 42000 …line 1 column 63 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 18446744073709551616
< error
| 1064 42000 …line 1 column 63 near "18446744073709551616" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 99999999999999999999
< error
| 1064 42000 …line 1 column 63 near "99999999999999999999" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1.5
< error
| 1064 42000 …line 1 column 46 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1e1
< error
| 1064 42000 …line 1 column 46 near "1e1" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483647 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483647 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483648 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483648 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775807 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775807 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775808 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775808 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 18446744073709551615 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 18446744073709551615 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 18446744073709551616 (a)
< error
| 1064 42000 …line 1 column 72 near "18446744073709551616 (a)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 99999999999999999999 (a)
< error
| 1064 42000 …line 1 column 72 near "99999999999999999999 (a)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 000000000000000000000001 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1.5 (a)
< error
| 1064 42000 …line 1 column 55 near "1.5 (a)" 
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1e1 (a)
< error
| 1064 42000 …line 1 column 55 near "1e1 (a)" 
