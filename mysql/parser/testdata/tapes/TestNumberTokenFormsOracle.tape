-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestNumberTokenFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a POINT SRID 4294967295)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 4294967296)
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 4294967296)
| ; SELECT FROM
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294967
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294967.5
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294967.5
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294968
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 4294968
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 0
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 0.0001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 99999999999999999999.5
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 99999999999999999999.5
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1e30
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1e30
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_HEARTBEAT_PERIOD = 4294968
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_HEARTBEAT_PERIOD = 4294968
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_HEARTBEAT_PERIOD = 99999999999999999999
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_HEARTBEAT_PERIOD = 99999999999999999999
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483647.9
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648.0
< error
| 1729 HY000 The requested value 2147483648.0 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648.0
| ; SELECT FROM
< error
| 1729 HY000 The requested value 2147483648.0 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 1e30
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x80000000
< error
| 1729 HY000 The requested value 0x80000000 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x80000000
| ; SELECT FROM
< error
| 1729 HY000 The requested value 0x80000000 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_DELAY = 2147483648
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_DELAY = 2147483648
| ; SELECT FROM
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_DELAY = 99999999999999999999
< error
| 1729 HY000 The requested value 99999999999999999999 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_DELAY = 99999999999999999999
| ; SELECT FROM
< error
| 1729 HY000 The requested value 99999999999999999999 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| RESET MASTER TO 0
< error
| 3567 HY000 The requested value '0' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 0
| ; SELECT FROM
< error
| 3567 HY000 The requested value '0' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2000000000
< ok 0 0
> exec ^
: string
| RESET MASTER TO 2000000001
< error
| 3567 HY000 The requested value '2000000001' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2000000001
| ; SELECT FROM
< error
| 3567 HY000 The requested value '2000000001' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2147483646
< error
| 3567 HY000 The requested value '2147483646' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2147483646
| ; SELECT FROM
< error
| 3567 HY000 The requested value '2147483646' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| CREATE TABLE t (a CHAR(0x10))
< error
| 1064 42000 …near '0x10))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(1.))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999.5))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999.5))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775806))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775806))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999.5))
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999.5))
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column 'cast as char' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 18446744073709551615K
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 18446744073709551615K
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 99999999999999999999K
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 99999999999999999999K
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 4M
< ok 0 0
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999K
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999K
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 4K
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 4K
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 4K
< error
| 1064 42000 …near '4K' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END
| ; SELECT FROM
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0.5 BEGIN END; END
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0.5 BEGIN END; END
| ; SELECT FROM
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 LIMIT 0x10
< error
| 1064 42000 …near '0x10' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483647K
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648K
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648K
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648G
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648G
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = K
< ok 0 0
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648M
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648M
| ; SELECT FROM
< error
| 1532 HY000 The size number was correct but we don't allow the digit part to be more than 2 billion
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE 4M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE 16384
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE 16384
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE 4M MAX_SIZE 1G EXTENT_SIZE 1M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE 4M MAX_SIZE 1G EXTENT_SIZE 1M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 4M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE 4M UNDO_BUFFER_SIZE 4M
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE 4M UNDO_BUFFER_SIZE 4M
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END
| ; SELECT FROM
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0.9; END
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0.9; END
| ; SELECT FROM
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.9; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.9; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(1.5) CHARACTER SET utf8mb4)
< ok 0 0
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1.5)) RETURNS VARCHAR(2.5) RETURN a
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1.5)) RETURNS VARCHAR(2.5) RETURN a
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a YEAR(.5))
< error
| 1818 HY000 Invalid display width. Use YEAR instead.
> exec ^
: string
| CREATE TABLE t (a YEAR(.5))
| ; SELECT FROM
< error
| 1818 HY000 Invalid display width. Use YEAR instead.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a VARCHAR(.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(.5))
< ok 0 0
> exec ^
: string
| SELECT CAST(a AS DECIMAL(.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS CHAR(1.))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(1.))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR .5 BEGIN END; END
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR .5 BEGIN END; END
| ; SELECT FROM
< error
| 1525 HY000 Incorrect CONDITION value: '0'
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647
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
| EXPLAIN FOR CONNECTION .5
< error
| 1064 42000 Only integers allowed as number here near '.5' at line 1
> exec ^
: string
| RESET MASTER TO .5
< error
| 1064 42000 Only integers allowed as number here near '.5' at line 1
> exec ^
: string
| RESET MASTER TO 0
< error
| 3567 HY000 The requested value '0' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 0
| ; SELECT FROM
< error
| 3567 HY000 The requested value '0' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 254 BUCKETS
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 'x'
< error
| 1064 42000 …near ''x'' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967295
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967296
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '4294967296' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65535
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65536
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '65536' at line 1
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648
| ; SELECT FROM
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1.5)))
< error
| 1064 42000 …near '1.5)))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARCHAR(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BINARY(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BINARY(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BINARY(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BINARY(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483647))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483648))
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARBINARY(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a VARBINARY(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a VARBINARY(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BLOB(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BLOB(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BLOB(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BLOB(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TEXT(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a TEXT(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TEXT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TEXT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIT(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIT(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(2147483647))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(2147483648))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a INT(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a INT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIGINT(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483647))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483648))
< error
| 1439 42000 Display width out of range for column 'a' (max = 255)
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a BIGINT(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a BIGINT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a BIGINT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NCHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NCHAR(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NCHAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a NUMERIC(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a NUMERIC(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FIXED(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FIXED(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FIXED(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FIXED(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'a' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'a table function argument'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483647) PATH '$')) jt
< error
| 1074 42000 Column length too big for column 'b' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483647) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483648) PATH '$')) jt
< error
| 1074 42000 Column length too big for column 'b' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483648) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775807) PATH '$')) jt
< error
| 1439 42000 Display width out of range for column 'b' (max = 4294967295)
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775807) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775808) PATH '$')) jt
< error
| 1439 42000 Display width out of range for column 'b' (max = 4294967295)
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775808) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551615) PATH '$')) jt
< error
| 1439 42000 Display width out of range for column 'b' (max = 4294967295)
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551615) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551616) PATH '$')) jt
< error
| 1439 42000 Display width out of range for column 'b' (max = 4294967295)
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551616) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(99999999999999999999) PATH '$')) jt
< error
| 1439 42000 Display width out of range for column 'b' (max = 4294967295)
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(99999999999999999999) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(000000000000000000000001) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'a table function argument'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(000000000000000000000001) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1.5) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'a table function argument'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1.5) PATH '$')) jt
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1e1) PATH '$')) jt
< error
| 1064 42000 …near '1e1) PATH '$')) jt' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1,0))
< error
| 1064 42000 …near '1e1,0))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(1e1,0))
< error
| 1064 42000 …near '1e1,0))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483648,0))
< error
| 1064 42000 …near '2147483648,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775807,0))
< error
| 1064 42000 …near '9223372036854775807,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775808,0))
< error
| 1064 42000 …near '9223372036854775808,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551615,0))
< error
| 1064 42000 …near '18446744073709551615,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551616,0))
< error
| 1064 42000 …near '18446744073709551616,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(99999999999999999999,0))
< error
| 1064 42000 …near '99999999999999999999,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(000000000000000000000001,0))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1.5,0))
< error
| 1064 42000 …near '1.5,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1e1,0))
< error
| 1064 42000 …near '1e1,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,2147483647))
< error
| 1425 42000 Too big scale 2147483647 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(10,1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a REAL(10,2147483647))
< error
| 1425 42000 Too big scale 2147483647 specified for column 'a'. Maximum is 30.
> exec ^
: string
| CREATE TABLE t (a REAL(10,2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a REAL(10,2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a REAL(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a REAL(10,1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 6.
> exec ^
: string
| CREATE TABLE t (a TIME(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TIME(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 6.
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a DATETIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 6.
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| SELECT NOW(1)
< ok 0 0
> exec ^
: string
| SELECT NOW(2147483647)
< error
| 1426 42000 Too-big precision 255 specified for 'now'. Maximum is 6.
> exec ^
: string
| SELECT NOW(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT NOW(2147483648)
< error
| 1064 42000 …near '2147483648)' at line 1
> exec ^
: string
| SELECT NOW(9223372036854775807)
< error
| 1064 42000 …near '9223372036854775807)' at line 1
> exec ^
: string
| SELECT NOW(9223372036854775808)
< error
| 1064 42000 …near '9223372036854775808)' at line 1
> exec ^
: string
| SELECT NOW(18446744073709551615)
< error
| 1064 42000 …near '18446744073709551615)' at line 1
> exec ^
: string
| SELECT NOW(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| SELECT NOW(99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| SELECT NOW(000000000000000000000001)
< ok 0 0
> exec ^
: string
| SELECT NOW(1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| SELECT NOW(1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| SELECT CURTIME(1)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(2147483647)
< error
| 1426 42000 Too-big precision 255 specified for 'curtime'. Maximum is 6.
> exec ^
: string
| SELECT CURTIME(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CURTIME(2147483648)
< error
| 1064 42000 …near '2147483648)' at line 1
> exec ^
: string
| SELECT CURTIME(9223372036854775807)
< error
| 1064 42000 …near '9223372036854775807)' at line 1
> exec ^
: string
| SELECT CURTIME(9223372036854775808)
< error
| 1064 42000 …near '9223372036854775808)' at line 1
> exec ^
: string
| SELECT CURTIME(18446744073709551615)
< error
| 1064 42000 …near '18446744073709551615)' at line 1
> exec ^
: string
| SELECT CURTIME(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| SELECT CURTIME(99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| SELECT CURTIME(000000000000000000000001)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| SELECT CURTIME(1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| SELECT SYSDATE(1)
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(2147483647)
< error
| 1426 42000 Too-big precision 255 specified for 'sysdate'. Maximum is 6.
> exec ^
: string
| SELECT SYSDATE(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SYSDATE(2147483648)
< error
| 1064 42000 …near '2147483648)' at line 1
> exec ^
: string
| SELECT SYSDATE(9223372036854775807)
< error
| 1064 42000 …near '9223372036854775807)' at line 1
> exec ^
: string
| SELECT SYSDATE(9223372036854775808)
< error
| 1064 42000 …near '9223372036854775808)' at line 1
> exec ^
: string
| SELECT SYSDATE(18446744073709551615)
< error
| 1064 42000 …near '18446744073709551615)' at line 1
> exec ^
: string
| SELECT SYSDATE(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| SELECT SYSDATE(99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| SELECT SYSDATE(000000000000000000000001)
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| SELECT SYSDATE(1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(1)
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(2147483647)
< error
| 1426 42000 Too-big precision 255 specified for 'utc_timestamp'. Maximum is 6.
> exec ^
: string
| SELECT UTC_TIMESTAMP(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT UTC_TIMESTAMP(2147483648)
< error
| 1064 42000 …near '2147483648)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(9223372036854775807)
< error
| 1064 42000 …near '9223372036854775807)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(9223372036854775808)
< error
| 1064 42000 …near '9223372036854775808)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(18446744073709551615)
< error
| 1064 42000 …near '18446744073709551615)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(000000000000000000000001)
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| SELECT UTC_TIMESTAMP(1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(1)
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(2147483647)
< error
| 1426 42000 Too-big precision 255 specified for 'now'. Maximum is 6.
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(2147483648)
< error
| 1064 42000 …near '2147483648)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(9223372036854775807)
< error
| 1064 42000 …near '9223372036854775807)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(9223372036854775808)
< error
| 1064 42000 …near '9223372036854775808)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(18446744073709551615)
< error
| 1064 42000 …near '18446744073709551615)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(000000000000000000000001)
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(1))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(2147483647))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(000000000000000000000001))
< error
| 1067 42000 Invalid default value for 'a'
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1))
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483647))
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(000000000000000000000001))
< error
| 1294 HY000 Invalid ON UPDATE clause for 'a' column
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(1))
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a(2147483647))
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(000000000000000000000001))
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
< error
| 1064 42000 …near '2147483648)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775807)))
< error
| 1064 42000 …near '9223372036854775807)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775808)))
< error
| 1064 42000 …near '9223372036854775808)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(18446744073709551615)))
< error
| 1064 42000 …near '18446744073709551615)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(18446744073709551616)))
< error
| 1064 42000 …near '18446744073709551616)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(99999999999999999999)))
< error
| 1064 42000 …near '99999999999999999999)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(000000000000000000000001)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1.5)))
< error
| 1064 42000 …near '1.5)))' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(1e1)))
< error
| 1064 42000 …near '1e1)))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 LINES
< error
| 1064 42000 …near '2147483648 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 LINES
< error
| 1064 42000 …near '9223372036854775807 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 LINES
< error
| 1064 42000 …near '9223372036854775808 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 LINES
< error
| 1064 42000 …near '18446744073709551615 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551616 LINES
< error
| 1064 42000 …near '18446744073709551616 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 99999999999999999999 LINES
< error
| 1064 42000 …near '99999999999999999999 LINES' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1.5 LINES
< error
| 1064 42000 …near '1.5 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1e1 LINES
< error
| 1064 42000 …near '1e1 LINES' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1 ROWS
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1 ROWS
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 ROWS
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 ROWS
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 ROWS
< error
| 1064 42000 …near '2147483648 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 ROWS
< error
| 1064 42000 …near '9223372036854775807 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 ROWS
< error
| 1064 42000 …near '9223372036854775808 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 ROWS
< error
| 1064 42000 …near '18446744073709551615 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551616 ROWS
< error
| 1064 42000 …near '18446744073709551616 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 99999999999999999999 ROWS
< error
| 1064 42000 …near '99999999999999999999 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 000000000000000000000001 ROWS
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 000000000000000000000001 ROWS
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1.5 ROWS
< error
| 1064 42000 …near '1.5 ROWS' at line 1
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1e1 ROWS
< error
| 1064 42000 …near '1e1 ROWS' at line 1
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
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 000000000000000000000001
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SHOW PROFILE FOR QUERY 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 2147483647
< error
| 3652 HY000 Invalid cpu id 2147483647
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-2147483647
< error
| 3652 HY000 Invalid cpu id 2147483647
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER VCPU = 0-1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE RESOURCE GROUP g TYPE = USER THREAD_PRIORITY = -1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 2147483647
< error
| 3652 HY000 Invalid cpu id 2147483647
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g VCPU = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| ALTER RESOURCE GROUP g THREAD_PRIORITY = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551616))
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(99999999999999999999))
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1.5))
< error
| 1064 42000 Only integers allowed as number here near '1.5))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1e1))
< error
| 1064 42000 Only integers allowed as number here near '1e1))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483647))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483648))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775807))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775808))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551615))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551616))
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(99999999999999999999))
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1.5))
< error
| 1064 42000 Only integers allowed as number here near '1.5))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1e1))
< error
| 1064 42000 Only integers allowed as number here near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AUTOEXTEND_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
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
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
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
| 1064 42000 Only integers allowed as number here near '18446744073709551616 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 000000000000000000000001 (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1.5 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 Only integers allowed as number here near '1.5 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1e1 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 Only integers allowed as number here near '1e1 (PARTITION p0 VALUES LESS THAN (10))' at line 1
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
| 1064 42000 Only integers allowed as number here near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 99999999999999999999)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1e1)
< error
| 1064 42000 Only integers allowed as number here near '1e1)' at line 1
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
| 1064 42000 Only integers allowed as number here near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 99999999999999999999)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MIN_ROWS = 1e1)
< error
| 1064 42000 Only integers allowed as number here near '1e1)' at line 1
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
| 1064 42000 Only integers allowed as number here near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 99999999999999999999)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 NODEGROUP = 1e1)
< error
| 1064 42000 Only integers allowed as number here near '1e1)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483647
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 2147483648
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER TABLE t ADD PARTITION PARTITIONS 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483647
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 2147483648
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER TABLE t COALESCE PARTITION 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
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
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SELECT 1 LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SELECT 1 LIMIT 000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
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
| 1064 42000 …near '18446744073709551616, 1' at line 1
> exec ^
: string
| SELECT 1 LIMIT 99999999999999999999, 1
< error
| 1064 42000 …near '99999999999999999999, 1' at line 1
> exec ^
: string
| SELECT 1 LIMIT 000000000000000000000001, 1
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1.5, 1
< error
| 1064 42000 …near '1.5, 1' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1e1, 1
< error
| 1064 42000 …near '1e1, 1' at line 1
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
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 000000000000000000000001
< ok 0 0
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SELECT 1 LIMIT 1 OFFSET 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| DELETE FROM t LIMIT 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 2147483647
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 2147483648
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775807
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775808
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551615
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| DELETE FROM t LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| DELETE FROM t LIMIT 000000000000000000000001
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| DELETE FROM t LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DELETE FROM t LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| DELETE FROM t LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483647
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483648
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775807
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775808
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551615
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 000000000000000000000001
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| HANDLER t READ FIRST LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| HANDLER t READ FIRST LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SHOW WARNINGS LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW WARNINGS LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SHOW WARNINGS LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 1
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 2147483647
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 2147483648
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 9223372036854775807
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 9223372036854775808
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 18446744073709551615
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 000000000000000000000001
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 1
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 2147483647
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 2147483648
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 9223372036854775807
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 9223372036854775808
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 18446744073709551615
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 000000000000000000000001
< ok 0 0
> exec ^
: string
| SHOW PROFILE LIMIT 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| SHOW PROFILE LIMIT 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| EXPLAIN FOR CONNECTION 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| EXPLAIN FOR CONNECTION 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| EXPLAIN FOR CONNECTION 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
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
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE USER u PASSWORD HISTORY 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
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
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 18446744073709551615 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 18446744073709551616 DAY
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 99999999999999999999 DAY
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 000000000000000000000001 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 1.5 DAY
< error
| 1064 42000 Only integers allowed as number here near '1.5 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD REUSE INTERVAL 1e1 DAY
< error
| 1064 42000 Only integers allowed as number here near '1e1 DAY' at line 1
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
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER USER u PASSWORD HISTORY 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| SET RESOURCE GROUP g FOR 1
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 2147483647
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 2147483648
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 9223372036854775807
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 9223372036854775808
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 18446744073709551615
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| SET RESOURCE GROUP g FOR 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| SET RESOURCE GROUP g FOR 000000000000000000000001
< error
| 1227 42000 Access denied; you need (at least one of) the RESOURCE_GROUP_ADMIN OR RESOURCE_GROUP_USER privilege(s) for this operation
> exec ^
: string
| SET RESOURCE GROUP g FOR 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SET RESOURCE GROUP g FOR 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| SET RESOURCE GROUP g FOR 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 2147483647 NAME 'n' DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 2147483647 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 2147483648 NAME 'n' DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 2147483648 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 9223372036854775807 NAME 'n' DEFINITION 'd'
< error
| 1690 22003 SRID value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 9223372036854775807 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 9223372036854775808 NAME 'n' DEFINITION 'd'
< error
| 1690 22003 SRID value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 9223372036854775808 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 18446744073709551615 NAME 'n' DEFINITION 'd'
< error
| 1690 22003 SRID value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 18446744073709551615 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 18446744073709551616 NAME 'n' DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616 NAME 'n' DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 99999999999999999999 NAME 'n' DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 NAME 'n' DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 000000000000000000000001 NAME 'n' DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 000000000000000000000001 NAME 'n' DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1.5 NAME 'n' DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '1.5 NAME 'n' DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1e1 NAME 'n' DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '1e1 NAME 'n' DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 1 DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 1 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 2147483647 DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 2147483647 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 2147483648 DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 2147483648 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 9223372036854775807 DEFINITION 'd'
< error
| 1690 22003 IDENTIFIED BY value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 9223372036854775807 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 9223372036854775808 DEFINITION 'd'
< error
| 1690 22003 IDENTIFIED BY value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 9223372036854775808 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 18446744073709551615 DEFINITION 'd'
< error
| 1690 22003 IDENTIFIED BY value is out of range in 'CREATE SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 18446744073709551615 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 18446744073709551616 DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616 DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 99999999999999999999 DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 000000000000000000000001 DEFINITION 'd'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 000000000000000000000001 DEFINITION 'd'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 1.5 DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '1.5 DEFINITION 'd'' at line 1
> exec ^
: string
| CREATE SPATIAL REFERENCE SYSTEM 1 NAME 'n' ORGANIZATION 'o' IDENTIFIED BY 1e1 DEFINITION 'd'
< error
| 1064 42000 Only integers allowed as number here near '1e1 DEFINITION 'd'' at line 1
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 9223372036854775807
< error
| 1690 22003 SRID value is out of range in 'DROP SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 9223372036854775808
< error
| 1690 22003 SRID value is out of range in 'DROP SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 18446744073709551615
< error
| 1690 22003 SRID value is out of range in 'DROP SPATIAL REFERENCE SYSTEM'
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| DROP SPATIAL REFERENCE SYSTEM 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER TABLESPACE ts AUTOEXTEND_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483647 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483647 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483648 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483648 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775807 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775807 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775808 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775808 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551615 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551615 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551616 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551616 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 99999999999999999999 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 99999999999999999999 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 000000000000000000000001 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 000000000000000000000001 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1.5 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1.5 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1e1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1e1 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483647 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483647 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483648 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483648 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775807 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775807 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775808 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775808 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551615 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551615 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551616 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551616 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -99999999999999999999 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -99999999999999999999 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -000000000000000000000001 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -000000000000000000000001 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1.5 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1.5 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1e1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1e1 ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483647
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483648
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551616
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 99999999999999999999
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1.5
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1e1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483647
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483648
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551616
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -99999999999999999999
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1.5
< ok 0 0
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1e1
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1.5
< ok 0 0
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 99999999999999999999)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1.5)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1e1)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) TABLE_CHECKSUM = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1e1
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483647
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483648
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551616
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 99999999999999999999
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1.5
< ok 0 0
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1e1
< ok 0 0
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
| SHOW BINLOG EVENTS FROM 1
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 2147483647
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 2147483648
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 9223372036854775807
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 9223372036854775808
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 18446744073709551615
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 18446744073709551616
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 99999999999999999999
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 000000000000000000000001
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 1.5
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 1.5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW BINLOG EVENTS FROM 1e1
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW BINLOG EVENTS FROM 1e1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 2147483647
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 2147483648
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 9223372036854775807
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 9223372036854775808
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 18446744073709551615
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 18446744073709551616
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 99999999999999999999
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 000000000000000000000001
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1.5
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1.5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1e1
< error
| 1227 42000 Access denied; you need (at least one of) the REPLICATION SLAVE privilege(s) for this operation
> exec ^
: string
| SHOW RELAYLOG EVENTS FROM 1e1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
< ok 0 0
> exec ^
: string
| KILL 99999999999999999999
< ok 0 0
> exec ^
: string
| KILL 000000000000000000000001
< ok 0 0
> exec ^
: string
| KILL 1.5
< ok 0 0
> exec ^
: string
| KILL 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| KILL QUERY 99999999999999999999
< ok 0 0
> exec ^
: string
| KILL QUERY 000000000000000000000001
< ok 0 0
> exec ^
: string
| KILL QUERY 1.5
< ok 0 0
> exec ^
: string
| KILL QUERY 1e1
< ok 0 0
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483647 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483647 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483648 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483648 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775807 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775807 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775808 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775808 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551615 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551615 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551616 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551616 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 99999999999999999999 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 99999999999999999999 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 000000000000000000000001 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 000000000000000000000001 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1.5 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1e1 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1e1 BEGIN END; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483647; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483647; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483648; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483648; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775807; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775807; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775808; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775808; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551615; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551615; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551616; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551616; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 99999999999999999999; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 99999999999999999999; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 000000000000000000000001; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 000000000000000000000001; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.5; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1.5; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1e1; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1e1; END
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483647
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483648
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775807
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775808
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551615
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551616
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 99999999999999999999
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1.5
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1.5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1e1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1e1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 2147483647 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 2147483647 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 2147483648 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 2147483648 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 9223372036854775807 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 9223372036854775807 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 9223372036854775808 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 9223372036854775808 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 18446744073709551615 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 18446744073709551615 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 18446744073709551616 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 18446744073709551616 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 99999999999999999999 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 99999999999999999999 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 000000000000000000000001 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 000000000000000000000001 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1.5 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1.5 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1e1 @a = MESSAGE_TEXT
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| GET DIAGNOSTICS CONDITION 1e1 @a = MESSAGE_TEXT
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (2147483647)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (2147483648)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (9223372036854775807)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (9223372036854775808)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (18446744073709551615)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (18446744073709551616)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (99999999999999999999)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (000000000000000000000001)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1.5)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1e1)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_PORT = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_LOG_POS = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO RELAY_LOG_POS = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_CONNECT_RETRY = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_RETRY_COUNT = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_SSL = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_AUTO_POSITION = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_ZSTD_COMPRESSION_LEVEL = 1e1
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 1
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 1.5
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_PORT = 1e1
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 1
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| CHANGE REPLICATION SOURCE TO SOURCE_LOG_POS = 1e1
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 1
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| START SLAVE UNTIL MASTER_LOG_FILE = 'f', MASTER_LOG_POS = 1e1
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 1
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL SOURCE_LOG_FILE = 'f', SOURCE_LOG_POS = 1e1
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 1
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 2147483647
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 2147483648
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 9223372036854775807
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 9223372036854775808
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 18446744073709551615
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 18446744073709551616
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 99999999999999999999
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 1.5
< ok 0 0
> exec ^
: string
| START REPLICA UNTIL RELAY_LOG_FILE = 'f', RELAY_LOG_POS = 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 1.5
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 1.5
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_UPDATES_PER_HOUR 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1.5
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 18446744073709551616
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1.5
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1e1
< ok 0 0
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
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 18446744073709551616
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 99999999999999999999
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 1.5
< ok 0 0
> exec ^
: string
| ALTER USER u WITH MAX_QUERIES_PER_HOUR 1e1
< ok 0 0
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 2147483647)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 2147483648)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 2147483648)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 9223372036854775807)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 9223372036854775808)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 18446744073709551615)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 18446744073709551616)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 99999999999999999999)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 000000000000000000000001)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1.5)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1.5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1e1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE SERVER s FOREIGN DATA WRAPPER mysql OPTIONS (PORT 1e1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 2147483647)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 2147483648)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 2147483648)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 9223372036854775807)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 9223372036854775808)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 18446744073709551615)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 18446744073709551616)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 99999999999999999999)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 000000000000000000000001)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1.5)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1.5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1e1)
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER SERVER s OPTIONS (PORT 1e1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483647 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483647 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483648 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483648 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775807 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775807 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775808 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775808 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551615 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551615 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551616 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551616 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 99999999999999999999 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 99999999999999999999 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 000000000000000000000001 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 000000000000000000000001 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1.5 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1.5 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1e1 PRECEDING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1e1 PRECEDING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483647 PRECEDING AND 2147483647 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483647 PRECEDING AND 2147483647 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483648 PRECEDING AND 2147483648 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483648 PRECEDING AND 2147483648 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775807 PRECEDING AND 9223372036854775807 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775807 PRECEDING AND 9223372036854775807 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775808 PRECEDING AND 9223372036854775808 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775808 PRECEDING AND 9223372036854775808 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551615 PRECEDING AND 18446744073709551615 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551615 PRECEDING AND 18446744073709551615 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551616 PRECEDING AND 18446744073709551616 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551616 PRECEDING AND 18446744073709551616 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 99999999999999999999 PRECEDING AND 99999999999999999999 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 99999999999999999999 PRECEDING AND 99999999999999999999 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 000000000000000000000001 PRECEDING AND 000000000000000000000001 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 000000000000000000000001 PRECEDING AND 000000000000000000000001 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1.5 PRECEDING AND 1.5 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1.5 PRECEDING AND 1.5 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1e1 PRECEDING AND 1e1 FOLLOWING) FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1e1 PRECEDING AND 1e1 FOLLOWING) FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 2147483647 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 2147483647 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 2147483648 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 2147483648 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 9223372036854775807 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 9223372036854775807 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 9223372036854775808 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 9223372036854775808 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 18446744073709551615 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 18446744073709551615 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 18446744073709551616 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 18446744073709551616 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 99999999999999999999 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 99999999999999999999 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 000000000000000000000001 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 000000000000000000000001 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1.5 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1.5 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1e1 DAY DO SELECT 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE EVENT e ON SCHEDULE EVERY 1e1 DAY DO SELECT 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483647
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483647
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483648
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483648
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775807
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775808
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775808
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551615
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551615
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551616
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551616
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 99999999999999999999
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 99999999999999999999
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 000000000000000000000001
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 1.5
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1.5
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t ORDER BY 1e1
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t ORDER BY 1e1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483647)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483648)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483648)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775807)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775808)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551615)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551616)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (99999999999999999999)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (000000000000000000000001)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (1.5)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1.5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT a FROM t WHERE a IN (1e1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT a FROM t WHERE a IN (1e1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (2147483647)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (2147483648)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (2147483648)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775807)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775808)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551615)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551616)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (99999999999999999999)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (000000000000000000000001)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (1.5)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1.5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (1e1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1e1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(1)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(2147483647)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(2147483647)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(2147483648)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(2147483648)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(9223372036854775807)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(9223372036854775807)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(9223372036854775808)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(9223372036854775808)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(18446744073709551615)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(18446744073709551615)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(18446744073709551616)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(18446744073709551616)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(99999999999999999999)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(99999999999999999999)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(000000000000000000000001)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(000000000000000000000001)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(1.5)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(1.5)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL p(1e1)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.p does not exist
> exec ^
: string
| CALL p(1e1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 1)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 2147483647)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 2147483648)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 9223372036854775807)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 9223372036854775808)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 18446744073709551615)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 18446744073709551616)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 99999999999999999999)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 1.5)
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO IGNORE_SERVER_IDS = (1, 1e1)
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647,0))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647,0))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …near ',0))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1e1,0))
< error
| 1064 42000 …near '1e1,0))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483647))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483647))
| ; SELECT FROM
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column '').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS TIME(2147483647))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS TIME(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS TIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 6.
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483648))
< error
| 1064 42000 …near '2147483648))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775807))
< error
| 1064 42000 …near '9223372036854775807))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775808))
< error
| 1064 42000 …near '9223372036854775808))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551615))
< error
| 1064 42000 …near '18446744073709551615))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DATETIME(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DATETIME(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483647 (a)
< error
| 1064 42000 …near '2147483647 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 2147483648 (a)
< error
| 1064 42000 …near '2147483648 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775807 (a)
< error
| 1064 42000 …near '9223372036854775807 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 9223372036854775808 (a)
< error
| 1064 42000 …near '9223372036854775808 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 18446744073709551615 (a)
< error
| 1064 42000 …near '18446744073709551615 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 18446744073709551616 (a)
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 99999999999999999999 (a)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 000000000000000000000001 (a)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1.5 (a)
< error
| 1064 42000 Only integers allowed as number here near '1.5 (a)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY KEY ALGORITHM = 1e1 (a)
< error
| 1064 42000 Only integers allowed as number here near '1e1 (a)' at line 1
> exec ^
: string
| RESET MASTER TO 1
< ok 0 0
> exec ^
: string
| RESET MASTER TO 2147483647
< error
| 3567 HY000 The requested value '2147483647' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2147483647
| ; SELECT FROM
< error
| 3567 HY000 The requested value '2147483647' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2147483648
< error
| 3567 HY000 The requested value '2147483648' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 2147483648
| ; SELECT FROM
< error
| 3567 HY000 The requested value '2147483648' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 9223372036854775807
< error
| 3567 HY000 The requested value '9223372036854775807' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 9223372036854775807
| ; SELECT FROM
< error
| 3567 HY000 The requested value '9223372036854775807' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 9223372036854775808
< error
| 3567 HY000 The requested value '9223372036854775808' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 9223372036854775808
| ; SELECT FROM
< error
| 3567 HY000 The requested value '9223372036854775808' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 18446744073709551615
< error
| 3567 HY000 The requested value '18446744073709551615' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 18446744073709551615
| ; SELECT FROM
< error
| 3567 HY000 The requested value '18446744073709551615' for the next binary log index is out of range. Please use a value between '1' and '2000000000'.
> exec ^
: string
| RESET MASTER TO 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| RESET MASTER TO 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| RESET MASTER TO 000000000000000000000001
< ok 0 0
> exec ^
: string
| RESET MASTER TO 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| RESET MASTER TO 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483647 DAY
< error
| 1525 HY000 Incorrect DAY value: '2147483647'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483647 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '2147483647'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483648 DAY
< error
| 1525 HY000 Incorrect DAY value: '2147483648'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 2147483648 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '2147483648'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775807 DAY
< error
| 1525 HY000 Incorrect DAY value: '9223372036854775807'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775807 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '9223372036854775807'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775808 DAY
< error
| 1525 HY000 Incorrect DAY value: '9223372036854775808'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 9223372036854775808 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '9223372036854775808'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 18446744073709551615 DAY
< error
| 1525 HY000 Incorrect DAY value: '18446744073709551615'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 18446744073709551615 DAY
| ; SELECT FROM
< error
| 1525 HY000 Incorrect DAY value: '18446744073709551615'
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 18446744073709551616 DAY
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 99999999999999999999 DAY
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 000000000000000000000001 DAY
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1.5 DAY
< error
| 1064 42000 Only integers allowed as number here near '1.5 DAY' at line 1
> exec ^
: string
| CREATE USER u PASSWORD EXPIRE INTERVAL 1e1 DAY
< error
| 1064 42000 Only integers allowed as number here near '1e1 DAY' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483647
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '2147483647'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483647
| ; SELECT FROM
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '2147483647'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483648
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '2147483648'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 2147483648
| ; SELECT FROM
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '2147483648'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775807
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '9223372036854775807'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775807
| ; SELECT FROM
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '9223372036854775807'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775808
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '9223372036854775808'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 9223372036854775808
| ; SELECT FROM
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '9223372036854775808'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 18446744073709551615
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '18446744073709551615'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 18446744073709551615
| ; SELECT FROM
< error
| 1525 HY000 Incorrect FAILED_LOGIN_ATTEMPTS value: '18446744073709551615'
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE USER u FAILED_LOGIN_ATTEMPTS 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483647
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '2147483647'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483647
| ; SELECT FROM
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '2147483647'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483648
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '2147483648'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 2147483648
| ; SELECT FROM
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '2147483648'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775807
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '9223372036854775807'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775807
| ; SELECT FROM
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '9223372036854775807'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775808
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '9223372036854775808'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 9223372036854775808
| ; SELECT FROM
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '9223372036854775808'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 18446744073709551615
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '18446744073709551615'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 18446744073709551615
| ; SELECT FROM
< error
| 1525 HY000 Incorrect PASSWORD_LOCK_TIME value: '18446744073709551615'
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 18446744073709551616
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 99999999999999999999
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1.5
< error
| 1064 42000 Only integers allowed as number here near '1.5' at line 1
> exec ^
: string
| CREATE USER u PASSWORD_LOCK_TIME 1e1
< error
| 1064 42000 Only integers allowed as number here near '1e1' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1 BUCKETS
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS
| ; SELECT FROM
< error
| 1690 22003 Number of buckets value is out of range in 'ANALYZE TABLE'
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483648 BUCKETS
< error
| 1064 42000 …near '2147483648 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775807 BUCKETS
< error
| 1064 42000 …near '9223372036854775807 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775808 BUCKETS
< error
| 1064 42000 …near '9223372036854775808 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551615 BUCKETS
< error
| 1064 42000 …near '18446744073709551615 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551616 BUCKETS
< error
| 1064 42000 …near '18446744073709551616 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 99999999999999999999 BUCKETS
< error
| 1064 42000 …near '99999999999999999999 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 000000000000000000000001 BUCKETS
< ok 0 0
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1.5 BUCKETS
< error
| 1064 42000 …near '1.5 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1e1 BUCKETS
< error
| 1064 42000 …near '1e1 BUCKETS' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -2147483648 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807))
< error
| 1426 42000 Too-big precision -1 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808))
< error
| 1426 42000 Too-big precision 0 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 0 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615))
< error
| 1426 42000 Too-big precision -1 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'a'. Maximum is 65.
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(1))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483647))
< error
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483647))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 2147483647 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483648))
< error
| 1426 42000 Too-big precision -2147483648 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483648))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -2147483648 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775807))
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775807))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775808))
< error
| 1426 42000 Too-big precision 0 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775808))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 0 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551615))
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551615))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551616))
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551616))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(99999999999999999999))
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(99999999999999999999))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision -1 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| SELECT CAST(a AS FLOAT(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(000000000000000000000001))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT CAST(a AS FLOAT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1)) RETURNS INT RETURN 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483647)) RETURNS INT RETURN 1
< error
| 1074 42000 Column length too big for column '' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483647)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column '' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483648)) RETURNS INT RETURN 1
< error
| 1074 42000 Column length too big for column '' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483648)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1074 42000 Column length too big for column '' (max = 255); use BLOB or TEXT instead
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775807)) RETURNS INT RETURN 1
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775807)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775808)) RETURNS INT RETURN 1
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775808)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551615)) RETURNS INT RETURN 1
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551615)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551616)) RETURNS INT RETURN 1
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551616)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(99999999999999999999)) RETURNS INT RETURN 1
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(99999999999999999999)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1439 42000 Display width out of range for column '' (max = 4294967295)
> exec ^
: string
| CREATE FUNCTION f(a CHAR(000000000000000000000001)) RETURNS INT RETURN 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION f(a CHAR(000000000000000000000001)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1.5)) RETURNS INT RETURN 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1.5)) RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1e1)) RETURNS INT RETURN 1
< error
| 1064 42000 …near '1e1)) RETURNS INT RETURN 1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483647
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '2147483647' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483648
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '2147483648' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775807
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775808
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551615
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 The valid range for key_block_size is [0,65535]. Error near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483647
< error
| 1064 42000 …near '2147483647' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483647
< error
| 1064 42000 …near '2147483647' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483647
< error
| 1064 42000 …near '2147483647' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483647
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '2147483647' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483648
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '2147483648' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775807
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775808
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551615
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551616
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 99999999999999999999
< error
| 1064 42000 The valid range for stats_sample_pages is [1, 65535]. Error near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 2147483647
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 2147483647
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 2147483648
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 2147483648
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 9223372036854775807
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 9223372036854775807
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 9223372036854775808
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 9223372036854775808
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 18446744073709551615
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 18446744073709551615
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 18446744073709551616
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 18446744073709551616
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 99999999999999999999
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 99999999999999999999
| ; SELECT FROM
< error
| 1624 HY000 The requested value for the heartbeat period is either negative or exceeds the maximum allowed (4294967 seconds).
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_HEARTBEAT_PERIOD = 1e1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 1
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483647
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 2147483648
| ; SELECT FROM
< error
| 1729 HY000 The requested value 2147483648 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 9223372036854775807
< error
| 1729 HY000 The requested value 9223372036854775807 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 9223372036854775807
| ; SELECT FROM
< error
| 1729 HY000 The requested value 9223372036854775807 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 9223372036854775808
< error
| 1729 HY000 The requested value 9223372036854775808 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 9223372036854775808
| ; SELECT FROM
< error
| 1729 HY000 The requested value 9223372036854775808 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 18446744073709551615
< error
| 1729 HY000 The requested value 18446744073709551615 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 18446744073709551615
| ; SELECT FROM
< error
| 1729 HY000 The requested value 18446744073709551615 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 18446744073709551616
< error
| 1729 HY000 The requested value 18446744073709551616 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 18446744073709551616
| ; SELECT FROM
< error
| 1729 HY000 The requested value 18446744073709551616 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 99999999999999999999
< error
| 1729 HY000 The requested value 99999999999999999999 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 99999999999999999999
| ; SELECT FROM
< error
| 1729 HY000 The requested value 99999999999999999999 for the source delay exceeds the maximum 2147483647
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 1.5
< ok 0 0
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 1e1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 2147483647)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 2147483648)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775807)
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775807)
| ; SELECT FROM
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775808)
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775808)
| ; SELECT FROM
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 18446744073709551615)
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 18446744073709551615)
| ; SELECT FROM
< error
| 1690 22003 SRID value is out of range in 'SRID'
> exec ^
: string
| CREATE TABLE t (a POINT SRID 18446744073709551616)
< error
| 1064 42000 Only integers allowed as number here near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 99999999999999999999)
< error
| 1064 42000 Only integers allowed as number here near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 000000000000000000000001)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1.5)
< error
| 1064 42000 Only integers allowed as number here near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1e1)
< error
| 1064 42000 Only integers allowed as number here near '1e1)' at line 1
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
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '9223372036854775807' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 9223372036854775808
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551615
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551616
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 99999999999999999999
< error
| 1064 42000 The valid range for avg_row_length is [0,4294967295]. Error near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1.5
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1e1
< ok 0 0
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| ALTER DATABASE d READ ONLY = 2147483647
< error
| 1064 42000 …near '2147483647' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 2147483648
< error
| 1064 42000 …near '2147483648' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 9223372036854775807
< error
| 1064 42000 …near '9223372036854775807' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 9223372036854775808
< error
| 1064 42000 …near '9223372036854775808' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 18446744073709551615
< error
| 1064 42000 …near '18446744073709551615' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| ALTER DATABASE d READ ONLY = 000000000000000000000001
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| ALTER DATABASE d READ ONLY = 000000000000000000000001
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
