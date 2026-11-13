-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_ForeignKeyDiffNonEmptyCorrect default#1
> exec
| DROP DATABASE IF EXISTS fknec_f_add_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_f_add_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_add_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_add_fk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_f_add_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_t_add_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_t_add_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_add_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_add_fk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_t_add_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_t_add_fk`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_f_drop_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_f_drop_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_drop_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_drop_fk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_f_drop_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_f_drop_fk`.`p`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_t_drop_fk
< ok 0 0
> exec
| CREATE DATABASE fknec_t_drop_fk DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_drop_fk
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_drop_fk.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_t_drop_fk.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_f_modify_action
< ok 0 0
> exec
| CREATE DATABASE fknec_f_modify_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_modify_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_modify_action.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_f_modify_action.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_f_modify_action`.`p`(`id`) ON UPDATE RESTRICT ON DELETE CASCADE 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_t_modify_action
< ok 0 0
> exec
| CREATE DATABASE fknec_t_modify_action DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_modify_action
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_modify_action.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_t_modify_action.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_t_modify_action`.`p`(`id`) ON UPDATE RESTRICT ON DELETE SET NULL 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_f_modify_refcol
< ok 0 0
> exec
| CREATE DATABASE fknec_f_modify_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_f_modify_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_f_modify_refcol.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_f_modify_refcol.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_f_modify_refcol`.`p`(`a`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS fknec_t_modify_refcol
< ok 0 0
> exec
| CREATE DATABASE fknec_t_modify_refcol DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE fknec_t_modify_refcol
< ok 0 0
> exec
| CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))
< ok 0 0
> exec
| CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(b))
< ok 0 0
> query
| SHOW CREATE TABLE fknec_t_modify_refcol.p
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| p
: bytes
| CREATE TABLE `p` (
|   `id` int(11) NOT NULL,
|   `a` int(11) NOT NULL,
|   `b` int(11) NOT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `ua` (`a`) BLOCK_SIZE 16384 LOCAL,
|   UNIQUE KEY `ub` (`b`) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE fknec_t_modify_refcol.c
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| c
: bytes
| CREATE TABLE `c` (
|   `id` int(11) NOT NULL,
|   `pid` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `fknec_t_modify_refcol`.`p`(`b`) ON UPDATE RESTRICT ON DELETE RESTRICT 
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
