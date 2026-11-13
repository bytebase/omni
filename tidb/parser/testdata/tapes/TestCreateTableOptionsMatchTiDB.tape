-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestCreateTableOptionsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| create table t (c int) table_checksum = 0
< ok 0 0
> exec ^
: string
| create table t (c int) table_checksum 1
< ok 0 0
> exec ^
: string
| create table t (c int) TABLE_CHECKSUM = 2
< ok 0 0
> exec ^
: string
| create table t (c int) checksum = 1
< ok 0 0
> exec ^
: string
| create table t (c int) checksum 0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_TOPN=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_TOPN 0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_SAMPLE_RATE=0.1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_SAMPLE_RATE 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE='all'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE 'list'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE='whatever'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b int) STATS_COL_LIST='a,b'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int, b int) STATS_COL_LIST 'a'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=1 STATS_TOPN=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE 'all' STATS_COL_LIST 'a'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1 PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (200))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=0
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) stats_auto_recalc default
< ok 0 0
> exec ^
: string
| create table t (a int serial default value)
< ok 0 0
> exec ^
: string
| create table t (a bigint serial default value)
< ok 0 0
> exec ^
: string
| create table t (a varchar(3) serial default value)
< error
| 1105 HY000 Incorrect column specifier for column 'a'
> exec ^
: string
| create table t (a varchar(3) serial default value)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int serial default value null)
< ok 0 0
> exec ^
: string
| create table t (a int serial default value not null)
< ok 0 0
> exec ^
: string
| create table t (a int serial default value primary key)
< ok 0 0
> exec ^
: string
| create table t (a int unique serial default value)
< ok 0 0
> exec ^
: string
| create table t (a int comment 'x' serial default value)
< ok 0 0
> exec ^
: string
| create table t (a int serial default value serial default value)
< ok 0 0
> exec ^
: string
| create table t (a int serial default value, b int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) ZEROFILL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR UNSIGNED)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) unsigned zerofill)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR(4) not null)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR ZEROFILL)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a YEAR(4) ZEROFILL)
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(x'61'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(b'01100001'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x6162))
< ok 0 0
> exec ^
: string
| create table t (c1 set(0b0110000101100010))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x61) charset utf8mb4)
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x61, 'b') character set utf8mb4)
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x6))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0x616))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0b1))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(0b101))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(x'0b61'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(x'0b'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(x'0B61'))
< ok 0 0
> exec ^
: string
| create table t (c1 set(x'0b61'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(b'1'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(b'101'))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(x''))
< ok 0 0
> exec ^
: string
| create table t (c1 enum(b''))
< ok 0 0
> exec ^
: string
| create table t (c1 enum('a', 0x61))
< ok 0 0
> exec ^
: string
| create table `123` (123a1 int)
< ok 0 0
> exec ^
: string
| create table t (123a1 int)
< ok 0 0
> exec ^
: string
| create table t (1a int)
< ok 0 0
> exec ^
: string
| create table t (123abc int)
< ok 0 0
> exec ^
: string
| create table t (12e int)
< ok 0 0
> exec ^
: string
| create table t (1_2 int)
< ok 0 0
> exec ^
: string
| create table t (123a1 int, 456b2 int)
< ok 0 0
> exec ^
: string
| create table 1t (a int)
< ok 0 0
> exec ^
: string
| create table 9t.t1 (a int)
< ok 0 0
> exec ^
: string
| create global temporary table t (a int, b varchar(255)) on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t (a int) ON COMMIT PRESERVE ROWS
< ok 0 0
> exec ^
: string
| create global temporary table if not exists t (a int) on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t (a int, b varchar(255)) partition by hash(a) partitions 10 on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t (a int) comment 'x' on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t like t1 on commit delete rows
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| create global temporary table t like t1 on commit delete rows
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop global temporary table t
< ok 0 0
> exec ^
: string
| drop global temporary table if exists t
< ok 0 0
> exec ^
: string
| drop global temporary table t1, t2
< ok 0 0
> exec ^
: string
| create temporary table t (a int, b varchar(255))
< ok 0 0
> exec ^
: string
| drop temporary table t
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE 'innodb'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE = innodb
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE innodb
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) NODEGROUP 1
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) NODEGROUP = 2
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) AFFINITY 's'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) AFFINITY = 's'
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) AFFINITY foo
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s STORAGE ENGINE 's'
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s STORAGE ENGINE = innodb
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s NODEGROUP 1
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s NODEGROUP = 2
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AFFINITY 's'
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AFFINITY foo
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s STORAGE MEMORY
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s ENGINE 'innodb' NODEGROUP 1 AFFINITY 's'
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s START 1 STORAGE ENGINE 'innodb', AFFINITY 's'
< ok 0 0
> exec ^
: string
| CREATE SEQUENCE s AFFINITY `select`
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE `select`
< error
| 1286 42000 Unknown storage engine 'select'
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE `select`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (c int) table_checksum
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| create table t (c int) table_checksum =
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| create table t (c int) table_checksum = 'abc'
< error
| 1064 42000 …line 1 column 45 near "'abc'" 
> exec ^
: string
| create table t (c int) table_checksum = -1
< error
| 1064 42000 …line 1 column 41 near "-1" 
> exec ^
: string
| create table t (c int) TABLE_CHECKSUM = DEFAULT
< error
| 1064 42000 …line 1 column 47 near "DEFAULT" 
> exec ^
: string
| create table t (c int) table_checksum = 1.5
< error
| 1064 42000 …line 1 column 43 near "1.5" 
> exec ^
: string
| create table t (c int) checksum = 'abc'
< error
| 1064 42000 …line 1 column 39 near "'abc'" 
> exec ^
: string
| create table t (c int) checksum = DEFAULT
< error
| 1064 42000 …line 1 column 41 near "DEFAULT" 
> exec ^
: string
| create table t (c int) checksum = 1.5
< error
| 1064 42000 …line 1 column 37 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=
< error
| 1064 42000 …line 1 column 37 near "" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS='abc'
< error
| 1064 42000 …line 1 column 42 near "'abc'" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS="1"
< error
| 1064 42000 …line 1 column 40 near ""1"" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=1.5
< error
| 1064 42000 …line 1 column 40 near "1.5" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=-1
< error
| 1064 42000 …line 1 column 38 near "-1" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_BUCKETS=DEFAULT
< error
| 1064 42000 …line 1 column 44 near "DEFAULT" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_TOPN
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_TOPN='abc'
< error
| 1064 42000 …line 1 column 39 near "'abc'" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_SAMPLE_RATE='abc'
< error
| 1064 42000 …line 1 column 46 near "'abc'" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_SAMPLE_RATE=-0.1
< error
| 1064 42000 …line 1 column 42 near "-0.1" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_SAMPLE_RATE=DEFAULT
< error
| 1064 42000 …line 1 column 48 near "DEFAULT" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE=
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE=1
< error
| 1064 42000 …line 1 column 41 near "1" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_COL_CHOICE=all
< error
| 1064 42000 …line 1 column 43 near "all" 
> exec ^
: string
| CREATE TABLE t (a int, b int) STATS_COL_LIST=1
< error
| 1064 42000 …line 1 column 46 near "1" 
> exec ^
: string
| create table t (a int) STATS_OPTIONS='str'
< error
| 1064 42000 …line 1 column 36 near "STATS_OPTIONS='str'" 
> exec ^
: string
| CREATE TABLE t (a int) stats_auto_recalc 2
< error
| 1064 42000 …line 1 column 42 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a int) stats_auto_recalc = 10
< error
| 1064 42000 …line 1 column 45 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'
< error
| 1064 42000 …line 1 column 46 near "'abc'" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC='1'
< error
| 1064 42000 …line 1 column 44 near "'1'" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=-1
< error
| 1064 42000 …line 1 column 42 near "-1" 
> exec ^
: string
| CREATE TABLE t (a int) STATS_AUTO_RECALC=
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| create table t (a int serial)
< error
| 1064 42000 …line 1 column 29 near ")" 
> exec ^
: string
| create table t (a int serial default)
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec ^
: string
| create table t (a int serial value)
< error
| 1064 42000 …line 1 column 34 near "value)" 
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_YEAR())
< error
| 1064 42000 …line 1 column 32 near "))" 
> exec ^
: string
| CREATE TABLE t (a YEAR())
< error
| 1064 42000 …line 1 column 24 near "))" 
> exec ^
: string
| CREATE TABLE t (a int())
< error
| 1064 42000 …line 1 column 23 near "))" 
> exec ^
: string
| CREATE TABLE t (a char())
< error
| 1064 42000 …line 1 column 24 near "))" 
> exec ^
: string
| CREATE TABLE t (a datetime())
< error
| 1064 42000 …line 1 column 28 near "))" 
> exec ^
: string
| CREATE TABLE t (a bit())
< error
| 1064 42000 …line 1 column 23 near "))" 
> exec ^
: string
| CREATE TABLE t (a blob())
< error
| 1064 42000 …line 1 column 24 near "))" 
> exec ^
: string
| CREATE TABLE t (a binary())
< error
| 1064 42000 …line 1 column 26 near "))" 
> exec ^
: string
| CREATE TABLE t (a decimal())
< error
| 1064 42000 …line 1 column 27 near "))" 
> exec ^
: string
| CREATE TABLE t (a float())
< error
| 1064 42000 …line 1 column 25 near "))" 
> exec ^
: string
| CREATE TABLE t (a decimal(10,))
< error
| 1064 42000 …line 1 column 30 near "))" 
> exec ^
: string
| CREATE TABLE t (a decimal(,2))
< error
| 1064 42000 …line 1 column 27 near ",2))" 
> exec ^
: string
| CREATE TABLE t (a float(10,))
< error
| 1064 42000 …line 1 column 28 near "))" 
> exec ^
: string
| CREATE TABLE t (a int(,))
< error
| 1064 42000 …line 1 column 23 near ",))" 
> exec ^
: string
| CREATE TABLE t (a SQL_TSI_DAY)
< error
| 1064 42000 …line 1 column 29 near "SQL_TSI_DAY)" 
> exec ^
: string
| create table t (a year(4) binary)
< error
| 1064 42000 …line 1 column 32 near "binary)" 
> exec ^
: string
| create table t (c1 enum())
< error
| 1064 42000 …line 1 column 25 near "))" 
> exec ^
: string
| create table t (c1 set())
< error
| 1064 42000 …line 1 column 24 near "))" 
> exec ^
: string
| create table t (c1 enum(1))
< error
| 1064 42000 …line 1 column 25 near "1))" 
> exec ^
: string
| create table t (c1 enum(1.5))
< error
| 1064 42000 …line 1 column 27 near "1.5))" 
> exec ^
: string
| create table t (c1 enum('a',))
< error
| 1064 42000 …line 1 column 29 near "))" 
> exec ^
: string
| create table t (c1 enum(0x))
< error
| 1064 42000 …line 1 column 26 near "0x))" 
> exec ^
: string
| create table t (c1 enum(0b))
< error
| 1064 42000 …line 1 column 26 near "0b))" 
> exec ^
: string
| create table t (c1 enum(0xZZ))
< error
| 1064 42000 …line 1 column 28 near "0xZZ))" 
> exec ^
: string
| create table t (c1 enum(x'6'))
< error
| 1064 42000 …line 1 column 28 near "x'6'))"hex literal: invalid hexadecimal format, must even numbers, but 1 
> exec ^
: string
| create table t (c1 enum(x'616'))
< error
| 1064 42000 …line 1 column 30 near "x'616'))"hex literal: invalid hexadecimal format, must even numbers, but 3 
> exec ^
: string
| create table t (c1 enum(x'0x61'))
< error
| 1064 42000 …line 1 column 27 near "x'0x61'))" 
> exec ^
: string
| create table t (c1 enum(b'0b01'))
< error
| 1064 42000 …line 1 column 27 near "b'0b01'))" 
> exec ^
: string
| create table 123 (123a1 int)
< error
| 1064 42000 …line 1 column 16 near "123 (123a1 int)" 
> exec ^
: string
| create table t (1e1 int)
< error
| 1064 42000 …line 1 column 19 near "1e1 int)" 
> exec ^
: string
| create table t (0x61 int)
< error
| 1064 42000 …line 1 column 20 near "0x61 int)" 
> exec ^
: string
| create global temporary table t (a int, b varchar(255))
< error
| 1064 42000 …line 1 column 55 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create temporary table t (a int, b varchar(255)) on commit delete rows
< error
| 1064 42000 …line 1 column 70 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t (a int, b varchar(255)) on commit delete rows
< error
| 1064 42000 …line 1 column 60 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t (a int) on commit
< error
| 1064 42000 …line 1 column 49 near "" 
> exec ^
: string
| create global temporary table t (a int) on commit rows
< error
| 1064 42000 …line 1 column 54 near "rows" 
> exec ^
: string
| create global temporary table t (a int) on commit delete
< error
| 1064 42000 …line 1 column 56 near "" 
> exec ^
: string
| create global temporary table t (a int) on commit preserve
< error
| 1064 42000 …line 1 column 58 near "" 
> exec ^
: string
| create global temporary table t (a int) on commit delete rows comment 'x'
< error
| 1064 42000 …line 1 column 69 near "comment 'x'" 
> exec ^
: string
| create global temporary table t (a int) on commit preserve rows partition by hash(a)
< error
| 1064 42000 …line 1 column 73 near "partition by hash(a)" 
> exec ^
: string
| create global temporary table t like t1
< error
| 1064 42000 …line 1 column 39 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global table t (a int)
< error
| 1064 42000 …line 1 column 19 near "table t (a int)" 
> exec ^
: string
| drop global table t
< error
| 1064 42000 …line 1 column 17 near "table t" 
> exec ^
: string
| create global temporary view v as select 1
< error
| 1064 42000 …line 1 column 28 near "view v as select 1" 
> exec ^
: string
| create global temporary index i on t (a)
< error
| 1064 42000 …line 1 column 29 near "index i on t (a)" 
> exec ^
: string
| create global temporary database d
< error
| 1064 42000 …line 1 column 32 near "database d" 
> exec ^
: string
| create global temporary sequence s
< error
| 1064 42000 …line 1 column 32 near "sequence s" 
> exec ^
: string
| drop global temporary database d
< error
| 1064 42000 …line 1 column 30 near "database d" 
> exec ^
: string
| drop global temporary view v
< error
| 1064 42000 …line 1 column 26 near "view v" 
> exec ^
: string
| drop global temporary index i on t
< error
| 1064 42000 …line 1 column 27 near "index i on t" 
> exec ^
: string
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1064 42000 …line 1 column 50 near "SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'" 
> exec ^
: string
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE=
< error
| 1064 42000 …line 1 column 50 near "SECONDARY_ENGINE_ATTRIBUTE=" 
> exec ^
: string
| CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE
< error
| 1064 42000 …line 1 column 50 near "SECONDARY_ENGINE_ATTRIBUTE" 
> exec ^
: string
| CREATE TABLE t (id INT) TABLESPACE ts1 SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1064 42000 …line 1 column 65 near "SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'" 
> exec ^
: string
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 1064 42000 …line 1 column 49 near "SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')" 
> exec ^
: string
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE=)
< error
| 1064 42000 …line 1 column 49 near "SECONDARY_ENGINE_ATTRIBUTE=)" 
> exec ^
: string
| CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE)
< error
| 1064 42000 …line 1 column 49 near "SECONDARY_ENGINE_ATTRIBUTE)" 
> exec ^
: string
| CREATE TABLE t (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'
< error
| 1064 42000 …line 1 column 40 near "ENGINE_ATTRIBUTE='{"key":"value"}'" 
> exec ^
: string
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE='{"key":"value"}')
< error
| 1064 42000 …line 1 column 39 near "ENGINE_ATTRIBUTE='{"key":"value"}')" 
> exec ^
: string
| CREATE TABLE t (id INT ENGINE_ATTRIBUTE)
< error
| 1064 42000 …line 1 column 39 near "ENGINE_ATTRIBUTE)" 
> exec ^
: string
| ALTER TABLE t SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'
< error
| 1064 42000 …line 1 column 40 near "SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'" 
> exec ^
: string
| alter table t ENGINE_ATTRIBUTE='{"a":1}'
< error
| 1064 42000 …line 1 column 30 near "ENGINE_ATTRIBUTE='{"a":1}'" 
> exec ^
: string
| alter table t add column c int SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'
< error
| 1064 42000 …line 1 column 57 near "SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'" 
> exec ^
: string
| CREATE SEQUENCE s STORAGE ENGINE 1
< error
| 1064 42000 …line 1 column 34 near "1" 
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE 1
< error
| 1064 42000 …line 1 column 39 near "1" 
> exec ^
: string
| CREATE SEQUENCE s AFFINITY 1
< error
| 1064 42000 …line 1 column 28 near "1" 
> exec ^
: string
| CREATE TABLE t (a int) AFFINITY 1
< error
| 1064 42000 …line 1 column 33 near "1" 
> exec ^
: string
| CREATE SEQUENCE s NODEGROUP 's'
< error
| 1064 42000 …line 1 column 31 near "'s'" 
> exec ^
: string
| CREATE TABLE t (a int) NODEGROUP 's'
< error
| 1064 42000 …line 1 column 36 near "'s'" 
> exec ^
: string
| CREATE SEQUENCE s STORAGE ENGINE
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| CREATE SEQUENCE s AFFINITY
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| CREATE SEQUENCE s NODEGROUP
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| CREATE SEQUENCE s AFFINITY SELECT
< error
| 1064 42000 …line 1 column 33 near "SELECT" 
> exec ^
: string
| CREATE TABLE t (a int) STORAGE ENGINE SELECT
< error
| 1064 42000 …line 1 column 44 near "SELECT" 
> exec ^
: string
| CREATE SEQUENCE s AFFINITY = SELECT
< error
| 1064 42000 …line 1 column 35 near "SELECT" 
> exec ^
: string
| CREATE SEQUENCE s STORAGE 's'
< error
| 1064 42000 …line 1 column 29 near "'s'" 
> exec ^
: string
| CREATE SEQUENCE s TTL + INTERVAL FALSE SECOND_MICROSECOND
< error
| 1064 42000 …line 1 column 23 near "+ INTERVAL FALSE SECOND_MICROSECOND" 
> exec ^
: string
| CREATE SEQUENCE s STATS_PERSISTENT DEFINED
< error
| 1064 42000 …line 1 column 42 near "DEFINED" 
> exec ^
: string
| ALTER TABLE t STATS_PERSISTENT 's'
< error
| 1064 42000 …line 1 column 34 near "'s'" 
