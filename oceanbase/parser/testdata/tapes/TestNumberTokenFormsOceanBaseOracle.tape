-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestNumberTokenFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| CREATE TABLE t (a POINT SRID 4294967295)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 4294967296)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 0x10)
< error
| 1064 42000 …near '0x10)' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 0 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1024 BUCKETS
< error
| 1046 3D000 No database selected
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1025 BUCKETS' at line 1
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x7fffffff
< error
| 1064 42000 …near 'CHANGE MASTER TO MASTER_DELAY = 0x7fffffff' at line 1
> exec ^
: string
| CHANGE MASTER TO MASTER_DELAY = 0x80000000
< error
| 1064 42000 …near 'CHANGE MASTER TO MASTER_DELAY = 0x80000000' at line 1
> exec ^
: string
| RESET MASTER TO 0x10
< error
| 1064 42000 …near 'RESET MASTER TO 0x10' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 0x10
< error
| 1064 42000 …near '0x10' at line 1
> exec ^
: string
| CREATE TABLE t (a CHAR(0x10))
< error
| 1064 42000 …near '0x10))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999.5))
< error
| 1064 42000 …near '99999999999999999999.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967295))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(4294967296))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775806))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551615.5))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999.5))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 1.5)
< error
| 1064 42000 …near 'MAX_ROWS = 1.5)' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 0x10
< error
| 1064 42000 …near '0x10' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_QUERIES_PER_HOUR 0x10
< error
| 1064 42000 …near 'MAX_QUERIES_PER_HOUR 0x10' at line 1
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 0 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 4294967296 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| SELECT 1 LIMIT 0x10
< error
| 1064 42000 …near '0x10' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) (PARTITION p0 MAX_ROWS = 0x10)
< error
| 1064 42000 …near 'MAX_ROWS = 0x10)' at line 1
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 0; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near '1.5,0))' at line 1
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
| SELECT 1 FROM t WHERE a = 99999999999999999999 OR a = -99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT 18446744073709551616e0, 18446744073709551616.0, 0x10, 99999999999999999999abc FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(.5))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(.5))
< error
| 1064 42000 …near '.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS CHAR(1.))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near '1.5,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 2147483647
< ok 0 0
> exec ^
: string
| EXPLAIN FOR CONNECTION .5
< error
| 1064 42000 …near 'FOR CONNECTION .5' at line 1
> exec ^
: string
| RESET MASTER TO .5
< error
| 1064 42000 …near 'RESET MASTER TO .5' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 254 BUCKETS
< error
| 1046 3D000 No database selected
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2048 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2049 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 65535 BUCKETS' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 16384 INITIAL_SIZE 10
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967295
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 4294967296
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65535
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 65536
< error
| 1046 3D000 No database selected
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
| CREATE TABLE t (a BIT(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIT(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIT(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIT(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIGINT(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIGINT(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIGINT(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BIGINT(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a BIGINT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a YEAR(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a YEAR(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a YEAR(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a YEAR(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a YEAR(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483647,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(2147483648,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775807,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(9223372036854775808,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551615,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …near '18446744073709551616,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …near '99999999999999999999,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(000000000000000000000001,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1.5,0))
< error
| 1064 42000 …near '1.5,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(1e1,0))
< error
| 1064 42000 …near '1e1,0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DECIMAL(10,18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
| CREATE TABLE t (a NUMERIC(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a NUMERIC(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a NUMERIC(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a NUMERIC(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a FIXED(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| CREATE TABLE t (a FIXED(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| CREATE TABLE t (a FIXED(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FIXED(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| CREATE TABLE t (a FIXED(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(1,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615,0))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(1e1,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(10,18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483647,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(2147483648,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775807,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(9223372036854775808,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551615,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(18446744073709551616,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(99999999999999999999,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(000000000000000000000001,0))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1.5,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a DOUBLE(1e1,0))
< error
| 1064 42000 …near ',0))' at line 1
> exec ^
: string
| CREATE TABLE t (a TIME(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIME(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIME(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIME(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIME(18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME(18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TIMESTAMP(18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
| SELECT CAST(a AS DECIMAL(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483647,0))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(2147483648,0))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775807,0))
< error
| 1426 42000 Too big precision 32738 specified for column 'a'. Maximum is 65.
> exec ^
: string
| SELECT CAST(a AS DECIMAL(9223372036854775808,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551615,0))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(18446744073709551616,0))
< error
| 1064 42000 …near '18446744073709551616,0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(99999999999999999999,0))
< error
| 1064 42000 …near '99999999999999999999,0))' at line 1
> exec ^
: string
| SELECT CAST(a AS DECIMAL(000000000000000000000001,0))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(1.5,0))
< error
| 1064 42000 …near '1.5,0))' at line 1
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
| SELECT CAST(a AS DECIMAL(10,2147483647))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,2147483648))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775807))
< error
| 1427 42000 For float(M,D), double(M,D) or decimal(M,D), M must be >= D (column 'a').
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DECIMAL(10,18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
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
| SELECT CAST(a AS FLOAT(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(18446744073709551616))
< error
| 1064 42000 …near '18446744073709551616))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(99999999999999999999))
< error
| 1064 42000 …near '99999999999999999999))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS FLOAT(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec ^
: string
| SELECT CAST(a AS FLOAT(1e1))
< error
| 1064 42000 …near '1e1))' at line 1
> exec ^
: string
| SELECT CAST(a AS TIME(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS TIME(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
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
| SELECT CAST(a AS DATETIME(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775807))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(9223372036854775808))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS DATETIME(18446744073709551615))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
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
| SELECT NOW(1)
< ok 0 0
> exec ^
: string
| SELECT NOW(2147483647)
< ok 0 0
> exec ^
: string
| SELECT NOW(2147483648)
< ok 0 0
> exec ^
: string
| SELECT NOW(9223372036854775807)
< ok 0 0
> exec ^
: string
| SELECT NOW(9223372036854775808)
< ok 0 0
> exec ^
: string
| SELECT NOW(18446744073709551615)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| SELECT CURTIME(2147483648)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(9223372036854775807)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(9223372036854775808)
< ok 0 0
> exec ^
: string
| SELECT CURTIME(18446744073709551615)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(2147483648)
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(9223372036854775807)
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(9223372036854775808)
< ok 0 0
> exec ^
: string
| SELECT SYSDATE(18446744073709551615)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(2147483648)
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(9223372036854775807)
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(9223372036854775808)
< ok 0 0
> exec ^
: string
| SELECT UTC_TIMESTAMP(18446744073709551615)
< ok 0 0
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
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(2147483648)
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(9223372036854775807)
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(9223372036854775808)
< ok 0 0
> exec ^
: string
| SELECT CURRENT_TIMESTAMP(18446744073709551615)
< ok 0 0
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
| DO WEIGHT_STRING(a LEVEL 1)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 2147483647)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 2147483648)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 9223372036854775807)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 9223372036854775808)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 18446744073709551615)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 18446744073709551616)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 18446744073709551616)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 99999999999999999999)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 99999999999999999999)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 000000000000000000000001)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1.5)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1.5)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1e1)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1e1)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 1)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 2147483647)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 2147483648)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 9223372036854775807)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 9223372036854775808)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 18446744073709551615)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 18446744073709551616)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1 - 18446744073709551616)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 99999999999999999999)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1 - 99999999999999999999)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 000000000000000000000001)
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 1.5)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1 - 1.5)' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a LEVEL 1 - 1e1)
< error
| 1064 42000 …near 'WEIGHT_STRING(a LEVEL 1 - 1e1)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a POINT SRID 1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME DEFAULT NOW(18446744073709551615))
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a DATETIME ON UPDATE NOW(18446744073709551615))
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a(18446744073709551615))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(2147483647)))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(2147483648)))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775807)))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(9223372036854775808)))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9), KEY (a(18446744073709551615)))
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
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
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE INDEX i ON t (a) KEY_BLOCK_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 18446744073709551616)
< error
| 1064 42000 …near '18446744073709551616)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 99999999999999999999)
< error
| 1064 42000 …near '99999999999999999999)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1.5)
< error
| 1064 42000 …near '1.5)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, KEY (a) KEY_BLOCK_SIZE 1e1)
< error
| 1064 42000 …near '1e1)' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) KEY_BLOCK_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MAX_ROWS = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) MIN_ROWS = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 18446744073709551615
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PACK_KEYS = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 18446744073709551615
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_AUTO_RECALC = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 18446744073709551615
< error
| 1046 3D000 No database selected
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_PERSISTENT = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) STATS_SAMPLE_PAGES = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY HASH (a) PARTITIONS 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2147483647 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 2147483648 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 9223372036854775807 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 9223372036854775808 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 18446744073709551615 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 18446744073709551616 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …near '18446744073709551616 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …near '99999999999999999999 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 000000000000000000000001 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1.5 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …near '1.5 (PARTITION p0 VALUES LESS THAN (10))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT, b INT) PARTITION BY RANGE (a) SUBPARTITION BY HASH (b) SUBPARTITIONS 1e1 (PARTITION p0 VALUES LESS THAN (10))
< error
| 1064 42000 …near '1e1 (PARTITION p0 VALUES LESS THAN (10))' at line 1
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
| 1046 3D000 No database selected
> exec ^
: string
| DELETE FROM t LIMIT 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| DELETE FROM t LIMIT 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| DELETE FROM t LIMIT 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| DELETE FROM t LIMIT 18446744073709551615
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| UPDATE t SET a = 1 LIMIT 18446744073709551615
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
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
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 1 LINES
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 LINES
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 LINES
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 LINES
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 LINES
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 LINES
< error
| 1017 HY000 File not exist
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
| 1017 HY000 File not exist
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
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483647 ROWS
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 2147483648 ROWS
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775807 ROWS
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 9223372036854775808 ROWS
< error
| 1017 HY000 File not exist
> exec ^
: string
| LOAD DATA INFILE 'f' INTO TABLE t IGNORE 18446744073709551615 ROWS
< error
| 1017 HY000 File not exist
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
| 1017 HY000 File not exist
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
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_CONNECTIONS_PER_HOUR 1e1
< error
| 1064 42000 …near '1e1' at line 1
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
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE USER u WITH MAX_USER_CONNECTIONS 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' FILE_BLOCK_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' AUTOEXTEND_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' MAX_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' EXTENT_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLESPACE ts ADD DATAFILE 'f' NODEGROUP = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' UNDO_BUFFER_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' REDO_BUFFER_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483647
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 2147483648
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775807
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 9223372036854775808
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 18446744073709551615
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 000000000000000000000001
< ok 0 0
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE LOGFILE GROUP lg ADD UNDOFILE 'f' NODEGROUP = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483647
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 2147483648
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775807
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 9223372036854775808
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 18446744073709551615
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 000000000000000000000001
< error
| 1235 0A000 statement type not supported
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| ALTER TABLESPACE ts INITIAL_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483647
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 2147483648
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775807
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 9223372036854775808
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551615
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 000000000000000000000001
< ok 0 0
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| ALTER LOGFILE GROUP lg ADD UNDOFILE 'f' INITIAL_SIZE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483647 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 2147483648 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775807 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 9223372036854775808 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551615 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 18446744073709551616 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 99999999999999999999 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 000000000000000000000001 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1.5 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT 1e1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483647 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -2147483648 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775807 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -9223372036854775808 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551615 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -18446744073709551616 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -99999999999999999999 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -000000000000000000000001 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1.5 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' DEFAULT -1e1 ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
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
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 18446744073709551616)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 99999999999999999999)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1.5)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT 1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -18446744073709551616)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -99999999999999999999)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1.5)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT DEFAULT -1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 1e1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t ALTER COLUMN a SET DEFAULT -1e1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AUTO_INCREMENT = 1e1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| ALTER TABLE t AUTO_INCREMENT = 1e1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1.5))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) PARTITION BY RANGE (a) (PARTITION p0 VALUES LESS THAN (1e1))
< error
| 1046 3D000 No database selected
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
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483647
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 2147483648
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775807
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 9223372036854775808
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551615
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 18446744073709551616
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 99999999999999999999
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 000000000000000000000001
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1.5
< ok 0 0
> exec ^
: string
| SIGNAL SQLSTATE '45000' SET MYSQL_ERRNO = 1e1
< ok 0 0
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483647 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 2147483648 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775807 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 9223372036854775808 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551615 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 18446744073709551616 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 99999999999999999999 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 000000000000000000000001 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1.5 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS 1e1 PRECEDING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483647 PRECEDING AND 2147483647 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 2147483648 PRECEDING AND 2147483648 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775807 PRECEDING AND 9223372036854775807 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 9223372036854775808 PRECEDING AND 9223372036854775808 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551615 PRECEDING AND 18446744073709551615 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 18446744073709551616 PRECEDING AND 18446744073709551616 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 99999999999999999999 PRECEDING AND 99999999999999999999 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 000000000000000000000001 PRECEDING AND 000000000000000000000001 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1.5 PRECEDING AND 1.5 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT SUM(a) OVER (ROWS BETWEEN 1e1 PRECEDING AND 1e1 FOLLOWING) FROM t
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT INTERVAL 1 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 2147483647 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 2147483648 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 9223372036854775807 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 9223372036854775808 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 18446744073709551615 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 18446744073709551616 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 99999999999999999999 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 000000000000000000000001 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 1.5 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT INTERVAL 1e1 DAY + a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT a FROM t ORDER BY 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t ORDER BY 1e1
< error
| 1046 3D000 No database selected
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
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (18446744073709551616)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (99999999999999999999)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (1.5)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t WHERE a IN (1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (18446744073709551616)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (99999999999999999999)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (1.5)
< error
| 1046 3D000 No database selected
> exec ^
: string
| INSERT INTO t VALUES (1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(2147483647)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(2147483648)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(9223372036854775807)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(9223372036854775808)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(18446744073709551615)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(18446744073709551616)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(99999999999999999999)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(000000000000000000000001)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(1.5)
< error
| 1046 3D000 No database selected
> exec ^
: string
| CALL p(1e1)
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 9223372036854775808
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 18446744073709551615
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 18446744073709551616
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 99999999999999999999
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1.5
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT a FROM t FOR UPDATE WAIT 1e1
< error
| 1046 3D000 No database selected
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483647))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(2147483648))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775807))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(9223372036854775808))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(9223372036854775808))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551615))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(18446744073709551615))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(18446744073709551616))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(18446744073709551616))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(99999999999999999999))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(99999999999999999999))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(000000000000000000000001))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1.5))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(1.5))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS CHAR(1e1))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS CHAR(1e1))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483647))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(2147483648))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775807))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(9223372036854775808))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(9223372036854775808))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551615))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(18446744073709551615))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(18446744073709551616))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(18446744073709551616))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(99999999999999999999))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(99999999999999999999))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(000000000000000000000001))
< error
| 1327 42000 Undeclared variable: a
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1.5))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(1.5))' at line 1
> exec ^
: string
| DO WEIGHT_STRING(a AS BINARY(1e1))
< error
| 1064 42000 …near 'WEIGHT_STRING(a AS BINARY(1e1))' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 9223372036854775808
< error
| 1064 42000 …near 'AVG_ROW_LENGTH = 9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551615
< error
| 1064 42000 …near 'AVG_ROW_LENGTH = 18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) AVG_ROW_LENGTH = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 9223372036854775808
< error
| 1064 42000 …near 'CHECKSUM = 9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 18446744073709551615
< error
| 1064 42000 …near 'CHECKSUM = 18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) CHECKSUM = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 2147483647
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 2147483648
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 9223372036854775807
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 9223372036854775808
< error
| 1064 42000 …near 'DELAY_KEY_WRITE = 9223372036854775808' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 18446744073709551615
< error
| 1064 42000 …near 'DELAY_KEY_WRITE = 18446744073709551615' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 18446744073709551616
< error
| 1064 42000 …near '18446744073709551616' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 99999999999999999999
< error
| 1064 42000 …near '99999999999999999999' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 000000000000000000000001
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1.5
< error
| 1064 42000 …near '1.5' at line 1
> exec ^
: string
| CREATE TABLE t (a INT) DELAY_KEY_WRITE = 1e1
< error
| 1064 42000 …near '1e1' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 1
< ok 0 0
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483647
< ok 0 0
> exec ^
: string
| HANDLER t READ FIRST LIMIT 2147483648
< ok 0 0
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775807
< ok 0 0
> exec ^
: string
| HANDLER t READ FIRST LIMIT 9223372036854775808
< error
| 1064 42000 …near 'LIMIT 9223372036854775808' at line 1
> exec ^
: string
| HANDLER t READ FIRST LIMIT 18446744073709551615
< error
| 1064 42000 …near 'LIMIT 18446744073709551615' at line 1
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
< ok 0 0
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
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 1 BUCKETS
< error
| 1046 3D000 No database selected
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483647 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483648 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 2147483648 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775807 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775807 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775808 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 9223372036854775808 BUCKETS' at line 1
> exec ^
: string
| ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551615 BUCKETS
< error
| 1064 42000 …near 'ANALYZE TABLE t UPDATE HISTOGRAM ON a WITH 18446744073709551615 BUCKETS' at line 1
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
< error
| 1046 3D000 No database selected
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
| CREATE TABLE t (a CHAR(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARCHAR(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BINARY(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a VARBINARY(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a BLOB(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a TEXT(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NCHAR(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a NATIONAL VARCHAR(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a CHAR VARYING(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(1))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483647))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(2147483648))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775807))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(9223372036854775808))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551615))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(18446744073709551616))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(99999999999999999999))
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE TABLE t (a FLOAT(000000000000000000000001))
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT CAST(a AS CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NCHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NCHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS NATIONAL CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CAST(a AS BINARY(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT CONVERT(a, CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(1))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(2147483647))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(2147483648))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(9223372036854775807))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(9223372036854775808))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(18446744073709551615))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(18446744073709551616))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(99999999999999999999))
< error
| 1439 42000 Display width out of range for column 'cast of char' (max = 4294967295)
> exec ^
: string
| SELECT JSON_VALUE(a, '$' RETURNING CHAR(000000000000000000000001))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(1) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483647) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(2147483648) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775807) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(9223372036854775808) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551615) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(18446744073709551616) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(99999999999999999999) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| SELECT * FROM JSON_TABLE(a, '$' COLUMNS (b CHAR(000000000000000000000001) PATH '$')) jt
< error
| 1054 42S22 Unknown column 'a' in 'from clause'
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 1 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483647 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 2147483648 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775807 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 9223372036854775808 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551615 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 18446744073709551616 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 99999999999999999999 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE CONTINUE HANDLER FOR 000000000000000000000001 BEGIN END; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 1; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483647; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 2147483648; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775807; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 9223372036854775808; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551615; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 18446744073709551616; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 99999999999999999999; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE PROCEDURE p() BEGIN DECLARE c CONDITION FOR 000000000000000000000001; END
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION f(a CHAR(1)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483647)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(2147483648)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775807)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(9223372036854775808)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551615)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(18446744073709551616)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(99999999999999999999)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
> exec ^
: string
| CREATE FUNCTION f(a CHAR(000000000000000000000001)) RETURNS INT RETURN 1
< error
| 1046 3D000 No database selected
