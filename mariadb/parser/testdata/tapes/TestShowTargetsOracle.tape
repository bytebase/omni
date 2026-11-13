-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestShowTargetsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SHOW
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW FULL
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW EXTENDED
< error
| 1064 42000 …near 'EXTENDED' at line 1
> exec ^
: string
| SHOW CREATE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW TABLE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW GLOBAL
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW SESSION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW MASTER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW SLAVE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW REPLICA
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW BINARY
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW BINLOG
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW RELAYLOG
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW STORAGE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW OPEN
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW CHARACTER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW ENGINE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW PROCEDURE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW FUNCTION
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW COUNT
< error
| 1064 42000 …near 'COUNT' at line 1
> exec ^
: string
| SHOW COUNT(*)
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW COUNT ERRORS
< error
| 1064 42000 …near 'COUNT ERRORS' at line 1
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
| 1064 42000 …near 'TABLES' at line 1
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
< ok 0 0
> exec ^
: string
| SHOW FULL FIELDS FROM t
< ok 0 0
> exec ^
: string
| SHOW FULL PROCESSLIST
< ok 0 0
> exec ^
: string
| SHOW FULL STATUS
< error
| 1064 42000 …near 'STATUS' at line 1
> exec ^
: string
| SHOW COLUMNS FROM t
< ok 0 0
> exec ^
: string
| SHOW COLUMNS FROM t FROM test
< ok 0 0
> exec ^
: string
| SHOW COLUMNS
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW COLUMNS t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| SHOW FIELDS FROM t
< ok 0 0
> exec ^
: string
| SHOW EXTENDED COLUMNS FROM t
< error
| 1064 42000 …near 'EXTENDED COLUMNS FROM t' at line 1
> exec ^
: string
| SHOW EXTENDED FIELDS FROM t
< error
| 1064 42000 …near 'EXTENDED FIELDS FROM t' at line 1
> exec ^
: string
| SHOW EXTENDED FULL COLUMNS FROM t
< error
| 1064 42000 …near 'EXTENDED FULL COLUMNS FROM t' at line 1
> exec ^
: string
| SHOW EXTENDED FULL FIELDS FROM t
< error
| 1064 42000 …near 'EXTENDED FULL FIELDS FROM t' at line 1
> exec ^
: string
| SHOW EXTENDED TABLES
< error
| 1064 42000 …near 'EXTENDED TABLES' at line 1
> exec ^
: string
| SHOW EXTENDED FULL TABLES
< error
| 1064 42000 …near 'EXTENDED FULL TABLES' at line 1
> exec ^
: string
| SHOW EXTENDED INDEX FROM t
< error
| 1064 42000 …near 'EXTENDED INDEX FROM t' at line 1
> exec ^
: string
| SHOW EXTENDED KEYS FROM t
< error
| 1064 42000 …near 'EXTENDED KEYS FROM t' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'd' at line 1
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW CREATE PROCEDURE pr
< ok 0 0
> exec ^
: string
| SHOW CREATE FUNCTION fn
< ok 0 0
> exec ^
: string
| SHOW CREATE TRIGGER tg
< ok 0 0
> exec ^
: string
| SHOW CREATE EVENT ev
< ok 0 0
> exec ^
: string
| SHOW CREATE USER u
< ok 0 0
> exec ^
: string
| SHOW CREATE USER
< ok 0 0
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
| 1064 42000 …near 'SELECT FROM' at line 2
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
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SHOW CREATE PLACEMENT POLICY x
< error
| 1064 42000 …near 'PLACEMENT POLICY x' at line 1
> exec ^
: string
| SHOW CREATE PLACEMENT POLICY
< error
| 1064 42000 …near 'PLACEMENT POLICY' at line 1
> exec ^
: string
| SHOW CREATE RESOURCE GROUP rg
< error
| 1064 42000 …near 'RESOURCE GROUP rg' at line 1
> exec ^
: string
| SHOW CREATE RESOURCE GROUP default
< error
| 1064 42000 …near 'RESOURCE GROUP default' at line 1
> exec ^
: string
| SHOW CREATE RESOURCE GROUP
< error
| 1064 42000 …near 'RESOURCE GROUP' at line 1
> exec ^
: string
| SHOW CREATE RESOURCE rg
< error
| 1064 42000 …near 'RESOURCE rg' at line 1
> exec ^
: string
| SHOW INDEX FROM t
< ok 0 0
> exec ^
: string
| SHOW INDEX IN t IN test
< ok 0 0
> exec ^
: string
| SHOW INDEX
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW INDEX t
< error
| 1064 42000 …near 't' at line 1
> exec ^
: string
| SHOW KEY FROM t
< error
| 1064 42000 …near 'KEY FROM t' at line 1
> exec ^
: string
| SHOW KEYS FROM t
< ok 0 0
> exec ^
: string
| SHOW INDEXES FROM t
< ok 0 0
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
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec ^
: string
| SHOW SESSION BINDINGS
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec ^
: string
| SHOW BINDINGS
< error
| 1064 42000 …near 'BINDINGS' at line 1
> exec ^
: string
| SHOW BINDING_CACHE STATUS
< error
| 1064 42000 …near 'BINDING_CACHE STATUS' at line 1
> exec ^
: string
| SHOW BINDING_CACHE
< error
| 1064 42000 …near 'BINDING_CACHE' at line 1
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
< ok 0 0
> exec ^
: string
| SHOW WARNINGS LIMIT 5, 10
< ok 0 0
> exec ^
: string
| SHOW ERRORS
< ok 0 0
> exec ^
: string
| SHOW ERRORS LIMIT 10
< ok 0 0
> exec ^
: string
| SHOW ENGINES
< ok 0 0
> exec ^
: string
| SHOW STORAGE ENGINES
< ok 0 0
> exec ^
: string
| SHOW ENGINE INNODB STATUS
< ok 0 0
> exec ^
: string
| SHOW ENGINE INNODB MUTEX
< ok 0 0
> exec ^
: string
| SHOW ENGINE INNODB LOGS
< ok 0 0
> exec ^
: string
| SHOW ENGINE INNODB
< error
| 1064 42000 …near '' at line 1
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
< ok 0 0
> exec ^
: string
| SHOW BINARY LOGS
< ok 0 0
> exec ^
: string
| SHOW BINARY LOG STATUS
< error
| 1064 42000 …near 'LOG STATUS' at line 1
> exec ^
: string
| SHOW BINARY LOG
< error
| 1064 42000 …near 'LOG' at line 1
> exec ^
: string
| SHOW BINLOG STATUS
< ok 0 0
> exec ^
: string
| SHOW SLAVE STATUS
< ok 0 0
> exec ^
: string
| SHOW SLAVE HOSTS
< ok 0 0
> exec ^
: string
| SHOW REPLICA STATUS
< ok 0 0
> exec ^
: string
| SHOW REPLICA HOSTS
< ok 0 0
> exec ^
: string
| SHOW REPLICAS
< error
| 1064 42000 …near 'REPLICAS' at line 1
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL 'c'
< ok 0 0
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL c
< error
| 1064 42000 …near 'c' at line 1
> exec ^
: string
| SHOW SLAVE STATUS FOR
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW SLAVE STATUS FOR CHANNEL
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW REPLICA STATUS FOR CHANNEL 'c'
< ok 0 0
> exec ^
: string
| SHOW REPLICA STATUS FOR CHANNEL c
< error
| 1064 42000 …near 'c' at line 1
> exec ^
: string
| SHOW BINLOG EVENT
< error
| 1064 42000 …near 'EVENT' at line 1
> exec ^
: string
| SHOW BINLOG EVENTS
< ok 0 0
> exec ^
: string
| SHOW BINLOG EVENTS IN 'f'
< ok 0 0
> exec ^
: string
| SHOW BINLOG EVENTS FROM 4
< ok 0 0
> exec ^
: string
| SHOW BINLOG EVENTS LIMIT 5, 10
< ok 0 0
> exec ^
: string
| SHOW RELAYLOG EVENT
< error
| 1064 42000 …near 'EVENT' at line 1
> exec ^
: string
| SHOW RELAYLOG EVENTS
< ok 0 0
> exec ^
: string
| SHOW RELAYLOG EVENTS FOR CHANNEL 'c'
< ok 0 0
> exec ^
: string
| SHOW RELAYLOG EVENTS FOR CHANNEL c
< error
| 1064 42000 …near 'c' at line 1
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
| 1064 42000 …near 'TRIGGER' at line 1
> exec ^
: string
| SHOW TRIGGERS
< ok 0 0
> exec ^
: string
| SHOW EVENT
< error
| 1064 42000 …near 'EVENT' at line 1
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
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
> exec ^
: string
| SHOW FUNCTION STATUS
< ok 0 0
> exec ^
: string
| SHOW FUNCTION CODE fn
< error
| 1289 HY000 The 'SHOW PROCEDURE|FUNCTION CODE' feature is disabled; you need MariaDB built with '--with-debug' to have it working
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
| 1064 42000 …near 'OFFSET 2' at line 1
> exec ^
: string
| SHOW PROFILE CPU LIMIT 1,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW PROFILE CPU OFFSET 2
< error
| 1064 42000 …near '' at line 1
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
| 1064 42000 …near '' at line 1
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
< error
| 1064 42000 …near '@v' at line 1
> exec ^
: string
| SHOW DATABASES LIKE concat('a','%')
< error
| 1064 42000 …near 'concat('a','%')' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| SHOW DATABASES LIKE -1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| SHOW DATABASES LIKE NULL
< error
| 1064 42000 …near 'NULL' at line 1
> exec ^
: string
| SHOW DATABASES LIKE ('t')
< error
| 1064 42000 …near '('t')' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE utf8mb4_bin
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE 'utf8mb4_bin'
< error
| 1064 42000 …near 'COLLATE 'utf8mb4_bin'' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE binary
< error
| 1064 42000 …near 'COLLATE binary' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE
< error
| 1064 42000 …near 'COLLATE' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't%' COLLATE 1
< error
| 1064 42000 …near 'COLLATE 1' at line 1
> exec ^
: string
| SHOW TABLES LIKE 't%' COLLATE utf8mb4_bin
< error
| 1064 42000 …near 'COLLATE utf8mb4_bin' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 'a' || 'b'
< error
| 1064 42000 …near '|| 'b'' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 'a' + 1
< error
| 1064 42000 …near '+ 1' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 'a' -> '$.x'
< error
| 1064 42000 …near '-> '$.x'' at line 1
> exec ^
: string
| SHOW DATABASES LIKE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't' LIKE 'u'
< error
| 1064 42000 …near 'LIKE 'u'' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't' AND 1
< error
| 1064 42000 …near 'AND 1' at line 1
> exec ^
: string
| SHOW DATABASES LIKE 't' WHERE 1=1
< error
| 1064 42000 …near 'WHERE 1=1' at line 1
> exec ^
: string
| SHOW DATABASES WHERE 1=1
< ok 0 0
> exec ^
: string
| SHOW DATABASES WHERE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW DATABASES WHERE 1=1 LIKE 't'
< ok 0 0
> exec ^
: string
| SHOW TABLES LIKE 't' WHERE 1=1
< error
| 1064 42000 …near 'WHERE 1=1' at line 1
> exec ^
: string
| SHOW PLACEMENT
< error
| 1064 42000 …near 'PLACEMENT' at line 1
> exec ^
: string
| SHOW PLACEMENT LIKE 'POLICY foo%'
< error
| 1064 42000 …near 'PLACEMENT LIKE 'POLICY foo%'' at line 1
> exec ^
: string
| SHOW PLACEMENT WHERE Target='TABLE test.t1'
< error
| 1064 42000 …near 'PLACEMENT WHERE Target='TABLE test.t1'' at line 1
> exec ^
: string
| SHOW PLACEMENT LABELS
< error
| 1064 42000 …near 'PLACEMENT LABELS' at line 1
> exec ^
: string
| SHOW PLACEMENT LABELS LIKE '%zone%'
< error
| 1064 42000 …near 'PLACEMENT LABELS LIKE '%zone%'' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR SCHEMA db1
< error
| 1064 42000 …near 'PLACEMENT FOR SCHEMA db1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE db1.tb1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE db1.tb1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION p1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR
< error
| 1064 42000 …near 'PLACEMENT FOR' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR DB db1
< error
| 1064 42000 …near 'PLACEMENT FOR DB db1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR PARTITION p1
< error
| 1064 42000 …near 'PLACEMENT FOR PARTITION p1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 LIKE 't'
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 LIKE 't'' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1 LIKE 't'
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1 LIKE 't'' at line 1
> exec ^
: string
| SHOW TABLE t1 REGIONS
< error
| 1064 42000 …near 't1 REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 REGIONS WHERE a=1
< error
| 1064 42000 …near 't1 REGIONS WHERE a=1' at line 1
> exec ^
: string
| SHOW TABLE t1 REGIONS LIKE 'x'
< error
| 1064 42000 …near 't1 REGIONS LIKE 'x'' at line 1
> exec ^
: string
| SHOW TABLE test.t1 REGIONS
< error
| 1064 42000 …near 'test.t1 REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) REGIONS
< error
| 1064 42000 …near 't1 PARTITION (p0) REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) REGIONS WHERE a=1
< error
| 1064 42000 …near 't1 PARTITION (p0,p1) REGIONS WHERE a=1' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION p0 REGIONS
< error
| 1064 42000 …near 't1 PARTITION p0 REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION () REGIONS
< error
| 1064 42000 …near 't1 PARTITION () REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 INDEX idx1 REGIONS
< error
| 1064 42000 …near 't1 INDEX idx1 REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2
< error
| 1064 42000 …near 't1 PARTITION (p0) INDEX idx1 REGIONS WHERE a=2' at line 1
> exec ^
: string
| SHOW TABLE t1 INDEX REGIONS
< error
| 1064 42000 …near 't1 INDEX REGIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 INDEX idx1
< error
| 1064 42000 …near 't1 INDEX idx1' at line 1
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS
< error
| 1064 42000 …near 't1 DISTRIBUTIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 DISTRIBUTIONS WHERE a=1
< error
| 1064 42000 …near 't1 DISTRIBUTIONS WHERE a=1' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0,p1) DISTRIBUTIONS
< error
| 1064 42000 …near 't1 PARTITION (p0,p1) DISTRIBUTIONS' at line 1
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID
< error
| 1064 42000 …near 't1 NEXT_ROW_ID' at line 1
> exec ^
: string
| SHOW TABLE test.t1 NEXT_ROW_ID
< error
| 1064 42000 …near 'test.t1 NEXT_ROW_ID' at line 1
> exec ^
: string
| SHOW TABLE t1 PARTITION (p0) NEXT_ROW_ID
< error
| 1064 42000 …near 't1 PARTITION (p0) NEXT_ROW_ID' at line 1
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID WHERE a=1
< error
| 1064 42000 …near 't1 NEXT_ROW_ID WHERE a=1' at line 1
> exec ^
: string
| SHOW TABLE t1
< error
| 1064 42000 …near 't1' at line 1
> exec ^
: string
| SHOW TABLE t1 REGIONS REGIONS
< error
| 1064 42000 …near 't1 REGIONS REGIONS' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOBS
< error
| 1064 42000 …near 'DISTRIBUTION JOBS' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOBS WHERE id > 0
< error
| 1064 42000 …near 'DISTRIBUTION JOBS WHERE id > 0' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOB 1
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOB 'x'
< error
| 1064 42000 …near 'DISTRIBUTION JOB 'x'' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOB
< error
| 1064 42000 …near 'DISTRIBUTION JOB' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOB 1 WHERE id > 0
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1 WHERE id > 0' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOBS 1
< error
| 1064 42000 …near 'DISTRIBUTION JOBS 1' at line 1
> exec ^
: string
| SHOW DISTRIBUTION
< error
| 1064 42000 …near 'DISTRIBUTION' at line 1
> exec ^
: string
| SHOW SESSION_STATES
< error
| 1064 42000 …near 'SESSION_STATES' at line 1
> exec ^
: string
| SHOW STATS_EXTENDED
< error
| 1064 42000 …near 'STATS_EXTENDED' at line 1
> exec ^
: string
| SHOW STATS_META
< error
| 1064 42000 …near 'STATS_META' at line 1
> exec ^
: string
| SHOW STATS_META WHERE table_name = 't'
< error
| 1064 42000 …near 'STATS_META WHERE table_name = 't'' at line 1
> exec ^
: string
| SHOW STATS_META LIKE 't'
< error
| 1064 42000 …near 'STATS_META LIKE 't'' at line 1
> exec ^
: string
| SHOW STATS_HISTOGRAMS
< error
| 1064 42000 …near 'STATS_HISTOGRAMS' at line 1
> exec ^
: string
| SHOW STATS_TOPN
< error
| 1064 42000 …near 'STATS_TOPN' at line 1
> exec ^
: string
| SHOW STATS_BUCKETS
< error
| 1064 42000 …near 'STATS_BUCKETS' at line 1
> exec ^
: string
| SHOW STATS_HEALTHY
< error
| 1064 42000 …near 'STATS_HEALTHY' at line 1
> exec ^
: string
| SHOW STATS_LOCKED
< error
| 1064 42000 …near 'STATS_LOCKED' at line 1
> exec ^
: string
| SHOW HISTOGRAMS_IN_FLIGHT
< error
| 1064 42000 …near 'HISTOGRAMS_IN_FLIGHT' at line 1
> exec ^
: string
| SHOW COLUMN_STATS_USAGE
< error
| 1064 42000 …near 'COLUMN_STATS_USAGE' at line 1
> exec ^
: string
| SHOW ANALYZE STATUS
< error
| 1064 42000 …near 'STATUS' at line 1
> exec ^
: string
| SHOW ANALYZE STATUS WHERE table_name = 't'
< error
| 1064 42000 …near 'STATUS WHERE table_name = 't'' at line 1
> exec ^
: string
| SHOW ANALYZE
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| SHOW BUILTINS
< error
| 1064 42000 …near 'BUILTINS' at line 1
> exec ^
: string
| SHOW BUILTINS LIKE 'x'
< error
| 1064 42000 …near 'BUILTINS LIKE 'x'' at line 1
> exec ^
: string
| SHOW BACKUPS
< error
| 1064 42000 …near 'BACKUPS' at line 1
> exec ^
: string
| SHOW BACKUPS WHERE start_time > now()
< error
| 1064 42000 …near 'BACKUPS WHERE start_time > now()' at line 1
> exec ^
: string
| SHOW RESTORES
< error
| 1064 42000 …near 'RESTORES' at line 1
> exec ^
: string
| SHOW RESTORES LIKE 'r0001'
< error
| 1064 42000 …near 'RESTORES LIKE 'r0001'' at line 1
> exec ^
: string
| SHOW CONFIG
< error
| 1064 42000 …near 'CONFIG' at line 1
> exec ^
: string
| SHOW CONFIG WHERE type='tidb'
< error
| 1064 42000 …near 'CONFIG WHERE type='tidb'' at line 1
> exec ^
: string
| SHOW DATABASES TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec ^
: string
| SHOW TABLES TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec ^
: string
| SHOW CREATE TABLE t TABLE t2
< error
| 1064 42000 …near 'TABLE t2' at line 1
> exec ^
: string
| SHOW TABLE STATUS TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec ^
: string
| SHOW VARIABLES LIKE 'x%' TABLE t
< error
| 1064 42000 …near 'TABLE t' at line 1
> exec ^
: string
| SHOW PLACEMENT TABLE t
< error
| 1064 42000 …near 'PLACEMENT TABLE t' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR DATABASE db1 TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR DATABASE db1 TABLE tb1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR SCHEMA db1 TABLE tb1
< error
| 1064 42000 …near 'PLACEMENT FOR SCHEMA db1 TABLE tb1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 TABLE tb2
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 TABLE tb2' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 DATABASE db1
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 DATABASE db1' at line 1
> exec ^
: string
| SHOW PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2
< error
| 1064 42000 …near 'PLACEMENT FOR TABLE tb1 PARTITION p1 PARTITION p2' at line 1
> exec ^
: string
| SHOW TABLE t1 REGIONS TABLE t2
< error
| 1064 42000 …near 't1 REGIONS TABLE t2' at line 1
> exec ^
: string
| SHOW TABLE t1 NEXT_ROW_ID TABLE t2
< error
| 1064 42000 …near 't1 NEXT_ROW_ID TABLE t2' at line 1
> exec ^
: string
| SHOW BUILTINS TABLE t
< error
| 1064 42000 …near 'BUILTINS TABLE t' at line 1
> exec ^
: string
| SHOW DISTRIBUTION JOB 1 TABLE t
< error
| 1064 42000 …near 'DISTRIBUTION JOB 1 TABLE t' at line 1
> exec ^
: string
| SHOW STATS_META TABLE t
< error
| 1064 42000 …near 'STATS_META TABLE t' at line 1
