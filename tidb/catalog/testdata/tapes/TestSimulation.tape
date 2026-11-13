-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSimulation default#1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT TIFLASH REPLICA
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT PARTITION ident
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident TIFLASH REPLICA
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE PARTITIONING
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME PARTITIONING
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SPLIT MAXVALUE PARTITION LESS THAN(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MERGE FIRST PARTITION LESS THAN(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident ATTRIBUTES DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE TTL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TTL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET TIFLASH REPLICA 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET HYPO TIFLASH REPLICA 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD(ident TINYINT)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD INDEX(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION PARTITIONS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LAST PARTITION LESS THAN(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COALESCE PARTITION 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PRIMARY KEY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PARTITION ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FIRST PARTITION LESS THAN(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP STATS_EXTENDED ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TRUNCATE PARTITION ALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident OPTIMIZE PARTITION ALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REPAIR PARTITION ALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT PARTITION ALL TABLESPACE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD PARTITION ALL TABLESPACE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT TABLESPACE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD TABLESPACE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REBUILD PARTITION ALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP KEY ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP FOREIGN KEY ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISABLE KEYS
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENABLE KEYS
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MODIFY ident TINYINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHANGE ident ident TINYINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT FALSE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident DROP DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME COLUMN ident TO ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TO ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME AS ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME KEY ident TO ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LOCK DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ WRITE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITH VALIDATION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITHOUT VALIDATION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_LOAD
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_UNLOAD
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident ENFORCED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CHECK ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident VISIBLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CACHE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NOCACHE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION ident INTO(PARTITION ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITH VALIDATION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITHOUT VALIDATION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY CLUSTERED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY NONCLUSTERED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY LOCAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY GLOBAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM COPY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INPLACE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INSTANT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ ONLY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP INDEX ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT zzz(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY zzz(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP COLUMN ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT FIRST
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AFTER ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE, FORCE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT FULLTEXT(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT zzz FULLTEXT(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT ident FULLTEXT(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident, ident TO ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SERIAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT PRIMARY KEY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident NOT ENFORCED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident) ENFORCED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NOT NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NOT NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_INCREMENT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE GLOBAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE LOCAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE KEY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SERIAL DEFAULT VALUE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLLATE 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1, 1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DISK
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE MEMORY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT FIXED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT GENERATED ALWAYS AS(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) VIRTUAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) STORED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PRIMARY KEY(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD UNIQUE(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FOREIGN KEY(ident) REFERENCES ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH FULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH PARTIAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH SIMPLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT ON UPDATE RESTRICT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT ON DELETE RESTRICT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE CASCADE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE NO ACTION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT FALSE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NEXT VALUE FOR ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT ident()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURDATE()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(NEXT VALUE FOR ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1e0
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1.5
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident DEPENDENCY(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CORRELATION(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident ON ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident, ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT((ident))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT ALGORITHM DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT LOCK DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX ident ON ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX ident ON ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ident ON ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident LIKE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION SYSTEM_TIME
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY HASH(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY ALGORITHM = 1()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident MICROSECOND)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(MICROSECOND)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) NULL PARTITION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) MAXVALUE PARTITION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) FIRST PARTITION LESS THAN(ident) LAST PARTITION LESS THAN(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY LINEAR KEY()
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident, PARTITION ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident, SUBPARTITION ident))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INSERT_METHOD 'binary'
< error
| 8233 HY000 CREATE/ALTER table with insert method option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DATA DIRECTORY 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INDEX DIRECTORY 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MAX_ROWS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MIN_ROWS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLESPACE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NODEGROUP 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PLACEMENT POLICY 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN MAXVALUE)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident DEFAULT)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident HISTORY)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident CURRENT)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident IGNORE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident REPLACE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident AS
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT * UNION SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident WITH ident AS(SELECT *) SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS WITH ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(LIKE ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = 'binary' VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident, ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH CASCADED CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH LOCAL CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident
< error
| 1051 42S02 Unknown table 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TEMPORARY TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW ident
< error
| 1051 42S02 Unknown table 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW IF EXISTS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident RESTRICT
< error
| 1051 42S02 Unknown table 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident CASCADE
< error
| 1051 42S02 Unknown table 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(ident))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(ident))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE, MAXVALUE))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT, DEFAULT))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED, DEFAULT))
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) VISIBLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) KEY_BLOCK_SIZE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING BTREE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) WITH PARSER ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) CLUSTERED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) GLOBAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) LOCAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY zzz(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY USING BTREE(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY BTREE(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY ident TYPE BTREE(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident USING BTREE ON ident(ident)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) TYPE BTREE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HASH
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING RTREE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HYPO
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HNSW
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident INVISIBLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident, ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident ASC
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident DESC
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(FULLTEXT(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(INDEX(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CONSTRAINT ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `ident` tinyint DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
< row
: bytes
| t
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `ident` tinytext DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SHOW CREATE TABLE `t`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT, ident TINYINT)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT, ident TINYINT)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHARSET 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COLLATE 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_INCREMENT 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_ID_CACHE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_RANDOM_BASE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AVG_ROW_LENGTH 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONNECTION 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECKSUM 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLE_CHECKSUM 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PASSWORD 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPRESSION 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident KEY_BLOCK_SIZE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DELAY_KEY_WRITE 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_BUCKETS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_TOPN 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_CHOICE 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_LIST 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SHARD_ROW_ID_BITS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PRE_SPLIT_REGIONS 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PACK_KEYS DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE MEMORY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE DISK
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE NULL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident UNION()
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TTL ident + INTERVAL FALSE MICROSECOND
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AFFINITY 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE AUTO_INCREMENT 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT 1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident COMMENT 'binary'
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary' COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary', COMMENT 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE ident
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DYNAMIC
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT FIXED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPRESSED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT REDUNDANT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPACT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_DEFAULT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_FAST
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SMALL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZLIB
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZSTD
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_QUICKLZ
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_LZMA
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SNAPPY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_UNCOMPRESSED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SMALLINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MIDDLEINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT1
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT2
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT3
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT4
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT8
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INTEGER
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIGINT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOLEAN
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NUMERIC
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FIXED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident REAL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE PRECISION
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT4
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT8
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHARACTER(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NVARCHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BINARY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARBINARY(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYBLOB
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary')
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SET('binary')
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident JSON
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHARACTER VARYING
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHAR VARYING
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHARACTER
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHARACTER(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHARACTER(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER VARYING(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER VARYING(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR VARYING(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR VARYING(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARYING(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident YEAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SQL_TSI_YEAR
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BLOB
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMBLOB
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGBLOB
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARBINARY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TEXT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMTEXT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGTEXT
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT ASCII
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT BYTE
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATETIME
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIMESTAMP
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIME
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNSIGNED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SIGNED
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ZEROFILL
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1, 1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary' BINARY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR < FLOAT >
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY CHARSET 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary' COLLATE 'binary'
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(0x1)
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(b'1')
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE t (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary', 'binary')
< error
| 1146 42S02 Table 'simulation.ident' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT TIFLASH REPLICA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT PARTITION ident
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident TIFLASH REPLICA
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES 'binary'
< error
| 8237 HY000 Invalid attributes: attributes should be in format 'key=value': binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS DEFAULT
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS 'binary'
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME PARTITIONING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SPLIT MAXVALUE PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported SPLIT LAST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MERGE FIRST PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported MERGE FIRST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident ATTRIBUTES DEFAULT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| TTL
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `TTL`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| TTL
: bytes
| CREATE TABLE `TTL` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET TIFLASH REPLICA 1
< error
| 1105 HY000 the tiflash replica count: 1 should be less than the total tiflash server count: 0
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET HYPO TIFLASH REPLICA 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET DEFAULT
< error
| 1105 HY000 Unknown charset 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD(ident TINYINT)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD INDEX(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION PARTITIONS 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LAST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ALL
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP ident
< error
| 1091 42000 Can't DROP 'ident'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PRIMARY KEY
< error
| 8200 HY000 Unsupported drop primary key when the table is using clustered index
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FIRST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP STATS_EXTENDED ident
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident OPTIMIZE PARTITION ALL
< error
| 8200 HY000 Unsupported optimize partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REPAIR PARTITION ALL
< error
| 8200 HY000 Unsupported repair partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REBUILD PARTITION ALL
< error
| 8200 HY000 Unsupported rebuild partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP KEY ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP FOREIGN KEY ident
< error
| 1091 42000 Can't DROP 'ident'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MODIFY ident TINYINT
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHANGE ident ident TINYINT
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT FALSE
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident DROP DEFAULT
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME COLUMN ident TO ident
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME AS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME KEY ident TO ident
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ WRITE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITH VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITHOUT VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_LOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_UNLOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CHECK ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident VISIBLE
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /* CACHED ON */
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION ident INTO(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ident
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITH VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY CLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY NONCLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP INDEX ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP COLUMN ident
< error
| 1091 42000 Can't DROP 'ident'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT FIRST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `ident` tinyint DEFAULT NULL,
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AFTER ident
< error
| 1054 42S22 Unknown column 'ident' in 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE, FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT zzz FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT ident FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident, ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SERIAL
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT PRIMARY KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident NOT ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident) ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NOT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint NOT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NOT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint NOT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_INCREMENT
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SERIAL DEFAULT VALUE
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'ident' column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL COMMENT 'binary',
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLLATE 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1, 1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DISK
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE MEMORY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT GENERATED ALWAYS AS(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) VIRTUAL
< error
| 1054 42S22 Unknown column 'ident' in 'generated column function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) STORED
< error
| 3106 HY000 'Adding generated stored column through ALTER TABLE' is not supported for generated columns.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PRIMARY KEY(ident)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD UNIQUE(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FOREIGN KEY(ident) REFERENCES ident
< error
| 1239 42000 Incorrect foreign key definition for 'fk_1': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH FULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH PARTIAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH SIMPLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT ON UPDATE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT ON DELETE RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET NULL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE NO ACTION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT FALSE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT '0',
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NEXT VALUE FOR ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT ident()
< error
| 3770 HY000 Default value expression of column 'ident' contains a disallowed function: `ident`.
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURDATE()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(NEXT VALUE FOR ident)
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1e0
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1.5
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident DEPENDENCY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CORRELATION(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident ON ident(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident, ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT((ident))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT ALGORITHM DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT LOCK DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX ident ON ident(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident LIKE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION SYSTEM_TIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY HASH(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY ALGORITHM = 1()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident MICROSECOND)
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(MICROSECOND)
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) NULL PARTITION
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) MAXVALUE PARTITION
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) FIRST PARTITION LESS THAN(ident) LAST PARTITION LESS THAN(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'partition function'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY LINEAR KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident, PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident, SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident COMMENT 'binary'
< error
| 1105 HY000 unknown partition option
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INSERT_METHOD 'binary'
< error
| 8233 HY000 CREATE/ALTER table with insert method option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DATA DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INDEX DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MAX_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MIN_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLESPACE ident
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NODEGROUP 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PLACEMENT POLICY 'binary'
< error
| 8239 HY000 Unknown placement policy 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident DEFAULT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident HISTORY)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident CURRENT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident IGNORE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident REPLACE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident AS
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT * UNION SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident WITH ident AS(SELECT *) SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS WITH ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(LIKE ident)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = 'binary' VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident, ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH CASCADED CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH LOCAL CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TEMPORARY TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW IF EXISTS ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) WITH PARSER ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) CLUSTERED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) LOCAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY zzz(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY USING BTREE(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY BTREE(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY ident TYPE BTREE(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident USING BTREE ON ident(ident)
< error
| 1072 42000 column does not exist: ident
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) TYPE BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HASH
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING RTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HYPO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HNSW
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident INVISIBLE
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident, ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident ASC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident DESC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(FULLTEXT(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(INDEX(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CONSTRAINT ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_ID_CACHE 1
< error
| 1105 HY000 Can't Alter AUTO_ID_CACHE between 1 and non-1, the underlying implementation is different
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_RANDOM_BASE 1
< error
| 8216 HY000 Invalid auto random: alter auto_random_base of a non auto_random table
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AVG_ROW_LENGTH 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONNECTION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLE_CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PASSWORD 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPRESSION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident KEY_BLOCK_SIZE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DELAY_KEY_WRITE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_BUCKETS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_TOPN 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_CHOICE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_LIST 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SHARD_ROW_ID_BITS 1
< error
| 8200 HY000 Unsupported shard_row_id_bits for table with primary key as row id
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PRE_SPLIT_REGIONS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PACK_KEYS DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE MEMORY
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE DISK
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE NULL
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident UNION()
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TTL ident + INTERVAL FALSE MICROSECOND
< error
| 1054 42S22 Unknown column 'ident' in 'TTL config'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AFFINITY 'binary'
< error
| 8266 HY000 Invalid AFFINITY 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident COMMENT 'binary'
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary' COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary', COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_FAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SMALL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZLIB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZSTD
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_QUICKLZ
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_LZMA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SNAPPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_UNCOMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` date DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` float DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` bit(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SMALLINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` smallint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MIDDLEINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT2
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` smallint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT3
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT8
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` bigint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INTEGER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIGINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` bigint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOLEAN
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NUMERIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` decimal(10,0) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident REAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` double DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` double DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE PRECISION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` double DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT4
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` float DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT8
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` double DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NVARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARBINARY(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varbinary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyblob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinytext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary')
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` enum('binary') DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SET('binary')
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` set('binary') DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident JSON
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` json DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` vector DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHARACTER VARYING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHAR VARYING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHARACTER
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARYING(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` varchar(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident YEAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` year(4) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SQL_TSI_YEAR
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` year(4) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` blob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumblob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGBLOB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` longblob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARBINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumblob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` text DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` mediumtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGTEXT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` longtext DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT ASCII
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinytext CHARACTER SET latin1 COLLATE latin1_bin DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT BYTE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyblob DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATETIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` datetime DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIMESTAMP
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` timestamp NULL DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` time DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNSIGNED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint unsigned DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SIGNED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ZEROFILL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` tinyint(3) unsigned zerofill DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` decimal(1,0) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1, 1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` decimal(1,1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary' BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR < FLOAT >
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` vector DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` char(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY CHARSET 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `ident` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary' COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(0x1)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
|q "CREATE TABLE `ident` (\n  `b` int NOT NULL,\n  `c` int DEFAULT NULL,\n  `d` int DEFAULT NULL,\n  `ident` enum('\x01') DEFAULT NULL,\n  PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,\n  KEY `idx_d` (`d`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(b'1')
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
|q "CREATE TABLE `ident` (\n  `b` int NOT NULL,\n  `c` int DEFAULT NULL,\n  `d` int DEFAULT NULL,\n  `ident` enum('\x01') DEFAULT NULL,\n  PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,\n  KEY `idx_d` (`d`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary', 'binary')
< error
| 1291 HY000 Column 'ident' has duplicated value 'binary' in ENUM
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT TIFLASH REPLICA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT PARTITION ident
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident TIFLASH REPLICA
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES 'binary'
< error
| 8237 HY000 Invalid attributes: attributes should be in format 'key=value': binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS DEFAULT
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS 'binary'
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME PARTITIONING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SPLIT MAXVALUE PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported SPLIT LAST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MERGE FIRST PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported MERGE FIRST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident ATTRIBUTES DEFAULT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| TTL
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `TTL`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| TTL
: bytes
| CREATE TABLE `TTL` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET TIFLASH REPLICA 1
< error
| 1105 HY000 the tiflash replica count: 1 should be less than the total tiflash server count: 0
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SET HYPO TIFLASH REPLICA 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET DEFAULT
< error
| 1105 HY000 Unknown charset 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD(ident TINYINT)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD INDEX(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION PARTITIONS 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LAST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ALL
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PRIMARY KEY
< error
| 8200 HY000 Unsupported drop primary key when the table is using clustered index
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FIRST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP STATS_EXTENDED ident
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident OPTIMIZE PARTITION ALL
< error
| 8200 HY000 Unsupported optimize partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REPAIR PARTITION ALL
< error
| 8200 HY000 Unsupported repair partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REBUILD PARTITION ALL
< error
| 8200 HY000 Unsupported rebuild partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP KEY ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP FOREIGN KEY ident
< error
| 1091 42000 Can't DROP 'ident'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DISABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MODIFY ident TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` tinyint DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHANGE ident ident TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` tinyint DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT FALSE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT '0',
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'expression'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME COLUMN ident TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME AS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME KEY ident TO ident
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ WRITE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITH VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident WITHOUT VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_LOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_UNLOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CHECK ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident VISIBLE
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /* CACHED ON */
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION ident INTO(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ident
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITH VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY CLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY NONCLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP INDEX ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP COLUMN ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT FIRST
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AFTER ident
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE, FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT zzz FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT ident FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident, ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SERIAL
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT PRIMARY KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident NOT ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident) ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NOT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NOT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_INCREMENT
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SERIAL DEFAULT VALUE
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COMMENT 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLLATE 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1, 1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DISK
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE MEMORY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT FIXED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT GENERATED ALWAYS AS(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) VIRTUAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) STORED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PRIMARY KEY(ident)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD UNIQUE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   UNIQUE KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FOREIGN KEY(ident) REFERENCES ident
< error
| 1239 42000 Incorrect foreign key definition for 'fk_1': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH FULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH PARTIAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH SIMPLE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE CASCADE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE NO ACTION
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT FALSE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NEXT VALUE FOR ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT ident()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURDATE()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(NEXT VALUE FOR ident)
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1e0
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1.5
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident DEPENDENCY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CORRELATION(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident ON ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident, ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT((ident))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT ALGORITHM DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT LOCK DEFAULT
< error
| 1091 42000 index ident doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX ident ON ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   UNIQUE KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident LIKE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION SYSTEM_TIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY HASH(ident)
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY ALGORITHM = 1()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident MICROSECOND)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(MICROSECOND)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) NULL PARTITION
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) MAXVALUE PARTITION
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) FIRST PARTITION LESS THAN(ident) LAST PARTITION LESS THAN(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'expression'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY LINEAR KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident, PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident, SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident COMMENT 'binary'
< error
| 1105 HY000 unknown partition option
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INSERT_METHOD 'binary'
< error
| 8233 HY000 CREATE/ALTER table with insert method option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DATA DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident INDEX DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MAX_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident MIN_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLESPACE ident
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident NODEGROUP 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PLACEMENT POLICY 'binary'
< error
| 8239 HY000 Unknown placement policy 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident DEFAULT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident HISTORY)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident CURRENT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident IGNORE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident REPLACE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident AS
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT * UNION SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident WITH ident AS(SELECT *) SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS WITH ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(LIKE ident)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = 'binary' VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident, ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH CASCADED CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH LOCAL CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TEMPORARY TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW IF EXISTS ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP TABLE ident CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) WITH PARSER ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) CLUSTERED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) LOCAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `zzz` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY USING BTREE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY BTREE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `BTREE` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY ident TYPE BTREE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX ident USING BTREE ON ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`),
|   KEY `ident` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) TYPE BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HASH
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING RTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HYPO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HNSW
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident INVISIBLE
< error
| 1176 42000 Key 'ident' doesn't exist in table 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident, ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident ASC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident DESC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(FULLTEXT(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(INDEX(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CONSTRAINT ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_ID_CACHE 1
< error
| 1105 HY000 Can't Alter AUTO_ID_CACHE between 1 and non-1, the underlying implementation is different
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_RANDOM_BASE 1
< error
| 8216 HY000 Invalid auto random: alter auto_random_base of a non auto_random table
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AVG_ROW_LENGTH 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONNECTION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLE_CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PASSWORD 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPRESSION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident KEY_BLOCK_SIZE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident DELAY_KEY_WRITE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_BUCKETS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_TOPN 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_CHOICE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_LIST 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SHARD_ROW_ID_BITS 1
< error
| 8200 HY000 Unsupported shard_row_id_bits for table with primary key as row id
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PRE_SPLIT_REGIONS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident PACK_KEYS DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE MEMORY
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE DISK
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE NULL
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident UNION()
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident TTL ident + INTERVAL FALSE MICROSECOND
< error
| 8148 HY000 Field 'ident' is of a not supported type for TTL config, expect DATETIME, DATE or TIMESTAMP
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident AFFINITY 'binary'
< error
| 8266 HY000 Invalid AFFINITY 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE ident COMMENT 'binary'
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary' COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary', COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| TRUNCATE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_FAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SMALL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZLIB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZSTD
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_QUICKLZ
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_LZMA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SNAPPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_UNCOMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SMALLINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MIDDLEINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT1
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT2
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT3
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT4
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT8
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INTEGER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIGINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOLEAN
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NUMERIC
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FIXED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident REAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE PRECISION
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT4
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT8
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NVARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARBINARY(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SET('binary')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident JSON
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHARACTER VARYING
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHAR VARYING
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident YEAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SQL_TSI_YEAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARBINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT ASCII
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT BYTE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATETIME
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIME
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNSIGNED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SIGNED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ZEROFILL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1, 1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary' BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR < FLOAT >
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY CHARSET 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary' COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(0x1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(b'1')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary', 'binary')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ANALYZE PARTITION ident INDEX zzz
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT TIFLASH REPLICA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE COMPACT PARTITION ident
< error
| 1146 42S02 Table 'simulation.COMPACT' doesn't exist
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPACT PARTITION ident TIFLASH REPLICA
< error
| 1105 HY000 table:ident is not a partition table, but user specify partition name list:[ident]
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ATTRIBUTES 'binary'
< error
| 8237 HY000 Invalid attributes: attributes should be in format 'key=value': binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS DEFAULT
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_OPTIONS 'binary'
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME PARTITIONING
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `PARTITIONING`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| PARTITIONING
: bytes
| CREATE TABLE `PARTITIONING` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SPLIT MAXVALUE PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported SPLIT LAST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident MERGE FIRST PARTITION LESS THAN(ident)
< error
| 8200 HY000 Unsupported MERGE FIRST PARTITION
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident ATTRIBUTES DEFAULT
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REMOVE TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TTL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| TTL
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `TTL`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| TTL
: bytes
| CREATE TABLE `TTL` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SET TIFLASH REPLICA 1
< error
| 1105 HY000 the tiflash replica count: 1 should be less than the total tiflash server count: 0
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SET HYPO TIFLASH REPLICA 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET DEFAULT
< error
| 1105 HY000 Unknown charset 
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD(ident TINYINT)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD INDEX(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   KEY `ident_2` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION PARTITIONS 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident LAST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CARDINALITY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ALL
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PRIMARY KEY
< error
| 8200 HY000 Unsupported drop primary key when the table is using clustered index
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP PARTITION ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident FIRST PARTITION LESS THAN(ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP STATS_EXTENDED ident
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident OPTIMIZE PARTITION ALL
< error
| 8200 HY000 Unsupported optimize partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REPAIR PARTITION ALL
< error
| 8200 HY000 Unsupported repair partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD PARTITION ALL TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident IMPORT TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DISCARD TABLESPACE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REBUILD PARTITION ALL
< error
| 8200 HY000 Unsupported rebuild partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP KEY ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP FOREIGN KEY ident
< error
| 1091 42000 Can't DROP 'ident'; check that column/key exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DISABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ENABLE KEYS
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident MODIFY ident TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` tinyint DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CHANGE ident ident TINYINT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` tinyint DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT FALSE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT '0',
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident SET DEFAULT(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'expression'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER ident DROP DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME COLUMN ident TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME AS ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident RENAME KEY ident TO ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident READ WRITE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident WITH VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident WITHOUT VALIDATION
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_LOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_UNLOAD
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CHECK ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin /* CACHED ON */
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident NOCACHE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident REORGANIZE PARTITION ident INTO(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECK PARTITION ident
< error
| 8200 HY000 Unsupported check partition
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITH VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident EXCHANGE PARTITION ident WITH TABLE ident WITHOUT VALIDATION
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY CLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY NONCLUSTERED
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT KEY GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM COPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INPLACE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALGORITHM INSTANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident READ ONLY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP INDEX ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT KEY zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP COLUMN ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT FIRST
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AFTER ident
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE, FORCE
< error
| 8200 HY000 Unsupported Unsupported/unknown ALTER TABLE specification
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT zzz FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CONSTRAINT ident FULLTEXT(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| RENAME TABLE ident TO ident, ident TO ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SERIAL
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT PRIMARY KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint PRIMARY KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER CHECK ident NOT ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD CHECK(ident) ENFORCED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NOT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident) NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT CHECK(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NOT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_INCREMENT
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE GLOBAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE LOCAL
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNIQUE KEY
< error
| 8200 HY000 unsupported add column 'ident' constraint UNIQUE KEY when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SERIAL DEFAULT VALUE
< error
| 8200 HY000 unsupported add column 'ident' constraint AUTO_INCREMENT when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ON UPDATE CURRENT_TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COMMENT 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLLATE 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AUTO_RANDOM(1, 1)
< error
| 8216 HY000 Invalid auto random: unsupported add column 'ident' constraint AUTO_RANDOM when altering 'simulation.ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE DISK
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT STORAGE MEMORY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT FIXED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT COLUMN_FORMAT DYNAMIC
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT GENERATED ALWAYS AS(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) VIRTUAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT AS(ident) STORED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT NULL NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PRIMARY KEY(ident)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   KEY `ident_2` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD UNIQUE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   UNIQUE KEY `ident_2` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FOREIGN KEY(ident) REFERENCES ident
< error
| 1239 42000 Incorrect foreign key definition for 'fk_1': Key reference and table reference don't match
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH FULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH PARTIAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident MATCH SIMPLE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE RESTRICT ON UPDATE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON UPDATE RESTRICT ON DELETE RESTRICT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE CASCADE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET NULL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE NO ACTION
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident ON DELETE SET DEFAULT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT FALSE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NEXT VALUE FOR ident
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT ident()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(CURRENT_TIMESTAMP)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_TIMESTAMP(1)
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURDATE()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT(NEXT VALUE FOR ident)
< error
| 1347 HY000 'simulation.ident' is not SEQUENCE
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIME
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT LOCALTIMESTAMP
< error
| 1067 42000 Invalid default value for 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT DEFAULT CURRENT_DATE()
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1e0
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_RATE 1.5
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident DEPENDENCY(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD STATS_EXTENDED ident CORRELATION(ident)
< error
| 1105 HY000 Extended statistics feature is not generally available now, and tidb_enable_extended_stats is OFF
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE INDEX ident ON ident(ident)
< error
| 1061 42000 Duplicate key name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT REFERENCES ident(ident)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident, ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT((ident))
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident LOCK DEFAULT ALGORITHM DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP INDEX ident ON ident ALGORITHM DEFAULT LOCK DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX ident ON ident(ident)
< error
| 1061 42000 Duplicate key name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE SPATIAL INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ident ON ident(ident)
< error
| 8200 HY000 FULLTEXT and SPATIAL index is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident LIKE ident
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION SYSTEM_TIME
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY HASH(ident)
< error
| 1503 HY000 A CLUSTERED INDEX must include all columns in the table's partitioning function
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY KEY ALGORITHM = 1()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident MICROSECOND)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(MICROSECOND)
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) NULL PARTITION
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) MAXVALUE PARTITION
< error
| 8200 HY000 Unsupported INTERVAL partitioning, currently requires FIRST and LAST partitions to be defined
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY RANGE(ident) INTERVAL(ident) FIRST PARTITION LESS THAN(ident) LAST PARTITION LESS THAN(ident)
< error
| 1054 42S22 Unknown column 'ident' in 'expression'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION BY LINEAR KEY()
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
| PARTITION BY KEY (`b`) PARTITIONS 1
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident, PARTITION ident)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident(SUBPARTITION ident, SUBPARTITION ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PARTITION ident COMMENT 'binary'
< error
| 1105 HY000 unknown partition option
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE ENGINE 'binary'
< error
| 1286 42000 Unknown storage engine 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident INSERT_METHOD 'binary'
< error
| 8233 HY000 CREATE/ALTER table with insert method option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DATA DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident INDEX DIRECTORY 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident MAX_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident MIN_ROWS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLESPACE ident
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident NODEGROUP 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PLACEMENT POLICY 'binary'
< error
| 8239 HY000 Unknown placement policy 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN MAXVALUE)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident DEFAULT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident HISTORY)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident CURRENT)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident IGNORE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident REPLACE
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident AS
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT * UNION SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident WITH ident AS(SELECT *) SELECT *
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * UNION SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS WITH ident AS(SELECT *) SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS(SELECT *)
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(LIKE ident)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE OR REPLACE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE ALGORITHM = UNDEFINED VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE ALGORITHM = MERGE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE ALGORITHM = TEMPTABLE VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE DEFINER = 'binary' VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE SQL SECURITY DEFINER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE SQL SECURITY INVOKER VIEW ident AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident(ident, ident) AS SELECT *
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH CASCADED CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE VIEW ident AS SELECT * WITH LOCAL CHECK OPTION
< error
| 1051 42S02 Unknown table ''
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TEMPORARY TABLE ident
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP VIEW ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP VIEW IF EXISTS ident
< error
| 1347 HY000 'simulation.ident' is not VIEW
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP TABLE ident RESTRICT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| DROP TABLE ident CASCADE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(ident))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAXVALUE, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES LESS THAN(MAX_CONNECTIONS_PER_HOUR, MAXVALUE))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFAULT, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD PARTITION(PARTITION ident VALUES IN(DEFINED, DEFAULT))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER IGNORE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT ident(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) VISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) KEY_BLOCK_SIZE 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) WITH PARSER ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) CLUSTERED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) GLOBAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) LOCAL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY zzz(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   KEY `zzz` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY USING BTREE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   KEY `ident_2` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY BTREE(ident)
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`),
|   KEY `BTREE` (`ident`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD KEY ident TYPE BTREE(ident)
< error
| 1061 42000 Duplicate key name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE INDEX ident USING BTREE ON ident(ident)
< error
| 1061 42000 Duplicate key name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) TYPE BTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HASH
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING RTREE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HYPO
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD FULLTEXT(ident) USING HNSW
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ALTER INDEX ident INVISIBLE
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident, ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident ASC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ORDER BY ident DESC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| TRUNCATE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(FULLTEXT(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(INDEX(ident))
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DROP CONSTRAINT ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYINT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident(ident TINYTEXT, ident TINYINT)
< error
| 1050 42S01 Table 'simulation.ident' already exists
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CHARSET 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_ID_CACHE 1
< error
| 1105 HY000 Can't Alter AUTO_ID_CACHE between 1 and non-1, the underlying implementation is different
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident AUTO_RANDOM_BASE 1
< error
| 8216 HY000 Invalid auto random: alter auto_random_base of a non auto_random table
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident AVG_ROW_LENGTH 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CONNECTION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident TABLE_CHECKSUM 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PASSWORD 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMPRESSION 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident KEY_BLOCK_SIZE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident DELAY_KEY_WRITE 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_AUTO_RECALC DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_SAMPLE_PAGES DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_BUCKETS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_TOPN 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_CHOICE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_COL_LIST 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SHARD_ROW_ID_BITS 1
< error
| 8200 HY000 Unsupported shard_row_id_bits for table with primary key as row id
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PRE_SPLIT_REGIONS 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident PACK_KEYS DEFAULT
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE MEMORY
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STORAGE DISK
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE NULL
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident SECONDARY_ENGINE 'binary'
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident UNION()
< error
| 8232 HY000 CREATE/ALTER table with union option is not supported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident TTL ident + INTERVAL FALSE MICROSECOND
< error
| 8148 HY000 Field 'ident' is of a not supported type for TTL config, expect DATETIME, DATE or TIMESTAMP
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident AFFINITY 'binary'
< error
| 8266 HY000 Invalid AFFINITY 'binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident FORCE AUTO_INCREMENT 1
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident STATS_PERSISTENT 1
< error
| 8200 HY000 This type of ALTER TABLE is currently unsupported
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| CREATE TABLE ident COMMENT 'binary'
< error
| 1113 42000 A table must have at least 1 column
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary' COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident COMMENT 'binary', COMMENT 'binary'
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='binary'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| TRUNCATE TABLE ident
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT DYNAMIC
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT FIXED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT REDUNDANT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT COMPACT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_DEFAULT
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_FAST
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SMALL
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZLIB
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_ZSTD
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_QUICKLZ
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_LZMA
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_SNAPPY
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ROW_FORMAT TOKUDB_UNCOMPRESSED
< ok 0 0
> query
| SHOW FULL TABLES
< result
< column "Tables_in_simulation" "VARCHAR"
< column "Table_type" "VARCHAR"
< row
: bytes
| ident
: bytes
| BASE TABLE
> query
| SHOW CREATE TABLE `ident`
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ident
: bytes
| CREATE TABLE `ident` (
|   `b` int NOT NULL,
|   `ident` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   KEY `ident` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SMALLINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MIDDLEINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT1
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT2
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT3
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT4
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INT8
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident INTEGER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BIGINT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BOOLEAN
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NUMERIC
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FIXED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident REAL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DOUBLE PRECISION
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT4
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident FLOAT8
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NVARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VARBINARY(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SET('binary')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident JSON
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHARACTER VARYING
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG CHAR VARYING
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARCHARACTER
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARCHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHARACTER(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NATIONAL CHAR(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident NCHAR VARYING(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident YEAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident SQL_TSI_YEAR
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident BLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGBLOB
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONG VARBINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident MEDIUMTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident LONGTEXT
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT ASCII
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT UNICODE
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYTEXT BYTE
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DATETIME
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIMESTAMP
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TIME
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT UNSIGNED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT SIGNED
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident TINYINT ZEROFILL
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident DECIMAL(1, 1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER CHARSET 'binary' BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident VECTOR < FLOAT >
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident CHARACTER BINARY CHARSET 'binary'
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident CONVERT TO CHARSET 'binary' COLLATE 'binary'
< error
| 8200 HY000 Unsupported modify charset from utf8mb4 to binary
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(0x1)
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM(b'1')
< error
| 1060 42S21 Duplicate column name 'ident'
> exec
| DROP DATABASE IF EXISTS simulation; CREATE DATABASE simulation; USE simulation; CREATE TABLE ident (b INT PRIMARY KEY, ident INT, d INT, KEY ident (d))
< ok 0 0
> exec
| ALTER TABLE ident ADD ident ENUM('binary', 'binary')
< error
| 1060 42S21 Duplicate column name 'ident'
