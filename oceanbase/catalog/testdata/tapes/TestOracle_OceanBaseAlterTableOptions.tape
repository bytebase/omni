-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseAlterTableOptions/a_ROW_FORMAT_=_COMPRESSED_table_takes_semistructured_encoding default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT) ROW_FORMAT = COMPRESSED
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'encoding'
< ok 0 0

== TestOracle_OceanBaseAlterTableOptions/a_dynamic_partition_policy default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT)
< ok 0 0
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'minute')
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(EXPIRE_TIME = '1 hours')
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = '+13:01')
< error
| 1298 HY000 Unknown or incorrect time zone
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = 's', TIME_UNIT = 's')
< error
| 1298 HY000 Unknown or incorrect time zone
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 's', TIME_ZONE = 's')
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour' TIME_UNIT = 'hour')
< error
| 1210 HY000 Invalid argument
> exec
| ALTER TABLE a DYNAMIC_PARTITION_POLICY(TIME_ZONE = '-12:59', PRECREATE_TIME = '+1day')
< error
| 1235 0A000 specify dynamic partition on non partitioned table is not supported
> exec
| CREATE TABLE h (b INT PRIMARY KEY) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE h DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour')
< error
| 1235 0A000 specify dynamic partition on table not partitioned by range is not supported
> exec
| CREATE TABLE r (b DATETIME PRIMARY KEY) PARTITION BY RANGE COLUMNS(b) (PARTITION p VALUES LESS THAN ('2026-01-01'))
< ok 0 0
> exec
| ALTER TABLE r DYNAMIC_PARTITION_POLICY(TIME_UNIT = 'hour')
< ok 0 0

== TestOracle_OceanBaseAlterTableOptions/options_first,_in_order default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT)
< ok 0 0
> exec
| ALTER TABLE a ADD b INT, BLOCK_SIZE 1
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| ALTER TABLE a DROP COLUMN zz, DUPLICATE_SCOPE 's'
< error
| 1210 HY000 Incorrect arguments to duplicate_scope
> exec
| ALTER TABLE a ADD KEY(zz), TABLE_MODE 's'
< error
| 1235 0A000 Table mode s is not supported
> exec
| ALTER TABLE a TABLE_MODE 's', LOCATION 's'
< error
| 1235 0A000 Table mode s is not supported
> exec
| ALTER TABLE a LOCATION 's', TABLE_MODE 's'
< error
| 4016 HY000 Internal error

== TestOracle_OceanBaseAlterTableOptions/values default#1
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT)
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE 'queuing'
< ok 0 0
> exec
| ALTER TABLE a TABLE_MODE ''
< error
| 1235 0A000 Table mode (null) is not supported
> exec
| ALTER TABLE a DUPLICATE_SCOPE 'CLUSTER'
< ok 0 0
> exec
| ALTER TABLE a DUPLICATE_SCOPE ''
< error
| 1210 HY000 Incorrect arguments to duplicate_scope
> exec
| ALTER TABLE a DUPLICATE_READ_CONSISTENCY 'weak'
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE ''
< ok 0 0
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'none'
< error
| 1210 HY000 Incorrect arguments to semistruct encoding type, should be none or encoding
> exec
| ALTER TABLE a SEMISTRUCT_ENCODING_TYPE 'encoding'
< error
| 1235 0A000 semistruct_encoding is not support if cs encoding is not set not supported
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE 'ORDER'
< ok 0 0
> exec
| ALTER TABLE a AUTO_INCREMENT_MODE ''
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your OceanBase version for the right syntax to use
> exec
| ALTER TABLE a KV_ATTRIBUTES '{"Hbase": {}}'
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1024
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1048576
< ok 0 0
> exec
| ALTER TABLE a BLOCK_SIZE 1048577
< error
| 5258 HY000 Invalid block size, block size should between 1024 and 1048576
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 4
< ok 0 0
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 3
< error
| 4179 HY000 target storage format version cannot be smaller than current version not allowed
> exec
| ALTER TABLE a STORAGE_FORMAT_VERSION 5
< error
| 1210 HY000 Incorrect arguments to invalid storage format version
> exec
| ALTER TABLE a TABLET_SIZE 0
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 4194304
< ok 0 0
> exec
| ALTER TABLE a TABLET_SIZE 3145728
< error
| 4147 HY000 Invalid config
> exec
| ALTER TABLE a ORGANIZATION HEAP
< error
| 1235 0A000 organization in the alter table statement not supported
> exec
| ALTER TABLE a FORMAT(TYPE = 'csv')
< error
| 4016 HY000 Internal error
