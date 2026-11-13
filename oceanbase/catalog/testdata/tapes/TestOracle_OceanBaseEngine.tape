-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseEngine default#1
> exec
| DROP DATABASE IF EXISTS obeng_unstated
< ok 0 0
> exec
| CREATE DATABASE obeng_unstated
< ok 1 0
> exec
| USE obeng_unstated
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a))
< ok 0 0
> query
| SHOW CREATE TABLE obeng_unstated.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'
: string
| obeng_unstated
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> exec
| DROP DATABASE IF EXISTS obeng_innodb
< ok 0 0
> exec
| CREATE DATABASE obeng_innodb
< ok 1 0
> exec
| USE obeng_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a)) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE obeng_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'
: string
| obeng_innodb
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> exec
| DROP DATABASE IF EXISTS obeng_myisam
< ok 0 0
> exec
| CREATE DATABASE obeng_myisam
< ok 1 0
> exec
| USE obeng_myisam
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a)) ENGINE=MyISAM
< ok 0 0
> query
| SHOW CREATE TABLE obeng_myisam.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'
: string
| obeng_myisam
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> exec
| DROP DATABASE IF EXISTS obeng_memory
< ok 0 0
> exec
| CREATE DATABASE obeng_memory
< ok 1 0
> exec
| USE obeng_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a)) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE obeng_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'
: string
| obeng_memory
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
> exec
| DROP DATABASE IF EXISTS obeng_nosuch
< ok 0 0
> exec
| CREATE DATABASE obeng_nosuch
< ok 1 0
> exec
| USE obeng_nosuch
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, KEY ka (a)) ENGINE=nosuch
< ok 0 0
> query
| SHOW CREATE TABLE obeng_nosuch.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `ka` (`a`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SELECT ENGINE FROM information_schema.TABLES WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't'
: string
| obeng_nosuch
< result
< column "ENGINE" "VARCHAR"
< row
: bytes
| InnoDB
