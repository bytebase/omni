-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_InnoDBHashRoundTrip default#1
> exec
| DROP DATABASE IF EXISTS hash_innodb_unique_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_unique_hash
< ok 1 0
> exec
| USE hash_innodb_unique_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_unique_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_innodb_key_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_key_hash
< ok 1 0
> exec
| USE hash_innodb_key_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_key_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_innodb_default_engine_hash
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_default_engine_hash
< ok 1 0
> exec
| USE hash_innodb_default_engine_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH)
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_default_engine_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_innodb_key_btree
< ok 0 0
> exec
| CREATE DATABASE hash_innodb_key_btree
< ok 1 0
> exec
| USE hash_innodb_key_btree
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING BTREE) ENGINE=InnoDB
< ok 0 0
> query
| SHOW CREATE TABLE hash_innodb_key_btree.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_memory_unique_hash
< ok 0 0
> exec
| CREATE DATABASE hash_memory_unique_hash
< ok 1 0
> exec
| USE hash_memory_unique_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), UNIQUE KEY uk (name) USING HASH) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE hash_memory_unique_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_memory_key_hash
< ok 0 0
> exec
| CREATE DATABASE hash_memory_key_hash
< ok 1 0
> exec
| USE hash_memory_key_hash
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20), KEY k (name) USING HASH) ENGINE=MEMORY
< ok 0 0
> query
| SHOW CREATE TABLE hash_memory_key_hash.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_alter_add_unique_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_unique_hash_innodb
< ok 1 0
> exec
| USE hash_alter_add_unique_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;ALTER TABLE t ADD UNIQUE uk (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_unique_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_alter_add_unique_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_unique_hash_memory
< ok 1 0
> exec
| USE hash_alter_add_unique_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;ALTER TABLE t ADD UNIQUE uk (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_unique_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_alter_add_key_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_key_hash_innodb
< ok 1 0
> exec
| USE hash_alter_add_key_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;ALTER TABLE t ADD KEY k (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_key_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_alter_add_key_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_alter_add_key_hash_memory
< ok 1 0
> exec
| USE hash_alter_add_key_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;ALTER TABLE t ADD KEY k (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_alter_add_key_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_create_index_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_create_index_hash_innodb
< ok 1 0
> exec
| USE hash_create_index_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;CREATE INDEX k ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_index_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_create_index_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_create_index_hash_memory
< ok 1 0
> exec
| USE hash_create_index_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;CREATE INDEX k ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_index_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `k` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_create_unique_index_hash_innodb
< ok 0 0
> exec
| CREATE DATABASE hash_create_unique_index_hash_innodb
< ok 1 0
> exec
| USE hash_create_unique_index_hash_innodb
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB;CREATE UNIQUE INDEX uk ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_unique_index_hash_innodb.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS hash_create_unique_index_hash_memory
< ok 0 0
> exec
| CREATE DATABASE hash_create_unique_index_hash_memory
< ok 1 0
> exec
| USE hash_create_unique_index_hash_memory
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20)) ENGINE=MEMORY;CREATE UNIQUE INDEX uk ON t (name) USING HASH
< ok 0 0
> query
| SHOW CREATE TABLE hash_create_unique_index_hash_memory.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `name` varchar(20) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `uk` (`name`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
