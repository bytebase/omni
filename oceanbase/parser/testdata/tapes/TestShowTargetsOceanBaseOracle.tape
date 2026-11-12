-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestShowTargetsOceanBaseOracle default#1
> exec
| DROP DATABASE IF EXISTS obshowt
< ok 0 0
> exec
| CREATE DATABASE obshowt
< ok 1 0
> exec
| USE obshowt
< ok 0 0
> exec
| SHOW
< error
| 1064 42000 …near 'SHOW' at line 1
> exec
| SHOW FULL
< error
| 1064 42000 …near 'FULL' at line 1
> exec
| SHOW EXTENDED
< error
| 1064 42000 …near 'EXTENDED' at line 1
> exec
| SHOW CREATE
< error
| 1064 42000 …near 'CREATE' at line 1
> exec
| SHOW TABLE
< error
| 1064 42000 …near 'TABLE' at line 1
> exec
| SHOW GLOBAL
< error
| 1064 42000 …near 'GLOBAL' at line 1
> exec
| SHOW SESSION
< error
| 1064 42000 …near 'SESSION' at line 1
> exec
| SHOW MASTER
< error
| 1064 42000 …near 'MASTER' at line 1
> exec
| SHOW SLAVE
< error
| 1064 42000 …near 'SLAVE' at line 1
> exec
| SHOW REPLICA
< error
| 1064 42000 …near 'REPLICA' at line 1
> exec
| SHOW BINARY
< error
| 1064 42000 …near 'BINARY' at line 1
> exec
| SHOW BINLOG
< error
| 1064 42000 …near 'BINLOG' at line 1
> exec
| SHOW RELAYLOG
< error
| 1064 42000 …near 'RELAYLOG' at line 1
> exec
| SHOW STORAGE
< error
| 1064 42000 …near 'STORAGE' at line 1
> exec
| SHOW OPEN
< error
| 1064 42000 …near 'OPEN' at line 1
> exec
| SHOW CHARACTER
< error
| 1064 42000 …near 'CHARACTER' at line 1
> exec
| SHOW ENGINE
< error
| 1064 42000 …near 'ENGINE' at line 1
> exec
| SHOW PROCEDURE
< error
| 1064 42000 …near 'PROCEDURE' at line 1
> exec
| SHOW FUNCTION
< error
| 1064 42000 …near 'FUNCTION' at line 1
> exec
| SHOW COUNT
< error
| 1064 42000 …near 'COUNT' at line 1
> exec
| SHOW COUNT(*)
< error
| 1064 42000 …near ')' at line 1
> exec
| SHOW COUNT ERRORS
< error
| 1064 42000 …near 'ERRORS' at line 1
> exec
| SHOW COUNT(*) ERRORS
< ok 0 0
> exec
| SHOW COUNT(*) WARNINGS
< ok 0 0
> exec
| SHOW COUNT(*) TABLES
< error
| 1064 42000 …near 'TABLES' at line 1
> exec
| SHOW DATABASES
< ok 0 0
> exec
| SHOW SCHEMAS
< ok 0 0
> exec
| SHOW TABLES
< ok 0 0
> exec
| SHOW TABLES FROM test
< ok 0 0
> exec
| SHOW TABLES IN test
< ok 0 0
> exec
| SHOW FULL TABLES
< ok 0 0
> exec
| SHOW FULL COLUMNS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW FULL FIELDS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW FULL PROCESSLIST
< ok 0 0
> exec
| SHOW FULL STATUS
< error
| 1064 42000 …near 'STATUS' at line 1
> exec
| SHOW COLUMNS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW COLUMNS FROM t FROM test
< error
| 1146 42S02 Table doesn't exist
> exec
| SHOW COLUMNS
< error
| 1064 42000 …near 'COLUMNS' at line 1
> exec
| SHOW COLUMNS t
< error
| 1064 42000 …near 't' at line 1
> exec
| SHOW FIELDS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED COLUMNS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED FIELDS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED FULL COLUMNS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED FULL FIELDS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED TABLES
< ok 0 0
> exec
| SHOW EXTENDED FULL TABLES
< ok 0 0
> exec
| SHOW EXTENDED INDEX FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW EXTENDED KEYS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW CREATE TABLE t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW CREATE DATABASE d
< error
| 1049 42000 Unknown database 'd'
> exec
| SHOW CREATE DATABASE IF NOT EXISTS d
< error
| 1049 42000 Unknown database 'd'
> exec
| SHOW CREATE DATABASE IF d
< error
| 1064 42000 …near 'd' at line 1
> exec
| SHOW CREATE SCHEMA d
< error
| 1049 42000 Unknown database 'd'
> exec
| SHOW CREATE VIEW v
< error
| 1146 42S02 Table 'obshowt.v' doesn't exist
> exec
| SHOW CREATE PROCEDURE pr
< error
| 1305 42000 procedure/function does not exist
> exec
| SHOW CREATE FUNCTION fn
< error
| 1305 42000 procedure/function does not exist
> exec
| SHOW CREATE TRIGGER tg
< error
| 1360 HY000 Trigger does not exist
> exec
| SHOW CREATE EVENT ev
< error
| 1064 42000 …near 'EVENT ev' at line 1
> exec
| SHOW CREATE USER u
< error
| 1396 HY000 Operation SHOW CREATE USER failed for `u`@`%`
> exec
| SHOW CREATE USER
< error
| 1064 42000 …near 'USER' at line 1
> exec
| SHOW CREATE SEQUENCE seq
< error
| 1064 42000 …near 'SEQUENCE seq' at line 1
> exec
| SHOW CREATE SEQUENCE test.seq
< error
| 1064 42000 …near 'SEQUENCE test.seq' at line 1
> exec
| SHOW CREATE PLACEMENT POLICY x
< error
| 1064 42000 …near 'PLACEMENT POLICY x' at line 1
> exec
| SHOW CREATE PLACEMENT POLICY
< error
| 1064 42000 …near 'PLACEMENT POLICY' at line 1
> exec
| SHOW CREATE RESOURCE GROUP rg
< error
| 1064 42000 …near 'RESOURCE GROUP rg' at line 1
> exec
| SHOW CREATE RESOURCE GROUP default
< error
| 1064 42000 …near 'RESOURCE GROUP default' at line 1
> exec
| SHOW CREATE RESOURCE GROUP
< error
| 1064 42000 …near 'RESOURCE GROUP' at line 1
> exec
| SHOW CREATE RESOURCE rg
< error
| 1064 42000 …near 'RESOURCE rg' at line 1
> exec
| SHOW INDEX FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW INDEX IN t IN test
< error
| 1146 42S02 Table doesn't exist
> exec
| SHOW INDEX
< error
| 1064 42000 …near 'INDEX' at line 1
> exec
| SHOW INDEX t
< error
| 1064 42000 …near 't' at line 1
> exec
| SHOW KEY FROM t
< error
| 1064 42000 …near 'KEY FROM t' at line 1
> exec
| SHOW KEYS FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW INDEXES FROM t
< error
| 1146 42S02 Table 'obshowt.t' doesn't exist
> exec
| SHOW GLOBAL VARIABLES
< ok 0 0
> exec
| SHOW GLOBAL STATUS
< ok 0 0
> exec
| SHOW SESSION VARIABLES
< ok 0 0
> exec
| SHOW SESSION STATUS
< ok 0 0
> exec
| SHOW GLOBAL BINDINGS
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec
| SHOW SESSION BINDINGS
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec
| SHOW BINDINGS
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec
| SHOW BINDING_CACHE STATUS
< error
| 1064 42000 …near 'BINDING_CACHE STATUS' at line 1
> exec
| SHOW BINDING_CACHE
< error
| 1064 42000 …near 'BINDING_CACHE' at line 1
> exec
| SHOW VARIABLES
< ok 0 0
> exec
| SHOW VARIABLES LIKE 'x%'
< ok 0 0
> exec
| SHOW STATUS
< ok 0 0
> exec
| SHOW WARNINGS
< ok 0 0
> exec
| SHOW WARNINGS LIMIT 10
< ok 0 0
> exec
| SHOW WARNINGS LIMIT 5, 10
< ok 0 0
> exec
| SHOW ERRORS
< ok 0 0
> exec
| SHOW ERRORS LIMIT 10
< ok 0 0
> exec
| SHOW ENGINES
< ok 0 0
> exec
| SHOW STORAGE ENGINES
< ok 0 0
> exec
| SHOW ENGINE INNODB STATUS
< ok 0 0
> exec
| SHOW ENGINE INNODB MUTEX
< ok 0 0
> exec
| SHOW ENGINE INNODB LOGS
< error
| 1064 42000 …near 'LOGS' at line 1
> exec
| SHOW ENGINE INNODB
< error
| 1064 42000 …near 'INNODB' at line 1
> exec
| SHOW PLUGINS
< ok 0 0
> exec
| SHOW MASTER STATUS
< error
| 1064 42000 …near 'MASTER STATUS' at line 1
> exec
| SHOW MASTER LOGS
< error
| 1064 42000 …near 'MASTER LOGS' at line 1
> exec
| SHOW BINARY LOGS
< error
| 1064 42000 …near 'BINARY LOGS' at line 1
> exec
| SHOW BINARY LOG STATUS
< error
| 1064 42000 …near 'BINARY LOG STATUS' at line 1
> exec
| SHOW BINARY LOG
< error
| 1064 42000 …near 'BINARY LOG' at line 1
> exec
| SHOW BINLOG STATUS
< error
| 1064 42000 …near 'BINLOG STATUS' at line 1
> exec
| SHOW SLAVE STATUS
< error
| 1064 42000 …near 'SLAVE STATUS' at line 1
> exec
| SHOW SLAVE HOSTS
< error
| 1064 42000 …near 'SLAVE HOSTS' at line 1
> exec
| SHOW REPLICA STATUS
< error
| 1064 42000 …near 'REPLICA STATUS' at line 1
> exec
| SHOW REPLICA HOSTS
< error
| 1064 42000 …near 'REPLICA HOSTS' at line 1
> exec
| SHOW REPLICAS
< error
| 1064 42000 …near 'REPLICAS' at line 1
> exec
| SHOW SLAVE STATUS FOR CHANNEL 'c'
< error
| 1064 42000 …near 'SLAVE STATUS FOR CHANNEL 'c'' at line 1
> exec
| SHOW SLAVE STATUS FOR CHANNEL c
< error
| 1064 42000 …near 'SLAVE STATUS FOR CHANNEL c' at line 1
> exec
| SHOW SLAVE STATUS FOR
< error
| 1064 42000 …near 'SLAVE STATUS FOR' at line 1
> exec
| SHOW SLAVE STATUS FOR CHANNEL
< error
| 1064 42000 …near 'SLAVE STATUS FOR CHANNEL' at line 1
> exec
| SHOW REPLICA STATUS FOR CHANNEL 'c'
< error
| 1064 42000 …near 'REPLICA STATUS FOR CHANNEL 'c'' at line 1
> exec
| SHOW REPLICA STATUS FOR CHANNEL c
< error
| 1064 42000 …near 'REPLICA STATUS FOR CHANNEL c' at line 1
> exec
| SHOW BINLOG EVENT
< error
| 1064 42000 …near 'BINLOG EVENT' at line 1
> exec
| SHOW BINLOG EVENTS
< error
| 1064 42000 …near 'BINLOG EVENTS' at line 1
> exec
| SHOW BINLOG EVENTS IN 'f'
< error
| 1064 42000 …near 'BINLOG EVENTS IN 'f'' at line 1
> exec
| SHOW BINLOG EVENTS FROM 4
< error
| 1064 42000 …near 'BINLOG EVENTS FROM 4' at line 1
> exec
| SHOW BINLOG EVENTS LIMIT 5, 10
< error
| 1064 42000 …near 'BINLOG EVENTS LIMIT 5, 10' at line 1
> exec
| SHOW RELAYLOG EVENT
< error
| 1064 42000 …near 'RELAYLOG EVENT' at line 1
> exec
| SHOW RELAYLOG EVENTS
< error
| 1064 42000 …near 'RELAYLOG EVENTS' at line 1
> exec
| SHOW RELAYLOG EVENTS FOR CHANNEL 'c'
< error
| 1064 42000 …near 'RELAYLOG EVENTS FOR CHANNEL 'c'' at line 1
> exec
| SHOW RELAYLOG EVENTS FOR CHANNEL c
< error
| 1064 42000 …near 'RELAYLOG EVENTS FOR CHANNEL c' at line 1
> exec
| SHOW TABLE STATUS
< ok 0 0
> exec
| SHOW TABLE STATUS FROM test
< ok 0 0
> exec
| SHOW TABLE STATUS LIKE 't'
< ok 0 0
> exec
| SHOW TRIGGER
< error
| 1064 42000 …near 'TRIGGER' at line 1
> exec
| SHOW TRIGGERS
< ok 0 0
> exec
| SHOW EVENT
< error
| 1064 42000 …near 'EVENT' at line 1
> exec
| SHOW EVENTS
< error
| 1064 42000 …near 'EVENTS' at line 1
> exec
| SHOW PROCEDURE STATUS
< ok 0 0
> exec
| SHOW PROCEDURE CODE pr
< error
| 1305 42000 procedure/function does not exist
> exec
| SHOW FUNCTION STATUS
< ok 0 0
> exec
| SHOW FUNCTION CODE fn
< error
| 1305 42000 procedure/function does not exist
> exec
| SHOW OPEN TABLES
< ok 0 0
> exec
| SHOW PRIVILEGES
< ok 0 0
> exec
| SHOW PROFILES
< ok 0 0
> exec
| SHOW PROFILE
< ok 0 0
> exec
| SHOW PROFILE CPU
< error
| 1064 42000 …near 'CPU' at line 1
> exec
| SHOW PROFILE CPU, MEMORY FOR QUERY 2
< error
| 1064 42000 …near 'CPU, MEMORY FOR QUERY 2' at line 1
> exec
| SHOW PROFILE BLOCK IO
< error
| 1064 42000 …near 'BLOCK IO' at line 1
> exec
| SHOW PROFILE FOR QUERY 2
< error
| 1064 42000 …near 'FOR QUERY 2' at line 1
> exec
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1
< error
| 1064 42000 …near 'CPU FOR QUERY 2 LIMIT 1' at line 1
> exec
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,1
< error
| 1064 42000 …near 'CPU FOR QUERY 2 LIMIT 1,1' at line 1
> exec
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1 OFFSET 2
< error
| 1064 42000 …near 'CPU FOR QUERY 2 LIMIT 1 OFFSET 2' at line 1
> exec
| SHOW PROFILE CPU LIMIT 1,1 OFFSET 2
< error
| 1064 42000 …near 'CPU LIMIT 1,1 OFFSET 2' at line 1
> exec
| SHOW PROFILE CPU LIMIT 1,
< error
| 1064 42000 …near 'CPU LIMIT 1,' at line 1
> exec
| SHOW PROFILE CPU OFFSET 2
< error
| 1064 42000 …near 'CPU OFFSET 2' at line 1
> exec
| SHOW CHARACTER SET
< ok 0 0
> exec
| SHOW CHARSET
< ok 0 0
> exec
| SHOW CHAR SET
< ok 0 0
> exec
| SHOW CHAR
< error
| 1064 42000 …near 'CHAR' at line 1
> exec
| SHOW COLLATION
< ok 0 0
> exec
| SHOW GRANTS
< ok 0 0
> exec
| SHOW GRANTS FOR u
< error
| 1133 42000 Can not find any matching row in the user table
> exec
| SHOW PROCESSLIST
< ok 0 0
> exec
| SHOW DATABASES LIKE 't'
< ok 0 0
> exec
| SHOW DATABASES LIKE @v
< error
| 1064 42000 …near '@v' at line 1
> exec
| SHOW DATABASES LIKE concat('a','%')
< error
| 1064 42000 …near 'concat('a','%')' at line 1
> exec
| SHOW DATABASES LIKE 1
< error
| 1064 42000 …near '1' at line 1
> exec
| SHOW DATABASES LIKE -1
< error
| 1064 42000 …near '-1' at line 1
> exec
| SHOW DATABASES LIKE NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec
| SHOW DATABASES LIKE ('t')
< error
| 1064 42000 …near '('t')' at line 1
> exec
| SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin' at line 1
> exec
| SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'
< error
| 1064 42000 …near 'COLLATE 'utf8mb4_bin'' at line 1
> exec
| SHOW DATABASES LIKE 't%' COLLATE binary
< error
| 1064 42000 …near 'COLLATE binary' at line 1
> exec
| SHOW DATABASES LIKE 't%' COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec
| SHOW DATABASES LIKE 't%' COLLATE 1
< error
| 1064 42000 …near 'COLLATE 1' at line 1
> exec
| SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin' at line 1
> exec
| SHOW DATABASES LIKE 'a' || 'b'
< error
| 1064 42000 …near '|| 'b'' at line 1
> exec
| SHOW DATABASES LIKE 'a' + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec
| SHOW DATABASES LIKE 'a' -> '$.x'
< error
| 1064 42000 …near '-> '$.x'' at line 1
> exec
| SHOW DATABASES LIKE
< error
| 1064 42000 …near 'LIKE' at line 1
> exec
| SHOW DATABASES LIKE 't' LIKE 'u'
< error
| 1064 42000 …near 'LIKE 'u'' at line 1
> exec
| SHOW DATABASES LIKE 't' AND 1
< error
| 1064 42000 …near 'AND 1' at line 1
> exec
| SHOW DATABASES LIKE 't' WHERE 1=1
< error
| 1064 42000 …near 'WHERE 1=1' at line 1
> exec
| SHOW DATABASES WHERE 1=1
< ok 0 0
> exec
| SHOW DATABASES WHERE
< error
| 1064 42000 …near 'WHERE' at line 1
> exec
| SHOW DATABASES WHERE 1=1 LIKE 't'
< ok 0 0
> exec
| SHOW TABLES LIKE 't' WHERE 1=1
< error
| 1064 42000 …near 'WHERE 1=1' at line 1
> exec
| SHOW PLACEMENT
< error
| 1064 42000 …near 'PLACEMENT' at line 1
> exec
| SHOW PLACEMENT LIKE 'POLICY foo%'
< error
| 1064 42000 …near 'PLACEMENT LIKE 'POLICY foo%'' at line 1
> exec
| SHOW PLACEMENT WHERE Target='TABLE test.t1'
< error
| 1064 42000 …near 'PLACEMENT WHERE Target='TABLE test.t1'' at line 1
> exec
| SHOW PLACEMENT LABELS
< error
| 1064 42000 …near 'PLACEMENT LABELS' at line 1
> exec
| SHOW PLACEMENT LABELS LIKE '%zone%'
< error
| 1064 42000 …near 'PLACEMENT LABELS LIKE '%zone%'' at line 1
> exec
| SHOW PLACEMENT FOR DATABASE db1
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1' at line 1
> exec
| SHOW PLACEMENT FOR SCHEMA db1
< error
| 1064 42000 …near 'PLACEMENT FOR SCHEMA db1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE db1.tb1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE db1.tb1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION p1' at line 1
> exec
| SHOW PLACEMENT FOR
< error
| 1064 42000 …near 'PLACEMENT FOR' at line 1
> exec
| SHOW PLACEMENT FOR DB db1
< error
| 1064 42000 …near 'PLACEMENT FOR DB db1' at line 1
> exec
| SHOW PLACEMENT FOR PARTITION p1
< error
| 1064 42000 …near 'PLACEMENT FOR PARTITION p1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 PARTITION
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 LIKE 't'
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 LIKE 't'' at line 1
> exec
| SHOW PLACEMENT FOR DATABASE db1 LIKE 't'
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1 LIKE 't'' at line 1
> exec
| SHOW TABLE t1 REGIONS
< error
| 1064 42000 …near 't1 REGIONS' at line 1
> exec
| SHOW TABLE t1 REGIONS WHERE a=1
< error
| 1064 42000 …near 't1 REGIONS WHERE a=1' at line 1
> exec
| SHOW TABLE t1 REGIONS LIKE 'x'
< error
| 1064 42000 …near 't1 REGIONS LIKE 'x'' at line 1
> exec
| SHOW TABLE test.t1 REGIONS
< error
| 1064 42000 …near 'test.t1 REGIONS' at line 1
> exec
| SHOW TABLE t1 PARTITION (p0) REGIONS
< error
| 1064 42000 …near 't1 PARTITION (p0) REGIONS' at line 1
> exec
| SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1
< error
| 1064 42000 …near 't1 PARTITION (p0,p1) REGIONS WHERE a=1' at line 1
> exec
| SHOW TABLE t1 PARTITION p0 REGIONS
< error
| 1064 42000 …near 't1 PARTITION p0 REGIONS' at line 1
> exec
| SHOW TABLE t1 PARTITION () REGIONS
< error
| 1064 42000 …near 't1 PARTITION () REGIONS' at line 1
> exec
| SHOW TABLE t1 INDEX idx1 REGIONS
< error
| 1064 42000 …near 't1 INDEX idx1 REGIONS' at line 1
> exec
| SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2
< error
| 1064 42000 …near 't1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2' at line 1
> exec
| SHOW TABLE t1 INDEX REGIONS
< error
| 1064 42000 …near 't1 INDEX REGIONS' at line 1
> exec
| SHOW TABLE t1 INDEX idx1
< error
| 1064 42000 …near 't1 INDEX idx1' at line 1
> exec
| SHOW TABLE t1 DISTRIBUTIONS
< error
| 1064 42000 …near 't1 DISTRIBUTIONS' at line 1
> exec
| SHOW TABLE t1 DISTRIBUTIONS WHERE a=1
< error
| 1064 42000 …near 't1 DISTRIBUTIONS WHERE a=1' at line 1
> exec
| SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS
< error
| 1064 42000 …near 't1 PARTITION (p0,p1) DISTRIBUTIONS' at line 1
> exec
| SHOW TABLE t1 NEXT_ROW_ID
< error
| 1064 42000 …near 't1 NEXT_ROW_ID' at line 1
> exec
| SHOW TABLE test.t1 NEXT_ROW_ID
< error
| 1064 42000 …near 'test.t1 NEXT_ROW_ID' at line 1
> exec
| SHOW TABLE t1 PARTITION (p0) NEXT_ROW_ID
< error
| 1064 42000 …near 't1 PARTITION (p0) NEXT_ROW_ID' at line 1
> exec
| SHOW TABLE t1 NEXT_ROW_ID WHERE a=1
< error
| 1064 42000 …near 't1 NEXT_ROW_ID WHERE a=1' at line 1
> exec
| SHOW TABLE t1
< error
| 1064 42000 …near 't1' at line 1
> exec
| SHOW TABLE t1 REGIONS REGIONS
< error
| 1064 42000 …near 't1 REGIONS REGIONS' at line 1
> exec
| SHOW DISTRIBUTION JOBS
< error
| 1064 42000 …near 'DISTRIBUTION JOBS' at line 1
> exec
| SHOW DISTRIBUTION JOBS WHERE id > 0
< error
| 1064 42000 …near 'DISTRIBUTION JOBS WHERE id > 0' at line 1
> exec
| SHOW DISTRIBUTION JOB 1
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1' at line 1
> exec
| SHOW DISTRIBUTION JOB 'x'
< error
| 1064 42000 …near 'DISTRIBUTION JOB 'x'' at line 1
> exec
| SHOW DISTRIBUTION JOB
< error
| 1064 42000 …near 'DISTRIBUTION JOB' at line 1
> exec
| SHOW DISTRIBUTION JOB 1 WHERE id > 0
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1 WHERE id > 0' at line 1
> exec
| SHOW DISTRIBUTION JOBS 1
< error
| 1064 42000 …near 'DISTRIBUTION JOBS 1' at line 1
> exec
| SHOW DISTRIBUTION
< error
| 1064 42000 …near 'DISTRIBUTION' at line 1
> exec
| SHOW SESSION_STATES
< error
| 1064 42000 …near 'SESSION_STATES' at line 1
> exec
| SHOW STATS_EXTENDED
< error
| 1064 42000 …near 'STATS_EXTENDED' at line 1
> exec
| SHOW STATS_META
< error
| 1064 42000 …near 'STATS_META' at line 1
> exec
| SHOW STATS_META WHERE table_name = 't'
< error
| 1064 42000 …near 'STATS_META WHERE table_name = 't'' at line 1
> exec
| SHOW STATS_META LIKE 't'
< error
| 1064 42000 …near 'STATS_META LIKE 't'' at line 1
> exec
| SHOW STATS_HISTOGRAMS
< error
| 1064 42000 …near 'STATS_HISTOGRAMS' at line 1
> exec
| SHOW STATS_TOPN
< error
| 1064 42000 …near 'STATS_TOPN' at line 1
> exec
| SHOW STATS_BUCKETS
< error
| 1064 42000 …near 'STATS_BUCKETS' at line 1
> exec
| SHOW STATS_HEALTHY
< error
| 1064 42000 …near 'STATS_HEALTHY' at line 1
> exec
| SHOW STATS_LOCKED
< error
| 1064 42000 …near 'STATS_LOCKED' at line 1
> exec
| SHOW HISTOGRAMS_IN_FLIGHT
< error
| 1064 42000 …near 'HISTOGRAMS_IN_FLIGHT' at line 1
> exec
| SHOW COLUMN_STATS_USAGE
< error
| 1064 42000 …near 'COLUMN_STATS_USAGE' at line 1
> exec
| SHOW ANALYZE STATUS
< error
| 1064 42000 …near 'ANALYZE STATUS' at line 1
> exec
| SHOW ANALYZE STATUS WHERE table_name = 't'
< error
| 1064 42000 …near 'ANALYZE STATUS WHERE table_name = 't'' at line 1
> exec
| SHOW ANALYZE
< error
| 1064 42000 …near 'ANALYZE' at line 1
> exec
| SHOW BUILTINS
< error
| 1064 42000 …near 'BUILTINS' at line 1
> exec
| SHOW BUILTINS LIKE 'x'
< error
| 1064 42000 …near 'BUILTINS LIKE 'x'' at line 1
> exec
| SHOW BACKUPS
< error
| 1064 42000 …near 'BACKUPS' at line 1
> exec
| SHOW BACKUPS WHERE start_time > now()
< error
| 1064 42000 …near 'BACKUPS WHERE start_time > now()' at line 1
> exec
| SHOW RESTORES
< error
| 1064 42000 …near 'RESTORES' at line 1
> exec
| SHOW RESTORES LIKE 'r0001'
< error
| 1064 42000 …near 'RESTORES LIKE 'r0001'' at line 1
> exec
| SHOW CONFIG
< error
| 1064 42000 …near 'CONFIG' at line 1
> exec
| SHOW CONFIG WHERE type='tidb'
< error
| 1064 42000 …near 'CONFIG WHERE type='tidb'' at line 1
> exec
| SHOW DATABASES TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| SHOW TABLES TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| SHOW CREATE TABLE t TABLE t2
< error
| 1064 42000 …near 'TABLE t2' at line 1
> exec
| SHOW TABLE STATUS TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| SHOW VARIABLES LIKE 'x%' TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec
| SHOW PLACEMENT TABLE t
< error
| 1064 42000 …near 'PLACEMENT TABLE t' at line 1
> exec
| SHOW PLACEMENT FOR DATABASE db1 TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1 TABLE tb1' at line 1
> exec
| SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR SCHEMA db1 TABLE tb1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 TABLE tb2
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 TABLE tb2' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 DATABASE db1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 DATABASE db1' at line 1
> exec
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2' at line 1
> exec
| SHOW TABLE t1 REGIONS TABLE t2
< error
| 1064 42000 …near 't1 REGIONS TABLE t2' at line 1
> exec
| SHOW TABLE t1 NEXT_ROW_ID TABLE t2
< error
| 1064 42000 …near 't1 NEXT_ROW_ID TABLE t2' at line 1
> exec
| SHOW BUILTINS TABLE t
< error
| 1064 42000 …near 'BUILTINS TABLE t' at line 1
> exec
| SHOW DISTRIBUTION JOB 1 TABLE t
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1 TABLE t' at line 1
> exec
| SHOW STATS_META TABLE t
< error
| 1064 42000 …near 'STATS_META TABLE t' at line 1
> exec
| DROP DATABASE IF EXISTS obshowt
< ok 0 0
