-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseAlterAddPrimaryKey default#1
> exec
| DROP DATABASE IF EXISTS obpkvis_0
< ok 0 0
> exec
| CREATE DATABASE obpkvis_0
< ok 1 0
> exec
| USE obpkvis_0
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_0.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_0
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_0.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_0
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_1
< ok 0 0
> exec
| CREATE DATABASE obpkvis_1
< ok 1 0
> exec
| USE obpkvis_1
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_1.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD CONSTRAINT pkc PRIMARY KEY (id) INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_1
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_1.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_1
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_2
< ok 0 0
> exec
| CREATE DATABASE obpkvis_2
< ok 1 0
> exec
| USE obpkvis_2
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_2.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) INVISIBLE VISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_2
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_2.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_2
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_3
< ok 0 0
> exec
| CREATE DATABASE obpkvis_3
< ok 1 0
> exec
| USE obpkvis_3
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_3.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) VISIBLE INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_3
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_3.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_3
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_4
< ok 0 0
> exec
| CREATE DATABASE obpkvis_4
< ok 1 0
> exec
| USE obpkvis_4
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_4.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) GLOBAL INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_4
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_4.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_4
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_5
< ok 0 0
> exec
| CREATE DATABASE obpkvis_5
< ok 1 0
> exec
| USE obpkvis_5
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_5.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) COMMENT 'alter pk comment'
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_5
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_5.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_5
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_6
< ok 0 0
> exec
| CREATE DATABASE obpkvis_6
< ok 1 0
> exec
| USE obpkvis_6
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_6.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD CONSTRAINT PRIMARY KEY (id) COMMENT 'x' INVISIBLE
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_6
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_6.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_6
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_7
< ok 0 0
> exec
| CREATE DATABASE obpkvis_7
< ok 1 0
> exec
| USE obpkvis_7
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_7.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) KEY_BLOCK_SIZE 8 COMMENT 'pkc' USING HASH
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_7
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_7.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_7
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_8
< ok 0 0
> exec
| CREATE DATABASE obpkvis_8
< ok 1 0
> exec
| USE obpkvis_8
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_8.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) STORING (a)
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_8
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_8.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_8
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_9
< ok 0 0
> exec
| CREATE DATABASE obpkvis_9
< ok 1 0
> exec
| USE obpkvis_9
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_9.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD CONSTRAINT pkc PRIMARY KEY (id) COMMENT 'x' STORING (a) STORING (id)
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_9
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_9.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_9
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_10
< ok 0 0
> exec
| CREATE DATABASE obpkvis_10
< ok 1 0
> exec
| USE obpkvis_10
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_10.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) STORING (nosuch)
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 't' AND INDEX_NAME = 'PRIMARY'
: string
| obpkvis_10
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> query
| SHOW CREATE TABLE obpkvis_10.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_10
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_create
< ok 0 0
> exec
| CREATE DATABASE obpkvis_create
< ok 1 0
> exec
| USE obpkvis_create
< ok 0 0
> exec
| CREATE TABLE t (id INT, PRIMARY KEY (id) COMMENT 'pk comment')
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_create.t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`) COMMENT 'pk comment'
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS obpkvis_create
< ok 1 0
> exec
| DROP DATABASE IF EXISTS obpkvis_alter
< ok 0 0
> exec
| CREATE DATABASE obpkvis_alter
< ok 1 0
> exec
| USE obpkvis_alter
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> query
| SHOW CREATE TABLE obpkvis_alter.t
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
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| ALTER TABLE t ALTER INDEX `PRIMARY` INVISIBLE
< error
| 1176 42000 Key 'PRIMARY' doesn't exist in table 't'
> exec
| ALTER TABLE t ALTER INDEX `PRIMARY` VISIBLE
< error
| 1176 42000 Key 'PRIMARY' doesn't exist in table 't'
> exec
| DROP DATABASE IF EXISTS obpkvis_alter
< ok 1 0
