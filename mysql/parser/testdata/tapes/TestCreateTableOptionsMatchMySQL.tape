-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestCreateTableOptionsMatchMySQL default#1
> exec
| create table t (c int) table_checksum = 0
< ok 0 0
> exec
| create table t (c int) table_checksum 1
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c int) table_checksum 1
< ok 0 0
> exec
| create table t (c int) TABLE_CHECKSUM = 2
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c int) TABLE_CHECKSUM = 2
< ok 0 0
> exec
| create table t (c int) table_checksum = 1.5
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c int) table_checksum = 1.5
< ok 0 0
> exec
| create table t (c int) table_checksum
< error
| 1064 42000 …near '' at line 1
> exec
| create table t (c int) table_checksum =
< error
| 1064 42000 …near '' at line 1
> exec
| create table t (c int) table_checksum = 'abc'
< error
| 1064 42000 …near ''abc'' at line 1
> exec
| create table t (c int) table_checksum = -1
< error
| 1064 42000 …near '-1' at line 1
> exec
| create table t (c int) TABLE_CHECKSUM = DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| create table t (c int) checksum = 1
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c int) checksum = 1
< ok 0 0
> exec
| create table t (c int) checksum = 1.5
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c int) checksum = 1.5
< ok 0 0
> exec
| create table t (c int) checksum = 'abc'
< error
| 1064 42000 …near ''abc'' at line 1
> exec
| create table t (c int) checksum = DEFAULT
< error
| 1064 42000 …near 'DEFAULT' at line 1
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC=0
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=0
< ok 0 0
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC=1
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=1
< ok 0 0
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC=DEFAULT
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=DEFAULT
< ok 0 0
> exec
| CREATE TABLE t (a int) stats_auto_recalc default
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a int) stats_auto_recalc default
< ok 0 0
> exec
| CREATE TABLE t (a int) stats_auto_recalc 2
< error
| 1064 42000 …near '2' at line 1
> exec
| CREATE TABLE t (a int) stats_auto_recalc = 10
< error
| 1064 42000 …near '10' at line 1
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'
< error
| 1064 42000 …near ''abc'' at line 1
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC='1'
< error
| 1064 42000 …near ''1'' at line 1
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC=-1
< error
| 1064 42000 …near '-1' at line 1
> exec
| CREATE TABLE t (a int) STATS_AUTO_RECALC=
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (a int) STATS_BUCKETS=1
< error
| 1064 42000 …near 'STATS_BUCKETS=1' at line 1
> exec
| CREATE TABLE t (a int) STATS_TOPN=1
< error
| 1064 42000 …near 'STATS_TOPN=1' at line 1
> exec
| CREATE TABLE t (a int) STATS_SAMPLE_RATE=0.1
< error
| 1064 42000 …near 'STATS_SAMPLE_RATE=0.1' at line 1
> exec
| CREATE TABLE t (a int) STATS_COL_CHOICE='all'
< error
| 1064 42000 …near 'STATS_COL_CHOICE='all'' at line 1
> exec
| CREATE TABLE t (a int, b int) STATS_COL_LIST='a,b'
< error
| 1064 42000 …near 'STATS_COL_LIST='a,b'' at line 1
> exec
| create table t (a int serial default value)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value)
< ok 0 0
> exec
| create table t (a bigint serial default value)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a bigint serial default value)
< ok 0 0
> exec
| create table t (a varchar(3) serial default value)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a varchar(3) serial default value)
< error
| 1063 42000 Incorrect column specifier for column 'a'
> exec ^
: string
| create table t (a varchar(3) serial default value)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| create table t (a int serial default value null)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value null)
< ok 0 0
> exec
| create table t (a int serial default value not null)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value not null)
< ok 0 0
> exec
| create table t (a int serial default value primary key)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value primary key)
< ok 0 0
> exec
| create table t (a int unique serial default value)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int unique serial default value)
< ok 0 0
> exec
| create table t (a int comment 'x' serial default value)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int comment 'x' serial default value)
< ok 0 0
> exec
| create table t (a int serial default value serial default value)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value serial default value)
< ok 0 0
> exec
| create table t (a int serial default value, b int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (a int serial default value, b int)
< ok 0 0
> exec
| create table t (a int serial)
< error
| 1064 42000 …near ')' at line 1
> exec
| create table t (a int serial default)
< error
| 1064 42000 …near ')' at line 1
> exec
| create table t (a int serial value)
< error
| 1064 42000 …near 'value)' at line 1
> exec
| CREATE TABLE t (a SQL_TSI_YEAR)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR(4) ZEROFILL)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR UNSIGNED)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR(4) unsigned zerofill)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) unsigned zerofill)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR(4) not null)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) not null)
< ok 0 0
> exec
| CREATE TABLE t (a YEAR ZEROFILL)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a YEAR ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t (a YEAR(4) ZEROFILL)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a YEAR(4) ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t (a SQL_TSI_YEAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| create table t (a year(4) binary)
< error
| 1064 42000 …near 'binary)' at line 1
> exec
| CREATE TABLE t (a YEAR())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a int())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a char())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a datetime())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a bit())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a blob())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a binary())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a decimal())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a float())
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a decimal(10,))
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a decimal(,2))
< error
| 1064 42000 …near ',2))' at line 1
> exec
| CREATE TABLE t (a float(10,))
< error
| 1064 42000 …near '))' at line 1
> exec
| CREATE TABLE t (a int(,))
< error
| 1064 42000 …near ',))' at line 1
> exec
| CREATE TABLE t (a int(1))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a int(1))
< ok 0 0
> exec
| CREATE TABLE t (a double(10,2))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (a double(10,2))
< ok 0 0
> exec
| create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))
< ok 0 0
> exec
| create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))
< ok 0 0
> exec
| create table t (c1 enum(x'61'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(x'61'))
< ok 0 0
> exec
| create table t (c1 enum(b'01100001'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(b'01100001'))
< ok 0 0
> exec
| create table t (c1 enum(0x6162))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x6162))
< ok 0 0
> exec
| create table t (c1 set(0b0110000101100010))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 set(0b0110000101100010))
< ok 0 0
> exec
| create table t (c1 enum(0x61) charset utf8mb4)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x61) charset utf8mb4)
< ok 0 0
> exec
| create table t (c1 enum(0x61, 'b') character set utf8mb4)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x61, 'b') character set utf8mb4)
< ok 0 0
> exec
| create table t (c1 enum())
< error
| 1064 42000 …near '))' at line 1
> exec
| create table t (c1 set())
< error
| 1064 42000 …near '))' at line 1
> exec
| create table t (c1 enum(1))
< error
| 1064 42000 …near '1))' at line 1
> exec
| create table t (c1 enum(1.5))
< error
| 1064 42000 …near '1.5))' at line 1
> exec
| create table t (c1 enum('a',))
< error
| 1064 42000 …near '))' at line 1
> exec
| create table t (c1 enum(0x))
< error
| 1064 42000 …near '0x))' at line 1
> exec
| create table t (c1 enum(0b))
< error
| 1064 42000 …near '0b))' at line 1
> exec
| create table t (c1 enum(0xZZ))
< error
| 1064 42000 …near '0xZZ))' at line 1
> exec
| create table t (c1 enum(0x6))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x6))
< ok 0 0
> exec
| create table t (c1 enum(0x616))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0x616))
< ok 0 0
> exec
| create table t (c1 enum(0b1))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0b1))
< ok 0 0
> exec
| create table t (c1 enum(0b101))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(0b101))
< ok 0 0
> exec
| create table t (c1 enum(x'0b61'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(x'0b61'))
< ok 0 0
> exec
| create table t (c1 enum(x'0b'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(x'0b'))
< ok 0 0
> exec
| create table t (c1 enum(x'0B61'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(x'0B61'))
< ok 0 0
> exec
| create table t (c1 set(x'0b61'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 set(x'0b61'))
< ok 0 0
> exec
| create table t (c1 enum(b'1'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(b'1'))
< ok 0 0
> exec
| create table t (c1 enum(b'101'))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(b'101'))
< ok 0 0
> exec
| create table t (c1 enum(x''))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(x''))
< ok 0 0
> exec
| create table t (c1 enum(b''))
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (c1 enum(b''))
< ok 0 0
> exec
| create table t (c1 enum(x'6'))
< error
| 1064 42000 …near 'x'6'))' at line 1
> exec
| create table t (c1 enum(x'616'))
< error
| 1064 42000 …near 'x'616'))' at line 1
> exec
| create table t (c1 enum(x'0x61'))
< error
| 1064 42000 …near 'x'0x61'))' at line 1
> exec
| create table t (c1 enum(b'0b01'))
< error
| 1064 42000 …near 'b'0b01'))' at line 1
> exec
| create table `123` (123a1 int)
< ok 0 0
> exec
| create table t (123a1 int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (123a1 int)
< ok 0 0
> exec
| create table t (1a int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (1a int)
< ok 0 0
> exec
| create table t (123abc int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (123abc int)
< ok 0 0
> exec
| create table t (12e int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (12e int)
< ok 0 0
> exec
| create table t (1_2 int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (1_2 int)
< ok 0 0
> exec
| create table t (123a1 int, 456b2 int)
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table t (123a1 int, 456b2 int)
< ok 0 0
> exec
| create table 1t (a int)
< ok 0 0
> exec
| create table 9t.t1 (a int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table 't1'
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| create table 9t.t1 (a int)
< error
| 1142 42000 CREATE command denied to user 'omni_test'@'172.17.0.1' for table 't1'
> exec ^
: string
| create table 9t.t1 (a int)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| create table 123 (123a1 int)
< error
| 1064 42000 …near '123 (123a1 int)' at line 1
> exec
| create table t (1e1 int)
< error
| 1064 42000 …near '1e1 int)' at line 1
> exec
| create table t (0x61 int)
< error
| 1064 42000 …near '0x61 int)' at line 1
> exec
| select 1e1 from t
< ok 0 0
> exec
| select 1E+10 from t
< ok 0 0
> exec
| select 1.5e3 from t
< ok 0 0
> exec
| select 1e+ from t
< error
| 1064 42000 …near 'from t' at line 1
> exec
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< ok 0 0
> exec
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE=
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE
< error
| 1064 42000 …near '' at line 1
> exec
| CREATE TABLE t (id INT) TABLESPACE ts1 SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 3510 HY000 Tablespace ts1 doesn't exist.
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (id INT) TABLESPACE ts1 SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< ok 0 0
> exec
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 1050 42S01 Table 't' already exists
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')
< ok 0 0
> exec
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE=)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 3981 HY000 Storage engine 'InnoDB' does not support ENGINE_ATTRIBUTE.
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'
< ok 0 0
> exec
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 3981 HY000 Storage engine 'InnoDB' does not support ENGINE_ATTRIBUTE.
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE='{"key":"value"}')
< ok 0 0
> exec
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE)
< error
| 1064 42000 …near ')' at line 1
> exec
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE=)
< error
| 1064 42000 …near ')' at line 1
> exec
| create global temporary table t (a int) on commit delete rows
< error
| 1064 42000 …near 'global temporary table t (a int) on commit delete rows' at line 1
> exec
| create temporary table t (a int) on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec
| create table t (a int) on commit delete rows
< error
| 1064 42000 …near 'on commit delete rows' at line 1
> exec
| create global table t (a int)
< error
| 1064 42000 …near 'global table t (a int)' at line 1
> exec
| drop global temporary table t
< error
| 1064 42000 …near 'global temporary table t' at line 1
> exec
| create global temporary view v as select 1
< error
| 1064 42000 …near 'global temporary view v as select 1' at line 1
> exec
| create global temporary database d
< error
| 1064 42000 …near 'global temporary database d' at line 1
> exec
| drop global temporary database d
< error
| 1064 42000 …near 'global temporary database d' at line 1
> exec
| create temporary table t (a int)
< ok 0 0
> exec
| drop temporary table t
< ok 0 0
