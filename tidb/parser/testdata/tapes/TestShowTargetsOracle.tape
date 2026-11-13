-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestShowTargetsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SHOW
< error
| 1064 42000 …line 1 column 4 near "" 
> exec ^
: string
| SHOW FULL
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| SHOW EXTENDED
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| SHOW CREATE
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SHOW TABLE
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| SHOW GLOBAL
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SHOW SESSION
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| SHOW MASTER
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SHOW SLAVE
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| SHOW REPLICA
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| SHOW BINARY
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SHOW BINLOG
< error
| 1064 42000 …line 1 column 11 near "BINLOG" 
> exec ^
: string
| SHOW RELAYLOG
< error
| 1064 42000 …line 1 column 13 near "RELAYLOG" 
> exec ^
: string
| SHOW STORAGE
< error
| 1064 42000 …line 1 column 12 near "STORAGE" 
> exec ^
: string
| SHOW OPEN
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| SHOW CHARACTER
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| SHOW ENGINE
< error
| 1064 42000 …line 1 column 11 near "ENGINE" 
> exec ^
: string
| SHOW PROCEDURE
< error
| 1064 42000 …line 1 column 14 near "" 
> exec ^
: string
| SHOW FUNCTION
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| SHOW COUNT
< error
| 1064 42000 …line 1 column 10 near "COUNT" 
> exec ^
: string
| SHOW COUNT(*)
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| SHOW COUNT ERRORS
< error
| 1064 42000 …line 1 column 10 near "COUNT ERRORS" 
> exec ^
: string
| SHOW COUNT(*) ERRORS
< ok 0 0
> exec ^
: string
| SHOW COUNT(*) WARNINGS
< ok 0 0
> exec ^
: string
| SHOW COUNT(*) TABLES
< error
| 1064 42000 …line 1 column 20 near "TABLES" 
> exec ^
: string
| SHOW DATABASES
< ok 0 0
> exec ^
: string
| SHOW SCHEMAS
< ok 0 0
> exec ^
: string
| SHOW TABLES
< ok 0 0
> exec ^
: string
| SHOW TABLES FROM test
< ok 0 0
> exec ^
: string
| SHOW TABLES IN test
< ok 0 0
> exec ^
: string
| SHOW FULL TABLES
< ok 0 0
> exec ^
: string
| SHOW FULL COLUMNS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW FULL COLUMNS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW FULL FIELDS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW FULL FIELDS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW FULL PROCESSLIST
< ok 0 0
> exec ^
: string
| SHOW FULL STATUS
< error
| 1064 42000 …line 1 column 16 near "STATUS" 
> exec ^
: string
| SHOW COLUMNS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW COLUMNS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW COLUMNS FROM t FROM test
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW COLUMNS FROM t FROM test
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW COLUMNS
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| SHOW COLUMNS t
< error
| 1064 42000 …line 1 column 14 near "t" 
> exec ^
: string
| SHOW FIELDS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW FIELDS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW EXTENDED COLUMNS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW EXTENDED COLUMNS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW EXTENDED FIELDS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW EXTENDED FIELDS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW EXTENDED FULL COLUMNS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW EXTENDED FULL COLUMNS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW EXTENDED FULL FIELDS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW EXTENDED FULL FIELDS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW EXTENDED TABLES
< error
| 1064 42000 …line 1 column 20 near "TABLES" 
> exec ^
: string
| SHOW EXTENDED FULL TABLES
< error
| 1064 42000 …line 1 column 25 near "TABLES" 
> exec ^
: string
| SHOW EXTENDED INDEX FROM t
< error
| 1064 42000 …line 1 column 19 near "INDEX FROM t" 
> exec ^
: string
| SHOW EXTENDED KEYS FROM t
< error
| 1064 42000 …line 1 column 18 near "KEYS FROM t" 
> exec ^
: string
| SHOW CREATE TABLE t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW CREATE TABLE t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW CREATE DATABASE d
< ok 0 0
> exec ^
: string
| SHOW CREATE DATABASE IF NOT EXISTS d
< ok 0 0
> exec ^
: string
| SHOW CREATE DATABASE IF d
< error
| 1064 42000 …line 1 column 25 near "d" 
> exec ^
: string
| SHOW CREATE SCHEMA d
< ok 0 0
> exec ^
: string
| SHOW CREATE VIEW v
< error
| 1146 42S02 Table 'test.v' doesn't exist
> exec ^
: string
| SHOW CREATE VIEW v
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW CREATE PROCEDURE pr
< ok 0 0
> exec ^
: string
| SHOW CREATE FUNCTION fn
< error
| 1064 42000 …line 1 column 20 near "FUNCTION fn" 
> exec ^
: string
| SHOW CREATE TRIGGER tg
< error
| 1064 42000 …line 1 column 19 near "TRIGGER tg" 
> exec ^
: string
| SHOW CREATE EVENT ev
< error
| 1064 42000 …line 1 column 17 near "EVENT ev" 
> exec ^
: string
| SHOW CREATE USER u
< ok 0 0
> exec ^
: string
| SHOW CREATE USER
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| SHOW CREATE SEQUENCE seq
< error
| 1146 42S02 Table 'test.seq' doesn't exist
> exec ^
: string
| SHOW CREATE SEQUENCE seq
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW CREATE SEQUENCE test.seq
< error
| 1146 42S02 Table 'test.seq' doesn't exist
> exec ^
: string
| SHOW CREATE SEQUENCE test.seq
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW CREATE PLACEMENT POLICY x
< ok 0 0
> exec ^
: string
| SHOW CREATE PLACEMENT POLICY
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| SHOW CREATE RESOURCE GROUP rg
< ok 0 0
> exec ^
: string
| SHOW CREATE RESOURCE GROUP default
< ok 0 0
> exec ^
: string
| SHOW CREATE RESOURCE GROUP
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| SHOW CREATE RESOURCE rg
< error
| 1064 42000 …line 1 column 23 near "rg" 
> exec ^
: string
| SHOW INDEX FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW INDEX FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW INDEX IN t IN test
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW INDEX IN t IN test
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW INDEX
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| SHOW INDEX t
< error
| 1064 42000 …line 1 column 12 near "t" 
> exec ^
: string
| SHOW KEY FROM t
< error
| 1064 42000 …line 1 column 8 near "KEY FROM t" 
> exec ^
: string
| SHOW KEYS FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW KEYS FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW INDEXES FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SHOW INDEXES FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW GLOBAL VARIABLES
< ok 0 0
> exec ^
: string
| SHOW GLOBAL STATUS
< ok 0 0
> exec ^
: string
| SHOW SESSION VARIABLES
< ok 0 0
> exec ^
: string
| SHOW SESSION STATUS
< ok 0 0
> exec ^
: string
| SHOW GLOBAL BINDINGS
< ok 0 0
> exec ^
: string
| SHOW SESSION BINDINGS
< ok 0 0
> exec ^
: string
| SHOW BINDINGS
< ok 0 0
> exec ^
: string
| SHOW BINDING_CACHE STATUS
< ok 0 0
> exec ^
: string
| SHOW BINDING_CACHE
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| SHOW VARIABLES
< ok 0 0
> exec ^
: string
| SHOW VARIABLES LIKE 'x%'
< ok 0 0
> exec ^
: string
| SHOW STATUS
< ok 0 0
> exec ^
: string
| SHOW WARNINGS
< ok 0 0
> exec ^
: string
| SHOW WARNINGS LIMIT 10
< error
| 1064 42000 …line 1 column 19 near "LIMIT 10" 
> exec ^
: string
| SHOW WARNINGS LIMIT 5, 10
< error
| 1064 42000 …line 1 column 19 near "LIMIT 5, 10" 
> exec ^
: string
| SHOW ERRORS
< ok 0 0
> exec ^
: string
| SHOW ERRORS LIMIT 10
< error
| 1064 42000 …line 1 column 17 near "LIMIT 10" 
> exec ^
: string
| SHOW ENGINES
< ok 0 0
> exec ^
: string
| SHOW STORAGE ENGINES
< error
| 1064 42000 …line 1 column 12 near "STORAGE ENGINES" 
> exec ^
: string
| SHOW ENGINE INNODB STATUS
< error
| 1064 42000 …line 1 column 11 near "ENGINE INNODB STATUS" 
> exec ^
: string
| SHOW ENGINE INNODB MUTEX
< error
| 1064 42000 …line 1 column 11 near "ENGINE INNODB MUTEX" 
> exec ^
: string
| SHOW ENGINE INNODB LOGS
< error
| 1064 42000 …line 1 column 11 near "ENGINE INNODB LOGS" 
> exec ^
: string
| SHOW ENGINE INNODB
< error
| 1064 42000 …line 1 column 11 near "ENGINE INNODB" 
> exec ^
: string
| SHOW PLUGINS
< ok 0 0
> exec ^
: string
| SHOW MASTER STATUS
< ok 0 0
> exec ^
: string
| SHOW MASTER LOGS
< error
| 1064 42000 …line 1 column 16 near "LOGS" 
> exec ^
: string
| SHOW BINARY LOGS
< error
| 1064 42000 …line 1 column 16 near "LOGS" 
> exec ^
: string
| SHOW BINARY LOG STATUS
< ok 0 0
> exec ^
: string
| SHOW BINARY LOG
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| SHOW BINLOG STATUS
< error
| 1064 42000 …line 1 column 11 near "BINLOG STATUS" 
> exec ^
: string
| SHOW SLAVE STATUS
< error
| 1235 42000 This version of TiDB doesn't yet support 'SHOW {REPLICA | SLAVE} STATUS'
> exec ^
: string
| SHOW SLAVE STATUS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW SLAVE HOSTS
< error
| 1064 42000 …line 1 column 16 near "HOSTS" 
> exec ^
: string
| SHOW REPLICA STATUS
< error
| 1235 42000 This version of TiDB doesn't yet support 'SHOW {REPLICA | SLAVE} STATUS'
> exec ^
: string
| SHOW REPLICA STATUS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW REPLICA HOSTS
< error
| 1064 42000 …line 1 column 18 near "HOSTS" 
> exec ^
: string
| SHOW REPLICAS
< error
| 1064 42000 …line 1 column 13 near "REPLICAS" 
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL 'c'
< error
| 1064 42000 …line 1 column 21 near "FOR CHANNEL 'c'" 
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL c
< error
| 1064 42000 …line 1 column 21 near "FOR CHANNEL c" 
> exec ^
: string
| SHOW SLAVE STATUS FOR
< error
| 1064 42000 …line 1 column 21 near "FOR" 
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL
< error
| 1064 42000 …line 1 column 21 near "FOR CHANNEL" 
> exec ^
: string
| SHOW REPLICA STATUS FOR CHANNEL 'c'
< error
| 1064 42000 …line 1 column 23 near "FOR CHANNEL 'c'" 
> exec ^
: string
| SHOW REPLICA STATUS FOR CHANNEL c
< error
| 1064 42000 …line 1 column 23 near "FOR CHANNEL c" 
> exec ^
: string
| SHOW BINLOG EVENT
< error
| 1064 42000 …line 1 column 11 near "BINLOG EVENT" 
> exec ^
: string
| SHOW BINLOG EVENTS
< error
| 1064 42000 …line 1 column 11 near "BINLOG EVENTS" 
> exec ^
: string
| SHOW BINLOG EVENTS IN 'f'
< error
| 1064 42000 …line 1 column 11 near "BINLOG EVENTS IN 'f'" 
> exec ^
: string
| SHOW BINLOG EVENTS FROM 4
< error
| 1064 42000 …line 1 column 11 near "BINLOG EVENTS FROM 4" 
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 5, 10
< error
| 1064 42000 …line 1 column 11 near "BINLOG EVENTS LIMIT 5, 10" 
> exec ^
: string
| SHOW RELAYLOG EVENT
< error
| 1064 42000 …line 1 column 13 near "RELAYLOG EVENT" 
> exec ^
: string
| SHOW RELAYLOG EVENTS
< error
| 1064 42000 …line 1 column 13 near "RELAYLOG EVENTS" 
> exec ^
: string
| SHOW RELAYLOG EVENTS FOR CHANNEL 'c'
< error
| 1064 42000 …line 1 column 13 near "RELAYLOG EVENTS FOR CHANNEL 'c'" 
> exec ^
: string
| SHOW RELAYLOG EVENTS FOR CHANNEL c
< error
| 1064 42000 …line 1 column 13 near "RELAYLOG EVENTS FOR CHANNEL c" 
> exec ^
: string
| SHOW TABLE STATUS
< ok 0 0
> exec ^
: string
| SHOW TABLE STATUS FROM test
< ok 0 0
> exec ^
: string
| SHOW TABLE STATUS LIKE 't'
< ok 0 0
> exec ^
: string
| SHOW TRIGGER
< error
| 1064 42000 …line 1 column 12 near "TRIGGER" 
> exec ^
: string
| SHOW TRIGGERS
< ok 0 0
> exec ^
: string
| SHOW EVENT
< error
| 1064 42000 …line 1 column 10 near "EVENT" 
> exec ^
: string
| SHOW EVENTS
< ok 0 0
> exec ^
: string
| SHOW PROCEDURE STATUS
< ok 0 0
> exec ^
: string
| SHOW PROCEDURE CODE pr
< error
| 1064 42000 …line 1 column 19 near "CODE pr" 
> exec ^
: string
| SHOW FUNCTION STATUS
< ok 0 0
> exec ^
: string
| SHOW FUNCTION CODE fn
< error
| 1064 42000 …line 1 column 18 near "CODE fn" 
> exec ^
: string
| SHOW OPEN TABLES
< ok 0 0
> exec ^
: string
| SHOW PRIVILEGES
< ok 0 0
> exec ^
: string
| SHOW PROFILES
< ok 0 0
> exec ^
: string
| SHOW PROFILE
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU, MEMORY FOR QUERY 2
< ok 0 0
> exec ^
: string
| SHOW PROFILE BLOCK IO
< ok 0 0
> exec ^
: string
| SHOW PROFILE FOR QUERY 2
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,1
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU FOR QUERY 2 LIMIT 1 OFFSET 2
< ok 0 0
> exec ^
: string
| SHOW PROFILE CPU LIMIT 1,1 OFFSET 2
< error
| 1064 42000 …line 1 column 33 near "OFFSET 2" 
> exec ^
: string
| SHOW PROFILE CPU LIMIT 1,
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| SHOW PROFILE CPU OFFSET 2
< error
| 1064 42000 …line 1 column 23 near "OFFSET 2" 
> exec ^
: string
| SHOW CHARACTER SET
< ok 0 0
> exec ^
: string
| SHOW CHARSET
< ok 0 0
> exec ^
: string
| SHOW CHAR SET
< ok 0 0
> exec ^
: string
| SHOW CHAR
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| SHOW COLLATION
< ok 0 0
> exec ^
: string
| SHOW GRANTS
< ok 0 0
> exec ^
: string
| SHOW GRANTS FOR u
< ok 0 0
> exec ^
: string
| SHOW PROCESSLIST
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 't'
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE @v
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE concat('a','%')
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 1
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE -1
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE NULL
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE ('t')
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE binary
< error
| 1253 42000 COLLATION 'binary' is not valid for CHARACTER SET 'utf8mb4'
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE 1
< error
| 1064 42000 …line 1 column 34 near "1" 
> exec ^
: string
| SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin
< ok 0 0
> exec ^
: string
| SHOW DATABASES LIKE 'a' || 'b'
< error
| 1064 42000 …line 1 column 26 near "|| 'b'" 
> exec ^
: string
| SHOW DATABASES LIKE 'a' + 1
< error
| 1064 42000 …line 1 column 25 near "+ 1" 
> exec ^
: string
| SHOW DATABASES LIKE 'a' -> '$.x'
< error
| 1064 42000 …line 1 column 26 near "-> '$.x'" 
> exec ^
: string
| SHOW DATABASES LIKE
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| SHOW DATABASES LIKE 't' LIKE 'u'
< error
| 1064 42000 …line 1 column 28 near "LIKE 'u'" 
> exec ^
: string
| SHOW DATABASES LIKE 't' AND 1
< error
| 1064 42000 …line 1 column 27 near "AND 1" 
> exec ^
: string
| SHOW DATABASES LIKE 't' WHERE 1=1
< error
| 1064 42000 …line 1 column 29 near "WHERE 1=1" 
> exec ^
: string
| SHOW DATABASES WHERE 1=1
< ok 0 0
> exec ^
: string
| SHOW DATABASES WHERE
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| SHOW DATABASES WHERE 1=1 LIKE 't'
< ok 0 0
> exec ^
: string
| SHOW TABLES LIKE 't' WHERE 1=1
< error
| 1064 42000 …line 1 column 26 near "WHERE 1=1" 
> exec ^
: string
| SHOW PLACEMENT
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT LIKE 'POLICY foo%'
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT WHERE Target='TABLE test.t1'
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT LABELS
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT LABELS LIKE '%zone%'
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT FOR SCHEMA db1
< ok 0 0
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1
< error
| 1146 42S02 Table 'test.tb1' doesn't exist
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE db1.tb1
< error
| 1146 42S02 Table 'db1.tb1' doesn't exist
> exec ^
: string
| SHOW PLACEMENT FOR TABLE db1.tb1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1
< error
| 1146 42S02 Table 'test.tb1' doesn't exist
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW PLACEMENT FOR
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| SHOW PLACEMENT FOR DB db1
< error
| 1064 42000 …line 1 column 21 near "DB db1" 
> exec ^
: string
| SHOW PLACEMENT FOR PARTITION p1
< error
| 1064 42000 …line 1 column 28 near "PARTITION p1" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION
< error
| 1064 42000 …line 1 column 38 near "" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 LIKE 't'
< error
| 1064 42000 …line 1 column 33 near "LIKE 't'" 
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1 LIKE 't'
< error
| 1064 42000 …line 1 column 36 near "LIKE 't'" 
> exec ^
: string
| SHOW TABLE t1 REGIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 REGIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 REGIONS WHERE a=1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 REGIONS WHERE a=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 REGIONS LIKE 'x'
< error
| 1064 42000 …line 1 column 26 near "LIKE 'x'" 
> exec ^
: string
| SHOW TABLE test.t1 REGIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE test.t1 REGIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) REGIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) REGIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION p0 REGIONS
< error
| 1064 42000 …line 1 column 26 near "p0 REGIONS" 
> exec ^
: string
| SHOW TABLE t1 PARTITION () REGIONS
< error
| 1064 42000 …line 1 column 26 near ") REGIONS" 
> exec ^
: string
| SHOW TABLE t1 INDEX idx1 REGIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 INDEX idx1 REGIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 INDEX REGIONS
< error
| 1064 42000 …line 1 column 27 near "" 
> exec ^
: string
| SHOW TABLE t1 INDEX idx1
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS WHERE a=1
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS WHERE a=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE test.t1 NEXT_ROW_ID
< error
| 1146 42S02 Table 'test.t1' doesn't exist
> exec ^
: string
| SHOW TABLE test.t1 NEXT_ROW_ID
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) NEXT_ROW_ID
< error
| 1064 42000 …line 1 column 40 near "NEXT_ROW_ID" 
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID WHERE a=1
< error
| 1064 42000 …line 1 column 31 near "WHERE a=1" 
> exec ^
: string
| SHOW TABLE t1
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| SHOW TABLE t1 REGIONS REGIONS
< error
| 1064 42000 …line 1 column 29 near "REGIONS" 
> exec ^
: string
| SHOW DISTRIBUTION JOBS
< ok 0 0
> exec ^
: string
| SHOW DISTRIBUTION JOBS WHERE id > 0
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| SHOW DISTRIBUTION JOBS WHERE id > 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW DISTRIBUTION JOB 1
< ok 0 0
> exec ^
: string
| SHOW DISTRIBUTION JOB 'x'
< error
| 1064 42000 …line 1 column 25 near "'x'" 
> exec ^
: string
| SHOW DISTRIBUTION JOB
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| SHOW DISTRIBUTION JOB 1 WHERE id > 0
< error
| 1064 42000 …line 1 column 29 near "WHERE id > 0" 
> exec ^
: string
| SHOW DISTRIBUTION JOBS 1
< error
| 1064 42000 …line 1 column 24 near "1" 
> exec ^
: string
| SHOW IMPORT JOBS
< ok 0 0
> exec ^
: string
| SHOW IMPORT JOBS WHERE id > 0
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| SHOW IMPORT JOBS WHERE id > 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW IMPORT JOB 1
< ok 0 0
> exec ^
: string
| SHOW IMPORT JOB 'x'
< error
| 1064 42000 …line 1 column 19 near "'x'" 
> exec ^
: string
| SHOW IMPORT JOB
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| SHOW IMPORT JOB 1 WHERE id > 0
< error
| 1064 42000 …line 1 column 23 near "WHERE id > 0" 
> exec ^
: string
| SHOW IMPORT JOBS 1
< error
| 1064 42000 …line 1 column 18 near "1" 
> exec ^
: string
| SHOW IMPORT
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| SHOW AFFINITY
< ok 0 0
> exec ^
: string
| SHOW AFFINITY WHERE id > 0
< error
| 1054 42S22 Unknown column 'id' in 'where clause'
> exec ^
: string
| SHOW AFFINITY WHERE id > 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW AFFINITY LIKE 'x'
< ok 0 0
> exec ^
: string
| SHOW DISTRIBUTION
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| SHOW SESSION_STATES
< ok 0 0
> exec ^
: string
| SHOW STATS_EXTENDED
< ok 0 0
> exec ^
: string
| SHOW STATS_META
< ok 0 0
> exec ^
: string
| SHOW STATS_META WHERE table_name = 't'
< ok 0 0
> exec ^
: string
| SHOW STATS_META LIKE 't'
< ok 0 0
> exec ^
: string
| SHOW STATS_HISTOGRAMS
< ok 0 0
> exec ^
: string
| SHOW STATS_TOPN
< ok 0 0
> exec ^
: string
| SHOW STATS_BUCKETS
< ok 0 0
> exec ^
: string
| SHOW STATS_HEALTHY
< ok 0 0
> exec ^
: string
| SHOW STATS_LOCKED
< ok 0 0
> exec ^
: string
| SHOW HISTOGRAMS_IN_FLIGHT
< ok 0 0
> exec ^
: string
| SHOW COLUMN_STATS_USAGE
< ok 0 0
> exec ^
: string
| SHOW ANALYZE STATUS
< ok 0 0
> exec ^
: string
| SHOW ANALYZE STATUS WHERE table_name = 't'
< ok 0 0
> exec ^
: string
| SHOW ANALYZE
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| SHOW BUILTINS
< ok 0 0
> exec ^
: string
| SHOW BUILTINS LIKE 'x'
< error
| 1064 42000 …line 1 column 18 near "LIKE 'x'" 
> exec ^
: string
| SHOW BACKUPS
< ok 0 0
> exec ^
: string
| SHOW BACKUPS WHERE start_time > now()
< error
| 1054 42S22 Unknown column 'start_time' in 'where clause'
> exec ^
: string
| SHOW BACKUPS WHERE start_time > now()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SHOW RESTORES
< ok 0 0
> exec ^
: string
| SHOW RESTORES LIKE 'r0001'
< ok 0 0
> exec ^
: string
| SHOW CONFIG
< ok 0 0
> exec ^
: string
| SHOW CONFIG WHERE type='tidb'
< ok 0 0
> exec ^
: string
| SHOW DATABASES TABLE t
< error
| 1064 42000 …line 1 column 20 near "TABLE t" 
> exec ^
: string
| SHOW TABLES TABLE t
< error
| 1064 42000 …line 1 column 17 near "TABLE t" 
> exec ^
: string
| SHOW CREATE TABLE t TABLE t2
< error
| 1064 42000 …line 1 column 25 near "TABLE t2" 
> exec ^
: string
| SHOW TABLE STATUS TABLE t
< error
| 1064 42000 …line 1 column 23 near "TABLE t" 
> exec ^
: string
| SHOW VARIABLES LIKE 'x%' TABLE t
< error
| 1064 42000 …line 1 column 30 near "TABLE t" 
> exec ^
: string
| SHOW PLACEMENT TABLE t
< error
| 1064 42000 …line 1 column 20 near "TABLE t" 
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1 TABLE tb1
< error
| 1064 42000 …line 1 column 37 near "TABLE tb1" 
> exec ^
: string
| SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1
< error
| 1064 42000 …line 1 column 35 near "TABLE tb1" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 TABLE tb2
< error
| 1064 42000 …line 1 column 34 near "TABLE tb2" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 DATABASE db1
< error
| 1064 42000 …line 1 column 37 near "DATABASE db1" 
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2
< error
| 1064 42000 …line 1 column 51 near "PARTITION p2" 
> exec ^
: string
| SHOW TABLE t1 REGIONS TABLE t2
< error
| 1064 42000 …line 1 column 27 near "TABLE t2" 
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID TABLE t2
< error
| 1064 42000 …line 1 column 31 near "TABLE t2" 
> exec ^
: string
| SHOW BUILTINS TABLE t
< error
| 1064 42000 …line 1 column 19 near "TABLE t" 
> exec ^
: string
| SHOW DISTRIBUTION JOB 1 TABLE t
< error
| 1064 42000 …line 1 column 29 near "TABLE t" 
> exec ^
: string
| SHOW IMPORT JOB 1 TABLE t
< error
| 1064 42000 …line 1 column 23 near "TABLE t" 
> exec ^
: string
| SHOW STATS_META TABLE t
< error
| 1064 42000 …line 1 column 21 near "TABLE t" 
